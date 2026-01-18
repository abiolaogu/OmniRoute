// Package domain contains the core domain models for the Payment Service.
// Following DDD principles with aggregates, entities, and value objects.
// Implements credit scoring, payment orchestration, and settlement engine.
package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Value Objects
// ============================================================================

// PaymentStatus represents the payment lifecycle
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

// PaymentMethod represents the payment channel
type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodMobileMoney  PaymentMethod = "mobile_money"
	PaymentMethodUSSD         PaymentMethod = "ussd"
	PaymentMethodWallet       PaymentMethod = "wallet"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodCredit       PaymentMethod = "credit"
	PaymentMethodPOS          PaymentMethod = "pos"
)

// CustomerType for credit classification
type CustomerType string

const (
	CustomerTypeConsumer    CustomerType = "consumer"
	CustomerTypeRetailer    CustomerType = "retailer"
	CustomerTypeWholesaler  CustomerType = "wholesaler"
	CustomerTypeDistributor CustomerType = "distributor"
	CustomerTypeEnterprise  CustomerType = "enterprise"
)

// CreditTier represents the credit score tier
type CreditTier string

const (
	CreditTierPremium    CreditTier = "premium"    // 800-1000
	CreditTierStandard   CreditTier = "standard"   // 600-799
	CreditTierLimited    CreditTier = "limited"    // 400-599
	CreditTierRestricted CreditTier = "restricted" // 200-399
	CreditTierNoCredit   CreditTier = "no_credit"  // <200
)

// InvoiceStatus represents invoice lifecycle
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusIssued    InvoiceStatus = "issued"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	InvoiceStatusWriteOff  InvoiceStatus = "write_off"
)

// CollectionStatus represents collection task status
type CollectionStatus string

const (
	CollectionPending    CollectionStatus = "pending"
	CollectionAssigned   CollectionStatus = "assigned"
	CollectionInProgress CollectionStatus = "in_progress"
	CollectionCompleted  CollectionStatus = "completed"
	CollectionFailed     CollectionStatus = "failed"
	CollectionEscalated  CollectionStatus = "escalated"
)

// Money represents a monetary value (value object)
type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// Add adds two Money values
func (m Money) Add(other Money) Money {
	return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}
}

// Sub subtracts two Money values
func (m Money) Sub(other Money) Money {
	return Money{Amount: m.Amount.Sub(other.Amount), Currency: m.Currency}
}

// ============================================================================
// Credit Scoring Components (Value Objects)
// ============================================================================

// CreditScoreComponent represents a component of the credit score
type CreditScoreComponent struct {
	Category     string        `json:"category"`
	MaxPoints    int           `json:"max_points"`
	EarnedPoints int           `json:"earned_points"`
	Details      []ScoreDetail `json:"details"`
}

// ScoreDetail provides breakdown of score calculation
type ScoreDetail struct {
	Factor      string `json:"factor"`
	Value       string `json:"value"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

// ============================================================================
// Aggregates and Entities
// ============================================================================

// Payment is the aggregate root for payment management
type Payment struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	OrderID    uuid.UUID `json:"order_id"`

	// Payment details
	Method   PaymentMethod   `json:"method"`
	Status   PaymentStatus   `json:"status"`
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`

	// Provider
	Provider    string `json:"provider"`
	ProviderRef string `json:"provider_ref"`

	// Wallet deduction
	WalletDeducted decimal.Decimal `json:"wallet_deducted"`

	// Fallback tracking
	FallbackUsed     bool   `json:"fallback_used"`
	OriginalProvider string `json:"original_provider,omitempty"`

	// Failure info
	FailureReason string `json:"failure_reason,omitempty"`
	FailureCode   string `json:"failure_code,omitempty"`

	// Timestamps
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreditLimit represents a customer's credit limit
type CreditLimit struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`

	// Limit details
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	UtilizedAmount  decimal.Decimal `json:"utilized_amount"`
	AvailableAmount decimal.Decimal `json:"available_amount"`

	// Terms
	PaymentTermsDays int `json:"payment_terms_days"`

	// Validity
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	ReviewDate time.Time `json:"review_date"`

	// Status
	IsActive     bool   `json:"is_active"`
	IsFrozen     bool   `json:"is_frozen"`
	FreezeReason string `json:"freeze_reason,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreditScore represents a customer's credit assessment
type CreditScore struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`

	// Score breakdown (total 1000 points)
	TotalScore int                    `json:"total_score"`
	Tier       CreditTier             `json:"tier"`
	Components []CreditScoreComponent `json:"components"`

	// Transaction History (350 points)
	TransactionScore int `json:"transaction_score"`

	// Payment Behavior (350 points)
	PaymentScore int `json:"payment_score"`

	// Business Profile (200 points)
	BusinessScore int `json:"business_score"`

	// External Signals (100 points)
	ExternalScore int `json:"external_score"`

	// Timestamps
	CalculatedAt time.Time `json:"calculated_at"`
	ValidUntil   time.Time `json:"valid_until"`
}

// Wallet represents a customer or worker digital wallet
type Wallet struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	OwnerType string    `json:"owner_type"` // customer, worker

	// Balance
	Balance          decimal.Decimal `json:"balance"`
	HeldBalance      decimal.Decimal `json:"held_balance"`
	AvailableBalance decimal.Decimal `json:"available_balance"`
	Currency         string          `json:"currency"`

	// Status
	IsActive bool `json:"is_active"`
	IsFrozen bool `json:"is_frozen"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WalletTransaction represents a wallet movement
type WalletTransaction struct {
	ID           uuid.UUID              `json:"id"`
	WalletID     uuid.UUID              `json:"wallet_id"`
	Type         string                 `json:"type"` // credit, debit, hold, release
	Amount       decimal.Decimal        `json:"amount"`
	BalanceAfter decimal.Decimal        `json:"balance_after"`
	Reference    string                 `json:"reference"`
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// Invoice represents a customer invoice
type Invoice struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	OrderID    uuid.UUID `json:"order_id"`

	// Invoice details
	InvoiceNumber string        `json:"invoice_number"`
	Status        InvoiceStatus `json:"status"`

	// Amounts
	SubTotal    decimal.Decimal `json:"sub_total"`
	TaxAmount   decimal.Decimal `json:"tax_amount"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	PaidAmount  decimal.Decimal `json:"paid_amount"`
	BalanceDue  decimal.Decimal `json:"balance_due"`
	Currency    string          `json:"currency"`

	// Dates
	IssueDate time.Time  `json:"issue_date"`
	DueDate   time.Time  `json:"due_date"`
	PaidDate  *time.Time `json:"paid_date,omitempty"`

	// Items
	Items []InvoiceItem `json:"items"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	ProductID   uuid.UUID       `json:"product_id"`
	SKU         string          `json:"sku"`
	Description string          `json:"description"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	TaxRate     decimal.Decimal `json:"tax_rate"`
	LineTotal   decimal.Decimal `json:"line_total"`
}

// CollectionTask represents a cash collection task
type CollectionTask struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	InvoiceID        uuid.UUID  `json:"invoice_id"`
	CustomerID       uuid.UUID  `json:"customer_id"`
	AssignedWorkerID *uuid.UUID `json:"assigned_worker_id,omitempty"`

	Status CollectionStatus `json:"status"`

	// Collection details
	AmountDue        decimal.Decimal `json:"amount_due"`
	AmountCollected  decimal.Decimal `json:"amount_collected"`
	CollectionMethod string          `json:"collection_method"` // cash, transfer, pos, mobile_money

	// Location
	Location  string  `json:"location"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	// Scheduling
	ScheduledDate time.Time  `json:"scheduled_date"`
	AttemptCount  int        `json:"attempt_count"`
	LastAttemptAt *time.Time `json:"last_attempt_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaymentPlan represents an installment payment plan
type PaymentPlan struct {
	ID         uuid.UUID   `json:"id"`
	TenantID   uuid.UUID   `json:"tenant_id"`
	CustomerID uuid.UUID   `json:"customer_id"`
	InvoiceIDs []uuid.UUID `json:"invoice_ids"`

	// Plan details
	TotalAmount  decimal.Decimal `json:"total_amount"`
	Installments []Installment   `json:"installments"`
	Status       string          `json:"status"` // pending_approval, active, completed, defaulted

	// Approval
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Installment represents a single payment in a plan
type Installment struct {
	Number     int             `json:"number"`
	Amount     decimal.Decimal `json:"amount"`
	DueDate    time.Time       `json:"due_date"`
	PaidAmount decimal.Decimal `json:"paid_amount"`
	PaidAt     *time.Time      `json:"paid_at,omitempty"`
	Status     string          `json:"status"` // pending, paid, overdue
}

// Settlement represents a merchant payout
type Settlement struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	MerchantID uuid.UUID `json:"merchant_id"`

	// Settlement details
	GrossAmount decimal.Decimal `json:"gross_amount"`
	Fees        decimal.Decimal `json:"fees"`
	NetAmount   decimal.Decimal `json:"net_amount"`
	Currency    string          `json:"currency"`

	// Status
	Status string `json:"status"` // pending, processing, completed, failed

	// Bank details
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`

	// Timestamps
	ScheduledFor time.Time  `json:"scheduled_for"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ============================================================================
// Domain Events
// ============================================================================

// PaymentCompletedEvent is raised when a payment succeeds
type PaymentCompletedEvent struct {
	PaymentID uuid.UUID       `json:"payment_id"`
	OrderID   uuid.UUID       `json:"order_id"`
	Amount    decimal.Decimal `json:"amount"`
	Method    PaymentMethod   `json:"method"`
	Timestamp time.Time       `json:"timestamp"`
}

// CreditLimitUpdatedEvent is raised when credit limit changes
type CreditLimitUpdatedEvent struct {
	CustomerID uuid.UUID       `json:"customer_id"`
	OldLimit   decimal.Decimal `json:"old_limit"`
	NewLimit   decimal.Decimal `json:"new_limit"`
	Reason     string          `json:"reason"`
	Timestamp  time.Time       `json:"timestamp"`
}

// InvoiceOverdueEvent is raised when an invoice becomes overdue
type InvoiceOverdueEvent struct {
	InvoiceID   uuid.UUID       `json:"invoice_id"`
	CustomerID  uuid.UUID       `json:"customer_id"`
	Amount      decimal.Decimal `json:"amount"`
	DaysPastDue int             `json:"days_past_due"`
	Timestamp   time.Time       `json:"timestamp"`
}

// ============================================================================
// Business Logic Methods
// ============================================================================

// CreditTierFromScore returns the credit tier for a given score
func CreditTierFromScore(score int) CreditTier {
	switch {
	case score >= 800:
		return CreditTierPremium
	case score >= 600:
		return CreditTierStandard
	case score >= 400:
		return CreditTierLimited
	case score >= 200:
		return CreditTierRestricted
	default:
		return CreditTierNoCredit
	}
}

// CalculateAvailable returns the available credit amount
func (cl *CreditLimit) CalculateAvailable() decimal.Decimal {
	return cl.Amount.Sub(cl.UtilizedAmount)
}

// CanUtilize checks if the credit limit can cover the requested amount
func (cl *CreditLimit) CanUtilize(amount decimal.Decimal) bool {
	if !cl.IsActive || cl.IsFrozen {
		return false
	}
	now := time.Now()
	if now.Before(cl.ValidFrom) || now.After(cl.ValidTo) {
		return false
	}
	return cl.CalculateAvailable().GreaterThanOrEqual(amount)
}

// CalculateTotal calculates total score and tier
func (cs *CreditScore) CalculateTotal() {
	cs.TotalScore = cs.TransactionScore + cs.PaymentScore + cs.BusinessScore + cs.ExternalScore
	cs.Tier = CreditTierFromScore(cs.TotalScore)
}

// CalculateAvailable returns the available wallet balance
func (w *Wallet) CalculateAvailable() decimal.Decimal {
	return w.Balance.Sub(w.HeldBalance)
}

// CanDebit checks if the wallet can cover a debit
func (w *Wallet) CanDebit(amount decimal.Decimal) bool {
	if !w.IsActive || w.IsFrozen {
		return false
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return false
	}
	return w.CalculateAvailable().GreaterThanOrEqual(amount)
}

// CanTransitionTo checks if a payment can transition to the target status
func (p *Payment) CanTransitionTo(target PaymentStatus) bool {
	validTransitions := map[PaymentStatus][]PaymentStatus{
		PaymentStatusPending:    {PaymentStatusProcessing, PaymentStatusCancelled},
		PaymentStatusProcessing: {PaymentStatusCompleted, PaymentStatusFailed},
		PaymentStatusCompleted:  {PaymentStatusRefunded},
		PaymentStatusFailed:     {PaymentStatusPending}, // retry
		PaymentStatusCancelled:  {},
		PaymentStatusRefunded:   {},
	}

	allowed, ok := validTransitions[p.Status]
	if !ok {
		return false
	}
	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

// CalculateBalanceDue calculates the remaining balance due
func (i *Invoice) CalculateBalanceDue() decimal.Decimal {
	return i.TotalAmount.Sub(i.PaidAmount)
}

// IsOverdue checks if the invoice is past due
func (i *Invoice) IsOverdue() bool {
	if i.Status == InvoiceStatusPaid || i.Status == InvoiceStatusCancelled {
		return false
	}
	return time.Now().After(i.DueDate)
}

// IsFullyCollected checks if the full amount has been collected
func (ct *CollectionTask) IsFullyCollected() bool {
	return ct.AmountCollected.GreaterThanOrEqual(ct.AmountDue)
}

// CalculateStatus calculates the current status of an installment
func (i *Installment) CalculateStatus() string {
	if i.PaidAmount.GreaterThanOrEqual(i.Amount) {
		return "paid"
	}
	if time.Now().After(i.DueDate) {
		return "overdue"
	}
	return "pending"
}

// CalculateNetAmount calculates the net settlement amount after fees
func (s *Settlement) CalculateNetAmount() decimal.Decimal {
	return s.GrossAmount.Sub(s.Fees)
}

// CreditLimitMultiplier returns the credit limit multiplier for a customer type
func (ct CustomerType) CreditLimitMultiplier() float64 {
	switch ct {
	case CustomerTypeConsumer:
		return 0.5
	case CustomerTypeRetailer:
		return 1.0
	case CustomerTypeWholesaler:
		return 2.0
	case CustomerTypeDistributor:
		return 3.0
	case CustomerTypeEnterprise:
		return 5.0
	default:
		return 1.0
	}
}
