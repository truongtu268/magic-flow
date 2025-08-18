package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// workflowVersionRepository implements WorkflowVersionRepository interface
type workflowVersionRepository struct {
	db *gorm.DB
}

// NewWorkflowVersionRepository creates a new workflow version repository
func NewWorkflowVersionRepository(db *gorm.DB) WorkflowVersionRepository {
	return &workflowVersionRepository{db: db}
}

// Create creates a new workflow version
func (r *workflowVersionRepository) Create(version *models.WorkflowVersion) error {
	return r.db.Create(version).Error
}

// GetByID retrieves a workflow version by ID
func (r *workflowVersionRepository) GetByID(id uuid.UUID) (*models.WorkflowVersion, error) {
	var version models.WorkflowVersion
	err := r.db.Where("id = ?", id).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetByWorkflowID retrieves all versions for a specific workflow
func (r *workflowVersionRepository) GetByWorkflowID(workflowID uuid.UUID, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{}).Where("workflow_id = ?", workflowID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("version DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetByWorkflowIDAndVersion retrieves a specific version of a workflow
func (r *workflowVersionRepository) GetByWorkflowIDAndVersion(workflowID uuid.UUID, version int) (*models.WorkflowVersion, error) {
	var workflowVersion models.WorkflowVersion
	err := r.db.Where("workflow_id = ? AND version = ?", workflowID, version).First(&workflowVersion).Error
	if err != nil {
		return nil, err
	}
	return &workflowVersion, nil
}

// GetLatestVersion retrieves the latest version of a workflow
func (r *workflowVersionRepository) GetLatestVersion(workflowID uuid.UUID) (*models.WorkflowVersion, error) {
	var version models.WorkflowVersion
	err := r.db.Where("workflow_id = ?", workflowID).Order("version DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetActiveVersion retrieves the currently active version of a workflow
func (r *workflowVersionRepository) GetActiveVersion(workflowID uuid.UUID) (*models.WorkflowVersion, error) {
	var version models.WorkflowVersion
	err := r.db.Where("workflow_id = ? AND is_active = ?", workflowID, true).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// Update updates a workflow version
func (r *workflowVersionRepository) Update(version *models.WorkflowVersion) error {
	return r.db.Save(version).Error
}

// Delete deletes a workflow version
func (r *workflowVersionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.WorkflowVersion{}, "id = ?", id).Error
}

// SetActiveVersion sets a specific version as active and deactivates others
func (r *workflowVersionRepository) SetActiveVersion(workflowID uuid.UUID, version int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deactivate all versions for this workflow
		if err := tx.Model(&models.WorkflowVersion{}).Where("workflow_id = ?", workflowID).Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the specified version
		return tx.Model(&models.WorkflowVersion{}).Where("workflow_id = ? AND version = ?", workflowID, version).Update("is_active", true).Error
	})
}

// GetNextVersionNumber retrieves the next version number for a workflow
func (r *workflowVersionRepository) GetNextVersionNumber(workflowID uuid.UUID) (int, error) {
	var maxVersion struct {
		MaxVersion *int `gorm:"column:max_version"`
	}

	err := r.db.Model(&models.WorkflowVersion{}).Select("MAX(version) as max_version").Where("workflow_id = ?", workflowID).Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}

	if maxVersion.MaxVersion == nil {
		return 1, nil
	}

	return *maxVersion.MaxVersion + 1, nil
}

// List retrieves workflow versions with filtering and pagination
func (r *workflowVersionRepository) List(limit, offset int, workflowID *uuid.UUID, isActive *bool) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{})

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetVersionsByTimeRange retrieves versions within a time range
func (r *workflowVersionRepository) GetVersionsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{})

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetVersionsByCreator retrieves versions created by a specific user
func (r *workflowVersionRepository) GetVersionsByCreator(createdBy string, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// CompareVersions retrieves two versions for comparison
func (r *workflowVersionRepository) CompareVersions(workflowID uuid.UUID, version1, version2 int) (*models.WorkflowVersion, *models.WorkflowVersion, error) {
	var v1, v2 models.WorkflowVersion

	// Get first version
	err := r.db.Where("workflow_id = ? AND version = ?", workflowID, version1).First(&v1).Error
	if err != nil {
		return nil, nil, err
	}

	// Get second version
	err = r.db.Where("workflow_id = ? AND version = ?", workflowID, version2).First(&v2).Error
	if err != nil {
		return nil, nil, err
	}

	return &v1, &v2, nil
}

// GetVersionHistory retrieves the complete version history for a workflow
func (r *workflowVersionRepository) GetVersionHistory(workflowID uuid.UUID) ([]*models.WorkflowVersion, error) {
	var versions []*models.WorkflowVersion
	err := r.db.Where("workflow_id = ?", workflowID).Order("version ASC").Find(&versions).Error
	return versions, err
}

// GetVersionStats retrieves version statistics for a workflow
func (r *workflowVersionRepository) GetVersionStats(workflowID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total versions
	var totalVersions int64
	if err := r.db.Model(&models.WorkflowVersion{}).Where("workflow_id = ?", workflowID).Count(&totalVersions).Error; err != nil {
		return nil, err
	}
	stats["total_versions"] = totalVersions

	// Active version
	activeVersion, err := r.GetActiveVersion(workflowID)
	if err == nil {
		stats["active_version"] = activeVersion.Version
	} else {
		stats["active_version"] = nil
	}

	// Latest version
	latestVersion, err := r.GetLatestVersion(workflowID)
	if err == nil {
		stats["latest_version"] = latestVersion.Version
	} else {
		stats["latest_version"] = nil
	}

	// Version creation frequency (versions per day in the last 30 days)
	thirtyDaysAgo := time.Now().UTC().Add(-30 * 24 * time.Hour)
	var recentVersions int64
	if err := r.db.Model(&models.WorkflowVersion{}).Where("workflow_id = ? AND created_at >= ?", workflowID, thirtyDaysAgo).Count(&recentVersions).Error; err != nil {
		return nil, err
	}
	stats["recent_versions_30d"] = recentVersions

	return stats, nil
}

// RollbackToVersion creates a new version based on an existing version (rollback)
func (r *workflowVersionRepository) RollbackToVersion(workflowID uuid.UUID, targetVersion int, createdBy string) (*models.WorkflowVersion, error) {
	// Get the target version to rollback to
	targetVersionData, err := r.GetByWorkflowIDAndVersion(workflowID, targetVersion)
	if err != nil {
		return nil, err
	}

	// Get the next version number
	nextVersion, err := r.GetNextVersionNumber(workflowID)
	if err != nil {
		return nil, err
	}

	// Create new version with the target version's definition
	newVersion := &models.WorkflowVersion{
		ID:                uuid.New(),
		WorkflowID:        workflowID,
		Version:           nextVersion,
		Definition:        targetVersionData.Definition,
		ChangeDescription: "Rollback to version " + string(rune(targetVersion)),
		CreatedBy:         createdBy,
		IsActive:          false,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	err = r.Create(newVersion)
	if err != nil {
		return nil, err
	}

	return newVersion, nil
}

// GetVersionDiff calculates the differences between two versions
func (r *workflowVersionRepository) GetVersionDiff(workflowID uuid.UUID, version1, version2 int) (map[string]interface{}, error) {
	v1, v2, err := r.CompareVersions(workflowID, version1, version2)
	if err != nil {
		return nil, err
	}

	diff := map[string]interface{}{
		"version_1": map[string]interface{}{
			"version":     v1.Version,
			"created_at":  v1.CreatedAt,
			"created_by":  v1.CreatedBy,
			"description": v1.ChangeDescription,
			"definition":  v1.Definition,
		},
		"version_2": map[string]interface{}{
			"version":     v2.Version,
			"created_at":  v2.CreatedAt,
			"created_by":  v2.CreatedBy,
			"description": v2.ChangeDescription,
			"definition":  v2.Definition,
		},
		"has_changes": v1.Definition != v2.Definition,
	}

	return diff, nil
}

// CleanupOldVersions deletes old versions while keeping a specified number of recent versions
func (r *workflowVersionRepository) CleanupOldVersions(workflowID uuid.UUID, keepCount int) (int64, error) {
	// Get versions to keep (most recent ones)
	var versionsToKeep []int
	err := r.db.Model(&models.WorkflowVersion{}).Select("version").Where("workflow_id = ?", workflowID).Order("version DESC").Limit(keepCount).Pluck("version", &versionsToKeep).Error
	if err != nil {
		return 0, err
	}

	if len(versionsToKeep) == 0 {
		return 0, nil
	}

	// Delete versions not in the keep list
	result := r.db.Where("workflow_id = ? AND version NOT IN ?", workflowID, versionsToKeep).Delete(&models.WorkflowVersion{})

	return result.RowsAffected, result.Error
}

// GetVersionsByStatus retrieves versions by their deployment status
func (r *workflowVersionRepository) GetVersionsByStatus(status string, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	// Note: This assumes there's a status field in the WorkflowVersion model
	// If not, this method might need to be adjusted based on the actual model structure
	query := r.db.Model(&models.WorkflowVersion{})

	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetActiveVersions retrieves all currently active versions across all workflows
func (r *workflowVersionRepository) GetActiveVersions(limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{}).Where("is_active = ?", true)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// SearchVersions searches versions by change description
func (r *workflowVersionRepository) SearchVersions(searchTerm string, limit, offset int) ([]*models.WorkflowVersion, int64, error) {
	var versions []*models.WorkflowVersion
	var total int64

	query := r.db.Model(&models.WorkflowVersion{}).Where("change_description ILIKE ?", "%"+searchTerm+"%")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get versions with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetVersionCreationTrend retrieves version creation trend data
func (r *workflowVersionRepository) GetVersionCreationTrend(workflowID *uuid.UUID, startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Determine the date truncation based on interval
	var dateTrunc string
	switch interval {
	case "hour":
		dateTrunc = "date_trunc('hour', created_at)"
	case "day":
		dateTrunc = "date_trunc('day', created_at)"
	case "week":
		dateTrunc = "date_trunc('week', created_at)"
	case "month":
		dateTrunc = "date_trunc('month', created_at)"
	default:
		dateTrunc = "date_trunc('day', created_at)"
	}

	query := r.db.Model(&models.WorkflowVersion{}).Select(dateTrunc + " as time_bucket, COUNT(*) as count")

	if workflowID != nil {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query = query.Group("time_bucket").Order("time_bucket ASC")

	var rawResults []struct {
		TimeBucket time.Time `gorm:"column:time_bucket"`
		Count      int64     `gorm:"column:count"`
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert to map format
	for _, result := range rawResults {
		results = append(results, map[string]interface{}{
			"time":  result.TimeBucket,
			"count": result.Count,
		})
	}

	return results, nil
}