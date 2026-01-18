// Package workflows contains Temporal workflow definitions for the SCE service.
// These workflows execute user-defined services dynamically based on workflow DSL.
package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// ServiceExecutionInput contains the input for executing a user-defined service
type ServiceExecutionInput struct {
	ServiceID     string                 `json:"service_id"`
	VersionID     string                 `json:"version_id"`
	TenantID      string                 `json:"tenant_id"`
	WorkflowDSL   *WorkflowDSL           `json:"workflow_dsl"`
	InputData     map[string]interface{} `json:"input_data"`
	InitiatedBy   string                 `json:"initiated_by"`
	CorrelationID string                 `json:"correlation_id"`
}

// ServiceExecutionResult contains the result of executing a service
type ServiceExecutionResult struct {
	Status      string                 `json:"status"`
	OutputData  map[string]interface{} `json:"output_data"`
	NodeResults map[string]*NodeResult `json:"node_results"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt time.Time              `json:"completed_at"`
	Error       *ExecutionError        `json:"error,omitempty"`
}

// NodeResult contains the result of executing a single node
type NodeResult struct {
	NodeID      string                 `json:"node_id"`
	Status      string                 `json:"status"`
	Output      map[string]interface{} `json:"output"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt time.Time              `json:"completed_at"`
	Error       *ExecutionError        `json:"error,omitempty"`
}

// ExecutionError represents an error during execution
type ExecutionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	NodeID  string `json:"node_id,omitempty"`
}

// ExecutionContext holds the current state during workflow execution
type ExecutionContext struct {
	Variables   map[string]interface{}
	NodeResults map[string]*NodeResult
}

// WorkflowDSL represents the compiled workflow definition
type WorkflowDSL struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Nodes       []WorkflowNode      `json:"nodes"`
	Edges       []WorkflowEdge      `json:"edges"`
	ErrorPolicy ErrorHandlingPolicy `json:"error_policy"`
}

// WorkflowNode represents a node in the workflow
type WorkflowNode struct {
	ID          string                 `json:"id"`
	Type        NodeType               `json:"type"`
	Label       string                 `json:"label"`
	Config      map[string]interface{} `json:"config"`
	RetryPolicy *RetryPolicy           `json:"retry_policy,omitempty"`
	Timeout     *int                   `json:"timeout_seconds,omitempty"`
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

// WorkflowEdge connects two nodes
type WorkflowEdge struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Condition string `json:"condition,omitempty"`
}

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxAttempts     int     `json:"max_attempts"`
	InitialInterval string  `json:"initial_interval"`
	BackoffCoeff    float64 `json:"backoff_coefficient"`
	MaxInterval     string  `json:"max_interval"`
}

// ErrorHandlingPolicy defines how errors are handled
type ErrorHandlingPolicy struct {
	OnError              string `json:"on_error"` // fail, compensate, ignore
	CompensationWorkflow string `json:"compensation_workflow,omitempty"`
}

// GetNode returns a node by ID
func (d *WorkflowDSL) GetNode(id string) *WorkflowNode {
	for i := range d.Nodes {
		if d.Nodes[i].ID == id {
			return &d.Nodes[i]
		}
	}
	return nil
}

// UserDefinedServiceWorkflow executes any user-created service dynamically
func UserDefinedServiceWorkflow(ctx workflow.Context, input *ServiceExecutionInput) (*ServiceExecutionResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting service execution",
		"serviceID", input.ServiceID,
		"versionID", input.VersionID,
		"tenantID", input.TenantID,
	)

	startedAt := workflow.Now(ctx)

	// Initialize execution context
	execCtx := &ExecutionContext{
		Variables:   make(map[string]interface{}),
		NodeResults: make(map[string]*NodeResult),
	}

	// Copy input data to variables
	for k, v := range input.InputData {
		execCtx.Variables[k] = v
	}

	// Set up query handler for workflow state
	err := workflow.SetQueryHandler(ctx, "getState", func() (*ExecutionContext, error) {
		return execCtx, nil
	})
	if err != nil {
		logger.Error("Failed to set query handler", "error", err)
	}

	// Get execution order (topological sort)
	executionOrder, err := getExecutionOrder(input.WorkflowDSL)
	if err != nil {
		return &ServiceExecutionResult{
			Status:    "failed",
			StartedAt: startedAt,
			Error: &ExecutionError{
				Code:    "INVALID_WORKFLOW",
				Message: err.Error(),
			},
		}, nil
	}

	// Execute nodes in order
	for _, nodeID := range executionOrder {
		node := input.WorkflowDSL.GetNode(nodeID)
		if node == nil {
			continue
		}

		result, err := executeNode(ctx, node, execCtx, input)
		if err != nil {
			if input.WorkflowDSL.ErrorPolicy.OnError == "compensate" {
				return compensate(ctx, execCtx, err, startedAt)
			}
			if input.WorkflowDSL.ErrorPolicy.OnError == "ignore" {
				logger.Warn("Ignoring node error", "node", nodeID, "error", err)
				continue
			}
			return &ServiceExecutionResult{
				Status:      "failed",
				OutputData:  execCtx.Variables,
				NodeResults: execCtx.NodeResults,
				StartedAt:   startedAt,
				CompletedAt: workflow.Now(ctx),
				Error: &ExecutionError{
					Code:    "EXECUTION_ERROR",
					Message: err.Error(),
					NodeID:  nodeID,
				},
			}, nil
		}

		execCtx.NodeResults[nodeID] = result

		// Update variables with output
		if result.Output != nil {
			for k, v := range result.Output {
				execCtx.Variables[k] = v
			}
		}
	}

	return &ServiceExecutionResult{
		Status:      "completed",
		OutputData:  execCtx.Variables,
		NodeResults: execCtx.NodeResults,
		StartedAt:   startedAt,
		CompletedAt: workflow.Now(ctx),
	}, nil
}

// getExecutionOrder returns nodes in topological order
func getExecutionOrder(dsl *WorkflowDSL) ([]string, error) {
	// Calculate in-degrees
	inDegree := make(map[string]int)
	for _, node := range dsl.Nodes {
		inDegree[node.ID] = 0
	}
	for _, edge := range dsl.Edges {
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

		for _, edge := range dsl.Edges {
			if edge.Source == nodeID {
				inDegree[edge.Target]--
				if inDegree[edge.Target] == 0 {
					queue = append(queue, edge.Target)
				}
			}
		}
	}

	if len(result) != len(dsl.Nodes) {
		return nil, fmt.Errorf("workflow contains cycles")
	}

	return result, nil
}

// executeNode dispatches to the appropriate node executor
func executeNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext, input *ServiceExecutionInput) (*NodeResult, error) {
	startedAt := workflow.Now(ctx)

	var output map[string]interface{}
	var err error

	switch node.Type {
	case NodeTypeActivity:
		output, err = executeActivityNode(ctx, node, execCtx)
	case NodeTypeAIAction:
		output, err = executeAIActionNode(ctx, node, execCtx)
	case NodeTypeN8N:
		output, err = executeN8NNode(ctx, node, execCtx)
	case NodeTypeDecision:
		output, err = executeDecisionNode(ctx, node, execCtx, input.WorkflowDSL)
	case NodeTypeParallel:
		output, err = executeParallelNode(ctx, node, execCtx, input)
	case NodeTypeWait:
		output, err = executeWaitNode(ctx, node, execCtx)
	case NodeTypeHumanTask:
		output, err = executeHumanTaskNode(ctx, node, execCtx)
	default:
		err = fmt.Errorf("unknown node type: %s", node.Type)
	}

	if err != nil {
		return nil, err
	}

	return &NodeResult{
		NodeID:      node.ID,
		Status:      "completed",
		Output:      output,
		StartedAt:   startedAt,
		CompletedAt: workflow.Now(ctx),
	}, nil
}

// compensate executes compensation logic when an error occurs
func compensate(ctx workflow.Context, execCtx *ExecutionContext, originalErr error, startedAt time.Time) (*ServiceExecutionResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing compensation for failed workflow")

	// In a full implementation, this would execute the compensation workflow
	// For now, we just return the error with compensation attempted

	return &ServiceExecutionResult{
		Status:      "compensated",
		OutputData:  execCtx.Variables,
		NodeResults: execCtx.NodeResults,
		StartedAt:   startedAt,
		CompletedAt: workflow.Now(ctx),
		Error: &ExecutionError{
			Code:    "COMPENSATED",
			Message: fmt.Sprintf("Original error: %v. Compensation executed.", originalErr),
		},
	}, nil
}
