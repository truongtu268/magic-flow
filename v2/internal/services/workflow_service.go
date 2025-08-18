package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/internal/engine"
	"magic-flow/v2/pkg/models"
)

// WorkflowService handles workflow business logic
type WorkflowService struct {
	repos  *database.RepositoryManager
	engine *engine.Engine
	parser *engine.WorkflowParser
	logger *logrus.Logger
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(repos *database.RepositoryManager, engine *engine.Engine, logger *logrus.Logger) *WorkflowService {
	return &WorkflowService{
		repos:  repos,
		engine: engine,
		parser: engine.NewWorkflowParser(),
		logger: logger,
	}
}

// CreateWorkflow creates a new workflow
func (s *WorkflowService) CreateWorkflow(req *CreateWorkflowRequest) (*models.Workflow, error) {
	// Parse workflow definition
	var workflow *models.Workflow
	var err error

	if req.YAMLDefinition != "" {
		workflow, err = s.parser.ParseYAML([]byte(req.YAMLDefinition))
	} else if req.JSONDefinition != "" {
		workflow, err = s.parser.ParseJSON([]byte(req.JSONDefinition))
	} else {
		return nil, fmt.Errorf("either YAML or JSON definition is required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow definition: %w", err)
	}

	// Override with request data
	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	workflow.CreatedBy = req.CreatedBy
	workflow.UpdatedBy = req.CreatedBy

	// Validate workflow
	if err := s.parser.ValidateWorkflow(workflow); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	// Save to database
	if err := s.repos.Workflow.Create(workflow); err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"workflow_id":   workflow.ID,
		"workflow_name": workflow.Name,
		"created_by":    workflow.CreatedBy,
	}).Info("Workflow created")

	return workflow, nil
}

// GetWorkflow retrieves a workflow by ID
func (s *WorkflowService) GetWorkflow(id uuid.UUID) (*models.Workflow, error) {
	workflow, err := s.repos.Workflow.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}
	return workflow, nil
}

// ListWorkflows retrieves workflows with pagination
func (s *WorkflowService) ListWorkflows(req *ListWorkflowsRequest) ([]*models.Workflow, int64, error) {
	workflows, total, err := s.repos.Workflow.List(req.Limit, req.Offset, req.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list workflows: %w", err)
	}
	return workflows, total, nil
}

// UpdateWorkflow updates an existing workflow
func (s *WorkflowService) UpdateWorkflow(id uuid.UUID, req *UpdateWorkflowRequest) (*models.Workflow, error) {
	// Get existing workflow
	workflow, err := s.repos.Workflow.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Update fields
	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if req.Status != "" {
		workflow.Status = models.WorkflowStatus(req.Status)
	}

	// Update definition if provided
	if req.YAMLDefinition != "" || req.JSONDefinition != "" {
		var updatedWorkflow *models.Workflow
		if req.YAMLDefinition != "" {
			updatedWorkflow, err = s.parser.ParseYAML([]byte(req.YAMLDefinition))
		} else {
			updatedWorkflow, err = s.parser.ParseJSON([]byte(req.JSONDefinition))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse workflow definition: %w", err)
		}

		// Validate updated workflow
		if err := s.parser.ValidateWorkflow(updatedWorkflow); err != nil {
			return nil, fmt.Errorf("workflow validation failed: %w", err)
		}

		workflow.Definition = updatedWorkflow.Definition
	}

	workflow.UpdatedBy = req.UpdatedBy
	workflow.UpdatedAt = time.Now().UTC()

	// Save changes
	if err := s.repos.Workflow.Update(workflow); err != nil {
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"workflow_id":   workflow.ID,
		"workflow_name": workflow.Name,
		"updated_by":    workflow.UpdatedBy,
	}).Info("Workflow updated")

	return workflow, nil
}

// DeleteWorkflow deletes a workflow
func (s *WorkflowService) DeleteWorkflow(id uuid.UUID) error {
	// Check if workflow exists
	_, err := s.repos.Workflow.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("workflow not found")
		}
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	// Check for active executions
	executions, err := s.repos.Execution.GetActiveExecutions()
	if err != nil {
		return fmt.Errorf("failed to check active executions: %w", err)
	}

	for _, execution := range executions {
		if execution.WorkflowID == id {
			return fmt.Errorf("cannot delete workflow with active executions")
		}
	}

	// Delete workflow
	if err := s.repos.Workflow.Delete(id); err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"workflow_id": id,
	}).Info("Workflow deleted")

	return nil
}

// ValidateWorkflow validates a workflow definition
func (s *WorkflowService) ValidateWorkflow(req *ValidateWorkflowRequest) (*ValidationResult, error) {
	var workflow *models.Workflow
	var err error

	if req.YAMLDefinition != "" {
		workflow, err = s.parser.ParseYAML([]byte(req.YAMLDefinition))
	} else if req.JSONDefinition != "" {
		workflow, err = s.parser.ParseJSON([]byte(req.JSONDefinition))
	} else {
		return nil, fmt.Errorf("either YAML or JSON definition is required")
	}

	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Parse error: %s", err.Error()))
		return result, nil
	}

	if err := s.parser.ValidateWorkflow(workflow); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Validation error: %s", err.Error()))
	}

	// Additional business logic validations
	if len(workflow.Definition.Spec.Steps) > 100 {
		result.Warnings = append(result.Warnings, "Workflow has more than 100 steps, consider breaking it down")
	}

	return result, nil
}

// ExecuteWorkflow executes a workflow
func (s *WorkflowService) ExecuteWorkflow(req *ExecuteWorkflowRequest) (*models.Execution, error) {
	// Get workflow
	workflow, err := s.repos.Workflow.GetByID(req.WorkflowID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	if workflow.Status != models.WorkflowStatusActive {
		return nil, fmt.Errorf("workflow is not active")
	}

	// Create execution record
	execution := &models.Execution{
		ID:          uuid.New(),
		WorkflowID:  req.WorkflowID,
		Status:      models.ExecutionStatusPending,
		TriggerType: models.TriggerType(req.TriggerType),
		TriggerData: req.TriggerData,
		Input:       req.Input,
		Context:     req.Context,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Save execution
	if err := s.repos.Execution.Create(execution); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Submit to engine for execution
	if err := s.engine.ExecuteWorkflow(workflow, execution); err != nil {
		// Update execution status to failed
		s.repos.Execution.UpdateStatus(execution.ID, models.ExecutionStatusFailed)
		return nil, fmt.Errorf("failed to execute workflow: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"execution_id": execution.ID,
		"workflow_id":  workflow.ID,
		"trigger_type": req.TriggerType,
		"created_by":   req.CreatedBy,
	}).Info("Workflow execution started")

	return execution, nil
}

// Request/Response types
type CreateWorkflowRequest struct {
	Name           string `json:"name" validate:"required,max=255"`
	Description    string `json:"description,omitempty"`
	YAMLDefinition string `json:"yaml_definition,omitempty"`
	JSONDefinition string `json:"json_definition,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
}

type UpdateWorkflowRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Status         string `json:"status,omitempty"`
	YAMLDefinition string `json:"yaml_definition,omitempty"`
	JSONDefinition string `json:"json_definition,omitempty"`
	UpdatedBy      string `json:"updated_by,omitempty"`
}

type ListWorkflowsRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Status string `json:"status,omitempty"`
}

type ValidateWorkflowRequest struct {
	YAMLDefinition string `json:"yaml_definition,omitempty"`
	JSONDefinition string `json:"json_definition,omitempty"`
}

type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type ExecuteWorkflowRequest struct {
	WorkflowID  uuid.UUID              `json:"workflow_id" validate:"required"`
	TriggerType string                 `json:"trigger_type" validate:"required"`
	TriggerData map[string]interface{} `json:"trigger_data,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	CreatedBy   string                 `json:"created_by,omitempty"`
}