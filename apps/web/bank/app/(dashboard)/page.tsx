// OmniRoute Bank Portal - Comprehensive Financial Institution Dashboard
// Credit Management, ATC Collections, Real-time Reconciliation

'use client';

import React, { useState } from 'react';
import {
    Card, Row, Col, Statistic, Table, Typography, Space, Tag, Progress,
    Tabs, Timeline, Select, Button, Avatar, Badge, Tooltip, Calendar,
    Modal, Input, Drawer, List, Segmented, Steps, Descriptions, Rate,
    message, Popconfirm, Empty, Divider
} from 'antd';
import {
    BankOutlined, DollarOutlined, ClockCircleOutlined, WarningOutlined,
    CheckCircleOutlined, SyncOutlined, FileSearchOutlined, UserOutlined,
    CalendarOutlined, LineChartOutlined, SafetyOutlined, AuditOutlined,
    FundOutlined, PieChartOutlined, BarChartOutlined, ThunderboltOutlined,
    AlertOutlined, CreditCardOutlined, TeamOutlined, GlobalOutlined,
    FileProtectOutlined, HistoryOutlined, PhoneOutlined, MailOutlined,
    EnvironmentOutlined, IdcardOutlined, RobotOutlined, RiseOutlined,
    FallOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { Dayjs } from 'dayjs';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { TabPane } = Tabs;

// =============================================================================
// TYPES
// =============================================================================

interface LoanApplication {
    id: string;
    applicationNumber: string;
    businessName: string;
    ownerName: string;
    phone: string;
    email: string;
    amount: number;
    purpose: string;
    tenure: number;
    creditScore: number;
    riskBand: 'prime' | 'near_prime' | 'subprime' | 'deep_subprime';
    status: 'new' | 'kyc_pending' | 'under_review' | 'credit_check' | 'approved' | 'disbursed' | 'rejected';
    kycStatus: 'pending' | 'in_progress' | 'verified' | 'failed';
    transactionHistory: {
        ordersLast90Days: number;
        totalGMV: number;
        avgOrderValue: number;
        paymentReliability: number;
    };
    appliedAt: string;
    assignedTo?: string;
}

interface ATCMandate {
    id: string;
    loanId: string;
    customerId: string;
    customerName: string;
    mandateAmount: number;
    frequency: 'daily' | 'weekly' | 'monthly';
    status: 'active' | 'pending_setup' | 'failed' | 'cancelled' | 'expired';
    bankName: string;
    accountNumber: string;
    nextDebitDate: string;
    remainingBalance: number;
    totalRepaid: number;
    consecutiveFailures: number;
}

interface CollectionEvent {
    date: string;
    customerId: string;
    customerName: string;
    amount: number;
    status: 'scheduled' | 'processing' | 'success' | 'failed' | 'retry';
    failReason?: string;
}

interface ReconciliationEntry {
    id: string;
    date: string;
    expectedAmount: number;
    actualAmount: number;
    variance: number;
    status: 'matched' | 'unmatched' | 'pending';
    source: string;
}

// =============================================================================
// MOCK DATA
// =============================================================================

const loanApplications: LoanApplication[] = [
    {
        id: '1',
        applicationNumber: 'LN-2026-0089',
        businessName: 'Mama Ngozi Ventures',
        ownerName: 'Ngozi Okafor',
        phone: '+234 803 456 7890',
        email: 'ngozi@mamangozi.com',
        amount: 500000,
        purpose: 'Inventory Expansion',
        tenure: 6,
        creditScore: 720,
        riskBand: 'prime',
        status: 'new',
        kycStatus: 'pending',
        transactionHistory: { ordersLast90Days: 145, totalGMV: 4500000, avgOrderValue: 31034, paymentReliability: 96 },
        appliedAt: '2026-01-18T10:30:00Z',
    },
    {
        id: '2',
        applicationNumber: 'LN-2026-0088',
        businessName: 'ChidiMart Distribution',
        ownerName: 'Chidi Nnamdi',
        phone: '+234 805 123 4567',
        email: 'chidi@chidimart.ng',
        amount: 2500000,
        purpose: 'Vehicle Purchase',
        tenure: 12,
        creditScore: 685,
        riskBand: 'near_prime',
        status: 'credit_check',
        kycStatus: 'verified',
        transactionHistory: { ordersLast90Days: 280, totalGMV: 12500000, avgOrderValue: 44642, paymentReliability: 89 },
        appliedAt: '2026-01-17T14:15:00Z',
        assignedTo: 'Analyst A',
    },
    {
        id: '3',
        applicationNumber: 'LN-2026-0087',
        businessName: 'Lagos Fresh Foods',
        ownerName: 'Tunde Bakare',
        phone: '+234 701 987 6543',
        email: 'tunde@lagosfresh.com',
        amount: 1000000,
        purpose: 'Cold Storage Equipment',
        tenure: 9,
        creditScore: 752,
        riskBand: 'prime',
        status: 'approved',
        kycStatus: 'verified',
        transactionHistory: { ordersLast90Days: 340, totalGMV: 18000000, avgOrderValue: 52941, paymentReliability: 98 },
        appliedAt: '2026-01-16T09:00:00Z',
        assignedTo: 'Analyst B',
    },
];

const atcMandates: ATCMandate[] = [
    { id: '1', loanId: 'LN-2025-0045', customerId: 'C001', customerName: 'Lagos Retail Hub', mandateAmount: 125000, frequency: 'monthly', status: 'active', bankName: 'GTBank', accountNumber: '****4523', nextDebitDate: '2026-01-25', remainingBalance: 875000, totalRepaid: 375000, consecutiveFailures: 0 },
    { id: '2', loanId: 'LN-2025-0032', customerId: 'C002', customerName: 'Kano Wholesale Ltd', mandateAmount: 85000, frequency: 'weekly', status: 'failed', bankName: 'First Bank', accountNumber: '****7891', nextDebitDate: '2026-01-22', remainingBalance: 510000, totalRepaid: 340000, consecutiveFailures: 2 },
    { id: '3', loanId: 'LN-2025-0078', customerId: 'C003', customerName: 'Abuja Fresh Mart', mandateAmount: 200000, frequency: 'monthly', status: 'active', bankName: 'UBA', accountNumber: '****2345', nextDebitDate: '2026-01-28', remainingBalance: 1400000, totalRepaid: 600000, consecutiveFailures: 0 },
    { id: '4', loanId: 'LN-2025-0091', customerId: 'C004', customerName: 'Port Harcourt Foods', mandateAmount: 150000, frequency: 'monthly', status: 'pending_setup', bankName: 'Access Bank', accountNumber: '****6789', nextDebitDate: '2026-02-01', remainingBalance: 1500000, totalRepaid: 0, consecutiveFailures: 0 },
];

// =============================================================================
// UTILITIES
// =============================================================================

const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
    }).format(amount);
};

const getRiskBandColor = (band: LoanApplication['riskBand']): string => {
    const colors: Record<LoanApplication['riskBand'], string> = {
        prime: '#52c41a',
        near_prime: '#1890ff',
        subprime: '#fa8c16',
        deep_subprime: '#ff4d4f',
    };
    return colors[band];
};

const getStatusSteps = (status: LoanApplication['status']): number => {
    const steps: Record<LoanApplication['status'], number> = {
        new: 0, kyc_pending: 1, under_review: 2, credit_check: 3, approved: 4, disbursed: 5, rejected: -1,
    };
    return steps[status];
};

// =============================================================================
// COMPONENTS
// =============================================================================

const CreditScoreGauge = ({ score, riskBand }: { score: number; riskBand: LoanApplication['riskBand'] }) => {
    const percentage = ((score - 300) / 550) * 100;
    return (
        <div style={{ textAlign: 'center' }}>
            <Progress
                type="dashboard"
                percent={percentage}
                strokeColor={getRiskBandColor(riskBand)}
                format={() => (
                    <div>
                        <Text style={{ fontSize: 28, fontWeight: 700, color: getRiskBandColor(riskBand) }}>{score}</Text>
                        <br />
                        <Tag color={getRiskBandColor(riskBand)}>{riskBand.replace('_', ' ').toUpperCase()}</Tag>
                    </div>
                )}
                size={180}
            />
        </div>
    );
};

const LoanApplicationCard = ({ application, onSelect }: { application: LoanApplication; onSelect: (app: LoanApplication) => void }) => (
    <Card
        hoverable
        onClick={() => onSelect(application)}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 16 }}
    >
        <Row gutter={16} align="middle">
            <Col flex="auto">
                <Space direction="vertical" size={0}>
                    <Space>
                        <Text strong>{application.businessName}</Text>
                        <Tag>{application.applicationNumber}</Tag>
                    </Space>
                    <Text type="secondary">{application.ownerName} • {application.purpose}</Text>
                    <Space style={{ marginTop: 8 }}>
                        <Tag color={getRiskBandColor(application.riskBand)}>{application.creditScore}</Tag>
                        <Text strong>{formatNaira(application.amount)}</Text>
                        <Text type="secondary">• {application.tenure} months</Text>
                    </Space>
                </Space>
            </Col>
            <Col>
                <Progress
                    type="circle"
                    percent={application.transactionHistory.paymentReliability}
                    size={50}
                    format={(p) => `${p}%`}
                    strokeColor={application.transactionHistory.paymentReliability >= 90 ? '#52c41a' : '#fa8c16'}
                />
            </Col>
        </Row>
    </Card>
);

const LoanKanban = ({ applications, onSelect }: { applications: LoanApplication[]; onSelect: (app: LoanApplication) => void }) => {
    const columns = [
        { key: 'new', title: '📥 New Applications', color: '#1890ff' },
        { key: 'kyc_pending', title: '🔍 KYC Review', color: '#722ed1' },
        { key: 'credit_check', title: '📊 Credit Assessment', color: '#fa8c16' },
        { key: 'approved', title: '✅ Approved', color: '#52c41a' },
    ];

    return (
        <Row gutter={16} style={{ overflowX: 'auto' }}>
            {columns.map((col) => {
                const apps = applications.filter(a => a.status === col.key || (col.key === 'kyc_pending' && (a.status === 'kyc_pending' || a.status === 'under_review')));
                return (
                    <Col key={col.key} style={{ minWidth: 320 }}>
                        <Card
                            title={
                                <Space>
                                    <Text style={{ color: col.color }}>{col.title}</Text>
                                    <Badge count={apps.length} style={{ background: col.color }} />
                                </Space>
                            }
                            bodyStyle={{ padding: 12, maxHeight: 500, overflowY: 'auto' }}
                        >
                            {apps.length === 0 ? (
                                <Empty description="No applications" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                            ) : (
                                apps.map((app) => <LoanApplicationCard key={app.id} application={app} onSelect={onSelect} />)
                            )}
                        </Card>
                    </Col>
                );
            })}
        </Row>
    );
};

const ATCCollectionCalendar = ({ mandates }: { mandates: ATCMandate[] }) => {
    const getListData = (value: Dayjs) => {
        const dateStr = value.format('YYYY-MM-DD');
        return mandates.filter(m => m.nextDebitDate === dateStr);
    };

    const dateCellRender = (value: Dayjs) => {
        const listData = getListData(value);
        return listData.length > 0 ? (
            <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                {listData.slice(0, 2).map((item) => (
                    <li key={item.id}>
                        <Badge
                            status={item.status === 'active' ? 'success' : item.status === 'failed' ? 'error' : 'warning'}
                            text={<Text style={{ fontSize: 10 }} ellipsis>{formatNaira(item.mandateAmount)}</Text>}
                        />
                    </li>
                ))}
                {listData.length > 2 && (
                    <li><Text type="secondary" style={{ fontSize: 10 }}>+{listData.length - 2} more</Text></li>
                )}
            </ul>
        ) : null;
    };

    return (
        <Card title={<Space><CalendarOutlined /> Collection Calendar</Space>}>
            <Calendar
                fullscreen={false}
                cellRender={(current, info) => info.type === 'date' ? dateCellRender(current) : null}
            />
        </Card>
    );
};

const ATCMandateTable = ({ mandates }: { mandates: ATCMandate[] }) => {
    const columns: ColumnsType<ATCMandate> = [
        {
            title: 'Customer',
            dataIndex: 'customerName',
            key: 'customerName',
            render: (name, record) => (
                <Space>
                    <Avatar style={{ background: record.status === 'active' ? '#52c41a' : record.status === 'failed' ? '#ff4d4f' : '#fa8c16' }}>
                        {name[0]}
                    </Avatar>
                    <div>
                        <Text strong>{name}</Text>
                        <br />
                        <Text type="secondary" style={{ fontSize: 12 }}>{record.loanId}</Text>
                    </div>
                </Space>
            ),
        },
        {
            title: 'Mandate',
            key: 'mandate',
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{formatNaira(record.mandateAmount)}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>{record.frequency}</Text>
                </Space>
            ),
        },
        {
            title: 'Bank Account',
            key: 'bank',
            render: (_, record) => `${record.bankName} ${record.accountNumber}`,
        },
        {
            title: 'Progress',
            key: 'progress',
            render: (_, record) => {
                const total = record.remainingBalance + record.totalRepaid;
                const progress = (record.totalRepaid / total) * 100;
                return (
                    <Tooltip title={`${formatNaira(record.totalRepaid)} of ${formatNaira(total)} repaid`}>
                        <Progress percent={progress} size="small" format={() => `${progress.toFixed(0)}%`} />
                    </Tooltip>
                );
            },
        },
        {
            title: 'Next Debit',
            dataIndex: 'nextDebitDate',
            key: 'nextDebitDate',
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status, record) => (
                <Space direction="vertical" size={0}>
                    <Tag color={status === 'active' ? 'success' : status === 'failed' ? 'error' : status === 'pending_setup' ? 'processing' : 'default'}>
                        {status.replace('_', ' ').toUpperCase()}
                    </Tag>
                    {record.consecutiveFailures > 0 && (
                        <Text type="danger" style={{ fontSize: 11 }}>
                            <WarningOutlined /> {record.consecutiveFailures} failures
                        </Text>
                    )}
                </Space>
            ),
        },
        {
            title: 'Actions',
            key: 'actions',
            render: (_, record) => (
                <Space>
                    {record.status === 'failed' && (
                        <Button size="small" type="primary" icon={<SyncOutlined />}>Retry</Button>
                    )}
                    <Button size="small" icon={<FileSearchOutlined />}>Details</Button>
                </Space>
            ),
        },
    ];

    return <Table columns={columns} dataSource={mandates} rowKey="id" pagination={false} />;
};

// =============================================================================
// MAIN COMPONENT
// =============================================================================

export default function BankDashboard() {
    const [activeTab, setActiveTab] = useState('overview');
    const [selectedApplication, setSelectedApplication] = useState<LoanApplication | null>(null);

    const portfolioStats = {
        totalDisbursed: 245600000,
        activeLoans: 156,
        collectionsToday: 12500000,
        collectionsCollected: 11750000,
        defaultRate: 2.3,
        pendingApplications: 23,
    };

    return (
        <div style={{ padding: 24, background: '#f0f2f5', minHeight: '100vh' }}>
            {/* Header */}
            <div style={{ marginBottom: 24 }}>
                <Title level={2} style={{ margin: 0 }}>
                    <BankOutlined style={{ marginRight: 12, color: '#1890ff' }} />
                    Bank Command Center
                </Title>
                <Text type="secondary">Loan portfolio management, ATC collections, and compliance</Text>
            </div>

            {/* Key Metrics */}
            <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic title="Total Disbursed" value={formatNaira(portfolioStats.totalDisbursed)} valueStyle={{ color: '#1890ff', fontSize: 20 }} />
                    </Card>
                </Col>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic title="Active Loans" value={portfolioStats.activeLoans} prefix={<FundOutlined />} />
                    </Card>
                </Col>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic
                            title="Today's Collections"
                            value={((portfolioStats.collectionsCollected / portfolioStats.collectionsToday) * 100).toFixed(1)}
                            suffix="%"
                            valueStyle={{ color: '#52c41a' }}
                            prefix={<RiseOutlined />}
                        />
                        <Progress percent={(portfolioStats.collectionsCollected / portfolioStats.collectionsToday) * 100} showInfo={false} strokeColor="#52c41a" />
                    </Card>
                </Col>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic title="Default Rate" value={portfolioStats.defaultRate} suffix="%" valueStyle={{ color: '#ff4d4f' }} prefix={<WarningOutlined />} />
                    </Card>
                </Col>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic title="Pending Review" value={portfolioStats.pendingApplications} prefix={<ClockCircleOutlined />} valueStyle={{ color: '#fa8c16' }} />
                    </Card>
                </Col>
                <Col xs={12} sm={8} lg={4}>
                    <Card>
                        <Statistic title="AI Risk Alerts" value={5} prefix={<RobotOutlined />} valueStyle={{ color: '#722ed1' }} />
                    </Card>
                </Col>
            </Row>

            {/* Main Tabs */}
            <Tabs activeKey={activeTab} onChange={setActiveTab}>
                <TabPane tab={<span><FundOutlined /> Loan Applications</span>} key="applications">
                    <LoanKanban applications={loanApplications} onSelect={setSelectedApplication} />
                </TabPane>
                <TabPane tab={<span><CreditCardOutlined /> ATC Mandates</span>} key="atc">
                    <Row gutter={[16, 16]}>
                        <Col xs={24} lg={16}>
                            <Card title="Active Mandates">
                                <ATCMandateTable mandates={atcMandates} />
                            </Card>
                        </Col>
                        <Col xs={24} lg={8}>
                            <ATCCollectionCalendar mandates={atcMandates} />
                        </Col>
                    </Row>
                </TabPane>
                <TabPane tab={<span><HistoryOutlined /> Reconciliation</span>} key="reconciliation">
                    <Card title="Daily Reconciliation">
                        <Empty description="Select a date to view reconciliation" />
                    </Card>
                </TabPane>
                <TabPane tab={<span><SafetyOutlined /> Compliance</span>} key="compliance">
                    <Row gutter={[16, 16]}>
                        <Col span={8}>
                            <Card>
                                <Statistic title="KYC Verified" value={98} suffix="%" prefix={<CheckCircleOutlined />} valueStyle={{ color: '#52c41a' }} />
                            </Card>
                        </Col>
                        <Col span={8}>
                            <Card>
                                <Statistic title="AML Alerts" value={3} prefix={<AlertOutlined />} valueStyle={{ color: '#ff4d4f' }} />
                            </Card>
                        </Col>
                        <Col span={8}>
                            <Card>
                                <Statistic title="Pending Reviews" value={12} prefix={<FileProtectOutlined />} />
                            </Card>
                        </Col>
                    </Row>
                </TabPane>
            </Tabs>

            {/* Loan Application Detail Drawer */}
            <Drawer
                title={`Loan Application: ${selectedApplication?.applicationNumber}`}
                open={selectedApplication !== null}
                onClose={() => setSelectedApplication(null)}
                width={640}
            >
                {selectedApplication && (
                    <Space direction="vertical" style={{ width: '100%' }} size="large">
                        <Steps
                            current={getStatusSteps(selectedApplication.status)}
                            size="small"
                            items={[
                                { title: 'New' },
                                { title: 'KYC' },
                                { title: 'Review' },
                                { title: 'Credit Check' },
                                { title: 'Approved' },
                                { title: 'Disbursed' },
                            ]}
                        />

                        <Divider />

                        <Row gutter={24}>
                            <Col span={12}>
                                <Descriptions column={1} size="small">
                                    <Descriptions.Item label="Business">{selectedApplication.businessName}</Descriptions.Item>
                                    <Descriptions.Item label="Owner">{selectedApplication.ownerName}</Descriptions.Item>
                                    <Descriptions.Item label="Phone">{selectedApplication.phone}</Descriptions.Item>
                                    <Descriptions.Item label="Email">{selectedApplication.email}</Descriptions.Item>
                                    <Descriptions.Item label="Purpose">{selectedApplication.purpose}</Descriptions.Item>
                                </Descriptions>
                            </Col>
                            <Col span={12}>
                                <CreditScoreGauge score={selectedApplication.creditScore} riskBand={selectedApplication.riskBand} />
                            </Col>
                        </Row>

                        <Card title="Transaction History on OmniRoute" size="small">
                            <Row gutter={16}>
                                <Col span={6}>
                                    <Statistic title="Orders (90d)" value={selectedApplication.transactionHistory.ordersLast90Days} />
                                </Col>
                                <Col span={6}>
                                    <Statistic title="Total GMV" value={formatNaira(selectedApplication.transactionHistory.totalGMV)} />
                                </Col>
                                <Col span={6}>
                                    <Statistic title="Avg Order" value={formatNaira(selectedApplication.transactionHistory.avgOrderValue)} />
                                </Col>
                                <Col span={6}>
                                    <Statistic title="Reliability" value={selectedApplication.transactionHistory.paymentReliability} suffix="%" valueStyle={{ color: '#52c41a' }} />
                                </Col>
                            </Row>
                        </Card>

                        <Space>
                            <Button type="primary" size="large">Approve Loan</Button>
                            <Button danger size="large">Reject</Button>
                            <Button size="large">Request More Info</Button>
                        </Space>
                    </Space>
                )}
            </Drawer>
        </div>
    );
}
