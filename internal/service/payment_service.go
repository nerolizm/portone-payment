package service

import (
	"encoding/json"

	v1 "github.com/nerolizm/portone-payment/internal/infrastructure/http/v1"
	"github.com/nerolizm/portone-payment/internal/model"
)

type PaymentServiceInterface interface {
	CancelPayment(impUid string) ([]byte, error)
}

type PaymentService struct {
	client v1.PaymentClientInterface
}

func NewPaymentService(client v1.PaymentClientInterface) PaymentServiceInterface {
	return &PaymentService{
		client: client,
	}
}

func (s *PaymentService) CancelPayment(impUid string) ([]byte, error) {
	// 1. 토큰 발급
	tokenBody, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	var tokenResp model.Response
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, err
	}

	tokenData, err := tokenResp.GetTokenData()
	if err != nil {
		return nil, err
	}

	// 2. 결제 상태 확인
	paymentBody, err := s.client.GetPaymentStatus(impUid, tokenData.AccessToken)
	if err != nil {
		return nil, err
	}

	var paymentResp model.Response
	if err := json.Unmarshal(paymentBody, &paymentResp); err != nil {
		return nil, err
	}

	paymentData, err := paymentResp.GetPaymentData()
	if err != nil {
		return nil, err
	}

	// 결제 승인 상태가 아니라면 단건조회 응답 반환
	if paymentData.Status != model.PaymentStatusPaid {
		return paymentBody, nil
	}

	// 3. 결제 취소 요청
	return s.client.RequestCancelPayment(impUid, tokenData.AccessToken)
}
