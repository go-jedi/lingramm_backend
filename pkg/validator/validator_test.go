package validator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestValidator(t *testing.T) {
	v := New()

	type person struct {
		Name string `json:"name" validate:"required,min=1"`
	}

	var (
		ctx = context.TODO()
	)

	tests := []struct {
		name string
		in   person
		err  error
	}{
		{
			name: "ok",
			in: person{
				Name: "ok",
			},
			err: nil,
		},
		{
			name: "validation error",
			in: person{
				Name: "",
			},
			err: errors.New("Key: 'person.Name' Error:Field validation for 'Name' failed on the 'required' tag"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := v.StructCtx(ctx, test.in)

			if test.err != nil {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), test.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
