package server_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/aria3ppp/rag-server/internal/pkg/server"
	test_port "github.com/aria3ppp/rag-server/internal/pkg/test/port"
	"github.com/aria3ppp/rag-server/pkg/wait"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestServerStart(t *testing.T) {
	type deps struct {
		config         server.Config
		logger         *slog.Logger
		grpcClientConn *grpc.ClientConn
		grpcServer     *grpc.Server
		httpServer     *http.Server
	}

	type want struct {
		err bool
	}

	tests := []struct {
		name                string
		setup               func(t *testing.T) (ctx context.Context, ctxCancel func(), deps deps, close func())
		ensureServerRunning func(t *testing.T, config *server.Config)
		want                want
	}{
		{
			name: "ok",
			setup: func(t *testing.T) (context.Context, func(), deps, func()) {
				closeList := make([]func(), 0)

				config := server.Config{
					GRPCPort:                uint16(test_port.GetFreePort(t)),
					HTTPPort:                uint16(test_port.GetFreePort(t)),
					GracefulShutdownTimeout: 30 * time.Second,
				}
				grpcServer := grpc.NewServer()
				healthServer := health.NewServer()
				grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

				httpServer := &http.Server{
					Addr: fmt.Sprintf(":%d", config.HTTPPort),
				}

				grpcClientConn, err := grpc.NewClient(
					fmt.Sprintf(":%d", config.GRPCPort),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}

				deps := deps{
					config:         config,
					logger:         slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
					grpcClientConn: grpcClientConn,
					grpcServer:     grpcServer,
					httpServer:     httpServer,
				}

				close := func() {
					for _, c := range slices.Backward(closeList) {
						c()
					}
				}

				ctx, ctxCancel := context.WithCancel(context.Background())

				return ctx, ctxCancel, deps, close
			},
			ensureServerRunning: func(t *testing.T, config *server.Config) {
				grpcClientConn, err := grpc.NewClient(
					fmt.Sprintf(":%d", config.GRPCPort),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}

				healthClient := grpc_health_v1.NewHealthClient(grpcClientConn)

				wait.Until(t, &wait.Opts{Interval: 500 * time.Millisecond, MaxRetries: 20}, func() error {
					resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
					if err != nil {
						return err
					}

					if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
						return fmt.Errorf("health status: %s", resp.Status)
					}

					return nil
				})
			},
			want: want{err: false},
		},
		{
			name: "fail to start gRPC server",
			setup: func(t *testing.T) (context.Context, func(), deps, func()) {
				closeList := make([]func(), 0)

				config := server.Config{
					GRPCPort:                uint16(test_port.GetFreePort(t)),
					HTTPPort:                uint16(test_port.GetFreePort(t)),
					GracefulShutdownTimeout: 30 * time.Second,
				}

				// listen to the grpc port
				listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}
				closeList = append(closeList, func() { listener.Close() })

				grpcServer := grpc.NewServer()
				httpServer := &http.Server{
					Addr: fmt.Sprintf(":%d", config.HTTPPort),
				}

				grpcClientConn, err := grpc.NewClient(
					fmt.Sprintf(":%d", config.GRPCPort),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}

				deps := deps{
					config:         config,
					logger:         slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
					grpcClientConn: grpcClientConn,
					grpcServer:     grpcServer,
					httpServer:     httpServer,
				}

				close := func() {
					for _, c := range slices.Backward(closeList) {
						c()
					}
				}

				ctx, ctxCancel := context.Background(), func() {}

				return ctx, ctxCancel, deps, close
			},
			ensureServerRunning: func(t *testing.T, config *server.Config) {},
			want:                want{err: true},
		},
		{
			name: "fail to start HTTP server",
			setup: func(t *testing.T) (context.Context, func(), deps, func()) {
				closeList := make([]func(), 0)

				config := server.Config{
					GRPCPort:                uint16(test_port.GetFreePort(t)),
					HTTPPort:                uint16(test_port.GetFreePort(t)),
					GracefulShutdownTimeout: 30 * time.Second,
				}

				// listen to the http port
				listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.HTTPPort))
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}
				closeList = append(closeList, func() { listener.Close() })

				grpcServer := grpc.NewServer()
				httpServer := &http.Server{
					Addr: fmt.Sprintf(":%d", config.HTTPPort),
				}

				grpcClientConn, err := grpc.NewClient(
					fmt.Sprintf(":%d", config.GRPCPort),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					t.Fatal(cmp.Diff(err, nil))
				}

				deps := deps{
					config:         config,
					logger:         slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
					grpcClientConn: grpcClientConn,
					grpcServer:     grpcServer,
					httpServer:     httpServer,
				}

				close := func() {
					for _, c := range slices.Backward(closeList) {
						c()
					}
				}

				ctx, ctxCancel := context.Background(), func() {}

				return ctx, ctxCancel, deps, close
			},
			ensureServerRunning: func(t *testing.T, config *server.Config) {},
			want:                want{err: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, ctxCancel, deps, setupClose := tt.setup(t)
			t.Cleanup(setupClose)

			srv := server.New(
				deps.config,
				deps.logger,
				deps.grpcClientConn,
				deps.grpcServer,
				deps.httpServer,
			)

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()

				err := srv.Start(ctx)
				if (err != nil) != tt.want.err {
					t.Error(cmp.Diff(err, nil))
				}
			}()

			tt.ensureServerRunning(t, &deps.config)
			ctxCancel()

			wg.Wait()
		})
	}
}
