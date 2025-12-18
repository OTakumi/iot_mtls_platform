// Package id provides a value object for handling unique identifiers.
package id

import (
	"fmt"

	"github.com/google/uuid"
)

// ID is a value object for unique identifiers.
type ID struct {
	Value uuid.UUID
}

// NewID creates a new ID.
func NewID() ID {
	return ID{Value: uuid.New()}
}

// Parse parses a UUID from a string.
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
