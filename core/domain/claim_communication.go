package domain

import (
	"time"
)

// ClaimCommunication represents multi-channel communication log
// Table: claim_communications
// Reference: seed/db/claims_database_schema.sql:456-476
type ClaimCommunication struct {
	// Primary Key
	ID string `db:"id" json:"id"`

	// Reference
	ClaimID string `db:"claim_id" json:"claim_id"`

	// Communication Information - BR-CLM-DC-019
	CommunicationType string `db:"communication_type" json:"communication_type"` // REGISTRATION, DOCUMENT_STATUS, INVESTIGATION, APPROVAL, PAYMENT
	Channel           string `db:"channel" json:"channel"` // SMS, EMAIL, WHATSAPP, PUSH, POSTAL

	// Recipient Information
	RecipientName   *string `db:"recipient_name" json:"recipient_name,omitempty"`
	RecipientMobile *string `db:"recipient_mobile" json:"recipient_mobile,omitempty"`
	RecipientEmail  *string `db:"recipient_email" json:"recipient_email,omitempty"`

	// Message Content
	TemplateID     *string `db:"template_id" json:"template_id,omitempty"`
	MessageContent *string `db:"message_content" json:"message_content,omitempty"`

	// Delivery Information
	SentAt            time.Time  `db:"sent_at" json:"sent_at"`
	DeliveryStatus    string     `db:"delivery_status" json:"delivery_status"` // SENT, DELIVERED, FAILED, PENDING
	DeliveryTimestamp *time.Time `db:"delivery_timestamp" json:"delivery_timestamp,omitempty"`
	FailureReason     *string    `db:"failure_reason" json:"failure_reason,omitempty"`
	ProviderReference *string    `db:"provider_reference" json:"provider_reference,omitempty"`

	// Additional Data
	Metadata map[string]interface{} `db:"metadata" json:"metadata,omitempty"`

	// Audit Fields
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the database table name for ClaimCommunication
func (ClaimCommunication) TableName() string {
	return "claim_communications"
}
