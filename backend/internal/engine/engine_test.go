package engine

import (
	"context"
	"testing"
	"time"

	"github.com/yesoreyeram/flow/backend/internal/models"
)

func TestEngine(t *testing.T) {
	engine := NewEngine()

	t.Run("ExecuteSimpleWorkflow", func(t *testing.T) {
		workflow := &models.Workflow{
			ID:   "test-workflow",
			Name: "Test Workflow",
			Nodes: []models.Node{
				{
					ID:   "node-1",
					Type: models.NodeTypeTransform,
					Data: models.NodeData{
						Label: "Transform",
						Config: map[string]interface{}{
							"code":     "return input",
							"language": "javascript",
						},
					},
				},
			},
			Edges:     []models.Edge{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		input := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		execution, err := engine.ExecuteWorkflow(ctx, workflow, input)

		if err != nil {
			t.Fatalf("Workflow execution failed: %v", err)
		}

		if execution.Status != models.StatusCompleted {
			t.Errorf("Expected status completed, got %s", execution.Status)
		}

		if len(execution.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(execution.Results))
		}
	})

	t.Run("BuildExecutionOrder", func(t *testing.T) {
		workflow := &models.Workflow{
			Nodes: []models.Node{
				{ID: "node-1"},
				{ID: "node-2"},
				{ID: "node-3"},
			},
		}

		order, err := engine.buildExecutionOrder(workflow)
		if err != nil {
			t.Fatalf("Failed to build execution order: %v", err)
		}

		if len(order) != 3 {
			t.Errorf("Expected 3 nodes in order, got %d", len(order))
		}
	})
}
