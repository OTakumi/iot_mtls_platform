// Package entity defines the domain entities.
package entity

import (
	"maps"
	"time"

	"github.com/google/uuid"
)

// Device represents a device entity, with HardwareID being mandatory and Name and Metadata being optional.
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
	Metadata JSONBMap `gorm:"type:jsonb;default:'{}'"`

	// CreatedAt, UpdatedAtという命名にすることでGORMが自動でタイムスタンプを追加する
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDevice creates a new Device.
func NewDevice(hardwareID string, name *string, metadata map[string]any) (*Device, error) {
	if hardwareID == "" {
		return nil, ErrHardwareIDEmpty
	}

	newMetadata := make(JSONBMap)
	maps.Copy(newMetadata, metadata)

	newDevice := &Device{
		ID:         uuid.Nil,
		HardwareID: hardwareID,
		Name:       "", // Default to empty string, to be overwritten if name is provided
		Metadata:   newMetadata,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	if name != nil {
		newDevice.Name = *name
	}

	return newDevice, nil
}
