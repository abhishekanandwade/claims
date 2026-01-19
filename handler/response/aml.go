package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== AML/CFT RESPONSE DTOs ====================

// AMLTriggerDetectionResponse represents the response for AML trigger detection
// POST /aml/detect-trigger
// Reference: FR-CLM-AML-001, BR-CLM-AML-001, BR-CLM-AML-002, BR-CLM-AML-003
type AMLTriggerDetectionResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	TriggerDetected           bool     `json:"trigger_detected"`
	TriggerTypes              []string `json:"trigger_types,omitempty"`
	RiskLevel                 string   `json:"risk_level,omitempty"` // LOW, MEDIUM, HIGH, CRITICAL
	RiskScore                 *float64 `json:"risk_score,omitempty"`
	AlertID                   *string  `json:"alert_id,omitempty"`
	TriggerReasons            []string `json:"trigger_reasons,omitempty"`
	RecommendedActions        []string `json:"recommended_actions,omitempty"`
	TransactionBlocked        bool     `json:"transaction_blocked"`
	FilingRequired            bool     `json:"filing_required,omitempty"` // STR, CTR filing required
	FilingType                *string  `json:"filing_type,omitempty"`     // STR, CTR, CCR, NTR
}

// AMLAlertResponse represents a single AML alert
type AMLAlertResponse struct {
	AlertID                 string  `json:"alert_id"`
	PolicyID                string  `json:"policy_id"`
	CustomerID              *string `json:"customer_id,omitempty"`
	TransactionType         string  `json:"transaction_type"`
	TransactionAmount       *float64 `json:"transaction_amount,omitempty"`
	TransactionDate         string  `json:"transaction_date"`
	TriggerCode             string  `json:"trigger_code"`
	RiskLevel               string  `json:"risk_level"`            // LOW, MEDIUM, HIGH, CRITICAL
	RiskScore               *int    `json:"risk_score,omitempty"`  // 0-100
	AlertStatus             string  `json:"alert_status"`          // FLAGGED, UNDER_REVIEW, FILED, CLOSED
	AlertDescription        *string `json:"alert_description,omitempty"`
	FilingRequired          bool    `json:"filing_required"`
	FilingType              *string `json:"filing_type,omitempty"` // STR, CTR, CCR, NTR
	FilingStatus            *string `json:"filing_status,omitempty"`
	FilingReference         *string `json:"filing_reference,omitempty"`
	TransactionBlocked      bool    `json:"transaction_blocked"`
	ReviewDecision          *string `json:"review_decision,omitempty"`
	OfficerRemarks          *string `json:"officer_remarks,omitempty"`
	ReviewedBy              *string `json:"reviewed_by,omitempty"`
	PANNumber               *string `json:"pan_number,omitempty"`
	PANVerified             *bool   `json:"pan_verified,omitempty"`
	NomineeChangeDetected   bool    `json:"nominee_change_detected"`
	CreatedAt               string  `json:"created_at"`
	UpdatedAt               string  `json:"updated_at"`
}

// NewAMLAlertResponse creates a new AMLAlertResponse from domain.AMLAlert
func NewAMLAlertResponse(alert domain.AMLAlert) AMLAlertResponse {
	resp := AMLAlertResponse{
		AlertID:               alert.AlertID,
		PolicyID:              alert.PolicyID,
		CustomerID:            alert.CustomerID,
		TransactionType:       alert.TransactionType,
		TransactionAmount:     alert.TransactionAmount,
		TransactionDate:       alert.TransactionDate.Format("2006-01-02"),
		TriggerCode:           alert.TriggerCode,
		RiskLevel:             alert.RiskLevel,
		RiskScore:             alert.RiskScore,
		AlertStatus:           alert.AlertStatus,
		AlertDescription:      alert.AlertDescription,
		FilingRequired:        alert.FilingRequired,
		TransactionBlocked:    alert.TransactionBlocked,
		NomineeChangeDetected: alert.NomineeChangeDetected,
		CreatedAt:             alert.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             alert.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if alert.FilingType != nil {
		resp.FilingType = alert.FilingType
	}
	if alert.FilingStatus != nil {
		resp.FilingStatus = alert.FilingStatus
	}
	if alert.FilingReference != nil {
		resp.FilingReference = alert.FilingReference
	}
	if alert.ReviewDecision != nil {
		resp.ReviewDecision = alert.ReviewDecision
	}
	if alert.OfficerRemarks != nil {
		resp.OfficerRemarks = alert.OfficerRemarks
	}
	if alert.ReviewedBy != nil {
		resp.ReviewedBy = alert.ReviewedBy
	}
	if alert.PANNumber != nil {
		resp.PANNumber = alert.PANNumber
	}
	if alert.PANVerified != nil {
		resp.PANVerified = alert.PANVerified
	}

	return resp
}

// AMLAlertsListResponse represents the response for listing AML alerts
type AMLAlertsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Data                      []AMLAlertResponse `json:"data"`
}

// NewAMLAlertsListResponse creates a new AMLAlertsListResponse
func NewAMLAlertsListResponse(alerts []domain.AMLAlert, total int64, skip, limit int) *AMLAlertsListResponse {
	data := make([]AMLAlertResponse, len(alerts))
	for i, alert := range alerts {
		data[i] = NewAMLAlertResponse(alert)
	}

	return &AMLAlertsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse:     port.NewMetaDataResponse(uint64(skip), uint64(limit), uint64(total)),
		Data:                 data,
	}
}

// AMLAlertGeneratedResponse represents the response when an AML alert is generated
type AMLAlertGeneratedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	AlertID                   string `json:"alert_id"`
	TriggerCode               string `json:"trigger_code"`
	RiskLevel                 string `json:"risk_level"`
	RiskScore                 float64 `json:"risk_score"`
	FilingRequired            bool    `json:"filing_required"`
	TransactionBlocked        bool    `json:"transaction_blocked"`
	Message                   string  `json:"message"`
}

// RiskScoreCalculationResponse represents the response for risk score calculation
// POST /aml/{alert_id}/calculate-risk-score
// Reference: BR-CLM-AML-004 (Risk Scoring Algorithm)
type RiskScoreCalculationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	AlertID                   string  `json:"alert_id"`
	RiskScore                 float64 `json:"risk_score"`           // 0-100
	RiskLevel                 string  `json:"risk_level"`           // LOW, MEDIUM, HIGH, CRITICAL
	RiskFactors               []RiskFactor `json:"risk_factors"`
	CalculationBreakdown      []RiskFactor  `json:"calculation_breakdown"`
	RecommendedActions        []string `json:"recommended_actions"`
	FilingRequired            bool     `json:"filing_required"`
	TransactionBlocked        bool     `json:"transaction_blocked"`
}

// RiskFactor represents a single risk factor
type RiskFactor struct {
	Factor      string  `json:"factor"`
	Weight      float64 `json:"weight"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Reference   string  `json:"reference,omitempty"` // Business rule reference
}

// AMLAlertDetailsResponse represents the detailed AML alert information
// GET /aml/{alert_id}/details
type AMLAlertDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	AMLAlertResponse          `json:",inline"`
	TriggerDetails            TriggerDetailsData `json:"trigger_details"`
	RiskAnalysis              RiskAnalysisData   `json:"risk_analysis"`
	FilingInformation         *FilingInfoData    `json:"filing_information,omitempty"`
	TransactionHistory        []TransactionHistoryItem `json:"transaction_history,omitempty"`
	CustomerHistory           *CustomerHistoryData `json:"customer_history,omitempty"`
}

// TriggerDetailsData contains detailed trigger information
type TriggerDetailsData struct {
	TriggerCode          string   `json:"trigger_code"`
	TriggerName          string   `json:"trigger_name"`
	TriggerDescription   *string  `json:"trigger_description,omitempty"`
	TriggerCategory      string   `json:"trigger_category"`      // CASH_THRESHOLD, PAN_MISMATCH, NOMINEE_CHANGE, etc.
	ApplicableRules      []string `json:"applicable_rules"`      // BR-CLM-AML-001, etc.
	Severity             string   `json:"severity"`               // LOW, MEDIUM, HIGH, CRITICAL
	RegulationReference  string   `json:"regulation_reference"`   // PMLA 2002 Section X
}

// RiskAnalysisData contains risk analysis information
type RiskAnalysisData struct {
	OverallRiskScore    float64           `json:"overall_risk_score"`
	RiskLevel           string            `json:"risk_level"`
	RiskFactors         []RiskFactor      `json:"risk_factors"`
	RiskTrend           string            `json:"risk_trend"`            // INCREASING, STABLE, DECREASING
	PeersRiskLevel      string            `json:"peers_risk_level"`      // Risk level compared to peers
	IndustryBenchmark   float64           `json:"industry_benchmark"`    // Industry average risk score
	RecommendedActions  []string          `json:"recommended_actions"`
}

// FilingInfoData contains STR/CTR filing information
type FilingInfoData struct {
	FilingType      string  `json:"filing_type"`      // STR, CTR, CCR, NTR
	FilingStatus    string  `json:"filing_status"`    // PENDING, FILED, ACKNOWLEDGED, REJECTED
	FilingDeadline  string  `json:"filing_deadline"`
	FilingDate      *string `json:"filing_date,omitempty"`
	ReportingAgency string  `json:"reporting_agency"` // FINNET, FINGATE
	FilingReference *string `json:"filing_reference,omitempty"`
	Acknowledgement *string `json:"acknowledgement,omitempty"`
}

// TransactionHistoryItem represents a transaction in history
type TransactionHistoryItem struct {
	TransactionID   string  `json:"transaction_id"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transaction_date"`
	PaymentMode     string  `json:"payment_mode"`
	Status          string  `json:"status"`
	AlertTriggered  bool    `json:"alert_triggered"`
	AlertID         *string `json:"alert_id,omitempty"`
}

// CustomerHistoryData contains customer risk history
type CustomerHistoryData struct {
	CustomerID            string  `json:"customer_id"`
	TotalAlerts          int     `json:"total_alerts"`
	HighRiskAlerts       int     `json:"high_risk_alerts"`
	CriticalRiskAlerts   int     `json:"critical_risk_alerts"`
	TotalTransactions    int     `json:"total_transactions"`
	SuspiciousActivities int     `json:"suspicious_activities"`
	AverageRiskScore     float64 `json:"average_risk_score"`
	RiskTrend            string  `json:"risk_trend"` // INCREASING, STABLE, DECREASING
}

// AMLAlertReviewResponse represents the response for AML alert review
// POST /aml/{alert_id}/review
// Reference: BR-CLM-AML-005 (Alert Review)
type AMLAlertReviewResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	AlertID                   string  `json:"alert_id"`
	ReviewDecision            string `json:"review_decision"`
	ReviewStatus              string `json:"review_status"` // APPROVED, REJECTED, ESCALATED
	OfficerID                 string `json:"officer_id"`
	OfficerRemarks            string `json:"officer_remarks"`
	EscalationLevel           *string `json:"escalation_level,omitempty"`
	TransactionBlocked        bool    `json:"transaction_blocked"`
	FilingRequired            bool    `json:"filing_required"`
	NextSteps                 []string `json:"next_steps"`
	ReviewedAt                string  `json:"reviewed_at"`
}

// AMLReportFiledResponse represents the response when STR/CTR is filed
// POST /aml/{alert_id}/file-report
// Reference: BR-CLM-AML-006 (STR Filing Within 7 Days)
// Reference: BR-CLM-AML-007 (CTR Filing Monthly)
type AMLReportFiledResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	AlertID                   string  `json:"alert_id"`
	ReportType                string  `json:"report_type"`                // STR, CTR, CCR, NTR
	FilingReference           string  `json:"filing_reference"`
	ReportingAgency           string  `json:"reporting_agency"`           // FINNET, FINGATE
	FilingStatus              string  `json:"filing_status"`              // FILED, PENDING, FAILED
	FilingDate                string  `json:"filing_date"`
	Acknowledgement           *string `json:"acknowledgement,omitempty"`
	TransactionBlocked        bool    `json:"transaction_blocked"`
	AlertStatus               string  `json:"alert_status"`               // FILED, CLOSED
	Message                   string  `json:"message"`
}

// AMLAlertQueueResponse represents the response for AML alert queue
// This is the 7th endpoint for listing alerts requiring action
type AMLAlertQueueResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Data                      []AMLAlertResponse `json:"data"`
	Summary                  AMLQueueSummary `json:"summary"`
}

// AMLQueueSummary contains summary statistics for the AML queue
type AMLQueueSummary struct {
	TotalAlerts        int64   `json:"total_alerts"`
	PendingReview      int64   `json:"pending_review"`
	HighRisk           int64   `json:"high_risk"`
	CriticalRisk       int64   `json:"critical_risk"`
	FilingOverdue      int64   `json:"filing_overdue"`
	TransactionBlocked int64   `json:"transaction_blocked"`
	AverageRiskScore   float64 `json:"average_risk_score"`
}

// NewAMLAlertQueueResponse creates a new AMLAlertQueueResponse
func NewAMLAlertQueueResponse(alerts []domain.AMLAlert, total int64, skip, limit int, summary AMLQueueSummary) *AMLAlertQueueResponse {
	data := make([]AMLAlertResponse, len(alerts))
	for i, alert := range alerts {
		data[i] = NewAMLAlertResponse(alert)
	}

	return &AMLAlertQueueResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse:     port.NewMetaDataResponse(uint64(skip), uint64(limit), uint64(total)),
		Data:                 data,
		Summary:              summary,
	}
}
