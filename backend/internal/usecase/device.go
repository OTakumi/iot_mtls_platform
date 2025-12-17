package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
)

// エンティティに関するアプリケーション固有のビジネスロジックを定義するインターフェース
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

// DeviceUsecaseインターフェースの実装
type deviceUsecase struct {
	deviceRepo repository.DeviceRepository
}

// deviceUsecaseの新しいインスタンスを生成
func NewDeviceUsecase(repo repository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{deviceRepo: repo}
}

// 新しいDeviceを登録
func (uc *deviceUsecase) CreateDevice(ctx context.Context, input CreateDeviceInput) (*DeviceOutput, error) {
	// ドメインエンティティの生成
	// IDはDB側で生成されるため、ここではuuid.Nilで問題ない
	var namePtr *string
	if input.Name != "" {
		namePtr = &input.Name
	}
	device, err := entity.NewDevice(input.HardwareID, namePtr, input.Metadata)
	if err != nil {
		return nil, err // エンティティ生成時のバリデーションエラーなど
	}

	// リポジトリを介して保存
	if err := uc.deviceRepo.Save(ctx, device); err != nil {
		return nil, err // DB保存時のエラー
	}

	return NewDeviceOutput(device), nil // 出力DTOに変換して返す
}

// 指定されたIDのDeviceを取得
func (uc *deviceUsecase) GetDevice(ctx context.Context, id uuid.UUID) (*DeviceOutput, error) {
	device, err := uc.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found") // デバイスが見つからない場合
	}
	return NewDeviceOutput(device), nil
}

// 全てのDeviceを取得
func (uc *deviceUsecase) ListDevices(ctx context.Context) ([]*DeviceOutput, error) {
	devices, err := uc.deviceRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var outputs []*DeviceOutput
	for _, device := range devices {
		outputs = append(outputs, NewDeviceOutput(device))
	}
	return outputs, nil
}

// 既存のDeviceを更新
func (uc *deviceUsecase) UpdateDevice(ctx context.Context, input UpdateDeviceInput) (*DeviceOutput, error) {
	// 更新対象のDeviceを検索
	device, err := uc.deviceRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found")
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
	if err := uc.deviceRepo.Save(ctx, device); err != nil {
		return nil, err
	}

	return NewDeviceOutput(device), nil
}

// 指定されたIDのDeviceを削除
func (uc *deviceUsecase) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	// 削除対象の存在チェック（任意、リポジトリ層でNotFoundエラーを返さないなら必須）
	// ここでは存在チェックなしで直接削除を試みる
	if err := uc.deviceRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
