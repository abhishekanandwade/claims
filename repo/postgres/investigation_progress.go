package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// InvestigationProgressRepository handles investigation progress data operations
type InvestigationProgressRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewInvestigationProgressRepository creates a new investigation progress repository
func NewInvestigationProgressRepository(db *dblib.DB, cfg *config.Config) *InvestigationProgressRepository {
	return &InvestigationProgressRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement investigation progress repository methods
// This is a placeholder - will be implemented in Phase 3
