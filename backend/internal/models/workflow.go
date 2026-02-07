package models

import (
	"time"
)

// NodeType represents the type of a workflow node
type NodeType string

const (
	NodeTypeHTTPRequest NodeType = "httpRequest"
	NodeTypeTransform   NodeType = "transform"
	NodeTypeCondition   NodeType = "condition"
	NodeTypeTrigger     NodeType = "trigger"
	NodeTypeWebhook     NodeType = "webhook"
)

// HTTPMethod represents HTTP request methods
type HTTPMethod string

const (
	MethodGET    HTTPMethod = "GET"
	MethodPOST   HTTPMethod = "POST"
	MethodPUT    HTTPMethod = "PUT"
	MethodDELETE HTTPMethod = "DELETE"
	MethodPATCH  HTTPMethod = "PATCH"
)

// Node represents a workflow node
type Node struct {
	ID       string                 `json:"id"`
	Type     NodeType               `json:"type"`
	Position Position               `json:"position"`
	Data     NodeData               `json:"data"`
}

// Position represents the position of a node on the canvas
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeData contains the node's data
type NodeData struct {
	Label   string                 `json:"label"`
	Config  map[string]interface{} `json:"config"`
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}

// Edge represents a connection between two nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// Workflow represents a complete workflow
type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Nodes       []Node    `json:"nodes"`
	Edges       []Edge    `json:"edges"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Version     int       `json:"version"`
}

// ExecutionStatus represents the status of an execution
type ExecutionStatus string

const (
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
)

// ExecutionResult represents the result of a node execution
type ExecutionResult struct {
	NodeID   string                 `json:"nodeId"`
	Status   string                 `json:"status"`
	Output   map[string]interface{} `json:"output,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Duration int64                  `json:"duration,omitempty"` // milliseconds
}

// WorkflowExecution represents a workflow execution instance
type WorkflowExecution struct {
	ID          string            `json:"id"`
	WorkflowID  string            `json:"workflowId"`
	Status      ExecutionStatus   `json:"status"`
	StartedAt   time.Time         `json:"startedAt"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
	Results     []ExecutionResult `json:"results"`
}

// HTTPRequestConfig represents configuration for HTTP request node
type HTTPRequestConfig struct {
	URL     string            `json:"url"`
	Method  HTTPMethod        `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
}

// TransformConfig represents configuration for transform node
type TransformConfig struct {
	Code     string `json:"code"`
	Language string `json:"language"` // "javascript" or "jq"
}

// ConditionConfig represents configuration for condition node
type ConditionConfig struct {
	Conditions []Condition `json:"conditions"`
	Combinator string      `json:"combinator"` // "AND" or "OR"
}

// Condition represents a single condition
type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "equals", "notEquals", "contains", "greaterThan", "lessThan"
	Value    interface{} `json:"value"`
}
