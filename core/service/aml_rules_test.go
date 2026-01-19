package service

import (
	"context"
	"testing"
	"time"
)

// TestAMLTriggerCodes tests that all AML trigger codes are properly defined
func TestAMLTriggerCodes(t *testing.T) {
	tests := []struct {
		name        string
		triggerCode AMLTriggerCode
		expected    string
	}{
		// Core AML Triggers
		{"AML_001 Cash Threshold", AML_001_CashThreshold, "AML_001"},
		{"AML_002 PAN Mismatch", AML_002_PANMismatch, "AML_002"},
		{"AML_003 Nominee Change", AML_003_NomineeChange, "AML_003"},
		{"AML_004 Frequent Surrenders", AML_004_FrequentSurrenders, "AML_004"},
		{"AML_005 Refund Without Bond", AML_005_RefundWithoutBond, "AML_005"},

		// AML Compliance Rules
		{"AML_006 STR Filing", AML_006_STRFilingTimeline, "AML_006"},
		{"AML_007 CTR Filing", AML_007_CTRFilingSchedule, "AML_007"},
		{"AML_008 CTR Aggregate", AML_008_CTRAggregateMonitoring, "AML_008"},
		{"AML_009 Third Party PAN", AML_009_ThirdPartyPAN, "AML_009"},
		{"AML_010 Regulatory Reporting", AML_010_RegulatoryReporting, "AML_010"},
		{"AML_011 Negative List", AML_011_NegativeListScreening, "AML_011"},
		{"AML_012 Beneficial Ownership", AML_012_BeneficialOwnership, "AML_012"},

		// Transaction Pattern Detection
		{"AML_013 Structured Deposits", AML_013_StructuredDeposits, "AML_013"},
		{"AML_014 Rapid Transaction Flow", AML_014_RapidTransactionFlow, "AML_014"},
		{"AML_015 Circular Transfers", AML_015_CircularTransfers, "AML_015"},
		{"AML_016 High Value First Premium", AML_016_HighValueFirstPremium, "AML_016"},
		{"AML_017 Frequent Policy Changes", AML_017_FrequentPolicyChanges, "AML_017"},
		{"AML_018 Early Surrender Pattern", AML_018_EarlySurrenderPattern, "AML_018"},
		{"AML_019 Multiple Payment Sources", AML_019_MultiplePaymentSources, "AML_019"},
		{"AML_020 Geographical Anomaly", AML_020_GeographicalAnomaly, "AML_020"},

		// Customer Behavior Patterns
		{"AML_021 Unusual Activity Spike", AML_021_UnusualActivitySpike, "AML_021"},
		{"AML_022 Inconsistent Income", AML_022_InconsistentIncomeProfile, "AML_022"},
		{"AML_023 High Risk Jurisdiction", AML_023_HighRiskJurisdiction, "AML_023"},
		{"AML_024 Non Resident", AML_024_NonResidentCustomer, "AML_024"},
		{"AML_025 PEP Family", AML_025_PEPFamilyMember, "AML_025"},
		{"AML_026 Shadow Director", AML_026_ShadowDirectorPattern, "AML_026"},
		{"AML_027 Shell Company", AML_027_ShellCompanyIndicators, "AML_027"},
		{"AML_028 Dormant Activation", AML_028_DormantActivation, "AML_028"},
		{"AML_029 Anomalous Settlement", AML_029_AnomalousSettlementPattern, "AML_029"},
		{"AML_030 International Wire", AML_030_InternationalWireTransfer, "AML_030"},

		// Claim and Payout Patterns
		{"AML_031 Rapid Claim Filing", AML_031_RapidClaimFiling, "AML_031"},
		{"AML_032 Multiple Claims", AML_032_MultipleClaimsShortPeriod, "AML_032"},
		{"AML_033 Claim Amount Anomaly", AML_033_ClaimAmountAnomaly, "AML_033"},
		{"AML_035 Suspicious Beneficiary", AML_035_SuspiciousBeneficiaryChange, "AML_035"},
		{"AML_036 Third Party Claimant", AML_036_ThirdPartyClaimant, "AML_036"},
		{"AML_037 Overdue Claim", AML_037_OverdueClaimFiling, "AML_037"},
		{"AML_039 Frequent Claim Contact", AML_039_FrequentClaimContact, "AML_039"},

		// Agent and Channel Patterns
		{"AML_041 Agent High Volume", AML_041_AgentHighVolume, "AML_041"},
		{"AML_042 Agent Cluster", AML_042_AgentClusterPattern, "AML_042"},
		{"AML_043 Channel Anomaly", AML_043_ChannelAnomaly, "AML_043"},
		{"AML_044 Agent Rapid Turnover", AML_044_AgentRapidTurnover, "AML_044"},
		{"AML_046 Fronting Pattern", AML_046_FrontingPattern, "AML_046"},

		// Product and Feature Patterns
		{"AML_051 High Risk Product", AML_051_HighRiskProductSelection, "AML_051"},
		{"AML_052 Premium Financing", AML_052_PremiumFinancingAbuse, "AML_052"},
		{"AML_053 Policy Loan", AML_053_PolicyLoanAnomaly, "AML_053"},
		{"AML_054 Withdrawal Pattern", AML_054_WithdrawalPattern, "AML_054"},
		{"AML_055 Rider Changes", AML_055_RiderFrequentChanges, "AML_055"},
		{"AML_057 Multiple Policies", AML_057_MultiplePoliciesSameLife, "AML_057"},
		{"AML_058 Over Insurance", AML_058_OverInsurancePattern, "AML_058"},
		{"AML_059 Shortlived Policy", AML_059_ShortlivedPolicyPattern, "AML_059"},
		{"AML_060 Unusual Beneficiary", AML_060_UnusualBeneficiaryDesignation, "AML_060"},

		// Technical and System Patterns
		{"AML_061 IP Address", AML_061_IPAddressAnomaly, "AML_061"},
		{"AML_062 Device Fingerprint", AML_062_DeviceFingerprint, "AML_062"},
		{"AML_063 Bot Activity", AML_063_BotActivityIndicator, "AML_063"},
		{"AML_064 Session Anomaly", AML_064_SessionAnomaly, "AML_064"},
		{"AML_065 Data Inconsistency", AML_065_DataInconsistency, "AML_065"},
		{"AML_067 Synthetic Identity", AML_067_SyntheticIdentity, "AML_067"},
		{"AML_068 Account Takeover", AML_068_AccountTakeover, "AML_068"},
		{"AML_069 Multiple Identity", AML_069_MultipleIdentityUsage, "AML_069"},
		{"AML_070 Anomaly Score", AML_070_AnomalyScoreThreshold, "AML_070"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.triggerCode) != tt.expected {
				t.Errorf("Trigger code = %s, want %s", tt.triggerCode, tt.expected)
			}
		})
	}
}

// TestEvaluateCashThreshold tests AML_001: High Cash Premium Alert
func TestEvaluateCashThreshold(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		paymentMode       string
		amount            float64
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
		expectedFiling    bool
		expectedFilingType FilingType
	}{
		{
			name:              "Cash transaction above threshold",
			paymentMode:       "CASH",
			amount:            60000.0,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
			expectedFiling:    true,
			expectedFilingType: FilingTypeCTR,
		},
		{
			name:              "Cash transaction below threshold",
			paymentMode:       "CASH",
			amount:            40000.0,
			shouldTrigger:     false,
		},
		{
			name:           "Non-cash transaction above threshold",
			paymentMode:    "NEFT",
			amount:         100000.0,
			shouldTrigger:  false,
		},
		{
			name:              "Cash transaction exactly at threshold",
			paymentMode:       "CASH",
			amount:            50001.0,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
			expectedFiling:    true,
			expectedFilingType: FilingTypeCTR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID: "TXN001",
				PaymentMode:    tt.paymentMode,
				Amount:         tt.amount,
			}

			result := engine.EvaluateCashThreshold(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger {
				if result.RiskLevel != tt.expectedRiskLevel {
					t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
				}
				if result.FilingRequired != tt.expectedFiling {
					t.Errorf("FilingRequired = %v, want %v", result.FilingRequired, tt.expectedFiling)
				}
				if result.FilingType != tt.expectedFilingType {
					t.Errorf("FilingType = %s, want %s", result.FilingType, tt.expectedFilingType)
				}
			}
		})
	}
}

// TestEvaluatePANMismatch tests AML_002: PAN Mismatch Alert
func TestEvaluatePANMismatch(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		panVerified       bool
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name:              "PAN not verified",
			panVerified:       false,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelMedium,
		},
		{
			name:          "PAN verified",
			panVerified:   true,
			shouldTrigger: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID: "TXN001",
				CustomerID:     "CUST001",
				PANVerified:    tt.panVerified,
			}

			result := engine.EvaluatePANMismatch(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateNomineeChangePostDeath tests AML_003: Nominee Change Post Death
func TestEvaluateNomineeChangePostDeath(t *testing.T) {
	engine := NewAMLRuleEngine()

	deathDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nomineeChangeBefore := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	nomineeChangeAfter := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                  string
		nomineeChangeDate      *time.Time
		deathDate             *time.Time
		shouldTrigger         bool
		expectedRiskLevel     RiskLevel
		expectedBlock         bool
		expectedFilingType    FilingType
	}{
		{
			name:              "Nominee changed after death",
			nomineeChangeDate:  &nomineeChangeAfter,
			deathDate:         &deathDate,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelCritical,
			expectedBlock:     true,
			expectedFilingType: FilingTypeSTR,
		},
		{
			name:             "Nominee changed before death",
			nomineeChangeDate: &nomineeChangeBefore,
			deathDate:        &deathDate,
			shouldTrigger:    false,
		},
		{
			name:             "No nominee change",
			nomineeChangeDate: nil,
			deathDate:        &deathDate,
			shouldTrigger:    false,
		},
		{
			name:             "No death date",
			nomineeChangeDate: &nomineeChangeAfter,
			deathDate:        nil,
			shouldTrigger:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:     "TXN001",
				PolicyID:          "POL001",
				NomineeChangeDate:  tt.nomineeChangeDate,
				DeathDate:         tt.deathDate,
			}

			result := engine.EvaluateNomineeChangePostDeath(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger {
				if result.RiskLevel != tt.expectedRiskLevel {
					t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
				}
				if result.TransactionBlocked != tt.expectedBlock {
					t.Errorf("TransactionBlocked = %v, want %v", result.TransactionBlocked, tt.expectedBlock)
				}
				if result.FilingType != tt.expectedFilingType {
					t.Errorf("FilingType = %s, want %s", result.FilingType, tt.expectedFilingType)
				}
			}
		})
	}
}

// TestEvaluateFrequentSurrenders tests AML_004: Frequent Surrenders
func TestEvaluateFrequentSurrenders(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		surrenderCount    int
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name:              "More than 3 surrenders",
			surrenderCount:    5,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelMedium,
		},
		{
			name:              "Exactly 3 surrenders",
			surrenderCount:    3,
			shouldTrigger:     false,
		},
		{
			name:           "Less than 3 surrenders",
			surrenderCount: 1,
			shouldTrigger:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:   "TXN001",
				CustomerID:      "CUST001",
				SurrenderCount:   tt.surrenderCount,
			}

			result := engine.EvaluateFrequentSurrenders(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateRefundWithoutBond tests AML_005: Refund Without Bond Delivery
func TestEvaluateRefundWithoutBond(t *testing.T) {
	engine := NewAMLRuleEngine()

	bondDispatchDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	refundBefore := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	refundAfter := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		refundDate        *time.Time
		bondDispatchDate  *time.Time
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name:             "Refund before bond dispatch",
			refundDate:       &refundBefore,
			bondDispatchDate: &bondDispatchDate,
			shouldTrigger:    true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name:             "Refund after bond dispatch",
			refundDate:       &refundAfter,
			bondDispatchDate: &bondDispatchDate,
			shouldTrigger:    false,
		},
		{
			name:             "No refund date",
			refundDate:       nil,
			bondDispatchDate: &bondDispatchDate,
			shouldTrigger:    false,
		},
		{
			name:             "No bond dispatch date",
			refundDate:       &refundBefore,
			bondDispatchDate: nil,
			shouldTrigger:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:     "TXN001",
				PolicyID:          "POL001",
				RefundDate:        tt.refundDate,
				BondDispatchDate:  tt.bondDispatchDate,
			}

			result := engine.EvaluateRefundWithoutBond(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateCTRFilingSchedule tests AML_007: CTR Filing Schedule
func TestEvaluateCTRFilingSchedule(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name                   string
		dailyCashAggregate     float64
		shouldTrigger          bool
		expectedRiskLevel      RiskLevel
		expectedFilingType     FilingType
	}{
		{
			name:               "Daily aggregate exceeds threshold",
			dailyCashAggregate: 1200000.0,
			shouldTrigger:      true,
			expectedRiskLevel:  RiskLevelCritical,
			expectedFilingType: FilingTypeCTR,
		},
		{
			name:               "Daily aggregate below threshold",
			dailyCashAggregate: 500000.0,
			shouldTrigger:      false,
		},
		{
			name:               "Daily aggregate exactly at threshold",
			dailyCashAggregate: 1000001.0,
			shouldTrigger:      true,
			expectedRiskLevel:  RiskLevelCritical,
			expectedFilingType: FilingTypeCTR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:      "TXN001",
				CustomerID:         "CUST001",
				DailyCashAggregate: tt.dailyCashAggregate,
			}

			result := engine.EvaluateCTRFilingSchedule(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger {
				if result.RiskLevel != tt.expectedRiskLevel {
					t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
				}
				if result.FilingType != tt.expectedFilingType {
					t.Errorf("FilingType = %s, want %s", result.FilingType, tt.expectedFilingType)
				}
			}
		})
	}
}

// TestEvaluateThirdPartyPANVerification tests AML_009: Third-Party PAN Verification
func TestEvaluateThirdPartyPANVerification(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name                   string
		thirdPartyPayment      bool
		thirdPartyPANVerified  bool
		shouldTrigger          bool
		expectedRiskLevel      RiskLevel
		expectedBlock          bool
	}{
		{
			name:                  "Third party payment with unverified PAN",
			thirdPartyPayment:     true,
			thirdPartyPANVerified:  false,
			shouldTrigger:         true,
			expectedRiskLevel:     RiskLevelCritical,
			expectedBlock:         true,
		},
		{
			name:                  "Third party payment with verified PAN",
			thirdPartyPayment:     true,
			thirdPartyPANVerified:  true,
			shouldTrigger:         false,
		},
		{
			name:              "No third party payment",
			thirdPartyPayment: false,
			shouldTrigger:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:         "TXN001",
				PolicyID:              "POL001",
				ThirdPartyPayment:     tt.thirdPartyPayment,
				ThirdPartyPANVerified: tt.thirdPartyPANVerified,
			}

			result := engine.EvaluateThirdPartyPANVerification(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger {
				if result.RiskLevel != tt.expectedRiskLevel {
					t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
				}
				if result.TransactionBlocked != tt.expectedBlock {
					t.Errorf("TransactionBlocked = %v, want %v", result.TransactionBlocked, tt.expectedBlock)
				}
			}
		})
	}
}

// TestEvaluateBeneficialOwnershipVerification tests AML_012: Beneficial Ownership Verification
func TestEvaluateBeneficialOwnershipVerification(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		customerType      string
		beneficialOwners  []string
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name:              "Company with beneficial owners",
			customerType:      "COMPANY",
			beneficialOwners:  []string{"OWNER1", "OWNER2"},
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name:              "Trust with beneficial owners",
			customerType:      "TRUST",
			beneficialOwners:  []string{"TRUSTEE1"},
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name:             "NGO with beneficial owners",
			customerType:     "NGO",
			beneficialOwners: []string{"DIRECTOR1", "DIRECTOR2"},
			shouldTrigger:    true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name:             "Individual customer",
			customerType:     "INDIVIDUAL",
			beneficialOwners: []string{},
			shouldTrigger:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := TransactionContext{
				TransactionID:     "TXN001",
				CustomerID:        "CUST001",
				CustomerType:      tt.customerType,
				BeneficialOwners:  tt.beneficialOwners,
			}

			result := engine.EvaluateBeneficialOwnershipVerification(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateStructuredDeposits tests AML_013: Structured Deposits (Smurfing)
func TestEvaluateStructuredDeposits(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		previousTransactions []float64
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name: "Multiple structured deposits",
			previousTransactions: []float64{48000, 48500, 48200, 48800},
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name: "No structured deposits",
			previousTransactions: []float64{60000, 30000, 20000},
			shouldTrigger:     false,
		},
		{
			name: "No transactions",
			previousTransactions: []float64{},
			shouldTrigger:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ExtendedTransactionContext{
				TransactionContext: TransactionContext{
					TransactionID: "TXN001",
					CustomerID:    "CUST001",
				},
				PreviousTransactions: tt.previousTransactions,
			}

			result := engine.EvaluateStructuredDeposits(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateRapidTransactionFlow tests AML_014: Rapid Transaction Flow
func TestEvaluateRapidTransactionFlow(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name                string
		transactionFrequency int
		shouldTrigger       bool
		expectedRiskLevel   RiskLevel
	}{
		{
			name:                "Very high daily frequency",
			transactionFrequency: 28,
			shouldTrigger:       true,
			expectedRiskLevel:   RiskLevelHigh,
		},
		{
			name:                "Moderate weekly frequency",
			transactionFrequency: 15,
			shouldTrigger:       true,
			expectedRiskLevel:   RiskLevelMedium,
		},
		{
			name:                "Normal frequency",
			transactionFrequency: 5,
			shouldTrigger:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ExtendedTransactionContext{
				TransactionContext: TransactionContext{
					TransactionID: "TXN001",
					CustomerID:    "CUST001",
				},
				TransactionFrequency: tt.transactionFrequency,
			}

			result := engine.EvaluateRapidTransactionFlow(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestEvaluateHighValueFirstPremium tests AML_016: High-Value First Premium
func TestEvaluateHighValueFirstPremium(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name              string
		premium           float64
		income            float64
		shouldTrigger     bool
		expectedRiskLevel RiskLevel
	}{
		{
			name:              "Premium much higher than income",
			premium:           500000.0,
			income:            50000.0,
			shouldTrigger:     true,
			expectedRiskLevel: RiskLevelHigh,
		},
		{
			name:              "Premium reasonable for income",
			premium:           25000.0,
			income:            500000.0,
			shouldTrigger:     false,
		},
		{
			name:           "Zero income",
			premium:        50000.0,
			income:         0.0,
			shouldTrigger:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ExtendedTransactionContext{
				TransactionContext: TransactionContext{
					TransactionID: "TXN001",
					CustomerID:    "CUST001",
				},
				PolicyPremium: tt.premium,
				CustomerIncome: tt.income,
			}

			result := engine.EvaluateHighValueFirstPremium(ctx)

			if result.Triggered != tt.shouldTrigger {
				t.Errorf("Triggered = %v, want %v", result.Triggered, tt.shouldTrigger)
			}

			if tt.shouldTrigger && result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", result.RiskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestCalculateRiskScore tests the risk score calculation function
func TestCalculateRiskScore(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name         string
		results      []AMLTriggerResult
		minScore     float64
		maxScore     float64
	}{
		{
			name:     "No triggers",
			results:  []AMLTriggerResult{},
			minScore: 0.0,
			maxScore: 0.0,
		},
		{
			name: "Mixed risk levels",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelLow},
				{Triggered: true, RiskLevel: RiskLevelMedium},
				{Triggered: true, RiskLevel: RiskLevelHigh},
			},
			minScore: 49.99, // Average of 25, 50, 75 = 50
			maxScore: 50.01,
		},
		{
			name: "All critical",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelCritical},
				{Triggered: true, RiskLevel: RiskLevelCritical},
			},
			minScore: 99.99, // Average of 100, 100 = 100
			maxScore: 100.0,
		},
		{
			name: "Some triggered, some not",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelHigh},
				{Triggered: false},
				{Triggered: true, RiskLevel: RiskLevelLow},
			},
			minScore: 33.32, // Average of 75, 0, 25 (only triggered results) = 100/3 = 33.33
			maxScore: 33.34,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := engine.CalculateRiskScore(tt.results)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Score = %f, want between %f and %f", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

// TestDetermineOverallRiskLevel tests overall risk level determination
func TestDetermineOverallRiskLevel(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name             string
		results          []AMLTriggerResult
		expectedRiskLevel RiskLevel
	}{
		{
			name:             "No triggers",
			results:          []AMLTriggerResult{},
			expectedRiskLevel: RiskLevelLow,
		},
		{
			name: "One low risk",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelLow},
			},
			expectedRiskLevel: RiskLevelLow,
		},
		{
			name: "Mixed with critical",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelLow},
				{Triggered: true, RiskLevel: RiskLevelCritical},
				{Triggered: true, RiskLevel: RiskLevelMedium},
			},
			expectedRiskLevel: RiskLevelCritical,
		},
		{
			name: "All high risk",
			results: []AMLTriggerResult{
				{Triggered: true, RiskLevel: RiskLevelHigh},
				{Triggered: true, RiskLevel: RiskLevelHigh},
			},
			expectedRiskLevel: RiskLevelHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			riskLevel := engine.DetermineOverallRiskLevel(tt.results)

			if riskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %s, want %s", riskLevel, tt.expectedRiskLevel)
			}
		})
	}
}

// TestIsSTRFilingRequired tests STR filing requirement check
func TestIsSTRFilingRequired(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name         string
		results      []AMLTriggerResult
		expectedReq  bool
	}{
		{
			name:        "No STR filing required",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeCTR},
			},
			expectedReq: false,
		},
		{
			name:        "STR filing required",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeSTR},
			},
			expectedReq: true,
		},
		{
			name:        "Mixed filing types",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeCTR},
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeSTR},
			},
			expectedReq: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required := engine.IsSTRFilingRequired(tt.results)

			if required != tt.expectedReq {
				t.Errorf("STR Filing Required = %v, want %v", required, tt.expectedReq)
			}
		})
	}
}

// TestIsCTRFilingRequired tests CTR filing requirement check
func TestIsCTRFilingRequired(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name         string
		results      []AMLTriggerResult
		expectedReq  bool
	}{
		{
			name:        "No CTR filing required",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeSTR},
			},
			expectedReq: false,
		},
		{
			name:        "CTR filing required",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeCTR},
			},
			expectedReq: true,
		},
		{
			name:        "Multiple CTR triggers",
			results: []AMLTriggerResult{
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeCTR},
				{Triggered: true, FilingRequired: true, FilingType: FilingTypeCTR},
			},
			expectedReq: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required := engine.IsCTRFilingRequired(tt.results)

			if required != tt.expectedReq {
				t.Errorf("CTR Filing Required = %v, want %v", required, tt.expectedReq)
			}
		})
	}
}

// TestShouldBlockTransaction tests transaction blocking check
func TestShouldBlockTransaction(t *testing.T) {
	engine := NewAMLRuleEngine()

	tests := []struct {
		name         string
		results      []AMLTriggerResult
		expectedBlock bool
	}{
		{
			name: "No blocking required",
			results: []AMLTriggerResult{
				{Triggered: true, TransactionBlocked: false},
			},
			expectedBlock: false,
		},
		{
			name: "One blocked transaction",
			results: []AMLTriggerResult{
				{Triggered: true, TransactionBlocked: true},
			},
			expectedBlock: true,
		},
		{
			name: "Mixed blocking status",
			results: []AMLTriggerResult{
				{Triggered: true, TransactionBlocked: false},
				{Triggered: true, TransactionBlocked: true},
			},
			expectedBlock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldBlock := engine.ShouldBlockTransaction(tt.results)

			if shouldBlock != tt.expectedBlock {
				t.Errorf("Should Block Transaction = %v, want %v", shouldBlock, tt.expectedBlock)
			}
		})
	}
}

// TestEvaluateAllTriggers tests the full trigger evaluation workflow
func TestEvaluateAllTriggers(t *testing.T) {
	engine := NewAMLRuleEngine()

	ctx := context.Background()

	deathDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nomineeChangeAfter := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	txContext := TransactionContext{
		TransactionID:      "TXN001",
		PolicyID:           "POL001",
		CustomerID:         "CUST001",
		Amount:             60000.0,
		PaymentMode:        "CASH",
		Timestamp:          time.Now(),
		PANVerified:        false,
		NomineeChangeDate:  &nomineeChangeAfter,
		DeathDate:         &deathDate,
		SurrenderCount:     5,
		DailyCashAggregate: 1200000.0,
	}

	results, err := engine.EvaluateAllTriggers(ctx, txContext)

	if err != nil {
		t.Fatalf("EvaluateAllTriggers failed: %v", err)
	}

	// Verify we got multiple triggers
	triggeredCount := 0
	for _, result := range results {
		if result.Triggered {
			triggeredCount++
		}
	}

	if triggeredCount == 0 {
		t.Error("Expected at least one trigger to fire, got none")
	}

	// Verify expected triggers are present
	triggerCodes := make(map[AMLTriggerCode]bool)
	for _, result := range results {
		triggerCodes[result.TriggerCode] = result.Triggered
	}

	expectedTriggers := []AMLTriggerCode{
		AML_001_CashThreshold,
		AML_002_PANMismatch,
		AML_003_NomineeChange,
		AML_004_FrequentSurrenders,
		AML_007_CTRFilingSchedule,
		AML_008_CTRAggregateMonitoring,
	}

	for _, expectedTrigger := range expectedTriggers {
		if !triggerCodes[expectedTrigger] {
			t.Errorf("Expected trigger %s to fire, but it didn't", expectedTrigger)
		}
	}

	// Verify overall risk level is critical (due to AML_003)
	overallRisk := engine.DetermineOverallRiskLevel(results)
	if overallRisk != RiskLevelCritical {
		t.Errorf("Expected overall risk level to be CRITICAL, got %s", overallRisk)
	}

	// Verify STR filing is required (due to AML_003)
	strRequired := engine.IsSTRFilingRequired(results)
	if !strRequired {
		t.Error("Expected STR filing to be required")
	}

	// Verify CTR filing is required (due to AML_001, AML_007, AML_008)
	ctrRequired := engine.IsCTRFilingRequired(results)
	if !ctrRequired {
		t.Error("Expected CTR filing to be required")
	}

	// Verify transaction should be blocked (due to AML_003)
	shouldBlock := engine.ShouldBlockTransaction(results)
	if !shouldBlock {
		t.Error("Expected transaction to be blocked")
	}
}
