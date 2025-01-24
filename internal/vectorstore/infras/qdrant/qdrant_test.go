package qdrant_test

import (
	"context"
	"io"
	"log/slog"
	"math"
	"testing"

	test_server "github.com/aria3ppp/rag-server/internal/pkg/test/server"
	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	qdrant_infras "github.com/aria3ppp/rag-server/internal/vectorstore/infras/qdrant"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"go.opentelemetry.io/otel/trace"
	otel_trace_noop "go.opentelemetry.io/otel/trace/noop"
)

func TestNewVectorRepo(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx    context.Context
		config *config.Config
		tracer trace.Tracer
		logger *slog.Logger
	}

	type want struct {
		nilRepo bool
		err     bool
	}

	type testCase struct {
		name     string
		clientFn func(*qdrant.Client) error
		input    input
		want     want
	}
	testCases := []testCase{
		func() testCase {
			canceledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			return testCase{
				name: "failed to qdrant client healthcheck",
				clientFn: func(c *qdrant.Client) error {
					return nil
				},
				input: input{
					ctx: canceledCtx,
					config: &config.Config{
						QdrantConfig: config.QdrantConfig{
							Host:           "localhost",
							CollectionName: "collection",
							VectorSize:     1,
						},
					},
					tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
					logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
				},
				want: want{
					nilRepo: true,
					err:     true,
				},
			}
		}(),
		{
			name: "failed to get collection info",
			clientFn: func(c *qdrant.Client) error {
				return nil
			},
			input: input{
				ctx: context.Background(),
				config: &config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "",
						VectorSize:     1,
					},
				},
				tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
				logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			},
			want: want{
				nilRepo: true,
				err:     true,
			},
		},
		{
			name: "ok",
			clientFn: func(c *qdrant.Client) error {
				return nil
			},
			input: input{
				ctx: context.Background(),
				config: &config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     1,
					},
				},
				tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
				logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			},
			want: want{
				nilRepo: false,
				err:     false,
			},
		},
		func() testCase {
			collectionName := "collection"

			return testCase{
				name: "ok_collection_exists",
				clientFn: func(client *qdrant.Client) error {
					return client.CreateCollection(context.Background(), &qdrant.CreateCollection{CollectionName: collectionName})
				},
				input: input{
					ctx: context.Background(),
					config: &config.Config{
						QdrantConfig: config.QdrantConfig{
							Host:           "localhost",
							CollectionName: collectionName,
							VectorSize:     1,
						},
					},
					tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
					logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
				},
				want: want{
					nilRepo: false,
					err:     false,
				},
			}
		}(),
		{
			name: "failed to create collection",
			clientFn: func(c *qdrant.Client) error {
				return nil
			},
			input: input{
				ctx: context.Background(),
				config: &config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     0,
					},
				},
				tracer: otel_trace_noop.NewTracerProvider().Tracer(""),
				logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			},
			want: want{
				nilRepo: true,
				err:     true,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qdrantGRPCPort, cleanup := test_server.SetupQdrantServer(t)
			t.Cleanup(cleanup)

			client, err := qdrant.NewClient(&qdrant.Config{
				Port: qdrantGRPCPort,
			})
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			if err := tt.clientFn(client); err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			tt.input.config.QdrantConfig.GRPCPort = uint16(qdrantGRPCPort)
			repo, err := qdrant_infras.NewVectorRepo(
				tt.input.ctx,
				tt.input.config,
				tt.input.tracer,
				tt.input.logger,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if (repo == nil) != tt.want.nilRepo {
				t.Fatal(cmp.Diff(repo, nil))
			}
		})
	}
}

func Test_QdrantRepo_Insert(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx        context.Context
		embeddings []*domain.VectorRepoInsertEmbedding
	}

	type want struct {
		err bool
	}

	type testCase struct {
		name     string
		config   config.Config
		clientFn func(*qdrant.Client) error
		input    input
		want     want
	}
	testCases := []testCase{
		func() testCase {
			embeddings := []*domain.VectorRepoInsertEmbedding{
				{
					ID:     uuid.NewString(),
					Vector: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
					Metadata: map[string]any{
						"k": complex64(0),
					},
				},
			}

			return testCase{
				name: "failed to convert to qdrant map",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     1,
					},
				},
				clientFn: func(client *qdrant.Client) error {
					return nil
				},
				input: input{
					ctx:        context.Background(),
					embeddings: embeddings,
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			embeddings := []*domain.VectorRepoInsertEmbedding{}

			return testCase{
				name: "failed to qdrant client upsert",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     1,
					},
				},
				clientFn: func(c *qdrant.Client) error {
					return nil
				},
				input: input{
					ctx:        context.Background(),
					embeddings: embeddings,
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			embeddings := []*domain.VectorRepoInsertEmbedding{
				{
					ID:       uuid.NewString(),
					Vector:   []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
					Metadata: map[string]any{"key1": "value1"},
				},
				{
					ID:       uuid.NewString(),
					Vector:   []float32{9, 8, 7, 6, 5, 4, 3, 2, 1},
					Metadata: map[string]any{"key2": "value2"},
				},
				{
					ID:       uuid.NewString(),
					Vector:   []float32{1, 1, 1, 1, 1, 1, 1, 1, 1},
					Metadata: map[string]any{"key3": "value3"},
				},
			}

			return testCase{
				name: "ok",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     1,
					},
				},
				clientFn: func(c *qdrant.Client) error {
					return nil
				},
				input: input{
					ctx:        context.Background(),
					embeddings: embeddings,
				},
				want: want{
					err: false,
				},
			}
		}(),
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			qdrantGRPCPort, cleanup := test_server.SetupQdrantServer(t)
			t.Cleanup(cleanup)

			client, err := qdrant.NewClient(&qdrant.Config{
				Port: qdrantGRPCPort,
			})
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			if err := tt.clientFn(client); err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			tt.config.QdrantConfig.GRPCPort = uint16(qdrantGRPCPort)
			repo, err := qdrant_infras.NewVectorRepo(
				ctx,
				&tt.config,
				otel_trace_noop.NewTracerProvider().Tracer(""),
				slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			)
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			err = repo.Insert(
				tt.input.ctx,
				tt.input.embeddings,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}
		})
	}
}

func Test_QdrantRepo_Query(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx   context.Context
		query *domain.VectorRepoQueryInput
	}

	type want struct {
		results []*domain.VectorRepoQueryResult
		err     bool
	}

	type testCase struct {
		name     string
		config   config.Config
		clientFn func(*qdrant.Client) error
		input    input
		want     want
	}
	testCases := []testCase{
		func() testCase {
			return testCase{
				name: "failed to qdrant client query",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: "collection",
						VectorSize:     1,
					},
				},
				clientFn: func(c *qdrant.Client) error {
					return nil
				},
				input: input{
					ctx: context.Background(),
					query: &domain.VectorRepoQueryInput{
						Vector:   []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
						TopK:     10,
						MinScore: 0.1,
						Filter:   nil,
					},
				},
				want: want{
					results: nil,
					err:     true,
				},
			}
		}(),
		func() testCase {
			collectionName := "collection"

			id1 := uuid.NewString()
			id2 := uuid.NewString()
			id3 := uuid.NewString()
			vector1 := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vector2 := []float32{9, 8, 7, 6, 5, 4, 3, 2, 1}
			vector3 := []float32{1, 1, 1, 1, 1, 1, 1, 1, 1}
			vectorSize := len(vector1)
			metadata1 := map[string]any{
				"k0": nil,
				"k1": true,
				"k2": int64(2),
				"k3": 3.0,
				"k4": "v4",
				"k5": map[string]any{
					"k1": int64(1),
					"k2": "v2",
				},
				"k6": []any{"v1", "v2", "v3", "v4", "v5", "v6"},
			}
			metadata2 := map[string]any{
				"k0": "v0",
				"k1": false,
				"k2": int64(3),
				"k3": 4.0,
				"k4": "v5",
			}
			metadata3 := map[string]any{
				"k0": "v1",
				"k1": true,
				"k2": int64(4),
				"k3": 5.0,
				"k4": "v6",
			}

			return testCase{
				name: "ok_multiple_documents",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: collectionName,
						VectorSize:     vectorSize,
					},
				},
				clientFn: func(client *qdrant.Client) error {
					if err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
						CollectionName: collectionName,
						VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
							Size:     uint64(vectorSize),
							Distance: qdrant.Distance_Cosine,
						}),
					}); err != nil {
						return err
					}

					if _, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
						Wait:           qdrant.PtrOf(true),
						CollectionName: collectionName,
						Points: []*qdrant.PointStruct{
							{
								Id:      qdrant.NewID(id1),
								Vectors: qdrant.NewVectors(vector1...),
								Payload: qdrant.NewValueMap(metadata1),
							},
							{
								Id:      qdrant.NewID(id2),
								Vectors: qdrant.NewVectors(vector2...),
								Payload: qdrant.NewValueMap(metadata2),
							},
							{
								Id:      qdrant.NewID(id3),
								Vectors: qdrant.NewVectors(vector3...),
								Payload: qdrant.NewValueMap(metadata3),
							},
						},
					}); err != nil {
						return err
					}

					return nil
				},
				input: input{
					ctx: context.Background(),
					query: &domain.VectorRepoQueryInput{
						Vector:   vector1,
						TopK:     10,
						MinScore: 0.1,
						Filter:   nil,
					},
				},
				want: want{
					results: []*domain.VectorRepoQueryResult{
						{
							ID:       id1,
							Score:    1,
							Vector:   cosineNormalize(vector1),
							Metadata: metadata1,
						},
						{
							ID:       id2,
							Score:    0.57894737,
							Vector:   cosineNormalize(vector2),
							Metadata: metadata2,
						},
						{
							ID:       id3,
							Score:    0.88852334,
							Vector:   cosineNormalize(vector3),
							Metadata: metadata3,
						},
					},
					err: false,
				},
			}
		}(),
		func() testCase {
			collectionName := "collection"

			id1 := uuid.NewString()
			id2 := uuid.NewString()
			id3 := uuid.NewString()
			vector1 := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vector2 := []float32{9, 8, 7, 6, 5, 4, 3, 2, 1}
			vector3 := []float32{1, 1, 1, 1, 1, 1, 1, 1, 1}
			vectorSize := len(vector1)
			metadata1 := map[string]any{
				"k0": nil,
				"k1": true,
				"k2": int64(2),
				"k3": 3.0,
				"k4": "v4",
			}
			metadata2 := map[string]any{
				"k0": "v0",
				"k1": false,
				"k2": int64(3),
				"k3": 4.0,
				"k4": "v5",
			}
			metadata3 := map[string]any{
				"k0": "v1",
				"k1": true,
				"k2": int64(4),
				"k3": 5.0,
				"k4": "v6",
			}

			return testCase{
				name: "ok_min_score_filters_out_documents",
				config: config.Config{
					QdrantConfig: config.QdrantConfig{
						Host:           "localhost",
						CollectionName: collectionName,
						VectorSize:     vectorSize,
					},
				},
				clientFn: func(client *qdrant.Client) error {
					if err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
						CollectionName: collectionName,
						VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
							Size:     uint64(vectorSize),
							Distance: qdrant.Distance_Cosine,
						}),
					}); err != nil {
						return err
					}

					if _, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
						Wait:           qdrant.PtrOf(true),
						CollectionName: collectionName,
						Points: []*qdrant.PointStruct{
							{
								Id:      qdrant.NewID(id1),
								Vectors: qdrant.NewVectors(vector1...),
								Payload: qdrant.NewValueMap(metadata1),
							},
							{
								Id:      qdrant.NewID(id2),
								Vectors: qdrant.NewVectors(vector2...),
								Payload: qdrant.NewValueMap(metadata2),
							},
							{
								Id:      qdrant.NewID(id3),
								Vectors: qdrant.NewVectors(vector3...),
								Payload: qdrant.NewValueMap(metadata3),
							},
						},
					}); err != nil {
						return err
					}

					return nil
				},
				input: input{
					ctx: context.Background(),
					query: &domain.VectorRepoQueryInput{
						Vector:   vector1,
						TopK:     10,
						MinScore: 0.9,
						Filter:   nil,
					},
				},
				want: want{
					results: []*domain.VectorRepoQueryResult{
						{
							ID:       id1,
							Score:    1,
							Vector:   cosineNormalize(vector1),
							Metadata: metadata1,
						},
					},
					err: false,
				},
			}
		}(),
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			qdrantGRPCPort, cleanup := test_server.SetupQdrantServer(t)
			t.Cleanup(cleanup)

			client, err := qdrant.NewClient(&qdrant.Config{
				Port: qdrantGRPCPort,
			})
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			if err := tt.clientFn(client); err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			tt.config.QdrantConfig.GRPCPort = uint16(qdrantGRPCPort)
			repo, err := qdrant_infras.NewVectorRepo(
				ctx,
				&tt.config,
				otel_trace_noop.NewTracerProvider().Tracer(""),
				slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			)
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}

			result, err := repo.Query(
				tt.input.ctx,
				tt.input.query,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			sortSlicesOpts := cmpopts.SortSlices(func(a, b *domain.VectorRepoQueryResult) bool { return a.ID < b.ID })

			if !cmp.Equal(result, tt.want.results, sortSlicesOpts) {
				t.Fatal(cmp.Diff(result, tt.want.results, sortSlicesOpts))
			}
		})
	}
}

func cosineNormalize(vector []float32) []float32 {
	normalized := make([]float32, len(vector))
	var magnitude float64
	for _, v := range vector {
		magnitude += float64(v) * float64(v)
	}
	magnitudeSqrt := float32(math.Sqrt(magnitude))
	if magnitude > 0 {
		for i, v := range vector {
			normalized[i] = v / magnitudeSqrt
		}
	}
	return normalized
}
