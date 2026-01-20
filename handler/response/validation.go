package response

import (
	"gitlab.cept.gov.in/pli/claims-api/core/port"
)

// ==================== VALIDATION SERVICE RESPONSE DTOs ====================

// PANValidationResponse represents the response for PAN validation
// POST /validate/pan
// Reference: VR-CLM-VAL-001 (PAN validation via NSDL/Customer Service)
type PANValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      PANValidationData `json:"data"`
}

// PANValidationData contains PAN validation details
type PANValidationData struct {
	Valid     bool    `json:"valid"`                 // true if PAN is valid
	PANNumber string  `json:"pan_number"`            // PAN number validated
	NameOnPAN *string `json:"name_on_pan,omitempty"` // Name as per PAN records
	PANStatus *string `json:"pan_status,omitempty"`  // ACTIVE, INACTIVE, INVALID
	PANType   *string `json:"pan_type,omitempty"`    // INDIVIDUAL, COMPANY, TRUST, etc.
}

// BankAccountValidationResponse represents the response for bank account validation
// POST /validate/bank-account
// Reference: VR-CLM-VAL-002 (Bank account validation via CBS/PFMS)
type BankAccountValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      BankAccountValidationData `json:"data"`
}

// BankAccountValidationData contains bank account validation details
type BankAccountValidationData struct {
	Valid               bool    `json:"valid"`                           // true if account is valid
	BankAccountNumber   string  `json:"bank_account_number"`             // Account number validated
	BankIFSC            string  `json:"bank_ifsc"`                       // IFSC code validated
	AccountHolderName   *string `json:"account_holder_name,omitempty"`   // Name as per bank records
	BankName            *string `json:"bank_name,omitempty"`             // Bank name
	BranchName          *string `json:"branch_name,omitempty"`           // Branch name
	AccountType         *string `json:"account_type,omitempty"`          // SAVINGS, CURRENT, NRE, NRO
	AccountStatus       *string `json:"account_status,omitempty"`        // ACTIVE, INACTIVE, CLOSED
	NameMatchPercentage *int    `json:"name_match_percentage,omitempty"` // Name match score (0-100)
	ValidationMethod    string  `json:"validation_method"`               // CBS, PFMS, PENNY_DROP
}

// DeathDateValidationResponse represents the response for death date validation
// POST /validate/death-date
// Reference: BR-CLM-DC-001 (Investigation trigger logic)
type DeathDateValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DeathDateValidationData `json:"data"`
}

// DeathDateValidationData contains death date validation details
type DeathDateValidationData struct {
	Valid                 bool     `json:"valid"`                         // true if death date is valid
	DeathDate             string   `json:"death_date"`                    // Death date validated
	PolicyID              string   `json:"policy_id"`                     // Policy ID
	PolicyIssueDate       *string  `json:"policy_issue_date,omitempty"`   // Policy issue date
	PolicyRevivalDate     *string  `json:"policy_revival_date,omitempty"` // Last revival date
	InvestigationRequired bool     `json:"investigation_required"`        // true if death within 3 years
	DaysFromIssue         *int     `json:"days_from_issue,omitempty"`     // Days from policy issue
	DaysFromRevival       *int     `json:"days_from_revival,omitempty"`   // Days from last revival
	ValidationMessages    []string `json:"validation_messages,omitempty"` // Validation messages
}

// IFSCValidationResponse represents the response for IFSC validation
// GET /validate/ifsc/{ifsc_code}
// Reference: VR-CLM-VAL-003 (IFSC validation via RBI IFSC code bank)
type IFSCValidationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      IFSCValidationData `json:"data"`
}

// IFSCValidationData contains IFSC validation details
type IFSCValidationData struct {
	Valid      bool    `json:"valid"`                 // true if IFSC is valid
	IFSCCode   string  `json:"ifsc_code"`             // IFSC code validated
	BankName   *string `json:"bank_name,omitempty"`   // Bank name
	BranchName *string `json:"branch_name,omitempty"` // Branch name
	Address    *string `json:"address,omitempty"`     // Branch address
	City       *string `json:"city,omitempty"`        // City
	State      *string `json:"state,omitempty"`       // State
	District   *string `json:"district,omitempty"`    // District
	PINCode    *string `json:"pin_code,omitempty"`    // PIN code
	MICRCode   *string `json:"micr_code,omitempty"`   // MICR code
}

// DynamicFormFieldsResponse represents the response for dynamic form fields
// GET /forms/death-claim/fields
// Reference: DFC-001 (Dynamic document checklist based on death type)
type DynamicFormFieldsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DynamicFormFieldsData `json:"data"`
}

// DynamicFormFieldsData contains dynamic form fields based on death type
type DynamicFormFieldsData struct {
	DeathType       string              `json:"death_type"`                 // Death type (NATURAL, ACCIDENTAL, etc.)
	FormFields      []FormField         `json:"form_fields,omitempty"`      // Dynamic form fields
	DocumentList    []DocumentChecklist `json:"document_list,omitempty"`    // Required documents
	ValidationRules []ValidationRule    `json:"validation_rules,omitempty"` // Validation rules
}

// FormField represents a dynamic form field
type FormField struct {
	FieldName     string      `json:"field_name"`               // Field name
	FieldType     string      `json:"field_type"`               // TEXT, NUMBER, DATE, SELECT, CHECKBOX
	Required      bool        `json:"required"`                 // true if mandatory
	Label         string      `json:"label"`                    // Field label
	Placeholder   *string     `json:"placeholder,omitempty"`    // Placeholder text
	Options       []string    `json:"options,omitempty"`        // Options for SELECT field
	DefaultValue  interface{} `json:"default_value,omitempty"`  // Default value
	Validation    *string     `json:"validation,omitempty"`     // Validation rules (regex, min, max, etc.)
	DisplayOrder  int         `json:"display_order"`            // Display order in form
	ConditionalOn *string     `json:"conditional_on,omitempty"` // Show only if this field has specific value
}

// DocumentChecklist represents a document in the checklist
type DocumentChecklist struct {
	DocumentType   string  `json:"document_type"`              // Document type
	DocumentName   string  `json:"document_name"`              // Document name
	Required       bool    `json:"required"`                   // true if mandatory
	DocumentFormat *string `json:"document_format,omitempty"`  // PDF, JPG, PNG, etc.
	MaxFileSizeMB  *int    `json:"max_file_size_mb,omitempty"` // Max file size in MB
	Description    *string `json:"description,omitempty"`      // Document description
	DisplayOrder   int     `json:"display_order"`              // Display order
	ConditionalOn  *string `json:"conditional_on,omitempty"`   // Required only for specific death type
}

// ValidationRule represents a validation rule for the death type
type ValidationRule struct {
	RuleName    string  `json:"rule_name"`           // Rule name
	RuleType    string  `json:"rule_type"`           // FIELD, DOCUMENT, BUSINESS
	Description string  `json:"description"`         // Rule description
	Severity    string  `json:"severity"`            // ERROR, WARNING, INFO
	Condition   *string `json:"condition,omitempty"` // Condition for applying rule
	Message     string  `json:"message"`             // Validation message
}

// ==================== HELPER FUNCTIONS ====================

// NewPANValidationResponse creates a new PAN validation response
func NewPANValidationResponse(valid bool, panNumber string, nameOnPAN, panStatus, panType *string) *PANValidationResponse {
	return &PANValidationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: PANValidationData{
			Valid:     valid,
			PANNumber: panNumber,
			NameOnPAN: nameOnPAN,
			PANStatus: panStatus,
			PANType:   panType,
		},
	}
}

// NewBankAccountValidationResponse creates a new bank account validation response
func NewBankAccountValidationResponse(valid bool, accountNumber, ifsc string, accountHolderName, bankName, branchName, accountType, accountStatus *string, nameMatchPercentage *int, validationMethod string) *BankAccountValidationResponse {
	return &BankAccountValidationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: BankAccountValidationData{
			Valid:               valid,
			BankAccountNumber:   accountNumber,
			BankIFSC:            ifsc,
			AccountHolderName:   accountHolderName,
			BankName:            bankName,
			BranchName:          branchName,
			AccountType:         accountType,
			AccountStatus:       accountStatus,
			NameMatchPercentage: nameMatchPercentage,
			ValidationMethod:    validationMethod,
		},
	}
}

// NewDeathDateValidationResponse creates a new death date validation response
func NewDeathDateValidationResponse(valid bool, deathDate, policyID string, investigationRequired bool, validationMessages []string) *DeathDateValidationResponse {
	return &DeathDateValidationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: DeathDateValidationData{
			Valid:                 valid,
			DeathDate:             deathDate,
			PolicyID:              policyID,
			InvestigationRequired: investigationRequired,
			ValidationMessages:    validationMessages,
		},
	}
}

// NewIFSCValidationResponse creates a new IFSC validation response
func NewIFSCValidationResponse(valid bool, ifscCode string, bankName, branchName, address, city, state, district, pinCode, micrCode *string) *IFSCValidationResponse {
	return &IFSCValidationResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: IFSCValidationData{
			Valid:      valid,
			IFSCCode:   ifscCode,
			BankName:   bankName,
			BranchName: branchName,
			Address:    address,
			City:       city,
			State:      state,
			District:   district,
			PINCode:    pinCode,
			MICRCode:   micrCode,
		},
	}
}

// NewDynamicFormFieldsResponse creates a new dynamic form fields response
func NewDynamicFormFieldsResponse(deathType string, formFields []FormField, documentList []DocumentChecklist, validationRules []ValidationRule) *DynamicFormFieldsResponse {
	return &DynamicFormFieldsResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: DynamicFormFieldsData{
			DeathType:       deathType,
			FormFields:      formFields,
			DocumentList:    documentList,
			ValidationRules: validationRules,
		},
	}
}
