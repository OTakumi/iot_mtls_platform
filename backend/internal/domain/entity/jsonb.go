package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONBMap is a custom type for `map[string]any` to handle JSONB database columns.
type JSONBMap map[string]any

// Scan implements the sql.Scanner interface, allowing the type to read from a database.
func (j *JSONBMap) Scan(value any) error {
	if value == nil {
		*j = make(JSONBMap)

		return nil
	}

	var source []byte

	switch v := value.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return ErrUnsupportedTypeForJSONBMapScan
	}

	if len(source) == 0 {
		*j = make(JSONBMap)

		return nil
	}

	err := json.Unmarshal(source, j)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONBMap: %w", err)
	}

	return nil
}

// Value implements the driver.Valuer interface, allowing the type to be written to a database.
func (j *JSONBMap) Value() (driver.Value, error) {
	if j == nil || *j == nil {
		return nil, nil //nolint:nilnil
	}
	// Empty maps are marshaled as empty JSON objects.
	if len(*j) == 0 {
		return "{}", nil
	}

	bytes, err := json.Marshal(*j)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBMap: %w", err)
	}

	return bytes, nil
}
