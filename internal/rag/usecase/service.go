package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aria3ppp/rag-server/internal/rag/domain"
)

type usecase struct {
	embedder    Embedder
	aiModel     AIModel
	vectorStore VectorStore
	clock       Clock
}

var _ UseCase = (*usecase)(nil)

func NewUseCase(
	embedder Embedder,
	aiModel AIModel,
	vectorStore VectorStore,
	clock Clock,
) *usecase {
	return &usecase{
		embedder:    embedder,
		aiModel:     aiModel,
		vectorStore: vectorStore,
		clock:       clock,
	}
}

func (uc *usecase) QuerySync(ctx context.Context, query string) (string, error) {
	queryContext, err := uc.createQueryContext(ctx, query)
	if err != nil {
		return "", err
	}

	aiModelQueryInput := &domain.AIModelStreamQueryInput{
		Prompt:  query,
		Context: queryContext,
	}

	var completion strings.Builder
	err = uc.aiModel.StreamQuery(ctx, aiModelQueryInput, func(ctx context.Context, chunk []byte, done bool) error {
		_, err := completion.Write(chunk)
		return err
	})
	if err != nil {
		return "", err
	}

	return completion.String(), nil
}

func (uc *usecase) QueryAsync(ctx context.Context, query string, handler func(ctx context.Context, event *domain.QueryAsyncEvent) error) error {
	queryContext, err := uc.createQueryContext(ctx, query)
	if err != nil {
		return err
	}

	aiModelStreamQueryInput := &domain.AIModelStreamQueryInput{
		Prompt:  query,
		Context: queryContext,
	}

	if err := uc.aiModel.StreamQuery(ctx, aiModelStreamQueryInput, func(ctx context.Context, chunk []byte, done bool) error {
		var queryAsyncEvent *domain.QueryAsyncEvent

		if done {
			queryAsyncEvent = &domain.QueryAsyncEvent{
				Content:   "",
				CreatedAt: uc.clock.TimeNow(),
				Done:      true,
			}
		} else {
			queryAsyncEvent = &domain.QueryAsyncEvent{
				Content:   string(chunk),
				CreatedAt: uc.clock.TimeNow(),
				Done:      done,
			}
		}

		if err := handler(ctx, queryAsyncEvent); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (uc *usecase) createQueryContext(ctx context.Context, query string) (string, error) {
	embeddings, err := uc.embedder.Embed(ctx, []string{query})
	if err != nil {
		return "", err
	}

	if len(embeddings) != 1 {
		return "", fmt.Errorf("unexpected number of embeddings: expected 1 got %d", len(embeddings))
	}

	vectorStoreSearchInput := &domain.VectorStoreSearchInput{
		Vector: embeddings[0],
		TopK:   1,
		Filter: map[string]any{},
	}

	vectorStoreSearchResults, err := uc.vectorStore.Search(ctx, vectorStoreSearchInput)
	if err != nil {
		return "", err
	}

	var queryContext string

	if len(vectorStoreSearchResults) == 0 {
		queryContext = ""
	} else {
		content, ok := vectorStoreSearchResults[1].Metadata["content"]
		if !ok {
			return "", errors.New("vector store result metadata have no field content")
		}

		queryContext, ok = content.(string)
		if !ok {
			return "", errors.New("vector store result content is not string")
		}
	}

	return queryContext, nil
}
