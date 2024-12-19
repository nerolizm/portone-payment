package model

import (
	"encoding/json"
	"fmt"
)

func (r *Response) GetTokenData() (TokenData, error) {
	if r.Response == nil {
		return TokenData{}, fmt.Errorf("response is nil")
	}

	var tokenData TokenData
	err := json.Unmarshal(r.Response, &tokenData)
	return tokenData, err
}

func (r *Response) GetPaymentData() (PaymentData, error) {
	if r.Response == nil {
		return PaymentData{}, fmt.Errorf("response is nil")
	}

	var paymentData PaymentData
	err := json.Unmarshal(r.Response, &paymentData)
	return paymentData, err
}
