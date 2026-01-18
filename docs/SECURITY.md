# OmniRoute - Security Documentation

## Security Overview

OmniRoute implements defense-in-depth security across all layers of the platform.

---

## Authentication & Authorization

### JWT (JSON Web Tokens)

```json
{
  "sub": "user-uuid",
  "tenant_id": "tenant-uuid",
  "roles": ["retailer", "order_manager"],
  "permissions": ["orders:read", "orders:write", "products:read"],
  "iat": 1705612800,
  "exp": 1705699200
}
```

### Role-Based Access Control (RBAC)

| Role | Permissions |
|------|-------------|
| **admin** | Full system access |
| **tenant_admin** | Full tenant access |
| **retailer** | Orders, products, payments |
| **worker** | Tasks, earnings, proofs |
| **driver** | Routes, deliveries |
| **readonly** | View only |

### Hasura Permission Example

```yaml
# products table permissions
select_permissions:
  role: retailer
  permission:
    filter:
      tenant_id:
        _eq: X-Hasura-Tenant-Id
    columns: [id, name, sku, price, description]

insert_permissions:
  role: retailer
  permission:
    check:
      tenant_id:
        _eq: X-Hasura-Tenant-Id
    columns: [name, sku, price, description]
```

---

## Data Protection

### Encryption

| Layer | Method |
|-------|--------|
| **In Transit** | TLS 1.3 |
| **At Rest** | AES-256 (YugabyteDB) |
| **Application** | Argon2id (passwords) |
| **API Keys** | SHA-256 hashing |

### PII Handling

```go
// PII fields are encrypted before storage
type Customer struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    Email     EncryptedString `json:"email"`      // Encrypted
    Phone     EncryptedString `json:"phone"`      // Encrypted
    NIN       EncryptedString `json:"nin"`        // Encrypted (National ID)
    Name      string          `json:"name"`
    CreatedAt time.Time
}
```

---

## API Security

### Rate Limiting

```go
// Rate limit configuration
type RateLimitConfig struct {
    RequestsPerSecond int           // 100
    BurstSize         int           // 200
    WindowDuration    time.Duration // 1 minute
}

// Per-tenant limits
var tenantLimits = map[TenantTier]RateLimitConfig{
    TierFree:       {100, 200, time.Minute},
    TierStandard:   {500, 1000, time.Minute},
    TierEnterprise: {5000, 10000, time.Minute},
}
```

### Input Validation

```go
func (r *CreateOrderRequest) Validate() error {
    return validation.ValidateStruct(r,
        validation.Field(&r.CustomerID, validation.Required, is.UUID),
        validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
        validation.Field(&r.DeliveryAddress, validation.Required),
        validation.Field(&r.TotalAmount, validation.Required, validation.Min(0)),
    )
}
```

### SQL Injection Prevention

```go
// Parameterized queries ONLY
func (r *Repository) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
    query := `SELECT * FROM orders WHERE id = $1 AND tenant_id = $2`
    return r.db.Query(ctx, query, id, getTenantID(ctx))
}
```

---

## Infrastructure Security

### Network Policies

```yaml
# Allow only specific traffic
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: pricing-engine-policy
spec:
  podSelector:
    matchLabels:
      app: pricing-engine
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: hasura
    ports:
    - port: 8081
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - port: 5432
```

### Secret Management

```yaml
# External Secrets Operator
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: omniroute-secrets
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: ClusterSecretStore
  target:
    name: omniroute-secrets
  data:
  - secretKey: DATABASE_URL
    remoteRef:
      key: omniroute/database
      property: url
```

---

## Audit Logging

### Audit Event Structure

```go
type AuditEvent struct {
    ID          uuid.UUID              `json:"id"`
    TenantID    uuid.UUID              `json:"tenant_id"`
    UserID      uuid.UUID              `json:"user_id"`
    Action      string                 `json:"action"`      // CREATE, UPDATE, DELETE
    Resource    string                 `json:"resource"`    // orders, payments
    ResourceID  string                 `json:"resource_id"`
    IPAddress   string                 `json:"ip_address"`
    UserAgent   string                 `json:"user_agent"`
    Changes     map[string]interface{} `json:"changes"`
    Timestamp   time.Time              `json:"timestamp"`
}
```

### Logged Events

| Category | Events |
|----------|--------|
| Authentication | Login, Logout, Failed Login, Password Change |
| Authorization | Permission Denied, Role Change |
| Data Access | Read Sensitive Data, Export Data |
| Data Modification | Create, Update, Delete |
| Administrative | User Management, Config Change |
| Payment | Transaction, Refund, Settlement |

---

## Compliance

### GDPR

- [x] Data Subject Rights (Access, Rectification, Erasure)
- [x] Consent Management
- [x] Data Processing Agreements
- [x] Privacy Impact Assessments
- [x] 72-hour Breach Notification

### NDPR (Nigeria)

- [x] Local data residency options
- [x] Consent-based processing
- [x] Audit trail maintenance

### PCI-DSS

- [x] No card data storage
- [x] Tokenization via payment processors
- [x] TLS everywhere
- [x] Access controls

---

## Vulnerability Management

### Dependency Scanning

```yaml
# GitHub Actions
- name: Run Trivy vulnerability scanner
  uses: aquasecurity/trivy-action@master
  with:
    scan-type: 'fs'
    ignore-unfixed: true
    severity: 'CRITICAL,HIGH'
```

### Container Scanning

```yaml
# Trivy container scan
- name: Scan Docker image
  run: |
    trivy image --severity HIGH,CRITICAL \
      registry.omniroute.io/pricing-engine:${{ github.sha }}
```

### SAST (Static Analysis)

```yaml
# GoSec for Go code
- name: Run GoSec
  run: gosec -fmt=sarif -out=results.sarif ./...
```

---

## Incident Response

### Severity Levels

| Level | Response Time | Examples |
|-------|---------------|----------|
| **P1 Critical** | 15 min | Data breach, full outage |
| **P2 High** | 1 hour | Payment failure, partial outage |
| **P3 Medium** | 4 hours | Performance degradation |
| **P4 Low** | 24 hours | Minor bugs, cosmetic issues |

### Response Procedure

1. **Detect**: Automated alerts via monitoring
2. **Triage**: Assess severity and impact
3. **Contain**: Isolate affected systems
4. **Eradicate**: Remove threat/fix issue
5. **Recover**: Restore normal operations
6. **Review**: Post-incident analysis

---

## Security Checklist

### Pre-Deployment

- [ ] All secrets in secret manager
- [ ] TLS certificates valid
- [ ] Network policies applied
- [ ] Vulnerability scan passed
- [ ] Penetration test completed

### Ongoing

- [ ] Weekly dependency updates
- [ ] Monthly access reviews
- [ ] Quarterly security assessments
- [ ] Annual penetration testing
