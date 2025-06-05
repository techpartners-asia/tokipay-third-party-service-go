package tokipay

import "time"

// Standard TokiPay API Response Structure
type TokiPayResponse[T any] struct {
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
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
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
	Status      string      `json:"status"` // PENDING, APPROVED, EXPIRED, CANCELLED
	TransNumber string      `json:"transNumber,omitempty"`
	Fee         float64     `json:"fee,omitempty"`
	VATDetails  *VATDetails `json:"vatDetails,omitempty"`
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

// Generic Response Types
type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
} 