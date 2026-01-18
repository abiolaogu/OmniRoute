// OmniRoute Platform Owner Portal - Super Admin Dashboard
// Ultimate Platform Control Center with AI Insights

'use client';

import React, { useState, useEffect } from 'react';
import {
    Card, Row, Col, Statistic, Table, Typography, Space, Tag, Progress,
    Tabs, Timeline, Select, Button, Avatar, Badge, Tooltip, Alert,
    Modal, Input, Drawer, List, Segmented, Dropdown, notification
} from 'antd';
import {
    DashboardOutlined, DollarOutlined, ShoppingCartOutlined, UserOutlined,
    RiseOutlined, FallOutlined, BankOutlined, TruckOutlined, ShopOutlined,
    SettingOutlined, BellOutlined, ThunderboltOutlined, GlobalOutlined,
    SafetyOutlined, RobotOutlined, LineChartOutlined, TeamOutlined,
    CloudServerOutlined, ApiOutlined, DatabaseOutlined, EyeOutlined,
    WarningOutlined, CheckCircleOutlined, ClockCircleOutlined, SyncOutlined,
    FireOutlined, StarOutlined, CrownOutlined, AimOutlined, RadarChartOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text, Paragraph } = Typography;

// =============================================================================
// TYPES
// =============================================================================

interface PlatformMetrics {
    gmv: number;
    gmvChange: number;
    revenue: number;
    revenueChange: number;
    activeUsers: number;
    activeUsersChange: number;
    orders: number;
    ordersChange: number;
    avgOrderValue: number;
    transactions: number;
    successRate: number;
}

interface TenantHealth {
    id: string;
    name: string;
    type: 'manufacturer' | 'distributor' | 'retailer' | 'warehouse';
    status: 'healthy' | 'warning' | 'critical';
    gmv: number;
    orders: number;
    healthScore: number;
    issues: string[];
}

interface SystemHealth {
    service: string;
    status: 'operational' | 'degraded' | 'down';
    latency: number;
    uptime: number;
    requests: number;
}

interface AIInsight {
    id: string;
    type: 'opportunity' | 'risk' | 'trend' | 'anomaly';
    title: string;
    description: string;
    impact: 'high' | 'medium' | 'low';
    action?: string;
    createdAt: string;
}

// =============================================================================
// MOCK DATA
// =============================================================================

const platformMetrics: PlatformMetrics = {
    gmv: 2456789000,
    gmvChange: 18.5,
    revenue: 73703670,
    revenueChange: 22.3,
    activeUsers: 12456,
    activeUsersChange: 15.2,
    orders: 45678,
    ordersChange: 12.8,
    avgOrderValue: 53780,
    transactions: 67890,
    successRate: 98.7,
};

const tenantHealth: TenantHealth[] = [
    { id: '1', name: 'Nestle Nigeria', type: 'manufacturer', status: 'healthy', gmv: 450000000, orders: 2340, healthScore: 98, issues: [] },
    { id: '2', name: 'Dangote Foods', type: 'manufacturer', status: 'healthy', gmv: 380000000, orders: 1890, healthScore: 95, issues: [] },
    { id: '3', name: 'Chi Limited', type: 'manufacturer', status: 'warning', gmv: 210000000, orders: 1230, healthScore: 72, issues: ['High return rate', 'Payment delays'] },
    { id: '4', name: 'Lagos Mega Distribution', type: 'distributor', status: 'healthy', gmv: 320000000, orders: 4560, healthScore: 94, issues: [] },
    { id: '5', name: 'Kano Central Wholesale', type: 'distributor', status: 'critical', gmv: 89000000, orders: 890, healthScore: 45, issues: ['Credit overdue', 'Low fulfillment rate', 'Customer complaints'] },
];

const systemHealth: SystemHealth[] = [
    { service: 'Order Service', status: 'operational', latency: 45, uptime: 99.99, requests: 125000 },
    { service: 'Payment Gateway', status: 'operational', latency: 120, uptime: 99.98, requests: 89000 },
    { service: 'Inventory Service', status: 'operational', latency: 38, uptime: 99.97, requests: 234000 },
    { service: 'AI Gateway', status: 'degraded', latency: 890, uptime: 99.5, requests: 45000 },
    { service: 'Notification Service', status: 'operational', latency: 23, uptime: 99.99, requests: 567000 },
];

const aiInsights: AIInsight[] = [
    { id: '1', type: 'opportunity', title: 'Cross-sell Opportunity in Lagos', description: 'Retailers buying beverages have 78% probability of buying snacks. Bundle recommendation pending.', impact: 'high', action: 'Create Bundle', createdAt: '2026-01-18T10:00:00Z' },
    { id: '2', type: 'risk', title: 'Credit Risk Alert: Kano Region', description: '5 distributors showing payment delay patterns. Recommend credit limit review.', impact: 'high', action: 'Review Credits', createdAt: '2026-01-18T09:30:00Z' },
    { id: '3', type: 'trend', title: 'Rising Demand: Cooking Oil', description: 'Demand forecasting shows 34% increase in cooking oil orders next week. Stock optimization suggested.', impact: 'medium', createdAt: '2026-01-18T08:00:00Z' },
    { id: '4', type: 'anomaly', title: 'Unusual Order Pattern Detected', description: 'Retailer ID #4523 placing orders 5x normal volume. Possible fraud or bulk purchase.', impact: 'medium', action: 'Investigate', createdAt: '2026-01-18T07:45:00Z' },
];

// =============================================================================
// UTILITIES
// =============================================================================

const formatNaira = (amount: number): string => {
    if (amount >= 1000000000) {
        return `₦${(amount / 1000000000).toFixed(2)}B`;
    }
    if (amount >= 1000000) {
        return `₦${(amount / 1000000).toFixed(2)}M`;
    }
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
    }).format(amount);
};

const getStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
        healthy: 'success', operational: 'success',
        warning: 'warning', degraded: 'warning',
        critical: 'error', down: 'error',
    };
    return colors[status] || 'default';
};

const getInsightIcon = (type: AIInsight['type']): React.ReactNode => {
    const icons: Record<AIInsight['type'], React.ReactNode> = {
        opportunity: <StarOutlined style={{ color: '#52c41a' }} />,
        risk: <WarningOutlined style={{ color: '#ff4d4f' }} />,
        trend: <LineChartOutlined style={{ color: '#1890ff' }} />,
        anomaly: <AimOutlined style={{ color: '#fa8c16' }} />,
    };
    return icons[type];
};

// =============================================================================
// COMPONENTS
// =============================================================================

const GlowCard = ({ children, color = '#1890ff', glow = false }: { children: React.ReactNode; color?: string; glow?: boolean }) => (
    <Card
        style={{
            background: `linear-gradient(135deg, ${color}08 0%, ${color}15 100%)`,
            border: `1px solid ${color}30`,
            boxShadow: glow ? `0 0 20px ${color}30` : 'none',
            transition: 'all 0.3s ease',
        }}
        hoverable
    >
        {children}
    </Card>
);

const MetricCard = ({
    title,
    value,
    change,
    icon,
    color,
    prefix = '',
    suffix = ''
}: {
    title: string;
    value: string | number;
    change?: number;
    icon: React.ReactNode;
    color: string;
    prefix?: string;
    suffix?: string;
}) => (
    <GlowCard color={color}>
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Text type="secondary">{title}</Text>
                <div style={{
                    background: `${color}20`,
                    borderRadius: 8,
                    padding: 8,
                    display: 'flex'
                }}>
                    {icon}
                </div>
            </div>
            <Statistic
                value={value}
                prefix={prefix}
                suffix={suffix}
                valueStyle={{ fontSize: 28, fontWeight: 700, color }}
            />
            {change !== undefined && (
                <Space>
                    <Tag color={change >= 0 ? 'success' : 'error'} icon={change >= 0 ? <RiseOutlined /> : <FallOutlined />}>
                        {Math.abs(change)}%
                    </Tag>
                    <Text type="secondary" style={{ fontSize: 12 }}>vs last month</Text>
                </Space>
            )}
        </Space>
    </GlowCard>
);

const LiveTicker = () => {
    const [orders, setOrders] = useState(45678);

    useEffect(() => {
        const interval = setInterval(() => {
            setOrders(prev => prev + Math.floor(Math.random() * 3));
        }, 3000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div style={{
            background: 'linear-gradient(90deg, #1a365d 0%, #2d3748 100%)',
            padding: '8px 24px',
            display: 'flex',
            gap: 48,
            alignItems: 'center',
            color: 'white',
            overflow: 'hidden'
        }}>
            <Space>
                <Badge status="processing" />
                <Text style={{ color: 'white' }}>LIVE</Text>
            </Space>
            <Text style={{ color: 'rgba(255,255,255,0.9)' }}>
                🛒 Orders Today: <Text strong style={{ color: '#52c41a' }}>{orders.toLocaleString()}</Text>
            </Text>
            <Text style={{ color: 'rgba(255,255,255,0.9)' }}>
                💰 GMV: <Text strong style={{ color: '#ffd700' }}>{formatNaira(platformMetrics.gmv)}</Text>
            </Text>
            <Text style={{ color: 'rgba(255,255,255,0.9)' }}>
                👥 Active Users: <Text strong style={{ color: '#1890ff' }}>{platformMetrics.activeUsers.toLocaleString()}</Text>
            </Text>
            <Text style={{ color: 'rgba(255,255,255,0.9)' }}>
                ✨ Success Rate: <Text strong style={{ color: '#52c41a' }}>{platformMetrics.successRate}%</Text>
            </Text>
        </div>
    );
};

const AICommandCenter = ({ insights }: { insights: AIInsight[] }) => (
    <Card
        title={
            <Space>
                <RobotOutlined style={{ color: '#722ed1', fontSize: 20 }} />
                <span>AI Command Center</span>
                <Tag color="purple">Powered by Claude</Tag>
            </Space>
        }
        extra={<Button type="link">All Insights →</Button>}
    >
        <List
            itemLayout="horizontal"
            dataSource={insights}
            renderItem={(item) => (
                <List.Item
                    actions={item.action ? [
                        <Button type="primary" size="small" key="action">{item.action}</Button>
                    ] : []}
                >
                    <List.Item.Meta
                        avatar={
                            <Avatar
                                style={{
                                    background: item.impact === 'high' ? '#ff4d4f20' : item.impact === 'medium' ? '#fa8c1620' : '#52c41a20',
                                }}
                            >
                                {getInsightIcon(item.type)}
                            </Avatar>
                        }
                        title={
                            <Space>
                                <Text strong>{item.title}</Text>
                                <Tag color={item.impact === 'high' ? 'red' : item.impact === 'medium' ? 'orange' : 'green'}>
                                    {item.impact.toUpperCase()}
                                </Tag>
                            </Space>
                        }
                        description={item.description}
                    />
                </List.Item>
            )}
        />
    </Card>
);

const TenantHealthMatrix = ({ tenants }: { tenants: TenantHealth[] }) => {
    const columns: ColumnsType<TenantHealth> = [
        {
            title: 'Tenant',
            dataIndex: 'name',
            key: 'name',
            render: (name, record) => (
                <Space>
                    <Avatar style={{ background: getStatusColor(record.status) === 'success' ? '#52c41a' : getStatusColor(record.status) === 'warning' ? '#fa8c16' : '#ff4d4f' }}>
                        {name[0]}
                    </Avatar>
                    <div>
                        <Text strong>{name}</Text>
                        <br />
                        <Tag>{record.type}</Tag>
                    </div>
                </Space>
            ),
        },
        {
            title: 'Health Score',
            dataIndex: 'healthScore',
            key: 'healthScore',
            render: (score) => (
                <Tooltip title={`Score: ${score}/100`}>
                    <Progress
                        percent={score}
                        size="small"
                        status={score >= 80 ? 'success' : score >= 50 ? 'normal' : 'exception'}
                        format={(p) => `${p}`}
                    />
                </Tooltip>
            ),
        },
        {
            title: 'GMV',
            dataIndex: 'gmv',
            key: 'gmv',
            render: (v) => formatNaira(v),
            align: 'right',
        },
        {
            title: 'Orders',
            dataIndex: 'orders',
            key: 'orders',
            align: 'right',
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status, record) => (
                <Space direction="vertical" size="small">
                    <Tag color={getStatusColor(status)}>{status.toUpperCase()}</Tag>
                    {record.issues.length > 0 && (
                        <Tooltip title={record.issues.join(', ')}>
                            <Text type="danger" style={{ fontSize: 12 }}>
                                {record.issues.length} issue(s)
                            </Text>
                        </Tooltip>
                    )}
                </Space>
            ),
        },
        {
            title: 'Actions',
            key: 'actions',
            render: () => (
                <Space>
                    <Button size="small" icon={<EyeOutlined />} />
                    <Button size="small" icon={<SettingOutlined />} />
                </Space>
            ),
        },
    ];

    return (
        <Card
            title={
                <Space>
                    <RadarChartOutlined style={{ color: '#1890ff', fontSize: 20 }} />
                    <span>Tenant Health Matrix</span>
                </Space>
            }
            extra={<Button type="link">Manage Tenants →</Button>}
        >
            <Table
                columns={columns}
                dataSource={tenants}
                rowKey="id"
                pagination={false}
                size="small"
            />
        </Card>
    );
};

const SystemStatusPanel = ({ services }: { services: SystemHealth[] }) => (
    <Card
        title={
            <Space>
                <CloudServerOutlined style={{ color: '#52c41a', fontSize: 20 }} />
                <span>System Status</span>
                <Badge status="success" text="All Systems Operational" />
            </Space>
        }
    >
        <Row gutter={[8, 8]}>
            {services.map((svc) => (
                <Col span={24} key={svc.service}>
                    <div style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '8px 12px',
                        background: '#f5f5f5',
                        borderRadius: 8,
                        borderLeft: `3px solid ${svc.status === 'operational' ? '#52c41a' : svc.status === 'degraded' ? '#fa8c16' : '#ff4d4f'}`
                    }}>
                        <Space>
                            <Badge status={svc.status === 'operational' ? 'success' : svc.status === 'degraded' ? 'warning' : 'error'} />
                            <Text>{svc.service}</Text>
                        </Space>
                        <Space size="large">
                            <Text type="secondary">{svc.latency}ms</Text>
                            <Text type="secondary">{svc.uptime}%</Text>
                            <Text type="secondary">{(svc.requests / 1000).toFixed(0)}K req</Text>
                        </Space>
                    </div>
                </Col>
            ))}
        </Row>
    </Card>
);

// =============================================================================
// MAIN COMPONENT
// =============================================================================

export default function PlatformOwnerDashboard() {
    const [period, setPeriod] = useState<'today' | 'week' | 'month' | 'quarter'>('month');
    const [showAIAssistant, setShowAIAssistant] = useState(false);

    return (
        <div style={{ background: '#f0f2f5', minHeight: '100vh' }}>
            {/* Live Ticker */}
            <LiveTicker />

            {/* Header */}
            <div style={{ padding: '24px 24px 0' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                    <div>
                        <Space align="center">
                            <CrownOutlined style={{ fontSize: 32, color: '#d69e2e' }} />
                            <div>
                                <Title level={2} style={{ margin: 0 }}>Platform Command Center</Title>
                                <Text type="secondary">Complete oversight of the OmniRoute ecosystem</Text>
                            </div>
                        </Space>
                    </div>
                    <Space>
                        <Segmented
                            options={[
                                { label: 'Today', value: 'today' },
                                { label: 'Week', value: 'week' },
                                { label: 'Month', value: 'month' },
                                { label: 'Quarter', value: 'quarter' },
                            ]}
                            value={period}
                            onChange={(v) => setPeriod(v as any)}
                        />
                        <Button icon={<RobotOutlined />} type="primary" onClick={() => setShowAIAssistant(true)}>
                            AI Assistant
                        </Button>
                        <Badge count={5}>
                            <Button icon={<BellOutlined />} />
                        </Badge>
                    </Space>
                </div>

                {/* Primary Metrics */}
                <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                    <Col xs={24} sm={12} lg={6}>
                        <MetricCard
                            title="Gross Merchandise Value"
                            value={formatNaira(platformMetrics.gmv)}
                            change={platformMetrics.gmvChange}
                            icon={<DollarOutlined style={{ color: '#52c41a', fontSize: 20 }} />}
                            color="#52c41a"
                        />
                    </Col>
                    <Col xs={24} sm={12} lg={6}>
                        <MetricCard
                            title="Platform Revenue (3%)"
                            value={formatNaira(platformMetrics.revenue)}
                            change={platformMetrics.revenueChange}
                            icon={<BankOutlined style={{ color: '#1890ff', fontSize: 20 }} />}
                            color="#1890ff"
                        />
                    </Col>
                    <Col xs={24} sm={12} lg={6}>
                        <MetricCard
                            title="Active Users"
                            value={platformMetrics.activeUsers.toLocaleString()}
                            change={platformMetrics.activeUsersChange}
                            icon={<TeamOutlined style={{ color: '#722ed1', fontSize: 20 }} />}
                            color="#722ed1"
                        />
                    </Col>
                    <Col xs={24} sm={12} lg={6}>
                        <MetricCard
                            title="Total Orders"
                            value={platformMetrics.orders.toLocaleString()}
                            change={platformMetrics.ordersChange}
                            icon={<ShoppingCartOutlined style={{ color: '#fa8c16', fontSize: 20 }} />}
                            color="#fa8c16"
                        />
                    </Col>
                </Row>

                {/* AI Command Center & System Status */}
                <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                    <Col xs={24} lg={16}>
                        <AICommandCenter insights={aiInsights} />
                    </Col>
                    <Col xs={24} lg={8}>
                        <SystemStatusPanel services={systemHealth} />
                    </Col>
                </Row>

                {/* Tenant Health Matrix */}
                <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                    <Col span={24}>
                        <TenantHealthMatrix tenants={tenantHealth} />
                    </Col>
                </Row>

                {/* Quick Actions */}
                <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                    <Col span={24}>
                        <Card title="Quick Actions">
                            <Row gutter={[16, 16]}>
                                {[
                                    { icon: <UserOutlined />, label: 'Onboard Tenant', color: '#1890ff' },
                                    { icon: <SafetyOutlined />, label: 'Risk Review', color: '#ff4d4f' },
                                    { icon: <ApiOutlined />, label: 'API Settings', color: '#722ed1' },
                                    { icon: <DatabaseOutlined />, label: 'Data Export', color: '#52c41a' },
                                    { icon: <GlobalOutlined />, label: 'Region Config', color: '#fa8c16' },
                                    { icon: <ThunderboltOutlined />, label: 'Performance', color: '#13c2c2' },
                                ].map((action, i) => (
                                    <Col xs={12} sm={8} md={4} key={i}>
                                        <Button
                                            type="default"
                                            icon={action.icon}
                                            style={{ width: '100%', height: 80, flexDirection: 'column' }}
                                        >
                                            <div style={{ marginTop: 8 }}>{action.label}</div>
                                        </Button>
                                    </Col>
                                ))}
                            </Row>
                        </Card>
                    </Col>
                </Row>
            </div>

            {/* AI Assistant Drawer */}
            <Drawer
                title={
                    <Space>
                        <RobotOutlined />
                        <span>OmniRoute AI Assistant</span>
                    </Space>
                }
                open={showAIAssistant}
                onClose={() => setShowAIAssistant(false)}
                width={480}
            >
                <Space direction="vertical" style={{ width: '100%' }} size="large">
                    <Alert
                        message="AI-Powered Insights"
                        description="Ask me anything about your platform performance, tenant health, or get recommendations."
                        type="info"
                        showIcon
                    />
                    <Input.TextArea
                        placeholder="e.g., 'Which tenants need attention?' or 'What's driving revenue growth?'"
                        rows={4}
                    />
                    <Button type="primary" block icon={<ThunderboltOutlined />}>
                        Get AI Insights
                    </Button>
                    <Card size="small" title="Suggested Questions">
                        <Space direction="vertical" size="small">
                            <Button type="link" size="small">📊 Show me underperforming regions</Button>
                            <Button type="link" size="small">💰 Revenue forecast for next quarter</Button>
                            <Button type="link" size="small">⚠️ Identify high-risk credit accounts</Button>
                            <Button type="link" size="small">📈 Top growth opportunities</Button>
                        </Space>
                    </Card>
                </Space>
            </Drawer>
        </div>
    );
}
