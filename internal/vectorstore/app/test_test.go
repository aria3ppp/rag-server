package app_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	test_port "github.com/aria3ppp/rag-server/internal/pkg/test/port"
	test_server "github.com/aria3ppp/rag-server/internal/pkg/test/server"
	"github.com/aria3ppp/rag-server/internal/vectorstore/app"
	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/pkg/wait"

	goccy_json "github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
	otel_trace_noop "go.opentelemetry.io/otel/trace/noop"
)

const qdrantServer = "qdrant_server"

var embedderPort int

func setupVectorStoreApp(t *testing.T, waitOpts *wait.Opts) (grpcGatewayPort uint16) {
	t.Helper()
	ctx := context.Background()

	ports, cleanup := test_server.SetupServers(t, map[string]test_server.TestServerFunc{
		qdrantServer: test_server.SetupQdrantServer,
	})
	t.Cleanup(cleanup)

	embedderBaseURL := fmt.Sprintf("http://localhost:%d/v1", embedderPort)
	embedderEmbeddingSize := test_server.GetEmbedderEmbeddingSize(t, embedderBaseURL)

	config := &config.Config{
		ServerConfig: config.ServerConfig{
			GRPCConfig: config.GRPCConfig{
				Port: uint16(test_port.GetFreePort(t)),
			},
			GatewayConfig: config.GatewayConfig{
				Port:           uint16(test_port.GetFreePort(t)),
				AllowedOrigins: []string{},
			},
			GracefulShutdownTimeout: 30 * time.Second,
		},
		EmbedderConfig: config.EmbedderConfig{
			BaseURL: embedderBaseURL,
		},
		QdrantConfig: config.QdrantConfig{
			Host:           "localhost",
			CollectionName: "collection",
			GRPCPort:       uint16(ports[qdrantServer]),
			VectorSize:     embedderEmbeddingSize,
		},
	}

	slogHandler := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	tracer := otel_trace_noop.NewTracerProvider().Tracer("")

	app, err := app.New(
		ctx,
		config,
		slogHandler,
		tracer,
		http.DefaultClient,
	)
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}

	var appWG sync.WaitGroup
	appWG.Add(1)
	go func() {
		defer appWG.Done()

		if err := app.Start(ctx); err != nil {
			t.Fatal(cmp.Diff(err, nil))
		}
	}()
	t.Cleanup(func() {
		defer appWG.Wait()

		if err := app.Shutdown(context.Background()); err != nil {
			t.Fatal(cmp.Diff(err, nil))
		}
	})

	if dt, ok := t.Deadline(); ok {
		d := time.Until(dt)
		if d > waitOpts.Interval {
			waitOpts.MaxRetries = int(int64(d) / int64(waitOpts.Interval))
		}
	}

	wait.Until(t, waitOpts, func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", int(config.ServerConfig.GatewayConfig.Port)))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var health struct {
			Status string `json:"status"`
		}
		if err := goccy_json.Unmarshal(body, &health); err != nil {
			return err
		}

		if health.Status != "SERVING" {
			return fmt.Errorf("health status: %s", health.Status)
		}

		return nil
	})

	return config.ServerConfig.GatewayConfig.Port
}

func TestMain(m *testing.M) {
	const embedderTestcontainer = "embedder_testcontainer"

	ports, cleanup := test_server.SetupServers(
		test_server.NewFatalizer(runtime.FuncForPC(func() uintptr { pc, _, _, _ := runtime.Caller(1); return pc }()).Name()),
		map[string]test_server.TestServerFunc{
			embedderTestcontainer: test_server.SetupEmbedderServer,
		},
	)

	embedderPort = ports[embedderTestcontainer]

	exitCode := m.Run()
	cleanup()

	os.Exit(exitCode)
}
