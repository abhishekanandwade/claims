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

const appealTable = "appeals"

// Create inserts a new appeal
// Reference: BR-CLM-DC-005 (90-day appeal window), seed/db/claims_database_schema.sql:304-328
func (r *AppealRepository) Create(ctx context.Context, data domain.Appeal) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(appealTable).
		Columns(
			"appeal_number", "claim_id", "appellant_name", "appellant_contact",
			"grounds_of_appeal", "supporting_documents",
			"condonation_request", "condonation_reason",
			"submission_date", "appeal_deadline", "appellate_authority_id",
			"status", "decision", "reasoned_order", "revised_claim_amount", "decision_date",
		).
		Values(
			data.AppealNumber, data.ClaimID, data.AppellantName, data.AppellantContact,
			data.GroundsOfAppeal, data.SupportingDocuments,
			data.CondonationRequest, data.CondonationReason,
			data.SubmissionDate, data.AppealDeadline, data.AppellateAuthorityID,
			data.Status, data.Decision, data.ReasonedOrder, data.RevisedClaimAmount, data.DecisionDate,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	return result, err
}

// FindByID retrieves an appeal by ID
func (r *AppealRepository) FindByID(ctx context.Context, appealID string) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"id": appealID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByAppealNumber retrieves an appeal by appeal number
func (r *AppealRepository) FindByAppealNumber(ctx context.Context, appealNumber string) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"appeal_number": appealNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByClaimID retrieves appeals by claim ID
func (r *AppealRepository) FindByClaimID(ctx context.Context, claimID string) ([]domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("submission_date DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return results, err
	}
	return results, nil
}

// List retrieves appeals with filters and pagination
func (r *AppealRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.Appeal, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query
	query := sq.Select("*").
		From(appealTable)

	// Apply filters
	if status, ok := filters["status"]; ok && status != nil {
		query = query.Where(sq.Eq{"status": status})
	}
	if claimID, ok := filters["claim_id"]; ok && claimID != nil {
		query = query.Where(sq.Eq{"claim_id": claimID})
	}
	if appellateAuthorityID, ok := filters["appellate_authority_id"]; ok && appellateAuthorityID != nil {
		query = query.Where(sq.Eq{"appellate_authority_id": appellateAuthorityID})
	}
	if condonationRequested, ok := filters["condonation_requested"]; ok && condonationRequested != nil {
		query = query.Where(sq.Eq{"condonation_request": condonationRequested})
	}
	if startDate, ok := filters["start_date"]; ok && startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			query = query.Where(sq.GtOrEq{"submission_date": t})
		}
	}
	if endDate, ok := filters["end_date"]; ok && endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			query = query.Where(sq.LtOrEq{"submission_date": t})
		}
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(appealTable)
	if _, ok := filters["status"]; ok && filters["status"] != nil {
		countQuery = countQuery.Where(sq.Eq{"status": filters["status"]})
	}
	if _, ok := filters["claim_id"]; ok && filters["claim_id"] != nil {
		countQuery = countQuery.Where(sq.Eq{"claim_id": filters["claim_id"]})
	}
	if _, ok := filters["appellate_authority_id"]; ok && filters["appellate_authority_id"] != nil {
		countQuery = countQuery.Where(sq.Eq{"appellate_authority_id": filters["appellate_authority_id"]})
	}
	if _, ok := filters["condonation_requested"]; ok && filters["condonation_requested"] != nil {
		countQuery = countQuery.Where(sq.Eq{"condonation_request": filters["condonation_requested"]})
	}
	if _, ok := filters["start_date"]; ok && filters["start_date"] != nil {
		if t, ok := filters["start_date"].(time.Time); ok {
			countQuery = countQuery.Where(sq.GtOrEq{"submission_date": t})
		}
	}
	if _, ok := filters["end_date"]; ok && filters["end_date"] != nil {
		if t, ok := filters["end_date"].(time.Time); ok {
			countQuery = countQuery.Where(sq.LtOrEq{"submission_date": t})
		}
	}
	countQuery = countQuery.PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count appeals: %w", err)
	}

	// Apply pagination and sorting
	query = query.OrderBy("submission_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return results, 0, err
	}
	return results, totalCount, nil
}

// Update updates an appeal
func (r *AppealRepository) Update(ctx context.Context, appealID string, updates map[string]interface{}) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(appealTable)

	// Build update fields dynamically
	if status, ok := updates["status"]; ok {
		query = query.Set("status", status)
	}
	if appellateAuthorityID, ok := updates["appellate_authority_id"]; ok {
		query = query.Set("appellate_authority_id", appellateAuthorityID)
	}
	if decision, ok := updates["decision"]; ok {
		query = query.Set("decision", decision)
	}
	if reasonedOrder, ok := updates["reasoned_order"]; ok {
		query = query.Set("reasoned_order", reasonedOrder)
	}
	if revisedClaimAmount, ok := updates["revised_claim_amount"]; ok {
		query = query.Set("revised_claim_amount", revisedClaimAmount)
	}
	if decisionDate, ok := updates["decision_date"]; ok {
		query = query.Set("decision_date", decisionDate)
	}
	if supportingDocuments, ok := updates["supporting_documents"]; ok {
		query = query.Set("supporting_documents", supportingDocuments)
	}

	query = query.Set("updated_at", time.Now()).
		Where(sq.Eq{"id": appealID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateStatus updates the status of an appeal
// Reference: BR-CLM-DC-007 (45-day decision timeline)
func (r *AppealRepository) UpdateStatus(ctx context.Context, appealID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(appealTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": appealID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	return err
}

// RecordDecision records an appeal decision
// Reference: BR-CLM-DC-007 (45-day decision timeline)
func (r *AppealRepository) RecordDecision(ctx context.Context, appealID string, decision, reasonedOrder string, revisedClaimAmount *float64) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(appealTable).
		Set("decision", decision).
		Set("reasoned_order", reasonedOrder).
		Set("revised_claim_amount", revisedClaimAmount).
		Set("decision_date", time.Now()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": appealID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetPendingReviewAppeals retrieves appeals pending review
func (r *AppealRepository) GetPendingReviewAppeals(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.Appeal, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query for pending appeals
	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"status": "SUBMITTED"})

	// Apply additional filters
	if appellateAuthorityID, ok := filters["appellate_authority_id"]; ok && appellateAuthorityID != nil {
		query = query.Where(sq.Eq{"appellate_authority_id": appellateAuthorityID})
	}
	if overdueOnly, ok := filters["overdue_only"]; ok && overdueOnly == true {
		query = query.Where(sq.Lt{"appeal_deadline": time.Now()})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(appealTable).Where(sq.Eq{"status": "SUBMITTED"})
	if _, ok := filters["appellate_authority_id"]; ok && filters["appellate_authority_id"] != nil {
		countQuery = countQuery.Where(sq.Eq{"appellate_authority_id": filters["appellate_authority_id"]})
	}
	if _, ok := filters["overdue_only"]; ok && filters["overdue_only"] == true {
		countQuery = countQuery.Where(sq.Lt{"appeal_deadline": time.Now()})
	}
	countQuery = countQuery.PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending appeals: %w", err)
	}

	// Apply pagination and sorting
	query = query.OrderBy("submission_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return results, 0, err
	}
	return results, totalCount, nil
}

// GetOverdueAppeals retrieves appeals with breached decision deadline
// Reference: BR-CLM-DC-007 (45-day decision timeline)
func (r *AppealRepository) GetOverdueAppeals(ctx context.Context) ([]domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"status": "SUBMITTED"}).
		Where(sq.Lt{"appeal_deadline": time.Now()}).
		OrderBy("appeal_deadline ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return results, err
	}
	return results, nil
}

// CheckAppealEligibility checks if a claim is eligible for appeal
// Reference: BR-CLM-DC-005 (90-day appeal window)
func (r *AppealRepository) CheckAppealEligibility(ctx context.Context, claimID string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// Check if claim exists and is rejected
	claimQuery := sq.Select("status", "updated_at").
		From("claims").
		Where(sq.Eq{"id": claimID}).
		PlaceholderFormat(sq.Dollar)

	claimResult, err := dblib.SelectOne(ctx, r.db, claimQuery, pgx.RowToStructByPos[struct {
		Status    string
		UpdatedAt time.Time
	}])
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "Claim not found", nil
		}
		return false, "", fmt.Errorf("failed to check claim: %w", err)
	}

	// Check if claim is rejected
	if claimResult.Status != "REJECTED" {
		return false, "Claim is not in rejected status", nil
	}

	// Check 90-day appeal window (BR-CLM-DC-005)
	ninetyDaysLater := claimResult.UpdatedAt.AddDate(0, 0, 90)
	if time.Now().After(ninetyDaysLater) {
		return false, "Appeal window (90 days) has expired", nil
	}

	// Check if appeal already exists for this claim
	existingAppealQuery := sq.Select("COUNT(*)").
		From(appealTable).
		Where(sq.Eq{"claim_id": claimID}).
		Where(sq.NotEq{"status": "DISMISSED"}).
		PlaceholderFormat(sq.Dollar)

	count, err := dblib.SelectOne(ctx, r.db, existingAppealQuery, pgx.RowTo[int64])
	if err != nil {
		return false, "", fmt.Errorf("failed to check existing appeals: %w", err)
	}

	if count > 0 {
		return false, "An appeal is already pending for this claim", nil
	}

	return true, "Eligible for appeal", nil
}

// GetAppealStats retrieves appeal statistics
func (r *AppealRepository) GetAppealStats(ctx context.Context) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select(
		"COUNT(*) FILTER (WHERE status = 'SUBMITTED') as submitted",
		"COUNT(*) FILTER (WHERE status = 'UNDER_REVIEW') as under_review",
		"COUNT(*) FILTER (WHERE status = 'ALLOWED') as allowed",
		"COUNT(*) FILTER (WHERE status = 'DISMISSED') as dismissed",
		"COUNT(*) FILTER (WHERE status = 'SUBMITTED' AND appeal_deadline < NOW()) as overdue",
	).
		From(appealTable).
		PlaceholderFormat(sq.Dollar)

	type Stats struct {
		Submitted   int64
		UnderReview int64
		Allowed     int64
		Dismissed   int64
		Overdue     int64
	}

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[Stats])
	if err != nil {
		return nil, fmt.Errorf("failed to get appeal stats: %w", err)
	}

	stats := map[string]int64{
		"submitted":    result.Submitted,
		"under_review": result.UnderReview,
		"allowed":      result.Allowed,
		"dismissed":    result.Dismissed,
		"overdue":      result.Overdue,
	}

	return stats, nil
}

// Delete soft deletes an appeal (sets deleted_at if supported, otherwise hard delete)
func (r *AppealRepository) Delete(ctx context.Context, appealID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(appealTable).
		Where(sq.Eq{"id": appealID}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Delete(ctx, r.db, query)
	return err
}

// GenerateAppealNumber generates a unique appeal number
// Format: APL{YYYY}{DDDD}
func (r *AppealRepository) GenerateAppealNumber(ctx context.Context) (string, error) {
	year := time.Now().Format("2006")

	// Count appeals in current year
	countQuery := sq.Select("COUNT(*)").
		From(appealTable).
		Where(sq.Like{"appeal_number": fmt.Sprintf("APL%s%%", year)}).
		PlaceholderFormat(sq.Dollar)

	count, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return "", fmt.Errorf("failed to generate appeal number: %w", err)
	}

	// Generate sequential number with 4 digits
	sequence := count + 1
	appealNumber := fmt.Sprintf("APL%s%04d", year, sequence)

	return appealNumber, nil
}

// ValidateAppealDeadline validates if appeal is within deadline
// Reference: BR-CLM-DC-005 (90-day appeal window)
func (r *AppealRepository) ValidateAppealDeadline(ctx context.Context, claimID string) (bool, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// Get claim rejection date
	claimQuery := sq.Select("updated_at").
		From("claims").
		Where(sq.Eq{"id": claimID}).
		Where(sq.Eq{"status": "REJECTED"}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, claimQuery, pgx.RowTo[time.Time])
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, time.Time{}, fmt.Errorf("claim not found or not rejected")
		}
		return false, time.Time{}, err
	}

	rejectionDate := result
	deadline := rejectionDate.AddDate(0, 0, 90) // 90 days from rejection
	isWithinDeadline := time.Now().Before(deadline) || time.Now().Equal(deadline)

	return isWithinDeadline, deadline, nil
}

// CalculateDecisionDeadline calculates the decision deadline for an appeal
// Reference: BR-CLM-DC-007 (45-day decision timeline)
func (r *AppealRepository) CalculateDecisionDeadline(submissionDate time.Time) time.Time {
	return submissionDate.AddDate(0, 0, 45) // 45 days from submission
}

// AssignAppellateAuthority assigns an appellate authority to an appeal
func (r *AppealRepository) AssignAppellateAuthority(ctx context.Context, appealID, authorityID string) (domain.Appeal, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(appealTable).
		Set("appellate_authority_id", authorityID).
		Set("status", "UNDER_REVIEW").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": appealID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAppealsByAuthority retrieves appeals assigned to an appellate authority
func (r *AppealRepository) GetAppealsByAuthority(ctx context.Context, authorityID string, skip, limit int64) ([]domain.Appeal, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query
	query := sq.Select("*").
		From(appealTable).
		Where(sq.Eq{"appellate_authority_id": authorityID})

	// Get total count
	countQuery := sq.Select("COUNT(*)").
		From(appealTable).
		Where(sq.Eq{"appellate_authority_id": authorityID}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count appeals: %w", err)
	}

	// Apply pagination and sorting
	query = query.OrderBy("submission_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Appeal])
	if err != nil {
		return results, 0, err
	}
	return results, totalCount, nil
}
