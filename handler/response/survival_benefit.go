package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== SURVIVAL BENEFIT RESPONSE DTOs ====================

// SurvivalBenefitClaimRegistrationResponse represents the response for SB claim registration
// Reference: FRS-SB-03, BR-CLM-SB-001
type SurvivalBenefitClaimRegistrationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string `json:"claim_id"`
	ClaimNumber               string `json:"claim_number"`
	AcknowledgementNumber     string `json:"acknowledgement_number"`
	SubmissionDate            string `json:"submission_date"`
	EstimatedSettlementDate   string `json:"estimated_settlement_date"`
	WorkflowState             *WorkflowStateResponse `json:"workflow_state,omitempty"`
}

// SurvivalBenefitClaimResponse represents a survival benefit claim
// Reference: FRS-SB-01, FRS-SB-15
type SurvivalBenefitClaimResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string `json:"claim_id"`
	ClaimNumber               string `json:"claim_number"`
	PolicyID                  string `json:"policy_id"`
	ClaimType                 string `json:"claim_type"`
	Status                    string `json:"status"`

	// Claimant Details
	ClaimantName         string `json:"claimant_name"`
	ClaimantRelation     string `json:"claimant_relation,omitempty"`
	ClaimantPhone        string `json:"claimant_phone,omitempty"`
	ClaimantEmail        string `json:"claimant_email,omitempty"`

	// Financial Details
	SBAmount         *float64 `json:"sb_amount,omitempty"`
	ApprovedAmount   *float64 `json:"approved_amount,omitempty"`
	PaymentMode      string   `json:"payment_mode,omitempty"`

	// Bank Details
	BankAccountNumber     string `json:"bank_account_number,omitempty"`
	BankIFSCCode          string `json:"bank_ifsc_code,omitempty"`
	BankAccountHolderName string `json:"bank_account_holder_name,omitempty"`
	BankName              string `json:"bank_name,omitempty"`
	BankVerified          bool   `json:"bank_verified"`

	// SLA Tracking
	SLADueDate     string `json:"sla_due_date"`
	SLAStatus      string `json:"sla_status"`
	DaysRemaining  *int64  `json:"days_remaining,omitempty"`

	// Dates
	SubmissionDate           string  `json:"submission_date"`
	ApprovalDate             string  `json:"approval_date,omitempty"`
	DisbursementDate         string  `json:"disbursement_date,omitempty"`
	EstimatedSettlementDate  string  `json:"estimated_settlement_date"`

	// Workflow
	InvestigationRequired bool   `json:"investigation_required"`
	DocumentCompleteness  string `json:"document_completeness,omitempty"`

	// DigiLocker
	UseDigiLocker bool `json:"use_digiLocker"`
}

// SurvivalBenefitClaimsListResponse represents a paginated list of SB claims
type SurvivalBenefitClaimsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Claims                    []SurvivalBenefitClaimResponse `json:"claims"`
}

// SBEligibilityValidationResponse represents the response for SB eligibility validation
// Reference: FRS-SB-02
type SBEligibilityValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Eligible                  bool     `json:"eligible"`
	SBDueDate                 string   `json:"sb_due_date"`
	SBAmount                  float64  `json:"sb_amount"`
	EligibilityReasons        []string `json:"eligibility_reasons"`
	IneligibilityReasons      []string `json:"ineligibility_reasons,omitempty"`
	PolicyDetails             *PolicyDetailsResponse `json:"policy_details,omitempty"`
	ClaimantDetails           *ClaimantDetailsResponse `json:"claimant_details,omitempty"`
	NextSBDueDate             *string `json:"next_sb_due_date,omitempty"`
}

// SurvivalBenefitPreFillDataResponse represents the pre-fill data for SB claim submission
// Reference: FRS-SB-03
type SurvivalBenefitPreFillDataResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyDetails             PolicyDetailsResponse     `json:"policy_details"`
	ClaimantDetails           ClaimantDetailsResponse   `json:"claimant_details"`
	BankDetails               *BankValidationData       `json:"bank_details,omitempty"`
	SBDueDetails              SBDueDetailsResponse      `json:"sb_due_details"`
	DocumentChecklist         []DocumentChecklistItem   `json:"document_checklist"`
}

// SBDueDetailsResponse represents survival benefit due details
type SBDueDetailsResponse struct {
	SBDueDate       string  `json:"sb_due_date"`
	SBAmount        float64 `json:"sb_amount"`
	SBInstallmentNo int     `json:"sb_installment_no"`
	TotalSBDue      float64 `json:"total_sb_due"`
	SBPaid          float64 `json:"sb_paid"`
	SBPending       float64 `json:"sb_pending"`
	IsFirstSB       bool    `json:"is_first_sb"`
	IsFinalSB       bool    `json:"is_final_sb"`
	NextSBDueDate   *string `json:"next_sb_due_date,omitempty"`
}

// SBDueReportResponse represents the survival benefit due report
// Reference: FRS-SB-01
type SBDueReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	ReportDate                string                   `json:"report_date"`
	TotalSBDue                int64                    `json:"total_sb_due"`
	TotalSBAmount             float64                  `json:"total_sb_amount"`
	SBDueByMonth              []SBDueByMonthResponse   `json:"sb_due_by_month"`
	SBDueByDivision           []SBDueByDivisionResponse `json:"sb_due_by_division"`
}

// SBDueByMonthResponse represents SB due grouped by month
type SBDueByMonthResponse struct {
	Month       string  `json:"month"`
	Year        int     `json:"year"`
	Count       int64   `json:"count"`
	TotalAmount float64 `json:"total_amount"`
}

// SBDueByDivisionResponse represents SB due grouped by division
type SBDueByDivisionResponse struct {
	DivisionCode string  `json:"division_code"`
	DivisionName string  `json:"division_name"`
	Count        int64   `json:"count"`
	TotalAmount  float64 `json:"total_amount"`
}

// ==================== HELPER FUNCTIONS ====================

// NewSurvivalBenefitClaimResponse creates a new SB claim response from domain model
// Reference: FRS-SB-15
func NewSurvivalBenefitClaimResponse(claim domain.Claim) SurvivalBenefitClaimResponse {
	response := SurvivalBenefitClaimResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		ClaimID:              claim.ID,
		ClaimNumber:          claim.ClaimNumber,
		PolicyID:             claim.PolicyID,
		ClaimType:            claim.ClaimType,
		Status:               claim.Status,
		ClaimantName:         claim.ClaimantName,
		BankVerified:         claim.BankVerified,
		InvestigationRequired: claim.InvestigationRequired,
		SLADueDate:           claim.SLADueDate.Format("2006-01-02 15:04:05"),
		SLAStatus:            claim.SLAStatus,
		SubmissionDate:       claim.CreatedAt.Format("2006-01-02 15:04:05"),
		UseDigiLocker:        false, // Set based on claim metadata
	}

	// Optional fields
	if claim.ClaimantRelation != nil {
		response.ClaimantRelation = *claim.ClaimantRelation
	}

	if claim.ClaimantPhone != nil {
		response.ClaimantPhone = *claim.ClaimantPhone
	}

	if claim.ClaimantEmail != nil {
		response.ClaimantEmail = *claim.ClaimantEmail
	}

	if claim.ClaimAmount != nil {
		response.SBAmount = claim.ClaimAmount
	}

	if claim.ApprovedAmount != nil {
		response.ApprovedAmount = claim.ApprovedAmount
	}

	if claim.PaymentMode != nil {
		response.PaymentMode = *claim.PaymentMode
	}

	if claim.BankAccountNumber != nil {
		response.BankAccountNumber = *claim.BankAccountNumber
	}

	if claim.BankIFSCCode != nil {
		response.BankIFSCCode = *claim.BankIFSCCode
	}

	if claim.BankAccountHolderName != nil {
		response.BankAccountHolderName = *claim.BankAccountHolderName
	}

	if claim.BankName != nil {
		response.BankName = *claim.BankName
	}

	if claim.WorkflowState != nil {
		response.DocumentCompleteness = *claim.WorkflowState
	}

	if claim.ApprovalDate != nil {
		response.ApprovalDate = claim.ApprovalDate.Format("2006-01-02 15:04:05")
	}

	if claim.DisbursementDate != nil {
		response.DisbursementDate = claim.DisbursementDate.Format("2006-01-02 15:04:05")
	}

	// Calculate estimated settlement date (7 days from submission as per BR-CLM-SB-001)
	if claim.SLADueDate.After(claim.CreatedAt) {
		response.EstimatedSettlementDate = claim.SLADueDate.Format("2006-01-02 15:04:05")
	}

	// Calculate days remaining
	if claim.SLADueDate.After(time.Now()) {
		days := int64(claim.SLADueDate.Sub(time.Now()).Hours() / 24)
		response.DaysRemaining = &days
	}

	return response
}

// NewSurvivalBenefitClaimsResponse creates a new SB claims list response
func NewSurvivalBenefitClaimsResponse(claims []domain.Claim, total int64, skip, limit int) SurvivalBenefitClaimsListResponse {
	claimResponses := make([]SurvivalBenefitClaimResponse, 0, len(claims))
	for _, claim := range claims {
		claimResponses = append(claimResponses, NewSurvivalBenefitClaimResponse(claim))
	}

	return SurvivalBenefitClaimsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse:     port.NewMetaDataResponse(uint64(total), uint64(skip), uint64(limit)),
		Claims:               claimResponses,
	}
}

// calculateSurvivalBenefitSLAStatus calculates SLA status for survival benefit claims
// Reference: BR-CLM-SB-001 (7 days SLA)
func calculateSurvivalBenefitSLAStatus(slaDueDate time.Time) string {
	now := time.Now()
	if slaDueDate.Before(now) {
		return "RED" // SLA breached
	}

	totalDuration := slaDueDate.Sub(now.Add(-7 * 24 * time.Hour)) // 7 days total SLA
	elapsed := now.Add(-7 * 24 * time.Hour).Sub(now)
	percentage := float64(elapsed) / float64(totalDuration) * 100

	if percentage < 70 {
		return "GREEN"
	} else if percentage < 90 {
		return "YELLOW"
	} else if percentage <= 100 {
		return "ORANGE"
	}
	return "RED"
}
