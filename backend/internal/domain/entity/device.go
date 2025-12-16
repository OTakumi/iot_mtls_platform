package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// デバイス情報はHardware_IDを必須とし、Name, MetadataはOptionとする
type Device struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// デバイス固有のIDを設定する
	// MACアドレスやデバイスIDなど
	HardwareID string `gorm:"uniqueIndex;not null;column:hardware_id"`

	// NameはNull許容
	Name string `gorm:"default:null"`

	// metadataには、デバイスの仕様、設置場所、デバイスバージョンなど、
	// デバイスによって異なる情報が入力される
	// このmetadataの内容を検索対象とすることも考えられるため、json形式で保存できるようにする
	// デフォルト値は空のjson {}
	Metadata map[string]any `gorm:"type:jsonb;default:'{}'"`

	// CreatedAt, UpdatedAtという命名にすることでGORMが自動でタイムスタンプを追加する
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDevice(hardwareId string, name *string, metadata map[string]interface{}) (*Device, error) {
	if hardwareId == "" {
		return nil, errors.New("hardware id cannot be empty")
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	newDevice := &Device{
		HardwareID: hardwareId,
		Metadata:   metadata,
	}

	if name != nil {
		newDevice.Name = *name
	}

	return newDevice, nil
}
