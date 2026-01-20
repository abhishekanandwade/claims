package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	*serverHandler.Base
	notificationClient *repo.NotificationClient
	claimRepo          *repo.ClaimRepository
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationClient *repo.NotificationClient, claimRepo *repo.ClaimRepository) *NotificationHandler {
	base := serverHandler.New("Notifications").
		SetPrefix("/v1").
		AddPrefix("")
	return &NotificationHandler{
		Base:               base,
		notificationClient: notificationClient,
		claimRepo:          claimRepo,
	}
}

// Routes defines all routes for this handler
func (h *NotificationHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Notifications (5 endpoints)
		serverRoute.POST("/notifications/send", h.SendNotification).Name("Send Notification"),
		serverRoute.POST("/notifications/send-batch", h.SendBatchNotifications).Name("Send Batch Notifications"),
		serverRoute.POST("/feedback/generate-link", h.GenerateFeedbackLink).Name("Generate Feedback Link"),
		serverRoute.GET("/notifications/:notification_id/status", h.GetNotificationStatus).Name("Get Notification Status"),
		serverRoute.POST("/notifications/:notification_id/resend", h.ResendNotification).Name("Resend Notification"),
	}
}

// SendNotification sends a notification via SMS/Email/WhatsApp
// POST /notifications/send
// Reference: BR-CLM-DC-019 (Communication triggers)
// Reference: Multi-channel notification support (SMS, Email, WhatsApp, Push)
func (h *NotificationHandler) SendNotification(sctx *serverRoute.Context, req SendNotificationRequest) (*resp.NotificationSentResponse, error) {
	log.Info(sctx.Ctx, "Sending notification", map[string]interface{}{
		"notification_type": req.NotificationType,
		"claim_id":          req.ClaimID,
		"recipient":         req.Recipient.Name,
		"channels":          req.Channels,
	})

	// Generate notification ID
	notificationID, err := generateNotificationID()
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate notification ID: %v", err)
		return nil, fmt.Errorf("failed to generate notification ID: %w", err)
	}

	// Prepare notification request
	notificationReq := &repo.NotificationRequest{
		NotificationID:   notificationID,
		NotificationType: req.NotificationType,
		ClaimID:          req.ClaimID,
		RecipientName:    req.Recipient.Name,
		RecipientMobile:  req.Recipient.Mobile,
		RecipientEmail:   req.Recipient.Email,
		Channels:         req.Channels,
		CustomMessage:    req.CustomMessage,
	}

	// Send notification via notification client
	channelsSent, channelsFailed, err := h.notificationClient.SendNotification(sctx.Ctx, notificationReq)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to send notification: %v", err)
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	// Log notification sent
	log.Info(sctx.Ctx, "Notification sent successfully", map[string]interface{}{
		"notification_id": notificationID,
		"channels_sent":   channelsSent,
		"channels_failed": channelsFailed,
	})

	return resp.NewNotificationSentResponse(notificationID, "SENT", channelsSent, channelsFailed), nil
}

// SendBatchNotifications sends batch notifications
// POST /notifications/send-batch
// Reference: Batch notifications for bulk operations
func (h *NotificationHandler) SendBatchNotifications(sctx *serverRoute.Context, req SendBatchNotificationsRequest) (*resp.BatchNotificationsSentResponse, error) {
	log.Info(sctx.Ctx, "Sending batch notifications", map[string]interface{}{
		"total_notifications": len(req.Notifications),
	})

	// Generate batch ID
	batchID, err := generateBatchID()
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate batch ID: %v", err)
		return nil, fmt.Errorf("failed to generate batch ID: %w", err)
	}

	// Process notifications in batch
	results := make([]resp.NotificationResult, len(req.Notifications))
	successful := 0
	failed := 0

	for i, notificationReq := range req.Notifications {
		// Generate notification ID
		notificationID, err := generateNotificationID()
		if err != nil {
			log.Error(sctx.Ctx, "Failed to generate notification ID: %v", err)
			results[i] = resp.NotificationResult{
				NotificationID: "",
				Status:         "FAILED",
				ChannelsSent:   []string{},
				ChannelsFailed: notificationReq.Channels,
				Error:          stringPtr("Failed to generate notification ID"),
			}
			failed++
			continue
		}

		// Prepare notification request
		req := &repo.NotificationRequest{
			NotificationID:   notificationID,
			NotificationType: notificationReq.NotificationType,
			ClaimID:          notificationReq.ClaimID,
			RecipientName:    notificationReq.Recipient.Name,
			RecipientMobile:  notificationReq.Recipient.Mobile,
			RecipientEmail:   notificationReq.Recipient.Email,
			Channels:         notificationReq.Channels,
			CustomMessage:    notificationReq.CustomMessage,
		}

		// Send notification
		channelsSent, channelsFailed, err := h.notificationClient.SendNotification(sctx.Ctx, req)
		if err != nil {
			log.Error(sctx.Ctx, "Failed to send notification: %v", err)
			results[i] = resp.NotificationResult{
				NotificationID: notificationID,
				Status:         "FAILED",
				ChannelsSent:   []string{},
				ChannelsFailed: notificationReq.Channels,
				Error:          stringPtr(err.Error()),
			}
			failed++
		} else {
			results[i] = resp.NotificationResult{
				NotificationID: notificationID,
				Status:         "SENT",
				ChannelsSent:   channelsSent,
				ChannelsFailed: channelsFailed,
			}
			successful++
		}
	}

	// Log batch notification results
	log.Info(sctx.Ctx, "Batch notifications sent", map[string]interface{}{
		"batch_id":   batchID,
		"total":      len(req.Notifications),
		"successful": successful,
		"failed":     failed,
	})

	return resp.NewBatchNotificationsSentResponse(batchID, len(req.Notifications), successful, failed, results), nil
}

// GenerateFeedbackLink generates a customer feedback link
// POST /feedback/generate-link
// Reference: BR-CLM-DC-020 (Customer feedback)
func (h *NotificationHandler) GenerateFeedbackLink(sctx *serverRoute.Context, req GenerateFeedbackLinkRequest) (*resp.FeedbackLinkGeneratedResponse, error) {
	log.Info(sctx.Ctx, "Generating feedback link", map[string]interface{}{
		"claim_id": req.ClaimID,
	})

	// Verify claim exists
	_, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		log.Error(sctx.Ctx, "Claim not found: %v", err)
		return nil, fmt.Errorf("claim not found: %w", err)
	}

	// Generate feedback ID and token
	feedbackID, err := generateFeedbackID()
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate feedback ID: %v", err)
		return nil, fmt.Errorf("failed to generate feedback ID: %w", err)
	}

	feedbackToken, err := generateFeedbackToken()
	if err != nil {
		log.Error(sctx.Ctx, "Failed to generate feedback token: %v", err)
		return nil, fmt.Errorf("failed to generate feedback token: %w", err)
	}

	// Determine expiry date
	expiryDays := 7
	if req.ExpiryDays != nil {
		expiryDays = *req.ExpiryDays
	}
	expiryDate := time.Now().AddDate(0, 0, expiryDays).Format("2006-01-02 15:04:05")

	// Determine feedback type
	feedbackType := "CLAIM_PROCESSING"
	if req.FeedbackType != nil {
		feedbackType = *req.FeedbackType
	}

	// Generate feedback URL (placeholder URL)
	// TODO: Replace with actual feedback form URL from configuration
	feedbackURL := fmt.Sprintf("https://pli.gov.in/feedback/%s?token=%s", feedbackID, feedbackToken)

	// Log feedback link generation
	log.Info(sctx.Ctx, "Feedback link generated successfully", map[string]interface{}{
		"feedback_id":   feedbackID,
		"claim_id":      req.ClaimID,
		"feedback_type": feedbackType,
		"expiry_date":   expiryDate,
	})

	return resp.NewFeedbackLinkGeneratedResponse(feedbackID, feedbackURL, feedbackToken, expiryDate, req.ClaimID, feedbackType), nil
}

// GetNotificationStatus gets the status of a notification
// GET /notifications/{notification_id}/status
// Reference: Notification tracking and monitoring
func (h *NotificationHandler) GetNotificationStatus(sctx *serverRoute.Context, req NotificationIDUri) (*resp.NotificationStatusResponse, error) {
	log.Info(sctx.Ctx, "Getting notification status", map[string]interface{}{
		"notification_id": req.NotificationID,
	})

	// TODO: Get notification status from database or notification service
	// Placeholder implementation
	status := "DELIVERED"
	channels := []resp.ChannelStatus{
		{
			Channel:     "SMS",
			Status:      "DELIVERED",
			DeliveredAt: stringPtr(time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")),
			RetryCount:  0,
		},
		{
			Channel:     "EMAIL",
			Status:      "DELIVERED",
			DeliveredAt: stringPtr(time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")),
			RetryCount:  0,
		},
	}

	return &resp.NotificationStatusResponse{
		NotificationID: req.NotificationID,
		Status:         status,
		Channels:       channels,
		CreatedAt:      time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
		UpdatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// ResendNotification resends a failed notification
// POST /notifications/{notification_id}/resend
// Reference: Notification retry mechanism
func (h *NotificationHandler) ResendNotification(sctx *serverRoute.Context, req NotificationIDUri) (*resp.NotificationSentResponse, error) {
	log.Info(sctx.Ctx, "Resending notification", map[string]interface{}{
		"notification_id": req.NotificationID,
	})

	// TODO: Get notification details from database
	// TODO: Resend notification via notification client
	// Placeholder implementation
	channelsSent := []string{"SMS", "EMAIL"}
	channelsFailed := []string{}

	log.Info(sctx.Ctx, "Notification resent successfully", map[string]interface{}{
		"notification_id": req.NotificationID,
		"channels_sent":   channelsSent,
	})

	return resp.NewNotificationSentResponse(req.NotificationID, "SENT", channelsSent, channelsFailed), nil
}

// NotificationIDUri represents the notification_id URI parameter
// GET /notifications/{notification_id}/status
// POST /notifications/{notification_id}/resend
type NotificationIDUri struct {
	NotificationID string `uri:"notification_id" validate:"required"`
}

// ==================== HELPER FUNCTIONS ====================

// generateNotificationID generates a unique notification ID
func generateNotificationID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("NOTIF%s", hex.EncodeToString(b)[:16].ToUpper()), nil
}

// generateBatchID generates a unique batch ID
func generateBatchID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("BATCH%s", hex.EncodeToString(b)[:16].ToUpper()), nil
}

// generateFeedbackID generates a unique feedback ID
func generateFeedbackID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("FDBK%s", hex.EncodeToString(b)[:16].ToUpper()), nil
}

// generateFeedbackToken generates a unique feedback token
func generateFeedbackToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Helper function for pointer to string
func stringPtr(s string) *string {
	return &s
}
