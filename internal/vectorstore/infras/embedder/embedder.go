package embedder

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
	"github.com/tmc/langchaingo/llms/openai"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type embedder struct {
	llmClient *openai.LLM
	config    *config.EmbedderConfig
	tracer    trace.Tracer
	logger    *slog.Logger
}

var _ usecase.Embedder = (*embedder)(nil)

func NewEmbedder(
	ctx context.Context,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
	httpClient *http.Client,
) (*embedder, error) {
	llmClient, err := openai.New(
		openai.WithBaseURL(config.EmbedderConfig.BaseURL),
		openai.WithToken("OPENAI_API_KEY"),
		openai.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	embeddings, err := llmClient.CreateEmbedding(ctx, []string{""})
	if err != nil {
		return nil, err
	}

	if embeddingSize := len(embeddings[0]); embeddingSize != config.QdrantConfig.VectorSize {
		return nil, fmt.Errorf("invalid embedding size: got=%d, want=%d", embeddingSize, config.QdrantConfig.VectorSize)
	}

	return &embedder{
		llmClient: llmClient,
		config:    &config.EmbedderConfig,
		tracer:    tracer,
		logger:    logger,
	}, nil
}

func (e *embedder) Embed(ctx context.Context, texts []string) (_ [][]float32, err error) {
	ctx, span := e.tracer.Start(ctx, "embedder.Embed")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	embeddings, err := e.llmClient.CreateEmbedding(ctx, texts)
	if err != nil {
		e.logger.ErrorContext(ctx, "failed to llm client create embedding", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to llm client create embedding: %w", err)
	}

	return embeddings, nil
}
