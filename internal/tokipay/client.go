package tokipay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents a TokiPay API client
type Client struct {
	BaseURL     string
	Username    string
	Password    string
	MerchantID  string
	APIKey      string
	AccessToken string
	TokenExpiry time.Time
	HTTPClient  *http.Client
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

	// Use Basic Authentication as per documentation
	auth := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.Password))
	url := fmt.Sprintf("%s/third-party-service/v1/auth/token", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

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

	var apiResponse Response[TokenResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	c.AccessToken = apiResponse.Data.AccessToken
	// Default to 2 weeks if ExpiresIn is not provided
	expiresIn := apiResponse.Data.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 2 * 7 * 24 * 60 * 60 // 2 weeks in seconds
	}
	c.TokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return nil
}

// CreateQRPayment creates a new QR payment
func (c *Client) CreateQRPayment(req QRPaymentRequest) (*QRPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	// Set merchant ID
	req.MerchantID = c.MerchantID

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/merchant-qr", c.BaseURL)
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
	httpReq.Header.Set("api-key", c.APIKey)

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

	var apiResponse Response[QRPaymentResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}

// CreateMobilePayment creates a mobile payment request
func (c *Client) CreateMobilePayment(req MobilePaymentRequest) (*MobilePaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	// Set merchant ID and defaults
	req.MerchantID = c.MerchantID
	if req.CountryCode == "" {
		req.CountryCode = "+976"
	}
	if req.Type == "" {
		req.Type = "THIRD_PARTY_PAY"
	}

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/phone-number", c.BaseURL)
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
	httpReq.Header.Set("api-key", c.APIKey)

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

	var apiResponse Response[MobilePaymentResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}

// CreateDeeplinkPayment creates a deeplink payment request
func (c *Client) CreateDeeplinkPayment(req DeeplinkPaymentRequest) (*DeeplinkPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	// Set merchant ID and type
	req.MerchantID = c.MerchantID
	req.Type = "THIRD_PARTY_PAY"

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/deeplink", c.BaseURL)
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
	httpReq.Header.Set("api-key", c.APIKey)

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

	var apiResponse Response[DeeplinkPaymentResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}

// CheckPaymentStatus checks the status of a payment
func (c *Client) CheckPaymentStatus(requestID string) (*PaymentStatusResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/status?requestId=%s", c.BaseURL, requestID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Accept", "*/*")
	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("api-key", c.APIKey)

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

	var apiResponse Response[PaymentStatusResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}

// CancelPayment cancels a payment request
func (c *Client) CancelPayment(requestID string) error {
	if err := c.GetAccessToken(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/%s", c.BaseURL, requestID)
	httpReq, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("api-key", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
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

	return nil
}

// RefundPayment refunds a payment
func (c *Client) RefundPayment(req RefundRequest) (*RefundResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	// Set merchant ID
	req.MerchantID = c.MerchantID

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/refund", c.BaseURL)
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

	var apiResponse Response[RefundResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}

// RegisterVAT registers VAT details for an organization
func (c *Client) RegisterVAT(req VATRegistrationRequest) (*VATRegistrationResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/third-party-service/v1/payment-request/vat", c.BaseURL)
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
	httpReq.Header.Set("api-key", c.APIKey)

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

	var apiResponse Response[VATRegistrationResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &apiResponse.Data, nil
}
