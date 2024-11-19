package domain

import (
	"context"

	vectorstore_error "github.com/aria3ppp/rag-server/internal/vectorstore/error"

	validatorPkg "github.com/go-playground/validator/v10"
)

var validator = validatorPkg.New(validatorPkg.WithRequiredStructEnabled())

type InsertTextsInputText struct {
	Text     string         `validate:"required,min=100,max=2000"`
	Metadata map[string]any `validate:"-"`
}

type InsertTextsInput struct {
	Texts []*InsertTextsInputText `validate:"required,min=1"`
}

func (input *InsertTextsInput) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, input); err != nil {
		if _, ok := err.(validatorPkg.ValidationErrors); ok {
			return vectorstore_error.NewError(err.Error())
		}
		return err
	}
	return nil
}

type SearchTextInput struct {
	Text   string         `validate:"required"`
	TopK   int            `validate:"min=1,max=100"`
	Filter map[string]any `validate:"-"`
}

func (input *SearchTextInput) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, input); err != nil {
		if _, ok := err.(validatorPkg.ValidationErrors); ok {
			return vectorstore_error.NewError(err.Error())
		}
		return err
	}
	return nil
}

type SearchTextResult struct {
	SimilarTexts []*SearchTextResultItem
}

type SearchTextResultItem struct {
	Text     string
	Score    float32
	Metadata map[string]any
}
