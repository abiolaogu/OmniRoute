// Package domain tests - Following XP Test-First Development
package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceDefinition(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (TenantID, ServiceName, ServiceCategory, UserID)
		wantErr  bool
		validate func(*testing.T, *ServiceDefinition)
	}{
		{
			name: "creates valid service definition",
			setup: func() (TenantID, ServiceName, ServiceCategory, UserID) {
				return TenantID(uuid.New()), ServiceName("Order Processor"), CategoryOrder, UserID(uuid.New())
			},
			wantErr: false,
			validate: func(t *testing.T, sd *ServiceDefinition) {
				assert.Equal(t, StatusDraft, sd.Status())
				assert.NotEmpty(t, sd.ID().String())
				assert.Len(t, sd.Events(), 1) // ServiceCreatedEvent
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenantID, name, category, userID := tt.setup()
			sd := NewServiceDefinition(tenantID, name, category, userID)

			require.NotNil(t, sd)
			tt.validate(t, sd)
		})
	}
}

func TestServiceName_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "Order Processor", false},
		{"too short", "AB", true},
		{"too long", string(make([]byte, 101)), true},
		{"minimum length", "ABC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServiceName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceDefinition_AddWorkflowNode(t *testing.T) {
	sd := createTestService(t)

	node := WorkflowNode{
		ID:       "start-node",
		Type:     "start",
		Label:    "Start",
		Position: Position{X: 100, Y: 100},
	}

	err := sd.AddWorkflowNode(node)
	require.NoError(t, err)
	assert.Len(t, sd.Workflow().Nodes, 1)

	// Test duplicate node ID
	err = sd.AddWorkflowNode(node)
	assert.Error(t, err)
}

func TestServiceDefinition_AddWorkflowNode_NotDraft(t *testing.T) {
	sd := createPublishedService(t)

	node := WorkflowNode{ID: "new-node", Type: "activity"}
	err := sd.AddWorkflowNode(node)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "draft status")
}

func TestServiceDefinition_Publish(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *ServiceDefinition
		wantErr bool
		errMsg  string
	}{
		{
			name: "publishes valid service",
			setup: func() *ServiceDefinition {
				sd := createTestService(t)
				addStartNode(sd)
				return sd
			},
			wantErr: false,
		},
		{
			name: "fails without workflow nodes",
			setup: func() *ServiceDefinition {
				return createTestService(t)
			},
			wantErr: true,
			errMsg:  "without workflow nodes",
		},
		{
			name: "fails without start node",
			setup: func() *ServiceDefinition {
				sd := createTestService(t)
				sd.AddWorkflowNode(WorkflowNode{ID: "n1", Type: "activity"})
				return sd
			},
			wantErr: true,
			errMsg:  "start node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := tt.setup()
			err := sd.Publish("Initial release", UserID(uuid.New()))

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, StatusPublished, sd.Status())
			}
		})
	}
}

func TestServiceDefinition_StatusTransitions(t *testing.T) {
	// Draft -> Published -> Deprecated -> Archived
	sd := createTestService(t)
	addStartNode(sd)

	// Publish
	err := sd.Publish("v1", UserID(uuid.New()))
	require.NoError(t, err)
	assert.Equal(t, StatusPublished, sd.Status())

	// Deprecate
	err = sd.Deprecate("Replaced by v2", nil)
	require.NoError(t, err)
	assert.Equal(t, StatusDeprecated, sd.Status())

	// Archive (with 0 active instances)
	err = sd.Archive(0)
	require.NoError(t, err)
	assert.Equal(t, StatusArchived, sd.Status())
}

func TestServiceDefinition_CannotArchiveWithActiveInstances(t *testing.T) {
	sd := createTestService(t)
	addStartNode(sd)
	sd.Publish("v1", UserID(uuid.New()))
	sd.Deprecate("test", nil)

	err := sd.Archive(5) // 5 active instances
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "active workflow instances")
}

func TestWorkflowGraph_Validate(t *testing.T) {
	tests := []struct {
		name    string
		graph   WorkflowGraph
		wantErr bool
	}{
		{
			name:    "empty graph",
			graph:   WorkflowGraph{},
			wantErr: true,
		},
		{
			name: "no start node",
			graph: WorkflowGraph{
				Nodes: []WorkflowNode{{ID: "n1", Type: "activity"}},
			},
			wantErr: true,
		},
		{
			name: "valid graph with start",
			graph: WorkflowGraph{
				Nodes: []WorkflowNode{{ID: "start", Type: "start"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.graph.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func createTestService(t *testing.T) *ServiceDefinition {
	name, _ := NewServiceName("Test Service")
	return NewServiceDefinition(
		TenantID(uuid.New()),
		name,
		CategoryOrder,
		UserID(uuid.New()),
	)
}

func addStartNode(sd *ServiceDefinition) {
	sd.AddWorkflowNode(WorkflowNode{
		ID:       "start",
		Type:     "start",
		Label:    "Start",
		Position: Position{X: 100, Y: 100},
	})
}

func createPublishedService(t *testing.T) *ServiceDefinition {
	sd := createTestService(t)
	addStartNode(sd)
	err := sd.Publish("v1", UserID(uuid.New()))
	require.NoError(t, err)
	return sd
}
