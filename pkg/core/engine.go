package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/truongtu268/magic-flow/pkg/config"
	"github.com/truongtu268/magic-flow/pkg/errors"
	"github.com/truongtu268/magic-flow/pkg/messaging"
	"github.com/truongtu268/magic-flow/pkg/storage"
)

// WorkflowEngine is the main engine implementation
type WorkflowEngine struct {
	config           *config.Config
	storage          storage.WorkflowStorage
	messaging        messaging.MessageQueue
	pubsub           messaging.PubSubService
	middlewareChain  *MiddlewareChain
	logger           Logger
	eventHandlers    map[WorkflowEventType][]WorkflowEventHandler
	runningWorkflows sync.Map
	shutdownChan     chan struct{}
	shutdownOnce     sync.Once
	mu               sync.RWMutex
}

// EngineConfig holds configuration for the workflow engine
type EngineConfig struct {
	Config    *config.Config
	Storage   storage.WorkflowStorage
	Messaging messaging.MessageQueue
	PubSub    messaging.PubSubService
	Logger    Logger
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(cfg *EngineConfig) (*WorkflowEngine, error) {
	if cfg == nil {
		return nil, errors.New(errors.ErrConfigMissing, "engine config is required")
	}
	
	if cfg.Storage == nil {
		return nil, errors.New(errors.ErrConfigMissing, "storage is required")
	}
	
	if cfg.Config == nil {
		return nil, errors.New(errors.ErrConfigMissing, "config is required")
	}
	
	if cfg.Logger == nil {
		cfg.Logger = &DefaultLogger{}
	}
	
	engine := &WorkflowEngine{
		config:          cfg.Config,
		storage:         cfg.Storage,
		messaging:       cfg.Messaging,
		pubsub:          cfg.PubSub,
		logger:          cfg.Logger,
		eventHandlers:   make(map[WorkflowEventType][]WorkflowEventHandler),
		shutdownChan:    make(chan struct{}),
		middlewareChain: NewMiddlewareChain(),
	}
	
	// Add default middleware
	engine.addDefaultMiddleware()
	
	return engine, nil
}

// Execute executes a workflow
func (e *WorkflowEngine) Execute(ctx context.Context, workflowID string, steps []Step, data WorkflowData) error {
	if len(steps) == 0 {
		return errors.New(errors.ErrValidationFailed, "workflow must have at least one step")
	}
	
	// Create workflow context
	workflowCtx := NewWorkflowContext(ctx, workflowID, "default", data, NewDefaultWorkflowMetadata())
	workflowCtx.SetStatus(WorkflowStatusRunning)
	
	// Store workflow in running workflows
	e.runningWorkflows.Store(workflowID, workflowCtx)
	defer e.runningWorkflows.Delete(workflowID)
	
	// Emit workflow started event
	e.emitEvent(WorkflowEventStarted, workflowCtx)
	
	// Execute workflow with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, e.config.Engine.WorkflowTimeout)
	defer cancel()
	
	err := e.executeWorkflow(ctxWithTimeout, workflowCtx, steps)
	
	if err != nil {
		workflowCtx.SetStatus(WorkflowStatusFailed)
		workflowCtx.SetError(err)
		e.emitEvent(WorkflowEventFailed, workflowCtx)
		e.logger.Error("Workflow execution failed", map[string]interface{}{
			"workflow_id": workflowID,
			"error":       err.Error(),
		})
		return err
	}
	
	workflowCtx.SetStatus(WorkflowStatusCompleted)
	e.emitEvent(WorkflowEventCompleted, workflowCtx)
	e.logger.Info("Workflow execution completed", map[string]interface{}{
		"workflow_id": workflowID,
		"duration":    time.Since(workflowCtx.StartTime),
	})
	
	return nil
}

// ExecuteStep executes a single step
func (e *WorkflowEngine) ExecuteStep(ctx context.Context, step Step, workflowCtx *WorkflowContext) error {
	// Set current step
	workflowCtx.SetCurrentStep(step.GetName())
	
	// Execute step through middleware chain
	stepHandler := func(ctx *WorkflowContext) (*string, error) {
		return step.Execute(ctx)
	}
	_, err := e.middlewareChain.Execute(workflowCtx, stepHandler)
	
	if err != nil {
		e.logger.Error("Step execution failed", map[string]interface{}{
			"workflow_id": workflowCtx.GetWorkflowID(),
			"step_name":   step.GetName(),
			"error":       err.Error(),
		})
		return errors.NewStepFailedError(step.GetName(), err)
	}
	
	e.logger.Debug("Step executed successfully", map[string]interface{}{
		"workflow_id": workflowCtx.GetWorkflowID(),
		"step_name":   step.GetName(),
	})
	
	return nil
}

// AddEventHandler adds an event handler for a specific event type
func (e *WorkflowEngine) AddEventHandler(eventType WorkflowEventType, handler WorkflowEventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.eventHandlers[eventType] = append(e.eventHandlers[eventType], handler)
}

// RemoveEventHandler removes an event handler
func (e *WorkflowEngine) RemoveEventHandler(eventType WorkflowEventType, handler WorkflowEventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	handlers := e.eventHandlers[eventType]
	for i, h := range handlers {
		if &h == &handler {
			e.eventHandlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// AddMiddleware adds middleware to the engine
func (e *WorkflowEngine) AddMiddleware(middleware Middleware) {
	e.middlewareChain.Add(middleware)
}

// GetWorkflowStatus returns the status of a workflow
func (e *WorkflowEngine) GetWorkflowStatus(workflowID string) (WorkflowStatus, error) {
	// Check running workflows first
	if ctx, ok := e.runningWorkflows.Load(workflowID); ok {
		workflowCtx := ctx.(*WorkflowContext)
		return workflowCtx.GetStatus(), nil
	}
	
	// Check storage
	// Check storage for workflow record
	record, err := e.storage.GetWorkflowRecord(context.Background(), workflowID)
	if err != nil {
		return WorkflowStatusUnknown, errors.NewWorkflowNotFoundError(workflowID)
	}
	return record.Status, nil
	
	return WorkflowStatusUnknown, errors.NewWorkflowNotFoundError(workflowID)
}

// CancelWorkflow cancels a running workflow
func (e *WorkflowEngine) CancelWorkflow(workflowID string) error {
	if ctx, ok := e.runningWorkflows.Load(workflowID); ok {
		workflowCtx := ctx.(*WorkflowContext)
		workflowCtx.SetStatus(WorkflowStatusCancelled)
		e.emitEvent(WorkflowEventCancelled, workflowCtx)
		e.logger.Info("Workflow cancelled", map[string]interface{}{
			"workflow_id": workflowID,
		})
		return nil
	}
	
	return errors.NewWorkflowNotFoundError(workflowID)
}

// GetRunningWorkflows returns a list of currently running workflows
func (e *WorkflowEngine) GetRunningWorkflows() []string {
	var workflowIDs []string
	e.runningWorkflows.Range(func(key, value interface{}) bool {
		workflowIDs = append(workflowIDs, key.(string))
		return true
	})
	return workflowIDs
}

// Shutdown gracefully shuts down the engine
func (e *WorkflowEngine) Shutdown(ctx context.Context) error {
	var shutdownErr error
	e.shutdownOnce.Do(func() {
		e.logger.Info("Shutting down workflow engine", nil)
		
		// Signal shutdown
		close(e.shutdownChan)
		
		// Wait for running workflows to complete or timeout
		shutdownCtx, cancel := context.WithTimeout(ctx, e.config.Engine.GracefulShutdownTimeout)
		defer cancel()
		
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-shutdownCtx.Done():
				// Force shutdown
				e.logger.Warn("Force shutdown due to timeout", nil)
				return
			case <-ticker.C:
				// Check if all workflows are done
				runningCount := 0
				e.runningWorkflows.Range(func(key, value interface{}) bool {
					runningCount++
					return true
				})
				
				if runningCount == 0 {
					e.logger.Info("All workflows completed, shutdown complete", nil)
					return
				}
				
				e.logger.Debug("Waiting for workflows to complete", map[string]interface{}{
					"running_count": runningCount,
				})
			}
		}
	})
	
	return shutdownErr
}

// Private methods

func (e *WorkflowEngine) executeWorkflow(ctx context.Context, workflowCtx *WorkflowContext, steps []Step) error {
	for i, step := range steps {
		select {
		case <-ctx.Done():
			return errors.NewWorkflowTimeoutError(workflowCtx.GetWorkflowID(), e.config.Engine.WorkflowTimeout)
		case <-e.shutdownChan:
			return errors.New(errors.ErrWorkflowCancelled, "workflow cancelled due to engine shutdown")
		default:
		}
		
		// Check if workflow was cancelled
		if workflowCtx.GetStatus() == WorkflowStatusCancelled {
			return errors.New(errors.ErrWorkflowCancelled, "workflow was cancelled")
		}
		
		// Set step order
		workflowCtx.StepOrder = i
		
		// Execute step
		if err := e.ExecuteStep(ctx, step, workflowCtx); err != nil {
			// TODO: Add recovery mechanism when needed
			return err
		}
		
		// Set next step if not the last step
		if i < len(steps)-1 {
			workflowCtx.SetNextStep(steps[i+1].GetName())
		}
	}
	
	return nil
}

func (e *WorkflowEngine) emitEvent(eventType WorkflowEventType, workflowCtx *WorkflowContext) {
	e.mu.RLock()
	handlers := e.eventHandlers[eventType]
	e.mu.RUnlock()
	
	if len(handlers) == 0 {
		return
	}
	
	event := &WorkflowEvent{
		ID:          uuid.New().String(),
		Type:        eventType,
		WorkflowID:  workflowCtx.GetWorkflowID(),
		Timestamp:   time.Now(),
		Data:        workflowCtx.Data.GetAll(),
	}
	
	// Execute handlers in goroutines to avoid blocking
	for _, handler := range handlers {
		go func(h WorkflowEventHandler) {
			defer func() {
				if r := recover(); r != nil {
					e.logger.Error("Event handler panicked", map[string]interface{}{
						"event_type": eventType,
						"panic":      r,
					})
				}
			}()
			
			if err := h(event); err != nil {
				e.logger.Error("Event handler failed", map[string]interface{}{
					"event_type": eventType,
					"error":      err.Error(),
				})
			}
		}(handler)
	}
}

func (e *WorkflowEngine) addDefaultMiddleware() {
	// Add basic logging middleware
	e.middlewareChain.Add(&LoggingMiddleware{Logger: e.logger})
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(message string, fields map[string]interface{}) {
	fmt.Printf("[DEBUG] %s %v\n", message, fields)
}

func (l *DefaultLogger) Info(message string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s %v\n", message, fields)
}

func (l *DefaultLogger) Warn(message string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s %v\n", message, fields)
}

func (l *DefaultLogger) Error(message string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s %v\n", message, fields)
}