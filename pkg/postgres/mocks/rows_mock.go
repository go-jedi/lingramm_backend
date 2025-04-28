package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type RowsMock struct {
	CloseFunc             func()
	ErrFunc               func() error
	CommandTagFunc        func() pgconn.CommandTag
	FieldDescriptionsFunc func() []pgconn.FieldDescription
	NextFunc              func() bool
	ScanFunc              func(dest ...any) error
	ValuesFunc            func() ([]any, error)
	RawValuesFunc         func() [][]byte
	ConnFunc              func() *pgx.Conn
}

func (m *RowsMock) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}

func (m *RowsMock) Err() error {
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}

func (m *RowsMock) CommandTag() pgconn.CommandTag {
	if m.CommandTagFunc != nil {
		return m.CommandTagFunc()
	}
	return pgconn.NewCommandTag("")
}

func (m *RowsMock) FieldDescriptions() []pgconn.FieldDescription {
	if m.FieldDescriptionsFunc != nil {
		return m.FieldDescriptionsFunc()
	}
	return nil
}

func (m *RowsMock) Next() bool {
	if m.NextFunc != nil {
		return m.NextFunc()
	}
	return false
}

func (m *RowsMock) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

func (m *RowsMock) Values() ([]any, error) {
	if m.ValuesFunc != nil {
		return m.ValuesFunc()
	}
	return nil, nil
}

func (m *RowsMock) RawValues() [][]byte {
	if m.RawValuesFunc != nil {
		return m.RawValuesFunc()
	}
	return nil
}

func (m *RowsMock) Conn() *pgx.Conn {
	if m.ConnFunc != nil {
		return m.ConnFunc()
	}
	return nil
}
