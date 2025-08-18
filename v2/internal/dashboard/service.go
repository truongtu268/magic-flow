package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// Service provides dashboard functionality
type Service struct {
	repoManager database.RepositoryManager
	metricsCollector *MetricsCollector
	realtimeManager *RealtimeManager
}

// NewService creates a new dashboard service
func NewService(repoManager database.RepositoryManager) *Service {
	return &Service{
		repoManager: repoManager,
		metricsCollector: NewMetricsCollector(repoManager),
		realtimeManager: NewRealtimeManager(),
	}
}

// DashboardOverview represents the main dashboard overview data
type DashboardOverview struct {
	WorkflowStats    WorkflowStats    `json:"workflow_stats"`
	ExecutionStats   ExecutionStats   `json:"execution_stats"`
	PerformanceStats PerformanceStats `json:"performance_stats"`
	SystemStats      SystemStats      `json:"system_stats"`
	RecentActivity   []ActivityItem   `json:"recent_activity"`
	Alerts           []AlertItem      `json:"alerts"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// WorkflowStats represents workflow statistics
type WorkflowStats struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Inactive int64 `json:"inactive"`
	Draft    int64 `json:"draft"`
}

// ExecutionStats represents execution statistics
type ExecutionStats struct {
	Total       int64   `json:"total"`
	Running     int64   `json:"running"`
	Completed   int64   `json:"completed"`
	Failed      int64   `json:"failed"`
	SuccessRate float64 `json:"success_rate"`
	Today       int64   `json:"today"`
	ThisWeek    int64   `json:"this_week"`
	ThisMonth   int64   `json:"this_month"`
}

// PerformanceStats represents performance statistics
type PerformanceStats struct {
	AverageExecutionTime time.Duration `json:"average_execution_time"`
	MedianExecutionTime  time.Duration `json:"median_execution_time"`
	Throughput           float64       `json:"throughput"` // executions per hour
	ErrorRate            float64       `json:"error_rate"`
	P95ExecutionTime     time.Duration `json:"p95_execution_time"`
	P99ExecutionTime     time.Duration `json:"p99_execution_time"`
}

// SystemStats represents system statistics
type SystemStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	ActiveUsers int64   `json:"active_users"`
	APIRequests int64   `json:"api_requests"`
	Uptime      string  `json:"uptime"`
}

// ActivityItem represents a recent activity item
type ActivityItem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	UserName    string    `json:"user_name"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AlertItem represents an alert item
type AlertItem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Resolved    bool      `json:"resolved"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// GetDashboardOverview returns the main dashboard overview
func (s *Service) GetDashboardOverview(ctx context.Context) (*DashboardOverview, error) {
	// Get workflow stats
	workflowStats, err := s.getWorkflowStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow stats: %w", err)
	}

	// Get execution stats
	executionStats, err := s.getExecutionStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution stats: %w", err)
	}

	// Get performance stats
	performanceStats, err := s.getPerformanceStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance stats: %w", err)
	}

	// Get system stats
	systemStats, err := s.getSystemStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system stats: %w", err)
	}

	// Get recent activity
	recentActivity, err := s.getRecentActivity(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}

	// Get alerts
	alerts, err := s.getActiveAlerts(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	return &DashboardOverview{
		WorkflowStats:    *workflowStats,
		ExecutionStats:   *executionStats,
		PerformanceStats: *performanceStats,
		SystemStats:      *systemStats,
		RecentActivity:   recentActivity,
		Alerts:           alerts,
		UpdatedAt:        time.Now(),
	}, nil
}

// getWorkflowStats retrieves workflow statistics
func (s *Service) getWorkflowStats(ctx context.Context) (*WorkflowStats, error) {
	workflowRepo := s.repoManager.WorkflowRepository()

	total, err := workflowRepo.Count(ctx, nil)
	if err != nil {
		return nil, err
	}

	active, err := workflowRepo.CountByStatus(ctx, models.WorkflowStatusActive)
	if err != nil {
		return nil, err
	}

	inactive, err := workflowRepo.CountByStatus(ctx, models.WorkflowStatusInactive)
	if err != nil {
		return nil, err
	}

	draft, err := workflowRepo.CountByStatus(ctx, models.WorkflowStatusDraft)
	if err != nil {
		return nil, err
	}

	return &WorkflowStats{
		Total:    total,
		Active:   active,
		Inactive: inactive,
		Draft:    draft,
	}, nil
}

// getExecutionStats retrieves execution statistics
func (s *Service) getExecutionStats(ctx context.Context) (*ExecutionStats, error) {
	executionRepo := s.repoManager.ExecutionRepository()

	total, err := executionRepo.Count(ctx, nil)
	if err != nil {
		return nil, err
	}

	running, err := executionRepo.CountByStatus(ctx, models.ExecutionStatusRunning)
	if err != nil {
		return nil, err
	}

	completed, err := executionRepo.CountByStatus(ctx, models.ExecutionStatusCompleted)
	if err != nil {
		return nil, err
	}

	failed, err := executionRepo.CountByStatus(ctx, models.ExecutionStatusFailed)
	if err != nil {
		return nil, err
	}

	// Calculate success rate
	successRate := float64(0)
	if total > 0 {
		successRate = float64(completed) / float64(total) * 100
	}

	// Get time-based stats
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	today, err := executionRepo.CountByTimeRange(ctx, startOfDay, now)
	if err != nil {
		return nil, err
	}

	thisWeek, err := executionRepo.CountByTimeRange(ctx, startOfWeek, now)
	if err != nil {
		return nil, err
	}

	thisMonth, err := executionRepo.CountByTimeRange(ctx, startOfMonth, now)
	if err != nil {
		return nil, err
	}

	return &ExecutionStats{
		Total:       total,
		Running:     running,
		Completed:   completed,
		Failed:      failed,
		SuccessRate: successRate,
		Today:       today,
		ThisWeek:    thisWeek,
		ThisMonth:   thisMonth,
	}, nil
}

// getPerformanceStats retrieves performance statistics
func (s *Service) getPerformanceStats(ctx context.Context) (*PerformanceStats, error) {
	executionRepo := s.repoManager.ExecutionRepository()

	// Get execution statistics for the last 24 hours
	now := time.Now()
	last24h := now.Add(-24 * time.Hour)

	stats, err := executionRepo.GetExecutionStatistics(ctx, &last24h, &now)
	if err != nil {
		return nil, err
	}

	// Calculate throughput (executions per hour)
	throughput := float64(stats.TotalExecutions) / 24.0

	// Calculate error rate
	errorRate := float64(0)
	if stats.TotalExecutions > 0 {
		errorRate = float64(stats.FailedExecutions) / float64(stats.TotalExecutions) * 100
	}

	return &PerformanceStats{
		AverageExecutionTime: stats.AverageDuration,
		MedianExecutionTime:  stats.MedianDuration,
		Throughput:           throughput,
		ErrorRate:            errorRate,
		P95ExecutionTime:     stats.P95Duration,
		P99ExecutionTime:     stats.P99Duration,
	}, nil
}

// getSystemStats retrieves system statistics
func (s *Service) getSystemStats(ctx context.Context) (*SystemStats, error) {
	// This would typically integrate with system monitoring tools
	// For now, we'll return mock data
	return &SystemStats{
		CPUUsage:    45.2,
		MemoryUsage: 67.8,
		DiskUsage:   23.4,
		ActiveUsers: 15,
		APIRequests: 1234,
		Uptime:      "5d 12h 34m",
	}, nil
}

// getRecentActivity retrieves recent activity items
func (s *Service) getRecentActivity(ctx context.Context, limit int) ([]ActivityItem, error) {
	// This would typically query an activity log table
	// For now, we'll return mock data
	return []ActivityItem{
		{
			ID:          uuid.New(),
			Type:        "workflow_created",
			Title:       "New Workflow Created",
			Description: "User created workflow 'Data Processing Pipeline'",
			UserID:      uuid.New(),
			UserName:    "John Doe",
			Timestamp:   time.Now().Add(-5 * time.Minute),
		},
		{
			ID:          uuid.New(),
			Type:        "execution_completed",
			Title:       "Execution Completed",
			Description: "Workflow 'Email Campaign' completed successfully",
			UserID:      uuid.New(),
			UserName:    "Jane Smith",
			Timestamp:   time.Now().Add(-15 * time.Minute),
		},
	}, nil
}

// getActiveAlerts retrieves active alerts
func (s *Service) getActiveAlerts(ctx context.Context, limit int) ([]AlertItem, error) {
	alertRepo := s.repoManager.AlertRepository()

	alerts, err := alertRepo.GetActiveAlerts(ctx, limit)
	if err != nil {
		return nil, err
	}

	alertItems := make([]AlertItem, len(alerts))
	for i, alert := range alerts {
		alertItems[i] = AlertItem{
			ID:          alert.ID,
			Type:        string(alert.Type),
			Severity:    string(alert.Severity),
			Title:       alert.Name,
			Description: alert.Description,
			Timestamp:   alert.CreatedAt,
			Resolved:    !alert.Enabled,
		}
	}

	return alertItems, nil
}

// GetWorkflowMetrics returns metrics for a specific workflow
func (s *Service) GetWorkflowMetrics(ctx context.Context, workflowID uuid.UUID, timeRange string) (*WorkflowMetrics, error) {
	return s.metricsCollector.GetWorkflowMetrics(ctx, workflowID, timeRange)
}

// GetExecutionMetrics returns metrics for executions
func (s *Service) GetExecutionMetrics(ctx context.Context, filters ExecutionMetricsFilters) (*ExecutionMetrics, error) {
	return s.metricsCollector.GetExecutionMetrics(ctx, filters)
}

// GetSystemMetrics returns system-wide metrics
func (s *Service) GetSystemMetrics(ctx context.Context, timeRange string) (*SystemMetrics, error) {
	return s.metricsCollector.GetSystemMetrics(ctx, timeRange)
}

// SubscribeToRealtimeUpdates subscribes to real-time dashboard updates
func (s *Service) SubscribeToRealtimeUpdates(ctx context.Context, clientID string) (<-chan RealtimeUpdate, error) {
	return s.realtimeManager.Subscribe(ctx, clientID)
}

// UnsubscribeFromRealtimeUpdates unsubscribes from real-time updates
func (s *Service) UnsubscribeFromRealtimeUpdates(clientID string) {
	s.realtimeManager.Unsubscribe(clientID)
}

// PublishRealtimeUpdate publishes a real-time update to all subscribers
func (s *Service) PublishRealtimeUpdate(update RealtimeUpdate) {
	s.realtimeManager.Publish(update)
}

// GetDashboardConfig returns dashboard configuration
func (s *Service) GetDashboardConfig(ctx context.Context, userID uuid.UUID) (*DashboardConfig, error) {
	dashboardRepo := s.repoManager.DashboardRepository()

	// Get user's dashboard configuration
	dashboards, err := dashboardRepo.GetByCreator(ctx, userID, 1, 0)
	if err != nil {
		return nil, err
	}

	config := &DashboardConfig{
		UserID:    userID,
		Theme:     "light",
		Language:  "en",
		Timezone:  "UTC",
		RefreshInterval: 30,
		Widgets:   []WidgetConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(dashboards) > 0 {
		// Load configuration from existing dashboard
		dashboard := dashboards[0]
		config.ID = dashboard.ID
		config.Name = dashboard.Name
		config.CreatedAt = dashboard.CreatedAt
		config.UpdatedAt = dashboard.UpdatedAt
	}

	return config, nil
}

// UpdateDashboardConfig updates dashboard configuration
func (s *Service) UpdateDashboardConfig(ctx context.Context, config *DashboardConfig) error {
	dashboardRepo := s.repoManager.DashboardRepository()

	if config.ID == uuid.Nil {
		// Create new dashboard
		dashboard := &models.Dashboard{
			Name:        config.Name,
			Description: "User dashboard configuration",
			CreatorID:   config.UserID,
			IsPublic:    false,
			Config:      map[string]interface{}{"widgets": config.Widgets},
		}

		err := dashboardRepo.Create(ctx, dashboard)
		if err != nil {
			return err
		}

		config.ID = dashboard.ID
	} else {
		// Update existing dashboard
		dashboard, err := dashboardRepo.GetByID(ctx, config.ID)
		if err != nil {
			return err
		}

		dashboard.Name = config.Name
		dashboard.Config = map[string]interface{}{"widgets": config.Widgets}

		err = dashboardRepo.Update(ctx, dashboard)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetHealthStatus returns the health status of the dashboard service
func (s *Service) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]ServiceHealth),
	}

	// Check database health
	if err := s.repoManager.HealthCheck(ctx); err != nil {
		status.Services["database"] = ServiceHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
		status.Status = "degraded"
	} else {
		status.Services["database"] = ServiceHealth{
			Status: "healthy",
		}
	}

	// Check metrics collector health
	if err := s.metricsCollector.HealthCheck(ctx); err != nil {
		status.Services["metrics"] = ServiceHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
		status.Status = "degraded"
	} else {
		status.Services["metrics"] = ServiceHealth{
			Status: "healthy",
		}
	}

	// Check realtime manager health
	status.Services["realtime"] = ServiceHealth{
		Status: "healthy",
	}

	return status, nil
}

// Close closes the dashboard service and cleans up resources
func (s *Service) Close() error {
	if s.realtimeManager != nil {
		s.realtimeManager.Close()
	}
	return nil
}