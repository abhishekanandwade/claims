package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ========================================
// POLICY SERVICE RESPONSE DTOS
// ========================================

// ExtendedPolicyDetailsResponse represents complete policy details with additional fields
// GET /policies/{id}
// Reference: INT-CLM-004 (Policy Service Integration)
type ExtendedPolicyDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string  `json:"policy_id"`
	PolicyNumber              string  `json:"policy_number"`
	PolicyType                string  `json:"policy_type"`
	SumAssured                float64 `json:"sum_assured"`
	IssueDate                 string  `json:"issue_date"` // YYYY-MM-DD format
	MaturityDate              string  `json:"maturity_date"` // YYYY-MM-DD format
	PolicyStatus              string  `json:"policy_status"`
	PolicyHolderName          string  `json:"policy_holder_name"`
	PolicyHolderDOB           string  `json:"policy_holder_dob"` // YYYY-MM-DD format
	NomineeName               *string `json:"nominee_name,omitempty"`
	NomineeRelation           *string `json:"nominee_relation,omitempty"`
	AgentCode                 *string `json:"agent_code,omitempty"`
	BranchCode                *string `json:"branch_code,omitempty"`
}

// PolicyEligibilityResponse represents policy claim eligibility check result
// GET /policies/{policy_id}/claim-eligibility
// Reference: INT-CLM-004 (Policy Service Integration)
type PolicyEligibilityResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Eligible                  bool     `json:"eligible"`
	PolicyStatus              string   `json:"policy_status"`
	EligibilityChecks         map[string]interface{} `json:"eligibility_checks"`
	IneligibilityReasons      []string `json:"ineligibility_reasons,omitempty"`
	ClaimType                 string   `json:"claim_type"`
	PolicyID                  string   `json:"policy_id"`
}

// PolicyBenefitCalculationResponse represents benefit calculation inputs from policy
// GET /policies/{id}/benefit-calculation
// Reference: INT-CLM-004 (Policy Service Integration)
type PolicyBenefitCalculationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string   `json:"policy_id"`
	SumAssured                float64  `json:"sum_assured"`
	ReversionaryBonus         float64  `json:"reversionary_bonus"`
	TerminalBonus             *float64 `json:"terminal_bonus,omitempty"`
	OutstandingLoan           float64  `json:"outstanding_loan"`
	UnpaidPremiums            float64  `json:"unpaid_premiums"`
	AccruedBonuses            float64  `json:"accrued_bonuses"`
	TotalDeductions           float64  `json:"total_deductions"` // Outstanding loan + unpaid premiums
	NetBenefit                float64  `json:"net_benefit"` // Sum assured + bonuses - deductions
}

// MaturityAmountResponse represents maturity amount calculation
// GET /policies/{policy_id}/maturity-amount
// Reference: BR-CLM-MC-001 (Maturity claim processing)
type MaturityAmountResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string                        `json:"policy_id"`
	MaturityAmount            float64                       `json:"maturity_amount"`
	MaturityDate              string                        `json:"maturity_date"` // YYYY-MM-DD format
	Breakdown                 MaturityCalculationBreakdown   `json:"breakdown"`
}

// MaturityCalculationBreakdown represents maturity amount breakdown
type MaturityCalculationBreakdown struct {
	SumAssured        float64 `json:"sum_assured"`
	Bonuses           float64 `json:"bonuses"`
	ReversionaryBonus float64 `json:"reversionary_bonus"`
	TerminalBonus     float64 `json:"terminal_bonus"`
	Interest          float64 `json:"interest"`
	TotalDeductions   float64 `json:"total_deductions"`
	NetAmount         float64 `json:"net_amount"`
}

// AccruedBonusesResponse represents accrued bonuses for policy
// GET /bonuses/{policy_id}/accrued
// Reference: INT-CLM-004 (Policy Service Integration)
type AccruedBonusesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string  `json:"policy_id"`
	ReversionaryBonus         float64 `json:"reversionary_bonus"`
	TerminalBonus             float64 `json:"terminal_bonus"`
	TotalBonus                float64 `json:"total_bonus"`
	BonusRate                 *float64 `json:"bonus_rate,omitempty"` // Bonus rate per ₹1000 SA
	BonusYear                 *int     `json:"bonus_year,omitempty"` // Bonus year
}

// OutstandingLoanResponse represents outstanding loan details
// GET /loans/{policy_id}/outstanding
// Reference: INT-CLM-004 (Policy Service Integration)
type OutstandingLoanResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string  `json:"policy_id"`
	OutstandingPrincipal      float64 `json:"outstanding_principal"`
	AccruedInterest           float64 `json:"accrued_interest"`
	TotalOutstanding          float64 `json:"total_outstanding"`
	LoanInterestRate          *float64 `json:"loan_interest_rate,omitempty"` // Interest rate %
	LoanTakenDate             *string `json:"loan_taken_date,omitempty"` // YYYY-MM-DD format
	LoanDueDate               *string `json:"loan_due_date,omitempty"` // YYYY-MM-DD format
}

// UnpaidPremiumsResponse represents unpaid premium details
// GET /premiums/{policy_id}/unpaid
// Reference: INT-CLM-004 (Policy Service Integration)
type UnpaidPremiumsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string              `json:"policy_id"`
	TotalUnpaid               float64             `json:"total_unpaid"`
	PremiumsDue               []PremiumDueItem    `json:"premiums_due"`
	OutstandingPremiumCount   int                 `json:"outstanding_premium_count"`
	LastPremiumDueDate        *string             `json:"last_premium_due_date,omitempty"` // YYYY-MM-DD format
}

// PremiumDueItem represents a single unpaid premium
type PremiumDueItem struct {
	DueDate      string  `json:"due_date"` // YYYY-MM-DD format
	Amount       float64 `json:"amount"`
	PremiumType  string  `json:"premium_type"` // REGULAR, REVIVAL, etc.
	OverdueDays  *int    `json:"overdue_days,omitempty"`
}

// FreeLookRefundCalculationExtendedResponse represents free look refund calculation with extended fields
// POST /policies/{policy_id}/freelook-refund-calculation
// Reference: BR-CLM-BOND-003 (Refund calculation)
type FreeLookRefundCalculationExtendedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                  string                   `json:"policy_id"`
	PremiumPaid               float64                  `json:"premium_paid"`
	Deductions                RefundDeductions         `json:"deductions"`
	RefundAmount              float64                  `json:"refund_amount"`
	RefundPercentage          float64                  `json:"refund_percentage"` // % of premium refunded
	CancellationDate          string                   `json:"cancellation_date"` // YYYY-MM-DD format
	FreeLookDaysUsed          int                      `json:"free_look_days_used"`
	FreeLookDaysRemaining     int                      `json:"free_look_days_remaining"`
}

// RefundDeductions represents deductions from free look refund
// Reference: BR-CLM-BOND-003 (Refund = Premium - Risk Premium - Stamp Duty - Medical - Other)
type RefundDeductions struct {
	StampDuty               float64 `json:"stamp_duty"` // 0.1% of premium
	MedicalExamCharges      float64 `json:"medical_exam_charges"` // 5% of premium (if medical done)
	ProportionateRiskPremium float64 `json:"proportionate_risk_premium"` // 10% of premium (approx)
	OtherCharges            float64 `json:"other_charges"` // 1% of premium (admin charges)
	Total                   float64 `json:"total"` // Sum of all deductions
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// NewExtendedPolicyDetailsResponse creates a new extended policy details response
func NewExtendedPolicyDetailsResponse(policyID, policyNumber, policyType string, sumAssured float64,
	issueDate, maturityDate, policyStatus, policyHolderName, policyHolderDOB string,
	nomineeName, nomineeRelation, agentCode, branchCode *string) *ExtendedPolicyDetailsResponse {

	return &ExtendedPolicyDetailsResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Policy details retrieved successfully",
		},
		PolicyID:         policyID,
		PolicyNumber:     policyNumber,
		PolicyType:       policyType,
		SumAssured:       sumAssured,
		IssueDate:        issueDate,
		MaturityDate:     maturityDate,
		PolicyStatus:     policyStatus,
		PolicyHolderName: policyHolderName,
		PolicyHolderDOB:  policyHolderDOB,
		NomineeName:      nomineeName,
		NomineeRelation:  nomineeRelation,
		AgentCode:        agentCode,
		BranchCode:       branchCode,
	}
}

// NewPolicyEligibilityResponse creates a new policy eligibility response
func NewPolicyEligibilityResponse(policyID, claimType, policyStatus string, eligible bool,
	eligibilityChecks map[string]interface{}, ineligibilityReasons []string) *PolicyEligibilityResponse {

	message := "Policy is eligible for claim"
	if !eligible {
		message = "Policy is not eligible for claim"
	}

	return &PolicyEligibilityResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    message,
		},
		PolicyID:             policyID,
		ClaimType:            claimType,
		Eligible:             eligible,
		PolicyStatus:         policyStatus,
		EligibilityChecks:    eligibilityChecks,
		IneligibilityReasons: ineligibilityReasons,
	}
}

// NewPolicyBenefitCalculationResponse creates a new benefit calculation response
func NewPolicyBenefitCalculationResponse(policyID string, sumAssured, reversionaryBonus, terminalBonus,
	outstandingLoan, unpaidPremiums, accruedBonuses float64) *PolicyBenefitCalculationResponse {

	totalDeductions := outstandingLoan + unpaidPremiums
	netBenefit := sumAssured + reversionaryBonus + terminalBonus - totalDeductions

	return &PolicyBenefitCalculationResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Benefit calculation retrieved successfully",
		},
		PolicyID:          policyID,
		SumAssured:        sumAssured,
		ReversionaryBonus: reversionaryBonus,
		TerminalBonus:     &terminalBonus,
		OutstandingLoan:   outstandingLoan,
		UnpaidPremiums:    unpaidPremiums,
		AccruedBonuses:    accruedBonuses,
		TotalDeductions:   totalDeductions,
		NetBenefit:        netBenefit,
	}
}

// NewMaturityAmountResponse creates a new maturity amount response
func NewMaturityAmountResponse(policyID string, maturityAmount float64, maturityDate string,
	breakdown MaturityCalculationBreakdown) *MaturityAmountResponse {

	return &MaturityAmountResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Maturity amount calculated successfully",
		},
		PolicyID:       policyID,
		MaturityAmount: maturityAmount,
		MaturityDate:   maturityDate,
		Breakdown:      breakdown,
	}
}

// NewAccruedBonusesResponse creates a new accrued bonuses response
func NewAccruedBonusesResponse(policyID string, reversionaryBonus, terminalBonus, totalBonus float64,
	bonusRate *float64, bonusYear *int) *AccruedBonusesResponse {

	return &AccruedBonusesResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Accrued bonuses retrieved successfully",
		},
		PolicyID:          policyID,
		ReversionaryBonus: reversionaryBonus,
		TerminalBonus:     terminalBonus,
		TotalBonus:        totalBonus,
		BonusRate:         bonusRate,
		BonusYear:         bonusYear,
	}
}

// NewOutstandingLoanResponse creates a new outstanding loan response
func NewOutstandingLoanResponse(policyID string, outstandingPrincipal, accruedInterest, totalOutstanding float64,
	loanInterestRate *float64, loanTakenDate, loanDueDate *string) *OutstandingLoanResponse {

	return &OutstandingLoanResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Outstanding loan retrieved successfully",
		},
		PolicyID:             policyID,
		OutstandingPrincipal: outstandingPrincipal,
		AccruedInterest:      accruedInterest,
		TotalOutstanding:     totalOutstanding,
		LoanInterestRate:     loanInterestRate,
		LoanTakenDate:        loanTakenDate,
		LoanDueDate:          loanDueDate,
	}
}

// NewUnpaidPremiumsResponse creates a new unpaid premiums response
func NewUnpaidPremiumsResponse(policyID string, totalUnpaid float64, premiumsDue []PremiumDueItem,
	outstandingPremiumCount int, lastPremiumDueDate *string) *UnpaidPremiumsResponse {

	return &UnpaidPremiumsResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Unpaid premiums retrieved successfully",
		},
		PolicyID:                policyID,
		TotalUnpaid:             totalUnpaid,
		PremiumsDue:             premiumsDue,
		OutstandingPremiumCount: outstandingPremiumCount,
		LastPremiumDueDate:      lastPremiumDueDate,
	}
}

// NewFreeLookRefundCalculationExtendedResponse creates a new free look refund calculation response
func NewFreeLookRefundCalculationExtendedResponse(policyID string, premiumPaid, refundAmount, refundPercentage float64,
	deductions RefundDeductions, cancellationDate string, freeLookDaysUsed, freeLookDaysRemaining int) *FreeLookRefundCalculationExtendedResponse {

	return &FreeLookRefundCalculationExtendedResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Free look refund calculated successfully",
		},
		PolicyID:              policyID,
		PremiumPaid:           premiumPaid,
		Deductions:            deductions,
		RefundAmount:          refundAmount,
		RefundPercentage:      refundPercentage,
		CancellationDate:      cancellationDate,
		FreeLookDaysUsed:      freeLookDaysUsed,
		FreeLookDaysRemaining: freeLookDaysRemaining,
	}
}
