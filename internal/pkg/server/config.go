package server

import "time"

type Config struct {
	GRPCPort                uint16
	HTTPPort                uint16
	GracefulShutdownTimeout time.Duration
}
