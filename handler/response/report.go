package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	"time"
)

// ========================================
// REPORTING & ANALYTICS - RESPONSE DTOS
// ========================================

// ClaimReportGeneratedResponse represents the response for claim report generation
// POST /reports/claims/generate
// Reference: FR-CLM-RPT-001
type ClaimReportGeneratedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ReportID                  string    `json:"report_id"`
	ReportType                string    `json:"report_type"`
	GeneratedAt               string    `json:"generated_at"`
	TotalRecords              int64     `json:"total_records"`
	ReportURL                 *string   `json:"report_url,omitempty"`
	ExpiresAt                 *string   `json:"expires_at,omitempty"`
}

// ClaimReportData represents the detailed claim report data
type ClaimReportData struct {
	ClaimNumber     string  `json:"claim_number"`
	PolicyID        string  `json:"policy_id"`
	ClaimType       string  `json:"claim_type"`
	Status          string  `json:"status"`
	ClaimAmount     float64 `json:"claim_amount"`
	ApprovedAmount  *float64 `json:"approved_amount,omitempty"`
	ClaimantName    string  `json:"claimant_name"`
	RegistrationDate string `json:"registration_date"`
	DaysPending     *int    `json:"days_pending,omitempty"`
	SLAStatus       *string `json:"sla_status,omitempty"`
	Division        *string `json:"division,omitempty"`
	District        *string `json:"district,omitempty"`
}

// DashboardReportResponse represents the response for dashboard report generation
// POST /reports/dashboard/generate
// Reference: FR-CLM-RPT-002
type DashboardReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ReportID                  string             `json:"report_id"`
	ReportType                string             `json:"report_type"`
	ReportDate                string             `json:"report_date"`
	GeneratedAt               string             `json:"generated_at"`
	Metrics                   DashboardMetrics   `json:"metrics"`
}

// DashboardMetrics represents the dashboard metrics
type DashboardMetrics struct {
	ClaimsRegistered      int64               `json:"claims_registered"`
	ClaimsApproved        int64               `json:"claims_approved"`
	ClaimsRejected        int64               `json:"claims_rejected"`
	ClaimsPaid            int64               `json:"claims_paid"`
	PendingApprovals      int64               `json:"pending_approvals"`
	OverdueClaims         int64               `json:"overdue_claims"`
	AvgProcessingTime     float64             `json:"avg_processing_time"` // In days
	PaymentDisbursed      float64             `json:"payment_disbursed"`    // In INR
	SLAComplianceRate     float64             `json:"sla_compliance_rate"`  // Percentage
	ClaimsByType          ClaimsByType        `json:"claims_by_type"`
	ClaimsByStatus        ClaimsByStatus      `json:"claims_by_status"`
	ClaimsByDivision      []ClaimsByDivision  `json:"claims_by_division,omitempty"`
	TopPerformingDivisions []TopDivision      `json:"top_performing_divisions,omitempty"`
}

// ClaimsByType represents claims breakdown by type
type ClaimsByType struct {
	Death            int64 `json:"death"`
	Maturity         int64 `json:"maturity"`
	SurvivalBenefit  int64 `json:"survival_benefit"`
	FreeLook         int64 `json:"free_look"`
}

// ClaimsByStatus represents claims breakdown by status
type ClaimsByStatus struct {
	Registered      int64 `json:"registered"`
	UnderInvestigation int64 `json:"under_investigation"`
	PendingApproval int64 `json:"pending_approval"`
	Approved        int64 `json:"approved"`
	Rejected        int64 `json:"rejected"`
	Paid            int64 `json:"paid"`
	Closed          int64 `json:"closed"`
}

// ClaimsByDivision represents claims by division
type ClaimsByDivision struct {
	Division        string  `json:"division"`
	TotalClaims     int64   `json:"total_claims"`
	ApprovedClaims  int64   `json:"approved_claims"`
	PendingClaims   int64   `json:"pending_claims"`
	SLACompliance   float64 `json:"sla_compliance"` // Percentage
}

// TopDivision represents top performing division
type TopDivision struct {
	Division        string  `json:"division"`
	ClaimsProcessed int64   `json:"claims_processed"`
	AvgProcessingTime float64 `json:"avg_processing_time"` // In days
	SLACompliance   float64 `json:"sla_compliance"`        // Percentage
}

// ClaimStatisticsResponse represents the response for claim statistics
// GET /reports/statistics
// Reference: FR-CLM-RPT-003
type ClaimStatisticsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	StartDate                 string              `json:"start_date"`
	EndDate                   string              `json:"end_date"`
	GroupBy                   string              `json:"group_by"`
	Statistics                []StatisticsData    `json:"statistics"`
	Summary                   StatisticsSummary    `json:"summary"`
}

// StatisticsData represents statistics data point
type StatisticsData struct {
	Key         string  `json:"key"`         // day, week, month, division, etc.
	Value       string  `json:"value"`       // Formatted date or division name
	Count       int64   `json:"count"`
	Amount      float64 `json:"amount"`      // Total claim amount
	Approved    int64   `json:"approved"`
	Rejected    int64   `json:"rejected"`
	Pending     int64   `json:"pending"`
	AvgTime     float64 `json:"avg_time"`    // Average processing time in days
}

// StatisticsSummary represents summary statistics
type StatisticsSummary struct {
	TotalClaims      int64   `json:"total_claims"`
	TotalAmount      float64 `json:"total_amount"`
	ApprovalRate     float64 `json:"approval_rate"`     // Percentage
	RejectionRate    float64 `json:"rejection_rate"`    // Percentage
	AvgProcessingTime float64 `json:"avg_processing_time"` // In days
}

// SlaComplianceReportResponse represents the response for SLA compliance report
// GET /reports/sla-compliance
// Reference: BR-CLM-DC-003, BR-CLM-DC-004, BR-CLM-DC-021
type SlaComplianceReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	StartDate                 string              `json:"start_date"`
	EndDate                   string              `json:"end_date"`
	SlaType                   string              `json:"sla_type"`
	TotalRecords              int64               `json:"total_records"`
	ComplianceRate            float64             `json:"compliance_rate"` // Percentage
	BreachCount               int64               `json:"breach_count"`
	BreachRate                float64             `json:"breach_rate"`      // Percentage
	ComplianceData            []SlaComplianceData `json:"compliance_data"`
	Summary                   SlaComplianceSummary `json:"summary"`
}

// SlaComplianceData represents SLA compliance data point
type SlaComplianceData struct {
	ClaimNumber      string  `json:"claim_number"`
	PolicyID         string  `json:"policy_id"`
	ClaimType        string  `json:"claim_type"`
	Status           string  `json:"status"`
	DueDate          string  `json:"due_date"`
	CompletedDate    *string `json:"completed_date,omitempty"`
	DaysTaken        *int    `json:"days_taken,omitempty"`
	SlaLimit         int     `json:"sla_limit"`      // In days
	SlaStatus        string  `json:"sla_status"`     // GREEN, YELLOW, ORANGE, RED
	Breached         bool    `json:"breached"`
	BreachDays       *int    `json:"breach_days,omitempty"`
	Division         *string `json:"division,omitempty"`
}

// SlaComplianceSummary represents SLA compliance summary
type SlaComplianceSummary struct {
	OnTimeCount     int64   `json:"on_time_count"`
	DelayedCount    int64   `json:"delayed_count"`
	BreachCount     int64   `json:"breach_count"`
	AvgProcessingTime float64 `json:"avg_processing_time"` // In days
	MaxBreachDays   int     `json:"max_breach_days"`
}

// PaymentReportResponse represents the response for payment report
// GET /reports/payments
// Reference: BR-CLM-PAY-001
type PaymentReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	StartDate                 string            `json:"start_date"`
	EndDate                   string            `json:"end_date"`
	TotalRecords              int64             `json:"total_records"`
	TotalAmount               float64           `json:"total_amount"`
	SuccessCount              int64             `json:"success_count"`
	FailedCount               int64             `json:"failed_count"`
	PendingCount              int64             `json:"pending_count"`
	PaymentData               []PaymentReportData `json:"payment_data"`
	Summary                   PaymentReportSummary `json:"summary"`
}

// PaymentReportData represents payment report data point
type PaymentReportData struct {
	PaymentID        string  `json:"payment_id"`
	ClaimNumber      string  `json:"claim_number"`
	PolicyID         string  `json:"policy_id"`
	ClaimantName     string  `json:"claimant_name"`
	Amount           float64 `json:"amount"`
	PaymentMode      string  `json:"payment_mode"`
	BankAccount      *string `json:"bank_account,omitempty"`
	BankIFSC         *string `json:"bank_ifsc,omitempty"`
	Status           string  `json:"status"`
	UTRNumber        *string `json:"utr_number,omitempty"`
	InitiatedDate    string  `json:"initiated_date"`
	CompletedDate    *string `json:"completed_date,omitempty"`
	DaysTaken        *int    `json:"days_taken,omitempty"`
	Division         *string `json:"division,omitempty"`
}

// PaymentReportSummary represents payment report summary
type PaymentReportSummary struct {
	TotalDisbursed    float64 `json:"total_disbursed"`
	TotalPending      float64 `json:"total_pending"`
	TotalFailed       float64 `json:"total_failed"`
	SuccessRate       float64 `json:"success_rate"` // Percentage
	AvgProcessingTime float64 `json:"avg_processing_time"` // In days
	PaymentsByMode    []PaymentByMode `json:"payments_by_mode"`
}

// PaymentByMode represents payment breakdown by mode
type PaymentByMode struct {
	PaymentMode string  `json:"payment_mode"`
	Count       int64   `json:"count"`
	Amount      float64 `json:"amount"`
	Percentage  float64 `json:"percentage"`
}

// InvestigationReportResponse represents the response for investigation report
// GET /reports/investigation
// Reference: BR-CLM-DC-002
type InvestigationReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	StartDate                 string                          `json:"start_date"`
	EndDate                   string                          `json:"end_date"`
	TotalRecords              int64                           `json:"total_records"`
	InvestigationData         []InvestigationReportListData   `json:"investigation_data"`
	Summary                   InvestigationReportSummary       `json:"summary"`
}

// InvestigationReportListData represents investigation report data point for list reports
type InvestigationReportListData struct {
	InvestigationID    string  `json:"investigation_id"`
	ClaimNumber        string  `json:"claim_number"`
	PolicyID           string  `json:"policy_id"`
	InvestigatorID     string  `json:"investigator_id"`
	InvestigatorName   *string `json:"investigator_name,omitempty"`
	Status             string  `json:"status"`
	AssignedDate       string  `json:"assigned_date"`
	DueDate            string  `json:"due_date"`
	CompletedDate      *string `json:"completed_date,omitempty"`
	Outcome            *string `json:"outcome,omitempty"` // CLEAR, SUSPECT, FRAUD
	DaysTaken          *int    `json:"days_taken,omitempty"`
	SLAStatus          *string `json:"sla_status,omitempty"` // GREEN, YELLOW, RED
	Breached           *bool   `json:"breached,omitempty"`
	ReinvestigationCount *int  `json:"reinvestigation_count,omitempty"`
}

// InvestigationReportSummary represents investigation report summary
type InvestigationReportSummary struct {
	TotalInvestigations int64   `json:"total_investigations"`
	CompletedCount      int64   `json:"completed_count"`
	PendingCount        int64   `json:"pending_count"`
	OverdueCount        int64   `json:"overdue_count"`
	AvgDuration         float64 `json:"avg_duration"` // In days
	Outcomes            InvestigationOutcomes `json:"outcomes"`
	TopInvestigators    []TopInvestigator `json:"top_investigators,omitempty"`
}

// InvestigationOutcomes represents investigation outcomes breakdown
type InvestigationOutcomes struct {
	ClearCount  int64   `json:"clear_count"`
	SuspectCount int64  `json:"suspect_count"`
	FraudCount  int64   `json:"fraud_count"`
	PendingCount int64  `json:"pending_count"`
}

// TopInvestigator represents top performing investigator
type TopInvestigator struct {
	InvestigatorID   string  `json:"investigator_id"`
	InvestigatorName *string `json:"investigator_name,omitempty"`
	TotalAssigned    int64   `json:"total_assigned"`
	TotalCompleted   int64   `json:"total_completed"`
	AvgDuration      float64 `json:"avg_duration"` // In days
	SLACompliance    float64 `json:"sla_compliance"` // Percentage
}

// AuditTrailResponse represents the response for audit trail
// GET /reports/audit-trail
// Reference: Security requirements - Audit trail
type AuditTrailResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	EntityID                  string          `json:"entity_id"`
	EntityType                string          `json:"entity_type"`
	StartDate                 string          `json:"start_date"`
	EndDate                   string          `json:"end_date"`
	TotalRecords              int64           `json:"total_records"`
	AuditEntries              []AuditEntry    `json:"audit_entries"`
}

// AuditEntry represents an audit trail entry
type AuditEntry struct {
	ID              string    `json:"id"`
	EntityID        string    `json:"entity_id"`
	EntityType      string    `json:"entity_type"`
	Action          string    `json:"action"`
	ActionType      string    `json:"action_type"`
	OldValue        *string   `json:"old_value,omitempty"`
	NewValue        *string   `json:"new_value,omitempty"`
	ChangedBy       *string   `json:"changed_by,omitempty"`
	ChangedByName   *string   `json:"changed_by_name,omitempty"`
	ChangedAt       string    `json:"changed_at"`
	IPAddress       *string   `json:"ip_address,omitempty"`
	UserAgent       *string   `json:"user_agent,omitempty"`
	SessionID       *string   `json:"session_id,omitempty"`
	Reason          *string   `json:"reason,omitempty"`
	Metadata        *string   `json:"metadata,omitempty"` // JSON metadata
}

// ReportExportedResponse represents the response for report export
// POST /reports/export
// Reference: FR-CLM-RPT-004
type ReportExportedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ExportID                  string  `json:"export_id"`
	ReportID                  string  `json:"report_id"`
	Format                    string  `json:"format"`
	Status                    string  `json:"status"`
	GeneratedAt               string  `json:"generated_at"`
	DownloadURL               *string `json:"download_url,omitempty"`
	ExpiresAt                 *string `json:"expires_at,omitempty"`
	EmailSent                 *bool   `json:"email_sent,omitempty"`
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// NewClaimReportGeneratedResponse creates a new claim report response
func NewClaimReportGeneratedResponse(reportID, reportType string, totalRecords int64, reportURL *string) *ClaimReportGeneratedResponse {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Report expires in 24 hours

	return &ClaimReportGeneratedResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Report generated successfully",
		},
		ReportID:     reportID,
		ReportType:   reportType,
		GeneratedAt:  now.Format("2006-01-02 15:04:05"),
		TotalRecords: totalRecords,
		ReportURL:    reportURL,
		ExpiresAt:    stringPtr(expiresAt.Format("2006-01-02 15:04:05")),
	}
}

// NewDashboardReportResponse creates a new dashboard report response
func NewDashboardReportResponse(reportID, reportType, reportDate string, metrics DashboardMetrics) *DashboardReportResponse {
	return &DashboardReportResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Dashboard report generated successfully",
		},
		ReportID:    reportID,
		ReportType:  reportType,
		ReportDate:  reportDate,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Metrics:     metrics,
	}
}

// NewClaimStatisticsResponse creates a new claim statistics response
func NewClaimStatisticsResponse(startDate, endDate, groupBy string, statistics []StatisticsData, summary StatisticsSummary) *ClaimStatisticsResponse {
	return &ClaimStatisticsResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Statistics retrieved successfully",
		},
		StartDate:  startDate,
		EndDate:    endDate,
		GroupBy:    groupBy,
		Statistics: statistics,
		Summary:    summary,
	}
}

// NewSlaComplianceReportResponse creates a new SLA compliance report response
func NewSlaComplianceReportResponse(startDate, endDate, slaType string, totalRecords int64, complianceRate, breachRate float64, breachCount int64, data []SlaComplianceData, summary SlaComplianceSummary) *SlaComplianceReportResponse {
	return &SlaComplianceReportResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "SLA compliance report generated successfully",
		},
		StartDate:      startDate,
		EndDate:        endDate,
		SlaType:        slaType,
		TotalRecords:   totalRecords,
		ComplianceRate: complianceRate,
		BreachCount:    breachCount,
		BreachRate:     breachRate,
		ComplianceData: data,
		Summary:        summary,
	}
}

// NewPaymentReportResponse creates a new payment report response
func NewPaymentReportResponse(startDate, endDate string, totalRecords int64, totalAmount float64, successCount, failedCount, pendingCount int64, data []PaymentReportData, summary PaymentReportSummary) *PaymentReportResponse {
	return &PaymentReportResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Payment report generated successfully",
		},
		StartDate:    startDate,
		EndDate:      endDate,
		TotalRecords: totalRecords,
		TotalAmount:  totalAmount,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		PendingCount: pendingCount,
		PaymentData:  data,
		Summary:      summary,
	}
}

// NewInvestigationReportResponse creates a new investigation report response
func NewInvestigationReportResponse(startDate, endDate string, totalRecords int64, data []InvestigationReportData, summary InvestigationReportSummary) *InvestigationReportResponse {
	return &InvestigationReportResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Investigation report generated successfully",
		},
		StartDate:         startDate,
		EndDate:           endDate,
		TotalRecords:      totalRecords,
		InvestigationData: data,
		Summary:           summary,
	}
}

// NewAuditTrailResponse creates a new audit trail response
func NewAuditTrailResponse(entityID, entityType, startDate, endDate string, totalRecords int64, entries []AuditEntry) *AuditTrailResponse {
	return &AuditTrailResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Audit trail retrieved successfully",
		},
		EntityID:     entityID,
		EntityType:   entityType,
		StartDate:    startDate,
		EndDate:      endDate,
		TotalRecords: totalRecords,
		AuditEntries: entries,
	}
}

// NewReportExportedResponse creates a new report export response
func NewReportExportedResponse(exportID, reportID, format, status string, downloadURL *string, emailSent *bool) *ReportExportedResponse {
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // Export expires in 7 days

	return &ReportExportedResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Report exported successfully",
		},
		ExportID:    exportID,
		ReportID:    reportID,
		Format:      format,
		Status:      status,
		GeneratedAt: now.Format("2006-01-02 15:04:05"),
		DownloadURL: downloadURL,
		ExpiresAt:   stringPtr(expiresAt.Format("2006-01-02 15:04:05")),
		EmailSent:   emailSent,
	}
}

// stringPtr returns a pointer to string
func stringPtr(s string) *string {
	return &s
}
