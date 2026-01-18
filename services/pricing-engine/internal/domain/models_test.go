// Package domain_test contains unit tests for the domain models
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/omniroute/pricing-engine/internal/domain"
)

func TestPriceRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.PriceRequest
		wantErr bool
	}{
		{
			name: "valid request with items",
			req: domain.PriceRequest{
				TenantID:   uuid.New(),
				CustomerID: uuid.New(),
				Items: []domain.PriceRequestItem{
					{ProductID: uuid.New(), Quantity: 10},
				},
				Currency:  "NGN",
				Channel:   "web",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty items",
			req: domain.PriceRequest{
				TenantID:   uuid.New(),
				CustomerID: uuid.New(),
				Items:      []domain.PriceRequestItem{},
				Currency:   "NGN",
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			req: domain.PriceRequest{
				TenantID:   uuid.New(),
				CustomerID: uuid.New(),
				Items: []domain.PriceRequestItem{
					{ProductID: uuid.New(), Quantity: 0},
				},
				Currency: "NGN",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PriceRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomerType_String(t *testing.T) {
	tests := []struct {
		ct   domain.CustomerType
		want string
	}{
		{domain.CustomerTypeConsumer, "consumer"},
		{domain.CustomerTypeRetailer, "retailer"},
		{domain.CustomerTypeWholesaler, "wholesaler"},
		{domain.CustomerTypeDistributor, "distributor"},
		{domain.CustomerTypeEnterprise, "enterprise"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := string(tt.ct); got != tt.want {
				t.Errorf("CustomerType = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceSource_Priority(t *testing.T) {
	// Contract prices should have highest priority
	sources := []domain.PriceSource{
		domain.PriceSourceBase,
		domain.PriceSourcePriceList,
		domain.PriceSourcePromotion,
		domain.PriceSourceContract,
	}

	expected := []int{0, 50, 70, 100}

	for i, source := range sources {
		if got := source.Priority(); got != expected[i] {
			t.Errorf("PriceSource.Priority() for %v = %v, want %v", source, got, expected[i])
		}
	}
}

func TestPricingMethod_Calculate(t *testing.T) {
	basePrice := decimal.NewFromFloat(100.00)

	tests := []struct {
		name     string
		method   domain.PricingMethod
		discount decimal.Decimal
		expected decimal.Decimal
	}{
		{
			name:     "fixed price",
			method:   domain.PricingMethodFixed,
			discount: decimal.NewFromFloat(80.00),
			expected: decimal.NewFromFloat(80.00),
		},
		{
			name:     "10% discount",
			method:   domain.PricingMethodDiscountPercent,
			discount: decimal.NewFromFloat(10),
			expected: decimal.NewFromFloat(90.00),
		},
		{
			name:     "fixed amount discount",
			method:   domain.PricingMethodDiscountAmount,
			discount: decimal.NewFromFloat(15.00),
			expected: decimal.NewFromFloat(85.00),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result decimal.Decimal
			switch tt.method {
			case domain.PricingMethodFixed:
				result = tt.discount
			case domain.PricingMethodDiscountPercent:
				discount := basePrice.Mul(tt.discount).Div(decimal.NewFromInt(100))
				result = basePrice.Sub(discount)
			case domain.PricingMethodDiscountAmount:
				result = basePrice.Sub(tt.discount)
			}

			if !result.Equal(tt.expected) {
				t.Errorf("PricingMethod calculation = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProduct_HasVariants(t *testing.T) {
	product := &domain.Product{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		SKU:       "TEST-001",
		Name:      "Test Product",
		BasePrice: decimal.NewFromFloat(100.00),
	}

	// Product without variants
	if product.HasVariants() {
		t.Error("Product.HasVariants() should return false when no variants")
	}

	// Add variants
	product.Variants = []domain.ProductVariant{
		{ID: uuid.New(), SKU: "TEST-001-A", Name: "Variant A"},
	}

	if !product.HasVariants() {
		t.Error("Product.HasVariants() should return true when variants exist")
	}
}

func TestPromotion_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		promotion domain.Promotion
		expected  bool
	}{
		{
			name: "active promotion",
			promotion: domain.Promotion{
				IsActive:  true,
				ValidFrom: now.Add(-24 * time.Hour),
				ValidTo:   now.Add(24 * time.Hour),
			},
			expected: true,
		},
		{
			name: "expired promotion",
			promotion: domain.Promotion{
				IsActive:  true,
				ValidFrom: now.Add(-48 * time.Hour),
				ValidTo:   now.Add(-24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "future promotion",
			promotion: domain.Promotion{
				IsActive:  true,
				ValidFrom: now.Add(24 * time.Hour),
				ValidTo:   now.Add(48 * time.Hour),
			},
			expected: false,
		},
		{
			name: "deactivated promotion",
			promotion: domain.Promotion{
				IsActive:  false,
				ValidFrom: now.Add(-24 * time.Hour),
				ValidTo:   now.Add(24 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.promotion.IsActiveAt(now); got != tt.expected {
				t.Errorf("Promotion.IsActiveAt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestVolumeDiscount_FindApplicableTier(t *testing.T) {
	vd := &domain.VolumeDiscount{
		ID:   uuid.New(),
		Name: "Bulk Discount",
		Tiers: []domain.VolumeDiscountTier{
			{MinQuantity: 1, MaxQuantity: 10, DiscountPercent: decimal.NewFromFloat(0)},
			{MinQuantity: 11, MaxQuantity: 50, DiscountPercent: decimal.NewFromFloat(5)},
			{MinQuantity: 51, MaxQuantity: 100, DiscountPercent: decimal.NewFromFloat(10)},
			{MinQuantity: 101, MaxQuantity: 0, DiscountPercent: decimal.NewFromFloat(15)}, // 0 = unlimited
		},
	}

	tests := []struct {
		quantity         int
		expectedDiscount float64
	}{
		{5, 0},
		{15, 5},
		{75, 10},
		{200, 15},
	}

	for _, tt := range tests {
		tier := vd.FindTierForQuantity(tt.quantity)
		if tier == nil {
			t.Errorf("FindTierForQuantity(%d) returned nil", tt.quantity)
			continue
		}

		expected := decimal.NewFromFloat(tt.expectedDiscount)
		if !tier.DiscountPercent.Equal(expected) {
			t.Errorf("FindTierForQuantity(%d) discount = %v, want %v",
				tt.quantity, tier.DiscountPercent, expected)
		}
	}
}

func TestTaxRate_CalculateTax(t *testing.T) {
	taxRate := &domain.TaxRate{
		ID:       uuid.New(),
		Name:     "VAT",
		Category: "standard",
		Rate:     decimal.NewFromFloat(7.5),
		Country:  "NG",
	}

	amount := decimal.NewFromFloat(1000.00)
	expectedTax := decimal.NewFromFloat(75.00)

	tax := taxRate.CalculateTax(amount)

	if !tax.Equal(expectedTax) {
		t.Errorf("TaxRate.CalculateTax(%v) = %v, want %v", amount, tax, expectedTax)
	}
}

func TestPriceResponse_CalculateTotals(t *testing.T) {
	response := &domain.PriceResponse{
		TenantID:   uuid.New(),
		CustomerID: uuid.New(),
		Currency:   "NGN",
		Items: []domain.PriceResponseItem{
			{
				ProductID: uuid.New(),
				Quantity:  10,
				UnitPrice: decimal.NewFromFloat(100.00),
				LineTotal: decimal.NewFromFloat(1000.00),
			},
			{
				ProductID: uuid.New(),
				Quantity:  5,
				UnitPrice: decimal.NewFromFloat(200.00),
				LineTotal: decimal.NewFromFloat(1000.00),
			},
		},
	}

	response.CalculateTotals()

	expectedSubtotal := decimal.NewFromFloat(2000.00)
	if !response.SubTotal.Equal(expectedSubtotal) {
		t.Errorf("PriceResponse.SubTotal = %v, want %v", response.SubTotal, expectedSubtotal)
	}
}

func TestContractPrice_IsValid(t *testing.T) {
	now := time.Now()

	contract := &domain.ContractPrice{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		ProductID:  uuid.New(),
		Price:      decimal.NewFromFloat(85.00),
		ValidFrom:  now.Add(-24 * time.Hour),
		ValidTo:    now.Add(24 * time.Hour),
		IsActive:   true,
	}

	if !contract.IsValidAt(now) {
		t.Error("ContractPrice.IsValidAt() should return true for valid contract")
	}

	// Test expired contract
	contract.ValidTo = now.Add(-1 * time.Hour)
	if contract.IsValidAt(now) {
		t.Error("ContractPrice.IsValidAt() should return false for expired contract")
	}
}
