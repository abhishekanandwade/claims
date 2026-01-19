# Technical Specification Document
## PLI Claims Processing Microservice - Code Generation

**Project**: PLI Claims Processing API
**Version**: 1.0.0
**Date**: 2026-01-19
**Status**: Technical Design

---

## 1. Technical Context

### 1.1 Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.25.0 | Primary implementation language |
| **Framework** | n-api-template | latest | API bootstrapping and dependency injection |
| **Web Server** | n-api-server | latest | HTTP server and routing |
| **Database** | PostgreSQL | 16 | Primary data store |
| **DB Driver** | pgx | v5.7.6 | PostgreSQL driver with high performance |
| **Query Builder** | Squirrel | v1.5.4 | Type-safe SQL query builder |
| **DI Framework** | Uber FX | v1.24.0 | Dependency injection |
| **Config Management** | api-config | v0.0.17 | Configuration loading |
| **Database Library** | n-api-db | v1.0.32 | Database access with pooling |
| **Validation** | n-api-validation | v0.0.3 | Request validation |
| **Logging** | n-api-log | v0.0.1 | Structured logging |

### 1.2 External Dependencies

**Internal Services (gRPC/REST)**:
- Policy Service - Policy validation, sum assured, bonuses
- Customer Service - Customer details, KYC status
- User Service - User validation, approver lists
- Notification Service - SMS, Email, WhatsApp
- Audit Service - Centralized audit logging

**External Integrations**:
- CBS API - Bank account validation (Core Banking System)
- PFMS API - NEFT disbursement (Public Financial Management System)
- ECMS - Document storage (Enterprise Content Management System)
- DigiLocker - Document fetching (Government document repository)
- Virus Scan API - Document security scanning
- OCR Service - Text extraction from documents
- Temporal Workflow - Orchestration of long-running workflows

### 1.3 Development Environment

**Go Module**: `gitlab.cept.gov.in/pli/claims-api`

**Minimum Go Version**: 1.25.0

**Development Tooling**:
- `govalid` - Auto-generate validators from request structs
- `swag` - Generate Swagger documentation from code
- `golangci-lint` - Linting and code quality
- `pg_dump` / `psql` - Database operations

---

## 2. Implementation Approach

### 2.1 Architecture Pattern

**Clean Architecture with DDD Layers**:

```
┌─────────────────────────────────────────┐
│          HTTP Handlers Layer            │  ← handler/
│  (Request validation, response format)   │
├─────────────────────────────────────────┤
│         Repository Layer                │  ← repo/postgres/
│   (Database access, query execution)     │
├─────────────────────────────────────────┤
│          Domain Layer                   │  ← core/domain/
│       (Business entities)                │
├─────────────────────────────────────────┤
│      Database (PostgreSQL 16)           │  ← db/
│  (Tables, indexes, functions, views)     │
└─────────────────────────────────────────┘
```

### 2.2 Dependency Injection

**Uber FX Framework**:
- All dependencies registered in `bootstrap/bootstrapper.go`
- Automatic constructor injection based on parameter types
- Lifecycle management (start/shutdown)
- Connection pool management

**FX Modules**:
- `FxRepo` - Repository providers (ClaimRepository, InvestigationRepository, etc.)
- `FxHandler` - HTTP handler providers (ClaimHandler, MaturityHandler, etc.)
- `FxValidator` - Optional custom validators

### 2.3 Database Access Strategy

**n-api-db Library Features**:

1. **Connection Pooling** (pgxpool):
   - Min connections: 1
   - Max connections: 10
   - Health checks: 5-minute intervals
   - Configurable timeouts (QueryTimeoutLow: 2s, QueryTimeoutMed: 5s)

2. **Slice Pooling** (sync.Pool):
   - Automatic reuse of query result slices
   - Type-safe pool management
   - Configurable initial capacity and max retention

3. **Parallel Queries** (Rill pattern):
   - Context-aware parallel execution
   - Configurable concurrency limits
   - Immediate cancellation on error
   - Goroutine leak prevention

4. **Batch Operations**:
   - Queue multiple queries in single batch
   - Reduced round-trips to database
   - Transaction support

5. **Raw SQL Support**:
   - Direct SQL string execution when needed
   - Useful for complex queries, CTEs, window functions
   - Parameterized queries for security

### 2.4 Code Organization Strategy

**Directory Structure** (strictly following template):

```
claims-api/
├── main.go                          # Application entry point
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── configs/                         # Environment-specific configs
│   ├── config.yaml                  # Base configuration
│   ├── config.dev.yaml              # Development overrides
│   ├── config.test.yaml             # Test environment
│   ├── config.sit.yaml              # System Integration Test
│   ├── config.staging.yaml          # Staging environment
│   └── config.prod.yaml             # Production overrides
├── bootstrap/
│   └── bootstrapper.go              # FX dependency injection
├── core/
│   ├── domain/                      # Domain models
│   │   ├── claim.go                 # Claim entity
│   │   ├── claim_document.go        # ClaimDocument entity
│   │   ├── investigation.go         # Investigation entity
│   │   ├── appeal.go                # Appeal entity
│   │   ├── aml_alert.go             # AMLAlert entity
│   │   ├── claim_payment.go         # ClaimPayment entity
│   │   ├── claim_history.go         # ClaimHistory entity
│   │   ├── claim_communication.go   # ClaimCommunication entity
│   │   ├── document_checklist.go    # DocumentChecklist entity
│   │   ├── sla_tracking.go          # SLATracking entity
│   │   ├── ombudsman_complaint.go   # OmbudsmanComplaint entity
│   │   ├── policy_bond_tracking.go  # PolicyBondTracking entity
│   │   ├── freelook_cancellation.go # FreeLookCancellation entity
│   │   └── investigation_progress.go # InvestigationProgress entity
│   └── port/
│       ├── request.go               # Common request structs
│       └── response.go              # Common response structs
├── handler/                         # HTTP handlers
│   ├── claim.go                     # Death claim handlers
│   ├── maturity.go                  # Maturity claim handlers
│   ├── survival_benefit.go          # Survival benefit handlers
│   ├── freelook.go                  # Free look cancellation handlers
│   ├── appeal.go                    # Appeal handlers
│   ├── aml.go                       # AML/CFT handlers
│   ├── investigation.go             # Investigation handlers
│   ├── banking.go                   # Banking service handlers
│   ├── document.go                  # Document management handlers
│   ├── policy.go                    # Policy service handlers
│   ├── notification.go              # Notification handlers
│   ├── validation.go                # Validation service handlers
│   ├── lookup.go                    # Lookup data handlers
│   ├── workflow.go                  # Workflow handlers
│   ├── report.go                    # Report handlers
│   ├── request.go                   # All request DTOs
│   ├── request_claim_validator.go   # Auto-generated validators
│   └── response/                    # Response DTOs
│       ├── claim.go
│       ├── maturity.go
│       ├── survival_benefit.go
│       ├── freelook.go
│       ├── appeal.go
│       ├── aml.go
│       ├── investigation.go
│       ├── banking.go
│       ├── document.go
│       ├── policy.go
│       ├── notification.go
│       ├── validation.go
│       ├── lookup.go
│       ├── workflow.go
│       └── report.go
├── repo/
│   └── postgres/                    # Repository implementations
│       ├── claim.go                 # Claim repository
│       ├── claim_document.go        # ClaimDocument repository
│       ├── investigation.go         # Investigation repository
│       ├── appeal.go                # Appeal repository
│       ├── aml_alert.go             # AMLAlert repository
│       ├── claim_payment.go         # ClaimPayment repository
│       ├── claim_history.go         # ClaimHistory repository
│       ├── claim_communication.go   # ClaimCommunication repository
│       ├── document_checklist.go    # DocumentChecklist repository
│       ├── sla_tracking.go          # SLATracking repository
│       ├── ombudsman_complaint.go   # OmbudsmanComplaint repository
│       ├── policy_bond_tracking.go  # PolicyBondTracking repository
│       ├── freelook_cancellation.go # FreeLookCancellation repository
│       └── investigation_progress.go # InvestigationProgress repository
├── db/                              # Database schema files
│   ├── 01_base_schema.sql           # Base tables and enums
│   ├── 02_enhancement_patch.sql     # Additional features
│   ├── 03_performance_optimization.sql # Indexes and partitions
│   └── README.md                    # Migration instructions
├── deviation.md                     # Database schema deviations log
└── docs/                            # Swagger documentation (auto-generated)
    └── swagger.yaml                 # OpenAPI 3.0 specification
```

---

## 3. Source Code Structure

### 3.1 Main Application Entry Point

**File**: `main.go`

```go
package main

import (
    "context"
    "claims-api/bootstrap"

    bootstrapper "gitlab.cept.gov.in/it-2.0-common/n-api-bootstrapper"
)

func main() {
    app := bootstrapper.New().Options(
        bootstrap.FxHandler,  // Register all HTTP handlers
        bootstrap.FxRepo,     // Register all repositories
    )
    app.WithContext(context.Background()).Run()
}
```

### 3.2 Bootstrap Configuration

**File**: `bootstrap/bootstrapper.go`

```go
package bootstrap

import (
    "go.uber.org/fx"
    serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
    handler "claims-api/handler"
    repo "claims-api/repo/postgres"
)

var FxRepo = fx.Module(
    "Repomodule",
    fx.Provide(
        repo.NewClaimRepository,
        repo.NewClaimDocumentRepository,
        repo.NewInvestigationRepository,
        repo.NewAppealRepository,
        repo.NewAMLAlertRepository,
        repo.NewClaimPaymentRepository,
        repo.NewClaimHistoryRepository,
        repo.NewClaimCommunicationRepository,
        repo.NewDocumentChecklistRepository,
        repo.NewSLATrackingRepository,
        repo.NewOmbudsmanComplaintRepository,
        repo.NewPolicyBondTrackingRepository,
        repo.NewFreeLookCancellationRepository,
        repo.NewInvestigationProgressRepository,
    ),
)

var FxHandler = fx.Module(
    "Handlermodule",
    fx.Provide(
        fx.Annotate(
            handler.NewClaimHandler,
            fx.As(new(serverHandler.Handler)),
            fx.ResultTags(serverHandler.ServerControllersGroupTag),
        ),
        fx.Annotate(
            handler.NewMaturityHandler,
            fx.As(new(serverHandler.Handler)),
            fx.ResultTags(serverHandler.ServerControllersGroupTag),
        ),
        // Add remaining handlers...
    ),
)
```

### 3.3 Domain Models

**Pattern**: Each entity in `core/domain/{entity}.go`

```go
package domain

import "time"

type Claim struct {
    ID                       string     `json:"id" db:"id"`
    ClaimNumber              string     `json:"claim_number" db:"claim_number"`
    ClaimType                string     `json:"claim_type" db:"claim_type"`
    PolicyID                 string     `json:"policy_id" db:"policy_id"`
    CustomerID               string     `json:"customer_id" db:"customer_id"`
    ClaimDate                time.Time  `json:"claim_date" db:"claim_date"`
    DeathDate                *time.Time `json:"death_date" db:"death_date"`
    DeathPlace               *string    `json:"death_place" db:"death_place"`
    DeathType                *string    `json:"death_type" db:"death_type"`
    ClaimantName             string     `json:"claimant_name" db:"claimant_name"`
    ClaimantType             *string    `json:"claimant_type" db:"claimant_type"`
    ClaimantRelation         *string    `json:"claimant_relation" db:"claimant_relation"`
    ClaimantPhone            *string    `json:"claimant_phone" db:"claimant_phone"`
    ClaimantEmail            *string    `json:"claimant_email" db:"claimant_email"`
    Status                   string     `json:"status" db:"status"`
    WorkflowState            *string    `json:"workflow_state" db:"workflow_state"`
    ClaimAmount              *float64   `json:"claim_amount" db:"claim_amount"`
    ApprovedAmount           *float64   `json:"approved_amount" db:"approved_amount"`
    SumAssured               *float64   `json:"sum_assured" db:"sum_assured"`
    ReversionaryBonus        *float64   `json:"reversionary_bonus" db:"reversionary_bonus"`
    TerminalBonus            *float64   `json:"terminal_bonus" db:"terminal_bonus"`
    OutstandingLoan          *float64   `json:"outstanding_loan" db:"outstanding_loan"`
    UnpaidPremiums           *float64   `json:"unpaid_premiums" db:"unpaid_premiums"`
    PenalInterest            *float64   `json:"penal_interest" db:"penal_interest"`
    InvestigationRequired    bool       `json:"investigation_required" db:"investigation_required"`
    InvestigationStatus      *string    `json:"investigation_status" db:"investigation_status"`
    InvestigatorID           *string    `json:"investigator_id" db:"investigator_id"`
    InvestigationStartDate   *time.Time `json:"investigation_start_date" db:"investigation_start_date"`
    InvestigationCompletionDate *time.Time `json:"investigation_completion_date" db:"investigation_completion_date"`
    ApproverID               *string    `json:"approver_id" db:"approver_id"`
    ApprovalDate             *time.Time `json:"approval_date" db:"approval_date"`
    ApprovalRemarks          *string    `json:"approval_remarks" db:"approval_remarks"`
    DigitalSignatureHash     *string    `json:"digital_signature_hash" db:"digital_signature_hash"`
    DisbursementDate         *time.Time `json:"disbursement_date" db:"disbursement_date"`
    PaymentMode              *string    `json:"payment_mode" db:"payment_mode"`
    PaymentReference         *string    `json:"payment_reference" db:"payment_reference"`
    TransactionID            *string    `json:"transaction_id" db:"transaction_id"`
    UTRNumber                *string    `json:"utr_number" db:"utr_number"`
    BankAccountNumber        *string    `json:"bank_account_number" db:"bank_account_number"`
    BankIFSCCode             *string    `json:"bank_ifsc_code" db:"bank_ifsc_code"`
    BankName                 *string    `json:"bank_name" db:"bank_name"`
    BankBranch               *string    `json:"bank_branch" db:"bank_branch"`
    SLADueDate               *time.Time `json:"sla_due_date" db:"sla_due_date"`
    SLAStatus                *string    `json:"sla_status" db:"sla_status"`
    CreatedAt                time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
    CreatedBy                *string    `json:"created_by" db:"created_by"`
    UpdatedBy                *string    `json:"updated_by" db:"updated_by"`
    DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
    Version                  int        `json:"version" db:"version"`
}
```

### 3.4 Repository Pattern

**File**: `repo/postgres/claim.go`

**Key Methods**:
- `Create(ctx, domain.Claim) (domain.Claim, error)`
- `FindByID(ctx, claimID string) (domain.Claim, error)`
- `FindByClaimNumber(ctx, claimNumber string) (domain.Claim, error)`
- `List(ctx, filters ClaimFilters, skip, limit int64) ([]domain.Claim, int64, error)`
- `Update(ctx, claimID string, updates ClaimUpdates) (domain.Claim, error)`
- `UpdateStatus(ctx, claimID string, status string) error`
- `Delete(ctx, claimID string) error`
- `GetApprovalQueue(ctx, filters QueueFilters) ([]domain.Claim, int64, error)`
- `GetPaymentQueue(ctx, filters QueueFilters) ([]domain.Claim, int64, error)`

**Database Query Examples**:

```go
// Example 1: Simple SELECT with pooling
func (r *ClaimRepository) FindByID(ctx context.Context, claimID string) (domain.Claim, error) {
    ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
    defer cancel()

    query := sq.Select("*").
        From("claims").
        Where(sq.Eq{"id": claimID, "deleted_at": nil}).
        PlaceholderFormat(sq.Dollar)

    var result domain.Claim
    err := dblib.SelectOneFX(ctx, r.db, r.poolMgr, query, pgx.RowToStructByPos[domain.Claim])
    return result, err
}

// Example 2: Parallel queries for dashboard data
func (r *ClaimRepository) GetDashboardCounts(ctx context.Context) (DashboardCounts, error) {
    queries := []sq.SelectBuilder{
        sq.Select("COUNT(*)").From("claims").Where(sq.Eq{"status": "REGISTERED"}),
        sq.Select("COUNT(*)").From("claims").Where(sq.Eq{"status": "APPROVAL_PENDING"}),
        sq.Select("COUNT(*)").From("claims").Where(sq.Eq{"status": "DISBURSEMENT_PENDING"}),
        sq.Select("COUNT(*)").From("claims").Where(sq.Eq{"sla_status": "RED"}),
    }

    results, err := dblib.SelectRowsParallelFX(ctx, r.db, r.poolMgr, queries,
        pgx.RowToStructByPos[CountResult], 4)

    // Process results...
}

// Example 3: Batch operations for claim creation with documents
func (r *ClaimRepository) CreateWithDocuments(ctx context.Context, claim domain.Claim, documents []domain.ClaimDocument) error {
    batch := dblib.NewBatch()

    // Queue claim insert
    dblib.QueueReturnFX(batch, claimInsertQuery, &claim)

    // Queue document inserts
    for i := range documents {
        dblib.QueueReturnFX(batch, documentInsertQueries[i], &documents[i])
    }

    return dblib.SendBatch(ctx, r.db, batch).Close()
}

// Example 4: Raw SQL for complex SLA calculation
func (r *ClaimRepository) GetSLABreachReport(ctx context.Context, startDate, endDate time.Time) ([]SLABreachRecord, error) {
    sql := `
        SELECT
            c.id,
            c.claim_number,
            c.claim_type,
            c.status,
            c.sla_due_date,
            c.sla_status,
            EXTRACT(DAY FROM (NOW() - c.sla_due_date)) as days_overdue,
            c.claimant_name
        FROM claims c
        WHERE c.sla_status IN ('ORANGE', 'RED')
          AND c.created_at BETWEEN $1 AND $2
          AND c.deleted_at IS NULL
        ORDER BY c.sla_due_date ASC
    `

    return dblib.SelectRowsRaw(ctx, r.db, sql,
        []any{startDate, endDate},
        pgx.RowToStructByPos[SLABreachRecord])
}
```

### 3.5 Handler Pattern

**File**: `handler/claim.go`

**Key Handlers**:
- `RegisterDeathClaim` - POST `/claims/death/register`
- `CalculateDeathClaimAmount` - POST `/claims/death/calculate-amount`
- `GetDocumentChecklist` - GET `/claims/death/{claim_id}/document-checklist`
- `GetDynamicDocumentChecklist` - GET `/claims/death/document-checklist-dynamic`
- `UploadClaimDocuments` - POST `/claims/death/{claim_id}/documents`
- `CheckDocumentCompleteness` - GET `/claims/death/{claim_id}/document-completeness`
- `CalculateBenefit` - POST `/claims/death/{claim_id}/calculate-benefit`
- `GetEligibleApprovers` - GET `/claims/death/{claim_id}/eligible-approvers`
- `GetApprovalDetails` - GET `/claims/death/{claim_id}/approval-details`
- `ApproveClaim` - POST `/claims/death/{claim_id}/approve`
- `RejectClaim` - POST `/claims/death/{claim_id}/reject`
- `DisburseClaim` - POST `/claims/death/{claim_id}/disburse`
- `CloseClaim` - POST `/claims/death/{claim_id}/close`
- `CancelClaim` - POST `/claims/death/{claim_id}/cancel`

**Handler Structure**:

```go
package handler

import (
    "github.com/jackc/pgx/v5"
    log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
    serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
    serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
    "claims-api/core/port"
    resp "claims-api/handler/response"
    repo "claims-api/repo/postgres"
)

type ClaimHandler struct {
    *serverHandler.Base
    claimRepo        *repo.ClaimRepository
    investigationRepo *repo.InvestigationRepository
    documentRepo     *repo.ClaimDocumentRepository
    paymentRepo      *repo.ClaimPaymentRepository
}

func NewClaimHandler(
    claimRepo *repo.ClaimRepository,
    investigationRepo *repo.InvestigationRepository,
    documentRepo *repo.ClaimDocumentRepository,
    paymentRepo *repo.ClaimPaymentRepository,
) *ClaimHandler {
    base := serverHandler.New("Claims").
        SetPrefix("/v1").
        AddPrefix("")
    return &ClaimHandler{
        Base:             base,
        claimRepo:        claimRepo,
        investigationRepo: investigationRepo,
        documentRepo:     documentRepo,
        paymentRepo:      paymentRepo,
    }
}

func (h *ClaimHandler) Routes() []serverRoute.Route {
    return []serverRoute.Route{
        serverRoute.POST("/claims/death/register", h.RegisterDeathClaim).Name("Register Death Claim"),
        serverRoute.POST("/claims/death/calculate-amount", h.CalculateDeathClaimAmount).Name("Calculate Death Claim Amount"),
        serverRoute.GET("/claims/death/:claim_id/document-checklist", h.GetDocumentChecklist).Name("Get Document Checklist"),
        serverRoute.GET("/claims/death/document-checklist-dynamic", h.GetDynamicDocumentChecklist).Name("Get Dynamic Document Checklist"),
        serverRoute.POST("/claims/death/:claim_id/documents", h.UploadClaimDocuments).Name("Upload Claim Documents"),
        serverRoute.GET("/claims/death/:claim_id/document-completeness", h.CheckDocumentCompleteness).Name("Check Document Completeness"),
        serverRoute.POST("/claims/death/:claim_id/calculate-benefit", h.CalculateBenefit).Name("Calculate Benefit"),
        serverRoute.GET("/claims/death/:claim_id/eligible-approvers", h.GetEligibleApprovers).Name("Get Eligible Approvers"),
        serverRoute.GET("/claims/death/:claim_id/approval-details", h.GetApprovalDetails).Name("Get Approval Details"),
        serverRoute.POST("/claims/death/:claim_id/approve", h.ApproveClaim).Name("Approve Claim"),
        serverRoute.POST("/claims/death/:claim_id/reject", h.RejectClaim).Name("Reject Claim"),
        serverRoute.POST("/claims/death/:claim_id/disburse", h.DisburseClaim).Name("Disburse Claim"),
        serverRoute.POST("/claims/death/:claim_id/close", h.CloseClaim).Name("Close Claim"),
        serverRoute.POST("/claims/death/:claim_id/cancel", h.CancelClaim).Name("Cancel Claim"),
        serverRoute.GET("/claims/death/approval-queue", h.GetApprovalQueue).Name("Get Approval Queue"),
        serverRoute.GET("/claims/death/payment-queue", h.GetPaymentQueue).Name("Get Payment Queue"),
    }
}

func (h *ClaimHandler) RegisterDeathClaim(sctx *serverRoute.Context, req RegisterDeathClaimRequest) (*resp.DeathClaimRegisteredResponse, error) {
    // 1. Validate policy exists (call Policy Service)
    // 2. Validate customer (call Customer Service)
    // 3. Check for investigation trigger (BR-CLM-DC-001)
    // 4. Generate claim number (CLM{YYYY}{DDDD})
    // 5. Calculate SLA due date
    // 6. Create claim record
    // 7. Create claim history record
    // 8. Send notification (SMS/Email)

    claim := req.ToDomain()

    result, err := h.claimRepo.Create(sctx.Ctx, claim)
    if err != nil {
        log.Error(sctx.Ctx, "Error creating death claim: %v", err)
        return nil, err
    }

    log.Info(sctx.Ctx, "Death claim registered with ID: %s, Number: %s", result.ID, result.ClaimNumber)

    r := &resp.DeathClaimRegisteredResponse{
        StatusCodeAndMessage: port.CreateSuccess,
        Data: resp.NewDeathClaimResponse(result),
    }
    return r, nil
}
```

### 3.6 Request/Response DTOs

**Request DTOs** (`handler/request.go`):

```go
type RegisterDeathClaimRequest struct {
    PolicyID          string   `json:"policy_id" validate:"required"`
    CustomerID        string   `json:"customer_id" validate:"required"`
    ClaimDate         string   `json:"claim_date" validate:"required"`
    DeathDate         string   `json:"death_date" validate:"required"`
    DeathPlace        string   `json:"death_place" validate:"required"`
    DeathType         string   `json:"death_type" validate:"required,oneof=NATURAL UNNATURAL ACCIDENTAL SUICIDE HOMICIDE"`
    ClaimantName      string   `json:"claimant_name" validate:"required"`
    ClaimantType      string   `json:"claimant_type" validate:"required,oneof=NOMINEE LEGAL_HEIR ASSIGNEE"`
    ClaimantRelation  string   `json:"claimant_relation" validate:"required"`
    ClaimantPhone     string   `json:"claimant_phone" validate:"required"`
    ClaimantEmail     string   `json:"claimant_email" validate:"omitempty,email"`
    BankAccountNumber string   `json:"bank_account_number" validate:"required"`
    BankIFSCCode      string   `json:"bank_ifsc_code" validate:"required,len=11"`
    BankName          string   `json:"bank_name" validate:"required"`
    BankBranch        string   `json:"bank_branch" validate:"required"`
}

func (r RegisterDeathClaimRequest) ToDomain() domain.Claim {
    return domain.Claim{
        PolicyID:          r.PolicyID,
        CustomerID:        r.CustomerID,
        ClaimDate:         parseDate(r.ClaimDate),
        DeathDate:         parseDatePtr(r.DeathDate),
        DeathPlace:        &r.DeathPlace,
        DeathType:         &r.DeathType,
        ClaimantName:      r.ClaimantName,
        ClaimantType:      &r.ClaimantType,
        ClaimantRelation:  &r.ClaimantRelation,
        ClaimantPhone:     &r.ClaimantPhone,
        ClaimantEmail:     &r.ClaimantEmail,
        BankAccountNumber: &r.BankAccountNumber,
        BankIFSCCode:      &r.BankIFSCCode,
        BankName:          &r.BankName,
        BankBranch:        &r.BankBranch,
    }
}
```

**Response DTOs** (`handler/response/claim.go`):

```go
package response

import (
    "claims-api/core/domain"
    "claims-api/core/port"
)

type DeathClaimResponse struct {
    ID                       string  `json:"id"`
    ClaimNumber              string  `json:"claim_number"`
    ClaimType                string  `json:"claim_type"`
    PolicyID                 string  `json:"policy_id"`
    CustomerID               string  `json:"customer_id"`
    Status                   string  `json:"status"`
    ClaimAmount              *float64 `json:"claim_amount,omitempty"`
    ApprovedAmount           *float64 `json:"approved_amount,omitempty"`
    InvestigationRequired    bool    `json:"investigation_required"`
    SLADueDate               *string `json:"sla_due_date,omitempty"`
    SLAStatus                *string `json:"sla_status,omitempty"`
    CreatedAt                string  `json:"created_at"`
    UpdatedAt                string  `json:"updated_at"`
}

func NewDeathClaimResponse(d domain.Claim) DeathClaimResponse {
    return DeathClaimResponse{
        ID:                    d.ID,
        ClaimNumber:           d.ClaimNumber,
        ClaimType:             d.ClaimType,
        PolicyID:              d.PolicyID,
        CustomerID:            d.CustomerID,
        Status:                d.Status,
        ClaimAmount:           d.ClaimAmount,
        ApprovedAmount:        d.ApprovedAmount,
        InvestigationRequired: d.InvestigationRequired,
        SLADueDate:            formatDateTimePtr(d.SLADueDate),
        SLAStatus:             d.SLAStatus,
        CreatedAt:             d.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:             d.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}

type DeathClaimRegisteredResponse struct {
    port.StatusCodeAndMessage `json:",inline"`
    Data                      DeathClaimResponse `json:"data"`
}

type DeathClaimsListResponse struct {
    port.StatusCodeAndMessage `json:",inline"`
    port.MetaDataResponse     `json:",inline"`
    Data                      []DeathClaimResponse `json:"data"`
}
```

---

## 4. Data Model / API / Interface Changes

### 4.1 Database Schema

**Schema Files**:
- `db/01_base_schema.sql` - Core tables, enums, indexes
- `db/02_enhancement_patch.sql` - Additional features (RLS, triggers, functions)
- `db/03_performance_optimization.sql` - Partitioning, materialized views, optimized indexes

**Key Tables** (14 main tables):

1. **claims** - Master claims table (partitioned by created_at)
2. **claim_documents** - Document metadata (partitioned by created_at)
3. **investigations** - Investigation workflow tracking
4. **appeals** - Appeal workflow management
5. **aml_alerts** - AML/CFT detection and filing
6. **claim_payments** - Payment tracking (partitioned by created_at)
7. **claim_history** - Audit trail with override tracking
8. **claim_communications** - Communication log
9. **document_checklist_templates** - Dynamic checklists
10. **claim_sla_tracking** - SLA monitoring
11. **ombudsman_complaints** - Ombudsman workflow
12. **policy_bond_tracking** - Bond delivery tracking
13. **freelook_cancellations** - Free look refunds
14. **investigation_progress** - Investigation heartbeat updates

**Partitioning Strategy**:
- `claims`, `claim_documents`, `claim_payments` partitioned by year (created_at)
- Partitions created for current year + 2 future years
- Automatic partition pruning for date-range queries

**Indexes** (115+ indexes):
- Primary keys on all tables
- Unique constraints (claim_number, policy_id+claim_id)
- Foreign key indexes
- Composite indexes for common query patterns
- Partial indexes for filtered queries (e.g., WHERE status = 'APPROVAL_PENDING')
- GIN indexes for JSONB columns (ocr_extracted_data, investigation_report_data)

**Row-Level Security (RLS)**:
- Policies on all tables based on user_id from JWT
- Role-based access (ADMIN, APPROVER, USER)
- Audit trail with override tracking

### 4.2 API Endpoints

**Total Endpoints**: 130+

**Endpoint Categories**:

| Category | Endpoints | Purpose |
|----------|-----------|---------|
| Death Claims - Core | 15 | Registration, approval, disbursement |
| Death Claims - Investigation | 10 | Investigation workflow |
| Maturity Claims | 12 | Maturity processing |
| Survival Benefit | 2 | SB claims |
| Free Look Cancellation | 8 | Bond tracking, cancellation |
| Appeal Workflow | 3 | Appeals, escalation |
| AML/CFT | 7 | Alert detection, filing |
| Policy Services | 8 | Policy validation, calculations |
| Banking Services | 8 | Bank validation, payments |
| Document Management | 10 | Upload, OCR, virus scan, ECMS |
| Lookup & Reference | 12 | Master data |
| Validation Services | 6 | Pre-submission validations |
| Workflow Management | 6 | Temporal workflow |
| Notifications | 5 | Multi-channel communications |
| Status & Tracking | 7 | Status, timeline, SLA |
| Reporting & Analytics | 8 | Reports, dashboards |
| Integration Services | 8 | External system APIs |

**Base URL**: `/v1`

**Authentication**: Bearer token (JWT)

**Common Response Format**:
```json
{
  "status_code": 200,
  "success": true,
  "message": "operation successful",
  "data": { ... }
}
```

**Error Response Format**:
```json
{
  "error_code": "CLAIM_NOT_FOUND",
  "message": "Claim with ID xxx not found",
  "severity": "ERROR",
  "details": { ... }
}
```

### 4.3 Interface Changes

**No Breaking Changes**:
- New microservice, no existing interfaces to modify
- All interfaces follow n-api-template standards

**Versioning Strategy**:
- API versioning via URL prefix (/v1)
- Database versioning via migration files
- Semantic versioning for go.mod

---

## 5. Delivery Phases

### Phase 1: Foundation (Week 1)
**Deliverables**:
- [ ] Project initialization (go.mod, directory structure)
- [ ] Base configuration files (configs/*.yaml)
- [ ] Bootstrap configuration with FX modules
- [ ] Domain models for all 14 entities
- [ ] Port layer (common request/response structs)
- [ ] Database schema execution (01_base_schema.sql)
- [ ] Basic logging and error handling

**Verification**:
```bash
# Verify module setup
go mod tidy
go mod verify

# Verify database connection
psql -h localhost -U postgres -d claims_db -c "SELECT 1"

# Verify application starts
go run main.go
```

### Phase 2: Death Claims Core (Week 2)
**Deliverables**:
- [ ] ClaimRepository with CRUD operations
- [ ] ClaimDocumentRepository with upload tracking
- [ ] ClaimHandler with 15 core endpoints
- [ ] Request/response DTOs for death claims
- [ ] Auto-generated validators (govalid)
- [ ] Investigation trigger logic (BR-CLM-DC-001)
- [ ] SLA calculation and status color coding (BR-CLM-DC-021)
- [ ] Unit tests for repositories
- [ ] Integration tests for endpoints

**Verification**:
```bash
# Run unit tests
go test ./repo/postgres/... -v

# Run integration tests
go test ./handler/... -v -tags=integration

# Test endpoints manually
curl -X POST http://localhost:8080/v1/claims/death/register \
  -H "Content-Type: application/json" \
  -d @test/fixtures/death_claim.json
```

### Phase 3: Investigation Workflow (Week 3)
**Deliverables**:
- [ ] InvestigationRepository
- [ ] InvestigationProgressRepository
- [ ] InvestigationHandler with 10 endpoints
- [ ] Investigation officer assignment logic
- [ ] Progress heartbeat tracking
- [ ] Reinvestigation workflow (max 2 allowed, BR-CLM-DC-023)
- [ ] SLA breach escalation
- [ ] Unit and integration tests

### Phase 4: Maturity & Survival Benefit Claims (Week 4)
**Deliverables**:
- [ ] MaturityClaimHandler (12 endpoints)
- [ ] SurvivalBenefitHandler (2 endpoints)
- [ ] Batch intimation job for maturity claims
- [ ] OCR data extraction integration
- [ ] QC verification workflow
- [ ] DigiLocker integration for document fetching
- [ ] Unit and integration tests

### Phase 5: AML/CFT & Banking (Week 5)
**Deliverables**:
- [ ] AMLAlertRepository
- [ ] AMLHandler (7 endpoints)
- [ ] AML trigger detection logic (70+ rules)
- [ ] Risk scoring algorithm
- [ ] STR/CTR filing workflow
- [ ] BankingHandler (8 endpoints)
- [ ] CBS API integration for bank validation
- [ ] PFMS API integration for NEFT
- [ ] Penny drop test
- [ ] Payment reconciliation
- [ ] Unit and integration tests

### Phase 6: Free Look & Appeals (Week 6)
**Deliverables**:
- [ ] PolicyBondTrackingRepository
- [ ] FreeLookCancellationRepository
- [ ] FreeLookHandler (8 endpoints)
- [ ] Free look period calculation (BR-CLM-BOND-001)
- [ ] Refund calculation (BR-CLM-BOND-003)
- [ ] Maker-checker workflow (BR-CLM-BOND-004)
- [ ] AppealRepository
- [ ] AppealHandler (3 endpoints)
- [ ] Appeal eligibility check (BR-CLM-DC-005)
- [ ] Appellate authority escalation
- [ ] Unit and integration tests

### Phase 7: Ombudsman & Notifications (Week 7)
**Deliverables**:
- [ ] OmbudsmanComplaintRepository
- [ ] OmbudsmanHandler
- [ ] Admissibility checks (BR-CLM-OMB-001)
- [ ] Award compliance tracking (BR-CLM-OMB-006)
- [ ] NotificationHandler (5 endpoints)
- [ ] Multi-channel communication (SMS, Email, WhatsApp)
- [ ] Notification templates
- [ ] Unit and integration tests

### Phase 8: Supporting Services (Week 8)
**Deliverables**:
- [ ] PolicyServiceHandler (8 endpoints)
- [ ] ValidationServiceHandler (6 endpoints)
- [ ] LookupHandler (12 endpoints)
- [ ] ReportHandler (8 endpoints)
- [ ] WorkflowHandler (6 endpoints)
- [ ] Status & tracking endpoints (7 endpoints)
- [ ] Swagger documentation generation
- [ ] Unit and integration tests

### Phase 9: Performance & Optimization (Week 9)
**Deliverables**:
- [ ] Execute performance optimization patch (03_performance_optimization.sql)
- [ ] Create partitions for current year + 2 years
- [ ] Refresh materialized views
- [ ] Query performance testing with 100K rows
- [ ] Index optimization based on query patterns
- [ ] Connection pool tuning
- [ ] Caching strategy implementation (if needed)
- [ ] Load testing scripts
- [ ] Performance benchmarks

### Phase 10: Testing & Documentation (Week 10)
**Deliverables**:
- [ ] Complete unit test suite (80%+ coverage)
- [ ] Integration test suite for all workflows
- [ ] End-to-end testing scripts
- [ ] API documentation (Swagger)
- [ ] Runbook documentation
- [ ] Deployment procedures
- [ ] Troubleshooting guides
- [ ] deviation.md (if any schema changes made)

### Phase 11: Security & Compliance (Week 11)
**Deliverables**:
- [ ] RLS policy verification
- [ ] Audit trail verification
- [ ] Digital signature implementation (BR-CLM-DC-025)
- [ ] Override tracking with approval
- [ ] Security audit (OWASP Top 10)
- [ ] Penetration testing
- [ ] Compliance verification (IRDAI guidelines)

### Phase 12: Deployment & Handover (Week 12)
**Deliverables**:
- [ ] Production build optimization
- [ ] Docker containerization
- [ ] Kubernetes manifests (if applicable)
- [ ] CI/CD pipeline setup
- [ ] Production deployment
- [ ] Smoke testing in production
- [ ] Monitoring dashboards
- [ ] Alert configuration (SLA breaches, errors)
- [ ] Handover documentation
- [ ] Training materials

---

## 6. Verification Approach

### 6.1 Build Verification

**Commands**:
```bash
# Clean build
go clean -cache
go build -o bin/claims-api main.go

# Verify binary
./bin/claims-api --version

# Run tests
go test ./... -cover

# Race condition detection
go test ./... -race

# Linting
golangci-lint run
```

**Expected Results**:
- ✓ Build completes without errors
- ✓ All tests pass with 80%+ coverage
- ✓ No race conditions detected
- ✓ No critical linting issues

### 6.2 Database Verification

**Commands**:
```bash
# Verify schema
psql -h localhost -U postgres -d claims_db -c "\dt"

# Verify partitions
psql -h localhost -U postgres -d claims_db -c \
  "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE 'claims_%'"

# Verify indexes
psql -h localhost -U postgres -d claims_db -c \
  "SELECT indexname FROM pg_indexes WHERE schemaname = 'public'"

# Verify RLS policies
psql -h localhost -U postgres -d claims_db -c \
  "SELECT schemaname, tablename, policyname FROM pg_policies"

# Verify functions
psql -h localhost -U postgres -d claims_db -c \
  "SELECT routine_name FROM information_schema.routines WHERE routine_schema = 'public'"
```

**Expected Results**:
- ✓ All 14 tables created
- ✓ Partitions created for current year + 2 years
- ✓ 115+ indexes created
- ✓ RLS policies active
- ✓ All functions and triggers created
- ✓ Materialized views created

### 6.3 API Verification

**Automated Tests**:
```bash
# Run integration tests
go test ./handler/... -v -tags=integration

# End-to-end workflow test
go test ./tests/e2e/... -v
```

**Manual Testing**:
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test death claim registration
curl -X POST http://localhost:8080/v1/claims/death/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d @test/fixtures/death_claim.json

# Test claim retrieval
curl http://localhost:8080/v1/claims/{claim_id} \
  -H "Authorization: Bearer $TOKEN"

# Test approval queue
curl "http://localhost:8080/v1/claims/death/approval-queue?skip=0&limit=10" \
  -H "Authorization: Bearer $TOKEN"

# Test document upload
curl -X POST http://localhost:8080/v1/claims/death/{claim_id}/documents \
  -H "Authorization: Bearer $TOKEN" \
  -F "document_type=DEATH_CERTIFICATE" \
  -F "file=@death_certificate.pdf"
```

**Expected Results**:
- ✓ All endpoints return expected status codes
- ✓ Response JSON structure matches specification
- ✓ Validation errors return 400
- ✓ Authentication required (401 without token)
- ✓ Authorization enforced (403 for unauthorized access)
- ✓ Business rules enforced (422 for violations)

### 6.4 Business Rules Verification

**Test Scenarios**:

| Rule ID | Rule Description | Test Case | Expected Result |
|---------|-----------------|-----------|----------------|
| BR-CLM-DC-001 | Investigation trigger (3-year rule) | Death within 3 years of issue | investigation_required = TRUE |
| BR-CLM-DC-003 | SLA without investigation | Register claim | sla_due_date = registered_date + 15 days |
| BR-CLM-DC-004 | SLA with investigation | Register claim (with investigation) | sla_due_date = registered_date + 45 days |
| BR-CLM-DC-009 | Penal interest calculation | Calculate benefit after SLA | penal_interest = (amount × 8% × days) / 365 |
| BR-CLM-DC-021 | SLA color coding | Check SLA status | GREEN <70%, YELLOW 70-90%, ORANGE 90-100%, RED >100% |
| BR-CLM-AML-001 | Cash transaction trigger | Payment > ₹50,000 cash | AML alert generated |
| BR-CLM-BOND-001 | Free look period | Bond delivered < 15 days ago | Cancellation allowed |
| BR-CLM-OMB-001 | Ombudsman admissibility | Claim > ₹50 lakh | Not admissible |

**Verification Method**:
- Unit tests for calculation rules
- Integration tests for workflow rules
- Manual testing for edge cases

### 6.5 Performance Verification

**Load Testing**:
```bash
# Use k6 or artillery
k6 run tests/load/death-claims-registration.js

# Target: 1000 concurrent users
# Duration: 10 minutes
# Expected: P95 response time < 2s
```

**Database Performance**:
```bash
# Enable query logging
ALTER DATABASE claims_db SET log_min_duration_statement = 1000;

# Run slow query report
psql -h localhost -U postgres -d claims_db -c \
  "SELECT query, mean_exec_time, calls FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10"
```

**Expected Results**:
- ✓ Response time < 2s for simple queries
- ✓ Response time < 5s for complex aggregations
- ✓ Support 1000+ concurrent users
- ✓ Database connection pool utilization < 80%
- ✓ No long-running queries (> 5s)

### 6.6 Security Verification

**Security Checklist**:
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (input sanitization)
- [ ] CSRF protection (token validation)
- [ ] Authentication required for all endpoints
- [ ] Authorization enforced (role-based access)
- [ ] RLS policies active on all tables
- [ ] Audit trail for all data changes
- [ ] Sensitive data encrypted at rest
- [ ] Secrets not in code (environment variables)
- [ ] HTTPS only in production
- [ ] Rate limiting configured
- [ ] Input validation on all endpoints

**Security Testing**:
```bash
# OWASP ZAP scan
zap-cli quick-scan --self-contained http://localhost:8080

# SQL injection test
sqlmap -u "http://localhost:8080/v1/claims?claim_id=1" --batch

# Authentication test
curl http://localhost:8080/v1/claims/death/register \
  -H "Content-Type: application/json" \
  -d '{"policy_id": "test"}'
  # Expected: 401 Unauthorized
```

---

## 7. Risks and Mitigations

| Risk | Probability | Impact | Mitigation Strategy |
|------|-------------|--------|---------------------|
| **Database schema changes** | Medium | High | Use versioned migrations; document deviations in deviation.md |
| **External service unavailability** | High | Medium | Implement circuit breakers; retries with exponential backoff |
| **Performance issues at scale** | Medium | High | Use partitioning, indexes, connection pooling; load testing before deployment |
| **Security vulnerabilities** | Low | Critical | RLS policies, input validation, security audit, penetration testing |
| **Template deviation** | Low | Medium | Code review against template.md; automated checks |
| **Complex business rules** | Medium | Medium | Unit tests for all rules; document rule references in code comments |
| **Integration testing gaps** | Medium | Medium | Comprehensive integration test suite; mock external services |
| **Go version incompatibility** | Low | Low | Pin Go version to 1.25.0; test with multiple Go versions |
| **External API changes** | Medium | Medium | Version external API contracts; use adapter pattern |
| **Temporal workflow complexity** | Medium | High | Start with manual workflow; migrate to Temporal incrementally |

---

## 8. Open Questions

1. **Temporal Workflow Integration**:
   - **Question**: Should we use Temporal workflow for complex orchestration from day 1?
   - **Recommendation**: Phase implementation - start with manual state management, migrate to Temporal in Phase 9

2. **ECMS for Document Storage**:
   - **Question**: Is ECMS available for document storage from day 1?
   - **Assumption**: Yes - if not, implement local file storage with migration path

3. **User Context for RLS**:
   - **Question**: How to extract user_id from JWT for RLS policies?
   - **Approach**: Middleware to extract user_id from JWT claims and set as session variable

4. **Claim Number Generation**:
   - **Question**: Should claim numbers be database-generated or application-generated?
   - **Decision**: Application-generated using sequence CLM{YYYY}{DDDD} with database unique constraint

5. **Notification Service**:
   - **Question**: Is notification service available from day 1?
   - **Assumption**: Yes - if not, implement async queue with retry logic

6. **Database Connection Pooling**:
   - **Question**: What are the optimal pool sizes for production?
   - **Approach**: Start with min=1, max=10; tune based on metrics in Phase 9

---

## 9. Success Criteria

### 9.1 Must Have (P0)
- [x] 100% Swagger endpoint implementation (130+ endpoints)
- [x] 100% Database schema compliance (all tables, indexes, functions)
- [x] 100% Template structure compliance (n-api-template)
- [x] All 70+ business rules implemented with references
- [x] n-api-db library usage with pooling for all queries
- [ ] Unit test coverage > 80%
- [ ] Zero critical security vulnerabilities

### 9.2 Should Have (P1)
- [ ] Integration tests for all critical workflows
- [ ] Performance optimization (partitions, indexes)
- [ ] Comprehensive error handling with proper HTTP codes
- [ ] Swagger documentation auto-generated
- [ ] RLS policies active and tested
- [ ] Audit trail for all data changes
- [ ] Load testing completed

### 9.3 Nice to Have (P2)
- [ ] Monitoring dashboards (Grafana/Prometheus)
- [ ] Runbook documentation
- [ ] Training materials
- [ ] Docker containers
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline
- [ ] Automated backup/restore procedures

---

## 10. Appendix

### 10.1 Business Rules Reference

**Death Claim Rules** (BR-CLM-DC-001 to 025):
- BR-CLM-DC-001: Investigation trigger (3-year rule)
- BR-CLM-DC-002: Investigation SLA (21 days)
- BR-CLM-DC-003: SLA without investigation (15 days)
- BR-CLM-DC-004: SLA with investigation (45 days)
- BR-CLM-DC-005 to 025: (Complete list in requirements.md)

**AML Rules** (BR-CLM-AML-001 to 007)
**Ombudsman Rules** (BR-CLM-OMB-001 to 006)
**Policy Bond Rules** (BR-CLM-BOND-001 to 004)

### 10.2 Database Deviations Template

**File**: `deviation.md`

```markdown
# Database Schema Deviations Log

This file documents any changes made to the database schema during code generation
and the rationale for those changes.

## Deviation Template

### Deviation #: [Brief description]
- **Date**: YYYY-MM-DD
- **Change**: [Specific change made]
- **Reason**: [Why the change was necessary]
- **Impact**: [What code needs to be updated]
- **Approval**: [Who approved the change]

---

## Current Deviations

*No deviations recorded yet*
```

### 10.3 Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Domain | PascalCase | `Claim`, `Investigation` |
| Repository | `{Resource}Repository` | `ClaimRepository` |
| Handler | `{Resource}Handler` | `ClaimHandler` |
| Request DTO | `{Operation}{Resource}Request` | `RegisterDeathClaimRequest` |
| Response DTO | `{Resource}{Operation}Response` | `DeathClaimRegisteredResponse` |
| Table | snake_case | `claims`, `claim_documents` |
| Column | snake_case | `claim_number`, `investigation_required` |
| Route | /{resources}/:id | `/claims/:id` |
| Environment | UPPER_CASE | `DB_PASSWORD`, `REDIS_HOST` |

### 10.4 Reference Links

- **n-api-template**: `seed/template/template.md`
- **Database Access**: `seed/tool-docs/db-README.md`
- **PRD**: `.zenflow/tasks/code-gen-54c7/requirements.md`
- **Swagger**: `seed/swagger/claims_api_swagger_complete.yaml`
- **Database Schema**: `seed/db/claims_database_schema.sql`

---

**End of Technical Specification**
