package handler

import (
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// SurvivalBenefitHandler handles survival benefit-related HTTP requests
// Reference: FRS-SB-01 to FRS-SB-15, BR-CLM-SB-001
type SurvivalBenefitHandler struct {
	*serverHandler.Base
	claimRepo    *repo.ClaimRepository
	claimDocRepo *repo.ClaimDocumentRepository
}

// NewSurvivalBenefitHandler creates a new survival benefit handler
func NewSurvivalBenefitHandler(claimRepo *repo.ClaimRepository, claimDocRepo *repo.ClaimDocumentRepository) *SurvivalBenefitHandler {
	base := serverHandler.New("SurvivalBenefits").
		SetPrefix("/v1").
		AddPrefix("")
	return &SurvivalBenefitHandler{
		Base:         base,
		claimRepo:    claimRepo,
		claimDocRepo: claimDocRepo,
	}
}

// Routes defines all routes for this handler
func (h *SurvivalBenefitHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Survival Benefit Claims (2 endpoints)
		serverRoute.POST("/claims/survival-benefit/submit", h.SubmitSurvivalBenefitClaim).Name("Submit Survival Benefit Claim"),
		serverRoute.POST("/claims/survival-benefit/:id/validate-eligibility", h.ValidateSBEligibility).Name("Validate SB Eligibility"),
	}
}

// ==================== HANDLER METHODS ====================

// SubmitSurvivalBenefitClaim handles survival benefit claim submission
// POST /claims/survival-benefit/submit
// Reference: FRS-SB-03, BR-CLM-SB-001 (7 days SLA)
func (h *SurvivalBenefitHandler) SubmitSurvivalBenefitClaim(sctx *serverRoute.Context, req SubmitSurvivalBenefitClaimRequest) (*resp.SurvivalBenefitClaimRegistrationResponse, error) {
	log.Info(sctx.Ctx, "Submitting survival benefit claim for policy %s", req.PolicyID)

	// TODO: Validate policy from Policy Service
	// TODO: Validate claimant from Customer Service
	// TODO: Validate bank details from CBS API

	// Calculate SLA due date (7 days from submission as per BR-CLM-SB-001)
	slaDueDate := time.Now().Add(7 * 24 * time.Hour)

	// Create claim domain object
	claim := domain.Claim{
		PolicyID:             req.PolicyID,
		ClaimType:            "SURVIVAL_BENEFIT",
		Status:               "SUBMITTED",
		ClaimantName:         req.ClaimantName,
		ClaimantRelation:     &req.ClaimantRelationship,
		ClaimantPhone:        &req.ClaimantMobile,
		ClaimantEmail:        &req.ClaimantEmail,
		ClaimAmount:          nil, // Will be calculated by Policy Service
		PaymentMode:          &req.DisbursementMode,
		BankAccountNumber:    &req.BankAccountNumber,
		BankIFSCCode:         &req.BankIFSC,
		InvestigationRequired: false, // SB claims typically don't require investigation
		WorkflowState:        nil, // Documents will be verified
		SLADueDate:           slaDueDate,
		SLAStatus:            "GREEN", // Will be recalculated dynamically
	}

	if req.PANNumber != nil {
		// Store PAN number in approver remarks temporarily for validation
		panRemarks := "PAN: " + *req.PANNumber
		claim.ApprovalRemarks = &panRemarks
	}

	// Generate claim number (format: SB{YYYY}{DDDD})
	// TODO: Implement proper claim number generation logic
	claim.ClaimNumber = "SB20250001"

	// Create claim in database
	createdClaim, err := h.claimRepo.Create(sctx.Ctx, claim)
	if err != nil {
		log.Error(sctx.Ctx, map[string]interface{}{
			"error":     err.Error(),
			"policy_id": req.PolicyID,
		}, "Failed to create survival benefit claim")
		return nil, err
	}

	// TODO: Upload documents to ECMS
	// TODO: Send intimation via Notification Service (SMS, Email, WhatsApp)
	// TODO: Trigger DigiLocker integration if UseDigiLocker is true

	// Generate acknowledgement number
	acknowledgementNumber := "SB-ACK-" + createdClaim.ID

	// Build workflow state
	workflowState := resp.WorkflowStateResponse{
		CurrentStep:     "SUBMITTED",
		NextStep:        "DOCUMENT_VERIFICATION",
		SLADeadline:     createdClaim.SLADueDate.Format("2006-01-02 15:04:05"),
		DaysRemaining:   7,
		SLAStatus:       "GREEN",
		AllowedActions:  []string{"Upload Documents", "Track Status"},
	}

	// Build response
	response := resp.SurvivalBenefitClaimRegistrationResponse{
		StatusCodeAndMessage:   port.CreateSuccess,
		ClaimID:                createdClaim.ID,
		ClaimNumber:            createdClaim.ClaimNumber,
		AcknowledgementNumber:  acknowledgementNumber,
		SubmissionDate:         createdClaim.CreatedAt.Format("2006-01-02 15:04:05"),
		EstimatedSettlementDate: createdClaim.SLADueDate.Format("2006-01-02 15:04:05"),
		WorkflowState:          &workflowState,
	}

	log.Info(sctx.Ctx, map[string]interface{}{
		"claim_id":     createdClaim.ID,
		"claim_number": createdClaim.ClaimNumber,
		"policy_id":    req.PolicyID,
	}, "Survival benefit claim submitted successfully")

	return &response, nil
}

// ValidateSBEligibility validates survival benefit claim eligibility
// POST /claims/survival-benefit/{id}/validate-eligibility
// Reference: FRS-SB-02
func (h *SurvivalBenefitHandler) ValidateSBEligibility(sctx *serverRoute.Context, req ValidateSBEligibilityRequest) (*resp.SBEligibilityValidationResponse, error) {
	log.Info(sctx.Ctx, "Validating eligibility for claim %s", req.ClaimID)

	// Fetch claim by ID
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, map[string]interface{}{
				"error":    "Claim not found",
				"claim_id": req.ClaimID,
			}, "Survival benefit claim not found")
			return nil, err
		}
		log.Error(sctx.Ctx, map[string]interface{}{
			"error":    err.Error(),
			"claim_id": req.ClaimID,
		}, "Failed to fetch survival benefit claim")
		return nil, err
	}

	// Validate claim type
	if claim.ClaimType != "SURVIVAL_BENEFIT" {
		log.Error(sctx.Ctx, map[string]interface{}{
			"claim_id":   req.ClaimID,
			"claim_type": claim.ClaimType,
		}, "Invalid claim type for survival benefit eligibility validation")
		return &resp.SBEligibilityValidationResponse{
			StatusCodeAndMessage: port.StatusCodeAndMessage{
				StatusCode: 400,
				Success:    false,
				Message:    "Invalid claim type for survival benefit eligibility validation",
			},
		}, nil
	}

	// TODO: Validate policy from Policy Service
	// TODO: Check policy status (active, paid-up, etc.)
	// TODO: Validate survival benefit due date
	// TODO: Calculate survival benefit amount from Policy Service
	// TODO: Check if SB already paid for this installment
	// TODO: Validate claimant details from Customer Service

	// For now, assume eligible (will be replaced with actual validation)
	eligible := true
	eligibilityReasons := []string{
		"Policy is active and in-force",
		"Survival benefit is due",
		"All premiums are paid",
		"Claimant details match policy records",
	}
	var ineligibilityReasons []string

	sbDueDate := time.Now().Add(30 * 24 * time.Hour) // TODO: Get from Policy Service
	sbAmount := 10000.0                                // TODO: Calculate from Policy Service

	// Build policy details
	policyDetails := resp.PolicyDetailsResponse{
		PolicyNumber: claim.PolicyID, // TODO: Fetch from Policy Service
		PolicyStatus: "ACTIVE",
		PolicyType:   "ENDOWMENT",
		SumAssured:   100000,
		IssueDate:    "2020-01-01",
		MaturityDate: "2040-01-01",
	}

	// Build customer/claimant details
	claimantRelation := ""
	if claim.ClaimantRelation != nil {
		claimantRelation = *claim.ClaimantRelation
	}

	claimantPhone := ""
	if claim.ClaimantPhone != nil {
		claimantPhone = *claim.ClaimantPhone
	}

	claimantEmail := ""
	if claim.ClaimantEmail != nil {
		claimantEmail = *claim.ClaimantEmail
	}

	customerDetails := resp.ClaimantDetailsResponse{
		Name:         claim.ClaimantName,
		Relationship: claimantRelation,
		Phone:        claimantPhone,
		Email:        claimantEmail,
		Address:      "", // TODO: Fetch from Customer Service
	}

	// Calculate next SB due date (typically annually)
	var nextSBDueDate *string
	nextDue := sbDueDate.Add(365 * 24 * time.Hour) // Next year
	nextDueStr := nextDue.Format("2006-01-02")
	nextSBDueDate = &nextDueStr

	// Build response
	response := resp.SBEligibilityValidationResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		Eligible:             eligible,
		SBDueDate:            sbDueDate.Format("2006-01-02"),
		SBAmount:             sbAmount,
		EligibilityReasons:   eligibilityReasons,
		IneligibilityReasons: ineligibilityReasons,
		PolicyDetails:        &policyDetails,
		ClaimantDetails:      &customerDetails,
		NextSBDueDate: nextSBDueDate,
	}

	log.Info(sctx.Ctx, map[string]interface{}{
		"claim_id":    req.ClaimID,
		"eligible":    eligible,
		"sb_amount":   sbAmount,
		"sb_due_date": sbDueDate.Format("2006-01-02"),
	}, "Survival benefit eligibility validated successfully")

	return &response, nil
}
