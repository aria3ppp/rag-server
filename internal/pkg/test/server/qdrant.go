package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupQdrantServer(f Fatalizer) (grpcPort int, cleanup func()) {
	f.Helper()

	const (
		internalGRPCPort = "6334/tcp"
		internalHTTPPort = "6333/tcp"
	)

	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "qdrant/qdrant:v1.12.4",
				ExposedPorts: []string{internalGRPCPort, internalHTTPPort},
				WaitingFor: wait.NewHTTPStrategy("/healthz").
					WithPort(internalHTTPPort).
					WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK }),
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

	grpcPort = lo.Must(strconv.Atoi(containerJSON.NetworkSettings.Ports[internalGRPCPort][0].HostPort))
	httpPort := lo.Must(strconv.Atoi(containerJSON.NetworkSettings.Ports[internalHTTPPort][0].HostPort))

	fmt.Printf("%s: [qdrant] grpc port: %d, http port: %d\n", f.Name(), grpcPort, httpPort)

	return grpcPort, cleanup
}
