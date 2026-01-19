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
	pfmsClient       *repo.PFMSClient
}

// NewBankingHandler creates a new banking handler
func NewBankingHandler(claimPaymentRepo *repo.ClaimPaymentRepository, claimRepo *repo.ClaimRepository, cbsClient *repo.CBSClient, pfmsClient *repo.PFMSClient) *BankingHandler {
	base := serverHandler.New("Banking").
		SetPrefix("/v1").
		AddPrefix("")
	return &BankingHandler{
		Base:             base,
		claimPaymentRepo: claimPaymentRepo,
		claimRepo:        claimRepo,
		cbsClient:        cbsClient,
		pfmsClient:       pfmsClient,
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
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: VR-CLM-API-002 (CBS/PFMS Bank Account API)
func (h *BankingHandler) ValidateViaPFMS(sctx *serverRoute.Context, req BankValidationRequest) (*resp.ExtendedBankValidationResponse, error) {
	// Call PFMS API for bank account validation
	pfmsReq := repo.PFMSBankValidationRequest{
		AccountNumber:     req.AccountNumber,
		IFSCCode:          req.IFSCCode,
		AccountHolderName: req.AccountHolderName,
	}

	pfmsResp, err := h.pfmsClient.ValidateBankAccount(sctx.Ctx, pfmsReq)
	if err != nil {
		log.Error(sctx.Ctx, "PFMS bank validation failed: %v", err)
		// Return error response with validation failure
		response := resp.NewBankValidationResponse(
			false, // valid
			req.AccountNumber,
			req.AccountHolderName,
			"",
			"PFMS_API",
			0.0, // nameMatchPercentage
		)
		response.StatusCode = 500
		response.Success = false
		response.Message = "Bank validation via PFMS API failed"
		return response, err
	}

	// Map PFMS API response to ExtendedBankValidationResponse
	validationMethod := "PFMS_API"
	response := resp.NewBankValidationResponse(
		pfmsResp.Valid, // valid
		pfmsResp.AccountNumber,
		pfmsResp.AccountHolderName,
		pfmsResp.BankName,
		validationMethod,
		pfmsResp.NameMatchPercentage,
	)

	// Set additional fields
	response.Data.IFSCCode = pfmsResp.IFSCCode
	response.Data.AccountStatus = &pfmsResp.AccountStatus
	response.Data.AccountType = &pfmsResp.AccountType
	response.Data.BranchName = &pfmsResp.BranchName
	response.Data.City = &pfmsResp.City
	response.Data.State = &pfmsResp.State
	response.Data.PINCode = &pfmsResp.PINCode
	response.Data.MICRCode = &pfmsResp.MICRCode

	log.Info(sctx.Ctx, "Bank account validated via PFMS: account=%s, valid=%v, name_match=%.2f%%, bank=%s",
		pfmsResp.AccountNumber, pfmsResp.Valid, pfmsResp.NameMatchPercentage, pfmsResp.BankName)

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
// Reference: INT-CLM-017 (PFMS API Integration)
// Reference: INT-CLM-018 (PFMS Integration for NEFT)
func (h *BankingHandler) InitiateNEFTTransfer(sctx *serverRoute.Context, req InitiateNEFTTransferRequest) (*resp.NEFTTransferInitiatedResponse, error) {
	// Validate claim ID
	if req.ClaimID == nil || *req.ClaimID == "" {
		return nil, fmt.Errorf("claim_id is required")
	}

	// Validate claim status (must be APPROVED)
	claim, err := h.claimRepo.FindByID(sctx.Ctx, *req.ClaimID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to find claim: %v", err)
		return nil, err
	}

	if claim.Status != "APPROVED" {
		response := &resp.NEFTTransferInitiatedResponse{
			StatusCodeAndMessage: port.StatusCodeAndMessage{
				StatusCode: 400,
				Success:    false,
				Message:    "Claim must be in APPROVED status to initiate disbursement",
			},
		}
		return response, nil
	}

	// Prepare PFMS NEFT request
	paymentReference := fmt.Sprintf("CLAIM-%s-%d", claim.ClaimNumber, time.Now().Unix())
	if req.ReferenceID != nil {
		paymentReference = *req.ReferenceID
	}

	pfmsReq := repo.PFMSNEFTTransferRequest{
		BeneficiaryAccount: req.AccountNumber,
		BeneficiaryIFSC:    req.IFSCCode,
		BeneficiaryName:    req.BeneficiaryName,
		Amount:             req.Amount,
		PaymentReference:   paymentReference,
		Purpose:            "CLAIM_DISBURSEMENT",
		ClaimNumber:        claim.ClaimNumber,
		PolicyNumber:       claim.PolicyID, // Assuming PolicyID is the policy number
	}

	// Call PFMS API to initiate NEFT transfer
	pfmsResp, err := h.pfmsClient.InitiateNEFTTransfer(sctx.Ctx, pfmsReq)
	if err != nil {
		log.Error(sctx.Ctx, "PFMS NEFT transfer failed: %v", err)
		response := &resp.NEFTTransferInitiatedResponse{
			StatusCodeAndMessage: port.StatusCodeAndMessage{
				StatusCode: 500,
				Success:    false,
				Message:    "NEFT transfer initiation failed",
			},
		}
		return response, err
	}

	// TODO: Record payment in claim_payments table
	// TODO: Update claim status to DISBURSED (or DISBURSING if processing)

	// Map PFMS API response to NEFTTransferInitiatedResponse
	response := &resp.NEFTTransferInitiatedResponse{
		StatusCodeAndMessage: port.StatusCodeAndMessage{
			StatusCode: 200,
			Success:    pfmsResp.Success,
			Message:    "NEFT transfer initiated successfully",
		},
		PaymentID:     pfmsResp.TransactionID,
		TransactionID: pfmsResp.TransactionID,
		ReferenceID:   pfmsResp.ReferenceNumber,
		UTR:           &pfmsResp.UTR,
		Amount:        pfmsResp.Amount,
		InitiatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Status:        pfmsResp.Status,
		BeneficiaryName: pfmsResp.BeneficiaryName,
	}

	log.Info(sctx.Ctx, "NEFT transfer initiated: claim_id=%s, payment_id=%s, utr=%s, amount=%.2f, status=%s",
		*req.ClaimID, pfmsResp.TransactionID, pfmsResp.UTR, pfmsResp.Amount, pfmsResp.Status)

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
// Reference: INT-CLM-017 (PFMS API Integration)
func (h *BankingHandler) GetPaymentStatus(sctx *serverRoute.Context, req PaymentIDUri) (*resp.PaymentStatusResponse, error) {
	// TODO: Query payment from claim_payments table
	// For now, fetch directly from PFMS API

	// Fetch payment status from PFMS API
	pfmsResp, err := h.pfmsClient.GetPaymentStatus(sctx.Ctx, req.PaymentID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch payment status from PFMS: %v", err)
		return nil, err
	}

	// Map PFMS response to PaymentStatusResponse
	response := resp.NewPaymentStatusResponse(
		pfmsResp.TransactionID,
		pfmsResp.Status,
		pfmsResp.Amount,
	)

	// Set additional fields from PFMS response
	response.TransactionID = &pfmsResp.TransactionID
	response.PaymentReference = &pfmsResp.ReferenceNumber
	response.AccountNumber = &pfmsResp.BeneficiaryAccount
	response.BeneficiaryName = &pfmsResp.BeneficiaryName

	if pfmsResp.InitiatedAt != nil {
		initiatedAt := pfmsResp.InitiatedAt.Format("2006-01-02 15:04:05")
		response.InitiatedAt = &initiatedAt
	}
	if pfmsResp.CompletedAt != nil {
		completedAt := pfmsResp.CompletedAt.Format("2006-01-02 15:04:05")
		response.CompletedAt = &completedAt
	}
	if pfmsResp.UTR != "" {
		response.UtrNumber = &pfmsResp.UTR
	}
	if pfmsResp.FailureReason != "" {
		response.FailureReason = &pfmsResp.FailureReason
	}

	// TODO: Update payment status in claim_payments table if changed

	log.Info(sctx.Ctx, "Payment status retrieved: payment_id=%s, status=%s, amount=%.2f",
		req.PaymentID, pfmsResp.Status, pfmsResp.Amount)

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
