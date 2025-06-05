package tokipay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TokiPayClient represents the TokiPay third-party service client
type TokiPayClient struct {
	BaseURL      string
	Username     string
	Password     string
	MerchantID   string
	APIKey       string
	AccessToken  string
	TokenExpiry  time.Time
	HTTPClient   *http.Client
}

// TokiPay interface defines all available methods
type TokiPay interface {
	// Authentication
	GetAccessToken() error
	
	// Payment Methods
	CreateQRPayment(req QRPaymentRequest) (*QRPaymentResponse, error)
	CreateMobilePayment(req MobilePaymentRequest) (*MobilePaymentResponse, error)
	CreateDeeplinkPayment(req DeeplinkPaymentRequest) (*DeeplinkPaymentResponse, error)
	
	// Payment Management
	CheckPaymentStatus(requestID string) (*PaymentStatusResponse, error)
	CancelPayment(requestID string) error
	RefundPayment(req RefundRequest) (*RefundResponse, error)
	
	// VAT Management
	RegisterVAT(req VATRegistrationRequest) (*VATRegistrationResponse, error)
}

// New creates a new TokiPay client instance
func New(baseURL, username, password, merchantID string) TokiPay {
	return &TokiPayClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		MerchantID: merchantID,
		APIKey:     ThirdPartyAPIKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAccessToken retrieves and stores the access token
func (c *TokiPayClient) GetAccessToken() error {
	// Check if token is still valid
	if c.AccessToken != "" && time.Now().Before(c.TokenExpiry) {
		return nil
	}

	// Create basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.Password))
	
	req, err := http.NewRequest("GET", c.BaseURL+TokenEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokiPayResponse[TokenResponse]
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if tokenResp.Code != 200 {
		return fmt.Errorf("token request failed: %s", tokenResp.Error.Message)
	}

	c.AccessToken = tokenResp.Data.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(TokenExpiryDuration) * time.Second)

	return nil
}

// CreateQRPayment creates a QR payment request
func (c *TokiPayClient) CreateQRPayment(req QRPaymentRequest) (*QRPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID

	var resp TokiPayResponse[QRPaymentResponse]
	if err := c.makeRequest("POST", QRPaymentEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("QR payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CreateMobilePayment creates a mobile payment request
func (c *TokiPayClient) CreateMobilePayment(req MobilePaymentRequest) (*MobilePaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID
	if req.CountryCode == "" {
		req.CountryCode = DefaultCountryCode
	}
	if req.Type == "" {
		req.Type = TypeThirdPartyPay
	}

	var resp TokiPayResponse[MobilePaymentResponse]
	if err := c.makeRequest("POST", MobilePaymentEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("mobile payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CreateDeeplinkPayment creates a deeplink payment request
func (c *TokiPayClient) CreateDeeplinkPayment(req DeeplinkPaymentRequest) (*DeeplinkPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID
	req.Type = TypeThirdPartyPay

	var resp TokiPayResponse[DeeplinkPaymentResponse]
	if err := c.makeRequest("POST", DeeplinkEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("deeplink payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CheckPaymentStatus checks the status of a payment
func (c *TokiPayClient) CheckPaymentStatus(requestID string) (*PaymentStatusResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s?requestId=%s", StatusEndpoint, requestID)
	
	var resp TokiPayResponse[PaymentStatusResponse]
	if err := c.makeRequest("GET", endpoint, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("payment status check failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CancelPayment cancels a payment request
func (c *TokiPayClient) CancelPayment(requestID string) error {
	if err := c.GetAccessToken(); err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/%s", CancelEndpoint, requestID)
	
	var resp TokiPayResponse[interface{}]
	if err := c.makeRequest("PATCH", endpoint, nil, &resp); err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("payment cancellation failed: %s", resp.Error.Message)
	}

	return nil
}

// RefundPayment processes a refund
func (c *TokiPayClient) RefundPayment(req RefundRequest) (*RefundResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID

	var resp TokiPayResponse[RefundResponse]
	if err := c.makeRequest("POST", RefundEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("refund request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// RegisterVAT registers organization VAT details
func (c *TokiPayClient) RegisterVAT(req VATRegistrationRequest) (*VATRegistrationResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	var resp TokiPayResponse[VATRegistrationResponse]
	if err := c.makeRequest("POST", VATEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("VAT registration failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// makeRequest is a helper method to make HTTP requests
func (c *TokiPayClient) makeRequest(method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
} 