package versioning

import (
	"time"

	"github.com/google/uuid"

	"magic-flow/v2/pkg/models"
)

// ChangeType represents the type of change in a version
type ChangeType string

const (
	ChangeTypeMajor ChangeType = "major" // Breaking changes
	ChangeTypeMinor ChangeType = "minor" // New features, backward compatible
	ChangeTypePatch ChangeType = "patch" // Bug fixes, backward compatible
)

// VersionChanges represents the changes being made in a new version
type VersionChanges struct {
	ChangeType      ChangeType             `json:"change_type"`
	Summary         string                 `json:"summary"`
	Details         string                 `json:"details"`
	NewDefinition   map[string]interface{} `json:"new_definition"`
	CreatedBy       uuid.UUID              `json:"created_by"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// MigrationPlan represents a plan for migrating between versions
type MigrationPlan struct {
	ID              uuid.UUID              `json:"id"`
	FromVersionID   *uuid.UUID             `json:"from_version_id,omitempty"`
	ToVersionID     uuid.UUID              `json:"to_version_id"`
	MigrationSteps  []MigrationStep        `json:"migration_steps"`
	RollbackSteps   []MigrationStep        `json:"rollback_steps"`
	EstimatedTime   time.Duration          `json:"estimated_time"`
	RiskLevel       RiskLevel              `json:"risk_level"`
	Prerequisites   []string               `json:"prerequisites,omitempty"`
	Validations     []ValidationRule       `json:"validations,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// MigrationStep represents a single step in a migration
type MigrationStep struct {
	ID          uuid.UUID              `json:"id"`
	Order       int                    `json:"order"`
	Type        MigrationStepType      `json:"type"`
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Conditions  []string               `json:"conditions,omitempty"`
	Rollback    *MigrationStep         `json:"rollback,omitempty"`
	Timeout     time.Duration          `json:"timeout"`
	RetryCount  int                    `json:"retry_count"`
}

// MigrationStepType represents the type of migration step
type MigrationStepType string

const (
	MigrationStepTypeSchemaUpdate    MigrationStepType = "schema_update"
	MigrationStepTypeDataMigration   MigrationStepType = "data_migration"
	MigrationStepTypeConfigUpdate    MigrationStepType = "config_update"
	MigrationStepTypeValidation      MigrationStepType = "validation"
	MigrationStepTypeCleanup         MigrationStepType = "cleanup"
	MigrationStepTypeBackup          MigrationStepType = "backup"
	MigrationStepTypeNotification    MigrationStepType = "notification"
	MigrationStepTypeCustom          MigrationStepType = "custom"
)

// RiskLevel represents the risk level of a migration
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// ValidationRule represents a validation rule for migrations
type ValidationRule struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        ValidationType         `json:"type"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Required    bool                   `json:"required"`
}

// ValidationType represents the type of validation
type ValidationType string

const (
	ValidationTypeSchema       ValidationType = "schema"
	ValidationTypeData         ValidationType = "data"
	ValidationTypeCompatibility ValidationType = "compatibility"
	ValidationTypePerformance  ValidationType = "performance"
	ValidationTypeSecurity     ValidationType = "security"
	ValidationTypeCustom       ValidationType = "custom"
)

// VersionComparison represents a comparison between two versions
type VersionComparison struct {
	Version1    *models.WorkflowVersion `json:"version1"`
	Version2    *models.WorkflowVersion `json:"version2"`
	Differences []VersionDifference     `json:"differences"`
	Summary     ComparisonSummary       `json:"summary"`
	GeneratedAt time.Time               `json:"generated_at"`
}

// VersionDifference represents a difference between two versions
type VersionDifference struct {
	Type        DifferenceType `json:"type"`
	Path        string         `json:"path"`
	Description string         `json:"description"`
	OldValue    interface{}    `json:"old_value,omitempty"`
	NewValue    interface{}    `json:"new_value,omitempty"`
	Impact      ImpactLevel    `json:"impact"`
}

// DifferenceType represents the type of difference
type DifferenceType string

const (
	DifferenceTypeAdded    DifferenceType = "added"
	DifferenceTypeRemoved  DifferenceType = "removed"
	DifferenceTypeModified DifferenceType = "modified"
	DifferenceTypeMoved    DifferenceType = "moved"
)

// ImpactLevel represents the impact level of a difference
type ImpactLevel string

const (
	ImpactLevelLow      ImpactLevel = "low"
	ImpactLevelMedium   ImpactLevel = "medium"
	ImpactLevelHigh     ImpactLevel = "high"
	ImpactLevelCritical ImpactLevel = "critical"
)

// ComparisonSummary provides a summary of the comparison
type ComparisonSummary struct {
	TotalDifferences int                        `json:"total_differences"`
	ByType           map[DifferenceType]int     `json:"by_type"`
	ByImpact         map[ImpactLevel]int        `json:"by_impact"`
	Compatibility    CompatibilityLevel         `json:"compatibility"`
	Recommendations  []string                   `json:"recommendations,omitempty"`
}

// CompatibilityLevel represents the compatibility level between versions
type CompatibilityLevel string

const (
	CompatibilityLevelFull    CompatibilityLevel = "full"    // Fully compatible
	CompatibilityLevelPartial CompatibilityLevel = "partial" // Partially compatible
	CompatibilityLevelNone    CompatibilityLevel = "none"    // Not compatible
)

// CompatibilityMatrix represents compatibility between all versions
type CompatibilityMatrix struct {
	WorkflowID  uuid.UUID                              `json:"workflow_id"`
	Versions    []*models.WorkflowVersion              `json:"versions"`
	Matrix      map[string]map[string]CompatibilityLevel `json:"matrix"`
	GeneratedAt time.Time                              `json:"generated_at"`
}

// RollbackRecord represents a record of a rollback operation
type RollbackRecord struct {
	ID            uuid.UUID      `json:"id"`
	WorkflowID    uuid.UUID      `json:"workflow_id"`
	FromVersionID uuid.UUID      `json:"from_version_id"`
	ToVersionID   uuid.UUID      `json:"to_version_id"`
	Reason        string         `json:"reason"`
	ExecutedBy    uuid.UUID      `json:"executed_by"`
	ExecutedAt    time.Time      `json:"executed_at"`
	MigrationPlan *MigrationPlan `json:"migration_plan"`
	Status        RollbackStatus `json:"status"`
	Error         string         `json:"error,omitempty"`
	Duration      time.Duration  `json:"duration"`
}

// RollbackStatus represents the status of a rollback operation
type RollbackStatus string

const (
	RollbackStatusPending   RollbackStatus = "pending"
	RollbackStatusRunning   RollbackStatus = "running"
	RollbackStatusCompleted RollbackStatus = "completed"
	RollbackStatusFailed    RollbackStatus = "failed"
	RollbackStatusCancelled RollbackStatus = "cancelled"
)

// MigrationExecution represents the execution of a migration
type MigrationExecution struct {
	ID            uuid.UUID         `json:"id"`
	PlanID        uuid.UUID         `json:"plan_id"`
	WorkflowID    uuid.UUID         `json:"workflow_id"`
	FromVersionID *uuid.UUID        `json:"from_version_id,omitempty"`
	ToVersionID   uuid.UUID         `json:"to_version_id"`
	Status        MigrationStatus   `json:"status"`
	StartedAt     time.Time         `json:"started_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	Duration      time.Duration     `json:"duration"`
	StepResults   []StepResult      `json:"step_results"`
	Error         string            `json:"error,omitempty"`
	ExecutedBy    uuid.UUID         `json:"executed_by"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// MigrationStatus represents the status of a migration execution
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"
	MigrationStatusRunning   MigrationStatus = "running"
	MigrationStatusCompleted MigrationStatus = "completed"
	MigrationStatusFailed    MigrationStatus = "failed"
	MigrationStatusCancelled MigrationStatus = "cancelled"
	MigrationStatusRolledBack MigrationStatus = "rolled_back"
)

// StepResult represents the result of executing a migration step
type StepResult struct {
	StepID      uuid.UUID     `json:"step_id"`
	Status      StepStatus    `json:"status"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Duration    time.Duration `json:"duration"`
	Output      string        `json:"output,omitempty"`
	Error       string        `json:"error,omitempty"`
	RetryCount  int           `json:"retry_count"`
}

// StepStatus represents the status of a migration step
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
	StepStatusRetrying  StepStatus = "retrying"
)

// VersioningConfig represents configuration for the versioning system
type VersioningConfig struct {
	AutoVersioning       bool                   `json:"auto_versioning"`
	VersioningStrategy   VersioningStrategy     `json:"versioning_strategy"`
	MigrationTimeout     time.Duration          `json:"migration_timeout"`
	MaxRollbackDepth     int                    `json:"max_rollback_depth"`
	BackupBeforeMigration bool                   `json:"backup_before_migration"`
	ValidationRules      []ValidationRule       `json:"validation_rules"`
	NotificationSettings NotificationSettings   `json:"notification_settings"`
	RetentionPolicy      RetentionPolicy        `json:"retention_policy"`
}

// VersioningStrategy represents the strategy for versioning
type VersioningStrategy string

const (
	VersioningStrategySemantic   VersioningStrategy = "semantic"   // Semantic versioning (x.y.z)
	VersioningStrategyTimestamp  VersioningStrategy = "timestamp"  // Timestamp-based versioning
	VersioningStrategySequential VersioningStrategy = "sequential" // Sequential numbering
	VersioningStrategyCustom     VersioningStrategy = "custom"     // Custom versioning scheme
)

// NotificationSettings represents notification settings for versioning events
type NotificationSettings struct {
	Enabled           bool     `json:"enabled"`
	Channels          []string `json:"channels"`
	OnVersionCreated  bool     `json:"on_version_created"`
	OnVersionActivated bool     `json:"on_version_activated"`
	OnMigrationStart  bool     `json:"on_migration_start"`
	OnMigrationComplete bool     `json:"on_migration_complete"`
	OnMigrationFailed bool     `json:"on_migration_failed"`
	OnRollback        bool     `json:"on_rollback"`
}

// RetentionPolicy represents the retention policy for versions
type RetentionPolicy struct {
	MaxVersions       int           `json:"max_versions"`
	RetentionPeriod   time.Duration `json:"retention_period"`
	KeepActiveVersions bool          `json:"keep_active_versions"`
	KeepTaggedVersions bool          `json:"keep_tagged_versions"`
	ArchiveOldVersions bool          `json:"archive_old_versions"`
}

// VersionTag represents a tag applied to a version
type VersionTag struct {
	ID          uuid.UUID `json:"id"`
	VersionID   uuid.UUID `json:"version_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// VersionMetrics represents metrics for version management
type VersionMetrics struct {
	WorkflowID        uuid.UUID `json:"workflow_id"`
	TotalVersions     int       `json:"total_versions"`
	ActiveVersion     string    `json:"active_version"`
	LastVersionDate   time.Time `json:"last_version_date"`
	MigrationCount    int       `json:"migration_count"`
	RollbackCount     int       `json:"rollback_count"`
	AverageMigrationTime time.Duration `json:"average_migration_time"`
	SuccessRate       float64   `json:"success_rate"`
	VersionFrequency  map[string]int `json:"version_frequency"`
}