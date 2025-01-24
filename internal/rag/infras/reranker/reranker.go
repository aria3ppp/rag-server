package reranker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/aria3ppp/rag-server/internal/rag/config"
	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"
	goccy_json "github.com/goccy/go-json"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"
)

type reranker struct {
	httpClient *http.Client
	config     *config.RerankerConfig
	tracer     trace.Tracer
	logger     *slog.Logger
}

var _ usecase.Reranker = (*reranker)(nil)

func NewReranker(
	ctx context.Context,
	config *config.Config,
	tracer trace.Tracer,
	logger *slog.Logger,
	httpClient *http.Client,
) (*reranker, error) {
	reqBodyBytes, err := json.Marshal(&rerankerRerankRequest{
		Query:     "",
		TopN:      1,
		Documents: []string{""},
	})
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Post(fmt.Sprintf("%s/rerank", config.RerankerConfig.BaseURL), "application/json", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respBody string
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			respBody = ""
		}

		respBody = string(respBodyBytes)
		return nil, fmt.Errorf("reranker got status code %d: %s", resp.StatusCode, respBody)
	}

	return &reranker{
		httpClient: httpClient,
		config:     &config.RerankerConfig,
		tracer:     tracer,
		logger:     logger,
	}, nil
}

func (r *reranker) Rerank(ctx context.Context, input *domain.RerankerRerankInput) (_ []*domain.RerankerRerankResult, err error) {
	var reqBodyBytes []byte
	reqBodyBytes, err = json.Marshal(&rerankerRerankRequest{
		Query:     input.Query,
		TopN:      input.TopN,
		Documents: input.Documents,
	})
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/rerank", r.config.BaseURL),
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return nil, err
	}

	httpResponse, err := r.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reranker got status code %d: %s", httpResponse.StatusCode, respBodyBytes)
	}

	var response rerankerRerankResponse
	if err = goccy_json.Unmarshal(respBodyBytes, &response); err != nil {
		return nil, err
	}

	results := lo.Map(response.Results, func(item *rerankerRerankResponseResult, _ int) *domain.RerankerRerankResult {
		return &domain.RerankerRerankResult{
			Index:    item.Index,
			Document: input.Documents[item.Index],
			Score:    item.RelevanceScore,
		}
	})

	return results, nil
}
