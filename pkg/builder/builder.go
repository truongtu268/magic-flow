package builder

import (
	"fmt"
	"sync"

	"github.com/truongtu268/magic-flow/pkg/core"
)

// WorkflowTemplate represents a workflow template
type WorkflowTemplate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Steps       map[string]string `json:"steps"` // step_name -> step_type
	StartStep   string            `json:"start_step"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// WorkflowBuilder provides a fluent interface for building workflows
type WorkflowBuilder struct {
	name        string
	description string
	steps       []core.Step
	startStep   string
	middlewares []core.Middleware
	metadata    map[string]interface{}
	mu          sync.RWMutex
}

// NewWorkflowBuilder creates a new workflow builder
func NewWorkflowBuilder(name string) *WorkflowBuilder {
	return &WorkflowBuilder{
		name:     name,
		steps:    make([]core.Step, 0),
		metadata: make(map[string]interface{}),
	}
}

// WithDescription sets the workflow description
func (wb *WorkflowBuilder) WithDescription(description string) *WorkflowBuilder {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.description = description
	return wb
}

// WithStartStep sets the starting step
func (wb *WorkflowBuilder) WithStartStep(stepName string) *WorkflowBuilder {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.startStep = stepName
	return wb
}

// AddStep adds a step to the workflow
func (wb *WorkflowBuilder) AddStep(step core.Step) *WorkflowBuilder {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.steps = append(wb.steps, step)
	return wb
}

// AddFunctionStep adds a function step to the workflow
func (wb *WorkflowBuilder) AddFunctionStep(name, description string, executeFunc func(ctx *core.WorkflowContext) (*string, error)) *WorkflowBuilder {
	step := core.NewFunctionStep(name, description, executeFunc)
	return wb.AddStep(step)
}

// AddConditionalStep adds a conditional step to the workflow
func (wb *WorkflowBuilder) AddConditionalStep(name, description string, conditionFunc func(ctx *core.WorkflowContext) (bool, error), trueStep, falseStep string) *WorkflowBuilder {
	step := core.NewConditionalStep(name, description, conditionFunc, trueStep, falseStep)
	return wb.AddStep(step)
}

// AddWaitStep adds a wait step to the workflow
func (wb *WorkflowBuilder) AddWaitStep(name, description, triggerKey, nextStep string) *WorkflowBuilder {
	step := core.NewWaitStep(name, description, triggerKey, nextStep)
	return wb.AddStep(step)
}

// AddParallelStep adds a parallel step to the workflow
func (wb *WorkflowBuilder) AddParallelStep(name, description string, steps []string, nextStep string) *WorkflowBuilder {
	step := core.NewParallelStep(name, description, steps, nextStep)
	return wb.AddStep(step)
}

// AddRetryStep adds a retry step to the workflow
func (wb *WorkflowBuilder) AddRetryStep(name, description string, wrappedStep core.Step, maxRetries int) *WorkflowBuilder {
	step := core.NewRetryStep(name, description, wrappedStep, maxRetries)
	return wb.AddStep(step)
}

// AddMiddleware adds middleware to the workflow
func (wb *WorkflowBuilder) AddMiddleware(middleware core.Middleware) *WorkflowBuilder {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.middlewares = append(wb.middlewares, middleware)
	return wb
}

// WithMetadata adds metadata to the workflow
func (wb *WorkflowBuilder) WithMetadata(key string, value interface{}) *WorkflowBuilder {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.metadata[key] = value
	return wb
}

// Build builds the workflow and registers it with the engine
func (wb *WorkflowBuilder) Build(engine core.Engine) error {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	
	if wb.name == "" {
		return fmt.Errorf("workflow name is required")
	}
	
	if len(wb.steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	
	if wb.startStep == "" {
		return fmt.Errorf("start step is required")
	}
	
	// Validate that start step exists
	startStepExists := false
	for _, step := range wb.steps {
		if step.GetName() == wb.startStep {
			startStepExists = true
			break
		}
	}
	if !startStepExists {
		return fmt.Errorf("start step '%s' not found in workflow steps", wb.startStep)
	}
	
	// Register all steps with the engine
	for _, step := range wb.steps {
		if err := engine.RegisterStep(step); err != nil {
			return fmt.Errorf("failed to register step '%s': %w", step.GetName(), err)
		}
	}
	
	// Add all middleware to the engine
	for _, middleware := range wb.middlewares {
		engine.AddMiddleware(middleware)
	}
	
	return nil
}

// GetSteps returns all steps in the workflow
func (wb *WorkflowBuilder) GetSteps() []core.Step {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	steps := make([]core.Step, len(wb.steps))
	copy(steps, wb.steps)
	return steps
}

// GetStartStep returns the start step name
func (wb *WorkflowBuilder) GetStartStep() string {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	return wb.startStep
}

// GetName returns the workflow name
func (wb *WorkflowBuilder) GetName() string {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	return wb.name
}

// GetDescription returns the workflow description
func (wb *WorkflowBuilder) GetDescription() string {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	return wb.description
}

// GetMetadata returns the workflow metadata
func (wb *WorkflowBuilder) GetMetadata() map[string]interface{} {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	metadata := make(map[string]interface{})
	for k, v := range wb.metadata {
		metadata[k] = v
	}
	return metadata
}

// ToTemplate converts the workflow builder to a template
func (wb *WorkflowBuilder) ToTemplate(version string) *WorkflowTemplate {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	
	steps := make(map[string]string)
	for _, step := range wb.steps {
		steps[step.GetName()] = fmt.Sprintf("%T", step)
	}
	
	return &WorkflowTemplate{
		Name:        wb.name,
		Description: wb.description,
		Version:     version,
		Steps:       steps,
		StartStep:   wb.startStep,
		Metadata:    wb.GetMetadata(),
	}
}

// WorkflowRegistry manages workflow templates and builders
type WorkflowRegistry struct {
	templates map[string]*WorkflowTemplate
	builders  map[string]*WorkflowBuilder
	mu        sync.RWMutex
}

// NewWorkflowRegistry creates a new workflow registry
func NewWorkflowRegistry() *WorkflowRegistry {
	return &WorkflowRegistry{
		templates: make(map[string]*WorkflowTemplate),
		builders:  make(map[string]*WorkflowBuilder),
	}
}

// RegisterTemplate registers a workflow template
func (wr *WorkflowRegistry) RegisterTemplate(template *WorkflowTemplate) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}
	
	wr.templates[template.Name] = template
	return nil
}

// GetTemplate retrieves a workflow template by name
func (wr *WorkflowRegistry) GetTemplate(name string) (*WorkflowTemplate, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	template, exists := wr.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	
	return template, nil
}

// ListTemplates returns all registered template names
func (wr *WorkflowRegistry) ListTemplates() []string {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	names := make([]string, 0, len(wr.templates))
	for name := range wr.templates {
		names = append(names, name)
	}
	return names
}

// RegisterBuilder registers a workflow builder
func (wr *WorkflowRegistry) RegisterBuilder(builder *WorkflowBuilder) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	name := builder.GetName()
	if name == "" {
		return fmt.Errorf("builder name is required")
	}
	
	wr.builders[name] = builder
	return nil
}

// GetBuilder retrieves a workflow builder by name
func (wr *WorkflowRegistry) GetBuilder(name string) (*WorkflowBuilder, error) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	builder, exists := wr.builders[name]
	if !exists {
		return nil, fmt.Errorf("builder '%s' not found", name)
	}
	
	return builder, nil
}

// ListBuilders returns all registered builder names
func (wr *WorkflowRegistry) ListBuilders() []string {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	names := make([]string, 0, len(wr.builders))
	for name := range wr.builders {
		names = append(names, name)
	}
	return names
}

// CreateBuilderFromTemplate creates a workflow builder from a template
func (wr *WorkflowRegistry) CreateBuilderFromTemplate(templateName string) (*WorkflowBuilder, error) {
	template, err := wr.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}
	
	builder := NewWorkflowBuilder(template.Name)
	builder.WithDescription(template.Description)
	builder.WithStartStep(template.StartStep)
	
	// Add metadata from template
	for key, value := range template.Metadata {
		builder.WithMetadata(key, value)
	}
	
	return builder, nil
}

// RemoveTemplate removes a template from the registry
func (wr *WorkflowRegistry) RemoveTemplate(name string) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	if _, exists := wr.templates[name]; !exists {
		return fmt.Errorf("template '%s' not found", name)
	}
	
	delete(wr.templates, name)
	return nil
}

// RemoveBuilder removes a builder from the registry
func (wr *WorkflowRegistry) RemoveBuilder(name string) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	if _, exists := wr.builders[name]; !exists {
		return fmt.Errorf("builder '%s' not found", name)
	}
	
	delete(wr.builders, name)
	return nil
}

// Clear removes all templates and builders
func (wr *WorkflowRegistry) Clear() {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	wr.templates = make(map[string]*WorkflowTemplate)
	wr.builders = make(map[string]*WorkflowBuilder)
}