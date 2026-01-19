package handler

import (
	"fmt"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// BankingHandler handles banking and payment-related HTTP requests
type BankingHandler struct {
	*serverHandler.Base
	claimPaymentRepo *repo.ClaimPaymentRepository
	claimRepo        *repo.ClaimRepository
	cbsClient        *repo.CBSClient
}

// NewBankingHandler creates a new banking handler
func NewBankingHandler(claimPaymentRepo *repo.ClaimPaymentRepository, claimRepo *repo.ClaimRepository, cbsClient *repo.CBSClient) *BankingHandler {
	base := serverHandler.New("Banking").
		SetPrefix("/v1").
		AddPrefix("")
	return &BankingHandler{
		Base:             base,
		claimPaymentRepo: claimPaymentRepo,
		claimRepo:        claimRepo,
		cbsClient:        cbsClient,
	}
}

// Routes defines all routes for this handler
func (h *BankingHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Banking & Payment Services (8 endpoints)
		serverRoute.POST("/banking/validate-account", h.ValidateBankAccount).Name("Validate Bank Account"),
		serverRoute.POST("/banking/validate-account-cbs", h.ValidateViaCBS).Name("Validate via CBS API"),
		serverRoute.POST("/banking/validate-account-pfms", h.ValidateViaPFMS).Name("Validate via PFMS API"),
		serverRoute.POST("/banking/penny-drop", h.PerformPennyDrop).Name("Perform Penny Drop Test"),
		serverRoute.POST("/banking/neft-transfer", h.InitiateNEFTTransfer).Name("Initiate NEFT Transfer"),
		serverRoute.POST("/banking/payment-reconciliation", h.ReconcilePayments).Name("Reconcile Payments"),
		serverRoute.GET("/banking/payment-status/:payment_id", h.GetPaymentStatus).Name("Get Payment Status"),
		serverRoute.POST("/banking/generate-voucher", h.GeneratePaymentVoucher).Name("Generate Payment Voucher"),
	}
}

// ValidateBankAccount validates bank account via CBS/PFMS
// POST /banking/validate-account
// Reference: BR-CLM-DC-010 (Payment Disbursement Workflow)
// Reference: Integration with CBS API and PFMS API
func (h *BankingHandler) ValidateBankAccount(sctx *serverRoute.Context, req BankValidationRequest) (*resp.ExtendedBankValidationResponse, error) {
	// TODO: Integration with CBS/PFMS API for bank validation
	// TODO: Validate account number and IFSC code
	// TODO: Verify account holder name with name match percentage
	// TODO: Check account status (ACTIVE, INACTIVE, CLOSED)

	// Placeholder response
	validationMethod := "CBS_API"
	if req.ValidationMethod != nil {
		validationMethod = *req.ValidationMethod
	}

	response := resp.NewBankValidationResponse(
		true, // valid
		req.AccountNumber,
		req.AccountHolderName,
		"State Bank of India", // TODO: Get actual bank name from IFSC
		validationMethod,
		100.0, // nameMatchPercentage
	)

	// Set additional fields
	response.Data.IFSCCode = req.IFSCCode
	accountStatus := "ACTIVE"
	response.Data.AccountStatus = &accountStatus
	accountType := "SAVINGS"
	response.Data.AccountType = &accountType

	log.Info(sctx.Ctx, "Bank account validated for account: %s", req.AccountNumber)
	return response, nil
}

// ValidateViaCBS validates bank account via CBS API
// POST /banking/validate-account-cbs
// Reference: Integration with CBS API (Core Banking System)
// Reference: INT-CLM-016 (CBS API Integration)
// Reference: FR-CLM-MC-010 (Bank Account Validation API-based)
// Reference: VR-CLM-API-002 (CBS/PFMS Bank Account API)
func (h *BankingHandler) ValidateViaCBS(sctx *serverRoute.Context, req BankValidationRequest) (*resp.ExtendedBankValidationResponse, error) {
	// Call CBS API for bank account validation
	cbsReq := repo.CBSAccountValidationRequest{
		AccountNumber:     req.AccountNumber,
		IFSCCode:          req.IFSCCode,
		AccountHolderName: req.AccountHolderName,
	}

	cbsResp, err := h.cbsClient.ValidateBankAccount(sctx.Ctx, cbsReq)
	if err != nil {
		log.Error(sctx.Ctx, "CBS API validation failed: %v", err)
		// Return error response with validation failure
		response := resp.NewBankValidationResponse(
			false, // valid
			req.AccountNumber,
			req.AccountHolderName,
			"",
			"CBS_API",
			0.0, // nameMatchPercentage
		)
		response.StatusCode = 500
		response.Success = false
		response.Message = "Bank validation via CBS API failed"
		return response, err
	}

	// Map CBS API response to ExtendedBankValidationResponse
	response := resp.NewBankValidationResponse(
		cbsResp.Valid,
		cbsResp.AccountNumber,
		cbsResp.AccountHolderName,
		cbsResp.BankName,
		"CBS_API",
		cbsResp.NameMatchPercentage,
	)

	response.Data.IFSCCode = cbsResp.IFSCCode
	response.Data.AccountStatus = &cbsResp.AccountStatus
	response.Data.AccountType = &cbsResp.AccountType
	response.Data.BranchName = &cbsResp.BranchName
	response.Data.City = &cbsResp.City
	// Additional fields from CBS response
	if cbsResp.State != "" {
		response.Data.State = &cbsResp.State
	}
	if cbsResp.PINCode != "" {
		response.Data.PINCode = &cbsResp.PINCode
	}
	if cbsResp.MICRCode != "" {
		response.Data.MICRCode = &cbsResp.MICRCode
	}

	log.Info(sctx.Ctx, "Bank account validated via CBS for account: %s, valid=%v, name_match=%.2f%%",
		req.AccountNumber, cbsResp.Valid, cbsResp.NameMatchPercentage)

	return response, nil
}

// ValidateViaPFMS validates bank account via PFMS API
// POST /banking/validate-account-pfms
// Reference: Integration with PFMS API (Public Financial Management System)
func (h *BankingHandler) ValidateViaPFMS(sctx *serverRoute.Context, req BankValidationRequest) (*resp.ExtendedBankValidationResponse, error) {
	// TODO: Call PFMS API endpoint
	// TODO: Parse PFMS API response
	// TODO: Extract PFMS validation details

	// Placeholder response
	response := resp.NewBankValidationResponse(
		true, // valid
		req.AccountNumber,
		req.AccountHolderName,
		"State Bank of India", // TODO: Get from PFMS API
		"PFMS_API",
		100.0, // nameMatchPercentage
	)

	response.Data.IFSCCode = req.IFSCCode
	accountStatus := "ACTIVE"
	response.Data.AccountStatus = &accountStatus
	accountType := "SAVINGS"
	response.Data.AccountType = &accountType

	log.Info(sctx.Ctx, "Bank account validated via PFMS for account: %s", req.AccountNumber)
	return response, nil
}

// PerformPennyDrop performs penny drop test (1 rupee transfer)
// POST /banking/penny-drop
// Reference: BR-CLM-DC-010 (Bank Account Validation)
// Reference: INT-CLM-016 (CBS API Integration)
func (h *BankingHandler) PerformPennyDrop(sctx *serverRoute.Context, req BankValidationRequest) (*resp.ExtendedBankValidationResponse, error) {
	// Generate reference ID for penny drop
	referenceID := fmt.Sprintf("CLAIM-PennyDrop-%d", time.Now().UnixNano())

	// Call CBS API for penny drop
	cbsReq := repo.CBSPennyDropRequest{
		AccountNumber:     req.AccountNumber,
		IFSCCode:          req.IFSCCode,
		AccountHolderName: req.AccountHolderName,
		Amount:            1.0, // Standard penny drop amount
		ReferenceID:       referenceID,
	}

	cbsResp, err := h.cbsClient.PerformPennyDrop(sctx.Ctx, cbsReq)
	if err != nil {
		log.Error(sctx.Ctx, "CBS penny drop failed: %v", err)
		// Return error response with validation failure
		response := resp.NewBankValidationResponse(
			false, // valid
			req.AccountNumber,
			req.AccountHolderName,
			"",
			"PENNY_DROP",
			0.0, // nameMatchPercentage
		)
		response.StatusCode = 500
		response.Success = false
		response.Message = "Penny drop test via CBS API failed"
		return response, err
	}

	// Penny drop successful, reverse the transaction
	if cbsResp.Success && cbsResp.Status == "CREDITED" {
		err = h.cbsClient.ReversePennyDrop(sctx.Ctx, cbsResp.TransactionID)
		if err != nil {
			log.Error(sctx.Ctx, "Failed to reverse penny drop: %v", err)
			// Don't fail the response, just log the error
			// The reversal can be done manually later
		}
	}

	// Map CBS API response to ExtendedBankValidationResponse
	validationMethod := "PENNY_DROP"
	response := resp.NewBankValidationResponse(
		cbsResp.Success, // valid
		cbsResp.AccountNumber,
		cbsResp.AccountHolderName,
		"", // Bank name not available in penny drop response
		validationMethod,
		cbsResp.NameMatchPercentage,
	)

	response.Data.IFSCCode = req.IFSCCode
	accountStatus := "ACTIVE"
	if cbsResp.Status == "FAILED" {
		accountStatus = "INACTIVE"
	}
	response.Data.AccountStatus = &accountStatus

	log.Info(sctx.Ctx, "Penny drop test completed for account: %s, success=%v, name_match=%.2f%%, txn_id=%s",
		req.AccountNumber, cbsResp.Success, cbsResp.NameMatchPercentage, cbsResp.TransactionID)

	return response, nil
}

// InitiateNEFTTransfer initiates NEFT transfer
// POST /banking/neft-transfer
// Reference: BR-CLM-DC-010 (Disbursement Workflow)
func (h *BankingHandler) InitiateNEFTTransfer(sctx *serverRoute.Context, req InitiateNEFTTransferRequest) (*resp.NEFTTransferInitiatedResponse, error) {
	// TODO: Validate claim status (must be APPROVED)
	// TODO: Integration with PFMS API for NEFT
	// TODO: Generate payment_id
	// TODO: Initiate NEFT transfer
	// TODO: Record payment in claim_payments table
	// TODO: Update claim status to DISBURSED

	// Placeholder response
	response := &resp.NEFTTransferInitiatedResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "NEFT transfer initiated successfully",
		},
		PaymentID:     "PAY-2025-000001", // TODO: Generate actual payment ID
		TransactionID: "TXN-2025-000001", // TODO: Get from banking gateway
		ReferenceID:   "REF-2025-000001",
		Amount:        req.Amount,
		InitiatedAt:   "2025-01-20 12:00:00",
		Status:        "PROCESSING",
	}

	log.Info(sctx.Ctx, "NEFT transfer initiated for amount: %f", req.Amount)
	return response, nil
}

// ReconcilePayments performs daily payment reconciliation
// POST /banking/payment-reconciliation
// Reference: BR-CLM-PAY-001 (Daily Reconciliation)
func (h *BankingHandler) ReconcilePayments(sctx *serverRoute.Context, req ReconcilePaymentsRequest) (*resp.PaymentReconciliationResponse, error) {
	// TODO: Query all payments for reconciliation date
	// TODO: Fetch payment status from banking gateway
	// TODO: Match expected vs actual amounts
	// TODO: Identify mismatched transactions
	// TODO: Generate reconciliation report
	// TODO: Update payment statuses in database

	// Placeholder response
	response := &resp.PaymentReconciliationResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    true,
			Message:    "Payment reconciliation completed successfully",
		},
		ReconciliationDate:  req.ReconciliationDate,
		TotalPayments:      1000,
		SuccessfulPayments: 950,
		FailedPayments:     30,
		PendingPayments:    20,
		TotalAmountReconciled: 50000000.00,
		ReconciliationSummary: resp.ReconciliationSummaryData{
			MatchedCount:       950,
			MatchedAmount:      47500000.00,
			UnmatchedCount:     50,
			UnmatchedAmount:    2500000.00,
			ReconciliationRate: 95.0,
		},
		ReconciledAt: "2025-01-20 18:00:00",
	}

	log.Info(sctx.Ctx, "Payment reconciliation completed for date: %s", req.ReconciliationDate)
	return response, nil
}

// GetPaymentStatus retrieves real-time payment status
// GET /banking/payment-status/:payment_id
// Reference: BR-CLM-DC-010 (Payment Tracking)
func (h *BankingHandler) GetPaymentStatus(sctx *serverRoute.Context, req PaymentIDUri) (*resp.PaymentStatusResponse, error) {
	// TODO: Query payment from claim_payments table
	// TODO: If status is PROCESSING/PENDING, fetch latest status from banking gateway
	// TODO: Update payment status if changed
	// TODO: Return payment details

	// Placeholder response
	response := resp.NewPaymentStatusResponse(
		req.PaymentID,
		"SUCCESS", // INITIATED, PROCESSING, SUCCESS, FAILED, CANCELLED
		100000.00,
	)

	// Set additional fields
	transactionID := "TXN-2025-000001"
	response.TransactionID = &transactionID
	paymentReference := "PFMS-REF-001"
	response.PaymentReference = &paymentReference
	paymentDate := "2025-01-20 14:30:00"
	response.PaymentDate = &paymentDate
	bankName := "State Bank of India"
	response.BankName = &bankName
	accountNumber := "1234567890"
	response.AccountNumber = &accountNumber
	ifscCode := "SBIN0001234"
	response.IFSCCode = &ifscCode
	beneficiaryName := "John Doe"
	response.BeneficiaryName = &beneficiaryName
	paymentMode := "NEFT"
	response.PaymentMode = &paymentMode
	initiatedAt := "2025-01-20 12:00:00"
	response.InitiatedAt = &initiatedAt
	completedAt := "2025-01-20 14:30:00"
	response.CompletedAt = &completedAt
	utrNumber := "UTR123456789012"
	response.UtrNumber = &utrNumber

	log.Info(sctx.Ctx, "Payment status retrieved for payment_id: %s", req.PaymentID)
	return response, nil
}

// GeneratePaymentVoucher generates payment voucher for accounting
// POST /banking/generate-voucher
// Reference: BR-CLM-PAY-002 (Voucher Generation)
func (h *BankingHandler) GeneratePaymentVoucher(sctx *serverRoute.Context, req GeneratePaymentVoucherRequest) (*resp.PaymentVoucherResponse, error) {
	// TODO: Validate claim and payment existence
	// TODO: Fetch payment details from claim_payments table
	// TODO: Generate voucher number (VOU{YYYY}{DDDD})
	// TODO: Create voucher PDF with stamp (if include_stamp=true)
	// TODO: Upload voucher to ECMS
	// TODO: Return voucher URL

	// Placeholder response
	voucherType := "PAYMENT"
	if req.VoucherType != nil {
		voucherType = *req.VoucherType
	}

	response := resp.NewPaymentVoucherResponse(
		"VOU-2025-000001", // TODO: Generate actual voucher number
		req.PaymentID,
		req.ClaimID,
		100000.00, // amount
		"John Doe", // beneficiaryName
		"1234567890", // accountNumber
		"State Bank of India", // bankName
	)

	response.VoucherType = voucherType

	// Set voucher details
	accountingHead := "PLI-Claims-Payment"
	budgetHead := "PLI-2024-25"
	financialYear := "2024-25"
	voucherData := resp.VoucherDetails{
		PaymentDate:     "2025-01-20",
		AuthorizationBy:  nil, // TODO: Get from claim approval
		AccountingHead:  &accountingHead,
		BudgetHead:      &budgetHead,
		FinancialYear:   &financialYear,
		NetAmount:       100000.00,
		Remarks:         nil,
		SupportingDocsCount: 5,
	}
	response.VoucherData = voucherData

	voucherURL := "https://ecms.pli.gov.in/vouchers/VOU-2025-000001.pdf"
	response.VoucherURL = &voucherURL

	log.Info(sctx.Ctx, "Payment voucher generated for payment_id: %s", req.PaymentID)
	return response, nil
}
