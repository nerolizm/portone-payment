package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nerolizm/portone-payment/internal/config"
	_http "github.com/nerolizm/portone-payment/internal/infrastructure/http"
)

const (
	tokenPath   = "/users/getToken"
	paymentPath = "/payments/%s"
	cancelPath  = "/payments/cancel"
)

type PaymentClientInterface interface {
	GetAccessToken() ([]byte, error)
	GetPaymentStatus(impUid, accessToken string) ([]byte, error)
	RequestCancelPayment(impUid, token string) ([]byte, error)
}

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient() PaymentClientInterface {
	return &Client{
		client:  &http.Client{},
		baseURL: "https://api.iamport.kr",
	}
}

func (c *Client) GetAccessToken() ([]byte, error) {
	requestBody := map[string]string{
		"imp_key":    config.Env.ImpKey,
		"imp_secret": config.Env.ImpSecret,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// resp 로깅
	fmt.Println("requestBody", string(jsonBody))

	resp, err := c.client.Post(c.baseURL+tokenPath, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, _http.NewHTTPError(
			resp.StatusCode,
			"failed to get access token",
		)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) GetPaymentStatus(impUid, accessToken string) ([]byte, error) {
	paymentURL := fmt.Sprintf(c.baseURL+paymentPath, impUid)

	req, err := http.NewRequest("GET", paymentURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, _http.NewHTTPError(
			resp.StatusCode,
			fmt.Sprintf("failed to get payment status for imp_uid %s", impUid),
		)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) RequestCancelPayment(impUid, token string) ([]byte, error) {
	cancelBody := map[string]string{
		"imp_uid": impUid,
	}

	jsonBody, err := json.Marshal(cancelBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+cancelPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, _http.NewHTTPError(
			resp.StatusCode,
			fmt.Sprintf("failed to cancel payment for imp_uid %s", impUid),
		)
	}

	return io.ReadAll(resp.Body)
}
