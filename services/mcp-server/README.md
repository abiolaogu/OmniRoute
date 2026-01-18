# OmniRoute MCP Server

## Overview

The MCP (Model Context Protocol) server provides a standardized interface for AI assistants to interact with the OmniRoute Commerce Platform. It exposes platform capabilities as tools, resources, and prompts that AI models can use to help users with commerce operations.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        AI Assistants                             │
│               (Claude, GPT, Gemini, Custom)                      │
└─────────────────────────┬───────────────────────────────────────┘
                          │ MCP Protocol (JSON-RPC)
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      MCP Server (:8200)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │    Tools    │  │  Resources  │  │   Prompts   │              │
│  │  (15 tools) │  │(10 sources) │  │(4 templates)│              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
└─────────┼────────────────┼────────────────┼─────────────────────┘
          │                │                │
          ▼                ▼                ▼
    ┌─────────────────────────────────────────────┐
    │            Integration Layer                 │
    │  ┌─────────┐  ┌─────────┐  ┌─────────────┐  │
    │  │ Hasura  │  │   n8n   │  │  Temporal   │  │
    │  │ GraphQL │  │Workflows│  │  Workflows  │  │
    │  └────┬────┘  └────┬────┘  └──────┬──────┘  │
    └───────┼────────────┼──────────────┼─────────┘
            │            │              │
            ▼            ▼              ▼
    ┌───────────────────────────────────────────────┐
    │              OmniRoute Services               │
    │  Pricing | Orders | Inventory | Customers     │
    │  Fleet | Analytics | AI/ML Services           │
    └───────────────────────────────────────────────┘
```

## Available Tools (15)

| Tool | Description | Integration |
|------|-------------|-------------|
| `create_order` | Create new orders | Hasura → order-service |
| `get_order` | Get order details | Hasura |
| `cancel_order` | Cancel orders | Hasura → Temporal |
| `search_products` | Search product catalog | Hasura |
| `get_product` | Get product details | Hasura |
| `check_inventory` | Check stock levels | Hasura |
| `get_customer` | Get customer profile | Hasura |
| `get_customer_credit` | Get credit info | Hasura → credit-scoring |
| `calculate_price` | Dynamic pricing | pricing-engine |
| `get_sales_report` | Sales analytics | analytics-service |
| `get_demand_forecast` | AI forecast | forecasting-service |
| `trigger_workflow` | Trigger n8n workflow | n8n |
| `start_order_workflow` | Start order processing | Temporal |
| `start_credit_review` | Credit review workflow | Temporal |
| `track_vehicle` | Vehicle tracking | fleet-service |
| `assign_worker` | Assign gig worker | gig-platform |
| `get_recommendations` | AI recommendations | recommendations-service |

## Available Resources (10)

| URI | Description |
|-----|-------------|
| `omniroute://catalog/categories` | Product categories |
| `omniroute://catalog/products` | Product catalog |
| `omniroute://orders/recent` | Recent orders |
| `omniroute://customers/segments` | Customer segments |
| `omniroute://inventory/low-stock` | Low stock items |
| `omniroute://analytics/dashboard` | Dashboard metrics |
| `omniroute://workers/available` | Available workers |
| `omniroute://fleet/active` | Active vehicles |
| `omniroute://config/pricing-rules` | Pricing rules |
| `omniroute://config/workflows` | Workflow definitions |

## Available Prompts (4)

| Prompt | Description |
|--------|-------------|
| `order_assistant` | Help customers with orders |
| `inventory_analyst` | Analyze inventory and recommend restocking |
| `sales_reporter` | Generate sales insights |
| `credit_advisor` | Evaluate credit applications |

## API Endpoints

### MCP Protocol (JSON-RPC)

```bash
# Initialize connection
POST /mcp/initialize

# List available tools
POST /mcp/tools/list

# Call a tool
POST /mcp/tools/call
{
  "name": "create_order",
  "arguments": {
    "customer_id": "uuid",
    "items": [{"product_id": "uuid", "quantity": 10}]
  }
}

# List resources
POST /mcp/resources/list

# Read resource
POST /mcp/resources/read
{
  "uri": "omniroute://analytics/dashboard"
}

# Get prompt
POST /mcp/prompts/get
{
  "name": "order_assistant",
  "arguments": {"customer_id": "uuid"}
}
```

### SSE Streaming

```bash
# Connect to SSE stream for real-time updates
GET /mcp/sse
```

### REST API

```bash
# Proxy GraphQL queries to Hasura
POST /api/v1/query

# Trigger workflow (auto-routes to n8n or Temporal)
POST /api/v1/workflow/trigger
```

## Integration Details

### Hasura GraphQL

The MCP server uses Hasura as the primary data layer:
- All read operations go through Hasura GraphQL
- Mutations are executed via Hasura with proper permissions
- Real-time subscriptions for resource updates

### n8n Workflows

Simple, event-driven automations:
- Low stock alerts
- Customer onboarding flows
- Payment reconciliation
- Scheduled reports

### Temporal Workflows

Complex, long-running business processes:
- Order processing (multi-step with compensations)
- Credit review (with manual approval)
- Worker payouts (batch processing)

## Configuration

```yaml
# config/config.yaml
hasura:
  endpoint: http://hasura:8080/v1/graphql
  admin_secret: ${HASURA_ADMIN_SECRET}

n8n:
  endpoint: http://n8n:5678
  api_key: ${N8N_API_KEY}

temporal:
  host: temporal:7233
  namespace: omniroute
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | 8200 |
| `HASURA_ENDPOINT` | Hasura GraphQL URL | http://hasura:8080/v1/graphql |
| `HASURA_ADMIN_SECRET` | Hasura admin secret | - |
| `N8N_ENDPOINT` | n8n API URL | http://n8n:5678 |
| `N8N_API_KEY` | n8n API key | - |
| `TEMPORAL_HOST` | Temporal server | temporal:7233 |
| `TEMPORAL_NAMESPACE` | Temporal namespace | omniroute |

## Running Locally

```bash
cd services/mcp-server
go run cmd/server/main.go
```

## Docker

```bash
docker build -t omniroute/mcp-server .
docker run -p 8200:8200 omniroute/mcp-server
```

## Example Usage

### From Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "omniroute": {
      "command": "curl",
      "args": ["-N", "http://localhost:8200/mcp/sse"]
    }
  }
}
```

### From Custom Application

```python
import httpx

# Initialize
response = httpx.post("http://localhost:8200/mcp/initialize")
print(response.json())

# List tools
response = httpx.post("http://localhost:8200/mcp/tools/list")
tools = response.json()["tools"]

# Call a tool
response = httpx.post(
    "http://localhost:8200/mcp/tools/call",
    json={
        "name": "search_products",
        "arguments": {"query": "rice", "limit": 10}
    }
)
print(response.json())
```
