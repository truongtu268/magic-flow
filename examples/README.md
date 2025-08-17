# Magic Flow Examples

This directory contains examples demonstrating how to use the Magic Flow workflow engine library. Each example showcases different features and use cases.

## Examples Overview

### 1. Basic Workflow (`basic_workflow/`)

A simple example demonstrating the fundamental concepts of Magic Flow:

- Creating a workflow engine
- Defining workflow steps using `FunctionStep`
- Executing a linear workflow
- Working with workflow data and step results

**Key Features Demonstrated:**
- Basic workflow execution
- Data validation and processing
- Step result storage
- Error handling

**Run the example:**
```bash
cd examples/basic_workflow
go run main.go
```

### 2. Advanced Workflow (`advanced_workflow/`)

A comprehensive example showcasing advanced features:

- Custom configuration and logging
- Middleware implementation (timing, validation)
- Conditional step execution
- Workflow metadata and tags
- Complex data processing scenarios

**Key Features Demonstrated:**
- Custom middleware creation
- Conditional workflow branching
- Workflow metadata management
- Performance monitoring
- Advanced error handling

**Run the example:**
```bash
cd examples/advanced_workflow
go run main.go
```

## Core Concepts

### Workflow Engine

The workflow engine is the central component that orchestrates step execution:

```go
// Create configuration
cfg := config.DefaultConfig()

// Create engine
engineConfig := &core.EngineConfig{
    Config: cfg,
    Logger: &core.DefaultLogger{},
}

engine, err := core.NewWorkflowEngine(engineConfig)
```

### Workflow Context

The workflow context carries state throughout execution:

```go
// Create workflow data
workflowData := core.NewDefaultWorkflowData()
workflowData.Set("key", "value")

// Create context
wCtx := core.NewWorkflowContext(
    context.Background(),
    "workflow-id",
    "Workflow Name",
    workflowData,
    nil, // metadata (optional)
)
```

### Steps

Steps are the building blocks of workflows. Use `FunctionStep` for simple function-based steps:

```go
step := core.NewFunctionStep(
    "step_name",
    "Step description",
    func(ctx *core.WorkflowContext) (*string, error) {
        // Step logic here
        return nil, nil // Return next step name or nil
    },
)
```

### Conditional Steps

For branching logic, use `ConditionalStep`:

```go
step := core.NewConditionalStep(
    "conditional_step",
    "Conditional logic",
    func(ctx *core.WorkflowContext) (bool, error) {
        // Return true/false based on condition
        return true, nil
    },
    "true_step",  // Step to execute if condition is true
    "false_step", // Step to execute if condition is false
)
```

### Middleware

Middleware provides cross-cutting concerns like logging, timing, and validation:

```go
type CustomMiddleware struct{}

func (m *CustomMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
    // Pre-processing
    result, err := next(ctx)
    // Post-processing
    return result, err
}

// Add to engine
engine.AddMiddleware(&CustomMiddleware{})
```

### Data Access

Safely access workflow data using helper functions:

```go
// Get string value
value, err := core.GetString(workflowCtx.Data, "key")

// Get numeric values
intVal, err := core.GetInt(workflowCtx.Data, "number")
floatVal, err := core.GetFloat64(workflowCtx.Data, "amount")

// Get boolean value
boolVal, err := core.GetBool(workflowCtx.Data, "flag")
```

### Step Results

Store and retrieve step execution results:

```go
// Store step result
workflowCtx.SetStepResult("step_name", map[string]interface{}{
    "processed_at": time.Now(),
    "success": true,
})

// Retrieve step result
result, exists := workflowCtx.GetStepResult("step_name")

// Get all step results
allResults := workflowCtx.GetAllStepResults()
```

### Metadata and Tags

Use metadata for workflow categorization and metrics:

```go
metadata := core.NewDefaultWorkflowMetadata()

// Add tags
metadata.AddTag("production")
metadata.AddTag("payment")

// Set custom fields
metadata.SetCustomField("priority", "high")
metadata.SetCustomField("department", "finance")

// Set execution metrics
metadata.SetExecutionMetric("start_time", time.Now())
```

## Best Practices

### 1. Error Handling

Always handle errors gracefully in step functions:

```go
func myStep(ctx *core.WorkflowContext) (*string, error) {
    value, err := core.GetString(ctx.Data, "required_field")
    if err != nil {
        return nil, fmt.Errorf("validation failed: %v", err)
    }
    
    // Process value...
    
    return nil, nil
}
```

### 2. Step Naming

Use descriptive names for steps and maintain consistency:

```go
// Good
core.NewFunctionStep("validate_payment_details", "Validates payment information", validatePayment)
core.NewFunctionStep("process_payment", "Processes the payment transaction", processPayment)
core.NewFunctionStep("send_confirmation", "Sends payment confirmation email", sendConfirmation)

// Avoid
core.NewFunctionStep("step1", "Step 1", step1)
core.NewFunctionStep("step2", "Step 2", step2)
```

### 3. Data Validation

Validate input data early in the workflow:

```go
func validateInputStep(ctx *core.WorkflowContext) (*string, error) {
    // Validate all required fields
    requiredFields := []string{"user_id", "amount", "currency"}
    
    for _, field := range requiredFields {
        if !ctx.Data.Has(field) {
            return nil, fmt.Errorf("required field missing: %s", field)
        }
    }
    
    return nil, nil
}
```

### 4. Logging and Monitoring

Use middleware for consistent logging and monitoring:

```go
type LoggingMiddleware struct {
    Logger core.Logger
}

func (lm *LoggingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
    start := time.Now()
    
    lm.Logger.Info("Step started", map[string]interface{}{
        "workflow_id": ctx.GetWorkflowID(),
        "step_order": ctx.GetStepOrder(),
    })
    
    result, err := next(ctx)
    duration := time.Since(start)
    
    if err != nil {
        lm.Logger.Error("Step failed", map[string]interface{}{
            "workflow_id": ctx.GetWorkflowID(),
            "step_order": ctx.GetStepOrder(),
            "duration_ms": duration.Milliseconds(),
            "error": err.Error(),
        })
    } else {
        lm.Logger.Info("Step completed", map[string]interface{}{
            "workflow_id": ctx.GetWorkflowID(),
            "step_order": ctx.GetStepOrder(),
            "duration_ms": duration.Milliseconds(),
        })
    }
    
    return result, err
}
```

### 5. Configuration Management

Use configuration for environment-specific settings:

```go
cfg := config.DefaultConfig()

// Adjust timeouts based on environment
if os.Getenv("ENV") == "production" {
    cfg.Engine.StepTimeout = 60 * time.Second
    cfg.Engine.WorkflowTimeout = 10 * time.Minute
} else {
    cfg.Engine.StepTimeout = 30 * time.Second
    cfg.Engine.WorkflowTimeout = 5 * time.Minute
}

cfg.Engine.EnableMetrics = true
cfg.Engine.EnableTracing = true
```

## Common Patterns

### Sequential Processing

```go
steps := []core.Step{
    core.NewFunctionStep("step1", "First step", step1),
    core.NewFunctionStep("step2", "Second step", step2),
    core.NewFunctionStep("step3", "Third step", step3),
}
```

### Conditional Branching

```go
steps := []core.Step{
    core.NewFunctionStep("validate", "Validate input", validate),
    core.NewConditionalStep("check_type", "Check type", checkType, "process_a", "process_b"),
    core.NewFunctionStep("process_a", "Process type A", processA),
    core.NewFunctionStep("process_b", "Process type B", processB),
    core.NewFunctionStep("finalize", "Finalize", finalize),
}
```

### Error Recovery

```go
func resilientStep(ctx *core.WorkflowContext) (*string, error) {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := performOperation(ctx)
        if err == nil {
            return nil, nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
        }
    }
    
    return nil, fmt.Errorf("operation failed after %d retries", maxRetries)
}
```

## Next Steps

1. **Explore the examples**: Run both basic and advanced examples to understand the concepts
2. **Read the documentation**: Check the main README and package documentation
3. **Build your workflow**: Start with a simple workflow and gradually add complexity
4. **Add persistence**: Integrate with storage backends for workflow persistence
5. **Implement monitoring**: Add metrics and observability to your workflows

For more information, see the main project documentation and API reference.