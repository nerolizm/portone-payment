package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nerolizm/portone-payment/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CancelPayment(impUid string) ([]byte, error) {
	args := m.Called(impUid)
	return args.Get(0).([]byte), args.Error(1)
}

func TestPaymentHandler_HandlePaymentCancel(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		impUid       string
		setupMock    func(*MockPaymentService)
		wantStatus   int
		wantCode     int
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:   "successful cancellation",
			method: http.MethodPost,
			impUid: "imp_test_123",
			setupMock: func(m *MockPaymentService) {
				m.On("CancelPayment", "imp_test_123").Return(
					[]byte(`{"code":0,"message":"success","response":{"status":"cancelled"}}`),
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantCode:   0,
		},
		{
			name:   "method not allowed",
			method: http.MethodGet,
			impUid: "imp_test_123",
			setupMock: func(m *MockPaymentService) {
				// No mock setup needed for this case
			},
			wantStatus:   http.StatusMethodNotAllowed,
			wantError:    true,
			wantErrorMsg: "Method not allowed",
		},
		{
			name:   "invalid request body - empty imp_uid",
			method: http.MethodPost,
			impUid: "",
			setupMock: func(m *MockPaymentService) {
				m.On("CancelPayment", "").Return(
					[]byte(`{"code":-1,"message":"imp_uid is empty"}`),
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantCode:   -1,
			wantError:  true,
		},
		{
			name:   "service error",
			method: http.MethodPost,
			impUid: "imp_test_123",
			setupMock: func(m *MockPaymentService) {
				m.On("CancelPayment", "imp_test_123").Return(
					[]byte{},
					fmt.Errorf("service error"),
				)
			},
			wantStatus:   http.StatusInternalServerError,
			wantError:    true,
			wantErrorMsg: "service error",
		},
		{
			name:   "payment not found",
			method: http.MethodPost,
			impUid: "imp_test_123",
			setupMock: func(m *MockPaymentService) {
				m.On("CancelPayment", "imp_test_123").Return(
					[]byte(`{"code":-1,"message":"Payment not found"}`),
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantCode:   -1,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentService)
			tt.setupMock(mockService)
			handler := NewPaymentHandler(mockService)

			// Create request
			reqBody := model.CancelRequest{
				ImpUid: tt.impUid,
			}
			jsonBody, err := json.Marshal(reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(tt.method, "/cancel-payment", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handler.HandlePaymentCancel(w, req)

			// Assert status code
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError {
				if tt.wantErrorMsg != "" {
					assert.Contains(t, w.Body.String(), tt.wantErrorMsg)
				}
				if tt.wantCode != 0 {
					var response model.Response
					err = json.NewDecoder(w.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantCode, response.Code)
				}
			} else {
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
				var response model.Response
				err = json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCode, response.Code)
			}

			// Verify all mocked calls were made
			mockService.AssertExpectations(t)
		})
	}
}
