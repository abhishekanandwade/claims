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

// InvestigationRepository handles investigation data operations
// Reference: seed/db/claims_database_schema.sql:247-280
// Reference: seed/tool-docs/db-README.md - n-api-db patterns
type InvestigationRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewInvestigationRepository creates a new investigation repository
func NewInvestigationRepository(db *dblib.DB, cfg *config.Config) *InvestigationRepository {
	return &InvestigationRepository{
		db:  db,
		cfg: cfg,
	}
}

const investigationTable = "investigations"

// Create inserts a new investigation
// Reference: BR-CLM-DC-001 (Investigation trigger), BR-CLM-DC-002 (21-day SLA)
func (r *InvestigationRepository) Create(ctx context.Context, data domain.Investigation) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(investigationTable).
		Columns(
			"investigation_id", "claim_id", "assigned_by", "investigator_id",
			"investigator_rank", "jurisdiction", "auto_assigned",
			"assignment_date", "due_date", "status", "progress_percentage",
			"investigation_outcome", "cause_of_death", "cause_of_death_verified",
			"hospital_records_verified", "detailed_findings", "recommendation",
			"report_document_id", "submitted_at", "reviewed_by", "reviewed_at",
			"review_decision", "reviewer_remarks", "reinvestigation_count",
		).
		Values(
			data.InvestigationID, data.ClaimID, data.AssignedBy, data.InvestigatorID,
			data.InvestigatorRank, data.Jurisdiction, data.AutoAssigned,
			data.AssignmentDate, data.DueDate, data.Status, data.ProgressPercentage,
			data.InvestigationOutcome, data.CauseOfDeath, data.CauseOfDeathVerified,
			data.HospitalRecordsVerified, data.DetailedFindings, data.Recommendation,
			data.ReportDocumentID, data.SubmittedAt, data.ReviewedBy, data.ReviewedAt,
			data.ReviewDecision, data.ReviewerRemarks, data.ReinvestigationCount,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// FindByID retrieves an investigation by ID
func (r *InvestigationRepository) FindByID(ctx context.Context, id string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByInvestigationID retrieves an investigation by investigation_id
func (r *InvestigationRepository) FindByInvestigationID(ctx context.Context, investigationID string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Eq{"investigation_id": investigationID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByClaimID retrieves all investigations for a claim
func (r *InvestigationRepository) FindByClaimID(ctx context.Context, claimID string) ([]domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("assignment_date DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return results, err
}

// List retrieves all investigations with pagination and filters
func (r *InvestigationRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64, orderBy, sortType string) ([]domain.Investigation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build count query
	countQuery := sq.Select("COUNT(*)").
		From(investigationTable).
		PlaceholderFormat(sq.Dollar)

	// Apply filters to count query
	if status, ok := filters["status"]; ok {
		countQuery = countQuery.Where(sq.Eq{"status": status})
	}
	if investigatorID, ok := filters["investigator_id"]; ok {
		countQuery = countQuery.Where(sq.Eq{"investigator_id": investigatorID})
	}
	if claimID, ok := filters["claim_id"]; ok {
		countQuery = countQuery.Where(sq.Eq{"claim_id": claimID})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	query := sq.Select("*").
		From(investigationTable).
		OrderBy(orderBy + " " + sortType).
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	// Apply filters to data query
	if status, ok := filters["status"]; ok {
		query = query.Where(sq.Eq{"status": status})
	}
	if investigatorID, ok := filters["investigator_id"]; ok {
		query = query.Where(sq.Eq{"investigator_id": investigatorID})
	}
	if claimID, ok := filters["claim_id"]; ok {
		query = query.Where(sq.Eq{"claim_id": claimID})
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates an investigation by ID
func (r *InvestigationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	// Apply updates dynamically
	for key, value := range updates {
		query = query.Set(key, value)
	}

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// UpdateStatus updates investigation status
func (r *InvestigationRepository) UpdateStatus(ctx context.Context, id string, status string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// UpdateProgress updates investigation progress percentage
// Reference: BR-CLM-DC-002 (Progress tracking)
func (r *InvestigationRepository) UpdateProgress(ctx context.Context, id string, progressPercentage int) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationTable).
		Set("progress_percentage", progressPercentage).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// GetActiveInvestigations retrieves all active investigations
// Reference: BR-CLM-DC-002 (Active investigation monitoring)
func (r *InvestigationRepository) GetActiveInvestigations(ctx context.Context, skip, limit int64) ([]domain.Investigation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(investigationTable).
		Where(sq.NotEq{"status": "COMPLETED"}).
		Where(sq.NotEq{"status": "CANCELLED"}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := sq.Select("*").
		From(investigationTable).
		Where(sq.NotEq{"status": "COMPLETED"}).
		Where(sq.NotEq{"status": "CANCELLED"}).
		OrderBy("assignment_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// GetOverdueInvestigations retrieves investigations past due date
// Reference: BR-CLM-DC-002 (21-day SLA breach)
func (r *InvestigationRepository) GetOverdueInvestigations(ctx context.Context) ([]domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Lt{"due_date": time.Now()}).
		Where(sq.NotEq{"status": "COMPLETED"}).
		Where(sq.NotEq{"status": "CANCELLED"}).
		OrderBy("due_date ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return results, err
}

// GetPendingInvestigationClaims retrieves claims pending investigation
// Reference: BR-CLM-DC-001 (Investigation required claims)
func (r *InvestigationRepository) GetPendingInvestigationClaims(ctx context.Context, skip, limit int64) ([]domain.Investigation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(investigationTable).
		Where(sq.Eq{"status": "ASSIGNED"}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Eq{"status": "ASSIGNED"}).
		OrderBy("assignment_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// SubmitReport submits investigation report
// Reference: BR-CLM-DC-002 (Report submission)
func (r *InvestigationRepository) SubmitReport(ctx context.Context, id string, outcome string, findings string, recommendation string, reportDocumentID string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationTable).
		Set("investigation_outcome", outcome).
		Set("detailed_findings", findings).
		Set("recommendation", recommendation).
		Set("report_document_id", reportDocumentID).
		Set("submitted_at", time.Now()).
		Set("status", "SUBMITTED").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// ReviewReport reviews investigation report
// Reference: BR-CLM-DC-011 (Review within 5 days)
func (r *InvestigationRepository) ReviewReport(ctx context.Context, id string, reviewedBy string, decision string, remarks string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationTable).
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", time.Now()).
		Set("review_decision", decision).
		Set("reviewer_remarks", remarks).
		Set("status", "REVIEWED").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// TriggerReinvestigation triggers a reinvestigation
// Reference: BR-CLM-DC-012 (Max 2 reinvestigations, 14 days each)
func (r *InvestigationRepository) TriggerReinvestigation(ctx context.Context, id string) (domain.Investigation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// First, check reinvestigation count
	investigation, err := r.FindByID(ctx, id)
	if err != nil {
		return domain.Investigation{}, err
	}

	// BR-CLM-DC-012: Max 2 reinvestigations
	if investigation.ReinvestigationCount >= 2 {
		return domain.Investigation{}, pgx.ErrNoRows // Will be handled as error by handler
	}

	// Increment reinvestigation count and reset for new investigation
	newDueDate := time.Now().AddDate(0, 0, 14) // 14 days for reinvestigation

	query := sq.Update(investigationTable).
		Set("reinvestigation_count", investigation.ReinvestigationCount+1).
		Set("due_date", newDueDate).
		Set("status", "ASSIGNED").
		Set("progress_percentage", 0).
		Set("investigation_outcome", nil).
		Set("detailed_findings", nil).
		Set("recommendation", nil).
		Set("submitted_at", nil).
		Set("reviewed_by", nil).
		Set("reviewed_at", nil).
		Set("review_decision", nil).
		Set("reviewer_remarks", nil).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	return result, err
}

// GetInvestigationsByInvestigator retrieves investigations assigned to a specific investigator
func (r *InvestigationRepository) GetInvestigationsByInvestigator(ctx context.Context, investigatorID string, skip, limit int64) ([]domain.Investigation, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(investigationTable).
		Where(sq.Eq{"investigator_id": investigatorID}).
		Where(sq.NotEq{"status": "COMPLETED"}).
		Where(sq.NotEq{"status": "CANCELLED"}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := sq.Select("*").
		From(investigationTable).
		Where(sq.Eq{"investigator_id": investigatorID}).
		Where(sq.NotEq{"status": "COMPLETED"}).
		Where(sq.NotEq{"status": "CANCELLED"}).
		OrderBy("assignment_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.Investigation])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Delete deletes an investigation by ID (soft delete not implemented, use UpdateStatus to CANCELLED)
func (r *InvestigationRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(investigationTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Delete(ctx, r.db, query)
	return err
}
