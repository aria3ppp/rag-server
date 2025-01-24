package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	vectorstorev1 "github.com/aria3ppp/rag-server/gen/go/vectorstore/v1"
	vectorstore_openapiv2 "github.com/aria3ppp/rag-server/gen/openapiv2/vectorstore"
	"github.com/aria3ppp/rag-server/internal/pkg/server"
	vectorstore_grpc_server "github.com/aria3ppp/rag-server/internal/vectorstore/app/grpc_server"
	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/embedder"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/qdrant"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/uuid"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
	template_app "github.com/aria3ppp/rag-server/pkg/app"
	"github.com/aria3ppp/rag-server/pkg/profile"

	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func New(
	ctx context.Context,
	config *config.Config,
	slogHandler slog.Handler,
	tracer trace.Tracer,
	httpClient *http.Client,
) (*template_app.App, error) {
	logger := slog.New(slogHandler)

	embedder, err := embedder.NewEmbedder(
		ctx,
		config,
		tracer,
		logger,
		httpClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to embedder.NewEmbedder: %w", err)
	}

	idGenerator := uuid.NewIDGenerator()

	vectorRepo, err := qdrant.NewVectorRepo(
		ctx,
		config,
		tracer,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to qdrant.NewVectorRepo: %w", err)
	}

	useCase := usecase.NewUseCase(
		embedder,
		idGenerator,
		vectorRepo,
		config,
		tracer,
		logger,
	)

	vectorStoreGRPCServer := vectorstore_grpc_server.NewGRPCServer(
		useCase,
		tracer,
		logger,
	)

	healthServer := health.NewServer()
	// healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	grpcServer := grpc.NewServer()

	vectorstorev1.RegisterVectorStoreServiceServer(grpcServer, vectorStoreGRPCServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	if profile.IsDebug {
		reflection.Register(grpcServer)
	}

	grpcClientConn, err := grpc.NewClient(
		fmt.Sprintf(":%d", config.ServerConfig.GRPCConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to grpc.NewClient: %w", err)
	}

	mux := grpc_gateway_runtime.NewServeMux(
		grpc_gateway_runtime.WithHealthEndpointAt(
			grpc_health_v1.NewHealthClient(grpcClientConn),
			"/healthz",
		),
	)
	mux.HandlePath(http.MethodGet, "/{version}/{file}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFileFS(w, r, vectorstore_openapiv2.EmbeddedFS, filepath.Join(pathParams["version"], pathParams["file"]))
	})

	if err := vectorstorev1.RegisterVectorStoreServiceHandler(ctx, mux, grpcClientConn); err != nil {
		return nil, fmt.Errorf("failed to vectorstorev1.RegisterVectorStoreServiceHandler: %w", err)
	}

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   config.ServerConfig.GatewayConfig.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler(mux)

	// create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerConfig.GatewayConfig.Port),
		Handler: corsHandler,
	}

	server := server.New(
		server.Config{
			GRPCPort:                config.ServerConfig.GRPCConfig.Port,
			HTTPPort:                config.ServerConfig.GatewayConfig.Port,
			GracefulShutdownTimeout: config.ServerConfig.GracefulShutdownTimeout,
		},
		logger,
		grpcClientConn,
		grpcServer,
		httpServer,
	)

	return template_app.New(server.Start, logger), nil
}
