package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aria3ppp/rag-server/internal/rag/config"
	"github.com/aria3ppp/rag-server/internal/rag/domain"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type usecase struct {
	vectorStore VectorStore
	reranker    Reranker
	llm         LLM
	clock       Clock
	config      *config.Config
	tracer      trace.Tracer
	logger      *slog.Logger
}

var _ UseCase = (*usecase)(nil)

func NewUseCase(
	vectorStore VectorStore,
	reranker Reranker,
	llm LLM,
	clock Clock,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
) *usecase {
	return &usecase{
		vectorStore: vectorStore,
		reranker:    reranker,
		llm:         llm,
		clock:       clock,
		config:      config,
		tracer:      tracer,
		logger:      logger,
	}
}

func (uc *usecase) Query(ctx context.Context, input *domain.QueryInput) (_ *domain.QueryResult, err error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.Query")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var (
		completion strings.Builder
		t0         *int64
		tEnd       int64
	)

	streamInput := &domain.QueryStreamInput{
		Query:    input.Query,
		Messages: input.Messages,
	}

	uc.QueryStream(ctx, streamInput, func(event *domain.QueryStreamResultEvent) (continueRunning bool) {
		if _, err = completion.WriteString(event.Content); err != nil {
			return false
		}

		if event.Error != nil {
			err = event.Error
			return false
		}

		if t0 == nil {
			t0 = &event.CreatedAtMS
		}
		tEnd = event.CreatedAtMS

		return true
	})
	if err != nil {
		return nil, err
	}

	return &domain.QueryResult{
		Content:     completion.String(),
		CreatedInMS: (tEnd - *t0),
	}, nil
}

func (uc *usecase) QueryStream(ctx context.Context, input *domain.QueryStreamInput, handler func(event *domain.QueryStreamResultEvent) (continueRunning bool)) {
	var err error

	ctx, span := uc.tracer.Start(ctx, "usecase.QueryStream")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())

			handler(&domain.QueryStreamResultEvent{
				Content:     "",
				CreatedAtMS: uc.clock.TimeNow().UnixMilli(),
				StopReason:  domain.StopReasonError,
				Error:       err,
			})
		}
	}()

	//
	// validate input
	//
	if err = input.Validate(ctx); err != nil {
		return
	}

	//
	// search vector store with top_k 10
	//

	vectorStoreSearchInput := &domain.VectorStoreSearchInput{
		Text:     input.Query,
		TopK:     5,   // TODO: make it configurable?
		MinScore: 0.4, // TODO: make it configurable?
		Filter:   map[string]any{},
	}

	var vectorStoreSearchResults []*domain.VectorStoreSearchResult
	vectorStoreSearchResults, err = uc.vectorStore.Search(ctx, vectorStoreSearchInput)
	if err != nil {
		return
	}

	var retrievedDocument string

	if len(vectorStoreSearchResults) == 1 {
		retrievedDocument = vectorStoreSearchResults[0].Text
	} else if len(vectorStoreSearchResults) > 1 {
		//
		// rerank search results
		//

		rerankInput := &domain.RerankerRerankInput{
			Query: input.Query,
			Documents: lo.Map(vectorStoreSearchResults, func(r *domain.VectorStoreSearchResult, _ int) string {
				return r.Text
			}),
			TopN: 1, // TODO: make it configurable?
		}

		var rerankResult []*domain.RerankerRerankResult
		rerankResult, err = uc.reranker.Rerank(ctx, rerankInput)
		if err != nil {
			return
		}

		maxRerankResult := lo.MaxBy(rerankResult, func(a *domain.RerankerRerankResult, b *domain.RerankerRerankResult) bool { return a.Score > b.Score })
		retrievedDocument = lo.IfF(maxRerankResult != nil, func() string { return maxRerankResult.Document }).ElseF(func() string { return "" })
	}

	//
	// prompt llm with retrieved document
	//

	chat := append(
		input.Messages,
		&domain.Message{
			Role:    domain.RoleAssistant,
			Content: retrievedDocument,
		},
		&domain.Message{
			Role:    domain.RoleUser,
			Content: input.Query,
		},
	)

	uc.llm.StreamCompletion(ctx, chat, func(completionChunk string, handlerErr error) (continueRunning bool) {
		err = handlerErr

		if err != nil {
			return false
		}

		return handler(&domain.QueryStreamResultEvent{
			Content:     completionChunk,
			CreatedAtMS: uc.clock.TimeNow().UnixMilli(),
			StopReason:  domain.StopReasonUnspecified,
			Error:       nil,
		})

	})

	handler(&domain.QueryStreamResultEvent{
		Content:     "",
		CreatedAtMS: uc.clock.TimeNow().UnixMilli(),
		StopReason:  domain.StopReasonDone,
		Error:       nil,
	})

	return
}
