package domain

import (
	"time"
)

// PolicyBondTracking represents policy bond dispatch and delivery tracking
// Table: policy_bond_tracking
// Reference: seed/db/claims_database_schema.sql:584-615
type PolicyBondTracking struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// References
	PolicyID string `db:"policy_id" json:"policy_id"`

	// Bond Information
	BondNumber string `db:"bond_number" json:"bond_number"`
	BondType   string `db:"bond_type" json:"bond_type"` // PHYSICAL, ELECTRONIC

	// Dispatch and Delivery
	PrintDate     *time.Time `db:"print_date" json:"print_date,omitempty"`
	DispatchDate  *time.Time `db:"dispatch_date" json:"dispatch_date,omitempty"`
	TrackingNumber *string   `db:"tracking_number" json:"tracking_number,omitempty"`
	DeliveryDate  *time.Time `db:"delivery_date" json:"delivery_date,omitempty"`
	DeliveryStatus *string   `db:"delivery_status" json:"delivery_status,omitempty"` // PENDING, DISPATCHED, DELIVERED, FAILED, RETURNED
	DeliveryAttemptCount int   `db:"delivery_attempt_count" json:"delivery_attempt_count"` // Max 3

	// Proof of Delivery
	PODReference            *string `db:"pod_reference" json:"pod_reference,omitempty"`
	RecipientName           *string `db:"recipient_name" json:"recipient_name,omitempty"`
	RecipientSignatureCaptured bool  `db:"recipient_signature_captured" json:"recipient_signature_captured"`

	// Delivery Failure Handling
	UndeliveredReason *string `db:"undelivered_reason" json:"undelivered_reason,omitempty"` // BR-CLM-BOND-002
	EscalationTriggered bool    `db:"escalation_triggered" json:"escalation_triggered"`
	EscalationDate      *time.Time `db:"escalation_date" json:"escalation_date,omitempty"`

	// Customer Interaction
	CustomerContacted bool      `db:"customer_contacted" json:"customer_contacted"`
	AddressVerified   bool      `db:"address_verified" json:"address_verified"`
	RedeliveryRequested bool    `db:"redelivery_requested" json:"redelivery_requested"`

	// Free Look Period - BR-CLM-BOND-001
	FreeLookPeriodStartDate *time.Time `db:"freelook_period_start_date" json:"freelook_period_start_date,omitempty"` // Physical: delivery date, Electronic: issue date
	FreeLookPeriodEndDate   *time.Time `db:"freelook_period_end_date" json:"freelook_period_end_date,omitempty"` // Physical: 15 days, Electronic: 30 days
	FreeLookCancellationSubmitted bool `db:"freelook_cancellation_submitted" json:"freelook_cancellation_submitted"`
	FreeLookCancellationID *string    `db:"freelook_cancellation_id" json:"freelook_cancellation_id,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the database table name for PolicyBondTracking
func (PolicyBondTracking) TableName() string {
	return "policy_bond_tracking"
}
