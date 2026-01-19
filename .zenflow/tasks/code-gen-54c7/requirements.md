# Product Requirements Document (PRD)
## Claims Processing Microservice - Code Generation

**Project**: PLI Claims Processing API
**Version**: 1.0.0
**Date**: 2026-01-19
**Status**: Requirements Analysis

---

## 1. Executive Summary

### 1.1 Purpose
Generate production-ready Go code for a Claims Processing Microservice for India Post Life Insurance (PLI) based on comprehensive business requirements, API specifications, and database schema.

### 1.2 Scope
The Claims Processing Microservice will handle the complete lifecycle of insurance claims including:
- Death Claims (with investigation workflow)
- Maturity Claims
- Survival Benefit Claims
- Free Look Cancellations
- AML/CFT Detection and Compliance
- Appeals and Ombudsman Integration
- Policy Bond Tracking

### 1.3 Key Constraints
1. **Template Compliance**: Must strictly follow `seed/template/template.md` structure
2. **API Contract**: Must implement all 130+ endpoints defined in Swagger specification
3. **Database Schema**: Must use the provided PostgreSQL DDL without modifications
4. **Business Rules**: Must implement all 70+ business rules from SRS documents
5. **Database Access**: Must use n-api-db library with pooling and batch operations

---

## 2. Business Requirements

### 2.1 Claim Types and Workflows

#### 2.1.1 Death Claims
**Workflow**:
1. Registration → Document Verification → Investigation (conditional) → Calculation → Approval → Disbursement → Closure

**Key Business Rules**:
- BR-CLM-DC-001: Auto-trigger investigation if death within 3 years of policy issue/revival
- BR-CLM-DC-002: Investigation SLA is 21 days
- BR-CLM-DC-003: SLA 15 days (no investigation) or 45 days (with investigation)
- BR-CLM-DC-009: Penal interest @ 8% p.a. for SLA breaches
- BR-CLM-DC-013: Maximum 2 reinvestigations allowed
- BR-CLM-DC-021: Color-coded SLA (GREEN <70%, YELLOW 70-90%, ORANGE 90-100%, RED >100%)

**Key Features**:
- Dynamic document checklist based on death type and nomination status
- Investigation officer assignment with jurisdiction-based routing
- Calculation with sum assured, bonuses, loans, and unpaid premiums
- Approval hierarchy based on claim amount
- Multi-mode disbursement (NEFT > POSB > Cheque)

#### 2.1.2 Maturity Claims
**Workflow**:
1. Intimation (60 days before) → Pre-fill Data → Submission → QC Verification → Approval → Disbursement

**Key Features**:
- Automated maturity intimation batch
- OCR-based data extraction from documents
- Online/offline submission modes
- 7-day SLA for processing

#### 2.1.3 Survival Benefit Claims
**Workflow**:
1. Eligibility Validation → Submission → Verification → Disbursement

**Key Features**:
- DigiLocker integration for document fetching
- Optional digital document submission

#### 2.1.4 Free Look Cancellations
**Workflow**:
1. Bond Tracking → Cancellation Request → Refund Calculation → Maker-Checker Approval → Disbursement

**Key Business Rules**:
- BR-CLM-BOND-001: 15 days (physical bond) or 30 days (electronic) free look period
- BR-CLM-BOND-003: Refund = Premium - (risk premium + stamp duty + medical + other)
- BR-CLM-BOND-004: Maker-checker workflow required

### 2.2 AML/CFT Requirements

**Triggers**:
- BR-CLM-AML-001: Cash transactions over ₹50,000
- BR-CLM-AML-003: Nominee change after policyholder death
- BR-CLM-AML-006: STR filing within 7 days
- BR-CLM-AML-007: CTR filing monthly

**Risk Levels**: LOW, MEDIUM, HIGH, CRITICAL

**Filing Types**: STR (Suspicious Transaction Report), CTR (Currency Transaction Report)

### 2.3 Appeals and Ombudsman

**Appeal Workflow**:
- BR-CLM-DC-005: 90-day appeal window after rejection
- BR-CLM-DC-007: 45-day SLA for appellate decision
- Escalation to next higher authority in approval hierarchy

**Ombudsman Workflow**:
- BR-CLM-OMB-001: Admissibility checks (₹50 lakh cap, 1-year limitation)
- BR-CLM-OMB-005: Award amount cap ₹50 lakh
- BR-CLM-OMB-006: 30-day compliance timeline

---

## 3. Functional Requirements

### 3.1 API Endpoints (130+)

#### 3.1.1 Death Claims Core (15 endpoints)
- POST `/claims/death/register` - Register death claim
- POST `/claims/death/calculate-amount` - Pre-calculate benefit
- GET `/claims/death/{claim_id}/document-checklist` - Get checklist
- GET `/claims/death/document-checklist-dynamic` - Dynamic checklist
- POST `/claims/death/{claim_id}/documents` - Upload documents
- GET `/claims/death/{claim_id}/document-completeness` - Check completeness
- POST `/claims/death/{claim_id}/send-reminder` - Send reminders
- POST `/claims/death/{claim_id}/calculate-benefit` - Calculate benefit
- GET `/claims/death/{claim_id}/eligible-approvers` - Get approvers
- GET `/claims/death/{claim_id}/approval-details` - Get details for approval
- POST `/claims/death/{claim_id}/approve` - Approve claim
- POST `/claims/death/{claim_id}/reject` - Reject claim
- POST `/claims/death/{claim_id}/disburse` - Initiate payment
- POST `/claims/death/{claim_id}/close` - Close claim
- POST `/claims/death/{claim_id}/cancel` - Cancel claim

#### 3.1.2 Investigation Workflow (10 endpoints)
- POST `/claims/death/{claim_id}/investigation/assign-officer` - Assign investigator
- GET `/claims/death/pending-investigation` - Get queue
- GET `/claims/death/{claim_id}/investigation/{investigation_id}/details` - Get details
- POST `/claims/death/{claim_id}/investigation/{investigation_id}/progress-update` - Heartbeat update
- POST `/claims/death/{claim_id}/investigation/{investigation_id}/submit-report` - Submit report
- POST `/claims/death/{claim_id}/investigation/{investigation_id}/review` - Review report
- POST `/claims/death/{id}/investigation/trigger-reinvestigation` - Reinvestigate
- POST `/claims/death/{id}/investigation/escalate-sla-breach` - Escalate
- POST `/claims/death/{id}/manual-review/assign` - Assign for manual review
- POST `/claims/death/{id}/reject-fraud` - Reject for fraud

#### 3.1.3 Queue Management (2 endpoints)
- GET `/claims/death/approval-queue` - Get approval queue
- GET `/claims/death/payment-queue` - Get payment queue

#### 3.1.4 Maturity Claims (12 endpoints)
- POST `/claims/maturity/send-intimation-batch` - Batch intimation
- POST `/claims/maturity/generate-due-report` - Generate report
- GET `/claims/maturity/pre-fill-data` - Get pre-filled data
- POST `/claims/maturity/submit` - Submit claim
- POST `/claims/maturity/{claim_id}/validate-documents` - Validate documents
- POST `/claims/maturity/{claim_id}/extract-ocr-data` - OCR extraction
- POST `/claims/maturity/{claim_id}/qc-verify` - QC verification
- POST `/claims/maturity/{claim_id}/validate-bank` - Validate bank
- GET `/claims/maturity/{claim_id}/approval-details` - Get details
- POST `/claims/maturity/{claim_id}/approve` - Approve
- POST `/claims/maturity/{claim_id}/disburse` - Disburse
- POST `/claims/maturity/{claim_id}/close` - Close

#### 3.1.5 Survival Benefit (2 endpoints)
- POST `/claims/survival-benefit/submit` - Submit SB claim
- POST `/claims/survival-benefit/{id}/validate-eligibility` - Validate eligibility

#### 3.1.6 Policy Services (8 endpoints)
- GET `/policies/{id}/details` - Get policy details
- GET `/policies/{policy_id}/claim-eligibility` - Check eligibility
- GET `/policies/{id}/benefit-calculation` - Get calculation inputs
- GET `/policies/{policy_id}/maturity-amount` - Get maturity amount
- GET `/bonuses/{policy_id}/accrued` - Get accrued bonuses
- GET `/loans/{policy_id}/outstanding` - Get outstanding loan
- GET `/premiums/{policy_id}/unpaid` - Get unpaid premiums
- POST `/policies/{policy_id}/freelook-refund-calculation` - Calculate refund

#### 3.1.7 Banking Services (8 endpoints)
- POST `/banking/validate-account` - Validate account
- POST `/banking/validate-account-cbs` - Validate via CBS
- POST `/banking/validate-account-pfms` - Validate via PFMS
- POST `/banking/penny-drop` - Penny drop test
- POST `/banking/neft-transfer` - Initiate NEFT
- POST `/banking/payment-reconciliation` - Reconcile payments
- GET `/payments/{payment_id}/status` - Get payment status
- POST `/webhooks/banking/payment-confirmation` - Payment webhook

#### 3.1.8 AML/CFT (7 endpoints)
- POST `/aml/detect-trigger` - Detect triggers
- POST `/aml/{alert_id}/generate-alert` - Generate alert
- POST `/aml/{alert_id}/calculate-risk-score` - Calculate risk
- GET `/aml/{alert_id}/details` - Get alert details
- POST `/aml/{alert_id}/review` - Review alert
- POST `/aml/{alert_id}/file-report` - File STR/CTR

#### 3.1.9 Document Management (5 endpoints)
- POST `/virus-scan/scan` - Scan for viruses
- POST `/ecms/upload` - Upload to ECMS
- POST `/ecms/archive-claim` - Archive claim
- POST `/documents/generate-sanction-letter` - Generate sanction letter
- POST `/documents/generate-rejection-letter` - Generate rejection letter

#### 3.1.10 Notifications (5 endpoints)
- POST `/notifications/send` - Send notification
- POST `/notifications/send-batch` - Batch notifications
- POST `/feedback/generate-link` - Generate feedback link

#### 3.1.11 Status and Tracking (7 endpoints)
- GET `/claims/{claim_id}/status` - Get claim status
- GET `/claims/{claim_id}/sla-countdown` - Get SLA countdown
- GET `/claims/{claim_id}/payment-status` - Get payment status
- GET `/claims/{claim_id}/timeline` - Get timeline
- GET `/claims/{claim_id}/communications` - Get communications
- POST `/claims/{claim_id}/calculate-penal-interest` - Calculate penal interest
- GET `/customers/{customer_id}/claims` - Get customer claims

#### 3.1.12 Lookup Data (7 endpoints)
- GET `/lookup/claimant-relationships` - Get relationships
- GET `/lookup/death-types` - Get death types
- GET `/lookup/document-types` - Get document types
- GET `/lookup/rejection-reasons` - Get rejection reasons
- GET `/lookup/investigation-officers` - Get investigators
- GET `/lookup/approvers` - Get approvers
- GET `/lookup/payment-modes` - Get payment modes

#### 3.1.13 Validation Services (6 endpoints)
- POST `/validate/pan` - Validate PAN
- POST `/validate/bank-account` - Validate bank account
- POST `/validate/death-date` - Validate death date
- GET `/validate/ifsc/{ifsc_code}` - Validate IFSC
- GET `/forms/death-claim/fields` - Get form fields

### 3.2 Database Schema Compliance

**Tables** (14 main tables):
1. `claims` - Master claims table (partitioned)
2. `claim_documents` - Document storage (partitioned)
3. `investigations` - Investigation workflow
4. `appeals` - Appeal workflow
5. `aml_alerts` - AML/CFT detection
6. `claim_payments` - Payment tracking (partitioned)
7. `claim_history` - Audit trail
8. `claim_communications` - Communication log
9. `document_checklist_templates` - Dynamic checklists
10. `claim_sla_tracking` - SLA monitoring
11. `ombudsman_complaints` - Ombudsman workflow
12. `policy_bond_tracking` - Bond delivery
13. `freelook_cancellations` - Free look refunds
14. `investigation_progress` - Investigation heartbeat

**Key Features**:
- Partitioning by date (claims, documents, payments)
- 115+ indexes for performance
- Row-level security (RLS)
- Audit trail with override tracking
- Materialized views for analytics
- Business rule functions (penal interest, SLA status)

### 3.3 Business Rule Implementations

**Calculation Rules**:
- Death Claim: SA + Reversionary Bonus + Terminal Bonus - Outstanding Loan - Unpaid Premiums
- Penal Interest: (Claim Amount × 8% × Days Delayed) / 365
- Free Look Refund: Premium - (Pro-rata Risk Premium + Stamp Duty + Medical + Other)

**SLA Rules**:
- Death Claim (no investigation): 15 days
- Death Claim (with investigation): 45 days
- Investigation: 21 days
- Maturity Claim: 7 days
- Appeal Decision: 45 days
- Ombudsman Compliance: 30 days

**Approval Hierarchy**:
- Up to ₹1 lakh: Postmaster
- ₹1-5 lakh: Superintendent
- ₹5-10 lakh: Senior Superintendent
- Above ₹10 lakh: Chief Postmaster General

---

## 4. Non-Functional Requirements

### 4.1 Performance Requirements
- Response time < 2s for simple queries
- Response time < 5s for complex aggregations
- Support 1000+ concurrent users
- Handle 100K claims in database
- Partitioning for tables > 100K rows

### 4.2 Scalability Requirements
- Horizontal scaling via microservices
- Database partitioning by date
- Connection pooling (min 1, max 10)
- Slice pooling for query results

### 4.3 Security Requirements
- JWT-based authentication
- Row-level security (RLS) policies
- Role-based access control (ADMIN, APPROVER, USER)
- Audit trail for all changes
- Digital signatures for overrides

### 4.4 Reliability Requirements
- 99.9% uptime
- Graceful degradation
- Database transaction support
- Retry mechanisms for external APIs

### 4.5 Maintainability Requirements
- Clean architecture (Handler → Repository → DB)
- Domain-driven design
- Comprehensive logging
- Error handling with proper HTTP codes

---

## 5. Technical Requirements

### 5.1 Technology Stack
- **Language**: Go 1.25.0
- **Framework**: n-api-template (Uber FX for DI)
- **Database**: PostgreSQL 16
- **Database Library**: n-api-db (with pgx driver)
- **Query Builder**: Squirrel
- **API Documentation**: Swagger/OpenAPI 3.0

### 5.2 Database Access Patterns
**Use Case 1**: Simple SELECT with pooling
```go
db.SelectRowsFX(ctx, database, poolMgr, query, pgx.RowToStructByPos[Claim])
```

**Use Case 2**: Parallel queries (Rill pattern)
```go
queries := []sq.SelectBuilder{...}
db.SelectRowsParallelFX(ctx, database, poolMgr, queries, scanFn, concurrency)
```

**Use Case 3**: Batch operations
```go
batch := db.NewBatch()
db.QueueReturnFX(batch, query1, &result1)
db.QueueReturnFX(batch, query2, &result2)
database.SendBatch(ctx, batch.Batch).Close()
```

**Use Case 4**: Raw SQL when needed
```go
db.SelectRowsRaw(ctx, database, "SELECT * FROM claims WHERE ...", args, scanFn)
```

### 5.3 Code Structure Compliance

**Must Follow Template** (`seed/template/template.md`):
```
project/
├── main.go                        # Application entry point
├── go.mod                         # Go module
├── configs/                       # Environment configs
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── bootstrap/
│   └── bootstrapper.go            # FX dependency injection
├── core/
│   ├── domain/
│   │   └── claim.go               # Domain models
│   └── port/
│       ├── request.go             # Common request structs
│       └── response.go            # Common response structs
├── handler/
│   ├── claim.go                   # HTTP handlers
│   ├── request.go                 # Request DTOs
│   ├── request_claim_validator.go # Auto-generated validators
│   └── response/
│       └── claim.go               # Response DTOs
├── repo/
│   └── postgres/
│       └── claim.go               # Repository/data access
└── db/
    └── claim.sql                  # Database schema
```

### 5.4 Naming Conventions
- **Domain**: `Claim` (PascalCase)
- **Repository**: `ClaimRepository`
- **Handler**: `ClaimHandler`
- **Request DTOs**: `CreateClaimRequest`, `UpdateClaimRequest`
- **Response DTOs**: `ClaimResponse`, `ClaimCreateResponse`
- **Functions**: `Create`, `FindByID`, `List`, `Update`, `Delete`
- **Routes**: `/claims` (plural lowercase)

---

## 6. Integration Requirements

### 6.1 External Systems
1. **Policy Service** - Validate policy, get SA, bonuses
2. **Customer Service** - Get customer details
3. **User Service** - Validate users, get approvers
4. **CBS API** - Bank account validation
5. **PFMS API** - Payment gateway
6. **ECMS** - Document management
7. **DigiLocker** - Document fetching
8. **Notification Service** - SMS/Email/WhatsApp

### 6.2 Internal Services
1. **Temporal Workflow** - Orchestrate long-running workflows
2. **Audit Service** - Centralized audit logging
3. **Configuration Service** - Dynamic configuration

---

## 7. Data Model Requirements

### 7.1 Core Entities

**Claim** (Main Entity)
```
- id (UUID, PK)
- claim_number (VARCHAR, UNIQUE)
- claim_type (ENUM: DEATH, MATURITY, SURVIVAL_BENEFIT, FREELOOK)
- policy_id (UUID, FK to policy service)
- customer_id (UUID, FK to customer service)
- status (ENUM: 15 states)
- claim_amount (NUMERIC)
- approved_amount (NUMERIC)
- investigation_required (BOOLEAN)
- sla_due_date (TIMESTAMP)
- sla_status (ENUM: GREEN, YELLOW, ORANGE, RED)
- created_at, updated_at (TIMESTAMP)
- created_by, updated_by (UUID)
- deleted_at (TIMESTAMP - soft delete)
- version (INTEGER - optimistic locking)
```

**ClaimDocument**
```
- id (UUID, PK)
- claim_id (UUID, FK)
- document_type (VARCHAR)
- document_url (TEXT)
- verified (BOOLEAN)
- ocr_extracted_data (JSONB)
```

**Investigation**
```
- id (UUID, PK)
- investigation_id (VARCHAR, UNIQUE)
- claim_id (UUID, FK)
- investigator_id (UUID)
- status (VARCHAR)
- investigation_outcome (ENUM: CLEAR, SUSPECT, FRAUD)
- due_date (TIMESTAMP)
```

**AMLAlert**
```
- id (UUID, PK)
- alert_id (VARCHAR, UNIQUE)
- trigger_code (VARCHAR)
- risk_level (ENUM: LOW, MEDIUM, HIGH, CRITICAL)
- alert_status (ENUM: FLAGGED, UNDER_REVIEW, FILED, CLOSED)
- filing_type (ENUM: STR, CTR, CCR, NTR)
```

**ClaimPayment**
```
- id (UUID, PK)
- payment_id (VARCHAR, UNIQUE)
- claim_id (UUID, FK)
- payment_amount (NUMERIC)
- payment_mode (ENUM: AUTO_NEFT, POSB_TRANSFER, CHEQUE)
- payment_status (VARCHAR)
- utr_number (VARCHAR)
```

### 7.2 Views for Operations
1. `v_active_claims` - Active claims with SLA countdown
2. `v_investigation_queue` - Priority-sorted investigation queue
3. `v_approval_queue` - SLA-prioritized approval queue
4. `v_sla_breach_report` - SLA breach analytics
5. `v_payment_queue` - Approved claims pending payment
6. `v_aml_high_risk_alerts` - High/critical risk alerts
7. `v_ombudsman_compliance_pending` - Awards pending compliance

---

## 8. Validation Requirements

### 8.1 Request Validation
- Use `validate` tags on request structs
- Auto-generate validators using `govalid`
- Validate business rules at handler level
- Return 400 for validation errors

### 8.2 Business Logic Validation
- Claim type specific validations (death_date required for DEATH claims)
- SLA calculations and color-coding
- Approval hierarchy limits
- Investigation trigger conditions
- Document completeness checks

---

## 9. Error Handling Requirements

### 9.1 HTTP Status Codes
- **200 OK** - Successful GET, PUT, DELETE
- **201 Created** - Successful POST
- **400 Bad Request** - Validation errors
- **403 Forbidden** - Insufficient authority
- **404 Not Found** - Resource not found
- **422 Unprocessable Entity** - Business rule violation
- **500 Internal Server Error** - Server errors

### 9.2 Error Response Format
```json
{
  "error_code": "CLAIM_NOT_FOUND",
  "message": "Claim with ID xxx not found",
  "severity": "ERROR",
  "details": {}
}
```

### 9.3 Logging Requirements
- **Info**: Successful operations (create, update, delete)
- **Error**: All errors with context (claim_id, user_id, operation)
- Use structured logging with `log.Error(ctx, "message: %v", err)`
- Include request ID for tracing

---

## 10. Testing Requirements

### 10.1 Unit Tests
- Repository layer tests (with mock DB)
- Handler tests (with mock repository)
- Business logic tests
- Validation tests

### 10.2 Integration Tests
- Database operations (with test DB)
- External API mocks
- End-to-end workflow tests

### 10.3 Performance Tests
- Query performance with 100K rows
- Concurrent request handling
- Partition pruning verification

---

## 11. Deployment Requirements

### 11.1 Configuration
- Environment-specific configs (dev, test, prod)
- Database credentials via environment variables
- Feature flags (tracing, caching)

### 11.2 Database Migrations
- Execute schema in order: base DDL → enhancement patch → optimization patch
- Create partitions for future years before year-end
- Refresh materialized views daily

### 11.3 Monitoring
- SLA breach alerts
- Database connection pool metrics
- API response times
- Error rates

---

## 12. Documentation Requirements

### 12.1 Code Documentation
- Godoc comments on exported functions
- Business rule comments with references (BR-CLM-XXX)
- Inline comments for complex logic

### 12.2 API Documentation
- Auto-generate Swagger from code
- Include examples for all endpoints
- Document error responses

### 12.3 Runbook Documentation
- Deployment procedures
- Troubleshooting guides
- Database maintenance procedures

---

## 13. Success Criteria

### 13.1 Must Have (P0)
- [ ] 100% Swagger endpoint implementation
- [ ] 100% Database schema compliance
- [ ] 100% Template structure compliance
- [ ] All business rules implemented
- [ ] n-api-db library usage with pooling

### 13.2 Should Have (P1)
- [ ] Unit tests for all repositories
- [ ] Integration tests for critical workflows
- [ ] Performance optimization (indexes, partitions)
- [ ] Comprehensive error handling

### 13.3 Nice to Have (P2)
- [ ] Swagger documentation auto-generation
- [ ] Load testing scripts
- [ ] Monitoring dashboards
- [ ] Runbook documentation

---

## 14. Assumptions and Dependencies

### 14.1 Assumptions
1. PostgreSQL 16 is available and configured
2. External services (Policy, Customer, User) are accessible via gRPC/REST
3. ECMS is available for document storage
4. Banking gateways (CBS, PFMS) are accessible
5. Notification service is operational

### 14.2 Dependencies
1. **Go 1.25.0** - Minimum Go version
2. **n-api-bootstrapper** - Application bootstrapping
3. **n-api-server** - HTTP server framework
4. **n-api-db** - Database access library
5. **api-config** - Configuration management
6. **n-api-validation** - Request validation
7. **n-api-log** - Structured logging
8. **Uber FX** - Dependency injection

---

## 15. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Database schema changes | High | Use versioned migrations, document deviations |
| External service unavailability | Medium | Implement circuit breakers, retries |
| Performance issues | High | Use partitioning, indexes, connection pooling |
| Security vulnerabilities | High | Use RLS, input validation, audit logging |
| Template deviations | Medium | Code review against template.md |

---

## 16. Open Questions

1. **Claim Number Generation**: Should it be auto-generated or provided by user?
   - **Assumption**: System generates using sequence: CLM{YYYY}{DDDD}

2. **Temporal Workflow Integration**: Should we use Temporal or custom workflow?
   - **Recommendation**: Use Temporal for complex workflows (investigation, appeal)

3. **Document Storage**: ECMS or local storage?
   - **Assumption**: ECMS via integration service

4. **User Context**: How to get user_id for RLS?
   - **Assumption**: From JWT token in request context

---

## Appendix A: Business Rules Reference

### Death Claim Rules (BR-CLM-DC-001 to 025)
- BR-CLM-DC-001: Investigation trigger (3-year rule)
- BR-CLM-DC-002: Investigation SLA (21 days)
- BR-CLM-DC-003: SLA without investigation (15 days)
- BR-CLM-DC-004: SLA with investigation (45 days)
- BR-CLM-DC-005: Appeal window (90 days)
- BR-CLM-DC-006: Appeal decision SLA (45 days)
- BR-CLM-DC-007: Condonation of delay
- BR-CLM-DC-008: Calculation formula
- BR-CLM-DC-009: Penal interest (8% p.a.)
- BR-CLM-DC-010: Auto-return (22 days)
- BR-CLM-DC-011: Document checklist
- BR-CLM-DC-012: Document completeness
- BR-CLM-DC-013: Conditional documents (unnatural death)
- BR-CLM-DC-014: Nomination absence documents
- BR-CLM-DC-015: Base mandatory documents
- BR-CLM-DC-016: Manual override audit
- BR-CLM-DC-017: Payment mode priority
- BR-CLM-DC-018: Bank validation
- BR-CLM-DC-019: Communication triggers
- BR-CLM-DC-020: SLA color coding
- BR-CLM-DC-021: SLA status calculation
- BR-CLM-DC-022: Approval hierarchy
- BR-CLM-DC-023: Reinvestigation limit
- BR-CLM-DC-024: Fraud detection
- BR-CLM-DC-025: Digital signatures

### AML Rules (BR-CLM-AML-001 to 007)
- BR-CLM-AML-001: Cash threshold (₹50k)
- BR-CLM-AML-002: Multiple transactions
- BR-CLM-AML-003: Nominee change detection
- BR-CLM-AML-004: High-risk jurisdictions
- BR-CLM-AML-005: Politically exposed persons
- BR-CLM-AML-006: STR timeline (7 days)
- BR-CLM-AML-007: CTR timeline (monthly)

### Ombudsman Rules (BR-CLM-OMB-001 to 006)
- BR-CLM-OMB-001: Admissibility (₹50L cap, 1-year limit)
- BR-CLM-OMB-002: Jurisdiction
- BR-CLM-OMB-003: Conflict of interest
- BR-CLM-OMB-004: Mediation attempt
- BR-CLM-OMB-005: Award cap (₹50L)
- BR-CLM-OMB-006: Compliance (30 days)

### Policy Bond Rules (BR-CLM-BOND-001 to 004)
- BR-CLM-BOND-001: Free look period (15/30 days)
- BR-CLM-BOND-002: Delivery failure escalation
- BR-CLM-BOND-003: Refund calculation
- BR-CLM-BOND-004: Maker-checker workflow

---

## Appendix B: Database Deviations Log

**File**: `deviation.md`

This file will document any changes made to the database schema during code generation and the rationale for those changes.

**Example**:
```
### Deviation 1: Added column to claims table
- **Date**: 2026-01-19
- **Change**: Added `workflow_state VARCHAR(50)` to claims table
- **Reason**: Required for Temporal workflow integration
- **Impact**: Need to update domain model and queries
```

---

**End of PRD**
