package models

import (
	"encoding/json"
	"testing"
)

func TestWorkflowSerialization(t *testing.T) {
	workflow := Workflow{
		ID:          "test-1",
		Name:        "Test Workflow",
		Description: "A test workflow",
		Nodes: []Node{
			{
				ID:   "node-1",
				Type: NodeTypeHTTPRequest,
				Position: Position{
					X: 100,
					Y: 200,
				},
				Data: NodeData{
					Label: "HTTP Request",
					Config: map[string]interface{}{
						"url":    "https://api.example.com",
						"method": "GET",
					},
				},
			},
		},
		Edges: []Edge{},
	}

	// Test JSON serialization
	data, err := json.Marshal(workflow)
	if err != nil {
		t.Fatalf("Failed to marshal workflow: %v", err)
	}

	// Test JSON deserialization
	var decoded Workflow
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal workflow: %v", err)
	}

	if decoded.ID != workflow.ID {
		t.Errorf("Expected ID %s, got %s", workflow.ID, decoded.ID)
	}

	if len(decoded.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(decoded.Nodes))
	}
}

func TestNodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
	}{
		{"HTTP Request", NodeTypeHTTPRequest},
		{"Transform", NodeTypeTransform},
		{"Condition", NodeTypeCondition},
		{"Trigger", NodeTypeTrigger},
		{"Webhook", NodeTypeWebhook},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nodeType == "" {
				t.Error("Node type should not be empty")
			}
		})
	}
}
