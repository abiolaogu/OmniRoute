// Package allocation implements intelligent task-to-worker matching
// This is OmniRoute's "Uber for Delivery" brain
package allocation

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/billyronks/omniroute/gig-platform/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("allocation-engine")

// AllocationEngine handles intelligent task assignment
type AllocationEngine struct {
	workerRepo      WorkerRepository
	taskRepo        TaskRepository
	offerRepo       OfferRepository
	geoService      GeoService
	pricingEngine   PricingEngine
	notifier        WorkerNotifier
	config          AllocationConfig
	logger          *zap.Logger
	
	// Real-time worker state
	workerLocations sync.Map // map[uuid.UUID]*WorkerState
	
	// Metrics
	allocationCount int64
	avgMatchTime    time.Duration
}

// AllocationConfig holds engine configuration
type AllocationConfig struct {
	// Offer settings
	OfferTimeoutSeconds     int     `json:"offer_timeout_seconds"`
	MaxConcurrentOffers     int     `json:"max_concurrent_offers"`
	OfferBroadcastRadius    float64 `json:"offer_broadcast_radius"` // km
	
	// Matching weights
	DistanceWeight          float64 `json:"distance_weight"`
	RatingWeight            float64 `json:"rating_weight"`
	ExperienceWeight        float64 `json:"experience_weight"`
	AcceptanceRateWeight    float64 `json:"acceptance_rate_weight"`
	OnTimeRateWeight        float64 `json:"on_time_rate_weight"`
	LoadBalanceWeight       float64 `json:"load_balance_weight"`
	
	// Constraints
	MaxWorkerDistance       float64 `json:"max_worker_distance"` // km
	MinWorkerRating         float64 `json:"min_worker_rating"`
	MaxTasksPerWorker       int     `json:"max_tasks_per_worker"`
	
	// Surge pricing
	EnableSurgePricing      bool    `json:"enable_surge_pricing"`
	SurgeThreshold          float64 `json:"surge_threshold"` // Demand/supply ratio
	MaxSurgeMultiplier      float64 `json:"max_surge_multiplier"`
	
	// AI optimization
	EnableAIOptimization    bool    `json:"enable_ai_optimization"`
	MLModelEndpoint         string  `json:"ml_model_endpoint"`
}

// WorkerRepository interface for worker data access
type WorkerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GigWorker, error)
	GetAvailableInRadius(ctx context.Context, lat, lng, radiusKm float64, workerTypes []domain.WorkerType) ([]*domain.GigWorker, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.GigWorker, error)
	UpdateAvailability(ctx context.Context, workerID uuid.UUID, availability domain.WorkerAvailability) error
}

// TaskRepository interface for task data access
type TaskRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetPendingTasks(ctx context.Context, tenantID uuid.UUID) ([]*domain.Task, error)
	UpdateStatus(ctx context.Context, taskID uuid.UUID, status domain.TaskStatus) error
	AssignWorker(ctx context.Context, taskID, workerID uuid.UUID) error
}

// OfferRepository interface for offer data access
type OfferRepository interface {
	Create(ctx context.Context, offer *domain.TaskOffer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TaskOffer, error)
	GetActiveForTask(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskOffer, error)
	GetActiveForWorker(ctx context.Context, workerID uuid.UUID) ([]*domain.TaskOffer, error)
	UpdateStatus(ctx context.Context, offerID uuid.UUID, status domain.OfferStatus) error
	ExpireOldOffers(ctx context.Context) error
}

// GeoService interface for geographic calculations
type GeoService interface {
	CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 // Returns km
	CalculateETA(origin, dest domain.GeoPoint, vehicleType string) (int, error) // Returns minutes
	GetDrivingRoute(origin, dest domain.GeoPoint) ([]domain.GeoPoint, float64, error)
	IsWithinServiceArea(lat, lng float64, areaID uuid.UUID) bool
}

// PricingEngine interface for calculating earnings
type PricingEngine interface {
	CalculateTaskEarning(ctx context.Context, task *domain.Task, worker *domain.GigWorker, distance float64) (*EarningCalculation, error)
}

// EarningCalculation holds earning breakdown
type EarningCalculation struct {
	BaseEarning       decimal.Decimal `json:"base_earning"`
	DistanceEarning   decimal.Decimal `json:"distance_earning"`
	WeightEarning     decimal.Decimal `json:"weight_earning"`
	TimeEarning       decimal.Decimal `json:"time_earning"`
	SurgeMultiplier   decimal.Decimal `json:"surge_multiplier"`
	BonusEarning      decimal.Decimal `json:"bonus_earning"`
	TotalEarning      decimal.Decimal `json:"total_earning"`
}

// WorkerNotifier interface for sending notifications to workers
type WorkerNotifier interface {
	SendTaskOffer(ctx context.Context, workerID uuid.UUID, offer *domain.TaskOffer, task *domain.Task) error
	SendTaskUpdate(ctx context.Context, workerID uuid.UUID, task *domain.Task, message string) error
}

// WorkerState holds real-time worker state
type WorkerState struct {
	WorkerID       uuid.UUID
	Location       domain.GeoPoint
	Availability   domain.WorkerAvailability
	ActiveTasks    int
	LastHeartbeat  time.Time
	CurrentTaskID  *uuid.UUID
}

// NewAllocationEngine creates a new allocation engine
func NewAllocationEngine(
	workerRepo WorkerRepository,
	taskRepo TaskRepository,
	offerRepo OfferRepository,
	geoService GeoService,
	pricingEngine PricingEngine,
	notifier WorkerNotifier,
	config AllocationConfig,
	logger *zap.Logger,
) *AllocationEngine {
	return &AllocationEngine{
		workerRepo:    workerRepo,
		taskRepo:      taskRepo,
		offerRepo:     offerRepo,
		geoService:    geoService,
		pricingEngine: pricingEngine,
		notifier:      notifier,
		config:        config,
		logger:        logger,
	}
}

// AllocateTask attempts to assign a task to the best available worker
func (e *AllocationEngine) AllocateTask(ctx context.Context, taskID uuid.UUID, strategy domain.AllocationStrategy) (*AllocationResult, error) {
	ctx, span := tracer.Start(ctx, "allocation.allocate_task",
		trace.WithAttributes(
			attribute.String("task_id", taskID.String()),
			attribute.String("strategy", string(strategy)),
		))
	defer span.End()
	
	startTime := time.Now()
	
	// Get task details
	task, err := e.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	
	// Get task location
	taskLocation := task.Dropoff.Address
	if task.Pickup != nil {
		taskLocation = task.Pickup.Address
	}
	
	// Find eligible workers
	eligibleWorkers, err := e.findEligibleWorkers(ctx, task, taskLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to find eligible workers: %w", err)
	}
	
	if len(eligibleWorkers) == 0 {
		return &AllocationResult{
			Success: false,
			Message: "No eligible workers available",
		}, nil
	}
	
	// Score and rank workers based on strategy
	rankedWorkers := e.rankWorkers(ctx, eligibleWorkers, task, strategy)
	
	// Execute allocation based on strategy
	var result *AllocationResult
	switch strategy {
	case domain.AllocationStrategyNearest:
		result = e.allocateToNearest(ctx, task, rankedWorkers)
	case domain.AllocationStrategyBroadcast:
		result = e.broadcastToWorkers(ctx, task, rankedWorkers)
	case domain.AllocationStrategyAIOptimized:
		result = e.allocateWithAI(ctx, task, rankedWorkers)
	default:
		result = e.allocateToNearest(ctx, task, rankedWorkers)
	}
	
	// Record metrics
	e.allocationCount++
	duration := time.Since(startTime)
	span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))
	
	e.logger.Info("Task allocation completed",
		zap.String("task_id", taskID.String()),
		zap.String("strategy", string(strategy)),
		zap.Bool("success", result.Success),
		zap.Int("eligible_workers", len(eligibleWorkers)),
		zap.Duration("duration", duration),
	)
	
	return result, nil
}

// AllocationResult contains the result of an allocation attempt
type AllocationResult struct {
	Success       bool              `json:"success"`
	TaskID        uuid.UUID         `json:"task_id"`
	WorkerID      *uuid.UUID        `json:"worker_id,omitempty"`
	OfferID       *uuid.UUID        `json:"offer_id,omitempty"`
	Strategy      string            `json:"strategy"`
	Message       string            `json:"message,omitempty"`
	OffersCreated int               `json:"offers_created,omitempty"`
	Earning       *EarningCalculation `json:"earning,omitempty"`
}

// CandidateWorker holds worker info with scoring
type CandidateWorker struct {
	Worker        *domain.GigWorker
	Distance      float64 // km from task
	ETA           int     // minutes
	Score         float64
	Earning       *EarningCalculation
	ScoreBreakdown map[string]float64
}

// findEligibleWorkers finds workers who can handle the task
func (e *AllocationEngine) findEligibleWorkers(ctx context.Context, task *domain.Task, location domain.Address) ([]*CandidateWorker, error) {
	ctx, span := tracer.Start(ctx, "allocation.find_eligible_workers")
	defer span.End()
	
	// Determine required worker types
	workerTypes := e.determineWorkerTypes(task)
	
	// Find workers within radius
	workers, err := e.workerRepo.GetAvailableInRadius(
		ctx,
		location.Latitude,
		location.Longitude,
		e.config.MaxWorkerDistance,
		workerTypes,
	)
	if err != nil {
		return nil, err
	}
	
	candidates := make([]*CandidateWorker, 0, len(workers))
	
	for _, worker := range workers {
		// Check eligibility constraints
		if !e.isWorkerEligible(worker, task) {
			continue
		}
		
		// Get real-time worker state
		state := e.getWorkerState(worker.ID)
		if state == nil || state.Availability != domain.WorkerAvailabilityOnline {
			continue
		}
		
		// Calculate distance
		distance := e.geoService.CalculateDistance(
			state.Location.Latitude, state.Location.Longitude,
			location.Latitude, location.Longitude,
		)
		
		if distance > e.config.MaxWorkerDistance {
			continue
		}
		
		// Calculate ETA
		eta, err := e.geoService.CalculateETA(state.Location, domain.GeoPoint{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		}, e.getVehicleType(worker))
		if err != nil {
			e.logger.Warn("Failed to calculate ETA", zap.Error(err))
			eta = int(distance * 3) // Fallback estimate: 3 min per km
		}
		
		// Calculate potential earning
		earning, err := e.pricingEngine.CalculateTaskEarning(ctx, task, worker, distance)
		if err != nil {
			e.logger.Warn("Failed to calculate earning", zap.Error(err))
			continue
		}
		
		candidates = append(candidates, &CandidateWorker{
			Worker:   worker,
			Distance: distance,
			ETA:      eta,
			Earning:  earning,
		})
	}
	
	span.SetAttributes(attribute.Int("candidates_found", len(candidates)))
	return candidates, nil
}

// isWorkerEligible checks if a worker meets task requirements
func (e *AllocationEngine) isWorkerEligible(worker *domain.GigWorker, task *domain.Task) bool {
	// Check worker status
	if worker.Status != domain.WorkerStatusActive {
		return false
	}
	
	// Check verification
	if worker.VerificationStatus != domain.VerificationApproved {
		return false
	}
	
	// Check rating
	if worker.Rating < e.config.MinWorkerRating {
		return false
	}
	
	// Check task type preference
	if !e.workerAcceptsTaskType(worker, task.Type) {
		return false
	}
	
	// Check capacity (for delivery tasks)
	if task.TotalWeight > 0 && worker.Vehicle != nil {
		if task.TotalWeight > worker.Vehicle.Capacity {
			return false
		}
	}
	
	// Check if worker accepts COD
	if task.CollectionAmount.GreaterThan(decimal.Zero) && !worker.TaskPreferences.AcceptCOD {
		return false
	}
	
	return true
}

// workerAcceptsTaskType checks if worker handles this task type
func (e *AllocationEngine) workerAcceptsTaskType(worker *domain.GigWorker, taskType domain.TaskType) bool {
	if len(worker.TaskPreferences.PreferredTaskTypes) == 0 {
		return true // No preference means accepts all
	}
	
	for _, t := range worker.TaskPreferences.PreferredTaskTypes {
		if t == taskType {
			return true
		}
	}
	
	return false
}

// determineWorkerTypes determines what types of workers can handle the task
func (e *AllocationEngine) determineWorkerTypes(task *domain.Task) []domain.WorkerType {
	switch task.Type {
	case domain.TaskTypeDelivery:
		if task.TotalWeight > 50 {
			return []domain.WorkerType{domain.WorkerTypeDriver}
		}
		if task.TotalWeight > 10 {
			return []domain.WorkerType{domain.WorkerTypeDriver, domain.WorkerTypeRider}
		}
		return []domain.WorkerType{domain.WorkerTypeDriver, domain.WorkerTypeRider, domain.WorkerTypeCyclist}
	
	case domain.TaskTypeCollection:
		return []domain.WorkerType{domain.WorkerTypeCollector, domain.WorkerTypeDriver, domain.WorkerTypeRider}
	
	case domain.TaskTypeSurvey:
		return []domain.WorkerType{domain.WorkerTypeSurveyor, domain.WorkerTypeWalker}
	
	case domain.TaskTypeMerchandising:
		return []domain.WorkerType{domain.WorkerTypeMerchandiser, domain.WorkerTypeWalker}
	
	default:
		return []domain.WorkerType{domain.WorkerTypeDriver, domain.WorkerTypeRider}
	}
}

// rankWorkers scores and ranks workers based on strategy
func (e *AllocationEngine) rankWorkers(ctx context.Context, candidates []*CandidateWorker, task *domain.Task, strategy domain.AllocationStrategy) []*CandidateWorker {
	for _, candidate := range candidates {
		candidate.Score, candidate.ScoreBreakdown = e.calculateWorkerScore(candidate, task, strategy)
	}
	
	// Sort by score (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})
	
	return candidates
}

// calculateWorkerScore computes a worker's suitability score
func (e *AllocationEngine) calculateWorkerScore(candidate *CandidateWorker, task *domain.Task, strategy domain.AllocationStrategy) (float64, map[string]float64) {
	breakdown := make(map[string]float64)
	worker := candidate.Worker
	
	// Distance score (closer is better) - normalize to 0-100
	maxDist := e.config.MaxWorkerDistance
	distScore := 100 * (1 - math.Min(candidate.Distance/maxDist, 1))
	breakdown["distance"] = distScore
	
	// Rating score - normalize to 0-100
	ratingScore := worker.Rating * 20 // 5.0 rating = 100
	breakdown["rating"] = ratingScore
	
	// Experience score (completed tasks) - log scale
	expScore := math.Min(100, math.Log(float64(worker.CompletedTasks+1))*20)
	breakdown["experience"] = expScore
	
	// Acceptance rate score
	acceptScore := worker.AcceptanceRate * 100
	breakdown["acceptance"] = acceptScore
	
	// On-time rate score
	onTimeScore := worker.OnTimeRate * 100
	breakdown["on_time"] = onTimeScore
	
	// Load balance score (fewer active tasks is better)
	state := e.getWorkerState(worker.ID)
	activeTasks := 0
	if state != nil {
		activeTasks = state.ActiveTasks
	}
	loadScore := 100 * (1 - float64(activeTasks)/float64(e.config.MaxTasksPerWorker))
	breakdown["load_balance"] = loadScore
	
	// Calculate weighted score
	var totalScore float64
	
	switch strategy {
	case domain.AllocationStrategyNearest:
		// Prioritize distance
		totalScore = distScore*0.5 + ratingScore*0.2 + acceptScore*0.15 + onTimeScore*0.15
	
	case domain.AllocationStrategyBestRated:
		// Prioritize rating
		totalScore = distScore*0.2 + ratingScore*0.4 + expScore*0.2 + acceptScore*0.1 + onTimeScore*0.1
	
	case domain.AllocationStrategyLoadBalanced:
		// Prioritize even distribution
		totalScore = distScore*0.2 + ratingScore*0.2 + loadScore*0.4 + acceptScore*0.1 + onTimeScore*0.1
	
	default:
		// Default balanced approach
		totalScore = distScore*e.config.DistanceWeight +
			ratingScore*e.config.RatingWeight +
			expScore*e.config.ExperienceWeight +
			acceptScore*e.config.AcceptanceRateWeight +
			onTimeScore*e.config.OnTimeRateWeight +
			loadScore*e.config.LoadBalanceWeight
	}
	
	return totalScore, breakdown
}

// allocateToNearest assigns to the highest-scoring worker directly
func (e *AllocationEngine) allocateToNearest(ctx context.Context, task *domain.Task, candidates []*CandidateWorker) *AllocationResult {
	if len(candidates) == 0 {
		return &AllocationResult{
			Success: false,
			TaskID:  task.ID,
			Message: "No eligible workers",
		}
	}
	
	// Take the best candidate
	best := candidates[0]
	
	// Create offer
	offer := &domain.TaskOffer{
		ID:            uuid.New(),
		TaskID:        task.ID,
		WorkerID:      best.Worker.ID,
		OfferedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(time.Duration(e.config.OfferTimeoutSeconds) * time.Second),
		Status:        domain.OfferStatusPending,
		BaseEarning:   best.Earning.BaseEarning,
		BonusEarning:  best.Earning.BonusEarning,
		EstimatedTime: best.ETA,
		Distance:      best.Distance,
	}
	
	if err := e.offerRepo.Create(ctx, offer); err != nil {
		e.logger.Error("Failed to create offer", zap.Error(err))
		return &AllocationResult{
			Success: false,
			TaskID:  task.ID,
			Message: "Failed to create offer",
		}
	}
	
	// Notify worker
	if err := e.notifier.SendTaskOffer(ctx, best.Worker.ID, offer, task); err != nil {
		e.logger.Error("Failed to notify worker", zap.Error(err))
	}
	
	// Update task status
	e.taskRepo.UpdateStatus(ctx, task.ID, domain.TaskStatusOffered)
	
	return &AllocationResult{
		Success:       true,
		TaskID:        task.ID,
		WorkerID:      &best.Worker.ID,
		OfferID:       &offer.ID,
		Strategy:      "nearest",
		OffersCreated: 1,
		Earning:       best.Earning,
	}
}

// broadcastToWorkers sends offers to multiple workers simultaneously
func (e *AllocationEngine) broadcastToWorkers(ctx context.Context, task *domain.Task, candidates []*CandidateWorker) *AllocationResult {
	if len(candidates) == 0 {
		return &AllocationResult{
			Success: false,
			TaskID:  task.ID,
			Message: "No eligible workers",
		}
	}
	
	// Limit broadcast count
	maxBroadcast := e.config.MaxConcurrentOffers
	if len(candidates) < maxBroadcast {
		maxBroadcast = len(candidates)
	}
	
	var wg sync.WaitGroup
	offersCreated := 0
	
	for i := 0; i < maxBroadcast; i++ {
		candidate := candidates[i]
		wg.Add(1)
		
		go func(c *CandidateWorker) {
			defer wg.Done()
			
			offer := &domain.TaskOffer{
				ID:            uuid.New(),
				TaskID:        task.ID,
				WorkerID:      c.Worker.ID,
				OfferedAt:     time.Now(),
				ExpiresAt:     time.Now().Add(time.Duration(e.config.OfferTimeoutSeconds) * time.Second),
				Status:        domain.OfferStatusPending,
				BaseEarning:   c.Earning.BaseEarning,
				BonusEarning:  c.Earning.BonusEarning,
				EstimatedTime: c.ETA,
				Distance:      c.Distance,
			}
			
			if err := e.offerRepo.Create(ctx, offer); err != nil {
				e.logger.Error("Failed to create broadcast offer", zap.Error(err))
				return
			}
			
			if err := e.notifier.SendTaskOffer(ctx, c.Worker.ID, offer, task); err != nil {
				e.logger.Error("Failed to notify worker in broadcast", zap.Error(err))
			}
			
			offersCreated++
		}(candidate)
	}
	
	wg.Wait()
	
	// Update task status
	e.taskRepo.UpdateStatus(ctx, task.ID, domain.TaskStatusOffered)
	
	return &AllocationResult{
		Success:       offersCreated > 0,
		TaskID:        task.ID,
		Strategy:      "broadcast",
		OffersCreated: offersCreated,
		Message:       fmt.Sprintf("Broadcasted to %d workers", offersCreated),
	}
}

// allocateWithAI uses ML model for optimal allocation
func (e *AllocationEngine) allocateWithAI(ctx context.Context, task *domain.Task, candidates []*CandidateWorker) *AllocationResult {
	if !e.config.EnableAIOptimization {
		return e.allocateToNearest(ctx, task, candidates)
	}
	
	// TODO: Call ML model endpoint for optimal worker selection
	// For now, fall back to standard allocation
	return e.allocateToNearest(ctx, task, candidates)
}

// AcceptOffer handles a worker accepting a task offer
func (e *AllocationEngine) AcceptOffer(ctx context.Context, offerID uuid.UUID, workerID uuid.UUID) (*domain.Task, error) {
	ctx, span := tracer.Start(ctx, "allocation.accept_offer")
	defer span.End()
	
	// Get offer
	offer, err := e.offerRepo.GetByID(ctx, offerID)
	if err != nil {
		return nil, fmt.Errorf("offer not found: %w", err)
	}
	
	// Validate
	if offer.WorkerID != workerID {
		return nil, fmt.Errorf("offer not for this worker")
	}
	
	if offer.Status != domain.OfferStatusPending {
		return nil, fmt.Errorf("offer is no longer pending")
	}
	
	if time.Now().After(offer.ExpiresAt) {
		return nil, fmt.Errorf("offer has expired")
	}
	
	// Check if task still available (first-come-first-served)
	task, err := e.taskRepo.GetByID(ctx, offer.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	
	if task.Status != domain.TaskStatusOffered && task.Status != domain.TaskStatusPending {
		return nil, fmt.Errorf("task is no longer available")
	}
	
	// Accept the offer
	now := time.Now()
	offer.Status = domain.OfferStatusAccepted
	offer.RespondedAt = &now
	
	if err := e.offerRepo.UpdateStatus(ctx, offerID, domain.OfferStatusAccepted); err != nil {
		return nil, fmt.Errorf("failed to update offer: %w", err)
	}
	
	// Assign worker to task
	if err := e.taskRepo.AssignWorker(ctx, task.ID, workerID); err != nil {
		return nil, fmt.Errorf("failed to assign worker: %w", err)
	}
	
	// Update task status
	task.Status = domain.TaskStatusAccepted
	task.AssignedWorkerID = &workerID
	task.AssignedAt = &now
	task.AcceptedAt = &now
	
	if err := e.taskRepo.UpdateStatus(ctx, task.ID, domain.TaskStatusAccepted); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	
	// Cancel other pending offers for this task
	e.cancelOtherOffers(ctx, task.ID, offerID)
	
	// Update worker state
	e.updateWorkerState(workerID, func(state *WorkerState) {
		state.ActiveTasks++
		state.CurrentTaskID = &task.ID
	})
	
	e.logger.Info("Offer accepted",
		zap.String("offer_id", offerID.String()),
		zap.String("worker_id", workerID.String()),
		zap.String("task_id", task.ID.String()),
	)
	
	return task, nil
}

// DeclineOffer handles a worker declining a task offer
func (e *AllocationEngine) DeclineOffer(ctx context.Context, offerID uuid.UUID, workerID uuid.UUID, reason string) error {
	offer, err := e.offerRepo.GetByID(ctx, offerID)
	if err != nil {
		return err
	}
	
	if offer.WorkerID != workerID {
		return fmt.Errorf("offer not for this worker")
	}
	
	now := time.Now()
	offer.Status = domain.OfferStatusDeclined
	offer.RespondedAt = &now
	offer.DeclineReason = reason
	
	if err := e.offerRepo.UpdateStatus(ctx, offerID, domain.OfferStatusDeclined); err != nil {
		return err
	}
	
	// Check if all offers for this task are declined/expired
	activeOffers, err := e.offerRepo.GetActiveForTask(ctx, offer.TaskID)
	if err != nil {
		return err
	}
	
	if len(activeOffers) == 0 {
		// Re-allocate the task
		go func() {
			ctx := context.Background()
			_, err := e.AllocateTask(ctx, offer.TaskID, domain.AllocationStrategyBroadcast)
			if err != nil {
				e.logger.Error("Failed to re-allocate task", zap.Error(err))
			}
		}()
	}
	
	return nil
}

// cancelOtherOffers cancels all other pending offers for a task
func (e *AllocationEngine) cancelOtherOffers(ctx context.Context, taskID, acceptedOfferID uuid.UUID) {
	offers, err := e.offerRepo.GetActiveForTask(ctx, taskID)
	if err != nil {
		e.logger.Error("Failed to get offers for cancellation", zap.Error(err))
		return
	}
	
	for _, offer := range offers {
		if offer.ID != acceptedOfferID && offer.Status == domain.OfferStatusPending {
			e.offerRepo.UpdateStatus(ctx, offer.ID, domain.OfferStatusCancelled)
		}
	}
}

// UpdateWorkerLocation updates a worker's real-time location
func (e *AllocationEngine) UpdateWorkerLocation(workerID uuid.UUID, location domain.GeoPoint) {
	e.updateWorkerState(workerID, func(state *WorkerState) {
		state.Location = location
		state.LastHeartbeat = time.Now()
	})
}

// SetWorkerAvailability updates a worker's availability status
func (e *AllocationEngine) SetWorkerAvailability(workerID uuid.UUID, availability domain.WorkerAvailability) {
	e.updateWorkerState(workerID, func(state *WorkerState) {
		state.Availability = availability
	})
}

// getWorkerState retrieves worker's real-time state
func (e *AllocationEngine) getWorkerState(workerID uuid.UUID) *WorkerState {
	if value, ok := e.workerLocations.Load(workerID); ok {
		return value.(*WorkerState)
	}
	return nil
}

// updateWorkerState updates worker state with a function
func (e *AllocationEngine) updateWorkerState(workerID uuid.UUID, updateFn func(*WorkerState)) {
	value, _ := e.workerLocations.LoadOrStore(workerID, &WorkerState{
		WorkerID:      workerID,
		Availability:  domain.WorkerAvailabilityOffline,
		ActiveTasks:   0,
		LastHeartbeat: time.Now(),
	})
	state := value.(*WorkerState)
	updateFn(state)
}

// getVehicleType returns the worker's vehicle type
func (e *AllocationEngine) getVehicleType(worker *domain.GigWorker) string {
	if worker.Vehicle != nil {
		return worker.Vehicle.Type
	}
	return "motorcycle" // Default
}

// import statement was missed
import "fmt"
