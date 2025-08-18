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

// ExecutionService handles execution business logic
type ExecutionService struct {
	repos  *database.RepositoryManager
	engine *engine.Engine
	logger *logrus.Logger
}

// NewExecutionService creates a new execution service
func NewExecutionService(repos *database.RepositoryManager, engine *engine.Engine, logger *logrus.Logger) *ExecutionService {
	return &ExecutionService{
		repos:  repos,
		engine: engine,
		logger: logger,
	}
}

// GetExecution retrieves an execution by ID
func (s *ExecutionService) GetExecution(id uuid.UUID) (*models.Execution, error) {
	execution, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}
	return execution, nil
}

// ListExecutions retrieves executions with pagination and filtering
func (s *ExecutionService) ListExecutions(req *ListExecutionsRequest) ([]*models.Execution, int64, error) {
	executions, total, err := s.repos.Execution.List(req.Limit, req.Offset, req.WorkflowID, req.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list executions: %w", err)
	}
	return executions, total, nil
}

// GetExecutionStatus retrieves the current status of an execution
func (s *ExecutionService) GetExecutionStatus(id uuid.UUID) (*ExecutionStatusResponse, error) {
	execution, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	// Get step executions
	stepExecutions, err := s.repos.StepExecution.GetByExecutionID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get step executions: %w", err)
	}

	// Calculate progress
	totalSteps := len(stepExecutions)
	completedSteps := 0
	for _, step := range stepExecutions {
		if step.Status == models.StepExecutionStatusCompleted {
			completedSteps++
		}
	}

	progress := 0.0
	if totalSteps > 0 {
		progress = float64(completedSteps) / float64(totalSteps) * 100
	}

	return &ExecutionStatusResponse{
		ID:             execution.ID,
		WorkflowID:     execution.WorkflowID,
		Status:         string(execution.Status),
		Progress:       progress,
		TotalSteps:     totalSteps,
		CompletedSteps: completedSteps,
		StartedAt:      execution.StartedAt,
		CompletedAt:    execution.CompletedAt,
		Error:          execution.Error,
		StepExecutions: stepExecutions,
	}, nil
}

// GetExecutionResults retrieves the results of an execution
func (s *ExecutionService) GetExecutionResults(id uuid.UUID) (*ExecutionResultsResponse, error) {
	execution, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	if execution.Status != models.ExecutionStatusCompleted && execution.Status != models.ExecutionStatusFailed {
		return nil, fmt.Errorf("execution is not completed")
	}

	// Get step executions with results
	stepExecutions, err := s.repos.StepExecution.GetByExecutionID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get step executions: %w", err)
	}

	return &ExecutionResultsResponse{
		ID:             execution.ID,
		WorkflowID:     execution.WorkflowID,
		Status:         string(execution.Status),
		Output:         execution.Output,
		Error:          execution.Error,
		StartedAt:      execution.StartedAt,
		CompletedAt:    execution.CompletedAt,
		Duration:       s.calculateDuration(execution.StartedAt, execution.CompletedAt),
		StepResults:    s.buildStepResults(stepExecutions),
	}, nil
}

// GetExecutionLogs retrieves logs for an execution
func (s *ExecutionService) GetExecutionLogs(id uuid.UUID, req *GetExecutionLogsRequest) (*ExecutionLogsResponse, error) {
	// Check if execution exists
	_, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	// Get execution events (which serve as logs)
	events, total, err := s.repos.ExecutionEvent.GetByExecutionID(id, req.Limit, req.Offset, req.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution events: %w", err)
	}

	// Convert events to log entries
	logs := make([]LogEntry, len(events))
	for i, event := range events {
		logs[i] = LogEntry{
			Timestamp: event.Timestamp,
			Level:     string(event.Type),
			Message:   event.Message,
			StepID:    event.StepID,
			Data:      event.Data,
		}
	}

	return &ExecutionLogsResponse{
		ExecutionID: id,
		Logs:        logs,
		Total:       total,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}, nil
}

// GetExecutionEvents retrieves events for an execution
func (s *ExecutionService) GetExecutionEvents(id uuid.UUID, req *GetExecutionEventsRequest) ([]*models.ExecutionEvent, int64, error) {
	// Check if execution exists
	_, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("execution not found")
		}
		return nil, 0, fmt.Errorf("failed to get execution: %w", err)
	}

	events, total, err := s.repos.ExecutionEvent.GetByExecutionID(id, req.Limit, req.Offset, req.EventType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get execution events: %w", err)
	}

	return events, total, nil
}

// CancelExecution cancels a running execution
func (s *ExecutionService) CancelExecution(id uuid.UUID, cancelledBy string) error {
	execution, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("execution not found")
		}
		return fmt.Errorf("failed to get execution: %w", err)
	}

	if execution.Status != models.ExecutionStatusRunning && execution.Status != models.ExecutionStatusPending {
		return fmt.Errorf("execution cannot be cancelled in current status: %s", execution.Status)
	}

	// Cancel in engine
	if err := s.engine.CancelExecution(id); err != nil {
		s.logger.WithError(err).WithField("execution_id", id).Warn("Failed to cancel execution in engine")
	}

	// Update status in database
	if err := s.repos.Execution.UpdateStatus(id, models.ExecutionStatusCancelled); err != nil {
		return fmt.Errorf("failed to update execution status: %w", err)
	}

	// Create cancellation event
	event := &models.ExecutionEvent{
		ID:          uuid.New(),
		ExecutionID: id,
		Type:        models.EventTypeExecutionCancelled,
		Message:     fmt.Sprintf("Execution cancelled by %s", cancelledBy),
		Timestamp:   time.Now().UTC(),
		Data: map[string]interface{}{
			"cancelled_by": cancelledBy,
		},
	}

	if err := s.repos.ExecutionEvent.Create(event); err != nil {
		s.logger.WithError(err).Warn("Failed to create cancellation event")
	}

	s.logger.WithFields(logrus.Fields{
		"execution_id": id,
		"cancelled_by": cancelledBy,
	}).Info("Execution cancelled")

	return nil
}

// RetryExecution retries a failed execution
func (s *ExecutionService) RetryExecution(id uuid.UUID, retryBy string) (*models.Execution, error) {
	originalExecution, err := s.repos.Execution.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	if originalExecution.Status != models.ExecutionStatusFailed {
		return nil, fmt.Errorf("only failed executions can be retried")
	}

	// Get workflow
	workflow, err := s.repos.Workflow.GetByID(originalExecution.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	if workflow.Status != models.WorkflowStatusActive {
		return nil, fmt.Errorf("workflow is not active")
	}

	// Create new execution
	newExecution := &models.Execution{
		ID:               uuid.New(),
		WorkflowID:       originalExecution.WorkflowID,
		Status:           models.ExecutionStatusPending,
		TriggerType:      originalExecution.TriggerType,
		TriggerData:      originalExecution.TriggerData,
		Input:            originalExecution.Input,
		Context:          originalExecution.Context,
		ParentExecutionID: &originalExecution.ID,
		CreatedBy:        retryBy,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Save new execution
	if err := s.repos.Execution.Create(newExecution); err != nil {
		return nil, fmt.Errorf("failed to create retry execution: %w", err)
	}

	// Submit to engine for execution
	if err := s.engine.ExecuteWorkflow(workflow, newExecution); err != nil {
		// Update execution status to failed
		s.repos.Execution.UpdateStatus(newExecution.ID, models.ExecutionStatusFailed)
		return nil, fmt.Errorf("failed to execute retry workflow: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"original_execution_id": originalExecution.ID,
		"new_execution_id":      newExecution.ID,
		"workflow_id":           workflow.ID,
		"retry_by":              retryBy,
	}).Info("Execution retried")

	return newExecution, nil
}

// Helper methods
func (s *ExecutionService) calculateDuration(startedAt, completedAt *time.Time) *time.Duration {
	if startedAt == nil || completedAt == nil {
		return nil
	}
	duration := completedAt.Sub(*startedAt)
	return &duration
}

func (s *ExecutionService) buildStepResults(stepExecutions []*models.StepExecution) []StepResult {
	results := make([]StepResult, len(stepExecutions))
	for i, step := range stepExecutions {
		results[i] = StepResult{
			StepID:      step.StepID,
			Status:      string(step.Status),
			Output:      step.Output,
			Error:       step.Error,
			StartedAt:   step.StartedAt,
			CompletedAt: step.CompletedAt,
			Duration:    s.calculateDuration(step.StartedAt, step.CompletedAt),
		}
	}
	return results
}

// Request/Response types
type ListExecutionsRequest struct {
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	Status     string     `json:"status,omitempty"`
}

type ExecutionStatusResponse struct {
	ID             uuid.UUID                `json:"id"`
	WorkflowID     uuid.UUID                `json:"workflow_id"`
	Status         string                   `json:"status"`
	Progress       float64                  `json:"progress"`
	TotalSteps     int                      `json:"total_steps"`
	CompletedSteps int                      `json:"completed_steps"`
	StartedAt      *time.Time               `json:"started_at"`
	CompletedAt    *time.Time               `json:"completed_at"`
	Error          *string                  `json:"error"`
	StepExecutions []*models.StepExecution  `json:"step_executions"`
}

type ExecutionResultsResponse struct {
	ID          uuid.UUID               `json:"id"`
	WorkflowID  uuid.UUID               `json:"workflow_id"`
	Status      string                  `json:"status"`
	Output      map[string]interface{}  `json:"output"`
	Error       *string                 `json:"error"`
	StartedAt   *time.Time              `json:"started_at"`
	CompletedAt *time.Time              `json:"completed_at"`
	Duration    *time.Duration          `json:"duration"`
	StepResults []StepResult            `json:"step_results"`
}

type StepResult struct {
	StepID      string                 `json:"step_id"`
	Status      string                 `json:"status"`
	Output      map[string]interface{} `json:"output"`
	Error       *string                `json:"error"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Duration    *time.Duration         `json:"duration"`
}

type GetExecutionLogsRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Level  string `json:"level,omitempty"`
}

type ExecutionLogsResponse struct {
	ExecutionID uuid.UUID  `json:"execution_id"`
	Logs        []LogEntry `json:"logs"`
	Total       int64      `json:"total"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	StepID    *string                `json:"step_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type GetExecutionEventsRequest struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	EventType string `json:"event_type,omitempty"`
}