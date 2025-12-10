package id

import (
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if id.Value.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("NewID() returned a zero UUID")
	}
}

func TestParse(t *testing.T) {
	t.Run("valid uuid", func(t *testing.T) {
		validUUID := "123e4567-e89b-12d3-a456-426614174000"
		id, err := Parse(validUUID)
		if err != nil {
			t.Errorf("Parse() returned an error for a valid UUID: %v", err)
		}
		if id.Value.String() != validUUID {
			t.Errorf("Parse() returned an incorrect UUID. got %v, want %v", id.Value.String(), validUUID)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		invalidUUID := "invalid-uuid"
		_, err := Parse(invalidUUID)
		if err == nil {
			t.Error("Parse() did not return an error for an invalid UUID")
		}
	})
}
