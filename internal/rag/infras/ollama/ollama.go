package ollama

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// ollamaQueryer implements AIQueryer for Ollama models
type ollamaQueryer struct {
	ollamaLLM *ollama.LLM
	model     string
}

func NewQueryer(ollamaLLM *ollama.LLM, model string) *ollamaQueryer {
	return &ollamaQueryer{ollamaLLM: ollamaLLM, model: model}
}

func (o *ollamaQueryer) StreamQuery(ctx context.Context, prompt string, handler func(ctx context.Context, chunk []byte, done bool) error) error {
	_, err := o.ollamaLLM.Call(ctx, prompt,
		llms.WithModel(o.model),
		llms.WithTemperature(0.7),
		llms.WithTopP(0.9),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			return handler(ctx, chunk, false)
		}),
	)
	if err != nil {
		return err
	}

	if err := handler(ctx, nil, true); err != nil {
		return err
	}

	return nil
}
