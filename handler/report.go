package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"

	"gitlab.cept.gov.in/pli/claims-api/handler/response"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"

	"gitlab.cept.gov.in/it-2.0-common/n-api-log"
)

// ReportHandler handles report generation and analytics endpoints
// Reference: FR-CLM-RPT-001 to FR-CLM-RPT-004
type ReportHandler struct {
	ClaimRepo         *domain.ClaimRepository
	PaymentRepo       *domain.ClaimPaymentRepository
	InvestigationRepo *domain.InvestigationRepository
	HistoryRepo       *domain.ClaimHistoryRepository
}

// NewReportHandler creates a new ReportHandler instance
func NewReportHandler(
	claimRepo *domain.ClaimRepository,
	paymentRepo *domain.ClaimPaymentRepository,
	investigationRepo *domain.InvestigationRepository,
	historyRepo *domain.ClaimHistoryRepository,
) *ReportHandler {
	return &ReportHandler{
		ClaimRepo:         claimRepo,
		PaymentRepo:       paymentRepo,
		InvestigationRepo: investigationRepo,
		HistoryRepo:       historyRepo,
	}
}

// ========================================
// ROUTE REGISTRATION
// ========================================

var (
	// POST /reports/claims/generate - Generate claim reports
	GenerateClaimReportRoute = serverRoute.Route{
		Method:  "POST",
		Path:    "/reports/claims/generate",
		Handler: "GenerateClaimReport",
	}

	// POST /reports/dashboard/generate - Generate dashboard reports
	GenerateDashboardReportRoute = serverRoute.Route{
		Method:  "POST",
		Path:    "/reports/dashboard/generate",
		Handler: "GenerateDashboardReport",
	}

	// GET /reports/statistics - Get claim statistics
	GetClaimStatisticsRoute = serverRoute.Route{
		Method:  "GET",
		Path:    "/reports/statistics",
		Handler: "GetClaimStatistics",
	}

	// GET /reports/sla-compliance - Get SLA compliance report
	GetSlaComplianceReportRoute = serverRoute.Route{
		Method:  "GET",
		Path:    "/reports/sla-compliance",
		Handler: "GetSlaComplianceReport",
	}

	// GET /reports/payments - Get payment report
	GetPaymentReportRoute = serverRoute.Route{
		Method:  "GET",
		Path:    "/reports/payments",
		Handler: "GetPaymentReport",
	}

	// GET /reports/investigation - Get investigation report
	GetInvestigationReportRoute = serverRoute.Route{
		Method:  "GET",
		Path:    "/reports/investigation",
		Handler: "GetInvestigationReport",
	}

	// GET /reports/audit-trail/:entity_id/:entity_type - Get audit trail
	GetAuditTrailRoute = serverRoute.Route{
		Method:  "GET",
		Path:    "/reports/audit-trail/:entity_id/:entity_type",
		Handler: "GetAuditTrail",
	}

	// POST /reports/export - Export report
	ExportReportRoute = serverRoute.Route{
		Method:  "POST",
		Path:    "/reports/export",
		Handler: "ExportReport",
	}
)

// RegisterRoutes registers all report routes
func (h *ReportHandler) RegisterRoutes(r serverHandler.RouteRegistrar) {
	r.RegisterRoute(&GenerateClaimReportRoute)
	r.RegisterRoute(&GenerateDashboardReportRoute)
	r.RegisterRoute(&GetClaimStatisticsRoute)
	r.RegisterRoute(&GetSlaComplianceReportRoute)
	r.RegisterRoute(&GetPaymentReportRoute)
	r.RegisterRoute(&GetInvestigationReportRoute)
	r.RegisterRoute(&GetAuditTrailRoute)
	r.RegisterRoute(&ExportReportRoute)
}

// ========================================
// HANDLER METHODS
// ========================================

// GenerateClaimReport generates various types of claim reports
// POST /reports/claims/generate
// Reference: FR-CLM-RPT-001
func (h *ReportHandler) GenerateClaimReport(sctx *serverRoute.Context, req GenerateClaimReportRequest) (*resp.ClaimReportGeneratedResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Generating claim report: type=%s, start=%s, end=%s", req.ReportType, req.StartDate, req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		log.Error(ctx, "Invalid start date format: %v", err)
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		log.Error(ctx, "Invalid end date format: %v", err)
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Build filters based on report type
	filters := map[string]interface{}{
		"created_at >=": startDate,
		"created_at <=": endDate,
	}

	if req.ClaimType != nil {
		filters["claim_type"] = *req.ClaimType
	}

	if req.Division != nil {
		filters["division"] = *req.Division
	}

	if req.District != nil {
		filters["district"] = *req.District
	}

	// Add status filter based on report type
	switch req.ReportType {
	case "CLAIMS_PENDING":
		filters["status"] = []string{"REGISTERED", "UNDER_INVESTIGATION", "PENDING_APPROVAL"}
	case "CLAIMS_APPROVED":
		filters["status"] = "APPROVED"
	case "CLAIMS_REJECTED":
		filters["status"] = "REJECTED"
	case "CLAIMS_PAID":
		filters["status"] = "PAID"
	case "SLA_BREACH":
		filters["sla_breached"] = true
	case "INVESTIGATION_SUMMARY":
		filters["investigation_required"] = true
	}

	// Fetch claims from repository
	claims, total, err := h.ClaimRepo.List(ctx, filters, 0, 10000, "created_at", "DESC")
	if err != nil {
		log.Error(ctx, "Failed to fetch claims for report: %v", err)
		return nil, fmt.Errorf("failed to fetch claims: %w", err)
	}

	// Generate report ID
	reportID, _ := generateReportID()

	// TODO: Generate actual report file (PDF/Excel)
	// This would integrate with a report generation service
	var reportURL *string
	if req.IncludeDetails != nil && *req.IncludeDetails {
		url := fmt.Sprintf("https://reports.pli.gov.in/%s.pdf", reportID)
		reportURL = &url
	}

	log.Info(ctx, "Claim report generated: report_id=%s, total_records=%d", reportID, total)

	return resp.NewClaimReportGeneratedResponse(reportID, req.ReportType, total, reportURL), nil
}

// GenerateDashboardReport generates dashboard reports with key metrics
// POST /reports/dashboard/generate
// Reference: FR-CLM-RPT-002
func (h *ReportHandler) GenerateDashboardReport(sctx *serverRoute.Context, req GenerateDashboardReportRequest) (*resp.DashboardReportResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Generating dashboard report: type=%s, date=%s", req.ReportType, req.ReportDate)

	// Parse report date
	reportDate, err := time.Parse("2006-01-02", req.ReportDate)
	if err != nil {
		log.Error(ctx, "Invalid report date format: %v", err)
		return nil, fmt.Errorf("invalid report date format: %w", err)
	}

	// Calculate date range based on report type
	var startDate, endDate time.Time
	switch req.ReportType {
	case "DAILY":
		startDate = reportDate
		endDate = reportDate.Add(24 * time.Hour)
	case "WEEKLY":
		weekday := reportDate.Weekday()
		startDate = reportDate.AddDate(0, 0, -int(weekday))
		endDate = startDate.AddDate(0, 0, 7)
	case "MONTHLY":
		startDate = time.Date(reportDate.Year(), reportDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
	case "QUARTERLY":
		quarter := (reportDate.Month() - 1) / 3
		startDate = time.Date(reportDate.Year(), quarter*3+1, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 3, 0)
	}

	// Fetch dashboard metrics
	filters := map[string]interface{}{
		"created_at >=": startDate,
		"created_at <=": endDate,
	}

	claims, total, err := h.ClaimRepo.List(ctx, filters, 0, 100000, "created_at", "DESC")
	if err != nil {
		log.Error(ctx, "Failed to fetch claims for dashboard: %v", err)
		return nil, fmt.Errorf("failed to fetch claims: %w", err)
	}

	// Calculate metrics
	metrics := h.calculateDashboardMetrics(ctx, claims)

	// Generate report ID
	reportID, _ := generateReportID()

	log.Info(ctx, "Dashboard report generated: report_id=%s, total_claims=%d", reportID, total)

	return resp.NewDashboardReportResponse(reportID, req.ReportType, req.ReportDate, metrics), nil
}

// GetClaimStatistics retrieves claim statistics grouped by specified dimension
// GET /reports/statistics
// Reference: FR-CLM-RPT-003
func (h *ReportHandler) GetClaimStatistics(sctx *serverRoute.Context, req GetClaimStatisticsRequest) (*resp.ClaimStatisticsResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Fetching claim statistics: start=%s, end=%s, group_by=%s", req.StartDate, req.EndDate, req.GroupBy)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// TODO: Implement statistics aggregation based on GroupBy parameter
	// This would require either:
	// 1. Complex SQL queries with GROUP BY
	// 2. Application-side aggregation
	// For now, return summary statistics

	filters := map[string]interface{}{
		"created_at >=": startDate,
		"created_at <=": endDate,
	}

	claims, total, err := h.ClaimRepo.List(ctx, filters, 0, 100000, "created_at", "DESC")
	if err != nil {
		log.Error(ctx, "Failed to fetch claims for statistics: %v", err)
		return nil, fmt.Errorf("failed to fetch claims: %w", err)
	}

	// Calculate summary statistics
	summary := h.calculateSummaryStatistics(ctx, claims)

	// Group data (placeholder - actual grouping would be more complex)
	statistics := []resp.StatisticsData{
		{
			Key:   req.GroupBy,
			Value: req.StartDate + " to " + req.EndDate,
			Count: total,
			Amount: summary.TotalAmount,
			Approved: int64(summary.ApprovalRate * float64(total) / 100),
			Rejected: int64(summary.RejectionRate * float64(total) / 100),
			Pending: total - int64(summary.ApprovalRate*float64(total)/100) - int64(summary.RejectionRate*float64(total)/100),
			AvgTime: summary.AvgProcessingTime,
		},
	}

	return resp.NewClaimStatisticsResponse(req.StartDate, req.EndDate, req.GroupBy, statistics, summary), nil
}

// GetSlaComplianceReport retrieves SLA compliance report
// GET /reports/sla-compliance
// Reference: BR-CLM-DC-003, BR-CLM-DC-004, BR-CLM-DC-021
func (h *ReportHandler) GetSlaComplianceReport(sctx *serverRoute.Context, req GetSlaReportRequest) (*resp.SlaComplianceReportResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Fetching SLA compliance report: sla_type=%s, start=%s, end=%s", req.SlaType, req.StartDate, req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Build filters
	filters := map[string]interface{}{
		"created_at >=": startDate,
		"created_at <=": endDate,
	}

	if req.Division != nil {
		filters["division"] = *req.Division
	}

	if req.BreachedOnly != nil && *req.BreachedOnly {
		filters["sla_breached"] = true
	}

	// Fetch claims
	claims, total, err := h.ClaimRepo.List(ctx, filters, 0, 100000, "created_at", "DESC")
	if err != nil {
		log.Error(ctx, "Failed to fetch claims for SLA report: %v", err)
		return nil, fmt.Errorf("failed to fetch claims: %w", err)
	}

	// Calculate SLA compliance data
	complianceData := h.calculateSlaCompliance(ctx, claims)

	// Calculate summary
	breachCount := int64(0)
	for _, data := range complianceData {
		if data.Breached {
			breachCount++
		}
	}

	complianceRate := float64(total-breachCount) / float64(total) * 100
	breachRate := float64(breachCount) / float64(total) * 100

	summary := resp.SlaComplianceSummary{
		OnTimeCount:  total - breachCount,
		DelayedCount: 0, // TODO: Calculate delayed count
		BreachCount:  breachCount,
	}

	return resp.NewSlaComplianceReportResponse(req.StartDate, req.EndDate, req.SlaType, total, complianceRate, breachRate, breachCount, complianceData, summary), nil
}

// GetPaymentReport retrieves payment report
// GET /reports/payments
// Reference: BR-CLM-PAY-001
func (h *ReportHandler) GetPaymentReport(sctx *serverRoute.Context, req GetPaymentReportRequest) (*resp.PaymentReportResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Fetching payment report: start=%s, end=%s", req.StartDate, req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// TODO: Fetch payment data from ClaimPaymentRepository
	// For now, return placeholder response

	totalRecords := int64(0)
	totalAmount := float64(0)
	successCount := int64(0)
	failedCount := int64(0)
	pendingCount := int64(0)

	paymentData := []resp.PaymentReportData{}
	summary := resp.PaymentReportSummary{}

	return resp.NewPaymentReportResponse(req.StartDate, req.EndDate, totalRecords, totalAmount, successCount, failedCount, pendingCount, paymentData, summary), nil
}

// GetInvestigationReport retrieves investigation report
// GET /reports/investigation
// Reference: BR-CLM-DC-002
func (h *ReportHandler) GetInvestigationReport(sctx *serverRoute.Context, req GetInvestigationReportRequest) (*resp.InvestigationReportResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Fetching investigation report: start=%s, end=%s", req.StartDate, req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Build filters
	filters := map[string]interface{}{
		"created_at >=": startDate,
		"created_at <=": endDate,
	}

	if req.InvestigationStatus != nil {
		filters["status"] = *req.InvestigationStatus
	}

	if req.Outcome != nil {
		filters["outcome"] = *req.Outcome
	}

	if req.InvestigatorID != nil {
		filters["investigator_id"] = *req.InvestigatorID
	}

	// Fetch investigations
	investigations, total, err := h.InvestigationRepo.List(ctx, filters, 0, 10000, "created_at", "DESC")
	if err != nil {
		log.Error(ctx, "Failed to fetch investigations for report: %v", err)
		return nil, fmt.Errorf("failed to fetch investigations: %w", err)
	}

	// Transform to report data
	investigationData := make([]resp.InvestigationReportListData, len(investigations))
	for i, inv := range investigations {
		investigationData[i] = resp.InvestigationReportListData{
			InvestigationID: inv.InvestigationID,
			ClaimNumber:     inv.ClaimID, // TODO: Fetch claim number
			PolicyID:        inv.PolicyID,
			InvestigatorID:  inv.InvestigatorID,
			Status:          inv.Status,
			AssignedDate:    inv.AssignedDate.Format("2006-01-02"),
			DueDate:         inv.DueDate.Format("2006-01-02"),
			// TODO: Add more fields
		}
	}

	// Calculate summary
	summary := resp.InvestigationReportSummary{}

	return resp.NewInvestigationReportResponse(req.StartDate, req.EndDate, total, investigationData, summary), nil
}

// GetAuditTrail retrieves audit trail for an entity
// GET /reports/audit-trail/:entity_id/:entity_type
// Reference: Security requirements - Audit trail
func (h *ReportHandler) GetAuditTrail(sctx *serverRoute.Context, req GetAuditTrailRequest) (*resp.AuditTrailResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Fetching audit trail: entity_id=%s, entity_type=%s", req.EntityID, req.EntityType)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// TODO: Fetch audit trail from claim_history table
	// For now, return placeholder response
	auditEntries := []resp.AuditEntry{}
	totalRecords := int64(0)

	return resp.NewAuditTrailResponse(req.EntityID, req.EntityType, req.StartDate, req.EndDate, totalRecords, auditEntries), nil
}

// ExportReport exports a report in specified format
// POST /reports/export
// Reference: FR-CLM-RPT-004
func (h *ReportHandler) ExportReport(sctx *serverRoute.Context, req ExportReportRequest) (*resp.ReportExportedResponse, error) {
	ctx := sctx.Ctx

	log.Info(ctx, "Exporting report: report_id=%s, format=%s", req.ReportID, req.Format)

	// Generate export ID
	exportID, _ := generateReportID()

	// TODO: Generate export file in specified format (PDF/Excel/CSV)
	// This would integrate with a report generation service

	var downloadURL *string
	url := fmt.Sprintf("https://reports.pli.gov.in/exports/%s.%s", exportID, req.Format)
	downloadURL = &url

	var emailSent *bool
	if req.EmailTo != nil {
		sent := true // TODO: Send email with report attachment
		emailSent = &sent
	}

	status := "COMPLETED"

	return resp.NewReportExportedResponse(exportID, req.ReportID, req.Format, status, downloadURL, emailSent), nil
}

// ========================================
// HELPER METHODS
// ========================================

// calculateDashboardMetrics calculates dashboard metrics from claims
func (h *ReportHandler) calculateDashboardMetrics(ctx context.Context, claims []domain.Claim) resp.DashboardMetrics {
	metrics := resp.DashboardMetrics{}

	for _, claim := range claims {
		metrics.ClaimsRegistered++

		switch claim.Status {
		case "APPROVED":
			metrics.ClaimsApproved++
		case "REJECTED":
			metrics.ClaimsRejected++
		case "PAID":
			metrics.ClaimsPaid++
		case "REGISTERED", "UNDER_INVESTIGATION", "PENDING_APPROVAL":
			metrics.PendingApprovals++
		}

		if claim.SLABreached != nil && *claim.SLABreached {
			metrics.OverdueClaims++
		}

		if claim.ClaimType == "DEATH" {
			metrics.ClaimsByType.Death++
		} else if claim.ClaimType == "MATURITY" {
			metrics.ClaimsByType.Maturity++
		} else if claim.ClaimType == "SURVIVAL_BENEFIT" {
			metrics.ClaimsByType.SurvivalBenefit++
		} else if claim.ClaimType == "FREELOOK" {
			metrics.ClaimsByType.FreeLook++
		}
	}

	// TODO: Calculate average processing time, payment disbursed, SLA compliance

	return metrics
}

// calculateSummaryStatistics calculates summary statistics
func (h *ReportHandler) calculateSummaryStatistics(ctx context.Context, claims []domain.Claim) resp.StatisticsSummary {
	summary := resp.StatisticsSummary{
		TotalClaims: int64(len(claims)),
	}

	approvedCount := int64(0)
	rejectedCount := int64(0)
	totalAmount := float64(0)

	for _, claim := range claims {
		totalAmount += claim.ClaimAmount

		if claim.Status == "APPROVED" || claim.Status == "PAID" {
			approvedCount++
		} else if claim.Status == "REJECTED" {
			rejectedCount++
		}
	}

	summary.TotalAmount = totalAmount

	if summary.TotalClaims > 0 {
		summary.ApprovalRate = float64(approvedCount) / float64(summary.TotalClaims) * 100
		summary.RejectionRate = float64(rejectedCount) / float64(summary.TotalClaims) * 100
	}

	return summary
}

// calculateSlaCompliance calculates SLA compliance data from claims
func (h *ReportHandler) calculateSlaCompliance(ctx context.Context, claims []domain.Claim) []resp.SlaComplianceData {
	complianceData := make([]resp.SlaComplianceData, 0, len(claims))

	for _, claim := range claims {
		data := resp.SlaComplianceData{
			ClaimNumber: claim.ClaimNumber,
			PolicyID:    claim.PolicyID,
			ClaimType:   claim.ClaimType,
			Status:      claim.Status,
			// TODO: Calculate due date, completion date, days taken, SLA status
		}

		if claim.SLAStatus != nil {
			data.SlaStatus = *claim.SLAStatus
		}

		if claim.SLABreached != nil {
			data.Breached = *claim.SLABreached
		}

		if claim.Division != nil {
			data.Division = claim.Division
		}

		complianceData = append(complianceData, data)
	}

	return complianceData
}

// generateReportID generates a unique report ID
func generateReportID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "RPT" + hex.EncodeToString(bytes), nil
}
