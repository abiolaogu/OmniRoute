// OmniRoute Admin Portal - Dashboard Page
// Platform Administration Dashboard

'use client';

import { useState } from 'react';
import { Card, Row, Col, Statistic, Table, Typography, Space, Tag, Select, DatePicker } from 'antd';
import {
    ShoppingCartOutlined,
    DollarOutlined,
    UserOutlined,
    RiseOutlined,
    FallOutlined,
    ShoppingOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

// Types
interface DashboardStats {
    revenue: number;
    revenueChange: number;
    orders: number;
    ordersChange: number;
    customers: number;
    customersChange: number;
    avgOrderValue: number;
    avgOrderValueChange: number;
}

interface RecentOrder {
    id: string;
    orderNumber: string;
    customer: string;
    total: number;
    status: 'pending' | 'processing' | 'shipped' | 'delivered';
    createdAt: string;
}

interface TopProduct {
    id: string;
    name: string;
    sold: number;
    revenue: number;
    image?: string;
}

// Mock Data
const stats: DashboardStats = {
    revenue: 45678900,
    revenueChange: 12.5,
    orders: 1234,
    ordersChange: 8.2,
    customers: 567,
    customersChange: 15.3,
    avgOrderValue: 37000,
    avgOrderValueChange: -2.1,
};

const recentOrders: RecentOrder[] = [
    { id: '1', orderNumber: 'ORD-2026-0001', customer: 'Mama Ngozi Stores', total: 125000, status: 'pending', createdAt: '2026-01-18T10:30:00Z' },
    { id: '2', orderNumber: 'ORD-2026-0002', customer: 'Chidi Wholesale', total: 890000, status: 'processing', createdAt: '2026-01-18T09:15:00Z' },
    { id: '3', orderNumber: 'ORD-2026-0003', customer: 'Lagos Retail Hub', total: 456000, status: 'shipped', createdAt: '2026-01-18T08:00:00Z' },
    { id: '4', orderNumber: 'ORD-2026-0004', customer: 'Abuja SuperMart', total: 234000, status: 'delivered', createdAt: '2026-01-17T16:45:00Z' },
    { id: '5', orderNumber: 'ORD-2026-0005', customer: 'Port Harcourt Mini', total: 78000, status: 'pending', createdAt: '2026-01-17T14:20:00Z' },
];

const topProducts: TopProduct[] = [
    { id: '1', name: 'Peak Milk 400g (Carton)', sold: 450, revenue: 1125000 },
    { id: '2', name: 'Indomie Chicken (Carton)', sold: 380, revenue: 760000 },
    { id: '3', name: 'Golden Penny Rice 50kg', sold: 120, revenue: 7800000 },
    { id: '4', name: 'Dangote Sugar 50kg', sold: 200, revenue: 5200000 },
    { id: '5', name: 'Kings Oil 5L', sold: 280, revenue: 980000 },
];

// Formatters
const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(amount);
};

const getStatusColor = (status: RecentOrder['status']): string => {
    const colors: Record<RecentOrder['status'], string> = {
        pending: 'warning',
        processing: 'processing',
        shipped: 'blue',
        delivered: 'success',
    };
    return colors[status];
};

// Components
const StatCard = ({
    title,
    value,
    change,
    prefix,
    icon,
    color
}: {
    title: string;
    value: string | number;
    change: number;
    prefix?: string;
    icon: React.ReactNode;
    color: string;
}) => (
    <Card hoverable>
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
            <Space>
                <div style={{
                    background: `${color}15`,
                    borderRadius: 8,
                    padding: 8,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                }}>
                    {icon}
                </div>
                <Text type="secondary">{title}</Text>
            </Space>
            <Statistic
                value={value}
                prefix={prefix}
                valueStyle={{ fontSize: 28, fontWeight: 600 }}
            />
            <Space>
                {change >= 0 ? (
                    <Tag color="success" icon={<RiseOutlined />}>{change}%</Tag>
                ) : (
                    <Tag color="error" icon={<FallOutlined />}>{Math.abs(change)}%</Tag>
                )}
                <Text type="secondary" style={{ fontSize: 12 }}>vs last month</Text>
            </Space>
        </Space>
    </Card>
);

// Main Component
export default function AdminDashboard() {
    const [period, setPeriod] = useState('month');

    const orderColumns: ColumnsType<RecentOrder> = [
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
        {
            title: 'Date',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: (date) => new Date(date).toLocaleDateString('en-NG'),
        },
    ];

    const productColumns: ColumnsType<TopProduct> = [
        {
            title: 'Product',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: 'Sold',
            dataIndex: 'sold',
            key: 'sold',
            align: 'right',
        },
        {
            title: 'Revenue',
            dataIndex: 'revenue',
            key: 'revenue',
            render: (value) => formatNaira(value),
            align: 'right',
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            {/* Header */}
            <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <Title level={2} style={{ margin: 0 }}>Dashboard</Title>
                    <Text type="secondary">Welcome back! Here's what's happening with your platform.</Text>
                </div>
                <Space>
                    <Select
                        value={period}
                        onChange={setPeriod}
                        options={[
                            { value: 'today', label: 'Today' },
                            { value: 'week', label: 'This Week' },
                            { value: 'month', label: 'This Month' },
                            { value: 'quarter', label: 'This Quarter' },
                            { value: 'year', label: 'This Year' },
                        ]}
                        style={{ width: 140 }}
                    />
                    <RangePicker />
                </Space>
            </div>

            {/* Stats Cards */}
            <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="Total Revenue"
                        value={formatNaira(stats.revenue)}
                        change={stats.revenueChange}
                        icon={<DollarOutlined style={{ fontSize: 20, color: '#52c41a' }} />}
                        color="#52c41a"
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="Total Orders"
                        value={stats.orders.toLocaleString()}
                        change={stats.ordersChange}
                        icon={<ShoppingCartOutlined style={{ fontSize: 20, color: '#1890ff' }} />}
                        color="#1890ff"
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="New Customers"
                        value={stats.customers.toLocaleString()}
                        change={stats.customersChange}
                        icon={<UserOutlined style={{ fontSize: 20, color: '#722ed1' }} />}
                        color="#722ed1"
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="Avg Order Value"
                        value={formatNaira(stats.avgOrderValue)}
                        change={stats.avgOrderValueChange}
                        icon={<ShoppingOutlined style={{ fontSize: 20, color: '#fa8c16' }} />}
                        color="#fa8c16"
                    />
                </Col>
            </Row>

            {/* Tables */}
            <Row gutter={[16, 16]}>
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
                <Col xs={24} lg={10}>
                    <Card title="Top Selling Products" extra={<a href="/products">View All</a>}>
                        <Table
                            columns={productColumns}
                            dataSource={topProducts}
                            rowKey="id"
                            pagination={false}
                            size="small"
                        />
                    </Card>
                </Col>
            </Row>
        </div>
    );
}
