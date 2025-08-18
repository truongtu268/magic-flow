package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"magic-flow/v2/pkg/api"
)

// Handlers provides HTTP handlers for dashboard endpoints
type Handlers struct {
	service  *Service
	upgrader websocket.Upgrader
}

// NewHandlers creates new dashboard handlers
func NewHandlers(service *Service) *Handlers {
	return &Handlers{
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// RegisterRoutes registers dashboard routes
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	// Dashboard overview
	router.GET("/overview", h.GetDashboardOverview)
	router.GET("/health", h.GetHealthStatus)

	// Metrics endpoints
	metricsGroup := router.Group("/metrics")
	{
		metricsGroup.GET("/workflows/:id", h.GetWorkflowMetrics)
		metricsGroup.GET("/executions", h.GetExecutionMetrics)
		metricsGroup.GET("/system", h.GetSystemMetrics)
	}

	// Configuration endpoints
	configGroup := router.Group("/config")
	{
		configGroup.GET("/", h.GetDashboardConfig)
		configGroup.PUT("/", h.UpdateDashboardConfig)
		configGroup.GET("/templates", h.GetDashboardTemplates)
		configGroup.POST("/templates", h.CreateDashboardTemplate)
	}

	// Widget endpoints
	widgetGroup := router.Group("/widgets")
	{
		widgetGroup.GET("/templates", h.GetWidgetTemplates)
		widgetGroup.POST("/templates", h.CreateWidgetTemplate)
		widgetGroup.GET("/data/:type", h.GetWidgetData)
	}

	// Real-time endpoints
	router.GET("/ws", h.HandleWebSocket)
	router.GET("/realtime/status", h.GetRealtimeStatus)

	// Analytics endpoints
	analyticsGroup := router.Group("/analytics")
	{
		analyticsGroup.GET("/dashboard/:id", h.GetDashboardAnalytics)
		analyticsGroup.GET("/usage", h.GetUsageAnalytics)
	}

	// Export/Import endpoints
	router.POST("/export", h.ExportDashboard)
	router.POST("/import", h.ImportDashboard)
}

// GetDashboardOverview returns the main dashboard overview
func (h *Handlers) GetDashboardOverview(c *gin.Context) {
	overview, err := h.service.GetDashboardOverview(c.Request.Context())
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dashboard overview", err)
		return
	}

	api.SuccessResponse(c, overview)
}

// GetHealthStatus returns the health status of the dashboard service
func (h *Handlers) GetHealthStatus(c *gin.Context) {
	status, err := h.service.GetHealthStatus(c.Request.Context())
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get health status", err)
		return
	}

	api.SuccessResponse(c, status)
}

// GetWorkflowMetrics returns metrics for a specific workflow
func (h *Handlers) GetWorkflowMetrics(c *gin.Context) {
	workflowIDStr := c.Param("id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid workflow ID", err)
		return
	}

	timeRange := c.DefaultQuery("time_range", "24h")

	metrics, err := h.service.GetWorkflowMetrics(c.Request.Context(), workflowID, timeRange)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get workflow metrics", err)
		return
	}

	api.SuccessResponse(c, metrics)
}

// GetExecutionMetrics returns execution metrics with filters
func (h *Handlers) GetExecutionMetrics(c *gin.Context) {
	var filters ExecutionMetricsFilters

	// Parse query parameters
	if workflowIDStr := c.Query("workflow_id"); workflowIDStr != "" {
		workflowID, err := uuid.Parse(workflowIDStr)
		if err != nil {
			api.ErrorResponse(c, http.StatusBadRequest, "Invalid workflow ID", err)
			return
		}
		filters.WorkflowID = &workflowID
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			api.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
			return
		}
		filters.UserID = &userID
	}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			api.ErrorResponse(c, http.StatusBadRequest, "Invalid start time format", err)
			return
		}
		filters.StartTime = &startTime
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			api.ErrorResponse(c, http.StatusBadRequest, "Invalid end time format", err)
			return
		}
		filters.EndTime = &endTime
	}

	filters.TimeRange = c.DefaultQuery("time_range", "24h")

	metrics, err := h.service.GetExecutionMetrics(c.Request.Context(), filters)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get execution metrics", err)
		return
	}

	api.SuccessResponse(c, metrics)
}

// GetSystemMetrics returns system-wide metrics
func (h *Handlers) GetSystemMetrics(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")

	metrics, err := h.service.GetSystemMetrics(c.Request.Context(), timeRange)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get system metrics", err)
		return
	}

	api.SuccessResponse(c, metrics)
}

// GetDashboardConfig returns dashboard configuration for the current user
func (h *Handlers) GetDashboardConfig(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	config, err := h.service.GetDashboardConfig(c.Request.Context(), userUUID)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dashboard config", err)
		return
	}

	api.SuccessResponse(c, config)
}

// UpdateDashboardConfig updates dashboard configuration
func (h *Handlers) UpdateDashboardConfig(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var config DashboardConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	config.UserID = userUUID
	config.UpdatedAt = time.Now()

	err := h.service.UpdateDashboardConfig(c.Request.Context(), &config)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to update dashboard config", err)
		return
	}

	api.SuccessResponse(c, config)
}

// GetDashboardTemplates returns available dashboard templates
func (h *Handlers) GetDashboardTemplates(c *gin.Context) {
	// Mock data for now
	templates := []DashboardTemplate{
		{
			ID:          uuid.New(),
			Name:        "Workflow Monitoring",
			Description: "Monitor workflow executions and performance",
			Category:    "monitoring",
			Tags:        []string{"workflows", "monitoring", "performance"},
			Popular:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "System Analytics",
			Description: "Analyze system performance and resource usage",
			Category:    "analytics",
			Tags:        []string{"system", "analytics", "resources"},
			Popular:     false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	api.SuccessResponse(c, templates)
}

// CreateDashboardTemplate creates a new dashboard template
func (h *Handlers) CreateDashboardTemplate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var template DashboardTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	template.ID = uuid.New()
	template.CreatedBy = userUUID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	// In a real implementation, save to database
	api.SuccessResponse(c, template)
}

// GetWidgetTemplates returns available widget templates
func (h *Handlers) GetWidgetTemplates(c *gin.Context) {
	// Mock data for now
	templates := []WidgetTemplate{
		{
			ID:          uuid.New(),
			Name:        "Execution Count",
			Description: "Display total execution count",
			Type:        WidgetTypeMetric,
			Category:    "metrics",
			Tags:        []string{"executions", "count"},
			Popular:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Success Rate Chart",
			Description: "Chart showing execution success rate over time",
			Type:        WidgetTypeChart,
			Category:    "charts",
			Tags:        []string{"success", "rate", "chart"},
			Popular:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	api.SuccessResponse(c, templates)
}

// CreateWidgetTemplate creates a new widget template
func (h *Handlers) CreateWidgetTemplate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var template WidgetTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	template.ID = uuid.New()
	template.CreatedBy = userUUID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	// In a real implementation, save to database
	api.SuccessResponse(c, template)
}

// GetWidgetData returns data for a specific widget type
func (h *Handlers) GetWidgetData(c *gin.Context) {
	widgetType := c.Param("type")
	timeRange := c.DefaultQuery("time_range", "24h")

	// Mock data based on widget type
	var data interface{}

	switch widgetType {
	case "execution_count":
		data = map[string]interface{}{
			"value": 1234,
			"change": "+5.2%",
			"trend": "up",
		}
	case "success_rate":
		data = map[string]interface{}{
			"value": 95.8,
			"change": "+2.1%",
			"trend": "up",
		}
	case "active_workflows":
		data = map[string]interface{}{
			"value": 42,
			"change": "+3",
			"trend": "up",
		}
	default:
		api.ErrorResponse(c, http.StatusBadRequest, "Unknown widget type", nil)
		return
	}

	api.SuccessResponse(c, map[string]interface{}{
		"type":       widgetType,
		"time_range": timeRange,
		"data":       data,
		"updated_at": time.Now(),
	})
}

// HandleWebSocket handles WebSocket connections for real-time updates
func (h *Handlers) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Failed to upgrade connection", err)
		return
	}
	defer conn.Close()

	// Generate client ID
	clientID := uuid.New().String()

	// Subscribe to real-time updates
	updates, err := h.service.SubscribeToRealtimeUpdates(c.Request.Context(), clientID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	defer h.service.UnsubscribeFromRealtimeUpdates(clientID)

	// Handle WebSocket messages
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Handle incoming messages if needed
		}
	}()

	// Send updates to client
	for update := range updates {
		data, err := json.Marshal(update)
		if err != nil {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}

// GetRealtimeStatus returns the status of real-time connections
func (h *Handlers) GetRealtimeStatus(c *gin.Context) {
	status := map[string]interface{}{
		"connected_clients": h.service.realtimeManager.GetConnectedClients(),
		"status":            "active",
		"timestamp":         time.Now(),
	}

	api.SuccessResponse(c, status)
}

// GetDashboardAnalytics returns analytics for a specific dashboard
func (h *Handlers) GetDashboardAnalytics(c *gin.Context) {
	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid dashboard ID", err)
		return
	}

	// Mock analytics data
	analytics := DashboardAnalytics{
		DashboardID:       dashboardID,
		ViewCount:         1250,
		UniqueViewers:     85,
		AverageViewTime:   5*time.Minute + 30*time.Second,
		BounceRate:        15.2,
		MostViewedWidget:  "execution_count",
		LeastViewedWidget: "system_health",
		PeakUsageHour:     14, // 2 PM
		LastViewed:        time.Now().Add(-2 * time.Hour),
		CreatedAt:         time.Now().AddDate(0, -1, 0),
		UpdatedAt:         time.Now(),
	}

	api.SuccessResponse(c, analytics)
}

// GetUsageAnalytics returns overall usage analytics
func (h *Handlers) GetUsageAnalytics(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "7d")

	// Mock usage analytics
	usage := map[string]interface{}{
		"time_range":        timeRange,
		"total_dashboards":  25,
		"active_dashboards": 18,
		"total_users":       12,
		"active_users":      8,
		"total_views":       5420,
		"average_session":   "8m 45s",
		"popular_widgets": []string{
			"execution_count",
			"success_rate",
			"workflow_list",
		},
		"generated_at": time.Now(),
	}

	api.SuccessResponse(c, usage)
}

// ExportDashboard exports dashboard configuration
func (h *Handlers) ExportDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var request struct {
		DashboardID uuid.UUID `json:"dashboard_id" binding:"required"`
		IncludeData bool      `json:"include_data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Get dashboard config
	config, err := h.service.GetDashboardConfig(c.Request.Context(), userUUID)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dashboard config", err)
		return
	}

	// Create export
	export := DashboardExport{
		Version:    "1.0",
		ExportedAt: time.Now(),
		ExportedBy: userUUID,
		Dashboard:  *config,
		Metadata: map[string]interface{}{
			"include_data": request.IncludeData,
		},
	}

	api.SuccessResponse(c, export)
}

// ImportDashboard imports dashboard configuration
func (h *Handlers) ImportDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		api.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var export DashboardExport
	if err := c.ShouldBindJSON(&export); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "Invalid export data", err)
		return
	}

	// Update dashboard config with imported data
	config := export.Dashboard
	config.ID = uuid.New() // Generate new ID
	config.UserID = userUUID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Generate new IDs for widgets
	for i := range config.Widgets {
		config.Widgets[i].ID = uuid.New()
		config.Widgets[i].CreatedAt = time.Now()
		config.Widgets[i].UpdatedAt = time.Now()
	}

	err := h.service.UpdateDashboardConfig(c.Request.Context(), &config)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "Failed to import dashboard", err)
		return
	}

	api.SuccessResponse(c, map[string]interface{}{
		"message":      "Dashboard imported successfully",
		"dashboard_id": config.ID,
		"imported_at":  time.Now(),
	})
}