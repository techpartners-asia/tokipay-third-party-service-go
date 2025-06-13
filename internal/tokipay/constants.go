package tokipay

// API Endpoints
const (
	// Base URLs
	ProductionBaseURL = "https://ms-api.toki.mn"
	TestBaseURL       = "https://qams-api.toki.mn"

	// API Key
	ThirdPartyAPIKey = "third_party_pay"

	// Payment Statuses
	PaymentStatusPending   = "PENDING"
	PaymentStatusApproved  = "APPROVED"
	PaymentStatusExpired   = "EXPIRED"
	PaymentStatusCancelled = "CANCELLED"
	PaymentStatusRefunded  = "REFUNDED"

	// VAT Types
	VATTypeIndividual = "INDIVIDUAL"
	VATTypeCompany    = "COMPANY"
)
