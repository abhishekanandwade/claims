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

// InvestigationProgressRepository handles investigation progress data operations
// Reference: seed/db/claims_database_schema.sql:287-299
// Reference: seed/tool-docs/db-README.md - n-api-db patterns
type InvestigationProgressRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewInvestigationProgressRepository creates a new investigation progress repository
func NewInvestigationProgressRepository(db *dblib.DB, cfg *config.Config) *InvestigationProgressRepository {
	return &InvestigationProgressRepository{
		db:  db,
		cfg: cfg,
	}
}

const investigationProgressTable = "investigation_progress"

// Create inserts a new investigation progress update
// Reference: BR-CLM-DC-002 (Progress tracking for investigations)
func (r *InvestigationProgressRepository) Create(ctx context.Context, data domain.InvestigationProgress) (domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(investigationProgressTable).
		Columns(
			"investigation_id", "update_date", "progress_percentage",
			"checklist_items_completed", "remarks", "estimated_completion_date",
			"updated_by",
		).
		Values(
			data.InvestigationID, data.UpdateDate, data.ProgressPercentage,
			data.ChecklistItemsCompleted, data.Remarks, data.EstimatedCompletionDate,
			data.UpdatedBy,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	return result, err
}

// FindByID retrieves an investigation progress by ID
func (r *InvestigationProgressRepository) FindByID(ctx context.Context, id string) (domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationProgressTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByInvestigationID retrieves all progress updates for an investigation
func (r *InvestigationProgressRepository) FindByInvestigationID(ctx context.Context, investigationID string) ([]domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationProgressTable).
		Where(sq.Eq{"investigation_id": investigationID}).
		OrderBy("update_date DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	return results, err
}

// LatestProgress retrieves the latest progress update for an investigation
func (r *InvestigationProgressRepository) LatestProgress(ctx context.Context, investigationID string) (domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationProgressTable).
		Where(sq.Eq{"investigation_id": investigationID}).
		OrderBy("update_date DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	if err != nil {
		return result, err
	}
	return result, nil
}

// List retrieves all investigation progress updates with pagination
func (r *InvestigationProgressRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.InvestigationProgress, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build count query
	countQuery := sq.Select("COUNT(*)").
		From(investigationProgressTable).
		PlaceholderFormat(sq.Dollar)

	// Apply filters to count query
	if investigationID, ok := filters["investigation_id"]; ok {
		countQuery = countQuery.Where(sq.Eq{"investigation_id": investigationID})
	}
	if startDate, ok := filters["start_date"]; ok {
		countQuery = countQuery.Where(sq.GtOrEq{"update_date": startDate})
	}
	if endDate, ok := filters["end_date"]; ok {
		countQuery = countQuery.Where(sq.LtOrEq{"update_date": endDate})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	query := sq.Select("*").
		From(investigationProgressTable).
		OrderBy("update_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	// Apply filters to data query
	if investigationID, ok := filters["investigation_id"]; ok {
		query = query.Where(sq.Eq{"investigation_id": investigationID})
	}
	if startDate, ok := filters["start_date"]; ok {
		query = query.Where(sq.GtOrEq{"update_date": startDate})
	}
	if endDate, ok := filters["end_date"]; ok {
		query = query.Where(sq.LtOrEq{"update_date": endDate})
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates an investigation progress by ID
func (r *InvestigationProgressRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(investigationProgressTable).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	// Apply updates dynamically
	for key, value := range updates {
		query = query.Set(key, value)
	}

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	return result, err
}

// UpdateProgress updates progress percentage and remarks
// Reference: BR-CLM-DC-002 (Heartbeat updates for long-running investigations)
func (r *InvestigationProgressRepository) UpdateProgress(ctx context.Context, investigationID string, progressPercentage int, checklistItems []string, remarks string, estimatedCompletion *time.Time, updatedBy string) (domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(investigationProgressTable).
		Columns(
			"investigation_id", "update_date", "progress_percentage",
			"checklist_items_completed", "remarks", "estimated_completion_date",
			"updated_by",
		).
		Values(
			investigationID, time.Now(), progressPercentage,
			checklistItems, remarks, estimatedCompletion,
			updatedBy,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	return result, err
}

// Delete deletes an investigation progress by ID
func (r *InvestigationProgressRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(investigationProgressTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Delete(ctx, r.db, query)
	return err
}

// GetProgressTimeline retrieves progress updates for an investigation with date range
func (r *InvestigationProgressRepository) GetProgressTimeline(ctx context.Context, investigationID string, startDate, endDate time.Time) ([]domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(investigationProgressTable).
		Where(sq.Eq{"investigation_id": investigationID}).
		Where(sq.GtOrEq{"update_date": startDate}).
		Where(sq.LtOrEq{"update_date": endDate}).
		OrderBy("update_date ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
	return results, err
}

// BatchCreate creates multiple investigation progress updates in a single transaction
// Reference: seed/tool-docs/db-README.md - Batch operations
func (r *InvestigationProgressRepository) BatchCreate(ctx context.Context, data []domain.InvestigationProgress) ([]domain.InvestigationProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Use transaction for batch insert
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	results := make([]domain.InvestigationProgress, 0, len(data))

	for _, item := range data {
		query := sq.Insert(investigationProgressTable).
			Columns(
				"investigation_id", "update_date", "progress_percentage",
				"checklist_items_completed", "remarks", "estimated_completion_date",
				"updated_by",
			).
			Values(
				item.InvestigationID, item.UpdateDate, item.ProgressPercentage,
				item.ChecklistItemsCompleted, item.Remarks, item.EstimatedCompletionDate,
				item.UpdatedBy,
			).
			Suffix("RETURNING *").
			PlaceholderFormat(sq.Dollar)

		result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.InvestigationProgress])
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return results, nil
}
