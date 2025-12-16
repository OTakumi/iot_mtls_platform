package entity

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// TestNewDeviceはNewDeviceコンストラクタ関数の様々な挙動をテストする
func TestNewDevice(t *testing.T) {
	// success - with name and metadata: 全ての必須およびオプション引数が適切に提供された場合の正常系テスト
	t.Run("success - with name and metadata", func(t *testing.T) {
		hardwareID := "test-hw-id-1"
		nameStr := "test-device-1"
		// ポインタとして名前を渡す
		name := &nameStr
		// シンプルなメタデータを渡す
		metadata := map[string]any{"key": "value"}

		// NewDeviceを呼び出し、デバイスとエラーを生成
		device, err := NewDevice(hardwareID, name, metadata)
		// 関数がエラーを返さないことを確認
		if err != nil {
			t.Fatalf("NewDevice() returned an unexpected error: %v", err)
		}

		// Deviceインスタンスがnilではないこと
		if device == nil {
			t.Fatal("NewDevice() returned a nil device")
		}

		// HardwareIDフィールドが入力値と一致すること
		if device.HardwareID != hardwareID {
			t.Errorf("expected HardwareID %q, got %q", hardwareID, device.HardwareID)
		}

		// Nameフィールドが入力値（ポインタが指す文字列）と一致すること
		if device.Name != *name {
			t.Errorf("expected Name %q, got %q", *name, device.Name)
		}

		// Metadataフィールドが入力マップとディープイコールであること
		if !reflect.DeepEqual(device.Metadata, metadata) {
			t.Errorf("expected Metadata %+v, got %+v", metadata, device.Metadata)
		}

		// IDフィールドが、仕様通りuuid.Nil（ゼロ値）であること
		// DB永続化まではIDが生成されないという仕様変更を反映
		if device.ID != uuid.Nil {
			t.Errorf("expected ID to be nil UUID, but got %q", device.ID.String())
		}
	})

	// success - with nil name and nil metadata
	// オプション引数がnilで提供された場合の正常系テスト
	t.Run("success - with nil name and nil metadata", func(t *testing.T) {
		hardwareID := "test-hw-id-2"

		// nameとmetadataをnilでNewDeviceを呼び出す
		device, err := NewDevice(hardwareID, nil, nil)
		// 関数がエラーを返さないこと
		if err != nil {
			t.Fatalf("NewDevice() returned an unexpected error: %v", err)
		}
		// Deviceインスタンスがnilではないこと
		if device == nil {
			t.Fatal("NewDevice() returned a nil device")
		}

		// Nameフィールドがnil入力に対して空文字列("")であること
		if device.Name != "" {
			t.Errorf("expected Name to be an empty string for nil input, got %q", device.Name)
		}
		// Metadataフィールドがnil入力に対してnilではなく、空のマップとして初期化されていること
		if device.Metadata == nil {
			t.Error("expected Metadata to be an empty map for nil input, but it was nil")
		}
		// Metadataフィールドが空のマップ（要素数0）であること
		if len(device.Metadata) != 0 {
			t.Errorf("expected Metadata to be an empty map, but it has %d items", len(device.Metadata))
		}
		// IDフィールドがuuid.Nilであること
		// DB永続化まではIDが生成されない
		if device.ID != uuid.Nil {
			t.Errorf("expected ID to be nil UUID, but got %q", device.ID.String())
		}
	})

	// success - with complex metadata
	// 複雑なメタデータが提供された場合の正常系テスト
	t.Run("success - with complex metadata", func(t *testing.T) {
		hardwareID := "test-hw-id-3"
		nameStr := "env-sensor-1"
		// ネストされたマップや異なる型の値を含む複雑なメタデータを定義
		metadata := map[string]any{
			"type": "env_sensor",
			"hardware": map[string]any{
				"model":        "Raspberry Pi 4B",
				"revision":     "1.2",
				"manufacturer": "Sony UK",
			},
			"location": map[string]any{
				"building": "Factory-A",
				"floor":    float64(2), // JSONの数値はGoではfloat64としてパースされる
				"zone":     "shipping_area",
			},
			"firmware": map[string]any{
				"version":     "2.4.1",
				"last_update": "2024-12-01T10:00:00Z",
			},
			"config": map[string]any{
				"sync_interval_sec":    float64(60),
				"alert_threshold_temp": 40.0,
			},
		}

		// NewDeviceを呼び出し
		device, err := NewDevice(hardwareID, &nameStr, metadata)
		// 関数がエラーを返さないこと
		if err != nil {
			t.Fatalf("NewDevice() returned an unexpected error: %v", err)
		}

		// 生成されたDevice.Metadataが入力された複雑なマップと完全に一致すること
		if !reflect.DeepEqual(device.Metadata, metadata) {
			t.Errorf("Metadata is not set correctly for complex data.\nGot:  %+v\nWant: %+v", device.Metadata, metadata)
		}
	})

	// failure - empty hardware_id
	// 必須のhardwareIDが空文字列の場合の異常系テスト
	t.Run("failure - empty hardware_id", func(t *testing.T) {
		nameStr := "some-name"
		expectedError := "hardware id cannot be empty" // 期待されるエラーメッセージ

		// hardwareIDを空文字列でNewDeviceを呼び出す
		device, err := NewDevice("", &nameStr, nil)

		// 関数がエラーを返すこと
		if err == nil {
			t.Fatal("NewDevice() did not return an error for empty hardware_id")
		}
		// 無効な入力に対してDeviceインスタンスがnilであること
		if device != nil {
			t.Errorf("NewDevice() returned a non-nil device for an invalid call: %+v", device)
		}
		// 返されたエラーメッセージが期待されるエラーメッセージと正確に一致すること
		if err.Error() != expectedError {
			t.Errorf("NewDevice() returned wrong error message. got %q, want %q", err.Error(), expectedError)
		}
	})
}

