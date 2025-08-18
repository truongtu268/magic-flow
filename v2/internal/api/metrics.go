package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magic-flow/v2/pkg/models"
	"github.com/sirupsen/logrus"
)

// getWorkflowMetrics gets workflow metrics
func (h *Handler) getWorkflowMetrics(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Parse filters
	filters := map[string]interface{}{
		"start_time": start,
		"end_time":   end,
	}

	if workflowID := c.Query("workflow_id"); workflowID != "" {
		if id, err := uuid.Parse(workflowID); err == nil {
			filters["workflow_id"] = id
		}
	}

	if environment := c.Query("environment"); environment != "" {
		filters["environment"] = environment
	}

	if metricName := c.Query("metric"); metricName != "" {
		filters["name"] = metricName
	}

	aggregation := c.DefaultQuery("aggregation", "avg")
	interval := c.DefaultQuery("interval", "5m")

	// Get workflow metrics
	metrics, err := h.services.MetricsService.GetWorkflowMetrics(filters, aggregation, interval)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get workflow metrics", err)
		return
	}

	h.successResponse(c, metrics)
}

// getWorkflowMetricsById gets metrics for a specific workflow
func (h *Handler) getWorkflowMetricsById(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Get workflow metrics
	metrics, err := h.services.MetricsService.GetWorkflowMetricsByID(id, start, end)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get workflow metrics", err)
		return
	}

	// Get workflow statistics
	stats, err := h.services.MetricsService.GetWorkflowStats(id, start, end)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get workflow stats", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics":   metrics,
		"stats":     stats,
		"timestamp": time.Now().UTC(),
	})
}

// getSystemMetrics gets system metrics
func (h *Handler) getSystemMetrics(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Parse filters
	filters := map[string]interface{}{
		"start_time": start,
		"end_time":   end,
	}

	if component := c.Query("component"); component != "" {
		filters["component"] = component
	}

	if instance := c.Query("instance"); instance != "" {
		filters["instance"] = instance
	}

	if environment := c.Query("environment"); environment != "" {
		filters["environment"] = environment
	}

	if metricName := c.Query("metric"); metricName != "" {
		filters["name"] = metricName
	}

	aggregation := c.DefaultQuery("aggregation", "avg")
	interval := c.DefaultQuery("interval", "5m")

	// Get system metrics
	metrics, err := h.services.MetricsService.GetSystemMetrics(filters, aggregation, interval)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get system metrics", err)
		return
	}

	h.successResponse(c, metrics)
}

// recordCustomMetric records a custom metric
func (h *Handler) recordCustomMetric(c *gin.Context) {
	var request MetricRequest
	if err := h.validateRequestBody(c, &request); err != nil {
		return
	}

	// Validate metric type
	if request.Type == "" {
		request.Type = models.MetricTypeGauge
	}

	// Create business metric
	metric := &models.BusinessMetric{
		Name:        request.Name,
		Type:        request.Type,
		Category:    models.MetricCategoryCustom,
		Value:       request.Value,
		Unit:        request.Unit,
		Labels:      request.Labels,
		Metadata:    request.Metadata,
		WorkflowID:  request.WorkflowID,
		ExecutionID: request.ExecutionID,
		Timestamp:   time.Now(),
	}

	// Record metric
	if err := h.services.MetricsService.RecordBusinessMetric(metric); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to record metric", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"metric_name": request.Name,
		"value":       request.Value,
		"user_id":     h.getUserID(c),
	}).Info("Custom metric recorded")

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Metric recorded successfully",
		"timestamp": time.Now().UTC(),
	})
}

// getCustomMetrics gets custom metrics
func (h *Handler) getCustomMetrics(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Parse filters
	filters := map[string]interface{}{
		"start_time": start,
		"end_time":   end,
		"category":   models.MetricCategoryCustom,
	}

	if metricName := c.Query("metric"); metricName != "" {
		filters["name"] = metricName
	}

	if workflowID := c.Query("workflow_id"); workflowID != "" {
		if id, err := uuid.Parse(workflowID); err == nil {
			filters["workflow_id"] = id
		}
	}

	if executionID := c.Query("execution_id"); executionID != "" {
		if id, err := uuid.Parse(executionID); err == nil {
			filters["execution_id"] = id
		}
	}

	if customerID := c.Query("customer_id"); customerID != "" {
		filters["customer_id"] = customerID
	}

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		filters["tenant_id"] = tenantID
	}

	aggregation := c.DefaultQuery("aggregation", "avg")
	interval := c.DefaultQuery("interval", "5m")

	// Get custom metrics
	metrics, err := h.services.MetricsService.GetBusinessMetrics(filters, aggregation, interval)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get custom metrics", err)
		return
	}

	h.successResponse(c, metrics)
}

// getMetricAggregations gets metric aggregations
func (h *Handler) getMetricAggregations(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Parse filters
	filters := map[string]interface{}{
		"start_time": start,
		"end_time":   end,
	}

	if metricName := c.Query("metric"); metricName != "" {
		filters["name"] = metricName
	}

	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	if aggregation := c.Query("aggregation"); aggregation != "" {
		filters["aggregation"] = aggregation
	}

	if interval := c.Query("interval"); interval != "" {
		filters["interval"] = interval
	}

	// Parse pagination
	page, limit := h.parsePagination(c)

	// Get metric aggregations
	aggregations, total, err := h.services.MetricsService.GetMetricAggregations(filters, page, limit)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get metric aggregations", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       aggregations,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// getDashboardOverview gets dashboard overview data
func (h *Handler) getDashboardOverview(c *gin.Context) {
	// Get overview data from services
	overview, err := h.services.DashboardService.GetOverview()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get dashboard overview", err)
		return
	}

	h.successResponse(c, overview)
}

// getWorkflowStatusSummary gets workflow status summary
func (h *Handler) getWorkflowStatusSummary(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Get workflow status summary
	summary, err := h.services.DashboardService.GetWorkflowStatusSummary(start, end)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get workflow status summary", err)
		return
	}

	h.successResponse(c, summary)
}

// getSystemHealth gets system health status
func (h *Handler) getSystemHealth(c *gin.Context) {
	// Get system health from services
	health, err := h.services.DashboardService.GetSystemHealth()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get system health", err)
		return
	}

	h.successResponse(c, health)
}

// getLiveMetrics gets live metrics for real-time dashboard
func (h *Handler) getLiveMetrics(c *gin.Context) {
	// Get live metrics from metrics collector
	metrics := h.metricsCollector.GetLiveMetrics()

	c.JSON(http.StatusOK, gin.H{
		"metrics":   metrics,
		"timestamp": time.Now().UTC(),
	})
}

// createAlert creates a new alert
func (h *Handler) createAlert(c *gin.Context) {
	var alert models.Alert
	if err := h.validateRequestBody(c, &alert); err != nil {
		return
	}

	// Set creator
	alert.CreatedBy = h.getUserID(c)

	// Create alert
	createdAlert, err := h.services.AlertService.Create(&alert)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create alert", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"alert_id": createdAlert.ID,
		"name":     createdAlert.Name,
		"user_id":  h.getUserID(c),
	}).Info("Alert created")

	c.JSON(http.StatusCreated, gin.H{
		"data":      createdAlert,
		"timestamp": time.Now().UTC(),
	})
}

// listAlerts lists all alerts
func (h *Handler) listAlerts(c *gin.Context) {
	page, limit := h.parsePagination(c)

	// Parse filters
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if severity := c.Query("severity"); severity != "" {
		filters["severity"] = severity
	}
	if enabled := c.Query("enabled"); enabled != "" {
		filters["enabled"] = enabled == "true"
	}

	// Get alerts
	alerts, total, err := h.services.AlertService.List(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list alerts", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       alerts,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// getAlert gets an alert by ID
func (h *Handler) getAlert(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	alert, err := h.services.AlertService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Alert not found", err)
		return
	}

	h.successResponse(c, alert)
}

// updateAlert updates an alert
func (h *Handler) updateAlert(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var updateData models.Alert
	if err := h.validateRequestBody(c, &updateData); err != nil {
		return
	}

	updateData.ID = id

	// Update alert
	updatedAlert, err := h.services.AlertService.Update(&updateData)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update alert", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"alert_id": id,
		"name":     updatedAlert.Name,
		"user_id":  h.getUserID(c),
	}).Info("Alert updated")

	h.successResponse(c, updatedAlert)
}

// deleteAlert deletes an alert
func (h *Handler) deleteAlert(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Delete alert
	if err := h.services.AlertService.Delete(id); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to delete alert", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"alert_id": id,
		"user_id":  h.getUserID(c),
	}).Info("Alert deleted")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Alert deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// enableAlert enables an alert
func (h *Handler) enableAlert(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Enable alert
	if err := h.services.AlertService.Enable(id); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to enable alert", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"alert_id": id,
		"user_id":  h.getUserID(c),
	}).Info("Alert enabled")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Alert enabled successfully",
		"timestamp": time.Now().UTC(),
	})
}

// disableAlert disables an alert
func (h *Handler) disableAlert(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Disable alert
	if err := h.services.AlertService.Disable(id); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to disable alert", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"alert_id": id,
		"user_id":  h.getUserID(c),
	}).Info("Alert disabled")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Alert disabled successfully",
		"timestamp": time.Now().UTC(),
	})
}

// getAlertEvents gets alert events
func (h *Handler) getAlertEvents(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	page, limit := h.parsePagination(c)

	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid time range", err)
		return
	}

	// Get alert events
	events, total, err := h.services.AlertService.GetEvents(id, start, end, page, limit)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get alert events", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       events,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}