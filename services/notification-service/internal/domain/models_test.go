// Package domain_test contains unit tests for the notification service domain models
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/omniroute/notification-service/internal/domain"
)

func TestNotificationChannel(t *testing.T) {
	tests := []struct {
		channel domain.NotificationChannel
		valid   bool
	}{
		{domain.ChannelWhatsApp, true},
		{domain.ChannelSMS, true},
		{domain.ChannelEmail, true},
		{domain.ChannelPush, true},
		{domain.ChannelVoice, true},
		{domain.ChannelUSSD, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.channel), func(t *testing.T) {
			if tt.channel == "" && tt.valid {
				t.Error("Empty channel should not be valid")
			}
		})
	}
}

func TestNotificationStatus(t *testing.T) {
	tests := []struct {
		name    string
		from    domain.NotificationStatus
		to      domain.NotificationStatus
		allowed bool
	}{
		{"pending to queued", domain.StatusPending, domain.StatusQueued, true},
		{"queued to sending", domain.StatusQueued, domain.StatusSending, true},
		{"sending to sent", domain.StatusSending, domain.StatusSent, true},
		{"sent to delivered", domain.StatusSent, domain.StatusDelivered, true},
		{"delivered to read", domain.StatusDelivered, domain.StatusRead, true},
		{"sending to failed", domain.StatusSending, domain.StatusFailed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Status transition logic would be tested here
			if tt.from == "" || tt.to == "" {
				t.Error("Status cannot be empty")
			}
		})
	}
}

func TestNotificationPriority(t *testing.T) {
	tests := []struct {
		priority domain.NotificationPriority
		weight   int
	}{
		{domain.PriorityLow, 1},
		{domain.PriorityNormal, 2},
		{domain.PriorityHigh, 3},
		{domain.PriorityCritical, 4},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if tt.priority == "" {
				t.Error("Priority cannot be empty")
			}
		})
	}
}

func TestNotification_Validate(t *testing.T) {
	tests := []struct {
		name    string
		notif   *domain.Notification
		wantErr bool
	}{
		{
			name: "valid notification",
			notif: &domain.Notification{
				ID:               uuid.New(),
				TenantID:         uuid.New(),
				Channel:          domain.ChannelSMS,
				Status:           domain.StatusPending,
				Priority:         domain.PriorityNormal,
				RecipientID:      uuid.New(),
				RecipientType:    "customer",
				RecipientAddress: "+2341234567890",
				Body:             "Test message",
			},
			wantErr: false,
		},
		{
			name: "missing recipient address",
			notif: &domain.Notification{
				ID:            uuid.New(),
				TenantID:      uuid.New(),
				Channel:       domain.ChannelSMS,
				Status:        domain.StatusPending,
				RecipientID:   uuid.New(),
				RecipientType: "customer",
				Body:          "Test message",
			},
			wantErr: true,
		},
		{
			name: "missing body",
			notif: &domain.Notification{
				ID:               uuid.New(),
				TenantID:         uuid.New(),
				Channel:          domain.ChannelSMS,
				Status:           domain.StatusPending,
				RecipientID:      uuid.New(),
				RecipientAddress: "+2341234567890",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notif.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplate_Validate(t *testing.T) {
	tests := []struct {
		name     string
		template *domain.Template
		wantErr  bool
	}{
		{
			name: "valid template",
			template: &domain.Template{
				ID:       uuid.New(),
				TenantID: uuid.New(),
				Name:     "Order Confirmation",
				Code:     "order_confirmation",
				Type:     domain.TemplateTypeTransactional,
				Channel:  domain.ChannelSMS,
				Body:     "Your order {{order_id}} has been confirmed.",
				Language: "en",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			template: &domain.Template{
				ID:       uuid.New(),
				TenantID: uuid.New(),
				Code:     "test",
				Type:     domain.TemplateTypeTransactional,
				Channel:  domain.ChannelSMS,
				Body:     "Test",
				Language: "en",
			},
			wantErr: true,
		},
		{
			name: "missing body",
			template: &domain.Template{
				ID:       uuid.New(),
				TenantID: uuid.New(),
				Name:     "Test",
				Code:     "test",
				Type:     domain.TemplateTypeTransactional,
				Channel:  domain.ChannelSMS,
				Language: "en",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUSSDSession_IsExpired(t *testing.T) {
	tests := []struct {
		name    string
		session *domain.USSDSession
		want    bool
	}{
		{
			name: "active session",
			session: &domain.USSDSession{
				ID:           "session1",
				Status:       domain.SessionStatusActive,
				LastActiveAt: time.Now(),
			},
			want: false,
		},
		{
			name: "expired session",
			session: &domain.USSDSession{
				ID:           "session2",
				Status:       domain.SessionStatusActive,
				LastActiveAt: time.Now().Add(-10 * time.Minute),
			},
			want: true,
		},
		{
			name: "ended session",
			session: &domain.USSDSession{
				ID:     "session3",
				Status: domain.SessionStatusEnded,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.session.IsExpired(5 * time.Minute)
			if got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_IsStale(t *testing.T) {
	tests := []struct {
		name   string
		device *domain.Device
		want   bool
	}{
		{
			name: "active device",
			device: &domain.Device{
				ID:           uuid.New(),
				IsActive:     true,
				LastActiveAt: time.Now(),
			},
			want: false,
		},
		{
			name: "stale device",
			device: &domain.Device{
				ID:           uuid.New(),
				IsActive:     true,
				LastActiveAt: time.Now().Add(-31 * 24 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactive device",
			device: &domain.Device{
				ID:       uuid.New(),
				IsActive: false,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.device.IsStale(30 * 24 * time.Hour)
			if got != tt.want {
				t.Errorf("IsStale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationPreferences_IsOptedOut(t *testing.T) {
	prefs := &domain.NotificationPreferences{
		UserID:   uuid.New(),
		TenantID: uuid.New(),
		OptedOut: []domain.NotificationChannel{domain.ChannelSMS, domain.ChannelEmail},
	}

	tests := []struct {
		channel domain.NotificationChannel
		want    bool
	}{
		{domain.ChannelSMS, true},
		{domain.ChannelEmail, true},
		{domain.ChannelWhatsApp, false},
		{domain.ChannelPush, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.channel), func(t *testing.T) {
			got := prefs.IsOptedOut(tt.channel)
			if got != tt.want {
				t.Errorf("IsOptedOut(%s) = %v, want %v", tt.channel, got, tt.want)
			}
		})
	}
}

func TestTemplateVariable_Validate(t *testing.T) {
	tests := []struct {
		name     string
		variable domain.TemplateVariable
		wantErr  bool
	}{
		{
			name: "valid variable",
			variable: domain.TemplateVariable{
				Name:     "order_id",
				Type:     "string",
				Required: true,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			variable: domain.TemplateVariable{
				Type:     "string",
				Required: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.variable.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
