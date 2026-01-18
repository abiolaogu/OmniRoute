// Package domain contains the core domain models for the Service Creation Environment.
// This follows DDD principles with aggregates, value objects, and domain events.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Errors
// ============================================================================

var (
	ErrInvalidServiceName    = errors.New("service name must be 3-100 characters")
	ErrEmptyWorkflow         = errors.New("workflow must have at least one node")
	ErrCannotPublishArchived = errors.New("cannot publish archived service")
	ErrCannotArchiveDraft    = errors.New("cannot archive draft service - deprecate first")
	ErrCyclicWorkflow        = errors.New("workflow contains cycles")
	ErrInvalidEdge           = errors.New("edge references non-existent node")
	ErrDecisionNeedsMultiple = errors.New("decision node requires at least 2 outgoing edges")
)

// ============================================================================
// Value Objects
// ============================================================================

// ServiceID is a strongly-typed identifier for services
type ServiceID uuid.UUID

func NewServiceID() ServiceID {
	return ServiceID(uuid.New())
}

func (s ServiceID) String() string {
	return uuid.UUID(s).String()
}

// TenantID is a strongly-typed identifier for tenants
type TenantID uuid.UUID

func (t TenantID) String() string {
	return uuid.UUID(t).String()
}

// VersionID is a strongly-typed identifier for versions
type VersionID uuid.UUID

func NewVersionID() VersionID {
	return VersionID(uuid.New())
}

// UserID is a strongly-typed identifier for users
type UserID uuid.UUID

// ServiceName is a validated service name
type ServiceName string

func (n ServiceName) Validate() error {
	if len(n) < 3 || len(n) > 100 {
		return ErrInvalidServiceName
	}
	return nil
}

// ServiceStatus represents the lifecycle status of a service
type ServiceStatus string

const (
	ServiceStatusDraft      ServiceStatus = "draft"
	ServiceStatusPublished  ServiceStatus = "published"
	ServiceStatusDeprecated ServiceStatus = "deprecated"
	ServiceStatusArchived   ServiceStatus = "archived"
)

// CanTransitionTo checks if a status transition is valid
func (s ServiceStatus) CanTransitionTo(target ServiceStatus) bool {
	transitions := map[ServiceStatus][]ServiceStatus{
		ServiceStatusDraft:      {ServiceStatusPublished},
		ServiceStatusPublished:  {ServiceStatusDeprecated},
		ServiceStatusDeprecated: {ServiceStatusArchived, ServiceStatusPublished},
		ServiceStatusArchived:   {},
	}

	allowed, ok := transitions[s]
	if !ok {
		return false
	}
	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

// ServiceCategory classifies services by their function
type ServiceCategory string

const (
	CategoryAutomation   ServiceCategory = "automation"
	CategoryIntegration  ServiceCategory = "integration"
	CategoryNotification ServiceCategory = "notification"
	CategoryPayment      ServiceCategory = "payment"
	CategoryFulfillment  ServiceCategory = "fulfillment"
	CategoryAnalytics    ServiceCategory = "analytics"
	CategoryCustom       ServiceCategory = "custom"
)

// ============================================================================
// Aggregate Root: ServiceDefinition
// ============================================================================

// ServiceDefinition is the aggregate root for service creation
type ServiceDefinition struct {
	ID            ServiceID        `json:"id"`
	TenantID      TenantID         `json:"tenant_id"`
	Name          ServiceName      `json:"name"`
	Description   string           `json:"description"`
	Category      ServiceCategory  `json:"category"`
	Workflow      WorkflowGraph    `json:"workflow"`
	Versions      []ServiceVersion `json:"versions"`
	ActiveVersion VersionID        `json:"active_version"`
	Status        ServiceStatus    `json:"status"`
	CreatedBy     UserID           `json:"created_by"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	PublishedAt   *time.Time       `json:"published_at,omitempty"`

	// Domain events (transient)
	events []DomainEvent
}

// ServiceVersion represents a specific version of a service
type ServiceVersion struct {
	ID            VersionID     `json:"id"`
	ServiceID     ServiceID     `json:"service_id"`
	VersionNumber int           `json:"version_number"`
	Workflow      WorkflowGraph `json:"workflow"`
	ReleaseNotes  string        `json:"release_notes"`
	CreatedAt     time.Time     `json:"created_at"`
}

// NewServiceDefinition creates a new service with validation
func NewServiceDefinition(tenantID TenantID, name ServiceName, description string, category ServiceCategory, createdBy UserID) (*ServiceDefinition, error) {
	if err := name.Validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	service := &ServiceDefinition{
		ID:          NewServiceID(),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Category:    category,
		Status:      ServiceStatusDraft,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		Workflow:    WorkflowGraph{},
		Versions:    []ServiceVersion{},
	}

	service.events = append(service.events, ServiceCreatedEvent{
		ServiceID: service.ID,
		TenantID:  tenantID,
		Name:      string(name),
		Category:  string(category),
		CreatedBy: createdBy,
		Timestamp: now,
	})

	return service, nil
}

// Publish transitions the service to published status
func (s *ServiceDefinition) Publish() error {
	if s.Status == ServiceStatusArchived {
		return ErrCannotPublishArchived
	}
	if len(s.Workflow.Nodes) == 0 {
		return ErrEmptyWorkflow
	}

	// Validate workflow
	if err := s.Workflow.Validate(); err != nil {
		return err
	}

	s.Status = ServiceStatusPublished
	now := time.Now().UTC()
	s.PublishedAt = &now
	s.UpdatedAt = now

	s.events = append(s.events, ServicePublishedEvent{
		ServiceID:   s.ID,
		TenantID:    s.TenantID,
		Version:     s.ActiveVersion,
		PublishedAt: now,
	})

	return nil
}

// Deprecate marks the service as deprecated
func (s *ServiceDefinition) Deprecate(reason string) error {
	if !s.Status.CanTransitionTo(ServiceStatusDeprecated) {
		return errors.New("cannot deprecate service in current status")
	}

	s.Status = ServiceStatusDeprecated
	s.UpdatedAt = time.Now().UTC()

	s.events = append(s.events, ServiceDeprecatedEvent{
		ServiceID: s.ID,
		Reason:    reason,
		Timestamp: s.UpdatedAt,
	})

	return nil
}

// Archive permanently archives the service
func (s *ServiceDefinition) Archive() error {
	if s.Status == ServiceStatusDraft {
		return ErrCannotArchiveDraft
	}
	if !s.Status.CanTransitionTo(ServiceStatusArchived) {
		return errors.New("cannot archive service in current status")
	}

	s.Status = ServiceStatusArchived
	s.UpdatedAt = time.Now().UTC()

	s.events = append(s.events, ServiceArchivedEvent{
		ServiceID: s.ID,
		Timestamp: s.UpdatedAt,
	})

	return nil
}

// AddVersion creates a new version of the service
func (s *ServiceDefinition) AddVersion(workflow WorkflowGraph, notes string) (VersionID, error) {
	if err := workflow.Validate(); err != nil {
		return VersionID{}, err
	}

	versionNum := s.nextVersionNumber()
	version := ServiceVersion{
		ID:            NewVersionID(),
		ServiceID:     s.ID,
		VersionNumber: versionNum,
		Workflow:      workflow,
		ReleaseNotes:  notes,
		CreatedAt:     time.Now().UTC(),
	}

	s.Versions = append(s.Versions, version)
	s.Workflow = workflow
	s.ActiveVersion = version.ID
	s.UpdatedAt = time.Now().UTC()

	s.events = append(s.events, VersionReleasedEvent{
		ServiceID: s.ID,
		VersionID: version.ID,
		Number:    versionNum,
		Timestamp: version.CreatedAt,
	})

	return version.ID, nil
}

// UpdateWorkflow updates the current workflow (draft mode only)
func (s *ServiceDefinition) UpdateWorkflow(workflow WorkflowGraph) error {
	if s.Status != ServiceStatusDraft {
		return errors.New("can only update workflow in draft status")
	}

	if err := workflow.Validate(); err != nil {
		return err
	}

	s.Workflow = workflow
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// PullEvents returns and clears all domain events
func (s *ServiceDefinition) PullEvents() []DomainEvent {
	events := s.events
	s.events = nil
	return events
}

func (s *ServiceDefinition) nextVersionNumber() int {
	if len(s.Versions) == 0 {
		return 1
	}
	return s.Versions[len(s.Versions)-1].VersionNumber + 1
}

// ============================================================================
// Workflow Graph Value Object
// ============================================================================

// WorkflowGraph represents the visual workflow as a directed acyclic graph
type WorkflowGraph struct {
	Nodes       []WorkflowNode      `json:"nodes"`
	Edges       []WorkflowEdge      `json:"edges"`
	Variables   []WorkflowVariable  `json:"variables"`
	Triggers    []WorkflowTrigger   `json:"triggers"`
	ErrorPolicy ErrorHandlingPolicy `json:"error_policy"`
}

// Validate ensures the workflow is a valid DAG
func (g *WorkflowGraph) Validate() error {
	if len(g.Nodes) == 0 {
		return ErrEmptyWorkflow
	}

	// Build node map
	nodeMap := make(map[string]bool)
	for _, node := range g.Nodes {
		nodeMap[node.ID] = true
	}

	// Validate edges reference existing nodes
	for _, edge := range g.Edges {
		if !nodeMap[edge.Source] || !nodeMap[edge.Target] {
			return ErrInvalidEdge
		}
	}

	// Check for cycles using DFS
	if g.hasCycle() {
		return ErrCyclicWorkflow
	}

	// Validate decision nodes have multiple outputs
	for _, node := range g.Nodes {
		if node.Type == NodeTypeDecision {
			outCount := 0
			for _, edge := range g.Edges {
				if edge.Source == node.ID {
					outCount++
				}
			}
			if outCount < 2 {
				return ErrDecisionNeedsMultiple
			}
		}
	}

	return nil
}

func (g *WorkflowGraph) hasCycle() bool {
	visited := make(map[string]int) // 0=unvisited, 1=visiting, 2=visited

	var dfs func(nodeID string) bool
	dfs = func(nodeID string) bool {
		if visited[nodeID] == 1 {
			return true // Back edge = cycle
		}
		if visited[nodeID] == 2 {
			return false // Already fully processed
		}

		visited[nodeID] = 1
		for _, edge := range g.Edges {
			if edge.Source == nodeID {
				if dfs(edge.Target) {
					return true
				}
			}
		}
		visited[nodeID] = 2
		return false
	}

	for _, node := range g.Nodes {
		if visited[node.ID] == 0 {
			if dfs(node.ID) {
				return true
			}
		}
	}
	return false
}

// TopologicalSort returns nodes in execution order
func (g *WorkflowGraph) TopologicalSort() ([]string, error) {
	if g.hasCycle() {
		return nil, ErrCyclicWorkflow
	}

	// Calculate in-degrees
	inDegree := make(map[string]int)
	for _, node := range g.Nodes {
		inDegree[node.ID] = 0
	}
	for _, edge := range g.Edges {
		inDegree[edge.Target]++
	}

	// Start with zero in-degree nodes
	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var result []string
	for len(queue) > 0 {
		nodeID := queue[0]
		queue = queue[1:]
		result = append(result, nodeID)

		for _, edge := range g.Edges {
			if edge.Source == nodeID {
				inDegree[edge.Target]--
				if inDegree[edge.Target] == 0 {
					queue = append(queue, edge.Target)
				}
			}
		}
	}

	return result, nil
}

// GetNodeByID returns a node by its ID
func (g *WorkflowGraph) GetNodeByID(id string) (*WorkflowNode, bool) {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i], true
		}
	}
	return nil, false
}

// GetOutgoingEdges returns all edges from a node
func (g *WorkflowGraph) GetOutgoingEdges(nodeID string) []WorkflowEdge {
	var edges []WorkflowEdge
	for _, edge := range g.Edges {
		if edge.Source == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// WorkflowNode represents a single node in the workflow
type WorkflowNode struct {
	ID          string       `json:"id"`
	Type        NodeType     `json:"type"`
	Label       string       `json:"label"`
	Position    Position     `json:"position"`
	Config      NodeConfig   `json:"config"`
	RetryPolicy *RetryPolicy `json:"retry_policy,omitempty"`
	Timeout     *int         `json:"timeout_seconds,omitempty"`
}

// NodeType defines the type of workflow node
type NodeType string

const (
	NodeTypeActivity  NodeType = "activity"
	NodeTypeSubflow   NodeType = "subflow"
	NodeTypeAIAction  NodeType = "ai_action"
	NodeTypeN8N       NodeType = "n8n"
	NodeTypeDecision  NodeType = "decision"
	NodeTypeParallel  NodeType = "parallel"
	NodeTypeWait      NodeType = "wait"
	NodeTypeHumanTask NodeType = "human_task"
)

// Position represents x,y coordinates in the visual editor
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeConfig holds type-specific configuration
type NodeConfig struct {
	ActivityName    string                 `json:"activity_name,omitempty"`
	WorkflowID      string                 `json:"workflow_id,omitempty"`
	Provider        string                 `json:"provider,omitempty"`
	Model           string                 `json:"model,omitempty"`
	PromptTemplate  string                 `json:"prompt_template,omitempty"`
	Temperature     float64                `json:"temperature,omitempty"`
	MaxTokens       int                    `json:"max_tokens,omitempty"`
	UseLocalModel   bool                   `json:"use_local_model,omitempty"`
	N8NWorkflowID   string                 `json:"n8n_workflow_id,omitempty"`
	WebhookPath     string                 `json:"webhook_path,omitempty"`
	Condition       string                 `json:"condition,omitempty"`
	Branches        []string               `json:"branches,omitempty"`
	WaitType        string                 `json:"wait_type,omitempty"`
	DurationSeconds int                    `json:"duration_seconds,omitempty"`
	SignalName      string                 `json:"signal_name,omitempty"`
	TaskType        string                 `json:"task_type,omitempty"`
	FormSchema      map[string]interface{} `json:"form_schema,omitempty"`
	InputMapping    map[string]string      `json:"input_mapping,omitempty"`
	OutputVariable  string                 `json:"output_variable,omitempty"`
}

// RetryPolicy defines retry behavior for nodes
type RetryPolicy struct {
	MaxAttempts     int     `json:"max_attempts"`
	InitialInterval string  `json:"initial_interval"`
	BackoffCoeff    float64 `json:"backoff_coefficient"`
	MaxInterval     string  `json:"max_interval"`
}

// WorkflowEdge connects two nodes
type WorkflowEdge struct {
	ID           string         `json:"id"`
	Source       string         `json:"source"`
	Target       string         `json:"target"`
	SourceHandle string         `json:"source_handle,omitempty"`
	TargetHandle string         `json:"target_handle,omitempty"`
	Condition    *EdgeCondition `json:"condition,omitempty"`
}

// EdgeCondition defines when an edge should be followed
type EdgeCondition struct {
	Expression string `json:"expression"`
	Label      string `json:"label,omitempty"`
}

// WorkflowVariable represents a workflow variable
type WorkflowVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Required     bool        `json:"required"`
	Description  string      `json:"description,omitempty"`
}

// WorkflowTrigger defines how a workflow is triggered
type WorkflowTrigger struct {
	Type   TriggerType            `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// TriggerType defines the trigger mechanism
type TriggerType string

const (
	TriggerTypeManual   TriggerType = "manual"
	TriggerTypeSchedule TriggerType = "schedule"
	TriggerTypeWebhook  TriggerType = "webhook"
	TriggerTypeEvent    TriggerType = "event"
)

// ErrorHandlingPolicy defines how errors are handled
type ErrorHandlingPolicy struct {
	OnError                string `json:"on_error"` // fail, compensate, ignore, retry
	CompensationWorkflowID string `json:"compensation_workflow_id,omitempty"`
	NotificationChannel    string `json:"notification_channel,omitempty"`
}

// ============================================================================
// Domain Events
// ============================================================================

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// ServiceCreatedEvent is raised when a new service is created
type ServiceCreatedEvent struct {
	ServiceID ServiceID `json:"service_id"`
	TenantID  TenantID  `json:"tenant_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedBy UserID    `json:"created_by"`
	Timestamp time.Time `json:"timestamp"`
}

func (e ServiceCreatedEvent) EventType() string     { return "service.created" }
func (e ServiceCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

// ServicePublishedEvent is raised when a service is published
type ServicePublishedEvent struct {
	ServiceID   ServiceID `json:"service_id"`
	TenantID    TenantID  `json:"tenant_id"`
	Version     VersionID `json:"version"`
	PublishedAt time.Time `json:"published_at"`
}

func (e ServicePublishedEvent) EventType() string     { return "service.published" }
func (e ServicePublishedEvent) OccurredAt() time.Time { return e.PublishedAt }

// ServiceDeprecatedEvent is raised when a service is deprecated
type ServiceDeprecatedEvent struct {
	ServiceID ServiceID `json:"service_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

func (e ServiceDeprecatedEvent) EventType() string     { return "service.deprecated" }
func (e ServiceDeprecatedEvent) OccurredAt() time.Time { return e.Timestamp }

// ServiceArchivedEvent is raised when a service is archived
type ServiceArchivedEvent struct {
	ServiceID ServiceID `json:"service_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (e ServiceArchivedEvent) EventType() string     { return "service.archived" }
func (e ServiceArchivedEvent) OccurredAt() time.Time { return e.Timestamp }

// VersionReleasedEvent is raised when a new version is released
type VersionReleasedEvent struct {
	ServiceID ServiceID `json:"service_id"`
	VersionID VersionID `json:"version_id"`
	Number    int       `json:"number"`
	Timestamp time.Time `json:"timestamp"`
}

func (e VersionReleasedEvent) EventType() string     { return "service.version_released" }
func (e VersionReleasedEvent) OccurredAt() time.Time { return e.Timestamp }
