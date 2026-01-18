// Package api provides HTTP handlers for the pricing engine
package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/omniroute/pricing-engine/internal/domain"
	"github.com/omniroute/pricing-engine/internal/engine"
)

// PricingHandler handles HTTP requests for pricing operations
type PricingHandler struct {
	engine *engine.PricingEngine
	logger *zap.Logger
}

// NewPricingHandler creates a new pricing handler
func NewPricingHandler(e *engine.PricingEngine, l *zap.Logger) *PricingHandler {
	return &PricingHandler{engine: e, logger: l}
}

// RegisterRoutes registers the handler routes
func (h *PricingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/api/v1/prices", h.CalculatePrices)
	mux.HandleFunc("/api/v1/prices/batch", h.BatchCalculatePrices)
}

// HealthCheck handles health check requests
func (h *PricingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// CalculatePrices handles price calculation requests
func (h *PricingHandler) CalculatePrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.PriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.engine.CalculatePrices(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to calculate prices", zap.Error(err))
		http.Error(w, "Price calculation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BatchCalculatePrices handles batch price calculation requests
func (h *PricingHandler) BatchCalculatePrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requests []domain.PriceRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		h.logger.Error("failed to decode batch request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	responses := make([]*domain.PriceResponse, 0, len(requests))
	for _, req := range requests {
		response, err := h.engine.CalculatePrices(r.Context(), &req)
		if err != nil {
			h.logger.Warn("failed to calculate price in batch", zap.Error(err))
			continue
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
