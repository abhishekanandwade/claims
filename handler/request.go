package handler

import (
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ========================================
// DEATH CLAIMS - REQUEST DTOS
// ========================================

// RegisterDeathClaimRequest represents the request body for registering a death claim
// POST /claims/death/register
// Reference: FR-CLM-DC-001, BR-CLM-DC-001, WF-CLM-DC-001
type RegisterDeathClaimRequest struct {
	PolicyID         string  `json:"policy_id" validate:"required"`
	DeathDate        string  `json:"death_date" validate:"required"` // YYYY-MM-DD format
	DeathPlace       string  `json:"death_place" validate:"required"`
	DeathType        string  `json:"death_type" validate:"required,oneof=NATURAL UNNATURAL ACCIDENTAL SUICIDE HOMICIDE"`
	ClaimantName     string  `json:"claimant_name" validate:"required"`
	ClaimantType     string  `json:"claimant_type" validate:"required,oneof=NOMINEE LEGAL_HEIR ASSIGNEE"`
	ClaimantRelation *string `json:"claimant_relation,omitempty"`
	ClaimantPhone    *string `json:"claimant_phone,omitempty" validate:"omitempty,len=10"`
	ClaimantEmail    *string `json:"claimant_email,omitempty" validate:"omitempty,email"`
}

// ToDomain converts request to domain model
func (r RegisterDeathClaimRequest) ToDomain() domain.Claim {
	return domain.Claim{
		PolicyID:         r.PolicyID,
		ClaimType:        "DEATH",
		ClaimantName:     r.ClaimantName,
		ClaimantType:     &r.ClaimantType,
		ClaimantRelation: r.ClaimantRelation,
		ClaimantPhone:    r.ClaimantPhone,
		ClaimantEmail:    r.ClaimantEmail,
		DeathPlace:       &r.DeathPlace,
		DeathType:        &r.DeathType,
		Status:           "REGISTERED",
	}
}

// CalculateDeathClaimAmountRequest represents the request for pre-calculating death claim amount
// POST /claims/death/calculate-amount
// Reference: CALC-001
type CalculateDeathClaimAmountRequest struct {
	PolicyID  string `json:"policy_id" validate:"required"`
	DeathDate string `json:"death_date" validate:"required"` // YYYY-MM-DD format
}

// ClaimIDUri represents the URI parameter for claim_id
type ClaimIDUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetDocumentChecklistRequest represents the request for getting document checklist
// GET /claims/death/{claim_id}/document-checklist
// Reference: FR-CLM-DC-002, VR-CLM-DC-001 to VR-CLM-DC-007
type GetDocumentChecklistRequest struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetDynamicDocumentChecklistUri represents the query parameters for dynamic document checklist
// GET /claims/death/document-checklist-dynamic
// Reference: DFC-001
type GetDynamicDocumentChecklistUri struct {
	DeathType        string `form:"death_type" validate:"required,oneof=NATURAL UNNATURAL ACCIDENTAL SUICIDE HOMICIDE"`
	NominationStatus string `form:"nomination_status" validate:"required,oneof=PRESENT ABSENT"`
	PolicyType       string `form:"policy_type" validate:"required,oneof=STANDARD ENDOWMENT WHOLE_LIFE TERM"`
}

// UploadClaimDocumentsRequest represents the request for uploading claim documents
// POST /claims/death/{claim_id}/documents
type UploadClaimDocumentsRequest struct {
	ClaimID       string `uri:"claim_id" validate:"required"`
	DocumentType  string `json:"document_type" validate:"required"`
	FileName      string `json:"file_name" validate:"required"`
	FileContent   []byte `json:"file_content" validate:"required"`
	DocumentName  string `json:"document_name" validate:"required"`
	DocumentSubType *string `json:"document_sub_type,omitempty"`
}

// CheckDocumentCompletenessUri represents the URI parameter for checking document completeness
// GET /claims/death/{claim_id}/document-completeness
type CheckDocumentCompletenessUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// SendDocumentReminderRequest represents the request for sending document submission reminder
// POST /claims/death/{claim_id}/send-reminder
type SendDocumentReminderRequest struct {
	ClaimID      string   `uri:"claim_id" validate:"required"`
	ReminderType string   `json:"reminder_type" validate:"required,oneof=MISSING_DOCUMENTS SLA_WARNING FINAL_NOTICE"`
	Channels     []string `json:"channels" validate:"omitempty,dive,oneof=SMS EMAIL WHATSAPP"`
}

// CalculateBenefitRequest represents the request for calculating death claim benefit
// POST /claims/death/{claim_id}/calculate-benefit
// Reference: FR-CLM-DC-004, BR-CLM-DC-008
type CalculateBenefitRequest struct {
	ClaimID            string  `uri:"claim_id" validate:"required"`
	ManualOverride     bool    `json:"manual_override"`
	OverrideReason     *string `json:"override_reason,omitempty"`
	SupervisorApproval *string `json:"supervisor_approval,omitempty"`
}

// OverrideCalculationRequest represents the request for overriding automated calculation
// POST /claims/death/{claim_id}/calculation/override
type OverrideCalculationRequest struct {
	ClaimID        string  `uri:"claim_id" validate:"required"`
	OverrideAmount float64 `json:"override_amount" validate:"required,gt=0"`
	OverrideReason string  `json:"override_reason" validate:"required,max=500"`
	SupervisorID   string  `json:"supervisor_id" validate:"required"`
}

// ApproveCalculationUri represents the URI parameter for approving calculation
// POST /claims/death/{claim_id}/calculation/approve
type ApproveCalculationUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetEligibleApproversUri represents the URI parameter for getting eligible approvers
// GET /claims/death/{claim_id}/eligible-approvers
// Reference: BR-CLM-DC-015
type GetEligibleApproversUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetApprovalDetailsUri represents the URI parameter for getting approval details
// GET /claims/death/{claim_id}/approval-details
type GetApprovalDetailsUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetFraudRedFlagsUri represents the URI parameter for getting fraud detection red flags
// GET /claims/death/{claim_id}/red-flags
type GetFraudRedFlagsUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// ApproveClaimRequest represents the request for approving death claim
// POST /claims/death/{claim_id}/approve
// Reference: BR-CLM-DC-005, BR-CLM-DC-015, BR-CLM-DC-025
type ApproveClaimRequest struct {
	ClaimID             string  `uri:"claim_id" validate:"required"`
	ApprovalRemarks     string  `json:"approval_remarks" validate:"required,max=1000"`
	DigitalSignature    bool    `json:"digital_signature"`
	ConditionApproval   *string `json:"condition_approval,omitempty"` // If approved with conditions
	InvestigationWaiver *string `json:"investigation_waiver,omitempty"` // Supervisor waiver for investigation
}

// RejectClaimRequest represents the request for rejecting death claim
// POST /claims/death/{claim_id}/reject
type RejectClaimRequest struct {
	ClaimID                   string  `uri:"claim_id" validate:"required"`
	RejectionReason           string  `json:"rejection_reason" validate:"required,oneof=FRAUD_DETECTED MATERIAL_SUPPRESSION POLICY_INVALID INCOMPLETE_DOCUMENTS INVESTIGATION_ADVERSE OTHER"`
	DetailedJustification     string  `json:"detailed_justification" validate:"required,max=2000"`
	AppealRightsCommunicated  bool    `json:"appeal_rights_communicated"`
	InvestigationReportID     *string `json:"investigation_report_id,omitempty"` // Required if rejection due to investigation
}

// ValidateBankAccountRequest represents the request for validating bank account
// POST /claims/death/{claim_id}/validate-bank-account
type ValidateBankAccountRequest struct {
	ClaimID          string `uri:"claim_id" validate:"required"`
	ValidationMethod string `json:"validation_method" validate:"omitempty,oneof=CBS_API PFMS_API PENNY_DROP"`
}

// DisburseClaimRequest represents the request for initiating payment disbursement
// POST /claims/death/{claim_id}/disburse
// Reference: FR-CLM-DC-010, BR-CLM-DC-010
type DisburseClaimRequest struct {
	ClaimID        string  `uri:"claim_id" validate:"required"`
	PaymentMode    string `json:"payment_mode" validate:"required,oneof=AUTO_NEFT POSB_TRANSFER CHEQUE"`
	PaymentDetails *string `json:"payment_details,omitempty"` // Additional payment details if needed
}

// CloseClaimRequest represents the request for closing claim
// POST /claims/death/{claim_id}/close
type CloseClaimRequest struct {
	ClaimID           string  `uri:"claim_id" validate:"required"`
	ClosureReason     string  `json:"closure_reason" validate:"required,oneof=PAYMENT_COMPLETED CLAIM_WITHDRAWN APPEAL_REJECTED"`
	ArchivalRequired  bool    `json:"archival_required"`
}

// CancelClaimRequest represents the request for cancelling claim (claimant withdrawal)
// POST /claims/death/{claim_id}/cancel
type CancelClaimRequest struct {
	ClaimID             string  `uri:"claim_id" validate:"required"`
	CancellationReason  string  `json:"cancellation_reason" validate:"required,max=500"`
	RequestedBy         string  `json:"requested_by" validate:"required"`
	CanResubmit         bool    `json:"can_resubmit"`
}

// ReturnClaimRequest represents the request for returning claim to claimant
// POST /claims/death/{claim_id}/return
type ReturnClaimRequest struct {
	ClaimID            string  `uri:"claim_id" validate:"required"`
	ReturnReason       string  `json:"return_reason" validate:"required,oneof=DOCUMENTS_NOT_RECEIVED SLA_EXPIRED"`
	ResubmitInstructions *string `json:"resubmit_instructions,omitempty" validate:"omitempty,max=1000"`
}

// RequestFeedbackUri represents the URI parameter for requesting customer feedback
// POST /claims/death/{claim_id}/request-feedback
type RequestFeedbackUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// ListClaimsParams represents query parameters for listing death claims
// GET /claims/death/approval-queue (and other list endpoints)
type ListClaimsParams struct {
	port.MetadataRequest
	// Add additional filters
	Status         string `form:"status" validate:"omitempty"`
	ClaimType      string `form:"claim_type" validate:"omitempty"`
	PolicyID       string `form:"policy_id" validate:"omitempty"`
	CustomerID     string `form:"customer_id" validate:"omitempty"`
	ClaimantName   string `form:"claimant_name" validate:"omitempty"`
	StartDate      string `form:"start_date" validate:"omitempty"` // YYYY-MM-DD
	EndDate        string `form:"end_date" validate:"omitempty"`   // YYYY-MM-DD
	InvestigationStatus string `form:"investigation_status" validate:"omitempty,oneof=CLEAR SUSPECT FRAUD"`
	SLAStatus      string `form:"sla_status" validate:"omitempty,oneof=GREEN YELLOW ORANGE RED"`
}

// GetPendingInvestigationClaimsUri represents query parameters for pending investigation claims
// GET /claims/death/pending-investigation
type GetPendingInvestigationClaimsUri struct {
	port.MetadataRequest
	Jurisdiction string `form:"jurisdiction" validate:"omitempty"`
	SLAStatus    string `form:"sla_status" validate:"omitempty,oneof=GREEN YELLOW RED"`
}

// InvestigationIDUri represents the URI parameter for investigation_id
type InvestigationIDUri struct {
	ClaimID        string `uri:"claim_id" validate:"required"`
	InvestigationID string `uri:"investigation_id" validate:"required"`
}

// AssignInvestigationRequest represents the request for assigning investigation officer
// POST /claims/death/{claim_id}/investigation/assign-officer
// Reference: BR-CLM-DC-002
type AssignInvestigationRequest struct {
	ClaimID        string  `uri:"claim_id" validate:"required"`
	InvestigatorID string  `json:"investigator_id" validate:"required"`
	Priority       string `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	AssignmentType string `json:"assignment_type" validate:"omitempty,oneof=AUTO MANUAL"` // Default: AUTO
}

// InvestigationProgressRequest represents the request for submitting investigation progress
// POST /claims/death/{claim_id}/investigation/{investigation_id}/progress-update
type InvestigationProgressRequest struct {
	ClaimID        string  `uri:"claim_id" validate:"required"`
	InvestigationID string `uri:"investigation_id" validate:"required"`
	ProgressNotes  string `json:"progress_notes" validate:"required,max=2000"`
	NextSteps      *string `json:"next_steps,omitempty" validate:"omitempty,max=500"`
	Percentage     int    `json:"percentage" validate:"required,min=0,max=100"`
}

// SubmitInvestigationReportRequest represents the request for submitting final investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/submit-report
type SubmitInvestigationReportRequest struct {
	ClaimID         string  `uri:"claim_id" validate:"required"`
	InvestigationID string  `uri:"investigation_id" validate:"required"`
	ReportOutcome   string `json:"report_outcome" validate:"required,oneof=CLEAR SUSPECT FRAUD"`
	Findings        string  `json:"findings" validate:"required,max=5000"`
	Evidence        *string `json:"evidence,omitempty" validate:"omitempty,max=2000"`
	Recommendation  string  `json:"recommendation" validate:"required,max=1000"`
	ReportDate      string  `json:"report_date" validate:"required"` // YYYY-MM-DD
}

// ReviewInvestigationReportRequest represents the request for reviewing investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/review
type ReviewInvestigationReportRequest struct {
	ClaimID         string  `uri:"claim_id" validate:"required"`
	InvestigationID string  `uri:"investigation_id" validate:"required"`
	ReviewDecision  string `json:"review_decision" validate:"required,oneof=ACCEPT REINVESTIGATE ESCALATE"`
	ReviewerRemarks string `json:"reviewer_remarks" validate:"required,max=2000"`
	ReinvestigationReason *string `json:"reinvestigation_reason,omitempty" validate:"omitempty,max=1000"`
}

// TriggerReinvestigationRequest represents the request for triggering reinvestigation
// POST /claims/death/{id}/investigation/trigger-reinvestigation
// Reference: BR-CLM-DC-013 (max 2 times)
type TriggerReinvestigationRequest struct {
	ClaimID             string   `uri:"id" validate:"required"`
	ReinvestigationReason string  `json:"reinvestigation_reason" validate:"required,max=1000"`
	SpecificFocusAreas  []string `json:"specific_focus_areas" validate:"omitempty,dive,max=200"`
}

// EscalateInvestigationSLAUri represents the URI parameter for escalating investigation SLA breach
// POST /claims/death/{id}/investigation/escalate-sla-breach
type EscalateInvestigationSLAUri struct {
	ClaimID string `uri:"id" validate:"required"`
}

// AssignManualReviewRequest represents the request for assigning claim for manual review
// POST /claims/death/{id}/manual-review/assign
type AssignManualReviewRequest struct {
	ClaimID    string `uri:"id" validate:"required"`
	ReviewerID string `json:"reviewer_id" validate:"required"`
	Priority   string `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH URGENT"`
}

// RejectClaimForFraudRequest represents the request for rejecting claim based on fraud
// POST /claims/death/{id}/reject-fraud
type RejectClaimForFraudRequest struct {
	ClaimID              string  `uri:"id" validate:"required"`
	FraudEvidence        string  `json:"fraud_evidence" validate:"required,max=5000"`
	InvestigationReportID string `json:"investigation_report_id" validate:"required"`
	LegalActionRequired  bool    `json:"legal_action_required"`
}

// CheckAppealEligibilityUri represents the URI parameter for checking appeal eligibility
// GET /claims/death/{claim_id}/appeal-eligibility
// Reference: BR-CLM-DC-005 (90-day window)
type CheckAppealEligibilityUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// GetAppellateAuthorityUri represents the URI parameter for getting appellate authority
// GET /claims/death/{claim_id}/appellate-authority
type GetAppellateAuthorityUri struct {
	ClaimID string `uri:"claim_id" validate:"required"`
}

// SubmitAppealRequest represents the request for submitting appeal
// POST /claims/death/{claim_id}/appeal
// Reference: BR-CLM-DC-005 (90-day window), BR-CLM-DC-006 (45-day SLA)
type SubmitAppealRequest struct {
	ClaimID             string  `uri:"claim_id" validate:"required"`
	AppealGrounds       string  `json:"appeal_grounds" validate:"required,max=2000"`
	AdditionalDocuments *string `json:"additional_documents,omitempty" validate:"omitempty,max=1000"`
	AppealType          string `json:"appeal_type" validate:"required,oneof=RECONSIDERATION APPELLATE_AUTHORITY OMBUDSMAN"`
}

// AppealIDUri represents the URI parameter for appeal_id
type AppealIDUri struct {
	ClaimID  string `uri:"claim_id" validate:"required"`
	AppealID string `uri:"appeal_id" validate:"required"`
}

// RecordAppealDecisionRequest represents the request for recording appellate authority decision
// POST /claims/death/{claim_id}/appeal/{appeal_id}/decision
type RecordAppealDecisionRequest struct {
	ClaimID        string  `uri:"claim_id" validate:"required"`
	AppealID       string `uri:"appeal_id" validate:"required"`
	Decision       string `json:"decision" validate:"required,oneof=APPEAL_ACCEPTED APPEAL_REJECTED PARTIAL_ACCEPTANCE"`
	ReasonedOrder  string  `json:"reasoned_order" validate:"required,max=5000"`
	ModificationDetails *string `json:"modification_details,omitempty" validate:"omitempty,max=2000"`
}
