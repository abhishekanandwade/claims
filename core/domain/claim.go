package domain

import (
	"time"
)

// Claim represents the master claims entity
// Table: claims (partitioned by created_at)
// Reference: E-CLM-DC-001, seed/db/claims_database_schema.sql:110-183
type Claim struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	ClaimNumber string `db:"claim_number" json:"claim_number"`
	ClaimType   string `db:"claim_type" json:"claim_type"` // DEATH, MATURITY, SURVIVAL_BENEFIT, FREELOOK
	PolicyID    string `db:"policy_id" json:"policy_id"`
	CustomerID  string `db:"customer_id" json:"customer_id"`

	// Dates
	ClaimDate  time.Time  `db:"claim_date" json:"claim_date"`
	DeathDate  *time.Time `db:"death_date" json:"death_date,omitempty"`
	DeathPlace *string    `db:"death_place" json:"death_place,omitempty"`
	DeathType  *string    `db:"death_type" json:"death_type,omitempty"` // NATURAL, UNNATURAL, ACCIDENTAL, SUICIDE, HOMICIDE

	// Claimant Information
	ClaimantName     string  `db:"claimant_name" json:"claimant_name"`
	ClaimantType     *string `db:"claimant_type" json:"claimant_type,omitempty"` // NOMINEE, LEGAL_HEIR, ASSIGNEE
	ClaimantRelation *string `db:"claimant_relation" json:"claimant_relation,omitempty"`
	ClaimantPhone    *string `db:"claimant_phone" json:"claimant_phone,omitempty"`
	ClaimantEmail    *string `db:"claimant_email" json:"claimant_email,omitempty"`

	// Status and Workflow
	Status               string     `db:"status" json:"status"` // REGISTERED, DOCUMENT_PENDING, DOCUMENT_VERIFIED, etc.
	WorkflowState        *string    `db:"workflow_state" json:"workflow_state,omitempty"`
	InvestigationRequired bool      `db:"investigation_required" json:"investigation_required"`
	InvestigationStatus  *string    `db:"investigation_status" json:"investigation_status,omitempty"` // CLEAR, SUSPECT, FRAUD
	InvestigatorID       *string    `db:"investigator_id" json:"investigator_id,omitempty"`
	InvestigationStartDate *time.Time `db:"investigation_start_date" json:"investigation_start_date,omitempty"`
	InvestigationCompletionDate *time.Time `db:"investigation_completion_date" json:"investigation_completion_date,omitempty"`

	// Financial Information
	ClaimAmount      *float64 `db:"claim_amount" json:"claim_amount,omitempty"`
	ApprovedAmount   *float64 `db:"approved_amount" json:"approved_amount,omitempty"`
	SumAssured       *float64 `db:"sum_assured" json:"sum_assured,omitempty"`
	ReversionaryBonus *float64 `db:"reversionary_bonus" json:"reversionary_bonus,omitempty"`
	TerminalBonus    *float64 `db:"terminal_bonus" json:"terminal_bonus,omitempty"`
	OutstandingLoan  *float64 `db:"outstanding_loan" json:"outstanding_loan,omitempty"`
	UnpaidPremiums   *float64 `db:"unpaid_premiums" json:"unpaid_premiums,omitempty"`
	PenalInterest    *float64 `db:"penal_interest" json:"penal_interest,omitempty"` // BR-CLM-DC-009

	// Approval Information
	ApproverID      *string    `db:"approver_id" json:"approver_id,omitempty"`
	ApprovalDate    *time.Time `db:"approval_date" json:"approval_date,omitempty"`
	ApprovalRemarks *string    `db:"approval_remarks" json:"approval_remarks,omitempty"`
	DigitalSignatureHash *string `db:"digital_signature_hash" json:"digital_signature_hash,omitempty"` // BR-CLM-DC-025

	// Payment Information
	DisbursementDate   *time.Time `db:"disbursement_date" json:"disbursement_date,omitempty"`
	PaymentMode        *string    `db:"payment_mode" json:"payment_mode,omitempty"` // AUTO_NEFT, POSB_TRANSFER, CHEQUE
	PaymentReference   *string    `db:"payment_reference" json:"payment_reference,omitempty"`
	TransactionID      *string    `db:"transaction_id" json:"transaction_id,omitempty"`
	UTRNumber          *string    `db:"utr_number" json:"utr_number,omitempty"`

	// Bank Information
	BankAccountNumber     *string `db:"bank_account_number" json:"bank_account_number,omitempty"`
	BankIFSCCode          *string `db:"bank_ifsc_code" json:"bank_ifsc_code,omitempty"`
	BankAccountHolderName *string `db:"bank_account_holder_name" json:"bank_account_holder_name,omitempty"`
	BankName              *string `db:"bank_name" json:"bank_name,omitempty"`
	BankVerified          bool    `db:"bank_verified" json:"bank_verified"`
	BankVerificationMethod *string `db:"bank_verification_method" json:"bank_verification_method,omitempty"`

	// Rejection Information
	RejectionReason *string `db:"rejection_reason" json:"rejection_reason,omitempty"`
	RejectionCode   *string `db:"rejection_code" json:"rejection_code,omitempty"`

	// Appeal Information
	AppealSubmitted bool   `db:"appeal_submitted" json:"appeal_submitted"`
	AppealID        *string `db:"appeal_id" json:"appeal_id,omitempty"` // BR-CLM-DC-005

	// SLA Information - BR-CLM-DC-003/004/021
	SLADueDate    time.Time  `db:"sla_due_date" json:"sla_due_date"`
	SLABreached   bool       `db:"sla_breached" json:"sla_breached"`
	SLABreachDays int        `db:"sla_breach_days" json:"sla_breach_days"`
	SLAStatus     string     `db:"sla_status" json:"sla_status"` // GREEN, YELLOW, ORANGE, RED

	// Closure Information
	ClosureDate   *time.Time `db:"closure_date" json:"closure_date,omitempty"`
	ClosureReason *string    `db:"closure_reason" json:"closure_reason,omitempty"`

	// Additional Data
	Metadata     map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	SearchVector string                 `db:"search_vector" json:"-"` // Exclude from JSON response

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Version   int        `db:"version" json:"version"`
}

// TableName returns the database table name for Claim
func (Claim) TableName() string {
	return "claims"
}
