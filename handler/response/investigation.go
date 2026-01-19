package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ========================================
// INVESTIGATION RESPONSE DTOS
// ========================================

// InvestigationResponse represents an investigation in API responses
type InvestigationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationData `json:"data"`
}

// InvestigationData contains investigation details
type InvestigationData struct {
	InvestigationID           string     `json:"investigation_id"`
	ClaimID                   string     `json:"claim_id"`
	ClaimNumber               string     `json:"claim_number,omitempty"`
	PolicyID                  string     `json:"policy_id,omitempty"`
	AssignedBy                string     `json:"assigned_by"`
	InvestigatorID            string     `json:"investigator_id"`
	InvestigatorName          string     `json:"investigator_name,omitempty"`
	InvestigatorRank          string     `json:"investigator_rank,omitempty"`
	Jurisdiction              string     `json:"jurisdiction,omitempty"`
	AutoAssigned              bool       `json:"auto_assigned"`
	AssignmentDate            string     `json:"assignment_date"`
	DueDate                   string     `json:"due_date"`
	Status                    string     `json:"status"`
	ProgressPercentage        int        `json:"progress_percentage"`
	InvestigationOutcome      string     `json:"investigation_outcome,omitempty"`
	CauseOfDeath              string     `json:"cause_of_death,omitempty"`
	CauseOfDeathVerified      bool       `json:"cause_of_death_verified,omitempty"`
	HospitalRecordsVerified   bool       `json:"hospital_records_verified,omitempty"`
	DetailedFindings          string     `json:"detailed_findings,omitempty"`
	Recommendation            string     `json:"recommendation,omitempty"`
	ReportDocumentID          string     `json:"report_document_id,omitempty"`
	SubmittedAt               string     `json:"submitted_at,omitempty"`
	ReviewedBy                string     `json:"reviewed_by,omitempty"`
	ReviewedAt                string     `json:"reviewed_at,omitempty"`
	ReviewDecision            string     `json:"review_decision,omitempty"`
	ReviewerRemarks           string     `json:"reviewer_remarks,omitempty"`
	ReinvestigationCount      int        `json:"reinvestigation_count"`
	SLAStatus                 string     `json:"sla_status"` // GREEN, YELLOW, ORANGE, RED
	DaysRemaining             int        `json:"days_remaining"`
	CreatedAt                 string     `json:"created_at"`
	UpdatedAt                 string     `json:"updated_at"`
}

// InvestigationAssignedResponse represents response after assigning investigation officer
// POST /claims/death/{claim_id}/investigation/assign-officer
// Reference: BR-CLM-DC-002 (21-day SLA)
type InvestigationAssignedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationAssignmentData `json:"data"`
}

// InvestigationAssignmentData contains assignment details
type InvestigationAssignmentData struct {
	InvestigationID string `json:"investigation_id"`
	ClaimID         string `json:"claim_id"`
	InvestigatorID  string `json:"investigator_id"`
	AssignmentDate  string `json:"assignment_date"`
	DueDate         string `json:"due_date"`
	Status          string `json:"status"`
	Message         string `json:"message"`
}

// PendingInvestigationsResponse represents response for pending investigation claims list
// GET /claims/death/pending-investigation
type PendingInvestigationsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Data                      []PendingInvestigationClaim `json:"data"`
}

// PendingInvestigationClaim represents a claim pending investigation assignment
type PendingInvestigationClaim struct {
	ClaimID              string `json:"claim_id"`
	ClaimNumber          string `json:"claim_number"`
	PolicyID             string `json:"policy_id"`
	CustomerName         string `json:"customer_name,omitempty"`
	ClaimDate            string `json:"claim_date"`
	DeathDate            string `json:"death_date"`
	DeathType            string `json:"death_type"`
	InvestigationRequired bool  `json:"investigation_required"`
	Priority             string `json:"priority"` // LOW, MEDIUM, HIGH, URGENT
	Jurisdiction         string `json:"jurisdiction,omitempty"`
	SLAStatus            string `json:"sla_status"`
	DaysInQueue          int    `json:"days_in_queue"`
}

// InvestigationDetailsResponse represents investigation assignment details
// GET /claims/death/{claim_id}/investigation/{investigation_id}/details
type InvestigationDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationDetailData `json:"data"`
}

// InvestigationDetailData contains detailed investigation information
type InvestigationDetailData struct {
	InvestigationResponse
	ClaimDetails          ClaimSummary              `json:"claim_details"`
	InvestigationChecklist []string                 `json:"investigation_checklist"`
	ProgressTimeline      []InvestigationProgress   `json:"progress_timeline,omitempty"`
	EvidenceDocuments     []DocumentInfo            `json:"evidence_documents,omitempty"`
}

// ClaimSummary contains basic claim information
type ClaimSummary struct {
	ClaimID            string  `json:"claim_id"`
	ClaimNumber        string  `json:"claim_number"`
	PolicyID           string  `json:"policy_id"`
	CustomerID         string  `json:"customer_id"`
	CustomerName       string  `json:"customer_name"`
	ClaimType          string  `json:"claim_type"`
	ClaimDate          string  `json:"claim_date"`
	DeathDate          string  `json:"death_date,omitempty"`
	DeathType          string  `json:"death_type,omitempty"`
	ClaimAmount        *float64 `json:"claim_amount,omitempty"`
	InvestigationRequired bool  `json:"investigation_required"`
	Status             string  `json:"status"`
}

// InvestigationProgress represents a progress update
type InvestigationProgress struct {
	ID                   string    `json:"id"`
	InvestigationID      string    `json:"investigation_id"`
	ProgressNotes        string    `json:"progress_notes"`
	NextSteps            string    `json:"next_steps,omitempty"`
	Percentage           int       `json:"percentage"`
	ChecklistItemsCompleted []string `json:"checklist_items_completed,omitempty"`
	EstimatedCompletionDate string `json:"estimated_completion_date,omitempty"`
	RecordedAt           string    `json:"recorded_at"`
	RecordedBy           string    `json:"recorded_by,omitempty"`
}

// DocumentInfo represents document information
type DocumentInfo struct {
	DocumentID     string `json:"document_id"`
	DocumentType   string `json:"document_type"`
	DocumentName   string `json:"document_name"`
	UploadDate     string `json:"upload_date"`
	Verified       bool   `json:"verified"`
	VirusScanStatus string `json:"virus_scan_status,omitempty"`
}

// InvestigationProgressUpdateResponse represents response for progress update
// POST /claims/death/{claim_id}/investigation/{investigation_id}/progress-update
// Reference: BR-CLM-DC-002 (heartbeat tracking)
type InvestigationProgressUpdateResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationProgressUpdateData `json:"data"`
}

// InvestigationProgressUpdateData contains progress update confirmation
type InvestigationProgressUpdateData struct {
	ProgressID      string `json:"progress_id"`
	InvestigationID string `json:"investigation_id"`
	Percentage      int    `json:"percentage"`
	RecordedAt      string `json:"recorded_at"`
	Message         string `json:"message"`
}

// InvestigationReportSubmittedResponse represents response after submitting investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/submit-report
// Reference: BR-CLM-DC-011 (review within 5 days)
type InvestigationReportSubmittedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationReportData `json:"data"`
}

// InvestigationReportData contains submitted report details
type InvestigationReportData struct {
	InvestigationID      string `json:"investigation_id"`
	ClaimID              string `json:"claim_id"`
	ReportOutcome        string `json:"report_outcome"` // CLEAR, SUSPECT, FRAUD
	SubmittedAt          string `json:"submitted_at"`
	Status               string `json:"status"`
	ReviewDueDate        string `json:"review_due_date"` // 5-day SLA
	Message              string `json:"message"`
}

// InvestigationReviewResponse represents response for investigation report review
// POST /claims/death/{claim_id}/investigation/{investigation_id}/review
// Reference: BR-CLM-DC-011 (5-day review SLA)
type InvestigationReviewResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationReviewData `json:"data"`
}

// InvestigationReviewData contains review decision details
type InvestigationReviewData struct {
	InvestigationID string `json:"investigation_id"`
	ReviewDecision  string `json:"review_decision"` // ACCEPT, REINVESTIGATE, ESCALATE
	ReviewedAt      string `json:"reviewed_at"`
	ReviewedBy      string `json:"reviewed_by"`
	Status          string `json:"status"`
	Message         string `json:"message"`
	NextAction      string `json:"next_action,omitempty"`
}

// ReinvestigationTriggeredResponse represents response for reinvestigation trigger
// POST /claims/death/{id}/investigation/trigger-reinvestigation
// Reference: BR-CLM-DC-013 (max 2 times, 14-day SLA)
type ReinvestigationTriggeredResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 ReinvestigationData `json:"data"`
}

// ReinvestigationData contains reinvestigation details
type ReinvestigationData struct {
	NewInvestigationID   string `json:"new_investigation_id"`
	OriginalInvestigationID string `json:"original_investigation_id"`
	ReinvestigationCount int    `json:"reinvestigation_count"`
	NewDueDate          string `json:"new_due_date"`
	Status              string `json:"status"`
	Message             string `json:"message"`
}

// InvestigationSLAEscalationResponse represents response for SLA breach escalation
// POST /claims/death/{id}/investigation/escalate-sla-breach
// Reference: BR-CLM-DC-002 (escalation hierarchy)
type InvestigationSLAEscalationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 InvestigationSLAEscalationData `json:"data"`
}

// InvestigationSLAEscalationData contains escalation details
type InvestigationSLAEscalationData struct {
	InvestigationID string `json:"investigation_id"`
	EscalationLevel string `json:"escalation_level"`
	EscalatedTo     string `json:"escalated_to"`
	EscalatedAt     string `json:"escalated_at"`
	Status          string `json:"status"`
	Message         string `json:"message"`
}

// ManualReviewAssignedResponse represents response for manual review assignment
// POST /claims/death/{id}/manual-review/assign
// Reference: BR-CLM-DC-011 (SUSPECT outcome handling)
type ManualReviewAssignedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 ManualReviewAssignmentData `json:"data"`
}

// ManualReviewAssignmentData contains manual review assignment details
type ManualReviewAssignmentData struct {
	ClaimID     string `json:"claim_id"`
	ReviewerID  string `json:"reviewer_id"`
	Priority    string `json:"priority"`
	AssignedAt  string `json:"assigned_at"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

// ClaimRejectedForFraudResponse represents response for rejecting claim based on fraud
// POST /claims/death/{id}/reject-fraud
// Reference: BR-CLM-DC-020 (fraud rejection)
type ClaimRejectedForFraudResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                 FraudRejectionData `json:"data"`
}

// FraudRejectionData contains fraud rejection details
type FraudRejectionData struct {
	ClaimID              string `json:"claim_id"`
	RejectionDate        string `json:"rejection_date"`
	InvestigationReportID string `json:"investigation_report_id"`
	LegalActionRequired  bool   `json:"legal_action_required"`
	LegalCaseID          string `json:"legal_case_id,omitempty"`
	Status               string `json:"status"`
	Message              string `json:"message"`
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// NewInvestigationResponse creates an investigation response from domain model
func NewInvestigationResponse(inv domain.Investigation, slaStatus string, daysRemaining int) InvestigationResponse {
	return InvestigationResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: InvestigationData{
			InvestigationID:           inv.InvestigationID,
			ClaimID:                   inv.ClaimID,
			AssignedBy:                inv.AssignedBy,
			InvestigatorID:            inv.InvestigatorID,
			InvestigatorRank:          getStringValue(inv.InvestigatorRank),
			Jurisdiction:              getStringValue(inv.Jurisdiction),
			AutoAssigned:              inv.AutoAssigned,
			AssignmentDate:            formatTime(inv.AssignmentDate),
			DueDate:                   formatTime(inv.DueDate),
			Status:                    inv.Status,
			ProgressPercentage:        int(inv.ProgressPercentage),
			InvestigationOutcome:      getStringValue(inv.InvestigationOutcome),
			CauseOfDeath:              getStringValue(inv.CauseOfDeath),
			CauseOfDeathVerified:      inv.CauseOfDeathVerified,
			HospitalRecordsVerified:   inv.HospitalRecordsVerified,
			DetailedFindings:          getStringValue(inv.DetailedFindings),
			Recommendation:            getStringValue(inv.Recommendation),
			ReportDocumentID:          getStringValue(inv.ReportDocumentID),
			SubmittedAt:               formatTimePtr(inv.SubmittedAt),
			ReviewedBy:                getStringValue(inv.ReviewedBy),
			ReviewedAt:                formatTimePtr(inv.ReviewedAt),
			ReviewDecision:            getStringValue(inv.ReviewDecision),
			ReviewerRemarks:           getStringValue(inv.ReviewerRemarks),
			ReinvestigationCount:      int(inv.ReinvestigationCount),
			SLAStatus:                 slaStatus,
			DaysRemaining:             daysRemaining,
			CreatedAt:                 formatTime(inv.CreatedAt),
			UpdatedAt:                 formatTime(inv.UpdatedAt),
		},
	}
}

// NewPendingInvestigationsResponse creates a response for pending investigation claims list
func NewPendingInvestigationsResponse(claims []domain.Claim, total int64, skip, limit int) PendingInvestigationsResponse {
	data := make([]PendingInvestigationClaim, len(claims))
	for i, claim := range claims {
		slaStatus := calculateSLAStatus(claim.SLAStatus, claim.SLADueDate)
		daysInQueue := int(time.Since(claim.CreatedAt).Hours() / 24)

		data[i] = PendingInvestigationClaim{
			ClaimID:              claim.ID,
			ClaimNumber:          claim.ClaimNumber,
			PolicyID:             claim.PolicyID,
			ClaimDate:            formatTime(claim.ClaimDate),
			DeathDate:            formatTimePtr(claim.DeathDate),
			DeathType:            getStringValue(claim.DeathType),
			InvestigationRequired: claim.InvestigationRequired,
			Priority:             "MEDIUM", // TODO: Calculate based on claim age and type
			Jurisdiction:         "",      // TODO: Get from claim metadata
			SLAStatus:            slaStatus,
			DaysInQueue:          daysInQueue,
		}
	}

	return PendingInvestigationsResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		MetaDataResponse:     port.NewMetaDataResponse(uint64(skip), uint64(limit), uint64(total)),
		Data:                 data,
	}
}

// Helper function to safely get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to format time as string
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// Helper function to format time pointer as string
func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
