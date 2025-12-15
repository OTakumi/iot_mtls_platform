package id

import (
	"fmt"

	"github.com/google/uuid"
)

type ID struct {
	Value uuid.UUID
}

func NewID() ID {
	return ID{Value: uuid.New()}
}

// parse uuid from string
func Parse(s string) (ID, error) {
	parsed, e := uuid.Parse(s)

	if e != nil {
		return ID{}, fmt.Errorf("invalid uuid format: %w", e)
	}

	return ID{Value: parsed}, nil
}

func (i ID) String() string {
	return i.Value.String()
}
