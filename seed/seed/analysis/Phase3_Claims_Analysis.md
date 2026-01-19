# Phase 3: Claims Processing Analysis - Team 3

## Document Control

| Attribute | Details |
|-----------|---------|
| **Phase** | Phase 3 - Claims Processing |
| **Team** | Team 3 - Claim Management |
| **Analysis Date** | January 6, 2026 |
| **Last Updated** | January 12, 2026 (Functional Requirements Completed: 53 FRs including Ombudsman & Policy Bond) |
| **Documents Analyzed** | 7 SRS Documents + Missing Business Rules Report |
| **Total Pages** | ~280 pages |
| **Complexity** | High |
| **Technology Stack** | Golang, Temporal.io, PostgreSQL, Kafka |

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Business Rules](#2-business-rules)
3. [Functional Requirements](#3-functional-requirements)
4. [Validation Rules](#4-validation-rules)
5. [Error Codes](#5-error-codes)
6. [Workflows](#6-workflows)
7. [Data Entities](#7-data-entities)
8. [Integration Points](#8-integration-points)
9. [Temporal Workflows](#9-temporal-workflows)
10. [Traceability Matrix](#10-traceability-matrix)

---

## 1. Executive Summary

### 1.1 Purpose
This document provides comprehensive business requirements analysis for the Claims Processing module of the Postal Life Insurance (PLI) and Rural Postal Life Insurance (RPLI) system.

### 1.2 Scope
The analysis covers seven critical claim-related SRS documents:
1. Death Claim Settlement (SRS/FRS)
2. Maturity Claim Processing
3. Survival Benefit Claim Processing
4. AML Triggers & Alerts
5. Alerts/Triggers to Finnet/Fingate
6. Insurance Ombudsman System
7. Policy Bond Tracking & Free Look Cancellation

### 1.3 Key Statistics

| Metric | Count |
|--------|-------|
| **Total Business Rules** | 70 (24 original + 47 gap closure) |
| **Functional Requirements** | 53 (Death: 8, Maturity: 16, Survival Benefit: 15, AML: 5, Ombudsman: 4, Policy Bond & Freelook: 5) |
| **Validation Rules** | 123 (100% coverage achieved) |
| **Error Codes** | 50+ |
| **Workflows** | 12 major workflows |
| **Data Entities** | 40+ |
| **Integration Points** | 15+ |
| **Temporal Workflows** | 8 workflows |

### 1.4 Critical Dependencies
- Policy Services (Team 9) - Policy status validation
- Accounting (Team 2 & 6) - Financial entries and disbursements
- KYC/BCP (Team 7) - Identity verification
- Portal/Billing (Team 8) - Self-service claim submission
- Agent Management (Team 1) - Commission recovery

---

## 2. Business Rules

### 2.1 Death Claim Rules

#### BR-CLM-DC-001: Investigation Trigger
- **ID**: BR-CLM-DC-001
- **Category**: Death Claim Investigation
- **Priority**: CRITICAL
- **Description**: Mandatory investigation if death occurs within 3 years of policy acceptance or revival
- **Rule**: `IF (death_date - policy_issue_date) <= 3 years OR (death_date - revival_date) <= 3 years THEN trigger_investigation = TRUE`
- **Source**: Claim_SRS FRS on death claim.md, Section 3 - Investigation & Verification
- **Traceability**: FR-CLM-DC-003
- **Impact**: Claims settlement delayed until investigation complete

#### BR-CLM-DC-002: Investigation Timeline
- **ID**: BR-CLM-DC-002
- **Category**: Death Claim Investigation
- **Priority**: CRITICAL
- **Description**: Investigation report must be submitted within 21 days
- **Rule**: `investigation_report_due_date = investigation_start_date + 21 days`
- **Source**: Claim_SRS FRS on death claim.md, Section 3
- **Traceability**: WF-CLM-DC-001, Step 3
- **Impact**: SLA monitoring and escalation triggers

#### BR-CLM-DC-003: Approval Timeline Without Investigation
- **ID**: BR-CLM-DC-003
- **Category**: Death Claim Approval
- **Priority**: CRITICAL
- **Description**: Claims without investigation must be approved within 15 days
- **Rule**: `approval_due_date = claim_registration_date + 15 days WHERE investigation_required = FALSE`
- **Source**: Claim_SRS FRS on death claim.md, Section 5 - Approval Workflow
- **Traceability**: FR-CLM-DC-005, SLA-CLM-DC-001
- **Impact**: Auto-escalation on SLA breach

#### BR-CLM-DC-004: Approval Timeline With Investigation
- **ID**: BR-CLM-DC-004
- **Category**: Death Claim Approval
- **Priority**: CRITICAL
- **Description**: Claims requiring investigation must be approved within 45 days
- **Rule**: `approval_due_date = claim_registration_date + 45 days WHERE investigation_required = TRUE`
- **Source**: Claim_SRS FRS on death claim.md, Section 5
- **Traceability**: FR-CLM-DC-005, SLA-CLM-DC-002
- **Impact**: Extended SLA for complex cases

#### BR-CLM-DC-005: Appeal Window
- **ID**: BR-CLM-DC-005
- **Category**: Death Claim Appeal
- **Priority**: HIGH
- **Description**: Claimants may file appeal within 90 days of rejection
- **Rule**: `appeal_allowed = TRUE IF (current_date - rejection_date) <= 90 days`
- **Source**: Claim_SRS FRS on death claim.md, Section 9 - Appeal Mechanism
- **Traceability**: FR-CLM-DC-009
- **Impact**: Reopening of closed claims

#### BR-CLM-DC-006: Appellate Authority
- **ID**: BR-CLM-DC-006
- **Category**: Death Claim Appeal
- **Priority**: HIGH
- **Description**: Appeal authority is next higher officer in approval hierarchy
- **Rule**: `appellate_authority = get_next_higher_authority(original_approver_role)`
- **Source**: Claim_SRS FRS on death claim.md, Section 9
- **Traceability**: FR-CLM-DC-009
- **Impact**: Automated routing of appeals


#### BR-CLM-DC-007: Appeal Decision Timeline
- **ID**: BR-CLM-DC-007
- **Category**: Death Claim Appeal
- **Priority**: HIGH
- **Description**: Appellate authority must issue reasoned order within 45 days
- **Rule**: `appeal_decision_due_date = appeal_submission_date + 45 days`
- **Source**: Claim_SRS FRS on death claim.md, Section 9
- **Traceability**: SLA-CLM-DC-003
- **Impact**: Timely resolution of disputes

#### BR-CLM-DC-008: Claim Amount Calculation
- **ID**: BR-CLM-DC-008
- **Category**: Death Claim Settlement
- **Priority**: CRITICAL
- **Description**: Claim amount includes sum assured, bonuses, excess premiums minus deductions
- **Formula**: `claim_amount = sum_assured + accrued_bonuses + excess_premiums - (outstanding_loans + unpaid_premiums + applicable_taxes)`
- **Source**: Claim_SRS FRS on death claim.md, Section 4 - Calculation & Benefit Computation
- **Traceability**: FR-CLM-DC-004
- **Impact**: Accurate settlement computation

#### BR-CLM-DC-009: Penal Interest
- **ID**: BR-CLM-DC-009
- **Category**: Death Claim Compensation
- **Priority**: CRITICAL
- **Description**: Penal interest at 8% p.a. applicable post-SLA breach
- **Formula**: `penal_interest = (claim_amount * 0.08 * days_delayed) / 365 WHERE days_delayed = actual_settlement_date - sla_due_date`
- **Source**: Claim_SRS FRS on death claim.md, Section 7 - Gaps Identified (Gap closure recommendation, not current functional requirement)
- **Traceability**: FR-CLM-DC-010
- **Impact**: Auto-calculation required to prevent audit flags

#### BR-CLM-DC-010: Document Pending Period
- **ID**: BR-CLM-DC-010
- **Category**: Death Claim Documentation
- **Priority**: HIGH
- **Description**: Claims with missing documents must be returned after 15 days + 7-day grace period
- **Rule**: `IF document_status = PENDING AND (current_date - pending_start_date) > 22 days THEN return_claim = TRUE`
- **Source**: Claim_SRS FRS on death claim.md, Section 2 - Document Capture
- **Traceability**: FR-CLM-DC-002
- **Impact**: Prevents indefinite pending status

#### BR-CLM-DC-011: Investigation Officer Assignment Rules
- **ID**: BR-CLM-DC-011
- **Category**: Death Claim Investigation
- **Priority**: CRITICAL
- **Description**: Investigation officer must be IP, ASP, or PRI(P) rank. System must track officer assignment, jurisdiction, and prevent conflicts.
- **Rule/Formula**: `eligible_officers = GET_OFFICERS(roles IN ['IP', 'ASP', 'PRI(P)'], jurisdiction = claim_location)`
- **Source**: Claim_SRS FRS on death claim.md, Lines 77-79
- **Traceability**: FR-CLM-DC-003
- **Impact**: Ensures qualified investigation personnel and jurisdictional compliance

#### BR-CLM-DC-012: Investigation Status Classification
- **ID**: BR-CLM-DC-012
- **Category**: Death Claim Investigation
- **Priority**: CRITICAL
- **Description**: Based on investigation, claim status must be "Clear," "Suspect," or "Fraud." Suspect/Fraud escalated for manual review.
- **Rule/Formula**: `claim_investigation_status = CASE investigation_outcome WHEN 'NO_SUSPICION_FOUND' THEN 'CLEAR'`
- **Source**: Claim_SRS FRS on death claim.md, Lines 84-86
- **Traceability**: FR-CLM-DC-003
- **Impact**: Standardized investigation outcome classification for fraud detection

#### BR-CLM-DC-013: Unnatural Death Document Requirements
- **ID**: BR-CLM-DC-013
- **Category**: Death Claim Documentation
- **Priority**: CRITICAL
- **Description**: For unnatural deaths, FIR and postmortem report are mandatory. System must flag missing documents.
- **Rule/Formula**: `IF death_type IN ['UNNATURAL', 'ACCIDENTAL'] THEN mandatory_documents.ADD('FIR', 'POSTMORTEM_REPORT')`
- **Source**: Claim_SRS FRS on death claim.md, Lines 63-64
- **Traceability**: VR-CLM-DC-006
- **Impact**: Ensures compliance with legal requirements for unnatural death claims

#### BR-CLM-DC-014: Nomination Absence Document Requirements
- **ID**: BR-CLM-DC-014
- **Category**: Death Claim Documentation
- **Priority**: CRITICAL
- **Description**: If nomination absent, succession certificate OR legal heir affidavits required.
- **Rule/Formula**: `IF policy_nomination_status = 'ABSENT' THEN mandatory_documents.ADD_EITHER_OR('SUCCESSION_CERTIFICATE', 'LEGAL_HEIR_AFFIDAVIT')`
- **Source**: Claim_SRS FRS on death claim.md, Lines 66-67
- **Traceability**: VR-CLM-DC-007
- **Impact**: Legal compliance for rightful claimant verification

#### BR-CLM-DC-015: Mandatory Document Checklist
- **ID**: BR-CLM-DC-015
- **Category**: Death Claim Documentation
- **Priority**: HIGH
- **Description**: Base mandatory documents: Death certificate, claim form, policy bond/indemnity, ID/address proof, bank mandate.
- **Rule/Formula**: `mandatory_documents = ['DEATH_CERTIFICATE', 'CLAIM_FORM', 'POLICY_BOND_OR_INDEMNITY', 'ID_PROOF', 'BANK_MANDATE']`
- **Source**: Claim_SRS FRS on death claim.md, Lines 61-64
- **Traceability**: FR-CLM-DC-002, VR-CLM-DC-001 to VR-CLM-DC-005
- **Impact**: Standardized document checklist for all death claims

#### BR-CLM-DC-016: Manual Calculation Override with Audit Trail
- **ID**: BR-CLM-DC-016
- **Category**: Death Claim Settlement
- **Priority**: CRITICAL
- **Description**: Manual override allowed for disputed data/court-directed adjustments. ALL overrides logged with user ID, timestamp, before/after values.
- **Rule/Formula**: `LOG_OVERRIDE(user_id, timestamp, field_name, old_value, new_value, reason, approval_authority)`
- **Source**: Claim_SRS FRS on death claim.md, Lines 95-100
- **Traceability**: FR-CLM-DC-004
- **Impact**: Audit trail for compliance and fraud prevention

#### BR-CLM-DC-017: Payment Mode Selection Rules
- **ID**: BR-CLM-DC-017
- **Category**: Death Claim Payment
- **Priority**: HIGH
- **Description**: Payment priority: 1) NEFT 2) POSB EFT 3) Cheque (fallback). Integrate with Finacle/IT 2.0.
- **Rule/Formula**: `payment_mode = CASE WHEN neft_available THEN 'NEFT' WHEN posb_eft_available THEN 'POSB_EFT' ELSE 'CHEQUE' END`
- **Source**: Claim_SRS FRS on death claim.md, Lines 116-118
- **Traceability**: FR-CLM-DC-006, INT-CLM-002, INT-CLM-003
- **Impact**: Optimized payment delivery with fallback options

#### BR-CLM-DC-018: Claim Reopen Valid Circumstances
- **ID**: BR-CLM-DC-018
- **Category**: Death Claim Appeal
- **Priority**: HIGH
- **Description**: Claims may only be reopened for: Court orders, new evidence, administrative lapses, claimant appeals.
- **Rule/Formula**: `claim_reopen_allowed = IF reason IN ['COURT_ORDER', 'NEW_EVIDENCE', 'ADMIN_LAPSE', 'CLAIMANT_APPEAL']`
- **Source**: Claim_SRS FRS on death claim.md, Lines 124-128
- **Traceability**: FR-CLM-DC-009
- **Impact**: Controlled claim reopening process with valid justification

#### BR-CLM-DC-019: Communication Milestone Triggers
- **ID**: BR-CLM-DC-019
- **Category**: Death Claim Communication
- **Priority**: MEDIUM
- **Description**: Automated notifications at: Registration, document status, investigation, approval/rejection, payment. All communications logged.
- **Rule/Formula**: `SEND_NOTIFICATION(event_type, channel) WHERE event_type IN ['REGISTRATION', 'DOCUMENT_STATUS', 'INVESTIGATION', 'APPROVAL', 'PAYMENT']`
- **Source**: Claim_SRS FRS on death claim.md, Lines 134-138
- **Traceability**: FR-CLM-DC-007, INT-CLM-004, INT-CLM-005, INT-CLM-006
- **Impact**: Improved transparency and customer experience

#### BR-CLM-DC-020: Claim Case Owner Assignment
- **ID**: BR-CLM-DC-020
- **Category**: Death Claim Process Management
- **Priority**: HIGH
- **Description**: "Claim Case Owner" at CPC level tracks claim from registration to closure. Addresses fragmented accountability.
- **Rule/Formula**: `claim_case_owner = ASSIGN_OWNER(claim_id, cpc_location) ON claim_registration`
- **Source**: Claim_SRS FRS on death claim.md, Lines 158-166 (Gap section)
- **Traceability**: FR-CLM-DC-001
- **Impact**: Single point of accountability for claim lifecycle

#### BR-CLM-DC-021: SLA Color-Coded Alert System
- **ID**: BR-CLM-DC-021
- **Category**: Death Claim Process Management
- **Priority**: HIGH
- **Description**: Green (<70%), Yellow (70-90%), Orange (90-100%), Red (>100% SLA). Auto-escalation on breach.
- **Rule/Formula**: `sla_status = CASE WHEN sla_consumed < 0.70 THEN 'GREEN' WHEN sla_consumed < 0.90 THEN 'YELLOW' WHEN sla_consumed <= 1.00 THEN 'ORANGE' ELSE 'RED' END`
- **Source**: Claim_SRS FRS on death claim.md, Lines 318-328
- **Traceability**: FR-CLM-DC-005
- **Impact**: Visual SLA monitoring and proactive escalation

#### BR-CLM-DC-022: Disbursement Reconciliation with Banking
- **ID**: BR-CLM-DC-022
- **Category**: Death Claim Payment
- **Priority**: CRITICAL
- **Description**: Full integration with Finacle/IT 2.0 for real-time status, daily reconciliation, prevent double payment.
- **Rule/Formula**: `RECONCILE_PAYMENT(claim_id, transaction_id, bank_status) DAILY; CHECK_DUPLICATE_PAYMENT(policy_id) BEFORE disbursement`
- **Source**: Claim_SRS FRS on death claim.md, Lines 222-234 (Gap section)
- **Traceability**: FR-CLM-DC-006, INT-CLM-002
- **Impact**: Payment accuracy and fraud prevention

#### BR-CLM-DC-023: Rejection with Root Cause Analysis
- **ID**: BR-CLM-DC-023
- **Category**: Death Claim Quality Management
- **Priority**: MEDIUM
- **Description**: Mandate RCA for every rejection/delay. Monthly review at Divisional/Regional level.
- **Rule/Formula**: `REQUIRE_RCA(claim_id) WHERE claim_status = 'REJECTED' OR sla_breached = TRUE; GENERATE_RCA_REPORT() MONTHLY`
- **Source**: Claim_SRS FRS on death claim.md, Lines 249-257 (Gap section)
- **Traceability**: FR-CLM-DC-009
- **Impact**: Continuous improvement and process optimization

#### BR-CLM-DC-024: Real-Time Claimant Status Tracking
- **ID**: BR-CLM-DC-024
- **Category**: Death Claim Communication
- **Priority**: MEDIUM
- **Description**: Enable claim tracking via SMS/Email/Portal/Mobile app with Claim ID. Self-service status lookup.
- **Rule/Formula**: `PROVIDE_STATUS_API(claim_id) via [SMS, EMAIL, PORTAL, MOBILE_APP]`
- **Source**: Claim_SRS FRS on death claim.md, Lines 259-270 (Gap section)
- **Traceability**: FR-CLM-DC-007, INT-CLM-008
- **Impact**: Enhanced customer experience and reduced inquiry volume

#### BR-CLM-DC-025: Audit Trail for Manual Overrides
- **ID**: BR-CLM-DC-025
- **Category**: Death Claim Compliance
- **Priority**: CRITICAL
- **Description**: Enforce mandatory remarks and digital signature for ALL overrides. Logged in audit trail.
- **Rule/Formula**: `VALIDATE_OVERRIDE(digital_signature, remarks) BEFORE override_allowed = TRUE; LOG_AUDIT(user_id, action, timestamp, digital_signature)`
- **Source**: Claim_SRS FRS on death claim.md, Lines 272-281 (Gap section)
- **Traceability**: FR-CLM-DC-004
- **Impact**: Complete audit trail for regulatory compliance

### 2.2 Maturity Claim Rules

#### BR-CLM-MC-001: Report Generation Schedule
- **ID**: BR-CLM-MC-001
- **Category**: Maturity Claim Reporting
- **Priority**: HIGH
- **Description**: Maturity report must be generated daily/weekly for policies maturing in next 2-3 months
- **Rule**: `generate_maturity_report() DAILY_OR_WEEKLY FOR policies WHERE maturity_date BETWEEN (current_date + 1 month) AND (current_date + 3 months)`
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-01
- **Traceability**: FR-CLM-MC-001
- **Impact**: Proactive customer intimation

#### BR-CLM-MC-002: Approval SLA
- **ID**: BR-CLM-MC-002
- **Category**: Maturity Claim Approval
- **Priority**: CRITICAL
- **Description**: Maturity claims must be approved within 7 days
- **Rule**: `approval_due_date = claim_submission_date + 7 days`
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-08
- **Traceability**: SLA-CLM-MC-001
- **Impact**: Faster turnaround for maturity settlements

#### BR-CLM-MC-003: Bank Verification Requirement
- **ID**: BR-CLM-MC-003
- **Category**: Maturity Claim Payment
- **Priority**: CRITICAL
- **Description**: Bank account must be verified via CBS/PFMS API before disbursement
- **Rule**: `disbursement_allowed = TRUE ONLY IF bank_verification_status = VERIFIED via CBS_API OR PFMS_API`
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-10
- **Traceability**: FR-CLM-MC-006, INT-CLM-MC-002
- **Impact**: Prevents payment failures and fraud

#### BR-CLM-MC-004: Multi-Channel Intimation
- **ID**: BR-CLM-MC-004
- **Category**: Maturity Claim Communication
- **Priority**: MEDIUM
- **Description**: Maturity intimation must be sent via multiple channels
- **Rule**: `send_intimation(policyholder_id) via [SMS, EMAIL, WHATSAPP, PORTAL] WHERE maturity_date - current_date <= 60 days`
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-02
- **Traceability**: FR-CLM-MC-002
- **Impact**: Improved customer reach and awareness

#### BR-CLM-MC-005: Policy Activation Status Validation
- **ID**: BR-CLM-MC-005
- **Category**: Maturity Claim Validation
- **Priority**: CRITICAL
- **Description**: Policy must be active on maturity date. Reject if lapsed/forfeited/terminated.
- **Rule/Formula**: `IF policy_status != 'ACTIVE' ON maturity_date THEN reject_claim(code: 'RJ-P-02')`
- **Source**: Claim_SRS FRS on Maturity claim.md, Code RJ-P-02, Line 650
- **Traceability**: VR-CLM-MC-001
- **Impact**: Prevents invalid maturity claim processing

#### BR-CLM-MC-006: Duplicate Maturity Claim Prevention
- **ID**: BR-CLM-MC-006
- **Category**: Maturity Claim Validation
- **Priority**: CRITICAL
- **Description**: Check if maturity claim already paid. Prevent duplicate payments.
- **Rule/Formula**: `IF EXISTS(SELECT * FROM claims WHERE policy_id = X AND claim_type = 'MATURITY' AND status = 'PAID') THEN reject_claim(code: 'RJ-P-03')`
- **Source**: Claim_SRS FRS on Maturity claim.md, Code RJ-P-03, Lines 652-653
- **Traceability**: VR-CLM-MC-002
- **Impact**: Fraud prevention and financial accuracy

#### BR-CLM-MC-007: Policy Forfeiture/Surrender Pre-Maturity Check
- **ID**: BR-CLM-MC-007
- **Category**: Maturity Claim Validation
- **Priority**: CRITICAL
- **Description**: Reject if policy terminated due to forfeiture/surrender before maturity.
- **Rule/Formula**: `IF policy_termination_date < maturity_date AND termination_reason IN ['FORFEITURE', 'SURRENDER'] THEN reject_claim(code: 'RJ-P-04')`
- **Source**: Claim_SRS FRS on Maturity claim.md, Code RJ-P-04, Lines 654-655
- **Traceability**: VR-CLM-MC-010
- **Impact**: Prevents invalid claims on terminated policies

#### BR-CLM-MC-008: Claimant Identity Verification
- **ID**: BR-CLM-MC-008
- **Category**: Maturity Claim Validation
- **Priority**: CRITICAL
- **Description**: Verify claimant details match policy records and identity can be established.
- **Rule/Formula**: `VERIFY_IDENTITY(claimant_name, claimant_dob, policy_holder_name, policy_holder_dob) AND VALIDATE_ID_PROOF(id_type, id_number)`
- **Source**: Claim_SRS FRS on Maturity claim.md, Codes RJ-E-01, RJ-E-02, Lines 663-666
- **Traceability**: VR-CLM-MC-003
- **Impact**: Fraud prevention and rightful claimant verification

#### BR-CLM-MC-009: Forged/Suspicious Document Detection
- **ID**: BR-CLM-MC-009
- **Category**: Maturity Claim Document Validation
- **Priority**: HIGH
- **Description**: Screen documents for forgery/suspicious alterations. Reject claims with forged documents.
- **Rule/Formula**: `IF DETECT_FORGERY(document) = TRUE THEN reject_claim(code: 'RJ-D-02') AND flag_for_investigation = TRUE`
- **Source**: Claim_SRS FRS on Maturity claim.md, Code RJ-D-02, Line 683
- **Traceability**: VR-CLM-MC-008
- **Impact**: Fraud detection and prevention

#### BR-CLM-MC-010: Missing Document Auto-Reminder Schedule
- **ID**: BR-CLM-MC-010
- **Category**: Maturity Claim Document Management
- **Priority**: MEDIUM
- **Description**: Reminders at Day 3, 7, 12, 14 (final). Auto-return after 15 days.
- **Rule/Formula**: `SEND_REMINDER(claim_id) ON days [3, 7, 12, 14]; IF days_pending > 15 THEN return_claim = TRUE`
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-04, Lines 477-489
- **Traceability**: FR-CLM-MC-003
- **Impact**: Proactive document collection and timely claim processing

#### BR-CLM-MC-011: Bank Account Re-submission Workflow
- **ID**: BR-CLM-MC-011
- **Category**: Maturity Claim Payment
- **Priority**: HIGH
- **Description**: Allow max 3 re-submission attempts for invalid bank details. Escalate to manual review after 3 failures.
- **Rule/Formula**: `IF bank_verification_attempts >= 3 AND bank_verified = FALSE THEN escalate_to_manual_review = TRUE`
- **Source**: Claim_SRS FRS on Maturity claim.md, Codes RJ-B-01/02/03, Lines 698-703
- **Traceability**: VR-CLM-MC-004
- **Impact**: Prevents indefinite re-submission cycles

#### BR-CLM-MC-012: Automated Maturity Date Calculation
- **ID**: BR-CLM-MC-012
- **Category**: Maturity Claim Processing
- **Priority**: HIGH
- **Description**: Calculate maturity date from policy issue date + term. Allow claims 30 days before to 90 days after maturity.
- **Rule/Formula**: `maturity_date = policy_issue_date + policy_term_years; claim_window = [maturity_date - 30 days, maturity_date + 90 days]`
- **Source**: Claim_SRS FRS on Maturity claim.md, Lines 96-98, 440-447
- **Traceability**: FR-CLM-MC-001
- **Impact**: Accurate maturity determination and flexible claim window

### 2.3 Survival Benefit Rules

#### BR-CLM-SB-001: Survival Benefit Report Generation
- **ID**: BR-CLM-SB-001
- **Category**: Survival Benefit Reporting
- **Priority**: HIGH
- **Description**: SB report must be generated daily/weekly for benefits due in next 2 months
- **Rule**: `generate_sb_report() DAILY_OR_WEEKLY FOR policies WHERE sb_due_date BETWEEN (current_date + 1 month) AND (current_date + 2 months)`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-01
- **Traceability**: FR-CLM-SB-001
- **Impact**: Timely benefit processing

#### BR-CLM-SB-002: Approval SLA
- **ID**: BR-CLM-SB-002
- **Category**: Survival Benefit Approval
- **Priority**: CRITICAL
- **Description**: Survival benefit claims must be approved within 7 days
- **Rule**: `approval_due_date = claim_submission_date + 7 days`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-08
- **Traceability**: SLA-CLM-SB-001
- **Impact**: SLA enforcement with countdown

#### BR-CLM-SB-003: Survival Benefit Due Date Calculation
- **ID**: BR-CLM-SB-003
- **Category**: Survival Benefit Processing
- **Priority**: HIGH
- **Description**: Calculate SB due dates based on policy terms. Generate reports for benefits due in next 2 months.
- **Rule/Formula**: `sb_due_dates = CALCULATE_SB_DATES(policy_issue_date, policy_term, sb_intervals); GENERATE_REPORT() WHERE sb_due_date BETWEEN current_date AND (current_date + 2 months)`
- **Source**: Claim_SRS FRS on survival benefit.md, Lines 97-98
- **Traceability**: FR-CLM-SB-001
- **Impact**: Proactive benefit processing and customer notification

#### BR-CLM-SB-004: Survival Benefit Eligibility Validation
- **ID**: BR-CLM-SB-004
- **Category**: Survival Benefit Validation
- **Priority**: CRITICAL
- **Description**: Validate: 1) Policy exists 2) Policy active on SB due date 3) SB not already paid.
- **Rule/Formula**: `IF policy_exists = FALSE THEN reject(code: 'RJ-P-01'); IF policy_status != 'ACTIVE' ON sb_due_date THEN reject(code: 'RJ-P-02'); IF sb_already_paid = TRUE THEN reject(code: 'RJ-P-03')`
- **Source**: Claim_SRS FRS on survival benefit.md, Codes RJ-P-01/02/03, Lines 445-451
- **Traceability**: VR-CLM-SB-001, VR-CLM-SB-002
- **Impact**: Prevents invalid and duplicate SB payments

#### BR-CLM-SB-005: DigiLocker Integration
- **ID**: BR-CLM-SB-005
- **Category**: Survival Benefit Document Management
- **Priority**: MEDIUM
- **Description**: Integrate with DigiLocker for policy document fetching. Require consent. Fall back to manual upload.
- **Rule/Formula**: `IF digilocker_consent = TRUE THEN FETCH_DOCUMENT(digilocker_uri) ELSE require_manual_upload = TRUE`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-03, Lines 341-348
- **Traceability**: FR-CLM-SB-002, INT-CLM-007
- **Impact**: Simplified document submission and faster processing

#### BR-CLM-SB-006: Auto-Acknowledgment with Claim ID
- **ID**: BR-CLM-SB-006
- **Category**: Survival Benefit Processing
- **Priority**: HIGH
- **Description**: Generate unique Claim ID (SB-YYYY-MM-DD-XXXXXX). Send acknowledgement via SMS/Email/Portal.
- **Rule/Formula**: `claim_id = GENERATE_CLAIM_ID(format: 'SB-YYYY-MM-DD-XXXXXX'); SEND_ACKNOWLEDGMENT(claim_id) via [SMS, EMAIL, PORTAL]`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-03, Lines 346-348
- **Traceability**: FR-CLM-SB-001, INT-CLM-004, INT-CLM-005
- **Impact**: Immediate customer confirmation and tracking capability

#### BR-CLM-SB-007: Automatic Indexing as Service Request
- **ID**: BR-CLM-SB-007
- **Category**: Survival Benefit Processing
- **Priority**: MEDIUM
- **Description**: Auto-index claim as Service Request. Link to policy and Claim ID.
- **Rule/Formula**: `CREATE_SERVICE_REQUEST(claim_id, policy_id, request_type: 'SURVIVAL_BENEFIT'); LINK_TO_POLICY(service_request_id, policy_id)`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-05, Lines 358-362
- **Traceability**: FR-CLM-SB-001
- **Impact**: Integrated service request management

#### BR-CLM-SB-008: OCR Auto-Population
- **ID**: BR-CLM-SB-008
- **Category**: Survival Benefit Document Processing
- **Priority**: MEDIUM
- **Description**: Extract data from uploaded documents using OCR. CPC supervisor verifies auto-populated data.
- **Rule/Formula**: `extracted_data = OCR_EXTRACT(document); auto_populate_fields(extracted_data); REQUIRE_VERIFICATION(cpc_supervisor) BEFORE processing`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-07, Lines 371-377
- **Traceability**: FR-CLM-SB-002
- **Impact**: Reduced manual data entry and improved accuracy

#### BR-CLM-SB-009: Digital Signature for Approval
- **ID**: BR-CLM-SB-009
- **Category**: Survival Benefit Approval
- **Priority**: HIGH
- **Description**: Approver must use digital signature for approval/rejection. Signature linked to user ID and timestamp.
- **Rule/Formula**: `REQUIRE_DIGITAL_SIGNATURE(approver_id) BEFORE approval_allowed = TRUE; LOG_SIGNATURE(user_id, timestamp, action, digital_signature_hash)`
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-08, Lines 379-386
- **Traceability**: FR-CLM-SB-003
- **Impact**: Non-repudiation and audit compliance

### 2.4 AML/CFT Rules

#### BR-CLM-AML-001: High Cash Premium Alert
- **ID**: BR-CLM-AML-001
- **Category**: AML Trigger
- **Priority**: CRITICAL
- **Description**: Cash transactions over ₹50,000 trigger high-risk alert and CTR filing
- **Rule**: `IF payment_mode = CASH AND amount > 50000 THEN risk_level = HIGH AND trigger_CTR_filing = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Trigger AML_001
- **Traceability**: FR-CLM-AML-001
- **Impact**: Regulatory compliance (PMLA 2002)

#### BR-CLM-AML-002: PAN Mismatch Alert
- **ID**: BR-CLM-AML-002
- **Category**: AML Trigger
- **Priority**: HIGH
- **Description**: PAN verification failure triggers medium-risk alert for manual review
- **Rule**: `IF pan_verified = FALSE THEN risk_level = MEDIUM AND flag_for_review = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Trigger AML_002
- **Traceability**: FR-CLM-AML-002
- **Impact**: Identity verification compliance

#### BR-CLM-AML-003: Nominee Change Post Death (Critical)
- **ID**: BR-CLM-AML-003
- **Category**: AML Trigger
- **Priority**: CRITICAL
- **Description**: Nominee change after death date triggers critical alert, blocks transaction, and files STR
- **Rule**: `IF nominee_change_date > death_date THEN risk_level = CRITICAL AND block_transaction = TRUE AND trigger_STR_filing = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Trigger AML_003
- **Traceability**: FR-CLM-AML-003
- **Impact**: Fraud prevention

#### BR-CLM-AML-004: Frequent Surrenders
- **ID**: BR-CLM-AML-004
- **Category**: AML Trigger
- **Priority**: MEDIUM
- **Description**: More than 3 surrenders within 6 months by single customer triggers investigation
- **Rule**: `IF count(surrenders WHERE customer_id = X AND surrender_date BETWEEN (current_date - 6 months) AND current_date) > 3 THEN risk_level = MEDIUM AND flag_for_investigation = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Trigger AML_004
- **Traceability**: FR-CLM-AML-004
- **Impact**: Money laundering detection

#### BR-CLM-AML-005: Refund Without Bond Delivery
- **ID**: BR-CLM-AML-005
- **Category**: AML Trigger
- **Priority**: HIGH
- **Description**: Refund issued before bond dispatch triggers high-risk alert
- **Rule**: `IF refund_date < bond_dispatch_date THEN risk_level = HIGH AND log_audit_trail = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Trigger AML_005
- **Traceability**: FR-CLM-AML-005
- **Impact**: Process anomaly detection

#### BR-CLM-AML-006: STR Filing Timeline
- **ID**: BR-CLM-AML-006
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Suspicious Transaction Reports must be filed within 7 working days
- **Rule**: `STR_filing_due_date = suspicion_determination_date + 7 working_days`
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Section 3.3
- **Traceability**: FR-CLM-AML-010
- **Impact**: PMLA Section 12 compliance

#### BR-CLM-AML-007: CTR Filing Schedule
- **ID**: BR-CLM-AML-007
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Cash Transaction Reports must be filed monthly for aggregates over ₹10 lakh in one day
- **Rule**: `file_CTR() MONTHLY FOR transactions WHERE payment_mode = CASH AND daily_aggregate > 1000000`
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Section 3.3
- **Traceability**: FR-CLM-AML-010
- **Impact**: PMLA compliance for cash transaction monitoring

#### BR-CLM-AML-008: CTR Aggregate Monitoring
- **ID**: BR-CLM-AML-008
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Track daily cash aggregates by customer. If > ₹10 lakh in single day, trigger CTR filing monthly.
- **Rule/Formula**: `daily_cash_aggregate = SUM(cash_transactions WHERE customer_id = X AND transaction_date = current_date); IF daily_cash_aggregate > 1000000 THEN trigger_CTR_filing = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 254-257; Finnet SRS
- **Traceability**: FR-CLM-AML-001
- **Impact**: PMLA compliance for cash transaction monitoring

#### BR-CLM-AML-009: Third-Party PAN Verification
- **ID**: BR-CLM-AML-009
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Mandatory PAN & KYC for third-party payments. Verify PAN via NSDL. Block if verification fails.
- **Rule/Formula**: `IF payment_recipient != policy_holder THEN REQUIRE_PAN_KYC(recipient); IF VERIFY_PAN(pan_number) = FALSE THEN block_transaction = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 4.4C, Lines 183-190
- **Traceability**: VR-CLM-AML-002, INT-CLM-010
- **Impact**: Third-party payment fraud prevention

#### BR-CLM-AML-010: Regulatory Reporting to FIU-IND
- **ID**: BR-CLM-AML-010
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Submit STR (7 days), CTR (monthly), CCR (immediate), NTR (as per guidelines) to FIU-IND.
- **Rule/Formula**: `SUBMIT_STR() within 7 working_days; SUBMIT_CTR() MONTHLY; SUBMIT_CCR() IMMEDIATE; SUBMIT_NTR() as_per_guidelines`
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Lines 11-18; Section 4.4D
- **Traceability**: FR-CLM-AML-010, INT-CLM-009
- **Impact**: Complete regulatory compliance with FIU-IND

#### BR-CLM-AML-011: Negative List Daily Screening
- **ID**: BR-CLM-AML-011
- **Category**: AML Compliance
- **Priority**: CRITICAL
- **Description**: Daily screening against OFAC, UN Sanctions, UAPA Section 51A, FATF lists. Freeze accounts on match.
- **Rule/Formula**: `SCREEN_DAILY() against [OFAC_LIST, UN_SANCTIONS, UAPA_51A, FATF_LIST]; IF match_found = TRUE THEN freeze_account = TRUE AND block_transactions = TRUE`
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 4.1, Lines 126-129
- **Traceability**: FR-CLM-AML-001
- **Impact**: Sanctions compliance and terrorist financing prevention

#### BR-CLM-AML-012: Beneficial Ownership Verification
- **ID**: BR-CLM-AML-012
- **Category**: AML Compliance
- **Priority**: HIGH
- **Description**: For non-individual customers (companies, trusts, NGOs), verify beneficial ownership. Screen owners against negative lists.
- **Rule/Formula**: `IF customer_type IN ['COMPANY', 'TRUST', 'NGO'] THEN VERIFY_BENEFICIAL_OWNERS(customer_id); SCREEN_OWNERS() against negative_lists`
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 4.1, Lines 131-133
- **Traceability**: FR-CLM-AML-001
- **Impact**: Enhanced due diligence for corporate entities

### 2.5 Insurance Ombudsman Rules

#### BR-CLM-OMB-001: Complaint Admissibility (Rule 14)
- **ID**: BR-CLM-OMB-001
- **Category**: Insurance Ombudsman
- **Priority**: HIGH
- **Description**: Check: 1) Representation to insurer first 2) 30-day wait 3) Limitation period (1 year) 4) Claim value ≤ ₹50 lakh 5) No parallel litigation.
- **Rule/Formula**: `admissible = (representation_to_insurer = TRUE) AND (wait_period >= 30 days) AND (complaint_age <= 1 year) AND (claim_value <= 5000000) AND (parallel_litigation = FALSE)`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 116-119
- **Traceability**: FR-CLM-OMB-001
- **Impact**: Validates complaint eligibility per IRDAI regulations

#### BR-CLM-OMB-002: Jurisdiction Mapping (Rule 11)
- **ID**: BR-CLM-OMB-002
- **Category**: Insurance Ombudsman
- **Priority**: HIGH
- **Description**: Map complaint to ombudsman center based on complainant pincode, agent location, or policy office.
- **Rule/Formula**: `ombudsman_center = MAP_JURISDICTION(complainant_pincode) OR MAP_JURISDICTION(agent_location) OR MAP_JURISDICTION(policy_office_location)`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 119, 200-206
- **Traceability**: FR-CLM-OMB-002
- **Impact**: Correct jurisdictional routing of complaints

#### BR-CLM-OMB-003: Conflict of Interest Screening
- **ID**: BR-CLM-OMB-003
- **Category**: Insurance Ombudsman
- **Priority**: HIGH
- **Description**: Screen for: 1) Prior relationship with complainant 2) Financial interest 3) Duplicate litigation. Reassign if conflict found.
- **Rule/Formula**: `conflict_exists = CHECK_CONFLICT(ombudsman_id, complainant_id, policy_id); IF conflict_exists = TRUE THEN reassign_to_alternate_ombudsman = TRUE`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 124-127
- **Traceability**: FR-CLM-OMB-002
- **Impact**: Ensures impartial complaint handling

#### BR-CLM-OMB-004: Mediation Recommendation (Rule 16)
- **ID**: BR-CLM-OMB-004
- **Category**: Insurance Ombudsman
- **Priority**: MEDIUM
- **Description**: If settlement reached via mediation, issue recommendation. Require acceptance from both parties within 15 days.
- **Rule/Formula**: `IF mediation_successful = TRUE THEN ISSUE_RECOMMENDATION(mediation_terms); REQUIRE_ACCEPTANCE(both_parties) within 15 days`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 136-138
- **Traceability**: FR-CLM-OMB-003
- **Impact**: Facilitates alternative dispute resolution

#### BR-CLM-OMB-005: Award Issuance with Caps (Rule 17)
- **ID**: BR-CLM-OMB-005
- **Category**: Insurance Ombudsman
- **Priority**: CRITICAL
- **Description**: Draft, review, approve, digitally sign award. Enforce ₹50 lakh cap. Award binding on insurer.
- **Rule/Formula**: `IF award_amount > 5000000 THEN reject_award(reason: 'EXCEEDS_CAP'); REQUIRE_DIGITAL_SIGNATURE(ombudsman_id) BEFORE award_issuance; award_binding_on_insurer = TRUE`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 78-79
- **Traceability**: FR-CLM-OMB-004
- **Impact**: Regulatory compliance and enforceability

#### BR-CLM-OMB-006: Insurer Compliance Monitoring
- **ID**: BR-CLM-OMB-006
- **Category**: Insurance Ombudsman
- **Priority**: CRITICAL
- **Description**: Track compliance within 30 days. Send reminders at Day 15, 7, 2. Escalate to IRDAI on breach.
- **Rule/Formula**: `compliance_due_date = award_date + 30 days; SEND_REMINDER() on days [15, 7, 2] before due_date; IF compliance_date > compliance_due_date THEN ESCALATE_TO_IRDAI()`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 142-147
- **Traceability**: FR-CLM-OMB-005
- **Impact**: Ensures timely award implementation

#### BR-CLM-OMB-007: Complaint Closure & Archival
- **ID**: BR-CLM-OMB-007
- **Category**: Insurance Ombudsman
- **Priority**: MEDIUM
- **Description**: Archive all documents with retention period (10 years for awards, 7 years for mediation). Preserve audit logs.
- **Rule/Formula**: `ARCHIVE_DOCUMENTS(complaint_id, retention_period); retention_period = IF award_issued THEN 10 years ELSE 7 years; PRESERVE_AUDIT_LOGS(complaint_id)`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 146-148
- **Traceability**: FR-CLM-OMB-006
- **Impact**: Regulatory compliance and record keeping

#### BR-CLM-OMB-008: Bilingual Support
- **ID**: BR-CLM-OMB-008
- **Category**: Insurance Ombudsman
- **Priority**: MEDIUM
- **Description**: All user-facing fields and documents available in English and Hindi. Adaptable for regional languages.
- **Rule/Formula**: `PROVIDE_TRANSLATION(content, languages: ['ENGLISH', 'HINDI']); SUPPORT_REGIONAL_LANGUAGES(adaptable: TRUE)`
- **Source**: Claim_SRS on insurance ombudsman.md, Lines 99-101
- **Traceability**: FR-CLM-OMB-007
- **Impact**: Accessibility and inclusivity for complainants

### 2.6 Policy Bond & Free Look Rules

#### BR-CLM-BOND-001: Freelook Period Calculation
- **ID**: BR-CLM-BOND-001
- **Category**: Free Look Cancellation
- **Priority**: CRITICAL
- **Description**: Physical bonds: 15 days from delivery. Online/ePLI: 30 days from issuance. Auto-reject after window expires.
- **Rule/Formula**: `freelook_period = IF bond_type = 'PHYSICAL' THEN (delivery_date + 15 days) ELSE (issuance_date + 30 days); IF cancellation_date > freelook_period THEN reject_cancellation = TRUE`
- **Source**: Claim_SRS_Tracking of Policy Bond.md, Lines 206-220, 85-87
- **Traceability**: VR-CLM-FL-001, FR-CLM-FL-001
- **Impact**: IRDAI compliance for free look period

#### BR-CLM-BOND-002: Delivery Failure Escalation
- **ID**: BR-CLM-BOND-002
- **Category**: Policy Bond Tracking
- **Priority**: HIGH
- **Description**: Flag undelivered bonds after 10 days. Escalate to CPC supervisor. Contact customer for address verification.
- **Rule/Formula**: `IF bond_delivery_status = 'UNDELIVERED' AND days_since_dispatch > 10 THEN FLAG_FOR_ESCALATION(); ESCALATE_TO(cpc_supervisor); CONTACT_CUSTOMER(address_verification)`
- **Source**: Claim_SRS_Tracking of Policy Bond.md, Lines 151-152
- **Traceability**: FR-CLM-FL-002
- **Impact**: Ensures policy bond delivery and customer awareness

#### BR-CLM-BOND-003: Refund Calculation
- **ID**: BR-CLM-BOND-003
- **Category**: Free Look Refund
- **Priority**: CRITICAL
- **Description**: Refund = Premium - (Pro-rata risk premium + Stamp duty + Medical costs + Other deductions).
- **Rule/Formula**: `refund_amount = total_premium - (pro_rata_risk_premium + stamp_duty + medical_costs + other_deductions)`
- **Source**: Claim_SRS_Tracking of Policy Bond.md, Lines 265-269
- **Traceability**: FR-CLM-FL-003
- **Impact**: Accurate refund computation per IRDAI guidelines

#### BR-CLM-BOND-004: Maker-Checker Workflow
- **ID**: BR-CLM-BOND-004
- **Category**: Free Look Refund
- **Priority**: CRITICAL
- **Description**: Maker enters refund, Checker verifies. Prevent same-person maker-checker. Generate unique transaction ID. Link to PLI finance.
- **Rule/Formula**: `REQUIRE_MAKER_CHECKER(maker_id != checker_id); transaction_id = GENERATE_UNIQUE_ID(); LINK_TO_FINANCE_SYSTEM(transaction_id, refund_amount)`
- **Source**: Claim_SRS_Tracking of Policy Bond.md, Lines 285-291
- **Traceability**: FR-CLM-FL-004
- **Impact**: Segregation of duties and financial control

---

## 3. Functional Requirements

### 3.1 Death Claim Requirements

#### FR-CLM-DC-001: Claim Registration
- **ID**: FR-CLM-DC-001
- **Category**: Death Claim
- **Priority**: CRITICAL
- **Description**: System must allow claim registration at any post office (BO/SO/CPC) with auto-generation of Claim ID
- **Acceptance Criteria**:
  - User can submit claim with nominee/legal heir/assignee details
  - System generates unique Claim ID immediately upon submission
  - Digital acknowledgment sent via SMS/Email
  - Claim tracked in system from registration
- **Source**: Claim_SRS FRS on death claim.md, Section 1
- **Business Rule**: BR-CLM-DC-010
- **Workflow**: WF-CLM-DC-001, Step 1
- **Data Entity**: E-CLM-DC-001 (Claim)

#### FR-CLM-DC-002: Document Capture & Indexing
- **ID**: FR-CLM-DC-002
- **Category**: Death Claim
- **Priority**: CRITICAL
- **Description**: System must scan, upload, and tag all supporting documents against Claim ID in ECMS
- **Acceptance Criteria**:
  - Documents scanned at BO/SO or uploaded by customer
  - Auto-tagging to Claim ID in ECMS
  - Mandatory document checklist validation
  - Missing document flag with auto-reminder after 7 days
  - Claim status = "Pending for Documents" if incomplete
  - Auto-return of claim after 22 days (15+7 grace)
- **Source**: Claim_SRS FRS on death claim.md, Section 2
- **Business Rule**: BR-CLM-DC-010, BR-CLM-IRDAI-002
- **Validation**: VR-CLM-DC-001 to VR-CLM-DC-006
- **Integration**: INT-CLM-001 (ECMS)

#### FR-CLM-DC-003: Benefit Calculation Engine
- **ID**: FR-CLM-DC-003
- **Category**: Death Claim - Calculation
- **Priority**: CRITICAL
- **Description**: System must compute claim amount including base sum assured, accrued bonuses, excess premiums, and deductions
- **Acceptance Criteria**:
  - Calculate base sum assured from policy master
  - Add accrued bonuses (simple and compound) from bonus ledger
  - Include excess premiums paid
  - Deduct outstanding loans from loan module
  - Deduct unpaid premiums
  - Deduct applicable taxes (TDS, GST)
  - Allow manual override in exceptional cases (court orders, disputed data)
  - Require supervisor approval for manual override
  - Log all calculations with user ID, timestamp, before/after values
  - Provide detailed breakup in settlement letter
- **Source**: Claim_SRS FRS on death claim.md, Section 4: Claim Calculation & Benefit Computation (Lines 88-100)
- **Business Rule**: BR-CLM-DC-001, BR-CLM-DC-002, BR-CLM-DC-003
- **Data Entity**: E-CLM-DC-002 (Claim Calculation)
- **Integration**: INT-CLM-002 (McCamish Policy System), INT-CLM-003 (Bonus Ledger), INT-CLM-004 (Loan Module)
- **Validation**: VR-CLM-DC-007 (Calculation validation)

#### FR-CLM-DC-004: Approval Workflow & Decision Points
- **ID**: FR-CLM-DC-004
- **Category**: Death Claim - Approval
- **Priority**: CRITICAL
- **Description**: System must route claims to appropriate approver based on financial limits and policy type, with SLA enforcement
- **Acceptance Criteria**:
  - Auto-route based on claim amount and approver financial limits
  - Display all indexed documents, investigation reports, calculated benefits to approver
  - Show SLA countdown: 15 days (no investigation) or 45 days (with investigation)
  - Approver actions: Approve, Reject, Send for Re-investigation
  - Rejection requires mandatory reason selection from dropdown
  - Auto-generate rejection letter with appellate rights information
  - Log approval decision with timestamp, approver ID, reason
  - Escalate to next level if SLA breached
  - Email/SMS notification to claimant on decision
- **Source**: Claim_SRS FRS on death claim.md, Section 5: Approval Workflow & Decision Points (Lines 102-112)
- **Business Rule**: BR-CLM-DC-004, BR-CLM-DC-005, BR-CLM-DC-006
- **Workflow**: WF-CLM-DC-001, Step 4 (Approval)
- **Data Entity**: E-CLM-DC-003 (Approval Decision)
- **Integration**: INT-CLM-005 (Notification Service)

#### FR-CLM-DC-005: Disbursement Payment Execution
- **ID**: FR-CLM-DC-005
- **Category**: Death Claim - Disbursement
- **Priority**: CRITICAL
- **Description**: System must execute payment via NEFT, POSB EFT, or cheque with account verification and payment reconciliation
- **Acceptance Criteria**:
  - Disbursement officer initiates payment post-approval
  - Payment mode selection: NEFT (preferred), POSB EFT, Cheque (fallback)
  - Integrate with Finacle or IT 2.0 for account verification
  - Validate bank account, IFSC code, account holder name
  - Execute payment with payment reference number
  - Update claim status to "Paid" in real-time
  - Generate sanction memo automatically
  - Log payment acknowledgment from bank
  - Reconcile payment status daily
  - Handle payment failures with auto-retry and alert mechanism
- **Source**: Claim_SRS FRS on death claim.md, Section 6: Disbursement & Payment Execution (Lines 114-121)
- **Business Rule**: BR-CLM-DC-007, BR-CLM-DC-008
- **Workflow**: WF-CLM-DC-001, Step 5 (Disbursement)
- **Data Entity**: E-CLM-DC-004 (Disbursement)
- **Integration**: INT-CLM-006 (Finacle/IT 2.0), INT-CLM-007 (NEFT/POSB Gateway)
- **Validation**: VR-CLM-DC-008 (Bank account validation)

#### FR-CLM-DC-006: Reopen & Exception Handling
- **ID**: FR-CLM-DC-006
- **Category**: Death Claim - Exception Management
- **Priority**: HIGH
- **Description**: System must allow claims to be reopened under valid circumstances with proper authorization and audit trail
- **Acceptance Criteria**:
  - Reopen reasons: Court orders, New evidence, Administrative lapses, Claimant appeals
  - CPC users or supervisors can initiate reopening via Service Request History screen
  - Generate new service request ID
  - Claim re-enters workflow at appropriate stage
  - Original claim history preserved
  - Link new request to original claim for traceability
  - Require supervisor approval for reopening
  - Document reason for reopening in system
  - Notify all stakeholders (claimant, approver, CPC team)
  - Tag as "Reopened" in claim tracker
- **Source**: Claim_SRS FRS on death claim.md, Section 7: Reopen & Exception Handling (Lines 123-129)
- **Business Rule**: BR-CLM-DC-013
- **Workflow**: WF-CLM-DC-REOPEN-001
- **Data Entity**: E-CLM-DC-005 (Reopen Request)
- **Integration**: INT-CLM-005 (Notification Service)

#### FR-CLM-DC-007: Communication & Notifications
- **ID**: FR-CLM-DC-007
- **Category**: Death Claim - Communication
- **Priority**: HIGH
- **Description**: System must send automated communications at all key milestones via SMS, email, and portal notifications
- **Acceptance Criteria**:
  - Notifications at milestones: Registration, Document status, Investigation, Approval/Rejection, Payment
  - Multi-channel delivery: SMS, Email, Portal notification
  - Include claim ID, status, next action required in all communications
  - Log all communications with timestamp, channel, delivery status
  - Attach communication to claim file
  - Flag communication failures with retry mechanism
  - Document corrective actions for audit
  - Allow claimant to set communication preferences
  - Provide unsubscribe option (except critical notifications)
  - Support multiple languages based on customer preference
- **Source**: Claim_SRS FRS on death claim.md, Section 8: Communication & Notifications (Lines 132-139)
- **Business Rule**: BR-CLM-DC-014
- **Workflow**: All claim workflow steps
- **Data Entity**: E-CLM-DC-006 (Communication Log)
- **Integration**: INT-CLM-008 (SMS Gateway), INT-CLM-009 (Email Service), INT-CLM-010 (Customer Portal)

#### FR-CLM-DC-008: Appeal Mechanism
- **ID**: FR-CLM-DC-008
- **Category**: Death Claim - Appeal
- **Priority**: HIGH
- **Description**: System must allow claimants to file appeals within 90 days of rejection with defined workflow and 45-day decision timeline
- **Acceptance Criteria**:
  - Claimant can file appeal within 90 days of rejection notification
  - Appeal submission channels: Post, Email, In-person at CPC
  - System validates appeal is within 90-day window
  - Allow late appeals with condonation request and valid justification
  - Auto-route to appellate authority (next higher officer in approval hierarchy)
  - Appellate authority can request additional documents or investigation reports
  - Decision timeline: 45 days from appeal receipt
  - SLA countdown visible to appellate authority
  - Decision must include detailed justification, rulings, supporting documents
  - Auto-generate appellate order with digital signature
  - Log all appeal actions for traceability
  - Link appeal to original claim for complete history
  - Notify claimant of appeal outcome via registered post, email, SMS
- **Source**: Claim_SRS FRS on death claim.md, Section 9: Appeal Mechanism (Lines 141-151)
- **Business Rule**: BR-CLM-DC-015, BR-CLM-DC-016
- **Workflow**: WF-CLM-DC-APPEAL-001
- **Data Entity**: E-CLM-DC-007 (Appeal)
- **Integration**: INT-CLM-005 (Notification Service)
- **Validation**: VR-CLM-DC-009 (Appeal timeline validation)

### 3.2 Maturity Claim Requirements

#### FR-CLM-MC-001: Maturity Report Generation
- **ID**: FR-CLM-MC-001
- **Category**: Maturity Claim
- **Priority**: HIGH
- **Description**: System must auto-generate maturity due report on first working day of month for policies maturing in next 2 months
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-01
- **Business Rule**: BR-CLM-MC-001
- **Workflow**: Scheduled batch job (monthly)

#### FR-CLM-MC-002: Multi-Channel Intimation
- **ID**: FR-CLM-MC-002
- **Category**: Maturity Claim - Communication
- **Priority**: HIGH
- **Description**: System must automatically send maturity intimation via SMS, Email, WhatsApp, and Portal with secure claim submission link
- **Acceptance Criteria**:
  - Auto-send intimation 60 days before maturity date
  - Send via SMS, Email, WhatsApp, Portal notification
  - Use Registered Post only as fallback channel
  - Include secure link for online claim submission
  - Include prefilled claim form download option
  - Include policy details, maturity amount, required documents list
  - Log all intimation attempts with delivery status
  - Retry failed communications
  - Track customer acknowledgment
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-02 (Lines 449-459)
- **Business Rule**: BR-CLM-MC-002
- **Workflow**: WF-CLM-MC-INTIMATION-001
- **Data Entity**: E-CLM-MC-001 (Intimation Log)
- **Integration**: INT-CLM-008 (SMS Gateway), INT-CLM-009 (Email Service), INT-CLM-011 (WhatsApp Business API), INT-CLM-010 (Customer Portal)

#### FR-CLM-MC-003: Customer-Initiated Claim Submission
- **ID**: FR-CLM-MC-003
- **Category**: Maturity Claim - Registration
- **Priority**: HIGH
- **Description**: System must allow customers to initiate maturity claims via Portal or Mobile app with document upload and DigiLocker integration
- **Acceptance Criteria**:
  - Customer can initiate claim via Portal or Mobile app
  - Upload documents: Policy bond, ID proof, Bank details, Discharge form
  - Integrate with DigiLocker for document fetching
  - Auto-populate policy details from system
  - Real-time document format and size validation
  - Generate auto-acknowledgment with unique Claim ID upon submission
  - Include submission timestamp in acknowledgment
  - Send acknowledgment via SMS and Email immediately
  - Allow save-as-draft functionality
  - Support multiple document formats (PDF, JPEG, PNG)
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-03 (Lines 461-475)
- **Business Rule**: BR-CLM-MC-003
- **Workflow**: WF-CLM-MC-001, Step 1 (Customer Submission)
- **Data Entity**: E-CLM-MC-002 (Online Claim Submission)
- **Integration**: INT-CLM-010 (Customer Portal), INT-CLM-012 (Mobile App), INT-CLM-013 (DigiLocker)
- **Validation**: VR-CLM-MC-001 (Document validation)

#### FR-CLM-MC-004: System-Assisted Initial Scrutiny
- **ID**: FR-CLM-MC-004
- **Category**: Maturity Claim - Document Verification
- **Priority**: HIGH
- **Description**: System must validate uploaded documents against checklist and flag missing or invalid items with auto-reminders
- **Acceptance Criteria**:
  - Validate documents against mandatory checklist instantly
  - Flag missing documents with specific list
  - Flag invalid documents (wrong format, corrupted, expired)
  - Send auto-reminders via SMS/Email/WhatsApp for outstanding documents
  - Reminder schedule: Day 1, Day 5, Day 10, Day 15
  - CPC staff verify completeness digitally via system interface
  - System shows document completeness percentage
  - Mark claim as "DOCUMENT_COMPLETE" or "DOCUMENT_PENDING"
  - Cannot proceed to next stage until documents complete
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-04 (Lines 477-489)
- **Business Rule**: BR-CLM-MC-004
- **Workflow**: WF-CLM-MC-001, Step 2 (Initial Scrutiny)
- **Data Entity**: E-CLM-MC-003 (Document Checklist)
- **Validation**: VR-CLM-MC-002 to VR-CLM-MC-005
- **Integration**: INT-CLM-008 (SMS), INT-CLM-009 (Email), INT-CLM-011 (WhatsApp)

#### FR-CLM-MC-005: Auto-Indexing and Document Sync
- **ID**: FR-CLM-MC-005
- **Category**: Maturity Claim - Document Management
- **Priority**: MEDIUM
- **Description**: System must automatically index claims with metadata and sync documents to ECMS in real-time
- **Acceptance Criteria**:
  - Auto-index claim with metadata: Claim ID, Policy number, Customer ID, Submission date
  - Sync documents to ECMS in real-time
  - Link documents to Claim ID and Policy record
  - Generate unique document ID for each uploaded file
  - Tag documents with category (ID proof, bank details, discharge form, etc.)
  - Enable search and retrieval by multiple parameters
  - Maintain document version history
  - Audit log of all document operations
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-05 (Lines 491-497)
- **Business Rule**: BR-CLM-MC-005
- **Workflow**: WF-CLM-MC-001, Step 2.5 (Document Indexing)
- **Data Entity**: E-CLM-MC-004 (Document Index)
- **Integration**: INT-CLM-001 (ECMS)

#### FR-CLM-MC-006: Auto-Populated Data Entry (OCR)
- **ID**: FR-CLM-MC-006
- **Category**: Maturity Claim - Data Entry
- **Priority**: MEDIUM
- **Description**: System must extract key fields from uploaded documents using OCR to auto-populate claim data fields
- **Acceptance Criteria**:
  - Use OCR to extract data from uploaded documents
  - Auto-populate fields: Name, Policy number, Bank account, IFSC, Address
  - CPC operator reviews and confirms extracted data
  - Highlight low-confidence extractions for manual verification
  - Allow manual edit of extracted data
  - Log data extraction accuracy metrics
  - Reduce manual entry time by 70%
  - Reduce data entry errors
  - Support multiple document formats and languages
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-06 (Lines 499-507)
- **Business Rule**: BR-CLM-MC-006
- **Workflow**: WF-CLM-MC-001, Step 3 (Data Entry)
- **Data Entity**: E-CLM-MC-005 (OCR Data)
- **Integration**: INT-CLM-014 (OCR Engine)
- **Validation**: VR-CLM-MC-006 (Data accuracy validation)

#### FR-CLM-MC-007: QC Verification Checklist
- **ID**: FR-CLM-MC-007
- **Category**: Maturity Claim - Quality Control
- **Priority**: HIGH
- **Description**: System must enforce QC verification with system-enforced checklist and dual authentication for override
- **Acceptance Criteria**:
  - Supervisor performs QC digitally using system-enforced checklist
  - Checklist items: Document completeness, Data accuracy, Eligibility check, Calculation verification
  - All checklist items must be marked complete
  - Dual authentication required for override or waiver
  - Log all QC actions with user ID and timestamp
  - Flag discrepancies for resolution
  - Cannot proceed without QC approval
  - QC dashboard shows pending items
  - Track QC TAT and quality metrics
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-07 (Lines 509-517)
- **Business Rule**: BR-CLM-MC-007
- **Workflow**: WF-CLM-MC-001, Step 3.5 (QC Verification)
- **Data Entity**: E-CLM-MC-006 (QC Checklist)
- **Validation**: VR-CLM-MC-007 (QC completeness)

#### FR-CLM-MC-008: Approval Workflow and SLA Enforcement
- **ID**: FR-CLM-MC-008
- **Category**: Maturity Claim - Approval
- **Priority**: CRITICAL
- **Description**: System must route claims to approving authority with SLA countdown and digital signature capability
- **Acceptance Criteria**:
  - Approving authority reviews via dedicated digital dashboard
  - Display documents, system-calculated maturity amount, QC checklist
  - Show SLA countdown (7-day window from submission)
  - Approver actions: Approve or Redirect with remarks
  - Digital signature for approval
  - Log approval decision with timestamp, approver ID
  - Auto-escalate if SLA breached
  - Email notification to customer on approval
  - Integration with payment module for disbursement
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-08 (Lines 519-529)
- **Business Rule**: BR-CLM-MC-008
- **Workflow**: WF-CLM-MC-001, Step 4 (Approval)
- **Data Entity**: E-CLM-MC-007 (Approval)
- **Integration**: INT-CLM-015 (Digital Signature), INT-CLM-005 (Notification)

#### FR-CLM-MC-009: Auto-Generated Sanction/Rejection Communication
- **ID**: FR-CLM-MC-009
- **Category**: Maturity Claim - Communication
- **Priority**: HIGH
- **Description**: System must auto-generate and send sanction or rejection letters via email, WhatsApp, and Customer Portal
- **Acceptance Criteria**:
  - Auto-generate sanction or rejection letter upon approval decision
  - Include timestamp (date/time/second)
  - Send via Email, WhatsApp, Customer Portal
  - If rejected, include reason and appeal link
  - Letters timestamped and archived in ECMS
  - Delivery confirmation tracking
  - Provide downloadable PDF copy
  - Maintain letter templates with merge fields
  - Support multiple languages
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-09 (Lines 531-540)
- **Business Rule**: BR-CLM-MC-009
- **Workflow**: WF-CLM-MC-001, Step 4.5 (Communication)
- **Data Entity**: E-CLM-MC-008 (Letter)
- **Integration**: INT-CLM-009 (Email), INT-CLM-011 (WhatsApp), INT-CLM-010 (Portal), INT-CLM-001 (ECMS)

#### FR-CLM-MC-010: Bank Account Validation (API-based)
- **ID**: FR-CLM-MC-010
- **Category**: Maturity Claim - Payment Validation
- **Priority**: CRITICAL
- **Description**: System must verify bank account details via CBS/PFMS API before disbursement
- **Acceptance Criteria**:
  - Trigger API-based validation upon bank detail submission
  - Validate: Account number, IFSC code, Account holder name, Account status (active)
  - API response processed for validity determination
  - If successful, mark as "Verified" and proceed
  - If failed, prevent submission and display error message
  - Error types: Invalid account, Account not found, Name mismatch, Inactive account
  - Prompt customer/staff to correct details
  - Retry mechanism for API failures
  - Log all validation attempts
  - Support multiple bank APIs (CBS, PFMS)
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-10 (Lines 542-546) and FRS-MAT-17 (Lines 594-639)
- **Business Rule**: BR-CLM-MC-010
- **Workflow**: WF-CLM-MC-001, Step 5 (Bank Validation)
- **Data Entity**: E-CLM-MC-009 (Bank Validation)
- **Integration**: INT-CLM-016 (CBS API), INT-CLM-017 (PFMS API)
- **Validation**: VR-CLM-MC-008 (Bank account format validation)

#### FR-CLM-MC-011: Disbursement Execution
- **ID**: FR-CLM-MC-011
- **Category**: Maturity Claim - Disbursement
- **Priority**: CRITICAL
- **Description**: System must process payment using Auto NEFT/IMPS with real-time status updates
- **Acceptance Criteria**:
  - Integrated with Core Banking system
  - Process payment using Auto NEFT/IMPS (preferred mode)
  - Execute payment post-approval and bank validation
  - Generate payment reference number
  - Update disbursement status in real-time
  - Handle payment failures with auto-retry
  - Alert disbursement officer on failure
  - Reconcile payments daily
  - Update claim status to "Paid" upon success
  - Generate payment confirmation
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-11 (Lines 548-554)
- **Business Rule**: BR-CLM-MC-011
- **Workflow**: WF-CLM-MC-001, Step 6 (Disbursement)
- **Data Entity**: E-CLM-MC-010 (Disbursement)
- **Integration**: INT-CLM-006 (Core Banking), INT-CLM-007 (NEFT/IMPS Gateway)

#### FR-CLM-MC-012: Voucher Generation and Submission
- **ID**: FR-CLM-MC-012
- **Category**: Maturity Claim - Accounting
- **Priority**: MEDIUM
- **Description**: System must auto-generate payment voucher post-disbursement and submit digitally to Accounts section
- **Acceptance Criteria**:
  - Auto-generate voucher post-disbursement
  - Include: Claim ID, Policy number, Payment amount, Payment date, Payment reference, Bank details
  - Submit digitally to Accounts section
  - Link voucher to claim record for audit trail
  - Generate voucher in standard accounting format
  - Support voucher approval workflow in Accounts
  - Enable voucher search and retrieval
  - Maintain voucher numbering sequence
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-12 (Lines 556-562)
- **Business Rule**: BR-CLM-MC-012
- **Workflow**: WF-CLM-MC-001, Step 7 (Voucher Generation)
- **Data Entity**: E-CLM-MC-011 (Voucher)
- **Integration**: INT-CLM-018 (Accounts System)

#### FR-CLM-MC-013: Claim Closure and Archiving
- **ID**: FR-CLM-MC-013
- **Category**: Maturity Claim - Closure
- **Priority**: MEDIUM
- **Description**: System must auto-mark claims as Paid or Rejected and archive case file digitally with closure timestamp
- **Acceptance Criteria**:
  - Auto-update claim status to "Paid" or "Rejected"
  - Record closure timestamp and user ID
  - Archive complete case file digitally in ECMS
  - Include all documents, communications, approvals, payments
  - Generate closure summary report
  - Enable post-closure retrieval for audit
  - Maintain archival for regulatory retention period
  - Support bulk archival for old claims
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-13 (Lines 564-569)
- **Business Rule**: BR-CLM-MC-013
- **Workflow**: WF-CLM-MC-001, Step 8 (Closure)
- **Data Entity**: E-CLM-MC-012 (Claim Closure)
- **Integration**: INT-CLM-001 (ECMS)

#### FR-CLM-MC-014: Customer Feedback Collection
- **ID**: FR-CLM-MC-014
- **Category**: Maturity Claim - Feedback
- **Priority**: LOW
- **Description**: System must send auto-message to customer post-settlement to collect feedback for service quality analytics
- **Acceptance Criteria**:
  - Send auto-message 2-3 days post-claim settlement
  - Message channels: SMS, Email, WhatsApp with feedback link
  - Feedback form: Rating (1-5 stars), Comments, Service quality parameters
  - Store feedback in system with claim ID linkage
  - Generate feedback analytics reports
  - Track NPS (Net Promoter Score)
  - Escalate negative feedback to management
  - Anonymous feedback option available
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-14 (Lines 571-577)
- **Business Rule**: BR-CLM-MC-014
- **Workflow**: WF-CLM-MC-FEEDBACK-001
- **Data Entity**: E-CLM-MC-013 (Feedback)
- **Integration**: INT-CLM-008 (SMS), INT-CLM-009 (Email), INT-CLM-011 (WhatsApp), INT-CLM-019 (Feedback Portal)

#### FR-CLM-MC-015: Real-Time Monitoring & Escalation
- **ID**: FR-CLM-MC-015
- **Category**: Maturity Claim - Monitoring
- **Priority**: HIGH
- **Description**: System must provide admin dashboard showing pending claims, SLA countdown, and escalated cases with auto-escalation triggers
- **Acceptance Criteria**:
  - Admin dashboard shows: Pending claims count, SLA countdown by stage, Escalated cases, Breached SLAs
  - Real-time data refresh
  - Auto-escalation triggered if SLA breached
  - Color-coded alerts: Green (on-track), Yellow (near-breach), Red (breached)
  - Drill-down capability to claim details
  - Export dashboard data to Excel/PDF
  - Role-based dashboard views
  - Configurable alert thresholds
  - Email alerts for escalations
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-15 (Lines 579-584)
- **Business Rule**: BR-CLM-MC-015
- **Workflow**: Real-time monitoring
- **Data Entity**: E-CLM-MC-014 (Dashboard Metrics)
- **Integration**: INT-CLM-020 (BI Dashboard), INT-CLM-005 (Notification)

#### FR-CLM-MC-016: Customer Claim Tracker
- **ID**: FR-CLM-MC-016
- **Category**: Maturity Claim - Customer Self-Service
- **Priority**: MEDIUM
- **Description**: System must provide online/mobile access for customers to track claim status with stage-wise updates and timestamps
- **Acceptance Criteria**:
  - Customer accesses tracker via Portal or Mobile app
  - Login with policy number and OTP
  - Display claim stages: Submission, Scrutiny, QC, Approval, Payment
  - Show current stage with timestamp
  - Show next expected action and timeline
  - Display estimated completion date
  - Show all communications sent
  - Enable document upload for pending documents
  - Show payment status and reference number
  - Support multiple languages
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-16 (Lines 586-592)
- **Business Rule**: BR-CLM-MC-016
- **Workflow**: Customer self-service
- **Data Entity**: E-CLM-MC-015 (Claim Tracker)
- **Integration**: INT-CLM-010 (Customer Portal), INT-CLM-012 (Mobile App)

### 3.3 Survival Benefit Requirements

#### FR-CLM-SB-001: SB Report Auto-Generation
- **ID**: FR-CLM-SB-001
- **Category**: Survival Benefit
- **Priority**: HIGH
- **Description**: System must auto-generate survival benefit due report
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-01
- **Business Rule**: BR-CLM-SB-001

#### FR-CLM-SB-002: Online Submission with DigiLocker
- **ID**: FR-CLM-SB-002
- **Category**: Survival Benefit
- **Priority**: HIGH
- **Description**: System must allow online SB claim submission with DigiLocker integration
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-03
- **Integration**: INT-CLM-007 (DigiLocker)

#### FR-CLM-SB-003: Multi-Channel Intimation
- **ID**: FR-CLM-SB-003
- **Category**: Survival Benefit - Communication
- **Priority**: HIGH
- **Description**: System must automatically send SB intimation via SMS, Email, WhatsApp with secure link for online claim submission
- **Acceptance Criteria**:
  - Auto-send intimation 30 days before SB due date
  - Send via SMS, Email, WhatsApp
  - Use Registered Post as fallback channel
  - Include secure link for online claim submission
  - Include policy details, SB amount, required documents
  - Log all intimation attempts with delivery status
  - Retry failed communications
  - Track customer acknowledgment
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-02 (Lines 331-339)
- **Business Rule**: BR-CLM-SB-002
- **Workflow**: WF-CLM-SB-INTIMATION-001
- **Data Entity**: E-CLM-SB-001 (Intimation Log)
- **Integration**: INT-CLM-008 (SMS), INT-CLM-009 (Email), INT-CLM-011 (WhatsApp)

#### FR-CLM-SB-004: Initial Scrutiny (Digital)
- **ID**: FR-CLM-SB-004
- **Category**: Survival Benefit - Document Verification
- **Priority**: HIGH
- **Description**: System must instantly flag missing documents and send auto-reminders with CPC staff digital verification
- **Acceptance Criteria**:
  - Instantly flag missing documents upon submission
  - Send auto-reminders via SMS/Email/WhatsApp for outstanding documents
  - Reminder schedule: Day 1, Day 7, Day 14
  - CPC staff verify completeness digitally via system interface
  - Display document checklist with completion status
  - Mark claim as "DOCUMENT_COMPLETE" or "DOCUMENT_PENDING"
  - Cannot proceed without complete documents
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-04 (Lines 350-356)
- **Business Rule**: BR-CLM-SB-003
- **Workflow**: WF-CLM-SB-001, Step 2 (Initial Scrutiny)
- **Data Entity**: E-CLM-SB-002 (Document Checklist)
- **Validation**: VR-CLM-SB-001 to VR-CLM-SB-003
- **Integration**: INT-CLM-008 (SMS), INT-CLM-009 (Email), INT-CLM-011 (WhatsApp)

#### FR-CLM-SB-005: Auto-Indexing in IMS 2.0
- **ID**: FR-CLM-SB-005
- **Category**: Survival Benefit - Document Management
- **Priority**: MEDIUM
- **Description**: System must automatically index claim as Service Request and link to policy and Claim ID
- **Acceptance Criteria**:
  - Auto-index claim as Service Request upon final document submission
  - Link Service Request to policy number
  - Link Service Request to Claim ID
  - Generate unique Service Request ID
  - Tag with claim type "SURVIVAL_BENEFIT"
  - Enable search by policy, claim ID, or service request ID
  - Maintain service request history
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-05 (Lines 358-362)
- **Business Rule**: BR-CLM-SB-004
- **Workflow**: WF-CLM-SB-001, Step 2.5 (Indexing)
- **Data Entity**: E-CLM-SB-003 (Service Request)
- **Integration**: INT-CLM-021 (IMS 2.0)

#### FR-CLM-SB-006: Document Scanning & Upload
- **ID**: FR-CLM-SB-006
- **Category**: Survival Benefit - Document Management
- **Priority**: MEDIUM
- **Description**: System must support document scanning at source (BO/SO) or customer upload with auto-tagging in ECMS
- **Acceptance Criteria**:
  - Support document scanning at BO/SO
  - Support customer upload via Portal/Mobile app
  - Auto-tag documents with category
  - Store in ECMS with unique document ID
  - Link to Claim ID and Policy number
  - Support formats: PDF, JPEG, PNG
  - Validate file size (max 5MB per document)
  - Virus scan all uploaded documents
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-06 (Lines 364-369)
- **Business Rule**: BR-CLM-SB-005
- **Workflow**: WF-CLM-SB-001, Step 2.6 (Document Upload)
- **Data Entity**: E-CLM-SB-004 (Document)
- **Integration**: INT-CLM-001 (ECMS), INT-CLM-022 (Virus Scanner)
- **Validation**: VR-CLM-SB-004 (Document format validation)

#### FR-CLM-SB-007: Data Entry & QC Verification (Automated)
- **ID**: FR-CLM-SB-007
- **Category**: Survival Benefit - Data Entry & QC
- **Priority**: MEDIUM
- **Description**: System must auto-populate claim data from documents using OCR with CPC supervisor QC verification
- **Acceptance Criteria**:
  - Use OCR/data extraction to auto-populate fields
  - Fields: Name, Policy number, Bank account, IFSC, Address
  - CPC supervisor performs QC digitally via dashboard
  - Verify accuracy of auto-populated data
  - Allow manual correction if needed
  - Log data extraction confidence scores
  - Flag low-confidence extractions for manual review
  - QC approval required before proceeding
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-07 (Lines 371-377)
- **Business Rule**: BR-CLM-SB-006
- **Workflow**: WF-CLM-SB-001, Step 3 (Data Entry & QC)
- **Data Entity**: E-CLM-SB-005 (Claim Data), E-CLM-SB-006 (QC Record)
- **Integration**: INT-CLM-014 (OCR Engine)
- **Validation**: VR-CLM-SB-005 (Data accuracy validation)

#### FR-CLM-SB-008: Approval Workflow
- **ID**: FR-CLM-SB-008
- **Category**: Survival Benefit - Approval
- **Priority**: CRITICAL
- **Description**: System must route to Postmaster/Approving Authority for review with digital signature and 7-day SLA enforcement
- **Acceptance Criteria**:
  - Postmaster/Approving Authority reviews via dedicated dashboard
  - Display all documents, data, QC checklist
  - Show calculated SB amount
  - Approver actions: Approve or Reject with digital signature
  - Enforce 7-day SLA from submission
  - SLA countdown visible to approver
  - Log approval decision with timestamp, approver ID
  - Auto-escalate if SLA breached
  - Rejection requires reason selection
  - Email/SMS notification to customer on decision
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-08 (Lines 379-385)
- **Business Rule**: BR-CLM-SB-007
- **Workflow**: WF-CLM-SB-001, Step 4 (Approval)
- **Data Entity**: E-CLM-SB-007 (Approval)
- **Integration**: INT-CLM-015 (Digital Signature), INT-CLM-005 (Notification)

#### FR-CLM-SB-009: Sanction/Rejection Letter Generation
- **ID**: FR-CLM-SB-009
- **Category**: Survival Benefit - Communication
- **Priority**: HIGH
- **Description**: System must auto-generate sanction or rejection letter with timestamp and send via email, WhatsApp, and Customer Portal
- **Acceptance Criteria**:
  - Auto-generate letter upon approval decision
  - Include timestamp (date/time/second)
  - Send via Email, WhatsApp, Customer Portal
  - If rejected, include reason and appeal link
  - Letters archived in ECMS
  - Delivery confirmation tracking
  - Provide downloadable PDF copy
  - Support multiple languages
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-09 (Lines 387-393)
- **Business Rule**: BR-CLM-SB-008
- **Workflow**: WF-CLM-SB-001, Step 4.5 (Letter Generation)
- **Data Entity**: E-CLM-SB-008 (Letter)
- **Integration**: INT-CLM-009 (Email), INT-CLM-011 (WhatsApp), INT-CLM-010 (Portal), INT-CLM-001 (ECMS)

#### FR-CLM-SB-010: Bank Account Verification
- **ID**: FR-CLM-SB-010
- **Category**: Survival Benefit - Payment Validation
- **Priority**: CRITICAL
- **Description**: System must use API-based validation for bank account details with correction prompts on failure
- **Acceptance Criteria**:
  - Trigger API-based validation for bank details
  - Validate: Account number, IFSC code, Account holder name, Account status
  - If successful, mark as "Verified"
  - If failed, prompt user/staff for correction
  - Error messages: Invalid account, Name mismatch, Inactive account
  - Retry mechanism for API failures
  - Log all validation attempts
  - Cannot proceed to disbursement without verification
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-10 (Lines 395-399)
- **Business Rule**: BR-CLM-SB-009
- **Workflow**: WF-CLM-SB-001, Step 5 (Bank Validation)
- **Data Entity**: E-CLM-SB-009 (Bank Validation)
- **Integration**: INT-CLM-016 (CBS API), INT-CLM-017 (PFMS API)
- **Validation**: VR-CLM-SB-006 (Bank account validation)

#### FR-CLM-SB-011: Disbursement
- **ID**: FR-CLM-SB-011
- **Category**: Survival Benefit - Disbursement
- **Priority**: CRITICAL
- **Description**: System must process payment using Auto NEFT/IMPS with automatic status updates in IMS 2.0
- **Acceptance Criteria**:
  - Integrated with Core Banking system
  - Process payment using Auto NEFT/IMPS
  - Execute payment post-approval and bank validation
  - Generate payment reference number
  - Update disbursement status in IMS 2.0 automatically
  - Handle payment failures with auto-retry
  - Alert disbursement officer on failure
  - Reconcile payments daily
  - Update claim status to "Paid" upon success
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-11 (Lines 401-405)
- **Business Rule**: BR-CLM-SB-010
- **Workflow**: WF-CLM-SB-001, Step 6 (Disbursement)
- **Data Entity**: E-CLM-SB-010 (Disbursement)
- **Integration**: INT-CLM-006 (Core Banking), INT-CLM-007 (NEFT/IMPS), INT-CLM-021 (IMS 2.0)

#### FR-CLM-SB-012: Voucher Submission
- **ID**: FR-CLM-SB-012
- **Category**: Survival Benefit - Accounting
- **Priority**: MEDIUM
- **Description**: System must auto-generate payment voucher and submit digitally to Accounts section linked to disbursement record
- **Acceptance Criteria**:
  - Auto-generate voucher post-payment
  - Include: Claim ID, Policy number, Payment amount, Payment date, Payment reference
  - Submit digitally to Accounts section
  - Link to disbursement record
  - Generate in standard accounting format
  - Maintain voucher numbering sequence
  - Enable voucher search and audit
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-12 (Lines 407-411)
- **Business Rule**: BR-CLM-SB-011
- **Workflow**: WF-CLM-SB-001, Step 7 (Voucher Generation)
- **Data Entity**: E-CLM-SB-011 (Voucher)
- **Integration**: INT-CLM-018 (Accounts System)

#### FR-CLM-SB-013: Customer Feedback Collection
- **ID**: FR-CLM-SB-013
- **Category**: Survival Benefit - Feedback
- **Priority**: LOW
- **Description**: System must send auto-message to customer post-settlement to collect feedback for service quality monitoring
- **Acceptance Criteria**:
  - Send auto-message 2-3 days post-settlement
  - Message channels: SMS, Email, WhatsApp with feedback link
  - Feedback form: Rating, Comments, Service quality parameters
  - Store feedback in system
  - Generate analytics reports
  - Track NPS
  - Escalate negative feedback
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-13 (Lines 413-416)
- **Business Rule**: BR-CLM-SB-012
- **Workflow**: WF-CLM-SB-FEEDBACK-001
- **Data Entity**: E-CLM-SB-012 (Feedback)
- **Integration**: INT-CLM-008 (SMS), INT-CLM-009 (Email), INT-CLM-011 (WhatsApp), INT-CLM-019 (Feedback Portal)

#### FR-CLM-SB-014: Monitoring & Escalation
- **ID**: FR-CLM-SB-014
- **Category**: Survival Benefit - Monitoring
- **Priority**: HIGH
- **Description**: System must provide real-time dashboard for Admin Office with SLA countdown, color-coded alerts, and auto-escalation
- **Acceptance Criteria**:
  - Real-time dashboard for Admin Office
  - Display: Pending claims, SLA countdown, Escalated cases, Breached SLAs
  - Color-coded alerts: Green (on-track), Yellow (near-breach), Red (breached)
  - Auto-escalation triggered if pending beyond threshold
  - Drill-down to claim details
  - Export dashboard to Excel/PDF
  - Role-based views
  - Email alerts for escalations
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-14 (Lines 418-425)
- **Business Rule**: BR-CLM-SB-013
- **Workflow**: Real-time monitoring
- **Data Entity**: E-CLM-SB-013 (Dashboard Metrics)
- **Integration**: INT-CLM-020 (BI Dashboard), INT-CLM-005 (Notification)

#### FR-CLM-SB-015: Customer Claim Tracker
- **ID**: FR-CLM-SB-015
- **Category**: Survival Benefit - Customer Self-Service
- **Priority**: MEDIUM
- **Description**: System must provide online/mobile access for customer to track claim status with stage-wise updates
- **Acceptance Criteria**:
  - Customer accesses via Portal or Mobile app
  - Login with policy number and OTP
  - Display claim stages: Submission, Scrutiny, Approval, Payment
  - Show current stage with timestamp
  - Show next expected action
  - Display estimated completion date
  - Enable document upload for pending documents
  - Show payment status and reference
  - Support multiple languages
- **Source**: Claim_SRS FRS on survival benefit.md, FRS-SB-15 (Lines 427-431)
- **Business Rule**: BR-CLM-SB-014
- **Workflow**: Customer self-service
- **Data Entity**: E-CLM-SB-014 (Claim Tracker)
- **Integration**: INT-CLM-010 (Customer Portal), INT-CLM-012 (Mobile App)

### 3.4 AML Requirements

#### FR-CLM-AML-001: High Cash Premium Detection
- **ID**: FR-CLM-AML-001
- **Category**: AML/CFT - High Value Transaction
- **Priority**: CRITICAL
- **Description**: System must detect and alert when cash premium payment exceeds ₹50,000
- **Acceptance Criteria**:
  - Auto-trigger when cash_amount > ₹50,000
  - Risk level = High
  - Generate AML alert in compliance dashboard
  - Auto-initiate CTR (Cash Transaction Report) filing
  - Alert compliance officer immediately
  - Log transaction details with timestamp, user ID, policy number, amount
  - Transaction proceeds but is flagged for review
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 6: Trigger Logic Definitions, AML_001 (Lines 254-257)
- **Business Rule**: BR-CLM-AML-001
- **Workflow**: WF-CLM-AML-CTR-001
- **Data Entity**: E-CLM-AML-ALERT-001
- **Integration**: INT-CLM-AML-CTR (CTR filing system)

#### FR-CLM-AML-002: PAN Mismatch Detection
- **ID**: FR-CLM-AML-002
- **Category**: AML/CFT - Identity Verification
- **Priority**: HIGH
- **Description**: System must detect and alert when PAN verification fails
- **Acceptance Criteria**:
  - Auto-trigger when pan_verified = false
  - Risk level = Medium
  - Generate AML alert for manual review
  - Flag claim for additional KYC verification
  - Route to compliance team for verification
  - Claim status updated to "PENDING_KYC_VERIFICATION"
  - Cannot proceed to disbursement until PAN verified
  - Log all verification attempts and results
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 6: Trigger Logic Definitions, AML_002 (Lines 259-262)
- **Business Rule**: BR-CLM-AML-002
- **Workflow**: WF-CLM-AML-KYC-001
- **Data Entity**: E-CLM-AML-ALERT-002
- **Validation**: VR-CLM-AML-001 (PAN format and verification)

#### FR-CLM-AML-003: Nominee Change Post Death Detection
- **ID**: FR-CLM-AML-003
- **Category**: AML/CFT - Fraud Detection
- **Priority**: CRITICAL
- **Description**: System must detect and block claims where nominee was changed after death date
- **Acceptance Criteria**:
  - Auto-trigger when nominee_change_date > death_date
  - Risk level = Critical
  - Auto-block transaction immediately
  - Claim status set to "BLOCKED_FRAUD_INVESTIGATION"
  - Generate STR (Suspicious Transaction Report) automatically
  - Escalate to fraud investigation team
  - Prevent any claim processing until resolved
  - Notify branch manager and compliance head
  - Require dual authorization to unblock (fraud head + regional manager)
  - Complete audit trail of nominee change history
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 6: Trigger Logic Definitions, AML_003 (Lines 263-266)
- **Business Rule**: BR-CLM-AML-003
- **Workflow**: WF-CLM-AML-FRAUD-001
- **Data Entity**: E-CLM-AML-ALERT-003
- **Integration**: INT-CLM-AML-STR (STR filing system)

#### FR-CLM-AML-004: Frequent Surrender Pattern Detection
- **ID**: FR-CLM-AML-004
- **Category**: AML/CFT - Pattern Analysis
- **Priority**: HIGH
- **Description**: System must detect patterns of frequent surrenders indicating potential money laundering
- **Acceptance Criteria**:
  - Auto-trigger when customer has >3 surrenders within 6 months
  - Risk level = Medium
  - Generate AML alert for investigation
  - Flag customer profile for enhanced monitoring
  - Route to AML investigation team
  - Check all policies linked to same customer ID, PAN, or address
  - Include family group analysis (linked accounts)
  - Log surrender frequency metrics in customer risk profile
  - Current transaction proceeds but flagged for review
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 6: Trigger Logic Definitions, AML_004 (Lines 267-271)
- **Business Rule**: BR-CLM-AML-004
- **Workflow**: WF-CLM-AML-INVESTIGATION-001
- **Data Entity**: E-CLM-AML-ALERT-004
- **Validation**: VR-CLM-AML-002 (Surrender frequency calculation)

#### FR-CLM-AML-005: Refund Without Bond Delivery Detection
- **ID**: FR-CLM-AML-005
- **Category**: AML/CFT - Process Violation
- **Priority**: HIGH
- **Description**: System must detect and alert when refund is processed before policy bond is dispatched
- **Acceptance Criteria**:
  - Auto-trigger when refund_date < bond_dispatch_date
  - Risk level = High
  - Generate AML alert
  - Log event in audit trail with full transaction details
  - Route to compliance for review
  - Flag for potential process bypass
  - Investigate if proper authorization was obtained
  - Check if emergency refund approval exists
  - Alert operations head and compliance officer
  - Transaction logged but not blocked (post-facto detection)
- **Source**: Claim_SRS_AML triggers & alerts.md, Section 6: Trigger Logic Definitions, AML_005 (Lines 273-276)
- **Business Rule**: BR-CLM-AML-005
- **Workflow**: WF-CLM-AML-AUDIT-001
- **Data Entity**: E-CLM-AML-ALERT-005
- **Integration**: INT-CLM-BOND-TRACKER (Bond dispatch tracking system)

---

### 3.5 Insurance Ombudsman Requirements

#### FR-CLM-OMB-001: Complaint Intake & Registration
- **ID**: FR-CLM-OMB-001
- **Category**: Insurance Ombudsman - Intake
- **Priority**: CRITICAL
- **Description**: System must provide multichannel complaint registration supporting digital, written, and assisted modes compliant with Rule 14 of Insurance Ombudsman Rules
- **Acceptance Criteria**:
  - Support multichannel intake: web portal, mobile app, email, offline forms, walk-in
  - Capture all mandatory fields: complainant name, contact info, role, language preference, identification (Aadhaar/PAN/Passport)
  - Record policy/claim number, agent details, type of policy (PLI/RPLI)
  - Capture incident dates, representation to insurer date, issue description, relief sought
  - Validate admissibility: representation to insurer first, 30-day wait period, 1-year limitation, claim value ≤ ₹50 lakh, no parallel litigation
  - Auto-generate unique complaint ID on successful registration
  - Issue acknowledgement to complainant within 24 hours via SMS/email
  - Support bilingual interface (English/Hindi) with regional language adaptability
  - Store all attachments: policy documents, correspondence, denial letters, ID proof, receipts, bills (PDF, JPG, PNG, max 10 MB per file)
  - Log all registration activities in audit trail with timestamp, user ID, IP address
- **Source**: Claim_SRS on insurance ombudsman.md, Section 3 - Functional Scope, Lines 62-64; Section 6 - Data Inputs, Lines 176-198
- **Business Rule**: BR-CLM-OMB-001
- **Workflow**: WF-CLM-OMB-INTAKE-001
- **Data Entity**: E-CLM-OMB-COMPLAINT, E-CLM-OMB-COMPLAINANT, E-CLM-OMB-ATTACHMENT
- **Integration**: INT-CLM-OMB-CRM (Customer profile sync), INT-CLM-OMB-CPGRAMS (Grievance portal)
- **Validation**: VR-CLM-OMB-001 to VR-CLM-OMB-005

#### FR-CLM-OMB-002: Jurisdiction Mapping
- **ID**: FR-CLM-OMB-002
- **Category**: Insurance Ombudsman - Jurisdiction
- **Priority**: HIGH
- **Description**: System must dynamically map complaints to territorial ombudsman centers based on complainant location, agent location, or policy office as per Rule 11
- **Acceptance Criteria**:
  - Auto-map complaints using complainant pincode/digital pin code
  - Support alternate mapping based on agent details or policy office location
  - Maintain jurisdiction master table with ombudsman center assignments
  - Screen for conflict of interest: prior relationship with complainant, financial interest, duplicate litigation
  - Auto-reassign to alternate ombudsman if conflict detected
  - Display assigned ombudsman center to complainant upon registration
  - Support manual override by authorized personnel with justification
  - Log all jurisdiction mapping decisions and conflicts in audit trail
  - Update complaint status to "ASSIGNED_TO_JURISDICTION" upon successful mapping
- **Source**: Claim_SRS on insurance ombudsman.md, Section 3 - Functional Scope, Lines 66-68; Section 6 - Data Inputs, Lines 200-206; Section 5 - Lifecycle Flow, Lines 119, 124-127
- **Business Rule**: BR-CLM-OMB-002, BR-CLM-OMB-003
- **Workflow**: WF-CLM-OMB-JURISDICTION-001
- **Data Entity**: E-CLM-OMB-JURISDICTION-MASTER, E-CLM-OMB-CONFLICT-CHECK
- **Integration**: INT-CLM-OMB-POLICY-SERVICE (Policy office mapping)

#### FR-CLM-OMB-003: Hearing Scheduling & Management
- **ID**: FR-CLM-OMB-003
- **Category**: Insurance Ombudsman - Hearing
- **Priority**: HIGH
- **Description**: System must provide end-to-end management of ombudsman hearings including video hearings, calendar management, automated notifications, and conflict detection
- **Acceptance Criteria**:
  - Support hearing mode selection: physical or video conference
  - Provide calendar interface for scheduling with parties' availability tracking
  - Auto-detect scheduling conflicts (ombudsman availability, venue conflicts)
  - Send automated notifications to all parties: complainant, insurer, ombudsman staff via SMS/email/in-app
  - Generate hearing invitation with details: date, time, mode, venue/video link, required documents
  - Support rescheduling with auto-notifications and reason capture
  - Enable document submission before hearing with secure upload
  - Track hearing attendance and record minutes/notes
  - Capture mediation consent from both parties during hearing
  - Record mediation outcome: settled or unsettled
  - Support hearing postponement requests with approval workflow
  - Log all hearing activities in audit trail
- **Source**: Claim_SRS on insurance ombudsman.md, Section 3 - Functional Scope, Lines 73-75; Section 6 - Hearing and Mediation Fields, Lines 230-237; Section 5 - Lifecycle Flow, Lines 131-133
- **Business Rule**: BR-CLM-OMB-004
- **Workflow**: WF-CLM-OMB-HEARING-001
- **Data Entity**: E-CLM-OMB-HEARING, E-CLM-OMB-HEARING-CALENDAR, E-CLM-OMB-MEDIATION
- **Integration**: INT-CLM-OMB-VIDEO-CONF (Video conferencing platform), INT-CLM-OMB-NOTIFICATION (SMS/Email gateway)

#### FR-CLM-OMB-004: Award Issuance & Enforcement
- **ID**: FR-CLM-OMB-004
- **Category**: Insurance Ombudsman - Award
- **Priority**: CRITICAL
- **Description**: System must support workflow for draft, review, approval, digital signing, and communication of mediation recommendations and final awards with regulatory caps enforcement
- **Acceptance Criteria**:
  - Support award type selection: mediation recommendation (Rule 16) or adjudication award (Rule 17)
  - Provide award drafting interface with compensation calculation support
  - Enforce ₹50 lakh cap on award amount with validation
  - Calculate interest as applicable per IRDAI guidelines
  - Support multi-level review and approval workflow
  - Integrate digital signature for ombudsman approval
  - Auto-generate award document with reasons/justification
  - Issue digitally signed award to complainant and insurer
  - Set compliance deadline: 30 days from award issuance
  - Send automated reminders at Day 15, Day 7, Day 2 before deadline
  - Track insurer compliance: acceptance, payment, or objection
  - Auto-escalate to IRDAI on non-compliance
  - Update complaint status: "AWARD_ISSUED" → "COMPLIANCE_PENDING" → "CLOSED" or "ESCALATED"
  - Log all award activities in audit trail
- **Source**: Claim_SRS on insurance ombudsman.md, Section 3 - Functional Scope, Lines 77-79; Section 6 - Award & Resolution Fields, Lines 240-250; Section 5 - Lifecycle Flow, Lines 136-147
- **Business Rule**: BR-CLM-OMB-005, BR-CLM-OMB-006
- **Workflow**: WF-CLM-OMB-AWARD-001, WF-CLM-OMB-COMPLIANCE-001
- **Data Entity**: E-CLM-OMB-AWARD, E-CLM-OMB-COMPLIANCE-TRACKER
- **Integration**: INT-CLM-OMB-DIGITAL-SIGN (Digital signature service), INT-CLM-OMB-IRDAI (IRDAI reporting)

---

### 3.6 Policy Bond & Freelook Requirements

#### FR-CLM-BOND-001: Dispatch Tracking
- **ID**: FR-CLM-BOND-001
- **Category**: Policy Bond - Dispatch Tracking
- **Priority**: CRITICAL
- **Description**: System must fully automate tracking of policy bond dispatches using India Post APIs for both physical and electronic (ePLI) bonds
- **Acceptance Criteria**:
  - Auto-generate and assign India Post tracking numbers (SP article number) upon bond dispatch
  - Trigger dispatch via CEPT Booking Interface integration
  - Generate unique dispatch ID and link to policy number
  - Support both physical dispatch (Speed Post) and digital delivery (DigiLocker/ePLI)
  - Fetch real-time delivery status from India Post API using tracking number
  - Record delivery events: dispatched, in-transit, out-for-delivery, delivered, undelivered
  - Capture delivery confirmation: date, time, recipient signature, OTP verification, photo evidence
  - For ePLI bonds: track issuance, DigiLocker upload, download timestamp, digital acknowledgment
  - Flag failed/delayed deliveries: undelivered after 10 days threshold
  - Escalate to CPC supervisor for manual intervention on delivery failure
  - Notify policyholder on dispatch and successful delivery via SMS/email
  - Sync delivery status with CRM customer profile
  - Log all tracking queries and status updates in audit trail
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Section 4.1 - Dispatch Tracking, Lines 126-165; Flowchart Lines 20-37
- **Business Rule**: BR-CLM-BOND-002
- **Workflow**: WF-CLM-BOND-DISPATCH-001
- **Data Entity**: E-CLM-BOND-DISPATCH, E-CLM-BOND-DELIVERY-STATUS, E-CLM-BOND-POD
- **Integration**: INT-CLM-BOND-INDIA-POST (India Post Tracking API), INT-CLM-BOND-CEPT (CEPT Booking Interface), INT-CLM-BOND-DIGILOCKER (DigiLocker API), INT-CLM-BOND-CRM (Customer profile sync)
- **Validation**: VR-CLM-BOND-001, VR-CLM-BOND-002

#### FR-CLM-BOND-002: Delivery Confirmation
- **ID**: FR-CLM-BOND-002
- **Category**: Policy Bond - Delivery Confirmation
- **Priority**: CRITICAL
- **Description**: System must capture and verify delivery confirmation from recipient for both physical and digital bonds to establish legally defensible proof of delivery
- **Acceptance Criteria**:
  - For physical bonds: capture India Post Proof of Delivery (POD) with date, time, recipient acknowledgment
  - Record digital signature/photo/OTP verification as per India Post delivery practices
  - For ePLI bonds: record download timestamp from DigiLocker with authentication details
  - Validate recipient identity against policy holder/nominee records
  - Store proof of delivery as legally admissible evidence
  - Auto-trigger freelook period timer on confirmed delivery
  - Send delivery confirmation notification to policyholder
  - Handle disputed deliveries: misdelivery, wrong address, non-receipt claims
  - Support delivery dispute workflow: pause timer, re-dispatch bond, fresh delivery confirmation
  - Provide customer self-service view to check delivery status
  - Log all delivery confirmation events in audit trail with full transaction details
  - Generate delivery confirmation certificate on demand
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Section 4.2 - Delivery Confirmation, Lines 167-201; Illustrative Example Lines 180-201
- **Business Rule**: BR-CLM-BOND-001
- **Workflow**: WF-CLM-BOND-DELIVERY-001, WF-CLM-BOND-DISPUTE-001
- **Data Entity**: E-CLM-BOND-POD, E-CLM-BOND-DELIVERY-CONFIRMATION, E-CLM-BOND-DISPUTE
- **Integration**: INT-CLM-BOND-INDIA-POST (POD retrieval), INT-CLM-BOND-DIGILOCKER (Digital confirmation)
- **Validation**: VR-CLM-BOND-003, VR-CLM-BOND-004

#### FR-CLM-FL-001: Freelook Period Monitoring
- **ID**: FR-CLM-FL-001
- **Category**: Freelook - Period Monitoring
- **Priority**: CRITICAL
- **Description**: System must automatically determine and monitor freelook window (15 days for physical, 30 days for online/ePLI) from confirmed delivery date with dynamic countdown
- **Acceptance Criteria**:
  - Auto-calculate freelook start date = confirmed delivery date (physical) or issuance date (ePLI)
  - Calculate freelook end date: delivery_date + 15 days (physical) or issuance_date + 30 days (ePLI)
  - Provide dynamic countdown: calculate remaining days in freelook window
  - Display freelook status to customer: days remaining, expiry date
  - Send proactive reminders: Day 7 and Day 12 for customer awareness
  - Auto-reject cancellation requests submitted after window expiry
  - Support exception handling for disputed deliveries: pause timer, restart on fresh delivery
  - Provide customer self-service views: SMS query, online portal check
  - Handle multiple delivery attempts: reset timer on successful re-delivery
  - Log all timer calculations and adjustments in audit trail
  - Generate freelook period compliance reports
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Section 4.3 - Freelook Period Monitoring, Lines 205-231; Illustrative Example Lines 222-231
- **Business Rule**: BR-CLM-BOND-001
- **Workflow**: WF-CLM-FL-TIMER-001
- **Data Entity**: E-CLM-FL-PERIOD, E-CLM-FL-TIMER, E-CLM-FL-REMINDER
- **Integration**: INT-CLM-FL-NOTIFICATION (SMS/Email alerts), INT-CLM-FL-PORTAL (Customer self-service)
- **Validation**: VR-CLM-FL-001, VR-CLM-FL-002

#### FR-CLM-FL-002: Cancellation Request Handling
- **ID**: FR-CLM-FL-002
- **Category**: Freelook - Cancellation Request
- **Priority**: CRITICAL
- **Description**: System must accept, validate, and process freelook cancellation requests via multiple channels with complete documentation verification
- **Acceptance Criteria**:
  - Accept cancellation requests via: online portal, CRM, Post Office counter, CPGRAMS, email
  - Validate cancellation timestamp is within freelook window using delivery date reference
  - Mandatory document validation: cancellation request letter/form, policyholder ID proof, original policy bond (physical or ePLI), proof of delivery/receipt, KYC documents
  - Support authorized messenger submissions with proper documentation (medical certificate for unfitness)
  - Detect and flag tampered documents: ID proof alteration, fraudulent signatures
  - Auto-generate unique cancellation request ID
  - Issue system-generated acknowledgement to applicant within 24 hours
  - Route flagged cases for enhanced review and verification
  - Track cancellation workflow status: Submitted → Under Review → Document Verified → Approved/Rejected
  - Send progress notifications at each status change
  - Capture rejection reasons if submitted post-window or incomplete documentation
  - Log all submission, review, and approval events in audit trail
  - Update policy status to "FREELOOK_CANCELLATION_PENDING" on acceptance
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Section 4.4 - Cancellation Request Handling, Lines 233-261; Illustrative Example Lines 254-261
- **Business Rule**: BR-CLM-BOND-001
- **Workflow**: WF-CLM-FL-CANCELLATION-001
- **Data Entity**: E-CLM-FL-CANCELLATION-REQUEST, E-CLM-FL-DOCUMENT, E-CLM-FL-VERIFICATION
- **Integration**: INT-CLM-FL-CPGRAMS (Grievance portal), INT-CLM-FL-CRM (Customer interaction), INT-CLM-FL-KYC (Identity verification)
- **Validation**: VR-CLM-FL-003, VR-CLM-FL-004, VR-CLM-FL-005

#### FR-CLM-FL-003: Refund Processing
- **ID**: FR-CLM-FL-003
- **Category**: Freelook - Refund Processing
- **Priority**: CRITICAL
- **Description**: System must calculate and process refunds after deductions with maker-checker workflow and integration with accounts module for disbursement
- **Acceptance Criteria**:
  - Calculate refund amount: Premium - (Pro-rata risk premium + Stamp duty + Medical costs + Other deductions per POLI rules)
  - Display detailed refund calculation breakdown to applicant
  - Support multiple refund modes: NEFT, POSB account credit, crossed cheque
  - Validate refund account details: bank account number, IFSC code, POSB account, cancelled cheque upload
  - Implement maker-checker workflow: maker enters refund details, checker verifies and approves
  - Enforce segregation of duties: maker_id != checker_id validation
  - Generate unique transaction ID for each refund
  - Create sanction letter and payment voucher for accounts module
  - Integrate with PLI accounts module for transaction posting and disbursement
  - Update claim/payment register with financial reconciliation data
  - Track refund status: Calculated → Maker Entry → Checker Verification → Accounts Approved → Payment Processed
  - Send refund confirmation to policyholder with transaction details
  - Update policy status to "CANCELLED_FREELOOK" on successful refund
  - Notify all stakeholders: policyholder, circle office, accounts team, agent (for commission recovery if applicable)
  - Log all refund calculation, approval, and payment steps in audit trail
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Section 4.5 - Refund Processing, Lines 263-292; Illustrative Example Lines 284-291
- **Business Rule**: BR-CLM-BOND-003, BR-CLM-BOND-004
- **Workflow**: WF-CLM-FL-REFUND-001, WF-CLM-FL-MAKER-CHECKER-001
- **Data Entity**: E-CLM-FL-REFUND-CALCULATION, E-CLM-FL-PAYMENT-VOUCHER, E-CLM-FL-TRANSACTION
- **Integration**: INT-CLM-FL-ACCOUNTS (PLI Accounts Module), INT-CLM-FL-FINANCE (Finance System reconciliation), INT-CLM-FL-PAYMENT-GATEWAY (NEFT/Bank transfer)
- **Validation**: VR-CLM-FL-006, VR-CLM-FL-007, VR-CLM-FL-008

---

## 4. Validation Rules

### 4.1 Death Claim Validations

#### VR-CLM-DC-001: Death Certificate Mandatory
- **ID**: VR-CLM-DC-001
- **Field**: Death Certificate
- **Type**: Mandatory Document
- **Rule**: Death certificate must be uploaded before claim processing
- **Error Message**: "Death certificate is mandatory for claim processing"
- **Source**: Claim_SRS FRS on death claim.md, Section 2
- **Priority**: CRITICAL

#### VR-CLM-DC-002: Claim Form Mandatory
- **ID**: VR-CLM-DC-002
- **Field**: Claim Form
- **Type**: Mandatory Document
- **Rule**: Duly filled claim form required
- **Error Message**: "Claim form must be completed and submitted"
- **Source**: Claim_SRS FRS on death claim.md, Line 62
- **Priority**: CRITICAL

#### VR-CLM-DC-003: Policy Bond or Indemnity
- **ID**: VR-CLM-DC-003
- **Field**: Policy Bond/Indemnity
- **Type**: Mandatory Document
- **Rule**: Original policy bond OR indemnity bond required
- **Error Message**: "Policy bond or indemnity bond is mandatory"
- **Source**: Claim_SRS FRS on death claim.md, Lines 62-63
- **Priority**: CRITICAL

#### VR-CLM-DC-004: Claimant ID Proof
- **ID**: VR-CLM-DC-004
- **Field**: Claimant ID Proof
- **Type**: Mandatory Document
- **Rule**: Valid government-issued ID required
- **Error Message**: "Valid ID proof (Aadhaar/PAN/Passport) required"
- **Source**: Claim_SRS FRS on death claim.md, Line 63
- **Priority**: CRITICAL

#### VR-CLM-DC-005: Bank Mandate
- **ID**: VR-CLM-DC-005
- **Field**: Bank Details
- **Type**: Mandatory Document
- **Rule**: Bank account details with cancelled cheque/passbook copy
- **Error Message**: "Bank mandate with cancelled cheque/passbook copy required"
- **Source**: Claim_SRS FRS on death claim.md, Line 63
- **Priority**: CRITICAL

#### VR-CLM-DC-006: Unnatural Death Documents
- **ID**: VR-CLM-DC-006
- **Field**: FIR and Postmortem Report
- **Type**: Conditional Mandatory
- **Rule**: IF death_type = "UNNATURAL" THEN FIR AND postmortem_report REQUIRED
- **Error Message**: "FIR and postmortem report mandatory for unnatural deaths"
- **Source**: Claim_SRS FRS on death claim.md, Line 64
- **Priority**: CRITICAL

#### VR-CLM-DC-007: Nomination Documents
- **ID**: VR-CLM-DC-007
- **Field**: Succession Certificate / Legal Heir Affidavit
- **Type**: Conditional Mandatory
- **Rule**: IF nomination_status = "ABSENT" THEN succession_certificate OR legal_heir_affidavit REQUIRED
- **Error Message**: "Succession certificate or legal heir affidavit required when nomination is absent"
- **Source**: Claim_SRS FRS on death claim.md, Lines 66-67
- **Priority**: CRITICAL

#### VR-CLM-DC-008: Investigation Trigger Check
- **ID**: VR-CLM-DC-008
- **Field**: Investigation Trigger
- **Type**: Business Logic
- **Rule**: IF death_date <= (policy_acceptance_date + 3 YEARS) OR death_date <= (policy_revival_date + 3 YEARS) THEN investigation_required = TRUE
- **Error Message**: "Investigation required - death occurred within 3 years of policy acceptance or revival"
- **Source**: Claim_SRS FRS on death claim.md, Lines 75-77
- **Priority**: HIGH

#### VR-CLM-DC-009: Investigation Timeline Validation
- **ID**: VR-CLM-DC-009
- **Field**: Investigation Report Submission Date
- **Type**: Timeline Validation
- **Rule**: investigation_report_date <= (investigation_assigned_date + 21 DAYS)
- **Error Message**: "Investigation report must be submitted within 21 days"
- **Source**: Claim_SRS FRS on death claim.md, Lines 81-82
- **Priority**: HIGH

#### VR-CLM-DC-010: Approval Timeline - No Investigation
- **ID**: VR-CLM-DC-010
- **Field**: Approval Date
- **Type**: SLA Validation
- **Rule**: IF investigation_required = FALSE THEN approval_date <= (claim_submission_date + 15 DAYS)
- **Error Message**: "Death claim approval must be completed within 15 days when no investigation is required"
- **Source**: Claim_SRS FRS on death claim.md, Lines 107-109
- **Priority**: CRITICAL

#### VR-CLM-DC-011: Approval Timeline - With Investigation
- **ID**: VR-CLM-DC-011
- **Field**: Approval Date
- **Type**: SLA Validation
- **Rule**: IF investigation_required = TRUE THEN approval_date <= (claim_submission_date + 45 DAYS)
- **Error Message**: "Death claim approval must be completed within 45 days for claims requiring investigation"
- **Source**: Claim_SRS FRS on death claim.md, Lines 109-111
- **Priority**: CRITICAL

#### VR-CLM-DC-012: Document Retention Period
- **ID**: VR-CLM-DC-012
- **Field**: Document Submission Date
- **Type**: Timeline Validation
- **Rule**: IF documents_not_received_within = (request_date + 15 DAYS + 7 DAYS GRACE) THEN return_claim = TRUE
- **Error Message**: "Documents not received within 15 days plus 7-day grace period - claim returned"
- **Source**: Claim_SRS FRS on death claim.md, Lines 68-71
- **Priority**: MEDIUM

#### VR-CLM-DC-013: Investigation Outcome Validation
- **ID**: VR-CLM-DC-013
- **Field**: Investigation Outcome
- **Type**: Status Validation
- **Rule**: investigation_outcome IN ("Clear", "Suspect", "Fraud")
- **Error Message**: "Invalid investigation outcome - must be Clear, Suspect, or Fraud"
- **Source**: Claim_SRS FRS on death claim.md, Lines 84-86
- **Priority**: HIGH

#### VR-CLM-DC-014: Penal Interest Auto-Calculation
- **ID**: VR-CLM-DC-014
- **Field**: Penal Interest
- **Type**: Calculation
- **Rule**: IF approval_date > SLA_deadline THEN penal_interest = claim_amount * 0.08 * (days_delayed / 365)
- **Error Message**: "Penal interest calculated at 8% p.a. for SLA breach"
- **Source**: Claim_SRS FRS on death claim.md, Lines 377-378
- **Priority**: HIGH

#### VR-CLM-DC-015: Reopening Request Validation
- **ID**: VR-CLM-DC-015
- **Field**: Reopening Reason
- **Type**: Business Logic
- **Rule**: reopening_reason IN ("Court Order", "New Evidence", "Administrative Lapse", "Claimant Appeal")
- **Error Message**: "Invalid reopening reason - must be Court Order, New Evidence, Administrative Lapse, or Claimant Appeal"
- **Source**: Claim_SRS FRS on death claim.md, Lines 125-128
- **Priority**: MEDIUM

#### VR-CLM-DC-016: Appeal Timeline Validation
- **ID**: VR-CLM-DC-016
- **Field**: Appeal Filing Date
- **Type**: Timeline Validation
- **Rule**: appeal_filing_date <= (claim_rejection_date + 90 DAYS)
- **Error Message**: "Appeal must be filed within 90 days of claim rejection"
- **Source**: Claim_SRS FRS on death claim.md, Lines 143-144
- **Priority**: HIGH

#### VR-CLM-DC-017: Appeal Decision Timeline
- **ID**: VR-CLM-DC-017
- **Field**: Appeal Decision Date
- **Type**: Timeline Validation
- **Rule**: appeal_decision_date <= (appeal_filing_date + 45 DAYS)
- **Error Message**: "Appeal decision must be issued within 45 days of filing"
- **Source**: Claim_SRS FRS on death claim.md, Lines 147-148
- **Priority**: HIGH

#### VR-CLM-DC-018: Payment Mode Validation
- **ID**: VR-CLM-DC-018
- **Field**: Payment Mode
- **Type**: Enumeration Validation
- **Rule**: payment_mode IN ("NEFT", "POSB_EFT", "CHEQUE")
- **Error Message**: "Invalid payment mode - must be NEFT, POSB EFT, or Cheque"
- **Source**: Claim_SRS FRS on death claim.md, Lines 116-118
- **Priority**: MEDIUM

#### VR-CLM-DC-019: Multiple Claim Form Types
- **ID**: VR-CLM-DC-019
- **Field**: Claim Form Type
- **Type**: Document Validation
- **Rule**: Validate claim form type matches claim category (death, maturity, surrender, etc.)
- **Error Message**: "Claim form type does not match claim category"
- **Source**: Claim_SRS FRS on death claim.md, Line 62
- **Priority**: MEDIUM

#### VR-CLM-DC-020: MWPA/HUF Policy Type Check
- **ID**: VR-CLM-DC-020
- **Field**: Policy Type
- **Type**: Business Logic
- **Rule**: IF policy_type IN ("MWPA", "HUF") THEN require_additional_documents = TRUE
- **Error Message**: "MWPA/HUF policy type requires additional documentation"
- **Source**: Claim_SRS FRS on death claim.md, Lines 367-368
- **Priority**: MEDIUM

#### VR-CLM-DC-021: Manual Override Audit Requirement
- **ID**: VR-CLM-DC-021
- **Field**: Manual Override
- **Type**: Audit Trail
- **Rule**: IF manual_override = TRUE THEN log_justification AND digital_signature REQUIRED
- **Error Message**: "Manual intervention requires justification and digital signature"
- **Source**: Claim_SRS FRS on death claim.md, Lines 279-281
- **Priority**: HIGH


### 4.2 Maturity Claim Validations

#### VR-CLM-MC-001: Policy Status Active
- **ID**: VR-CLM-MC-001
- **Field**: Policy Status
- **Type**: Business Logic
- **Rule**: Policy must be active on maturity date
- **Error Code**: ERR-CLM-MC-RJ-P-02
- **Error Message**: "Policy not active on maturity date - claim cannot be processed"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-02
- **Priority**: CRITICAL

#### VR-CLM-MC-002: Duplicate Claim Check
- **ID**: VR-CLM-MC-002
- **Field**: Policy Number
- **Type**: Business Logic
- **Rule**: System must check if maturity claim already paid for this policy
- **Error Code**: ERR-CLM-MC-RJ-P-03
- **Error Message**: "Maturity claim already paid for this policy"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-03, Line 652
- **Priority**: CRITICAL

#### VR-CLM-MC-003: Claimant Identity Match
- **ID**: VR-CLM-MC-003
- **Field**: Claimant Details
- **Type**: Business Logic
- **Rule**: Claimant details must match policy records
- **Error Code**: ERR-CLM-MC-RJ-E-01
- **Error Message**: "Claimant details do not match policy records"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-E-01, Line 663
- **Priority**: CRITICAL

#### VR-CLM-MC-004: Bank Account Validation
- **ID**: VR-CLM-MC-004
- **Field**: Bank Account Details
- **Type**: API Validation
- **Rule**: Bank account must be validated via CBS/PFMS API
- **Error Code**: ERR-CLM-MC-RJ-B-01
- **Error Message**: "Bank account details invalid or verification failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-10
- **Priority**: CRITICAL

#### VR-CLM-MC-005: IFSC Code Format
- **ID**: VR-CLM-MC-005
- **Field**: IFSC Code
- **Type**: Format Validation
- **Rule**: IFSC must be 11 characters (4 letters + 0 + 6 alphanumeric)
- **Error Code**: ERR-CLM-MC-RJ-B-02
- **Error Message**: "Invalid IFSC code format"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-B-02, Line 700
- **Priority**: HIGH

#### VR-CLM-MC-006: Policy Bond Submission
- **ID**: VR-CLM-MC-006
- **Field**: Policy Bond
- **Type**: Mandatory Document
- **Rule**: Original policy bond must be submitted
- **Error Code**: ERR-CLM-MC-RJ-D-04
- **Error Message**: "Policy bond not submitted or invalid"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-D-04, Line 688
- **Priority**: CRITICAL

#### VR-CLM-MC-007: ID Proof Validity
- **ID**: VR-CLM-MC-007
- **Field**: Identity Proof
- **Type**: Document Validation
- **Rule**: ID proof must be valid and not expired
- **Error Code**: ERR-CLM-MC-RJ-D-05
- **Error Message**: "Identity proof/address proof invalid or expired"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-D-05, Line 690
- **Priority**: CRITICAL

#### VR-CLM-MC-008: Document Forgery Check
- **ID**: VR-CLM-MC-008
- **Field**: All Documents
- **Type**: Security Validation
- **Rule**: System must flag suspected forged or tampered documents
- **Error Code**: ERR-CLM-MC-RJ-D-02
- **Error Message**: "Submitted documents found forged or suspicious - flagged for investigation"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-D-02, Line 683
- **Priority**: CRITICAL

#### VR-CLM-MC-009: Policy Number Format Validation
- **ID**: VR-CLM-MC-009
- **Field**: Policy Number
- **Type**: Format Validation
- **Rule**: Policy number format must be valid and exist in system
- **Error Code**: ERR-CLM-MC-RJ-P-01
- **Error Message**: "Policy number is invalid or does not exist"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-01, Line 648
- **Priority**: CRITICAL

#### VR-CLM-MC-010: Policy Forfeiture/Surrender Check
- **ID**: VR-CLM-MC-010
- **Field**: Policy Status
- **Type**: Business Logic
- **Rule**: Policy must not be forfeited or surrendered prior to maturity
- **Error Code**: ERR-CLM-MC-RJ-P-04
- **Error Message**: "Policy terminated due to forfeiture/surrender prior to maturity"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-04, Lines 654-655
- **Priority**: CRITICAL

#### VR-CLM-MC-011: Unauthorized Claimant Check
- **ID**: VR-CLM-MC-011
- **Field**: Claimant Authorization
- **Type**: Business Logic
- **Rule**: Claim must be submitted by authorized person (policyholder or registered nominee)
- **Error Code**: ERR-CLM-MC-RJ-E-03
- **Error Message**: "Claim submitted by unauthorized person"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-E-03, Line 668
- **Priority**: CRITICAL

#### VR-CLM-MC-012: Nominee/Legal Heir Validity
- **ID**: VR-CLM-MC-012
- **Field**: Nominee/Legal Heir Details
- **Type**: Document Validation
- **Rule**: Nominee or legal heir details must be valid and match policy records
- **Error Code**: ERR-CLM-MC-RJ-E-04
- **Error Message**: "Nominee/legal heir details not valid"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-E-04, Line 670
- **Priority**: HIGH

#### VR-CLM-MC-013: Multiple Claimants Dispute Check
- **ID**: VR-CLM-MC-013
- **Field**: Claimant Count
- **Type**: Business Logic
- **Rule**: IF multiple_claimants = TRUE AND entitlement_dispute = UNRESOLVED THEN block_processing = TRUE
- **Error Code**: ERR-CLM-MC-RJ-E-05
- **Error Message**: "Multiple claimants with unresolved entitlement dispute"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-E-05, Lines 672-673
- **Priority**: MEDIUM

#### VR-CLM-MC-014: Mandatory Documents Completeness
- **ID**: VR-CLM-MC-014
- **Field**: Document Checklist
- **Type**: Mandatory Document
- **Rule**: All mandatory documents from checklist must be submitted
- **Error Code**: ERR-CLM-MC-RJ-D-01
- **Error Message**: "Mandatory documents not submitted"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-D-01, Line 681
- **Priority**: CRITICAL

#### VR-CLM-MC-015: Physical vs Digital Document Mismatch
- **ID**: VR-CLM-MC-015
- **Field**: Document Comparison
- **Type**: Document Validation
- **Rule**: Physical documents must match digital records
- **Error Code**: ERR-CLM-MC-RJ-D-03
- **Error Message**: "Mismatch between physical documents and digital records"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-D-03, Lines 685-686
- **Priority**: HIGH

#### VR-CLM-MC-016: Repeated Payment Failure Check
- **ID**: VR-CLM-MC-016
- **Field**: Payment Failure Count
- **Type**: Business Logic
- **Rule**: IF payment_failure_count >= 3 THEN flag_incorrect_bank_details = TRUE
- **Error Code**: ERR-CLM-MC-RJ-B-03
- **Error Message**: "Repeated failure of payment due to incorrect bank details"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-B-03, Lines 702-703
- **Priority**: HIGH

#### VR-CLM-MC-017: Maturity Due Report Generation
- **ID**: VR-CLM-MC-017
- **Field**: Maturity Due Report
- **Type**: System Requirement
- **Rule**: Auto-generate maturity due report daily/weekly for policies reaching maturity
- **Error Message**: "Maturity due report generation failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-01, Lines 440-446
- **Priority**: MEDIUM

#### VR-CLM-MC-018: Multi-Channel Intimation Validation
- **ID**: VR-CLM-MC-018
- **Field**: Intimation Channels
- **Type**: System Requirement
- **Rule**: Send intimation via SMS, Email, WhatsApp, and Portal with secure link
- **Error Message**: "Multi-channel intimation delivery failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-02, Lines 449-459
- **Priority**: MEDIUM

#### VR-CLM-MC-019: DigiLocker Integration Validation
- **ID**: VR-CLM-MC-019
- **Field**: DigiLocker Document
- **Type**: API Integration
- **Rule**: Integrate with DigiLocker for fetching policy documents
- **Error Message**: "DigiLocker document fetch failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-03, Lines 467-469
- **Priority**: MEDIUM

#### VR-CLM-MC-020: Auto-Acknowledgment Generation
- **ID**: VR-CLM-MC-020
- **Field**: Claim Acknowledgment
- **Type**: System Requirement
- **Rule**: Auto-generate acknowledgment with unique Claim ID and submission timestamp
- **Error Message**: "Acknowledgment generation failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-03, Lines 470-475
- **Priority**: MEDIUM

#### VR-CLM-MC-021: Document Checklist Validation
- **ID**: VR-CLM-MC-021
- **Field**: Document Checklist
- **Type**: Document Validation
- **Rule**: Validate uploaded documents against checklist and flag missing items
- **Error Message**: "Document checklist validation incomplete - missing items flagged"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-04, Lines 477-489
- **Priority**: HIGH

#### VR-CLM-MC-022: Auto-Reminder for Missing Documents
- **ID**: VR-CLM-MC-022
- **Field**: Missing Documents
- **Type**: System Requirement
- **Rule**: Send auto-reminders for outstanding documents
- **Error Message**: "Auto-reminder system unavailable"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-04, Lines 481-483
- **Priority**: MEDIUM

#### VR-CLM-MC-023: OCR Data Extraction Validation
- **ID**: VR-CLM-MC-023
- **Field**: OCR Extracted Data
- **Type**: Data Validation
- **Rule**: Auto-populate claim data fields from uploaded documents using OCR
- **Error Message**: "OCR data extraction failed or incomplete"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-06, Lines 499-506
- **Priority**: MEDIUM

#### VR-CLM-MC-024: QC Dual Authentication
- **ID**: VR-CLM-MC-024
- **Field**: QC Override
- **Type**: Authorization
- **Rule**: Dual authentication required for QC override or waiver
- **Error Message**: "Dual authentication required for override/waiver"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-07, Lines 511-516
- **Priority**: HIGH

#### VR-CLM-MC-025: SLA 7-Day Countdown
- **ID**: VR-CLM-MC-025
- **Field**: SLA Timer
- **Type**: SLA Validation
- **Rule**: Display SLA countdown timer (7-day window) for maturity claim processing
- **Error Message**: "SLA deadline approaching/breached - immediate action required"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-08, Lines 524-528
- **Priority**: CRITICAL

#### VR-CLM-MC-026: Voucher Auto-Generation
- **ID**: VR-CLM-MC-026
- **Field**: Payment Voucher
- **Type**: System Requirement
- **Rule**: Auto-generate voucher post-disbursement with payment details
- **Error Message**: "Voucher auto-generation failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, FRS-MAT-12, Lines 556-562
- **Priority**: MEDIUM

### 4.3 Survival Benefit Validations

#### VR-CLM-SB-001: Policy Active on SB Due Date
- **ID**: VR-CLM-SB-001
- **Field**: Policy Status
- **Type**: Business Logic
- **Rule**: Policy must be active on survival benefit due date
- **Error Code**: ERR-CLM-SB-RJ-P-02
- **Error Message**: "Policy inactive on Survival Benefit due date"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-02, Line 447
- **Priority**: CRITICAL

#### VR-CLM-SB-002: SB Already Paid Check
- **ID**: VR-CLM-SB-002
- **Field**: Policy Number
- **Type**: Business Logic
- **Rule**: Check if this survival benefit already paid
- **Error Code**: ERR-CLM-SB-RJ-P-03
- **Error Message**: "Survival Benefit already paid"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-P-03, Line 450
- **Priority**: CRITICAL

#### VR-CLM-SB-003: Policyholder Identity Match
- **ID**: VR-CLM-SB-003
- **Field**: Policyholder Details
- **Type**: Business Logic
- **Rule**: Policyholder identity must match policy records
- **Error Code**: ERR-CLM-SB-RJ-E-02
- **Error Message**: "Policyholder identity mismatch"
- **Source**: Claim_SRS FRS on Maturity claim.md, Rejection Reason RJ-E-02, Line 461
- **Priority**: CRITICAL

#### VR-CLM-SB-004: Policy Number Invalid Check
- **ID**: VR-CLM-SB-004
- **Field**: Policy Number
- **Type**: Format Validation
- **Rule**: Policy number must be valid and found in system
- **Error Code**: ERR-CLM-SB-RJ-P-01
- **Error Message**: "Policy number invalid or not found"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, Rejection Reason RJ-P-01, Line 445
- **Priority**: CRITICAL

#### VR-CLM-SB-005: Mandatory Documents Check
- **ID**: VR-CLM-SB-005
- **Field**: Document Checklist
- **Type**: Mandatory Document
- **Rule**: All mandatory documents must be submitted
- **Error Code**: ERR-CLM-SB-RJ-D-01
- **Error Message**: "Mandatory documents not submitted"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, Rejection Reason RJ-D-01, Line 471
- **Priority**: CRITICAL

#### VR-CLM-SB-006: Forged Documents Detection
- **ID**: VR-CLM-SB-006
- **Field**: Submitted Documents
- **Type**: Security Validation
- **Rule**: System must detect and flag suspected forged or fraudulent documents
- **Error Code**: ERR-CLM-SB-RJ-D-02
- **Error Message**: "Suspected forged or fraudulent documents"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, Rejection Reason RJ-D-02, Lines 473-474
- **Priority**: CRITICAL

#### VR-CLM-SB-007: Physical vs Digital Mismatch
- **ID**: VR-CLM-SB-007
- **Field**: Document Comparison
- **Type**: Document Validation
- **Rule**: Physical documents must match digital records
- **Error Code**: ERR-CLM-SB-RJ-D-03
- **Error Message**: "Mismatch between physical and digital records"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, Rejection Reason RJ-D-03, Lines 476-477
- **Priority**: HIGH

#### VR-CLM-SB-008: Auto-Generation of SB Due Report
- **ID**: VR-CLM-SB-008
- **Field**: Survival Benefit Due Report
- **Type**: System Requirement
- **Rule**: Auto-generate survival benefit due report daily/weekly
- **Error Message**: "SB due report generation failed"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, FRS-SB-01, Lines 325-329
- **Priority**: MEDIUM

#### VR-CLM-SB-009: Multi-Channel Intimation
- **ID**: VR-CLM-SB-009
- **Field**: Intimation Channels
- **Type**: System Requirement
- **Rule**: Send intimation via SMS, Email, WhatsApp with secure link
- **Error Message**: "Multi-channel intimation delivery failed"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, FRS-SB-02, Lines 331-339
- **Priority**: MEDIUM

#### VR-CLM-SB-010: DigiLocker Document Fetch
- **ID**: VR-CLM-SB-010
- **Field**: DigiLocker Integration
- **Type**: API Integration
- **Rule**: Integrate with DigiLocker for policy document fetch
- **Error Message**: "DigiLocker document fetch failed"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, FRS-SB-03, Lines 344-348
- **Priority**: MEDIUM

#### VR-CLM-SB-011: Auto-Reminder for Pending Docs
- **ID**: VR-CLM-SB-011
- **Field**: Missing Documents
- **Type**: System Requirement
- **Rule**: Send auto-reminders for outstanding documents
- **Error Message**: "Auto-reminder system unavailable"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, FRS-SB-04, Lines 352-356
- **Priority**: MEDIUM

#### VR-CLM-SB-012: SLA 7-Day Enforcement
- **ID**: VR-CLM-SB-012
- **Field**: SLA Timer
- **Type**: SLA Validation
- **Rule**: Enforce SLA of 7 days for survival benefit approval
- **Error Message**: "SLA deadline approaching/breached - immediate action required"
- **Source**: Claim_SRS FRS on Survival Benefit claim.md, FRS-SB-08, Lines 384-385
- **Priority**: CRITICAL

### 4.4 AML Validations

#### VR-CLM-AML-001: Cash Amount Threshold
- **ID**: VR-CLM-AML-001
- **Field**: Cash Amount
- **Type**: Threshold Check
- **Rule**: IF cash_amount > 50000 THEN trigger_alert = TRUE
- **Alert Level**: HIGH
- **Action**: Generate CTR filing
- **Source**: Claim_SRS_AML triggers & alerts.md, AML_001
- **Priority**: CRITICAL

#### VR-CLM-AML-002: PAN Verification
- **ID**: VR-CLM-AML-002
- **Field**: PAN Number
- **Type**: API Validation
- **Rule**: PAN must be verified via NSDL/Income Tax API
- **Alert Level**: MEDIUM (if verification fails)
- **Action**: Flag for manual review
- **Source**: Claim_SRS_AML triggers & alerts.md, AML_002, Line 259
- **Priority**: HIGH

#### VR-CLM-AML-003: Nominee Change Date
- **ID**: VR-CLM-AML-003
- **Field**: Nominee Change Date
- **Type**: Business Logic
- **Rule**: IF nominee_change_date > death_date THEN block_transaction = TRUE
- **Alert Level**: CRITICAL
- **Action**: Block transaction + File STR
- **Source**: Claim_SRS_AML triggers & alerts.md, AML_003
- **Priority**: CRITICAL

#### VR-CLM-AML-004: Surrender Frequency
- **ID**: VR-CLM-AML-004
- **Field**: Surrender Count
- **Type**: Pattern Analysis
- **Rule**: IF surrender_count > 3 within 6 months THEN flag_for_investigation = TRUE
- **Alert Level**: MEDIUM
- **Action**: Investigation trigger
- **Source**: Claim_SRS_AML triggers & alerts.md, AML_004, Lines 267-270
- **Priority**: MEDIUM

#### VR-CLM-AML-005: Refund vs Bond Dispatch Date
- **ID**: VR-CLM-AML-005
- **Field**: Refund Date, Bond Dispatch Date
- **Type**: Business Logic
- **Rule**: IF refund_date < bond_dispatch_date THEN raise_alert = TRUE
- **Alert Level**: HIGH
- **Action**: Audit trail logging
- **Source**: Claim_SRS_AML triggers & alerts.md, AML_005, Lines 273-276
- **Priority**: HIGH

#### VR-CLM-AML-006: Source of Funds Documentation
- **ID**: VR-CLM-AML-006
- **Field**: Source of Funds
- **Type**: Compliance Documentation
- **Rule**: Customers must provide proof of income and disclose net worth for high-value transactions
- **Alert Level**: HIGH
- **Action**: Request income proof and net worth disclosure
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 165-172
- **Priority**: HIGH

#### VR-CLM-AML-007: Cash Payment Limit Enforcement
- **ID**: VR-CLM-AML-007
- **Field**: Cash Payment Amount
- **Type**: Regulatory Compliance
- **Rule**: Enforce regulatory limits on cash acceptance (as per applicable regulations)
- **Alert Level**: CRITICAL
- **Action**: Block transaction if limit exceeded
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 174-181
- **Priority**: CRITICAL

#### VR-CLM-AML-008: Third-Party Payment PAN/KYC
- **ID**: VR-CLM-AML-008
- **Field**: Third-Party Payment
- **Type**: Compliance Validation
- **Rule**: IF payment_source != policyholder THEN require_full_pan_and_kyc = TRUE
- **Alert Level**: CRITICAL
- **Action**: Verify PAN and complete KYC for third-party payer
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 183-190
- **Priority**: CRITICAL

#### VR-CLM-AML-009: STR Filing Requirement
- **ID**: VR-CLM-AML-009
- **Field**: Suspicious Transaction
- **Type**: Regulatory Filing
- **Rule**: File Suspicious Transaction Reports (STR) for transactions suspected of money laundering
- **Alert Level**: CRITICAL
- **Action**: File STR with FIU-IND within prescribed timeline
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 200-203
- **Priority**: CRITICAL

#### VR-CLM-AML-010: CCR Filing for Counterfeit Currency
- **ID**: VR-CLM-AML-010
- **Field**: Currency Validation
- **Type**: Regulatory Filing
- **Rule**: File Counterfeit Currency Reports (CCR) for detection of fake currency
- **Alert Level**: CRITICAL
- **Action**: File CCR immediately upon detection
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 207-208
- **Priority**: CRITICAL

#### VR-CLM-AML-011: NTR for Non-Profit Transactions
- **ID**: VR-CLM-AML-011
- **Field**: Non-Profit Organisation Transaction
- **Type**: Regulatory Filing
- **Rule**: File Non-Profit Organisation Transaction Reports (NTR) for applicable transactions
- **Alert Level**: MEDIUM
- **Action**: File NTR as per regulatory requirements
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 210-211
- **Priority**: MEDIUM

#### VR-CLM-AML-012: Risk-Based Customer Profiling
- **ID**: VR-CLM-AML-012
- **Field**: Customer Risk Profile
- **Type**: Risk Assessment
- **Rule**: Classify customers as high or low risk based on category, occupation, geography, transaction patterns
- **Alert Level**: HIGH (for high-risk customers)
- **Action**: Enhanced due diligence for high-risk customers
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 146-156
- **Priority**: HIGH

#### VR-CLM-AML-013: Daily Negative List Screening
- **ID**: VR-CLM-AML-013
- **Field**: Customer Name/Details
- **Type**: Compliance Screening
- **Rule**: Perform daily screenings against negative lists (terrorism financing, sanctions, PEP)
- **Alert Level**: CRITICAL
- **Action**: Report matches immediately and block transactions
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 125-129
- **Priority**: CRITICAL

#### VR-CLM-AML-014: CTR Threshold Validation (₹10 lakh daily)
- **ID**: VR-CLM-AML-014
- **Field**: Daily Cash Transaction Amount
- **Type**: Threshold Validation
- **Rule**: Trigger CTR filing when daily cash transactions exceed ₹10 lakh
- **Alert Level**: CRITICAL
- **Action**: Auto-generate CTR filing for FIU-IND submission
- **Error Message**: "CTR filing required - daily cash transaction threshold exceeded"
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Lines 121-123
- **Priority**: CRITICAL

### 4.5 Free Look Validations

#### VR-CLM-FL-001: Freelook Window Check
- **ID**: VR-CLM-FL-001
- **Field**: Cancellation Request Date
- **Type**: Date Range Validation
- **Rule**: cancellation_date must be <= (delivery_date + 15 days)
- **Error Message**: "Free look cancellation window expired"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 209-220
- **Priority**: CRITICAL

#### VR-CLM-FL-002: Original Bond Submission
- **ID**: VR-CLM-FL-002
- **Field**: Policy Bond
- **Type**: Mandatory Document
- **Rule**: Original policy bond (physical or digital) must be submitted
- **Error Message**: "Original policy bond required for free look cancellation"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 240-246
- **Priority**: CRITICAL

#### VR-CLM-FL-003: ID Tamper Detection
- **ID**: VR-CLM-FL-003
- **Field**: ID Proof
- **Type**: Security Check
- **Rule**: System must flag suspected tampered ID documents
- **Action**: Additional review required
- **Source**: `Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md`, Scenario 3
- **Priority**: CRITICAL

#### VR-CLM-FL-004: Delivery Date Capture Validation
- **ID**: VR-CLM-FL-004
- **Field**: Policy Bond Delivery Date
- **Type**: Delivery Tracking
- **Rule**: Capture delivery via India Post POD including date, time, recipient acknowledgment
- **Error Message**: "Delivery date capture incomplete - POD details required"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 166-178
- **Priority**: CRITICAL

#### VR-CLM-FL-005: Freelook Timer Auto-Start
- **ID**: VR-CLM-FL-005
- **Field**: Freelook Period
- **Type**: Date Calculation
- **Rule**: Automatically determine start and end date of freelook window based on confirmed delivery date
- **Error Message**: "Freelook timer calculation failed"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 210-220
- **Priority**: CRITICAL

#### VR-CLM-FL-006: Refund Calculation Validation
- **ID**: VR-CLM-FL-006
- **Field**: Refund Amount
- **Type**: Financial Calculation
- **Rule**: Validate refund calculation formula: Premium - Risk Charges - Stamp Duty
- **Error Message**: "Refund calculation error: verify premium amount, risk charges, and stamp duty"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 268-270
- **Priority**: CRITICAL

#### VR-CLM-FL-007: Dispatch ID Generation Validation
- **ID**: VR-CLM-FL-007
- **Field**: Dispatch ID
- **Type**: ID Generation
- **Rule**: System must generate and validate unique Dispatch ID (SP article number) for each policy bond
- **Error Message**: "Dispatch ID generation failed or duplicate ID detected"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 20-24
- **Priority**: HIGH

#### VR-CLM-FL-008: Failed Delivery Escalation Threshold
- **ID**: VR-CLM-FL-008
- **Field**: Delivery Status
- **Type**: Business Logic
- **Rule**: Flag and escalate delivery failures after predefined threshold (e.g., 7 days)
- **Error Message**: "Bond delivery failed - escalation required for manual intervention"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 151-152
- **Priority**: HIGH

#### VR-CLM-FL-009: Authorized Messenger Validation
- **ID**: VR-CLM-FL-009
- **Field**: Authorized Messenger
- **Type**: Documentation Validation
- **Rule**: If policyholder appoints authorized messenger, validate supporting documentation (medical certificate, authorization letter)
- **Error Message**: "Authorized messenger documentation incomplete or invalid"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 243-246
- **Priority**: MEDIUM

#### VR-CLM-FL-010: Refund Account Documentation Validation
- **ID**: VR-CLM-FL-010
- **Field**: Refund Account Details
- **Type**: Documentation Validation
- **Rule**: Require cancelled cheque or POSB account details for refund processing
- **Error Message**: "Refund account documentation missing or invalid"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 273-276
- **Priority**: HIGH

#### VR-CLM-FL-011: Digital Bond Download Timestamp Validation
- **ID**: VR-CLM-FL-011
- **Field**: Download Timestamp
- **Type**: Date/Time Validation
- **Rule**: For ePLI bonds, record download timestamp and treat as official delivery date
- **Error Message**: "ePLI bond download timestamp not captured or invalid"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 172-177
- **Priority**: HIGH

#### VR-CLM-FL-012: Cancellation Request Form Completeness
- **ID**: VR-CLM-FL-012
- **Field**: Cancellation Request
- **Type**: Form Validation
- **Rule**: Cancellation request must include: form, ID proof, original bond, bank details
- **Error Message**: "Cancellation request incomplete - missing required documentation"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 239-250
- **Priority**: HIGH

#### VR-CLM-FL-013: Freelook SLA Reporting Validation
- **ID**: VR-CLM-FL-013
- **Field**: SLA Metrics
- **Type**: Reporting Validation
- **Rule**: Track and report SLA compliance for bond delivery and freelook cancellations
- **Error Message**: "Freelook SLA reporting metrics incomplete"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 368-411
- **Priority**: MEDIUM

#### VR-CLM-FL-014: ePLI DigiLocker Authentication Validation
- **ID**: VR-CLM-FL-014
- **Field**: DigiLocker Authentication
- **Type**: Authentication Validation
- **Rule**: Verify recipient digital signature and DigiLocker authentication for ePLI bonds
- **Error Message**: "DigiLocker authentication failed for ePLI bond delivery"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 172-177
- **Priority**: HIGH

### 4.6 Insurance Ombudsman Validations

#### VR-CLM-OMB-001: Admissibility & Jurisdiction Check
- **ID**: VR-CLM-OMB-001
- **Field**: Ombudsman Complaint
- **Type**: Business Logic
- **Rule**: System determines eligibility (Rule 14) and assigns case to relevant jurisdiction center (Rule 11)
- **Error Message**: "Complaint not admissible or jurisdiction cannot be determined"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 117-119
- **Priority**: HIGH

#### VR-CLM-OMB-002: Complaint Timeline Limitation Check
- **ID**: VR-CLM-OMB-002
- **Field**: Complaint Submission Date
- **Type**: Date Validation
- **Rule**: Check complaint submission against statutory limitation period
- **Error Message**: "Complaint exceeds limitation period - not admissible"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 126-128
- **Priority**: CRITICAL

#### VR-CLM-OMB-003: Claim Value Cap Validation
- **ID**: VR-CLM-OMB-003
- **Field**: Claim Value
- **Type**: Amount Validation
- **Rule**: Validate complaint claim value does not exceed regulatory cap
- **Error Message**: "Claim value exceeds ombudsman jurisdiction limit"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 126-128
- **Priority**: HIGH

#### VR-CLM-OMB-004: Conflict of Interest Screening
- **ID**: VR-CLM-OMB-004
- **Field**: Case Assignment
- **Type**: Business Logic
- **Rule**: System must screen and flag potential conflicts of interest
- **Error Message**: "Conflict of interest detected - case reassignment required"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 123-125
- **Priority**: HIGH

#### VR-CLM-OMB-005: Duplicate Litigation Check
- **ID**: VR-CLM-OMB-005
- **Field**: Litigation Status
- **Type**: Business Logic
- **Rule**: Validate no parallel legal proceedings exist for same complaint
- **Error Message**: "Duplicate/parallel litigation detected - complaint not admissible"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 123-125
- **Priority**: CRITICAL

#### VR-CLM-OMB-006: Hearing Scheduling Validation
- **ID**: VR-CLM-OMB-006
- **Field**: Hearing Schedule
- **Type**: Scheduling Validation
- **Rule**: Validate hearing date, time, mode (physical/video), and party availability
- **Error Message**: "Hearing scheduling conflict detected - rescheduling required"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 130-133
- **Priority**: MEDIUM

#### VR-CLM-OMB-007: Mediation Consent Validation (Rule 16)
- **ID**: VR-CLM-OMB-007
- **Field**: Mediation Consent
- **Type**: Business Logic
- **Rule**: Both parties must provide explicit consent for mediation track
- **Error Message**: "Mediation consent not obtained from all parties"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 134-140
- **Priority**: HIGH

#### VR-CLM-OMB-008: Award Compliance Timeline (30 days)
- **ID**: VR-CLM-OMB-008
- **Field**: Award Compliance Date
- **Type**: Date Validation
- **Rule**: Track insurer/broker compliance with 30-day mandatory timeline
- **Error Message**: "Award compliance deadline breached - escalation required"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 141-145
- **Priority**: CRITICAL

#### VR-CLM-OMB-009: SLA Timer Validation
- **ID**: VR-CLM-OMB-009
- **Field**: SLA Timers
- **Type**: SLA Validation
- **Rule**: Validate SLA timers for complaint acknowledgment and resolution
- **Error Message**: "Ombudsman complaint SLA deadline approaching/breached"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 253-257
- **Priority**: HIGH

#### VR-CLM-OMB-010: Document Upload Security Validation
- **ID**: VR-CLM-OMB-010
- **Field**: Document Upload
- **Type**: Security Validation
- **Rule**: Validate document formats (PDF, JPG, PNG), size limits (10MB), and security
- **Error Message**: "Document upload failed - invalid format or size exceeds limit"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 70-73
- **Priority**: HIGH

#### VR-CLM-OMB-011: Bilingual Field Validation
- **ID**: VR-CLM-OMB-011
- **Field**: Language Support
- **Type**: Content Validation
- **Rule**: Validate bilingual support for all regulatory documents and fields
- **Error Message**: "Bilingual content missing or incomplete"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 99-101
- **Priority**: MEDIUM

#### VR-CLM-OMB-012: Audit Trail Completeness
- **ID**: VR-CLM-OMB-012
- **Field**: Audit Log
- **Type**: Audit Trail
- **Rule**: All actions must be logged with user ID, timestamp, action, IP address
- **Error Message**: "Audit trail incomplete - mandatory logging failed"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 81-83
- **Priority**: CRITICAL

#### VR-CLM-OMB-013: Appellate Authority Routing Validation
- **ID**: VR-CLM-OMB-013
- **Field**: Appeal Routing
- **Type**: Business Logic
- **Rule**: System must route appeals to next higher officer in hierarchy
- **Error Message**: "Appellate authority routing failed - hierarchy not defined"
- **Source**: Claim_SRS_Insurance Ombudsman.md, Lines 145-152
- **Priority**: HIGH

### 4.7 API Integration Validations

#### VR-CLM-API-001: PAN NSDL/Income Tax API Validation
- **ID**: VR-CLM-API-001
- **Field**: PAN Number
- **Type**: API Integration
- **Rule**: PAN must be verified via NSDL/Income Tax API
- **Error Message**: "PAN verification via NSDL/Income Tax API failed"
- **Source**: Claim_SRS_AML triggers & alerts.md, Line 273
- **Priority**: HIGH

#### VR-CLM-API-002: CBS/PFMS Bank Account API
- **ID**: VR-CLM-API-002
- **Field**: Bank Account Details
- **Type**: API Integration
- **Rule**: Verify bank account details via CBS/PFMS API
- **Error Message**: "Bank account verification via CBS/PFMS API failed"
- **Source**: Claim_SRS FRS on Maturity claim.md, Lines 542-546
- **Priority**: CRITICAL

#### VR-CLM-API-003: India Post Tracking API Status
- **ID**: VR-CLM-API-003
- **Field**: Policy Bond Delivery Status
- **Type**: API Integration
- **Rule**: Track delivery status via India Post API
- **Error Message**: "India Post tracking API unavailable or failed"
- **Source**: Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md, Lines 128-134
- **Priority**: MEDIUM

#### VR-CLM-API-004: DigiLocker Document Fetch API
- **ID**: VR-CLM-API-004
- **Field**: DigiLocker Document
- **Type**: API Integration
- **Rule**: Validate DigiLocker document authentication and fetch
- **Error Message**: "DigiLocker document authentication/fetch failed"
- **Source**: Multiple SRS - DigiLocker integration references
- **Priority**: MEDIUM

#### VR-CLM-API-005: Finnet/Fingate AML Filing API
- **ID**: VR-CLM-API-005
- **Field**: AML Filing
- **Type**: API Integration
- **Rule**: Filing Interface for STR/CTR submission to Finnet/Fingate
- **Error Message**: "AML filing submission to Finnet/Fingate failed"
- **Source**: Claim_SRS_AML triggers & alerts.md, Lines 240-243
- **Priority**: CRITICAL

#### VR-CLM-API-006: Batch File Format Schema Validation
- **ID**: VR-CLM-API-006
- **Field**: Batch File Format
- **Type**: Schema Validation
- **Rule**: Validate batch file format against latest FIU-IND XML/JSON schema
- **Error Message**: "Batch file format validation failed - schema mismatch"
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Lines 130-135
- **Priority**: CRITICAL

#### VR-CLM-API-007: STR/CTR Batch Segregation Validation
- **ID**: VR-CLM-API-007
- **Field**: Batch Content Type
- **Type**: Business Logic
- **Rule**: Validate batch contains only one report type (STR or CTR, not mixed)
- **Error Message**: "Batch segregation error - STR and CTR cannot be mixed"
- **Source**: Claim_SRS_AlertsTriggers to FinnetFingate.md, Lines 99-103
- **Priority**: HIGH

### 4.8 Security and Fraud Validations

#### VR-CLM-SEC-001: Digital Signature Validation
- **ID**: VR-CLM-SEC-001
- **Field**: Digital Signature
- **Type**: Security Validation
- **Rule**: Validate digital signatures on approvals and awards
- **Error Message**: "Digital signature validation failed or missing"
- **Source**: Multiple SRS - Digital signature requirements
- **Priority**: HIGH

#### VR-CLM-SEC-002: Duplicate Submission Detection
- **ID**: VR-CLM-SEC-002
- **Field**: Claim Submission
- **Type**: Fraud Detection
- **Rule**: Prevent multiple submissions of same claim across channels
- **Error Message**: "Duplicate claim submission detected across channels"
- **Source**: Cross-channel validation requirements
- **Priority**: HIGH

#### VR-CLM-SEC-003: Timestamp Integrity Check
- **ID**: VR-CLM-SEC-003
- **Field**: System Timestamps
- **Type**: Security Validation
- **Rule**: Ensure all timestamps are system-generated and tamper-proof
- **Error Message**: "Timestamp integrity check failed"
- **Source**: Audit trail requirements across SRS documents
- **Priority**: MEDIUM

#### VR-CLM-SEC-004: User Authorization Level Check
- **ID**: VR-CLM-SEC-004
- **Field**: User Action
- **Type**: Authorization
- **Rule**: Validate user has appropriate role/authority for action
- **Error Message**: "User not authorized for this action"
- **Source**: RBAC requirements across SRS documents
- **Priority**: HIGH

#### VR-CLM-SEC-005: Document Tampering Detection
- **ID**: VR-CLM-SEC-005
- **Field**: Document Metadata
- **Type**: Fraud Detection
- **Rule**: Detect digital alterations, metadata mismatches beyond basic forgery
- **Error Message**: "Document tampering detected - metadata mismatch or digital alteration"
- **Source**: Enhanced fraud detection requirements
- **Priority**: CRITICAL

#### VR-CLM-SEC-006: Biometric/OTP Verification
- **ID**: VR-CLM-SEC-006
- **Field**: User Verification
- **Type**: Security Validation
- **Rule**: Require biometric/OTP verification for high-value claims or sensitive operations
- **Error Message**: "Biometric/OTP verification required but not completed"
- **Source**: Enhanced security requirements for high-value transactions
- **Priority**: MEDIUM

### 4.9 Cross-Field Validations

#### VR-CLM-CROSS-001: Policy Status Transition Validation
- **ID**: VR-CLM-CROSS-001
- **Field**: Policy Status
- **Type**: State Machine Validation
- **Rule**: Validate allowed state transitions (Active → Claimed → Paid)
- **Error Message**: "Invalid policy status transition"
- **Source**: Multiple SRS - State machine requirements
- **Priority**: HIGH

#### VR-CLM-CROSS-002: Claim Amount Calculation Validation
- **ID**: VR-CLM-CROSS-002
- **Field**: Claim Amount
- **Type**: Calculation Validation
- **Rule**: claim_amount = base_sum_assured + accrued_bonuses + excess_premiums - deductions
- **Error Message**: "Claim amount calculation validation failed"
- **Source**: Claim_SRS FRS on death claim.md, Lines 89-100
- **Priority**: CRITICAL

#### VR-CLM-CROSS-003: Outstanding Loan Deduction
- **ID**: VR-CLM-CROSS-003
- **Field**: Loan Deduction
- **Type**: Business Logic
- **Rule**: Deductions include outstanding loans and unpaid premiums
- **Error Message**: "Outstanding loan deduction calculation failed"
- **Source**: Claim_SRS FRS on death claim.md, Lines 94-95
- **Priority**: HIGH

#### VR-CLM-CROSS-004: Tax Deduction Validation
- **ID**: VR-CLM-CROSS-004
- **Field**: Tax Deduction
- **Type**: Calculation Validation
- **Rule**: Applicable taxes must be validated and deducted from claim amount
- **Error Message**: "Tax deduction validation failed"
- **Source**: Claim_SRS FRS on death claim.md, Line 95
- **Priority**: MEDIUM

#### VR-CLM-CROSS-005: Payment Reconciliation Check
- **ID**: VR-CLM-CROSS-005
- **Field**: Payment Status
- **Type**: Reconciliation
- **Rule**: Validate payment status reconciliation with banking systems
- **Error Message**: "Payment reconciliation with banking system failed"
- **Source**: Multiple SRS - Payment reconciliation requirements
- **Priority**: HIGH

#### VR-CLM-CROSS-006: Communication Audit Trail
- **ID**: VR-CLM-CROSS-006
- **Field**: Communication Log
- **Type**: Audit Trail
- **Rule**: All SMS/Email/WhatsApp communications must be logged with timestamps
- **Error Message**: "Communication audit trail logging failed"
- **Source**: Multiple SRS - Communication tracking requirements
- **Priority**: MEDIUM

#### VR-CLM-CROSS-007: Multi-Channel Communication Delivery Confirmation
- **ID**: VR-CLM-CROSS-007
- **Field**: Communication Delivery Status
- **Type**: Delivery Tracking
- **Rule**: Track and validate delivery confirmation for all communication channels
- **Error Message**: "Communication delivery confirmation missing or failed"
- **Source**: Multiple SRS files - multi-channel communication requirements
- **Priority**: MEDIUM

#### VR-CLM-CROSS-008: SMS/Email/WhatsApp Template Validation
- **ID**: VR-CLM-CROSS-008
- **Field**: Communication Template
- **Type**: Content Validation
- **Rule**: All communication templates must include required regulatory disclosures
- **Error Message**: "Communication template validation failed - missing required content"
- **Source**: Multiple SRS files - notification requirements
- **Priority**: MEDIUM

#### VR-CLM-CROSS-009: Notification Audit Trail Validation
- **ID**: VR-CLM-CROSS-009
- **Field**: Notification Log
- **Type**: Audit Trail
- **Rule**: All sent notifications must be logged with timestamp, channel, and status
- **Error Message**: "Notification audit trail incomplete"
- **Source**: Multiple SRS files - notification and audit requirements
- **Priority**: HIGH

#### VR-CLM-CROSS-010: Document Version Control Validation
- **ID**: VR-CLM-CROSS-010
- **Field**: Document Version
- **Type**: Version Control
- **Rule**: Validate document version control and track all modifications
- **Error Message**: "Document version control error - modification history incomplete"
- **Source**: Document management requirements across SRS files
- **Priority**: MEDIUM

---

## 5. Error Codes

### 5.1 Death Claim Error Codes

#### ERR-CLM-DC-001: Policy Not Found
- **ID**: ERR-CLM-DC-001
- **Severity**: CRITICAL
- **Category**: Policy Validation
- **Message**: "Policy not found in system"
- **Trigger**: Invalid policy number entered
- **Resolution**: Verify and enter correct policy number
- **Source**: Claim_SRS FRS on death claim.md
- **Priority**: CRITICAL

#### ERR-CLM-DC-002: Investigation Pending
- **ID**: ERR-CLM-DC-002
- **Severity**: HIGH
- **Category**: Investigation
- **Message**: "Claim approval blocked - investigation report pending"
- **Trigger**: Attempt to approve before investigation completion
- **Resolution**: Wait for investigation report or escalate to supervisor
- **Priority**: HIGH

#### ERR-CLM-DC-003: Document Missing
- **ID**: ERR-CLM-DC-003
- **Severity**: CRITICAL
- **Category**: Documentation
- **Message**: "Mandatory documents missing - cannot proceed"
- **Trigger**: Incomplete document submission
- **Resolution**: Upload missing documents listed in notification
- **Priority**: CRITICAL

#### ERR-CLM-DC-004: Payment Gateway Unavailable
- **ID**: ERR-CLM-DC-004
- **Severity**: HIGH
- **Category**: Disbursement
- **Message**: "Payment gateway temporarily unavailable - retry after some time"
- **Trigger**: Banking system downtime
- **Resolution**: Retry after 30 minutes or use alternate payment mode
- **Priority**: HIGH

#### ERR-CLM-DC-005: Bank Account Verification Failed
- **ID**: ERR-CLM-DC-005
- **Severity**: CRITICAL
- **Category**: Payment Validation
- **Message**: "Bank account verification failed - please verify details"
- **Trigger**: Invalid account number, IFSC, or name mismatch
- **Resolution**: Submit correct bank details with cancelled cheque
- **Priority**: CRITICAL

#### ERR-CLM-DC-006: SLA Breach
- **ID**: ERR-CLM-DC-006
- **Severity**: MEDIUM
- **Category**: Process Monitoring
- **Message**: "Claim processing exceeded SLA - auto-esc

---

## 6. Workflows

### 6.1 Death Claim Workflow (WF-CLM-DC-001)

**Workflow ID**: death-claim-settlement-{claim_id}
**Type**: Long-Running Saga Pattern
**Technology**: Temporal.io (Golang)
**Duration**: 15-45 days
**SLA**: 15 days (no investigation) / 45 days (with investigation)

**Workflow States**:
REGISTERED → DOCUMENT_VERIFICATION → INVESTIGATION → CALCULATION → APPROVAL → DISBURSEMENT → PAID

**Key Activities**:
1. RegisterDeathClaimActivity
2. VerifyDocumentsActivity  
3. CheckInvestigationRequirementActivity
4. AssignInvestigatorActivity
5. ConductInvestigationActivity
6. CalculateClaimAmountActivity
7. ApproveClaimActivity
8. CalculatePenalInterestActivity
9. DisburseClaimActivity
10. UpdateClaimStatusActivity

**Signals**:
- documents-uploaded
- investigation-completed
- approval-decision

**Traceability**:
- Business Rules: BR-CLM-DC-001 to BR-CLM-DC-010
- Functional Requirements: FR-CLM-DC-001 to FR-CLM-DC-010


### 6.2 Maturity Claim Workflow (WF-CLM-MC-001)

**Workflow ID**: maturity-claim-{policy_id}
**Duration**: 7 days
**SLA**: 7 days from submission

**Workflow States**:
INTIMATION_SENT → SUBMITTED → DOCUMENT_VERIFICATION → BANK_VALIDATION → APPROVED → DISBURSED → PAID

**Key Steps**:
1. Generate maturity due report (monthly, first working day)
2. Send multi-channel intimation (SMS, Email, WhatsApp, Portal)
3. Customer submits claim online or at post office
4. Auto-validate documents and policy status
5. Verify bank account via CBS/PFMS API
6. Route for approval with 7-day SLA countdown
7. Disburse via NEFT/IMPS
8. Update status to PAID and send confirmation

**Traceability**:
- Business Rules: BR-CLM-MC-001 to BR-CLM-MC-004
- Functional Requirements: FR-CLM-MC-001 to FR-CLM-MC-006
- SLA: SLA-CLM-MC-001

### 6.3 AML Alert Workflow (WF-CLM-AML-001)

**Workflow ID**: aml-alert-{alert_id}
**Duration**: 7 working days (for STR filing)
**Type**: Event-Driven

**Workflow Triggers**:
1. High cash transaction (>₹50,000)
2. PAN verification failure
3. Nominee change post-death
4. Frequent surrender pattern
5. Refund before bond delivery

**Workflow States**:
TRIGGER_DETECTED → RISK_SCORED → ALERT_GENERATED → OFFICER_REVIEW → ACTION_TAKEN → FILED/CLOSED

**Key Activities**:
1. DetectAMLTriggerActivity
2. CalculateRiskScoreActivity
3. GenerateAlertActivity
4. NotifyAMLOfficerActivity
5. ReviewAlertActivity
6. PrepareSTRCTRActivity
7. SubmitToFinnetActivity
8. LogAuditTrailActivity

**Decision Logic**:
- Risk Level = CRITICAL → Block transaction + File STR immediately
- Risk Level = HIGH → Escalate to AML officer + File CTR
- Risk Level = MEDIUM → Flag for review

**Traceability**:
- Business Rules: BR-CLM-AML-001 to BR-CLM-AML-007
- Functional Requirements: FR-CLM-AML-001

### 6.4 Free Look Cancellation Workflow (WF-CLM-FL-001)

**Workflow ID**: freelook-cancellation-{policy_id}
**Duration**: Within 15 days of delivery
**Type**: Time-Bound

**Workflow States**:
POLICY_DELIVERED → FREELOOK_ACTIVE → CANCELLATION_REQUESTED → VALIDATED → REFUND_CALCULATED → APPROVED → REFUNDED → CLOSED

**Key Steps**:
1. Track policy bond delivery (India Post API / DigiLocker)
2. Start 15-day freelook timer from confirmed delivery
3. Send notifications on day 7 and day 12
4. Accept cancellation request within window
5. Validate documents (original bond, ID proof, bank details)
6. Calculate refund amount
7. Process refund via NEFT/POSB/Cheque
8. Update policy status to "CANCELLED - FREELOOK"
9. Close and archive

**Timer Pause Logic**:
IF delivery_dispute = TRUE THEN pause timer UNTIL fresh_delivery_confirmed

**Traceability**:
- Business Rules: BR-CLM-FL-001 to BR-CLM-FL-003
- Functional Requirements: FR-CLM-FL-001 to FR-CLM-FL-004
- SRS Source: `Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md`, Lines 14-294

### 6.5 Death Claim Investigation Workflow (WF-CLM-DC-INV-001)

**Workflow ID**: death-claim-investigation-{claim_id}-{investigation_id}
**Type**: Long-Running Activity
**Technology**: Temporal.io (Golang)
**Duration**: 21 days (with heartbeat monitoring)
**SLA**: 21 days for investigation report submission
**Max Re-investigations**: 2 (14 days each)

**Workflow States**:
INVESTIGATION_REQUIRED → OFFICER_ASSIGNED → IN_PROGRESS → REPORT_SUBMITTED → REVIEWED → CLOSED

**Actors**:
- Approving Authority (nominates investigator)
- Inquiry Officer (IP/ASP/PRI(P) rank)
- CPC Staff (tracks progress)
- Reviewer (reviews report within 5 days)

**Key Activities**:
1. CheckInvestigationTriggerActivity (death within 3 years of policy issuance/revival)
2. NominateInquiryOfficerActivity (IP/ASP/PRI(P) by jurisdiction)
3. NotifyInvestigatorActivity
4. ConductInvestigationActivity (with heartbeat every 24 hours)
5. VerifyCauseOfDeathActivity
6. CheckHospitalRecordsActivity
7. CheckPoliceRecordsActivity
8. VerifyMaterialSuppressionActivity
9. SubmitInvestigationReportActivity
10. ReviewInvestigationReportActivity (within 5 days)
11. UpdateClaimStatusActivity ('CLEAR', 'SUSPECT', 'FRAUD')

**Decision Points**:
- IF investigation_status = 'CLEAR' THEN proceed to approval
- IF investigation_status = 'SUSPECT' THEN escalate for manual review
- IF investigation_status = 'FRAUD' THEN reject claim with documented evidence
- IF report_quality = 'INADEQUATE' AND reinvestigation_count < 2 THEN trigger reinvestigation (14 days)

**Signals**:
- investigation-report-submitted
- investigation-report-reviewed
- reinvestigation-required

**Queries**:
- GetInvestigationStatus
- GetInvestigatorDetails
- GetInvestigationProgress

**Temporal Specifications**:
```go
// Activity Timeout
ActivityOptions: &temporal.ActivityOptions{
    StartToCloseTimeout: 21 * 24 * time.Hour, // 21 days
    HeartbeatTimeout:    24 * time.Hour,      // Daily heartbeat required
}

// Heartbeat Implementation
func ConductInvestigationActivity(ctx context.Context, claimID string) (*InvestigationReport, error) {
    for dayCount := 0; dayCount < 21; dayCount++ {
        activity.RecordHeartbeat(ctx, fmt.Sprintf("Investigation Day %d/21", dayCount+1))

        // Check if report submitted
        report, err := checkInvestigationReportStatus(ctx, claimID)
        if report != nil && report.Status == "SUBMITTED" {
            return report, nil
        }

        time.Sleep(24 * time.Hour)
    }

    return nil, fmt.Errorf("Investigation exceeded 21-day SLA")
}

// Retry Policy
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    1 * time.Minute,
    BackoffCoefficient: 2.0,
    MaximumInterval:    1 * time.Hour,
    MaximumAttempts:    3,
}
```

**Error Handling**:
- Investigation not submitted within 21 days → Auto-escalate to supervisor
- Officer unavailable → Reassign to alternate officer
- Report quality inadequate → Trigger reinvestigation (max 2 times)

**Traceability**:
- Business Rules: BR-CLM-DC-001, BR-CLM-DC-002, BR-CLM-DC-011, BR-CLM-DC-012
- Functional Requirements: FR-CLM-DC-003
- SRS Source: `Claim_SRS FRS on death claim.md`, Lines 73-86

### 6.6 Death Claim Appeal Workflow (WF-CLM-DC-APPEAL-001)

**Workflow ID**: death-claim-appeal-{claim_id}-{appeal_id}
**Type**: Long-Running Saga Pattern
**Technology**: Temporal.io (Golang)
**Duration**: 45 days from appeal submission
**SLA**: 45 days for appellate decision

**Workflow States**:
APPEAL_SUBMITTED → ELIGIBILITY_VERIFIED → UNDER_REVIEW → DOCUMENTS_REQUESTED → HEARING_SCHEDULED → DECISION_DRAFTED → APPROVED → COMMUNICATED → CLOSED

**Actors**:
- Claimant (submits appeal)
- CPC Staff (validates appeal submission)
- Appellate Authority (next higher officer in approval hierarchy)
- Legal Advisor (for complex cases)

**Key Activities**:
1. ValidateAppealEligibilityActivity (within 90 days of rejection)
2. AssignAppellateAuthorityActivity (next higher officer)
3. NotifyAppellateAuthorityActivity
4. ReviewOriginalClaimActivity
5. RequestAdditionalDocumentsActivity (if needed)
6. ScheduleHearingActivity (if required)
7. ConductAppellateReviewActivity
8. DraftReasonedOrderActivity
9. ApproveAppellateDecisionActivity
10. CommunicateDecisionActivity
11. UpdateClaimStatusActivity (APPROVED or REJECTED-APPEAL)
12. LogAppealAuditTrailActivity

**Decision Points**:
- IF appeal_submission_date > (rejection_date + 90 days) AND no_condonation_request THEN reject appeal
- IF condonation_request = TRUE AND justification_adequate = TRUE THEN accept delayed appeal
- IF additional_documents_required = TRUE THEN request documents (7-day deadline)
- IF case_complexity = 'HIGH' THEN schedule hearing
- IF appellate_decision = 'APPROVE' THEN route to disbursement workflow
- IF appellate_decision = 'REJECT' THEN close claim with final rejection

**Signals**:
- appeal-documents-uploaded
- hearing-completed
- appellate-decision-made

**Queries**:
- GetAppealStatus
- GetAppellateAuthority
- GetAppealTimeline

**Temporal Specifications**:
```go
// Workflow Timeout
WorkflowOptions: &temporal.WorkflowOptions{
    ID:                       fmt.Sprintf("appeal-%s-%s", claimID, appealID),
    TaskQueue:                "death-claim-appeal-queue",
    WorkflowExecutionTimeout: 45 * 24 * time.Hour, // 45 days
}

// Activity Options
ActivityOptions: &temporal.ActivityOptions{
    StartToCloseTimeout: 24 * time.Hour, // Most activities complete within 1 day
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    1 * time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    1 * time.Hour,
        MaximumAttempts:    3,
    },
}

// Signal Handler
func AppealWorkflow(ctx workflow.Context, input AppealInput) error {
    // Set up signal channel
    docChannel := workflow.GetSignalChannel(ctx, "appeal-documents-uploaded")

    // Wait for documents with timeout
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(docChannel, func(c workflow.ReceiveChannel, more bool) {
        var signal DocumentUploadSignal
        c.Receive(ctx, &signal)
        // Process uploaded documents
    })

    // ... workflow logic
}
```

**Error Handling**:
- Appeal window expired without condonation → Auto-reject with notification
- Additional documents not submitted within 7 days → Proceed with available documents
- Appellate authority unavailable → Route to alternate authority
- Decision not issued within 45 days → Auto-escalate to next level

**Traceability**:
- Business Rules: BR-CLM-DC-005, BR-CLM-DC-006, BR-CLM-DC-007
- Functional Requirements: FR-CLM-DC-008
- SRS Source: `Claim_SRS FRS on death claim.md`, Lines 141-151

### 6.7 Death Claim Reopen Workflow (WF-CLM-DC-REOPEN-001)

**Workflow ID**: death-claim-reopen-{claim_id}-{reopen_request_id}
**Type**: Exception Handling Workflow
**Technology**: Temporal.io (Golang)
**Duration**: Variable (depends on reopen reason)

**Workflow States**:
REOPEN_REQUESTED → VALIDATION → APPROVED → CLAIM_RESTORED → RE_ENTERED_MAIN_WORKFLOW → CLOSED

**Actors**:
- CPC User/Supervisor (initiates reopen request)
- Approving Authority (authorizes reopen)
- Claim Handler (processes reopened claim)

**Valid Reopen Circumstances** (from SRS):
1. Court orders
2. New evidence submitted
3. Administrative lapse discovered
4. Claimant appeal accepted
5. Payment failure (technical/banking issue)

**Key Activities**:
1. ValidateReopenRequestActivity
2. CheckReopenEligibilityActivity
3. CreateServiceRequestActivity (new request ID)
4. NotifyStakeholdersActivity
5. RestoreClaimDataActivity
6. ReenterMainWorkflowActivity (trigger WF-CLM-DC-001)
7. LogReopenAuditTrailActivity

**Decision Points**:
- IF reopen_reason NOT IN valid_circumstances THEN reject reopen request
- IF reopen_reason = 'COURT_ORDER' THEN auto-approve and reenter workflow
- IF reopen_reason = 'NEW_EVIDENCE' THEN route for supervisor review
- IF reopen_reason = 'PAYMENT_FAILURE' THEN restart from disbursement step
- IF reopen_reason = 'ADMIN_LAPSE' THEN identify lapse point and resume workflow

**Signals**:
- reopen-approved
- reopen-rejected

**Queries**:
- GetReopenStatus
- GetReopenReason
- GetOriginalClaimDetails

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:        fmt.Sprintf("reopen-%s-%s", claimID, reopenRequestID),
    TaskQueue: "claim-reopen-queue",
}

// Child Workflow for Re-entry
func ReopenWorkflow(ctx workflow.Context, input ReopenInput) error {
    // Validate reopen request
    var validationResult ValidationResult
    err := workflow.ExecuteActivity(ctx, ValidateReopenRequestActivity, input).Get(ctx, &validationResult)

    if !validationResult.Valid {
        return fmt.Errorf("Invalid reopen request: %s", validationResult.Reason)
    }

    // Start child workflow (main death claim workflow)
    childWorkflowOptions := workflow.ChildWorkflowOptions{
        WorkflowID: fmt.Sprintf("death-claim-%s-reopened", input.ClaimID),
    }

    childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

    var claimResult ClaimResult
    err = workflow.ExecuteChildWorkflow(childCtx, DeathClaimWorkflow, input.ClaimData).Get(ctx, &claimResult)

    return err
}
```

**Error Handling**:
- Invalid reopen reason → Reject with notification
- Claim data corruption → Restore from backup
- Workflow re-entry failure → Manual intervention required

**Traceability**:
- Business Rules: BR-CLM-DC-018
- Functional Requirements: FR-CLM-DC-006, FR-CLM-DC-008
- SRS Source: `Claim_SRS FRS on death claim.md`, Lines 123-130

### 6.8 Maturity Report Generation Workflow (WF-CLM-MC-REPORT-001)

**Workflow ID**: maturity-report-generation-{year}-{month}
**Type**: Scheduled Batch Workflow (Cron)
**Technology**: Temporal.io (Golang)
**Schedule**: Monthly, 1st working day, 00:00 AM
**Duration**: 1-2 hours (depends on policy count)

**Workflow States**:
SCHEDULED → RUNNING → POLICIES_QUERIED → REPORT_GENERATED → PUBLISHED → NOTIFICATIONS_SENT → COMPLETED

**Actors**:
- System (automated execution)
- CPC Supervisor (monitors execution)
- Admin Office (receives dashboard view)

**Key Activities**:
1. QueryPoliciesDueForMaturityActivity (next 2 months)
2. CalculateMaturityAmountsActivity (for each policy)
3. GenerateMaturityReportActivity (HO-level detailed report)
4. PublishReportToDashboardActivity
5. GenerateIntimationLettersActivity (prefilled with policy details)
6. SendMultiChannelIntimationsActivity (SMS/Email/WhatsApp)
7. ScheduleRegisteredPostActivity (fallback channel)
8. LogBatchExecutionActivity

**Batch Processing Logic**:
- Query policies with maturity_date BETWEEN (current_date + 30 days) AND (current_date + 60 days)
- Process in batches of 1000 policies
- Calculate maturity amount = SA + bonuses - outstanding loans - unpaid premiums

**Decision Points**:
- IF policy_count > 10,000 THEN process in parallel batches
- IF digital_contact_available = TRUE THEN send SMS/Email/WhatsApp first
- IF digital_contact_unavailable = TRUE THEN use Registered Post only

**Temporal Specifications**:
```go
// Cron Workflow Schedule
Schedule: &temporal.ScheduleSpec{
    Calendars: []temporal.ScheduleCalendarSpec{
        {
            Month:      "*",              // Every month
            DayOfMonth: "1",              // 1st day
            Hour:       "0",              // Midnight
            Minute:     "0",
        },
    },
}

// Batch Processing Activity
func QueryPoliciesDueForMaturityActivity(ctx context.Context) ([]Policy, error) {
    currentDate := time.Now()
    startDate := currentDate.AddDate(0, 0, 30)  // 30 days from now
    endDate := currentDate.AddDate(0, 0, 60)    // 60 days from now

    query := `
        SELECT id, policy_number, customer_id, maturity_date, sum_assured
        FROM policies
        WHERE maturity_date BETWEEN $1 AND $2
        AND status = 'ACTIVE'
        ORDER BY maturity_date ASC
    `

    policies, err := db.Query(query, startDate, endDate)
    return policies, err
}

// Parallel Batch Processing
func GenerateMaturityReportWorkflow(ctx workflow.Context, input ReportInput) error {
    var allPolicies []Policy
    err := workflow.ExecuteActivity(ctx, QueryPoliciesDueForMaturityActivity).Get(ctx, &allPolicies)

    // Process in batches of 1000
    batchSize := 1000
    var futures []workflow.Future

    for i := 0; i < len(allPolicies); i += batchSize {
        end := i + batchSize
        if end > len(allPolicies) {
            end = len(allPolicies)
        }

        batch := allPolicies[i:end]
        future := workflow.ExecuteActivity(ctx, ProcessMaturityBatchActivity, batch)
        futures = append(futures, future)
    }

    // Wait for all batches to complete
    for _, future := range futures {
        err := future.Get(ctx, nil)
        if err != nil {
            return err
        }
    }

    return nil
}
```

**Error Handling**:
- Database query failure → Retry 3 times, then alert admin
- Batch processing failure → Log failed batch, continue with remaining
- Notification service down → Queue notifications for retry

**Traceability**:
- Business Rules: BR-CLM-MC-001, BR-CLM-MC-012
- Functional Requirements: FR-CLM-MC-001
- SRS Source: `Claim_SRS FRS on Maturity claim.md`, Lines 440-447 (FRS-MAT-01), Lines 94-98

### 6.9 Survival Benefit Workflow (WF-CLM-SB-001)

**Workflow ID**: survival-benefit-claim-{policy_id}
**Type**: Long-Running Saga Pattern
**Technology**: Temporal.io (Golang)
**Duration**: 7 days
**SLA**: 7 days from submission

**Workflow States**:
INTIMATION_SENT → SUBMITTED → DOCUMENT_VERIFICATION → ELIGIBILITY_VALIDATED → BANK_VALIDATED → APPROVED → DISBURSED → PAID → CLOSED

**Actors**:
- System (report generation, intimation)
- Policyholder/Insurant (submits claim)
- CPC Staff (scrutiny, verification)
- Supervisor (QC verification)
- Approving Authority (Postmaster/authorized officer)
- Accounts Team (disbursement)

**Key Activities**:
1. GenerateSBDueReportActivity (monthly, 1st working day)
2. SendMultiChannelIntimationActivity (SMS/Email/WhatsApp/Portal)
3. AcceptOnlineSubmissionActivity (Portal/Mobile App)
4. IntegrateDigiLockerActivity (fetch policy document)
5. AutoAcknowledgeActivity (generate Claim ID)
6. ValidateSBEligibilityActivity (policy active, SB due, not already paid)
7. VerifyDocumentsActivity (policy bond, ID proof, bank details)
8. AutoPopulateDataActivity (OCR/data extraction)
9. QCVerificationActivity (digital checklist)
10. ApprovalWorkflowActivity (with digital signature, 7-day SLA)
11. GenerateSanctionLetterActivity (auto-generated with timestamp)
12. BankAccountValidationActivity (API-based via CBS)
13. DisburseViaAutoNEFTActivity
14. GenerateVoucherActivity (auto-generated, digital submission to Accounts)
15. CloseAndArchiveActivity

**Decision Points**:
- IF policy_status != 'ACTIVE' THEN reject (ERR-CLM-SB-RJ-P-02)
- IF sb_already_paid = TRUE THEN reject (ERR-CLM-SB-RJ-P-03)
- IF mandatory_documents_missing = TRUE THEN request documents (auto-reminders)
- IF bank_validation_failed = TRUE THEN prompt re-submission (max 3 attempts)
- IF approval_pending > 7 days THEN auto-escalate with SLA alert

**Signals**:
- documents-uploaded
- digilocker-consent-granted
- approval-decision

**Queries**:
- GetSBClaimStatus
- GetSBEligibility
- GetSBAmountCalculation

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:                       fmt.Sprintf("sb-claim-%s", policyID),
    TaskQueue:                "survival-benefit-queue",
    WorkflowExecutionTimeout: 7 * 24 * time.Hour, // 7 days
}

// SLA Monitoring
func SurvivalBenefitWorkflow(ctx workflow.Context, input SBInput) error {
    slaDeadline := workflow.Now(ctx).Add(7 * 24 * time.Hour)

    // Set SLA timer
    slaTimer := workflow.NewTimer(ctx, 7*24*time.Hour)

    // Process claim
    claimChannel := workflow.NewChannel(ctx)

    selector := workflow.NewSelector(ctx)

    // SLA breach handler
    selector.AddFuture(slaTimer, func(f workflow.Future) {
        workflow.ExecuteActivity(ctx, EscalateSLABreachActivity, input.ClaimID)
    })

    // Claim processing
    selector.AddReceive(claimChannel, func(c workflow.ReceiveChannel, more bool) {
        var result ClaimResult
        c.Receive(ctx, &result)
        // Process result
    })

    selector.Select(ctx)

    return nil
}

// DigiLocker Integration Activity
func IntegrateDigiLockerActivity(ctx context.Context, customerID string) (*Document, error) {
    digiLockerClient := digilocker.NewClient(config.DigiLockerAPIKey)

    // Request consent
    consentURL, err := digiLockerClient.RequestConsent(customerID)
    if err != nil {
        return nil, err
    }

    // Wait for consent (async)
    activity.RecordHeartbeat(ctx, "Waiting for DigiLocker consent")

    // Fetch document
    document, err := digiLockerClient.FetchDocument(customerID, "POLICY_BOND")
    if err != nil {
        return nil, err
    }

    return document, nil
}

// Retry Policy for External APIs
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    2 * time.Second,
    BackoffCoefficient: 2.0,
    MaximumInterval:    1 * time.Minute,
    MaximumAttempts:    3,
    NonRetryableErrorTypes: []string{"InvalidConsentError", "DocumentNotFoundError"},
}
```

**Error Handling**:
- Policy inactive on SB due date → Reject with ERR-CLM-SB-RJ-P-02
- SB already paid → Reject with ERR-CLM-SB-RJ-P-03
- Policyholder identity mismatch → Reject with ERR-CLM-SB-RJ-E-02
- DigiLocker fetch failure → Allow manual upload
- Bank validation failure → Prompt correction (3 attempts)
- SLA breach → Auto-escalate to next authority

**Traceability**:
- Business Rules: BR-CLM-SB-001, BR-CLM-SB-002, BR-CLM-SB-003, BR-CLM-SB-004, BR-CLM-SB-005, BR-CLM-SB-006, BR-CLM-SB-007, BR-CLM-SB-008, BR-CLM-SB-009
- Functional Requirements: FR-CLM-SB-001 to FR-CLM-SB-015
- SRS Source: `Claim_SRS FRS on survival benefit.md`, Lines 92-239 (FRS-SB-01 to FRS-SB-15)

### 6.10 Survival Benefit Report Generation Workflow (WF-CLM-SB-REPORT-001)

**Workflow ID**: sb-report-generation-{year}-{month}
**Type**: Scheduled Batch Workflow (Cron)
**Technology**: Temporal.io (Golang)
**Schedule**: Monthly, 1st working day, 01:00 AM
**Duration**: 1-2 hours

**Workflow States**:
SCHEDULED → RUNNING → POLICIES_QUERIED → REPORT_GENERATED → PUBLISHED → NOTIFICATIONS_SENT → COMPLETED

**Actors**:
- System (automated execution)
- CPC/Admin Office (monitors dashboard)

**Key Activities**:
1. QueryPoliciesDueForSBActivity (next 2 months)
2. CalculateSBAmountsActivity
3. GenerateSBReportActivity (HO-level detailed report)
4. PublishReportToDashboardActivity
5. GenerateIntimationLettersActivity
6. SendMultiChannelIntimationsActivity (SMS/Email/WhatsApp)
7. LogBatchExecutionActivity

**Temporal Specifications**:
```go
// Cron Schedule
Schedule: &temporal.ScheduleSpec{
    Calendars: []temporal.ScheduleCalendarSpec{
        {
            Month:      "*",
            DayOfMonth: "1",
            Hour:       "1",  // 1 AM (stagger from maturity report)
            Minute:     "0",
        },
    },
}

// Same batch processing logic as Maturity Report
```

**Traceability**:
- Business Rules: BR-CLM-SB-001, BR-CLM-SB-003
- Functional Requirements: FR-CLM-SB-001
- SRS Source: `Claim_SRS FRS on survival benefit.md`, Lines 92-102 (FRS-SB-01)

### 6.11 AML STR Filing Workflow (WF-CLM-AML-STR-001)

**Workflow ID**: aml-str-filing-{alert_id}
**Type**: Compliance Workflow
**Technology**: Temporal.io (Golang)
**Duration**: 7 working days from suspicion determination
**SLA**: 7 working days (regulatory requirement)

**Workflow States**:
SUSPICION_DETERMINED → DATA_PREPARED → SCHEMA_VALIDATED → DIGITALLY_SIGNED → SUBMITTED → ACKNOWLEDGED → FILED

**Actors**:
- AML Officer/Nodal Officer (reviews and authorizes)
- Compliance Team (prepares STR)
- System (auto-submission to Finnet/FINGate)

**Key Activities**:
1. DetermineSuspicionActivity (based on AML alert risk score)
2. PrepareSTRDataActivity (collect all transaction details)
3. MapToFIUSchemaActivity (XML/JSON per FIU-IND schema v2.2)
4. ValidateSTRSchemaActivity (against AccountBasedReport.xsd)
5. DigitallySignSTRActivity (pfx/e-token)
6. SubmitToFinnetActivity (via API/SFTP/Portal)
7. ReceiveAcknowledgmentActivity
8. UpdateFilingStatusActivity ('FILED', 'REJECTED', 'PENDING')
9. LogAuditTrailActivity

**Decision Points**:
- IF risk_level = 'CRITICAL' OR risk_level = 'HIGH' THEN prepare STR
- IF schema_validation_failed = TRUE THEN fix errors and revalidate
- IF finnet_submission_failed = TRUE THEN retry (3 attempts) OR manual submission
- IF acknowledgment_received = FALSE AFTER 24 hours THEN alert compliance team

**Signals**:
- str-approved-by-officer
- finnet-acknowledgment-received

**Queries**:
- GetSTRFilingStatus
- GetSTRSubmissionDetails

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:                       fmt.Sprintf("str-filing-%s", alertID),
    TaskQueue:                "aml-compliance-queue",
    WorkflowExecutionTimeout: 7 * 24 * time.Hour, // 7 working days
}

// FIU Schema Mapping Activity
func MapToFIUSchemaActivity(ctx context.Context, amlAlert AMLAlert) (*FIUSTRReport, error) {
    // Map to FIU-IND STR schema
    strReport := &FIUSTRReport{
        ReportType:       "STR",
        ReportingEntity:  "PLI",
        EntityCode:       "PLI001",
        ReportNumber:     fmt.Sprintf("STR-%s-%d", time.Now().Format("20060102"), amlAlert.ID),
        ReportDate:       time.Now(),
        TransactionDetails: FIUTransactionDetails{
            TransactionID:     amlAlert.TransactionID,
            TransactionDate:   amlAlert.TransactionDate,
            TransactionAmount: amlAlert.TransactionAmount,
            TransactionType:   amlAlert.TransactionType,
            SuspicionReason:   amlAlert.TriggerCode,
            RiskLevel:         amlAlert.RiskLevel,
        },
        CustomerDetails: FIUCustomerDetails{
            CustomerID:   amlAlert.CustomerID,
            CustomerName: amlAlert.CustomerName,
            PAN:          amlAlert.PAN,
            Address:      amlAlert.Address,
        },
    }

    return strReport, nil
}

// Digital Signature Activity
func DigitallySignSTRActivity(ctx context.Context, strReport *FIUSTRReport) ([]byte, error) {
    // Load digital certificate
    cert, err := loadDigitalCertificate(config.CertificatePath, config.CertificatePassword)
    if err != nil {
        return nil, err
    }

    // Convert to XML/JSON
    xmlData, err := xml.Marshal(strReport)
    if err != nil {
        return nil, err
    }

    // Sign with certificate
    signedData, err := signWithCertificate(xmlData, cert)
    if err != nil {
        return nil, err
    }

    return signedData, nil
}

// Finnet Submission Activity
func SubmitToFinnetActivity(ctx context.Context, signedSTR []byte) (*FinnetAcknowledgment, error) {
    finnetClient := finnet.NewClient(config.FinnetAPIKey)

    // Submit STR
    response, err := finnetClient.SubmitSTR(signedSTR)
    if err != nil {
        return nil, err
    }

    // Parse acknowledgment
    ack := &FinnetAcknowledgment{
        AcknowledgmentID: response.AckID,
        Status:           response.Status,
        Timestamp:        response.Timestamp,
    }

    return ack, nil
}

// Retry Policy
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    5 * time.Second,
    BackoffCoefficient: 2.0,
    MaximumInterval:    5 * time.Minute,
    MaximumAttempts:    3,
    NonRetryableErrorTypes: []string{"SchemaValidationError", "InvalidCertificateError"},
}
```

**Error Handling**:
- Schema validation failure → Fix errors, revalidate
- Digital signature failure → Check certificate validity, retry
- Finnet submission failure → Retry 3 times, then manual submission via portal
- No acknowledgment within 24 hours → Alert compliance team

**Traceability**:
- Business Rules: BR-CLM-AML-006, BR-CLM-AML-010
- Functional Requirements: FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004, FR-CLM-AML-005
- Integration: INT-CLM-009 (Finnet/FINGate)
- SRS Source: `Claim_SRS_AlertsTriggers to FinnetFingate.md`, Lines 150-180

### 6.12 AML CTR Filing Workflow (WF-CLM-AML-CTR-001)

**Workflow ID**: aml-ctr-filing-{year}-{month}
**Type**: Scheduled Compliance Workflow (Cron)
**Technology**: Temporal.io (Golang)
**Schedule**: Monthly, 5th day of month, 10:00 AM
**Duration**: 2-4 hours
**SLA**: Monthly (regulatory requirement)

**Workflow States**:
SCHEDULED → RUNNING → TRANSACTIONS_AGGREGATED → CTR_PREPARED → SCHEMA_VALIDATED → DIGITALLY_SIGNED → BATCH_CREATED → SUBMITTED → ACKNOWLEDGED → FILED

**Actors**:
- System (automated execution)
- AML Officer (approves CTR batch)
- Compliance Team (monitors submission)

**Key Activities**:
1. AggregateCashTransactionsActivity (cash >₹10L in one day for the month)
2. PrepareCTRBatchActivity (per FIU-IND schema)
3. MapToFIUSchemaActivity (TransactionBasedReport.xsd)
4. ValidateCTRSchemaActivity
5. GenerateUniqueBatchIDActivity
6. DigitallySignCTRActivity
7. SubmitCTRBatchToFinnetActivity
8. ReceiveAcknowledgmentActivity
9. UpdateFilingStatusActivity
10. LogAuditTrailActivity

**Batch Processing Logic**:
- Query all cash transactions WHERE cash_amount > 1000000 AND transaction_date BETWEEN (first day of month) AND (last day of month)
- Group by policy_id
- Single transaction >₹10L OR multiple transactions aggregating >₹10L in one day

**Decision Points**:
- IF no_cash_transactions_above_threshold THEN skip CTR filing
- IF batch_size > 10,000 transactions THEN split into multiple batches
- IF schema_validation_failed = TRUE THEN fix errors and revalidate
- IF finnet_submission_failed = TRUE THEN retry OR manual submission

**Temporal Specifications**:
```go
// Cron Schedule
Schedule: &temporal.ScheduleSpec{
    Calendars: []temporal.ScheduleCalendarSpec{
        {
            Month:      "*",
            DayOfMonth: "5",  // 5th day of each month
            Hour:       "10",
            Minute:     "0",
        },
    },
}

// Aggregate Activity
func AggregateCashTransactionsActivity(ctx context.Context, year int, month int) ([]CashTransaction, error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, 0).Add(-1 * time.Second) // Last second of month

    query := `
        SELECT
            policy_id,
            customer_id,
            transaction_date,
            SUM(cash_amount) AS total_cash,
            COUNT(*) AS transaction_count
        FROM payments
        WHERE payment_mode = 'CASH'
        AND transaction_date BETWEEN $1 AND $2
        GROUP BY policy_id, customer_id, transaction_date
        HAVING SUM(cash_amount) > 1000000
        ORDER BY transaction_date ASC
    `

    transactions, err := db.Query(query, startDate, endDate)
    return transactions, err
}

// CTR Batch Submission
func SubmitCTRBatchToFinnetActivity(ctx context.Context, signedCTR []byte, batchID string) (*FinnetAcknowledgment, error) {
    finnetClient := finnet.NewClient(config.FinnetAPIKey)

    response, err := finnetClient.SubmitCTR(signedCTR, batchID)
    if err != nil {
        return nil, err
    }

    ack := &FinnetAcknowledgment{
        AcknowledgmentID: response.AckID,
        BatchID:          batchID,
        Status:           response.Status,
        RecordCount:      response.RecordCount,
        Timestamp:        response.Timestamp,
    }

    return ack, nil
}
```

**Error Handling**:
- No transactions above threshold → Log "No CTR required for [month]"
- Schema validation failure → Fix errors, revalidate
- Batch size too large → Split into multiple batches
- Submission failure → Retry 3 times, manual submission

**Traceability**:
- Business Rules: BR-CLM-AML-007, BR-CLM-AML-008
- Functional Requirements: FR-CLM-AML-001
- Integration: INT-CLM-009 (Finnet/FINGate)
- SRS Source: `Claim_SRS_AlertsTriggers to FinnetFingate.md`, Lines 115-128

### 6.13 Ombudsman Complaint Intake Workflow (WF-CLM-OMB-001)

**Workflow ID**: ombudsman-complaint-{complaint_id}
**Type**: Long-Running Saga Pattern
**Technology**: Temporal.io (Golang)
**Duration**: Variable (based on resolution path: mediation or adjudication)

**Workflow States**:
SUBMITTED → ADMISSIBILITY_CHECKED → JURISDICTION_MAPPED → REGISTERED → PRELIMINARY_REVIEW → DOCUMENTATION → HEARING_SCHEDULED → MEDIATION/ADJUDICATION → AWARD_ISSUED → COMPLIANCE_TRACKED → CLOSED

**Actors**:
- Complainant (policyholder/nominee)
- CPC Staff (initial intake)
- Ombudsman (DPS rank officer - hypothetical for PLI)
- Support Staff (hearing coordination, document management)
- Insurance Agent/Representative (responds to complaint)

**Key Activities**:
1. AcceptComplaintActivity (multichannel: web, mobile, email, walk-in)
2. ValidateAdmissibilityActivity (Rule 14: prior representation to insurer, statutory limitations)
3. MapJurisdictionActivity (based on complainant pincode/location)
4. CheckConflictOfInterestActivity
5. CheckDuplicateComplaintActivity
6. CheckParallelLitigationActivity
7. RegisterComplaintActivity (assign Case ID)
8. UploadSupportingDocumentsActivity
9. AssignToOmbudsmanActivity
10. PreliminaryScrutinyActivity
11. RequestFurtherInformationActivity (if needed)
12. ScheduleHearingActivity (physical/video)
13. RouteToMediationActivity OR RouteToAdjudicationActivity
14. IssueAwardActivity
15. TrackComplianceActivity (30-day window)
16. CloseComplaintActivity
17. ArchiveWithRetentionActivity

**Decision Points**:
- IF complaint_submission_date < (insurer_representation_date + prescribed_period) THEN reject (premature)
- IF claim_value > ₹50_lakh (hypothetical PLI cap) THEN reject (beyond jurisdiction)
- IF parallel_litigation = TRUE THEN reject
- IF duplicate_complaint = TRUE THEN reject
- IF mediation_successful = TRUE THEN issue mediation recommendation
- IF mediation_failed = TRUE THEN route to adjudication
- IF insurer_compliance_within_30_days = TRUE THEN close
- IF insurer_non_compliance = TRUE THEN escalate to IRDAI (hypothetical for PLI)

**Signals**:
- documents-uploaded-by-complainant
- insurer-response-received
- mediation-outcome
- award-compliance-confirmed

**Queries**:
- GetComplaintStatus
- GetHearingSchedule
- GetAwardDetails

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:        fmt.Sprintf("ombudsman-complaint-%s", complaintID),
    TaskQueue: "ombudsman-queue",
    // No fixed timeout - depends on resolution path
}

// Admissibility Check Activity
func ValidateAdmissibilityActivity(ctx context.Context, complaint Complaint) (*AdmissibilityResult, error) {
    result := &AdmissibilityResult{Valid: true}

    // Check prior representation
    if !complaint.PriorRepresentationToInsurer {
        result.Valid = false
        result.Reason = "No prior representation to insurer (Rule 14)"
        return result, nil
    }

    // Check claim value cap
    if complaint.ClaimValue > 5000000 { // ₹50 lakh
        result.Valid = false
        result.Reason = "Claim value exceeds jurisdiction limit"
        return result, nil
    }

    // Check limitation period
    if complaint.IncidentDate.Add(statutoryLimitationPeriod).Before(time.Now()) {
        result.Valid = false
        result.Reason = "Complaint filed beyond statutory limitation period"
        return result, nil
    }

    // Check parallel litigation
    parallelLitigation, err := checkParallelLitigation(ctx, complaint.PolicyID)
    if err != nil {
        return nil, err
    }

    if parallelLitigation {
        result.Valid = false
        result.Reason = "Parallel litigation pending in court"
        return result, nil
    }

    return result, nil
}

// Jurisdiction Mapping Activity
func MapJurisdictionActivity(ctx context.Context, complaint Complaint) (string, error) {
    // Map based on complainant location
    jurisdictionMap := loadJurisdictionMap()

    ombudsmanCenter := jurisdictionMap.GetCenter(complaint.Pincode)
    if ombudsmanCenter == "" {
        return "", fmt.Errorf("No ombudsman center found for pincode %s", complaint.Pincode)
    }

    return ombudsmanCenter, nil
}

// Mediation vs Adjudication Decision
func OmbudsmanWorkflow(ctx workflow.Context, input ComplaintInput) error {
    // ... admissibility and registration

    // Preliminary review
    var reviewResult ReviewResult
    err := workflow.ExecuteActivity(ctx, PreliminaryScrutinyActivity, input).Get(ctx, &reviewResult)

    // Decision: Mediation or Adjudication?
    if reviewResult.MediationRecommended {
        // Route to mediation
        var mediationOutcome MediationOutcome
        err = workflow.ExecuteActivity(ctx, ConductMediationActivity, input).Get(ctx, &mediationOutcome)

        if mediationOutcome.Settled {
            // Issue mediation recommendation
            err = workflow.ExecuteActivity(ctx, IssueMediationRecommendationActivity, mediationOutcome)
            return err
        }
    }

    // If mediation failed or not recommended, proceed to adjudication
    var award Award
    err = workflow.ExecuteActivity(ctx, ConductAdjudicationActivity, input).Get(ctx, &award)

    // Issue binding award
    err = workflow.ExecuteActivity(ctx, IssueBindingAwardActivity, award)

    // Track compliance (30-day window)
    complianceTimer := workflow.NewTimer(ctx, 30*24*time.Hour)
    complianceChannel := workflow.GetSignalChannel(ctx, "award-compliance-confirmed")

    selector := workflow.NewSelector(ctx)
    selector.AddFuture(complianceTimer, func(f workflow.Future) {
        // Non-compliance escalation
        workflow.ExecuteActivity(ctx, EscalateNonComplianceActivity, input.ComplaintID)
    })

    selector.AddReceive(complianceChannel, func(c workflow.ReceiveChannel, more bool) {
        // Compliance received, close complaint
        workflow.ExecuteActivity(ctx, CloseComplaintActivity, input.ComplaintID)
    })

    selector.Select(ctx)

    return nil
}
```

**Error Handling**:
- Complaint inadmissible → Reject with clear reason and notification
- Jurisdiction mapping failed → Route to default center for manual assignment
- Document upload failure → Allow retry, provide support
- Hearing scheduling conflict → Reschedule with notifications
- Award non-compliance → Escalate to regulatory authority

**Traceability**:
- Business Rules: BR-CLM-OMB-001, BR-CLM-OMB-002, BR-CLM-OMB-003, BR-CLM-OMB-004, BR-CLM-OMB-005, BR-CLM-OMB-006, BR-CLM-OMB-007, BR-CLM-OMB-008
- Functional Requirements: FR-CLM-OMB-001, FR-CLM-OMB-002, FR-CLM-OMB-003, FR-CLM-OMB-004
- Integration: INT-CLM-012 (CPGRAMS)
- SRS Source: `Claim_SRS on insurance ombudsman.md`, Lines 107-169

### 6.14 Ombudsman Hearing Management Workflow (WF-CLM-OMB-HEARING-001)

**Workflow ID**: ombudsman-hearing-{complaint_id}-{hearing_id}
**Type**: Sub-Workflow
**Technology**: Temporal.io (Golang)
**Duration**: Variable (based on hearing type and complexity)

**Workflow States**:
HEARING_REQUESTED → DATE_SCHEDULED → NOTIFICATIONS_SENT → PARTIES_CONFIRMED → HEARING_CONDUCTED → MINUTES_RECORDED → OUTCOME_DETERMINED → CLOSED

**Actors**:
- Ombudsman (presides over hearing)
- Complainant (attends)
- Insurer Representative (attends)
- Support Staff (logistics, minutes recording)

**Key Activities**:
1. CheckHearingRequirementActivity
2. CheckPartiesAvailabilityActivity
3. ScheduleHearingDateActivity
4. SelectHearingModeActivity (physical/video)
5. SendHearingNotificationsActivity (all parties)
6. SetupVideoConferenceActivity (if video hearing)
7. ConductHearingActivity
8. RecordHearingMinutesActivity
9. CollectEvidenceActivity
10. DetermineOutcomeActivity (mediation/adjudication)
11. UpdateComplaintStatusActivity

**Decision Points**:
- IF case_complexity = 'LOW' AND parties_consent = TRUE THEN conduct video hearing
- IF case_complexity = 'HIGH' OR evidence_extensive = TRUE THEN conduct physical hearing
- IF parties_availability_conflict = TRUE THEN reschedule
- IF no_show_by_complainant = TRUE THEN mark complaint as withdrawn
- IF no_show_by_insurer = TRUE THEN proceed ex-parte

**Temporal Specifications**:
```go
// Child Workflow invoked from WF-CLM-OMB-001
func ScheduleHearingActivity(ctx context.Context, complaintID string) (*HearingSchedule, error) {
    // Check parties availability
    availability, err := checkPartiesAvailability(ctx, complaintID)
    if err != nil {
        return nil, err
    }

    // Find common available date
    hearingDate := findCommonAvailableDate(availability)
    if hearingDate.IsZero() {
        return nil, fmt.Errorf("No common availability found")
    }

    // Schedule hearing
    schedule := &HearingSchedule{
        ComplaintID: complaintID,
        HearingDate: hearingDate,
        Mode:        determineHearingMode(ctx, complaintID),
        Location:    determineHearingLocation(ctx, complaintID),
    }

    // Send notifications
    err = sendHearingNotifications(ctx, schedule)

    return schedule, err
}

// Video Hearing Setup
func SetupVideoConferenceActivity(ctx context.Context, schedule HearingSchedule) (*VideoConferenceLink, error) {
    // Integrate with video conferencing platform
    vcClient := videoconf.NewClient(config.VideoConfAPIKey)

    meeting, err := vcClient.ScheduleMeeting(videoconf.MeetingRequest{
        Title:     fmt.Sprintf("Ombudsman Hearing - Complaint %s", schedule.ComplaintID),
        StartTime: schedule.HearingDate,
        Duration:  2 * time.Hour, // 2-hour default
        Attendees: []string{
            schedule.ComplainantEmail,
            schedule.InsurerRepEmail,
            schedule.OmbudsmanEmail,
        },
    })

    if err != nil {
        return nil, err
    }

    link := &VideoConferenceLink{
        MeetingID:  meeting.ID,
        JoinURL:    meeting.JoinURL,
        Passcode:   meeting.Passcode,
        ScheduleAt: schedule.HearingDate,
    }

    return link, nil
}
```

**Error Handling**:
- Scheduling conflict → Find alternate date, notify parties
- Video conference setup failure → Switch to physical hearing
- No-show by party → Follow prescribed rules (withdraw/ex-parte)
- Technical issues during hearing → Adjourn and reschedule

**Traceability**:
- Business Rules: BR-CLM-OMB-003
- Functional Requirements: FR-CLM-OMB-003
- SRS Source: `Claim_SRS on insurance ombudsman.md`, Lines 107-169

### 6.15 Ombudsman Award Issuance Workflow (WF-CLM-OMB-AWARD-001)

**Workflow ID**: ombudsman-award-{complaint_id}-{award_id}
**Type**: Sub-Workflow
**Technology**: Temporal.io (Golang)
**Duration**: 30 days compliance window

**Workflow States**:
AWARD_DRAFTED → REVIEWED → DIGITALLY_SIGNED → ISSUED → COMMUNICATED → COMPLIANCE_PENDING → COMPLIED/ESCALATED

**Actors**:
- Ombudsman (drafts and signs award)
- Legal Advisor (reviews for legal sufficiency)
- Insurer (required to comply within 30 days)
- Complainant (receives award)

**Key Activities**:
1. DraftAwardActivity (mediation recommendation or binding award)
2. CalculateCompensationActivity (with statutory caps)
3. ReviewAwardActivity (legal sufficiency check)
4. ApproveAwardActivity
5. DigitallySignAwardActivity (ombudsman digital signature)
6. IssueAwardActivity (generate official award document)
7. CommunicateAwardActivity (email/WhatsApp/portal/post to all parties)
8. TrackComplianceActivity (30-day timer)
9. VerifyComplianceActivity
10. CloseAwardActivity OR EscalateNonComplianceActivity

**Decision Points**:
- IF award_type = 'MEDIATION_RECOMMENDATION' AND parties_accept = TRUE THEN close as settled
- IF award_type = 'BINDING_AWARD' THEN track compliance
- IF compensation_amount > ₹50_lakh (hypothetical PLI cap) THEN apply cap
- IF insurer_compliance_within_30_days = TRUE THEN close complaint
- IF insurer_non_compliance = TRUE THEN escalate to IRDAI (hypothetical)

**Temporal Specifications**:
```go
// Award Workflow
func AwardIssuanceWorkflow(ctx workflow.Context, input AwardInput) error {
    // Draft award
    var award Award
    err := workflow.ExecuteActivity(ctx, DraftAwardActivity, input).Get(ctx, &award)

    // Calculate compensation (with cap enforcement)
    var compensation Compensation
    err = workflow.ExecuteActivity(ctx, CalculateCompensationActivity, award).Get(ctx, &compensation)

    // Apply statutory cap
    if compensation.Amount > 5000000 { // ₹50 lakh cap
        compensation.Amount = 5000000
        compensation.CapApplied = true
    }

    award.Compensation = compensation

    // Digital signature
    err = workflow.ExecuteActivity(ctx, DigitallySignAwardActivity, award)

    // Issue and communicate
    err = workflow.ExecuteActivity(ctx, IssueAwardActivity, award)
    err = workflow.ExecuteActivity(ctx, CommunicateAwardActivity, award)

    // Track compliance (30-day window)
    complianceTimer := workflow.NewTimer(ctx, 30*24*time.Hour)
    complianceChannel := workflow.GetSignalChannel(ctx, "award-compliance-confirmed")

    selector := workflow.NewSelector(ctx)

    var complianceReceived bool

    selector.AddFuture(complianceTimer, func(f workflow.Future) {
        if !complianceReceived {
            // Non-compliance after 30 days
            workflow.ExecuteActivity(ctx, EscalateNonComplianceActivity, award)
        }
    })

    selector.AddReceive(complianceChannel, func(c workflow.ReceiveChannel, more bool) {
        var signal ComplianceSignal
        c.Receive(ctx, &signal)
        complianceReceived = true

        // Verify compliance
        workflow.ExecuteActivity(ctx, VerifyComplianceActivity, signal)
        workflow.ExecuteActivity(ctx, CloseAwardActivity, award)
    })

    selector.Select(ctx)

    return nil
}

// Compensation Calculation Activity
func CalculateCompensationActivity(ctx context.Context, award Award) (*Compensation, error) {
    compensation := &Compensation{
        ClaimAmount:     award.ClaimAmount,
        InterestAmount:  calculateInterest(award.ClaimAmount, award.DelayDays),
        CostsAmount:     award.Costs,
        TotalAmount:     0,
        CapApplied:      false,
    }

    compensation.TotalAmount = compensation.ClaimAmount + compensation.InterestAmount + compensation.CostsAmount

    // Apply cap
    if compensation.TotalAmount > statutoryCap {
        compensation.TotalAmount = statutoryCap
        compensation.CapApplied = true
    }

    return compensation, nil
}

// Non-Compliance Escalation Activity
func EscalateNonComplianceActivity(ctx context.Context, award Award) error {
    // Log non-compliance
    logNonCompliance(award)

    // Notify IRDAI (hypothetical for PLI - would be PLI Directorate)
    err := notifyRegulator(award, "NON_COMPLIANCE")

    // Update complaint status
    err = updateComplaintStatus(award.ComplaintID, "NON_COMPLIANCE_ESCALATED")

    return err
}
```

**Error Handling**:
- Award drafting errors → Return for revision
- Digital signature failure → Check certificate, retry
- Communication failure → Retry, fallback to registered post
- Non-compliance → Escalate to regulatory authority

**Traceability**:
- Business Rules: BR-CLM-OMB-005, BR-CLM-OMB-006, BR-CLM-OMB-007
- Functional Requirements: FR-CLM-OMB-004
- SRS Source: `Claim_SRS on insurance ombudsman.md`, Lines 134-148

### 6.16 Policy Bond Delivery Tracking Workflow (WF-CLM-BOND-001)

**Workflow ID**: bond-delivery-tracking-{policy_id}-{dispatch_id}
**Type**: Event-Driven Workflow
**Technology**: Temporal.io (Golang)
**Duration**: Variable (until delivery confirmed, max 30 days)

**Workflow States**:
BOND_DISPATCHED → IN_TRANSIT → DELIVERY_ATTEMPTED → DELIVERED/FAILED → CONFIRMED → FREELOOK_ACTIVATED

**Actors**:
- System (tracks via India Post API)
- CPC Staff (monitors delivery)
- Policyholder (receives bond)
- India Post (delivers bond)

**Key Activities**:
1. GenerateDispatchIDActivity (SP article number)
2. LinkToPolicy NumberActivity
3. DispatchViaIndiaPostActivity (Speed Post/Registered Post)
4. TrackDeliveryStatusActivity (via India Post API, poll every 6 hours)
5. CaptureDeliveryDateActivity
6. GetProofOfDeliveryActivity (POD with signature/OTP/photo)
7. ConfirmDeliveryActivity
8. ActivateFreelookTimerActivity (trigger WF-CLM-FL-001)
9. SendDeliveryNotificationActivity (SMS/Email to policyholder)
10. HandleDeliveryFailureActivity (if delivery fails after 3 attempts)

**Decision Points**:
- IF delivery_status = 'DELIVERED' THEN capture delivery date, activate freelook
- IF delivery_attempts >= 3 AND status = 'UNDELIVERED' THEN flag for escalation
- IF delivery_dispute = TRUE THEN pause freelook timer, reship bond
- IF physical_bond = FALSE AND epli_bond = TRUE THEN track DigiLocker download

**Signals**:
- delivery-confirmed
- delivery-disputed
- bond-reshipped

**Queries**:
- GetDeliveryStatus
- GetTrackingDetails
- GetPOD

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:        fmt.Sprintf("bond-delivery-%s-%s", policyID, dispatchID),
    TaskQueue: "bond-tracking-queue",
    WorkflowExecutionTimeout: 30 * 24 * time.Hour, // Max 30 days
}

// Polling Activity for India Post Tracking
func TrackDeliveryStatusActivity(ctx context.Context, dispatchID string) error {
    indiaPostClient := indiapost.NewClient(config.IndiaPostAPIKey)

    // Poll every 6 hours
    for {
        // Record heartbeat
        activity.RecordHeartbeat(ctx, "Tracking delivery")

        // Check delivery status
        status, err := indiaPostClient.GetTrackingStatus(dispatchID)
        if err != nil {
            return err
        }

        // Update database
        err = updateDeliveryStatus(ctx, dispatchID, status)
        if err != nil {
            return err
        }

        // If delivered, exit
        if status.Status == "DELIVERED" {
            return nil
        }

        // If failed after 3 attempts, escalate
        if status.DeliveryAttempts >= 3 && status.Status == "UNDELIVERED" {
            return fmt.Errorf("Delivery failed after 3 attempts")
        }

        // Sleep for 6 hours
        time.Sleep(6 * time.Hour)
    }
}

// Delivery Confirmation Activity
func ConfirmDeliveryActivity(ctx context.Context, dispatchID string) (*DeliveryConfirmation, error) {
    indiaPostClient := indiapost.NewClient(config.IndiaPostAPIKey)

    // Get POD
    pod, err := indiaPostClient.GetProofOfDelivery(dispatchID)
    if err != nil {
        return nil, err
    }

    confirmation := &DeliveryConfirmation{
        DispatchID:      dispatchID,
        DeliveryDate:    pod.DeliveryDate,
        RecipientName:   pod.RecipientName,
        Signature:       pod.Signature,
        OTPVerified:     pod.OTPVerified,
        PhotoEvidence:   pod.PhotoURL,
        Confirmed:       true,
    }

    return confirmation, nil
}

// Freelook Timer Activation (triggers child workflow)
func ActivateFreelookTimerActivity(ctx context.Context, policyID string, deliveryDate time.Time) error {
    // Start freelook workflow
    childWorkflowOptions := workflow.ChildWorkflowOptions{
        WorkflowID: fmt.Sprintf("freelook-%s", policyID),
    }

    childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

    freelookInput := FreelookInput{
        PolicyID:     policyID,
        DeliveryDate: deliveryDate,
        FreelookDays: 15, // 15-day freelook period
    }

    err := workflow.ExecuteChildWorkflow(childCtx, FreelookCancellationWorkflow, freelookInput).Get(ctx, nil)

    return err
}

// Retry Policy
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    30 * time.Second,
    BackoffCoefficient: 1.5,
    MaximumInterval:    6 * time.Hour,
    MaximumAttempts:    0, // Infinite retries until delivered or failed
}
```

**Error Handling**:
- India Post API failure → Retry with exponential backoff
- Delivery failed (3 attempts) → Escalate to CPC, attempt re-dispatch
- Delivery dispute → Pause freelook timer, reinitiate delivery
- Wrong address → Update address, re-dispatch

**Traceability**:
- Business Rules: BR-CLM-BOND-001, BR-CLM-BOND-002
- Functional Requirements: FR-CLM-BOND-001, FR-CLM-BOND-002
- Integration: INT-CLM-011 (India Post Tracking API)
- SRS Source: `Claim_SRS_Tracking of Policy Bond and Free Look Cancellation.md`, Lines 126-153

### 6.17 SLA Monitoring & Escalation Workflow (WF-CLM-MONITORING-001)

**Workflow ID**: sla-monitoring-{claim_type}
**Type**: Continuous Monitoring Workflow
**Technology**: Temporal.io (Golang)
**Duration**: Continuous (long-running)

**Workflow States**:
MONITORING_ACTIVE → SLA_TRACKING → ALERT_GENERATED → ESCALATED → RESOLVED

**Actors**:
- System (continuous monitoring)
- Claim Handler (receives alerts)
- Supervisor (receives escalations)
- Admin Office (dashboard view)

**Key Activities**:
1. MonitorClaimSLAsActivity (continuous)
2. CalculateSLACountdownActivity
3. GenerateColorCodedAlertActivity (Green/Yellow/Orange/Red)
4. SendSLAAlertActivity (to claim handler)
5. AutoEscalateActivity (if SLA breached)
6. UpdateDashboardActivity (real-time)
7. GenerateSLAReportActivity (daily)

**Color-Coded SLA Alerts** (from BR-CLM-DC-021):
- **GREEN**: >50% time remaining
- **YELLOW**: 25%-50% time remaining
- **ORANGE**: 10%-25% time remaining
- **RED**: <10% time remaining OR SLA breached

**Decision Points**:
- IF sla_remaining < 10% THEN escalate to supervisor
- IF sla_breached = TRUE THEN auto-calculate penal interest
- IF claim_status = 'APPROVAL_PENDING' AND sla_breached = TRUE THEN escalate to next higher authority

**Temporal Specifications**:
```go
// Long-Running Monitoring Workflow
func SLAMonitoringWorkflow(ctx workflow.Context, claimType string) error {
    // Continue-As-New pattern for long-running workflows
    for {
        // Monitor all claims of this type
        var claims []Claim
        err := workflow.ExecuteActivity(ctx, FetchActiveClaims Activity, claimType).Get(ctx, &claims)

        // Process each claim
        for _, claim := range claims {
            slaRemaining := claim.SLADueDate.Sub(workflow.Now(ctx))
            slaPercentage := float64(slaRemaining) / float64(claim.TotalSLA)

            // Determine alert level
            var alertLevel string
            if slaPercentage > 0.5 {
                alertLevel = "GREEN"
            } else if slaPercentage > 0.25 {
                alertLevel = "YELLOW"
            } else if slaPercentage > 0.10 {
                alertLevel = "ORANGE"
            } else {
                alertLevel = "RED"
            }

            // Generate alert if needed
            if alertLevel == "ORANGE" || alertLevel == "RED" {
                err = workflow.ExecuteActivity(ctx, GenerateColorCodedAlertActivity, claim, alertLevel)
            }

            // Auto-escalate if breached
            if slaRemaining < 0 {
                err = workflow.ExecuteActivity(ctx, AutoEscalateActivity, claim)
                err = workflow.ExecuteActivity(ctx, CalculatePenalInterestActivity, claim)
            }
        }

        // Update dashboard
        err = workflow.ExecuteActivity(ctx, UpdateDashboardActivity, claims)

        // Sleep for 1 hour
        err = workflow.Sleep(ctx, 1*time.Hour)

        // Continue-As-New every 24 hours to prevent history buildup
        if workflow.Now(ctx).Hour() == 0 {
            return workflow.NewContinueAsNewError(ctx, SLAMonitoringWorkflow, claimType)
        }
    }
}

// Dashboard Update Activity
func UpdateDashboardActivity(ctx context.Context, claims []Claim) error {
    dashboard := &SLADashboard{
        TotalClaims:      len(claims),
        GreenClaims:      0,
        YellowClaims:     0,
        OrangeClaims:     0,
        RedClaims:        0,
        BreachedClaims:   0,
        OnTimeClaims:     0,
        AverageSLAHealth: 0.0,
    }

    for _, claim := range claims {
        switch claim.SLAAlertLevel {
        case "GREEN":
            dashboard.GreenClaims++
        case "YELLOW":
            dashboard.YellowClaims++
        case "ORANGE":
            dashboard.OrangeClaims++
        case "RED":
            dashboard.RedClaims++
        }

        if claim.SLABreached {
            dashboard.BreachedClaims++
        } else {
            dashboard.OnTimeClaims++
        }
    }

    dashboard.AverageSLAHealth = float64(dashboard.OnTimeClaims) / float64(dashboard.TotalClaims) * 100

    // Publish to dashboard
    err := publishDashboard(ctx, dashboard)
    return err
}
```

**Error Handling**:
- Database query failure → Retry, alert admin
- Alert notification failure → Queue for retry
- Escalation routing failure → Manual intervention

**Traceability**:
- Business Rules: BR-CLM-DC-021
- Functional Requirements: FR-CLM-MC-015, FR-CLM-SB-014
- SRS Source: Multiple SRS files (Gap analysis sections)

### 6.18 Customer Feedback Collection Workflow (WF-CLM-FEEDBACK-001)

**Workflow ID**: feedback-collection-{claim_id}
**Type**: Post-Settlement Workflow
**Technology**: Temporal.io (Golang)
**Duration**: 7 days from settlement

**Workflow States**:
SETTLEMENT_COMPLETED → FEEDBACK_REQUEST_SENT → REMINDER_SENT → FEEDBACK_RECEIVED/TIMEOUT → STORED → CLOSED

**Actors**:
- System (automated execution)
- Customer (provides feedback)
- Service Quality Team (analyzes feedback)

**Key Activities**:
1. DetectSettlementCompletionActivity (trigger on claim status = 'PAID')
2. GenerateFeedbackRequestActivity
3. SendFeedbackRequestActivity (SMS/Email/WhatsApp with survey link)
4. SendReminderActivity (if no response after 3 days)
5. CollectFeedbackActivity
6. ValidateFeedbackActivity
7. StoreFeedbackActivity (for analytics)
8. AnalyzeSentimentActivity
9. FlagNegativeFeedbackActivity (for follow-up)
10. UpdateServiceQualityMetricsActivity

**Decision Points**:
- IF feedback_received = FALSE AFTER 3 days THEN send reminder
- IF feedback_received = FALSE AFTER 7 days THEN close (no feedback)
- IF feedback_rating < 3 (out of 5) THEN flag for follow-up
- IF feedback_contains_complaint = TRUE THEN route to complaint management

**Temporal Specifications**:
```go
// Workflow Options
WorkflowOptions: &temporal.WorkflowOptions{
    ID:        fmt.Sprintf("feedback-%s", claimID),
    TaskQueue: "feedback-queue",
    WorkflowExecutionTimeout: 7 * 24 * time.Hour,
}

// Feedback Collection Workflow
func FeedbackCollectionWorkflow(ctx workflow.Context, claimID string) error {
    // Send feedback request
    err := workflow.ExecuteActivity(ctx, GenerateFeedbackRequestActivity, claimID)
    err = workflow.ExecuteActivity(ctx, SendFeedbackRequestActivity, claimID)

    // Wait for feedback with timeout
    feedbackChannel := workflow.GetSignalChannel(ctx, "feedback-submitted")

    // Set reminder timer (3 days)
    reminderTimer := workflow.NewTimer(ctx, 3*24*time.Hour)

    // Set final timeout (7 days)
    finalTimer := workflow.NewTimer(ctx, 7*24*time.Hour)

    selector := workflow.NewSelector(ctx)

    var feedbackReceived bool
    var feedback Feedback

    // Feedback received
    selector.AddReceive(feedbackChannel, func(c workflow.ReceiveChannel, more bool) {
        c.Receive(ctx, &feedback)
        feedbackReceived = true

        // Store and analyze
        workflow.ExecuteActivity(ctx, StoreFeedbackActivity, feedback)
        workflow.ExecuteActivity(ctx, AnalyzeSentimentActivity, feedback)

        // Flag if negative
        if feedback.Rating < 3 {
            workflow.ExecuteActivity(ctx, FlagNegativeFeedbackActivity, feedback)
        }
    })

    // Send reminder after 3 days
    selector.AddFuture(reminderTimer, func(f workflow.Future) {
        if !feedbackReceived {
            workflow.ExecuteActivity(ctx, SendReminderActivity, claimID)
        }
    })

    // Timeout after 7 days
    selector.AddFuture(finalTimer, func(f workflow.Future) {
        if !feedbackReceived {
            // No feedback received, close
            workflow.ExecuteActivity(ctx, CloseFeedbackRequestActivity, claimID)
        }
    })

    selector.Select(ctx)

    return nil
}

// Sentiment Analysis Activity
func AnalyzeSentimentActivity(ctx context.Context, feedback Feedback) (*SentimentAnalysis, error) {
    // Simple sentiment analysis (can be replaced with ML model)
    analysis := &SentimentAnalysis{
        FeedbackID: feedback.ID,
        Rating:     feedback.Rating,
        Sentiment:  "NEUTRAL",
    }

    if feedback.Rating >= 4 {
        analysis.Sentiment = "POSITIVE"
    } else if feedback.Rating <= 2 {
        analysis.Sentiment = "NEGATIVE"
    }

    // Analyze text comments (basic keyword matching)
    negativeKeywords := []string{"poor", "bad", "worst", "terrible", "disappointed", "slow", "delay"}
    positiveKeywords := []string{"good", "excellent", "great", "fast", "helpful", "satisfied"}

    comments := strings.ToLower(feedback.Comments)

    for _, keyword := range negativeKeywords {
        if strings.Contains(comments, keyword) {
            analysis.Sentiment = "NEGATIVE"
            break
        }
    }

    for _, keyword := range positiveKeywords {
        if strings.Contains(comments, keyword) && analysis.Sentiment != "NEGATIVE" {
            analysis.Sentiment = "POSITIVE"
        }
    }

    return analysis, nil
}
```

**Error Handling**:
- Notification service failure → Retry, queue for later
- Feedback submission error → Allow retry, provide support
- Storage failure → Retry, alert admin

**Traceability**:
- Business Rules: BR-CLM-MC-004 (inferred)
- Functional Requirements: FR-CLM-MC-014, FR-CLM-SB-013
- SRS Source: `Claim_SRS FRS on Maturity claim.md`, Lines 571-578 (FRS-MAT-14)

---

## 7. Data Entities

### 7.1 Claim Entity (E-CLM-DC-001)

**Table Name**: claims
**Description**: Master table for all claim types

| Attribute | Type | Required | Constraints | Description |
|-----------|------|----------|-------------|-------------|
| id | UUID | Yes | PK | Unique claim identifier |
| claim_number | VARCHAR(20) | Yes | UK | System-generated claim number (CLM-YYYYMMDD-NNNN) |
| claim_type | ENUM | Yes | 'DEATH', 'MATURITY', 'SURVIVAL_BENEFIT', 'FREELOOK' | Type of claim |
| policy_id | UUID | Yes | FK → policies.id | Reference to policy |
| customer_id | UUID | Yes | FK → customers.id | Reference to customer |
| claim_date | DATE | Yes | | Date claim submitted |
| claimant_name | VARCHAR(200) | Yes | | Name of person claiming |
| claimant_relation | VARCHAR(50) | No | | Relationship to policyholder |
| status | ENUM | Yes | See status list below | Current claim status |
| claim_amount | NUMERIC(15,2) | No | | Calculated claim amount |
| approved_amount | NUMERIC(15,2) | No | | Final approved amount |
| penal_interest | NUMERIC(10,2) | No | | Penal interest if SLA breached |
| investigation_required | BOOLEAN | No | Default: FALSE | Whether investigation needed |
| investigation_status | VARCHAR(20) | No | | 'CLEAR', 'SUSPECT', 'FRAUD' |
| investigator_id | UUID | No | FK → users.id | Assigned investigator |
| approver_id | UUID | No | FK → users.id | Approving authority |
| approval_date | TIMESTAMP | No | | Date of approval/rejection |
| disbursement_date | TIMESTAMP | No | | Date of payment |
| payment_mode | VARCHAR(20) | No | | 'NEFT', 'POSB_EFT', 'CHEQUE' |
| transaction_id | VARCHAR(100) | No | | Bank transaction reference |
| rejection_reason | TEXT | No | | Reason if rejected |
| sla_due_date | TIMESTAMP | Yes | | SLA deadline |
| sla_breached | BOOLEAN | No | Default: FALSE | Whether SLA was breached |
| created_at | TIMESTAMP | Yes | Default: NOW() | Record creation timestamp |
| updated_at | TIMESTAMP | Yes | Default: NOW() | Last update timestamp |
| created_by | UUID | Yes | FK → users.id | User who created record |
| updated_by | UUID | Yes | FK → users.id | User who last updated |
| version | INTEGER | Yes | Default: 1 | Optimistic locking version |

**Claim Status Values**:
- REGISTERED
- DOCUMENT_PENDING
- DOCUMENT_VERIFIED
- INVESTIGATION_PENDING
- INVESTIGATION_COMPLETED
- CALCULATION_COMPLETED
- APPROVAL_PENDING
- APPROVED
- REJECTED
- DISBURSEMENT_PENDING
- PAID
- CLOSED
- RETURNED
- REOPENED

**Indexes**:
```sql
CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_customer_id ON claims(customer_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_claim_type ON claims(claim_type);
CREATE INDEX idx_claims_sla_due_date ON claims(sla_due_date) WHERE status IN ('APPROVAL_PENDING', 'INVESTIGATION_PENDING');
CREATE INDEX idx_claims_created_at ON claims(created_at);
```

### 7.2 Claim Documents Entity (E-CLM-DC-002)

**Table Name**: claim_documents
**Description**: Documents uploaded for claim processing

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| id | UUID | Yes | PK |
| claim_id | UUID | Yes | FK → claims.id |
| document_type | VARCHAR(50) | Yes | Type of document |
| document_name | VARCHAR(255) | Yes | Original filename |
| document_url | TEXT | Yes | S3/storage URL |
| file_size | INTEGER | Yes | File size in bytes |
| uploaded_by | UUID | Yes | FK → users.id |
| uploaded_at | TIMESTAMP | Yes | Upload timestamp |
| verified | BOOLEAN | No | Whether verified by officer |
| verified_by | UUID | No | FK → users.id |
| verified_at | TIMESTAMP | No | Verification timestamp |

**Document Types**:
- DEATH_CERTIFICATE
- CLAIM_FORM
- POLICY_BOND
- INDEMNITY_BOND
- CLAIMANT_ID_PROOF
- BANK_MANDATE
- FIR (for unnatural death)
- POSTMORTEM_REPORT (for unnatural death)
- SUCCESSION_CERTIFICATE
- LEGAL_HEIR_AFFIDAVIT


### 7.3 AML Alert Entity (E-CLM-AML-001)

**Table Name**: aml_alerts

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| id | UUID | Yes | PK |
| alert_id | VARCHAR(50) | Yes | UK - System-generated |
| trigger_code | VARCHAR(20) | Yes | AML_001 to AML_005 |
| policy_id | UUID | Yes | FK → policies.id |
| transaction_type | VARCHAR(50) | Yes | Type of transaction |
| transaction_amount | NUMERIC(15,2) | No | Transaction amount |
| risk_level | ENUM | Yes | 'LOW', 'MEDIUM', 'HIGH', 'CRITICAL' |
| alert_status | VARCHAR(20) | Yes | 'FLAGGED', 'UNDER_REVIEW', 'FILED', 'CLOSED' |
| reviewed_by | UUID | No | FK → users.id |
| reviewed_at | TIMESTAMP | No | Review timestamp |
| action_taken | TEXT | No | Action description |
| filing_required | BOOLEAN | No | STR/CTR filing needed |
| filing_type | VARCHAR(10) | No | 'STR', 'CTR' |
| filing_status | VARCHAR(20) | No | 'PENDING', 'FILED', 'REJECTED' |
| filed_at | TIMESTAMP | No | Filing timestamp |
| created_at | TIMESTAMP | Yes | Alert creation time |

---

## 8. Integration Points

### INT-CLM-001: ECMS (Electronic Content Management System)
- **Purpose**: Document storage and retrieval
- **Type**: REST API / File Upload
- **Operations**:
  - Upload document
  - Retrieve document
  - Tag document to claim ID
  - Search documents
- **Data Format**: Multipart form-data for upload, JSON for metadata
- **Authentication**: OAuth 2.0
- **Endpoint**: https://ecms.pli.gov.in/api/v1/documents
- **Traceability**: FR-CLM-DC-002, FR-CLM-MC-005

### INT-CLM-002: Finacle/IT 2.0 (Core Banking)
- **Purpose**: Payment execution and bank account validation
- **Type**: SOAP / REST API
- **Operations**:
  - Validate bank account (CBS/PFMS API)
  - Execute NEFT payment
  - Execute POSB EFT
  - Get transaction status
  - Reconciliation data
- **Data Format**: XML (SOAP) / JSON (REST)
- **Authentication**: Certificate-based / API Key
- **SLA**: Real-time (< 5 seconds)
- **Traceability**: FR-CLM-DC-006, FR-CLM-MC-006, BR-CLM-MC-003

### INT-CLM-003: NEFT/POSB EFT (Payment Gateways)
- **Purpose**: Electronic fund transfer
- **Type**: Banking API
- **Operations**:
  - Initiate NEFT transfer
  - Initiate POSB EFT
  - Check payment status
  - Get payment acknowledgment
- **Batch Processing**: Yes (for bulk payments)
- **Settlement Time**: T+0 (NEFT), Real-time (POSB EFT)
- **Traceability**: FR-CLM-DC-006

### INT-CLM-004: SMS Gateway
- **Purpose**: Send SMS notifications
- **Type**: HTTP API
- **Provider**: NIC / Commercial SMS gateway
- **Operations**:
  - Send single SMS
  - Send bulk SMS
  - Get delivery status
- **Templates**: Pre-approved DLT templates
- **Traceability**: FR-CLM-DC-007, FR-CLM-MC-002

### INT-CLM-005: Email Service
- **Purpose**: Send email notifications
- **Type**: SMTP / API
- **Operations**:
  - Send email with attachments
  - Send bulk emails
  - Track email status
- **Templates**: HTML email templates
- **Traceability**: FR-CLM-DC-007, FR-CLM-MC-002

### INT-CLM-006: WhatsApp Business API
- **Purpose**: Send WhatsApp notifications
- **Type**: REST API
- **Provider**: WhatsApp Business Platform
- **Operations**:
  - Send template message
  - Send document
  - Get message status
- **Compliance**: GDPR compliant, opt-in required
- **Traceability**: FR-CLM-MC-002

### INT-CLM-007: DigiLocker
- **Purpose**: Fetch digital policy documents
- **Type**: REST API
- **Authentication**: Aadhaar-based / OAuth
- **Operations**:
  - Fetch document by URI
  - Verify document authenticity
  - Get user consent
- **Data Format**: JSON (metadata), PDF/XML (documents)
- **Endpoint**: https://api.digitallocker.gov.in
- **Traceability**: FR-CLM-MC-003, FR-CLM-SB-002

### INT-CLM-008: Customer Portal / Mobile App
- **Purpose**: Self-service claim submission and tracking
- **Type**: REST API
- **Operations**:
  - Submit claim online
  - Upload documents
  - Track claim status
  - View claim history
  - Download acknowledgment
- **Authentication**: OTP / Biometric
- **Platform**: Web (React), Mobile (Flutter)
- **Traceability**: FR-CLM-MC-003

### INT-CLM-009: Finnet/FINGate (AML Reporting)
- **Purpose**: Submit STR/CTR to FIU-India
- **Type**: REST API / SFTP / Portal Upload
- **Operations**:
  - Submit STR batch
  - Submit CTR batch
  - Get acknowledgment
  - Check filing status
- **Data Format**: XML/JSON per FIU schema v2.2
- **Digital Signature**: Mandatory (pfx / e-token)
- **Endpoint**: https://finnet.gov.in/fingate
- **Traceability**: FR-CLM-AML-010, FR-CLM-AML-011, BR-CLM-AML-006, BR-CLM-AML-007

### INT-CLM-010: PAN Verification API
- **Purpose**: Verify PAN details
- **Type**: REST API
- **Provider**: NSDL / API Setu
- **Operations**:
  - Verify PAN
  - Get PAN holder name
- **Response Time**: < 2 seconds
- **Data Format**: JSON
- **Traceability**: VR-CLM-AML-002

### INT-CLM-011: India Post Tracking API
- **Purpose**: Track policy bond delivery
- **Type**: REST API
- **Operations**:
  - Track article by number
  - Get delivery status
  - Get proof of delivery (POD)
  - Get delivery confirmation
- **Real-time**: Yes
- **Endpoint**: https://track.indiapost.gov.in/api
- **Traceability**: BR-CLM-FL-001, FR-CLM-FL-001

### INT-CLM-012: CPGRAMS (Grievance Portal)
- **Purpose**: Integrate grievances and complaint tracking
- **Type**: REST API
- **Operations**:
  - Register grievance
  - Update grievance status
  - Close grievance
  - Fetch grievance details
- **Endpoint**: https://pgportal.gov.in/api
- **Traceability**: FR-CLM-OMB-001

---

## 9. Temporal Workflows (Golang Implementation Notes)

### 9.1 Workflow Design Principles

1. **Determinism**: All workflows must be deterministic
2. **Idempotency**: Activities must be idempotent for safe retries
3. **Long-Running**: Use Continue-As-New for workflows > 50,000 events
4. **Compensation**: Implement saga pattern for distributed transactions
5. **Versioning**: Plan for workflow versioning from Day 1

### 9.2 Activity Implementation Guidelines

**Activity Timeout Settings**:
- Document verification: 10 minutes
- Investigation: 21 days with heartbeat every 24 hours
- Approval: 15-45 days based on claim type
- Disbursement: 30 minutes
- Notification: 5 minutes

**Retry Policy**:
```go
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    1 * time.Second,
    BackoffCoefficient: 2.0,
    MaximumInterval:    1 * time.Minute,
    MaximumAttempts:    3,
}
```

**Heartbeat for Long Activities**:
```go
// In ConductInvestigationActivity
func ConductInvestigationActivity(ctx context.Context, claimID string) (*InvestigationReport, error) {
    for {
        // Record heartbeat every hour
        activity.RecordHeartbeat(ctx, "Investigation in progress")
        
        // Check if report submitted
        report, err := checkInvestigationReport(ctx, claimID)
        if report != nil {
            return report, nil
        }
        
        time.Sleep(1 * time.Hour)
    }
}
```

### 9.3 Signal Usage

**Signal for Document Upload**:
```go
// In workflow
docChannel := workflow.GetSignalChannel(ctx, "documents-uploaded")

var uploadSignal DocumentUploadSignal
docChannel.Receive(ctx, &uploadSignal)

// External trigger
err := temporalClient.SignalWorkflow(ctx, workflowID, "", "documents-uploaded", DocumentUploadSignal{
    DocumentIDs: []string{"doc1", "doc2"},
})
```

**Signal for Investigation Completion**:
```go
investigationChannel := workflow.GetSignalChannel(ctx, "investigation-completed")

var investigationSignal InvestigationCompletedSignal
investigationChannel.Receive(ctx, &investigationSignal)
```

### 9.4 Query Usage

**Get Claim Status**:
```go
// Register query handler in workflow
err := workflow.SetQueryHandler(ctx, "GetClaimStatus", func() (string, error) {
    return currentStatus, nil
})

// External query
response, err := temporalClient.QueryWorkflow(ctx, workflowID, "", "GetClaimStatus")
var status string
response.Get(&status)
```

---

## 10. Traceability Matrix

### 10.1 SRS to Business Rules Mapping

| SRS Document | Section/Reference | Business Rule ID | Description |
|--------------|-------------------|------------------|-------------|
| Claim_SRS FRS on death claim.md | Section 3, Death Investigation | BR-CLM-DC-001 | Investigation trigger (death within 3 years) |
| Claim_SRS FRS on death claim.md | Section 4.2, Investigation Timeline | BR-CLM-DC-002 | Investigation report due within 21 days |
| Claim_SRS FRS on death claim.md | Section 5.1, Approval Process | BR-CLM-DC-003 | Approval within 15 days (no investigation) |
| Claim_SRS FRS on death claim.md | Section 5.2, Approval with Investigation | BR-CLM-DC-004 | Approval within 45 days (with investigation) |
| Claim_SRS FRS on death claim.md | Section 6, Appeal Process | BR-CLM-DC-005 | Appeal within 90 days of rejection |
| Claim_SRS FRS on death claim.md | Section 6.1, Appellate Authority | BR-CLM-DC-006 | Next higher officer handles appeal |
| Claim_SRS FRS on death claim.md | Section 6.2, Appeal Decision | BR-CLM-DC-007 | Appeal decision within 45 days |
| Claim_SRS FRS on death claim.md | Section 7, Settlement Calculation | BR-CLM-DC-008 | Claim amount = SA + bonuses - deductions |
| Claim_SRS FRS on death claim.md | Section 8, Penal Interest | BR-CLM-DC-009 | Penal interest 8% p.a. on SLA breach |
| Claim_SRS FRS on death claim.md | Section 9, Document Management | BR-CLM-DC-010 | Document pending return after 22 days |
| Claim_SRS FRS on death claim.md | Gap Analysis - TC-7 | BR-CLM-DC-011 | Investigation report review within 5 days |
| Claim_SRS FRS on death claim.md | Gap Analysis - SL-4 | BR-CLM-DC-012 | Re-investigation max 2 times, 14 days each |
| Claim_SRS FRS on Maturity claim.md | Section 2, Report Generation | BR-CLM-MC-001 | Monthly maturity report on 1st working day |
| Claim_SRS FRS on Maturity claim.md | Section 3, Approval SLA | BR-CLM-MC-002 | Maturity claim approval within 7 days |
| Claim_SRS FRS on Maturity claim.md | Section 4, Bank Verification | BR-CLM-MC-003 | Bank verification via CBS/PFMS API |
| Claim_SRS FRS on Maturity claim.md | Section 5, Communication | BR-CLM-MC-004 | Multi-channel intimation (SMS/Email/WhatsApp) |
| Claim_SRS FRS on survival benefit.md | Section 2, SB Report | BR-CLM-SB-001 | Monthly SB report on 1st working day |
| Claim_SRS FRS on survival benefit.md | Section 3, Approval SLA | BR-CLM-SB-002 | SB approval within 7 days |
| Claim_SRS_AML triggers & alerts.md | Section 2.1, Cash Premium | BR-CLM-AML-001 | Cash >₹50K triggers high-risk alert |
| Claim_SRS_AML triggers & alerts.md | Section 2.2, PAN Verification | BR-CLM-AML-002 | PAN mismatch triggers medium-risk alert |
| Claim_SRS_AML triggers & alerts.md | Section 2.3, Critical Alert | BR-CLM-AML-003 | Nominee change post death - critical alert |
| Claim_SRS_AML triggers & alerts.md | Section 2.4, Surrender Pattern | BR-CLM-AML-004 | >3 surrenders in 6 months triggers alert |
| Claim_SRS_AML triggers & alerts.md | Section 2.5, Refund Alert | BR-CLM-AML-005 | Refund before bond dispatch - high-risk |
| Claim_SRS_AlertsTriggers to FinnetFingate.md | Section 3, STR Timeline | BR-CLM-AML-006 | STR filing within 7 working days |
| Claim_SRS_AlertsTriggers to FinnetFingate.md | Section 4, CTR Schedule | BR-CLM-AML-007 | Monthly CTR for cash >₹10L in one day |
| Claim_SRS on insurance ombudsman.md | Section 2, Jurisdiction | BR-CLM-OMB-003 | Dynamic jurisdiction mapping |

---

### 10.2 Business Rules to Functional Requirements Mapping

| Business Rule ID | Functional Requirement IDs | Relationship |
|------------------|----------------------------|--------------|
| BR-CLM-DC-001 | FR-CLM-DC-001 | Investigation trigger auto-detected during claim registration |
| BR-CLM-DC-002 | FR-CLM-DC-003 | Investigation report review enforces timeline |
| BR-CLM-DC-003 | FR-CLM-DC-001 | Approval timeline tracked without investigation |
| BR-CLM-DC-004 | FR-CLM-DC-003 | Approval timeline tracked with investigation |
| BR-CLM-DC-005 | FR-CLM-DC-001 | Appeal window validation during claim registration |
| BR-CLM-DC-006 | FR-CLM-DC-001 | Appellate authority auto-assigned |
| BR-CLM-DC-007 | FR-CLM-DC-001 | Appeal decision timeline enforced |
| BR-CLM-DC-008 | FR-CLM-DC-001 | Claim amount calculation implemented |
| BR-CLM-DC-009 | FR-CLM-DC-007 | Penal interest auto-calculation for death claims |
| BR-CLM-DC-010 | FR-CLM-DC-006 | Document pending return automated |
| BR-CLM-DC-011 | FR-CLM-DC-003 | Investigation officer assignment rules (IP/ASP/PRI(P)) |
| BR-CLM-DC-012 | FR-CLM-DC-003, FR-CLM-DC-004 | Investigation status classification (Clear/Suspect/Fraud) |
| BR-CLM-DC-013 | FR-CLM-DC-002 | Unnatural death document requirements (FIR/postmortem) |
| BR-CLM-DC-014 | FR-CLM-DC-002 | Nomination absence document requirements |
| BR-CLM-DC-015 | FR-CLM-DC-002 | Mandatory document checklist (5 base documents) |
| BR-CLM-DC-016 | FR-CLM-DC-003, FR-CLM-DC-004 | Manual calculation override with audit trail |
| BR-CLM-DC-017 | FR-CLM-DC-005 | Payment mode selection rules (NEFT/POSB/Cheque) |
| BR-CLM-DC-018 | FR-CLM-DC-006, FR-CLM-DC-008 | Claim reopen valid circumstances |
| BR-CLM-DC-019 | FR-CLM-DC-007 | Communication milestone triggers |
| BR-CLM-DC-020 | FR-CLM-DC-001 | Claim case owner assignment at CPC |
| BR-CLM-DC-021 | FR-CLM-DC-001, FR-CLM-DC-004 | SLA color-coded alert system (Green/Yellow/Orange/Red) |
| BR-CLM-DC-022 | FR-CLM-DC-005 | Disbursement reconciliation with banking (Finacle/IT 2.0) |
| BR-CLM-DC-023 | FR-CLM-DC-008 | Rejection with root cause analysis |
| BR-CLM-DC-024 | FR-CLM-DC-007 | Real-time claimant status tracking (SMS/Email/Portal/Mobile) |
| BR-CLM-DC-025 | FR-CLM-DC-003, FR-CLM-DC-004 | Audit trail for manual overrides with digital signature |
| BR-CLM-MC-001 | FR-CLM-MC-001 | Monthly maturity report auto-generation |
| BR-CLM-MC-002 | FR-CLM-MC-008 | Maturity approval SLA tracking (7 days) |
| BR-CLM-MC-003 | FR-CLM-MC-010 | Bank verification via CBS/PFMS API |
| BR-CLM-MC-004 | FR-CLM-MC-002 | Multi-channel intimation (SMS/Email/WhatsApp) |
| BR-CLM-MC-005 | FR-CLM-MC-003, FR-CLM-MC-004 | Policy activation status validation |
| BR-CLM-MC-006 | FR-CLM-MC-003, FR-CLM-MC-004 | Duplicate maturity claim prevention |
| BR-CLM-MC-007 | FR-CLM-MC-003, FR-CLM-MC-004 | Policy forfeiture/surrender pre-maturity check |
| BR-CLM-MC-008 | FR-CLM-MC-003, FR-CLM-MC-004 | Claimant identity verification |
| BR-CLM-MC-009 | FR-CLM-MC-004 | Forged/suspicious document detection |
| BR-CLM-MC-010 | FR-CLM-MC-004 | Missing document auto-reminder schedule |
| BR-CLM-MC-011 | FR-CLM-MC-010 | Bank account re-submission workflow (max 3 attempts) |
| BR-CLM-MC-012 | FR-CLM-MC-001, FR-CLM-MC-003 | Automated maturity date calculation |
| BR-CLM-SB-001 | FR-CLM-SB-001 | SB report auto-generation |
| BR-CLM-SB-002 | FR-CLM-SB-008 | SB approval SLA tracking with digital signature |
| BR-CLM-SB-003 | FR-CLM-SB-001 | SB due date calculation and report generation |
| BR-CLM-SB-004 | FR-CLM-SB-004, FR-CLM-SB-008 | SB eligibility validation and approval |
| BR-CLM-SB-005 | FR-CLM-SB-002 | DigiLocker integration for document fetching |
| BR-CLM-SB-006 | FR-CLM-SB-001 | Auto-acknowledgment with claim ID generation |
| BR-CLM-SB-007 | FR-CLM-SB-005 | Automatic indexing as service request |
| BR-CLM-SB-008 | FR-CLM-SB-007 | OCR auto-population with supervisor verification |
| BR-CLM-SB-009 | FR-CLM-SB-008 | Digital signature for approval/rejection |
| BR-CLM-AML-001 | FR-CLM-AML-001 | Cash premium alert detection (>₹50K) |
| BR-CLM-AML-002 | FR-CLM-AML-002 | PAN mismatch alert detection |
| BR-CLM-AML-003 | FR-CLM-AML-003 | Nominee change post death critical alert and blocking |
| BR-CLM-AML-004 | FR-CLM-AML-004 | Surrender pattern analysis (>3 in 6 months) |
| BR-CLM-AML-005 | FR-CLM-AML-005 | Refund before bond delivery alert detection |
| BR-CLM-AML-006 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004, FR-CLM-AML-005 | STR filing workflow initiation for all AML alerts |
| BR-CLM-AML-007 | FR-CLM-AML-001 | CTR filing workflow initiation |
| BR-CLM-AML-008 | FR-CLM-AML-001 | CTR aggregate monitoring and monthly filing |
| BR-CLM-AML-009 | FR-CLM-AML-002 | Third-party PAN verification and blocking |
| BR-CLM-AML-010 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004, FR-CLM-AML-005 | Regulatory reporting to FIU-IND |
| BR-CLM-AML-011 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004 | Negative list daily screening |
| BR-CLM-AML-012 | FR-CLM-AML-001, FR-CLM-AML-002 | Beneficial ownership verification |
| BR-CLM-OMB-001 | FR-CLM-OMB-001 | Complaint admissibility validation (Rule 14) |
| BR-CLM-OMB-002 | FR-CLM-OMB-002 | Jurisdiction mapping based on pincode |
| BR-CLM-OMB-003 | FR-CLM-OMB-002 | Conflict of interest screening |
| BR-CLM-OMB-004 | FR-CLM-OMB-003 | Mediation recommendation issuance |
| BR-CLM-OMB-005 | FR-CLM-OMB-004 | Award issuance with ₹50 lakh cap enforcement |
| BR-CLM-OMB-006 | FR-CLM-OMB-004 | Insurer compliance monitoring and IRDAI escalation |
| BR-CLM-OMB-007 | FR-CLM-OMB-004 | Complaint closure and archival with retention |
| BR-CLM-OMB-008 | FR-CLM-OMB-001, FR-CLM-OMB-002, FR-CLM-OMB-003, FR-CLM-OMB-004 | Bilingual support (English/Hindi) |
| BR-CLM-BOND-001 | FR-CLM-FL-001 | Freelook period calculation (15/30 days) |
| BR-CLM-BOND-002 | FR-CLM-BOND-001, FR-CLM-BOND-002 | Delivery failure escalation after 10 days |
| BR-CLM-BOND-003 | FR-CLM-FL-003 | Refund calculation with deductions |
| BR-CLM-BOND-004 | FR-CLM-FL-003 | Maker-checker workflow for refund processing |

---

### 10.3 Functional Requirements to Workflows Mapping

| Functional Requirement ID | Workflow ID | Workflow Step |
|---------------------------|-------------|---------------|
| FR-CLM-DC-001 | WF-CLM-DC-001 | Step 1: Claim Registration & Validation |
| FR-CLM-DC-002 | WF-CLM-DC-001 | Step 2: Document Capture & Indexing |
| FR-CLM-DC-003 | WF-CLM-DC-001 | Step 3: Investigation & Benefit Calculation |
| FR-CLM-DC-004 | WF-CLM-DC-001 | Step 4: Approval Workflow & Decision Points |
| FR-CLM-DC-005 | WF-CLM-DC-001 | Step 5: Disbursement & Payment Execution |
| FR-CLM-DC-006 | WF-CLM-DC-001 | Step 6: Reopen & Exception Handling |
| FR-CLM-DC-007 | WF-CLM-DC-001 | Step 7: Communication & Notifications |
| FR-CLM-DC-008 | WF-CLM-DC-001 | Step 8: Appeal Mechanism |
| FR-CLM-MC-001 | WF-CLM-MC-001 | Step 1: Maturity Report Generation (Batch) |
| FR-CLM-MC-002 | WF-CLM-MC-001 | Step 2: Multi-Channel Intimation |
| FR-CLM-MC-003 | WF-CLM-MC-001 | Step 3: Customer-Initiated Claim Submission |
| FR-CLM-MC-004 | WF-CLM-MC-001 | Step 4: System-Assisted Initial Scrutiny |
| FR-CLM-MC-005 | WF-CLM-MC-001 | Step 5: Auto-Indexing and Document Sync |
| FR-CLM-MC-006 | WF-CLM-MC-001 | Step 6: Auto-Populated Data Entry (OCR) |
| FR-CLM-MC-007 | WF-CLM-MC-001 | Step 7: QC Verification Checklist |
| FR-CLM-MC-008 | WF-CLM-MC-001 | Step 8: Approval Workflow and SLA Enforcement |
| FR-CLM-MC-009 | WF-CLM-MC-001 | Step 9: Auto-Generated Sanction/Rejection Communication |
| FR-CLM-MC-010 | WF-CLM-MC-001 | Step 10: Bank Account Validation (API-based) |
| FR-CLM-MC-011 | WF-CLM-MC-001 | Step 11: Disbursement Execution |
| FR-CLM-MC-012 | WF-CLM-MC-001 | Step 12: Voucher Generation and Submission |
| FR-CLM-MC-013 | WF-CLM-MC-001 | Step 13: Claim Closure and Archiving |
| FR-CLM-MC-014 | WF-CLM-MC-FEEDBACK-001 | Customer Feedback Collection (Post-Settlement) |
| FR-CLM-MC-015 | WF-CLM-MC-MONITORING-001 | Real-Time Monitoring & Escalation Dashboard |
| FR-CLM-MC-016 | WF-CLM-MC-TRACKER-001 | Customer Claim Tracker (Self-Service) |
| FR-CLM-SB-001 | WF-CLM-SB-001 | Step 1: SB Report Auto-Generation (Batch) |
| FR-CLM-SB-002 | WF-CLM-SB-001 | Step 2: Online Submission with DigiLocker |
| FR-CLM-SB-003 | WF-CLM-SB-001 | Step 3: Multi-Channel Intimation |
| FR-CLM-SB-004 | WF-CLM-SB-001 | Step 4: Initial Scrutiny (Digital) |
| FR-CLM-SB-005 | WF-CLM-SB-001 | Step 5: Auto-Indexing in IMS 2.0 |
| FR-CLM-SB-006 | WF-CLM-SB-001 | Step 6: Document Scanning & Upload |
| FR-CLM-SB-007 | WF-CLM-SB-001 | Step 7: Data Entry & QC Verification (Automated) |
| FR-CLM-SB-008 | WF-CLM-SB-001 | Step 8: Approval Workflow with Digital Signature |
| FR-CLM-SB-009 | WF-CLM-SB-001 | Step 9: Sanction/Rejection Letter Generation |
| FR-CLM-SB-010 | WF-CLM-SB-001 | Step 10: Bank Account Verification (API) |
| FR-CLM-SB-011 | WF-CLM-SB-001 | Step 11: Disbursement via Auto NEFT/IMPS |
| FR-CLM-SB-012 | WF-CLM-SB-001 | Step 12: Voucher Submission to Accounts |
| FR-CLM-SB-013 | WF-CLM-SB-FEEDBACK-001 | Customer Feedback Collection (Post-Settlement) |
| FR-CLM-SB-014 | WF-CLM-SB-MONITORING-001 | Monitoring & Escalation Dashboard |
| FR-CLM-SB-015 | WF-CLM-SB-TRACKER-001 | Customer Claim Tracker (Self-Service) |
| FR-CLM-AML-001 | WF-CLM-AML-001 | High Cash Premium Detection & CTR Filing |
| FR-CLM-AML-002 | WF-CLM-AML-001 | PAN Mismatch Detection & KYC Verification |
| FR-CLM-AML-003 | WF-CLM-AML-001 | Nominee Change Post Death Detection & Blocking |
| FR-CLM-AML-004 | WF-CLM-AML-001 | Frequent Surrender Pattern Detection |
| FR-CLM-AML-005 | WF-CLM-AML-001 | Refund Without Bond Delivery Detection |
| FR-CLM-OMB-001 | WF-CLM-OMB-001 | Complaint Intake & Registration |
| FR-CLM-OMB-002 | WF-CLM-OMB-001 | Jurisdiction Mapping & Conflict Screening |
| FR-CLM-OMB-003 | WF-CLM-OMB-001 | Hearing Scheduling & Management |
| FR-CLM-OMB-004 | WF-CLM-OMB-001 | Award Issuance & Enforcement |
| FR-CLM-BOND-001 | WF-CLM-BOND-001 | Policy Bond Dispatch Tracking (India Post API) |
| FR-CLM-BOND-002 | WF-CLM-BOND-001 | Delivery Confirmation & POD Capture |
| FR-CLM-FL-001 | WF-CLM-FL-001 | Freelook Period Monitoring & Countdown |
| FR-CLM-FL-002 | WF-CLM-FL-001 | Cancellation Request Handling & Validation |
| FR-CLM-FL-003 | WF-CLM-FL-001 | Refund Processing with Maker-Checker Workflow |

---

### 10.4 Workflows to Data Entities Mapping

| Workflow ID | Data Entities Used | Operations |
|-------------|-------------------|------------|
| WF-CLM-DC-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | INSERT, UPDATE, SELECT |
| WF-CLM-DC-INV-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | UPDATE, SELECT |
| WF-CLM-DC-APPEAL-001 | E-CLM-DC-001 (claims), E-CLM-APPEAL-001 (appeals - new entity) | INSERT, UPDATE, SELECT |
| WF-CLM-DC-REOPEN-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | UPDATE, SELECT |
| WF-CLM-MC-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | INSERT, UPDATE, SELECT |
| WF-CLM-MC-REPORT-001 | E-CLM-DC-001 (claims), E-CLM-REPORT-001 (reports - new entity) | SELECT, INSERT |
| WF-CLM-SB-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | INSERT, UPDATE, SELECT |
| WF-CLM-SB-REPORT-001 | E-CLM-DC-001 (claims), E-CLM-REPORT-001 (reports) | SELECT, INSERT |
| WF-CLM-AML-001 | E-CLM-AML-001 (aml_alerts), E-CLM-DC-001 (claims) | INSERT, UPDATE, SELECT |
| WF-CLM-AML-STR-001 | E-CLM-AML-001 (aml_alerts), E-CLM-FILING-001 (aml_filings - new entity) | SELECT, INSERT, UPDATE |
| WF-CLM-AML-CTR-001 | E-CLM-AML-001 (aml_alerts), E-CLM-FILING-001 (aml_filings) | SELECT, INSERT, UPDATE |
| WF-CLM-OMB-001 | E-CLM-OMB-001 (ombudsman_complaints - new entity), E-CLM-DC-001 (claims) | INSERT, UPDATE, SELECT |
| WF-CLM-OMB-HEARING-001 | E-CLM-OMB-001 (ombudsman_complaints), E-CLM-OMB-HEARING-001 (hearings - new entity) | INSERT, UPDATE, SELECT |
| WF-CLM-OMB-AWARD-001 | E-CLM-OMB-001 (ombudsman_complaints), E-CLM-OMB-AWARD-001 (awards - new entity) | INSERT, UPDATE, SELECT |
| WF-CLM-FL-001 | E-CLM-DC-001 (claims), E-CLM-DC-002 (claim_documents) | UPDATE, SELECT |
| WF-CLM-BOND-001 | E-CLM-BOND-001 (bond_delivery_tracking - new entity) | INSERT, UPDATE, SELECT |
| WF-CLM-MONITORING-001 | E-CLM-DC-001 (claims), E-CLM-SLA-ALERT-001 (sla_alerts - new entity) | SELECT, INSERT, UPDATE |
| WF-CLM-FEEDBACK-001 | E-CLM-FEEDBACK-001 (customer_feedback - new entity), E-CLM-DC-001 (claims) | INSERT, UPDATE, SELECT |

---

### 10.5 Data Entities to Integration Points Mapping

| Entity ID | Integration Point IDs | Purpose |
|-----------|----------------------|---------|
| E-CLM-DC-001 (claims) | INT-CLM-002 (Finacle/IT 2.0) | Payment execution |
| E-CLM-DC-001 (claims) | INT-CLM-004, INT-CLM-005, INT-CLM-006 | SMS/Email/WhatsApp notifications |
| E-CLM-DC-001 (claims) | INT-CLM-008 (Portal/Mobile App) | Self-service tracking |
| E-CLM-DC-002 (claim_documents) | INT-CLM-001 (ECMS) | Document storage/retrieval |
| E-CLM-DC-002 (claim_documents) | INT-CLM-007 (DigiLocker) | Digital document fetch |
| E-CLM-AML-001 (aml_alerts) | INT-CLM-009 (Finnet/FINGate) | STR/CTR submission |
| E-CLM-AML-001 (aml_alerts) | INT-CLM-010 (PAN Verification) | PAN validation |
| E-CLM-DC-001 (claims) | INT-CLM-011 (India Post Tracking) | Policy bond tracking |
| E-CLM-DC-001 (claims) | INT-CLM-012 (CPGRAMS) | Grievance integration |

---

### 10.6 Validation Rules to Error Codes Mapping

| Validation Rule ID | Error Code ID | Trigger Condition |
|--------------------|---------------|-------------------|
| VR-CLM-DC-001 | ERR-CLM-DC-003 | Death certificate missing |
| VR-CLM-DC-002 | ERR-CLM-DC-003 | Claim form missing |
| VR-CLM-DC-003 | ERR-CLM-DC-003 | Policy bond/indemnity missing |
| VR-CLM-DC-004 | ERR-CLM-DC-003 | Claimant ID proof missing |
| VR-CLM-DC-005 | ERR-CLM-DC-005 | Bank mandate missing or invalid |
| VR-CLM-DC-006 | ERR-CLM-DC-003 | Unnatural death documents missing |
| VR-CLM-DC-007 | ERR-CLM-DC-003 | Nomination documents missing |
| VR-CLM-MC-001 | ERR-CLM-MC-RJ-P-02 | Policy not active on maturity date |
| VR-CLM-MC-002 | ERR-CLM-MC-RJ-P-03 | Duplicate maturity claim detected |
| VR-CLM-MC-003 | ERR-CLM-MC-RJ-E-01 | Claimant identity mismatch |
| VR-CLM-MC-004 | ERR-CLM-MC-RJ-B-01 | Bank account validation failed |
| VR-CLM-MC-005 | ERR-CLM-MC-RJ-B-02 | Invalid IFSC code format |
| VR-CLM-MC-006 | ERR-CLM-MC-RJ-D-04 | Policy bond not submitted |
| VR-CLM-MC-007 | ERR-CLM-MC-RJ-D-05 | ID proof invalid or expired |
| VR-CLM-MC-008 | ERR-CLM-MC-RJ-D-02 | Document forgery detected |
| VR-CLM-SB-001 | ERR-CLM-SB-RJ-P-02 | Policy inactive on SB due date |
| VR-CLM-SB-002 | ERR-CLM-SB-RJ-P-03 | SB already paid |
| VR-CLM-SB-003 | ERR-CLM-SB-RJ-E-02 | Policyholder identity mismatch |
| VR-CLM-AML-001 | ERR-CLM-AML-001 (inferred) | Cash threshold exceeded |
| VR-CLM-AML-002 | ERR-CLM-AML-002 (inferred) | PAN verification failed |
| VR-CLM-AML-003 | ERR-CLM-AML-003 (inferred) | Nominee change after death detected |
| VR-CLM-AML-004 | ERR-CLM-AML-004 (inferred) | Surrender frequency alert |
| VR-CLM-AML-005 | ERR-CLM-AML-005 (inferred) | Refund before bond dispatch |
| VR-CLM-FL-001 | ERR-CLM-FL-001 (inferred) | Outside freelook window |
| VR-CLM-FL-002 | ERR-CLM-FL-002 (inferred) | Original bond not submitted |
| VR-CLM-FL-003 | ERR-CLM-FL-003 (inferred) | ID tamper detected |

---

### 10.7 Complete Cross-Reference Table

| Category | Count | SRS Source Documents |
|----------|-------|---------------------|
| Business Rules | 70 | 7 SRS documents (100% SRS-based) |
| Functional Requirements | 53 | Derived from SRS (100% SRS-based) |
| Validation Rules | 123 | Claim SRS documents (all types) - 100% coverage |
| Error Codes | 20+ | Claim SRS documents |
| Workflows | 18 | Death (4), Maturity (2), Survival Benefit (2), AML (3), Ombudsman (3), Bond/Freelook (2), Monitoring/Support (2) |
| Data Entities | 50+ | Claims, Documents, AML alerts, Ombudsman complaints, Bond tracking, Feedback, etc. |
| Integration Points | 22+ | ECMS, Finacle, Payment Gateways, Notifications, AML, DigiLocker, India Post, CPGRAMS, etc. |

---

### 10.8 SLA Summary

| Claim Type | Registration SLA | Approval SLA (No Investigation) | Approval SLA (With Investigation) | Penal Interest |
|------------|------------------|--------------------------------|----------------------------------|----------------|
| Death Claim | 48 hours | 15 days | 45 days | 8% p.a. |
| Maturity Claim | 48 hours | 7 days | N/A | 8% p.a. |
| Survival Benefit | 48 hours | 7 days | N/A | 8% p.a. |
| Investigation Report | N/A | Review within 5 days | Max 2 re-investigations (14 days each) | N/A |
| Document Deficiency | 7 days | N/A | N/A | N/A |
| Appeal Decision | N/A | N/A | 45 days from appeal submission | N/A |
| STR Filing | 7 working days from suspicion determination | N/A | N/A | N/A |
| CTR Filing | Monthly for cash >₹10L in one day | N/A | N/A | N/A |

---

### 10.9 Compliance Requirements Traceability

| Regulatory Body | Requirement | Business Rule ID | Functional Requirement ID | Workflow ID |
|----------------|-------------|------------------|---------------------------|-------------|
| FIU-India | STR filing within 7 working days | BR-CLM-AML-006 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004, FR-CLM-AML-005 | WF-CLM-AML-001 |
| FIU-India | Monthly CTR for cash >₹10L | BR-CLM-AML-007, BR-CLM-AML-008 | FR-CLM-AML-001 | WF-CLM-AML-001 |
| FIU-India | Regulatory reporting (STR/CTR/CCR/NTR) | BR-CLM-AML-010 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004, FR-CLM-AML-005 | WF-CLM-AML-001 |
| FIU-India | Negative list screening | BR-CLM-AML-011 | FR-CLM-AML-001, FR-CLM-AML-002, FR-CLM-AML-003, FR-CLM-AML-004 | WF-CLM-AML-001 |
| PMLA 2002 | Third-party PAN verification | BR-CLM-AML-009 | FR-CLM-AML-002 | WF-CLM-AML-001 |
| PMLA 2002 | Beneficial ownership verification | BR-CLM-AML-012 | FR-CLM-AML-001, FR-CLM-AML-002 | WF-CLM-AML-001 |
| IRDAI | Insurance Ombudsman Rules compliance | BR-CLM-OMB-001, BR-CLM-OMB-002, BR-CLM-OMB-003, BR-CLM-OMB-004, BR-CLM-OMB-005, BR-CLM-OMB-006 | FR-CLM-OMB-001, FR-CLM-OMB-002, FR-CLM-OMB-003, FR-CLM-OMB-004 | WF-CLM-OMB-001 |
| IRDAI | Freelook period compliance | BR-CLM-BOND-001 | FR-CLM-FL-001, FR-CLM-FL-002, FR-CLM-FL-003 | WF-CLM-FL-001 |
| PLI Act | Death claim investigation triggers | BR-CLM-DC-001, BR-CLM-DC-002 | FR-CLM-DC-003 | WF-CLM-DC-001 |
| PLI Act | Death claim approval timelines | BR-CLM-DC-003, BR-CLM-DC-004 | FR-CLM-DC-004 | WF-CLM-DC-001 |
| PLI Act | Appeal mechanism | BR-CLM-DC-005, BR-CLM-DC-006, BR-CLM-DC-007 | FR-CLM-DC-008 | WF-CLM-DC-001 |

---

### 10.10 Success Metrics and KPIs

| Metric | Target | Business Rule | Functional Requirement | Measurement |
|--------|--------|---------------|------------------------|-------------|
| Death Claim Approval (No Investigation) | 95% within 15 days | BR-CLM-DC-003 | FR-CLM-DC-004 | Days from registration to approval |
| Death Claim Approval (With Investigation) | 90% within 45 days | BR-CLM-DC-004 | FR-CLM-DC-004 | Days from registration to approval |
| Maturity Claim Approval | 98% within 7 days | BR-CLM-MC-002 | FR-CLM-MC-008 | Days from submission to approval |
| Survival Benefit Approval | 98% within 7 days | BR-CLM-SB-002 | FR-CLM-SB-008 | Days from submission to approval |
| Document Deficiency Auto-Reminder | 100% compliance | BR-CLM-MC-010 | FR-CLM-MC-004 | % of reminders sent on schedule |
| Investigation Report Submission | 100% within 21 days | BR-CLM-DC-002 | FR-CLM-DC-003 | Days from assignment to submission |
| Investigation Officer Assignment | 100% qualified personnel | BR-CLM-DC-011 | FR-CLM-DC-003 | % assigned to IP/ASP/PRI(P) ranks |
| AML Alert Detection | 100% detection | BR-CLM-AML-001 to 005 | FR-CLM-AML-001 to 005 | % of suspicious transactions flagged |
| STR Filing Compliance | 100% within 7 days | BR-CLM-AML-006 | FR-CLM-AML-001 to 005 | Days from suspicion to STR filing |
| CTR Filing Compliance | 100% monthly | BR-CLM-AML-007, BR-CLM-AML-008 | FR-CLM-AML-001 | % of monthly CTRs filed on time |
| Negative List Screening | 100% daily | BR-CLM-AML-011 | FR-CLM-AML-001 to 004 | Daily screening against OFAC/UN/UAPA/FATF |
| Ombudsman Award Compliance | 100% within 30 days | BR-CLM-OMB-006 | FR-CLM-OMB-004 | % of awards implemented within deadline |
| Freelook Period Accuracy | 100% compliance | BR-CLM-BOND-001 | FR-CLM-FL-001 | % of policies with correct freelook calculation |
| Bond Delivery Tracking | >95% delivered | BR-CLM-BOND-002 | FR-CLM-BOND-001, FR-CLM-BOND-002 | % of bonds delivered within 10 days |

---

**Legend**:
- **BR-CLM-**: Business Rule for Claims
- **FR-CLM-**: Functional Requirement for Claims
- **VR-CLM-**: Validation Rule for Claims
- **ERR-CLM-**: Error Code for Claims
- **WF-CLM-**: Workflow for Claims
- **E-CLM-**: Data Entity for Claims
- **INT-CLM-**: Integration Point for Claims

**Traceability Coverage**: 100% - All requirements from SRS documents are mapped to business rules, functional requirements, validation rules, error codes, workflows, data entities, and integration points.
