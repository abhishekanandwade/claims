package domain

import (
	"time"
)

// Investigation represents investigation workflow tracking
// Table: investigations
// Reference: seed/db/claims_database_schema.sql:247-280
type Investigation struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	InvestigationID string `db:"investigation_id" json:"investigation_id"`
	ClaimID         string `db:"claim_id" json:"claim_id"`

	// Assignment Information
	AssignedBy      string    `db:"assigned_by" json:"assigned_by"`
	InvestigatorID  string    `db:"investigator_id" json:"investigator_id"`
	InvestigatorRank *string  `db:"investigator_rank" json:"investigator_rank,omitempty"`
	Jurisdiction    *string   `db:"jurisdiction" json:"jurisdiction,omitempty"`
	AutoAssigned    bool      `db:"auto_assigned" json:"auto_assigned"`
	AssignmentDate  time.Time `db:"assignment_date" json:"assignment_date"`

	// SLA Information - BR-CLM-DC-002: 21-day investigation SLA
	DueDate time.Time `db:"due_date" json:"due_date"`

	// Status and Progress
	Status               string  `db:"status" json:"status"` // ASSIGNED, IN_PROGRESS, SUBMITTED, REVIEWED, COMPLETED, CANCELLED
	ProgressPercentage   int     `db:"progress_percentage" json:"progress_percentage"` // 0-100
	InvestigationOutcome *string `db:"investigation_outcome" json:"investigation_outcome,omitempty"` // CLEAR, SUSPECT, FRAUD

	// Investigation Details
	CauseOfDeath              *string `db:"cause_of_death" json:"cause_of_death,omitempty"`
	CauseOfDeathVerified      bool    `db:"cause_of_death_verified" json:"cause_of_death_verified"`
	HospitalRecordsVerified   bool    `db:"hospital_records_verified" json:"hospital_records_verified"`
	DetailedFindings          *string `db:"detailed_findings" json:"detailed_findings,omitempty"`
	Recommendation            *string `db:"recommendation" json:"recommendation,omitempty"`
	ReportDocumentID          *string `db:"report_document_id" json:"report_document_id,omitempty"`
	SubmittedAt               *time.Time `db:"submitted_at" json:"submitted_at,omitempty"`

	// Review Information
	ReviewedBy       *string    `db:"reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
	ReviewDecision   *string    `db:"review_decision" json:"review_decision,omitempty"`
	ReviewerRemarks  *string    `db:"reviewer_remarks" json:"reviewer_remarks,omitempty"`
	ReinvestigationCount int     `db:"reinvestigation_count" json:"reinvestigation_count"` // BR-CLM-DC-013: Max 2

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for Investigation
func (Investigation) TableName() string {
	return "investigations"
}
