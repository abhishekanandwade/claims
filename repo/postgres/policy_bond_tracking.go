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

// PolicyBondTrackingRepository handles policy bond tracking data operations
// Reference: seed/db/claims_database_schema.sql:584-615
// Reference: seed/tool-docs/db-README.md - n-api-db patterns
type PolicyBondTrackingRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewPolicyBondTrackingRepository creates a new policy bond tracking repository
func NewPolicyBondTrackingRepository(db *dblib.DB, cfg *config.Config) *PolicyBondTrackingRepository {
	return &PolicyBondTrackingRepository{
		db:  db,
		cfg: cfg,
	}
}

const policyBondTrackingTable = "policy_bond_tracking"

// Create inserts a new policy bond tracking record
// Reference: BR-CLM-BOND-001 (Free look period calculation)
func (r *PolicyBondTrackingRepository) Create(ctx context.Context, data domain.PolicyBondTracking) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(policyBondTrackingTable).
		Columns(
			"policy_id", "bond_number", "bond_type",
			"print_date", "dispatch_date", "tracking_number",
			"delivery_date", "delivery_status", "delivery_attempt_count",
			"pod_reference", "recipient_name", "recipient_signature_captured",
			"undelivered_reason", "escalation_triggered", "escalation_date",
			"customer_contacted", "address_verified", "redelivery_requested",
			"freelook_period_start_date", "freelook_period_end_date",
			"freelook_cancellation_submitted", "freelook_cancellation_id",
		).
		Values(
			data.PolicyID, data.BondNumber, data.BondType,
			data.PrintDate, data.DispatchDate, data.TrackingNumber,
			data.DeliveryDate, data.DeliveryStatus, data.DeliveryAttemptCount,
			data.PODReference, data.RecipientName, data.RecipientSignatureCaptured,
			data.UndeliveredReason, data.EscalationTriggered, data.EscalationDate,
			data.CustomerContacted, data.AddressVerified, data.RedeliveryRequested,
			data.FreeLookPeriodStartDate, data.FreeLookPeriodEndDate,
			data.FreeLookCancellationSubmitted, data.FreeLookCancellationID,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	return result, err
}

// FindByID retrieves a policy bond tracking record by ID
func (r *PolicyBondTrackingRepository) FindByID(ctx context.Context, id string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByPolicyID retrieves policy bond tracking record by policy ID
func (r *PolicyBondTrackingRepository) FindByPolicyID(ctx context.Context, policyID string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.Eq{"policy_id": policyID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByBondNumber retrieves policy bond tracking record by bond number
func (r *PolicyBondTrackingRepository) FindByBondNumber(ctx context.Context, bondNumber string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.Eq{"bond_number": bondNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByTrackingNumber retrieves policy bond tracking record by tracking number
func (r *PolicyBondTrackingRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.Eq{"tracking_number": trackingNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return result, err
	}
	return result, nil
}

// List retrieves policy bond tracking records with filters and pagination
func (r *PolicyBondTrackingRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.PolicyBondTracking, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		PlaceholderFormat(sq.Dollar)

	// Apply filters
	if policyID, ok := filters["policy_id"]; ok && policyID != nil {
		query = query.Where(sq.Eq{"policy_id": policyID})
	}
	if bondType, ok := filters["bond_type"]; ok && bondType != nil {
		query = query.Where(sq.Eq{"bond_type": bondType})
	}
	if deliveryStatus, ok := filters["delivery_status"]; ok && deliveryStatus != nil {
		query = query.Where(sq.Eq{"delivery_status": deliveryStatus})
	}
	if escalationTriggered, ok := filters["escalation_triggered"]; ok && escalationTriggered != nil {
		query = query.Where(sq.Eq{"escalation_triggered": escalationTriggered})
	}
	if startDate, ok := filters["start_date"]; ok && startDate != nil {
		query = query.Where(sq.GtOrEq{"dispatch_date": startDate})
	}
	if endDate, ok := filters["end_date"]; ok && endDate != nil {
		query = query.Where(sq.LtOrEq{"dispatch_date": endDate})
	}

	// Get total count
	countQuery := sq.Select("COUNT(*)").From(policyBondTrackingTable)
	// Apply same filters to count query
	if policyID, ok := filters["policy_id"]; ok && policyID != nil {
		countQuery = countQuery.Where(sq.Eq{"policy_id": policyID})
	}
	if bondType, ok := filters["bond_type"]; ok && bondType != nil {
		countQuery = countQuery.Where(sq.Eq{"bond_type": bondType})
	}
	if deliveryStatus, ok := filters["delivery_status"]; ok && deliveryStatus != nil {
		countQuery = countQuery.Where(sq.Eq{"delivery_status": deliveryStatus})
	}
	if escalationTriggered, ok := filters["escalation_triggered"]; ok && escalationTriggered != nil {
		countQuery = countQuery.Where(sq.Eq{"escalation_triggered": escalationTriggered})
	}
	if startDate, ok := filters["start_date"]; ok && startDate != nil {
		countQuery = countQuery.Where(sq.GtOrEq{"dispatch_date": startDate})
	}
	if endDate, ok := filters["end_date"]; ok && endDate != nil {
		countQuery = countQuery.Where(sq.LtOrEq{"dispatch_date": endDate})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	query = query.OrderBy("created_at DESC").Limit(uint64(limit)).Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates a policy bond tracking record
func (r *PolicyBondTrackingRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(policyBondTrackingTable).
		Where(sq.Eq{"id": id}).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar)

	for key, value := range updates {
		query = query.Set(key, value)
	}

	query = query.Suffix("RETURNING *")

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	return result, err
}

// UpdateDeliveryStatus updates delivery status and related fields
// Reference: BR-CLM-BOND-002 (Delivery failure escalation)
func (r *PolicyBondTrackingRepository) UpdateDeliveryStatus(ctx context.Context, id string, deliveryStatus string, deliveryDate *time.Time, trackingNumber *string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	updates := map[string]interface{}{
		"delivery_status": deliveryStatus,
		"updated_at":      time.Now(),
	}

	if deliveryDate != nil {
		updates["delivery_date"] = *deliveryDate
	}
	if trackingNumber != nil {
		updates["tracking_number"] = *trackingNumber
	}

	// Update free look period start date based on delivery for physical bonds
	// BR-CLM-BOND-001: Physical bonds - free look starts from delivery date
	if deliveryDate != nil {
		updates["freelook_period_start_date"] = *deliveryDate
	}

	return r.Update(ctx, id, updates)
}

// UpdateEscalation updates escalation status for undelivered bonds
// Reference: BR-CLM-BOND-002 (Escalation after 10 days)
func (r *PolicyBondTrackingRepository) UpdateEscalation(ctx context.Context, id string, escalationTriggered bool, escalationDate *time.Time) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	updates := map[string]interface{}{
		"escalation_triggered": escalationTriggered,
		"updated_at":           time.Now(),
	}

	if escalationDate != nil {
		updates["escalation_date"] = *escalationDate
	}

	return r.Update(ctx, id, updates)
}

// UpdatePOD updates proof of delivery information
func (r *PolicyBondTrackingRepository) UpdatePOD(ctx context.Context, id string, podReference *string, recipientName *string, signatureCaptured bool) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	updates := map[string]interface{}{
		"pod_reference":               podReference,
		"recipient_name":              recipientName,
		"recipient_signature_captured": signatureCaptured,
		"updated_at":                  time.Now(),
	}

	return r.Update(ctx, id, updates)
}

// UpdateCustomerInteraction updates customer contact and address verification status
func (r *PolicyBondTrackingRepository) UpdateCustomerInteraction(ctx context.Context, id string, customerContacted bool, addressVerified bool, redeliveryRequested bool) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	updates := map[string]interface{}{
		"customer_contacted":   customerContacted,
		"address_verified":     addressVerified,
		"redelivery_requested": redeliveryRequested,
		"updated_at":           time.Now(),
	}

	return r.Update(ctx, id, updates)
}

// LinkFreeLookCancellation links a free look cancellation to bond tracking
func (r *PolicyBondTrackingRepository) LinkFreeLookCancellation(ctx context.Context, policyID string, cancellationID string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(policyBondTrackingTable).
		Where(sq.Eq{"policy_id": policyID}).
		Set("freelook_cancellation_submitted", true).
		Set("freelook_cancellation_id", cancellationID).
		Set("updated_at", time.Now()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	return result, err
}

// GetUndeliveredBonds retrieves bonds that were not delivered
// Reference: BR-CLM-BOND-002 (Undelivered bonds after 10 days)
func (r *PolicyBondTrackingRepository) GetUndeliveredBonds(ctx context.Context, daysSinceDispatch int) ([]domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	cutoffDate := time.Now().AddDate(0, 0, -daysSinceDispatch)

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.And{
			sq.Eq{"delivery_status": "FAILED"},
			sq.LtOrEq{"dispatch_date": cutoffDate},
			sq.Eq{"escalation_triggered": false},
		}).
		OrderBy("dispatch_date ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetBondsRequiringEscalation retrieves bonds that need escalation
func (r *PolicyBondTrackingRepository) GetBondsRequiringEscalation(ctx context.Context) ([]domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// BR-CLM-BOND-002: Escalate after 10 days of delivery failure
	return r.GetUndeliveredBonds(ctx, 10)
}

// GetActiveFreeLookPeriodBonds retrieves bonds within free look period
// Reference: BR-CLM-BOND-001 (Free look period: 15 days physical, 30 days electronic)
func (r *PolicyBondTrackingRepository) GetActiveFreeLookPeriodBonds(ctx context.Context) ([]domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.And{
			sq.NotEq{"delivery_status": "FAILED"},
			sq.NotEq{"freelook_period_end_date": nil},
			sq.GtOrEq{"freelook_period_end_date": time.Now()},
		}).
		OrderBy("freelook_period_end_date ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetExpiredFreeLookPeriodBonds retrieves bonds past free look period
func (r *PolicyBondTrackingRepository) GetExpiredFreeLookPeriodBonds(ctx context.Context) ([]domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(policyBondTrackingTable).
		Where(sq.And{
			sq.NotEq{"delivery_status": "FAILED"},
			sq.NotEq{"freelook_period_end_date": nil},
			sq.Lt{"freelook_period_end_date": time.Now()},
			sq.Eq{"freelook_cancellation_submitted": false},
		}).
		OrderBy("freelook_period_end_date DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Delete deletes a policy bond tracking record
func (r *PolicyBondTrackingRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(policyBondTrackingTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Delete(ctx, r.db, query)
	return err
}

// GetDeliveryStats returns delivery statistics
func (r *PolicyBondTrackingRepository) GetDeliveryStats(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select(
		"delivery_status",
		"COUNT(*) as count",
	).
		From(policyBondTrackingTable).
		Where(sq.And{
			sq.GtOrEq{"dispatch_date": startDate},
			sq.LtOrEq{"dispatch_date": endDate},
		}).
		GroupBy("delivery_status").
		PlaceholderFormat(sq.Dollar)

	type Result struct {
		DeliveryStatus string `db:"delivery_status"`
		Count          int64  `db:"count"`
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[Result])
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.DeliveryStatus] = result.Count
	}

	return stats, nil
}

// IncrementDeliveryAttempt increments delivery attempt count
func (r *PolicyBondTrackingRepository) IncrementDeliveryAttempt(ctx context.Context, id string) (domain.PolicyBondTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(policyBondTrackingTable).
		Where(sq.Eq{"id": id}).
		Set("delivery_attempt_count", sq.Expr("delivery_attempt_count + 1")).
		Set("updated_at", time.Now()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.PolicyBondTracking])
	return result, err
}
