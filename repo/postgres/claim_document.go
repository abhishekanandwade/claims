package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// ClaimDocumentRepository handles claim document data operations
type ClaimDocumentRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimDocumentRepository creates a new claim document repository
func NewClaimDocumentRepository(db *dblib.DB, cfg *config.Config) *ClaimDocumentRepository {
	return &ClaimDocumentRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement claim document repository methods
// This is a placeholder - will be implemented in Phase 2
