package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magic-flow/v2/pkg/models"
	"github.com/sirupsen/logrus"
)

// createWorkflow creates a new workflow
func (h *Handler) createWorkflow(c *gin.Context) {
	var workflow models.Workflow
	if err := h.validateRequestBody(c, &workflow); err != nil {
		return
	}

	// Set creator
	userID := h.getUserID(c)
	workflow.CreatedBy = userID

	// Validate workflow
	if err := workflow.Validate(); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Workflow validation failed", err)
		return
	}

	// Create workflow
	createdWorkflow, err := h.services.WorkflowService.Create(&workflow)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create workflow", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"workflow_id": createdWorkflow.ID,
		"name":        createdWorkflow.Name,
		"user_id":     userID,
	}).Info("Workflow created")

	c.JSON(http.StatusCreated, gin.H{
		"data":      createdWorkflow,
		"timestamp": time.Now().UTC(),
	})
}

// listWorkflows lists all workflows with pagination and filtering
func (h *Handler) listWorkflows(c *gin.Context) {
	page, limit := h.parsePagination(c)
	
	// Parse filters
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if tags := c.Query("tags"); tags != "" {
		filters["tags"] = tags
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	// Get workflows
	workflows, total, err := h.services.WorkflowService.List(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list workflows", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       workflows,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// getWorkflow gets a workflow by ID
func (h *Handler) getWorkflow(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	workflow, err := h.services.WorkflowService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	h.successResponse(c, workflow)
}

// updateWorkflow updates an existing workflow
func (h *Handler) updateWorkflow(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var updateData models.Workflow
	if err := h.validateRequestBody(c, &updateData); err != nil {
		return
	}

	// Get existing workflow
	existingWorkflow, err := h.services.WorkflowService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	// Update fields
	updateData.ID = id
	updateData.CreatedBy = existingWorkflow.CreatedBy
	updateData.CreatedAt = existingWorkflow.CreatedAt
	updateData.UpdatedBy = h.getUserID(c)

	// Validate updated workflow
	if err := updateData.Validate(); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Workflow validation failed", err)
		return
	}

	// Update workflow
	updatedWorkflow, err := h.services.WorkflowService.Update(&updateData)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update workflow", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"workflow_id": id,
		"name":        updatedWorkflow.Name,
		"user_id":     h.getUserID(c),
	}).Info("Workflow updated")

	h.successResponse(c, updatedWorkflow)
}

// deleteWorkflow deletes a workflow
func (h *Handler) deleteWorkflow(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Check if workflow has active executions
	activeExecutions, err := h.services.ExecutionService.CountActiveByWorkflowID(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to check active executions", err)
		return
	}

	if activeExecutions > 0 {
		h.errorResponse(c, http.StatusConflict, "Cannot delete workflow with active executions", nil)
		return
	}

	// Delete workflow
	if err := h.services.WorkflowService.Delete(id); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to delete workflow", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"workflow_id": id,
		"user_id":     h.getUserID(c),
	}).Info("Workflow deleted")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Workflow deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// validateWorkflow validates a workflow definition
func (h *Handler) validateWorkflow(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Get workflow
	workflow, err := h.services.WorkflowService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	// Validate workflow using the engine
	validationResult, err := h.workflowEngine.ValidateWorkflow(workflow)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to validate workflow", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      validationResult.Valid,
		"errors":     validationResult.Errors,
		"warnings":   validationResult.Warnings,
		"timestamp":  time.Now().UTC(),
	})
}

// executeWorkflow executes a workflow
func (h *Handler) executeWorkflow(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	var request ExecutionRequest
	if err := h.validateRequestBody(c, &request); err != nil {
		return
	}

	// Get workflow
	workflow, err := h.services.WorkflowService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Workflow not found", err)
		return
	}

	// Check if workflow is active
	if workflow.Status != models.WorkflowStatusActive {
		h.errorResponse(c, http.StatusBadRequest, "Workflow is not active", nil)
		return
	}

	// Create execution
	execution := &models.Execution{
		WorkflowID:    id,
		WorkflowName:  workflow.Name,
		Status:        models.ExecutionStatusPending,
		Input:         request.Input,
		Environment:   request.Environment,
		Tags:          request.Tags,
		Priority:      request.Priority,
		TriggeredBy:   h.getUserID(c),
		TriggerType:   models.TriggerTypeManual,
	}

	if request.ScheduledAt != nil {
		execution.ScheduledAt = request.ScheduledAt
		execution.TriggerType = models.TriggerTypeScheduled
	}

	// Save execution
	createdExecution, err := h.services.ExecutionService.Create(execution)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create execution", err)
		return
	}

	// Submit to workflow engine
	if request.ScheduledAt == nil || request.ScheduledAt.Before(time.Now()) {
		if err := h.workflowEngine.SubmitExecution(createdExecution); err != nil {
			// Update execution status to failed
			createdExecution.Fail(err, "ENGINE_SUBMIT_ERROR")
			h.services.ExecutionService.Update(createdExecution)
			
			h.errorResponse(c, http.StatusInternalServerError, "Failed to submit execution", err)
			return
		}
	}

	logrus.WithFields(logrus.Fields{
		"execution_id": createdExecution.ID,
		"workflow_id":  id,
		"user_id":      h.getUserID(c),
	}).Info("Workflow execution started")

	c.JSON(http.StatusCreated, gin.H{
		"data":      createdExecution,
		"timestamp": time.Now().UTC(),
	})
}

// getExecution gets an execution by ID
func (h *Handler) getExecution(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	execution, err := h.services.ExecutionService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Execution not found", err)
		return
	}

	h.successResponse(c, execution)
}

// getExecutionStatus gets execution status
func (h *Handler) getExecutionStatus(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	execution, err := h.services.ExecutionService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Execution not found", err)
		return
	}

	// Get step executions
	stepExecutions, err := h.services.ExecutionService.GetStepExecutions(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get step executions", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution":       execution,
		"step_executions": stepExecutions,
		"progress":        execution.GetProgress(),
		"timestamp":       time.Now().UTC(),
	})
}

// getExecutionResults gets execution results
func (h *Handler) getExecutionResults(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	execution, err := h.services.ExecutionService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Execution not found", err)
		return
	}

	if !execution.IsFinished() {
		h.errorResponse(c, http.StatusBadRequest, "Execution is not finished", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution": execution,
		"output":    execution.Output,
		"error":     execution.Error,
		"timestamp": time.Now().UTC(),
	})
}

// listExecutions lists executions with pagination and filtering
func (h *Handler) listExecutions(c *gin.Context) {
	page, limit := h.parsePagination(c)
	
	// Parse filters
	filters := map[string]interface{}{}
	if workflowID := c.Query("workflow_id"); workflowID != "" {
		if id, err := uuid.Parse(workflowID); err == nil {
			filters["workflow_id"] = id
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if triggeredBy := c.Query("triggered_by"); triggeredBy != "" {
		filters["triggered_by"] = triggeredBy
	}
	if environment := c.Query("environment"); environment != "" {
		filters["environment"] = environment
	}

	// Parse time range
	if start, end, err := h.parseTimeRange(c); err == nil {
		filters["start_time"] = start
		filters["end_time"] = end
	}

	// Get executions
	executions, total, err := h.services.ExecutionService.List(page, limit, filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to list executions", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       executions,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}

// cancelExecution cancels an execution
func (h *Handler) cancelExecution(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	execution, err := h.services.ExecutionService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Execution not found", err)
		return
	}

	if execution.IsFinished() {
		h.errorResponse(c, http.StatusBadRequest, "Execution is already finished", nil)
		return
	}

	// Cancel execution in engine
	if err := h.workflowEngine.CancelExecution(id); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to cancel execution", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"execution_id": id,
		"user_id":      h.getUserID(c),
	}).Info("Execution cancelled")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Execution cancelled successfully",
		"timestamp": time.Now().UTC(),
	})
}

// retryExecution retries a failed execution
func (h *Handler) retryExecution(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	execution, err := h.services.ExecutionService.GetByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Execution not found", err)
		return
	}

	if execution.Status != models.ExecutionStatusFailed {
		h.errorResponse(c, http.StatusBadRequest, "Only failed executions can be retried", nil)
		return
	}

	// Create new execution for retry
	retryExecution := &models.Execution{
		WorkflowID:       execution.WorkflowID,
		WorkflowName:     execution.WorkflowName,
		WorkflowVersionID: execution.WorkflowVersionID,
		Status:           models.ExecutionStatusPending,
		Input:            execution.Input,
		Environment:      execution.Environment,
		Tags:             execution.Tags,
		Priority:         execution.Priority,
		TriggeredBy:      h.getUserID(c),
		TriggerType:      models.TriggerTypeRetry,
		ParentExecutionID: &execution.ID,
	}

	// Save retry execution
	createdExecution, err := h.services.ExecutionService.Create(retryExecution)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create retry execution", err)
		return
	}

	// Submit to workflow engine
	if err := h.workflowEngine.SubmitExecution(createdExecution); err != nil {
		createdExecution.Fail(err, "ENGINE_SUBMIT_ERROR")
		h.services.ExecutionService.Update(createdExecution)
		
		h.errorResponse(c, http.StatusInternalServerError, "Failed to submit retry execution", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"execution_id":       createdExecution.ID,
		"original_execution": id,
		"user_id":            h.getUserID(c),
	}).Info("Execution retried")

	c.JSON(http.StatusCreated, gin.H{
		"data":      createdExecution,
		"timestamp": time.Now().UTC(),
	})
}

// getExecutionLogs gets execution logs
func (h *Handler) getExecutionLogs(c *gin.Context) {
	id, err := h.parseUUID(c, "id")
	if err != nil {
		return
	}

	// Parse pagination for logs
	page, limit := h.parsePagination(c)
	level := c.Query("level") // debug, info, warn, error

	// Get logs from execution service
	logs, total, err := h.services.ExecutionService.GetLogs(id, page, limit, level)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get execution logs", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListResponse{
		Data:       logs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Timestamp:  time.Now().UTC(),
	})
}