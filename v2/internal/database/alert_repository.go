package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// alertRepository implements AlertRepository interface
type alertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

// Create creates a new alert
func (r *alertRepository) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

// GetByID retrieves an alert by ID
func (r *alertRepository) GetByID(id uuid.UUID) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Where("id = ?", id).First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// List retrieves alerts with filtering and pagination
func (r *alertRepository) List(limit, offset int, enabled *bool, severity string) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{})

	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// Update updates an alert
func (r *alertRepository) Update(alert *models.Alert) error {
	return r.db.Save(alert).Error
}

// Delete deletes an alert
func (r *alertRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Alert{}, "id = ?", id).Error
}

// GetEnabledAlerts retrieves all enabled alerts
func (r *alertRepository) GetEnabledAlerts() ([]*models.Alert, error) {
	var alerts []*models.Alert
	err := r.db.Where("enabled = ?", true).Find(&alerts).Error
	return alerts, err
}

// GetAlertsBySeverity retrieves alerts by severity level
func (r *alertRepository) GetAlertsBySeverity(severity models.AlertSeverity, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{}).Where("severity = ?", severity)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// GetAlertsByWorkflow retrieves alerts for a specific workflow
func (r *alertRepository) GetAlertsByWorkflow(workflowID uuid.UUID, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	// Note: This assumes alerts can be associated with workflows
	// The actual implementation might need to be adjusted based on the Alert model structure
	query := r.db.Model(&models.Alert{}).Where("workflow_id = ?", workflowID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// GetAlertsByCreator retrieves alerts created by a specific user
func (r *alertRepository) GetAlertsByCreator(createdBy string, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// SearchAlerts searches alerts by name or description
func (r *alertRepository) SearchAlerts(searchTerm string, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{}).Where("name ILIKE ? OR description ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// EnableAlert enables an alert
func (r *alertRepository) EnableAlert(id uuid.UUID) error {
	return r.db.Model(&models.Alert{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableAlert disables an alert
func (r *alertRepository) DisableAlert(id uuid.UUID) error {
	return r.db.Model(&models.Alert{}).Where("id = ?", id).Update("enabled", false).Error
}

// UpdateLastTriggered updates the last triggered time for an alert
func (r *alertRepository) UpdateLastTriggered(id uuid.UUID, triggeredAt time.Time) error {
	return r.db.Model(&models.Alert{}).Where("id = ?", id).Update("last_triggered_at", triggeredAt).Error
}

// GetAlertStats retrieves alert statistics
func (r *alertRepository) GetAlertStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total alerts
	var total int64
	if err := r.db.Model(&models.Alert{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Enabled alerts
	var enabled int64
	if err := r.db.Model(&models.Alert{}).Where("enabled = ?", true).Count(&enabled).Error; err != nil {
		return nil, err
	}
	stats["enabled"] = enabled

	// Disabled alerts
	stats["disabled"] = total - enabled

	// Count by severity
	severityStats := make(map[string]int64)
	severities := []models.AlertSeverity{
		models.AlertSeverityLow,
		models.AlertSeverityMedium,
		models.AlertSeverityHigh,
		models.AlertSeverityCritical,
	}

	for _, severity := range severities {
		var count int64
		if err := r.db.Model(&models.Alert{}).Where("severity = ?", severity).Count(&count).Error; err != nil {
			return nil, err
		}
		severityStats[string(severity)] = count
	}
	stats["by_severity"] = severityStats

	// Recently triggered alerts (last 24 hours)
	twentyFourHoursAgo := time.Now().UTC().Add(-24 * time.Hour)
	var recentlyTriggered int64
	if err := r.db.Model(&models.Alert{}).Where("last_triggered_at >= ?", twentyFourHoursAgo).Count(&recentlyTriggered).Error; err != nil {
		return nil, err
	}
	stats["recently_triggered_24h"] = recentlyTriggered

	return stats, nil
}

// GetRecentlyTriggeredAlerts retrieves alerts that were triggered recently
func (r *alertRepository) GetRecentlyTriggeredAlerts(since time.Time, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{}).Where("last_triggered_at >= ?", since)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("last_triggered_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// GetAlertsByTimeRange retrieves alerts created within a time range
func (r *alertRepository) GetAlertsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	query := r.db.Model(&models.Alert{})

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

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// BulkUpdateAlertStatus updates the enabled status of multiple alerts
func (r *alertRepository) BulkUpdateAlertStatus(ids []uuid.UUID, enabled bool) error {
	return r.db.Model(&models.Alert{}).Where("id IN ?", ids).Update("enabled", enabled).Error
}

// GetAlertsByMetricType retrieves alerts for a specific metric type
func (r *alertRepository) GetAlertsByMetricType(metricType string, limit, offset int) ([]*models.Alert, int64, error) {
	var alerts []*models.Alert
	var total int64

	// Note: This assumes alerts have a metric_type field or similar
	// The actual implementation might need to be adjusted based on the Alert model structure
	query := r.db.Model(&models.Alert{}).Where("conditions ->> 'metric' = ?", metricType)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get alerts with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// CreateAlertEvent creates a new alert event
func (r *alertRepository) CreateAlertEvent(event *models.AlertEvent) error {
	return r.db.Create(event).Error
}

// GetAlertEvents retrieves events for a specific alert
func (r *alertRepository) GetAlertEvents(alertID uuid.UUID, limit, offset int) ([]*models.AlertEvent, int64, error) {
	var events []*models.AlertEvent
	var total int64

	query := r.db.Model(&models.AlertEvent{}).Where("alert_id = ?", alertID)

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

// GetAlertEventByID retrieves an alert event by ID
func (r *alertRepository) GetAlertEventByID(id uuid.UUID) (*models.AlertEvent, error) {
	var event models.AlertEvent
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetRecentAlertEvents retrieves recent alert events across all alerts
func (r *alertRepository) GetRecentAlertEvents(limit int) ([]*models.AlertEvent, error) {
	var events []*models.AlertEvent
	err := r.db.Limit(limit).Order("created_at DESC").Find(&events).Error
	return events, err
}

// GetAlertEventsByTimeRange retrieves alert events within a time range
func (r *alertRepository) GetAlertEventsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.AlertEvent, int64, error) {
	var events []*models.AlertEvent
	var total int64

	query := r.db.Model(&models.AlertEvent{})

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

// GetAlertEventStats retrieves alert event statistics
func (r *alertRepository) GetAlertEventStats(alertID *uuid.UUID, startTime, endTime *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.AlertEvent{})

	if alertID != nil {
		query = query.Where("alert_id = ?", *alertID)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Total events
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by status
	statusStats := make(map[string]int64)
	statuses := []models.AlertEventStatus{
		models.AlertEventStatusTriggered,
		models.AlertEventStatusResolved,
		models.AlertEventStatusAcknowledged,
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

	return stats, nil
}

// CleanupOldAlertEvents deletes alert events older than the specified duration
func (r *alertRepository) CleanupOldAlertEvents(olderThan time.Duration) (int64, error) {
	threshold := time.Now().UTC().Add(-olderThan)

	result := r.db.Where("created_at < ?", threshold).Delete(&models.AlertEvent{})

	return result.RowsAffected, result.Error
}

// GetAlertTriggerFrequency retrieves alert trigger frequency data
func (r *alertRepository) GetAlertTriggerFrequency(alertID *uuid.UUID, startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
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

	query := r.db.Model(&models.AlertEvent{}).Select(dateTrunc + " as time_bucket, COUNT(*) as count").Where("status = ?", models.AlertEventStatusTriggered)

	if alertID != nil {
		query = query.Where("alert_id = ?", *alertID)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query = query.Group("time_bucket").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time `gorm:"column:time_bucket"`
		Count      int64     `gorm:"column:count"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":  result.TimeBucket,
			"count": result.Count,
		})
	}

	return results, nil
}

// GetMostTriggeredAlerts retrieves the most frequently triggered alerts
func (r *alertRepository) GetMostTriggeredAlerts(limit int, startTime, endTime *time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("alert_events ae").Select("ae.alert_id, a.name, COUNT(*) as trigger_count").Joins("JOIN alerts a ON ae.alert_id = a.id").Where("ae.status = ?", models.AlertEventStatusTriggered)

	if startTime != nil {
		query = query.Where("ae.created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("ae.created_at <= ?", *endTime)
	}

	query = query.Group("ae.alert_id, a.name").Order("trigger_count DESC").Limit(limit)

	var rawResults []struct {
		AlertID      uuid.UUID `gorm:"column:alert_id"`
		Name         string    `gorm:"column:name"`
		TriggerCount int64     `gorm:"column:trigger_count"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"alert_id":      result.AlertID,
			"name":          result.Name,
			"trigger_count": result.TriggerCount,
		})
	}

	return results, nil
}