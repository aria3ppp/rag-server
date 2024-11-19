package domain

type VectorStoreSearchInput struct {
	Vector []float32
	TopK   int
	Filter map[string]any
}

type VectorStoreSearchResult struct {
	ID       string
	Score    float32
	Vector   []float32
	Metadata map[string]any
}
