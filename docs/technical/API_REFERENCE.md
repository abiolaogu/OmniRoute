# OmniRoute API Reference
## Complete API Documentation

---

## Introduction

The OmniRoute API provides programmatic access to all platform features. This reference documents all available endpoints, authentication methods, and best practices.

### Base URLs

| Environment | URL |
|-------------|-----|
| Production | `https://api.omniroute.io/v1` |
| Staging | `https://api-staging.omniroute.io/v1` |
| Sandbox | `https://api-sandbox.omniroute.io/v1` |

### API Standards
- **Protocol**: HTTPS only
- **Format**: JSON
- **Encoding**: UTF-8
- **Date Format**: ISO 8601 (e.g., `2026-01-18T20:30:00Z`)
- **Currency**: ISO 4217 codes (e.g., `NGN`, `USD`)

---

## Authentication

### API Key Authentication
```http
GET /api/v1/orders HTTP/1.1
Host: api.omniroute.io
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json
```

### OAuth 2.0
For applications requiring user-level access, use OAuth 2.0:

```http
POST /oauth/token HTTP/1.1
Host: api.omniroute.io
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=AUTHORIZATION_CODE&
client_id=YOUR_CLIENT_ID&
client_secret=YOUR_CLIENT_SECRET&
redirect_uri=YOUR_REDIRECT_URI
```

### Token Refresh
```http
POST /oauth/token HTTP/1.1
Host: api.omniroute.io
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&
refresh_token=YOUR_REFRESH_TOKEN&
client_id=YOUR_CLIENT_ID
```

---

## Rate Limiting

| Tier | Requests/Minute | Requests/Day |
|------|-----------------|--------------|
| Free | 60 | 1,000 |
| Standard | 300 | 10,000 |
| Enterprise | 1,000 | 100,000 |
| Custom | Unlimited | Custom |

**Rate Limit Headers:**
```
X-RateLimit-Limit: 300
X-RateLimit-Remaining: 299
X-RateLimit-Reset: 1705609800
```

---

## Common Patterns

### Pagination
```http
GET /api/v1/orders?page=2&per_page=50
```

**Response:**
```json
{
  "data": [...],
  "meta": {
    "current_page": 2,
    "per_page": 50,
    "total_pages": 10,
    "total_count": 500
  }
}
```

### Filtering
```http
GET /api/v1/orders?status=pending&created_after=2026-01-01
```

### Sorting
```http
GET /api/v1/orders?sort=-created_at,total_amount
```
(Prefix with `-` for descending order)

### Field Selection
```http
GET /api/v1/orders?fields=id,status,total_amount
```

---

## Core Endpoints

### Orders API

#### List Orders
```http
GET /api/v1/orders
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Filter by status |
| `customer_id` | uuid | Filter by customer |
| `created_after` | datetime | Created after date |
| `created_before` | datetime | Created before date |

**Response:**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "order_number": "ORD-2026-00001",
      "status": "pending",
      "customer_id": "...",
      "items": [...],
      "subtotal": 45000.00,
      "tax": 3375.00,
      "total": 48375.00,
      "currency": "NGN",
      "created_at": "2026-01-18T10:30:00Z"
    }
  ]
}
```

#### Create Order
```http
POST /api/v1/orders
Content-Type: application/json

{
  "customer_id": "...",
  "shipping_address_id": "...",
  "items": [
    {
      "product_id": "...",
      "quantity": 10,
      "unit_price": 4500.00
    }
  ],
  "payment_method": "credit",
  "notes": "Deliver before noon"
}
```

#### Get Order
```http
GET /api/v1/orders/{order_id}
```

#### Update Order
```http
PATCH /api/v1/orders/{order_id}
```

#### Cancel Order
```http
POST /api/v1/orders/{order_id}/cancel
```

---

### Products API

#### List Products
```http
GET /api/v1/products
```

#### Get Product
```http
GET /api/v1/products/{product_id}
```

#### Product Inventory
```http
GET /api/v1/products/{product_id}/inventory
```

---

### Customers API

#### List Customers
```http
GET /api/v1/customers
```

#### Create Customer
```http
POST /api/v1/customers
```

#### Get Customer Credit
```http
GET /api/v1/customers/{customer_id}/credit
```

---

### Payments API

#### Initiate Payment
```http
POST /api/v1/payments
```

**Request:**
```json
{
  "order_id": "...",
  "amount": 48375.00,
  "currency": "NGN",
  "payment_method": "bank_transfer",
  "callback_url": "https://your-app.com/payment-callback"
}
```

#### Verify Payment
```http
GET /api/v1/payments/{payment_id}/verify
```

---

### Inventory API

#### Get Inventory Levels
```http
GET /api/v1/inventory
```

#### Adjust Inventory
```http
POST /api/v1/inventory/adjustments
```

#### Transfer Inventory
```http
POST /api/v1/inventory/transfers
```

---

### Shipments API

#### Create Shipment
```http
POST /api/v1/shipments
```

#### Track Shipment
```http
GET /api/v1/shipments/{shipment_id}/tracking
```

---

## GraphQL API

In addition to REST, OmniRoute provides a GraphQL API via Hasura.

**Endpoint:** `https://api.omniroute.io/v1/graphql`

### Example Query
```graphql
query GetOrderWithItems($orderId: uuid!) {
  orders_by_pk(id: $orderId) {
    id
    order_number
    status
    total
    customer {
      name
      email
    }
    items {
      product {
        name
        sku
      }
      quantity
      unit_price
    }
  }
}
```

### Example Mutation
```graphql
mutation CreateOrder($input: orders_insert_input!) {
  insert_orders_one(object: $input) {
    id
    order_number
  }
}
```

### Subscriptions
```graphql
subscription OrderUpdates($customerId: uuid!) {
  orders(where: {customer_id: {_eq: $customerId}}) {
    id
    status
    updated_at
  }
}
```

---

## Error Handling

### Error Response Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "quantity",
        "message": "must be greater than 0"
      }
    ],
    "request_id": "req_abc123"
  }
}
```

### Error Codes
| Code | HTTP Status | Description |
|------|-------------|-------------|
| `AUTHENTICATION_ERROR` | 401 | Invalid or missing credentials |
| `AUTHORIZATION_ERROR` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `RATE_LIMIT_ERROR` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Webhooks

### Webhook Events
| Event | Description |
|-------|-------------|
| `order.created` | New order placed |
| `order.updated` | Order status changed |
| `order.cancelled` | Order cancelled |
| `payment.completed` | Payment successful |
| `payment.failed` | Payment failed |
| `shipment.created` | Shipment created |
| `shipment.delivered` | Shipment delivered |
| `inventory.low` | Stock below threshold |

### Webhook Payload
```json
{
  "id": "evt_abc123",
  "type": "order.created",
  "created_at": "2026-01-18T10:30:00Z",
  "data": {
    "order_id": "...",
    "order_number": "ORD-2026-00001",
    "total": 48375.00
  }
}
```

### Webhook Verification
Verify webhook signatures using HMAC-SHA256:
```
X-OmniRoute-Signature: sha256=abc123...
```

---

## SDKs

### Node.js
```bash
npm install @omniroute/sdk
```

```javascript
const OmniRoute = require('@omniroute/sdk');

const client = new OmniRoute.Client({
  apiKey: 'YOUR_API_KEY',
  environment: 'production'
});

const orders = await client.orders.list({ status: 'pending' });
```

### Python
```bash
pip install omniroute
```

```python
from omniroute import OmniRouteClient

client = OmniRouteClient(api_key='YOUR_API_KEY')
orders = client.orders.list(status='pending')
```

### Go
```bash
go get github.com/omniroute/omniroute-go
```

```go
import "github.com/omniroute/omniroute-go"

client := omniroute.NewClient("YOUR_API_KEY")
orders, err := client.Orders.List(ctx, &omniroute.OrderListParams{
    Status: "pending",
})
```

---

## Best Practices

1. **Use Idempotency Keys** for POST requests
2. **Implement Retry Logic** with exponential backoff
3. **Cache Responses** where appropriate
4. **Use Webhooks** instead of polling
5. **Handle Rate Limits** gracefully
6. **Validate Webhook Signatures** always
7. **Use Pagination** for large datasets
8. **Log Request IDs** for debugging
