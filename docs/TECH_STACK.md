# OmniRoute Commerce Platform - Tech Stack & Architecture

## Executive Summary

OmniRoute is a **comprehensive B2B/B2C commerce and logistics platform** designed for emerging markets, particularly Africa. It combines e-commerce, gig economy workforce management, financial services, and last-mile logistics into a unified platform.

---

## Total Lines of Code

| Language | Lines | Files | Purpose |
|----------|-------|-------|---------|
| **Go** | 20,377 | 49 | Backend microservices |
| **Python** | 493 | 2 | Route optimization (OR-Tools) |
| **Dart** | 1,335 | 1 | Mobile app (Flutter) |
| **SQL** | 1,376 | 4 | Database migrations |
| **YAML** | 1,882 | 4 | Configuration, Docker, K8s |
| **Markdown** | 17,268 | 14 | Documentation |
| **Total** | **42,731** | 74 | |

---

## Technology Stack

### 🔧 Backend Services

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.22+ | Primary backend language |
| **Python** | 3.11+ | AI/ML services, route optimization |
| **Gin** | 1.9+ | HTTP framework |
| **gRPC** | 1.60+ | Service-to-service communication |

### 🗄️ Data Layer

| Technology | Version | Purpose |
|------------|---------|---------|
| **PostgreSQL** | 16 | Primary database (development) |
| **YugabyteDB** | 2.20 | Distributed database (production) |
| **Redis** | 7 | Caching (development) |
| **DragonflyDB** | 1.13 | High-performance cache (production) |
| **TimescaleDB** | - | Time-series data (telemetry) |

### 📨 Messaging & Events

| Technology | Version | Purpose |
|------------|---------|---------|
| **Apache Kafka** | 3.6 | Event streaming (development) |
| **Redpanda** | 23.3 | Kafka-compatible (production) |
| **Redis Streams** | - | Real-time events |

### 🌐 API Layer

| Technology | Version | Purpose |
|------------|---------|---------|
| **Hasura** | 2.36 | GraphQL API gateway |
| **OpenAPI** | 3.1 | REST API specification |
| **GraphQL** | - | Flexible data queries |

### ⚙️ Orchestration & Workflows

| Technology | Version | Purpose |
|------------|---------|---------|
| **Temporal** | 1.22 | Workflow orchestration |
| **n8n** | 1.22 | No-code workflow automation |
| **MCP** | 2024-11 | Model Context Protocol for AI |
| **Kubernetes** | 1.28+ | Container orchestration |
| **Docker** | 24+ | Containerization |
| **Helm** | 3.x | K8s package management |

### 📊 Observability

| Technology | Version | Purpose |
|------------|---------|---------|
| **OpenTelemetry** | 1.22 | Distributed tracing |
| **Prometheus** | 2.48 | Metrics collection |
| **Grafana** | 10.2 | Dashboards & visualization |
| **Jaeger** | 1.52 | Trace visualization |
| **Zap** | 1.26 | Structured logging |

### 🤖 AI/ML

| Technology | Purpose |
|------------|---------|
| **Google OR-Tools** | Route optimization (VRP) |
| **Prophet** | Demand forecasting |
| **XGBoost** | Feature-based predictions |
| **PyTorch** | Deep learning models |
| **Feast** | Feature store |
| **MLflow** | ML experiment tracking |

### 📱 Mobile

| Technology | Purpose |
|------------|---------|
| **Flutter** | Cross-platform mobile apps |
| **Dart** | Mobile app language |

### 🔐 Security

| Technology | Purpose |
|------------|---------|
| **JWT** | Authentication tokens |
| **OAuth 2.0** | Authorization |
| **TLS 1.3** | Transport encryption |
| **Argon2** | Password hashing |

### ☁️ Infrastructure

| Technology | Purpose |
|------------|---------|
| **Terraform** | Infrastructure as Code |
| **GitHub Actions** | CI/CD pipelines |
| **MinIO** | Object storage (S3-compatible) |
| **Nginx** | Reverse proxy, load balancing |

---

## Microservices Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           API GATEWAY (Hasura GraphQL)                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Pricing    │  │     Gig      │  │ Notification │  │   Payment    │    │
│  │   Engine     │  │  Platform    │  │   Service    │  │   Service    │    │
│  │   :8081      │  │   :8082      │  │    :8083     │  │    :8084     │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │     SCE      │  │    Bank      │  │     ATC      │  │    Route     │    │
│  │   Service    │  │   Gateway    │  │   Service    │  │  Optimizer   │    │
│  │   :8085      │  │   :8086      │  │    :8087     │  │    :8088     │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                         DATA LAYER                                           │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────────────────────┐ │
│  │  YugabyteDB    │  │  DragonflyDB   │  │         Redpanda              │ │
│  │  (PostgreSQL)  │  │   (Redis)      │  │    (Event Streaming)          │ │
│  └────────────────┘  └────────────────┘  └────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Services Overview

| Service | Port | Language | Responsibility |
|---------|------|----------|----------------|
| **pricing-engine** | 8081 | Go | Dynamic pricing, promotions, discounts |
| **gig-platform** | 8082 | Go | Worker management, task assignment |
| **notification-service** | 8083 | Go | Multi-channel notifications |
| **payment-service** | 8084 | Go | Payments, invoicing, credit |
| **sce-service** | 8085 | Go | Service Creation Environment |
| **bank-gateway** | 8086 | Go | Bank integrations, virtual accounts |
| **atc-service** | 8087 | Go | Authority to Collect (B2B) |
| **route-optimizer** | 8088 | Python | VRP solving, route planning |
| **wms-service** | 8089 | Go | Warehouse management |

---

## Shared Packages

| Package | Purpose |
|---------|---------|
| `pkg/config` | Centralized configuration |
| `pkg/cache` | DragonflyDB/Redis client |
| `pkg/database` | YugabyteDB/PostgreSQL client |
| `pkg/messaging` | Redpanda/Kafka client |
| `pkg/telemetry` | OpenTelemetry tracing |
| `pkg/migrate` | Database migrations |

---

## Design Principles

### 1. Domain-Driven Design (DDD)
- **Aggregates**: Product, Order, Worker, Payment
- **Value Objects**: Money, Location, TimeWindow
- **Domain Events**: OrderCreated, PaymentCompleted
- **Repositories**: Interface-based data access

### 2. Extreme Programming (XP)
- Test-Driven Development
- Continuous Integration
- Pair Programming Ready
- Simple Design

### 3. Legacy Modernization
- Event-Driven Architecture
- Database-per-Service
- API Gateway Pattern
- Strangler Fig Pattern

---

## Deployment Profiles

### Development
```bash
docker-compose up
```
- PostgreSQL, Redis, Kafka
- All services locally

### Distributed (Staging)
```bash
docker-compose --profile distributed up
```
- YugabyteDB, DragonflyDB, Redpanda
- Hasura, Temporal

### Production
```bash
docker-compose --profile full up
```
- Full infrastructure stack
- Multi-region ready
