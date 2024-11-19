package usecase

import (
	"context"
	"fmt"

	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"

	"github.com/samber/lo"
)

type usecase struct {
	embedder    Embedder
	vectorRepo  VectorRepo
	idGenerator IDGenerator
}

var _ UseCase = (*usecase)(nil)

func NewUseCase(
	embedder Embedder,
	idGenerator IDGenerator,
	vectorRepo VectorRepo,
) *usecase {
	return &usecase{
		embedder:    embedder,
		idGenerator: idGenerator,
		vectorRepo:  vectorRepo,
	}
}

func (uc *usecase) InsertTexts(ctx context.Context, input *domain.InsertTextsInput) error {
	if err := input.Validate(ctx); err != nil {
		return err
	}

	textsString := lo.Map(input.Texts, func(item *domain.InsertTextsInputText, _ int) string { return item.Text })

	embeddings, err := uc.embedder.Embed(ctx, textsString)
	if err != nil {
		return err
	}

	if len(embeddings) != len(input.Texts) {
		return fmt.Errorf("invalid embeddings length: texts length = %d, embeddings length = %d", len(input.Texts), len(embeddings))
	}

	vectorRepoInsertEmbeddings := make([]*domain.VectorRepoInsertEmbedding, 0, len(input.Texts))

	for index, text := range input.Texts {
		id, err := uc.idGenerator.NewID()
		if err != nil {
			return err
		}

		metadata := text.Metadata
		metadata["text"] = text.Text

		vectorRepoInsertEmbeddings = append(
			vectorRepoInsertEmbeddings,
			&domain.VectorRepoInsertEmbedding{
				ID:       id,
				Vector:   embeddings[index],
				Metadata: metadata,
			},
		)
	}

	if err := uc.vectorRepo.Insert(ctx, vectorRepoInsertEmbeddings); err != nil {
		return err
	}

	return nil
}

func (uc *usecase) SearchText(ctx context.Context, input *domain.SearchTextInput) (*domain.SearchTextResult, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, err
	}

	vectors, err := uc.embedder.Embed(ctx, []string{input.Text})
	if err != nil {
		return nil, err
	}

	if len(vectors) != 1 {
		return nil, fmt.Errorf("invalid vectors length: vector length must be 1 got %d", len(vectors))
	}

	vectorRepoQueryInput := &domain.VectorRepoQueryInput{
		Vector: vectors[0],
		TopK:   input.TopK,
		Filter: input.Filter,
	}

	queryResults, err := uc.vectorRepo.Query(ctx, vectorRepoQueryInput)
	if err != nil {
		return nil, err
	}

	if len(queryResults) > input.TopK {
		return nil, fmt.Errorf("query results length couldn't exceed topk: query result length = %d, topk = %d", len(queryResults), input.TopK)
	}

	similarTexts := make([]*domain.SearchTextResultItem, 0, len(queryResults))

	for _, qr := range queryResults {
		textAny, exists := qr.Metadata["text"]
		if !exists {
			return nil, fmt.Errorf("metadata field text not exist for record with id %s", qr.ID)
		}

		text, assertionOk := textAny.(string)
		if !assertionOk {
			return nil, fmt.Errorf("metadata field text is not a string")
		}

		delete(qr.Metadata, "text")

		similarTexts = append(
			similarTexts,
			&domain.SearchTextResultItem{
				Text:     text,
				Score:    qr.Score,
				Metadata: qr.Metadata,
			},
		)
	}

	searchTextResults := &domain.SearchTextResult{
		SimilarTexts: similarTexts,
	}

	return searchTextResults, nil
}
