-- Inventory Service Database Schema
-- Stock levels, Reservations, Movements, Batches, Stocktakes

-- Warehouses table
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_warehouses_tenant ON warehouses(tenant_id);

-- Warehouse locations (zones, aisles, racks, bins)
CREATE TABLE IF NOT EXISTS warehouse_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    zone VARCHAR(50),
    aisle VARCHAR(20),
    rack VARCHAR(20),
    level VARCHAR(20),
    bin VARCHAR(20),
    location_type VARCHAR(20) CHECK (location_type IN ('receiving', 'storage', 'picking', 'packing', 'shipping')),
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(warehouse_id, code)
);

CREATE INDEX idx_locations_warehouse ON warehouse_locations(warehouse_id);

-- Stock levels
CREATE TABLE IF NOT EXISTS stock_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    location_id UUID REFERENCES warehouse_locations(id),
    product_id UUID NOT NULL,
    variant_id UUID,
    quantity_on_hand INT NOT NULL DEFAULT 0,
    quantity_reserved INT DEFAULT 0,
    quantity_available INT GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
    reorder_point INT,
    reorder_quantity INT,
    last_counted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(warehouse_id, product_id, variant_id, location_id)
);

CREATE INDEX idx_stock_tenant ON stock_levels(tenant_id);
CREATE INDEX idx_stock_warehouse ON stock_levels(warehouse_id);
CREATE INDEX idx_stock_product ON stock_levels(product_id);
CREATE INDEX idx_stock_low ON stock_levels(tenant_id, quantity_on_hand) WHERE quantity_on_hand <= reorder_point;

-- Stock reservations
CREATE TABLE IF NOT EXISTS stock_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    stock_level_id UUID NOT NULL REFERENCES stock_levels(id) ON DELETE CASCADE,
    order_id UUID,
    order_item_id UUID,
    quantity INT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'released', 'cancelled')),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reservations_order ON stock_reservations(order_id);
CREATE INDEX idx_reservations_status ON stock_reservations(status);

-- Stock movements
CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL,
    variant_id UUID,
    movement_type VARCHAR(30) CHECK (movement_type IN (
        'purchase_receive', 'sale_dispatch', 'transfer_in', 'transfer_out',
        'adjustment_add', 'adjustment_remove', 'return_receive', 'damage_write_off'
    )),
    quantity INT NOT NULL,
    from_location_id UUID REFERENCES warehouse_locations(id),
    to_location_id UUID REFERENCES warehouse_locations(id),
    reference_type VARCHAR(30),
    reference_id UUID,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_movements_tenant ON stock_movements(tenant_id);
CREATE INDEX idx_movements_product ON stock_movements(product_id);
CREATE INDEX idx_movements_created ON stock_movements(created_at DESC);

-- Batches (for lot tracking)
CREATE TABLE IF NOT EXISTS batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL,
    batch_number VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    quantity_remaining INT NOT NULL,
    manufactured_date DATE,
    expiry_date DATE,
    supplier_id UUID,
    purchase_price DECIMAL(15,2),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, product_id, batch_number)
);

CREATE INDEX idx_batches_product ON batches(product_id);
CREATE INDEX idx_batches_expiry ON batches(expiry_date);

-- Stocktakes
CREATE TABLE IF NOT EXISTS stocktakes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    stocktake_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'in_progress', 'completed', 'cancelled')),
    stocktake_type VARCHAR(20) CHECK (stocktake_type IN ('full', 'cycle', 'spot')),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stocktakes_warehouse ON stocktakes(warehouse_id);

-- Stocktake counts
CREATE TABLE IF NOT EXISTS stocktake_counts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stocktake_id UUID NOT NULL REFERENCES stocktakes(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    variant_id UUID,
    location_id UUID REFERENCES warehouse_locations(id),
    expected_quantity INT NOT NULL,
    counted_quantity INT,
    variance INT GENERATED ALWAYS AS (counted_quantity - expected_quantity) STORED,
    counted_by UUID,
    counted_at TIMESTAMPTZ,
    notes TEXT
);

CREATE INDEX idx_counts_stocktake ON stocktake_counts(stocktake_id);
