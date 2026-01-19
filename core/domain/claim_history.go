package domain

import (
	"time"
)

// ClaimHistory represents the complete audit trail for all claim changes
// Table: claim_history
// Reference: seed/db/claims_database_schema.sql:424-450
type ClaimHistory struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Reference
	ClaimID string `db:"claim_id" json:"claim_id"`

	// Action Information
	ActionType        string `db:"action_type" json:"action_type"`
	ActionDescription *string `db:"action_description" json:"action_description,omitempty"`

	// Status Change Tracking
	OldStatus *string `db:"old_status" json:"old_status,omitempty"`
	NewStatus *string `db:"new_status" json:"new_status,omitempty"`

	// Value Change Tracking
	OldValues map[string]interface{} `db:"old_values" json:"old_values,omitempty"`
	NewValues map[string]interface{} `db:"new_values" json:"new_values,omitempty"`

	// Override Information - BR-CLM-DC-016/025
	OverrideApplied    bool    `db:"override_applied" json:"override_applied"`
	OverrideReason     *string `db:"override_reason" json:"override_reason,omitempty"`
	OverrideField      *string `db:"override_field" json:"override_field,omitempty"`
	OverrideOldValue   *string `db:"override_old_value" json:"override_old_value,omitempty"`
	OverrideNewValue   *string `db:"override_new_value" json:"override_new_value,omitempty"`
	DigitalSignatureHash *string `db:"digital_signature_hash" json:"digital_signature_hash,omitempty"`

	// Performed By Information
	PerformedBy string    `db:"performed_by" json:"performed_by"`
	PerformedAt time.Time `db:"performed_at" json:"performed_at"`

	// Request Metadata
	IPAddress  *string `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent  *string `db:"user_agent" json:"user_agent,omitempty"`
}

// TableName returns the database table name for ClaimHistory
func (ClaimHistory) TableName() string {
	return "claim_history"
}
