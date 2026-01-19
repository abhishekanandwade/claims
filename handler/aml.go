package handler

import (
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// AMLHandler handles AML/CFT-related HTTP requests
type AMLHandler struct {
	*serverHandler.Base
	svc *repo.AMLAlertRepository
}

// NewAMLHandler creates a new AML handler
func NewAMLHandler(svc *repo.AMLAlertRepository) *AMLHandler {
	base := serverHandler.New("AML").
		SetPrefix("/v1").
		AddPrefix("")
	return &AMLHandler{
		Base: base,
		svc:  svc,
	}
}

// Routes defines all routes for this handler
func (h *AMLHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// AML/CFT Endpoints (7 endpoints)
		serverRoute.POST("/aml/detect-trigger", h.DetectAMLTrigger).Name("Detect AML Trigger"),
		serverRoute.POST("/aml/:alert_id/generate-alert", h.GenerateAMLAlert).Name("Generate AML Alert"),
		serverRoute.POST("/aml/:alert_id/calculate-risk-score", h.CalculateAMLRiskScore).Name("Calculate AML Risk Score"),
		serverRoute.GET("/aml/:alert_id/details", h.GetAMLAlertDetails).Name("Get AML Alert Details"),
		serverRoute.POST("/aml/:alert_id/review", h.ReviewAMLAlert).Name("Review AML Alert"),
		serverRoute.POST("/aml/:alert_id/file-report", h.FileAMLReport).Name("File AML Report"),
		serverRoute.GET("/aml/queue/pending-review", h.GetPendingReviewQueue).Name("Get Pending Review Queue"),
	}
}

// DetectAMLTrigger detects AML trigger conditions during transactions
// POST /aml/detect-trigger
// Reference: FR-CLM-AML-001, BR-CLM-AML-001 (High Cash Premium Alert)
// Reference: BR-CLM-AML-002 (PAN Mismatch Alert)
// Reference: BR-CLM-AML-003 (Nominee Change Post Death)
func (h *AMLHandler) DetectAMLTrigger(sctx *serverRoute.Context, req DetectAMLTriggerRequest) (*resp.AMLTriggerDetectionResponse, error) {
	log.Info(sctx.Ctx, "Detecting AML triggers for transaction_type: %s, amount: %.2f", req.TransactionType, req.TransactionAmount)

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid transaction_date format: %v", err)
		return nil, err
	}

	// Evaluate AML triggers using business rules
	// This is a placeholder - in Task 5.3, we'll implement 70+ AML rules
	triggerDetected := false
	var triggerTypes []string
	var triggerReasons []string
	var recommendedActions []string
	riskLevel := "LOW"
	var riskScore float64 = 0.0
	filingRequired := false
	var filingType *string
	transactionBlocked := false

	// BR-CLM-AML-001: Cash transactions over ₹50,000 trigger high-risk alert and CTR filing
	if req.PaymentMode != nil && *req.PaymentMode == "CASH" && req.TransactionAmount > 50000 {
		triggerDetected = true
		triggerTypes = append(triggerTypes, "CASH_THRESHOLD")
		triggerReasons = append(triggerReasons, "Cash transaction over ₹50,000")
		riskLevel = "HIGH"
		riskScore = 75.0
		filingRequired = true
		ft := "CTR"
		filingType = &ft
		recommendedActions = append(recommendedActions, "File CTR within monthly deadline", "Review transaction documentation")
	}

	// BR-CLM-AML-002: PAN verification failure triggers medium-risk alert
	// TODO: Integrate with Customer Service for PAN verification
	if req.PANNumber == nil {
		triggerDetected = true
		if len(triggerTypes) == 0 {
			triggerTypes = append(triggerTypes, "PAN_MISSING")
		} else {
			triggerTypes = append(triggerTypes, "PAN_MISSING")
		}
		triggerReasons = append(triggerReasons, "PAN number not provided")
		if riskLevel == "LOW" {
			riskLevel = "MEDIUM"
			riskScore = 50.0
		}
		recommendedActions = append(recommendedActions, "Verify PAN with Customer Service", "Request PAN documentation")
	}

	// Additional trigger checks will be implemented in Task 5.3
	// BR-CLM-AML-003: Nominee change after policyholder death
	// BR-CLM-AML-004: Multiple high-value transactions
	// BR-CLM-AML-005: Geographical risk indicators
	// ... and 65+ more rules

	// Determine if transaction should be blocked
	if riskLevel == "CRITICAL" {
		transactionBlocked = true
	}

	// If trigger detected, create alert
	var alertID *string
	if triggerDetected {
		// Convert riskScore to int for domain model
		riskScoreInt := int(riskScore)

		// Generate a default PolicyID if not provided (for transaction-level alerts)
		policyID := "SYSTEM"
		if req.PolicyID != nil {
			policyID = *req.PolicyID
		}

		alert := domain.AMLAlert{
			PolicyID:            policyID,
			CustomerID:          req.CustomerID,
			TransactionType:     req.TransactionType,
			TransactionAmount:   &req.TransactionAmount,
			TransactionDate:     transactionDate,
			TriggerCode:         triggerTypes[0], // Primary trigger
			AlertDescription:    &triggerReasons[0],
			RiskLevel:           riskLevel,
			RiskScore:           &riskScoreInt,
			AlertStatus:         "FLAGGED",
			FilingRequired:      filingRequired,
			FilingType:          filingType,
			TransactionBlocked:  transactionBlocked,
		}

		createdAlert, err := h.svc.Create(sctx.Ctx, alert)
		if err != nil {
			log.Error(sctx.Ctx, "Error creating AML alert: %v", err)
			return nil, err
		}

		alertID = &createdAlert.AlertID
		log.Info(sctx.Ctx, "AML alert created: %s", createdAlert.AlertID)
	}

	response := &resp.AMLTriggerDetectionResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		TriggerDetected:      triggerDetected,
		TriggerTypes:         triggerTypes,
		RiskLevel:            riskLevel,
		RiskScore:            &riskScore,
		AlertID:              alertID,
		TriggerReasons:       triggerReasons,
		RecommendedActions:   recommendedActions,
		TransactionBlocked:   transactionBlocked,
		FilingRequired:       filingRequired,
		FilingType:           filingType,
	}

	return response, nil
}

// GenerateAMLAlert generates an AML alert from detected trigger
// POST /aml/{alert_id}/generate-alert
func (h *AMLHandler) GenerateAMLAlert(sctx *serverRoute.Context, req AlertIDUri) (*resp.AMLAlertGeneratedResponse, error) {
	log.Info(sctx.Ctx, "Generating AML alert: %s", req.AlertID)

	// This endpoint is typically called internally after DetectAMLTrigger
	// For now, we'll retrieve the alert and return it
	alert, err := h.svc.FindByAlertID(sctx.Ctx, req.AlertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "AML alert not found: %s", req.AlertID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding AML alert: %v", err)
		return nil, err
	}

	response := &resp.AMLAlertGeneratedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		AlertID:              alert.AlertID,
		TriggerCode:          alert.TriggerCode,
		RiskLevel:            alert.RiskLevel,
		RiskScore:            float64(*alert.RiskScore),
		FilingRequired:       alert.FilingRequired,
		TransactionBlocked:   alert.TransactionBlocked,
		Message:              "AML alert generated successfully",
	}

	return response, nil
}

// CalculateAMLRiskScore calculates risk score for an AML alert
// POST /aml/{alert_id}/calculate-risk-score
// Reference: BR-CLM-AML-004 (Risk Scoring Algorithm)
func (h *AMLHandler) CalculateAMLRiskScore(sctx *serverRoute.Context, req AlertIDUri) (*resp.RiskScoreCalculationResponse, error) {
	log.Info(sctx.Ctx, "Calculating risk score for alert: %s", req.AlertID)

	// Retrieve alert
	alert, err := h.svc.FindByAlertID(sctx.Ctx, req.AlertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "AML alert not found: %s", req.AlertID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding AML alert: %v", err)
		return nil, err
	}

	// TODO: Implement comprehensive risk scoring algorithm in Task 5.3
	// For now, use existing risk score from alert
	riskScore := 0.0
	if alert.RiskScore != nil {
		riskScore = float64(*alert.RiskScore)
	}

	// Determine risk level based on score
	riskLevel := "LOW"
	if riskScore >= 80 {
		riskLevel = "CRITICAL"
	} else if riskScore >= 60 {
		riskLevel = "HIGH"
	} else if riskScore >= 40 {
		riskLevel = "MEDIUM"
	}

	// Risk factors (placeholder - will be implemented in Task 5.3)
	riskFactors := []resp.RiskFactor{
		{
			Factor:      "Transaction Amount",
			Weight:      0.3,
			Score:       riskScore * 0.3,
			Description: "High-value transaction detected",
			Reference:   "BR-CLM-AML-001",
		},
		{
			Factor:      "Payment Mode",
			Weight:      0.2,
			Score:       riskScore * 0.2,
			Description: "Cash payment risk",
			Reference:   "BR-CLM-AML-001",
		},
	}

	// Determine recommended actions based on risk level
	recommendedActions := []string{}
	if riskLevel == "CRITICAL" {
		recommendedActions = append(recommendedActions, "Block transaction immediately", "Initiate STR filing within 7 days", "Escalate to compliance officer")
	} else if riskLevel == "HIGH" {
		recommendedActions = append(recommendedActions, "Review transaction documentation", "Consider blocking transaction", "Prepare for STR filing")
	} else if riskLevel == "MEDIUM" {
		recommendedActions = append(recommendedActions, "Conduct manual review", "Verify customer identity", "Monitor future transactions")
	}

	response := &resp.RiskScoreCalculationResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		AlertID:              alert.AlertID,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		RiskFactors:          riskFactors,
		CalculationBreakdown: riskFactors,
		RecommendedActions:   recommendedActions,
		FilingRequired:       alert.FilingRequired,
		TransactionBlocked:   alert.TransactionBlocked,
	}

	return response, nil
}

// GetAMLAlertDetails retrieves detailed information about an AML alert
// GET /aml/{alert_id}/details
func (h *AMLHandler) GetAMLAlertDetails(sctx *serverRoute.Context, req AlertIDUri) (*resp.AMLAlertDetailsResponse, error) {
	log.Info(sctx.Ctx, "Retrieving AML alert details: %s", req.AlertID)

	// Retrieve alert
	alert, err := h.svc.FindByAlertID(sctx.Ctx, req.AlertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "AML alert not found: %s", req.AlertID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding AML alert: %v", err)
		return nil, err
	}

	// Build alert response
	alertResp := resp.NewAMLAlertResponse(alert)

	// Build trigger details
	triggerDetails := resp.TriggerDetailsData{
		TriggerCode:         alert.TriggerCode,
		TriggerName:         getTriggerName(alert.TriggerCode),
		TriggerDescription:  alert.AlertDescription,
		TriggerCategory:     getTriggerCategory(alert.TriggerCode),
		ApplicableRules:     getApplicableRules(alert.TriggerCode),
		Severity:            alert.RiskLevel,
		RegulationReference: "PMLA 2002",
	}

	// Build risk analysis
	// Convert RiskScore from *int to float64
	riskScore := 0.0
	if alert.RiskScore != nil {
		riskScore = float64(*alert.RiskScore)
	}

	riskAnalysis := resp.RiskAnalysisData{
		OverallRiskScore: riskScore,
		RiskLevel:        alert.RiskLevel,
		RiskFactors:      []resp.RiskFactor{}, // Will be populated in Task 5.3
		RiskTrend:        "STABLE",
		PeersRiskLevel:   "MEDIUM",
		IndustryBenchmark: 45.0,
		RecommendedActions: []string{
			"Review transaction details",
			"Verify customer documentation",
		},
	}

	// Build filing information if applicable
	var filingInfo *resp.FilingInfoData
	if alert.FilingRequired && alert.FilingType != nil {
		filingStatus := "PENDING"
		if alert.FilingStatus != nil {
			filingStatus = *alert.FilingStatus
		}

		// Calculate filing deadline based on filing type
		// STR: 7 days, CTR: 30 days
		filingDeadline := ""
		if alert.FilingType != nil {
			if *alert.FilingType == "STR" {
				filingDeadline = alert.CreatedAt.AddDate(0, 0, 7).Format("2006-01-02")
			} else if *alert.FilingType == "CTR" {
				filingDeadline = alert.CreatedAt.AddDate(0, 0, 30).Format("2006-01-02")
			}
		}

		filingInfo = &resp.FilingInfoData{
			FilingType:      *alert.FilingType,
			FilingStatus:    filingStatus,
			FilingDeadline:  filingDeadline,
			ReportingAgency: "FINNET", // Default, will be configured
		}
	}

	// TODO: Fetch transaction history (will be implemented in later phases)
	transactionHistory := []resp.TransactionHistoryItem{}

	// TODO: Fetch customer risk history (will be implemented in later phases)
	var customerHistory *resp.CustomerHistoryData

	response := &resp.AMLAlertDetailsResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		AMLAlertResponse:     alertResp,
		TriggerDetails:       triggerDetails,
		RiskAnalysis:         riskAnalysis,
		FilingInformation:    filingInfo,
		TransactionHistory:   transactionHistory,
		CustomerHistory:      customerHistory,
	}

	return response, nil
}

// ReviewAMLAlert reviews an AML alert and takes action
// POST /aml/{alert_id}/review
// Reference: BR-CLM-AML-005 (Alert Review)
func (h *AMLHandler) ReviewAMLAlert(sctx *serverRoute.Context, req ReviewAMLAlertRequest) (*resp.AMLAlertReviewResponse, error) {
	log.Info(sctx.Ctx, "Reviewing AML alert: %s, decision: %s", req.AlertID, req.ReviewDecision)

	// Retrieve alert
	alert, err := h.svc.FindByAlertID(sctx.Ctx, req.AlertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "AML alert not found: %s", req.AlertID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding AML alert: %v", err)
		return nil, err
	}

	// Update alert with review decision
	// Determine alert status and next steps based on review decision
	var nextSteps []string
	var alertStatus string
	var transactionBlocked bool = alert.TransactionBlocked
	filingRequired := alert.FilingRequired
	reviewStatus := "APPROVED"

	switch req.ReviewDecision {
	case "CLEAR":
		alertStatus = "CLOSED"
		transactionBlocked = false
		reviewStatus = "APPROVED"
		nextSteps = append(nextSteps, "Alert closed - no further action required", "Transaction can proceed")

	case "FILE_STR":
		alertStatus = "FILED"
		reviewStatus = "APPROVED"
		nextSteps = append(nextSteps, "File STR within 7 days", "Block transaction if not already blocked", "Document all supporting evidence")

	case "FILE_CTR":
		alertStatus = "FILED"
		reviewStatus = "APPROVED"
		nextSteps = append(nextSteps, "File CTR in monthly batch", "Update CTR log")

	case "BLOCK_TRANSACTION":
		alertStatus = "FILED"
		transactionBlocked = true
		reviewStatus = "APPROVED"
		nextSteps = append(nextSteps, "Transaction blocked", "Notify user of block", "Escalate to compliance officer")

	case "ESCALATE":
		alertStatus = "UNDER_REVIEW"
		reviewStatus = "ESCALATED"
		nextSteps = append(nextSteps, "Escalated to senior officer", "Await review")
	}

	// Update alert using repository method
	_, err = h.svc.Update(sctx.Ctx, req.AlertID,
		nil, // riskLevel
		nil, // riskScore
		&alertStatus, // alertStatus
		nil, // alertDescription
		&req.OfficerID, // reviewedBy
		&req.ReviewDecision, // reviewDecision
		&req.OfficerRemarks, // officerRemarks
		nil, // actionTaken
		&transactionBlocked, // transactionBlocked
		nil, // filingRequired
		nil, // filingType
		nil, // filingStatus
		nil, // filingReference
		nil, // filedAt
		nil, // filedBy
	)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating AML alert: %v", err)
		return nil, err
	}

	response := &resp.AMLAlertReviewResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		AlertID:              req.AlertID,
		ReviewDecision:       req.ReviewDecision,
		ReviewStatus:         reviewStatus,
		OfficerID:            req.OfficerID,
		OfficerRemarks:       req.OfficerRemarks,
		EscalationLevel:      req.EscalationLevel,
		TransactionBlocked:   transactionBlocked,
		FilingRequired:       filingRequired,
		NextSteps:            nextSteps,
		ReviewedAt:           time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Info(sctx.Ctx, "AML alert review completed: %s, status: %s", req.AlertID, reviewStatus)
	return response, nil
}

// FileAMLReport files STR/CTR with regulatory authorities
// POST /aml/{alert_id}/file-report
// Reference: BR-CLM-AML-006 (STR Filing Within 7 Days)
// Reference: BR-CLM-AML-007 (CTR Filing Monthly)
func (h *AMLHandler) FileAMLReport(sctx *serverRoute.Context, req FileAMLReportRequest) (*resp.AMLReportFiledResponse, error) {
	log.Info(sctx.Ctx, "Filing %s for alert: %s", req.ReportType, req.AlertID)

	// Parse filing date
	filingDate, err := time.Parse("2006-01-02", req.FilingDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid filing_date format: %v", err)
		return nil, err
	}

	// Retrieve alert
	alert, err := h.svc.FindByAlertID(sctx.Ctx, req.AlertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "AML alert not found: %s", req.AlertID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding AML alert: %v", err)
		return nil, err
	}

	// TODO: Integrate with FINNET/FINGATE API for actual filing
	// For now, we'll mark as filed
	filingStatus := "FILED"
	acknowledgement := req.FilingReference + "-ACK"

	// Update alert
	_, err = h.svc.UpdateFiling(sctx.Ctx, req.AlertID, req.ReportType, req.FilingReference, req.FiledBy)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating AML alert filing: %v", err)
		return nil, err
	}

	// If report type is STR, close alert after filing
	alertStatus := "FILED"
	if req.ReportType == "STR" {
		alertStatus = "CLOSED"
		err = h.svc.UpdateStatus(sctx.Ctx, req.AlertID, alertStatus, req.FiledBy, nil)
		if err != nil {
			log.Error(sctx.Ctx, "Error updating alert status: %v", err)
			return nil, err
		}
	}

	response := &resp.AMLReportFiledResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		AlertID:              req.AlertID,
		ReportType:           req.ReportType,
		FilingReference:      req.FilingReference,
		ReportingAgency:      req.ReportingAgency,
		FilingStatus:         filingStatus,
		FilingDate:           filingDate.Format("2006-01-02"),
		Acknowledgement:      &acknowledgement,
		TransactionBlocked:   alert.TransactionBlocked,
		AlertStatus:          alertStatus,
		Message:              req.ReportType + " filed successfully with " + req.ReportingAgency,
	}

	log.Info(sctx.Ctx, "%s filed successfully for alert: %s", req.ReportType, req.AlertID)
	return response, nil
}

// GetPendingReviewQueue retrieves AML alerts pending review
// GET /aml/queue/pending-review
// This is the 7th endpoint for queue management
func (h *AMLHandler) GetPendingReviewQueue(sctx *serverRoute.Context, req port.MetadataRequest) (*resp.AMLAlertQueueResponse, error) {
	log.Info(sctx.Ctx, "Retrieving AML pending review queue")

	// Retrieve pending review alerts
	alerts, err := h.svc.GetPendingReviewAlerts(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Error retrieving pending review alerts: %v", err)
		return nil, err
	}

	// For now, return all alerts without pagination
	// TODO: Add pagination support to repository methods
	total := int64(len(alerts))

	// Calculate summary statistics
	summary, err := h.calculateQueueSummary(sctx)
	if err != nil {
		log.Error(sctx.Ctx, "Error calculating queue summary: %v", err)
		return nil, err
	}

	response := resp.NewAMLAlertQueueResponse(alerts, total, int(req.Skip), int(req.Limit), summary)
	return response, nil
}

// Helper functions

// getTriggerName returns the human-readable name for a trigger code
func getTriggerName(triggerCode string) string {
	triggerNames := map[string]string{
		"CASH_THRESHOLD":       "High Cash Transaction",
		"PAN_MISSING":          "PAN Not Provided",
		"PAN_MISMATCH":         "PAN Verification Failed",
		"NOMINEE_CHANGE":       "Nominee Change After Death",
		"MULTIPLE_TRANS":       "Multiple High-Value Transactions",
		"GEO_RISK":             "Geographical Risk",
		"SUSPICIOUS_PATTERN":   "Suspicious Transaction Pattern",
	}

	if name, ok := triggerNames[triggerCode]; ok {
		return name
	}
	return "Unknown Trigger"
}

// getTriggerCategory returns the category for a trigger code
func getTriggerCategory(triggerCode string) string {
	categories := map[string]string{
		"CASH_THRESHOLD":     "CASH_THRESHOLD",
		"PAN_MISSING":        "KYC_VERIFICATION",
		"PAN_MISMATCH":       "KYC_VERIFICATION",
		"NOMINEE_CHANGE":     "FRAUD_INDICATOR",
		"MULTIPLE_TRANS":     "PATTERN_ANALYSIS",
		"GEO_RISK":           "GEOGRAPHICAL_RISK",
		"SUSPICIOUS_PATTERN": "PATTERN_ANALYSIS",
	}

	if category, ok := categories[triggerCode]; ok {
		return category
	}
	return "OTHER"
}

// getApplicableRules returns the applicable business rules for a trigger code
func getApplicableRules(triggerCode string) []string {
	rules := map[string][]string{
		"CASH_THRESHOLD":     {"BR-CLM-AML-001"},
		"PAN_MISSING":        {"BR-CLM-AML-002"},
		"PAN_MISMATCH":       {"BR-CLM-AML-002"},
		"NOMINEE_CHANGE":     {"BR-CLM-AML-003"},
		"MULTIPLE_TRANS":     {"BR-CLM-AML-004"},
		"GEO_RISK":           {"BR-CLM-AML-005"},
		"SUSPICIOUS_PATTERN": {"BR-CLM-AML-004"},
	}

	if ruleList, ok := rules[triggerCode]; ok {
		return ruleList
	}
	return []string{}
}

// calculateQueueSummary calculates summary statistics for the AML queue
func (h *AMLHandler) calculateQueueSummary(sctx *serverRoute.Context) (resp.AMLQueueSummary, error) {
	// Get pending review alerts
	pendingReviewAlerts, err := h.svc.GetPendingReviewAlerts(sctx.Ctx)
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Get high-risk alerts
	highRiskAlerts, err := h.svc.GetHighRiskAlerts(sctx.Ctx)
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Get alerts requiring filing (any filing type)
	_, err = h.svc.GetAlertsRequiringFiling(sctx.Ctx, "")
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Get overdue filing alerts (7 days overdue)
	overdueFilingAlerts, err := h.svc.GetOverdueFilingAlerts(sctx.Ctx, 7)
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Get blocked transactions
	blockedTransactions, err := h.svc.GetBlockedTransactions(sctx.Ctx)
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Calculate statistics
	totalAlerts := int64(len(pendingReviewAlerts))
	pendingReview := int64(len(pendingReviewAlerts))
	highRisk := int64(len(highRiskAlerts))
	criticalRisk := int64(0) // Will be calculated from highRisk alerts
	filingOverdue := int64(len(overdueFilingAlerts))
	transactionBlocked := int64(len(blockedTransactions))

	// Count critical risk alerts
	for _, alert := range highRiskAlerts {
		if alert.RiskLevel == "CRITICAL" {
			criticalRisk++
		}
	}

	// Get risk score distribution for average calculation
	riskDist, err := h.svc.GetRiskScoreDistribution(sctx.Ctx)
	if err != nil {
		return resp.AMLQueueSummary{}, err
	}

	// Calculate average risk score
	var totalScore int64
	var totalCount int64
	for score, count := range riskDist {
		// Parse score from key (e.g., "HIGH", "MEDIUM", "LOW")
		// For now, use a simple approximation
		switch score {
		case "CRITICAL":
			totalScore += 90 * count
		case "HIGH":
			totalScore += 70 * count
		case "MEDIUM":
			totalScore += 50 * count
		case "LOW":
			totalScore += 30 * count
		}
		totalCount += count
	}

	averageRiskScore := 0.0
	if totalCount > 0 {
		averageRiskScore = float64(totalScore) / float64(totalCount)
	}

	summary := resp.AMLQueueSummary{
		TotalAlerts:        totalAlerts,
		PendingReview:      pendingReview,
		HighRisk:           highRisk,
		CriticalRisk:       criticalRisk,
		FilingOverdue:      filingOverdue,
		TransactionBlocked: transactionBlocked,
		AverageRiskScore:   averageRiskScore,
	}

	return summary, nil
}
