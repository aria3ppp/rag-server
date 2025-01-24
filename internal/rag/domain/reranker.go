package domain

type RerankerRerankInput struct {
	Query     string
	Documents []string
	TopN      int
}

type RerankerRerankResult struct {
	Index    int
	Document string
	Score    float32
}
