package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Common response types
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Timestamp  time.Time   `json:"timestamp"`
}

type ErrorResponse struct {
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      string    `json:"code,omitempty"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Workflow request types
type CreateWorkflowRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Definition  map[string]interface{} `json:"definition" binding:"required"`
	InputSchema map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`
	Config      map[string]interface{} `json:"config"`
	Tags        []string               `json:"tags"`
	IsActive    bool                   `json:"is_active"`
}

type UpdateWorkflowRequest struct {
	Name        *string                 `json:"name"`
	Description *string                 `json:"description"`
	Definition  *map[string]interface{} `json:"definition"`
	InputSchema *map[string]interface{} `json:"input_schema"`
	OutputSchema *map[string]interface{} `json:"output_schema"`
	Config      *map[string]interface{} `json:"config"`
	Tags        *[]string               `json:"tags"`
	IsActive    *bool                   `json:"is_active"`
}

type ExecuteWorkflowRequest struct {
	Input       map[string]interface{} `json:"input"`
	Config      map[string]interface{} `json:"config"`
	TriggerType string                 `json:"trigger_type"`
	TriggerData map[string]interface{} `json:"trigger_data"`
	Async       bool                   `json:"async"`
	Timeout     *int                   `json:"timeout"`
	RetryPolicy map[string]interface{} `json:"retry_policy"`
}

// Code generation request types
type CodeGenRequest struct {
	WorkflowID uuid.UUID              `json:"workflow_id" binding:"required"`
	Language   string                 `json:"language" binding:"required"`
	Template   string                 `json:"template"`
	Options    map[string]interface{} `json:"options"`
}

type CodeGenJob struct {
	ID          uuid.UUID              `json:"id"`
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	Language    string                 `json:"language"`
	Template    string                 `json:"template"`
	Status      string                 `json:"status"`
	Progress    int                    `json:"progress"`
	Error       string                 `json:"error,omitempty"`
	ErrorCode   string                 `json:"error_code,omitempty"`
	Artifacts   []string               `json:"artifacts,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Options     map[string]interface{} `json:"options"`
}

func (j *CodeGenJob) IsCompleted() bool {
	return j.Status == "completed"
}

func (j *CodeGenJob) IsFailed() bool {
	return j.Status == "failed"
}

func (j *CodeGenJob) GetDurationSeconds() int64 {
	if j.CompletedAt == nil {
		return 0
	}
	return int64(j.CompletedAt.Sub(j.CreatedAt).Seconds())
}

type CodeGenTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Language    string                 `json:"language"`
	Category    string                 `json:"category"`
	Version     string                 `json:"version"`
	Author      string                 `json:"author"`
	Tags        []string               `json:"tags"`
	Options     map[string]interface{} `json:"options"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Metrics request types
type RecordMetricRequest struct {
	Name       string                 `json:"name" binding:"required"`
	Value      float64                `json:"value" binding:"required"`
	Type       string                 `json:"type" binding:"required"`
	Labels     map[string]string      `json:"labels"`
	Metadata   map[string]interface{} `json:"metadata"`
	Timestamp  *time.Time             `json:"timestamp"`
}

type MetricAggregationRequest struct {
	MetricName   string            `json:"metric_name" binding:"required"`
	Aggregation  string            `json:"aggregation" binding:"required"`
	GroupBy      []string          `json:"group_by"`
	Filters      map[string]string `json:"filters"`
	StartTime    time.Time         `json:"start_time" binding:"required"`
	EndTime      time.Time         `json:"end_time" binding:"required"`
	Interval     string            `json:"interval"`
}

// Alert request types
type CreateAlertRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Type        string                 `json:"type" binding:"required"`
	Condition   map[string]interface{} `json:"condition" binding:"required"`
	Actions     []map[string]interface{} `json:"actions" binding:"required"`
	Severity    string                 `json:"severity"`
	Tags        []string               `json:"tags"`
	Config      map[string]interface{} `json:"config"`
	IsEnabled   bool                   `json:"is_enabled"`
}

type UpdateAlertRequest struct {
	Name        *string                  `json:"name"`
	Description *string                  `json:"description"`
	Type        *string                  `json:"type"`
	Condition   *map[string]interface{}  `json:"condition"`
	Actions     *[]map[string]interface{} `json:"actions"`
	Severity    *string                  `json:"severity"`
	Tags        *[]string                `json:"tags"`
	Config      *map[string]interface{}  `json:"config"`
	IsEnabled   *bool                    `json:"is_enabled"`
}

// Dashboard request types
type CreateDashboardRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Description string                   `json:"description"`
	Layout      map[string]interface{}   `json:"layout" binding:"required"`
	Widgets     []map[string]interface{} `json:"widgets" binding:"required"`
	Filters     map[string]interface{}   `json:"filters"`
	Settings    map[string]interface{}   `json:"settings"`
	IsPublic    bool                     `json:"is_public"`
	Tags        []string                 `json:"tags"`
}

type UpdateDashboardRequest struct {
	Name        *string                   `json:"name"`
	Description *string                   `json:"description"`
	Layout      *map[string]interface{}   `json:"layout"`
	Widgets     *[]map[string]interface{} `json:"widgets"`
	Filters     *map[string]interface{}   `json:"filters"`
	Settings    *map[string]interface{}   `json:"settings"`
	IsPublic    *bool                     `json:"is_public"`
	Tags        *[]string                 `json:"tags"`
}

// Version request types
type CreateVersionRequest struct {
	Version         string                 `json:"version" binding:"required"`
	Description     string                 `json:"description"`
	Changelog       string                 `json:"changelog"`
	BreakingChanges bool                   `json:"breaking_changes"`
	Definition      map[string]interface{} `json:"definition"`
	InputSchema     map[string]interface{} `json:"input_schema"`
	OutputSchema    map[string]interface{} `json:"output_schema"`
	Config          map[string]interface{} `json:"config"`
}

type RollbackVersionRequest struct {
	Reason        string `json:"reason"`
	ForceRollback bool   `json:"force_rollback"`
}

type DeployVersionRequest struct {
	Environment string                 `json:"environment" binding:"required"`
	Strategy    string                 `json:"strategy"`
	Config      map[string]interface{} `json:"config"`
}

// Response types for dashboard data
type DashboardOverview struct {
	TotalWorkflows      int64                  `json:"total_workflows"`
	ActiveWorkflows     int64                  `json:"active_workflows"`
	TotalExecutions     int64                  `json:"total_executions"`
	SuccessfulExecutions int64                 `json:"successful_executions"`
	FailedExecutions    int64                  `json:"failed_executions"`
	RunningExecutions   int64                  `json:"running_executions"`
	AverageExecutionTime float64               `json:"average_execution_time"`
	Throughput          float64                `json:"throughput"`
	ErrorRate           float64                `json:"error_rate"`
	SystemHealth        string                 `json:"system_health"`
	RecentActivity      []map[string]interface{} `json:"recent_activity"`
	TopWorkflows        []map[string]interface{} `json:"top_workflows"`
	Alerts              []map[string]interface{} `json:"alerts"`
	Timestamp           time.Time              `json:"timestamp"`
}

type WorkflowStatusSummary struct {
	WorkflowID    uuid.UUID `json:"workflow_id"`
	WorkflowName  string    `json:"workflow_name"`
	TotalRuns     int64     `json:"total_runs"`
	SuccessfulRuns int64    `json:"successful_runs"`
	FailedRuns    int64     `json:"failed_runs"`
	RunningRuns   int64     `json:"running_runs"`
	AverageTime   float64   `json:"average_time"`
	LastRun       *time.Time `json:"last_run"`
	Status        string    `json:"status"`
}

type SystemHealth struct {
	OverallStatus   string                 `json:"overall_status"`
	Components      map[string]interface{} `json:"components"`
	CPUUsage        float64                `json:"cpu_usage"`
	MemoryUsage     float64                `json:"memory_usage"`
	DiskUsage       float64                `json:"disk_usage"`
	DatabaseStatus  string                 `json:"database_status"`
	CacheStatus     string                 `json:"cache_status"`
	QueueStatus     string                 `json:"queue_status"`
	ActiveConnections int64                `json:"active_connections"`
	Uptime          int64                  `json:"uptime"`
	Version         string                 `json:"version"`
	Timestamp       time.Time              `json:"timestamp"`
}

type LiveMetrics struct {
	ExecutionsPerSecond float64                `json:"executions_per_second"`
	ErrorsPerSecond     float64                `json:"errors_per_second"`
	AverageLatency      float64                `json:"average_latency"`
	Throughput          float64                `json:"throughput"`
	ActiveExecutions    int64                  `json:"active_executions"`
	QueuedExecutions    int64                  `json:"queued_executions"`
	CustomMetrics       map[string]interface{} `json:"custom_metrics"`
	Timestamp           time.Time              `json:"timestamp"`
}

// Pagination helpers
type PaginationParams struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// Time range helpers
type TimeRangeParams struct {
	StartTime *time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   *time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
}

func (t *TimeRangeParams) GetTimeRange() (time.Time, time.Time, error) {
	var start, end time.Time

	if t.StartTime != nil {
		start = *t.StartTime
	} else {
		// Default to 24 hours ago
		start = time.Now().UTC().Add(-24 * time.Hour)
	}

	if t.EndTime != nil {
		end = *t.EndTime
	} else {
		// Default to now
		end = time.Now().UTC()
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("start_time cannot be after end_time")
	}

	return start, end, nil
}