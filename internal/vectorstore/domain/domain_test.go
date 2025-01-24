package domain_test

import (
	"context"
	"strings"
	"testing"

	internal_error "github.com/aria3ppp/rag-server/internal/pkg/error"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

func Test_InsertTextsInput_Validate(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx context.Context
	}

	type want struct {
		err                 bool
		validationErr       bool
		validationErrString string
	}

	type testCase struct {
		name         string
		domainObject *domain.InsertTextsInput
		input        input
		want         want
	}
	testCases := []testCase{
		{
			name: "ok",
			domainObject: &domain.InsertTextsInput{
				Texts: []*domain.InsertTextsInputText{
					{
						Text:     strings.Repeat("x", 100),
						Metadata: nil,
					},
				},
			},
			input: input{
				ctx: context.Background(),
			},
			want: want{
				err:                 false,
				validationErr:       false,
				validationErrString: "",
			},
		},
		{
			name: "validation_error_texts",
			domainObject: &domain.InsertTextsInput{
				Texts: []*domain.InsertTextsInputText{},
			},
			input: input{
				ctx: context.Background(),
			},
			want: want{
				err:           true,
				validationErr: true,
				validationErrString: func() string {
					d := &domain.InsertTextsInput{
						Texts: []*domain.InsertTextsInputText{},
					}
					validator := validatorPkg.New(validatorPkg.WithRequiredStructEnabled())
					err := validator.StructCtx(context.Background(), d)
					validationErr, ok := err.(validatorPkg.ValidationErrors)
					if !ok {
						panic("validator.ValidationErrors didn't happen")
					}
					return internal_error.NewValidationError(validationErr).Error()
				}(),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.domainObject.Validate(tt.input.ctx)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if tt.want.validationErr {
				validationErr, ok := err.(*internal_error.ValidationError)
				if !ok {
					t.Fatal(cmp.Diff(ok, true))
				}

				if !cmp.Equal(validationErr.Error(), tt.want.validationErrString) {
					t.Fatal(cmp.Diff(validationErr.Error(), tt.want.validationErrString))
				}
			}
		})
	}
}

func Test_SearchTextInput_Validate(t *testing.T) {
	t.Parallel()

	type input struct {
		ctx context.Context
	}

	type want struct {
		err                 bool
		validationErr       bool
		validationErrString string
	}

	type testCase struct {
		name         string
		domainObject *domain.SearchTextInput
		input        input
		want         want
	}
	testCases := []testCase{
		{
			name: "ok",
			domainObject: &domain.SearchTextInput{
				Text:     "t",
				TopK:     1,
				MinScore: 0.1,
				Filter:   nil,
			},
			input: input{
				ctx: context.Background(),
			},
			want: want{
				err:                 false,
				validationErr:       false,
				validationErrString: "",
			},
		},
		{
			name: "validation_error_texts",
			domainObject: &domain.SearchTextInput{
				Text:     "",
				TopK:     0,
				MinScore: 0.0,
				Filter:   nil,
			},
			input: input{
				ctx: context.Background(),
			},
			want: want{
				err:           true,
				validationErr: true,
				validationErrString: func() string {
					d := &domain.SearchTextInput{
						Text:     "",
						TopK:     0,
						MinScore: 0.0,
						Filter:   nil,
					}
					validator := validatorPkg.New(validatorPkg.WithRequiredStructEnabled())
					err := validator.StructCtx(context.Background(), d)
					validationErr, ok := err.(validatorPkg.ValidationErrors)
					if !ok {
						panic("validator.ValidationErrors didn't happen")
					}
					return internal_error.NewValidationError(validationErr).Error()
				}(),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.domainObject.Validate(tt.input.ctx)
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			if tt.want.validationErr {
				validationErr, ok := err.(*internal_error.ValidationError)
				if !ok {
					t.Fatal(cmp.Diff(ok, true))
				}

				if !cmp.Equal(validationErr.Error(), tt.want.validationErrString) {
					t.Fatal(cmp.Diff(validationErr.Error(), tt.want.validationErrString))
				}
			}
		})
	}
}
