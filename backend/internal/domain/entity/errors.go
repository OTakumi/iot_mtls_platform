package entity

import "errors"

var (
	// ErrHardwareIDEmpty is returned when a hardware ID is empty.
	ErrHardwareIDEmpty = errors.New("hardware id cannot be empty")
	// ErrUnsupportedTypeForJSONBMapScan is returned when an unsupported type is used for scanning a JSONBMap.
	ErrUnsupportedTypeForJSONBMapScan = errors.New("unsupported type for JSONBMap Scan")
	ErrDeviceNotFound                 = errors.New("device not found")
)
