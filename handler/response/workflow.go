package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	"time"
)

// ========================================
// WORKFLOW MANAGEMENT - RESPONSE DTOS
// ========================================

// WorkflowDetailsResponse represents the response for workflow details
// GET /workflows/death-claim/{workflow_id}
// Reference: INT-CLM-020 (Temporal integration)
type WorkflowDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	WorkflowID                string                 `json:"workflow_id"`
	WorkflowType              string                 `json:"workflow_type"` // e.g., "death_claim", "maturity_claim", "investigation"
	Status                    string                 `json:"status"`         // e.g., "RUNNING", "COMPLETED", "FAILED", "CANCELLED"
	StartedAt                 string                 `json:"started_at"`     // RFC3339 format
	CompletedAt               *string                `json:"completed_at,omitempty"` // RFC3339 format
	ClaimID                   *string                `json:"claim_id,omitempty"`
	History                   []WorkflowEvent        `json:"history,omitempty"`
	CurrentActivity           *string                `json:"current_activity,omitempty"`
	Metadata                  map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowEvent represents a single event in the workflow history
type WorkflowEvent struct {
	EventID   string                 `json:"event_id"`
	EventName string                 `json:"event_name"` // e.g., "activity_scheduled", "activity_completed", "workflow_signaled"
	Timestamp string                 `json:"timestamp"` // RFC3339 format
	Details   map[string]interface{} `json:"details,omitempty"`
}

// WorkflowSignalResponse represents the response for sending a signal to a workflow
// POST /workflows/{workflow_id}/signal
// Reference: INT-CLM-020 (Temporal integration)
type WorkflowSignalResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	WorkflowID                string `json:"workflow_id"`
	SignalName                string `json:"signal_name"`
	SignalSentAt              string `json:"signal_sent_at"` // RFC3339 format
	SignalID                  string `json:"signal_id"`      // Unique identifier for the signal
}

// NewWorkflowDetailsResponse creates a new WorkflowDetailsResponse from workflow details
func NewWorkflowDetailsResponse(
	workflowID string,
	workflowType string,
	status string,
	startedAt time.Time,
	completedAt *time.Time,
	claimID *string,
	history []WorkflowEvent,
	currentActivity *string,
	metadata map[string]interface{},
) *WorkflowDetailsResponse {
	resp := &WorkflowDetailsResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Workflow details retrieved successfully",
		},
		WorkflowID:      workflowID,
		WorkflowType:    workflowType,
		Status:          status,
		StartedAt:       startedAt.Format("2006-01-02 15:04:05"),
		ClaimID:         claimID,
		History:         history,
		CurrentActivity: currentActivity,
		Metadata:        metadata,
	}

	if completedAt != nil {
		formattedCompletedAt := completedAt.Format("2006-01-02 15:04:05")
		resp.CompletedAt = &formattedCompletedAt
	}

	return resp
}

// NewWorkflowSignalResponse creates a new WorkflowSignalResponse after sending a signal
func NewWorkflowSignalResponse(workflowID string, signalName string, signalID string) *WorkflowSignalResponse {
	return &WorkflowSignalResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Signal sent to workflow successfully",
		},
		WorkflowID:   workflowID,
		SignalName:   signalName,
		SignalSentAt: time.Now().Format("2006-01-02 15:04:05"),
		SignalID:     signalID,
	}
}
