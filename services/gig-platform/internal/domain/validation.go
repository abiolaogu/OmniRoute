// Package domain contains validation methods for gig platform domain models.
package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrWorkerIDRequired        = errors.New("worker ID is required")
	ErrTenantIDRequired        = errors.New("tenant ID is required")
	ErrUserIDRequired          = errors.New("user ID is required")
	ErrWorkerNameRequired      = errors.New("worker name is required")
	ErrWorkerPhoneRequired     = errors.New("worker phone is required")
	ErrTaskIDRequired          = errors.New("task ID is required")
	ErrOrderIDRequired         = errors.New("order ID is required")
	ErrInvalidLocation         = errors.New("invalid location coordinates")
	ErrInvalidRating           = errors.New("rating must be between 1 and 5")
	ErrInvalidPayout           = errors.New("payout must be positive")
	ErrWorkerNotVerified       = errors.New("worker is not verified")
	ErrWorkerNotAvailable      = errors.New("worker is not available")
	ErrTaskAlreadyAssigned     = errors.New("task is already assigned")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// Validate validates a GigWorker
func (w *GigWorker) Validate() error {
	if w.ID == uuid.Nil {
		return ErrWorkerIDRequired
	}
	if w.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if w.UserID == uuid.Nil {
		return ErrUserIDRequired
	}
	if w.FirstName == "" || w.LastName == "" {
		return ErrWorkerNameRequired
	}
	if w.Phone == "" {
		return ErrWorkerPhoneRequired
	}
	return nil
}

// Validate validates a Task
func (t *Task) Validate() error {
	if t.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if t.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if t.OrderID == uuid.Nil {
		return ErrOrderIDRequired
	}
	if err := t.PickupLocation.Validate(); err != nil {
		return err
	}
	if err := t.DropoffLocation.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates a Location
func (l Location) Validate() error {
	if l.Latitude < -90 || l.Latitude > 90 {
		return ErrInvalidLocation
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return ErrInvalidLocation
	}
	return nil
}

// Validate validates a TaskOffer
func (o *TaskOffer) Validate() error {
	if o.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if o.TaskID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if o.WorkerID == uuid.Nil {
		return ErrWorkerIDRequired
	}
	if o.PayoutAmount.IsNegative() {
		return ErrInvalidPayout
	}
	return nil
}

// Validate validates an Allocation
func (a *Allocation) Validate() error {
	if a.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if a.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if a.TaskID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if a.WorkerID == uuid.Nil {
		return ErrWorkerIDRequired
	}
	if a.TotalPayout.IsNegative() {
		return ErrInvalidPayout
	}
	return nil
}

// Validate validates an Earning
func (e *Earning) Validate() error {
	if e.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if e.WorkerID == uuid.Nil {
		return ErrWorkerIDRequired
	}
	if e.Amount.IsNegative() {
		return ErrInvalidPayout
	}
	if !e.IsValidType() {
		return errors.New("invalid earning type")
	}
	return nil
}

// Validate validates a Payout
func (p *Payout) Validate() error {
	if p.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if p.WorkerID == uuid.Nil {
		return ErrWorkerIDRequired
	}
	if p.Amount.IsNegative() || p.Amount.IsZero() {
		return ErrInvalidPayout
	}
	return nil
}

// Validate validates a TaskProof
func (tp *TaskProof) Validate() error {
	if tp.ID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if tp.TaskID == uuid.Nil {
		return ErrTaskIDRequired
	}
	if tp.Type == "" {
		return errors.New("proof type is required")
	}
	if tp.URL == "" {
		return errors.New("proof URL is required")
	}
	return nil
}
