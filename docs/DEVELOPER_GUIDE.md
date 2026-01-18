# OmniRoute - Developer Guide

## Getting Started

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.22+ | Backend development |
| Python | 3.11+ | ML services |
| Docker | 24+ | Containerization |
| Docker Compose | 2.20+ | Local orchestration |
| Node.js | 20+ | Frontend tools |
| Flutter | 3.16+ | Mobile development |

### Quick Start

```bash
# Clone repository
git clone https://github.com/omniroute/omniroute.git
cd omniroute

# Start infrastructure
docker-compose up -d postgres redis kafka

# Run migrations
make migrate-up

# Start a service (e.g., pricing-engine)
cd services/pricing-engine
go run cmd/server/main.go
```

---

## Project Structure

```
omniroute/
├── apps/                       # Client applications
│   └── retailer-mobile/        # Flutter mobile app
├── docs/                       # Documentation
├── infrastructure/             # IaC configurations
│   ├── docker/                 # Dockerfiles
│   ├── kubernetes/             # K8s manifests
│   └── terraform/              # Terraform modules
├── pkg/                        # Shared Go packages
│   ├── cache/                  # DragonflyDB client
│   ├── config/                 # Configuration
│   ├── database/               # YugabyteDB client
│   ├── messaging/              # Redpanda client
│   ├── migrate/                # Migration runner
│   └── telemetry/              # OpenTelemetry
├── services/                   # Microservices
│   ├── pricing-engine/         # Go
│   ├── gig-platform/           # Go
│   ├── notification-service/   # Go
│   ├── payment-service/        # Go
│   ├── sce-service/            # Go
│   ├── bank-gateway/           # Go
│   ├── atc-service/            # Go
│   ├── route-optimizer/        # Python
│   └── wms-service/            # Go
├── docker-compose.yml          # Local development
├── Makefile                    # Build automation
└── README.md                   # Project overview
```

---

## Service Structure (Go)

Each Go service follows this structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
├── internal/
│   ├── api/
│   │   └── handlers.go         # HTTP handlers
│   ├── domain/
│   │   ├── models.go           # Domain models
│   │   ├── repository.go       # Repository interfaces
│   │   └── validation.go       # Validation methods
│   ├── repository/
│   │   └── postgres.go         # Repository implementation
│   ├── service/
│   │   └── service.go          # Business logic
│   └── cache/
│       └── redis.go            # Cache implementation
├── migrations/
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── go.mod
└── go.sum
```

---

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/my-feature
```

### 2. Implement with TDD

```go
// 1. Write failing test
func TestCalculatePrice_WithDiscount_ReturnsDiscountedPrice(t *testing.T) {
    engine := NewPricingEngine()
    price := engine.Calculate(product, customer)
    assert.Equal(t, expected, price)
}

// 2. Write minimal code to pass
func (e *PricingEngine) Calculate(p Product, c Customer) Money {
    return p.BasePrice.Subtract(p.Discount)
}

// 3. Refactor
```

### 3. Run Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

### 4. Lint and Format

```bash
# Format
gofmt -w .

# Lint
golangci-lint run

# Security scan
gosec ./...
```

### 5. Create Pull Request

```bash
git push origin feature/my-feature
# Create PR via GitHub
```

---

## API Development

### GraphQL (Hasura)

Most APIs are exposed via Hasura GraphQL:

```graphql
# Query
query GetProducts($tenantId: uuid!) {
  products(where: { tenant_id: { _eq: $tenantId } }) {
    id
    name
    price
    variants {
      id
      sku
      price
    }
  }
}

# Mutation via Action
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
    status
    total
  }
}
```

### REST (Service Direct)

Internal and webhook endpoints use REST:

```go
func (h *Handler) CreatePayment(c *gin.Context) {
    var req CreatePaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    payment, err := h.service.CreatePayment(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, payment)
}
```

---

## Database Migrations

### Create Migration

```bash
# Create new migration
make migrate-create NAME=add_customer_credit_limit
```

This creates:
- `migrations/002_add_customer_credit_limit.up.sql`
- `migrations/002_add_customer_credit_limit.down.sql`

### Write Migration

```sql
-- 002_add_customer_credit_limit.up.sql
ALTER TABLE customers 
ADD COLUMN credit_limit DECIMAL(15,2) DEFAULT 0.00;

CREATE INDEX idx_customers_credit_limit 
ON customers(tenant_id, credit_limit);

-- 002_add_customer_credit_limit.down.sql
DROP INDEX IF EXISTS idx_customers_credit_limit;
ALTER TABLE customers DROP COLUMN credit_limit;
```

### Run Migrations

```bash
# Apply all pending
make migrate-up

# Rollback last
make migrate-down

# Force to specific version
make migrate-force VERSION=1
```

---

## Adding a New Service

### 1. Create Directory Structure

```bash
mkdir -p services/new-service/{cmd/server,internal/{api,domain,repository,service}}
```

### 2. Initialize Go Module

```bash
cd services/new-service
go mod init github.com/omniroute/new-service
```

### 3. Create Main Entry Point

```go
// cmd/server/main.go
package main

func main() {
    cfg := config.Load()
    logger := logging.NewLogger(cfg.LogLevel)
    
    db, err := database.NewPool(cfg.DatabaseURL)
    if err != nil {
        logger.Fatal("database connection failed", zap.Error(err))
    }
    
    repo := repository.NewPostgres(db)
    svc := service.New(repo)
    handler := api.NewHandler(svc)
    
    router := gin.New()
    router.GET("/health", handler.Health)
    // ... routes
    
    router.Run(":" + cfg.Port)
}
```

### 4. Add to Docker Compose

```yaml
# docker-compose.yml
new-service:
  build:
    context: ./services/new-service
    dockerfile: ../../infrastructure/docker/Dockerfile.go
  ports:
    - "8090:8090"
  environment:
    DATABASE_URL: postgres://...
```

---

## Testing Strategy

### Unit Tests
- Test individual functions/methods
- Mock dependencies
- Fast execution

### Integration Tests
- Test with real database
- Use test containers
- `//go:build integration`

### E2E Tests
- Test full workflows
- Use Docker Compose
- Cypress for UI

```go
// Unit test example
func TestPricingEngine_Calculate(t *testing.T) {
    mockRepo := &MockRepository{}
    engine := NewPricingEngine(mockRepo)
    
    mockRepo.On("GetPriceList", mock.Anything).Return(priceList, nil)
    
    price, err := engine.Calculate(ctx, request)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, price)
}
```

---

## Debugging

### Local Debugging

```bash
# Run with Delve
dlv debug cmd/server/main.go
```

### Docker Debugging

```bash
# View logs
docker-compose logs -f pricing-engine

# Shell into container
docker-compose exec pricing-engine sh

# Check connectivity
docker-compose exec pricing-engine ping postgres
```

### Tracing

Access Jaeger UI at http://localhost:16686 to view traces.

---

## Configuration

All services use environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP port | 8080 |
| `DATABASE_URL` | PostgreSQL URL | - |
| `REDIS_URL` | Redis URL | - |
| `KAFKA_BROKERS` | Kafka brokers | - |
| `LOG_LEVEL` | Log level | info |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Jaeger endpoint | - |

---

## Commit Convention

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

Example:
```
feat(pricing): add volume discount calculation

Implements tiered pricing based on order quantity.
Closes #123
```
