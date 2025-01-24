package config

import "time"

type Config struct {
	ServerConfig   ServerConfig
	EmbedderConfig EmbedderConfig
	QdrantConfig   QdrantConfig
}

type ServerConfig struct {
	GRPCConfig              GRPCConfig
	GatewayConfig           GatewayConfig
	GracefulShutdownTimeout time.Duration `env:"VECTORSTORE_SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"30s"`
}

type GRPCConfig struct {
	Port uint16 `env:"VECTORSTORE_SERVER_GRPC_PORT" envDefault:"9091"`
}

type GatewayConfig struct {
	Port           uint16   `env:"VECTORSTORE_SERVER_GATEWAY_PORT" envDefault:"8080"`
	AllowedOrigins []string `env:"VECTORSTORE_SERVER_GATEWAY_ALLOWED_ORIGINS"`
}

type EmbedderConfig struct {
	BaseURL string `env:"EMBEDDER_BASEURL,notEmpty"`
}

type QdrantConfig struct {
	Host           string `env:"QDRANT_HOST,notEmpty"`
	GRPCPort       uint16 `env:"QDRANT_GRPC_PORT,notEmpty"`
	CollectionName string `env:"QDRANT_COLLECTION_NAME,notEmpty"`
	VectorSize     int    `env:"QDRANT_VECTOR_SIZE,notEmpty"`
}
