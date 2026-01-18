// Package activities contains AI/LLM integration activities for the SCE service.
package activities

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// AIActivities contains AI-related activity implementations
type AIActivities struct {
	// In production: LLM clients, local model endpoints, etc.
}

// NewAIActivities creates a new AIActivities instance
func NewAIActivities() *AIActivities {
	return &AIActivities{}
}

// LLMInput contains inputs for LLM calls
type LLMInput struct {
	Provider    string    `json:"provider"` // anthropic, openai, google
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	System      string    `json:"system,omitempty"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResult contains LLM call results
type LLMResult struct {
	Content      string `json:"content"`
	Model        string `json:"model"`
	TokensUsed   int    `json:"tokens_used"`
	FinishReason string `json:"finish_reason"`
}

// CallLLM invokes a cloud LLM (Claude, GPT-4, Gemini)
func (a *AIActivities) CallLLM(ctx context.Context, input *LLMInput) (*LLMResult, error) {
	ctx, span := tracer.Start(ctx, "CallLLM")
	defer span.End()
	span.SetAttributes(attribute.String("provider", input.Provider))
	span.SetAttributes(attribute.String("model", input.Model))

	// In production, this would call the actual LLM APIs
	// For now, return a mock response

	var content string
	switch input.Provider {
	case "anthropic":
		content = fmt.Sprintf("[Claude %s] Response to: %s", input.Model, getLastMessage(input.Messages))
	case "openai":
		content = fmt.Sprintf("[GPT %s] Response to: %s", input.Model, getLastMessage(input.Messages))
	case "google":
		content = fmt.Sprintf("[Gemini %s] Response to: %s", input.Model, getLastMessage(input.Messages))
	default:
		content = fmt.Sprintf("[%s/%s] Response generated", input.Provider, input.Model)
	}

	return &LLMResult{
		Content:      content,
		Model:        input.Model,
		TokensUsed:   len(content) / 4, // Approximate
		FinishReason: "stop",
	}, nil
}

// LocalModelInput contains inputs for local model inference
type LocalModelInput struct {
	Model       string  `json:"model"` // llama-3.3-70b, mistral-large, qwen-coder
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// LocalModelInference invokes a local model via vLLM
func (a *AIActivities) LocalModelInference(ctx context.Context, input *LocalModelInput) (*LLMResult, error) {
	ctx, span := tracer.Start(ctx, "LocalModelInference")
	defer span.End()
	span.SetAttributes(attribute.String("model", input.Model))

	// In production, this would call the local vLLM endpoint
	return &LLMResult{
		Content:      fmt.Sprintf("[Local %s] Response to prompt", input.Model),
		Model:        input.Model,
		TokensUsed:   100,
		FinishReason: "stop",
	}, nil
}

// EmbeddingInput contains inputs for embedding generation
type EmbeddingInput struct {
	Model string   `json:"model"`
	Texts []string `json:"texts"`
}

// EmbeddingResult contains embedding results
type EmbeddingResult struct {
	Embeddings [][]float64 `json:"embeddings"`
	Model      string      `json:"model"`
	Dimensions int         `json:"dimensions"`
}

// GenerateEmbedding creates vector embeddings
func (a *AIActivities) GenerateEmbedding(ctx context.Context, input *EmbeddingInput) (*EmbeddingResult, error) {
	ctx, span := tracer.Start(ctx, "GenerateEmbedding")
	defer span.End()
	span.SetAttributes(attribute.Int("text_count", len(input.Texts)))

	// In production, this would call the embedding model
	embeddings := make([][]float64, len(input.Texts))
	dimensions := 1536 // OpenAI ada-002 dimensions

	for i := range embeddings {
		embeddings[i] = make([]float64, dimensions)
		// Mock: fill with zeros
	}

	return &EmbeddingResult{
		Embeddings: embeddings,
		Model:      input.Model,
		Dimensions: dimensions,
	}, nil
}

// StructuredOutputInput contains inputs for structured output generation
type StructuredOutputInput struct {
	Provider string      `json:"provider"`
	Model    string      `json:"model"`
	Prompt   string      `json:"prompt"`
	Schema   interface{} `json:"schema"` // JSON Schema
}

// StructuredOutputResult contains structured output results
type StructuredOutputResult struct {
	Data  interface{} `json:"data"`
	Valid bool        `json:"valid"`
	Model string      `json:"model"`
}

// StructuredOutput generates structured data from LLM
func (a *AIActivities) StructuredOutput(ctx context.Context, input *StructuredOutputInput) (*StructuredOutputResult, error) {
	ctx, span := tracer.Start(ctx, "StructuredOutput")
	defer span.End()

	// In production, this would use function calling or JSON mode
	return &StructuredOutputResult{
		Data:  map[string]interface{}{"generated": true},
		Valid: true,
		Model: input.Model,
	}, nil
}

// ============================================================================
// n8n Integration Activities
// ============================================================================

// N8NActivities contains n8n-related activity implementations
type N8NActivities struct {
	// In production: n8n API client
	baseURL string
}

// NewN8NActivities creates a new N8NActivities instance
func NewN8NActivities(baseURL string) *N8NActivities {
	return &N8NActivities{baseURL: baseURL}
}

// N8NExecutionInput contains inputs for n8n workflow execution
type N8NExecutionInput struct {
	WorkflowID        string                 `json:"workflow_id"`
	WebhookPath       string                 `json:"webhook_path,omitempty"`
	InputData         map[string]interface{} `json:"input_data"`
	WaitForCompletion bool                   `json:"wait_for_completion"`
	Timeout           int                    `json:"timeout_seconds"`
}

// N8NExecutionResult contains n8n execution results
type N8NExecutionResult struct {
	ExecutionID string                 `json:"execution_id"`
	Status      string                 `json:"status"`
	Data        map[string]interface{} `json:"data,omitempty"`
	StartedAt   string                 `json:"started_at"`
	FinishedAt  string                 `json:"finished_at,omitempty"`
}

// ExecuteN8NWorkflow triggers an n8n workflow
func (a *N8NActivities) ExecuteN8NWorkflow(ctx context.Context, input *N8NExecutionInput) (*N8NExecutionResult, error) {
	ctx, span := tracer.Start(ctx, "ExecuteN8NWorkflow")
	defer span.End()
	span.SetAttributes(attribute.String("workflow_id", input.WorkflowID))

	// In production, this would call the n8n API
	now := time.Now()

	return &N8NExecutionResult{
		ExecutionID: fmt.Sprintf("exec_%d", now.UnixNano()),
		Status:      "success",
		Data:        input.InputData,
		StartedAt:   now.Format(time.RFC3339),
		FinishedAt:  now.Add(time.Second).Format(time.RFC3339),
	}, nil
}

// N8NWorkflowMeta contains n8n workflow metadata
type N8NWorkflowMeta struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	Tags        []string `json:"tags"`
}

// ListN8NWorkflows returns available n8n workflows for tenant
func (a *N8NActivities) ListN8NWorkflows(ctx context.Context, tenantID string) ([]N8NWorkflowMeta, error) {
	ctx, span := tracer.Start(ctx, "ListN8NWorkflows")
	defer span.End()
	span.SetAttributes(attribute.String("tenant_id", tenantID))

	// In production, this would query n8n API
	return []N8NWorkflowMeta{
		{ID: "wf1", Name: "Send Email", Description: "Send transactional email", Active: true, Tags: []string{"notification"}},
		{ID: "wf2", Name: "Sync CRM", Description: "Sync data with CRM", Active: true, Tags: []string{"integration"}},
		{ID: "wf3", Name: "Process Order", Description: "Order processing workflow", Active: true, Tags: []string{"automation"}},
	}, nil
}

// Helper functions

func getLastMessage(messages []Message) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[len(messages)-1].Content
}
