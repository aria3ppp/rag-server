package usecase

import (
	"context"
	"time"

	"github.com/aria3ppp/rag-server/internal/rag/domain"
)

type (
	Embedder interface {
		Embed(ctx context.Context, texts []string) ([][]float32, error)
	}

	// TODO: set system and user prompt messages
	AIModel interface {
		StreamQuery(ctx context.Context, prompt *domain.AIModelStreamQueryInput, handler func(ctx context.Context, chunk []byte, done bool) error) error
	}

	VectorStore interface {
		Search(ctx context.Context, query *domain.VectorStoreSearchInput) ([]*domain.VectorStoreSearchResult, error)
	}

	Clock interface {
		TimeNow() time.Time
	}

	UseCase interface {
		QuerySync(ctx context.Context, query string) (string, error)
		QueryAsync(ctx context.Context, query string, handler func(ctx context.Context, event *domain.QueryAsyncEvent) error) error
	}
)
