package repository

import (
	"errors"
	"sync"

	"github.com/yesoreyeram/flow/backend/internal/models"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// Repository defines the interface for workflow storage
type Repository interface {
	// Workflow operations
	GetWorkflow(id string) (*models.Workflow, error)
	ListWorkflows() ([]*models.Workflow, error)
	CreateWorkflow(workflow *models.Workflow) error
	UpdateWorkflow(workflow *models.Workflow) error
	DeleteWorkflow(id string) error

	// Execution operations
	GetExecution(id string) (*models.WorkflowExecution, error)
	ListExecutions(workflowID string) ([]*models.WorkflowExecution, error)
	CreateExecution(execution *models.WorkflowExecution) error
	UpdateExecution(execution *models.WorkflowExecution) error
}

// InMemoryRepository implements Repository interface using in-memory storage
type InMemoryRepository struct {
	workflows  map[string]*models.Workflow
	executions map[string]*models.WorkflowExecution
	mu         sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		workflows:  make(map[string]*models.Workflow),
		executions: make(map[string]*models.WorkflowExecution),
	}
}

func (r *InMemoryRepository) GetWorkflow(id string) (*models.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflow, ok := r.workflows[id]
	if !ok {
		return nil, ErrNotFound
	}
	return workflow, nil
}

func (r *InMemoryRepository) ListWorkflows() ([]*models.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflows := make([]*models.Workflow, 0, len(r.workflows))
	for _, w := range r.workflows {
		workflows = append(workflows, w)
	}
	return workflows, nil
}

func (r *InMemoryRepository) CreateWorkflow(workflow *models.Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[workflow.ID]; exists {
		return ErrAlreadyExists
	}

	r.workflows[workflow.ID] = workflow
	return nil
}

func (r *InMemoryRepository) UpdateWorkflow(workflow *models.Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[workflow.ID]; !exists {
		return ErrNotFound
	}

	r.workflows[workflow.ID] = workflow
	return nil
}

func (r *InMemoryRepository) DeleteWorkflow(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[id]; !exists {
		return ErrNotFound
	}

	delete(r.workflows, id)
	return nil
}

func (r *InMemoryRepository) GetExecution(id string) (*models.WorkflowExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	execution, ok := r.executions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return execution, nil
}

func (r *InMemoryRepository) ListExecutions(workflowID string) ([]*models.WorkflowExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	executions := make([]*models.WorkflowExecution, 0)
	for _, e := range r.executions {
		if e.WorkflowID == workflowID {
			executions = append(executions, e)
		}
	}
	return executions, nil
}

func (r *InMemoryRepository) CreateExecution(execution *models.WorkflowExecution) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.executions[execution.ID]; exists {
		return ErrAlreadyExists
	}

	r.executions[execution.ID] = execution
	return nil
}

func (r *InMemoryRepository) UpdateExecution(execution *models.WorkflowExecution) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.executions[execution.ID]; !exists {
		return ErrNotFound
	}

	r.executions[execution.ID] = execution
	return nil
}
