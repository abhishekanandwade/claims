package domain

import (
	"time"
)

// FreeLookCancellation represents free look cancellation and refund processing
// Table: freelook_cancellations
// Reference: seed/db/claims_database_schema.sql:621-654
type FreeLookCancellation struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	CancellationNumber string `db:"cancellation_number" json:"cancellation_number"`
	PolicyID           string `db:"policy_id" json:"policy_id"`
	BondTrackingID     *string `db:"bond_tracking_id" json:"bond_tracking_id,omitempty"`

	// Cancellation Request
	CancellationRequestDate time.Time `db:"cancellation_request_date" json:"cancellation_request_date"`
	CancellationReason      string    `db:"cancellation_reason" json:"cancellation_reason"`
	FreeLookPeriodValid     bool       `db:"freelook_period_valid" json:"freelook_period_valid"`
	RejectionReason         *string    `db:"rejection_reason" json:"rejection_reason,omitempty"`

	// Refund Calculation - BR-CLM-BOND-003
	TotalPremium        float64 `db:"total_premium" json:"total_premium"`
	ProRataRiskPremium  float64 `db:"pro_rata_risk_premium" json:"pro_rata_risk_premium"`
	StampDuty           float64 `db:"stamp_duty" json:"stamp_duty"`
	MedicalCosts        *float64 `db:"medical_costs" json:"medical_costs,omitempty"`
	OtherDeductions     *float64 `db:"other_deductions" json:"other_deductions,omitempty"`
	RefundAmount        float64 `db:"refund_amount" json:"refund_amount"` // Premium - (risk premium + stamp duty + medical + other)

	// Maker-Checker Workflow - BR-CLM-BOND-004
	MakerID               string    `db:"maker_id" json:"maker_id"`
	MakerEntryDate        time.Time `db:"maker_entry_date" json:"maker_entry_date"`
	CheckerID             *string   `db:"checker_id" json:"checker_id,omitempty"`
	CheckerVerificationDate *time.Time `db:"checker_verification_date" json:"checker_verification_date,omitempty"`
	MakerCheckerApproved  bool      `db:"maker_checker_approved" json:"maker_checker_approved"`

	// Refund Processing
	RefundTransactionID *string    `db:"refund_transaction_id" json:"refund_transaction_id,omitempty"`
	RefundStatus        string     `db:"refund_status" json:"refund_status"` // PENDING, PROCESSING, SUCCESS, FAILED
	RefundDate           *time.Time `db:"refund_date" json:"refund_date,omitempty"`
	LinkedToFinance      bool       `db:"linked_to_finance" json:"linked_to_finance"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for FreeLookCancellation
func (FreeLookCancellation) TableName() string {
	return "freelook_cancellations"
}
