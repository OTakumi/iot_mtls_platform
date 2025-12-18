package persistence_test

import (
	"context"
	"testing"

	"backend/internal/domain/entity"
	"backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestDeviceGormRepository_Integration performs integration tests for
// the GORM repository's CRUD operations against a real database.
func TestDeviceGormRepository_Integration(t *testing.T) {
	// testDB is initialized in main_test.go.
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	repo := persistence.NewDeviceGormRepository(testDB)
	ctx := context.Background()

	// Save (Create)
	t.Run("Save(Create) - Creates a new device", func(t *testing.T) {
		cleanupTable(t)

		// Prepare test data.
		deviceName := "test-device"
		device, err := entity.NewDevice("hw-create-01", &deviceName, map[string]any{"key": "val"})
		require.NoError(t, err)

		// Execute the method.
		err = repo.Save(ctx, device)
		// Assert the results.
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, device.ID, "ID should be assigned after saving")
		// Re-fetch from DB and verify the created data.
		foundDevice, err := repo.FindByID(ctx, device.ID)
		require.NoError(t, err)
		assert.NotNil(t, foundDevice)
		assert.Equal(t, device.HardwareID, foundDevice.HardwareID)
		assert.Equal(t, device.Name, foundDevice.Name)
		assert.Equal(t, device.Metadata, foundDevice.Metadata)
	})

	// FindByID
	t.Run("FindByID - Finds an existing device by ID", func(t *testing.T) {
		cleanupTable(t)

		// Prepare and save test data.
		deviceName := "find-me"
		savedDevice, err := entity.NewDevice("hw-find-01", &deviceName, nil)
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// Execute.
		foundDevice, err := repo.FindByID(ctx, savedDevice.ID)

		// Assert.
		require.NoError(t, err)
		assert.NotNil(t, foundDevice)
		assert.Equal(t, savedDevice.ID, foundDevice.ID)
		assert.Equal(t, savedDevice.HardwareID, foundDevice.HardwareID)
	})

	t.Run("FindByID - Returns error for non-existent ID", func(t *testing.T) {
		cleanupTable(t)

		// Execute.
		nonExistentID := uuid.New()
		foundDevice, err := repo.FindByID(ctx, nonExistentID)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "should return gorm.ErrRecordNotFound")
		assert.Nil(t, foundDevice)
	})

	// Save (Update)
	t.Run("Save(Update) - Updates an existing device", func(t *testing.T) {
		cleanupTable(t)

		// Prepare and save test data.
		originalName := "original-name"
		savedDevice, err := entity.NewDevice("hw-update-01", &originalName, map[string]any{"status": "inactive"})
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// Update the fields.
		updatedName := "updated-name"
		savedDevice.Name = updatedName
		savedDevice.Metadata = map[string]any{"status": "active", "version": 2}

		// Execute.
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// Re-fetch from the DB and verify.
		foundDevice, err := repo.FindByID(ctx, savedDevice.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedName, foundDevice.Name)
		assert.Equal(t, entity.JSONBMap{"status": "active", "version": float64(2)},
			foundDevice.Metadata,
			"Note: JSON numbers are unmarshaled as float64 in Go.",
		)
	})

	// Delete
	t.Run("Delete - Deletes an existing device", func(t *testing.T) {
		cleanupTable(t)

		// Prepare and save test data.
		deviceName := "delete-me"
		savedDevice, err := entity.NewDevice("hw-delete-01", &deviceName, nil)
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// Execute.
		err = repo.Delete(ctx, savedDevice.ID)
		require.NoError(t, err)

		// Verify that re-fetching from the DB fails (record not found).
		_, err = repo.FindByID(ctx, savedDevice.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
