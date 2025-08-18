package versioning

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"magic-flow/v2/pkg/models"
)

// Validator handles validation of workflow versions and migrations
type Validator struct {
	config *ValidationConfig
}

// ValidationConfig contains configuration for validation
type ValidationConfig struct {
	StrictMode          bool     `json:"strict_mode"`
	AllowedStepTypes    []string `json:"allowed_step_types"`
	MaxStepsPerWorkflow int      `json:"max_steps_per_workflow"`
	RequiredFields      []string `json:"required_fields"`
	CustomValidators    []string `json:"custom_validators"`
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		config: &ValidationConfig{
			StrictMode:          false,
			AllowedStepTypes:    []string{"http", "script", "condition", "loop", "parallel", "custom"},
			MaxStepsPerWorkflow: 100,
			RequiredFields:      []string{"name", "steps"},
		},
	}
}

// ValidateVersion validates a new version before creation
func (v *Validator) ValidateVersion(ctx context.Context, workflow *models.Workflow, changes VersionChanges) error {
	// Validate change type
	if err := v.validateChangeType(changes.ChangeType); err != nil {
		return fmt.Errorf("invalid change type: %w", err)
	}

	// Validate workflow definition
	if err := v.validateWorkflowDefinition(changes.NewDefinition); err != nil {
		return fmt.Errorf("invalid workflow definition: %w", err)
	}

	// Validate compatibility with existing executions
	if err := v.validateExecutionCompatibility(ctx, workflow, changes); err != nil {
		return fmt.Errorf("execution compatibility check failed: %w", err)
	}

	// Validate schema changes
	if err := v.validateSchemaChanges(workflow.Definition, changes.NewDefinition, changes.ChangeType); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// Validate business rules
	if err := v.validateBusinessRules(ctx, workflow, changes); err != nil {
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	return nil
}

// ValidateMigrationPlan validates a migration plan before execution
func (v *Validator) ValidateMigrationPlan(ctx context.Context, plan *MigrationPlan) error {
	// Validate plan structure
	if err := v.validatePlanStructure(plan); err != nil {
		return fmt.Errorf("invalid plan structure: %w", err)
	}

	// Validate migration steps
	for i, step := range plan.MigrationSteps {
		if err := v.validateMigrationStep(step, i); err != nil {
			return fmt.Errorf("invalid migration step %d: %w", i, err)
		}
	}

	// Validate rollback steps
	for i, step := range plan.RollbackSteps {
		if err := v.validateMigrationStep(step, i); err != nil {
			return fmt.Errorf("invalid rollback step %d: %w", i, err)
		}
	}

	// Validate risk assessment
	if err := v.validateRiskAssessment(plan); err != nil {
		return fmt.Errorf("risk assessment validation failed: %w", err)
	}

	return nil
}

// ValidateRollback validates a rollback operation
func (v *Validator) ValidateRollback(ctx context.Context, currentVersion, targetVersion *models.WorkflowVersion) error {
	// Validate version compatibility
	if err := v.validateRollbackCompatibility(currentVersion, targetVersion); err != nil {
		return fmt.Errorf("rollback compatibility check failed: %w", err)
	}

	// Validate data integrity requirements
	if err := v.validateDataIntegrity(ctx, currentVersion, targetVersion); err != nil {
		return fmt.Errorf("data integrity validation failed: %w", err)
	}

	// Validate business constraints
	if err := v.validateRollbackConstraints(ctx, currentVersion, targetVersion); err != nil {
		return fmt.Errorf("rollback constraint validation failed: %w", err)
	}

	return nil
}

// Private validation methods

func (v *Validator) validateChangeType(changeType ChangeType) error {
	validTypes := []ChangeType{ChangeTypeMajor, ChangeTypeMinor, ChangeTypePatch}
	for _, validType := range validTypes {
		if changeType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid change type: %s", changeType)
}

func (v *Validator) validateWorkflowDefinition(definition map[string]interface{}) error {
	// Validate required fields
	for _, field := range v.config.RequiredFields {
		if _, exists := definition[field]; !exists {
			return fmt.Errorf("required field missing: %s", field)
		}
	}

	// Validate workflow name
	name, ok := definition["name"].(string)
	if !ok {
		return fmt.Errorf("workflow name must be a string")
	}
	if err := v.validateWorkflowName(name); err != nil {
		return fmt.Errorf("invalid workflow name: %w", err)
	}

	// Validate steps
	steps, ok := definition["steps"].([]interface{})
	if !ok {
		return fmt.Errorf("steps must be an array")
	}
	if err := v.validateSteps(steps); err != nil {
		return fmt.Errorf("invalid steps: %w", err)
	}

	// Validate inputs if present
	if inputs, exists := definition["inputs"]; exists {
		if err := v.validateInputsOutputs(inputs, "inputs"); err != nil {
			return fmt.Errorf("invalid inputs: %w", err)
		}
	}

	// Validate outputs if present
	if outputs, exists := definition["outputs"]; exists {
		if err := v.validateInputsOutputs(outputs, "outputs"); err != nil {
			return fmt.Errorf("invalid outputs: %w", err)
		}
	}

	return nil
}

func (v *Validator) validateWorkflowName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("workflow name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("workflow name too long (max 100 characters)")
	}

	// Check for valid characters
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("workflow name contains invalid characters")
	}

	return nil
}

func (v *Validator) validateSteps(steps []interface{}) error {
	if len(steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	if len(steps) > v.config.MaxStepsPerWorkflow {
		return fmt.Errorf("too many steps (max %d)", v.config.MaxStepsPerWorkflow)
	}

	stepNames := make(map[string]bool)
	for i, stepInterface := range steps {
		step, ok := stepInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("step %d is not a valid object", i)
		}

		// Validate step structure
		if err := v.validateStep(step, i); err != nil {
			return fmt.Errorf("step %d validation failed: %w", i, err)
		}

		// Check for duplicate step names
		if name, exists := step["name"].(string); exists {
			if stepNames[name] {
				return fmt.Errorf("duplicate step name: %s", name)
			}
			stepNames[name] = true
		}
	}

	return nil
}

func (v *Validator) validateStep(step map[string]interface{}, index int) error {
	// Validate required step fields
	requiredStepFields := []string{"name", "type"}
	for _, field := range requiredStepFields {
		if _, exists := step[field]; !exists {
			return fmt.Errorf("required field missing: %s", field)
		}
	}

	// Validate step name
	name, ok := step["name"].(string)
	if !ok {
		return fmt.Errorf("step name must be a string")
	}
	if err := v.validateStepName(name); err != nil {
		return fmt.Errorf("invalid step name: %w", err)
	}

	// Validate step type
	stepType, ok := step["type"].(string)
	if !ok {
		return fmt.Errorf("step type must be a string")
	}
	if err := v.validateStepType(stepType); err != nil {
		return fmt.Errorf("invalid step type: %w", err)
	}

	// Validate step-specific configuration
	if err := v.validateStepConfiguration(step, stepType); err != nil {
		return fmt.Errorf("invalid step configuration: %w", err)
	}

	return nil
}

func (v *Validator) validateStepName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("step name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("step name too long (max 50 characters)")
	}

	// Check for valid characters
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("step name contains invalid characters")
	}

	return nil
}

func (v *Validator) validateStepType(stepType string) error {
	for _, allowedType := range v.config.AllowedStepTypes {
		if stepType == allowedType {
			return nil
		}
	}
	return fmt.Errorf("unsupported step type: %s", stepType)
}

func (v *Validator) validateStepConfiguration(step map[string]interface{}, stepType string) error {
	switch stepType {
	case "http":
		return v.validateHTTPStep(step)
	case "script":
		return v.validateScriptStep(step)
	case "condition":
		return v.validateConditionStep(step)
	case "loop":
		return v.validateLoopStep(step)
	case "parallel":
		return v.validateParallelStep(step)
	default:
		return nil // Custom steps don't have specific validation
	}
}

func (v *Validator) validateHTTPStep(step map[string]interface{}) error {
	config, exists := step["config"]
	if !exists {
		return fmt.Errorf("HTTP step requires config")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("HTTP step config must be an object")
	}

	// Validate required HTTP fields
	if _, exists := configMap["url"]; !exists {
		return fmt.Errorf("HTTP step requires url")
	}
	if _, exists := configMap["method"]; !exists {
		return fmt.Errorf("HTTP step requires method")
	}

	return nil
}

func (v *Validator) validateScriptStep(step map[string]interface{}) error {
	config, exists := step["config"]
	if !exists {
		return fmt.Errorf("script step requires config")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("script step config must be an object")
	}

	// Validate required script fields
	if _, exists := configMap["script"]; !exists {
		return fmt.Errorf("script step requires script")
	}

	return nil
}

func (v *Validator) validateConditionStep(step map[string]interface{}) error {
	config, exists := step["config"]
	if !exists {
		return fmt.Errorf("condition step requires config")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("condition step config must be an object")
	}

	// Validate required condition fields
	if _, exists := configMap["condition"]; !exists {
		return fmt.Errorf("condition step requires condition")
	}

	return nil
}

func (v *Validator) validateLoopStep(step map[string]interface{}) error {
	config, exists := step["config"]
	if !exists {
		return fmt.Errorf("loop step requires config")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("loop step config must be an object")
	}

	// Validate required loop fields
	if _, exists := configMap["items"]; !exists {
		return fmt.Errorf("loop step requires items")
	}

	return nil
}

func (v *Validator) validateParallelStep(step map[string]interface{}) error {
	config, exists := step["config"]
	if !exists {
		return fmt.Errorf("parallel step requires config")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("parallel step config must be an object")
	}

	// Validate required parallel fields
	if _, exists := configMap["branches"]; !exists {
		return fmt.Errorf("parallel step requires branches")
	}

	return nil
}

func (v *Validator) validateInputsOutputs(schema interface{}, schemaType string) error {
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%s must be an object", schemaType)
	}

	// Validate each field in the schema
	for fieldName, fieldDef := range schemaMap {
		if err := v.validateSchemaField(fieldName, fieldDef); err != nil {
			return fmt.Errorf("invalid field %s: %w", fieldName, err)
		}
	}

	return nil
}

func (v *Validator) validateSchemaField(fieldName string, fieldDef interface{}) error {
	fieldMap, ok := fieldDef.(map[string]interface{})
	if !ok {
		return fmt.Errorf("field definition must be an object")
	}

	// Validate field type
	fieldType, exists := fieldMap["type"]
	if !exists {
		return fmt.Errorf("field type is required")
	}

	fieldTypeStr, ok := fieldType.(string)
	if !ok {
		return fmt.Errorf("field type must be a string")
	}

	validTypes := []string{"string", "number", "integer", "boolean", "array", "object"}
	validType := false
	for _, vt := range validTypes {
		if fieldTypeStr == vt {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid field type: %s", fieldTypeStr)
	}

	return nil
}

func (v *Validator) validateExecutionCompatibility(ctx context.Context, workflow *models.Workflow, changes VersionChanges) error {
	// Check if there are running executions
	// This would typically query the database for active executions
	// For now, we'll implement a basic check

	if changes.ChangeType == ChangeTypeMajor {
		// Major changes might break running executions
		return fmt.Errorf("major changes require all executions to be completed before deployment")
	}

	return nil
}

func (v *Validator) validateSchemaChanges(oldDef, newDef map[string]interface{}, changeType ChangeType) error {
	// Validate input schema changes
	if err := v.validateSchemaCompatibility(oldDef["inputs"], newDef["inputs"], "inputs", changeType); err != nil {
		return fmt.Errorf("input schema compatibility check failed: %w", err)
	}

	// Validate output schema changes
	if err := v.validateSchemaCompatibility(oldDef["outputs"], newDef["outputs"], "outputs", changeType); err != nil {
		return fmt.Errorf("output schema compatibility check failed: %w", err)
	}

	return nil
}

func (v *Validator) validateSchemaCompatibility(oldSchema, newSchema interface{}, schemaType string, changeType ChangeType) error {
	// For patch changes, schema should be fully backward compatible
	if changeType == ChangeTypePatch {
		if !v.isSchemaBinaryCompatible(oldSchema, newSchema) {
			return fmt.Errorf("patch changes must maintain full %s schema compatibility", schemaType)
		}
	}

	// For minor changes, new fields can be added but existing fields should not be removed
	if changeType == ChangeTypeMinor {
		if !v.isSchemaBackwardCompatible(oldSchema, newSchema) {
			return fmt.Errorf("minor changes must maintain backward %s schema compatibility", schemaType)
		}
	}

	// Major changes can break compatibility but should be documented
	if changeType == ChangeTypeMajor {
		// Major changes are allowed but we might want to validate they're intentional
		if v.config.StrictMode && v.isSchemaBinaryCompatible(oldSchema, newSchema) {
			return fmt.Errorf("major version bump detected but no breaking %s schema changes found", schemaType)
		}
	}

	return nil
}

func (v *Validator) isSchemaBinaryCompatible(oldSchema, newSchema interface{}) bool {
	// Simplified binary compatibility check
	// In a real implementation, this would do deep schema comparison
	return fmt.Sprintf("%v", oldSchema) == fmt.Sprintf("%v", newSchema)
}

func (v *Validator) isSchemaBackwardCompatible(oldSchema, newSchema interface{}) bool {
	// Simplified backward compatibility check
	// This should check that all fields in oldSchema exist in newSchema
	// and that their types are compatible

	oldMap, ok1 := oldSchema.(map[string]interface{})
	newMap, ok2 := newSchema.(map[string]interface{})

	if !ok1 || !ok2 {
		return true // If either is not a map, consider compatible
	}

	// Check that all old fields exist in new schema
	for fieldName, oldFieldDef := range oldMap {
		newFieldDef, exists := newMap[fieldName]
		if !exists {
			return false // Field was removed
		}

		// Check field type compatibility
		if !v.areFieldTypesCompatible(oldFieldDef, newFieldDef) {
			return false
		}
	}

	return true
}

func (v *Validator) areFieldTypesCompatible(oldField, newField interface{}) bool {
	// Simplified field type compatibility check
	oldMap, ok1 := oldField.(map[string]interface{})
	newMap, ok2 := newField.(map[string]interface{})

	if !ok1 || !ok2 {
		return true
	}

	oldType, _ := oldMap["type"].(string)
	newType, _ := newMap["type"].(string)

	return oldType == newType
}

func (v *Validator) validateBusinessRules(ctx context.Context, workflow *models.Workflow, changes VersionChanges) error {
	// Implement business-specific validation rules
	// This could include checks for:
	// - Workflow naming conventions
	// - Step count limits
	// - Resource usage constraints
	// - Security requirements
	// - Compliance requirements

	return nil
}

func (v *Validator) validatePlanStructure(plan *MigrationPlan) error {
	if plan.ID == (uuid.UUID{}) {
		return fmt.Errorf("migration plan must have an ID")
	}

	if plan.ToVersionID == (uuid.UUID{}) {
		return fmt.Errorf("migration plan must specify target version")
	}

	if len(plan.MigrationSteps) == 0 {
		return fmt.Errorf("migration plan must have at least one step")
	}

	return nil
}

func (v *Validator) validateMigrationStep(step MigrationStep, index int) error {
	if step.ID == (uuid.UUID{}) {
		return fmt.Errorf("migration step must have an ID")
	}

	if step.Order != index+1 {
		return fmt.Errorf("migration step order mismatch: expected %d, got %d", index+1, step.Order)
	}

	if step.Description == "" {
		return fmt.Errorf("migration step must have a description")
	}

	if step.Action == "" {
		return fmt.Errorf("migration step must have an action")
	}

	if step.Timeout <= 0 {
		return fmt.Errorf("migration step must have a positive timeout")
	}

	return nil
}

func (v *Validator) validateRiskAssessment(plan *MigrationPlan) error {
	validRiskLevels := []RiskLevel{RiskLevelLow, RiskLevelMedium, RiskLevelHigh, RiskLevelCritical}
	validRisk := false
	for _, level := range validRiskLevels {
		if plan.RiskLevel == level {
			validRisk = true
			break
		}
	}

	if !validRisk {
		return fmt.Errorf("invalid risk level: %s", plan.RiskLevel)
	}

	// Validate that high-risk migrations have appropriate safeguards
	if plan.RiskLevel == RiskLevelHigh || plan.RiskLevel == RiskLevelCritical {
		if len(plan.Prerequisites) == 0 {
			return fmt.Errorf("high-risk migrations must have prerequisites")
		}
		if len(plan.Validations) == 0 {
			return fmt.Errorf("high-risk migrations must have validations")
		}
	}

	return nil
}

func (v *Validator) validateRollbackCompatibility(currentVersion, targetVersion *models.WorkflowVersion) error {
	if currentVersion.WorkflowID != targetVersion.WorkflowID {
		return fmt.Errorf("versions belong to different workflows")
	}

	// Check if rollback is to a newer version (which doesn't make sense)
	if strings.Compare(currentVersion.Version, targetVersion.Version) < 0 {
		return fmt.Errorf("cannot rollback to a newer version")
	}

	return nil
}

func (v *Validator) validateDataIntegrity(ctx context.Context, currentVersion, targetVersion *models.WorkflowVersion) error {
	// Validate that rollback won't cause data loss or corruption
	// This would typically involve checking:
	// - Schema compatibility
	// - Data migration requirements
	// - Execution state compatibility

	return nil
}

func (v *Validator) validateRollbackConstraints(ctx context.Context, currentVersion, targetVersion *models.WorkflowVersion) error {
	// Validate business constraints for rollback
	// This could include:
	// - Time-based constraints (e.g., can't rollback after X days)
	// - Execution-based constraints (e.g., can't rollback if certain executions exist)
	// - Approval requirements for rollbacks

	return nil
}