package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// executionEventRepository implements ExecutionEventRepository interface
type executionEventRepository struct {
	db *gorm.DB
}

// NewExecutionEventRepository creates a new execution event repository
func NewExecutionEventRepository(db *gorm.DB) ExecutionEventRepository {
	return &executionEventRepository{db: db}
}

// Create creates a new execution event
func (r *executionEventRepository) Create(event *models.ExecutionEvent) error {
	return r.db.Create(event).Error
}

// GetByID retrieves an execution event by ID
func (r *executionEventRepository) GetByID(id uuid.UUID) (*models.ExecutionEvent, error) {
	var event models.ExecutionEvent
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByExecutionID retrieves all events for a specific execution
func (r *executionEventRepository) GetByExecutionID(executionID uuid.UUID, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("execution_id = ?", executionID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at ASC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetByWorkflowID retrieves events for all executions of a specific workflow
func (r *executionEventRepository) GetByWorkflowID(workflowID uuid.UUID, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("workflow_id = ?", workflowID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetByEventType retrieves events by event type
func (r *executionEventRepository) GetByEventType(eventType models.ExecutionEventType, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("event_type = ?", eventType)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetByTimeRange retrieves events within a time range
func (r *executionEventRepository) GetByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{})

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

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// List retrieves events with filtering and pagination
func (r *executionEventRepository) List(limit, offset int, workflowID, executionID *uuid.UUID, eventType string, startTime, endTime *time.Time) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{})

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if executionID != nil {
		query = query.Where("execution_id = ?", *executionID)
	}

	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

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

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// Delete deletes an execution event
func (r *executionEventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExecutionEvent{}, "id = ?", id).Error
}

// GetErrorEvents retrieves error events within a time range
func (r *executionEventRepository) GetErrorEvents(startTime, endTime *time.Time, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("event_type IN ?", []models.ExecutionEventType{
		models.ExecutionEventTypeExecutionFailed,
		models.ExecutionEventTypeStepFailed,
		models.ExecutionEventTypeError,
	})

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

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetEventsByStepID retrieves events for a specific step
func (r *executionEventRepository) GetEventsByStepID(executionID uuid.UUID, stepID string, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("execution_id = ? AND step_id = ?", executionID, stepID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at ASC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetLatestEventByType retrieves the latest event of a specific type for an execution
func (r *executionEventRepository) GetLatestEventByType(executionID uuid.UUID, eventType models.ExecutionEventType) (*models.ExecutionEvent, error) {
	var event models.ExecutionEvent
	err := r.db.Where("execution_id = ? AND event_type = ?", executionID, eventType).Order("created_at DESC").First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEventStats retrieves event statistics
func (r *executionEventRepository) GetEventStats(workflowID, executionID *uuid.UUID, startTime, endTime *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.ExecutionEvent{})

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if executionID != nil {
		query = query.Where("execution_id = ?", *executionID)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Count by event type
	typeStats := make(map[string]int64)
	eventTypes := []models.ExecutionEventType{
		models.ExecutionEventTypeExecutionStarted,
		models.ExecutionEventTypeExecutionCompleted,
		models.ExecutionEventTypeExecutionFailed,
		models.ExecutionEventTypeExecutionCancelled,
		models.ExecutionEventTypeStepStarted,
		models.ExecutionEventTypeStepCompleted,
		models.ExecutionEventTypeStepFailed,
		models.ExecutionEventTypeStepSkipped,
		models.ExecutionEventTypeError,
		models.ExecutionEventTypeWarning,
		models.ExecutionEventTypeInfo,
	}

	for _, eventType := range eventTypes {
		var count int64
		typeQuery := query.Session(&gorm.Session{})
		if err := typeQuery.Where("event_type = ?", eventType).Count(&count).Error; err != nil {
			return nil, err
		}
		typeStats[string(eventType)] = count
	}

	stats["by_type"] = typeStats

	// Total count
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Error rate
	errorCount := typeStats[string(models.ExecutionEventTypeExecutionFailed)] + typeStats[string(models.ExecutionEventTypeStepFailed)] + typeStats[string(models.ExecutionEventTypeError)]
	errorRate := 0.0
	if total > 0 {
		errorRate = float64(errorCount) / float64(total) * 100
	}
	stats["error_rate"] = errorRate

	return stats, nil
}

// GetRecentEvents retrieves the most recent events
func (r *executionEventRepository) GetRecentEvents(limit int) ([]*models.ExecutionEvent, error) {
	var events []*models.ExecutionEvent
	err := r.db.Limit(limit).Order("created_at DESC").Find(&events).Error
	return events, err
}

// GetEventsByLevel retrieves events by log level (info, warning, error)
func (r *executionEventRepository) GetEventsByLevel(level string, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	var eventTypes []models.ExecutionEventType
	switch level {
	case "error":
		eventTypes = []models.ExecutionEventType{
			models.ExecutionEventTypeExecutionFailed,
			models.ExecutionEventTypeStepFailed,
			models.ExecutionEventTypeError,
		}
	case "warning":
		eventTypes = []models.ExecutionEventType{
			models.ExecutionEventTypeWarning,
		}
	case "info":
		eventTypes = []models.ExecutionEventType{
			models.ExecutionEventTypeExecutionStarted,
			models.ExecutionEventTypeExecutionCompleted,
			models.ExecutionEventTypeStepStarted,
			models.ExecutionEventTypeStepCompleted,
			models.ExecutionEventTypeStepSkipped,
			models.ExecutionEventTypeInfo,
		}
	default:
		// Return all events if level is not recognized
		eventTypes = []models.ExecutionEventType{
			models.ExecutionEventTypeExecutionStarted,
			models.ExecutionEventTypeExecutionCompleted,
			models.ExecutionEventTypeExecutionFailed,
			models.ExecutionEventTypeExecutionCancelled,
			models.ExecutionEventTypeStepStarted,
			models.ExecutionEventTypeStepCompleted,
			models.ExecutionEventTypeStepFailed,
			models.ExecutionEventTypeStepSkipped,
			models.ExecutionEventTypeError,
			models.ExecutionEventTypeWarning,
			models.ExecutionEventTypeInfo,
		}
	}

	query := r.db.Model(&models.ExecutionEvent{}).Where("event_type IN ?", eventTypes)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// SearchEvents searches events by message content
func (r *executionEventRepository) SearchEvents(searchTerm string, limit, offset int) ([]*models.ExecutionEvent, int64, error) {
	var events []*models.ExecutionEvent
	var total int64

	query := r.db.Model(&models.ExecutionEvent{}).Where("message ILIKE ?", "%"+searchTerm+"%")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get events with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// CleanupOldEvents deletes events older than the specified duration
func (r *executionEventRepository) CleanupOldEvents(olderThan time.Duration) (int64, error) {
	threshold := time.Now().UTC().Add(-olderThan)

	result := r.db.Where("created_at < ?", threshold).Delete(&models.ExecutionEvent{})

	return result.RowsAffected, result.Error
}

// BulkCreate creates multiple execution events in a single transaction
func (r *executionEventRepository) BulkCreate(events []*models.ExecutionEvent) error {
	if len(events) == 0 {
		return nil
	}

	return r.db.CreateInBatches(events, 100).Error
}

// GetEventTimeline retrieves events for creating a timeline view
func (r *executionEventRepository) GetEventTimeline(executionID uuid.UUID) ([]*models.ExecutionEvent, error) {
	var events []*models.ExecutionEvent
	err := r.db.Where("execution_id = ?", executionID).Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetEventsByExecutionAndStep retrieves events for a specific execution and step
func (r *executionEventRepository) GetEventsByExecutionAndStep(executionID uuid.UUID, stepID string) ([]*models.ExecutionEvent, error) {
	var events []*models.ExecutionEvent
	err := r.db.Where("execution_id = ? AND step_id = ?", executionID, stepID).Order("created_at ASC").Find(&events).Error
	return events, err
}

// CountEventsByType counts events by type within a time range
func (r *executionEventRepository) CountEventsByType(eventType models.ExecutionEventType, startTime, endTime *time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.ExecutionEvent{}).Where("event_type = ?", eventType)

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetEventFrequency retrieves event frequency data for charts
func (r *executionEventRepository) GetEventFrequency(startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Determine the date truncation based on interval
	var dateTrunc string
	switch interval {
	case "hour":
		dateTrunc = "date_trunc('hour', created_at)"
	case "day":
		dateTrunc = "date_trunc('day', created_at)"
	case "week":
		dateTrunc = "date_trunc('week', created_at)"
	case "month":
		dateTrunc = "date_trunc('month', created_at)"
	default:
		dateTrunc = "date_trunc('hour', created_at)"
	}

	query := r.db.Model(&models.ExecutionEvent{}).Select(dateTrunc + " as time_bucket, event_type, COUNT(*) as count")

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query = query.Group("time_bucket, event_type").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time                   `gorm:"column:time_bucket"`
		EventType  models.ExecutionEventType `gorm:"column:event_type"`
		Count      int64                      `gorm:"column:count"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":       result.TimeBucket,
			"event_type": string(result.EventType),
			"count":      result.Count,
		})
	}

	return results, nil
}