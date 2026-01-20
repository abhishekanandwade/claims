package repo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	nlog "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	"gitlab.cept.gov.in/pli/claims-api/configs"
)

// NotificationClient handles sending notifications via multiple channels
// Integrates with:
// 1. SMS Gateway (for SMS notifications)
// 2. Email Service (for email notifications)
// 3. WhatsApp Business API (for WhatsApp notifications)
//
// Reference: BR-CLM-DC-019 (Communication triggers)
// Reference: INT-CLM-009 (Email Service), INT-CLM-011 (WhatsApp)
type NotificationClient struct {
	config            *configs.Config
	httpClient        *http.Client
	smsGatewayURL     string
	emailServiceURL   string
	whatsappAPIURL    string
	authenticationKey string
	smsEnabled        bool
	emailEnabled      bool
	whatsappEnabled   bool
}

// NewNotificationClient creates a new notification client
func NewNotificationClient(config *configs.Config) *NotificationClient {
	// Create HTTP client with timeout and TLS config
	httpClient := &http.Client{
		Timeout: time.Duration(config.APIClients.NotificationService.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Set to true only for testing
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &NotificationClient{
		config:            config,
		httpClient:        httpClient,
		smsGatewayURL:     config.APIClients.NotificationService.BaseURL + "/sms/send",
		emailServiceURL:   config.APIClients.NotificationService.BaseURL + "/email/send",
		whatsappAPIURL:    config.APIClients.NotificationService.BaseURL + "/whatsapp/send",
		authenticationKey: config.APIClients.NotificationService.APIKey,
		smsEnabled:        config.APIClients.NotificationService.SMSEnabled,
		emailEnabled:      config.APIClients.NotificationService.EmailEnabled,
		whatsappEnabled:   config.APIClients.NotificationService.WhatsappEnabled,
	}
}

// SendMaturityIntimation sends maturity intimation notification to policyholder
// This implements the NotificationService interface used by batch jobs
// Reference: FR-CLM-MC-002, BR-CLM-MC-002
func (c *NotificationClient) SendMaturityIntimation(ctx context.Context, req *MaturityIntimationRequest) error {
	nlog.Info(ctx, "Sending maturity intimation", map[string]interface{}{
		"policy_id":       req.PolicyID,
		"policy_number":   req.PolicyNumber,
		"customer_id":     req.CustomerID,
		"customer_name":   req.CustomerName,
		"maturity_date":   req.MaturityDate.Format("2006-01-02"),
		"maturity_amount": req.MaturityAmount,
		"channels":        req.Channels,
	})

	// Send notifications via each configured channel
	for _, channel := range req.Channels {
		switch channel {
		case "SMS":
			if err := c.sendSMS(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send SMS", map[string]interface{}{
					"policy_id": req.PolicyID,
					"error":     err.Error(),
				})
				return fmt.Errorf("SMS send failed: %w", err)
			}

		case "EMAIL":
			if err := c.sendEmail(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send Email", map[string]interface{}{
					"policy_id": req.PolicyID,
					"error":     err.Error(),
				})
				return fmt.Errorf("Email send failed: %w", err)
			}

		case "WHATSAPP":
			if err := c.sendWhatsApp(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send WhatsApp", map[string]interface{}{
					"policy_id": req.PolicyID,
					"error":     err.Error(),
				})
				return fmt.Errorf("WhatsApp send failed: %w", err)
			}

		default:
			nlog.Warn(ctx, "Unknown notification channel", map[string]interface{}{
				"channel":   channel,
				"policy_id": req.PolicyID,
			})
		}
	}

	return nil
}

// sendSMS sends SMS notification
// Integrates with SMS Gateway API
// Reference: BR-CLM-DC-019, INT-CLM-010 (SMS Gateway)
func (c *NotificationClient) sendSMS(ctx context.Context, req *MaturityIntimationRequest) error {
	if !c.smsEnabled {
		nlog.Warn(ctx, "SMS notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending SMS", map[string]interface{}{
		"phone":     req.Phone,
		"policy_id": req.PolicyID,
	})

	// Prepare SMS message
	message := fmt.Sprintf(
		"Dear %s, your policy %s is maturing on %s. Maturity amount: %.2f. Please submit documents to initiate claim process.",
		req.CustomerName,
		req.PolicyNumber,
		req.MaturityDate.Format("2006-01-02"),
		req.MaturityAmount,
	)

	// Prepare SMS Gateway API request
	smsRequest := map[string]interface{}{
		"phone":      req.Phone,
		"message":    message,
		"template":   "MATURITY_INTIMATION",
		"policy_id":  req.PolicyID,
		"customer_id": req.CustomerID,
	}

	// Marshal request body
	requestBody, err := json.Marshal(smsRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.smsGatewayURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)
	httpReq.Header.Set("X-API-Key", c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("SMS gateway API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		nlog.Error(ctx, "SMS gateway returned non-OK status", map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
			"phone":       req.Phone,
		})
		return fmt.Errorf("SMS gateway returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var smsResponse struct {
		Success   bool   `json:"success"`
		MessageID string `json:"message_id"`
		Error     string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &smsResponse); err != nil {
		return fmt.Errorf("failed to parse SMS response: %w", err)
	}

	if !smsResponse.Success {
		return fmt.Errorf("SMS sending failed: %s", smsResponse.Error)
	}

	nlog.Info(ctx, "SMS sent successfully", map[string]interface{}{
		"phone":       req.Phone,
		"message_id":  smsResponse.MessageID,
		"policy_id":   req.PolicyID,
	})

	return nil
}

// sendEmail sends email notification
// Integrates with Email Service API
// Reference: BR-CLM-DC-019, INT-CLM-009 (Email Service)
func (c *NotificationClient) sendEmail(ctx context.Context, req *MaturityIntimationRequest) error {
	if !c.emailEnabled {
		nlog.Warn(ctx, "Email notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending Email", map[string]interface{}{
		"email":     req.Email,
		"policy_id": req.PolicyID,
	})

	// Prepare email data
	emailRequest := map[string]interface{}{
		"to":       req.Email,
		"subject":  fmt.Sprintf("Maturity Intimation - Policy %s", req.PolicyNumber),
		"template": "MATURITY_INTIMATION",
		"data": map[string]interface{}{
			"customer_name":       req.CustomerName,
			"policy_number":       req.PolicyNumber,
			"maturity_date":       req.MaturityDate.Format("2006-01-02"),
			"maturity_amount":     fmt.Sprintf("%.2f", req.MaturityAmount),
			"documents_required":  []string{"Policy bond", "NEFT form", "Cancelled cheque", "ID proof"},
			"claim_submission_url": "https://pli.gov.in/claims/submit",
		},
	}

	// Marshal request body
	requestBody, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.emailServiceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)
	httpReq.Header.Set("X-API-Key", c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("email service API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read email response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		nlog.Error(ctx, "Email service returned non-OK status", map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
			"email":       req.Email,
		})
		return fmt.Errorf("email service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var emailResponse struct {
		Success    bool   `json:"success"`
		MessageID  string `json:"message_id"`
		Error      string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &emailResponse); err != nil {
		return fmt.Errorf("failed to parse email response: %w", err)
	}

	if !emailResponse.Success {
		return fmt.Errorf("email sending failed: %s", emailResponse.Error)
	}

	nlog.Info(ctx, "Email sent successfully", map[string]interface{}{
		"email":      req.Email,
		"message_id": emailResponse.MessageID,
		"policy_id":  req.PolicyID,
	})

	return nil
}

// sendWhatsApp sends WhatsApp notification
// Integrates with WhatsApp Business API
// Reference: BR-CLM-DC-019, INT-CLM-011 (WhatsApp Business API)
func (c *NotificationClient) sendWhatsApp(ctx context.Context, req *MaturityIntimationRequest) error {
	if !c.whatsappEnabled {
		nlog.Warn(ctx, "WhatsApp notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending WhatsApp", map[string]interface{}{
		"phone":     req.Phone,
		"policy_id": req.PolicyID,
	})

	// Prepare WhatsApp Business API request
	whatsappRequest := map[string]interface{}{
		"phone": req.Phone,
		"template": map[string]interface{}{
			"name": "maturity_intimation",
			"language": map[string]interface{}{
				"code": "en",
			},
		},
		"components": []map[string]interface{}{
			{
				"type": "body",
				"parameters": []map[string]interface{}{
					{
						"type": "text",
						"text": req.CustomerName,
					},
					{
						"type": "text",
						"text": req.PolicyNumber,
					},
					{
						"type": "text",
						"text": req.MaturityDate.Format("2006-01-02"),
					},
					{
						"type": "text",
						"text": fmt.Sprintf("%.2f", req.MaturityAmount),
					},
				},
			},
		},
	}

	// Marshal request body
	requestBody, err := json.Marshal(whatsappRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal WhatsApp request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.whatsappAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("WhatsApp API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read WhatsApp response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		nlog.Error(ctx, "WhatsApp API returned non-OK status", map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
			"phone":       req.Phone,
		})
		return fmt.Errorf("WhatsApp API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var whatsappResponse struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &whatsappResponse); err != nil {
		return fmt.Errorf("failed to parse WhatsApp response: %w", err)
	}

	if whatsappResponse.Error != nil {
		return fmt.Errorf("WhatsApp sending failed: %s", whatsappResponse.Error.Message)
	}

	messageID := ""
	if len(whatsappResponse.Messages) > 0 {
		messageID = whatsappResponse.Messages[0].ID
	}

	nlog.Info(ctx, "WhatsApp message sent successfully", map[string]interface{}{
		"phone":       req.Phone,
		"message_id":  messageID,
		"policy_id":   req.PolicyID,
	})

	return nil
}

// SendNotification sends a generic notification via multiple channels
// This is a generic method used by NotificationHandler
// Reference: BR-CLM-DC-019 (Communication triggers)
func (c *NotificationClient) SendNotification(ctx context.Context, req *NotificationRequest) ([]string, []string, error) {
	nlog.Info(ctx, "Sending notification", map[string]interface{}{
		"notification_id":   req.NotificationID,
		"notification_type": req.NotificationType,
		"claim_id":          req.ClaimID,
		"recipient_name":    req.RecipientName,
		"channels":          req.Channels,
	})

	channelsSent := []string{}
	channelsFailed := []string{}

	// Send notifications via each configured channel
	for _, channel := range req.Channels {
		switch channel {
		case "SMS":
			if req.RecipientMobile == "" {
				nlog.Warn(ctx, "Mobile number not provided, skipping SMS", map[string]interface{}{
					"notification_id": req.NotificationID,
				})
				channelsFailed = append(channelsFailed, "SMS")
				continue
			}

			if err := c.sendGenericSMS(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send SMS", map[string]interface{}{
					"notification_id": req.NotificationID,
					"error":           err.Error(),
				})
				channelsFailed = append(channelsFailed, "SMS")
			} else {
				channelsSent = append(channelsSent, "SMS")
			}

		case "EMAIL":
			if req.RecipientEmail == "" {
				nlog.Warn(ctx, "Email not provided, skipping Email", map[string]interface{}{
					"notification_id": req.NotificationID,
				})
				channelsFailed = append(channelsFailed, "EMAIL")
				continue
			}

			if err := c.sendGenericEmail(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send Email", map[string]interface{}{
					"notification_id": req.NotificationID,
					"error":           err.Error(),
				})
				channelsFailed = append(channelsFailed, "EMAIL")
			} else {
				channelsSent = append(channelsSent, "EMAIL")
			}

		case "WHATSAPP":
			if req.RecipientMobile == "" {
				nlog.Warn(ctx, "Mobile number not provided, skipping WhatsApp", map[string]interface{}{
					"notification_id": req.NotificationID,
				})
				channelsFailed = append(channelsFailed, "WHATSAPP")
				continue
			}

			if err := c.sendGenericWhatsApp(ctx, req); err != nil {
				nlog.Error(ctx, "Failed to send WhatsApp", map[string]interface{}{
					"notification_id": req.NotificationID,
					"error":           err.Error(),
				})
				channelsFailed = append(channelsFailed, "WHATSAPP")
			} else {
				channelsSent = append(channelsSent, "WHATSAPP")
			}

		case "PUSH":
			// TODO: Implement push notifications
			nlog.Warn(ctx, "Push notifications not implemented yet", map[string]interface{}{
				"notification_id": req.NotificationID,
			})
			channelsFailed = append(channelsFailed, "PUSH")

		default:
			nlog.Warn(ctx, "Unknown notification channel", map[string]interface{}{
				"channel":         channel,
				"notification_id": req.NotificationID,
			})
			channelsFailed = append(channelsFailed, channel)
		}
	}

	return channelsSent, channelsFailed, nil
}

// sendGenericSMS sends generic SMS notification
// Integrates with SMS Gateway API
// Reference: BR-CLM-DC-019, INT-CLM-010 (SMS Gateway)
func (c *NotificationClient) sendGenericSMS(ctx context.Context, req *NotificationRequest) error {
	if !c.smsEnabled {
		nlog.Warn(ctx, "SMS notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending generic SMS", map[string]interface{}{
		"notification_id": req.NotificationID,
		"phone":           req.RecipientMobile,
	})

	// Prepare message
	message := ""
	if req.CustomMessage != nil && *req.CustomMessage != "" {
		message = *req.CustomMessage
	} else {
		// Default message based on notification type
		switch req.NotificationType {
		case "CLAIM_REGISTERED":
			message = fmt.Sprintf("Dear %s, your claim %s has been registered successfully. We will process it shortly.", req.RecipientName, getStringValue(req.ClaimID, ""))
		case "CLAIM_APPROVED":
			message = fmt.Sprintf("Dear %s, your claim %s has been approved. Payment will be processed soon.", req.RecipientName, getStringValue(req.ClaimID, ""))
		case "CLAIM_REJECTED":
			message = fmt.Sprintf("Dear %s, your claim %s has been rejected. Please check your email for details.", req.RecipientName, getStringValue(req.ClaimID, ""))
		case "DOCUMENT_REQUIRED":
			message = fmt.Sprintf("Dear %s, additional documents are required for your claim %s. Please check your email.", req.RecipientName, getStringValue(req.ClaimID, ""))
		case "PAYMENT_PROCESSED":
			message = fmt.Sprintf("Dear %s, payment for your claim %s has been processed. Amount will be credited shortly.", req.RecipientName, getStringValue(req.ClaimID, ""))
		default:
			message = fmt.Sprintf("Dear %s, there is an update on your claim %s. Please check your email for details.", req.RecipientName, getStringValue(req.ClaimID, ""))
		}
	}

	// Prepare SMS Gateway API request
	smsRequest := map[string]interface{}{
		"phone":      req.RecipientMobile,
		"message":    message,
		"template":   req.NotificationType,
		"metadata": map[string]interface{}{
			"notification_id": req.NotificationID,
			"claim_id":        getStringValue(req.ClaimID, ""),
		},
	}

	// Marshal request body
	requestBody, err := json.Marshal(smsRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.smsGatewayURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)
	httpReq.Header.Set("X-API-Key", c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("SMS gateway API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		nlog.Error(ctx, "SMS gateway returned non-OK status", map[string]interface{}{
			"status_code":     resp.StatusCode,
			"response":        string(respBody),
			"phone":           req.RecipientMobile,
			"notification_id": req.NotificationID,
		})
		return fmt.Errorf("SMS gateway returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var smsResponse struct {
		Success   bool   `json:"success"`
		MessageID string `json:"message_id"`
		Error     string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &smsResponse); err != nil {
		return fmt.Errorf("failed to parse SMS response: %w", err)
	}

	if !smsResponse.Success {
		return fmt.Errorf("SMS sending failed: %s", smsResponse.Error)
	}

	nlog.Info(ctx, "SMS sent successfully", map[string]interface{}{
		"phone":           req.RecipientMobile,
		"message_id":      smsResponse.MessageID,
		"notification_id": req.NotificationID,
	})

	return nil
}

// sendGenericEmail sends generic email notification
// Integrates with Email Service API
// Reference: BR-CLM-DC-019, INT-CLM-009 (Email Service)
func (c *NotificationClient) sendGenericEmail(ctx context.Context, req *NotificationRequest) error {
	if !c.emailEnabled {
		nlog.Warn(ctx, "Email notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending generic Email", map[string]interface{}{
		"notification_id": req.NotificationID,
		"email":           req.RecipientEmail,
	})

	// Prepare email subject
	subject := ""
	switch req.NotificationType {
	case "CLAIM_REGISTERED":
		subject = fmt.Sprintf("Claim Registered - %s", getStringValue(req.ClaimID, ""))
	case "CLAIM_APPROVED":
		subject = fmt.Sprintf("Claim Approved - %s", getStringValue(req.ClaimID, ""))
	case "CLAIM_REJECTED":
		subject = fmt.Sprintf("Claim Rejected - %s", getStringValue(req.ClaimID, ""))
	case "DOCUMENT_REQUIRED":
		subject = fmt.Sprintf("Documents Required - %s", getStringValue(req.ClaimID, ""))
	case "PAYMENT_PROCESSED":
		subject = fmt.Sprintf("Payment Processed - %s", getStringValue(req.ClaimID, ""))
	default:
		subject = fmt.Sprintf("Update on your Claim - %s", getStringValue(req.ClaimID, ""))
	}

	// Prepare email data
	emailRequest := map[string]interface{}{
		"to":       req.RecipientEmail,
		"subject":  subject,
		"template": req.NotificationType,
		"data": map[string]interface{}{
			"recipient_name": req.RecipientName,
			"claim_id":       getStringValue(req.ClaimID, ""),
			"custom_message": getStringValue(req.CustomMessage, ""),
		},
	}

	// Marshal request body
	requestBody, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.emailServiceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)
	httpReq.Header.Set("X-API-Key", c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("email service API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read email response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		nlog.Error(ctx, "Email service returned non-OK status", map[string]interface{}{
			"status_code":     resp.StatusCode,
			"response":        string(respBody),
			"email":           req.RecipientEmail,
			"notification_id": req.NotificationID,
		})
		return fmt.Errorf("email service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var emailResponse struct {
		Success   bool   `json:"success"`
		MessageID string `json:"message_id"`
		Error     string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &emailResponse); err != nil {
		return fmt.Errorf("failed to parse email response: %w", err)
	}

	if !emailResponse.Success {
		return fmt.Errorf("email sending failed: %s", emailResponse.Error)
	}

	nlog.Info(ctx, "Email sent successfully", map[string]interface{}{
		"email":           req.RecipientEmail,
		"message_id":      emailResponse.MessageID,
		"notification_id": req.NotificationID,
	})

	return nil
}

// sendGenericWhatsApp sends generic WhatsApp notification
// Integrates with WhatsApp Business API
// Reference: BR-CLM-DC-019, INT-CLM-011 (WhatsApp Business API)
func (c *NotificationClient) sendGenericWhatsApp(ctx context.Context, req *NotificationRequest) error {
	if !c.whatsappEnabled {
		nlog.Warn(ctx, "WhatsApp notifications are disabled", nil)
		return nil
	}

	nlog.Info(ctx, "Sending generic WhatsApp", map[string]interface{}{
		"notification_id": req.NotificationID,
		"phone":           req.RecipientMobile,
	})

	// Determine template based on notification type
	template := "claim_update"
	switch req.NotificationType {
	case "CLAIM_REGISTERED":
		template = "claim_registered"
	case "CLAIM_APPROVED":
		template = "claim_approved"
	case "CLAIM_REJECTED":
		template = "claim_rejected"
	case "DOCUMENT_REQUIRED":
		template = "document_required"
	case "PAYMENT_PROCESSED":
		template = "payment_processed"
	}

	// Prepare WhatsApp Business API request
	whatsappRequest := map[string]interface{}{
		"phone": req.RecipientMobile,
		"template": map[string]interface{}{
			"name": template,
			"language": map[string]interface{}{
				"code": "en",
			},
		},
		"components": []map[string]interface{}{
			{
				"type": "body",
				"parameters": []map[string]interface{}{
					{
						"type": "text",
						"text": req.RecipientName,
					},
					{
						"type": "text",
						"text": getStringValue(req.ClaimID, ""),
					},
				},
			},
		},
	}

	// Marshal request body
	requestBody, err := json.Marshal(whatsappRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal WhatsApp request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.whatsappAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authenticationKey)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("WhatsApp API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read WhatsApp response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		nlog.Error(ctx, "WhatsApp API returned non-OK status", map[string]interface{}{
			"status_code":     resp.StatusCode,
			"response":        string(respBody),
			"phone":           req.RecipientMobile,
			"notification_id": req.NotificationID,
		})
		return fmt.Errorf("WhatsApp API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var whatsappResponse struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &whatsappResponse); err != nil {
		return fmt.Errorf("failed to parse WhatsApp response: %w", err)
	}

	if whatsappResponse.Error != nil {
		return fmt.Errorf("WhatsApp sending failed: %s", whatsappResponse.Error.Message)
	}

	messageID := ""
	if len(whatsappResponse.Messages) > 0 {
		messageID = whatsappResponse.Messages[0].ID
	}

	nlog.Info(ctx, "WhatsApp message sent successfully", map[string]interface{}{
		"phone":           req.RecipientMobile,
		"message_id":      messageID,
		"notification_id": req.NotificationID,
	})

	return nil
}

// MaturityIntimationRequest represents a request for maturity intimation
// This is duplicated here to avoid import cycle
type MaturityIntimationRequest struct {
	PolicyID       string    `json:"policy_id"`
	PolicyNumber   string    `json:"policy_number"`
	CustomerID     string    `json:"customer_id"`
	CustomerName   string    `json:"customer_name"`
	MaturityDate   time.Time `json:"maturity_date"`
	MaturityAmount float64   `json:"maturity_amount"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Channels       []string  `json:"channels"`
}

// NotificationRequest represents a generic notification request
type NotificationRequest struct {
	NotificationID   string   `json:"notification_id"`
	NotificationType string   `json:"notification_type"`
	ClaimID          *string  `json:"claim_id,omitempty"`
	RecipientName    string   `json:"recipient_name"`
	RecipientMobile  string   `json:"recipient_mobile,omitempty"`
	RecipientEmail   string   `json:"recipient_email,omitempty"`
	Channels         []string `json:"channels"`
	CustomMessage    *string  `json:"custom_message,omitempty"`
}

// Helper function to safely get string value from pointer
func getStringValue(ptr *string, defaultVal string) string {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
