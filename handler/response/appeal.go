package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// AppealEligibilityData represents the data for appeal eligibility check
// Reference: FR-CLM-DC-023, BR-CLM-DC-005 (90-day appeal window)
type AppealEligibilityData struct {
	ClaimID                string    `json:"claim_id"`
	IsEligible             bool      `json:"is_eligible"`
	AppealWindowEnds       string    `json:"appeal_window_ends"`       // Format: "2006-01-02 15:04:05"
	DaysRemaining          int       `json:"days_remaining"`
	RejectionReason        *string   `json:"rejection_reason,omitempty"`        // Reason for claim rejection
	CondonationRequired    bool      `json:"condonation_required"`              // True if appeal window expired
	CondonationReason      *string   `json:"condonation_reason,omitempty"`       // Reason if condonation needed
	AppealTypesAvailable   []string  `json:"appeal_types_available"`            // RECONSIDERATION, APPELLATE_AUTHORITY, OMBUDSMAN
}

// AppealEligibilityResponse represents the response for appeal eligibility check
// GET /claims/death/{claim_id}/appeal-eligibility
type AppealEligibilityResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      AppealEligibilityData `json:"data"`
}

// AppellateAuthorityData represents the data for appellate authority
// Reference: BR-CLM-DC-005 (Escalation to next higher authority)
type AppellateAuthorityData struct {
	AuthorityID      string  `json:"authority_id"`
	AuthorityName    string  `json:"authority_name"`
	AuthorityLevel   string  `json:"authority_level"`    // LEVEL_1, LEVEL_2, LEVEL_3, LEVEL_4
	AuthorityType    string  `json:"authority_type"`    // DIVISION_HEAD, ZONAL_MANAGER, etc.
	Designation      string  `json:"designation"`
	Department       string  `json:"department"`
	Contact          string  `json:"contact,omitempty"`
	Email            string  `json:"email,omitempty"`
	Location         string  `json:"location,omitempty"`
	Jurisdiction     string  `json:"jurisdiction"`       // Geographic or functional jurisdiction
	MaxClaimAmount   float64 `json:"max_claim_amount"`   // Maximum claim amount they can approve
}

// AppellateAuthorityResponse represents the response for appellate authority details
// GET /claims/death/{claim_id}/appellate-authority
type AppellateAuthorityResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      AppellateAuthorityData `json:"data"`
}

// AppealSubmittedData represents the data for appeal submission
// Reference: BR-CLM-DC-005 (90-day window), BR-CLM-DC-007 (45-day SLA)
type AppealSubmittedData struct {
	AppealID              string    `json:"appeal_id"`
	AppealNumber          string    `json:"appeal_number"`
	ClaimID               string    `json:"claim_id"`
	SubmissionDate        string    `json:"submission_date"`         // Format: "2006-01-02 15:04:05"
	AppealDeadline        string    `json:"appeal_deadline"`         // Decision deadline (45 days)
	AppellateAuthority    string    `json:"appellate_authority"`     // Name of assigned authority
	CurrentStatus         string    `json:"current_status"`          // SUBMITTED, UNDER_REVIEW, PENDING_DECISION
	ExpectedDecisionBy    string    `json:"expected_decision_by"`    // Format: "2006-01-02 15:04:05"
	AppealSLAStatus       string    `json:"appeal_sla_status"`       // GREEN, YELLOW, ORANGE, RED
	TrackingURL           *string   `json:"tracking_url,omitempty"`  // URL to track appeal status
	NextSteps             []string  `json:"next_steps"`              // What happens next
	DocumentsRequired     []string  `json:"documents_required"`      // List of required documents (if any)
}

// AppealSubmissionResponse represents the response for appeal submission
// POST /claims/death/{claim_id}/appeal
type AppealSubmissionResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      AppealSubmittedData `json:"data"`
}

// AppealDetailsData represents the full details of an appeal
type AppealDetailsData struct {
	AppealID              string     `json:"appeal_id"`
	AppealNumber          string     `json:"appeal_number"`
	ClaimID               string     `json:"claim_id"`
	ClaimNumber           string     `json:"claim_number"`
	AppellantName         string     `json:"appellant_name"`
	AppellantContact      string     `json:"appellant_contact"`
	GroundsOfAppeal       string     `json:"grounds_of_appeal"`
	SupportingDocuments   *string    `json:"supporting_documents,omitempty"`
	SubmissionDate        string     `json:"submission_date"`         // Format: "2006-01-02 15:04:05"
	AppealDeadline        string     `json:"appeal_deadline"`         // Decision deadline
	AppealType            string     `json:"appeal_type"`             // RECONSIDERATION, APPELLATE_AUTHORITY, OMBUDSMAN
	AppellateAuthority    string     `json:"appellate_authority"`     // Name of assigned authority
	CurrentStatus         string     `json:"current_status"`          // SUBMITTED, UNDER_REVIEW, DECIDED, etc.
	Decision              *string    `json:"decision,omitempty"`      // APPEAL_ACCEPTED, APPEAL_REJECTED, PARTIAL_ACCEPTANCE
	ReasonedOrder         *string    `json:"reasoned_order,omitempty"` // Detailed decision with reasons
	RevisedClaimAmount    *float64   `json:"revised_claim_amount,omitempty"` // If appeal accepted with different amount
	DecisionDate          *string    `json:"decision_date,omitempty"` // Format: "2006-01-02 15:04:05"
	AppealSLAStatus       string     `json:"appeal_sla_status"`       // GREEN, YELLOW, ORANGE, RED
	DaysUntilDeadline     int        `json:"days_until_deadline"`
	CondonationRequested  bool       `json:"condonation_requested"`
	CondonationReason     *string    `json:"condonation_reason,omitempty"`
	CondonationApproved   *bool      `json:"condonation_approved,omitempty"`
	AppealHistory         []AppealHistoryItem `json:"appeal_history,omitempty"`
}

// AppealHistoryItem represents a single event in appeal history
type AppealHistoryItem struct {
	EventType        string    `json:"event_type"`        // SUBMITTED, UNDER_REVIEW, DECIDED, etc.
	EventDate        string    `json:"event_date"`        // Format: "2006-01-02 15:04:05"
	PerformedBy      string    `json:"performed_by"`      // User or system
	Remarks          *string   `json:"remarks,omitempty"`
}

// AppealDetailsResponse represents the response for appeal details
// GET /claims/death/{claim_id}/appeal/{appeal_id}
type AppealDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      AppealDetailsData `json:"data"`
}

// AppealDecisionData represents the data for appeal decision recording
// Reference: BR-CLM-DC-006 (45-day SLA)
type AppealDecisionData struct {
	AppealID           string  `json:"appeal_id"`
	AppealNumber       string  `json:"appeal_number"`
	Decision           string  `json:"decision"`                      // APPEAL_ACCEPTED, APPEAL_REJECTED, PARTIAL_ACCEPTANCE
	ReasonedOrder      string  `json:"reasoned_order"`                // Detailed decision with reasons
	DecisionDate       string  `json:"decision_date"`                 // Format: "2006-01-02 15:04:05"
	DecisionBy         string  `json:"decision_by"`                   // Authority who made decision
	RevisedClaimAmount *float64 `json:"revised_claim_amount,omitempty"` // If accepted with different amount
	ModificationDetails *string `json:"modification_details,omitempty"` // What was changed
	ClaimStatusUpdate  string  `json:"claim_status_update"`           // New claim status after decision
	NextSteps          []string `json:"next_steps"`                   // What happens next
}

// AppealDecisionResponse represents the response for appeal decision recording
// POST /claims/death/{claim_id}/appeal/{appeal_id}/decision
type AppealDecisionResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      AppealDecisionData `json:"data"`
}

// AppealResponse represents the response for a single appeal
type AppealResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      domain.Appeal `json:"data"`
}

// AppealsListResponse represents the response for list of appeals
type AppealsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      []domain.Appeal `json:"data"`
	port.MetaDataResponse     `json:",inline"`
}

// Helper functions

// NewAppealResponse creates a new appeal response from domain model
func NewAppealResponse(appeal domain.Appeal) *AppealResponse {
	return &AppealResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Appeal details retrieved successfully",
		},
		Data: appeal,
	}
}

// NewAppealsListResponse creates a new appeals list response
func NewAppealsListResponse(appeals []domain.Appeal, total int64, skip, limit int) *AppealsListResponse {
	return &AppealsListResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Appeals retrieved successfully",
		},
		Data: appeals,
		MetaDataResponse: port.MetaDataResponse{
			Skip:                 uint64(skip),
			Limit:                uint64(limit),
			TotalRecordsCount:    int(total),
			ReturnedRecordsCount: uint64(len(appeals)),
		},
	}
}

// CalculateAppealSLAStatus calculates the SLA status for an appeal
// Reference: BR-CLM-DC-007 (45-day decision timeline)
func CalculateAppealSLAStatus(submissionDate, deadline string) string {
	submission, err := time.Parse("2006-01-02T15:04:05Z", submissionDate)
	if err != nil {
		submission, err = time.Parse("2006-01-02 15:04:05", submissionDate)
		if err != nil {
			return "UNKNOWN"
		}
	}

	deadlineTime, err := time.Parse("2006-01-02T15:04:05Z", deadline)
	if err != nil {
		deadlineTime, err = time.Parse("2006-01-02 15:04:05", deadline)
		if err != nil {
			return "UNKNOWN"
		}
	}

	now := time.Now()
	totalDuration := deadlineTime.Sub(submission)
	elapsed := now.Sub(submission)

	if totalDuration <= 0 {
		return "RED"
	}

	percentageUsed := float64(elapsed) / float64(totalDuration) * 100

	// SLA Status Calculation
	// GREEN: < 70% of time used
	// YELLOW: 70% - 90% of time used
	// ORANGE: 90% - 100% of time used
	// RED: > 100% of time used (breached)
	switch {
	case percentageUsed >= 100:
		return "RED"
	case percentageUsed >= 90:
		return "ORANGE"
	case percentageUsed >= 70:
		return "YELLOW"
	default:
		return "GREEN"
	}
}
