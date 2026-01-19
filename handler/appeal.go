package handler

import (
	"fmt"
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

// AppealHandler handles appeal-related HTTP requests
// Reference: seed/swagger/claims_api_swagger_complete.yaml (Death Claims - Appeal section)
type AppealHandler struct {
	*serverHandler.Base
	appealRepo *repo.AppealRepository
	claimRepo  *repo.ClaimRepository
}

// NewAppealHandler creates a new appeal handler
func NewAppealHandler(appealRepo *repo.AppealRepository, claimRepo *repo.ClaimRepository) *AppealHandler {
	base := serverHandler.New("Appeals").
		SetPrefix("/v1").
		AddPrefix("")
	return &AppealHandler{
		Base:       base,
		appealRepo: appealRepo,
		claimRepo:  claimRepo,
	}
}

// Routes defines all routes for this handler
// Reference: seed/swagger/claims_api_swagger_complete.yaml:887-1008
func (h *AppealHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Appeal Workflow (3 endpoints)
		serverRoute.GET("/claims/death/:claim_id/appeal-eligibility", h.CheckAppealEligibility).Name("Check Appeal Eligibility"),
		serverRoute.GET("/claims/death/:claim_id/appellate-authority", h.GetAppellateAuthority).Name("Get Appellate Authority"),
		serverRoute.POST("/claims/death/:claim_id/appeal", h.SubmitAppeal).Name("Submit Appeal"),
		serverRoute.POST("/claims/death/:claim_id/appeal/:appeal_id/decision", h.RecordAppealDecision).Name("Record Appeal Decision"),
	}
}

// CheckAppealEligibility checks if a claim is eligible for appeal
// GET /claims/death/{claim_id}/appeal-eligibility
// Reference: FR-CLM-DC-023, BR-CLM-DC-005 (90-day appeal window)
func (h *AppealHandler) CheckAppealEligibility(sctx *serverRoute.Context, req *CheckAppealEligibilityUri) (*resp.AppealEligibilityResponse, error) {
	log.Info(sctx.Ctx, "Checking appeal eligibility for claim_id: %v", req.ClaimID)

	// Fetch claim details
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %v", req.ClaimID)
			return nil, &port.ErrorResponse{
				StatusCode: 404,
				Success:    false,
				Message:    "Claim not found",
			}
		}
		log.Error(sctx.Ctx, "Failed to fetch claim: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to fetch claim details",
		}
	}

	// Check if claim is in rejected status
	if claim.Status != "REJECTED" {
		return &resp.AppealEligibilityResponse{
			StatusCodeAndMessage: port.StatusCodeAndMessage{
				StatusCode: 400,
				Success:    false,
				Message:    "Appeal can only be filed for rejected claims",
			},
			Data: resp.AppealEligibilityData{
				ClaimID:              claim.ID,
				IsEligible:           false,
				RejectionReason:      claim.RejectionReason,
				AppealTypesAvailable: []string{},
			},
		}, nil
	}

	// Check appeal eligibility using repository
	// Reference: BR-CLM-DC-005 (90-day appeal window)
	eligible, err := h.appealRepo.CheckAppealEligibility(sctx.Ctx, req.ClaimID)
	if err != nil && err != pgx.ErrNoRows {
		log.Error(sctx.Ctx, "Failed to check appeal eligibility: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to check appeal eligibility",
		}
	}

	// Calculate days remaining
	daysRemaining := 0
	appealWindowEnds := ""
	condonationRequired := false
	var condonationReason *string

	if claim.RejectedAt != nil {
		rejectionDate := claim.RejectedAt
		deadline := rejectionDate.AddDate(0, 0, 90) // 90 days from rejection (BR-CLM-DC-005)
		appealWindowEnds = deadline.Format("2006-01-02 15:04:05")
		daysRemaining = int(time.Until(deadline).Hours() / 24)

		// Check if appeal window has expired
		if time.Now().After(deadline) {
			condonationRequired = true
			reason := "Appeal window has expired. Condonation of delay is required."
			condonationReason = &reason
		}
	}

	// Determine available appeal types
	appealTypes := []string{"RECONSIDERATION", "APPELLATE_AUTHORITY"}
	if claim.ApprovalLevel != nil && *claim.ApprovalLevel >= 3 {
		// Ombudsman option available for higher-level rejections
		appealTypes = append(appealTypes, "OMBUDSMAN")
	}

	return &resp.AppealEligibilityResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Appeal eligibility checked successfully",
		},
		Data: resp.AppealEligibilityData{
			ClaimID:              claim.ID,
			IsEligible:           eligible.IsEligible,
			AppealWindowEnds:     appealWindowEnds,
			DaysRemaining:        daysRemaining,
			RejectionReason:      claim.RejectionReason,
			CondonationRequired:  condonationRequired,
			CondonationReason:    condonationReason,
			AppealTypesAvailable: appealTypes,
		},
	}, nil
}

// GetAppellateAuthority retrieves the appellate authority details for a claim
// GET /claims/death/{claim_id}/appellate-authority
// Reference: BR-CLM-DC-005 (Escalation to next higher authority in approval hierarchy)
func (h *AppealHandler) GetAppellateAuthority(sctx *serverRoute.Context, req *GetAppellateAuthorityUri) (*resp.AppellateAuthorityResponse, error) {
	log.Info(sctx.Ctx, "Getting appellate authority for claim_id: %v", req.ClaimID)

	// Fetch claim details
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %v", req.ClaimID)
			return nil, &port.ErrorResponse{
				StatusCode: 404,
				Success:    false,
				Message:    "Claim not found",
			}
		}
		log.Error(sctx.Ctx, "Failed to fetch claim: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to fetch claim details",
		}
	}

	// Determine appellate authority based on approval level
	// Reference: BR-CLM-DC-022 (Approval hierarchy - 4 levels)
	var authority resp.AppellateAuthorityData

	currentLevel := 0
	if claim.ApprovalLevel != nil {
		currentLevel = int(*claim.ApprovalLevel)
	}

	// Escalate to next higher authority
	switch currentLevel {
	case 1:
		// Approved by Division Head -> Escalate to Zonal Manager
		authority = resp.AppellateAuthorityData{
			AuthorityID:    "LEVEL_2",
			AuthorityName:  "Zonal Manager",
			AuthorityLevel: "LEVEL_2",
			AuthorityType:  "ZONAL_MANAGER",
			Designation:    "Zonal Manager",
			Department:     "Operations",
			Location:       getStringValue(claim.ZonalOffice),
			MaxClaimAmount: 1000000.0, // ₹10 lakh
		}
	case 2:
		// Approved by Zonal Manager -> Escalate to Regional Director
		authority = resp.AppellateAuthorityData{
			AuthorityID:    "LEVEL_3",
			AuthorityName:  "Regional Director",
			AuthorityLevel: "LEVEL_3",
			AuthorityType:  "REGIONAL_DIRECTOR",
			Designation:    "Regional Director",
			Department:     "Claims",
			Location:       getStringValue(claim.RegionalOffice),
			MaxClaimAmount: 5000000.0, // ₹50 lakh
		}
	case 3:
		// Approved by Regional Director -> Escalate to Chief Claims Manager
		authority = resp.AppellateAuthorityData{
			AuthorityID:    "LEVEL_4",
			AuthorityName:  "Chief Claims Manager",
			AuthorityLevel: "LEVEL_4",
			AuthorityType:  "CHIEF_CLAIMS_MANAGER",
			Designation:    "Chief Claims Manager",
			Department:     "Corporate Claims",
			Location:       "Head Office",
			MaxClaimAmount: 10000000.0, // ₹1 crore
		}
	default:
		// No approval level or already at highest level -> Default to Division Head
		authority = resp.AppellateAuthorityData{
			AuthorityID:    "LEVEL_1",
			AuthorityName:  "Division Head",
			AuthorityLevel: "LEVEL_1",
			AuthorityType:  "DIVISION_HEAD",
			Designation:    "Division Head",
			Department:     "Claims Division",
			Location:       getStringValue(claim.DivisionOffice),
			MaxClaimAmount: 500000.0, // ₹5 lakh
		}
	}

	return &resp.AppellateAuthorityResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Appellate authority details retrieved successfully",
		},
		Data: authority,
	}, nil
}

// SubmitAppeal submits an appeal for a rejected claim
// POST /claims/death/{claim_id}/appeal
// Reference: BR-CLM-DC-005 (90-day window), BR-CLM-DC-007 (45-day SLA)
func (h *AppealHandler) SubmitAppeal(sctx *serverRoute.Context, req *SubmitAppealRequest) (*resp.AppealSubmissionResponse, error) {
	log.Info(sctx.Ctx, "Submitting appeal for claim_id: %v", req.ClaimID)

	// Validate claim is rejected
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %v", req.ClaimID)
			return nil, &port.ErrorResponse{
				StatusCode: 404,
				Success:    false,
				Message:    "Claim not found",
			}
		}
		log.Error(sctx.Ctx, "Failed to fetch claim: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to fetch claim details",
		}
	}

	if claim.Status != "REJECTED" {
		return nil, &port.ErrorResponse{
			StatusCode: 400,
			Success:    false,
			Message:    "Appeal can only be filed for rejected claims",
		}
	}

	// Validate appeal deadline
	// Reference: BR-CLM-DC-005 (90-day appeal window)
	if claim.RejectedAt != nil {
		deadline := claim.RejectedAt.AddDate(0, 0, 90)
		if time.Now().After(deadline) && !req.CondonationRequest {
			return nil, &port.ErrorResponse{
				StatusCode: 400,
				Success:    false,
				Message:    "Appeal window has expired. Please request condonation of delay.",
			}
		}
	}

	// Check for duplicate appeals
	eligible, err := h.appealRepo.CheckAppealEligibility(sctx.Ctx, req.ClaimID)
	if err != nil && err != pgx.ErrNoRows {
		log.Error(sctx.Ctx, "Failed to check duplicate appeals: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to validate appeal eligibility",
		}
	}

	if !eligible.IsEligible {
		return nil, &port.ErrorResponse{
			StatusCode: 400,
			Success:    false,
			Message:    "An appeal is already pending for this claim",
		}
	}

	// Generate appeal number
	appealNumber, err := h.appealRepo.GenerateAppealNumber(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate appeal number: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to generate appeal number",
		}
	}

	// Calculate appeal deadline (45 days from submission)
	// Reference: BR-CLM-DC-007 (45-day SLA for decision)
	submissionDate := time.Now()
	appealDeadline := submissionDate.AddDate(0, 0, 45)

	// Determine appellate authority
	appellateAuthorityID := ""
	if claim.ApprovalLevel != nil {
		currentLevel := int(*claim.ApprovalLevel)
		appellateAuthorityID = fmt.Sprintf("LEVEL_%d", currentLevel+1)
	} else {
		appellateAuthorityID = "LEVEL_1"
	}

	// TODO: Get user context from JWT token
	// appellantName := sctx.User.Name
	// appellantContact := sctx.User.Contact
	appellantName := "Claimant" // Placeholder
	appellantContact := ""      // Placeholder

	// Create appeal domain object
	appeal := domain.Appeal{
		AppealNumber:         appealNumber,
		ClaimID:              claim.ID,
		AppellantName:        appellantName,
		AppellantContact:     appellantContact,
		GroundsOfAppeal:      req.AppealGrounds,
		SupportingDocuments:  req.AdditionalDocuments,
		SubmissionDate:       &submissionDate,
		AppealDeadline:       &appealDeadline,
		AppellateAuthorityID: &appellateAuthorityID,
		Status:               "SUBMITTED",
	}

	// Create appeal in database
	createdAppeal, err := h.appealRepo.Create(sctx.Ctx, appeal)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to create appeal: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to submit appeal",
		}
	}

	expectedDecisionBy := appealDeadline.Format("2006-01-02 15:04:05")
	appealSLAStatus := resp.CalculateAppealSLAStatus(submissionDate.Format(time.RFC3339), appealDeadline.Format(time.RFC3339))

	// TODO: Fetch authority name from appellate_authority_id
	// For now, use placeholder
	appellateAuthorityName := "Appellate Authority"

	return &resp.AppealSubmissionResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 201,
			Success:    true,
			Message:    "Appeal submitted successfully",
		},
		Data: resp.AppealSubmittedData{
			AppealID:           createdAppeal.ID,
			AppealNumber:       appealNumber,
			ClaimID:            claim.ID,
			SubmissionDate:     submissionDate.Format("2006-01-02 15:04:05"),
			AppealDeadline:     appealDeadline.Format("2006-01-02 15:04:05"),
			AppellateAuthority: appellateAuthorityName,
			CurrentStatus:      createdAppeal.Status,
			ExpectedDecisionBy: expectedDecisionBy,
			AppealSLAStatus:    appealSLAStatus,
			TrackingURL:        nil, // TODO: Generate tracking URL
			NextSteps: []string{
				"Appeal will be reviewed by appellate authority",
				"You will be notified of the decision within 45 days",
				"Track appeal status using appeal number: " + appealNumber,
			},
			DocumentsRequired: []string{}, // TODO: Check if additional documents needed
		},
	}, nil
}

// RecordAppealDecision records the appellate authority's decision on an appeal
// POST /claims/death/{claim_id}/appeal/{appeal_id}/decision
// Reference: BR-CLM-DC-006 (45-day SLA)
func (h *AppealHandler) RecordAppealDecision(sctx *serverRoute.Context, req *RecordAppealDecisionRequest) (*resp.AppealDecisionResponse, error) {
	log.Info(sctx.Ctx, "Recording appeal decision for appeal_id: %v, claim_id: %v", req.AppealID, req.ClaimID)

	// Fetch appeal details
	appeal, err := h.appealRepo.FindByID(sctx.Ctx, req.AppealID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Appeal not found: %v", req.AppealID)
			return nil, &port.ErrorResponse{
				StatusCode: 404,
				Success:    false,
				Message:    "Appeal not found",
			}
		}
		log.Error(sctx.Ctx, "Failed to fetch appeal: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to fetch appeal details",
		}
	}

	// Check if appeal is already decided
	if appeal.Status == "DECIDED" {
		return nil, &port.ErrorResponse{
			StatusCode: 400,
			Success:    false,
			Message:    "Appeal has already been decided",
		}
	}

	// Record appeal decision
	decisionDate := time.Now()

	// TODO: Get user context from JWT token
	// decisionBy := sctx.User.Name
	decisionBy := "Appellate Authority" // Placeholder

	// Update appeal with decision
	updatedAppeal, err := h.appealRepo.RecordDecision(sctx.Ctx, req.AppealID, req.Decision, req.ReasonedOrder, req.ModificationDetails, decisionBy, &decisionDate)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to record appeal decision: %v", err)
		return nil, &port.ErrorResponse{
			StatusCode: 500,
			Success:    false,
			Message:    "Failed to record appeal decision",
		}
	}

	// Update claim status based on appeal decision
	// Reference: BR-CLM-DC-020 (Appeal outcomes)
	claimStatusUpdate := ""
	var revisedAmount *float64
	nextSteps := []string{}

	switch req.Decision {
	case "APPEAL_ACCEPTED":
		// Appeal accepted - revert claim and reprocess
		claimStatusUpdate = "UNDER_REVIEW"
		nextSteps = append(nextSteps,
			"Claim will be re-evaluated based on appeal decision",
			"You will be notified of the revised claim amount",
			"Disbursement will be processed after re-evaluation",
		)
	case "PARTIAL_ACCEPTANCE":
		// Partial acceptance - update claim amount
		claimStatusUpdate = "APPROVED"
		revisedAmount = updatedAppeal.RevisedClaimAmount
		nextSteps = append(nextSteps,
			"Claim approved with revised amount",
			"Disbursement will be processed within 7 days",
			"You will receive payment via NEFT to registered bank account",
		)
	case "APPEAL_REJECTED":
		// Appeal rejected - claim remains rejected
		claimStatusUpdate = "REJECTED"
		nextSteps = append(nextSteps,
			"Appeal has been rejected",
			"You may approach Ombudsman within 90 days",
			"Contact customer service for further assistance",
		)
	}

	// Update claim status
	_, err = h.claimRepo.UpdateStatus(sctx.Ctx, req.ClaimID, claimStatusUpdate)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to update claim status: %v", err)
		// Log error but don't fail the request
	}

	return &resp.AppealDecisionResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Appeal decision recorded successfully",
		},
		Data: resp.AppealDecisionData{
			AppealID:            appeal.ID,
			AppealNumber:        appeal.AppealNumber,
			Decision:            req.Decision,
			ReasonedOrder:       req.ReasonedOrder,
			DecisionDate:        decisionDate.Format("2006-01-02 15:04:05"),
			DecisionBy:          decisionBy,
			RevisedClaimAmount:  revisedAmount,
			ModificationDetails: req.ModificationDetails,
			ClaimStatusUpdate:   claimStatusUpdate,
			NextSteps:           nextSteps,
		},
	}, nil
}

// Helper function to safely get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
