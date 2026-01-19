package repo

import (

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// ClaimPaymentRepository handles claim payment data operations
type ClaimPaymentRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimPaymentRepository creates a new claim payment repository
func NewClaimPaymentRepository(db *dblib.DB, cfg *config.Config) *ClaimPaymentRepository {
	return &ClaimPaymentRepository{
		db:  db,
		cfg: cfg,
	}
}

// TODO: Implement claim payment repository methods
