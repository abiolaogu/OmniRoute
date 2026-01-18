// Package domain contains repository interfaces for the Bank Gateway service.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// BankConnectionRepository defines operations for bank connection persistence
type BankConnectionRepository interface {
	// FindByID retrieves a bank connection by ID
	FindByID(ctx context.Context, tenantID, connectionID uuid.UUID) (*BankConnection, error)

	// FindByBankCode retrieves a bank connection by bank code
	FindByBankCode(ctx context.Context, tenantID uuid.UUID, bankCode string) (*BankConnection, error)

	// FindByType retrieves bank connections by type
	FindByType(ctx context.Context, tenantID uuid.UUID, connType BankConnectionType) ([]*BankConnection, error)

	// FindActive retrieves all active bank connections
	FindActive(ctx context.Context, tenantID uuid.UUID) ([]*BankConnection, error)

	// Save persists a bank connection
	Save(ctx context.Context, connection *BankConnection) error

	// Update updates a bank connection
	Update(ctx context.Context, connection *BankConnection) error

	// UpdateStatus updates connection status
	UpdateStatus(ctx context.Context, connectionID uuid.UUID, status ConnectionStatus) error
}

// PaymentTransactionRepository defines operations for payment persistence
type PaymentTransactionRepository interface {
	// FindByID retrieves a payment by ID
	FindByID(ctx context.Context, paymentID uuid.UUID) (*PaymentTransaction, error)

	// FindByReference retrieves a payment by reference
	FindByReference(ctx context.Context, tenantID uuid.UUID, reference string) (*PaymentTransaction, error)

	// FindByExternalReference retrieves a payment by external reference
	FindByExternalReference(ctx context.Context, externalRef string) (*PaymentTransaction, error)

	// FindByStatus retrieves payments by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status PaymentStatus, limit, offset int) ([]*PaymentTransaction, error)

	// FindByInitiator retrieves payments initiated by a user
	FindByInitiator(ctx context.Context, tenantID, userID uuid.UUID, from, to time.Time) ([]*PaymentTransaction, error)

	// FindPending retrieves pending payments
	FindPending(ctx context.Context, limit int) ([]*PaymentTransaction, error)

	// Save persists a payment transaction
	Save(ctx context.Context, payment *PaymentTransaction) error

	// Update updates a payment transaction
	Update(ctx context.Context, payment *PaymentTransaction) error

	// UpdateStatus updates payment status
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status PaymentStatus, bankResponse map[string]interface{}) error
}

// VirtualAccountRepository defines operations for virtual account persistence
type VirtualAccountRepository interface {
	// FindByID retrieves a virtual account by ID
	FindByID(ctx context.Context, accountID uuid.UUID) (*VirtualAccount, error)

	// FindByAccountNumber retrieves a virtual account by account number
	FindByAccountNumber(ctx context.Context, accountNumber string) (*VirtualAccount, error)

	// FindByOwner retrieves virtual accounts for an owner
	FindByOwner(ctx context.Context, tenantID, ownerID uuid.UUID, ownerType AccountOwnerType) ([]*VirtualAccount, error)

	// FindActive retrieves active virtual accounts
	FindActive(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*VirtualAccount, error)

	// Save persists a virtual account
	Save(ctx context.Context, account *VirtualAccount) error

	// Update updates a virtual account
	Update(ctx context.Context, account *VirtualAccount) error

	// UpdateBalance updates account balance
	UpdateBalance(ctx context.Context, accountID uuid.UUID, newBalance interface{}) error

	// UpdateStatus updates account status
	UpdateStatus(ctx context.Context, accountID uuid.UUID, status AccountStatus) error
}

// BulkPaymentBatchRepository defines operations for bulk payment batch persistence
type BulkPaymentBatchRepository interface {
	// FindByID retrieves a batch by ID
	FindByID(ctx context.Context, batchID uuid.UUID) (*BulkPaymentBatch, error)

	// FindByReference retrieves a batch by reference
	FindByReference(ctx context.Context, tenantID uuid.UUID, reference string) (*BulkPaymentBatch, error)

	// FindByStatus retrieves batches by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status BatchStatus, limit, offset int) ([]*BulkPaymentBatch, error)

	// FindByInitiator retrieves batches initiated by a user
	FindByInitiator(ctx context.Context, tenantID, userID uuid.UUID) ([]*BulkPaymentBatch, error)

	// Save persists a batch
	Save(ctx context.Context, batch *BulkPaymentBatch) error

	// Update updates a batch
	Update(ctx context.Context, batch *BulkPaymentBatch) error

	// UpdateCounts updates success/fail counts
	UpdateCounts(ctx context.Context, batchID uuid.UUID, successCount, failCount int) error
}

// ReconciliationRepository defines operations for reconciliation
type ReconciliationRepository interface {
	// GetUnreconciledPayments retrieves payments not yet reconciled
	GetUnreconciledPayments(ctx context.Context, tenantID uuid.UUID, from, to time.Time) ([]*PaymentTransaction, error)

	// MarkReconciled marks a payment as reconciled
	MarkReconciled(ctx context.Context, paymentID uuid.UUID, statementRef string) error

	// GetReconciliationSummary retrieves reconciliation summary
	GetReconciliationSummary(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*ReconciliationSummary, error)
}

// ReconciliationSummary represents a reconciliation report summary
type ReconciliationSummary struct {
	TotalTransactions int                  `json:"total_transactions"`
	TotalAmount       interface{}          `json:"total_amount"`
	MatchedCount      int                  `json:"matched_count"`
	UnmatchedCount    int                  `json:"unmatched_count"`
	Discrepancies     []ReconciliationItem `json:"discrepancies,omitempty"`
}

// ReconciliationItem represents a single reconciliation item
type ReconciliationItem struct {
	TransactionID  uuid.UUID   `json:"transaction_id"`
	ExpectedAmount interface{} `json:"expected_amount"`
	ActualAmount   interface{} `json:"actual_amount"`
	Difference     interface{} `json:"difference"`
}
