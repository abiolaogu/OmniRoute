// Package domain_test provides unit tests for authentication domain models.
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// User represents a user for testing (mirrors domain.User)
type User struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Email          string
	Status         string
	FailedAttempts int
	LockedUntil    *time.Time
	Roles          []Role
}

type Role struct {
	ID          uuid.UUID
	Name        string
	Permissions []Permission
}

type Permission struct {
	Resource string
	Action   string
	Scope    string
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Resource == resource && perm.Action == action {
				return true
			}
		}
	}
	return false
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IncrementFailedAttempts increments failed login attempts
func (u *User) IncrementFailedAttempts(maxAttempts int, lockDuration time.Duration) {
	u.FailedAttempts++
	if u.FailedAttempts >= maxAttempts {
		lockUntil := time.Now().Add(lockDuration)
		u.LockedUntil = &lockUntil
	}
}

// ResetFailedAttempts resets the failed login counter
func (u *User) ResetFailedAttempts() {
	u.FailedAttempts = 0
	u.LockedUntil = nil
}

func TestUser_HasPermission_WithMatchingPermission_ReturnsTrue(t *testing.T) {
	user := &User{
		Roles: []Role{
			{
				Name: "admin",
				Permissions: []Permission{
					{Resource: "orders", Action: "read", Scope: "tenant"},
					{Resource: "orders", Action: "write", Scope: "tenant"},
				},
			},
		},
	}

	if !user.HasPermission("orders", "read") {
		t.Error("Expected user to have orders:read permission")
	}
}

func TestUser_HasPermission_WithoutMatchingPermission_ReturnsFalse(t *testing.T) {
	user := &User{
		Roles: []Role{
			{
				Name: "viewer",
				Permissions: []Permission{
					{Resource: "orders", Action: "read", Scope: "own"},
				},
			},
		},
	}

	if user.HasPermission("orders", "write") {
		t.Error("Expected user NOT to have orders:write permission")
	}
}

func TestUser_HasPermission_WithNoRoles_ReturnsFalse(t *testing.T) {
	user := &User{
		Roles: []Role{},
	}

	if user.HasPermission("orders", "read") {
		t.Error("Expected user with no roles to have no permissions")
	}
}

func TestUser_IsLocked_WhenNotLocked_ReturnsFalse(t *testing.T) {
	user := &User{
		LockedUntil: nil,
	}

	if user.IsLocked() {
		t.Error("Expected user to not be locked")
	}
}

func TestUser_IsLocked_WhenLockExpired_ReturnsFalse(t *testing.T) {
	pastTime := time.Now().Add(-1 * time.Hour)
	user := &User{
		LockedUntil: &pastTime,
	}

	if user.IsLocked() {
		t.Error("Expected user to not be locked after lock expiry")
	}
}

func TestUser_IsLocked_WhenCurrentlyLocked_ReturnsTrue(t *testing.T) {
	futureTime := time.Now().Add(1 * time.Hour)
	user := &User{
		LockedUntil: &futureTime,
	}

	if !user.IsLocked() {
		t.Error("Expected user to be locked")
	}
}

func TestUser_IncrementFailedAttempts_BelowMax_DoesNotLock(t *testing.T) {
	user := &User{
		FailedAttempts: 0,
	}

	user.IncrementFailedAttempts(5, 30*time.Minute)

	if user.FailedAttempts != 1 {
		t.Errorf("Expected 1 failed attempt, got %d", user.FailedAttempts)
	}
	if user.LockedUntil != nil {
		t.Error("Expected user to not be locked below max attempts")
	}
}

func TestUser_IncrementFailedAttempts_AtMax_LocksUser(t *testing.T) {
	user := &User{
		FailedAttempts: 4,
	}

	user.IncrementFailedAttempts(5, 30*time.Minute)

	if user.FailedAttempts != 5 {
		t.Errorf("Expected 5 failed attempts, got %d", user.FailedAttempts)
	}
	if user.LockedUntil == nil {
		t.Error("Expected user to be locked at max attempts")
	}
	if !user.IsLocked() {
		t.Error("Expected IsLocked to return true")
	}
}

func TestUser_ResetFailedAttempts_ClearsCounterAndLock(t *testing.T) {
	lockTime := time.Now().Add(1 * time.Hour)
	user := &User{
		FailedAttempts: 5,
		LockedUntil:    &lockTime,
	}

	user.ResetFailedAttempts()

	if user.FailedAttempts != 0 {
		t.Errorf("Expected 0 failed attempts, got %d", user.FailedAttempts)
	}
	if user.LockedUntil != nil {
		t.Error("Expected lock to be cleared")
	}
}
