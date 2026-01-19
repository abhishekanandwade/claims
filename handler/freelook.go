package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// FreeLookHandler handles free look-related HTTP requests
type FreeLookHandler struct {
	*serverHandler.Base
	svc *repo.FreeLookCancellationRepository
}

// NewFreeLookHandler creates a new free look handler
func NewFreeLookHandler(svc *repo.FreeLookCancellationRepository) *FreeLookHandler {
	base := serverHandler.New("FreeLook").
		SetPrefix("/v1").
		AddPrefix("")
	return &FreeLookHandler{
		Base: base,
		svc:  svc,
	}
}

// TODO: Implement free look handler routes and methods
// This is a placeholder - will be implemented in Phase 6
