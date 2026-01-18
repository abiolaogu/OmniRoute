// Package domain contains the core domain models for the Notification Service.
// Following DDD principles with aggregates, entities, and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Domain Errors
// ============================================================================

var (
	ErrRecipientRequired    = errors.New("recipient address is required")
	ErrBodyRequired         = errors.New("message body is required")
	ErrTemplateNameRequired = errors.New("template name is required")
	ErrTemplateBodyRequired = errors.New("template body is required")
	ErrVariableNameRequired = errors.New("variable name is required")
	ErrInvalidChannel       = errors.New("invalid notification channel")
	ErrSessionExpired       = errors.New("USSD session has expired")
)

// ============================================================================
// Value Objects
// ============================================================================

// NotificationChannel represents the delivery channel
type NotificationChannel string

const (
	ChannelWhatsApp NotificationChannel = "whatsapp"
	ChannelSMS      NotificationChannel = "sms"
	ChannelEmail    NotificationChannel = "email"
	ChannelPush     NotificationChannel = "push"
	ChannelVoice    NotificationChannel = "voice"
	ChannelUSSD     NotificationChannel = "ussd"
)

// NotificationStatus represents the lifecycle status
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusQueued    NotificationStatus = "queued"
	StatusSending   NotificationStatus = "sending"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusRead      NotificationStatus = "read"
	StatusFailed    NotificationStatus = "failed"
	StatusCancelled NotificationStatus = "cancelled"
)

// NotificationPriority represents the urgency level
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// TemplateType represents the type of message template
type TemplateType string

const (
	TemplateTypeTransactional TemplateType = "transactional"
	TemplateTypeMarketing     TemplateType = "marketing"
	TemplateTypeOTP           TemplateType = "otp"
	TemplateTypeReminder      TemplateType = "reminder"
	TemplateTypeAlert         TemplateType = "alert"
)

// SessionStatus for USSD sessions
type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusEnded   SessionStatus = "ended"
	SessionStatusTimeout SessionStatus = "timeout"
)

// ============================================================================
// Aggregates and Entities
// ============================================================================

// Notification is the aggregate root for message management
type Notification struct {
	ID       uuid.UUID            `json:"id"`
	TenantID uuid.UUID            `json:"tenant_id"`
	Channel  NotificationChannel  `json:"channel"`
	Status   NotificationStatus   `json:"status"`
	Priority NotificationPriority `json:"priority"`

	// Recipient
	RecipientID      uuid.UUID `json:"recipient_id"`
	RecipientType    string    `json:"recipient_type"`
	RecipientAddress string    `json:"recipient_address"`

	// Content
	TemplateID *uuid.UUID             `json:"template_id,omitempty"`
	Subject    string                 `json:"subject,omitempty"`
	Body       string                 `json:"body"`
	Data       map[string]interface{} `json:"data,omitempty"`

	// Metadata
	CorrelationID string `json:"correlation_id,omitempty"`
	ExternalID    string `json:"external_id,omitempty"`

	// Timing
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`

	// Retry
	Attempts    int    `json:"attempts"`
	MaxAttempts int    `json:"max_attempts"`
	LastError   string `json:"last_error,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeliveryStatus provides detailed delivery information
type DeliveryStatus struct {
	NotificationID uuid.UUID          `json:"notification_id"`
	Status         NotificationStatus `json:"status"`
	ProviderStatus string             `json:"provider_status,omitempty"`
	DeliveredAt    *time.Time         `json:"delivered_at,omitempty"`
	ReadAt         *time.Time         `json:"read_at,omitempty"`
	FailureReason  string             `json:"failure_reason,omitempty"`
}

// Template represents a message template
type Template struct {
	ID       uuid.UUID           `json:"id"`
	TenantID uuid.UUID           `json:"tenant_id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	Type     TemplateType        `json:"type"`
	Channel  NotificationChannel `json:"channel"`

	// Content
	Subject   string             `json:"subject,omitempty"`
	Body      string             `json:"body"`
	Variables []TemplateVariable `json:"variables"`

	// Localization
	Language     string            `json:"language"`
	Translations map[string]string `json:"translations,omitempty"`

	// Status
	IsActive   bool       `json:"is_active"`
	IsApproved bool       `json:"is_approved"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TemplateVariable defines a variable in a template
type TemplateVariable struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"default_value,omitempty"`
}

// USSDSession represents an active USSD session
type USSDSession struct {
	ID          string        `json:"id"`
	TenantID    uuid.UUID     `json:"tenant_id"`
	PhoneNumber string        `json:"phone_number"`
	ServiceCode string        `json:"service_code"`
	Status      SessionStatus `json:"status"`

	// Navigation
	CurrentMenuID string   `json:"current_menu_id"`
	MenuPath      []string `json:"menu_path"`

	// State
	State map[string]interface{} `json:"state"`

	// Timestamps
	StartedAt    time.Time  `json:"started_at"`
	LastActiveAt time.Time  `json:"last_active_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
}

// USSDMenu represents a USSD menu configuration
type USSDMenu struct {
	ID         uuid.UUID    `json:"id"`
	TenantID   uuid.UUID    `json:"tenant_id"`
	Code       string       `json:"code"`
	Title      string       `json:"title"`
	Message    string       `json:"message"`
	Options    []USSDOption `json:"options"`
	IsTerminal bool         `json:"is_terminal"`
	Action     string       `json:"action,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// USSDOption represents a menu option
type USSDOption struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	NextMenuID string `json:"next_menu_id,omitempty"`
	Action     string `json:"action,omitempty"`
}

// WhatsAppTemplate represents a WhatsApp Business template
type WhatsAppTemplate struct {
	ID         uuid.UUID           `json:"id"`
	TenantID   uuid.UUID           `json:"tenant_id"`
	ExternalID string              `json:"external_id"`
	Name       string              `json:"name"`
	Category   string              `json:"category"`
	Language   string              `json:"language"`
	Status     string              `json:"status"`
	Components []TemplateComponent `json:"components"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// TemplateComponent represents a WhatsApp template component
type TemplateComponent struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
	Text   string `json:"text,omitempty"`
}

// Device represents a registered push notification device
type Device struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	UserID       uuid.UUID `json:"user_id"`
	Token        string    `json:"token"`
	Platform     string    `json:"platform"`
	AppVersion   string    `json:"app_version"`
	DeviceModel  string    `json:"device_model"`
	IsActive     bool      `json:"is_active"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	UserID   uuid.UUID                                 `json:"user_id"`
	TenantID uuid.UUID                                 `json:"tenant_id"`
	Channels map[NotificationChannel]ChannelPreference `json:"channels"`
	OptedOut []NotificationChannel                     `json:"opted_out"`
}

// ChannelPreference represents preferences for a specific channel
type ChannelPreference struct {
	Enabled    bool        `json:"enabled"`
	QuietHours *QuietHours `json:"quiet_hours,omitempty"`
	Categories []string    `json:"categories,omitempty"`
}

// QuietHours represents do-not-disturb hours
type QuietHours struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Timezone string `json:"timezone"`
}

// ============================================================================
// Domain Events
// ============================================================================

// NotificationSentEvent is raised when a notification is sent
type NotificationSentEvent struct {
	NotificationID uuid.UUID           `json:"notification_id"`
	Channel        NotificationChannel `json:"channel"`
	RecipientID    uuid.UUID           `json:"recipient_id"`
	Timestamp      time.Time           `json:"timestamp"`
}

// NotificationDeliveredEvent is raised when a notification is delivered
type NotificationDeliveredEvent struct {
	NotificationID uuid.UUID `json:"notification_id"`
	DeliveredAt    time.Time `json:"delivered_at"`
}

// NotificationFailedEvent is raised when a notification fails
type NotificationFailedEvent struct {
	NotificationID uuid.UUID `json:"notification_id"`
	Reason         string    `json:"reason"`
	Timestamp      time.Time `json:"timestamp"`
}
