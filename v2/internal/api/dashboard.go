package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// getDashboardOverview gets dashboard overview data
func (h *Handler) getDashboardOverview(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		// Default to last 24 hours
		end = time.Now().UTC()
		start = end.Add(-24 * time.Hour)
	}

	// Get overview data
	overview, err := h.services.DashboardService.GetOverview(start, end)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get dashboard overview", err)
		return
	}

	h.successResponse(c, overview)
}

// getWorkflowStatus gets workflow status summary
func (h *Handler) getWorkflowStatus(c *gin.Context) {
	// Parse time range
	start, end, err := h.parseTimeRange(c)
	if err != nil {
		// Default to last 24 hours
		end = time.Now().UTC()
		start = end.Add(-24 * time.Hour)
	}

	// Get workflow status
	status, err := h.services.DashboardService.GetWorkflowStatus(start, end)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get workflow status", err)
		return
	}

	h.successResponse(c, status)
}

// getSystemHealth gets system health metrics
func (h *Handler) getSystemHealth(c *gin.Context) {
	// Get system health
	health, err := h.services.DashboardService.GetSystemHealth()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get system health", err)
		return
	}

	h.successResponse(c, health)
}

// getLiveMetrics gets live metrics data
func (h *Handler) getLiveMetrics(c *gin.Context) {
	// Parse metric types
	metricTypes := c.QueryArray("metrics")
	if len(metricTypes) == 0 {
		// Default metrics
		metricTypes = []string{"executions", "errors", "latency", "throughput"}
	}

	// Get live metrics
	metrics, err := h.services.DashboardService.GetLiveMetrics(metricTypes)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get live metrics", err)
		return
	}

	h.successResponse(c, metrics)
}

// createDashboard creates a new dashboard
func (h *Handler) createDashboard(c *gin.Context) {
	var dashboard struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Layout      map[string]interface{} `json:"layout" binding:"required"`
		Widgets     []map[string]interface{} `json:"widgets" binding:"required"`
		Filters     map[string]interface{} `json:"filters"`
		Settings    map[string]interface{} `json:"settings"`
		IsPublic    bool                   `json:"is_public"`
		Tags        []string               `json:"tags"`
	}

	if err := h.validateRequestBody(c, &dashboard); err != nil {
		return
	}

	// Create dashboard
	createdDashboard, err := h.services.DashboardService.CreateDashboard(
		dashboard.Name,
		dashboard.Description,
		dashboard.Layout,
		dashboard.Widgets,
		dashboard.Filters,
		dashboard.Settings,
		dashboard.IsPublic,
		dashboard.Tags,
		h.getUserID(c),
	)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create dashboard", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"dashboard_id": createdDashboard.ID,
		"name":         dashboard.Name,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard created")

	c.JSON(http.StatusCreated, gin.H{
		"data":      createdDashboard,
		"timestamp": time.Now().UTC(),
	})
}

// listDashboards lists dashboards
func (h *Handler) listDashboards(c *gin.Context) {
	page, limit := h.parsePagination(c)

	// Parse filters
	filters := map[string]interface{}{}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if isPublic := c.Query("is_public"); isPublic != "" {
		filters["is_public"] = isPublic == "true"
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filters["tags"] = tags
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	// Get dashboards
	dashboards, total, err := h.services.DashboardService.ListDashboards(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list dashboards", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       dashboards,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// getDashboard gets a specific dashboard
func (h *Handler) getDashboard(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Get dashboard
	dashboard, err := h.services.DashboardService.GetDashboard(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Dashboard not found", err)
		return
	}

	h.successResponse(c, dashboard)
}

// updateDashboard updates a dashboard
func (h *Handler) updateDashboard(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var updates struct {
		Name        *string                 `json:"name"`
		Description *string                 `json:"description"`
		Layout      *map[string]interface{} `json:"layout"`
		Widgets     *[]map[string]interface{} `json:"widgets"`
		Filters     *map[string]interface{} `json:"filters"`
		Settings    *map[string]interface{} `json:"settings"`
		IsPublic    *bool                   `json:"is_public"`
		Tags        *[]string               `json:"tags"`
	}

	if err := h.validateRequestBody(c, &updates); err != nil {
		return
	}

	// Update dashboard
	updatedDashboard, err := h.services.DashboardService.UpdateDashboard(id, updates, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update dashboard", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"dashboard_id": id,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard updated")

	h.successResponse(c, updatedDashboard)
}

// deleteDashboard deletes a dashboard
func (h *Handler) deleteDashboard(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Delete dashboard
	if err := h.services.DashboardService.DeleteDashboard(id, h.getUserID(c)); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to delete dashboard", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"dashboard_id": id,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard deleted")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Dashboard deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// shareDashboard shares a dashboard
func (h *Handler) shareDashboard(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var shareData struct {
		Users       []string `json:"users"`
		Permissions string   `json:"permissions" binding:"required"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := h.validateRequestBody(c, &shareData); err != nil {
		return
	}

	// Share dashboard
	shareInfo, err := h.services.DashboardService.ShareDashboard(id, shareData.Users, shareData.Permissions, shareData.ExpiresAt, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to share dashboard", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"dashboard_id": id,
		"users":        shareData.Users,
		"permissions":  shareData.Permissions,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard shared")

	c.JSON(http.StatusOK, gin.H{
		"data":      shareInfo,
		"message":   "Dashboard shared successfully",
		"timestamp": time.Now().UTC(),
	})
}

// exportDashboard exports a dashboard
func (h *Handler) exportDashboard(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "json"
	}

	// Export dashboard
	exportData, filename, contentType, err := h.services.DashboardService.ExportDashboard(id, format)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to export dashboard", err)
		return
	}

	// Set headers for file download
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", string(rune(len(exportData))))

	// Write file data
	c.Data(http.StatusOK, contentType, exportData)

	logrus.WithFields(logrus.Fields{
		"dashboard_id": id,
		"format":       format,
		"filename":     filename,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard exported")
}

// importDashboard imports a dashboard
func (h *Handler) importDashboard(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "File is required", err)
		return
	}
	defer file.Close()

	overwrite := c.PostForm("overwrite") == "true"

	// Import dashboard
	dashboard, err := h.services.DashboardService.ImportDashboard(file, header.Filename, overwrite, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to import dashboard", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"dashboard_id": dashboard.ID,
		"filename":     header.Filename,
		"overwrite":    overwrite,
		"user_id":      h.getUserID(c),
	}).Info("Dashboard imported")

	c.JSON(http.StatusCreated, gin.H{
		"data":      dashboard,
		"message":   "Dashboard imported successfully",
		"timestamp": time.Now().UTC(),
	})
}