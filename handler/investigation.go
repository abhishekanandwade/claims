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

// InvestigationHandler handles investigation-related HTTP requests
type InvestigationHandler struct {
	*serverHandler.Base
	svc         *repo.InvestigationRepository
	claimSvc    *repo.ClaimRepository
	progressSvc *repo.InvestigationProgressRepository
}

// NewInvestigationHandler creates a new investigation handler
func NewInvestigationHandler(
	svc *repo.InvestigationRepository,
	claimSvc *repo.ClaimRepository,
	progressSvc *repo.InvestigationProgressRepository,
) *InvestigationHandler {
	base := serverHandler.New("Investigations").
		SetPrefix("/v1").
		AddPrefix("")
	return &InvestigationHandler{
		Base:        base,
		svc:         svc,
		claimSvc:    claimSvc,
		progressSvc: progressSvc,
	}
}

// Routes defines all routes for this handler
func (h *InvestigationHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Investigation Workflow (10 endpoints)
		serverRoute.POST("/claims/death/:claim_id/investigation/assign-officer", h.AssignInvestigationOfficer).Name("Assign Investigation Officer"),
		serverRoute.GET("/claims/death/pending-investigation", h.GetPendingInvestigationClaims).Name("Get Pending Investigation Claims"),
		serverRoute.GET("/claims/death/:claim_id/investigation/:investigation_id/details", h.GetInvestigationDetails).Name("Get Investigation Details"),
		serverRoute.POST("/claims/death/:claim_id/investigation/:investigation_id/progress-update", h.SubmitInvestigationProgress).Name("Submit Investigation Progress"),
		serverRoute.POST("/claims/death/:claim_id/investigation/:investigation_id/submit-report", h.SubmitInvestigationReport).Name("Submit Investigation Report"),
		serverRoute.POST("/claims/death/:claim_id/investigation/:investigation_id/review", h.ReviewInvestigationReport).Name("Review Investigation Report"),
		serverRoute.POST("/claims/death/:id/investigation/trigger-reinvestigation", h.TriggerReinvestigation).Name("Trigger Reinvestigation"),
		serverRoute.POST("/claims/death/:id/investigation/escalate-sla-breach", h.EscalateInvestigationSLA).Name("Escalate Investigation SLA Breach"),
		serverRoute.POST("/claims/death/:id/manual-review/assign", h.AssignManualReview).Name("Assign Manual Review"),
		serverRoute.POST("/claims/death/:id/reject-fraud", h.RejectClaimForFraud).Name("Reject Claim For Fraud"),
	}
}

// AssignInvestigationOfficer assigns an investigation officer to a claim
// POST /claims/death/{claim_id}/investigation/assign-officer
// Reference: BR-CLM-DC-002 (21-day SLA)
func (h *InvestigationHandler) AssignInvestigationOfficer(sctx *serverRoute.Context, req AssignInvestigationRequest) (*resp.InvestigationAssignedResponse, error) {
	// Verify claim exists and requires investigation
	claim, err := h.claimSvc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding claim: %v", err)
		return nil, err
	}

	if !claim.InvestigationRequired {
		log.Error(sctx.Ctx, "Claim does not require investigation: %s", req.ClaimID)
		return nil, pgx.ErrNoRows
	}

	// Check if investigation already exists
	existingInv, _ := h.svc.FindByClaimID(sctx.Ctx, req.ClaimID)
	if len(existingInv) > 0 {
		// Check if there's an active investigation
		for _, inv := range existingInv {
			if inv.Status == "ASSIGNED" || inv.Status == "IN_PROGRESS" {
				log.Error(sctx.Ctx, "Active investigation already exists: %s", inv.InvestigationID)
				return nil, pgx.ErrNoRows
			}
		}
	}

	// Calculate due date (21 days from assignment)
	// Reference: BR-CLM-DC-002
	assignmentDate := time.Now()
	dueDate := assignmentDate.AddDate(0, 0, 21)

	// Determine assignment type (default to AUTO)
	assignmentType := req.AssignmentType
	if assignmentType == "" {
		assignmentType = "AUTO"
	}
	autoAssigned := assignmentType == "AUTO"

	// TODO: Get assigned_by from user context
	assignedBy := "SYSTEM"

	// TODO: Get investigator rank and jurisdiction from user service
	investigatorRank := "FIELD_OFFICER"
	jurisdiction := "REGIONAL"

	// Create investigation
	investigationData := domain.Investigation{
		InvestigationID:    generateInvestigationID(),
		ClaimID:           req.ClaimID,
		AssignedBy:        assignedBy,
		InvestigatorID:    req.InvestigatorID,
		InvestigatorRank:  &investigatorRank,
		Jurisdiction:      &jurisdiction,
		AutoAssigned:      autoAssigned,
		AssignmentDate:    assignmentDate,
		DueDate:           dueDate,
		Status:            "ASSIGNED",
		ProgressPercentage: 0,
		ReinvestigationCount: 0,
	}

	result, err := h.svc.Create(sctx.Ctx, investigationData)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating investigation: %v", err)
		return nil, err
	}

	// Update claim investigation status
	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"investigation_status": "ASSIGNED",
		"investigator_id":      req.InvestigatorID,
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim investigation status: %v", err)
	}

	log.Info(sctx.Ctx, "Investigation assigned: ID=%s, ClaimID=%s, InvestigatorID=%s, DueDate=%s",
		result.InvestigationID, req.ClaimID, req.InvestigatorID, dueDate.Format("2006-01-02"))

	// Build response
	r := &resp.InvestigationAssignedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: resp.InvestigationAssignmentData{
			InvestigationID: result.InvestigationID,
			ClaimID:         req.ClaimID,
			InvestigatorID:  req.InvestigatorID,
			AssignmentDate:  assignmentDate.Format("2006-01-02 15:04:05"),
			DueDate:         dueDate.Format("2006-01-02"),
			Status:          "ASSIGNED",
			Message:         "Investigation officer assigned successfully",
		},
	}
	return r, nil
}

// GetPendingInvestigationClaims retrieves list of claims pending investigation assignment
// GET /claims/death/pending-investigation
func (h *InvestigationHandler) GetPendingInvestigationClaims(sctx *serverRoute.Context, req GetPendingInvestigationClaimsUri) (*resp.PendingInvestigationsResponse, error) {
	// Build filters for query
	filters := map[string]interface{}{
		"investigation_required": true,
		"status":                 "REGISTERED",
	}

	if req.Jurisdiction != "" {
		filters["jurisdiction"] = req.Jurisdiction
	}

	if req.SLAStatus != "" {
		// TODO: Filter by SLA status
		// This requires calculating SLA for each claim
	}

	// Get claims with pagination
	claims, total, err := h.claimSvc.List(sctx.Ctx, filters, int64(req.Skip), int64(req.Limit), req.OrderBy, req.SortType)
	if err != nil {
		log.Error(sctx.Ctx, "Error retrieving pending investigation claims: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Retrieved %d pending investigation claims (total: %d)", len(claims), total)

	// Build response
	r := resp.NewPendingInvestigationsResponse(claims, total, int(req.Skip), int(req.Limit))
	return &r, nil
}

// GetInvestigationDetails retrieves investigation assignment details
// GET /claims/death/{claim_id}/investigation/{investigation_id}/details
func (h *InvestigationHandler) GetInvestigationDetails(sctx *serverRoute.Context, req InvestigationIDUri) (*resp.InvestigationDetailsResponse, error) {
	// Get investigation
	investigation, err := h.svc.FindByInvestigationID(sctx.Ctx, req.InvestigationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Investigation not found: %s", req.InvestigationID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding investigation: %v", err)
		return nil, err
	}

	// Get claim details
	claim, err := h.claimSvc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		log.Error(sctx.Ctx, "Error finding claim: %v", err)
		return nil, err
	}

	// Calculate SLA status
	slaStatus, daysRemaining := calculateInvestigationSLA(investigation.DueDate, time.Now())

	// Get progress timeline
	startDate := time.Time{} // Zero time for no start date filter
	endDate := time.Time{}   // Zero time for no end date filter
	progressTimeline, _ := h.progressSvc.GetProgressTimeline(sctx.Ctx, req.InvestigationID, startDate, endDate)

	// TODO: Get investigation checklist based on death type
	checklist := getInvestigationChecklist(claim.DeathType)

	log.Info(sctx.Ctx, "Retrieved investigation details: %s", req.InvestigationID)

	// Build response
	r := &resp.InvestigationDetailsResponse{
		StatusCodeAndMessage: port.ReadSuccess,
		Data: resp.InvestigationDetailData{
			InvestigationResponse: resp.NewInvestigationResponse(investigation, slaStatus, daysRemaining),
			ClaimDetails: resp.ClaimSummary{
				ClaimID:              claim.ID,
				ClaimNumber:          claim.ClaimNumber,
				PolicyID:             claim.PolicyID,
				CustomerID:           claim.CustomerID,
				ClaimType:            claim.ClaimType,
				ClaimDate:            claim.ClaimDate.Format("2006-01-02 15:04:05"),
				DeathType:            getStringValue(claim.DeathType),
				DeathDate:            formatTimePtr(claim.DeathDate),
				InvestigationRequired: claim.InvestigationRequired,
				Status:               claim.Status,
			},
			InvestigationChecklist: checklist,
			ProgressTimeline:      convertProgressToResponse(progressTimeline),
		},
	}
	return r, nil
}

// SubmitInvestigationProgress submits investigation progress update (heartbeat)
// POST /claims/death/{claim_id}/investigation/{investigation_id}/progress-update
// Reference: BR-CLM-DC-002 (heartbeat tracking)
func (h *InvestigationHandler) SubmitInvestigationProgress(sctx *serverRoute.Context, req InvestigationProgressRequest) (*resp.InvestigationProgressUpdateResponse, error) {
	// Verify investigation exists
	_, err := h.svc.FindByInvestigationID(sctx.Ctx, req.InvestigationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Investigation not found: %s", req.InvestigationID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding investigation: %v", err)
		return nil, err
	}

	// Create progress record
	progressData := domain.InvestigationProgress{
		InvestigationID:     req.InvestigationID,
		UpdateDate:          time.Now(),
		ProgressPercentage:  req.Percentage,
		Remarks:             req.ProgressNotes,
		EstimatedCompletionDate: nil,
		// TODO: Get updated_by from user context
		UpdatedBy: "INVESTIGATOR_USER_ID",
	}

	progressResult, err := h.progressSvc.Create(sctx.Ctx, progressData)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating progress record: %v", err)
		return nil, err
	}

	// Update investigation progress
	_, err = h.svc.UpdateProgress(sctx.Ctx, req.InvestigationID, req.Percentage)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating investigation progress: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Investigation progress submitted: ID=%s, Percentage=%d",
		progressResult.ID, req.Percentage)

	// Build response
	r := &resp.InvestigationProgressUpdateResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.InvestigationProgressUpdateData{
			ProgressID:      progressResult.ID,
			InvestigationID: req.InvestigationID,
			Percentage:      req.Percentage,
			RecordedAt:      time.Now().Format("2006-01-02 15:04:05"),
			Message:         "Progress recorded successfully",
		},
	}
	return r, nil
}

// SubmitInvestigationReport submits final investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/submit-report
// Reference: BR-CLM-DC-011 (review within 5 days)
func (h *InvestigationHandler) SubmitInvestigationReport(sctx *serverRoute.Context, req SubmitInvestigationReportRequest) (*resp.InvestigationReportSubmittedResponse, error) {
	// Verify investigation exists
	_, err := h.svc.FindByInvestigationID(sctx.Ctx, req.InvestigationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Investigation not found: %s", req.InvestigationID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding investigation: %v", err)
		return nil, err
	}

	// Parse report date for validation
	_, err = time.Parse("2006-01-02", req.ReportDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid report date format: %v", err)
		return nil, err
	}

	// Submit report
	_, err = h.svc.SubmitReport(sctx.Ctx, req.InvestigationID, req.ReportOutcome, req.Findings,
		req.Recommendation, "REPORT_DOC_ID")
	if err != nil {
		log.Error(sctx.Ctx, "Error submitting investigation report: %v", err)
		return nil, err
	}

	// Calculate review due date (5 days from submission)
	// Reference: BR-CLM-DC-011
	reviewDueDate := time.Now().AddDate(0, 0, 5)

	// Update claim investigation status based on outcome
	claimStatus := "UNDER_REVIEW"
	if req.ReportOutcome == "FRAUD" {
		claimStatus = "FRAUD_DETECTED"
	} else if req.ReportOutcome == "CLEAR" {
		claimStatus = "INVESTIGATION_COMPLETE"
	}

	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"investigation_status": claimStatus,
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim investigation status: %v", err)
	}

	log.Info(sctx.Ctx, "Investigation report submitted: ID=%s, Outcome=%s, ReviewDueDate=%s",
		req.InvestigationID, req.ReportOutcome, reviewDueDate.Format("2006-01-02"))

	// Build response
	r := &resp.InvestigationReportSubmittedResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.InvestigationReportData{
			InvestigationID: req.InvestigationID,
			ClaimID:         req.ClaimID,
			ReportOutcome:   req.ReportOutcome,
			SubmittedAt:     time.Now().Format("2006-01-02 15:04:05"),
			Status:          "SUBMITTED_FOR_REVIEW",
			ReviewDueDate:   reviewDueDate.Format("2006-01-02"),
			Message:         "Investigation report submitted successfully. Review due within 5 days.",
		},
	}
	return r, nil
}

// ReviewInvestigationReport reviews investigation report
// POST /claims/death/{claim_id}/investigation/{investigation_id}/review
// Reference: BR-CLM-DC-011 (5-day review SLA)
func (h *InvestigationHandler) ReviewInvestigationReport(sctx *serverRoute.Context, req ReviewInvestigationReportRequest) (*resp.InvestigationReviewResponse, error) {
	// Verify investigation exists
	_, err := h.svc.FindByInvestigationID(sctx.Ctx, req.InvestigationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Investigation not found: %s", req.InvestigationID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding investigation: %v", err)
		return nil, err
	}

	// TODO: Get reviewed_by from user context
	reviewedBy := "REVIEWER_USER_ID"

	// Review report
	_, err = h.svc.ReviewReport(sctx.Ctx, req.InvestigationID, reviewedBy, req.ReviewDecision,
		req.ReviewerRemarks)
	if err != nil {
		log.Error(sctx.Ctx, "Error reviewing investigation report: %v", err)
		return nil, err
	}

	// Determine next action based on review decision
	nextAction := ""
	claimStatus := ""

	switch req.ReviewDecision {
	case "ACCEPT":
		nextAction = "Proceed to claim approval workflow"
		claimStatus = "INVESTIGATION_ACCEPTED"
	case "REINVESTIGATE":
		nextAction = "Trigger reinvestigation workflow"
		claimStatus = "REINVESTIGATION_REQUIRED"
	case "ESCALATE":
		nextAction = "Escalate to higher authority"
		claimStatus = "ESCALATED"
	}

	// Update claim status
	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"investigation_status": claimStatus,
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim investigation status: %v", err)
	}

	log.Info(sctx.Ctx, "Investigation report reviewed: ID=%s, Decision=%s, Reviewer=%s",
		req.InvestigationID, req.ReviewDecision, reviewedBy)

	// Build response
	r := &resp.InvestigationReviewResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.InvestigationReviewData{
			InvestigationID: req.InvestigationID,
			ReviewDecision:  req.ReviewDecision,
			ReviewedAt:      time.Now().Format("2006-01-02 15:04:05"),
			ReviewedBy:      reviewedBy,
			Status:          "REVIEWED",
			Message:         "Investigation report reviewed successfully",
			NextAction:      nextAction,
		},
	}
	return r, nil
}

// TriggerReinvestigation triggers reinvestigation (max 2 times per BR-CLM-DC-013)
// POST /claims/death/{id}/investigation/trigger-reinvestigation
// Reference: BR-CLM-DC-013 (max 2 times, 14-day SLA)
func (h *InvestigationHandler) TriggerReinvestigation(sctx *serverRoute.Context, req TriggerReinvestigationRequest) (*resp.ReinvestigationTriggeredResponse, error) {
	// Get existing investigations for claim
	existingInv, err := h.svc.FindByClaimID(sctx.Ctx, req.ClaimID)
	if err != nil {
		log.Error(sctx.Ctx, "Error finding existing investigations: %v", err)
		return nil, err
	}

	if len(existingInv) == 0 {
		log.Error(sctx.Ctx, "No existing investigation found for claim: %s", req.ClaimID)
		return nil, pgx.ErrNoRows
	}

	// Get latest investigation
	latestInv := existingInv[len(existingInv)-1]

	// Check reinvestigation limit (max 2)
	if latestInv.ReinvestigationCount >= 2 {
		log.Error(sctx.Ctx, "Maximum reinvestigation limit reached (2): %s", req.ClaimID)
		return nil, pgx.ErrNoRows
	}

	// Trigger reinvestigation
	newInv, err := h.svc.TriggerReinvestigation(sctx.Ctx, latestInv.InvestigationID)
	if err != nil {
		log.Error(sctx.Ctx, "Error triggering reinvestigation: %v", err)
		return nil, err
	}

	// Update claim status
	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"investigation_status": "REINVESTIGATION_ASSIGNED",
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim investigation status: %v", err)
	}

	log.Info(sctx.Ctx, "Reinvestigation triggered: ClaimID=%s, NewInvID=%s, Count=%d",
		req.ClaimID, newInv.InvestigationID, latestInv.ReinvestigationCount+1)

	// Build response
	r := &resp.ReinvestigationTriggeredResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: resp.ReinvestigationData{
			NewInvestigationID:     newInv.InvestigationID,
			OriginalInvestigationID: latestInv.InvestigationID,
			ReinvestigationCount:    int(latestInv.ReinvestigationCount + 1),
			NewDueDate:             newInv.DueDate.Format("2006-01-02"),
			Status:                 "ASSIGNED",
			Message:                "Reinvestigation triggered successfully",
		},
	}
	return r, nil
}

// EscalateInvestigationSLA escalates investigation due to SLA breach
// POST /claims/death/{id}/investigation/escalate-sla-breach
// Reference: BR-CLM-DC-002 (escalation hierarchy)
func (h *InvestigationHandler) EscalateInvestigationSLA(sctx *serverRoute.Context, req EscalateInvestigationSLAUri) (*resp.InvestigationSLAEscalationResponse, error) {
	// Get claim
	_, err := h.claimSvc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding claim: %v", err)
		return nil, err
	}

	// Get investigation
	investigations, err := h.svc.FindByClaimID(sctx.Ctx, req.ClaimID)
	if err != nil || len(investigations) == 0 {
		log.Error(sctx.Ctx, "No investigation found for claim: %s", req.ClaimID)
		return nil, pgx.ErrNoRows
	}

	inv := investigations[len(investigations)-1]

	// Determine escalation level based on SLA breach
	escalationLevel := "LEVEL_1"
	escalatedTo := "DIVISION_HEAD"

	if inv.ReinvestigationCount > 0 {
		escalationLevel = "LEVEL_2"
		escalatedTo = "ZONAL_MANAGER"
	}

	// TODO: Send escalation notification to escalated_to

	log.Info(sctx.Ctx, "Investigation SLA breach escalated: ClaimID=%s, InvID=%s, Level=%s, To=%s",
		req.ClaimID, inv.InvestigationID, escalationLevel, escalatedTo)

	// Build response
	r := &resp.InvestigationSLAEscalationResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.InvestigationSLAEscalationData{
			InvestigationID: inv.InvestigationID,
			EscalationLevel: escalationLevel,
			EscalatedTo:     escalatedTo,
			EscalatedAt:     time.Now().Format("2006-01-02 15:04:05"),
			Status:          "ESCALATED",
			Message:         "Investigation escalated due to SLA breach",
		},
	}
	return r, nil
}

// AssignManualReview assigns claim for manual review (SUSPECT outcome)
// POST /claims/death/{id}/manual-review/assign
// Reference: BR-CLM-DC-011 (SUSPECT outcome handling)
func (h *InvestigationHandler) AssignManualReview(sctx *serverRoute.Context, req AssignManualReviewRequest) (*resp.ManualReviewAssignedResponse, error) {
	// Verify claim exists
	_, err := h.claimSvc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding claim: %v", err)
		return nil, err
	}

	// Update claim status for manual review
	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"status":       "MANUAL_REVIEW",
		"approver_id":  req.ReviewerID,
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim for manual review: %v", err)
		return nil, err
	}

	// TODO: Send notification to reviewer

	log.Info(sctx.Ctx, "Manual review assigned: ClaimID=%s, ReviewerID=%s, Priority=%s",
		req.ClaimID, req.ReviewerID, req.Priority)

	// Build response
	r := &resp.ManualReviewAssignedResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ManualReviewAssignmentData{
			ClaimID:    req.ClaimID,
			ReviewerID: req.ReviewerID,
			Priority:   req.Priority,
			AssignedAt: time.Now().Format("2006-01-02 15:04:05"),
			Status:     "ASSIGNED",
			Message:    "Manual review assigned successfully",
		},
	}
	return r, nil
}

// RejectClaimForFraud rejects claim based on fraud investigation
// POST /claims/death/{id}/reject-fraud
// Reference: BR-CLM-DC-020 (fraud rejection)
func (h *InvestigationHandler) RejectClaimForFraud(sctx *serverRoute.Context, req RejectClaimForFraudRequest) (*resp.ClaimRejectedForFraudResponse, error) {
	// Verify claim exists
	claim, err := h.claimSvc.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error finding claim: %v", err)
		return nil, err
	}

	// Update claim status
	rejectionDate := time.Now()
	_, err = h.claimSvc.Update(sctx.Ctx, req.ClaimID, map[string]interface{}{
		"status":             "REJECTED_FRAUD",
		"rejection_code":     "FRAUD_DETECTED",
		"rejection_reason":   req.FraudEvidence,
		"workflow_state":     "CLOSED",
	})
	if err != nil {
		log.Error(sctx.Ctx, "Error updating claim for fraud rejection: %v", err)
		return nil, err
	}

	// TODO: Create legal case if required
	legalCaseID := ""
	if req.LegalActionRequired {
		// TODO: Integrate with legal case management system
		legalCaseID = "LEGAL_" + claim.ClaimNumber
	}

	log.Info(sctx.Ctx, "Claim rejected for fraud: ClaimID=%s, LegalAction=%s, CaseID=%s",
		req.ClaimID, boolToString(req.LegalActionRequired), legalCaseID)

	// Build response
	r := &resp.ClaimRejectedForFraudResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.FraudRejectionData{
			ClaimID:              req.ClaimID,
			RejectionDate:        rejectionDate.Format("2006-01-02 15:04:05"),
			InvestigationReportID: req.InvestigationReportID,
			LegalActionRequired:  req.LegalActionRequired,
			LegalCaseID:          legalCaseID,
			Status:               "REJECTED",
			Message:              "Claim rejected due to fraud",
		},
	}
	return r, nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// generateInvestigationID generates a unique investigation ID
func generateInvestigationID() string {
	// TODO: Implement proper ID generation
	// Format: INV{YYYY}{DDDD}
	return "INV2025001"
}

// calculateInvestigationSLA calculates SLA status and days remaining
func calculateInvestigationSLA(dueDate time.Time, currentDate time.Time) (string, int) {
	if dueDate.IsZero() {
		return "UNKNOWN", 0
	}

	totalDuration := dueDate.Sub(currentDate)
	daysRemaining := int(totalDuration.Hours() / 24)

	if daysRemaining < 0 {
		return "RED", 0
	}

	// Calculate percentage of SLA used
	// Assuming 21-day SLA for investigation
	slaDays := 21
	daysUsed := slaDays - daysRemaining
	percentageUsed := float64(daysUsed) / float64(slaDays) * 100

	switch {
	case percentageUsed < 70:
		return "GREEN", daysRemaining
	case percentageUsed < 90:
		return "YELLOW", daysRemaining
	case percentageUsed < 100:
		return "ORANGE", daysRemaining
	default:
		return "RED", daysRemaining
	}
}

// getInvestigationChecklist returns investigation checklist based on death type
func getInvestigationChecklist(deathType *string) []string {
	// TODO: Implement dynamic checklist based on death type
	// Reference: BR-CLM-DC-011
	baseChecklist := []string{
		"Verify death certificate",
		"Confirm cause of death",
		"Check hospital records",
		"Interview family members",
		"Verify claimant identity",
		"Check policy status",
	}

	if deathType != nil {
		switch *deathType {
		case "ACCIDENTAL":
			baseChecklist = append(baseChecklist, []string{
				"Obtain FIR copy",
				"Verify police report",
				"Check post-mortem report",
			}...)
		case "UNNATURAL":
			baseChecklist = append(baseChecklist, []string{
				"Obtain post-mortem report",
				"Verify police investigation",
				"Check forensic evidence",
			}...)
		case "SUICIDE":
			baseChecklist = append(baseChecklist, []string{
				"Verify suicide note (if any)",
				"Check police investigation",
				"Verify mental health records",
			}...)
		}
	}

	return baseChecklist
}

// convertProgressToResponse converts domain progress to response format
func convertProgressToResponse(progress []domain.InvestigationProgress) []resp.InvestigationProgress {
	result := make([]resp.InvestigationProgress, len(progress))
	for i, p := range progress {
		result[i] = resp.InvestigationProgress{
			ID:                   p.ID,
			InvestigationID:      p.InvestigationID,
			ProgressNotes:        p.Remarks,
			NextSteps:            "", // Not in domain model
			Percentage:           p.ProgressPercentage,
			ChecklistItemsCompleted: p.ChecklistItemsCompleted,
			EstimatedCompletionDate: formatTimePtr(p.EstimatedCompletionDate),
			RecordedAt:           p.UpdateDate.Format("2006-01-02 15:04:05"),
			RecordedBy:           p.UpdatedBy,
		}
	}
	return result
}

// Helper function to format time pointer as string
func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// Helper function to convert bool to string
func boolToString(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}
