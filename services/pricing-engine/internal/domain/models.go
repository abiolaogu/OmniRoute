// Package domain contains all the core domain models for the pricing engine
package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CustomerType represents the classification of a customer
type CustomerType string

const (
	CustomerTypeConsumer    CustomerType = "consumer"
	CustomerTypeRetailer    CustomerType = "retailer"
	CustomerTypeWholesaler  CustomerType = "wholesaler"
	CustomerTypeDistributor CustomerType = "distributor"
	CustomerTypeEnterprise  CustomerType = "enterprise"
)

// CustomerTier represents the loyalty/value tier of a customer
type CustomerTier string

const (
	CustomerTierStandard CustomerTier = "standard"
	CustomerTierSilver   CustomerTier = "silver"
	CustomerTierGold     CustomerTier = "gold"
	CustomerTierPlatinum CustomerTier = "platinum"
)

// PricingMethod defines how a price is calculated
type PricingMethod string

const (
	PricingMethodFixed           PricingMethod = "fixed"
	PricingMethodDiscountPercent PricingMethod = "discount_percent"
	PricingMethodDiscountAmount  PricingMethod = "discount_amount"
	PricingMethodMargin          PricingMethod = "margin"
	PricingMethodCostPlus        PricingMethod = "cost_plus"
)

// PriceSource indicates where a calculated price came from
type PriceSource string

const (
	PriceSourceBase      PriceSource = "base"
	PriceSourcePriceList PriceSource = "price_list"
	PriceSourceContract  PriceSource = "contract"
	PriceSourcePromotion PriceSource = "promotion"
	PriceSourceVolume    PriceSource = "volume_discount"
	PriceSourceCombined  PriceSource = "combined"
)

// Priority returns the priority of the price source (higher = more important)
func (ps PriceSource) Priority() int {
	switch ps {
	case PriceSourceContract:
		return 100
	case PriceSourcePromotion:
		return 70
	case PriceSourceVolume:
		return 60
	case PriceSourcePriceList:
		return 50
	case PriceSourceBase:
		return 0
	default:
		return 0
	}
}

// Product represents a product in the catalog
type Product struct {
	ID            uuid.UUID        `json:"id"`
	TenantID      uuid.UUID        `json:"tenant_id"`
	SKU           string           `json:"sku"`
	Name          string           `json:"name"`
	CategoryID    uuid.UUID        `json:"category_id"`
	Brand         string           `json:"brand"`
	BasePrice     decimal.Decimal  `json:"base_price"`
	CostPrice     decimal.Decimal  `json:"cost_price"`
	Currency      string           `json:"currency"`
	TaxCategory   string           `json:"tax_category"`
	MinOrderQty   int              `json:"min_order_quantity"`
	MaxOrderQty   int              `json:"max_order_quantity"`
	OrderMultiple int              `json:"order_multiple"`
	Status        string           `json:"status"`
	Visibility    string           `json:"visibility"`
	Variants      []ProductVariant `json:"variants,omitempty"`
}

// HasVariants returns true if the product has variants
func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

// ProductVariant represents a variant of a product
type ProductVariant struct {
	ID              uuid.UUID         `json:"id"`
	ProductID       uuid.UUID         `json:"product_id"`
	SKU             string            `json:"sku"`
	Name            string            `json:"name"`
	Attributes      map[string]string `json:"attributes"`
	PriceAdjustment decimal.Decimal   `json:"price_adjustment"`
	CostPrice       decimal.Decimal   `json:"cost_price"`
	IsDefault       bool              `json:"is_default"`
	Status          string            `json:"status"`
}

// Customer represents a buyer in the system
type Customer struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	Type          CustomerType    `json:"type"`
	Tier          CustomerTier    `json:"tier"`
	Segment       string          `json:"segment"`
	BusinessName  string          `json:"business_name"`
	TerritoryID   uuid.UUID       `json:"territory_id"`
	PriceListID   uuid.UUID       `json:"price_list_id"`
	DiscountTier  string          `json:"discount_tier"`
	CreditLimit   decimal.Decimal `json:"credit_limit"`
	CreditUsed    decimal.Decimal `json:"credit_used"`
	PaymentTerms  int             `json:"payment_terms"`
	LifetimeValue decimal.Decimal `json:"lifetime_value"`
	OrderCount    int             `json:"order_count"`
	Status        string          `json:"status"`
}

// PriceList represents a collection of prices for specific customer segments
type PriceList struct {
	ID               uuid.UUID      `json:"id"`
	TenantID         uuid.UUID      `json:"tenant_id"`
	Name             string         `json:"name"`
	Code             string         `json:"code"`
	CustomerTypes    []CustomerType `json:"customer_types"`
	TerritoryIDs     []uuid.UUID    `json:"territory_ids"`
	CustomerSegments []string       `json:"customer_segments"`
	Currency         string         `json:"currency"`
	ValidFrom        time.Time      `json:"valid_from"`
	ValidTo          time.Time      `json:"valid_to"`
	Priority         int            `json:"priority"`
	IsActive         bool           `json:"is_active"`
}

// PriceListItem represents a specific price for a product in a price list
type PriceListItem struct {
	ID              uuid.UUID       `json:"id"`
	PriceListID     uuid.UUID       `json:"price_list_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	VariantID       uuid.UUID       `json:"variant_id"`
	PricingMethod   PricingMethod   `json:"pricing_method"`
	Price           decimal.Decimal `json:"price"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
	DiscountAmount  decimal.Decimal `json:"discount_amount"`
	MarginPercent   decimal.Decimal `json:"margin_percent"`
	MinQuantity     int             `json:"min_quantity"`
	MaxQuantity     int             `json:"max_quantity"`
	ValidFrom       time.Time       `json:"valid_from"`
	ValidTo         time.Time       `json:"valid_to"`
}

// VolumeDiscount represents quantity-based discount rules
type VolumeDiscount struct {
	ID            uuid.UUID            `json:"id"`
	TenantID      uuid.UUID            `json:"tenant_id"`
	Name          string               `json:"name"`
	AppliesTo     string               `json:"applies_to"`
	CategoryIDs   []uuid.UUID          `json:"category_ids"`
	ProductIDs    []uuid.UUID          `json:"product_ids"`
	BrandNames    []string             `json:"brand_names"`
	CustomerTypes []CustomerType       `json:"customer_types"`
	CustomerTiers []CustomerTier       `json:"customer_tiers"`
	Tiers         []VolumeDiscountTier `json:"tiers"`
	CanCombine    bool                 `json:"can_combine"`
	Priority      int                  `json:"priority"`
	ValidFrom     time.Time            `json:"valid_from"`
	ValidTo       time.Time            `json:"valid_to"`
	IsActive      bool                 `json:"is_active"`
}

// FindTierForQuantity finds the applicable tier for a given quantity
func (vd *VolumeDiscount) FindTierForQuantity(quantity int) *VolumeDiscountTier {
	for i := range vd.Tiers {
		tier := &vd.Tiers[i]
		if quantity >= tier.MinQuantity {
			if tier.MaxQuantity == 0 || quantity <= tier.MaxQuantity {
				return tier
			}
		}
	}
	return nil
}

// VolumeDiscountTier represents a single tier in a volume discount
type VolumeDiscountTier struct {
	MinQuantity     int             `json:"min_quantity"`
	MaxQuantity     int             `json:"max_quantity"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
	DiscountAmount  decimal.Decimal `json:"discount_amount"`
}

// ContractPrice represents customer-specific negotiated pricing
type ContractPrice struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	CustomerID        uuid.UUID       `json:"customer_id"`
	ProductID         uuid.UUID       `json:"product_id"`
	VariantID         uuid.UUID       `json:"variant_id"`
	Price             decimal.Decimal `json:"price"`
	MinQuantity       int             `json:"min_quantity"`
	ValidFrom         time.Time       `json:"valid_from"`
	ValidTo           time.Time       `json:"valid_to"`
	ContractReference string          `json:"contract_reference"`
	IsActive          bool            `json:"is_active"`
}

// IsValidAt checks if the contract price is valid at the given time
func (cp *ContractPrice) IsValidAt(t time.Time) bool {
	if !cp.IsActive {
		return false
	}
	return !t.Before(cp.ValidFrom) && !t.After(cp.ValidTo)
}

// Promotion represents a promotional pricing rule
type Promotion struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	Name          string          `json:"name"`
	Code          string          `json:"code"`
	Type          PromotionType   `json:"type"`
	AppliesTo     string          `json:"applies_to"`
	ProductIDs    []uuid.UUID     `json:"product_ids"`
	CategoryIDs   []uuid.UUID     `json:"category_ids"`
	CustomerTypes []CustomerType  `json:"customer_types"`
	CustomerTiers []CustomerTier  `json:"customer_tiers"`
	DiscountType  DiscountType    `json:"discount_type"`
	DiscountValue decimal.Decimal `json:"discount_value"`
	MinOrderValue decimal.Decimal `json:"min_order_value"`
	MinQuantity   int             `json:"min_quantity"`
	MaxUsage      int             `json:"max_usage"`
	UsageCount    int             `json:"usage_count"`
	BuyQuantity   int             `json:"buy_quantity"`
	GetQuantity   int             `json:"get_quantity"`
	CanCombine    bool            `json:"can_combine"`
	Priority      int             `json:"priority"`
	ValidFrom     time.Time       `json:"valid_from"`
	ValidTo       time.Time       `json:"valid_to"`
	IsActive      bool            `json:"is_active"`
}

// IsActiveAt checks if the promotion is active at the given time
func (p *Promotion) IsActiveAt(t time.Time) bool {
	if !p.IsActive {
		return false
	}
	return !t.Before(p.ValidFrom) && !t.After(p.ValidTo)
}

// PromotionType defines the kind of promotion
type PromotionType string

const (
	PromotionTypeDiscount     PromotionType = "discount"
	PromotionTypeBuyXGetY     PromotionType = "buy_x_get_y"
	PromotionTypeBundle       PromotionType = "bundle"
	PromotionTypeFreeShipping PromotionType = "free_shipping"
)

// DiscountType defines how a discount is applied
type DiscountType string

const (
	DiscountTypePercent DiscountType = "percent"
	DiscountTypeFixed   DiscountType = "fixed"
)

// PriceRequest represents a request to calculate prices
type PriceRequest struct {
	TenantID   uuid.UUID          `json:"tenant_id"`
	CustomerID uuid.UUID          `json:"customer_id"`
	Items      []PriceRequestItem `json:"items"`
	Currency   string             `json:"currency"`
	Channel    string             `json:"channel"`
	Timestamp  time.Time          `json:"timestamp"`
}

// Validate validates the price request
func (r *PriceRequest) Validate() error {
	if len(r.Items) == 0 {
		return errors.New("items cannot be empty")
	}
	for _, item := range r.Items {
		if item.Quantity <= 0 {
			return errors.New("item quantity must be positive")
		}
	}
	return nil
}

// PriceRequestItem represents a single item in a price request
type PriceRequestItem struct {
	ProductID uuid.UUID `json:"product_id"`
	VariantID uuid.UUID `json:"variant_id"`
	Quantity  int       `json:"quantity"`
}

// PriceResponse represents the calculated prices
type PriceResponse struct {
	TenantID          uuid.UUID           `json:"tenant_id"`
	CustomerID        uuid.UUID           `json:"customer_id"`
	Items             []PriceResponseItem `json:"items"`
	SubTotal          decimal.Decimal     `json:"subtotal"`
	TotalDiscount     decimal.Decimal     `json:"total_discount"`
	TaxTotal          decimal.Decimal     `json:"tax_total"`
	GrandTotal        decimal.Decimal     `json:"grand_total"`
	Currency          string              `json:"currency"`
	AppliedPromotions []AppliedPromotion  `json:"applied_promotions"`
	CalculatedAt      time.Time           `json:"calculated_at"`
}

// CalculateTotals calculates the subtotal and totals from items
func (r *PriceResponse) CalculateTotals() {
	r.SubTotal = decimal.Zero
	r.TotalDiscount = decimal.Zero
	r.TaxTotal = decimal.Zero

	for _, item := range r.Items {
		r.SubTotal = r.SubTotal.Add(item.LineTotal)
		r.TotalDiscount = r.TotalDiscount.Add(item.DiscountAmount.Mul(decimal.NewFromInt(int64(item.Quantity))))
		r.TaxTotal = r.TaxTotal.Add(item.TaxAmount)
	}

	r.GrandTotal = r.SubTotal.Add(r.TaxTotal)
}

// PriceResponseItem represents the calculated price for a single item
type PriceResponseItem struct {
	ProductID       uuid.UUID        `json:"product_id"`
	VariantID       uuid.UUID        `json:"variant_id"`
	SKU             string           `json:"sku"`
	Name            string           `json:"name"`
	Quantity        int              `json:"quantity"`
	BasePrice       decimal.Decimal  `json:"base_price"`
	UnitPrice       decimal.Decimal  `json:"unit_price"`
	OriginalPrice   decimal.Decimal  `json:"original_price"`
	DiscountAmount  decimal.Decimal  `json:"discount_amount"`
	DiscountPercent decimal.Decimal  `json:"discount_percent"`
	TaxAmount       decimal.Decimal  `json:"tax_amount"`
	LineTotal       decimal.Decimal  `json:"line_total"`
	PriceSource     PriceSource      `json:"price_source"`
	PriceSourceID   uuid.UUID        `json:"price_source_id"`
	PriceBreakdown  []PriceComponent `json:"price_breakdown"`
}

// PriceComponent represents a single component of the final price
type PriceComponent struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	SourceID   uuid.UUID       `json:"source_id"`
	Amount     decimal.Decimal `json:"amount"`
	IsDiscount bool            `json:"is_discount"`
	Priority   int             `json:"priority"`
}

// AppliedPromotion represents a promotion that was applied to the order
type AppliedPromotion struct {
	PromotionID    uuid.UUID       `json:"promotion_id"`
	Name           string          `json:"name"`
	Code           string          `json:"code"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	AppliedToItems []uuid.UUID     `json:"applied_to_items"`
}

// TaxRate represents a tax rate configuration
type TaxRate struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	Name       string          `json:"name"`
	Code       string          `json:"code"`
	Rate       decimal.Decimal `json:"rate"`
	Category   string          `json:"category"`
	Country    string          `json:"country"`
	State      string          `json:"state"`
	IsCompound bool            `json:"is_compound"`
	IsActive   bool            `json:"is_active"`
}

// CalculateTax calculates the tax amount for a given amount
func (tr *TaxRate) CalculateTax(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(tr.Rate).Div(decimal.NewFromInt(100)).Round(2)
}

// ============================================================================
// Haversine distance calculation for location-based pricing
// ============================================================================

// Location represents a geographical point
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DistanceTo calculates the distance in kilometers to another location
func (l Location) DistanceTo(other Location) float64 {
	const earthRadius = 6371.0 // km

	lat1 := l.Latitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	dLat := (other.Latitude - l.Latitude) * math.Pi / 180
	dLon := (other.Longitude - l.Longitude) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
