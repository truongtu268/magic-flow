package versioning

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// Service provides a unified interface for versioning operations
type Service struct {
	manager   *Manager
	migrator  *Migrator
	validator *Validator
	config    *VersioningConfig
}

// NewService creates a new versioning service
func NewService(repoManager database.RepositoryManager, config *VersioningConfig) *Service {
	if config == nil {
		config = &VersioningConfig{
			AutoVersioning:        false,
			VersioningStrategy:    VersioningStrategySemantic,
			MigrationTimeout:      30 * time.Minute,
			MaxRollbackDepth:      10,
			BackupBeforeMigration: true,
			ValidationRules:       []ValidationRule{},
			NotificationSettings: NotificationSettings{
				Enabled:             true,
				Channels:            []string{"email", "webhook"},
				OnVersionCreated:    true,
				OnVersionActivated:  true,
				OnMigrationStart:    true,
				OnMigrationComplete: true,
				OnMigrationFailed:   true,
				OnRollback:          true,
			},
			RetentionPolicy: RetentionPolicy{
				MaxVersions:        50,
				RetentionPeriod:    365 * 24 * time.Hour, // 1 year
				KeepActiveVersions: true,
				KeepTaggedVersions: true,
				ArchiveOldVersions: true,
			},
		}
	}

	return &Service{
		manager:   NewManager(repoManager),
		migrator:  NewMigrator(repoManager),
		validator: NewValidator(),
		config:    config,
	}
}

// CreateVersion creates a new version with validation and notifications
func (s *Service) CreateVersion(ctx context.Context, workflowID uuid.UUID, changes VersionChanges) (*models.WorkflowVersion, error) {
	// Validate the version before creation
	if err := s.manager.ValidateVersion(ctx, workflowID, changes); err != nil {
		return nil, fmt.Errorf("version validation failed: %w", err)
	}

	// Create the version
	version, err := s.manager.CreateVersion(ctx, workflowID, changes)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Send notification if enabled
	if s.config.NotificationSettings.Enabled && s.config.NotificationSettings.OnVersionCreated {
		go s.sendVersionCreatedNotification(ctx, version)
	}

	// Apply retention policy if auto-versioning is enabled
	if s.config.AutoVersioning {
		go s.applyRetentionPolicy(ctx, workflowID)
	}

	return version, nil
}

// ActivateVersion activates a version with migration and notifications
func (s *Service) ActivateVersion(ctx context.Context, versionID uuid.UUID) error {
	// Send migration start notification
	if s.config.NotificationSettings.Enabled && s.config.NotificationSettings.OnMigrationStart {
		go s.sendMigrationStartNotification(ctx, versionID)
	}

	// Activate the version
	err := s.manager.ActivateVersion(ctx, versionID)
	if err != nil {
		// Send migration failed notification
		if s.config.NotificationSettings.Enabled && s.config.NotificationSettings.OnMigrationFailed {
			go s.sendMigrationFailedNotification(ctx, versionID, err)
		}
		return fmt.Errorf("failed to activate version: %w", err)
	}

	// Send activation and migration complete notifications
	if s.config.NotificationSettings.Enabled {
		if s.config.NotificationSettings.OnVersionActivated {
			go s.sendVersionActivatedNotification(ctx, versionID)
		}
		if s.config.NotificationSettings.OnMigrationComplete {
			go s.sendMigrationCompleteNotification(ctx, versionID)
		}
	}

	return nil
}

// RollbackToVersion performs a rollback with validation and notifications
func (s *Service) RollbackToVersion(ctx context.Context, workflowID, targetVersionID uuid.UUID, reason string) error {
	// Validate rollback
	versionRepo := s.manager.repoManager.WorkflowVersionRepository()
	currentVersion, err := versionRepo.GetActiveVersion(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	targetVersion, err := versionRepo.GetByID(ctx, targetVersionID)
	if err != nil {
		return fmt.Errorf("failed to get target version: %w", err)
	}

	if err := s.validator.ValidateRollback(ctx, currentVersion, targetVersion); err != nil {
		return fmt.Errorf("rollback validation failed: %w", err)
	}

	// Perform rollback
	err = s.manager.RollbackToVersion(ctx, workflowID, targetVersionID, reason)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Send rollback notification
	if s.config.NotificationSettings.Enabled && s.config.NotificationSettings.OnRollback {
		go s.sendRollbackNotification(ctx, workflowID, targetVersionID, reason)
	}

	return nil
}

// GetVersionHistory returns version history with optional filtering
func (s *Service) GetVersionHistory(ctx context.Context, workflowID uuid.UUID, filters *VersionFilters) ([]*models.WorkflowVersion, error) {
	versions, err := s.manager.GetVersionHistory(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version history: %w", err)
	}

	// Apply filters if provided
	if filters != nil {
		versions = s.applyVersionFilters(versions, filters)
	}

	return versions, nil
}

// CompareVersions compares two versions with detailed analysis
func (s *Service) CompareVersions(ctx context.Context, version1ID, version2ID uuid.UUID) (*VersionComparison, error) {
	comparison, err := s.manager.CompareVersions(ctx, version1ID, version2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to compare versions: %w", err)
	}

	// Enhance comparison with additional analysis
	comparison.Summary = s.generateComparisonSummary(comparison.Differences)

	return comparison, nil
}

// GetCompatibilityMatrix returns compatibility matrix with recommendations
func (s *Service) GetCompatibilityMatrix(ctx context.Context, workflowID uuid.UUID) (*CompatibilityMatrix, error) {
	matrix, err := s.manager.GetCompatibilityMatrix(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get compatibility matrix: %w", err)
	}

	return matrix, nil
}

// CreateMigrationPlan creates a detailed migration plan with risk assessment
func (s *Service) CreateMigrationPlan(ctx context.Context, fromVersionID, toVersionID uuid.UUID) (*MigrationPlan, error) {
	plan, err := s.manager.GetMigrationPlan(ctx, fromVersionID, toVersionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration plan: %w", err)
	}

	// Validate the migration plan
	if err := s.validator.ValidateMigrationPlan(ctx, plan); err != nil {
		return nil, fmt.Errorf("migration plan validation failed: %w", err)
	}

	return plan, nil
}

// GetVersionMetrics returns comprehensive version metrics
func (s *Service) GetVersionMetrics(ctx context.Context, workflowID uuid.UUID, timeRange *TimeRange) (*VersionMetrics, error) {
	versions, err := s.manager.GetVersionHistory(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version history: %w", err)
	}

	// Filter versions by time range if provided
	if timeRange != nil {
		versions = s.filterVersionsByTimeRange(versions, timeRange)
	}

	// Calculate comprehensive metrics
	metrics := s.calculateDetailedMetrics(workflowID, versions)

	return metrics, nil
}

// TagVersion adds a tag to a version
func (s *Service) TagVersion(ctx context.Context, versionID uuid.UUID, tag *VersionTag) error {
	// Validate tag
	if err := s.validateVersionTag(tag); err != nil {
		return fmt.Errorf("invalid version tag: %w", err)
	}

	// Save tag (this would use a repository)
	// For now, this is a placeholder
	return nil
}

// GetVersionTags returns all tags for a version
func (s *Service) GetVersionTags(ctx context.Context, versionID uuid.UUID) ([]*VersionTag, error) {
	// Get tags from repository
	// For now, return empty slice
	return []*VersionTag{}, nil
}

// ArchiveOldVersions archives old versions based on retention policy
func (s *Service) ArchiveOldVersions(ctx context.Context, workflowID uuid.UUID) error {
	versions, err := s.manager.GetVersionHistory(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get version history: %w", err)
	}

	// Apply retention policy
	versionsToArchive := s.identifyVersionsToArchive(versions)

	// Archive versions
	for _, version := range versionsToArchive {
		if err := s.archiveVersion(ctx, version); err != nil {
			return fmt.Errorf("failed to archive version %s: %w", version.Version, err)
		}
	}

	return nil
}

// GetServiceHealth returns the health status of the versioning service
func (s *Service) GetServiceHealth(ctx context.Context) *ServiceHealth {
	health := &ServiceHealth{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks: map[string]interface{}{
			"database":  "connected",
			"validator": "operational",
			"migrator":  "operational",
		},
	}

	// Perform health checks
	// This would typically check database connectivity, etc.

	return health
}

// Private helper methods

func (s *Service) applyVersionFilters(versions []*models.WorkflowVersion, filters *VersionFilters) []*models.WorkflowVersion {
	filteredVersions := make([]*models.WorkflowVersion, 0)

	for _, version := range versions {
		// Apply change type filter
		if filters.ChangeType != "" && version.ChangeType != string(filters.ChangeType) {
			continue
		}

		// Apply active status filter
		if filters.ActiveOnly && !version.IsActive {
			continue
		}

		// Apply date range filter
		if filters.FromDate != nil && version.CreatedAt.Before(*filters.FromDate) {
			continue
		}
		if filters.ToDate != nil && version.CreatedAt.After(*filters.ToDate) {
			continue
		}

		filteredVersions = append(filteredVersions, version)
	}

	return filteredVersions
}

func (s *Service) generateComparisonSummary(differences []VersionDifference) ComparisonSummary {
	summary := ComparisonSummary{
		TotalDifferences: len(differences),
		ByType:           make(map[DifferenceType]int),
		ByImpact:         make(map[ImpactLevel]int),
		Compatibility:    CompatibilityLevelFull,
		Recommendations:  []string{},
	}

	// Count differences by type and impact
	for _, diff := range differences {
		summary.ByType[diff.Type]++
		summary.ByImpact[diff.Impact]++

		// Determine overall compatibility
		if diff.Impact == ImpactLevelHigh || diff.Impact == ImpactLevelCritical {
			summary.Compatibility = CompatibilityLevelNone
		} else if diff.Impact == ImpactLevelMedium && summary.Compatibility == CompatibilityLevelFull {
			summary.Compatibility = CompatibilityLevelPartial
		}
	}

	// Generate recommendations
	if summary.Compatibility == CompatibilityLevelNone {
		summary.Recommendations = append(summary.Recommendations, "Major version bump required due to breaking changes")
	} else if summary.Compatibility == CompatibilityLevelPartial {
		summary.Recommendations = append(summary.Recommendations, "Minor version bump recommended")
	} else {
		summary.Recommendations = append(summary.Recommendations, "Patch version bump is sufficient")
	}

	return summary
}

func (s *Service) filterVersionsByTimeRange(versions []*models.WorkflowVersion, timeRange *TimeRange) []*models.WorkflowVersion {
	filteredVersions := make([]*models.WorkflowVersion, 0)

	for _, version := range versions {
		if version.CreatedAt.After(timeRange.Start) && version.CreatedAt.Before(timeRange.End) {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func (s *Service) calculateDetailedMetrics(workflowID uuid.UUID, versions []*models.WorkflowVersion) *VersionMetrics {
	metrics := &VersionMetrics{
		WorkflowID:       workflowID,
		TotalVersions:    len(versions),
		VersionFrequency: make(map[string]int),
	}

	if len(versions) == 0 {
		return metrics
	}

	// Calculate metrics from versions
	for _, version := range versions {
		if version.IsActive {
			metrics.ActiveVersion = version.Version
		}
		if version.CreatedAt.After(metrics.LastVersionDate) {
			metrics.LastVersionDate = version.CreatedAt
		}
		metrics.VersionFrequency[version.ChangeType]++
	}

	// Calculate success rate and average migration time
	// These would typically come from migration execution records
	metrics.SuccessRate = 0.95 // Placeholder
	metrics.AverageMigrationTime = 5 * time.Minute // Placeholder

	return metrics
}

func (s *Service) validateVersionTag(tag *VersionTag) error {
	if tag.Name == "" {
		return fmt.Errorf("tag name is required")
	}
	if len(tag.Name) > 50 {
		return fmt.Errorf("tag name too long (max 50 characters)")
	}
	return nil
}

func (s *Service) identifyVersionsToArchive(versions []*models.WorkflowVersion) []*models.WorkflowVersion {
	versionsToArchive := make([]*models.WorkflowVersion, 0)

	// Apply retention policy logic
	if len(versions) > s.config.RetentionPolicy.MaxVersions {
		// Archive oldest versions beyond the limit
		excessCount := len(versions) - s.config.RetentionPolicy.MaxVersions
		for i := 0; i < excessCount; i++ {
			version := versions[i]
			// Don't archive active versions if policy says to keep them
			if s.config.RetentionPolicy.KeepActiveVersions && version.IsActive {
				continue
			}
			versionsToArchive = append(versionsToArchive, version)
		}
	}

	// Archive versions older than retention period
	cutoffDate := time.Now().Add(-s.config.RetentionPolicy.RetentionPeriod)
	for _, version := range versions {
		if version.CreatedAt.Before(cutoffDate) {
			if s.config.RetentionPolicy.KeepActiveVersions && version.IsActive {
				continue
			}
			versionsToArchive = append(versionsToArchive, version)
		}
	}

	return versionsToArchive
}

func (s *Service) archiveVersion(ctx context.Context, version *models.WorkflowVersion) error {
	// Archive version logic
	// This would typically move the version to an archive storage
	return nil
}

func (s *Service) applyRetentionPolicy(ctx context.Context, workflowID uuid.UUID) {
	// Apply retention policy in background
	if err := s.ArchiveOldVersions(ctx, workflowID); err != nil {
		// Log error but don't fail the main operation
		// In a real implementation, this would use a proper logger
	}
}

// Notification methods (simplified implementations)

func (s *Service) sendVersionCreatedNotification(ctx context.Context, version *models.WorkflowVersion) {
	// Send notification about version creation
	// This would integrate with notification service
}

func (s *Service) sendVersionActivatedNotification(ctx context.Context, versionID uuid.UUID) {
	// Send notification about version activation
}

func (s *Service) sendMigrationStartNotification(ctx context.Context, versionID uuid.UUID) {
	// Send notification about migration start
}

func (s *Service) sendMigrationCompleteNotification(ctx context.Context, versionID uuid.UUID) {
	// Send notification about migration completion
}

func (s *Service) sendMigrationFailedNotification(ctx context.Context, versionID uuid.UUID, err error) {
	// Send notification about migration failure
}

func (s *Service) sendRollbackNotification(ctx context.Context, workflowID, targetVersionID uuid.UUID, reason string) {
	// Send notification about rollback
}

// Additional types for service functionality

type VersionFilters struct {
	ChangeType ChangeType `json:"change_type,omitempty"`
	ActiveOnly bool       `json:"active_only,omitempty"`
	FromDate   *time.Time `json:"from_date,omitempty"`
	ToDate     *time.Time `json:"to_date,omitempty"`
}

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ServiceHealth struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
}