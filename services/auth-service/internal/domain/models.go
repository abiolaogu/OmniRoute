// Package domain provides DDD domain models for the authentication service.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// User represents a user in the system (Aggregate Root)
type User struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	Email            string
	Phone            string
	PasswordHash     string
	FirstName        string
	LastName         string
	Status           UserStatus
	EmailVerified    bool
	PhoneVerified    bool
	TwoFactorEnabled bool
	TwoFactorSecret  string
	LastLoginAt      *time.Time
	LastLoginIP      string
	FailedAttempts   int
	LockedUntil      *time.Time
	Roles            []Role
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Role represents a role that can be assigned to users
type Role struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	IsSystem    bool // System roles cannot be modified
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission represents an individual permission
type Permission struct {
	ID          uuid.UUID
	Resource    string // e.g., "orders", "products"
	Action      string // e.g., "read", "write", "delete"
	Scope       string // e.g., "own", "tenant", "all"
	Description string
}

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Domain    string
	LogoURL   string
	Settings  TenantSettings
	Status    TenantStatus
	Plan      TenantPlan
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TenantStatus represents tenant status
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusDeleted   TenantStatus = "deleted"
)

// TenantPlan represents subscription plan
type TenantPlan string

const (
	TenantPlanFree       TenantPlan = "free"
	TenantPlanStarter    TenantPlan = "starter"
	TenantPlanPro        TenantPlan = "pro"
	TenantPlanEnterprise TenantPlan = "enterprise"
)

// TenantSettings holds tenant-specific configurations
type TenantSettings struct {
	Currency          string
	Timezone          string
	Language          string
	DateFormat        string
	EnabledFeatures   []string
	CustomBranding    bool
	WhiteLabelEnabled bool
	APIRateLimit      int
}

// Session represents an active user session
type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	TenantID     uuid.UUID
	AccessToken  string
	RefreshToken string
	DeviceInfo   DeviceInfo
	IPAddress    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// DeviceInfo stores information about the client device
type DeviceInfo struct {
	UserAgent   string
	DeviceType  string
	OS          string
	Browser     string
	IsMobile    bool
	Fingerprint string
}

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	UserID      uuid.UUID
	Name        string
	KeyHash     string // Hashed key, actual key shown only on creation
	Prefix      string // First few chars for identification
	Permissions []string
	LastUsedAt  *time.Time
	ExpiresAt   *time.Time
	CreatedAt   time.Time
}

// OAuthClient represents an OAuth 2.0 client application
type OAuthClient struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ClientID       string
	ClientSecret   string
	Name           string
	RedirectURIs   []string
	GrantTypes     []string
	Scopes         []string
	IsConfidential bool
	CreatedAt      time.Time
}

// AuditLog represents an authentication-related audit event
type AuditLog struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	UserID     *uuid.UUID
	Action     string
	Resource   string
	ResourceID string
	IPAddress  string
	UserAgent  string
	Details    map[string]interface{}
	Success    bool
	CreatedAt  time.Time
}

// Credentials for authentication
type Credentials struct {
	Email    string
	Password string
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// JWTClaims represents claims in a JWT token
type JWTClaims struct {
	UserID      uuid.UUID `json:"sub"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Email       string    `json:"email"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	SessionID   uuid.UUID `json:"session_id"`
	IssuedAt    int64     `json:"iat"`
	ExpiresAt   int64     `json:"exp"`
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
