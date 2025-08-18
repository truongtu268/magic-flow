package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusDraft      WorkflowStatus = "draft"
	WorkflowStatusActive     WorkflowStatus = "active"
	WorkflowStatusInactive   WorkflowStatus = "inactive"
	WorkflowStatusDeprecated WorkflowStatus = "deprecated"
	WorkflowStatusArchived   WorkflowStatus = "archived"
)

// Workflow represents a workflow definition
type Workflow struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null" validate:"required,min=1,max=255"`
	Description string         `json:"description" gorm:"type:text"`
	Version     string         `json:"version" gorm:"not null" validate:"required"`
	Status      WorkflowStatus `json:"status" gorm:"default:'draft'" validate:"required"`
	
	// Metadata
	Tags      []string `json:"tags" gorm:"type:text[]"`
	Owner     string   `json:"owner" validate:"required"`
	CreatedBy string   `json:"created_by" validate:"required"`
	
	// Workflow definition
	Definition WorkflowDefinition `json:"definition" gorm:"type:jsonb"`
	
	// Schema definitions
	InputSchema  JSONSchema `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema JSONSchema `json:"output_schema" gorm:"type:jsonb"`
	
	// Configuration
	Config WorkflowConfig `json:"config" gorm:"type:jsonb"`
	
	// Versioning
	VersionInfo VersionInfo `json:"version_info" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Executions []Execution `json:"-" gorm:"foreignKey:WorkflowID"`
	Versions   []WorkflowVersion `json:"-" gorm:"foreignKey:WorkflowID"`
}

// WorkflowDefinition represents the YAML workflow definition
type WorkflowDefinition struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   WorkflowMetadata       `json:"metadata" yaml:"metadata"`
	Spec       WorkflowSpec           `json:"spec" yaml:"spec"`
}

// WorkflowMetadata contains workflow metadata
type WorkflowMetadata struct {
	Name        string            `json:"name" yaml:"name"`
	Version     string            `json:"version" yaml:"version"`
	Description string            `json:"description" yaml:"description"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// WorkflowSpec contains the workflow specification
type WorkflowSpec struct {
	InputSchema  JSONSchema    `json:"input_schema" yaml:"input_schema"`
	OutputSchema JSONSchema    `json:"output_schema" yaml:"output_schema"`
	Steps        []WorkflowStep `json:"steps" yaml:"steps"`
	ErrorHandling ErrorHandling `json:"error_handling,omitempty" yaml:"error_handling,omitempty"`
	RetryPolicy   RetryPolicy   `json:"retry_policy,omitempty" yaml:"retry_policy,omitempty"`
	Timeout       string        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	FeatureFlags  map[string]bool `json:"feature_flags,omitempty" yaml:"feature_flags,omitempty"`
}

// WorkflowStep represents a single step in the workflow
type WorkflowStep struct {
	Name        string                 `json:"name" yaml:"name"`
	Type        string                 `json:"type" yaml:"type"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	DependsOn   []string               `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	Condition   string                 `json:"condition,omitempty" yaml:"condition,omitempty"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
	DataMapping DataMapping            `json:"data_mapping,omitempty" yaml:"data_mapping,omitempty"`
	ErrorHandling ErrorHandling        `json:"error_handling,omitempty" yaml:"error_handling,omitempty"`
	RetryPolicy RetryPolicy            `json:"retry_policy,omitempty" yaml:"retry_policy,omitempty"`
	Timeout     string                 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// DataMapping represents data transformation between steps
type DataMapping struct {
	Input  map[string]string `json:"input,omitempty" yaml:"input,omitempty"`
	Output map[string]string `json:"output,omitempty" yaml:"output,omitempty"`
}

// ErrorHandling represents error handling configuration
type ErrorHandling struct {
	Strategy    string `json:"strategy" yaml:"strategy"` // continue, stop, retry
	MaxRetries  int    `json:"max_retries,omitempty" yaml:"max_retries,omitempty"`
	RetryDelay  string `json:"retry_delay,omitempty" yaml:"retry_delay,omitempty"`
	FallbackStep string `json:"fallback_step,omitempty" yaml:"fallback_step,omitempty"`
}

// RetryPolicy represents retry configuration
type RetryPolicy struct {
	MaxAttempts int    `json:"max_attempts" yaml:"max_attempts"`
	Delay       string `json:"delay" yaml:"delay"`
	Backoff     string `json:"backoff,omitempty" yaml:"backoff,omitempty"` // linear, exponential
	MaxDelay    string `json:"max_delay,omitempty" yaml:"max_delay,omitempty"`
}

// JSONSchema represents a JSON schema definition
type JSONSchema struct {
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
	AdditionalProperties interface{} `json:"additionalProperties,omitempty"`
}

// WorkflowConfig represents workflow configuration
type WorkflowConfig struct {
	Timeout         string            `json:"timeout,omitempty"`
	MaxConcurrency  int               `json:"max_concurrency,omitempty"`
	RetryPolicy     RetryPolicy       `json:"retry_policy,omitempty"`
	ErrorHandling   ErrorHandling     `json:"error_handling,omitempty"`
	Notifications   []Notification    `json:"notifications,omitempty"`
	Webhooks        []Webhook         `json:"webhooks,omitempty"`
	Environment     map[string]string `json:"environment,omitempty"`
}

// Notification represents a notification configuration
type Notification struct {
	Type      string            `json:"type"` // email, slack, webhook
	Events    []string          `json:"events"`
	Config    map[string]string `json:"config"`
	Enabled   bool              `json:"enabled"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Events  []string          `json:"events"`
	Enabled bool              `json:"enabled"`
}

// VersionInfo represents version-specific information
type VersionInfo struct {
	CreatedAt       time.Time         `json:"created_at"`
	CreatedBy       string            `json:"created_by"`
	Changelog       string            `json:"changelog,omitempty"`
	BreakingChanges bool              `json:"breaking_changes"`
	MigrationRequired bool            `json:"migration_required"`
	Compatibility   CompatibilityInfo `json:"compatibility"`
	Dependencies    []Dependency      `json:"dependencies,omitempty"`
	Rollback        RollbackInfo      `json:"rollback"`
}

// CompatibilityInfo represents compatibility information
type CompatibilityInfo struct {
	MinPlatformVersion string   `json:"min_platform_version"`
	MaxPlatformVersion string   `json:"max_platform_version"`
	DeprecatedFeatures []string `json:"deprecated_features,omitempty"`
	RemovedFeatures    []string `json:"removed_features,omitempty"`
}

// Dependency represents a workflow dependency
type Dependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Required bool   `json:"required"`
}

// RollbackInfo represents rollback information
type RollbackInfo struct {
	SafeRollbackVersions  []string `json:"safe_rollback_versions"`
	RollbackNotes         string   `json:"rollback_notes,omitempty"`
	DataMigrationRequired bool     `json:"data_migration_required"`
}

// BeforeCreate sets the ID before creating
func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Workflow model
func (Workflow) TableName() string {
	return "workflows"
}

// Validate validates the workflow
func (w *Workflow) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	
	if w.Version == "" {
		return fmt.Errorf("workflow version is required")
	}
	
	if w.Owner == "" {
		return fmt.Errorf("workflow owner is required")
	}
	
	if w.CreatedBy == "" {
		return fmt.Errorf("workflow created_by is required")
	}
	
	// Validate workflow definition
	if len(w.Definition.Spec.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	
	return nil
}

// GetFullName returns the full name including version
func (w *Workflow) GetFullName() string {
	return fmt.Sprintf("%s:%s", w.Name, w.Version)
}

// IsActive returns true if the workflow is active
func (w *Workflow) IsActive() bool {
	return w.Status == WorkflowStatusActive
}

// IsDeprecated returns true if the workflow is deprecated
func (w *Workflow) IsDeprecated() bool {
	return w.Status == WorkflowStatusDeprecated
}

// ToJSON converts the workflow to JSON
func (w *Workflow) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// FromJSON populates the workflow from JSON
func (w *Workflow) FromJSON(data []byte) error {
	return json.Unmarshal(data, w)
}