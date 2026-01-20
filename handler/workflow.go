package handler

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
)

// WorkflowHandler handles Temporal workflow-related HTTP requests
// Reference: INT-CLM-020 (Temporal integration)
type WorkflowHandler struct {
	*serverHandler.Base
	// TODO: Add Temporal client dependency when Temporal integration is implemented
	// temporalClient client.Client
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler() *WorkflowHandler {
	base := serverHandler.New("Workflows").
		SetPrefix("/v1").
		AddPrefix("")
	return &WorkflowHandler{
		Base: base,
	}
}

// Routes defines all routes for this handler
func (h *WorkflowHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Workflow Management (2 endpoints)
		serverRoute.GET("/workflows/death-claim/:workflow_id", h.GetWorkflowDetails).Name("Get Workflow Details"),
		serverRoute.POST("/workflows/:workflow_id/signal", h.SignalWorkflow).Name("Signal Workflow"),
	}
}

// GetWorkflowDetails retrieves the details of a Temporal workflow instance
// GET /workflows/death-claim/{workflow_id}
// Reference: INT-CLM-020 (Temporal integration)
func (h *WorkflowHandler) GetWorkflowDetails(sctx *serverRoute.Context, uri WorkflowIDUri) (*resp.WorkflowDetailsResponse, error) {
	log.Info(sctx.Ctx, "Fetching workflow details for workflow_id: %s", uri.WorkflowID)

	// TODO: Integrate with Temporal to fetch actual workflow details
	// For now, returning a placeholder response
	// Actual implementation will:
	// 1. Query Temporal workflow execution using temporalClient.DescribeWorkflowExecution()
	// 2. Fetch workflow history using temporalClient.GetWorkflowHistory()
	// 3. Map Temporal workflow state to response DTO

	// Placeholder workflow details
	workflowType := "death_claim"
	status := "RUNNING"
	startedAt := time.Now().Add(-24 * time.Hour) // Started 24 hours ago
	var completedAt *time.Time = nil             // Still running
	claimID := "CLM20260001"
	currentActivity := "waiting_for_investigation_report"

	// Placeholder workflow history
	history := []resp.WorkflowEvent{
		{
			EventID:   "evt_001",
			EventName: "workflow_started",
			Timestamp: startedAt.Format("2006-01-02 15:04:05"),
			Details: map[string]interface{}{
				"triggered_by": "claim_registration",
				"claim_id":     claimID,
			},
		},
		{
			EventID:   "evt_002",
			EventName: "activity_scheduled",
			Timestamp: startedAt.Add(1 * time.Minute).Format("2006-01-02 15:04:05"),
			Details: map[string]interface{}{
				"activity": "assign_investigation_officer",
			},
		},
		{
			EventID:   "evt_003",
			EventName: "activity_completed",
			Timestamp: startedAt.Add(2 * time.Hour).Format("2006-01-02 15:04:05"),
			Details: map[string]interface{}{
				"activity": "assign_investigation_officer",
				"result":   "investigation_assigned",
			},
		},
	}

	// Placeholder metadata
	metadata := map[string]interface{}{
		"claim_type":        "DEATH",
		"investigation_id":  "INV20260001",
		"sla_deadline":      startedAt.Add(45 * 24 * time.Hour).Format("2006-01-02 15:04:05"), // 45 days from start
		"current_step":      "INVESTIGATION",
		"next_step":         "APPROVAL",
		"workflow_version":  "v1.0.0",
		"temporal_task_queue": "death-claim-task-queue",
	}

	// Build response
	r := resp.NewWorkflowDetailsResponse(
		uri.WorkflowID,
		workflowType,
		status,
		startedAt,
		completedAt,
		&claimID,
		history,
		&currentActivity,
		metadata,
	)

	log.Info(sctx.Ctx, "Workflow details retrieved for workflow_id: %s, status: %s", uri.WorkflowID, status)
	return r, nil
}

// SignalWorkflow sends a signal to a running Temporal workflow instance
// POST /workflows/{workflow_id}/signal
// Reference: INT-CLM-020 (Temporal integration)
func (h *WorkflowHandler) SignalWorkflow(sctx *serverRoute.Context, uri WorkflowIDUri, req SignalWorkflowRequest) (*resp.WorkflowSignalResponse, error) {
	log.Info(sctx.Ctx, "Sending signal '%s' to workflow_id: %s", req.SignalName, uri.WorkflowID)

	// TODO: Integrate with Temporal to send actual signal
	// For now, returning a placeholder response
	// Actual implementation will:
	// 1. Validate signal name against allowed signals for the workflow type
	// 2. Validate signal data structure (schema validation)
	// 3. Send signal using temporalClient.SignalWorkflow()
	// 4. Return signal ID for tracking

	// Validate signal name (basic validation)
	allowedSignals := map[string]bool{
		"investigation_complete":       true,
		"document_uploaded":           true,
		"approval_received":           true,
		"rejection_received":          true,
		"payment_processed":           true,
		"investigation_escalated":     true,
		"document_reminder_sent":      true,
		"sla_breach_warning":          true,
		"cancel_workflow":             true,
		"manual_override_requested":   true,
	}

	if !allowedSignals[req.SignalName] {
		log.Error(sctx.Ctx, "Invalid signal name: %s", req.SignalName)
		return nil, &port.DomainError{
			StatusCode: 400,
			Success:    false,
			Message:    "Invalid signal name",
			ErrorCode:  "INVALID_SIGNAL",
		}
	}

	// Generate unique signal ID
	signalIDBytes := make([]byte, 16)
	if _, err := rand.Read(signalIDBytes); err != nil {
		log.Error(sctx.Ctx, "Error generating signal ID: %v", err)
		return nil, err
	}
	signalID := "SIG" + hex.EncodeToString(signalIDBytes)

	// Log signal data for debugging
	log.Info(sctx.Ctx, "Signal data: %+v", req.SignalData)

	// Build response
	r := resp.NewWorkflowSignalResponse(
		uri.WorkflowID,
		req.SignalName,
		signalID,
	)

	log.Info(sctx.Ctx, "Signal sent successfully. workflow_id: %s, signal_name: %s, signal_id: %s", uri.WorkflowID, req.SignalName, signalID)
	return r, nil
}
