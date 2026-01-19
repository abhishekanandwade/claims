package main

import (
	"context"
	"fmt"
	"os"
	"time"

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
	nlog "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	"gitlab.cept.gov.in/pli/claims-api/batch"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// MaturityIntimationBatchRunner is a standalone executable for running the maturity intimation batch job
// This can be run independently or scheduled via cron/Kubernetes CronJob
//
// Usage:
//   go run cmd/batch_runner/main.go
//
// Environment variables:
//   CONFIG_FILE - Path to config file (default: configs/config.yaml)
//
// Reference: FR-CLM-MC-002, BR-CLM-MC-002
func main() {
	ctx := context.Background()

	// Load configuration
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		nlog.Error(ctx, "Failed to load configuration", map[string]interface{}{
			"error": err.Error(),
			"path":  configPath,
		})
		os.Exit(1)
	}

	nlog.Info(ctx, "Configuration loaded successfully", map[string]interface{}{
		"config_path": configPath,
		"appname":     cfg.GetString("appname"),
	})

	// Initialize database connection
	db, err := dblib.NewDB(cfg)
	if err != nil {
		nlog.Error(ctx, "Failed to initialize database", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
	defer db.Close()

	nlog.Info(ctx, "Database connection established", map[string]interface{}{
		"host":     cfg.GetString("db.host"),
		"port":     cfg.GetString("db.port"),
		"database": cfg.GetString("db.database"),
	})

	// Initialize repositories
	claimRepo := repo.NewClaimRepository(db, cfg)

	// Initialize notification client
	notificationClient := repo.NewNotificationClient()

	// Load batch job configuration
	batchConfig := &batch.BatchConfig{
		MaturityIntimationDaysInAdvance: cfg.GetInt("batch.maturity_intimation.days_in_advance"),
		BatchSize:                        cfg.GetInt("batch.maturity_intimation.batch_size"),
		NotificationChannels:             cfg.GetStringSlice("batch.maturity_intimation.notification_channels"),
	}

	// Initialize batch job
	job := batch.NewMaturityIntimationJob(claimRepo, notificationClient, batchConfig)

	// Check if job is enabled
	if !cfg.GetBool("batch.maturity_intimation.enabled") {
		nlog.Info(ctx, "Maturity intimation batch job is disabled", map[string]interface{}{
			"config_key": "batch.maturity_intimation.enabled",
		})
		os.Exit(0)
	}

	// Run the batch job
	runCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	nlog.Info(ctx, "Starting maturity intimation batch job", map[string]interface{}{
		"days_in_advance": batchConfig.MaturityIntimationDaysInAdvance,
		"batch_size":      batchConfig.BatchSize,
		"channels":        batchConfig.NotificationChannels,
	})

	startTime := time.Now()
	if err := job.Run(runCtx); err != nil {
		nlog.Error(ctx, "Batch job failed", map[string]interface{}{
			"error":    err.Error(),
			"duration": time.Since(startTime).Seconds(),
		})
		os.Exit(1)
	}

	duration := time.Since(startTime)
	nlog.Info(ctx, "Batch job completed successfully", map[string]interface{}{
		"duration_seconds": duration.Seconds(),
	})

	fmt.Printf("\n✅ Maturity intimation batch job completed in %.2f seconds\n", duration.Seconds())
}
