-- Fleet Service Database Schema
-- Vehicles, Drivers, Telemetry, Maintenance, Fuel, Geofencing

-- Vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    registration_number VARCHAR(50) NOT NULL,
    vin VARCHAR(50),
    make VARCHAR(100),
    model VARCHAR(100),
    year INT,
    vehicle_type VARCHAR(30) CHECK (vehicle_type IN ('motorcycle', 'car', 'van', 'truck', 'pickup')),
    fuel_type VARCHAR(20) CHECK (fuel_type IN ('petrol', 'diesel', 'electric', 'hybrid', 'cng')),
    capacity_kg DECIMAL(10,2),
    capacity_volume_m3 DECIMAL(10,2),
    
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'in_use', 'maintenance', 'inactive')),
    assigned_driver_id UUID,
    
    current_latitude DECIMAL(10,8),
    current_longitude DECIMAL(11,8),
    last_location_at TIMESTAMPTZ,
    
    odometer_km DECIMAL(12,2) DEFAULT 0,
    purchase_date DATE,
    purchase_price DECIMAL(15,2),
    insurance_expiry DATE,
    inspection_expiry DATE,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, registration_number)
);

CREATE INDEX idx_vehicles_tenant ON vehicles(tenant_id);
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_vehicles_driver ON vehicles(assigned_driver_id);

-- Drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID, -- Link to auth service user
    driver_number VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    license_number VARCHAR(100) NOT NULL,
    license_class VARCHAR(20),
    license_expiry DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    
    -- Performance
    safety_score DECIMAL(5,2) DEFAULT 100,
    total_trips INT DEFAULT 0,
    total_distance_km DECIMAL(12,2) DEFAULT 0,
    
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(50),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, driver_number)
);

CREATE INDEX idx_drivers_tenant ON drivers(tenant_id);
CREATE INDEX idx_drivers_status ON drivers(status);

-- Telemetry data (time-series)
CREATE TABLE IF NOT EXISTS vehicle_telemetry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    altitude_m DECIMAL(8,2),
    speed_kmh DECIMAL(6,2),
    heading_degrees DECIMAL(5,2),
    
    engine_rpm INT,
    fuel_level_percent DECIMAL(5,2),
    odometer_km DECIMAL(12,2),
    engine_temperature_c DECIMAL(5,2),
    
    ignition_on BOOLEAN,
    doors_locked BOOLEAN,
    
    raw_data JSONB
);

CREATE INDEX idx_telemetry_vehicle_time ON vehicle_telemetry(vehicle_id, timestamp DESC);

-- Trips table
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    driver_id UUID REFERENCES drivers(id),
    
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'cancelled')),
    purpose VARCHAR(50),
    
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    start_location JSONB,
    end_location JSONB,
    
    distance_km DECIMAL(10,2),
    duration_minutes INT,
    idle_time_minutes INT,
    max_speed_kmh DECIMAL(6,2),
    fuel_used_liters DECIMAL(8,2),
    
    route_polyline TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_trips_vehicle ON trips(vehicle_id);
CREATE INDEX idx_trips_driver ON trips(driver_id);
CREATE INDEX idx_trips_start ON trips(start_time DESC);

-- Maintenance records
CREATE TABLE IF NOT EXISTS maintenance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    
    maintenance_type VARCHAR(30) CHECK (maintenance_type IN (
        'oil_change', 'tire_rotation', 'brake_service', 'inspection',
        'repair', 'scheduled', 'unscheduled'
    )),
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled')),
    
    scheduled_date DATE,
    completed_date DATE,
    odometer_at_service DECIMAL(12,2),
    
    description TEXT,
    vendor VARCHAR(255),
    cost DECIMAL(15,2),
    invoice_number VARCHAR(100),
    
    next_service_date DATE,
    next_service_km DECIMAL(12,2),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_maintenance_vehicle ON maintenance_records(vehicle_id);
CREATE INDEX idx_maintenance_status ON maintenance_records(status);

-- Fuel records
CREATE TABLE IF NOT EXISTS fuel_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES drivers(id),
    
    fuel_type VARCHAR(20),
    quantity_liters DECIMAL(8,2) NOT NULL,
    price_per_liter DECIMAL(10,2),
    total_cost DECIMAL(15,2),
    
    odometer_km DECIMAL(12,2),
    is_full_tank BOOLEAN DEFAULT TRUE,
    station_name VARCHAR(255),
    station_location JSONB,
    
    receipt_url TEXT,
    notes TEXT,
    
    refuel_date TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_fuel_vehicle ON fuel_records(vehicle_id);
CREATE INDEX idx_fuel_date ON fuel_records(refuel_date DESC);

-- Geofences
CREATE TABLE IF NOT EXISTS geofences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    geofence_type VARCHAR(20) DEFAULT 'polygon' CHECK (geofence_type IN ('polygon', 'circle')),
    
    -- For circle
    center_latitude DECIMAL(10,8),
    center_longitude DECIMAL(11,8),
    radius_meters DECIMAL(10,2),
    
    -- For polygon (GeoJSON)
    boundary JSONB,
    
    alert_on_enter BOOLEAN DEFAULT TRUE,
    alert_on_exit BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_geofences_tenant ON geofences(tenant_id);

-- Geofence violations
CREATE TABLE IF NOT EXISTS geofence_violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    geofence_id UUID NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    driver_id UUID REFERENCES drivers(id),
    violation_type VARCHAR(10) CHECK (violation_type IN ('enter', 'exit')),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_violations_geofence ON geofence_violations(geofence_id);
CREATE INDEX idx_violations_vehicle ON geofence_violations(vehicle_id);
