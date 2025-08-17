package core

import (
	"fmt"
)

// BaseStep provides a base implementation for steps
type BaseStep struct {
	Name        string
	Description string
}

// GetName returns the step name
func (s *BaseStep) GetName() string {
	return s.Name
}

// GetDescription returns the step description
func (s *BaseStep) GetDescription() string {
	return s.Description
}

// Execute is a placeholder implementation that should be overridden
func (s *BaseStep) Execute(ctx *WorkflowContext) (*string, error) {
	return nil, fmt.Errorf("execute method not implemented for step %s", s.Name)
}

// NewBaseStep creates a new base step
func NewBaseStep(name, description string) *BaseStep {
	return &BaseStep{
		Name:        name,
		Description: description,
	}
}

// FunctionStep wraps a function as a step
type FunctionStep struct {
	*BaseStep
	ExecuteFunc func(ctx *WorkflowContext) (*string, error)
}

// Execute runs the wrapped function
func (s *FunctionStep) Execute(ctx *WorkflowContext) (*string, error) {
	if s.ExecuteFunc == nil {
		return nil, fmt.Errorf("execute function not set for step %s", s.Name)
	}
	return s.ExecuteFunc(ctx)
}

// NewFunctionStep creates a new function step
func NewFunctionStep(name, description string, executeFunc func(ctx *WorkflowContext) (*string, error)) *FunctionStep {
	return &FunctionStep{
		BaseStep:    NewBaseStep(name, description),
		ExecuteFunc: executeFunc,
	}
}

// ConditionalStep executes different logic based on a condition
type ConditionalStep struct {
	*BaseStep
	ConditionFunc func(ctx *WorkflowContext) (bool, error)
	TrueStep      string
	FalseStep     string
}

// Execute evaluates the condition and returns the appropriate next step
func (s *ConditionalStep) Execute(ctx *WorkflowContext) (*string, error) {
	if s.ConditionFunc == nil {
		return nil, fmt.Errorf("condition function not set for step %s", s.Name)
	}
	
	condition, err := s.ConditionFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("condition evaluation failed for step %s: %w", s.Name, err)
	}
	
	if condition {
		if s.TrueStep != "" {
			return &s.TrueStep, nil
		}
	} else {
		if s.FalseStep != "" {
			return &s.FalseStep, nil
		}
	}
	
	return nil, nil // End workflow
}

// NewConditionalStep creates a new conditional step
func NewConditionalStep(name, description string, conditionFunc func(ctx *WorkflowContext) (bool, error), trueStep, falseStep string) *ConditionalStep {
	return &ConditionalStep{
		BaseStep:      NewBaseStep(name, description),
		ConditionFunc: conditionFunc,
		TrueStep:      trueStep,
		FalseStep:     falseStep,
	}
}

// WaitStep represents a step that waits for external trigger
type WaitStep struct {
	*BaseStep
	TriggerKey string
	NextStep   string
}

// Execute sets the workflow to waiting status
func (s *WaitStep) Execute(ctx *WorkflowContext) (*string, error) {
	ctx.SetWaiting(s.TriggerKey)
	if s.NextStep != "" {
		return &s.NextStep, nil
	}
	return nil, nil
}

// NewWaitStep creates a new wait step
func NewWaitStep(name, description, triggerKey, nextStep string) *WaitStep {
	return &WaitStep{
		BaseStep:   NewBaseStep(name, description),
		TriggerKey: triggerKey,
		NextStep:   nextStep,
	}
}

// ParallelStep executes multiple steps in parallel
type ParallelStep struct {
	*BaseStep
	Steps    []string
	NextStep string
}

// Execute marks the step for parallel execution
func (s *ParallelStep) Execute(ctx *WorkflowContext) (*string, error) {
	// Store parallel steps in metadata for the engine to handle
	ctx.Metadata.SetExecutionMetric("parallel_steps", s.Steps)
	ctx.Metadata.SetExecutionMetric("parallel_next_step", s.NextStep)
	
	if s.NextStep != "" {
		return &s.NextStep, nil
	}
	return nil, nil
}

// NewParallelStep creates a new parallel step
func NewParallelStep(name, description string, steps []string, nextStep string) *ParallelStep {
	return &ParallelStep{
		BaseStep: NewBaseStep(name, description),
		Steps:    steps,
		NextStep: nextStep,
	}
}

// RetryStep wraps another step with retry logic
type RetryStep struct {
	*BaseStep
	WrappedStep Step
	MaxRetries  int
	RetryCount  int
}

// Execute runs the wrapped step with retry logic
func (s *RetryStep) Execute(ctx *WorkflowContext) (*string, error) {
	var lastErr error
	
	for attempt := 0; attempt <= s.MaxRetries; attempt++ {
		nextStep, err := s.WrappedStep.Execute(ctx)
		if err == nil {
			s.RetryCount = attempt
			return nextStep, nil
		}
		
		lastErr = err
		if attempt < s.MaxRetries {
			// Store retry attempt in metadata
			ctx.Metadata.SetExecutionMetric(fmt.Sprintf("retry_attempt_%d", attempt+1), err.Error())
		}
	}
	
	s.RetryCount = s.MaxRetries + 1
	return nil, fmt.Errorf("step %s failed after %d retries: %w", s.Name, s.MaxRetries, lastErr)
}

// NewRetryStep creates a new retry step
func NewRetryStep(name, description string, wrappedStep Step, maxRetries int) *RetryStep {
	return &RetryStep{
		BaseStep:    NewBaseStep(name, description),
		WrappedStep: wrappedStep,
		MaxRetries:  maxRetries,
		RetryCount:  0,
	}
}