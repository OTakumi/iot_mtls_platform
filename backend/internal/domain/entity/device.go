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

	// HardwareID is a unique identifier for the device.
	// e.g., MAC address or serial number.
	HardwareID string `gorm:"uniqueIndex;not null;column:hardware_id"`

	// Name is an optional, user-friendly name for the device.
	Name string `gorm:"default:null"`

	// Metadata contains device-specific information, such as specifications,
	// installation location, and firmware version.
	// Its content is searchable, so it is stored in JSONB format.
	// Defaults to an empty JSON object '{}'.
	Metadata JSONBMap `gorm:"type:jsonb;default:'{}'"`

	// CreatedAt and UpdatedAt are automatically managed by GORM.
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
		Name:       "", // Default to an empty string, to be overwritten if a name is provided.
		Metadata:   newMetadata,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	if name != nil {
		newDevice.Name = *name
	}

	return newDevice, nil
}
