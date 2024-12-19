package service

import (
	"encoding/json"
	"fmt"

	v1 "github.com/nerolizm/portone-payment/internal/infrastructure/http/v1"
	"github.com/nerolizm/portone-payment/internal/model"
)

type PaymentService struct {
	client *v1.Client
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		client: v1.NewClient(),
	}
}

func (s *PaymentService) CancelPayment(impUid string) ([]byte, error) {
	// 1. 토큰 발급
	tokenBody, err := s.client.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token for imp_uid %s: %w", impUid, err)
	}

	var tokenResp model.Response
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response for imp_uid %s: %w", impUid, err)
	}

	tokenData, err := tokenResp.GetTokenData()
	if err != nil {
		return nil, fmt.Errorf("failed to get token data for imp_uid %s: %w", impUid, err)
	}

	// 2. 결제 상태 확인
	paymentBody, err := s.client.GetPaymentStatus(impUid, tokenData.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment status for imp_uid %s: %w", impUid, err)
	}

	var paymentResp model.Response
	if err := json.Unmarshal(paymentBody, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response for imp_uid %s: %w", impUid, err)
	}

	paymentData, err := paymentResp.GetPaymentData()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment data for imp_uid %s: %w", impUid, err)
	}

	// 결제 승인 상태가 아니라면 단건조회 응답 반환
	if paymentData.Status != model.PaymentStatusPaid {
		return paymentBody, nil
	}

	// 3. 결제 취소 요청
	cancelBody, err := s.client.RequestCancelPayment(impUid, tokenData.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment for imp_uid %s: %w", impUid, err)
	}

	var cancelResp model.Response
	if err := json.Unmarshal(cancelBody, &cancelResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel response for imp_uid %s: %w", impUid, err)
	}

	return cancelBody, nil
}
