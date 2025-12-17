package entity

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// TestNewDeviceはNewDeviceコンストラクタ関数の様々な挙動をテストする
func TestNewDevice(t *testing.T) {
	// args構造体はNewDevice関数に渡す引数をまとめる
	type args struct {
		hardwareID string
		name       *string
		metadata   map[string]any
	}

	// want構造体は期待されるDeviceの状態を定義する
	type want struct {
		hardwareID string
		name       string
		metadata   JSONBMap
		id         uuid.UUID
	}

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

	tests := []struct {
		name       string // テストケース名
		desc       string // 詳細な意図や観点
		args       args
		want       want
		wantErr    bool
		wantErrMsg string // 期待されるエラーメッセージ
	}{
		{
			name: "success: with name and metadata",
			desc: "全ての必須およびオプション引数が適切に提供された場合の正常系テスト",
			args: args{
				hardwareID: "test-hw-id-1",
				name:       &nameStr,
				metadata:   map[string]any{"key": "value"},
			},
			want: want{
				hardwareID: "test-hw-id-1",
				name:       "test-device-1",
				metadata:   JSONBMap{"key": "value"},
				id:         uuid.Nil,
			},
			wantErr: false,
		},
		{
			name: "success: with nil name and nil metadata",
			desc: "オプション引数がnilで提供された場合の正常系テスト",
			args: args{
				hardwareID: "test-hw-id-2",
				name:       nil,
				metadata:   nil,
			},
			want: want{
				hardwareID: "test-hw-id-2",
				name:       "",
				metadata:   JSONBMap{},
				id:         uuid.Nil,
			},
			wantErr: false,
		},
		{
			name: "success: with complex metadata",
			desc: "複雑なメタデータが提供された場合の正常系テスト",
			args: args{
				hardwareID: "test-hw-id-3",
				name:       &nameStrComplex,
				metadata:   complexMetadata,
			},
			want: want{
				hardwareID: "test-hw-id-3",
				name:       "env-sensor-1",
				metadata:   JSONBMap(complexMetadata),
				id:         uuid.Nil,
			},
			wantErr: false,
		},
		{
			name: "failure: empty hardware_id",
			desc: "必須のhardwareIDが空文字列の場合の異常系テスト",
			args: args{
				hardwareID: "",
				name:       &nameStrEmptyID,
				metadata:   nil,
			},
			want:       want{}, // エラーケースなのでwantは重要ではない
			wantErr:    true,
			wantErrMsg: "hardware id cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := NewDevice(tt.args.hardwareID, tt.args.name, tt.args.metadata)

			// Assert - Error
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("NewDevice() error msg = %q, wantErrMsg %q", err.Error(), tt.wantErrMsg)
				}
				if got != nil {
					t.Errorf("NewDevice() got = %v, want nil for error case", got)
				}
				return // エラーケースのテストはここで終了
			}

			// Assert - Success (got is not nil)
			if got == nil {
				t.Fatal("NewDevice() returned a nil device for a success case")
			}

			// Assert - Success (Fields)
			if got.HardwareID != tt.want.hardwareID {
				t.Errorf("NewDevice() HardwareID = %v, want %v", got.HardwareID, tt.want.hardwareID)
			}
			if got.Name != tt.want.name {
				t.Errorf("NewDevice() Name = %v, want %v", got.Name, tt.want.name)
			}
			if !reflect.DeepEqual(got.Metadata, tt.want.metadata) {
				t.Errorf("NewDevice() Metadata = %v, want %v", got.Metadata, tt.want.metadata)
			}
			if got.ID != tt.want.id {
				t.Errorf("NewDevice() ID = %v, want %v", got.ID, tt.want.id)
			}
		})
	}
}
