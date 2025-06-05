package tokipay

import "fmt"

// CreateQRPayment creates a QR payment request
func (c *Client) CreateQRPayment(req QRPaymentRequest) (*QRPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID

	var resp Response[QRPaymentResponse]
	if err := c.makeRequest("POST", QRPaymentEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("QR payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CreateMobilePayment creates a mobile payment request
func (c *Client) CreateMobilePayment(req MobilePaymentRequest) (*MobilePaymentResponse, error) {
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

	var resp Response[MobilePaymentResponse]
	if err := c.makeRequest("POST", MobilePaymentEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("mobile payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CreateDeeplinkPayment creates a deeplink payment request
func (c *Client) CreateDeeplinkPayment(req DeeplinkPaymentRequest) (*DeeplinkPaymentResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID
	req.Type = TypeThirdPartyPay

	var resp Response[DeeplinkPaymentResponse]
	if err := c.makeRequest("POST", DeeplinkEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("deeplink payment request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CheckPaymentStatus checks the status of a payment
func (c *Client) CheckPaymentStatus(requestID string) (*PaymentStatusResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s?requestId=%s", StatusEndpoint, requestID)
	
	var resp Response[PaymentStatusResponse]
	if err := c.makeRequest("GET", endpoint, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("payment status check failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// CancelPayment cancels a payment request
func (c *Client) CancelPayment(requestID string) error {
	if err := c.GetAccessToken(); err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/%s", CancelEndpoint, requestID)
	
	var resp Response[interface{}]
	if err := c.makeRequest("PATCH", endpoint, nil, &resp); err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("payment cancellation failed: %s", resp.Error.Message)
	}

	return nil
}

// RefundPayment processes a refund
func (c *Client) RefundPayment(req RefundRequest) (*RefundResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	req.MerchantID = c.MerchantID

	var resp Response[RefundResponse]
	if err := c.makeRequest("POST", RefundEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("refund request failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
}

// RegisterVAT registers organization VAT details
func (c *Client) RegisterVAT(req VATRegistrationRequest) (*VATRegistrationResponse, error) {
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	var resp Response[VATRegistrationResponse]
	if err := c.makeRequest("POST", VATEndpoint, req, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("VAT registration failed: %s", resp.Error.Message)
	}

	return &resp.Data, nil
} 