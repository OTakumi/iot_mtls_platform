package repository

import (
	"context"

	"backend/internal/domain/entity"

	"github.com/google/uuid"
)

type DeviceRepository interface {
	Save(ctx context.Context, device *entity.Device) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error)
	FindAll(ctx context.Context) ([]*entity.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
