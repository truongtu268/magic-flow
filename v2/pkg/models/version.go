package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VersionStatus represents the status of a workflow version
type VersionStatus string

const (
	VersionStatusDevelopment VersionStatus = "development"
	VersionStatusTesting     VersionStatus = "testing"
	VersionStatusStaging     VersionStatus = "staging"
	VersionStatusProduction  VersionStatus = "production"
	VersionStatusDeprecated  VersionStatus = "deprecated"
	VersionStatusArchived    VersionStatus = "archived"
)

// MigrationStrategy represents the migration strategy
type MigrationStrategy string

const (
	MigrationStrategyBlueGreen MigrationStrategy = "blue_green"
	MigrationStrategyCanary    MigrationStrategy = "canary"
	MigrationStrategyRolling   MigrationStrategy = "rolling"
	MigrationStrategyImmediate MigrationStrategy = "immediate"
)

// WorkflowVersion represents a specific version of a workflow
type WorkflowVersion struct {
	ID         uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkflowID uuid.UUID     `json:"workflow_id" gorm:"type:uuid;not null;index"`
	Version    string        `json:"version" gorm:"not null;index"`
	Status     VersionStatus `json:"status" gorm:"default:'development'"`
	
	// Version metadata
	Description     string    `json:"description"`
	Changelog       string    `json:"changelog"`
	BreakingChanges bool      `json:"breaking_changes"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	
	// Workflow definition for this version
	Definition WorkflowDefinition `json:"definition" gorm:"type:jsonb"`
	
	// Schema definitions
	InputSchema  JSONSchema `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema JSONSchema `json:"output_schema" gorm:"type:jsonb"`
	
	// Version-specific configuration
	Config VersionConfig `json:"config" gorm:"type:jsonb"`
	
	// Migration information
	Migration MigrationInfo `json:"migration" gorm:"type:jsonb"`
	
	// Compatibility information
	Compatibility CompatibilityInfo `json:"compatibility" gorm:"type:jsonb"`
	
	// Dependencies
	Dependencies []Dependency `json:"dependencies" gorm:"type:jsonb"`
	
	// Rollback information
	Rollback RollbackInfo `json:"rollback" gorm:"type:jsonb"`
	
	// Timestamps
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Workflow   Workflow     `json:"-" gorm:"foreignKey:WorkflowID"`
	Executions []Execution  `json:"-" gorm:"foreignKey:WorkflowVersionID"`
	Deployments []Deployment `json:"-" gorm:"foreignKey:VersionID"`
}

// VersionConfig represents version-specific configuration
type VersionConfig struct {
	Timeout         string            `json:"timeout,omitempty"`
	MaxConcurrency  int               `json:"max_concurrency,omitempty"`
	RetryPolicy     RetryPolicy       `json:"retry_policy,omitempty"`
	ErrorHandling   ErrorHandling     `json:"error_handling,omitempty"`
	Environment     map[string]string `json:"environment,omitempty"`
	FeatureFlags    map[string]bool   `json:"feature_flags,omitempty"`
	ResourceLimits  ResourceLimits    `json:"resource_limits,omitempty"`
}

// ResourceLimits represents resource limits for the version
type ResourceLimits struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	Disk   string `json:"disk,omitempty"`
}

// MigrationInfo represents migration information
type MigrationInfo struct {
	Required          bool              `json:"required"`
	Strategy          MigrationStrategy `json:"strategy"`
	FromVersion       string            `json:"from_version,omitempty"`
	DataMigration     DataMigration     `json:"data_migration,omitempty"`
	TrafficRouting    TrafficRouting    `json:"traffic_routing,omitempty"`
	HealthChecks      []HealthCheck     `json:"health_checks,omitempty"`
	RollbackTriggers  []RollbackTrigger `json:"rollback_triggers,omitempty"`
	Validation        MigrationValidation `json:"validation,omitempty"`
	Notifications     []Notification    `json:"notifications,omitempty"`
}

// DataMigration represents data migration configuration
type DataMigration struct {
	Required        bool              `json:"required"`
	Scripts         []MigrationScript `json:"scripts,omitempty"`
	Validation      []ValidationRule  `json:"validation,omitempty"`
	Backup          BackupConfig      `json:"backup,omitempty"`
	Timeout         string            `json:"timeout,omitempty"`
	Parallel        bool              `json:"parallel"`
}

// MigrationScript represents a migration script
type MigrationScript struct {
	Name           string            `json:"name"`
	Type           string            `json:"type"` // sql, go, custom
	Script         string            `json:"script"`
	RollbackScript string            `json:"rollback_script,omitempty"`
	Timeout        string            `json:"timeout,omitempty"`
	DependsOn      []string          `json:"depends_on,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Name           string                 `json:"name"`
	Type           string                 `json:"type"` // sql, custom, function
	Query          string                 `json:"query,omitempty"`
	Function       string                 `json:"function,omitempty"`
	ExpectedResult interface{}            `json:"expected_result,omitempty"`
	Timeout        string                 `json:"timeout,omitempty"`
	Config         map[string]interface{} `json:"config,omitempty"`
}

// BackupConfig represents backup configuration
type BackupConfig struct {
	Enabled   bool   `json:"enabled"`
	Retention string `json:"retention"`
	Storage   string `json:"storage"`
	Encryption bool  `json:"encryption"`
	Compression bool `json:"compression"`
}

// TrafficRouting represents traffic routing configuration
type TrafficRouting struct {
	Phases []RoutingPhase `json:"phases"`
}

// RoutingPhase represents a phase in traffic routing
type RoutingPhase struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Duration    string            `json:"duration"`
	Traffic     map[string]int    `json:"traffic"` // version -> percentage
	Criteria    []RoutingCriteria `json:"criteria,omitempty"`
	Actions     []string          `json:"actions,omitempty"`
}

// RoutingCriteria represents routing criteria
type RoutingCriteria struct {
	Type   string `json:"type"`   // header, customer_tier, random
	Header string `json:"header,omitempty"`
	Value  string `json:"value,omitempty"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	Name         string `json:"name"`
	Type         string `json:"type"` // http, tcp, custom
	Path         string `json:"path,omitempty"`
	Port         int    `json:"port,omitempty"`
	InitialDelay string `json:"initial_delay"`
	Period       string `json:"period"`
	Timeout      string `json:"timeout"`
	FailureThreshold int `json:"failure_threshold"`
}

// RollbackTrigger represents a rollback trigger
type RollbackTrigger struct {
	Metric     string      `json:"metric"`
	Threshold  interface{} `json:"threshold"`
	Comparison string      `json:"comparison"` // greater_than, less_than, equals
	Duration   string      `json:"duration"`
	Action     string      `json:"action"` // rollback, alert
}

// MigrationValidation represents migration validation configuration
type MigrationValidation struct {
	PreMigration  []ValidationRule `json:"pre_migration,omitempty"`
	PostMigration []ValidationRule `json:"post_migration,omitempty"`
	Timeout       string           `json:"timeout,omitempty"`
}

// Deployment represents a deployment of a workflow version
type Deployment struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VersionID uuid.UUID        `json:"version_id" gorm:"type:uuid;not null;index"`
	Status    DeploymentStatus `json:"status" gorm:"default:'pending'"`
	
	// Deployment information
	Environment string                 `json:"environment"`
	Strategy    MigrationStrategy      `json:"strategy"`
	Config      map[string]interface{} `json:"config" gorm:"type:jsonb"`
	
	// Timing
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Duration    int64      `json:"duration"`
	
	// Error information
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	
	// Deployment metadata
	DeployedBy string                 `json:"deployed_by"`
	Metadata   map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Version WorkflowVersion `json:"-" gorm:"foreignKey:VersionID"`
}

// DeploymentStatus represents the status of a deployment
type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusDeploying  DeploymentStatus = "deploying"
	DeploymentStatusDeployed   DeploymentStatus = "deployed"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusRollingBack DeploymentStatus = "rolling_back"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

// VersionComparison represents a comparison between two versions
type VersionComparison struct {
	Workflow    string           `json:"workflow"`
	FromVersion string           `json:"from_version"`
	ToVersion   string           `json:"to_version"`
	Changes     VersionChanges   `json:"changes"`
	Compatibility CompatibilityCheck `json:"compatibility"`
}

// VersionChanges represents changes between versions
type VersionChanges struct {
	AddedSteps    []StepChange   `json:"added_steps,omitempty"`
	ModifiedSteps []StepChange   `json:"modified_steps,omitempty"`
	RemovedSteps  []StepChange   `json:"removed_steps,omitempty"`
	SchemaChanges SchemaChanges  `json:"schema_changes"`
	ConfigChanges ConfigChanges  `json:"config_changes"`
}

// StepChange represents a change to a workflow step
type StepChange struct {
	Name        string                   `json:"name"`
	Type        string                   `json:"type,omitempty"`
	Position    int                      `json:"position,omitempty"`
	Description string                   `json:"description,omitempty"`
	Changes     []FieldChange            `json:"changes,omitempty"`
}

// FieldChange represents a change to a field
type FieldChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

// SchemaChanges represents changes to input/output schemas
type SchemaChanges struct {
	InputSchema  SchemaFieldChanges `json:"input_schema"`
	OutputSchema SchemaFieldChanges `json:"output_schema"`
}

// SchemaFieldChanges represents changes to schema fields
type SchemaFieldChanges struct {
	AddedFields    []SchemaField `json:"added_fields,omitempty"`
	RemovedFields  []SchemaField `json:"removed_fields,omitempty"`
	ModifiedFields []SchemaField `json:"modified_fields,omitempty"`
}

// SchemaField represents a schema field
type SchemaField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ConfigChanges represents changes to configuration
type ConfigChanges struct {
	AddedConfig    map[string]interface{} `json:"added_config,omitempty"`
	RemovedConfig  map[string]interface{} `json:"removed_config,omitempty"`
	ModifiedConfig map[string]FieldChange `json:"modified_config,omitempty"`
}

// CompatibilityCheck represents compatibility check results
type CompatibilityCheck struct {
	BreakingChanges   bool     `json:"breaking_changes"`
	MigrationRequired bool     `json:"migration_required"`
	RollbackSafe      bool     `json:"rollback_safe"`
	Warnings          []string `json:"warnings,omitempty"`
	Errors            []string `json:"errors,omitempty"`
}

// BeforeCreate sets the ID before creating
func (wv *WorkflowVersion) BeforeCreate(tx *gorm.DB) error {
	if wv.ID == uuid.Nil {
		wv.ID = uuid.New()
	}
	return nil
}

// BeforeCreate sets the ID before creating
func (d *Deployment) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the WorkflowVersion model
func (WorkflowVersion) TableName() string {
	return "workflow_versions"
}

// TableName returns the table name for the Deployment model
func (Deployment) TableName() string {
	return "deployments"
}

// Validate validates the workflow version
func (wv *WorkflowVersion) Validate() error {
	if wv.Version == "" {
		return fmt.Errorf("version is required")
	}
	
	if wv.CreatedBy == "" {
		return fmt.Errorf("created_by is required")
	}
	
	// Validate workflow definition
	if len(wv.Definition.Spec.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	
	return nil
}

// IsProduction returns true if the version is in production
func (wv *WorkflowVersion) IsProduction() bool {
	return wv.Status == VersionStatusProduction
}

// IsDeprecated returns true if the version is deprecated
func (wv *WorkflowVersion) IsDeprecated() bool {
	return wv.Status == VersionStatusDeprecated
}

// IsArchived returns true if the version is archived
func (wv *WorkflowVersion) IsArchived() bool {
	return wv.Status == VersionStatusArchived
}

// CanDeploy returns true if the version can be deployed
func (wv *WorkflowVersion) CanDeploy() bool {
	return wv.Status == VersionStatusTesting || wv.Status == VersionStatusStaging
}

// CanRollback returns true if the version can be rolled back to
func (wv *WorkflowVersion) CanRollback() bool {
	return wv.Status == VersionStatusProduction && !wv.IsArchived()
}

// GetFullName returns the full name including version
func (wv *WorkflowVersion) GetFullName() string {
	return fmt.Sprintf("%s:%s", wv.Definition.Metadata.Name, wv.Version)
}

// ToJSON converts the workflow version to JSON
func (wv *WorkflowVersion) ToJSON() ([]byte, error) {
	return json.Marshal(wv)
}

// FromJSON populates the workflow version from JSON
func (wv *WorkflowVersion) FromJSON(data []byte) error {
	return json.Unmarshal(data, wv)
}

// Start marks the deployment as started
func (d *Deployment) Start() {
	now := time.Now()
	d.Status = DeploymentStatusDeploying
	d.StartedAt = &now
}

// Complete marks the deployment as completed
func (d *Deployment) Complete() {
	now := time.Now()
	d.Status = DeploymentStatusDeployed
	d.CompletedAt = &now
	
	if d.StartedAt != nil {
		d.Duration = now.Sub(*d.StartedAt).Milliseconds()
	}
}

// Fail marks the deployment as failed
func (d *Deployment) Fail(err error, errorCode string) {
	now := time.Now()
	d.Status = DeploymentStatusFailed
	d.CompletedAt = &now
	d.Error = err.Error()
	d.ErrorCode = errorCode
	
	if d.StartedAt != nil {
		d.Duration = now.Sub(*d.StartedAt).Milliseconds()
	}
}

// Rollback marks the deployment as rolling back
func (d *Deployment) Rollback() {
	d.Status = DeploymentStatusRollingBack
}

// RollbackComplete marks the rollback as completed
func (d *Deployment) RollbackComplete() {
	now := time.Now()
	d.Status = DeploymentStatusRolledBack
	d.CompletedAt = &now
	
	if d.StartedAt != nil {
		d.Duration = now.Sub(*d.StartedAt).Milliseconds()
	}
}

// IsDeploying returns true if the deployment is in progress
func (d *Deployment) IsDeploying() bool {
	return d.Status == DeploymentStatusDeploying
}

// IsDeployed returns true if the deployment is completed
func (d *Deployment) IsDeployed() bool {
	return d.Status == DeploymentStatusDeployed
}

// IsFailed returns true if the deployment failed
func (d *Deployment) IsFailed() bool {
	return d.Status == DeploymentStatusFailed
}

// IsRollingBack returns true if the deployment is rolling back
func (d *Deployment) IsRollingBack() bool {
	return d.Status == DeploymentStatusRollingBack
}

// IsRolledBack returns true if the deployment was rolled back
func (d *Deployment) IsRolledBack() bool {
	return d.Status == DeploymentStatusRolledBack
}

// IsFinished returns true if the deployment is in a terminal state
func (d *Deployment) IsFinished() bool {
	return d.IsDeployed() || d.IsFailed() || d.IsRolledBack()
}

// GetDurationSeconds returns the duration in seconds
func (d *Deployment) GetDurationSeconds() float64 {
	return float64(d.Duration) / 1000.0
}