package model

import "encoding/json"

type PaymentStatus string

const (
	PaymentStatusReady  PaymentStatus = "ready"  // 브라우저 창 이탈, 가상계좌 발급 완료 등 미결제 상태
	PaymentStatusPaid   PaymentStatus = "paid"   // 결제완료
	PaymentStatusFailed PaymentStatus = "failed" // 신용카드 한도 초과, 체크카드 잔액 부족, 브라우저 창 종료 또는 취소 버튼 클릭 등 결제실패 상태
)

type Response struct {
	Code     int             `json:"code,omitempty"`
	Message  string          `json:"message,omitempty"`
	Response json.RawMessage `json:"response,omitempty"`
}

// 결제내역 단건조회 API 응답 데이터
// https://developers.portone.io/api/rest-v1/payment?v=v1#get%20%2Fpayments%2F%7Bimp_uid%7D
type PaymentData struct {
	Status PaymentStatus `json:"status"`
}

// 토큰 발급 API 응답 데이터
// https://developers.portone.io/api/rest-v1/auth?v=v1#post%20%2Fusers%2FgetToken
type TokenData struct {
	AccessToken string `json:"access_token"`
}

// 결제 취소 API 요청 데이터
// https://developers.portone.io/api/rest-v1/payment?v=v1#post%20%2Fpayments%2Fcancel
type CancelRequest struct {
	ImpUid string `json:"imp_uid"`
}
