package usecase

import (
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/entity"
)

// CreateDeviceInput is the input data for creating a Device.
type CreateDeviceInput struct {
	HardwareID string         // 必須
	Name       string         // オプショナル
	Metadata   map[string]any // オプショナル
}

// UpdateDeviceInput is the input data for updating a Device.
type UpdateDeviceInput struct {
	ID       uuid.UUID
	Name     *string        // オプショナル、nilで渡された場合は更新しない
	Metadata map[string]any // オプショナル、nilで渡された場合は更新しない
}

// DeviceOutput is the output data for displaying Device information.
type DeviceOutput struct {
	ID         uuid.UUID      `json:"id"`
	HardwareID string         `json:"hardwareId"`
	Name       string         `json:"name,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

// NewDeviceOutput creates a new DeviceOutput from an entity.
func NewDeviceOutput(device *entity.Device) *DeviceOutput {
	return &DeviceOutput{
		ID:         device.ID,
		HardwareID: device.HardwareID,
		Name:       device.Name,
		Metadata:   device.Metadata,
		CreatedAt:  device.CreatedAt,
		UpdatedAt:  device.UpdatedAt,
	}
}
