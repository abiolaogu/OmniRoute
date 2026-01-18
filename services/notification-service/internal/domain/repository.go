// Package domain contains repository interfaces for the Notification Service.
// Following DDD principles, repository interfaces are defined in the domain layer.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationRepository defines operations for notification persistence
type NotificationRepository interface {
	// FindByID retrieves a notification by ID
	FindByID(ctx context.Context, notificationID uuid.UUID) (*Notification, error)

	// FindByRecipient retrieves notifications for a recipient
	FindByRecipient(ctx context.Context, tenantID, recipientID uuid.UUID, limit, offset int) ([]*Notification, error)

	// FindByStatus retrieves notifications by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status NotificationStatus, limit int) ([]*Notification, error)

	// FindPending retrieves pending notifications for sending
	FindPending(ctx context.Context, limit int) ([]*Notification, error)

	// FindScheduled retrieves scheduled notifications due for sending
	FindScheduled(ctx context.Context, before time.Time, limit int) ([]*Notification, error)

	// Save persists a notification
	Save(ctx context.Context, notification *Notification) error

	// Update updates a notification
	Update(ctx context.Context, notification *Notification) error

	// UpdateStatus updates notification status
	UpdateStatus(ctx context.Context, notificationID uuid.UUID, status NotificationStatus) error

	// MarkDelivered marks a notification as delivered
	MarkDelivered(ctx context.Context, notificationID uuid.UUID, deliveredAt time.Time) error

	// MarkRead marks a notification as read
	MarkRead(ctx context.Context, notificationID uuid.UUID, readAt time.Time) error

	// IncrementAttempts increments the attempt count and sets last error
	IncrementAttempts(ctx context.Context, notificationID uuid.UUID, lastError string) error
}

// TemplateRepository defines operations for template persistence
type TemplateRepository interface {
	// FindByID retrieves a template by ID
	FindByID(ctx context.Context, templateID uuid.UUID) (*Template, error)

	// FindByCode retrieves a template by code
	FindByCode(ctx context.Context, tenantID uuid.UUID, code string, channel NotificationChannel) (*Template, error)

	// FindByChannel retrieves templates by channel
	FindByChannel(ctx context.Context, tenantID uuid.UUID, channel NotificationChannel) ([]*Template, error)

	// FindActive retrieves active templates
	FindActive(ctx context.Context, tenantID uuid.UUID, templateType TemplateType) ([]*Template, error)

	// Save persists a template
	Save(ctx context.Context, template *Template) error

	// Update updates a template
	Update(ctx context.Context, template *Template) error

	// Delete removes a template
	Delete(ctx context.Context, templateID uuid.UUID) error

	// Approve marks a template as approved
	Approve(ctx context.Context, templateID uuid.UUID) error
}

// USSDSessionRepository defines operations for USSD session persistence
type USSDSessionRepository interface {
	// FindByID retrieves a session by ID
	FindByID(ctx context.Context, sessionID string) (*USSDSession, error)

	// FindByPhoneNumber retrieves active session for a phone number
	FindByPhoneNumber(ctx context.Context, tenantID uuid.UUID, phoneNumber string) (*USSDSession, error)

	// Save persists a session
	Save(ctx context.Context, session *USSDSession) error

	// Update updates a session
	Update(ctx context.Context, session *USSDSession) error

	// End ends a session
	End(ctx context.Context, sessionID string) error

	// CleanupExpired removes expired sessions
	CleanupExpired(ctx context.Context, before time.Time) (int, error)
}

// USSDMenuRepository defines operations for USSD menu persistence
type USSDMenuRepository interface {
	// FindByID retrieves a menu by ID
	FindByID(ctx context.Context, menuID uuid.UUID) (*USSDMenu, error)

	// FindByCode retrieves a menu by code
	FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*USSDMenu, error)

	// FindRoot retrieves the root menu for a service code
	FindRoot(ctx context.Context, tenantID uuid.UUID, serviceCode string) (*USSDMenu, error)

	// FindAll retrieves all menus for a tenant
	FindAll(ctx context.Context, tenantID uuid.UUID) ([]*USSDMenu, error)

	// Save persists a menu
	Save(ctx context.Context, menu *USSDMenu) error

	// Update updates a menu
	Update(ctx context.Context, menu *USSDMenu) error

	// Delete removes a menu
	Delete(ctx context.Context, menuID uuid.UUID) error
}

// DeviceRepository defines operations for device persistence
type DeviceRepository interface {
	// FindByID retrieves a device by ID
	FindByID(ctx context.Context, deviceID uuid.UUID) (*Device, error)

	// FindByUserID retrieves devices for a user
	FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*Device, error)

	// FindByToken retrieves a device by token
	FindByToken(ctx context.Context, token string) (*Device, error)

	// FindActive retrieves active devices for a user
	FindActive(ctx context.Context, tenantID, userID uuid.UUID) ([]*Device, error)

	// Save persists a device
	Save(ctx context.Context, device *Device) error

	// Update updates a device
	Update(ctx context.Context, device *Device) error

	// Delete removes a device
	Delete(ctx context.Context, deviceID uuid.UUID) error

	// UpdateLastActive updates the last active timestamp
	UpdateLastActive(ctx context.Context, deviceID uuid.UUID) error

	// DeactivateStale deactivates stale devices
	DeactivateStale(ctx context.Context, before time.Time) (int, error)
}

// PreferencesRepository defines operations for notification preferences
type PreferencesRepository interface {
	// FindByUserID retrieves preferences for a user
	FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*NotificationPreferences, error)

	// Save persists preferences
	Save(ctx context.Context, prefs *NotificationPreferences) error

	// Update updates preferences
	Update(ctx context.Context, prefs *NotificationPreferences) error

	// OptOut records an opt-out for a channel
	OptOut(ctx context.Context, tenantID, userID uuid.UUID, channel NotificationChannel) error

	// OptIn records an opt-in for a channel
	OptIn(ctx context.Context, tenantID, userID uuid.UUID, channel NotificationChannel) error
}

// WhatsAppTemplateRepository defines operations for WhatsApp template persistence
type WhatsAppTemplateRepository interface {
	// FindByID retrieves a template by ID
	FindByID(ctx context.Context, templateID uuid.UUID) (*WhatsAppTemplate, error)

	// FindByExternalID retrieves a template by external ID
	FindByExternalID(ctx context.Context, externalID string) (*WhatsAppTemplate, error)

	// FindApproved retrieves approved templates
	FindApproved(ctx context.Context, tenantID uuid.UUID) ([]*WhatsAppTemplate, error)

	// Save persists a template
	Save(ctx context.Context, template *WhatsAppTemplate) error

	// Update updates a template
	Update(ctx context.Context, template *WhatsAppTemplate) error

	// UpdateStatus updates template status
	UpdateStatus(ctx context.Context, templateID uuid.UUID, status string) error
}
