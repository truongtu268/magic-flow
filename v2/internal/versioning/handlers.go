package versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// Handlers provides HTTP handlers for versioning operations
type Handlers struct {
	manager *Manager
}

// NewHandlers creates a new handlers instance
func NewHandlers(repoManager database.RepositoryManager) *Handlers {
	return &Handlers{
		manager: NewManager(repoManager),
	}
}

// CreateVersionRequest represents a request to create a new version
type CreateVersionRequest struct {
	ChangeType    ChangeType             `json:"change_type" binding:"required"`
	Summary       string                 `json:"summary" binding:"required"`
	Details       string                 `json:"details"`
	Definition    map[string]interface{} `json:"definition" binding:"required"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ActivateVersionRequest represents a request to activate a version
type ActivateVersionRequest struct {
	VersionID uuid.UUID `json:"version_id" binding:"required"`
}

// RollbackRequest represents a request to rollback to a previous version
type RollbackRequest struct {
	TargetVersionID uuid.UUID `json:"target_version_id" binding:"required"`
	Reason          string    `json:"reason" binding:"required"`
}

// CompareVersionsRequest represents a request to compare two versions
type CompareVersionsRequest struct {
	Version1ID uuid.UUID `json:"version1_id" binding:"required"`
	Version2ID uuid.UUID `json:"version2_id" binding:"required"`
}

// CreateVersion creates a new version of a workflow
// @Summary Create workflow version
// @Description Create a new version of a workflow with specified changes
// @Tags versioning
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param request body CreateVersionRequest true "Version creation request"
// @Success 201 {object} models.WorkflowVersion
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions [post]
func (h *Handlers) CreateVersion(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (would be set by auth middleware)
	userID := h.getUserIDFromContext(c)

	// Create version changes
	changes := VersionChanges{
		ChangeType:    req.ChangeType,
		Summary:       req.Summary,
		Details:       req.Details,
		NewDefinition: req.Definition,
		CreatedBy:     userID,
		Metadata:      req.Metadata,
	}

	// Create the version
	version, err := h.manager.CreateVersion(c.Request.Context(), workflowID, changes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, version)
}

// GetVersionHistory gets the version history for a workflow
// @Summary Get workflow version history
// @Description Get the complete version history for a workflow
// @Tags versioning
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param limit query int false "Limit number of versions returned"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.WorkflowVersion
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions [get]
func (h *Handlers) GetVersionHistory(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Get pagination parameters
	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get version history
	versions, err := h.manager.GetVersionHistory(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply pagination
	start := offset
	if start > len(versions) {
		start = len(versions)
	}
	end := start + limit
	if end > len(versions) {
		end = len(versions)
	}

	paginatedVersions := versions[start:end]

	c.JSON(http.StatusOK, gin.H{
		"versions": paginatedVersions,
		"total":    len(versions),
		"limit":    limit,
		"offset":   offset,
	})
}

// GetVersion gets a specific version by ID
// @Summary Get workflow version
// @Description Get a specific version of a workflow by version ID
// @Tags versioning
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param version_id path string true "Version ID"
// @Success 200 {object} models.WorkflowVersion
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/{version_id} [get]
func (h *Handlers) GetVersion(c *gin.Context) {
	versionIDStr := c.Param("version_id")
	versionID, err := uuid.Parse(versionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	// Get the version (this would use a repository method)
	// For now, we'll return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Version details",
		"version_id": versionID,
	})
}

// ActivateVersion activates a specific version
// @Summary Activate workflow version
// @Description Activate a specific version of a workflow
// @Tags versioning
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param request body ActivateVersionRequest true "Version activation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/activate [post]
func (h *Handlers) ActivateVersion(c *gin.Context) {
	var req ActivateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Activate the version
	err := h.manager.ActivateVersion(c.Request.Context(), req.VersionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Version activated successfully",
		"version_id": req.VersionID,
		"activated_at": time.Now(),
	})
}

// RollbackToVersion rolls back a workflow to a previous version
// @Summary Rollback workflow version
// @Description Rollback a workflow to a previous version
// @Tags versioning
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param request body RollbackRequest true "Rollback request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/rollback [post]
func (h *Handlers) RollbackToVersion(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var req RollbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform rollback
	err = h.manager.RollbackToVersion(c.Request.Context(), workflowID, req.TargetVersionID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rollback completed successfully",
		"workflow_id": workflowID,
		"target_version_id": req.TargetVersionID,
		"reason": req.Reason,
		"rolled_back_at": time.Now(),
	})
}

// CompareVersions compares two versions of a workflow
// @Summary Compare workflow versions
// @Description Compare two versions of a workflow and return differences
// @Tags versioning
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param request body CompareVersionsRequest true "Version comparison request"
// @Success 200 {object} VersionComparison
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/compare [post]
func (h *Handlers) CompareVersions(c *gin.Context) {
	var req CompareVersionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Compare versions
	comparison, err := h.manager.CompareVersions(c.Request.Context(), req.Version1ID, req.Version2ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

// GetCompatibilityMatrix gets the compatibility matrix for all versions
// @Summary Get version compatibility matrix
// @Description Get compatibility information between all versions of a workflow
// @Tags versioning
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Success 200 {object} CompatibilityMatrix
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/compatibility [get]
func (h *Handlers) GetCompatibilityMatrix(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Get compatibility matrix
	matrix, err := h.manager.GetCompatibilityMatrix(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrix)
}

// GetMigrationPlan gets the migration plan between two versions
// @Summary Get migration plan
// @Description Get the migration plan for upgrading between two versions
// @Tags versioning
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param from_version query string true "From version ID"
// @Param to_version query string true "To version ID"
// @Success 200 {object} MigrationPlan
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/migration-plan [get]
func (h *Handlers) GetMigrationPlan(c *gin.Context) {
	fromVersionStr := c.Query("from_version")
	toVersionStr := c.Query("to_version")

	if fromVersionStr == "" || toVersionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both from_version and to_version are required"})
		return
	}

	fromVersionID, err := uuid.Parse(fromVersionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_version ID"})
		return
	}

	toVersionID, err := uuid.Parse(toVersionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_version ID"})
		return
	}

	// Get migration plan
	plan, err := h.manager.GetMigrationPlan(c.Request.Context(), fromVersionID, toVersionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// ValidateVersion validates a version before creation
// @Summary Validate workflow version
// @Description Validate a workflow version definition before creation
// @Tags versioning
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param request body CreateVersionRequest true "Version validation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/validate [post]
func (h *Handlers) ValidateVersion(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := h.getUserIDFromContext(c)

	// Create version changes for validation
	changes := VersionChanges{
		ChangeType:    req.ChangeType,
		Summary:       req.Summary,
		Details:       req.Details,
		NewDefinition: req.Definition,
		CreatedBy:     userID,
		Metadata:      req.Metadata,
	}

	// Validate the version
	err = h.manager.ValidateVersion(c.Request.Context(), workflowID, changes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"message": "Version definition is valid",
	})
}

// GetVersionMetrics gets metrics for version management
// @Summary Get version metrics
// @Description Get metrics and statistics for workflow version management
// @Tags versioning
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Success 200 {object} VersionMetrics
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflows/{workflow_id}/versions/metrics [get]
func (h *Handlers) GetVersionMetrics(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Get version history to calculate metrics
	versions, err := h.manager.GetVersionHistory(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate metrics
	metrics := h.calculateVersionMetrics(workflowID, versions)

	c.JSON(http.StatusOK, metrics)
}

// Helper methods

func (h *Handlers) getUserIDFromContext(c *gin.Context) uuid.UUID {
	// This would typically extract user ID from JWT token or session
	// For now, return a placeholder UUID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.New() // Placeholder
	}

	if userID, ok := userIDStr.(uuid.UUID); ok {
		return userID
	}

	if userIDString, ok := userIDStr.(string); ok {
		if userID, err := uuid.Parse(userIDString); err == nil {
			return userID
		}
	}

	return uuid.New() // Fallback
}

func (h *Handlers) calculateVersionMetrics(workflowID uuid.UUID, versions []*models.WorkflowVersion) *VersionMetrics {
	metrics := &VersionMetrics{
		WorkflowID:    workflowID,
		TotalVersions: len(versions),
		VersionFrequency: make(map[string]int),
	}

	if len(versions) == 0 {
		return metrics
	}

	// Find active version and latest version date
	for _, version := range versions {
		if version.IsActive {
			metrics.ActiveVersion = version.Version
		}
		if version.CreatedAt.After(metrics.LastVersionDate) {
			metrics.LastVersionDate = version.CreatedAt
		}

		// Count change types
		metrics.VersionFrequency[version.ChangeType]++
	}

	// Calculate success rate (simplified - would need actual migration data)
	metrics.SuccessRate = 0.95 // Placeholder

	// Calculate average migration time (simplified)
	metrics.AverageMigrationTime = 5 * time.Minute // Placeholder

	return metrics
}

// RegisterRoutes registers all versioning routes
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		workflows := v1.Group("/workflows/:workflow_id")
		{
			versions := workflows.Group("/versions")
			{
				versions.POST("", h.CreateVersion)
				versions.GET("", h.GetVersionHistory)
				versions.GET("/:version_id", h.GetVersion)
				versions.POST("/activate", h.ActivateVersion)
				versions.POST("/rollback", h.RollbackToVersion)
				versions.POST("/compare", h.CompareVersions)
				versions.GET("/compatibility", h.GetCompatibilityMatrix)
				versions.GET("/migration-plan", h.GetMigrationPlan)
				versions.POST("/validate", h.ValidateVersion)
				versions.GET("/metrics", h.GetVersionMetrics)
			}
		}
	}
}