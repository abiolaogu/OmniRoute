// Package engine contains the core pricing calculation logic
package engine

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/omniroute/pricing-engine/internal/domain"
)

// PricingEngine is the main pricing calculation engine
type PricingEngine struct {
	productRepo      ProductRepository
	customerRepo     CustomerRepository
	priceListRepo    PriceListRepository
	volumeDiscountRepo VolumeDiscountRepository
	contractPriceRepo ContractPriceRepository
	promotionRepo    PromotionRepository
	taxRepo          TaxRepository
	cache            PriceCache
	logger           *zap.Logger
	
	// Configuration
	config           *PricingConfig
}

// PricingConfig holds configuration for the pricing engine
type PricingConfig struct {
	EnableCaching        bool
	CacheTTL             time.Duration
	MaxConcurrentCalcs   int
	EnableVolumeDiscounts bool
	EnablePromotions     bool
	EnableContractPricing bool
	RoundingPrecision    int32
	RoundingMode         string // "up", "down", "half_up", "half_down"
}

// DefaultConfig returns the default pricing configuration
func DefaultConfig() *PricingConfig {
	return &PricingConfig{
		EnableCaching:        true,
		CacheTTL:             5 * time.Minute,
		MaxConcurrentCalcs:   100,
		EnableVolumeDiscounts: true,
		EnablePromotions:     true,
		EnableContractPricing: true,
		RoundingPrecision:    2,
		RoundingMode:         "half_up",
	}
}

// Repository interfaces
type ProductRepository interface {
	GetByID(ctx context.Context, tenantID, productID uuid.UUID) (*domain.Product, error)
	GetByIDs(ctx context.Context, tenantID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]*domain.Product, error)
	GetVariant(ctx context.Context, productID, variantID uuid.UUID) (*domain.ProductVariant, error)
	GetVariants(ctx context.Context, productID uuid.UUID, variantIDs []uuid.UUID) (map[uuid.UUID]*domain.ProductVariant, error)
}

type CustomerRepository interface {
	GetByID(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Customer, error)
}

type PriceListRepository interface {
	GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, now time.Time) ([]*domain.PriceList, error)
	GetItems(ctx context.Context, priceListIDs []uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID][]*domain.PriceListItem, error)
}

type VolumeDiscountRepository interface {
	GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, productIDs []uuid.UUID, now time.Time) ([]*domain.VolumeDiscount, error)
}

type ContractPriceRepository interface {
	GetForCustomer(ctx context.Context, customerID uuid.UUID, productIDs []uuid.UUID, now time.Time) (map[uuid.UUID]*domain.ContractPrice, error)
}

type PromotionRepository interface {
	GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, productIDs []uuid.UUID, now time.Time) ([]*domain.Promotion, error)
}

type TaxRepository interface {
	GetRates(ctx context.Context, tenantID uuid.UUID, categories []string, country, state string) (map[string]*domain.TaxRate, error)
}

type PriceCache interface {
	Get(ctx context.Context, key string) (*domain.PriceResponse, error)
	Set(ctx context.Context, key string, response *domain.PriceResponse, ttl time.Duration) error
	Invalidate(ctx context.Context, patterns []string) error
}

// NewPricingEngine creates a new pricing engine instance
func NewPricingEngine(
	productRepo ProductRepository,
	customerRepo CustomerRepository,
	priceListRepo PriceListRepository,
	volumeDiscountRepo VolumeDiscountRepository,
	contractPriceRepo ContractPriceRepository,
	promotionRepo PromotionRepository,
	taxRepo TaxRepository,
	cache PriceCache,
	logger *zap.Logger,
	config *PricingConfig,
) *PricingEngine {
	if config == nil {
		config = DefaultConfig()
	}
	
	return &PricingEngine{
		productRepo:        productRepo,
		customerRepo:       customerRepo,
		priceListRepo:      priceListRepo,
		volumeDiscountRepo: volumeDiscountRepo,
		contractPriceRepo:  contractPriceRepo,
		promotionRepo:      promotionRepo,
		taxRepo:            taxRepo,
		cache:              cache,
		logger:             logger,
		config:             config,
	}
}

// CalculatePrices calculates prices for a given request
func (e *PricingEngine) CalculatePrices(ctx context.Context, req *domain.PriceRequest) (*domain.PriceResponse, error) {
	startTime := time.Now()
	
	// Validate request
	if err := e.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid price request: %w", err)
	}
	
	// Check cache
	if e.config.EnableCaching {
		cacheKey := e.generateCacheKey(req)
		if cached, err := e.cache.Get(ctx, cacheKey); err == nil && cached != nil {
			e.logger.Debug("cache hit", zap.String("key", cacheKey))
			return cached, nil
		}
	}
	
	// Load required data concurrently
	data, err := e.loadPricingData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to load pricing data: %w", err)
	}
	
	// Calculate prices for each item
	response := &domain.PriceResponse{
		TenantID:          req.TenantID,
		CustomerID:        req.CustomerID,
		Items:             make([]domain.PriceResponseItem, 0, len(req.Items)),
		Currency:          req.Currency,
		AppliedPromotions: make([]domain.AppliedPromotion, 0),
		CalculatedAt:      time.Now(),
	}
	
	// Process items with concurrent calculation for large requests
	if len(req.Items) > 10 {
		response.Items = e.calculateItemsPricesConcurrent(ctx, req, data)
	} else {
		for _, item := range req.Items {
			calculated := e.calculateItemPrice(ctx, item, req, data)
			response.Items = append(response.Items, calculated)
		}
	}
	
	// Apply order-level promotions
	if e.config.EnablePromotions {
		e.applyOrderPromotions(ctx, response, data)
	}
	
	// Calculate totals
	e.calculateTotals(response)
	
	// Apply tax
	e.applyTax(ctx, response, data)
	
	// Round final values
	e.roundPrices(response)
	
	// Cache result
	if e.config.EnableCaching {
		cacheKey := e.generateCacheKey(req)
		if err := e.cache.Set(ctx, cacheKey, response, e.config.CacheTTL); err != nil {
			e.logger.Warn("failed to cache price response", zap.Error(err))
		}
	}
	
	e.logger.Info("price calculation completed",
		zap.Int("items", len(response.Items)),
		zap.Duration("duration", time.Since(startTime)),
	)
	
	return response, nil
}

// pricingData holds all the loaded data needed for price calculation
type pricingData struct {
	customer        *domain.Customer
	products        map[uuid.UUID]*domain.Product
	variants        map[uuid.UUID]*domain.ProductVariant
	priceLists      []*domain.PriceList
	priceListItems  map[uuid.UUID][]*domain.PriceListItem
	volumeDiscounts []*domain.VolumeDiscount
	contractPrices  map[uuid.UUID]*domain.ContractPrice
	promotions      []*domain.Promotion
	taxRates        map[string]*domain.TaxRate
}

// loadPricingData loads all required data for price calculation
func (e *PricingEngine) loadPricingData(ctx context.Context, req *domain.PriceRequest) (*pricingData, error) {
	data := &pricingData{
		products: make(map[uuid.UUID]*domain.Product),
		variants: make(map[uuid.UUID]*domain.ProductVariant),
	}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 7)
	
	// Extract product and variant IDs
	productIDs := make([]uuid.UUID, 0, len(req.Items))
	variantIDs := make([]uuid.UUID, 0)
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
		if item.VariantID != uuid.Nil {
			variantIDs = append(variantIDs, item.VariantID)
		}
	}
	
	// Load customer
	wg.Add(1)
	go func() {
		defer wg.Done()
		customer, err := e.customerRepo.GetByID(ctx, req.TenantID, req.CustomerID)
		if err != nil {
			errChan <- fmt.Errorf("failed to load customer: %w", err)
			return
		}
		mu.Lock()
		data.customer = customer
		mu.Unlock()
	}()
	
	// Load products
	wg.Add(1)
	go func() {
		defer wg.Done()
		products, err := e.productRepo.GetByIDs(ctx, req.TenantID, productIDs)
		if err != nil {
			errChan <- fmt.Errorf("failed to load products: %w", err)
			return
		}
		mu.Lock()
		data.products = products
		mu.Unlock()
	}()
	
	// Load variants if any
	if len(variantIDs) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, item := range req.Items {
				if item.VariantID != uuid.Nil {
					variant, err := e.productRepo.GetVariant(ctx, item.ProductID, item.VariantID)
					if err != nil {
						e.logger.Warn("failed to load variant", 
							zap.String("variant_id", item.VariantID.String()),
							zap.Error(err))
						continue
					}
					mu.Lock()
					data.variants[item.VariantID] = variant
					mu.Unlock()
				}
			}
		}()
	}
	
	// Wait for customer to be loaded first (needed for other queries)
	wg.Wait()
	
	// Check for errors so far
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	
	if data.customer == nil {
		return nil, fmt.Errorf("customer not found: %s", req.CustomerID)
	}
	
	// Load price lists
	wg.Add(1)
	go func() {
		defer wg.Done()
		priceLists, err := e.priceListRepo.GetApplicable(ctx, req.TenantID, data.customer, req.Timestamp)
		if err != nil {
			errChan <- fmt.Errorf("failed to load price lists: %w", err)
			return
		}
		
		if len(priceLists) > 0 {
			priceListIDs := make([]uuid.UUID, len(priceLists))
			for i, pl := range priceLists {
				priceListIDs[i] = pl.ID
			}
			
			items, err := e.priceListRepo.GetItems(ctx, priceListIDs, productIDs)
			if err != nil {
				errChan <- fmt.Errorf("failed to load price list items: %w", err)
				return
			}
			
			mu.Lock()
			data.priceLists = priceLists
			data.priceListItems = items
			mu.Unlock()
		}
	}()
	
	// Load volume discounts
	if e.config.EnableVolumeDiscounts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			discounts, err := e.volumeDiscountRepo.GetApplicable(ctx, req.TenantID, data.customer, productIDs, req.Timestamp)
			if err != nil {
				e.logger.Warn("failed to load volume discounts", zap.Error(err))
				return
			}
			mu.Lock()
			data.volumeDiscounts = discounts
			mu.Unlock()
		}()
	}
	
	// Load contract prices
	if e.config.EnableContractPricing {
		wg.Add(1)
		go func() {
			defer wg.Done()
			contracts, err := e.contractPriceRepo.GetForCustomer(ctx, req.CustomerID, productIDs, req.Timestamp)
			if err != nil {
				e.logger.Warn("failed to load contract prices", zap.Error(err))
				return
			}
			mu.Lock()
			data.contractPrices = contracts
			mu.Unlock()
		}()
	}
	
	// Load promotions
	if e.config.EnablePromotions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			promos, err := e.promotionRepo.GetApplicable(ctx, req.TenantID, data.customer, productIDs, req.Timestamp)
			if err != nil {
				e.logger.Warn("failed to load promotions", zap.Error(err))
				return
			}
			mu.Lock()
			data.promotions = promos
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	// Check for any remaining errors
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	
	return data, nil
}

// calculateItemPrice calculates the price for a single item
func (e *PricingEngine) calculateItemPrice(
	ctx context.Context,
	item domain.PriceRequestItem,
	req *domain.PriceRequest,
	data *pricingData,
) domain.PriceResponseItem {
	product, ok := data.products[item.ProductID]
	if !ok {
		e.logger.Warn("product not found", zap.String("product_id", item.ProductID.String()))
		return domain.PriceResponseItem{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		}
	}
	
	result := domain.PriceResponseItem{
		ProductID:      item.ProductID,
		VariantID:      item.VariantID,
		SKU:            product.SKU,
		Name:           product.Name,
		Quantity:       item.Quantity,
		BasePrice:      product.BasePrice,
		OriginalPrice:  product.BasePrice,
		PriceBreakdown: make([]domain.PriceComponent, 0),
	}
	
	// Adjust base price for variant
	if item.VariantID != uuid.Nil {
		if variant, ok := data.variants[item.VariantID]; ok {
			result.BasePrice = result.BasePrice.Add(variant.PriceAdjustment)
			result.OriginalPrice = result.BasePrice
			result.SKU = variant.SKU
			if variant.Name != "" {
				result.Name = fmt.Sprintf("%s - %s", product.Name, variant.Name)
			}
		}
	}
	
	result.PriceBreakdown = append(result.PriceBreakdown, domain.PriceComponent{
		Type:       "base",
		Name:       "Base Price",
		Amount:     result.BasePrice,
		IsDiscount: false,
		Priority:   0,
	})
	
	// Determine the best price from various sources
	bestPrice := result.BasePrice
	bestSource := domain.PriceSourceBase
	var bestSourceID uuid.UUID
	
	// Check contract price first (highest priority)
	if e.config.EnableContractPricing && data.contractPrices != nil {
		if contract, ok := data.contractPrices[item.ProductID]; ok {
			if item.Quantity >= contract.MinQuantity {
				bestPrice = contract.Price
				bestSource = domain.PriceSourceContract
				bestSourceID = contract.ID
				
				result.PriceBreakdown = append(result.PriceBreakdown, domain.PriceComponent{
					Type:       "contract",
					Name:       fmt.Sprintf("Contract Price (%s)", contract.ContractReference),
					SourceID:   contract.ID,
					Amount:     contract.Price,
					IsDiscount: false,
					Priority:   100,
				})
			}
		}
	}
	
	// Check price lists (if no contract price)
	if bestSource == domain.PriceSourceBase && data.priceLists != nil {
		priceListPrice, priceListID := e.getBestPriceListPrice(item, data)
		if priceListPrice.LessThan(bestPrice) {
			bestPrice = priceListPrice
			bestSource = domain.PriceSourcePriceList
			bestSourceID = priceListID
			
			result.PriceBreakdown = append(result.PriceBreakdown, domain.PriceComponent{
				Type:       "price_list",
				Name:       "Price List",
				SourceID:   priceListID,
				Amount:     priceListPrice,
				IsDiscount: false,
				Priority:   50,
			})
		}
	}
	
	result.UnitPrice = bestPrice
	result.PriceSource = bestSource
	result.PriceSourceID = bestSourceID
	
	// Calculate discount from original price
	if result.OriginalPrice.GreaterThan(result.UnitPrice) {
		result.DiscountAmount = result.OriginalPrice.Sub(result.UnitPrice)
		result.DiscountPercent = result.DiscountAmount.Div(result.OriginalPrice).Mul(decimal.NewFromInt(100))
	}
	
	// Apply volume discounts
	if e.config.EnableVolumeDiscounts && data.volumeDiscounts != nil {
		volumeDiscount := e.calculateVolumeDiscount(product, item.Quantity, data)
		if volumeDiscount.GreaterThan(decimal.Zero) {
			result.DiscountAmount = result.DiscountAmount.Add(volumeDiscount)
			result.UnitPrice = result.UnitPrice.Sub(volumeDiscount)
			
			result.PriceBreakdown = append(result.PriceBreakdown, domain.PriceComponent{
				Type:       "volume_discount",
				Name:       "Volume Discount",
				Amount:     volumeDiscount.Neg(),
				IsDiscount: true,
				Priority:   60,
			})
		}
	}
	
	// Apply item-level promotions
	if e.config.EnablePromotions && data.promotions != nil {
		promoDiscount := e.calculateItemPromotion(product, item.Quantity, data)
		if promoDiscount.GreaterThan(decimal.Zero) {
			result.DiscountAmount = result.DiscountAmount.Add(promoDiscount)
			result.UnitPrice = result.UnitPrice.Sub(promoDiscount)
			
			result.PriceBreakdown = append(result.PriceBreakdown, domain.PriceComponent{
				Type:       "promotion",
				Name:       "Promotional Discount",
				Amount:     promoDiscount.Neg(),
				IsDiscount: true,
				Priority:   70,
			})
		}
	}
	
	// Calculate line total
	result.LineTotal = result.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
	
	return result
}

// calculateItemsPricesConcurrent calculates prices for multiple items concurrently
func (e *PricingEngine) calculateItemsPricesConcurrent(
	ctx context.Context,
	req *domain.PriceRequest,
	data *pricingData,
) []domain.PriceResponseItem {
	results := make([]domain.PriceResponseItem, len(req.Items))
	
	sem := make(chan struct{}, e.config.MaxConcurrentCalcs)
	var wg sync.WaitGroup
	
	for i, item := range req.Items {
		wg.Add(1)
		go func(idx int, itm domain.PriceRequestItem) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			
			results[idx] = e.calculateItemPrice(ctx, itm, req, data)
		}(i, item)
	}
	
	wg.Wait()
	return results
}

// getBestPriceListPrice finds the best price from applicable price lists
func (e *PricingEngine) getBestPriceListPrice(
	item domain.PriceRequestItem,
	data *pricingData,
) (decimal.Decimal, uuid.UUID) {
	product := data.products[item.ProductID]
	if product == nil {
		return decimal.Zero, uuid.Nil
	}
	
	bestPrice := product.BasePrice
	var bestPriceListID uuid.UUID
	
	// Sort price lists by priority (higher priority first)
	sortedPriceLists := make([]*domain.PriceList, len(data.priceLists))
	copy(sortedPriceLists, data.priceLists)
	sort.Slice(sortedPriceLists, func(i, j int) bool {
		return sortedPriceLists[i].Priority > sortedPriceLists[j].Priority
	})
	
	for _, pl := range sortedPriceLists {
		items, ok := data.priceListItems[pl.ID]
		if !ok {
			continue
		}
		
		for _, pli := range items {
			if pli.ProductID != item.ProductID {
				continue
			}
			
			// Check variant match
			if item.VariantID != uuid.Nil && pli.VariantID != uuid.Nil && pli.VariantID != item.VariantID {
				continue
			}
			
			// Check quantity requirements
			if item.Quantity < pli.MinQuantity {
				continue
			}
			if pli.MaxQuantity > 0 && item.Quantity > pli.MaxQuantity {
				continue
			}
			
			// Calculate price based on pricing method
			var calculatedPrice decimal.Decimal
			switch pli.PricingMethod {
			case domain.PricingMethodFixed:
				calculatedPrice = pli.Price
			case domain.PricingMethodDiscountPercent:
				discount := bestPrice.Mul(pli.DiscountPercent).Div(decimal.NewFromInt(100))
				calculatedPrice = bestPrice.Sub(discount)
			case domain.PricingMethodDiscountAmount:
				calculatedPrice = bestPrice.Sub(pli.DiscountAmount)
			case domain.PricingMethodMargin:
				if product.CostPrice.GreaterThan(decimal.Zero) {
					calculatedPrice = product.CostPrice.Mul(decimal.NewFromInt(1).Add(pli.MarginPercent.Div(decimal.NewFromInt(100))))
				}
			default:
				continue
			}
			
			if calculatedPrice.LessThan(bestPrice) && calculatedPrice.GreaterThan(decimal.Zero) {
				bestPrice = calculatedPrice
				bestPriceListID = pl.ID
			}
		}
	}
	
	return bestPrice, bestPriceListID
}

// calculateVolumeDiscount calculates volume-based discounts
func (e *PricingEngine) calculateVolumeDiscount(
	product *domain.Product,
	quantity int,
	data *pricingData,
) decimal.Decimal {
	var totalDiscount decimal.Decimal
	
	for _, vd := range data.volumeDiscounts {
		if !e.volumeDiscountApplies(vd, product, data.customer) {
			continue
		}
		
		// Find applicable tier
		for _, tier := range vd.Tiers {
			if quantity >= tier.MinQuantity && (tier.MaxQuantity == 0 || quantity <= tier.MaxQuantity) {
				if tier.DiscountPercent.GreaterThan(decimal.Zero) {
					discount := product.BasePrice.Mul(tier.DiscountPercent).Div(decimal.NewFromInt(100))
					totalDiscount = totalDiscount.Add(discount)
				} else if tier.DiscountAmount.GreaterThan(decimal.Zero) {
					totalDiscount = totalDiscount.Add(tier.DiscountAmount)
				}
				
				// If can't combine, take the first matching discount
				if !vd.CanCombine {
					return totalDiscount
				}
				break
			}
		}
	}
	
	return totalDiscount
}

// volumeDiscountApplies checks if a volume discount applies to a product/customer
func (e *PricingEngine) volumeDiscountApplies(
	vd *domain.VolumeDiscount,
	product *domain.Product,
	customer *domain.Customer,
) bool {
	// Check customer type
	if len(vd.CustomerTypes) > 0 {
		found := false
		for _, ct := range vd.CustomerTypes {
			if ct == customer.Type {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check customer tier
	if len(vd.CustomerTiers) > 0 {
		found := false
		for _, tier := range vd.CustomerTiers {
			if tier == customer.Tier {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check product/category/brand
	switch vd.AppliesTo {
	case "all":
		return true
	case "product":
		for _, pid := range vd.ProductIDs {
			if pid == product.ID {
				return true
			}
		}
		return false
	case "category":
		for _, cid := range vd.CategoryIDs {
			if cid == product.CategoryID {
				return true
			}
		}
		return false
	case "brand":
		for _, brand := range vd.BrandNames {
			if brand == product.Brand {
				return true
			}
		}
		return false
	}
	
	return false
}

// calculateItemPromotion calculates item-level promotional discounts
func (e *PricingEngine) calculateItemPromotion(
	product *domain.Product,
	quantity int,
	data *pricingData,
) decimal.Decimal {
	var totalDiscount decimal.Decimal
	
	for _, promo := range data.promotions {
		if promo.Type != domain.PromotionTypeDiscount {
			continue
		}
		
		if !e.promotionApplies(promo, product, data.customer) {
			continue
		}
		
		if quantity < promo.MinQuantity {
			continue
		}
		
		var discount decimal.Decimal
		switch promo.DiscountType {
		case domain.DiscountTypePercent:
			discount = product.BasePrice.Mul(promo.DiscountValue).Div(decimal.NewFromInt(100))
		case domain.DiscountTypeFixed:
			discount = promo.DiscountValue
		}
		
		if discount.GreaterThan(decimal.Zero) {
			totalDiscount = totalDiscount.Add(discount)
		}
		
		if !promo.CanCombine {
			break
		}
	}
	
	return totalDiscount
}

// promotionApplies checks if a promotion applies to a product/customer
func (e *PricingEngine) promotionApplies(
	promo *domain.Promotion,
	product *domain.Product,
	customer *domain.Customer,
) bool {
	// Check customer type
	if len(promo.CustomerTypes) > 0 {
		found := false
		for _, ct := range promo.CustomerTypes {
			if ct == customer.Type {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check product applicability
	switch promo.AppliesTo {
	case "all":
		return true
	case "product":
		for _, pid := range promo.ProductIDs {
			if pid == product.ID {
				return true
			}
		}
		return false
	case "category":
		for _, cid := range promo.CategoryIDs {
			if cid == product.CategoryID {
				return true
			}
		}
		return false
	}
	
	return false
}

// applyOrderPromotions applies order-level promotions
func (e *PricingEngine) applyOrderPromotions(
	ctx context.Context,
	response *domain.PriceResponse,
	data *pricingData,
) {
	for _, promo := range data.promotions {
		if promo.Type == domain.PromotionTypeDiscount {
			continue // Already handled at item level
		}
		
		// Calculate order subtotal for threshold check
		subtotal := decimal.Zero
		for _, item := range response.Items {
			subtotal = subtotal.Add(item.LineTotal)
		}
		
		if promo.MinOrderValue.GreaterThan(decimal.Zero) && subtotal.LessThan(promo.MinOrderValue) {
			continue
		}
		
		appliedPromo := domain.AppliedPromotion{
			PromotionID:    promo.ID,
			Name:           promo.Name,
			Code:           promo.Code,
			AppliedToItems: make([]uuid.UUID, 0),
		}
		
		switch promo.Type {
		case domain.PromotionTypeBuyXGetY:
			discount := e.applyBuyXGetY(response, promo, data)
			appliedPromo.DiscountAmount = discount
		case domain.PromotionTypeFreeShipping:
			appliedPromo.DiscountAmount = decimal.Zero // Handled elsewhere
		}
		
		if appliedPromo.DiscountAmount.GreaterThan(decimal.Zero) {
			response.AppliedPromotions = append(response.AppliedPromotions, appliedPromo)
		}
	}
}

// applyBuyXGetY applies buy X get Y promotions
func (e *PricingEngine) applyBuyXGetY(
	response *domain.PriceResponse,
	promo *domain.Promotion,
	data *pricingData,
) decimal.Decimal {
	var totalDiscount decimal.Decimal
	
	for i := range response.Items {
		item := &response.Items[i]
		product := data.products[item.ProductID]
		
		if !e.promotionApplies(promo, product, data.customer) {
			continue
		}
		
		// Calculate free items
		sets := item.Quantity / promo.BuyQuantity
		freeItems := sets * promo.GetQuantity
		
		if freeItems > 0 {
			discount := item.UnitPrice.Mul(decimal.NewFromInt(int64(freeItems)))
			totalDiscount = totalDiscount.Add(discount)
		}
	}
	
	return totalDiscount
}

// calculateTotals calculates order totals
func (e *PricingEngine) calculateTotals(response *domain.PriceResponse) {
	response.Subtotal = decimal.Zero
	response.TotalDiscount = decimal.Zero
	
	for _, item := range response.Items {
		response.Subtotal = response.Subtotal.Add(item.LineTotal)
		response.TotalDiscount = response.TotalDiscount.Add(
			item.DiscountAmount.Mul(decimal.NewFromInt(int64(item.Quantity))),
		)
	}
	
	// Add order-level promotion discounts
	for _, promo := range response.AppliedPromotions {
		response.TotalDiscount = response.TotalDiscount.Add(promo.DiscountAmount)
	}
	
	response.GrandTotal = response.Subtotal.Sub(response.TotalDiscount)
}

// applyTax applies tax to the order
func (e *PricingEngine) applyTax(
	ctx context.Context,
	response *domain.PriceResponse,
	data *pricingData,
) {
	// Simplified tax calculation - would need more complex logic for production
	response.TaxTotal = decimal.Zero
	
	for i := range response.Items {
		product := data.products[response.Items[i].ProductID]
		if product == nil {
			continue
		}
		
		// Default VAT rate for Nigeria: 7.5%
		taxRate := decimal.NewFromFloat(0.075)
		
		if data.taxRates != nil {
			if rate, ok := data.taxRates[product.TaxCategory]; ok {
				taxRate = rate.Rate.Div(decimal.NewFromInt(100))
			}
		}
		
		response.Items[i].TaxAmount = response.Items[i].LineTotal.Mul(taxRate)
		response.TaxTotal = response.TaxTotal.Add(response.Items[i].TaxAmount)
	}
	
	response.GrandTotal = response.GrandTotal.Add(response.TaxTotal)
}

// roundPrices rounds all monetary values
func (e *PricingEngine) roundPrices(response *domain.PriceResponse) {
	precision := e.config.RoundingPrecision
	
	for i := range response.Items {
		response.Items[i].UnitPrice = response.Items[i].UnitPrice.Round(precision)
		response.Items[i].DiscountAmount = response.Items[i].DiscountAmount.Round(precision)
		response.Items[i].TaxAmount = response.Items[i].TaxAmount.Round(precision)
		response.Items[i].LineTotal = response.Items[i].LineTotal.Round(precision)
	}
	
	response.Subtotal = response.Subtotal.Round(precision)
	response.TotalDiscount = response.TotalDiscount.Round(precision)
	response.TaxTotal = response.TaxTotal.Round(precision)
	response.GrandTotal = response.GrandTotal.Round(precision)
}

// validateRequest validates a price request
func (e *PricingEngine) validateRequest(req *domain.PriceRequest) error {
	if req.TenantID == uuid.Nil {
		return fmt.Errorf("tenant_id is required")
	}
	if req.CustomerID == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	
	for i, item := range req.Items {
		if item.ProductID == uuid.Nil {
			return fmt.Errorf("item %d: product_id is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be positive", i)
		}
	}
	
	if req.Currency == "" {
		req.Currency = "NGN"
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}
	
	return nil
}

// generateCacheKey generates a cache key for a price request
func (e *PricingEngine) generateCacheKey(req *domain.PriceRequest) string {
	// Create a deterministic key based on the request parameters
	key := fmt.Sprintf("price:%s:%s:%s",
		req.TenantID.String(),
		req.CustomerID.String(),
		req.Currency,
	)
	
	// Sort items by product ID for consistent key generation
	sortedItems := make([]domain.PriceRequestItem, len(req.Items))
	copy(sortedItems, req.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].ProductID.String() < sortedItems[j].ProductID.String()
	})
	
	for _, item := range sortedItems {
		key += fmt.Sprintf(":%s:%s:%d",
			item.ProductID.String(),
			item.VariantID.String(),
			item.Quantity,
		)
	}
	
	return key
}
