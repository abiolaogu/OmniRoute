// Package repository provides PostgreSQL data access for the Gig Platform service.
// Handles gig workers, tasks, allocations, earnings, routes, and real-time operations.
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
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("gig-platform/repository")

// Common errors
var (
	ErrNotFound           = errors.New("record not found")
	ErrDuplicateKey       = errors.New("duplicate key violation")
	ErrOptimisticLock     = errors.New("optimistic lock conflict")
	ErrInvalidInput       = errors.New("invalid input")
	ErrWorkerNotAvailable = errors.New("worker not available")
	ErrTaskAlreadyAssigned = errors.New("task already assigned")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// ============================================================================
// Repository Interface
// ============================================================================

// Repository defines the data access interface for the Gig Platform
type Repository interface {
	// Worker operations
	CreateWorker(ctx context.Context, worker *GigWorker) error
	GetWorker(ctx context.Context, id uuid.UUID) (*GigWorker, error)
	GetWorkerByUserID(ctx context.Context, userID uuid.UUID) (*GigWorker, error)
	UpdateWorker(ctx context.Context, worker *GigWorker) error
	ListWorkers(ctx context.Context, filter WorkerFilter) (*WorkerList, error)
	UpdateWorkerStatus(ctx context.Context, id uuid.UUID, status WorkerStatus) error
	UpdateWorkerLocation(ctx context.Context, id uuid.UUID, lat, lng float64) error
	GetNearbyWorkers(ctx context.Context, lat, lng float64, radiusKm float64, filter NearbyFilter) ([]*GigWorker, error)
	
	// Task operations
	CreateTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context, id uuid.UUID) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	ListTasks(ctx context.Context, filter TaskFilter) (*TaskList, error)
	GetTasksByWorker(ctx context.Context, workerID uuid.UUID, status []TaskStatus) ([]*Task, error)
	GetPendingTasks(ctx context.Context, limit int) ([]*Task, error)
	AssignTask(ctx context.Context, taskID, workerID uuid.UUID) error
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status TaskStatus, metadata map[string]interface{}) error
	
	// Task proof operations
	AddTaskProof(ctx context.Context, proof *TaskProof) error
	GetTaskProofs(ctx context.Context, taskID uuid.UUID) ([]*TaskProof, error)
	
	// Allocation operations
	CreateAllocation(ctx context.Context, allocation *Allocation) error
	GetAllocation(ctx context.Context, id uuid.UUID) (*Allocation, error)
	UpdateAllocation(ctx context.Context, allocation *Allocation) error
	GetAllocationsByTask(ctx context.Context, taskID uuid.UUID) ([]*Allocation, error)
	GetAllocationsByWorker(ctx context.Context, workerID uuid.UUID, status []AllocationStatus) ([]*Allocation, error)
	
	// Earnings operations
	CreateEarning(ctx context.Context, earning *Earning) error
	GetEarning(ctx context.Context, id uuid.UUID) (*Earning, error)
	GetEarningsByWorker(ctx context.Context, workerID uuid.UUID, filter EarningFilter) (*EarningList, error)
	GetEarningsSummary(ctx context.Context, workerID uuid.UUID, from, to time.Time) (*EarningsSummary, error)
	
	// Payout operations
	CreatePayout(ctx context.Context, payout *Payout) error
	GetPayout(ctx context.Context, id uuid.UUID) (*Payout, error)
	GetPayoutsByWorker(ctx context.Context, workerID uuid.UUID, filter PayoutFilter) (*PayoutList, error)
	UpdatePayoutStatus(ctx context.Context, id uuid.UUID, status PayoutStatus, reference string) error
	
	// Route operations
	CreateRoute(ctx context.Context, route *Route) error
	GetRoute(ctx context.Context, id uuid.UUID) (*Route, error)
	GetActiveRouteByWorker(ctx context.Context, workerID uuid.UUID) (*Route, error)
	UpdateRoute(ctx context.Context, route *Route) error
	AddRouteStop(ctx context.Context, routeID uuid.UUID, stop *RouteStop) error
	UpdateRouteStop(ctx context.Context, routeID, stopID uuid.UUID, status StopStatus) error
	
	// Performance operations
	GetWorkerPerformance(ctx context.Context, workerID uuid.UUID, from, to time.Time) (*WorkerPerformance, error)
	UpdateWorkerRating(ctx context.Context, workerID uuid.UUID, rating float64) error
	
	// Analytics operations
	GetTaskStats(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*TaskStats, error)
	GetWorkerStats(ctx context.Context, tenantID uuid.UUID) (*WorkerStats, error)
	
	// Transaction support
	WithTx(ctx context.Context, fn func(Repository) error) error
}

// ============================================================================
// Domain Types
// ============================================================================

type WorkerStatus string
const (
	WorkerStatusPending   WorkerStatus = "pending"
	WorkerStatusActive    WorkerStatus = "active"
	WorkerStatusInactive  WorkerStatus = "inactive"
	WorkerStatusSuspended WorkerStatus = "suspended"
	WorkerStatusOffline   WorkerStatus = "offline"
	WorkerStatusOnline    WorkerStatus = "online"
	WorkerStatusBusy      WorkerStatus = "busy"
)

type WorkerType string
const (
	WorkerTypeDriver      WorkerType = "driver"
	WorkerTypeRider       WorkerType = "rider"
	WorkerTypeSalesRep    WorkerType = "sales_rep"
	WorkerTypeWarehouse   WorkerType = "warehouse"
	WorkerTypeMerchandiser WorkerType = "merchandiser"
)

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

type TaskType string
const (
	TaskTypeDelivery     TaskType = "delivery"
	TaskTypePickup       TaskType = "pickup"
	TaskTypeSalesVisit   TaskType = "sales_visit"
	TaskTypeCollection   TaskType = "collection"
	TaskTypeMerchandising TaskType = "merchandising"
	TaskTypeAudit        TaskType = "audit"
)

type AllocationStatus string
const (
	AllocationStatusPending  AllocationStatus = "pending"
	AllocationStatusOffered  AllocationStatus = "offered"
	AllocationStatusAccepted AllocationStatus = "accepted"
	AllocationStatusRejected AllocationStatus = "rejected"
	AllocationStatusExpired  AllocationStatus = "expired"
)

type EarningType string
const (
	EarningTypeTaskCompletion EarningType = "task_completion"
	EarningTypeBonus          EarningType = "bonus"
	EarningTypeTip            EarningType = "tip"
	EarningTypeIncentive      EarningType = "incentive"
	EarningTypeAdjustment     EarningType = "adjustment"
)

type PayoutStatus string
const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

type StopStatus string
const (
	StopStatusPending   StopStatus = "pending"
	StopStatusArrived   StopStatus = "arrived"
	StopStatusCompleted StopStatus = "completed"
	StopStatusSkipped   StopStatus = "skipped"
)

// GigWorker represents a gig economy worker
type GigWorker struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	UserID            uuid.UUID              `json:"user_id"`
	Type              WorkerType             `json:"type"`
	Status            WorkerStatus           `json:"status"`
	FirstName         string                 `json:"first_name"`
	LastName          string                 `json:"last_name"`
	Email             string                 `json:"email"`
	Phone             string                 `json:"phone"`
	ProfilePhotoURL   *string                `json:"profile_photo_url,omitempty"`
	Rating            float64                `json:"rating"`
	TotalTasks        int                    `json:"total_tasks"`
	CompletedTasks    int                    `json:"completed_tasks"`
	CurrentLat        *float64               `json:"current_lat,omitempty"`
	CurrentLng        *float64               `json:"current_lng,omitempty"`
	LastLocationAt    *time.Time             `json:"last_location_at,omitempty"`
	VehicleType       *string                `json:"vehicle_type,omitempty"`
	VehiclePlate      *string                `json:"vehicle_plate,omitempty"`
	LicenseNumber     *string                `json:"license_number,omitempty"`
	BankAccountID     *uuid.UUID             `json:"bank_account_id,omitempty"`
	WalletID          *uuid.UUID             `json:"wallet_id,omitempty"`
	ZoneIDs           []uuid.UUID            `json:"zone_ids"`
	Skills            []string               `json:"skills"`
	MaxConcurrentTasks int                   `json:"max_concurrent_tasks"`
	AvailableFrom     *time.Time             `json:"available_from,omitempty"`
	AvailableTo       *time.Time             `json:"available_to,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Version           int                    `json:"version"`
}

// Task represents a work task
type Task struct {
	ID                 uuid.UUID              `json:"id"`
	TenantID           uuid.UUID              `json:"tenant_id"`
	Type               TaskType               `json:"type"`
	Status             TaskStatus             `json:"status"`
	Priority           int                    `json:"priority"`
	Title              string                 `json:"title"`
	Description        *string                `json:"description,omitempty"`
	OrderID            *uuid.UUID             `json:"order_id,omitempty"`
	CustomerID         *uuid.UUID             `json:"customer_id,omitempty"`
	AssignedWorkerID   *uuid.UUID             `json:"assigned_worker_id,omitempty"`
	RequiredWorkerType WorkerType             `json:"required_worker_type"`
	RequiredSkills     []string               `json:"required_skills"`
	PickupLat          *float64               `json:"pickup_lat,omitempty"`
	PickupLng          *float64               `json:"pickup_lng,omitempty"`
	PickupAddress      *string                `json:"pickup_address,omitempty"`
	DeliveryLat        *float64               `json:"delivery_lat,omitempty"`
	DeliveryLng        *float64               `json:"delivery_lng,omitempty"`
	DeliveryAddress    *string                `json:"delivery_address,omitempty"`
	ScheduledAt        *time.Time             `json:"scheduled_at,omitempty"`
	StartedAt          *time.Time             `json:"started_at,omitempty"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
	DeadlineAt         *time.Time             `json:"deadline_at,omitempty"`
	EstimatedDuration  *int                   `json:"estimated_duration_minutes,omitempty"`
	ActualDuration     *int                   `json:"actual_duration_minutes,omitempty"`
	EstimatedDistance  *float64               `json:"estimated_distance_km,omitempty"`
	ActualDistance     *float64               `json:"actual_distance_km,omitempty"`
	BasePay            decimal.Decimal        `json:"base_pay"`
	BonusPay           decimal.Decimal        `json:"bonus_pay"`
	TipAmount          decimal.Decimal        `json:"tip_amount"`
	Currency           string                 `json:"currency"`
	Notes              *string                `json:"notes,omitempty"`
	FailureReason      *string                `json:"failure_reason,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	Version            int                    `json:"version"`
}

// TaskProof represents proof of task completion (photos, signatures, etc.)
type TaskProof struct {
	ID        uuid.UUID              `json:"id"`
	TaskID    uuid.UUID              `json:"task_id"`
	Type      string                 `json:"type"` // photo, signature, barcode, document
	URL       string                 `json:"url"`
	Caption   *string                `json:"caption,omitempty"`
	Lat       *float64               `json:"lat,omitempty"`
	Lng       *float64               `json:"lng,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// Allocation represents a task allocation attempt to a worker
type Allocation struct {
	ID          uuid.UUID        `json:"id"`
	TenantID    uuid.UUID        `json:"tenant_id"`
	TaskID      uuid.UUID        `json:"task_id"`
	WorkerID    uuid.UUID        `json:"worker_id"`
	Status      AllocationStatus `json:"status"`
	Score       float64          `json:"score"`
	Distance    *float64         `json:"distance_km,omitempty"`
	ETA         *int             `json:"eta_minutes,omitempty"`
	OfferedAt   time.Time        `json:"offered_at"`
	ExpiresAt   time.Time        `json:"expires_at"`
	RespondedAt *time.Time       `json:"responded_at,omitempty"`
	Reason      *string          `json:"reason,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Earning represents worker earnings
type Earning struct {
	ID          uuid.UUID              `json:"id"`
	TenantID    uuid.UUID              `json:"tenant_id"`
	WorkerID    uuid.UUID              `json:"worker_id"`
	TaskID      *uuid.UUID             `json:"task_id,omitempty"`
	Type        EarningType            `json:"type"`
	Amount      decimal.Decimal        `json:"amount"`
	Currency    string                 `json:"currency"`
	Description string                 `json:"description"`
	PayoutID    *uuid.UUID             `json:"payout_id,omitempty"`
	IsPaidOut   bool                   `json:"is_paid_out"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	EarnedAt    time.Time              `json:"earned_at"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Payout represents a payout to a worker
type Payout struct {
	ID             uuid.UUID              `json:"id"`
	TenantID       uuid.UUID              `json:"tenant_id"`
	WorkerID       uuid.UUID              `json:"worker_id"`
	Amount         decimal.Decimal        `json:"amount"`
	Currency       string                 `json:"currency"`
	Status         PayoutStatus           `json:"status"`
	Method         string                 `json:"method"` // bank_transfer, mobile_money, wallet
	BankAccountID  *uuid.UUID             `json:"bank_account_id,omitempty"`
	WalletID       *uuid.UUID             `json:"wallet_id,omitempty"`
	Reference      string                 `json:"reference"`
	ProviderRef    *string                `json:"provider_reference,omitempty"`
	FailureReason  *string                `json:"failure_reason,omitempty"`
	ProcessedAt    *time.Time             `json:"processed_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// Route represents an optimized delivery route
type Route struct {
	ID               uuid.UUID     `json:"id"`
	TenantID         uuid.UUID     `json:"tenant_id"`
	WorkerID         uuid.UUID     `json:"worker_id"`
	Status           string        `json:"status"` // pending, active, completed, cancelled
	TotalDistance    float64       `json:"total_distance_km"`
	TotalDuration    int           `json:"total_duration_minutes"`
	EstimatedEndAt   *time.Time    `json:"estimated_end_at,omitempty"`
	StartedAt        *time.Time    `json:"started_at,omitempty"`
	CompletedAt      *time.Time    `json:"completed_at,omitempty"`
	Stops            []*RouteStop  `json:"stops"`
	OptimizationScore float64      `json:"optimization_score"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// RouteStop represents a stop in a route
type RouteStop struct {
	ID              uuid.UUID              `json:"id"`
	RouteID         uuid.UUID              `json:"route_id"`
	TaskID          uuid.UUID              `json:"task_id"`
	Sequence        int                    `json:"sequence"`
	Status          StopStatus             `json:"status"`
	Lat             float64                `json:"lat"`
	Lng             float64                `json:"lng"`
	Address         string                 `json:"address"`
	DistanceFromPrev float64               `json:"distance_from_prev_km"`
	DurationFromPrev int                   `json:"duration_from_prev_minutes"`
	EstimatedArrival *time.Time            `json:"estimated_arrival,omitempty"`
	ActualArrival    *time.Time            `json:"actual_arrival,omitempty"`
	CompletedAt      *time.Time            `json:"completed_at,omitempty"`
	Notes            *string               `json:"notes,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ============================================================================
// Filter and Result Types
// ============================================================================

type WorkerFilter struct {
	TenantID   uuid.UUID
	Status     []WorkerStatus
	Type       []WorkerType
	ZoneID     *uuid.UUID
	Skills     []string
	Search     *string
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}

type WorkerList struct {
	Workers []*GigWorker `json:"workers"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
}

type NearbyFilter struct {
	Status []WorkerStatus
	Type   []WorkerType
	Skills []string
	Limit  int
}

type TaskFilter struct {
	TenantID         uuid.UUID
	Status           []TaskStatus
	Type             []TaskType
	WorkerID         *uuid.UUID
	CustomerID       *uuid.UUID
	OrderID          *uuid.UUID
	ScheduledFrom    *time.Time
	ScheduledTo      *time.Time
	CreatedFrom      *time.Time
	CreatedTo        *time.Time
	Limit            int
	Offset           int
	SortBy           string
	SortOrder        string
}

type TaskList struct {
	Tasks  []*Task `json:"tasks"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type EarningFilter struct {
	Type      []EarningType
	IsPaidOut *bool
	From      *time.Time
	To        *time.Time
	Limit     int
	Offset    int
}

type EarningList struct {
	Earnings []*Earning `json:"earnings"`
	Total    int        `json:"total"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type EarningsSummary struct {
	WorkerID         uuid.UUID       `json:"worker_id"`
	TotalEarnings    decimal.Decimal `json:"total_earnings"`
	PendingPayout    decimal.Decimal `json:"pending_payout"`
	PaidOut          decimal.Decimal `json:"paid_out"`
	TaskCompletionAmt decimal.Decimal `json:"task_completion_amount"`
	BonusAmount      decimal.Decimal `json:"bonus_amount"`
	TipAmount        decimal.Decimal `json:"tip_amount"`
	IncentiveAmount  decimal.Decimal `json:"incentive_amount"`
	TaskCount        int             `json:"task_count"`
	Currency         string          `json:"currency"`
	FromDate         time.Time       `json:"from_date"`
	ToDate           time.Time       `json:"to_date"`
}

type PayoutFilter struct {
	Status []PayoutStatus
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

type PayoutList struct {
	Payouts []*Payout `json:"payouts"`
	Total   int       `json:"total"`
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
}

type WorkerPerformance struct {
	WorkerID           uuid.UUID       `json:"worker_id"`
	TotalTasks         int             `json:"total_tasks"`
	CompletedTasks     int             `json:"completed_tasks"`
	FailedTasks        int             `json:"failed_tasks"`
	CancelledTasks     int             `json:"cancelled_tasks"`
	CompletionRate     float64         `json:"completion_rate"`
	OnTimeRate         float64         `json:"on_time_rate"`
	AverageRating      float64         `json:"average_rating"`
	TotalEarnings      decimal.Decimal `json:"total_earnings"`
	AverageTaskTime    int             `json:"average_task_time_minutes"`
	TotalDistanceTraveled float64      `json:"total_distance_km"`
	FromDate           time.Time       `json:"from_date"`
	ToDate             time.Time       `json:"to_date"`
}

type TaskStats struct {
	TenantID        uuid.UUID `json:"tenant_id"`
	TotalTasks      int       `json:"total_tasks"`
	PendingTasks    int       `json:"pending_tasks"`
	AssignedTasks   int       `json:"assigned_tasks"`
	InProgressTasks int       `json:"in_progress_tasks"`
	CompletedTasks  int       `json:"completed_tasks"`
	FailedTasks     int       `json:"failed_tasks"`
	CancelledTasks  int       `json:"cancelled_tasks"`
	AvgCompletionTime int     `json:"avg_completion_time_minutes"`
	FromDate        time.Time `json:"from_date"`
	ToDate          time.Time `json:"to_date"`
}

type WorkerStats struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	TotalWorkers   int       `json:"total_workers"`
	ActiveWorkers  int       `json:"active_workers"`
	OnlineWorkers  int       `json:"online_workers"`
	BusyWorkers    int       `json:"busy_workers"`
	ByType         map[WorkerType]int `json:"by_type"`
}

// ============================================================================
// PostgreSQL Implementation
// ============================================================================

type postgresRepository struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) conn(ctx context.Context) interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.pool
}

// pgconn import
type pgconn struct{}
type CommandTag = interface{ RowsAffected() int64 }

// ============================================================================
// Worker Operations
// ============================================================================

func (r *postgresRepository) CreateWorker(ctx context.Context, worker *GigWorker) error {
	ctx, span := tracer.Start(ctx, "repository.CreateWorker")
	defer span.End()
	
	if worker.ID == uuid.Nil {
		worker.ID = uuid.New()
	}
	worker.CreatedAt = time.Now().UTC()
	worker.UpdatedAt = worker.CreatedAt
	worker.Version = 1
	
	zoneIDs, _ := json.Marshal(worker.ZoneIDs)
	skills, _ := json.Marshal(worker.Skills)
	metadata, _ := json.Marshal(worker.Metadata)
	
	query := `
		INSERT INTO gig_workers (
			id, tenant_id, user_id, type, status, first_name, last_name, email, phone,
			profile_photo_url, rating, total_tasks, completed_tasks, current_location,
			last_location_at, vehicle_type, vehicle_plate, license_number,
			bank_account_id, wallet_id, zone_ids, skills, max_concurrent_tasks,
			available_from, available_to, metadata, created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			ST_SetSRID(ST_MakePoint($14, $15), 4326), $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)`
	
	var lat, lng float64
	if worker.CurrentLat != nil && worker.CurrentLng != nil {
		lat, lng = *worker.CurrentLat, *worker.CurrentLng
	}
	
	_, err := r.pool.Exec(ctx, query,
		worker.ID, worker.TenantID, worker.UserID, worker.Type, worker.Status,
		worker.FirstName, worker.LastName, worker.Email, worker.Phone,
		worker.ProfilePhotoURL, worker.Rating, worker.TotalTasks, worker.CompletedTasks,
		lng, lat, worker.LastLocationAt, worker.VehicleType, worker.VehiclePlate,
		worker.LicenseNumber, worker.BankAccountID, worker.WalletID,
		zoneIDs, skills, worker.MaxConcurrentTasks,
		worker.AvailableFrom, worker.AvailableTo, metadata,
		worker.CreatedAt, worker.UpdatedAt, worker.Version,
	)
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateKey
		}
		span.RecordError(err)
		return fmt.Errorf("create worker: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetWorker(ctx context.Context, id uuid.UUID) (*GigWorker, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWorker")
	defer span.End()
	span.SetAttributes(attribute.String("worker.id", id.String()))
	
	query := `
		SELECT id, tenant_id, user_id, type, status, first_name, last_name, email, phone,
			profile_photo_url, rating, total_tasks, completed_tasks,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng,
			last_location_at, vehicle_type, vehicle_plate, license_number,
			bank_account_id, wallet_id, zone_ids, skills, max_concurrent_tasks,
			available_from, available_to, metadata, created_at, updated_at, version
		FROM gig_workers
		WHERE id = $1`
	
	worker := &GigWorker{}
	var zoneIDs, skills, metadata []byte
	var lat, lng sql.NullFloat64
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&worker.ID, &worker.TenantID, &worker.UserID, &worker.Type, &worker.Status,
		&worker.FirstName, &worker.LastName, &worker.Email, &worker.Phone,
		&worker.ProfilePhotoURL, &worker.Rating, &worker.TotalTasks, &worker.CompletedTasks,
		&lat, &lng, &worker.LastLocationAt, &worker.VehicleType, &worker.VehiclePlate,
		&worker.LicenseNumber, &worker.BankAccountID, &worker.WalletID,
		&zoneIDs, &skills, &worker.MaxConcurrentTasks,
		&worker.AvailableFrom, &worker.AvailableTo, &metadata,
		&worker.CreatedAt, &worker.UpdatedAt, &worker.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		span.RecordError(err)
		return nil, fmt.Errorf("get worker: %w", err)
	}
	
	if lat.Valid && lng.Valid {
		worker.CurrentLat = &lat.Float64
		worker.CurrentLng = &lng.Float64
	}
	
	json.Unmarshal(zoneIDs, &worker.ZoneIDs)
	json.Unmarshal(skills, &worker.Skills)
	json.Unmarshal(metadata, &worker.Metadata)
	
	return worker, nil
}

func (r *postgresRepository) GetWorkerByUserID(ctx context.Context, userID uuid.UUID) (*GigWorker, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWorkerByUserID")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, user_id, type, status, first_name, last_name, email, phone,
			profile_photo_url, rating, total_tasks, completed_tasks,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng,
			last_location_at, vehicle_type, vehicle_plate, license_number,
			bank_account_id, wallet_id, zone_ids, skills, max_concurrent_tasks,
			available_from, available_to, metadata, created_at, updated_at, version
		FROM gig_workers
		WHERE user_id = $1`
	
	worker := &GigWorker{}
	var zoneIDs, skills, metadata []byte
	var lat, lng sql.NullFloat64
	
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&worker.ID, &worker.TenantID, &worker.UserID, &worker.Type, &worker.Status,
		&worker.FirstName, &worker.LastName, &worker.Email, &worker.Phone,
		&worker.ProfilePhotoURL, &worker.Rating, &worker.TotalTasks, &worker.CompletedTasks,
		&lat, &lng, &worker.LastLocationAt, &worker.VehicleType, &worker.VehiclePlate,
		&worker.LicenseNumber, &worker.BankAccountID, &worker.WalletID,
		&zoneIDs, &skills, &worker.MaxConcurrentTasks,
		&worker.AvailableFrom, &worker.AvailableTo, &metadata,
		&worker.CreatedAt, &worker.UpdatedAt, &worker.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get worker by user id: %w", err)
	}
	
	if lat.Valid && lng.Valid {
		worker.CurrentLat = &lat.Float64
		worker.CurrentLng = &lng.Float64
	}
	
	json.Unmarshal(zoneIDs, &worker.ZoneIDs)
	json.Unmarshal(skills, &worker.Skills)
	json.Unmarshal(metadata, &worker.Metadata)
	
	return worker, nil
}

func (r *postgresRepository) UpdateWorker(ctx context.Context, worker *GigWorker) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateWorker")
	defer span.End()
	
	zoneIDs, _ := json.Marshal(worker.ZoneIDs)
	skills, _ := json.Marshal(worker.Skills)
	metadata, _ := json.Marshal(worker.Metadata)
	
	query := `
		UPDATE gig_workers SET
			type = $3, status = $4, first_name = $5, last_name = $6, email = $7, phone = $8,
			profile_photo_url = $9, rating = $10, vehicle_type = $11, vehicle_plate = $12,
			license_number = $13, bank_account_id = $14, wallet_id = $15,
			zone_ids = $16, skills = $17, max_concurrent_tasks = $18,
			available_from = $19, available_to = $20, metadata = $21,
			updated_at = $22, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`
	
	worker.UpdatedAt = time.Now().UTC()
	
	err := r.pool.QueryRow(ctx, query,
		worker.ID, worker.Version, worker.Type, worker.Status,
		worker.FirstName, worker.LastName, worker.Email, worker.Phone,
		worker.ProfilePhotoURL, worker.Rating, worker.VehicleType, worker.VehiclePlate,
		worker.LicenseNumber, worker.BankAccountID, worker.WalletID,
		zoneIDs, skills, worker.MaxConcurrentTasks,
		worker.AvailableFrom, worker.AvailableTo, metadata, worker.UpdatedAt,
	).Scan(&worker.Version)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptimisticLock
		}
		span.RecordError(err)
		return fmt.Errorf("update worker: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) ListWorkers(ctx context.Context, filter WorkerFilter) (*WorkerList, error) {
	ctx, span := tracer.Start(ctx, "repository.ListWorkers")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, s)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if len(filter.Type) > 0 {
		placeholders := make([]string, len(filter.Type))
		for i, t := range filter.Type {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, t)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if filter.Search != nil && *filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)",
			argNum, argNum, argNum, argNum,
		))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM gig_workers WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count workers: %w", err)
	}
	
	// Build sort
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	
	// Fetch records
	query := fmt.Sprintf(`
		SELECT id, tenant_id, user_id, type, status, first_name, last_name, email, phone,
			profile_photo_url, rating, total_tasks, completed_tasks,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng,
			last_location_at, vehicle_type, vehicle_plate, license_number,
			bank_account_id, wallet_id, zone_ids, skills, max_concurrent_tasks,
			available_from, available_to, metadata, created_at, updated_at, version
		FROM gig_workers
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		whereClause, sortBy, sortOrder, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list workers: %w", err)
	}
	defer rows.Close()
	
	workers := make([]*GigWorker, 0)
	for rows.Next() {
		worker := &GigWorker{}
		var zoneIDs, skills, metadata []byte
		var lat, lng sql.NullFloat64
		
		err := rows.Scan(
			&worker.ID, &worker.TenantID, &worker.UserID, &worker.Type, &worker.Status,
			&worker.FirstName, &worker.LastName, &worker.Email, &worker.Phone,
			&worker.ProfilePhotoURL, &worker.Rating, &worker.TotalTasks, &worker.CompletedTasks,
			&lat, &lng, &worker.LastLocationAt, &worker.VehicleType, &worker.VehiclePlate,
			&worker.LicenseNumber, &worker.BankAccountID, &worker.WalletID,
			&zoneIDs, &skills, &worker.MaxConcurrentTasks,
			&worker.AvailableFrom, &worker.AvailableTo, &metadata,
			&worker.CreatedAt, &worker.UpdatedAt, &worker.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan worker: %w", err)
		}
		
		if lat.Valid && lng.Valid {
			worker.CurrentLat = &lat.Float64
			worker.CurrentLng = &lng.Float64
		}
		json.Unmarshal(zoneIDs, &worker.ZoneIDs)
		json.Unmarshal(skills, &worker.Skills)
		json.Unmarshal(metadata, &worker.Metadata)
		
		workers = append(workers, worker)
	}
	
	return &WorkerList{
		Workers: workers,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}

func (r *postgresRepository) UpdateWorkerStatus(ctx context.Context, id uuid.UUID, status WorkerStatus) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateWorkerStatus")
	defer span.End()
	
	query := `UPDATE gig_workers SET status = $2, updated_at = $3 WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id, status, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("update worker status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) UpdateWorkerLocation(ctx context.Context, id uuid.UUID, lat, lng float64) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateWorkerLocation")
	defer span.End()
	
	query := `
		UPDATE gig_workers 
		SET current_location = ST_SetSRID(ST_MakePoint($2, $3), 4326),
			last_location_at = $4,
			updated_at = $4
		WHERE id = $1`
	
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, query, id, lng, lat, now)
	if err != nil {
		return fmt.Errorf("update worker location: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) GetNearbyWorkers(ctx context.Context, lat, lng float64, radiusKm float64, filter NearbyFilter) ([]*GigWorker, error) {
	ctx, span := tracer.Start(ctx, "repository.GetNearbyWorkers")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	
	// PostGIS distance calculation in meters
	conditions = append(conditions, `ST_DWithin(
		current_location::geography,
		ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
		$3
	)`)
	args = append(args, lng, lat, radiusKm*1000)
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
			args = append(args, s)
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if len(filter.Type) > 0 {
		placeholders := make([]string, len(filter.Type))
		for i, t := range filter.Type {
			placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
			args = append(args, t)
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	limit := 50
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, user_id, type, status, first_name, last_name, email, phone,
			profile_photo_url, rating, total_tasks, completed_tasks,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng,
			last_location_at, vehicle_type, vehicle_plate, license_number,
			bank_account_id, wallet_id, zone_ids, skills, max_concurrent_tasks,
			available_from, available_to, metadata, created_at, updated_at, version,
			ST_Distance(
				current_location::geography,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
			) / 1000 as distance_km
		FROM gig_workers
		WHERE %s
		ORDER BY distance_km ASC
		LIMIT $%d`, whereClause, len(args)+1)
	
	args = append(args, limit)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get nearby workers: %w", err)
	}
	defer rows.Close()
	
	workers := make([]*GigWorker, 0)
	for rows.Next() {
		worker := &GigWorker{}
		var zoneIDs, skills, metadata []byte
		var workerLat, workerLng sql.NullFloat64
		var distanceKm float64
		
		err := rows.Scan(
			&worker.ID, &worker.TenantID, &worker.UserID, &worker.Type, &worker.Status,
			&worker.FirstName, &worker.LastName, &worker.Email, &worker.Phone,
			&worker.ProfilePhotoURL, &worker.Rating, &worker.TotalTasks, &worker.CompletedTasks,
			&workerLat, &workerLng, &worker.LastLocationAt, &worker.VehicleType, &worker.VehiclePlate,
			&worker.LicenseNumber, &worker.BankAccountID, &worker.WalletID,
			&zoneIDs, &skills, &worker.MaxConcurrentTasks,
			&worker.AvailableFrom, &worker.AvailableTo, &metadata,
			&worker.CreatedAt, &worker.UpdatedAt, &worker.Version,
			&distanceKm,
		)
		if err != nil {
			return nil, fmt.Errorf("scan worker: %w", err)
		}
		
		if workerLat.Valid && workerLng.Valid {
			worker.CurrentLat = &workerLat.Float64
			worker.CurrentLng = &workerLng.Float64
		}
		json.Unmarshal(zoneIDs, &worker.ZoneIDs)
		json.Unmarshal(skills, &worker.Skills)
		json.Unmarshal(metadata, &worker.Metadata)
		
		workers = append(workers, worker)
	}
	
	return workers, nil
}

// ============================================================================
// Task Operations
// ============================================================================

func (r *postgresRepository) CreateTask(ctx context.Context, task *Task) error {
	ctx, span := tracer.Start(ctx, "repository.CreateTask")
	defer span.End()
	
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	task.CreatedAt = time.Now().UTC()
	task.UpdatedAt = task.CreatedAt
	task.Version = 1
	
	if task.Status == "" {
		task.Status = TaskStatusPending
	}
	
	requiredSkills, _ := json.Marshal(task.RequiredSkills)
	metadata, _ := json.Marshal(task.Metadata)
	
	query := `
		INSERT INTO tasks (
			id, tenant_id, type, status, priority, title, description,
			order_id, customer_id, assigned_worker_id, required_worker_type, required_skills,
			pickup_location, pickup_address, delivery_location, delivery_address,
			scheduled_at, started_at, completed_at, deadline_at,
			estimated_duration_minutes, actual_duration_minutes,
			estimated_distance_km, actual_distance_km,
			base_pay, bonus_pay, tip_amount, currency,
			notes, failure_reason, metadata, created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			ST_SetSRID(ST_MakePoint($13, $14), 4326),
			$15,
			ST_SetSRID(ST_MakePoint($16, $17), 4326),
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36
		)`
	
	var pickupLng, pickupLat, deliveryLng, deliveryLat float64
	if task.PickupLat != nil && task.PickupLng != nil {
		pickupLat, pickupLng = *task.PickupLat, *task.PickupLng
	}
	if task.DeliveryLat != nil && task.DeliveryLng != nil {
		deliveryLat, deliveryLng = *task.DeliveryLat, *task.DeliveryLng
	}
	
	_, err := r.pool.Exec(ctx, query,
		task.ID, task.TenantID, task.Type, task.Status, task.Priority, task.Title, task.Description,
		task.OrderID, task.CustomerID, task.AssignedWorkerID, task.RequiredWorkerType, requiredSkills,
		pickupLng, pickupLat, task.PickupAddress,
		deliveryLng, deliveryLat, task.DeliveryAddress,
		task.ScheduledAt, task.StartedAt, task.CompletedAt, task.DeadlineAt,
		task.EstimatedDuration, task.ActualDuration, task.EstimatedDistance, task.ActualDistance,
		task.BasePay, task.BonusPay, task.TipAmount, task.Currency,
		task.Notes, task.FailureReason, metadata, task.CreatedAt, task.UpdatedAt, task.Version,
	)
	
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("create task: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	ctx, span := tracer.Start(ctx, "repository.GetTask")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, type, status, priority, title, description,
			order_id, customer_id, assigned_worker_id, required_worker_type, required_skills,
			ST_Y(pickup_location::geometry), ST_X(pickup_location::geometry), pickup_address,
			ST_Y(delivery_location::geometry), ST_X(delivery_location::geometry), delivery_address,
			scheduled_at, started_at, completed_at, deadline_at,
			estimated_duration_minutes, actual_duration_minutes,
			estimated_distance_km, actual_distance_km,
			base_pay, bonus_pay, tip_amount, currency,
			notes, failure_reason, metadata, created_at, updated_at, version
		FROM tasks
		WHERE id = $1`
	
	task := &Task{}
	var requiredSkills, metadata []byte
	var pickupLat, pickupLng, deliveryLat, deliveryLng sql.NullFloat64
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.TenantID, &task.Type, &task.Status, &task.Priority, &task.Title, &task.Description,
		&task.OrderID, &task.CustomerID, &task.AssignedWorkerID, &task.RequiredWorkerType, &requiredSkills,
		&pickupLat, &pickupLng, &task.PickupAddress,
		&deliveryLat, &deliveryLng, &task.DeliveryAddress,
		&task.ScheduledAt, &task.StartedAt, &task.CompletedAt, &task.DeadlineAt,
		&task.EstimatedDuration, &task.ActualDuration, &task.EstimatedDistance, &task.ActualDistance,
		&task.BasePay, &task.BonusPay, &task.TipAmount, &task.Currency,
		&task.Notes, &task.FailureReason, &metadata, &task.CreatedAt, &task.UpdatedAt, &task.Version,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		span.RecordError(err)
		return nil, fmt.Errorf("get task: %w", err)
	}
	
	if pickupLat.Valid && pickupLng.Valid {
		task.PickupLat = &pickupLat.Float64
		task.PickupLng = &pickupLng.Float64
	}
	if deliveryLat.Valid && deliveryLng.Valid {
		task.DeliveryLat = &deliveryLat.Float64
		task.DeliveryLng = &deliveryLng.Float64
	}
	
	json.Unmarshal(requiredSkills, &task.RequiredSkills)
	json.Unmarshal(metadata, &task.Metadata)
	
	return task, nil
}

func (r *postgresRepository) UpdateTask(ctx context.Context, task *Task) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateTask")
	defer span.End()
	
	requiredSkills, _ := json.Marshal(task.RequiredSkills)
	metadata, _ := json.Marshal(task.Metadata)
	
	var pickupLng, pickupLat, deliveryLng, deliveryLat float64
	if task.PickupLat != nil && task.PickupLng != nil {
		pickupLat, pickupLng = *task.PickupLat, *task.PickupLng
	}
	if task.DeliveryLat != nil && task.DeliveryLng != nil {
		deliveryLat, deliveryLng = *task.DeliveryLat, *task.DeliveryLng
	}
	
	query := `
		UPDATE tasks SET
			type = $3, status = $4, priority = $5, title = $6, description = $7,
			order_id = $8, customer_id = $9, assigned_worker_id = $10,
			required_worker_type = $11, required_skills = $12,
			pickup_location = ST_SetSRID(ST_MakePoint($13, $14), 4326),
			pickup_address = $15,
			delivery_location = ST_SetSRID(ST_MakePoint($16, $17), 4326),
			delivery_address = $18,
			scheduled_at = $19, started_at = $20, completed_at = $21, deadline_at = $22,
			estimated_duration_minutes = $23, actual_duration_minutes = $24,
			estimated_distance_km = $25, actual_distance_km = $26,
			base_pay = $27, bonus_pay = $28, tip_amount = $29, currency = $30,
			notes = $31, failure_reason = $32, metadata = $33,
			updated_at = $34, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`
	
	task.UpdatedAt = time.Now().UTC()
	
	err := r.pool.QueryRow(ctx, query,
		task.ID, task.Version,
		task.Type, task.Status, task.Priority, task.Title, task.Description,
		task.OrderID, task.CustomerID, task.AssignedWorkerID,
		task.RequiredWorkerType, requiredSkills,
		pickupLng, pickupLat, task.PickupAddress,
		deliveryLng, deliveryLat, task.DeliveryAddress,
		task.ScheduledAt, task.StartedAt, task.CompletedAt, task.DeadlineAt,
		task.EstimatedDuration, task.ActualDuration, task.EstimatedDistance, task.ActualDistance,
		task.BasePay, task.BonusPay, task.TipAmount, task.Currency,
		task.Notes, task.FailureReason, metadata, task.UpdatedAt,
	).Scan(&task.Version)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptimisticLock
		}
		span.RecordError(err)
		return fmt.Errorf("update task: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) ListTasks(ctx context.Context, filter TaskFilter) (*TaskList, error) {
	ctx, span := tracer.Start(ctx, "repository.ListTasks")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, s)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if len(filter.Type) > 0 {
		placeholders := make([]string, len(filter.Type))
		for i, t := range filter.Type {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, t)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if filter.WorkerID != nil {
		conditions = append(conditions, fmt.Sprintf("assigned_worker_id = $%d", argNum))
		args = append(args, *filter.WorkerID)
		argNum++
	}
	
	if filter.OrderID != nil {
		conditions = append(conditions, fmt.Sprintf("order_id = $%d", argNum))
		args = append(args, *filter.OrderID)
		argNum++
	}
	
	if filter.ScheduledFrom != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at >= $%d", argNum))
		args = append(args, *filter.ScheduledFrom)
		argNum++
	}
	
	if filter.ScheduledTo != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at <= $%d", argNum))
		args = append(args, *filter.ScheduledTo)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count tasks: %w", err)
	}
	
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, type, status, priority, title, description,
			order_id, customer_id, assigned_worker_id, required_worker_type, required_skills,
			ST_Y(pickup_location::geometry), ST_X(pickup_location::geometry), pickup_address,
			ST_Y(delivery_location::geometry), ST_X(delivery_location::geometry), delivery_address,
			scheduled_at, started_at, completed_at, deadline_at,
			estimated_duration_minutes, actual_duration_minutes,
			estimated_distance_km, actual_distance_km,
			base_pay, bonus_pay, tip_amount, currency,
			notes, failure_reason, metadata, created_at, updated_at, version
		FROM tasks
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`, whereClause, sortBy, sortOrder, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()
	
	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		var requiredSkills, metadata []byte
		var pickupLat, pickupLng, deliveryLat, deliveryLng sql.NullFloat64
		
		err := rows.Scan(
			&task.ID, &task.TenantID, &task.Type, &task.Status, &task.Priority, &task.Title, &task.Description,
			&task.OrderID, &task.CustomerID, &task.AssignedWorkerID, &task.RequiredWorkerType, &requiredSkills,
			&pickupLat, &pickupLng, &task.PickupAddress,
			&deliveryLat, &deliveryLng, &task.DeliveryAddress,
			&task.ScheduledAt, &task.StartedAt, &task.CompletedAt, &task.DeadlineAt,
			&task.EstimatedDuration, &task.ActualDuration, &task.EstimatedDistance, &task.ActualDistance,
			&task.BasePay, &task.BonusPay, &task.TipAmount, &task.Currency,
			&task.Notes, &task.FailureReason, &metadata, &task.CreatedAt, &task.UpdatedAt, &task.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		
		if pickupLat.Valid && pickupLng.Valid {
			task.PickupLat = &pickupLat.Float64
			task.PickupLng = &pickupLng.Float64
		}
		if deliveryLat.Valid && deliveryLng.Valid {
			task.DeliveryLat = &deliveryLat.Float64
			task.DeliveryLng = &deliveryLng.Float64
		}
		
		json.Unmarshal(requiredSkills, &task.RequiredSkills)
		json.Unmarshal(metadata, &task.Metadata)
		
		tasks = append(tasks, task)
	}
	
	return &TaskList{
		Tasks:  tasks,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

func (r *postgresRepository) GetTasksByWorker(ctx context.Context, workerID uuid.UUID, status []TaskStatus) ([]*Task, error) {
	ctx, span := tracer.Start(ctx, "repository.GetTasksByWorker")
	defer span.End()
	
	var args []interface{}
	args = append(args, workerID)
	
	statusClause := ""
	if len(status) > 0 {
		placeholders := make([]string, len(status))
		for i, s := range status {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, s)
		}
		statusClause = fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, type, status, priority, title, description,
			order_id, customer_id, assigned_worker_id, required_worker_type, required_skills,
			ST_Y(pickup_location::geometry), ST_X(pickup_location::geometry), pickup_address,
			ST_Y(delivery_location::geometry), ST_X(delivery_location::geometry), delivery_address,
			scheduled_at, started_at, completed_at, deadline_at,
			estimated_duration_minutes, actual_duration_minutes,
			estimated_distance_km, actual_distance_km,
			base_pay, bonus_pay, tip_amount, currency,
			notes, failure_reason, metadata, created_at, updated_at, version
		FROM tasks
		WHERE assigned_worker_id = $1%s
		ORDER BY COALESCE(scheduled_at, created_at) ASC`, statusClause)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get tasks by worker: %w", err)
	}
	defer rows.Close()
	
	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		var requiredSkills, metadata []byte
		var pickupLat, pickupLng, deliveryLat, deliveryLng sql.NullFloat64
		
		err := rows.Scan(
			&task.ID, &task.TenantID, &task.Type, &task.Status, &task.Priority, &task.Title, &task.Description,
			&task.OrderID, &task.CustomerID, &task.AssignedWorkerID, &task.RequiredWorkerType, &requiredSkills,
			&pickupLat, &pickupLng, &task.PickupAddress,
			&deliveryLat, &deliveryLng, &task.DeliveryAddress,
			&task.ScheduledAt, &task.StartedAt, &task.CompletedAt, &task.DeadlineAt,
			&task.EstimatedDuration, &task.ActualDuration, &task.EstimatedDistance, &task.ActualDistance,
			&task.BasePay, &task.BonusPay, &task.TipAmount, &task.Currency,
			&task.Notes, &task.FailureReason, &metadata, &task.CreatedAt, &task.UpdatedAt, &task.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		
		if pickupLat.Valid && pickupLng.Valid {
			task.PickupLat = &pickupLat.Float64
			task.PickupLng = &pickupLng.Float64
		}
		if deliveryLat.Valid && deliveryLng.Valid {
			task.DeliveryLat = &deliveryLat.Float64
			task.DeliveryLng = &deliveryLng.Float64
		}
		
		json.Unmarshal(requiredSkills, &task.RequiredSkills)
		json.Unmarshal(metadata, &task.Metadata)
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *postgresRepository) GetPendingTasks(ctx context.Context, limit int) ([]*Task, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPendingTasks")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, type, status, priority, title, description,
			order_id, customer_id, assigned_worker_id, required_worker_type, required_skills,
			ST_Y(pickup_location::geometry), ST_X(pickup_location::geometry), pickup_address,
			ST_Y(delivery_location::geometry), ST_X(delivery_location::geometry), delivery_address,
			scheduled_at, started_at, completed_at, deadline_at,
			estimated_duration_minutes, actual_duration_minutes,
			estimated_distance_km, actual_distance_km,
			base_pay, bonus_pay, tip_amount, currency,
			notes, failure_reason, metadata, created_at, updated_at, version
		FROM tasks
		WHERE status = 'pending' AND assigned_worker_id IS NULL
		ORDER BY priority DESC, created_at ASC
		LIMIT $1`
	
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get pending tasks: %w", err)
	}
	defer rows.Close()
	
	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		var requiredSkills, metadata []byte
		var pickupLat, pickupLng, deliveryLat, deliveryLng sql.NullFloat64
		
		err := rows.Scan(
			&task.ID, &task.TenantID, &task.Type, &task.Status, &task.Priority, &task.Title, &task.Description,
			&task.OrderID, &task.CustomerID, &task.AssignedWorkerID, &task.RequiredWorkerType, &requiredSkills,
			&pickupLat, &pickupLng, &task.PickupAddress,
			&deliveryLat, &deliveryLng, &task.DeliveryAddress,
			&task.ScheduledAt, &task.StartedAt, &task.CompletedAt, &task.DeadlineAt,
			&task.EstimatedDuration, &task.ActualDuration, &task.EstimatedDistance, &task.ActualDistance,
			&task.BasePay, &task.BonusPay, &task.TipAmount, &task.Currency,
			&task.Notes, &task.FailureReason, &metadata, &task.CreatedAt, &task.UpdatedAt, &task.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		
		if pickupLat.Valid && pickupLng.Valid {
			task.PickupLat = &pickupLat.Float64
			task.PickupLng = &pickupLng.Float64
		}
		if deliveryLat.Valid && deliveryLng.Valid {
			task.DeliveryLat = &deliveryLat.Float64
			task.DeliveryLng = &deliveryLng.Float64
		}
		
		json.Unmarshal(requiredSkills, &task.RequiredSkills)
		json.Unmarshal(metadata, &task.Metadata)
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *postgresRepository) AssignTask(ctx context.Context, taskID, workerID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "repository.AssignTask")
	defer span.End()
	
	query := `
		UPDATE tasks 
		SET assigned_worker_id = $2, 
			status = 'assigned',
			updated_at = $3
		WHERE id = $1 AND status = 'pending' AND assigned_worker_id IS NULL
		RETURNING id`
	
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, taskID, workerID, time.Now().UTC()).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskAlreadyAssigned
		}
		return fmt.Errorf("assign task: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status TaskStatus, metadata map[string]interface{}) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateTaskStatus")
	defer span.End()
	
	now := time.Now().UTC()
	
	var startedAt, completedAt *time.Time
	switch status {
	case TaskStatusInProgress:
		startedAt = &now
	case TaskStatusCompleted, TaskStatusFailed:
		completedAt = &now
	}
	
	query := `
		UPDATE tasks SET
			status = $2,
			started_at = COALESCE($3, started_at),
			completed_at = COALESCE($4, completed_at),
			updated_at = $5
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, status, startedAt, completedAt, now)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

// ============================================================================
// Task Proof Operations
// ============================================================================

func (r *postgresRepository) AddTaskProof(ctx context.Context, proof *TaskProof) error {
	ctx, span := tracer.Start(ctx, "repository.AddTaskProof")
	defer span.End()
	
	if proof.ID == uuid.Nil {
		proof.ID = uuid.New()
	}
	proof.CreatedAt = time.Now().UTC()
	
	metadata, _ := json.Marshal(proof.Metadata)
	
	query := `
		INSERT INTO task_proofs (id, task_id, type, url, caption, location, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8, $9)`
	
	var lng, lat float64
	if proof.Lat != nil && proof.Lng != nil {
		lat, lng = *proof.Lat, *proof.Lng
	}
	
	_, err := r.pool.Exec(ctx, query,
		proof.ID, proof.TaskID, proof.Type, proof.URL, proof.Caption,
		lng, lat, metadata, proof.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("add task proof: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetTaskProofs(ctx context.Context, taskID uuid.UUID) ([]*TaskProof, error) {
	ctx, span := tracer.Start(ctx, "repository.GetTaskProofs")
	defer span.End()
	
	query := `
		SELECT id, task_id, type, url, caption,
			ST_Y(location::geometry), ST_X(location::geometry),
			metadata, created_at
		FROM task_proofs
		WHERE task_id = $1
		ORDER BY created_at ASC`
	
	rows, err := r.pool.Query(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task proofs: %w", err)
	}
	defer rows.Close()
	
	proofs := make([]*TaskProof, 0)
	for rows.Next() {
		proof := &TaskProof{}
		var metadata []byte
		var lat, lng sql.NullFloat64
		
		err := rows.Scan(
			&proof.ID, &proof.TaskID, &proof.Type, &proof.URL, &proof.Caption,
			&lat, &lng, &metadata, &proof.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan proof: %w", err)
		}
		
		if lat.Valid && lng.Valid {
			proof.Lat = &lat.Float64
			proof.Lng = &lng.Float64
		}
		json.Unmarshal(metadata, &proof.Metadata)
		
		proofs = append(proofs, proof)
	}
	
	return proofs, nil
}

// ============================================================================
// Allocation Operations
// ============================================================================

func (r *postgresRepository) CreateAllocation(ctx context.Context, allocation *Allocation) error {
	ctx, span := tracer.Start(ctx, "repository.CreateAllocation")
	defer span.End()
	
	if allocation.ID == uuid.Nil {
		allocation.ID = uuid.New()
	}
	allocation.CreatedAt = time.Now().UTC()
	allocation.UpdatedAt = allocation.CreatedAt
	
	query := `
		INSERT INTO task_allocations (
			id, tenant_id, task_id, worker_id, status, score, distance_km,
			eta_minutes, offered_at, expires_at, responded_at, reason,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	
	_, err := r.pool.Exec(ctx, query,
		allocation.ID, allocation.TenantID, allocation.TaskID, allocation.WorkerID,
		allocation.Status, allocation.Score, allocation.Distance, allocation.ETA,
		allocation.OfferedAt, allocation.ExpiresAt, allocation.RespondedAt, allocation.Reason,
		allocation.CreatedAt, allocation.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create allocation: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetAllocation(ctx context.Context, id uuid.UUID) (*Allocation, error) {
	ctx, span := tracer.Start(ctx, "repository.GetAllocation")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, task_id, worker_id, status, score, distance_km,
			eta_minutes, offered_at, expires_at, responded_at, reason,
			created_at, updated_at
		FROM task_allocations
		WHERE id = $1`
	
	allocation := &Allocation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&allocation.ID, &allocation.TenantID, &allocation.TaskID, &allocation.WorkerID,
		&allocation.Status, &allocation.Score, &allocation.Distance, &allocation.ETA,
		&allocation.OfferedAt, &allocation.ExpiresAt, &allocation.RespondedAt, &allocation.Reason,
		&allocation.CreatedAt, &allocation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get allocation: %w", err)
	}
	
	return allocation, nil
}

func (r *postgresRepository) UpdateAllocation(ctx context.Context, allocation *Allocation) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateAllocation")
	defer span.End()
	
	allocation.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE task_allocations SET
			status = $2, responded_at = $3, reason = $4, updated_at = $5
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query,
		allocation.ID, allocation.Status, allocation.RespondedAt,
		allocation.Reason, allocation.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update allocation: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) GetAllocationsByTask(ctx context.Context, taskID uuid.UUID) ([]*Allocation, error) {
	ctx, span := tracer.Start(ctx, "repository.GetAllocationsByTask")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, task_id, worker_id, status, score, distance_km,
			eta_minutes, offered_at, expires_at, responded_at, reason,
			created_at, updated_at
		FROM task_allocations
		WHERE task_id = $1
		ORDER BY score DESC`
	
	rows, err := r.pool.Query(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("get allocations by task: %w", err)
	}
	defer rows.Close()
	
	allocations := make([]*Allocation, 0)
	for rows.Next() {
		allocation := &Allocation{}
		err := rows.Scan(
			&allocation.ID, &allocation.TenantID, &allocation.TaskID, &allocation.WorkerID,
			&allocation.Status, &allocation.Score, &allocation.Distance, &allocation.ETA,
			&allocation.OfferedAt, &allocation.ExpiresAt, &allocation.RespondedAt, &allocation.Reason,
			&allocation.CreatedAt, &allocation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}
	
	return allocations, nil
}

func (r *postgresRepository) GetAllocationsByWorker(ctx context.Context, workerID uuid.UUID, status []AllocationStatus) ([]*Allocation, error) {
	ctx, span := tracer.Start(ctx, "repository.GetAllocationsByWorker")
	defer span.End()
	
	var args []interface{}
	args = append(args, workerID)
	
	statusClause := ""
	if len(status) > 0 {
		placeholders := make([]string, len(status))
		for i, s := range status {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, s)
		}
		statusClause = fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, task_id, worker_id, status, score, distance_km,
			eta_minutes, offered_at, expires_at, responded_at, reason,
			created_at, updated_at
		FROM task_allocations
		WHERE worker_id = $1%s
		ORDER BY offered_at DESC`, statusClause)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get allocations by worker: %w", err)
	}
	defer rows.Close()
	
	allocations := make([]*Allocation, 0)
	for rows.Next() {
		allocation := &Allocation{}
		err := rows.Scan(
			&allocation.ID, &allocation.TenantID, &allocation.TaskID, &allocation.WorkerID,
			&allocation.Status, &allocation.Score, &allocation.Distance, &allocation.ETA,
			&allocation.OfferedAt, &allocation.ExpiresAt, &allocation.RespondedAt, &allocation.Reason,
			&allocation.CreatedAt, &allocation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}
	
	return allocations, nil
}

// ============================================================================
// Earnings Operations
// ============================================================================

func (r *postgresRepository) CreateEarning(ctx context.Context, earning *Earning) error {
	ctx, span := tracer.Start(ctx, "repository.CreateEarning")
	defer span.End()
	
	if earning.ID == uuid.Nil {
		earning.ID = uuid.New()
	}
	earning.CreatedAt = time.Now().UTC()
	
	metadata, _ := json.Marshal(earning.Metadata)
	
	query := `
		INSERT INTO worker_earnings (
			id, tenant_id, worker_id, task_id, type, amount, currency,
			description, payout_id, is_paid_out, metadata, earned_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	
	_, err := r.pool.Exec(ctx, query,
		earning.ID, earning.TenantID, earning.WorkerID, earning.TaskID, earning.Type,
		earning.Amount, earning.Currency, earning.Description, earning.PayoutID,
		earning.IsPaidOut, metadata, earning.EarnedAt, earning.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create earning: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetEarning(ctx context.Context, id uuid.UUID) (*Earning, error) {
	ctx, span := tracer.Start(ctx, "repository.GetEarning")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, worker_id, task_id, type, amount, currency,
			description, payout_id, is_paid_out, metadata, earned_at, created_at
		FROM worker_earnings
		WHERE id = $1`
	
	earning := &Earning{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&earning.ID, &earning.TenantID, &earning.WorkerID, &earning.TaskID, &earning.Type,
		&earning.Amount, &earning.Currency, &earning.Description, &earning.PayoutID,
		&earning.IsPaidOut, &metadata, &earning.EarnedAt, &earning.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get earning: %w", err)
	}
	
	json.Unmarshal(metadata, &earning.Metadata)
	
	return earning, nil
}

func (r *postgresRepository) GetEarningsByWorker(ctx context.Context, workerID uuid.UUID, filter EarningFilter) (*EarningList, error) {
	ctx, span := tracer.Start(ctx, "repository.GetEarningsByWorker")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("worker_id = $%d", argNum))
	args = append(args, workerID)
	argNum++
	
	if len(filter.Type) > 0 {
		placeholders := make([]string, len(filter.Type))
		for i, t := range filter.Type {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, t)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if filter.IsPaidOut != nil {
		conditions = append(conditions, fmt.Sprintf("is_paid_out = $%d", argNum))
		args = append(args, *filter.IsPaidOut)
		argNum++
	}
	
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("earned_at >= $%d", argNum))
		args = append(args, *filter.From)
		argNum++
	}
	
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("earned_at <= $%d", argNum))
		args = append(args, *filter.To)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM worker_earnings WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count earnings: %w", err)
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, worker_id, task_id, type, amount, currency,
			description, payout_id, is_paid_out, metadata, earned_at, created_at
		FROM worker_earnings
		WHERE %s
		ORDER BY earned_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get earnings: %w", err)
	}
	defer rows.Close()
	
	earnings := make([]*Earning, 0)
	for rows.Next() {
		earning := &Earning{}
		var metadata []byte
		
		err := rows.Scan(
			&earning.ID, &earning.TenantID, &earning.WorkerID, &earning.TaskID, &earning.Type,
			&earning.Amount, &earning.Currency, &earning.Description, &earning.PayoutID,
			&earning.IsPaidOut, &metadata, &earning.EarnedAt, &earning.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan earning: %w", err)
		}
		
		json.Unmarshal(metadata, &earning.Metadata)
		earnings = append(earnings, earning)
	}
	
	return &EarningList{
		Earnings: earnings,
		Total:    total,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

func (r *postgresRepository) GetEarningsSummary(ctx context.Context, workerID uuid.UUID, from, to time.Time) (*EarningsSummary, error) {
	ctx, span := tracer.Start(ctx, "repository.GetEarningsSummary")
	defer span.End()
	
	query := `
		SELECT 
			worker_id,
			COALESCE(SUM(amount), 0) as total_earnings,
			COALESCE(SUM(CASE WHEN is_paid_out = false THEN amount ELSE 0 END), 0) as pending_payout,
			COALESCE(SUM(CASE WHEN is_paid_out = true THEN amount ELSE 0 END), 0) as paid_out,
			COALESCE(SUM(CASE WHEN type = 'task_completion' THEN amount ELSE 0 END), 0) as task_completion_amount,
			COALESCE(SUM(CASE WHEN type = 'bonus' THEN amount ELSE 0 END), 0) as bonus_amount,
			COALESCE(SUM(CASE WHEN type = 'tip' THEN amount ELSE 0 END), 0) as tip_amount,
			COALESCE(SUM(CASE WHEN type = 'incentive' THEN amount ELSE 0 END), 0) as incentive_amount,
			COUNT(DISTINCT task_id) as task_count,
			MIN(currency) as currency
		FROM worker_earnings
		WHERE worker_id = $1 AND earned_at >= $2 AND earned_at <= $3
		GROUP BY worker_id`
	
	summary := &EarningsSummary{
		WorkerID: workerID,
		FromDate: from,
		ToDate:   to,
	}
	
	err := r.pool.QueryRow(ctx, query, workerID, from, to).Scan(
		&summary.WorkerID, &summary.TotalEarnings, &summary.PendingPayout,
		&summary.PaidOut, &summary.TaskCompletionAmt, &summary.BonusAmount,
		&summary.TipAmount, &summary.IncentiveAmount, &summary.TaskCount, &summary.Currency,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return empty summary
			summary.Currency = "NGN"
			return summary, nil
		}
		return nil, fmt.Errorf("get earnings summary: %w", err)
	}
	
	return summary, nil
}

// ============================================================================
// Payout Operations
// ============================================================================

func (r *postgresRepository) CreatePayout(ctx context.Context, payout *Payout) error {
	ctx, span := tracer.Start(ctx, "repository.CreatePayout")
	defer span.End()
	
	if payout.ID == uuid.Nil {
		payout.ID = uuid.New()
	}
	payout.CreatedAt = time.Now().UTC()
	payout.UpdatedAt = payout.CreatedAt
	
	metadata, _ := json.Marshal(payout.Metadata)
	
	query := `
		INSERT INTO worker_payouts (
			id, tenant_id, worker_id, amount, currency, status, method,
			bank_account_id, wallet_id, reference, provider_reference,
			failure_reason, processed_at, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	
	_, err := r.pool.Exec(ctx, query,
		payout.ID, payout.TenantID, payout.WorkerID, payout.Amount, payout.Currency,
		payout.Status, payout.Method, payout.BankAccountID, payout.WalletID,
		payout.Reference, payout.ProviderRef, payout.FailureReason, payout.ProcessedAt,
		metadata, payout.CreatedAt, payout.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create payout: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) GetPayout(ctx context.Context, id uuid.UUID) (*Payout, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPayout")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, worker_id, amount, currency, status, method,
			bank_account_id, wallet_id, reference, provider_reference,
			failure_reason, processed_at, metadata, created_at, updated_at
		FROM worker_payouts
		WHERE id = $1`
	
	payout := &Payout{}
	var metadata []byte
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&payout.ID, &payout.TenantID, &payout.WorkerID, &payout.Amount, &payout.Currency,
		&payout.Status, &payout.Method, &payout.BankAccountID, &payout.WalletID,
		&payout.Reference, &payout.ProviderRef, &payout.FailureReason, &payout.ProcessedAt,
		&metadata, &payout.CreatedAt, &payout.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get payout: %w", err)
	}
	
	json.Unmarshal(metadata, &payout.Metadata)
	
	return payout, nil
}

func (r *postgresRepository) GetPayoutsByWorker(ctx context.Context, workerID uuid.UUID, filter PayoutFilter) (*PayoutList, error) {
	ctx, span := tracer.Start(ctx, "repository.GetPayoutsByWorker")
	defer span.End()
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	conditions = append(conditions, fmt.Sprintf("worker_id = $%d", argNum))
	args = append(args, workerID)
	argNum++
	
	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argNum)
			args = append(args, s)
			argNum++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argNum))
		args = append(args, *filter.From)
		argNum++
	}
	
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argNum))
		args = append(args, *filter.To)
		argNum++
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM worker_payouts WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count payouts: %w", err)
	}
	
	query := fmt.Sprintf(`
		SELECT id, tenant_id, worker_id, amount, currency, status, method,
			bank_account_id, wallet_id, reference, provider_reference,
			failure_reason, processed_at, metadata, created_at, updated_at
		FROM worker_payouts
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)
	
	args = append(args, filter.Limit, filter.Offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get payouts: %w", err)
	}
	defer rows.Close()
	
	payouts := make([]*Payout, 0)
	for rows.Next() {
		payout := &Payout{}
		var metadata []byte
		
		err := rows.Scan(
			&payout.ID, &payout.TenantID, &payout.WorkerID, &payout.Amount, &payout.Currency,
			&payout.Status, &payout.Method, &payout.BankAccountID, &payout.WalletID,
			&payout.Reference, &payout.ProviderRef, &payout.FailureReason, &payout.ProcessedAt,
			&metadata, &payout.CreatedAt, &payout.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan payout: %w", err)
		}
		
		json.Unmarshal(metadata, &payout.Metadata)
		payouts = append(payouts, payout)
	}
	
	return &PayoutList{
		Payouts: payouts,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}

func (r *postgresRepository) UpdatePayoutStatus(ctx context.Context, id uuid.UUID, status PayoutStatus, reference string) error {
	ctx, span := tracer.Start(ctx, "repository.UpdatePayoutStatus")
	defer span.End()
	
	now := time.Now().UTC()
	
	var processedAt *time.Time
	if status == PayoutStatusCompleted || status == PayoutStatusFailed {
		processedAt = &now
	}
	
	query := `
		UPDATE worker_payouts SET
			status = $2,
			provider_reference = COALESCE($3, provider_reference),
			processed_at = COALESCE($4, processed_at),
			updated_at = $5
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id, status, reference, processedAt, now)
	if err != nil {
		return fmt.Errorf("update payout status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	// If completed, mark earnings as paid
	if status == PayoutStatusCompleted {
		_, err = r.pool.Exec(ctx,
			`UPDATE worker_earnings SET is_paid_out = true, payout_id = $1 WHERE payout_id = $1`,
			id,
		)
		if err != nil {
			return fmt.Errorf("update earnings paid status: %w", err)
		}
	}
	
	return nil
}

// ============================================================================
// Route Operations
// ============================================================================

func (r *postgresRepository) CreateRoute(ctx context.Context, route *Route) error {
	ctx, span := tracer.Start(ctx, "repository.CreateRoute")
	defer span.End()
	
	if route.ID == uuid.Nil {
		route.ID = uuid.New()
	}
	route.CreatedAt = time.Now().UTC()
	route.UpdatedAt = route.CreatedAt
	
	query := `
		INSERT INTO delivery_routes (
			id, tenant_id, worker_id, status, total_distance_km, total_duration_minutes,
			estimated_end_at, started_at, completed_at, optimization_score,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	
	_, err := r.pool.Exec(ctx, query,
		route.ID, route.TenantID, route.WorkerID, route.Status,
		route.TotalDistance, route.TotalDuration, route.EstimatedEndAt,
		route.StartedAt, route.CompletedAt, route.OptimizationScore,
		route.CreatedAt, route.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create route: %w", err)
	}
	
	// Insert stops
	for i, stop := range route.Stops {
		stop.RouteID = route.ID
		stop.Sequence = i + 1
		if err := r.AddRouteStop(ctx, route.ID, stop); err != nil {
			return err
		}
	}
	
	return nil
}

func (r *postgresRepository) GetRoute(ctx context.Context, id uuid.UUID) (*Route, error) {
	ctx, span := tracer.Start(ctx, "repository.GetRoute")
	defer span.End()
	
	query := `
		SELECT id, tenant_id, worker_id, status, total_distance_km, total_duration_minutes,
			estimated_end_at, started_at, completed_at, optimization_score,
			created_at, updated_at
		FROM delivery_routes
		WHERE id = $1`
	
	route := &Route{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&route.ID, &route.TenantID, &route.WorkerID, &route.Status,
		&route.TotalDistance, &route.TotalDuration, &route.EstimatedEndAt,
		&route.StartedAt, &route.CompletedAt, &route.OptimizationScore,
		&route.CreatedAt, &route.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get route: %w", err)
	}
	
	// Get stops
	stopsQuery := `
		SELECT id, route_id, task_id, sequence, status,
			ST_Y(location::geometry), ST_X(location::geometry), address,
			distance_from_prev_km, duration_from_prev_minutes,
			estimated_arrival, actual_arrival, completed_at, notes, metadata
		FROM route_stops
		WHERE route_id = $1
		ORDER BY sequence ASC`
	
	rows, err := r.pool.Query(ctx, stopsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("get route stops: %w", err)
	}
	defer rows.Close()
	
	route.Stops = make([]*RouteStop, 0)
	for rows.Next() {
		stop := &RouteStop{}
		var metadata []byte
		
		err := rows.Scan(
			&stop.ID, &stop.RouteID, &stop.TaskID, &stop.Sequence, &stop.Status,
			&stop.Lat, &stop.Lng, &stop.Address,
			&stop.DistanceFromPrev, &stop.DurationFromPrev,
			&stop.EstimatedArrival, &stop.ActualArrival, &stop.CompletedAt, &stop.Notes, &metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("scan route stop: %w", err)
		}
		
		json.Unmarshal(metadata, &stop.Metadata)
		route.Stops = append(route.Stops, stop)
	}
	
	return route, nil
}

func (r *postgresRepository) GetActiveRouteByWorker(ctx context.Context, workerID uuid.UUID) (*Route, error) {
	ctx, span := tracer.Start(ctx, "repository.GetActiveRouteByWorker")
	defer span.End()
	
	query := `
		SELECT id FROM delivery_routes
		WHERE worker_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1`
	
	var routeID uuid.UUID
	err := r.pool.QueryRow(ctx, query, workerID).Scan(&routeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active route: %w", err)
	}
	
	return r.GetRoute(ctx, routeID)
}

func (r *postgresRepository) UpdateRoute(ctx context.Context, route *Route) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateRoute")
	defer span.End()
	
	route.UpdatedAt = time.Now().UTC()
	
	query := `
		UPDATE delivery_routes SET
			status = $2, total_distance_km = $3, total_duration_minutes = $4,
			estimated_end_at = $5, started_at = $6, completed_at = $7,
			optimization_score = $8, updated_at = $9
		WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query,
		route.ID, route.Status, route.TotalDistance, route.TotalDuration,
		route.EstimatedEndAt, route.StartedAt, route.CompletedAt,
		route.OptimizationScore, route.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update route: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) AddRouteStop(ctx context.Context, routeID uuid.UUID, stop *RouteStop) error {
	ctx, span := tracer.Start(ctx, "repository.AddRouteStop")
	defer span.End()
	
	if stop.ID == uuid.Nil {
		stop.ID = uuid.New()
	}
	stop.RouteID = routeID
	
	metadata, _ := json.Marshal(stop.Metadata)
	
	query := `
		INSERT INTO route_stops (
			id, route_id, task_id, sequence, status, location, address,
			distance_from_prev_km, duration_from_prev_minutes,
			estimated_arrival, actual_arrival, completed_at, notes, metadata
		) VALUES (
			$1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326),
			$8, $9, $10, $11, $12, $13, $14, $15
		)`
	
	_, err := r.pool.Exec(ctx, query,
		stop.ID, stop.RouteID, stop.TaskID, stop.Sequence, stop.Status,
		stop.Lng, stop.Lat, stop.Address,
		stop.DistanceFromPrev, stop.DurationFromPrev,
		stop.EstimatedArrival, stop.ActualArrival, stop.CompletedAt, stop.Notes, metadata,
	)
	if err != nil {
		return fmt.Errorf("add route stop: %w", err)
	}
	
	return nil
}

func (r *postgresRepository) UpdateRouteStop(ctx context.Context, routeID, stopID uuid.UUID, status StopStatus) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateRouteStop")
	defer span.End()
	
	now := time.Now().UTC()
	
	var actualArrival, completedAt *time.Time
	switch status {
	case StopStatusArrived:
		actualArrival = &now
	case StopStatusCompleted, StopStatusSkipped:
		completedAt = &now
	}
	
	query := `
		UPDATE route_stops SET
			status = $3,
			actual_arrival = COALESCE($4, actual_arrival),
			completed_at = COALESCE($5, completed_at)
		WHERE route_id = $1 AND id = $2`
	
	result, err := r.pool.Exec(ctx, query, routeID, stopID, status, actualArrival, completedAt)
	if err != nil {
		return fmt.Errorf("update route stop: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

// ============================================================================
// Performance & Analytics Operations
// ============================================================================

func (r *postgresRepository) GetWorkerPerformance(ctx context.Context, workerID uuid.UUID, from, to time.Time) (*WorkerPerformance, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWorkerPerformance")
	defer span.End()
	
	query := `
		SELECT 
			$1::uuid as worker_id,
			COUNT(*) as total_tasks,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_tasks,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_tasks,
			COALESCE(
				COUNT(CASE WHEN status = 'completed' THEN 1 END)::float / NULLIF(COUNT(*), 0),
				0
			) as completion_rate,
			COALESCE(
				COUNT(CASE WHEN status = 'completed' AND completed_at <= deadline_at THEN 1 END)::float / 
				NULLIF(COUNT(CASE WHEN status = 'completed' THEN 1 END), 0),
				0
			) as on_time_rate,
			AVG(actual_duration_minutes) as avg_task_time,
			SUM(actual_distance_km) as total_distance
		FROM tasks
		WHERE assigned_worker_id = $1 AND created_at >= $2 AND created_at <= $3`
	
	performance := &WorkerPerformance{
		WorkerID: workerID,
		FromDate: from,
		ToDate:   to,
	}
	
	var avgTaskTime, totalDistance sql.NullFloat64
	
	err := r.pool.QueryRow(ctx, query, workerID, from, to).Scan(
		&performance.WorkerID,
		&performance.TotalTasks,
		&performance.CompletedTasks,
		&performance.FailedTasks,
		&performance.CancelledTasks,
		&performance.CompletionRate,
		&performance.OnTimeRate,
		&avgTaskTime,
		&totalDistance,
	)
	if err != nil {
		return nil, fmt.Errorf("get worker performance: %w", err)
	}
	
	if avgTaskTime.Valid {
		performance.AverageTaskTime = int(avgTaskTime.Float64)
	}
	if totalDistance.Valid {
		performance.TotalDistanceTraveled = totalDistance.Float64
	}
	
	// Get worker rating
	var rating float64
	err = r.pool.QueryRow(ctx, "SELECT rating FROM gig_workers WHERE id = $1", workerID).Scan(&rating)
	if err == nil {
		performance.AverageRating = rating
	}
	
	// Get earnings
	earningsSummary, _ := r.GetEarningsSummary(ctx, workerID, from, to)
	if earningsSummary != nil {
		performance.TotalEarnings = earningsSummary.TotalEarnings
	}
	
	return performance, nil
}

func (r *postgresRepository) UpdateWorkerRating(ctx context.Context, workerID uuid.UUID, rating float64) error {
	ctx, span := tracer.Start(ctx, "repository.UpdateWorkerRating")
	defer span.End()
	
	query := `UPDATE gig_workers SET rating = $2, updated_at = $3 WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, workerID, rating, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("update worker rating: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	
	return nil
}

func (r *postgresRepository) GetTaskStats(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*TaskStats, error) {
	ctx, span := tracer.Start(ctx, "repository.GetTaskStats")
	defer span.End()
	
	query := `
		SELECT 
			COUNT(*) as total_tasks,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
			COUNT(CASE WHEN status = 'assigned' THEN 1 END) as assigned_tasks,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_tasks,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_tasks,
			COALESCE(AVG(actual_duration_minutes) FILTER (WHERE status = 'completed'), 0) as avg_completion_time
		FROM tasks
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3`
	
	stats := &TaskStats{
		TenantID: tenantID,
		FromDate: from,
		ToDate:   to,
	}
	
	var avgCompletionTime float64
	
	err := r.pool.QueryRow(ctx, query, tenantID, from, to).Scan(
		&stats.TotalTasks,
		&stats.PendingTasks,
		&stats.AssignedTasks,
		&stats.InProgressTasks,
		&stats.CompletedTasks,
		&stats.FailedTasks,
		&stats.CancelledTasks,
		&avgCompletionTime,
	)
	if err != nil {
		return nil, fmt.Errorf("get task stats: %w", err)
	}
	
	stats.AvgCompletionTime = int(avgCompletionTime)
	
	return stats, nil
}

func (r *postgresRepository) GetWorkerStats(ctx context.Context, tenantID uuid.UUID) (*WorkerStats, error) {
	ctx, span := tracer.Start(ctx, "repository.GetWorkerStats")
	defer span.End()
	
	query := `
		SELECT 
			COUNT(*) as total_workers,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_workers,
			COUNT(CASE WHEN status = 'online' THEN 1 END) as online_workers,
			COUNT(CASE WHEN status = 'busy' THEN 1 END) as busy_workers
		FROM gig_workers
		WHERE tenant_id = $1`
	
	stats := &WorkerStats{
		TenantID: tenantID,
		ByType:   make(map[WorkerType]int),
	}
	
	err := r.pool.QueryRow(ctx, query, tenantID).Scan(
		&stats.TotalWorkers,
		&stats.ActiveWorkers,
		&stats.OnlineWorkers,
		&stats.BusyWorkers,
	)
	if err != nil {
		return nil, fmt.Errorf("get worker stats: %w", err)
	}
	
	// Get breakdown by type
	typeQuery := `
		SELECT type, COUNT(*) 
		FROM gig_workers 
		WHERE tenant_id = $1 
		GROUP BY type`
	
	rows, err := r.pool.Query(ctx, typeQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("get worker type breakdown: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var workerType WorkerType
		var count int
		if err := rows.Scan(&workerType, &count); err != nil {
			return nil, fmt.Errorf("scan worker type: %w", err)
		}
		stats.ByType[workerType] = count
	}
	
	return stats, nil
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
