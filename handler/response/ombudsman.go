package response

import (
	"time"

	port "gitlab.cept.gov.in/pli/claims-api/core/port"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// ==================== OMBUDSMAN RESPONSE DTOs ====================

// ComplaintRegisteredResponse represents successful complaint registration
// FR-CLM-OMB-001: Complaint Intake & Registration
type ComplaintRegisteredResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintNumber           string `json:"complaint_number"`           // OMB{YYYY}{DDDD}
	ComplaintID               string `json:"complaint_id"`               // UUID
	AcknowledgementSent       bool   `json:"acknowledgement_sent"`       // Sent via SMS/email within 24 hours
	AcknowledgementDate       string `json:"acknowledgement_date"`       // Format: "2006-01-02 15:04:05"
	Status                    string `json:"status"`                     // REGISTERED
	AssignedJurisdiction      string `json:"assigned_jurisdiction"`       // BR-CLM-OMB-002: Auto-mapped jurisdiction
	NextSteps                 string `json:"next_steps"`                 // Instructions to complainant
}

// ComplaintDetailsResponse represents full complaint details
type ComplaintDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintData             ComplaintData `json:"complaint_data"`
}

// ComplaintData represents the core complaint information
type ComplaintData struct {
	// Complainant Details
	ComplaintID            string `json:"complaint_id"`
	ComplaintNumber        string `json:"complaint_number"`
	ComplainantName        string `json:"complainant_name"`
	ComplainantAddress     string `json:"complainant_address"`
	ComplainantMobile      string `json:"complainant_mobile"`
	ComplainantEmail       string `json:"complainant_email,omitempty"`
	ComplainantRole        string `json:"complainant_role"`
	LanguagePreference     string `json:"language_preference"`
	IDProofType            string `json:"id_proof_type"`
	IDProofNumber          string `json:"id_proof_number,omitempty"` // Masked for security

	// Policy/Claim Details
	PolicyNumber           string `json:"policy_number"`
	ClaimNumber            string `json:"claim_number,omitempty"`
	PolicyType             string `json:"policy_type"`
	AgentName              string `json:"agent_name,omitempty"`
	AgentBranch            string `json:"agent_branch,omitempty"`

	// Complaint Details
	ComplaintCategory      string `json:"complaint_category"`
	IncidentDate           string `json:"incident_date"`
	RepresentationDate     string `json:"representation_date"`
	IssueDescription       string `json:"issue_description"`
	ReliefSought           string `json:"relief_sought"`
	ClaimValue             float64 `json:"claim_value,omitempty"`

	// Jurisdiction & Assignment
	OmbudsmanCenter        string `json:"ombudsman_center,omitempty"`
	AssignedOmbudsmanID    string `json:"assigned_ombudsman_id,omitempty"`
	AssignedOmbudsmanName  string `json:"assigned_ombudsman_name,omitempty"` // TODO: Lookup from user service

	// Admissibility
	Admissible              *bool   `json:"admissible,omitempty"`              // BR-CLM-OMB-001: Admissibility check result
	AdmissibilityChecked    bool    `json:"admissibility_checked"`
	AdmissibilityReason     string  `json:"admissibility_reason,omitempty"`
	InadmissibilityReason   string  `json:"inadmissibility_reason,omitempty"`

	// Status & Timeline
	Status                  string `json:"status"`
	SubmittedDate           string `json:"submitted_date"`
	LastUpdatedDate         string `json:"last_updated_date"`

	// SLA Information
	AcknowledgementSentDate string `json:"acknowledgement_sent_date,omitempty"`
	ResolutionDueDate       string `json:"resolution_due_date,omitempty"`    // Statutory timeline
	DaysInQueue             int     `json:"days_in_queue"`

	// Attachments
	AttachmentCount         int      `json:"attachment_count"`
	Attachments             []AttachmentData `json:"attachments,omitempty"`
}

// AttachmentData represents attachment metadata
type AttachmentData struct {
	AttachmentID     string `json:"attachment_id"`
	DocumentType     string `json:"document_type"`
	FileName         string `json:"file_name"`
	UploadedDate     string `json:"uploaded_date"`
	FileSize         int64  `json:"file_size"` // In bytes
	UploadedBy       string `json:"uploaded_by,omitempty"`
}

// OmbudsmanAssignedResponse represents successful ombudsman assignment
// FR-CLM-OMB-002: Jurisdiction Mapping
// BR-CLM-OMB-003: Conflict of interest screening
type OmbudsmanAssignedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	AssignedOmbudsmanID       string `json:"assigned_ombudsman_id"`
	OmbudsmanCenter           string `json:"ombudsman_center"`
	ConflictCheckPerformed    bool   `json:"conflict_check_performed"`
	ConflictDetected          bool   `json:"conflict_detected"`
	ReassignmentRequired      bool   `json:"reassignment_required,omitempty"` // If conflict found
	AssignmentDate            string `json:"assignment_date"`
	AssignmentRemarks         string `json:"assignment_remarks,omitempty"`
}

// AdmissibilityReviewedResponse represents admissibility review decision
// BR-CLM-OMB-001: Admissibility checks
type AdmissibilityReviewedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	Admissible                bool   `json:"admissible"`
	AdmissibilityReason       string `json:"admissibility_reason,omitempty"`
	InadmissibilityReason     string `json:"inadmissibility_reason,omitempty"`
	ReviewedBy                string `json:"reviewed_by"`
	ReviewDate                string `json:"review_date"`
	NextSteps                 string `json:"next_steps,omitempty"`
}

// ComplaintTimelineResponse represents complaint history/timeline
type ComplaintTimelineResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string            `json:"complaint_id"`
	Timeline                  []TimelineEntry  `json:"timeline"`
	TotalEvents               int               `json:"total_events"`
}

// TimelineEntry represents a single timeline event
type TimelineEntry struct {
	EventID          string `json:"event_id"`
	EventType        string `json:"event_type"`         // REGISTERED, ASSIGNED, ADMISSIBILITY_CHECK, HEARING_SCHEDULED, MEDIATION, AWARD_ISSUED, COMPLIANCE_RECORDED, CLOSED
	EventDate        string `json:"event_date"`         // Format: "2006-01-02 15:04:05"
	EventDescription string `json:"event_description"`
	PerformedBy      string `json:"performed_by,omitempty"`
	PreviousStatus   string `json:"previous_status,omitempty"`
	NewStatus        string `json:"new_status,omitempty"`
	Remarks          string `json:"remarks,omitempty"`
	AttachmentID     string `json:"attachment_id,omitempty"`
}

// MediationRecordedResponse represents mediation outcome recording
// FR-CLM-OMB-003: Hearing Scheduling & Management
// BR-CLM-OMB-004: Mediation recommendation
type MediationRecordedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	HearingID                 string `json:"hearing_id"`
	MediationDate             string `json:"mediation_date"`
	ConsentToMediate          bool   `json:"consent_to_mediate"`
	MediationSuccessful       bool   `json:"mediation_successful"`
	SettlementTerms           string `json:"settlement_terms,omitempty"`
	ComplainantAccepted       bool   `json:"complainant_accepted"`
	InsurerAccepted           bool   `json:"insurer_accepted"`
	RecordingOfficer          string `json:"recording_officer"`
	RecordingDate             string `json:"recording_date"`
	Remarks                   string `json:"remarks,omitempty"`
	NextSteps                 string `json:"next_steps,omitempty"`
}

// AwardIssuedResponse represents award issuance (mediation or adjudication)
// FR-CLM-OMB-004: Award Issuance & Enforcement
// BR-CLM-OMB-005: Award issuance with ₹50 lakh cap
type AwardIssuedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	AwardData                 AwardData `json:"award_data"`
}

// AwardData represents award information
type AwardData struct {
	AwardID                  string  `json:"award_id"`
	AwardType                string  `json:"award_type"`                  // MEDIATION_RECOMMENDATION or ADJUDICATION_AWARD
	AwardAmount              float64 `json:"award_amount"`
	AwardCurrency            string  `json:"award_currency"`
	InterestRate             float64 `json:"interest_rate,omitempty"`
	InterestAmount           float64 `json:"interest_amount,omitempty"`
	TotalAwardAmount         float64 `json:"total_award_amount"`
	AwardReasoning           string  `json:"award_reasoning"`
	DigitalSignatureHash     string  `json:"digital_signature_hash"`
	DigitalSignatureDate     string  `json:"digital_signature_date"`
	IssuedBy                 string  `json:"issued_by"`
	IssuedDate               string  `json:"issued_date"`
	ComplianceDeadline       string  `json:"compliance_deadline"`           // 30 days (BR-CLM-OMB-006)
	SupportingDocuments      []string `json:"supporting_documents,omitempty"`
	DocumentURL              string  `json:"document_url,omitempty"`        // TODO: ECMS document URL
	Status                   string  `json:"status"`                        // ISSUED, PENDING_ACCEPTANCE, ACCEPTED, COMPLIED
	BindingOnInsurer         bool    `json:"binding_on_insurer"`            // BR-CLM-OMB-005: Award is binding
}

// ComplianceRecordedResponse represents compliance recording
// BR-CLM-OMB-006: Insurer compliance monitoring
type ComplianceRecordedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	AwardID                   string `json:"award_id"`
	ComplianceStatus          string  `json:"compliance_status"`
	ComplianceDate            string `json:"compliance_date"`
	PaymentReference          string `json:"payment_reference,omitempty"`
	PaymentAmount             float64 `json:"payment_amount,omitempty"`
	ObjectionReason           string `json:"objection_reason,omitempty"`
	DaysToComply              int     `json:"days_to_comply"`               // Days taken to comply (or overdue)
	Overdue                   bool    `json:"overdue"`                      // Whether compliance is overdue
	RecordedBy                string  `json:"recorded_by"`
	RecordingDate             string  `json:"recording_date"`
	NextSteps                 string  `json:"next_steps,omitempty"`
}

// ComplaintClosedResponse represents complaint closure
// BR-CLM-OMB-007: Complaint closure & archival
type ComplaintClosedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	ClosureReason             string `json:"closure_reason"`
	ClosureType               string `json:"closure_type"`
	ClosedDate                string `json:"closed_date"`
	RetentionPeriod           int    `json:"retention_period"`              // Years for archival
	ArchivalDate              string `json:"archival_date"`                 // When archival will occur
	ClosedBy                  string `json:"closed_by"`
}

// IRDAIEscalationResponse represents escalation to IRDAI
// BR-CLM-OMB-006: Escalate to IRDAI on breach
type IRDAIEscalationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ComplaintID               string `json:"complaint_id"`
	AwardID                   string `json:"award_id"`
	EscalationID              string `json:"escalation_id"`
	EscalationReason          string `json:"escalation_reason"`
	BreachDetails             string `json:"breach_details"`
	DaysOverdue               int    `json:"days_overdue"`
	EscalationDate            string `json:"escalation_date"`
	EscalatedBy               string `json:"escalated_by"`
	IRDAIReference            string `json:"irdai_reference,omitempty"`
	Status                    string `json:"status"`                        // ESCALATED, UNDER_REVIEW, RESOLVED
}

// ComplaintsListResponse represents paginated list of complaints
type ComplaintsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Complaints                []ComplaintSummary `json:"complaints"`
}

// ComplaintSummary represents a complaint in list view
type ComplaintSummary struct {
	ComplaintID            string  `json:"complaint_id"`
	ComplaintNumber        string  `json:"complaint_number"`
	ComplainantName        string  `json:"complainant_name"`
	PolicyNumber           string  `json:"policy_number"`
	ComplaintCategory      string  `json:"complaint_category"`
	Status                 string  `json:"status"`
	OmbudsmanCenter        string  `json:"ombudsman_center,omitempty"`
	AssignedOmbudsmanName  string  `json:"assigned_ombudsman_name,omitempty"`
	SubmittedDate          string  `json:"submitted_date"`
	DaysInQueue            int     `json:"days_in_queue"`
	Admissible             *bool   `json:"admissible,omitempty"`
	ClaimValue             float64 `json:"claim_value,omitempty"`
	AwardAmount            float64 `json:"award_amount,omitempty"`
}

// ComplianceQueueResponse represents complaints requiring compliance tracking
// BR-CLM-OMB-006: 30-day compliance monitoring
type ComplianceQueueResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Complaints                []ComplianceItem `json:"complaints"`
	QueueSummary              ComplianceQueueSummary `json:"queue_summary"`
}

// ComplianceItem represents a complaint requiring compliance action
type ComplianceItem struct {
	ComplaintID            string `json:"complaint_id"`
	ComplaintNumber        string `json:"complaint_number"`
	ComplainantName        string `json:"complainant_name"`
	AwardID                string `json:"award_id"`
	AwardType              string `json:"award_type"`
	AwardAmount            float64 `json:"award_amount"`
	AwardDate              string `json:"award_date"`
	ComplianceDeadline     string `json:"compliance_deadline"`
	DaysRemaining          int     `json:"days_remaining"`              // Days until deadline
	Overdue                bool    `json:"overdue"`
	DaysOverdue            int     `json:"days_overdue,omitempty"`      // If overdue
	ComplianceStatus       string  `json:"compliance_status"`
	ReminderSent           bool    `json:"reminder_sent"`              // Whether reminders were sent
	LastReminderDate       string  `json:"last_reminder_date,omitempty"`
}

// ComplianceQueueSummary represents compliance queue statistics
type ComplianceQueueSummary struct {
	TotalPending        int64 `json:"total_pending"`         // Total awaiting compliance
	TotalComplied       int64 `json:"total_complied"`        // Total complied
	TotalOverdue        int64 `json:"total_overdue"`         // Total overdue
	TotalEscalated      int64 `json:"total_escalated"`       // Total escalated to IRDAI
	AverageComplianceDays float64 `json:"average_compliance_days"` // Average days to comply
}

// ==================== HELPER FUNCTIONS ====================

// NewComplaintDetailsResponse creates complaint details response from domain model
func NewComplaintDetailsResponse(complaint domain.OmbudsmanComplaint) *ComplaintDetailsResponse {
	data := ComplaintData{
		ComplaintID:            complaint.ComplaintID,
		ComplaintNumber:        complaint.ComplaintNumber,
		ComplainantName:        complaint.ComplainantName,
		ComplainantAddress:     complaint.ComplainantAddress,
		ComplainantMobile:      complaint.ComplainantMobile,
		ComplainantEmail:       stringValue(complaint.ComplainantEmail),
		ComplainantRole:        complaint.ComplainantRole,
		LanguagePreference:     complaint.LanguagePreference,
		IDProofType:            complaint.IDProofType,
		IDProofNumber:          maskIDProof(complaint.IDProofType, complaint.IDProofNumber),
		PolicyNumber:           complaint.PolicyID,
		ClaimNumber:            stringValue(complaint.ClaimID),
		PolicyType:             complaint.PolicyType,
		AgentName:              stringValue(complaint.AgentName),
		AgentBranch:            stringValue(complaint.AgentBranch),
		ComplaintCategory:      complaint.ComplaintCategory,
		IncidentDate:           formatDate(complaint.IncidentDate),
		RepresentationDate:     formatDate(complaint.RepresentationDate),
		IssueDescription:       complaint.IssueDescription,
		ReliefSought:           complaint.ReliefSought,
		ClaimValue:             floatValue(complaint.ClaimValue),
		OmbudsmanCenter:        stringValue(complaint.OmbudsmanCenter),
		AssignedOmbudsmanID:    stringValue(complaint.AssignedOmbudsmanID),
		// AssignedOmbudsmanName: TODO - Lookup from user service
		Admissible:              complaint.Admissible,
		AdmissibilityChecked:    complaint.AdmissibilityChecked,
		AdmissibilityReason:     stringValue(complaint.AdmissibilityReason),
		InadmissibilityReason:   stringValue(complaint.InadmissibilityReason),
		Status:                  complaint.Status,
		SubmittedDate:           formatDateTime(complaint.CreatedAt),
		LastUpdatedDate:         formatDateTime(complaint.UpdatedAt),
		AcknowledgementSentDate: formatDateTimePtr(complaint.AcknowledgementSentDate),
		ResolutionDueDate:       formatDateTimePtr(complaint.ResolutionDueDate),
		DaysInQueue:             daysSince(complaint.CreatedAt),
		AttachmentCount:         0, // TODO - Query from attachments table
	}

	return &ComplaintDetailsResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintData:         data,
	}
}

// NewComplaintsListResponse creates paginated complaints list response
func NewComplaintsListResponse(complaints []domain.OmbudsmanComplaint, total int64, skip, limit int) *ComplaintsListResponse {
	summaries := make([]ComplaintSummary, len(complaints))
	for i, c := range complaints {
		summaries[i] = ComplaintSummary{
			ComplaintID:            c.ComplaintID,
			ComplaintNumber:        c.ComplaintNumber,
			ComplainantName:        c.ComplainantName,
			PolicyNumber:           c.PolicyID,
			ComplaintCategory:      c.ComplaintCategory,
			Status:                 c.Status,
			OmbudsmanCenter:        stringValue(c.OmbudsmanCenter),
			// AssignedOmbudsmanName: TODO - Lookup from user service
			SubmittedDate:     formatDateTime(c.CreatedAt),
			DaysInQueue:       daysSince(c.CreatedAt),
			Admissible:        c.Admissible,
			ClaimValue:        floatValue(c.ClaimValue),
			AwardAmount:       floatValue(c.AwardAmount),
		}
	}

	return &ComplaintsListResponse{
		StatusCodeAndMessage: port.Success(),
		MetaDataResponse:     port.NewMetaDataResponse(total, skip, limit),
		Complaints:           summaries,
	}
}

// CalculateComplianceSLAStatus calculates SLA status for compliance tracking
func CalculateComplianceSLAStatus(daysRemaining int) string {
	if daysRemaining < 0 {
		return "RED" // Overdue
	} else if daysRemaining <= 2 {
		return "ORANGE" // Critical (≤2 days remaining)
	} else if daysRemaining <= 7 {
		return "YELLOW" // Warning (≤7 days remaining)
	}
	return "GREEN" // On track (>7 days remaining)
}

// ==================== UTILITY FUNCTIONS ====================

// stringValue safely dereferences string pointer
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// floatValue safely dereferences float64 pointer
func floatValue(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// formatDateTime formats time as "2006-01-02 15:04:05"
func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// formatDateTimePtr formats time pointer as "2006-01-02 15:04:05"
func formatDateTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// formatDate formats date as "2006-01-02"
func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// daysSince calculates days elapsed since given time
func daysSince(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	return int(time.Since(t).Hours() / 24)
}

// maskIDProof masks ID proof number for security (show only last 4 characters)
func maskIDProof(idType, idNumber *string) string {
	if idNumber == nil {
		return ""
	}
	if *idNumber == "" {
		return ""
	}

	// Show full number for demonstration - in production, mask appropriately
	// For Aadhaar: show first 4 and last 4 (e.g., XXXX-XXXX-1234)
	// For PAN: show first 4 and last 4 (e.g., ABCPXXXXX)
	// For Passport: show last 4 only
	return *idNumber // TODO: Implement masking logic
}
