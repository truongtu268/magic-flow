package versioning

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"magic-flow/v2/internal/database"
	"magic-flow/v2/pkg/models"
)

// Migrator handles migration operations between workflow versions
type Migrator struct {
	repoManager database.RepositoryManager
	config      *VersioningConfig
}

// NewMigrator creates a new migrator instance
func NewMigrator(repoManager database.RepositoryManager) *Migrator {
	return &Migrator{
		repoManager: repoManager,
		config: &VersioningConfig{
			MigrationTimeout:     30 * time.Minute,
			MaxRollbackDepth:     10,
			BackupBeforeMigration: true,
		},
	}
}

// CreateMigrationPlan creates a migration plan between two versions
func (m *Migrator) CreateMigrationPlan(ctx context.Context, fromVersion, toVersion *models.WorkflowVersion) (*MigrationPlan, error) {
	plan := &MigrationPlan{
		ID:            uuid.New(),
		ToVersionID:   toVersion.ID,
		CreatedAt:     time.Now(),
		EstimatedTime: 5 * time.Minute, // Default estimate
		RiskLevel:     RiskLevelLow,
	}

	if fromVersion != nil {
		plan.FromVersionID = &fromVersion.ID
	}

	// Analyze differences and create migration steps
	steps, err := m.analyzeDifferences(fromVersion, toVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze differences: %w", err)
	}

	plan.MigrationSteps = steps
	plan.RollbackSteps = m.createRollbackSteps(steps)

	// Calculate risk level and estimated time
	plan.RiskLevel = m.calculateRiskLevel(steps)
	plan.EstimatedTime = m.estimateMigrationTime(steps)

	// Add validations
	plan.Validations = m.createValidationRules(fromVersion, toVersion)

	// Add prerequisites
	plan.Prerequisites = m.identifyPrerequisites(fromVersion, toVersion)

	return plan, nil
}

// CreateRollbackPlan creates a rollback plan from current to target version
func (m *Migrator) CreateRollbackPlan(ctx context.Context, currentVersion, targetVersion *models.WorkflowVersion) (*MigrationPlan, error) {
	// Create a reverse migration plan
	forwardPlan, err := m.CreateMigrationPlan(ctx, targetVersion, currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create forward plan for rollback: %w", err)
	}

	// Create rollback plan by reversing the forward plan
	rollbackPlan := &MigrationPlan{
		ID:              uuid.New(),
		FromVersionID:   &currentVersion.ID,
		ToVersionID:     targetVersion.ID,
		MigrationSteps:  m.reverseSteps(forwardPlan.RollbackSteps),
		RollbackSteps:   m.reverseSteps(forwardPlan.MigrationSteps),
		EstimatedTime:   forwardPlan.EstimatedTime,
		RiskLevel:       m.adjustRiskLevelForRollback(forwardPlan.RiskLevel),
		Prerequisites:   []string{"backup_verification", "execution_pause"},
		Validations:     m.createRollbackValidations(currentVersion, targetVersion),
		CreatedAt:       time.Now(),
	}

	return rollbackPlan, nil
}

// ExecuteMigration executes a migration plan
func (m *Migrator) ExecuteMigration(ctx context.Context, fromVersion, toVersion *models.WorkflowVersion) error {
	// Create migration plan
	plan, err := m.CreateMigrationPlan(ctx, fromVersion, toVersion)
	if err != nil {
		return fmt.Errorf("failed to create migration plan: %w", err)
	}

	// Create migration execution record
	execution := &MigrationExecution{
		ID:         uuid.New(),
		PlanID:     plan.ID,
		ToVersionID: toVersion.ID,
		Status:     MigrationStatusPending,
		StartedAt:  time.Now(),
		StepResults: make([]StepResult, 0, len(plan.MigrationSteps)),
	}

	if fromVersion != nil {
		execution.FromVersionID = &fromVersion.ID
		execution.WorkflowID = fromVersion.WorkflowID
	} else {
		execution.WorkflowID = toVersion.WorkflowID
	}

	// Execute migration steps
	err = m.executeMigrationSteps(ctx, execution, plan)
	if err != nil {
		execution.Status = MigrationStatusFailed
		execution.Error = err.Error()
		m.saveMigrationExecution(ctx, execution)
		return fmt.Errorf("migration execution failed: %w", err)
	}

	// Mark as completed
	now := time.Now()
	execution.Status = MigrationStatusCompleted
	execution.CompletedAt = &now
	execution.Duration = now.Sub(execution.StartedAt)

	// Save execution record
	err = m.saveMigrationExecution(ctx, execution)
	if err != nil {
		return fmt.Errorf("failed to save migration execution: %w", err)
	}

	return nil
}

// ExecuteRollback executes a rollback plan
func (m *Migrator) ExecuteRollback(ctx context.Context, plan *MigrationPlan) error {
	// Create rollback execution record
	execution := &MigrationExecution{
		ID:            uuid.New(),
		PlanID:        plan.ID,
		FromVersionID: plan.FromVersionID,
		ToVersionID:   plan.ToVersionID,
		Status:        MigrationStatusPending,
		StartedAt:     time.Now(),
		StepResults:   make([]StepResult, 0, len(plan.MigrationSteps)),
	}

	// Execute rollback steps
	err := m.executeMigrationSteps(ctx, execution, plan)
	if err != nil {
		execution.Status = MigrationStatusFailed
		execution.Error = err.Error()
		m.saveMigrationExecution(ctx, execution)
		return fmt.Errorf("rollback execution failed: %w", err)
	}

	// Mark as completed
	now := time.Now()
	execution.Status = MigrationStatusCompleted
	execution.CompletedAt = &now
	execution.Duration = now.Sub(execution.StartedAt)

	// Save execution record
	err = m.saveMigrationExecution(ctx, execution)
	if err != nil {
		return fmt.Errorf("failed to save rollback execution: %w", err)
	}

	return nil
}

// Private helper methods

func (m *Migrator) analyzeDifferences(fromVersion, toVersion *models.WorkflowVersion) ([]MigrationStep, error) {
	steps := []MigrationStep{}

	// If this is the first version, create initialization steps
	if fromVersion == nil {
		steps = append(steps, MigrationStep{
			ID:          uuid.New(),
			Order:       1,
			Type:        MigrationStepTypeSchemaUpdate,
			Description: "Initialize workflow schema",
			Action:      "create_workflow_schema",
			Timeout:     5 * time.Minute,
			RetryCount:  3,
		})
		return steps, nil
	}

	// Parse workflow definitions
	fromDef := fromVersion.Definition
	toDef := toVersion.Definition

	// Compare steps
	stepChanges := m.compareSteps(fromDef, toDef)
	for i, change := range stepChanges {
		steps = append(steps, MigrationStep{
			ID:          uuid.New(),
			Order:       i + 1,
			Type:        MigrationStepTypeSchemaUpdate,
			Description: change.Description,
			Action:      change.Action,
			Parameters:  change.Parameters,
			Timeout:     2 * time.Minute,
			RetryCount:  2,
		})
	}

	// Compare inputs/outputs
	ioChanges := m.compareInputsOutputs(fromDef, toDef)
	for i, change := range ioChanges {
		steps = append(steps, MigrationStep{
			ID:          uuid.New(),
			Order:       len(steps) + i + 1,
			Type:        MigrationStepTypeDataMigration,
			Description: change.Description,
			Action:      change.Action,
			Parameters:  change.Parameters,
			Timeout:     5 * time.Minute,
			RetryCount:  3,
		})
	}

	// Add validation step
	steps = append(steps, MigrationStep{
		ID:          uuid.New(),
		Order:       len(steps) + 1,
		Type:        MigrationStepTypeValidation,
		Description: "Validate migrated workflow",
		Action:      "validate_workflow",
		Timeout:     1 * time.Minute,
		RetryCount:  1,
	})

	return steps, nil
}

func (m *Migrator) compareSteps(fromDef, toDef map[string]interface{}) []StepChange {
	changes := []StepChange{}

	// Get steps from both definitions
	fromSteps, _ := fromDef["steps"].([]interface{})
	toSteps, _ := toDef["steps"].([]interface{})

	// Simple comparison - in a real implementation, this would be more sophisticated
	if len(fromSteps) != len(toSteps) {
		changes = append(changes, StepChange{
			Description: fmt.Sprintf("Step count changed from %d to %d", len(fromSteps), len(toSteps)),
			Action:      "update_step_count",
			Parameters: map[string]interface{}{
				"from_count": len(fromSteps),
				"to_count":   len(toSteps),
			},
		})
	}

	return changes
}

func (m *Migrator) compareInputsOutputs(fromDef, toDef map[string]interface{}) []StepChange {
	changes := []StepChange{}

	// Compare inputs
	fromInputs, _ := fromDef["inputs"].(map[string]interface{})
	toInputs, _ := toDef["inputs"].(map[string]interface{})

	if !m.mapsEqual(fromInputs, toInputs) {
		changes = append(changes, StepChange{
			Description: "Input schema changed",
			Action:      "migrate_input_schema",
			Parameters: map[string]interface{}{
				"from_inputs": fromInputs,
				"to_inputs":   toInputs,
			},
		})
	}

	// Compare outputs
	fromOutputs, _ := fromDef["outputs"].(map[string]interface{})
	toOutputs, _ := toDef["outputs"].(map[string]interface{})

	if !m.mapsEqual(fromOutputs, toOutputs) {
		changes = append(changes, StepChange{
			Description: "Output schema changed",
			Action:      "migrate_output_schema",
			Parameters: map[string]interface{}{
				"from_outputs": fromOutputs,
				"to_outputs":   toOutputs,
			},
		})
	}

	return changes
}

func (m *Migrator) createRollbackSteps(migrationSteps []MigrationStep) []MigrationStep {
	rollbackSteps := make([]MigrationStep, 0, len(migrationSteps))

	// Create rollback steps in reverse order
	for i := len(migrationSteps) - 1; i >= 0; i-- {
		step := migrationSteps[i]
		rollbackStep := MigrationStep{
			ID:          uuid.New(),
			Order:       len(migrationSteps) - i,
			Type:        step.Type,
			Description: fmt.Sprintf("Rollback: %s", step.Description),
			Action:      m.getRollbackAction(step.Action),
			Parameters:  m.getRollbackParameters(step.Parameters),
			Timeout:     step.Timeout,
			RetryCount:  step.RetryCount,
		}
		rollbackSteps = append(rollbackSteps, rollbackStep)
	}

	return rollbackSteps
}

func (m *Migrator) calculateRiskLevel(steps []MigrationStep) RiskLevel {
	// Calculate risk based on step types and complexity
	highRiskSteps := 0
	for _, step := range steps {
		if step.Type == MigrationStepTypeDataMigration || step.Type == MigrationStepTypeSchemaUpdate {
			highRiskSteps++
		}
	}

	if highRiskSteps > 5 {
		return RiskLevelCritical
	} else if highRiskSteps > 3 {
		return RiskLevelHigh
	} else if highRiskSteps > 1 {
		return RiskLevelMedium
	}
	return RiskLevelLow
}

func (m *Migrator) estimateMigrationTime(steps []MigrationStep) time.Duration {
	totalTime := time.Duration(0)
	for _, step := range steps {
		totalTime += step.Timeout
	}
	return totalTime
}

func (m *Migrator) createValidationRules(fromVersion, toVersion *models.WorkflowVersion) []ValidationRule {
	rules := []ValidationRule{
		{
			ID:          uuid.New(),
			Name:        "schema_validation",
			Description: "Validate workflow schema",
			Type:        ValidationTypeSchema,
			Required:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "compatibility_check",
			Description: "Check version compatibility",
			Type:        ValidationTypeCompatibility,
			Required:    true,
		},
	}

	return rules
}

func (m *Migrator) identifyPrerequisites(fromVersion, toVersion *models.WorkflowVersion) []string {
	prerequisites := []string{}

	if m.config.BackupBeforeMigration {
		prerequisites = append(prerequisites, "create_backup")
	}

	prerequisites = append(prerequisites, "pause_executions", "validate_permissions")

	return prerequisites
}

func (m *Migrator) executeMigrationSteps(ctx context.Context, execution *MigrationExecution, plan *MigrationPlan) error {
	execution.Status = MigrationStatusRunning

	for _, step := range plan.MigrationSteps {
		stepResult := StepResult{
			StepID:    step.ID,
			Status:    StepStatusPending,
			StartedAt: time.Now(),
		}

		// Execute step
		err := m.executeStep(ctx, step)
		if err != nil {
			stepResult.Status = StepStatusFailed
			stepResult.Error = err.Error()
			execution.StepResults = append(execution.StepResults, stepResult)
			return fmt.Errorf("step %s failed: %w", step.Description, err)
		}

		now := time.Now()
		stepResult.Status = StepStatusCompleted
		stepResult.CompletedAt = &now
		stepResult.Duration = now.Sub(stepResult.StartedAt)
		execution.StepResults = append(execution.StepResults, stepResult)
	}

	return nil
}

func (m *Migrator) executeStep(ctx context.Context, step MigrationStep) error {
	// This is a simplified implementation
	// In a real system, this would execute the actual migration logic
	switch step.Action {
	case "create_workflow_schema":
		return m.createWorkflowSchema(ctx, step.Parameters)
	case "update_step_count":
		return m.updateStepCount(ctx, step.Parameters)
	case "migrate_input_schema":
		return m.migrateInputSchema(ctx, step.Parameters)
	case "migrate_output_schema":
		return m.migrateOutputSchema(ctx, step.Parameters)
	case "validate_workflow":
		return m.validateWorkflow(ctx, step.Parameters)
	default:
		return fmt.Errorf("unknown migration action: %s", step.Action)
	}
}

// Migration action implementations (simplified)

func (m *Migrator) createWorkflowSchema(ctx context.Context, params map[string]interface{}) error {
	// Implementation for creating workflow schema
	return nil
}

func (m *Migrator) updateStepCount(ctx context.Context, params map[string]interface{}) error {
	// Implementation for updating step count
	return nil
}

func (m *Migrator) migrateInputSchema(ctx context.Context, params map[string]interface{}) error {
	// Implementation for migrating input schema
	return nil
}

func (m *Migrator) migrateOutputSchema(ctx context.Context, params map[string]interface{}) error {
	// Implementation for migrating output schema
	return nil
}

func (m *Migrator) validateWorkflow(ctx context.Context, params map[string]interface{}) error {
	// Implementation for validating workflow
	return nil
}

// Helper methods

func (m *Migrator) reverseSteps(steps []MigrationStep) []MigrationStep {
	reversed := make([]MigrationStep, len(steps))
	for i, step := range steps {
		reversed[len(steps)-1-i] = step
		reversed[len(steps)-1-i].Order = i + 1
	}
	return reversed
}

func (m *Migrator) adjustRiskLevelForRollback(originalRisk RiskLevel) RiskLevel {
	// Rollbacks are generally riskier
	switch originalRisk {
	case RiskLevelLow:
		return RiskLevelMedium
	case RiskLevelMedium:
		return RiskLevelHigh
	case RiskLevelHigh:
		return RiskLevelCritical
	default:
		return RiskLevelCritical
	}
}

func (m *Migrator) createRollbackValidations(currentVersion, targetVersion *models.WorkflowVersion) []ValidationRule {
	return []ValidationRule{
		{
			ID:          uuid.New(),
			Name:        "rollback_compatibility",
			Description: "Validate rollback compatibility",
			Type:        ValidationTypeCompatibility,
			Required:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "data_integrity",
			Description: "Validate data integrity after rollback",
			Type:        ValidationTypeData,
			Required:    true,
		},
	}
}

func (m *Migrator) getRollbackAction(action string) string {
	// Map migration actions to their rollback equivalents
	switch action {
	case "create_workflow_schema":
		return "remove_workflow_schema"
	case "update_step_count":
		return "revert_step_count"
	case "migrate_input_schema":
		return "revert_input_schema"
	case "migrate_output_schema":
		return "revert_output_schema"
	default:
		return "revert_" + action
	}
}

func (m *Migrator) getRollbackParameters(params map[string]interface{}) map[string]interface{} {
	// Swap from/to parameters for rollback
	rollbackParams := make(map[string]interface{})
	for k, v := range params {
		switch k {
		case "from_count":
			rollbackParams["to_count"] = v
		case "to_count":
			rollbackParams["from_count"] = v
		case "from_inputs":
			rollbackParams["to_inputs"] = v
		case "to_inputs":
			rollbackParams["from_inputs"] = v
		case "from_outputs":
			rollbackParams["to_outputs"] = v
		case "to_outputs":
			rollbackParams["from_outputs"] = v
		default:
			rollbackParams[k] = v
		}
	}
	return rollbackParams
}

func (m *Migrator) mapsEqual(map1, map2 map[string]interface{}) bool {
	// Simple map comparison - in a real implementation, this would be more sophisticated
	if len(map1) != len(map2) {
		return false
	}

	for k, v1 := range map1 {
		v2, exists := map2[k]
		if !exists {
			return false
		}

		// Simple value comparison
		v1Bytes, _ := json.Marshal(v1)
		v2Bytes, _ := json.Marshal(v2)
		if string(v1Bytes) != string(v2Bytes) {
			return false
		}
	}

	return true
}

func (m *Migrator) saveMigrationExecution(ctx context.Context, execution *MigrationExecution) error {
	// Save migration execution to database
	// This would typically use a dedicated repository
	return nil
}

// StepChange represents a change that needs to be migrated
type StepChange struct {
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters"`
}