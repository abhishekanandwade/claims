package handler

import (
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// ClaimHandler handles claim-related HTTP requests
type ClaimHandler struct {
	*serverHandler.Base
	svc *repo.ClaimRepository
}

// NewClaimHandler creates a new claim handler
func NewClaimHandler(svc *repo.ClaimRepository) *ClaimHandler {
	base := serverHandler.New("Claims").
		SetPrefix("/v1").
		AddPrefix("")
	return &ClaimHandler{
		Base: base,
		svc:  svc,
	}
}

// Routes defines all routes for this handler
func (h *ClaimHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Death Claims - Core Operations (15 endpoints)
		serverRoute.POST("/claims/death/register", h.RegisterDeathClaim).Name("Register Death Claim"),
		serverRoute.POST("/claims/death/calculate-amount", h.CalculateDeathClaimAmount).Name("Calculate Death Claim Amount"),
		serverRoute.GET("/claims/death/:claim_id/document-checklist", h.GetDocumentChecklist).Name("Get Document Checklist"),
		serverRoute.GET("/claims/death/document-checklist-dynamic", h.GetDynamicDocumentChecklist).Name("Get Dynamic Document Checklist"),
		serverRoute.POST("/claims/death/:claim_id/documents", h.UploadClaimDocuments).Name("Upload Claim Documents"),
		serverRoute.GET("/claims/death/:claim_id/document-completeness", h.CheckDocumentCompleteness).Name("Check Document Completeness"),
		serverRoute.POST("/claims/death/:claim_id/calculate-benefit", h.CalculateBenefit).Name("Calculate Benefit"),
		serverRoute.GET("/claims/death/:claim_id/eligible-approvers", h.GetEligibleApprovers).Name("Get Eligible Approvers"),
		serverRoute.GET("/claims/death/:claim_id/approval-details", h.GetApprovalDetails).Name("Get Approval Details"),
		serverRoute.POST("/claims/death/:claim_id/approve", h.ApproveClaim).Name("Approve Claim"),
		serverRoute.POST("/claims/death/:claim_id/reject", h.RejectClaim).Name("Reject Claim"),
		serverRoute.POST("/claims/death/:claim_id/disburse", h.DisburseClaim).Name("Disburse Claim"),
		serverRoute.POST("/claims/death/:claim_id/close", h.CloseClaim).Name("Close Claim"),
		serverRoute.POST("/claims/death/:claim_id/cancel", h.CancelClaim).Name("Cancel Claim"),
		serverRoute.GET("/claims/death/approval-queue", h.GetApprovalQueue).Name("Get Approval Queue"),
	}
}

// RegisterDeathClaim registers a new death claim
// POST /claims/death/register
// Reference: FR-CLM-DC-001, BR-CLM-DC-001, WF-CLM-DC-001
func (h *ClaimHandler) RegisterDeathClaim(sctx *serverRoute.Context, req RegisterDeathClaimRequest) (*resp.DeathClaimRegisteredResponse, error) {
	// Parse death date
	deathDate, err := time.Parse("2006-01-02", req.DeathDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid death date format: %v", err)
		return nil, err
	}

	// Convert request to domain model
	data := req.ToDomain()
	data.DeathDate = &deathDate

	// Call repository to create claim
	result, err := h.svc.Create(sctx.Ctx, data)
	if err != nil {
		log.Error(sctx.Ctx, "Error registering death claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Death claim registered with ID: %s, Claim Number: %s", result.ID, result.ClaimNumber)

	// Build response
	r := &resp.DeathClaimRegisteredResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: resp.DeathClaimRegistrationData{
			ClaimID:               result.ID,
			ClaimNumber:           result.ClaimNumber,
			Status:                result.Status,
			AcknowledgmentNumber:  "ACK" + result.ClaimNumber, // TODO: Generate proper acknowledgment number
			InvestigationRequired: result.InvestigationRequired,
		},
	}
	return r, nil
}

// CalculateDeathClaimAmount pre-calculates death claim amount
// POST /claims/death/calculate-amount
// Reference: CALC-001
func (h *ClaimHandler) CalculateDeathClaimAmount(sctx *serverRoute.Context, req CalculateDeathClaimAmountRequest) (*resp.ClaimAmountCalculationResponse, error) {
	// TODO: Implement benefit calculation logic
	// This will integrate with policy service for actual calculation
	// For now, returning a placeholder response

	sumAssured := 500000.00
	reversionaryBonus := 100000.00
	terminalBonus := 50000.00
	outstandingLoan := 0.0
	unpaidPremiums := 0.0

	r := &resp.ClaimAmountCalculationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.ClaimCalculationData{
			SumAssured:      sumAssured,
			AccruedBonuses:  reversionaryBonus + terminalBonus,
			Deductions:      outstandingLoan + unpaidPremiums,
			NetClaimAmount:  sumAssured + reversionaryBonus + terminalBonus - outstandingLoan - unpaidPremiums,
			CalculationBreakdown: resp.ClaimCalculationBreakdown{
				SumAssured:        sumAssured,
				ReversionaryBonus: reversionaryBonus,
				TerminalBonus:     terminalBonus,
				OutstandingLoan:   outstandingLoan,
				UnpaidPremiums:    unpaidPremiums,
			},
		},
	}
	return r, nil
}

// GetDocumentChecklist retrieves document checklist for a claim
// GET /claims/death/{claim_id}/document-checklist
// Reference: FR-CLM-DC-002, VR-CLM-DC-001 to VR-CLM-DC-007
func (h *ClaimHandler) GetDocumentChecklist(sctx *serverRoute.Context, req ClaimIDUri) (*resp.DocumentChecklistResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Fetch document checklist from document_checklist table
	// For now, returning a placeholder response
	r := &resp.DocumentChecklistResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.DocumentChecklistData{
			ClaimID:   claim.ID,
			Documents: []resp.DocumentChecklistItem{},
		},
	}
	return r, nil
}

// GetDynamicDocumentChecklist retrieves context-aware document checklist
// GET /claims/death/document-checklist-dynamic
// Reference: DFC-001
func (h *ClaimHandler) GetDynamicDocumentChecklist(sctx *serverRoute.Context, req GetDynamicDocumentChecklistUri) (*resp.DynamicDocumentChecklistResponse, error) {
	// TODO: Implement dynamic checklist logic based on death type, nomination status, policy type
	// For now, returning a placeholder response
	r := &resp.DynamicDocumentChecklistResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.DynamicDocumentChecklistData{
			MandatoryBase:    10, // Base mandatory documents
			ConditionalAdded: 0,  // Additional documents based on context
			TotalMandatory:   10,
			Documents:        []resp.DocumentChecklistItem{},
		},
	}
	return r, nil
}

// UploadClaimDocuments uploads documents for a claim
// POST /claims/death/{claim_id}/documents
func (h *ClaimHandler) UploadClaimDocuments(sctx *serverRoute.Context, req UploadClaimDocumentsRequest) (*resp.DocumentChecklistResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Implement document upload to ECMS
	// TODO: Create claim_document record
	log.Info(sctx.Ctx, "Document uploaded for claim ID: %s, Document Type: %s", req.ClaimID, req.DocumentType)

	// Return document checklist showing the uploaded document
	r := &resp.DocumentChecklistResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: resp.DocumentChecklistData{
			ClaimID:   claim.ID,
			Documents: []resp.DocumentChecklistItem{
				{
					DocumentType: req.DocumentType,
					Mandatory:    false,
					Description:  req.DocumentName,
					Uploaded:     true,
					DocumentID:   "DOC-" + claim.ID, // TODO: Use actual document ID from ECMS
					UploadedAt:   time.Now().Format("2006-01-02 15:04:05"),
				},
			},
		},
	}
	return r, nil
}

// CheckDocumentCompleteness checks if all mandatory documents are submitted
// GET /claims/death/{claim_id}/document-completeness
func (h *ClaimHandler) CheckDocumentCompleteness(sctx *serverRoute.Context, req ClaimIDUri) (*resp.DocumentCompletenessResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Implement document completeness check
	r := &resp.DocumentCompletenessResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.DocumentCompletenessData{
			ClaimID:            claim.ID,
			DocumentsComplete:  false,
			MandatoryCount:     10,
			UploadedCount:      0,
			VerifiedCount:      0,
			MissingDocuments:   []string{},
			PendingVerification: []string{},
		},
	}
	return r, nil
}

// CalculateBenefit calculates the death claim benefit amount
// POST /claims/death/{claim_id}/calculate-benefit
// Reference: FR-CLM-DC-004, BR-CLM-DC-008
func (h *ClaimHandler) CalculateBenefit(sctx *serverRoute.Context, req ClaimIDUri) (*resp.BenefitCalculationResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Implement benefit calculation
	sumAssured := 500000.00
	reversionaryBonus := 100000.00
	terminalBonus := 50000.00
	outstandingLoan := 0.0
	unpaidPremiums := 0.0

	r := &resp.BenefitCalculationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.ClaimCalculationData{
			ClaimID:        claim.ID,
			SumAssured:      sumAssured,
			AccruedBonuses:  reversionaryBonus + terminalBonus,
			Deductions:      outstandingLoan + unpaidPremiums,
			NetClaimAmount:  sumAssured + reversionaryBonus + terminalBonus - outstandingLoan - unpaidPremiums,
			CalculationBreakdown: resp.ClaimCalculationBreakdown{
				SumAssured:        sumAssured,
				ReversionaryBonus: reversionaryBonus,
				TerminalBonus:     terminalBonus,
				OutstandingLoan:   outstandingLoan,
				UnpaidPremiums:    unpaidPremiums,
			},
		},
	}
	return r, nil
}

// GetEligibleApprovers retrieves eligible approvers for a claim
// GET /claims/death/{claim_id}/eligible-approvers
func (h *ClaimHandler) GetEligibleApprovers(sctx *serverRoute.Context, req ClaimIDUri) (*resp.EligibleApproversResponse, error) {
	// Verify claim exists
	_, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Fetch eligible approvers based on claim amount and user hierarchy
	r := &resp.EligibleApproversResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.EligibleApproversData{
			ClaimAmount:       500000.00, // TODO: Get from claim calculation
			RequiredAuthority: "ASSISTANT_DIRECTOR", // TODO: Calculate based on claim amount
			EligibleApprovers: []resp.ApproverInfo{},
		},
	}
	return r, nil
}

// GetApprovalDetails retrieves approval workflow details for a claim
// GET /claims/death/{claim_id}/approval-details
func (h *ClaimHandler) GetApprovalDetails(sctx *serverRoute.Context, req ClaimIDUri) (*resp.ApprovalDetailsResponse, error) {
	// Verify claim exists
	_, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Fetch approval workflow details
	r := &resp.ApprovalDetailsResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: resp.ApprovalDetailsData{
			ClaimID:         req.ClaimID,
			PolicyDetails:   resp.PolicyDetailsResponse{}, // TODO: Fetch from policy service
			ClaimantDetails: resp.ClaimantDetailsResponse{}, // TODO: Populate
			Calculation:     resp.ClaimCalculationData{}, // TODO: Calculate
			Documents:       []resp.DocumentChecklistItem{}, // TODO: Fetch documents
		},
	}
	return r, nil
}

// ApproveClaim approves a death claim
// POST /claims/death/{claim_id}/approve
// Reference: BR-CLM-DC-005
func (h *ClaimHandler) ApproveClaim(sctx *serverRoute.Context, req ApproveClaimRequest) (*resp.ClaimApprovalResponse, error) {
	// Verify claim exists
	_, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// Update claim status to APPROVED
	// TODO: Get approver ID from context/session
	// TODO: Get approved amount from calculation
	updates := map[string]interface{}{
		"status":           "APPROVED",
		"approval_remarks": req.ApprovalRemarks,
	}

	updatedClaim, err := h.svc.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error approving claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Claim approved with ID: %s", req.ClaimID)

	r := &resp.ClaimApprovalResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ClaimApprovalData{
			ClaimID:          updatedClaim.ID,
			ApprovalDecision: "APPROVED",
			Approver:         "SYSTEM", // TODO: Get from context
			ApprovedAt:       time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	return r, nil
}

// RejectClaim rejects a death claim
// POST /claims/death/{claim_id}/reject
// Reference: BR-CLM-DC-020
func (h *ClaimHandler) RejectClaim(sctx *serverRoute.Context, req RejectClaimRequest) (*resp.ClaimRejectionResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// Update claim status to REJECTED
	updates := map[string]interface{}{
		"status":           "REJECTED",
		"rejection_reason": &req.DetailedJustification,
	}

	_, err = h.svc.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error rejecting claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Claim rejected with ID: %s, Reason: %s", req.ClaimID, req.RejectionReason)

	r := &resp.ClaimRejectionResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ClaimRejectionData{
			ClaimID:                 claim.ID,
			RejectionReason:         req.RejectionReason,
			DetailedJustification:   req.DetailedJustification,
			AppealRightsCommunicated: req.AppealRightsCommunicated,
			AppealDeadline:          "", // TODO: Calculate appeal deadline (30 days)
		},
	}
	return r, nil
}

// DisburseClaim disburses payment for an approved claim
// POST /claims/death/{claim_id}/disburse
// Reference: BR-CLM-DC-010
func (h *ClaimHandler) DisburseClaim(sctx *serverRoute.Context, req DisburseClaimRequest) (*resp.ClaimDisbursementResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// TODO: Implement disbursement via PFMS integration
	// Update claim with disbursement details
	updates := map[string]interface{}{
		"status":            "DISBURSED",
		"payment_mode":      &req.PaymentMode,
		"payment_reference": req.PaymentDetails, // Using payment_details as reference
		"disbursement_date": time.Now(),        // Set disbursement date to now
	}

	_, err = h.svc.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error disbursing claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Claim disbursed with ID: %s, Mode: %s", req.ClaimID, req.PaymentMode)

	r := &resp.ClaimDisbursementResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ClaimDisbursementData{
			PaymentID:          "PAY-" + claim.ID, // TODO: Generate proper payment ID
			PaymentReference:   getStringValue(req.PaymentDetails),
			Status:             "DISBURSED",
			InitiatedAt:        time.Now().Format("2006-01-02 15:04:05"),
			PaymentMode:        req.PaymentMode,
			Amount:             500000.00, // TODO: Get from approved amount
			BankAccountNumber:  "", // TODO: Get from claim
			IfscCode:           "", // TODO: Get from claim
			EstimatedCreditDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"), // T+2 days
		},
	}
	return r, nil
}

// CloseClaim closes a claim
// POST /claims/death/{claim_id}/close
func (h *ClaimHandler) CloseClaim(sctx *serverRoute.Context, req CloseClaimRequest) (*resp.ClaimCloseResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// Update claim status to CLOSED
	updates := map[string]interface{}{
		"status":         "CLOSED",
		"closure_reason": &req.ClosureReason,
		"closure_date":   time.Now(), // Set closure date to now
	}

	_, err = h.svc.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error closing claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Claim closed with ID: %s, Reason: %s", req.ClaimID, req.ClosureReason)

	r := &resp.ClaimCloseResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ClaimCloseData{
			ClaimID:  claim.ID,
			Status:   "CLOSED",
			ClosedAt: time.Now().Format("2006-01-02 15:04:05"),
			ClosedBy: "SYSTEM", // TODO: Get from context
			Remarks:  req.ClosureReason,
		},
	}
	return r, nil
}

// CancelClaim cancels a claim
// POST /claims/death/{claim_id}/cancel
func (h *ClaimHandler) CancelClaim(sctx *serverRoute.Context, req CancelClaimRequest) (*resp.ClaimCancelResponse, error) {
	// Verify claim exists
	claim, err := h.svc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found with ID: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim by ID: %v", err)
		return nil, err
	}

	// Update claim status to CANCELLED
	updates := map[string]interface{}{
		"status":           "CANCELLED",
		"rejection_reason": &req.CancellationReason,
	}

	_, err = h.svc.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error cancelling claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Claim cancelled with ID: %s, Reason: %s", req.ClaimID, req.CancellationReason)

	r := &resp.ClaimCancelResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ClaimCancelData{
			ClaimID:     claim.ID,
			Status:      "CANCELLED",
			CancelledAt: time.Now().Format("2006-01-02 15:04:05"),
			CancelledBy: req.RequestedBy,
			Reason:      req.CancellationReason,
		},
	}
	return r, nil
}

// GetApprovalQueue retrieves claims pending approval
// GET /claims/death/approval-queue
// Reference: BR-CLM-DC-005
func (h *ClaimHandler) GetApprovalQueue(sctx *serverRoute.Context, req ListClaimsParams) (*resp.DeathClaimsListResponse, error) {
	// Build filters from request parameters
	filters := make(map[string]interface{})
	if req.Status != "" {
		filters["status"] = req.Status
	}

	// Convert uint64 to int64 for repository
	skip := int64(req.Skip)
	limit := int64(req.Limit)

	// Call repository to get approval queue
	results, totalCount, err := h.svc.GetApprovalQueue(sctx.Ctx, filters, skip, limit)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching approval queue: %v", err)
		return nil, err
	}

	// Build response
	r := &resp.DeathClaimsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse: port.NewMetaDataResponse(req.Skip, req.Limit, uint64(totalCount)),
		Data: resp.NewDeathClaimsResponse(results),
	}
	return r, nil
}

// Helper function to safely get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
