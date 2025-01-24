package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	test_server "github.com/aria3ppp/rag-server/internal/pkg/test/server"
	"github.com/google/go-cmp/cmp"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/openai"
	"google.golang.org/protobuf/encoding/protojson"
)

func _TestSetupQdrantServer(t *testing.T) {
	// setup qdrant server
	grpcPort, cleanup := test_server.SetupQdrantServer(t)
	t.Cleanup(cleanup)

	// create qdrant client
	client, err := qdrant.NewClient(&qdrant.Config{
		Port: grpcPort,
	})
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}

	// create the collection
	cn := "cn"
	if err := client.CreateCollection(
		context.Background(),
		&qdrant.CreateCollection{
			CollectionName: cn,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1,
				Distance: qdrant.Distance_Cosine,
			}),
		},
	); err != nil {
		t.Fatal(err)
	}

	t0 := time.Now()

	if _, err := client.Upsert(
		context.Background(),
		&qdrant.UpsertPoints{
			// Wait:           lo.ToPtr(true),
			CollectionName: cn,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewIDNum(0),
					Vectors: qdrant.NewVectors(1),
					// Payload: qdrant.NewValueMap(map[string]any{"k": "v"}),
				},
			},
		},
	); err != nil {
		t.Fatal(err)
	}

	// compute the upsert runtime duration
	d := time.Since(t0)
	fmt.Println("upsert took:", d.String())

	// fetch the point
	ps, err := client.Get(
		context.Background(),
		&qdrant.GetPoints{
			CollectionName: cn,
			Ids:            []*qdrant.PointId{qdrant.NewIDNum(0)},
			WithPayload:    qdrant.NewWithPayload(true),
			WithVectors:    qdrant.NewWithVectors(true),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// print the point
	for _, p := range ps {
		bytes, err := protojson.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(bytes))
	}
}

func _TestSetupEmbedderServer(t *testing.T) {
	// setup embedder server
	httpPort, cleanup := test_server.SetupEmbedderServer(t)
	t.Cleanup(cleanup)

	llm, err := openai.New(
		openai.WithBaseURL(fmt.Sprintf("http://localhost:%d/v1", httpPort)),
		openai.WithToken("OPENAI_API_KEY"),
	)
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}

	texts := []string{
		"hello",
		"world",
	}

	embeddings, err := llm.CreateEmbedding(
		context.Background(),
		texts,
	)
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}

	for i, e := range embeddings {
		fmt.Println(texts[i], ":", e)
	}
}
