package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// IValidator defines the interface for the validator.
//
//go:generate mockery --name=IValidator --output=mocks --case=underscore
type IValidator interface {
	Struct(s interface{}) error
	StructCtx(ctx context.Context, s interface{}) (err error)
}

// Validator is a wrapper around validator.Validate that implements IValidator.
type Validator struct {
	*validator.Validate
}

// Make sure Validator implements IValidator.
var _ IValidator = (*Validator)(nil)

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{validator.New(validator.WithRequiredStructEnabled())}
}
