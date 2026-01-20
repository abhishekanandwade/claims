package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/service"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// StatusHandler handles status and tracking-related HTTP requests
type StatusHandler struct {
	*Base
	claimRepo                 *repo.ClaimRepository
	claimHistoryRepo          *repo.ClaimHistoryRepository
	claimPaymentRepo          *repo.ClaimPaymentRepository
	investigationRepo         *repo.InvestigationRepository
	investigationProgressRepo *repo.InvestigationProgressRepository
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(
	claimRepo *repo.ClaimRepository,
	claimHistoryRepo *repo.ClaimHistoryRepository,
	claimPaymentRepo *repo.ClaimPaymentRepository,
	investigationRepo *repo.InvestigationRepository,
	investigationProgressRepo *repo.InvestigationProgressRepository,
) *StatusHandler {
	base := NewBase("Status & Tracking")
	return &StatusHandler{
		Base:                      base,
		claimRepo:                 claimRepo,
		claimHistoryRepo:          claimHistoryRepo,
		claimPaymentRepo:          claimPaymentRepo,
		investigationRepo:         investigationRepo,
		investigationProgressRepo: investigationProgressRepo,
	}
}

// RegisterRoutes registers all status and tracking routes
func (h *StatusHandler) RegisterRoutes() {
	hRoutes := h.Routes()

	// Claim Status Endpoints
	hRoutes.GET("/claims/:claim_id/status", h.GetClaimStatus).Name("Get Claim Status")
	hRoutes.GET("/claims/:claim_id/sla-countdown", h.GetSLACountdown).Name("Get SLA Countdown")
	hRoutes.GET("/claims/:claim_id/payment-status", h.GetClaimPaymentStatus).Name("Get Claim Payment Status")
	hRoutes.GET("/claims/:claim_id/timeline", h.GetClaimTimeline).Name("Get Claim Timeline")

	// Investigation Progress Status
	hRoutes.GET("/claims/death/:claim_id/investigation/:investigation_id/progress-status", h.GetInvestigationProgressStatus).Name("Get Investigation Progress Status")
}

// ========================================
// CLAIM STATUS ENDPOINTS
// ========================================

// GetClaimStatus retrieves the current status of a claim
// GET /claims/{claim_id}/status
// Reference: BR-CLM-DC-021 (SLA status tracking)
func (h *StatusHandler) GetClaimStatus(sctx *serverRoute.Context, req ClaimIDUri) (*resp.ClaimStatusResponse, error) {
	// Find claim by ID
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("claim not found")
		}
		log.Error(sctx.Ctx, "failed to get claim status: %v", err)
		return nil, fmt.Errorf("failed to get claim status")
	}

	// Build workflow state
	workflowState := resp.WorkflowState{
		CurrentStep: getCurrentStepForStatus(claim.Status),
		Status:      claim.Status,
		UpdatedAt:   claim.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Calculate SLA status
	slaStatus := service.CalculateSLAStatusForClaim(claim)
	workflowState.SLAStatus = slaStatus
	workflowState.SLADeadline = claim.SLADueDate.Format("2006-01-02 15:04:05")

	// Get next step and allowed actions based on current status
	nextStep, allowedActions := getNextStepAndActions(claim.Status, claim.InvestigationRequired, claim.InvestigationStatus)
	workflowState.NextStep = nextStep
	workflowState.AllowedActions = allowedActions

	return resp.NewClaimStatusResponse(claim.ID, claim.Status, workflowState, claim.UpdatedAt), nil
}

// GetSLACountdown retrieves real-time SLA countdown for a claim
// GET /claims/{claim_id}/sla-countdown
// Reference: BR-CLM-DC-003 (SLA without investigation), BR-CLM-DC-004 (SLA with investigation)
func (h *StatusHandler) GetSLACountdown(sctx *serverRoute.Context, req ClaimIDUri) (*resp.SLACountdownResponse, error) {
	// Find claim by ID
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("claim not found")
		}
		log.Error(sctx.Ctx, "failed to get SLA countdown: %v", err)
		return nil, fmt.Errorf("failed to get SLA countdown")
	}

	now := time.Now()
	var slaType string
	var totalDays int
	var deadline time.Time

	// Determine SLA type and calculate days
	if claim.InvestigationRequired && claim.InvestigationStatus == "IN_PROGRESS" {
		// Investigation SLA (21 days)
		slaType = "INVESTIGATION_SLA"
		totalDays = 21
		// TODO: Get investigation start date from investigation table
		deadline = claim.SLADueDate
	} else if claim.ClaimType == "MATURITY" {
		// Maturity claim SLA (7 days)
		slaType = "CLAIM_SLA"
		totalDays = 7
		deadline = claim.SLADueDate
	} else if claim.InvestigationRequired {
		// Death claim with investigation (45 days)
		slaType = "CLAIM_SLA"
		totalDays = 45
		deadline = claim.SLADueDate
	} else {
		// Death claim without investigation (15 days)
		slaType = "CLAIM_SLA"
		totalDays = 15
		deadline = claim.SLADueDate
	}

	// Calculate elapsed and remaining days
	elapsedDays := int(now.Sub(claim.ClaimDate).Hours() / 24)
	remainingDays := totalDays - elapsedDays
	if remainingDays < 0 {
		remainingDays = 0
	}

	// Calculate SLA status (GREEN, YELLOW, RED)
	slaPercentage := float64(elapsedDays) / float64(totalDays) * 100
	slaStatus := service.CalculateSLAStatus(elapsedDays, totalDays)

	// Get allowed actions based on claim status
	_, allowedActions := getNextStepAndActions(claim.Status, claim.InvestigationRequired, claim.InvestigationStatus)

	return resp.NewSLACountdownResponse(slaType, totalDays, elapsedDays, remainingDays, deadline, slaStatus, allowedActions), nil
}

// GetClaimPaymentStatus retrieves payment status for a claim
// GET /claims/{claim_id}/payment-status
// Reference: BR-CLM-DC-010 (Payment disbursement workflow)
func (h *StatusHandler) GetClaimPaymentStatus(sctx *serverRoute.Context, req ClaimIDUri) (*resp.ClaimPaymentStatusResponse, error) {
	// Find claim by ID
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("claim not found")
		}
		log.Error(sctx.Ctx, "failed to get claim payment status: %v", err)
		return nil, fmt.Errorf("failed to get claim payment status")
	}

	// Build payment status response
	paymentData := make(map[string]interface{})

	// Check if claim has payment reference
	if claim.PaymentReference != "" {
		paymentData["payment_reference"] = claim.PaymentReference
	}

	if claim.TransactionID != "" {
		paymentData["transaction_id"] = claim.TransactionID
	}

	if claim.UTRNumber != "" {
		paymentData["utr_number"] = claim.UTRNumber
	}

	if claim.DisbursementDate != nil {
		paymentData["completed_at"] = claim.DisbursementDate
		paymentData["payment_status"] = "COMPLETED"
	} else if claim.Status == "APPROVED" {
		paymentData["payment_status"] = "PENDING"
	} else if claim.Status == "DISBURSED" {
		paymentData["payment_status"] = "DISBURSED"
	} else {
		paymentData["payment_status"] = "NOT_INITIATED"
	}

	// If claim has approved amount, include it
	if claim.ApprovedAmount != nil {
		paymentData["amount"] = *claim.ApprovedAmount
	}

	return resp.NewClaimPaymentStatusResponse(claim.ID, paymentData), nil
}

// GetClaimTimeline retrieves complete timeline of claim events
// GET /claims/{claim_id}/timeline
// Reference: BR-CLM-DC-019 (Communication triggers)
func (h *StatusHandler) GetClaimTimeline(sctx *serverRoute.Context, req ClaimIDUri) (*resp.ClaimTimelineResponse, error) {
	// Get claim history/timeline from database
	histories, err := h.claimHistoryRepo.GetTimeline(sctx.Ctx, req.ClaimID)
	if err != nil {
		log.Error(sctx.Ctx, "failed to get claim timeline: %v", err)
		return nil, fmt.Errorf("failed to get claim timeline")
	}

	// Convert history records to timeline events
	events := make([]resp.TimelineEvent, 0, len(histories))
	for _, history := range histories {
		event := resp.TimelineEvent{
			ID:          history.ID.String(),
			Timestamp:   history.CreatedAt.Format("2006-01-02 15:04:05"),
			EventType:   history.ActionType,
			EntityID:    history.EntityID.String(),
			Description: getEventDescription(history.ActionType, history.EntityType),
			ChangedBy:   history.ChangedBy,
		}

		// Add metadata if available
		if history.Metadata != "" {
			// TODO: Parse JSON metadata into map
		}

		events = append(events, event)
	}

	return resp.NewClaimTimelineResponse(req.ClaimID, events), nil
}

// ========================================
// INVESTIGATION PROGRESS STATUS
// ========================================

// GetInvestigationProgressStatus retrieves investigation progress status
// GET /claims/death/{claim_id}/investigation/{investigation_id}/progress-status
// Reference: BR-CLM-DC-002 (Investigation SLA)
func (h *StatusHandler) GetInvestigationProgressStatus(sctx *serverRoute.Context, req InvestigationIDUri) (*resp.InvestigationProgressStatusResponse, error) {
	// Find investigation by ID
	investigation, err := h.investigationRepo.FindByID(sctx.Ctx, req.InvestigationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("investigation not found")
		}
		log.Error(sctx.Ctx, "failed to get investigation progress status: %v", err)
		return nil, fmt.Errorf("failed to get investigation progress status")
	}

	// Get latest progress update
	progressList, err := h.investigationProgressRepo.GetProgressTimeline(sctx.Ctx, investigation.ID.String(), nil, nil, 1, 1)
	if err != nil {
		log.Error(sctx.Ctx, "failed to get investigation progress: %v", err)
		return nil, fmt.Errorf("failed to get investigation progress")
	}

	var lastHeartbeat time.Time
	var progressPercentage float64
	var completedItems, totalItems int

	if len(progressList) > 0 {
		progress := progressList[0]
		lastHeartbeat = progress.CreatedAt
		progressPercentage = calculateProgressPercentage(progress)
		completedItems = len(progress.ChecklistItemsCompleted)
		// TODO: Get total checklist items from investigation template
		totalItems = 10
	}

	// Calculate SLA status for investigation (21 days)
	investigationSLAStatus := service.CalculateSLAStatus(
		int(time.Since(investigation.AssignedAt).Hours()/24),
		21,
	)

	return resp.NewInvestigationProgressStatusResponse(
		investigation.ID.String(),
		investigation.ClaimID.String(),
		investigation.Status,
		progressPercentage,
		lastHeartbeat,
		investigation.EstimatedCompletionDate,
		completedItems,
		totalItems,
		investigationSLAStatus,
	), nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// getCurrentStepForStatus maps claim status to current workflow step
func getCurrentStepForStatus(status string) string {
	stepMapping := map[string]string{
		"REGISTERED":             "Claim Registration",
		"UNDER_INVESTIGATION":    "Investigation",
		"INVESTIGATION_COMPLETE": "Investigation Review",
		"APPROVAL_PENDING":       "Approval",
		"CALCULATION_PENDING":    "Benefit Calculation",
		"APPROVED":               "Payment Processing",
		"DISBURSED":              "Disbursement Complete",
		"CLOSED":                 "Claim Closed",
		"REJECTED":               "Claim Rejected",
		"CANCELLED":              "Claim Cancelled",
	}

	if step, ok := stepMapping[status]; ok {
		return step
	}
	return "Unknown"
}

// getNextStepAndActions determines the next step and allowed actions based on claim status
func getNextStepAndActions(status string, investigationRequired bool, investigationStatus string) (string, []string) {
	statusMapping := map[string][]string{
		"REGISTERED":             {"Submit Documents", "Request Information"},
		"UNDER_INVESTIGATION":    {"View Investigation", "Submit Report"},
		"INVESTIGATION_COMPLETE": {"Review Investigation", "Approve/Reject Investigation"},
		"APPROVAL_PENDING":       {"Approve", "Reject", "Request Information"},
		"CALCULATION_PENDING":    {"Calculate Benefit", "Approve Calculation"},
		"APPROVED":               {"Initiate Payment", "Validate Bank Account"},
		"DISBURSED":              {"View Payment Details", "Close Claim"},
		"CLOSED":                 {"View Claim Details", "Download Documents"},
		"REJECTED":               {"View Rejection Reason", "Submit Appeal"},
		"CANCELLED":              {"View Cancellation Reason"},
	}

	nextStepMapping := map[string]string{
		"REGISTERED":             "Document Upload",
		"UNDER_INVESTIGATION":    "Investigation Completion",
		"INVESTIGATION_COMPLETE": "Investigation Review",
		"APPROVAL_PENDING":       "Approval Decision",
		"CALCULATION_PENDING":    "Calculation Review",
		"APPROVED":               "Payment Disbursement",
		"DISBURSED":              "Claim Closure",
		"CLOSED":                 "Claim Complete",
		"REJECTED":               "Appeal Period",
		"CANCELLED":              "Claim Cancelled",
	}

	actions, _ := statusMapping[status]
	nextStep, _ := nextStepMapping[status]

	return nextStep, actions
}

// getEventDescription generates human-readable description for timeline events
func getEventDescription(actionType, entityType string) string {
	descriptions := map[string]string{
		"CLAIM_REGISTERED":       "Claim was registered",
		"STATUS_CHANGED":         fmt.Sprintf("%s status was updated", entityType),
		"DOCUMENT_UPLOADED":      "Documents were uploaded",
		"DOCUMENT_VERIFIED":      "Documents were verified",
		"INVESTIGATION_ASSIGNED": "Investigation officer was assigned",
		"INVESTIGATION_STARTED":  "Investigation was initiated",
		"INVESTIGATION_COMPLETE": "Investigation was completed",
		"APPROVAL_GRANTED":       "Claim was approved",
		"APPROVAL_REJECTED":      "Claim was rejected",
		"PAYMENT_INITIATED":      "Payment disbursement was initiated",
		"PAYMENT_COMPLETED":      "Payment was successfully disbursed",
		"CLAIM_CLOSED":           "Claim was closed",
	}

	if desc, ok := descriptions[actionType]; ok {
		return desc
	}
	return fmt.Sprintf("%s: %s", entityType, actionType)
}

// calculateProgressPercentage calculates investigation progress percentage
func calculateProgressPercentage(progress repo.InvestigationProgress) float64 {
	// TODO: Implement actual progress calculation based on checklist items
	// For now, return a placeholder value
	return 50.0
}
