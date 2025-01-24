package reranker

type rerankerRerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	TopN      int      `json:"top_n"`
	Documents []string `json:"documents"`
}

type rerankerRerankResponse struct {
	Results []*rerankerRerankResponseResult `json:"results"`
}

type rerankerRerankResponseResult struct {
	Index          int     `json:"index"`
	RelevanceScore float32 `json:"relevance_score"`
}
