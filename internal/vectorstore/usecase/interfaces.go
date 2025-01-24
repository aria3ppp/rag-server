package usecase

//go:generate mockgen -destination=mocks/mocks.go -package=mocks -typed . Embedder,IDGenerator,VectorRepo,UseCase

import (
	"context"

	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
)

type (
	Embedder interface {
		Embed(ctx context.Context, texts []string) ([][]float32, error)
	}

	IDGenerator interface {
		NewID() (string, error)
	}

	VectorRepo interface {
		Insert(ctx context.Context, embeddings []*domain.VectorRepoInsertEmbedding) error
		Query(ctx context.Context, query *domain.VectorRepoQueryInput) ([]*domain.VectorRepoQueryResult, error)
	}

	UseCase interface {
		InsertTexts(ctx context.Context, input *domain.InsertTextsInput) error
		SearchText(ctx context.Context, input *domain.SearchTextInput) (*domain.SearchTextResult, error)
	}
)
