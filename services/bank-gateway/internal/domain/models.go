// Package domain contains the core domain models for the Bank Gateway service.
// Following DDD principles for financial integration with banks and mobile money.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Domain errors
var (
	ErrBankConnectionRequired = errors.New("bank connection is required")
	ErrInvalidBankCode        = errors.New("invalid bank code")
	ErrInvalidAccountNumber   = errors.New("invalid account number")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrPaymentLimitExceeded   = errors.New("payment limit exceeded")
	ErrDuplicatePayment       = errors.New("duplicate payment reference")
	ErrPaymentExpired         = errors.New("payment has expired")
	ErrAccountNotVerified     = errors.New("account not verified")
	ErrInvalidPaymentStatus   = errors.New("invalid payment status transition")
	ErrBatchSizeExceeded      = errors.New("batch size exceeded maximum allowed")
	ErrVirtualAccountExpired  = errors.New("virtual account has expired")
)

// ============================================================================
// Value Objects
// ============================================================================

// BankConnectionType represents the type of bank connection
type BankConnectionType string

const (
	ConnectionNIBSSNIP      BankConnectionType = "NIBSS_NIP"
	ConnectionNIBSSNEFT     BankConnectionType = "NIBSS_NEFT"
	ConnectionRTGS          BankConnectionType = "RTGS"
	ConnectionMobileMoney   BankConnectionType = "MOBILE_MONEY"
	ConnectionOpenBanking   BankConnectionType = "OPEN_BANKING"
	ConnectionCardProcessor BankConnectionType = "CARD_PROCESSOR"
)

// ConnectionStatus represents the status of a bank connection
type ConnectionStatus string

const (
	ConnectionStatusActive      ConnectionStatus = "active"
	ConnectionStatusInactive    ConnectionStatus = "inactive"
	ConnectionStatusMaintenance ConnectionStatus = "maintenance"
	ConnectionStatusDegraded    ConnectionStatus = "degraded"
)

// PaymentType represents types of payments
type PaymentType string

const (
	PaymentTypeSupplier     PaymentType = "SUPPLIER_PAYMENT"
	PaymentTypeWorkerPayout PaymentType = "WORKER_PAYOUT"
	PaymentTypeRefund       PaymentType = "CUSTOMER_REFUND"
	PaymentTypeCollection   PaymentType = "COLLECTION"
	PaymentTypeInternal     PaymentType = "INTERNAL_TRANSFER"
	PaymentTypeBulk         PaymentType = "BULK_DISBURSEMENT"
)

// PaymentDirection represents the direction of payment
type PaymentDirection string

const (
	PaymentDirectionInbound  PaymentDirection = "inbound"
	PaymentDirectionOutbound PaymentDirection = "outbound"
)

// PaymentRail represents the payment rail used
type PaymentRail string

const (
	RailNIBSSNIP     PaymentRail = "NIBSS_NIP"
	RailNIBSSNEFT    PaymentRail = "NIBSS_NEFT"
	RailRTGS         PaymentRail = "RTGS"
	RailMobileMoney  PaymentRail = "MOBILE_MONEY"
	RailCardPayment  PaymentRail = "CARD"
	RailInternalBook PaymentRail = "INTERNAL"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending          PaymentStatus = "PENDING"
	PaymentStatusAwaitingApproval PaymentStatus = "AWAITING_APPROVAL"
	PaymentStatusProcessing       PaymentStatus = "PROCESSING"
	PaymentStatusSentToBank       PaymentStatus = "SENT_TO_BANK"
	PaymentStatusConfirmed        PaymentStatus = "CONFIRMED"
	PaymentStatusFailed           PaymentStatus = "FAILED"
	PaymentStatusReversed         PaymentStatus = "REVERSED"
	PaymentStatusExpired          PaymentStatus = "EXPIRED"
)

// AccountOwnerType represents the owner type for virtual accounts
type AccountOwnerType string

const (
	OwnerCustomer AccountOwnerType = "CUSTOMER"
	OwnerSupplier AccountOwnerType = "SUPPLIER"
	OwnerWorker   AccountOwnerType = "WORKER"
	OwnerGroup    AccountOwnerType = "GROUP"
	OwnerEscrow   AccountOwnerType = "ESCROW"
)

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusFrozen  AccountStatus = "frozen"
	AccountStatusClosed  AccountStatus = "closed"
	AccountStatusExpired AccountStatus = "expired"
)

// AccountDetails represents bank account details
type AccountDetails struct {
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name,omitempty"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Provider      string `json:"provider,omitempty"`
}

// Validate validates account details
func (a AccountDetails) Validate() error {
	if a.AccountNumber == "" {
		return ErrInvalidAccountNumber
	}
	if a.BankCode == "" && a.Phone == "" {
		return ErrInvalidBankCode
	}
	return nil
}

// ============================================================================
// Aggregates
// ============================================================================

// BankConnection represents a connection to a bank or payment provider
type BankConnection struct {
	ID                  uuid.UUID              `json:"id"`
	TenantID            uuid.UUID              `json:"tenant_id"`
	BankCode            string                 `json:"bank_code"`
	BankName            string                 `json:"bank_name"`
	ConnectionType      BankConnectionType     `json:"connection_type"`
	APIEndpoint         string                 `json:"api_endpoint"`
	Status              ConnectionStatus       `json:"status"`
	LastHealthCheck     *time.Time             `json:"last_health_check,omitempty"`
	SupportedOperations []string               `json:"supported_operations"`
	RateLimits          map[string]interface{} `json:"rate_limits,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// PaymentTransaction represents a payment transaction
type PaymentTransaction struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	Reference         string                 `json:"reference"`
	ExternalReference string                 `json:"external_reference,omitempty"`
	Type              PaymentType            `json:"type"`
	Direction         PaymentDirection       `json:"direction"`
	Rail              PaymentRail            `json:"rail"`
	SourceAccount     AccountDetails         `json:"source_account"`
	DestAccount       AccountDetails         `json:"destination_account"`
	Amount            decimal.Decimal        `json:"amount"`
	Currency          string                 `json:"currency"`
	Fee               decimal.Decimal        `json:"fee"`
	Status            PaymentStatus          `json:"status"`
	StatusHistory     []PaymentStatusChange  `json:"status_history"`
	InitiatedBy       uuid.UUID              `json:"initiated_by"`
	ApprovedBy        *uuid.UUID             `json:"approved_by,omitempty"`
	Narration         string                 `json:"narration,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	BankResponse      map[string]interface{} `json:"bank_response,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	CompletedAt       *time.Time             `json:"completed_at,omitempty"`

	// Domain events
	events []DomainEvent
}

// PaymentStatusChange represents a status change in payment history
type PaymentStatusChange struct {
	Status    PaymentStatus `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Details   string        `json:"details,omitempty"`
}

// VirtualAccount represents a virtual account for collections
type VirtualAccount struct {
	ID              uuid.UUID              `json:"id"`
	TenantID        uuid.UUID              `json:"tenant_id"`
	AccountNumber   string                 `json:"account_number"`
	BankCode        string                 `json:"bank_code"`
	AccountName     string                 `json:"account_name"`
	OwnerID         uuid.UUID              `json:"owner_id"`
	OwnerType       AccountOwnerType       `json:"owner_type"`
	Currency        string                 `json:"currency"`
	Balance         decimal.Decimal        `json:"balance"`
	Status          AccountStatus          `json:"status"`
	ExpiryDate      *time.Time             `json:"expiry_date,omitempty"`
	CollectionRules map[string]interface{} `json:"collection_rules,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// BulkPaymentBatch represents a batch of payments
type BulkPaymentBatch struct {
	ID               uuid.UUID            `json:"id"`
	TenantID         uuid.UUID            `json:"tenant_id"`
	Reference        string               `json:"reference"`
	Name             string               `json:"name"`
	TotalAmount      decimal.Decimal      `json:"total_amount"`
	TotalCount       int                  `json:"total_count"`
	SuccessfulCount  int                  `json:"successful_count"`
	FailedCount      int                  `json:"failed_count"`
	Status           BatchStatus          `json:"status"`
	Payments         []PaymentTransaction `json:"payments,omitempty"`
	InitiatedBy      uuid.UUID            `json:"initiated_by"`
	ApprovalWorkflow *uuid.UUID           `json:"approval_workflow,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	CompletedAt      *time.Time           `json:"completed_at,omitempty"`
}

// BatchStatus represents the status of a batch
type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "PENDING"
	BatchStatusProcessing BatchStatus = "PROCESSING"
	BatchStatusCompleted  BatchStatus = "COMPLETED"
	BatchStatusFailed     BatchStatus = "FAILED"
	BatchStatusPartial    BatchStatus = "PARTIAL"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// PaymentInitiatedEvent is raised when a payment is initiated
type PaymentInitiatedEvent struct {
	PaymentID uuid.UUID
	Reference string
	Amount    decimal.Decimal
	Timestamp time.Time
}

func (e PaymentInitiatedEvent) EventType() string     { return "payment.initiated" }
func (e PaymentInitiatedEvent) OccurredAt() time.Time { return e.Timestamp }

// PaymentCompletedEvent is raised when a payment is completed
type PaymentCompletedEvent struct {
	PaymentID         uuid.UUID
	Reference         string
	ExternalReference string
	Status            PaymentStatus
	Timestamp         time.Time
}

func (e PaymentCompletedEvent) EventType() string     { return "payment.completed" }
func (e PaymentCompletedEvent) OccurredAt() time.Time { return e.Timestamp }

// VirtualAccountCreditedEvent is raised when a virtual account receives funds
type VirtualAccountCreditedEvent struct {
	AccountID     uuid.UUID
	AccountNumber string
	Amount        decimal.Decimal
	Reference     string
	Timestamp     time.Time
}

func (e VirtualAccountCreditedEvent) EventType() string     { return "virtual_account.credited" }
func (e VirtualAccountCreditedEvent) OccurredAt() time.Time { return e.Timestamp }

// ============================================================================
// Aggregate Methods
// ============================================================================

// NewPaymentTransaction creates a new payment transaction
func NewPaymentTransaction(tenantID uuid.UUID, paymentType PaymentType,
	source, dest AccountDetails, amount decimal.Decimal, currency string,
	initiatedBy uuid.UUID, narration string) (*PaymentTransaction, error) {

	if err := source.Validate(); err != nil {
		return nil, err
	}
	if err := dest.Validate(); err != nil {
		return nil, err
	}
	if amount.IsNegative() || amount.IsZero() {
		return nil, errors.New("amount must be positive")
	}

	now := time.Now()
	tx := &PaymentTransaction{
		ID:            uuid.New(),
		TenantID:      tenantID,
		Reference:     generateReference(),
		Type:          paymentType,
		Direction:     PaymentDirectionOutbound,
		SourceAccount: source,
		DestAccount:   dest,
		Amount:        amount,
		Currency:      currency,
		Status:        PaymentStatusPending,
		InitiatedBy:   initiatedBy,
		Narration:     narration,
		CreatedAt:     now,
		StatusHistory: []PaymentStatusChange{
			{Status: PaymentStatusPending, Timestamp: now},
		},
	}

	tx.events = append(tx.events, PaymentInitiatedEvent{
		PaymentID: tx.ID,
		Reference: tx.Reference,
		Amount:    tx.Amount,
		Timestamp: now,
	})

	return tx, nil
}

// UpdateStatus updates the payment status
func (p *PaymentTransaction) UpdateStatus(newStatus PaymentStatus, details string) error {
	if !p.isValidTransition(newStatus) {
		return ErrInvalidPaymentStatus
	}

	now := time.Now()
	p.Status = newStatus
	p.StatusHistory = append(p.StatusHistory, PaymentStatusChange{
		Status:    newStatus,
		Timestamp: now,
		Details:   details,
	})

	if newStatus == PaymentStatusConfirmed || newStatus == PaymentStatusFailed {
		p.CompletedAt = &now
		p.events = append(p.events, PaymentCompletedEvent{
			PaymentID:         p.ID,
			Reference:         p.Reference,
			ExternalReference: p.ExternalReference,
			Status:            newStatus,
			Timestamp:         now,
		})
	}

	return nil
}

func (p *PaymentTransaction) isValidTransition(newStatus PaymentStatus) bool {
	validTransitions := map[PaymentStatus][]PaymentStatus{
		PaymentStatusPending:          {PaymentStatusAwaitingApproval, PaymentStatusProcessing, PaymentStatusFailed},
		PaymentStatusAwaitingApproval: {PaymentStatusProcessing, PaymentStatusFailed, PaymentStatusExpired},
		PaymentStatusProcessing:       {PaymentStatusSentToBank, PaymentStatusFailed},
		PaymentStatusSentToBank:       {PaymentStatusConfirmed, PaymentStatusFailed},
		PaymentStatusConfirmed:        {PaymentStatusReversed},
	}

	allowed, ok := validTransitions[p.Status]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

// GetEvents returns and clears domain events
func (p *PaymentTransaction) GetEvents() []DomainEvent {
	events := p.events
	p.events = nil
	return events
}

// IsExpired checks if a virtual account is expired
func (v *VirtualAccount) IsExpired() bool {
	if v.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*v.ExpiryDate)
}

// Credit adds funds to a virtual account
func (v *VirtualAccount) Credit(amount decimal.Decimal, reference string) error {
	if v.Status != AccountStatusActive {
		return errors.New("account is not active")
	}
	if v.IsExpired() {
		return ErrVirtualAccountExpired
	}
	v.Balance = v.Balance.Add(amount)
	return nil
}

func generateReference() string {
	return "PAY" + uuid.New().String()[:8]
}
