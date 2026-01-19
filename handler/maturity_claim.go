package handler

import (
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/domain"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// MaturityClaimHandler handles maturity claim-related HTTP requests
type MaturityClaimHandler struct {
	*serverHandler.Base
	claimRepo    *repo.ClaimRepository
	claimDocRepo *repo.ClaimDocumentRepository
}

// NewMaturityClaimHandler creates a new maturity claim handler
func NewMaturityClaimHandler(claimRepo *repo.ClaimRepository, claimDocRepo *repo.ClaimDocumentRepository) *MaturityClaimHandler {
	base := serverHandler.New("MaturityClaims").
		SetPrefix("/v1").
		AddPrefix("")
	return &MaturityClaimHandler{
		Base:         base,
		claimRepo:    claimRepo,
		claimDocRepo: claimDocRepo,
	}
}

// Routes defines all routes for this handler
func (h *MaturityClaimHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		// Maturity Claims (12 endpoints)
		serverRoute.POST("/claims/maturity/send-intimation-batch", h.SendMaturityIntimationBatch).Name("Send Maturity Intimation Batch"),
		serverRoute.POST("/claims/maturity/generate-due-report", h.GenerateMaturityDueReport).Name("Generate Maturity Due Report"),
		serverRoute.GET("/claims/maturity/pre-fill-data", h.GetMaturityPreFillData).Name("Get Maturity Pre-Fill Data"),
		serverRoute.POST("/claims/maturity/submit", h.SubmitMaturityClaim).Name("Submit Maturity Claim"),
		serverRoute.POST("/claims/maturity/:claim_id/validate-documents", h.ValidateMaturityDocuments).Name("Validate Maturity Documents"),
		serverRoute.POST("/claims/maturity/:claim_id/extract-ocr-data", h.ExtractOCRData).Name("Extract OCR Data"),
		serverRoute.POST("/claims/maturity/:claim_id/qc-verify", h.QCVerifyMaturityClaim).Name("QC Verify Maturity Claim"),
		serverRoute.POST("/claims/maturity/:claim_id/validate-bank", h.ValidateMaturityBankAccount).Name("Validate Maturity Bank Account"),
		serverRoute.GET("/claims/maturity/:claim_id/approval-details", h.GetMaturityApprovalDetails).Name("Get Maturity Approval Details"),
		serverRoute.POST("/claims/maturity/:claim_id/approve", h.ApproveMaturityClaim).Name("Approve Maturity Claim"),
		serverRoute.POST("/claims/maturity/:claim_id/disburse", h.DisburseMaturityClaim).Name("Disburse Maturity Claim"),
		serverRoute.POST("/claims/maturity/:claim_id/generate-voucher", h.GenerateMaturityVoucher).Name("Generate Maturity Voucher"),
	}
}

// SendMaturityIntimationBatch sends intimation for policies maturing in 60 days
// POST /claims/maturity/send-intimation-batch
// Reference: FR-CLM-MC-002, BR-CLM-MC-002 (60 days before maturity)
func (h *MaturityClaimHandler) SendMaturityIntimationBatch(sctx *serverRoute.Context, req SendMaturityIntimationBatchRequest) (*resp.MaturityIntimationBatchResponse, error) {
	// TODO: Implement batch intimation logic
	// 1. Query policies maturing between maturity_date_from and maturity_date_to
	// 2. Filter by channels specified
	// 3. Send notifications via Notification Service
	// 4. Track intimations sent

	log.Info(sctx.Ctx, "Sending maturity intimation batch from %s to %s via %v", req.MaturityDateFrom, req.MaturityDateTo, req.Channels)

	// Placeholder response
	r := &resp.MaturityIntimationBatchResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		TotalPolicies:        0,
		IntimationsSent:      0,
		Failed:               0,
	}

	return r, nil
}

// GenerateMaturityDueReport generates monthly maturity due report
// POST /claims/maturity/generate-due-report
func (h *MaturityClaimHandler) GenerateMaturityDueReport(sctx *serverRoute.Context, req GenerateMaturityDueReportRequest) (*resp.MaturityDueReportResponse, error) {
	// TODO: Implement report generation logic
	// 1. Query policies maturing in the specified month/year
	// 2. Aggregate data (total policies, total amount, etc.)
	// 3. Generate report in PDF/Excel format
	// 4. Upload to document storage
	// 5. Return report URL

	log.Info(sctx.Ctx, "Generating maturity due report for %d/%d", req.ReportMonth, req.ReportYear)

	// Placeholder response
	r := &resp.MaturityDueReportResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		ReportID:             "RPT" + time.Now().Format("20060102150405"),
		ReportURL:            "https://docs.example.com/reports/maturity_due_report.pdf",
		TotalPolicies:        0,
		TotalAmount:          0.0,
	}

	return r, nil
}

// GetMaturityPreFillData retrieves pre-filled data for maturity claim form
// GET /claims/maturity/pre-fill-data
func (h *MaturityClaimHandler) GetMaturityPreFillData(sctx *serverRoute.Context, req GetMaturityPreFillDataRequest) (*resp.MaturityPreFillDataResponse, error) {
	// TODO: Implement pre-fill data retrieval
	// 1. Validate token from intimation email/SMS
	// 2. Query policy details from Policy Service
	// 3. Query customer details from Customer Service
	// 4. Query bank details on record
	// 5. Calculate maturity amount

	log.Info(sctx.Ctx, "Getting pre-fill data for policy %s with token %s", req.PolicyID, req.Token)

	// Placeholder response
	r := &resp.MaturityPreFillDataResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		PolicyID:             req.PolicyID,
		CustomerName:         "John Doe",
		MaturityDate:         "2024-03-15",
		MaturityAmount:       500000.00,
		RegisteredMobile:     "9876543210",
		RegisteredEmail:      "john.doe@example.com",
		BankDetailsOnRecord: &resp.BankDetails{
			BankName:          "State Bank of India",
			BankAccountNumber: "1234567890",
			BankIFSC:          "SBIN0001234",
			BankAccountType:   "Savings",
			AccountHolderName: "John Doe",
		},
	}

	return r, nil
}

// SubmitMaturityClaim submits a new maturity claim
// POST /claims/maturity/submit
// Reference: BR-CLM-MC-001 (7 days SLA)
func (h *MaturityClaimHandler) SubmitMaturityClaim(sctx *serverRoute.Context, req SubmitMaturityClaimRequest) (*resp.MaturityClaimRegistrationResponse, error) {
	// Calculate SLA due date (7 days from submission)
	slaDueDate := time.Now().AddDate(0, 0, 7)

	// Create domain claim object
	data := domain.Claim{
		PolicyID:          req.PolicyID,
		ClaimType:         "MATURITY",
		ClaimantName:      req.ClaimantName,
		ClaimantPhone:     &req.ClaimantMobile,
		ClaimantEmail:     &req.ClaimantEmail,
		PaymentMode:       &req.DisbursementMode,
		BankAccountNumber: &req.BankAccountNumber,
		BankIFSCCode:      &req.BankIFSC,
		Status:            "SUBMITTED",
		SLADueDate:        slaDueDate,
	}

	if req.ClaimantRelationship != "" {
		data.ClaimantRelation = &req.ClaimantRelationship
	}

	// Call repository to create claim
	result, err := h.claimRepo.Create(sctx.Ctx, data)
	if err != nil {
		log.Error(sctx.Ctx, "Error submitting maturity claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Maturity claim submitted with ID: %s, Claim Number: %s", result.ID, result.ClaimNumber)

	// TODO: Upload documents to ECMS
	// TODO: Trigger OCR extraction for documents

	// Build response
	r := &resp.MaturityClaimRegistrationResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		ClaimID:              result.ID,
		ClaimNumber:          result.ClaimNumber,
		PolicyID:             result.PolicyID,
		SubmissionDate:       result.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if result.SLADueDate != (time.Time{}) {
		r.SLADueDate = result.SLADueDate.Format("2006-01-02 15:04:05")
	}

	return r, nil
}

// ValidateMaturityDocuments validates uploaded documents against checklist
// POST /claims/maturity/{claim_id}/validate-documents
func (h *MaturityClaimHandler) ValidateMaturityDocuments(sctx *serverRoute.Context, req ClaimIDUri) (*resp.DocumentsValidatedResponse, error) {
	// TODO: Implement document validation logic
	// 1. Query document checklist for maturity claims
	// 2. Query uploaded documents for the claim
	// 3. Verify mandatory documents are present
	// 4. Verify document verification status
	// 5. Return validation result

	log.Info(sctx.Ctx, "Validating documents for claim %s", req.ClaimID)

	// Placeholder response
	r := &resp.DocumentsValidatedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Valid:                true,
		MandatoryDocuments:   3,
		MandatoryVerified:    3,
		OptionalDocuments:    2,
		OptionalVerified:     1,
		MissingDocuments:     []string{},
	}

	return r, nil
}

// ExtractOCRData extracts data from documents using OCR
// POST /claims/maturity/{claim_id}/extract-ocr-data
func (h *MaturityClaimHandler) ExtractOCRData(sctx *serverRoute.Context, req ExtractOCRDataRequest) (*resp.OCRDataExtractedResponse, error) {
	// TODO: Implement OCR extraction logic
	// 1. Call ECMS OCR service for each document
	// 2. Extract fields (policy number, claimant name, bank details, etc.)
	// 3. Return extracted data with confidence score
	// 4. Map extracted fields to claim form fields

	log.Info(sctx.Ctx, "Extracting OCR data for claim %s, documents: %v", req.ClaimID, req.DocumentIDs)

	// Placeholder response
	r := &resp.OCRDataExtractedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		ExtractedData: map[string]interface{}{
			"policy_number":        "POL123456",
			"claimant_name":        "John Doe",
			"bank_account_number":  "1234567890",
			"bank_ifsc":            "SBIN0001234",
		},
		ConfidenceScore: 0.95,
		FieldsExtracted: []string{"policy_number", "claimant_name", "bank_account_number", "bank_ifsc"},
	}

	return r, nil
}

// QCVerifyMaturityClaim performs QC verification of OCR data
// POST /claims/maturity/{claim_id}/qc-verify
func (h *MaturityClaimHandler) QCVerifyMaturityClaim(sctx *serverRoute.Context, req QCVerifyMaturityClaimRequest) (*resp.QCVerificationResponse, error) {
	// TODO: Implement QC verification logic
	// 1. Update claim with QC status
	// 2. Store corrections if any
	// 3. Store QC remarks
	// 4. Trigger re-processing if corrections required
	// 5. Move to next workflow step if approved

	log.Info(sctx.Ctx, "QC verifying claim %s with status %s", req.ClaimID, req.QCStatus)

	// Update claim status
	updates := map[string]interface{}{
		"status":   "QC_VERIFIED",
		"qc_status": req.QCStatus,
	}

	if req.QCRemarks != nil {
		updates["qc_remarks"] = *req.QCRemarks
	}

	_, err := h.claimRepo.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating QC verification: %v", err)
		return nil, err
	}

	// Build response
	r := &resp.QCVerificationResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		QCStatus:             req.QCStatus,
		QCRemarks:            req.QCRemarks,
		QCVerifiedBy:         "qc_user@example.com", // TODO: Get from user context
		QCVerificationDate:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return r, nil
}

// ValidateMaturityBankAccount validates bank account for disbursement
// POST /claims/maturity/{claim_id}/validate-bank
func (h *MaturityClaimHandler) ValidateMaturityBankAccount(sctx *serverRoute.Context, req ClaimIDUri) (*resp.BankValidationData, error) {
	// TODO: Implement bank validation logic
	// 1. Get claim details including bank account
	// 2. Call CBS API for bank account validation
	// 3. Perform penny drop test
	// 4. Return validation result

	log.Info(sctx.Ctx, "Validating bank account for claim %s", req.ClaimID)

	// Placeholder response - using existing BankValidationData from claim.go
	r := &resp.BankValidationData{
		Valid:               true,
		AccountNumber:       "1234567890",
		AccountHolderName:   "John Doe",
		BankName:            "State Bank of India",
		ValidationMethod:    "CBS",
		NameMatchPercentage: 100.0,
	}

	return r, nil
}

// GetMaturityApprovalDetails retrieves claim details for approval
// GET /claims/maturity/{claim_id}/approval-details
func (h *MaturityClaimHandler) GetMaturityApprovalDetails(sctx *serverRoute.Context, req ClaimIDUri) (*resp.MaturityApprovalDetailsResponse, error) {
	// Fetch claim details
	claim, err := h.claimRepo.FindByID(sctx.Ctx, req.ClaimID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(sctx.Ctx, "Claim not found: %s", req.ClaimID)
			return nil, err
		}
		log.Error(sctx.Ctx, "Error fetching claim: %v", err)
		return nil, err
	}

	// TODO: Get document details from ClaimDocumentRepository
	// TODO: Get approval history from ClaimHistoryRepository
	// TODO: Get eligible approvers from User Service
	// TODO: Calculate benefit amount from Policy Service

	log.Info(sctx.Ctx, "Retrieved approval details for claim %s", req.ClaimID)

	// Build response with safe pointer handling
	r := &resp.MaturityApprovalDetailsResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		ClaimID:              claim.ID,
		ClaimNumber:          claim.ClaimNumber,
		PolicyID:             claim.PolicyID,
		ClaimantDetails: resp.MaturityClaimantData{
			Name: claim.ClaimantName,
		},
		DocumentDetails: resp.MaturityDocumentData{
			TotalDocuments:    3,
			VerifiedDocuments: 3,
			PendingDocuments:  0,
			DocumentChecklist: []resp.DocumentChecklistItem{},
		},
		CalculationDetails: resp.MaturityCalculationData{
			SumAssured:       500000.00,
			TotalAmount:      500000.00,
			NetPayableAmount: 500000.00,
		},
		ApprovalHistory:      []resp.ApprovalHistoryItem{},
		CurrentApprovalLevel: "LEVEL_1",
		EligibleApprovers:    []resp.ApproverInfo{},
	}

	// Handle optional fields safely
	if claim.ClaimantRelation != nil {
		r.ClaimantDetails.Relationship = *claim.ClaimantRelation
	}
	if claim.ClaimantPhone != nil {
		r.ClaimantDetails.Mobile = *claim.ClaimantPhone
	}
	if claim.ClaimantEmail != nil {
		r.ClaimantDetails.Email = *claim.ClaimantEmail
	}
	if claim.PaymentMode != nil {
		r.DisbursementDetails.DisbursementMode = *claim.PaymentMode
	}
	if claim.BankAccountNumber != nil {
		r.DisbursementDetails.BankAccountNumber = *claim.BankAccountNumber
	}
	if claim.BankIFSCCode != nil {
		r.DisbursementDetails.BankIFSC = *claim.BankIFSCCode
	}
	if claim.ClaimAmount != nil {
		r.DisbursementDetails.DisbursementAmount = *claim.ClaimAmount
		r.CalculationDetails.SumAssured = *claim.ClaimAmount
		r.CalculationDetails.TotalAmount = *claim.ClaimAmount
		r.CalculationDetails.NetPayableAmount = *claim.ClaimAmount
	}

	return r, nil
}

// ApproveMaturityClaim approves or rejects a maturity claim
// POST /claims/maturity/{claim_id}/approve
func (h *MaturityClaimHandler) ApproveMaturityClaim(sctx *serverRoute.Context, req ApproveMaturityClaimRequest) (*resp.MaturityClaimApprovedResponse, error) {
	// Determine new status
	newStatus := "APPROVED"
	if req.ApprovalStatus == "REJECTED" {
		newStatus = "REJECTED"
	}

	// Update claim
	updates := map[string]interface{}{
		"status":           newStatus,
		"approved_amount":  req.ApprovalAmount,
		"approver_id":      req.ApproverID,
		"approval_remarks": req.ApprovalRemarks,
	}

	_, err := h.claimRepo.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error approving maturity claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Maturity claim %s %s with amount %f", req.ClaimID, req.ApprovalStatus, req.ApprovalAmount)

	// TODO: Create claim history entry
	// TODO: Trigger disbursement workflow if approved
	// TODO: Send notification to claimant

	// Build response
	r := &resp.MaturityClaimApprovedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		ClaimID:              req.ClaimID,
		ApprovalStatus:       req.ApprovalStatus,
		ApprovalAmount:       req.ApprovalAmount,
		ApprovalRemarks:      req.ApprovalRemarks,
		ApprovedBy:           req.ApproverID,
		ApprovalDate:         time.Now().Format("2006-01-02 15:04:05"),
		ApprovalLevel:        req.ApprovalLevel,
	}

	return r, nil
}

// DisburseMaturityClaim initiates disbursement for approved claim
// POST /claims/maturity/{claim_id}/disburse
func (h *MaturityClaimHandler) DisburseMaturityClaim(sctx *serverRoute.Context, req DisburseMaturityClaimRequest) (*resp.MaturityClaimDisbursementInitiatedResponse, error) {
	// Parse disbursement date
	disbursementDate, err := time.Parse("2006-01-02", req.DisbursementDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid disbursement date format: %v", err)
		return nil, err
	}

	// Update claim status to DISBURSED
	updates := map[string]interface{}{
		"status":            "DISBURSED",
		"disbursement_date": disbursementDate,
		"payment_mode":       req.DisbursementMode,
		"payment_reference": req.ReferenceNumber,
	}

	if req.UTRNumber != nil {
		updates["utr_number"] = *req.UTRNumber
	}

	_, err = h.claimRepo.Update(sctx.Ctx, req.ClaimID, updates)
	if err != nil {
		log.Error(sctx.Ctx, "Error disbursing maturity claim: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Maturity claim %s disbursed with amount %f via %s", req.ClaimID, req.DisbursementAmount, req.DisbursementMode)

	// TODO: Create claim payment record
	// TODO: Call PFMS/CBS API for NEFT transfer
	// TODO: Send disbursement confirmation to claimant

	// Build response
	r := &resp.MaturityClaimDisbursementInitiatedResponse{
		StatusCodeAndMessage:  port.CreateSuccess,
		ClaimID:               req.ClaimID,
		DisbursementID:        "DIS" + time.Now().Format("20060102150405"),
		DisbursementAmount:    req.DisbursementAmount,
		DisbursementMode:      req.DisbursementMode,
		ReferenceNumber:       req.ReferenceNumber,
		EstimatedTransferDate: disbursementDate.AddDate(0, 0, 1).Format("2006-01-02"),
	}

	return r, nil
}

// GenerateMaturityVoucher generates payment voucher
// POST /claims/maturity/{claim_id}/generate-voucher
func (h *MaturityClaimHandler) GenerateMaturityVoucher(sctx *serverRoute.Context, req ClaimIDUri) (*resp.MaturityVoucherGeneratedResponse, error) {
	// TODO: Implement voucher generation logic
	// 1. Fetch claim and disbursement details
	// 2. Generate voucher number (format: VOU{YYYY}{DDDD})
	// 3. Create voucher PDF
	// 4. Upload to document storage
	// 5. Return voucher URL

	log.Info(sctx.Ctx, "Generating voucher for claim %s", req.ClaimID)

	// Placeholder response
	r := &resp.MaturityVoucherGeneratedResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		VoucherNumber:        "VOU" + time.Now().Format("2006") + "0001",
		VoucherDate:          time.Now().Format("2006-01-02"),
		VoucherURL:           "https://docs.example.com/vouchers/maturity_voucher.pdf",
		ClaimID:              req.ClaimID,
		DisbursementAmount:   500000.00,
	}

	return r, nil
}
