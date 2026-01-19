package service

import (
	"fmt"
	"math"
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// BusinessRulesService implements business logic for death claims
// Reference: seed/analysis/business_rules.md, .zenflow/tasks/code-gen-54c7/requirements.md:668-693
type BusinessRulesService struct{}

// NewBusinessRulesService creates a new business rules service
func NewBusinessRulesService() *BusinessRulesService {
	return &BusinessRulesService{}
}

// SLA Status Constants
// Reference: BR-CLM-DC-021 (Color-coded SLA)
const (
	SLAStatusGreen   = "GREEN"   // < 70% of SLA consumed
	SLAStatusYellow  = "YELLOW"  // 70-90% of SLA consumed
	SLAStatusOrange  = "ORANGE"  // 90-100% of SLA consumed
	SLAStatusRed     = "RED"     // > 100% of SLA consumed (breached)
)

// Death Type Constants
const (
	DeathTypeNatural      = "NATURAL"
	DeathTypeAccidental   = "ACCIDENTAL"
	DeathTypeUnnatural    = "UNNATURAL"
)

// =============================================================================
// BR-CLM-DC-001: Investigation trigger (3-year rule)
// Reference: requirements.md:44, 669
// =============================================================================

// ShouldTriggerInvestigation determines if a claim requires investigation
// Business Rule: Auto-trigger investigation if death within 3 years of policy issue/revival
// Reference: BR-CLM-DC-001
func (s *BusinessRulesService) ShouldTriggerInvestigation(
	policyIssueDate time.Time,
	policyRevivalDate *time.Time,
	deathDate time.Time,
) (bool, string) {
	// Check if death occurred within 3 years of policy issue
	threeYearsAfterIssue := policyIssueDate.AddDate(3, 0, 0)
	if deathDate.Before(threeYearsAfterIssue) {
		reason := fmt.Sprintf("Death within 3 years of policy issue date (%s)",
			policyIssueDate.Format("2006-01-02"))
		return true, reason
	}

	// Check if policy was revived and death occurred within 3 years of revival
	if policyRevivalDate != nil {
		threeYearsAfterRevival := policyRevivalDate.AddDate(3, 0, 0)
		if deathDate.Before(threeYearsAfterRevival) {
			reason := fmt.Sprintf("Death within 3 years of policy revival date (%s)",
				policyRevivalDate.Format("2006-01-02"))
			return true, reason
		}
	}

	return false, ""
}

// =============================================================================
// BR-CLM-DC-002: Investigation SLA (21 days)
// Reference: requirements.md:45, 670
// =============================================================================

// CalculateInvestigationSLA calculates the investigation due date
// Business Rule: Investigation SLA is 21 days from assignment
// Reference: BR-CLM-DC-002
func (s *BusinessRulesService) CalculateInvestigationSLA(investigationStartDate time.Time) time.Time {
	return investigationStartDate.AddDate(0, 0, 21)
}

// IsInvestigationOverdue checks if investigation is past SLA
// Reference: BR-CLM-DC-002
func (s *BusinessRulesService) IsInvestigationOverdue(investigationStartDate, currentDate time.Time) bool {
	dueDate := s.CalculateInvestigationSLA(investigationStartDate)
	return currentDate.After(dueDate)
}

// GetInvestigationDaysRemaining returns days remaining for investigation SLA
// Reference: BR-CLM-DC-002
func (s *BusinessRulesService) GetInvestigationDaysRemaining(investigationStartDate, currentDate time.Time) int {
	dueDate := s.CalculateInvestigationSLA(investigationStartDate)
	remaining := int(dueDate.Sub(currentDate).Hours() / 24)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// =============================================================================
// BR-CLM-DC-003: SLA without investigation (15 days)
// Reference: requirements.md:46, 671
// =============================================================================

// CalculateClaimSLAWithoutInvestigation calculates SLA for claims without investigation
// Business Rule: SLA is 15 days from claim registration
// Reference: BR-CLM-DC-003
func (s *BusinessRulesService) CalculateClaimSLAWithoutInvestigation(claimDate time.Time) time.Time {
	return claimDate.AddDate(0, 0, 15)
}

// =============================================================================
// BR-CLM-DC-004: SLA with investigation (45 days)
// Reference: requirements.md:46, 672
// =============================================================================

// CalculateClaimSLAWithInvestigation calculates SLA for claims with investigation
// Business Rule: SLA is 45 days from claim registration
// Reference: BR-CLM-DC-004
func (s *BusinessRulesService) CalculateClaimSLAWithInvestigation(claimDate time.Time) time.Time {
	return claimDate.AddDate(0, 0, 45)
}

// =============================================================================
// Combined SLA Calculation
// =============================================================================

// CalculateClaimSLADueDate calculates the SLA due date based on investigation requirement
// Reference: BR-CLM-DC-003, BR-CLM-DC-004
func (s *BusinessRulesService) CalculateClaimSLADueDate(claimDate time.Time, investigationRequired bool) time.Time {
	if investigationRequired {
		return s.CalculateClaimSLAWithInvestigation(claimDate)
	}
	return s.CalculateClaimSLAWithoutInvestigation(claimDate)
}

// =============================================================================
// BR-CLM-DC-009: Penal interest calculation (8% p.a.)
// Reference: requirements.md:47, 677
// =============================================================================

// CalculatePenalInterest calculates penal interest for SLA breaches
// Business Rule: Penal interest @ 8% p.a. for SLA breaches
// Formula: (Claim Amount * 8% * Breach Days) / 365
// Reference: BR-CLM-DC-009
func (s *BusinessRulesService) CalculatePenalInterest(claimAmount float64, breachDays int) float64 {
	if claimAmount <= 0 || breachDays <= 0 {
		return 0.0
	}

	// Formula: Claim Amount * (8/100) * (Breach Days / 365)
	penalInterest := claimAmount * 0.08 * float64(breachDays) / 365.0

	// Round to 2 decimal places
	return math.Round(penalInterest*100) / 100
}

// CalculatePenalInterestWithDate calculates penal interest using actual dates
// Reference: BR-CLM-DC-009
func (s *BusinessRulesService) CalculatePenalInterestWithDate(
	claimAmount float64,
	slaDueDate time.Time,
	settlementDate time.Time,
) (float64, int) {
	if settlementDate.Before(slaDueDate) || settlementDate.Equal(slaDueDate) {
		return 0.0, 0
	}

	breachDays := int(settlementDate.Sub(slaDueDate).Hours() / 24)
	penalInterest := s.CalculatePenalInterest(claimAmount, breachDays)

	return penalInterest, breachDays
}

// =============================================================================
// BR-CLM-DC-021: SLA color coding (GREEN/YELLOW/ORANGE/RED)
// Reference: requirements.md:49, 688
// =============================================================================

// CalculateSLAStatus determines the SLA status based on time elapsed
// Business Rule:
//   - GREEN: < 70% of SLA consumed
//   - YELLOW: 70-90% of SLA consumed
//   - ORANGE: 90-100% of SLA consumed
//   - RED: > 100% of SLA consumed (breached)
// Reference: BR-CLM-DC-021
func (s *BusinessRulesService) CalculateSLAStatus(claimDate, slaDueDate, currentDate time.Time) string {
	if currentDate.After(slaDueDate) {
		return SLAStatusRed
	}

	totalDuration := slaDueDate.Sub(claimDate).Hours()
	elapsedDuration := currentDate.Sub(claimDate).Hours()

	if totalDuration <= 0 {
		return SLAStatusRed
	}

	percentageConsumed := (elapsedDuration / totalDuration) * 100

	switch {
	case percentageConsumed < 70:
		return SLAStatusGreen
	case percentageConsumed < 90:
		return SLAStatusYellow
	case percentageConsumed < 100:
		return SLAStatusOrange
	default:
		return SLAStatusRed
	}
}

// CalculateSLAStatusForClaim calculates SLA status for a claim domain object
// Reference: BR-CLM-DC-021
func (s *BusinessRulesService) CalculateSLAStatusForClaim(claim domain.Claim) string {
	currentDate := time.Now()
	return s.CalculateSLAStatus(claim.ClaimDate, claim.SLADueDate, currentDate)
}

// CalculateSLAPercentageRemaining calculates the percentage of SLA remaining
// Reference: BR-CLM-DC-021
func (s *BusinessRulesService) CalculateSLAPercentageRemaining(claimDate, slaDueDate, currentDate time.Time) float64 {
	if currentDate.After(slaDueDate) {
		return 0.0
	}

	totalDuration := slaDueDate.Sub(claimDate).Hours()
	elapsedDuration := currentDate.Sub(claimDate).Hours()

	if totalDuration <= 0 {
		return 0.0
	}

	percentageConsumed := (elapsedDuration / totalDuration) * 100
	percentageRemaining := 100.0 - percentageConsumed

	if percentageRemaining < 0 {
		return 0.0
	}

	return math.Round(percentageRemaining*100) / 100
}

// =============================================================================
// Claim Amount Calculation
// Reference: BR-CLM-DC-008 (Calculation formula)
// =============================================================================

// CalculateClaimAmount calculates the claim amount
// Business Rule: Sum Assured + Reversionary Bonus + Terminal Bonus - Outstanding Loan - Unpaid Premiums
// Reference: BR-CLM-DC-008
func (s *BusinessRulesService) CalculateClaimAmount(
	sumAssured float64,
	reversionaryBonus float64,
	terminalBonus float64,
	outstandingLoan float64,
	unpaidPremiums float64,
) float64 {
	total := sumAssured + reversionaryBonus + terminalBonus - outstandingLoan - unpaidPremiums

	// Ensure amount is not negative
	if total < 0 {
		return 0.0
	}

	return math.Round(total*100) / 100
}

// CalculateClaimAmountForDomain calculates claim amount from domain object
// Reference: BR-CLM-DC-008
func (s *BusinessRulesService) CalculateClaimAmountForDomain(claim domain.Claim) float64 {
	sumAssured := 0.0
	if claim.SumAssured != nil {
		sumAssured = *claim.SumAssured
	}

	reversionaryBonus := 0.0
	if claim.ReversionaryBonus != nil {
		reversionaryBonus = *claim.ReversionaryBonus
	}

	terminalBonus := 0.0
	if claim.TerminalBonus != nil {
		terminalBonus = *claim.TerminalBonus
	}

	outstandingLoan := 0.0
	if claim.OutstandingLoan != nil {
		outstandingLoan = *claim.OutstandingLoan
	}

	unpaidPremiums := 0.0
	if claim.UnpaidPremiums != nil {
		unpaidPremiums = *claim.UnpaidPremiums
	}

	return s.CalculateClaimAmount(
		sumAssured,
		reversionaryBonus,
		terminalBonus,
		outstandingLoan,
		unpaidPremiums,
	)
}

// =============================================================================
// Document Checklist Rules
// Reference: BR-CLM-DC-011, BR-CLM-DC-012, BR-CLM-DC-013, BR-CLM-DC-014, BR-CLM-DC-015
// =============================================================================

// GetRequiredDocumentTypes returns required document types based on claim characteristics
// Business Rules:
//   - BR-CLM-DC-015: Base mandatory documents (death certificate, claim form, ID proof)
//   - BR-CLM-DC-013: Conditional documents for unnatural death (post-mortem, FIR, police report)
//   - BR-CLM-DC-014: Additional documents if no nomination (succession certificate, legal heir certificate)
// Reference: BR-CLM-DC-011, BR-CLM-DC-013, BR-CLM-DC-014, BR-CLM-DC-015
func (s *BusinessRulesService) GetRequiredDocumentTypes(
	deathType string,
	hasNomination bool,
) []string {
	baseDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
		"CLAIMANT_ID_PROOF",
		"CLAIMANT_PHOTOGRAPH",
		"POLICY_BOND_OR_INDEMNITY_BOND",
		"BANK_ACCOUNT_PROOF",
	}

	// BR-CLM-DC-013: Conditional documents for unnatural death
	if deathType == DeathTypeUnnatural || deathType == DeathTypeAccidental {
		baseDocuments = append(baseDocuments,
			"POST_MORTEM_REPORT",
			"FIR_COPY",
			"POLICE_REPORT",
			"AUTOPSY_REPORT",
		)
	}

	// BR-CLM-DC-014: Additional documents if no nomination
	if !hasNomination {
		baseDocuments = append(baseDocuments,
			"SUCCESSION_CERTIFICATE",
			"LEGAL_HEIR_CERTIFICATE",
			"NO_NOMINATION_AFFIDAVIT",
		)
	}

	return baseDocuments
}

// IsDocumentComplete checks if all required documents are uploaded
// Reference: BR-CLM-DC-012
func (s *BusinessRulesService) IsDocumentComplete(
	uploadedDocuments []string,
	requiredDocuments []string,
) bool {
	uploadedMap := make(map[string]bool)
	for _, doc := range uploadedDocuments {
		uploadedMap[doc] = true
	}

	for _, required := range requiredDocuments {
		if !uploadedMap[required] {
			return false
		}
	}

	return true
}

// =============================================================================
// Approval Hierarchy
// Reference: BR-CLM-DC-022
// =============================================================================

// GetApprovalLevel determines the approval level based on claim amount
// Business Rule: Approval hierarchy based on claim amount ranges
// Reference: BR-CLM-DC-022
func (s *BusinessRulesService) GetApprovalLevel(claimAmount float64) string {
	switch {
	case claimAmount <= 100000:
		return "LEVEL_1" // Assistant Divisional Manager
	case claimAmount <= 500000:
		return "LEVEL_2" // Divisional Manager
	case claimAmount <= 1000000:
		return "LEVEL_3" // Senior Divisional Manager
	default:
		return "LEVEL_4" // Regional Office or above
	}
}

// =============================================================================
// Reinvestigation Limit
// Reference: BR-CLM-DC-023
// =============================================================================

// CanReinvestigate checks if reinvestigation is allowed
// Business Rule: Maximum 2 reinvestigations allowed
// Reference: BR-CLM-DC-023
func (s *BusinessRulesService) CanReinvestigate(currentReinvestigationCount int) bool {
	return currentReinvestigationCount < 2
}

// =============================================================================
// Payment Mode Priority
// Reference: BR-CLM-DC-017
// =============================================================================

// GetPreferredPaymentMode returns the preferred payment mode
// Business Rule: Payment mode priority - NEFT > POSB > Cheque
// Reference: BR-CLM-DC-017
func (s *BusinessRulesService) GetPreferredPaymentMode(
	hasNEFT bool,
	hasPOSB bool,
) string {
	if hasNEFT {
		return "NEFT"
	}
	if hasPOSB {
		return "POSB"
	}
	return "CHEQUE"
}

// =============================================================================
// Communication Triggers
// Reference: BR-CLM-DC-019
// =============================================================================

// ShouldSendDocumentReminder checks if document reminder should be sent
// Business Rule: Send reminder if documents pending after 7 days of registration
// Reference: BR-CLM-DC-019
func (s *BusinessRulesService) ShouldSendDocumentReminder(claimDate, currentDate time.Time, isDocumentComplete bool) bool {
	if isDocumentComplete {
		return false
	}

	daysSinceRegistration := int(currentDate.Sub(claimDate).Hours() / 24)
	return daysSinceRegistration >= 7
}

// ShouldSendSLABreachedWarning checks if SLA breach warning should be sent
// Business Rule: Send warning when SLA status is ORANGE (90% consumed)
// Reference: BR-CLM-DC-019, BR-CLM-DC-021
func (s *BusinessRulesService) ShouldSendSLABreachedWarning(slaStatus string) bool {
	return slaStatus == SLAStatusOrange
}

// =============================================================================
// Helper Functions
// =============================================================================

// IsValidDeathType checks if death type is valid
func (s *BusinessRulesService) IsValidDeathType(deathType string) bool {
	validTypes := map[string]bool{
		DeathTypeNatural:    true,
		DeathTypeAccidental: true,
		DeathTypeUnnatural:  true,
	}
	return validTypes[deathType]
}

// IsValidClaimType checks if claim type is valid
func (s *BusinessRulesService) IsValidClaimType(claimType string) bool {
	validTypes := map[string]bool{
		"DEATH":           true,
		"MATURITY":        true,
		"SURVIVAL_BENEFIT": true,
		"FREELOOK":        true,
	}
	return validTypes[claimType]
}
