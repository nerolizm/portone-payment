package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetAccessToken() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockClient) GetPaymentStatus(impUid, token string) ([]byte, error) {
	args := m.Called(impUid, token)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockClient) RequestCancelPayment(impUid, token string) ([]byte, error) {
	args := m.Called(impUid, token)
	return args.Get(0).([]byte), args.Error(1)
}

func TestPaymentService_CancelPayment(t *testing.T) {
	mockClient := new(MockClient)
	service := NewPaymentService(mockClient)

	tests := []struct {
		name    string
		impUid  string
		setup   func()
		want    []byte
		wantErr bool
	}{
		{
			name:   "successful cancellation - paid status",
			impUid: "test_imp_uid",
			setup: func() {
				mockClient.On("GetAccessToken").Return(
					[]byte(`{"response":{"access_token":"test_token"}}`), nil,
				)
				mockClient.On("GetPaymentStatus", "test_imp_uid", "test_token").Return(
					[]byte(`{"response":{"status":"paid"}}`), nil,
				)
				mockClient.On("RequestCancelPayment", "test_imp_uid", "test_token").Return(
					[]byte(`{"response":{"status":"cancelled"}}`), nil,
				)
			},
			want:    []byte(`{"response":{"status":"cancelled"}}`),
			wantErr: false,
		},
		{
			name:   "no cancellation - not paid status",
			impUid: "test_imp_uid",
			setup: func() {
				mockClient.On("GetAccessToken").Return(
					[]byte(`{"response":{"access_token":"test_token"}}`), nil,
				)
				mockClient.On("GetPaymentStatus", "test_imp_uid", "test_token").Return(
					[]byte(`{"response":{"status":"ready"}}`), nil,
				)
			},
			want:    []byte(`{"response":{"status":"ready"}}`),
			wantErr: false,
		},
		{
			name:   "error - GetAccessToken returns 401",
			impUid: "test_imp_uid",
			setup: func() {
				mockClient.On("GetAccessToken").Return(
					[]byte(`{"code":401,"message":"Invalid API credentials"}`),
					fmt.Errorf("unauthorized access"),
				)
			},
			wantErr: true,
		},
		{
			name:   "error - GetPaymentStatus returns 401",
			impUid: "test_imp_uid",
			setup: func() {
				mockClient.On("GetAccessToken").Return(
					[]byte(`{"response":{"access_token":"test_token"}}`), nil,
				)
				mockClient.On("GetPaymentStatus", "test_imp_uid", "test_token").Return(
					[]byte(`{"code":401,"message":"Invalid access token"}`),
					fmt.Errorf("unauthorized access"),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			tt.setup()

			got, err := service.CancelPayment(tt.impUid)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, string(tt.want), string(got))
			mockClient.AssertExpectations(t)
		})
	}
}
