# Magic Flow

A powerful, flexible, and extensible workflow engine library for Go that enables you to build complex business processes with ease.

## Features

- 🚀 **High Performance**: Optimized for concurrent workflow execution
- 🔧 **Modular Architecture**: Clean separation of concerns with well-defined interfaces
- 🎯 **Type Safety**: Strongly typed interfaces and comprehensive error handling
- 🔄 **Middleware Support**: Extensible middleware chain for cross-cutting concerns
- 📊 **Built-in Monitoring**: Metrics, logging, and tracing capabilities
- 🔀 **Conditional Logic**: Support for branching and conditional step execution
- 💾 **Persistence Ready**: Pluggable storage backends for workflow state
- 🔧 **Configuration Management**: Flexible configuration system
- 📝 **Comprehensive Testing**: High test coverage with unit and integration tests

## Quick Start

### Installation

```bash
go get github.com/truongtu268/magic-flow
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/truongtu268/magic-flow/pkg/config"
    "github.com/truongtu268/magic-flow/pkg/core"
)

func main() {
    // Create configuration
    cfg := config.DefaultConfig()

    // Create workflow engine
    engineConfig := &core.EngineConfig{
        Config: cfg,
        Logger: &core.DefaultLogger{},
    }

    engine, err := core.NewWorkflowEngine(engineConfig)
    if err != nil {
        log.Fatalf("Failed to create engine: %v", err)
    }

    // Create workflow context
    workflowData := core.NewDefaultWorkflowData()
    workflowData.Set("input", "Hello, Magic Flow!")

    wCtx := core.NewWorkflowContext(
        context.Background(),
        "my-workflow-001",
        "My First Workflow",
        workflowData,
        nil,
    )

    // Define workflow steps
    steps := []core.Step{
        core.NewFunctionStep("process", "Process data", func(ctx *core.WorkflowContext) (*string, error) {
            input, _ := core.GetString(ctx.Data, "input")
            result := fmt.Sprintf("Processed: %s", input)
            ctx.SetData("result", result)
            return nil, nil
        }),
    }

    // Execute workflow
    err = engine.Execute(context.Background(), wCtx.WorkflowID, steps, wCtx.Data)
    if err != nil {
        log.Fatalf("Workflow execution failed: %v", err)
    }

    fmt.Println("Workflow completed successfully!")
}
```

## Architecture

Magic Flow follows a clean, modular architecture:

```
pkg/
├── config/          # Configuration management
├── core/            # Core workflow engine and interfaces
├── errors/          # Custom error types and handling
├── events/          # Event system for workflow notifications
├── messaging/       # Message queue and pub/sub interfaces
├── recovery/        # Workflow recovery and resilience
└── storage/         # Storage interfaces and implementations
```

### Core Components

- **Workflow Engine**: Orchestrates step execution and manages workflow lifecycle
- **Steps**: Individual units of work that can be composed into workflows
- **Middleware**: Cross-cutting concerns like logging, timing, and validation
- **Context**: Carries workflow state and data throughout execution
- **Storage**: Pluggable persistence layer for workflow state
- **Events**: Notification system for workflow lifecycle events

## Examples

The `examples/` directory contains comprehensive examples:

- **Basic Workflow**: Simple linear workflow execution
- **Advanced Workflow**: Conditional logic, middleware, and complex data processing

Run the examples:

```bash
# Basic example
cd examples/basic_workflow
go run main.go

# Advanced example
cd examples/advanced_workflow
go run main.go
```

See the [Examples README](examples/README.md) for detailed documentation.

## Key Concepts

### Workflow Steps

Steps are the building blocks of workflows. Use `FunctionStep` for simple function-based steps:

```go
step := core.NewFunctionStep(
    "validate_input",
    "Validates user input",
    func(ctx *core.WorkflowContext) (*string, error) {
        // Step logic here
        return nil, nil // Return next step name or nil
    },
)
```

### Conditional Logic

Implement branching logic with `ConditionalStep`:

```go
step := core.NewConditionalStep(
    "check_amount",
    "Check payment amount",
    func(ctx *core.WorkflowContext) (bool, error) {
        amount, _ := core.GetFloat64(ctx.Data, "amount")
        return amount >= 100.0, nil
    },
    "large_payment",  // Step for true condition
    "small_payment",  // Step for false condition
)
```

### Middleware

Add cross-cutting concerns with middleware:

```go
type TimingMiddleware struct {
    Logger core.Logger
}

func (tm *TimingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
    start := time.Now()
    result, err := next(ctx)
    duration := time.Since(start)
    
    tm.Logger.Info("Step completed", map[string]interface{}{
        "duration_ms": duration.Milliseconds(),
    })
    
    return result, err
}

// Add to engine
engine.AddMiddleware(&TimingMiddleware{Logger: logger})
```

### Data Access

Safely access workflow data:

```go
// Type-safe data access
value, err := core.GetString(ctx.Data, "key")
intVal, err := core.GetInt(ctx.Data, "number")
floatVal, err := core.GetFloat64(ctx.Data, "amount")
boolVal, err := core.GetBool(ctx.Data, "flag")

// Set data
ctx.SetData("result", "processed value")

// Store step results
ctx.SetStepResult("step_name", map[string]interface{}{
    "processed_at": time.Now(),
    "success": true,
})
```

## Configuration

Configure the engine for different environments:

```go
cfg := config.DefaultConfig()

// Adjust timeouts
cfg.Engine.StepTimeout = 30 * time.Second
cfg.Engine.WorkflowTimeout = 5 * time.Minute

// Enable features
cfg.Engine.EnableMetrics = true
cfg.Engine.EnableTracing = true

// Configure logging
cfg.Logging.Level = "info"
cfg.Logging.Format = "json"
```

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/core
```

## Building

Build the library:

```bash
# Build all packages
go build ./...

# Build examples
go build ./examples/basic_workflow
go build ./examples/advanced_workflow
```

## Contributing

We welcome contributions! Please see our contributing guidelines for details.

### Development Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run tests: `go test ./...`
4. Build examples: `go build ./examples/...`

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write comprehensive tests for new features
- Document public APIs with clear examples

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [ ] Enhanced persistence backends (PostgreSQL, MongoDB, Redis)
- [ ] Distributed workflow execution
- [ ] Web UI for workflow monitoring
- [ ] GraphQL API for workflow management
- [ ] Advanced scheduling capabilities
- [ ] Workflow versioning and migration tools

## Support

For questions, issues, or contributions:

- Create an issue on GitHub
- Check the examples and documentation
- Review the test suite for usage patterns

---

**Magic Flow** - Build powerful workflows with confidence! 🚀
