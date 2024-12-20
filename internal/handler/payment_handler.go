package handler

import (
	"encoding/json"
	"net/http"

	_http "github.com/nerolizm/portone-payment/internal/infrastructure/http"
	"github.com/nerolizm/portone-payment/internal/model"
	"github.com/nerolizm/portone-payment/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PaymentHandlerInterface interface {
	HandlePaymentCancel(w http.ResponseWriter, r *http.Request)
}

type PaymentHandler struct {
	service service.PaymentServiceInterface
	logger  zerolog.Logger
}

func NewPaymentHandler(service service.PaymentServiceInterface) PaymentHandlerInterface {
	return &PaymentHandler{
		service: service,
		logger:  log.With().Str("handler", "payment").Logger(),
	}
}

func (h *PaymentHandler) HandlePaymentCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cancelReq model.CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&cancelReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug().Str("imp_uid", cancelReq.ImpUid).Msg("Payment cancel request received")

	response, err := h.service.CancelPayment(cancelReq.ImpUid)
	if err != nil {
		h.logger.Error().Err(err).Str("imp_uid", cancelReq.ImpUid).Msg("Failed to cancel payment")

		if httpErr, ok := err.(*_http.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
