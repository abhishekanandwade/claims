# Multi-Channel Notification Service Integration

## Overview
This document describes the multi-channel notification service integration for the PLI Claims Processing API. The notification system supports SMS, Email, and WhatsApp notifications for various claim-related events.

## Architecture

### Notification Client
The `NotificationClient` in `repo/postgres/notification_client.go` handles all notification operations. It integrates with external notification service APIs via HTTP/REST.

### Supported Channels
1. **SMS** - SMS Gateway integration for text messages
2. **Email** - Email Service integration for rich HTML emails
3. **WhatsApp** - WhatsApp Business API for template-based messages
4. **Push** - Push notifications (placeholder for future implementation)

## Configuration

### Configuration File (configs/config.yaml)
```yaml
api_clients:
  notification_service:
    enabled: true
    base_url: "https://notification-service.pli.gov.in"
    api_key: ""  # Load from environment variable: NOTIFICATION_API_KEY
    timeout: 15  # Request timeout in seconds
    sms_enabled: true
    email_enabled: true
    whatsapp_enabled: false  # Set to true to enable WhatsApp
```

### Environment Variables
The following environment variables should be set in production:
- `NOTIFICATION_API_KEY` - API key for notification service authentication
- `DB_PASSWORD` - Database password (for security)

## API Integration Details

### 1. SMS Gateway Integration

**Endpoint**: `POST /sms/send`

**Request Format**:
```json
{
  "phone": "+919876543210",
  "message": "Dear Customer, your claim has been registered successfully.",
  "template": "CLAIM_REGISTERED",
  "policy_id": "POL123456",
  "customer_id": "CUST789"
}
```

**Response Format**:
```json
{
  "success": true,
  "message_id": "sms_msg_123456",
  "error": null
}
```

**Error Handling**:
- HTTP status codes: 200 (success), 4xx/5xx (failure)
- Error messages returned in response body
- All errors logged with context

### 2. Email Service Integration

**Endpoint**: `POST /email/send`

**Request Format**:
```json
{
  "to": "customer@example.com",
  "subject": "Claim Registered - CLM2025001",
  "template": "CLAIM_REGISTERED",
  "data": {
    "recipient_name": "John Doe",
    "claim_id": "CLM2025001",
    "custom_message": ""
  }
}
```

**Response Format**:
```json
{
  "success": true,
  "message_id": "email_msg_789012",
  "error": null
}
```

**Email Templates**:
- `CLAIM_REGISTERED` - Claim registration confirmation
- `CLAIM_APPROVED` - Claim approval notification
- `CLAIM_REJECTED` - Claim rejection with reasons
- `DOCUMENT_REQUIRED` - Document completeness reminder
- `PAYMENT_PROCESSED` - Payment disbursement notification
- `MATURITY_INTIMATION` - Maturity claim intimation (batch job)

### 3. WhatsApp Business API Integration

**Endpoint**: `POST /whatsapp/send`

**Request Format**:
```json
{
  "phone": "+919876543210",
  "template": {
    "name": "claim_registered",
    "language": {
      "code": "en"
    }
  },
  "components": [
    {
      "type": "body",
      "parameters": [
        {
          "type": "text",
          "text": "John Doe"
        },
        {
          "type": "text",
          "text": "CLM2025001"
        }
      ]
    }
  ]
}
```

**Response Format**:
```json
{
  "messages": [
    {
      "id": "wamid.HBgLNDE5MTIzNDU2Nzg5FQIAERgSMzg1QTlCNkE2RTlFRTdFNDc="
    }
  ],
  "error": null
}
```

**WhatsApp Templates**:
- `claim_registered` - Claim registration notification
- `claim_approved` - Claim approval notification
- `claim_rejected` - Claim rejection notification
- `document_required` - Document requirement reminder
- `payment_processed` - Payment processed notification
- `maturity_intimation` - Maturity intimation (batch job)

## Notification Types

### Claim Lifecycle Notifications

1. **CLAIM_REGISTERED**
   - Triggered when a new claim is registered
   - Sent to claimant via SMS, Email, WhatsApp
   - Message includes claim number and next steps

2. **CLAIM_APPROVED**
   - Triggered when a claim is approved
   - Sent to claimant via SMS, Email, WhatsApp
   - Message includes approved amount and payment timeline

3. **CLAIM_REJECTED**
   - Triggered when a claim is rejected
   - Sent to claimant via Email (detailed), SMS (brief)
   - Message includes rejection reasons and appeal rights

4. **DOCUMENT_REQUIRED**
   - Triggered when documents are incomplete
   - Sent to claimant via SMS, Email, WhatsApp
   - Message lists required documents

5. **PAYMENT_PROCESSED**
   - Triggered when payment is disbursed
   - Sent to claimant via SMS, Email
   - Message includes payment amount and UTR number

### Batch Notifications

1. **MATURITY_INTIMATION**
   - Triggered by batch job (daily at 9:00 AM)
   - Sent to policyholders whose policies mature in 60 days
   - Multi-channel: SMS, Email, WhatsApp (if enabled)
   - Reference: BR-CLM-MC-002 (60-day advance intimation)

## Business Rules

### BR-CLM-DC-019: Communication Triggers
- Document reminder sent when documents are incomplete
- SLA breach warning sent when claim exceeds 70% of SLA
- Payment notification sent when disbursement is initiated

### Notification Channel Priority
1. **Primary Channel**: Email (detailed information)
2. **Secondary Channel**: SMS (immediate alert)
3. **Optional Channel**: WhatsApp (rich experience)

### Channel Selection Logic
```
IF email_available AND email_enabled:
    Send Email
IF mobile_available AND sms_enabled:
    Send SMS
IF mobile_available AND whatsapp_enabled:
    Send WhatsApp
```

## Implementation Details

### HTTP Client Configuration
```go
httpClient := &http.Client{
    Timeout: time.Duration(config.Timeout) * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false,  // Enforce TLS verification
        },
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

### Authentication
- **Method**: Bearer Token + API Key (dual authentication)
- **Headers**:
  - `Authorization: Bearer <api_key>`
  - `X-API-Key: <api_key>`

### Retry Logic
- Currently no automatic retry (can be added)
- Errors are logged and propagated to caller
- Caller can implement retry logic

### Graceful Degradation
- If one channel fails, others continue
- Partial failures tracked and reported
- `channelsSent` and `channelsFailed` arrays returned

## Usage Examples

### Example 1: Send Single Notification
```go
req := &NotificationRequest{
    NotificationID:   "notif_123",
    NotificationType: "CLAIM_REGISTERED",
    ClaimID:          &claimID,
    RecipientName:    "John Doe",
    RecipientMobile:  "+919876543210",
    RecipientEmail:   "john@example.com",
    Channels:         []string{"SMS", "EMAIL"},
}

channelsSent, channelsFailed, err := notificationClient.SendNotification(ctx, req)
if err != nil {
    log.Error("Failed to send notification", err)
}
```

### Example 2: Send Batch Notifications (Maturity Intimation)
```go
req := &MaturityIntimationRequest{
    PolicyID:       "POL123",
    PolicyNumber:   "123456789",
    CustomerID:     "CUST456",
    CustomerName:   "Jane Doe",
    MaturityDate:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
    MaturityAmount: 500000.00,
    Phone:          "+919876543210",
    Email:          "jane@example.com",
    Channels:       []string{"SMS", "EMAIL", "WHATSAPP"},
}

err := notificationClient.SendMaturityIntimation(ctx, req)
if err != nil {
    log.Error("Failed to send maturity intimation", err)
}
```

### Example 3: Multi-Channel Notification with Custom Message
```go
customMsg := "Your claim requires additional documents. Please submit death certificate."
req := &NotificationRequest{
    NotificationID:   "notif_456",
    NotificationType: "DOCUMENT_REQUIRED",
    ClaimID:          &claimID,
    RecipientName:    "John Doe",
    RecipientMobile:  "+919876543210",
    RecipientEmail:   "john@example.com",
    Channels:         []string{"SMS", "EMAIL", "WHATSAPP"},
    CustomMessage:    &customMsg,
}

channelsSent, channelsFailed, err := notificationClient.SendNotification(ctx, req)
```

## Monitoring and Logging

### Log Levels
- **Info**: Successful notification delivery
- **Warn**: Disabled channels, missing contact info
- **Error**: API failures, timeout errors

### Key Metrics
- Notification delivery rate (success/total)
- Channel-wise delivery statistics
- Average response time
- Error rate by channel

### Monitoring Endpoints
TODO: Implement metrics endpoint for notification statistics
```
GET /notifications/stats
Response: {
  "total_sent": 1000,
  "sms_sent": 800,
  "email_sent": 950,
  "whatsapp_sent": 600,
  "success_rate": 95.5
}
```

## Security Considerations

### API Key Management
- API keys loaded from environment variables
- Never hardcode API keys in source code
- Rotate API keys regularly (90 days recommended)

### Data Privacy
- Phone numbers masked in logs (last 4 digits only)
- Email addresses masked in logs (first 2 characters visible)
- No sensitive claim data in notification messages

### HTTPS/TLS
- All API calls use HTTPS
- TLS certificate verification enforced
- No insecure connections allowed

## Troubleshooting

### Common Issues

1. **SMS Gateway Timeout**
   - Check network connectivity
   - Verify SMS gateway URL is accessible
   - Check API key validity

2. **Email Service Returns 401 Unauthorized**
   - Verify API key is correct
   - Check if API key has expired
   - Ensure API key has email permissions

3. **WhatsApp Template Rejected**
   - Verify template name is correct
   - Check template is approved in WhatsApp Business Manager
   - Ensure parameter count matches template definition

4. **Partial Channel Failure**
   - Check individual channel status in response
   - Verify contact information is valid
   - Check channel-specific configuration

### Debug Mode
Enable debug logging by setting log level to DEBUG:
```go
nlog.SetLevel("DEBUG")
```

## Future Enhancements

1. **Push Notifications**
   - Implement mobile app push notifications
   - FCM (Firebase Cloud Messaging) integration
   - APNS (Apple Push Notification Service) integration

2. **Notification Preferences**
   - Allow customers to opt-in/opt-out of channels
   - Store channel preferences in database
   - Respect Do Not Disturb (DND) registry

3. **Advanced Features**
   - Notification scheduling
   - Retry logic with exponential backoff
   - Webhook notifications for third-party integrations
   - Multi-language support

4. **Analytics**
   - Notification delivery dashboard
   - Channel performance metrics
   - Customer engagement tracking

## References

- **Business Rules**: BR-CLM-DC-019 (Communication triggers)
- **Integrations**: INT-CLM-009 (Email), INT-CLM-010 (SMS), INT-CLM-011 (WhatsApp)
- **Template**: seed/template/template.md
- **Swagger**: seed/swagger/notification.yaml

## Contact

For questions or issues related to notification service integration:
- Development Team: dev-team@pli.gov.in
- Notification Service Provider: support@notification-service.pli.gov.in
