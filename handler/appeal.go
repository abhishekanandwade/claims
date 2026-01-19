package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// AppealHandler handles appeal-related HTTP requests
type AppealHandler struct {
	*serverHandler.Base
	svc *repo.AppealRepository
}

// NewAppealHandler creates a new appeal handler
func NewAppealHandler(svc *repo.AppealRepository) *AppealHandler {
	base := serverHandler.New("Appeals").
		SetPrefix("/v1").
		AddPrefix("")
	return &AppealHandler{
		Base: base,
		svc:  svc,
	}
}

// TODO: Implement appeal handler routes and methods
// This is a placeholder - will be implemented in Phase 6
