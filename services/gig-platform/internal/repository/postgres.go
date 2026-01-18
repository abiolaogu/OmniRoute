// Package repository provides database access implementations
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/omniroute/pricing-engine/internal/domain"
)

// PostgresProductRepository implements ProductRepository using PostgreSQL
type PostgresProductRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresProductRepository creates a new PostgreSQL product repository
func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{pool: pool}
}

// GetByID retrieves a single product by ID
func (r *PostgresProductRepository) GetByID(ctx context.Context, tenantID, productID uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, sku, name, category_id, brand, base_price, cost_price,
			   currency, tax_category, min_order_quantity, max_order_quantity,
			   order_multiple, status, visibility
		FROM products
		WHERE tenant_id = $1 AND id = $2 AND status = 'active'
	`
	
	var p domain.Product
	var categoryID, costPrice *string
	var maxOrderQty *int
	
	err := r.pool.QueryRow(ctx, query, tenantID, productID).Scan(
		&p.ID, &p.TenantID, &p.SKU, &p.Name, &categoryID, &p.Brand,
		&p.BasePrice, &costPrice, &p.Currency, &p.TaxCategory,
		&p.MinOrderQty, &maxOrderQty, &p.OrderMultiple, &p.Status, &p.Visibility,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found: %s", productID)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	if categoryID != nil {
		p.CategoryID, _ = uuid.Parse(*categoryID)
	}
	if costPrice != nil {
		p.CostPrice, _ = decimal.NewFromString(*costPrice)
	}
	if maxOrderQty != nil {
		p.MaxOrderQty = *maxOrderQty
	}
	
	return &p, nil
}

// GetByIDs retrieves multiple products by their IDs
func (r *PostgresProductRepository) GetByIDs(ctx context.Context, tenantID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]*domain.Product, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID]*domain.Product), nil
	}
	
	query := `
		SELECT id, tenant_id, sku, name, category_id, brand, base_price, cost_price,
			   currency, tax_category, min_order_quantity, max_order_quantity,
			   order_multiple, status, visibility
		FROM products
		WHERE tenant_id = $1 AND id = ANY($2) AND status = 'active'
	`
	
	rows, err := r.pool.Query(ctx, query, tenantID, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()
	
	products := make(map[uuid.UUID]*domain.Product)
	
	for rows.Next() {
		var p domain.Product
		var categoryID, costPrice *string
		var maxOrderQty *int
		
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.SKU, &p.Name, &categoryID, &p.Brand,
			&p.BasePrice, &costPrice, &p.Currency, &p.TaxCategory,
			&p.MinOrderQty, &maxOrderQty, &p.OrderMultiple, &p.Status, &p.Visibility,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		
		if categoryID != nil {
			p.CategoryID, _ = uuid.Parse(*categoryID)
		}
		if costPrice != nil {
			p.CostPrice, _ = decimal.NewFromString(*costPrice)
		}
		if maxOrderQty != nil {
			p.MaxOrderQty = *maxOrderQty
		}
		
		products[p.ID] = &p
	}
	
	return products, nil
}

// GetVariant retrieves a product variant
func (r *PostgresProductRepository) GetVariant(ctx context.Context, productID, variantID uuid.UUID) (*domain.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, name, attributes, price_adjustment,
			   cost_price, is_default, status
		FROM product_variants
		WHERE product_id = $1 AND id = $2 AND status = 'active'
	`
	
	var v domain.ProductVariant
	var attributes []byte
	var costPrice *string
	
	err := r.pool.QueryRow(ctx, query, productID, variantID).Scan(
		&v.ID, &v.ProductID, &v.SKU, &v.Name, &attributes,
		&v.PriceAdjustment, &costPrice, &v.IsDefault, &v.Status,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("variant not found: %s", variantID)
		}
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}
	
	if costPrice != nil {
		v.CostPrice, _ = decimal.NewFromString(*costPrice)
	}
	
	return &v, nil
}

// GetVariants retrieves multiple variants
func (r *PostgresProductRepository) GetVariants(ctx context.Context, productID uuid.UUID, variantIDs []uuid.UUID) (map[uuid.UUID]*domain.ProductVariant, error) {
	if len(variantIDs) == 0 {
		return make(map[uuid.UUID]*domain.ProductVariant), nil
	}
	
	query := `
		SELECT id, product_id, sku, name, attributes, price_adjustment,
			   cost_price, is_default, status
		FROM product_variants
		WHERE product_id = $1 AND id = ANY($2) AND status = 'active'
	`
	
	rows, err := r.pool.Query(ctx, query, productID, variantIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()
	
	variants := make(map[uuid.UUID]*domain.ProductVariant)
	
	for rows.Next() {
		var v domain.ProductVariant
		var attributes []byte
		var costPrice *string
		
		err := rows.Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.Name, &attributes,
			&v.PriceAdjustment, &costPrice, &v.IsDefault, &v.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}
		
		if costPrice != nil {
			v.CostPrice, _ = decimal.NewFromString(*costPrice)
		}
		
		variants[v.ID] = &v
	}
	
	return variants, nil
}

// PostgresCustomerRepository implements CustomerRepository using PostgreSQL
type PostgresCustomerRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresCustomerRepository creates a new PostgreSQL customer repository
func NewPostgresCustomerRepository(pool *pgxpool.Pool) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{pool: pool}
}

// GetByID retrieves a customer by ID
func (r *PostgresCustomerRepository) GetByID(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, tenant_id, type, tier, segment, business_name, territory_id,
			   price_list_id, discount_tier, credit_limit, credit_used,
			   payment_terms, lifetime_value, order_count, status
		FROM customers
		WHERE tenant_id = $1 AND id = $2 AND status = 'active'
	`
	
	var c domain.Customer
	var territoryID, priceListID *string
	
	err := r.pool.QueryRow(ctx, query, tenantID, customerID).Scan(
		&c.ID, &c.TenantID, &c.Type, &c.Tier, &c.Segment, &c.BusinessName,
		&territoryID, &priceListID, &c.DiscountTier, &c.CreditLimit,
		&c.CreditUsed, &c.PaymentTerms, &c.LifetimeValue, &c.OrderCount, &c.Status,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("customer not found: %s", customerID)
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	if territoryID != nil {
		c.TerritoryID, _ = uuid.Parse(*territoryID)
	}
	if priceListID != nil {
		c.PriceListID, _ = uuid.Parse(*priceListID)
	}
	
	return &c, nil
}

// PostgresPriceListRepository implements PriceListRepository using PostgreSQL
type PostgresPriceListRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresPriceListRepository creates a new PostgreSQL price list repository
func NewPostgresPriceListRepository(pool *pgxpool.Pool) *PostgresPriceListRepository {
	return &PostgresPriceListRepository{pool: pool}
}

// GetApplicable retrieves all applicable price lists for a customer
func (r *PostgresPriceListRepository) GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, now time.Time) ([]*domain.PriceList, error) {
	query := `
		SELECT id, tenant_id, name, code, customer_types, territories,
			   customer_segments, currency, valid_from, valid_to, priority, is_active
		FROM price_lists
		WHERE tenant_id = $1 
		  AND is_active = true
		  AND (valid_from IS NULL OR valid_from <= $2)
		  AND (valid_to IS NULL OR valid_to >= $2)
		  AND (
			  customer_types = '{}' 
			  OR $3::varchar = ANY(customer_types)
		  )
		  AND (
			  territories = '{}' 
			  OR $4::uuid = ANY(territories)
		  )
		ORDER BY priority DESC
	`
	
	rows, err := r.pool.Query(ctx, query, tenantID, now, string(customer.Type), customer.TerritoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query price lists: %w", err)
	}
	defer rows.Close()
	
	var priceLists []*domain.PriceList
	
	for rows.Next() {
		var pl domain.PriceList
		var customerTypes, territories, segments []string
		var validFrom, validTo *time.Time
		
		err := rows.Scan(
			&pl.ID, &pl.TenantID, &pl.Name, &pl.Code, &customerTypes,
			&territories, &segments, &pl.Currency, &validFrom, &validTo,
			&pl.Priority, &pl.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price list: %w", err)
		}
		
		// Convert string arrays to typed arrays
		for _, ct := range customerTypes {
			pl.CustomerTypes = append(pl.CustomerTypes, domain.CustomerType(ct))
		}
		for _, t := range territories {
			if tid, err := uuid.Parse(t); err == nil {
				pl.TerritoryIDs = append(pl.TerritoryIDs, tid)
			}
		}
		pl.CustomerSegments = segments
		
		if validFrom != nil {
			pl.ValidFrom = *validFrom
		}
		if validTo != nil {
			pl.ValidTo = *validTo
		}
		
		priceLists = append(priceLists, &pl)
	}
	
	// Also add customer's specific price list if assigned
	if customer.PriceListID != uuid.Nil {
		found := false
		for _, pl := range priceLists {
			if pl.ID == customer.PriceListID {
				found = true
				break
			}
		}
		
		if !found {
			pl, err := r.getByID(ctx, customer.PriceListID)
			if err == nil && pl != nil {
				priceLists = append([]*domain.PriceList{pl}, priceLists...)
			}
		}
	}
	
	return priceLists, nil
}

// getByID retrieves a single price list by ID
func (r *PostgresPriceListRepository) getByID(ctx context.Context, id uuid.UUID) (*domain.PriceList, error) {
	query := `
		SELECT id, tenant_id, name, code, customer_types, territories,
			   customer_segments, currency, valid_from, valid_to, priority, is_active
		FROM price_lists
		WHERE id = $1 AND is_active = true
	`
	
	var pl domain.PriceList
	var customerTypes, territories, segments []string
	var validFrom, validTo *time.Time
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pl.ID, &pl.TenantID, &pl.Name, &pl.Code, &customerTypes,
		&territories, &segments, &pl.Currency, &validFrom, &validTo,
		&pl.Priority, &pl.IsActive,
	)
	
	if err != nil {
		return nil, err
	}
	
	for _, ct := range customerTypes {
		pl.CustomerTypes = append(pl.CustomerTypes, domain.CustomerType(ct))
	}
	for _, t := range territories {
		if tid, err := uuid.Parse(t); err == nil {
			pl.TerritoryIDs = append(pl.TerritoryIDs, tid)
		}
	}
	pl.CustomerSegments = segments
	
	if validFrom != nil {
		pl.ValidFrom = *validFrom
	}
	if validTo != nil {
		pl.ValidTo = *validTo
	}
	
	return &pl, nil
}

// GetItems retrieves price list items for given price lists and products
func (r *PostgresPriceListRepository) GetItems(ctx context.Context, priceListIDs []uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID][]*domain.PriceListItem, error) {
	if len(priceListIDs) == 0 || len(productIDs) == 0 {
		return make(map[uuid.UUID][]*domain.PriceListItem), nil
	}
	
	query := `
		SELECT id, price_list_id, product_id, variant_id, pricing_method,
			   price, discount_percent, discount_amount, margin_percent,
			   min_quantity, max_quantity, valid_from, valid_to
		FROM price_list_items
		WHERE price_list_id = ANY($1) AND product_id = ANY($2)
		ORDER BY price_list_id, min_quantity
	`
	
	rows, err := r.pool.Query(ctx, query, priceListIDs, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query price list items: %w", err)
	}
	defer rows.Close()
	
	items := make(map[uuid.UUID][]*domain.PriceListItem)
	
	for rows.Next() {
		var pli domain.PriceListItem
		var variantID *string
		var price, discountPercent, discountAmount, marginPercent *string
		var maxQty *int
		var validFrom, validTo *time.Time
		
		err := rows.Scan(
			&pli.ID, &pli.PriceListID, &pli.ProductID, &variantID, &pli.PricingMethod,
			&price, &discountPercent, &discountAmount, &marginPercent,
			&pli.MinQuantity, &maxQty, &validFrom, &validTo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price list item: %w", err)
		}
		
		if variantID != nil {
			pli.VariantID, _ = uuid.Parse(*variantID)
		}
		if price != nil {
			pli.Price, _ = decimal.NewFromString(*price)
		}
		if discountPercent != nil {
			pli.DiscountPercent, _ = decimal.NewFromString(*discountPercent)
		}
		if discountAmount != nil {
			pli.DiscountAmount, _ = decimal.NewFromString(*discountAmount)
		}
		if marginPercent != nil {
			pli.MarginPercent, _ = decimal.NewFromString(*marginPercent)
		}
		if maxQty != nil {
			pli.MaxQuantity = *maxQty
		}
		if validFrom != nil {
			pli.ValidFrom = *validFrom
		}
		if validTo != nil {
			pli.ValidTo = *validTo
		}
		
		items[pli.PriceListID] = append(items[pli.PriceListID], &pli)
	}
	
	return items, nil
}

// PostgresContractPriceRepository implements ContractPriceRepository using PostgreSQL
type PostgresContractPriceRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresContractPriceRepository creates a new PostgreSQL contract price repository
func NewPostgresContractPriceRepository(pool *pgxpool.Pool) *PostgresContractPriceRepository {
	return &PostgresContractPriceRepository{pool: pool}
}

// GetForCustomer retrieves contract prices for a customer
func (r *PostgresContractPriceRepository) GetForCustomer(ctx context.Context, customerID uuid.UUID, productIDs []uuid.UUID, now time.Time) (map[uuid.UUID]*domain.ContractPrice, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID]*domain.ContractPrice), nil
	}
	
	query := `
		SELECT id, tenant_id, customer_id, product_id, variant_id, price,
			   min_quantity, valid_from, valid_to, contract_reference, is_active
		FROM contract_prices
		WHERE customer_id = $1 
		  AND product_id = ANY($2)
		  AND is_active = true
		  AND valid_from <= $3
		  AND valid_to >= $3
		ORDER BY product_id, min_quantity DESC
	`
	
	rows, err := r.pool.Query(ctx, query, customerID, productIDs, now)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract prices: %w", err)
	}
	defer rows.Close()
	
	contracts := make(map[uuid.UUID]*domain.ContractPrice)
	
	for rows.Next() {
		var cp domain.ContractPrice
		var variantID *string
		
		err := rows.Scan(
			&cp.ID, &cp.TenantID, &cp.CustomerID, &cp.ProductID, &variantID,
			&cp.Price, &cp.MinQuantity, &cp.ValidFrom, &cp.ValidTo,
			&cp.ContractReference, &cp.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract price: %w", err)
		}
		
		if variantID != nil {
			cp.VariantID, _ = uuid.Parse(*variantID)
		}
		
		// Keep only the first (highest min_quantity) contract per product
		if _, exists := contracts[cp.ProductID]; !exists {
			contracts[cp.ProductID] = &cp
		}
	}
	
	return contracts, nil
}

// PostgresVolumeDiscountRepository implements VolumeDiscountRepository using PostgreSQL
type PostgresVolumeDiscountRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresVolumeDiscountRepository creates a new PostgreSQL volume discount repository
func NewPostgresVolumeDiscountRepository(pool *pgxpool.Pool) *PostgresVolumeDiscountRepository {
	return &PostgresVolumeDiscountRepository{pool: pool}
}

// GetApplicable retrieves applicable volume discounts
func (r *PostgresVolumeDiscountRepository) GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, productIDs []uuid.UUID, now time.Time) ([]*domain.VolumeDiscount, error) {
	query := `
		SELECT id, tenant_id, name, applies_to, category_ids, product_ids,
			   brand_names, customer_types, customer_tiers, tiers,
			   can_combine, priority, valid_from, valid_to, is_active
		FROM volume_discounts
		WHERE tenant_id = $1
		  AND is_active = true
		  AND (valid_from IS NULL OR valid_from <= $2)
		  AND (valid_to IS NULL OR valid_to >= $2)
		  AND (
			  customer_types = '{}' 
			  OR $3::varchar = ANY(customer_types)
		  )
		  AND (
			  customer_tiers = '{}' 
			  OR $4::varchar = ANY(customer_tiers)
		  )
		ORDER BY priority DESC
	`
	
	rows, err := r.pool.Query(ctx, query, tenantID, now, string(customer.Type), string(customer.Tier))
	if err != nil {
		return nil, fmt.Errorf("failed to query volume discounts: %w", err)
	}
	defer rows.Close()
	
	var discounts []*domain.VolumeDiscount
	
	for rows.Next() {
		var vd domain.VolumeDiscount
		var categoryIDs, productIDsArr, brandNames, customerTypes, customerTiers []string
		var tiersJSON []byte
		var validFrom, validTo *time.Time
		
		err := rows.Scan(
			&vd.ID, &vd.TenantID, &vd.Name, &vd.AppliesTo, &categoryIDs,
			&productIDsArr, &brandNames, &customerTypes, &customerTiers, &tiersJSON,
			&vd.CanCombine, &vd.Priority, &validFrom, &validTo, &vd.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan volume discount: %w", err)
		}
		
		// Parse arrays
		for _, cid := range categoryIDs {
			if id, err := uuid.Parse(cid); err == nil {
				vd.CategoryIDs = append(vd.CategoryIDs, id)
			}
		}
		for _, pid := range productIDsArr {
			if id, err := uuid.Parse(pid); err == nil {
				vd.ProductIDs = append(vd.ProductIDs, id)
			}
		}
		vd.BrandNames = brandNames
		for _, ct := range customerTypes {
			vd.CustomerTypes = append(vd.CustomerTypes, domain.CustomerType(ct))
		}
		for _, tier := range customerTiers {
			vd.CustomerTiers = append(vd.CustomerTiers, domain.CustomerTier(tier))
		}
		
		if validFrom != nil {
			vd.ValidFrom = *validFrom
		}
		if validTo != nil {
			vd.ValidTo = *validTo
		}
		
		// Parse tiers JSON - simplified for now
		// In production, use proper JSON unmarshaling
		
		discounts = append(discounts, &vd)
	}
	
	return discounts, nil
}

// PostgresPromotionRepository implements PromotionRepository using PostgreSQL
type PostgresPromotionRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresPromotionRepository creates a new PostgreSQL promotion repository
func NewPostgresPromotionRepository(pool *pgxpool.Pool) *PostgresPromotionRepository {
	return &PostgresPromotionRepository{pool: pool}
}

// GetApplicable retrieves applicable promotions
func (r *PostgresPromotionRepository) GetApplicable(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, productIDs []uuid.UUID, now time.Time) ([]*domain.Promotion, error) {
	query := `
		SELECT id, tenant_id, name, code, type, applies_to, product_ids,
			   category_ids, customer_types, customer_tiers, discount_type,
			   discount_value, min_order_value, min_quantity, max_usage,
			   usage_count, buy_quantity, get_quantity, can_combine,
			   priority, valid_from, valid_to, is_active
		FROM promotions
		WHERE tenant_id = $1
		  AND is_active = true
		  AND valid_from <= $2
		  AND valid_to >= $2
		  AND (max_usage = 0 OR usage_count < max_usage)
		  AND (
			  customer_types = '{}' 
			  OR $3::varchar = ANY(customer_types)
		  )
		ORDER BY priority DESC
	`
	
	rows, err := r.pool.Query(ctx, query, tenantID, now, string(customer.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to query promotions: %w", err)
	}
	defer rows.Close()
	
	var promotions []*domain.Promotion
	
	for rows.Next() {
		var p domain.Promotion
		var productIDsArr, categoryIDsArr, customerTypes, customerTiers []string
		var minOrderValue *string
		
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.Name, &p.Code, &p.Type, &p.AppliesTo,
			&productIDsArr, &categoryIDsArr, &customerTypes, &customerTiers,
			&p.DiscountType, &p.DiscountValue, &minOrderValue, &p.MinQuantity,
			&p.MaxUsage, &p.UsageCount, &p.BuyQuantity, &p.GetQuantity,
			&p.CanCombine, &p.Priority, &p.ValidFrom, &p.ValidTo, &p.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan promotion: %w", err)
		}
		
		// Parse arrays
		for _, pid := range productIDsArr {
			if id, err := uuid.Parse(pid); err == nil {
				p.ProductIDs = append(p.ProductIDs, id)
			}
		}
		for _, cid := range categoryIDsArr {
			if id, err := uuid.Parse(cid); err == nil {
				p.CategoryIDs = append(p.CategoryIDs, id)
			}
		}
		for _, ct := range customerTypes {
			p.CustomerTypes = append(p.CustomerTypes, domain.CustomerType(ct))
		}
		for _, tier := range customerTiers {
			p.CustomerTiers = append(p.CustomerTiers, domain.CustomerTier(tier))
		}
		
		if minOrderValue != nil {
			p.MinOrderValue, _ = decimal.NewFromString(*minOrderValue)
		}
		
		promotions = append(promotions, &p)
	}
	
	return promotions, nil
}

// PostgresTaxRepository implements TaxRepository using PostgreSQL
type PostgresTaxRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresTaxRepository creates a new PostgreSQL tax repository
func NewPostgresTaxRepository(pool *pgxpool.Pool) *PostgresTaxRepository {
	return &PostgresTaxRepository{pool: pool}
}

// GetRates retrieves tax rates for given categories and location
func (r *PostgresTaxRepository) GetRates(ctx context.Context, tenantID uuid.UUID, categories []string, country, state string) (map[string]*domain.TaxRate, error) {
	query := `
		SELECT id, tenant_id, name, code, rate, category, country, state,
			   is_compound, is_active
		FROM tax_rates
		WHERE tenant_id = $1
		  AND is_active = true
		  AND category = ANY($2)
		  AND country = $3
		  AND (state = '' OR state = $4)
		ORDER BY state DESC NULLS LAST
	`
	
	rows, err := r.pool.Query(ctx, query, tenantID, categories, country, state)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax rates: %w", err)
	}
	defer rows.Close()
	
	rates := make(map[string]*domain.TaxRate)
	
	for rows.Next() {
		var tr domain.TaxRate
		
		err := rows.Scan(
			&tr.ID, &tr.TenantID, &tr.Name, &tr.Code, &tr.Rate,
			&tr.Category, &tr.Country, &tr.State, &tr.IsCompound, &tr.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax rate: %w", err)
		}
		
		// Keep only the most specific rate per category
		if _, exists := rates[tr.Category]; !exists {
			rates[tr.Category] = &tr
		}
	}
	
	return rates, nil
}
