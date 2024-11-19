package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	vectorstorev1 "github.com/aria3ppp/rag-server/gen/go/vectorstore/v1"
	"github.com/aria3ppp/rag-server/gen/openapiv2/vectorstore"
	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/fastembed"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/qdrant"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/uuid"
	"github.com/aria3ppp/rag-server/internal/vectorstore/router"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
	"github.com/aria3ppp/rag-server/pkg/logger"
	"github.com/aria3ppp/rag-server/pkg/profile"

	"github.com/caarlos0/env/v11"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	logger := logger.NewLogger(os.Stdout)

	var config config.Config
	err := env.Parse(&config)
	if err != nil {
		logger.Error("failed to parse env configs", slog.Any("error", err))
		os.Exit(1)
	}

	embedder, err := fastembed.NewEmbedder(
		ctx,
		http.DefaultClient,
		config.FastembedConfig.BaseURL,
	)
	if err != nil {
		logger.Error("failed to fastembed.NewEmbedder", slog.Any("error", err))
		os.Exit(1)
	}

	idGenerator := uuid.NewIDGenerator()

	vectorRepo, err := qdrant.NewVectorRepo(
		ctx,
		config.QdrantConfig.Host,
		int(config.QdrantConfig.GRPCPort),
		config.QdrantConfig.Collection,
	)
	if err != nil {
		logger.Error("failed to qdrant.NewVectorRepo", slog.Any("error", err))
		os.Exit(1)
	}

	useCase := usecase.NewUseCase(
		embedder,
		idGenerator,
		vectorRepo,
	)

	vectorStoreGRPCService := router.NewVectorStoreGRPCService(useCase)

	healthServer := health.NewServer()
	// healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	grpcServer := grpc.NewServer()

	vectorstorev1.RegisterVectorStoreServiceServer(grpcServer, vectorStoreGRPCService)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	if profile.IsDebug {
		reflection.Register(grpcServer)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcClientConn, err := grpc.NewClient(fmt.Sprintf(":%d", config.GRPCServerConfig.Port), opts...)
	if err != nil {
		logger.Error("failed to grpc.NewClient", slog.Any("error", err))
		os.Exit(1)
	}

	defer func() {
		go func() {
			<-ctx.Done()
			if err := grpcClientConn.Close(); err != nil {
				logger.Error("failed to grpcClientConn.Close", slog.Any("error", err))
			}
		}()
	}()

	mux := runtime.NewServeMux(
		runtime.WithHealthEndpointAt(
			grpc_health_v1.NewHealthClient(grpcClientConn),
			"/healthz",
		),
	)
	mux.HandlePath(http.MethodGet, "/{version}/openapiv2/{file}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		version := pathParams["version"]
		file := pathParams["file"]
		http.ServeFileFS(w, r, vectorstore.EmbeddedFS, filepath.Join(version, file))
	})

	if err := vectorstorev1.RegisterVectorStoreServiceHandler(ctx, mux, grpcClientConn); err != nil {
		logger.Error("failed to vectorstorev1.RegisterVectorStoreServiceHandler", slog.Any("error", err))
		os.Exit(1)
	}

	// Start gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCServerConfig.Port))
	if err != nil {
		logger.Error("failed to grpc listen", slog.Uint64("port", uint64(config.GRPCServerConfig.Port)), slog.Any("error", err))
		os.Exit(1)
	}

	// Start servers
	go func() {
		logger.Info("starting grpc server", slog.Uint64("port", uint64(config.GRPCServerConfig.Port)))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("failed to serve grpc", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("starting http server", slog.Uint64("port", uint64(config.GRPCGatewayConfig.Port)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.GRPCGatewayConfig.Port), mux); err != nil {
		logger.Error("failed to serve http", slog.Any("error", err))
		os.Exit(1)
	}

}
