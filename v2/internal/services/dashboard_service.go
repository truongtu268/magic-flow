package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	repos  *database.RepositoryManager
	logger *logrus.Logger
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(repos *database.RepositoryManager, logger *logrus.Logger) *DashboardService {
	return &DashboardService{
		repos:  repos,
		logger: logger,
	}
}

// CreateDashboard creates a new dashboard
func (s *DashboardService) CreateDashboard(req *CreateDashboardRequest) (*models.Dashboard, error) {
	// Validate dashboard configuration
	if err := s.validateDashboardConfig(req.Config); err != nil {
		return nil, fmt.Errorf("invalid dashboard configuration: %w", err)
	}

	dashboard := &models.Dashboard{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		IsPublic:    req.IsPublic,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repos.Dashboard.Create(dashboard); err != nil {
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"dashboard_id":   dashboard.ID,
		"dashboard_name": dashboard.Name,
		"created_by":     dashboard.CreatedBy,
		"is_public":      dashboard.IsPublic,
	}).Info("Dashboard created")

	return dashboard, nil
}

// GetDashboard retrieves a dashboard by ID
func (s *DashboardService) GetDashboard(id uuid.UUID) (*models.Dashboard, error) {
	dashboard, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dashboard not found")
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	return dashboard, nil
}

// ListDashboards retrieves dashboards with pagination and filtering
func (s *DashboardService) ListDashboards(req *ListDashboardsRequest) ([]*models.Dashboard, int64, error) {
	dashboards, total, err := s.repos.Dashboard.List(req.Limit, req.Offset, req.IsPublic, req.CreatedBy)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list dashboards: %w", err)
	}
	return dashboards, total, nil
}

// UpdateDashboard updates an existing dashboard
func (s *DashboardService) UpdateDashboard(id uuid.UUID, req *UpdateDashboardRequest) (*models.Dashboard, error) {
	dashboard, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dashboard not found")
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	// Update fields
	if req.Name != "" {
		dashboard.Name = req.Name
	}
	if req.Description != "" {
		dashboard.Description = req.Description
	}
	if req.Config != nil {
		if err := s.validateDashboardConfig(req.Config); err != nil {
			return nil, fmt.Errorf("invalid dashboard configuration: %w", err)
		}
		dashboard.Config = req.Config
	}
	if req.IsPublic != nil {
		dashboard.IsPublic = *req.IsPublic
	}

	dashboard.UpdatedBy = req.UpdatedBy
	dashboard.UpdatedAt = time.Now().UTC()

	if err := s.repos.Dashboard.Update(dashboard); err != nil {
		return nil, fmt.Errorf("failed to update dashboard: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"dashboard_id":   dashboard.ID,
		"dashboard_name": dashboard.Name,
		"updated_by":     dashboard.UpdatedBy,
	}).Info("Dashboard updated")

	return dashboard, nil
}

// DeleteDashboard deletes a dashboard
func (s *DashboardService) DeleteDashboard(id uuid.UUID) error {
	// Check if dashboard exists
	_, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("dashboard not found")
		}
		return fmt.Errorf("failed to get dashboard: %w", err)
	}

	if err := s.repos.Dashboard.Delete(id); err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"dashboard_id": id,
	}).Info("Dashboard deleted")

	return nil
}

// ShareDashboard creates a shareable link for a dashboard
func (s *DashboardService) ShareDashboard(id uuid.UUID, req *ShareDashboardRequest) (*ShareDashboardResponse, error) {
	dashboard, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dashboard not found")
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	// Generate share token
	shareToken := uuid.New().String()

	// Create share configuration
	shareConfig := map[string]interface{}{
		"dashboard_id": id,
		"share_token":  shareToken,
		"expires_at":   req.ExpiresAt,
		"permissions":  req.Permissions,
		"shared_by":    req.SharedBy,
		"created_at":   time.Now().UTC(),
	}

	// Store share configuration (in a real implementation, you'd store this in a separate table)
	s.logger.WithFields(logrus.Fields{
		"dashboard_id": id,
		"share_token":  shareToken,
		"shared_by":    req.SharedBy,
	}).Info("Dashboard shared")

	return &ShareDashboardResponse{
		ShareToken: shareToken,
		ShareURL:   fmt.Sprintf("/shared/dashboard/%s", shareToken),
		ExpiresAt:  req.ExpiresAt,
	}, nil
}

// ExportDashboard exports a dashboard configuration
func (s *DashboardService) ExportDashboard(id uuid.UUID) (*ExportDashboardResponse, error) {
	dashboard, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dashboard not found")
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	// Create export data
	exportData := DashboardExport{
		Version:     "1.0",
		Name:        dashboard.Name,
		Description: dashboard.Description,
		Config:      dashboard.Config,
		ExportedAt:  time.Now().UTC(),
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dashboard export: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"dashboard_id":   dashboard.ID,
		"dashboard_name": dashboard.Name,
	}).Info("Dashboard exported")

	return &ExportDashboardResponse{
		Filename: fmt.Sprintf("%s-dashboard-export.json", dashboard.Name),
		Content:  string(jsonData),
		Size:     len(jsonData),
	}, nil
}

// ImportDashboard imports a dashboard from exported data
func (s *DashboardService) ImportDashboard(req *ImportDashboardRequest) (*models.Dashboard, error) {
	// Parse import data
	var importData DashboardExport
	if err := json.Unmarshal([]byte(req.Content), &importData); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	// Validate configuration
	if err := s.validateDashboardConfig(importData.Config); err != nil {
		return nil, fmt.Errorf("invalid dashboard configuration in import: %w", err)
	}

	// Create dashboard
	dashboard := &models.Dashboard{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: importData.Description,
		Config:      importData.Config,
		IsPublic:    false, // Imported dashboards are private by default
		CreatedBy:   req.ImportedBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repos.Dashboard.Create(dashboard); err != nil {
		return nil, fmt.Errorf("failed to create imported dashboard: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"dashboard_id":   dashboard.ID,
		"dashboard_name": dashboard.Name,
		"imported_by":    req.ImportedBy,
	}).Info("Dashboard imported")

	return dashboard, nil
}

// GetDashboardData retrieves data for dashboard widgets
func (s *DashboardService) GetDashboardData(id uuid.UUID, req *GetDashboardDataRequest) (*DashboardDataResponse, error) {
	dashboard, err := s.repos.Dashboard.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dashboard not found")
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	// Extract widgets from dashboard config
	widgets, ok := dashboard.Config["widgets"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid dashboard configuration: widgets not found")
	}

	// Collect data for each widget
	widgetData := make(map[string]interface{})
	for _, widget := range widgets {
		widgetMap, ok := widget.(map[string]interface{})
		if !ok {
			continue
		}

		widgetID, ok := widgetMap["id"].(string)
		if !ok {
			continue
		}

		data, err := s.getWidgetData(widgetMap, req.TimeRange)
		if err != nil {
			s.logger.WithError(err).WithField("widget_id", widgetID).Warn("Failed to get widget data")
			widgetData[widgetID] = map[string]interface{}{"error": err.Error()}
		} else {
			widgetData[widgetID] = data
		}
	}

	return &DashboardDataResponse{
		DashboardID: id,
		WidgetData:  widgetData,
		TimeRange:   req.TimeRange,
		LastUpdated: time.Now().UTC(),
	}, nil
}

// Helper methods
func (s *DashboardService) validateDashboardConfig(config map[string]interface{}) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Check for required fields
	if _, ok := config["widgets"]; !ok {
		return fmt.Errorf("widgets field is required in config")
	}

	// Validate widgets
	widgets, ok := config["widgets"].([]interface{})
	if !ok {
		return fmt.Errorf("widgets must be an array")
	}

	for i, widget := range widgets {
		widgetMap, ok := widget.(map[string]interface{})
		if !ok {
			return fmt.Errorf("widget %d must be an object", i)
		}

		if _, ok := widgetMap["id"]; !ok {
			return fmt.Errorf("widget %d must have an id field", i)
		}

		if _, ok := widgetMap["type"]; !ok {
			return fmt.Errorf("widget %d must have a type field", i)
		}
	}

	return nil
}

func (s *DashboardService) getWidgetData(widget map[string]interface{}, timeRange *TimeRange) (interface{}, error) {
	widgetType, ok := widget["type"].(string)
	if !ok {
		return nil, fmt.Errorf("widget type must be a string")
	}

	switch widgetType {
	case "metric":
		return s.getMetricWidgetData(widget, timeRange)
	case "chart":
		return s.getChartWidgetData(widget, timeRange)
	case "table":
		return s.getTableWidgetData(widget, timeRange)
	case "status":
		return s.getStatusWidgetData(widget, timeRange)
	default:
		return nil, fmt.Errorf("unsupported widget type: %s", widgetType)
	}
}

func (s *DashboardService) getMetricWidgetData(widget map[string]interface{}, timeRange *TimeRange) (interface{}, error) {
	// This is a simplified implementation
	// In a real system, you would fetch actual metrics based on widget configuration
	return map[string]interface{}{
		"value":     42.5,
		"unit":      "requests/sec",
		"trend":     "up",
		"change":    "+5.2%",
		"timestamp": time.Now().UTC(),
	}, nil
}

func (s *DashboardService) getChartWidgetData(widget map[string]interface{}, timeRange *TimeRange) (interface{}, error) {
	// This is a simplified implementation
	// In a real system, you would fetch actual time series data
	dataPoints := make([]map[string]interface{}, 0)
	for i := 0; i < 24; i++ {
		dataPoints = append(dataPoints, map[string]interface{}{
			"timestamp": time.Now().UTC().Add(-time.Duration(23-i) * time.Hour),
			"value":     float64(i*2 + 10),
		})
	}

	return map[string]interface{}{
		"data":   dataPoints,
		"labels": []string{"Executions", "Success Rate"},
	}, nil
}

func (s *DashboardService) getTableWidgetData(widget map[string]interface{}, timeRange *TimeRange) (interface{}, error) {
	// This is a simplified implementation
	return map[string]interface{}{
		"headers": []string{"Workflow", "Status", "Last Run", "Success Rate"},
		"rows": [][]interface{}{
			{"Data Processing", "Active", "2 minutes ago", "98.5%"},
			{"Email Campaign", "Active", "5 minutes ago", "99.2%"},
			{"Report Generation", "Inactive", "1 hour ago", "95.8%"},
		},
	}, nil
}

func (s *DashboardService) getStatusWidgetData(widget map[string]interface{}, timeRange *TimeRange) (interface{}, error) {
	// This is a simplified implementation
	return map[string]interface{}{
		"status":  "healthy",
		"message": "All systems operational",
		"details": map[string]interface{}{
			"database":    "healthy",
			"api_server":  "healthy",
			"queue":       "healthy",
			"last_check": time.Now().UTC(),
		},
	}, nil
}

// Request/Response types
type CreateDashboardRequest struct {
	Name        string                 `json:"name" validate:"required,max=255"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config" validate:"required"`
	IsPublic    bool                   `json:"is_public"`
	CreatedBy   string                 `json:"created_by,omitempty"`
}

type UpdateDashboardRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	IsPublic    *bool                  `json:"is_public,omitempty"`
	UpdatedBy   string                 `json:"updated_by,omitempty"`
}

type ListDashboardsRequest struct {
	Limit     int     `json:"limit"`
	Offset    int     `json:"offset"`
	IsPublic  *bool   `json:"is_public,omitempty"`
	CreatedBy string  `json:"created_by,omitempty"`
}

type ShareDashboardRequest struct {
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Permissions map[string]interface{} `json:"permissions,omitempty"`
	SharedBy    string                 `json:"shared_by,omitempty"`
}

type ShareDashboardResponse struct {
	ShareToken string     `json:"share_token"`
	ShareURL   string     `json:"share_url"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type ExportDashboardResponse struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Size     int    `json:"size"`
}

type ImportDashboardRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Content    string `json:"content" validate:"required"`
	ImportedBy string `json:"imported_by,omitempty"`
}

type DashboardExport struct {
	Version     string                 `json:"version"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	ExportedAt  time.Time              `json:"exported_at"`
}

type GetDashboardDataRequest struct {
	TimeRange *TimeRange `json:"time_range,omitempty"`
}

type DashboardDataResponse struct {
	DashboardID uuid.UUID              `json:"dashboard_id"`
	WidgetData  map[string]interface{} `json:"widget_data"`
	TimeRange   *TimeRange             `json:"time_range"`
	LastUpdated time.Time              `json:"last_updated"`
}