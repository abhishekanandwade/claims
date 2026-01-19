package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// ClaimHistoryRepository handles claim history data operations
type ClaimHistoryRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimHistoryRepository creates a new claim history repository
func NewClaimHistoryRepository(db *dblib.DB, cfg *config.Config) *ClaimHistoryRepository {
	return &ClaimHistoryRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement claim history repository methods
