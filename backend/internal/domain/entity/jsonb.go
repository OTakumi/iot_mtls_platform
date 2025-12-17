package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONBMap は JSONB カラム用のカスタムマップ型です。
// map[string]any のエイリアスであり、sql.Scanner および driver.Valuer インターフェースを実装します。
type JSONBMap map[string]any

// Scan はデータベースの値をJSONBMapにスキャンします。
func (j *JSONBMap) Scan(value interface{}) error {
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
		return errors.New("unsupported type for JSONBMap Scan")
	}

	if len(source) == 0 {
		*j = make(JSONBMap)
		return nil
	}

	return json.Unmarshal(source, j)
}

// Value はJSONBMapの値をデータベースに保存できる形式に変換します。
func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	// 空のマップは空のJSONオブジェクトとして保存
	if len(j) == 0 {
		return "{}", nil
	}

	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("JSONBMapのMarshalに失敗: %w", err)
	}
	return bytes, nil
}
