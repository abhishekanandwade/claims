-- ============================================
-- PLI Claims Processing Database Schema
-- Database: claims_db
-- PostgreSQL Version: 16
-- Service: Claims Microservice
-- Generated: 2026-01-19
-- ============================================

-- ============================================
-- EXTENSIONS
-- ============================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================
-- ENUM TYPES
-- ============================================

CREATE TYPE claim_type_enum AS ENUM (
    'DEATH',
    'MATURITY',
    'SURVIVAL_BENEFIT',
    'FREELOOK'
);

CREATE TYPE claim_status_enum AS ENUM (
    'REGISTERED',
    'DOCUMENT_PENDING',
    'DOCUMENT_VERIFIED',
    'INVESTIGATION_PENDING',
    'INVESTIGATION_COMPLETED',
    'CALCULATION_COMPLETED',
    'APPROVAL_PENDING',
    'APPROVED',
    'REJECTED',
    'DISBURSEMENT_PENDING',
    'PAID',
    'CLOSED',
    'RETURNED',
    'REOPENED'
);

CREATE TYPE death_type_enum AS ENUM (
    'NATURAL',
    'UNNATURAL',
    'ACCIDENTAL',
    'SUICIDE',
    'HOMICIDE'
);

CREATE TYPE claimant_type_enum AS ENUM (
    'NOMINEE',
    'LEGAL_HEIR',
    'ASSIGNEE'
);

CREATE TYPE payment_mode_enum AS ENUM (
    'AUTO_NEFT',
    'POSB_TRANSFER',
    'CHEQUE'
);

CREATE TYPE investigation_status_enum AS ENUM (
    'CLEAR',
    'SUSPECT',
    'FRAUD'
);

CREATE TYPE investigation_outcome_enum AS ENUM (
    'CLEAR',
    'SUSPECT',
    'FRAUD'
);

CREATE TYPE aml_risk_level_enum AS ENUM (
    'LOW',
    'MEDIUM',
    'HIGH',
    'CRITICAL'
);

CREATE TYPE aml_alert_status_enum AS ENUM (
    'FLAGGED',
    'UNDER_REVIEW',
    'FILED',
    'CLOSED'
);

CREATE TYPE aml_filing_type_enum AS ENUM (
    'STR',
    'CTR',
    'CCR',
    'NTR'
);

CREATE TYPE sla_status_enum AS ENUM (
    'GREEN',
    'YELLOW',
    'ORANGE',
    'RED'
);

-- ============================================
-- MAIN TABLES
-- ============================================

-- E-CLM-DC-001: Claims table (partitioned by created_at)
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(20) UNIQUE NOT NULL,
    claim_type claim_type_enum NOT NULL,
    policy_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    claim_date DATE NOT NULL,
    death_date DATE,
    death_place VARCHAR(255),
    death_type death_type_enum,
    claimant_name VARCHAR(200) NOT NULL,
    claimant_type claimant_type_enum,
    claimant_relation VARCHAR(50),
    claimant_phone VARCHAR(20),
    claimant_email VARCHAR(255),
    status claim_status_enum NOT NULL DEFAULT 'REGISTERED',
    workflow_state VARCHAR(50),
    claim_amount NUMERIC(15,2),
    approved_amount NUMERIC(15,2),
    sum_assured NUMERIC(15,2),
    reversionary_bonus NUMERIC(12,2),
    terminal_bonus NUMERIC(12,2),
    outstanding_loan NUMERIC(12,2),
    unpaid_premiums NUMERIC(12,2),
    penal_interest NUMERIC(10,2),
    investigation_required BOOLEAN DEFAULT FALSE,
    investigation_status investigation_status_enum,
    investigator_id UUID,
    investigation_start_date DATE,
    investigation_completion_date DATE,
    approver_id UUID,
    approval_date TIMESTAMP WITH TIME ZONE,
    approval_remarks TEXT,
    digital_signature_hash VARCHAR(255),
    disbursement_date TIMESTAMP WITH TIME ZONE,
    payment_mode payment_mode_enum,
    payment_reference VARCHAR(100),
    transaction_id VARCHAR(100),
    utr_number VARCHAR(50),
    bank_account_number VARCHAR(30),
    bank_ifsc_code VARCHAR(11),
    bank_account_holder_name VARCHAR(200),
    bank_name VARCHAR(100),
    bank_verified BOOLEAN DEFAULT FALSE,
    bank_verification_method VARCHAR(20),
    rejection_reason TEXT,
    rejection_code VARCHAR(20),
    appeal_submitted BOOLEAN DEFAULT FALSE,
    appeal_id UUID,
    sla_due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    sla_breached BOOLEAN DEFAULT FALSE,
    sla_breach_days INTEGER DEFAULT 0,
    sla_status sla_status_enum DEFAULT 'GREEN',
    closure_date TIMESTAMP WITH TIME ZONE,
    closure_reason VARCHAR(50),
    metadata JSONB,
    search_vector tsvector,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- BR-CLM-DC-001: Investigation trigger validation
    CONSTRAINT chk_death_date_valid CHECK (death_date IS NULL OR death_date <= CURRENT_DATE),
    CONSTRAINT chk_claim_amount_positive CHECK (claim_amount IS NULL OR claim_amount > 0),
    CONSTRAINT chk_approved_amount_positive CHECK (approved_amount IS NULL OR approved_amount > 0),
    CONSTRAINT chk_penal_interest_non_negative CHECK (penal_interest IS NULL OR penal_interest >= 0),
    CONSTRAINT chk_death_fields_for_death_claim CHECK (
        (claim_type = 'DEATH' AND death_date IS NOT NULL) OR
        (claim_type != 'DEATH' AND death_date IS NULL)
    )
) PARTITION BY RANGE (created_at);

COMMENT ON TABLE claims IS 'E-CLM-DC-001: Master table for all claim types (Death, Maturity, Survival Benefit, Freelook)';
COMMENT ON COLUMN claims.investigation_required IS 'BR-CLM-DC-001: Auto-set to TRUE if death within 3 years of policy issue/revival';
COMMENT ON COLUMN claims.penal_interest IS 'BR-CLM-DC-009: 8% p.a. calculated on SLA breach';
COMMENT ON COLUMN claims.sla_due_date IS 'BR-CLM-DC-003/004: 15 days without investigation, 45 days with investigation';

-- Create partitions for claims table (yearly partitions)
CREATE TABLE claims_2024 PARTITION OF claims
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE claims_2025 PARTITION OF claims
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE claims_2026 PARTITION OF claims
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE claims_default PARTITION OF claims DEFAULT;

-- E-CLM-DC-002: Claim Documents table (partitioned by uploaded_date)
CREATE TABLE claim_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_name VARCHAR(255) NOT NULL,
    document_url TEXT NOT NULL,
    ecms_reference_id VARCHAR(100),
    file_size INTEGER NOT NULL,
    file_hash VARCHAR(255),
    content_type VARCHAR(100),
    is_mandatory BOOLEAN DEFAULT FALSE,
    uploaded_by UUID NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    virus_scanned BOOLEAN DEFAULT FALSE,
    virus_scan_status VARCHAR(20),
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID,
    verified_at TIMESTAMP WITH TIME ZONE,
    verification_remarks TEXT,
    ocr_extracted_data JSONB,
    ocr_confidence_score NUMERIC(5,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_file_size_positive CHECK (file_size > 0),
    CONSTRAINT chk_ocr_confidence_range CHECK (ocr_confidence_score IS NULL OR (ocr_confidence_score >= 0 AND ocr_confidence_score <= 100))
) PARTITION BY RANGE (uploaded_at);

COMMENT ON TABLE claim_documents IS 'E-CLM-DC-002: Documents uploaded for claim processing with OCR support';
COMMENT ON COLUMN claim_documents.document_type IS 'BR-CLM-DC-013/014/015: Document types based on death type and nomination status';

CREATE TABLE claim_documents_2024 PARTITION OF claim_documents
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE claim_documents_2025 PARTITION OF claim_documents
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE claim_documents_2026 PARTITION OF claim_documents
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE claim_documents_default PARTITION OF claim_documents DEFAULT;

-- Investigation Assignment table
CREATE TABLE investigations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    investigation_id VARCHAR(50) UNIQUE NOT NULL,
    claim_id UUID NOT NULL,
    assigned_by UUID NOT NULL,
    investigator_id UUID NOT NULL,
    investigator_rank VARCHAR(20),
    jurisdiction VARCHAR(100),
    auto_assigned BOOLEAN DEFAULT FALSE,
    assignment_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ASSIGNED',
    progress_percentage INTEGER DEFAULT 0,
    investigation_outcome investigation_outcome_enum,
    cause_of_death VARCHAR(255),
    cause_of_death_verified BOOLEAN DEFAULT FALSE,
    hospital_records_verified BOOLEAN DEFAULT FALSE,
    detailed_findings TEXT,
    recommendation VARCHAR(50),
    report_document_id UUID,
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_by UUID,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    review_decision VARCHAR(20),
    reviewer_remarks TEXT,
    reinvestigation_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-DC-002: 21-day investigation SLA
    CONSTRAINT chk_progress_range CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    -- BR-CLM-DC-013: Max 2 reinvestigations
    CONSTRAINT chk_reinvestigation_limit CHECK (reinvestigation_count <= 2)
);

COMMENT ON TABLE investigations IS 'Investigation workflow tracking with SLA monitoring';
COMMENT ON COLUMN investigations.due_date IS 'BR-CLM-DC-002: Investigation report due within 21 days';
COMMENT ON COLUMN investigations.reinvestigation_count IS 'BR-CLM-DC-013: Maximum 2 reinvestigations allowed';

-- Investigation Progress Updates table
CREATE TABLE investigation_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    investigation_id UUID NOT NULL,
    update_date DATE NOT NULL,
    progress_percentage INTEGER NOT NULL,
    checklist_items_completed TEXT[],
    remarks TEXT NOT NULL,
    estimated_completion_date DATE,
    updated_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_progress_valid CHECK (progress_percentage >= 0 AND progress_percentage <= 100)
);

COMMENT ON TABLE investigation_progress IS 'Heartbeat updates for long-running investigations';

-- Appeals table
CREATE TABLE appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    appeal_number VARCHAR(50) UNIQUE NOT NULL,
    claim_id UUID NOT NULL,
    appellant_name VARCHAR(200) NOT NULL,
    appellant_contact JSONB,
    grounds_of_appeal TEXT NOT NULL,
    supporting_documents UUID[],
    condonation_request BOOLEAN DEFAULT FALSE,
    condonation_reason TEXT,
    submission_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    appeal_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    appellate_authority_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED',
    decision VARCHAR(30),
    reasoned_order TEXT,
    revised_claim_amount NUMERIC(15,2),
    decision_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-DC-005: 90-day appeal window
    -- BR-CLM-DC-007: 45-day decision timeline
    CONSTRAINT chk_revised_amount_positive CHECK (revised_claim_amount IS NULL OR revised_claim_amount > 0)
);

COMMENT ON TABLE appeals IS 'Appeal workflow for rejected claims';
COMMENT ON COLUMN appeals.submission_date IS 'BR-CLM-DC-005: Must be within 90 days of rejection';
COMMENT ON COLUMN appeals.appeal_deadline IS 'BR-CLM-DC-007: 45-day SLA for appellate decision';

-- E-CLM-AML-001: AML Alerts table
CREATE TABLE aml_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alert_id VARCHAR(50) UNIQUE NOT NULL,
    trigger_code VARCHAR(20) NOT NULL,
    policy_id UUID NOT NULL,
    customer_id UUID,
    transaction_type VARCHAR(50) NOT NULL,
    transaction_amount NUMERIC(15,2),
    transaction_date DATE NOT NULL,
    payment_mode VARCHAR(20),
    risk_level aml_risk_level_enum NOT NULL,
    risk_score INTEGER,
    alert_status aml_alert_status_enum NOT NULL DEFAULT 'FLAGGED',
    alert_description TEXT,
    trigger_details JSONB,
    reviewed_by UUID,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    review_decision VARCHAR(20),
    officer_remarks TEXT,
    action_taken TEXT,
    transaction_blocked BOOLEAN DEFAULT FALSE,
    filing_required BOOLEAN DEFAULT FALSE,
    filing_type aml_filing_type_enum,
    filing_status VARCHAR(20),
    filing_reference VARCHAR(100),
    filed_at TIMESTAMP WITH TIME ZONE,
    filed_by UUID,
    pan_number VARCHAR(10),
    pan_verified BOOLEAN,
    nominee_change_detected BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-AML-001: Cash over 50k triggers CTR
    -- BR-CLM-AML-003: Nominee change post-death blocks transaction
    CONSTRAINT chk_risk_score_range CHECK (risk_score IS NULL OR (risk_score >= 0 AND risk_score <= 100))
);

COMMENT ON TABLE aml_alerts IS 'E-CLM-AML-001: AML/CFT alert detection and tracking';
COMMENT ON COLUMN aml_alerts.trigger_code IS 'BR-CLM-AML-001 to 005: AML_001 to AML_005 trigger codes';
COMMENT ON COLUMN aml_alerts.filing_required IS 'BR-CLM-AML-006/007: STR within 7 days, CTR monthly';

-- Claim Payments table (partitioned by payment_date)
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
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_payment_amount_positive CHECK (payment_amount > 0),
    CONSTRAINT chk_retry_count_limit CHECK (retry_count <= 3)
) PARTITION BY RANGE (created_at);

COMMENT ON TABLE claim_payments IS 'Payment disbursement tracking with reconciliation';
COMMENT ON COLUMN claim_payments.payment_mode IS 'BR-CLM-DC-017: NEFT > POSB_TRANSFER > CHEQUE priority';

CREATE TABLE claim_payments_2024 PARTITION OF claim_payments
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE claim_payments_2025 PARTITION OF claim_payments
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE claim_payments_2026 PARTITION OF claim_payments
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE claim_payments_default PARTITION OF claim_payments DEFAULT;

-- Claim History / Audit Trail table
CREATE TABLE claim_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    action_description TEXT,
    old_status claim_status_enum,
    new_status claim_status_enum,
    old_values JSONB,
    new_values JSONB,
    override_applied BOOLEAN DEFAULT FALSE,
    override_reason TEXT,
    override_field VARCHAR(100),
    override_old_value TEXT,
    override_new_value TEXT,
    digital_signature_hash VARCHAR(255),
    performed_by UUID NOT NULL,
    performed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,

    -- BR-CLM-DC-016: Manual override audit trail
    -- BR-CLM-DC-025: Digital signature for overrides
    CONSTRAINT chk_override_requires_reason CHECK (
        (override_applied = FALSE) OR
        (override_applied = TRUE AND override_reason IS NOT NULL AND digital_signature_hash IS NOT NULL)
    )
);

COMMENT ON TABLE claim_history IS 'Complete audit trail for all claim changes';
COMMENT ON COLUMN claim_history.override_applied IS 'BR-CLM-DC-016/025: Manual override with mandatory remarks and digital signature';

-- Communication Log table
CREATE TABLE claim_communications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    communication_type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    recipient_name VARCHAR(200),
    recipient_mobile VARCHAR(20),
    recipient_email VARCHAR(255),
    template_id VARCHAR(50),
    message_content TEXT,
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    delivery_status VARCHAR(20) DEFAULT 'SENT',
    delivery_timestamp TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    provider_reference VARCHAR(100),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-DC-019: Communication milestone triggers
    CONSTRAINT chk_channel_valid CHECK (channel IN ('SMS', 'EMAIL', 'WHATSAPP', 'PUSH', 'POSTAL'))
);

COMMENT ON TABLE claim_communications IS 'Multi-channel communication log';
COMMENT ON COLUMN claim_communications.communication_type IS 'BR-CLM-DC-019: REGISTRATION, DOCUMENT_STATUS, INVESTIGATION, APPROVAL, PAYMENT';

-- Document Checklist Template table
CREATE TABLE document_checklist_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_type claim_type_enum NOT NULL,
    death_type death_type_enum,
    nomination_status VARCHAR(20),
    policy_type VARCHAR(50),
    document_type VARCHAR(50) NOT NULL,
    document_description TEXT,
    is_mandatory BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INTEGER,
    validation_rules JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-DC-015: Base mandatory documents
    -- BR-CLM-DC-013/014: Conditional documents
    CONSTRAINT uq_checklist_template UNIQUE (claim_type, death_type, nomination_status, document_type)
);

COMMENT ON TABLE document_checklist_templates IS 'Dynamic document checklist based on claim context';
COMMENT ON COLUMN document_checklist_templates.is_mandatory IS 'BR-CLM-DC-015: Base mandatory or conditional based on death_type/nomination';

-- SLA Tracking table
CREATE TABLE claim_sla_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    sla_type VARCHAR(50) NOT NULL,
    sla_start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    sla_due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    sla_total_days INTEGER NOT NULL,
    sla_elapsed_days INTEGER NOT NULL DEFAULT 0,
    sla_remaining_days INTEGER NOT NULL,
    sla_status sla_status_enum NOT NULL DEFAULT 'GREEN',
    sla_breach BOOLEAN DEFAULT FALSE,
    sla_breach_date TIMESTAMP WITH TIME ZONE,
    sla_completion_date TIMESTAMP WITH TIME ZONE,
    escalation_triggered BOOLEAN DEFAULT FALSE,
    escalation_level INTEGER DEFAULT 0,
    last_escalation_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-DC-021: Color-coded SLA system
    CONSTRAINT chk_sla_days_positive CHECK (sla_total_days > 0)
);

COMMENT ON TABLE claim_sla_tracking IS 'Real-time SLA monitoring with color-coded alerts';
COMMENT ON COLUMN claim_sla_tracking.sla_status IS 'BR-CLM-DC-021: GREEN (<70%), YELLOW (70-90%), ORANGE (90-100%), RED (>100%)';

-- Ombudsman Complaints table
CREATE TABLE ombudsman_complaints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    complaint_number VARCHAR(50) UNIQUE NOT NULL,
    claim_id UUID,
    policy_id UUID NOT NULL,
    complainant_name VARCHAR(200) NOT NULL,
    complainant_contact JSONB NOT NULL,
    complaint_description TEXT NOT NULL,
    complaint_category VARCHAR(50),
    claim_value NUMERIC(15,2),
    representation_to_insurer_date DATE,
    wait_period_completed BOOLEAN DEFAULT FALSE,
    limitation_period_valid BOOLEAN DEFAULT FALSE,
    parallel_litigation BOOLEAN DEFAULT FALSE,
    admissible BOOLEAN,
    inadmissibility_reason TEXT,
    ombudsman_center VARCHAR(100),
    jurisdiction_basis VARCHAR(50),
    assigned_ombudsman_id UUID,
    conflict_of_interest BOOLEAN DEFAULT FALSE,
    status VARCHAR(30) NOT NULL DEFAULT 'SUBMITTED',
    mediation_attempted BOOLEAN DEFAULT FALSE,
    mediation_successful BOOLEAN,
    mediation_terms TEXT,
    recommendation_issued BOOLEAN DEFAULT FALSE,
    recommendation_date DATE,
    award_issued BOOLEAN DEFAULT FALSE,
    award_number VARCHAR(50),
    award_amount NUMERIC(15,2),
    award_date DATE,
    award_digitally_signed BOOLEAN DEFAULT FALSE,
    compliance_due_date DATE,
    compliance_status VARCHAR(20),
    compliance_date DATE,
    escalated_to_irdai BOOLEAN DEFAULT FALSE,
    closure_date DATE,
    archival_date DATE,
    retention_period_years INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-OMB-001: Admissibility checks
    -- BR-CLM-OMB-005: 50 lakh cap
    CONSTRAINT chk_claim_value_cap CHECK (claim_value IS NULL OR claim_value <= 5000000),
    CONSTRAINT chk_award_amount_cap CHECK (award_amount IS NULL OR award_amount <= 5000000)
);

COMMENT ON TABLE ombudsman_complaints IS 'Insurance Ombudsman complaint lifecycle management';
COMMENT ON COLUMN ombudsman_complaints.claim_value IS 'BR-CLM-OMB-001: Must be <= ₹50 lakh for admissibility';
COMMENT ON COLUMN ombudsman_complaints.compliance_due_date IS 'BR-CLM-OMB-006: 30 days from award date';

-- Policy Bond Tracking table
CREATE TABLE policy_bond_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL,
    bond_number VARCHAR(50) UNIQUE NOT NULL,
    bond_type VARCHAR(20) NOT NULL,
    print_date DATE,
    dispatch_date DATE,
    tracking_number VARCHAR(50),
    delivery_date DATE,
    delivery_status VARCHAR(30),
    delivery_attempt_count INTEGER DEFAULT 0,
    pod_reference VARCHAR(100),
    recipient_name VARCHAR(200),
    recipient_signature_captured BOOLEAN DEFAULT FALSE,
    undelivered_reason VARCHAR(255),
    escalation_triggered BOOLEAN DEFAULT FALSE,
    escalation_date DATE,
    customer_contacted BOOLEAN DEFAULT FALSE,
    address_verified BOOLEAN DEFAULT FALSE,
    redelivery_requested BOOLEAN DEFAULT FALSE,
    freelook_period_start_date DATE,
    freelook_period_end_date DATE,
    freelook_cancellation_submitted BOOLEAN DEFAULT FALSE,
    freelook_cancellation_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-BOND-001: Free look period calculation
    -- BR-CLM-BOND-002: Delivery failure escalation
    CONSTRAINT chk_bond_type_valid CHECK (bond_type IN ('PHYSICAL', 'ELECTRONIC')),
    CONSTRAINT chk_delivery_attempts_limit CHECK (delivery_attempt_count <= 3)
);

COMMENT ON TABLE policy_bond_tracking IS 'Policy bond dispatch and delivery tracking';
COMMENT ON COLUMN policy_bond_tracking.freelook_period_end_date IS 'BR-CLM-BOND-001: Physical: 15 days from delivery, Electronic: 30 days from issuance';

-- Freelook Cancellation table
CREATE TABLE freelook_cancellations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cancellation_number VARCHAR(50) UNIQUE NOT NULL,
    policy_id UUID NOT NULL,
    bond_tracking_id UUID,
    cancellation_request_date DATE NOT NULL,
    cancellation_reason TEXT NOT NULL,
    freelook_period_valid BOOLEAN DEFAULT FALSE,
    rejection_reason TEXT,
    total_premium NUMERIC(12,2) NOT NULL,
    pro_rata_risk_premium NUMERIC(12,2) NOT NULL,
    stamp_duty NUMERIC(8,2) NOT NULL,
    medical_costs NUMERIC(8,2),
    other_deductions NUMERIC(8,2),
    refund_amount NUMERIC(12,2) NOT NULL,
    maker_id UUID NOT NULL,
    maker_entry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    checker_id UUID,
    checker_verification_date TIMESTAMP WITH TIME ZONE,
    maker_checker_approved BOOLEAN DEFAULT FALSE,
    refund_transaction_id VARCHAR(100),
    refund_status VARCHAR(20) DEFAULT 'PENDING',
    refund_date DATE,
    linked_to_finance BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- BR-CLM-BOND-003: Refund calculation
    -- BR-CLM-BOND-004: Maker-checker workflow
    CONSTRAINT chk_maker_checker_different CHECK (maker_id != checker_id),
    CONSTRAINT chk_refund_amount_valid CHECK (
        refund_amount = total_premium - (pro_rata_risk_premium + stamp_duty + COALESCE(medical_costs, 0) + COALESCE(other_deductions, 0))
    )
);

COMMENT ON TABLE freelook_cancellations IS 'Free look cancellation and refund processing';
COMMENT ON COLUMN freelook_cancellations.refund_amount IS 'BR-CLM-BOND-003: Premium - (risk premium + stamp duty + medical + other)';
COMMENT ON COLUMN freelook_cancellations.maker_id IS 'BR-CLM-BOND-004: Maker-checker segregation of duties';

-- ============================================
-- INDEXES
-- ============================================

-- Claims table indexes
CREATE INDEX idx_claims_policy_id ON claims(policy_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_claims_customer_id ON claims(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_claims_status ON claims(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_claims_claim_type ON claims(claim_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_claims_claim_number ON claims(claim_number);
CREATE INDEX idx_claims_sla_due_date ON claims(sla_due_date) WHERE status IN ('APPROVAL_PENDING', 'INVESTIGATION_PENDING') AND deleted_at IS NULL;
CREATE INDEX idx_claims_sla_status ON claims(sla_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_claims_sla_breached ON claims(sla_breached) WHERE sla_breached = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_claims_death_date ON claims(death_date) WHERE death_date IS NOT NULL;
CREATE INDEX idx_claims_investigation_required ON claims(investigation_required) WHERE investigation_required = TRUE;
CREATE INDEX idx_claims_created_at ON claims(created_at);
CREATE INDEX idx_claims_approval_pending ON claims(status, approved_amount) WHERE status = 'APPROVAL_PENDING' AND deleted_at IS NULL;
CREATE INDEX idx_claims_disbursement_pending ON claims(status) WHERE status = 'DISBURSEMENT_PENDING' AND deleted_at IS NULL;
CREATE INDEX idx_claims_metadata ON claims USING gin(metadata);
CREATE INDEX idx_claims_search_vector ON claims USING gin(search_vector);

-- Composite index for common queries
CREATE INDEX idx_claims_status_type_sla ON claims(status, claim_type, sla_status) WHERE deleted_at IS NULL;

-- Claim Documents indexes
CREATE INDEX idx_claim_documents_claim_id ON claim_documents(claim_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_claim_documents_type ON claim_documents(document_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_claim_documents_verified ON claim_documents(verified) WHERE verified = FALSE;
CREATE INDEX idx_claim_documents_uploaded_at ON claim_documents(uploaded_at);
CREATE INDEX idx_claim_documents_virus_scan ON claim_documents(virus_scanned) WHERE virus_scanned = FALSE;
CREATE INDEX idx_claim_documents_ocr ON claim_documents USING gin(ocr_extracted_data);

-- Investigations indexes
CREATE INDEX idx_investigations_claim_id ON investigations(claim_id);
CREATE INDEX idx_investigations_investigator_id ON investigations(investigator_id);
CREATE INDEX idx_investigations_status ON investigations(status);
CREATE INDEX idx_investigations_due_date ON investigations(due_date) WHERE status NOT IN ('COMPLETED', 'CANCELLED');
CREATE INDEX idx_investigations_jurisdiction ON investigations(jurisdiction);
CREATE INDEX idx_investigations_outcome ON investigations(investigation_outcome) WHERE investigation_outcome IS NOT NULL;

-- Investigation Progress indexes
CREATE INDEX idx_investigation_progress_investigation_id ON investigation_progress(investigation_id);
CREATE INDEX idx_investigation_progress_update_date ON investigation_progress(update_date);

-- Appeals indexes
CREATE INDEX idx_appeals_claim_id ON appeals(claim_id);
CREATE INDEX idx_appeals_status ON appeals(status);
CREATE INDEX idx_appeals_submission_date ON appeals(submission_date);
CREATE INDEX idx_appeals_appeal_deadline ON appeals(appeal_deadline) WHERE status = 'SUBMITTED';

-- AML Alerts indexes
CREATE INDEX idx_aml_alerts_policy_id ON aml_alerts(policy_id);
CREATE INDEX idx_aml_alerts_customer_id ON aml_alerts(customer_id);
CREATE INDEX idx_aml_alerts_status ON aml_alerts(alert_status);
CREATE INDEX idx_aml_alerts_risk_level ON aml_alerts(risk_level);
CREATE INDEX idx_aml_alerts_trigger_code ON aml_alerts(trigger_code);
CREATE INDEX idx_aml_alerts_transaction_date ON aml_alerts(transaction_date);
CREATE INDEX idx_aml_alerts_filing_required ON aml_alerts(filing_required) WHERE filing_required = TRUE;
CREATE INDEX idx_aml_alerts_pan ON aml_alerts(pan_number) WHERE pan_number IS NOT NULL;
CREATE INDEX idx_aml_alerts_trigger_details ON aml_alerts USING gin(trigger_details);

-- Claim Payments indexes
CREATE INDEX idx_claim_payments_claim_id ON claim_payments(claim_id);
CREATE INDEX idx_claim_payments_status ON claim_payments(payment_status);
CREATE INDEX idx_claim_payments_payment_date ON claim_payments(payment_date) WHERE payment_date IS NOT NULL;
CREATE INDEX idx_claim_payments_utr ON claim_payments(utr_number) WHERE utr_number IS NOT NULL;
CREATE INDEX idx_claim_payments_reconciliation ON claim_payments(reconciliation_status) WHERE reconciliation_status = 'PENDING';
CREATE INDEX idx_claim_payments_created_at ON claim_payments(created_at);

-- Claim History indexes
CREATE INDEX idx_claim_history_claim_id ON claim_history(claim_id);
CREATE INDEX idx_claim_history_performed_at ON claim_history(performed_at);
CREATE INDEX idx_claim_history_performed_by ON claim_history(performed_by);
CREATE INDEX idx_claim_history_override ON claim_history(override_applied) WHERE override_applied = TRUE;

-- Communications indexes
CREATE INDEX idx_claim_communications_claim_id ON claim_communications(claim_id);
CREATE INDEX idx_claim_communications_sent_at ON claim_communications(sent_at);
CREATE INDEX idx_claim_communications_delivery_status ON claim_communications(delivery_status);
CREATE INDEX idx_claim_communications_channel ON claim_communications(channel);

-- SLA Tracking indexes
CREATE INDEX idx_claim_sla_tracking_claim_id ON claim_sla_tracking(claim_id);
CREATE INDEX idx_claim_sla_tracking_status ON claim_sla_tracking(sla_status);
CREATE INDEX idx_claim_sla_tracking_breach ON claim_sla_tracking(sla_breach) WHERE sla_breach = TRUE;
CREATE INDEX idx_claim_sla_tracking_due_date ON claim_sla_tracking(sla_due_date);

-- Ombudsman Complaints indexes
CREATE INDEX idx_ombudsman_complaints_claim_id ON ombudsman_complaints(claim_id);
CREATE INDEX idx_ombudsman_complaints_policy_id ON ombudsman_complaints(policy_id);
CREATE INDEX idx_ombudsman_complaints_status ON ombudsman_complaints(status);
CREATE INDEX idx_ombudsman_complaints_compliance ON ombudsman_complaints(compliance_status) WHERE compliance_status != 'COMPLETED';

-- Policy Bond Tracking indexes
CREATE INDEX idx_policy_bond_tracking_policy_id ON policy_bond_tracking(policy_id);
CREATE INDEX idx_policy_bond_tracking_tracking_number ON policy_bond_tracking(tracking_number);
CREATE INDEX idx_policy_bond_tracking_delivery_status ON policy_bond_tracking(delivery_status);
CREATE INDEX idx_policy_bond_tracking_escalation ON policy_bond_tracking(escalation_triggered) WHERE escalation_triggered = TRUE;

-- Freelook Cancellations indexes
CREATE INDEX idx_freelook_cancellations_policy_id ON freelook_cancellations(policy_id);
CREATE INDEX idx_freelook_cancellations_refund_status ON freelook_cancellations(refund_status);
CREATE INDEX idx_freelook_cancellations_maker_checker ON freelook_cancellations(maker_checker_approved) WHERE maker_checker_approved = FALSE;

-- ============================================
-- FUNCTIONS
-- ============================================

-- Function to update updated_at timestamp and version
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    IF TG_TABLE_NAME IN ('claims', 'investigations', 'appeals', 'ombudsman_complaints') THEN
        NEW.version = OLD.version + 1;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Function to update full-text search vector
CREATE OR REPLACE FUNCTION update_claim_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.claim_number, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.claimant_name, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(NEW.death_place, '')), 'C');
    RETURN NEW;
END;
$$ language 'plpgsql';

-- BR-CLM-DC-001: Function to check investigation requirement
CREATE OR REPLACE FUNCTION check_investigation_requirement()
RETURNS TRIGGER AS $$
DECLARE
    policy_issue_date DATE;
    policy_revival_date DATE;
BEGIN
    IF NEW.claim_type = 'DEATH' AND NEW.death_date IS NOT NULL THEN
        -- Fetch policy dates (would be from policy service in real implementation)
        -- SELECT issue_date, last_revival_date INTO policy_issue_date, policy_revival_date
        -- FROM policies WHERE id = NEW.policy_id;

        -- For now, setting investigation_required based on death_type
        IF NEW.death_type IN ('UNNATURAL', 'ACCIDENTAL', 'SUICIDE', 'HOMICIDE') THEN
            NEW.investigation_required := TRUE;
        END IF;

        -- Check 3-year rule (would need policy dates)
        -- IF (NEW.death_date - policy_issue_date) <= INTERVAL '3 years' OR
        --    (policy_revival_date IS NOT NULL AND (NEW.death_date - policy_revival_date) <= INTERVAL '3 years') THEN
        --     NEW.investigation_required := TRUE;
        -- END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- BR-CLM-DC-009: Function to calculate penal interest
CREATE OR REPLACE FUNCTION calculate_penal_interest(
    p_claim_amount NUMERIC,
    p_sla_due_date TIMESTAMP WITH TIME ZONE,
    p_actual_settlement_date TIMESTAMP WITH TIME ZONE
) RETURNS NUMERIC AS $$
DECLARE
    v_days_delayed INTEGER;
    v_penal_interest NUMERIC;
BEGIN
    IF p_actual_settlement_date <= p_sla_due_date THEN
        RETURN 0;
    END IF;

    v_days_delayed := EXTRACT(DAY FROM (p_actual_settlement_date - p_sla_due_date));
    v_penal_interest := (p_claim_amount * 0.08 * v_days_delayed) / 365;

    RETURN ROUND(v_penal_interest, 2);
END;
$$ LANGUAGE 'plpgsql' IMMUTABLE;

-- BR-CLM-DC-010: Function to auto-return pending claims
CREATE OR REPLACE FUNCTION auto_return_pending_documents()
RETURNS TABLE(claim_id UUID, claim_number VARCHAR, days_pending INTEGER) AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.claim_number,
        EXTRACT(DAY FROM (NOW() - c.created_at))::INTEGER
    FROM claims c
    WHERE c.status = 'DOCUMENT_PENDING'
    AND c.deleted_at IS NULL
    AND EXTRACT(DAY FROM (NOW() - c.created_at)) > 22;
END;
$$ LANGUAGE 'plpgsql' STABLE;

-- BR-CLM-DC-021: Function to calculate SLA status
CREATE OR REPLACE FUNCTION calculate_sla_status(
    p_sla_start_date TIMESTAMP WITH TIME ZONE,
    p_sla_due_date TIMESTAMP WITH TIME ZONE
) RETURNS sla_status_enum AS $$
DECLARE
    v_total_duration NUMERIC;
    v_elapsed_duration NUMERIC;
    v_consumption_percentage NUMERIC;
BEGIN
    v_total_duration := EXTRACT(EPOCH FROM (p_sla_due_date - p_sla_start_date));
    v_elapsed_duration := EXTRACT(EPOCH FROM (NOW() - p_sla_start_date));
    v_consumption_percentage := (v_elapsed_duration / v_total_duration);

    IF v_consumption_percentage < 0.70 THEN
        RETURN 'GREEN';
    ELSIF v_consumption_percentage < 0.90 THEN
        RETURN 'YELLOW';
    ELSIF v_consumption_percentage <= 1.00 THEN
        RETURN 'ORANGE';
    ELSE
        RETURN 'RED';
    END IF;
END;
$$ LANGUAGE 'plpgsql' IMMUTABLE;

-- Function to validate workflow state transition
CREATE OR REPLACE FUNCTION validate_workflow_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Add workflow state validation logic here
    -- This would check if the status transition is valid
    -- For example: REGISTERED -> DOCUMENT_PENDING -> DOCUMENT_VERIFIED -> etc.

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Function to log claim status changes
CREATE OR REPLACE FUNCTION log_claim_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO claim_history (
            claim_id,
            action_type,
            action_description,
            old_status,
            new_status,
            old_values,
            new_values,
            performed_by
        ) VALUES (
            NEW.id,
            'STATUS_CHANGE',
            'Claim status changed from ' || OLD.status || ' to ' || NEW.status,
            OLD.status,
            NEW.status,
            to_jsonb(OLD),
            to_jsonb(NEW),
            NEW.updated_by
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- ============================================
-- TRIGGERS
-- ============================================

-- Triggers for updated_at column
CREATE TRIGGER update_claims_updated_at
    BEFORE UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_claim_documents_updated_at
    BEFORE UPDATE ON claim_documents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_investigations_updated_at
    BEFORE UPDATE ON investigations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_appeals_updated_at
    BEFORE UPDATE ON appeals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_aml_alerts_updated_at
    BEFORE UPDATE ON aml_alerts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_claim_payments_updated_at
    BEFORE UPDATE ON claim_payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ombudsman_complaints_updated_at
    BEFORE UPDATE ON ombudsman_complaints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_policy_bond_tracking_updated_at
    BEFORE UPDATE ON policy_bond_tracking
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_freelook_cancellations_updated_at
    BEFORE UPDATE ON freelook_cancellations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for full-text search vector update
CREATE TRIGGER update_claims_search_vector
    BEFORE INSERT OR UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION update_claim_search_vector();

-- BR-CLM-DC-001: Trigger for investigation requirement check
CREATE TRIGGER check_investigation_requirement_trigger
    BEFORE INSERT OR UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION check_investigation_requirement();

-- Trigger for claim status change logging
CREATE TRIGGER log_claim_status_change_trigger
    AFTER UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION log_claim_status_change();

-- Trigger for workflow state validation
CREATE TRIGGER validate_workflow_transition_trigger
    BEFORE UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION validate_workflow_transition();

-- ============================================
-- VIEWS
-- ============================================

-- v_active_claims: Active claims with SLA tracking
CREATE OR REPLACE VIEW v_active_claims AS
SELECT
    c.id,
    c.claim_number,
    c.claim_type,
    c.policy_id,
    c.customer_id,
    c.claimant_name,
    c.status,
    c.claim_amount,
    c.approved_amount,
    c.investigation_required,
    c.sla_due_date,
    c.sla_status,
    c.sla_breached,
    EXTRACT(DAY FROM (c.sla_due_date - NOW()))::INTEGER as days_until_sla_breach,
    EXTRACT(DAY FROM (NOW() - c.created_at))::INTEGER as claim_age_days,
    c.created_at,
    c.updated_at
FROM claims c
WHERE c.status NOT IN ('CLOSED', 'PAID')
AND c.deleted_at IS NULL
ORDER BY c.sla_due_date ASC;

COMMENT ON VIEW v_active_claims IS 'Active claims with SLA countdown and aging';

-- v_investigation_queue: Claims pending investigation assignment
CREATE OR REPLACE VIEW v_investigation_queue AS
SELECT
    c.id as claim_id,
    c.claim_number,
    c.policy_id,
    c.customer_id,
    c.claimant_name,
    c.death_date,
    c.death_type,
    c.claim_amount,
    c.created_at as claim_registration_date,
    c.sla_due_date,
    c.sla_status,
    EXTRACT(DAY FROM (NOW() - c.created_at))::INTEGER as pending_days,
    CASE
        WHEN c.sla_status = 'RED' THEN 1
        WHEN c.sla_status = 'ORANGE' THEN 2
        WHEN c.sla_status = 'YELLOW' THEN 3
        ELSE 4
    END as priority_order
FROM claims c
WHERE c.investigation_required = TRUE
AND c.status = 'INVESTIGATION_PENDING'
AND c.investigator_id IS NULL
AND c.deleted_at IS NULL
ORDER BY priority_order ASC, c.created_at ASC;

COMMENT ON VIEW v_investigation_queue IS 'Claims pending investigation officer assignment';

-- v_approval_queue: Claims pending approval
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
    CASE
        WHEN c.sla_status = 'RED' THEN 1
        WHEN c.sla_status = 'ORANGE' THEN 2
        WHEN c.sla_status = 'YELLOW' THEN 3
        ELSE 4
    END as priority_order,
    c.created_at
FROM claims c
WHERE c.status = 'APPROVAL_PENDING'
AND c.deleted_at IS NULL
ORDER BY priority_order ASC, c.created_at ASC;

COMMENT ON VIEW v_approval_queue IS 'Claims pending approval sorted by SLA priority';

-- v_sla_breach_report: SLA breach statistics
CREATE OR REPLACE VIEW v_sla_breach_report AS
SELECT
    c.claim_type,
    c.status,
    COUNT(*) as total_claims,
    COUNT(*) FILTER (WHERE c.sla_breached = TRUE) as breached_count,
    COUNT(*) FILTER (WHERE c.sla_status = 'RED') as red_sla_count,
    COUNT(*) FILTER (WHERE c.sla_status = 'ORANGE') as orange_sla_count,
    COUNT(*) FILTER (WHERE c.sla_status = 'YELLOW') as yellow_sla_count,
    COUNT(*) FILTER (WHERE c.sla_status = 'GREEN') as green_sla_count,
    ROUND(AVG(EXTRACT(DAY FROM (c.updated_at - c.created_at)))::NUMERIC, 2) as avg_processing_days,
    ROUND((COUNT(*) FILTER (WHERE c.sla_breached = TRUE)::NUMERIC / NULLIF(COUNT(*), 0)) * 100, 2) as breach_percentage
FROM claims c
WHERE c.deleted_at IS NULL
GROUP BY c.claim_type, c.status
ORDER BY c.claim_type, c.status;

COMMENT ON VIEW v_sla_breach_report IS 'SLA breach analytics by claim type and status';

-- v_payment_queue: Approved claims pending payment
CREATE OR REPLACE VIEW v_payment_queue AS
SELECT
    c.id as claim_id,
    c.claim_number,
    c.claim_type,
    c.policy_id,
    c.claimant_name,
    c.approved_amount,
    c.bank_account_number,
    c.bank_ifsc_code,
    c.bank_account_holder_name,
    c.bank_verified,
    c.approval_date,
    EXTRACT(DAY FROM (NOW() - c.approval_date))::INTEGER as days_since_approval,
    c.payment_mode
FROM claims c
WHERE c.status = 'DISBURSEMENT_PENDING'
AND c.deleted_at IS NULL
ORDER BY c.approval_date ASC;

COMMENT ON VIEW v_payment_queue IS 'Claims approved and pending payment disbursement';

-- v_aml_high_risk_alerts: High risk AML alerts pending review
CREATE OR REPLACE VIEW v_aml_high_risk_alerts AS
SELECT
    a.id,
    a.alert_id,
    a.trigger_code,
    a.policy_id,
    a.customer_id,
    a.transaction_type,
    a.transaction_amount,
    a.risk_level,
    a.alert_status,
    a.transaction_blocked,
    a.filing_required,
    a.filing_type,
    EXTRACT(DAY FROM (NOW() - a.created_at))::INTEGER as alert_age_days,
    a.created_at
FROM aml_alerts a
WHERE a.risk_level IN ('HIGH', 'CRITICAL')
AND a.alert_status IN ('FLAGGED', 'UNDER_REVIEW')
ORDER BY
    CASE a.risk_level
        WHEN 'CRITICAL' THEN 1
        WHEN 'HIGH' THEN 2
        ELSE 3
    END,
    a.created_at ASC;

COMMENT ON VIEW v_aml_high_risk_alerts IS 'High and critical risk AML alerts requiring immediate attention';

-- v_ombudsman_compliance_pending: Ombudsman awards pending compliance
CREATE OR REPLACE VIEW v_ombudsman_compliance_pending AS
SELECT
    o.id,
    o.complaint_number,
    o.claim_id,
    o.policy_id,
    o.complainant_name,
    o.award_number,
    o.award_amount,
    o.award_date,
    o.compliance_due_date,
    EXTRACT(DAY FROM (o.compliance_due_date - NOW()))::INTEGER as days_until_due,
    o.compliance_status,
    CASE
        WHEN o.compliance_due_date < NOW() THEN 'OVERDUE'
        WHEN EXTRACT(DAY FROM (o.compliance_due_date - NOW())) <= 7 THEN 'URGENT'
        WHEN EXTRACT(DAY FROM (o.compliance_due_date - NOW())) <= 15 THEN 'WARNING'
        ELSE 'ON_TRACK'
    END as urgency_status
FROM ombudsman_complaints o
WHERE o.award_issued = TRUE
AND o.compliance_status != 'COMPLETED'
AND o.compliance_due_date IS NOT NULL
ORDER BY o.compliance_due_date ASC;

COMMENT ON VIEW v_ombudsman_compliance_pending IS 'BR-CLM-OMB-006: Ombudsman awards pending compliance within 30 days';

-- ============================================
-- ROW-LEVEL SECURITY
-- ============================================

-- Enable RLS on claims table
ALTER TABLE claims ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see their own assigned claims
CREATE POLICY claim_user_policy ON claims
    FOR SELECT
    USING (
        created_by = current_setting('app.current_user_id')::UUID OR
        updated_by = current_setting('app.current_user_id')::UUID OR
        approver_id = current_setting('app.current_user_id')::UUID OR
        investigator_id = current_setting('app.current_user_id')::UUID
    );

-- Policy: Approvers can see claims in approval queue
CREATE POLICY claim_approver_policy ON claims
    FOR SELECT
    USING (
        status = 'APPROVAL_PENDING' AND
        current_setting('app.user_role', true) IN ('APPROVER', 'SUPERVISOR', 'ADMIN')
    );

-- Policy: Admins can see all claims
CREATE POLICY claim_admin_policy ON claims
    FOR ALL
    USING (current_setting('app.user_role', true) = 'ADMIN');

-- Enable RLS on other sensitive tables
ALTER TABLE claim_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE aml_alerts ENABLE ROW LEVEL SECURITY;
ALTER TABLE claim_history ENABLE ROW LEVEL SECURITY;

-- Payment policies
CREATE POLICY payment_admin_policy ON claim_payments
    FOR ALL
    USING (current_setting('app.user_role', true) IN ('ADMIN', 'FINANCE_OFFICER'));

-- AML policies
CREATE POLICY aml_compliance_policy ON aml_alerts
    FOR ALL
    USING (current_setting('app.user_role', true) IN ('ADMIN', 'COMPLIANCE_OFFICER', 'AML_OFFICER'));

-- Audit trail policy (read-only for most users)
CREATE POLICY audit_read_policy ON claim_history
    FOR SELECT
    USING (current_setting('app.user_role', true) IN ('ADMIN', 'AUDITOR', 'COMPLIANCE_OFFICER'));

-- ============================================
-- INITIAL DATA / SEED DATA
-- ============================================

-- Insert document checklist templates for death claims
INSERT INTO document_checklist_templates (claim_type, death_type, nomination_status, document_type, document_description, is_mandatory, display_order) VALUES
-- BR-CLM-DC-015: Base mandatory documents
('DEATH', NULL, NULL, 'DEATH_CERTIFICATE', 'Death certificate issued by competent authority', TRUE, 1),
('DEATH', NULL, NULL, 'CLAIM_FORM', 'Duly filled claim form', TRUE, 2),
('DEATH', NULL, NULL, 'POLICY_BOND_OR_INDEMNITY', 'Original policy bond or indemnity bond', TRUE, 3),
('DEATH', NULL, NULL, 'CLAIMANT_ID_PROOF', 'ID proof of claimant (Aadhaar/PAN/Passport)', TRUE, 4),
('DEATH', NULL, NULL, 'BANK_MANDATE', 'Bank account details and cancelled cheque', TRUE, 5),

-- BR-CLM-DC-013: Unnatural death documents
('DEATH', 'UNNATURAL', NULL, 'FIR', 'First Information Report from police', TRUE, 6),
('DEATH', 'UNNATURAL', NULL, 'POSTMORTEM_REPORT', 'Post-mortem report', TRUE, 7),
('DEATH', 'ACCIDENTAL', NULL, 'FIR', 'First Information Report from police', TRUE, 6),
('DEATH', 'ACCIDENTAL', NULL, 'POSTMORTEM_REPORT', 'Post-mortem report', TRUE, 7),

-- BR-CLM-DC-014: Nomination absence documents
('DEATH', NULL, 'ABSENT', 'SUCCESSION_CERTIFICATE', 'Succession certificate from court', TRUE, 8),
('DEATH', NULL, 'ABSENT', 'LEGAL_HEIR_AFFIDAVIT', 'Legal heir affidavits', TRUE, 9);

-- ============================================
-- GRANTS AND PERMISSIONS
-- ============================================

-- Create roles
-- CREATE ROLE claims_service_user;
-- CREATE ROLE claims_admin;
-- CREATE ROLE claims_approver;
-- CREATE ROLE claims_investigator;
-- CREATE ROLE compliance_officer;
-- CREATE ROLE auditor;

-- Grant permissions
-- GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO claims_service_user;
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO claims_admin;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO auditor;
-- GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO claims_service_user;

-- ============================================
-- MAINTENANCE PROCEDURES
-- ============================================

-- Function to archive old claims (older than 7 years)
CREATE OR REPLACE FUNCTION archive_old_claims()
RETURNS INTEGER AS $$
DECLARE
    archived_count INTEGER;
BEGIN
    -- Move to archive table or update deleted_at
    UPDATE claims
    SET deleted_at = NOW()
    WHERE closure_date < (NOW() - INTERVAL '7 years')
    AND deleted_at IS NULL;

    GET DIAGNOSTICS archived_count = ROW_COUNT;
    RETURN archived_count;
END;
$$ LANGUAGE 'plpgsql';

-- Function to cleanup old audit logs (older than 10 years)
CREATE OR REPLACE FUNCTION cleanup_old_audit_logs()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM claim_history
    WHERE performed_at < (NOW() - INTERVAL '10 years');

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE 'plpgsql';

-- ============================================
-- PERFORMANCE OPTIMIZATION
-- ============================================

-- Analyze tables for query optimization
ANALYZE claims;
ANALYZE claim_documents;
ANALYZE investigations;
ANALYZE aml_alerts;
ANALYZE claim_payments;

-- ============================================
-- SCHEMA VERSION TRACKING
-- ============================================

CREATE TABLE schema_versions (
    version VARCHAR(20) PRIMARY KEY,
    description TEXT,
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    applied_by VARCHAR(100)
);

INSERT INTO schema_versions (version, description, applied_by)
VALUES ('1.0.0', 'Initial claims database schema with all entities, indexes, triggers, and views', 'system');

-- ============================================
-- END OF SCHEMA
-- ============================================
