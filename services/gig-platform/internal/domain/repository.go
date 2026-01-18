// Package domain contains repository interfaces for the Gig Platform service.
// Following DDD principles, repository interfaces are defined in the domain layer.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// WorkerRepository defines operations for worker persistence
type WorkerRepository interface {
	// FindByID retrieves a worker by ID
	FindByID(ctx context.Context, tenantID, workerID uuid.UUID) (*GigWorker, error)

	// FindByUserID retrieves a worker by user ID
	FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*GigWorker, error)

	// FindAvailable retrieves available workers
	FindAvailable(ctx context.Context, tenantID uuid.UUID, workerType WorkerType, limit int) ([]*GigWorker, error)

	// FindNearLocation retrieves workers near a location
	FindNearLocation(ctx context.Context, tenantID uuid.UUID, lat, lng float64, radiusKm float64, limit int) ([]*GigWorker, error)

	// FindByStatus retrieves workers by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status WorkerStatus, limit, offset int) ([]*GigWorker, error)

	// Save persists a worker
	Save(ctx context.Context, worker *GigWorker) error

	// Update updates a worker
	Update(ctx context.Context, worker *GigWorker) error

	// UpdateStatus updates worker status
	UpdateStatus(ctx context.Context, workerID uuid.UUID, status WorkerStatus) error

	// UpdateLocation updates worker location
	UpdateLocation(ctx context.Context, workerID uuid.UUID, location Location) error
}

// TaskRepository defines operations for task persistence
type TaskRepository interface {
	// FindByID retrieves a task by ID
	FindByID(ctx context.Context, tenantID, taskID uuid.UUID) (*Task, error)

	// FindByOrderID retrieves tasks for an order
	FindByOrderID(ctx context.Context, tenantID, orderID uuid.UUID) ([]*Task, error)

	// FindByWorkerID retrieves tasks assigned to a worker
	FindByWorkerID(ctx context.Context, tenantID, workerID uuid.UUID, statuses []TaskStatus) ([]*Task, error)

	// FindPending retrieves pending tasks
	FindPending(ctx context.Context, tenantID uuid.UUID, taskType TaskType, limit int) ([]*Task, error)

	// FindOverdue retrieves overdue tasks
	FindOverdue(ctx context.Context, tenantID uuid.UUID) ([]*Task, error)

	// Save persists a task
	Save(ctx context.Context, task *Task) error

	// Update updates a task
	Update(ctx context.Context, task *Task) error

	// UpdateStatus updates task status
	UpdateStatus(ctx context.Context, taskID uuid.UUID, status TaskStatus) error

	// AssignWorker assigns a worker to a task
	AssignWorker(ctx context.Context, taskID, workerID uuid.UUID) error
}

// AllocationRepository defines operations for allocation persistence
type AllocationRepository interface {
	// FindByID retrieves an allocation by ID
	FindByID(ctx context.Context, allocationID uuid.UUID) (*Allocation, error)

	// FindByTaskID retrieves allocations for a task
	FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*Allocation, error)

	// FindByWorkerID retrieves allocations for a worker
	FindByWorkerID(ctx context.Context, workerID uuid.UUID, status AllocationStatus) ([]*Allocation, error)

	// FindPendingOffers retrieves pending offers for a worker
	FindPendingOffers(ctx context.Context, workerID uuid.UUID) ([]*TaskOffer, error)

	// Save persists an allocation
	Save(ctx context.Context, allocation *Allocation) error

	// Update updates an allocation
	Update(ctx context.Context, allocation *Allocation) error

	// SaveOffer persists a task offer
	SaveOffer(ctx context.Context, offer *TaskOffer) error

	// UpdateOfferStatus updates offer status
	UpdateOfferStatus(ctx context.Context, offerID uuid.UUID, status OfferStatus) error
}

// EarningRepository defines operations for earning persistence
type EarningRepository interface {
	// FindByWorkerID retrieves earnings for a worker
	FindByWorkerID(ctx context.Context, workerID uuid.UUID, from, to time.Time) ([]*Earning, error)

	// FindByTaskID retrieves earnings for a task
	FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*Earning, error)

	// Save persists an earning
	Save(ctx context.Context, earning *Earning) error

	// SumByWorker calculates total earnings for a worker in a period
	SumByWorker(ctx context.Context, workerID uuid.UUID, from, to time.Time) (interface{}, error)
}

// PayoutRepository defines operations for payout persistence
type PayoutRepository interface {
	// FindByID retrieves a payout by ID
	FindByID(ctx context.Context, payoutID uuid.UUID) (*Payout, error)

	// FindByWorkerID retrieves payouts for a worker
	FindByWorkerID(ctx context.Context, workerID uuid.UUID, limit, offset int) ([]*Payout, error)

	// FindPending retrieves pending payouts
	FindPending(ctx context.Context, limit int) ([]*Payout, error)

	// Save persists a payout
	Save(ctx context.Context, payout *Payout) error

	// Update updates a payout
	Update(ctx context.Context, payout *Payout) error
}

// TaskProofRepository defines operations for task proof persistence
type TaskProofRepository interface {
	// FindByTaskID retrieves proofs for a task
	FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]*TaskProof, error)

	// Save persists a task proof
	Save(ctx context.Context, proof *TaskProof) error
}

// LocationHistoryRepository defines operations for location history
type LocationHistoryRepository interface {
	// Save persists a location history entry
	Save(ctx context.Context, workerID uuid.UUID, location Location, timestamp time.Time) error

	// FindByWorker retrieves location history for a worker
	FindByWorker(ctx context.Context, workerID uuid.UUID, from, to time.Time) ([]LocationHistory, error)
}

// LocationHistory represents a historical location entry
type LocationHistory struct {
	WorkerID  uuid.UUID `json:"worker_id"`
	Location  Location  `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}
