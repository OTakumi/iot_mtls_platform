// Package persistence provides the persistence layer of the application.
package persistence

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
)

// DeviceGormRepository is the GORM implementation of the DeviceRepository.
type DeviceGormRepository struct {
	db *gorm.DB
}

// NewDeviceGormRepository creates a new instance of DeviceGormRepository.
//
//nolint:ireturn
func NewDeviceGormRepository(db *gorm.DB) repository.DeviceRepository {
	return &DeviceGormRepository{db: db}
}

// Save creates a new device or updates an existing one.
func (r *DeviceGormRepository) Save(ctx context.Context, device *entity.Device) error {
	// GORM's Save method handles both creation (if primary key is zero) and update.
	return r.db.WithContext(ctx).Save(device).Error
}

// FindByID finds a device by its UUID.
func (r *DeviceGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	var device entity.Device
	// It returns `gorm.ErrRecordNotFound` if no record is found.
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// FindByHardwareID finds a device by its hardware ID.
func (r *DeviceGormRepository) FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error) {
	var device entity.Device
	// It returns `gorm.ErrRecordNotFound` if no record is found.
	err := r.db.WithContext(ctx).Where("hardware_id = ?", hardwareID).First(&device).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// FindAll retrieves all devices.
func (r *DeviceGormRepository) FindAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	// It returns an empty slice if no devices are found.
	err := r.db.WithContext(ctx).Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Delete removes a device by its UUID.
func (r *DeviceGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// It deletes a record by its primary key.
	// If the record to be deleted is not found, GORM does not return an error, but RowsAffected will be 0.
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Device{}) //nolint:exhaustruct
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrDeviceNotFound
	}

	return nil
}
