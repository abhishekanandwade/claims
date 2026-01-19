package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
)

// BankingHandler handles banking-related HTTP requests
type BankingHandler struct {
	*serverHandler.Base
}

// NewBankingHandler creates a new banking handler
func NewBankingHandler() *BankingHandler {
	base := serverHandler.New("Banking").
		SetPrefix("/v1").
		AddPrefix("")
	return &BankingHandler{
		Base: base,
	}
}

// TODO: Implement banking handler routes and methods
// This is a placeholder - will be implemented in Phase 5
