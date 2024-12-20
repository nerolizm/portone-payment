package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nerolizm/portone-payment/internal/config"

	"github.com/nerolizm/portone-payment/internal/handler"
	v1 "github.com/nerolizm/portone-payment/internal/infrastructure/http/v1"
	"github.com/nerolizm/portone-payment/internal/model"
	"github.com/nerolizm/portone-payment/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestPaymentCancellation_Integration(t *testing.T) {
	currentDir, err := os.Getwd()
	assert.NoError(t, err)

	projectRoot := filepath.Join(currentDir, "../..")
	err = os.Chdir(projectRoot)
	assert.NoError(t, err)

	defer func() {
		err := os.Chdir(currentDir)
		assert.NoError(t, err)
	}()

	err = config.Init()
	assert.NoError(t, err)

	client := v1.NewClient()
	paymentService := service.NewPaymentService(client)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	tests := []struct {
		name       string
		impUid     string
		wantStatus int
		wantCode   int
		wantError  bool
	}{
		{
			name:       "non-existent payment",
			impUid:     "imp_test_123",
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := model.CancelRequest{
				ImpUid: tt.impUid,
			}
			jsonBody, err := json.Marshal(reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/cancel-payment", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			paymentHandler.HandlePaymentCancel(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
