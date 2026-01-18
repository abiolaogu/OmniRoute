# OmniRoute Platform - Technical Architecture

## Executive Summary

OmniRoute is a B2B FMCG commerce platform that combines:
- **Multi-channel Commerce**: Web, Mobile, WhatsApp, USSD, Voice ordering
- **Embedded Finance**: Trade credit, digital wallets, instant payments
- **Gig Economy Logistics**: "Uber for Delivery" with intelligent routing
- **Real-time Operations**: Event-driven architecture with Kafka

Built on proven open-source technologies and leveraging existing BillyRonks repositories.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                  API GATEWAY                                     │
│                              (Kong / Traefik)                                   │
│     ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐            │
│     │   Web   │  │ Mobile  │  │WhatsApp │  │  USSD   │  │  API    │            │
│     │ Portal  │  │   App   │  │  Bot    │  │ Gateway │  │ Clients │            │
│     └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘            │
└──────────┼───────────┼───────────┼───────────┼───────────┼──────────────────────┘
           │           │           │           │           │
           └───────────┴───────────┴─────┬─────┴───────────┘
                                         │
           ┌─────────────────────────────┼─────────────────────────────┐
           │                             ▼                             │
           │  ┌─────────────────────────────────────────────────────┐  │
           │  │              MICROSERVICES LAYER                    │  │
           │  │                                                     │  │
           │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │  │
           │  │  │ Pricing  │ │  Order   │ │Inventory │ │Customer│ │  │
           │  │  │ Engine   │ │ Service  │ │ Service  │ │Service │ │  │
           │  │  │  (Go)    │ │  (Go)    │ │  (Go)    │ │  (Go)  │ │  │
           │  │  └────┬─────┘ └────┬─────┘ └────┬─────┘ └───┬────┘ │  │
           │  │       │            │            │           │      │  │
           │  │  ┌────┴──────┐ ┌───┴────┐ ┌────┴─────┐ ┌───┴────┐ │  │
           │  │  │ Payment   │ │  Gig   │ │Notifica- │ │Catalog │ │  │
           │  │  │ Service   │ │Platform│ │tion Svc  │ │Service │ │  │
           │  │  │  (Go)     │ │  (Go)  │ │  (Go)    │ │  (Go)  │ │  │
           │  │  └───────────┘ └────────┘ └──────────┘ └────────┘ │  │
           │  │                                                     │  │
           │  └─────────────────────────────────────────────────────┘  │
           │                             │                             │
           └─────────────────────────────┼─────────────────────────────┘
                                         │
           ┌─────────────────────────────┼─────────────────────────────┐
           │              EVENT BUS (Apache Kafka)                     │
           │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐ │
           │  │ orders   │ │inventory │ │ payments │ │ gig-workers  │ │
           │  └──────────┘ └──────────┘ └──────────┘ └──────────────┘ │
           └─────────────────────────────┬─────────────────────────────┘
                                         │
           ┌─────────────────────────────┼─────────────────────────────┐
           │             EXTERNAL INTEGRATIONS                         │
           │                             │                             │
           │  ┌──────────┐  ┌──────────┐ │ ┌──────────┐ ┌──────────┐  │
           │  │   VAS    │  │opensase- │ │ │ Fineract │ │   Maps   │  │
           │  │Messaging │  │payments  │ │ │ Banking  │ │   API    │  │
           │  │(WhatsApp │  │ (Rust)   │ │ │ (Java)   │ │ (Google) │  │
           │  │SMS,USSD) │  │          │ │ │          │ │          │  │
           │  └──────────┘  └──────────┘ │ └──────────┘ └──────────┘  │
           │                             │                             │
           └─────────────────────────────┼─────────────────────────────┘
                                         │
           ┌─────────────────────────────┼─────────────────────────────┐
           │                  DATA LAYER                               │
           │                             │                             │
           │  ┌──────────┐  ┌──────────┐ │ ┌──────────┐ ┌──────────┐  │
           │  │PostgreSQL│  │  Redis   │ │ │  Druid   │ │  MinIO   │  │
           │  │  + Citus │  │ Cluster  │ │ │Analytics │ │  S3      │  │
           │  └──────────┘  └──────────┘ │ └──────────┘ └──────────┘  │
           │                             │                             │
           └─────────────────────────────┴─────────────────────────────┘
```

---

## Technology Stack

### Core Services (Go 1.22+)
| Service | Purpose | Key Libraries |
|---------|---------|---------------|
| **Pricing Engine** | Dynamic B2B pricing, promotions | decimal, gRPC |
| **Order Service** | Order lifecycle management | CQRS pattern |
| **Inventory Service** | Real-time stock management | Optimistic locking |
| **Customer Service** | Retailer/distributor management | KYC integration |
| **Payment Service** | Payment processing, wallets | Fineract SDK |
| **Gig Platform** | Worker management, task allocation | ML integration |
| **Notification Service** | Multi-channel messaging | VAS integration |
| **Catalog Service** | Product management | Elasticsearch |

### External Integrations (from existing repos)
| System | Source | Purpose |
|--------|--------|---------|
| **VAS Platform** | VAS repo | WhatsApp, SMS, USSD messaging |
| **opensase-payments** | opensase-payments repo | Paystack/Flutterwave integration |
| **Apache Fineract** | Global-FinTech repo | Core banking, credit management |
| **JPOS** | Global-FinTech repo | Card processing (future) |
| **n8n** | eCommerce repo | Workflow orchestration |

### Data Infrastructure
| Technology | Purpose | Configuration |
|------------|---------|---------------|
| **PostgreSQL 16 + Citus** | Transactional data | Sharded by tenant_id |
| **Redis 7 Cluster** | Caching, real-time state | Worker locations |
| **Apache Kafka 3.6** | Event streaming | 6+ partitions per topic |
| **Apache Druid** | Real-time analytics | OLAP queries |
| **MinIO** | Object storage | S3-compatible |

---

## Service Details

### 1. Pricing Engine

**Purpose**: Calculates prices for B2B transactions with complex rules.

**Features**:
- Customer-specific pricing
- Volume discounts
- Time-based promotions
- Regional pricing
- Channel-specific rates

**Key Endpoints**:
```
POST /api/v1/prices/calculate
POST /api/v1/prices/bulk-calculate
GET  /api/v1/price-lists/{id}
POST /api/v1/promotions/validate
```

**Performance**: <10ms P99 latency, 50K+ calculations/second

### 2. Order Service

**Purpose**: Manages B2B order lifecycle from creation to delivery.

**Features**:
- Multi-channel order intake (Web, WhatsApp, USSD, API)
- Trade credit integration
- Inventory reservation
- Fulfillment orchestration
- Partial delivery support

**Event Flow**:
```
order.created → inventory.reserved → order.confirmed
    → order.assigned_to_driver → order.delivered
```

### 3. Gig Platform

**Purpose**: Manages delivery workforce and task assignment.

**Features**:
- Real-time worker tracking
- Intelligent task allocation (nearest, best-rated, AI-optimized)
- Route optimization
- Payment collection
- Instant earnings payout

**Allocation Strategies**:
1. **Nearest**: Assign to closest available worker
2. **Best-Rated**: Prioritize highest-rated workers
3. **Broadcast**: Offer to multiple workers simultaneously
4. **AI-Optimized**: ML-based optimal matching

### 4. Notification Service

**Purpose**: Multi-channel customer communication via VAS platform.

**Channels**:
- WhatsApp Business API (templates, interactive messages)
- SMS (bulk, transactional)
- USSD (interactive menus for ordering)
- Voice/IVR (order reminders)
- Push notifications
- Email

**USSD Menu Structure**:
```
*123# → Main Menu
  1. Place Order → Categories → Products → Quantity → Confirm
  2. Check Order Status → Order list
  3. View Balance → Credit/Wallet info
  4. Make Payment → Payment options
  5. Contact Support → Callback/Call
```

### 5. Payment Service

**Purpose**: Payment processing with embedded finance.

**Features**:
- Payment gateway integration (Paystack, Flutterwave)
- Digital wallets (Fineract savings accounts)
- Trade credit (Fineract loan module)
- Cash collection by gig workers
- Settlement processing

**Credit Assessment Flow**:
```
1. Customer applies for credit
2. System gathers order/payment history
3. AI model scores creditworthiness
4. Auto-approve or refer to manual review
5. Credit limit assigned in Fineract
6. Customer uses credit for orders
7. Repayment tracked and reported
```

---

## Event-Driven Architecture

### Kafka Topics

| Topic | Publishers | Consumers |
|-------|------------|-----------|
| `omniroute.orders` | Order Service | Inventory, Payment, Notification, Analytics |
| `omniroute.inventory` | Inventory Service | Order, Notification, Analytics |
| `omniroute.payments` | Payment Service | Order, Notification, Credit |
| `omniroute.gig-workers` | Gig Platform | Order, Notification, Analytics |
| `omniroute.notifications` | All Services | Notification Service |
| `omniroute.analytics` | All Services | Analytics Service, Druid |

### Event Contracts

All events follow a standard structure:
```json
{
  "id": "evt_uuid",
  "type": "order.created",
  "source": "order-service",
  "tenant_id": "tenant_uuid",
  "data": { ... },
  "metadata": {
    "correlation_id": "corr_uuid",
    "user_id": "user_uuid",
    "channel": "whatsapp"
  },
  "timestamp": "2026-01-18T12:00:00Z",
  "version": 1
}
```

### Event Processing Patterns

1. **Choreography**: Services react to events independently
2. **Saga Pattern**: Multi-step transactions with compensation
3. **CQRS**: Separate read/write models for performance
4. **Event Sourcing**: Full audit trail for financial transactions

---

## Database Design

### Multi-Tenancy Strategy

- **Shared Schema**: All tenants share tables
- **Row-Level Security**: `tenant_id` in every table
- **Citus Sharding**: Distributed by `tenant_id` for scale

### Key Tables

| Schema | Table | Purpose | Sharding |
|--------|-------|---------|----------|
| public | tenants | Tenant configuration | Reference |
| public | customers | Retailer/distributor data | tenant_id |
| public | products | Product catalog | tenant_id |
| public | orders | Order records | tenant_id |
| public | inventory_items | Stock levels | tenant_id |
| public | gig_workers | Delivery workforce | tenant_id |
| public | payments | Payment transactions | tenant_id |
| public | wallets | Digital wallets | tenant_id |
| public | credit_accounts | Trade credit accounts | tenant_id |

### Indexes

Critical indexes for performance:
- `orders(tenant_id, status, created_at DESC)` - Dashboard queries
- `inventory_items(warehouse_id, product_id)` - Stock checks
- `gig_workers(tenant_id, availability)` - Worker allocation
- `payments(tenant_id, status, created_at DESC)` - Payment reports

---

## Deployment Architecture

### Kubernetes Cluster

```
┌────────────────────────────────────────────────────────────────┐
│                    KUBERNETES CLUSTER                          │
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                    INGRESS (Nginx)                       │  │
│  │                 *.omniroute.com                          │  │
│  └─────────────────────────────────────────────────────────┘  │
│                              │                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                 API GATEWAY (Kong)                       │  │
│  │     Rate limiting, Auth, Routing, Observability         │  │
│  └─────────────────────────────────────────────────────────┘  │
│                              │                                 │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │ pricing  │ │  order   │ │inventory │ │   Other Services │  │
│  │ engine   │ │ service  │ │ service  │ │                  │  │
│  │ 3 pods   │ │ 5 pods   │ │ 3 pods   │ │                  │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘  │
│                              │                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                    DATA PLANE                            │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────┐ │  │
│  │  │Postgres │  │  Redis  │  │  Kafka  │  │   MinIO     │ │  │
│  │  │HA Proxy │  │ Sentinel│  │ Cluster │  │ (S3)        │ │  │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────────┘ │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                 OBSERVABILITY                            │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────┐ │  │
│  │  │Prometheus│ │ Grafana │  │  Jaeger │  │   Loki      │ │  │
│  │  │ metrics │  │  dashb  │  │  traces │  │   logs      │ │  │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────────┘ │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

### Resource Allocation

| Service | CPU | Memory | Replicas | HPA Target |
|---------|-----|--------|----------|------------|
| Pricing Engine | 500m | 512Mi | 3-10 | 70% CPU |
| Order Service | 1000m | 1Gi | 5-20 | 70% CPU |
| Inventory Service | 500m | 512Mi | 3-10 | 70% CPU |
| Gig Platform | 1000m | 1Gi | 5-15 | 70% CPU |
| Notification Service | 500m | 512Mi | 3-10 | Queue depth |
| Payment Service | 500m | 512Mi | 3-8 | 70% CPU |

---

## Security Architecture

### Authentication & Authorization

1. **API Gateway**: JWT validation, rate limiting
2. **Service-to-Service**: mTLS with service mesh
3. **User Auth**: Keycloak (OAuth2/OIDC)
4. **API Keys**: For external integrations

### Data Security

1. **Encryption at Rest**: PostgreSQL TDE, Redis encryption
2. **Encryption in Transit**: TLS 1.3 everywhere
3. **Secrets Management**: HashiCorp Vault
4. **PII Handling**: Tokenization, masking

### Compliance

1. **NDPR**: Nigeria Data Protection Regulation
2. **PCI DSS**: For payment processing
3. **SOC 2**: Audit trails, access controls

---

## Performance Targets

| Metric | Target | Current Architecture Support |
|--------|--------|------------------------------|
| API Latency (P99) | <200ms | Caching, read replicas |
| Order Processing | <5s end-to-end | Event-driven, async |
| Message Delivery | <30s | VAS platform SLA |
| Payment Completion | <10s | Direct provider integration |
| Inventory Update | Real-time | Kafka streaming |
| Worker Location | <5s updates | Redis pub/sub |
| System Availability | 99.9% | Multi-AZ, auto-scaling |
| TPS (peak) | 10,000+ | Horizontal scaling |

---

## Development Workflow

### Local Development

```bash
# Start infrastructure
docker-compose up -d postgres redis kafka

# Start services
make run-pricing-engine
make run-order-service
# ... etc

# Or start everything
docker-compose up
```

### Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Load tests
make test-load
```

### CI/CD Pipeline (GitHub Actions + ArgoCD)

```
Push → Build → Test → Scan → Push Image → Update Manifests → ArgoCD Sync
```

---

## Monitoring & Observability

### Metrics (Prometheus)

- Request rates, latencies, errors
- Business metrics (orders, revenue)
- Resource utilization

### Logging (Loki)

- Structured JSON logs
- Correlation IDs for tracing
- 30-day retention

### Tracing (Jaeger)

- Distributed tracing
- OpenTelemetry integration
- Service dependency mapping

### Alerting

| Alert | Condition | Severity |
|-------|-----------|----------|
| High Error Rate | >1% 5xx errors | Critical |
| Slow Response | P99 >500ms | Warning |
| Low Inventory | Stock <reorder point | Warning |
| Payment Failure | >5% failure rate | Critical |
| Worker Shortage | <10% available | Warning |

---

## Repository Integration Summary

| Repository | Components Used | Integration Method |
|------------|-----------------|-------------------|
| **VAS** | WhatsApp, SMS, USSD modules | HTTP API |
| **Global-FinTech** | Fineract, JPOS, Hyperledger | SDK/API |
| **eCommerce** | Kafka events, n8n workflows | Event contracts |
| **opensase-payments** | Payment gateway integration | HTTP API |
| **AI-Agents** | Credit scoring, support bot | gRPC/HTTP |
| **Northflank-Alternative** | Deployment scripts | CI/CD pipeline |

---

## Next Steps

### Phase 1: Foundation (Months 1-2)
- [ ] Deploy core services (Pricing, Order, Inventory)
- [ ] Integrate opensase-payments
- [ ] Set up Kafka event bus
- [ ] Configure VAS messaging

### Phase 2: Finance (Months 3-4)
- [ ] Deploy Fineract for wallets/credit
- [ ] Implement credit scoring
- [ ] Build settlement system
- [ ] Mobile money integration

### Phase 3: Gig Platform (Months 5-6)
- [ ] Launch gig worker mobile app
- [ ] Implement task allocation engine
- [ ] Add route optimization
- [ ] Enable instant payouts

### Phase 4: Scale (Months 7-12)
- [ ] Multi-region deployment
- [ ] AI-powered features
- [ ] Additional markets
- [ ] Partner integrations

---

## Contact

- **Architecture**: architecture@omniroute.com
- **DevOps**: devops@omniroute.com
- **Support**: support@omniroute.com
