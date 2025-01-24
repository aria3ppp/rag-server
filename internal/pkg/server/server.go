package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type Server struct {
	config         Config
	logger         *slog.Logger
	grpcClientConn *grpc.ClientConn
	grpcServer     *grpc.Server
	httpServer     *http.Server
}

func New(
	config Config,
	logger *slog.Logger,
	grpcClientConn *grpc.ClientConn,
	grpcServer *grpc.Server,
	httpServer *http.Server,
) *Server {
	return &Server{
		config:         config,
		logger:         logger,
		grpcClientConn: grpcClientConn,
		grpcServer:     grpcServer,
		httpServer:     httpServer,
	}
}

func (s *Server) Start(ctx context.Context) (err error) {
	defer func() {
		// Close the gRPC client connection during shutdown
		if grpcClientConnCloseErr := s.grpcClientConn.Close(); grpcClientConnCloseErr != nil {
			err = fmt.Errorf("gRPC client connection close error: %w", grpcClientConnCloseErr)
		}
	}()

	// create gRPC listener
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to net.listen on port %d: %w", s.config.GRPCPort, err)
	}

	// error channel to collect errors from both servers
	errChan := make(chan error, 2)

	// start gRPC server in a goroutine
	go func() {
		s.logger.InfoContext(ctx, "starting gRPC server", slog.Uint64("port", uint64(s.config.GRPCPort)))
		if err := s.grpcServer.Serve(grpcListener); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// start HTTP server in a goroutine
	go func() {
		s.logger.InfoContext(ctx, "starting HTTP server", slog.Uint64("port", uint64(s.config.HTTPPort)))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// wait for any error or context cancellation
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), s.config.GracefulShutdownTimeout)
		defer cancel()

		if err := s.shutdown(ctx); err != nil {
			return err
		}
		return ctx.Err()
	}
}

func (s *Server) shutdown(ctx context.Context) error {
	errChan := make(chan error, 2)

	// shutdown HTTP server
	go func() {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("HTTP server shutdown error: %w", err)
			return
		}
		errChan <- nil
	}()

	// shutdown gRPC server
	go func() {
		stopped := make(chan struct{})
		go func() {
			s.grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-ctx.Done():
			s.grpcServer.Stop()
			errChan <- fmt.Errorf("gRPC server shutdown timeout: %w", ctx.Err())
		case <-stopped:
			errChan <- nil
		}
	}()

	var shutdownErr error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			shutdownErr = errors.Join(shutdownErr, err)
		}
	}

	return shutdownErr
}
