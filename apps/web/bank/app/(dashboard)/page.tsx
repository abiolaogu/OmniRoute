// OmniRoute Bank Portal - Dashboard Page
// Financial Institution Loan & ATC Management

'use client';

import { useState } from 'react';
import { Card, Row, Col, Statistic, Table, Typography, Space, Tag, Progress, Tabs, Timeline, Select, Button } from 'antd';
import {
    BankOutlined,
    DollarOutlined,
    ClockCircleOutlined,
    WarningOutlined,
    CheckCircleOutlined,
    SyncOutlined,
    FileSearchOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;

// Types
interface LoanApplication {
    id: string;
    applicationNumber: string;
    businessName: string;
    amount: number;
    creditScore: number;
    riskBand: 'prime' | 'near_prime' | 'subprime' | 'deep_subprime';
    status: 'pending' | 'under_review' | 'approved' | 'disbursed' | 'rejected';
    appliedAt: string;
}

interface ATCMandate {
    id: string;
    customerId: string;
    customerName: string;
    mandateAmount: number;
    status: 'active' | 'pending' | 'expired' | 'cancelled';
    nextCollection: string;
    bank: string;
    accountLast4: string;
}

interface CollectionSummary {
    totalExpected: number;
    totalCollected: number;
    failedCount: number;
    successRate: number;
}

// Mock Data
const loanApplications: LoanApplication[] = [
    { id: '1', applicationNumber: 'LN-2026-0042', businessName: 'Mama Ngozi Stores', amount: 500000, creditScore: 720, riskBand: 'prime', status: 'pending', appliedAt: '2026-01-18T10:30:00Z' },
    { id: '2', applicationNumber: 'LN-2026-0041', businessName: 'Chidi Wholesale Ltd', amount: 2500000, creditScore: 680, riskBand: 'near_prime', status: 'under_review', appliedAt: '2026-01-18T09:15:00Z' },
    { id: '3', applicationNumber: 'LN-2026-0040', businessName: 'Lagos Retail Hub', amount: 1000000, creditScore: 750, riskBand: 'prime', status: 'approved', appliedAt: '2026-01-17T16:00:00Z' },
    { id: '4', applicationNumber: 'LN-2026-0039', businessName: 'Abuja SuperMart', amount: 750000, creditScore: 590, riskBand: 'subprime', status: 'rejected', appliedAt: '2026-01-17T14:20:00Z' },
    { id: '5', applicationNumber: 'LN-2026-0038', businessName: 'Port Harcourt Foods', amount: 1500000, creditScore: 710, riskBand: 'prime', status: 'disbursed', appliedAt: '2026-01-16T11:00:00Z' },
];

const atcMandates: ATCMandate[] = [
    { id: '1', customerId: 'C001', customerName: 'Lagos Retail Hub', mandateAmount: 125000, status: 'active', nextCollection: '2026-01-25', bank: 'GTBank', accountLast4: '4523' },
    { id: '2', customerId: 'C002', customerName: 'Chidi Wholesale Ltd', mandateAmount: 280000, status: 'active', nextCollection: '2026-01-28', bank: 'First Bank', accountLast4: '7891' },
    { id: '3', customerId: 'C003', customerName: 'Port Harcourt Foods', mandateAmount: 95000, status: 'pending', nextCollection: '2026-01-30', bank: 'UBA', accountLast4: '2345' },
    { id: '4', customerId: 'C004', customerName: 'Ibadan Mart', mandateAmount: 150000, status: 'expired', nextCollection: '-', bank: 'Access Bank', accountLast4: '6789' },
];

const collectionSummary: CollectionSummary = {
    totalExpected: 12500000,
    totalCollected: 11750000,
    failedCount: 23,
    successRate: 94,
};

// Formatters
const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(amount);
};

const getStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
        pending: 'warning',
        under_review: 'processing',
        approved: 'success',
        disbursed: 'blue',
        rejected: 'error',
        active: 'success',
        expired: 'default',
        cancelled: 'error',
    };
    return colors[status] || 'default';
};

const getRiskBandColor = (band: LoanApplication['riskBand']): string => {
    const colors: Record<LoanApplication['riskBand'], string> = {
        prime: 'green',
        near_prime: 'blue',
        subprime: 'orange',
        deep_subprime: 'red',
    };
    return colors[band];
};

// Main Component
export default function BankDashboard() {
    const [activeTab, setActiveTab] = useState('overview');

    const loanColumns: ColumnsType<LoanApplication> = [
        {
            title: 'Application',
            dataIndex: 'applicationNumber',
            key: 'applicationNumber',
            render: (text) => <Text strong>{text}</Text>,
        },
        {
            title: 'Business',
            dataIndex: 'businessName',
            key: 'businessName',
        },
        {
            title: 'Amount',
            dataIndex: 'amount',
            key: 'amount',
            render: (value) => formatNaira(value),
            align: 'right',
        },
        {
            title: 'Credit Score',
            dataIndex: 'creditScore',
            key: 'creditScore',
            render: (score, record) => (
                <Space>
                    <Text strong>{score}</Text>
                    <Tag color={getRiskBandColor(record.riskBand)}>{record.riskBand.replace('_', ' ').toUpperCase()}</Tag>
                </Space>
            ),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => (
                <Tag color={getStatusColor(status)}>
                    {status.replace('_', ' ').toUpperCase()}
                </Tag>
            ),
        },
        {
            title: 'Action',
            key: 'action',
            render: (_, record) => (
                <Button type="link" icon={<FileSearchOutlined />}>
                    Review
                </Button>
            ),
        },
    ];

    const atcColumns: ColumnsType<ATCMandate> = [
        {
            title: 'Customer',
            dataIndex: 'customerName',
            key: 'customerName',
        },
        {
            title: 'Mandate Amount',
            dataIndex: 'mandateAmount',
            key: 'mandateAmount',
            render: (value) => formatNaira(value),
            align: 'right',
        },
        {
            title: 'Bank Account',
            key: 'account',
            render: (_, record) => (
                <Text>{record.bank} ****{record.accountLast4}</Text>
            ),
        },
        {
            title: 'Next Collection',
            dataIndex: 'nextCollection',
            key: 'nextCollection',
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
            <div style={{ marginBottom: 24 }}>
                <Title level={2} style={{ margin: 0 }}>
                    <BankOutlined style={{ marginRight: 12 }} />
                    Bank Dashboard
                </Title>
                <Text type="secondary">Loan portfolio, ATC mandates, and collections overview</Text>
            </div>

            {/* Key Metrics */}
            <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card>
                        <Statistic
                            title="Active Loans"
                            value={156}
                            prefix={<DollarOutlined />}
                            valueStyle={{ color: '#1890ff' }}
                        />
                        <Text type="secondary">₦245.6M disbursed</Text>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card>
                        <Statistic
                            title="Pending Applications"
                            value={23}
                            prefix={<ClockCircleOutlined />}
                            valueStyle={{ color: '#fa8c16' }}
                        />
                        <Text type="secondary">₦34.2M requested</Text>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card>
                        <Statistic
                            title="Collection Rate"
                            value={collectionSummary.successRate}
                            suffix="%"
                            prefix={<CheckCircleOutlined />}
                            valueStyle={{ color: '#52c41a' }}
                        />
                        <Progress percent={collectionSummary.successRate} showInfo={false} strokeColor="#52c41a" />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card>
                        <Statistic
                            title="Failed Collections"
                            value={collectionSummary.failedCount}
                            prefix={<WarningOutlined />}
                            valueStyle={{ color: '#ff4d4f' }}
                        />
                        <Text type="secondary">₦2.3M pending retry</Text>
                    </Card>
                </Col>
            </Row>

            {/* Tabs */}
            <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={[
                    {
                        key: 'overview',
                        label: 'Loan Applications',
                        children: (
                            <Card>
                                <Table
                                    columns={loanColumns}
                                    dataSource={loanApplications}
                                    rowKey="id"
                                    pagination={{ pageSize: 5 }}
                                />
                            </Card>
                        ),
                    },
                    {
                        key: 'atc',
                        label: 'ATC Mandates',
                        children: (
                            <Card>
                                <Table
                                    columns={atcColumns}
                                    dataSource={atcMandates}
                                    rowKey="id"
                                    pagination={{ pageSize: 5 }}
                                />
                            </Card>
                        ),
                    },
                    {
                        key: 'collections',
                        label: 'Today\'s Collections',
                        children: (
                            <Card>
                                <Row gutter={[16, 16]}>
                                    <Col span={12}>
                                        <Card type="inner" title="Expected">
                                            <Statistic value={formatNaira(collectionSummary.totalExpected)} />
                                        </Card>
                                    </Col>
                                    <Col span={12}>
                                        <Card type="inner" title="Collected">
                                            <Statistic
                                                value={formatNaira(collectionSummary.totalCollected)}
                                                valueStyle={{ color: '#52c41a' }}
                                            />
                                        </Card>
                                    </Col>
                                </Row>
                                <Timeline style={{ marginTop: 24 }}>
                                    <Timeline.Item color="green" dot={<CheckCircleOutlined />}>
                                        Lagos Retail Hub - ₦125,000 collected at 09:15 AM
                                    </Timeline.Item>
                                    <Timeline.Item color="green" dot={<CheckCircleOutlined />}>
                                        Chidi Wholesale Ltd - ₦280,000 collected at 09:30 AM
                                    </Timeline.Item>
                                    <Timeline.Item color="red" dot={<WarningOutlined />}>
                                        Abuja SuperMart - ₦150,000 FAILED (Insufficient funds)
                                    </Timeline.Item>
                                    <Timeline.Item color="blue" dot={<SyncOutlined spin />}>
                                        Port Harcourt Foods - ₦95,000 processing...
                                    </Timeline.Item>
                                </Timeline>
                            </Card>
                        ),
                    },
                ]}
            />
        </div>
    );
}
