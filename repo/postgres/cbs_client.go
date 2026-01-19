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

// CBSClient handles integration with CBS (Core Banking System) API
// Reference: INT-CLM-016 (CBS API Integration)
// Reference: BR-CLM-MC-003 (Bank Verification Requirement)
// Reference: FR-CLM-MC-010 (Bank Account Validation API-based)
type CBSClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	timeout    time.Duration
}

// CBSAccountValidationRequest represents the request payload for CBS account validation
type CBSAccountValidationRequest struct {
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
	AccountHolderName string `json:"account_holder_name"`
}

// CBSAccountValidationResponse represents the response from CBS account validation API
type CBSAccountValidationResponse struct {
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

// CBSPennyDropRequest represents the request payload for CBS penny drop test
type CBSPennyDropRequest struct {
	AccountNumber     string  `json:"account_number"`
	IFSCCode          string  `json:"ifsc_code"`
	AccountHolderName string  `json:"account_holder_name"`
	Amount            float64 `json:"amount"` // Usually 1.0 for penny drop
	ReferenceID       string  `json:"reference_id"`
}

// CBSPennyDropResponse represents the response from CBS penny drop API
type CBSPennyDropResponse struct {
	Success           bool    `json:"success"`
	TransactionID     string  `json:"transaction_id"`
	ReferenceID       string  `json:"reference_id"`
	AccountNumber     string  `json:"account_number"`
	AccountHolderName string  `json:"account_holder_name"`
	NameMatchPercentage float64 `json:"name_match_percentage"`
	Amount            float64 `json:"amount"`
	CreditDate        string  `json:"credit_date"`
	Status            string  `json:"status"` // CREDITED, PENDING, FAILED
	ResponseCode      string  `json:"response_code"`
	ResponseMessage   string  `json:"response_message"`
}

// CBSAPIError represents an error response from CBS API
type CBSAPIError struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

// Error implements the error interface
func (e *CBSAPIError) Error() string {
	return fmt.Sprintf("CBS API Error [%d]: %s - %s", e.StatusCode, e.ErrorCode, e.Message)
}

// NewCBSClient creates a new CBS API client
// Reference: seed/template/template.md - External Service Client Pattern
func NewCBSClient(cfg *config.Config) *CBSClient {
	// Read CBS API configuration from config
	timeout := 30 * time.Second
	if cfg != nil {
		if cfg.GetInt("api_clients.cbs.timeout") > 0 {
			timeout = time.Duration(cfg.GetInt("api_clients.cbs.timeout")) * time.Second
		}
	}

	baseURL := "https://cbs-api.pli.gov.in/api/v1" // Default CBS API endpoint
	if cfg != nil {
		if cfg.GetString("api_clients.cbs.base_url") != "" {
			baseURL = cfg.GetString("api_clients.cbs.base_url")
		}
	}

	apiKey := "" // API key loaded from config or environment
	if cfg != nil {
		apiKey = cfg.GetString("api_clients.cbs.api_key")
	}

	return &CBSClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
		timeout: timeout,
	}
}

// ValidateBankAccount validates bank account details via CBS API
// Reference: FR-CLM-MC-010 (Bank Account Validation API-based)
// Reference: FR-CLM-SB-010 (Bank Account Validation for Survival Benefit)
// Reference: VR-CLM-API-002 (CBS/PFMS Bank Account API)
// Reference: BR-CLM-MC-003 (Bank Verification Requirement)
//
// Business Rules:
// - BR-CLM-DC-018: Bank validation must be completed before disbursement
// - BR-CLM-MC-003: Bank account must be verified via CBS/PFMS API before disbursement
//
// Validation Checks:
// 1. Account number format and existence
// 2. IFSC code validity
// 3. Account holder name match (with percentage)
// 4. Account status (ACTIVE, INACTIVE, CLOSED)
// 5. Account type (SAVINGS, CURRENT, NRE, NRO)
func (c *CBSClient) ValidateBankAccount(ctx context.Context, req CBSAccountValidationRequest) (*CBSAccountValidationResponse, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Prepare request payload
	requestBody, err := json.Marshal(req)
	if err != nil {
		log.Error(ctx, "Failed to marshal CBS request: %v", err)
		return nil, fmt.Errorf("failed to marshal CBS request: %w", err)
	}

	// Build CBS API endpoint URL
	url := fmt.Sprintf("%s/bank/validate", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestBody))
	if err != nil {
		log.Error(ctx, "Failed to create CBS HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create CBS HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	// Execute HTTP request
	log.Info(ctx, "Calling CBS API for bank validation: account=%s, ifsc=%s", req.AccountNumber, req.IFSCCode)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, "CBS API request failed: %v", err)
		return nil, fmt.Errorf("CBS API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Failed to read CBS API response: %v", err)
		return nil, fmt.Errorf("failed to read CBS API response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		var apiErr CBSAPIError
		if err := json.Unmarshal(responseBody, &apiErr); err == nil {
			log.Error(ctx, "CBS API returned error: %v", apiErr)
			return nil, &apiErr
		}
		err := fmt.Errorf("CBS API returned status %d: %s", resp.StatusCode, string(responseBody))
		log.Error(ctx, "CBS API error: %v", err)
		return nil, err
	}

	// Parse successful response
	var validationResp CBSAccountValidationResponse
	if err := json.Unmarshal(responseBody, &validationResp); err != nil {
		log.Error(ctx, "Failed to parse CBS API response: %v", err)
		return nil, fmt.Errorf("failed to parse CBS API response: %w", err)
	}

	log.Info(ctx, "CBS API validation successful: valid=%v, name_match=%.2f%%, status=%s",
		validationResp.Valid, validationResp.NameMatchPercentage, validationResp.AccountStatus)

	return &validationResp, nil
}

// PerformPennyDrop performs penny drop test via CBS API
// Reference: BR-CLM-DC-010 (Bank Account Validation)
//
// Penny Drop Process:
// 1. Initiate 1 rupee transfer to beneficiary account
// 2. Wait for credit confirmation (usually within minutes)
// 3. Verify account holder name from credit transaction
// 4. Calculate name match percentage
// 5. Automatically reverse the penny drop amount
//
// This is the most reliable method for bank account validation as it verifies:
// - Account exists and is active
// - Account can receive credits
// - Account holder name matches exactly
func (c *CBSClient) PerformPennyDrop(ctx context.Context, req CBSPennyDropRequest) (*CBSPennyDropResponse, error) {
	// Create context with timeout (penny drop may take longer)
	ctx, cancel := context.WithTimeout(ctx, c.timeout*2) // Double timeout for penny drop
	defer cancel()

	// Set default amount if not provided
	if req.Amount <= 0 {
		req.Amount = 1.0 // Standard penny drop amount
	}

	// Generate reference ID if not provided
	if req.ReferenceID == "" {
		req.ReferenceID = fmt.Sprintf("PennyDrop-%d", time.Now().UnixNano())
	}

	// Prepare request payload
	requestBody, err := json.Marshal(req)
	if err != nil {
		log.Error(ctx, "Failed to marshal CBS penny drop request: %v", err)
		return nil, fmt.Errorf("failed to marshal CBS penny drop request: %w", err)
	}

	// Build CBS API endpoint URL
	url := fmt.Sprintf("%s/bank/penny-drop", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestBody))
	if err != nil {
		log.Error(ctx, "Failed to create CBS penny drop HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create CBS penny drop HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	// Execute HTTP request
	log.Info(ctx, "Calling CBS API for penny drop: account=%s, ifsc=%s, amount=%.2f, ref=%s",
		req.AccountNumber, req.IFSCCode, req.Amount, req.ReferenceID)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, "CBS penny drop API request failed: %v", err)
		return nil, fmt.Errorf("CBS penny drop API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Failed to read CBS penny drop API response: %v", err)
		return nil, fmt.Errorf("failed to read CBS penny drop API response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		var apiErr CBSAPIError
		if err := json.Unmarshal(responseBody, &apiErr); err == nil {
			log.Error(ctx, "CBS penny drop API returned error: %v", apiErr)
			return nil, &apiErr
		}
		err := fmt.Errorf("CBS penny drop API returned status %d: %s", resp.StatusCode, string(responseBody))
		log.Error(ctx, "CBS penny drop API error: %v", err)
		return nil, err
	}

	// Parse successful response
	var pennyDropResp CBSPennyDropResponse
	if err := json.Unmarshal(responseBody, &pennyDropResp); err != nil {
		log.Error(ctx, "Failed to parse CBS penny drop API response: %v", err)
		return nil, fmt.Errorf("failed to parse CBS penny drop API response: %w", err)
	}

	log.Info(ctx, "CBS penny drop API successful: success=%v, txn_id=%s, status=%s, name_match=%.2f%%",
		pennyDropResp.Success, pennyDropResp.TransactionID, pennyDropResp.Status, pennyDropResp.NameMatchPercentage)

	return &pennyDropResp, nil
}

// GetPennyDropStatus checks the status of a penny drop transaction
// This can be used to poll for completion if the penny drop is processed asynchronously
func (c *CBSClient) GetPennyDropStatus(ctx context.Context, referenceID string) (*CBSPennyDropResponse, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build CBS API endpoint URL
	url := fmt.Sprintf("%s/bank/penny-drop/status/%s", c.baseURL, referenceID)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error(ctx, "Failed to create CBS penny drop status HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create CBS penny drop status HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	// Execute HTTP request
	log.Info(ctx, "Calling CBS API for penny drop status: ref=%s", referenceID)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, "CBS penny drop status API request failed: %v", err)
		return nil, fmt.Errorf("CBS penny drop status API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Failed to read CBS penny drop status API response: %v", err)
		return nil, fmt.Errorf("failed to read CBS penny drop status API response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		var apiErr CBSAPIError
		if err := json.Unmarshal(responseBody, &apiErr); err == nil {
			log.Error(ctx, "CBS penny drop status API returned error: %v", apiErr)
			return nil, &apiErr
		}
		err := fmt.Errorf("CBS penny drop status API returned status %d: %s", resp.StatusCode, string(responseBody))
		log.Error(ctx, "CBS penny drop status API error: %v", err)
		return nil, err
	}

	// Parse successful response
	var pennyDropResp CBSPennyDropResponse
	if err := json.Unmarshal(responseBody, &pennyDropResp); err != nil {
		log.Error(ctx, "Failed to parse CBS penny drop status API response: %v", err)
		return nil, fmt.Errorf("failed to parse CBS penny drop status API response: %w", err)
	}

	log.Info(ctx, "CBS penny drop status API successful: status=%s, name_match=%.2f%%",
		pennyDropResp.Status, pennyDropResp.NameMatchPercentage)

	return &pennyDropResp, nil
}

// ReversePennyDrop reverses a penny drop transaction (refund the amount)
// This should be called after successful penny drop verification
func (c *CBSClient) ReversePennyDrop(ctx context.Context, transactionID string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build CBS API endpoint URL
	url := fmt.Sprintf("%s/bank/penny-drop/reverse/%s", c.baseURL, transactionID)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		log.Error(ctx, "Failed to create CBS penny drop reverse HTTP request: %v", err)
		return fmt.Errorf("failed to create CBS penny drop reverse HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	// Execute HTTP request
	log.Info(ctx, "Calling CBS API to reverse penny drop: txn_id=%s", transactionID)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, "CBS penny drop reverse API request failed: %v", err)
		return fmt.Errorf("CBS penny drop reverse API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("CBS penny drop reverse API returned status %d: %s", resp.StatusCode, string(responseBody))
		log.Error(ctx, "CBS penny drop reverse API error: %v", err)
		return err
	}

	log.Info(ctx, "CBS penny drop reversed successfully: txn_id=%s", transactionID)
	return nil
}

// HealthCheck performs a health check on the CBS API
// This can be used for monitoring and alerting
func (c *CBSClient) HealthCheck(ctx context.Context) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Short timeout for health check
	defer cancel()

	// Build CBS API health endpoint URL
	url := fmt.Sprintf("%s/health", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error(ctx, "Failed to create CBS health check HTTP request: %v", err)
		return fmt.Errorf("failed to create CBS health check HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	// Execute HTTP request
	log.Info(ctx, "Calling CBS API health check")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, "CBS health check API request failed: %v", err)
		return fmt.Errorf("CBS health check API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("CBS health check returned status %d", resp.StatusCode)
		log.Error(ctx, "CBS health check error: %v", err)
		return err
	}

	log.Info(ctx, "CBS API health check successful")
	return nil
}
