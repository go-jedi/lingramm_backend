package mocks

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type RowMock struct {
	mock.Mock
}

func NewMockRow(t *testing.T) *RowMock {
	m := &RowMock{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *RowMock) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}
