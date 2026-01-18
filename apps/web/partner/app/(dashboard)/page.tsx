// OmniRoute Partner Portal - Dashboard Page
// Role-based B2B Partner Dashboard

'use client';

import { useState } from 'react';
import { Card, Row, Col, Statistic, Table, Typography, Space, Tag, Progress, Avatar, List, Button, Divider } from 'antd';
import {
    ShoppingCartOutlined,
    DollarOutlined,
    InboxOutlined,
    TruckOutlined,
    UserOutlined,
    RiseOutlined,
    BellOutlined,
    PlusOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;

// Types
interface Order {
    id: string;
    orderNumber: string;
    customer: string;
    items: number;
    total: number;
    status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered';
    createdAt: string;
}

interface LowStockItem {
    id: string;
    name: string;
    sku: string;
    currentStock: number;
    reorderPoint: number;
    daysUntilStockout: number;
}

interface CreditInfo {
    limit: number;
    used: number;
    available: number;
    nextPaymentDate: string;
    nextPaymentAmount: number;
}

// Mock Data
const recentOrders: Order[] = [
    { id: '1', orderNumber: 'PO-2026-0089', customer: 'Kano Retailers Coop', items: 15, total: 450000, status: 'pending', createdAt: '2026-01-18T10:30:00Z' },
    { id: '2', orderNumber: 'PO-2026-0088', customer: 'Ibadan Marts Network', items: 8, total: 280000, status: 'confirmed', createdAt: '2026-01-18T09:15:00Z' },
    { id: '3', orderNumber: 'PO-2026-0087', customer: 'Enugu Fresh Stores', items: 22, total: 890000, status: 'processing', createdAt: '2026-01-17T16:00:00Z' },
    { id: '4', orderNumber: 'PO-2026-0086', customer: 'Abuja Central Market', items: 5, total: 125000, status: 'shipped', createdAt: '2026-01-17T14:20:00Z' },
];

const lowStockItems: LowStockItem[] = [
    { id: '1', name: 'Peak Milk 400g', sku: 'PM-400-CTN', currentStock: 45, reorderPoint: 100, daysUntilStockout: 3 },
    { id: '2', name: 'Indomie Chicken (Carton)', sku: 'IND-CHK-CTN', currentStock: 120, reorderPoint: 200, daysUntilStockout: 5 },
    { id: '3', name: 'Golden Penny Spaghetti', sku: 'GPS-500-CTN', currentStock: 80, reorderPoint: 150, daysUntilStockout: 4 },
];

const creditInfo: CreditInfo = {
    limit: 5000000,
    used: 1850000,
    available: 3150000,
    nextPaymentDate: '2026-01-25',
    nextPaymentAmount: 500000,
};

// Formatters
const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
    }).format(amount);
};

const getStatusColor = (status: Order['status']): string => {
    const colors: Record<Order['status'], string> = {
        pending: 'warning',
        confirmed: 'processing',
        processing: 'blue',
        shipped: 'cyan',
        delivered: 'success',
    };
    return colors[status];
};

// Main Component
export default function PartnerDashboard() {
    const partnerType = 'distributor'; // Would come from auth context

    const orderColumns: ColumnsType<Order> = [
        {
            title: 'Order',
            dataIndex: 'orderNumber',
            key: 'orderNumber',
            render: (text) => <Text strong>{text}</Text>,
        },
        {
            title: 'Customer',
            dataIndex: 'customer',
            key: 'customer',
        },
        {
            title: 'Items',
            dataIndex: 'items',
            key: 'items',
            align: 'center',
        },
        {
            title: 'Total',
            dataIndex: 'total',
            key: 'total',
            render: (value) => formatNaira(value),
            align: 'right',
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => (
                <Tag color={getStatusColor(status)}>
                    {status.toUpperCase()}
                </Tag>
            ),
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            {/* Header */}
            <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <Title level={2} style={{ margin: 0 }}>Partner Dashboard</Title>
                    <Text type="secondary">Welcome back! Here's your business overview.</Text>
                </div>
                <Space>
                    <Button icon={<PlusOutlined />} type="primary">New Order</Button>
                    <Button icon={<BellOutlined />}>Notifications</Button>
                </Space>
            </div>

            {/* Key Metrics */}
            <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card hoverable>
                        <Statistic
                            title="Today's Sales"
                            value={formatNaira(1850000)}
                            prefix={<DollarOutlined />}
                            valueStyle={{ color: '#52c41a' }}
                        />
                        <Space style={{ marginTop: 8 }}>
                            <Tag color="success" icon={<RiseOutlined />}>12%</Tag>
                            <Text type="secondary">vs yesterday</Text>
                        </Space>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card hoverable>
                        <Statistic
                            title="Pending Orders"
                            value={23}
                            prefix={<ShoppingCartOutlined />}
                            valueStyle={{ color: '#fa8c16' }}
                        />
                        <Text type="secondary">₦4.2M value</Text>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card hoverable>
                        <Statistic
                            title="Low Stock Items"
                            value={lowStockItems.length}
                            prefix={<InboxOutlined />}
                            valueStyle={{ color: '#ff4d4f' }}
                        />
                        <Text type="secondary">Action required</Text>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card hoverable>
                        <Statistic
                            title="In Transit"
                            value={8}
                            prefix={<TruckOutlined />}
                            valueStyle={{ color: '#1890ff' }}
                        />
                        <Text type="secondary">5 arriving today</Text>
                    </Card>
                </Col>
            </Row>

            <Row gutter={[16, 16]}>
                {/* Recent Orders */}
                <Col xs={24} lg={14}>
                    <Card title="Recent Orders" extra={<a href="/orders">View All</a>}>
                        <Table
                            columns={orderColumns}
                            dataSource={recentOrders}
                            rowKey="id"
                            pagination={false}
                            size="small"
                        />
                    </Card>
                </Col>

                {/* Side Panel */}
                <Col xs={24} lg={10}>
                    {/* Credit Facility */}
                    <Card title="Credit Facility" style={{ marginBottom: 16 }}>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Text type="secondary">Credit Limit</Text>
                                <Title level={4} style={{ margin: '4px 0' }}>{formatNaira(creditInfo.limit)}</Title>
                            </Col>
                            <Col span={12}>
                                <Text type="secondary">Available</Text>
                                <Title level={4} style={{ margin: '4px 0', color: '#52c41a' }}>{formatNaira(creditInfo.available)}</Title>
                            </Col>
                        </Row>
                        <Progress
                            percent={(creditInfo.used / creditInfo.limit) * 100}
                            strokeColor="#1890ff"
                            format={() => `${formatNaira(creditInfo.used)} used`}
                        />
                        <Divider />
                        <Space direction="vertical" style={{ width: '100%' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                <Text type="secondary">Next Payment</Text>
                                <Text strong>{creditInfo.nextPaymentDate}</Text>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                <Text type="secondary">Amount Due</Text>
                                <Text strong style={{ color: '#fa8c16' }}>{formatNaira(creditInfo.nextPaymentAmount)}</Text>
                            </div>
                        </Space>
                    </Card>

                    {/* Low Stock Alerts */}
                    <Card
                        title={
                            <Space>
                                <InboxOutlined style={{ color: '#ff4d4f' }} />
                                <span>Low Stock Alerts</span>
                            </Space>
                        }
                        extra={<a href="/inventory/alerts">View All</a>}
                    >
                        <List
                            itemLayout="horizontal"
                            dataSource={lowStockItems}
                            renderItem={(item) => (
                                <List.Item
                                    actions={[<Button size="small" type="link">Reorder</Button>]}
                                >
                                    <List.Item.Meta
                                        title={item.name}
                                        description={
                                            <Space>
                                                <Tag color="error">{item.currentStock} left</Tag>
                                                <Text type="secondary">{item.daysUntilStockout} days until stockout</Text>
                                            </Space>
                                        }
                                    />
                                </List.Item>
                            )}
                        />
                    </Card>
                </Col>
            </Row>
        </div>
    );
}
