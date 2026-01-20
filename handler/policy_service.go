package handler

import (
	"fmt"
	"math"
	"time"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/handler/response"
)

// PolicyServiceHandler handles policy service-related HTTP requests
type PolicyServiceHandler struct {
	*serverHandler.Base
	// TODO: Add Policy Service client when integration is implemented
	// policyServiceClient *PolicyServiceClient
}

// NewPolicyServiceHandler creates a new policy service handler
func NewPolicyServiceHandler() *PolicyServiceHandler {
	base := serverHandler.New("PolicyService").
		SetPrefix("/v1").
		AddPrefix("")
	return &PolicyServiceHandler{
		Base: base,
	}
}

// RegisterRoutes registers all policy service routes
func (h *PolicyServiceHandler) RegisterRoutes() {
	// Policy Details and Eligibility
	h.Add(serverRoute.New("/policies/:id", h.GetPolicyDetails, "GET"))
	h.Add(serverRoute.New("/policies/:policy_id/claim-eligibility", h.CheckPolicyClaimEligibility, "GET"))

	// Benefit Calculation
	h.Add(serverRoute.New("/policies/:id/benefit-calculation", h.GetPolicyBenefitCalculation, "GET"))
	h.Add(serverRoute.New("/policies/:policy_id/maturity-amount", h.GetMaturityAmount, "GET"))

	// Financial Components
	h.Add(serverRoute.New("/bonuses/:policy_id/accrued", h.GetAccruedBonuses, "GET"))
	h.Add(serverRoute.New("/loans/:policy_id/outstanding", h.GetOutstandingLoan, "GET"))
	h.Add(serverRoute.New("/premiums/:policy_id/unpaid", h.GetUnpaidPremiums, "GET"))

	// Free Look Refund Calculation
	h.Add(serverRoute.New("/policies/:policy_id/freelook-refund-calculation", h.CalculateFreeLookRefund, "POST"))
}

// PolicyIDUri represents URI parameters for /policies/:id
type PolicyIDUri struct {
	ID string `uri:"id" validate:"required"`
}

// PolicyIDUriForPolicyID represents URI parameters for /policies/:policy_id
type PolicyIDUriForPolicyID struct {
	PolicyID string `uri:"policy_id" validate:"required"`
}

// ========================================
// POLICY DETAILS AND ELIGIBILITY ENDPOINTS
// ========================================

// GetPolicyDetails retrieves complete policy details
// GET /policies/:id
// Reference: INT-CLM-002 (McCamish Policy System Integration)
func (h *PolicyServiceHandler) GetPolicyDetails(sctx *serverRoute.Context, req *PolicyIDUri) (*response.ExtendedPolicyDetailsResponse, error) {
	// TODO: Integrate with Policy Service (INT-CLM-002) to fetch real policy details
	// For now, return mock data

	return response.NewExtendedPolicyDetailsResponse(
		req.ID,
		"POL123456789",
		"ENDOWMENT",
		500000.0,
		"2015-06-15",
		"2035-06-15",
		"INFORCE",
		"John Doe",
		"1985-03-20",
		stringPtr("Jane Doe"),
		stringPtr("Spouse"),
		stringPtr("AGT001"),
		stringPtr("BRCH001"),
	), nil
}

// CheckPolicyClaimEligibility checks if policy is eligible for claim
// GET /policies/:policy_id/claim-eligibility?claim_type=DEATH
// Reference: INT-CLM-002 (Policy Service Integration)
func (h *PolicyServiceHandler) CheckPolicyClaimEligibility(sctx *serverRoute.Context, req *CheckPolicyClaimEligibilityRequest) (*response.PolicyEligibilityResponse, error) {
	// TODO: Integrate with Policy Service (INT-CLM-002) for real eligibility checks
	// For now, implement basic eligibility logic

	// Mock eligibility checks
	eligibilityChecks := make(map[string]interface{})
	var ineligibilityReasons []string
	eligible := true
	policyStatus := "INFORCE"

	// Check 1: Policy status
	eligibilityChecks["policy_status"] = map[string]interface{}{
		"status":  policyStatus,
		"eligible": policyStatus == "INFORCE",
	}
	if policyStatus != "INFORCE" {
		eligible = false
		ineligibilityReasons = append(ineligibilityReasons, "Policy is not in-force status")
	}

	// Check 2: Claim type specific validations
	switch req.ClaimType {
	case "DEATH":
		eligibilityChecks["death_claim_validation"] = map[string]interface{}{
			"requires_investigation": false,
			"investigation_rule":     "BR-CLM-DC-001 (3-year rule)",
		}
	case "MATURITY":
		// Check if policy has reached maturity date
		eligibilityChecks["maturity_validation"] = map[string]interface{}{
			"maturity_date_reached": true,
			"days_to_maturity":      0,
		}
	case "SURVIVAL_BENEFIT":
		eligibilityChecks["survival_benefit_validation"] = map[string]interface{}{
			"due_installments": 2,
			"eligible":         true,
		}
	case "FREELOOK":
		eligibilityChecks["freelook_validation"] = map[string]interface{}{
			"within_freelook_period": true,
			"days_remaining":         12,
		}
	}

	// Check 3: No duplicate claims
	eligibilityChecks["duplicate_claim_check"] = map[string]interface{}{
		"has_active_claim": false,
		"eligible":         true,
	}

	return response.NewPolicyEligibilityResponse(
		req.PolicyID,
		req.ClaimType,
		policyStatus,
		eligible,
		eligibilityChecks,
		ineligibilityReasons,
	), nil
}

// ========================================
// BENEFIT CALCULATION ENDPOINTS
// ========================================

// GetPolicyBenefitCalculation retrieves benefit calculation inputs from policy
// GET /policies/:id/benefit-calculation
// Reference: INT-CLM-002 (Policy Service), INT-CLM-003 (Bonus Ledger), INT-CLM-004 (Loan Module)
func (h *PolicyServiceHandler) GetPolicyBenefitCalculation(sctx *serverRoute.Context, req *PolicyIDUri) (*response.PolicyBenefitCalculationResponse, error) {
	// TODO: Integrate with Policy Service (INT-CLM-002), Bonus Ledger (INT-CLM-003), Loan Module (INT-CLM-004)
	// For now, return mock data

	// Mock benefit calculation data
	sumAssured := 500000.0
	reversionaryBonus := 75000.0
	terminalBonus := 25000.0
	outstandingLoan := 50000.0
	unpaidPremiums := 10000.0
	accruedBonuses := reversionaryBonus + terminalBonus

	return response.NewPolicyBenefitCalculationResponse(
		req.ID,
		sumAssured,
		reversionaryBonus,
		terminalBonus,
		outstandingLoan,
		unpaidPremiums,
		accruedBonuses,
	), nil
}

// GetMaturityAmount calculates maturity amount for policy
// GET /policies/:policy_id/maturity-amount
// Reference: BR-CLM-MC-001, CALC-001 (Maturity claim calculation)
func (h *PolicyServiceHandler) GetMaturityAmount(sctx *serverRoute.Context, req *PolicyIDUri) (*response.MaturityAmountResponse, error) {
	// TODO: Integrate with Policy Service (INT-CLM-002) and Bonus Ledger (INT-CLM-003)
	// For now, implement basic maturity amount calculation

	// Mock maturity calculation
	sumAssured := 500000.0
	reversionaryBonus := 75000.0
	terminalBonus := 25000.0
	accruedInterest := 15000.0

	maturityAmount := sumAssured + reversionaryBonus + terminalBonus + accruedInterest
	maturityDate := "2035-06-15"

	// Create breakdown
	breakdown := response.MaturityCalculationBreakdown{
		SumAssured:        sumAssured,
		Bonuses:           reversionaryBonus + terminalBonus,
		ReversionaryBonus: reversionaryBonus,
		TerminalBonus:     terminalBonus,
		Interest:          accruedInterest,
		TotalDeductions:   0.0,
		NetAmount:         maturityAmount,
	}

	return response.NewMaturityAmountResponse(
		req.PolicyID,
		maturityAmount,
		maturityDate,
		breakdown,
	), nil
}

// ========================================
// FINANCIAL COMPONENTS ENDPOINTS
// ========================================

// GetAccruedBonuses retrieves accrued bonuses for policy
// GET /bonuses/{policy_id}/accrued
// Reference: INT-CLM-003 (Bonus Ledger Integration)
func (h *PolicyServiceHandler) GetAccruedBonuses(sctx *serverRoute.Context, req *PolicyIDUri) (*response.AccruedBonusesResponse, error) {
	// TODO: Integrate with Bonus Ledger (INT-CLM-003)
	// For now, return mock data

	reversionaryBonus := 75000.0
	terminalBonus := 25000.0
	totalBonus := reversionaryBonus + terminalBonus
	bonusRate := float64Ptr(50.0) // ₹50 per ₹1000 SA
	bonusYear := intPtr(2024)

	return response.NewAccruedBonusesResponse(
		req.PolicyID,
		reversionaryBonus,
		terminalBonus,
		totalBonus,
		bonusRate,
		bonusYear,
	), nil
}

// GetOutstandingLoan retrieves outstanding loan amount
// GET /loans/{policy_id}/outstanding
// Reference: INT-CLM-004 (Loan Module Integration)
func (h *PolicyServiceHandler) GetOutstandingLoan(sctx *serverRoute.Context, req *PolicyIDUri) (*response.OutstandingLoanResponse, error) {
	// TODO: Integrate with Loan Module (INT-CLM-004)
	// For now, return mock data

	outstandingPrincipal := 50000.0
	accruedInterest := 5000.0
	totalOutstanding := outstandingPrincipal + accruedInterest
	loanInterestRate := float64Ptr(10.5) // 10.5% p.a.
	loanTakenDate := stringPtr("2023-01-15")
	loanDueDate := stringPtr("2025-01-15")

	return response.NewOutstandingLoanResponse(
		req.PolicyID,
		outstandingPrincipal,
		accruedInterest,
		totalOutstanding,
		loanInterestRate,
		loanTakenDate,
		loanDueDate,
	), nil
}

// GetUnpaidPremiums retrieves unpaid premium amount
// GET /premiums/{policy_id}/unpaid
// Reference: INT-CLM-002 (Policy Service Integration)
func (h *PolicyServiceHandler) GetUnpaidPremiums(sctx *serverRoute.Context, req *PolicyIDUri) (*response.UnpaidPremiumsResponse, error) {
	// TODO: Integrate with Policy Service (INT-CLM-002)
	// For now, return mock data

	totalUnpaid := 10000.0

	// Mock premium due items
	premiumsDue := []response.PremiumDueItem{
		{
			DueDate:      "2024-12-15",
			Amount:       5000.0,
			PremiumType:  "REGULAR",
			OverdueDays:  intPtr(35),
		},
		{
			DueDate:      "2025-01-15",
			Amount:       5000.0,
			PremiumType:  "REGULAR",
			OverdueDays:  intPtr(5),
		},
	}

	outstandingPremiumCount := len(premiumsDue)
	lastPremiumDueDate := stringPtr("2025-01-15")

	return response.NewUnpaidPremiumsResponse(
		req.PolicyID,
		totalUnpaid,
		premiumsDue,
		outstandingPremiumCount,
		lastPremiumDueDate,
	), nil
}

// ========================================
// FREE LOOK REFUND CALCULATION ENDPOINT
// ========================================

// CalculateFreeLookRefund calculates free look refund amount
// POST /policies/{policy_id}/freelook-refund-calculation
// Reference: BR-CLM-BOND-003 (Refund calculation)
func (h *PolicyServiceHandler) CalculateFreeLookRefund(sctx *serverRoute.Context, req *CalculateFreeLookRefundRequest) (*response.FreeLookRefundCalculationExtendedResponse, error) {
	// Parse dates
	cancellationDate, err := time.Parse("2006-01-02", req.CancellationDate)
	if err != nil {
		return nil, fmt.Errorf("invalid cancellation_date format: %w", err)
	}

	// Determine free look period based on bond type
	// BR-CLM-BOND-001: 15 days for physical, 30 days for electronic
	var freeLookDays int
	var freeLookStartDate time.Time

	if req.BondType == "PHYSICAL" {
		// Physical bonds: 15 days from delivery date
		if req.DeliveryDate == nil {
			return nil, fmt.Errorf("delivery_date is required for PHYSICAL bonds")
		}
		deliveryDate, err := time.Parse("2006-01-02", *req.DeliveryDate)
		if err != nil {
			return nil, fmt.Errorf("invalid delivery_date format: %w", err)
		}
		freeLookStartDate = deliveryDate
		freeLookDays = 15
	} else {
		// Electronic bonds: 30 days from issue date
		// TODO: Get issue date from Policy Service
		freeLookStartDate = cancellationDate.AddDate(0, 0, -30) // Mock: assume issued 30 days ago
		freeLookDays = 30
	}

	// Calculate free look days used
	daysUsed := int(cancellationDate.Sub(freeLookStartDate).Hours()/24) + 1
	freeLookDaysUsed := daysUsed
	if freeLookDaysUsed < 0 {
		freeLookDaysUsed = 0
	}
	freeLookDaysRemaining := freeLookDays - freeLookDaysUsed

	// Check if within free look period
	if freeLookDaysUsed > freeLookDays {
		return nil, fmt.Errorf("cancellation is beyond free look period (%d days)", freeLookDays)
	}

	// Calculate deductions as per BR-CLM-BOND-003
	// Refund = Premium - (Risk Premium + Stamp Duty + Medical + Other)
	stampDuty := req.PremiumPaid * 0.001 // 0.1% for stamp duty
	medicalExamCharges := req.PremiumPaid * 0.05 // 5% for medical (if done)
	proportionateRiskPremium := req.PremiumPaid * 0.10 // 10% for risk premium
	otherCharges := req.PremiumPaid * 0.01 // 1% for admin charges

	totalDeductions := stampDuty + medicalExamCharges + proportionateRiskPremium + otherCharges
	refundAmount := req.PremiumPaid - totalDeductions
	refundPercentage := (refundAmount / req.PremiumPaid) * 100

	// Round to 2 decimal places
	refundAmount = math.Round(refundAmount*100) / 100
	refundPercentage = math.Round(refundPercentage*100) / 100

	deductions := response.RefundDeductions{
		StampDuty:               math.Round(stampDuty*100) / 100,
		MedicalExamCharges:      math.Round(medicalExamCharges*100) / 100,
		ProportionateRiskPremium: math.Round(proportionateRiskPremium*100) / 100,
		OtherCharges:            math.Round(otherCharges*100) / 100,
		Total:                   math.Round(totalDeductions*100) / 100,
	}

	return response.NewFreeLookRefundCalculationExtendedResponse(
		req.PolicyID,
		math.Round(req.PremiumPaid*100)/100,
		refundAmount,
		refundPercentage,
		deductions,
		req.CancellationDate,
		freeLookDaysUsed,
		freeLookDaysRemaining,
	), nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// stringPtr returns a pointer to string
func stringPtr(s string) *string {
	return &s
}

// float64Ptr returns a pointer to float64
func float64Ptr(f float64) *float64 {
	return &f
}

// intPtr returns a pointer to int
func intPtr(i int) *int {
	return &i
}
