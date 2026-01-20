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
	ClaimID         string  `uri:"claim_id" validate:"required"`
	DocumentType    string  `json:"document_type" validate:"required"`
	FileName        string  `json:"file_name" validate:"required"`
	FileContent     []byte  `json:"file_content" validate:"required"`
	DocumentName    string  `json:"document_name" validate:"required"`
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
	ConditionApproval   *string `json:"condition_approval,omitempty"`   // If approved with conditions
	InvestigationWaiver *string `json:"investigation_waiver,omitempty"` // Supervisor waiver for investigation
}

// RejectClaimRequest represents the request for rejecting death claim
// POST /claims/death/{claim_id}/reject
type RejectClaimRequest struct {
	ClaimID                  string  `uri:"claim_id" validate:"required"`
	RejectionReason          string  `json:"rejection_reason" validate:"required,oneof=FRAUD_DETECTED MATERIAL_SUPPRESSION POLICY_INVALID INCOMPLETE_DOCUMENTS INVESTIGATION_ADVERSE OTHER"`
	DetailedJustification    string  `json:"detailed_justification" validate:"required,max=2000"`
	AppealRightsCommunicated bool    `json:"appeal_rights_communicated"`
	InvestigationReportID    *string `json:"investigation_report_id,omitempty"` // Required if rejection due to investigation
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
	PaymentMode    string  `json:"payment_mode" validate:"required,oneof=AUTO_NEFT POSB_TRANSFER CHEQUE"`
	PaymentDetails *string `json:"payment_details,omitempty"` // Additional payment details if needed
}

// CloseClaimRequest represents the request for closing claim
// POST /claims/death/{claim_id}/close
type CloseClaimRequest struct {
	ClaimID          string `uri:"claim_id" validate:"required"`
	ClosureReason    string `json:"closure_reason" validate:"required,oneof=PAYMENT_COMPLETED CLAIM_WITHDRAWN APPEAL_REJECTED"`
	ArchivalRequired bool   `json:"archival_required"`
}

// CancelClaimRequest represents the request for cancelling claim (claimant withdrawal)
// POST /claims/death/{claim_id}/cancel
type CancelClaimRequest struct {
	ClaimID            string `uri:"claim_id" validate:"required"`
	CancellationReason string `json:"cancellation_reason" validate:"required,max=500"`
	RequestedBy        string `json:"requested_by" validate:"required"`
	CanResubmit        bool   `json:"can_resubmit"`
}

// ReturnClaimRequest represents the request for returning claim to claimant
// POST /claims/death/{claim_id}/return
type ReturnClaimRequest struct {
	ClaimID              string  `uri:"claim_id" validate:"required"`
	ReturnReason         string  `json:"return_reason" validate:"required,oneof=DOCUMENTS_NOT_RECEIVED SLA_EXPIRED"`
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
	Status              string `form:"status" validate:"omitempty"`
	ClaimType           string `form:"claim_type" validate:"omitempty"`
	PolicyID            string `form:"policy_id" validate:"omitempty"`
	CustomerID          string `form:"customer_id" validate:"omitempty"`
	ClaimantName        string `form:"claimant_name" validate:"omitempty"`
	StartDate           string `form:"start_date" validate:"omitempty"` // YYYY-MM-DD
	EndDate             string `form:"end_date" validate:"omitempty"`   // YYYY-MM-DD
	InvestigationStatus string `form:"investigation_status" validate:"omitempty,oneof=CLEAR SUSPECT FRAUD"`
	SLAStatus           string `form:"sla_status" validate:"omitempty,oneof=GREEN YELLOW ORANGE RED"`
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
	ClaimID         string `uri:"claim_id" validate:"required"`
	InvestigationID string `uri:"investigation_id" validate:"required"`
}

// AssignInvestigationRequest represents the request for assigning investigation officer
// POST /claims/death/{claim_id}/investigation/assign-officer
// Reference: BR-CLM-DC-002
type AssignInvestigationRequest struct {
	ClaimID        string `uri:"claim_id" validate:"required"`
	InvestigatorID string `json:"investigator_id" validate:"required"`
	Priority       string `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	AssignmentType string `json:"assignment_type" validate:"omitempty,oneof=AUTO MANUAL"` // Default: AUTO
}

// InvestigationProgressRequest represents the request for submitting investigation progress
// POST /claims/death/{claim_id}/investigation/{investigation_id}/progress-update
type InvestigationProgressRequest struct {
	ClaimID         string  `uri:"claim_id" validate:"required"`
	InvestigationID string  `uri:"investigation_id" validate:"required"`
	ProgressNotes   string  `json:"progress_notes" validate:"required,max=2000"`
	NextSteps       *string `json:"next_steps,omitempty" validate:"omitempty,max=500"`
	Percentage      int     `json:"percentage" validate:"required,min=0,max=100"`
}

// SubmitInvestigationReportRequest represents the request for submitting final investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/submit-report
type SubmitInvestigationReportRequest struct {
	ClaimID         string  `uri:"claim_id" validate:"required"`
	InvestigationID string  `uri:"investigation_id" validate:"required"`
	ReportOutcome   string  `json:"report_outcome" validate:"required,oneof=CLEAR SUSPECT FRAUD"`
	Findings        string  `json:"findings" validate:"required,max=5000"`
	Evidence        *string `json:"evidence,omitempty" validate:"omitempty,max=2000"`
	Recommendation  string  `json:"recommendation" validate:"required,max=1000"`
	ReportDate      string  `json:"report_date" validate:"required"` // YYYY-MM-DD
}

// ReviewInvestigationReportRequest represents the request for reviewing investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/review
type ReviewInvestigationReportRequest struct {
	ClaimID               string  `uri:"claim_id" validate:"required"`
	InvestigationID       string  `uri:"investigation_id" validate:"required"`
	ReviewDecision        string  `json:"review_decision" validate:"required,oneof=ACCEPT REINVESTIGATE ESCALATE"`
	ReviewerRemarks       string  `json:"reviewer_remarks" validate:"required,max=2000"`
	ReinvestigationReason *string `json:"reinvestigation_reason,omitempty" validate:"omitempty,max=1000"`
}

// TriggerReinvestigationRequest represents the request for triggering reinvestigation
// POST /claims/death/{id}/investigation/trigger-reinvestigation
// Reference: BR-CLM-DC-013 (max 2 times)
type TriggerReinvestigationRequest struct {
	ClaimID               string   `uri:"id" validate:"required"`
	ReinvestigationReason string   `json:"reinvestigation_reason" validate:"required,max=1000"`
	SpecificFocusAreas    []string `json:"specific_focus_areas" validate:"omitempty,dive,max=200"`
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
	ClaimID               string `uri:"id" validate:"required"`
	FraudEvidence         string `json:"fraud_evidence" validate:"required,max=5000"`
	InvestigationReportID string `json:"investigation_report_id" validate:"required"`
	LegalActionRequired   bool   `json:"legal_action_required"`
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
	AppealType          string  `json:"appeal_type" validate:"required,oneof=RECONSIDERATION APPELLATE_AUTHORITY OMBUDSMAN"`
}

// AppealIDUri represents the URI parameter for appeal_id
type AppealIDUri struct {
	ClaimID  string `uri:"claim_id" validate:"required"`
	AppealID string `uri:"appeal_id" validate:"required"`
}

// RecordAppealDecisionRequest represents the request for recording appellate authority decision
// POST /claims/death/{claim_id}/appeal/{appeal_id}/decision
type RecordAppealDecisionRequest struct {
	ClaimID             string  `uri:"claim_id" validate:"required"`
	AppealID            string  `uri:"appeal_id" validate:"required"`
	Decision            string  `json:"decision" validate:"required,oneof=APPEAL_ACCEPTED APPEAL_REJECTED PARTIAL_ACCEPTANCE"`
	ReasonedOrder       string  `json:"reasoned_order" validate:"required,max=5000"`
	ModificationDetails *string `json:"modification_details,omitempty" validate:"omitempty,max=2000"`
}

// ========================================
// MATURITY CLAIM REQUEST DTOS
// ========================================

// SendMaturityIntimationBatchRequest represents the request for sending maturity intimation batch
// POST /claims/maturity/send-intimation-batch
// Reference: FR-CLM-MC-002, BR-CLM-MC-002 (60 days before maturity)
type SendMaturityIntimationBatchRequest struct {
	MaturityDateFrom string   `json:"maturity_date_from" validate:"required" example:"2024-03-01"`
	MaturityDateTo   string   `json:"maturity_date_to" validate:"required" example:"2024-03-31"`
	Channels         []string `json:"channels" validate:"required,dive,oneof=SMS EMAIL WHATSAPP"`
}

// GenerateMaturityDueReportRequest represents the request for generating maturity due report
// POST /claims/maturity/generate-due-report
type GenerateMaturityDueReportRequest struct {
	ReportMonth int `json:"report_month" validate:"required,min=1,max=12"`
	ReportYear  int `json:"report_year" validate:"required,min=2020,max=2100"`
}

// GetMaturityPreFillDataRequest represents the request for getting pre-filled data
// GET /claims/maturity/pre-fill-data
type GetMaturityPreFillDataRequest struct {
	PolicyID string `uri:"policy_id" validate:"required"`
	Token    string `uri:"token" validate:"required"`
}

// SubmitMaturityClaimRequest represents the request for submitting maturity claim
// POST /claims/maturity/submit
// Reference: BR-CLM-MC-001 (7 days SLA)
type SubmitMaturityClaimRequest struct {
	PolicyID             string   `json:"policy_id" validate:"required"`
	ClaimantName         string   `json:"claimant_name" validate:"required,max=200"`
	ClaimantRelationship string   `json:"claimant_relationship" validate:"required,max=100"`
	ClaimantMobile       string   `json:"claimant_mobile" validate:"required,len=10"`
	ClaimantEmail        string   `json:"claimant_email" validate:"required,email"`
	DisbursementMode     string   `json:"disbursement_mode" validate:"required,oneof=NEFT POSB CHEQUE"`
	BankAccountNumber    string   `json:"bank_account_number" validate:"required,max=50"`
	BankIFSC             string   `json:"bank_ifsc" validate:"required,len=11"`
	BankAccountType      string   `json:"bank_account_type" validate:"required,oneof=Savings Current"`
	Documents            []string `json:"documents" validate:"required,dive,max=100"`
	IsNRI                bool     `json:"is_nri"`
	NRICountry           *string  `json:"nri_country,omitempty" validate:"omitempty,max=100"`
	IsPANAvailable       bool     `json:"is_pan_available"`
	PANNumber            *string  `json:"pan_number,omitempty" validate:"omitempty,len=10"`
	Acknowledgement      bool     `json:"acknowledgement" validate:"required"`
}

// ExtractOCRDataRequest represents the request for extracting OCR data
// POST /claims/maturity/{claim_id}/extract-ocr-data
type ExtractOCRDataRequest struct {
	ClaimID     string   `uri:"claim_id" validate:"required"`
	DocumentIDs []string `json:"document_ids" validate:"required,dive,max=100"`
}

// QCVerifyMaturityClaimRequest represents the request for QC verification
// POST /claims/maturity/{claim_id}/qc-verify
type QCVerifyMaturityClaimRequest struct {
	ClaimID     string                 `uri:"claim_id" validate:"required"`
	QCStatus    string                 `json:"qc_status" validate:"required,oneof=APPROVED REJECTED CORRECTIONS_REQUIRED"`
	Corrections map[string]interface{} `json:"corrections,omitempty"`
	QCRemarks   *string                `json:"qc_remarks,omitempty" validate:"omitempty,max=2000"`
}

// ApproveMaturityClaimRequest represents the request for approving maturity claim
// POST /claims/maturity/{claim_id}/approve
type ApproveMaturityClaimRequest struct {
	ClaimID             string  `uri:"claim_id" validate:"required"`
	ApprovalStatus      string  `json:"approval_status" validate:"required,oneof=APPROVED REJECTED"`
	ApprovalAmount      float64 `json:"approval_amount" validate:"required,gt=0"`
	ApprovalRemarks     string  `json:"approval_remarks" validate:"required,max=2000"`
	ApproverID          string  `json:"approver_id" validate:"required"`
	ApprovalLevel       string  `json:"approval_level" validate:"required,oneof=LEVEL_1 LEVEL_2 LEVEL_3 LEVEL_4"`
	CalculationOverride bool    `json:"calculation_override"`
	OverrideReason      *string `json:"override_reason,omitempty" validate:"omitempty,max=1000"`
}

// DisburseMaturityClaimRequest represents the request for disbursement
// POST /claims/maturity/{claim_id}/disburse
type DisburseMaturityClaimRequest struct {
	ClaimID            string  `uri:"claim_id" validate:"required"`
	DisbursementAmount float64 `json:"disbursement_amount" validate:"required,gt=0"`
	DisbursementMode   string  `json:"disbursement_mode" validate:"required,oneof=NEFT POSB CHEQUE"`
	ReferenceNumber    string  `json:"reference_number" validate:"required,max=100"`
	DisbursementDate   string  `json:"disbursement_date" validate:"required"`
	BankAccountNumber  string  `json:"bank_account_number" validate:"required,max=50"`
	BankIFSC           string  `json:"bank_ifsc" validate:"required,len=11"`
	DisburseTo         string  `json:"disburse_to" validate:"required,max=200"`
	UTRNumber          *string `json:"utr_number,omitempty" validate:"omitempty,max=50"`
	ChequeNumber       *string `json:"cheque_number,omitempty" validate:"omitempty,max=20"`
	ChequeDate         *string `json:"cheque_date,omitempty" validate:"omitempty"`
}

// CloseMaturityClaimRequest represents the request for closing maturity claim
// POST /claims/maturity/{claim_id}/close
type CloseMaturityClaimRequest struct {
	ClaimID       string `uri:"claim_id" validate:"required"`
	ClosureDate   string `json:"closure_date" validate:"required"`
	ClosureReason string `json:"closure_reason" validate:"required,max=1000"`
	ClosedBy      string `json:"closed_by" validate:"required"`
}

// RequestMaturityFeedbackRequest represents the request for requesting feedback
// POST /claims/maturity/{claim_id}/request-feedback
type RequestMaturityFeedbackRequest struct {
	ClaimID     string `uri:"claim_id" validate:"required"`
	Channel     string `json:"channel" validate:"required,oneof=SMS EMAIL WHATSAPP"`
	FeedbackURL string `json:"feedback_url" validate:"required,url"`
}

// ==================== SURVIVAL BENEFIT REQUEST DTOs ====================

// SubmitSurvivalBenefitClaimRequest represents the request for submitting survival benefit claim
// POST /claims/survival-benefit/submit
// Reference: FRS-SB-03, BR-CLM-SB-001 (7 days SLA)
type SubmitSurvivalBenefitClaimRequest struct {
	PolicyID             string   `json:"policy_id" validate:"required"`
	ClaimantName         string   `json:"claimant_name" validate:"required,max=200"`
	ClaimantRelationship string   `json:"claimant_relationship" validate:"required,max=100"`
	ClaimantMobile       string   `json:"claimant_mobile" validate:"required,len=10"`
	ClaimantEmail        string   `json:"claimant_email" validate:"required,email"`
	DisbursementMode     string   `json:"disbursement_mode" validate:"required,oneof=NEFT POSB CHEQUE"`
	BankAccountNumber    string   `json:"bank_account_number" validate:"required,max=50"`
	BankIFSC             string   `json:"bank_ifsc" validate:"required,len=11"`
	BankAccountType      string   `json:"bank_account_type" validate:"required,oneof=Savings Current"`
	Documents            []string `json:"documents" validate:"required,dive,max=100"`
	UseDigiLocker        bool     `json:"use_digiLocker"`
	IsNRI                bool     `json:"is_nri"`
	NRICountry           *string  `json:"nri_country,omitempty" validate:"omitempty,max=100"`
	IsPANAvailable       bool     `json:"is_pan_available"`
	PANNumber            *string  `json:"pan_number,omitempty" validate:"omitempty,len=10"`
	Acknowledgement      bool     `json:"acknowledgement" validate:"required"`
}

// ValidateSBEligibilityRequest represents the request for validating SB eligibility
// POST /claims/survival-benefit/{id}/validate-eligibility
type ValidateSBEligibilityRequest struct {
	ClaimID string `uri:"id" validate:"required"`
}

// ==================== AML/CFT REQUEST DTOs ====================

// DetectAMLTriggerRequest represents the request for detecting AML trigger conditions
// POST /aml/detect-trigger
// Reference: FR-CLM-AML-001, BR-CLM-AML-001 (High Cash Premium Alert)
// Reference: BR-CLM-AML-002 (PAN Mismatch Alert)
// Reference: BR-CLM-AML-003 (Nominee Change Post Death)
type DetectAMLTriggerRequest struct {
	TransactionType   string  `json:"transaction_type" validate:"required,oneof=PREMIUM CLAIM DISBURSEMENT REFUND"`
	TransactionAmount float64 `json:"transaction_amount" validate:"required,gt=0"`
	PolicyID          *string `json:"policy_id,omitempty" validate:"omitempty,max=50"`
	CustomerID        *string `json:"customer_id,omitempty" validate:"omitempty,max=50"`
	PaymentMode       *string `json:"payment_mode,omitempty" validate:"omitempty,oneof=CASH CHEQUE NEFT POSB DD"`
	PANNumber         *string `json:"pan_number,omitempty" validate:"omitempty,len=10"`
	BankAccountNumber *string `json:"bank_account_number,omitempty" validate:"omitempty,max=50"`
	NomineeID         *string `json:"nominee_id,omitempty" validate:"omitempty,max=50"`
	TransactionDate   string  `json:"transaction_date" validate:"required"` // YYYY-MM-DD format
}

// AlertIDUri represents the alert_id URI parameter
type AlertIDUri struct {
	AlertID string `uri:"alert_id" validate:"required"`
}

// ReviewAMLAlertRequest represents the request for reviewing AML alert
// POST /aml/{alert_id}/review
// Reference: BR-CLM-AML-004 (Risk Scoring Algorithm)
// Reference: BR-CLM-AML-005 (Alert Review)
type ReviewAMLAlertRequest struct {
	AlertID         string  `uri:"alert_id" validate:"required"`
	ReviewDecision  string  `json:"review_decision" validate:"required,oneof=CLEAR FILE_STR FILE_CTR BLOCK_TRANSACTION ESCALATE"`
	OfficerRemarks  string  `json:"officer_remarks" validate:"required,max=2000"`
	EscalationLevel *string `json:"escalation_level,omitempty" validate:"omitempty,oneof=LEVEL_1 LEVEL_2 LEVEL_3"`
	OfficerID       string  `json:"officer_id" validate:"required"`
}

// FileAMLReportRequest represents the request for filing STR/CTR with regulatory authorities
// POST /aml/{alert_id}/file-report
// Reference: BR-CLM-AML-006 (STR Filing Within 7 Days)
// Reference: BR-CLM-AML-007 (CTR Filing Monthly)
type FileAMLReportRequest struct {
	AlertID         string   `uri:"alert_id" validate:"required"`
	ReportType      string   `json:"report_type" validate:"required,oneof=STR CTR CCR NTR"`
	ReportingAgency string   `json:"reporting_agency" validate:"required,oneof=FINNET FINGATE"`
	FilingReference string   `json:"filing_reference" validate:"required,max=100"`
	ReportDetails   string   `json:"report_details" validate:"required,max=5000"`
	Attachments     []string `json:"attachments,omitempty" validate:"omitempty,dive,max=200"`
	SupportingDocs  []string `json:"supporting_docs,omitempty" validate:"omitempty,dive,max=200"`
	FiledBy         string   `json:"filed_by" validate:"required"`
	FilingDate      string   `json:"filing_date" validate:"required"` // YYYY-MM-DD format
	Acknowledgement bool     `json:"acknowledgement" validate:"required"`
}

// ==================== BANKING & PAYMENT REQUEST DTOs ====================

// BankValidationRequest represents the request for bank account validation
// POST /banking/validate-account
// POST /banking/validate-account-cbs
// POST /banking/validate-account-pfms
// POST /banking/penny-drop
// Reference: BR-CLM-DC-010 (Payment Disbursement Workflow)
// Reference: Integration with CBS API and PFMS API
type BankValidationRequest struct {
	AccountNumber     string  `json:"account_number" validate:"required,max=50"`
	IFSCCode          string  `json:"ifsc_code" validate:"required,max=20"`
	AccountHolderName string  `json:"account_holder_name" validate:"required,max=200"`
	ValidationMethod  *string `json:"validation_method,omitempty" validate:"omitempty,oneof=CBS_API PFMS_API PENNY_DROP"`
	ClaimID           *string `json:"claim_id,omitempty" validate:"omitempty,max=50"`
}

// InitiateNEFTTransferRequest represents the request for NEFT transfer
// POST /banking/neft-transfer
// Reference: BR-CLM-DC-010 (Disbursement Workflow)
type InitiateNEFTTransferRequest struct {
	AccountNumber   string  `json:"account_number" validate:"required,max=50"`
	IFSCCode        string  `json:"ifsc_code" validate:"required,max=20"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	BeneficiaryName string  `json:"beneficiary_name" validate:"required,max=200"`
	ReferenceID     *string `json:"reference_id,omitempty" validate:"omitempty,max=50"`
	ClaimID         *string `json:"claim_id,omitempty" validate:"omitempty,max=50"`
	PaymentMode     *string `json:"payment_mode,omitempty" validate:"omitempty,oneof=NEFT RTGS IMPS"`
	Remarks         *string `json:"remarks,omitempty" validate:"omitempty,max=500"`
}

// ReconcilePaymentsRequest represents the request for daily payment reconciliation
// POST /banking/payment-reconciliation
// Reference: BR-CLM-PAY-001 (Daily Reconciliation)
type ReconcilePaymentsRequest struct {
	ReconciliationDate string `json:"reconciliation_date" validate:"required"` // YYYY-MM-DD format
	IncludeFailed      *bool  `json:"include_failed,omitempty" validate:"omitempty"`
	IncludePending     *bool  `json:"include_pending,omitempty" validate:"omitempty"`
}

// PaymentIDUri represents the payment_id URI parameter
// GET /banking/payment-status/{payment_id}
type PaymentIDUri struct {
	PaymentID string `uri:"payment_id" validate:"required"`
}

// PaymentConfirmationWebhookRequest represents the webhook request from banking gateway
// POST /webhooks/banking/payment-confirmation
// Reference: Integration with Banking Gateway
type PaymentConfirmationWebhookRequest struct {
	PaymentID     string   `json:"payment_id" validate:"required"`
	Status        string   `json:"status" validate:"required,oneof=SUCCESS FAILED PENDING"`
	TransactionID string   `json:"transaction_id" validate:"required,max=100"`
	Amount        *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Timestamp     string   `json:"timestamp" validate:"required"` // YYYY-MM-DD HH:MM:SS format
	FailureReason *string  `json:"failure_reason,omitempty" validate:"omitempty,max=500"`
	BankReference *string  `json:"bank_reference,omitempty" validate:"omitempty,max=100"`
}

// GeneratePaymentVoucherRequest represents the request for generating payment voucher
// POST /banking/generate-voucher
// Reference: BR-CLM-PAY-002 (Voucher Generation)
type GeneratePaymentVoucherRequest struct {
	ClaimID      string  `json:"claim_id" validate:"required,max=50"`
	PaymentID    string  `json:"payment_id" validate:"required,max=50"`
	VoucherType  *string `json:"voucher_type,omitempty" validate:"omitempty,oneof=PAYMENT RECEIPT JOURNAL"`
	IncludeStamp *bool   `json:"include_stamp,omitempty" validate:"omitempty"`
}

// CBSAccountValidationRequest represents the request for CBS account validation
// This is used internally by the BankingHandler to call CBS API
type CBSAccountValidationRequest struct {
	AccountNumber     string `json:"account_number" validate:"required"`
	IFSCCode          string `json:"ifsc_code" validate:"required"`
	AccountHolderName string `json:"account_holder_name" validate:"required"`
}

// CBSPennyDropRequest represents the request for CBS penny drop test
// This is used internally by the BankingHandler to call CBS API
type CBSPennyDropRequest struct {
	AccountNumber     string  `json:"account_number" validate:"required"`
	IFSCCode          string  `json:"ifsc_code" validate:"required"`
	AccountHolderName string  `json:"account_holder_name" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	ReferenceID       string  `json:"reference_id" validate:"required"`
}

// ========================================
// FREE LOOK CANCELLATION - REQUEST DTOS
// ========================================

// TrackPolicyBondRequest represents the request for tracking policy bond delivery
// POST /policy-bond/track
// Reference: FR-CLM-BOND-001, BR-CLM-BOND-001
type TrackPolicyBondRequest struct {
	PolicyID       string  `json:"policy_id" validate:"required"`
	BondType       string  `json:"bond_type" validate:"required,oneof=PHYSICAL ELECTRONIC"`
	DispatchDate   string  `json:"dispatch_date" validate:"required"` // YYYY-MM-DD format
	DispatchNumber *string `json:"dispatch_number,omitempty" validate:"omitempty,max=50"`
	CourierName    *string `json:"courier_name,omitempty" validate:"omitempty,max=100"`
}

// UpdateBondDeliveryRequest represents the request for updating bond delivery status
// POST /policy-bond/{bond_id}/delivery-status
// Reference: FR-CLM-BOND-002, BR-CLM-BOND-002
type UpdateBondDeliveryRequest struct {
	BondID           string  `json:"bond_id" validate:"required"`
	DeliveryStatus   string  `json:"delivery_status" validate:"required,oneof=DELIVERED UNDELIVERED PENDING RETURNED"`
	DeliveryDate     string  `json:"delivery_date" validate:"required"` // YYYY-MM-DD format
	DeliveryMethod   *string `json:"delivery_method,omitempty" validate:"omitempty,oneof=POST_OFFICE COURIER DIGILOCKER EMAIL"`
	PODImageURL      *string `json:"pod_image_url,omitempty" validate:"omitempty,url"`
	ReceiverName     *string `json:"receiver_name,omitempty" validate:"omitempty,max=100"`
	ReceiverRelation *string `json:"receiver_relation,omitempty" validate:"omitempty,max=50"`
	Remarks          *string `json:"remarks,omitempty" validate:"omitempty,max=500"`
}

// BondIDUri represents the bond_id URI parameter
// GET /policy-bond/{bond_id}/details
type BondIDUri struct {
	BondID string `uri:"bond_id" validate:"required"`
}

// PolicyIDUri represents the policy_id URI parameter
// GET /policy-bond/policy/{policy_id}
// GET /freelook/policy/{policy_id}/eligibility
type PolicyIDUri struct {
	PolicyID string `uri:"policy_id" validate:"required"`
}

// SubmitFreeLookCancellationRequest represents the request for submitting free look cancellation
// POST /freelook/cancellation/submit
// Reference: FR-CLM-FL-002, BR-CLM-BOND-001, VR-CLM-FL-001
type SubmitFreeLookCancellationRequest struct {
	PolicyID           string   `json:"policy_id" validate:"required"`
	CancellationReason string   `json:"cancellation_reason" validate:"required,max=500"`
	Channel            string   `json:"channel" validate:"required,oneof=ONLINE PORTAL POST_OFFICE CPGRAMS EMAIL PHONE"`
	CancellationDate   string   `json:"cancellation_date" validate:"required"` // YYYY-MM-DD format
	BondSubmitted      bool     `json:"bond_submitted" validate:"required"`
	BondType           *string  `json:"bond_type,omitempty" validate:"omitempty,oneof=PHYSICAL ELECTRONIC"`
	DocumentURLs       []string `json:"document_urls" validate:"required,min=1"`
	ClaimantName       string   `json:"claimant_name" validate:"required"`
	ClaimantPhone      *string  `json:"claimant_phone,omitempty" validate:"omitempty,len=10"`
	ClaimantEmail      *string  `json:"claimant_email,omitempty" validate:"omitempty,email"`
	BankAccountNumber  *string  `json:"bank_account_number,omitempty" validate:"omitempty,max=30"`
	BankIFSCCode       *string  `json:"bank_ifsc_code,omitempty" validate:"omitempty,max=20"`
	RefundAmount       *float64 `json:"refund_amount,omitempty" validate:"omitempty,gt=0"`
}

// CancellationIDUri represents the cancellation_id URI parameter
// GET /freelook/cancellation/{cancellation_id}/details
type CancellationIDUri struct {
	CancellationID string `uri:"cancellation_id" validate:"required"`
}

// ReviewFreeLookCancellationRequest represents the request for reviewing free look cancellation (maker-checker)
// POST /freelook/cancellation/{cancellation_id}/review
// Reference: BR-CLM-BOND-004 (Maker-Checker Workflow)
type ReviewFreeLookCancellationRequest struct {
	CancellationID string   `json:"cancellation_id" validate:"required"`
	ReviewAction   string   `json:"review_action" validate:"required,oneof=APPROVE REJECT"`
	ReviewComments string   `json:"review_comments" validate:"required,max=1000"`
	CheckedBy      string   `json:"checked_by" validate:"required"`
	OverrideAmount *float64 `json:"override_amount,omitempty" validate:"omitempty,gt=0"`
	OverrideReason *string  `json:"override_reason,omitempty" validate:"omitempty,max=500"`
}

// ProcessFreeLookRefundRequest represents the request for processing free look refund
// POST /freelook/cancellation/{cancellation_id}/process-refund
// Reference: FR-CLM-FL-003, BR-CLM-BOND-003
type ProcessFreeLookRefundRequest struct {
	CancellationID  string `json:"cancellation_id" validate:"required"`
	RefundMode      string `json:"refund_mode" validate:"required,oneof=NEFT RTGS POSB CHEQUE"`
	ReferenceNumber string `json:"reference_number" validate:"required,max=50"`
	ProcessedBy     string `json:"processed_by" validate:"required"`
	FinanceApproved bool   `json:"finance_approved" validate:"required"`
}

// FreeLookCancellationIDUri represents the cancellation_id URI parameter for refund status
// GET /freelook/cancellation/{cancellation_id}/refund-status
type FreeLookCancellationIDUri struct {
	CancellationID string `uri:"cancellation_id" validate:"required"`
}

// ==================== OMBUDSMAN REQUEST DTOs ====================

// SubmitOmbudsmanComplaintRequest represents a new ombudsman complaint submission
// FR-CLM-OMB-001: Complaint Intake & Registration
// BR-CLM-OMB-001: Admissibility checks
type SubmitOmbudsmanComplaintRequest struct {
	// Complainant Details
	ComplainantName    string `json:"complainant_name" validate:"required,max=200"`
	ComplainantAddress string `json:"complainant_address" validate:"required,max=500"`
	ComplainantMobile  string `json:"complainant_mobile" validate:"required,len=10"`
	ComplainantEmail   string `json:"complainant_email" validate:"omitempty,email,max=100"`
	ComplainantRole    string `json:"complainant_role" validate:"required,oneof=POLICYHOLDER NOMINEE LEGAL_HEIR ASSIGNEE AUTHORIZED_REPRESENTATIVE"`
	LanguagePreference string `json:"language_preference" validate:"required,oneof=ENGLISH HINDI"`
	IDProofType        string `json:"id_proof_type" validate:"required,oneof=AADHAAR PAN PASSPORT OTHER"`
	IDProofNumber      string `json:"id_proof_number" validate:"required,max=50"`

	// Policy/Claim Details
	PolicyNumber string `json:"policy_number" validate:"required,max=50"`
	ClaimNumber  string `json:"claim_number,omitempty" validate:"omitempty,max=50"`
	PolicyType   string `json:"policy_type" validate:"required,oneof=PLI RPLI"`
	AgentName    string `json:"agent_name,omitempty" validate:"omitempty,max=200"`
	AgentBranch  string `json:"agent_branch,omitempty" validate:"omitempty,max=100"`

	// Complaint Details
	IncidentDate       string  `json:"incident_date" validate:"required"`
	RepresentationDate string  `json:"representation_date" validate:"required"` // Date representation made to insurer
	ComplaintCategory  string  `json:"complaint_category" validate:"required,oneof=CLAIM_DELAY PARTIAL_REPUDIATION FULL_REPUDIATION PREMIUM_DISPUTE POLICY_MISREPRESENTATION NON_ISSUANCE_OF_POLICY POLICY_SERVICING OTHER"`
	IssueDescription   string  `json:"issue_description" validate:"required,max=5000"`
	ReliefSought       string  `json:"relief_sought" validate:"required,max=2000"`
	ClaimValue         float64 `json:"claim_value,omitempty" validate:"omitempty,min=0"`

	// Additional Information
	ParallelLitigation bool     `json:"parallel_litigation" validate:"required"` // BR-CLM-OMB-001: No parallel litigation check
	IsEmergency        bool     `json:"is_emergency" validate:"required"`
	Channel            string   `json:"channel" validate:"required,oneof=WEB_PORTAL MOBILE_APP EMAIL OFFLINE_FORM WALK_IN"`
	AttachmentIDs      []string `json:"attachment_ids,omitempty" validate:"omitempty,dive,max=50"`
}

// ComplaintIDUri represents the complaint_id URI parameter
// GET /ombudsman/{complaint_id}/details
// POST /ombudsman/{complaint_id}/assign
// POST /ombudsman/{complaint_id}/admissibility
// GET /ombudsman/{complaint_id}/timeline
type ComplaintIDUri struct {
	ComplaintID string `uri:"complaint_id" validate:"required"`
}

// AssignOmbudsmanRequest represents assigning an ombudsman to a complaint
// FR-CLM-OMB-002: Jurisdiction Mapping
// BR-CLM-OMB-002: Jurisdiction mapping
// BR-CLM-OMB-003: Conflict of interest screening
type AssignOmbudsmanRequest struct {
	OmbudsmanID     string `json:"ombudsman_id" validate:"required"`
	OmbudsmanCenter string `json:"ombudsman_center" validate:"required"`
	ConflictCheck   bool   `json:"conflict_check" validate:"required"`                     // Whether conflict screening was performed
	OverrideReason  string `json:"override_reason,omitempty" validate:"omitempty,max=500"` // If manual override
}

// ReviewAdmissibilityRequest represents admissibility review decision
// BR-CLM-OMB-001: Admissibility checks
type ReviewAdmissibilityRequest struct {
	Admissible            bool   `json:"admissible" validate:"required"`
	AdmissibilityReason   string `json:"admissibility_reason" validate:"required,max=1000"`
	InadmissibilityReason string `json:"inadmissibility_reason,omitempty" validate:"omitempty,max=1000"` // If not admissible
	ReviewedBy            string `json:"reviewed_by" validate:"required"`
}

// RecordMediationRequest represents recording mediation outcome
// FR-CLM-OMB-003: Hearing Scheduling & Management
// BR-CLM-OMB-004: Mediation recommendation
type RecordMediationRequest struct {
	HearingID           string `json:"hearing_id" validate:"required"`
	MediationDate       string `json:"mediation_date" validate:"required"`
	ConsentToMediate    bool   `json:"consent_to_mediate" validate:"required"`
	MediationSuccessful bool   `json:"mediation_successful" validate:"required"`
	SettlementTerms     string `json:"settlement_terms,omitempty" validate:"omitempty,max=5000"`
	ComplainantAccepted bool   `json:"complainant_accepted" validate:"required"`
	InsurerAccepted     bool   `json:"insurer_accepted" validate:"required"`
	RecordingOfficer    string `json:"recording_officer" validate:"required"`
	Remarks             string `json:"remarks,omitempty" validate:"omitempty,max=2000"`
}

// IssueAwardRequest represents issuing an award (mediation or adjudication)
// FR-CLM-OMB-004: Award Issuance & Enforcement
// BR-CLM-OMB-005: Award issuance with ₹50 lakh cap
type IssueAwardRequest struct {
	AwardType            string   `json:"award_type" validate:"required,oneof=MEDIATION_RECOMMENDATION ADJUDICATION_AWARD"`
	AwardAmount          float64  `json:"award_amount" validate:"required,min=0"`                    // BR-CLM-OMB-005: Must be ≤ ₹50 lakh
	AwardCurrency        string   `json:"award_currency" validate:"required,len=3"`                  // INR
	InterestRate         float64  `json:"interest_rate,omitempty" validate:"omitempty,min=0,max=20"` // Interest % if applicable
	InterestAmount       float64  `json:"interest_amount,omitempty" validate:"omitempty,min=0"`
	TotalAwardAmount     float64  `json:"total_award_amount" validate:"required,min=0"` // Including interest
	AwardReasoning       string   `json:"award_reasoning" validate:"required,max=10000"`
	DigitalSignature     string   `json:"digital_signature" validate:"required"` // Ombudsman's digital signature hash
	DigitalSignatureDate string   `json:"digital_signature_date" validate:"required"`
	SupportingDocuments  []string `json:"supporting_documents,omitempty" validate:"omitempty,dive,max=50"`
	ComplianceDeadline   string   `json:"compliance_deadline" validate:"required"`                        // 30 days from award date (BR-CLM-OMB-006)
	ReminderSchedule     []string `json:"reminder_schedule,omitempty" validate:"omitempty,dive,datetime"` // Days 15, 7, 2 before deadline
	IssuedBy             string   `json:"issued_by" validate:"required"`
}

// RecordComplianceRequest represents recording insurer compliance with award
// BR-CLM-OMB-006: Insurer compliance monitoring
type RecordComplianceRequest struct {
	ComplianceStatus string  `json:"compliance_status" validate:"required,oneof=ACCEPTED PAYMENT_INITIATED PAYMENT_COMPLETED OBJECTION_FILED ESCALATED"`
	ComplianceDate   string  `json:"compliance_date" validate:"required"`
	PaymentReference string  `json:"payment_reference,omitempty" validate:"omitempty,max=100"` // UTR/transaction reference
	PaymentAmount    float64 `json:"payment_amount,omitempty" validate:"omitempty,min=0"`
	ObjectionReason  string  `json:"objection_reason,omitempty" validate:"omitempty,max=5000"` // If objected
	RecordedBy       string  `json:"recorded_by" validate:"required"`
}

// CloseComplaintRequest represents closing an ombudsman complaint
// BR-CLM-OMB-007: Complaint closure & archival
type CloseComplaintRequest struct {
	ClosureReason   string `json:"closure_reason" validate:"required,max=2000"`
	ClosureType     string `json:"closure_type" validate:"required,oneof=AWARD_FULFILLED CONDONATION_WITHDRAWN IRDAI_ESCALATED OTHER"`
	RetentionPeriod int    `json:"retention_period" validate:"required"` // Years: 10 for awards, 7 for mediation (BR-CLM-OMB-007)
	ClosedBy        string `json:"closed_by" validate:"required"`
}

// EscalateToIRDAIRequest represents escalating non-compliance to IRDAI
// BR-CLM-OMB-006: Escalate to IRDAI on breach
type EscalateToIRDAIRequest struct {
	EscalationReason string `json:"escalation_reason" validate:"required,max=2000"`
	BreachDetails    string `json:"breach_details" validate:"required,max=5000"`
	DaysOverdue      int    `json:"days_overdue" validate:"required,min=0"`
	EscalationDate   string `json:"escalation_date" validate:"required"`
	EscalatedBy      string `json:"escalated_by" validate:"required"`
	IRDAIReference   string `json:"irdai_reference,omitempty" validate:"omitempty,max=100"`
}

// ==================== NOTIFICATION REQUEST DTOs ====================

// SendNotificationRequest represents the request for sending a notification
// POST /notifications/send
// Reference: BR-CLM-DC-019 (Communication triggers)
type SendNotificationRequest struct {
	NotificationType string        `json:"notification_type" validate:"required"` // CLAIM_REGISTERED, CLAIM_APPROVED, CLAIM_REJECTED, DOCUMENT_REQUIRED, PAYMENT_PROCESSED, MATURITY_DUE, etc.
	ClaimID          *string       `json:"claim_id,omitempty" validate:"omitempty,max=50"`
	Recipient        RecipientInfo `json:"recipient" validate:"required"`
	Channels         []string      `json:"channels" validate:"required,min=1,dive,oneof=SMS EMAIL WHATSAPP PUSH"`
	CustomMessage    *string       `json:"custom_message,omitempty" validate:"omitempty,max=5000"`
}

// RecipientInfo represents notification recipient details
type RecipientInfo struct {
	Name   string `json:"name" validate:"required,max=200"`
	Mobile string `json:"mobile,omitempty" validate:"omitempty,max=15"`
	Email  string `json:"email,omitempty" validate:"omitempty,email,max=200"`
}

// SendBatchNotificationsRequest represents the request for sending batch notifications
// POST /notifications/send-batch
// Reference: Batch notifications for bulk operations
type SendBatchNotificationsRequest struct {
	Notifications []SendNotificationRequest `json:"notifications" validate:"required,min=1,max=1000,dive"`
}

// GenerateFeedbackLinkRequest represents the request for generating customer feedback link
// POST /feedback/generate-link
// Reference: BR-CLM-DC-020 (Customer feedback)
type GenerateFeedbackLinkRequest struct {
	ClaimID       string  `json:"claim_id" validate:"required,max=50"`
	CustomerEmail *string `json:"customer_email,omitempty" validate:"omitempty,email,max=200"`
	FeedbackType  *string `json:"feedback_type,omitempty" validate:"omitempty,oneof=CLAIM_PROCESSING CUSTOMER_SERVICE DOCUMENT_QUALITY"` // Default: CLAIM_PROCESSING
	ExpiryDays    *int    `json:"expiry_days,omitempty" validate:"omitempty,min=1,max=30"`                                               // Default: 7 days
}

// ==================== POLICY SERVICE REQUEST DTOs ====================

// CheckPolicyClaimEligibilityRequest represents the request for checking policy claim eligibility
// GET /policies/{policy_id}/claim-eligibility
// Reference: INT-CLM-004 (Policy Service Integration)
type CheckPolicyClaimEligibilityRequest struct {
	PolicyID  string `json:"policy_id" validate:"required"`                                                 // From URI parameter
	ClaimType string `json:"claim_type" validate:"required,oneof=DEATH MATURITY SURVIVAL_BENEFIT FREELOOK"` // From query parameter
}

// CalculateFreeLookRefundRequest represents the request for calculating free look refund
// POST /policies/{policy_id}/freelook-refund-calculation
// Reference: BR-CLM-BOND-003 (Refund calculation)
type CalculateFreeLookRefundRequest struct {
	PolicyID         string  `json:"policy_id" validate:"required"` // From URI parameter
	PremiumPaid      float64 `json:"premium_paid" validate:"required,gt=0"`
	CancellationDate string  `json:"cancellation_date" validate:"required"` // YYYY-MM-DD format
	BondType         string  `json:"bond_type" validate:"required,oneof=PHYSICAL ELECTRONIC"`
	DeliveryDate     *string `json:"delivery_date,omitempty"` // YYYY-MM-DD format (required for PHYSICAL bonds)
}

// ==================== VALIDATION SERVICE REQUEST DTOs ====================

// ValidatePANRequest represents the request for validating PAN number
// POST /validate/pan
// Reference: VR-CLM-VAL-001 (PAN validation)
type ValidatePANRequest struct {
	PANNumber string `json:"pan_number" validate:"required,len=10"` // 10-character alphanumeric PAN
}

// ValidateBankAccountServiceRequest represents the request for validating bank account
// POST /validate/bank-account
// Reference: VR-CLM-VAL-002 (Bank account validation)
type ValidateBankAccountServiceRequest struct {
	BankAccountNumber string `json:"bank_account_number" validate:"required,max=50"`
	BankIFSC          string `json:"bank_ifsc" validate:"required,len=11"` // 11-character IFSC code
	AccountHolderName string `json:"account_holder_name,omitempty" validate:"omitempty,max=200"`
	ValidationMethod  string `json:"validation_method,omitempty" validate:"omitempty,oneof=CBS PFMS PENNY_DROP"` // Default: CBS
}

// ValidateDeathDateRequest represents the request for validating death date against policy dates
// POST /validate/death-date
// Reference: BR-CLM-DC-001 (Investigation trigger based on death date)
type ValidateDeathDateRequest struct {
	PolicyID  string `json:"policy_id" validate:"required"`
	DeathDate string `json:"death_date" validate:"required"` // YYYY-MM-DD format
}

// IFSCCodeUri represents the IFSC code URI parameter
// GET /validate/ifsc/{ifsc_code}
// Reference: VR-CLM-VAL-003 (IFSC validation)
type IFSCCodeUri struct {
	IFSCCode string `uri:"ifsc_code" validate:"required,len=11"` // 11-character IFSC code
}

// GetDynamicFormFieldsRequest represents the request for getting dynamic form fields based on death type
// GET /forms/death-claim/fields
// Reference: DFC-001 (Dynamic document checklist)
type GetDynamicFormFieldsRequest struct {
	DeathType string `form:"death_type" validate:"required,oneof=NATURAL ACCIDENTAL UNNATURAL SUICIDE"` // Death type for dynamic fields
}

// ==================== LOOKUP & REFERENCE REQUEST DTOs ====================

// GetDocumentTypesRequest represents the request for getting document types
// GET /lookup/document-types
// Reference: DFC-001 (Dynamic document checklist)
type GetDocumentTypesRequest struct {
	ClaimType        string  `form:"claim_type" validate:"required,oneof=DEATH MATURITY SURVIVAL_BENEFIT FREELOOK"`
	DeathType        *string `form:"death_type,omitempty" validate:"omitempty,oneof=NATURAL UNNATURAL ACCIDENTAL SUICIDE"`
	NominationStatus *string `form:"nomination_status,omitempty" validate:"omitempty,oneof=NOMINATED NOT_NOMINATED"`
}

// GetRejectionReasonsRequest represents the request for getting rejection reasons
// GET /lookup/rejection-reasons
// Reference: BR-CLM-DC-020 (Claim rejection with appeal rights)
type GetRejectionReasonsRequest struct {
	ClaimType string `form:"claim_type" validate:"required,oneof=DEATH MATURITY SURVIVAL_BENEFIT FREELOOK"`
}

// GetInvestigationOfficersRequest represents the request for getting investigation officers
// GET /lookup/investigation-officers
// Reference: BR-CLM-DC-002 (Investigation assignment)
type GetInvestigationOfficersRequest struct {
	Jurisdiction  string  `form:"jurisdiction" validate:"required,max=100"`
	Rank          *string `form:"rank,omitempty" validate:"omitempty,max=100"`
	AvailableOnly bool    `form:"available_only"` // Default: true
}

// GetApproversListRequest represents the request for getting approvers
// GET /lookup/approvers
// Reference: BR-CLM-DC-022 (Approval hierarchy)
type GetApproversListRequest struct {
	ClaimAmount float64 `form:"claim_amount" validate:"required,gt=0"`
	Location    string  `form:"location" validate:"required,max=200"`
}
