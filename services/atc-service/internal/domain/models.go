// Package domain contains the core domain models for the Authority to Collect (ATC) service.
// ATC enables distributors, wholesalers, and manufacturers to delegate collection rights
// to downstream partners in the FMCG distribution chain.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Domain errors
var (
	ErrATCNotFound             = errors.New("ATC grant not found")
	ErrATCExpired              = errors.New("ATC grant has expired")
	ErrATCInactive             = errors.New("ATC grant is not active")
	ErrATCLimitExceeded        = errors.New("ATC cumulative limit exceeded")
	ErrATCAmountExceeded       = errors.New("amount exceeds maximum per-transaction limit")
	ErrInvalidGrantor          = errors.New("grantor not authorized to create this ATC")
	ErrInvalidGrantee          = errors.New("grantee not eligible for this ATC")
	ErrDuplicateATC            = errors.New("duplicate ATC grant exists")
	ErrCollectionNotAllowed    = errors.New("collection not allowed by ATC scope")
	ErrSettlementFailed        = errors.New("settlement processing failed")
	ErrInvalidCommissionConfig = errors.New("invalid commission configuration")
)

// ============================================================================
// Value Objects
// ============================================================================

// PartyType represents the type of party in the distribution chain
type PartyType string

const (
	PartyManufacturer PartyType = "MANUFACTURER"
	PartyDistributor  PartyType = "DISTRIBUTOR"
	PartyWholesaler   PartyType = "WHOLESALER"
	PartyRetailer     PartyType = "RETAILER"
	PartyWorker       PartyType = "WORKER"
	PartyPlatform     PartyType = "PLATFORM"
)

// ATCScope represents the scope of collection authority
type ATCScope string

const (
	ScopeAllProducts       ATCScope = "ALL_PRODUCTS"
	ScopeProductCategory   ATCScope = "PRODUCT_CATEGORY"
	ScopeSpecificProducts  ATCScope = "SPECIFIC_PRODUCTS"
	ScopeSpecificCustomers ATCScope = "SPECIFIC_CUSTOMERS"
	ScopeGeographicArea    ATCScope = "GEOGRAPHIC_AREA"
)

// CollectionType represents allowed collection methods
type CollectionType string

const (
	CollectionCash                  CollectionType = "CASH"
	CollectionMobileMoney           CollectionType = "MOBILE_MONEY"
	CollectionBankTransfer          CollectionType = "BANK_TRANSFER"
	CollectionCheque                CollectionType = "CHEQUE"
	CollectionCreditAgainstPurchase CollectionType = "CREDIT_AGAINST_PURCHASE"
	CollectionAllMethods            CollectionType = "ALL_METHODS"
)

// CommissionType represents the type of commission structure
type CommissionType string

const (
	CommissionPercentage       CommissionType = "PERCENTAGE"
	CommissionFlatFee          CommissionType = "FLAT_FEE"
	CommissionTieredPercentage CommissionType = "TIERED_PERCENTAGE"
	CommissionHybrid           CommissionType = "HYBRID"
)

// SettlementFrequency represents how often settlements occur
type SettlementFrequency string

const (
	SettlementInstant  SettlementFrequency = "INSTANT"
	SettlementDaily    SettlementFrequency = "DAILY"
	SettlementWeekly   SettlementFrequency = "WEEKLY"
	SettlementBiweekly SettlementFrequency = "BIWEEKLY"
	SettlementMonthly  SettlementFrequency = "MONTHLY"
	SettlementOnDemand SettlementFrequency = "ON_DEMAND"
)

// ATCStatus represents the status of an ATC grant
type ATCStatus string

const (
	ATCStatusDraft           ATCStatus = "DRAFT"
	ATCStatusPendingApproval ATCStatus = "PENDING_APPROVAL"
	ATCStatusActive          ATCStatus = "ACTIVE"
	ATCStatusSuspended       ATCStatus = "SUSPENDED"
	ATCStatusExpired         ATCStatus = "EXPIRED"
	ATCStatusRevoked         ATCStatus = "REVOKED"
)

// SettlementStatus represents the status of a settlement
type SettlementStatus string

const (
	SettlementPending    SettlementStatus = "PENDING"
	SettlementScheduled  SettlementStatus = "SCHEDULED"
	SettlementProcessing SettlementStatus = "PROCESSING"
	SettlementSettled    SettlementStatus = "SETTLED"
	SettlementFailed     SettlementStatus = "FAILED"
	SettlementDisputed   SettlementStatus = "DISPUTED"
)

// AccountDetails represents bank/mobile money account details
type AccountDetails struct {
	BankCode      string `json:"bank_code,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	AccountName   string `json:"account_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Provider      string `json:"provider,omitempty"`
}

// GeoPoint represents a geographic point
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ============================================================================
// Aggregates
// ============================================================================

// ATCGrant represents an Authority to Collect grant
type ATCGrant struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Reference string    `json:"reference"`

	// Grantor (who gives authority)
	GrantorID   uuid.UUID `json:"grantor_id"`
	GrantorType PartyType `json:"grantor_type"`
	GrantorName string    `json:"grantor_name"`

	// Grantee (who receives authority)
	GranteeID   uuid.UUID `json:"grantee_id"`
	GranteeType PartyType `json:"grantee_type"`
	GranteeName string    `json:"grantee_name"`

	// Scope of authority
	Scope          ATCScope         `json:"scope"`
	CollectionType CollectionType   `json:"collection_type"`
	ScopeConfig    []ATCScopeConfig `json:"scope_config,omitempty"`

	// Financial terms
	Currency            string          `json:"currency"`
	MaxAmount           decimal.Decimal `json:"max_amount"`           // Per-transaction limit
	CumulativeLimit     decimal.Decimal `json:"cumulative_limit"`     // Total collection limit
	CumulativeCollected decimal.Decimal `json:"cumulative_collected"` // Running total

	// Commission/Margin structure
	CommissionType  CommissionType   `json:"commission_type"`
	CommissionRate  decimal.Decimal  `json:"commission_rate"` // Percentage
	CommissionFlat  decimal.Decimal  `json:"commission_flat"` // Fixed amount
	MinCommission   decimal.Decimal  `json:"min_commission"`
	MaxCommission   decimal.Decimal  `json:"max_commission"`
	CommissionTiers []CommissionTier `json:"commission_tiers,omitempty"`

	// Settlement terms
	SettlementFrequency SettlementFrequency `json:"settlement_frequency"`
	SettlementDay       int                 `json:"settlement_day"`        // Day of week/month
	SettlementDelayDays int                 `json:"settlement_delay_days"` // Days after collection
	SettlementAccount   AccountDetails      `json:"settlement_account"`

	// Validity
	Status        ATCStatus  `json:"status"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`

	// Hierarchy
	ParentATCID *uuid.UUID  `json:"parent_atc_id,omitempty"`
	ChildATCs   []uuid.UUID `json:"child_atcs,omitempty"`

	// Metadata
	TermsDocumentURL string     `json:"terms_document_url,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ATCScopeConfig represents scope configuration
type ATCScopeConfig struct {
	ID          uuid.UUID `json:"id"`
	ATCGrantID  uuid.UUID `json:"atc_grant_id"`
	ScopeType   string    `json:"scope_type"` // PRODUCT_ID, PRODUCT_CATEGORY, CUSTOMER_ID, etc.
	ScopeValues []string  `json:"scope_values"`
	CreatedAt   time.Time `json:"created_at"`
}

// CommissionTier represents a tiered commission structure
type CommissionTier struct {
	ID             uuid.UUID        `json:"id"`
	ATCGrantID     uuid.UUID        `json:"atc_grant_id"`
	TierNumber     int              `json:"tier_number"`
	MinAmount      decimal.Decimal  `json:"min_amount"`
	MaxAmount      *decimal.Decimal `json:"max_amount,omitempty"`
	CommissionRate decimal.Decimal  `json:"commission_rate"`
	CreatedAt      time.Time        `json:"created_at"`
}

// ATCCollection represents a collection made under an ATC grant
type ATCCollection struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	ATCGrantID uuid.UUID `json:"atc_grant_id"`

	// Collection details
	Reference         string         `json:"reference"`
	CollectedFromID   uuid.UUID      `json:"collected_from_id"`
	CollectedFromType PartyType      `json:"collected_from_type"`
	CollectedFromName string         `json:"collected_from_name"`
	CollectionMethod  CollectionType `json:"collection_method"`
	PaymentReference  string         `json:"payment_reference,omitempty"`

	// Amounts
	GrossAmount      decimal.Decimal `json:"gross_amount"`
	CommissionAmount decimal.Decimal `json:"commission_amount"`
	NetAmount        decimal.Decimal `json:"net_amount"` // Amount to settle to grantor

	// Related entities
	OrderID   *uuid.UUID `json:"order_id,omitempty"`
	InvoiceID *uuid.UUID `json:"invoice_id,omitempty"`

	// Settlement tracking
	SettlementStatus  SettlementStatus `json:"settlement_status"`
	SettlementBatchID *uuid.UUID       `json:"settlement_batch_id,omitempty"`
	SettledAt         *time.Time       `json:"settled_at,omitempty"`

	// Metadata
	CollectionPoint *GeoPoint  `json:"collection_point,omitempty"`
	CollectedBy     *uuid.UUID `json:"collected_by,omitempty"` // Worker who collected
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ATCSettlementBatch represents a settlement batch
type ATCSettlementBatch struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Reference string    `json:"reference"`

	ATCGrantID            uuid.UUID `json:"atc_grant_id"`
	SettlementPeriodStart time.Time `json:"settlement_period_start"`
	SettlementPeriodEnd   time.Time `json:"settlement_period_end"`

	TotalCollections    int             `json:"total_collections"`
	GrossAmount         decimal.Decimal `json:"gross_amount"`
	TotalCommission     decimal.Decimal `json:"total_commission"`
	NetSettlementAmount decimal.Decimal `json:"net_settlement_amount"`

	Collections []ATCCollection `json:"collections,omitempty"`

	// Payment details
	PaymentTransactionID *uuid.UUID `json:"payment_transaction_id,omitempty"`
	PaymentStatus        string     `json:"payment_status"`

	// Reconciliation
	ReconciliationStatus string          `json:"reconciliation_status"`
	DiscrepancyAmount    decimal.Decimal `json:"discrepancy_amount"`
	DiscrepancyNotes     string          `json:"discrepancy_notes,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	SettledAt *time.Time `json:"settled_at,omitempty"`
}

// ============================================================================
// Aggregate Methods
// ============================================================================

// NewATCGrant creates a new ATC grant
func NewATCGrant(tenantID, grantorID uuid.UUID, grantorType PartyType, grantorName string,
	granteeID uuid.UUID, granteeType PartyType, granteeName string,
	scope ATCScope, collectionType CollectionType, currency string) (*ATCGrant, error) {

	if !isValidGrantorGranteeRelation(grantorType, granteeType) {
		return nil, ErrInvalidGrantor
	}

	now := time.Now()
	return &ATCGrant{
		ID:                  uuid.New(),
		TenantID:            tenantID,
		Reference:           generateATCReference(),
		GrantorID:           grantorID,
		GrantorType:         grantorType,
		GrantorName:         grantorName,
		GranteeID:           granteeID,
		GranteeType:         granteeType,
		GranteeName:         granteeName,
		Scope:               scope,
		CollectionType:      collectionType,
		Currency:            currency,
		CumulativeCollected: decimal.Zero,
		Status:              ATCStatusDraft,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// Validate validates the ATC grant
func (a *ATCGrant) Validate() error {
	if a.GrantorID == uuid.Nil {
		return errors.New("grantor ID is required")
	}
	if a.GranteeID == uuid.Nil {
		return errors.New("grantee ID is required")
	}
	if a.CommissionType == CommissionTieredPercentage && len(a.CommissionTiers) == 0 {
		return ErrInvalidCommissionConfig
	}
	return nil
}

// IsActive checks if the ATC grant is currently active
func (a *ATCGrant) IsActive() bool {
	if a.Status != ATCStatusActive {
		return false
	}
	now := time.Now()
	if now.Before(a.EffectiveFrom) {
		return false
	}
	if a.EffectiveTo != nil && now.After(*a.EffectiveTo) {
		return false
	}
	return true
}

// CanCollect checks if a collection can be made under this ATC
func (a *ATCGrant) CanCollect(amount decimal.Decimal) error {
	if !a.IsActive() {
		return ErrATCInactive
	}
	if !a.MaxAmount.IsZero() && amount.GreaterThan(a.MaxAmount) {
		return ErrATCAmountExceeded
	}
	if !a.CumulativeLimit.IsZero() {
		remaining := a.CumulativeLimit.Sub(a.CumulativeCollected)
		if amount.GreaterThan(remaining) {
			return ErrATCLimitExceeded
		}
	}
	return nil
}

// CalculateCommission calculates the commission for a given amount
func (a *ATCGrant) CalculateCommission(amount decimal.Decimal) decimal.Decimal {
	var commission decimal.Decimal

	switch a.CommissionType {
	case CommissionPercentage:
		commission = amount.Mul(a.CommissionRate).Div(decimal.NewFromInt(100))
	case CommissionFlatFee:
		commission = a.CommissionFlat
	case CommissionTieredPercentage:
		commission = a.calculateTieredCommission(amount)
	case CommissionHybrid:
		percentageCommission := amount.Mul(a.CommissionRate).Div(decimal.NewFromInt(100))
		commission = percentageCommission.Add(a.CommissionFlat)
	}

	// Apply min/max bounds
	if !a.MinCommission.IsZero() && commission.LessThan(a.MinCommission) {
		commission = a.MinCommission
	}
	if !a.MaxCommission.IsZero() && commission.GreaterThan(a.MaxCommission) {
		commission = a.MaxCommission
	}

	return commission
}

func (a *ATCGrant) calculateTieredCommission(amount decimal.Decimal) decimal.Decimal {
	for _, tier := range a.CommissionTiers {
		if amount.GreaterThanOrEqual(tier.MinAmount) {
			if tier.MaxAmount == nil || amount.LessThanOrEqual(*tier.MaxAmount) {
				return amount.Mul(tier.CommissionRate).Div(decimal.NewFromInt(100))
			}
		}
	}
	return decimal.Zero
}

// RecordCollection records a collection and updates cumulative amounts
func (a *ATCGrant) RecordCollection(amount decimal.Decimal) {
	a.CumulativeCollected = a.CumulativeCollected.Add(amount)
	a.UpdatedAt = time.Now()
}

// Activate activates the ATC grant
func (a *ATCGrant) Activate(approvedBy uuid.UUID) error {
	if err := a.Validate(); err != nil {
		return err
	}
	now := time.Now()
	a.Status = ATCStatusActive
	a.ApprovedBy = &approvedBy
	a.ApprovedAt = &now
	a.UpdatedAt = now
	return nil
}

// Suspend suspends the ATC grant
func (a *ATCGrant) Suspend() {
	a.Status = ATCStatusSuspended
	a.UpdatedAt = time.Now()
}

// Revoke revokes the ATC grant
func (a *ATCGrant) Revoke() {
	a.Status = ATCStatusRevoked
	a.UpdatedAt = time.Now()
}

// NewATCCollection creates a new collection under an ATC grant
func NewATCCollection(atc *ATCGrant, collectedFromID uuid.UUID, collectedFromType PartyType,
	collectedFromName string, method CollectionType, grossAmount decimal.Decimal) (*ATCCollection, error) {

	if err := atc.CanCollect(grossAmount); err != nil {
		return nil, err
	}

	commission := atc.CalculateCommission(grossAmount)
	netAmount := grossAmount.Sub(commission)

	return &ATCCollection{
		ID:                uuid.New(),
		TenantID:          atc.TenantID,
		ATCGrantID:        atc.ID,
		Reference:         generateCollectionReference(),
		CollectedFromID:   collectedFromID,
		CollectedFromType: collectedFromType,
		CollectedFromName: collectedFromName,
		CollectionMethod:  method,
		GrossAmount:       grossAmount,
		CommissionAmount:  commission,
		NetAmount:         netAmount,
		SettlementStatus:  SettlementPending,
		CreatedAt:         time.Now(),
	}, nil
}

// Helper functions

func isValidGrantorGranteeRelation(grantor, grantee PartyType) bool {
	// Valid relationships in distribution chain
	validRelations := map[PartyType][]PartyType{
		PartyManufacturer: {PartyDistributor},
		PartyDistributor:  {PartyWholesaler, PartyRetailer, PartyWorker},
		PartyWholesaler:   {PartyRetailer, PartyWorker},
		PartyPlatform:     {PartyDistributor, PartyWholesaler, PartyRetailer, PartyWorker},
	}

	allowed, ok := validRelations[grantor]
	if !ok {
		return false
	}
	for _, g := range allowed {
		if g == grantee {
			return true
		}
	}
	return false
}

func generateATCReference() string {
	return "ATC" + uuid.New().String()[:8]
}

func generateCollectionReference() string {
	return "COL" + uuid.New().String()[:8]
}
