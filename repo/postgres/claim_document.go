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

// ClaimDocumentRepository handles claim document data operations
type ClaimDocumentRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClaimDocumentRepository creates a new claim document repository
func NewClaimDocumentRepository(db *dblib.DB, cfg *config.Config) *ClaimDocumentRepository {
	return &ClaimDocumentRepository{
		db:  db,
		cfg: cfg,
	}
}

const claimDocumentTable = "claim_documents"

// Create inserts a new claim document
// Reference: E-CLM-DC-002, seed/db/claims_database_schema.sql:203-244
func (r *ClaimDocumentRepository) Create(ctx context.Context, data domain.ClaimDocument) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Insert(claimDocumentTable).
		Columns(
			"claim_id", "document_type", "document_name", "document_url",
			"ecms_reference_id", "file_size", "file_hash", "content_type",
			"is_mandatory", "uploaded_by", "uploaded_at",
			"virus_scanned", "virus_scan_status",
			"verified", "verified_by", "verified_at", "verification_remarks",
			"ocr_extracted_data", "ocr_confidence_score",
			"created_at", "updated_at",
		).
		Values(
			data.ClaimID, data.DocumentType, data.DocumentName, data.DocumentURL,
			data.ECMSReferenceID, data.FileSize, data.FileHash, data.ContentType,
			data.IsMandatory, data.UploadedBy, data.UploadedAt,
			data.VirusScanned, data.VirusScanStatus,
			data.Verified, data.VerifiedBy, data.VerifiedAt, data.VerificationRemarks,
			data.OCRExtractedData, data.OCRConfidenceScore,
			data.CreatedAt, data.UpdatedAt,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.InsertReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	return result, err
}

// CreateBatch inserts multiple claim documents by calling Create in a loop
// Note: For true batch operations, use database transactions or batch APIs
func (r *ClaimDocumentRepository) CreateBatch(ctx context.Context, documents []domain.ClaimDocument) ([]domain.ClaimDocument, error) {
	results := make([]domain.ClaimDocument, 0, len(documents))

	for _, doc := range documents {
		result, err := r.Create(ctx, doc)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// FindByID retrieves a claim document by ID
func (r *ClaimDocumentRepository) FindByID(ctx context.Context, documentID string) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.Eq{"id": documentID}).
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.SelectOne(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindByClaimID retrieves all documents for a specific claim
// Returns documents ordered by upload date (newest first)
func (r *ClaimDocumentRepository) FindByClaimID(ctx context.Context, claimID string) ([]domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.Eq{"claim_id": claimID}).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindByClaimIDAndType retrieves documents for a claim filtered by document type
func (r *ClaimDocumentRepository) FindByClaimIDAndType(ctx context.Context, claimID string, documentType string) ([]domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"document_type": documentType},
		}).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// List retrieves claim documents with filters and pagination
func (r *ClaimDocumentRepository) List(ctx context.Context, filters map[string]interface{}, skip, limit int64, orderBy, sortType string) ([]domain.ClaimDocument, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query with filters
	baseQuery := sq.Select("*").From(claimDocumentTable).PlaceholderFormat(sq.Dollar)

	// Apply filters if provided
	for key, value := range filters {
		baseQuery = baseQuery.Where(sq.Eq{key: value})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").From(claimDocumentTable).PlaceholderFormat(sq.Dollar)
	for key, value := range filters {
		countQuery = countQuery.Where(sq.Eq{key: value})
	}

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination and sorting
	query := baseQuery.OrderBy(orderBy+" "+sortType).
		Limit(uint64(limit)).
		Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// Update updates claim document fields
func (r *ClaimDocumentRepository) Update(ctx context.Context, documentID string, updates map[string]interface{}) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimDocumentTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": documentID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	// Apply updates dynamically
	for key, value := range updates {
		query = query.Set(key, value)
	}

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	return result, err
}

// UpdateVerification updates document verification status
// Reference: BR-CLM-DC-006 (Document verification workflow)
func (r *ClaimDocumentRepository) UpdateVerification(ctx context.Context, documentID string, verified bool, verifiedBy string, remarks *string) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	now := time.Now()

	query := sq.Update(claimDocumentTable).
		Set("verified", verified).
		Set("verified_by", verifiedBy).
		Set("verified_at", now).
		Set("verification_remarks", remarks).
		Set("updated_at", now).
		Where(sq.Eq{"id": documentID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	return result, err
}

// UpdateVirusScan updates virus scan status
// Reference: BR-CLM-DC-011 (Virus scanning requirement)
func (r *ClaimDocumentRepository) UpdateVirusScan(ctx context.Context, documentID string, scanned bool, status *string) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimDocumentTable).
		Set("virus_scanned", scanned).
		Set("virus_scan_status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": documentID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	return result, err
}

// UpdateOCRData updates OCR extracted data and confidence score
func (r *ClaimDocumentRepository) UpdateOCRData(ctx context.Context, documentID string, ocrData map[string]interface{}, confidenceScore *float64) (domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Update(claimDocumentTable).
		Set("ocr_extracted_data", ocrData).
		Set("ocr_confidence_score", confidenceScore).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": documentID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	result, err := dblib.UpdateReturning(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	return result, err
}

// MarkAsDeleted soft deletes a claim document
func (r *ClaimDocumentRepository) MarkAsDeleted(ctx context.Context, documentID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	now := time.Now()

	query := sq.Update(claimDocumentTable).
		Set("deleted_at", now).
		Set("updated_at", now).
		Where(sq.Eq{"id": documentID}).
		PlaceholderFormat(sq.Dollar)

	_, err := dblib.Update(ctx, r.db, query)
	return err
}

// GetUnverifiedDocuments retrieves documents pending verification
// Reference: BR-CLM-DC-006 (Document verification queue)
func (r *ClaimDocumentRepository) GetUnverifiedDocuments(ctx context.Context, filters map[string]interface{}, skip, limit int64) ([]domain.ClaimDocument, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	// Build base query for unverified documents
	baseQuery := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.Eq{"verified": false}).
		Where(sq.Eq{"virus_scanned": true}). // Only show scanned documents
		PlaceholderFormat(sq.Dollar)

	// Apply additional filters
	for key, value := range filters {
		baseQuery = baseQuery.Where(sq.Eq{key: value})
	}

	// Count query
	countQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.Eq{"verified": false}).
		Where(sq.Eq{"virus_scanned": true}).
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
		OrderBy("uploaded_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip))

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// GetDocumentsPendingVirusScan retrieves documents pending virus scan
// Reference: BR-CLM-DC-011 (Virus scanning requirement)
func (r *ClaimDocumentRepository) GetDocumentsPendingVirusScan(ctx context.Context, skip, limit int64) ([]domain.ClaimDocument, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	query := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.Eq{"virus_scanned": false}).
		OrderBy("uploaded_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	countQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.Eq{"virus_scanned": false}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, countQuery, pgx.RowTo[int64])
	if err != nil {
		return nil, 0, err
	}

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// GetMandatoryDocuments retrieves mandatory documents for a claim
func (r *ClaimDocumentRepository) GetMandatoryDocuments(ctx context.Context, claimID string) ([]domain.ClaimDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	query := sq.Select("*").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"is_mandatory": true},
		}).
		OrderBy("document_type ASC").
		PlaceholderFormat(sq.Dollar)

	results, err := dblib.SelectRows(ctx, r.db, query, pgx.RowToStructByPos[domain.ClaimDocument])
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CheckDocumentCompleteness checks if all mandatory documents are uploaded for a claim
// Returns (isComplete, uploadedCount, requiredCount, error)
// Reference: BR-CLM-DC-006 (Document completeness check)
func (r *ClaimDocumentRepository) CheckDocumentCompleteness(ctx context.Context, claimID string) (bool, int64, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// Count uploaded mandatory documents
	uploadedQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"is_mandatory": true},
		}).
		PlaceholderFormat(sq.Dollar)

	uploadedCount, err := dblib.SelectOne(ctx, r.db, uploadedQuery, pgx.RowTo[int64])
	if err != nil {
		return false, 0, 0, err
	}

	// Get required document types from document_checklist table
	requiredQuery := sq.Select("COUNT(*)").
		From("document_checklist").
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"is_mandatory": true},
		}).
		PlaceholderFormat(sq.Dollar)

	requiredCount, err := dblib.SelectOne(ctx, r.db, requiredQuery, pgx.RowTo[int64])
	if err != nil {
		return false, 0, 0, err
	}

	isComplete := uploadedCount == requiredCount && requiredCount > 0

	return isComplete, uploadedCount, requiredCount, nil
}

// GetDocumentStats retrieves document statistics for a claim
// Returns (totalDocuments, verifiedCount, pendingCount, mandatoryCount)
func (r *ClaimDocumentRepository) GetDocumentStats(ctx context.Context, claimID string) (int64, int64, int64, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	// Query total documents
	totalQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.Eq{"claim_id": claimID}).
		PlaceholderFormat(sq.Dollar)

	totalCount, err := dblib.SelectOne(ctx, r.db, totalQuery, pgx.RowTo[int64])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Query verified documents
	verifiedQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"verified": true},
		}).
		PlaceholderFormat(sq.Dollar)

	verifiedCount, err := dblib.SelectOne(ctx, r.db, verifiedQuery, pgx.RowTo[int64])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Query pending (unverified) documents
	pendingQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"verified": false},
		}).
		PlaceholderFormat(sq.Dollar)

	pendingCount, err := dblib.SelectOne(ctx, r.db, pendingQuery, pgx.RowTo[int64])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Query mandatory documents
	mandatoryQuery := sq.Select("COUNT(*)").
		From(claimDocumentTable).
		Where(sq.And{
			sq.Eq{"claim_id": claimID},
			sq.Eq{"is_mandatory": true},
		}).
		PlaceholderFormat(sq.Dollar)

	mandatoryCount, err := dblib.SelectOne(ctx, r.db, mandatoryQuery, pgx.RowTo[int64])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return totalCount, verifiedCount, pendingCount, mandatoryCount, nil
}

// BatchUpdateVerification updates verification status for multiple documents
// Note: For true batch operations, use database transactions or batch APIs
func (r *ClaimDocumentRepository) BatchUpdateVerification(ctx context.Context, documentIDs []string, verified bool, verifiedBy string) ([]domain.ClaimDocument, error) {
	results := make([]domain.ClaimDocument, 0, len(documentIDs))

	for _, docID := range documentIDs {
		// Use empty string for nil remarks
		var remarks *string
		result, err := r.UpdateVerification(ctx, docID, verified, verifiedBy, remarks)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
