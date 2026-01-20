package handler

import (
	"fmt"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
	"gitlab.cept.gov.in/pli/claims-api/core/port"
	resp "gitlab.cept.gov.in/pli/claims-api/handler/response"
	repo "gitlab.cept.gov.in/pli/claims-api/repo/postgres"
)

// ValidationServiceHandler handles validation service endpoints
type ValidationServiceHandler struct {
	*serverHandler.Base
	cbsClient      *repo.CBSClient
	pfmsClient     *repo.PFMSClient
	claimRepo      *repo.ClaimRepository
	policyBondRepo *repo.PolicyBondTrackingRepository
}

// NewValidationServiceHandler creates a new validation service handler
func NewValidationServiceHandler(cbsClient *repo.CBSClient, pfmsClient *repo.PFMSClient, claimRepo *repo.ClaimRepository, policyBondRepo *repo.PolicyBondTrackingRepository) *ValidationServiceHandler {
	base := serverHandler.New("ValidationService").
		SetPrefix("/v1").
		AddPrefix("")
	return &ValidationServiceHandler{
		Base:           base,
		cbsClient:      cbsClient,
		pfmsClient:     pfmsClient,
		claimRepo:      claimRepo,
		policyBondRepo: policyBondRepo,
	}
}

// Routes defines all routes for this handler
func (h *ValidationServiceHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/validate/pan", h.ValidatePAN).Name("Validate PAN"),
		serverRoute.POST("/validate/bank-account", h.ValidateBankAccountService).Name("Validate Bank Account"),
		serverRoute.POST("/validate/death-date", h.ValidateDeathDate).Name("Validate Death Date"),
		serverRoute.GET("/validate/ifsc/:ifsc_code", h.ValidateIFSC).Name("Validate IFSC"),
		serverRoute.GET("/forms/death-claim/fields", h.GetDynamicFormFields).Name("Get Dynamic Form Fields"),
	}
}

// ValidatePAN validates PAN number
// POST /validate/pan
// Reference: VR-CLM-VAL-001 (PAN validation via NSDL/Customer Service API)
func (h *ValidationServiceHandler) ValidatePAN(sctx *serverRoute.Context, req ValidatePANRequest) (*resp.PANValidationResponse, error) {
	log.Info(sctx.Ctx, "Validating PAN number: %s", req.PANNumber)

	// PAN validation rules (as per Income Tax Department)
	// 1. Length must be 10 characters
	// 2. Format: 5 letters + 4 digits + 1 letter
	// 3. Case-insensitive
	panRegex := regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)
	panUpper := regexp.MustCompile(`^[a-zA-Z]{5}[0-9]{4}[a-zA-Z]{1}$`)

	isValid := false
	var panStatus, panType *string
	var nameOnPAN *string

	// Validate format
	if panRegex.MatchString(req.PANNumber) {
		isValid = true
		status := "ACTIVE"
		panStatus = &status
		typeStr := "INDIVIDUAL"
		panType = &typeStr
		log.Info(sctx.Ctx, "PAN format is valid: %s", req.PANNumber)
	} else if panUpper.MatchString(req.PANNumber) {
		// Convert to uppercase for validation
		isValid = true
		status := "ACTIVE"
		panStatus = &status
		typeStr := "INDIVIDUAL"
		panType = &typeStr
		log.Info(sctx.Ctx, "PAN format is valid (case-insensitive): %s", req.PANNumber)
	} else {
		isValid = false
		status := "INVALID"
		panStatus = &status
		log.Error(sctx.Ctx, "Invalid PAN format: %s", req.PANNumber)
	}

	// TODO: Integrate with Customer Service API for real PAN validation
	// INT-CLM-005: Customer Service Integration
	// This will verify PAN against NSDL database and get actual name

	r := resp.NewPANValidationResponse(isValid, req.PANNumber, nameOnPAN, panStatus, panType)
	return r, nil
}

// ValidateBankAccountService validates bank account via CBS/PFMS
// POST /validate/bank-account
// Reference: VR-CLM-VAL-002 (Bank account validation via CBS/PFMS API)
func (h *ValidationServiceHandler) ValidateBankAccountService(sctx *serverRoute.Context, req ValidateBankAccountServiceRequest) (*resp.BankAccountValidationResponse, error) {
	log.Info(sctx.Ctx, "Validating bank account: %s with IFSC: %s using method: %s", req.BankAccountNumber, req.BankIFSC, req.ValidationMethod)

	// Determine validation method (default: CBS)
	validationMethod := req.ValidationMethod
	if validationMethod == "" {
		validationMethod = "CBS"
	}

	var isValid bool
	var accountHolderName, bankName, branchName, accountType, accountStatus *string
	var nameMatchPercentage *int
	var err error

	// Route to appropriate validation method
	switch validationMethod {
	case "CBS":
		// Call CBS API for validation
		cbsResp, cbsErr := h.cbsClient.ValidateBankAccount(sctx.Ctx, repo.CBSAccountValidationRequest{
			AccountNumber: req.BankAccountNumber,
			IFSC:          req.BankIFSC,
		})

		if cbsErr != nil {
			log.Error(sctx.Ctx, "CBS API error: %v", cbsErr)
			// Return invalid response on error
			isValid = false
		} else {
			isValid = cbsResp.Valid
			accountHolderName = &cbsResp.AccountHolderName
			bankName = &cbsResp.BankName
			branchName = &cbsResp.BranchName
			accountType = &cbsResp.AccountType
			accountStatus = &cbsResp.AccountStatus
			nameMatchPercentage = &cbsResp.NameMatchPercentage
		}

	case "PFMS":
		// Call PFMS API for validation
		pfmsResp, pfmsErr := h.pfmsClient.ValidateBankAccount(sctx.Ctx, req.BankAccountNumber, req.BankIFSC)

		if pfmsErr != nil {
			log.Error(sctx.Ctx, "PFMS API error: %v", pfmsErr)
			isValid = false
		} else {
			isValid = pfmsResp.Valid
			accountHolderName = &pfmsResp.AccountHolderName
			bankName = &pfmsResp.BankName
			branchName = &pfmsResp.BranchName
			accountType = &pfmsResp.AccountType
			accountStatus = &pfmsResp.AccountStatus
			nameMatchPercentage = &pfmsResp.NameMatchPercentage
		}

	case "PENNY_DROP":
		// Perform penny drop test via CBS
		pennyResp, pennyErr := h.cbsClient.PerformPennyDrop(sctx.Ctx, repo.CBSPennyDropRequest{
			AccountNumber: req.BankAccountNumber,
			IFSC:          req.BankIFSC,
		})

		if pennyErr != nil {
			log.Error(sctx.Ctx, "Penny drop error: %v", pennyErr)
			isValid = false
		} else {
			isValid = pennyResp.Valid
			accountHolderName = &pennyResp.AccountHolderName
			bankName = &pennyResp.BankName
			branchName = &pennyResp.BranchName
			accountStatus = &pennyResp.AccountStatus
			// Name match is 100% if penny drop succeeds
			match := 100
			nameMatchPercentage = &match
		}

	default:
		err = fmt.Errorf("invalid validation method: %s", validationMethod)
		log.Error(sctx.Ctx, "Error: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Bank account validation result: valid=%t, bank=%s, branch=%s", isValid, getStringValue(bankName), getStringValue(branchName))

	r := resp.NewBankAccountValidationResponse(isValid, req.BankAccountNumber, req.BankIFSC, accountHolderName, bankName, branchName, accountType, accountStatus, nameMatchPercentage, validationMethod)
	return r, nil
}

// ValidateDeathDate validates death date against policy dates
// POST /validate/death-date
// Reference: BR-CLM-DC-001 (Investigation trigger based on death within 3 years)
func (h *ValidationServiceHandler) ValidateDeathDate(sctx *serverRoute.Context, req ValidateDeathDateRequest) (*resp.DeathDateValidationResponse, error) {
	log.Info(sctx.Ctx, "Validating death date: %s for policy: %s", req.DeathDate, req.PolicyID)

	// Parse death date
	deathDate, err := time.Parse("2006-01-02", req.DeathDate)
	if err != nil {
		log.Error(sctx.Ctx, "Invalid death date format: %v", err)
		return nil, fmt.Errorf("invalid death date format: %w", err)
	}

	// TODO: Integrate with Policy Service to get policy details
	// INT-CLM-002: Policy Service Integration
	// For now, using mock data

	// Mock policy dates (should come from Policy Service)
	var policyIssueDate, policyRevivalDate *string
	var investigationRequired bool
	var validationMessages []string

	// Calculate days from policy issue (mock)
	// BR-CLM-DC-001: Investigation required if death within 3 years of policy issue/revival
	// For demonstration, assume policy was issued 2 years ago
	daysFromIssue := 730 // 2 years
	daysFromIssuePtr := &daysFromIssue

	policyIssue := "2022-01-15"
	policyIssueDate = &policyIssue

	// Check if death within 3 years (1095 days)
	if daysFromIssue <= 1095 {
		investigationRequired = true
		validationMessages = append(validationMessages, "Death occurred within 3 years of policy issue date. Investigation required as per BR-CLM-DC-001.")
		log.Info(sctx.Ctx, "Investigation required: death within 3 years of policy issue (days: %d)", daysFromIssue)
	} else {
		investigationRequired = false
		validationMessages = append(validationMessages, "Death occurred after 3 years of policy issue date. No investigation required.")
		log.Info(sctx.Ctx, "No investigation required: death after 3 years of policy issue (days: %d)", daysFromIssue)
	}

	// Validate death date is not in future
	if deathDate.After(time.Now()) {
		validationMessages = append(validationMessages, "Death date cannot be in the future.")
		log.Error(sctx.Ctx, "Invalid death date: future date detected")
	}

	r := resp.NewDeathDateValidationResponse(true, req.DeathDate, req.PolicyID, investigationRequired, validationMessages)
	// Add additional fields
	r.Data.PolicyIssueDate = policyIssueDate
	r.Data.PolicyRevivalDate = policyRevivalDate
	r.Data.DaysFromIssue = daysFromIssuePtr

	return r, nil
}

// ValidateIFSC validates IFSC code and retrieves bank details
// GET /validate/ifsc/{ifsc_code}
// Reference: VR-CLM-VAL-003 (IFSC validation via RBI IFSC code bank)
func (h *ValidationServiceHandler) ValidateIFSC(sctx *serverRoute.Context, req IFSCCodeUri) (*resp.IFSCValidationResponse, error) {
	log.Info(sctx.Ctx, "Validating IFSC code: %s", req.IFSCCode)

	// IFSC validation rules (as per RBI)
	// 1. Length must be 11 characters
	// 2. Format: 4 letters (bank code) + 0 (zero) + 6 alphanumeric (branch code)
	// 3. Case-insensitive
	ifscRegex := regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)
	ifscUpper := regexp.MustCompile(`^[a-zA-Z]{4}0[a-zA-Z0-9]{6}$`)

	isValid := ifscRegex.MatchString(req.IFSCCode) || ifscUpper.MatchString(req.IFSCCode)

	var bankName, branchName, address, city, state, district, pinCode, micrCode *string

	if isValid {
		// Extract bank code from IFSC (first 4 characters)
		bankCode := req.IFSCCode[0:4]

		// TODO: Integrate with RBI IFSC code database or use API
		// For now, using mock data
		bank := "Sample Bank"
		bankName = &bank

		branch := "Sample Branch"
		branchName = &branch

		addr := "123, Sample Street, Sample Area"
		address = &addr

		c := "MUMBAI"
		city = &c

		s := "MAHARASHTRA"
		state = &s

		d := "MUMBAI CITY"
		district = &d

		pin := "400001"
		pinCode = &pin

		micr := "400001001"
		micrCode = &micr

		log.Info(sctx.Ctx, "IFSC code is valid: %s (Bank: %s, Branch: %s)", req.IFSCCode, bank, branch)
	} else {
		log.Error(sctx.Ctx, "Invalid IFSC code format: %s", req.IFSCCode)
	}

	r := resp.NewIFSCValidationResponse(isValid, req.IFSCCode, bankName, branchName, address, city, state, district, pinCode, micrCode)
	return r, nil
}

// GetDynamicFormFields retrieves dynamic form fields based on death type
// GET /forms/death-claim/fields
// Reference: DFC-001 (Dynamic document checklist and form fields)
func (h *ValidationServiceHandler) GetDynamicFormFields(sctx *serverRoute.Context, req GetDynamicFormFieldsRequest) (*resp.DynamicFormFieldsResponse, error) {
	log.Info(sctx.Ctx, "Getting dynamic form fields for death type: %s", req.DeathType)

	var formFields []resp.FormField
	var documentList []resp.DocumentChecklist
	var validationRules []resp.ValidationRule

	// Define form fields based on death type
	switch req.DeathType {
	case "NATURAL":
		// Natural death form fields
		formFields = []resp.FormField{
			{
				FieldName:    "death_certificate_number",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "Death Certificate Number",
				Placeholder:  getStringPtr("Enter death certificate number"),
				DisplayOrder: 1,
			},
			{
				FieldName:    "place_of_death",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "Place of Death",
				Placeholder:  getStringPtr("Enter place of death"),
				DisplayOrder: 2,
			},
			{
				FieldName:    "cause_of_death",
				FieldType:    "SELECT",
				Required:     true,
				Label:        "Cause of Death",
				Options:      []string{"Illness", "Natural Causes", "Old Age", "Other"},
				DisplayOrder: 3,
			},
		}

		documentList = []resp.DocumentChecklist{
			{
				DocumentType:   "DEATH_CERTIFICATE",
				DocumentName:   "Death Certificate",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Original death certificate issued by municipal authority"),
				DisplayOrder:   1,
			},
			{
				DocumentType:   "CLAIMANT_ID_PROOF",
				DocumentName:   "Claimant ID Proof",
				Required:       true,
				DocumentFormat: getStringPtr("PDF,JPG,PNG"),
				MaxFileSizeMB:  getIntPtr(2),
				Description:    getStringPtr("Aadhaar Card / PAN Card / Passport"),
				DisplayOrder:   2,
			},
		}

	case "ACCIDENTAL":
		// Accidental death form fields
		formFields = []resp.FormField{
			{
				FieldName:    "fir_number",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "FIR Number",
				Placeholder:  getStringPtr("Enter FIR number"),
				DisplayOrder: 1,
			},
			{
				FieldName:    "accident_date",
				FieldType:    "DATE",
				Required:     true,
				Label:        "Accident Date",
				DisplayOrder: 2,
			},
			{
				FieldName:    "accident_place",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "Accident Place",
				Placeholder:  getStringPtr("Enter place of accident"),
				DisplayOrder: 3,
			},
			{
				FieldName:    "police_report_number",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "Police Report Number",
				Placeholder:  getStringPtr("Enter police report number"),
				DisplayOrder: 4,
			},
		}

		documentList = []resp.DocumentChecklist{
			{
				DocumentType:   "DEATH_CERTIFICATE",
				DocumentName:   "Death Certificate",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Original death certificate mentioning accidental death"),
				DisplayOrder:   1,
			},
			{
				DocumentType:   "FIR_COPY",
				DocumentName:   "FIR Copy",
				Required:       true,
				DocumentFormat: getStringPtr("PDF,JPG,PNG"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("First Information Report (FIR) copy"),
				DisplayOrder:   2,
			},
			{
				DocumentType:   "POLICE_REPORT",
				DocumentName:   "Police Report",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Final police report/charge sheet"),
				DisplayOrder:   3,
			},
			{
				DocumentType:   "POST_MORTEM_REPORT",
				DocumentName:   "Post Mortem Report",
				Required:       false,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Post mortem report if conducted"),
				DisplayOrder:   4,
			},
		}

	case "UNNATURAL":
		// Unnatural death form fields
		formFields = []resp.FormField{
			{
				FieldName:    "fir_number",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "FIR Number",
				Placeholder:  getStringPtr("Enter FIR number"),
				DisplayOrder: 1,
			},
			{
				FieldName:    "incident_date",
				FieldType:    "DATE",
				Required:     true,
				Label:        "Incident Date",
				DisplayOrder: 2,
			},
			{
				FieldName:    "cause_of_death",
				FieldType:    "SELECT",
				Required:     true,
				Label:        "Cause of Death",
				Options:      []string{"Homicide", "Suicide", "Poisoning", "Burns", "Drowning", "Other"},
				DisplayOrder: 3,
			},
		}

		documentList = []resp.DocumentChecklist{
			{
				DocumentType:   "DEATH_CERTIFICATE",
				DocumentName:   "Death Certificate",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Original death certificate"),
				DisplayOrder:   1,
			},
			{
				DocumentType:   "FIR_COPY",
				DocumentName:   "FIR Copy",
				Required:       true,
				DocumentFormat: getStringPtr("PDF,JPG,PNG"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("First Information Report (FIR) copy"),
				DisplayOrder:   2,
			},
			{
				DocumentType:   "POST_MORTEM_REPORT",
				DocumentName:   "Post Mortem Report",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Post mortem report is mandatory for unnatural death"),
				DisplayOrder:   3,
			},
			{
				DocumentType:   "POLICE_INVESTIGATION_REPORT",
				DocumentName:   "Police Investigation Report",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Complete police investigation report"),
				DisplayOrder:   4,
			},
		}

	case "SUICIDE":
		// Suicide form fields
		formFields = []resp.FormField{
			{
				FieldName:    "fir_number",
				FieldType:    "TEXT",
				Required:     true,
				Label:        "FIR Number",
				Placeholder:  getStringPtr("Enter FIR number"),
				DisplayOrder: 1,
			},
			{
				FieldName:    "suicide_note_available",
				FieldType:    "CHECKBOX",
				Required:     true,
				Label:        "Suicide Note Available",
				DisplayOrder: 2,
			},
			{
				FieldName:    "incident_date",
				FieldType:    "DATE",
				Required:     true,
				Label:        "Incident Date",
				DisplayOrder: 3,
			},
		}

		documentList = []resp.DocumentChecklist{
			{
				DocumentType:   "DEATH_CERTIFICATE",
				DocumentName:   "Death Certificate",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Original death certificate mentioning suicide"),
				DisplayOrder:   1,
			},
			{
				DocumentType:   "FIR_COPY",
				DocumentName:   "FIR Copy",
				Required:       true,
				DocumentFormat: getStringPtr("PDF,JPG,PNG"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("First Information Report (FIR) copy"),
				DisplayOrder:   2,
			},
			{
				DocumentType:   "POST_MORTEM_REPORT",
				DocumentName:   "Post Mortem Report",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Post mortem report is mandatory for suicide"),
				DisplayOrder:   3,
			},
			{
				DocumentType:   "POLICE_INVESTIGATION_REPORT",
				DocumentName:   "Police Investigation Report",
				Required:       true,
				DocumentFormat: getStringPtr("PDF"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Complete police investigation report with cause of death"),
				DisplayOrder:   4,
			},
			{
				DocumentType:   "SUICIDE_NOTE",
				DocumentName:   "Suicide Note",
				Required:       false,
				DocumentFormat: getStringPtr("PDF,JPG,PNG"),
				MaxFileSizeMB:  getIntPtr(5),
				Description:    getStringPtr("Suicide note if available"),
				DisplayOrder:   5,
			},
		}
	}

	// Common validation rules
	validationRules = []resp.ValidationRule{
		{
			RuleName:    "death_date_validation",
			RuleType:    "BUSINESS",
			Description: "Death date must not be in the future",
			Severity:    "ERROR",
			Message:     "Death date cannot be a future date",
		},
		{
			RuleName:    "document_mandatory",
			RuleType:    "DOCUMENT",
			Description: "All mandatory documents must be uploaded",
			Severity:    "ERROR",
			Message:     "Please upload all required documents before submission",
		},
		{
			RuleName:    "claimant_validation",
			RuleType:    "FIELD",
			Description: "Claimant details must match nominee records",
			Severity:    "WARNING",
			Message:     "Claimant will be verified against nominee records",
		},
	}

	// Add investigation rule for accidental/unnatural/suicide deaths
	if req.DeathType == "ACCIDENTAL" || req.DeathType == "UNNATURAL" || req.DeathType == "SUICIDE" {
		investigationRule := resp.ValidationRule{
			RuleName:    "investigation_mandatory",
			RuleType:    "BUSINESS",
			Description: "Investigation is mandatory for non-natural deaths",
			Severity:    "INFO",
			Message:     "This claim will be referred to the investigation department as per BR-CLM-DC-001",
		}
		validationRules = append(validationRules, investigationRule)
	}

	log.Info(sctx.Ctx, "Returning %d form fields, %d documents, %d validation rules for death type: %s",
		len(formFields), len(documentList), len(validationRules), req.DeathType)

	r := resp.NewDynamicFormFieldsResponse(req.DeathType, formFields, documentList, validationRules)
	return r, nil
}

// ==================== HELPER FUNCTIONS ====================

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getStringPtr returns a pointer to the given string
func getStringPtr(s string) *string {
	return &s
}

// getIntPtr returns a pointer to the given int
func getIntPtr(i int) *int {
	return &i
}

// ==================== ERROR HANDLING ====================

// handleError processes pgx errors and returns appropriate HTTP status
func (h *ValidationServiceHandler) handleError(err error) error {
	if err == pgx.ErrNoRows {
		return fmt.Errorf("record not found")
	}
	return err
}
