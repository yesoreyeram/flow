package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yesoreyeram/flow/backend/internal/config"
	"github.com/yesoreyeram/flow/backend/internal/engine"
	"github.com/yesoreyeram/flow/backend/internal/models"
	"github.com/yesoreyeram/flow/backend/internal/repository"
)

// Server represents the API server
type Server struct {
	config *config.Config
	repo   repository.Repository
	engine *engine.Engine
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, repo repository.Repository) http.Handler {
	s := &Server{
		config: cfg,
		repo:   repo,
		engine: engine.NewEngine(),
	}

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/api/health", s.handleHealth)

	// Workflow endpoints
	mux.HandleFunc("/api/workflows", s.handleWorkflows)
	mux.HandleFunc("/api/workflows/", s.handleWorkflowByID)

	// Execution endpoints
	mux.HandleFunc("/api/executions/", s.handleExecutionByID)

	// Apply middleware
	handler := corsMiddleware(cfg.CorsOrigins)(mux)
	handler = loggingMiddleware(handler)
	handler = recoveryMiddleware(handler)

	return handler
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleWorkflows(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listWorkflows(w, r)
	case http.MethodPost:
		s.createWorkflow(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleWorkflowByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := r.URL.Path[len("/api/workflows/"):]
	
	// Check for execute endpoint
	if len(id) > 8 && id[len(id)-8:] == "/execute" {
		workflowID := id[:len(id)-8]
		s.executeWorkflow(w, r, workflowID)
		return
	}
	
	// Check for executions endpoint
	if len(id) > 11 && id[len(id)-11:] == "/executions" {
		workflowID := id[:len(id)-11]
		s.listWorkflowExecutions(w, r, workflowID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getWorkflow(w, r, id)
	case http.MethodPut:
		s.updateWorkflow(w, r, id)
	case http.MethodDelete:
		s.deleteWorkflow(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := s.repo.ListWorkflows()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list workflows")
		return
	}

	respondJSON(w, http.StatusOK, workflows)
}

func (s *Server) getWorkflow(w http.ResponseWriter, r *http.Request, id string) {
	workflow, err := s.repo.GetWorkflow(id)
	if err == repository.ErrNotFound {
		respondError(w, http.StatusNotFound, "Workflow not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get workflow")
		return
	}

	respondJSON(w, http.StatusOK, workflow)
}

func (s *Server) createWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set metadata
	now := time.Now()
	workflow.ID = fmt.Sprintf("wf-%d", now.UnixNano())
	workflow.CreatedAt = now
	workflow.UpdatedAt = now
	workflow.Version = 1

	if err := s.repo.CreateWorkflow(&workflow); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create workflow")
		return
	}

	respondJSON(w, http.StatusCreated, workflow)
}

func (s *Server) updateWorkflow(w http.ResponseWriter, r *http.Request, id string) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	workflow.ID = id
	workflow.UpdatedAt = time.Now()
	workflow.Version++

	if err := s.repo.UpdateWorkflow(&workflow); err == repository.ErrNotFound {
		respondError(w, http.StatusNotFound, "Workflow not found")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update workflow")
		return
	}

	respondJSON(w, http.StatusOK, workflow)
}

func (s *Server) deleteWorkflow(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.repo.DeleteWorkflow(id); err == repository.ErrNotFound {
		respondError(w, http.StatusNotFound, "Workflow not found")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete workflow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) executeWorkflow(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workflow, err := s.repo.GetWorkflow(id)
	if err == repository.ErrNotFound {
		respondError(w, http.StatusNotFound, "Workflow not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get workflow")
		return
	}

	var input struct {
		Input map[string]interface{} `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		input.Input = make(map[string]interface{})
	}

	execution, err := s.engine.ExecuteWorkflow(r.Context(), workflow, input.Input)
	if err != nil {
		// Execution still returns even if failed
		if execution != nil {
			_ = s.repo.CreateExecution(execution)
		}
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Workflow execution failed: %v", err))
		return
	}

	if err := s.repo.CreateExecution(execution); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save execution")
		return
	}

	respondJSON(w, http.StatusOK, execution)
}

func (s *Server) handleExecutionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/executions/"):]
	execution, err := s.repo.GetExecution(id)
	if err == repository.ErrNotFound {
		respondError(w, http.StatusNotFound, "Execution not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get execution")
		return
	}

	respondJSON(w, http.StatusOK, execution)
}

func (s *Server) listWorkflowExecutions(w http.ResponseWriter, r *http.Request, workflowID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	executions, err := s.repo.ListExecutions(workflowID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list executions")
		return
	}

	respondJSON(w, http.StatusOK, executions)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"message": message,
		"code":    fmt.Sprintf("%d", status),
	})
}
