package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// ClaimHandler handles claim-related HTTP requests
type ClaimHandler struct {
	*serverHandler.Base
	svc *repo.ClaimRepository
}

// NewClaimHandler creates a new claim handler
func NewClaimHandler(svc *repo.ClaimRepository) *ClaimHandler {
	base := serverHandler.New("Claims").
		SetPrefix("/v1").
		AddPrefix("")
	return &ClaimHandler{
		Base: base,
		svc:  svc,
	}
}

// TODO: Implement claim handler routes and methods
// This is a placeholder - will be implemented in Phase 2
