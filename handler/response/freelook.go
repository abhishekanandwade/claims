package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ========================================
// FREE LOOK CANCELLATION - RESPONSE DTOS
// ========================================

// PolicyBondTrackedResponse represents response for tracking policy bond
// POST /policy-bond/track
// Reference: FR-CLM-BOND-001
type PolicyBondTrackedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ID                       string `json:"id"`
	TrackingNumber           string `json:"tracking_number"`
	PolicyID                 string `json:"policy_id"`
	BondType                 string `json:"bond_type"`
	DispatchDate             string `json:"dispatch_date"`          // YYYY-MM-DD HH:MM:SS format
	EstimatedDeliveryDate    string `json:"estimated_delivery_date"` // YYYY-MM-DD HH:MM:SS format
	CreatedAt                string `json:"created_at"`             // YYYY-MM-DD HH:MM:SS format
}

// BondDeliveryUpdatedResponse represents response for updating bond delivery status
// POST /policy-bond/{bond_id}/delivery-status
// Reference: FR-CLM-BOND-002
type BondDeliveryUpdatedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	ID                       string `json:"id"`
	DeliveryStatus           string `json:"delivery_status"`
	DeliveryDate             *string `json:"delivery_date,omitempty"`        // YYYY-MM-DD HH:MM:SS format
	FreeLookStartDate        string `json:"free_look_start_date"` // YYYY-MM-DD HH:MM:SS format
	FreeLookEndDate          string `json:"free_look_end_date"`   // YYYY-MM-DD HH:MM:SS format
	DaysRemaining            int    `json:"days_remaining"`
	UpdatedAt                string `json:"updated_at"`           // YYYY-MM-DD HH:MM:SS format
}

// PolicyBondDetailsResponse represents response for policy bond details
// GET /policy-bond/{bond_id}/details
type PolicyBondDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                 string  `json:"policy_id"`
	BondType                 string  `json:"bond_type"`
	DispatchDate             string  `json:"dispatch_date,omitempty"`               // YYYY-MM-DD HH:MM:SS format
	TrackingNumber           *string `json:"tracking_number,omitempty"`
	DeliveryStatus           string  `json:"delivery_status"`
	DeliveryDate             *string `json:"delivery_date,omitempty"`     // YYYY-MM-DD HH:MM:SS format
	DeliveryAttempt          *int    `json:"delivery_attempt,omitempty"`
	FreeLookStartDate        *string `json:"free_look_start_date,omitempty"` // YYYY-MM-DD HH:MM:SS format
	FreeLookEndDate          *string `json:"free_look_end_date,omitempty"`   // YYYY-MM-DD HH:MM:SS format
	FreeLookStatus           *string `json:"free_look_status,omitempty"`     // ACTIVE, EXPIRED, CANCELLED
	DaysRemaining            *int    `json:"days_remaining,omitempty"`
	EscalationTriggered      bool    `json:"escalation_triggered"`
	CreatedAt                string  `json:"created_at"`    // YYYY-MM-DD HH:MM:SS format
	UpdatedAt                string  `json:"updated_at"`    // YYYY-MM-DD HH:MM:SS format
}

// PolicyBondsListResponse represents response for listing policy bonds
// GET /policy-bond/policy/{policy_id}
type PolicyBondsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	port.MetaDataResponse     `json:",inline"`
	Bonds                     []PolicyBondDetailsResponse `json:"bonds"`
	TotalBonds                int                         `json:"total_bonds"`
}

// FreeLookEligibilityResponse represents response for checking free look eligibility
// GET /freelook/policy/{policy_id}/eligibility
// Reference: BR-CLM-BOND-001, VR-CLM-FL-001
type FreeLookEligibilityResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                 string  `json:"policy_id"`
	Eligible                 bool    `json:"eligible"`
	FreeLookStatus           string  `json:"free_look_status"` // ACTIVE, EXPIRED, NOT_STARTED
	BondType                 string  `json:"bond_type"`         // PHYSICAL, ELECTRONIC
	FreeLookStartDate        *string `json:"free_look_start_date,omitempty"` // YYYY-MM-DD HH:MM:SS format
	FreeLookEndDate          *string `json:"free_look_end_date,omitempty"`   // YYYY-MM-DD HH:MM:SS format
	DaysRemaining            *int    `json:"days_remaining,omitempty"`
	Reason                   *string `json:"reason,omitempty"`
}

// FreeLookRefundCalculationResponse represents response for calculating free look refund
// POST /policies/{policy_id}/freelook-refund-calculation
// Reference: FR-CLM-FL-003, BR-CLM-BOND-003
type FreeLookRefundCalculationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	PolicyID                 string  `json:"policy_id"`
	PremiumPaid              float64 `json:"premium_paid"`
	Deductions               DeductionBreakdown `json:"deductions"`
	RefundAmount             float64 `json:"refund_amount"`
	CalculatedAt             string  `json:"calculated_at"` // YYYY-MM-DD HH:MM:SS format
}

// DeductionBreakdown represents the breakdown of deductions for free look refund
type DeductionBreakdown struct {
	ProportionateRiskPremium float64 `json:"proportionate_risk_premium"`
	StampDuty               float64 `json:"stamp_duty"`
	MedicalExamCharges       float64 `json:"medical_exam_charges"`
	OtherCharges             float64 `json:"other_charges"`
	TotalDeductions          float64 `json:"total_deductions"`
}

// FreeLookCancellationSubmittedResponse represents response for submitting cancellation
// POST /freelook/cancellation/submit
// Reference: FR-CLM-FL-002
type FreeLookCancellationSubmittedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	CancellationID            string `json:"cancellation_id"`
	PolicyID                  string `json:"policy_id"`
	CancellationNumber        string `json:"cancellation_number"`
	CancellationDate          string `json:"cancellation_date"`     // YYYY-MM-DD HH:MM:SS format
	Status                    string `json:"status"`                // SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, PROCESSED
	RefundAmount              *float64 `json:"refund_amount,omitempty"`
	SubmittedAt               string  `json:"submitted_at"`          // YYYY-MM-DD HH:MM:SS format
}

// FreeLookCancellationDetailsResponse represents response for cancellation details
// GET /freelook/cancellation/{cancellation_id}/details
type FreeLookCancellationDetailsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	CancellationID            string  `json:"cancellation_id"`
	CancellationNumber        string  `json:"cancellation_number"`
	PolicyID                  string  `json:"policy_id"`
	BondID                    *string `json:"bond_id,omitempty"`
	CancellationReason        string  `json:"cancellation_reason"`
	Channel                   string  `json:"channel"`                    // ONLINE, PORTAL, POST_OFFICE, CPGRAMS, EMAIL, PHONE
	CancellationDate          string  `json:"cancellation_date"`          // YYYY-MM-DD HH:MM:SS format
	Status                    string  `json:"status"`                     // SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, PROCESSED
	PremiumPaid               *float64 `json:"premium_paid,omitempty"`
	Deductions                *DeductionBreakdown `json:"deductions,omitempty"`
	RefundAmount              *float64 `json:"refund_amount,omitempty"`
	RefundStatus              *string `json:"refund_status,omitempty"`    // PENDING, PROCESSING, SUCCESS, FAILED
	RefundMode                *string `json:"refund_mode,omitempty"`      // NEFT, RTGS, POSB, CHEQUE
	ReferenceNumber           *string `json:"reference_number,omitempty"`
	MakerID                   *string `json:"maker_id,omitempty"`
	CheckerID                 *string `json:"checker_id,omitempty"`
	ReviewedAt                *string `json:"reviewed_at,omitempty"`      // YYYY-MM-DD HH:MM:SS format
	ReviewComments            *string `json:"review_comments,omitempty"`
	OverrideAmount            *float64 `json:"override_amount,omitempty"`
	OverrideReason            *string `json:"override_reason,omitempty"`
	RefundProcessedAt         *string `json:"refund_processed_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	ClaimantName              string  `json:"claimant_name"`
	ClaimantPhone             *string `json:"claimant_phone,omitempty"`
	ClaimantEmail             *string `json:"claimant_email,omitempty"`
	BankAccountNumber         *string `json:"bank_account_number,omitempty"`
	BankIFSCCode              *string `json:"bank_ifsc_code,omitempty"`
	DocumentURLs              []string `json:"document_urls"`
	CreatedAt                 string  `json:"created_at"` // YYYY-MM-DD HH:MM:SS format
	UpdatedAt                 string  `json:"updated_at"` // YYYY-MM-DD HH:MM:SS format
}

// FreeLookCancellationReviewResponse represents response for reviewing cancellation
// POST /freelook/cancellation/{cancellation_id}/review
// Reference: BR-CLM-BOND-004 (Maker-Checker Workflow)
type FreeLookCancellationReviewResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	CancellationID            string `json:"cancellation_id"`
	Status                    string `json:"status"` // APPROVED, REJECTED
	RefundAmount              *float64 `json:"refund_amount,omitempty"`
	ReviewedAt                string `json:"reviewed_at"` // YYYY-MM-DD HH:MM:SS format
	ReviewedBy                string `json:"reviewed_by"`
	ReviewComments            string `json:"review_comments"`
}

// FreeLookRefundProcessedResponse represents response for processing refund
// POST /freelook/cancellation/{cancellation_id}/process-refund
// Reference: FR-CLM-FL-003
type FreeLookRefundProcessedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	CancellationID            string `json:"cancellation_id"`
	PolicyID                  string `json:"policy_id"`
	RefundStatus              string `json:"refund_status"` // PENDING, PROCESSING, SUCCESS, FAILED
	RefundAmount              float64 `json:"refund_amount"`
	RefundMode                string `json:"refund_mode"` // NEFT, RTGS, POSB, CHEQUE
	ReferenceNumber           string `json:"reference_number"`
	ProcessedAt               string  `json:"processed_at"` // YYYY-MM-DD HH:MM:SS format
	ProcessedBy               string  `json:"processed_by"`
}

// FreeLookRefundStatusResponse represents response for refund status
// GET /freelook/cancellation/{cancellation_id}/refund-status
type FreeLookRefundStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	CancellationID            string  `json:"cancellation_id"`
	PolicyID                  string  `json:"policy_id"`
	Status                    string  `json:"status"` // SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, PROCESSED
	RefundStatus              string  `json:"refund_status"` // PENDING, PROCESSING, SUCCESS, FAILED
	RefundAmount              *float64 `json:"refund_amount,omitempty"`
	RefundMode                *string  `json:"refund_mode,omitempty"` // NEFT, RTGS, POSB, CHEQUE
	ReferenceNumber           *string  `json:"reference_number,omitempty"`
	UTRNumber                 *string  `json:"utr_number,omitempty"`
	ProcessedAt               *string  `json:"processed_at,omitempty"` // YYYY-MM-DD HH:MM:SS format
	FailedReason              *string  `:"failed_reason,omitempty"`
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// NewPolicyBondDetailsResponse creates a new policy bond details response from domain model
func NewPolicyBondDetailsResponse(bond domain.PolicyBondTracking) PolicyBondDetailsResponse {
	response := PolicyBondDetailsResponse{
		PolicyID:                 bond.PolicyID,
		BondType:                 bond.BondType,
		DeliveryAttempt:          &bond.DeliveryAttemptCount,
		EscalationTriggered:      bond.EscalationTriggered,
		CreatedAt:                formatTime(bond.CreatedAt),
		UpdatedAt:                formatTime(bond.UpdatedAt),
	}

	if bond.DispatchDate != nil {
		formatted := formatTime(*bond.DispatchDate)
		response.DispatchDate = formatted
	}

	if bond.TrackingNumber != nil {
		response.TrackingNumber = bond.TrackingNumber
	}

	if bond.DeliveryStatus != nil {
		response.DeliveryStatus = *bond.DeliveryStatus
	}

	if bond.DeliveryDate != nil {
		formatted := formatTime(*bond.DeliveryDate)
		response.DeliveryDate = &formatted
	}

	if bond.FreeLookPeriodStartDate != nil {
		formatted := formatTime(*bond.FreeLookPeriodStartDate)
		response.FreeLookStartDate = &formatted
	}

	if bond.FreeLookPeriodEndDate != nil {
		formatted := formatTime(*bond.FreeLookPeriodEndDate)
		response.FreeLookEndDate = &formatted
	}

	// Calculate days remaining
	if bond.FreeLookPeriodEndDate != nil {
		daysRemaining := int(time.Until(*bond.FreeLookPeriodEndDate).Hours() / 24)
		response.DaysRemaining = &daysRemaining

		// Determine free look status
		if daysRemaining < 0 {
			status := "EXPIRED"
			response.FreeLookStatus = &status
		} else if daysRemaining <= 3 {
			status := "EXPIRING_SOON"
			response.FreeLookStatus = &status
		} else {
			status := "ACTIVE"
			response.FreeLookStatus = &status
		}
	}

	return response
}

// NewFreeLookCancellationDetailsResponse creates a new cancellation details response from domain model
func NewFreeLookCancellationDetailsResponse(cancellation domain.FreeLookCancellation) FreeLookCancellationDetailsResponse {
	response := FreeLookCancellationDetailsResponse{
		CancellationID:     cancellation.CancellationID,
		CancellationNumber: cancellation.CancellationNumber,
		PolicyID:           cancellation.PolicyID,
		CancellationReason: cancellation.CancellationReason,
		Channel:            cancellation.Channel,
		CancellationDate:   formatTime(cancellation.CancellationDate),
		Status:             cancellation.Status,
		ClaimantName:       cancellation.ClaimantName,
		ClaimantPhone:      cancellation.ClaimantPhone,
		ClaimantEmail:      cancellation.ClaimantEmail,
		BankAccountNumber:  cancellation.BankAccountNumber,
		BankIFSCCode:       cancellation.BankIFSCCode,
		DocumentURLs:       cancellation.DocumentURLs,
		CreatedAt:          formatTime(cancellation.CreatedAt),
		UpdatedAt:          formatTime(cancellation.UpdatedAt),
	}

	if cancellation.BondID != nil {
		response.BondID = cancellation.BondID
	}

	if cancellation.PremiumPaid != nil {
		response.PremiumPaid = cancellation.PremiumPaid
	}

	if cancellation.ProportionateRiskPremium != nil && cancellation.StampDuty != nil &&
		cancellation.MedicalExamCharges != nil && cancellation.OtherCharges != nil {
		response.Deductions = &DeductionBreakdown{
			ProportionateRiskPremium: *cancellation.ProportionateRiskPremium,
			StampDuty:               *cancellation.StampDuty,
			MedicalExamCharges:       *cancellation.MedicalExamCharges,
			OtherCharges:             *cancellation.OtherCharges,
			TotalDeductions:          *cancellation.ProportionateRiskPremium + *cancellation.StampDuty +
				*cancellation.MedicalExamCharges + *cancellation.OtherCharges,
		}
	}

	if cancellation.RefundAmount != nil {
		response.RefundAmount = cancellation.RefundAmount
	}

	if cancellation.RefundStatus != nil {
		response.RefundStatus = cancellation.RefundStatus
	}

	if cancellation.RefundMode != nil {
		response.RefundMode = cancellation.RefundMode
	}

	if cancellation.ReferenceNumber != nil {
		response.ReferenceNumber = cancellation.ReferenceNumber
	}

	if cancellation.MakerID != nil {
		response.MakerID = cancellation.MakerID
	}

	if cancellation.CheckerID != nil {
		response.CheckerID = cancellation.CheckerID
	}

	if cancellation.ReviewedAt != nil {
		formatted := formatTime(*cancellation.ReviewedAt)
		response.ReviewedAt = &formatted
	}

	if cancellation.ReviewComments != nil {
		response.ReviewComments = cancellation.ReviewComments
	}

	if cancellation.OverrideAmount != nil {
		response.OverrideAmount = cancellation.OverrideAmount
	}

	if cancellation.OverrideReason != nil {
		response.OverrideReason = cancellation.OverrideReason
	}

	if cancellation.RefundProcessedAt != nil {
		formatted := formatTime(*cancellation.RefundProcessedAt)
		response.RefundProcessedAt = &formatted
	}

	return response
}

// calculateFreeLookSLAStatus calculates SLA status for free look period
// Returns: ACTIVE, EXPIRING_SOON, EXPIRED
func calculateFreeLookSLAStatus(daysRemaining *int) string {
	if daysRemaining == nil {
		return "NOT_STARTED"
	}

	if *daysRemaining < 0 {
		return "EXPIRED"
	}

	if *daysRemaining <= 3 {
		return "EXPIRING_SOON"
	}

	return "ACTIVE"
}

