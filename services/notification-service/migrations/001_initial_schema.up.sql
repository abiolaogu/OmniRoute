-- OmniRoute Notification Service Database Schema
-- Migration: 001_initial_schema
-- Description: Creates the core notification tables

-- ============================================================================
-- Extensions
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================================
-- Notification Tables
-- ============================================================================

CREATE TYPE notification_channel AS ENUM (
    'whatsapp',
    'sms',
    'email',
    'push',
    'voice',
    'ussd'
);

CREATE TYPE notification_status AS ENUM (
    'pending',
    'queued',
    'sending',
    'sent',
    'delivered',
    'read',
    'failed',
    'cancelled'
);

CREATE TYPE notification_priority AS ENUM (
    'low',
    'normal',
    'high',
    'critical'
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    channel notification_channel NOT NULL,
    status notification_status NOT NULL DEFAULT 'pending',
    priority notification_priority NOT NULL DEFAULT 'normal',
    
    -- Recipient
    recipient_id UUID NOT NULL,
    recipient_type VARCHAR(50) NOT NULL, -- customer, worker, admin
    recipient_address VARCHAR(255) NOT NULL, -- phone, email, etc.
    
    -- Content
    template_id UUID,
    subject VARCHAR(500),
    body TEXT NOT NULL,
    data JSONB DEFAULT '{}',
    
    -- Metadata
    correlation_id VARCHAR(100),
    external_id VARCHAR(100),
    
    -- Timing
    scheduled_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    read_at TIMESTAMP WITH TIME ZONE,
    
    -- Retry
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    last_error TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_tenant ON notifications(tenant_id);
CREATE INDEX idx_notifications_recipient ON notifications(recipient_id);
CREATE INDEX idx_notifications_status ON notifications(tenant_id, status);
CREATE INDEX idx_notifications_channel ON notifications(tenant_id, channel);
CREATE INDEX idx_notifications_scheduled ON notifications(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_notifications_correlation ON notifications(correlation_id);

-- ============================================================================
-- Delivery Status Table (for webhooks and updates)
-- ============================================================================

CREATE TABLE delivery_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    status notification_status NOT NULL,
    provider_status VARCHAR(100),
    provider_message TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_delivery_notification ON delivery_status(notification_id);

-- ============================================================================
-- Template Tables
-- ============================================================================

CREATE TYPE template_type AS ENUM (
    'transactional',
    'marketing',
    'otp',
    'reminder',
    'alert'
);

CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    type template_type NOT NULL,
    channel notification_channel NOT NULL,
    
    -- Content
    subject VARCHAR(500),
    body TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    
    -- Localization
    language VARCHAR(10) DEFAULT 'en',
    translations JSONB DEFAULT '{}',
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_approved BOOLEAN DEFAULT false,
    approved_at TIMESTAMP WITH TIME ZONE,
    approved_by UUID,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, code, channel, language)
);

CREATE INDEX idx_templates_tenant ON templates(tenant_id);
CREATE INDEX idx_templates_code ON templates(tenant_id, code);
CREATE INDEX idx_templates_channel ON templates(tenant_id, channel);

-- ============================================================================
-- WhatsApp Templates
-- ============================================================================

CREATE TABLE whatsapp_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    external_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- AUTHENTICATION, MARKETING, UTILITY
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected
    components JSONB DEFAULT '[]',
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, name, language)
);

CREATE INDEX idx_wa_templates_tenant ON whatsapp_templates(tenant_id);
CREATE INDEX idx_wa_templates_status ON whatsapp_templates(tenant_id, status);

-- ============================================================================
-- USSD Session Tables
-- ============================================================================

CREATE TYPE session_status AS ENUM (
    'active',
    'ended',
    'timeout'
);

CREATE TABLE ussd_sessions (
    id VARCHAR(100) PRIMARY KEY,
    tenant_id UUID NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    service_code VARCHAR(20) NOT NULL,
    status session_status NOT NULL DEFAULT 'active',
    
    -- Navigation
    current_menu_id VARCHAR(100),
    menu_path VARCHAR(1000)[] DEFAULT '{}',
    
    -- State
    state JSONB DEFAULT '{}',
    
    -- Timestamps
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_ussd_sessions_phone ON ussd_sessions(phone_number);
CREATE INDEX idx_ussd_sessions_active ON ussd_sessions(status, last_active_at);

CREATE TABLE ussd_menus (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    options JSONB DEFAULT '[]',
    is_terminal BOOLEAN DEFAULT false,
    action VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_ussd_menus_tenant ON ussd_menus(tenant_id);

-- ============================================================================
-- Device Registration (for Push Notifications)
-- ============================================================================

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    platform VARCHAR(20) NOT NULL, -- ios, android, web
    app_version VARCHAR(50),
    device_model VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(token)
);

CREATE INDEX idx_devices_user ON devices(tenant_id, user_id);
CREATE INDEX idx_devices_token ON devices(token);
CREATE INDEX idx_devices_active ON devices(is_active, platform);

-- ============================================================================
-- Notification Preferences
-- ============================================================================

CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    channel notification_channel NOT NULL,
    enabled BOOLEAN DEFAULT true,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    quiet_hours_timezone VARCHAR(50),
    categories VARCHAR(100)[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, user_id, channel)
);

CREATE INDEX idx_prefs_user ON notification_preferences(tenant_id, user_id);

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

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_templates_updated_at BEFORE UPDATE ON templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wa_templates_updated_at BEFORE UPDATE ON whatsapp_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ussd_menus_updated_at BEFORE UPDATE ON ussd_menus FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_devices_updated_at BEFORE UPDATE ON devices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_prefs_updated_at BEFORE UPDATE ON notification_preferences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
