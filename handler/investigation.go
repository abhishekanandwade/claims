package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// InvestigationHandler handles investigation-related HTTP requests
type InvestigationHandler struct {
	*serverHandler.Base
	svc *repo.InvestigationRepository
}

// NewInvestigationHandler creates a new investigation handler
func NewInvestigationHandler(svc *repo.InvestigationRepository) *InvestigationHandler {
	base := serverHandler.New("Investigations").
		SetPrefix("/v1").
		AddPrefix("")
	return &InvestigationHandler{
		Base: base,
		svc:  svc,
	}
}

// TODO: Implement investigation handler routes and methods
// This is a placeholder - will be implemented in Phase 3
