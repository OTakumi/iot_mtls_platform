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

// GORMリポジトリのCRUD操作を実際のDBで検証する結合テスト.
func TestDeviceGormRepository_Integration(t *testing.T) {
	// testDBは main_test.go で初期化される
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	repo := persistence.NewDeviceGormRepository(testDB)
	ctx := context.Background()

	// Save (Create)
	t.Run("Save(Create) - Creates a new device", func(t *testing.T) {
		cleanupTable(t)

		// テストデータ作成
		deviceName := "test-device"
		device, err := entity.NewDevice("hw-create-01", &deviceName, map[string]any{"key": "val"})
		require.NoError(t, err)

		// 実行
		err = repo.Save(ctx, device)

		// 検証
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, device.ID, "ID should be assigned after saving")

		// DBから再取得して検証
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

		// テストデータ作成 & 保存
		deviceName := "find-me"
		savedDevice, err := entity.NewDevice("hw-find-01", &deviceName, nil)
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// 実行
		foundDevice, err := repo.FindByID(ctx, savedDevice.ID)

		// 検証
		require.NoError(t, err)
		assert.NotNil(t, foundDevice)
		assert.Equal(t, savedDevice.ID, foundDevice.ID)
		assert.Equal(t, savedDevice.HardwareID, foundDevice.HardwareID)
	})

	t.Run("FindByID - Returns error for non-existent ID", func(t *testing.T) {
		cleanupTable(t)

		// 実行
		nonExistentID := uuid.New()
		foundDevice, err := repo.FindByID(ctx, nonExistentID)

		// 検証
		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "should return gorm.ErrRecordNotFound")
		assert.Nil(t, foundDevice)
	})

	// Save (Update)
	t.Run("Save(Update) - Updates an existing device", func(t *testing.T) {
		cleanupTable(t)

		// テストデータ作成 & 保存
		originalName := "original-name"
		savedDevice, err := entity.NewDevice("hw-update-01", &originalName, map[string]any{"status": "inactive"})
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// フィールドを更新
		updatedName := "updated-name"
		savedDevice.Name = updatedName
		savedDevice.Metadata = map[string]any{"status": "active", "version": 2}

		// 実行
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// DBから再取得して検証
		foundDevice, err := repo.FindByID(ctx, savedDevice.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedName, foundDevice.Name)
		assert.Equal(t, entity.JSONBMap{"status": "active", "version": float64(2)},
			foundDevice.Metadata,
			"JSONB number type becomes float64",
		)
	})

	// Delete
	t.Run("Delete - Deletes an existing device", func(t *testing.T) {
		cleanupTable(t)

		// テストデータ作成 & 保存
		deviceName := "delete-me"
		savedDevice, err := entity.NewDevice("hw-delete-01", &deviceName, nil)
		require.NoError(t, err)
		err = repo.Save(ctx, savedDevice)
		require.NoError(t, err)

		// 実行
		err = repo.Delete(ctx, savedDevice.ID)
		require.NoError(t, err)

		// DBから再取得して、見つからないことを確認
		_, err = repo.FindByID(ctx, savedDevice.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
