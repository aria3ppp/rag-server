package qdrant

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"

	"github.com/qdrant/go-client/qdrant"
)

type qdrantRepo struct {
	client         *qdrant.Client
	collectionName string
}

var _ usecase.VectorRepo = (*qdrantRepo)(nil)

func NewVectorRepo(ctx context.Context, host string, port int, collection string) (*qdrantRepo, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}

	return &qdrantRepo{client: client, collectionName: collection}, nil
}

func (repo *qdrantRepo) Insert(ctx context.Context, embeddings []*domain.VectorRepoInsertEmbedding) error {
	points := make([]*qdrant.PointStruct, len(embeddings))
	for i, emb := range embeddings {
		payload, err := convertToQdrantMap(emb.Metadata)
		if err != nil {
			return err
		}

		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewID(emb.ID),
			Vectors: qdrant.NewVectors(emb.Vector...),
			Payload: payload,
		}
	}

	_, err := repo.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: repo.collectionName,
		Points:         points,
	})

	if err != nil {
		return fmt.Errorf("failed to insert embeddings: %v", err)
	}

	return nil
}

func convertToQdrantValue(v interface{}) (*qdrant.Value, error) {
	if v == nil {
		return &qdrant.Value{Kind: &qdrant.Value_NullValue{}}, nil
	}

	switch value := v.(type) {
	case bool:
		return &qdrant.Value{Kind: &qdrant.Value_BoolValue{BoolValue: value}}, nil
	case int:
		return &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(value)}}, nil
	case int64:
		return &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: value}}, nil
	case float32:
		return &qdrant.Value{Kind: &qdrant.Value_DoubleValue{DoubleValue: float64(value)}}, nil
	case float64:
		return &qdrant.Value{Kind: &qdrant.Value_DoubleValue{DoubleValue: value}}, nil
	case string:
		return &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: value}}, nil
	case []interface{}:
		listValue := &qdrant.ListValue{Values: make([]*qdrant.Value, len(value))}
		for i, item := range value {
			convertedItem, err := convertToQdrantValue(item)
			if err != nil {
				return nil, fmt.Errorf("error converting list item at index %d: %w", i, err)
			}
			listValue.Values[i] = convertedItem
		}
		return &qdrant.Value{Kind: &qdrant.Value_ListValue{ListValue: listValue}}, nil
	case map[string]interface{}:
		structValue := &qdrant.Struct{Fields: make(map[string]*qdrant.Value)}
		for k, v := range value {
			convertedValue, err := convertToQdrantValue(v)
			if err != nil {
				return nil, fmt.Errorf("error converting struct field '%s': %w", k, err)
			}
			structValue.Fields[k] = convertedValue
		}
		return &qdrant.Value{Kind: &qdrant.Value_StructValue{StructValue: structValue}}, nil
	default:
		// Handle other slice types
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice {
			listValue := &qdrant.ListValue{Values: make([]*qdrant.Value, rv.Len())}
			for i := 0; i < rv.Len(); i++ {
				convertedItem, err := convertToQdrantValue(rv.Index(i).Interface())
				if err != nil {
					return nil, fmt.Errorf("error converting slice item at index %d: %w", i, err)
				}
				listValue.Values[i] = convertedItem
			}
			return &qdrant.Value{Kind: &qdrant.Value_ListValue{ListValue: listValue}}, nil
		}
		// For unsupported types, return an error
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func convertToQdrantMap(input map[string]any) (map[string]*qdrant.Value, error) {
	result := make(map[string]*qdrant.Value, len(input))
	for key, value := range input {
		convertedValue, err := convertToQdrantValue(value)
		if err != nil {
			return nil, fmt.Errorf("error converting field '%s': %w", key, err)
		}
		result[key] = convertedValue
	}
	return result, nil
}

func (repo *qdrantRepo) Query(ctx context.Context, query *domain.VectorRepoQueryInput) ([]*domain.VectorRepoQueryResult, error) {
	searchParams := &qdrant.QueryPoints{
		CollectionName: repo.collectionName,
		Query:          qdrant.NewQueryDense(query.Vector),
		Limit:          qdrant.PtrOf(uint64(query.TopK)),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	}

	response, err := repo.client.Query(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("qdrant search failed: %v", err)
	}

	results := make([]*domain.VectorRepoQueryResult, 0, len(response))
	for _, point := range response {
		queryResult := &domain.VectorRepoQueryResult{}
		queryResult.ID = point.Id.String()
		queryResult.Score = point.GetScore()
		queryResult.Vector = point.GetVectors().GetVector().GetData()
		queryResult.Metadata, err = convertFromQdrantMap(point.Payload)
		if err != nil {
			return nil, err
		}

		results = append(results, queryResult)
	}

	return results, nil
}

func convertFromQdrantValue(value *qdrant.Value) (any, error) {
	if value == nil {
		return nil, nil
	}

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
		list := make([]interface{}, len(v.ListValue.Values))
		for i, item := range v.ListValue.Values {
			mapValue, err := convertFromQdrantValue(item)
			if err != nil {
				return nil, err
			}
			list[i] = mapValue
		}
		return list, nil
	case *qdrant.Value_StructValue:
		m := make(map[string]interface{})
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
