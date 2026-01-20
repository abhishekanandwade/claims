package handler

import (
	"context"
	"fmt"
	"time"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
)

// OmbudsmanHandler handles ombudsman complaint operations
// FR-CLM-OMB-001 to FR-CLM-OMB-004: Ombudsman complaint management
// BR-CLM-OMB-001 to BR-CLM-OMB-008: Ombudsman business rules
type OmbudsmanHandler struct {
	*serverHandler.Base
	ombudsmanRepo *repo.OmbudsmanComplaintRepository
}

// NewOmbudsmanHandler creates a new OmbudsmanHandler instance
func NewOmbudsmanHandler(ombudsmanRepo *repo.OmbudsmanComplaintRepository) *OmbudsmanHandler {
	base := serverHandler.New("Ombudsman").
		SetPrefix("/v1").
		AddPrefix("")
	return &OmbudsmanHandler{
		Base:          base,
		ombudsmanRepo: ombudsmanRepo,
	}
}

// RegisterRoutes registers Ombudsman routes with the server
func (h *OmbudsmanHandler) RegisterRoutes(router serverHandler.RouterAPI) {
	// ==================== COMPLAINT INTAKE & REGISTRATION ====================
	router.POST("/ombudsman/complaint/submit", h.SubmitComplaint, nil, &serverRoute.RouteConfig{
		Summary:   "Submit new ombudsman complaint",
		Reference: "FR-CLM-OMB-001, BR-CLM-OMB-001",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== COMPLAINT DETAILS & TRACKING ====================
	router.GET("/ombudsman/complaint/{complaint_id}/details", h.GetComplaintDetails, nil, &serverRoute.RouteConfig{
		Summary:   "Get complaint details by ID",
		Reference: "FR-CLM-OMB-001",
		Tags:      []string{"Ombudsman"},
	})

	router.GET("/ombudsman/complaint/{complaint_id}/timeline", h.GetComplaintTimeline, nil, &serverRoute.RouteConfig{
		Summary:   "Get complaint timeline/history",
		Reference: "FR-CLM-OMB-001",
		Tags:      []string{"Ombudsman"},
	})

	router.GET("/ombudsman/complaints/list", h.ListComplaints, nil, &serverRoute.RouteConfig{
		Summary:   "List all complaints with filters",
		Reference: "FR-CLM-OMB-001",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== JURISDICTION & ASSIGNMENT ====================
	router.POST("/ombudsman/complaint/{complaint_id}/assign", h.AssignOmbudsman, nil, &serverRoute.RouteConfig{
		Summary:   "Assign ombudsman to complaint",
		Reference: "FR-CLM-OMB-002, BR-CLM-OMB-002, BR-CLM-OMB-003",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== ADMISSIBILITY REVIEW ====================
	router.POST("/ombudsman/complaint/{complaint_id}/admissibility", h.ReviewAdmissibility, nil, &serverRoute.RouteConfig{
		Summary:   "Review complaint admissibility",
		Reference: "BR-CLM-OMB-001",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== MEDIATION & HEARING ====================
	router.POST("/ombudsman/complaint/{complaint_id}/mediation", h.RecordMediation, nil, &serverRoute.RouteConfig{
		Summary:   "Record mediation outcome",
		Reference: "FR-CLM-OMB-003, BR-CLM-OMB-004",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== AWARD ISSUANCE & ENFORCEMENT ====================
	router.POST("/ombudsman/complaint/{complaint_id}/award", h.IssueAward, nil, &serverRoute.RouteConfig{
		Summary:   "Issue award (mediation or adjudication)",
		Reference: "FR-CLM-OMB-004, BR-CLM-OMB-005",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== COMPLIANCE MONITORING ====================
	router.POST("/ombudsman/complaint/{complaint_id}/compliance", h.RecordCompliance, nil, &serverRoute.RouteConfig{
		Summary:   "Record insurer compliance with award",
		Reference: "BR-CLM-OMB-006",
		Tags:      []string{"Ombudsman"},
	})

	router.GET("/ombudsman/compliance/queue", h.GetComplianceQueue, nil, &serverRoute.RouteConfig{
		Summary:   "Get complaints requiring compliance tracking",
		Reference: "BR-CLM-OMB-006",
		Tags:      []string{"Ombudsman"},
	})

	// ==================== ESCALATION & CLOSURE ====================
	router.POST("/ombudsman/complaint/{complaint_id}/escalate-irdai", h.EscalateToIRDAI, nil, &serverRoute.RouteConfig{
		Summary:   "Escalate non-compliance to IRDAI",
		Reference: "BR-CLM-OMB-006",
		Tags:      []string{"Ombudsman"},
	})

	router.POST("/ombudsman/complaint/{complaint_id}/close", h.CloseComplaint, nil, &serverRoute.RouteConfig{
		Summary:   "Close ombudsman complaint",
		Reference: "BR-CLM-OMB-007",
		Tags:      []string{"Ombudsman"},
	})
}

// ==================== COMPLAINT INTAKE & REGISTRATION ====================

// SubmitComplaint handles new ombudsman complaint submission
// FR-CLM-OMB-001: Complaint Intake & Registration
// BR-CLM-OMB-001: Admissibility checks (representation to insurer, 30-day wait, 1-year limitation, ₹50 lakh cap, no parallel litigation)
func (h *OmbudsmanHandler) SubmitComplaint(sctx *serverRoute.Context, req *SubmitOmbudsmanComplaintRequest) (*response.ComplaintRegisteredResponse, error) {
	log.Info(sctx.Ctx, "Submitting ombudsman complaint: policy_number=%s, complaint_category=%s", req.PolicyNumber, req.ComplaintCategory)

	// TODO: Check for duplicate complaints (same policy + complainant within 30 days)

	// Parse dates
	incidentDate, err := time.Parse("2006-01-02", req.IncidentDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid incident_date format: %v", err)
		return nil, fmt.Errorf("invalid incident_date format: %w", err)
	}

	representationDate, err := time.Parse("2006-01-02", req.RepresentationDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid representation_date format: %v", err)
		return nil, fmt.Errorf("invalid representation_date format: %w", err)
	}

	// BR-CLM-OMB-001: Check 30-day wait period (representation to insurer)
	daysSinceRepresentation := int(time.Since(representationDate).Hours() / 24)
	if daysSinceRepresentation < 30 {
		log.Error(sctx.Ctx, "Complaint submitted before 30-day wait period: days_since=%d", daysSinceRepresentation)
		// Still allow submission but flag for review - actual validation in admissibility check
	}

	// TODO: Map jurisdiction based on complainant location (BR-CLM-OMB-002)
	// For now, use placeholder - will be updated in AssignOmbudsman
	ombudsmanCenter := "PENDING_JURISDICTION_MAPPING"

	// Calculate resolution due date (statutory timeline - varies by complaint type)
	// Typically 30-90 days from registration depending on complexity
	resolutionDueDate := time.Now().AddDate(0, 0, 90) // 90-day default timeline

	// Create complaint domain object
	complaint := &domain.OmbudsmanComplaint{
		ComplainantName:        req.ComplainantName,
		ComplainantAddress:     req.ComplainantAddress,
		ComplainantMobile:      req.ComplainantMobile,
		ComplainantEmail:       &req.ComplainantEmail,
		ComplainantRole:        req.ComplainantRole,
		LanguagePreference:     req.LanguagePreference,
		IDProofType:            req.IDProofType,
		IDProofNumber:          req.IDProofNumber,
		PolicyID:               req.PolicyNumber,
		ClaimID:                &req.ClaimNumber,
		PolicyType:             req.PolicyType,
		AgentName:              &req.AgentName,
		AgentBranch:            &req.AgentBranch,
		ComplaintCategory:      req.ComplaintCategory,
		IncidentDate:           incidentDate,
		RepresentationDate:     representationDate,
		IssueDescription:       req.IssueDescription,
		ReliefSought:           req.ReliefSought,
		ClaimValue:             &req.ClaimValue,
		OmbudsmanCenter:        &ombudsmanCenter,
		Status:                 "REGISTERED",
		Admissible:             nil, // Will be determined in admissibility review
		AdmissibilityChecked:   false,
		ParallelLitigation:     req.ParallelLitigation, // BR-CLM-OMB-001: No parallel litigation
		IsEmergency:            req.IsEmergency,
		Channel:                req.Channel,
		ResolutionDueDate:      &resolutionDueDate,
		AcknowledgementSent:    false,
		AcknowledgementSentDate: nil,
		EscalatedToIRDAI:       false,
		RetentionPeriod:        7, // BR-CLM-OMB-007: 7 years for mediation (10 for awards)
	}

	// Generate complaint number (OMB{YYYY}{DDDD})
	complaintNumber, err := h.ombudsmanRepo.GenerateComplaintNumber(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate complaint number: %v", err)
		return nil, fmt.Errorf("failed to generate complaint number: %w", err)
	}
	complaint.ComplaintNumber = complaintNumber

	// Create complaint in database
	createdComplaint, err := h.ombudsmanRepo.Create(sctx.Ctx, *complaint)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to create complaint: %v", err)
		return nil, fmt.Errorf("failed to create complaint: %w", err)
	}

	log.Info(sctx.Ctx, "Complaint created successfully: complaint_id=%s, complaint_number=%s", createdComplaint.ComplaintID, createdComplaint.ComplaintNumber)

	// TODO: Send acknowledgement via SMS/email within 24 hours (FR-CLM-OMB-001)
	// TODO: Upload attachments to ECMS and link to complaint

	// BR-CLM-OMB-002: Auto-map to jurisdiction center
	// TODO: Call jurisdiction mapping service

	// Build response
	resp := &response.ComplaintRegisteredResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintNumber:       createdComplaint.ComplaintNumber,
		ComplaintID:           createdComplaint.ComplaintID,
		AcknowledgementSent:   false, // TODO: Set to true after sending acknowledgement
		AcknowledgementDate:   "",    // TODO: Set after sending acknowledgement
		Status:                createdComplaint.Status,
		AssignedJurisdiction:  ombudsmanCenter,
		NextSteps:             "Your complaint has been registered. You will receive an acknowledgement within 24 hours. The ombudsman center will review your complaint and contact you for further proceedings.",
	}

	return resp, nil
}

// ==================== COMPLAINT DETAILS & TRACKING ====================

// GetComplaintDetails retrieves full complaint details
func (h *OmbudsmanHandler) GetComplaintDetails(sctx *serverRoute.Context, req *ComplaintIDUri) (*response.ComplaintDetailsResponse, error) {
	log.Info(sctx.Ctx, "Getting complaint details: complaint_id=%s", req.ComplaintID)

	// Fetch complaint from database
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, req.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Build response
	resp := response.NewComplaintDetailsResponse(complaint)
	return resp, nil
}

// GetComplaintTimeline retrieves complaint history/timeline
func (h *OmbudsmanHandler) GetComplaintTimeline(sctx *serverRoute.Context, req *ComplaintIDUri) (*response.ComplaintTimelineResponse, error) {
	log.Info(sctx.Ctx, "Getting complaint timeline: complaint_id=%s", req.ComplaintID)

	// TODO: Implement timeline retrieval from audit trail / claim_history table
	// For now, return empty timeline
	timeline := []response.TimelineEntry{}

	// Build response
	resp := &response.ComplaintTimelineResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          req.ComplaintID,
		Timeline:             timeline,
		TotalEvents:          len(timeline),
	}

	return resp, nil
}

// ListComplaints retrieves paginated list of complaints with filters
func (h *OmbudsmanHandler) ListComplaints(sctx *serverRoute.Context, req *MetadataRequest) (*response.ComplaintsListResponse, error) {
	log.Info(sctx.Ctx, "Listing complaints: skip=%d, limit=%d", req.Skip, req.Limit)

	// TODO: Extract filter parameters from query string
	// Filters: status, ombudsman_center, complaint_category, admissible, date ranges

	// For now, fetch all complaints with pagination
	filters := map[string]interface{}{}
	complaints, total, err := h.ombudsmanRepo.List(sctx.Ctx, filters, req.Skip, req.Limit)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to list complaints: %v", err)
		return nil, fmt.Errorf("failed to list complaints: %w", err)
	}

	// Build response
	resp := response.NewComplaintsListResponse(complaints, total, req.Skip, req.Limit)
	return resp, nil
}

// ==================== JURISDICTION & ASSIGNMENT ====================

// AssignOmbudsman assigns an ombudsman to a complaint
// FR-CLM-OMB-002: Jurisdiction Mapping
// BR-CLM-OMB-002: Map complaint to territorial ombudsman center
// BR-CLM-OMB-003: Conflict of interest screening
func (h *OmbudsmanHandler) AssignOmbudsman(sctx *serverRoute.Context, uri *ComplaintIDUri, req *AssignOmbudsmanRequest) (*response.OmbudsmanAssignedResponse, error) {
	log.Info(sctx.Ctx, "Assigning ombudsman: complaint_id=%s, ombudsman_id=%s, ombudsman_center=%s", uri.ComplaintID, req.OmbudsmanID, req.OmbudsmanCenter)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Check if conflict detected
	conflictDetected := false
	reassignmentRequired := false
	assignmentRemarks := "Assigned to ombudsman"

	if req.ConflictCheck {
		// TODO: BR-CLM-OMB-003: Check for conflict of interest
		// Check for: prior relationship with complainant, financial interest, duplicate litigation
		// For now, assume no conflict
		conflictDetected = false
		if conflictDetected {
			reassignmentRequired = true
			assignmentRemarks = "Conflict detected - reassignment required"
		}
	}

	// Update complaint with ombudsman assignment
	updates := map[string]interface{}{
		"assigned_ombudsman_id": req.OmbudsmanID,
		"ombudsman_center":      req.OmbudsmanCenter,
		"status":                "ASSIGNED_TO_JURISDICTION",
	}

	updatedComplaint, err := h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to assign ombudsman: %v", err)
		return nil, fmt.Errorf("failed to assign ombudsman: %w", err)
	}

	log.Info(sctx.Ctx, "Ombudsman assigned successfully: complaint_id=%s, ombudsman_id=%s", updatedComplaint.ComplaintID, req.OmbudsmanID)

	// Build response
	resp := &response.OmbudsmanAssignedResponse{
		StatusCodeAndMessage:  port.Success(),
		ComplaintID:            uri.ComplaintID,
		AssignedOmbudsmanID:    req.OmbudsmanID,
		OmbudsmanCenter:        req.OmbudsmanCenter,
		ConflictCheckPerformed: req.ConflictCheck,
		ConflictDetected:       conflictDetected,
		ReassignmentRequired:   reassignmentRequired,
		AssignmentDate:         time.Now().Format("2006-01-02 15:04:05"),
		AssignmentRemarks:      assignmentRemarks,
	}

	return resp, nil
}

// ==================== ADMISSIBILITY REVIEW ====================

// ReviewAdmissibility reviews and updates complaint admissibility
// BR-CLM-OMB-001: Admissibility checks (representation to insurer, 30-day wait, 1-year limitation, ₹50 lakh cap, no parallel litigation)
func (h *OmbudsmanHandler) ReviewAdmissibility(sctx *serverRoute.Context, uri *ComplaintIDUri, req *ReviewAdmissibilityRequest) (*response.AdmissibilityReviewedResponse, error) {
	log.Info(sctx.Ctx, "Reviewing admissibility: complaint_id=%s, admissible=%t", uri.ComplaintID, req.Admissible)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Update admissibility
	updates := map[string]interface{}{
		"admissible":              &req.Admissible,
		"admissibility_checked":   true,
		"admissibility_reason":    req.AdmissibilityReason,
		"inadmissibility_reason":  &req.InadmissibilityReason,
	}

	// Update status based on admissibility
	if req.Admissible {
		updates["status"] = "ADMISSIBLE"
	} else {
		updates["status"] = "INADMISSIBLE"
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to update admissibility: %v", err)
		return nil, fmt.Errorf("failed to update admissibility: %w", err)
	}

	log.Info(sctx.Ctx, "Admissibility updated successfully: complaint_id=%s, admissible=%t", uri.ComplaintID, req.Admissible)

	// Build response
	nextSteps := ""
	if req.Admissible {
		nextSteps = "Complaint is admissible. Proceeding with hearing and mediation process."
	} else {
		nextSteps = "Complaint is inadmissible. Complainant will be informed with reasons. Complaint may be closed or appealed."
	}

	resp := &response.AdmissibilityReviewedResponse{
		StatusCodeAndMessage:   port.Success(),
		ComplaintID:            uri.ComplaintID,
		Admissible:             req.Admissible,
		AdmissibilityReason:    req.AdmissibilityReason,
		InadmissibilityReason:  req.InadmissibilityReason,
		ReviewedBy:             req.ReviewedBy,
		ReviewDate:             time.Now().Format("2006-01-02 15:04:05"),
		NextSteps:              nextSteps,
	}

	return resp, nil
}

// ==================== MEDIATION & HEARING ====================

// RecordMediation records mediation outcome
// FR-CLM-OMB-003: Hearing Scheduling & Management
// BR-CLM-OMB-004: Mediation recommendation (Rule 16)
func (h *OmbudsmanHandler) RecordMediation(sctx *serverRoute.Context, uri *ComplaintIDUri, req *RecordMediationRequest) (*response.MediationRecordedResponse, error) {
	log.Info(sctx.Ctx, "Recording mediation: complaint_id=%s, mediation_successful=%t", uri.ComplaintID, req.MediationSuccessful)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// TODO: Record mediation in hearing/mediation table (separate table)

	// Update complaint status based on mediation outcome
	status := "MEDIATION_COMPLETED"
	nextSteps := ""

	if req.MediationSuccessful && req.ComplainantAccepted && req.InsurerAccepted {
		status = "MEDIATION_SETTLED"
		nextSteps = "Mediation successful. Both parties accepted. Mediation recommendation will be issued."
	} else if !req.MediationSuccessful {
		status = "MEDIATION_FAILED"
		nextSteps = "Mediation unsuccessful. Proceeding to adjudication (award)."
	} else {
		nextSteps = "Mediation recorded. Awaiting acceptance from parties."
	}

	updates := map[string]interface{}{
		"status": status,
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to record mediation: %v", err)
		return nil, fmt.Errorf("failed to record mediation: %w", err)
	}

	log.Info(sctx.Ctx, "Mediation recorded successfully: complaint_id=%s, status=%s", uri.ComplaintID, status)

	// Build response
	mediationDate, _ := time.Parse("2006-01-02", req.MediationDate)

	resp := &response.MediationRecordedResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          uri.ComplaintID,
		HearingID:            req.HearingID,
		MediationDate:        mediationDate.Format("2006-01-02 15:04:05"),
		ConsentToMediate:     req.ConsentToMediate,
		MediationSuccessful:  req.MediationSuccessful,
		SettlementTerms:      req.SettlementTerms,
		ComplainantAccepted:  req.ComplainantAccepted,
		InsurerAccepted:      req.InsurerAccepted,
		RecordingOfficer:     req.RecordingOfficer,
		RecordingDate:        time.Now().Format("2006-01-02 15:04:05"),
		Remarks:              req.Remarks,
		NextSteps:            nextSteps,
	}

	return resp, nil
}

// ==================== AWARD ISSUANCE & ENFORCEMENT ====================

// IssueAward issues an award (mediation recommendation or adjudication award)
// FR-CLM-OMB-004: Award Issuance & Enforcement
// BR-CLM-OMB-005: Award issuance with ₹50 lakh cap
func (h *OmbudsmanHandler) IssueAward(sctx *serverRoute.Context, uri *ComplaintIDUri, req *IssueAwardRequest) (*response.AwardIssuedResponse, error) {
	log.Info(sctx.Ctx, "Issuing award: complaint_id=%s, award_type=%s, award_amount=%.2f", uri.ComplaintID, req.AwardType, req.TotalAwardAmount)

	// BR-CLM-OMB-005: Check ₹50 lakh cap
	if req.TotalAwardAmount > 5000000 {
		log.Error(sctx.Ctx, "Award amount exceeds ₹50 lakh cap: amount=%.2f", req.TotalAwardAmount)
		return nil, fmt.Errorf("award amount cannot exceed ₹50 lakh (5,000,000) as per BR-CLM-OMB-005")
	}

	// Parse dates
	digitalSignatureDate, err := time.Parse("2006-01-02", req.DigitalSignatureDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid digital_signature_date format: %v", err)
		return nil, fmt.Errorf("invalid digital_signature_date format: %w", err)
	}

	complianceDeadline, err := time.Parse("2006-01-02", req.ComplianceDeadline)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid compliance_deadline format: %v", err)
		return nil, fmt.Errorf("invalid compliance_deadline format: %w", err)
	}

	// BR-CLM-OMB-006: 30-day compliance timeline
	expectedComplianceDeadline := time.Now().AddDate(0, 0, 30)
	daysDiff := int(complianceDeadline.Sub(expectedComplianceDeadline).Hours() / 24)
	if daysDiff < -7 || daysDiff > 7 {
		log.Warn(sctx.Ctx, "Compliance deadline deviates from 30-day standard: days_diff=%d", daysDiff)
	}

	// TODO: Issue award in award table (separate table)
	awardID := fmt.Sprintf("AWD%s%04d", time.Now().Format("20060102"), 1) // Placeholder

	// Update complaint with award details
	updates := map[string]interface{}{
		"status":                 "AWARD_ISSUED",
		"award_type":             req.AwardType,
		"award_amount":           req.AwardAmount,
		"total_award_amount":     req.TotalAwardAmount,
		"digital_signature":      req.DigitalSignature,
		"digital_signature_date": digitalSignatureDate,
		"compliance_deadline":    complianceDeadline,
		"retention_period":       10, // BR-CLM-OMB-007: 10 years for awards
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to issue award: %v", err)
		return nil, fmt.Errorf("failed to issue award: %w", err)
	}

	log.Info(sctx.Ctx, "Award issued successfully: complaint_id=%s, award_id=%s", uri.ComplaintID, awardID)

	// TODO: Send award document to complainant and insurer
	// TODO: Schedule reminders for compliance (Day 15, 7, 2 before deadline)

	// Build response
	resp := &response.AwardIssuedResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          uri.ComplaintID,
		AwardData: response.AwardData{
			AwardID:              awardID,
			AwardType:            req.AwardType,
			AwardAmount:          req.AwardAmount,
			AwardCurrency:        req.AwardCurrency,
			InterestRate:         req.InterestRate,
			InterestAmount:       req.InterestAmount,
			TotalAwardAmount:     req.TotalAwardAmount,
			AwardReasoning:       req.AwardReasoning,
			DigitalSignatureHash: req.DigitalSignature,
			DigitalSignatureDate: digitalSignatureDate.Format("2006-01-02 15:04:05"),
			IssuedBy:             req.IssuedBy,
			IssuedDate:           time.Now().Format("2006-01-02 15:04:05"),
			ComplianceDeadline:   complianceDeadline.Format("2006-01-02 15:04:05"),
			SupportingDocuments:  req.SupportingDocuments,
			// DocumentURL:         TODO - ECMS document URL
			Status:           "ISSUED",
			BindingOnInsurer: true, // BR-CLM-OMB-005: Award is binding on insurer
		},
	}

	return resp, nil
}

// ==================== COMPLIANCE MONITORING ====================

// RecordCompliance records insurer compliance with award
// BR-CLM-OMB-006: Insurer compliance monitoring (30-day timeline)
func (h *OmbudsmanHandler) RecordCompliance(sctx *serverRoute.Context, uri *ComplaintIDUri, req *RecordComplianceRequest) (*response.ComplianceRecordedResponse, error) {
	log.Info(sctx.Ctx, "Recording compliance: complaint_id=%s, compliance_status=%s", uri.ComplaintID, req.ComplianceStatus)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Parse compliance date
	complianceDate, err := time.Parse("2006-01-02", req.ComplianceDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid compliance_date format: %v", err)
		return nil, fmt.Errorf("invalid compliance_date format: %w", err)
	}

	// Calculate days to comply (or overdue)
	daysToComply := 0
	overdue := false
	if complaint.ComplianceDeadline != nil {
		daysToComply = int(complianceDate.Sub(*complaint.ComplianceDeadline).Hours() / 24)
		overdue = daysToComply > 0
	}

	// Update complaint status based on compliance
	status := "COMPLIANCE_RECORDED"
	nextSteps := ""

	switch req.ComplianceStatus {
	case "ACCEPTED", "PAYMENT_INITIATED", "PAYMENT_COMPLETED":
		status = "COMPLIED"
		nextSteps = "Insurer complied with award. Complaint will be closed after verification."
	case "OBJECTION_FILED":
		status = "OBJECTION_PENDING"
		nextSteps = "Insurer objected to award. Review required. May escalate to IRDAI."
	case "ESCALATED":
		status = "ESCALATED_TO_IRDAI"
		nextSteps = "Complaint escalated to IRDAI for further action."
	default:
		nextSteps = "Compliance recorded. Status updated."
	}

	updates := map[string]interface{}{
		"status":            status,
		"compliance_status": req.ComplianceStatus,
		"compliance_date":   complianceDate,
	}

	if req.PaymentReference != "" {
		updates["payment_reference"] = req.PaymentReference
	}
	if req.PaymentAmount > 0 {
		paymentAmount := req.PaymentAmount
		updates["payment_amount"] = paymentAmount
	}
	if req.ObjectionReason != "" {
		updates["objection_reason"] = req.ObjectionReason
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to record compliance: %v", err)
		return nil, fmt.Errorf("failed to record compliance: %w", err)
	}

	log.Info(sctx.Ctx, "Compliance recorded successfully: complaint_id=%s, status=%s", uri.ComplaintID, status)

	// Build response
	resp := &response.ComplianceRecordedResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          uri.ComplaintID,
		// AwardID:              TODO - Fetch from award table
		ComplianceStatus:  req.ComplianceStatus,
		ComplianceDate:    complianceDate.Format("2006-01-02 15:04:05"),
		PaymentReference:  req.PaymentReference,
		PaymentAmount:     req.PaymentAmount,
		ObjectionReason:   req.ObjectionReason,
		DaysToComply:      daysToComply,
		Overdue:           overdue,
		RecordedBy:        req.RecordedBy,
		RecordingDate:     time.Now().Format("2006-01-02 15:04:05"),
		NextSteps:         nextSteps,
	}

	return resp, nil
}

// GetComplianceQueue retrieves complaints requiring compliance tracking
// BR-CLM-OMB-006: 30-day compliance monitoring with reminders
func (h *OmbudsmanHandler) GetComplianceQueue(sctx *serverRoute.Context, req *MetadataRequest) (*response.ComplianceQueueResponse, error) {
	log.Info(sctx.Ctx, "Getting compliance queue: skip=%d, limit=%d", req.Skip, req.Limit)

	// TODO: Fetch complaints with award issued and pending compliance
	// For now, return empty queue
	complianceItems := []response.ComplianceItem{}
	queueSummary := response.ComplianceQueueSummary{
		TotalPending:         0,
		TotalComplied:        0,
		TotalOverdue:         0,
		TotalEscalated:       0,
		AverageComplianceDays: 0.0,
	}

	// Build response
	resp := &response.ComplianceQueueResponse{
		StatusCodeAndMessage: port.Success(),
		MetaDataResponse:     port.NewMetaDataResponse(int64(len(complianceItems)), req.Skip, req.Limit),
		Complaints:           complianceItems,
		QueueSummary:         queueSummary,
	}

	return resp, nil
}

// ==================== ESCALATION & CLOSURE ====================

// EscalateToIRDAI escalates non-compliance to IRDAI
// BR-CLM-OMB-006: Escalate to IRDAI on breach
func (h *OmbudsmanHandler) EscalateToIRDAI(sctx *serverRoute.Context, uri *ComplaintIDUri, req *EscalateToIRDAIRequest) (*response.IRDAIEscalationResponse, error) {
	log.Info(sctx.Ctx, "Escalating to IRDAI: complaint_id=%s, days_overdue=%d", uri.ComplaintID, req.DaysOverdue)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Parse escalation date
	escalationDate, err := time.Parse("2006-01-02", req.EscalationDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid escalation_date format: %v", err)
		return nil, fmt.Errorf("invalid escalation_date format: %w", err)
	}

	// Generate escalation ID
	escalationID := fmt.Sprintf("ESC%s%04d", time.Now().Format("20060102"), 1) // Placeholder

	// Update complaint with escalation details
	updates := map[string]interface{}{
		"status":              "ESCALATED_TO_IRDAI",
		"escalated_to_irdai":  true,
		"escalation_id":       escalationID,
		"irdai_reference":     &req.IRDAIReference,
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to escalate to IRDAI: %v", err)
		return nil, fmt.Errorf("failed to escalate to IRDAI: %w", err)
	}

	log.Info(sctx.Ctx, "Escalated to IRDAI successfully: complaint_id=%s, escalation_id=%s", uri.ComplaintID, escalationID)

	// TODO: Send escalation report to IRDAI via integration

	// Build response
	resp := &response.IRDAIEscalationResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          uri.ComplaintID,
		// AwardID:              TODO - Fetch from award table
		EscalationID:     escalationID,
		EscalationReason: req.EscalationReason,
		BreachDetails:    req.BreachDetails,
		DaysOverdue:      req.DaysOverdue,
		EscalationDate:   escalationDate.Format("2006-01-02 15:04:05"),
		EscalatedBy:      req.EscalatedBy,
		IRDAIReference:   req.IRDAIReference,
		Status:           "ESCALATED",
	}

	return resp, nil
}

// CloseComplaint closes an ombudsman complaint
// BR-CLM-OMB-007: Complaint closure & archival
func (h *OmbudsmanHandler) CloseComplaint(sctx *serverRoute.Context, uri *ComplaintIDUri, req *CloseComplaintRequest) (*response.ComplaintClosedResponse, error) {
	log.Info(sctx.Ctx, "Closing complaint: complaint_id=%s, closure_type=%s", uri.ComplaintID, req.ClosureType)

	// Fetch complaint
	complaint, err := h.ombudsmanRepo.FindByID(sctx.Ctx, uri.ComplaintID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch complaint: %v", err)
		return nil, fmt.Errorf("failed to fetch complaint: %w", err)
	}

	// Update complaint as closed
	closedDate := time.Now()
	archivalDate := closedDate.AddDate(req.RetentionPeriod, 0, 0) // BR-CLM-OMB-007: Retention period

	updates := map[string]interface{}{
		"status":            "CLOSED",
		"closure_reason":    req.ClosureReason,
		"closure_type":      req.ClosureType,
		"closed_date":       closedDate,
		"archival_date":     archivalDate,
		"retention_period":  req.RetentionPeriod,
	}

	_, err = h.ombudsmanRepo.Update(sctx.Ctx, uri.ComplaintID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to close complaint: %v", err)
		return nil, fmt.Errorf("failed to close complaint: %w", err)
	}

	log.Info(sctx.Ctx, "Complaint closed successfully: complaint_id=%s, archival_date=%s", uri.ComplaintID, archivalDate.Format("2006-01-02"))

	// TODO: Archive complaint documents and audit logs

	// Build response
	resp := &response.ComplaintClosedResponse{
		StatusCodeAndMessage: port.Success(),
		ComplaintID:          uri.ComplaintID,
		ClosureReason:        req.ClosureReason,
		ClosureType:          req.ClosureType,
		ClosedDate:           closedDate.Format("2006-01-02 15:04:05"),
		RetentionPeriod:      req.RetentionPeriod,
		ArchivalDate:         archivalDate.Format("2006-01-02"),
		ClosedBy:             req.ClosedBy,
	}

	return resp, nil
}

