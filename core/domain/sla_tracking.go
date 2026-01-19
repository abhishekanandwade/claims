package domain

import (
	"time"
)

// ClaimSLATracking represents real-time SLA monitoring with color-coded alerts
// Table: claim_sla_tracking
// Reference: seed/db/claims_database_schema.sql:505-526
type ClaimSLATracking struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Reference
	ClaimID string `db:"claim_id" json:"claim_id"`

	// SLA Information
	SLAType          string    `db:"sla_type" json:"sla_type"` // DOCUMENT_COLLECTION, INVESTIGATION, APPROVAL, PAYMENT
	SLAStartDate     time.Time `db:"sla_start_date" json:"sla_start_date"`
	SLADueDate       time.Time `db:"sla_due_date" json:"sla_due_date"`
	SLATotalDays     int       `db:"sla_total_days" json:"sla_total_days"`
	SLAElapsedDays   int       `db:"sla_elapsed_days" json:"sla_elapsed_days"`
	SLARemainingDays int       `db:"sla_remaining_days" json:"sla_remaining_days"`

	// SLA Status - BR-CLM-DC-021
	SLAStatus  string     `db:"sla_status" json:"sla_status"` // GREEN, YELLOW, ORANGE, RED
	SLABreach  bool       `db:"sla_breach" json:"sla_breach"`
	SLABreachDate *time.Time `db:"sla_breach_date" json:"sla_breach_date,omitempty"`
	SLACompletionDate *time.Time `db:"sla_completion_date" json:"sla_completion_date,omitempty"`

	// Escalation Information
	EscalationTriggered  bool       `db:"escalation_triggered" json:"escalation_triggered"`
	EscalationLevel      int        `db:"escalation_level" json:"escalation_level"`
	LastEscalationDate   *time.Time `db:"last_escalation_date" json:"last_escalation_date,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for ClaimSLATracking
func (ClaimSLATracking) TableName() string {
	return "claim_sla_tracking"
}
