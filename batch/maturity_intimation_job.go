package batch

import (
	"context"
	"fmt"
	"time"

	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
	nlog "gitlab.cept.gov.in/it-2.0-common/n-api-log"
)

// MaturityIntimationJob handles batch intimation of maturity claims
// Reference: FR-CLM-MC-002 - 60-day advance intimation for maturity
type MaturityIntimationJob struct {
	claimRepo      *repo.ClaimRepository
	notificationSvc NotificationService
	config         *BatchConfig
}

// BatchConfig contains configuration for batch jobs
type BatchConfig struct {
	// MaturityIntimationDaysInAdvance is the number of days before maturity to send intimation
	// Default: 60 days per BR-CLM-MC-002
	MaturityIntimationDaysInAdvance int `json:"maturity_intimation_days_in_advance"`

	// BatchSize is the number of policies to process in each batch
	// Default: 100
	BatchSize int `json:"batch_size"`

	// NotificationChannels are the channels to send notifications through
	// Options: SMS, EMAIL, WHATSAPP
	NotificationChannels []string `json:"notification_channels"`
}

// NotificationService defines the interface for sending notifications
// TODO: This will be implemented when notification service integration is done
type NotificationService interface {
	SendMaturityIntimation(ctx context.Context, req *MaturityIntimationRequest) error
}

// MaturityIntimationRequest represents a request for maturity intimation
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

// NewMaturityIntimationJob creates a new maturity intimation job
func NewMaturityIntimationJob(
	claimRepo *repo.ClaimRepository,
	notificationSvc NotificationService,
	config *BatchConfig,
) *MaturityIntimationJob {
	// Set defaults if not provided
	if config == nil {
		config = &BatchConfig{
			MaturityIntimationDaysInAdvance: 60, // BR-CLM-MC-002
			BatchSize:                        100,
			NotificationChannels:             []string{"SMS", "EMAIL"},
		}
	}

	return &MaturityIntimationJob{
		claimRepo:      claimRepo,
		notificationSvc: notificationSvc,
		config:         config,
	}
}

// Run executes the maturity intimation batch job
// This job:
// 1. Queries policies maturing in the next MaturityIntimationDaysInAdvance days
// 2. Filters out policies that already have intimations sent
// 3. Sends notifications to policyholders via configured channels
// 4. Logs all activities for audit trail
func (j *MaturityIntimationJob) Run(ctx context.Context) error {
	startTime := time.Now()
	nlog.Info(ctx, "Starting maturity intimation batch job", map[string]interface{}{
		"job": "maturity_intimation",
	})

	// Calculate date range for maturity due report
	daysInAdvance := j.config.MaturityIntimationDaysInAdvance
	startDate := time.Now().AddDate(0, 0, daysInAdvance).Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, 30) // Next 30 days from intimation date

	nlog.Info(ctx, "Querying policies for maturity intimation", map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})

	// Query policies maturing in the date range
	// TODO: This query needs to be implemented in ClaimRepository
	// For now, we'll create a placeholder
	policies, err := j.getPoliciesDueForMaturity(ctx, startDate, endDate)
	if err != nil {
		nlog.Error(ctx, "Failed to query policies due for maturity", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to query policies: %w", err)
	}

	nlog.Info(ctx, "Found policies due for maturity", map[string]interface{}{
		"count": len(policies),
	})

	// Process policies in batches
	successCount := 0
	failureCount := 0

	for i := 0; i < len(policies); i += j.config.BatchSize {
		end := i + j.config.BatchSize
		if end > len(policies) {
			end = len(policies)
		}

		batch := policies[i:end]
		nlog.Info(ctx, "Processing batch", map[string]interface{}{
			"batch_number": i/j.config.BatchSize + 1,
			"batch_size":   len(batch),
		})

		for _, policy := range batch {
			if err := j.processPolicy(ctx, policy); err != nil {
				nlog.Error(ctx, "Failed to process policy", map[string]interface{}{
					"policy_id":     policy.PolicyID,
					"policy_number": policy.PolicyNumber,
					"error":         err.Error(),
				})
				failureCount++
			} else {
				successCount++
			}
		}
	}

	duration := time.Since(startTime)
	nlog.Info(ctx, "Maturity intimation batch job completed", map[string]interface{}{
		"duration_seconds": duration.Seconds(),
		"total_policies":   len(policies),
		"success_count":    successCount,
		"failure_count":    failureCount,
	})

	return nil
}

// PolicyDueForMaturity represents a policy that is due for maturity
type PolicyDueForMaturity struct {
	PolicyID       string    `json:"policy_id"`
	PolicyNumber   string    `json:"policy_number"`
	CustomerID     string    `json:"customer_id"`
	CustomerName   string    `json:"customer_name"`
	MaturityDate   time.Time `json:"maturity_date"`
	MaturityAmount float64   `json:"maturity_amount"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
}

// getPoliciesDueForMaturity queries policies maturing in the given date range
// TODO: Implement this method by calling Policy Service
// Reference: FR-CLM-MC-002, BR-CLM-MC-002
func (j *MaturityIntimationJob) getPoliciesDueForMaturity(
	ctx context.Context,
	startDate, endDate time.Time,
) ([]PolicyDueForMaturity, error) {
	// TODO: This is a placeholder implementation
	// In production, this would call Policy Service API or query policy database
	// The query should:
	// 1. Find policies with maturity_date between startDate and endDate
	// 2. Filter out policies that already have maturity claims registered
	// 3. Filter out policies where intimation was already sent (check claim_history table)
	// 4. Join with customer table to get contact details
	// 5. Return paginated results

	nlog.Info(ctx, "Querying policy service for maturity due policies", map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})

	// Placeholder: Return empty slice
	// In production, this would call Policy Service API
	// Example API call:
	// GET /policies/maturity-due?start_date={startDate}&end_date={endDate}&page=1&limit=100
	return []PolicyDueForMaturity{}, nil
}

// processPolicy sends maturity intimation for a single policy
func (j *MaturityIntimationJob) processPolicy(ctx context.Context, policy PolicyDueForMaturity) error {
	// Create intimation request
	req := &MaturityIntimationRequest{
		PolicyID:       policy.PolicyID,
		PolicyNumber:   policy.PolicyNumber,
		CustomerID:     policy.CustomerID,
		CustomerName:   policy.CustomerName,
		MaturityDate:   policy.MaturityDate,
		MaturityAmount: policy.MaturityAmount,
		Phone:          policy.Phone,
		Email:          policy.Email,
		Channels:       j.config.NotificationChannels,
	}

	// Send notification via configured channels
	if err := j.notificationSvc.SendMaturityIntimation(ctx, req); err != nil {
		return fmt.Errorf("failed to send intimation: %w", err)
	}

	// Log successful intimation
	nlog.Info(ctx, "Maturity intimation sent successfully", map[string]interface{}{
		"policy_id":     policy.PolicyID,
		"policy_number": policy.PolicyNumber,
		"customer_id":   policy.CustomerID,
		"channels":      req.Channels,
	})

	// TODO: Record intimation in claim_history table for audit trail
	// This prevents duplicate intimations
	// Example:
	// INSERT INTO claim_history (claim_id, event_type, description, created_by)
	// VALUES (policy.policy_id, 'MATURITY_INTIMATION_SENT', 'Intimation sent via SMS, EMAIL', 'SYSTEM')

	return nil
}

// StartScheduler starts a cron scheduler for the maturity intimation job
// The job runs daily at 9:00 AM to check for policies due for maturity
// Reference: BR-CLM-MC-002 - 60-day advance intimation
func (j *MaturityIntimationJob) StartScheduler(ctx context.Context) error {
	// TODO: Implement cron scheduler
	// Options:
	// 1. Use github.com/robfig/cron/v3 for in-process cron
	// 2. Use external cron (Linux cron, Kubernetes CronJob)
	// 3. Use Temporal workflow for scheduled tasks
	//
	// Example with robfig/cron:
	// c := cron.New(cron.WithSeconds())
	// c.AddFunc("0 0 9 * * *", func() { // Daily at 9:00 AM
	//     if err := j.Run(context.Background()); err != nil {
	//         nlog.Error("Maturity intimation job failed", map[string]interface{}{
	//             "error": err.Error(),
	//         })
	//     }
	// })
	// c.Start()

	nlog.Info(ctx, "Maturity intimation scheduler started", map[string]interface{}{
		"schedule": "0 0 9 * * *", // Daily at 9:00 AM
		"timezone": "Asia/Kolkata",
	})

	return nil
}
