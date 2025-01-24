package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type usecase struct {
	embedder    Embedder
	vectorRepo  VectorRepo
	idGenerator IDGenerator
	config      *config.Config
	tracer      trace.Tracer
	logger      *slog.Logger
}

var _ UseCase = (*usecase)(nil)

func NewUseCase(
	embedder Embedder,
	idGenerator IDGenerator,
	vectorRepo VectorRepo,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
) *usecase {
	return &usecase{
		embedder:    embedder,
		idGenerator: idGenerator,
		vectorRepo:  vectorRepo,
		config:      config,
		tracer:      tracer,
		logger:      logger,
	}
}

func (uc *usecase) InsertTexts(ctx context.Context, input *domain.InsertTextsInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.InsertTexts")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if err := input.Validate(ctx); err != nil {
		uc.logger.ErrorContext(ctx, "failed to validate input", slog.String("error", err.Error()))
		return err
	}

	textsString := lo.Map(input.Texts, func(item *domain.InsertTextsInputText, _ int) string { return item.Text })

	embeddings, err := uc.embedder.Embed(ctx, textsString)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to embed text", slog.String("error", err.Error()))
		return err
	}

	if len(embeddings) != len(input.Texts) {
		uc.logger.ErrorContext(ctx, "invalid embeddings length", slog.Int("texts length", len(input.Texts)), slog.Int("embeddings length", len(embeddings)))
		return fmt.Errorf("invalid embeddings length: texts length = %d, embeddings length = %d", len(input.Texts), len(embeddings))
	}

	vectorRepoInsertEmbeddings := make([]*domain.VectorRepoInsertEmbedding, 0, len(input.Texts))

	for index, text := range input.Texts {
		id, err := uc.idGenerator.NewID()
		if err != nil {
			uc.logger.ErrorContext(ctx, "failed to generate new id", slog.String("error", err.Error()))
			return err
		}

		metadata := lo.Assign(
			text.Metadata,
			map[string]any{"text": text.Text},
		)

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
		uc.logger.ErrorContext(ctx, "failed to repo insert", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (uc *usecase) SearchText(ctx context.Context, input *domain.SearchTextInput) (_ *domain.SearchTextResult, err error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SearchText")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if err := input.Validate(ctx); err != nil {
		uc.logger.ErrorContext(ctx, "failed to validate input", slog.String("error", err.Error()))
		return nil, err
	}

	vectors, err := uc.embedder.Embed(ctx, []string{input.Text})
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to embed text", slog.String("error", err.Error()))
		return nil, err
	}

	if len(vectors) != 1 {
		uc.logger.ErrorContext(ctx, "invalid vectors length", slog.Int("must", 1), slog.Int("got", len(vectors)))
		return nil, fmt.Errorf("invalid vectors length: vector length must be 1 got %d", len(vectors))
	}

	vectorRepoQueryInput := &domain.VectorRepoQueryInput{
		Vector:   vectors[0],
		TopK:     input.TopK,
		MinScore: input.MinScore,
		Filter:   input.Filter,
	}

	queryResults, err := uc.vectorRepo.Query(ctx, vectorRepoQueryInput)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to repo query", slog.String("error", err.Error()))
		return nil, err
	}

	if len(queryResults) > input.TopK {
		uc.logger.ErrorContext(ctx, "query results length couldn't exceed topk", slog.Int("query result length", len(queryResults)), slog.Int("topk", input.TopK))
		return nil, fmt.Errorf("query results length couldn't exceed topk: query result length = %d, topk = %d", len(queryResults), input.TopK)
	}

	similarTexts := make([]*domain.SearchTextResultItem, 0, len(queryResults))

	for _, qr := range queryResults {
		textAny, exists := qr.Metadata["text"] // SAFETY: querying a nil map won't panic
		if !exists {
			uc.logger.ErrorContext(ctx, "metadata field text not exist for the record", slog.String("record id", qr.ID))
			return nil, fmt.Errorf("metadata field text not exist for record with id %s", qr.ID)
		}

		text, assertionOk := textAny.(string)
		if !assertionOk {
			uc.logger.ErrorContext(ctx, "metadata field text is not a string", slog.String("got type", fmt.Sprintf("%T", textAny)))
			return nil, fmt.Errorf("metadata field text is not a string: got type %T", textAny)
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
