package storage

import (
	"context"
	"time"

	"github.com/truongtu268/magic-flow/pkg/events"
)

// WorkflowRecord represents a stored workflow record
type WorkflowRecord struct {
	ID           string                 `json:"id" db:"id"`
	WorkflowName string                 `json:"workflow_name" db:"workflow_name"`
	CurrentStep  string                 `json:"current_step" db:"current_step"`
	NextStep     *string                `json:"next_step" db:"next_step"`
	Status       events.WorkflowStatus    `json:"status" db:"status"`
	Data         map[string]interface{} `json:"data" db:"data"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	StepResults  map[string]interface{} `json:"step_results" db:"step_results"`
	StartTime    time.Time              `json:"start_time" db:"start_time"`
	EndTime      *time.Time             `json:"end_time" db:"end_time"`
	Error        *string                `json:"error" db:"error"`
	StepOrder    int                    `json:"step_order" db:"step_order"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// WorkflowStorage defines the interface for persisting workflow state
type WorkflowStorage interface {
	// CreateWorkflowRecord creates a new workflow record
	CreateWorkflowRecord(ctx context.Context, record *WorkflowRecord) error
	// GetWorkflowRecord retrieves a workflow record by ID
	GetWorkflowRecord(ctx context.Context, workflowID string) (*WorkflowRecord, error)
	// UpdateWorkflowRecord updates an existing workflow record
	UpdateWorkflowRecord(ctx context.Context, record *WorkflowRecord) error
	// DeleteWorkflowRecord deletes a workflow record
	DeleteWorkflowRecord(ctx context.Context, workflowID string) error
	// ListWorkflowRecords lists workflow records with optional filtering
	ListWorkflowRecords(ctx context.Context, filter *WorkflowFilter) ([]*WorkflowRecord, error)
	// GetWaitingWorkflows retrieves workflows in waiting status
	GetWaitingWorkflows(ctx context.Context) ([]*WorkflowRecord, error)
	// GetWaitingWorkflowsByTrigger retrieves workflows waiting for a specific trigger
	GetWaitingWorkflowsByTrigger(ctx context.Context, triggerKey string) ([]*WorkflowRecord, error)
	// Close closes the storage connection
	Close() error
}

// EnhancedWorkflowStorage extends WorkflowStorage with additional features
type EnhancedWorkflowStorage interface {
	WorkflowStorage
	// GetWorkflowsByStatus gets workflows by status
	GetWorkflowsByStatus(ctx context.Context, status events.WorkflowStatus) ([]*WorkflowRecord, error)
	// GetWorkflowHistory retrieves workflow execution history
	GetWorkflowHistory(ctx context.Context, workflowID string) ([]*WorkflowRecord, error)
	// BulkUpdateWorkflows updates multiple workflows in a transaction
	BulkUpdateWorkflows(ctx context.Context, records []*WorkflowRecord) error
	// GetWorkflowStats returns workflow statistics
	GetWorkflowStats(ctx context.Context) (*WorkflowStats, error)
	// ArchiveCompletedWorkflows archives old completed workflows
	ArchiveCompletedWorkflows(ctx context.Context, olderThan time.Time) (int, error)
}

// WorkflowFilter defines filtering options for listing workflows
type WorkflowFilter struct {
	Status       *events.WorkflowStatus `json:"status,omitempty"`
	WorkflowName *string              `json:"workflow_name,omitempty"`
	StartTime    *time.Time           `json:"start_time,omitempty"`
	EndTime      *time.Time           `json:"end_time,omitempty"`
	Limit        *int                 `json:"limit,omitempty"`
	Offset       *int                 `json:"offset,omitempty"`
	OrderBy      *string              `json:"order_by,omitempty"`
	OrderDesc    bool                 `json:"order_desc,omitempty"`
}

// WorkflowStats represents workflow statistics
type WorkflowStats struct {
	TotalWorkflows     int64                          `json:"total_workflows"`
	StatusCounts       map[events.WorkflowStatus]int64  `json:"status_counts"`
	WorkflowNameCounts map[string]int64               `json:"workflow_name_counts"`
	AverageExecutionTime time.Duration               `json:"average_execution_time"`
	LastUpdated        time.Time                     `json:"last_updated"`
}

// BackgroundJob represents a background job
type BackgroundJob struct {
	ID          string                 `json:"id" db:"id"`
	Type        string                 `json:"type" db:"type"`
	Payload     map[string]interface{} `json:"payload" db:"payload"`
	Status      events.JobStatus         `json:"status" db:"status"`
	ScheduledAt time.Time              `json:"scheduled_at" db:"scheduled_at"`
	StartedAt   *time.Time             `json:"started_at" db:"started_at"`
	CompletedAt *time.Time             `json:"completed_at" db:"completed_at"`
	Error       *string                `json:"error" db:"error"`
	RetryCount  int                    `json:"retry_count" db:"retry_count"`
	MaxRetries  int                    `json:"max_retries" db:"max_retries"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// BackgroundJobProcessor defines the interface for processing background jobs
type BackgroundJobProcessor interface {
	// EnqueueJob adds a job to the processing queue
	EnqueueJob(ctx context.Context, job *BackgroundJob) error
	// ProcessJobs starts processing jobs from the queue
	ProcessJobs(ctx context.Context) error
	// GetJob retrieves a job by ID
	GetJob(ctx context.Context, jobID string) (*BackgroundJob, error)
	// UpdateJobStatus updates the status of a job
	UpdateJobStatus(ctx context.Context, jobID string, status events.JobStatus, error *string) error
	// GetPendingJobs retrieves pending jobs
	GetPendingJobs(ctx context.Context, limit int) ([]*BackgroundJob, error)
	// RegisterJobHandler registers a handler for a specific job type
	RegisterJobHandler(jobType string, handler JobHandler) error
	// Stop stops the job processor
	Stop() error
}

// JobHandler defines the interface for handling specific job types
type JobHandler interface {
	// Handle processes a background job
	Handle(ctx context.Context, job *BackgroundJob) error
	// GetJobType returns the job type this handler processes
	GetJobType() string
}

// TransactionManager defines the interface for managing database transactions
type TransactionManager interface {
	// BeginTransaction starts a new transaction
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction defines the interface for database transactions
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
	// GetContext returns the transaction context
	GetContext() context.Context
}

// StorageConfig defines configuration for storage implementations
type StorageConfig struct {
	DatabaseURL      string        `json:"database_url"`
	MaxConnections   int           `json:"max_connections"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	QueryTimeout     time.Duration `json:"query_timeout"`
	RetryAttempts    int           `json:"retry_attempts"`
	RetryDelay       time.Duration `json:"retry_delay"`
	EnableMetrics    bool          `json:"enable_metrics"`
	TablePrefix      string        `json:"table_prefix"`
}

// DefaultStorageConfig returns a default storage configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		DatabaseURL:       "postgres://localhost:5432/magic_flow?sslmode=disable",
		MaxConnections:    10,
		ConnectionTimeout: 30 * time.Second,
		QueryTimeout:      30 * time.Second,
		RetryAttempts:     3,
		RetryDelay:        1 * time.Second,
		EnableMetrics:     true,
		TablePrefix:       "magic_flow_",
	}
}