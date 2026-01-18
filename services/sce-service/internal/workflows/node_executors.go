// Package workflows contains node executors for the SCE service.
package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// executeActivityNode executes a Temporal activity
func executeActivityNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)

	activityName, ok := node.Config["activity_name"].(string)
	if !ok {
		return nil, fmt.Errorf("activity node %s missing activity_name", node.ID)
	}

	logger.Info("Executing activity node", "node", node.ID, "activity", activityName)

	// Build activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}

	if node.Timeout != nil {
		ao.StartToCloseTimeout = time.Duration(*node.Timeout) * time.Second
	}

	if node.RetryPolicy != nil {
		ao.RetryPolicy = &workflow.RetryPolicy{
			MaximumAttempts: int32(node.RetryPolicy.MaxAttempts),
		}
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// Build input from mapping
	input := buildActivityInput(node.Config, execCtx.Variables)

	// Execute activity
	var result map[string]interface{}
	err := workflow.ExecuteActivity(ctx, activityName, input).Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("activity %s failed: %w", activityName, err)
	}

	return result, nil
}

// executeAIActionNode executes an AI/LLM action
func executeAIActionNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing AI action node", "node", node.ID)

	provider := getStringConfig(node.Config, "provider", "anthropic")
	model := getStringConfig(node.Config, "model", "claude-3-sonnet")
	promptTemplate := getStringConfig(node.Config, "prompt_template", "")
	useLocal := getBoolConfig(node.Config, "use_local_model", false)

	// Interpolate prompt with variables
	prompt := interpolateTemplate(promptTemplate, execCtx.Variables)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 120 * time.Second,
		TaskQueue:           "omniroute-ai",
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	input := map[string]interface{}{
		"provider":        provider,
		"model":           model,
		"prompt":          prompt,
		"use_local_model": useLocal,
		"max_tokens":      getIntConfig(node.Config, "max_tokens", 1000),
		"temperature":     getFloatConfig(node.Config, "temperature", 0.7),
	}

	var result map[string]interface{}
	activityName := "CallLLM"
	if useLocal {
		activityName = "LocalModelInference"
	}

	err := workflow.ExecuteActivity(ctx, activityName, input).Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("AI action failed: %w", err)
	}

	return result, nil
}

// executeN8NNode executes an n8n workflow
func executeN8NNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing n8n node", "node", node.ID)

	workflowID := getStringConfig(node.Config, "n8n_workflow_id", "")
	if workflowID == "" {
		return nil, fmt.Errorf("n8n node %s missing n8n_workflow_id", node.ID)
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 300 * time.Second,
		TaskQueue:           "omniroute-integration",
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	input := map[string]interface{}{
		"workflow_id":         workflowID,
		"input_data":          buildActivityInput(node.Config, execCtx.Variables),
		"wait_for_completion": true,
	}

	var result map[string]interface{}
	err := workflow.ExecuteActivity(ctx, "ExecuteN8NWorkflow", input).Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("n8n workflow failed: %w", err)
	}

	return result, nil
}

// executeDecisionNode evaluates a condition and returns the branch result
func executeDecisionNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext, dsl *WorkflowDSL) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing decision node", "node", node.ID)

	condition := getStringConfig(node.Config, "condition", "")

	// Evaluate condition (simplified - in production use expr library)
	result := evaluateCondition(condition, execCtx.Variables)

	return map[string]interface{}{
		"branch":    result,
		"condition": condition,
	}, nil
}

// executeParallelNode executes multiple branches in parallel
func executeParallelNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext, input *ServiceExecutionInput) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing parallel node", "node", node.ID)

	branches, ok := node.Config["branches"].([]interface{})
	if !ok || len(branches) == 0 {
		return nil, fmt.Errorf("parallel node %s has no branches", node.ID)
	}

	// Create selector for parallel execution
	selector := workflow.NewSelector(ctx)
	results := make(map[string]interface{})

	for i, branch := range branches {
		branchID, ok := branch.(string)
		if !ok {
			continue
		}

		branchNode := input.WorkflowDSL.GetNode(branchID)
		if branchNode == nil {
			continue
		}

		// Execute each branch as a goroutine
		branchCtx := workflow.WithValue(ctx, "branch", i)
		future := workflow.ExecuteActivity(branchCtx, "ExecuteBranch", map[string]interface{}{
			"node_id": branchID,
			"config":  branchNode.Config,
		})

		selector.AddFuture(future, func(f workflow.Future) {
			var branchResult map[string]interface{}
			if err := f.Get(ctx, &branchResult); err != nil {
				logger.Error("Branch failed", "branch", branchID, "error", err)
			} else {
				results[branchID] = branchResult
			}
		})
	}

	// Wait for all branches
	for range branches {
		selector.Select(ctx)
	}

	return results, nil
}

// executeWaitNode pauses execution for a specified duration or signal
func executeWaitNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing wait node", "node", node.ID)

	waitType := getStringConfig(node.Config, "wait_type", "timer")

	switch waitType {
	case "timer":
		durationSec := getIntConfig(node.Config, "duration_seconds", 60)
		err := workflow.Sleep(ctx, time.Duration(durationSec)*time.Second)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"waited": durationSec}, nil

	case "signal":
		signalName := getStringConfig(node.Config, "signal_name", "continue")
		signalChan := workflow.GetSignalChannel(ctx, signalName)

		var signalData interface{}
		signalChan.Receive(ctx, &signalData)

		return map[string]interface{}{"signal": signalName, "data": signalData}, nil

	default:
		return nil, fmt.Errorf("unknown wait type: %s", waitType)
	}
}

// executeHumanTaskNode waits for human approval
func executeHumanTaskNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Executing human task node", "node", node.ID)

	taskType := getStringConfig(node.Config, "task_type", "approval")
	title := getStringConfig(node.Config, "title", "Task")
	timeoutSec := getIntConfig(node.Config, "timeout_seconds", 86400) // 24 hours

	// Create human task signal channel
	signalName := fmt.Sprintf("human_task_%s", node.ID)
	signalChan := workflow.GetSignalChannel(ctx, signalName)

	// Wait for human response with timeout
	var response map[string]interface{}
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &response)
	})

	timerFuture := workflow.NewTimer(ctx, time.Duration(timeoutSec)*time.Second)
	selector.AddFuture(timerFuture, func(f workflow.Future) {
		response = map[string]interface{}{
			"timeout": true,
			"message": "Human task timed out",
		}
	})

	selector.Select(ctx)

	return map[string]interface{}{
		"task_type": taskType,
		"title":     title,
		"response":  response,
	}, nil
}

// Helper functions

func buildActivityInput(config map[string]interface{}, variables map[string]interface{}) map[string]interface{} {
	input := make(map[string]interface{})

	if mapping, ok := config["input_mapping"].(map[string]interface{}); ok {
		for key, varName := range mapping {
			if name, ok := varName.(string); ok {
				if val, exists := variables[name]; exists {
					input[key] = val
				}
			}
		}
	}

	return input
}

func interpolateTemplate(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = replaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

func replaceAll(s, old, new string) string {
	for {
		idx := findIndex(s, old)
		if idx == -1 {
			break
		}
		s = s[:idx] + new + s[idx+len(old):]
	}
	return s
}

func findIndex(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func evaluateCondition(condition string, variables map[string]interface{}) bool {
	// Simplified condition evaluation
	// In production, use an expression library like antonmedv/expr
	if val, ok := variables[condition]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func getStringConfig(config map[string]interface{}, key, defaultVal string) string {
	if val, ok := config[key].(string); ok {
		return val
	}
	return defaultVal
}

func getIntConfig(config map[string]interface{}, key string, defaultVal int) int {
	if val, ok := config[key].(float64); ok {
		return int(val)
	}
	if val, ok := config[key].(int); ok {
		return val
	}
	return defaultVal
}

func getFloatConfig(config map[string]interface{}, key string, defaultVal float64) float64 {
	if val, ok := config[key].(float64); ok {
		return val
	}
	return defaultVal
}

func getBoolConfig(config map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := config[key].(bool); ok {
		return val
	}
	return defaultVal
}
