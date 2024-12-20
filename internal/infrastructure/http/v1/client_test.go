package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetAccessToken(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		wantStatus int
	}{
		{
			name: "success",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":0,"message":"success","response":{"access_token":"test_token"}}`))
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "unauthorized",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.setupMock))
			defer server.Close()

			client := &Client{
				client:  &http.Client{},
				baseURL: server.URL,
			}

			resp, err := client.GetAccessToken()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, string(resp), "test_token")
			}
		})
	}
}

func TestClient_GetPaymentStatus(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		wantStatus int
	}{
		{
			name: "success",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":0,"message":"success","response":{"status":"paid"}}`))
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "unauthorized",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "not found",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.setupMock))
			defer server.Close()

			client := &Client{
				client:  &http.Client{},
				baseURL: server.URL,
			}

			resp, err := client.GetPaymentStatus("test_imp_uid", "test_token")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, string(resp), "paid")
			}
		})
	}
}

func TestClient_RequestCancelPayment(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		wantStatus int
	}{
		{
			name: "success",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":0,"message":"success","response":{"status":"cancelled"}}`))
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "unauthorized",
			setupMock: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.setupMock))
			defer server.Close()

			client := &Client{
				client:  &http.Client{},
				baseURL: server.URL,
			}

			resp, err := client.RequestCancelPayment("test_imp_uid", "test_token")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, string(resp), "cancelled")
			}
		})
	}
}
