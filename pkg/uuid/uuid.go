package uuid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// IUUID defines the interface for the uuid.
//
//go:generate mockery --name=IUUID --output=mocks --case=underscore
type IUUID interface {
	Generate() (string, error)
}

type UUID struct{}

func New() *UUID {
	return &UUID{}
}

type cryptoEntropy struct{}

func (cryptoEntropy) Read(p []byte) (int, error) {
	return rand.Read(p)
}

func (u *UUID) Generate() (string, error) {
	t := time.Now().UTC()

	e := cryptoEntropy{}

	id, err := ulid.New(ulid.Timestamp(t), e)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
