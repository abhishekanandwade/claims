package handler

import (
	"errors"
	"fmt"
	"time"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/handler/response"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
)

// FreeLookHandler handles free look cancellation and policy bond tracking requests
type FreeLookHandler struct {
	*serverHandler.Base
	policyBondRepo       *PolicyBondTrackingRepository
	freelookRepo         *FreeLookCancellationRepository
}

// NewFreeLookHandler creates a new free look handler
func NewFreeLookHandler(
	policyBondRepo *PolicyBondTrackingRepository,
	freelookRepo *FreeLookCancellationRepository,
) *FreeLookHandler {
	base := serverHandler.New("FreeLook").
		SetPrefix("/v1").
		AddPrefix("")

	return &FreeLookHandler{
		Base:             base,
		policyBondRepo:   policyBondRepo,
		freelookRepo:     freelookRepo,
	}
}

// RegisterRoutes registers all routes for the free look handler
func (h *FreeLookHandler) RegisterRoutes() []serverRoute.Route {
	return []serverRoute.Route{
		// Policy Bond Tracking Endpoints
		serverRoute.NewRoute().
			SetPath("/policy-bond/track").
			SetMethod("POST").
			SetAuthRequired(true).
			SetHandlerFunc(h.TrackPolicyBond),

		serverRoute.NewRoute().
			SetPath("/policy-bond/{bond_id}/delivery-status").
			SetMethod("POST").
			SetAuthRequired(true).
			SetHandlerFunc(h.UpdateBondDelivery),

		serverRoute.NewRoute().
			SetPath("/policy-bond/{bond_id}/details").
			SetMethod("GET").
			SetAuthRequired(true).
			SetHandlerFunc(h.GetBondDetails),

		serverRoute.NewRoute().
			SetPath("/policy-bond/policy/{policy_id}").
			SetMethod("GET").
			SetAuthRequired(true).
			SetHandlerFunc(h.GetBondsByPolicy),

		// Free Look Cancellation Endpoints
		serverRoute.NewRoute().
			SetPath("/freelook/policy/{policy_id}/eligibility").
			SetMethod("GET").
			SetAuthRequired(true).
			SetHandlerFunc(h.CheckFreeLookEligibility),

		serverRoute.NewRoute().
			SetPath("/freelook/cancellation/submit").
			SetMethod("POST").
			SetAuthRequired(true).
			SetHandlerFunc(h.SubmitFreeLookCancellation),

		serverRoute.NewRoute().
			SetPath("/freelook/cancellation/{cancellation_id}/details").
			SetMethod("GET").
			SetAuthRequired(true).
			SetHandlerFunc(h.GetCancellationDetails),

		serverRoute.NewRoute().
			SetPath("/freelook/cancellation/{cancellation_id}/review").
			SetMethod("POST").
			SetAuthRequired(true).
			SetHandlerFunc(h.ReviewFreeLookCancellation),
	}
}

// ========================================
// POLICY BOND TRACKING ENDPOINTS
// ========================================

// TrackPolicyBond tracks policy bond delivery
// POST /policy-bond/track
// Reference: FR-CLM-BOND-001, BR-CLM-BOND-001
func (h *FreeLookHandler) TrackPolicyBond(sctx *serverRoute.Context, req TrackPolicyBondRequest) (*response.PolicyBondTrackedResponse, error) {
	ctx := sctx.Ctx

	// Parse dispatch date
	dispatchDate, err := time.Parse("2006-01-02", req.DispatchDate)
	if err != nil {
		log.Error(ctx, "Invalid dispatch date format: %v", err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "INVALID_DATE", "Invalid dispatch date format. Use YYYY-MM-DD")
	}

	// Calculate estimated delivery date (7 days from dispatch)
	estimatedDeliveryDate := dispatchDate.AddDate(0, 0, 7)

	// Create domain model
	bond := domain.PolicyBondTracking{
		PolicyID:               req.PolicyID,
		BondType:               req.BondType,
		DispatchDate:           &dispatchDate,
		TrackingNumber:         req.DispatchNumber,
		DeliveryStatus:         strPtr("PENDING"),
		DeliveryAttemptCount:   0,
		EscalationTriggered:    false,
	}

	// For electronic bonds, free look starts from issuance date (dispatch date)
	// For physical bonds, free look starts from delivery date
	if req.BondType == "ELECTRONIC" {
		freeLookEndDate := dispatchDate.AddDate(0, 0, 30) // 30 days for electronic
		bond.FreeLookPeriodStartDate = &dispatchDate
		bond.FreeLookPeriodEndDate = &freeLookEndDate
	}

	// Generate tracking number
	trackingNumber := "BOND" + time.Now().Format("20060102") + generateRandomNumber(6)

	// Create bond tracking record
	createdBond, err := h.policyBondRepo.Create(ctx, bond)
	if err != nil {
		log.Error(ctx, "Failed to create policy bond tracking: %v", err)
		return nil, err
	}

	// Return response
	resp := &response.PolicyBondTrackedResponse{
		StatusCodeAndMessage: response.StatusOK("Policy bond tracking created successfully"),
		ID:                   createdBond.ID,
		TrackingNumber:       trackingNumber,
		PolicyID:             createdBond.PolicyID,
		BondType:             createdBond.BondType,
		DispatchDate:         formatTimePtr(createdBond.DispatchDate),
		EstimatedDeliveryDate: formatTime(estimatedDeliveryDate),
		CreatedAt:            formatTime(createdBond.CreatedAt),
	}

	log.Info(ctx, "Policy bond tracking created: id=%s, policy_id=%s", createdBond.ID, createdBond.PolicyID)

	return resp, nil
}

// UpdateBondDelivery updates bond delivery status
// POST /policy-bond/{bond_id}/delivery-status
// Reference: FR-CLM-BOND-002, BR-CLM-BOND-002
func (h *FreeLookHandler) UpdateBondDelivery(sctx *serverRoute.Context, req UpdateBondDeliveryRequest) (*response.BondDeliveryUpdatedResponse, error) {
	ctx := sctx.Ctx

	// Parse delivery date
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		log.Error(ctx, "Invalid delivery date format: %v", err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "INVALID_DATE", "Invalid delivery date format. Use YYYY-MM-DD")
	}

	// Get existing bond
	bond, err := h.policyBondRepo.FindByID(ctx, req.BondID)
	if err != nil {
		log.Error(ctx, "Bond not found: bond_id=%s, error=%v", req.BondID, err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 404, "BOND_NOT_FOUND", "Policy bond not found")
	}

	// Calculate free look period
	var freeLookStartDate, freeLookEndDate time.Time
	var daysRemaining int

	if bond.BondType == "PHYSICAL" {
		// 15 days for physical bonds
		freeLookStartDate = deliveryDate
		freeLookEndDate = deliveryDate.AddDate(0, 0, 15)
	} else {
		// 30 days for electronic bonds
		freeLookStartDate = deliveryDate
		freeLookEndDate = deliveryDate.AddDate(0, 0, 30)
	}

	// Calculate days remaining
	daysRemaining = int(time.Until(freeLookEndDate).Hours() / 24)

	// Determine free look status
	freeLookStatus := "ACTIVE"
	if daysRemaining < 0 {
		freeLookStatus = "EXPIRED"
		daysRemaining = 0
	}

	// Update bond
	updates := map[string]interface{}{
		"delivery_status":             req.DeliveryStatus,
		"delivery_date":               deliveryDate,
		"freelook_period_start_date":  freeLookStartDate,
		"freelook_period_end_date":    freeLookEndDate,
	}

	// Increment delivery attempt count if delivered
	if req.DeliveryStatus == "DELIVERED" {
		updates["delivery_attempt_count"] = bond.DeliveryAttemptCount + 1
	}

	updatedBond, err := h.policyBondRepo.Update(ctx, req.BondID, updates)
	if err != nil {
		log.Error(ctx, "Failed to update bond delivery: bond_id=%s, error=%v", req.BondID, err)
		return nil, err
	}

	// Return response
	resp := &response.BondDeliveryUpdatedResponse{
		StatusCodeAndMessage: response.StatusOK("Bond delivery status updated successfully"),
		ID:                   updatedBond.ID,
		DeliveryStatus:       *updatedBond.DeliveryStatus,
		FreeLookStartDate:    formatTimePtr(updatedBond.FreeLookPeriodStartDate),
		FreeLookEndDate:      formatTimePtr(updatedBond.FreeLookPeriodEndDate),
		DaysRemaining:        daysRemaining,
		UpdatedAt:            formatTime(updatedBond.UpdatedAt),
	}

	if updatedBond.DeliveryDate != nil {
		formatted := formatTime(*updatedBond.DeliveryDate)
		resp.DeliveryDate = &formatted
	}

	log.Info(ctx, "Bond delivery updated: id=%s, status=%s, days_remaining=%d",
		updatedBond.ID, *updatedBond.DeliveryStatus, daysRemaining)

	return resp, nil
}

// GetBondDetails retrieves bond details
// GET /policy-bond/{bond_id}/details
func (h *FreeLookHandler) GetBondDetails(sctx *serverRoute.Context, req BondIDUri) (*response.PolicyBondDetailsResponse, error) {
	ctx := sctx.Ctx

	bond, err := h.policyBondRepo.FindByID(ctx, req.BondID)
	if err != nil {
		log.Error(ctx, "Bond not found: bond_id=%s, error=%v", req.BondID, err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 404, "BOND_NOT_FOUND", "Policy bond not found")
	}

	resp := response.NewPolicyBondDetailsResponse(bond)
	resp.StatusCodeAndMessage = response.StatusOK("Bond details retrieved successfully")

	return &resp, nil
}

// GetBondsByPolicy retrieves bonds by policy ID
// GET /policy-bond/policy/{policy_id}
func (h *FreeLookHandler) GetBondsByPolicy(sctx *serverRoute.Context, req PolicyIDUri) (*response.PolicyBondsListResponse, error) {
	ctx := sctx.Ctx

	bonds, err := h.policyBondRepo.FindByPolicyID(ctx, req.PolicyID)
	if err != nil {
		log.Error(ctx, "Failed to fetch bonds: policy_id=%s, error=%v", req.PolicyID, err)
		return nil, err
	}

	// Convert to response
	bondResponses := make([]response.PolicyBondDetailsResponse, len(bonds))
	for i, bond := range bonds {
		bondResponses[i] = response.NewPolicyBondDetailsResponse(bond)
	}

	resp := &response.PolicyBondsListResponse{
		StatusCodeAndMessage: response.StatusOK("Policy bonds retrieved successfully"),
		MetaDataResponse: response.NewMetaDataResponse(len(bondResponses), 0, len(bondResponses), ""),
		Bonds:              bondResponses,
		TotalBonds:         len(bondResponses),
	}

	return resp, nil
}

// ========================================
// FREE LOOK CANCELLATION ENDPOINTS
// ========================================

// CheckFreeLookEligibility checks if policy is eligible for free look cancellation
// GET /freelook/policy/{policy_id}/eligibility
// Reference: BR-CLM-BOND-001, VR-CLM-FL-001
func (h *FreeLookHandler) CheckFreeLookEligibility(sctx *serverRoute.Context, req PolicyIDUri) (*response.FreeLookEligibilityResponse, error) {
	ctx := sctx.Ctx

	// Get bond tracking for policy
	bonds, err := h.policyBondRepo.FindByPolicyID(ctx, req.PolicyID)
	if err != nil || len(bonds) == 0 {
		return &response.FreeLookEligibilityResponse{
			StatusCodeAndMessage: response.StatusOK("Eligibility checked successfully"),
			PolicyID:             req.PolicyID,
			Eligible:             false,
			FreeLookStatus:       "NOT_STARTED",
			Reason:               strPtr("Policy bond not yet dispatched"),
		}, nil
	}

	// Get the most recent bond
	bond := bonds[0]

	// Check eligibility based on free look status
	eligible := false
	freeLookStatus := "NOT_STARTED"
	var daysRemaining *int
	var reason *string

	if bond.FreeLookPeriodEndDate != nil {
		// Calculate days remaining
		days := int(time.Until(*bond.FreeLookPeriodEndDate).Hours() / 24)
		daysRemaining = &days

		if days > 0 {
			freeLookStatus = "ACTIVE"
			eligible = true
		} else {
			freeLookStatus = "EXPIRED"
			reason = strPtr("Free look period has expired")
		}
	}

	resp := &response.FreeLookEligibilityResponse{
		StatusCodeAndMessage: response.StatusOK("Eligibility checked successfully"),
		PolicyID:             req.PolicyID,
		Eligible:             eligible,
		FreeLookStatus:       freeLookStatus,
		BondType:             bond.BondType,
		DaysRemaining:        daysRemaining,
		Reason:               reason,
	}

	if bond.FreeLookPeriodStartDate != nil {
		formatted := formatTime(*bond.FreeLookPeriodStartDate)
		resp.FreeLookStartDate = &formatted
	}

	if bond.FreeLookPeriodEndDate != nil {
		formatted := formatTime(*bond.FreeLookPeriodEndDate)
		resp.FreeLookEndDate = &formatted
	}

	return resp, nil
}

// SubmitFreeLookCancellation submits free look cancellation request
// POST /freelook/cancellation/submit
// Reference: FR-CLM-FL-002, BR-CLM-BOND-001, VR-CLM-FL-001
func (h *FreeLookHandler) SubmitFreeLookCancellation(sctx *serverRoute.Context, req SubmitFreeLookCancellationRequest) (*response.FreeLookCancellationSubmittedResponse, error) {
	ctx := sctx.Ctx

	// Parse cancellation date
	cancellationDate, err := time.Parse("2006-01-02", req.CancellationDate)
	if err != nil {
		log.Error(ctx, "Invalid cancellation date format: %v", err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "INVALID_DATE", "Invalid cancellation date format. Use YYYY-MM-DD")
	}

	// Get bond tracking to validate free look period
	bonds, err := h.policyBondRepo.FindByPolicyID(ctx, req.PolicyID)
	if err != nil || len(bonds) == 0 {
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "BOND_NOT_FOUND", "Policy bond not found")
	}

	bond := bonds[0]

	// Validate free look period (VR-CLM-FL-001)
	if bond.FreeLookPeriodEndDate != nil && cancellationDate.After(*bond.FreeLookPeriodEndDate) {
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "FREELOOK_EXPIRED", "Free look cancellation window has expired")
	}

	// Validate bond submission (VR-CLM-FL-002)
	if !req.BondSubmitted {
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "BOND_NOT_SUBMITTED", "Original policy bond must be submitted")
	}

	// Calculate refund amount (BR-CLM-BOND-003)
	// TODO: Integrate with Policy Service to get premium details
	premiumPaid := 0.0
	if req.RefundAmount != nil {
		premiumPaid = *req.RefundAmount
	}

	// Deductions: risk premium (10%) + stamp duty (0.1%) + medical charges (5%) + other charges (1%)
	proportionateRiskPremium := premiumPaid * 0.10
	stampDuty := premiumPaid * 0.001
	medicalExamCharges := premiumPaid * 0.05
	otherCharges := premiumPaid * 0.01
	totalDeductions := proportionateRiskPremium + stampDuty + medicalExamCharges + otherCharges
	refundAmount := premiumPaid - totalDeductions

	// Generate cancellation number
	cancellationNumber := "FLC" + time.Now().Format("20060102") + generateRandomNumber(6)

	// Create domain model
	cancellation := domain.FreeLookCancellation{
		CancellationNumber:      cancellationNumber,
		PolicyID:                req.PolicyID,
		BondID:                  &bond.ID,
		CancellationReason:      req.CancellationReason,
		Channel:                 req.Channel,
		CancellationDate:        cancellationDate,
		Status:                  "SUBMITTED",
		PremiumPaid:             &premiumPaid,
		ProportionateRiskPremium: &proportionateRiskPremium,
		StampDuty:               &stampDuty,
		MedicalExamCharges:      &medicalExamCharges,
		OtherCharges:            &otherCharges,
		RefundAmount:            &refundAmount,
		RefundStatus:            strPtr("PENDING"),
		ClaimantName:            req.ClaimantName,
		ClaimantPhone:           req.ClaimantPhone,
		ClaimantEmail:           req.ClaimantEmail,
		BankAccountNumber:       req.BankAccountNumber,
		BankIFSCCode:            req.BankIFSCCode,
		DocumentURLs:            req.DocumentURLs,
	}

	// TODO: Get maker ID from user context
	makerID := "MAKER_USER_001" // Placeholder
	cancellation.MakerID = &makerID

	// Create cancellation
	createdCancellation, err := h.freelookRepo.Create(ctx, cancellation)
	if err != nil {
		log.Error(ctx, "Failed to create free look cancellation: %v", err)
		return nil, err
	}

	// Return response
	resp := &response.FreeLookCancellationSubmittedResponse{
		StatusCodeAndMessage: response.StatusCreated("Free look cancellation submitted successfully"),
		CancellationID:       createdCancellation.CancellationID,
		PolicyID:             createdCancellation.PolicyID,
		CancellationNumber:   createdCancellation.CancellationNumber,
		CancellationDate:     formatTime(createdCancellation.CancellationDate),
		Status:               createdCancellation.Status,
		RefundAmount:         createdCancellation.RefundAmount,
		SubmittedAt:          formatTime(createdCancellation.CreatedAt),
	}

	log.Info(ctx, "Free look cancellation submitted: cancellation_id=%s, policy_id=%s, amount=%f",
		createdCancellation.CancellationID, createdCancellation.PolicyID, refundAmount)

	return resp, nil
}

// GetCancellationDetails retrieves cancellation details
// GET /freelook/cancellation/{cancellation_id}/details
func (h *FreeLookHandler) GetCancellationDetails(sctx *serverRoute.Context, req CancellationIDUri) (*response.FreeLookCancellationDetailsResponse, error) {
	ctx := sctx.Ctx

	cancellation, err := h.freelookRepo.FindByID(ctx, req.CancellationID)
	if err != nil {
		log.Error(ctx, "Cancellation not found: cancellation_id=%s, error=%v", req.CancellationID, err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 404, "CANCELLATION_NOT_FOUND", "Free look cancellation not found")
	}

	resp := response.NewFreeLookCancellationDetailsResponse(cancellation)
	resp.StatusCodeAndMessage = response.StatusOK("Cancellation details retrieved successfully")

	return &resp, nil
}

// ReviewFreeLookCancellation reviews free look cancellation (maker-checker workflow)
// POST /freelook/cancellation/{cancellation_id}/review
// Reference: BR-CLM-BOND-004 (Maker-Checker Workflow)
func (h *FreeLookHandler) ReviewFreeLookCancellation(sctx *serverRoute.Context, req ReviewFreeLookCancellationRequest) (*response.FreeLookCancellationReviewResponse, error) {
	ctx := sctx.Ctx

	// Get existing cancellation
	cancellation, err := h.freelookRepo.FindByID(ctx, req.CancellationID)
	if err != nil {
		log.Error(ctx, "Cancellation not found: cancellation_id=%s, error=%v", req.CancellationID, err)
		return nil, fmt.Errorf("HTTP %d: %s - %s", 404, "CANCELLATION_NOT_FOUND", "Free look cancellation not found")
	}

	// Validate maker-checker workflow (BR-CLM-BOND-004)
	if cancellation.MakerID != nil && *cancellation.MakerID == req.CheckedBy {
		return nil, fmt.Errorf("HTTP %d: %s - %s", 400, "MAKER_CHECKER_VALIDATION", "Maker and checker cannot be the same person")
	}

	// Update cancellation based on review action
	var newStatus string
	var refundAmount *float64

	if req.ReviewAction == "APPROVE" {
		newStatus = "APPROVED"
		refundAmount = cancellation.RefundAmount

		// Apply override if provided
		if req.OverrideAmount != nil {
			refundAmount = req.OverrideAmount
		}
	} else {
		newStatus = "REJECTED"
		refundAmount = nil
	}

	updates := map[string]interface{}{
		"status":           newStatus,
		"checker_id":       req.CheckedBy,
		"review_comments":  req.ReviewComments,
		"reviewed_at":      time.Now(),
	}

	if req.OverrideAmount != nil {
		updates["override_amount"] = req.OverrideAmount
	}

	if req.OverrideReason != nil {
		updates["override_reason"] = req.OverrideReason
	}

	if req.OverrideAmount != nil {
		updates["refund_amount"] = req.OverrideAmount
	}

	_, err = h.freelookRepo.Update(ctx, req.CancellationID, updates)
	if err != nil {
		log.Error(ctx, "Failed to update cancellation: cancellation_id=%s, error=%v", req.CancellationID, err)
		return nil, err
	}

	// Return response
	resp := &response.FreeLookCancellationReviewResponse{
		StatusCodeAndMessage: response.StatusOK("Cancellation reviewed successfully"),
		CancellationID:       req.CancellationID,
		Status:               newStatus,
		RefundAmount:         refundAmount,
		ReviewedAt:           formatTime(time.Now()),
		ReviewedBy:           req.CheckedBy,
		ReviewComments:       req.ReviewComments,
	}

	log.Info(ctx, "Free look cancellation reviewed: cancellation_id=%s, action=%s, checked_by=%s",
		req.CancellationID, req.ReviewAction, req.CheckedBy)

	return resp, nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// formatTime formats time to "YYYY-MM-DD HH:MM:SS" string
func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// formatTimePtr formats time pointer to "YYYY-MM-DD HH:MM:SS" string
func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// intPtr returns a pointer to int
func intPtr(i int) *int {
	return &i
}

// strPtr returns a pointer to string
func strPtr(s string) *string {
	return &s
}

// generateRandomNumber generates a random number string with specified length
func generateRandomNumber(length int) string {
	// Simple implementation - in production, use crypto/rand
	result := ""
	for i := 0; i < length; i++ {
		result += string(rune('0' + i%10))
	}
	return result
}
