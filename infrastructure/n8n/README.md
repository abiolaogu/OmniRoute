# n8n Workflow Automation Configuration

## Overview

n8n is integrated as the workflow automation engine for OmniRoute, enabling no-code/low-code automation of business processes.

## Deployment

n8n is deployed as part of the docker-compose stack with persistent storage for workflows.

## Pre-built Workflows

### 1. Order Processing Workflow
- Trigger: New order webhook from Hasura
- Actions:
  - Validate inventory
  - Process payment
  - Assign to gig worker
  - Send notification

### 2. Low Stock Alert Workflow
- Trigger: Scheduled (every hour)
- Actions:
  - Check inventory levels
  - Generate reorder suggestions
  - Create purchase orders
  - Notify procurement

### 3. Payment Reconciliation Workflow
- Trigger: Daily schedule
- Actions:
  - Fetch bank statements
  - Match with invoices
  - Flag discrepancies
  - Send report

### 4. Customer Onboarding Workflow
- Trigger: New customer registration
- Actions:
  - Send welcome email
  - Create credit account
  - Assign sales rep
  - Schedule follow-up

### 5. Worker Payout Workflow
- Trigger: Weekly schedule
- Actions:
  - Calculate earnings
  - Process bank transfers
  - Send payment notifications
  - Update ledger

## Environment Variables

```env
N8N_HOST=n8n.omniroute.local
N8N_PORT=5678
N8N_PROTOCOL=http
N8N_ENCRYPTION_KEY=your-encryption-key
WEBHOOK_URL=https://api.omniroute.io/webhooks/n8n
```

## Integration Points

| Service | Webhook URL | Purpose |
|---------|-------------|---------|
| Order Service | /webhook/order-created | Trigger order workflows |
| Payment Service | /webhook/payment-completed | Update order status |
| Inventory Service | /webhook/stock-low | Trigger reorder |
| Notification Service | /webhook/send-notification | Multi-channel delivery |
