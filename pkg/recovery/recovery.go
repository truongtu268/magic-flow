package recovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/truongtu268/magic-flow/pkg/core"
	"github.com/truongtu268/magic-flow/pkg/storage"
)

// RecoveryStrategy defines different recovery strategies
type RecoveryStrategy string

const (
	RecoveryStrategyRetry    RecoveryStrategy = "RETRY"
	RecoveryStrategySkip     RecoveryStrategy = "SKIP"
	RecoveryStrategyFail     RecoveryStrategy = "FAIL"
	RecoveryStrategyRestart  RecoveryStrategy = "RESTART"
	RecoveryStrategyCustom   RecoveryStrategy = "CUSTOM"
)

// RecoveryPolicy defines the recovery policy for workflows
type RecoveryPolicy struct {
	Strategy      RecoveryStrategy `json:"strategy"`
	MaxRetries    int              `json:"max_retries"`
	RetryDelay    time.Duration    `json:"retry_delay"`
	BackoffFactor float64          `json:"backoff_factor"`
	MaxDelay      time.Duration    `json:"max_delay"`
	Timeout       time.Duration    `json:"timeout"`
	CustomHandler RecoveryHandler  `json:"-"`
}

// DefaultRecoveryPolicy returns a default recovery policy
func DefaultRecoveryPolicy() *RecoveryPolicy {
	return &RecoveryPolicy{
		Strategy:      RecoveryStrategyRetry,
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		BackoffFactor: 2.0,
		MaxDelay:      30 * time.Second,
		Timeout:       5 * time.Minute,
	}
}

// RecoveryHandler defines the interface for custom recovery handlers
type RecoveryHandler interface {
	// Handle processes a failed workflow and returns recovery action
	Handle(ctx context.Context, workflowCtx *core.WorkflowContext, err error) (*RecoveryAction, error)
	// GetStrategy returns the recovery strategy this handler supports
	GetStrategy() RecoveryStrategy
}

// RecoveryAction defines the action to take during recovery
type RecoveryAction struct {
	Action     RecoveryStrategy       `json:"action"`
	NextStep   *string                `json:"next_step,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Delay      time.Duration          `json:"delay,omitempty"`
	Message    string                 `json:"message,omitempty"`
}

// WorkflowRecoveryManager manages workflow recovery operations
type WorkflowRecoveryManager struct {
	storage         storage.WorkflowStorage
	engine          core.Engine
	logger          core.Logger
	defaultPolicy   *RecoveryPolicy
	policies        map[string]*RecoveryPolicy // workflow-specific policies
	handlers        map[RecoveryStrategy]RecoveryHandler
	recoveryHistory map[string][]*RecoveryAttempt
	mu              sync.RWMutex
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

// RecoveryAttempt represents a recovery attempt
type RecoveryAttempt struct {
	ID          string           `json:"id"`
	WorkflowID  string           `json:"workflow_id"`
	Strategy    RecoveryStrategy `json:"strategy"`
	AttemptTime time.Time        `json:"attempt_time"`
	Success     bool             `json:"success"`
	Error       *string          `json:"error,omitempty"`
	Duration    time.Duration    `json:"duration"`
}

// NewWorkflowRecoveryManager creates a new workflow recovery manager
func NewWorkflowRecoveryManager(storage storage.WorkflowStorage, engine core.Engine, logger core.Logger) *WorkflowRecoveryManager {
	return &WorkflowRecoveryManager{
		storage:         storage,
		engine:          engine,
		logger:          logger,
		defaultPolicy:   DefaultRecoveryPolicy(),
		policies:        make(map[string]*RecoveryPolicy),
		handlers:        make(map[RecoveryStrategy]RecoveryHandler),
		recoveryHistory: make(map[string][]*RecoveryAttempt),
		stopChan:        make(chan struct{}),
	}
}

// SetDefaultPolicy sets the default recovery policy
func (rm *WorkflowRecoveryManager) SetDefaultPolicy(policy *RecoveryPolicy) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.defaultPolicy = policy
}

// SetWorkflowPolicy sets a recovery policy for a specific workflow
func (rm *WorkflowRecoveryManager) SetWorkflowPolicy(workflowName string, policy *RecoveryPolicy) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.policies[workflowName] = policy
}

// RegisterHandler registers a recovery handler for a strategy
func (rm *WorkflowRecoveryManager) RegisterHandler(handler RecoveryHandler) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.handlers[handler.GetStrategy()] = handler
}

// RecoverWorkflow attempts to recover a failed workflow
func (rm *WorkflowRecoveryManager) RecoverWorkflow(ctx context.Context, workflowID string) error {
	// Get workflow record
	record, err := rm.storage.GetWorkflowRecord(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow record: %w", err)
	}
	
	if record.Status != core.WorkflowStatusFailed {
		return fmt.Errorf("workflow %s is not in failed status", workflowID)
	}
	
	// Get recovery policy
	policy := rm.getPolicy(record.WorkflowName)
	
	// Create workflow context from record
	workflowCtx := rm.createWorkflowContextFromRecord(ctx, record)
	
	// Attempt recovery
	attempt := &RecoveryAttempt{
		ID:          fmt.Sprintf("%s-%d", workflowID, time.Now().Unix()),
		WorkflowID:  workflowID,
		Strategy:    policy.Strategy,
		AttemptTime: time.Now(),
	}
	
	startTime := time.Now()
	defer func() {
		attempt.Duration = time.Since(startTime)
		rm.addRecoveryAttempt(workflowID, attempt)
	}()
	
	// Execute recovery strategy
	var recoveryErr error
	switch policy.Strategy {
	case RecoveryStrategyRetry:
		recoveryErr = rm.retryWorkflow(ctx, workflowCtx, policy)
	case RecoveryStrategySkip:
		recoveryErr = rm.skipFailedStep(ctx, workflowCtx)
	case RecoveryStrategyFail:
		recoveryErr = rm.markWorkflowFailed(ctx, workflowCtx)
	case RecoveryStrategyRestart:
		recoveryErr = rm.restartWorkflow(ctx, workflowCtx)
	case RecoveryStrategyCustom:
		recoveryErr = rm.executeCustomRecovery(ctx, workflowCtx, policy)
	default:
		recoveryErr = fmt.Errorf("unknown recovery strategy: %s", policy.Strategy)
	}
	
	if recoveryErr != nil {
		attempt.Success = false
		errorMsg := recoveryErr.Error()
		attempt.Error = &errorMsg
		rm.logger.Error("Recovery failed", map[string]interface{}{"workflow_id": workflowID, "strategy": policy.Strategy, "error": recoveryErr.Error()})
		return recoveryErr
	}
	
	attempt.Success = true
	rm.logger.Info("Recovery successful", map[string]interface{}{"workflow_id": workflowID, "strategy": policy.Strategy})
	return nil
}

// StartRecoveryMonitor starts monitoring for failed workflows
func (rm *WorkflowRecoveryManager) StartRecoveryMonitor(ctx context.Context, interval time.Duration) {
	rm.wg.Add(1)
	go func() {
		defer rm.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-rm.stopChan:
				return
			case <-ticker.C:
				rm.monitorFailedWorkflows(ctx)
			}
		}
	}()
}

// Stop stops the recovery manager
func (rm *WorkflowRecoveryManager) Stop() {
	close(rm.stopChan)
	rm.wg.Wait()
}

// GetRecoveryHistory returns recovery history for a workflow
func (rm *WorkflowRecoveryManager) GetRecoveryHistory(workflowID string) []*RecoveryAttempt {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	history, exists := rm.recoveryHistory[workflowID]
	if !exists {
		return []*RecoveryAttempt{}
	}
	
	// Return a copy
	result := make([]*RecoveryAttempt, len(history))
	copy(result, history)
	return result
}

// Private methods

func (rm *WorkflowRecoveryManager) getPolicy(workflowName string) *RecoveryPolicy {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	if policy, exists := rm.policies[workflowName]; exists {
		return policy
	}
	return rm.defaultPolicy
}

func (rm *WorkflowRecoveryManager) createWorkflowContextFromRecord(ctx context.Context, record *storage.WorkflowRecord) *core.WorkflowContext {
	data := core.NewDefaultWorkflowDataWithMap(record.Data)
	metadata := core.NewDefaultWorkflowMetadataWithMap(record.Metadata)
	
	workflowCtx := core.NewWorkflowContext(ctx, record.ID, record.WorkflowName, data, metadata)
	workflowCtx.SetCurrentStep(record.CurrentStep)
	workflowCtx.SetStatus(record.Status)
	
	if record.NextStep != nil {
		workflowCtx.SetNextStep(*record.NextStep)
	}
	
	for stepName, result := range record.StepResults {
		workflowCtx.SetStepResult(stepName, result)
	}
	
	return workflowCtx
}

func (rm *WorkflowRecoveryManager) retryWorkflow(ctx context.Context, workflowCtx *core.WorkflowContext, policy *RecoveryPolicy) error {
	// Get retry count from metadata
	retryCount, _ := workflowCtx.Metadata.GetExecutionMetric("recovery_retry_count")
	count, ok := retryCount.(int)
	if !ok {
		count = 0
	}
	
	if count >= policy.MaxRetries {
		return fmt.Errorf("maximum retry attempts (%d) exceeded", policy.MaxRetries)
	}
	
	// Calculate delay with backoff
	delay := time.Duration(float64(policy.RetryDelay) * float64(count) * policy.BackoffFactor)
	if delay > policy.MaxDelay {
		delay = policy.MaxDelay
	}
	
	// Wait before retry
	if delay > 0 {
		time.Sleep(delay)
	}
	
	// Update retry count
	workflowCtx.Metadata.SetExecutionMetric("recovery_retry_count", count+1)
	workflowCtx.SetStatus(core.WorkflowStatusRunning)
	
	// Execute the failed step again
	nextStep, err := rm.engine.ExecuteStep(workflowCtx, workflowCtx.GetCurrentStep())
	if err != nil {
		return fmt.Errorf("retry failed: %w", err)
	}
	
	if nextStep != nil {
		workflowCtx.SetNextStep(*nextStep)
	} else {
		workflowCtx.Complete()
	}
	
	return rm.updateWorkflowRecord(ctx, workflowCtx)
}

func (rm *WorkflowRecoveryManager) skipFailedStep(ctx context.Context, workflowCtx *core.WorkflowContext) error {
	// Move to next step if available
	nextStep := workflowCtx.GetNextStep()
	if nextStep == nil {
		// No next step, complete the workflow
		workflowCtx.Complete()
	} else {
		workflowCtx.SetCurrentStep(*nextStep)
		workflowCtx.ClearNextStep()
		workflowCtx.SetStatus(core.WorkflowStatusRunning)
	}
	
	return rm.updateWorkflowRecord(ctx, workflowCtx)
}

func (rm *WorkflowRecoveryManager) markWorkflowFailed(ctx context.Context, workflowCtx *core.WorkflowContext) error {
	workflowCtx.SetStatus(core.WorkflowStatusFailed)
	return rm.updateWorkflowRecord(ctx, workflowCtx)
}

func (rm *WorkflowRecoveryManager) restartWorkflow(ctx context.Context, workflowCtx *core.WorkflowContext) error {
	// Reset workflow to initial state
	workflowCtx.SetStatus(core.WorkflowStatusPending)
	workflowCtx.ClearNextStep()
	
	// Clear step results and reset metadata
	for stepName := range workflowCtx.GetAllStepResults() {
		workflowCtx.SetStepResult(stepName, nil)
	}
	
	restartCount, _ := workflowCtx.Metadata.GetExecutionMetric("restart_count")
	count, ok := restartCount.(int)
	if !ok {
		count = 0
	}
	workflowCtx.Metadata.SetExecutionMetric("restart_count", count + 1)
	
	return rm.updateWorkflowRecord(ctx, workflowCtx)
}

func (rm *WorkflowRecoveryManager) executeCustomRecovery(ctx context.Context, workflowCtx *core.WorkflowContext, policy *RecoveryPolicy) error {
	if policy.CustomHandler == nil {
		return fmt.Errorf("custom recovery handler not set")
	}
	
	action, err := policy.CustomHandler.Handle(ctx, workflowCtx, workflowCtx.GetError())
	if err != nil {
		return fmt.Errorf("custom recovery handler failed: %w", err)
	}
	
	// Apply recovery action
	switch action.Action {
	case RecoveryStrategyRetry:
		return rm.retryWorkflow(ctx, workflowCtx, policy)
	case RecoveryStrategySkip:
		return rm.skipFailedStep(ctx, workflowCtx)
	case RecoveryStrategyFail:
		return rm.markWorkflowFailed(ctx, workflowCtx)
	case RecoveryStrategyRestart:
		return rm.restartWorkflow(ctx, workflowCtx)
	default:
		return fmt.Errorf("unknown recovery action: %s", action.Action)
	}
}

func (rm *WorkflowRecoveryManager) updateWorkflowRecord(ctx context.Context, workflowCtx *core.WorkflowContext) error {
	record := &storage.WorkflowRecord{
		ID:           workflowCtx.GetWorkflowID(),
		WorkflowName: workflowCtx.GetWorkflowName(),
		CurrentStep:  workflowCtx.GetCurrentStep(),
		NextStep:     workflowCtx.GetNextStep(),
		Status:       workflowCtx.GetStatus(),
		Data:         workflowCtx.Data.GetAll(),
		Metadata:     workflowCtx.Metadata.GetExecutionMetrics(),
		StepResults:  workflowCtx.GetAllStepResults(),
		UpdatedAt:    time.Now(),
	}
	
	if workflowCtx.GetError() != nil {
		errorMsg := workflowCtx.GetError().Error()
		record.Error = &errorMsg
	}
	
	return rm.storage.UpdateWorkflowRecord(ctx, record)
}

func (rm *WorkflowRecoveryManager) monitorFailedWorkflows(ctx context.Context) {
	// Get failed workflows
	filter := &storage.WorkflowFilter{
		Status: &[]core.WorkflowStatus{core.WorkflowStatusFailed}[0],
		Limit:  &[]int{100}[0],
	}
	
	records, err := rm.storage.ListWorkflowRecords(ctx, filter)
	if err != nil {
		rm.logger.Error("Failed to get failed workflows", map[string]interface{}{"error": err.Error()})
		return
	}
	
	// Attempt recovery for each failed workflow
	for _, record := range records {
		// Check if recovery should be attempted
		if rm.shouldAttemptRecovery(record) {
			if err := rm.RecoverWorkflow(ctx, record.ID); err != nil {
				rm.logger.Error("Auto-recovery failed", map[string]interface{}{"workflow_id": record.ID, "error": err.Error()})
			}
		}
	}
}

func (rm *WorkflowRecoveryManager) shouldAttemptRecovery(record *storage.WorkflowRecord) bool {
	// Check if enough time has passed since last failure
	minInterval := 5 * time.Minute
	if time.Since(record.UpdatedAt) < minInterval {
		return false
	}
	
	// Check recovery history
	history := rm.GetRecoveryHistory(record.ID)
	if len(history) > 0 {
		lastAttempt := history[len(history)-1]
		if time.Since(lastAttempt.AttemptTime) < minInterval {
			return false
		}
	}
	
	return true
}

func (rm *WorkflowRecoveryManager) addRecoveryAttempt(workflowID string, attempt *RecoveryAttempt) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if _, exists := rm.recoveryHistory[workflowID]; !exists {
		rm.recoveryHistory[workflowID] = make([]*RecoveryAttempt, 0)
	}
	
	rm.recoveryHistory[workflowID] = append(rm.recoveryHistory[workflowID], attempt)
	
	// Keep only last 10 attempts
	if len(rm.recoveryHistory[workflowID]) > 10 {
		rm.recoveryHistory[workflowID] = rm.recoveryHistory[workflowID][1:]
	}
}