package domain

import (
	"time"
)

// DocumentChecklistTemplate represents dynamic document checklist based on claim context
// Table: document_checklist_templates
// Reference: seed/db/claims_database_schema.sql:482-499
type DocumentChecklistTemplate struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Context Information
	ClaimType        string  `db:"claim_type" json:"claim_type"` // DEATH, MATURITY, SURVIVAL_BENEFIT, FREELOOK
	DeathType        *string `db:"death_type" json:"death_type,omitempty"` // NATURAL, UNNATURAL, ACCIDENTAL, SUICIDE, HOMICIDE
	NominationStatus *string `db:"nomination_status" json:"nomination_status,omitempty"` // PRESENT, ABSENT
	PolicyType       *string `db:"policy_type" json:"policy_type,omitempty"`

	// Document Information - BR-CLM-DC-015
	DocumentType        string  `db:"document_type" json:"document_type"`
	DocumentDescription *string `db:"document_description" json:"document_description,omitempty"`
	IsMandatory         bool    `db:"is_mandatory" json:"is_mandatory"` // Base mandatory or conditional
	DisplayOrder        *int    `db:"display_order" json:"display_order,omitempty"`

	// Validation Rules
	ValidationRules map[string]interface{} `db:"validation_rules" json:"validation_rules,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for DocumentChecklistTemplate
func (DocumentChecklistTemplate) TableName() string {
	return "document_checklist_templates"
}
