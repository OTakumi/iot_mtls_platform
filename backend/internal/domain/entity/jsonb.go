package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONBMap は JSONB カラム用のカスタムマップ型です。
// map[string]any のエイリアスであり、sql.Scanner および driver.Valuer インターフェースを実装します。
type JSONBMap map[string]any

// Scan はデータベースの値をJSONBMapにスキャンします。
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

// Value はJSONBMapの値をデータベースに保存できる形式に変換します。
func (j *JSONBMap) Value() (driver.Value, error) {
	if j == nil || *j == nil {
		return nil, nil //nolint:nilnil
	}
	// 空のマップは空のJSONオブジェクトとして保存
	if len(*j) == 0 {
		return "{}", nil
	}

	bytes, err := json.Marshal(*j)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBMap: %w", err)
	}

	return bytes, nil
}
