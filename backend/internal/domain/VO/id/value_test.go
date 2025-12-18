package id_test

import (
	"testing"

	"backend/internal/domain/VO/id"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	newID := id.NewID()
	if newID.Value.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("NewID() returned a zero UUID")
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	t.Run("valid uuid", func(t *testing.T) {
		t.Parallel()

		validUUID := "123e4567-e89b-12d3-a456-426614174000"

		parsedID, err := id.Parse(validUUID)
		if err != nil {
			t.Errorf("Parse() returned an error for a valid UUID: %v", err)
		}

		if parsedID.Value.String() != validUUID {
			t.Errorf("Parse() returned an incorrect UUID. got %v, want %v", parsedID.Value.String(), validUUID)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "invalid-uuid"

		_, err := id.Parse(invalidUUID)
		if err == nil {
			t.Error("Parse() did not return an error for an invalid UUID")
		}
	})
}
