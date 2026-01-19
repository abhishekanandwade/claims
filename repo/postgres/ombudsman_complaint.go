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

// OmbudsmanComplaintRepository handles ombudsman complaint data operations
// Reference: seed/db/claims_database_schema.sql:532-577
type OmbudsmanComplaintRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewOmbudsmanComplaintRepository creates a new ombudsman complaint repository
func NewOmbudsmanComplaintRepository(db *dblib.DB, cfg *config.Config) *OmbudsmanComplaintRepository {
	return &OmbudsmanComplaintRepository{
		db:  db,
		cfg: cfg,
	}
}

const ombudsmanComplaintTable = "ombudsman_complaints"

// Create inserts a new ombudsman complaint
// Reference: BR-CLM-OMB-001 (Admissibility checks), seed/db/claims_database_schema.sql:532-577
func (r *OmbudsmanComplaintRepository) Create(ctx context.Context, data domain.OmbudsmanComplaint) (domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(ombudsmanComplaintTable).
		Columns(
			"complaint_number", "claim_id", "policy_id",
			"complainant_name", "complainant_contact", "complaint_description", "complaint_category", "claim_value",
			"representation_to_insurer_date", "wait_period_completed", "limitation_period_valid", "parallel_litigation",
			"admissible", "inadmissibility_reason",
			"ombudsman_center", "jurisdiction_basis", "assigned_ombudsman_id", "conflict_of_interest",
			"status",
			"mediation_attempted", "mediation_successful", "mediation_terms",
			"recommendation_issued", "recommendation_date",
			"award_issued", "award_number", "award_amount", "award_date", "award_digitally_signed",
			"compliance_due_date", "compliance_status", "compliance_date", "escalated_to_irdai",
			"closure_date", "archival_date", "retention_period_years",
		).
		Values(
			data.ComplaintNumber, data.ClaimID, data.PolicyID,
			data.ComplainantName, data.ComplainantContact, data.ComplaintDescription, data.ComplaintCategory, data.ClaimValue,
			data.RepresentationToInsurerDate, data.WaitPeriodCompleted, data.LimitationPeriodValid, data.ParallelLitigation,
			data.Admissible, data.InadmissibilityReason,
			data.OmbudsmanCenter, data.JurisdictionBasis, data.AssignedOmbudsmanID, data.ConflictOfInterest,
			data.Status,
			data.MediationAttempted, data.MediationSuccessful, data.MediationTerms,
			data.RecommendationIssued, data.RecommendationDate,
			data.AwardIssued, data.AwardNumber, data.AwardAmount, data.AwardDate, data.AwardDigitallySigned,
			data.ComplianceDueDate, data.ComplianceStatus, data.ComplianceDate, data.EscalatedToIRDAI,
			data.ClosureDate, data.ArchivalDate, data.RetentionPeriodYears,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	return result, err
}

// FindByID retrieves an ombudsman complaint by ID
func (r *OmbudsmanComplaintRepository) FindByID(ctx context.Context, id string) (domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByComplaintNumber retrieves an ombudsman complaint by complaint number
func (r *OmbudsmanComplaintRepository) FindByComplaintNumber(ctx context.Context, complaintNumber string) (domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"complaint_number": complaintNumber}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByClaimID retrieves ombudsman complaints by claim ID
func (r *OmbudsmanComplaintRepository) FindByClaimID(ctx context.Context, claimID string) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByPolicyID retrieves ombudsman complaints by policy ID
func (r *OmbudsmanComplaintRepository) FindByPolicyID(ctx context.Context, policyID string) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"policy_id": policyID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// List retrieves ombudsman complaints with filters and pagination
func (r *OmbudsmanComplaintRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit uint64) ([]domain.OmbudsmanComplaint, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable)

	// Apply filters
	if status, ok := filters["status"]; ok && status != nil {
		query = query.Where(sq.Eq{"status": status})
	}
	if ombudsmanCenter, ok := filters["ombudsman_center"]; ok && ombudsmanCenter != nil {
		query = query.Where(sq.Eq{"ombudsman_center": ombudsmanCenter})
	}
	if assignedOmbudsmanID, ok := filters["assigned_ombudsman_id"]; ok && assignedOmbudsmanID != nil {
		query = query.Where(sq.Eq{"assigned_ombudsman_id": assignedOmbudsmanID})
	}
	if admissible, ok := filters["admissible"]; ok && admissible != nil {
		query = query.Where(sq.Eq{"admissible": admissible})
	}
	if complaintCategory, ok := filters["complaint_category"]; ok && complaintCategory != nil {
		query = query.Where(sq.Eq{"complaint_category": complaintCategory})
	}
	if escalationToIRDAI, ok := filters["escalated_to_irdai"]; ok && escalationToIRDAI != nil {
		query = query.Where(sq.Eq{"escalated_to_irdai": escalationToIRDAI})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").From(ombudsmanComplaintTable)
	if status, ok := filters["status"]; ok && status != nil {
		countQuery = countQuery.Where(sq.Eq{"status": status})
	}
	if ombudsmanCenter, ok := filters["ombudsman_center"]; ok && ombudsmanCenter != nil {
		countQuery = countQuery.Where(sq.Eq{"ombudsman_center": ombudsmanCenter})
	}
	if assignedOmbudsmanID, ok := filters["assigned_ombudsman_id"]; ok && assignedOmbudsmanID != nil {
		countQuery = countQuery.Where(sq.Eq{"assigned_ombudsman_id": assignedOmbudsmanID})
	}
	if admissible, ok := filters["admissible"]; ok && admissible != nil {
		countQuery = countQuery.Where(sq.Eq{"admissible": admissible})
	}
	countQuery = countQuery.PlaceholderFormat(sq.Dollar)

	// Get total count
	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	query = query.OrderBy("created_at DESC").
		Limit(limit).
		Offset(skip).
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return results, 0, err
	}

	return results, totalCount, nil
}

// Update updates an ombudsman complaint
func (r *OmbudsmanComplaintRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable)

	// Build dynamic update
	for col, val := range updates {
		query = query.Set(col, val)
	}

	query = query.Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		if err == pgx.ErrNoRows {
			return result, fmt.Errorf("ombudsman complaint not found")
		}
		return result, err
	}

	return result, nil
}

// UpdateStatus updates the status of an ombudsman complaint
func (r *OmbudsmanComplaintRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// UpdateAdmissibility updates admissibility checks
// Reference: BR-CLM-OMB-001 (Admissibility checks: ₹50 lakh cap, 1-year limitation)
func (r *OmbudsmanComplaintRepository) UpdateAdmissibility(ctx context.Context, id string, admissible bool, reason *string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("admissible", admissible).
		Set("inadmissibility_reason", reason).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// AssignOmbudsman assigns an ombudsman to a complaint
// Reference: BR-CLM-OMB-002 (Jurisdiction), BR-CLM-OMB-003 (Conflict of interest)
func (r *OmbudsmanComplaintRepository) AssignOmbudsman(ctx context.Context, id string, ombudsmanID string, center string, jurisdictionBasis string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("assigned_ombudsman_id", ombudsmanID).
		Set("ombudsman_center", center).
		Set("jurisdiction_basis", jurisdictionBasis).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// RecordMediation records mediation attempt and outcome
// Reference: BR-CLM-OMB-004 (Mediation attempt)
func (r *OmbudsmanComplaintRepository) RecordMediation(ctx context.Context, id string, successful bool, terms *string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("mediation_attempted", true).
		Set("mediation_successful", &successful).
		Set("mediation_terms", terms).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// IssueRecommendation issues a recommendation
func (r *OmbudsmanComplaintRepository) IssueRecommendation(ctx context.Context, id string, recommendationDate time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("recommendation_issued", true).
		Set("recommendation_date", recommendationDate).
		Set("status", "RECOMMENDATION").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// IssueAward issues an award
// Reference: BR-CLM-OMB-005 (Award amount cap: ₹50 lakh)
func (r *OmbudsmanComplaintRepository) IssueAward(ctx context.Context, id string, awardNumber string, awardAmount float64, awardDate time.Time, complianceDueDate time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("award_issued", true).
		Set("award_number", awardNumber).
		Set("award_amount", awardAmount).
		Set("award_date", awardDate).
		Set("compliance_due_date", complianceDueDate).
		Set("status", "AWARD").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// RecordCompliance records compliance with the award
// Reference: BR-CLM-OMB-006 (30-day compliance timeline)
func (r *OmbudsmanComplaintRepository) RecordCompliance(ctx context.Context, id string, complianceStatus string, complianceDate time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("compliance_status", complianceStatus).
		Set("compliance_date", complianceDate).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// EscalateToIRDAI escalates the complaint to IRDAI for non-compliance
// Reference: BR-CLM-OMB-006 (30-day compliance timeline)
func (r *OmbudsmanComplaintRepository) EscalateToIRDAI(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("escalated_to_irdai", true).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// CloseComplaint closes the complaint
func (r *OmbudsmanComplaintRepository) CloseComplaint(ctx context.Context, id string, closureDate time.Time, archivalDate *time.Time, retentionPeriod *int) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(ombudsmanComplaintTable).
		Set("status", "CLOSED").
		Set("closure_date", closureDate).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

	if archivalDate != nil {
		query = query.Set("archival_date", archivalDate)
	}
	if retentionPeriod != nil {
		query = query.Set("retention_period_years", retentionPeriod)
	}

	query = query.PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// GetPendingAdmissibilityChecks retrieves complaints pending admissibility review
// Reference: BR-CLM-OMB-001 (Admissibility checks)
func (r *OmbudsmanComplaintRepository) GetPendingAdmissibilityChecks(ctx context.Context) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"status": "SUBMITTED"}).
		Where(sq.Or{
			sq.Eq{"admissible": nil},
			sq.Eq{"wait_period_completed": false},
			sq.Eq{"limitation_period_valid": false},
		}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetPendingMediation retrieves complaints pending mediation
// Reference: BR-CLM-OMB-004 (Mediation attempt)
func (r *OmbudsmanComplaintRepository) GetPendingMediation(ctx context.Context) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"status": "UNDER_REVIEW"}).
		Where(sq.Eq{"mediation_attempted": false}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetPendingAwards retrieves complaints with recommendations but no award yet
func (r *OmbudsmanComplaintRepository) GetPendingAwards(ctx context.Context) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"status": "RECOMMENDATION"}).
		Where(sq.Eq{"award_issued": false}).
		OrderBy("recommendation_date ASC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetOverdueCompliance retrieves complaints with overdue compliance
// Reference: BR-CLM-OMB-006 (30-day compliance timeline)
func (r *OmbudsmanComplaintRepository) GetOverdueCompliance(ctx context.Context) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"status": "AWARD"}).
		Where(sq.Eq{"award_issued": true}).
		Where(sq.NotEq{"compliance_status": "COMPLIED"}).
		Where(sq.Lt{"compliance_due_date": time.Now()}).
		Where(sq.Eq{"escalated_to_irdai": false}).
		OrderBy("compliance_due_date ASC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetComplaintsByOmbudsman retrieves complaints assigned to a specific ombudsman
func (r *OmbudsmanComplaintRepository) GetComplaintsByOmbudsman(ctx context.Context, ombudsmanID string, status *string) ([]domain.OmbudsmanComplaint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"assigned_ombudsman_id": ombudsmanID})

	if status != nil {
		query = query.Where(sq.Eq{"status": *status})
	}

	query = query.OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.OmbudsmanComplaint])
	if err != nil {
		return result, err
	}
	return result, nil
}

// CheckAdmissibility validates admissibility criteria
// Reference: BR-CLM-OMB-001 (Admissibility checks: ₹50 lakh cap, 1-year limitation)
func (r *OmbudsmanComplaintRepository) CheckAdmissibility(ctx context.Context, claimValue float64, representationDate *time.Time, parallelLitigation bool) (admissible bool, reason *string, err error) {
	// Check claim value cap (₹50 lakh)
	if claimValue > 5000000 {
		falseVal := false
		reasonStr := "Claim value exceeds ₹50 lakh cap"
		return falseVal, &reasonStr, nil
	}

	// Check 1-year limitation period
	if representationDate != nil {
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		if representationDate.Before(oneYearAgo) {
			falseVal := false
			reasonStr := "Representation to insurer was more than 1 year ago"
			return falseVal, &reasonStr, nil
		}
	}

	// Check parallel litigation
	if parallelLitigation {
		falseVal := false
		reasonStr := "Parallel litigation is pending"
		return falseVal, &reasonStr, nil
	}

	// All checks passed
	trueVal := true
	return trueVal, nil, nil
}

// GenerateComplaintNumber generates a unique complaint number
// Format: OMB{YYYY}{DDDD}
func (r *OmbudsmanComplaintRepository) GenerateComplaintNumber(ctx context.Context) (string, error) {
	year := time.Now().Format("2006")

	// Get count of complaints in current year
	countQuery := sq.Select("COUNT(*)").
		From(ombudsmanComplaintTable).
		Where(sq.Like{"complaint_number": "OMB" + year + "%"}).
		PlaceholderFormat(sq.Dollar)

	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	count, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return "", err
	}

	// Generate sequential number (4 digits, zero-padded)
	seqNum := int(count) + 1
	complaintNumber := fmt.Sprintf("OMB%s%04d", year, seqNum)

	return complaintNumber, nil
}

// Delete soft deletes an ombudsman complaint (sets deletion timestamp if needed)
func (r *OmbudsmanComplaintRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Delete(ombudsmanComplaintTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Delete(ctx, r.db, query)
	return err
}

// GetComplaintStats retrieves statistics for ombudsman complaints
func (r *OmbudsmanComplaintRepository) GetComplaintStats(ctx context.Context) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select(
		"COUNT(*) as total",
		"COUNT(*) FILTER (WHERE status = 'SUBMITTED') as submitted",
		"COUNT(*) FILTER (WHERE status = 'UNDER_REVIEW') as under_review",
		"COUNT(*) FILTER (WHERE status = 'MEDIATION') as mediation",
		"COUNT(*) FILTER (WHERE status = 'RECOMMENDATION') as recommendation",
		"COUNT(*) FILTER (WHERE status = 'AWARD') as award",
		"COUNT(*) FILTER (WHERE status = 'CLOSED') as closed",
		"COUNT(*) FILTER (WHERE admissible = true) as admissible",
		"COUNT(*) FILTER (WHERE admissible = false) as inadmissible",
		"COUNT(*) FILTER (WHERE escalated_to_irdai = true) as escalated",
	).
		From(ombudsmanComplaintTable).
		PlaceholderFormat(sq.Dollar)

	type Stats struct {
		Total        int64 `db:"total"`
		Submitted    int64 `db:"submitted"`
		UnderReview  int64 `db:"under_review"`
		Mediation    int64 `db:"mediation"`
		Recommendation int64 `db:"recommendation"`
		Award        int64 `db:"award"`
		Closed       int64 `db:"closed"`
		Admissible   int64 `db:"admissible"`
		Inadmissible int64 `db:"inadmissible"`
		Escalated    int64 `db:"escalated"`
	}

	stats, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[Stats])
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"total":         stats.Total,
		"submitted":     stats.Submitted,
		"under_review":  stats.UnderReview,
		"mediation":     stats.Mediation,
		"recommendation": stats.Recommendation,
		"award":         stats.Award,
		"closed":        stats.Closed,
		"admissible":    stats.Admissible,
		"inadmissible":  stats.Inadmissible,
		"escalated":     stats.Escalated,
	}, nil
}

// IsDuplicateComplaint checks if a duplicate complaint exists for the same claim/policy
func (r *OmbudsmanComplaintRepository) IsDuplicateComplaint(ctx context.Context, claimID *string, policyID string, complainantName string, description string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("COUNT(*)").
		From(ombudsmanComplaintTable).
		Where(sq.Eq{"policy_id": policyID}).
		Where(sq.Eq{"complainant_name": complainantName})

	if claimID != nil {
		query = query.Where(sq.Eq{"claim_id": *claimID})
	}

	// Check if complaint was submitted in last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	query = query.Where(sq.Gt{"created_at": thirtyDaysAgo}).
		PlaceholderFormat(sq.Dollar)

	count, err := dblib.SelectOne(ctx, r.db, query, pgx.RowTo[int64])
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
