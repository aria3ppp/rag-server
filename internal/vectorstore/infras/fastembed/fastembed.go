package fastembed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
)

type fastembed struct {
	client *http.Client
	config *config.FastembedConfig
	logger *slog.Logger
}

var _ usecase.Embedder = (*fastembed)(nil)

func NewEmbedder(ctx context.Context, client *http.Client, config *config.Config, logger *slog.Logger) (*fastembed, error) {
	fastembed := &fastembed{client: client, config: &config.FastembedConfig, logger: logger}

	if err := fastembed.healthcheck(ctx, config); err != nil {
		return nil, err
	}

	return fastembed, nil
}

func (e *fastembed) healthcheck(ctx context.Context, config *config.Config) error {
	// httpRequest, err := http.NewRequestWithContext(
	// 	ctx,
	// 	http.MethodGet,
	// 	fmt.Sprintf("%s/healthcheck", e.config.BaseURL),
	// 	nil,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to create http request: %w", err)
	// }

	// httpResponse, err := e.client.Do(httpRequest)
	// if err != nil {
	// 	return fmt.Errorf("failed to do http request: %w", err)
	// }

	// defer httpResponse.Body.Close()

	// healthCheckResponse := &fastembedHealthCheckResponseDTO{}
	// if err := json.NewDecoder(httpResponse.Body).Decode(healthCheckResponse); err != nil {
	// 	return fmt.Errorf("failed to decode response body: %w", err)
	// }

	// if healthCheckResponse.Status != "OK" {
	// 	return fmt.Errorf("fastembed status not ok: %s", healthCheckResponse.Status)
	// }

	// return nil

	embeddingRequest := &fastembedEmbedRequestDTO{
		Documents: []string{""},
	}

	embeddingRequestJsonBytes, err := json.Marshal(embeddingRequest)
	if err != nil {
		return fmt.Errorf("failed marshaling json: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/embed_text", e.config.BaseURL),
		bytes.NewReader(embeddingRequestJsonBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	httpResponse, err := e.client.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("failed to do http request: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("embedder service returned status code: %d", httpResponse.StatusCode)
	}

	embedResponse := &fastembedEmbedResponseDTO{}
	if err := json.NewDecoder(httpResponse.Body).Decode(embedResponse); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	if embedResponse.EmbeddingSize != config.QdrantConfig.VectorSize {
		return fmt.Errorf("invalid embedding size: got=%d, want=%d", embedResponse.EmbeddingSize, config.QdrantConfig.VectorSize)
	}

	return nil
}

func (e *fastembed) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	embeddingRequest := &fastembedEmbedRequestDTO{
		Documents: texts,
	}

	embeddingRequestJsonBytes, err := json.Marshal(embeddingRequest)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling json: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/embed_text", e.config.BaseURL),
		bytes.NewReader(embeddingRequestJsonBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpResponse, err := e.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedder service returned status code: %d", httpResponse.StatusCode)
	}

	embedResponse := &fastembedEmbedResponseDTO{}
	if err := json.NewDecoder(httpResponse.Body).Decode(embedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return embedResponse.Embeddings, nil
}
