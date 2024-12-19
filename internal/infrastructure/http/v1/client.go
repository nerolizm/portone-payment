package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nerolizm/portone-payment/internal/config"
)

const (
	baseURL     = "https://api.iamport.kr"
	tokenPath   = "/users/getToken"
	paymentPath = "/payments/%s" // GET /payments/{imp_uid}
	cancelPath  = "/payments/cancel"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
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

	resp, err := c.client.Post(baseURL+tokenPath, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) GetPaymentStatus(impUid, accessToken string) ([]byte, error) {
	paymentURL := fmt.Sprintf(baseURL+paymentPath, impUid)

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

	req, err := http.NewRequest("POST", baseURL+cancelPath, bytes.NewBuffer(jsonBody))
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

	return io.ReadAll(resp.Body)
}
