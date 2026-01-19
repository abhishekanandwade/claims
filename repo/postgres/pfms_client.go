package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
)

// PFMSClient handles integration with PFMS (Public Financial Management System) API
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: BR-CLM-DC-010 (Payment Disbursement Workflow)
// Reference: BR-CLM-PAY-001 (Daily Payment Reconciliation)
type PFMSClient struct {
	httpClient    *http.Client
	baseURL       string
	apiKey        string
	timeout       time.Duration
	retryAttempts int
	retryDelay    time.Duration
}

// PFMSBankValidationRequest represents the request payload for PFMS bank validation
type PFMSBankValidationRequest struct {
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
	AccountHolderName string `json:"account_holder_name"`
}

// PFMSBankValidationResponse represents the response from PFMS bank validation API
type PFMSBankValidationResponse struct {
	Valid               bool    `json:"valid"`
	AccountNumber       string  `json:"account_number"`
	AccountHolderName   string  `json:"account_holder_name"`
	NameMatchPercentage float64 `json:"name_match_percentage"`
	BankName            string  `json:"bank_name"`
	IFSCCode            string  `json:"ifsc_code"`
	AccountStatus       string  `json:"account_status"` // ACTIVE, INACTIVE, CLOSED
	AccountType         string  `json:"account_type"`    // SAVINGS, CURRENT, NRE, NRO
	BranchName          string  `json:"branch_name"`
	City                string  `json:"city"`
	State               string  `json:"state"`
	PINCode             string  `json:"pincode"`
	MICRCode            string  `json:"micr_code"`
	ResponseCode        string  `json:"response_code"`
	ResponseMessage     string  `json:"response_message"`
}

// PFMSNEFTTransferRequest represents the request payload for NEFT transfer
type PFMSNEFTTransferRequest struct {
	BeneficiaryAccount   string  `json:"beneficiary_account"`
	BeneficiaryIFSC      string  `json:"beneficiary_ifsc"`
	BeneficiaryName      string  `json:"beneficiary_name"`
	Amount               float64 `json:"amount"`
	PaymentReference     string  `json:"payment_reference"`
	Purpose              string  `json:"purpose"`
	ClaimNumber          string  `json:"claim_number,omitempty"`
	PolicyNumber         string  `json:"policy_number,omitempty"`
	SchemeCode           string  `json:"scheme_code,omitempty"`
	DepartmentCode       string  `json:"department_code,omitempty"`
}

// PFMSNEFTTransferResponse represents the response from PFMS NEFT transfer API
type PFMSNEFTTransferResponse struct {
	Success         bool    `json:"success"`
	TransactionID   string  `json:"transaction_id"`
	ReferenceNumber string  `json:"reference_number"`
	UTR             string  `json:"utr"` // Unique Transaction Reference
	Status          string  `json:"status"` // INITIATED, PROCESSING, SUCCESS, FAILED
	Amount          float64 `json:"amount"`
	BeneficiaryName string  `json:"beneficiary_name"`
	ResponseCode    string  `json:"response_code"`
	ResponseMessage string  `json:"response_message"`
	ErrorDetails    string  `json:"error_details,omitempty"`
}

// PFMSPaymentStatusResponse represents the response from PFMS payment status API
type PFMSPaymentStatusResponse struct {
	TransactionID       string  `json:"transaction_id"`
	ReferenceNumber     string  `json:"reference_number"`
	UTR                 string  `json:"utr"`
	Status              string  `json:"status"` // INITIATED, PROCESSING, SUCCESS, FAILED, REVERSED
	Amount              float64 `json:"amount"`
	BeneficiaryAccount  string  `json:"beneficiary_account"`
	BeneficiaryName     string  `json:"beneficiary_name"`
	InitiatedAt         *time.Time `json:"initiated_at"`
	CompletedAt         *time.Time `json:"completed_at"`
	FailedAt            *time.Time `json:"failed_at,omitempty"`
	FailureReason       string  `json:"failure_reason,omitempty"`
	ResponseCode        string  `json:"response_code"`
	ResponseMessage     string  `json:"response_message"`
}

// PFMSAPIError represents an error from PFMS API
type PFMSAPIError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *PFMSAPIError) Error() string {
	return fmt.Sprintf("PFMS API error (status %d): %s - %s", e.StatusCode, e.Message, e.Details)
}

// NewPFMSClient creates a new PFMS API client
// Reference: config.yaml - api_clients.pfms section
func NewPFMSClient(cfg *config.Config) *PFMSClient {
	timeout := 30 * time.Second
	if val := cfg.Get("api_clients.pfms.timeout"); val != nil {
		if strVal, ok := val.(string); ok && strVal != "" {
			duration, err := time.ParseDuration(strVal + "s")
			if err == nil {
				timeout = duration
			}
		}
	}

	retryAttempts := 3
	if val := cfg.Get("api_clients.pfms.retry_attempts"); val != nil {
		if strVal, ok := val.(string); ok && strVal != "" {
			fmt.Sscanf(strVal, "%d", &retryAttempts)
		}
	}

	retryDelay := 1 * time.Second
	if val := cfg.Get("api_clients.pfms.retry_delay"); val != nil {
		if strVal, ok := val.(string); ok && strVal != "" {
			duration, err := time.ParseDuration(strVal + "s")
			if err == nil {
				retryDelay = duration
			}
		}
	}

	baseURL := ""
	if val := cfg.Get("api_clients.pfms.base_url"); val != nil {
		if strVal, ok := val.(string); ok {
			baseURL = strVal
		}
	}

	apiKey := ""
	if val := cfg.Get("api_clients.pfms.api_key"); val != nil {
		if strVal, ok := val.(string); ok {
			apiKey = strVal
		}
	}

	return &PFMSClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:       baseURL,
		apiKey:        apiKey,
		timeout:       timeout,
		retryAttempts: retryAttempts,
		retryDelay:    retryDelay,
	}
}

// ValidateBankAccount validates bank account via PFMS API
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: VR-CLM-API-002 (CBS/PFMS Bank Account API)
func (c *PFMSClient) ValidateBankAccount(ctx context.Context, req PFMSBankValidationRequest) (*PFMSBankValidationResponse, error) {
	// Check if PFMS integration is enabled
	if c.baseURL == "" {
		return nil, fmt.Errorf("PFMS API is not configured")
	}

	// Prepare request
	url := fmt.Sprintf("%s/bank/validate", c.baseURL)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with retries
	var lastErr error
	for attempt := 0; attempt < c.retryAttempts; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("X-API-Key", c.apiKey)

		// Send request
		log.Info(ctx, "Calling PFMS API for bank validation (attempt %d/%d): account=%s, ifsc=%s",
			attempt+1, c.retryAttempts, req.AccountNumber, req.IFSCCode)

		startTime := time.Now()
		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			log.Error(ctx, "PFMS API request failed (attempt %d/%d): %v", attempt+1, c.retryAttempts, err)
			time.Sleep(c.retryDelay)
			continue
		}

		duration := time.Since(startTime)
		log.Info(ctx, "PFMS API responded in %dms for account validation", duration.Milliseconds())

		// Read response body
		respBody, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			log.Error(ctx, "Failed to read PFMS API response body: %v", err)
			time.Sleep(c.retryDelay)
			continue
		}

		// Check HTTP status code
		if httpResp.StatusCode != http.StatusOK {
			lastErr = &PFMSAPIError{
				StatusCode: httpResp.StatusCode,
				Message:    httpResp.Status,
				Details:    string(respBody),
			}
			log.Error(ctx, "PFMS API returned non-OK status %d: %s", httpResp.StatusCode, string(respBody))
			return nil, lastErr
		}

		// Parse response
		var response PFMSBankValidationResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			lastErr = fmt.Errorf("failed to parse response: %w", err)
			log.Error(ctx, "Failed to parse PFMS API response: %v", err)
			return nil, lastErr
		}

		log.Info(ctx, "PFMS bank validation completed: valid=%v, account=%s, name_match=%.2f%%, bank=%s",
			response.Valid, response.AccountNumber, response.NameMatchPercentage, response.BankName)

		return &response, nil
	}

	return nil, fmt.Errorf("PFMS API request failed after %d attempts: %w", c.retryAttempts, lastErr)
}

// InitiateNEFTTransfer initiates NEFT transfer via PFMS API
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: BR-CLM-DC-010 (Payment Disbursement Workflow)
// Reference: INT-CLM-018 (PFMS Integration for NEFT)
func (c *PFMSClient) InitiateNEFTTransfer(ctx context.Context, req PFMSNEFTTransferRequest) (*PFMSNEFTTransferResponse, error) {
	// Check if PFMS integration is enabled
	if c.baseURL == "" {
		return nil, fmt.Errorf("PFMS API is not configured")
	}

	// Validate request
	if req.BeneficiaryAccount == "" || req.BeneficiaryIFSC == "" || req.Amount <= 0 {
		return nil, fmt.Errorf("invalid NEFT transfer request: missing required fields or invalid amount")
	}

	// Prepare request
	url := fmt.Sprintf("%s/payment/neft/transfer", c.baseURL)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with retries
	var lastErr error
	for attempt := 0; attempt < c.retryAttempts; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("X-API-Key", c.apiKey)

		// Send request
		log.Info(ctx, "Calling PFMS API for NEFT transfer (attempt %d/%d): beneficiary=%s, amount=%.2f, ref=%s",
			attempt+1, c.retryAttempts, req.BeneficiaryAccount, req.Amount, req.PaymentReference)

		startTime := time.Now()
		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			log.Error(ctx, "PFMS NEFT API request failed (attempt %d/%d): %v", attempt+1, c.retryAttempts, err)
			time.Sleep(c.retryDelay)
			continue
		}

		duration := time.Since(startTime)
		log.Info(ctx, "PFMS NEFT API responded in %dms", duration.Milliseconds())

		// Read response body
		respBody, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			log.Error(ctx, "Failed to read PFMS NEFT response body: %v", err)
			time.Sleep(c.retryDelay)
			continue
		}

		// Check HTTP status code
		if httpResp.StatusCode != http.StatusOK {
			lastErr = &PFMSAPIError{
				StatusCode: httpResp.StatusCode,
				Message:    httpResp.Status,
				Details:    string(respBody),
			}
			log.Error(ctx, "PFMS NEFT API returned non-OK status %d: %s", httpResp.StatusCode, string(respBody))
			return nil, lastErr
		}

		// Parse response
		var response PFMSNEFTTransferResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			lastErr = fmt.Errorf("failed to parse response: %w", err)
			log.Error(ctx, "Failed to parse PFMS NEFT response: %v", err)
			return nil, lastErr
		}

		log.Info(ctx, "PFMS NEFT transfer initiated: success=%v, txn_id=%s, utr=%s, status=%s, amount=%.2f",
			response.Success, response.TransactionID, response.UTR, response.Status, response.Amount)

		return &response, nil
	}

	return nil, fmt.Errorf("PFMS NEFT API request failed after %d attempts: %w", c.retryAttempts, lastErr)
}

// GetPaymentStatus retrieves payment status from PFMS API
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: BR-CLM-PAY-001 (Daily Payment Reconciliation)
func (c *PFMSClient) GetPaymentStatus(ctx context.Context, transactionID string) (*PFMSPaymentStatusResponse, error) {
	// Check if PFMS integration is enabled
	if c.baseURL == "" {
		return nil, fmt.Errorf("PFMS API is not configured")
	}

	// Validate transaction ID
	if transactionID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	// Prepare request
	url := fmt.Sprintf("%s/payment/status/%s", c.baseURL, transactionID)

	// Create HTTP request with retries
	var lastErr error
	for attempt := 0; attempt < c.retryAttempts; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("X-API-Key", c.apiKey)

		// Send request
		log.Info(ctx, "Calling PFMS API for payment status (attempt %d/%d): txn_id=%s",
			attempt+1, c.retryAttempts, transactionID)

		startTime := time.Now()
		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			log.Error(ctx, "PFMS status API request failed (attempt %d/%d): %v", attempt+1, c.retryAttempts, err)
			time.Sleep(c.retryDelay)
			continue
		}

		duration := time.Since(startTime)
		log.Info(ctx, "PFMS status API responded in %dms", duration.Milliseconds())

		// Read response body
		respBody, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			log.Error(ctx, "Failed to read PFMS status response body: %v", err)
			time.Sleep(c.retryDelay)
			continue
		}

		// Check HTTP status code
		if httpResp.StatusCode != http.StatusOK {
			lastErr = &PFMSAPIError{
				StatusCode: httpResp.StatusCode,
				Message:    httpResp.Status,
				Details:    string(respBody),
			}
			log.Error(ctx, "PFMS status API returned non-OK status %d: %s", httpResp.StatusCode, string(respBody))
			return nil, lastErr
		}

		// Parse response
		var response PFMSPaymentStatusResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			lastErr = fmt.Errorf("failed to parse response: %w", err)
			log.Error(ctx, "Failed to parse PFMS status response: %v", err)
			return nil, lastErr
		}

		log.Info(ctx, "PFMS payment status retrieved: txn_id=%s, utr=%s, status=%s, amount=%.2f",
			response.TransactionID, response.UTR, response.Status, response.Amount)

		return &response, nil
	}

	return nil, fmt.Errorf("PFMS status API request failed after %d attempts: %w", c.retryAttempts, lastErr)
}

// HealthCheck performs health check on PFMS API
// Reference: Monitoring and alerting requirements
func (c *PFMSClient) HealthCheck(ctx context.Context) error {
	if c.baseURL == "" {
		return fmt.Errorf("PFMS API is not configured")
	}

	url := fmt.Sprintf("%s/health", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("X-API-Key", c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("PFMS health check failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("PFMS health check returned status %d", httpResp.StatusCode)
	}

	log.Info(ctx, "PFMS API health check successful")
	return nil
}
