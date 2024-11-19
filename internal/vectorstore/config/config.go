package config

type Config struct {
	GRPCServerConfig  GRPCServerConfig
	GRPCGatewayConfig GRPCGatewayConfig
	FastembedConfig   FastembedConfig
	QdrantConfig      QdrantConfig
}

type GRPCServerConfig struct {
	Port uint16 `env:"VECTORSTORE_GRPC_SERVER_PORT" envDefault:"9091"`
}

type GRPCGatewayConfig struct {
	Port uint16 `env:"VECTORSTORE_GRPC_GATEWAY_PORT" envDefault:"8080"`
}

type FastembedConfig struct {
	BaseURL string `env:"FASTEMBED_BASEURL,notEmpty"`
}

type QdrantConfig struct {
	Host           string `env:"QDRANT_HOST,notEmpty"`
	GRPCPort       uint16 `env:"QDRANT_GRPC_PORT,notEmpty"`
	CollectionName string `env:"QDRANT_COLLECTION_NAME,notEmpty"`
	VectorSize     int    `env:"QDRANT_VECTOR_SIZE,notEmpty"`
}
