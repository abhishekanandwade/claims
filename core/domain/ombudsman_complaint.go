package domain

import (
	"time"
)

// OmbudsmanComplaint represents Insurance Ombudsman complaint lifecycle management
// Table: ombudsman_complaints
// Reference: seed/db/claims_database_schema.sql:532-577
type OmbudsmanComplaint struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Identifiers
	ComplaintNumber string `db:"complaint_number" json:"complaint_number"`
	ClaimID         *string `db:"claim_id" json:"claim_id,omitempty"`
	PolicyID        string `db:"policy_id" json:"policy_id"`

	// Complainant Information
	ComplainantName     string                 `db:"complainant_name" json:"complainant_name"`
	ComplainantContact  map[string]interface{} `db:"complainant_contact" json:"complainant_contact,omitempty"`
	ComplaintDescription string                `db:"complaint_description" json:"complaint_description"`
	ComplaintCategory   *string                `db:"complaint_category" json:"complaint_category,omitempty"`
	ClaimValue          *float64               `db:"claim_value" json:"claim_value,omitempty"` // BR-CLM-OMB-001: Must be <= ₹50 lakh

	// Admissibility Checks - BR-CLM-OMB-001
	RepresentationToInsurerDate  *time.Time `db:"representation_to_insurer_date" json:"representation_to_insurer_date,omitempty"`
	WaitPeriodCompleted          bool       `db:"wait_period_completed" json:"wait_period_completed"`
	LimitationPeriodValid        bool       `db:"limitation_period_valid" json:"limitation_period_valid"`
	ParallelLitigation           bool       `db:"parallel_litigation" json:"parallel_litigation"`
	Admissible                   *bool      `db:"admissible" json:"admissible,omitempty"`
	InadmissibilityReason        *string    `db:"inadmissibility_reason" json:"inadmissibility_reason,omitempty"`

	// Ombudsman Assignment
	OmbudsmanCenter       *string `db:"ombudsman_center" json:"ombudsman_center,omitempty"`
	JurisdictionBasis     *string `db:"jurisdiction_basis" json:"jurisdiction_basis,omitempty"`
	AssignedOmbudsmanID   *string `db:"assigned_ombudsman_id" json:"assigned_ombudsman_id,omitempty"`
	ConflictOfInterest    bool    `db:"conflict_of_interest" json:"conflict_of_interest"`

	// Status and Processing
	Status string `db:"status" json:"status"` // SUBMITTED, UNDER_REVIEW, MEDIATION, RECOMMENDATION, AWARD, CLOSED

	// Mediation Information
	MediationAttempted bool    `db:"mediation_attempted" json:"mediation_attempted"`
	MediationSuccessful *bool   `db:"mediation_successful" json:"mediation_successful,omitempty"`
	MediationTerms     *string `db:"mediation_terms" json:"mediation_terms,omitempty"`

	// Recommendation Information
	RecommendationIssued bool       `db:"recommendation_issued" json:"recommendation_issued"`
	RecommendationDate   *time.Time `db:"recommendation_date" json:"recommendation_date,omitempty"`

	// Award Information - BR-CLM-OMB-005
	AwardIssued        bool      `db:"award_issued" json:"award_issued"`
	AwardNumber        *string   `db:"award_number" json:"award_number,omitempty"`
	AwardAmount        *float64  `db:"award_amount" json:"award_amount,omitempty"` // BR-CLM-OMB-005: Max ₹50 lakh
	AwardDate          *time.Time `db:"award_date" json:"award_date,omitempty"`
	AwardDigitallySigned bool     `db:"award_digitally_signed" json:"award_digitally_signed"`

	// Compliance Information - BR-CLM-OMB-006
	ComplianceDueDate *time.Time `db:"compliance_due_date" json:"compliance_due_date,omitempty"` // 30 days from award date
	ComplianceStatus  *string    `db:"compliance_status" json:"compliance_status,omitempty"`
	ComplianceDate    *time.Time `db:"compliance_date" json:"compliance_date,omitempty"`
	EscalatedToIRDAI  bool       `db:"escalated_to_irdai" json:"escalated_to_irdai"`

	// Closure Information
	ClosureDate         *time.Time `db:"closure_date" json:"closure_date,omitempty"`
	ArchivalDate        *time.Time `db:"archival_date" json:"archival_date,omitempty"`
	RetentionPeriodYears *int      `db:"retention_period_years" json:"retention_period_years,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for OmbudsmanComplaint
func (OmbudsmanComplaint) TableName() string {
	return "ombudsman_complaints"
}
