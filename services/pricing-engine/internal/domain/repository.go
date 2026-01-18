// Package domain contains repository interfaces for the Pricing Engine service.
// Following DDD principles, repository interfaces are defined in the domain layer.
package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductRepository defines operations for product persistence
type ProductRepository interface {
	// FindByID retrieves a product by ID
	FindByID(ctx context.Context, tenantID, productID uuid.UUID) (*Product, error)

	// FindBySKU retrieves a product by SKU
	FindBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*Product, error)

	// FindByIDs retrieves multiple products by IDs
	FindByIDs(ctx context.Context, tenantID uuid.UUID, productIDs []uuid.UUID) ([]*Product, error)

	// FindByCategory retrieves products in a category
	FindByCategory(ctx context.Context, tenantID, categoryID uuid.UUID, limit, offset int) ([]*Product, error)

	// Save persists a product
	Save(ctx context.Context, product *Product) error

	// Update updates a product
	Update(ctx context.Context, product *Product) error

	// Delete removes a product
	Delete(ctx context.Context, tenantID, productID uuid.UUID) error
}

// CustomerRepository defines operations for customer persistence
type CustomerRepository interface {
	// FindByID retrieves a customer by ID
	FindByID(ctx context.Context, tenantID, customerID uuid.UUID) (*Customer, error)

	// FindBySegment retrieves customers in a segment
	FindBySegment(ctx context.Context, tenantID uuid.UUID, segment string, limit, offset int) ([]*Customer, error)

	// Save persists a customer
	Save(ctx context.Context, customer *Customer) error

	// Update updates a customer
	Update(ctx context.Context, customer *Customer) error

	// UpdateCreditUsage updates customer credit usage
	UpdateCreditUsage(ctx context.Context, tenantID, customerID uuid.UUID, creditUsed interface{}) error
}

// PriceListRepository defines operations for price list persistence
type PriceListRepository interface {
	// FindByID retrieves a price list by ID
	FindByID(ctx context.Context, tenantID, priceListID uuid.UUID) (*PriceList, error)

	// FindActiveForCustomer retrieves active price lists applicable to a customer
	FindActiveForCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*PriceList, error)

	// FindItemsForProducts retrieves price list items for products
	FindItemsForProducts(ctx context.Context, priceListID uuid.UUID, productIDs []uuid.UUID) ([]*PriceListItem, error)

	// Save persists a price list
	Save(ctx context.Context, priceList *PriceList) error

	// SaveItem persists a price list item
	SaveItem(ctx context.Context, item *PriceListItem) error
}

// ContractRepository defines operations for contract persistence
type ContractRepository interface {
	// FindActiveForCustomer retrieves active contracts for a customer
	FindActiveForCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Contract, error)

	// FindPricingForProducts retrieves contract pricing for products
	FindPricingForProducts(ctx context.Context, contractID uuid.UUID, productIDs []uuid.UUID) ([]*ContractPricing, error)

	// Save persists a contract
	Save(ctx context.Context, contract *Contract) error
}

// PromotionRepository defines operations for promotion persistence
type PromotionRepository interface {
	// FindActiveForProducts retrieves active promotions for products
	FindActiveForProducts(ctx context.Context, tenantID uuid.UUID, productIDs []uuid.UUID) ([]*Promotion, error)

	// FindActiveForCustomer retrieves active promotions for a customer
	FindActiveForCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Promotion, error)

	// Save persists a promotion
	Save(ctx context.Context, promotion *Promotion) error

	// IncrementUsage increments the usage count for a promotion
	IncrementUsage(ctx context.Context, promotionID uuid.UUID) error
}

// VolumeDiscountRepository defines operations for volume discount persistence
type VolumeDiscountRepository interface {
	// FindForProduct retrieves volume discounts for a product
	FindForProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*VolumeDiscount, error)

	// FindForCategory retrieves volume discounts for a category
	FindForCategory(ctx context.Context, tenantID, categoryID uuid.UUID) ([]*VolumeDiscount, error)

	// Save persists a volume discount
	Save(ctx context.Context, discount *VolumeDiscount) error
}

// TaxRuleRepository defines operations for tax rule persistence
type TaxRuleRepository interface {
	// FindForLocation retrieves tax rules for a location
	FindForLocation(ctx context.Context, tenantID uuid.UUID, country, state string) ([]*TaxRule, error)

	// FindForCategory retrieves tax rules for a tax category
	FindForCategory(ctx context.Context, tenantID uuid.UUID, taxCategory string) ([]*TaxRule, error)

	// Save persists a tax rule
	Save(ctx context.Context, rule *TaxRule) error
}

// PriceCalculationLogRepository defines operations for price calculation logging
type PriceCalculationLogRepository interface {
	// Save persists a price calculation log entry
	Save(ctx context.Context, log *PriceCalculationLog) error

	// FindByCorrelationID retrieves logs by correlation ID
	FindByCorrelationID(ctx context.Context, correlationID string) ([]*PriceCalculationLog, error)
}

// Contract represents a customer contract
type Contract struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

// ContractPricing represents contract-specific pricing
type ContractPricing struct {
	ContractID uuid.UUID   `json:"contract_id"`
	ProductID  uuid.UUID   `json:"product_id"`
	Price      interface{} `json:"price"`
}

// Promotion represents a promotional discount
type Promotion struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}

// VolumeDiscount represents quantity-based pricing
type VolumeDiscount struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	MinQty    int       `json:"min_qty"`
	MaxQty    int       `json:"max_qty"`
}

// TaxRule represents a tax calculation rule
type TaxRule struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	TaxCategory string    `json:"tax_category"`
	Rate        float64   `json:"rate"`
}

// PriceCalculationLog represents a price calculation audit log
type PriceCalculationLog struct {
	ID            uuid.UUID `json:"id"`
	CorrelationID string    `json:"correlation_id"`
}
