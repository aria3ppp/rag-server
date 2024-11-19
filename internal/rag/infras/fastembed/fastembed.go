package fastembed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/rag/usecase"
)

type FastembedEmbedRequest struct {
	Documents []string `json:"documents"`
}

type FastembedEmbedResponse struct {
	Embeddings    [][]float32 `json:"embeddings"`
	EmbeddingSize int         `json:"embedding_size"`
}

// fastembed implements Embedder
type fastembed struct {
	client  *http.Client
	baseURL string
}

var _ usecase.Embedder = (*fastembed)(nil)

func NewEmbedder(client *http.Client, baseURL string) *fastembed {
	return &fastembed{client: client, baseURL: baseURL}
}

func (e *fastembed) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	embeddingRequest := &FastembedEmbedRequest{
		Documents: texts,
	}

	embeddingRequestJsonBytes, err := json.Marshal(embeddingRequest)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling json: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		e.baseURL,
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

	embedResponse := &FastembedEmbedResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(embedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return embedResponse.Embeddings, nil
}
