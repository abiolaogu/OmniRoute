// Package domain contains repository interfaces for the SCE Service.
// Following DDD principles, repository interfaces are defined in the domain layer.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ServiceDefinitionRepository defines operations for service definition persistence
type ServiceDefinitionRepository interface {
	// FindByID retrieves a service definition by ID
	FindByID(ctx context.Context, tenantID uuid.UUID, serviceID ServiceID) (*ServiceDefinition, error)

	// FindByName retrieves a service definition by name
	FindByName(ctx context.Context, tenantID uuid.UUID, name ServiceName) (*ServiceDefinition, error)

	// FindByStatus retrieves service definitions by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status ServiceStatus, limit, offset int) ([]*ServiceDefinition, error)

	// FindByCategory retrieves service definitions by category
	FindByCategory(ctx context.Context, tenantID uuid.UUID, category ServiceCategory, limit, offset int) ([]*ServiceDefinition, error)

	// FindPublished retrieves published service definitions
	FindPublished(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*ServiceDefinition, error)

	// FindAll retrieves all service definitions for a tenant
	FindAll(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*ServiceDefinition, error)

	// Save persists a service definition
	Save(ctx context.Context, service *ServiceDefinition) error

	// Update updates a service definition
	Update(ctx context.Context, service *ServiceDefinition) error

	// Delete removes a service definition
	Delete(ctx context.Context, tenantID uuid.UUID, serviceID ServiceID) error

	// UpdateStatus updates service status
	UpdateStatus(ctx context.Context, serviceID ServiceID, status ServiceStatus) error

	// UpdateWorkflow updates the workflow graph
	UpdateWorkflow(ctx context.Context, serviceID ServiceID, workflow *WorkflowGraph) error
}

// ServiceVersionRepository defines operations for service version persistence
type ServiceVersionRepository interface {
	// FindByID retrieves a version by ID
	FindByID(ctx context.Context, versionID VersionID) (*ServiceVersion, error)

	// FindByServiceID retrieves versions for a service
	FindByServiceID(ctx context.Context, serviceID ServiceID) ([]*ServiceVersion, error)

	// FindActive retrieves the active version for a service
	FindActive(ctx context.Context, serviceID ServiceID) (*ServiceVersion, error)

	// FindByVersionNumber retrieves a specific version
	FindByVersionNumber(ctx context.Context, serviceID ServiceID, major, minor, patch int) (*ServiceVersion, error)

	// Save persists a version
	Save(ctx context.Context, version *ServiceVersion) error

	// Activate activates a version (deactivating others)
	Activate(ctx context.Context, versionID VersionID) error
}

// ServiceExecutionRepository defines operations for service execution persistence
type ServiceExecutionRepository interface {
	// FindByID retrieves an execution by ID
	FindByID(ctx context.Context, executionID string) (*ServiceExecution, error)

	// FindByServiceID retrieves executions for a service
	FindByServiceID(ctx context.Context, serviceID ServiceID, limit, offset int) ([]*ServiceExecution, error)

	// FindByTenant retrieves executions for a tenant
	FindByTenant(ctx context.Context, tenantID uuid.UUID, from, to time.Time, limit, offset int) ([]*ServiceExecution, error)

	// FindByStatus retrieves executions by status
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status string, limit int) ([]*ServiceExecution, error)

	// FindRunning retrieves running executions
	FindRunning(ctx context.Context, tenantID uuid.UUID) ([]*ServiceExecution, error)

	// Save persists an execution
	Save(ctx context.Context, execution *ServiceExecution) error

	// Update updates an execution
	Update(ctx context.Context, execution *ServiceExecution) error

	// UpdateStatus updates execution status
	UpdateStatus(ctx context.Context, executionID string, status string) error
}

// WorkflowTemplateRepository defines operations for workflow template persistence
type WorkflowTemplateRepository interface {
	// FindByID retrieves a template by ID
	FindByID(ctx context.Context, templateID uuid.UUID) (*WorkflowTemplate, error)

	// FindByCategory retrieves templates by category
	FindByCategory(ctx context.Context, category string) ([]*WorkflowTemplate, error)

	// FindAll retrieves all templates
	FindAll(ctx context.Context, limit, offset int) ([]*WorkflowTemplate, error)

	// Save persists a template
	Save(ctx context.Context, template *WorkflowTemplate) error

	// Update updates a template
	Update(ctx context.Context, template *WorkflowTemplate) error

	// Delete removes a template
	Delete(ctx context.Context, templateID uuid.UUID) error
}

// ServiceExecution represents a service execution record
type ServiceExecution struct {
	ID            string                 `json:"id"`
	ServiceID     ServiceID              `json:"service_id"`
	VersionID     VersionID              `json:"version_id"`
	TenantID      uuid.UUID              `json:"tenant_id"`
	Status        string                 `json:"status"`
	InputData     map[string]interface{} `json:"input_data"`
	OutputData    map[string]interface{} `json:"output_data"`
	Error         *string                `json:"error,omitempty"`
	StartedAt     time.Time              `json:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	WorkflowRunID string                 `json:"workflow_run_id"`
}

// ServiceVersion represents a version of a service
type ServiceVersion struct {
	ID        VersionID      `json:"id"`
	ServiceID ServiceID      `json:"service_id"`
	Major     int            `json:"major"`
	Minor     int            `json:"minor"`
	Patch     int            `json:"patch"`
	Workflow  *WorkflowGraph `json:"workflow"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy UserID         `json:"created_by"`
}

// WorkflowTemplate represents a reusable workflow template
type WorkflowTemplate struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Workflow    *WorkflowGraph `json:"workflow"`
	IsPublic    bool           `json:"is_public"`
	CreatedAt   time.Time      `json:"created_at"`
}
