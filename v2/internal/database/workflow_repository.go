package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// workflowRepository implements WorkflowRepository interface
type workflowRepository struct {
	db *gorm.DB
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *gorm.DB) WorkflowRepository {
	return &workflowRepository{db: db}
}

// Create creates a new workflow
func (r *workflowRepository) Create(workflow *models.Workflow) error {
	return r.db.Create(workflow).Error
}

// GetByID retrieves a workflow by ID
func (r *workflowRepository) GetByID(id uuid.UUID) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.Where("id = ?", id).First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// GetByName retrieves a workflow by name
func (r *workflowRepository) GetByName(name string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.Where("name = ?", name).First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// List retrieves workflows with pagination and filtering
func (r *workflowRepository) List(limit, offset int, status string) ([]*models.Workflow, int64, error) {
	var workflows []*models.Workflow
	var total int64

	query := r.db.Model(&models.Workflow{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get workflows with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

// Update updates a workflow
func (r *workflowRepository) Update(workflow *models.Workflow) error {
	return r.db.Save(workflow).Error
}

// Delete deletes a workflow
func (r *workflowRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Workflow{}, "id = ?", id).Error
}

// Count returns the total number of workflows
func (r *workflowRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Workflow{}).Count(&count).Error
	return count, err
}

// CountByStatus returns the number of workflows with a specific status
func (r *workflowRepository) CountByStatus(status models.WorkflowStatus) (int64, error) {
	var count int64
	err := r.db.Model(&models.Workflow{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetByTriggerType retrieves workflows by trigger type
func (r *workflowRepository) GetByTriggerType(triggerType models.TriggerType) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	err := r.db.Where("definition->>'$.spec.triggers[*].type' = ?", triggerType).Find(&workflows).Error
	return workflows, err
}

// Search searches workflows by name or description
func (r *workflowRepository) Search(query string, limit, offset int) ([]*models.Workflow, int64, error) {
	var workflows []*models.Workflow
	var total int64

	searchQuery := r.db.Model(&models.Workflow{}).Where(
		"name ILIKE ? OR description ILIKE ?",
		"%"+query+"%", "%"+query+"%",
	)

	// Get total count
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get workflows with pagination
	err := searchQuery.Limit(limit).Offset(offset).Order("created_at DESC").Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

// GetActiveWorkflows retrieves all active workflows
func (r *workflowRepository) GetActiveWorkflows() ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	err := r.db.Where("status = ?", models.WorkflowStatusActive).Find(&workflows).Error
	return workflows, err
}

// GetWorkflowsByCreator retrieves workflows created by a specific user
func (r *workflowRepository) GetWorkflowsByCreator(createdBy string, limit, offset int) ([]*models.Workflow, int64, error) {
	var workflows []*models.Workflow
	var total int64

	query := r.db.Model(&models.Workflow{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get workflows with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

// GetRecentlyUpdated retrieves recently updated workflows
func (r *workflowRepository) GetRecentlyUpdated(since time.Time, limit int) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	err := r.db.Where("updated_at > ?", since).Limit(limit).Order("updated_at DESC").Find(&workflows).Error
	return workflows, err
}

// UpdateStatus updates only the status of a workflow
func (r *workflowRepository) UpdateStatus(id uuid.UUID, status models.WorkflowStatus) error {
	return r.db.Model(&models.Workflow{}).Where("id = ?", id).Update("status", status).Error
}

// BulkUpdateStatus updates the status of multiple workflows
func (r *workflowRepository) BulkUpdateStatus(ids []uuid.UUID, status models.WorkflowStatus) error {
	return r.db.Model(&models.Workflow{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetWorkflowStats retrieves workflow statistics
func (r *workflowRepository) GetWorkflowStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count by status
	statusStats := make(map[string]int64)
	statuses := []models.WorkflowStatus{
		models.WorkflowStatusDraft,
		models.WorkflowStatusActive,
		models.WorkflowStatusInactive,
		models.WorkflowStatusArchived,
	}

	for _, status := range statuses {
		count, err := r.CountByStatus(status)
		if err != nil {
			return nil, err
		}
		statusStats[string(status)] = count
	}

	stats["by_status"] = statusStats

	// Total count
	total, err := r.Count()
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Recent activity (workflows updated in last 24 hours)
	since := time.Now().Add(-24 * time.Hour)
	recentlyUpdated, err := r.GetRecentlyUpdated(since, 100)
	if err != nil {
		return nil, err
	}
	stats["recently_updated_count"] = len(recentlyUpdated)

	return stats, nil
}