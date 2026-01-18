// Package api provides HTTP handlers for the pricing engine
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/omniroute/pricing-engine/internal/domain"
	"github.com/omniroute/pricing-engine/internal/engine"
)

// PricingHandler handles pricing-related HTTP requests
type PricingHandler struct {
	engine *engine.PricingEngine
	logger *zap.Logger
}

// NewPricingHandler creates a new pricing handler
func NewPricingHandler(engine *engine.PricingEngine, logger *zap.Logger) *PricingHandler {
	return &PricingHandler{
		engine: engine,
		logger: logger,
	}
}

// CalculatePriceRequest is the API request structure
type CalculatePriceRequest struct {
	CustomerID string              `json:"customer_id"`
	Items      []PriceItemRequest  `json:"items"`
	Currency   string              `json:"currency"`
	Channel    string              `json:"channel"`
}

// PriceItemRequest represents an item in the price request
type PriceItemRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id,omitempty"`
	Quantity  int    `json:"quantity"`
}

// CalculatePriceResponse is the API response structure
type CalculatePriceResponse struct {
	CustomerID        string               `json:"customer_id"`
	Items             []PriceItemResponse  `json:"items"`
	Summary           PriceSummary         `json:"summary"`
	AppliedPromotions []AppliedPromotion   `json:"applied_promotions,omitempty"`
	Currency          string               `json:"currency"`
	CalculatedAt      string               `json:"calculated_at"`
}

// PriceItemResponse represents a priced item in the response
type PriceItemResponse struct {
	ProductID       string             `json:"product_id"`
	VariantID       string             `json:"variant_id,omitempty"`
	SKU             string             `json:"sku"`
	Name            string             `json:"name"`
	Quantity        int                `json:"quantity"`
	BasePrice       string             `json:"base_price"`
	UnitPrice       string             `json:"unit_price"`
	OriginalPrice   string             `json:"original_price"`
	DiscountAmount  string             `json:"discount_amount"`
	DiscountPercent string             `json:"discount_percent"`
	TaxAmount       string             `json:"tax_amount"`
	LineTotal       string             `json:"line_total"`
	PriceSource     string             `json:"price_source"`
	PriceBreakdown  []PriceComponent   `json:"price_breakdown,omitempty"`
}

// PriceSummary contains order-level price totals
type PriceSummary struct {
	Subtotal      string `json:"subtotal"`
	TotalDiscount string `json:"total_discount"`
	TaxTotal      string `json:"tax_total"`
	GrandTotal    string `json:"grand_total"`
}

// PriceComponent represents a component of the price calculation
type PriceComponent struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Amount     string `json:"amount"`
	IsDiscount bool   `json:"is_discount"`
}

// AppliedPromotion represents an applied promotion
type AppliedPromotion struct {
	PromotionID    string   `json:"promotion_id"`
	Name           string   `json:"name"`
	Code           string   `json:"code,omitempty"`
	DiscountAmount string   `json:"discount_amount"`
	AppliedToItems []string `json:"applied_to_items,omitempty"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error       string            `json:"error"`
	Code        string            `json:"code"`
	Details     map[string]string `json:"details,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
}

// CalculatePrice handles the price calculation endpoint
func (h *PricingHandler) CalculatePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	
	// Extract tenant ID from context or header
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "tenant_required", "X-Tenant-ID header is required", requestID)
		return
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_tenant", "Invalid tenant ID format", requestID)
		return
	}
	
	// Parse request body
	var req CalculatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_json", "Invalid request body", requestID)
		return
	}
	
	// Validate request
	if req.CustomerID == "" {
		h.writeError(w, http.StatusBadRequest, "customer_required", "customer_id is required", requestID)
		return
	}
	
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_customer", "Invalid customer ID format", requestID)
		return
	}
	
	if len(req.Items) == 0 {
		h.writeError(w, http.StatusBadRequest, "items_required", "At least one item is required", requestID)
		return
	}
	
	// Convert to domain request
	domainReq := &domain.PriceRequest{
		TenantID:   tenantID,
		CustomerID: customerID,
		Items:      make([]domain.PriceRequestItem, 0, len(req.Items)),
		Currency:   req.Currency,
		Channel:    req.Channel,
		Timestamp:  time.Now(),
	}
	
	for _, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid_product", 
				"Invalid product ID: "+item.ProductID, requestID)
			return
		}
		
		var variantID uuid.UUID
		if item.VariantID != "" {
			variantID, err = uuid.Parse(item.VariantID)
			if err != nil {
				h.writeError(w, http.StatusBadRequest, "invalid_variant", 
					"Invalid variant ID: "+item.VariantID, requestID)
				return
			}
		}
		
		if item.Quantity <= 0 {
			h.writeError(w, http.StatusBadRequest, "invalid_quantity", 
				"Quantity must be positive", requestID)
			return
		}
		
		domainReq.Items = append(domainReq.Items, domain.PriceRequestItem{
			ProductID: productID,
			VariantID: variantID,
			Quantity:  item.Quantity,
		})
	}
	
	// Calculate prices
	result, err := h.engine.CalculatePrices(ctx, domainReq)
	if err != nil {
		h.logger.Error("price calculation failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.writeError(w, http.StatusInternalServerError, "calculation_failed", 
			"Failed to calculate prices", requestID)
		return
	}
	
	// Convert to API response
	response := h.convertToAPIResponse(result)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// convertToAPIResponse converts the domain response to API response
func (h *PricingHandler) convertToAPIResponse(result *domain.PriceResponse) *CalculatePriceResponse {
	response := &CalculatePriceResponse{
		CustomerID:   result.CustomerID.String(),
		Items:        make([]PriceItemResponse, 0, len(result.Items)),
		Currency:     result.Currency,
		CalculatedAt: result.CalculatedAt.Format(time.RFC3339),
		Summary: PriceSummary{
			Subtotal:      result.Subtotal.StringFixed(2),
			TotalDiscount: result.TotalDiscount.StringFixed(2),
			TaxTotal:      result.TaxTotal.StringFixed(2),
			GrandTotal:    result.GrandTotal.StringFixed(2),
		},
	}
	
	for _, item := range result.Items {
		apiItem := PriceItemResponse{
			ProductID:       item.ProductID.String(),
			SKU:             item.SKU,
			Name:            item.Name,
			Quantity:        item.Quantity,
			BasePrice:       item.BasePrice.StringFixed(2),
			UnitPrice:       item.UnitPrice.StringFixed(2),
			OriginalPrice:   item.OriginalPrice.StringFixed(2),
			DiscountAmount:  item.DiscountAmount.StringFixed(2),
			DiscountPercent: item.DiscountPercent.StringFixed(2),
			TaxAmount:       item.TaxAmount.StringFixed(2),
			LineTotal:       item.LineTotal.StringFixed(2),
			PriceSource:     string(item.PriceSource),
		}
		
		if item.VariantID != uuid.Nil {
			apiItem.VariantID = item.VariantID.String()
		}
		
		if len(item.PriceBreakdown) > 0 {
			apiItem.PriceBreakdown = make([]PriceComponent, 0, len(item.PriceBreakdown))
			for _, comp := range item.PriceBreakdown {
				apiItem.PriceBreakdown = append(apiItem.PriceBreakdown, PriceComponent{
					Type:       comp.Type,
					Name:       comp.Name,
					Amount:     comp.Amount.StringFixed(2),
					IsDiscount: comp.IsDiscount,
				})
			}
		}
		
		response.Items = append(response.Items, apiItem)
	}
	
	if len(result.AppliedPromotions) > 0 {
		response.AppliedPromotions = make([]AppliedPromotion, 0, len(result.AppliedPromotions))
		for _, promo := range result.AppliedPromotions {
			apiPromo := AppliedPromotion{
				PromotionID:    promo.PromotionID.String(),
				Name:           promo.Name,
				Code:           promo.Code,
				DiscountAmount: promo.DiscountAmount.StringFixed(2),
			}
			
			if len(promo.AppliedToItems) > 0 {
				apiPromo.AppliedToItems = make([]string, 0, len(promo.AppliedToItems))
				for _, id := range promo.AppliedToItems {
					apiPromo.AppliedToItems = append(apiPromo.AppliedToItems, id.String())
				}
			}
			
			response.AppliedPromotions = append(response.AppliedPromotions, apiPromo)
		}
	}
	
	return response
}

// writeError writes an error response
func (h *PricingHandler) writeError(w http.ResponseWriter, status int, code, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:     message,
		Code:      code,
		RequestID: requestID,
	})
}

// GetPrice handles the single product price lookup endpoint
func (h *PricingHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	
	// Extract parameters
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "tenant_required", "X-Tenant-ID header is required", requestID)
		return
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_tenant", "Invalid tenant ID format", requestID)
		return
	}
	
	query := r.URL.Query()
	
	customerIDStr := query.Get("customer_id")
	if customerIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "customer_required", "customer_id is required", requestID)
		return
	}
	
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_customer", "Invalid customer ID format", requestID)
		return
	}
	
	productIDStr := query.Get("product_id")
	if productIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "product_required", "product_id is required", requestID)
		return
	}
	
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_product", "Invalid product ID format", requestID)
		return
	}
	
	// Build request for single product
	domainReq := &domain.PriceRequest{
		TenantID:   tenantID,
		CustomerID: customerID,
		Items: []domain.PriceRequestItem{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
		Currency:  query.Get("currency"),
		Timestamp: time.Now(),
	}
	
	// Calculate price
	result, err := h.engine.CalculatePrices(ctx, domainReq)
	if err != nil {
		h.logger.Error("price lookup failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.writeError(w, http.StatusInternalServerError, "lookup_failed", 
			"Failed to lookup price", requestID)
		return
	}
	
	if len(result.Items) == 0 {
		h.writeError(w, http.StatusNotFound, "product_not_found", 
			"Product not found", requestID)
		return
	}
	
	// Return single item response
	item := result.Items[0]
	response := PriceItemResponse{
		ProductID:       item.ProductID.String(),
		SKU:             item.SKU,
		Name:            item.Name,
		Quantity:        item.Quantity,
		BasePrice:       item.BasePrice.StringFixed(2),
		UnitPrice:       item.UnitPrice.StringFixed(2),
		OriginalPrice:   item.OriginalPrice.StringFixed(2),
		DiscountAmount:  item.DiscountAmount.StringFixed(2),
		DiscountPercent: item.DiscountPercent.StringFixed(2),
		TaxAmount:       item.TaxAmount.StringFixed(2),
		LineTotal:       item.LineTotal.StringFixed(2),
		PriceSource:     string(item.PriceSource),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// BulkPriceRequest represents a bulk price lookup request
type BulkPriceRequest struct {
	CustomerID string   `json:"customer_id"`
	ProductIDs []string `json:"product_ids"`
	Currency   string   `json:"currency"`
}

// BulkPriceResponse represents a bulk price lookup response
type BulkPriceResponse struct {
	CustomerID string                     `json:"customer_id"`
	Prices     map[string]BulkPriceItem   `json:"prices"`
	Currency   string                     `json:"currency"`
}

// BulkPriceItem represents a single product's price in bulk response
type BulkPriceItem struct {
	ProductID   string `json:"product_id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	BasePrice   string `json:"base_price"`
	UnitPrice   string `json:"unit_price"`
	PriceSource string `json:"price_source"`
	Available   bool   `json:"available"`
}

// BulkGetPrices handles bulk price lookup
func (h *PricingHandler) BulkGetPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "tenant_required", "X-Tenant-ID header is required", requestID)
		return
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_tenant", "Invalid tenant ID format", requestID)
		return
	}
	
	var req BulkPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_json", "Invalid request body", requestID)
		return
	}
	
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_customer", "Invalid customer ID format", requestID)
		return
	}
	
	// Build domain request
	items := make([]domain.PriceRequestItem, 0, len(req.ProductIDs))
	for _, pidStr := range req.ProductIDs {
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			continue
		}
		items = append(items, domain.PriceRequestItem{
			ProductID: pid,
			Quantity:  1,
		})
	}
	
	domainReq := &domain.PriceRequest{
		TenantID:   tenantID,
		CustomerID: customerID,
		Items:      items,
		Currency:   req.Currency,
		Timestamp:  time.Now(),
	}
	
	result, err := h.engine.CalculatePrices(ctx, domainReq)
	if err != nil {
		h.logger.Error("bulk price lookup failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.writeError(w, http.StatusInternalServerError, "lookup_failed", 
			"Failed to lookup prices", requestID)
		return
	}
	
	// Build response
	response := BulkPriceResponse{
		CustomerID: req.CustomerID,
		Prices:     make(map[string]BulkPriceItem),
		Currency:   result.Currency,
	}
	
	for _, item := range result.Items {
		response.Prices[item.ProductID.String()] = BulkPriceItem{
			ProductID:   item.ProductID.String(),
			SKU:         item.SKU,
			Name:        item.Name,
			BasePrice:   item.BasePrice.StringFixed(2),
			UnitPrice:   item.UnitPrice.StringFixed(2),
			PriceSource: string(item.PriceSource),
			Available:   item.UnitPrice.GreaterThan(decimal.Zero),
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck handles the health check endpoint
func (h *PricingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutes registers all pricing routes
func (h *PricingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/prices/calculate", h.CalculatePrice)
	mux.HandleFunc("GET /api/v1/prices", h.GetPrice)
	mux.HandleFunc("POST /api/v1/prices/bulk", h.BulkGetPrices)
	mux.HandleFunc("GET /health", h.HealthCheck)
}
