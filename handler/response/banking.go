package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== BANKING & PAYMENT RESPONSE DTOs ====================

// ExtendedBankValidationResponse extends the base BankValidationResponse with additional fields
// Used for banking endpoints that require more detailed validation info
// POST /banking/validate-account
// POST /banking/validate-account-cbs
// POST /banking/validate-account-pfms
// POST /banking/penny-drop
// Reference: BR-CLM-DC-010 (Payment Disbursement Workflow)
type ExtendedBankValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                     ExtendedBankValidationData `json:"data"`
}

// ExtendedBankValidationData contains detailed bank validation results
type ExtendedBankValidationData struct {
	Valid               bool    `json:"valid"`
	AccountNumber       string  `json:"account_number"`
	AccountHolderName   string  `json:"account_holder_name"`
	BankName            string  `json:"bank_name,omitempty"`
	IFSCCode            string  `json:"ifsc_code,omitempty"`
	ValidationMethod    string  `json:"validation_method"`
	NameMatchPercentage float64 `json:"name_match_percentage,omitempty"`
	ValidationDate      string  `json:"validation_date,omitempty"` // YYYY-MM-DD HH:MM:SS format
	AccountStatus       *string `json:"account_status,omitempty"`  // ACTIVE, INACTIVE, CLOSED
	AccountType         *string `json:"account_type,omitempty"`    // SAVINGS, CURRENT, NRE
	BranchName          *string `json:"branch_name,omitempty"`
	City                *string `json:"city,omitempty"`
	State               *string `json:"state,omitempty"`
	PINCode             *string `json:"pincode,omitempty"`
	MICRCode            *string `json:"micr_code,omitempty"`
	FailureReason       *string `json:"failure_reason,omitempty"`
}

// NEFTTransferInitiatedResponse represents the response for NEFT transfer initiation
// POST /banking/neft-transfer
// Reference: BR-CLM-DC-010 (Disbursement Workflow)
type NEFTTransferInitiatedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PaymentID                 string  `json:"payment_id,omitempty"`
	TransactionID             string  `json:"transaction_id,omitempty"`
	ReferenceID               string  `json:"reference_id,omitempty"`
	Amount                    float64 `json:"amount,omitempty"`
	BankReference             *string `json:"bank_reference,omitempty"`
	InitiatedAt               string  `json:"initiated_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	ExpectedSettlement        *string `json:"expected_settlement,omitempty"` // YYYY-MM-DD HH:MM:SS format
	Status                    string  `json:"status,omitempty"`
}

// PaymentReconciliationResponse represents the response for payment reconciliation
// POST /banking/payment-reconciliation
// Reference: BR-CLM-PAY-001 (Daily Reconciliation)
type PaymentReconciliationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ReconciliationDate        string                    `json:"reconciliation_date,omitempty"` // YYYY-MM-DD format
	TotalPayments             int64                     `json:"total_payments,omitempty"`
	SuccessfulPayments        int64                     `json:"successful_payments,omitempty"`
	FailedPayments            int64                     `json:"failed_payments,omitempty"`
	PendingPayments           int64                     `json:"pending_payments,omitempty"`
	TotalAmountReconciled     float64                   `json:"total_amount_reconciled,omitempty"`
	MismatchedTransactions    []PaymentMismatchData     `json:"mismatched_transactions,omitempty"`
	ReconciliationSummary     ReconciliationSummaryData `json:"reconciliation_summary,omitempty"`
	ReconciledAt              string                    `json:"reconciled_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
}

// PaymentMismatchData represents mismatched transaction details
type PaymentMismatchData struct {
	PaymentID        string  `json:"payment_id,omitempty"`
	ExpectedAmount   float64 `json:"expected_amount,omitempty"`
	ActualAmount     *float64 `json:"actual_amount,omitempty"`
	MismatchReason   string  `json:"mismatch_reason,omitempty"`
	TransactionID    *string `json:"transaction_id,omitempty"`
	BankReference    *string `json:"bank_reference,omitempty"`
}

// ReconciliationSummaryData represents reconciliation summary
type ReconciliationSummaryData struct {
	MatchedCount      int64   `json:"matched_count,omitempty"`
	MatchedAmount     float64 `json:"matched_amount,omitempty"`
	UnmatchedCount    int64   `json:"unmatched_count,omitempty"`
	UnmatchedAmount   float64 `json:"unmatched_amount,omitempty"`
	ReconciliationRate float64 `json:"reconciliation_rate,omitempty"` // Percentage
}

// PaymentStatusResponse represents the payment status response
// GET /banking/payment-status/{payment_id}
type PaymentStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PaymentID                 string  `json:"payment_id,omitempty"`
	TransactionID             *string `json:"transaction_id,omitempty"`
	PaymentStatus             string  `json:"payment_status,omitempty"` // INITIATED, PROCESSING, SUCCESS, FAILED, CANCELLED
	PaymentReference          *string `json:"payment_reference,omitempty"`
	PaymentDate               *string `json:"payment_date,omitempty"` // YYYY-MM-DD HH:MM:SS format
	Amount                    float64 `json:"amount,omitempty"`
	BankName                  *string `json:"bank_name,omitempty"`
	AccountNumber             *string `json:"account_number,omitempty"`
	IFSCCode                  *string `json:"ifsc_code,omitempty"`
	BeneficiaryName           *string `json:"beneficiary_name,omitempty"`
	PaymentMode               *string `json:"payment_mode,omitempty"` // NEFT, RTGS, IMPS, POSB, CHEQUE
	InitiatedAt               *string `json:"initiated_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	CompletedAt               *string `json:"completed_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	FailureReason             *string `json:"failure_reason,omitempty"`
	UtrNumber                 *string `json:"utr_number,omitempty"` // Unified Payment Reference
}

// WebhookResponse represents the webhook acknowledgment response
// POST /webhooks/banking/payment-confirmation
type WebhookResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	WebhookReceived           bool   `json:"webhook_received,omitempty"`
	PaymentID                 string `json:"payment_id,omitempty"`
	ProcessedAt               string `json:"processed_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	Message                   string `json:"message,omitempty"`
}

// PaymentVoucherResponse represents the payment voucher generation response
// POST /banking/generate-voucher
// Reference: BR-CLM-PAY-002 (Voucher Generation)
type PaymentVoucherResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	VoucherNumber             string  `json:"voucher_number,omitempty"`
	VoucherDate               string  `json:"voucher_date,omitempty"` // YYYY-MM-DD format
	VoucherType               string  `json:"voucher_type,omitempty"` // PAYMENT, RECEIPT, JOURNAL
	VoucherURL                *string `json:"voucher_url,omitempty"`
	PaymentID                 string  `json:"payment_id,omitempty"`
	ClaimID                   string  `json:"claim_id,omitempty"`
	Amount                    float64 `json:"amount,omitempty"`
	BeneficiaryName           string  `json:"beneficiary_name,omitempty"`
	AccountNumber             string  `json:"account_number,omitempty"`
	BankName                  string  `json:"bank_name,omitempty"`
	GeneratedAt               string  `json:"generated_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	VoucherData               VoucherDetails `json:"voucher_data,omitempty"`
}

// VoucherDetails represents detailed voucher information
type VoucherDetails struct {
	PaymentDate         string  `json:"payment_date,omitempty"` // YYYY-MM-DD format
	AuthorizationBy     *string `json:"authorization_by,omitempty"`
	AuthorizationDate   *string `json:"authorization_date,omitempty"` // YYYY-MM-DD HH:MM:SS format
	AccountingHead      *string `json:"accounting_head,omitempty"`
	BudgetHead          *string `json:"budget_head,omitempty"`
	FinancialYear       *string `json:"financial_year,omitempty"`
	StampDuty           *float64 `json:"stamp_duty,omitempty"`
	NetAmount           float64 `json:"net_amount,omitempty"`
	Remarks             *string `json:"remarks,omitempty"`
	SupportingDocsCount int     `json:"supporting_docs_count,omitempty"`
}

// NewBankValidationResponse creates a new bank validation response
func NewBankValidationResponse(valid bool, accountNumber, accountHolderName, bankName, validationMethod string, nameMatchPercentage float64) *ExtendedBankValidationResponse {
	return &ExtendedBankValidationResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		Data: ExtendedBankValidationData{
			Valid:               valid,
			AccountNumber:       accountNumber,
			AccountHolderName:   accountHolderName,
			BankName:            bankName,
			ValidationMethod:    validationMethod,
			NameMatchPercentage: nameMatchPercentage,
			ValidationDate:      time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

// NewPaymentStatusResponse creates a new payment status response
func NewPaymentStatusResponse(paymentID, status string, amount float64) *PaymentStatusResponse {
	return &PaymentStatusResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		PaymentID:            paymentID,
		PaymentStatus:        status,
		Amount:               amount,
	}
}

// NewPaymentVoucherResponse creates a new payment voucher response
func NewPaymentVoucherResponse(voucherNumber, paymentID, claimID string, amount float64, beneficiaryName, accountNumber, bankName string) *PaymentVoucherResponse {
	return &PaymentVoucherResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		VoucherNumber:         voucherNumber,
		VoucherDate:           time.Now().Format("2006-01-02"),
		VoucherType:           "PAYMENT",
		PaymentID:             paymentID,
		ClaimID:               claimID,
		Amount:                amount,
		BeneficiaryName:       beneficiaryName,
		AccountNumber:         accountNumber,
		BankName:              bankName,
		GeneratedAt:           time.Now().Format("2006-01-02 15:04:05"),
	}
}
