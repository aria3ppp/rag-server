package vectorstore

import (
	"context"
	"fmt"
	"log/slog"

	vectorstore_v1 "github.com/aria3ppp/rag-server/gen/go/vectorstore/v1"
	"github.com/aria3ppp/rag-server/internal/rag/config"
	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"
	"github.com/samber/lo"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type vectorstore struct {
	client *grpc.ClientConn
	config *config.VectorStoreConfig
	tracer trace.Tracer
	logger *slog.Logger
}

var _ usecase.VectorStore = (*vectorstore)(nil)

func NewVectorStore(
	ctx context.Context,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
) (*vectorstore, error) {
	client, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", config.VectorStoreConfig.Host, config.VectorStoreConfig.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to grpc.NewClient: %w", err)
	}

	healthResponse, err := grpc_health_v1.NewHealthClient(client).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return nil, err
	}

	if healthResponse.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		return nil, fmt.Errorf("vectorstore health status: %s", healthResponse.GetStatus())
	}

	return &vectorstore{
		client: client,
		config: &config.VectorStoreConfig,
		tracer: tracer,
		logger: logger,
	}, nil
}

func (vs *vectorstore) Search(ctx context.Context, query *domain.VectorStoreSearchInput) ([]*domain.VectorStoreSearchResult, error) {
	filter, err := structpb.NewStruct(query.Filter)
	if err != nil {
		return nil, err
	}

	request := &vectorstore_v1.VectorStoreServiceSearchTextRequest{
		Text:     query.Text,
		TopK:     int64(query.TopK),
		MinScore: query.MinScore,
		Filter:   filter,
	}

	response, err := vectorstore_v1.NewVectorStoreServiceClient(vs.client).SearchText(ctx, request)
	if err != nil {
		return nil, err
	}

	result := lo.Map(response.GetSimilarTexts(), func(similarText *vectorstore_v1.VectorStoreServiceSearchTextResponseSimilarText, _ int) *domain.VectorStoreSearchResult {
		return &domain.VectorStoreSearchResult{
			Text:     similarText.GetText(),
			Score:    similarText.GetScore(),
			Metadata: similarText.GetMetadata().AsMap(),
		}
	})

	return result, nil
}

func (vs *vectorstore) Close() error {
	return vs.client.Close()
}
