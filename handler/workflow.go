package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
)

// WorkflowHandler handles workflow-related HTTP requests
type WorkflowHandler struct {
	*serverHandler.Base
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler() *WorkflowHandler {
	base := serverHandler.New("Workflow").
		SetPrefix("/v1").
		AddPrefix("")
	return &WorkflowHandler{
		Base: base,
	}
}

// TODO: Implement workflow handler routes and methods
// This is a placeholder - will be implemented in Phase 8
