# TokiPay Go Client

A Go client library for integrating with the TokiPay third-party payment service.

## Features

- QR Code Payment
- Mobile Payment
- Deeplink Payment
- Payment Status Check
- Payment Cancellation
- Refund Processing
- VAT Registration
- Organization VAT Handling

## Installation

```bash
go get github.com/yourusername/tokipay-go
```

## Configuration

The client requires the following configuration:

```go
client := tokipay.New(
    "https://qams-api.toki.mn", // Test environment
    "your_username",            // Provided by TokiPay team
    "your_password",            // Provided by TokiPay team
    "your_merchant_id",         // Provided by TokiPay team
)
```

### Environment Variables

For testing, you can use the following environment variables:

```bash
# Test Environment
TOKIPAY_TEST_BASE_URL=https://qams-api.toki.mn
TOKIPAY_TEST_USERNAME=your_test_username
TOKIPAY_TEST_PASSWORD=your_test_password
TOKIPAY_TEST_MERCHANT_ID=your_test_merchant_id

# Production Environment
TOKIPAY_PROD_BASE_URL=https://ms-api.toki.mn
TOKIPAY_PROD_USERNAME=your_prod_username
TOKIPAY_PROD_PASSWORD=your_prod_password
TOKIPAY_PROD_MERCHANT_ID=your_prod_merchant_id

# Test Configuration
TEST_ORDER_PREFIX=TEST_ORDER_
TEST_AMOUNT=1000
TEST_PHONE_NUMBER=99661234
TEST_COUNTRY_CODE=+976
```

## Usage Examples

### QR Payment

```go
qrReq := tokipay.QRPaymentRequest{
    SuccessURL: "https://yoursite.com/success",
    FailureURL: "https://yoursite.com/failure",
    OrderID:    "ORDER_12345",
    Amount:     1000.0,
    Notes:      "Test QR Payment",
}

qrResp, err := client.CreateQRPayment(qrReq)
if err != nil {
    log.Printf("QR Payment failed: %v", err)
} else {
    fmt.Printf("QR Code Data: %s\n", qrResp.RequestID)
}
```

### Mobile Payment

```go
mobileReq := tokipay.MobilePaymentRequest{
    SuccessURL:      "https://yoursite.com/success",
    FailureURL:      "https://yoursite.com/failure",
    OrderID:         "ORDER_12346",
    Amount:          2000.0,
    Notes:           "Test Mobile Payment",
    PhoneNo:         "99661234",
    CountryCode:     "+976",
    Type:            tokipay.TypeThirdPartyPay,
    SuccessText:     "Payment completed successfully!",
    ProductsInfo:    []string{"Coffee 1500 MNT", "Cookie 500 MNT"},
    PaymentCategory: "HOLD",
}

mobileResp, err := client.CreateMobilePayment(mobileReq)
if err != nil {
    log.Printf("Mobile Payment failed: %v", err)
} else {
    fmt.Printf("Request ID: %s\n", mobileResp.RequestID)
}
```

### Deeplink Payment

```go
deeplinkReq := tokipay.DeeplinkPaymentRequest{
    SuccessURL: "https://yoursite.com/success",
    FailureURL: "https://yoursite.com/failure",
    OrderID:    "ORDER_12347",
    Amount:     1500.0,
    Notes:      "Test Deeplink Payment",
}

deeplinkResp, err := client.CreateDeeplinkPayment(deeplinkReq)
if err != nil {
    log.Printf("Deeplink Payment failed: %v", err)
} else {
    fmt.Printf("Deeplink: %s\n", deeplinkResp.Deeplink)
}
```

### Check Payment Status

```go
statusResp, err := client.CheckPaymentStatus(requestID)
if err != nil {
    log.Printf("Status check failed: %v", err)
} else {
    fmt.Printf("Payment Status: %s\n", statusResp.Status)
    if statusResp.VATDetails != nil {
        fmt.Printf("VAT Type: %s, VAT ID: %s\n", statusResp.VATDetails.VATType, statusResp.VATDetails.VATID)
    }
}
```

### Cancel Payment

```go
err := client.CancelPayment(requestID)
if err != nil {
    log.Printf("Payment cancellation failed: %v", err)
} else {
    fmt.Printf("Payment cancelled successfully!\n")
}
```

### Refund Payment

```go
refundReq := tokipay.RefundRequest{
    TransNumber: "3425279",
    Amount:      "500", // Optional, full refund if not provided
}

refundResp, err := client.RefundPayment(refundReq)
if err != nil {
    log.Printf("Refund failed: %v", err)
} else {
    fmt.Printf("Refund Transaction Number: %s\n", refundResp.TransNumber)
}
```

### Register VAT Details

```go
vatReq := tokipay.VATRegistrationRequest{
    TransactionID: "3425279",
    DDTD:          "19910000004",
    TotalAmount:   "1000",
    VATAmount:     "90",
    CreatedDate:   "11/20/2024",
    MerchantName:  "Test Merchant",
    MerchantTIN:   "1234567",
}

vatResp, err := client.RegisterVAT(vatReq)
if err != nil {
    log.Printf("VAT registration failed: %v", err)
} else {
    fmt.Printf("VAT registered successfully!\n")
}
```

## Callback Handling

The client supports handling callbacks from TokiPay. When a payment is completed, TokiPay will send a callback to your success or failure URL with the following data:

```go
type CallbackRequest struct {
    OrderID       string  `json:"orderId"`
    RequestID     string  `json:"requestId"`
    Status        string  `json:"status"` // SUCCESS or FAILURE
    Amount        float64 `json:"amount"`
    Authorization string  `json:"authorization"`
}
```

For organization transactions, the callback will include additional headers:

```go
type CallbackHeaders struct {
    VATID   string `header:"VAT_ID"`
    VATType string `header:"VAT_TYPE"`
}
```

## Error Handling

All API calls return errors that should be handled appropriately. The client uses standard Go error handling patterns.

## Testing

To run the tests:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.