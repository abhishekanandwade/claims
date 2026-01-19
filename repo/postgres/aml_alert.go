package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// AMLAlertRepository handles AML/CFT alert data operations
// Reference: E-CLM-AML-001, seed/db/claims_database_schema.sql:335-371
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

const amlAlertTable = "aml_alerts"

// Create inserts a new AML alert
// Reference: BR-CLM-AML-001 to 005 (AML trigger rules)
func (r *AMLAlertRepository) Create(ctx context.Context, data domain.AMLAlert) (domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(amlAlertTable).
		Columns(
			"id", "alert_id", "trigger_code", "policy_id", "customer_id",
			"transaction_type", "transaction_amount", "transaction_date", "payment_mode",
			"risk_level", "risk_score", "alert_status", "alert_description", "trigger_details",
			"transaction_blocked",
			"filing_required", "filing_type", "filing_status", "filing_reference", "filed_at", "filed_by",
			"pan_number", "pan_verified", "nominee_change_detected",
			"created_at", "updated_at",
		).
		Values(
			data.ID, data.AlertID, data.TriggerCode, data.PolicyID, data.CustomerID,
			data.TransactionType, data.TransactionAmount, data.TransactionDate, data.PaymentMode,
			data.RiskLevel, data.RiskScore, data.AlertStatus, data.AlertDescription, data.TriggerDetails,
			data.TransactionBlocked,
			data.FilingRequired, data.FilingType, data.FilingStatus, data.FilingReference, data.FiledAt, data.FiledBy,
			data.PANNumber, data.PANVerified, data.NomineeChangeDetected,
			data.CreatedAt, data.UpdatedAt,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return result, err
}

// FindByID retrieves an AML alert by ID
func (r *AMLAlertRepository) FindByID(ctx context.Context, alertID string) (domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"id": alertID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByAlertID retrieves an AML alert by alert_id (business key)
func (r *AMLAlertRepository) FindByAlertID(ctx context.Context, alertID string) (domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"alert_id": alertID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByPolicyID retrieves all AML alerts for a policy
func (r *AMLAlertRepository) FindByPolicyID(ctx context.Context, policyID string) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"policy_id": policyID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// FindByCustomerID retrieves all AML alerts for a customer
func (r *AMLAlertRepository) FindByCustomerID(ctx context.Context, customerID string) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"customer_id": customerID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// List retrieves AML alerts with pagination and filtering
// Reference: BR-CLM-AML-001 to 005 (AML alert monitoring)
func (r *AMLAlertRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit uint64, orderBy, sortType string) ([]domain.AMLAlert, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build count query
	countQuery := sq.Select("COUNT(*)").
		From(amlAlertTable).
		PlaceholderFormat(sq.Dollar)

	// Build data query
	query := sq.Select("*").
		From(amlAlertTable).
		PlaceholderFormat(sq.Dollar)

	// Apply filters
	if riskLevel, ok := filters["risk_level"]; ok && riskLevel != nil {
		countQuery = countQuery.Where(sq.Eq{"risk_level": riskLevel})
		query = query.Where(sq.Eq{"risk_level": riskLevel})
	}

	if alertStatus, ok := filters["alert_status"]; ok && alertStatus != nil {
		countQuery = countQuery.Where(sq.Eq{"alert_status": alertStatus})
		query = query.Where(sq.Eq{"alert_status": alertStatus})
	}

	if triggerCode, ok := filters["trigger_code"]; ok && triggerCode != nil {
		countQuery = countQuery.Where(sq.Eq{"trigger_code": triggerCode})
		query = query.Where(sq.Eq{"trigger_code": triggerCode})
	}

	if policyID, ok := filters["policy_id"]; ok && policyID != nil {
		countQuery = countQuery.Where(sq.Eq{"policy_id": policyID})
		query = query.Where(sq.Eq{"policy_id": policyID})
	}

	if filingRequired, ok := filters["filing_required"]; ok && filingRequired != nil {
		countQuery = countQuery.Where(sq.Eq{"filing_required": filingRequired})
		query = query.Where(sq.Eq{"filing_required": filingRequired})
	}

	if transactionBlocked, ok := filters["transaction_blocked"]; ok && transactionBlocked != nil {
		countQuery = countQuery.Where(sq.Eq{"transaction_blocked": transactionBlocked})
		query = query.Where(sq.Eq{"transaction_blocked": transactionBlocked})
	}

	// Date range filter for transaction_date
	if startDate, ok := filters["start_date"]; ok && startDate != nil {
		countQuery = countQuery.Where(sq.GtOrEq{"transaction_date": startDate})
		query = query.Where(sq.GtOrEq{"transaction_date": startDate})
	}

	if endDate, ok := filters["end_date"]; ok && endDate != nil {
		countQuery = countQuery.Where(sq.LtOrEq{"transaction_date": endDate})
		query = query.Where(sq.LtOrEq{"transaction_date": endDate})
	}

	// Get total count
	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting and pagination
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sortType == "" {
		sortType = "DESC"
	}

	query = query.OrderBy(orderBy + " " + sortType).
		Limit(limit).
		Offset(skip)

	// Execute data query
	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates an AML alert by ID
func (r *AMLAlertRepository) Update(ctx context.Context, alertID string,
	riskLevel *string,
	riskScore *int,
	alertStatus *string,
	alertDescription *string,
	reviewedBy *string,
	reviewDecision *string,
	officerRemarks *string,
	actionTaken *string,
	transactionBlocked *bool,
	filingRequired *bool,
	filingType *string,
	filingStatus *string,
	filingReference *string,
	filedAt *time.Time,
	filedBy *string,
) (domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(amlAlertTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": alertID}).
		PlaceholderFormat(sq.Dollar)

	// Only update non-nil fields
	if riskLevel != nil {
		query = query.Set("risk_level", *riskLevel)
	}
	if riskScore != nil {
		query = query.Set("risk_score", *riskScore)
	}
	if alertStatus != nil {
		query = query.Set("alert_status", *alertStatus)
	}
	if alertDescription != nil {
		query = query.Set("alert_description", *alertDescription)
	}
	if reviewedBy != nil {
		query = query.Set("reviewed_by", *reviewedBy)
		query = query.Set("reviewed_at", time.Now())
	}
	if reviewDecision != nil {
		query = query.Set("review_decision", *reviewDecision)
	}
	if officerRemarks != nil {
		query = query.Set("officer_remarks", *officerRemarks)
	}
	if actionTaken != nil {
		query = query.Set("action_taken", *actionTaken)
	}
	if transactionBlocked != nil {
		query = query.Set("transaction_blocked", *transactionBlocked)
	}
	if filingRequired != nil {
		query = query.Set("filing_required", *filingRequired)
	}
	if filingType != nil {
		query = query.Set("filing_type", *filingType)
	}
	if filingStatus != nil {
		query = query.Set("filing_status", *filingStatus)
	}
	if filingReference != nil {
		query = query.Set("filing_reference", *filingReference)
	}
	if filedAt != nil {
		query = query.Set("filed_at", *filedAt)
	}
	if filedBy != nil {
		query = query.Set("filed_by", *filedBy)
	}

	query = query.Suffix("RETURNING *")

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return result, err
}

// UpdateStatus updates alert status
func (r *AMLAlertRepository) UpdateStatus(ctx context.Context, alertID string, status string, reviewedBy string, reviewDecision *string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(amlAlertTable).
		Set("alert_status", status).
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", time.Now()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": alertID}).
		PlaceholderFormat(sq.Dollar)

	if reviewDecision != nil {
		query = query.Set("review_decision", *reviewDecision)
	}

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// UpdateFiling updates filing information
// Reference: BR-CLM-AML-006 (STR filing within 7 days), BR-CLM-AML-007 (CTR filing monthly)
func (r *AMLAlertRepository) UpdateFiling(ctx context.Context, alertID string, filingType string, filingReference string, filedBy string) (domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(amlAlertTable).
		Set("filing_type", filingType).
		Set("filing_reference", filingReference).
		Set("filing_status", "FILED").
		Set("filed_at", time.Now()).
		Set("filed_by", filedBy).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": alertID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return result, err
}

// GetHighRiskAlerts retrieves all HIGH and CRITICAL risk alerts
// Reference: BR-CLM-AML-002 (High-risk customer review)
func (r *AMLAlertRepository) GetHighRiskAlerts(ctx context.Context) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Or{
			sq.Eq{"risk_level": "HIGH"},
			sq.Eq{"risk_level": "CRITICAL"},
		}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetPendingReviewAlerts retrieves alerts pending review
func (r *AMLAlertRepository) GetPendingReviewAlerts(ctx context.Context) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"alert_status": "FLAGGED"}).
		OrderBy("risk_level DESC, created_at ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetAlertsRequiringFiling retrieves alerts requiring STR/CTR filing
// Reference: BR-CLM-AML-006 (STR filing), BR-CLM-AML-007 (CTR filing)
func (r *AMLAlertRepository) GetAlertsRequiringFiling(ctx context.Context, filingType string) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.And{
			sq.Eq{"filing_required": true},
			sq.NotEq{"alert_status": "CLOSED"},
		}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	if filingType != "" {
		query = query.Where(sq.Eq{"filing_type": filingType})
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetOverdueFilingAlerts retrieves alerts with overdue filing
// Reference: BR-CLM-AML-006 (STR within 7 days)
func (r *AMLAlertRepository) GetOverdueFilingAlerts(ctx context.Context, daysOverdue int) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.And{
			sq.Eq{"filing_required": true},
			sq.NotEq{"filing_status": "FILED"},
			sq.LtOrEq{"created_at": sq.Expr("NOW() - INTERVAL '" + fmt.Sprintf("%d days", daysOverdue) + "'")},
		}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetBlockedTransactions retrieves alerts with blocked transactions
func (r *AMLAlertRepository) GetBlockedTransactions(ctx context.Context) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"transaction_blocked": true}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetCustomerRiskHistory retrieves AML alert history for a customer
func (r *AMLAlertRepository) GetCustomerRiskHistory(ctx context.Context, customerID string, months int) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.And{
			sq.Eq{"customer_id": customerID},
			sq.GtOrEq{"created_at": sq.Expr("NOW() - INTERVAL '" + fmt.Sprintf("%d months", months) + "'")},
		}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetAlertsByTriggerCode retrieves alerts by specific trigger code
// Reference: BR-CLM-AML-001 to 005 (Specific AML triggers)
func (r *AMLAlertRepository) GetAlertsByTriggerCode(ctx context.Context, triggerCode string) ([]domain.AMLAlert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(amlAlertTable).
		Where(sq.Eq{"trigger_code": triggerCode}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.AMLAlert])
	return results, err
}

// GetRiskScoreDistribution aggregates risk score distribution
func (r *AMLAlertRepository) GetRiskScoreDistribution(ctx context.Context) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("risk_level", "COUNT(*) as count").
		From(amlAlertTable).
		GroupBy("risk_level").
		PlaceholderFormat(sq.Dollar)

	type RiskCount struct {
		RiskLevel string `db:"risk_level"`
		Count     int64  `db:"count"`
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[RiskCount])
	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int64)
	for _, result := range results {
		distribution[result.RiskLevel] = result.Count
	}

	return distribution, nil
}

// GetAlertsStats retrieves AML alert statistics
func (r *AMLAlertRepository) GetAlertsStats(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build stats query with time range
	query := sq.Select(
		"COUNT(*) as total_alerts",
		"COUNT(*) FILTER (WHERE risk_level IN ('HIGH', 'CRITICAL')) as high_risk_alerts",
		"COUNT(*) FILTER (WHERE alert_status = 'FLAGGED') as pending_review",
		"COUNT(*) FILTER (WHERE filing_required = true) as requiring_filing",
		"COUNT(*) FILTER (WHERE filing_status = 'FILED') as filed",
		"COUNT(*) FILTER (WHERE transaction_blocked = true) as blocked_transactions",
	).
		From(amlAlertTable).
		Where(sq.And{
			sq.GtOrEq{"created_at": startDate},
			sq.LtOrEq{"created_at": endDate},
		}).
		PlaceholderFormat(sq.Dollar)

	type AlertStats struct {
		TotalAlerts        int64 `db:"total_alerts"`
		HighRiskAlerts     int64 `db:"high_risk_alerts"`
		PendingReview      int64 `db:"pending_review"`
		RequiringFiling    int64 `db:"requiring_filing"`
		Filed              int64 `db:"filed"`
		BlockedTransactions int64 `db:"blocked_transactions"`
	}

	stats, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[AlertStats])
	if err != nil {
		return nil, err
	}

	resultMap := map[string]int64{
		"total_alerts":         stats.TotalAlerts,
		"high_risk_alerts":     stats.HighRiskAlerts,
		"pending_review":       stats.PendingReview,
		"requiring_filing":     stats.RequiringFiling,
		"filed":                stats.Filed,
		"blocked_transactions": stats.BlockedTransactions,
	}

	return resultMap, nil
}

// BatchUpdateFilingStatus updates filing status for multiple alerts
// Reference: BR-CLM-AML-006/007 (Bulk filing updates)
func (r *AMLAlertRepository) BatchUpdateFilingStatus(ctx context.Context, alertIDs []string, filingStatus string, filedBy string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	if len(alertIDs) == 0 {
		return nil
	}

	query := sq.Update(amlAlertTable).
		Set("filing_status", filingStatus).
		Set("filed_at", time.Now()).
		Set("filed_by", filedBy).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": alertIDs}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// CheckDuplicateAlert checks if a similar alert already exists
func (r *AMLAlertRepository) CheckDuplicateAlert(ctx context.Context, policyID string, triggerCode string, transactionDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("COUNT(*)").
		From(amlAlertTable).
		Where(sq.And{
			sq.Eq{"policy_id": policyID},
			sq.Eq{"trigger_code": triggerCode},
			sq.Eq{"transaction_date": transactionDate},
		}).
		PlaceholderFormat(sq.Dollar)

	count, err := dblib.SelectOne(ctx, r.db, query, pgx.RowTo[int64])
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete soft deletes an AML alert (marks as closed)
func (r *AMLAlertRepository) Delete(ctx context.Context, alertID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(amlAlertTable).
		Set("alert_status", "CLOSED").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": alertID}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}
