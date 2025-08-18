package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"magic-flow/v2/pkg/models"
)

// WorkflowParser handles parsing and validation of workflow definitions
type WorkflowParser struct{}

// NewWorkflowParser creates a new workflow parser
func NewWorkflowParser() *WorkflowParser {
	return &WorkflowParser{}
}

// ParseYAML parses a YAML workflow definition into a Workflow model
func (p *WorkflowParser) ParseYAML(yamlContent []byte) (*models.Workflow, error) {
	var yamlWorkflow YAMLWorkflow
	if err := yaml.Unmarshal(yamlContent, &yamlWorkflow); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return p.convertToWorkflow(&yamlWorkflow)
}

// ParseJSON parses a JSON workflow definition into a Workflow model
func (p *WorkflowParser) ParseJSON(jsonContent []byte) (*models.Workflow, error) {
	var jsonWorkflow JSONWorkflow
	if err := json.Unmarshal(jsonContent, &jsonWorkflow); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return p.convertJSONToWorkflow(&jsonWorkflow)
}

// ValidateWorkflow validates a workflow definition
func (p *WorkflowParser) ValidateWorkflow(workflow *models.Workflow) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if workflow.Definition == nil {
		return fmt.Errorf("workflow definition is required")
	}

	if len(workflow.Definition.Spec.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	// Validate steps
	stepIDs := make(map[string]bool)
	for i, step := range workflow.Definition.Spec.Steps {
		if step.ID == "" {
			return fmt.Errorf("step %d: ID is required", i)
		}

		if stepIDs[step.ID] {
			return fmt.Errorf("step %d: duplicate step ID '%s'", i, step.ID)
		}
		stepIDs[step.ID] = true

		if step.Type == "" {
			return fmt.Errorf("step %d (%s): type is required", i, step.ID)
		}

		if err := p.validateStepType(step); err != nil {
			return fmt.Errorf("step %d (%s): %w", i, step.ID, err)
		}

		// Validate dependencies
		for _, dep := range step.DependsOn {
			if !stepIDs[dep] && dep != step.ID {
				// Check if dependency exists in previous steps
				found := false
				for j := 0; j < i; j++ {
					if workflow.Definition.Spec.Steps[j].ID == dep {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("step %d (%s): dependency '%s' not found", i, step.ID, dep)
				}
			}
		}
	}

	// Validate triggers
	if workflow.Definition.Spec.Triggers != nil {
		for i, trigger := range workflow.Definition.Spec.Triggers {
			if trigger.Type == "" {
				return fmt.Errorf("trigger %d: type is required", i)
			}
		}
	}

	return nil
}

func (p *WorkflowParser) validateStepType(step models.WorkflowStep) error {
	switch step.Type {
	case "http":
		return p.validateHTTPStep(step)
	case "script":
		return p.validateScriptStep(step)
	case "transform":
		return p.validateTransformStep(step)
	case "delay":
		return p.validateDelayStep(step)
	case "conditional":
		return p.validateConditionalStep(step)
	default:
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}

func (p *WorkflowParser) validateHTTPStep(step models.WorkflowStep) error {
	if step.Config == nil {
		return fmt.Errorf("HTTP step requires config")
	}

	url, ok := step.Config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("HTTP step requires 'url' in config")
	}

	method, ok := step.Config["method"].(string)
	if !ok || method == "" {
		method = "GET"
	}

	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	validMethod := false
	for _, vm := range validMethods {
		if strings.ToUpper(method) == vm {
			validMethod = true
			break
		}
	}
	if !validMethod {
		return fmt.Errorf("invalid HTTP method: %s", method)
	}

	return nil
}

func (p *WorkflowParser) validateScriptStep(step models.WorkflowStep) error {
	if step.Config == nil {
		return fmt.Errorf("script step requires config")
	}

	command, ok := step.Config["command"].(string)
	if !ok || command == "" {
		return fmt.Errorf("script step requires 'command' in config")
	}

	return nil
}

func (p *WorkflowParser) validateTransformStep(step models.WorkflowStep) error {
	if step.Config == nil {
		return fmt.Errorf("transform step requires config")
	}

	transformType, ok := step.Config["type"].(string)
	if !ok || transformType == "" {
		return fmt.Errorf("transform step requires 'type' in config")
	}

	validTypes := []string{"json", "filter", "map", "aggregate"}
	validType := false
	for _, vt := range validTypes {
		if transformType == vt {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid transform type: %s", transformType)
	}

	return nil
}

func (p *WorkflowParser) validateDelayStep(step models.WorkflowStep) error {
	if step.Config == nil {
		return fmt.Errorf("delay step requires config")
	}

	duration, ok := step.Config["duration"].(string)
	if !ok || duration == "" {
		return fmt.Errorf("delay step requires 'duration' in config")
	}

	return nil
}

func (p *WorkflowParser) validateConditionalStep(step models.WorkflowStep) error {
	if step.Config == nil {
		return fmt.Errorf("conditional step requires config")
	}

	condition, ok := step.Config["condition"].(string)
	if !ok || condition == "" {
		return fmt.Errorf("conditional step requires 'condition' in config")
	}

	return nil
}

func (p *WorkflowParser) convertToWorkflow(yamlWorkflow *YAMLWorkflow) (*models.Workflow, error) {
	workflow := &models.Workflow{
		ID:          uuid.New(),
		Name:        yamlWorkflow.Name,
		Description: yamlWorkflow.Description,
		Status:      models.WorkflowStatusDraft,
		Definition: &models.WorkflowDefinition{
			Version: yamlWorkflow.Version,
			Metadata: models.WorkflowMetadata{
				Name:        yamlWorkflow.Name,
				Description: yamlWorkflow.Description,
				Version:     yamlWorkflow.Version,
				Labels:      yamlWorkflow.Labels,
				Annotations: yamlWorkflow.Annotations,
			},
			Spec: models.WorkflowSpec{
				Steps:    make([]models.WorkflowStep, len(yamlWorkflow.Steps)),
				Triggers: convertYAMLTriggers(yamlWorkflow.Triggers),
			},
		},
	}

	// Convert steps
	for i, yamlStep := range yamlWorkflow.Steps {
		step := models.WorkflowStep{
			ID:          yamlStep.ID,
			Name:        yamlStep.Name,
			Description: yamlStep.Description,
			Type:        yamlStep.Type,
			Config:      yamlStep.Config,
			DependsOn:   yamlStep.DependsOn,
			Timeout:     yamlStep.Timeout,
			RetryPolicy: convertYAMLRetryPolicy(yamlStep.Retry),
		}

		// Convert input/output mappings
		if yamlStep.Input != nil {
			step.InputMapping = &models.DataMapping{
				Mappings: yamlStep.Input,
			}
		}
		if yamlStep.Output != nil {
			step.OutputMapping = &models.DataMapping{
				Mappings: yamlStep.Output,
			}
		}

		// Convert error handling
		if yamlStep.OnError != nil {
			step.ErrorHandling = &models.ErrorHandling{
				Strategy:    yamlStep.OnError.Strategy,
				FallbackStep: yamlStep.OnError.FallbackStep,
				IgnoreErrors: yamlStep.OnError.IgnoreErrors,
			}
		}

		workflow.Definition.Spec.Steps[i] = step
	}

	return workflow, nil
}

func (p *WorkflowParser) convertJSONToWorkflow(jsonWorkflow *JSONWorkflow) (*models.Workflow, error) {
	// Convert JSON workflow to YAML format first, then use existing conversion
	yamlWorkflow := &YAMLWorkflow{
		Name:        jsonWorkflow.Name,
		Description: jsonWorkflow.Description,
		Version:     jsonWorkflow.Version,
		Labels:      jsonWorkflow.Labels,
		Annotations: jsonWorkflow.Annotations,
		Steps:       make([]YAMLStep, len(jsonWorkflow.Steps)),
		Triggers:    convertJSONTriggers(jsonWorkflow.Triggers),
	}

	for i, jsonStep := range jsonWorkflow.Steps {
		yamlWorkflow.Steps[i] = YAMLStep{
			ID:          jsonStep.ID,
			Name:        jsonStep.Name,
			Description: jsonStep.Description,
			Type:        jsonStep.Type,
			Config:      jsonStep.Config,
			Input:       jsonStep.Input,
			Output:      jsonStep.Output,
			DependsOn:   jsonStep.DependsOn,
			Timeout:     jsonStep.Timeout,
			Retry:       convertJSONRetryPolicy(jsonStep.Retry),
			OnError:     convertJSONErrorHandling(jsonStep.OnError),
		}
	}

	return p.convertToWorkflow(yamlWorkflow)
}

func convertYAMLTriggers(yamlTriggers []YAMLTrigger) []models.WorkflowTrigger {
	if yamlTriggers == nil {
		return nil
	}

	triggers := make([]models.WorkflowTrigger, len(yamlTriggers))
	for i, yamlTrigger := range yamlTriggers {
		triggers[i] = models.WorkflowTrigger{
			Type:   yamlTrigger.Type,
			Config: yamlTrigger.Config,
		}
	}
	return triggers
}

func convertYAMLRetryPolicy(yamlRetry *YAMLRetryPolicy) *models.RetryPolicy {
	if yamlRetry == nil {
		return nil
	}

	return &models.RetryPolicy{
		MaxAttempts: yamlRetry.MaxAttempts,
		Delay:       yamlRetry.Delay,
		Backoff:     yamlRetry.Backoff,
		MaxDelay:    yamlRetry.MaxDelay,
	}
}

func convertJSONTriggers(jsonTriggers []JSONTrigger) []YAMLTrigger {
	if jsonTriggers == nil {
		return nil
	}

	triggers := make([]YAMLTrigger, len(jsonTriggers))
	for i, jsonTrigger := range jsonTriggers {
		triggers[i] = YAMLTrigger{
			Type:   jsonTrigger.Type,
			Config: jsonTrigger.Config,
		}
	}
	return triggers
}

func convertJSONRetryPolicy(jsonRetry *JSONRetryPolicy) *YAMLRetryPolicy {
	if jsonRetry == nil {
		return nil
	}

	return &YAMLRetryPolicy{
		MaxAttempts: jsonRetry.MaxAttempts,
		Delay:       jsonRetry.Delay,
		Backoff:     jsonRetry.Backoff,
		MaxDelay:    jsonRetry.MaxDelay,
	}
}

func convertJSONErrorHandling(jsonError *JSONErrorHandling) *YAMLErrorHandling {
	if jsonError == nil {
		return nil
	}

	return &YAMLErrorHandling{
		Strategy:     jsonError.Strategy,
		FallbackStep: jsonError.FallbackStep,
		IgnoreErrors: jsonError.IgnoreErrors,
	}
}

// YAML workflow definition structures
type YAMLWorkflow struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description,omitempty"`
	Version     string                 `yaml:"version"`
	Labels      map[string]string      `yaml:"labels,omitempty"`
	Annotations map[string]string      `yaml:"annotations,omitempty"`
	Steps       []YAMLStep             `yaml:"steps"`
	Triggers    []YAMLTrigger          `yaml:"triggers,omitempty"`
}

type YAMLStep struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Type        string                 `yaml:"type"`
	Config      map[string]interface{} `yaml:"config,omitempty"`
	Input       map[string]string      `yaml:"input,omitempty"`
	Output      map[string]string      `yaml:"output,omitempty"`
	DependsOn   []string               `yaml:"depends_on,omitempty"`
	Timeout     string                 `yaml:"timeout,omitempty"`
	Retry       *YAMLRetryPolicy       `yaml:"retry,omitempty"`
	OnError     *YAMLErrorHandling     `yaml:"on_error,omitempty"`
}

type YAMLTrigger struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config,omitempty"`
}

type YAMLRetryPolicy struct {
	MaxAttempts int    `yaml:"max_attempts,omitempty"`
	Delay       string `yaml:"delay,omitempty"`
	Backoff     string `yaml:"backoff,omitempty"`
	MaxDelay    string `yaml:"max_delay,omitempty"`
}

type YAMLErrorHandling struct {
	Strategy     string   `yaml:"strategy,omitempty"`
	FallbackStep string   `yaml:"fallback_step,omitempty"`
	IgnoreErrors []string `yaml:"ignore_errors,omitempty"`
}

// JSON workflow definition structures
type JSONWorkflow struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Steps       []JSONStep             `json:"steps"`
	Triggers    []JSONTrigger          `json:"triggers,omitempty"`
}

type JSONStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Input       map[string]string      `json:"input,omitempty"`
	Output      map[string]string      `json:"output,omitempty"`
	DependsOn   []string               `json:"depends_on,omitempty"`
	Timeout     string                 `json:"timeout,omitempty"`
	Retry       *JSONRetryPolicy       `json:"retry,omitempty"`
	OnError     *JSONErrorHandling     `json:"on_error,omitempty"`
}

type JSONTrigger struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type JSONRetryPolicy struct {
	MaxAttempts int    `json:"max_attempts,omitempty"`
	Delay       string `json:"delay,omitempty"`
	Backoff     string `json:"backoff,omitempty"`
	MaxDelay    string `json:"max_delay,omitempty"`
}

type JSONErrorHandling struct {
	Strategy     string   `json:"strategy,omitempty"`
	FallbackStep string   `json:"fallback_step,omitempty"`
	IgnoreErrors []string `json:"ignore_errors,omitempty"`
}