package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yesoreyeram/flow/backend/internal/models"
)

// Engine handles workflow execution
type Engine struct {
	httpClient *http.Client
}

// NewEngine creates a new workflow engine
func NewEngine() *Engine {
	return &Engine{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecuteWorkflow executes a workflow
func (e *Engine) ExecuteWorkflow(ctx context.Context, workflow *models.Workflow, input map[string]interface{}) (*models.WorkflowExecution, error) {
	execution := &models.WorkflowExecution{
		ID:         fmt.Sprintf("exec-%d", time.Now().UnixNano()),
		WorkflowID: workflow.ID,
		Status:     models.StatusRunning,
		StartedAt:  time.Now(),
		Results:    make([]models.ExecutionResult, 0),
	}

	// Build execution graph
	nodeOrder, err := e.buildExecutionOrder(workflow)
	if err != nil {
		execution.Status = models.StatusFailed
		now := time.Now()
		execution.CompletedAt = &now
		return execution, err
	}

	// Execute nodes in order
	nodeOutputs := make(map[string]map[string]interface{})
	nodeOutputs["input"] = input

	for _, nodeID := range nodeOrder {
		node := e.findNode(workflow, nodeID)
		if node == nil {
			continue
		}

		startTime := time.Now()
		output, err := e.executeNode(ctx, node, nodeOutputs)
		duration := time.Since(startTime).Milliseconds()

		result := models.ExecutionResult{
			NodeID:   nodeID,
			Duration: duration,
		}

		if err != nil {
			result.Status = "error"
			result.Error = err.Error()
			execution.Results = append(execution.Results, result)
			execution.Status = models.StatusFailed
			now := time.Now()
			execution.CompletedAt = &now
			return execution, err
		}

		result.Status = "success"
		result.Output = output
		execution.Results = append(execution.Results, result)
		nodeOutputs[nodeID] = output
	}

	execution.Status = models.StatusCompleted
	now := time.Now()
	execution.CompletedAt = &now

	return execution, nil
}

// buildExecutionOrder builds the execution order of nodes (topological sort)
func (e *Engine) buildExecutionOrder(workflow *models.Workflow) ([]string, error) {
	// Simple implementation: execute nodes in order they appear
	// In a production system, implement proper topological sort
	order := make([]string, 0, len(workflow.Nodes))
	for _, node := range workflow.Nodes {
		order = append(order, node.ID)
	}
	return order, nil
}

// findNode finds a node by ID
func (e *Engine) findNode(workflow *models.Workflow, nodeID string) *models.Node {
	for i := range workflow.Nodes {
		if workflow.Nodes[i].ID == nodeID {
			return &workflow.Nodes[i]
		}
	}
	return nil
}

// executeNode executes a single node
func (e *Engine) executeNode(ctx context.Context, node *models.Node, previousOutputs map[string]map[string]interface{}) (map[string]interface{}, error) {
	switch node.Type {
	case models.NodeTypeHTTPRequest:
		return e.executeHTTPRequest(ctx, node)
	case models.NodeTypeTransform:
		return e.executeTransform(ctx, node, previousOutputs)
	case models.NodeTypeCondition:
		return e.executeCondition(ctx, node, previousOutputs)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}

// executeHTTPRequest executes an HTTP request node
func (e *Engine) executeHTTPRequest(ctx context.Context, node *models.Node) (map[string]interface{}, error) {
	config := node.Data.Config

	url, ok := config["url"].(string)
	if !ok {
		return nil, errors.New("missing or invalid url")
	}

	method, ok := config["method"].(string)
	if !ok {
		method = "GET"
	}

	var body io.Reader
	if bodyData, ok := config["body"]; ok && bodyData != nil {
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// Execute request
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// If not JSON, return as string
		result = map[string]interface{}{
			"body":       string(respBody),
			"statusCode": resp.StatusCode,
		}
	} else {
		result["statusCode"] = resp.StatusCode
	}

	return result, nil
}

// executeTransform executes a transform node
func (e *Engine) executeTransform(ctx context.Context, node *models.Node, previousOutputs map[string]map[string]interface{}) (map[string]interface{}, error) {
	// Simple pass-through for now
	// In production, use a JavaScript runtime or jq library
	return map[string]interface{}{
		"transformed": true,
		"input":       previousOutputs,
	}, nil
}

// executeCondition executes a condition node
func (e *Engine) executeCondition(ctx context.Context, node *models.Node, previousOutputs map[string]map[string]interface{}) (map[string]interface{}, error) {
	// Simple condition evaluation
	// In production, implement full condition logic
	return map[string]interface{}{
		"result": true,
	}, nil
}
