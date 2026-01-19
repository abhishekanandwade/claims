package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// OmbudsmanHandler handles ombudsman-related HTTP requests
type OmbudsmanHandler struct {
	*serverHandler.Base
	svc *repo.OmbudsmanComplaintRepository
}

// NewOmbudsmanHandler creates a new ombudsman handler
func NewOmbudsmanHandler(svc *repo.OmbudsmanComplaintRepository) *OmbudsmanHandler {
	base := serverHandler.New("Ombudsman").
		SetPrefix("/v1").
		AddPrefix("")
	return &OmbudsmanHandler{
		Base: base,
		svc:  svc,
	}
}

// TODO: Implement ombudsman handler routes and methods
// This is a placeholder - will be implemented in Phase 7
