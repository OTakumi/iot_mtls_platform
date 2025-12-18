// Package persistence provides the persistence layer of the application.
package persistence

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
)

// DeviceGormRepository is a GORM implementation of the DeviceRepository interface.
type DeviceGormRepository struct {
	db *gorm.DB
}

// NewDeviceGormRepository creates a new instance of DeviceGormRepository.
//
//nolint:ireturn
func NewDeviceGormRepository(db *gorm.DB) repository.DeviceRepository {
	return &DeviceGormRepository{db: db}
}

// Save saves a new Device entity or updates an existing one.
func (r *DeviceGormRepository) Save(ctx context.Context, device *entity.Device) error {
	// GORMのSaveメソッドは、主キーが存在すれば更新、なければ新規作成
	return r.db.WithContext(ctx).Save(device).Error
}

// FindByID finds a Device entity by its UUID.
func (r *DeviceGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	var device entity.Device
	// IDで検索し、レコードが見つからない場合はgorm.ErrRecordNotFoundを返す
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// FindByHardwareID finds a Device entity by its HardwareID.
func (r *DeviceGormRepository) FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error) {
	var device entity.Device
	// HardwareIDで検索し、レコードが見つからない場合はgorm.ErrRecordNotFoundを返す
	err := r.db.WithContext(ctx).Where("hardware_id = ?", hardwareID).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// FindAll finds all Device entities.
func (r *DeviceGormRepository) FindAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	// 全てのレコードを取得
	err := r.db.WithContext(ctx).Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Delete deletes a Device entity by its UUID.
func (r *DeviceGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// IDでレコードを削除
	// 削除対象が見つからない場合もエラーとしない (DeletedAtを使用していないため)
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Device{}).Error //nolint:exhaustruct
}
