package core

import (
	"context"

	"github.com/truongtu268/magic-flow/pkg/events"
)

// Step represents a single workflow step
type Step interface {
	// Execute runs the step logic and returns the next step name
	Execute(ctx *WorkflowContext) (*string, error)
	// GetName returns the step name
	GetName() string
	// GetDescription returns the step description
	GetDescription() string
}

// Middleware provides cross-cutting concerns for step execution
type Middleware interface {
	// Handle processes the step with middleware logic
	Handle(ctx *WorkflowContext, next StepHandler) (*string, error)
}

// StepHandler is a function type for handling step execution
type StepHandler func(ctx *WorkflowContext) (*string, error)

// MiddlewareChain manages a chain of middleware for step execution
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

// Add adds a middleware to the chain
func (mc *MiddlewareChain) Add(middleware Middleware) {
	mc.middlewares = append(mc.middlewares, middleware)
}

// Execute executes the middleware chain
func (mc *MiddlewareChain) Execute(ctx *WorkflowContext, handler StepHandler) (*string, error) {
	if len(mc.middlewares) == 0 {
		return handler(ctx)
	}
	
	index := 0
	var next StepHandler
	next = func(ctx *WorkflowContext) (*string, error) {
		if index >= len(mc.middlewares) {
			return handler(ctx)
		}
		middleware := mc.middlewares[index]
		index++
		return middleware.Handle(ctx, next)
	}
	
	return next(ctx)
}

// LoggingMiddleware provides logging for step execution
type LoggingMiddleware struct {
	Logger Logger
}

// Handle implements the Middleware interface
func (lm *LoggingMiddleware) Handle(ctx *WorkflowContext, next StepHandler) (*string, error) {
	lm.Logger.Info("Executing step", map[string]interface{}{
		"workflow_id": ctx.GetWorkflowID(),
		"step_order":  ctx.StepOrder,
	})
	
	result, err := next(ctx)
	
	if err != nil {
		lm.Logger.Error("Step execution failed", map[string]interface{}{
			"workflow_id": ctx.GetWorkflowID(),
			"step_order":  ctx.StepOrder,
			"error":       err.Error(),
		})
	} else {
		lm.Logger.Info("Step execution completed", map[string]interface{}{
			"workflow_id": ctx.GetWorkflowID(),
			"step_order":  ctx.StepOrder,
			"result":      result,
		})
	}
	
	return result, err
}

// Logger defines the logging interface for the workflow engine
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
}

// WorkflowData represents the data structure for workflow execution
type WorkflowData interface {
	// Validate checks if the data structure is valid
	Validate() error
	// Convert converts the data to the target type
	Convert(target interface{}) error
	// GetAll returns all data as a map
	GetAll() map[string]interface{}
	// Get retrieves a value by key
	Get(key string) (interface{}, bool)
	// Set stores a value by key
	Set(key string, value interface{})
	// Delete removes a value by key
	Delete(key string)
	// Has checks if a key exists
	Has(key string) bool
	// Keys returns all keys
	Keys() []string
	// Clear removes all data
	Clear()
	// Size returns the number of items
	Size() int
	// ToMap returns data as a map
	ToMap() map[string]interface{}
	// FromMap loads data from a map
	FromMap(data map[string]interface{})
	// MustGet retrieves a value by key, panics if not found
	MustGet(key string) interface{}
}

// WorkflowMetadata represents metadata for workflow execution
type WorkflowMetadata interface {
	// GetExecutionMetrics returns execution metrics
	GetExecutionMetrics() map[string]interface{}
	// SetExecutionMetric sets an execution metric
	SetExecutionMetric(key string, value interface{})
	// GetExecutionMetric gets an execution metric
	GetExecutionMetric(key string) (interface{}, bool)
	// AddTag adds a tag
	AddTag(tag string)
	// HasTag checks if a tag exists
	HasTag(tag string) bool
	// GetTags returns all tags
	GetTags() []string
	// RemoveTag removes a tag
	RemoveTag(tag string)
	// SetCustomField sets a custom field
	SetCustomField(key string, value interface{})
	// GetCustomField gets a custom field
	GetCustomField(key string) (interface{}, bool)
	// GetCustomFields returns all custom fields
	GetCustomFields() map[string]interface{}
	// ToMap returns metadata as a map
	ToMap() map[string]interface{}
	// Validate checks if the metadata is valid
	Validate() error
}

// Engine defines the core workflow engine interface
type Engine interface {
	// RegisterStep registers a step with the engine
	RegisterStep(step Step) error
	// AddMiddleware adds middleware to the engine
	AddMiddleware(middleware Middleware)
	// ExecuteWorkflow executes a workflow starting from the given step
	ExecuteWorkflow(ctx context.Context, workflowID, workflowName, startStep string, data WorkflowData) (*WorkflowContext, error)
	// ExecuteStep executes a single step
	ExecuteStep(ctx *WorkflowContext, stepName string) (*string, error)
	// ListSteps returns all registered step names
	ListSteps() []string
	// Stop stops the engine
	Stop() error
}

// Re-export event types for convenience
type WorkflowEventType = events.WorkflowEventType
type WorkflowEvent = events.WorkflowEvent
type WorkflowEventHandler = events.WorkflowEventHandler

// Re-export status types for convenience
type WorkflowStatus = events.WorkflowStatus
type JobStatus = events.JobStatus

// Re-export event constants
const (
	WorkflowEventStarted      = events.WorkflowEventStarted
	WorkflowEventCompleted    = events.WorkflowEventCompleted
	WorkflowEventFailed       = events.WorkflowEventFailed
	WorkflowEventCancelled    = events.WorkflowEventCancelled
	WorkflowEventPaused       = events.WorkflowEventPaused
	WorkflowEventResumed      = events.WorkflowEventResumed
	WorkflowEventStepStarted  = events.WorkflowEventStepStarted
	WorkflowEventStepCompleted = events.WorkflowEventStepCompleted
	WorkflowEventStepFailed   = events.WorkflowEventStepFailed
)

// Re-export status constants
const (
	WorkflowStatusPending   = events.WorkflowStatusPending
	WorkflowStatusRunning   = events.WorkflowStatusRunning
	WorkflowStatusCompleted = events.WorkflowStatusCompleted
	WorkflowStatusFailed    = events.WorkflowStatusFailed
	WorkflowStatusCancelled = events.WorkflowStatusCancelled
	WorkflowStatusPaused    = events.WorkflowStatusPaused
	WorkflowStatusUnknown   = events.WorkflowStatusUnknown

	JobStatusPending   = events.JobStatusPending
	JobStatusRunning   = events.JobStatusRunning
	JobStatusCompleted = events.JobStatusCompleted
	JobStatusFailed    = events.JobStatusFailed
	JobStatusCancelled = events.JobStatusCancelled
)