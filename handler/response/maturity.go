package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/port"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// ========================================
// MATURITY CLAIM RESPONSE DTOS
// ========================================

// MaturityIntimationBatchResponse represents the response for batch intimation
// POST /claims/maturity/send-intimation-batch
type MaturityIntimationBatchResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	TotalPolicies             int `json:"total_policies"`
	IntimationsSent           int `json:"intimations_sent"`
	Failed                    int `json:"failed"`
}

// MaturityDueReportResponse represents the response for maturity due report
// POST /claims/maturity/generate-due-report
type MaturityDueReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ReportID                  string  `json:"report_id"`
	ReportURL                 string  `json:"report_url"`
	TotalPolicies             int     `json:"total_policies"`
	TotalAmount               float64 `json:"total_amount"`
}

// BankDetails represents bank account details
type BankDetails struct {
	BankName           string `json:"bank_name"`
	BankAccountNumber  string `json:"bank_account_number"`
	BankIFSC           string `json:"bank_ifsc"`
	BankAccountType    string `json:"bank_account_type"`
	AccountHolderName  string `json:"account_holder_name"`
}

// MaturityPreFillDataResponse represents the response for pre-filled data
// GET /claims/maturity/pre-fill-data
type MaturityPreFillDataResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string       `json:"policy_id"`
	CustomerName              string       `json:"customer_name"`
	MaturityDate              string       `json:"maturity_date"`
	MaturityAmount            float64      `json:"maturity_amount"`
	RegisteredMobile          string       `json:"registered_mobile"`
	RegisteredEmail           string       `json:"registered_email"`
	BankDetailsOnRecord       *BankDetails `json:"bank_details_on_record,omitempty"`
}

// MaturityClaimRegistrationResponse represents the response for maturity claim submission
// POST /claims/maturity/submit
// Reference: BR-CLM-MC-001 (7 days SLA)
type MaturityClaimRegistrationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string    `json:"claim_id"`
	ClaimNumber               string    `json:"claim_number"`
	PolicyID                  string    `json:"policy_id"`
	MaturityDate              string    `json:"maturity_date"`
	ClaimAmount               float64   `json:"claim_amount"`
	SLADueDate                string    `json:"sla_due_date"`
	SubmissionDate            string    `json:"submission_date"`
}

// MaturityClaimResponse represents a maturity claim details
type MaturityClaimResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string     `json:"claim_id"`
	ClaimNumber               string     `json:"claim_number"`
	PolicyID                  string     `json:"policy_id"`
	ClaimType                 string     `json:"claim_type"`
	ClaimantName              string     `json:"claimant_name"`
	ClaimantRelationship      string     `json:"claimant_relationship"`
	ClaimantMobile            string     `json:"claimant_mobile"`
	ClaimantEmail             string     `json:"claimant_email"`
	MaturityDate              string     `json:"maturity_date"`
	ClaimAmount               float64    `json:"claim_amount"`
	DisbursementMode          string     `json:"disbursement_mode"`
	BankAccountNumber         string     `json:"bank_account_number"`
	BankIFSC                  string     `json:"bank_ifsc"`
	BankAccountType           string     `json:"bank_account_type"`
	Status                    string     `json:"status"`
	SLAStatus                 string     `json:"sla_status"`
	SLADueDate                string     `json:"sla_due_date"`
	CreatedAt                 string     `json:"created_at"`
	UpdatedAt                 string     `json:"updated_at"`
}

// MaturityClaimsListResponse represents a list of maturity claims
type MaturityClaimsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:"metadata,omitempty"`
	Claims                    []MaturityClaimResponse `json:"claims"`
}

// DocumentsValidatedResponse represents the response for document validation
// POST /claims/maturity/{claim_id}/validate-documents
type DocumentsValidatedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Valid                    bool     `json:"valid"`
	MandatoryDocuments       int      `json:"mandatory_documents"`
	MandatoryVerified        int      `json:"mandatory_verified"`
	OptionalDocuments        int      `json:"optional_documents"`
	OptionalVerified         int      `json:"optional_verified"`
	MissingDocuments         []string `json:"missing_documents,omitempty"`
}

// OCRDataExtractedResponse represents the response for OCR data extraction
// POST /claims/maturity/{claim_id}/extract-ocr-data
type OCRDataExtractedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ExtractedData            map[string]interface{} `json:"extracted_data"`
	ConfidenceScore          float64                `json:"confidence_score"`
	FieldsExtracted          []string               `json:"fields_extracted"`
}

// QCVerificationResponse represents the response for QC verification
// POST /claims/maturity/{claim_id}/qc-verify
type QCVerificationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	QCStatus                 string  `json:"qc_status"`
	QCRemarks                *string `json:"qc_remarks,omitempty"`
	QCVerifiedBy             string  `json:"qc_verified_by"`
	QCVerificationDate       string  `json:"qc_verification_date"`
}

// MaturityApprovalDetailsResponse represents the response for approval details
// GET /claims/maturity/{claim_id}/approval-details
type MaturityApprovalDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                  string                 `json:"claim_id"`
	ClaimNumber              string                 `json:"claim_number"`
	PolicyID                 string                 `json:"policy_id"`
	MaturityDate             string                 `json:"maturity_date"`
	ClaimAmount              float64                `json:"claim_amount"`
	ClaimantDetails          MaturityClaimantData   `json:"claimant_details"`
	DisbursementDetails      MaturityDisbursementData `json:"disbursement_details"`
	DocumentDetails          MaturityDocumentData   `json:"document_details"`
	CalculationDetails       MaturityCalculationData `json:"calculation_details"`
	ApprovalHistory          []ApprovalHistoryItem  `json:"approval_history"`
	CurrentApprovalLevel     string                 `json:"current_approval_level"`
	EligibleApprovers        []ApproverInfo         `json:"eligible_approvers"`
}

// MaturityClaimantData represents claimant information
type MaturityClaimantData struct {
	Name           string `json:"name"`
	Relationship   string `json:"relationship"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	IsNRI          bool   `json:"is_nri"`
	NRICountry     *string `json:"nri_country,omitempty"`
	IsPANAvailable bool   `json:"is_pan_available"`
	PANNumber      *string `json:"pan_number,omitempty"`
}

// MaturityDisbursementData represents disbursement information
type MaturityDisbursementData struct {
	DisbursementMode    string  `json:"disbursement_mode"`
	BankAccountNumber   string  `json:"bank_account_number"`
	BankIFSC            string  `json:"bank_ifsc"`
	BankAccountType     string  `json:"bank_account_type"`
	AccountHolderName   string  `json:"account_holder_name"`
	DisbursementAmount  float64 `json:"disbursement_amount"`
}

// MaturityDocumentData represents document information
type MaturityDocumentData struct {
	TotalDocuments       int                  `json:"total_documents"`
	VerifiedDocuments    int                  `json:"verified_documents"`
	PendingDocuments     int                  `json:"pending_documents"`
	DocumentChecklist    []DocumentChecklistItem `json:"document_checklist"`
}

// MaturityCalculationData represents calculation information
type MaturityCalculationData struct {
	SumAssured          float64 `json:"sum_assured"`
	Bonuses             float64 `json:"bonuses"`
	AccruedBonus        float64 `json:"accrued_bonus"`
	TotalAmount         float64 `json:"total_amount"`
	LoanOutstanding     float64 `json:"loan_outstanding"`
	PremiumDue          float64 `json:"premium_due"`
	NetPayableAmount    float64 `json:"net_payable_amount"`
}

// ApprovalHistoryItem represents an approval history item
type ApprovalHistoryItem struct {
	ApprovalLevel    string `json:"approval_level"`
	ApproverID       string `json:"approver_id"`
	ApproverName     string `json:"approver_name"`
	ApprovalStatus   string `json:"approval_status"`
	ApprovalAmount   float64 `json:"approval_amount"`
	ApprovalRemarks  string `json:"approval_remarks"`
	ApprovalDate     string `json:"approval_date"`
}

// MaturityClaimApprovedResponse represents the response for claim approval
// POST /claims/maturity/{claim_id}/approve
type MaturityClaimApprovedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                  string  `json:"claim_id"`
	ApprovalStatus           string  `json:"approval_status"`
	ApprovalAmount           float64 `json:"approval_amount"`
	ApprovalRemarks          string  `json:"approval_remarks"`
	ApprovedBy               string  `json:"approved_by"`
	ApprovalDate             string  `json:"approval_date"`
	ApprovalLevel            string  `json:"approval_level"`
}

// MaturityClaimDisbursementInitiatedResponse represents the response for disbursement initiation
// POST /claims/maturity/{claim_id}/disburse
type MaturityClaimDisbursementInitiatedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                  string  `json:"claim_id"`
	DisbursementID           string  `json:"disbursement_id"`
	DisbursementAmount       float64 `json:"disbursement_amount"`
	DisbursementMode         string  `json:"disbursement_mode"`
	ReferenceNumber         string  `json:"reference_number"`
	EstimatedTransferDate    string  `json:"estimated_transfer_date"`
}

// MaturityVoucherGeneratedResponse represents the response for voucher generation
// POST /claims/maturity/{claim_id}/generate-voucher
type MaturityVoucherGeneratedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	VoucherNumber            string    `json:"voucher_number"`
	VoucherDate              string    `json:"voucher_date"`
	VoucherURL               string    `json:"voucher_url"`
	ClaimID                  string    `json:"claim_id"`
	DisbursementAmount       float64   `json:"disbursement_amount"`
}

// MaturityDocumentReminderSentResponse represents the response for document reminder
// POST /claims/maturity/{claim_id}/send-document-reminder
type MaturityDocumentReminderSentResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ReminderSent              bool   `json:"reminder_sent"`
	ReminderDate              string `json:"reminder_date"`
	Channel                   string `json:"channel"`
}

// MaturityClaimClosedResponse represents the response for claim closure
// POST /claims/maturity/{claim_id}/close
type MaturityClaimClosedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string `json:"claim_id"`
	ClosureDate               string `json:"closure_date"`
	ClosureReason             string `json:"closure_reason"`
	ClosedBy                  string `json:"closed_by"`
}

// MaturityFeedbackRequestedResponse represents the response for feedback request
// POST /claims/maturity/{claim_id}/request-feedback
type MaturityFeedbackRequestedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	FeedbackRequestSent       bool   `json:"feedback_request_sent"`
	FeedbackURL               string `json:"feedback_url"`
	Channel                   string `json:"channel"`
	SentDate                  string `json:"sent_date"`
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// NewMaturityClaimResponse creates a new MaturityClaimResponse from domain.Claim
func NewMaturityClaimResponse(claim domain.Claim) MaturityClaimResponse {
	response := MaturityClaimResponse{
		ClaimID:           claim.ID,
		ClaimNumber:       claim.ClaimNumber,
		PolicyID:          claim.PolicyID,
		ClaimType:         claim.ClaimType,
		ClaimantName:      claim.ClaimantName,
		Status:            claim.Status,
		CreatedAt:         claim.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         claim.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Handle optional fields safely
	if claim.ClaimantPhone != nil {
		response.ClaimantMobile = *claim.ClaimantPhone
	}
	if claim.ClaimantEmail != nil {
		response.ClaimantEmail = *claim.ClaimantEmail
	}
	if claim.ClaimAmount != nil {
		response.ClaimAmount = *claim.ClaimAmount
	}
	if claim.PaymentMode != nil {
		response.DisbursementMode = *claim.PaymentMode
	}
	if claim.BankAccountNumber != nil {
		response.BankAccountNumber = *claim.BankAccountNumber
	}
	if claim.BankIFSCCode != nil {
		response.BankIFSC = *claim.BankIFSCCode
	}
	if claim.SLADueDate != (time.Time{}) {
		response.SLADueDate = claim.SLADueDate.Format("2006-01-02 15:04:05")
	}

	// Calculate SLA status
	if claim.SLADueDate != (time.Time{}) {
		response.SLAStatus = calculateMaturitySLAStatus(claim.SLADueDate, claim.Status)
	}

	return response
}

// NewMaturityClaimsResponse creates a new MaturityClaimsListResponse
func NewMaturityClaimsResponse(claims []domain.Claim, total int64, skip, limit int64) MaturityClaimsListResponse {
	claimResponses := make([]MaturityClaimResponse, 0, len(claims))
	for _, claim := range claims {
		claimResponses = append(claimResponses, NewMaturityClaimResponse(claim))
	}

	return MaturityClaimsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse: port.NewMetaDataResponse(uint64(skip), uint64(limit), uint64(total)),
		Claims:           claimResponses,
	}
}

// calculateMaturitySLAStatus calculates SLA status for maturity claims
// Reference: BR-CLM-MC-001 (7 days SLA)
func calculateMaturitySLAStatus(slaDueDate time.Time, status string) string {
	now := time.Now()

	// If claim is already closed/approved/disbursed, SLA is not applicable
	if status == "CLOSED" || status == "APPROVED" || status == "DISBURSED" {
		return "COMPLETED"
	}

	// If SLA due date has passed
	if now.After(slaDueDate) {
		return "RED"
	}

	// Calculate percentage of SLA remaining
	totalDuration := slaDueDate.Sub(now)
	remainingPercentage := (totalDuration.Seconds() / (7 * 24 * time.Hour).Seconds()) * 100

	if remainingPercentage >= 70 {
		return "GREEN"
	} else if remainingPercentage >= 30 {
		return "YELLOW"
	} else {
		return "ORANGE"
	}
}
