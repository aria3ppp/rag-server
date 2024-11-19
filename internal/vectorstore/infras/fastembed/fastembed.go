package fastembed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
)

type fastembed struct {
	client  *http.Client
	baseURL string
}

var _ usecase.Embedder = (*fastembed)(nil)

func NewEmbedder(ctx context.Context, client *http.Client, baseURL string) (*fastembed, error) {
	fastembed := &fastembed{client: client, baseURL: baseURL}

	status, err := fastembed.healthcheck(ctx)
	if err != nil {
		return nil, err
	}

	if status != "OK" {
		return nil, fmt.Errorf("fastembed status not ok: %s", status)
	}

	return fastembed, nil
}

func (e *fastembed) healthcheck(ctx context.Context) (status string, err error) {
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/healthcheck", e.baseURL),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	httpResponse, err := e.client.Do(httpRequest)
	if err != nil {
		return "", fmt.Errorf("failed to do http request: %w", err)
	}

	defer httpResponse.Body.Close()

	healthCheckResponse := &fastembedHealthCheckResponseDTO{}
	if err := json.NewDecoder(httpResponse.Body).Decode(healthCheckResponse); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	return healthCheckResponse.Status, nil
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
		fmt.Sprintf("%s/embed-text", e.baseURL),
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

	embedResponse := &fastembedEmbedResponseDTO{}
	if err := json.NewDecoder(httpResponse.Body).Decode(embedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return embedResponse.Embeddings, nil
}
