package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// AlertService handles alert business logic
type AlertService struct {
	repos  *database.RepositoryManager
	logger *logrus.Logger
}

// NewAlertService creates a new alert service
func NewAlertService(repos *database.RepositoryManager, logger *logrus.Logger) *AlertService {
	return &AlertService{
		repos:  repos,
		logger: logger,
	}
}

// CreateAlert creates a new alert
func (s *AlertService) CreateAlert(req *CreateAlertRequest) (*models.Alert, error) {
	alert := &models.Alert{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Type:        models.AlertType(req.Type),
		Conditions:  req.Conditions,
		Actions:     req.Actions,
		Severity:    models.AlertSeverity(req.Severity),
		Enabled:     req.Enabled,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Validate alert conditions
	if err := s.validateAlertConditions(alert.Conditions); err != nil {
		return nil, fmt.Errorf("invalid alert conditions: %w", err)
	}

	// Validate alert actions
	if err := s.validateAlertActions(alert.Actions); err != nil {
		return nil, fmt.Errorf("invalid alert actions: %w", err)
	}

	if err := s.repos.Alert.Create(alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id":   alert.ID,
		"alert_name": alert.Name,
		"type":       alert.Type,
		"severity":   alert.Severity,
		"created_by": alert.CreatedBy,
	}).Info("Alert created")

	return alert, nil
}

// GetAlert retrieves an alert by ID
func (s *AlertService) GetAlert(id uuid.UUID) (*models.Alert, error) {
	alert, err := s.repos.Alert.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("alert not found")
		}
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	return alert, nil
}

// ListAlerts retrieves alerts with pagination and filtering
func (s *AlertService) ListAlerts(req *ListAlertsRequest) ([]*models.Alert, int64, error) {
	alerts, total, err := s.repos.Alert.List(req.Limit, req.Offset, req.Type, req.Severity, req.Enabled)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list alerts: %w", err)
	}
	return alerts, total, nil
}

// UpdateAlert updates an existing alert
func (s *AlertService) UpdateAlert(id uuid.UUID, req *UpdateAlertRequest) (*models.Alert, error) {
	alert, err := s.repos.Alert.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("alert not found")
		}
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	// Update fields
	if req.Name != "" {
		alert.Name = req.Name
	}
	if req.Description != "" {
		alert.Description = req.Description
	}
	if req.Type != "" {
		alert.Type = models.AlertType(req.Type)
	}
	if req.Conditions != nil {
		if err := s.validateAlertConditions(req.Conditions); err != nil {
			return nil, fmt.Errorf("invalid alert conditions: %w", err)
		}
		alert.Conditions = req.Conditions
	}
	if req.Actions != nil {
		if err := s.validateAlertActions(req.Actions); err != nil {
			return nil, fmt.Errorf("invalid alert actions: %w", err)
		}
		alert.Actions = req.Actions
	}
	if req.Severity != "" {
		alert.Severity = models.AlertSeverity(req.Severity)
	}
	if req.Enabled != nil {
		alert.Enabled = *req.Enabled
	}

	alert.UpdatedBy = req.UpdatedBy
	alert.UpdatedAt = time.Now().UTC()

	if err := s.repos.Alert.Update(alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id":   alert.ID,
		"alert_name": alert.Name,
		"updated_by": alert.UpdatedBy,
	}).Info("Alert updated")

	return alert, nil
}

// DeleteAlert deletes an alert
func (s *AlertService) DeleteAlert(id uuid.UUID) error {
	// Check if alert exists
	_, err := s.repos.Alert.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("alert not found")
		}
		return fmt.Errorf("failed to get alert: %w", err)
	}

	if err := s.repos.Alert.Delete(id); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id": id,
	}).Info("Alert deleted")

	return nil
}

// EnableAlert enables an alert
func (s *AlertService) EnableAlert(id uuid.UUID, enabledBy string) error {
	alert, err := s.repos.Alert.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("alert not found")
		}
		return fmt.Errorf("failed to get alert: %w", err)
	}

	if alert.Enabled {
		return fmt.Errorf("alert is already enabled")
	}

	alert.Enabled = true
	alert.UpdatedBy = enabledBy
	alert.UpdatedAt = time.Now().UTC()

	if err := s.repos.Alert.Update(alert); err != nil {
		return fmt.Errorf("failed to enable alert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id":   alert.ID,
		"alert_name": alert.Name,
		"enabled_by": enabledBy,
	}).Info("Alert enabled")

	return nil
}

// DisableAlert disables an alert
func (s *AlertService) DisableAlert(id uuid.UUID, disabledBy string) error {
	alert, err := s.repos.Alert.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("alert not found")
		}
		return fmt.Errorf("failed to get alert: %w", err)
	}

	if !alert.Enabled {
		return fmt.Errorf("alert is already disabled")
	}

	alert.Enabled = false
	alert.UpdatedBy = disabledBy
	alert.UpdatedAt = time.Now().UTC()

	if err := s.repos.Alert.Update(alert); err != nil {
		return fmt.Errorf("failed to disable alert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id":    alert.ID,
		"alert_name":  alert.Name,
		"disabled_by": disabledBy,
	}).Info("Alert disabled")

	return nil
}

// GetAlertEvents retrieves events for an alert
func (s *AlertService) GetAlertEvents(alertID uuid.UUID, req *GetAlertEventsRequest) ([]*models.AlertEvent, int64, error) {
	// Check if alert exists
	_, err := s.repos.Alert.GetByID(alertID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("alert not found")
		}
		return nil, 0, fmt.Errorf("failed to get alert: %w", err)
	}

	events, total, err := s.repos.Alert.GetAlertEvents(alertID, req.Limit, req.Offset, req.StartTime, req.EndTime)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get alert events: %w", err)
	}

	return events, total, nil
}

// TriggerAlert triggers an alert and creates an alert event
func (s *AlertService) TriggerAlert(alertID uuid.UUID, triggerData map[string]interface{}) error {
	alert, err := s.repos.Alert.GetByID(alertID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("alert not found")
		}
		return fmt.Errorf("failed to get alert: %w", err)
	}

	if !alert.Enabled {
		s.logger.WithField("alert_id", alertID).Debug("Alert is disabled, skipping trigger")
		return nil
	}

	// Create alert event
	event := &models.AlertEvent{
		ID:        uuid.New(),
		AlertID:   alertID,
		Message:   s.generateAlertMessage(alert, triggerData),
		Severity:  alert.Severity,
		Data:      triggerData,
		Timestamp: time.Now().UTC(),
	}

	if err := s.repos.Alert.CreateAlertEvent(event); err != nil {
		return fmt.Errorf("failed to create alert event: %w", err)
	}

	// Execute alert actions
	if err := s.executeAlertActions(alert, event); err != nil {
		s.logger.WithError(err).WithField("alert_id", alertID).Error("Failed to execute alert actions")
	}

	s.logger.WithFields(logrus.Fields{
		"alert_id":   alertID,
		"alert_name": alert.Name,
		"severity":   alert.Severity,
		"event_id":   event.ID,
	}).Info("Alert triggered")

	return nil
}

// EvaluateAlerts evaluates all enabled alerts against current metrics
func (s *AlertService) EvaluateAlerts() error {
	// Get all enabled alerts
	alerts, _, err := s.repos.Alert.List(1000, 0, "", "", &[]bool{true}[0])
	if err != nil {
		return fmt.Errorf("failed to get enabled alerts: %w", err)
	}

	for _, alert := range alerts {
		if err := s.evaluateAlert(alert); err != nil {
			s.logger.WithError(err).WithField("alert_id", alert.ID).Error("Failed to evaluate alert")
		}
	}

	return nil
}

// Helper methods
func (s *AlertService) validateAlertConditions(conditions map[string]interface{}) error {
	// Basic validation - ensure required fields exist
	if conditions == nil {
		return fmt.Errorf("conditions cannot be nil")
	}

	// Check for required condition fields
	if _, ok := conditions["metric"]; !ok {
		return fmt.Errorf("metric field is required in conditions")
	}

	if _, ok := conditions["operator"]; !ok {
		return fmt.Errorf("operator field is required in conditions")
	}

	if _, ok := conditions["threshold"]; !ok {
		return fmt.Errorf("threshold field is required in conditions")
	}

	// Validate operator
	operator, ok := conditions["operator"].(string)
	if !ok {
		return fmt.Errorf("operator must be a string")
	}

	validOperators := []string{">", ">=", "<", "<=", "==", "!="}
	validOperator := false
	for _, validOp := range validOperators {
		if operator == validOp {
			validOperator = true
			break
		}
	}

	if !validOperator {
		return fmt.Errorf("invalid operator: %s", operator)
	}

	return nil
}

func (s *AlertService) validateAlertActions(actions map[string]interface{}) error {
	if actions == nil {
		return fmt.Errorf("actions cannot be nil")
	}

	// Validate action types
	for actionType, actionConfig := range actions {
		switch actionType {
		case "email":
			if err := s.validateEmailAction(actionConfig); err != nil {
				return fmt.Errorf("invalid email action: %w", err)
			}
		case "webhook":
			if err := s.validateWebhookAction(actionConfig); err != nil {
				return fmt.Errorf("invalid webhook action: %w", err)
			}
		case "slack":
			if err := s.validateSlackAction(actionConfig); err != nil {
				return fmt.Errorf("invalid slack action: %w", err)
			}
		default:
			return fmt.Errorf("unsupported action type: %s", actionType)
		}
	}

	return nil
}

func (s *AlertService) validateEmailAction(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("email action config must be an object")
	}

	if _, ok := configMap["recipients"]; !ok {
		return fmt.Errorf("recipients field is required for email action")
	}

	return nil
}

func (s *AlertService) validateWebhookAction(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("webhook action config must be an object")
	}

	if _, ok := configMap["url"]; !ok {
		return fmt.Errorf("url field is required for webhook action")
	}

	return nil
}

func (s *AlertService) validateSlackAction(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("slack action config must be an object")
	}

	if _, ok := configMap["channel"]; !ok {
		return fmt.Errorf("channel field is required for slack action")
	}

	return nil
}

func (s *AlertService) generateAlertMessage(alert *models.Alert, triggerData map[string]interface{}) string {
	return fmt.Sprintf("Alert '%s' triggered: %s", alert.Name, alert.Description)
}

func (s *AlertService) executeAlertActions(alert *models.Alert, event *models.AlertEvent) error {
	// This is a simplified implementation
	// In a real system, you would implement actual email, webhook, and Slack integrations
	for actionType, actionConfig := range alert.Actions {
		switch actionType {
		case "email":
			s.logger.WithFields(logrus.Fields{
				"alert_id":  alert.ID,
				"event_id":  event.ID,
				"action":    "email",
				"config":    actionConfig,
			}).Info("Would send email notification")
		case "webhook":
			s.logger.WithFields(logrus.Fields{
				"alert_id":  alert.ID,
				"event_id":  event.ID,
				"action":    "webhook",
				"config":    actionConfig,
			}).Info("Would send webhook notification")
		case "slack":
			s.logger.WithFields(logrus.Fields{
				"alert_id":  alert.ID,
				"event_id":  event.ID,
				"action":    "slack",
				"config":    actionConfig,
			}).Info("Would send Slack notification")
		}
	}

	return nil
}

func (s *AlertService) evaluateAlert(alert *models.Alert) error {
	// This is a simplified implementation
	// In a real system, you would fetch current metrics and evaluate conditions
	s.logger.WithFields(logrus.Fields{
		"alert_id":   alert.ID,
		"alert_name": alert.Name,
	}).Debug("Evaluating alert conditions")

	return nil
}

// Request/Response types
type CreateAlertRequest struct {
	Name        string                 `json:"name" validate:"required,max=255"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type" validate:"required"`
	Conditions  map[string]interface{} `json:"conditions" validate:"required"`
	Actions     map[string]interface{} `json:"actions" validate:"required"`
	Severity    string                 `json:"severity" validate:"required"`
	Enabled     bool                   `json:"enabled"`
	CreatedBy   string                 `json:"created_by,omitempty"`
}

type UpdateAlertRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Actions     map[string]interface{} `json:"actions,omitempty"`
	Severity    string                 `json:"severity,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	UpdatedBy   string                 `json:"updated_by,omitempty"`
}

type ListAlertsRequest struct {
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Type     string  `json:"type,omitempty"`
	Severity string  `json:"severity,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

type GetAlertEventsRequest struct {
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}