package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// AMLAlertRepository handles AML alert data operations
type AMLAlertRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewAMLAlertRepository creates a new AML alert repository
func NewAMLAlertRepository(db *dblib.DB, cfg *config.Config) *AMLAlertRepository {
	return &AMLAlertRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement AML alert repository methods
