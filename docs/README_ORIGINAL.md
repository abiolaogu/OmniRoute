# 🚀 OmniRoute Commerce Platform

> **The Commerce Operating System for Emerging Markets**

OmniRoute is a unified commerce and distribution platform that connects manufacturers, distributors, wholesalers, retailers, logistics partners, finance providers, and gig workers in a single ecosystem.

---

## 📋 Table of Contents

- [Vision](#-vision)
- [Key Features](#-key-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [Documentation](#-documentation)
- [Technology Stack](#-technology-stack)
- [Project Structure](#-project-structure)
- [API Reference](#-api-reference)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Vision

Traditional commerce platforms solve only one dimension of the problem. OmniRoute is different:

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   TRADITIONAL                    OMNIROUTE                      │
│   ───────────                    ─────────                      │
│                                                                 │
│   Shopify = D2C only             ALL CHANNELS UNIFIED           │
│   TradeDepot = B2B marketplace   SAAS (You own customers)       │
│   FieldAssist = DMS only         COMMERCE + DISTRIBUTION        │
│   SAP = Too expensive            AFFORDABLE & FAST              │
│                                                                 │
│   Result: Fragmented             Result: ONE PLATFORM           │
│   experience, data silos         Complete visibility            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## ✨ Key Features

### For Manufacturers
- 📊 **Secondary Sales Visibility** - See what's selling at retail level
- 💰 **Multi-tier Pricing** - Different prices for different customer types
- 🎯 **Trade Promotions** - Run and track promotional campaigns
- 📈 **Demand Forecasting** - AI-powered demand prediction

### For Distributors
- 🚚 **Route Optimization** - Smart beat planning for sales reps
- 📱 **Van Sales App** - Mobile POS with offline capability
- 💳 **Credit Management** - Automated credit limits and collections
- 📦 **Inventory Sync** - Real-time stock visibility across locations

### For Retailers
- 🛒 **Easy Ordering** - Web, mobile, WhatsApp, USSD, or voice
- 💵 **Trade Credit** - Buy now, pay later
- 📋 **Smart Reorder** - AI suggests what to order
- 🏪 **Group Buying** - Join buying groups for better prices

### For Logistics Partners
- 🗺️ **Load Matching** - Connect with available shipments
- 📍 **Real-time Tracking** - GPS-enabled fleet visibility
- 💰 **Instant Payments** - Get paid when deliveries complete
- 📊 **Performance Analytics** - Track efficiency metrics

### For Finance Partners
- 📊 **Transaction Data** - Rich data for underwriting
- 🔒 **Secured Lending** - Loans secured by platform transactions
- 📈 **Portfolio Monitoring** - Real-time borrower health
- 💳 **Embedded Products** - BNPL, invoice financing, working capital

### For Gig Workers
- 🎯 **Multiple Task Types** - Delivery, sales, collections, audits
- 📈 **Career Progression** - Level up from Starter to Master
- 💰 **Instant Earnings** - Get paid same day
- 🎓 **Skills Training** - Free courses and certifications
- 🏥 **Benefits** - Insurance and savings programs

---

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │Consumer │ │Retailer │ │Sales Rep│ │Gig      │ │ Admin   │ │WhatsApp │  │
│  │App      │ │Portal   │ │App      │ │Worker   │ │Dashboard│ │Bot      │  │
│  │(Flutter)│ │(Next.js)│ │(Flutter)│ │(Flutter)│ │(Next.js)│ │(Node)   │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
├─────────────────────────────────────────────────────────────────────────────┤
│                              API GATEWAY (Kong)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           MICROSERVICES (Go/Rust)                           │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │Catalog  │ │Pricing  │ │Order    │ │Inventory│ │Customer │ │Gig      │  │
│  │Service  │ │Engine   │ │Service  │ │Service  │ │Service  │ │Platform │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │Payment  │ │Credit   │ │Route    │ │Delivery │ │Notifi-  │ │Analytics│  │
│  │Gateway  │ │Engine   │ │Optimizer│ │Tracking │ │cation   │ │Service  │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
├─────────────────────────────────────────────────────────────────────────────┤
│                         EVENT STREAMING (Kafka)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                              DATA LAYER                                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │PostgreSQL│ │Redis   │ │Timescale│ │Elastic  │ │ClickHse │ │MinIO    │  │
│  │(Primary)│ │(Cache)  │ │(Metrics)│ │(Search) │ │(OLAP)   │ │(Storage)│  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Local Development

```bash
# Clone the repository
git clone https://github.com/omniroute/platform.git
cd platform

# Start infrastructure
docker-compose up -d postgres redis kafka

# Run database migrations
make migrate

# Start the pricing engine
cd services/pricing-engine
go run cmd/server/main.go

# Start the API gateway
cd ../api-gateway
npm install
npm run dev

# Start the retailer portal
cd ../apps/retailer-portal
npm install
npm run dev
```

### Using Docker Compose (Full Stack)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Access services
# - API Gateway: http://localhost:8000
# - Retailer Portal: http://localhost:3000
# - Admin Dashboard: http://localhost:3001
```

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [Product Roadmap](./docs/PRODUCT_ROADMAP.md) | 2025-2027 feature roadmap |
| [Competitive Positioning](./docs/COMPETITIVE_POSITIONING.md) | Market analysis and differentiation |
| [API Design](./docs/API_DESIGN.md) | REST API specification |
| [Database Schema](./docs/DATABASE_SCHEMA.sql) | PostgreSQL schema |
| [Credit & Payment Module](./docs/CREDIT_PAYMENT_MODULE.md) | Financial services architecture |
| [Extreme Innovation](./docs/EXTREME_INNOVATION.md) | Breakthrough features |

---

## 🛠 Technology Stack

### Backend Services
| Service | Language | Framework |
|---------|----------|-----------|
| Pricing Engine | Go | Standard library |
| Order Service | Go | Chi |
| Inventory Service | Go | Chi |
| Payment Service | Go | Chi |
| Route Optimizer | Rust | Actix-web |
| Analytics Service | Python | FastAPI |
| Notification Service | Go | Chi |

### Frontend Applications
| Application | Framework | UI Library |
|-------------|-----------|------------|
| Retailer Portal | Next.js 14 | shadcn/ui |
| Admin Dashboard | Next.js 14 | shadcn/ui |
| Consumer App | Flutter | Material 3 |
| Sales Rep App | Flutter | Material 3 |
| Gig Worker App | Flutter | Material 3 |
| WhatsApp Bot | Node.js | Baileys |

### Infrastructure
| Component | Technology |
|-----------|------------|
| Container Orchestration | Kubernetes |
| API Gateway | Kong |
| Service Mesh | Istio |
| CI/CD | GitHub Actions |
| Infrastructure as Code | Terraform |
| Monitoring | Prometheus + Grafana |
| Logging | Loki |
| Tracing | Jaeger |

---

## 📁 Project Structure

```
omniroute/
├── services/                    # Backend microservices
│   ├── pricing-engine/          # Multi-tier pricing calculations
│   │   ├── cmd/server/          # Entry point
│   │   ├── internal/
│   │   │   ├── api/             # HTTP handlers
│   │   │   ├── domain/          # Domain models
│   │   │   ├── engine/          # Core pricing logic
│   │   │   ├── repository/      # Data access
│   │   │   └── cache/           # Caching layer
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── order-service/           # Order management
│   ├── inventory-service/       # Stock management
│   ├── customer-service/        # Customer & credit
│   ├── payment-service/         # Payment processing
│   ├── gig-platform/            # Gig worker management
│   ├── route-optimizer/         # Route optimization (Rust)
│   └── analytics-service/       # Reporting (Python)
│
├── apps/                        # Frontend applications
│   ├── retailer-portal/         # Next.js B2B portal
│   ├── admin-dashboard/         # Next.js admin
│   ├── consumer-app/            # Flutter consumer app
│   ├── sales-rep-app/           # Flutter sales rep app
│   ├── gig-worker-app/          # Flutter gig worker app
│   └── whatsapp-bot/            # Node.js WhatsApp bot
│
├── packages/                    # Shared libraries
│   ├── ui/                      # Shared UI components
│   ├── api-client/              # Generated API clients
│   └── common/                  # Shared utilities
│
├── infrastructure/              # IaC and deployment
│   ├── terraform/               # AWS/GCP infrastructure
│   ├── kubernetes/              # K8s manifests
│   ├── helm/                    # Helm charts
│   └── docker/                  # Docker configurations
│
├── docs/                        # Documentation
│   ├── api/                     # API documentation
│   ├── architecture/            # Architecture diagrams
│   └── guides/                  # Developer guides
│
├── scripts/                     # Utility scripts
├── docker-compose.yml           # Local development
├── Makefile                     # Build automation
└── README.md
```

---

## 📡 API Reference

### Base URL
```
Production: https://api.omniroute.com/v1
Staging:    https://api.staging.omniroute.com/v1
```

### Authentication
```bash
curl -X GET https://api.omniroute.com/v1/products \
  -H "Authorization: Bearer sk_live_xxxxx" \
  -H "X-Tenant-ID: tnnt_xxxxx"
```

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| **Products** |||
| GET | /products | List products |
| POST | /products | Create product |
| GET | /products/:id | Get product |
| PUT | /products/:id | Update product |
| **Orders** |||
| GET | /orders | List orders |
| POST | /orders | Create order |
| GET | /orders/:id | Get order |
| POST | /orders/:id/confirm | Confirm order |
| **Customers** |||
| GET | /customers | List customers |
| POST | /customers | Create customer |
| GET | /customers/:id/credit | Get credit info |
| **Pricing** |||
| POST | /prices/calculate | Calculate prices |
| GET | /prices | Get single price |
| POST | /prices/bulk | Bulk price lookup |
| **Gig Workers** |||
| GET | /gig-workers | List workers |
| GET | /gig-tasks | List tasks |
| POST | /gig-tasks/:id/complete | Complete task |

See [API Documentation](./docs/API_DESIGN.md) for complete reference.

---

## 🚢 Deployment

### Kubernetes Deployment

```bash
# Apply namespace and configs
kubectl apply -f infrastructure/kubernetes/namespace.yaml
kubectl apply -f infrastructure/kubernetes/configmaps.yaml
kubectl apply -f infrastructure/kubernetes/secrets.yaml

# Deploy services
kubectl apply -f infrastructure/kubernetes/services/

# Deploy ingress
kubectl apply -f infrastructure/kubernetes/ingress.yaml

# Check status
kubectl get pods -n omniroute
```

### Helm Deployment

```bash
# Add OmniRoute Helm repo
helm repo add omniroute https://charts.omniroute.com

# Install
helm install omniroute omniroute/platform \
  --namespace omniroute \
  --create-namespace \
  --values values.yaml
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `REDIS_URL` | Redis connection string | Yes |
| `KAFKA_BROKERS` | Kafka broker addresses | Yes |
| `JWT_SECRET` | JWT signing secret | Yes |
| `PAYSTACK_SECRET_KEY` | Paystack API key | Yes |

---

## 🧪 Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run e2e tests
make test-e2e

# Generate coverage report
make coverage
```

---

## 📈 Monitoring

### Metrics
- Prometheus metrics available at `/metrics` on all services
- Pre-built Grafana dashboards in `infrastructure/grafana/`

### Logging
- Structured JSON logging to stdout
- Aggregated via Loki

### Tracing
- OpenTelemetry instrumentation
- Jaeger for distributed tracing

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Go: Follow [Effective Go](https://golang.org/doc/effective_go)
- TypeScript: ESLint + Prettier
- Python: Black + isort

---

## 📄 License

This project is licensed under the MIT License - see [LICENSE](./LICENSE) for details.

---

## 🙏 Acknowledgments

- All our early customers who believed in the vision
- The open-source community for amazing tools
- Our investors for their support

---

## 📞 Contact

- **Website**: [omniroute.com](https://omniroute.com)
- **Email**: hello@omniroute.com
- **Twitter**: [@omniroute](https://twitter.com/omniroute)
- **LinkedIn**: [OmniRoute](https://linkedin.com/company/omniroute)

---

<p align="center">
  <strong>Built with ❤️ for Africa, by Africa</strong>
</p>
