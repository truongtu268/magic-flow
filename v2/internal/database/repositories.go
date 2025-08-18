package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// WorkflowRepository handles workflow data operations
type WorkflowRepository struct {
	db *gorm.DB
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(workflow *models.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowRepository) GetByID(id uuid.UUID) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.Preload("Versions").First(&workflow, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *WorkflowRepository) GetByName(name string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.Preload("Versions").First(&workflow, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *WorkflowRepository) List(limit, offset int, status string) ([]*models.Workflow, int64, error) {
	var workflows []*models.Workflow
	var total int64

	query := r.db.Model(&models.Workflow{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get workflows with pagination
	err := query.Preload("Versions").Limit(limit).Offset(offset).Order("created_at DESC").Find(&workflows).Error
	return workflows, total, err
}

func (r *WorkflowRepository) Update(workflow *models.Workflow) error {
	return r.db.Save(workflow).Error
}

func (r *WorkflowRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Workflow{}, "id = ?", id).Error
}

func (r *WorkflowRepository) UpdateStatus(id uuid.UUID, status models.WorkflowStatus) error {
	return r.db.Model(&models.Workflow{}).Where("id = ?", id).Update("status", status).Error
}

// ExecutionRepository handles execution data operations
type ExecutionRepository struct {
	db *gorm.DB
}

// NewExecutionRepository creates a new execution repository
func NewExecutionRepository(db *gorm.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

func (r *ExecutionRepository) Create(execution *models.Execution) error {
	return r.db.Create(execution).Error
}

func (r *ExecutionRepository) GetByID(id uuid.UUID) (*models.Execution, error) {
	var execution models.Execution
	err := r.db.Preload("Steps").Preload("Events").First(&execution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *ExecutionRepository) List(workflowID *uuid.UUID, limit, offset int, status string) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{})
	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Preload("Steps").Limit(limit).Offset(offset).Order("started_at DESC").Find(&executions).Error
	return executions, total, err
}

func (r *ExecutionRepository) Update(execution *models.Execution) error {
	return r.db.Save(execution).Error
}

func (r *ExecutionRepository) UpdateStatus(id uuid.UUID, status models.ExecutionStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	if status == models.ExecutionStatusCompleted || status == models.ExecutionStatusFailed || status == models.ExecutionStatusCancelled {
		updates["completed_at"] = time.Now().UTC()
	}

	return r.db.Model(&models.Execution{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ExecutionRepository) GetActiveExecutions() ([]*models.Execution, error) {
	var executions []*models.Execution
	err := r.db.Where("status IN ?", []models.ExecutionStatus{
		models.ExecutionStatusPending,
		models.ExecutionStatusRunning,
	}).Find(&executions).Error
	return executions, err
}

func (r *ExecutionRepository) GetExecutionStats(workflowID *uuid.UUID, from, to time.Time) (map[string]int64, error) {
	stats := make(map[string]int64)

	query := r.db.Model(&models.Execution{}).Where("started_at BETWEEN ? AND ?", from, to)
	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	// Total executions
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Status breakdown
	statuses := []models.ExecutionStatus{
		models.ExecutionStatusCompleted,
		models.ExecutionStatusFailed,
		models.ExecutionStatusCancelled,
		models.ExecutionStatusRunning,
		models.ExecutionStatusPending,
	}

	for _, status := range statuses {
		var count int64
		if err := query.Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[string(status)] = count
	}

	return stats, nil
}

// StepExecutionRepository handles step execution data operations
type StepExecutionRepository struct {
	db *gorm.DB
}

// NewStepExecutionRepository creates a new step execution repository
func NewStepExecutionRepository(db *gorm.DB) *StepExecutionRepository {
	return &StepExecutionRepository{db: db}
}

func (r *StepExecutionRepository) Create(stepExecution *models.StepExecution) error {
	return r.db.Create(stepExecution).Error
}

func (r *StepExecutionRepository) GetByExecutionID(executionID uuid.UUID) ([]*models.StepExecution, error) {
	var steps []*models.StepExecution
	err := r.db.Where("execution_id = ?", executionID).Order("started_at ASC").Find(&steps).Error
	return steps, err
}

func (r *StepExecutionRepository) UpdateStatus(id uuid.UUID, status models.StepStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	if status == models.StepStatusCompleted || status == models.StepStatusFailed || status == models.StepStatusSkipped {
		updates["completed_at"] = time.Now().UTC()
	}

	return r.db.Model(&models.StepExecution{}).Where("id = ?", id).Updates(updates).Error
}

// ExecutionEventRepository handles execution event data operations
type ExecutionEventRepository struct {
	db *gorm.DB
}

// NewExecutionEventRepository creates a new execution event repository
func NewExecutionEventRepository(db *gorm.DB) *ExecutionEventRepository {
	return &ExecutionEventRepository{db: db}
}

func (r *ExecutionEventRepository) Create(event *models.ExecutionEvent) error {
	return r.db.Create(event).Error
}

func (r *ExecutionEventRepository) GetByExecutionID(executionID uuid.UUID, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("execution_id = ?", executionID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&events).Error
	return events, total, err
}

// WorkflowVersionRepository handles workflow version data operations
type WorkflowVersionRepository struct {
	db *gorm.DB
}

// NewWorkflowVersionRepository creates a new workflow version repository
func NewWorkflowVersionRepository(db *gorm.DB) *WorkflowVersionRepository {
	return &WorkflowVersionRepository{db: db}
}

func (r *WorkflowVersionRepository) Create(version *models.WorkflowVersion) error {
	return r.db.Create(version).Error
}

func (r *WorkflowVersionRepository) GetByID(id uuid.UUID) (*models.WorkflowVersion, error) {
	var version models.WorkflowVersion
	err := r.db.Preload("Deployments").First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *WorkflowVersionRepository) GetByWorkflowID(workflowID uuid.UUID, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{}).Where("workflow_id = ?", workflowID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Preload("Deployments").Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	return versions, total, err
}

func (r *WorkflowVersionRepository) GetLatestVersion(workflowID uuid.UUID) (*models.WorkflowVersion, error) {
	var version models.WorkflowVersion
	err := r.db.Where("workflow_id = ?", workflowID).Order("created_at DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *WorkflowVersionRepository) UpdateStatus(id uuid.UUID, status models.WorkflowVersionStatus) error {
	return r.db.Model(&models.WorkflowVersion{}).Where("id = ?", id).Update("status", status).Error
}

// MetricsRepository handles metrics data operations
type MetricsRepository struct {
	db *gorm.DB
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	return &MetricsRepository{db: db}
}

func (r *MetricsRepository) CreateWorkflowMetric(metric *models.WorkflowMetric) error {
	return r.db.Create(metric).Error
}

func (r *MetricsRepository) CreateSystemMetric(metric *models.SystemMetric) error {
	return r.db.Create(metric).Error
}

func (r *MetricsRepository) CreateBusinessMetric(metric *models.BusinessMetric) error {
	return r.db.Create(metric).Error
}

func (r *MetricsRepository) GetWorkflowMetrics(workflowID *uuid.UUID, metricName string, from, to time.Time, limit, offset int) ([]*models.WorkflowMetric, int64, error) {
	var metrics []*models.WorkflowMetric
	var total int64

	query := r.db.Model(&models.WorkflowMetric{}).Where("timestamp BETWEEN ? AND ?", from, to)
	if workflowID != nil {
		query = query.Where("labels->>'workflow_id' = ?", workflowID.String())
	}
	if metricName != "" {
		query = query.Where("name = ?", metricName)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get metrics with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&metrics).Error
	return metrics, total, err
}

func (r *MetricsRepository) GetSystemMetrics(metricName string, from, to time.Time, limit, offset int) ([]*models.SystemMetric, int64, error) {
	var metrics []*models.SystemMetric
	var total int64

	query := r.db.Model(&models.SystemMetric{}).Where("timestamp BETWEEN ? AND ?", from, to)
	if metricName != "" {
		query = query.Where("name = ?", metricName)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get metrics with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&metrics).Error
	return metrics, total, err
}

// AlertRepository handles alert data operations
type AlertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

func (r *AlertRepository) GetByID(id uuid.UUID) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Preload("Events").First(&alert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) List(limit, offset int, enabled *bool) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{})
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	return alerts, total, err
}

func (r *AlertRepository) Update(alert *models.Alert) error {
	return r.db.Save(alert).Error
}

func (r *AlertRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Alert{}, "id = ?", id).Error
}

func (r *AlertRepository) CreateEvent(event *models.AlertEvent) error {
	return r.db.Create(event).Error
}

func (r *AlertRepository) GetEvents(alertID uuid.UUID, limit, offset int) ([]*models.AlertEvent, int64, error) {
	var events []*models.AlertEvent
	var total int64

	query := r.db.Model(&models.AlertEvent{}).Where("alert_id = ?", alertID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&events).Error
	return events, total, err
}

// DashboardRepository handles dashboard data operations
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) Create(dashboard *models.Dashboard) error {
	return r.db.Create(dashboard).Error
}

func (r *DashboardRepository) GetByID(id uuid.UUID) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.First(&dashboard, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (r *DashboardRepository) List(limit, offset int, createdBy string, isPublic *bool) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{})
	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	return dashboards, total, err
}

func (r *DashboardRepository) Update(dashboard *models.Dashboard) error {
	return r.db.Save(dashboard).Error
}

func (r *DashboardRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Dashboard{}, "id = ?", id).Error
}

// RepositoryManager manages all repositories
type RepositoryManager struct {
	Workflow        *WorkflowRepository
	Execution       *ExecutionRepository
	StepExecution   *StepExecutionRepository
	ExecutionEvent  *ExecutionEventRepository
	WorkflowVersion *WorkflowVersionRepository
	Metrics         *MetricsRepository
	Alert           *AlertRepository
	Dashboard       *DashboardRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		Workflow:        NewWorkflowRepository(db),
		Execution:       NewExecutionRepository(db),
		StepExecution:   NewStepExecutionRepository(db),
		ExecutionEvent:  NewExecutionEventRepository(db),
		WorkflowVersion: NewWorkflowVersionRepository(db),
		Metrics:         NewMetricsRepository(db),
		Alert:           NewAlertRepository(db),
		Dashboard:       NewDashboardRepository(db),
	}
}