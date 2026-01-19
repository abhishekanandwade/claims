package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// AMLHandler handles AML-related HTTP requests
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

// TODO: Implement AML handler routes and methods
// This is a placeholder - will be implemented in Phase 5
