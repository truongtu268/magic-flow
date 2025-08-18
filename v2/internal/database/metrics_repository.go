package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// metricsRepository implements MetricsRepository interface
type metricsRepository struct {
	db *gorm.DB
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository(db *gorm.DB) MetricsRepository {
	return &metricsRepository{db: db}
}

// CreateWorkflowMetric creates a new workflow metric
func (r *metricsRepository) CreateWorkflowMetric(metric *models.WorkflowMetric) error {
	return r.db.Create(metric).Error
}

// CreateSystemMetric creates a new system metric
func (r *metricsRepository) CreateSystemMetric(metric *models.SystemMetric) error {
	return r.db.Create(metric).Error
}

// CreateBusinessMetric creates a new business metric
func (r *metricsRepository) CreateBusinessMetric(metric *models.BusinessMetric) error {
	return r.db.Create(metric).Error
}

// GetWorkflowMetrics retrieves workflow metrics with filtering
func (r *metricsRepository) GetWorkflowMetrics(workflowID *uuid.UUID, metricType string, startTime, endTime *time.Time, limit, offset int) ([]*models.WorkflowMetric, int64, error) {
	var metrics []*models.WorkflowMetric
	var total int64

	query := r.db.Model(&models.WorkflowMetric{})

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get metrics with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&metrics).Error
	if err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}

// GetSystemMetrics retrieves system metrics with filtering
func (r *metricsRepository) GetSystemMetrics(metricType string, startTime, endTime *time.Time, limit, offset int) ([]*models.SystemMetric, int64, error) {
	var metrics []*models.SystemMetric
	var total int64

	query := r.db.Model(&models.SystemMetric{})

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get metrics with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&metrics).Error
	if err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}

// GetBusinessMetrics retrieves business metrics with filtering
func (r *metricsRepository) GetBusinessMetrics(metricName string, startTime, endTime *time.Time, limit, offset int) ([]*models.BusinessMetric, int64, error) {
	var metrics []*models.BusinessMetric
	var total int64

	query := r.db.Model(&models.BusinessMetric{})

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get metrics with pagination
	err := query.Limit(limit).Offset(offset).Order("timestamp DESC").Find(&metrics).Error
	if err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}

// GetAggregatedWorkflowMetrics retrieves aggregated workflow metrics
func (r *metricsRepository) GetAggregatedWorkflowMetrics(workflowID *uuid.UUID, metricType, aggregation string, startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Determine the date truncation based on interval
	var dateTrunc string
	switch interval {
	case "minute":
		dateTrunc = "date_trunc('minute', timestamp)"
	case "hour":
		dateTrunc = "date_trunc('hour', timestamp)"
	case "day":
		dateTrunc = "date_trunc('day', timestamp)"
	case "week":
		dateTrunc = "date_trunc('week', timestamp)"
	case "month":
		dateTrunc = "date_trunc('month', timestamp)"
	default:
		dateTrunc = "date_trunc('hour', timestamp)"
	}

	// Determine the aggregation function
	var aggFunc string
	switch aggregation {
	case "avg":
		aggFunc = "AVG(value)"
	case "sum":
		aggFunc = "SUM(value)"
	case "min":
		aggFunc = "MIN(value)"
	case "max":
		aggFunc = "MAX(value)"
	case "count":
		aggFunc = "COUNT(*)"
	default:
		aggFunc = "AVG(value)"
	}

	query := r.db.Model(&models.WorkflowMetric{}).Select(dateTrunc + " as time_bucket, " + aggFunc + " as value")

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	query = query.Group("time_bucket").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time `gorm:"column:time_bucket"`
		Value      float64   `gorm:"column:value"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":  result.TimeBucket,
			"value": result.Value,
		})
	}

	return results, nil
}

// GetAggregatedSystemMetrics retrieves aggregated system metrics
func (r *metricsRepository) GetAggregatedSystemMetrics(metricType, aggregation string, startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Determine the date truncation based on interval
	var dateTrunc string
	switch interval {
	case "minute":
		dateTrunc = "date_trunc('minute', timestamp)"
	case "hour":
		dateTrunc = "date_trunc('hour', timestamp)"
	case "day":
		dateTrunc = "date_trunc('day', timestamp)"
	case "week":
		dateTrunc = "date_trunc('week', timestamp)"
	case "month":
		dateTrunc = "date_trunc('month', timestamp)"
	default:
		dateTrunc = "date_trunc('hour', timestamp)"
	}

	// Determine the aggregation function
	var aggFunc string
	switch aggregation {
	case "avg":
		aggFunc = "AVG(value)"
	case "sum":
		aggFunc = "SUM(value)"
	case "min":
		aggFunc = "MIN(value)"
	case "max":
		aggFunc = "MAX(value)"
	case "count":
		aggFunc = "COUNT(*)"
	default:
		aggFunc = "AVG(value)"
	}

	query := r.db.Model(&models.SystemMetric{}).Select(dateTrunc + " as time_bucket, " + aggFunc + " as value")

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	query = query.Group("time_bucket").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time `gorm:"column:time_bucket"`
		Value      float64   `gorm:"column:value"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":  result.TimeBucket,
			"value": result.Value,
		})
	}

	return results, nil
}

// GetAggregatedBusinessMetrics retrieves aggregated business metrics
func (r *metricsRepository) GetAggregatedBusinessMetrics(metricName, aggregation string, startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Determine the date truncation based on interval
	var dateTrunc string
	switch interval {
	case "minute":
		dateTrunc = "date_trunc('minute', timestamp)"
	case "hour":
		dateTrunc = "date_trunc('hour', timestamp)"
	case "day":
		dateTrunc = "date_trunc('day', timestamp)"
	case "week":
		dateTrunc = "date_trunc('week', timestamp)"
	case "month":
		dateTrunc = "date_trunc('month', timestamp)"
	default:
		dateTrunc = "date_trunc('hour', timestamp)"
	}

	// Determine the aggregation function
	var aggFunc string
	switch aggregation {
	case "avg":
		aggFunc = "AVG(value)"
	case "sum":
		aggFunc = "SUM(value)"
	case "min":
		aggFunc = "MIN(value)"
	case "max":
		aggFunc = "MAX(value)"
	case "count":
		aggFunc = "COUNT(*)"
	default:
		aggFunc = "AVG(value)"
	}

	query := r.db.Model(&models.BusinessMetric{}).Select(dateTrunc + " as time_bucket, " + aggFunc + " as value")

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	query = query.Group("time_bucket").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time `gorm:"column:time_bucket"`
		Value      float64   `gorm:"column:value"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":  result.TimeBucket,
			"value": result.Value,
		})
	}

	return results, nil
}

// GetMetricAggregations retrieves metric aggregations
func (r *metricsRepository) GetMetricAggregations(limit, offset int) ([]*models.MetricAggregation, int64, error) {
	var aggregations []*models.MetricAggregation
	var total int64

	query := r.db.Model(&models.MetricAggregation{})

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get aggregations with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&aggregations).Error
	if err != nil {
		return nil, 0, err
	}

	return aggregations, total, nil
}

// CreateMetricAggregation creates a new metric aggregation
func (r *metricsRepository) CreateMetricAggregation(aggregation *models.MetricAggregation) error {
	return r.db.Create(aggregation).Error
}

// GetMetricAggregationByID retrieves a metric aggregation by ID
func (r *metricsRepository) GetMetricAggregationByID(id uuid.UUID) (*models.MetricAggregation, error) {
	var aggregation models.MetricAggregation
	err := r.db.Where("id = ?", id).First(&aggregation).Error
	if err != nil {
		return nil, err
	}
	return &aggregation, nil
}

// UpdateMetricAggregation updates a metric aggregation
func (r *metricsRepository) UpdateMetricAggregation(aggregation *models.MetricAggregation) error {
	return r.db.Save(aggregation).Error
}

// DeleteMetricAggregation deletes a metric aggregation
func (r *metricsRepository) DeleteMetricAggregation(id uuid.UUID) error {
	return r.db.Delete(&models.MetricAggregation{}, "id = ?", id).Error
}

// GetLatestSystemMetrics retrieves the latest system metrics for each metric type
func (r *metricsRepository) GetLatestSystemMetrics() ([]*models.SystemMetric, error) {
	var metrics []*models.SystemMetric

	// Get the latest metric for each metric type
	subquery := r.db.Model(&models.SystemMetric{}).Select("metric_type, MAX(timestamp) as max_timestamp").Group("metric_type")

	err := r.db.Table("system_metrics sm").Joins("INNER JOIN (?) latest ON sm.metric_type = latest.metric_type AND sm.timestamp = latest.max_timestamp", subquery).Find(&metrics).Error

	return metrics, err
}

// GetWorkflowMetricSummary retrieves workflow metric summary
func (r *metricsRepository) GetWorkflowMetricSummary(workflowID uuid.UUID, startTime, endTime *time.Time) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	query := r.db.Model(&models.WorkflowMetric{}).Where("workflow_id = ?", workflowID)

	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	// Get execution count
	var executionCount int64
	if err := query.Where("metric_type = ?", "execution_count").Count(&executionCount).Error; err != nil {
		return nil, err
	}
	summary["execution_count"] = executionCount

	// Get average execution time
	var avgExecutionTime struct {
		AvgTime *float64 `gorm:"column:avg_time"`
	}
	if err := query.Select("AVG(value) as avg_time").Where("metric_type = ?", "execution_time").Scan(&avgExecutionTime).Error; err != nil {
		return nil, err
	}
	summary["avg_execution_time"] = avgExecutionTime.AvgTime

	// Get success rate
	var successCount, totalCount int64
	if err := query.Where("metric_type = ? AND value = ?", "execution_status", 1).Count(&successCount).Error; err != nil {
		return nil, err
	}
	if err := query.Where("metric_type = ?", "execution_status").Count(&totalCount).Error; err != nil {
		return nil, err
	}

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successCount) / float64(totalCount) * 100
	}
	summary["success_rate"] = successRate

	return summary, nil
}

// GetSystemHealthMetrics retrieves system health metrics
func (r *metricsRepository) GetSystemHealthMetrics() (map[string]interface{}, error) {
	health := make(map[string]interface{})

	// Get latest CPU usage
	var cpuUsage struct {
		Value *float64 `gorm:"column:value"`
	}
	if err := r.db.Model(&models.SystemMetric{}).Select("value").Where("metric_type = ?", "cpu_usage").Order("timestamp DESC").Limit(1).Scan(&cpuUsage).Error; err == nil {
		health["cpu_usage"] = cpuUsage.Value
	}

	// Get latest memory usage
	var memoryUsage struct {
		Value *float64 `gorm:"column:value"`
	}
	if err := r.db.Model(&models.SystemMetric{}).Select("value").Where("metric_type = ?", "memory_usage").Order("timestamp DESC").Limit(1).Scan(&memoryUsage).Error; err == nil {
		health["memory_usage"] = memoryUsage.Value
	}

	// Get latest disk usage
	var diskUsage struct {
		Value *float64 `gorm:"column:value"`
	}
	if err := r.db.Model(&models.SystemMetric{}).Select("value").Where("metric_type = ?", "disk_usage").Order("timestamp DESC").Limit(1).Scan(&diskUsage).Error; err == nil {
		health["disk_usage"] = diskUsage.Value
	}

	// Get active connections
	var activeConnections struct {
		Value *float64 `gorm:"column:value"`
	}
	if err := r.db.Model(&models.SystemMetric{}).Select("value").Where("metric_type = ?", "active_connections").Order("timestamp DESC").Limit(1).Scan(&activeConnections).Error; err == nil {
		health["active_connections"] = activeConnections.Value
	}

	return health, nil
}

// CleanupOldMetrics deletes metrics older than the specified duration
func (r *metricsRepository) CleanupOldMetrics(olderThan time.Duration) (int64, error) {
	threshold := time.Now().UTC().Add(-olderThan)

	// Cleanup workflow metrics
	workflowResult := r.db.Where("timestamp < ?", threshold).Delete(&models.WorkflowMetric{})
	workflowDeleted := workflowResult.RowsAffected

	// Cleanup system metrics
	systemResult := r.db.Where("timestamp < ?", threshold).Delete(&models.SystemMetric{})
	systemDeleted := systemResult.RowsAffected

	// Cleanup business metrics
	businessResult := r.db.Where("timestamp < ?", threshold).Delete(&models.BusinessMetric{})
	businessDeleted := businessResult.RowsAffected

	totalDeleted := workflowDeleted + systemDeleted + businessDeleted

	// Return error if any of the deletions failed
	if workflowResult.Error != nil {
		return totalDeleted, workflowResult.Error
	}
	if systemResult.Error != nil {
		return totalDeleted, systemResult.Error
	}
	if businessResult.Error != nil {
		return totalDeleted, businessResult.Error
	}

	return totalDeleted, nil
}

// BulkCreateWorkflowMetrics creates multiple workflow metrics in batches
func (r *metricsRepository) BulkCreateWorkflowMetrics(metrics []*models.WorkflowMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	return r.db.CreateInBatches(metrics, 100).Error
}

// BulkCreateSystemMetrics creates multiple system metrics in batches
func (r *metricsRepository) BulkCreateSystemMetrics(metrics []*models.SystemMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	return r.db.CreateInBatches(metrics, 100).Error
}

// BulkCreateBusinessMetrics creates multiple business metrics in batches
func (r *metricsRepository) BulkCreateBusinessMetrics(metrics []*models.BusinessMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	return r.db.CreateInBatches(metrics, 100).Error
}

// GetMetricTypes retrieves all unique metric types
func (r *metricsRepository) GetMetricTypes() (map[string][]string, error) {
	types := make(map[string][]string)

	// Get workflow metric types
	var workflowTypes []string
	if err := r.db.Model(&models.WorkflowMetric{}).Distinct("metric_type").Pluck("metric_type", &workflowTypes).Error; err != nil {
		return nil, err
	}
	types["workflow"] = workflowTypes

	// Get system metric types
	var systemTypes []string
	if err := r.db.Model(&models.SystemMetric{}).Distinct("metric_type").Pluck("metric_type", &systemTypes).Error; err != nil {
		return nil, err
	}
	types["system"] = systemTypes

	// Get business metric names
	var businessNames []string
	if err := r.db.Model(&models.BusinessMetric{}).Distinct("metric_name").Pluck("metric_name", &businessNames).Error; err != nil {
		return nil, err
	}
	types["business"] = businessNames

	return types, nil
}

// GetMetricStatistics retrieves statistical information about metrics
func (r *metricsRepository) GetMetricStatistics(startTime, endTime *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.WorkflowMetric{})
	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}

	// Workflow metrics count
	var workflowCount int64
	if err := query.Count(&workflowCount).Error; err != nil {
		return nil, err
	}
	stats["workflow_metrics_count"] = workflowCount

	// System metrics count
	systemQuery := r.db.Model(&models.SystemMetric{})
	if startTime != nil {
		systemQuery = systemQuery.Where("timestamp >= ?", *startTime)
	}
	if endTime != nil {
		systemQuery = systemQuery.Where("timestamp <= ?", *endTime)
	}

	var systemCount int64
	if err := systemQuery.Count(&systemCount).Error; err != nil {
		return nil, err
	}
	stats["system_metrics_count"] = systemCount

	// Business metrics count
	businessQuery := r.db.Model(&models.BusinessMetric{})
	if startTime != nil {
		businessQuery = businessQuery.Where("timestamp >= ?", *startTime)
	}
	if endTime != nil {
		businessQuery = businessQuery.Where("timestamp <= ?", *endTime)
	}

	var businessCount int64
	if err := businessQuery.Count(&businessCount).Error; err != nil {
		return nil, err
	}
	stats["business_metrics_count"] = businessCount

	stats["total_metrics_count"] = workflowCount + systemCount + businessCount

	return stats, nil
}