package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// AppealRepository handles appeal data operations
type AppealRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewAppealRepository creates a new appeal repository
func NewAppealRepository(db *dblib.DB, cfg *config.Config) *AppealRepository {
	return &AppealRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement appeal repository methods
