package openai

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/rag/config"
	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type openaiLLM struct {
	client *openai.Client
	config *config.OpenAIConfig
	tracer trace.Tracer
	logger *slog.Logger
}

var _ usecase.LLM = (*openaiLLM)(nil)

func NewLLM(
	ctx context.Context,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
	httpClient *http.Client,
) (*openaiLLM, error) {
	client := openai.NewClient(
		option.WithBaseURL(config.OpenAIConfig.BaseURL),
		option.WithAPIKey(config.OpenAIConfig.APIKey),
		option.WithHTTPClient(httpClient),
	)

	if _, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.F(openai.EmbeddingNewParamsInputUnion(openai.EmbeddingNewParamsInputArrayOfStrings{""})),
	}); err != nil {
		return nil, err
	}

	return &openaiLLM{
		client: client,
		config: &config.OpenAIConfig,
		tracer: tracer,
		logger: logger,
	}, nil
}

func (llm *openaiLLM) StreamCompletion(ctx context.Context, chat []*domain.Message, completionHandler func(completionChunk string, err error) (continueRunning bool)) {
	var err error

	ctx, span := llm.tracer.Start(ctx, "openaiLLM.StreamCompletion")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(chat))
	for _, m := range chat {
		switch m.Role {
		case domain.RoleSystem:
			messages = append(messages, openai.SystemMessage(m.Content))
		case domain.RoleAssistant:
			messages = append(messages, openai.AssistantMessage(m.Content))
		case domain.RoleUser:
			messages = append(messages, openai.UserMessage(m.Content))
		}
	}

	stream := llm.client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Seed:     openai.Int(0),
		Model:    openai.F(llm.config.Model),
	})
	defer stream.Close()

	for stream.Next() {
		// Cancel the stream on ctx.Done
		select {
		case <-ctx.Done():
			err = ctx.Err()
			completionHandler("", err)
			return
		default:
		}

		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			if continueRunning := completionHandler(chunk.Choices[0].Delta.Content, nil); !continueRunning {
				return
			}
		}
	}

	return
}
