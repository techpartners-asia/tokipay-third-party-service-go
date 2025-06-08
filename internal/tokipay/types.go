package tokipay

// Response represents the standard TokiPay API response structure
type Response[T any] struct {
	Code      int       `json:"code"`
	Status    string    `json:"status"`
	Timestamp int64     `json:"timestamp"`
	Data      T         `json:"data"`
	Error     *APIError `json:"error"`
}

type APIError struct {
	Message string `json:"message"`
}

// Token Request/Response
type TokenRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	MerchantID string `json:"merchant_id"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

// QR Payment Request/Response
type QRPaymentRequest struct {
	SuccessURL string  `json:"successUrl" binding:"required"`
	FailureURL string  `json:"failureUrl" binding:"required"`
	OrderID    string  `json:"orderId" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Notes      string  `json:"notes,omitempty"`
	MerchantID string  `json:"merchantId" binding:"required"`
}

type QRPaymentResponse struct {
	RequestID     string `json:"requestId"`
	TransactionID string `json:"transactionId"`
}

// Mobile Payment Request/Response
type MobilePaymentRequest struct {
	SuccessURL      string   `json:"successUrl" binding:"required"`
	FailureURL      string   `json:"failureUrl" binding:"required"`
	OrderID         string   `json:"orderId" binding:"required"`
	MerchantID      string   `json:"merchantId" binding:"required"`
	Amount          float64  `json:"amount" binding:"required"`
	Notes           string   `json:"notes,omitempty"`
	PhoneNo         string   `json:"phoneNo" binding:"required"`
	CountryCode     string   `json:"countryCode" binding:"required"`
	Type            string   `json:"type" binding:"required"` // SPOS || THIRD_PARTY_PAY
	SuccessText     string   `json:"successText,omitempty"`
	EbarimtText     string   `json:"ebarimtText,omitempty"`
	ProductsInfo    []string `json:"productsInfo,omitempty"`
	PaymentCategory string   `json:"paymentCategory,omitempty"`
}

type MobilePaymentResponse struct {
	RequestID string `json:"requestId"`
}

// Deeplink Payment Request/Response
type DeeplinkPaymentRequest struct {
	SuccessURL string  `json:"successUrl" binding:"required"`
	FailureURL string  `json:"failureUrl" binding:"required"`
	OrderID    string  `json:"orderId" binding:"required"`
	MerchantID string  `json:"merchantId" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Notes      string  `json:"notes,omitempty"`
	Type       string  `json:"type" binding:"required"` // THIRD_PARTY_PAY
}

type DeeplinkPaymentResponse struct {
	Deeplink      string `json:"deeplink"`
	TransactionID string `json:"transactionId"`
}

// Payment Status Response
type PaymentStatusResponse struct {
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	PaidAmount    float64 `json:"paid_amount"`
	PaidDate      string  `json:"paid_date"`
}

type VATDetails struct {
	VATType string `json:"vatType"`
	VATID   string `json:"vatId"`
}

// Refund Request/Response
type RefundRequest struct {
	MerchantID  string `json:"merchantId" binding:"required"`
	TransNumber string `json:"transNumber" binding:"required"`
	Amount      string `json:"amount,omitempty"` // Optional, full refund if not provided
}

type RefundResponse struct {
	TransNumber      string `json:"transNumber"`
	Response         string `json:"response"`
	TxnNumber        string `json:"txnNumber"`
	TopupTransnumber string `json:"topupTransnumber"`
}

// VAT Registration Request/Response
type VATRegistrationRequest struct {
	TransactionID string `json:"transactionId" binding:"required"`
	DDTD          string `json:"DDTD" binding:"required"`
	TotalAmount   string `json:"totalAmount,omitempty"`
	VATAmount     string `json:"vatAmount,omitempty"`
	CreatedDate   string `json:"createdDate,omitempty"`
	MerchantName  string `json:"merchantName,omitempty"`
	MerchantTIN   string `json:"merchantTin,omitempty"`
}

type VATRegistrationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Callback Request from TokiPay
type CallbackRequest struct {
	OrderID       string  `json:"orderId"`
	RequestID     string  `json:"requestId"`
	Status        string  `json:"status"` // SUCCESS or FAILURE
	Amount        float64 `json:"amount"`
	Authorization string  `json:"authorization"`
}

// Callback Headers for organization transactions
type CallbackHeaders struct {
	VATID   string `header:"VAT_ID"`
	VATType string `header:"VAT_TYPE"`
}

// Cancel Invoice Request/Response
type CancelInvoiceRequest struct {
	Reason string `json:"reason"`
}

type CancelInvoiceResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// Refund Payment Request/Response
type RefundPaymentRequest struct {
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}

type RefundPaymentResponse struct {
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	RefundAmount  float64 `json:"refund_amount"`
	RefundDate    string  `json:"refund_date"`
}

// Error Response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
