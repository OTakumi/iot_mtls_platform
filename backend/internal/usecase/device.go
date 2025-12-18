package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
)

// DeviceUsecase エンティティに関するアプリケーション固有のビジネスロジックを定義するインターフェース.
type DeviceUsecase interface {
	// 新しいDeviceを登録
	CreateDevice(ctx context.Context, input CreateDeviceInput) (*DeviceOutput, error)
	// 指定されたIDのDeviceを取得
	GetDevice(ctx context.Context, id uuid.UUID) (*DeviceOutput, error)
	// 全てのDeviceを取得
	ListDevices(ctx context.Context) ([]*DeviceOutput, error)
	// 既存のDeviceを更新
	UpdateDevice(ctx context.Context, input UpdateDeviceInput) (*DeviceOutput, error)
	// 指定されたIDのDeviceを削除
	DeleteDevice(ctx context.Context, id uuid.UUID) error
}

// deviceUsecase DeviceUsecaseインターフェースの実装.
type deviceUsecase struct {
	deviceRepo repository.DeviceRepository
}

// NewDeviceUsecase deviceUsecaseの新しいインスタンスを生成.
//
//nolint:ireturn
func NewDeviceUsecase(repo repository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{deviceRepo: repo}
}

// CreateDevice 新しいDeviceを登録.
func (uc *deviceUsecase) CreateDevice(ctx context.Context, input CreateDeviceInput) (*DeviceOutput, error) {
	// ドメインエンティティの生成
	// IDはDB側で生成されるため、ここではuuid.Nilで問題ない
	var namePtr *string
	if input.Name != "" {
		namePtr = &input.Name
	}

	device, err := entity.NewDevice(input.HardwareID, namePtr, input.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create new device entity: %w", err)
	}

	// リポジトリを介して保存
	err = uc.deviceRepo.Save(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to save device: %w", err) // DB保存時のエラー
	}

	return NewDeviceOutput(device), nil // 出力DTOに変換して返す
}

// GetDevice 指定されたIDのDeviceを取得.
func (uc *deviceUsecase) GetDevice(ctx context.Context, id uuid.UUID) (*DeviceOutput, error) {
	device, err := uc.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find device by ID: %w", err)
	}

	if device == nil {
		return nil, entity.ErrDeviceNotFound
	}

	return NewDeviceOutput(device), nil
}

// ListDevices 全てのDeviceを取得.
func (uc *deviceUsecase) ListDevices(ctx context.Context) ([]*DeviceOutput, error) {
	devices, err := uc.deviceRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all devices: %w", err)
	}

	outputs := make([]*DeviceOutput, 0, len(devices))

	for _, device := range devices {
		outputs = append(outputs, NewDeviceOutput(device))
	}

	return outputs, nil
}

// UpdateDevice 既存のDeviceを更新.
func (uc *deviceUsecase) UpdateDevice(ctx context.Context, input UpdateDeviceInput) (*DeviceOutput, error) {
	// 更新対象のDeviceを検索
	device, err := uc.deviceRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find device by ID for update: %w", err)
	}

	if device == nil {
		return nil, entity.ErrDeviceNotFound
	}

	// エンティティの値を更新
	// HardwareIDはデバイス固有の情報のため更新不可
	// input.Nameは*stringなのでnilチェック
	if input.Name != nil {
		device.Name = *input.Name
	}

	// input.Metadataはmapなのでnilチェック
	if input.Metadata != nil {
		device.Metadata = input.Metadata
	}

	// リポジトリを介して更新
	err = uc.deviceRepo.Save(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return NewDeviceOutput(device), nil
}

// DeleteDevice 指定されたIDのDeviceを削除.
func (uc *deviceUsecase) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	// 削除対象の存在チェック
	device, err := uc.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find device by ID for deletion: %w", err) // FindByIDで発生したエラーを返す
	}

	if device == nil {
		return entity.ErrDeviceNotFound // デバイスが見つからない場合
	}

	// デバイスが存在する場合のみ削除を実行
	err = uc.deviceRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil
}
