package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ClaimStatusResponse represents the status of a claim
type ClaimStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string                `json:"claim_id"`
	Status                    string                `json:"status"`
	WorkflowState             WorkflowStateResponse `json:"workflow_state"`
	UpdatedAt                 string                `json:"updated_at"`
}

// SLACountdownResponse represents SLA countdown information
type SLACountdownResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	SLAType                   string   `json:"sla_type"` // CLAIM_SLA, INVESTIGATION_SLA, APPEAL_SLA
	TotalDays                 int      `json:"total_days"`
	ElapsedDays               int      `json:"elapsed_days"`
	RemainingDays             int      `json:"remaining_days"`
	Deadline                  string   `json:"deadline"`
	DaysRemaining             int      `json:"days_remaining"`
	SLAStatus                 string   `json:"sla_status"` // GREEN, YELLOW, RED
	AllowedActions            []string `json:"allowed_actions"`
}

// ClaimPaymentStatusResponse represents payment status for a claim
type ClaimPaymentStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string  `json:"claim_id"`
	PaymentID                 string  `json:"payment_id,omitempty"`
	PaymentStatus             string  `json:"payment_status,omitempty"`
	PaymentReference          string  `json:"payment_reference,omitempty"`
	TransactionID             string  `json:"transaction_id,omitempty"`
	UTRNumber                 string  `json:"utr_number,omitempty"`
	InitiatedAt               string  `json:"initiated_at,omitempty"`
	CompletedAt               string  `json:"completed_at,omitempty"`
	Amount                    float64 `json:"amount,omitempty"`
}

// ClaimTimelineResponse represents the complete claim timeline
type ClaimTimelineResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ClaimID                   string          `json:"claim_id"`
	TimelineEvents            []TimelineEvent `json:"timeline_events"`
	TotalEvents               int             `json:"total_events"`
}

// TimelineEvent represents a single event in the claim timeline
type TimelineEvent struct {
	ID            string                 `json:"id"`
	Timestamp     string                 `json:"timestamp"`
	EventType     string                 `json:"event_type"` // CREATED, UPDATED, STATUS_CHANGED, DOCUMENT_UPLOADED, etc.
	EntityID      string                 `json:"entity_id,omitempty"`
	EntityName    string                 `json:"entity_name,omitempty"`
	Description   string                 `json:"description"`
	ChangedFields []string               `json:"changed_fields,omitempty"`
	OldValue      string                 `json:"old_value,omitempty"`
	NewValue      string                 `json:"new_value,omitempty"`
	ChangedBy     string                 `json:"changed_by,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// InvestigationProgressStatusResponse represents investigation progress status
type InvestigationProgressStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	InvestigationID           string  `json:"investigation_id"`
	ClaimID                   string  `json:"claim_id"`
	Status                    string  `json:"status"`
	ProgressPercentage        float64 `json:"progress_percentage"`
	LastHeartbeat             string  `json:"last_heartbeat"`
	EstimatedCompletionDate   string  `json:"estimated_completion_date,omitempty"`
	CompletedChecklistItems   int     `json:"completed_checklist_items"`
	TotalChecklistItems       int     `json:"total_checklist_items"`
	SLAStatus                 string  `json:"sla_status"` // GREEN, YELLOW, RED
}

// NotificationDeliveryStatusResponse represents notification delivery status
type NotificationDeliveryStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	NotificationID            string                   `json:"notification_id"`
	Status                    string                   `json:"status"` // SENT, DELIVERED, FAILED, PENDING
	DeliveryAttempts          int                      `json:"delivery_attempts"`
	LastAttemptAt             string                   `json:"last_attempt_at,omitempty"`
	DeliveryStatus            map[string]ChannelStatus `json:"delivery_status"` // SMS, Email, WhatsApp
	CreatedAt                 string                   `json:"created_at"`
}

// EntityStatusResponse represents generic entity status
type EntityStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	EntityID                  string `json:"entity_id"`
	EntityType                string `json:"entity_type"` // CLAIM, DOCUMENT, INVESTIGATION, PAYMENT, APPEAL
	Status                    string `json:"status"`
	UpdatedAt                 string `json:"updated_at"`
	LastActivity              string `json:"last_activity,omitempty"`
}

// BulkStatusResponse represents status of multiple entities
type BulkStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Statuses                  []EntityStatusResponse `json:"statuses"`
	TotalEntities             int                    `json:"total_entities"`
	Successful                int                    `json:"successful"`
	Failed                    int                    `json:"failed"`
}

// Helper functions

// NewClaimStatusResponse creates a new claim status response
func NewClaimStatusResponse(claimID, status string, workflowState WorkflowStateResponse, updatedAt time.Time) *ClaimStatusResponse {
	return &ClaimStatusResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: "200",
			Success:    true,
			Message:    "Claim status retrieved successfully",
		},
		ClaimID:       claimID,
		Status:        status,
		WorkflowState: workflowState,
		UpdatedAt:     updatedAt.Format("2006-01-02 15:04:05"),
	}
}

// NewSLACountdownResponse creates a new SLA countdown response
func NewSLACountdown(slaType string, totalDays, elapsedDays, remainingDays int, deadline time.Time, slaStatus string, allowedActions []string) *SLACountdownResponse {
	return &SLACountdownResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: "200",
			Success:    true,
			Message:    "SLA countdown retrieved successfully",
		},
		SLAType:        slaType,
		TotalDays:      totalDays,
		ElapsedDays:    elapsedDays,
		RemainingDays:  remainingDays,
		Deadline:       deadline.Format("2006-01-02 15:04:05"),
		DaysRemaining:  remainingDays,
		SLAStatus:      slaStatus,
		AllowedActions: allowedActions,
	}
}

// NewClaimPaymentStatusResponse creates a new claim payment status response
func NewClaimPaymentStatusResponse(claimID string, paymentData map[string]interface{}) *ClaimPaymentStatusResponse {
	resp := &ClaimPaymentStatusResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: "200",
			Success:    true,
			Message:    "Payment status retrieved successfully",
		},
		ClaimID: claimID,
	}

	if paymentID, ok := paymentData["payment_id"].(string); ok {
		resp.PaymentID = paymentID
	}
	if paymentStatus, ok := paymentData["payment_status"].(string); ok {
		resp.PaymentStatus = paymentStatus
	}
	if paymentRef, ok := paymentData["payment_reference"].(string); ok {
		resp.PaymentReference = paymentRef
	}
	if txnID, ok := paymentData["transaction_id"].(string); ok {
		resp.TransactionID = txnID
	}
	if utr, ok := paymentData["utr_number"].(string); ok {
		resp.UTRNumber = utr
	}
	if initiatedAt, ok := paymentData["initiated_at"].(time.Time); ok {
		resp.InitiatedAt = initiatedAt.Format("2006-01-02 15:04:05")
	}
	if completedAt, ok := paymentData["completed_at"].(time.Time); ok {
		resp.CompletedAt = completedAt.Format("2006-01-02 15:04:05")
	}
	if amount, ok := paymentData["amount"].(float64); ok {
		resp.Amount = amount
	}

	return resp
}

// NewClaimTimelineResponse creates a new claim timeline response
func NewClaimTimelineResponse(claimID string, events []TimelineEvent) *ClaimTimelineResponse {
	return &ClaimTimelineResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: "200",
			Success:    true,
			Message:    "Claim timeline retrieved successfully",
		},
		ClaimID:        claimID,
		TimelineEvents: events,
		TotalEvents:    len(events),
	}
}

// NewInvestigationProgressStatusResponse creates a new investigation progress status response
func NewInvestigationProgressStatusResponse(investigationID, claimID, status string, progressPercentage float64, lastHeartbeat time.Time, estimatedCompletion *time.Time, completedItems, totalItems int, slaStatus string) *InvestigationProgressStatusResponse {
	resp := &InvestigationProgressStatusResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: "200",
			Success:    true,
			Message:    "Investigation progress status retrieved successfully",
		},
		InvestigationID:         investigationID,
		ClaimID:                 claimID,
		Status:                  status,
		ProgressPercentage:      progressPercentage,
		LastHeartbeat:           lastHeartbeat.Format("2006-01-02 15:04:05"),
		CompletedChecklistItems: completedItems,
		TotalChecklistItems:     totalItems,
		SLAStatus:               slaStatus,
	}

	if estimatedCompletion != nil {
		resp.EstimatedCompletionDate = estimatedCompletion.Format("2006-01-02 15:04:05")
	}

	return resp
}

// CalculateTimelineProgress calculates progress percentage from timeline events
func CalculateTimelineProgress(totalSteps, completedSteps int) float64 {
	if totalSteps == 0 {
		return 0.0
	}
	return float64(completedSteps) / float64(totalSteps) * 100
}
