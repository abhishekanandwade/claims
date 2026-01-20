package response

import (
	"time"

	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== NOTIFICATION RESPONSE DTOs ====================

// NotificationSentResponse represents response after sending notification
// POST /notifications/send
// Reference: BR-CLM-DC-019 (Communication triggers)
type NotificationSentResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	NotificationID            string   `json:"notification_id"` // Unique notification ID for tracking
	Status                    string   `json:"status"`          // SENT, FAILED, PENDING
	ChannelsSent              []string `json:"channels_sent"`   // Channels successfully sent
	ChannelsFailed            []string `json:"channels_failed,omitempty"` // Channels that failed
	SentAt                    string   `json:"sent_at"`         // ISO 8601 timestamp
}

// BatchNotificationsSentResponse represents response after sending batch notifications
// POST /notifications/send-batch
type BatchNotificationsSentResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	BatchID                   string `json:"batch_id"` // Batch ID for tracking
	TotalNotifications        int    `json:"total_notifications"`
	Successful                int    `json:"successful"`
	Failed                    int    `json:"failed"`
	NotificationResults       []NotificationResult `json:"notification_results"`
	SentAt                    string `json:"sent_at"` // ISO 8601 timestamp
}

// NotificationResult represents result of individual notification in batch
type NotificationResult struct {
	NotificationID string   `json:"notification_id"`
	Status         string   `json:"status"` // SENT, FAILED, PENDING
	ChannelsSent   []string `json:"channels_sent"`
	ChannelsFailed []string `json:"channels_failed,omitempty"`
	Error          *string  `json:"error,omitempty"`
}

// FeedbackLinkGeneratedResponse represents response after generating feedback link
// POST /feedback/generate-link
// Reference: BR-CLM-DC-020 (Customer feedback)
type FeedbackLinkGeneratedResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	FeedbackID                string `json:"feedback_id"`         // Unique feedback ID
	FeedbackURL               string `json:"feedback_url"`        // Full URL for feedback form
	FeedbackToken             string `json:"feedback_token"`      // Token for authentication
	ExpiryDate                string `json:"expiry_date"`         // ISO 8601 timestamp
	ClaimID                   string `json:"claim_id"`
	FeedbackType              string `json:"feedback_type"`       // CLAIM_PROCESSING, CUSTOMER_SERVICE, DOCUMENT_QUALITY
	Questions                 []FeedbackQuestion `json:"questions,omitempty"` // Pre-defined questions
	CreatedAt                 string `json:"created_at"`          // ISO 8601 timestamp
}

// FeedbackQuestion represents a feedback question
type FeedbackQuestion struct {
	QuestionID   string   `json:"question_id"`
	QuestionText string   `json:"question_text"`
	QuestionType string   `json:"question_type"` // RATING, TEXT, MULTIPLE_CHOICE
	Options      []string `json:"options,omitempty"` // For multiple choice questions
	Required     bool     `json:"required"`
	MaxRating    *int     `json:"max_rating,omitempty"` // For rating questions (1-5 scale)
}

// NotificationStatusResponse represents notification status details
// GET /notifications/{notification_id}/status (Additional endpoint for tracking)
type NotificationStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	NotificationID            string `json:"notification_id"`
	Status                    string `json:"status"` // SENT, DELIVERED, FAILED, PENDING
	Channels                  []ChannelStatus `json:"channels"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
}

// ChannelStatus represents status of notification per channel
type ChannelStatus struct {
	Channel      string    `json:"channel"` // SMS, EMAIL, WHATSAPP, PUSH
	Status       string    `json:"status"`  // SENT, DELIVERED, FAILED, PENDING
	DeliveredAt  *string   `json:"delivered_at,omitempty"`
	FailedAt     *string   `json:"failed_at,omitempty"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	RetryCount   int       `json:"retry_count"`
}

// ==================== HELPER FUNCTIONS ====================

// NewNotificationSentResponse creates a new notification sent response
func NewNotificationSentResponse(notificationID, status string, channelsSent, channelsFailed []string) *NotificationSentResponse {
	return &NotificationSentResponse{
		NotificationID: notificationID,
		Status:         status,
		ChannelsSent:   channelsSent,
		ChannelsFailed: channelsFailed,
		SentAt:         time.Now().Format("2006-01-02 15:04:05"),
	}
}

// NewBatchNotificationsSentResponse creates a new batch notifications sent response
func NewBatchNotificationsSentResponse(batchID string, total, successful, failed int, results []NotificationResult) *BatchNotificationsSentResponse {
	return &BatchNotificationsSentResponse{
		BatchID:             batchID,
		TotalNotifications:  total,
		Successful:          successful,
		Failed:              failed,
		NotificationResults: results,
		SentAt:              time.Now().Format("2006-01-02 15:04:05"),
	}
}

// NewFeedbackLinkGeneratedResponse creates a new feedback link generated response
func NewFeedbackLinkGeneratedResponse(feedbackID, feedbackURL, feedbackToken, expiryDate, claimID, feedbackType string) *FeedbackLinkGeneratedResponse {
	// Add default questions for claim processing feedback
	questions := []FeedbackQuestion{
		{
			QuestionID:   "Q1",
			QuestionText: "How satisfied are you with the claim processing experience?",
			QuestionType: "RATING",
			Required:     true,
			MaxRating:    intPtr(5),
		},
		{
			QuestionID:   "Q2",
			QuestionText: "How would you rate the timeliness of the claim settlement?",
			QuestionType: "RATING",
			Required:     true,
			MaxRating:    intPtr(5),
		},
		{
			QuestionID:   "Q3",
			QuestionText: "How would you rate the communication from our team?",
			QuestionType: "RATING",
			Required:     true,
			MaxRating:    intPtr(5),
		},
		{
			QuestionID:   "Q4",
			QuestionText: "Please share any additional feedback or suggestions:",
			QuestionType: "TEXT",
			Required:     false,
		},
	}

	return &FeedbackLinkGeneratedResponse{
		FeedbackID:    feedbackID,
		FeedbackURL:   feedbackURL,
		FeedbackToken: feedbackToken,
		ExpiryDate:    expiryDate,
		ClaimID:       claimID,
		FeedbackType:  feedbackType,
		Questions:     questions,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}
}

// Helper function for pointer to int
func intPtr(i int) *int {
	return &i
}
