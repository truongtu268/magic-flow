package versioning

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// Manager handles workflow versioning operations
type Manager struct {
	repoManager database.RepositoryManager
	migrator    *Migrator
	validator   *Validator
}

// NewManager creates a new versioning manager
func NewManager(repoManager database.RepositoryManager) *Manager {
	return &Manager{
		repoManager: repoManager,
		migrator:    NewMigrator(repoManager),
		validator:   NewValidator(),
	}
}

// CreateVersion creates a new version of a workflow
func (m *Manager) CreateVersion(ctx context.Context, workflowID uuid.UUID, changes VersionChanges) (*models.WorkflowVersion, error) {
	workflowRepo := m.repoManager.WorkflowRepository()
	versionRepo := m.repoManager.WorkflowVersionRepository()

	// Get current workflow
	workflow, err := workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Get current version
	currentVersion, err := versionRepo.GetLatestVersion(ctx, workflowID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	// Determine next version number
	nextVersionNumber := "1.0.0"
	if currentVersion != nil {
		nextVersionNumber = m.calculateNextVersion(currentVersion.Version, changes.ChangeType)
	}

	// Validate the new version
	if err := m.validator.ValidateVersion(ctx, workflow, changes); err != nil {
		return nil, fmt.Errorf("version validation failed: %w", err)
	}

	// Create new version
	newVersion := &models.WorkflowVersion{
		WorkflowID:    workflowID,
		Version:       nextVersionNumber,
		Definition:    changes.NewDefinition,
		ChangeType:    string(changes.ChangeType),
		ChangeSummary: changes.Summary,
		ChangeDetails: changes.Details,
		CreatedBy:     changes.CreatedBy,
		IsActive:      false, // Will be activated after successful migration
		Metadata:      changes.Metadata,
	}

	// Save the new version
	err = versionRepo.Create(ctx, newVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Create migration plan
	migrationPlan, err := m.migrator.CreateMigrationPlan(ctx, currentVersion, newVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration plan: %w", err)
	}

	newVersion.MigrationPlan = migrationPlan

	// Update version with migration plan
	err = versionRepo.Update(ctx, newVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to update version with migration plan: %w", err)
	}

	return newVersion, nil
}

// ActivateVersion activates a specific version of a workflow
func (m *Manager) ActivateVersion(ctx context.Context, versionID uuid.UUID) error {
	versionRepo := m.repoManager.WorkflowVersionRepository()
	workflowRepo := m.repoManager.WorkflowRepository()

	// Get the version to activate
	version, err := versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	// Get current active version
	currentVersion, err := versionRepo.GetActiveVersion(ctx, version.WorkflowID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get current active version: %w", err)
	}

	// Execute migration if needed
	if currentVersion != nil {
		err = m.migrator.ExecuteMigration(ctx, currentVersion, version)
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Start transaction
	tx := m.repoManager.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deactivate current version
	if currentVersion != nil {
		currentVersion.IsActive = false
		currentVersion.DeactivatedAt = &time.Time{}
		*currentVersion.DeactivatedAt = time.Now()
		err = versionRepo.UpdateWithTx(ctx, tx, currentVersion)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to deactivate current version: %w", err)
		}
	}

	// Activate new version
	version.IsActive = true
	version.ActivatedAt = &time.Time{}
	*version.ActivatedAt = time.Now()
	err = versionRepo.UpdateWithTx(ctx, tx, version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to activate new version: %w", err)
	}

	// Update workflow with new definition
	workflow, err := workflowRepo.GetByID(ctx, version.WorkflowID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	workflow.Definition = version.Definition
	workflow.Version = version.Version
	workflow.UpdatedAt = time.Now()

	err = workflowRepo.UpdateWithTx(ctx, tx, workflow)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackToVersion rolls back a workflow to a previous version
func (m *Manager) RollbackToVersion(ctx context.Context, workflowID uuid.UUID, targetVersionID uuid.UUID, reason string) error {
	versionRepo := m.repoManager.WorkflowVersionRepository()

	// Get target version
	targetVersion, err := versionRepo.GetByID(ctx, targetVersionID)
	if err != nil {
		return fmt.Errorf("failed to get target version: %w", err)
	}

	// Verify target version belongs to the workflow
	if targetVersion.WorkflowID != workflowID {
		return fmt.Errorf("target version does not belong to the specified workflow")
	}

	// Get current active version
	currentVersion, err := versionRepo.GetActiveVersion(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get current active version: %w", err)
	}

	// Check if rollback is possible
	if !m.canRollback(currentVersion, targetVersion) {
		return fmt.Errorf("rollback to version %s is not possible", targetVersion.Version)
	}

	// Create rollback migration plan
	rollbackPlan, err := m.migrator.CreateRollbackPlan(ctx, currentVersion, targetVersion)
	if err != nil {
		return fmt.Errorf("failed to create rollback plan: %w", err)
	}

	// Execute rollback migration
	err = m.migrator.ExecuteRollback(ctx, rollbackPlan)
	if err != nil {
		return fmt.Errorf("rollback execution failed: %w", err)
	}

	// Create rollback record
	rollbackRecord := &RollbackRecord{
		ID:              uuid.New(),
		WorkflowID:      workflowID,
		FromVersionID:   currentVersion.ID,
		ToVersionID:     targetVersionID,
		Reason:          reason,
		ExecutedAt:      time.Now(),
		MigrationPlan:   rollbackPlan,
		Status:          RollbackStatusCompleted,
	}

	// Save rollback record
	err = m.saveRollbackRecord(ctx, rollbackRecord)
	if err != nil {
		return fmt.Errorf("failed to save rollback record: %w", err)
	}

	// Activate target version
	err = m.ActivateVersion(ctx, targetVersionID)
	if err != nil {
		return fmt.Errorf("failed to activate target version: %w", err)
	}

	return nil
}

// GetVersionHistory returns the version history for a workflow
func (m *Manager) GetVersionHistory(ctx context.Context, workflowID uuid.UUID) ([]*models.WorkflowVersion, error) {
	versionRepo := m.repoManager.WorkflowVersionRepository()
	return versionRepo.GetVersionHistory(ctx, workflowID)
}

// CompareVersions compares two versions and returns the differences
func (m *Manager) CompareVersions(ctx context.Context, version1ID, version2ID uuid.UUID) (*VersionComparison, error) {
	versionRepo := m.repoManager.WorkflowVersionRepository()

	// Get both versions
	version1, err := versionRepo.GetByID(ctx, version1ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version 1: %w", err)
	}

	version2, err := versionRepo.GetByID(ctx, version2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version 2: %w", err)
	}

	// Ensure versions belong to the same workflow
	if version1.WorkflowID != version2.WorkflowID {
		return nil, fmt.Errorf("versions belong to different workflows")
	}

	// Compare definitions
	comparison := &VersionComparison{
		Version1:    version1,
		Version2:    version2,
		Differences: m.calculateDifferences(version1.Definition, version2.Definition),
		GeneratedAt: time.Now(),
	}

	return comparison, nil
}

// GetCompatibilityMatrix returns compatibility information between versions
func (m *Manager) GetCompatibilityMatrix(ctx context.Context, workflowID uuid.UUID) (*CompatibilityMatrix, error) {
	versions, err := m.GetVersionHistory(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version history: %w", err)
	}

	matrix := &CompatibilityMatrix{
		WorkflowID:  workflowID,
		Versions:    versions,
		Matrix:      make(map[string]map[string]CompatibilityLevel),
		GeneratedAt: time.Now(),
	}

	// Calculate compatibility between all version pairs
	for i, v1 := range versions {
		matrix.Matrix[v1.Version] = make(map[string]CompatibilityLevel)
		for j, v2 := range versions {
			if i == j {
				matrix.Matrix[v1.Version][v2.Version] = CompatibilityLevelFull
			} else {
				matrix.Matrix[v1.Version][v2.Version] = m.calculateCompatibility(v1, v2)
			}
		}
	}

	return matrix, nil
}

// GetMigrationPlan returns the migration plan for upgrading between versions
func (m *Manager) GetMigrationPlan(ctx context.Context, fromVersionID, toVersionID uuid.UUID) (*MigrationPlan, error) {
	versionRepo := m.repoManager.WorkflowVersionRepository()

	// Get both versions
	fromVersion, err := versionRepo.GetByID(ctx, fromVersionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get from version: %w", err)
	}

	toVersion, err := versionRepo.GetByID(ctx, toVersionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get to version: %w", err)
	}

	return m.migrator.CreateMigrationPlan(ctx, fromVersion, toVersion)
}

// ValidateVersion validates a version before creation
func (m *Manager) ValidateVersion(ctx context.Context, workflowID uuid.UUID, changes VersionChanges) error {
	workflowRepo := m.repoManager.WorkflowRepository()

	workflow, err := workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	return m.validator.ValidateVersion(ctx, workflow, changes)
}

// Helper methods

func (m *Manager) calculateNextVersion(currentVersion string, changeType ChangeType) string {
	// Parse semantic version and increment based on change type
	// This is a simplified implementation
	switch changeType {
	case ChangeTypeMajor:
		return incrementMajorVersion(currentVersion)
	case ChangeTypeMinor:
		return incrementMinorVersion(currentVersion)
	case ChangeTypePatch:
		return incrementPatchVersion(currentVersion)
	default:
		return incrementPatchVersion(currentVersion)
	}
}

func (m *Manager) canRollback(currentVersion, targetVersion *models.WorkflowVersion) bool {
	// Check if rollback is possible based on compatibility and business rules
	// This is a simplified implementation
	return m.calculateCompatibility(currentVersion, targetVersion) != CompatibilityLevelNone
}

func (m *Manager) calculateDifferences(def1, def2 map[string]interface{}) []VersionDifference {
	// Calculate differences between two workflow definitions
	// This is a simplified implementation
	differences := []VersionDifference{}

	// Compare steps
	steps1, ok1 := def1["steps"].([]interface{})
	steps2, ok2 := def2["steps"].([]interface{})

	if ok1 && ok2 {
		if len(steps1) != len(steps2) {
			differences = append(differences, VersionDifference{
				Type:        DifferenceTypeModified,
				Path:        "steps",
				Description: fmt.Sprintf("Step count changed from %d to %d", len(steps1), len(steps2)),
				OldValue:    len(steps1),
				NewValue:    len(steps2),
			})
		}
	}

	return differences
}

func (m *Manager) calculateCompatibility(v1, v2 *models.WorkflowVersion) CompatibilityLevel {
	// Calculate compatibility level between two versions
	// This is a simplified implementation
	if v1.ChangeType == string(ChangeTypeMajor) || v2.ChangeType == string(ChangeTypeMajor) {
		return CompatibilityLevelNone
	}
	if v1.ChangeType == string(ChangeTypeMinor) || v2.ChangeType == string(ChangeTypeMinor) {
		return CompatibilityLevelPartial
	}
	return CompatibilityLevelFull
}

func (m *Manager) saveRollbackRecord(ctx context.Context, record *RollbackRecord) error {
	// Save rollback record to database
	// This would typically use a dedicated repository
	return nil
}

// Version increment helpers

func incrementMajorVersion(version string) string {
	// Simplified version increment
	return "2.0.0" // This should parse and increment properly
}

func incrementMinorVersion(version string) string {
	// Simplified version increment
	return "1.1.0" // This should parse and increment properly
}

func incrementPatchVersion(version string) string {
	// Simplified version increment
	return "1.0.1" // This should parse and increment properly
}