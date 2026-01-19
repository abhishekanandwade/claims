package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// ClaimCommunicationRepository handles claim communication data operations
type ClaimCommunicationRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimCommunicationRepository creates a new claim communication repository
func NewClaimCommunicationRepository(db *dblib.DB, cfg *config.Config) *ClaimCommunicationRepository {
	return &ClaimCommunicationRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement claim communication repository methods
