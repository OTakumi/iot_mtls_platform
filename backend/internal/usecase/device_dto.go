package usecase

import (
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/entity"
)

// Device作成用の入力データ
type CreateDeviceInput struct {
	HardwareID string         // 必須
	Name       string         // オプショナル
	Metadata   map[string]any // オプショナル
}

// Device更新用の入力データ
type UpdateDeviceInput struct {
	ID       uuid.UUID
	Name     *string        // オプショナル、nilで渡された場合は更新しない
	Metadata map[string]any // オプショナル、nilで渡された場合は更新しない
}

// エンティティから出力DTOを生成
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

// Device情報を表示するための出力データ
type DeviceOutput struct {
	ID         uuid.UUID      `json:"id"`
	HardwareID string         `json:"hardware_id"`
	Name       string         `json:"name,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
