package prob

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aria3ppp/rag-server/internal/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func CheckToRunProbe(config server.Config) {
	var (
		probeType string
		mute      bool
	)

	flag.StringVar(&probeType, "probe", "", "probe type (http or grpc: other values skips the prob and run the server)")
	flag.BoolVar(&mute, "mute", false, "mute prob output")
	flag.Parse()

	var (
		probeFunc func(port uint16) (response string, err error)
		port      uint16
	)

	switch strings.ToLower(probeType) {
	case "http":
		probeFunc = runHTTPProbe
		port = config.HTTPPort
	case "grpc":
		probeFunc = runGRPCProbe
		port = config.GRPCPort
	default:
		return
	}

	response, err := probeFunc(port)
	if err != nil {
		if !mute {
			fmt.Fprintf(os.Stderr, "prob failed at %d: %s\n", time.Now().Unix(), err)
		}
		os.Exit(1)
	}

	if !mute {
		fmt.Printf("probe successful at %d: %s\n", time.Now().Unix(), strings.TrimRight(response, "\n"))
	}

	os.Exit(0)
}
func runHTTPProbe(httpPort uint16) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", httpPort))
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%s", responseBody)
	}

	return string(responseBody), nil
}
func runGRPCProbe(grpcPort uint16) (string, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", grpcPort), opts...)
	if err != nil {
		return "", fmt.Errorf("failed to connect to grpc server: %w", err)
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)

	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return "", fmt.Errorf("grpc check failed: %w", err)
	}

	resp.ProtoReflect().Descriptor()

	responseJSON, err := protojson.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal json: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return "", fmt.Errorf("%s", responseJSON)
	}

	return string(responseJSON), nil
}
