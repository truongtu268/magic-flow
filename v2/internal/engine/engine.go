package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"magic-flow/v2/pkg/models"
)

// Engine represents the workflow execution engine
type Engine struct {
	mu               sync.RWMutex
	executions       map[uuid.UUID]*ExecutionContext
	stepExecutors    map[string]StepExecutor
	eventHandlers    []EventHandler
	metrics          MetricsCollector
	logger           *logrus.Logger
	maxConcurrent    int
	currentExecutions int
	shutdownCh       chan struct{}
	wg               sync.WaitGroup
}

// ExecutionContext holds the context for a workflow execution
type ExecutionContext struct {
	Execution    *models.Execution
	Workflow     *models.Workflow
	Input        map[string]interface{}
	Output       map[string]interface{}
	Variables    map[string]interface{}
	StepResults  map[string]interface{}
	Context      context.Context
	Cancel       context.CancelFunc
	StartTime    time.Time
	EndTime      *time.Time
	CurrentStep  string
	RetryCount   int
	MaxRetries   int
	Timeout      time.Duration
	mu           sync.RWMutex
}

// StepExecutor interface for executing workflow steps
type StepExecutor interface {
	Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error)
	Validate(step *models.WorkflowStep) error
	GetType() string
}

// EventHandler interface for handling workflow events
type EventHandler interface {
	Handle(event *WorkflowEvent) error
	GetEventTypes() []string
}

// MetricsCollector interface for collecting execution metrics
type MetricsCollector interface {
	RecordExecution(execution *models.Execution)
	RecordStepExecution(step *models.StepExecution)
	RecordError(err error, context map[string]interface{})
	RecordMetric(name string, value float64, labels map[string]string)
}

// WorkflowEvent represents an event during workflow execution
type WorkflowEvent struct {
	Type        string                 `json:"type"`
	ExecutionID uuid.UUID              `json:"execution_id"`
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	StepID      string                 `json:"step_id,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`
	Error       string                 `json:"error,omitempty"`
}

// NewEngine creates a new workflow execution engine
func NewEngine(maxConcurrent int, metrics MetricsCollector, logger *logrus.Logger) *Engine {
	return &Engine{
		executions:    make(map[uuid.UUID]*ExecutionContext),
		stepExecutors: make(map[string]StepExecutor),
		eventHandlers: make([]EventHandler, 0),
		metrics:       metrics,
		logger:        logger,
		maxConcurrent: maxConcurrent,
		shutdownCh:    make(chan struct{}),
	}
}

// RegisterStepExecutor registers a step executor for a specific step type
func (e *Engine) RegisterStepExecutor(stepType string, executor StepExecutor) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stepExecutors[stepType] = executor
}

// RegisterEventHandler registers an event handler
func (e *Engine) RegisterEventHandler(handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventHandlers = append(e.eventHandlers, handler)
}

// ExecuteWorkflow executes a workflow
func (e *Engine) ExecuteWorkflow(ctx context.Context, workflow *models.Workflow, input map[string]interface{}, config map[string]interface{}) (*models.Execution, error) {
	// Check if we can accept more executions
	e.mu.Lock()
	if e.currentExecutions >= e.maxConcurrent {
		e.mu.Unlock()
		return nil, fmt.Errorf("maximum concurrent executions reached: %d", e.maxConcurrent)
	}
	e.currentExecutions++
	e.mu.Unlock()

	// Create execution record
	execution := &models.Execution{
		ID:         uuid.New(),
		WorkflowID: workflow.ID,
		Status:     models.ExecutionStatusRunning,
		Input:      input,
		Config:     config,
		StartedAt:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	// Create execution context
	execCtx, cancel := context.WithCancel(ctx)
	execContext := &ExecutionContext{
		Execution:   execution,
		Workflow:    workflow,
		Input:       input,
		Output:      make(map[string]interface{}),
		Variables:   make(map[string]interface{}),
		StepResults: make(map[string]interface{}),
		Context:     execCtx,
		Cancel:      cancel,
		StartTime:   time.Now(),
		MaxRetries:  3, // Default retry count
		Timeout:     30 * time.Minute, // Default timeout
	}

	// Apply configuration
	if timeout, ok := config["timeout"]; ok {
		if timeoutInt, ok := timeout.(int); ok {
			execContext.Timeout = time.Duration(timeoutInt) * time.Second
		}
	}
	if maxRetries, ok := config["max_retries"]; ok {
		if retriesInt, ok := maxRetries.(int); ok {
			execContext.MaxRetries = retriesInt
		}
	}

	// Set timeout if specified
	if execContext.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(execCtx, execContext.Timeout)
		execContext.Context = execCtx
		execContext.Cancel = cancel
	}

	// Store execution context
	e.mu.Lock()
	e.executions[execution.ID] = execContext
	e.mu.Unlock()

	// Start execution in goroutine
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		defer func() {
			e.mu.Lock()
			e.currentExecutions--
			delete(e.executions, execution.ID)
			e.mu.Unlock()
		}()

		e.executeWorkflowSteps(execContext)
	}()

	// Emit execution started event
	e.emitEvent(&WorkflowEvent{
		Type:        "execution.started",
		ExecutionID: execution.ID,
		WorkflowID:  workflow.ID,
		Timestamp:   time.Now().UTC(),
		Data: map[string]interface{}{
			"input":  input,
			"config": config,
		},
	})

	// Record metrics
	e.metrics.RecordExecution(execution)

	e.logger.WithFields(logrus.Fields{
		"execution_id": execution.ID,
		"workflow_id":  workflow.ID,
		"workflow_name": workflow.Name,
	}).Info("Workflow execution started")

	return execution, nil
}

// executeWorkflowSteps executes the workflow steps
func (e *Engine) executeWorkflowSteps(execContext *ExecutionContext) {
	defer execContext.Cancel()

	// Parse workflow definition
	var workflowDef models.WorkflowDefinition
	if err := yaml.Unmarshal([]byte(fmt.Sprintf("%v", execContext.Workflow.Definition)), &workflowDef); err != nil {
		e.failExecution(execContext, fmt.Errorf("failed to parse workflow definition: %w", err))
		return
	}

	// Initialize variables with input
	for key, value := range execContext.Input {
		execContext.Variables[key] = value
	}

	// Execute steps
	for _, step := range workflowDef.Steps {
		select {
		case <-execContext.Context.Done():
			e.cancelExecution(execContext, "execution cancelled or timed out")
			return
		default:
		}

		if err := e.executeStep(execContext, &step); err != nil {
			if step.ErrorHandling != nil && step.ErrorHandling.ContinueOnError {
				e.logger.WithFields(logrus.Fields{
					"execution_id": execContext.Execution.ID,
					"step_id":      step.ID,
					"error":        err.Error(),
				}).Warn("Step failed but continuing execution")
				continue
			}

			// Handle retries
			if step.ErrorHandling != nil && step.ErrorHandling.RetryPolicy != nil {
				if e.shouldRetry(execContext, &step, err) {
					e.retryStep(execContext, &step)
					continue
				}
			}

			e.failExecution(execContext, fmt.Errorf("step %s failed: %w", step.ID, err))
			return
		}
	}

	// Set output
	if workflowDef.Output != nil {
		execContext.Output = e.evaluateDataMapping(execContext, workflowDef.Output)
	} else {
		execContext.Output = execContext.Variables
	}

	e.completeExecution(execContext)
}

// executeStep executes a single workflow step
func (e *Engine) executeStep(execContext *ExecutionContext, step *models.WorkflowStep) error {
	execContext.mu.Lock()
	execContext.CurrentStep = step.ID
	execContext.mu.Unlock()

	// Create step execution record
	stepExecution := &models.StepExecution{
		ID:          uuid.New(),
		ExecutionID: execContext.Execution.ID,
		StepID:      step.ID,
		Status:      models.StepStatusRunning,
		StartedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Get step executor
	e.mu.RLock()
	executor, exists := e.stepExecutors[step.Type]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no executor found for step type: %s", step.Type)
	}

	// Prepare step input
	stepInput := make(map[string]interface{})
	if step.Input != nil {
		stepInput = e.evaluateDataMapping(execContext, step.Input)
	}

	// Emit step started event
	e.emitEvent(&WorkflowEvent{
		Type:        "step.started",
		ExecutionID: execContext.Execution.ID,
		WorkflowID:  execContext.Workflow.ID,
		StepID:      step.ID,
		Timestamp:   time.Now().UTC(),
		Data: map[string]interface{}{
			"step_type": step.Type,
			"input":     stepInput,
		},
	})

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"step_id":      step.ID,
		"step_type":    step.Type,
	}).Info("Executing workflow step")

	// Execute step
	startTime := time.Now()
	output, err := executor.Execute(execContext.Context, step, stepInput)
	duration := time.Since(startTime)

	if err != nil {
		stepExecution.Status = models.StepStatusFailed
		stepExecution.Error = err.Error()
		stepExecution.CompletedAt = &[]time.Time{time.Now().UTC()}[0]
		stepExecution.Duration = int64(duration.Seconds())

		// Emit step failed event
		e.emitEvent(&WorkflowEvent{
			Type:        "step.failed",
			ExecutionID: execContext.Execution.ID,
			WorkflowID:  execContext.Workflow.ID,
			StepID:      step.ID,
			Timestamp:   time.Now().UTC(),
			Error:       err.Error(),
			Data: map[string]interface{}{
				"duration": duration.Seconds(),
			},
		})

		e.metrics.RecordStepExecution(stepExecution)
		return err
	}

	// Step completed successfully
	stepExecution.Status = models.StepStatusCompleted
	stepExecution.Output = output
	stepExecution.CompletedAt = &[]time.Time{time.Now().UTC()}[0]
	stepExecution.Duration = int64(duration.Seconds())

	// Store step result
	execContext.mu.Lock()
	execContext.StepResults[step.ID] = output
	// Apply output mapping to variables
	if step.Output != nil {
		mappedOutput := e.evaluateDataMapping(execContext, step.Output)
		for key, value := range mappedOutput {
			execContext.Variables[key] = value
		}
	} else {
		// Default: merge output into variables
		for key, value := range output {
			execContext.Variables[key] = value
		}
	}
	execContext.mu.Unlock()

	// Emit step completed event
	e.emitEvent(&WorkflowEvent{
		Type:        "step.completed",
		ExecutionID: execContext.Execution.ID,
		WorkflowID:  execContext.Workflow.ID,
		StepID:      step.ID,
		Timestamp:   time.Now().UTC(),
		Data: map[string]interface{}{
			"output":   output,
			"duration": duration.Seconds(),
		},
	})

	e.metrics.RecordStepExecution(stepExecution)

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"step_id":      step.ID,
		"duration":     duration.Seconds(),
	}).Info("Workflow step completed")

	return nil
}

// shouldRetry determines if a step should be retried
func (e *Engine) shouldRetry(execContext *ExecutionContext, step *models.WorkflowStep, err error) bool {
	if step.ErrorHandling == nil || step.ErrorHandling.RetryPolicy == nil {
		return false
	}

	retryPolicy := step.ErrorHandling.RetryPolicy
	if execContext.RetryCount >= retryPolicy.MaxRetries {
		return false
	}

	// Check retry conditions if specified
	if len(retryPolicy.RetryOn) > 0 {
		// Simple error message matching
		errorMsg := err.Error()
		for _, condition := range retryPolicy.RetryOn {
			if condition == errorMsg {
				return true
			}
		}
		return false
	}

	return true
}

// retryStep retries a failed step
func (e *Engine) retryStep(execContext *ExecutionContext, step *models.WorkflowStep) {
	execContext.mu.Lock()
	execContext.RetryCount++
	execContext.mu.Unlock()

	retryPolicy := step.ErrorHandling.RetryPolicy
	delay := time.Duration(retryPolicy.InitialDelay) * time.Second
	if retryPolicy.BackoffMultiplier > 0 {
		for i := 0; i < execContext.RetryCount-1; i++ {
			delay = time.Duration(float64(delay) * retryPolicy.BackoffMultiplier)
		}
	}

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"step_id":      step.ID,
		"retry_count":  execContext.RetryCount,
		"delay":        delay.Seconds(),
	}).Info("Retrying workflow step")

	// Wait before retry
	select {
	case <-time.After(delay):
	case <-execContext.Context.Done():
		return
	}

	// Retry the step
	e.executeStep(execContext, step)
}

// completeExecution marks an execution as completed
func (e *Engine) completeExecution(execContext *ExecutionContext) {
	now := time.Now().UTC()
	execContext.EndTime = &now

	execContext.Execution.Status = models.ExecutionStatusCompleted
	execContext.Execution.Output = execContext.Output
	execContext.Execution.CompletedAt = &now
	execContext.Execution.Duration = int64(now.Sub(execContext.StartTime).Seconds())
	execContext.Execution.UpdatedAt = now

	// Emit execution completed event
	e.emitEvent(&WorkflowEvent{
		Type:        "execution.completed",
		ExecutionID: execContext.Execution.ID,
		WorkflowID:  execContext.Workflow.ID,
		Timestamp:   now,
		Data: map[string]interface{}{
			"output":   execContext.Output,
			"duration": execContext.Execution.Duration,
		},
	})

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"workflow_id":  execContext.Workflow.ID,
		"duration":     execContext.Execution.Duration,
	}).Info("Workflow execution completed")
}

// failExecution marks an execution as failed
func (e *Engine) failExecution(execContext *ExecutionContext, err error) {
	now := time.Now().UTC()
	execContext.EndTime = &now

	execContext.Execution.Status = models.ExecutionStatusFailed
	execContext.Execution.Error = err.Error()
	execContext.Execution.CompletedAt = &now
	execContext.Execution.Duration = int64(now.Sub(execContext.StartTime).Seconds())
	execContext.Execution.UpdatedAt = now

	// Emit execution failed event
	e.emitEvent(&WorkflowEvent{
		Type:        "execution.failed",
		ExecutionID: execContext.Execution.ID,
		WorkflowID:  execContext.Workflow.ID,
		Timestamp:   now,
		Error:       err.Error(),
		Data: map[string]interface{}{
			"duration": execContext.Execution.Duration,
		},
	})

	e.metrics.RecordError(err, map[string]interface{}{
		"execution_id": execContext.Execution.ID,
		"workflow_id":  execContext.Workflow.ID,
	})

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"workflow_id":  execContext.Workflow.ID,
		"error":        err.Error(),
		"duration":     execContext.Execution.Duration,
	}).Error("Workflow execution failed")
}

// cancelExecution cancels an execution
func (e *Engine) cancelExecution(execContext *ExecutionContext, reason string) {
	now := time.Now().UTC()
	execContext.EndTime = &now

	execContext.Execution.Status = models.ExecutionStatusCancelled
	execContext.Execution.Error = reason
	execContext.Execution.CompletedAt = &now
	execContext.Execution.Duration = int64(now.Sub(execContext.StartTime).Seconds())
	execContext.Execution.UpdatedAt = now

	// Emit execution cancelled event
	e.emitEvent(&WorkflowEvent{
		Type:        "execution.cancelled",
		ExecutionID: execContext.Execution.ID,
		WorkflowID:  execContext.Workflow.ID,
		Timestamp:   now,
		Data: map[string]interface{}{
			"reason":   reason,
			"duration": execContext.Execution.Duration,
		},
	})

	e.logger.WithFields(logrus.Fields{
		"execution_id": execContext.Execution.ID,
		"workflow_id":  execContext.Workflow.ID,
		"reason":       reason,
		"duration":     execContext.Execution.Duration,
	}).Info("Workflow execution cancelled")
}

// CancelExecution cancels a running execution
func (e *Engine) CancelExecution(executionID uuid.UUID) error {
	e.mu.RLock()
	execContext, exists := e.executions[executionID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	execContext.Cancel()
	return nil
}

// GetExecution gets an execution context
func (e *Engine) GetExecution(executionID uuid.UUID) (*ExecutionContext, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	execContext, exists := e.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	return execContext, nil
}

// ListExecutions lists all running executions
func (e *Engine) ListExecutions() []*ExecutionContext {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executions := make([]*ExecutionContext, 0, len(e.executions))
	for _, execContext := range e.executions {
		executions = append(executions, execContext)
	}

	return executions
}

// emitEvent emits a workflow event to all registered handlers
func (e *Engine) emitEvent(event *WorkflowEvent) {
	e.mu.RLock()
	handlers := make([]EventHandler, len(e.eventHandlers))
	copy(handlers, e.eventHandlers)
	e.mu.RUnlock()

	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.Handle(event); err != nil {
				e.logger.WithFields(logrus.Fields{
					"event_type": event.Type,
					"error":      err.Error(),
				}).Error("Failed to handle workflow event")
			}
		}(handler)
	}
}

// evaluateDataMapping evaluates data mapping expressions
func (e *Engine) evaluateDataMapping(execContext *ExecutionContext, mapping *models.DataMapping) map[string]interface{} {
	result := make(map[string]interface{})

	if mapping == nil {
		return result
	}

	execContext.mu.RLock()
	defer execContext.mu.RUnlock()

	// Simple implementation - in a real system, you'd want a more sophisticated expression evaluator
	for key, expr := range *mapping {
		if exprStr, ok := expr.(string); ok {
			// Handle variable references like ${variable_name}
			if len(exprStr) > 3 && exprStr[:2] == "${" && exprStr[len(exprStr)-1:] == "}" {
				varName := exprStr[2 : len(exprStr)-1]
				if value, exists := execContext.Variables[varName]; exists {
					result[key] = value
				} else if value, exists := execContext.StepResults[varName]; exists {
					result[key] = value
				}
			} else {
				// Literal value
				result[key] = expr
			}
		} else {
			// Non-string value
			result[key] = expr
		}
	}

	return result
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down workflow engine")

	// Signal shutdown
	close(e.shutdownCh)

	// Cancel all running executions
	e.mu.RLock()
	for _, execContext := range e.executions {
		execContext.Cancel()
	}
	e.mu.RUnlock()

	// Wait for all executions to complete or timeout
	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		e.logger.Info("All workflow executions completed")
		return nil
	case <-ctx.Done():
		e.logger.Warn("Shutdown timeout reached, forcing exit")
		return ctx.Err()
	}
}