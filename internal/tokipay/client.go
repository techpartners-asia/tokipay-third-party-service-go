package tokipay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a TokiPay API client
type Client struct {
	BaseURL      string
	Username     string
	Password     string
	MerchantID   string
	APIKey       string
	AccessToken  string
	TokenExpiry  time.Time
	HTTPClient   *http.Client
}

// New creates a new TokiPay client
func New(baseURL, username, password, merchantID string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		MerchantID: merchantID,
		APIKey:     ThirdPartyAPIKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// GetAccessToken retrieves an access token from the API
func (c *Client) GetAccessToken() error {
	// Check if we have a valid token
	if c.AccessToken != "" && time.Now().Before(c.TokenExpiry) {
		return nil
	}

	url := fmt.Sprintf("%s/api/v1/third-party/access-token", c.BaseURL)
	reqBody := TokenRequest{
		Username:   c.Username,
		Password:   c.Password,
		MerchantID: c.MerchantID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("error decoding error response: %v", err)
		}
		return fmt.Errorf("API error: %s", errResp.Message)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

// CreateQRPayment creates a new QR payment
func (c *Client) CreateQRPayment(req QRPaymentRequest) (*QRPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/third-party/qr-payment", c.BaseURL)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.Message)
	}

	var qrResp QRPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&qrResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &qrResp, nil
}

// CheckPaymentStatus checks the status of a payment
func (c *Client) CheckPaymentStatus(requestID string) (*PaymentStatusResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/third-party/check-payment/%s", c.BaseURL, requestID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.Message)
	}

	var statusResp PaymentStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &statusResp, nil
}

// CancelInvoice cancels an invoice
func (c *Client) CancelInvoice(requestID string, req CancelInvoiceRequest) (*CancelInvoiceResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/third-party/cancel-invoice/%s", c.BaseURL, requestID)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.Message)
	}

	var cancelResp CancelInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&cancelResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &cancelResp, nil
}

// RefundPayment refunds a payment
func (c *Client) RefundPayment(requestID string, req RefundPaymentRequest) (*RefundPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/third-party/refund-payment/%s", c.BaseURL, requestID)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.Message)
	}

	var refundResp RefundPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&refundResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &refundResp, nil
}

// RegisterVAT registers VAT details for an organization
func (c *Client) RegisterVAT(req VATRegistrationRequest) (*VATRegistrationResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/third-party/register-vat", c.BaseURL)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.Message)
	}

	var vatResp VATRegistrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&vatResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &vatResp, nil
} 