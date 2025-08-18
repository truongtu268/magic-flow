package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExecutionStatus represents the status of a workflow execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
	ExecutionStatusTimeout   ExecutionStatus = "timeout"
	ExecutionStatusPaused    ExecutionStatus = "paused"
)

// StepStatus represents the status of a workflow step execution
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
	StepStatusRetrying  StepStatus = "retrying"
)

// TriggerType represents how the execution was triggered
type TriggerType string

const (
	TriggerTypeManual    TriggerType = "manual"
	TriggerTypeAPI       TriggerType = "api"
	TriggerTypeScheduled TriggerType = "scheduled"
	TriggerTypeWebhook   TriggerType = "webhook"
	TriggerTypeEvent     TriggerType = "event"
)

// Execution represents a workflow execution instance
type Execution struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkflowID uuid.UUID       `json:"workflow_id" gorm:"type:uuid;not null;index"`
	Status     ExecutionStatus `json:"status" gorm:"default:'pending';index"`
	
	// Trigger information
	TriggerType TriggerType            `json:"trigger_type" gorm:"not null"`
	TriggerBy   string                 `json:"trigger_by"`
	TriggerData map[string]interface{} `json:"trigger_data" gorm:"type:jsonb"`
	
	// Input and output data
	InputData  map[string]interface{} `json:"input_data" gorm:"type:jsonb"`
	OutputData map[string]interface{} `json:"output_data" gorm:"type:jsonb"`
	
	// Execution context
	Context ExecutionContext `json:"context" gorm:"type:jsonb"`
	
	// Timing information
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Duration    int64      `json:"duration"` // Duration in milliseconds
	
	// Error information
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	
	// Metadata
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Workflow  Workflow       `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	Steps     []StepExecution `json:"steps,omitempty" gorm:"foreignKey:ExecutionID"`
	Events    []ExecutionEvent `json:"events,omitempty" gorm:"foreignKey:ExecutionID"`
}

// ExecutionContext represents the execution context
type ExecutionContext struct {
	ExecutionID   string                 `json:"execution_id"`
	WorkflowName  string                 `json:"workflow_name"`
	WorkflowVersion string               `json:"workflow_version"`
	Environment   map[string]string      `json:"environment,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	Secrets       map[string]string      `json:"secrets,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	TraceID       string                 `json:"trace_id,omitempty"`
}

// StepExecution represents the execution of a single workflow step
type StepExecution struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExecutionID uuid.UUID  `json:"execution_id" gorm:"type:uuid;not null;index"`
	StepName    string     `json:"step_name" gorm:"not null;index"`
	StepType    string     `json:"step_type" gorm:"not null"`
	Status      StepStatus `json:"status" gorm:"default:'pending'"`
	
	// Input and output data
	InputData  map[string]interface{} `json:"input_data" gorm:"type:jsonb"`
	OutputData map[string]interface{} `json:"output_data" gorm:"type:jsonb"`
	
	// Timing information
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Duration    int64      `json:"duration"` // Duration in milliseconds
	
	// Retry information
	Attempt     int `json:"attempt" gorm:"default:1"`
	MaxAttempts int `json:"max_attempts" gorm:"default:1"`
	
	// Error information
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	
	// Metadata
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Execution Execution `json:"-" gorm:"foreignKey:ExecutionID"`
}

// ExecutionEvent represents an event during workflow execution
type ExecutionEvent struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExecutionID uuid.UUID `json:"execution_id" gorm:"type:uuid;not null;index"`
	EventType   string    `json:"event_type" gorm:"not null;index"`
	StepName    string    `json:"step_name,omitempty" gorm:"index"`
	
	// Event data
	Data map[string]interface{} `json:"data" gorm:"type:jsonb"`
	
	// Timing
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	
	// Metadata
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Relationships
	Execution Execution `json:"-" gorm:"foreignKey:ExecutionID"`
}

// ExecutionMetrics represents execution metrics
type ExecutionMetrics struct {
	TotalExecutions     int64   `json:"total_executions"`
	SuccessfulExecutions int64  `json:"successful_executions"`
	FailedExecutions    int64   `json:"failed_executions"`
	AverageDuration     float64 `json:"average_duration"`
	SuccessRate         float64 `json:"success_rate"`
	Throughput          float64 `json:"throughput"`
}

// BeforeCreate sets the ID before creating
func (e *Execution) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// BeforeCreate sets the ID before creating
func (se *StepExecution) BeforeCreate(tx *gorm.DB) error {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return nil
}

// BeforeCreate sets the ID before creating
func (ee *ExecutionEvent) BeforeCreate(tx *gorm.DB) error {
	if ee.ID == uuid.Nil {
		ee.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Execution model
func (Execution) TableName() string {
	return "executions"
}

// TableName returns the table name for the StepExecution model
func (StepExecution) TableName() string {
	return "step_executions"
}

// TableName returns the table name for the ExecutionEvent model
func (ExecutionEvent) TableName() string {
	return "execution_events"
}

// Start marks the execution as started
func (e *Execution) Start() {
	now := time.Now()
	e.Status = ExecutionStatusRunning
	e.StartedAt = &now
}

// Complete marks the execution as completed
func (e *Execution) Complete(outputData map[string]interface{}) {
	now := time.Now()
	e.Status = ExecutionStatusCompleted
	e.CompletedAt = &now
	e.OutputData = outputData
	
	if e.StartedAt != nil {
		e.Duration = now.Sub(*e.StartedAt).Milliseconds()
	}
}

// Fail marks the execution as failed
func (e *Execution) Fail(err error, errorCode string) {
	now := time.Now()
	e.Status = ExecutionStatusFailed
	e.CompletedAt = &now
	e.Error = err.Error()
	e.ErrorCode = errorCode
	
	if e.StartedAt != nil {
		e.Duration = now.Sub(*e.StartedAt).Milliseconds()
	}
}

// Cancel marks the execution as cancelled
func (e *Execution) Cancel() {
	now := time.Now()
	e.Status = ExecutionStatusCancelled
	e.CompletedAt = &now
	
	if e.StartedAt != nil {
		e.Duration = now.Sub(*e.StartedAt).Milliseconds()
	}
}

// IsRunning returns true if the execution is running
func (e *Execution) IsRunning() bool {
	return e.Status == ExecutionStatusRunning
}

// IsCompleted returns true if the execution is completed
func (e *Execution) IsCompleted() bool {
	return e.Status == ExecutionStatusCompleted
}

// IsFailed returns true if the execution failed
func (e *Execution) IsFailed() bool {
	return e.Status == ExecutionStatusFailed
}

// IsCancelled returns true if the execution was cancelled
func (e *Execution) IsCancelled() bool {
	return e.Status == ExecutionStatusCancelled
}

// IsFinished returns true if the execution is in a terminal state
func (e *Execution) IsFinished() bool {
	return e.IsCompleted() || e.IsFailed() || e.IsCancelled() || e.Status == ExecutionStatusTimeout
}

// GetDurationSeconds returns the duration in seconds
func (e *Execution) GetDurationSeconds() float64 {
	return float64(e.Duration) / 1000.0
}

// Start marks the step execution as started
func (se *StepExecution) Start() {
	now := time.Now()
	se.Status = StepStatusRunning
	se.StartedAt = &now
}

// Complete marks the step execution as completed
func (se *StepExecution) Complete(outputData map[string]interface{}) {
	now := time.Now()
	se.Status = StepStatusCompleted
	se.CompletedAt = &now
	se.OutputData = outputData
	
	if se.StartedAt != nil {
		se.Duration = now.Sub(*se.StartedAt).Milliseconds()
	}
}

// Fail marks the step execution as failed
func (se *StepExecution) Fail(err error, errorCode string) {
	now := time.Now()
	se.Status = StepStatusFailed
	se.CompletedAt = &now
	se.Error = err.Error()
	se.ErrorCode = errorCode
	
	if se.StartedAt != nil {
		se.Duration = now.Sub(*se.StartedAt).Milliseconds()
	}
}

// Skip marks the step execution as skipped
func (se *StepExecution) Skip(reason string) {
	now := time.Now()
	se.Status = StepStatusSkipped
	se.CompletedAt = &now
	se.Error = reason
	
	if se.StartedAt != nil {
		se.Duration = now.Sub(*se.StartedAt).Milliseconds()
	}
}

// Retry marks the step execution for retry
func (se *StepExecution) Retry() {
	se.Status = StepStatusRetrying
	se.Attempt++
}

// CanRetry returns true if the step can be retried
func (se *StepExecution) CanRetry() bool {
	return se.Attempt < se.MaxAttempts
}

// IsRunning returns true if the step is running
func (se *StepExecution) IsRunning() bool {
	return se.Status == StepStatusRunning
}

// IsCompleted returns true if the step is completed
func (se *StepExecution) IsCompleted() bool {
	return se.Status == StepStatusCompleted
}

// IsFailed returns true if the step failed
func (se *StepExecution) IsFailed() bool {
	return se.Status == StepStatusFailed
}

// IsSkipped returns true if the step was skipped
func (se *StepExecution) IsSkipped() bool {
	return se.Status == StepStatusSkipped
}

// IsFinished returns true if the step is in a terminal state
func (se *StepExecution) IsFinished() bool {
	return se.IsCompleted() || se.IsFailed() || se.IsSkipped()
}

// GetDurationSeconds returns the duration in seconds
func (se *StepExecution) GetDurationSeconds() float64 {
	return float64(se.Duration) / 1000.0
}

// ToJSON converts the execution to JSON
func (e *Execution) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON populates the execution from JSON
func (e *Execution) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// AddEvent adds an event to the execution
func (e *Execution) AddEvent(eventType, stepName string, data map[string]interface{}) *ExecutionEvent {
	event := &ExecutionEvent{
		ID:          uuid.New(),
		ExecutionID: e.ID,
		EventType:   eventType,
		StepName:    stepName,
		Data:        data,
		Timestamp:   time.Now(),
	}
	
	e.Events = append(e.Events, *event)
	return event
}

// GetStepExecution returns the step execution by step name
func (e *Execution) GetStepExecution(stepName string) *StepExecution {
	for i := range e.Steps {
		if e.Steps[i].StepName == stepName {
			return &e.Steps[i]
		}
	}
	return nil
}

// GetCompletedSteps returns all completed steps
func (e *Execution) GetCompletedSteps() []StepExecution {
	var completed []StepExecution
	for _, step := range e.Steps {
		if step.IsCompleted() {
			completed = append(completed, step)
		}
	}
	return completed
}

// GetFailedSteps returns all failed steps
func (e *Execution) GetFailedSteps() []StepExecution {
	var failed []StepExecution
	for _, step := range e.Steps {
		if step.IsFailed() {
			failed = append(failed, step)
		}
	}
	return failed
}

// GetProgress returns the execution progress as a percentage
func (e *Execution) GetProgress() float64 {
	if len(e.Steps) == 0 {
		return 0.0
	}
	
	finished := 0
	for _, step := range e.Steps {
		if step.IsFinished() {
			finished++
		}
	}
	
	return float64(finished) / float64(len(e.Steps)) * 100.0
}