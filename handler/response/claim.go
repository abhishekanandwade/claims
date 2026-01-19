package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// DeathClaimResponse represents a death claim in API responses
type DeathClaimResponse struct {
	ID                      string     `json:"id"`
	PolicyID                string     `json:"policy_id"`
	CustomerID              string     `json:"customer_id"`
	ClaimNumber             string     `json:"claim_number"`
	ClaimType               string     `json:"claim_type"`
	ClaimDate               string     `json:"claim_date"`
	DeathDate               string     `json:"death_date,omitempty"`
	DeathPlace              string     `json:"death_place,omitempty"`
	DeathType               string     `json:"death_type,omitempty"`
	ClaimantName            string     `json:"claimant_name"`
	ClaimantType            string     `json:"claimant_type,omitempty"`
	ClaimantRelation        string     `json:"claimant_relation,omitempty"`
	ClaimantPhone           string     `json:"claimant_phone,omitempty"`
	ClaimantEmail           string     `json:"claimant_email,omitempty"`
	BankAccountNumber       string     `json:"bank_account_number,omitempty"`
	BankIfscCode            string     `json:"bank_ifsc_code,omitempty"`
	BankAccountHolderName   string     `json:"bank_account_holder_name,omitempty"`
	BankName                string     `json:"bank_name,omitempty"`
	BankVerified            bool       `json:"bank_verified"`
	InvestigationRequired   bool       `json:"investigation_required"`
	InvestigationStatus     string     `json:"investigation_status,omitempty"`
	InvestigatorID          string     `json:"investigator_id,omitempty"`
	ClaimAmount             *float64   `json:"claim_amount,omitempty"`
	ApprovedAmount          *float64   `json:"approved_amount,omitempty"`
	Status                  string     `json:"status"`
	WorkflowState           string     `json:"workflow_state,omitempty"`
	ApproverID              string     `json:"approver_id,omitempty"`
	ApprovalDate            string     `json:"approval_date,omitempty"`
	ApprovalRemarks         string     `json:"approval_remarks,omitempty"`
	RejectionReason         string     `json:"rejection_reason,omitempty"`
	RejectionCode           string     `json:"rejection_code,omitempty"`
	AppealSubmitted         bool       `json:"appeal_submitted"`
	AppealID                string     `json:"appeal_id,omitempty"`
	SLADueDate              string     `json:"sla_due_date"`
	SLABreached             bool       `json:"sla_breached"`
	SLABreachDays           int        `json:"sla_breach_days"`
	SLAStatus               string     `json:"sla_status"`
	PaymentMode             string     `json:"payment_mode,omitempty"`
	PaymentReference        string     `json:"payment_reference,omitempty"`
	TransactionID           string     `json:"transaction_id,omitempty"`
	UTRNumber               string     `json:"utr_number,omitempty"`
	DisbursementDate        string     `json:"disbursement_date,omitempty"`
	ClosureDate             string     `json:"closure_date,omitempty"`
	ClosureReason           string     `json:"closure_reason,omitempty"`
	CreatedAt               string     `json:"created_at"`
	UpdatedAt               string     `json:"updated_at"`
	CreatedBy               string     `json:"created_by"`
	UpdatedBy               string     `json:"updated_by"`
	Version                 int        `json:"version"`
}

// NewDeathClaimResponse converts domain Claim to response DTO
func NewDeathClaimResponse(d domain.Claim) DeathClaimResponse {
	r := DeathClaimResponse{
		ID:              d.ID,
		PolicyID:        d.PolicyID,
		CustomerID:      d.CustomerID,
		ClaimNumber:     d.ClaimNumber,
		ClaimType:       d.ClaimType,
		ClaimDate:       d.ClaimDate.Format("2006-01-02"),
		ClaimantName:    d.ClaimantName,
		BankVerified:    d.BankVerified,
		InvestigationRequired: d.InvestigationRequired,
		Status:          d.Status,
		AppealSubmitted: d.AppealSubmitted,
		SLADueDate:      d.SLADueDate.Format("2006-01-02"),
		SLABreached:     d.SLABreached,
		SLABreachDays:   d.SLABreachDays,
		SLAStatus:       d.SLAStatus,
		CreatedAt:       d.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       d.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:       d.CreatedBy,
		UpdatedBy:       d.UpdatedBy,
		Version:         d.Version,
	}

	// Add optional fields if present
	if d.DeathDate != nil {
		r.DeathDate = d.DeathDate.Format("2006-01-02")
	}
	if d.DeathPlace != nil {
		r.DeathPlace = *d.DeathPlace
	}
	if d.DeathType != nil {
		r.DeathType = *d.DeathType
	}
	if d.ClaimantType != nil {
		r.ClaimantType = *d.ClaimantType
	}
	if d.ClaimantRelation != nil {
		r.ClaimantRelation = *d.ClaimantRelation
	}
	if d.ClaimantPhone != nil {
		r.ClaimantPhone = *d.ClaimantPhone
	}
	if d.ClaimantEmail != nil {
		r.ClaimantEmail = *d.ClaimantEmail
	}
	if d.BankAccountNumber != nil {
		r.BankAccountNumber = *d.BankAccountNumber
	}
	if d.BankIFSCCode != nil {
		r.BankIfscCode = *d.BankIFSCCode
	}
	if d.BankAccountHolderName != nil {
		r.BankAccountHolderName = *d.BankAccountHolderName
	}
	if d.BankName != nil {
		r.BankName = *d.BankName
	}
	if d.InvestigationStatus != nil {
		r.InvestigationStatus = *d.InvestigationStatus
	}
	if d.InvestigatorID != nil {
		r.InvestigatorID = *d.InvestigatorID
	}
	if d.ClaimAmount != nil {
		r.ClaimAmount = d.ClaimAmount
	}
	if d.ApprovedAmount != nil {
		r.ApprovedAmount = d.ApprovedAmount
	}
	if d.WorkflowState != nil {
		r.WorkflowState = *d.WorkflowState
	}
	if d.ApproverID != nil {
		r.ApproverID = *d.ApproverID
	}
	if d.ApprovalDate != nil {
		r.ApprovalDate = d.ApprovalDate.Format("2006-01-02 15:04:05")
	}
	if d.ApprovalRemarks != nil {
		r.ApprovalRemarks = *d.ApprovalRemarks
	}
	if d.RejectionReason != nil {
		r.RejectionReason = *d.RejectionReason
	}
	if d.RejectionCode != nil {
		r.RejectionCode = *d.RejectionCode
	}
	if d.AppealID != nil {
		r.AppealID = *d.AppealID
	}
	if d.PaymentMode != nil {
		r.PaymentMode = *d.PaymentMode
	}
	if d.PaymentReference != nil {
		r.PaymentReference = *d.PaymentReference
	}
	if d.TransactionID != nil {
		r.TransactionID = *d.TransactionID
	}
	if d.UTRNumber != nil {
		r.UTRNumber = *d.UTRNumber
	}
	if d.DisbursementDate != nil {
		r.DisbursementDate = d.DisbursementDate.Format("2006-01-02")
	}
	if d.ClosureDate != nil {
		r.ClosureDate = d.ClosureDate.Format("2006-01-02")
	}
	if d.ClosureReason != nil {
		r.ClosureReason = *d.ClosureReason
	}

	return r
}

// NewDeathClaimsResponse converts slice of domain Claims to response DTOs
func NewDeathClaimsResponse(data []domain.Claim) []DeathClaimResponse {
	res := make([]DeathClaimResponse, 0, len(data))
	for _, d := range data {
		res = append(res, NewDeathClaimResponse(d))
	}
	return res
}

// DeathClaimRegisteredResponse represents the response for death claim registration
// Reference: FR-CLM-DC-001, BR-CLM-DC-001, WF-CLM-DC-001
type DeathClaimRegisteredResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DeathClaimRegistrationData `json:"data"`
}

// DeathClaimRegistrationData contains registration details
type DeathClaimRegistrationData struct {
	ClaimID                  string `json:"claim_id"`
	ClaimNumber              string `json:"claim_number"`
	Status                   string `json:"status"`
	AcknowledgmentNumber     string `json:"acknowledgment_number"`
	InvestigationRequired    bool   `json:"investigation_required"`
	InvestigationTriggerReason string `json:"investigation_trigger_reason,omitempty"`
	WorkflowState            WorkflowStateResponse `json:"workflow_state"`
}

// WorkflowStateResponse represents the current workflow state
type WorkflowStateResponse struct {
	CurrentStep  string   `json:"current_step"`
	NextStep     string   `json:"next_step,omitempty"`
	SLADeadline  string   `json:"sla_deadline"`
	DaysRemaining int     `json:"days_remaining"`
	SLAStatus    string   `json:"sla_status"` // GREEN, YELLOW, RED
	AllowedActions []string `json:"allowed_actions"`
}

// NewWorkflowStateResponse creates workflow state from domain Claim
func NewWorkflowStateResponse(claim domain.Claim, slaDeadline time.Time) WorkflowStateResponse {
	return WorkflowStateResponse{
		CurrentStep:   claim.Status,
		SLADeadline:   slaDeadline.Format("2006-01-02 15:04:05"),
		DaysRemaining: int(slaDeadline.Sub(time.Now()).Hours() / 24),
		SLAStatus:     calculateSLAStatus(claim.SLAStatus, slaDeadline),
	}
}

// calculateSLAStatus determines SLA status based on deadline and status
func calculateSLAStatus(status string, deadline time.Time) string {
	daysRemaining := int(deadline.Sub(time.Now()).Hours() / 24)

	if status == "COMPLETED" || status == "APPROVED" || status == "REJECTED" {
		return "GREEN"
	}

	if daysRemaining < 0 {
		return "RED"
	} else if daysRemaining <= 3 {
		return "YELLOW"
	}
	return "GREEN"
}

// DeathClaimsListResponse represents the response for listing death claims
type DeathClaimsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Data                      []DeathClaimResponse `json:"data"`
}

// ClaimAmountCalculationResponse represents the response for claim amount calculation
// Reference: CALC-001
type ClaimAmountCalculationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimCalculationData `json:"data"`
}

// ClaimCalculationData contains calculated claim amount details
type ClaimCalculationData struct {
	ClaimID              string  `json:"claim_id,omitempty"`
	SumAssured           float64 `json:"sum_assured"`
	AccruedBonuses       float64 `json:"accrued_bonuses"`
	Deductions           float64 `json:"deductions"`
	NetClaimAmount       float64 `json:"net_claim_amount"`
	CalculationBreakdown ClaimCalculationBreakdown `json:"calculation_breakdown"`
}

// ClaimCalculationBreakdown provides detailed calculation breakdown
type ClaimCalculationBreakdown struct {
	SumAssured         float64 `json:"sum_assured"`
	ReversionaryBonus  float64 `json:"reversionary_bonus"`
	TerminalBonus      float64 `json:"terminal_bonus"`
	OutstandingLoan    float64 `json:"outstanding_loan"`
	UnpaidPremiums     float64 `json:"unpaid_premiums"`
	PenalInterest      float64 `json:"penal_interest,omitempty"`
	OtherDeductions    float64 `json:"other_deductions,omitempty"`
}

// DocumentChecklistResponse represents the response for document checklist
// Reference: FR-CLM-DC-002, VR-CLM-DC-001 to VR-CLM-DC-007
type DocumentChecklistResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DocumentChecklistData `json:"data"`
}

// DocumentChecklistData contains document checklist details
type DocumentChecklistData struct {
	ClaimID  string                  `json:"claim_id"`
	Documents []DocumentChecklistItem `json:"documents"`
}

// DocumentChecklistItem represents a single document in checklist
type DocumentChecklistItem struct {
	DocumentType string `json:"document_type"`
	Mandatory    bool   `json:"mandatory"`
	Description  string `json:"description"`
	Uploaded     bool   `json:"uploaded"`
	DocumentID   string `json:"document_id,omitempty"`
	UploadedAt   string `json:"uploaded_at,omitempty"`
}

// DynamicDocumentChecklistResponse represents the response for dynamic document checklist
// Reference: DFC-001
type DynamicDocumentChecklistResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DynamicDocumentChecklistData `json:"data"`
}

// DynamicDocumentChecklistData contains context-aware document checklist
type DynamicDocumentChecklistData struct {
	MandatoryBase    int                     `json:"mandatory_base"`
	ConditionalAdded int                     `json:"conditional_added"`
	TotalMandatory   int                     `json:"total_mandatory"`
	Documents        []DocumentChecklistItem `json:"documents"`
}

// DocumentCompletenessResponse represents the response for document completeness check
type DocumentCompletenessResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DocumentCompletenessData `json:"data"`
}

// DocumentCompletenessData contains document completeness status
type DocumentCompletenessData struct {
	ClaimID            string   `json:"claim_id"`
	DocumentsComplete  bool     `json:"documents_complete"`
	MandatoryCount     int      `json:"mandatory_count"`
	UploadedCount      int      `json:"uploaded_count"`
	VerifiedCount      int      `json:"verified_count"`
	MissingDocuments   []string `json:"missing_documents,omitempty"`
	PendingVerification []string `json:"pending_verification,omitempty"`
}

// BenefitCalculationResponse represents the response for benefit calculation
// Reference: FR-CLM-DC-004, BR-CLM-DC-008
type BenefitCalculationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimCalculationData `json:"data"`
}

// EligibleApproversResponse represents the response for eligible approvers
// Reference: BR-CLM-DC-015
type EligibleApproversResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      EligibleApproversData `json:"data"`
}

// EligibleApproversData contains eligible approvers list
type EligibleApproversData struct {
	ClaimAmount          float64         `json:"claim_amount"`
	RequiredAuthority    string          `json:"required_authority"`
	EligibleApprovers    []ApproverInfo  `json:"eligible_approvers"`
}

// ApproverInfo represents an eligible approver
type ApproverInfo struct {
	UserID           string  `json:"user_id"`
	Name             string  `json:"name"`
	Designation      string  `json:"designation"`
	AuthorityLimit   float64 `json:"authority_limit"`
}

// ApprovalDetailsResponse represents the response for claim approval details
// Reference: Approver dashboard API
type ApprovalDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ApprovalDetailsData `json:"data"`
}

// ApprovalDetailsData contains complete claim details for approval
type ApprovalDetailsData struct {
	ClaimID            string                   `json:"claim_id"`
	PolicyDetails      PolicyDetailsResponse    `json:"policy_details"`
	ClaimantDetails    ClaimantDetailsResponse  `json:"claimant_details"`
	Calculation        ClaimCalculationData     `json:"calculation,omitempty"`
	Documents          []DocumentChecklistItem  `json:"documents"`
	InvestigationStatus string                  `json:"investigation_status,omitempty"`
	RedFlags           []RedFlag                `json:"red_flags,omitempty"`
}

// PolicyDetailsResponse represents policy information
type PolicyDetailsResponse struct {
	PolicyID     string  `json:"policy_id"`
	PolicyNumber string  `json:"policy_number"`
	PolicyType   string  `json:"policy_type"`
	SumAssured   float64 `json:"sum_assured"`
	IssueDate    string  `json:"issue_date"`
	MaturityDate string  `json:"maturity_date"`
	PolicyStatus string  `json:"policy_status"`
}

// ClaimantDetailsResponse represents claimant information
type ClaimantDetailsResponse struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	PAN          string `json:"pan,omitempty"`
	Aadhaar      string `json:"aadhaar,omitempty"`
}

// RedFlag represents a fraud detection indicator
type RedFlag struct {
	FlagType   string `json:"flag_type"`
	Severity   string `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	Description string `json:"description"`
}

// FraudRedFlagsResponse represents the response for fraud red flags
type FraudRedFlagsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      FraudRedFlagsData `json:"data"`
}

// FraudRedFlagsData contains fraud detection results
type FraudRedFlagsData struct {
	ClaimID   string    `json:"claim_id"`
	RiskScore int       `json:"risk_score"` // 0-100
	RedFlags  []RedFlag `json:"red_flags"`
}

// ClaimApprovalResponse represents the response for claim approval
type ClaimApprovalResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimApprovalData `json:"data"`
}

// ClaimApprovalData contains approval details
type ClaimApprovalData struct {
	ClaimID          string `json:"claim_id"`
	ApprovalDecision string `json:"approval_decision"`
	Approver         string `json:"approver"`
	ApprovedAt       string `json:"approved_at"`
	DigitalSignature string `json:"digital_signature,omitempty"`
}

// ClaimRejectionResponse represents the response for claim rejection
type ClaimRejectionResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimRejectionData `json:"data"`
}

// ClaimRejectionData contains rejection details
type ClaimRejectionData struct {
	ClaimID                 string `json:"claim_id"`
	RejectionReason         string `json:"rejection_reason"`
	DetailedJustification   string `json:"detailed_justification"`
	AppealRightsCommunicated bool  `json:"appeal_rights_communicated"`
	AppealDeadline          string `json:"appeal_deadline,omitempty"`
}

// BankValidationResponse represents the response for bank account validation
// Reference: Banking service integration
type BankValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      BankValidationData `json:"data"`
}

// BankValidationData contains bank validation results
type BankValidationData struct {
	Valid               bool    `json:"valid"`
	AccountNumber       string  `json:"account_number"`
	AccountHolderName   string  `json:"account_holder_name"`
	BankName            string  `json:"bank_name,omitempty"`
	ValidationMethod    string  `json:"validation_method"` // CBS, PENNY_DROP
	NameMatchPercentage float64 `json:"name_match_percentage,omitempty"`
}

// ClaimDisbursementResponse represents the response for claim disbursement
// Reference: BR-CLM-DC-010, WF-CLM-DC-004
type ClaimDisbursementResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimDisbursementData `json:"data"`
}

// ClaimDisbursementData contains disbursement details
type ClaimDisbursementData struct {
	PaymentID         string `json:"payment_id"`
	PaymentReference  string `json:"payment_reference"`
	Status            string `json:"status"`
	InitiatedAt       string `json:"initiated_at"`
	PaymentMode       string `json:"payment_mode"`
	Amount            float64 `json:"amount"`
	BankAccountNumber string `json:"bank_account_number"`
	IfscCode          string `json:"ifsc_code"`
	EstimatedCreditDate string `json:"estimated_credit_date,omitempty"`
}

// ClaimCloseResponse represents the response for claim closure
type ClaimCloseResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimCloseData `json:"data"`
}

// ClaimCloseData contains claim closure details
type ClaimCloseData struct {
	ClaimID    string `json:"claim_id"`
	Status     string `json:"status"`
	ClosedAt   string `json:"closed_at"`
	ClosedBy   string `json:"closed_by"`
	Remarks    string `json:"remarks,omitempty"`
}

// ClaimCancelResponse represents the response for claim cancellation
type ClaimCancelResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimCancelData `json:"data"`
}

// ClaimCancelData contains claim cancellation details
type ClaimCancelData struct {
	ClaimID     string `json:"claim_id"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelled_at"`
	CancelledBy string `json:"cancelled_by"`
	Reason      string `json:"reason,omitempty"`
}

// ClaimReturnResponse represents the response for claim return
type ClaimReturnResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ClaimReturnData `json:"data"`
}

// ClaimReturnData contains claim return details
type ClaimReturnData struct {
	ClaimID         string `json:"claim_id"`
	Status          string `json:"status"`
	ReturnedAt      string `json:"returned_at"`
	ReturnedBy      string `json:"returned_by"`
	ReturnReason    string `json:"return_reason"`
	AdditionalRequirements []string `json:"additional_requirements,omitempty"`
}
