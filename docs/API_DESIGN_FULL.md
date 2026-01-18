# OmniRoute Commerce Platform
# API Design & Integration Guide

---

## API Philosophy

OmniRoute follows **API-first design** principles:
- RESTful endpoints for CRUD operations
- GraphQL for complex queries (planned Phase 2)
- WebSockets for real-time updates
- Webhooks for event notifications

---

## Authentication & Authorization

### Authentication Methods

```yaml
# 1. API Key Authentication (Server-to-Server)
Authorization: Bearer sk_live_xxxxxxxxxxxxxxxxxxxxx

# 2. JWT Token (User Sessions)
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

# 3. OAuth 2.0 (Third-party Integrations)
Authorization: Bearer oauth_access_token_here
```

### Multi-Tenant Header

```yaml
# Required for all API calls
X-Tenant-ID: tnnt_xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### Rate Limits

| Plan | Requests/Min | Requests/Day |
|------|--------------|--------------|
| Free | 60 | 10,000 |
| Growth | 300 | 100,000 |
| Scale | 1,000 | 500,000 |
| Enterprise | Custom | Unlimited |

---

## Core API Endpoints

### Products API

```yaml
openapi: 3.0.3
info:
  title: OmniRoute Products API
  version: 1.0.0

paths:
  /api/v1/products:
    get:
      summary: List products
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
        - name: category_id
          in: query
          schema:
            type: string
            format: uuid
        - name: status
          in: query
          schema:
            type: string
            enum: [active, inactive, discontinued]
        - name: search
          in: query
          schema:
            type: string
        - name: sort
          in: query
          schema:
            type: string
            enum: [name, created_at, price, -name, -created_at, -price]
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Product'
                  meta:
                    $ref: '#/components/schemas/PaginationMeta'
    
    post:
      summary: Create product
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProductCreate'
      responses:
        '201':
          description: Product created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Product'

  /api/v1/products/{product_id}:
    get:
      summary: Get product by ID
      parameters:
        - name: product_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Product'
    
    put:
      summary: Update product
      parameters:
        - name: product_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProductUpdate'
      responses:
        '200':
          description: Product updated
    
    delete:
      summary: Delete product (soft delete)
      parameters:
        - name: product_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Product deleted

  /api/v1/products/{product_id}/variants:
    get:
      summary: List product variants
    post:
      summary: Create product variant

  /api/v1/products/{product_id}/prices:
    get:
      summary: Get all prices for product
      description: Returns prices across all price lists
    post:
      summary: Set product price in a price list

components:
  schemas:
    Product:
      type: object
      properties:
        id:
          type: string
          format: uuid
        tenant_id:
          type: string
          format: uuid
        sku:
          type: string
          maxLength: 100
        barcode:
          type: string
        name:
          type: string
          maxLength: 255
        description:
          type: string
        short_description:
          type: string
          maxLength: 500
        category_id:
          type: string
          format: uuid
        brand:
          type: string
        base_price:
          type: string
          format: decimal
          example: "1500.00"
        cost_price:
          type: string
          format: decimal
        currency:
          type: string
          default: "NGN"
        tax_category:
          type: string
        unit_of_measure:
          type: string
          default: "piece"
        units_per_case:
          type: integer
          default: 1
        min_order_quantity:
          type: integer
          default: 1
        max_order_quantity:
          type: integer
        order_multiple:
          type: integer
          default: 1
        weight_kg:
          type: number
        dimensions:
          type: object
          properties:
            length:
              type: number
            width:
              type: number
            height:
              type: number
            unit:
              type: string
        primary_image_url:
          type: string
          format: uri
        images:
          type: array
          items:
            type: object
            properties:
              url:
                type: string
              alt:
                type: string
              position:
                type: integer
        status:
          type: string
          enum: [active, inactive, discontinued, pending]
        visibility:
          type: string
          enum: [public, private, b2b_only, b2c_only]
        track_inventory:
          type: boolean
          default: true
        allow_backorder:
          type: boolean
          default: false
        tags:
          type: array
          items:
            type: string
        metadata:
          type: object
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    ProductCreate:
      type: object
      required:
        - sku
        - name
        - base_price
      properties:
        sku:
          type: string
        name:
          type: string
        base_price:
          type: string
        # ... other fields

    PaginationMeta:
      type: object
      properties:
        current_page:
          type: integer
        per_page:
          type: integer
        total_pages:
          type: integer
        total_count:
          type: integer
        has_next:
          type: boolean
        has_prev:
          type: boolean
```

### Orders API

```yaml
paths:
  /api/v1/orders:
    get:
      summary: List orders
      parameters:
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, confirmed, processing, delivered, cancelled]
        - name: customer_id
          in: query
          schema:
            type: string
            format: uuid
        - name: channel
          in: query
          schema:
            type: string
            enum: [web, mobile_app, whatsapp, ussd, sales_rep, api]
        - name: date_from
          in: query
          schema:
            type: string
            format: date
        - name: date_to
          in: query
          schema:
            type: string
            format: date
      responses:
        '200':
          description: List of orders
    
    post:
      summary: Create order
      description: |
        Create a new order. The pricing engine will automatically calculate:
        - Customer-specific prices
        - Volume discounts
        - Promotional discounts
        - Taxes
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/OrderCreate'
      responses:
        '201':
          description: Order created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Order'

  /api/v1/orders/{order_id}:
    get:
      summary: Get order details
    put:
      summary: Update order
    delete:
      summary: Cancel order

  /api/v1/orders/{order_id}/confirm:
    post:
      summary: Confirm order
      description: Move order from pending to confirmed status

  /api/v1/orders/{order_id}/fulfill:
    post:
      summary: Mark order as fulfilled
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                fulfillment_location_id:
                  type: string
                  format: uuid
                tracking_number:
                  type: string
                carrier:
                  type: string
                notes:
                  type: string

  /api/v1/orders/{order_id}/payments:
    get:
      summary: Get order payments
    post:
      summary: Record payment against order

components:
  schemas:
    OrderCreate:
      type: object
      required:
        - customer_id
        - items
      properties:
        customer_id:
          type: string
          format: uuid
        channel:
          type: string
          enum: [web, mobile_app, whatsapp, ussd, sales_rep, api]
          default: api
        items:
          type: array
          minItems: 1
          items:
            type: object
            required:
              - product_id
              - quantity
            properties:
              product_id:
                type: string
                format: uuid
              variant_id:
                type: string
                format: uuid
              quantity:
                type: integer
                minimum: 1
              notes:
                type: string
        shipping_address:
          $ref: '#/components/schemas/Address'
        billing_address:
          $ref: '#/components/schemas/Address'
        fulfillment_method:
          type: string
          enum: [delivery, pickup, van_delivery]
        requested_delivery_date:
          type: string
          format: date
        payment_method:
          type: string
          enum: [cash, bank_transfer, card, mobile_money, credit]
        coupon_codes:
          type: array
          items:
            type: string
        customer_notes:
          type: string
        metadata:
          type: object

    Order:
      type: object
      properties:
        id:
          type: string
          format: uuid
        order_number:
          type: string
          example: "ORD-2025-00001234"
        customer_id:
          type: string
          format: uuid
        customer:
          $ref: '#/components/schemas/CustomerSummary'
        channel:
          type: string
        status:
          type: string
        items:
          type: array
          items:
            $ref: '#/components/schemas/OrderItem'
        subtotal:
          type: string
          format: decimal
        discount_total:
          type: string
          format: decimal
        tax_total:
          type: string
          format: decimal
        shipping_total:
          type: string
          format: decimal
        grand_total:
          type: string
          format: decimal
        currency:
          type: string
        applied_promotions:
          type: array
          items:
            type: object
            properties:
              promotion_id:
                type: string
              name:
                type: string
              discount_amount:
                type: string
        payment_status:
          type: string
          enum: [pending, partial, paid, refunded]
        payment_method:
          type: string
        amount_paid:
          type: string
          format: decimal
        shipping_address:
          $ref: '#/components/schemas/Address'
        billing_address:
          $ref: '#/components/schemas/Address'
        fulfillment_method:
          type: string
        fulfillment_location_id:
          type: string
          format: uuid
        requested_delivery_date:
          type: string
          format: date
        promised_delivery_date:
          type: string
          format: date
        actual_delivery_date:
          type: string
          format: date
        placed_at:
          type: string
          format: date-time
        confirmed_at:
          type: string
          format: date-time
        shipped_at:
          type: string
          format: date-time
        delivered_at:
          type: string
          format: date-time
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time
```

### Customers API

```yaml
paths:
  /api/v1/customers:
    get:
      summary: List customers
      parameters:
        - name: type
          in: query
          schema:
            type: string
            enum: [consumer, retailer, wholesaler, distributor, enterprise]
        - name: tier
          in: query
          schema:
            type: string
            enum: [standard, silver, gold, platinum]
        - name: territory_id
          in: query
          schema:
            type: string
            format: uuid
        - name: assigned_rep_id
          in: query
          schema:
            type: string
            format: uuid
        - name: search
          in: query
          description: Search by name, email, phone, or business name
          schema:
            type: string
    post:
      summary: Create customer

  /api/v1/customers/{customer_id}:
    get:
      summary: Get customer details
    put:
      summary: Update customer
    delete:
      summary: Deactivate customer

  /api/v1/customers/{customer_id}/credit:
    get:
      summary: Get customer credit information
      responses:
        '200':
          content:
            application/json:
              schema:
                type: object
                properties:
                  credit_limit:
                    type: string
                    format: decimal
                  credit_used:
                    type: string
                    format: decimal
                  credit_available:
                    type: string
                    format: decimal
                  payment_terms:
                    type: integer
                    description: Days
                  credit_score:
                    type: integer
                  score_components:
                    type: object
                  last_review_date:
                    type: string
                    format: date
                  next_review_date:
                    type: string
                    format: date

  /api/v1/customers/{customer_id}/orders:
    get:
      summary: Get customer order history

  /api/v1/customers/{customer_id}/invoices:
    get:
      summary: Get customer invoices

  /api/v1/customers/{customer_id}/payments:
    get:
      summary: Get customer payment history
```

### Gig Workers API

```yaml
paths:
  /api/v1/gig-workers:
    get:
      summary: List gig workers
      parameters:
        - name: worker_type
          in: query
          schema:
            type: string
            enum: [delivery_rider, digital_sales_agent, field_auditor, collection_agent]
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, active, suspended, inactive]
        - name: is_online
          in: query
          schema:
            type: boolean
        - name: near_location
          in: query
          description: "Format: lat,lng"
          schema:
            type: string
        - name: radius_km
          in: query
          schema:
            type: number
            default: 10
        - name: min_level
          in: query
          schema:
            type: string
            enum: [starter, bronze, silver, gold, diamond, master]
    post:
      summary: Register new gig worker

  /api/v1/gig-workers/{worker_id}:
    get:
      summary: Get worker profile
    put:
      summary: Update worker profile

  /api/v1/gig-workers/{worker_id}/location:
    put:
      summary: Update worker location
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - latitude
                - longitude
              properties:
                latitude:
                  type: number
                longitude:
                  type: number
                accuracy:
                  type: number
                is_online:
                  type: boolean

  /api/v1/gig-workers/{worker_id}/tasks:
    get:
      summary: Get worker's tasks

  /api/v1/gig-workers/{worker_id}/earnings:
    get:
      summary: Get worker earnings
      parameters:
        - name: period
          in: query
          schema:
            type: string
            enum: [today, week, month, all]

  /api/v1/gig-tasks:
    get:
      summary: List available tasks
      parameters:
        - name: task_type
          in: query
          schema:
            type: string
        - name: status
          in: query
          schema:
            type: string
        - name: near_location
          in: query
          schema:
            type: string
        - name: scheduled_date
          in: query
          schema:
            type: string
            format: date
    post:
      summary: Create new task

  /api/v1/gig-tasks/{task_id}:
    get:
      summary: Get task details
    put:
      summary: Update task

  /api/v1/gig-tasks/{task_id}/assign:
    post:
      summary: Assign task to worker
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - worker_id
              properties:
                worker_id:
                  type: string
                  format: uuid
                assignment_method:
                  type: string
                  enum: [manual, claimed]

  /api/v1/gig-tasks/{task_id}/accept:
    post:
      summary: Worker accepts task

  /api/v1/gig-tasks/{task_id}/start:
    post:
      summary: Worker starts task

  /api/v1/gig-tasks/{task_id}/complete:
    post:
      summary: Complete task with proof
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                proof_photos:
                  type: array
                  items:
                    type: string
                    format: uri
                signature:
                  type: string
                  description: Base64 encoded signature image
                recipient_name:
                  type: string
                notes:
                  type: string
                collected_amount:
                  type: string
                  format: decimal
                payment_method:
                  type: string

  /api/v1/gig-tasks/{task_id}/rate:
    post:
      summary: Rate task completion
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - rating
              properties:
                rating:
                  type: integer
                  minimum: 1
                  maximum: 5
                feedback:
                  type: string
```

---

## Webhooks

### Webhook Events

```yaml
# Order Events
order.created          # New order placed
order.confirmed        # Order confirmed
order.processing       # Order being prepared
order.shipped          # Order shipped
order.delivered        # Order delivered
order.cancelled        # Order cancelled

# Payment Events
payment.initiated      # Payment started
payment.completed      # Payment successful
payment.failed         # Payment failed
payment.refunded       # Payment refunded

# Inventory Events
inventory.low_stock    # Stock below threshold
inventory.out_of_stock # Stock depleted
inventory.replenished  # Stock added

# Customer Events
customer.created       # New customer
customer.updated       # Customer profile updated
customer.credit_updated # Credit limit changed

# Gig Worker Events
gig.task_created       # New task available
gig.task_assigned      # Task assigned to worker
gig.task_completed     # Task completed
gig.task_failed        # Task failed

# Collection Events
collection.due_soon    # Payment due in 3 days
collection.overdue     # Payment past due
collection.collected   # Payment collected
```

### Webhook Payload Structure

```json
{
  "id": "evt_xxxxxxxxxxxxx",
  "type": "order.created",
  "created_at": "2025-01-18T10:30:00Z",
  "tenant_id": "tnnt_xxxxx",
  "data": {
    "order_id": "ord_xxxxx",
    "order_number": "ORD-2025-00001234",
    "customer_id": "cust_xxxxx",
    "grand_total": "150000.00",
    "currency": "NGN",
    "status": "pending"
  },
  "metadata": {
    "source": "mobile_app",
    "user_agent": "OmniRoute/1.0 iOS"
  }
}
```

### Webhook Security

```python
# Verify webhook signature (Python example)
import hmac
import hashlib

def verify_webhook_signature(payload: bytes, signature: str, secret: str) -> bool:
    """
    Verify that the webhook came from OmniRoute
    """
    expected = hmac.new(
        secret.encode('utf-8'),
        payload,
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(f"sha256={expected}", signature)

# Usage
@app.post("/webhooks/omniroute")
async def handle_webhook(request: Request):
    payload = await request.body()
    signature = request.headers.get("X-OmniRoute-Signature")
    
    if not verify_webhook_signature(payload, signature, WEBHOOK_SECRET):
        raise HTTPException(status_code=401, detail="Invalid signature")
    
    event = json.loads(payload)
    
    match event["type"]:
        case "order.created":
            await handle_new_order(event["data"])
        case "payment.completed":
            await handle_payment(event["data"])
        # ... handle other events
    
    return {"status": "ok"}
```

---

## SDK Examples

### Python SDK

```python
# Installation: pip install omniroute-sdk

from omniroute import OmniRouteClient

# Initialize client
client = OmniRouteClient(
    api_key="sk_live_xxxxxxxxxxxx",
    tenant_id="tnnt_xxxxxxxxxxxx"
)

# List products
products = client.products.list(
    status="active",
    category_id="cat_xxxxx",
    limit=50
)

for product in products:
    print(f"{product.sku}: {product.name} - {product.base_price}")

# Create order
order = client.orders.create(
    customer_id="cust_xxxxx",
    channel="api",
    items=[
        {"product_id": "prod_xxxxx", "quantity": 10},
        {"product_id": "prod_yyyyy", "quantity": 5},
    ],
    shipping_address={
        "address_line1": "123 Market Street",
        "city": "Lagos",
        "state": "Lagos",
        "country": "NGA"
    },
    fulfillment_method="delivery"
)

print(f"Order created: {order.order_number}")
print(f"Total: {order.currency} {order.grand_total}")

# Get customer credit
credit = client.customers.get_credit("cust_xxxxx")
print(f"Available credit: {credit.credit_available}")

# Calculate prices before ordering
prices = client.prices.calculate(
    customer_id="cust_xxxxx",
    items=[
        {"product_id": "prod_xxxxx", "quantity": 10},
        {"product_id": "prod_yyyyy", "quantity": 5},
    ]
)

for item in prices.items:
    print(f"{item.name}: {item.unit_price} x {item.quantity} = {item.line_total}")
    if item.discount_amount > 0:
        print(f"  Discount: -{item.discount_amount} ({item.price_source})")
```

### JavaScript/TypeScript SDK

```typescript
// Installation: npm install @omniroute/sdk

import { OmniRouteClient } from '@omniroute/sdk';

// Initialize client
const client = new OmniRouteClient({
  apiKey: 'sk_live_xxxxxxxxxxxx',
  tenantId: 'tnnt_xxxxxxxxxxxx',
});

// List products
const products = await client.products.list({
  status: 'active',
  categoryId: 'cat_xxxxx',
  limit: 50,
});

products.data.forEach(product => {
  console.log(`${product.sku}: ${product.name} - ${product.basePrice}`);
});

// Create order with async/await
async function createOrder() {
  try {
    const order = await client.orders.create({
      customerId: 'cust_xxxxx',
      channel: 'api',
      items: [
        { productId: 'prod_xxxxx', quantity: 10 },
        { productId: 'prod_yyyyy', quantity: 5 },
      ],
      shippingAddress: {
        addressLine1: '123 Market Street',
        city: 'Lagos',
        state: 'Lagos',
        country: 'NGA',
      },
      fulfillmentMethod: 'delivery',
    });

    console.log(`Order created: ${order.orderNumber}`);
    console.log(`Total: ${order.currency} ${order.grandTotal}`);
    
    return order;
  } catch (error) {
    if (error instanceof OmniRouteError) {
      console.error(`API Error: ${error.code} - ${error.message}`);
    }
    throw error;
  }
}

// Real-time updates with WebSocket
const ws = client.realtime.connect();

ws.subscribe('orders', (event) => {
  console.log(`Order ${event.data.orderNumber}: ${event.type}`);
});

ws.subscribe('inventory', (event) => {
  if (event.type === 'inventory.low_stock') {
    console.log(`Low stock alert: ${event.data.productId}`);
  }
});
```

### Go SDK

```go
// Installation: go get github.com/omniroute/omniroute-go

package main

import (
    "context"
    "fmt"
    "log"

    "github.com/omniroute/omniroute-go"
)

func main() {
    // Initialize client
    client := omniroute.NewClient(
        "sk_live_xxxxxxxxxxxx",
        omniroute.WithTenantID("tnnt_xxxxxxxxxxxx"),
    )

    ctx := context.Background()

    // List products
    products, err := client.Products.List(ctx, &omniroute.ProductListParams{
        Status: omniroute.String("active"),
        Limit:  omniroute.Int64(50),
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, product := range products.Data {
        fmt.Printf("%s: %s - %s\n", product.SKU, product.Name, product.BasePrice)
    }

    // Create order
    order, err := client.Orders.Create(ctx, &omniroute.OrderCreateParams{
        CustomerID: "cust_xxxxx",
        Channel:    "api",
        Items: []*omniroute.OrderItemParams{
            {ProductID: "prod_xxxxx", Quantity: 10},
            {ProductID: "prod_yyyyy", Quantity: 5},
        },
        ShippingAddress: &omniroute.AddressParams{
            AddressLine1: "123 Market Street",
            City:         "Lagos",
            State:        "Lagos",
            Country:      "NGA",
        },
        FulfillmentMethod: "delivery",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order created: %s\n", order.OrderNumber)
    fmt.Printf("Total: %s %s\n", order.Currency, order.GrandTotal)
}
```

---

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "invalid_request",
    "message": "The request was invalid or cannot be served.",
    "details": {
      "field": "items",
      "reason": "At least one item is required"
    },
    "request_id": "req_xxxxxxxxxxxxx",
    "documentation_url": "https://docs.omniroute.com/errors/invalid_request"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_request` | 400 | Request validation failed |
| `authentication_failed` | 401 | Invalid or missing API key |
| `permission_denied` | 403 | Not authorized for this action |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource state conflict |
| `rate_limited` | 429 | Too many requests |
| `internal_error` | 500 | Server error |
| `service_unavailable` | 503 | Service temporarily unavailable |

---

## Best Practices

### 1. Idempotency

```yaml
# Use idempotency keys for mutations
POST /api/v1/orders
X-Idempotency-Key: order_create_abc123_20250118
```

### 2. Pagination

```python
# Always use cursor-based pagination for large datasets
def fetch_all_orders(client):
    orders = []
    page = 1
    
    while True:
        response = client.orders.list(page=page, limit=100)
        orders.extend(response.data)
        
        if not response.meta.has_next:
            break
        page += 1
    
    return orders
```

### 3. Rate Limiting

```python
# Handle rate limits gracefully
import time
from omniroute.exceptions import RateLimitError

def safe_api_call(func, *args, max_retries=3, **kwargs):
    for attempt in range(max_retries):
        try:
            return func(*args, **kwargs)
        except RateLimitError as e:
            if attempt < max_retries - 1:
                wait_time = e.retry_after or (2 ** attempt)
                time.sleep(wait_time)
            else:
                raise
```

### 4. Webhook Reliability

```python
# Implement webhook retries with exponential backoff
# OmniRoute will retry failed webhooks:
# - Attempt 1: Immediate
# - Attempt 2: 1 minute
# - Attempt 3: 5 minutes
# - Attempt 4: 30 minutes
# - Attempt 5: 2 hours
# - Attempt 6: 24 hours

# Always return 2xx quickly, process async
@app.post("/webhooks/omniroute")
async def webhook(request: Request, background_tasks: BackgroundTasks):
    event = await request.json()
    background_tasks.add_task(process_webhook, event)
    return {"status": "received"}
```

---

*API Version: 1.0*
*Last Updated: January 2025*
*Documentation: https://docs.omniroute.com*
