package domain

import (
	"context"

	internal_error "github.com/aria3ppp/rag-server/internal/pkg/error"

	validatorPkg "github.com/go-playground/validator/v10"
)

var validator = validatorPkg.New(validatorPkg.WithRequiredStructEnabled())

type Role int8

const (
	RoleUnspecified Role = iota
	RoleSystem
	RoleAssistant
	RoleUser
)

type StopReason int8

const (
	StopReasonUnspecified StopReason = iota
	StopReasonDone
	StopReasonError
)

type Message struct {
	Role    Role
	Content string
}

type QueryInput struct {
	Query    string     `validate:"required,min=2,max=2000"`
	Messages []*Message `validate:"-"`
}

func (input *QueryInput) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, input); err != nil {
		if _, ok := err.(validatorPkg.ValidationErrors); ok {
			return internal_error.NewValidationError(err)
		}
		return err
	}
	return nil
}

type QueryResult struct {
	Content     string
	CreatedInMS int64
}

type QueryStreamInput struct {
	Query    string     `validate:"required,min=2,max=2000"`
	Messages []*Message `validate:"-"`
}

func (input *QueryStreamInput) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, input); err != nil {
		if _, ok := err.(validatorPkg.ValidationErrors); ok {
			return internal_error.NewValidationError(err)
		}
		return err
	}
	return nil
}

type QueryStreamResultEvent struct {
	Content     string
	CreatedAtMS int64
	StopReason  StopReason
	Error       error
}
