package service

import (
	"math"
	"testing"
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
)

// =============================================================================
// Test BR-CLM-DC-001: Investigation trigger (3-year rule)
// =============================================================================

func TestShouldTriggerInvestigation_WithinThreeYearsOfIssue(t *testing.T) {
	service := NewBusinessRulesService()

	policyIssueDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	var policyRevivalDate *time.Time = nil
	deathDate := time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC) // Within 3 years

	shouldTrigger, reason := service.ShouldTriggerInvestigation(policyIssueDate, policyRevivalDate, deathDate)

	if !shouldTrigger {
		t.Errorf("Expected investigation to trigger within 3 years of policy issue")
	}

	if reason == "" {
		t.Errorf("Expected reason to be provided")
	}
}

func TestShouldTriggerInvestigation_WithinThreeYearsOfRevival(t *testing.T) {
	service := NewBusinessRulesService()

	policyIssueDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	revivalDate := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC) // Within 3 years of revival

	shouldTrigger, reason := service.ShouldTriggerInvestigation(policyIssueDate, &revivalDate, deathDate)

	if !shouldTrigger {
		t.Errorf("Expected investigation to trigger within 3 years of policy revival")
	}

	if reason == "" {
		t.Errorf("Expected reason to be provided")
	}
}

func TestShouldTriggerInvestigation_OutsideThreeYears(t *testing.T) {
	service := NewBusinessRulesService()

	policyIssueDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	var policyRevivalDate *time.Time = nil
	deathDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) // > 3 years

	shouldTrigger, reason := service.ShouldTriggerInvestigation(policyIssueDate, policyRevivalDate, deathDate)

	if shouldTrigger {
		t.Errorf("Expected investigation NOT to trigger after 3 years of policy issue")
	}

	if reason != "" {
		t.Errorf("Expected no reason when investigation not required, got: %s", reason)
	}
}

// =============================================================================
// Test BR-CLM-DC-002: Investigation SLA (21 days)
// =============================================================================

func TestCalculateInvestigationSLA(t *testing.T) {
	service := NewBusinessRulesService()

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedDueDate := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC) // 21 days later

	dueDate := service.CalculateInvestigationSLA(startDate)

	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, dueDate)
	}
}

func TestIsInvestigationOverdue_NotOverdue(t *testing.T) {
	service := NewBusinessRulesService()

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // 14 days later

	isOverdue := service.IsInvestigationOverdue(startDate, currentDate)

	if isOverdue {
		t.Errorf("Expected investigation not to be overdue at 14 days")
	}
}

func TestIsInvestigationOverdue_Overdue(t *testing.T) {
	service := NewBusinessRulesService()

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC) // 22 days later

	isOverdue := service.IsInvestigationOverdue(startDate, currentDate)

	if !isOverdue {
		t.Errorf("Expected investigation to be overdue at 22 days")
	}
}

func TestGetInvestigationDaysRemaining(t *testing.T) {
	service := NewBusinessRulesService()

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC) // 9 days later

	expectedDaysRemaining := 12 // 21 - 9 = 12
	daysRemaining := service.GetInvestigationDaysRemaining(startDate, currentDate)

	if daysRemaining != expectedDaysRemaining {
		t.Errorf("Expected %d days remaining, got %d", expectedDaysRemaining, daysRemaining)
	}
}

// =============================================================================
// Test BR-CLM-DC-003: SLA without investigation (15 days)
// =============================================================================

func TestCalculateClaimSLAWithoutInvestigation(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days later

	dueDate := service.CalculateClaimSLAWithoutInvestigation(claimDate)

	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, dueDate)
	}
}

// =============================================================================
// Test BR-CLM-DC-004: SLA with investigation (45 days)
// =============================================================================

func TestCalculateClaimSLAWithInvestigation(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedDueDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC) // 45 days later

	dueDate := service.CalculateClaimSLAWithInvestigation(claimDate)

	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, dueDate)
	}
}

// =============================================================================
// Test Combined SLA Calculation
// =============================================================================

func TestCalculateClaimSLADueDate_WithInvestigation(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedDueDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC) // 45 days

	dueDate := service.CalculateClaimSLADueDate(claimDate, true)

	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, dueDate)
	}
}

func TestCalculateClaimSLADueDate_WithoutInvestigation(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days

	dueDate := service.CalculateClaimSLADueDate(claimDate, false)

	if !dueDate.Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, dueDate)
	}
}

// =============================================================================
// Test BR-CLM-DC-009: Penal interest calculation (8% p.a.)
// =============================================================================

func TestCalculatePenalInterest_ValidInputs(t *testing.T) {
	service := NewBusinessRulesService()

	claimAmount := 100000.0
	breachDays := 10

	// Formula: 100000 * 0.08 * 10 / 365 = 219.18
	expectedInterest := 219.18
	penalInterest := service.CalculatePenalInterest(claimAmount, breachDays)

	diff := math.Abs(penalInterest - expectedInterest)
	if diff > 0.01 {
		t.Errorf("Expected penal interest ~%.2f, got %.2f", expectedInterest, penalInterest)
	}
}

func TestCalculatePenalInterest_ZeroClaimAmount(t *testing.T) {
	service := NewBusinessRulesService()

	penalInterest := service.CalculatePenalInterest(0.0, 10)

	if penalInterest != 0.0 {
		t.Errorf("Expected zero penal interest for zero claim amount, got %.2f", penalInterest)
	}
}

func TestCalculatePenalInterest_ZeroBreachDays(t *testing.T) {
	service := NewBusinessRulesService()

	penalInterest := service.CalculatePenalInterest(100000.0, 0)

	if penalInterest != 0.0 {
		t.Errorf("Expected zero penal interest for zero breach days, got %.2f", penalInterest)
	}
}

func TestCalculatePenalInterestWithDate_NoBreach(t *testing.T) {
	service := NewBusinessRulesService()

	claimAmount := 100000.0
	slaDueDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	settlementDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC) // Before SLA

	penalInterest, breachDays := service.CalculatePenalInterestWithDate(claimAmount, slaDueDate, settlementDate)

	if penalInterest != 0.0 {
		t.Errorf("Expected zero penal interest when no breach, got %.2f", penalInterest)
	}

	if breachDays != 0 {
		t.Errorf("Expected zero breach days, got %d", breachDays)
	}
}

func TestCalculatePenalInterestWithDate_WithBreach(t *testing.T) {
	service := NewBusinessRulesService()

	claimAmount := 100000.0
	slaDueDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	settlementDate := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC) // 10 days after SLA

	penalInterest, breachDays := service.CalculatePenalInterestWithDate(claimAmount, slaDueDate, settlementDate)

	if breachDays != 10 {
		t.Errorf("Expected 10 breach days, got %d", breachDays)
	}

	expectedInterest := 219.18
	diff := math.Abs(penalInterest - expectedInterest)
	if diff > 0.01 {
		t.Errorf("Expected penal interest ~%.2f, got %.2f", expectedInterest, penalInterest)
	}
}

// =============================================================================
// Test BR-CLM-DC-021: SLA color coding (GREEN/YELLOW/ORANGE/RED)
// =============================================================================

func TestCalculateSLAStatus_Green(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	slaDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days total
	currentDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC) // 9 days elapsed (60%)

	status := service.CalculateSLAStatus(claimDate, slaDueDate, currentDate)

	if status != SLAStatusGreen {
		t.Errorf("Expected GREEN status, got %s", status)
	}
}

func TestCalculateSLAStatus_Yellow(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	slaDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days total
	currentDate := time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC) // 11 days elapsed (73%)

	status := service.CalculateSLAStatus(claimDate, slaDueDate, currentDate)

	if status != SLAStatusYellow {
		t.Errorf("Expected YELLOW status, got %s", status)
	}
}

func TestCalculateSLAStatus_Orange(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	slaDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days total
	currentDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // 14 days elapsed (93%)

	status := service.CalculateSLAStatus(claimDate, slaDueDate, currentDate)

	if status != SLAStatusOrange {
		t.Errorf("Expected ORANGE status, got %s", status)
	}
}

func TestCalculateSLAStatus_Red(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	slaDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days total
	currentDate := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC) // 16 days elapsed (>100%)

	status := service.CalculateSLAStatus(claimDate, slaDueDate, currentDate)

	if status != SLAStatusRed {
		t.Errorf("Expected RED status, got %s", status)
	}
}

func TestCalculateSLAPercentageRemaining(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	slaDueDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 15 days total
	currentDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC) // 9 days elapsed

	percentageRemaining := service.CalculateSLAPercentageRemaining(claimDate, slaDueDate, currentDate)

	expectedPercentage := 40.0 // 100 - 60 = 40%
	diff := math.Abs(percentageRemaining - expectedPercentage)
	if diff > 0.01 {
		t.Errorf("Expected ~%.2f%% remaining, got %.2f%%", expectedPercentage, percentageRemaining)
	}
}

func TestCalculateSLAStatusForClaim(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Now().Add(-10 * 24 * time.Hour)
	slaDueDate := time.Now().Add(5 * 24 * time.Hour)

	claim := domain.Claim{
		ClaimDate:  claimDate,
		SLADueDate: slaDueDate,
	}

	status := service.CalculateSLAStatusForClaim(claim)

	if status != SLAStatusGreen {
		t.Errorf("Expected GREEN status for claim, got %s", status)
	}
}

// =============================================================================
// Test Claim Amount Calculation (BR-CLM-DC-008)
// =============================================================================

func TestCalculateClaimAmount_AllComponents(t *testing.T) {
	service := NewBusinessRulesService()

	sumAssured := 500000.0
	reversionaryBonus := 50000.0
	terminalBonus := 20000.0
	outstandingLoan := 10000.0
	unpaidPremiums := 5000.0

	// 500000 + 50000 + 20000 - 10000 - 5000 = 555000
	expectedAmount := 555000.0
	claimAmount := service.CalculateClaimAmount(
		sumAssured,
		reversionaryBonus,
		terminalBonus,
		outstandingLoan,
		unpaidPremiums,
	)

	if claimAmount != expectedAmount {
		t.Errorf("Expected claim amount %.2f, got %.2f", expectedAmount, claimAmount)
	}
}

func TestCalculateClaimAmount_NegativeResult(t *testing.T) {
	service := NewBusinessRulesService()

	sumAssured := 10000.0
	reversionaryBonus := 0.0
	terminalBonus := 0.0
	outstandingLoan := 15000.0 // More than sum assured
	unpaidPremiums := 5000.0

	claimAmount := service.CalculateClaimAmount(
		sumAssured,
		reversionaryBonus,
		terminalBonus,
		outstandingLoan,
		unpaidPremiums,
	)

	if claimAmount != 0.0 {
		t.Errorf("Expected zero claim amount for negative result, got %.2f", claimAmount)
	}
}

// =============================================================================
// Test Document Checklist Rules
// =============================================================================

func TestGetRequiredDocumentTypes_NaturalDeathWithNomination(t *testing.T) {
	service := NewBusinessRulesService()

	documents := service.GetRequiredDocumentTypes(DeathTypeNatural, true)

	expectedDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
		"CLAIMANT_ID_PROOF",
		"CLAIMANT_PHOTOGRAPH",
		"POLICY_BOND_OR_INDEMNITY_BOND",
		"BANK_ACCOUNT_PROOF",
	}

	if len(documents) != len(expectedDocuments) {
		t.Errorf("Expected %d documents, got %d", len(expectedDocuments), len(documents))
	}

	for _, expected := range expectedDocuments {
		found := false
		for _, doc := range documents {
			if doc == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected document %s not found", expected)
		}
	}
}

func TestGetRequiredDocumentTypes_UnnaturalDeathWithNomination(t *testing.T) {
	service := NewBusinessRulesService()

	documents := service.GetRequiredDocumentTypes(DeathTypeUnnatural, true)

	expectedConditionalDocs := []string{
		"POST_MORTEM_REPORT",
		"FIR_COPY",
		"POLICE_REPORT",
		"AUTOPSY_REPORT",
	}

	for _, expected := range expectedConditionalDocs {
		found := false
		for _, doc := range documents {
			if doc == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected conditional document %s not found for unnatural death", expected)
		}
	}
}

func TestGetRequiredDocumentTypes_NoNomination(t *testing.T) {
	service := NewBusinessRulesService()

	documents := service.GetRequiredDocumentTypes(DeathTypeNatural, false)

	expectedNominationDocs := []string{
		"SUCCESSION_CERTIFICATE",
		"LEGAL_HEIR_CERTIFICATE",
		"NO_NOMINATION_AFFIDAVIT",
	}

	for _, expected := range expectedNominationDocs {
		found := false
		for _, doc := range documents {
			if doc == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected nomination-related document %s not found", expected)
		}
	}
}

func TestIsDocumentComplete_Complete(t *testing.T) {
	service := NewBusinessRulesService()

	uploadedDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
		"CLAIMANT_ID_PROOF",
		"CLAIMANT_PHOTOGRAPH",
		"POLICY_BOND_OR_INDEMNITY_BOND",
		"BANK_ACCOUNT_PROOF",
	}

	requiredDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
		"CLAIMANT_ID_PROOF",
		"CLAIMANT_PHOTOGRAPH",
		"POLICY_BOND_OR_INDEMNITY_BOND",
		"BANK_ACCOUNT_PROOF",
	}

	isComplete := service.IsDocumentComplete(uploadedDocuments, requiredDocuments)

	if !isComplete {
		t.Errorf("Expected documents to be complete")
	}
}

func TestIsDocumentComplete_Incomplete(t *testing.T) {
	service := NewBusinessRulesService()

	uploadedDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
	}

	requiredDocuments := []string{
		"DEATH_CERTIFICATE",
		"CLAIM_FORM",
		"CLAIMANT_ID_PROOF",
		"CLAIMANT_PHOTOGRAPH",
		"POLICY_BOND_OR_INDEMNITY_BOND",
		"BANK_ACCOUNT_PROOF",
	}

	isComplete := service.IsDocumentComplete(uploadedDocuments, requiredDocuments)

	if isComplete {
		t.Errorf("Expected documents to be incomplete")
	}
}

// =============================================================================
// Test Approval Hierarchy (BR-CLM-DC-022)
// =============================================================================

func TestGetApprovalLevel_Level1(t *testing.T) {
	service := NewBusinessRulesService()

	level := service.GetApprovalLevel(50000.0)

	if level != "LEVEL_1" {
		t.Errorf("Expected LEVEL_1, got %s", level)
	}
}

func TestGetApprovalLevel_Level2(t *testing.T) {
	service := NewBusinessRulesService()

	level := service.GetApprovalLevel(250000.0)

	if level != "LEVEL_2" {
		t.Errorf("Expected LEVEL_2, got %s", level)
	}
}

func TestGetApprovalLevel_Level3(t *testing.T) {
	service := NewBusinessRulesService()

	level := service.GetApprovalLevel(750000.0)

	if level != "LEVEL_3" {
		t.Errorf("Expected LEVEL_3, got %s", level)
	}
}

func TestGetApprovalLevel_Level4(t *testing.T) {
	service := NewBusinessRulesService()

	level := service.GetApprovalLevel(1500000.0)

	if level != "LEVEL_4" {
		t.Errorf("Expected LEVEL_4, got %s", level)
	}
}

// =============================================================================
// Test Reinvestigation Limit (BR-CLM-DC-023)
// =============================================================================

func TestCanReinvestigate_UnderLimit(t *testing.T) {
	service := NewBusinessRulesService()

	canReinvestigate := service.CanReinvestigate(1)

	if !canReinvestigate {
		t.Errorf("Expected reinvestigation to be allowed with 1 previous reinvestigation")
	}
}

func TestCanReinvestigate_AtLimit(t *testing.T) {
	service := NewBusinessRulesService()

	canReinvestigate := service.CanReinvestigate(2)

	if canReinvestigate {
		t.Errorf("Expected reinvestigation NOT to be allowed with 2 previous reinvestigations")
	}
}

// =============================================================================
// Test Payment Mode Priority (BR-CLM-DC-017)
// =============================================================================

func TestGetPreferredPaymentMode_NEFT(t *testing.T) {
	service := NewBusinessRulesService()

	mode := service.GetPreferredPaymentMode(true, true)

	if mode != "NEFT" {
		t.Errorf("Expected NEFT, got %s", mode)
	}
}

func TestGetPreferredPaymentMode_POSB(t *testing.T) {
	service := NewBusinessRulesService()

	mode := service.GetPreferredPaymentMode(false, true)

	if mode != "POSB" {
		t.Errorf("Expected POSB, got %s", mode)
	}
}

func TestGetPreferredPaymentMode_Cheque(t *testing.T) {
	service := NewBusinessRulesService()

	mode := service.GetPreferredPaymentMode(false, false)

	if mode != "CHEQUE" {
		t.Errorf("Expected CHEQUE, got %s", mode)
	}
}

// =============================================================================
// Test Communication Triggers (BR-CLM-DC-019)
// =============================================================================

func TestShouldSendDocumentReminder_ShouldSend(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Now().Add(-10 * 24 * time.Hour) // 10 days ago
	currentDate := time.Now()

	shouldSend := service.ShouldSendDocumentReminder(claimDate, currentDate, false)

	if !shouldSend {
		t.Errorf("Expected document reminder to be sent after 7 days")
	}
}

func TestShouldSendDocumentReminder_ShouldNotSend_TooEarly(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Now().Add(-5 * 24 * time.Hour) // 5 days ago
	currentDate := time.Now()

	shouldSend := service.ShouldSendDocumentReminder(claimDate, currentDate, false)

	if shouldSend {
		t.Errorf("Expected document reminder NOT to be sent before 7 days")
	}
}

func TestShouldSendDocumentReminder_ShouldNotSend_Complete(t *testing.T) {
	service := NewBusinessRulesService()

	claimDate := time.Now().Add(-10 * 24 * time.Hour) // 10 days ago
	currentDate := time.Now()

	shouldSend := service.ShouldSendDocumentReminder(claimDate, currentDate, true)

	if shouldSend {
		t.Errorf("Expected document reminder NOT to be sent when documents are complete")
	}
}

func TestShouldSendSLABreachedWarning(t *testing.T) {
	service := NewBusinessRulesService()

	shouldSend := service.ShouldSendSLABreachedWarning(SLAStatusOrange)

	if !shouldSend {
		t.Errorf("Expected SLA breach warning to be sent for ORANGE status")
	}
}

func TestShouldSendSLABreachedWarning_NotWarning(t *testing.T) {
	service := NewBusinessRulesService()

	shouldSend := service.ShouldSendSLABreachedWarning(SLAStatusGreen)

	if shouldSend {
		t.Errorf("Expected SLA breach warning NOT to be sent for GREEN status")
	}
}

// =============================================================================
// Test Helper Functions
// =============================================================================

func TestIsValidDeathType_Valid(t *testing.T) {
	service := NewBusinessRulesService()

	validTypes := []string{DeathTypeNatural, DeathTypeAccidental, DeathTypeUnnatural}

	for _, deathType := range validTypes {
		if !service.IsValidDeathType(deathType) {
			t.Errorf("Expected %s to be valid death type", deathType)
		}
	}
}

func TestIsValidDeathType_Invalid(t *testing.T) {
	service := NewBusinessRulesService()

	if service.IsValidDeathType("INVALID") {
		t.Errorf("Expected INVALID to be invalid death type")
	}
}

func TestIsValidClaimType_Valid(t *testing.T) {
	service := NewBusinessRulesService()

	validTypes := []string{"DEATH", "MATURITY", "SURVIVAL_BENEFIT", "FREELOOK"}

	for _, claimType := range validTypes {
		if !service.IsValidClaimType(claimType) {
			t.Errorf("Expected %s to be valid claim type", claimType)
		}
	}
}

func TestIsValidClaimType_Invalid(t *testing.T) {
	service := NewBusinessRulesService()

	if service.IsValidClaimType("INVALID") {
		t.Errorf("Expected INVALID to be invalid claim type")
	}
}
