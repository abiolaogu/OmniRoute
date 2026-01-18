// OmniRoute Partner Portal - Comprehensive B2B Multi-Role Dashboard
// Manufacturer, Distributor, Wholesaler, Retailer, Logistics, Warehouse Portals

'use client';

import React, { useState, useEffect } from 'react';
import {
    Card, Row, Col, Statistic, Table, Typography, Space, Tag, Progress,
    Tabs, Timeline, Select, Button, Avatar, Badge, Tooltip, Alert,
    Steps, Segmented, Drawer, List, Empty, Divider, Spin, Switch, Rate,
    notification, Modal, Menu, Dropdown, Input
} from 'antd';
import {
    ShopOutlined, DollarOutlined, ShoppingCartOutlined, UserOutlined,
    RiseOutlined, FallOutlined, TruckOutlined, InboxOutlined, BarChartOutlined,
    PlusOutlined, BellOutlined, SyncOutlined, SearchOutlined, FilterOutlined,
    ExportOutlined, QrcodeOutlined, EnvironmentOutlined, ClockCircleOutlined,
    CheckCircleOutlined, WarningOutlined, ThunderboltOutlined, TeamOutlined,
    StockOutlined, FundOutlined, LineChartOutlined, FireOutlined, StarOutlined,
    SendOutlined, EyeOutlined, EditOutlined, DeleteOutlined, SwapOutlined,
    BuildOutlined, CarOutlined, AimOutlined, RadarChartOutlined, RobotOutlined,
    BankOutlined, CreditCardOutlined, GoldOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

// =============================================================================
// TYPES
// =============================================================================

type PartnerRole = 'manufacturer' | 'distributor' | 'wholesaler' | 'retailer' | 'logistics' | 'warehouse';

interface RoleConfig {
    key: PartnerRole;
    name: string;
    icon: React.ReactNode;
    color: string;
    metrics: string[];
}

interface Order {
    id: string;
    orderNumber: string;
    type: 'inbound' | 'outbound';
    counterparty: string;
    items: number;
    total: number;
    status: 'draft' | 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled';
    createdAt: string;
    estimatedDelivery?: string;
}

interface InventoryItem {
    id: string;
    sku: string;
    name: string;
    category: string;
    quantity: number;
    reserved: number;
    available: number;
    reorderPoint: number;
    unitPrice: number;
    location: string;
    lastMovement: string;
    status: 'healthy' | 'low' | 'critical' | 'overstock';
    velocity: 'fast' | 'medium' | 'slow';
}

interface Customer {
    id: string;
    name: string;
    type: 'retailer' | 'wholesaler' | 'consumer';
    tier: 'standard' | 'silver' | 'gold' | 'platinum';
    totalOrders: number;
    totalSpent: number;
    creditLimit: number;
    creditUsed: number;
    lastOrder: string;
    rating: number;
}

interface Route {
    id: string;
    name: string;
    stops: number;
    totalValue: number;
    assignedAgent?: string;
    status: 'pending' | 'in_progress' | 'completed';
    completionRate: number;
}

// =============================================================================
// ROLE CONFIGURATIONS
// =============================================================================

const roleConfigs: RoleConfig[] = [
    { key: 'manufacturer', name: 'Manufacturer', icon: <BuildOutlined />, color: '#722ed1', metrics: ['production', 'distribution', 'inventory'] },
    { key: 'distributor', name: 'Distributor', icon: <TruckOutlined />, color: '#1890ff', metrics: ['sales', 'routes', 'credit'] },
    { key: 'wholesaler', name: 'Wholesaler', icon: <ShopOutlined />, color: '#52c41a', metrics: ['orders', 'inventory', 'pricing'] },
    { key: 'retailer', name: 'Retailer', icon: <ShoppingCartOutlined />, color: '#fa8c16', metrics: ['sales', 'inventory', 'restock'] },
    { key: 'logistics', name: '3PL Provider', icon: <CarOutlined />, color: '#13c2c2', metrics: ['trips', 'fleet', 'pod'] },
    { key: 'warehouse', name: 'Warehouse', icon: <InboxOutlined />, color: '#eb2f96', metrics: ['storage', 'fulfillment', 'wms'] },
];

// =============================================================================
// MOCK DATA
// =============================================================================

const orders: Order[] = [
    { id: '1', orderNumber: 'ORD-2026-0145', type: 'outbound', counterparty: 'Lagos Mega Retailers', items: 45, total: 2340000, status: 'confirmed', createdAt: '2026-01-18T10:30:00Z', estimatedDelivery: '2026-01-19' },
    { id: '2', orderNumber: 'ORD-2026-0144', type: 'inbound', counterparty: 'Nestle Nigeria PLC', items: 120, total: 8900000, status: 'processing', createdAt: '2026-01-18T09:00:00Z', estimatedDelivery: '2026-01-20' },
    { id: '3', orderNumber: 'ORD-2026-0143', type: 'outbound', counterparty: 'Kano Central Market', items: 28, total: 1560000, status: 'shipped', createdAt: '2026-01-17T14:00:00Z' },
    { id: '4', orderNumber: 'ORD-2026-0142', type: 'outbound', counterparty: 'Ibadan Wholesale Hub', items: 65, total: 3450000, status: 'delivered', createdAt: '2026-01-17T11:00:00Z' },
];

const inventory: InventoryItem[] = [
    { id: '1', sku: 'PM-400-CTN', name: 'Peak Milk 400g (Carton)', category: 'Dairy', quantity: 450, reserved: 120, available: 330, reorderPoint: 200, unitPrice: 2500, location: 'A-01-03', lastMovement: '2h ago', status: 'healthy', velocity: 'fast' },
    { id: '2', sku: 'IND-CHK-40', name: 'Indomie Chicken (Carton of 40)', category: 'Noodles', quantity: 85, reserved: 40, available: 45, reorderPoint: 100, unitPrice: 5500, location: 'A-02-01', lastMovement: '4h ago', status: 'low', velocity: 'fast' },
    { id: '3', sku: 'GPR-50KG', name: 'Golden Penny Rice 50kg', category: 'Grains', quantity: 25, reserved: 15, available: 10, reorderPoint: 50, unitPrice: 68000, location: 'B-01-01', lastMovement: '1h ago', status: 'critical', velocity: 'fast' },
    { id: '4', sku: 'DSG-50KG', name: 'Dangote Sugar 50kg', category: 'Sugar', quantity: 320, reserved: 80, available: 240, reorderPoint: 100, unitPrice: 48000, location: 'B-02-03', lastMovement: '6h ago', status: 'overstock', velocity: 'medium' },
];

const customers: Customer[] = [
    { id: '1', name: 'Mama Ngozi Stores', type: 'retailer', tier: 'gold', totalOrders: 156, totalSpent: 12500000, creditLimit: 500000, creditUsed: 180000, lastOrder: '2026-01-18', rating: 4.8 },
    { id: '2', name: 'Kano Wholesale Ltd', type: 'wholesaler', tier: 'platinum', totalOrders: 89, totalSpent: 45600000, creditLimit: 2000000, creditUsed: 1200000, lastOrder: '2026-01-17', rating: 4.5 },
    { id: '3', name: 'Ibadan Mini Mart', type: 'retailer', tier: 'silver', totalOrders: 45, totalSpent: 3200000, creditLimit: 200000, creditUsed: 150000, lastOrder: '2026-01-15', rating: 4.2 },
];

const routes: Route[] = [
    { id: '1', name: 'Lagos Island Route A', stops: 12, totalValue: 1850000, assignedAgent: 'John Okafor', status: 'in_progress', completionRate: 58 },
    { id: '2', name: 'Mainland Route B', stops: 18, totalValue: 2340000, assignedAgent: 'Mary Adebayo', status: 'pending', completionRate: 0 },
    { id: '3', name: 'Ikeja Express', stops: 8, totalValue: 980000, assignedAgent: 'Chidi Nnamdi', status: 'completed', completionRate: 100 },
];

// =============================================================================
// UTILITIES
// =============================================================================

const formatNaira = (amount: number): string => {
    if (amount >= 1000000) return `₦${(amount / 1000000).toFixed(1)}M`;
    if (amount >= 1000) return `₦${(amount / 1000).toFixed(0)}K`;
    return `₦${amount.toLocaleString()}`;
};

const getStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
        draft: 'default', pending: 'warning', confirmed: 'processing', processing: 'blue',
        shipped: 'cyan', delivered: 'success', cancelled: 'error',
        healthy: 'success', low: 'warning', critical: 'error', overstock: 'purple',
        in_progress: 'processing', completed: 'success',
    };
    return colors[status] || 'default';
};

// =============================================================================
// COMPONENTS
// =============================================================================

const QuickStatCard = ({ title, value, icon, color, trend, trendUp }: any) => (
    <Card hoverable bodyStyle={{ padding: 16 }}>
        <Space direction="vertical" size={0} style={{ width: '100%' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <Text type="secondary" style={{ fontSize: 13 }}>{title}</Text>
                <div style={{ background: `${color}15`, borderRadius: 8, padding: 6, display: 'flex' }}>
                    {React.cloneElement(icon, { style: { color, fontSize: 16 } })}
                </div>
            </div>
            <Title level={3} style={{ margin: '8px 0 4px', color }}>{value}</Title>
            {trend && (
                <Space size={4}>
                    <Tag color={trendUp ? 'success' : 'error'} style={{ fontSize: 11 }}>
                        {trendUp ? <RiseOutlined /> : <FallOutlined />} {trend}%
                    </Tag>
                    <Text type="secondary" style={{ fontSize: 11 }}>vs last week</Text>
                </Space>
            )}
        </Space>
    </Card>
);

const InventoryHealthView = ({ items }: { items: InventoryItem[] }) => {
    const healthy = items.filter(i => i.status === 'healthy').length;
    const low = items.filter(i => i.status === 'low').length;
    const critical = items.filter(i => i.status === 'critical').length;
    const overstock = items.filter(i => i.status === 'overstock').length;

    return (
        <Card title={<Space><InboxOutlined /> Inventory Health</Space>} extra={<Button type="link">Full Inventory →</Button>}>
            <Row gutter={16} style={{ marginBottom: 24 }}>
                <Col span={6}><Statistic title="Healthy" value={healthy} valueStyle={{ color: '#52c41a' }} prefix={<CheckCircleOutlined />} /></Col>
                <Col span={6}><Statistic title="Low Stock" value={low} valueStyle={{ color: '#fa8c16' }} prefix={<WarningOutlined />} /></Col>
                <Col span={6}><Statistic title="Critical" value={critical} valueStyle={{ color: '#ff4d4f' }} prefix={<FireOutlined />} /></Col>
                <Col span={6}><Statistic title="Overstock" value={overstock} valueStyle={{ color: '#722ed1' }} prefix={<StockOutlined />} /></Col>
            </Row>
            <List
                size="small"
                dataSource={items.filter(i => i.status !== 'healthy').slice(0, 5)}
                renderItem={(item) => (
                    <List.Item
                        actions={[
                            <Button size="small" type="primary" key="reorder">
                                {item.status === 'overstock' ? 'Discount' : 'Reorder'}
                            </Button>
                        ]}
                    >
                        <List.Item.Meta
                            avatar={<Avatar style={{ background: getStatusColor(item.status) === 'success' ? '#52c41a' : getStatusColor(item.status) === 'warning' ? '#fa8c16' : getStatusColor(item.status) === 'error' ? '#ff4d4f' : '#722ed1' }} icon={<InboxOutlined />} />}
                            title={item.name}
                            description={
                                <Space>
                                    <Tag color={getStatusColor(item.status)}>{item.status.toUpperCase()}</Tag>
                                    <Text type="secondary">{item.available} available • {item.location}</Text>
                                </Space>
                            }
                        />
                    </List.Item>
                )}
            />
        </Card>
    );
};

const OrdersTable = ({ orders, type }: { orders: Order[]; type: 'inbound' | 'outbound' | 'all' }) => {
    const filtered = type === 'all' ? orders : orders.filter(o => o.type === type);

    const columns: ColumnsType<Order> = [
        {
            title: 'Order',
            dataIndex: 'orderNumber',
            key: 'orderNumber',
            render: (text, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{text}</Text>
                    <Tag color={record.type === 'inbound' ? 'blue' : 'green'}>{record.type.toUpperCase()}</Tag>
                </Space>
            ),
        },
        { title: 'Counterparty', dataIndex: 'counterparty', key: 'counterparty' },
        { title: 'Items', dataIndex: 'items', key: 'items', align: 'center' },
        { title: 'Total', dataIndex: 'total', key: 'total', render: (v) => formatNaira(v), align: 'right' },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => <Tag color={getStatusColor(status)}>{status.toUpperCase()}</Tag>,
        },
        {
            title: 'Actions',
            key: 'actions',
            render: () => (
                <Space>
                    <Button size="small" icon={<EyeOutlined />} />
                    <Button size="small" icon={<EditOutlined />} />
                </Space>
            ),
        },
    ];

    return <Table columns={columns} dataSource={filtered} rowKey="id" size="small" />;
};

const CreditFacilityWidget = () => (
    <Card
        title={<Space><CreditCardOutlined style={{ color: '#1890ff' }} /> Credit Facility</Space>}
        extra={<Button type="link">Pay Now →</Button>}
    >
        <Row gutter={16}>
            <Col span={8}>
                <Statistic title="Credit Limit" value={formatNaira(5000000)} />
            </Col>
            <Col span={8}>
                <Statistic title="Available" value={formatNaira(3150000)} valueStyle={{ color: '#52c41a' }} />
            </Col>
            <Col span={8}>
                <Statistic title="Used" value={formatNaira(1850000)} valueStyle={{ color: '#fa8c16' }} />
            </Col>
        </Row>
        <Progress percent={37} strokeColor={{ '0%': '#52c41a', '100%': '#fa8c16' }} style={{ marginTop: 16 }} />
        <Alert
            message="Payment Due: Jan 25, 2026"
            description={`Amount: ${formatNaira(500000)}`}
            type="warning"
            showIcon
            style={{ marginTop: 16 }}
        />
    </Card>
);

const RouteManagement = ({ routes }: { routes: Route[] }) => (
    <Card title={<Space><EnvironmentOutlined style={{ color: '#1890ff' }} /> Route Management</Space>}>
        <List
            dataSource={routes}
            renderItem={(route) => (
                <List.Item
                    actions={[
                        <Button key="view" size="small" icon={<EyeOutlined />}>Track</Button>,
                    ]}
                >
                    <List.Item.Meta
                        avatar={
                            <Progress
                                type="circle"
                                percent={route.completionRate}
                                size={50}
                                format={(p) => `${p}%`}
                                status={route.status === 'completed' ? 'success' : 'active'}
                            />
                        }
                        title={
                            <Space>
                                <Text strong>{route.name}</Text>
                                <Tag color={getStatusColor(route.status)}>{route.status.replace('_', ' ').toUpperCase()}</Tag>
                            </Space>
                        }
                        description={
                            <Space split={<Divider type="vertical" />}>
                                <Text type="secondary">{route.stops} stops</Text>
                                <Text type="secondary">{formatNaira(route.totalValue)}</Text>
                                {route.assignedAgent && <Text>{route.assignedAgent}</Text>}
                            </Space>
                        }
                    />
                </List.Item>
            )}
        />
    </Card>
);

const CustomersList = ({ customers }: { customers: Customer[] }) => (
    <Card title={<Space><TeamOutlined /> Top Customers</Space>} extra={<Button type="link">All Customers →</Button>}>
        <List
            size="small"
            dataSource={customers}
            renderItem={(customer) => (
                <List.Item>
                    <List.Item.Meta
                        avatar={
                            <Avatar style={{ background: customer.tier === 'platinum' ? '#722ed1' : customer.tier === 'gold' ? '#d69e2e' : customer.tier === 'silver' ? '#8c8c8c' : '#1890ff' }}>
                                {customer.name[0]}
                            </Avatar>
                        }
                        title={
                            <Space>
                                <Text strong>{customer.name}</Text>
                                <Tag color={customer.tier === 'platinum' ? 'purple' : customer.tier === 'gold' ? 'gold' : customer.tier === 'silver' ? 'default' : 'blue'}>
                                    {customer.tier.toUpperCase()}
                                </Tag>
                            </Space>
                        }
                        description={
                            <Space split={<Divider type="vertical" />}>
                                <Text type="secondary">{customer.totalOrders} orders</Text>
                                <Text type="secondary">{formatNaira(customer.totalSpent)} spent</Text>
                                <Rate disabled defaultValue={customer.rating} style={{ fontSize: 12 }} />
                            </Space>
                        }
                    />
                </List.Item>
            )}
        />
    </Card>
);

// =============================================================================
// MAIN COMPONENT
// =============================================================================

export default function PartnerDashboard() {
    const [currentRole, setCurrentRole] = useState<PartnerRole>('distributor');
    const [orderTab, setOrderTab] = useState<'all' | 'inbound' | 'outbound'>('all');
    const roleConfig = roleConfigs.find(r => r.key === currentRole)!;

    // Dynamic metrics based on role
    const getMetrics = () => {
        switch (currentRole) {
            case 'manufacturer':
                return { m1: '12,450 units', m1Label: 'Production Today', m2: '₦45.6M', m2Label: 'Distributor Orders', m3: '23', m3Label: 'Active SKUs', m4: '98.5%', m4Label: 'Fulfillment Rate' };
            case 'distributor':
                return { m1: '₦8.9M', m1Label: 'Today\'s Sales', m2: '156', m2Label: 'Active Retailers', m3: '₦5M', m3Label: 'Credit Extended', m4: '85%', m4Label: 'Route Coverage' };
            case 'wholesaler':
                return { m1: '₦12.3M', m1Label: 'Today\'s Orders', m2: '45', m2Label: 'Pending Orders', m3: '₦2.1M', m3Label: 'Credit Available', m4: '92%', m4Label: 'Stock Availability' };
            case 'retailer':
                return { m1: '₦890K', m1Label: 'Today\'s Sales', m2: '234', m2Label: 'Transactions', m3: '12', m3Label: 'Low Stock Items', m4: '₦45K', m4Label: 'Avg Basket' };
            case 'logistics':
                return { m1: '45', m1Label: 'Active Trips', m2: '23', m2Label: 'Drivers Online', m3: '94%', m3Label: 'On-Time Rate', m4: '156', m4Label: 'PODs Today' };
            case 'warehouse':
                return { m1: '12,450', m1Label: 'Items in Stock', m2: '89', m2Label: 'Orders to Pick', m3: '45', m3Label: 'Inbound Today', m4: '98%', m4Label: 'Space Utilization' };
            default:
                return { m1: '0', m1Label: '', m2: '0', m2Label: '', m3: '0', m3Label: '', m4: '0', m4Label: '' };
        }
    };

    const metrics = getMetrics();

    return (
        <div style={{ background: '#f0f2f5', minHeight: '100vh' }}>
            {/* Role Switcher Header */}
            <div style={{
                background: `linear-gradient(135deg, ${roleConfig.color} 0%, ${roleConfig.color}dd 100%)`,
                padding: '16px 24px',
                color: 'white'
            }}>
                <Row align="middle" justify="space-between">
                    <Col>
                        <Space>
                            {roleConfig.icon}
                            <Title level={3} style={{ color: 'white', margin: 0 }}>{roleConfig.name} Dashboard</Title>
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <Segmented
                                options={roleConfigs.map(r => ({
                                    label: (
                                        <Space>
                                            {r.icon}
                                            <span style={{ display: 'inline' }}>{r.name}</span>
                                        </Space>
                                    ),
                                    value: r.key,
                                }))}
                                value={currentRole}
                                onChange={(v) => setCurrentRole(v as PartnerRole)}
                                style={{ background: 'rgba(255,255,255,0.2)' }}
                            />
                            <Badge count={5}>
                                <Button icon={<BellOutlined />} type="default" style={{ background: 'rgba(255,255,255,0.2)', border: 'none', color: 'white' }} />
                            </Badge>
                        </Space>
                    </Col>
                </Row>
            </div>

            <div style={{ padding: 24 }}>
                {/* Metrics Row */}
                <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                    <Col xs={12} sm={6}>
                        <QuickStatCard title={metrics.m1Label} value={metrics.m1} icon={<DollarOutlined />} color={roleConfig.color} trend={12.5} trendUp />
                    </Col>
                    <Col xs={12} sm={6}>
                        <QuickStatCard title={metrics.m2Label} value={metrics.m2} icon={<ShoppingCartOutlined />} color="#1890ff" trend={8.2} trendUp />
                    </Col>
                    <Col xs={12} sm={6}>
                        <QuickStatCard title={metrics.m3Label} value title={metrics.m3Label} value={metrics.m3} icon={<InboxOutlined />} color="#fa8c16" trend={-2.1} trendUp={false} />
                    </Col>
                    <Col xs={12} sm={6}>
                        <QuickStatCard title={metrics.m4Label} value={metrics.m4} icon={<LineChartOutlined />} color="#52c41a" />
                    </Col>
                </Row>

                {/* Quick Actions */}
                <Card bodyStyle={{ padding: 12 }} style={{ marginBottom: 24 }}>
                    <Space wrap>
                        <Button type="primary" icon={<PlusOutlined />}>New Order</Button>
                        <Button icon={<SyncOutlined />}>Sync Inventory</Button>
                        <Button icon={<QrcodeOutlined />}>Scan Barcode</Button>
                        <Button icon={<ExportOutlined />}>Export Report</Button>
                        <Button icon={<RobotOutlined />}>AI Suggestions</Button>
                    </Space>
                </Card>

                {/* Main Content */}
                <Row gutter={[16, 16]}>
                    <Col xs={24} lg={16}>
                        <Card
                            title={<Space><ShoppingCartOutlined /> Orders</Space>}
                            extra={
                                <Segmented
                                    options={[
                                        { label: 'All', value: 'all' },
                                        { label: 'Inbound', value: 'inbound' },
                                        { label: 'Outbound', value: 'outbound' },
                                    ]}
                                    value={orderTab}
                                    onChange={(v) => setOrderTab(v as any)}
                                />
                            }
                        >
                            <OrdersTable orders={orders} type={orderTab} />
                        </Card>

                        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
                            <Col span={12}>
                                <InventoryHealthView items={inventory} />
                            </Col>
                            <Col span={12}>
                                {currentRole === 'distributor' && <RouteManagement routes={routes} />}
                                {currentRole !== 'distributor' && <CustomersList customers={customers} />}
                            </Col>
                        </Row>
                    </Col>

                    <Col xs={24} lg={8}>
                        <Space direction="vertical" style={{ width: '100%' }} size={16}>
                            <CreditFacilityWidget />

                            <Card title={<Space><RobotOutlined /> AI Insights</Space>}>
                                <List
                                    size="small"
                                    dataSource={[
                                        { icon: <ThunderboltOutlined style={{ color: '#52c41a' }} />, text: 'Golden Penny Rice demand up 34% - stock up now' },
                                        { icon: <WarningOutlined style={{ color: '#ff4d4f' }} />, text: '2 retailers with overdue payments' },
                                        { icon: <StarOutlined style={{ color: '#d69e2e' }} />, text: 'Best selling: Peak Milk 400g (450 units)' },
                                    ]}
                                    renderItem={(item) => (
                                        <List.Item>
                                            <Space>{item.icon} <Text>{item.text}</Text></Space>
                                        </List.Item>
                                    )}
                                />
                            </Card>

                            <Card title={<Space><ClockCircleOutlined /> Recent Activity</Space>}>
                                <Timeline
                                    items={[
                                        { color: 'green', children: 'Order ORD-0145 confirmed • 10:30 AM' },
                                        { color: 'blue', children: 'Inventory sync completed • 10:00 AM' },
                                        { color: 'orange', children: 'Low stock alert: Indomie • 9:45 AM' },
                                        { color: 'green', children: 'Payment received ₦1.2M • 9:15 AM' },
                                    ]}
                                />
                            </Card>
                        </Space>
                    </Col>
                </Row>
            </div>
        </div>
    );
}
