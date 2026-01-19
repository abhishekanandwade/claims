#!/bin/bash
# Schema Verification Script
# This script verifies that all database objects were created correctly

# Database connection parameters
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-claims_db}"

echo "=========================================="
echo "PLI Claims Database Schema Verification"
echo "=========================================="
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo ""

# Function to run SQL and display results
run_sql() {
    local description=$1
    local sql=$2

    echo "--- $description ---"
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$sql"
    echo ""
}

# 1. Verify Tables
echo "=========================================="
echo "1. Verifying Tables (14 expected)"
echo "=========================================="
run_sql "Count tables" "SELECT COUNT(*) as table_count FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';"

run_sql "List tables" "\dt"

# Expected tables
expected_tables=(
    "claims"
    "claim_documents"
    "investigations"
    "investigation_progress"
    "appeals"
    "aml_alerts"
    "claim_payments"
    "claim_history"
    "claim_communications"
    "document_checklist_templates"
    "claim_sla_tracking"
    "ombudsman_complaints"
    "policy_bond_tracking"
    "freelook_cancellations"
    "schema_versions"
)

echo "Expected tables:"
for table in "${expected_tables[@]}"; do
    echo "  - $table"
done
echo ""

# 2. Verify Partitions
echo "=========================================="
echo "2. Verifying Partitions"
echo "=========================================="
run_sql "Claims partitions" "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE 'claims_%' ORDER BY tablename;"
run_sql "Claim documents partitions" "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE 'claim_documents_%' ORDER BY tablename;"
run_sql "Claim payments partitions" "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE 'claim_payments_%' ORDER BY tablename;"
echo ""

# 3. Verify Enums
echo "=========================================="
echo "3. Verifying Enums (12 expected)"
echo "=========================================="
run_sql "Count enums" "SELECT COUNT(*) as enum_count FROM pg_type WHERE typtype = 'e';"
run_sql "List enums" "\dT"
echo ""

# 4. Verify Indexes
echo "=========================================="
echo "4. Verifying Indexes (60+ expected)"
echo "=========================================="
run_sql "Count indexes" "SELECT COUNT(*) as index_count FROM pg_indexes WHERE schemaname = 'public';"
run_sql "Index summary" "SELECT tablename, COUNT(*) as index_count FROM pg_indexes WHERE schemaname = 'public' GROUP BY tablename ORDER BY tablename;"
echo ""

# 5. Verify Views
echo "=========================================="
echo "5. Verifying Views (7 expected)"
echo "=========================================="
run_sql "List views" "\dv"
echo ""

# Expected views
expected_views=(
    "v_active_claims"
    "v_investigation_queue"
    "v_approval_queue"
    "v_sla_breach_report"
    "v_payment_queue"
    "v_aml_high_risk_alerts"
    "v_ombudsman_compliance_pending"
)

echo "Expected views:"
for view in "${expected_views[@]}"; do
    echo "  - $view"
done
echo ""

# 6. Verify Functions
echo "=========================================="
echo "6. Verifying Functions"
echo "=========================================="
run_sql "List functions" "\df"
echo ""

# 7. Verify Triggers
echo "=========================================="
echo "7. Verifying Triggers"
echo "=========================================="
run_sql "Count triggers" "SELECT COUNT(*) as trigger_count FROM pg_trigger WHERE tgrelid::regclass::text IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public');"
run_sql "List triggers" "SELECT tgname, tgrelid::regclass::text as table_name FROM pg_trigger WHERE tgrelid::regclass::text IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') ORDER BY table_name, tgname;"
echo ""

# 8. Verify RLS Policies
echo "=========================================="
echo "8. Verifying Row-Level Security Policies"
echo "=========================================="
run_sql "RLS enabled tables" "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND relrowsecurity = true;"
run_sql "RLS policies" "SELECT schemaname, tablename, policyname FROM pg_policies ORDER BY tablename, policyname;"
echo ""

# 9. Verify Schema Version
echo "=========================================="
echo "9. Verifying Schema Version"
echo "=========================================="
run_sql "Schema version" "SELECT * FROM schema_versions ORDER BY applied_at DESC;"
echo ""

# 10. Test Query
echo "=========================================="
echo "10. Testing Basic Query"
echo "=========================================="
run_sql "Test query on claims" "SELECT COUNT(*) FROM claims;"
echo ""

# 11. Verify Extensions
echo "=========================================="
echo "11. Verifying Extensions"
echo "=========================================="
run_sql "Installed extensions" "SELECT extname, extversion FROM pg_extension ORDER BY extname;"
echo ""

# 12. Verify Constraints
echo "=========================================="
echo "12. Verifying Table Constraints"
echo "=========================================="
run_sql "Foreign keys" "SELECT COUNT(*) as fk_count FROM information_schema.table_constraints WHERE constraint_schema = 'public' AND constraint_type = 'FOREIGN KEY';"
run_sql "Check constraints" "SELECT COUNT(*) as check_count FROM information_schema.table_constraints WHERE constraint_schema = 'public' AND constraint_type = 'CHECK';"
run_sql "Unique constraints" "SELECT COUNT(*) as unique_count FROM information_schema.table_constraints WHERE constraint_schema = 'public' AND constraint_type = 'UNIQUE';"
echo ""

# 13. Verify Seed Data
echo "=========================================="
echo "13. Verifying Seed Data"
echo "=========================================="
run_sql "Document checklist templates" "SELECT COUNT(*) FROM document_checklist_templates;"
run_sql "Sample templates" "SELECT claim_type, death_type, document_type, is_mandatory FROM document_checklist_templates ORDER BY display_order LIMIT 10;"
echo ""

# 14. Verify Table Sizes
echo "=========================================="
echo "14. Table Sizes (Initial State)"
echo "=========================================="
run_sql "Table sizes" "SELECT tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
echo ""

echo "=========================================="
echo "Verification Complete!"
echo "=========================================="
echo ""
echo "Summary:"
echo "- Check that all 14 tables exist"
echo "- Check that partitions are created for 2024, 2025, 2026"
echo "- Check that all 12 enum types exist"
echo "- Check that 60+ indexes are created"
echo "- Check that all 7 views exist"
echo "- Check that functions and triggers are created"
echo "- Check that RLS policies are in place"
echo "- Check that schema version 1.0.0 is recorded"
echo ""
echo "If any counts are lower than expected, review the schema execution output for errors."
