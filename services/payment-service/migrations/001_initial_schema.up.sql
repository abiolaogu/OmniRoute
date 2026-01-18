-- OmniRoute Payment Service Database Schema
-- Migration: 001_initial_schema
-- Description: Creates the core payment and credit tables

-- ============================================================================
-- Extensions
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Credit Limit Tables
-- ============================================================================

CREATE TYPE credit_tier AS ENUM (
    'premium',     -- 800-1000
    'standard',    -- 600-799
    'limited',     -- 400-599
    'restricted',  -- 200-399
    'no_credit'    -- <200
);

CREATE TABLE credit_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    
    -- Limit details
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    utilized_amount DECIMAL(15,2) DEFAULT 0,
    available_amount DECIMAL(15,2) GENERATED ALWAYS AS (amount - utilized_amount) STORED,
    
    -- Terms
    payment_terms_days INTEGER DEFAULT 30,
    
    -- Validity
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_to TIMESTAMP WITH TIME ZONE NOT NULL,
    review_date TIMESTAMP WITH TIME ZONE,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_frozen BOOLEAN DEFAULT false,
    freeze_reason TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, customer_id)
);

CREATE INDEX idx_credit_limits_customer ON credit_limits(customer_id);
CREATE INDEX idx_credit_limits_tenant ON credit_limits(tenant_id);

-- ============================================================================
-- Credit Scores Table (1000 point scale)
-- ============================================================================

CREATE TABLE credit_scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    
    -- Score breakdown (total 1000 points)
    total_score INTEGER NOT NULL,
    tier credit_tier NOT NULL,
    
    -- Transaction History (350 points)
    transaction_score INTEGER DEFAULT 0,
    
    -- Payment Behavior (350 points)
    payment_score INTEGER DEFAULT 0,
    
    -- Business Profile (200 points)
    business_score INTEGER DEFAULT 0,
    
    -- External Signals (100 points)
    external_score INTEGER DEFAULT 0,
    
    -- Score components JSON for detailed breakdown
    components JSONB DEFAULT '[]',
    
    -- Timestamps
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    valid_until TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_credit_scores_customer ON credit_scores(customer_id);
CREATE INDEX idx_credit_scores_tenant ON credit_scores(tenant_id);
CREATE INDEX idx_credit_scores_calculated ON credit_scores(calculated_at);

-- ============================================================================
-- Wallet Tables
-- ============================================================================

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    owner_type VARCHAR(50) NOT NULL, -- customer, worker
    
    -- Balance
    balance DECIMAL(15,2) DEFAULT 0,
    held_balance DECIMAL(15,2) DEFAULT 0,
    available_balance DECIMAL(15,2) GENERATED ALWAYS AS (balance - held_balance) STORED,
    currency VARCHAR(3) DEFAULT 'NGN',
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_frozen BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, owner_id, owner_type)
);

CREATE INDEX idx_wallets_owner ON wallets(owner_id, owner_type);
CREATE INDEX idx_wallets_tenant ON wallets(tenant_id);

CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    type VARCHAR(50) NOT NULL, -- credit, debit, hold, release
    amount DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,
    reference VARCHAR(100),
    description TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_wallet_txn_wallet ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_txn_created ON wallet_transactions(created_at);
CREATE INDEX idx_wallet_txn_reference ON wallet_transactions(reference);

-- ============================================================================
-- Payment Tables
-- ============================================================================

CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'cancelled',
    'refunded'
);

CREATE TYPE payment_method AS ENUM (
    'card',
    'bank_transfer',
    'mobile_money',
    'ussd',
    'wallet',
    'cash',
    'credit',
    'pos'
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    order_id UUID NOT NULL,
    
    -- Payment details
    method payment_method NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    
    -- Provider
    provider VARCHAR(50),
    provider_ref VARCHAR(100),
    
    -- Wallet deduction
    wallet_deducted DECIMAL(15,2) DEFAULT 0,
    
    -- Fallback tracking
    fallback_used BOOLEAN DEFAULT false,
    original_provider VARCHAR(50),
    
    -- Failure info
    failure_reason TEXT,
    failure_code VARCHAR(50),
    
    -- Timestamps
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_tenant ON payments(tenant_id);
CREATE INDEX idx_payments_customer ON payments(customer_id);
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_provider_ref ON payments(provider_ref);

-- ============================================================================
-- Invoice Tables
-- ============================================================================

CREATE TYPE invoice_status AS ENUM (
    'draft',
    'issued',
    'paid',
    'overdue',
    'cancelled',
    'write_off'
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    order_id UUID,
    
    -- Invoice details
    invoice_number VARCHAR(50) NOT NULL,
    status invoice_status NOT NULL DEFAULT 'draft',
    
    -- Amounts
    sub_total DECIMAL(15,2) NOT NULL,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    balance_due DECIMAL(15,2) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,
    currency VARCHAR(3) DEFAULT 'NGN',
    
    -- Dates
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    paid_date DATE,
    
    -- Items
    items JSONB DEFAULT '[]',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, invoice_number)
);

CREATE INDEX idx_invoices_tenant ON invoices(tenant_id);
CREATE INDEX idx_invoices_customer ON invoices(customer_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date) WHERE status IN ('issued', 'overdue');
CREATE INDEX idx_invoices_number ON invoices(invoice_number);

-- ============================================================================
-- Collection Task Tables
-- ============================================================================

CREATE TYPE collection_status AS ENUM (
    'pending',
    'assigned',
    'in_progress',
    'completed',
    'failed',
    'escalated'
);

CREATE TABLE collection_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    customer_id UUID NOT NULL,
    assigned_worker_id UUID,
    
    status collection_status NOT NULL DEFAULT 'pending',
    
    -- Collection details
    amount_due DECIMAL(15,2) NOT NULL,
    amount_collected DECIMAL(15,2) DEFAULT 0,
    collection_method VARCHAR(50), -- cash, transfer, pos, mobile_money
    
    -- Location
    location TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    
    -- Scheduling
    scheduled_date DATE NOT NULL,
    attempt_count INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Notes
    notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_collections_tenant ON collection_tasks(tenant_id);
CREATE INDEX idx_collections_invoice ON collection_tasks(invoice_id);
CREATE INDEX idx_collections_worker ON collection_tasks(assigned_worker_id);
CREATE INDEX idx_collections_status ON collection_tasks(status);
CREATE INDEX idx_collections_scheduled ON collection_tasks(scheduled_date);

-- ============================================================================
-- Payment Plan Tables
-- ============================================================================

CREATE TABLE payment_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    invoice_ids UUID[] NOT NULL,
    
    -- Plan details
    total_amount DECIMAL(15,2) NOT NULL,
    installments JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(50) DEFAULT 'pending_approval', -- pending_approval, active, completed, defaulted
    
    -- Approval
    approved_by UUID,
    approved_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_plans_tenant ON payment_plans(tenant_id);
CREATE INDEX idx_plans_customer ON payment_plans(customer_id);
CREATE INDEX idx_plans_status ON payment_plans(status);

-- ============================================================================
-- Settlement Tables
-- ============================================================================

CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    
    -- Settlement details
    gross_amount DECIMAL(15,2) NOT NULL,
    fees DECIMAL(15,2) DEFAULT 0,
    net_amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    
    -- Status
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    
    -- Bank details
    bank_code VARCHAR(10),
    account_number VARCHAR(20),
    account_name VARCHAR(255),
    
    -- Reference
    reference VARCHAR(100),
    
    -- Timestamps
    scheduled_for TIMESTAMP WITH TIME ZONE NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_settlements_tenant ON settlements(tenant_id);
CREATE INDEX idx_settlements_merchant ON settlements(merchant_id);
CREATE INDEX idx_settlements_status ON settlements(status);
CREATE INDEX idx_settlements_scheduled ON settlements(scheduled_for);

-- ============================================================================
-- Credit History (for scoring)
-- ============================================================================

CREATE TABLE credit_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- invoice_paid, invoice_overdue, credit_increase, etc.
    amount DECIMAL(15,2),
    days_late INTEGER,
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_credit_history_customer ON credit_history(customer_id);
CREATE INDEX idx_credit_history_type ON credit_history(event_type);
CREATE INDEX idx_credit_history_created ON credit_history(created_at);

-- ============================================================================
-- Triggers
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_credit_limits_updated_at BEFORE UPDATE ON credit_limits FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_collections_updated_at BEFORE UPDATE ON collection_tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_plans_updated_at BEFORE UPDATE ON payment_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
