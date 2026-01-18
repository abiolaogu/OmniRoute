// Package api provides HTTP handlers for the Notification Service.
// Supports multi-channel messaging: WhatsApp, SMS, USSD, Push, Email, Voice
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"omniroute/services/notification-service/internal/domain"
)

var tracer = otel.Tracer("notification-service-api")

// ============================================================================
// Handler Dependencies
// ============================================================================

// NotificationService defines business logic operations
type NotificationService interface {
	// Core messaging
	Send(ctx context.Context, cmd SendNotificationCommand) (*domain.Notification, error)
	SendBulk(ctx context.Context, cmd SendBulkNotificationCommand) (*BulkSendResult, error)
	GetNotification(ctx context.Context, tenantID, id string) (*domain.Notification, error)
	ListNotifications(ctx context.Context, query NotificationListQuery) (*NotificationList, error)
	GetDeliveryStatus(ctx context.Context, tenantID, id string) (*domain.DeliveryStatus, error)
	CancelNotification(ctx context.Context, tenantID, id string) error

	// Templates
	CreateTemplate(ctx context.Context, cmd CreateTemplateCommand) (*domain.Template, error)
	UpdateTemplate(ctx context.Context, cmd UpdateTemplateCommand) (*domain.Template, error)
	GetTemplate(ctx context.Context, tenantID, id string) (*domain.Template, error)
	ListTemplates(ctx context.Context, query TemplateListQuery) (*TemplateList, error)
	DeleteTemplate(ctx context.Context, tenantID, id string) error
	RenderTemplate(ctx context.Context, tenantID, templateID string, variables map[string]interface{}) (string, error)

	// USSD
	HandleUSSDRequest(ctx context.Context, req USSDRequest) (*USSDResponse, error)
	GetUSSDSession(ctx context.Context, sessionID string) (*domain.USSDSession, error)
	CreateUSSDMenu(ctx context.Context, cmd CreateUSSDMenuCommand) (*domain.USSDMenu, error)
	UpdateUSSDMenu(ctx context.Context, cmd UpdateUSSDMenuCommand) (*domain.USSDMenu, error)
	GetUSSDMenu(ctx context.Context, tenantID, id string) (*domain.USSDMenu, error)
	ListUSSDMenus(ctx context.Context, tenantID string) ([]*domain.USSDMenu, error)

	// WhatsApp
	HandleWhatsAppWebhook(ctx context.Context, payload WhatsAppWebhookPayload) error
	SendWhatsAppTemplate(ctx context.Context, cmd SendWhatsAppTemplateCommand) (*domain.Notification, error)
	SendWhatsAppInteractive(ctx context.Context, cmd SendWhatsAppInteractiveCommand) (*domain.Notification, error)
	GetWhatsAppTemplates(ctx context.Context, tenantID string) ([]*domain.WhatsAppTemplate, error)
	SyncWhatsAppTemplates(ctx context.Context, tenantID string) error

	// Push notifications
	RegisterDevice(ctx context.Context, cmd RegisterDeviceCommand) (*domain.Device, error)
	UnregisterDevice(ctx context.Context, tenantID, deviceID string) error
	GetUserDevices(ctx context.Context, tenantID, userID string) ([]*domain.Device, error)
	SendPush(ctx context.Context, cmd SendPushCommand) (*domain.Notification, error)
	SendPushToTopic(ctx context.Context, cmd SendPushToTopicCommand) (*domain.Notification, error)
	SubscribeToTopic(ctx context.Context, tenantID, deviceID, topic string) error
	UnsubscribeFromTopic(ctx context.Context, tenantID, deviceID, topic string) error

	// Analytics
	GetChannelStats(ctx context.Context, tenantID string, startDate, endDate time.Time) (*ChannelStats, error)
	GetDeliveryReport(ctx context.Context, query DeliveryReportQuery) (*DeliveryReport, error)
	GetTemplatePerformance(ctx context.Context, tenantID, templateID string, days int) (*TemplatePerformance, error)

	// Preferences
	GetUserPreferences(ctx context.Context, tenantID, userID string) (*domain.NotificationPreferences, error)
	UpdateUserPreferences(ctx context.Context, cmd UpdatePreferencesCommand) (*domain.NotificationPreferences, error)
	OptOut(ctx context.Context, tenantID, userID, channel string) error
	OptIn(ctx context.Context, tenantID, userID, channel string) error
}

// Handler contains HTTP handlers for notifications
type Handler struct {
	service NotificationService
	logger  *zap.Logger
}

// NewHandler creates a new notification handler
func NewHandler(service NotificationService, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ============================================================================
// Request/Response Types
// ============================================================================

// APIResponse is standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIMeta contains response metadata
type APIMeta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// SendNotificationRequest for sending a notification
type SendNotificationRequest struct {
	Channel     string                 `json:"channel" binding:"required,oneof=sms whatsapp email push voice"`
	Recipient   RecipientRequest       `json:"recipient" binding:"required"`
	Content     ContentRequest         `json:"content"`
	TemplateID  string                 `json:"template_id,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Priority    string                 `json:"priority,omitempty" binding:"omitempty,oneof=low normal high urgent"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Reference   string                 `json:"reference,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
}

// RecipientRequest for notification recipient
type RecipientRequest struct {
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
	Name     string `json:"name,omitempty"`
}

// ContentRequest for notification content
type ContentRequest struct {
	// SMS/Voice
	Body string `json:"body,omitempty"`

	// Email
	Subject  string   `json:"subject,omitempty"`
	HTMLBody string   `json:"html_body,omitempty"`
	From     string   `json:"from,omitempty"`
	ReplyTo  string   `json:"reply_to,omitempty"`
	CC       []string `json:"cc,omitempty"`
	BCC      []string `json:"bcc,omitempty"`

	// Push
	Title    string            `json:"title,omitempty"`
	ImageURL string            `json:"image_url,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
	Actions  []PushAction      `json:"actions,omitempty"`
	Sound    string            `json:"sound,omitempty"`
	Badge    int               `json:"badge,omitempty"`

	// WhatsApp
	MediaURL  string `json:"media_url,omitempty"`
	MediaType string `json:"media_type,omitempty"` // image, video, document, audio
	Caption   string `json:"caption,omitempty"`
	FileName  string `json:"file_name,omitempty"`

	// Voice
	VoiceGender   string `json:"voice_gender,omitempty"`
	VoiceLanguage string `json:"voice_language,omitempty"`
	DTMF          bool   `json:"dtmf,omitempty"`
	MaxDigits     int    `json:"max_digits,omitempty"`
	CallbackURL   string `json:"callback_url,omitempty"`
}

// PushAction for push notification actions
type PushAction struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
}

// SendBulkNotificationRequest for bulk sending
type SendBulkNotificationRequest struct {
	Channel    string                 `json:"channel" binding:"required"`
	Recipients []RecipientRequest     `json:"recipients" binding:"required,min=1,max=1000"`
	Content    ContentRequest         `json:"content"`
	TemplateID string                 `json:"template_id,omitempty"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
	Priority   string                 `json:"priority,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	BatchSize  int                    `json:"batch_size,omitempty"`
}

// CreateTemplateRequest for creating a template
type CreateTemplateRequest struct {
	Name        string            `json:"name" binding:"required"`
	Channel     string            `json:"channel" binding:"required,oneof=sms whatsapp email push voice"`
	Content     TemplateContent   `json:"content" binding:"required"`
	Description string            `json:"description,omitempty"`
	Category    string            `json:"category,omitempty"`
	Language    string            `json:"language,omitempty"`
	Variables   []TemplateVar     `json:"variables,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TemplateContent for template content by channel
type TemplateContent struct {
	Body      string       `json:"body,omitempty"`
	Subject   string       `json:"subject,omitempty"`
	HTMLBody  string       `json:"html_body,omitempty"`
	Title     string       `json:"title,omitempty"`
	Header    *MediaHeader `json:"header,omitempty"`
	Footer    string       `json:"footer,omitempty"`
	Buttons   []Button     `json:"buttons,omitempty"`
}

// MediaHeader for WhatsApp template headers
type MediaHeader struct {
	Type     string `json:"type"` // text, image, video, document
	Text     string `json:"text,omitempty"`
	MediaURL string `json:"media_url,omitempty"`
}

// Button for WhatsApp/Push buttons
type Button struct {
	Type    string `json:"type"` // quick_reply, url, call, copy_code
	Text    string `json:"text"`
	URL     string `json:"url,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// TemplateVar describes a template variable
type TemplateVar struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, number, date, currency
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateTemplateRequest for updating a template
type UpdateTemplateRequest struct {
	Name        string            `json:"name,omitempty"`
	Content     *TemplateContent  `json:"content,omitempty"`
	Description string            `json:"description,omitempty"`
	Category    string            `json:"category,omitempty"`
	Language    string            `json:"language,omitempty"`
	Variables   []TemplateVar     `json:"variables,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Active      *bool             `json:"active,omitempty"`
}

// CreateUSSDMenuRequest for creating USSD menu
type CreateUSSDMenuRequest struct {
	Code        string         `json:"code" binding:"required"` // e.g., *123#
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description,omitempty"`
	RootScreen  USSDScreen     `json:"root_screen" binding:"required"`
	Screens     []USSDScreen   `json:"screens,omitempty"`
	Timeout     int            `json:"timeout,omitempty"` // seconds
	Active      bool           `json:"active"`
}

// USSDScreen represents a USSD menu screen
type USSDScreen struct {
	ID          string       `json:"id" binding:"required"`
	Type        string       `json:"type" binding:"required,oneof=menu input confirm end"`
	Title       string       `json:"title,omitempty"`
	Body        string       `json:"body" binding:"required"`
	Options     []USSDOption `json:"options,omitempty"`
	InputType   string       `json:"input_type,omitempty"` // text, number, phone, pin
	MinLength   int          `json:"min_length,omitempty"`
	MaxLength   int          `json:"max_length,omitempty"`
	Validation  string       `json:"validation,omitempty"` // regex
	NextScreen  string       `json:"next_screen,omitempty"`
	Action      string       `json:"action,omitempty"` // API action to call
	ActionData  interface{}  `json:"action_data,omitempty"`
}

// USSDOption for menu options
type USSDOption struct {
	Key        string      `json:"key" binding:"required"` // 1, 2, 3...
	Label      string      `json:"label" binding:"required"`
	NextScreen string      `json:"next_screen,omitempty"`
	Action     string      `json:"action,omitempty"`
	ActionData interface{} `json:"action_data,omitempty"`
}

// USSDWebhookRequest from USSD gateway
type USSDWebhookRequest struct {
	SessionID   string `json:"session_id" binding:"required"`
	ServiceCode string `json:"service_code" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	NetworkCode string `json:"network_code,omitempty"`
	Input       string `json:"input"`
	Type        string `json:"type"` // initiation, response, timeout, end
}

// USSDWebhookResponse to USSD gateway
type USSDWebhookResponse struct {
	SessionID string `json:"session_id"`
	Type      string `json:"type"` // CON (continue), END
	Message   string `json:"message"`
}

// WhatsAppWebhookRequest from WhatsApp Business API
type WhatsAppWebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      *struct {
						Body string `json:"body"`
					} `json:"text,omitempty"`
					Image *struct {
						Caption  string `json:"caption"`
						MimeType string `json:"mime_type"`
						SHA256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"image,omitempty"`
					Document *struct {
						Caption  string `json:"caption"`
						FileName string `json:"filename"`
						MimeType string `json:"mime_type"`
						SHA256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"document,omitempty"`
					Location *struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
						Name      string  `json:"name"`
						Address   string  `json:"address"`
					} `json:"location,omitempty"`
					Interactive *struct {
						Type        string `json:"type"`
						ButtonReply *struct {
							ID    string `json:"id"`
							Title string `json:"title"`
						} `json:"button_reply,omitempty"`
						ListReply *struct {
							ID          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description"`
						} `json:"list_reply,omitempty"`
					} `json:"interactive,omitempty"`
					Context *struct {
						From string `json:"from"`
						ID   string `json:"id"`
					} `json:"context,omitempty"`
				} `json:"messages"`
				Statuses []struct {
					ID          string `json:"id"`
					Status      string `json:"status"` // sent, delivered, read, failed
					Timestamp   string `json:"timestamp"`
					RecipientID string `json:"recipient_id"`
					Errors      []struct {
						Code    int    `json:"code"`
						Title   string `json:"title"`
						Message string `json:"message"`
					} `json:"errors,omitempty"`
				} `json:"statuses"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// SendWhatsAppTemplateRequest for WhatsApp template messages
type SendWhatsAppTemplateRequest struct {
	Phone        string              `json:"phone" binding:"required"`
	TemplateName string              `json:"template_name" binding:"required"`
	Language     string              `json:"language" binding:"required"`
	Components   []TemplateComponent `json:"components,omitempty"`
	Reference    string              `json:"reference,omitempty"`
}

// TemplateComponent for WhatsApp template parameters
type TemplateComponent struct {
	Type       string               `json:"type"` // header, body, button
	Parameters []TemplateParameter  `json:"parameters,omitempty"`
	SubType    string               `json:"sub_type,omitempty"` // for buttons: quick_reply, url
	Index      int                  `json:"index,omitempty"`    // button index
}

// TemplateParameter for template variable values
type TemplateParameter struct {
	Type     string         `json:"type"` // text, currency, date_time, image, document, video
	Text     string         `json:"text,omitempty"`
	Currency *CurrencyValue `json:"currency,omitempty"`
	DateTime *DateTimeValue `json:"date_time,omitempty"`
	Image    *MediaValue    `json:"image,omitempty"`
	Document *MediaValue    `json:"document,omitempty"`
	Video    *MediaValue    `json:"video,omitempty"`
}

// CurrencyValue for currency parameters
type CurrencyValue struct {
	FallbackValue string `json:"fallback_value"`
	Code          string `json:"code"`
	Amount1000    int64  `json:"amount_1000"` // Amount * 1000
}

// DateTimeValue for date/time parameters
type DateTimeValue struct {
	FallbackValue string `json:"fallback_value"`
}

// MediaValue for media parameters
type MediaValue struct {
	Link     string `json:"link,omitempty"`
	ID       string `json:"id,omitempty"`
	Caption  string `json:"caption,omitempty"`
	FileName string `json:"filename,omitempty"`
}

// SendWhatsAppInteractiveRequest for interactive messages
type SendWhatsAppInteractiveRequest struct {
	Phone           string              `json:"phone" binding:"required"`
	InteractiveType string              `json:"type" binding:"required,oneof=button list product product_list"`
	Header          *InteractiveHeader  `json:"header,omitempty"`
	Body            string              `json:"body" binding:"required"`
	Footer          string              `json:"footer,omitempty"`
	Action          InteractiveAction   `json:"action" binding:"required"`
	Reference       string              `json:"reference,omitempty"`
}

// InteractiveHeader for WhatsApp interactive header
type InteractiveHeader struct {
	Type     string `json:"type"` // text, image, video, document
	Text     string `json:"text,omitempty"`
	MediaURL string `json:"media_url,omitempty"`
}

// InteractiveAction for WhatsApp interactive actions
type InteractiveAction struct {
	// For buttons
	Buttons []InteractiveButton `json:"buttons,omitempty"`

	// For lists
	Button   string            `json:"button,omitempty"` // List button text
	Sections []ListSection     `json:"sections,omitempty"`

	// For products
	CatalogID         string   `json:"catalog_id,omitempty"`
	ProductRetailerID string   `json:"product_retailer_id,omitempty"`
	Sections          []ProductSection `json:"product_sections,omitempty"`
}

// InteractiveButton for button actions
type InteractiveButton struct {
	Type  string `json:"type"` // reply
	Reply struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"reply"`
}

// ListSection for list messages
type ListSection struct {
	Title string     `json:"title,omitempty"`
	Rows  []ListRow  `json:"rows"`
}

// ListRow for list items
type ListRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// ProductSection for product list messages
type ProductSection struct {
	Title              string   `json:"title,omitempty"`
	ProductRetailerIDs []string `json:"product_retailer_ids"`
}

// RegisterDeviceRequest for push device registration
type RegisterDeviceRequest struct {
	UserID       string            `json:"user_id" binding:"required"`
	Token        string            `json:"token" binding:"required"`
	Platform     string            `json:"platform" binding:"required,oneof=ios android web"`
	AppVersion   string            `json:"app_version,omitempty"`
	DeviceModel  string            `json:"device_model,omitempty"`
	OSVersion    string            `json:"os_version,omitempty"`
	Locale       string            `json:"locale,omitempty"`
	Timezone     string            `json:"timezone,omitempty"`
	Topics       []string          `json:"topics,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// SendPushRequest for direct push notification
type SendPushRequest struct {
	UserID   string            `json:"user_id,omitempty"`
	DeviceID string            `json:"device_id,omitempty"`
	Title    string            `json:"title" binding:"required"`
	Body     string            `json:"body" binding:"required"`
	ImageURL string            `json:"image_url,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
	Actions  []PushAction      `json:"actions,omitempty"`
	Sound    string            `json:"sound,omitempty"`
	Badge    int               `json:"badge,omitempty"`
	Priority string            `json:"priority,omitempty"`
	TTL      int               `json:"ttl,omitempty"` // seconds
	CollapseKey string         `json:"collapse_key,omitempty"`
}

// SendPushToTopicRequest for topic-based push
type SendPushToTopicRequest struct {
	Topic    string            `json:"topic" binding:"required"`
	Title    string            `json:"title" binding:"required"`
	Body     string            `json:"body" binding:"required"`
	ImageURL string            `json:"image_url,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
	Actions  []PushAction      `json:"actions,omitempty"`
}

// UpdatePreferencesRequest for notification preferences
type UpdatePreferencesRequest struct {
	UserID          string                 `json:"user_id" binding:"required"`
	Channels        map[string]bool        `json:"channels,omitempty"`        // channel -> enabled
	Categories      map[string]bool        `json:"categories,omitempty"`      // category -> enabled
	QuietHours      *QuietHours            `json:"quiet_hours,omitempty"`
	Frequency       map[string]string      `json:"frequency,omitempty"`       // category -> realtime/daily/weekly
	PreferredLang   string                 `json:"preferred_language,omitempty"`
}

// QuietHours for do-not-disturb settings
type QuietHours struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"start_time"` // HH:MM
	EndTime   string `json:"end_time"`   // HH:MM
	Timezone  string `json:"timezone"`
	Days      []int  `json:"days,omitempty"` // 0=Sunday
}

// SMSWebhookRequest from SMS gateway (delivery reports)
type SMSWebhookRequest struct {
	MessageID   string `json:"message_id"`
	Recipient   string `json:"recipient"`
	Status      string `json:"status"` // delivered, failed, expired, rejected
	ErrorCode   string `json:"error_code,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
	Timestamp   string `json:"timestamp"`
	NetworkCode string `json:"network_code,omitempty"`
}

// InboundSMSRequest for received SMS
type InboundSMSRequest struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Body        string `json:"body"`
	MessageID   string `json:"message_id"`
	Timestamp   string `json:"timestamp"`
	NetworkCode string `json:"network_code,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
}

// VoiceCallbackRequest for voice call events
type VoiceCallbackRequest struct {
	CallID      string `json:"call_id"`
	SessionID   string `json:"session_id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Status      string `json:"status"` // ringing, answered, completed, failed, busy, no_answer
	Direction   string `json:"direction"` // inbound, outbound
	Duration    int    `json:"duration,omitempty"` // seconds
	DTMFDigits  string `json:"dtmf_digits,omitempty"`
	RecordingURL string `json:"recording_url,omitempty"`
	Timestamp   string `json:"timestamp"`
}

// ============================================================================
// Command Types
// ============================================================================

type SendNotificationCommand struct {
	TenantID    string
	Channel     string
	Recipient   domain.Recipient
	Content     domain.Content
	TemplateID  string
	Variables   map[string]interface{}
	Priority    string
	ScheduledAt *time.Time
	ExpiresAt   *time.Time
	Reference   string
	Tags        []string
	Metadata    map[string]string
	SenderID    string
}

type SendBulkNotificationCommand struct {
	TenantID   string
	Channel    string
	Recipients []domain.Recipient
	Content    domain.Content
	TemplateID string
	Variables  map[string]interface{}
	Priority   string
	Tags       []string
	BatchSize  int
	SenderID   string
}

type CreateTemplateCommand struct {
	TenantID    string
	Name        string
	Channel     string
	Content     domain.TemplateContent
	Description string
	Category    string
	Language    string
	Variables   []domain.TemplateVariable
	Tags        []string
	Metadata    map[string]string
	CreatedBy   string
}

type UpdateTemplateCommand struct {
	TenantID    string
	TemplateID  string
	Name        string
	Content     *domain.TemplateContent
	Description string
	Category    string
	Language    string
	Variables   []domain.TemplateVariable
	Tags        []string
	Metadata    map[string]string
	Active      *bool
	UpdatedBy   string
}

type CreateUSSDMenuCommand struct {
	TenantID    string
	Code        string
	Name        string
	Description string
	RootScreen  domain.USSDScreen
	Screens     []domain.USSDScreen
	Timeout     int
	Active      bool
	CreatedBy   string
}

type UpdateUSSDMenuCommand struct {
	TenantID    string
	MenuID      string
	Name        string
	Description string
	RootScreen  *domain.USSDScreen
	Screens     []domain.USSDScreen
	Timeout     int
	Active      *bool
	UpdatedBy   string
}

type USSDRequest struct {
	SessionID   string
	ServiceCode string
	PhoneNumber string
	NetworkCode string
	Input       string
	Type        string
	TenantID    string
}

type WhatsAppWebhookPayload struct {
	TenantID string
	Raw      json.RawMessage
}

type SendWhatsAppTemplateCommand struct {
	TenantID     string
	Phone        string
	TemplateName string
	Language     string
	Components   []domain.TemplateComponent
	Reference    string
	SenderID     string
}

type SendWhatsAppInteractiveCommand struct {
	TenantID        string
	Phone           string
	InteractiveType string
	Header          *domain.InteractiveHeader
	Body            string
	Footer          string
	Action          domain.InteractiveAction
	Reference       string
	SenderID        string
}

type RegisterDeviceCommand struct {
	TenantID    string
	UserID      string
	Token       string
	Platform    string
	AppVersion  string
	DeviceModel string
	OSVersion   string
	Locale      string
	Timezone    string
	Topics      []string
	Metadata    map[string]string
}

type SendPushCommand struct {
	TenantID    string
	UserID      string
	DeviceID    string
	Title       string
	Body        string
	ImageURL    string
	Data        map[string]string
	Actions     []domain.PushAction
	Sound       string
	Badge       int
	Priority    string
	TTL         int
	CollapseKey string
	SenderID    string
}

type SendPushToTopicCommand struct {
	TenantID string
	Topic    string
	Title    string
	Body     string
	ImageURL string
	Data     map[string]string
	Actions  []domain.PushAction
	SenderID string
}

type UpdatePreferencesCommand struct {
	TenantID      string
	UserID        string
	Channels      map[string]bool
	Categories    map[string]bool
	QuietHours    *domain.QuietHours
	Frequency     map[string]string
	PreferredLang string
}

// ============================================================================
// Query/Result Types
// ============================================================================

type NotificationListQuery struct {
	TenantID  string
	Channel   string
	Status    string
	Recipient string
	Reference string
	Tags      []string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type NotificationList struct {
	Notifications []*domain.Notification
	TotalCount    int64
	Page          int
	PageSize      int
}

type TemplateListQuery struct {
	TenantID string
	Channel  string
	Category string
	Search   string
	Active   *bool
	Tags     []string
	Page     int
	PageSize int
}

type TemplateList struct {
	Templates  []*domain.Template
	TotalCount int64
	Page       int
	PageSize   int
}

type BulkSendResult struct {
	TotalRequested int
	TotalQueued    int
	TotalFailed    int
	BatchID        string
	Errors         []BulkSendError
}

type BulkSendError struct {
	Index   int
	Phone   string
	Error   string
}

type ChannelStats struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
	Channels  map[string]ChannelMetrics
}

type ChannelMetrics struct {
	Sent      int64
	Delivered int64
	Failed    int64
	Read      int64
	Clicked   int64
	Cost      float64
	Currency  string
}

type DeliveryReportQuery struct {
	TenantID  string
	Channel   string
	StartDate time.Time
	EndDate   time.Time
	GroupBy   string // hour, day, week
}

type DeliveryReport struct {
	Query   DeliveryReportQuery
	Data    []DeliveryReportData
	Summary DeliveryReportSummary
}

type DeliveryReportData struct {
	Period    time.Time
	Sent      int64
	Delivered int64
	Failed    int64
	Pending   int64
	Rate      float64
}

type DeliveryReportSummary struct {
	TotalSent      int64
	TotalDelivered int64
	TotalFailed    int64
	AvgDeliveryRate float64
	TotalCost      float64
}

type TemplatePerformance struct {
	TemplateID     string
	TemplateName   string
	TotalSent      int64
	TotalDelivered int64
	TotalFailed    int64
	DeliveryRate   float64
	OpenRate       float64
	ClickRate      float64
	DailyMetrics   []DailyMetric
}

type DailyMetric struct {
	Date      time.Time
	Sent      int64
	Delivered int64
	Opened    int64
	Clicked   int64
}

// ============================================================================
// Route Registration
// ============================================================================

// RegisterRoutes registers all notification routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Health checks
	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)

	api := r.Group("/api/v1")
	{
		// Core notifications
		notifications := api.Group("/notifications")
		{
			notifications.POST("", h.SendNotification)
			notifications.POST("/bulk", h.SendBulkNotification)
			notifications.GET("", h.ListNotifications)
			notifications.GET("/:id", h.GetNotification)
			notifications.GET("/:id/status", h.GetDeliveryStatus)
			notifications.DELETE("/:id", h.CancelNotification)
		}

		// SMS
		sms := api.Group("/sms")
		{
			sms.POST("/send", h.SendSMS)
			sms.POST("/webhook/delivery", h.SMSDeliveryWebhook)
			sms.POST("/webhook/inbound", h.InboundSMSWebhook)
		}

		// WhatsApp
		whatsapp := api.Group("/whatsapp")
		{
			whatsapp.POST("/send", h.SendWhatsApp)
			whatsapp.POST("/template", h.SendWhatsAppTemplate)
			whatsapp.POST("/interactive", h.SendWhatsAppInteractive)
			whatsapp.GET("/templates", h.GetWhatsAppTemplates)
			whatsapp.POST("/templates/sync", h.SyncWhatsAppTemplates)
			whatsapp.GET("/webhook", h.WhatsAppWebhookVerify)
			whatsapp.POST("/webhook", h.WhatsAppWebhook)
		}

		// USSD
		ussd := api.Group("/ussd")
		{
			ussd.POST("/callback", h.USSDCallback)
			ussd.GET("/sessions/:sessionId", h.GetUSSDSession)
			ussd.POST("/menus", h.CreateUSSDMenu)
			ussd.GET("/menus", h.ListUSSDMenus)
			ussd.GET("/menus/:id", h.GetUSSDMenu)
			ussd.PUT("/menus/:id", h.UpdateUSSDMenu)
			ussd.DELETE("/menus/:id", h.DeleteUSSDMenu)
		}

		// Email
		email := api.Group("/email")
		{
			email.POST("/send", h.SendEmail)
			email.POST("/webhook/events", h.EmailEventsWebhook)
		}

		// Push
		push := api.Group("/push")
		{
			push.POST("/send", h.SendPush)
			push.POST("/topic", h.SendPushToTopic)
			push.POST("/devices", h.RegisterDevice)
			push.DELETE("/devices/:id", h.UnregisterDevice)
			push.GET("/users/:userId/devices", h.GetUserDevices)
			push.POST("/devices/:id/topics/:topic", h.SubscribeToTopic)
			push.DELETE("/devices/:id/topics/:topic", h.UnsubscribeFromTopic)
		}

		// Voice
		voice := api.Group("/voice")
		{
			voice.POST("/call", h.InitiateCall)
			voice.POST("/webhook/events", h.VoiceEventsWebhook)
		}

		// Templates
		templates := api.Group("/templates")
		{
			templates.POST("", h.CreateTemplate)
			templates.GET("", h.ListTemplates)
			templates.GET("/:id", h.GetTemplate)
			templates.PUT("/:id", h.UpdateTemplate)
			templates.DELETE("/:id", h.DeleteTemplate)
			templates.POST("/:id/render", h.RenderTemplate)
			templates.GET("/:id/performance", h.GetTemplatePerformance)
		}

		// Preferences
		preferences := api.Group("/preferences")
		{
			preferences.GET("/:userId", h.GetUserPreferences)
			preferences.PUT("/:userId", h.UpdateUserPreferences)
			preferences.POST("/:userId/opt-out/:channel", h.OptOut)
			preferences.POST("/:userId/opt-in/:channel", h.OptIn)
		}

		// Analytics
		analytics := api.Group("/analytics")
		{
			analytics.GET("/channels", h.GetChannelStats)
			analytics.GET("/delivery", h.GetDeliveryReport)
		}
	}
}

// ============================================================================
// Health Handlers
// ============================================================================

// Health checks service health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"status":  "healthy",
			"service": "notification-service",
		},
	})
}

// Ready checks if service is ready to accept traffic
func (h *Handler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"status": "ready",
		},
	})
}

// ============================================================================
// Core Notification Handlers
// ============================================================================

// SendNotification sends a notification via specified channel
func (h *Handler) SendNotification(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendNotification")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("channel", req.Channel),
		attribute.String("tenant_id", tenantID),
	)

	cmd := SendNotificationCommand{
		TenantID:    tenantID,
		Channel:     req.Channel,
		Recipient:   h.mapRecipient(req.Recipient),
		Content:     h.mapContent(req.Content),
		TemplateID:  req.TemplateID,
		Variables:   req.Variables,
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
		ExpiresAt:   req.ExpiresAt,
		Reference:   req.Reference,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		SenderID:    userID,
	}

	notification, err := h.service.Send(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SendBulkNotification sends notifications to multiple recipients
func (h *Handler) SendBulkNotification(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendBulkNotification")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendBulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("channel", req.Channel),
		attribute.Int("recipient_count", len(req.Recipients)),
	)

	recipients := make([]domain.Recipient, len(req.Recipients))
	for i, r := range req.Recipients {
		recipients[i] = h.mapRecipient(r)
	}

	cmd := SendBulkNotificationCommand{
		TenantID:   tenantID,
		Channel:    req.Channel,
		Recipients: recipients,
		Content:    h.mapContent(req.Content),
		TemplateID: req.TemplateID,
		Variables:  req.Variables,
		Priority:   req.Priority,
		Tags:       req.Tags,
		BatchSize:  req.BatchSize,
		SenderID:   userID,
	}

	result, err := h.service.SendBulk(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetNotification retrieves a notification by ID
func (h *Handler) GetNotification(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetNotification")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	notification, err := h.service.GetNotification(ctx, tenantID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// ListNotifications lists notifications with filters
func (h *Handler) ListNotifications(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "ListNotifications")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	query := NotificationListQuery{
		TenantID:  tenantID,
		Channel:   c.Query("channel"),
		Status:    c.Query("status"),
		Recipient: c.Query("recipient"),
		Reference: c.Query("reference"),
		Page:      h.parseInt(c.Query("page"), 1),
		PageSize:  h.parseInt(c.Query("page_size"), 20),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	if tags := c.Query("tags"); tags != "" {
		query.Tags = strings.Split(tags, ",")
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			query.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			query.EndDate = &t
		}
	}

	result, err := h.service.ListNotifications(ctx, query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	totalPages := int(result.TotalCount) / query.PageSize
	if int(result.TotalCount)%query.PageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    result.Notifications,
		Meta: &APIMeta{
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalCount: result.TotalCount,
			TotalPages: totalPages,
		},
	})
}

// GetDeliveryStatus retrieves delivery status for a notification
func (h *Handler) GetDeliveryStatus(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetDeliveryStatus")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	status, err := h.service.GetDeliveryStatus(ctx, tenantID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    status,
	})
}

// CancelNotification cancels a scheduled notification
func (h *Handler) CancelNotification(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "CancelNotification")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	if err := h.service.CancelNotification(ctx, tenantID, id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"cancelled": true},
	})
}

// ============================================================================
// SMS Handlers
// ============================================================================

// SendSMS sends an SMS message
func (h *Handler) SendSMS(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendSMS")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := SendNotificationCommand{
		TenantID:   tenantID,
		Channel:    "sms",
		Recipient:  h.mapRecipient(req.Recipient),
		Content:    h.mapContent(req.Content),
		TemplateID: req.TemplateID,
		Variables:  req.Variables,
		Priority:   req.Priority,
		Reference:  req.Reference,
		Tags:       req.Tags,
		Metadata:   req.Metadata,
		SenderID:   userID,
	}

	notification, err := h.service.Send(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SMSDeliveryWebhook handles SMS delivery reports
func (h *Handler) SMSDeliveryWebhook(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SMSDeliveryWebhook")
	defer span.End()

	var req SMSWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	h.logger.Info("SMS delivery report received",
		zap.String("message_id", req.MessageID),
		zap.String("status", req.Status),
	)

	// Process delivery status update asynchronously
	// Implementation would update notification status in database

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// InboundSMSWebhook handles inbound SMS messages
func (h *Handler) InboundSMSWebhook(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "InboundSMSWebhook")
	defer span.End()

	var req InboundSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	h.logger.Info("Inbound SMS received",
		zap.String("from", req.From),
		zap.String("body", req.Body),
	)

	// Process inbound message - could trigger automated responses,
	// keyword handling, etc.

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ============================================================================
// WhatsApp Handlers
// ============================================================================

// SendWhatsApp sends a WhatsApp message
func (h *Handler) SendWhatsApp(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendWhatsApp")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := SendNotificationCommand{
		TenantID:   tenantID,
		Channel:    "whatsapp",
		Recipient:  h.mapRecipient(req.Recipient),
		Content:    h.mapContent(req.Content),
		TemplateID: req.TemplateID,
		Variables:  req.Variables,
		Priority:   req.Priority,
		Reference:  req.Reference,
		Tags:       req.Tags,
		Metadata:   req.Metadata,
		SenderID:   userID,
	}

	notification, err := h.service.Send(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SendWhatsAppTemplate sends a WhatsApp template message
func (h *Handler) SendWhatsAppTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendWhatsAppTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendWhatsAppTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	components := make([]domain.TemplateComponent, len(req.Components))
	for i, c := range req.Components {
		components[i] = h.mapTemplateComponent(c)
	}

	cmd := SendWhatsAppTemplateCommand{
		TenantID:     tenantID,
		Phone:        req.Phone,
		TemplateName: req.TemplateName,
		Language:     req.Language,
		Components:   components,
		Reference:    req.Reference,
		SenderID:     userID,
	}

	notification, err := h.service.SendWhatsAppTemplate(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SendWhatsAppInteractive sends a WhatsApp interactive message
func (h *Handler) SendWhatsAppInteractive(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendWhatsAppInteractive")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendWhatsAppInteractiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := SendWhatsAppInteractiveCommand{
		TenantID:        tenantID,
		Phone:           req.Phone,
		InteractiveType: req.InteractiveType,
		Header:          h.mapInteractiveHeader(req.Header),
		Body:            req.Body,
		Footer:          req.Footer,
		Action:          h.mapInteractiveAction(req.Action),
		Reference:       req.Reference,
		SenderID:        userID,
	}

	notification, err := h.service.SendWhatsAppInteractive(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// GetWhatsAppTemplates retrieves WhatsApp Business templates
func (h *Handler) GetWhatsAppTemplates(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetWhatsAppTemplates")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	templates, err := h.service.GetWhatsAppTemplates(ctx, tenantID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    templates,
	})
}

// SyncWhatsAppTemplates syncs templates from WhatsApp Business API
func (h *Handler) SyncWhatsAppTemplates(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SyncWhatsAppTemplates")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	if err := h.service.SyncWhatsAppTemplates(ctx, tenantID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"synced": true},
	})
}

// WhatsAppWebhookVerify handles WhatsApp webhook verification
func (h *Handler) WhatsAppWebhookVerify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	// Verify token should be configured per tenant
	expectedToken := c.GetString("whatsapp_verify_token")
	if expectedToken == "" {
		expectedToken = "omniroute_verify_token" // Default
	}

	if mode == "subscribe" && token == expectedToken {
		c.String(http.StatusOK, challenge)
		return
	}

	c.String(http.StatusForbidden, "Forbidden")
}

// WhatsAppWebhook handles WhatsApp webhook events
func (h *Handler) WhatsAppWebhook(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "WhatsAppWebhook")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	var req WhatsAppWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Store raw payload for processing
	rawPayload, _ := json.Marshal(req)
	payload := WhatsAppWebhookPayload{
		TenantID: tenantID,
		Raw:      rawPayload,
	}

	if err := h.service.HandleWhatsAppWebhook(ctx, payload); err != nil {
		h.logger.Error("Failed to process WhatsApp webhook",
			zap.Error(err),
		)
		// Still return 200 to acknowledge receipt
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ============================================================================
// USSD Handlers
// ============================================================================

// USSDCallback handles USSD gateway callbacks
func (h *Handler) USSDCallback(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "USSDCallback")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	var req USSDWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try form data (common for USSD gateways)
		req = USSDWebhookRequest{
			SessionID:   c.PostForm("sessionId"),
			ServiceCode: c.PostForm("serviceCode"),
			PhoneNumber: c.PostForm("phoneNumber"),
			NetworkCode: c.PostForm("networkCode"),
			Input:       c.PostForm("text"),
			Type:        c.DefaultPostForm("type", "response"),
		}
	}

	span.SetAttributes(
		attribute.String("session_id", req.SessionID),
		attribute.String("service_code", req.ServiceCode),
		attribute.String("phone_number", req.PhoneNumber),
	)

	ussdReq := USSDRequest{
		SessionID:   req.SessionID,
		ServiceCode: req.ServiceCode,
		PhoneNumber: req.PhoneNumber,
		NetworkCode: req.NetworkCode,
		Input:       req.Input,
		Type:        req.Type,
		TenantID:    tenantID,
	}

	response, err := h.service.HandleUSSDRequest(ctx, ussdReq)
	if err != nil {
		h.logger.Error("USSD handling failed",
			zap.Error(err),
			zap.String("session_id", req.SessionID),
		)
		// Return END response on error
		c.String(http.StatusOK, "END An error occurred. Please try again.")
		return
	}

	// Format response based on gateway expectations
	// CON = continue session, END = end session
	prefix := "CON "
	if response.Type == "END" {
		prefix = "END "
	}

	c.String(http.StatusOK, prefix+response.Message)
}

// GetUSSDSession retrieves USSD session details
func (h *Handler) GetUSSDSession(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetUSSDSession")
	defer span.End()

	sessionID := c.Param("sessionId")

	session, err := h.service.GetUSSDSession(ctx, sessionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    session,
	})
}

// CreateUSSDMenu creates a new USSD menu
func (h *Handler) CreateUSSDMenu(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "CreateUSSDMenu")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req CreateUSSDMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := CreateUSSDMenuCommand{
		TenantID:    tenantID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		RootScreen:  h.mapUSSDScreen(req.RootScreen),
		Screens:     h.mapUSSDScreens(req.Screens),
		Timeout:     req.Timeout,
		Active:      req.Active,
		CreatedBy:   userID,
	}

	menu, err := h.service.CreateUSSDMenu(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    menu,
	})
}

// ListUSSDMenus lists USSD menus
func (h *Handler) ListUSSDMenus(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "ListUSSDMenus")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	menus, err := h.service.ListUSSDMenus(ctx, tenantID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    menus,
	})
}

// GetUSSDMenu retrieves a USSD menu
func (h *Handler) GetUSSDMenu(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetUSSDMenu")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	menu, err := h.service.GetUSSDMenu(ctx, tenantID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    menu,
	})
}

// UpdateUSSDMenu updates a USSD menu
func (h *Handler) UpdateUSSDMenu(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "UpdateUSSDMenu")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req CreateUSSDMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	rootScreen := h.mapUSSDScreen(req.RootScreen)
	cmd := UpdateUSSDMenuCommand{
		TenantID:    tenantID,
		MenuID:      id,
		Name:        req.Name,
		Description: req.Description,
		RootScreen:  &rootScreen,
		Screens:     h.mapUSSDScreens(req.Screens),
		Timeout:     req.Timeout,
		Active:      &req.Active,
		UpdatedBy:   userID,
	}

	menu, err := h.service.UpdateUSSDMenu(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    menu,
	})
}

// DeleteUSSDMenu deletes a USSD menu
func (h *Handler) DeleteUSSDMenu(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "DeleteUSSDMenu")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	// Soft delete by deactivating
	active := false
	cmd := UpdateUSSDMenuCommand{
		TenantID: tenantID,
		MenuID:   id,
		Active:   &active,
	}

	_, err := h.service.UpdateUSSDMenu(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"deleted": true},
	})
}

// ============================================================================
// Email Handlers
// ============================================================================

// SendEmail sends an email
func (h *Handler) SendEmail(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendEmail")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := SendNotificationCommand{
		TenantID:    tenantID,
		Channel:     "email",
		Recipient:   h.mapRecipient(req.Recipient),
		Content:     h.mapContent(req.Content),
		TemplateID:  req.TemplateID,
		Variables:   req.Variables,
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
		Reference:   req.Reference,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		SenderID:    userID,
	}

	notification, err := h.service.Send(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// EmailEventsWebhook handles email provider webhooks (bounce, complaint, etc.)
func (h *Handler) EmailEventsWebhook(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "EmailEventsWebhook")
	defer span.End()

	// Handle various email events: delivered, bounced, complained, opened, clicked
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	h.logger.Info("Email event received",
		zap.Any("payload", payload),
	)

	// Process asynchronously
	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ============================================================================
// Push Notification Handlers
// ============================================================================

// RegisterDevice registers a device for push notifications
func (h *Handler) RegisterDevice(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "RegisterDevice")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := RegisterDeviceCommand{
		TenantID:    tenantID,
		UserID:      req.UserID,
		Token:       req.Token,
		Platform:    req.Platform,
		AppVersion:  req.AppVersion,
		DeviceModel: req.DeviceModel,
		OSVersion:   req.OSVersion,
		Locale:      req.Locale,
		Timezone:    req.Timezone,
		Topics:      req.Topics,
		Metadata:    req.Metadata,
	}

	device, err := h.service.RegisterDevice(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    device,
	})
}

// UnregisterDevice removes a device from push notifications
func (h *Handler) UnregisterDevice(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "UnregisterDevice")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	deviceID := c.Param("id")

	if err := h.service.UnregisterDevice(ctx, tenantID, deviceID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"unregistered": true},
	})
}

// GetUserDevices retrieves devices for a user
func (h *Handler) GetUserDevices(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetUserDevices")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.Param("userId")

	devices, err := h.service.GetUserDevices(ctx, tenantID, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    devices,
	})
}

// SendPush sends a push notification
func (h *Handler) SendPush(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendPush")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	senderID := c.GetString("user_id")

	var req SendPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	actions := make([]domain.PushAction, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = domain.PushAction{
			ID:    a.ID,
			Title: a.Title,
			URL:   a.URL,
		}
	}

	cmd := SendPushCommand{
		TenantID:    tenantID,
		UserID:      req.UserID,
		DeviceID:    req.DeviceID,
		Title:       req.Title,
		Body:        req.Body,
		ImageURL:    req.ImageURL,
		Data:        req.Data,
		Actions:     actions,
		Sound:       req.Sound,
		Badge:       req.Badge,
		Priority:    req.Priority,
		TTL:         req.TTL,
		CollapseKey: req.CollapseKey,
		SenderID:    senderID,
	}

	notification, err := h.service.SendPush(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SendPushToTopic sends a push notification to a topic
func (h *Handler) SendPushToTopic(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SendPushToTopic")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	senderID := c.GetString("user_id")

	var req SendPushToTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	actions := make([]domain.PushAction, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = domain.PushAction{
			ID:    a.ID,
			Title: a.Title,
			URL:   a.URL,
		}
	}

	cmd := SendPushToTopicCommand{
		TenantID: tenantID,
		Topic:    req.Topic,
		Title:    req.Title,
		Body:     req.Body,
		ImageURL: req.ImageURL,
		Data:     req.Data,
		Actions:  actions,
		SenderID: senderID,
	}

	notification, err := h.service.SendPushToTopic(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// SubscribeToTopic subscribes a device to a topic
func (h *Handler) SubscribeToTopic(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "SubscribeToTopic")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	deviceID := c.Param("id")
	topic := c.Param("topic")

	if err := h.service.SubscribeToTopic(ctx, tenantID, deviceID, topic); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"subscribed": true, "topic": topic},
	})
}

// UnsubscribeFromTopic unsubscribes a device from a topic
func (h *Handler) UnsubscribeFromTopic(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "UnsubscribeFromTopic")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	deviceID := c.Param("id")
	topic := c.Param("topic")

	if err := h.service.UnsubscribeFromTopic(ctx, tenantID, deviceID, topic); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"unsubscribed": true, "topic": topic},
	})
}

// ============================================================================
// Voice Handlers
// ============================================================================

// InitiateCall initiates an outbound voice call
func (h *Handler) InitiateCall(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "InitiateCall")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cmd := SendNotificationCommand{
		TenantID:   tenantID,
		Channel:    "voice",
		Recipient:  h.mapRecipient(req.Recipient),
		Content:    h.mapContent(req.Content),
		TemplateID: req.TemplateID,
		Variables:  req.Variables,
		Reference:  req.Reference,
		Metadata:   req.Metadata,
		SenderID:   userID,
	}

	notification, err := h.service.Send(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    notification,
	})
}

// VoiceEventsWebhook handles voice call events
func (h *Handler) VoiceEventsWebhook(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "VoiceEventsWebhook")
	defer span.End()

	var req VoiceCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	h.logger.Info("Voice call event received",
		zap.String("call_id", req.CallID),
		zap.String("status", req.Status),
	)

	// Process call status update
	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ============================================================================
// Template Handlers
// ============================================================================

// CreateTemplate creates a notification template
func (h *Handler) CreateTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "CreateTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	variables := make([]domain.TemplateVariable, len(req.Variables))
	for i, v := range req.Variables {
		variables[i] = domain.TemplateVariable{
			Name:        v.Name,
			Type:        v.Type,
			Required:    v.Required,
			Default:     v.Default,
			Description: v.Description,
		}
	}

	cmd := CreateTemplateCommand{
		TenantID:    tenantID,
		Name:        req.Name,
		Channel:     req.Channel,
		Content:     h.mapTemplateContent(req.Content),
		Description: req.Description,
		Category:    req.Category,
		Language:    req.Language,
		Variables:   variables,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		CreatedBy:   userID,
	}

	template, err := h.service.CreateTemplate(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    template,
	})
}

// ListTemplates lists notification templates
func (h *Handler) ListTemplates(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "ListTemplates")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	query := TemplateListQuery{
		TenantID: tenantID,
		Channel:  c.Query("channel"),
		Category: c.Query("category"),
		Search:   c.Query("search"),
		Page:     h.parseInt(c.Query("page"), 1),
		PageSize: h.parseInt(c.Query("page_size"), 20),
	}

	if active := c.Query("active"); active != "" {
		a := active == "true"
		query.Active = &a
	}

	if tags := c.Query("tags"); tags != "" {
		query.Tags = strings.Split(tags, ",")
	}

	result, err := h.service.ListTemplates(ctx, query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	totalPages := int(result.TotalCount) / query.PageSize
	if int(result.TotalCount)%query.PageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    result.Templates,
		Meta: &APIMeta{
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalCount: result.TotalCount,
			TotalPages: totalPages,
		},
	})
}

// GetTemplate retrieves a template
func (h *Handler) GetTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	template, err := h.service.GetTemplate(ctx, tenantID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    template,
	})
}

// UpdateTemplate updates a template
func (h *Handler) UpdateTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "UpdateTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	variables := make([]domain.TemplateVariable, len(req.Variables))
	for i, v := range req.Variables {
		variables[i] = domain.TemplateVariable{
			Name:        v.Name,
			Type:        v.Type,
			Required:    v.Required,
			Default:     v.Default,
			Description: v.Description,
		}
	}

	var content *domain.TemplateContent
	if req.Content != nil {
		c := h.mapTemplateContent(*req.Content)
		content = &c
	}

	cmd := UpdateTemplateCommand{
		TenantID:    tenantID,
		TemplateID:  id,
		Name:        req.Name,
		Content:     content,
		Description: req.Description,
		Category:    req.Category,
		Language:    req.Language,
		Variables:   variables,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		Active:      req.Active,
		UpdatedBy:   userID,
	}

	template, err := h.service.UpdateTemplate(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    template,
	})
}

// DeleteTemplate deletes a template
func (h *Handler) DeleteTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "DeleteTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	if err := h.service.DeleteTemplate(ctx, tenantID, id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"deleted": true},
	})
}

// RenderTemplate renders a template with variables
func (h *Handler) RenderTemplate(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "RenderTemplate")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")

	var req struct {
		Variables map[string]interface{} `json:"variables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	rendered, err := h.service.RenderTemplate(ctx, tenantID, id, req.Variables)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"rendered": rendered},
	})
}

// GetTemplatePerformance retrieves template performance metrics
func (h *Handler) GetTemplatePerformance(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetTemplatePerformance")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	id := c.Param("id")
	days := h.parseInt(c.DefaultQuery("days", "30"), 30)

	performance, err := h.service.GetTemplatePerformance(ctx, tenantID, id, days)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    performance,
	})
}

// ============================================================================
// Preferences Handlers
// ============================================================================

// GetUserPreferences retrieves notification preferences for a user
func (h *Handler) GetUserPreferences(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetUserPreferences")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.Param("userId")

	preferences, err := h.service.GetUserPreferences(ctx, tenantID, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    preferences,
	})
}

// UpdateUserPreferences updates notification preferences
func (h *Handler) UpdateUserPreferences(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "UpdateUserPreferences")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.Param("userId")

	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var quietHours *domain.QuietHours
	if req.QuietHours != nil {
		quietHours = &domain.QuietHours{
			Enabled:   req.QuietHours.Enabled,
			StartTime: req.QuietHours.StartTime,
			EndTime:   req.QuietHours.EndTime,
			Timezone:  req.QuietHours.Timezone,
			Days:      req.QuietHours.Days,
		}
	}

	cmd := UpdatePreferencesCommand{
		TenantID:      tenantID,
		UserID:        userID,
		Channels:      req.Channels,
		Categories:    req.Categories,
		QuietHours:    quietHours,
		Frequency:     req.Frequency,
		PreferredLang: req.PreferredLang,
	}

	preferences, err := h.service.UpdateUserPreferences(ctx, cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    preferences,
	})
}

// OptOut opts a user out of a channel
func (h *Handler) OptOut(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "OptOut")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.Param("userId")
	channel := c.Param("channel")

	if err := h.service.OptOut(ctx, tenantID, userID, channel); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"opted_out": true, "channel": channel},
	})
}

// OptIn opts a user back into a channel
func (h *Handler) OptIn(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "OptIn")
	defer span.End()

	tenantID := c.GetString("tenant_id")
	userID := c.Param("userId")
	channel := c.Param("channel")

	if err := h.service.OptIn(ctx, tenantID, userID, channel); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    gin.H{"opted_in": true, "channel": channel},
	})
}

// ============================================================================
// Analytics Handlers
// ============================================================================

// GetChannelStats retrieves channel statistics
func (h *Handler) GetChannelStats(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetChannelStats")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	startDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format(time.RFC3339)))
	endDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("end_date", time.Now().Format(time.RFC3339)))

	stats, err := h.service.GetChannelStats(ctx, tenantID, startDate, endDate)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetDeliveryReport retrieves delivery report
func (h *Handler) GetDeliveryReport(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "GetDeliveryReport")
	defer span.End()

	tenantID := c.GetString("tenant_id")

	startDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)))
	endDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("end_date", time.Now().Format(time.RFC3339)))

	query := DeliveryReportQuery{
		TenantID:  tenantID,
		Channel:   c.Query("channel"),
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   c.DefaultQuery("group_by", "day"),
	}

	report, err := h.service.GetDeliveryReport(ctx, query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    report,
	})
}

// ============================================================================
// Helper Methods
// ============================================================================

func (h *Handler) errorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	h.logger.Error("Handler error", zap.Error(err))

	// Map domain errors to HTTP status codes
	switch {
	case strings.Contains(err.Error(), "not found"):
		h.errorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
	case strings.Contains(err.Error(), "already exists"):
		h.errorResponse(c, http.StatusConflict, "ALREADY_EXISTS", err.Error())
	case strings.Contains(err.Error(), "invalid"):
		h.errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	case strings.Contains(err.Error(), "unauthorized"):
		h.errorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	case strings.Contains(err.Error(), "forbidden"):
		h.errorResponse(c, http.StatusForbidden, "FORBIDDEN", err.Error())
	case strings.Contains(err.Error(), "rate limit"):
		h.errorResponse(c, http.StatusTooManyRequests, "RATE_LIMITED", err.Error())
	default:
		h.errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
	}
}

func (h *Handler) parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

func (h *Handler) mapRecipient(r RecipientRequest) domain.Recipient {
	return domain.Recipient{
		Phone:    r.Phone,
		Email:    r.Email,
		UserID:   r.UserID,
		DeviceID: r.DeviceID,
		Name:     r.Name,
	}
}

func (h *Handler) mapContent(c ContentRequest) domain.Content {
	actions := make([]domain.PushAction, len(c.Actions))
	for i, a := range c.Actions {
		actions[i] = domain.PushAction{
			ID:    a.ID,
			Title: a.Title,
			URL:   a.URL,
		}
	}

	return domain.Content{
		Body:          c.Body,
		Subject:       c.Subject,
		HTMLBody:      c.HTMLBody,
		From:          c.From,
		ReplyTo:       c.ReplyTo,
		CC:            c.CC,
		BCC:           c.BCC,
		Title:         c.Title,
		ImageURL:      c.ImageURL,
		Data:          c.Data,
		Actions:       actions,
		Sound:         c.Sound,
		Badge:         c.Badge,
		MediaURL:      c.MediaURL,
		MediaType:     c.MediaType,
		Caption:       c.Caption,
		FileName:      c.FileName,
		VoiceGender:   c.VoiceGender,
		VoiceLanguage: c.VoiceLanguage,
		DTMF:          c.DTMF,
		MaxDigits:     c.MaxDigits,
		CallbackURL:   c.CallbackURL,
	}
}

func (h *Handler) mapTemplateContent(c TemplateContent) domain.TemplateContent {
	buttons := make([]domain.Button, len(c.Buttons))
	for i, b := range c.Buttons {
		buttons[i] = domain.Button{
			Type:    b.Type,
			Text:    b.Text,
			URL:     b.URL,
			Phone:   b.Phone,
			Payload: b.Payload,
		}
	}

	var header *domain.MediaHeader
	if c.Header != nil {
		header = &domain.MediaHeader{
			Type:     c.Header.Type,
			Text:     c.Header.Text,
			MediaURL: c.Header.MediaURL,
		}
	}

	return domain.TemplateContent{
		Body:     c.Body,
		Subject:  c.Subject,
		HTMLBody: c.HTMLBody,
		Title:    c.Title,
		Header:   header,
		Footer:   c.Footer,
		Buttons:  buttons,
	}
}

func (h *Handler) mapTemplateComponent(c TemplateComponent) domain.TemplateComponent {
	params := make([]domain.TemplateParameter, len(c.Parameters))
	for i, p := range c.Parameters {
		params[i] = domain.TemplateParameter{
			Type: p.Type,
			Text: p.Text,
		}
		if p.Currency != nil {
			params[i].Currency = &domain.CurrencyValue{
				FallbackValue: p.Currency.FallbackValue,
				Code:          p.Currency.Code,
				Amount1000:    p.Currency.Amount1000,
			}
		}
		if p.DateTime != nil {
			params[i].DateTime = &domain.DateTimeValue{
				FallbackValue: p.DateTime.FallbackValue,
			}
		}
		if p.Image != nil {
			params[i].Image = &domain.MediaValue{
				Link:    p.Image.Link,
				ID:      p.Image.ID,
				Caption: p.Image.Caption,
			}
		}
	}

	return domain.TemplateComponent{
		Type:       c.Type,
		Parameters: params,
		SubType:    c.SubType,
		Index:      c.Index,
	}
}

func (h *Handler) mapInteractiveHeader(ih *InteractiveHeader) *domain.InteractiveHeader {
	if ih == nil {
		return nil
	}
	return &domain.InteractiveHeader{
		Type:     ih.Type,
		Text:     ih.Text,
		MediaURL: ih.MediaURL,
	}
}

func (h *Handler) mapInteractiveAction(a InteractiveAction) domain.InteractiveAction {
	buttons := make([]domain.InteractiveButton, len(a.Buttons))
	for i, b := range a.Buttons {
		buttons[i] = domain.InteractiveButton{
			Type: b.Type,
			Reply: domain.ButtonReply{
				ID:    b.Reply.ID,
				Title: b.Reply.Title,
			},
		}
	}

	sections := make([]domain.ListSection, len(a.Sections))
	for i, s := range a.Sections {
		rows := make([]domain.ListRow, len(s.Rows))
		for j, r := range s.Rows {
			rows[j] = domain.ListRow{
				ID:          r.ID,
				Title:       r.Title,
				Description: r.Description,
			}
		}
		sections[i] = domain.ListSection{
			Title: s.Title,
			Rows:  rows,
		}
	}

	return domain.InteractiveAction{
		Buttons:           buttons,
		Button:            a.Button,
		Sections:          sections,
		CatalogID:         a.CatalogID,
		ProductRetailerID: a.ProductRetailerID,
	}
}

func (h *Handler) mapUSSDScreen(s USSDScreen) domain.USSDScreen {
	options := make([]domain.USSDOption, len(s.Options))
	for i, o := range s.Options {
		options[i] = domain.USSDOption{
			Key:        o.Key,
			Label:      o.Label,
			NextScreen: o.NextScreen,
			Action:     o.Action,
			ActionData: o.ActionData,
		}
	}

	return domain.USSDScreen{
		ID:         s.ID,
		Type:       s.Type,
		Title:      s.Title,
		Body:       s.Body,
		Options:    options,
		InputType:  s.InputType,
		MinLength:  s.MinLength,
		MaxLength:  s.MaxLength,
		Validation: s.Validation,
		NextScreen: s.NextScreen,
		Action:     s.Action,
		ActionData: s.ActionData,
	}
}

func (h *Handler) mapUSSDScreens(screens []USSDScreen) []domain.USSDScreen {
	result := make([]domain.USSDScreen, len(screens))
	for i, s := range screens {
		result[i] = h.mapUSSDScreen(s)
	}
	return result
}

// ============================================================================
// USSD Response Type
// ============================================================================

// USSDResponse represents the USSD service response
type USSDResponse struct {
	SessionID string `json:"session_id"`
	Type      string `json:"type"` // CON, END
	Message   string `json:"message"`
}
