// Package domain_test contains unit tests for SCE domain models
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/omniroute/sce-service/internal/domain"
)

func TestNewServiceDefinition(t *testing.T) {
	tests := []struct {
		name    string
		svcName domain.ServiceName
		wantErr bool
	}{
		{"valid name", "Order Processing Service", false},
		{"short name", "AB", true},
		{"valid short", "ABC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := domain.NewServiceDefinition(
				domain.TenantID(uuid.New()),
				tt.svcName,
				"Test service",
				domain.CategoryAutomation,
				domain.UserID(uuid.New()),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewServiceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if svc.Status != domain.ServiceStatusDraft {
					t.Error("Expected draft status for new service")
				}
				events := svc.PullEvents()
				if len(events) != 1 {
					t.Errorf("Expected 1 event, got %d", len(events))
				}
			}
		})
	}
}

func TestServiceDefinition_Publish(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *domain.ServiceDefinition
		wantErr bool
		errType error
	}{
		{
			name: "publish with workflow",
			setup: func() *domain.ServiceDefinition {
				svc := createTestService()
				svc.Workflow = createValidWorkflow()
				return svc
			},
			wantErr: false,
		},
		{
			name: "publish without workflow",
			setup: func() *domain.ServiceDefinition {
				return createTestService()
			},
			wantErr: true,
			errType: domain.ErrEmptyWorkflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setup()
			svc.PullEvents() // Clear creation event

			err := svc.Publish()

			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if svc.Status != domain.ServiceStatusPublished {
					t.Error("Expected published status")
				}
				if svc.PublishedAt == nil {
					t.Error("Expected PublishedAt to be set")
				}
				events := svc.PullEvents()
				if len(events) != 1 {
					t.Errorf("Expected 1 event, got %d", len(events))
				}
			}
		})
	}
}

func TestServiceDefinition_StatusTransitions(t *testing.T) {
	tests := []struct {
		name    string
		from    domain.ServiceStatus
		to      domain.ServiceStatus
		allowed bool
	}{
		{"draft to published", domain.ServiceStatusDraft, domain.ServiceStatusPublished, true},
		{"published to deprecated", domain.ServiceStatusPublished, domain.ServiceStatusDeprecated, true},
		{"deprecated to archived", domain.ServiceStatusDeprecated, domain.ServiceStatusArchived, true},
		{"deprecated back to published", domain.ServiceStatusDeprecated, domain.ServiceStatusPublished, true},
		{"draft to archived", domain.ServiceStatusDraft, domain.ServiceStatusArchived, false},
		{"archived to anything", domain.ServiceStatusArchived, domain.ServiceStatusPublished, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.allowed {
				t.Errorf("CanTransitionTo(%v, %v) = %v, want %v", tt.from, tt.to, got, tt.allowed)
			}
		})
	}
}

func TestServiceDefinition_AddVersion(t *testing.T) {
	svc := createTestService()
	workflow := createValidWorkflow()

	versionID, err := svc.AddVersion(workflow, "Initial release")

	if err != nil {
		t.Fatalf("AddVersion() error = %v", err)
	}

	if len(svc.Versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(svc.Versions))
	}

	if svc.ActiveVersion != versionID {
		t.Error("Expected active version to be updated")
	}

	if svc.Versions[0].VersionNumber != 1 {
		t.Errorf("Expected version 1, got %d", svc.Versions[0].VersionNumber)
	}

	// Add another version
	_, err = svc.AddVersion(workflow, "Bug fixes")
	if err != nil {
		t.Fatalf("AddVersion() error = %v", err)
	}

	if svc.Versions[1].VersionNumber != 2 {
		t.Errorf("Expected version 2, got %d", svc.Versions[1].VersionNumber)
	}
}

func TestWorkflowGraph_Validate(t *testing.T) {
	tests := []struct {
		name    string
		graph   domain.WorkflowGraph
		wantErr error
	}{
		{
			name:    "valid linear workflow",
			graph:   createValidWorkflow(),
			wantErr: nil,
		},
		{
			name: "empty workflow",
			graph: domain.WorkflowGraph{
				Nodes: []domain.WorkflowNode{},
			},
			wantErr: domain.ErrEmptyWorkflow,
		},
		{
			name:    "workflow with cycle",
			graph:   createCyclicWorkflow(),
			wantErr: domain.ErrCyclicWorkflow,
		},
		{
			name:    "invalid edge reference",
			graph:   createInvalidEdgeWorkflow(),
			wantErr: domain.ErrInvalidEdge,
		},
		{
			name:    "decision node without branches",
			graph:   createDecisionWithOneBranch(),
			wantErr: domain.ErrDecisionNeedsMultiple,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.graph.Validate()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err != tt.wantErr {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestWorkflowGraph_TopologicalSort(t *testing.T) {
	graph := domain.WorkflowGraph{
		Nodes: []domain.WorkflowNode{
			{ID: "start", Type: domain.NodeTypeActivity},
			{ID: "process", Type: domain.NodeTypeActivity},
			{ID: "end", Type: domain.NodeTypeActivity},
		},
		Edges: []domain.WorkflowEdge{
			{ID: "e1", Source: "start", Target: "process"},
			{ID: "e2", Source: "process", Target: "end"},
		},
	}

	order, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	// Verify start comes before process, process before end
	indexOf := func(s string) int {
		for i, id := range order {
			if id == s {
				return i
			}
		}
		return -1
	}

	if indexOf("start") >= indexOf("process") {
		t.Error("start should come before process")
	}
	if indexOf("process") >= indexOf("end") {
		t.Error("process should come before end")
	}
}

func TestWorkflowGraph_GetNodeByID(t *testing.T) {
	graph := createValidWorkflow()

	node, found := graph.GetNodeByID("node1")
	if !found {
		t.Error("Expected to find node1")
	}
	if node.ID != "node1" {
		t.Errorf("Expected node1, got %s", node.ID)
	}

	_, found = graph.GetNodeByID("nonexistent")
	if found {
		t.Error("Should not find nonexistent node")
	}
}

func TestDomainEvents(t *testing.T) {
	svc := createTestService()
	svc.Workflow = createValidWorkflow()

	// Check creation event
	events := svc.PullEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType() != "service.created" {
		t.Errorf("Expected service.created, got %s", events[0].EventType())
	}

	// Publish and check event
	svc.Publish()
	events = svc.PullEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType() != "service.published" {
		t.Errorf("Expected service.published, got %s", events[0].EventType())
	}

	// Deprecate and check event
	svc.Deprecate("End of life")
	events = svc.PullEvents()
	if events[0].EventType() != "service.deprecated" {
		t.Errorf("Expected service.deprecated, got %s", events[0].EventType())
	}
}

// Helper functions

func createTestService() *domain.ServiceDefinition {
	svc, _ := domain.NewServiceDefinition(
		domain.TenantID(uuid.New()),
		"Test Service",
		"A test service",
		domain.CategoryAutomation,
		domain.UserID(uuid.New()),
	)
	return svc
}

func createValidWorkflow() domain.WorkflowGraph {
	return domain.WorkflowGraph{
		Nodes: []domain.WorkflowNode{
			{
				ID:    "node1",
				Type:  domain.NodeTypeActivity,
				Label: "Start",
				Config: domain.NodeConfig{
					ActivityName: "StartActivity",
				},
			},
			{
				ID:    "node2",
				Type:  domain.NodeTypeActivity,
				Label: "Process",
				Config: domain.NodeConfig{
					ActivityName: "ProcessActivity",
				},
			},
			{
				ID:    "node3",
				Type:  domain.NodeTypeActivity,
				Label: "End",
				Config: domain.NodeConfig{
					ActivityName: "EndActivity",
				},
			},
		},
		Edges: []domain.WorkflowEdge{
			{ID: "e1", Source: "node1", Target: "node2"},
			{ID: "e2", Source: "node2", Target: "node3"},
		},
		Triggers: []domain.WorkflowTrigger{
			{Type: domain.TriggerTypeManual},
		},
		ErrorPolicy: domain.ErrorHandlingPolicy{
			OnError: "fail",
		},
	}
}

func createCyclicWorkflow() domain.WorkflowGraph {
	return domain.WorkflowGraph{
		Nodes: []domain.WorkflowNode{
			{ID: "a", Type: domain.NodeTypeActivity},
			{ID: "b", Type: domain.NodeTypeActivity},
			{ID: "c", Type: domain.NodeTypeActivity},
		},
		Edges: []domain.WorkflowEdge{
			{ID: "e1", Source: "a", Target: "b"},
			{ID: "e2", Source: "b", Target: "c"},
			{ID: "e3", Source: "c", Target: "a"}, // Creates cycle
		},
	}
}

func createInvalidEdgeWorkflow() domain.WorkflowGraph {
	return domain.WorkflowGraph{
		Nodes: []domain.WorkflowNode{
			{ID: "a", Type: domain.NodeTypeActivity},
		},
		Edges: []domain.WorkflowEdge{
			{ID: "e1", Source: "a", Target: "nonexistent"},
		},
	}
}

func createDecisionWithOneBranch() domain.WorkflowGraph {
	return domain.WorkflowGraph{
		Nodes: []domain.WorkflowNode{
			{ID: "start", Type: domain.NodeTypeActivity},
			{ID: "decision", Type: domain.NodeTypeDecision},
			{ID: "end", Type: domain.NodeTypeActivity},
		},
		Edges: []domain.WorkflowEdge{
			{ID: "e1", Source: "start", Target: "decision"},
			{ID: "e2", Source: "decision", Target: "end"}, // Only one branch
		},
	}
}

var _ = time.Now // Suppress unused import warning
