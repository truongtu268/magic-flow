package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// MetricsCollector handles metrics collection and aggregation
type MetricsCollector struct {
	repoManager database.RepositoryManager
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(repoManager database.RepositoryManager) *MetricsCollector {
	return &MetricsCollector{
		repoManager: repoManager,
	}
}

// WorkflowMetrics represents metrics for a specific workflow
type WorkflowMetrics struct {
	WorkflowID       uuid.UUID                `json:"workflow_id"`
	WorkflowName     string                   `json:"workflow_name"`
	TotalExecutions  int64                    `json:"total_executions"`
	SuccessfulRuns   int64                    `json:"successful_runs"`
	FailedRuns       int64                    `json:"failed_runs"`
	AverageRuntime   time.Duration            `json:"average_runtime"`
	SuccessRate      float64                  `json:"success_rate"`
	LastExecution    *time.Time               `json:"last_execution"`
	ExecutionTrend   []ExecutionTrendPoint    `json:"execution_trend"`
	PerformanceTrend []PerformanceTrendPoint  `json:"performance_trend"`
	ErrorBreakdown   map[string]int64         `json:"error_breakdown"`
	StepMetrics      []StepMetrics            `json:"step_metrics"`
	TimeRange        string                   `json:"time_range"`
	GeneratedAt      time.Time                `json:"generated_at"`
}

// ExecutionMetrics represents execution metrics with various filters
type ExecutionMetrics struct {
	TotalExecutions    int64                   `json:"total_executions"`
	RunningExecutions  int64                   `json:"running_executions"`
	CompletedExecutions int64                  `json:"completed_executions"`
	FailedExecutions   int64                   `json:"failed_executions"`
	CancelledExecutions int64                  `json:"cancelled_executions"`
	AverageRuntime     time.Duration           `json:"average_runtime"`
	MedianRuntime      time.Duration           `json:"median_runtime"`
	P95Runtime         time.Duration           `json:"p95_runtime"`
	P99Runtime         time.Duration           `json:"p99_runtime"`
	Throughput         float64                 `json:"throughput"`
	ErrorRate          float64                 `json:"error_rate"`
	ExecutionsByStatus map[string]int64        `json:"executions_by_status"`
	ExecutionsByHour   []HourlyExecutionCount  `json:"executions_by_hour"`
	TopFailedWorkflows []WorkflowFailureCount  `json:"top_failed_workflows"`
	TimeRange          string                  `json:"time_range"`
	GeneratedAt        time.Time               `json:"generated_at"`
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	TotalWorkflows     int64                    `json:"total_workflows"`
	ActiveWorkflows    int64                    `json:"active_workflows"`
	TotalExecutions    int64                    `json:"total_executions"`
	DailyExecutions    int64                    `json:"daily_executions"`
	WeeklyExecutions   int64                    `json:"weekly_executions"`
	MonthlyExecutions  int64                    `json:"monthly_executions"`
	SystemLoad         SystemLoadMetrics        `json:"system_load"`
	ResourceUsage      ResourceUsageMetrics     `json:"resource_usage"`
	APIMetrics         APIMetrics               `json:"api_metrics"`
	UserActivity       UserActivityMetrics      `json:"user_activity"`
	WorkflowTrends     []WorkflowTrendPoint     `json:"workflow_trends"`
	ExecutionTrends    []ExecutionTrendPoint    `json:"execution_trends"`
	TimeRange          string                   `json:"time_range"`
	GeneratedAt        time.Time                `json:"generated_at"`
}

// ExecutionTrendPoint represents a point in execution trend data
type ExecutionTrendPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Executions  int64     `json:"executions"`
	Successful  int64     `json:"successful"`
	Failed      int64     `json:"failed"`
	Cancelled   int64     `json:"cancelled"`
	AverageTime float64   `json:"average_time"`
}

// PerformanceTrendPoint represents a point in performance trend data
type PerformanceTrendPoint struct {
	Timestamp       time.Time `json:"timestamp"`
	AverageRuntime  float64   `json:"average_runtime"`
	MedianRuntime   float64   `json:"median_runtime"`
	P95Runtime      float64   `json:"p95_runtime"`
	Throughput      float64   `json:"throughput"`
	ErrorRate       float64   `json:"error_rate"`
}

// StepMetrics represents metrics for individual workflow steps
type StepMetrics struct {
	StepName        string        `json:"step_name"`
	StepType        string        `json:"step_type"`
	Executions      int64         `json:"executions"`
	Successful      int64         `json:"successful"`
	Failed          int64         `json:"failed"`
	AverageRuntime  time.Duration `json:"average_runtime"`
	SuccessRate     float64       `json:"success_rate"`
	CommonErrors    []string      `json:"common_errors"`
}

// HourlyExecutionCount represents execution count for a specific hour
type HourlyExecutionCount struct {
	Hour       time.Time `json:"hour"`
	Executions int64     `json:"executions"`
	Successful int64     `json:"successful"`
	Failed     int64     `json:"failed"`
}

// WorkflowFailureCount represents failure count for a workflow
type WorkflowFailureCount struct {
	WorkflowID   uuid.UUID `json:"workflow_id"`
	WorkflowName string    `json:"workflow_name"`
	FailureCount int64     `json:"failure_count"`
	ErrorRate    float64   `json:"error_rate"`
}

// SystemLoadMetrics represents system load metrics
type SystemLoadMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIO   float64 `json:"network_io"`
	LoadAverage float64 `json:"load_average"`
}

// ResourceUsageMetrics represents resource usage metrics
type ResourceUsageMetrics struct {
	ActiveConnections int64   `json:"active_connections"`
	DatabaseQueries   int64   `json:"database_queries"`
	CacheHitRate      float64 `json:"cache_hit_rate"`
	QueueSize         int64   `json:"queue_size"`
	WorkerUtilization float64 `json:"worker_utilization"`
}

// APIMetrics represents API usage metrics
type APIMetrics struct {
	TotalRequests    int64            `json:"total_requests"`
	RequestsPerHour  float64          `json:"requests_per_hour"`
	AverageLatency   time.Duration    `json:"average_latency"`
	P95Latency       time.Duration    `json:"p95_latency"`
	ErrorRate        float64          `json:"error_rate"`
	EndpointMetrics  []EndpointMetric `json:"endpoint_metrics"`
}

// EndpointMetric represents metrics for a specific API endpoint
type EndpointMetric struct {
	Endpoint       string        `json:"endpoint"`
	Method         string        `json:"method"`
	Requests       int64         `json:"requests"`
	AverageLatency time.Duration `json:"average_latency"`
	ErrorRate      float64       `json:"error_rate"`
}

// UserActivityMetrics represents user activity metrics
type UserActivityMetrics struct {
	ActiveUsers       int64 `json:"active_users"`
	DailyActiveUsers  int64 `json:"daily_active_users"`
	WeeklyActiveUsers int64 `json:"weekly_active_users"`
	NewUsers          int64 `json:"new_users"`
	UserSessions      int64 `json:"user_sessions"`
}

// WorkflowTrendPoint represents a point in workflow trend data
type WorkflowTrendPoint struct {
	Timestamp        time.Time `json:"timestamp"`
	TotalWorkflows   int64     `json:"total_workflows"`
	ActiveWorkflows  int64     `json:"active_workflows"`
	NewWorkflows     int64     `json:"new_workflows"`
	UpdatedWorkflows int64     `json:"updated_workflows"`
}

// ExecutionMetricsFilters represents filters for execution metrics
type ExecutionMetricsFilters struct {
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Status     *string    `json:"status,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	TimeRange  string     `json:"time_range,omitempty"`
}

// GetWorkflowMetrics retrieves metrics for a specific workflow
func (mc *MetricsCollector) GetWorkflowMetrics(ctx context.Context, workflowID uuid.UUID, timeRange string) (*WorkflowMetrics, error) {
	workflowRepo := mc.repoManager.WorkflowRepository()
	executionRepo := mc.repoManager.ExecutionRepository()
	metricsRepo := mc.repoManager.MetricsRepository()

	// Get workflow details
	workflow, err := workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Parse time range
	startTime, endTime := mc.parseTimeRange(timeRange)

	// Get execution statistics
	stats, err := executionRepo.GetExecutionStatistics(ctx, &startTime, &endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution statistics: %w", err)
	}

	// Get execution trend data
	executionTrend, err := mc.getExecutionTrend(ctx, workflowID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution trend: %w", err)
	}

	// Get performance trend data
	performanceTrend, err := mc.getPerformanceTrend(ctx, workflowID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance trend: %w", err)
	}

	// Get error breakdown
	errorBreakdown, err := mc.getErrorBreakdown(ctx, workflowID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get error breakdown: %w", err)
	}

	// Get step metrics
	stepMetrics, err := mc.getStepMetrics(ctx, workflowID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get step metrics: %w", err)
	}

	// Get last execution time
	lastExecution, err := executionRepo.GetLastExecutionTime(ctx, workflowID)
	if err != nil && err.Error() != "record not found" {
		return nil, fmt.Errorf("failed to get last execution time: %w", err)
	}

	// Calculate success rate
	successRate := float64(0)
	if stats.TotalExecutions > 0 {
		successRate = float64(stats.SuccessfulExecutions) / float64(stats.TotalExecutions) * 100
	}

	return &WorkflowMetrics{
		WorkflowID:       workflowID,
		WorkflowName:     workflow.Name,
		TotalExecutions:  stats.TotalExecutions,
		SuccessfulRuns:   stats.SuccessfulExecutions,
		FailedRuns:       stats.FailedExecutions,
		AverageRuntime:   stats.AverageDuration,
		SuccessRate:      successRate,
		LastExecution:    lastExecution,
		ExecutionTrend:   executionTrend,
		PerformanceTrend: performanceTrend,
		ErrorBreakdown:   errorBreakdown,
		StepMetrics:      stepMetrics,
		TimeRange:        timeRange,
		GeneratedAt:      time.Now(),
	}, nil
}

// GetExecutionMetrics retrieves execution metrics with filters
func (mc *MetricsCollector) GetExecutionMetrics(ctx context.Context, filters ExecutionMetricsFilters) (*ExecutionMetrics, error) {
	executionRepo := mc.repoManager.ExecutionRepository()

	// Parse time range
	timeRange := filters.TimeRange
	if timeRange == "" {
		timeRange = "24h"
	}
	startTime, endTime := mc.parseTimeRange(timeRange)

	// Apply time filters if provided
	if filters.StartTime != nil {
		startTime = *filters.StartTime
	}
	if filters.EndTime != nil {
		endTime = *filters.EndTime
	}

	// Get execution statistics
	stats, err := executionRepo.GetExecutionStatistics(ctx, &startTime, &endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution statistics: %w", err)
	}

	// Get executions by status
	executionsByStatus := map[string]int64{
		"running":   stats.RunningExecutions,
		"completed": stats.SuccessfulExecutions,
		"failed":    stats.FailedExecutions,
		"cancelled": stats.CancelledExecutions,
	}

	// Get hourly execution counts
	executionsByHour, err := mc.getHourlyExecutionCounts(ctx, startTime, endTime, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get hourly execution counts: %w", err)
	}

	// Get top failed workflows
	topFailedWorkflows, err := mc.getTopFailedWorkflows(ctx, startTime, endTime, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top failed workflows: %w", err)
	}

	// Calculate throughput (executions per hour)
	duration := endTime.Sub(startTime).Hours()
	throughput := float64(stats.TotalExecutions) / duration

	// Calculate error rate
	errorRate := float64(0)
	if stats.TotalExecutions > 0 {
		errorRate = float64(stats.FailedExecutions) / float64(stats.TotalExecutions) * 100
	}

	return &ExecutionMetrics{
		TotalExecutions:     stats.TotalExecutions,
		RunningExecutions:   stats.RunningExecutions,
		CompletedExecutions: stats.SuccessfulExecutions,
		FailedExecutions:    stats.FailedExecutions,
		CancelledExecutions: stats.CancelledExecutions,
		AverageRuntime:      stats.AverageDuration,
		MedianRuntime:       stats.MedianDuration,
		P95Runtime:          stats.P95Duration,
		P99Runtime:          stats.P99Duration,
		Throughput:          throughput,
		ErrorRate:           errorRate,
		ExecutionsByStatus:  executionsByStatus,
		ExecutionsByHour:    executionsByHour,
		TopFailedWorkflows:  topFailedWorkflows,
		TimeRange:           timeRange,
		GeneratedAt:         time.Now(),
	}, nil
}

// GetSystemMetrics retrieves system-wide metrics
func (mc *MetricsCollector) GetSystemMetrics(ctx context.Context, timeRange string) (*SystemMetrics, error) {
	workflowRepo := mc.repoManager.WorkflowRepository()
	executionRepo := mc.repoManager.ExecutionRepository()

	// Parse time range
	startTime, endTime := mc.parseTimeRange(timeRange)

	// Get workflow counts
	totalWorkflows, err := workflowRepo.Count(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get total workflows: %w", err)
	}

	activeWorkflows, err := workflowRepo.CountByStatus(ctx, models.WorkflowStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get active workflows: %w", err)
	}

	// Get execution counts
	totalExecutions, err := executionRepo.Count(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get total executions: %w", err)
	}

	// Get time-based execution counts
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	dailyExecutions, err := executionRepo.CountByTimeRange(ctx, startOfDay, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily executions: %w", err)
	}

	weeklyExecutions, err := executionRepo.CountByTimeRange(ctx, startOfWeek, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly executions: %w", err)
	}

	monthlyExecutions, err := executionRepo.CountByTimeRange(ctx, startOfMonth, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly executions: %w", err)
	}

	// Get trend data
	workflowTrends, err := mc.getWorkflowTrends(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow trends: %w", err)
	}

	executionTrends, err := mc.getExecutionTrend(ctx, uuid.Nil, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution trends: %w", err)
	}

	// Mock system metrics (would be replaced with actual system monitoring)
	systemLoad := SystemLoadMetrics{
		CPUUsage:    45.2,
		MemoryUsage: 67.8,
		DiskUsage:   23.4,
		NetworkIO:   12.5,
		LoadAverage: 1.2,
	}

	resourceUsage := ResourceUsageMetrics{
		ActiveConnections: 25,
		DatabaseQueries:   1500,
		CacheHitRate:      85.5,
		QueueSize:         10,
		WorkerUtilization: 75.0,
	}

	apiMetrics := APIMetrics{
		TotalRequests:   5000,
		RequestsPerHour: 208.3,
		AverageLatency:  150 * time.Millisecond,
		P95Latency:      500 * time.Millisecond,
		ErrorRate:       2.5,
	}

	userActivity := UserActivityMetrics{
		ActiveUsers:       15,
		DailyActiveUsers:  25,
		WeeklyActiveUsers: 45,
		NewUsers:          3,
		UserSessions:      35,
	}

	return &SystemMetrics{
		TotalWorkflows:    totalWorkflows,
		ActiveWorkflows:   activeWorkflows,
		TotalExecutions:   totalExecutions,
		DailyExecutions:   dailyExecutions,
		WeeklyExecutions:  weeklyExecutions,
		MonthlyExecutions: monthlyExecutions,
		SystemLoad:        systemLoad,
		ResourceUsage:     resourceUsage,
		APIMetrics:        apiMetrics,
		UserActivity:      userActivity,
		WorkflowTrends:    workflowTrends,
		ExecutionTrends:   executionTrends,
		TimeRange:         timeRange,
		GeneratedAt:       time.Now(),
	}, nil
}

// parseTimeRange parses a time range string and returns start and end times
func (mc *MetricsCollector) parseTimeRange(timeRange string) (time.Time, time.Time) {
	now := time.Now()
	var startTime time.Time

	switch timeRange {
	case "1h":
		startTime = now.Add(-1 * time.Hour)
	case "6h":
		startTime = now.Add(-6 * time.Hour)
	case "24h", "1d":
		startTime = now.Add(-24 * time.Hour)
	case "7d", "1w":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "30d", "1m":
		startTime = now.Add(-30 * 24 * time.Hour)
	case "90d", "3m":
		startTime = now.Add(-90 * 24 * time.Hour)
	default:
		startTime = now.Add(-24 * time.Hour)
	}

	return startTime, now
}

// Helper methods for getting specific metrics data

func (mc *MetricsCollector) getExecutionTrend(ctx context.Context, workflowID uuid.UUID, startTime, endTime time.Time) ([]ExecutionTrendPoint, error) {
	// This would query execution data and aggregate by time intervals
	// For now, return mock data
	return []ExecutionTrendPoint{
		{
			Timestamp:   startTime,
			Executions:  10,
			Successful:  8,
			Failed:      2,
			Cancelled:   0,
			AverageTime: 120.5,
		},
	}, nil
}

func (mc *MetricsCollector) getPerformanceTrend(ctx context.Context, workflowID uuid.UUID, startTime, endTime time.Time) ([]PerformanceTrendPoint, error) {
	// This would query performance data and aggregate by time intervals
	// For now, return mock data
	return []PerformanceTrendPoint{
		{
			Timestamp:      startTime,
			AverageRuntime: 120.5,
			MedianRuntime:  115.0,
			P95Runtime:     200.0,
			Throughput:     5.2,
			ErrorRate:      2.5,
		},
	}, nil
}

func (mc *MetricsCollector) getErrorBreakdown(ctx context.Context, workflowID uuid.UUID, startTime, endTime time.Time) (map[string]int64, error) {
	// This would query error data and aggregate by error type
	// For now, return mock data
	return map[string]int64{
		"timeout":           5,
		"validation_error":  3,
		"network_error":     2,
		"authentication":    1,
	}, nil
}

func (mc *MetricsCollector) getStepMetrics(ctx context.Context, workflowID uuid.UUID, startTime, endTime time.Time) ([]StepMetrics, error) {
	// This would query step execution data
	// For now, return mock data
	return []StepMetrics{
		{
			StepName:       "validate_input",
			StepType:       "validation",
			Executions:     100,
			Successful:     95,
			Failed:         5,
			AverageRuntime: 2 * time.Second,
			SuccessRate:    95.0,
			CommonErrors:   []string{"invalid_format", "missing_field"},
		},
	}, nil
}

func (mc *MetricsCollector) getHourlyExecutionCounts(ctx context.Context, startTime, endTime time.Time, filters ExecutionMetricsFilters) ([]HourlyExecutionCount, error) {
	// This would query execution data and aggregate by hour
	// For now, return mock data
	return []HourlyExecutionCount{
		{
			Hour:       startTime.Truncate(time.Hour),
			Executions: 25,
			Successful: 22,
			Failed:     3,
		},
	}, nil
}

func (mc *MetricsCollector) getTopFailedWorkflows(ctx context.Context, startTime, endTime time.Time, limit int) ([]WorkflowFailureCount, error) {
	// This would query execution data and aggregate failures by workflow
	// For now, return mock data
	return []WorkflowFailureCount{
		{
			WorkflowID:   uuid.New(),
			WorkflowName: "Data Processing Pipeline",
			FailureCount: 15,
			ErrorRate:    12.5,
		},
	}, nil
}

func (mc *MetricsCollector) getWorkflowTrends(ctx context.Context, startTime, endTime time.Time) ([]WorkflowTrendPoint, error) {
	// This would query workflow data and aggregate by time intervals
	// For now, return mock data
	return []WorkflowTrendPoint{
		{
			Timestamp:        startTime,
			TotalWorkflows:   50,
			ActiveWorkflows:  35,
			NewWorkflows:     2,
			UpdatedWorkflows: 5,
		},
	}, nil
}

// HealthCheck checks the health of the metrics collector
func (mc *MetricsCollector) HealthCheck(ctx context.Context) error {
	// Check if we can access the repository manager
	return mc.repoManager.HealthCheck(ctx)
}