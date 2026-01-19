package handler

import (
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	*serverHandler.Base
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	base := serverHandler.New("Notifications").
		SetPrefix("/v1").
		AddPrefix("")
	return &NotificationHandler{
		Base: base,
	}
}

// TODO: Implement notification handler routes and methods
// This is a placeholder - will be implemented in Phase 7
