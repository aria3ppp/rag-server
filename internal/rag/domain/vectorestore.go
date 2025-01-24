package domain

type VectorStoreSearchInput struct {
	Text     string
	TopK     int
	MinScore float32
	Filter   map[string]any
}

type VectorStoreSearchResult struct {
	Text     string
	Score    float32
	Metadata map[string]any
}
