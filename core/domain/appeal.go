package domain

import (
	"time"
)

// Appeal represents appeal workflow for rejected claims
// Table: appeals
// Reference: seed/db/claims_database_schema.sql:304-328
type Appeal struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	AppealNumber string `db:"appeal_number" json:"appeal_number"`
	ClaimID      string `db:"claim_id" json:"claim_id"`

	// Appellant Information
	AppellantName    string                 `db:"appellant_name" json:"appellant_name"`
	AppellantContact map[string]interface{} `db:"appellant_contact" json:"appellant_contact,omitempty"`

	// Appeal Details
	GroundsOfAppeal       string   `db:"grounds_of_appeal" json:"grounds_of_appeal"`
	SupportingDocuments   []string `db:"supporting_documents" json:"supporting_documents,omitempty"` // Array of document IDs
	CondonationRequest    bool     `db:"condonation_request" json:"condonation_request"`
	CondonationReason     *string  `db:"condonation_reason" json:"condonation_reason,omitempty"`
	SubmissionDate        time.Time `db:"submission_date" json:"submission_date"` // BR-CLM-DC-005: Within 90 days
	AppealDeadline        time.Time `db:"appeal_deadline" json:"appeal_deadline"` // BR-CLM-DC-007: 45-day decision timeline

	// Appellate Authority
	AppellateAuthorityID *string `db:"appellate_authority_id" json:"appellate_authority_id,omitempty"`

	// Status and Decision
	Status             string     `db:"status" json:"status"` // SUBMITTED, UNDER_REVIEW, ALLOWED, DISMISSED
	Decision           *string    `db:"decision" json:"decision,omitempty"`
	ReasonedOrder      *string    `db:"reasoned_order" json:"reasoned_order,omitempty"`
	RevisedClaimAmount *float64   `db:"revised_claim_amount" json:"revised_claim_amount,omitempty"`
	DecisionDate       *time.Time `db:"decision_date" json:"decision_date,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for Appeal
func (Appeal) TableName() string {
	return "appeals"
}
