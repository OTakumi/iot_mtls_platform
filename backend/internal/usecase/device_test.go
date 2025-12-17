package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"backend/internal/domain/entity"
	"backend/internal/usecase"
)

// DeviceRepositoryインターフェースのモック実装
type MockDeviceRepository struct {
	mock.Mock
}

// モックされたSaveメソッド
func (m *MockDeviceRepository) Save(ctx context.Context, device *entity.Device) error {
	args := m.Called(ctx, device)
	return args.Error(0)
}

// モックされたFindByIDメソッド
func (m *MockDeviceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Device), args.Error(1)
}

// モックされたFindByHardwareIDメソッド
func (m *MockDeviceRepository) FindByHardwareID(ctx context.Context, hardwareID string) (*entity.Device, error) {
	args := m.Called(ctx, hardwareID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Device), args.Error(1)
}

// モックされたFindAllメソッド
func (m *MockDeviceRepository) FindAll(ctx context.Context) ([]*entity.Device, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Device), args.Error(1)
}

// モックされたDeleteメソッド
func (m *MockDeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateDeviceメソッドのテスト
func TestCreateDevice(t *testing.T) {
	ctx := context.Background()

	type args struct {
		input usecase.CreateDeviceInput
	}
	type want struct {
		output *usecase.DeviceOutput
	}
	type mockSetup func(repo *MockDeviceRepository)

	tests := []struct {
		name       string
		desc       string
		args       args
		setupMock  mockSetup
		want       want
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系: 新しいデバイスを作成する",
			desc: "正常な入力で新しいデバイスが作成され、対応する出力DTOが返されることを確認する",
			args: args{
				input: usecase.CreateDeviceInput{
					HardwareID: "hw-create-001",
					Name:       "テストデバイス C",
					Metadata:   map[string]any{"os": "linux"},
				},
			},
			setupMock: func(repo *MockDeviceRepository) {
				// モックのSaveはIDを変更すべきではないため、Run関数を削除
				repo.On("Save", ctx, mock.AnythingOfType("*entity.Device")).
					Return(nil).
					Once()
			},
			want: want{
				output: &usecase.DeviceOutput{
					HardwareID: "hw-create-001",
					Name:       "テストデバイス C",
					Metadata:   map[string]any{"os": "linux"},
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: リポジトリがエラーを返す",
			desc: "リポジトリのSaveメソッドがエラーを返した場合、ユースケースもエラーを返すことを確認する",
			args: args{
				input: usecase.CreateDeviceInput{
					HardwareID: "hw-err-001",
					Name:       "エラーデバイス",
				},
			},
			setupMock: func(repo *MockDeviceRepository) {
				repo.On("Save", ctx, mock.AnythingOfType("*entity.Device")).Return(errors.New("repository save error")).Once()
			},
			want:       want{output: nil},
			wantErr:    true,
			wantErrMsg: "repository save error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDeviceRepository)
			uc := usecase.NewDeviceUsecase(mockRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			got, err := uc.CreateDevice(ctx, tt.args.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Equal(t, tt.want.output, got)
				mockRepo.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)

			assert.Equal(t, tt.want.output.HardwareID, got.HardwareID)
			assert.Equal(t, tt.want.output.Name, got.Name)
			assert.Equal(t, tt.want.output.Metadata, got.Metadata)
			assert.Equal(t, uuid.Nil, got.ID) // IDはDB側で生成されるため、ユースケース層が返すIDはuuid.Nilであることを期待

			mockRepo.AssertExpectations(t)
		})
	}
}

// GetDeviceメソッドのテスト
func TestGetDevice(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	type args struct {
		id uuid.UUID
	}
	type want struct {
		output *usecase.DeviceOutput
	}
	type mockSetup func(repo *MockDeviceRepository, args args)

	tests := []struct {
		name       string
		desc       string
		args       args
		setupMock  mockSetup
		want       want
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系: 既存のデバイスを取得する",
			desc: "存在するIDでデバイスを正常に取得できることを確認する",
			args: args{id: testID},
			setupMock: func(repo *MockDeviceRepository, args args) {
				repo.On("FindByID", ctx, args.id).Return(&entity.Device{
					ID:         args.id,
					HardwareID: "hw-get-001",
					Name:       "取得デバイス",
				}, nil).Once()
			},
			want: want{
				output: &usecase.DeviceOutput{
					ID:         testID,
					HardwareID: "hw-get-001",
					Name:       "取得デバイス",
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: デバイスが見つからない",
			desc: "存在しないIDで検索した場合に、'device not found'エラーが返されることを確認する",
			args: args{id: uuid.New()},
			setupMock: func(repo *MockDeviceRepository, args args) {
				repo.On("FindByID", ctx, args.id).Return(nil, nil).Once()
			},
			want:       want{output: nil},
			wantErr:    true,
			wantErrMsg: "device not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDeviceRepository)
			uc := usecase.NewDeviceUsecase(mockRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, tt.args)
			}

			got, err := uc.GetDevice(ctx, tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			// ID以外はwantと一致し、CreatedAt/UpdatedAtは無視する
			if tt.want.output != nil && got != nil {
				tt.want.output.CreatedAt = got.CreatedAt
				tt.want.output.UpdatedAt = got.UpdatedAt
			}

			assert.Equal(t, tt.want.output, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

// ListDevicesメソッドのテスト
func TestListDevices(t *testing.T) {
	ctx := context.Background()

	type want struct {
		output []*usecase.DeviceOutput
	}
	type mockSetup func(repo *MockDeviceRepository)

	// テストデータ
	device1 := &entity.Device{ID: uuid.New(), HardwareID: "hw-list-001", Name: "リストデバイス 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	device2 := &entity.Device{ID: uuid.New(), HardwareID: "hw-list-002", Name: "リストデバイス 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	tests := []struct {
		name       string
		desc       string
		setupMock  mockSetup
		want       want
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系: 全てのデバイスをリストする",
			desc: "複数のデバイスが存在する場合に、それらが全て正常にリストアップされることを確認する",
			setupMock: func(repo *MockDeviceRepository) {
				repo.On("FindAll", ctx).Return([]*entity.Device{device1, device2}, nil).Once()
			},
			want: want{
				output: []*usecase.DeviceOutput{
					usecase.NewDeviceOutput(device1),
					usecase.NewDeviceOutput(device2),
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: リポジトリがエラーを返す",
			desc: "リポジトリのFindAllメソッドがエラーを返した場合、ユースケースもエラーを返すことを確認する",
			setupMock: func(repo *MockDeviceRepository) {
				repo.On("FindAll", ctx).Return(nil, errors.New("db find all error")).Once()
			},
			want:       want{output: nil},
			wantErr:    true,
			wantErrMsg: "db find all error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDeviceRepository)
			uc := usecase.NewDeviceUsecase(mockRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			got, err := uc.ListDevices(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.output, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

// UpdateDeviceメソッドのテスト
func TestUpdateDevice(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()
	updatedName := "更新された名前"

	type args struct {
		input usecase.UpdateDeviceInput
	}
	type want struct {
		output *usecase.DeviceOutput
	}
	type mockSetup func(repo *MockDeviceRepository, args args)

	tests := []struct {
		name       string
		desc       string
		args       args
		setupMock  mockSetup
		want       want
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系: デバイス名とメタデータを更新する",
			desc: "既存のデバイスの名前とメタデータが正常に更新されることを確認する",
			args: args{
				input: usecase.UpdateDeviceInput{
					ID:       testID,
					Name:     &updatedName,
					Metadata: map[string]any{"status": "active"},
				},
			},
			setupMock: func(repo *MockDeviceRepository, args args) {
				existingDevice := &entity.Device{ID: args.input.ID, HardwareID: "hw-update-001"}
				repo.On("FindByID", ctx, args.input.ID).Return(existingDevice, nil).Once()
				repo.On("Save", ctx, mock.AnythingOfType("*entity.Device")).Return(nil).Once()
			},
			want: want{
				output: &usecase.DeviceOutput{
					ID:         testID,
					HardwareID: "hw-update-001",
					Name:       updatedName,
					Metadata:   map[string]any{"status": "active"},
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: 更新対象のデバイスが見つからない",
			desc: "更新対象のデバイスIDが存在しない場合に'device not found'エラーが返されることを確認する",
			args: args{input: usecase.UpdateDeviceInput{ID: uuid.New()}},
			setupMock: func(repo *MockDeviceRepository, args args) {
				repo.On("FindByID", ctx, args.input.ID).Return(nil, nil).Once()
			},
			want:       want{output: nil},
			wantErr:    true,
			wantErrMsg: "device not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDeviceRepository)
			uc := usecase.NewDeviceUsecase(mockRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, tt.args)
			}

			got, err := uc.UpdateDevice(ctx, tt.args.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			if tt.want.output != nil && got != nil {
				tt.want.output.CreatedAt = got.CreatedAt
				tt.want.output.UpdatedAt = got.UpdatedAt
			}

			assert.Equal(t, tt.want.output, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

// DeleteDeviceメソッドのテスト
func TestDeleteDevice(t *testing.T) {
	ctx := context.Background()

	testID := uuid.New()

	type args struct {
		id uuid.UUID
	}
	type mockSetup func(repo *MockDeviceRepository, args args)

	tests := []struct {
		name       string
		desc       string
		args       args
		setupMock  mockSetup
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系: デバイスを削除する",
			desc: "存在するデバイスIDを指定して、デバイスが正常に削除されることを確認する",
			args: args{id: testID}, // testID を使用
			setupMock: func(repo *MockDeviceRepository, args args) {
				// FindByIDが成功してデバイスを返すようにモック
				repo.On("FindByID", ctx, args.id).Return(&entity.Device{ID: args.id}, nil).Once()
				// Deleteが成功するようにモック
				repo.On("Delete", ctx, args.id).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "異常系: 削除対象のデバイスが見つからない",
			desc: "デバイスが見つからなかった場合に'device not found'エラーが返されることを確認する",
			args: args{id: uuid.New()}, // 新しいユニークなIDを使用
			setupMock: func(repo *MockDeviceRepository, args args) {
				// FindByIDがnilとnilエラーを返すようにモック (デバイスが見つからない場合)
				repo.On("FindByID", ctx, args.id).Return(nil, nil).Once()
			},
			wantErr:    true,
			wantErrMsg: "device not found",
		},
		{
			name: "異常系: FindByIDがエラーを返す",
			desc: "FindByIDメソッドがエラーを返した場合、ユースケースもそのエラーを返すことを確認する",
			args: args{id: uuid.New()}, // 新しいユニークなIDを使用
			setupMock: func(repo *MockDeviceRepository, args args) {
				// FindByIDがエラーを返すようにモック
				repo.On("FindByID", ctx, args.id).Return(nil, errors.New("FindByID db error")).Once()
			},
			wantErr:    true,
			wantErrMsg: "FindByID db error",
		},
		{
			name: "異常系: Deleteがエラーを返す",
			desc: "リポジトリのDeleteメソッドがエラーを返した場合、ユースケースもエラーを返すことを確認する",
			args: args{id: testID},
			setupMock: func(repo *MockDeviceRepository, args args) {
				// FindByIDが成功してデバイスを返すようにモック
				repo.On("FindByID", ctx, args.id).Return(&entity.Device{ID: args.id}, nil).Once()
				// Deleteがエラーを返すようにモック
				repo.On("Delete", ctx, args.id).Return(errors.New("db delete error")).Once()
			},
			wantErr:    true,
			wantErrMsg: "db delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDeviceRepository)
			uc := usecase.NewDeviceUsecase(mockRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, tt.args)
			}

			err := uc.DeleteDevice(ctx, tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
