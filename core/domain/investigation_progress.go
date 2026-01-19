package domain

import (
	"time"
)

// InvestigationProgress represents heartbeat updates for long-running investigations
// Table: investigation_progress
// Reference: seed/db/claims_database_schema.sql:287-299
type InvestigationProgress struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Reference
	InvestigationID string `db:"investigation_id" json:"investigation_id"`

	// Progress Update
	UpdateDate              time.Time  `db:"update_date" json:"update_date"`
	ProgressPercentage      int        `db:"progress_percentage" json:"progress_percentage"` // 0-100
	ChecklistItemsCompleted []string   `db:"checklist_items_completed" json:"checklist_items_completed,omitempty"`
	Remarks                 string     `db:"remarks" json:"remarks"`
	EstimatedCompletionDate *time.Time `db:"estimated_completion_date" json:"estimated_completion_date,omitempty"`

	// Updated By
	UpdatedBy string `db:"updated_by" json:"updated_by"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the database table name for InvestigationProgress
func (InvestigationProgress) TableName() string {
	return "investigation_progress"
}
