package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// stepExecutionRepository implements StepExecutionRepository interface
type stepExecutionRepository struct {
	db *gorm.DB
}

// NewStepExecutionRepository creates a new step execution repository
func NewStepExecutionRepository(db *gorm.DB) StepExecutionRepository {
	return &stepExecutionRepository{db: db}
}

// Create creates a new step execution
func (r *stepExecutionRepository) Create(stepExecution *models.StepExecution) error {
	return r.db.Create(stepExecution).Error
}

// GetByID retrieves a step execution by ID
func (r *stepExecutionRepository) GetByID(id uuid.UUID) (*models.StepExecution, error) {
	var stepExecution models.StepExecution
	err := r.db.Where("id = ?", id).First(&stepExecution).Error
	if err != nil {
		return nil, err
	}
	return &stepExecution, nil
}

// GetByExecutionID retrieves all step executions for a specific execution
func (r *stepExecutionRepository) GetByExecutionID(executionID uuid.UUID) ([]*models.StepExecution, error) {
	var stepExecutions []*models.StepExecution
	err := r.db.Where("execution_id = ?", executionID).Order("step_order ASC, created_at ASC").Find(&stepExecutions).Error
	return stepExecutions, err
}

// GetByExecutionIDAndStepID retrieves a specific step execution
func (r *stepExecutionRepository) GetByExecutionIDAndStepID(executionID uuid.UUID, stepID string) (*models.StepExecution, error) {
	var stepExecution models.StepExecution
	err := r.db.Where("execution_id = ? AND step_id = ?", executionID, stepID).First(&stepExecution).Error
	if err != nil {
		return nil, err
	}
	return &stepExecution, nil
}

// Update updates a step execution
func (r *stepExecutionRepository) Update(stepExecution *models.StepExecution) error {
	return r.db.Save(stepExecution).Error
}

// Delete deletes a step execution
func (r *stepExecutionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.StepExecution{}, "id = ?", id).Error
}

// UpdateStatus updates only the status of a step execution
func (r *stepExecutionRepository) UpdateStatus(id uuid.UUID, status models.StepExecutionStatus) error {
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	// Set completion time for terminal states
	if status == models.StepExecutionStatusCompleted || status == models.StepExecutionStatusFailed || status == models.StepExecutionStatusSkipped {
		updateData["completed_at"] = time.Now().UTC()
	}

	// Set start time for running state
	if status == models.StepExecutionStatusRunning {
		updateData["started_at"] = time.Now().UTC()
	}

	return r.db.Model(&models.StepExecution{}).Where("id = ?", id).Updates(updateData).Error
}

// GetByStatus retrieves step executions by status
func (r *stepExecutionRepository) GetByStatus(status models.StepExecutionStatus, limit, offset int) ([]*models.StepExecution, int64, error) {
	var stepExecutions []*models.StepExecution
	var total int64

	query := r.db.Model(&models.StepExecution{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get step executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&stepExecutions).Error
	if err != nil {
		return nil, 0, err
	}

	return stepExecutions, total, nil
}

// GetActiveStepExecutions retrieves all active (running/pending) step executions
func (r *stepExecutionRepository) GetActiveStepExecutions() ([]*models.StepExecution, error) {
	var stepExecutions []*models.StepExecution
	err := r.db.Where("status IN ?", []models.StepExecutionStatus{
		models.StepExecutionStatusPending,
		models.StepExecutionStatusRunning,
	}).Find(&stepExecutions).Error
	return stepExecutions, err
}

// GetFailedStepExecutions retrieves failed step executions within a time range
func (r *stepExecutionRepository) GetFailedStepExecutions(startTime, endTime *time.Time, limit, offset int) ([]*models.StepExecution, int64, error) {
	var stepExecutions []*models.StepExecution
	var total int64

	query := r.db.Model(&models.StepExecution{}).Where("status = ?", models.StepExecutionStatusFailed)

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

	// Get step executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&stepExecutions).Error
	if err != nil {
		return nil, 0, err
	}

	return stepExecutions, total, nil
}

// GetStepExecutionsByTimeRange retrieves step executions within a time range
func (r *stepExecutionRepository) GetStepExecutionsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.StepExecution, int64, error) {
	var stepExecutions []*models.StepExecution
	var total int64

	query := r.db.Model(&models.StepExecution{})

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

	// Get step executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&stepExecutions).Error
	if err != nil {
		return nil, 0, err
	}

	return stepExecutions, total, nil
}

// GetLongRunningStepExecutions retrieves step executions that have been running for longer than the specified duration
func (r *stepExecutionRepository) GetLongRunningStepExecutions(duration time.Duration) ([]*models.StepExecution, error) {
	var stepExecutions []*models.StepExecution
	threshold := time.Now().UTC().Add(-duration)

	err := r.db.Where("status = ? AND started_at < ?", models.StepExecutionStatusRunning, threshold).Find(&stepExecutions).Error
	return stepExecutions, err
}

// GetStepExecutionStats retrieves step execution statistics
func (r *stepExecutionRepository) GetStepExecutionStats(executionID *uuid.UUID, stepType string, startTime, endTime *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.StepExecution{})

	if executionID != nil {
		query = query.Where("execution_id = ?", *executionID)
	}

	if stepType != "" {
		query = query.Where("step_type = ?", stepType)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Count by status
	statusStats := make(map[string]int64)
	statuses := []models.StepExecutionStatus{
		models.StepExecutionStatusPending,
		models.StepExecutionStatusRunning,
		models.StepExecutionStatusCompleted,
		models.StepExecutionStatusFailed,
		models.StepExecutionStatusSkipped,
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
		successRate = float64(statusStats[string(models.StepExecutionStatusCompleted)]) / float64(total) * 100
	}
	stats["success_rate"] = successRate

	// Average execution time for completed step executions
	var avgDuration *time.Duration
	var result struct {
		AvgDuration *float64 `gorm:"column:avg_duration"`
	}

	err := query.Session(&gorm.Session{}).Select("AVG(EXTRACT(EPOCH FROM (completed_at - started_at))) as avg_duration").Where("status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", models.StepExecutionStatusCompleted).Scan(&result).Error
	if err == nil && result.AvgDuration != nil {
		duration := time.Duration(*result.AvgDuration) * time.Second
		avgDuration = &duration
	}
	stats["average_duration"] = avgDuration

	// Count by step type
	if stepType == "" {
		typeStats := make(map[string]int64)
		var typeResults []struct {
			StepType string `gorm:"column:step_type"`
			Count    int64  `gorm:"column:count"`
		}

		err := query.Session(&gorm.Session{}).Select("step_type, COUNT(*) as count").Group("step_type").Scan(&typeResults).Error
		if err == nil {
			for _, result := range typeResults {
				typeStats[result.StepType] = result.Count
			}
		}
		stats["by_type"] = typeStats
	}

	return stats, nil
}

// GetStepExecutionsByStepType retrieves step executions by step type
func (r *stepExecutionRepository) GetStepExecutionsByStepType(stepType string, limit, offset int) ([]*models.StepExecution, int64, error) {
	var stepExecutions []*models.StepExecution
	var total int64

	query := r.db.Model(&models.StepExecution{}).Where("step_type = ?", stepType)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get step executions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&stepExecutions).Error
	if err != nil {
		return nil, 0, err
	}

	return stepExecutions, total, nil
}

// GetRetryableStepExecutions retrieves failed step executions that can be retried
func (r *stepExecutionRepository) GetRetryableStepExecutions(limit int) ([]*models.StepExecution, error) {
	var stepExecutions []*models.StepExecution
	err := r.db.Where("status = ? AND retry_count < max_retries", models.StepExecutionStatusFailed).Limit(limit).Order("created_at DESC").Find(&stepExecutions).Error
	return stepExecutions, err
}

// IncrementRetryCount increments the retry count for a step execution
func (r *stepExecutionRepository) IncrementRetryCount(id uuid.UUID) error {
	return r.db.Model(&models.StepExecution{}).Where("id = ?", id).UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

// GetStepExecutionsByExecutionIDAndStatus retrieves step executions by execution ID and status
func (r *stepExecutionRepository) GetStepExecutionsByExecutionIDAndStatus(executionID uuid.UUID, status models.StepExecutionStatus) ([]*models.StepExecution, error) {
	var stepExecutions []*models.StepExecution
	err := r.db.Where("execution_id = ? AND status = ?", executionID, status).Order("step_order ASC, created_at ASC").Find(&stepExecutions).Error
	return stepExecutions, err
}

// GetPendingStepExecutions retrieves pending step executions for an execution
func (r *stepExecutionRepository) GetPendingStepExecutions(executionID uuid.UUID) ([]*models.StepExecution, error) {
	return r.GetStepExecutionsByExecutionIDAndStatus(executionID, models.StepExecutionStatusPending)
}

// GetRunningStepExecutions retrieves running step executions for an execution
func (r *stepExecutionRepository) GetRunningStepExecutions(executionID uuid.UUID) ([]*models.StepExecution, error) {
	return r.GetStepExecutionsByExecutionIDAndStatus(executionID, models.StepExecutionStatusRunning)
}

// GetCompletedStepExecutions retrieves completed step executions for an execution
func (r *stepExecutionRepository) GetCompletedStepExecutions(executionID uuid.UUID) ([]*models.StepExecution, error) {
	return r.GetStepExecutionsByExecutionIDAndStatus(executionID, models.StepExecutionStatusCompleted)
}

// GetFailedStepExecutionsForExecution retrieves failed step executions for an execution
func (r *stepExecutionRepository) GetFailedStepExecutionsForExecution(executionID uuid.UUID) ([]*models.StepExecution, error) {
	return r.GetStepExecutionsByExecutionIDAndStatus(executionID, models.StepExecutionStatusFailed)
}

// CleanupOldStepExecutions deletes step executions older than the specified duration
func (r *stepExecutionRepository) CleanupOldStepExecutions(olderThan time.Duration) (int64, error) {
	threshold := time.Now().UTC().Add(-olderThan)

	result := r.db.Where("created_at < ? AND status IN ?", threshold, []models.StepExecutionStatus{
		models.StepExecutionStatusCompleted,
		models.StepExecutionStatusFailed,
		models.StepExecutionStatusSkipped,
	}).Delete(&models.StepExecution{})

	return result.RowsAffected, result.Error
}

// BulkUpdateStatus updates the status of multiple step executions
func (r *stepExecutionRepository) BulkUpdateStatus(ids []uuid.UUID, status models.StepExecutionStatus) error {
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	// Set completion time for terminal states
	if status == models.StepExecutionStatusCompleted || status == models.StepExecutionStatusFailed || status == models.StepExecutionStatusSkipped {
		updateData["completed_at"] = time.Now().UTC()
	}

	// Set start time for running state
	if status == models.StepExecutionStatusRunning {
		updateData["started_at"] = time.Now().UTC()
	}

	return r.db.Model(&models.StepExecution{}).Where("id IN ?", ids).Updates(updateData).Error
}

// GetStepExecutionProgress calculates the progress of step executions for an execution
func (r *stepExecutionRepository) GetStepExecutionProgress(executionID uuid.UUID) (map[string]interface{}, error) {
	progress := make(map[string]interface{})

	// Count total steps
	var total int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ?", executionID).Count(&total).Error; err != nil {
		return nil, err
	}

	// Count completed steps
	var completed int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ? AND status = ?", executionID, models.StepExecutionStatusCompleted).Count(&completed).Error; err != nil {
		return nil, err
	}

	// Count failed steps
	var failed int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ? AND status = ?", executionID, models.StepExecutionStatusFailed).Count(&failed).Error; err != nil {
		return nil, err
	}

	// Count running steps
	var running int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ? AND status = ?", executionID, models.StepExecutionStatusRunning).Count(&running).Error; err != nil {
		return nil, err
	}

	// Count pending steps
	var pending int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ? AND status = ?", executionID, models.StepExecutionStatusPending).Count(&pending).Error; err != nil {
		return nil, err
	}

	// Count skipped steps
	var skipped int64
	if err := r.db.Model(&models.StepExecution{}).Where("execution_id = ? AND status = ?", executionID, models.StepExecutionStatusSkipped).Count(&skipped).Error; err != nil {
		return nil, err
	}

	progress["total"] = total
	progress["completed"] = completed
	progress["failed"] = failed
	progress["running"] = running
	progress["pending"] = pending
	progress["skipped"] = skipped

	// Calculate percentage
	percentage := 0.0
	if total > 0 {
		percentage = float64(completed+failed+skipped) / float64(total) * 100
	}
	progress["percentage"] = percentage

	return progress, nil
}