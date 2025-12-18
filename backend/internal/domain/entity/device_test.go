package entity_test

import (
	"reflect"
	"testing"

	"backend/internal/domain/entity"

	"github.com/google/uuid"
)

// want構造体は期待されるDeviceの状態を定義する.
type want struct {
	hardwareID string
	name       string
	metadata   entity.JSONBMap
	id         uuid.UUID
}

// TestNewDeviceはNewDeviceコンストラクタ関数の様々な挙動をテストする.
func TestNewDevice(t *testing.T) {
	t.Parallel()

	tests := getNewDeviceTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := entity.NewDevice(tt.args.hardwareID, tt.args.name, tt.args.metadata)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewDevice() expected error, but got nil")
				}

				if err.Error() != tt.wantErrMsg {
					t.Errorf("NewDevice() error msg = %q, wantErrMsg %q", err.Error(), tt.wantErrMsg)
				}

				if got != nil {
					t.Errorf("NewDevice() got = %v, want nil for error case", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("NewDevice() unexpected error: %v", err)
			}

			assertDevice(t, tt.want, got)
		})
	}
}

type newDeviceTestArgs struct {
	hardwareID string
	name       *string
	metadata   map[string]any
}

func getNewDeviceTestCases() []struct {
	name       string // テストケース名
	desc       string // 詳細な意図や観点
	args       newDeviceTestArgs
	want       want
	wantErr    bool
	wantErrMsg string // 期待されるエラーメッセージ
} {
	// テストケースで使用する変数を事前に定義
	nameStr := "test-device-1"
	nameStrComplex := "env-sensor-1"
	nameStrEmptyID := "some-name"

	// 複雑なメタデータの定義
	complexMetadata := map[string]any{
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

	return []struct {
		name       string // テストケース名
		desc       string // 詳細な意図や観点
		args       newDeviceTestArgs
		want       want
		wantErr    bool
		wantErrMsg string // 期待されるエラーメッセージ
	}{
		{
			name: "success: with name and metadata",
			desc: "Successful test case when all mandatory and optional arguments are provided correctly.",
			args: newDeviceTestArgs{
				hardwareID: "test-hw-id-1",
				name:       &nameStr,
				metadata:   map[string]any{"key": "value"},
			},
			want: want{
				hardwareID: "test-hw-id-1",
				name:       "test-device-1",
				metadata:   entity.JSONBMap{"key": "value"},
				id:         uuid.Nil,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "success: with nil name and nil metadata",
			desc: "Successful test case when optional arguments are provided as nil.",
			args: newDeviceTestArgs{
				hardwareID: "test-hw-id-2",
				name:       nil,
				metadata:   nil,
			},
			want: want{
				hardwareID: "test-hw-id-2",
				name:       "",
				metadata:   entity.JSONBMap{},
				id:         uuid.Nil,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "success: with complex metadata",
			desc: "Successful test case when complex metadata is provided.",
			args: newDeviceTestArgs{
				hardwareID: "test-hw-id-3",
				name:       &nameStrComplex,
				metadata:   complexMetadata,
			},
			want: want{
				hardwareID: "test-hw-id-3",
				name:       "env-sensor-1",
				metadata:   entity.JSONBMap(complexMetadata),
				id:         uuid.Nil,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "failure: empty hardware_id",
			desc: "Failure test case when mandatory hardwareID is an empty string.",
			args: newDeviceTestArgs{
				hardwareID: "",
				name:       &nameStrEmptyID,
				metadata:   nil,
			},
			want:       want{hardwareID: "", name: "", metadata: nil, id: uuid.Nil}, // want is not important for error cases
			wantErr:    true,
			wantErrMsg: "hardware id cannot be empty",
		},
	}
}

func assertDevice(t *testing.T, want want, got *entity.Device) {
	t.Helper()

	if got == nil {
		t.Fatal("NewDevice() returned a nil device for a success case")
	}

	if got.HardwareID != want.hardwareID {
		t.Errorf("NewDevice() HardwareID = %v, want %v", got.HardwareID, want.hardwareID)
	}

	if got.Name != want.name {
		t.Errorf("NewDevice() Name = %v, want %v", got.Name, want.name)
	}

	if !reflect.DeepEqual(got.Metadata, want.metadata) {
		t.Errorf("NewDevice() Metadata = %v, want %v", got.Metadata, want.metadata)
	}

	if got.ID != want.id {
		t.Errorf("NewDevice() ID = %v, want %v", got.ID, want.id)
	}
}
