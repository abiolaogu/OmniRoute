-- OmniRoute Gig Platform Database Schema
-- Migration: 001_initial_schema
-- Description: Creates the core gig platform tables

-- ============================================================================
-- Extensions
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- ============================================================================
-- Worker Tables
-- ============================================================================

CREATE TYPE worker_status AS ENUM (
    'available',
    'busy',
    'offline',
    'on_break',
    'in_transit'
);

CREATE TYPE worker_type AS ENUM (
    'delivery',
    'sales',
    'collection',
    'audit',
    'multi_role'
);

CREATE TYPE worker_level AS ENUM (
    'starter',
    'bronze',
    'silver',
    'gold',
    'diamond',
    'master'
);

CREATE TABLE gig_workers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type worker_type NOT NULL DEFAULT 'delivery',
    level worker_level NOT NULL DEFAULT 'starter',
    status worker_status NOT NULL DEFAULT 'offline',
    
    -- Personal info
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    
    -- Location (PostGIS)
    current_location GEOGRAPHY(POINT, 4326),
    home_location GEOGRAPHY(POINT, 4326),
    current_address TEXT,
    
    -- Transport
    vehicle_type VARCHAR(50),
    vehicle_plate VARCHAR(50),
    
    -- Performance
    rating DECIMAL(3,2) DEFAULT 0,
    total_tasks INTEGER DEFAULT 0,
    completed_tasks INTEGER DEFAULT 0,
    success_rate DECIMAL(5,2) DEFAULT 0,
    
    -- Earnings
    total_earnings DECIMAL(15,2) DEFAULT 0,
    wallet_balance DECIMAL(15,2) DEFAULT 0,
    
    -- Verification
    is_verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, user_id)
);

CREATE INDEX idx_workers_tenant ON gig_workers(tenant_id);
CREATE INDEX idx_workers_status ON gig_workers(tenant_id, status);
CREATE INDEX idx_workers_type ON gig_workers(tenant_id, type);
CREATE INDEX idx_workers_location ON gig_workers USING GIST(current_location);

-- ============================================================================
-- Task Tables
-- ============================================================================

CREATE TYPE task_status AS ENUM (
    'pending',
    'assigned',
    'accepted',
    'in_progress',
    'completed',
    'failed',
    'cancelled'
);

CREATE TYPE task_type AS ENUM (
    'delivery',
    'pickup',
    'collection',
    'sales_visit',
    'audit',
    'custom'
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    type task_type NOT NULL DEFAULT 'delivery',
    status task_status NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    
    -- Assignment
    assigned_worker_id UUID REFERENCES gig_workers(id),
    assigned_at TIMESTAMP WITH TIME ZONE,
    
    -- Locations
    pickup_location GEOGRAPHY(POINT, 4326),
    pickup_address TEXT,
    dropoff_location GEOGRAPHY(POINT, 4326),
    dropoff_address TEXT,
    
    -- Order reference
    order_id UUID,
    customer_id UUID,
    customer_name VARCHAR(255),
    customer_phone VARCHAR(50),
    
    -- Task details
    description TEXT,
    instructions TEXT,
    items JSONB DEFAULT '[]',
    
    -- Collection
    collection_amount DECIMAL(15,2) DEFAULT 0,
    collected_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Pricing
    base_payout DECIMAL(10,2) DEFAULT 0,
    bonus_payout DECIMAL(10,2) DEFAULT 0,
    surge_multiplier DECIMAL(4,2) DEFAULT 1.0,
    
    -- Timing
    scheduled_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    deadline TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tasks_tenant ON tasks(tenant_id);
CREATE INDEX idx_tasks_status ON tasks(tenant_id, status);
CREATE INDEX idx_tasks_worker ON tasks(assigned_worker_id);
CREATE INDEX idx_tasks_order ON tasks(order_id);
CREATE INDEX idx_tasks_scheduled ON tasks(tenant_id, scheduled_at);
CREATE INDEX idx_tasks_pickup ON tasks USING GIST(pickup_location);
CREATE INDEX idx_tasks_dropoff ON tasks USING GIST(dropoff_location);

-- ============================================================================
-- Task Offers Table
-- ============================================================================

CREATE TYPE offer_status AS ENUM (
    'pending',
    'accepted',
    'rejected',
    'expired'
);

CREATE TABLE task_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    status offer_status NOT NULL DEFAULT 'pending',
    payout_amount DECIMAL(10,2) NOT NULL,
    offered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    responded_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(task_id, worker_id)
);

CREATE INDEX idx_offers_task ON task_offers(task_id);
CREATE INDEX idx_offers_worker ON task_offers(worker_id, status);
CREATE INDEX idx_offers_expires ON task_offers(expires_at) WHERE status = 'pending';

-- ============================================================================
-- Allocations Table
-- ============================================================================

CREATE TYPE allocation_status AS ENUM (
    'offered',
    'accepted',
    'rejected',
    'expired',
    'cancelled'
);

CREATE TABLE allocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    task_id UUID NOT NULL REFERENCES tasks(id),
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    status allocation_status NOT NULL DEFAULT 'offered',
    
    -- Payout details
    base_payout DECIMAL(10,2) DEFAULT 0,
    bonus_payout DECIMAL(10,2) DEFAULT 0,
    total_payout DECIMAL(10,2) DEFAULT 0,
    
    -- Timestamps
    allocated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    accepted_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_allocations_task ON allocations(task_id);
CREATE INDEX idx_allocations_worker ON allocations(worker_id);

-- ============================================================================
-- Earnings Table
-- ============================================================================

CREATE TABLE earnings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    task_id UUID REFERENCES tasks(id),
    type VARCHAR(50) NOT NULL, -- task_payout, bonus, tip, penalty
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- pending, credited, paid
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_earnings_worker ON earnings(worker_id);
CREATE INDEX idx_earnings_task ON earnings(task_id);
CREATE INDEX idx_earnings_status ON earnings(status);

-- ============================================================================
-- Payouts Table
-- ============================================================================

CREATE TABLE payouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    payment_method VARCHAR(50) NOT NULL,
    reference VARCHAR(100),
    processed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payouts_worker ON payouts(worker_id);
CREATE INDEX idx_payouts_status ON payouts(status);

-- ============================================================================
-- Task Proofs Table
-- ============================================================================

CREATE TABLE task_proofs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- photo, signature, receipt
    url TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_proofs_task ON task_proofs(task_id);

-- ============================================================================
-- Worker Location History
-- ============================================================================

CREATE TABLE worker_location_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    accuracy DECIMAL(10,2),
    speed DECIMAL(10,2),
    heading DECIMAL(5,2),
    battery_level INTEGER,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_location_history_worker ON worker_location_history(worker_id, recorded_at);
CREATE INDEX idx_location_history_location ON worker_location_history USING GIST(location);

-- ============================================================================
-- Routes Table (for optimized delivery routes)
-- ============================================================================

CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    worker_id UUID NOT NULL REFERENCES gig_workers(id),
    task_ids UUID[] NOT NULL,
    optimized_order INTEGER[],
    total_distance_km DECIMAL(10,2),
    estimated_duration_mins INTEGER,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_routes_worker ON routes(worker_id);
CREATE INDEX idx_routes_status ON routes(status);

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

CREATE TRIGGER update_workers_updated_at BEFORE UPDATE ON gig_workers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
