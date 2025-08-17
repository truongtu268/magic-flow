package middleware

import (
	"fmt"
	"time"

	"github.com/truongtu268/magic-flow/pkg/core"
)



// TimingMiddleware measures step execution time
type TimingMiddleware struct {
	Logger core.Logger
}

// NewTimingMiddleware creates a new timing middleware
func NewTimingMiddleware(logger core.Logger) *TimingMiddleware {
	return &TimingMiddleware{
		Logger: logger,
	}
}

// Handle measures and logs step execution time
func (m *TimingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	stepName := ctx.GetCurrentStep()
	workflowID := ctx.GetWorkflowID()
	startTime := time.Now()
	
	nextStep, err := next(ctx)
	
	duration := time.Since(startTime)
	
	// Store timing information in metadata
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_duration", stepName), duration.String())
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_start_time", stepName), startTime)
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_end_time", stepName), time.Now())
	
	if m.Logger != nil {
		m.Logger.Info("Step timing", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "duration": duration.String()})
	}
	
	return nextStep, err
}

// ErrorHandlingMiddleware provides error handling and recovery
type ErrorHandlingMiddleware struct {
	Logger         core.Logger
	RecoveryFunc   func(ctx *core.WorkflowContext, err error) (*string, error)
	IgnoreErrors   []string
	RetryOnErrors  []string
	MaxRetries     int
}

// NewErrorHandlingMiddleware creates a new error handling middleware
func NewErrorHandlingMiddleware(logger core.Logger) *ErrorHandlingMiddleware {
	return &ErrorHandlingMiddleware{
		Logger:     logger,
		MaxRetries: 3,
	}
}

// Handle provides error handling and recovery
func (m *ErrorHandlingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	stepName := ctx.GetCurrentStep()
	workflowID := ctx.GetWorkflowID()
	
	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("panic in step %s: %v", stepName, r)
			if m.Logger != nil {
				m.Logger.Error("Step panic recovered", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "panic": r})
			}
			ctx.SetError(panicErr)
		}
	}()
	
	nextStep, err := next(ctx)
	
	if err != nil {
		// Check if error should be ignored
		for _, ignoreErr := range m.IgnoreErrors {
			if err.Error() == ignoreErr {
				if m.Logger != nil {
					m.Logger.Warn("Ignoring error", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "error": err.Error()})
				}
				return nextStep, nil
			}
		}
		
		// Check if error should trigger retry
		for _, retryErr := range m.RetryOnErrors {
			if err.Error() == retryErr {
				retryCount, _ := ctx.Metadata.GetExecutionMetric(fmt.Sprintf("step_%s_retry_count", stepName))
				count, ok := retryCount.(int)
				if !ok {
					count = 0
				}
				
				if count < m.MaxRetries {
					ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_retry_count", stepName), count+1)
					if m.Logger != nil {
						m.Logger.Warn("Retrying step", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "retry_count": count+1, "error": err.Error()})
					}
					return next(ctx) // Retry
				}
			}
		}
		
		// Use custom recovery function if provided
		if m.RecoveryFunc != nil {
			recoveryNext, recoveryErr := m.RecoveryFunc(ctx, err)
			if recoveryErr == nil {
				if m.Logger != nil {
					m.Logger.Info("Error recovered", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "original_error": err.Error()})
				}
				return recoveryNext, nil
			}
		}
		
		// Log error and return
		if m.Logger != nil {
			m.Logger.Error("Step execution failed", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "error": err.Error()})
		}
	}
	
	return nextStep, err
}

// ValidationMiddleware validates workflow context and data
type ValidationMiddleware struct {
	Logger         core.Logger
	ValidationFunc func(ctx *core.WorkflowContext) error
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware(logger core.Logger, validationFunc func(ctx *core.WorkflowContext) error) *ValidationMiddleware {
	return &ValidationMiddleware{
		Logger:         logger,
		ValidationFunc: validationFunc,
	}
}

// Handle validates workflow context before step execution
func (m *ValidationMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	stepName := ctx.GetCurrentStep()
	workflowID := ctx.GetWorkflowID()
	
	// Validate workflow data
	if err := ctx.Data.Validate(); err != nil {
		validationErr := fmt.Errorf("workflow data validation failed for step %s: %w", stepName, err)
		if m.Logger != nil {
			m.Logger.Error("Data validation failed", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "error": err.Error()})
		}
		return nil, validationErr
	}
	
	// Custom validation if provided
	if m.ValidationFunc != nil {
		if err := m.ValidationFunc(ctx); err != nil {
			validationErr := fmt.Errorf("custom validation failed for step %s: %w", stepName, err)
			if m.Logger != nil {
				m.Logger.Error("Custom validation failed", map[string]interface{}{"workflow_id": workflowID, "step": stepName, "error": err.Error()})
			}
			return nil, validationErr
		}
	}
	
	return next(ctx)
}

// MetricsMiddleware collects execution metrics
type MetricsMiddleware struct {
	Logger      core.Logger
	MetricsFunc func(stepName string, duration time.Duration, success bool)
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(logger core.Logger, metricsFunc func(stepName string, duration time.Duration, success bool)) *MetricsMiddleware {
	return &MetricsMiddleware{
		Logger:      logger,
		MetricsFunc: metricsFunc,
	}
}

// Handle collects execution metrics
func (m *MetricsMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	stepName := ctx.GetCurrentStep()
	startTime := time.Now()
	
	nextStep, err := next(ctx)
	
	duration := time.Since(startTime)
	success := err == nil
	
	// Store metrics in context
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_execution_count", stepName), 1)
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_success", stepName), success)
	ctx.Metadata.SetExecutionMetric(fmt.Sprintf("step_%s_duration_ms", stepName), duration.Milliseconds())
	
	// Call custom metrics function if provided
	if m.MetricsFunc != nil {
		m.MetricsFunc(stepName, duration, success)
	}
	
	return nextStep, err
}

// RateLimitingMiddleware provides rate limiting for step execution
type RateLimitingMiddleware struct {
	Logger       core.Logger
	RateLimit    int           // requests per second
	BurstLimit   int           // burst capacity
	WindowSize   time.Duration // time window for rate limiting
	lastExecution map[string]time.Time
}

// NewRateLimitingMiddleware creates a new rate limiting middleware
func NewRateLimitingMiddleware(logger core.Logger, rateLimit, burstLimit int, windowSize time.Duration) *RateLimitingMiddleware {
	return &RateLimitingMiddleware{
		Logger:        logger,
		RateLimit:     rateLimit,
		BurstLimit:    burstLimit,
		WindowSize:    windowSize,
		lastExecution: make(map[string]time.Time),
	}
}

// Handle applies rate limiting to step execution
func (m *RateLimitingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	stepName := ctx.GetCurrentStep()
	workflowID := ctx.GetWorkflowID()
	now := time.Now()
	
	// Check rate limit
	key := fmt.Sprintf("%s:%s", workflowID, stepName)
	if lastExec, exists := m.lastExecution[key]; exists {
		if now.Sub(lastExec) < m.WindowSize {
			rateLimitErr := fmt.Errorf("rate limit exceeded for step %s in workflow %s", stepName, workflowID)
			if m.Logger != nil {
				m.Logger.Warn("Rate limit exceeded", map[string]interface{}{"workflow_id": workflowID, "step": stepName})
			}
			return nil, rateLimitErr
		}
	}
	
	m.lastExecution[key] = now
	return next(ctx)
}