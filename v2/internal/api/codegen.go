package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// generateCode generates client code for a workflow
func (h *Handler) generateCode(c *gin.Context) {
	var request CodeGenRequest
	if err := h.validateRequestBody(c, &request); err != nil {
		return
	}

	// Validate workflow exists
	workflow, err := h.services.WorkflowService.GetByID(request.WorkflowID)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	// Validate language
	supportedLanguages := []string{"go", "typescript", "python", "java", "csharp"}
	validLanguage := false
	for _, lang := range supportedLanguages {
		if request.Language == lang {
			validLanguage = true
			break
		}
	}

	if !validLanguage {
		h.errorResponse(c, http.StatusBadRequest, "Unsupported language", nil)
		return
	}

	// Submit code generation job
	job, err := h.services.CodeGenService.GenerateCode(workflow, request.Language, request.Template, request.Options)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to start code generation", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"job_id":      job.ID,
		"workflow_id": request.WorkflowID,
		"language":    request.Language,
		"user_id":     h.getUserID(c),
	}).Info("Code generation job started")

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":    job.ID,
		"status":     job.Status,
		"message":    "Code generation job started",
		"timestamp": time.Now().UTC(),
	})
}

// getCodeGenStatus gets the status of a code generation job
func (h *Handler) getCodeGenStatus(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Get job status
	job, err := h.services.CodeGenService.GetJobStatus(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Code generation job not found", err)
		return
	}

	response := gin.H{
		"job_id":     job.ID,
		"status":      job.Status,
		"progress":    job.Progress,
		"created_at":  job.CreatedAt,
		"updated_at":  job.UpdatedAt,
		"timestamp":   time.Now().UTC(),
	}

	if job.CompletedAt != nil {
		response["completed_at"] = job.CompletedAt
		response["duration"] = job.GetDurationSeconds()
	}

	if job.Error != "" {
		response["error"] = job.Error
		response["error_code"] = job.ErrorCode
	}

	if job.IsCompleted() {
		response["download_url"] = "/api/v1/codegen/jobs/" + job.ID.String() + "/download"
		response["artifacts"] = job.Artifacts
	}

	c.JSON(http.StatusOK, response)
}

// downloadGeneratedCode downloads the generated code
func (h *Handler) downloadGeneratedCode(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Get job
	job, err := h.services.CodeGenService.GetJobStatus(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Code generation job not found", err)
		return
	}

	if !job.IsCompleted() {
		h.errorResponse(c, http.StatusBadRequest, "Code generation job is not completed", nil)
		return
	}

	if job.IsFailed() {
		h.errorResponse(c, http.StatusBadRequest, "Code generation job failed", nil)
		return
	}

	// Get generated code archive
	archiveData, filename, err := h.services.CodeGenService.GetGeneratedCode(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get generated code", err)
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", string(rune(len(archiveData))))

	// Write file data
	c.Data(http.StatusOK, "application/zip", archiveData)

	logrus.WithFields(logrus.Fields{
		"job_id":   id,
		"filename": filename,
		"user_id":  h.getUserID(c),
	}).Info("Generated code downloaded")
}

// listCodeGenTemplates lists available code generation templates
func (h *Handler) listCodeGenTemplates(c *gin.Context) {
	language := c.Query("language")
	category := c.Query("category")

	// Get templates
	templates, err := h.services.CodeGenService.ListTemplates(language, category)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list templates", err)
		return
	}

	h.successResponse(c, templates)
}

// listCodeGenJobs lists code generation jobs
func (h *Handler) listCodeGenJobs(c *gin.Context) {
	page, limit := h.parsePagination(c)

	// Parse filters
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if language := c.Query("language"); language != "" {
		filters["language"] = language
	}
	if workflowID := c.Query("workflow_id"); workflowID != "" {
		if id, err := uuid.Parse(workflowID); err == nil {
			filters["workflow_id"] = id
		}
	}

	// Parse time range
	if start, end, err := h.parseTimeRange(c); err == nil {
		filters["start_time"] = start
		filters["end_time"] = end
	}

	// Get jobs
	jobs, total, err := h.services.CodeGenService.ListJobs(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list code generation jobs", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       jobs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// createWorkflowVersion creates a new workflow version
func (h *Handler) createWorkflowVersion(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var versionData struct {
		Version         string                 `json:"version" binding:"required"`
		Description     string                 `json:"description"`
		Changelog       string                 `json:"changelog"`
		BreakingChanges bool                   `json:"breaking_changes"`
		Definition      map[string]interface{} `json:"definition"`
		InputSchema     map[string]interface{} `json:"input_schema"`
		OutputSchema    map[string]interface{} `json:"output_schema"`
		Config          map[string]interface{} `json:"config"`
	}

	if err := h.validateRequestBody(c, &versionData); err != nil {
		return
	}

	// Validate workflow exists
	workflow, err := h.services.WorkflowService.GetByID(workflowID)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	// Create version
	version, err := h.services.VersionService.CreateVersion(workflow, versionData.Version, versionData.Description, versionData.Changelog, versionData.BreakingChanges, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create workflow version", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"version_id":  version.ID,
		"workflow_id": workflowID,
		"version":     versionData.Version,
		"user_id":     h.getUserID(c),
	}).Info("Workflow version created")

	c.JSON(http.StatusCreated, gin.H{
		"data":      version,
		"timestamp": time.Now().UTC(),
	})
}

// listWorkflowVersions lists workflow versions
func (h *Handler) listWorkflowVersions(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	page, limit := h.parsePagination(c)

	// Parse filters
	filters := map[string]interface{}{
		"workflow_id": workflowID,
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Get versions
	versions, total, err := h.services.VersionService.ListVersions(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list workflow versions", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       versions,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// getWorkflowVersion gets a specific workflow version
func (h *Handler) getWorkflowVersion(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	versionStr := c.Param("version")
	if versionStr == "" {
		h.errorResponse(c, http.StatusBadRequest, "Version parameter is required", nil)
		return
	}

	// Get version
	version, err := h.services.VersionService.GetVersion(workflowID, versionStr)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow version not found", err)
		return
	}

	h.successResponse(c, version)
}

// rollbackWorkflowVersion rolls back to a specific workflow version
func (h *Handler) rollbackWorkflowVersion(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	versionStr := c.Param("version")
	if versionStr == "" {
		h.errorResponse(c, http.StatusBadRequest, "Version parameter is required", nil)
		return
	}

	var rollbackData struct {
		Reason      string `json:"reason"`
		ForceRollback bool `json:"force_rollback"`
	}

	if err := h.validateRequestBody(c, &rollbackData); err != nil {
		return
	}

	// Perform rollback
	rollbackInfo, err := h.services.VersionService.RollbackVersion(workflowID, versionStr, rollbackData.Reason, rollbackData.ForceRollback, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to rollback workflow version", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"workflow_id": workflowID,
		"version":     versionStr,
		"reason":      rollbackData.Reason,
		"user_id":     h.getUserID(c),
	}).Info("Workflow version rolled back")

	c.JSON(http.StatusOK, gin.H{
		"data":      rollbackInfo,
		"message":   "Workflow version rolled back successfully",
		"timestamp": time.Now().UTC(),
	})
}

// compareWorkflowVersions compares two workflow versions
func (h *Handler) compareWorkflowVersions(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	fromVersion := c.Param("from")
	toVersion := c.Param("to")

	if fromVersion == "" || toVersion == "" {
		h.errorResponse(c, http.StatusBadRequest, "Both from and to version parameters are required", nil)
		return
	}

	// Compare versions
	comparison, err := h.services.VersionService.CompareVersions(workflowID, fromVersion, toVersion)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to compare workflow versions", err)
		return
	}

	h.successResponse(c, comparison)
}

// deployWorkflowVersion deploys a workflow version
func (h *Handler) deployWorkflowVersion(c *gin.Context) {
	workflowID, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	versionStr := c.Param("version")
	if versionStr == "" {
		h.errorResponse(c, http.StatusBadRequest, "Version parameter is required", nil)
		return
	}

	var deploymentData struct {
		Environment string                 `json:"environment" binding:"required"`
		Strategy    string                 `json:"strategy"`
		Config      map[string]interface{} `json:"config"`
	}

	if err := h.validateRequestBody(c, &deploymentData); err != nil {
		return
	}

	// Deploy version
	deployment, err := h.services.VersionService.DeployVersion(workflowID, versionStr, deploymentData.Environment, deploymentData.Strategy, deploymentData.Config, h.getUserID(c))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to deploy workflow version", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"deployment_id": deployment.ID,
		"workflow_id":   workflowID,
		"version":       versionStr,
		"environment":   deploymentData.Environment,
		"user_id":       h.getUserID(c),
	}).Info("Workflow version deployment started")

	c.JSON(http.StatusAccepted, gin.H{
		"data":      deployment,
		"message":   "Workflow version deployment started",
		"timestamp": time.Now().UTC(),
	})
}