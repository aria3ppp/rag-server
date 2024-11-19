package fastembed

type fastembedEmbedRequestDTO struct {
	Documents []string `json:"documents"`
}

type fastembedEmbedResponseDTO struct {
	Embeddings    [][]float32 `json:"embeddings"`
	EmbeddingSize int         `json:"embedding_size"`
}

type fastembedHealthCheckResponseDTO struct {
	Status string `json:"status"`
}
