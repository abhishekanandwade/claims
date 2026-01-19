package domain

import (
	"time"
)

// ClaimDocument represents documents uploaded for claim processing
// Table: claim_documents (partitioned by uploaded_at)
// Reference: E-CLM-DC-002, seed/db/claims_database_schema.sql:203-244
type ClaimDocument struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// References
	ClaimID string `db:"claim_id" json:"claim_id"`

	// Document Information
	DocumentType string `db:"document_type" json:"document_type"` // BR-CLM-DC-013/014/015
	DocumentName string `db:"document_name" json:"document_name"`
	DocumentURL  string `db:"document_url" json:"document_url"`

	// ECMS Integration
	ECMSReferenceID *string `db:"ecms_reference_id" json:"ecms_reference_id,omitempty"`

	// File Metadata
	FileSize int     `db:"file_size" json:"file_size"`
	FileHash *string `db:"file_hash" json:"file_hash,omitempty"`
	ContentType *string `db:"content_type" json:"content_type,omitempty"`

	// Document Properties
	IsMandatory bool `db:"is_mandatory" json:"is_mandatory"`

	// Upload Information
	UploadedBy string    `db:"uploaded_by" json:"uploaded_by"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`

	// Virus Scanning
	VirusScanned     bool    `db:"virus_scanned" json:"virus_scanned"`
	VirusScanStatus  *string `db:"virus_scan_status" json:"virus_scan_status,omitempty"`

	// Verification Information
	Verified            bool       `db:"verified" json:"verified"`
	VerifiedBy          *string    `db:"verified_by" json:"verified_by,omitempty"`
	VerifiedAt          *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	VerificationRemarks *string    `db:"verification_remarks" json:"verification_remarks,omitempty"`

	// OCR Information
	OCRExtractedData    map[string]interface{} `db:"ocr_extracted_data" json:"ocr_extracted_data,omitempty"`
	OCRConfidenceScore  *float64               `db:"ocr_confidence_score" json:"ocr_confidence_score,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// TableName returns the database table name for ClaimDocument
func (ClaimDocument) TableName() string {
	return "claim_documents"
}
