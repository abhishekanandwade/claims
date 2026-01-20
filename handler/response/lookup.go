package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== LOOKUP & REFERENCE RESPONSE DTOs ====================

// ClaimantRelationshipResponse represents a claimant relationship
// GET /lookup/claimant-relationships
type ClaimantRelationshipResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Relationships             []RelationshipItem `json:"relationships"`
}

// RelationshipItem represents a single relationship item
type RelationshipItem struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// DeathTypesResponse represents death types list
// GET /lookup/death-types
type DeathTypesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	DeathTypes                []DeathTypeItem `json:"death_types"`
}

// DeathTypeItem represents a single death type item
type DeathTypeItem struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	RequiresFIR bool   `json:"requires_fir"` // Whether FIR is required for this death type
}

// DocumentTypesResponse represents document types for claim
// GET /lookup/document-types
type DocumentTypesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	DocumentTypes             []DocumentTypeItem `json:"document_types"`
}

// DocumentTypeItem represents a single document type item
type DocumentTypeItem struct {
	DocumentCode      string `json:"document_code"`
	DocumentName      string `json:"document_name"`
	Required          bool   `json:"required"`
	DocumentFormat    string `json:"document_format"` // PDF, JPEG, PNG
	MaxFileSize       int    `json:"max_file_size"`   // in MB
	Description       string `json:"description,omitempty"`
	ApplicableFor     string `json:"applicable_for"` // ALL, ACCIDENTAL, UNNATURAL, SUICIDE
	IsMandatory       bool   `json:"is_mandatory"`
	MandatoryForStage string `json:"mandatory_for_stage"` // REGISTRATION, APPROVAL, DISBURSEMENT
}

// RejectionReasonsResponse represents rejection reasons list
// GET /lookup/rejection-reasons
type RejectionReasonsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	RejectionReasons          []RejectionReasonItem `json:"rejection_reasons"`
}

// RejectionReasonItem represents a single rejection reason item
type RejectionReasonItem struct {
	ReasonCode         string `json:"reason_code"`
	Reason             string `json:"reason"`
	Description        string `json:"description,omitempty"`
	AppealAllowed      bool   `json:"appeal_allowed"`                // Whether appeal is allowed for this rejection
	AppellateAuthority string `json:"appellate_authority,omitempty"` // Authority for appeal
	Category           string `json:"category"`                      // DOCUMENTARY, MEDICAL, LEGAL, TECHNICAL, FRAUD
}

// InvestigationOfficersResponse represents investigation officers list
// GET /lookup/investigation-officers
type InvestigationOfficersResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Officers                  []InvestigationOfficerItem `json:"officers"`
}

// InvestigationOfficerItem represents a single investigation officer item
type InvestigationOfficerItem struct {
	OfficerID      string  `json:"officer_id"`
	OfficerName    string  `json:"officer_name"`
	Rank           string  `json:"rank"`
	Jurisdiction   string  `json:"jurisdiction"`
	Phone          string  `json:"phone,omitempty"`
	Email          string  `json:"email,omitempty"`
	Available      bool    `json:"available"`                // Whether officer is available for assignment
	ActiveCases    int     `json:"active_cases"`             // Number of active cases
	Experience     float64 `json:"experience"`               // Years of experience
	Specialization string  `json:"specialization,omitempty"` // e.g., "FRAUD", "MEDICAL", "LEGAL"
}

// ApproversListResponse represents eligible approvers for claim
// GET /lookup/approvers
type ApproversListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Approvers                 []ApproverItem `json:"approvers"`
}

// ApproverItem represents a single approver item
type ApproverItem struct {
	ApproverID      string  `json:"approver_id"`
	ApproverName    string  `json:"approver_name"`
	Designation     string  `json:"designation"`
	Department      string  `json:"department"`
	Location        string  `json:"location"`
	MinAmount       float64 `json:"min_amount"`       // Minimum claim amount for this approver
	MaxAmount       float64 `json:"max_amount"`       // Maximum claim amount for this approver
	ApprovalLevel   string  `json:"approval_level"`   // LEVEL_1, LEVEL_2, LEVEL_3, LEVEL_4
	Available       bool    `json:"available"`        // Whether approver is available
	ActiveApprovals int     `json:"active_approvals"` // Number of active approvals
	Phone           string  `json:"phone,omitempty"`
	Email           string  `json:"email,omitempty"`
}

// PaymentModesResponse represents available payment modes
// GET /lookup/payment-modes
type PaymentModesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PaymentModes              []PaymentModeItem `json:"payment_modes"`
}

// PaymentModeItem represents a single payment mode item
type PaymentModeItem struct {
	PaymentMode            string  `json:"payment_mode"` // NEFT, POSB, CHEQUE
	Description            string  `json:"description"`
	Enabled                bool    `json:"enabled"`                  // Whether payment mode is currently enabled
	MinAmount              float64 `json:"min_amount"`               // Minimum amount for this payment mode
	MaxAmount              float64 `json:"max_amount"`               // Maximum amount for this payment mode (0 for no limit)
	ProcessingTime         string  `json:"processing_time"`          // e.g., "1-2 working days"
	RequiresBankValidation bool    `json:"requires_bank_validation"` // Whether bank validation is required
	Charges                float64 `json:"charges"`                  // Processing charges, if any
}

// ApprovalHierarchyResponse represents approval hierarchy and financial limits
// GET /approvers/financial-limits
type ApprovalHierarchyResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Hierarchy                 []ApprovalLevelItem `json:"hierarchy"`
}

// ApprovalLevelItem represents a single approval level item
type ApprovalLevelItem struct {
	Designation      string  `json:"designation"`
	ApprovalLevel    string  `json:"approval_level"`    // LEVEL_1, LEVEL_2, LEVEL_3, LEVEL_4
	MinAmount        float64 `json:"min_amount"`        // Minimum claim amount for this level
	MaxAmount        float64 `json:"max_amount"`        // Maximum claim amount for this level (0 for no limit)
	RequiresApproval bool    `json:"requires_approval"` // Whether approval is required at this level
	CanOverride      bool    `json:"can_override"`      // Whether can override lower level decisions
	EscalationOrder  int     `json:"escalation_order"`  // Order for escalation (1 = lowest, 4 = highest)
}

// ==================== HELPER FUNCTIONS ====================

// NewClaimantRelationshipsResponse creates a new claimant relationships response
func NewClaimantRelationshipsResponse(relationships []RelationshipItem) *ClaimantRelationshipResponse {
	return &ClaimantRelationshipResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Claimant relationships retrieved successfully",
		},
		Relationships: relationships,
	}
}

// NewDeathTypesResponse creates a new death types response
func NewDeathTypesResponse(deathTypes []DeathTypeItem) *DeathTypesResponse {
	return &DeathTypesResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Death types retrieved successfully",
		},
		DeathTypes: deathTypes,
	}
}

// NewDocumentTypesResponse creates a new document types response
func NewDocumentTypesResponse(documentTypes []DocumentTypeItem) *DocumentTypesResponse {
	return &DocumentTypesResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Document types retrieved successfully",
		},
		DocumentTypes: documentTypes,
	}
}

// NewRejectionReasonsResponse creates a new rejection reasons response
func NewRejectionReasonsResponse(rejectionReasons []RejectionReasonItem) *RejectionReasonsResponse {
	return &RejectionReasonsResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Rejection reasons retrieved successfully",
		},
		RejectionReasons: rejectionReasons,
	}
}

// NewInvestigationOfficersResponse creates a new investigation officers response
func NewInvestigationOfficersResponse(officers []InvestigationOfficerItem) *InvestigationOfficersResponse {
	return &InvestigationOfficersResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Investigation officers retrieved successfully",
		},
		Officers: officers,
	}
}

// NewApproversListResponse creates a new approvers list response
func NewApproversListResponse(approvers []ApproverItem) *ApproversListResponse {
	return &ApproversListResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Approvers retrieved successfully",
		},
		Approvers: approvers,
	}
}

// NewPaymentModesResponse creates a new payment modes response
func NewPaymentModesResponse(paymentModes []PaymentModeItem) *PaymentModesResponse {
	return &PaymentModesResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Payment modes retrieved successfully",
		},
		PaymentModes: paymentModes,
	}
}

// NewApprovalHierarchyResponse creates a new approval hierarchy response
func NewApprovalHierarchyResponse(hierarchy []ApprovalLevelItem) *ApprovalHierarchyResponse {
	return &ApprovalHierarchyResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Approval hierarchy retrieved successfully",
		},
		Hierarchy: hierarchy,
	}
}
