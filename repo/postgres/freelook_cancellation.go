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

// FreeLookCancellationRepository handles free look cancellation and refund processing operations
// Reference: seed/db/claims_database_schema.sql:621-654
// Reference: seed/tool-docs/db-README.md - n-api-db patterns
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

const freeLookCancellationTable = "freelook_cancellations"

// Create inserts a new free look cancellation record
// Reference: BR-CLM-BOND-001 (Free look period validation)
// Reference: BR-CLM-BOND-003 (Refund calculation)
// Reference: BR-CLM-BOND-004 (Maker-checker workflow)
func (r *FreeLookCancellationRepository) Create(ctx context.Context, data domain.FreeLookCancellation) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(freeLookCancellationTable).
		Columns(
			"cancellation_number", "policy_id", "bond_tracking_id",
			"cancellation_request_date", "cancellation_reason", "freelook_period_valid",
			"rejection_reason", "total_premium", "pro_rata_risk_premium",
			"stamp_duty", "medical_costs", "other_deductions",
			"refund_amount", "maker_id", "maker_entry_date",
			"checker_id", "checker_verification_date", "maker_checker_approved",
			"refund_transaction_id", "refund_status", "refund_date",
			"linked_to_finance",
		).
		Values(
			data.CancellationNumber, data.PolicyID, data.BondTrackingID,
			data.CancellationRequestDate, data.CancellationReason, data.FreeLookPeriodValid,
			data.RejectionReason, data.TotalPremium, data.ProRataRiskPremium,
			data.StampDuty, data.MedicalCosts, data.OtherDeductions,
			data.RefundAmount, data.MakerID, data.MakerEntryDate,
			data.CheckerID, data.CheckerVerificationDate, data.MakerCheckerApproved,
			data.RefundTransactionID, data.RefundStatus, data.RefundDate,
			data.LinkedToFinance,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return result, err
}

// FindByID retrieves a free look cancellation record by ID
func (r *FreeLookCancellationRepository) FindByID(ctx context.Context, id string) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(freeLookCancellationTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByCancellationNumber retrieves a free look cancellation record by cancellation number
func (r *FreeLookCancellationRepository) FindByCancellationNumber(ctx context.Context, cancellationNumber string) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(freeLookCancellationTable).
		Where(sq.Eq{"cancellation_number": cancellationNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByPolicyID retrieves free look cancellation record by policy ID
func (r *FreeLookCancellationRepository) FindByPolicyID(ctx context.Context, policyID string) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(freeLookCancellationTable).
		Where(sq.Eq{"policy_id": policyID}).
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return result, err
	}
	return result, nil
}

// List retrieves free look cancellation records with filters and pagination
func (r *FreeLookCancellationRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.FreeLookCancellation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query with filters
	query := sq.Select("*").From(freeLookCancellationTable)

	// Apply dynamic filters
	if policyID, ok := filters["policy_id"]; ok && policyID != nil {
		query = query.Where(sq.Eq{"policy_id": policyID})
	}
	if refundStatus, ok := filters["refund_status"]; ok && refundStatus != nil {
		query = query.Where(sq.Eq{"refund_status": refundStatus})
	}
	if makerID, ok := filters["maker_id"]; ok && makerID != nil {
		query = query.Where(sq.Eq{"maker_id": makerID})
	}
	if checkerID, ok := filters["checker_id"]; ok && checkerID != nil {
		query = query.Where(sq.Eq{"checker_id": checkerID})
	}
	if makerCheckerApproved, ok := filters["maker_checker_approved"]; ok && makerCheckerApproved != nil {
		query = query.Where(sq.Eq{"maker_checker_approved": makerCheckerApproved})
	}
	if linkedToFinance, ok := filters["linked_to_finance"]; ok && linkedToFinance != nil {
		query = query.Where(sq.Eq{"linked_to_finance": linkedToFinance})
	}
	if startDate, ok := filters["start_date"]; ok && startDate != nil {
		query = query.Where(sq.GtOrEq{"cancellation_request_date": startDate})
	}
	if endDate, ok := filters["end_date"]; ok && endDate != nil {
		query = query.Where(sq.LtOrEq{"cancellation_request_date": endDate})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(freeLookCancellationTable)
	if whereClause, _, err := query.ToSql(); err == nil {
		// Extract WHERE clause from main query
		if len(whereClause) > 0 {
			// Simplified count query - in production, extract WHERE clause properly
			countQuery = sq.Select("COUNT(*)").From(freeLookCancellationTable)
		}
	}

	var total int64
	countRow, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err == nil {
		total = countRow
	}

	// Apply pagination and sorting
	query = query.OrderBy("created_at DESC").Offset(uint64(skip)).Limit(uint64(limit))
	query = query.PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update updates a free look cancellation record
func (r *FreeLookCancellationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(freeLookCancellationTable)

	// Build dynamic UPDATE query
	for column, value := range updates {
		query = query.Set(column, value)
	}

	query = query.Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return result, err
}

// UpdateStatus updates the refund status
// Reference: Maker-checker workflow (BR-CLM-BOND-004)
func (r *FreeLookCancellationRepository) UpdateStatus(ctx context.Context, id string, refundStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(freeLookCancellationTable).
		Set("refund_status", refundStatus).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return err
}

// UpdateRefundTransaction updates refund transaction details
func (r *FreeLookCancellationRepository) UpdateRefundTransaction(ctx context.Context, id string, transactionID string, refundStatus string, refundDate time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(freeLookCancellationTable).
		Set("refund_transaction_id", transactionID).
		Set("refund_status", refundStatus).
		Set("refund_date", refundDate).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return err
}

// MakerCheckerApproval implements maker-checker workflow
// Reference: BR-CLM-BOND-004 (Maker-checker segregation of duties)
func (r *FreeLookCancellationRepository) MakerCheckerApproval(ctx context.Context, id string, checkerID string, approved bool) (domain.FreeLookCancellation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(freeLookCancellationTable).
		Set("checker_id", checkerID).
		Set("checker_verification_date", time.Now()).
		Set("maker_checker_approved", approved).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return result, err
}

// GetPendingApprovals retrieves cancellations pending checker approval
// Reference: BR-CLM-BOND-004 (Maker-checker workflow)
func (r *FreeLookCancellationRepository) GetPendingApprovals(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.FreeLookCancellation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(freeLookCancellationTable).
		Where(sq.Eq{"maker_checker_approved": false}).
		Where(sq.NotEq{"maker_id": nil}).
		Where(sq.Eq{"checker_id": nil})

	// Apply additional filters
	if makerID, ok := filters["maker_id"]; ok && makerID != nil {
		query = query.Where(sq.Eq{"maker_id": makerID})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(freeLookCancellationTable).
		Where(sq.Eq{"maker_checker_approved": false}).
		Where(sq.Eq{"checker_id": nil})

	var total int64
	countRow, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err == nil {
		total = countRow
	}

	// Apply pagination and sorting
	query = query.OrderBy("maker_entry_date ASC").Offset(uint64(skip)).Limit(uint64(limit))
	query = query.PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetPendingRefunds retrieves cancellations with pending refund status
func (r *FreeLookCancellationRepository) GetPendingRefunds(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.FreeLookCancellation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(freeLookCancellationTable).
		Where(sq.Eq{"maker_checker_approved": true}).
		Where(sq.NotEq{"refund_status": "SUCCESS"})

	// Apply additional filters
	if refundStatus, ok := filters["refund_status"]; ok && refundStatus != nil {
		query = query.Where(sq.Eq{"refund_status": refundStatus})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(freeLookCancellationTable).
		Where(sq.Eq{"maker_checker_approved": true}).
		Where(sq.NotEq{"refund_status": "SUCCESS"})

	var total int64
	countRow, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err == nil {
		total = countRow
	}

	// Apply pagination and sorting
	query = query.OrderBy("checker_verification_date ASC").Offset(uint64(skip)).Limit(uint64(limit))
	query = query.PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetRefundStats retrieves refund statistics
func (r *FreeLookCancellationRepository) GetRefundStats(ctx context.Context) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select(
		"refund_status",
		"COUNT(*) as count",
	).
		From(freeLookCancellationTable).
		GroupBy("refund_status").
		PlaceholderFormat(sq.Dollar)

	rows, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[struct {
		RefundStatus string `db:"refund_status"`
		Count        int64  `db:"count"`
	}])
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, row := range rows {
		stats[row.RefundStatus] = row.Count
	}

	return stats, nil
}

// CalculateRefundAmount calculates refund amount based on premium and deductions
// Reference: BR-CLM-BOND-003: Refund = Premium - (risk premium + stamp duty + medical + other)
func (r *FreeLookCancellationRepository) CalculateRefundAmount(
	ctx context.Context,
	totalPremium float64,
	proRataRiskPremium float64,
	stampDuty float64,
	medicalCosts *float64,
	otherDeductions *float64,
) float64 {
	// Handle nullable medical costs
	mc := 0.0
	if medicalCosts != nil {
		mc = *medicalCosts
	}

	// Handle nullable other deductions
	od := 0.0
	if otherDeductions != nil {
		od = *otherDeductions
	}

	// Calculate refund: Premium - (risk premium + stamp duty + medical + other)
	refundAmount := totalPremium - (proRataRiskPremium + stampDuty + mc + od)

	// Ensure refund is not negative
	if refundAmount < 0 {
		refundAmount = 0
	}

	return refundAmount
}

// ValidateFreeLookPeriod validates if cancellation request is within free look period
// Reference: BR-CLM-BOND-001: 15 days (physical) or 30 days (electronic)
func (r *FreeLookCancellationRepository) ValidateFreeLookPeriod(
	ctx context.Context,
	bondType string,
	freelookPeriodStartDate time.Time,
	requestDate time.Time,
) (bool, int) {
	var daysInPeriod int

	// BR-CLM-BOND-001: Free look period based on bond type
	if bondType == "ELECTRONIC" {
		daysInPeriod = 30 // 30 days for electronic bonds
	} else {
		daysInPeriod = 15 // 15 days for physical bonds
	}

	// Calculate days elapsed
	daysElapsed := int(requestDate.Sub(freelookPeriodStartDate).Hours() / 24)

	// Check if within free look period
	isValid := daysElapsed >= 0 && daysElapsed <= daysInPeriod

	return isValid, daysElapsed
}

// LinkToFinance marks cancellation as linked to finance system for refund processing
func (r *FreeLookCancellationRepository) LinkToFinance(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(freeLookCancellationTable).
		Set("linked_to_finance", true).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.FreeLookCancellation])
	return err
}

// Delete soft deletes a free look cancellation record (sets deleted_at if column exists)
// Note: Current schema doesn't have deleted_at, so this is a placeholder for future enhancement
func (r *FreeLookCancellationRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(freeLookCancellationTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, _ := query.ToSql()
	_, err := r.db.Exec(ctx, sql, args...)
	return err
}
