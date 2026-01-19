package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
)

// SurvivalBenefitHandler handles survival benefit-related HTTP requests
type SurvivalBenefitHandler struct {
	*serverHandler.Base
}

// NewSurvivalBenefitHandler creates a new survival benefit handler
func NewSurvivalBenefitHandler() *SurvivalBenefitHandler {
	base := serverHandler.New("SurvivalBenefits").
		SetPrefix("/v1").
		AddPrefix("")
	return &SurvivalBenefitHandler{
		Base: base,
	}
}

// TODO: Implement survival benefit handler routes and methods
// This is a placeholder - will be implemented in Phase 4
