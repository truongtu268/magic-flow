package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// executionRepository implements ExecutionRepository interface
type executionRepository struct {
	db *gorm.DB
}

// NewExecutionRepository creates a new execution repository
func NewExecutionRepository(db *gorm.DB) ExecutionRepository {
	return &executionRepository{db: db}
}

// Create creates a new execution
func (r *executionRepository) Create(execution *models.Execution) error {
	return r.db.Create(execution).Error
}

// GetByID retrieves an execution by ID
func (r *executionRepository) GetByID(id uuid.UUID) (*models.Execution, error) {
	var execution models.Execution
	err := r.db.Where("id = ?", id).First(&execution).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

// List retrieves executions with pagination and filtering
func (r *executionRepository) List(limit, offset int, workflowID *uuid.UUID, status string) ([]*models.Execution, int64, error) {
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
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// Update updates an execution
func (r *executionRepository) Update(execution *models.Execution) error {
	return r.db.Save(execution).Error
}

// Delete deletes an execution
func (r *executionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Execution{}, "id = ?", id).Error
}

// GetByWorkflowID retrieves executions for a specific workflow
func (r *executionRepository) GetByWorkflowID(workflowID uuid.UUID, limit, offset int) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{}).Where("workflow_id = ?", workflowID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// GetByStatus retrieves executions by status
func (r *executionRepository) GetByStatus(status models.ExecutionStatus, limit, offset int) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// GetActiveExecutions retrieves all active (running/pending) executions
func (r *executionRepository) GetActiveExecutions() ([]*models.Execution, error) {
	var executions []*models.Execution
	err := r.db.Where("status IN ?", []models.ExecutionStatus{
		models.ExecutionStatusPending,
		models.ExecutionStatusRunning,
	}).Find(&executions).Error
	return executions, err
}

// UpdateStatus updates only the status of an execution
func (r *executionRepository) UpdateStatus(id uuid.UUID, status models.ExecutionStatus) error {
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	// Set completion time for terminal states
	if status == models.ExecutionStatusCompleted || status == models.ExecutionStatusFailed || status == models.ExecutionStatusCancelled {
		updateData["completed_at"] = time.Now().UTC()
	}

	// Set start time for running state
	if status == models.ExecutionStatusRunning {
		updateData["started_at"] = time.Now().UTC()
	}

	return r.db.Model(&models.Execution{}).Where("id = ?", id).Updates(updateData).Error
}

// GetExecutionsByTimeRange retrieves executions within a time range
func (r *executionRepository) GetExecutionsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{})

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// CountByTimeRange counts executions within a time range
func (r *executionRepository) CountByTimeRange(startTime, endTime *time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.Execution{})

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	err := query.Count(&count).Error
	return count, err
}

// CountByStatusAndTimeRange counts executions by status within a time range
func (r *executionRepository) CountByStatusAndTimeRange(status models.ExecutionStatus, startTime, endTime *time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.Execution{}).Where("status = ?", status)

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetLongRunningExecutions retrieves executions that have been running for longer than the specified duration
func (r *executionRepository) GetLongRunningExecutions(duration time.Duration) ([]*models.Execution, error) {
	var executions []*models.Execution
	threshold := time.Now().UTC().Add(-duration)

	err := r.db.Where("status = ? AND started_at < ?", models.ExecutionStatusRunning, threshold).Find(&executions).Error
	return executions, err
}

// GetFailedExecutions retrieves failed executions within a time range
func (r *executionRepository) GetFailedExecutions(startTime, endTime *time.Time, limit, offset int) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{}).Where("status = ?", models.ExecutionStatusFailed)

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// GetExecutionsByCreator retrieves executions created by a specific user
func (r *executionRepository) GetExecutionsByCreator(createdBy string, limit, offset int) ([]*models.Execution, int64, error) {
	var executions []*models.Execution
	var total int64

	query := r.db.Model(&models.Execution{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&executions).Error
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// GetExecutionStats retrieves execution statistics
func (r *executionRepository) GetExecutionStats(workflowID *uuid.UUID, startTime, endTime *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.Execution{})

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Count by status
	statusStats := make(map[string]int64)
	statuses := []models.ExecutionStatus{
		models.ExecutionStatusPending,
		models.ExecutionStatusRunning,
		models.ExecutionStatusCompleted,
		models.ExecutionStatusFailed,
		models.ExecutionStatusCancelled,
	}

	for _, status := range statuses {
		var count int64
		statusQuery := query.Session(&gorm.Session{})
		if err := statusQuery.Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		statusStats[string(status)] = count
	}

	stats["by_status"] = statusStats

	// Total count
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Success rate
	successRate := 0.0
	if total > 0 {
		successRate = float64(statusStats[string(models.ExecutionStatusCompleted)]) / float64(total) * 100
	}
	stats["success_rate"] = successRate

	// Average execution time for completed executions
	var avgDuration *time.Duration
	var result struct {
		AvgDuration *float64 `gorm:"column:avg_duration"`
	}

	err := query.Session(&gorm.Session{}).Select("AVG(EXTRACT(EPOCH FROM (completed_at - started_at))) as avg_duration").Where("status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", models.ExecutionStatusCompleted).Scan(&result).Error
	if err == nil && result.AvgDuration != nil {
		duration := time.Duration(*result.AvgDuration) * time.Second
		avgDuration = &duration
	}
	stats["average_duration"] = avgDuration

	return stats, nil
}

// CleanupOldExecutions deletes executions older than the specified duration
func (r *executionRepository) CleanupOldExecutions(olderThan time.Duration) (int64, error) {
	threshold := time.Now().UTC().Add(-olderThan)

	result := r.db.Where("created_at < ? AND status IN ?", threshold, []models.ExecutionStatus{
		models.ExecutionStatusCompleted,
		models.ExecutionStatusFailed,
		models.ExecutionStatusCancelled,
	}).Delete(&models.Execution{})

	return result.RowsAffected, result.Error
}

// GetRetryableExecutions retrieves failed executions that can be retried
func (r *executionRepository) GetRetryableExecutions(limit int) ([]*models.Execution, error) {
	var executions []*models.Execution
	err := r.db.Where("status = ? AND parent_execution_id IS NULL", models.ExecutionStatusFailed).Limit(limit).Order("created_at DESC").Find(&executions).Error
	return executions, err
}

// GetChildExecutions retrieves child executions (retries) for a parent execution
func (r *executionRepository) GetChildExecutions(parentID uuid.UUID) ([]*models.Execution, error) {
	var executions []*models.Execution
	err := r.db.Where("parent_execution_id = ?", parentID).Order("created_at ASC").Find(&executions).Error
	return executions, err
}