# Maturity Intimation Batch Job

## Overview

The maturity intimation batch job automatically sends notifications to policyholders whose policies are maturing in the next 60 days. This ensures timely intimation per IRDAI guidelines and helps policyholders initiate the claim process early.

**Reference**: FR-CLM-MC-002, BR-CLM-MC-002

## Architecture

```
┌─────────────────┐
│  Cron Scheduler │ (Daily at 9:00 AM)
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│ MaturityIntimationJob   │
└────────┬────────────────┘
         │
         ├──► Query Policies (Policy Service)
         │    └──► Policies maturing in 60-90 days
         │    └──► Filter: No claim registered
         │    └──► Filter: No intimation sent
         │
         ├──► Send Notifications (Notification Service)
         │    ├──► SMS
         │    ├──► Email
         │    └──► WhatsApp (optional)
         │
         └──► Record Audit Trail (Claim History)
              └──► Event: MATURITY_INTIMATION_SENT
```

## Configuration

### Config File (`configs/config.yaml`)

```yaml
batch:
  maturity_intimation:
    enabled: true                           # Enable/disable batch job
    schedule: "0 0 9 * * *"                 # Cron schedule (daily at 9:00 AM)
    days_in_advance: 60                      # Days before maturity to send intimation
    batch_size: 100                          # Number of policies to process per batch
    notification_channels:                   # Notification channels to use
      - "SMS"
      - "EMAIL"
      - "WHATSAPP"
```

### Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enabled` | boolean | `true` | Enable/disable the batch job |
| `schedule` | string | `"0 0 9 * * *"` | Cron schedule (second minute hour day month weekday) |
| `days_in_advance` | int | `60` | Days before maturity to send intimation (BR-CLM-MC-002) |
| `batch_size` | int | `100` | Number of policies to process per batch |
| `notification_channels` | []string | `["SMS", "EMAIL"]` | Channels to send notifications through |

## Business Rules

### BR-CLM-MC-002: 60-Day Advance Intimation

- **Rule**: Send intimation 60 days before maturity date
- **Rationale**: Allows sufficient time for policyholders to gather documents and submit claim
- **Implementation**: Batch job runs daily, queries policies with maturity_date in (today + 60 to today + 90) days

### Duplicate Prevention

- **Rule**: Do not send multiple intimations for the same policy
- **Implementation**: Check `claim_history` table for `MATURITY_INTIMATION_SENT` event before sending
- **Audit Trail**: Record every intimation in `claim_history` table with:
  - `event_type`: `MATURITY_INTIMATION_SENT`
  - `event_category`: `NOTIFICATION`
  - `performed_by`: `SYSTEM`
  - `performed_by_role`: `BATCH_JOB`

## Running the Batch Job

### Option 1: Standalone Executable

```bash
# Run manually
go run cmd/batch_runner/main.go

# Run with custom config
CONFIG_FILE=configs/config.prod.yaml go run cmd/batch_runner/main.go

# Build and run
go build -o bin/batch_runner cmd/batch_runner/main.go
./bin/batch_runner
```

### Option 2: Linux Cron

Add to crontab (`crontab -e`):

```cron
# Run daily at 9:00 AM
0 9 * * * cd /path/to/claims-api && ./bin/batch_runner >> /var/log/claims-api/batch.log 2>&1
```

### Option 3: Kubernetes CronJob

Create `k8s/cronjob-maturity-intimation.yaml`:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: maturity-intimation-job
spec:
  schedule: "0 9 * * *"  # Daily at 9:00 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: batch-runner
            image: claims-api:latest
            command: ["/bin/batch_runner"]
            env:
            - name: CONFIG_FILE
              value: "/config/config.prod.yaml"
            volumeMounts:
            - name: config
              mountPath: /config
          volumes:
          - name: config
            configMap:
              name: claims-api-config
          restartPolicy: OnFailure
```

Deploy:

```bash
kubectl apply -f k8s/cronjob-maturity-intimation.yaml
```

### Option 4: Temporal Workflow

**TODO**: Implement Temporal workflow for scheduled tasks (Phase 8)

## Monitoring

### Logs

The batch job logs all activities:

```json
{
  "level": "INFO",
  "message": "Starting maturity intimation batch job",
  "job": "maturity_intimation",
  "days_in_advance": 60,
  "batch_size": 100
}

{
  "level": "INFO",
  "message": "Found policies due for maturity",
  "count": 150
}

{
  "level": "INFO",
  "message": "Maturity intimation sent successfully",
  "policy_id": "123e4567-e89b-12d3-a456-426614174000",
  "policy_number": "POL123456",
  "customer_id": "CUST001",
  "channels": ["SMS", "EMAIL"]
}

{
  "level": "INFO",
  "message": "Maturity intimation batch job completed",
  "duration_seconds": 45.23,
  "total_policies": 150,
  "success_count": 148,
  "failure_count": 2
}
```

### Metrics

Monitor the following metrics:

| Metric | Description | Target |
|--------|-------------|--------|
| `batch_job_duration_seconds` | Total execution time | < 30 minutes |
| `batch_job_policies_processed` | Total policies processed | Varies |
| `batch_job_success_count` | Successful intimations | > 95% |
| `batch_job_failure_count` | Failed intimations | < 5% |

### Alerts

Set up alerts for:

1. **Job Failure**: Batch job exits with error
2. **High Failure Rate**: > 5% intimations fail
3. **Long Running Time**: Job takes > 30 minutes
4. **No Policies Processed**: Job completes but processes 0 policies

## Notification Content

### SMS Template

```
Dear {customer_name}, your policy {policy_number} is maturing on {maturity_date}. Maturity amount: {maturity_amount}. Please submit documents to initiate claim process. Regards, PLI
```

### Email Template

**Subject**: Maturity Intimation - Policy {policy_number}

**Body**:

```
Dear {customer_name},

We are pleased to inform you that your policy is approaching maturity.

Policy Details:
- Policy Number: {policy_number}
- Maturity Date: {maturity_date}
- Maturity Amount: {maturity_amount}

Documents Required:
1. Original Policy Bond
2. NEFT Mandate Form (duly filled)
3. Cancelled Cheque (for NEFT credit)
4. Identity Proof (Aadhaar/PAN)
5. Discharge Receipt

Please submit these documents at your nearest PLI office or via our online portal to initiate the claim process.

Note: As per IRDAI guidelines, maturity claims are processed within 7 days of receipt of complete documents.

If you have any queries, please contact our customer support.

Regards,
Postal Life Insurance
```

### WhatsApp Template

```
Dear {customer_name}, your PLI policy {policy_number} is maturing on {maturity_date}. Amount: {maturity_amount}. Submit documents: Policy bond, NEFT form, Cancelled cheque, ID proof. Claim processed within 7 days. Contact: 1800-XXX-XXXX
```

## Database Schema

### claim_history Table

Audit trail entries created by batch job:

```sql
INSERT INTO claim_history (
    id,
    policy_id,
    claim_id,
    event_type,
    event_category,
    description,
    old_value,
    new_value,
    performed_by,
    performed_by_role,
    ip_address,
    user_agent,
    created_at
) VALUES (
    'uuid-here',
    'policy-uuid',
    NULL,  -- No claim yet
    'MATURITY_INTIMATION_SENT',
    'NOTIFICATION',
    'Maturity intimation sent via SMS, EMAIL',
    NULL,
    '{"channels": ["SMS", "EMAIL"]}',
    'SYSTEM',
    'BATCH_JOB',
    NULL,
    NULL,
    NOW()
);
```

## Integration Points

### 1. Policy Service Integration

**TODO**: Implement in production

The batch job needs to query policies from Policy Service:

```http
GET /policies/maturity-due
Query Parameters:
  - start_date: string (ISO 8601 date)
  - end_date: string (ISO 8601 date)
  - page: int (default: 1)
  - limit: int (default: 100)

Response:
{
  "data": [
    {
      "policy_id": "uuid",
      "policy_number": "POL123456",
      "customer_id": "uuid",
      "customer_name": "John Doe",
      "maturity_date": "2025-03-15",
      "maturity_amount": 250000.00,
      "phone": "+91-9876543210",
      "email": "john.doe@example.com"
    }
  ],
  "pagination": {
    "total": 150,
    "page": 1,
    "limit": 100
  }
}
```

### 2. Notification Service Integration

**TODO**: Implement in production

#### SMS Gateway

```http
POST /sms/send
{
  "phone": "+91-9876543210",
  "message": "Dear John Doe, your policy POL123456 is maturing on 2025-03-15...",
  "template_id": "MATURITY_INTIMATION"
}
```

#### Email Service

```http
POST /email/send
{
  "to": "john.doe@example.com",
  "subject": "Maturity Intimation - Policy POL123456",
  "template": "MATURITY_INTIMATION",
  "data": {
    "customer_name": "John Doe",
    "policy_number": "POL123456",
    "maturity_date": "2025-03-15",
    "maturity_amount": 250000.00
  }
}
```

#### WhatsApp Business API

```http
POST /whatsapp/send
{
  "phone": "+91-9876543210",
  "template": "maturity_intimation",
  "parameters": {
    "customer_name": "John Doe",
    "policy_number": "POL123456",
    "maturity_date": "2025-03-15",
    "maturity_amount": "250000.00"
  }
}
```

## Testing

### Manual Testing

```bash
# Run batch job manually
go run cmd/batch_runner/main.go

# Check logs
tail -f /var/log/claims-api/batch.log

# Verify audit trail
psql -h localhost -U postgres -d claims_db -c "
  SELECT
    policy_id,
    event_type,
    description,
    created_at
  FROM claim_history
  WHERE event_type = 'MATURITY_INTIMATION_SENT'
  ORDER BY created_at DESC
  LIMIT 10;
"
```

### Integration Testing

```bash
# Create test policies due for maturity in 60-90 days
# Run batch job
# Verify notifications sent
# Check claim_history table
# Verify no duplicate intimations
```

## Troubleshooting

### Issue: Batch job processes 0 policies

**Possible Causes**:
1. No policies due for maturity in date range
2. Policy Service integration not working
3. All policies already have intimations sent

**Resolution**:
```sql
-- Check if intimations already sent
SELECT
  policy_id,
  COUNT(*) as intimation_count
FROM claim_history
WHERE event_type = 'MATURATION_INTIMATION_SENT'
GROUP BY policy_id;

-- Check if claims already registered
SELECT
  policy_id,
  claim_number,
  status
FROM claims
WHERE claim_type = 'MATURITY'
  AND policy_id IN (...);
```

### Issue: High failure rate (>5%)

**Possible Causes**:
1. Invalid phone numbers/emails
2. SMS/Email gateway down
3. Rate limiting by gateway

**Resolution**:
1. Check gateway logs
2. Validate contact details in Policy Service
3. Implement retry logic with exponential backoff
4. Configure rate limits in gateway

### Issue: Duplicate intimations sent

**Possible Causes**:
1. Audit trail check not working
2. Job run multiple times concurrently

**Resolution**:
1. Implement distributed lock (Redis)
2. Add unique constraint on (policy_id, event_type, event_date) in claim_history
3. Ensure only one job instance runs at a time

## Future Enhancements

1. **Retry Logic**: Implement exponential backoff for failed notifications
2. **Multi-language Support**: Send notifications in regional languages
3. **Preference Management**: Allow customers to choose notification channels
4. **Real-time Dashboard**: Monitor batch job execution in real-time
5. **A/B Testing**: Test different message templates for better engagement
6. **Escalation**: Alert support team if intimation fails multiple times
7. **Customer Portal**: Display intimation status in customer portal

## References

- FR-CLM-MC-002: Batch intimation for maturity claims
- BR-CLM-MC-002: 60-day advance intimation requirement
- BR-CLM-DC-019: Communication triggers
- seed/analysis/maturity_claim_analysis.md
- seed/swagger/maturity_claim.yaml
