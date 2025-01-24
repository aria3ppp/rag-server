package domain

type VectorRepoInsertEmbedding struct {
	ID       string
	Vector   []float32
	Metadata map[string]any
}

type VectorRepoQueryInput struct {
	Vector   []float32
	TopK     int
	MinScore float32
	Filter   map[string]any
}

type VectorRepoQueryResult struct {
	ID       string
	Score    float32
	Vector   []float32
	Metadata map[string]any
}
