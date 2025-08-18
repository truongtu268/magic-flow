package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// DatabaseEventHandler handles workflow events by storing them in the database
type DatabaseEventHandler struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewDatabaseEventHandler creates a new database event handler
func NewDatabaseEventHandler(db *gorm.DB, logger *logrus.Logger) *DatabaseEventHandler {
	return &DatabaseEventHandler{
		db:     db,
		logger: logger,
	}
}

func (h *DatabaseEventHandler) Handle(event *WorkflowEvent) error {
	// Convert event data to JSON
	eventDataJSON, err := json.Marshal(event.Data)
	if err != nil {
		h.logger.WithError(err).Error("Failed to marshal event data")
		eventDataJSON = []byte("{}")
	}

	// Create execution event record
	executionEvent := &models.ExecutionEvent{
		ExecutionID: event.ExecutionID,
		Type:        event.Type,
		StepID:      event.StepID,
		Timestamp:   event.Timestamp,
		Data:        string(eventDataJSON),
		Error:       event.Error,
		CreatedAt:   time.Now().UTC(),
	}

	// Save to database
	if err := h.db.Create(executionEvent).Error; err != nil {
		h.logger.WithFields(logrus.Fields{
			"execution_id": event.ExecutionID,
			"event_type":   event.Type,
			"error":        err.Error(),
		}).Error("Failed to save execution event")
		return err
	}

	// Update execution status based on event type
	if err := h.updateExecutionStatus(event); err != nil {
		h.logger.WithError(err).Error("Failed to update execution status")
		return err
	}

	return nil
}

func (h *DatabaseEventHandler) updateExecutionStatus(event *WorkflowEvent) error {
	var updates map[string]interface{}

	switch event.Type {
	case "execution.started":
		updates = map[string]interface{}{
			"status":     models.ExecutionStatusRunning,
			"started_at": event.Timestamp,
			"updated_at": time.Now().UTC(),
		}

	case "execution.completed":
		updates = map[string]interface{}{
			"status":       models.ExecutionStatusCompleted,
			"completed_at": event.Timestamp,
			"updated_at":   time.Now().UTC(),
		}
		if duration, ok := event.Data["duration"].(int64); ok {
			updates["duration"] = duration
		}
		if output, ok := event.Data["output"]; ok {
			updates["output"] = output
		}

	case "execution.failed":
		updates = map[string]interface{}{
			"status":       models.ExecutionStatusFailed,
			"error":        event.Error,
			"completed_at": event.Timestamp,
			"updated_at":   time.Now().UTC(),
		}
		if duration, ok := event.Data["duration"].(int64); ok {
			updates["duration"] = duration
		}

	case "execution.cancelled":
		updates = map[string]interface{}{
			"status":       models.ExecutionStatusCancelled,
			"error":        event.Error,
			"completed_at": event.Timestamp,
			"updated_at":   time.Now().UTC(),
		}
		if duration, ok := event.Data["duration"].(int64); ok {
			updates["duration"] = duration
		}

	default:
		// For step events, just update the timestamp
		updates = map[string]interface{}{
			"updated_at": time.Now().UTC(),
		}
	}

	if len(updates) > 0 {
		return h.db.Model(&models.Execution{}).Where("id = ?", event.ExecutionID).Updates(updates).Error
	}

	return nil
}

func (h *DatabaseEventHandler) GetEventTypes() []string {
	return []string{
		"execution.started",
		"execution.completed",
		"execution.failed",
		"execution.cancelled",
		"step.started",
		"step.completed",
		"step.failed",
	}
}

// MetricsEventHandler handles workflow events by recording metrics
type MetricsEventHandler struct {
	metrics MetricsCollector
	logger  *logrus.Logger
}

// NewMetricsEventHandler creates a new metrics event handler
func NewMetricsEventHandler(metrics MetricsCollector, logger *logrus.Logger) *MetricsEventHandler {
	return &MetricsEventHandler{
		metrics: metrics,
		logger:  logger,
	}
}

func (h *MetricsEventHandler) Handle(event *WorkflowEvent) error {
	labels := map[string]string{
		"workflow_id": event.WorkflowID.String(),
		"event_type":  event.Type,
	}

	if event.StepID != "" {
		labels["step_id"] = event.StepID
	}

	switch event.Type {
	case "execution.started":
		h.metrics.RecordMetric("workflow_executions_started_total", 1, labels)

	case "execution.completed":
		h.metrics.RecordMetric("workflow_executions_completed_total", 1, labels)
		if duration, ok := event.Data["duration"].(int64); ok {
			h.metrics.RecordMetric("workflow_execution_duration_seconds", float64(duration), labels)
		}

	case "execution.failed":
		h.metrics.RecordMetric("workflow_executions_failed_total", 1, labels)
		if duration, ok := event.Data["duration"].(int64); ok {
			h.metrics.RecordMetric("workflow_execution_duration_seconds", float64(duration), labels)
		}

	case "execution.cancelled":
		h.metrics.RecordMetric("workflow_executions_cancelled_total", 1, labels)

	case "step.started":
		h.metrics.RecordMetric("workflow_steps_started_total", 1, labels)

	case "step.completed":
		h.metrics.RecordMetric("workflow_steps_completed_total", 1, labels)
		if duration, ok := event.Data["duration"].(float64); ok {
			h.metrics.RecordMetric("workflow_step_duration_seconds", duration, labels)
		}

	case "step.failed":
		h.metrics.RecordMetric("workflow_steps_failed_total", 1, labels)
		if duration, ok := event.Data["duration"].(float64); ok {
			h.metrics.RecordMetric("workflow_step_duration_seconds", duration, labels)
		}
	}

	return nil
}

func (h *MetricsEventHandler) GetEventTypes() []string {
	return []string{
		"execution.started",
		"execution.completed",
		"execution.failed",
		"execution.cancelled",
		"step.started",
		"step.completed",
		"step.failed",
	}
}

// WebhookEventHandler handles workflow events by sending webhooks
type WebhookEventHandler struct {
	webhooks []models.Webhook
	client   *http.Client
	logger   *logrus.Logger
}

// NewWebhookEventHandler creates a new webhook event handler
func NewWebhookEventHandler(webhooks []models.Webhook, logger *logrus.Logger) *WebhookEventHandler {
	return &WebhookEventHandler{
		webhooks: webhooks,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (h *WebhookEventHandler) Handle(event *WorkflowEvent) error {
	for _, webhook := range h.webhooks {
		// Check if webhook is interested in this event type
		if !h.shouldSendWebhook(webhook, event) {
			continue
		}

		go h.sendWebhook(webhook, event)
	}

	return nil
}

func (h *WebhookEventHandler) shouldSendWebhook(webhook models.Webhook, event *WorkflowEvent) bool {
	// Check if webhook is enabled
	if !webhook.IsEnabled {
		return false
	}

	// Check event type filters
	if len(webhook.Events) > 0 {
		found := false
		for _, eventType := range webhook.Events {
			if eventType == event.Type {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check workflow filters
	if len(webhook.WorkflowIDs) > 0 {
		found := false
		for _, workflowID := range webhook.WorkflowIDs {
			if workflowID == event.WorkflowID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (h *WebhookEventHandler) sendWebhook(webhook models.Webhook, event *WorkflowEvent) {
	// Prepare webhook payload
	payload := map[string]interface{}{
		"event":        event,
		"webhook_id":   webhook.ID,
		"timestamp":    time.Now().UTC(),
	}

	// Add custom headers
	headers := make(map[string]string)
	for key, value := range webhook.Headers {
		headers[key] = value
	}
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = "Magic-Flow-Webhook/1.0"

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"webhook_id": webhook.ID,
			"error":      err.Error(),
		}).Error("Failed to marshal webhook payload")
		return
	}

	// Create request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"webhook_id": webhook.ID,
			"url":        webhook.URL,
			"error":      err.Error(),
		}).Error("Failed to create webhook request")
		return
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add signature if secret is configured
	if webhook.Secret != "" {
		signature := h.generateSignature(payloadBytes, webhook.Secret)
		req.Header.Set("X-Magic-Flow-Signature", signature)
	}

	// Send request with retries
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := h.client.Do(req)
		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"webhook_id": webhook.ID,
				"url":        webhook.URL,
				"attempt":    attempt + 1,
				"error":      err.Error(),
			}).Warn("Webhook request failed")

			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			h.logger.WithFields(logrus.Fields{
				"webhook_id":  webhook.ID,
				"url":         webhook.URL,
				"status_code": resp.StatusCode,
				"attempt":     attempt + 1,
			}).Info("Webhook sent successfully")
			return
		}

		h.logger.WithFields(logrus.Fields{
			"webhook_id":  webhook.ID,
			"url":         webhook.URL,
			"status_code": resp.StatusCode,
			"attempt":     attempt + 1,
		}).Warn("Webhook request returned error status")

		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
}

func (h *WebhookEventHandler) generateSignature(payload []byte, secret string) string {
	// Simple HMAC-SHA256 signature
	// In a real implementation, you'd use crypto/hmac
	return fmt.Sprintf("sha256=%x", payload) // Simplified for demo
}

func (h *WebhookEventHandler) GetEventTypes() []string {
	return []string{
		"execution.started",
		"execution.completed",
		"execution.failed",
		"execution.cancelled",
		"step.started",
		"step.completed",
		"step.failed",
	}
}

// LogEventHandler handles workflow events by logging them
type LogEventHandler struct {
	logger *logrus.Logger
}

// NewLogEventHandler creates a new log event handler
func NewLogEventHandler(logger *logrus.Logger) *LogEventHandler {
	return &LogEventHandler{
		logger: logger,
	}
}

func (h *LogEventHandler) Handle(event *WorkflowEvent) error {
	fields := logrus.Fields{
		"event_type":   event.Type,
		"execution_id": event.ExecutionID,
		"workflow_id":  event.WorkflowID,
		"timestamp":    event.Timestamp,
	}

	if event.StepID != "" {
		fields["step_id"] = event.StepID
	}

	if event.Error != "" {
		fields["error"] = event.Error
	}

	if len(event.Data) > 0 {
		fields["data"] = event.Data
	}

	switch event.Type {
	case "execution.started":
		h.logger.WithFields(fields).Info("Workflow execution started")
	case "execution.completed":
		h.logger.WithFields(fields).Info("Workflow execution completed")
	case "execution.failed":
		h.logger.WithFields(fields).Error("Workflow execution failed")
	case "execution.cancelled":
		h.logger.WithFields(fields).Warn("Workflow execution cancelled")
	case "step.started":
		h.logger.WithFields(fields).Debug("Workflow step started")
	case "step.completed":
		h.logger.WithFields(fields).Debug("Workflow step completed")
	case "step.failed":
		h.logger.WithFields(fields).Error("Workflow step failed")
	default:
		h.logger.WithFields(fields).Info("Workflow event")
	}

	return nil
}

func (h *LogEventHandler) GetEventTypes() []string {
	return []string{
		"execution.started",
		"execution.completed",
		"execution.failed",
		"execution.cancelled",
		"step.started",
		"step.completed",
		"step.failed",
	}
}