package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// MetricsService handles metrics business logic
type MetricsService struct {
	repos  *database.RepositoryManager
	logger *logrus.Logger
}

// NewMetricsService creates a new metrics service
func NewMetricsService(repos *database.RepositoryManager, logger *logrus.Logger) *MetricsService {
	return &MetricsService{
		repos:  repos,
		logger: logger,
	}
}

// GetWorkflowMetrics retrieves workflow metrics
func (s *MetricsService) GetWorkflowMetrics(req *GetWorkflowMetricsRequest) (*WorkflowMetricsResponse, error) {
	var metrics []*models.WorkflowMetrics
	var err error

	if req.WorkflowID != nil {
		metrics, err = s.repos.Metrics.GetWorkflowMetricsByWorkflowID(*req.WorkflowID, req.StartTime, req.EndTime)
	} else {
		metrics, err = s.repos.Metrics.GetWorkflowMetrics(req.StartTime, req.EndTime)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get workflow metrics: %w", err)
	}

	// Calculate aggregated metrics
	aggregated := s.aggregateWorkflowMetrics(metrics)

	return &WorkflowMetricsResponse{
		Metrics:    metrics,
		Aggregated: aggregated,
		TimeRange: TimeRange{
			Start: req.StartTime,
			End:   req.EndTime,
		},
	}, nil
}

// GetSystemMetrics retrieves system metrics
func (s *MetricsService) GetSystemMetrics(req *GetSystemMetricsRequest) (*SystemMetricsResponse, error) {
	metrics, err := s.repos.Metrics.GetSystemMetrics(req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get system metrics: %w", err)
	}

	// Calculate aggregated metrics
	aggregated := s.aggregateSystemMetrics(metrics)

	return &SystemMetricsResponse{
		Metrics:    metrics,
		Aggregated: aggregated,
		TimeRange: TimeRange{
			Start: req.StartTime,
			End:   req.EndTime,
		},
	}, nil
}

// GetBusinessMetrics retrieves business metrics
func (s *MetricsService) GetBusinessMetrics(req *GetBusinessMetricsRequest) (*BusinessMetricsResponse, error) {
	metrics, err := s.repos.Metrics.GetBusinessMetrics(req.StartTime, req.EndTime, req.MetricName)
	if err != nil {
		return nil, fmt.Errorf("failed to get business metrics: %w", err)
	}

	return &BusinessMetricsResponse{
		Metrics: metrics,
		TimeRange: TimeRange{
			Start: req.StartTime,
			End:   req.EndTime,
		},
	}, nil
}

// RecordBusinessMetric records a custom business metric
func (s *MetricsService) RecordBusinessMetric(req *RecordBusinessMetricRequest) error {
	metric := &models.BusinessMetrics{
		ID:         uuid.New(),
		Name:       req.Name,
		Value:      req.Value,
		Unit:       req.Unit,
		Tags:       req.Tags,
		WorkflowID: req.WorkflowID,
		Timestamp:  time.Now().UTC(),
	}

	if err := s.repos.Metrics.CreateBusinessMetric(metric); err != nil {
		return fmt.Errorf("failed to record business metric: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"metric_name": req.Name,
		"value":       req.Value,
		"unit":        req.Unit,
		"workflow_id": req.WorkflowID,
	}).Debug("Business metric recorded")

	return nil
}

// GetMetricAggregations retrieves metric aggregations
func (s *MetricsService) GetMetricAggregations(req *GetMetricAggregationsRequest) ([]*models.MetricAggregation, int64, error) {
	aggregations, total, err := s.repos.Metrics.GetMetricAggregations(
		req.Limit,
		req.Offset,
		req.MetricType,
		req.AggregationType,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get metric aggregations: %w", err)
	}

	return aggregations, total, nil
}

// CreateMetricAggregation creates a new metric aggregation
func (s *MetricsService) CreateMetricAggregation(req *CreateMetricAggregationRequest) (*models.MetricAggregation, error) {
	aggregation := &models.MetricAggregation{
		ID:              uuid.New(),
		MetricType:      req.MetricType,
		AggregationType: req.AggregationType,
		TimeWindow:      req.TimeWindow,
		Value:           req.Value,
		Tags:            req.Tags,
		Timestamp:       time.Now().UTC(),
	}

	if err := s.repos.Metrics.CreateMetricAggregation(aggregation); err != nil {
		return nil, fmt.Errorf("failed to create metric aggregation: %w", err)
	}

	return aggregation, nil
}

// GetDashboardOverview retrieves dashboard overview data
func (s *MetricsService) GetDashboardOverview() (*DashboardOverviewResponse, error) {
	// Get workflow counts
	totalWorkflows, err := s.repos.Workflow.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow count: %w", err)
	}

	activeWorkflows, err := s.repos.Workflow.CountByStatus(models.WorkflowStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get active workflow count: %w", err)
	}

	// Get execution counts for today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	totalExecutions, err := s.repos.Execution.CountByTimeRange(&today, &tomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution count: %w", err)
	}

	successfulExecutions, err := s.repos.Execution.CountByStatusAndTimeRange(models.ExecutionStatusCompleted, &today, &tomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to get successful execution count: %w", err)
	}

	failedExecutions, err := s.repos.Execution.CountByStatusAndTimeRange(models.ExecutionStatusFailed, &today, &tomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed execution count: %w", err)
	}

	runningExecutions, err := s.repos.Execution.CountByStatusAndTimeRange(models.ExecutionStatusRunning, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get running execution count: %w", err)
	}

	// Calculate success rate
	successRate := 0.0
	if totalExecutions > 0 {
		successRate = float64(successfulExecutions) / float64(totalExecutions) * 100
	}

	// Get recent executions
	recentExecutions, _, err := s.repos.Execution.List(10, 0, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get recent executions: %w", err)
	}

	return &DashboardOverviewResponse{
		Workflows: WorkflowOverview{
			Total:  totalWorkflows,
			Active: activeWorkflows,
		},
		Executions: ExecutionOverview{
			Today:       totalExecutions,
			Successful:  successfulExecutions,
			Failed:      failedExecutions,
			Running:     runningExecutions,
			SuccessRate: successRate,
		},
		RecentExecutions: recentExecutions,
		LastUpdated:      time.Now().UTC(),
	}, nil
}

// GetWorkflowStatusSummary retrieves workflow status summary
func (s *MetricsService) GetWorkflowStatusSummary() (*WorkflowStatusSummaryResponse, error) {
	statusMap := make(map[string]int64)

	// Get counts for each status
	statuses := []models.WorkflowStatus{
		models.WorkflowStatusDraft,
		models.WorkflowStatusActive,
		models.WorkflowStatusInactive,
		models.WorkflowStatusArchived,
	}

	for _, status := range statuses {
		count, err := s.repos.Workflow.CountByStatus(status)
		if err != nil {
			return nil, fmt.Errorf("failed to get count for status %s: %w", status, err)
		}
		statusMap[string(status)] = count
	}

	return &WorkflowStatusSummaryResponse{
		StatusCounts: statusMap,
		LastUpdated:  time.Now().UTC(),
	}, nil
}

// GetSystemHealth retrieves system health information
func (s *MetricsService) GetSystemHealth() (*SystemHealthResponse, error) {
	// Check database health
	dbHealthy := true
	if err := s.repos.Database.Health(); err != nil {
		dbHealthy = false
		s.logger.WithError(err).Warn("Database health check failed")
	}

	// Get latest system metrics
	latestMetrics, err := s.repos.Metrics.GetLatestSystemMetrics()
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get latest system metrics")
	}

	// Determine overall health
	overallHealthy := dbHealthy
	if latestMetrics != nil {
		// Check if CPU usage is too high (>90%)
		if latestMetrics.CPUUsage > 90 {
			overallHealthy = false
		}
		// Check if memory usage is too high (>90%)
		if latestMetrics.MemoryUsage > 90 {
			overallHealthy = false
		}
		// Check if disk usage is too high (>95%)
		if latestMetrics.DiskUsage > 95 {
			overallHealthy = false
		}
	}

	return &SystemHealthResponse{
		Healthy:       overallHealthy,
		Database:      dbHealthy,
		SystemMetrics: latestMetrics,
		LastChecked:   time.Now().UTC(),
	}, nil
}

// GetLiveMetrics retrieves live metrics
func (s *MetricsService) GetLiveMetrics() (*LiveMetricsResponse, error) {
	// Get current running executions
	runningExecutions, err := s.repos.Execution.GetActiveExecutions()
	if err != nil {
		return nil, fmt.Errorf("failed to get running executions: %w", err)
	}

	// Get latest system metrics
	latestSystemMetrics, err := s.repos.Metrics.GetLatestSystemMetrics()
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get latest system metrics")
	}

	// Get execution rate (executions per minute in last hour)
	lastHour := time.Now().UTC().Add(-time.Hour)
	now := time.Now().UTC()
	executionsLastHour, err := s.repos.Execution.CountByTimeRange(&lastHour, &now)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions in last hour: %w", err)
	}

	executionRate := float64(executionsLastHour) / 60.0 // per minute

	return &LiveMetricsResponse{
		RunningExecutions: len(runningExecutions),
		ExecutionRate:     executionRate,
		SystemMetrics:     latestSystemMetrics,
		Timestamp:         time.Now().UTC(),
	}, nil
}

// Helper methods
func (s *MetricsService) aggregateWorkflowMetrics(metrics []*models.WorkflowMetrics) *AggregatedWorkflowMetrics {
	if len(metrics) == 0 {
		return &AggregatedWorkflowMetrics{}
	}

	var totalExecutions, successfulExecutions, failedExecutions int64
	var totalDuration, avgDuration time.Duration

	for _, metric := range metrics {
		totalExecutions += metric.TotalExecutions
		successfulExecutions += metric.SuccessfulExecutions
		failedExecutions += metric.FailedExecutions
		totalDuration += metric.AverageExecutionTime
	}

	if len(metrics) > 0 {
		avgDuration = totalDuration / time.Duration(len(metrics))
	}

	successRate := 0.0
	if totalExecutions > 0 {
		successRate = float64(successfulExecutions) / float64(totalExecutions) * 100
	}

	return &AggregatedWorkflowMetrics{
		TotalExecutions:      totalExecutions,
		SuccessfulExecutions: successfulExecutions,
		FailedExecutions:     failedExecutions,
		SuccessRate:          successRate,
		AverageExecutionTime: avgDuration,
	}
}

func (s *MetricsService) aggregateSystemMetrics(metrics []*models.SystemMetrics) *AggregatedSystemMetrics {
	if len(metrics) == 0 {
		return &AggregatedSystemMetrics{}
	}

	var totalCPU, totalMemory, totalDisk float64
	var maxCPU, maxMemory, maxDisk float64

	for i, metric := range metrics {
		totalCPU += metric.CPUUsage
		totalMemory += metric.MemoryUsage
		totalDisk += metric.DiskUsage

		if i == 0 || metric.CPUUsage > maxCPU {
			maxCPU = metric.CPUUsage
		}
		if i == 0 || metric.MemoryUsage > maxMemory {
			maxMemory = metric.MemoryUsage
		}
		if i == 0 || metric.DiskUsage > maxDisk {
			maxDisk = metric.DiskUsage
		}
	}

	count := float64(len(metrics))
	return &AggregatedSystemMetrics{
		AverageCPUUsage:    totalCPU / count,
		AverageMemoryUsage: totalMemory / count,
		AverageDiskUsage:   totalDisk / count,
		MaxCPUUsage:        maxCPU,
		MaxMemoryUsage:     maxMemory,
		MaxDiskUsage:       maxDisk,
	}
}

// Request/Response types
type GetWorkflowMetricsRequest struct {
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
}

type WorkflowMetricsResponse struct {
	Metrics    []*models.WorkflowMetrics   `json:"metrics"`
	Aggregated *AggregatedWorkflowMetrics `json:"aggregated"`
	TimeRange  TimeRange                  `json:"time_range"`
}

type AggregatedWorkflowMetrics struct {
	TotalExecutions      int64         `json:"total_executions"`
	SuccessfulExecutions int64         `json:"successful_executions"`
	FailedExecutions     int64         `json:"failed_executions"`
	SuccessRate          float64       `json:"success_rate"`
	AverageExecutionTime time.Duration `json:"average_execution_time"`
}

type GetSystemMetricsRequest struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

type SystemMetricsResponse struct {
	Metrics    []*models.SystemMetrics   `json:"metrics"`
	Aggregated *AggregatedSystemMetrics `json:"aggregated"`
	TimeRange  TimeRange                `json:"time_range"`
}

type AggregatedSystemMetrics struct {
	AverageCPUUsage    float64 `json:"average_cpu_usage"`
	AverageMemoryUsage float64 `json:"average_memory_usage"`
	AverageDiskUsage   float64 `json:"average_disk_usage"`
	MaxCPUUsage        float64 `json:"max_cpu_usage"`
	MaxMemoryUsage     float64 `json:"max_memory_usage"`
	MaxDiskUsage       float64 `json:"max_disk_usage"`
}

type GetBusinessMetricsRequest struct {
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	MetricName string     `json:"metric_name,omitempty"`
}

type BusinessMetricsResponse struct {
	Metrics   []*models.BusinessMetrics `json:"metrics"`
	TimeRange TimeRange                 `json:"time_range"`
}

type RecordBusinessMetricRequest struct {
	Name       string                 `json:"name" validate:"required"`
	Value      float64                `json:"value" validate:"required"`
	Unit       string                 `json:"unit,omitempty"`
	Tags       map[string]interface{} `json:"tags,omitempty"`
	WorkflowID *uuid.UUID             `json:"workflow_id,omitempty"`
}

type GetMetricAggregationsRequest struct {
	Limit           int        `json:"limit"`
	Offset          int        `json:"offset"`
	MetricType      string     `json:"metric_type,omitempty"`
	AggregationType string     `json:"aggregation_type,omitempty"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	EndTime         *time.Time `json:"end_time,omitempty"`
}

type CreateMetricAggregationRequest struct {
	MetricType      string                 `json:"metric_type" validate:"required"`
	AggregationType string                 `json:"aggregation_type" validate:"required"`
	TimeWindow      string                 `json:"time_window" validate:"required"`
	Value           float64                `json:"value" validate:"required"`
	Tags            map[string]interface{} `json:"tags,omitempty"`
}

type TimeRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type DashboardOverviewResponse struct {
	Workflows        WorkflowOverview     `json:"workflows"`
	Executions       ExecutionOverview    `json:"executions"`
	RecentExecutions []*models.Execution  `json:"recent_executions"`
	LastUpdated      time.Time            `json:"last_updated"`
}

type WorkflowOverview struct {
	Total  int64 `json:"total"`
	Active int64 `json:"active"`
}

type ExecutionOverview struct {
	Today       int64   `json:"today"`
	Successful  int64   `json:"successful"`
	Failed      int64   `json:"failed"`
	Running     int64   `json:"running"`
	SuccessRate float64 `json:"success_rate"`
}

type WorkflowStatusSummaryResponse struct {
	StatusCounts map[string]int64 `json:"status_counts"`
	LastUpdated  time.Time        `json:"last_updated"`
}

type SystemHealthResponse struct {
	Healthy       bool                   `json:"healthy"`
	Database      bool                   `json:"database"`
	SystemMetrics *models.SystemMetrics  `json:"system_metrics"`
	LastChecked   time.Time              `json:"last_checked"`
}

type LiveMetricsResponse struct {
	RunningExecutions int                   `json:"running_executions"`
	ExecutionRate     float64               `json:"execution_rate"`
	SystemMetrics     *models.SystemMetrics `json:"system_metrics"`
	Timestamp         time.Time             `json:"timestamp"`
}