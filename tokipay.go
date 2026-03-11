package tokipay

import (
	"fmt"
	"time"

	"resty.dev/v3"
)

// TokiPayClient represents the TokiPay third-party service client
type TokiPayClient struct {
	BaseURL     string
	Username    string
	Password    string
	MerchantID  string
	APIKey      string
	AccessToken string
	TokenExpiry time.Time
	resty       *resty.Client
}

// TokiPay interface defines all available methods
type TokiPay interface {
	// Authentication
	GetAccessToken() error
	Close() error

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

// New creates a new TokiPay client instance. Call Close() when done to release resources.
func New(baseURL, username, password, merchantID string) TokiPay {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(30 * time.Second)
	return &TokiPayClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		MerchantID: merchantID,
		APIKey:     ThirdPartyAPIKey,
		resty:      client,
	}
}

// GetAccessToken retrieves and stores the access token
func (c *TokiPayClient) GetAccessToken() error {
	if c.AccessToken != "" && time.Now().Before(c.TokenExpiry) {
		return nil
	}

	var tokenResp TokiPayResponse[TokenResponse]
	resp, err := c.resty.R().
		SetBasicAuth(c.Username, c.Password).
		SetResult(&tokenResp).
		Get(TokenEndpoint)
	if err != nil {
		body := ""
		if resp != nil {
			body = truncate(resp.String(), 500)
		}
		return fmt.Errorf("failed to execute request: %w\nresponse: %s", err, body)
	}
	if resp.IsError() {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	if tokenResp.Code != 200 {
		return fmt.Errorf("token request failed: %s", errorMessage(tokenResp.Error))
	}

	c.AccessToken = tokenResp.Data.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(TokenExpiryDuration) * time.Second)
	return nil
}

// Close releases resources held by the client. Call when done using the client.
func (c *TokiPayClient) Close() error {
	return c.resty.Close()
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
		return nil, fmt.Errorf("QR payment request failed: %s", errorMessage(resp.Error))
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
		return nil, fmt.Errorf("mobile payment request failed: %s", errorMessage(resp.Error))
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
		return nil, fmt.Errorf("deeplink payment request failed: %s", errorMessage(resp.Error))
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
		return nil, fmt.Errorf("payment status check failed: %s", errorMessage(resp.Error))
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
		return fmt.Errorf("payment cancellation failed: %s", errorMessage(resp.Error))
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
		return nil, fmt.Errorf("refund request failed: %s", errorMessage(resp.Error))
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
		return nil, fmt.Errorf("VAT registration failed: %s", errorMessage(resp.Error))
	}

	return &resp.Data, nil
}

// makeRequest is a helper method to make HTTP requests
func (c *TokiPayClient) makeRequest(method, endpoint string, body interface{}, result interface{}) error {
	req := c.resty.R().
		SetHeader("api-key", c.APIKey).
		SetAuthToken(c.AccessToken).
		SetResult(result)
	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Execute(method, endpoint)
	if err != nil {
		body := ""
		if resp != nil {
			body = truncate(resp.String(), 500)
		}
		return fmt.Errorf("failed to execute request: %w\nresponse: %s", err, body)
	}
	if resp.IsError() {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode(), truncate(resp.String(), 500))
	}
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// errorMessage safely returns the API error message, avoiding nil pointer panic
func errorMessage(err *APIError) string {
	if err != nil {
		return err.Message
	}
	return "unknown error"
}
