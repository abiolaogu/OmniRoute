// Package domain_test contains unit tests for the Gig Platform domain models
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/omniroute/gig-platform/internal/domain"
)

func TestWorkerStatus_Transitions(t *testing.T) {
	tests := []struct {
		name    string
		from    domain.WorkerStatus
		to      domain.WorkerStatus
		allowed bool
	}{
		{"offline to available", domain.WorkerStatusOffline, domain.WorkerStatusAvailable, true},
		{"available to busy", domain.WorkerStatusAvailable, domain.WorkerStatusBusy, true},
		{"busy to available", domain.WorkerStatusBusy, domain.WorkerStatusAvailable, true},
		{"available to in_transit", domain.WorkerStatusAvailable, domain.WorkerStatusInTransit, true},
		{"on_break to available", domain.WorkerStatusOnBreak, domain.WorkerStatusAvailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := tt.from.CanTransitionTo(tt.to)
			if allowed != tt.allowed {
				t.Errorf("WorkerStatus.CanTransitionTo(%v, %v) = %v, want %v",
					tt.from, tt.to, allowed, tt.allowed)
			}
		})
	}
}

func TestWorkerLevel_String(t *testing.T) {
	levels := []struct {
		level domain.WorkerLevel
		want  string
	}{
		{domain.WorkerLevelStarter, "starter"},
		{domain.WorkerLevelBronze, "bronze"},
		{domain.WorkerLevelSilver, "silver"},
		{domain.WorkerLevelGold, "gold"},
		{domain.WorkerLevelDiamond, "diamond"},
		{domain.WorkerLevelMaster, "master"},
	}

	for _, tt := range levels {
		t.Run(tt.want, func(t *testing.T) {
			if got := string(tt.level); got != tt.want {
				t.Errorf("WorkerLevel = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGigWorker_CanAcceptTask(t *testing.T) {
	tests := []struct {
		name     string
		worker   domain.GigWorker
		taskType domain.TaskType
		expected bool
	}{
		{
			name: "delivery worker can accept delivery task",
			worker: domain.GigWorker{
				Status:     domain.WorkerStatusAvailable,
				Type:       domain.WorkerTypeDelivery,
				IsVerified: true,
			},
			taskType: domain.TaskTypeDelivery,
			expected: true,
		},
		{
			name: "busy worker cannot accept task",
			worker: domain.GigWorker{
				Status:     domain.WorkerStatusBusy,
				Type:       domain.WorkerTypeDelivery,
				IsVerified: true,
			},
			taskType: domain.TaskTypeDelivery,
			expected: false,
		},
		{
			name: "unverified worker cannot accept task",
			worker: domain.GigWorker{
				Status:     domain.WorkerStatusAvailable,
				Type:       domain.WorkerTypeDelivery,
				IsVerified: false,
			},
			taskType: domain.TaskTypeDelivery,
			expected: false,
		},
		{
			name: "multi_role worker can accept any task",
			worker: domain.GigWorker{
				Status:     domain.WorkerStatusAvailable,
				Type:       domain.WorkerTypeMultiRole,
				IsVerified: true,
			},
			taskType: domain.TaskTypeCollection,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.worker.CanAcceptTask(tt.taskType); got != tt.expected {
				t.Errorf("GigWorker.CanAcceptTask() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGigWorker_UpdateRating(t *testing.T) {
	worker := &domain.GigWorker{
		ID:             uuid.New(),
		Rating:         decimal.NewFromFloat(4.5),
		TotalTasks:     100,
		CompletedTasks: 95,
	}

	// Add a new 5-star rating
	newRating := decimal.NewFromFloat(5.0)
	worker.AddRating(newRating)

	// Rating should be recalculated
	if worker.Rating.LessThan(decimal.NewFromFloat(4.5)) {
		t.Error("Rating should increase after 5-star addition")
	}
}

func TestGigWorker_CalculateSuccessRate(t *testing.T) {
	worker := &domain.GigWorker{
		TotalTasks:     100,
		CompletedTasks: 95,
	}

	expectedRate := decimal.NewFromFloat(95.0)
	if rate := worker.CalculateSuccessRate(); !rate.Equal(expectedRate) {
		t.Errorf("CalculateSuccessRate() = %v, want %v", rate, expectedRate)
	}

	// Test zero tasks
	worker.TotalTasks = 0
	if rate := worker.CalculateSuccessRate(); !rate.Equal(decimal.Zero) {
		t.Errorf("CalculateSuccessRate() with zero tasks = %v, want 0", rate)
	}
}

func TestTask_StatusTransitions(t *testing.T) {
	tests := []struct {
		name    string
		from    domain.TaskStatus
		to      domain.TaskStatus
		allowed bool
	}{
		{"pending to assigned", domain.TaskStatusPending, domain.TaskStatusAssigned, true},
		{"assigned to accepted", domain.TaskStatusAssigned, domain.TaskStatusAccepted, true},
		{"accepted to in_progress", domain.TaskStatusAccepted, domain.TaskStatusInProgress, true},
		{"in_progress to completed", domain.TaskStatusInProgress, domain.TaskStatusCompleted, true},
		{"completed to pending", domain.TaskStatusCompleted, domain.TaskStatusPending, false},
		{"pending to completed", domain.TaskStatusPending, domain.TaskStatusCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &domain.Task{Status: tt.from}
			allowed := task.CanTransitionTo(tt.to)
			if allowed != tt.allowed {
				t.Errorf("Task.CanTransitionTo(%v, %v) = %v, want %v",
					tt.from, tt.to, allowed, tt.allowed)
			}
		})
	}
}

func TestTask_CalculatePayout(t *testing.T) {
	task := &domain.Task{
		BasePayout:      decimal.NewFromFloat(500.0),
		BonusPayout:     decimal.NewFromFloat(100.0),
		SurgeMultiplier: decimal.NewFromFloat(1.5),
	}

	// Expected: (500 + 100) * 1.5 = 900
	expected := decimal.NewFromFloat(900.0)
	if payout := task.CalculateTotalPayout(); !payout.Equal(expected) {
		t.Errorf("Task.CalculateTotalPayout() = %v, want %v", payout, expected)
	}
}

func TestTask_IsOverdue(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		deadline *time.Time
		status   domain.TaskStatus
		expected bool
	}{
		{
			name:     "no deadline set",
			deadline: nil,
			status:   domain.TaskStatusPending,
			expected: false,
		},
		{
			name:     "deadline in future",
			deadline: ptr(now.Add(24 * time.Hour)),
			status:   domain.TaskStatusPending,
			expected: false,
		},
		{
			name:     "deadline passed",
			deadline: ptr(now.Add(-1 * time.Hour)),
			status:   domain.TaskStatusPending,
			expected: true,
		},
		{
			name:     "completed task not overdue",
			deadline: ptr(now.Add(-1 * time.Hour)),
			status:   domain.TaskStatusCompleted,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &domain.Task{
				Deadline: tt.deadline,
				Status:   tt.status,
			}
			if got := task.IsOverdue(); got != tt.expected {
				t.Errorf("Task.IsOverdue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTaskOffer_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		expiresAt time.Time
		status    domain.OfferStatus
		expected  bool
	}{
		{
			name:      "not expired",
			expiresAt: now.Add(5 * time.Minute),
			status:    domain.OfferStatusPending,
			expected:  false,
		},
		{
			name:      "expired by time",
			expiresAt: now.Add(-1 * time.Minute),
			status:    domain.OfferStatusPending,
			expected:  true,
		},
		{
			name:      "already accepted",
			expiresAt: now.Add(-1 * time.Minute),
			status:    domain.OfferStatusAccepted,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offer := &domain.TaskOffer{
				ExpiresAt: tt.expiresAt,
				Status:    tt.status,
			}
			if got := offer.IsExpired(); got != tt.expected {
				t.Errorf("TaskOffer.IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLocation_DistanceTo(t *testing.T) {
	// Test distance calculation between two points
	// Lagos to Abuja approx 450km
	lagos := domain.Location{
		Latitude:  6.5244,
		Longitude: 3.3792,
	}

	abuja := domain.Location{
		Latitude:  9.0765,
		Longitude: 7.3986,
	}

	distance := lagos.DistanceTo(abuja)

	// Should be approximately 450km (allow some variance)
	if distance < 400 || distance > 500 {
		t.Errorf("Location.DistanceTo() = %vkm, expected ~450km", distance)
	}
}

func TestEarning_ValidTypes(t *testing.T) {
	validTypes := []string{"task_payout", "bonus", "tip", "penalty", "referral"}

	for _, typ := range validTypes {
		earning := &domain.Earning{Type: typ}
		if !earning.IsValidType() {
			t.Errorf("Earning.IsValidType() for %s = false, want true", typ)
		}
	}

	// Invalid type
	earning := &domain.Earning{Type: "invalid_type"}
	if earning.IsValidType() {
		t.Error("Earning.IsValidType() for invalid_type = true, want false")
	}
}

// Helper function to create pointer to time
func ptr(t time.Time) *time.Time {
	return &t
}
