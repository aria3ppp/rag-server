package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tmc/langchaingo/llms/openai"
)

func SetupEmbedderServer(f Fatalizer) (httpPort int, cleanup func()) {
	return setupLlamaCppServer(f, "embedder", []string{"-m", "/models/CompendiumLabs/bge-small-en-v1.5-gguf/bge-small-en-v1.5-f32.gguf", "--embedding"})
}

func setupLlamaCppServer(f Fatalizer, modelsType string, command []string) (httpPort int, cleanup func()) {
	f.Helper()

	const (
		internalHTTPPort = "8080/tcp"
	)

	_, currentFileDir, _, _ := runtime.Caller(0)
	currentFileDir = filepath.Dir(currentFileDir)

	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "ghcr.io/ggerganov/llama.cpp:server",
				Cmd:          command,
				Env:          map[string]string{"LLAMA_ARG_PORT": internalHTTPPort},
				ExposedPorts: []string{internalHTTPPort},
				WaitingFor: wait.NewHTTPStrategy("/health").
					WithPort(internalHTTPPort).
					WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK }).
					WithResponseMatcher(func(body io.Reader) bool {
						bytes, err := io.ReadAll(body)
						if err != nil {
							f.Fatal(cmp.Diff(err, nil))
							return false
						}
						return string(bytes) == `{"status":"ok"}`
					}),
				HostConfigModifier: func(hc *container.HostConfig) {
					// hc.Binds = []string{filepath.Join(currentFileDir, "..", "..", "..", "..", "models") + ":/models"}
					hc.Binds = []string{filepath.Join(currentFileDir, "..", "..", "..", "..", "models") + ":/models"}
				},
			},
			Started: true,
		},
	)
	if err != nil {
		f.Fatal(cmp.Diff(err, nil))
	}

	cleanup = func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			f.Fatal(cmp.Diff(err, nil))
		}
	}

	containerJSON, err := container.Inspect(context.Background())
	if err != nil {
		f.Fatal(cmp.Diff(err, nil))
	}

	httpPort = lo.Must(strconv.Atoi(containerJSON.NetworkSettings.Ports[internalHTTPPort][0].HostPort))
	fmt.Printf("%s: [%s] http port: %d\n", modelsType, f.Name(), httpPort)

	return httpPort, cleanup
}

func GetEmbedderEmbeddingSize(tb testing.TB, baseURL string) int {
	tb.Helper()

	llm, err := openai.New(
		openai.WithBaseURL(baseURL),
		openai.WithToken("OPENAI_API_KEY"),
	)
	if err != nil {
		tb.Fatal(cmp.Diff(err, nil))
	}

	embeddings, err := llm.CreateEmbedding(context.Background(), []string{""})
	if err != nil {
		tb.Fatal(cmp.Diff(err, nil))
	}

	return len(embeddings[0])
}
