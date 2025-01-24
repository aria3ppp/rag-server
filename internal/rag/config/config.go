package config

import "time"

type Config struct {
	ServerConfig      ServerConfig
	OpenAIConfig      OpenAIConfig
	RerankerConfig    RerankerConfig
	VectorStoreConfig VectorStoreConfig
}

type ServerConfig struct {
	GRPCConfig              GRPCConfig
	GatewayConfig           GatewayConfig
	GracefulShutdownTimeout time.Duration `env:"RAG_SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"30s"`
}

type GRPCConfig struct {
	Port uint16 `env:"RAG_SERVER_GRPC_PORT" envDefault:"9001"`
}

type GatewayConfig struct {
	Port           uint16   `env:"RAG_SERVER_GATEWAY_PORT" envDefault:"8000"`
	AllowedOrigins []string `env:"RAG_SERVER_GATEWAY_ALLOWED_ORIGINS"`
}

type OpenAIConfig struct {
	BaseURL string `env:"OPENAI_BASEURL,notEmpty"`
	APIKey  string `env:"OPENAI_APIKEY,notEmpty"`
	Model   string `env:"OPENAI_MODEL,notEmpty"`
}

type RerankerConfig struct {
	BaseURL string `env:"RERANKER_BASEURL,notEmpty"`
}

type VectorStoreConfig struct {
	Host     string `env:"VECTORSTORE_HOST,notEmpty"`
	GRPCPort uint16 `env:"VECTORSTORE_SERVER_GRPC_PORT,notEmpty"`
}
