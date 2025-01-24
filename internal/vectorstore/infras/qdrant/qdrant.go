package qdrant

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"

	"github.com/qdrant/go-client/qdrant"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpc_codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

type qdrantRepo struct {
	client *qdrant.Client
	config *config.QdrantConfig
	tracer trace.Tracer
	logger *slog.Logger
}

var _ usecase.VectorRepo = (*qdrantRepo)(nil)

func NewVectorRepo(
	ctx context.Context,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
) (*qdrantRepo, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: config.QdrantConfig.Host,
		Port: int(config.QdrantConfig.GRPCPort),
	})
	if err != nil {
		return nil, err
	}

	_, err = client.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := client.GetCollectionInfo(ctx, config.QdrantConfig.CollectionName); err != nil {
		if grpc_status.Code(err) != grpc_codes.NotFound {
			return nil, fmt.Errorf("failed to get collection info: %w", err)
		}

		createCollection := &qdrant.CreateCollection{
			CollectionName: config.QdrantConfig.CollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     uint64(config.QdrantConfig.VectorSize),
				Distance: qdrant.Distance_Cosine,
			}),
		}

		if err := client.CreateCollection(ctx, createCollection); err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	return &qdrantRepo{
		client: client,
		config: &config.QdrantConfig,
		tracer: tracer,
		logger: logger,
	}, nil
}

func (repo *qdrantRepo) Insert(ctx context.Context, embeddings []*domain.VectorRepoInsertEmbedding) (err error) {
	ctx, span := repo.tracer.Start(ctx, "qdrantRepo.Insert")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	points := make([]*qdrant.PointStruct, 0, len(embeddings))
	for _, embedding := range embeddings {
		payload, err := qdrant.TryValueMap(embedding.Metadata)
		if err != nil {
			repo.logger.ErrorContext(ctx, "failed to convert to qdrant map", slog.String("error", err.Error()))
			return err
		}

		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewID(embedding.ID),
			Vectors: qdrant.NewVectors(embedding.Vector...),
			Payload: payload,
		})
	}

	_, err = repo.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: repo.config.CollectionName,
		Points:         points,
	})
	if err != nil {
		repo.logger.ErrorContext(ctx, "failed to qdrant client upsert", slog.String("error", err.Error()))
		return fmt.Errorf("failed to qdrant client upsert: %v", err)
	}

	return nil
}

func (repo *qdrantRepo) Query(ctx context.Context, query *domain.VectorRepoQueryInput) (_ []*domain.VectorRepoQueryResult, err error) {
	ctx, span := repo.tracer.Start(ctx, "qdrantRepo.Query")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	searchParams := &qdrant.QueryPoints{
		CollectionName: repo.config.CollectionName,
		Query:          qdrant.NewQueryDense(query.Vector),
		Limit:          qdrant.PtrOf(uint64(query.TopK)),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
		Filter:         nil, // TODO: convert query.Filter to *qdrant.Filter
		ScoreThreshold: &query.MinScore,
	}

	response, err := repo.client.Query(ctx, searchParams)
	if err != nil {
		repo.logger.ErrorContext(ctx, "failed to qdrant client query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to qdrant client query: %v", err)
	}

	results := make([]*domain.VectorRepoQueryResult, 0, len(response))
	for _, point := range response {
		metadata, err := convertFromQdrantMap(point.Payload)
		if err != nil {
			repo.logger.ErrorContext(ctx, "failed to convert from qdrant map", slog.String("error", err.Error()))
			return nil, err
		}

		results = append(results, &domain.VectorRepoQueryResult{
			ID:       point.Id.GetUuid(),
			Score:    point.GetScore(),
			Vector:   point.GetVectors().GetVector().GetData(),
			Metadata: metadata,
		})
	}

	return results, nil
}

func convertFromQdrantMap(input map[string]*qdrant.Value) (map[string]any, error) {
	result := make(map[string]any, len(input))
	for key, value := range input {
		mapValue, err := convertFromQdrantValue(value)
		if err != nil {
			return nil, err
		}
		result[key] = mapValue
	}
	return result, nil
}

func convertFromQdrantValue(value *qdrant.Value) (any, error) {
	switch v := value.GetKind().(type) {
	case *qdrant.Value_NullValue:
		return nil, nil
	case *qdrant.Value_BoolValue:
		return v.BoolValue, nil
	case *qdrant.Value_IntegerValue:
		return v.IntegerValue, nil
	case *qdrant.Value_DoubleValue:
		return v.DoubleValue, nil
	case *qdrant.Value_StringValue:
		return v.StringValue, nil
	case *qdrant.Value_ListValue:
		list := make([]any, len(v.ListValue.Values))
		for i, item := range v.ListValue.Values {
			mapValue, err := convertFromQdrantValue(item)
			if err != nil {
				return nil, err
			}
			list[i] = mapValue
		}
		return list, nil
	case *qdrant.Value_StructValue:
		m := make(map[string]any)
		for k, val := range v.StructValue.Fields {
			mapValue, err := convertFromQdrantValue(val)
			if err != nil {
				return nil, err
			}
			m[k] = mapValue
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
