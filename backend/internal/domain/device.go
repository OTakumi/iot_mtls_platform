package domain

import (
	"errors"
	"time"

	"backend/internal/domain/datetime"
	"backend/internal/domain/id"

	"github.com/google/uuid"
)

type Device struct {
	id          uuid.UUID
	hardware_id string
	name        string
	metadata    []string
	created_at  time.Time
}

func NewDevice(hardware_id string, name string, metadata []string) (Device, error) {
	if hardware_id == "" {
		return Device{}, errors.New("hardware id cannot be empty")
	}

	if name == "" {
		return Device{}, errors.New("name cannot be empty")
	}

	newDevice := Device{
		id:          id.NewID().Value,
		hardware_id: hardware_id,
		name:        name,
		metadata:    metadata,
		created_at:  datetime.Now().Value,
	}

	return newDevice, nil
}
