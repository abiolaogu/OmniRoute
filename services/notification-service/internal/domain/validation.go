// Package domain contains validation methods for notification service domain models.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Validate validates a Notification
func (n *Notification) Validate() error {
	if n.ID == uuid.Nil {
		return ErrRecipientRequired
	}
	if n.TenantID == uuid.Nil {
		return ErrRecipientRequired
	}
	if n.RecipientAddress == "" {
		return ErrRecipientRequired
	}
	if n.Body == "" {
		return ErrBodyRequired
	}
	if !n.Channel.IsValid() {
		return ErrInvalidChannel
	}
	return nil
}

// IsValid checks if the notification channel is valid
func (c NotificationChannel) IsValid() bool {
	switch c {
	case ChannelWhatsApp, ChannelSMS, ChannelEmail, ChannelPush, ChannelVoice, ChannelUSSD:
		return true
	}
	return false
}

// Validate validates a Template
func (t *Template) Validate() error {
	if t.ID == uuid.Nil {
		return ErrTemplateNameRequired
	}
	if t.TenantID == uuid.Nil {
		return ErrTemplateNameRequired
	}
	if t.Name == "" {
		return ErrTemplateNameRequired
	}
	if t.Body == "" {
		return ErrTemplateBodyRequired
	}
	if t.Language == "" {
		return ErrTemplateNameRequired
	}
	return nil
}

// Validate validates a TemplateVariable
func (v TemplateVariable) Validate() error {
	if v.Name == "" {
		return ErrVariableNameRequired
	}
	return nil
}

// IsExpired checks if a USSD session has expired
func (s *USSDSession) IsExpired(timeout time.Duration) bool {
	if s.Status == SessionStatusEnded || s.Status == SessionStatusTimeout {
		return true
	}
	return time.Since(s.LastActiveAt) > timeout
}

// Validate validates a USSDSession
func (s *USSDSession) Validate() error {
	if s.ID == "" {
		return ErrRecipientRequired
	}
	if s.PhoneNumber == "" {
		return ErrRecipientRequired
	}
	if s.ServiceCode == "" {
		return ErrRecipientRequired
	}
	return nil
}

// Validate validates a USSDMenu
func (m *USSDMenu) Validate() error {
	if m.ID == uuid.Nil {
		return ErrTemplateNameRequired
	}
	if m.Code == "" {
		return ErrTemplateNameRequired
	}
	if m.Title == "" {
		return ErrTemplateNameRequired
	}
	if m.Message == "" {
		return ErrBodyRequired
	}
	return nil
}

// IsStale checks if a device is stale (inactive for too long)
func (d *Device) IsStale(timeout time.Duration) bool {
	if !d.IsActive {
		return true
	}
	return time.Since(d.LastActiveAt) > timeout
}

// Validate validates a Device
func (d *Device) Validate() error {
	if d.ID == uuid.Nil {
		return ErrRecipientRequired
	}
	if d.UserID == uuid.Nil {
		return ErrRecipientRequired
	}
	if d.Token == "" {
		return ErrRecipientRequired
	}
	if d.Platform == "" {
		return ErrRecipientRequired
	}
	return nil
}

// IsOptedOut checks if a user has opted out of a channel
func (p *NotificationPreferences) IsOptedOut(channel NotificationChannel) bool {
	for _, c := range p.OptedOut {
		if c == channel {
			return true
		}
	}
	return false
}

// Validate validates a WhatsAppTemplate
func (t *WhatsAppTemplate) Validate() error {
	if t.ID == uuid.Nil {
		return ErrTemplateNameRequired
	}
	if t.Name == "" {
		return ErrTemplateNameRequired
	}
	if t.ExternalID == "" {
		return ErrTemplateNameRequired
	}
	return nil
}
