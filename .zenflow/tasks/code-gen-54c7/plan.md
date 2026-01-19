# Full SDD workflow

## Configuration
- **Artifacts Path**: {@artifacts_path} → `.zenflow/tasks/{task_id}`

---

## Workflow Steps

### [x] Step: Requirements
<!-- chat-id: b8a9c594-263f-42f2-a66e-ccaad0345412 -->

Create a Product Requirements Document (PRD) based on the feature description.

1. Review existing codebase to understand current architecture and patterns
2. Analyze the feature definition and identify unclear aspects
3. Ask the user for clarifications on aspects that significantly impact scope or user experience
4. Make reasonable decisions for minor details based on context and conventions
5. If user can't clarify, make a decision, state the assumption, and continue

Save the PRD to `{@artifacts_path}/requirements.md`.

**Status**: ✅ Completed
**Output**: `.zenflow/tasks/code-gen-54c7/requirements.md`

**Key Deliverables**:
- Comprehensive PRD covering all 130+ API endpoints
- Business rules documentation (70+ rules)
- Database schema compliance requirements
- Integration requirements for 8 external systems
- Non-functional requirements (performance, security, scalability)
- Technical stack specification (Go 1.25.0, PostgreSQL 16, n-api-db)
- Code structure template compliance
- Testing and deployment requirements

**Assumptions Made**:
1. Claim number auto-generation: CLM{YYYY}{DDDD} format
2. Temporal workflow for complex orchestration
3. ECMS for document storage integration
4. JWT-based user context for RLS
5. PostgreSQL 16 with partitioning enabled

### [x] Step: Technical Specification
<!-- chat-id: ec8d45b2-a352-4dc0-b2cb-ad97a7a35160 -->

Create a technical specification based on the PRD in `{@artifacts_path}/requirements.md`.

1. Review existing codebase architecture and identify reusable components
2. Define the implementation approach

Save to `{@artifacts_path}/spec.md` with:
- Technical context (language, dependencies)
- Implementation approach referencing existing code patterns
- Source code structure changes
- Data model / API / interface changes
- Delivery phases (incremental, testable milestones)
- Verification approach using project lint/test commands

**Status**: ✅ Completed
**Output**: `.zenflow/tasks/code-gen-54c7/spec.md`

**Key Deliverables**:
- Complete technical stack specification (Go 1.25.0, PostgreSQL 16, n-api-db)
- Clean architecture with DDD layers (Handler → Repository → DB)
- Detailed code organization following template.md strictly
- Database access patterns with n-api-db (pooling, parallel queries, batch operations)
- 12 delivery phases with testable milestones (12-week implementation)
- Comprehensive verification approach (build, database, API, business rules, performance, security)
- Risk mitigation strategies for 10 identified risks
- Success criteria (P0: 100% implementation, P1: optimization, P2: monitoring)

### [x] Step: Planning
<!-- chat-id: 4c8fcf45-c3a3-4547-9d5c-c30fa04a22d0 -->

Create a detailed implementation plan based on `{@artifacts_path}/spec.md`.

**Status**: ✅ Completed
**Output**: `.zenflow/tasks/code-gen-54c7/plan.md` (this file)

---

## Detailed Implementation Plan

### Overview
This plan breaks down the PLI Claims Processing API implementation into 12 phases with concrete, actionable tasks. Each task includes verification steps and references to the template.md and db-README.md contracts.

**Total Scope**:
- 130+ API endpoints across 14 handlers
- 14 database tables with partitioning
- 70+ business rules
- 8 external service integrations
- 12-week implementation timeline

---

## Phase 1: Project Foundation Setup
**Duration**: Week 1
**Objective**: Establish project structure, dependencies, and database schema

### [x] Task 1.1: Initialize Go Module and Directory Structure
<!-- chat-id: 3103ca21-a4f6-4c43-b69c-96b24a53c5db -->
**Reference**: `seed/template/template.md` - Project Structure section

**Steps**:
1. Create `go.mod` with module `gitlab.cept.gov.in/pli/claims-api`
2. Set Go version to 1.25.0
3. Add core dependencies:
   - n-api-bootstrapper v0.0.14
   - n-api-server v0.0.17
   - n-api-db v1.0.32
   - n-api-log v0.0.1
   - n-api-validation v0.0.3
   - api-config v0.0.17
   - pgx v5.7.6
   - squirrel v1.5.4
   - uber-fx v1.24.0
4. Create directory structure:
   - `bootstrap/`
   - `configs/`
   - `core/domain/`
   - `core/port/`
   - `handler/`
   - `handler/response/`
   - `repo/postgres/`
   - `db/`

**Verification**:
```bash
go mod tidy
go mod verify
# Verify directories exist
ls -la bootstrap configs core handler repo db
```

---

### [x] Task 1.2: Create Configuration Files
<!-- chat-id: 149cd110-5c9c-4c05-a880-159b69f01dd9 -->
**Reference**: `seed/template/template.md` - Configuration Files section

**Steps**:
1. Create `configs/config.yaml` with base configuration:
   - appname: "claims-api"
   - trace: disabled (dev), enabled (prod)
   - cache: redis + local cache settings
   - db: connection pool settings
   - info: name and version for Swagger
2. Create environment-specific configs:
   - `config.dev.yaml`
   - `config.test.yaml`
   - `config.sit.yaml`
   - `config.staging.yaml`
   - `config.prod.yaml`

**Status**: ✅ Completed

**Verification**:
```bash
# Verify config files are valid YAML
# Check all required keys are present
```

---

### [x] Task 1.3: Create Port Layer (Request/Response Interfaces)
<!-- chat-id: 7adb7ee6-8ea7-41da-bc38-c34ea59c308c -->
**Reference**: `seed/template/template.md` - Port Layer section

**Steps**:
1. Create `core/port/request.go`:
   - MetadataRequest struct (pagination, sorting)
2. Create `core/port/response.go`:
   - StatusCodeAndMessage with predefined constants
   - FileResponse struct
   - MetaDataResponse struct
   - Helper functions (NewMetaDataResponse)

**Status**: ✅ Completed

**Verification**:
```bash
# Verify package compiles
go build ./core/port/...
```

---

### [x] Task 1.4: Create Domain Models (14 Entities)
<!-- chat-id: 5fcdde66-6544-4f80-83e4-0afa292a9ea7 -->
**Reference**: `seed/template/template.md` - Domain Model Pattern section
**Reference**: `seed/db/claims_database_schema.sql` - Table definitions

**Steps**:
Create domain models in `core/domain/`:
1. `claim.go` - Claim entity with all fields from claims table
2. `claim_document.go` - ClaimDocument entity
3. `investigation.go` - Investigation entity
4. `appeal.go` - Appeal entity
5. `aml_alert.go` - AMLAlert entity
6. `claim_payment.go` - ClaimPayment entity
7. `claim_history.go` - ClaimHistory entity
8. `claim_communication.go` - ClaimCommunication entity
9. `document_checklist.go` - DocumentChecklist entity
10. `sla_tracking.go` - SLATracking entity
11. `ombudsman_complaint.go` - OmbudsmanComplaint entity
12. `policy_bond_tracking.go` - PolicyBondTracking entity
13. `freelook_cancellation.go` - FreeLookCancellation entity
14. `investigation_progress.go` - InvestigationProgress entity

Each model must:
- Use `json` and `db` tags
- Match database column names exactly
- Include ID, CreatedAt, UpdatedAt, DeletedAt fields
- Use appropriate Go types (string for UUID, *string for optional, time.Time for dates)

**Status**: ✅ Completed

**Verification**:
```bash
go build ./core/domain/...
# Verify all fields match database schema
```

---

### [x] Task 1.5: Create Bootstrap Configuration
<!-- chat-id: 075635b7-e827-4931-bb80-a1e8ebcfbada -->
**Reference**: `seed/template/template.md` - Bootstrap Configuration section

**Steps**:
1. Create `bootstrap/bootstrapper.go`:
   - FxRepo module with all 14 repository providers
   - FxHandler module with all handler providers (placeholders for now)
2. Register repositories (will be implemented in later phases):
   - NewClaimRepository
   - NewClaimDocumentRepository
   - NewInvestigationRepository
   - (and 11 more)

**Status**: ✅ Completed

**Verification**:
```bash
go build ./bootstrap/...
```

**Key Deliverables**:
- Created `bootstrap/bootstrapper.go` with FxRepo and FxHandler modules
- Registered 14 repository providers in FxRepo:
  - Claim, ClaimDocument, ClaimPayment, ClaimHistory, ClaimCommunication
  - Investigation, InvestigationProgress
  - AMLAlert, Appeal, OmbudsmanComplaint
  - PolicyBondTracking, FreeLookCancellation
  - DocumentChecklist, SLATracking
- Registered 16 handler providers in FxHandler (all as placeholders):
  - Claim, Investigation, MaturityClaim, SurvivalBenefit
  - AML, Banking, FreeLook, Appeal
  - Ombudsman, Notification, PolicyService, ValidationService
  - Lookup, Report, Workflow, Status
- Created 14 placeholder repository files in `repo/postgres/`
- Created 16 placeholder handler files in `handler/`
- All code compiles successfully following template.md structure

---

### [ ] Task 1.6: Execute Database Schema
**Reference**: `seed/db/claims_database_schema.sql`

**Steps**:
1. Create `db/01_base_schema.sql` from seed file
2. Execute schema on PostgreSQL 16 database:
   ```bash
   psql -h localhost -U postgres -d claims_db -f db/01_base_schema.sql
   ```
3. Verify all 14 tables, enums, and indexes created
4. Create `db/README.md` with migration instructions

**Verification**:
```bash
psql -h localhost -U postgres -d claims_db -c "\dt"
psql -h localhost -U postgres -d claims_db -c "\dT"
psql -h localhost -U postgres -d claims_db -c "SELECT indexname FROM pg_indexes WHERE schemaname = 'public'"
```

---

### [ ] Task 1.7: Create Main Application Entry Point
**Reference**: `seed/template/template.md` - Main Application Entry Point section

**Steps**:
1. Create `main.go`:
   - Import bootstrap and bootstrapper packages
   - Create app with FxHandler and FxRepo modules
   - Run with context.Background()

**Verification**:
```bash
go run main.go
# Verify server starts (should fail on routes initially, but bootstrap should work)
```

---

## Phase 2: Death Claims Core Implementation
**Duration**: Week 2
**Objective**: Implement death claim registration and processing workflows

### [ ] Task 2.1: Create ClaimRepository
**Reference**: `seed/template/template.md` - Repository Pattern section
**Reference**: `seed/tool-docs/db-README.md` - Database access patterns

**Steps**:
1. Create `repo/postgres/claim.go`
2. Implement methods:
   - Create(ctx, domain.Claim) (domain.Claim, error)
   - FindByID(ctx, claimID string) (domain.Claim, error)
   - FindByClaimNumber(ctx, claimNumber string) (domain.Claim, error)
   - List(ctx, filters, skip, limit) ([]domain.Claim, int64, error)
   - Update(ctx, claimID string, updates) (domain.Claim, error)
   - UpdateStatus(ctx, claimID string, status string) error
   - GetApprovalQueue(ctx, filters) ([]domain.Claim, int64, error)
   - GetPaymentQueue(ctx, filters) ([]domain.Claim, int64, error)
3. Use n-api-db patterns:
   - dblib.SelectOneFX for single row
   - dblib.SelectRowsFX for multiple rows
   - dblib.Insert for inserts
   - dblib.Update for updates
   - Context timeout from config

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestClaimRepository
```

---

### [ ] Task 2.2: Create ClaimDocumentRepository
**Reference**: Same as Task 2.1

**Steps**:
1. Create `repo/postgres/claim_document.go`
2. Implement CRUD + bulk operations for documents

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestClaimDocumentRepository
```

---

### [ ] Task 2.3: Create Death Claim Request DTOs
**Reference**: `seed/template/template.md` - Request DTO Pattern section

**Steps**:
Add to `handler/request.go`:
1. RegisterDeathClaimRequest
2. CalculateDeathClaimAmountRequest
3. GetDocumentChecklistRequest
4. UploadClaimDocumentsRequest
5. CheckDocumentCompletenessUri
6. CalculateBenefitRequest
7. GetEligibleApproversUri
8. GetApprovalDetailsUri
9. ApproveClaimRequest
10. RejectClaimRequest
11. DisburseClaimRequest
12. CloseClaimRequest
13. CancelClaimRequest
14. ListClaimsParams

**Verification**:
```bash
cd handler && govalid
# Verify validators generated
```

---

### [ ] Task 2.4: Create Death Claim Response DTOs
**Reference**: `seed/template/template.md` - Response DTO Pattern section

**Steps**:
Create `handler/response/claim.go`:
1. DeathClaimResponse
2. DeathClaimRegisteredResponse
3. DeathClaimsListResponse
4. DocumentChecklistResponse
5. ClaimAmountCalculationResponse
6. DocumentCompletenessResponse
7. BenefitCalculationResponse
8. EligibleApproversResponse
9. ApprovalDetailsResponse
10. ClaimApprovalResponse
11. ClaimDisbursementResponse

**Verification**:
```bash
go build ./handler/response/...
```

---

### [ ] Task 2.5: Implement ClaimHandler with 15 Core Endpoints
**Reference**: `seed/template/template.md` - Handler Pattern section
**Reference**: `seed/swagger/` - Endpoint specifications

**Steps**:
Create `handler/claim.go` with routes:
1. POST `/claims/death/register` - RegisterDeathClaim
2. POST `/claims/death/calculate-amount` - CalculateDeathClaimAmount
3. GET `/claims/death/:claim_id/document-checklist` - GetDocumentChecklist
4. GET `/claims/death/document-checklist-dynamic` - GetDynamicDocumentChecklist
5. POST `/claims/death/:claim_id/documents` - UploadClaimDocuments
6. GET `/claims/death/:claim_id/document-completeness` - CheckDocumentCompleteness
7. POST `/claims/death/:claim_id/calculate-benefit` - CalculateBenefit
8. GET `/claims/death/:claim_id/eligible-approvers` - GetEligibleApprovers
9. GET `/claims/death/:claim_id/approval-details` - GetApprovalDetails
10. POST `/claims/death/:claim_id/approve` - ApproveClaim
11. POST `/claims/death/:claim_id/reject` - RejectClaim
12. POST `/claims/death/:claim_id/disburse` - DisburseClaim
13. POST `/claims/death/:claim_id/close` - CloseClaim
14. POST `/claims/death/:claim_id/cancel` - CancelClaim
15. GET `/claims/death/approval-queue` - GetApprovalQueue

Each handler must:
- Follow template pattern
- Inject ClaimRepository
- Use proper logging (log.Error, log.Info)
- Handle pgx.ErrNoRows for 404

**Verification**:
```bash
go test ./handler/... -v -tags=integration
# Test each endpoint manually with curl
```

---

### [ ] Task 2.6: Implement Business Rules for Death Claims
**Reference**: `.zenflow/tasks/code-gen-54c7/requirements.md` - Business Rules section

**Steps**:
Implement business rules in ClaimRepository or ClaimHandler:
1. BR-CLM-DC-001: Investigation trigger (3-year rule)
2. BR-CLM-DC-002: Investigation SLA (21 days)
3. BR-CLM-DC-003: SLA without investigation (15 days)
4. BR-CLM-DC-004: SLA with investigation (45 days)
5. BR-CLM-DC-009: Penal interest calculation
6. BR-CLM-DC-021: SLA color coding (GREEN/YELLOW/ORANGE/RED)

**Verification**:
```bash
go test ./... -v -run TestBusinessRules
```

---

### [ ] Task 2.7: Register ClaimHandler in Bootstrap
**Reference**: `seed/template/template.md` - Bootstrap Configuration section

**Steps**:
1. Update `bootstrap/bootstrapper.go`
2. Add ClaimRepository to FxRepo
3. Add ClaimHandler to FxHandler with fx.Annotate

**Verification**:
```bash
go run main.go
# Verify routes are registered
curl http://localhost:8080/v1/claims/death/approval-queue
```

---

## Phase 3: Investigation Workflow
**Duration**: Week 3
**Objective**: Implement investigation assignment and tracking

### [ ] Task 3.1: Create InvestigationRepository
**Reference**: `seed/template/template.md` - Repository Pattern section

**Steps**:
1. Create `repo/postgres/investigation.go`
2. Implement CRUD methods
3. Implement investigation-specific queries (active investigations, SLA tracking)

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestInvestigationRepository
```

---

### [ ] Task 3.2: Create InvestigationProgressRepository
**Reference**: Same as Task 3.1

**Steps**:
1. Create `repo/postgres/investigation_progress.go`
2. Implement heartbeat tracking methods
3. Implement progress update methods

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestInvestigationProgressRepository
```

---

### [ ] Task 3.3: Implement InvestigationHandler (10 endpoints)
**Reference**: `seed/swagger/` - Investigation endpoints

**Steps**:
Create `handler/investigation.go` with routes:
1. POST `/investigations/:investigation_id/assign` - AssignInvestigator
2. POST `/investigations/:investigation_id/progress` - UpdateProgress
3. GET `/investigations/:investigation_id/report` - GetInvestigationReport
4. POST `/investigations/:investigation_id/complete` - CompleteInvestigation
5. POST `/investigations/:investigation_id/reopen` - ReopenInvestigation
6. GET `/claims/:claim_id/investigations` - GetClaimInvestigations
7. POST `/investigations/:investigation_id/evidence` - UploadEvidence
8. GET `/investigations/pending` - GetPendingInvestigations
9. GET `/investigations/overdue` - GetOverdueInvestigations
10. GET `/investigations/:investigation_id/timeline` - GetInvestigationTimeline

**Verification**:
```bash
go test ./handler/... -v -tags=integration
```

---

### [ ] Task 3.4: Register InvestigationHandler in Bootstrap
**Reference**: `seed/template/template.md` - Bootstrap Configuration section

**Steps**:
1. Add InvestigationRepository to FxRepo
2. Add InvestigationProgressRepository to FxRepo
3. Add InvestigationHandler to FxHandler

**Verification**:
```bash
go run main.go
curl http://localhost:8080/v1/investigations/pending
```

---

## Phase 4: Maturity & Survival Benefit Claims
**Duration**: Week 4
**Objective**: Implement maturity and SB claim processing

### [ ] Task 4.1: Create MaturityClaimHandler (12 endpoints)
**Reference**: `seed/swagger/` - Maturity claim endpoints

**Steps**:
1. Create request/response DTOs in handler/request.go and handler/response/maturity.go
2. Create `handler/maturity.go` with 12 maturity claim endpoints
3. Implement batch intimation logic for maturity claims
4. Implement OCR data extraction integration

**Verification**:
```bash
go test ./handler/... -v -run TestMaturityHandler
```

---

### [ ] Task 4.2: Create SurvivalBenefitHandler (2 endpoints)
**Reference**: `seed/swagger/` - Survival benefit endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/survival_benefit.go` with 2 SB endpoints

**Verification**:
```bash
go test ./handler/... -v -run TestSurvivalBenefitHandler
```

---

### [ ] Task 4.3: Implement batch intimation job
**Reference**: Requirements document - batch processing requirements

**Steps**:
1. Create cron job or scheduled task for maturity claim intimations
2. Query policies maturing in next 30 days
3. Send notifications to policyholders

**Verification**:
```bash
# Test batch job manually
# Verify notifications sent
```

---

### [ ] Task 4.4: Register handlers in bootstrap
**Reference**: `seed/template/template.md` - Bootstrap Configuration section

**Steps**:
1. Add MaturityClaimHandler to FxHandler
2. Add SurvivalBenefitHandler to FxHandler

**Verification**:
```bash
go run main.go
```

---

## Phase 5: AML/CFT & Banking Services
**Duration**: Week 5
**Objective**: Implement AML detection and payment processing

### [ ] Task 5.1: Create AMLAlertRepository
**Reference**: `seed/template/template.md` - Repository Pattern section

**Steps**:
1. Create `repo/postgres/aml_alert.go`
2. Implement CRUD methods
3. Implement risk scoring queries

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestAMLAlertRepository
```

---

### [ ] Task 5.2: Implement AMLHandler (7 endpoints)
**Reference**: `seed/swagger/` - AML endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/aml.go` with 7 AML endpoints
3. Implement AML trigger detection logic (70+ rules from requirements)

**Verification**:
```bash
go test ./handler/... -v -run TestAMLHandler
```

---

### [ ] Task 5.3: Implement 70+ AML trigger rules
**Reference**: `.zenflow/tasks/code-gen-54c7/requirements.md` - AML business rules

**Steps**:
1. Implement each AML rule as a separate function
2. Create rule engine to evaluate all applicable rules
3. Log rule violations with BR-CLM-AML-* references

**Verification**:
```bash
go test ./... -v -run TestAMLRules
```

---

### [ ] Task 5.4: Create BankingHandler (8 endpoints)
**Reference**: `seed/swagger/` - Banking endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/banking.go` with 8 banking endpoints
3. Implement bank validation logic

**Verification**:
```bash
go test ./handler/... -v -run TestBankingHandler
```

---

### [ ] Task 5.5: Integrate CBS API for bank validation
**Reference**: Requirements document - CBS API integration

**Steps**:
1. Create CBS client in `repo/postgres/cbs_client.go`
2. Implement bank account validation
3. Implement penny drop test

**Verification**:
```bash
# Test CBS integration with test bank account
```

---

### [ ] Task 5.6: Integrate PFMS API for NEFT
**Reference**: Requirements document - PFMS API integration

**Steps**:
1. Create PFMS client in `repo/postgres/pfms_client.go`
2. Implement NEFT disbursement
3. Implement payment status tracking

**Verification**:
```bash
# Test PFMS integration with test disbursement
```

---

## Phase 6: Free Look & Appeals
**Duration**: Week 6
**Objective**: Implement free look cancellation and appeal workflow

### [ ] Task 6.1: Create PolicyBondTrackingRepository
**Reference**: `seed/template/template.md` - Repository Pattern section

**Steps**:
1. Create `repo/postgres/policy_bond_tracking.go`
2. Implement bond tracking CRUD

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestPolicyBondTrackingRepository
```

---

### [ ] Task 6.2: Create FreeLookCancellationRepository
**Reference**: Same as Task 6.1

**Steps**:
1. Create `repo/postgres/freelook_cancellation.go`
2. Implement CRUD + refund calculation

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestFreeLookCancellationRepository
```

---

### [ ] Task 6.3: Implement FreeLookHandler (8 endpoints)
**Reference**: `seed/swagger/` - Free look endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/freelook.go` with 8 endpoints
3. Implement free look period calculation (BR-CLM-BOND-001)
4. Implement refund calculation (BR-CLM-BOND-003)
5. Implement maker-checker workflow (BR-CLM-BOND-004)

**Verification**:
```bash
go test ./handler/... -v -run TestFreeLookHandler
```

---

### [ ] Task 6.4: Create AppealRepository
**Reference**: `seed/template/template.md` - Repository Pattern section

**Steps**:
1. Create `repo/postgres/appeal.go`
2. Implement appeal CRUD

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestAppealRepository
```

---

### [ ] Task 6.5: Implement AppealHandler (3 endpoints)
**Reference**: `seed/swagger/` - Appeal endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/appeal.go` with 3 endpoints
3. Implement appeal eligibility check (BR-CLM-DC-005)
4. Implement appellate authority escalation

**Verification**:
```bash
go test ./handler/... -v -run TestAppealHandler
```

---

## Phase 7: Ombudsman & Notifications
**Duration**: Week 7
**Objective**: Implement ombudsman complaints and notification service

### [ ] Task 7.1: Create OmbudsmanComplaintRepository
**Reference**: `seed/template/template.md` - Repository Pattern section

**Steps**:
1. Create `repo/postgres/ombudsman_complaint.go`
2. Implement complaint CRUD + workflow

**Verification**:
```bash
go test ./repo/postgres/... -v -run TestOmbudsmanComplaintRepository
```

---

### [ ] Task 7.2: Implement OmbudsmanHandler
**Reference**: `seed/swagger/` - Ombudsman endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/ombudsman.go`
3. Implement admissibility checks (BR-CLM-OMB-001)
4. Implement award compliance tracking (BR-CLM-OMB-006)

**Verification**:
```bash
go test ./handler/... -v -run TestOmbudsmanHandler
```

---

### [ ] Task 7.3: Create NotificationHandler (5 endpoints)
**Reference**: `seed/swagger/` - Notification endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/notification.go` with 5 endpoints

**Verification**:
```bash
go test ./handler/... -v -run TestNotificationHandler
```

---

### [ ] Task 7.4: Implement multi-channel communication (SMS, Email, WhatsApp)
**Reference**: Requirements document - Notification Service integration

**Steps**:
1. Create notification client in `repo/postgres/notification_client.go`
2. Implement SMS sending
3. Implement Email sending
4. Implement WhatsApp sending

**Verification**:
```bash
# Test notification sending with test recipient
```

---

## Phase 8: Supporting Services
**Duration**: Week 8
**Objective**: Implement lookup, validation, reporting, and workflow services

### [ ] Task 8.1: Implement PolicyServiceHandler (8 endpoints)
**Reference**: `seed/swagger/` - Policy service endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/policy.go` with 8 endpoints
3. Integrate with Policy Service (gRPC/REST)

**Verification**:
```bash
go test ./handler/... -v -run TestPolicyServiceHandler
```

---

### [ ] Task 8.2: Implement ValidationServiceHandler (6 endpoints)
**Reference**: `seed/swagger/` - Validation service endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/validation.go` with 6 endpoints

**Verification**:
```bash
go test ./handler/... -v -run TestValidationServiceHandler
```

---

### [ ] Task 8.3: Implement LookupHandler (12 endpoints)
**Reference**: `seed/swagger/` - Lookup endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/lookup.go` with 12 endpoints for master data

**Verification**:
```bash
go test ./handler/... -v -run TestLookupHandler
```

---

### [ ] Task 8.4: Implement ReportHandler (8 endpoints)
**Reference**: `seed/swagger/` - Report endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/report.go` with 8 endpoints
3. Implement complex queries for reports

**Verification**:
```bash
go test ./handler/... -v -run TestReportHandler
```

---

### [ ] Task 8.5: Implement WorkflowHandler (6 endpoints)
**Reference**: `seed/swagger/` - Workflow endpoints

**Steps**:
1. Create request/response DTOs
2. Create `handler/workflow.go` with 6 endpoints
3. Integrate with Temporal workflow

**Verification**:
```bash
go test ./handler/... -v -run TestWorkflowHandler
```

---

### [ ] Task 8.6: Implement Status & Tracking endpoints (7 endpoints)
**Reference**: `seed/swagger/` - Status endpoints

**Steps**:
1. Create request/response DTOs
2. Add status endpoints to appropriate handlers
3. Implement timeline tracking

**Verification**:
```bash
go test ./handler/... -v -run TestStatusEndpoints
```

---

## Phase 9: Performance & Optimization
**Duration**: Week 9
**Objective**: Execute performance optimizations and testing

### [ ] Task 9.1: Execute performance optimization patch (03_performance_optimization.sql)
**Reference**: Database schema - Performance optimization file

**Steps**:
1. Create `db/03_performance_optimization.sql`
2. Execute on database:
   ```bash
   psql -h localhost -U postgres -d claims_db -f db/03_performance_optimization.sql
   ```
3. Verify partitions created

**Verification**:
```bash
psql -h localhost -U postgres -d claims_db -c "SELECT tablename FROM pg_tables WHERE tablename LIKE 'claims_%'"
```

---

### [ ] Task 9.2: Create partitions for current year + 2 years
**Reference**: Performance optimization requirements

**Steps**:
1. Create partition creation script
2. Execute for current year + 2 future years

**Verification**:
```bash
psql -h localhost -U postgres -d claims_db -c "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE 'claims_%' ORDER BY tablename"
```

---

### [ ] Task 9.3: Refresh materialized views
**Reference**: Database schema - Materialized views

**Steps**:
1. Create refresh script for materialized views
2. Schedule periodic refresh

**Verification**:
```bash
psql -h localhost -U postgres -d claims_db -c "REFRESH MATERIALIZED VIEW CONCURRENTLY mv_dashboard_counts"
```

---

### [ ] Task 9.4: Query performance testing with 100K rows
**Reference**: Performance requirements

**Steps**:
1. Generate test data (100K claims)
2. Run performance tests on all queries
3. Identify slow queries (>2s)
4. Optimize with indexes

**Verification**:
```bash
# Run performance test suite
# Verify P95 response time < 2s
```

---

### [ ] Task 9.5: Connection pool tuning
**Reference**: `seed/tool-docs/db-README.md` - Pool management

**Steps**:
1. Analyze connection pool metrics
2. Adjust min/max connections in config
3. Test under load

**Verification**:
```bash
# Monitor pool utilization during load test
# Verify pool utilization < 80%
```

---

### [ ] Task 9.6: Load testing and benchmarking
**Reference**: Performance requirements

**Steps**:
1. Create load test scripts (k6 or artillery)
2. Test with 1000 concurrent users
3. Measure P50, P95, P99 response times
4. Identify bottlenecks

**Verification**:
```bash
k6 run tests/load/death-claims-registration.js
# Verify P95 < 2s
```

---

## Phase 10: Testing & Documentation
**Duration**: Week 10
**Objective**: Complete test suite and documentation

### [ ] Task 10.1: Complete unit test suite (80%+ coverage)
**Reference**: Testing requirements

**Steps**:
1. Write unit tests for all repositories
2. Write unit tests for all handlers
3. Mock external dependencies
4. Achieve 80%+ code coverage

**Verification**:
```bash
go test ./... -cover
# Verify coverage > 80%
```

---

### [ ] Task 10.2: Integration test suite for all workflows
**Reference**: Integration testing requirements

**Steps**:
1. Write integration tests for all critical workflows
2. Test with real database
3. Test external service integrations with mocks

**Verification**:
```bash
go test ./... -v -tags=integration
```

---

### [ ] Task 10.3: End-to-end testing scripts
**Reference**: E2E testing requirements

**Steps**:
1. Create E2E test scripts for all major workflows
2. Test complete claim lifecycle
3. Test error scenarios

**Verification**:
```bash
./tests/e2e/run_all.sh
```

---

### [ ] Task 10.4: Generate Swagger documentation
**Reference**: Swagger generation requirements

**Steps**:
1. Install swag tool
2. Generate swagger.yaml from code annotations
3. Verify all endpoints documented

**Verification**:
```bash
swag init
# Verify docs/swagger.yaml contains all 130+ endpoints
```

---

### [ ] Task 10.5: Create runbook documentation
**Reference**: Documentation requirements

**Steps**:
1. Create runbook with common operational procedures
2. Document troubleshooting steps
3. Add deployment procedures

**Verification**:
```bash
# Review runbook for completeness
```

---

### [ ] Task 10.6: Create deviation.md (if any schema changes)
**Reference**: Schema deviation tracking

**Steps**:
1. Review database schema vs. seed file
2. Document any changes made
3. Include rationale and impact

**Verification**:
```bash
# Verify deviation.md is accurate
```

---

## Phase 11: Security & Compliance
**Duration**: Week 11
**Objective**: Security hardening and compliance verification

### [ ] Task 11.1: Verify RLS policies are active
**Reference**: Security requirements - RLS

**Steps**:
1. Test RLS policies with different user roles
2. Verify data isolation

**Verification**:
```bash
psql -h localhost -U postgres -d claims_db -c "SELECT schemaname, tablename, policyname FROM pg_policies"
```

---

### [ ] Task 11.2: Verify audit trail functionality
**Reference**: Security requirements - Audit trail

**Steps**:
1. Test audit trail creation
2. Verify all data changes logged
3. Verify override tracking

**Verification**:
```bash
# Create test claim, modify it, verify audit trail
```

---

### [ ] Task 11.3: Implement digital signatures (BR-CLM-DC-025)
**Reference**: Business rule - Digital signatures

**Steps**:
1. Implement digital signature generation
2. Implement signature verification
3. Store signature hash

**Verification**:
```bash
go test ./... -v -run TestDigitalSignature
```

---

### [ ] Task 11.4: Security audit (OWASP Top 10)
**Reference**: Security requirements

**Steps**:
1. Run OWASP ZAP scan
2. Review and fix vulnerabilities
3. Verify no critical issues

**Verification**:
```bash
zap-cli quick-scan --self-contained http://localhost:8080
```

---

### [ ] Task 11.5: Penetration testing
**Reference**: Security requirements

**Steps**:
1. Conduct penetration testing
2. Document findings
3. Fix identified issues

**Verification**:
```bash
# Review penetration test report
```

---

### [ ] Task 11.6: Compliance verification (IRDAI guidelines)
**Reference**: Compliance requirements

**Steps**:
1. Verify compliance with IRDAI guidelines
2. Document compliance checklist
3. Address any gaps

**Verification**:
```bash
# Review compliance checklist
```

---

## Phase 12: Deployment & Handover
**Duration**: Week 12
**Objective**: Production deployment and documentation

### [ ] Task 12.1: Production build optimization
**Reference**: Deployment requirements

**Steps**:
1. Optimize Go build for production
2. Reduce binary size
3. Enable optimizations

**Verification**:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o bin/claims-api main.go
```

---

### [ ] Task 12.2: Docker containerization
**Reference**: Deployment requirements

**Steps**:
1. Create Dockerfile
2. Create .dockerignore
3. Test container build

**Verification**:
```bash
docker build -t claims-api:latest .
docker run claims-api:latest
```

---

### [ ] Task 12.3: Kubernetes manifests
**Reference**: Deployment requirements

**Steps**:
1. Create deployment.yaml
2. Create service.yaml
3. Create configmap.yaml
4. Create secret.yaml

**Verification**:
```bash
kubectl apply -f k8s/
kubectl get pods
```

---

### [ ] Task 12.4: CI/CD pipeline setup
**Reference**: Deployment requirements

**Steps**:
1. Create CI/CD pipeline (GitLab CI or GitHub Actions)
2. Configure build, test, deploy stages
3. Configure automated deployments

**Verification**:
```bash
# Trigger pipeline, verify all stages pass
```

---

### [ ] Task 12.5: Production deployment
**Reference**: Deployment requirements

**Steps**:
1. Deploy to production environment
2. Verify health endpoints
3. Monitor logs

**Verification**:
```bash
curl http://prod-server/health
```

---

### [ ] Task 12.6: Smoke testing in production
**Reference**: Deployment requirements

**Steps**:
1. Run smoke tests against production
2. Test critical endpoints
3. Verify integrations

**Verification**:
```bash
./tests/smoke/prod_smoke.sh
```

---

### [ ] Task 12.7: Monitoring dashboards
**Reference**: Monitoring requirements

**Steps**:
1. Create Grafana dashboards
2. Configure Prometheus metrics
3. Set up alerts

**Verification**:
```bash
# Verify dashboards are accessible
# Verify metrics are collected
```

---

### [ ] Task 12.8: Alert configuration
**Reference**: Monitoring requirements

**Steps**:
1. Configure alerts for:
   - SLA breaches
   - Error rates
   - Response times
   - Database connection issues

**Verification**:
```bash
# Test alert conditions
# Verify alerts are sent
```

---

### [ ] Task 12.9: Handover documentation
**Reference**: Documentation requirements

**Steps**:
1. Create comprehensive handover document
2. Document architecture decisions
3. Document operational procedures

**Verification**:
```bash
# Review handover document for completeness
```

---

### [ ] Task 12.10: Training materials
**Reference**: Training requirements

**Steps**:
1. Create training presentation
2. Create troubleshooting guide
3. Create FAQ document

**Verification**:
```bash
# Deliver training session
```

---

## Success Criteria

### P0 (Must Have)
- [ ] 100% Swagger endpoint implementation (130+ endpoints)
- [ ] 100% Database schema compliance (all tables, indexes, functions)
- [ ] 100% Template structure compliance (n-api-template)
- [ ] All 70+ business rules implemented with references
- [ ] n-api-db library usage with pooling for all queries

### P1 (Should Have)
- [ ] Integration tests for all critical workflows
- [ ] Performance optimization (partitions, indexes)
- [ ] Comprehensive error handling with proper HTTP codes
- [ ] Swagger documentation auto-generated
- [ ] RLS policies active and tested
- [ ] Unit test coverage > 80%

### P2 (Nice to Have)
- [ ] Monitoring dashboards
- [ ] Runbook documentation
- [ ] Training materials
- [ ] Docker containers
- [ ] Kubernetes manifests

---

## Implementation Notes

1. **Strict Template Adherence**: Every file must follow the template.md structure exactly
2. **Database Access**: Use n-api-db patterns from db-README.md (pooling, parallel queries, batch operations)
3. **Business Rules**: All rules must reference BR-CLM-* identifiers in comments
4. **Optimized Queries**: Analyze requirements and use batch/parallel queries where appropriate
5. **Schema Deviations**: Document any required schema changes in deviation.md with clear rationale

---

**End of Implementation Plan**
