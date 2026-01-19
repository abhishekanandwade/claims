package domain

import (
	"time"
)

// ClaimPayment represents payment disbursement tracking
// Table: claim_payments (partitioned by created_at)
// Reference: seed/db/claims_database_schema.sql:378-407
type ClaimPayment struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	PaymentID string `db:"payment_id" json:"payment_id"`
	ClaimID   string `db:"claim_id" json:"claim_id"`

	// Payment Information
	PaymentAmount     float64  `db:"payment_amount" json:"payment_amount"`
	PaymentMode       string   `db:"payment_mode" json:"payment_mode"` // BR-CLM-DC-017: AUTO_NEFT, POSB_TRANSFER, CHEQUE
	PaymentReference  *string  `db:"payment_reference" json:"payment_reference,omitempty"`
	UTRNumber         *string  `db:"utr_number" json:"utr_number,omitempty"`
	TransactionID     *string  `db:"transaction_id" json:"transaction_id,omitempty"`

	// Beneficiary Information
	BeneficiaryAccountNumber string  `db:"beneficiary_account_number" json:"beneficiary_account_number"`
	BeneficiaryIFSCCode      string  `db:"beneficiary_ifsc_code" json:"beneficiary_ifsc_code"`
	BeneficiaryName          string  `db:"beneficiary_name" json:"beneficiary_name"`
	BeneficiaryBankName      *string `db:"beneficiary_bank_name" json:"beneficiary_bank_name,omitempty"`

	// Initiation Information
	InitiatedBy string    `db:"initiated_by" json:"initiated_by"`
	InitiatedAt time.Time `db:"initiated_at" json:"initiated_at"`

	// Payment Status
	PaymentDate    *time.Time `db:"payment_date" json:"payment_date,omitempty"`
	PaymentStatus  string     `db:"payment_status" json:"payment_status"` // PENDING, PROCESSING, SUCCESS, FAILED
	FailureReason  *string    `db:"failure_reason" json:"failure_reason,omitempty"`
	RetryCount     int        `db:"retry_count" json:"retry_count"`

	// Reconciliation Information
	ReconciliationStatus *string    `db:"reconciliation_status" json:"reconciliation_status,omitempty"` // PENDING, MATCHED, MISMATCH
	ReconciledAt         *time.Time `db:"reconciled_at" json:"reconciled_at,omitempty"`
	VoucherNumber        *string    `db:"voucher_number" json:"voucher_number,omitempty"`
	VoucherDate          *time.Time `db:"voucher_date" json:"voucher_date,omitempty"`

	// Additional Data
	Metadata map[string]interface{} `db:"metadata" json:"metadata,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for ClaimPayment
func (ClaimPayment) TableName() string {
	return "claim_payments"
}
