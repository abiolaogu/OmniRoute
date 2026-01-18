// Package domain contains the core domain models for the SCE Service
// Implements DDD Aggregate Root pattern for Service Definitions
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// VALUE OBJECTS
// =============================================================================

// ServiceID is a unique identifier for a service definition
type ServiceID uuid.UUID

func NewServiceID() ServiceID { return ServiceID(uuid.New()) }

func (id ServiceID) String() string { return uuid.UUID(id).String() }

// TenantID identifies the tenant owner
type TenantID uuid.UUID

func (id TenantID) String() string { return uuid.UUID(id).String() }

// UserID identifies the user
type UserID uuid.UUID

func (id UserID) String() string { return uuid.UUID(id).String() }

// VersionID identifies a service version
type VersionID uuid.UUID

func NewVersionID() VersionID { return VersionID(uuid.New()) }

// ServiceName is a validated name for a service
type ServiceName string

func NewServiceName(name string) (ServiceName, error) {
	if len(name) < 3 || len(name) > 100 {
		return "", errors.New("service name must be between 3 and 100 characters")
	}
	return ServiceName(name), nil
}

// ServiceCategory categorizes services
type ServiceCategory string

const (
	CategoryOrder       ServiceCategory = "ORDER"
	CategoryPayment     ServiceCategory = "PAYMENT"
	CategoryInventory   ServiceCategory = "INVENTORY"
	CategoryLogistics   ServiceCategory = "LOGISTICS"
	CategoryAnalytics   ServiceCategory = "ANALYTICS"
	CategoryIntegration ServiceCategory = "INTEGRATION"
)

// ServiceStatus represents the lifecycle status
type ServiceStatus string

const (
	StatusDraft      ServiceStatus = "DRAFT"
	StatusPublished  ServiceStatus = "PUBLISHED"
	StatusDeprecated ServiceStatus = "DEPRECATED"
	StatusArchived   ServiceStatus = "ARCHIVED"
)

// =============================================================================
// DOMAIN EVENTS
// =============================================================================

// DomainEvent interface for all domain events
type DomainEvent interface {
	OccurredAt() time.Time
	AggregateID() string
}

type BaseEvent struct {
	occurredAt  time.Time
	aggregateID string
}

func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) AggregateID() string   { return e.aggregateID }

type ServiceCreatedEvent struct {
	BaseEvent
	Name      ServiceName
	Category  ServiceCategory
	CreatedBy UserID
}

type ServicePublishedEvent struct {
	BaseEvent
	Version   VersionID
	Changelog string
}

type ServiceDeprecatedEvent struct {
	BaseEvent
	Reason          string
	MigrationTarget *ServiceID
}

type WorkflowNodeAddedEvent struct {
	BaseEvent
	NodeID   string
	NodeType string
}

// =============================================================================
// AGGREGATE ROOT: SERVICE DEFINITION
// =============================================================================

// WorkflowNode represents a node in the workflow graph
type WorkflowNode struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Label       string                 `json:"label"`
	Config      map[string]interface{} `json:"config"`
	Position    Position               `json:"position"`
	RetryPolicy *RetryPolicy           `json:"retry_policy,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RetryPolicy struct {
	MaxAttempts        int           `json:"max_attempts"`
	InitialInterval    time.Duration `json:"initial_interval"`
	MaxInterval        time.Duration `json:"max_interval"`
	BackoffCoefficient float64       `json:"backoff_coefficient"`
}

// WorkflowEdge represents a connection between nodes
type WorkflowEdge struct {
	ID        string  `json:"id"`
	Source    string  `json:"source"`
	Target    string  `json:"target"`
	Condition *string `json:"condition,omitempty"`
	Label     *string `json:"label,omitempty"`
}

// WorkflowGraph represents the complete workflow
type WorkflowGraph struct {
	Nodes []WorkflowNode `json:"nodes"`
	Edges []WorkflowEdge `json:"edges"`
}

func (wg *WorkflowGraph) HasNodes() bool { return len(wg.Nodes) > 0 }

func (wg *WorkflowGraph) Validate() error {
	if !wg.HasNodes() {
		return errors.New("workflow must have at least one node")
	}
	// Check for start node
	hasStart := false
	for _, n := range wg.Nodes {
		if n.Type == "start" {
			hasStart = true
			break
		}
	}
	if !hasStart {
		return errors.New("workflow must have a start node")
	}
	return nil
}

// ServiceVersion represents a published version
type ServiceVersion struct {
	ID          VersionID     `json:"id"`
	Number      int           `json:"number"`
	Changelog   string        `json:"changelog"`
	Workflow    WorkflowGraph `json:"workflow"`
	PublishedAt time.Time     `json:"published_at"`
	PublishedBy UserID        `json:"published_by"`
}

// ServiceDefinition is the Aggregate Root
type ServiceDefinition struct {
	id            ServiceID
	tenantID      TenantID
	name          ServiceName
	description   string
	category      ServiceCategory
	workflow      WorkflowGraph
	versions      []ServiceVersion
	activeVersion *VersionID
	status        ServiceStatus
	createdBy     UserID
	createdAt     time.Time
	updatedAt     time.Time
	publishedAt   *time.Time

	// Domain events
	events []DomainEvent
}

// Constructor
func NewServiceDefinition(tenantID TenantID, name ServiceName, category ServiceCategory, createdBy UserID) *ServiceDefinition {
	sd := &ServiceDefinition{
		id:        NewServiceID(),
		tenantID:  tenantID,
		name:      name,
		category:  category,
		status:    StatusDraft,
		createdBy: createdBy,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		versions:  []ServiceVersion{},
		events:    []DomainEvent{},
	}

	sd.events = append(sd.events, ServiceCreatedEvent{
		BaseEvent: BaseEvent{occurredAt: time.Now(), aggregateID: sd.id.String()},
		Name:      name,
		Category:  category,
		CreatedBy: createdBy,
	})

	return sd
}

// Getters
func (sd *ServiceDefinition) ID() ServiceID           { return sd.id }
func (sd *ServiceDefinition) TenantID() TenantID      { return sd.tenantID }
func (sd *ServiceDefinition) Name() ServiceName       { return sd.name }
func (sd *ServiceDefinition) Status() ServiceStatus   { return sd.status }
func (sd *ServiceDefinition) Workflow() WorkflowGraph { return sd.workflow }
func (sd *ServiceDefinition) Events() []DomainEvent   { return sd.events }

// ClearEvents clears domain events after publishing
func (sd *ServiceDefinition) ClearEvents() { sd.events = nil }

// =============================================================================
// AGGREGATE METHODS (Business Logic with Invariant Enforcement)
// =============================================================================

// AddWorkflowNode adds a node to the workflow
func (sd *ServiceDefinition) AddWorkflowNode(node WorkflowNode) error {
	if sd.status != StatusDraft {
		return errors.New("can only modify workflow in draft status")
	}

	// Check for duplicate ID
	for _, n := range sd.workflow.Nodes {
		if n.ID == node.ID {
			return errors.New("node ID already exists")
		}
	}

	sd.workflow.Nodes = append(sd.workflow.Nodes, node)
	sd.updatedAt = time.Now()

	sd.events = append(sd.events, WorkflowNodeAddedEvent{
		BaseEvent: BaseEvent{occurredAt: time.Now(), aggregateID: sd.id.String()},
		NodeID:    node.ID,
		NodeType:  node.Type,
	})

	return nil
}

// AddWorkflowEdge adds an edge to the workflow
func (sd *ServiceDefinition) AddWorkflowEdge(edge WorkflowEdge) error {
	if sd.status != StatusDraft {
		return errors.New("can only modify workflow in draft status")
	}

	// Validate source and target exist
	sourceFound, targetFound := false, false
	for _, n := range sd.workflow.Nodes {
		if n.ID == edge.Source {
			sourceFound = true
		}
		if n.ID == edge.Target {
			targetFound = true
		}
	}

	if !sourceFound || !targetFound {
		return errors.New("edge references non-existent nodes")
	}

	sd.workflow.Edges = append(sd.workflow.Edges, edge)
	sd.updatedAt = time.Now()

	return nil
}

// Publish publishes a new version of the service
func (sd *ServiceDefinition) Publish(changelog string, publishedBy UserID) error {
	// Invariant: Cannot publish without workflow nodes
	if !sd.workflow.HasNodes() {
		return errors.New("cannot publish service without workflow nodes")
	}

	// Invariant: Validate workflow
	if err := sd.workflow.Validate(); err != nil {
		return err
	}

	// Invariant: Status transition
	if sd.status != StatusDraft && sd.status != StatusPublished {
		return errors.New("can only publish from draft or published status")
	}

	versionNumber := len(sd.versions) + 1
	versionID := NewVersionID()

	version := ServiceVersion{
		ID:          versionID,
		Number:      versionNumber,
		Changelog:   changelog,
		Workflow:    sd.workflow, // Snapshot
		PublishedAt: time.Now(),
		PublishedBy: publishedBy,
	}

	sd.versions = append(sd.versions, version)
	sd.activeVersion = &versionID
	sd.status = StatusPublished
	now := time.Now()
	sd.publishedAt = &now
	sd.updatedAt = now

	sd.events = append(sd.events, ServicePublishedEvent{
		BaseEvent: BaseEvent{occurredAt: time.Now(), aggregateID: sd.id.String()},
		Version:   versionID,
		Changelog: changelog,
	})

	return nil
}

// Deprecate marks the service as deprecated
func (sd *ServiceDefinition) Deprecate(reason string, migrationTarget *ServiceID) error {
	if sd.status != StatusPublished {
		return errors.New("can only deprecate published services")
	}

	sd.status = StatusDeprecated
	sd.updatedAt = time.Now()

	sd.events = append(sd.events, ServiceDeprecatedEvent{
		BaseEvent:       BaseEvent{occurredAt: time.Now(), aggregateID: sd.id.String()},
		Reason:          reason,
		MigrationTarget: migrationTarget,
	})

	return nil
}

// Archive archives the service
func (sd *ServiceDefinition) Archive(activeInstanceCount int) error {
	// Invariant: Cannot archive with active instances
	if activeInstanceCount > 0 {
		return errors.New("cannot archive service with active workflow instances")
	}

	if sd.status != StatusDeprecated {
		return errors.New("can only archive deprecated services")
	}

	sd.status = StatusArchived
	sd.updatedAt = time.Now()

	return nil
}
