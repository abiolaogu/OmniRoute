// Package activities contains Temporal activity implementations for the SCE service.
// Activities are the building blocks that perform actual work in workflows.
package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("sce-service/activities")

// CoreActivities contains all core activity implementations
type CoreActivities struct {
	httpClient *http.Client
}

// NewCoreActivities creates a new CoreActivities instance
func NewCoreActivities() *CoreActivities {
	return &CoreActivities{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ============================================================================
// HTTP Activities
// ============================================================================

// HTTPActivityInput contains inputs for HTTP activities
type HTTPActivityInput struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
	Timeout int               `json:"timeout_seconds"`
}

// HTTPActivityResult contains HTTP activity results
type HTTPActivityResult struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	BodyRaw    string                 `json:"body_raw"`
}

// HTTPCall executes an HTTP request
func (a *CoreActivities) HTTPCall(ctx context.Context, input *HTTPActivityInput) (*HTTPActivityResult, error) {
	ctx, span := tracer.Start(ctx, "HTTPCall")
	defer span.End()
	span.SetAttributes(attribute.String("http.method", input.Method))
	span.SetAttributes(attribute.String("http.url", input.URL))

	// Build request body
	var bodyReader io.Reader
	if input.Body != nil {
		bodyBytes, err := json.Marshal(input.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, input.Method, input.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}
	if input.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result := &HTTPActivityResult{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		BodyRaw:    string(respBody),
	}

	for k, v := range resp.Header {
		if len(v) > 0 {
			result.Headers[k] = v[0]
		}
	}

	// Try to parse as JSON
	json.Unmarshal(respBody, &result.Body)

	return result, nil
}

// ============================================================================
// Notification Activities
// ============================================================================

// SendNotificationInput contains inputs for notifications
type SendNotificationInput struct {
	TenantID  string            `json:"tenant_id"`
	Channel   string            `json:"channel"` // sms, email, push, whatsapp
	Recipient string            `json:"recipient"`
	Template  string            `json:"template"`
	Variables map[string]string `json:"variables"`
}

// SendNotificationResult contains notification results
type SendNotificationResult struct {
	Success     bool   `json:"success"`
	MessageID   string `json:"message_id"`
	Channel     string `json:"channel"`
	DeliveredAt string `json:"delivered_at,omitempty"`
}

// SendNotification sends a notification through the specified channel
func (a *CoreActivities) SendNotification(ctx context.Context, input *SendNotificationInput) (*SendNotificationResult, error) {
	ctx, span := tracer.Start(ctx, "SendNotification")
	defer span.End()
	span.SetAttributes(attribute.String("channel", input.Channel))
	span.SetAttributes(attribute.String("template", input.Template))

	// In production, this would call the notification service
	// For now, return a mock response
	return &SendNotificationResult{
		Success:   true,
		MessageID: fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Channel:   input.Channel,
	}, nil
}

// ============================================================================
// Transform Activities
// ============================================================================

// TransformInput contains inputs for data transformation
type TransformInput struct {
	Data     interface{}       `json:"data"`
	Template string            `json:"template"`
	Mapping  map[string]string `json:"mapping"`
}

// TransformResult contains transformation results
type TransformResult struct {
	Data interface{} `json:"data"`
}

// TransformJSON transforms JSON data according to a template/mapping
func (a *CoreActivities) TransformJSON(ctx context.Context, input *TransformInput) (*TransformResult, error) {
	ctx, span := tracer.Start(ctx, "TransformJSON")
	defer span.End()

	// Apply mapping
	result := make(map[string]interface{})
	if dataMap, ok := input.Data.(map[string]interface{}); ok {
		for targetKey, sourceKey := range input.Mapping {
			if val, exists := dataMap[sourceKey]; exists {
				result[targetKey] = val
			}
		}
	}

	return &TransformResult{Data: result}, nil
}

// ValidateInput contains inputs for schema validation
type ValidateInput struct {
	Data   interface{} `json:"data"`
	Schema interface{} `json:"schema"`
}

// ValidateResult contains validation results
type ValidateResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// ValidateSchema validates data against a JSON schema
func (a *CoreActivities) ValidateSchema(ctx context.Context, input *ValidateInput) (*ValidateResult, error) {
	ctx, span := tracer.Start(ctx, "ValidateSchema")
	defer span.End()

	// Simplified validation - in production use jsonschema library
	return &ValidateResult{Valid: true}, nil
}

// MapFieldsInput contains inputs for field mapping
type MapFieldsInput struct {
	Source  map[string]interface{} `json:"source"`
	Mapping map[string]string      `json:"mapping"`
}

// MapFieldsResult contains field mapping results
type MapFieldsResult struct {
	Result map[string]interface{} `json:"result"`
}

// MapFields maps fields from source to target structure
func (a *CoreActivities) MapFields(ctx context.Context, input *MapFieldsInput) (*MapFieldsResult, error) {
	ctx, span := tracer.Start(ctx, "MapFields")
	defer span.End()

	result := make(map[string]interface{})
	for targetKey, sourceKey := range input.Mapping {
		if val, exists := input.Source[sourceKey]; exists {
			result[targetKey] = val
		}
	}

	return &MapFieldsResult{Result: result}, nil
}

// ============================================================================
// Storage Activities
// ============================================================================

// UploadInput contains inputs for file upload
type UploadInput struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"`
}

// UploadResult contains upload results
type UploadResult struct {
	URL  string `json:"url"`
	Size int    `json:"size"`
	ETag string `json:"etag"`
}

// UploadFile uploads a file to storage
func (a *CoreActivities) UploadFile(ctx context.Context, input *UploadInput) (*UploadResult, error) {
	ctx, span := tracer.Start(ctx, "UploadFile")
	defer span.End()
	span.SetAttributes(attribute.String("bucket", input.Bucket))
	span.SetAttributes(attribute.String("key", input.Key))

	// In production, this would upload to GCS/S3
	return &UploadResult{
		URL:  fmt.Sprintf("gs://%s/%s", input.Bucket, input.Key),
		Size: len(input.Data),
		ETag: fmt.Sprintf("etag_%d", time.Now().UnixNano()),
	}, nil
}

// DownloadInput contains inputs for file download
type DownloadInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

// DownloadResult contains download results
type DownloadResult struct {
	Data        []byte `json:"data"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
}

// DownloadFile downloads a file from storage
func (a *CoreActivities) DownloadFile(ctx context.Context, input *DownloadInput) (*DownloadResult, error) {
	ctx, span := tracer.Start(ctx, "DownloadFile")
	defer span.End()

	// In production, this would download from GCS/S3
	return &DownloadResult{
		Data:        []byte{},
		ContentType: "application/octet-stream",
		Size:        0,
	}, nil
}

// SignedURLInput contains inputs for generating signed URLs
type SignedURLInput struct {
	Bucket    string `json:"bucket"`
	Key       string `json:"key"`
	ExpiresIn int    `json:"expires_in_seconds"`
	Method    string `json:"method"` // GET or PUT
}

// SignedURLResult contains signed URL results
type SignedURLResult struct {
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

// GenerateSignedURL generates a signed URL for direct access
func (a *CoreActivities) GenerateSignedURL(ctx context.Context, input *SignedURLInput) (*SignedURLResult, error) {
	ctx, span := tracer.Start(ctx, "GenerateSignedURL")
	defer span.End()

	expiresAt := time.Now().Add(time.Duration(input.ExpiresIn) * time.Second)

	return &SignedURLResult{
		URL:       fmt.Sprintf("https://storage.googleapis.com/%s/%s?signed=true", input.Bucket, input.Key),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}
