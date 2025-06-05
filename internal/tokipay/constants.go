package tokipay

// API Endpoints
const (
	// Base URLs
	ProductionBaseURL = "https://api.toki.mn"
	TestBaseURL      = "https://qams-api.toki.mn"

	// API Key
	ThirdPartyAPIKey = "third-party-api-key"

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