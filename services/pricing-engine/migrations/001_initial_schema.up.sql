-- OmniRoute Pricing Engine Database Schema
-- Migration: 001_initial_schema
-- Description: Creates the core pricing tables

-- ============================================================================
-- Extensions
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================================
-- Tenant Table (Multi-tenancy)
-- ============================================================================

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- Customer Tables
-- ============================================================================

CREATE TYPE customer_type AS ENUM (
    'consumer',
    'retailer', 
    'wholesaler',
    'distributor',
    'enterprise'
);

CREATE TYPE customer_tier AS ENUM (
    'bronze',
    'silver',
    'gold',
    'platinum',
    'diamond'
);

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    external_id VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    type customer_type NOT NULL DEFAULT 'retailer',
    tier customer_tier NOT NULL DEFAULT 'bronze',
    email VARCHAR(255),
    phone VARCHAR(50),
    country VARCHAR(2),
    state VARCHAR(100),
    city VARCHAR(100),
    credit_limit DECIMAL(15,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, external_id)
);

CREATE INDEX idx_customers_tenant ON customers(tenant_id);
CREATE INDEX idx_customers_type ON customers(tenant_id, type);
CREATE INDEX idx_customers_tier ON customers(tenant_id, tier);

-- ============================================================================
-- Product Tables
-- ============================================================================

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES categories(id),
    level INTEGER DEFAULT 0,
    path VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_categories_parent ON categories(tenant_id, parent_id);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    brand VARCHAR(255),
    category_id UUID REFERENCES categories(id),
    base_price DECIMAL(15,4) NOT NULL,
    cost_price DECIMAL(15,4),
    currency VARCHAR(3) DEFAULT 'NGN',
    unit_of_measure VARCHAR(50) DEFAULT 'piece',
    tax_category VARCHAR(50) DEFAULT 'standard',
    weight_kg DECIMAL(10,4),
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, sku)
);

CREATE INDEX idx_products_tenant ON products(tenant_id);
CREATE INDEX idx_products_category ON products(tenant_id, category_id);
CREATE INDEX idx_products_brand ON products(tenant_id, brand);
CREATE INDEX idx_products_sku ON products(tenant_id, sku);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_adjustment DECIMAL(15,4) DEFAULT 0,
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_variants_product ON product_variants(product_id);

-- ============================================================================
-- Price List Tables
-- ============================================================================

CREATE TYPE pricing_method AS ENUM (
    'fixed',
    'discount_percent',
    'discount_amount',
    'margin'
);

CREATE TABLE price_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    priority INTEGER DEFAULT 0,
    customer_types customer_type[],
    customer_tiers customer_tier[],
    customer_ids UUID[],
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_price_lists_tenant ON price_lists(tenant_id);
CREATE INDEX idx_price_lists_valid ON price_lists(tenant_id, valid_from, valid_to);

CREATE TABLE price_list_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    price_list_id UUID NOT NULL REFERENCES price_lists(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    pricing_method pricing_method NOT NULL DEFAULT 'fixed',
    price DECIMAL(15,4),
    discount_percent DECIMAL(8,4),
    discount_amount DECIMAL(15,4),
    margin_percent DECIMAL(8,4),
    min_quantity INTEGER DEFAULT 1,
    max_quantity INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_price_list_items_list ON price_list_items(price_list_id);
CREATE INDEX idx_price_list_items_product ON price_list_items(product_id);

-- ============================================================================
-- Volume Discount Tables
-- ============================================================================

CREATE TABLE volume_discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    applies_to VARCHAR(50) NOT NULL DEFAULT 'all', -- all, product, category, brand
    product_ids UUID[],
    category_ids UUID[],
    brand_names VARCHAR(255)[],
    customer_types customer_type[],
    customer_tiers customer_tier[],
    can_combine BOOLEAN DEFAULT false,
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_volume_discounts_tenant ON volume_discounts(tenant_id);

CREATE TABLE volume_discount_tiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    volume_discount_id UUID NOT NULL REFERENCES volume_discounts(id) ON DELETE CASCADE,
    min_quantity INTEGER NOT NULL,
    max_quantity INTEGER,
    discount_percent DECIMAL(8,4),
    discount_amount DECIMAL(15,4),
    tier_order INTEGER DEFAULT 0
);

CREATE INDEX idx_volume_tiers_discount ON volume_discount_tiers(volume_discount_id);

-- ============================================================================
-- Contract Price Tables
-- ============================================================================

CREATE TABLE contract_prices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    product_id UUID NOT NULL REFERENCES products(id),
    contract_reference VARCHAR(100),
    price DECIMAL(15,4) NOT NULL,
    min_quantity INTEGER DEFAULT 1,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_to TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_contract_prices_customer ON contract_prices(customer_id);
CREATE INDEX idx_contract_prices_product ON contract_prices(product_id);
CREATE INDEX idx_contract_prices_valid ON contract_prices(valid_from, valid_to);

-- ============================================================================
-- Promotion Tables
-- ============================================================================

CREATE TYPE promotion_type AS ENUM (
    'discount',
    'bundle',
    'buy_x_get_y',
    'gift_with_purchase',
    'free_shipping'
);

CREATE TYPE discount_type AS ENUM (
    'percent',
    'fixed'
);

CREATE TABLE promotions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    type promotion_type NOT NULL,
    discount_type discount_type,
    discount_value DECIMAL(15,4),
    applies_to VARCHAR(50) NOT NULL DEFAULT 'all',
    product_ids UUID[],
    category_ids UUID[],
    brand_names VARCHAR(255)[],
    customer_types customer_type[],
    min_quantity INTEGER DEFAULT 1,
    min_order_value DECIMAL(15,2),
    max_uses INTEGER,
    current_uses INTEGER DEFAULT 0,
    can_combine BOOLEAN DEFAULT false,
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_promotions_tenant ON promotions(tenant_id);
CREATE INDEX idx_promotions_code ON promotions(tenant_id, code);
CREATE INDEX idx_promotions_valid ON promotions(tenant_id, valid_from, valid_to);

-- ============================================================================
-- Tax Tables
-- ============================================================================

CREATE TABLE tax_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'standard',
    rate DECIMAL(8,4) NOT NULL,
    country VARCHAR(2) NOT NULL,
    state VARCHAR(100),
    priority INTEGER DEFAULT 0,
    is_compound BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_to TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tax_rates_tenant ON tax_rates(tenant_id);
CREATE INDEX idx_tax_rates_location ON tax_rates(tenant_id, country, state);
CREATE INDEX idx_tax_rates_category ON tax_rates(tenant_id, category);

-- ============================================================================
-- Audit Log (for price change history)
-- ============================================================================

CREATE TABLE price_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    old_value JSONB,
    new_value JSONB,
    changed_by UUID,
    reason VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_price_history_entity ON price_history(entity_type, entity_id);
CREATE INDEX idx_price_history_created ON price_history(created_at);

-- ============================================================================
-- Trigger for updated_at
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_price_lists_updated_at BEFORE UPDATE ON price_lists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_volume_discounts_updated_at BEFORE UPDATE ON volume_discounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_contract_prices_updated_at BEFORE UPDATE ON contract_prices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_promotions_updated_at BEFORE UPDATE ON promotions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tax_rates_updated_at BEFORE UPDATE ON tax_rates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
