package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aria3ppp/rag-server/internal/pkg/prob"
	"github.com/aria3ppp/rag-server/internal/pkg/server"
	vectorstore_app "github.com/aria3ppp/rag-server/internal/vectorstore/app"
	vectorstore_config "github.com/aria3ppp/rag-server/internal/vectorstore/config"
	otel_handler "github.com/aria3ppp/rag-server/pkg/logger/handler/otel"
	stacktrace_handler "github.com/aria3ppp/rag-server/pkg/logger/handler/stacktrace"
	"github.com/aria3ppp/rag-server/pkg/opentelemetry"
	"github.com/aria3ppp/rag-server/pkg/profile"

	"github.com/caarlos0/env/v11"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	var config vectorstore_config.ServerConfig
	if err := env.Parse(&config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse envs: %s\n", err)
		os.Exit(1)
	}

	prob.CheckToRunProbe(
		server.Config{
			GRPCPort: config.GRPCConfig.Port,
			HTTPPort: config.GatewayConfig.Port,
		},
	)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var slogHandler slog.Handler
	if profile.IsDebug {
		slogHandler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		)
	} else {
		slogHandler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			},
		)
	}
	slogHandler = otel_handler.NewHandlerWithTraceInfo(
		slogHandler,
		nil,
	)
	slogHandler = stacktrace_handler.NewStackTraceHandler(
		slogHandler,
		&stacktrace_handler.Opts{
			SkipFrames: 4,
		},
	)
	logger := slog.New(slogHandler)

	otelInitShutdown, err := opentelemetry.InitTracer()
	if err != nil {
		logger.ErrorContext(ctx, "failed to init tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer otelInitShutdown(ctx)

	tracer := otel.Tracer(
		"vectorstore",
		trace.WithInstrumentationVersion(otel.Version()),
	)

	var config vectorstore_config.Config
	if err := env.Parse(&config); err != nil {
		logger.ErrorContext(ctx, "failed to parse env configs", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// create and initialize vectorstore app
	app, err := vectorstore_app.New(
		ctx,
		&config,
		slogHandler,
		tracer,
		http.DefaultClient,
	)
	if err != nil {
		logger.ErrorContext(ctx, "failed to app new", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// start the app and wait for completion or interruption
	if err := app.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.ErrorContext(ctx, "app start failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
