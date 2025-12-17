package persistence

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
)

// DeviceRepositoryインターフェースのGORM実装
type DeviceGormRepository struct {
	db *gorm.DB
}

// DeviceGormRepositoryの新しいインスタンスを生成
func NewDeviceGormRepository(db *gorm.DB) repository.DeviceRepository {
	return &DeviceGormRepository{db: db}
}

// 新しいDeviceエンティティを保存または既存のエンティティを更新
func (r *DeviceGormRepository) Save(ctx context.Context, device *entity.Device) error {
	// GORMのSaveメソッドは、主キーが存在すれば更新、なければ新規作成
	return r.db.WithContext(ctx).Save(device).Error
}

// 指定されたUUIDを持つDeviceエンティティを検索
func (r *DeviceGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	var device entity.Device
	// IDで検索し、レコードが見つからない場合はgorm.ErrRecordNotFoundを返す
	if err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// 指定されたHardwareIDを持つDeviceエンティティを検索
func (r *DeviceGormRepository) FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error) {
	var device entity.Device
	// HardwareIDで検索し、レコードが見つからない場合はgorm.ErrRecordNotFoundを返す
	if err := r.db.WithContext(ctx).Where("hardware_id = ?", hardwareID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// 全てのDeviceエンティティを取得
func (r *DeviceGormRepository) FindAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	// 全てのレコードを取得
	if err := r.db.WithContext(ctx).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// 指定されたUUIDを持つDeviceエンティティを削除
func (r *DeviceGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// IDでレコードを削除
	// 削除対象が見つからない場合もエラーとしない (DeletedAtを使用していないため)
	return r.db.WithContext(ctx).Delete(&entity.Device{}, "id = ?", id).Error
}
