package repository

import (
	"testing"
	"time"

	"github.com/yesoreyeram/flow/backend/internal/models"
)

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()

	t.Run("CreateAndGetWorkflow", func(t *testing.T) {
		workflow := &models.Workflow{
			ID:        "test-1",
			Name:      "Test Workflow",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Nodes:     []models.Node{},
			Edges:     []models.Edge{},
		}

		err := repo.CreateWorkflow(workflow)
		if err != nil {
			t.Fatalf("Failed to create workflow: %v", err)
		}

		retrieved, err := repo.GetWorkflow("test-1")
		if err != nil {
			t.Fatalf("Failed to get workflow: %v", err)
		}

		if retrieved.ID != workflow.ID {
			t.Errorf("Expected ID %s, got %s", workflow.ID, retrieved.ID)
		}
	})

	t.Run("UpdateWorkflow", func(t *testing.T) {
		workflow := &models.Workflow{
			ID:        "test-2",
			Name:      "Original Name",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Nodes:     []models.Node{},
			Edges:     []models.Edge{},
		}

		repo.CreateWorkflow(workflow)

		workflow.Name = "Updated Name"
		err := repo.UpdateWorkflow(workflow)
		if err != nil {
			t.Fatalf("Failed to update workflow: %v", err)
		}

		retrieved, _ := repo.GetWorkflow("test-2")
		if retrieved.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got '%s'", retrieved.Name)
		}
	})

	t.Run("DeleteWorkflow", func(t *testing.T) {
		workflow := &models.Workflow{
			ID:        "test-3",
			Name:      "To Delete",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Nodes:     []models.Node{},
			Edges:     []models.Edge{},
		}

		repo.CreateWorkflow(workflow)

		err := repo.DeleteWorkflow("test-3")
		if err != nil {
			t.Fatalf("Failed to delete workflow: %v", err)
		}

		_, err = repo.GetWorkflow("test-3")
		if err != ErrNotFound {
			t.Error("Expected ErrNotFound after deletion")
		}
	})

	t.Run("ListWorkflows", func(t *testing.T) {
		workflows, err := repo.ListWorkflows()
		if err != nil {
			t.Fatalf("Failed to list workflows: %v", err)
		}

		if len(workflows) < 1 {
			t.Error("Expected at least one workflow")
		}
	})
}
