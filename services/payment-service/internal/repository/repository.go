// Package repository provides PostgreSQL data access for the Payment Service.
// Handles payments, wallets, credit accounts, disbursements, bank accounts, and reconciliation.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("payment-service/repository")

// Common errors
var (
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateKey        = errors.New("duplicate key violation")
	ErrOptimisticLock      = errors.New("optimistic lock conflict")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrAccountFrozen       = errors.New("account is frozen")
	ErrCreditLimitExceeded = errors.New("credit limit exceeded")
)

// ============================================================================
// Repository Interface
// ============================================================================

type Repository interface {
	// Payment operations
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPayment(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetPaymentByReference(ctx context.Context, reference string) (*Payment, error)
	UpdatePayment(ctx context.Context, payment *Payment) error
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status PaymentStatus, providerRef *string) error
	ListPayments(ctx context.Context, filter PaymentFilter) (*PaymentList, error)
	
	// Refund operations
	CreateRefund(ctx context.Context, refund *Refund) error
	GetRefund(ctx context.Context, id uuid.UUID) (*Refund, error)
	GetRefundsByPayment(ctx context.Context, paymentID uuid.UUID) ([]*Refund, error)
	UpdateRefundStatus(ctx context.Context, id uuid.UUID, status RefundStatus, providerRef *string) error
	
	// Wallet operations
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetWallet(ctx context.Context, id uuid.UUID) (*Wallet, error)
	GetWalletByCustomer(ctx context.Context, customerID uuid.UUID) (*Wallet, error)
	UpdateWallet(ctx context.Context, wallet *Wallet) error
	CreditWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error
	DebitWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error
	FreezeWallet(ctx context.Context, id uuid.UUID, reason string) error
	UnfreezeWallet(ctx context.Context, id uuid.UUID) error
	
	// Wallet transaction operations
	CreateWalletTransaction(ctx context.Context, tx *WalletTransaction) error
	GetWalletTransactions(ctx context.Context, walletID uuid.UUID, filter TransactionFilter) (*TransactionList, error)
	
	// Credit account operations
	CreateCreditAccount(ctx context.Context, account *CreditAccount) error
	GetCreditAccount(ctx context.Context, id uuid.UUID) (*CreditAccount, error)
	GetCreditAccountByCustomer(ctx context.Context, customerID uuid.UUID) (*CreditAccount, error)
	UpdateCreditAccount(ctx context.Context, account *CreditAccount) error
	UseCreditLine(ctx context.Context, id uuid.UUID, amount decimal.Decimal, orderID uuid.UUID) error
	RepayCreditLine(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error
	
	// Credit limit increase operations
	CreateCreditLimitRequest(ctx context.Context, request *CreditLimitRequest) error
	GetCreditLimitRequest(ctx context.Context, id uuid.UUID) (*CreditLimitRequest, error)
	UpdateCreditLimitRequest(ctx context.Context, request *CreditLimitRequest) error
	GetPendingCreditLimitRequests(ctx context.Context, tenantID uuid.UUID) ([]*CreditLimitRequest, error)
	
	// Credit transaction operations
	CreateCreditTransaction(ctx context.Context, tx *CreditTransaction) error
	GetCreditTransactions(ctx context.Context, accountID uuid.UUID, filter TransactionFilter) (*CreditTransactionList, error)
	GetCreditStatement(ctx context.Context, accountID uuid.UUID, from, to time.Time) (*CreditStatement, error)
	
	// Credit scoring
	GetCreditScore(ctx context.Context, customerID uuid.UUID) (*CreditScore, error)
	UpsertCreditScore(ctx context.Context, score *CreditScore) error
	GetCreditScoreFactors(ctx context.Context, customerID uuid.UUID) (*CreditScoreFactors, error)
	
	// Disbursement operations
	CreateDisbursement(ctx context.Context, disbursement *Disbursement) error
	GetDisbursement(ctx context.Context, id uuid.UUID) (*Disbursement, error)
	UpdateDisbursement(ctx context.Context, disbursement *Disbursement) error
	ListDisbursements(ctx context.Context, filter DisbursementFilter) (*DisbursementList, error)
	
	// Disbursement batch operations
	CreateDisbursementBatch(ctx context.Context, batch *DisbursementBatch) error
	GetDisbursementBatch(ctx context.Context, id uuid.UUID) (*DisbursementBatch, error)
	UpdateDisbursementBatch(ctx context.Context, batch *DisbursementBatch) error
	
	// Bank account operations
	CreateBankAccount(ctx context.Context, account *BankAccount) error
	GetBankAccount(ctx context.Context, id uuid.UUID) (*BankAccount, error)
	GetBankAccountsByCustomer(ctx context.Context, customerID uuid.UUID) ([]*BankAccount, error)
	UpdateBankAccount(ctx context.Context, account *BankAccount) error
	DeleteBankAccount(ctx context.Context, id uuid.UUID) error
	SetDefaultBankAccount(ctx context.Context, customerID, accountID uuid.UUID) error
	
	// Virtual account operations
	CreateVirtualAccount(ctx context.Context, account *VirtualAccount) error
	GetVirtualAccount(ctx context.Context, id uuid.UUID) (*VirtualAccount, error)
	GetVirtualAccountByNumber(ctx context.Context, accountNumber string) (*VirtualAccount, error)
	GetVirtualAccountByCustomer(ctx context.Context, customerID uuid.UUID) (*VirtualAccount, error)
	UpdateVirtualAccount(ctx context.Context, account *VirtualAccount) error
	
	// Reconciliation operations
	CreateReconciliationRecord(ctx context.Context, record *ReconciliationRecord) error
	GetReconciliationReport(ctx context.Context, filter ReconciliationFilter) (*ReconciliationReport, error)
	FlagPaymentForReview(ctx context.Context, paymentID uuid.UUID, reason string) error
	
	// Analytics
	GetPaymentStats(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*PaymentStats, error)
	GetRevenueReport(ctx context.Context, tenantID uuid.UUID, from, to time.Time, groupBy string) ([]*RevenueEntry, error)
	
	// Transaction support
	WithTx(ctx context.Context, fn func(Repository) error) error
}

// ============================================================================
// Domain Types
// ============================================================================

type PaymentStatus string
const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

type PaymentMethod string
const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodUSSD         PaymentMethod = "ussd"
	PaymentMethodMobileMoney  PaymentMethod = "mobile_money"
	PaymentMethodWallet       PaymentMethod = "wallet"
	PaymentMethodCreditLine   PaymentMethod = "credit_line"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodPOS          PaymentMethod = "pos"
)

type RefundStatus string
const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
	RefundStatusFailed    RefundStatus = "failed"
)

type WalletStatus string
const (
	WalletStatusActive WalletStatus = "active"
	WalletStatusFrozen WalletStatus = "frozen"
	WalletStatusClosed WalletStatus = "closed"
)

type CreditAccountStatus string
const (
	CreditAccountStatusActive    CreditAccountStatus = "active"
	CreditAccountStatusSuspended CreditAccountStatus = "suspended"
	CreditAccountStatusClosed    CreditAccountStatus = "closed"
)

type DisbursementStatus string
const (
	DisbursementStatusPending    DisbursementStatus = "pending"
	DisbursementStatusProcessing DisbursementStatus = "processing"
	DisbursementStatusCompleted  DisbursementStatus = "completed"
	DisbursementStatusFailed     DisbursementStatus = "failed"
	DisbursementStatusCancelled  DisbursementStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	OrderID           *uuid.UUID             `json:"order_id,omitempty"`
	CustomerID        uuid.UUID              `json:"customer_id"`
	Amount            decimal.Decimal        `json:"amount"`
	Currency          string                 `json:"currency"`
	Method            PaymentMethod          `json:"method"`
	Provider          string                 `json:"provider"`
	Status            PaymentStatus          `json:"status"`
	Reference         string                 `json:"reference"`
	ProviderReference *string                `json:"provider_reference,omitempty"`
	CardToken         *string                `json:"card_token,omitempty"`
	CardLast4         *string                `json:"card_last4,omitempty"`
	CardBrand         *string                `json:"card_brand,omitempty"`
	BankCode          *string                `json:"bank_code,omitempty"`
	BankAccountNumber *string                `json:"bank_account_number,omitempty"`
	MobileNumber      *string                `json:"mobile_number,omitempty"`
	MobileNetwork     *string                `json:"mobile_network,omitempty"`
	WalletID          *uuid.UUID             `json:"wallet_id,omitempty"`
	CreditAccountID   *uuid.UUID             `json:"credit_account_id,omitempty"`
	CollectedBy       *uuid.UUID             `json:"collected_by,omitempty"`
	CollectorType     *string                `json:"collector_type,omitempty"`
	Fee               decimal.Decimal        `json:"fee"`
	FailureReason     *string                `json:"failure_reason,omitempty"`
	RefundedAmount    decimal.Decimal        `json:"refunded_amount"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	PaidAt            *time.Time             `json:"paid_at,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Version           int                    `json:"version"`
}

// Refund represents a payment refund
type Refund struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	PaymentID         uuid.UUID       `json:"payment_id"`
	Amount            decimal.Decimal `json:"amount"`
	Currency          string          `json:"currency"`
	Status            RefundStatus    `json:"status"`
	Reason            string          `json:"reason"`
	Reference         string          `json:"reference"`
	ProviderReference *string         `json:"provider_reference,omitempty"`
	FailureReason     *string         `json:"failure_reason,omitempty"`
	ProcessedAt       *time.Time      `json:"processed_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Wallet represents a customer wallet
type Wallet struct {
	ID              uuid.UUID              `json:"id"`
	TenantID        uuid.UUID              `json:"tenant_id"`
	CustomerID      uuid.UUID              `json:"customer_id"`
	CustomerType    string                 `json:"customer_type"`
	Name            *string                `json:"name,omitempty"`
	Balance         decimal.Decimal        `json:"balance"`
	Currency        string                 `json:"currency"`
	Status          WalletStatus           `json:"status"`
	FreezeReason    *string                `json:"freeze_reason,omitempty"`
	FineractID      *string                `json:"fineract_id,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Version         int                    `json:"version"`
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID          uuid.UUID              `json:"id"`
	WalletID    uuid.UUID              `json:"wallet_id"`
	Type        string                 `json:"type"` // credit, debit
	Amount      decimal.Decimal        `json:"amount"`
	Currency    string                 `json:"currency"`
	Balance     decimal.Decimal        `json:"balance_after"`
	Source      string                 `json:"source"` // payment, refund, earning, transfer, adjustment
	SourceID    *uuid.UUID             `json:"source_id,omitempty"`
	Reference   string                 `json:"reference"`
	Description *string                `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CreditAccount represents a customer credit account
type CreditAccount struct {
	ID               uuid.UUID              `json:"id"`
	TenantID         uuid.UUID              `json:"tenant_id"`
	CustomerID       uuid.UUID              `json:"customer_id"`
	CustomerType     string                 `json:"customer_type"`
	CreditLimit      decimal.Decimal        `json:"credit_limit"`
	AvailableCredit  decimal.Decimal        `json:"available_credit"`
	OutstandingBalance decimal.Decimal      `json:"outstanding_balance"`
	Currency         string                 `json:"currency"`
	PaymentTermsDays int                    `json:"payment_terms_days"`
	InterestRate     decimal.Decimal        `json:"interest_rate"`
	GracePeriodDays  int                    `json:"grace_period_days"`
	BillingCycle     string                 `json:"billing_cycle"`
	Status           CreditAccountStatus    `json:"status"`
	OverdueAmount    decimal.Decimal        `json:"overdue_amount"`
	LastPaymentDate  *time.Time             `json:"last_payment_date,omitempty"`
	NextDueDate      *time.Time             `json:"next_due_date,omitempty"`
	FineractLoanID   *string                `json:"fineract_loan_id,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Version          int                    `json:"version"`
}

// CreditLimitRequest represents a credit limit increase request
type CreditLimitRequest struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	CreditAccountID  uuid.UUID       `json:"credit_account_id"`
	CurrentLimit     decimal.Decimal `json:"current_limit"`
	RequestedLimit   decimal.Decimal `json:"requested_limit"`
	Status           string          `json:"status"` // pending, approved, rejected
	Reason           *string         `json:"reason,omitempty"`
	ReviewedBy       *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewNotes      *string         `json:"review_notes,omitempty"`
	ReviewedAt       *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// CreditTransaction represents a credit account transaction
type CreditTransaction struct {
	ID              uuid.UUID              `json:"id"`
	CreditAccountID uuid.UUID              `json:"credit_account_id"`
	Type            string                 `json:"type"` // usage, repayment, fee, interest, adjustment
	Amount          decimal.Decimal        `json:"amount"`
	Currency        string                 `json:"currency"`
	OrderID         *uuid.UUID             `json:"order_id,omitempty"`
	PaymentMethod   *string                `json:"payment_method,omitempty"`
	Reference       string                 `json:"reference"`
	Description     *string                `json:"description,omitempty"`
	DueDate         *time.Time             `json:"due_date,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// CreditScore represents a customer's credit score
type CreditScore struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	CustomerID       uuid.UUID `json:"customer_id"`
	Score            int       `json:"score"`
	Rating           string    `json:"rating"` // excellent, good, fair, poor
	PaymentHistory   int       `json:"payment_history_score"`
	OrderFrequency   int       `json:"order_frequency_score"`
	CreditUtilization int      `json:"credit_utilization_score"`
	AccountAge       int       `json:"account_age_score"`
	LastCalculatedAt time.Time `json:"last_calculated_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreditScoreFactors contains the factors used to calculate credit score
type CreditScoreFactors struct {
	CustomerID        uuid.UUID       `json:"customer_id"`
	TotalOrders       int             `json:"total_orders"`
	PaidOnTime        int             `json:"paid_on_time"`
	PaidLate          int             `json:"paid_late"`
	AverageOrderValue decimal.Decimal `json:"average_order_value"`
	TotalSpent        decimal.Decimal `json:"total_spent"`
	CreditLimit       decimal.Decimal `json:"credit_limit"`
	CurrentBalance    decimal.Decimal `json:"current_balance"`
	AccountAgeDays    int             `json:"account_age_days"`
	OverdueCount      int             `json:"overdue_count"`
}

// CreditStatement represents a credit account statement
type CreditStatement struct {
	AccountID          uuid.UUID            `json:"account_id"`
	FromDate           time.Time            `json:"from_date"`
	ToDate             time.Time            `json:"to_date"`
	OpeningBalance     decimal.Decimal      `json:"opening_balance"`
	ClosingBalance     decimal.Decimal      `json:"closing_balance"`
	TotalUsage         decimal.Decimal      `json:"total_usage"`
	TotalRepayments    decimal.Decimal      `json:"total_repayments"`
	TotalFees          decimal.Decimal      `json:"total_fees"`
	TotalInterest      decimal.Decimal      `json:"total_interest"`
	Transactions       []*CreditTransaction `json:"transactions"`
}

// Disbursement represents a payout to a recipient
type Disbursement struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	BatchID           *uuid.UUID             `json:"batch_id,omitempty"`
	RecipientID       uuid.UUID              `json:"recipient_id"`
	RecipientType     string                 `json:"recipient_type"`
	Amount            decimal.Decimal        `json:"amount"`
	Currency          string                 `json:"currency"`
	DisbursementType  string                 `json:"disbursement_type"`
	Status            DisbursementStatus     `json:"status"`
	PayoutMethod      string                 `json:"payout_method"`
	BankAccountID     *uuid.UUID             `json:"bank_account_id,omitempty"`
	MobileNumber      *string                `json:"mobile_number,omitempty"`
	MobileNetwork     *string                `json:"mobile_network,omitempty"`
	WalletID          *uuid.UUID             `json:"wallet_id,omitempty"`
	Reference         string                 `json:"reference"`
	ProviderReference *string                `json:"provider_reference,omitempty"`
	Description       *string                `json:"description,omitempty"`
	FailureReason     *string                `json:"failure_reason,omitempty"`
	ScheduledFor      *time.Time             `json:"scheduled_for,omitempty"`
	ProcessedAt       *time.Time             `json:"processed_at,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// DisbursementBatch represents a batch of disbursements
type DisbursementBatch struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	TotalCount      int             `json:"total_count"`
	ProcessedCount  int             `json:"processed_count"`
	SuccessCount    int             `json:"success_count"`
	FailedCount     int             `json:"failed_count"`
	Status          string          `json:"status"` // pending, processing, completed, failed
	ProcessingMode  string          `json:"processing_mode"`
	ScheduledFor    *time.Time      `json:"scheduled_for,omitempty"`
	StartedAt       *time.Time      `json:"started_at,omitempty"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// BankAccount represents a customer's bank account
type BankAccount struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	CustomerID    uuid.UUID  `json:"customer_id"`
	BankCode      string     `json:"bank_code"`
	BankName      string     `json:"bank_name"`
	AccountNumber string     `json:"account_number"`
	AccountName   string     `json:"account_name"`
	IsVerified    bool       `json:"is_verified"`
	IsDefault     bool       `json:"is_default"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// VirtualAccount represents a virtual bank account
type VirtualAccount struct {
	ID              uuid.UUID              `json:"id"`
	TenantID        uuid.UUID              `json:"tenant_id"`
	CustomerID      uuid.UUID              `json:"customer_id"`
	Provider        string                 `json:"provider"`
	BankName        string                 `json:"bank_name"`
	AccountNumber   string                 `json:"account_number"`
	AccountName     string                 `json:"account_name"`
	ProviderRef     *string                `json:"provider_reference,omitempty"`
	Status          string                 `json:"status"` // active, closed
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// ReconciliationRecord represents a reconciliation entry
type ReconciliationRecord struct {
	ID              uuid.UUID              `json:"id"`
	TenantID        uuid.UUID              `json:"tenant_id"`
	PaymentID       uuid.UUID              `json:"payment_id"`
	Provider        string                 `json:"provider"`
	ProviderRef     string                 `json:"provider_reference"`
	ProviderAmount  decimal.Decimal        `json:"provider_amount"`
	SystemAmount    decimal.Decimal        `json:"system_amount"`
	Difference      decimal.Decimal        `json:"difference"`
	Status          string                 `json:"status"` // matched, discrepancy, flagged
	FlagReason      *string                `json:"flag_reason,omitempty"`
	ResolvedAt      *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy      *uuid.UUID             `json:"resolved_by,omitempty"`
	Notes           *string                `json:"notes,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// ============================================================================
// Filter and Result Types
// ============================================================================

type PaymentFilter struct {
	TenantID    uuid.UUID
	CustomerID  *uuid.UUID
	OrderID     *uuid.UUID
	Status      []PaymentStatus
	Method      []PaymentMethod
	Provider    *string
	Reference   *string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

type PaymentList struct {
	Payments []*Payment `json:"payments"`
	Total    int        `json:"total"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type TransactionFilter struct {
	Type     *string
	Source   *string
	FromDate *time.Time
	ToDate   *time.Time
	Limit    int
	Offset   int
}

type TransactionList struct {
	Transactions []*WalletTransaction `json:"transactions"`
	Total        int                  `json:"total"`
	Limit        int                  `json:"limit"`
	Offset       int                  `json:"offset"`
}

type CreditTransactionList struct {
	Transactions []*CreditTransaction `json:"transactions"`
	Total        int                  `json:"total"`
	Limit        int                  `json:"limit"`
	Offset       int                  `json:"offset"`
}

type DisbursementFilter struct {
	TenantID      uuid.UUID
	RecipientID   *uuid.UUID
	RecipientType *string
	BatchID       *uuid.UUID
	Status        []DisbursementStatus
	Type          *string
	FromDate      *time.Time
	ToDate        *time.Time
	Limit         int
	Offset        int
}

type DisbursementList struct {
	Disbursements []*Disbursement `json:"disbursements"`
	Total         int             `json:"total"`
	Limit         int             `json:"limit"`
	Offset        int             `json:"offset"`
}

type ReconciliationFilter struct {
	TenantID uuid.UUID
	Provider *string
	Status   *string
	FromDate *time.Time
	ToDate   *time.Time
}

type ReconciliationReport struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	FromDate        time.Time       `json:"from_date"`
	ToDate          time.Time       `json:"to_date"`
	TotalPayments   int             `json:"total_payments"`
	Matched         int             `json:"matched"`
	Discrepancies   int             `json:"discrepancies"`
	Flagged         int             `json:"flagged"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	DiscrepancyAmt  decimal.Decimal `json:"discrepancy_amount"`
}

type PaymentStats struct {
	TenantID      uuid.UUID                  `json:"tenant_id"`
	FromDate      time.Time                  `json:"from_date"`
	ToDate        time.Time                  `json:"to_date"`
	TotalPayments int64                      `json:"total_payments"`
	TotalAmount   decimal.Decimal            `json:"total_amount"`
	TotalFees     decimal.Decimal            `json:"total_fees"`
	ByStatus      map[PaymentStatus]int64    `json:"by_status"`
	ByMethod      map[PaymentMethod]int64    `json:"by_method"`
	ByProvider    map[string]int64           `json:"by_provider"`
	SuccessRate   float64                    `json:"success_rate"`
}

type RevenueEntry struct {
	Period  string          `json:"period"`
	Gross   decimal.Decimal `json:"gross"`
	Refunds decimal.Decimal `json:"refunds"`
	Net     decimal.Decimal `json:"net"`
	Fees    decimal.Decimal `json:"fees"`
	Count   int64           `json:"count"`
}

// ============================================================================
// PostgreSQL Implementation
// ============================================================================

type postgresRepository struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

// ============================================================================
// Payment Operations
// ============================================================================

func (r *postgresRepository) CreatePayment(ctx context.Context, payment *Payment) error {
	ctx, span := tracer.Start(ctx, "repository.CreatePayment")
	defer span.End()
	
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	payment.CreatedAt = time.Now().UTC()
	payment.UpdatedAt = payment.CreatedAt
	payment.Version = 1
	
	if payment.Status == "" {
		payment.Status = PaymentStatusPending
	}
	
	metadata, _ := json.Marshal(payment.Metadata)
	
	query := `
		INSERT INTO payments (
			id, tenant_id, order_id, customer_id, amount, currency, method, provider,
			status, reference, provider_reference, card_token, card_last4, card_brand,
			bank_code, bank_account_number, mobile_number, mobile_network,
			wallet_id, credit_account_id, collected_by, collector_type,
			fee, failure_reason, refunded_amount, metadata, paid_at,
			created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)`
	
	_, err := r.pool.Exec(ctx, query,
		payment.ID, payment.TenantID, payment.OrderID, payment.CustomerID,
		payment.Amount, payment.Currency, payment.Method, payment.Provider,
		payment.Status, payment.Reference, payment.ProviderReference,
		payment.CardToken, payment.CardLast4, payment.CardBrand,
		payment.BankCode, payment.BankAccountNumber, payment.MobileNumber, payment.MobileNetwork,
		payment.WalletID, payment.CreditAccountID, payment.CollectedBy, payment.CollectorType,
		payment.Fee, payment.FailureReason, payment.RefundedAmount, metadata, payment.PaidAt,
		payment.CreatedAt, payment.UpdatedAt, payment.Version,
	)
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateKey
		}
		span.RecordError(err)
		return fmt.Errorf("create payment: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetPayment(ctx context.Context, id uuid.UUID) (*Payment, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPayment")
	defer span.End()
	span.SetAttributes(attribute.String("payment.id", id.String()))
	
	query := `
		SELECT id, tenant_id, order_id, customer_id, amount, currency, method, provider,
			status, reference, provider_reference, card_token, card_last4, card_brand,
			bank_code, bank_account_number, mobile_number, mobile_network,
			wallet_id, credit_account_id, collected_by, collector_type,
			fee, failure_reason, refunded_amount, metadata, paid_at,
			created_at, updated_at, version
		FROM payments WHERE id = $1`
	
	return r.scanPayment(ctx, r.pool.QueryRow(ctx, query, id))
}

func (r *postgresRepository) GetPaymentByReference(ctx context.Context, reference string) (*Payment, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPaymentByReference")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, order_id, customer_id, amount, currency, method, provider,
			status, reference, provider_reference, card_token, card_last4, card_brand,
			bank_code, bank_account_number, mobile_number, mobile_network,
			wallet_id, credit_account_id, collected_by, collector_type,
			fee, failure_reason, refunded_amount, metadata, paid_at,
			created_at, updated_at, version
		FROM payments WHERE reference = $1`
	
	return r.scanPayment(ctx, r.pool.QueryRow(ctx, query, reference))
}

func (r *postgresRepository) scanPayment(ctx context.Context, row pgx.Row) (*Payment, error) {
	payment := &Payment{}
	var metadata []byte
	
	err := row.Scan(
		&payment.ID, &payment.TenantID, &payment.OrderID, &payment.CustomerID,
		&payment.Amount, &payment.Currency, &payment.Method, &payment.Provider,
		&payment.Status, &payment.Reference, &payment.ProviderReference,
		&payment.CardToken, &payment.CardLast4, &payment.CardBrand,
		&payment.BankCode, &payment.BankAccountNumber, &payment.MobileNumber, &payment.MobileNetwork,
		&payment.WalletID, &payment.CreditAccountID, &payment.CollectedBy, &payment.CollectorType,
		&payment.Fee, &payment.FailureReason, &payment.RefundedAmount, &metadata, &payment.PaidAt,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan payment: %w", err)
	}
	
	json.Unmarshal(metadata, &payment.Metadata)
	return payment, nil
}

func (r *postgresRepository) UpdatePayment(ctx context.Context, payment *Payment) error {
	ctx, span := tracer.Start(ctx, "repository.UpdatePayment")
	defer span.End()
	
	metadata, _ := json.Marshal(payment.Metadata)
	payment.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE payments SET
			status = $3, provider_reference = $4, card_token = $5, card_last4 = $6, card_brand = $7,
			fee = $8, failure_reason = $9, refunded_amount = $10, metadata = $11, paid_at = $12,
			updated_at = $13, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`
	
	err := r.pool.QueryRow(ctx, query,
		payment.ID, payment.Version, payment.Status, payment.ProviderReference,
		payment.CardToken, payment.CardLast4, payment.CardBrand,
		payment.Fee, payment.FailureReason, payment.RefundedAmount, metadata,
		payment.PaidAt, payment.UpdatedAt,
	).Scan(&payment.Version)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptimisticLock
		}
		return fmt.Errorf("update payment: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status PaymentStatus, providerRef *string) error {
	ctx, span := tracer.Start(ctx, "repository.UpdatePaymentStatus")
	defer span.End()
	
	now := time.Now().UTC()
	var paidAt *time.Time
	if status == PaymentStatusCompleted {
		paidAt = &now
	}
	
	query := `
		UPDATE payments SET
			status = $2,
			provider_reference = COALESCE($3, provider_reference),
			paid_at = COALESCE($4, paid_at),
			updated_at = $5
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, status, providerRef, paidAt, now)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) ListPayments(ctx context.Context, filter PaymentFilter) (*PaymentList, error) {
	ctx, span := tracer.Start(ctx, "repository.ListPayments")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++
	
	if filter.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argNum))
		args = append(args, *filter.CustomerID)
		argNum++
	}
	
	if filter.OrderID != nil {
		conditions = append(conditions, fmt.Sprintf("order_id = $%d", argNum))
		args = append(args, *filter.OrderID)
		argNum++
	}
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, s)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argNum))
		args = append(args, *filter.FromDate)
		argNum++
	}
	
	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argNum))
		args = append(args, *filter.ToDate)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM payments WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count payments: %w", err)
	}
	
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, order_id, customer_id, amount, currency, method, provider,
			status, reference, provider_reference, card_token, card_last4, card_brand,
			bank_code, bank_account_number, mobile_number, mobile_network,
			wallet_id, credit_account_id, collected_by, collector_type,
			fee, failure_reason, refunded_amount, metadata, paid_at,
			created_at, updated_at, version
		FROM payments WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`, whereClause, sortBy, sortOrder, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	defer rows.Close()
	
	payments := make([]*Payment, 0)
	for rows.Next() {
		payment := &Payment{}
		var metadata []byte
		
		err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.OrderID, &payment.CustomerID,
			&payment.Amount, &payment.Currency, &payment.Method, &payment.Provider,
			&payment.Status, &payment.Reference, &payment.ProviderReference,
			&payment.CardToken, &payment.CardLast4, &payment.CardBrand,
			&payment.BankCode, &payment.BankAccountNumber, &payment.MobileNumber, &payment.MobileNetwork,
			&payment.WalletID, &payment.CreditAccountID, &payment.CollectedBy, &payment.CollectorType,
			&payment.Fee, &payment.FailureReason, &payment.RefundedAmount, &metadata, &payment.PaidAt,
			&payment.CreatedAt, &payment.UpdatedAt, &payment.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		
		json.Unmarshal(metadata, &payment.Metadata)
		payments = append(payments, payment)
	}
	
	return &PaymentList{
		Payments: payments,
		Total:    total,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// ============================================================================
// Refund Operations
// ============================================================================

func (r *postgresRepository) CreateRefund(ctx context.Context, refund *Refund) error {
	ctx, span := tracer.Start(ctx, "repository.CreateRefund")
	defer span.End()
	
	if refund.ID == uuid.Nil {
		refund.ID = uuid.New()
	}
	refund.CreatedAt = time.Now().UTC()
	refund.UpdatedAt = refund.CreatedAt
	
	query := `
		INSERT INTO refunds (
			id, tenant_id, payment_id, amount, currency, status, reason,
			reference, provider_reference, failure_reason, processed_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	
	_, err := r.pool.Exec(ctx, query,
		refund.ID, refund.TenantID, refund.PaymentID, refund.Amount, refund.Currency,
		refund.Status, refund.Reason, refund.Reference, refund.ProviderReference,
		refund.FailureReason, refund.ProcessedAt, refund.CreatedAt, refund.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("create refund: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetRefund(ctx context.Context, id uuid.UUID) (*Refund, error) {
	ctx, span := tracer.Start(ctx, "repository.GetRefund")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, payment_id, amount, currency, status, reason,
			reference, provider_reference, failure_reason, processed_at,
			created_at, updated_at
		FROM refunds WHERE id = $1`
	
	refund := &Refund{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&refund.ID, &refund.TenantID, &refund.PaymentID, &refund.Amount, &refund.Currency,
		&refund.Status, &refund.Reason, &refund.Reference, &refund.ProviderReference,
		&refund.FailureReason, &refund.ProcessedAt, &refund.CreatedAt, &refund.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get refund: %w", err)
	}
	
	return refund, nil
}

func (r *postgresRepository) GetRefundsByPayment(ctx context.Context, paymentID uuid.UUID) ([]*Refund, error) {
	ctx, span := tracer.Start(ctx, "repository.GetRefundsByPayment")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, payment_id, amount, currency, status, reason,
			reference, provider_reference, failure_reason, processed_at,
			created_at, updated_at
		FROM refunds WHERE payment_id = $1
		ORDER BY created_at DESC`
	
	rows, err := r.pool.Query(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get refunds by payment: %w", err)
	}
	defer rows.Close()
	
	refunds := make([]*Refund, 0)
	for rows.Next() {
		refund := &Refund{}
		err := rows.Scan(
			&refund.ID, &refund.TenantID, &refund.PaymentID, &refund.Amount, &refund.Currency,
			&refund.Status, &refund.Reason, &refund.Reference, &refund.ProviderReference,
			&refund.FailureReason, &refund.ProcessedAt, &refund.CreatedAt, &refund.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan refund: %w", err)
		}
		refunds = append(refunds, refund)
	}
	
	return refunds, nil
}

func (r *postgresRepository) UpdateRefundStatus(ctx context.Context, id uuid.UUID, status RefundStatus, providerRef *string) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateRefundStatus")
	defer span.End()
	
	now := time.Now().UTC()
	var processedAt *time.Time
	if status == RefundStatusProcessed || status == RefundStatusFailed {
		processedAt = &now
	}
	
	query := `
		UPDATE refunds SET
			status = $2,
			provider_reference = COALESCE($3, provider_reference),
			processed_at = COALESCE($4, processed_at),
			updated_at = $5
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, status, providerRef, processedAt, now)
	if err != nil {
		return fmt.Errorf("update refund status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

// ============================================================================
// Wallet Operations
// ============================================================================

func (r *postgresRepository) CreateWallet(ctx context.Context, wallet *Wallet) error {
	ctx, span := tracer.Start(ctx, "repository.CreateWallet")
	defer span.End()
	
	if wallet.ID == uuid.Nil {
		wallet.ID = uuid.New()
	}
	wallet.CreatedAt = time.Now().UTC()
	wallet.UpdatedAt = wallet.CreatedAt
	wallet.Version = 1
	wallet.Status = WalletStatusActive
	
	metadata, _ := json.Marshal(wallet.Metadata)
	
	query := `
		INSERT INTO wallets (
			id, tenant_id, customer_id, customer_type, name, balance, currency,
			status, freeze_reason, fineract_id, metadata, created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	
	_, err := r.pool.Exec(ctx, query,
		wallet.ID, wallet.TenantID, wallet.CustomerID, wallet.CustomerType,
		wallet.Name, wallet.Balance, wallet.Currency, wallet.Status,
		wallet.FreezeReason, wallet.FineractID, metadata,
		wallet.CreatedAt, wallet.UpdatedAt, wallet.Version,
	)
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateKey
		}
		return fmt.Errorf("create wallet: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetWallet(ctx context.Context, id uuid.UUID) (*Wallet, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWallet")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, customer_id, customer_type, name, balance, currency,
			status, freeze_reason, fineract_id, metadata, created_at, updated_at, version
		FROM wallets WHERE id = $1`
	
	wallet := &Wallet{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&wallet.ID, &wallet.TenantID, &wallet.CustomerID, &wallet.CustomerType,
		&wallet.Name, &wallet.Balance, &wallet.Currency, &wallet.Status,
		&wallet.FreezeReason, &wallet.FineractID, &metadata,
		&wallet.CreatedAt, &wallet.UpdatedAt, &wallet.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	
	json.Unmarshal(metadata, &wallet.Metadata)
	return wallet, nil
}

func (r *postgresRepository) GetWalletByCustomer(ctx context.Context, customerID uuid.UUID) (*Wallet, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWalletByCustomer")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, customer_id, customer_type, name, balance, currency,
			status, freeze_reason, fineract_id, metadata, created_at, updated_at, version
		FROM wallets WHERE customer_id = $1`
	
	wallet := &Wallet{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, customerID).Scan(
		&wallet.ID, &wallet.TenantID, &wallet.CustomerID, &wallet.CustomerType,
		&wallet.Name, &wallet.Balance, &wallet.Currency, &wallet.Status,
		&wallet.FreezeReason, &wallet.FineractID, &metadata,
		&wallet.CreatedAt, &wallet.UpdatedAt, &wallet.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get wallet by customer: %w", err)
	}
	
	json.Unmarshal(metadata, &wallet.Metadata)
	return wallet, nil
}

func (r *postgresRepository) UpdateWallet(ctx context.Context, wallet *Wallet) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateWallet")
	defer span.End()
	
	metadata, _ := json.Marshal(wallet.Metadata)
	wallet.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE wallets SET
			name = $3, status = $4, freeze_reason = $5,
			fineract_id = $6, metadata = $7, updated_at = $8,
			version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`
	
	err := r.pool.QueryRow(ctx, query,
		wallet.ID, wallet.Version, wallet.Name, wallet.Status,
		wallet.FreezeReason, wallet.FineractID, metadata, wallet.UpdatedAt,
	).Scan(&wallet.Version)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptimisticLock
		}
		return fmt.Errorf("update wallet: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) CreditWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error {
	ctx, span := tracer.Start(ctx, "repository.CreditWallet")
	defer span.End()
	
	query := `
		UPDATE wallets SET
			balance = balance + $2,
			updated_at = $3,
			version = version + 1
		WHERE id = $1 AND status = 'active'
		RETURNING balance`
	
	var newBalance decimal.Decimal
	err := r.pool.QueryRow(ctx, query, id, amount, time.Now().UTC()).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountFrozen
		}
		return fmt.Errorf("credit wallet: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) DebitWallet(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error {
	ctx, span := tracer.Start(ctx, "repository.DebitWallet")
	defer span.End()
	
	query := `
		UPDATE wallets SET
			balance = balance - $2,
			updated_at = $3,
			version = version + 1
		WHERE id = $1 AND status = 'active' AND balance >= $2
		RETURNING balance`
	
	var newBalance decimal.Decimal
	err := r.pool.QueryRow(ctx, query, id, amount, time.Now().UTC()).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Could be frozen or insufficient balance
			wallet, getErr := r.GetWallet(ctx, id)
			if getErr != nil {
				return getErr
			}
			if wallet.Status != WalletStatusActive {
				return ErrAccountFrozen
			}
			return ErrInsufficientBalance
		}
		return fmt.Errorf("debit wallet: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) FreezeWallet(ctx context.Context, id uuid.UUID, reason string) error {
	ctx, span := tracer.Start(ctx, "repository.FreezeWallet")
	defer span.End()
	
	query := `
		UPDATE wallets SET
			status = 'frozen',
			freeze_reason = $2,
			updated_at = $3
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, reason, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("freeze wallet: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) UnfreezeWallet(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "repository.UnfreezeWallet")
	defer span.End()
	
	query := `
		UPDATE wallets SET
			status = 'active',
			freeze_reason = NULL,
			updated_at = $2
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("unfreeze wallet: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

// ============================================================================
// Wallet Transaction Operations
// ============================================================================

func (r *postgresRepository) CreateWalletTransaction(ctx context.Context, tx *WalletTransaction) error {
	ctx, span := tracer.Start(ctx, "repository.CreateWalletTransaction")
	defer span.End()
	
	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}
	tx.CreatedAt = time.Now().UTC()
	
	metadata, _ := json.Marshal(tx.Metadata)
	
	query := `
		INSERT INTO wallet_transactions (
			id, wallet_id, type, amount, currency, balance_after, source,
			source_id, reference, description, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	
	_, err := r.pool.Exec(ctx, query,
		tx.ID, tx.WalletID, tx.Type, tx.Amount, tx.Currency, tx.Balance,
		tx.Source, tx.SourceID, tx.Reference, tx.Description, metadata, tx.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("create wallet transaction: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetWalletTransactions(ctx context.Context, walletID uuid.UUID, filter TransactionFilter) (*TransactionList, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWalletTransactions")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("wallet_id = $%d", argNum))
	args = append(args, walletID)
	argNum++
	
	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argNum))
		args = append(args, *filter.Type)
		argNum++
	}
	
	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argNum))
		args = append(args, *filter.FromDate)
		argNum++
	}
	
	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argNum))
		args = append(args, *filter.ToDate)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM wallet_transactions WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count wallet transactions: %w", err)
	}
	
	query := fmt.Sprintf(`
		SELECT id, wallet_id, type, amount, currency, balance_after, source,
			source_id, reference, description, metadata, created_at
		FROM wallet_transactions WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get wallet transactions: %w", err)
	}
	defer rows.Close()
	
	transactions := make([]*WalletTransaction, 0)
	for rows.Next() {
		tx := &WalletTransaction{}
		var metadata []byte
		
		err := rows.Scan(
			&tx.ID, &tx.WalletID, &tx.Type, &tx.Amount, &tx.Currency, &tx.Balance,
			&tx.Source, &tx.SourceID, &tx.Reference, &tx.Description, &metadata, &tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wallet transaction: %w", err)
		}
		
		json.Unmarshal(metadata, &tx.Metadata)
		transactions = append(transactions, tx)
	}
	
	return &TransactionList{
		Transactions: transactions,
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
	}, nil
}

// ============================================================================
// Credit Account Operations
// ============================================================================

func (r *postgresRepository) CreateCreditAccount(ctx context.Context, account *CreditAccount) error {
	ctx, span := tracer.Start(ctx, "repository.CreateCreditAccount")
	defer span.End()
	
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = account.CreatedAt
	account.Version = 1
	account.Status = CreditAccountStatusActive
	account.AvailableCredit = account.CreditLimit
	
	metadata, _ := json.Marshal(account.Metadata)
	
	query := `
		INSERT INTO credit_accounts (
			id, tenant_id, customer_id, customer_type, credit_limit, available_credit,
			outstanding_balance, currency, payment_terms_days, interest_rate,
			grace_period_days, billing_cycle, status, overdue_amount,
			last_payment_date, next_due_date, fineract_loan_id, metadata,
			created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)`
	
	_, err := r.pool.Exec(ctx, query,
		account.ID, account.TenantID, account.CustomerID, account.CustomerType,
		account.CreditLimit, account.AvailableCredit, account.OutstandingBalance,
		account.Currency, account.PaymentTermsDays, account.InterestRate,
		account.GracePeriodDays, account.BillingCycle, account.Status, account.OverdueAmount,
		account.LastPaymentDate, account.NextDueDate, account.FineractLoanID, metadata,
		account.CreatedAt, account.UpdatedAt, account.Version,
	)
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateKey
		}
		return fmt.Errorf("create credit account: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetCreditAccount(ctx context.Context, id uuid.UUID) (*CreditAccount, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditAccount")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, customer_id, customer_type, credit_limit, available_credit,
			outstanding_balance, currency, payment_terms_days, interest_rate,
			grace_period_days, billing_cycle, status, overdue_amount,
			last_payment_date, next_due_date, fineract_loan_id, metadata,
			created_at, updated_at, version
		FROM credit_accounts WHERE id = $1`
	
	account := &CreditAccount{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&account.ID, &account.TenantID, &account.CustomerID, &account.CustomerType,
		&account.CreditLimit, &account.AvailableCredit, &account.OutstandingBalance,
		&account.Currency, &account.PaymentTermsDays, &account.InterestRate,
		&account.GracePeriodDays, &account.BillingCycle, &account.Status, &account.OverdueAmount,
		&account.LastPaymentDate, &account.NextDueDate, &account.FineractLoanID, &metadata,
		&account.CreatedAt, &account.UpdatedAt, &account.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get credit account: %w", err)
	}
	
	json.Unmarshal(metadata, &account.Metadata)
	return account, nil
}

func (r *postgresRepository) GetCreditAccountByCustomer(ctx context.Context, customerID uuid.UUID) (*CreditAccount, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditAccountByCustomer")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, customer_id, customer_type, credit_limit, available_credit,
			outstanding_balance, currency, payment_terms_days, interest_rate,
			grace_period_days, billing_cycle, status, overdue_amount,
			last_payment_date, next_due_date, fineract_loan_id, metadata,
			created_at, updated_at, version
		FROM credit_accounts WHERE customer_id = $1`
	
	account := &CreditAccount{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, customerID).Scan(
		&account.ID, &account.TenantID, &account.CustomerID, &account.CustomerType,
		&account.CreditLimit, &account.AvailableCredit, &account.OutstandingBalance,
		&account.Currency, &account.PaymentTermsDays, &account.InterestRate,
		&account.GracePeriodDays, &account.BillingCycle, &account.Status, &account.OverdueAmount,
		&account.LastPaymentDate, &account.NextDueDate, &account.FineractLoanID, &metadata,
		&account.CreatedAt, &account.UpdatedAt, &account.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get credit account by customer: %w", err)
	}
	
	json.Unmarshal(metadata, &account.Metadata)
	return account, nil
}

func (r *postgresRepository) UpdateCreditAccount(ctx context.Context, account *CreditAccount) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateCreditAccount")
	defer span.End()
	
	metadata, _ := json.Marshal(account.Metadata)
	account.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE credit_accounts SET
			credit_limit = $3, available_credit = $4, outstanding_balance = $5,
			payment_terms_days = $6, interest_rate = $7, grace_period_days = $8,
			billing_cycle = $9, status = $10, overdue_amount = $11,
			last_payment_date = $12, next_due_date = $13, fineract_loan_id = $14,
			metadata = $15, updated_at = $16, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`
	
	err := r.pool.QueryRow(ctx, query,
		account.ID, account.Version, account.CreditLimit, account.AvailableCredit,
		account.OutstandingBalance, account.PaymentTermsDays, account.InterestRate,
		account.GracePeriodDays, account.BillingCycle, account.Status, account.OverdueAmount,
		account.LastPaymentDate, account.NextDueDate, account.FineractLoanID,
		metadata, account.UpdatedAt,
	).Scan(&account.Version)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptimisticLock
		}
		return fmt.Errorf("update credit account: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) UseCreditLine(ctx context.Context, id uuid.UUID, amount decimal.Decimal, orderID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "repository.UseCreditLine")
	defer span.End()
	
	query := `
		UPDATE credit_accounts SET
			available_credit = available_credit - $2,
			outstanding_balance = outstanding_balance + $2,
			updated_at = $3,
			version = version + 1
		WHERE id = $1 AND status = 'active' AND available_credit >= $2
		RETURNING available_credit`
	
	var newAvailable decimal.Decimal
	err := r.pool.QueryRow(ctx, query, id, amount, time.Now().UTC()).Scan(&newAvailable)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			account, getErr := r.GetCreditAccount(ctx, id)
			if getErr != nil {
				return getErr
			}
			if account.Status != CreditAccountStatusActive {
				return ErrAccountFrozen
			}
			return ErrCreditLimitExceeded
		}
		return fmt.Errorf("use credit line: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) RepayCreditLine(ctx context.Context, id uuid.UUID, amount decimal.Decimal, reference string) error {
	ctx, span := tracer.Start(ctx, "repository.RepayCreditLine")
	defer span.End()
	
	now := time.Now().UTC()
	
	query := `
		UPDATE credit_accounts SET
			available_credit = LEAST(available_credit + $2, credit_limit),
			outstanding_balance = GREATEST(outstanding_balance - $2, 0),
			overdue_amount = GREATEST(overdue_amount - $2, 0),
			last_payment_date = $3,
			updated_at = $3,
			version = version + 1
		WHERE id = $1
		RETURNING outstanding_balance`
	
	var newBalance decimal.Decimal
	err := r.pool.QueryRow(ctx, query, id, amount, now).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("repay credit line: %w", err)
	}
	
	return nil
}

// ============================================================================
// Credit Limit Request Operations
// ============================================================================

func (r *postgresRepository) CreateCreditLimitRequest(ctx context.Context, request *CreditLimitRequest) error {
	ctx, span := tracer.Start(ctx, "repository.CreateCreditLimitRequest")
	defer span.End()
	
	if request.ID == uuid.Nil {
		request.ID = uuid.New()
	}
	request.CreatedAt = time.Now().UTC()
	request.UpdatedAt = request.CreatedAt
	request.Status = "pending"
	
	query := `
		INSERT INTO credit_limit_requests (
			id, tenant_id, credit_account_id, current_limit, requested_limit,
			status, reason, reviewed_by, review_notes, reviewed_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	
	_, err := r.pool.Exec(ctx, query,
		request.ID, request.TenantID, request.CreditAccountID, request.CurrentLimit,
		request.RequestedLimit, request.Status, request.Reason, request.ReviewedBy,
		request.ReviewNotes, request.ReviewedAt, request.CreatedAt, request.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("create credit limit request: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetCreditLimitRequest(ctx context.Context, id uuid.UUID) (*CreditLimitRequest, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditLimitRequest")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, credit_account_id, current_limit, requested_limit,
			status, reason, reviewed_by, review_notes, reviewed_at,
			created_at, updated_at
		FROM credit_limit_requests WHERE id = $1`
	
	request := &CreditLimitRequest{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&request.ID, &request.TenantID, &request.CreditAccountID, &request.CurrentLimit,
		&request.RequestedLimit, &request.Status, &request.Reason, &request.ReviewedBy,
		&request.ReviewNotes, &request.ReviewedAt, &request.CreatedAt, &request.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get credit limit request: %w", err)
	}
	
	return request, nil
}

func (r *postgresRepository) UpdateCreditLimitRequest(ctx context.Context, request *CreditLimitRequest) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateCreditLimitRequest")
	defer span.End()
	
	request.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE credit_limit_requests SET
			status = $2, reviewed_by = $3, review_notes = $4,
			reviewed_at = $5, updated_at = $6
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query,
		request.ID, request.Status, request.ReviewedBy, request.ReviewNotes,
		request.ReviewedAt, request.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("update credit limit request: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) GetPendingCreditLimitRequests(ctx context.Context, tenantID uuid.UUID) ([]*CreditLimitRequest, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPendingCreditLimitRequests")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, credit_account_id, current_limit, requested_limit,
			status, reason, reviewed_by, review_notes, reviewed_at,
			created_at, updated_at
		FROM credit_limit_requests
		WHERE tenant_id = $1 AND status = 'pending'
		ORDER BY created_at ASC`
	
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("get pending credit limit requests: %w", err)
	}
	defer rows.Close()
	
	requests := make([]*CreditLimitRequest, 0)
	for rows.Next() {
		request := &CreditLimitRequest{}
		err := rows.Scan(
			&request.ID, &request.TenantID, &request.CreditAccountID, &request.CurrentLimit,
			&request.RequestedLimit, &request.Status, &request.Reason, &request.ReviewedBy,
			&request.ReviewNotes, &request.ReviewedAt, &request.CreatedAt, &request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan credit limit request: %w", err)
		}
		requests = append(requests, request)
	}
	
	return requests, nil
}

// ============================================================================
// Credit Transaction Operations
// ============================================================================

func (r *postgresRepository) CreateCreditTransaction(ctx context.Context, tx *CreditTransaction) error {
	ctx, span := tracer.Start(ctx, "repository.CreateCreditTransaction")
	defer span.End()
	
	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}
	tx.CreatedAt = time.Now().UTC()
	
	metadata, _ := json.Marshal(tx.Metadata)
	
	query := `
		INSERT INTO credit_transactions (
			id, credit_account_id, type, amount, currency, order_id,
			payment_method, reference, description, due_date, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	
	_, err := r.pool.Exec(ctx, query,
		tx.ID, tx.CreditAccountID, tx.Type, tx.Amount, tx.Currency, tx.OrderID,
		tx.PaymentMethod, tx.Reference, tx.Description, tx.DueDate, metadata, tx.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("create credit transaction: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetCreditTransactions(ctx context.Context, accountID uuid.UUID, filter TransactionFilter) (*CreditTransactionList, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditTransactions")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("credit_account_id = $%d", argNum))
	args = append(args, accountID)
	argNum++
	
	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argNum))
		args = append(args, *filter.Type)
		argNum++
	}
	
	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argNum))
		args = append(args, *filter.FromDate)
		argNum++
	}
	
	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argNum))
		args = append(args, *filter.ToDate)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM credit_transactions WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count credit transactions: %w", err)
	}
	
	query := fmt.Sprintf(`
		SELECT id, credit_account_id, type, amount, currency, order_id,
			payment_method, reference, description, due_date, metadata, created_at
		FROM credit_transactions WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get credit transactions: %w", err)
	}
	defer rows.Close()
	
	transactions := make([]*CreditTransaction, 0)
	for rows.Next() {
		tx := &CreditTransaction{}
		var metadata []byte
		
		err := rows.Scan(
			&tx.ID, &tx.CreditAccountID, &tx.Type, &tx.Amount, &tx.Currency, &tx.OrderID,
			&tx.PaymentMethod, &tx.Reference, &tx.Description, &tx.DueDate, &metadata, &tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan credit transaction: %w", err)
		}
		
		json.Unmarshal(metadata, &tx.Metadata)
		transactions = append(transactions, tx)
	}
	
	return &CreditTransactionList{
		Transactions: transactions,
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
	}, nil
}

func (r *postgresRepository) GetCreditStatement(ctx context.Context, accountID uuid.UUID, from, to time.Time) (*CreditStatement, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditStatement")
	defer span.End()
	
	// Get opening balance
	openingQuery := `
		SELECT COALESCE(SUM(CASE WHEN type = 'usage' THEN amount ELSE -amount END), 0)
		FROM credit_transactions
		WHERE credit_account_id = $1 AND created_at < $2`
	
	var openingBalance decimal.Decimal
	err := r.pool.QueryRow(ctx, openingQuery, accountID, from).Scan(&openingBalance)
	if err != nil {
		return nil, fmt.Errorf("get opening balance: %w", err)
	}
	
	// Get transactions in period
	txList, err := r.GetCreditTransactions(ctx, accountID, TransactionFilter{
		FromDate: &from,
		ToDate:   &to,
		Limit:    10000,
		Offset:   0,
	})
	if err != nil {
		return nil, err
	}
	
	statement := &CreditStatement{
		AccountID:      accountID,
		FromDate:       from,
		ToDate:         to,
		OpeningBalance: openingBalance,
		Transactions:   txList.Transactions,
	}
	
	closingBalance := openingBalance
	for _, tx := range txList.Transactions {
		switch tx.Type {
		case "usage":
			statement.TotalUsage = statement.TotalUsage.Add(tx.Amount)
			closingBalance = closingBalance.Add(tx.Amount)
		case "repayment":
			statement.TotalRepayments = statement.TotalRepayments.Add(tx.Amount)
			closingBalance = closingBalance.Sub(tx.Amount)
		case "fee":
			statement.TotalFees = statement.TotalFees.Add(tx.Amount)
			closingBalance = closingBalance.Add(tx.Amount)
		case "interest":
			statement.TotalInterest = statement.TotalInterest.Add(tx.Amount)
			closingBalance = closingBalance.Add(tx.Amount)
		}
	}
	
	statement.ClosingBalance = closingBalance
	
	return statement, nil
}

// ============================================================================
// Credit Score Operations
// ============================================================================

func (r *postgresRepository) GetCreditScore(ctx context.Context, customerID uuid.UUID) (*CreditScore, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditScore")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, customer_id, score, rating, payment_history_score,
			order_frequency_score, credit_utilization_score, account_age_score,
			last_calculated_at, created_at, updated_at
		FROM credit_scores WHERE customer_id = $1`
	
	score := &CreditScore{}
	err := r.pool.QueryRow(ctx, query, customerID).Scan(
		&score.ID, &score.TenantID, &score.CustomerID, &score.Score, &score.Rating,
		&score.PaymentHistory, &score.OrderFrequency, &score.CreditUtilization,
		&score.AccountAge, &score.LastCalculatedAt, &score.CreatedAt, &score.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get credit score: %w", err)
	}
	
	return score, nil
}

func (r *postgresRepository) UpsertCreditScore(ctx context.Context, score *CreditScore) error {
	ctx, span := tracer.Start(ctx, "repository.UpsertCreditScore")
	defer span.End()
	
	if score.ID == uuid.Nil {
		score.ID = uuid.New()
	}
	score.UpdatedAt = time.Now().UTC()
	if score.CreatedAt.IsZero() {
		score.CreatedAt = score.UpdatedAt
	}
	score.LastCalculatedAt = score.UpdatedAt
	
	query := `
		INSERT INTO credit_scores (
			id, tenant_id, customer_id, score, rating, payment_history_score,
			order_frequency_score, credit_utilization_score, account_age_score,
			last_calculated_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (customer_id) DO UPDATE SET
			score = EXCLUDED.score,
			rating = EXCLUDED.rating,
			payment_history_score = EXCLUDED.payment_history_score,
			order_frequency_score = EXCLUDED.order_frequency_score,
			credit_utilization_score = EXCLUDED.credit_utilization_score,
			account_age_score = EXCLUDED.account_age_score,
			last_calculated_at = EXCLUDED.last_calculated_at,
			updated_at = EXCLUDED.updated_at`
	
	_, err := r.pool.Exec(ctx, query,
		score.ID, score.TenantID, score.CustomerID, score.Score, score.Rating,
		score.PaymentHistory, score.OrderFrequency, score.CreditUtilization,
		score.AccountAge, score.LastCalculatedAt, score.CreatedAt, score.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("upsert credit score: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetCreditScoreFactors(ctx context.Context, customerID uuid.UUID) (*CreditScoreFactors, error) {
	ctx, span := tracer.Start(ctx, "repository.GetCreditScoreFactors")
	defer span.End()
	
	// This would typically aggregate data from multiple tables
	// For now, returning placeholder implementation
	factors := &CreditScoreFactors{CustomerID: customerID}
	
	// Get order stats
	orderQuery := `
		SELECT COUNT(*), COALESCE(AVG(total_amount), 0), COALESCE(SUM(total_amount), 0)
		FROM orders WHERE customer_id = $1 AND status = 'completed'`
	
	err := r.pool.QueryRow(ctx, orderQuery, customerID).Scan(
		&factors.TotalOrders, &factors.AverageOrderValue, &factors.TotalSpent,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get order stats: %w", err)
	}
	
	// Get credit account info
	creditQuery := `
		SELECT credit_limit, outstanding_balance, 
			EXTRACT(DAY FROM NOW() - created_at)::int as age_days
		FROM credit_accounts WHERE customer_id = $1`
	
	err = r.pool.QueryRow(ctx, creditQuery, customerID).Scan(
		&factors.CreditLimit, &factors.CurrentBalance, &factors.AccountAgeDays,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get credit account: %w", err)
	}
	
	return factors, nil
}

// ============================================================================
// Disbursement Operations (Simplified)
// ============================================================================

func (r *postgresRepository) CreateDisbursement(ctx context.Context, d *Disbursement) error {
	ctx, span := tracer.Start(ctx, "repository.CreateDisbursement")
	defer span.End()
	
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	d.CreatedAt = time.Now().UTC()
	d.UpdatedAt = d.CreatedAt
	
	metadata, _ := json.Marshal(d.Metadata)
	
	query := `
		INSERT INTO disbursements (
			id, tenant_id, batch_id, recipient_id, recipient_type, amount, currency,
			disbursement_type, status, payout_method, bank_account_id, mobile_number,
			mobile_network, wallet_id, reference, provider_reference, description,
			failure_reason, scheduled_for, processed_at, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23
		)`
	
	_, err := r.pool.Exec(ctx, query,
		d.ID, d.TenantID, d.BatchID, d.RecipientID, d.RecipientType, d.Amount, d.Currency,
		d.DisbursementType, d.Status, d.PayoutMethod, d.BankAccountID, d.MobileNumber,
		d.MobileNetwork, d.WalletID, d.Reference, d.ProviderReference, d.Description,
		d.FailureReason, d.ScheduledFor, d.ProcessedAt, metadata, d.CreatedAt, d.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("create disbursement: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetDisbursement(ctx context.Context, id uuid.UUID) (*Disbursement, error) {
	ctx, span := tracer.Start(ctx, "repository.GetDisbursement")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, batch_id, recipient_id, recipient_type, amount, currency,
			disbursement_type, status, payout_method, bank_account_id, mobile_number,
			mobile_network, wallet_id, reference, provider_reference, description,
			failure_reason, scheduled_for, processed_at, metadata, created_at, updated_at
		FROM disbursements WHERE id = $1`
	
	d := &Disbursement{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.TenantID, &d.BatchID, &d.RecipientID, &d.RecipientType, &d.Amount, &d.Currency,
		&d.DisbursementType, &d.Status, &d.PayoutMethod, &d.BankAccountID, &d.MobileNumber,
		&d.MobileNetwork, &d.WalletID, &d.Reference, &d.ProviderReference, &d.Description,
		&d.FailureReason, &d.ScheduledFor, &d.ProcessedAt, &metadata, &d.CreatedAt, &d.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get disbursement: %w", err)
	}
	
	json.Unmarshal(metadata, &d.Metadata)
	return d, nil
}

func (r *postgresRepository) UpdateDisbursement(ctx context.Context, d *Disbursement) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateDisbursement")
	defer span.End()
	
	metadata, _ := json.Marshal(d.Metadata)
	d.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE disbursements SET
			status = $2, provider_reference = $3, failure_reason = $4,
			processed_at = $5, metadata = $6, updated_at = $7
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query,
		d.ID, d.Status, d.ProviderReference, d.FailureReason,
		d.ProcessedAt, metadata, d.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("update disbursement: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) ListDisbursements(ctx context.Context, filter DisbursementFilter) (*DisbursementList, error) {
	ctx, span := tracer.Start(ctx, "repository.ListDisbursements")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++
	
	if filter.RecipientID != nil {
		conditions = append(conditions, fmt.Sprintf("recipient_id = $%d", argNum))
		args = append(args, *filter.RecipientID)
		argNum++
	}
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, s)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM disbursements WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count disbursements: %w", err)
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, batch_id, recipient_id, recipient_type, amount, currency,
			disbursement_type, status, payout_method, bank_account_id, mobile_number,
			mobile_network, wallet_id, reference, provider_reference, description,
			failure_reason, scheduled_for, processed_at, metadata, created_at, updated_at
		FROM disbursements WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list disbursements: %w", err)
	}
	defer rows.Close()
	
	disbursements := make([]*Disbursement, 0)
	for rows.Next() {
		d := &Disbursement{}
		var metadata []byte
		
		err := rows.Scan(
			&d.ID, &d.TenantID, &d.BatchID, &d.RecipientID, &d.RecipientType, &d.Amount, &d.Currency,
			&d.DisbursementType, &d.Status, &d.PayoutMethod, &d.BankAccountID, &d.MobileNumber,
			&d.MobileNetwork, &d.WalletID, &d.Reference, &d.ProviderReference, &d.Description,
			&d.FailureReason, &d.ScheduledFor, &d.ProcessedAt, &metadata, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan disbursement: %w", err)
		}
		
		json.Unmarshal(metadata, &d.Metadata)
		disbursements = append(disbursements, d)
	}
	
	return &DisbursementList{
		Disbursements: disbursements,
		Total:         total,
		Limit:         filter.Limit,
		Offset:        filter.Offset,
	}, nil
}

// ============================================================================
// Bank Account Operations (Stub)
// ============================================================================

func (r *postgresRepository) CreateBankAccount(ctx context.Context, account *BankAccount) error {
	return nil // Stub
}

func (r *postgresRepository) GetBankAccount(ctx context.Context, id uuid.UUID) (*BankAccount, error) {
	return nil, ErrNotFound // Stub
}

func (r *postgresRepository) GetBankAccountsByCustomer(ctx context.Context, customerID uuid.UUID) ([]*BankAccount, error) {
	return []*BankAccount{}, nil // Stub
}

func (r *postgresRepository) UpdateBankAccount(ctx context.Context, account *BankAccount) error {
	return nil // Stub
}

func (r *postgresRepository) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	return nil // Stub
}

func (r *postgresRepository) SetDefaultBankAccount(ctx context.Context, customerID, accountID uuid.UUID) error {
	return nil // Stub
}

// ============================================================================
// Virtual Account Operations (Stub)
// ============================================================================

func (r *postgresRepository) CreateVirtualAccount(ctx context.Context, account *VirtualAccount) error {
	return nil // Stub
}

func (r *postgresRepository) GetVirtualAccount(ctx context.Context, id uuid.UUID) (*VirtualAccount, error) {
	return nil, ErrNotFound // Stub
}

func (r *postgresRepository) GetVirtualAccountByNumber(ctx context.Context, accountNumber string) (*VirtualAccount, error) {
	return nil, ErrNotFound // Stub
}

func (r *postgresRepository) GetVirtualAccountByCustomer(ctx context.Context, customerID uuid.UUID) (*VirtualAccount, error) {
	return nil, ErrNotFound // Stub
}

func (r *postgresRepository) UpdateVirtualAccount(ctx context.Context, account *VirtualAccount) error {
	return nil // Stub
}

// ============================================================================
// Batch Operations (Stub)
// ============================================================================

func (r *postgresRepository) CreateDisbursementBatch(ctx context.Context, batch *DisbursementBatch) error {
	return nil // Stub
}

func (r *postgresRepository) GetDisbursementBatch(ctx context.Context, id uuid.UUID) (*DisbursementBatch, error) {
	return nil, ErrNotFound // Stub
}

func (r *postgresRepository) UpdateDisbursementBatch(ctx context.Context, batch *DisbursementBatch) error {
	return nil // Stub
}

// ============================================================================
// Reconciliation Operations (Stub)
// ============================================================================

func (r *postgresRepository) CreateReconciliationRecord(ctx context.Context, record *ReconciliationRecord) error {
	return nil // Stub
}

func (r *postgresRepository) GetReconciliationReport(ctx context.Context, filter ReconciliationFilter) (*ReconciliationReport, error) {
	return &ReconciliationReport{}, nil // Stub
}

func (r *postgresRepository) FlagPaymentForReview(ctx context.Context, paymentID uuid.UUID, reason string) error {
	return nil // Stub
}

// ============================================================================
// Analytics Operations
// ============================================================================

func (r *postgresRepository) GetPaymentStats(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*PaymentStats, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPaymentStats")
	defer span.End()
	
	query := `
		SELECT 
			COUNT(*) as total,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(fee), 0) as total_fees,
			COUNT(*) FILTER (WHERE status = 'completed') as completed
		FROM payments
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3`
	
	stats := &PaymentStats{
		TenantID:   tenantID,
		FromDate:   from,
		ToDate:     to,
		ByStatus:   make(map[PaymentStatus]int64),
		ByMethod:   make(map[PaymentMethod]int64),
		ByProvider: make(map[string]int64),
	}
	
	var completed int64
	err := r.pool.QueryRow(ctx, query, tenantID, from, to).Scan(
		&stats.TotalPayments, &stats.TotalAmount, &stats.TotalFees, &completed,
	)
	if err != nil {
		return nil, fmt.Errorf("get payment stats: %w", err)
	}
	
	if stats.TotalPayments > 0 {
		stats.SuccessRate = float64(completed) / float64(stats.TotalPayments) * 100
	}
	
	return stats, nil
}

func (r *postgresRepository) GetRevenueReport(ctx context.Context, tenantID uuid.UUID, from, to time.Time, groupBy string) ([]*RevenueEntry, error) {
	ctx, span := tracer.Start(ctx, "repository.GetRevenueReport")
	defer span.End()
	
	var timeFormat string
	switch groupBy {
	case "week":
		timeFormat = "IYYY-IW"
	case "month":
		timeFormat = "YYYY-MM"
	default:
		timeFormat = "YYYY-MM-DD"
	}
	
	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(created_at, '%s') as period,
			COALESCE(SUM(amount), 0) as gross,
			COALESCE(SUM(refunded_amount), 0) as refunds,
			COALESCE(SUM(amount) - SUM(refunded_amount), 0) as net,
			COALESCE(SUM(fee), 0) as fees,
			COUNT(*) as count
		FROM payments
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3 AND status = 'completed'
		GROUP BY period
		ORDER BY period ASC`, timeFormat)
	
	rows, err := r.pool.Query(ctx, query, tenantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get revenue report: %w", err)
	}
	defer rows.Close()
	
	entries := make([]*RevenueEntry, 0)
	for rows.Next() {
		entry := &RevenueEntry{}
		err := rows.Scan(&entry.Period, &entry.Gross, &entry.Refunds, &entry.Net, &entry.Fees, &entry.Count)
		if err != nil {
			return nil, fmt.Errorf("scan revenue entry: %w", err)
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// ============================================================================
// Transaction Support
// ============================================================================

func (r *postgresRepository) WithTx(ctx context.Context, fn func(Repository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()
	
	txRepo := &postgresRepository{pool: r.pool, tx: tx}
	
	if err := fn(txRepo); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	
	return nil
}
