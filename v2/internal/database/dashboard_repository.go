package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// dashboardRepository implements DashboardRepository interface
type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// Create creates a new dashboard
func (r *dashboardRepository) Create(dashboard *models.Dashboard) error {
	return r.db.Create(dashboard).Error
}

// GetByID retrieves a dashboard by ID
func (r *dashboardRepository) GetByID(id uuid.UUID) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.Where("id = ?", id).First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// GetByName retrieves a dashboard by name
func (r *dashboardRepository) GetByName(name string) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.Where("name = ?", name).First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// List retrieves dashboards with filtering and pagination
func (r *dashboardRepository) List(limit, offset int, createdBy string, isPublic *bool) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{})

	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// Update updates a dashboard
func (r *dashboardRepository) Update(dashboard *models.Dashboard) error {
	return r.db.Save(dashboard).Error
}

// Delete deletes a dashboard
func (r *dashboardRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Dashboard{}, "id = ?", id).Error
}

// GetDashboardsByCreator retrieves dashboards created by a specific user
func (r *dashboardRepository) GetDashboardsByCreator(createdBy string, limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// GetPublicDashboards retrieves all public dashboards
func (r *dashboardRepository) GetPublicDashboards(limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{}).Where("is_public = ?", true)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// GetDashboardsByTimeRange retrieves dashboards created within a time range
func (r *dashboardRepository) GetDashboardsByTimeRange(startTime, endTime *time.Time, limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{})

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

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// SearchDashboards searches dashboards by name or description
func (r *dashboardRepository) SearchDashboards(searchTerm string, limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{}).Where("name ILIKE ? OR description ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// GetDashboardsByTag retrieves dashboards with a specific tag
func (r *dashboardRepository) GetDashboardsByTag(tag string, limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	// Note: This assumes tags are stored as JSON array in the dashboard model
	// The actual implementation might need to be adjusted based on the Dashboard model structure
	query := r.db.Model(&models.Dashboard{}).Where("tags @> ?", `["`+tag+`"]`)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// UpdateDashboardAccess updates dashboard access settings
func (r *dashboardRepository) UpdateDashboardAccess(id uuid.UUID, isPublic bool, shareToken string) error {
	return r.db.Model(&models.Dashboard{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_public":    isPublic,
		"share_token":  shareToken,
		"updated_at":   time.Now().UTC(),
	}).Error
}

// GetDashboardByShareToken retrieves a dashboard by its share token
func (r *dashboardRepository) GetDashboardByShareToken(shareToken string) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.Where("share_token = ?", shareToken).First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// UpdateDashboardLastAccessed updates the last accessed time for a dashboard
func (r *dashboardRepository) UpdateDashboardLastAccessed(id uuid.UUID, accessedAt time.Time) error {
	return r.db.Model(&models.Dashboard{}).Where("id = ?", id).Update("last_accessed_at", accessedAt).Error
}

// GetDashboardStats retrieves dashboard statistics
func (r *dashboardRepository) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total dashboards
	var total int64
	if err := r.db.Model(&models.Dashboard{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Public dashboards
	var public int64
	if err := r.db.Model(&models.Dashboard{}).Where("is_public = ?", true).Count(&public).Error; err != nil {
		return nil, err
	}
	stats["public"] = public

	// Private dashboards
	stats["private"] = total - public

	// Recently created dashboards (last 7 days)
	sevenDaysAgo := time.Now().UTC().Add(-7 * 24 * time.Hour)
	var recentlyCreated int64
	if err := r.db.Model(&models.Dashboard{}).Where("created_at >= ?", sevenDaysAgo).Count(&recentlyCreated).Error; err != nil {
		return nil, err
	}
	stats["recently_created_7d"] = recentlyCreated

	// Recently accessed dashboards (last 24 hours)
	twentyFourHoursAgo := time.Now().UTC().Add(-24 * time.Hour)
	var recentlyAccessed int64
	if err := r.db.Model(&models.Dashboard{}).Where("last_accessed_at >= ?", twentyFourHoursAgo).Count(&recentlyAccessed).Error; err != nil {
		return nil, err
	}
	stats["recently_accessed_24h"] = recentlyAccessed

	// Count by creator (top 10)
	var creatorStats []struct {
		CreatedBy string `gorm:"column:created_by"`
		Count     int64  `gorm:"column:count"`
	}

	if err := r.db.Model(&models.Dashboard{}).Select("created_by, COUNT(*) as count").Group("created_by").Order("count DESC").Limit(10).Scan(&creatorStats).Error; err != nil {
		return nil, err
	}

	creatorMap := make(map[string]int64)
	for _, stat := range creatorStats {
		creatorMap[stat.CreatedBy] = stat.Count
	}
	stats["by_creator"] = creatorMap

	return stats, nil
}

// GetMostAccessedDashboards retrieves the most frequently accessed dashboards
func (r *dashboardRepository) GetMostAccessedDashboards(limit int, startTime, endTime *time.Time) ([]*models.Dashboard, error) {
	var dashboards []*models.Dashboard

	query := r.db.Model(&models.Dashboard{})

	if startTime != nil {
		query = query.Where("last_accessed_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("last_accessed_at <= ?", *endTime)
	}

	// Note: This is a simplified implementation
	// In a real scenario, you might want to track access counts separately
	err := query.Where("last_accessed_at IS NOT NULL").Order("last_accessed_at DESC").Limit(limit).Find(&dashboards).Error

	return dashboards, err
}

// GetDashboardCreationTrend retrieves dashboard creation trend data
func (r *dashboardRepository) GetDashboardCreationTrend(startTime, endTime *time.Time, interval string) ([]map[string]interface{}, error) {
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

	query := r.db.Model(&models.Dashboard{}).Select(dateTrunc + " as time_bucket, COUNT(*) as count")

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

// CloneDashboard creates a copy of an existing dashboard
func (r *dashboardRepository) CloneDashboard(originalID uuid.UUID, newName, newDescription, createdBy string) (*models.Dashboard, error) {
	// Get the original dashboard
	original, err := r.GetByID(originalID)
	if err != nil {
		return nil, err
	}

	// Create a new dashboard with copied configuration
	newDashboard := &models.Dashboard{
		ID:          uuid.New(),
		Name:        newName,
		Description: newDescription,
		Config:      original.Config, // Copy the configuration
		IsPublic:    false,           // New dashboard is private by default
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Save the new dashboard
	err = r.Create(newDashboard)
	if err != nil {
		return nil, err
	}

	return newDashboard, nil
}

// BulkUpdateDashboardVisibility updates the visibility of multiple dashboards
func (r *dashboardRepository) BulkUpdateDashboardVisibility(ids []uuid.UUID, isPublic bool) error {
	return r.db.Model(&models.Dashboard{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"is_public":  isPublic,
		"updated_at": time.Now().UTC(),
	}).Error
}

// GetDashboardsByWidget retrieves dashboards that contain a specific widget type
func (r *dashboardRepository) GetDashboardsByWidget(widgetType string, limit, offset int) ([]*models.Dashboard, int64, error) {
	var dashboards []*models.Dashboard
	var total int64

	// Note: This assumes widgets are stored in the config JSON field
	// The actual implementation might need to be adjusted based on the Dashboard model structure
	query := r.db.Model(&models.Dashboard{}).Where("config -> 'widgets' @> ?", `[{"type":"`+widgetType+`"}]`)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}

// GetDashboardAccessLog retrieves access log for a dashboard (if implemented)
func (r *dashboardRepository) GetDashboardAccessLog(dashboardID uuid.UUID, limit, offset int) ([]map[string]interface{}, int64, error) {
	// Note: This is a placeholder implementation
	// In a real scenario, you would have a separate access log table
	// For now, we'll return empty results
	return []map[string]interface{}{}, 0, nil
}

// UpdateDashboardConfig updates only the configuration of a dashboard
func (r *dashboardRepository) UpdateDashboardConfig(id uuid.UUID, config map[string]interface{}) error {
	return r.db.Model(&models.Dashboard{}).Where("id = ?", id).Updates(map[string]interface{}{
		"config":     config,
		"updated_at": time.Now().UTC(),
	}).Error
}

// GetDashboardConfigHistory retrieves configuration change history (if implemented)
func (r *dashboardRepository) GetDashboardConfigHistory(dashboardID uuid.UUID, limit, offset int) ([]map[string]interface{}, int64, error) {
	// Note: This is a placeholder implementation
	// In a real scenario, you would have a separate config history table
	// For now, we'll return empty results
	return []map[string]interface{}{}, 0, nil
}

// ValidateDashboardOwnership checks if a user owns a dashboard
func (r *dashboardRepository) ValidateDashboardOwnership(dashboardID uuid.UUID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Dashboard{}).Where("id = ? AND created_by = ?", dashboardID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetSharedDashboards retrieves dashboards shared with a specific user
func (r *dashboardRepository) GetSharedDashboards(userID string, limit, offset int) ([]*models.Dashboard, int64, error) {
	// Note: This is a simplified implementation
	// In a real scenario, you might have a separate sharing table
	// For now, we'll return public dashboards not created by the user
	var dashboards []*models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{}).Where("is_public = ? AND created_by != ?", true, userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get dashboards with pagination
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&dashboards).Error
	if err != nil {
		return nil, 0, err
	}

	return dashboards, total, nil
}