package usecase

import (
	"context"
	"time"

	"github.com/aria3ppp/rag-server/internal/rag/domain"
)

type (
	Reranker interface {
		Rerank(ctx context.Context, input *domain.RerankerRerankInput) ([]*domain.RerankerRerankResult, error)
	}

	LLM interface {
		StreamCompletion(ctx context.Context, chat []*domain.Message, completionHandler func(completionChunk string, err error) (continueRunning bool))
	}

	VectorStore interface {
		Search(ctx context.Context, query *domain.VectorStoreSearchInput) ([]*domain.VectorStoreSearchResult, error)
	}

	Clock interface {
		TimeNow() time.Time
	}

	UseCase interface {
		QueryStream(ctx context.Context, input *domain.QueryStreamInput, handler func(event *domain.QueryStreamResultEvent) (continueRunning bool))
		Query(ctx context.Context, input *domain.QueryInput) (*domain.QueryResult, error)
	}
)
