package repository

import (
	"backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

// DeviceRepository Deviceエンティティの永続化を抽象化するインターフェース
type DeviceRepository interface {
	// 新しいDeviceエンティティを保存または既存のエンティティを更新
	Save(ctx context.Context, device *entity.Device) error
	// 指定されたUUIDを持つDeviceエンティティを検索
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	// 指定されたHardwareIDを持つDeviceエンティティを検索
	FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error)
	// 全てのDeviceエンティティを取得
	FindAll(ctx context.Context) ([]*entity.Device, error)
	// 指定されたUUIDを持つDeviceエンティティを削除
	Delete(ctx context.Context, id uuid.UUID) error
}
