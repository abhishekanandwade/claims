package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// FreeLookCancellationRepository handles free look cancellation data operations
type FreeLookCancellationRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewFreeLookCancellationRepository creates a new free look cancellation repository
func NewFreeLookCancellationRepository(db *dblib.DB, cfg *config.Config) *FreeLookCancellationRepository {
	return &FreeLookCancellationRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement free look cancellation repository methods
