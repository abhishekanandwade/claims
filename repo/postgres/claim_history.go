package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

// ClaimHistoryRepository handles claim history data operations
type ClaimHistoryRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimHistoryRepository creates a new claim history repository
func NewClaimHistoryRepository(db *dblib.DB, cfg *config.Config) *ClaimHistoryRepository {
	return &ClaimHistoryRepository{
		db:  db,
		cfg: cfg,
	}
}

// ClaimHistory represents a claim history record
type ClaimHistory struct {
	ID            uuid.UUID `json:"id"`
	ClaimID       uuid.UUID `json:"claim_id"`
	EntityType    string    `json:"entity_type"`    // CLAIM, DOCUMENT, INVESTIGATION, PAYMENT, APPEAL, etc.
	EntityID      uuid.UUID `json:"entity_id"`      // ID of the entity
	ActionType    string    `json:"action_type"`    // CREATED, UPDATED, DELETED, STATUS_CHANGED, etc.
	OldValue      string    `json:"old_value"`      // JSON string of old value
	NewValue      string    `json:"new_value"`      // JSON string of new value
	ChangedFields string    `json:"changed_fields"` // JSON array of field names
	ChangedBy     string    `json:"changed_by"`     // User ID who made the change
	ChangeReason  string    `json:"change_reason"`  // Reason for the change
	IPAddress     string    `json:"ip_address"`     // IP address of the user
	UserAgent     string    `json:"user_agent"`     // User agent string
	Metadata      string    `json:"metadata"`       // Additional metadata as JSON
	CreatedAt     time.Time `json:"created_at"`
}

// Create creates a new claim history record
func (r *ClaimHistoryRepository) Create(ctx context.Context, history ClaimHistory) (ClaimHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	history.ID = uuid.New()
	history.CreatedAt = time.Now()

	query := sq.
		Insert("claim_history").
		Columns(
			"id", "claim_id", "entity_type", "entity_id",
			"action_type", "old_value", "new_value", "changed_fields",
			"changed_by", "change_reason", "ip_address", "user_agent",
			"metadata", "created_at",
		).
		Values(
			history.ID, history.ClaimID, history.EntityType, history.EntityID,
			history.ActionType, history.OldValue, history.NewValue, history.ChangedFields,
			history.ChangedBy, history.ChangeReason, history.IPAddress, history.UserAgent,
			history.Metadata, history.CreatedAt,
		).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[ClaimHistory])
	if err != nil {
		return ClaimHistory{}, fmt.Errorf("failed to create claim history: %w", err)
	}

	return result, nil
}

// FindByID finds a claim history record by ID
func (r *ClaimHistoryRepository) FindByID(ctx context.Context, historyID string) (ClaimHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"id": historyID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ClaimHistory{}, fmt.Errorf("claim history not found")
		}
		return ClaimHistory{}, fmt.Errorf("failed to find claim history: %w", err)
	}

	return result, nil
}

// FindByClaimID retrieves all history records for a claim
func (r *ClaimHistoryRepository) FindByClaimID(ctx context.Context, claimID string, entityType string, skip, limit int) ([]ClaimHistory, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// Build base query
	baseQuery := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"claim_id": claimID})

	// Add entity type filter if specified
	if entityType != "" {
		baseQuery = baseQuery.Where(sq.Eq{"entity_type": entityType})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From("claim_history").Where(sq.Eq{"claim_id": claimID})
	if entityType != "" {
		countQuery = countQuery.Where(sq.Eq{"entity_type": entityType})
	}

	total, err := dblib.SelectOne[int64](
		ctx,
		r.db,
		countQuery.PlaceholderFormat(sq.Dollar),
		pgx.RowTo[int64],
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count claim history: %w", err)
	}

	// Get paginated results
	query := baseQuery.
		OrderBy("created_at DESC").
		Offset(uint64(skip)).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find claim history: %w", err)
	}

	return results, total, nil
}

// FindByEntityID retrieves history records for a specific entity
func (r *ClaimHistoryRepository) FindByEntityID(ctx context.Context, entityID string) ([]ClaimHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"entity_id": entityID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find entity history: %w", err)
	}

	return results, nil
}

// GetTimeline retrieves claim timeline with grouped events
func (r *ClaimHistoryRepository) GetTimeline(ctx context.Context, claimID string) ([]ClaimHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim timeline: %w", err)
	}

	return results, nil
}

// GetAuditTrail retrieves audit trail for a claim
func (r *ClaimHistoryRepository) GetAuditTrail(ctx context.Context, claimID string, actionType string, startDate, endDate time.Time, skip, limit int) ([]ClaimHistory, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query
	baseQuery := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"claim_id": claimID})

	// Add action type filter if specified
	if actionType != "" {
		baseQuery = baseQuery.Where(sq.Eq{"action_type": actionType})
	}

	// Add date range filters
	if !startDate.IsZero() {
		baseQuery = baseQuery.Where(sq.GtOrEq{"created_at": startDate})
	}
	if !endDate.IsZero() {
		baseQuery = baseQuery.Where(sq.LtOrEq{"created_at": endDate})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From("claim_history").Where(sq.Eq{"claim_id": claimID})
	if actionType != "" {
		countQuery = countQuery.Where(sq.Eq{"action_type": actionType})
	}
	if !startDate.IsZero() {
		countQuery = countQuery.Where(sq.GtOrEq{"created_at": startDate})
	}
	if !endDate.IsZero() {
		countQuery = countQuery.Where(sq.LtOrEq{"created_at": endDate})
	}

	total, err := dblib.SelectOne[int64](
		ctx,
		r.db,
		countQuery.PlaceholderFormat(sq.Dollar),
		pgx.RowTo[int64],
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit trail: %w", err)
	}

	// Get paginated results
	query := baseQuery.
		OrderBy("created_at DESC").
		Offset(uint64(skip)).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit trail: %w", err)
	}

	return results, total, nil
}

// GetRecentActivity retrieves recent activity for a claim
func (r *ClaimHistoryRepository) GetRecentActivity(ctx context.Context, claimID string, limit int) ([]ClaimHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Select("*").
		From("claim_history").
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows[ClaimHistory](
		ctx,
		r.db,
		query,
		pgx.RowToStructByPos[ClaimHistory],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}

	return results, nil
}

// BatchCreate creates multiple history records in a transaction
func (r *ClaimHistoryRepository) BatchCreate(ctx context.Context, histories []ClaimHistory) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	now := time.Now()
	for _, history := range histories {
		history.ID = uuid.New()
		history.CreatedAt = now

		query := sq.
			Insert("claim_history").
			Columns(
				"id", "claim_id", "entity_type", "entity_id",
				"action_type", "old_value", "new_value", "changed_fields",
				"changed_by", "change_reason", "ip_address", "user_agent",
				"metadata", "created_at",
			).
			Values(
				history.ID, history.ClaimID, history.EntityType, history.EntityID,
				history.ActionType, history.OldValue, history.NewValue, history.ChangedFields,
				history.ChangedBy, history.ChangeReason, history.IPAddress, history.UserAgent,
				history.Metadata, history.CreatedAt,
			).
			PlaceholderFormat(sq.Dollar)

		sql, args, _ := query.ToSql()
		batch.Queue(sql, args...)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(histories); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to create history record %d: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete deletes a claim history record by ID
func (r *ClaimHistoryRepository) Delete(ctx context.Context, historyID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Delete("claim_history").
		Where(sq.Eq{"id": historyID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, _ := query.ToSql()

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete claim history: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("claim history not found")
	}

	return nil
}

// GetHistoryStats gets statistics about claim history
func (r *ClaimHistoryRepository) GetHistoryStats(ctx context.Context, claimID string) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.
		Select("action_type", "COUNT(*) as count").
		From("claim_history").
		Where(sq.Eq{"claim_id": claimID}).
		GroupBy("action_type").
		PlaceholderFormat(sq.Dollar)

	sql, args, _ := query.ToSql()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get history stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var actionType string
		var count int64
		if err := rows.Scan(&actionType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan history stat: %w", err)
		}
		stats[actionType] = count
	}

	return stats, nil
}
