-- ============================================
-- Claims Database Performance Optimization Patch
-- PostgreSQL Version: 16
-- Purpose: Performance improvements based on production considerations
-- Apply AFTER: claims_schema_enhancement_patch.sql
-- Version: 1.0.2
-- ============================================

BEGIN;

-- ============================================
-- 1. FIX PARTITION STRATEGY
-- ============================================

-- Issue: claim_payments was partitioned by created_at instead of payment_date
-- This causes inefficient partition pruning for payment queries

-- Drop existing partitions
DROP TABLE IF EXISTS claim_payments_2024 CASCADE;
DROP TABLE IF EXISTS claim_payments_2025 CASCADE;
DROP TABLE IF EXISTS claim_payments_2026 CASCADE;
DROP TABLE IF EXISTS claim_payments_default CASCADE;

-- Drop and recreate claim_payments with correct partitioning
ALTER TABLE claim_payments RENAME TO claim_payments_old;

CREATE TABLE claim_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id VARCHAR(50) UNIQUE NOT NULL,
    claim_id UUID NOT NULL,
    payment_amount NUMERIC(15,2) NOT NULL,
    payment_mode payment_mode_enum NOT NULL,
    payment_reference VARCHAR(100),
    utr_number VARCHAR(50),
    transaction_id VARCHAR(100),
    beneficiary_account_number VARCHAR(30) NOT NULL,
    beneficiary_ifsc_code VARCHAR(11) NOT NULL,
    beneficiary_name VARCHAR(200) NOT NULL,
    beneficiary_bank_name VARCHAR(100),
    initiated_by UUID NOT NULL,
    initiated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    payment_date TIMESTAMP WITH TIME ZONE,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    failure_reason TEXT,
    retry_count INTEGER DEFAULT 0,
    reconciliation_status VARCHAR(20) DEFAULT 'PENDING',
    reconciled_at TIMESTAMP WITH TIME ZONE,
    voucher_number VARCHAR(50),
    voucher_date DATE,
    payment_remarks TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_payment_amount_positive CHECK (payment_amount > 0),
    CONSTRAINT chk_retry_count_limit CHECK (retry_count <= 3)
) PARTITION BY RANGE (COALESCE(payment_date, created_at));

COMMENT ON TABLE claim_payments IS 'Payment disbursement tracking with reconciliation - PARTITIONED BY payment_date for optimal query performance';
COMMENT ON COLUMN claim_payments.payment_mode IS 'BR-CLM-DC-017: NEFT > POSB_TRANSFER > CHEQUE priority';

-- Create partitions for claim_payments (by payment_date)
CREATE TABLE claim_payments_2024 PARTITION OF claim_payments
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE claim_payments_2025 PARTITION OF claim_payments
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE claim_payments_2026 PARTITION OF claim_payments
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE claim_payments_default PARTITION OF claim_payments DEFAULT;

-- Migrate data if exists
INSERT INTO claim_payments SELECT * FROM claim_payments_old;
DROP TABLE claim_payments_old;

-- ============================================
-- 2. ADD COMPOUND INDEXES FOR PERFORMANCE
-- ============================================

-- Index 1: Status + Created date (for reporting and dashboards)
CREATE INDEX idx_claims_status_created ON claims(status, created_at DESC)
    WHERE deleted_at IS NULL;

COMMENT ON INDEX idx_claims_status_created IS
    'Dashboard queries: Claims by status over time, trending analysis';

-- Index 2: Policy + Status (for policy-level queries)
CREATE INDEX idx_claims_policy_status ON claims(policy_id, status)
    WHERE deleted_at IS NULL;

COMMENT ON INDEX idx_claims_policy_status IS
    'Policy service integration: All claims for a policy grouped by status';

-- Index 3: Customer + Status (for customer portal)
CREATE INDEX idx_claims_customer_status ON claims(customer_id, status)
    WHERE deleted_at IS NULL;

COMMENT ON INDEX idx_claims_customer_status IS
    'Customer portal: My claims view filtered by status';

-- Index 4: Time-based reporting (critical for dashboards)
CREATE INDEX idx_claims_created_date_type ON claims(
    DATE(created_at), claim_type, status
) WHERE deleted_at IS NULL;

COMMENT ON INDEX idx_claims_created_date_type IS
    'Daily/monthly reports: Claims registered per day by type and status';

-- Index 5: Payment queue optimization
CREATE INDEX idx_claims_approval_date_status ON claims(approval_date DESC, status)
    WHERE status IN ('APPROVED', 'DISBURSEMENT_PENDING');

COMMENT ON INDEX idx_claims_approval_date_status IS
    'Payment queue: Recently approved claims awaiting disbursement';

-- Index 6: Investigation queue optimization
CREATE INDEX idx_investigations_assigned_status ON investigations(
    assignment_date, due_date, status
) WHERE status NOT IN ('COMPLETED', 'CANCELLED');

COMMENT ON INDEX idx_investigations_assigned_status IS
    'Investigation queue: Active investigations sorted by assignment and due date';

-- Index 7: SLA breach monitoring (critical for alerts)
CREATE INDEX idx_claims_sla_breach_status ON claims(sla_breached, sla_due_date)
    WHERE sla_breached = TRUE AND status NOT IN ('PAID', 'CLOSED', 'REJECTED');

COMMENT ON INDEX idx_claims_sla_breach_status IS
    'SLA breach alerts: Active claims that breached SLA, sorted by due date';

-- Index 8: Claim type + SLA status (for prioritization)
CREATE INDEX idx_claims_type_sla_status ON claims(claim_type, sla_status, created_at)
    WHERE status IN ('APPROVAL_PENDING', 'INVESTIGATION_PENDING') AND deleted_at IS NULL;

COMMENT ON INDEX idx_claims_type_sla_status IS
    'Priority queues: Claims grouped by type and SLA urgency';

-- Index 9: Approver workload (for load balancing)
CREATE INDEX idx_claims_approver_pending ON claims(approver_id, status, sla_due_date)
    WHERE status = 'APPROVAL_PENDING' AND deleted_at IS NULL;

COMMENT ON INDEX idx_claims_approver_pending IS
    'Approver dashboard: Workload distribution and SLA monitoring';

-- Index 10: Payment reconciliation
CREATE INDEX idx_claim_payments_recon ON claim_payments(
    reconciliation_status, payment_date
) WHERE reconciliation_status = 'PENDING';

COMMENT ON INDEX idx_claim_payments_recon IS
    'Daily reconciliation: Pending payments sorted by payment date';

-- Index 11: Document verification queue
CREATE INDEX idx_claim_documents_verification ON claim_documents(
    claim_id, verified, uploaded_at
) WHERE verified = FALSE AND deleted_at IS NULL;

COMMENT ON INDEX idx_claim_documents_verification IS
    'Document verification queue: Unverified documents per claim';

-- Index 12: AML alert review queue
CREATE INDEX idx_aml_alerts_review ON aml_alerts(
    risk_level, alert_status, created_at
) WHERE alert_status IN ('FLAGGED', 'UNDER_REVIEW');

COMMENT ON INDEX idx_aml_alerts_review IS
    'AML review queue: Active alerts sorted by risk level and age';

-- ============================================
-- 3. REMOVE DUPLICATE ENUM TYPE
-- ============================================

-- Issue: investigation_status_enum and investigation_outcome_enum are identical
-- Solution: Use single enum type

-- Update investigations table to use investigation_outcome_enum
ALTER TABLE investigations
    ALTER COLUMN investigation_outcome TYPE investigation_outcome_enum;

-- Drop the duplicate enum (if not used elsewhere)
-- Note: This will fail if investigation_status_enum is still in use
DO $$
BEGIN
    -- First, update any columns using investigation_status_enum
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'claims'
        AND column_name = 'investigation_status'
        AND udt_name = 'investigation_status_enum'
    ) THEN
        ALTER TABLE claims
            ALTER COLUMN investigation_status TYPE investigation_outcome_enum
            USING investigation_status::text::investigation_outcome_enum;
    END IF;

    -- Now drop the duplicate enum
    DROP TYPE IF EXISTS investigation_status_enum CASCADE;

    RAISE NOTICE 'Duplicate enum investigation_status_enum removed. Using investigation_outcome_enum.';
END $$;

-- Update comments
COMMENT ON COLUMN claims.investigation_status IS
    'BR-CLM-DC-012: Investigation status - CLEAR, SUSPECT, or FRAUD (uses investigation_outcome_enum)';

-- ============================================
-- 4. OPTIMIZE PENAL INTEREST TRIGGER
-- ============================================

-- Issue: Trigger was not created in original schema, add optimized version

-- Drop existing trigger if exists
DROP TRIGGER IF EXISTS trg_calculate_penal_interest ON claims;

-- Create optimized function that only calculates when needed
CREATE OR REPLACE FUNCTION auto_calculate_penal_interest()
RETURNS TRIGGER AS $$
DECLARE
    v_penal_interest NUMERIC;
BEGIN
    -- Only calculate penal interest when:
    -- 1. SLA is newly breached (wasn't breached before)
    -- 2. Claim is not in final status
    -- 3. Payment date is set (actual settlement happened)

    IF NEW.sla_breached = TRUE AND
       (OLD.sla_breached IS FALSE OR OLD.sla_breached IS NULL) AND
       NEW.status NOT IN ('PAID', 'CLOSED', 'REJECTED') THEN

        -- Calculate penal interest if we have necessary dates
        IF NEW.disbursement_date IS NOT NULL AND NEW.sla_due_date IS NOT NULL THEN
            v_penal_interest := calculate_penal_interest(
                NEW.approved_amount,
                NEW.sla_due_date,
                NEW.disbursement_date
            );

            NEW.penal_interest := v_penal_interest;

            -- Log to audit trail
            INSERT INTO claim_history (
                claim_id,
                action_type,
                action_description,
                old_values,
                new_values,
                performed_by,
                performed_at
            ) VALUES (
                NEW.id,
                'PENAL_INTEREST_CALCULATED',
                'BR-CLM-DC-009: Penal interest auto-calculated due to SLA breach',
                jsonb_build_object('penal_interest', OLD.penal_interest),
                jsonb_build_object('penal_interest', v_penal_interest),
                NEW.updated_by,
                NOW()
            );
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

COMMENT ON FUNCTION auto_calculate_penal_interest() IS
    'BR-CLM-DC-009: Auto-calculates 8% p.a. penal interest only when SLA is first breached';

-- Create trigger with WHEN condition for performance
CREATE TRIGGER trg_calculate_penal_interest
    BEFORE UPDATE ON claims
    FOR EACH ROW
    WHEN (
        -- Only fire when SLA status changes to breached
        NEW.sla_breached = TRUE AND
        (OLD.sla_breached IS FALSE OR OLD.sla_breached IS NULL) AND
        NEW.status NOT IN ('PAID', 'CLOSED', 'REJECTED')
    )
    EXECUTE FUNCTION auto_calculate_penal_interest();

COMMENT ON TRIGGER trg_calculate_penal_interest ON claims IS
    'BR-CLM-DC-009: Calculates 8% p.a. penal interest only when SLA is first breached (optimized with WHEN condition)';

-- ============================================
-- 5. ADD MISSING COVERING INDEXES
-- ============================================

-- Covering index for approval queue (includes all displayed columns)
CREATE INDEX idx_claims_approval_queue_covering ON claims(
    status, sla_status, created_at
) INCLUDE (claim_number, claim_type, claim_amount, investigation_status)
  WHERE status = 'APPROVAL_PENDING' AND deleted_at IS NULL;

COMMENT ON INDEX idx_claims_approval_queue_covering IS
    'Index-only scan for approval queue view with all display columns';

-- Covering index for payment queue
CREATE INDEX idx_claims_payment_queue_covering ON claims(
    status, approval_date
) INCLUDE (claim_number, approved_amount, bank_account_number, bank_ifsc_code)
  WHERE status = 'DISBURSEMENT_PENDING' AND deleted_at IS NULL;

COMMENT ON INDEX idx_claims_payment_queue_covering IS
    'Index-only scan for payment queue with bank details';

-- ============================================
-- 6. ADD STATISTICAL INDEXES FOR ANALYTICS
-- ============================================

-- Index for daily claim counts
CREATE INDEX idx_claims_daily_stats ON claims(
    DATE(created_at), claim_type
) WHERE deleted_at IS NULL;

-- Index for monthly settlement analysis
CREATE INDEX idx_claims_settlement_stats ON claims(
    DATE_TRUNC('month', disbursement_date), claim_type
) WHERE disbursement_date IS NOT NULL;

-- Index for SLA compliance reporting
CREATE INDEX idx_claims_sla_compliance ON claims(
    DATE_TRUNC('month', created_at), sla_breached
) WHERE deleted_at IS NULL;

-- ============================================
-- 7. OPTIMIZE VIEWS FOR BETTER PERFORMANCE
-- ============================================

-- Drop and recreate approval queue view with better index usage
DROP VIEW IF EXISTS v_approval_queue CASCADE;

CREATE OR REPLACE VIEW v_approval_queue AS
SELECT
    c.id as claim_id,
    c.claim_number,
    c.claim_type,
    c.policy_id,
    c.customer_id,
    c.claimant_name,
    c.status,
    c.claim_amount,
    c.investigation_required,
    c.investigation_status,
    c.sla_due_date,
    c.sla_status,
    EXTRACT(DAY FROM (c.sla_due_date - NOW()))::INTEGER as days_until_sla,
    CASE c.sla_status
        WHEN 'RED' THEN 1
        WHEN 'ORANGE' THEN 2
        WHEN 'YELLOW' THEN 3
        ELSE 4
    END as priority_order,
    c.created_at
FROM claims c
WHERE c.status = 'APPROVAL_PENDING'
AND c.deleted_at IS NULL
ORDER BY priority_order ASC, c.created_at ASC;

COMMENT ON VIEW v_approval_queue IS
    'Optimized approval queue using idx_claims_approval_queue_covering for index-only scans';

-- ============================================
-- 8. ADD MATERIALIZED VIEW FOR HEAVY ANALYTICS
-- ============================================

-- Materialized view for daily claim statistics
CREATE MATERIALIZED VIEW mv_daily_claim_stats AS
SELECT
    DATE(created_at) as claim_date,
    claim_type,
    status,
    COUNT(*) as claim_count,
    SUM(claim_amount) as total_amount,
    AVG(claim_amount) as avg_amount,
    COUNT(*) FILTER (WHERE sla_breached = TRUE) as sla_breach_count,
    COUNT(*) FILTER (WHERE investigation_required = TRUE) as investigation_count
FROM claims
WHERE deleted_at IS NULL
GROUP BY DATE(created_at), claim_type, status;

CREATE UNIQUE INDEX idx_mv_daily_claim_stats ON mv_daily_claim_stats(claim_date, claim_type, status);

COMMENT ON MATERIALIZED VIEW mv_daily_claim_stats IS
    'Daily aggregated statistics for dashboards - refresh daily via scheduled job';

-- Function to refresh materialized view
CREATE OR REPLACE FUNCTION refresh_daily_claim_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_daily_claim_stats;
END;
$$ LANGUAGE 'plpgsql';

COMMENT ON FUNCTION refresh_daily_claim_stats() IS
    'Refresh daily claim statistics - schedule to run at midnight';

-- ============================================
-- 9. ADD PARTITION MAINTENANCE FUNCTION
-- ============================================

-- Function to create next year's partitions automatically
CREATE OR REPLACE FUNCTION create_next_year_partitions()
RETURNS void AS $$
DECLARE
    v_next_year INTEGER;
    v_year_start DATE;
    v_year_end DATE;
BEGIN
    v_next_year := EXTRACT(YEAR FROM NOW()) + 1;
    v_year_start := make_date(v_next_year, 1, 1);
    v_year_end := make_date(v_next_year + 1, 1, 1);

    -- Create partitions for claims
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS claims_%s PARTITION OF claims FOR VALUES FROM (%L) TO (%L)',
        v_next_year, v_year_start, v_year_end
    );

    -- Create partitions for claim_documents
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS claim_documents_%s PARTITION OF claim_documents FOR VALUES FROM (%L) TO (%L)',
        v_next_year, v_year_start, v_year_end
    );

    -- Create partitions for claim_payments
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS claim_payments_%s PARTITION OF claim_payments FOR VALUES FROM (%L) TO (%L)',
        v_next_year, v_year_start, v_year_end
    );

    RAISE NOTICE 'Created partitions for year %', v_next_year;
END;
$$ LANGUAGE 'plpgsql';

COMMENT ON FUNCTION create_next_year_partitions() IS
    'Auto-create next year partitions - schedule to run in November';

-- ============================================
-- 10. UPDATE SCHEMA VERSION
-- ============================================

INSERT INTO schema_versions (version, description, applied_by)
VALUES ('1.0.2', 'Performance optimization: compound indexes, partition fix, enum dedup, optimized triggers', 'system')
ON CONFLICT (version) DO NOTHING;

-- ============================================
-- PERFORMANCE ANALYSIS QUERIES
-- ============================================

-- Analyze all tables for query planner
ANALYZE claims;
ANALYZE claim_documents;
ANALYZE claim_payments;
ANALYZE investigations;
ANALYZE aml_alerts;

-- Check partition sizes
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'claims_%' OR tablename LIKE 'claim_documents_%' OR tablename LIKE 'claim_payments_%'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check index usage statistics
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 20;

COMMIT;

-- ============================================
-- MAINTENANCE SCHEDULE RECOMMENDATIONS
-- ============================================

/*
Add these to your cron/scheduler:

1. Daily (midnight):
   - SELECT refresh_daily_claim_stats();
   - VACUUM ANALYZE claims;

2. Weekly (Sunday 2 AM):
   - VACUUM ANALYZE claim_documents;
   - VACUUM ANALYZE claim_payments;
   - REINDEX INDEX CONCURRENTLY idx_claims_status_created;

3. Monthly (1st day, 3 AM):
   - SELECT archive_old_claims();
   - VACUUM FULL claim_history;

4. Yearly (November 1st):
   - SELECT create_next_year_partitions();

5. Monitor continuously:
   - SELECT * FROM pg_stat_user_indexes WHERE idx_scan = 0; -- Unused indexes
   - SELECT * FROM pg_stat_user_tables WHERE n_live_tup > 100000; -- Large tables
*/

-- ============================================
-- PERFORMANCE TESTING QUERIES
-- ============================================

-- Test 1: Approval queue (should use idx_claims_approval_queue_covering)
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM v_approval_queue LIMIT 20;

-- Test 2: Payment queue (should use idx_claims_payment_queue_covering)
EXPLAIN (ANALYZE, BUFFERS)
SELECT claim_number, approved_amount, bank_account_number
FROM claims
WHERE status = 'DISBURSEMENT_PENDING' AND deleted_at IS NULL
ORDER BY approval_date DESC
LIMIT 20;

-- Test 3: Daily statistics (should use idx_claims_daily_stats)
EXPLAIN (ANALYZE, BUFFERS)
SELECT DATE(created_at) as date, claim_type, COUNT(*)
FROM claims
WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(created_at), claim_type;

-- Test 4: Partition pruning (should only scan relevant partition)
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*)
FROM claim_payments
WHERE payment_date BETWEEN '2025-01-01' AND '2025-12-31';

-- ============================================
-- END OF PERFORMANCE OPTIMIZATION PATCH
-- ============================================
