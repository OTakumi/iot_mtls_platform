package domain

import (
	"testing"
)

func TestNewDevice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hardwareID := "some-hardware-id"
		name := "some-name"
		metadata := []string{"meta1", "meta2"}

		device, err := NewDevice(hardwareID, name, metadata)

		if err != nil {
			t.Errorf("NewDevice() returned an error: %v", err)
		}

		if device.hardware_id != hardwareID {
			t.Errorf("hardware_id is not set correctly. got %s, want %s", device.hardware_id, hardwareID)
		}

		if device.name != name {
			t.Errorf("name is not set correctly. got %s, want %s", device.name, name)
		}

		if len(device.metadata) != len(metadata) {
			t.Errorf("metadata is not set correctly. got %d items, want %d", len(device.metadata), len(metadata))
		}

		if device.id.String() == "00000000-0000-0000-0000-000000000000" {
			t.Error("id is not set")
		}

		if device.created_at.IsZero() {
			t.Error("created_at is not set")
		}
	})

	t.Run("empty hardware_id", func(t *testing.T) {
		_, err := NewDevice("", "some-name", nil)
		if err == nil {
			t.Error("NewDevice() did not return an error for empty hardware_id")
		}

		expectedError := "hardware id cannot be empty"
		if err.Error() != expectedError {
			t.Errorf("NewDevice() returned wrong error. got %q, want %q", err.Error(), expectedError)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := NewDevice("some-hardware-id", "", nil)
		if err == nil {
			t.Error("NewDevice() did not return an error for empty name")
		}

		expectedError := "name cannot be empty"
		if err.Error() != expectedError {
			t.Errorf("NewDevice() returned wrong error. got %q, want %q", err.Error(), expectedError)
		}
	})
}
