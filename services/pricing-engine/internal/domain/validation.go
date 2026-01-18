// Package domain contains validation methods for pricing engine domain models.
package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrProductIDRequired     = errors.New("product ID is required")
	ErrTenantIDRequired      = errors.New("tenant ID is required")
	ErrProductSKURequired    = errors.New("product SKU is required")
	ErrProductNameRequired   = errors.New("product name is required")
	ErrInvalidBasePrice      = errors.New("base price must be positive")
	ErrInvalidQuantity       = errors.New("quantity must be positive")
	ErrCustomerIDRequired    = errors.New("customer ID is required")
	ErrPriceListNameRequired = errors.New("price list name is required")
	ErrInvalidDateRange      = errors.New("valid_from must be before valid_to")
	ErrEmptyPriceRequest     = errors.New("price request must have at least one item")
)

// Validate validates a Product
func (p *Product) Validate() error {
	if p.ID == uuid.Nil {
		return ErrProductIDRequired
	}
	if p.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if p.SKU == "" {
		return ErrProductSKURequired
	}
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if p.BasePrice.IsNegative() {
		return ErrInvalidBasePrice
	}
	return nil
}

// Validate validates a ProductVariant
func (v *ProductVariant) Validate() error {
	if v.ID == uuid.Nil {
		return ErrProductIDRequired
	}
	if v.ProductID == uuid.Nil {
		return ErrProductIDRequired
	}
	if v.SKU == "" {
		return ErrProductSKURequired
	}
	return nil
}

// Validate validates a Customer
func (c *Customer) Validate() error {
	if c.ID == uuid.Nil {
		return ErrCustomerIDRequired
	}
	if c.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	return nil
}

// Validate validates a PriceList
func (pl *PriceList) Validate() error {
	if pl.ID == uuid.Nil {
		return ErrProductIDRequired
	}
	if pl.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if pl.Name == "" {
		return ErrPriceListNameRequired
	}
	if !pl.ValidFrom.IsZero() && !pl.ValidTo.IsZero() && pl.ValidFrom.After(pl.ValidTo) {
		return ErrInvalidDateRange
	}
	return nil
}

// Validate validates a PriceListItem
func (pli *PriceListItem) Validate() error {
	if pli.ID == uuid.Nil {
		return ErrProductIDRequired
	}
	if pli.PriceListID == uuid.Nil {
		return ErrProductIDRequired
	}
	if pli.ProductID == uuid.Nil {
		return ErrProductIDRequired
	}
	return nil
}

// Validate validates a PriceRequest
func (pr *PriceRequest) Validate() error {
	if pr.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if pr.CustomerID == uuid.Nil {
		return ErrCustomerIDRequired
	}
	if len(pr.Items) == 0 {
		return ErrEmptyPriceRequest
	}
	for _, item := range pr.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates a PriceRequestItem
func (pri *PriceRequestItem) Validate() error {
	if pri.ProductID == uuid.Nil {
		return ErrProductIDRequired
	}
	if pri.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	return nil
}
