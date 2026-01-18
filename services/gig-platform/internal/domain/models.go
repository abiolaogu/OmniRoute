// Package domain contains the core domain models for the Gig Platform service.
// Following DDD principles with aggregates, entities, and value objects.
package domain

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Value Objects
// ============================================================================

// WorkerStatus represents the availability status of a gig worker
type WorkerStatus string

const (
	WorkerStatusAvailable WorkerStatus = "available"
	WorkerStatusBusy      WorkerStatus = "busy"
	WorkerStatusOffline   WorkerStatus = "offline"
	WorkerStatusOnBreak   WorkerStatus = "on_break"
	WorkerStatusInTransit WorkerStatus = "in_transit"
)

// WorkerType represents the type of work a worker can perform
type WorkerType string

const (
	WorkerTypeDelivery   WorkerType = "delivery"
	WorkerTypeSales      WorkerType = "sales"
	WorkerTypeCollection WorkerType = "collection"
	WorkerTypeAudit      WorkerType = "audit"
	WorkerTypeMultiRole  WorkerType = "multi_role"
)

// WorkerLevel represents the career progression level
type WorkerLevel string

const (
	WorkerLevelStarter WorkerLevel = "starter"
	WorkerLevelBronze  WorkerLevel = "bronze"
	WorkerLevelSilver  WorkerLevel = "silver"
	WorkerLevelGold    WorkerLevel = "gold"
	WorkerLevelDiamond WorkerLevel = "diamond"
	WorkerLevelMaster  WorkerLevel = "master"
)

// TaskStatus represents the lifecycle status of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusAccepted   TaskStatus = "accepted"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// TaskType represents the category of task
type TaskType string

const (
	TaskTypeDelivery   TaskType = "delivery"
	TaskTypePickup     TaskType = "pickup"
	TaskTypeCollection TaskType = "collection"
	TaskTypeSalesVisit TaskType = "sales_visit"
	TaskTypeAudit      TaskType = "audit"
	TaskTypeCustom     TaskType = "custom"
)

// AllocationStatus represents the status of a task allocation offer
type AllocationStatus string

const (
	AllocationStatusOffered   AllocationStatus = "offered"
	AllocationStatusAccepted  AllocationStatus = "accepted"
	AllocationStatusRejected  AllocationStatus = "rejected"
	AllocationStatusExpired   AllocationStatus = "expired"
	AllocationStatusCancelled AllocationStatus = "cancelled"
)

// OfferStatus for task offers
type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "pending"
	OfferStatusAccepted OfferStatus = "accepted"
	OfferStatusRejected OfferStatus = "rejected"
	OfferStatusExpired  OfferStatus = "expired"
)

// Location represents a geographical location (value object)
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
	City      string  `json:"city,omitempty"`
	State     string  `json:"state,omitempty"`
	Country   string  `json:"country,omitempty"`
}

// WorkerAvailability represents availability status with optional reason
type WorkerAvailability struct {
	Status    WorkerStatus `json:"status"`
	Reason    string       `json:"reason,omitempty"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ============================================================================
// Aggregates and Entities
// ============================================================================

// GigWorker is the aggregate root for worker management
type GigWorker struct {
	ID       uuid.UUID    `json:"id"`
	TenantID uuid.UUID    `json:"tenant_id"`
	UserID   uuid.UUID    `json:"user_id"`
	Type     WorkerType   `json:"type"`
	Level    WorkerLevel  `json:"level"`
	Status   WorkerStatus `json:"status"`

	// Personal info
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`

	// Location
	CurrentLocation Location `json:"current_location"`
	HomeLocation    Location `json:"home_location"`

	// Transport
	VehicleType  string `json:"vehicle_type"`
	VehiclePlate string `json:"vehicle_plate"`

	// Performance
	Rating         decimal.Decimal `json:"rating"`
	TotalTasks     int             `json:"total_tasks"`
	CompletedTasks int             `json:"completed_tasks"`
	SuccessRate    decimal.Decimal `json:"success_rate"`

	// Earnings
	TotalEarnings decimal.Decimal `json:"total_earnings"`
	WalletBalance decimal.Decimal `json:"wallet_balance"`

	// Verification
	IsVerified bool       `json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task is the aggregate root for task management
type Task struct {
	ID       uuid.UUID  `json:"id"`
	TenantID uuid.UUID  `json:"tenant_id"`
	Type     TaskType   `json:"type"`
	Status   TaskStatus `json:"status"`
	Priority int        `json:"priority"`

	// Assignment
	AssignedWorkerID *uuid.UUID `json:"assigned_worker_id,omitempty"`
	AssignedAt       *time.Time `json:"assigned_at,omitempty"`

	// Locations
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`

	// Order reference
	OrderID    uuid.UUID `json:"order_id"`
	CustomerID uuid.UUID `json:"customer_id"`

	// Task details
	Description  string     `json:"description"`
	Instructions string     `json:"instructions"`
	Items        []TaskItem `json:"items"`

	// Collection (for cash collection tasks)
	CollectionAmount decimal.Decimal `json:"collection_amount"`
	CollectedAmount  decimal.Decimal `json:"collected_amount"`

	// Pricing
	BasePayout      decimal.Decimal `json:"base_payout"`
	BonusPayout     decimal.Decimal `json:"bonus_payout"`
	SurgeMultiplier decimal.Decimal `json:"surge_multiplier"`

	// Timing
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskItem represents an item in a task
type TaskItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	Weight    float64   `json:"weight"`
}

// TaskOffer represents an offer made to a worker
type TaskOffer struct {
	ID           uuid.UUID       `json:"id"`
	TaskID       uuid.UUID       `json:"task_id"`
	WorkerID     uuid.UUID       `json:"worker_id"`
	Status       OfferStatus     `json:"status"`
	OfferedAt    time.Time       `json:"offered_at"`
	ExpiresAt    time.Time       `json:"expires_at"`
	RespondedAt  *time.Time      `json:"responded_at,omitempty"`
	PayoutAmount decimal.Decimal `json:"payout_amount"`
}

// Allocation represents a task assignment
type Allocation struct {
	ID       uuid.UUID        `json:"id"`
	TenantID uuid.UUID        `json:"tenant_id"`
	TaskID   uuid.UUID        `json:"task_id"`
	WorkerID uuid.UUID        `json:"worker_id"`
	Status   AllocationStatus `json:"status"`

	// Payout details
	BasePayout  decimal.Decimal `json:"base_payout"`
	BonusPayout decimal.Decimal `json:"bonus_payout"`
	TotalPayout decimal.Decimal `json:"total_payout"`

	// Timestamps
	AllocatedAt time.Time  `json:"allocated_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Earning represents worker earnings for completed tasks
type Earning struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	WorkerID    uuid.UUID       `json:"worker_id"`
	TaskID      uuid.UUID       `json:"task_id"`
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	PaidAt      *time.Time      `json:"paid_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Payout represents a payout to a worker
type Payout struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	WorkerID      uuid.UUID       `json:"worker_id"`
	Amount        decimal.Decimal `json:"amount"`
	Status        string          `json:"status"`
	PaymentMethod string          `json:"payment_method"`
	Reference     string          `json:"reference"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
	FailureReason string          `json:"failure_reason,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

// TaskProof represents proof of task completion
type TaskProof struct {
	ID        uuid.UUID              `json:"id"`
	TaskID    uuid.UUID              `json:"task_id"`
	Type      string                 `json:"type"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// ============================================================================
// Domain Events
// ============================================================================

// WorkerRegisteredEvent is raised when a new worker joins
type WorkerRegisteredEvent struct {
	WorkerID  uuid.UUID  `json:"worker_id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Type      WorkerType `json:"type"`
	Timestamp time.Time  `json:"timestamp"`
}

// TaskCreatedEvent is raised when a new task is created
type TaskCreatedEvent struct {
	TaskID    uuid.UUID `json:"task_id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Type      TaskType  `json:"type"`
	OrderID   uuid.UUID `json:"order_id"`
	Timestamp time.Time `json:"timestamp"`
}

// TaskAssignedEvent is raised when a task is assigned
type TaskAssignedEvent struct {
	TaskID    uuid.UUID       `json:"task_id"`
	WorkerID  uuid.UUID       `json:"worker_id"`
	Payout    decimal.Decimal `json:"payout"`
	Timestamp time.Time       `json:"timestamp"`
}

// TaskCompletedEvent is raised when a task is completed
type TaskCompletedEvent struct {
	TaskID    uuid.UUID     `json:"task_id"`
	WorkerID  uuid.UUID     `json:"worker_id"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// ============================================================================
// Business Logic Methods
// ============================================================================

// CanTransitionTo checks if a worker status can transition to another
func (s WorkerStatus) CanTransitionTo(target WorkerStatus) bool {
	validTransitions := map[WorkerStatus][]WorkerStatus{
		WorkerStatusOffline:   {WorkerStatusAvailable},
		WorkerStatusAvailable: {WorkerStatusBusy, WorkerStatusOnBreak, WorkerStatusInTransit, WorkerStatusOffline},
		WorkerStatusBusy:      {WorkerStatusAvailable, WorkerStatusInTransit},
		WorkerStatusOnBreak:   {WorkerStatusAvailable, WorkerStatusOffline},
		WorkerStatusInTransit: {WorkerStatusAvailable, WorkerStatusBusy},
	}

	allowed, ok := validTransitions[s]
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

// CanAcceptTask checks if a worker can accept a task of the given type
func (w *GigWorker) CanAcceptTask(taskType TaskType) bool {
	if w.Status != WorkerStatusAvailable {
		return false
	}
	if !w.IsVerified {
		return false
	}

	// Multi-role workers can accept any task
	if w.Type == WorkerTypeMultiRole {
		return true
	}

	// Check task type matches worker type
	switch w.Type {
	case WorkerTypeDelivery:
		return taskType == TaskTypeDelivery || taskType == TaskTypePickup
	case WorkerTypeSales:
		return taskType == TaskTypeSalesVisit
	case WorkerTypeCollection:
		return taskType == TaskTypeCollection
	case WorkerTypeAudit:
		return taskType == TaskTypeAudit
	}
	return false
}

// CalculateSuccessRate calculates the worker's success rate
func (w *GigWorker) CalculateSuccessRate() decimal.Decimal {
	if w.TotalTasks == 0 {
		return decimal.Zero
	}
	return decimal.NewFromInt(int64(w.CompletedTasks)).
		Div(decimal.NewFromInt(int64(w.TotalTasks))).
		Mul(decimal.NewFromInt(100))
}

// AddRating adds a new rating and recalculates the average
func (w *GigWorker) AddRating(newRating decimal.Decimal) {
	totalRatings := decimal.NewFromInt(int64(w.CompletedTasks))
	currentTotal := w.Rating.Mul(totalRatings)
	newTotal := currentTotal.Add(newRating)
	w.Rating = newTotal.Div(totalRatings.Add(decimal.NewFromInt(1)))
}

// CanTransitionTo checks if a task can transition to a new status
func (t *Task) CanTransitionTo(target TaskStatus) bool {
	validTransitions := map[TaskStatus][]TaskStatus{
		TaskStatusPending:    {TaskStatusAssigned, TaskStatusCancelled},
		TaskStatusAssigned:   {TaskStatusAccepted, TaskStatusPending, TaskStatusCancelled},
		TaskStatusAccepted:   {TaskStatusInProgress, TaskStatusCancelled},
		TaskStatusInProgress: {TaskStatusCompleted, TaskStatusFailed},
		TaskStatusCompleted:  {},
		TaskStatusFailed:     {},
		TaskStatusCancelled:  {},
	}

	allowed, ok := validTransitions[t.Status]
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

// CalculateTotalPayout calculates the total payout including surge
func (t *Task) CalculateTotalPayout() decimal.Decimal {
	basePlusSBonus := t.BasePayout.Add(t.BonusPayout)
	return basePlusSBonus.Mul(t.SurgeMultiplier)
}

// IsOverdue checks if the task is past its deadline
func (t *Task) IsOverdue() bool {
	if t.Deadline == nil {
		return false
	}
	if t.Status == TaskStatusCompleted || t.Status == TaskStatusCancelled {
		return false
	}
	return time.Now().After(*t.Deadline)
}

// IsExpired checks if the task offer has expired
func (o *TaskOffer) IsExpired() bool {
	if o.Status != OfferStatusPending {
		return false
	}
	return time.Now().After(o.ExpiresAt)
}

// DistanceTo calculates the distance in kilometers to another location using Haversine formula
func (l Location) DistanceTo(other Location) float64 {
	const earthRadius = 6371.0 // km

	lat1 := l.Latitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	dLat := (other.Latitude - l.Latitude) * math.Pi / 180
	dLon := (other.Longitude - l.Longitude) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// IsValidType checks if the earning type is valid
func (e *Earning) IsValidType() bool {
	validTypes := []string{"task_payout", "bonus", "tip", "penalty", "referral"}
	for _, t := range validTypes {
		if e.Type == t {
			return true
		}
	}
	return false
}
