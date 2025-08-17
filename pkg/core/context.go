package core

import (
	"context"
	"sync"
	"time"
)

// WorkflowContext represents the execution context of a workflow
type WorkflowContext struct {
	WorkflowID   string                 `json:"workflow_id"`
	WorkflowName string                 `json:"workflow_name"`
	CurrentStep  string                 `json:"current_step"`
	NextStep     *string                `json:"next_step"`
	Data         WorkflowData           `json:"data"`
	Metadata     WorkflowMetadata       `json:"metadata"`
	StepResults  map[string]interface{} `json:"step_results"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      *time.Time             `json:"end_time"`
	Status       WorkflowStatus         `json:"status"`
	Error        error                  `json:"error,omitempty"`
	StepOrder    int                    `json:"step_order"`
	ctx          context.Context        `json:"-"`
	mu           sync.RWMutex           `json:"-"`
}

// NewWorkflowContext creates a new workflow context
func NewWorkflowContext(ctx context.Context, workflowID, workflowName string, data WorkflowData, metadata WorkflowMetadata) *WorkflowContext {
	if data == nil {
		data = NewDefaultWorkflowData()
	}
	if metadata == nil {
		metadata = NewDefaultWorkflowMetadata()
	}
	
	return &WorkflowContext{
		WorkflowID:   workflowID,
		WorkflowName: workflowName,
		Data:         data,
		Metadata:     metadata,
		StepResults:  make(map[string]interface{}),
		StartTime:    time.Now(),
		Status:       WorkflowStatusPending,
		StepOrder:    0,
		ctx:          ctx,
	}
}

// NewWorkflowContextSimple creates a new workflow context with default metadata
func NewWorkflowContextSimple(workflowID, workflowName string, data WorkflowData) *WorkflowContext {
	return &WorkflowContext{
		WorkflowID:   workflowID,
		WorkflowName: workflowName,
		Data:         data,
		Metadata:     NewDefaultWorkflowMetadata(),
		StepResults:  make(map[string]interface{}),
		StartTime:    time.Now(),
		Status:       WorkflowStatusPending,
		StepOrder:    0,
		ctx:          context.Background(),
	}
}

// GetContext returns the underlying context
func (wc *WorkflowContext) GetContext() context.Context {
	return wc.ctx
}

// SetContext sets the underlying context
func (wc *WorkflowContext) SetContext(ctx context.Context) {
	wc.ctx = ctx
}

// GetWorkflowID returns the workflow ID
func (wc *WorkflowContext) GetWorkflowID() string {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.WorkflowID
}

// GetWorkflowName returns the workflow name
func (wc *WorkflowContext) GetWorkflowName() string {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.WorkflowName
}

// GetCurrentStep returns the current step
func (wc *WorkflowContext) GetCurrentStep() string {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.CurrentStep
}

// SetCurrentStep sets the current step
func (wc *WorkflowContext) SetCurrentStep(step string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.CurrentStep = step
}

// GetNextStep returns the next step
func (wc *WorkflowContext) GetNextStep() *string {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.NextStep
}

// SetNextStep sets the next step
func (wc *WorkflowContext) SetNextStep(step string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.NextStep = &step
}

// ClearNextStep clears the next step
func (wc *WorkflowContext) ClearNextStep() {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.NextStep = nil
}

// GetStatus returns the workflow status
func (wc *WorkflowContext) GetStatus() WorkflowStatus {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.Status
}

// SetStatus sets the workflow status
func (wc *WorkflowContext) SetStatus(status WorkflowStatus) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.Status = status
}

// GetError returns the workflow error
func (wc *WorkflowContext) GetError() error {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.Error
}

// SetError sets the workflow error
func (wc *WorkflowContext) SetError(err error) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.Error = err
	wc.Status = WorkflowStatusFailed
}

// Complete marks the workflow as completed
func (wc *WorkflowContext) Complete() {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	now := time.Now()
	wc.Status = WorkflowStatusCompleted
	wc.EndTime = &now
}

// Cancel marks the workflow as cancelled
func (wc *WorkflowContext) Cancel() {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	now := time.Now()
	wc.Status = WorkflowStatusCancelled
	wc.EndTime = &now
}

// SetWaiting marks the workflow as waiting
func (wc *WorkflowContext) SetWaiting(triggerKey string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.Status = WorkflowStatusPaused
	wc.Metadata.SetExecutionMetric("trigger_key", triggerKey)
	wc.Metadata.SetExecutionMetric("waiting_since", time.Now())
}

// GetData returns the workflow data
func (wc *WorkflowContext) GetData(key string) (interface{}, bool) {
	return wc.Data.Get(key)
}

// SetData sets workflow data
func (wc *WorkflowContext) SetData(key string, value interface{}) {
	wc.Data.Set(key, value)
}

// GetStepResult returns a step result
func (wc *WorkflowContext) GetStepResult(stepName string) (interface{}, bool) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	result, exists := wc.StepResults[stepName]
	return result, exists
}

// SetStepResult sets a step result
func (wc *WorkflowContext) SetStepResult(stepName string, result interface{}) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.StepResults[stepName] = result
}

// GetAllStepResults returns all step results
func (wc *WorkflowContext) GetAllStepResults() map[string]interface{} {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	results := make(map[string]interface{})
	for k, v := range wc.StepResults {
		results[k] = v
	}
	return results
}

// IncrementStepOrder increments and returns the current step order
func (wc *WorkflowContext) IncrementStepOrder() int {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.StepOrder++
	return wc.StepOrder
}

// GetStepOrder returns the current step order
func (wc *WorkflowContext) GetStepOrder() int {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.StepOrder
}