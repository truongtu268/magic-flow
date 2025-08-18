package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magic-flow/v2/internal/engine"
	"github.com/magic-flow/v2/internal/metrics"
	"github.com/magic-flow/v2/internal/services"
	"github.com/magic-flow/v2/pkg/models"
	"github.com/sirupsen/logrus"
)

// Handler represents the API handler
type Handler struct {
	services        *services.Container
	workflowEngine  *engine.Engine
	metricsCollector *metrics.Collector
}

// NewHandler creates a new API handler
func NewHandler(services *services.Container, workflowEngine *engine.Engine, metricsCollector *metrics.Collector) *Handler {
	return &Handler{
		services:        services,
		workflowEngine:  workflowEngine,
		metricsCollector: metricsCollector,
	}
}

// SetupRoutes sets up all API routes
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", h.healthCheck)
	router.GET("/ready", h.readinessCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Workflow management
		workflows := v1.Group("/workflows")
		{
			workflows.POST("", h.createWorkflow)
			workflows.GET("", h.listWorkflows)
			workflows.GET("/:id", h.getWorkflow)
			workflows.PUT("/:id", h.updateWorkflow)
			workflows.DELETE("/:id", h.deleteWorkflow)
			workflows.POST("/:id/validate", h.validateWorkflow)
		}

		// Workflow execution
		executions := v1.Group("/executions")
		{
			executions.POST("/workflows/:id/execute", h.executeWorkflow)
			executions.GET("/:id", h.getExecution)
			executions.GET("/:id/status", h.getExecutionStatus)
			executions.GET("/:id/results", h.getExecutionResults)
			executions.GET("/:id/events", h.streamExecutionEvents)
			executions.GET("", h.listExecutions)
			executions.POST("/:id/cancel", h.cancelExecution)
			executions.POST("/:id/retry", h.retryExecution)
			executions.GET("/:id/logs", h.getExecutionLogs)
		}

		// Metrics
		metrics := v1.Group("/metrics")
		{
			metrics.GET("/workflows", h.getWorkflowMetrics)
			metrics.GET("/workflows/:id", h.getWorkflowMetricsById)
			metrics.GET("/system", h.getSystemMetrics)
			metrics.POST("/custom", h.recordCustomMetric)
			metrics.GET("/custom", h.getCustomMetrics)
			metrics.GET("/aggregations", h.getMetricAggregations)
		}

		// Code generation
		codegen := v1.Group("/codegen")
		{
			codegen.POST("/generate", h.generateCode)
			codegen.GET("/jobs/:id", h.getCodeGenStatus)
			codegen.GET("/jobs/:id/download", h.downloadGeneratedCode)
			codegen.GET("/templates", h.listCodeGenTemplates)
			codegen.GET("/jobs", h.listCodeGenJobs)
		}

		// Version management
		versions := v1.Group("/versions")
		{
			versions.POST("/workflows/:id/versions", h.createWorkflowVersion)
			versions.GET("/workflows/:id/versions", h.listWorkflowVersions)
			versions.GET("/workflows/:id/versions/:version", h.getWorkflowVersion)
			versions.POST("/workflows/:id/versions/:version/rollback", h.rollbackWorkflowVersion)
			versions.GET("/workflows/:id/versions/:from/compare/:to", h.compareWorkflowVersions)
			versions.POST("/workflows/:id/versions/:version/deploy", h.deployWorkflowVersion)
		}

		// Dashboard
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/overview", h.getDashboardOverview)
			dashboard.GET("/workflows/status", h.getWorkflowStatusSummary)
			dashboard.GET("/system/health", h.getSystemHealth)
			dashboard.GET("/metrics/live", h.getLiveMetrics)
			dashboard.POST("/dashboards", h.createDashboard)
			dashboard.GET("/dashboards", h.listDashboards)
			dashboard.GET("/dashboards/:id", h.getDashboard)
			dashboard.PUT("/dashboards/:id", h.updateDashboard)
			dashboard.DELETE("/dashboards/:id", h.deleteDashboard)
		}

		// Alerts
		alerts := v1.Group("/alerts")
		{
			alerts.POST("", h.createAlert)
			alerts.GET("", h.listAlerts)
			alerts.GET("/:id", h.getAlert)
			alerts.PUT("/:id", h.updateAlert)
			alerts.DELETE("/:id", h.deleteAlert)
			alerts.POST("/:id/enable", h.enableAlert)
			alerts.POST("/:id/disable", h.disableAlert)
			alerts.GET("/:id/events", h.getAlertEvents)
		}
	}

	// WebSocket endpoints
	ws := router.Group("/ws")
	{
		ws.GET("/executions/:id", h.streamExecutionWebSocket)
		ws.GET("/metrics", h.streamMetricsWebSocket)
		ws.GET("/alerts", h.streamAlertsWebSocket)
	}

	// Static files for dashboard
	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/dashboard", "./web/index.html")
}

// Health check endpoint
func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "2.0.0",
	})
}

// Readiness check endpoint
func (h *Handler) readinessCheck(c *gin.Context) {
	// Check database connection
	if err := h.services.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  err.Error(),
		})
		return
	}

	// Check workflow engine
	if !h.workflowEngine.IsReady() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "workflow engine not ready",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
	})
}

// Error response helper
func (h *Handler) errorResponse(c *gin.Context, statusCode int, message string, err error) {
	logrus.WithError(err).Error(message)
	c.JSON(statusCode, gin.H{
		"error":     message,
		"timestamp": time.Now().UTC(),
	})
}

// Success response helper
func (h *Handler) successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"timestamp": time.Now().UTC(),
	})
}

// Parse UUID from path parameter
func (h *Handler) parseUUID(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid UUID format", err)
		return uuid.Nil, err
	}
	return id, nil
}

// Parse pagination parameters
func (h *Handler) parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	return page, limit
}

// Parse time range parameters
func (h *Handler) parseTimeRange(c *gin.Context) (time.Time, time.Time, error) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	
	var start, end time.Time
	var err error
	
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		start = time.Now().Add(-24 * time.Hour) // Default to last 24 hours
	}
	
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		end = time.Now()
	}
	
	return start, end, nil
}

// Validate request body
func (h *Handler) validateRequestBody(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return err
	}
	return nil
}

// Get user ID from context (for authentication)
func (h *Handler) getUserID(c *gin.Context) string {
	// This would typically come from JWT token or session
	// For now, return a default user ID
	return c.GetHeader("X-User-ID")
}

// Response structures
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Timestamp  time.Time   `json:"timestamp"`
}

type ExecutionRequest struct {
	Input       map[string]interface{} `json:"input"`
	Environment string                 `json:"environment,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Priority    string                 `json:"priority,omitempty"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

type CodeGenRequest struct {
	WorkflowID uuid.UUID `json:"workflow_id"`
	Language   string    `json:"language"`
	Template   string    `json:"template,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

type MetricRequest struct {
	Name        string                 `json:"name"`
	Type        models.MetricType      `json:"type"`
	Value       float64                `json:"value"`
	Unit        string                 `json:"unit,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	WorkflowID  *uuid.UUID             `json:"workflow_id,omitempty"`
	ExecutionID *uuid.UUID             `json:"execution_id,omitempty"`
}

type DashboardOverview struct {
	TotalWorkflows      int64                  `json:"total_workflows"`
	ActiveExecutions    int64                  `json:"active_executions"`
	CompletedToday      int64                  `json:"completed_today"`
	FailedToday         int64                  `json:"failed_today"`
	AverageExecutionTime float64               `json:"average_execution_time"`
	SystemHealth        string                 `json:"system_health"`
	RecentExecutions    []models.Execution     `json:"recent_executions"`
	TopWorkflows        []WorkflowStats        `json:"top_workflows"`
	Alerts              []models.Alert         `json:"active_alerts"`
}

type WorkflowStats struct {
	Workflow       models.Workflow `json:"workflow"`
	ExecutionCount int64           `json:"execution_count"`
	SuccessRate    float64         `json:"success_rate"`
	AverageTime    float64         `json:"average_time"`
}

type SystemHealth struct {
	Status      string                 `json:"status"`
	Components  map[string]string      `json:"components"`
	Metrics     map[string]interface{} `json:"metrics"`
	Timestamp   time.Time              `json:"timestamp"`
}