package embedder_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"testing"

	test_server "github.com/aria3ppp/rag-server/internal/pkg/test/server"
	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/infras/embedder"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/trace"
	otel_trace_noop "go.opentelemetry.io/otel/trace/noop"
)

func TestNewEmbedder(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx    context.Context
		config *config.Config
		tracer trace.Tracer
		logger *slog.Logger
	}

	type want struct {
		nilEmbedder bool
		err         bool
	}

	type testCase struct {
		name  string
		input input
		want  want
	}
	testCases := []testCase{
		// {
		// 	name: "invalid embedding size",
		// 	input: input{
		// 		ctx:    context.Background(),
		// 		client: &http.Client{},
		// 		config: &config.Config{},
		// 		tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
		// 		logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
		// 	},
		// 	want: want{
		// 		nilEmbedder: true,
		// 		err:         true,
		// 	},
		// },
		{
			name: "ok",
			input: input{
				ctx:    context.Background(),
				config: &config.Config{},
				tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
				logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			},
			want: want{
				nilEmbedder: false,
				err:         false,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			embedderBaseURL := fmt.Sprintf("http://localhost:%d/v1", embedderPort)
			embedderEmbeddingSize := test_server.GetEmbedderEmbeddingSize(t, embedderBaseURL)

			tt.input.config.EmbedderConfig.BaseURL = embedderBaseURL
			tt.input.config.QdrantConfig.VectorSize = embedderEmbeddingSize
			em, err := embedder.NewEmbedder(
				tt.input.ctx,
				tt.input.config,
				tt.input.tracer,
				tt.input.logger,
				http.DefaultClient,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if (em == nil) != tt.want.nilEmbedder {
				t.Fatal(cmp.Diff(err, nil))
			}
		})
	}
}

func Test_Embedder_Embed(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx   context.Context
		texts []string
	}

	type want struct {
		nilEmbeddings bool
		err           bool
	}

	type testCase struct {
		name   string
		config config.Config
		input  input
		want   want
	}
	testCases := []testCase{
		{
			name:   "failed to llm client create embedding",
			config: config.Config{},
			input: input{
				ctx:   context.Background(),
				texts: []string{},
			},
			want: want{
				nilEmbeddings: true,
				err:           true,
			},
		},
		{
			name:   "ok",
			config: config.Config{},
			input: input{
				ctx:   context.Background(),
				texts: []string{""},
			},
			want: want{
				nilEmbeddings: false,
				err:           false,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			embedderBaseURL := fmt.Sprintf("http://localhost:%d/v1", embedderPort)
			embedderEmbeddingSize := test_server.GetEmbedderEmbeddingSize(t, embedderBaseURL)

			tt.config.EmbedderConfig.BaseURL = embedderBaseURL
			tt.config.QdrantConfig.VectorSize = embedderEmbeddingSize
			em, err := embedder.NewEmbedder(
				ctx,
				&tt.config,
				otel_trace_noop.NewTracerProvider().Tracer(""),
				slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
				http.DefaultClient,
			)

			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			embeddings, err := em.Embed(tt.input.ctx, tt.input.texts)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if (embeddings == nil) != tt.want.nilEmbeddings {
				t.Fatal(cmp.Diff(embeddings == nil, tt.want.nilEmbeddings))
			}
		})
	}
}

var embedderPort int

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
