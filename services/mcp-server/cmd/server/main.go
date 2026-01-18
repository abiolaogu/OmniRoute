// Package main provides the MCP (Model Context Protocol) server for OmniRoute.
// This server exposes platform capabilities as tools and resources for AI assistants.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// MCP Protocol Version
const MCPVersion = "2024-11-05"

// ServerInfo provides server metadata
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Capabilities describes what the MCP server supports
type Capabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// Prompt represents an MCP prompt template
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// MCPServer holds the server state
type MCPServer struct {
	hasuraClient   *HasuraClient
	n8nClient      *N8NClient
	temporalClient *TemporalClient
	tools          []Tool
	resources      []Resource
	prompts        []Prompt
}

// Hasura Client
type HasuraClient struct {
	endpoint string
	secret   string
	client   *http.Client
}

func NewHasuraClient(endpoint, secret string) *HasuraClient {
	return &HasuraClient{
		endpoint: endpoint,
		secret:   secret,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (h *HasuraClient) Query(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", h.endpoint, io.NopCloser(nil))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", h.secret)
	_ = body // Would be used in actual request

	// Placeholder - actual implementation would make HTTP call
	return map[string]interface{}{"data": nil}, nil
}

// n8n Client
type N8NClient struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

func NewN8NClient(endpoint, apiKey string) *N8NClient {
	return &N8NClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (n *N8NClient) TriggerWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) error {
	// Trigger n8n workflow via webhook
	return nil
}

func (n *N8NClient) GetWorkflowStatus(ctx context.Context, executionID string) (map[string]interface{}, error) {
	return map[string]interface{}{"status": "completed"}, nil
}

// Temporal Client
type TemporalClient struct {
	hostPort  string
	namespace string
}

func NewTemporalClient(hostPort, namespace string) *TemporalClient {
	return &TemporalClient{
		hostPort:  hostPort,
		namespace: namespace,
	}
}

func (t *TemporalClient) StartWorkflow(ctx context.Context, workflowType string, input interface{}) (string, error) {
	// Start Temporal workflow
	return "workflow-run-id", nil
}

func (t *TemporalClient) GetWorkflowResult(ctx context.Context, workflowID, runID string) (interface{}, error) {
	return map[string]interface{}{"result": "success"}, nil
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	server := &MCPServer{
		hasuraClient:   NewHasuraClient(getEnv("HASURA_ENDPOINT", "http://hasura:8080/v1/graphql"), getEnv("HASURA_ADMIN_SECRET", "")),
		n8nClient:      NewN8NClient(getEnv("N8N_ENDPOINT", "http://n8n:5678"), getEnv("N8N_API_KEY", "")),
		temporalClient: NewTemporalClient(getEnv("TEMPORAL_HOST", "temporal:7233"), getEnv("TEMPORAL_NAMESPACE", "default")),
	}
	server.registerTools()
	server.registerResources()
	server.registerPrompts()
	return server
}

func (s *MCPServer) registerTools() {
	s.tools = []Tool{
		// Order Management Tools
		{
			Name:        "create_order",
			Description: "Create a new order in the OmniRoute platform",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_id":      map[string]interface{}{"type": "string", "description": "Customer UUID"},
					"items":            map[string]interface{}{"type": "array", "description": "Order items"},
					"shipping_address": map[string]interface{}{"type": "object", "description": "Delivery address"},
				},
				"required": []string{"customer_id", "items"},
			},
		},
		{
			Name:        "get_order",
			Description: "Get order details by ID or order number",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id":     map[string]interface{}{"type": "string", "description": "Order UUID"},
					"order_number": map[string]interface{}{"type": "string", "description": "Order number"},
				},
			},
		},
		{
			Name:        "cancel_order",
			Description: "Cancel an existing order",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id": map[string]interface{}{"type": "string", "description": "Order UUID"},
					"reason":   map[string]interface{}{"type": "string", "description": "Cancellation reason"},
				},
				"required": []string{"order_id"},
			},
		},

		// Product Tools
		{
			Name:        "search_products",
			Description: "Search for products in the catalog",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":       map[string]interface{}{"type": "string", "description": "Search query"},
					"category_id": map[string]interface{}{"type": "string", "description": "Filter by category"},
					"min_price":   map[string]interface{}{"type": "number", "description": "Minimum price"},
					"max_price":   map[string]interface{}{"type": "number", "description": "Maximum price"},
					"limit":       map[string]interface{}{"type": "integer", "description": "Result limit", "default": 20},
				},
			},
		},
		{
			Name:        "get_product",
			Description: "Get detailed product information",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "Product UUID"},
					"sku":        map[string]interface{}{"type": "string", "description": "Product SKU"},
				},
			},
		},
		{
			Name:        "check_inventory",
			Description: "Check product inventory/stock levels",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id":   map[string]interface{}{"type": "string", "description": "Product UUID"},
					"warehouse_id": map[string]interface{}{"type": "string", "description": "Warehouse UUID"},
				},
				"required": []string{"product_id"},
			},
		},

		// Customer Tools
		{
			Name:        "get_customer",
			Description: "Get customer profile and history",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_id": map[string]interface{}{"type": "string", "description": "Customer UUID"},
					"email":       map[string]interface{}{"type": "string", "description": "Customer email"},
					"phone":       map[string]interface{}{"type": "string", "description": "Customer phone"},
				},
			},
		},
		{
			Name:        "get_customer_credit",
			Description: "Get customer credit limit and usage",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_id": map[string]interface{}{"type": "string", "description": "Customer UUID"},
				},
				"required": []string{"customer_id"},
			},
		},

		// Pricing Tools
		{
			Name:        "calculate_price",
			Description: "Calculate dynamic price for products",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id":  map[string]interface{}{"type": "string", "description": "Product UUID"},
					"quantity":    map[string]interface{}{"type": "integer", "description": "Quantity"},
					"customer_id": map[string]interface{}{"type": "string", "description": "Customer UUID for personalized pricing"},
				},
				"required": []string{"product_id", "quantity"},
			},
		},

		// Analytics Tools
		{
			Name:        "get_sales_report",
			Description: "Get sales analytics and reports",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"period":      map[string]interface{}{"type": "string", "enum": []string{"today", "week", "month", "quarter", "year"}},
					"group_by":    map[string]interface{}{"type": "string", "enum": []string{"day", "week", "month", "product", "customer"}},
					"category_id": map[string]interface{}{"type": "string", "description": "Filter by category"},
				},
			},
		},
		{
			Name:        "get_demand_forecast",
			Description: "Get AI-powered demand forecast for products",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id":   map[string]interface{}{"type": "string", "description": "Product UUID"},
					"horizon_days": map[string]interface{}{"type": "integer", "description": "Forecast horizon in days", "default": 30},
				},
				"required": []string{"product_id"},
			},
		},

		// Workflow Tools (n8n integration)
		{
			Name:        "trigger_workflow",
			Description: "Trigger an automated workflow in n8n",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"workflow_name": map[string]interface{}{"type": "string", "enum": []string{"order_processing", "low_stock_alert", "payment_reconciliation", "customer_onboarding", "worker_payout"}},
					"payload":       map[string]interface{}{"type": "object", "description": "Workflow input data"},
				},
				"required": []string{"workflow_name"},
			},
		},

		// Temporal Workflow Tools
		{
			Name:        "start_order_workflow",
			Description: "Start a Temporal order processing workflow",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id":    map[string]interface{}{"type": "string", "description": "Order UUID"},
					"customer_id": map[string]interface{}{"type": "string", "description": "Customer UUID"},
				},
				"required": []string{"order_id", "customer_id"},
			},
		},
		{
			Name:        "start_credit_review",
			Description: "Start a credit review workflow for a customer",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_id":      map[string]interface{}{"type": "string", "description": "Customer UUID"},
					"requested_amount": map[string]interface{}{"type": "number", "description": "Requested credit amount"},
				},
				"required": []string{"customer_id", "requested_amount"},
			},
		},

		// Fleet Tools
		{
			Name:        "track_vehicle",
			Description: "Get real-time vehicle location and status",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"vehicle_id":          map[string]interface{}{"type": "string", "description": "Vehicle UUID"},
					"registration_number": map[string]interface{}{"type": "string", "description": "Vehicle registration"},
				},
			},
		},

		// Worker Tools
		{
			Name:        "assign_worker",
			Description: "Assign a gig worker to an order",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id":  map[string]interface{}{"type": "string", "description": "Order UUID"},
					"worker_id": map[string]interface{}{"type": "string", "description": "Specific worker UUID (optional)"},
				},
				"required": []string{"order_id"},
			},
		},

		// Recommendations
		{
			Name:        "get_recommendations",
			Description: "Get AI product recommendations for a customer",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"customer_id": map[string]interface{}{"type": "string", "description": "Customer UUID"},
					"limit":       map[string]interface{}{"type": "integer", "description": "Number of recommendations", "default": 10},
				},
				"required": []string{"customer_id"},
			},
		},
	}
}

func (s *MCPServer) registerResources() {
	s.resources = []Resource{
		// Data Resources
		{URI: "omniroute://catalog/categories", Name: "Product Categories", Description: "List of all product categories", MimeType: "application/json"},
		{URI: "omniroute://catalog/products", Name: "Products", Description: "Product catalog", MimeType: "application/json"},
		{URI: "omniroute://orders/recent", Name: "Recent Orders", Description: "Last 100 orders", MimeType: "application/json"},
		{URI: "omniroute://customers/segments", Name: "Customer Segments", Description: "Customer segmentation data", MimeType: "application/json"},
		{URI: "omniroute://inventory/low-stock", Name: "Low Stock Items", Description: "Products below reorder point", MimeType: "application/json"},
		{URI: "omniroute://analytics/dashboard", Name: "Dashboard Metrics", Description: "Key business metrics", MimeType: "application/json"},
		{URI: "omniroute://workers/available", Name: "Available Workers", Description: "Currently available gig workers", MimeType: "application/json"},
		{URI: "omniroute://fleet/active", Name: "Active Vehicles", Description: "Vehicles currently in operation", MimeType: "application/json"},

		// Configuration Resources
		{URI: "omniroute://config/pricing-rules", Name: "Pricing Rules", Description: "Active pricing rules and discounts", MimeType: "application/json"},
		{URI: "omniroute://config/workflows", Name: "Workflow Definitions", Description: "Available n8n workflows", MimeType: "application/json"},
	}
}

func (s *MCPServer) registerPrompts() {
	s.prompts = []Prompt{
		{
			Name:        "order_assistant",
			Description: "Help customers place and track orders",
			Arguments: []PromptArgument{
				{Name: "customer_id", Description: "Customer UUID", Required: true},
				{Name: "context", Description: "Conversation context", Required: false},
			},
		},
		{
			Name:        "inventory_analyst",
			Description: "Analyze inventory levels and suggest restocking",
			Arguments: []PromptArgument{
				{Name: "warehouse_id", Description: "Warehouse to analyze", Required: false},
				{Name: "category_id", Description: "Category to focus on", Required: false},
			},
		},
		{
			Name:        "sales_reporter",
			Description: "Generate sales insights and recommendations",
			Arguments: []PromptArgument{
				{Name: "period", Description: "Time period for analysis", Required: true},
				{Name: "focus_area", Description: "Specific area to analyze", Required: false},
			},
		},
		{
			Name:        "credit_advisor",
			Description: "Evaluate credit applications and provide recommendations",
			Arguments: []PromptArgument{
				{Name: "customer_id", Description: "Customer UUID", Required: true},
			},
		},
	}
}

func main() {
	port := getEnv("SERVER_PORT", "8200")
	mcpServer := NewMCPServer()

	router := gin.Default()

	// Health endpoints
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// MCP Protocol Endpoints (JSON-RPC over HTTP)
	mcp := router.Group("/mcp")
	{
		// Initialize handshake
		mcp.POST("/initialize", mcpServer.handleInitialize)

		// Tools
		mcp.POST("/tools/list", mcpServer.handleListTools)
		mcp.POST("/tools/call", mcpServer.handleCallTool)

		// Resources
		mcp.POST("/resources/list", mcpServer.handleListResources)
		mcp.POST("/resources/read", mcpServer.handleReadResource)

		// Prompts
		mcp.POST("/prompts/list", mcpServer.handleListPrompts)
		mcp.POST("/prompts/get", mcpServer.handleGetPrompt)
	}

	// SSE endpoint for streaming (MCP transport)
	router.GET("/mcp/sse", mcpServer.handleSSE)

	// REST API for direct access
	api := router.Group("/api/v1")
	{
		api.POST("/query", mcpServer.handleGraphQLProxy)
		api.POST("/workflow/trigger", mcpServer.handleWorkflowTrigger)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("MCP Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// MCP Protocol Handlers

func (s *MCPServer) handleInitialize(c *gin.Context) {
	response := map[string]interface{}{
		"protocolVersion": MCPVersion,
		"serverInfo": ServerInfo{
			Name:    "OmniRoute MCP Server",
			Version: "1.0.0",
		},
		"capabilities": Capabilities{
			Tools:     &ToolsCapability{ListChanged: true},
			Resources: &ResourcesCapability{Subscribe: true, ListChanged: true},
			Prompts:   &PromptsCapability{ListChanged: true},
		},
	}
	c.JSON(200, response)
}

func (s *MCPServer) handleListTools(c *gin.Context) {
	c.JSON(200, gin.H{"tools": s.tools})
}

func (s *MCPServer) handleCallTool(c *gin.Context) {
	var req struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := s.executeTool(c.Request.Context(), req.Name, req.Arguments)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"content": []map[string]interface{}{
			{"type": "text", "text": result},
		},
	})
}

func (s *MCPServer) executeTool(ctx context.Context, name string, args map[string]interface{}) (string, error) {
	switch name {
	case "create_order":
		return s.executeCreateOrder(ctx, args)
	case "get_order":
		return s.executeGetOrder(ctx, args)
	case "search_products":
		return s.executeSearchProducts(ctx, args)
	case "check_inventory":
		return s.executeCheckInventory(ctx, args)
	case "calculate_price":
		return s.executeCalculatePrice(ctx, args)
	case "trigger_workflow":
		return s.executeTriggerWorkflow(ctx, args)
	case "start_order_workflow":
		return s.executeStartOrderWorkflow(ctx, args)
	case "get_recommendations":
		return s.executeGetRecommendations(ctx, args)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *MCPServer) executeCreateOrder(ctx context.Context, args map[string]interface{}) (string, error) {
	// Execute via Hasura GraphQL
	mutation := `
		mutation CreateOrder($customer_id: uuid!, $items: [order_items_insert_input!]!) {
			insert_orders_one(object: {customer_id: $customer_id, items: {data: $items}}) {
				id
				order_number
				status
			}
		}
	`
	result, err := s.hasuraClient.Query(ctx, mutation, args)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) executeGetOrder(ctx context.Context, args map[string]interface{}) (string, error) {
	query := `
		query GetOrder($order_id: uuid, $order_number: String) {
			orders(where: {_or: [{id: {_eq: $order_id}}, {order_number: {_eq: $order_number}}]}) {
				id
				order_number
				status
				total_amount
				items { product_id quantity unit_price }
			}
		}
	`
	result, err := s.hasuraClient.Query(ctx, query, args)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) executeSearchProducts(ctx context.Context, args map[string]interface{}) (string, error) {
	query := `
		query SearchProducts($query: String, $category_id: uuid, $limit: Int) {
			products(where: {name: {_ilike: $query}, category_id: {_eq: $category_id}}, limit: $limit) {
				id sku name base_price category { name }
			}
		}
	`
	result, err := s.hasuraClient.Query(ctx, query, args)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) executeCheckInventory(ctx context.Context, args map[string]interface{}) (string, error) {
	query := `
		query CheckInventory($product_id: uuid!, $warehouse_id: uuid) {
			stock_levels(where: {product_id: {_eq: $product_id}, warehouse_id: {_eq: $warehouse_id}}) {
				warehouse_id quantity_on_hand quantity_reserved quantity_available
			}
		}
	`
	result, err := s.hasuraClient.Query(ctx, query, args)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) executeCalculatePrice(ctx context.Context, args map[string]interface{}) (string, error) {
	// Call pricing-engine service
	result := map[string]interface{}{
		"base_price":      1000.00,
		"discount_amount": 100.00,
		"tax_amount":      67.50,
		"total_price":     967.50,
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) executeTriggerWorkflow(ctx context.Context, args map[string]interface{}) (string, error) {
	workflowName := args["workflow_name"].(string)
	payload, _ := args["payload"].(map[string]interface{})

	err := s.n8nClient.TriggerWorkflow(ctx, workflowName, payload)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Workflow '%s' triggered successfully", workflowName), nil
}

func (s *MCPServer) executeStartOrderWorkflow(ctx context.Context, args map[string]interface{}) (string, error) {
	runID, err := s.temporalClient.StartWorkflow(ctx, "OrderProcessingWorkflow", args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Order workflow started with run ID: %s", runID), nil
}

func (s *MCPServer) executeGetRecommendations(ctx context.Context, args map[string]interface{}) (string, error) {
	// Call recommendations service
	result := map[string]interface{}{
		"recommendations": []map[string]interface{}{
			{"product_id": "prod-1", "score": 0.95, "reason": "frequently_bought"},
			{"product_id": "prod-2", "score": 0.88, "reason": "similar_customers"},
		},
	}
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

func (s *MCPServer) handleListResources(c *gin.Context) {
	c.JSON(200, gin.H{"resources": s.resources})
}

func (s *MCPServer) handleReadResource(c *gin.Context) {
	var req struct {
		URI string `json:"uri"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	content, err := s.readResource(c.Request.Context(), req.URI)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"contents": []map[string]interface{}{
			{"uri": req.URI, "mimeType": "application/json", "text": content},
		},
	})
}

func (s *MCPServer) readResource(ctx context.Context, uri string) (string, error) {
	// Fetch resource data based on URI
	switch uri {
	case "omniroute://catalog/categories":
		result, _ := s.hasuraClient.Query(ctx, "query { categories { id name slug } }", nil)
		jsonResult, _ := json.MarshalIndent(result, "", "  ")
		return string(jsonResult), nil
	case "omniroute://analytics/dashboard":
		result := map[string]interface{}{
			"gmv_today":        1500000,
			"orders_today":     245,
			"active_customers": 1200,
		}
		jsonResult, _ := json.MarshalIndent(result, "", "  ")
		return string(jsonResult), nil
	default:
		return "{}", nil
	}
}

func (s *MCPServer) handleListPrompts(c *gin.Context) {
	c.JSON(200, gin.H{"prompts": s.prompts})
}

func (s *MCPServer) handleGetPrompt(c *gin.Context) {
	var req struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	prompt := s.generatePrompt(req.Name, req.Arguments)
	c.JSON(200, gin.H{
		"description": prompt.Description,
		"messages": []map[string]interface{}{
			{"role": "user", "content": map[string]string{"type": "text", "text": prompt.Content}},
		},
	})
}

type GeneratedPrompt struct {
	Description string
	Content     string
}

func (s *MCPServer) generatePrompt(name string, args map[string]string) GeneratedPrompt {
	switch name {
	case "order_assistant":
		return GeneratedPrompt{
			Description: "Order assistance for customer",
			Content:     fmt.Sprintf("You are an order assistant for customer %s. Help them browse products, place orders, and track deliveries.", args["customer_id"]),
		}
	case "inventory_analyst":
		return GeneratedPrompt{
			Description: "Inventory analysis",
			Content:     "Analyze current inventory levels, identify low stock items, and recommend restocking quantities based on demand forecasts.",
		}
	case "sales_reporter":
		return GeneratedPrompt{
			Description: "Sales analysis",
			Content:     fmt.Sprintf("Generate a sales report for the %s period. Include key metrics, trends, and actionable recommendations.", args["period"]),
		}
	default:
		return GeneratedPrompt{Description: "Unknown prompt", Content: ""}
	}
}

func (s *MCPServer) handleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Send initial connection event
	c.SSEvent("message", gin.H{"type": "connection", "status": "connected"})
	c.Writer.Flush()

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.SSEvent("ping", gin.H{"timestamp": time.Now().Unix()})
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

func (s *MCPServer) handleGraphQLProxy(c *gin.Context) {
	var req struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := s.hasuraClient.Query(c.Request.Context(), req.Query, req.Variables)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (s *MCPServer) handleWorkflowTrigger(c *gin.Context) {
	var req struct {
		WorkflowType string                 `json:"workflow_type"`
		Input        map[string]interface{} `json:"input"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Determine which system to use (n8n or Temporal)
	if req.WorkflowType == "order_processing" || req.WorkflowType == "credit_review" || req.WorkflowType == "worker_payout" {
		// Use Temporal for complex, long-running workflows
		runID, err := s.temporalClient.StartWorkflow(c.Request.Context(), req.WorkflowType, req.Input)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"system": "temporal", "run_id": runID})
	} else {
		// Use n8n for simpler automation
		err := s.n8nClient.TriggerWorkflow(c.Request.Context(), req.WorkflowType, req.Input)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"system": "n8n", "status": "triggered"})
	}
}
