package domain

import (
	"time"
)

// AMLAlert represents AML/CFT alert detection and tracking
// Table: aml_alerts
// Reference: E-CLM-AML-001, seed/db/claims_database_schema.sql:335-371
type AMLAlert struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	AlertID         string `db:"alert_id" json:"alert_id"`
	TriggerCode     string `db:"trigger_code" json:"trigger_code"` // BR-CLM-AML-001 to 005: AML_001 to AML_005
	PolicyID        string `db:"policy_id" json:"policy_id"`
	CustomerID      *string `db:"customer_id" json:"customer_id,omitempty"`

	// Transaction Information
	TransactionType   string     `db:"transaction_type" json:"transaction_type"`
	TransactionAmount *float64   `db:"transaction_amount" json:"transaction_amount,omitempty"`
	TransactionDate   time.Time  `db:"transaction_date" json:"transaction_date"`
	PaymentMode       *string    `db:"payment_mode" json:"payment_mode,omitempty"`

	// Risk Assessment
	RiskLevel  string  `db:"risk_level" json:"risk_level"` // LOW, MEDIUM, HIGH, CRITICAL
	RiskScore  *int    `db:"risk_score" json:"risk_score,omitempty"` // 0-100
	AlertStatus string  `db:"alert_status" json:"alert_status"` // FLAGGED, UNDER_REVIEW, FILED, CLOSED
	AlertDescription *string `db:"alert_description" json:"alert_description,omitempty"`
	TriggerDetails map[string]interface{} `db:"trigger_details" json:"trigger_details,omitempty"`

	// Review Information
	ReviewedBy      *string    `db:"reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
	ReviewDecision  *string    `db:"review_decision" json:"review_decision,omitempty"`
	OfficerRemarks  *string    `db:"officer_remarks" json:"officer_remarks,omitempty"`
	ActionTaken     *string    `db:"action_taken" json:"action_taken,omitempty"`

	// Transaction Control
	TransactionBlocked bool `db:"transaction_blocked" json:"transaction_blocked"`

	// Filing Information - BR-CLM-AML-006/007
	FilingRequired bool    `db:"filing_required" json:"filing_required"`
	FilingType     *string `db:"filing_type" json:"filing_type,omitempty"` // STR, CTR, CCR, NTR
	FilingStatus   *string `db:"filing_status" json:"filing_status,omitempty"`
	FilingReference *string `db:"filing_reference" json:"filing_reference,omitempty"`
	FiledAt        *time.Time `db:"filed_at" json:"filed_at,omitempty"`
	FiledBy        *string   `db:"filed_by" json:"filed_by,omitempty"`

	// PAN Information
	PANNumber         *string `db:"pan_number" json:"pan_number,omitempty"`
	PANVerified       *bool   `db:"pan_verified" json:"pan_verified,omitempty"`
	NomineeChangeDetected bool `db:"nominee_change_detected" json:"nominee_change_detected"` // BR-CLM-AML-003

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for AMLAlert
func (AMLAlert) TableName() string {
	return "aml_alerts"
}
