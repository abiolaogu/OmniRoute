// Package domain contains validation methods for SCE service domain models.
package domain

import (
	"errors"
)

// Additional domain errors for validation
var (
	ErrServiceNameRequired  = errors.New("service name is required")
	ErrServiceNameTooShort  = errors.New("service name must be at least 3 characters")
	ErrServiceNameTooLong   = errors.New("service name must be at most 100 characters")
	ErrTenantIDRequired     = errors.New("tenant ID is required")
	ErrInvalidServiceStatus = errors.New("invalid service status")
	ErrWorkflowRequired     = errors.New("workflow is required for publishing")
	ErrInvalidNodeType      = errors.New("invalid node type")
	ErrNodeIDRequired       = errors.New("node ID is required")
	ErrEdgeSourceRequired   = errors.New("edge source is required")
	ErrEdgeTargetRequired   = errors.New("edge target is required")
	ErrSelfReferencingEdge  = errors.New("edge cannot reference the same node")
	ErrInvalidErrorPolicy   = errors.New("invalid error handling policy")
	ErrVersionRequired      = errors.New("version is required")
)

// Validate validates a ServiceName
func (n ServiceName) Validate() error {
	if n == "" {
		return ErrServiceNameRequired
	}
	if len(n) < 3 {
		return ErrServiceNameTooShort
	}
	if len(n) > 100 {
		return ErrServiceNameTooLong
	}
	return nil
}

// Validate validates a ServiceStatus
func (s ServiceStatus) Validate() error {
	switch s {
	case ServiceStatusDraft, ServiceStatusPublished, ServiceStatusDeprecated, ServiceStatusArchived:
		return nil
	}
	return ErrInvalidServiceStatus
}

// Validate validates a NodeType
func (t NodeType) Validate() error {
	switch t {
	case NodeTypeActivity, NodeTypeSubflow, NodeTypeAIAction, NodeTypeN8N,
		NodeTypeDecision, NodeTypeParallel, NodeTypeWait, NodeTypeHumanTask:
		return nil
	}
	return ErrInvalidNodeType
}

// Validate validates a WorkflowNode
func (n *WorkflowNode) Validate() error {
	if n.ID == "" {
		return ErrNodeIDRequired
	}
	if err := n.Type.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates a WorkflowEdge
func (e *WorkflowEdge) Validate() error {
	if e.Source == "" {
		return ErrEdgeSourceRequired
	}
	if e.Target == "" {
		return ErrEdgeTargetRequired
	}
	if e.Source == e.Target {
		return ErrSelfReferencingEdge
	}
	return nil
}

// Validate validates an ErrorHandlingPolicy
func (p *ErrorHandlingPolicy) Validate() error {
	switch p.OnError {
	case "fail", "compensate", "ignore", "retry":
		return nil
	case "":
		return nil // Default to fail
	}
	return ErrInvalidErrorPolicy
}

// Validate validates a WorkflowGraph
func (g *WorkflowGraph) Validate() error {
	// Validate nodes
	nodeIDs := make(map[string]bool)
	for _, node := range g.Nodes {
		if err := node.Validate(); err != nil {
			return err
		}
		nodeIDs[node.ID] = true
	}

	// Validate edges
	for _, edge := range g.Edges {
		if err := edge.Validate(); err != nil {
			return err
		}
		// Check that source and target nodes exist
		if !nodeIDs[edge.Source] {
			return errors.New("edge references non-existent source node: " + edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return errors.New("edge references non-existent target node: " + edge.Target)
		}
	}

	// Validate error policy
	if g.ErrorPolicy != nil {
		if err := g.ErrorPolicy.Validate(); err != nil {
			return err
		}
	}

	// Check for cycles
	if g.HasCycle() {
		return ErrCyclicWorkflow
	}

	return nil
}

// ValidateForPublish performs additional validation required for publishing
func (s *ServiceDefinition) ValidateForPublish() error {
	// Basic validation
	if err := s.Validate(); err != nil {
		return err
	}

	// Must have a workflow
	if s.CurrentWorkflow == nil {
		return ErrWorkflowRequired
	}

	// Validate workflow
	if err := s.CurrentWorkflow.Validate(); err != nil {
		return err
	}

	// Must have at least one node
	if len(s.CurrentWorkflow.Nodes) == 0 {
		return ErrEmptyWorkflow
	}

	return nil
}
