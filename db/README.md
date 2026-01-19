# PLI Claims Processing Database

This directory contains the database schema and migration scripts for the PLI Claims Processing API.

## Database Information

- **Database Name**: `claims_db`
- **PostgreSQL Version**: 16
- **Schema Version**: 1.0.0

## Prerequisites

1. PostgreSQL 16 installed and running
2. Database user with CREATE privileges
3. psql client installed

## Database Setup

### Step 1: Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE claims_db;

# Exit psql
\q
```

### Step 2: Execute Base Schema

```bash
# Execute the schema file
psql -h localhost -U postgres -d claims_db -f db/01_base_schema.sql

# Or with specific host and port
psql -h <hostname> -p <port> -U <username> -d claims_db -f db/01_base_schema.sql
```

### Step 3: Verify Schema Creation

```bash
# Connect to the database
psql -h localhost -U postgres -d claims_db

# List all tables
\dt

# Expected output should show 14 tables:
# - claims
# - claim_documents
# - investigations
# - investigation_progress
# - appeals
# - aml_alerts
# - claim_payments
# - claim_history
# - claim_communications
# - document_checklist_templates
# - claim_sla_tracking
# - ombudsman_complaints
# - policy_bond_tracking
# - freelook_cancellations

# List all enums
\dT

# Expected output should show 12 enum types:
# - claim_type_enum
# - claim_status_enum
# - death_type_enum
# - claimant_type_enum
# - payment_mode_enum
# - investigation_status_enum
# - investigation_outcome_enum
# - aml_risk_level_enum
# - aml_alert_status_enum
# - aml_filing_type_enum
# - sla_status_enum

# List all indexes
SELECT indexname FROM pg_indexes WHERE schemaname = 'public' ORDER BY indexname;

# Expected: 60+ indexes covering all tables

# List all views
\dv

# Expected output should show 7 views:
# - v_active_claims
# - v_investigation_queue
# - v_approval_queue
# - v_sla_breach_report
# - v_payment_queue
# - v_aml_high_risk_alerts
# - v_ombudsman_compliance_pending

# List all functions
\df

# Expected: Multiple functions including:
# - update_updated_at_column()
# - update_claim_search_vector()
# - check_investigation_requirement()
# - calculate_penal_interest()
# - auto_return_pending_documents()
# - calculate_sla_status()
# - validate_workflow_transition()
# - log_claim_status_change()
# - archive_old_claims()
# - cleanup_old_audit_logs()

# List all triggers
SELECT tgname FROM pg_trigger WHERE tgrelid::regclass IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') ORDER BY tgname;

# List all policies (RLS)
SELECT schemaname, tablename, policyname FROM pg_policies;

# Expected: RLS policies on claims, claim_payments, aml_alerts, claim_history

# Check schema version
SELECT * FROM schema_versions;

# Expected:
# version | description | applied_at | applied_by
# 1.0.0 | Initial claims database schema... | [timestamp] | system
```

## Schema Structure

### Tables

1. **claims** - Master table for all claim types (Death, Maturity, Survival Benefit, Freelook)
   - Partitioned by `created_at` (yearly partitions: 2024, 2025, 2026, default)
   - Supports all claim workflows with SLA tracking

2. **claim_documents** - Documents uploaded for claim processing
   - Partitioned by `uploaded_at` (yearly partitions)
   - Supports OCR data extraction
   - Virus scanning integration

3. **investigations** - Investigation workflow tracking
   - 21-day SLA enforcement
   - Progress tracking with heartbeat updates

4. **investigation_progress** - Heartbeat updates for investigations
   - Daily progress updates from investigators

5. **appeals** - Appeal workflow for rejected claims
   - 90-day appeal window
   - 45-day decision timeline

6. **aml_alerts** - AML/CFT alert detection and tracking
   - 70+ trigger rules support
   - STR/CTR filing tracking

7. **claim_payments** - Payment disbursement tracking
   - Partitioned by `created_at` (yearly partitions)
   - Support for NEFT, POSB_TRANSFER, CHEQUE

8. **claim_history** - Complete audit trail
   - All changes logged with old/new values
   - Override tracking with digital signatures

9. **claim_communications** - Multi-channel communication log
   - SMS, Email, WhatsApp, Push, Postal

10. **document_checklist_templates** - Dynamic document checklist
    - Context-aware based on claim type, death type, nomination

11. **claim_sla_tracking** - Real-time SLA monitoring
    - Color-coded alerts (GREEN, YELLOW, ORANGE, RED)

12. **ombudsman_complaints** - Insurance Ombudsman complaint lifecycle
    - Admissibility checks
    - Award compliance tracking (30 days)

13. **policy_bond_tracking** - Policy bond dispatch and delivery
    - Free look period calculation
    - Delivery failure escalation

14. **freelook_cancellations** - Free look cancellation and refund
    - Maker-checker workflow
    - Pro-rata refund calculation

### Key Features

#### Partitioning
- **claims**: Yearly partitions by `created_at`
- **claim_documents**: Yearly partitions by `uploaded_at`
- **claim_payments**: Yearly partitions by `created_at`

Current partitions: 2024, 2025, 2026, default

#### Row-Level Security (RLS)
- Claims table: User-based access control
- Payments: Admin and finance officers only
- AML alerts: Compliance officers only
- Audit trail: Read-only for auditors

#### Business Rules Enforcement
- **BR-CLM-DC-001**: Investigation trigger (3-year rule)
- **BR-CLM-DC-002**: Investigation SLA (21 days)
- **BR-CLM-DC-003**: SLA without investigation (15 days)
- **BR-CLM-DC-004**: SLA with investigation (45 days)
- **BR-CLM-DC-009**: Penal interest calculation (8% p.a.)
- **BR-CLM-DC-021**: SLA color coding system
- And 60+ more business rules

#### Triggers
- Auto-update `updated_at` timestamp
- Auto-increment `version` on updates
- Full-text search vector update
- Investigation requirement check
- Status change audit logging
- Workflow state validation

#### Functions
- `calculate_penal_interest()` - Calculate 8% p.a. penal interest on SLA breach
- `calculate_sla_status()` - Determine SLA color (GREEN/YELLOW/ORANGE/RED)
- `check_investigation_requirement()` - Auto-detect investigation need
- `auto_return_pending_documents()` - Return claims pending >22 days
- `archive_old_claims()` - Archive claims closed >7 years
- `cleanup_old_audit_logs()` - Cleanup audit logs >10 years

## Testing the Schema

### Test Data Insertion

```sql
-- Insert a test claim
INSERT INTO claims (
    claim_number,
    claim_type,
    policy_id,
    customer_id,
    claim_date,
    death_date,
    death_place,
    death_type,
    claimant_name,
    claimant_type,
    claimant_phone,
    claimant_email,
    sla_due_date,
    created_by,
    updated_by
) VALUES (
    'CLM20250001',
    'DEATH',
    uuid_generate_v4(),
    uuid_generate_v4(),
    CURRENT_DATE,
    CURRENT_DATE - INTERVAL '30 days',
    'Mumbai',
    'NATURAL',
    'John Doe',
    'NOMINEE',
    '+919876543210',
    'john.doe@example.com',
    NOW() + INTERVAL '15 days',
    uuid_generate_v4(),
    uuid_generate_v4()
);

-- Verify insertion
SELECT * FROM claims WHERE claim_number = 'CLM20250001';

-- Test trigger (search_vector should be auto-populated)
SELECT claim_number, claimant_name, search_vector FROM claims WHERE claim_number = 'CLM20250001';
```

### Test Views

```sql
-- View active claims
SELECT * FROM v_active_claims LIMIT 10;

-- View investigation queue
SELECT * FROM v_investigation_queue;

-- View approval queue
SELECT * FROM v_approval_queue;

-- View SLA breach report
SELECT * FROM v_sla_breach_report;

-- View payment queue
SELECT * FROM v_payment_queue;

-- View high-risk AML alerts
SELECT * FROM v_aml_high_risk_alerts;

-- View ombudsman compliance pending
SELECT * FROM v_ombudsman_compliance_pending;
```

## Maintenance

### Partition Management (Future Years)

Create partitions for future years as needed:

```sql
-- Example: Create 2027 partition
CREATE TABLE claims_2027 PARTITION OF claims
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');

CREATE TABLE claim_documents_2027 PARTITION OF claim_documents
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');

CREATE TABLE claim_payments_2027 PARTITION OF claim_payments
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');
```

### Vacuum and Analyze

```sql
-- Run vacuum and analyze periodically
VACUUM ANALYZE claims;
VACUUM ANALYZE claim_documents;
VACUUM ANALYZE investigations;
VACUUM ANALYZE aml_alerts;
VACUUM ANALYZE claim_payments;
```

### Archive Old Data

```sql
-- Archive claims closed more than 7 years ago
SELECT archive_old_claims();

-- Cleanup audit logs older than 10 years
SELECT cleanup_old_audit_logs();
```

## Schema Versioning

The `schema_versions` table tracks all schema changes:

```sql
SELECT * FROM schema_versions ORDER BY applied_at DESC;
```

## Connection String Examples

### Go Application (n-api-db)
```go
// configs/config.yaml
db:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: claims_db
  sslmode: disable
```

### Direct PostgreSQL Connection
```
postgresql://postgres:password@localhost:5432/claims_db
```

## Backup and Restore

### Backup
```bash
pg_dump -h localhost -U postgres -d claims_db -F c -f claims_db_backup.dump
```

### Restore
```bash
pg_restore -h localhost -U postgres -d claims_db -c claims_db_backup.dump
```

## Troubleshooting

### Issue: Partition does not exist for new data
**Solution**: Create partition for the required year

### Issue: RLS blocking access
**Solution**: Set user context before queries
```sql
SET app.current_user_id = 'your-user-id';
SET app.user_role = 'ADMIN';
```

### Issue: Permission denied on tables
**Solution**: Grant required permissions
```sql
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO claims_service_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO claims_service_user;
```

## Additional Resources

- PostgreSQL 16 Documentation: https://www.postgresql.org/docs/16/
- n-api-db Documentation: `seed/tool-docs/db-README.md`
- Business Rules: `.zenflow/tasks/code-gen-54c7/requirements.md`
- API Specification: `seed/swagger/`

## Notes

1. This schema is designed for PostgreSQL 16+ features (partitioning, RLS, JSONB)
2. All partitioned tables include a default partition for data outside defined ranges
3. RLS is enabled but commented-out role grants need to be uncommented for production
4. Document checklist templates are seeded with base death claim documents
5. The schema version is tracked in `schema_versions` table

---

**Last Updated**: 2026-01-19
**Schema Version**: 1.0.0
