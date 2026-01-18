-- Analytics Service Database Schema
-- Metrics, KPIs, Reports, Dashboards, Alerts

-- Metrics definitions
CREATE TABLE IF NOT EXISTS metric_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    unit VARCHAR(20),
    aggregation_type VARCHAR(20) CHECK (aggregation_type IN ('sum', 'avg', 'count', 'min', 'max', 'last')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Daily metrics (aggregated)
CREATE TABLE IF NOT EXISTS daily_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    metric_id UUID NOT NULL REFERENCES metric_definitions(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    value DECIMAL(20,4) NOT NULL,
    dimensions JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, metric_id, date, dimensions)
);

CREATE INDEX idx_daily_metrics_date ON daily_metrics(tenant_id, date DESC);
CREATE INDEX idx_daily_metrics_metric ON daily_metrics(metric_id);

-- KPI definitions
CREATE TABLE IF NOT EXISTS kpi_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    formula TEXT, -- SQL or calculation formula
    target_value DECIMAL(20,4),
    target_direction VARCHAR(10) CHECK (target_direction IN ('increase', 'decrease', 'maintain')),
    unit VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- KPI values
CREATE TABLE IF NOT EXISTS kpi_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kpi_id UUID NOT NULL REFERENCES kpi_definitions(id) ON DELETE CASCADE,
    period_type VARCHAR(20) CHECK (period_type IN ('daily', 'weekly', 'monthly', 'quarterly', 'yearly')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    value DECIMAL(20,4) NOT NULL,
    target DECIMAL(20,4),
    variance DECIMAL(20,4),
    variance_percent DECIMAL(10,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(kpi_id, period_type, period_start)
);

CREATE INDEX idx_kpi_values_period ON kpi_values(kpi_id, period_start DESC);

-- Report definitions
CREATE TABLE IF NOT EXISTS report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    report_type VARCHAR(30) CHECK (report_type IN ('sales', 'inventory', 'customer', 'worker', 'finance', 'custom')),
    config JSONB NOT NULL, -- Report configuration (columns, filters, etc.)
    is_scheduled BOOLEAN DEFAULT FALSE,
    schedule_cron VARCHAR(100),
    recipients TEXT[],
    format VARCHAR(10) DEFAULT 'pdf' CHECK (format IN ('pdf', 'excel', 'csv', 'json')),
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reports_tenant ON report_definitions(tenant_id);

-- Generated reports
CREATE TABLE IF NOT EXISTS generated_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_definition_id UUID REFERENCES report_definitions(id) ON DELETE SET NULL,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'generating' CHECK (status IN ('generating', 'completed', 'failed')),
    file_url TEXT,
    file_size_bytes BIGINT,
    row_count INT,
    error_message TEXT,
    generated_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_generated_reports_tenant ON generated_reports(tenant_id);
CREATE INDEX idx_generated_reports_created ON generated_reports(created_at DESC);

-- Dashboard definitions
CREATE TABLE IF NOT EXISTS dashboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    layout JSONB NOT NULL, -- Widget positions and sizes
    is_default BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Dashboard widgets
CREATE TABLE IF NOT EXISTS dashboard_widgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    widget_type VARCHAR(30) CHECK (widget_type IN ('chart', 'number', 'table', 'gauge', 'map', 'list')),
    title VARCHAR(255) NOT NULL,
    config JSONB NOT NULL, -- Chart type, data source, colors, etc.
    position_x INT DEFAULT 0,
    position_y INT DEFAULT 0,
    width INT DEFAULT 4,
    height INT DEFAULT 3,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Analytics alerts
CREATE TABLE IF NOT EXISTS analytics_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metric_id UUID REFERENCES metric_definitions(id),
    kpi_id UUID REFERENCES kpi_definitions(id),
    condition_type VARCHAR(20) CHECK (condition_type IN ('above', 'below', 'equals', 'change_percent')),
    threshold_value DECIMAL(20,4) NOT NULL,
    severity VARCHAR(10) CHECK (severity IN ('info', 'warning', 'critical')),
    notification_channels TEXT[],
    is_active BOOLEAN DEFAULT TRUE,
    last_triggered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Alert history
CREATE TABLE IF NOT EXISTS alert_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id UUID NOT NULL REFERENCES analytics_alerts(id) ON DELETE CASCADE,
    triggered_value DECIMAL(20,4) NOT NULL,
    message TEXT,
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_by UUID,
    acknowledged_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_alert_history_alert ON alert_history(alert_id);
CREATE INDEX idx_alert_history_created ON alert_history(created_at DESC);
