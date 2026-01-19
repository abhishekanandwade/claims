package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// ClaimRepository handles claim data operations
type ClaimRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimRepository creates a new claim repository
func NewClaimRepository(db *dblib.DB, cfg *config.Config) *ClaimRepository {
	return &ClaimRepository{
		db:  db,
		cfg: cfg,
	}
}

const claimTable = "claims"

// Create inserts a new claim
// Reference: BR-CLM-DC-001 (Claim registration), seed/db/claims_database_schema.sql:110-183
func (r *ClaimRepository) Create(ctx context.Context, data domain.Claim) (domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(claimTable).
		Columns(
			"claim_number", "claim_type", "policy_id", "customer_id",
			"claim_date", "death_date", "death_place", "death_type",
			"claimant_name", "claimant_type", "claimant_relation", "claimant_phone", "claimant_email",
			"status", "workflow_state", "investigation_required",
			"claim_amount", "approved_amount", "sum_assured",
			"reversionary_bonus", "terminal_bonus", "outstanding_loan", "unpaid_premiums",
			"sla_due_date", "sla_breached", "sla_breach_days", "sla_status",
			"bank_account_number", "bank_ifsc_code", "bank_account_holder_name", "bank_name", "bank_verified",
			"metadata", "created_by", "updated_by",
		).
		Values(
			data.ClaimNumber, data.ClaimType, data.PolicyID, data.CustomerID,
			data.ClaimDate, data.DeathDate, data.DeathPlace, data.DeathType,
			data.ClaimantName, data.ClaimantType, data.ClaimantRelation, data.ClaimantPhone, data.ClaimantEmail,
			data.Status, data.WorkflowState, data.InvestigationRequired,
			data.ClaimAmount, data.ApprovedAmount, data.SumAssured,
			data.ReversionaryBonus, data.TerminalBonus, data.OutstandingLoan, data.UnpaidPremiums,
			data.SLADueDate, data.SLABreached, data.SLABreachDays, data.SLAStatus,
			data.BankAccountNumber, data.BankIFSCCode, data.BankAccountHolderName, data.BankName, data.BankVerified,
			data.Metadata, data.CreatedBy, data.UpdatedBy,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	return result, err
}

// FindByID retrieves a claim by ID
func (r *ClaimRepository) FindByID(ctx context.Context, claimID string) (domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"id": claimID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByClaimNumber retrieves a claim by claim number
func (r *ClaimRepository) FindByClaimNumber(ctx context.Context, claimNumber string) (domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"claim_number": claimNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return result, err
	}
	return result, nil
}

// List retrieves claims with filters and pagination
func (r *ClaimRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64, orderBy, sortType string) ([]domain.Claim, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query with filters
	baseQuery := sq.Select("*").From(claimTable).PlaceholderFormat(sq.Dollar)

	// Apply filters if provided
	for key, value := range filters {
		baseQuery = baseQuery.Where(sq.Eq{key: value})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").From(claimTable).PlaceholderFormat(sq.Dollar)
	for key, value := range filters {
		countQuery = countQuery.Where(sq.Eq{key: value})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination and sorting
	query := baseQuery.OrderBy(orderBy + " " + sortType).
		Limit(uint64(limit)).
		Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates claim fields
func (r *ClaimRepository) Update(ctx context.Context, claimID string, updates map[string]interface{}) (domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": claimID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	// Apply updates dynamically
	for key, value := range updates {
		query = query.Set(key, value)
	}

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	return result, err
}

// UpdateStatus updates claim status
// Reference: BR-CLM-DC-021 (SLA color coding based on status)
func (r *ClaimRepository) UpdateStatus(ctx context.Context, claimID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": claimID}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// GetApprovalQueue retrieves claims pending approval
// Reference: BR-CLM-DC-005 (Approval workflow)
func (r *ClaimRepository) GetApprovalQueue(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.Claim, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query for approval queue
	baseQuery := sq.Select("*").
		From(claimTable).
		Where(sq.Or{
			sq.Eq{"status": "DOCUMENT_VERIFIED"},
			sq.Eq{"status": "INVESTIGATION_COMPLETED"},
			sq.Eq{"status": "APPROVAL_PENDING"},
		}).
		PlaceholderFormat(sq.Dollar)

	// Apply additional filters
	for key, value := range filters {
		baseQuery = baseQuery.Where(sq.Eq{key: value})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(claimTable).
		Where(sq.Or{
			sq.Eq{"status": "DOCUMENT_VERIFIED"},
			sq.Eq{"status": "INVESTIGATION_COMPLETED"},
			sq.Eq{"status": "APPROVAL_PENDING"},
		}).
		PlaceholderFormat(sq.Dollar)

	for key, value := range filters {
		countQuery = countQuery.Where(sq.Eq{key: value})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination
	query := baseQuery.
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// GetPaymentQueue retrieves claims approved and pending payment
// Reference: BR-CLM-DC-010 (Disbursement workflow)
func (r *ClaimRepository) GetPaymentQueue(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.Claim, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query for payment queue
	baseQuery := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"status": "APPROVED"}).
		Where(sq.Eq{"disbursement_date": nil}).
		PlaceholderFormat(sq.Dollar)

	// Apply additional filters
	for key, value := range filters {
		baseQuery = baseQuery.Where(sq.Eq{key: value})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(claimTable).
		Where(sq.Eq{"status": "APPROVED"}).
		Where(sq.Eq{"disbursement_date": nil}).
		PlaceholderFormat(sq.Dollar)

	for key, value := range filters {
		countQuery = countQuery.Where(sq.Eq{key: value})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination
	query := baseQuery.
		OrderBy("approval_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// FindByPolicyID retrieves claims by policy ID
func (r *ClaimRepository) FindByPolicyID(ctx context.Context, policyID string) ([]domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"policy_id": policyID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindByCustomerID retrieves claims by customer ID
func (r *ClaimRepository) FindByCustomerID(ctx context.Context, customerID string) ([]domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"customer_id": customerID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindByStatus retrieves claims by status
// Reference: BR-CLM-DC-021 (SLA status based filtering)
func (r *ClaimRepository) FindByStatus(ctx context.Context, status string, skip, limit int64) ([]domain.Claim, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"status": status}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	countQuery := sq.Select("COUNT(*)").
		From(claimTable).
		Where(sq.Eq{"status": status}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// UpdateSLAStatus updates claim SLA status
// Reference: BR-CLM-DC-003, BR-CLM-DC-004, BR-CLM-DC-021 (SLA tracking)
func (r *ClaimRepository) UpdateSLAStatus(ctx context.Context, claimID string, slaStatus string, slaBreached bool, breachDays int) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimTable).
		Set("sla_status", slaStatus).
		Set("sla_breached", slaBreached).
		Set("sla_breach_days", breachDays).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": claimID}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// GetOverdueSLAClaims retrieves claims with breached SLA
// Reference: BR-CLM-DC-021 (RED status claims)
func (r *ClaimRepository) GetOverdueSLAClaims(ctx context.Context, skip, limit int64) ([]domain.Claim, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"sla_breached": true}).
		OrderBy("sla_breach_days DESC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	countQuery := sq.Select("COUNT(*)").
		From(claimTable).
		Where(sq.Eq{"sla_breached": true}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// FindClaimsNeedingInvestigation retrieves claims requiring investigation
// Reference: BR-CLM-DC-001 (3-year rule triggering investigation)
func (r *ClaimRepository) FindClaimsNeedingInvestigation(ctx context.Context) ([]domain.Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimTable).
		Where(sq.Eq{"investigation_required": true}).
		Where(sq.NotEq{"investigation_status": "COMPLETED"}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Claim])
	if err != nil {
		return nil, err
	}

	return results, nil
}
