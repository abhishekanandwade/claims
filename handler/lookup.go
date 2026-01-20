package handler

import (
	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
)

// LookupHandler handles lookup and reference data HTTP requests
type LookupHandler struct {
	*serverHandler.Base
}

// NewLookupHandler creates a new lookup handler
func NewLookupHandler() *LookupHandler {
	base := serverHandler.New("Lookup").
		SetPrefix("/v1").
		AddPrefix("")
	return &LookupHandler{
		Base: base,
	}
}

// Routes defines all routes for this handler
func (h *LookupHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Claimant Relationships (1 endpoint)
		serverRoute.GET("/lookup/claimant-relationships", h.GetClaimantRelationships).Name("Get Claimant Relationships"),

		// Death Types (1 endpoint)
		serverRoute.GET("/lookup/death-types", h.GetDeathTypes).Name("Get Death Types"),

		// Document Types (1 endpoint)
		serverRoute.GET("/lookup/document-types", h.GetDocumentTypes).Name("Get Document Types"),

		// Rejection Reasons (1 endpoint)
		serverRoute.GET("/lookup/rejection-reasons", h.GetRejectionReasons).Name("Get Rejection Reasons"),

		// Investigation Officers (1 endpoint)
		serverRoute.GET("/lookup/investigation-officers", h.GetInvestigationOfficers).Name("Get Investigation Officers"),

		// Approvers (1 endpoint)
		serverRoute.GET("/lookup/approvers", h.GetApproversList).Name("Get Approvers List"),

		// Payment Modes (1 endpoint)
		serverRoute.GET("/lookup/payment-modes", h.GetPaymentModes).Name("Get Payment Modes"),

		// Approval Hierarchy (1 endpoint)
		serverRoute.GET("/approvers/financial-limits", h.GetApprovalHierarchy).Name("Get Approval Hierarchy"),

		// Additional lookup endpoints (4 endpoints)
		serverRoute.GET("/lookup/banks", h.GetBanksList).Name("Get Banks List"),

		serverRoute.GET("/lookup/branches", h.GetBranchesList).Name("Get Branches List"),

		serverRoute.GET("/lookup/states", h.GetStatesList).Name("Get States List"),

		serverRoute.GET("/lookup/districts", h.GetDistrictsList).Name("Get Districts List"),
	}
}

// ==================== HANDLER METHODS ====================

// GetClaimantRelationships handles GET /lookup/claimant-relationships
// Reference: Master data for claimant relationships
func (h *LookupHandler) GetClaimantRelationships(sctx *serverRoute.Context, req *struct{}) (*resp.ClaimantRelationshipResponse, error) {
	log.Info(sctx.Ctx, "GetClaimantRelationships: Fetching claimant relationships")

	// TODO: Integrate with master data service or database lookup table
	// For now, return static data as per business requirements
	relationships := []resp.RelationshipItem{
		{Code: "SELF", Description: "Self"},
		{Code: "SPOUSE", Description: "Spouse"},
		{Code: "SON", Description: "Son"},
		{Code: "DAUGHTER", Description: "Daughter"},
		{Code: "FATHER", Description: "Father"},
		{Code: "MOTHER", Description: "Mother"},
		{Code: "BROTHER", Description: "Brother"},
		{Code: "SISTER", Description: "Sister"},
		{Code: "GRANDFATHER", Description: "Grandfather"},
		{Code: "GRANDMOTHER", Description: "Grandmother"},
		{Code: "LEGAL_HEIR", Description: "Legal Heir"},
		{Code: "ASSIGNEE", Description: "Assignee"},
		{Code: "GUARDIAN", Description: "Guardian"},
		{Code: "TRUSTEE", Description: "Trustee"},
		{Code: "NOMINEE", Description: "Nominee"},
		{Code: "OTHER", Description: "Other"},
	}

	return resp.NewClaimantRelationshipsResponse(relationships), nil
}

// GetDeathTypes handles GET /lookup/death-types
// Reference: Master data for death types
func (h *LookupHandler) GetDeathTypes(sctx *serverRoute.Context, req *struct{}) (*resp.DeathTypesResponse, error) {
	log.Info(sctx.Ctx, "GetDeathTypes: Fetching death types")

	// TODO: Integrate with master data service or database lookup table
	// Death types with FIR requirement flags
	deathTypes := []resp.DeathTypeItem{
		{Code: "NATURAL", Description: "Natural Death", RequiresFIR: false},
		{Code: "ACCIDENTAL", Description: "Accidental Death", RequiresFIR: true},
		{Code: "UNNATURAL", Description: "Unnatural Death", RequiresFIR: true},
		{Code: "SUICIDE", Description: "Suicide", RequiresFIR: true},
		{Code: "HOMICIDE", Description: "Homicide", RequiresFIR: true},
	}

	return resp.NewDeathTypesResponse(deathTypes), nil
}

// GetDocumentTypes handles GET /lookup/document-types
// Reference: DFC-001 (Dynamic document checklist based on claim type, death type, nomination status)
func (h *LookupHandler) GetDocumentTypes(sctx *serverRoute.Context, req *GetDocumentTypesRequest) (*resp.DocumentTypesResponse, error) {
	log.Info(sctx.Ctx, "GetDocumentTypes: Fetching document types for claim_type=%s, death_type=%v, nomination_status=%v",
		req.ClaimType, req.DeathType, req.NominationStatus)

	// TODO: Integrate with document_checklist table or master data service
	// Return dynamic document types based on claim characteristics
	var documentTypes []resp.DocumentTypeItem

	switch req.ClaimType {
	case "DEATH":
		// Base documents for all death claims
		documentTypes = append(documentTypes, []resp.DocumentTypeItem{
			{
				DocumentCode:      "DEATH_CERTIFICATE",
				DocumentName:      "Death Certificate",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Original death certificate issued by municipal authority",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "CLAIM_FORM",
				DocumentName:      "Claim Form",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Duly filled claim form",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "ID_PROOF",
				DocumentName:      "Claimant ID Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Aadhaar Card / PAN Card / Voter ID",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
		}...)

		// Additional documents for accidental/unnatural death
		if req.DeathType != nil && (*req.DeathType == "ACCIDENTAL" || *req.DeathType == "UNNATURAL" || *req.DeathType == "HOMICIDE") {
			documentTypes = append(documentTypes, []resp.DocumentTypeItem{
				{
					DocumentCode:      "FIR_COPY",
					DocumentName:      "FIR Copy",
					Required:          true,
					DocumentFormat:    "PDF",
					MaxFileSize:       5,
					Description:       "First Information Report filed with police",
					ApplicableFor:     "ACCIDENTAL",
					IsMandatory:       true,
					MandatoryForStage: "REGISTRATION",
				},
				{
					DocumentCode:      "POST_MORTEM_REPORT",
					DocumentName:      "Post Mortem Report",
					Required:          true,
					DocumentFormat:    "PDF",
					MaxFileSize:       5,
					Description:       "Autopsy report if post mortem conducted",
					ApplicableFor:     "ACCIDENTAL",
					IsMandatory:       false,
					MandatoryForStage: "REGISTRATION",
				},
				{
					DocumentCode:      "POLICE_REPORT",
					DocumentName:      "Police Investigation Report",
					Required:          false,
					DocumentFormat:    "PDF",
					MaxFileSize:       10,
					Description:       "Final police report with investigation details",
					ApplicableFor:     "ACCIDENTAL",
					IsMandatory:       false,
					MandatoryForStage: "APPROVAL",
				},
			}...)
		}

		// Suicide-specific documents
		if req.DeathType != nil && *req.DeathType == "SUICIDE" {
			documentTypes = append(documentTypes, resp.DocumentTypeItem{
				DocumentCode:      "SUICIDE_NOTE",
				DocumentName:      "Suicide Note (if available)",
				Required:          false,
				DocumentFormat:    "PDF,JPEG,PNG",
				MaxFileSize:       5,
				Description:       "Any suicide note or last letter",
				ApplicableFor:     "SUICIDE",
				IsMandatory:       false,
				MandatoryForStage: "REGISTRATION",
			})
		}

		// Nomination-specific documents
		if req.NominationStatus != nil && *req.NominationStatus == "NOT_NOMINATED" {
			documentTypes = append(documentTypes, []resp.DocumentTypeItem{
				{
					DocumentCode:      "LEGAL_HEIR_CERTIFICATE",
					DocumentName:      "Legal Heir Certificate",
					Required:          true,
					DocumentFormat:    "PDF",
					MaxFileSize:       5,
					Description:       "Legal heir certificate from competent authority",
					ApplicableFor:     "ALL",
					IsMandatory:       true,
					MandatoryForStage: "REGISTRATION",
				},
				{
					DocumentCode:      "SUCCESSION_CERTIFICATE",
					DocumentName:      "Succession Certificate",
					Required:          false,
					DocumentFormat:    "PDF",
					MaxFileSize:       5,
					Description:       "Succession certificate from court",
					ApplicableFor:     "ALL",
					IsMandatory:       false,
					MandatoryForStage: "APPROVAL",
				},
			}...)
		}

	case "MATURITY":
		documentTypes = []resp.DocumentTypeItem{
			{
				DocumentCode:      "CLAIM_FORM",
				DocumentName:      "Maturity Claim Form",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Duly filled maturity claim form",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "POLICY_BOND",
				DocumentName:      "Original Policy Bond",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Original policy bond document",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "DISBURSEMENT",
			},
			{
				DocumentCode:      "ID_PROOF",
				DocumentName:      "Policyholder ID Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Aadhaar Card / PAN Card / Voter ID",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "BANK_ACCOUNT_PROOF",
				DocumentName:      "Bank Account Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Cancelled cheque / Bank passbook copy",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "DISBURSEMENT",
			},
		}

	case "SURVIVAL_BENEFIT":
		documentTypes = []resp.DocumentTypeItem{
			{
				DocumentCode:      "CLAIM_FORM",
				DocumentName:      "Survival Benefit Claim Form",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Duly filled survival benefit claim form",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "ID_PROOF",
				DocumentName:      "Policyholder ID Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Aadhaar Card / PAN Card / Voter ID",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "BANK_ACCOUNT_PROOF",
				DocumentName:      "Bank Account Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Cancelled cheque / Bank passbook copy",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "DISBURSEMENT",
			},
		}

	case "FREELOOK":
		documentTypes = []resp.DocumentTypeItem{
			{
				DocumentCode:      "CANCELLATION_REQUEST",
				DocumentName:      "Free Look Cancellation Request",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Duly filled free look cancellation request form",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "POLICY_BOND",
				DocumentName:      "Original Policy Bond",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       5,
				Description:       "Original policy bond to be returned",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "APPROVAL",
			},
			{
				DocumentCode:      "ID_PROOF",
				DocumentName:      "Policyholder ID Proof",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Aadhaar Card / PAN Card / Voter ID",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "REGISTRATION",
			},
			{
				DocumentCode:      "BANK_ACCOUNT_PROOF",
				DocumentName:      "Bank Account Proof for Refund",
				Required:          true,
				DocumentFormat:    "PDF",
				MaxFileSize:       2,
				Description:       "Cancelled cheque / Bank passbook copy for refund credit",
				ApplicableFor:     "ALL",
				IsMandatory:       true,
				MandatoryForStage: "DISBURSEMENT",
			},
		}
	}

	return resp.NewDocumentTypesResponse(documentTypes), nil
}

// GetRejectionReasons handles GET /lookup/rejection-reasons
// Reference: BR-CLM-DC-020 (Claim rejection with appeal rights)
func (h *LookupHandler) GetRejectionReasons(sctx *serverRoute.Context, req *GetRejectionReasonsRequest) (*resp.RejectionReasonsResponse, error) {
	log.Info(sctx.Ctx, "GetRejectionReasons: Fetching rejection reasons for claim_type=%s", req.ClaimType)

	// TODO: Integrate with master data service or database lookup table
	// Return rejection reasons based on claim type
	var rejectionReasons []resp.RejectionReasonItem

	// Common rejection reasons for all claim types
	rejectionReasons = append(rejectionReasons, []resp.RejectionReasonItem{
		{
			ReasonCode:         "INCOMPLETE_DOCUMENTS",
			Reason:             "Incomplete Documents",
			Description:        "Required documents are missing or incomplete",
			AppealAllowed:      true,
			AppellateAuthority: "Division Head",
			Category:           "DOCUMENTARY",
		},
		{
			ReasonCode:         "INVALID_DOCUMENTS",
			Reason:             "Invalid Documents",
			Description:        "Documents provided are not valid or expired",
			AppealAllowed:      true,
			AppellateAuthority: "Division Head",
			Category:           "DOCUMENTARY",
		},
		{
			ReasonCode:         "POLICY_LAPSED",
			Reason:             "Policy Lapsed",
			Description:        "Policy has lapsed due to non-payment of premiums",
			AppealAllowed:      true,
			AppellateAuthority: "Zonal Manager",
			Category:           "TECHNICAL",
		},
		{
			ReasonCode:         "CLAIM_TIME_BARRED",
			Reason:             "Claim Time Barred",
			Description:        "Claim filed beyond the limitation period",
			AppealAllowed:      true,
			AppellateAuthority: "Zonal Manager",
			Category:           "LEGAL",
		},
	}...)

	// Claim type-specific rejection reasons
	switch req.ClaimType {
	case "DEATH":
		rejectionReasons = append(rejectionReasons, []resp.RejectionReasonItem{
			{
				ReasonCode:         "CAUSE_OF_DEATH_NOT_COVERED",
				Reason:             "Cause of Death Not Covered",
				Description:        "Death cause is excluded under policy terms",
				AppealAllowed:      true,
				AppellateAuthority: "Claims Committee",
				Category:           "TECHNICAL",
			},
			{
				ReasonCode:         "FRAUDULENT_CLAIM",
				Reason:             "Fraudulent Claim",
				Description:        "Claim found to be fraudulent based on investigation",
				AppealAllowed:      false,
				AppellateAuthority: "",
				Category:           "FRAUD",
			},
			{
				ReasonCode:         "SUPPRESSION_OF_MATERIAL",
				Reason:             "Suppression of Material Fact",
				Description:        "Material facts were suppressed at policy issuance",
				AppealAllowed:      true,
				AppellateAuthority: "Claims Committee",
				Category:           "LEGAL",
			},
			{
				ReasonCode:         "NO_INSURABLE_INTEREST",
				Reason:             "No Insurable Interest",
				Description:        "Claimant has no insurable interest in the policy",
				AppealAllowed:      true,
				AppellateAuthority: "Division Head",
				Category:           "LEGAL",
			},
		}...)

	case "MATURITY", "SURVIVAL_BENEFIT":
		rejectionReasons = append(rejectionReasons, []resp.RejectionReasonItem{
			{
				ReasonCode:         "POLICY_NOT_ASSIGNED",
				Reason:             "Policy Not Assigned",
				Description:        "Policy has not been assigned/issued correctly",
				AppealAllowed:      true,
				AppellateAuthority: "Division Head",
				Category:           "TECHNICAL",
			},
			{
				ReasonCode:         "TITLE_DEFECT",
				Reason:             "Title Defect",
				Description:        "There is a defect in the title of the policy",
				AppealAllowed:      true,
				AppellateAuthority: "Claims Committee",
				Category:           "LEGAL",
			},
			{
				ReasonCode:         "DISPUTE_OVER_BENEFICIARY",
				Reason:             "Dispute Over Beneficiary",
				Description:        "There is a dispute regarding the rightful beneficiary",
				AppealAllowed:      true,
				AppellateAuthority: "Claims Committee",
				Category:           "LEGAL",
			},
		}...)

	case "FREELOOK":
		rejectionReasons = append(rejectionReasons, []resp.RejectionReasonItem{
			{
				ReasonCode:         "FREE_LOOK_PERIOD_EXPIRED",
				Reason:             "Free Look Period Expired",
				Description:        "Cancellation request received after free look period",
				AppealAllowed:      true,
				AppellateAuthority: "Division Head",
				Category:           "TECHNICAL",
			},
			{
				ReasonCode:         "BOND_NOT_RETURNED",
				Reason:             "Policy Bond Not Returned",
				Description:        "Original policy bond not returned within stipulated time",
				AppealAllowed:      true,
				AppellateAuthority: "Division Head",
				Category:           "DOCUMENTARY",
			},
		}...)
	}

	return resp.NewRejectionReasonsResponse(rejectionReasons), nil
}

// GetInvestigationOfficers handles GET /lookup/investigation-officers
// Reference: BR-CLM-DC-002 (Investigation assignment based on jurisdiction and availability)
func (h *LookupHandler) GetInvestigationOfficers(sctx *serverRoute.Context, req *GetInvestigationOfficersRequest) (*resp.InvestigationOfficersResponse, error) {
	log.Info(sctx.Ctx, "GetInvestigationOfficers: Fetching investigation officers for jurisdiction=%s, rank=%v, available_only=%t",
		req.Jurisdiction, req.Rank, req.AvailableOnly)

	// TODO: Integrate with User Service or HRMS for actual officer data
	// For now, return mock data
	officers := []resp.InvestigationOfficerItem{
		{
			OfficerID:      "INV001",
			OfficerName:    "Rajesh Kumar",
			Rank:           "Senior Investigator",
			Jurisdiction:   "Delhi Zone",
			Phone:          "9876543210",
			Email:          "rajesh.kumar@pli.gov.in",
			Available:      true,
			ActiveCases:    5,
			Experience:     15.0,
			Specialization: "FRAUD",
		},
		{
			OfficerID:      "INV002",
			OfficerName:    "Priya Sharma",
			Rank:           "Investigator",
			Jurisdiction:   "Delhi Zone",
			Phone:          "9876543211",
			Email:          "priya.sharma@pli.gov.in",
			Available:      true,
			ActiveCases:    3,
			Experience:     10.0,
			Specialization: "MEDICAL",
		},
		{
			OfficerID:      "INV003",
			OfficerName:    "Amit Singh",
			Rank:           "Senior Investigator",
			Jurisdiction:   "Delhi Zone",
			Phone:          "9876543212",
			Email:          "amit.singh@pli.gov.in",
			Available:      false,
			ActiveCases:    12,
			Experience:     20.0,
			Specialization: "LEGAL",
		},
	}

	// Filter by availability if requested
	if req.AvailableOnly {
		filteredOfficers := make([]resp.InvestigationOfficerItem, 0)
		for _, officer := range officers {
			if officer.Available {
				filteredOfficers = append(filteredOfficers, officer)
			}
		}
		officers = filteredOfficers
	}

	// Filter by rank if specified
	if req.Rank != nil {
		filteredOfficers := make([]resp.InvestigationOfficerItem, 0)
		for _, officer := range officers {
			if officer.Rank == *req.Rank {
				filteredOfficers = append(filteredOfficers, officer)
			}
		}
		officers = filteredOfficers
	}

	if len(officers) == 0 {
		return nil, pgx.ErrNoRows
	}

	return resp.NewInvestigationOfficersResponse(officers), nil
}

// GetApproversList handles GET /lookup/approvers
// Reference: BR-CLM-DC-022 (Approval hierarchy based on claim amount and location)
func (h *LookupHandler) GetApproversList(sctx *serverRoute.Context, req *GetApproversListRequest) (*resp.ApproversListResponse, error) {
	log.Info(sctx.Ctx, "GetApproversList: Fetching approvers for claim_amount=%f, location=%s",
		req.ClaimAmount, req.Location)

	// TODO: Integrate with User Service or Organization Service for actual approver data
	// Determine approval level based on claim amount (BR-CLM-DC-022)
	var approvalLevel string
	var minAmount, maxAmount float64

	switch {
	case req.ClaimAmount <= 50000:
		approvalLevel = "LEVEL_1"
		minAmount = 0
		maxAmount = 50000
	case req.ClaimAmount <= 200000:
		approvalLevel = "LEVEL_2"
		minAmount = 50001
		maxAmount = 200000
	case req.ClaimAmount <= 500000:
		approvalLevel = "LEVEL_3"
		minAmount = 200001
		maxAmount = 500000
	default:
		approvalLevel = "LEVEL_4"
		minAmount = 500001
		maxAmount = 0 // No upper limit
	}

	// Mock approvers data
	approvers := []resp.ApproverItem{
		{
			ApproverID:      "APP001",
			ApproverName:    "Suresh Patil",
			Designation:     "Assistant Division Manager",
			Department:      "Claims",
			Location:        req.Location,
			MinAmount:       minAmount,
			MaxAmount:       maxAmount,
			ApprovalLevel:   approvalLevel,
			Available:       true,
			ActiveApprovals: 8,
			Phone:           "9876543220",
			Email:           "suresh.patil@pli.gov.in",
		},
		{
			ApproverID:      "APP002",
			ApproverName:    "Meena Desai",
			Designation:     "Division Manager",
			Department:      "Claims",
			Location:        req.Location,
			MinAmount:       minAmount,
			MaxAmount:       maxAmount,
			ApprovalLevel:   approvalLevel,
			Available:       true,
			ActiveApprovals: 5,
			Phone:           "9876543221",
			Email:           "meena.desai@pli.gov.in",
		},
	}

	// Filter approvers by availability and matching criteria
	filteredApprovers := make([]resp.ApproverItem, 0)
	for _, approver := range approvers {
		if approver.Available &&
			approver.MinAmount <= req.ClaimAmount &&
			(approver.MaxAmount == 0 || approver.MaxAmount >= req.ClaimAmount) &&
			approver.Location == req.Location {
			filteredApprovers = append(filteredApprovers, approver)
		}
	}

	if len(filteredApprovers) == 0 {
		return nil, pgx.ErrNoRows
	}

	return resp.NewApproversListResponse(filteredApprovers), nil
}

// GetPaymentModes handles GET /lookup/payment-modes
// Reference: BR-CLM-DC-017 (Payment mode priority: NEFT > POSB > Cheque)
func (h *LookupHandler) GetPaymentModes(sctx *serverRoute.Context, req *struct{}) (*resp.PaymentModesResponse, error) {
	log.Info(sctx.Ctx, "GetPaymentModes: Fetching available payment modes")

	// TODO: Integrate with Configuration Service or master data
	paymentModes := []resp.PaymentModeItem{
		{
			PaymentMode:            "NEFT",
			Description:            "National Electronic Funds Transfer",
			Enabled:                true,
			MinAmount:              1.0,
			MaxAmount:              0, // No upper limit
			ProcessingTime:         "1-2 working days",
			RequiresBankValidation: true,
			Charges:                0, // No charges
		},
		{
			PaymentMode:            "POSB",
			Description:            "Pay Order at State Bank of India",
			Enabled:                true,
			MinAmount:              1.0,
			MaxAmount:              0, // No upper limit
			ProcessingTime:         "3-5 working days",
			RequiresBankValidation: true,
			Charges:                25.0, // Processing charges
		},
		{
			PaymentMode:            "CHEQUE",
			Description:            "Cheque Payment",
			Enabled:                true,
			MinAmount:              1.0,
			MaxAmount:              0, // No upper limit
			ProcessingTime:         "7-10 working days",
			RequiresBankValidation: true,
			Charges:                0, // No charges
		},
	}

	return resp.NewPaymentModesResponse(paymentModes), nil
}

// GetApprovalHierarchy handles GET /approvers/financial-limits
// Reference: BR-CLM-DC-022 (4-level approval hierarchy)
func (h *LookupHandler) GetApprovalHierarchy(sctx *serverRoute.Context, req *struct{}) (*resp.ApprovalHierarchyResponse, error) {
	log.Info(sctx.Ctx, "GetApprovalHierarchy: Fetching approval hierarchy and financial limits")

	// TODO: Integrate with Configuration Service or master data
	hierarchy := []resp.ApprovalLevelItem{
		{
			Designation:      "Assistant Division Manager",
			ApprovalLevel:    "LEVEL_1",
			MinAmount:        0,
			MaxAmount:        50000,
			RequiresApproval: true,
			CanOverride:      false,
			EscalationOrder:  1,
		},
		{
			Designation:      "Division Manager",
			ApprovalLevel:    "LEVEL_2",
			MinAmount:        50001,
			MaxAmount:        200000,
			RequiresApproval: true,
			CanOverride:      true,
			EscalationOrder:  2,
		},
		{
			Designation:      "Zonal Manager",
			ApprovalLevel:    "LEVEL_3",
			MinAmount:        200001,
			MaxAmount:        500000,
			RequiresApproval: true,
			CanOverride:      true,
			EscalationOrder:  3,
		},
		{
			Designation:      "Chief Claims Manager",
			ApprovalLevel:    "LEVEL_4",
			MinAmount:        500001,
			MaxAmount:        0, // No upper limit
			RequiresApproval: true,
			CanOverride:      true,
			EscalationOrder:  4,
		},
	}

	return resp.NewApprovalHierarchyResponse(hierarchy), nil
}

// GetBanksList handles GET /lookup/banks
// Reference: Master data for banks
func (h *LookupHandler) GetBanksList(sctx *serverRoute.Context, req *struct{}) (*resp.PaymentModesResponse, error) {
	log.Info(sctx.Ctx, "GetBanksList: Fetching banks list")

	// TODO: Integrate with Bank Master Service or RBI IFSC database
	// This is a placeholder - should return actual bank list
	return resp.NewPaymentModesResponse([]resp.PaymentModeItem{}), nil
}

// GetBranchesList handles GET /lookup/branches
// Reference: Master data for bank branches
func (h *LookupHandler) GetBranchesList(sctx *serverRoute.Context, req *struct{}) (*resp.PaymentModesResponse, error) {
	log.Info(sctx.Ctx, "GetBranchesList: Fetching branches list")

	// TODO: Integrate with Bank Master Service or RBI IFSC database
	// This is a placeholder - should return actual branch list based on IFSC/bank
	return resp.NewPaymentModesResponse([]resp.PaymentModeItem{}), nil
}

// GetStatesList handles GET /lookup/states
// Reference: Master data for Indian states
func (h *LookupHandler) GetStatesList(sctx *serverRoute.Context, req *struct{}) (*resp.PaymentModesResponse, error) {
	log.Info(sctx.Ctx, "GetStatesList: Fetching states list")

	// TODO: Integrate with Location Master Service
	// This is a placeholder - should return actual states list
	return resp.NewPaymentModesResponse([]resp.PaymentModeItem{}), nil
}

// GetDistrictsList handles GET /lookup/districts
// Reference: Master data for Indian districts
func (h *LookupHandler) GetDistrictsList(sctx *serverRoute.Context, req *struct{}) (*resp.PaymentModesResponse, error) {
	log.Info(sctx.Ctx, "GetDistrictsList: Fetching districts list")

	// TODO: Integrate with Location Master Service
	// This is a placeholder - should return actual districts list based on state
	return resp.NewPaymentModesResponse([]resp.PaymentModeItem{}), nil
}
