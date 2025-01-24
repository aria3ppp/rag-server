package usecase_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/aria3ppp/rag-server/internal/vectorstore/config"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase/mocks"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

type mockups struct {
	embedder    *mocks.MockEmbedder
	vectorRepo  *mocks.MockVectorRepo
	idGenerator *mocks.MockIDGenerator
}

func Test_UseCase_InsertTexts(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx   context.Context
		input *domain.InsertTextsInput
	}

	type want struct {
		err bool
	}

	type testCase struct {
		name   string
		mockFn func(mockups)
		input  input
		want   want
	}
	testCases := []testCase{
		{
			name:   "failed to validate input",
			mockFn: func(m mockups) {},
			input: input{
				ctx: context.Background(),
				input: &domain.InsertTextsInput{
					Texts: []*domain.InsertTextsInputText{
						{
							Text:     "",
							Metadata: nil,
						},
					},
				},
			},
			want: want{
				err: true,
			},
		},
		func() testCase {
			text := strings.Repeat("t", 100)
			var metadata map[string]any = nil

			return testCase{
				name: "failed to embed text",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{text}).Return(nil, errors.New("error")),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{
							{
								Text:     text,
								Metadata: metadata,
							},
						},
					},
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			text := strings.Repeat("t", 100)
			var metadata map[string]any = nil

			return testCase{
				name: "invalid embeddings length",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{text}).Return([][]float32{}, nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{
							{
								Text:     text,
								Metadata: metadata,
							},
						},
					},
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			text := strings.Repeat("t", 100)
			var metadata map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}

			return testCase{
				name: "failed to generate new id",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{text}).Return([][]float32{embedding}, nil),
						m.idGenerator.EXPECT().NewID().Return("", errors.New("error")),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{
							{
								Text:     text,
								Metadata: metadata,
							},
						},
					},
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			text := strings.Repeat("t", 100)
			var metadata map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			id := uuid.NewString()
			vectorstoreInsertEmbeddings := []*domain.VectorRepoInsertEmbedding{
				{
					ID:       id,
					Vector:   embedding,
					Metadata: lo.Assign(metadata, map[string]any{"text": text}),
				},
			}

			return testCase{
				name: "failed to repo insert",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{text}).Return([][]float32{embedding}, nil),
						m.idGenerator.EXPECT().NewID().Return(id, nil),
						m.vectorRepo.EXPECT().Insert(gomock.Any(), vectorstoreInsertEmbeddings).Return(errors.New("error")),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{
							{
								Text:     text,
								Metadata: metadata,
							},
						},
					},
				},
				want: want{
					err: true,
				},
			}
		}(),
		func() testCase {
			text := strings.Repeat("t", 100)
			var metadata map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			id := uuid.NewString()
			vectorstoreInsertEmbeddings := []*domain.VectorRepoInsertEmbedding{
				{
					ID:       id,
					Vector:   embedding,
					Metadata: lo.Assign(metadata, map[string]any{"text": text}),
				},
			}

			return testCase{
				name: "ok",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{text}).Return([][]float32{embedding}, nil),
						m.idGenerator.EXPECT().NewID().Return(id, nil),
						m.vectorRepo.EXPECT().Insert(gomock.Any(), vectorstoreInsertEmbeddings).Return(nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{
							{
								Text:     text,
								Metadata: metadata,
							},
						},
					},
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

			controller := gomock.NewController(t)
			m := mockups{
				embedder:    mocks.NewMockEmbedder(controller),
				vectorRepo:  mocks.NewMockVectorRepo(controller),
				idGenerator: mocks.NewMockIDGenerator(controller),
			}
			tt.mockFn(m)

			uc := usecase.NewUseCase(
				m.embedder,
				m.idGenerator,
				m.vectorRepo,
				&config.Config{},
				noop.NewTracerProvider().Tracer(""),
				slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			)

			err := uc.InsertTexts(
				tt.input.ctx,
				tt.input.input,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}
		})
	}
}

func Test_UseCase_SearchText(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx   context.Context
		input *domain.SearchTextInput
	}

	type want struct {
		result *domain.SearchTextResult
		err    bool
	}

	type testCase struct {
		name   string
		mockFn func(mockups)
		input  input
		want   want
	}
	testCases := []testCase{
		{
			name:   "failed to validate input",
			mockFn: func(m mockups) {},
			input: input{
				ctx: context.Background(),
				input: &domain.SearchTextInput{
					Text:     "",
					TopK:     0,
					MinScore: 0.0,
					Filter:   nil,
				},
			},
			want: want{
				result: nil,
				err:    true,
			},
		},
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			return testCase{
				name: "failed to embed text",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return(nil, errors.New("error")),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			return testCase{
				name: "invalid vectors length",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{}, nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vectorRepoQueryInput := &domain.VectorRepoQueryInput{
				Vector:   embedding,
				TopK:     topK,
				MinScore: minScore,
				Filter:   filter,
			}

			return testCase{
				name: "failed to repo query",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{embedding}, nil),
						m.vectorRepo.EXPECT().Query(gomock.Any(), vectorRepoQueryInput).Return(nil, errors.New("error")),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vectorRepoQueryInput := &domain.VectorRepoQueryInput{
				Vector:   embedding,
				TopK:     topK,
				MinScore: minScore,
				Filter:   filter,
			}
			id := uuid.NewString()
			text := strings.Repeat("t", 100)
			metadata := map[string]any{
				"text": text,
			}
			resultMetadata := lo.FromEntries(lo.Entries(metadata))
			delete(resultMetadata, "text")
			score := float32(1)
			vectorRepoQueryResult := &domain.VectorRepoQueryResult{
				ID:       id,
				Score:    score,
				Vector:   embedding,
				Metadata: metadata,
			}

			return testCase{
				name: "query results length couldn't exceed topk",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{embedding}, nil),
						m.vectorRepo.EXPECT().Query(gomock.Any(), vectorRepoQueryInput).Return(lo.Times(topK+1, func(_ int) *domain.VectorRepoQueryResult { return vectorRepoQueryResult }), nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vectorRepoQueryInput := &domain.VectorRepoQueryInput{
				Vector:   embedding,
				TopK:     topK,
				MinScore: minScore,
				Filter:   filter,
			}
			id := uuid.NewString()
			metadata := map[string]any{}
			score := float32(1)
			vectorRepoQueryResult := &domain.VectorRepoQueryResult{
				ID:       id,
				Score:    score,
				Vector:   embedding,
				Metadata: metadata,
			}

			return testCase{
				name: "metadata field text not exist",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{embedding}, nil),
						m.vectorRepo.EXPECT().Query(gomock.Any(), vectorRepoQueryInput).Return([]*domain.VectorRepoQueryResult{vectorRepoQueryResult}, nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vectorRepoQueryInput := &domain.VectorRepoQueryInput{
				Vector:   embedding,
				TopK:     topK,
				MinScore: minScore,
				Filter:   filter,
			}
			id := uuid.NewString()
			metadata := map[string]any{
				"text": false,
			}
			score := float32(1)
			vectorRepoQueryResult := &domain.VectorRepoQueryResult{
				ID:       id,
				Score:    score,
				Vector:   embedding,
				Metadata: metadata,
			}

			return testCase{
				name: "metadata field text is not a string",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{embedding}, nil),
						m.vectorRepo.EXPECT().Query(gomock.Any(), vectorRepoQueryInput).Return([]*domain.VectorRepoQueryResult{vectorRepoQueryResult}, nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: nil,
					err:    true,
				},
			}
		}(),
		func() testCase {
			queryText := "text"
			topK := 10
			minScore := float32(0.1)
			var filter map[string]any = nil

			embedding := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
			vectorRepoQueryInput := &domain.VectorRepoQueryInput{
				Vector:   embedding,
				TopK:     topK,
				MinScore: minScore,
				Filter:   filter,
			}
			id := uuid.NewString()
			text := strings.Repeat("t", 100)
			metadata := map[string]any{
				"text": text,
			}
			resultMetadata := lo.FromEntries(lo.Entries(metadata))
			delete(resultMetadata, "text")
			score := float32(1)
			vectorRepoQueryResult := &domain.VectorRepoQueryResult{
				ID:       id,
				Score:    score,
				Vector:   embedding,
				Metadata: metadata,
			}

			return testCase{
				name: "ok",
				mockFn: func(m mockups) {
					gomock.InOrder(
						m.embedder.EXPECT().Embed(gomock.Any(), []string{queryText}).Return([][]float32{embedding}, nil),
						m.vectorRepo.EXPECT().Query(gomock.Any(), vectorRepoQueryInput).Return([]*domain.VectorRepoQueryResult{vectorRepoQueryResult}, nil),
					)
				},
				input: input{
					ctx: context.Background(),
					input: &domain.SearchTextInput{
						Text:     queryText,
						TopK:     topK,
						MinScore: minScore,
						Filter:   filter,
					},
				},
				want: want{
					result: &domain.SearchTextResult{
						SimilarTexts: []*domain.SearchTextResultItem{
							{
								Text:     text,
								Score:    score,
								Metadata: resultMetadata,
							},
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

			controller := gomock.NewController(t)
			m := mockups{
				embedder:    mocks.NewMockEmbedder(controller),
				vectorRepo:  mocks.NewMockVectorRepo(controller),
				idGenerator: mocks.NewMockIDGenerator(controller),
			}
			tt.mockFn(m)

			uc := usecase.NewUseCase(
				m.embedder,
				m.idGenerator,
				m.vectorRepo,
				&config.Config{},
				noop.NewTracerProvider().Tracer(""),
				slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
			)

			result, err := uc.SearchText(
				tt.input.ctx,
				tt.input.input,
			)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if !cmp.Equal(result, tt.want.result) {
				t.Fatal(cmp.Diff(result, tt.want.result))
			}
		})
	}
}
