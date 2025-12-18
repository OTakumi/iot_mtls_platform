// Package repository defines the repository interfaces.
package repository

import (
	"backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

// DeviceRepository defines the interface for persisting Device entities.
type DeviceRepository interface {
	// Save creates a new Device or updates an existing one.
	Save(ctx context.Context, device *entity.Device) error
	// FindByID retrieves a Device by its UUID.
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	// FindByHardwareID retrieves a Device by its hardware ID.
	FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error)
	// FindAll retrieves all Device entities.
	FindAll(ctx context.Context) ([]*entity.Device, error)
	// Delete removes a Device by its UUID.
	Delete(ctx context.Context, id uuid.UUID) error
}
