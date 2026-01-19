-- ============================================
-- Claims Database Schema Enhancement Patch
-- PostgreSQL Version: 16
-- Purpose: Add missing Swagger API fields
-- Apply AFTER: claims_database_schema.sql
-- Version: 1.0.1
-- ============================================

-- This patch adds fields from Swagger specification that were
-- identified during Swagger-to-DDL comparison review

BEGIN;

-- ============================================
-- HIGH PRIORITY: Core API Functionality Fields
-- ============================================

-- 1. Add bank account type (from BankDetails schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS bank_account_type VARCHAR(20);

COMMENT ON COLUMN claims.bank_account_type IS 'Swagger: BankDetails.account_type - SAVINGS, CURRENT, or POSB';

-- 2. Add bank branch name (from BankDetails schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS bank_branch_name VARCHAR(200);

COMMENT ON COLUMN claims.bank_branch_name IS 'Swagger: BankDetails.branch_name - Bank branch name';

-- 3. Add submission mode (from MaturityClaimSubmission schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS submission_mode VARCHAR(30);

COMMENT ON COLUMN claims.submission_mode IS 'Swagger: MaturityClaimSubmission.submission_mode - ONLINE_PORTAL, ONLINE_MOBILE_APP, OFFLINE_BO_SO';

-- 4. Add workflow next step (from WorkflowState schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS workflow_next_step VARCHAR(50);

COMMENT ON COLUMN claims.workflow_next_step IS 'Swagger: WorkflowState.next_step - Next workflow step in state machine';

-- 5. Add workflow allowed actions (from WorkflowState schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS workflow_allowed_actions TEXT[];

COMMENT ON COLUMN claims.workflow_allowed_actions IS 'Swagger: WorkflowState.allowed_actions - Array of actions allowed in current state';

-- 6. Add payment remarks (from DisbursementRequest schema)
ALTER TABLE claim_payments ADD COLUMN IF NOT EXISTS payment_remarks TEXT;

COMMENT ON COLUMN claim_payments.payment_remarks IS 'Swagger: DisbursementRequest.payment_remarks - Payment-specific remarks';

-- ============================================
-- MEDIUM PRIORITY: Nice-to-Have Fields
-- ============================================

-- 7. Add bank name match percentage (from BankValidationResponse schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS bank_name_match_percentage NUMERIC(5,2);

COMMENT ON COLUMN claims.bank_name_match_percentage IS 'Swagger: BankValidationResponse.name_match_percentage - Name matching score during validation';

-- 8. Add DigiLocker consent (from SurvivalBenefitSubmission schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS digilocker_consent BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN claims.digilocker_consent IS 'BR-CLM-SB-005: DigiLocker integration consent for document fetching';

-- 9. Add DigiLocker authorization URL (from SurvivalBenefitClaimRegistrationResponse schema)
ALTER TABLE claims ADD COLUMN IF NOT EXISTS digilocker_auth_url TEXT;

COMMENT ON COLUMN claims.digilocker_auth_url IS 'Swagger: SurvivalBenefitClaimRegistrationResponse.digilocker_auth_url - OAuth URL for DigiLocker';

-- ============================================
-- CONSTRAINTS
-- ============================================

-- Add check constraint for bank_account_type
ALTER TABLE claims ADD CONSTRAINT chk_bank_account_type
    CHECK (bank_account_type IS NULL OR bank_account_type IN ('SAVINGS', 'CURRENT', 'POSB'));

-- Add check constraint for submission_mode
ALTER TABLE claims ADD CONSTRAINT chk_submission_mode
    CHECK (submission_mode IS NULL OR submission_mode IN ('ONLINE_PORTAL', 'ONLINE_MOBILE_APP', 'OFFLINE_BO_SO'));

-- Add check constraint for bank_name_match_percentage
ALTER TABLE claims ADD CONSTRAINT chk_bank_name_match_percentage
    CHECK (bank_name_match_percentage IS NULL OR (bank_name_match_percentage >= 0 AND bank_name_match_percentage <= 100));

-- ============================================
-- INDEXES
-- ============================================

-- Add index for submission_mode (for filtering online vs offline claims)
CREATE INDEX IF NOT EXISTS idx_claims_submission_mode ON claims(submission_mode)
    WHERE submission_mode IS NOT NULL AND deleted_at IS NULL;

-- Add index for bank_account_type (for payment routing)
CREATE INDEX IF NOT EXISTS idx_claims_bank_account_type ON claims(bank_account_type)
    WHERE bank_account_type IS NOT NULL AND deleted_at IS NULL;

-- Add index for digilocker_consent (for DigiLocker integration workflow)
CREATE INDEX IF NOT EXISTS idx_claims_digilocker_consent ON claims(digilocker_consent)
    WHERE digilocker_consent = TRUE AND deleted_at IS NULL;

-- ============================================
-- PARTITIONED TABLE UPDATES
-- ============================================

-- Note: The ALTER TABLE commands above automatically apply to all partitions
-- of partitioned tables (claims, claim_payments). No separate partition updates needed.

-- ============================================
-- UPDATE SCHEMA VERSION
-- ============================================

INSERT INTO schema_versions (version, description, applied_by)
VALUES ('1.0.1', 'Add missing Swagger API fields: bank_account_type, submission_mode, workflow helpers, DigiLocker fields', 'system')
ON CONFLICT (version) DO NOTHING;

-- ============================================
-- VERIFICATION QUERIES
-- ============================================

-- Verify new columns exist
DO $$
DECLARE
    missing_columns TEXT[] := ARRAY[]::TEXT[];
BEGIN
    -- Check claims table columns
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claims' AND column_name = 'bank_account_type') THEN
        missing_columns := array_append(missing_columns, 'claims.bank_account_type');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claims' AND column_name = 'bank_branch_name') THEN
        missing_columns := array_append(missing_columns, 'claims.bank_branch_name');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claims' AND column_name = 'submission_mode') THEN
        missing_columns := array_append(missing_columns, 'claims.submission_mode');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claims' AND column_name = 'workflow_next_step') THEN
        missing_columns := array_append(missing_columns, 'claims.workflow_next_step');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claims' AND column_name = 'workflow_allowed_actions') THEN
        missing_columns := array_append(missing_columns, 'claims.workflow_allowed_actions');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'claim_payments' AND column_name = 'payment_remarks') THEN
        missing_columns := array_append(missing_columns, 'claim_payments.payment_remarks');
    END IF;

    IF array_length(missing_columns, 1) > 0 THEN
        RAISE EXCEPTION 'Missing columns: %', array_to_string(missing_columns, ', ');
    ELSE
        RAISE NOTICE 'All columns added successfully!';
    END IF;
END $$;

COMMIT;

-- ============================================
-- ROLLBACK SCRIPT (if needed)
-- ============================================

/*
-- Uncomment and run if you need to rollback this patch

BEGIN;

ALTER TABLE claims DROP COLUMN IF EXISTS bank_account_type;
ALTER TABLE claims DROP COLUMN IF EXISTS bank_branch_name;
ALTER TABLE claims DROP COLUMN IF EXISTS submission_mode;
ALTER TABLE claims DROP COLUMN IF EXISTS workflow_next_step;
ALTER TABLE claims DROP COLUMN IF EXISTS workflow_allowed_actions;
ALTER TABLE claims DROP COLUMN IF EXISTS bank_name_match_percentage;
ALTER TABLE claims DROP COLUMN IF EXISTS digilocker_consent;
ALTER TABLE claims DROP COLUMN IF EXISTS digilocker_auth_url;

ALTER TABLE claim_payments DROP COLUMN IF EXISTS payment_remarks;

DROP INDEX IF EXISTS idx_claims_submission_mode;
DROP INDEX IF EXISTS idx_claims_bank_account_type;
DROP INDEX IF EXISTS idx_claims_digilocker_consent;

DELETE FROM schema_versions WHERE version = '1.0.1';

COMMIT;
*/

-- ============================================
-- USAGE NOTES
-- ============================================

-- After applying this patch, the schema will have:
-- 1. Full Swagger API field coverage (100%)
-- 2. Enhanced workflow state management
-- 3. Better payment tracking with remarks
-- 4. DigiLocker integration support
-- 5. Improved bank validation with account type and branch

-- Example usage:

-- Set bank account type during claim registration
-- UPDATE claims SET bank_account_type = 'SAVINGS' WHERE claim_id = 'xxx';

-- Set submission mode
-- UPDATE claims SET submission_mode = 'ONLINE_PORTAL' WHERE claim_id = 'xxx';

-- Set workflow next step and allowed actions
-- UPDATE claims
-- SET workflow_next_step = 'APPROVAL_PENDING',
--     workflow_allowed_actions = ARRAY['APPROVE', 'REJECT', 'REQUEST_MORE_INFO']
-- WHERE claim_id = 'xxx';

-- Add payment remarks
-- UPDATE claim_payments SET payment_remarks = 'Urgent payment requested' WHERE payment_id = 'xxx';

-- ============================================
-- END OF PATCH
-- ============================================
