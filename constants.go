package tokipay

const (
	// API Endpoints
	TokenEndpoint         = "/third-party-service/v1/auth/token"
	QRPaymentEndpoint     = "/third-party-service/v1/payment-request/merchant-qr"
	MobilePaymentEndpoint = "/third-party-service/v1/payment-request/phone-number"
	DeeplinkEndpoint      = "/third-party-service/v1/payment-request/deeplink"
	StatusEndpoint        = "/third-party-service/v1/payment-request/status"
	CancelEndpoint        = "/third-party-service/v1/payment-request"
	RefundEndpoint        = "/third-party-service/v1/payment-request/refund"
	VATEndpoint           = "/third-party-service/v1/payment-request/vat"

	// API Keys
	ThirdPartyAPIKey = "third_party_pay"

	// Payment Types
	TypeSPOS          = "SPOS"
	TypeThirdPartyPay = "THIRD_PARTY_PAY"

	// Payment Status
	StatusPending   = "PENDING"
	StatusApproved  = "APPROVED"
	StatusExpired   = "EXPIRED"
	StatusCancelled = "CANCELLED"
	StatusSuccess   = "SUCCESS"
	StatusFailure   = "FAILURE"

	// VAT Types
	VATTypeOrganization = "ORGANIZATION"

	// Default Country Code
	DefaultCountryCode = "+976"

	// Token expiry (2 weeks)
	TokenExpiryDuration = 2 * 7 * 24 * 60 * 60 // seconds
)
