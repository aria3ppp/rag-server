package qdrant

import (
	"context"
	"fmt"

	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"

	"github.com/qdrant/go-client/qdrant"
)

// qdrantStore implements VectorStore for Qdrant
type qdrantStore struct {
	client         *qdrant.Client
	collectionName string
}

var _ usecase.VectorStore = (*qdrantStore)(nil)

func NewQdrantStore(host string, port int, collection string) (*qdrantStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		return nil, err
	}

	return &qdrantStore{client: client, collectionName: collection}, nil
}

func (q *qdrantStore) Search(ctx context.Context, query *domain.VectorStoreSearchInput) ([]*domain.VectorStoreSearchResult, error) {
	searchParams := &qdrant.QueryPoints{
		CollectionName: q.collectionName,
		Query:          qdrant.NewQueryDense(query.Vector),
		Limit:          qdrant.PtrOf(uint64(query.TopK)),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	}

	response, err := q.client.Query(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("qdrant search failed: %v", err)
	}

	results := make([]*domain.VectorStoreSearchResult, 0, len(response))
	for _, point := range response {
		queryResult := &domain.VectorStoreSearchResult{}
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
