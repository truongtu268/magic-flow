package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/truongtu268/magic-flow/pkg/config"
	"github.com/truongtu268/magic-flow/pkg/core"
)

func main() {
	// Create configuration
	cfg := config.DefaultConfig()

	// Create workflow engine (simplified for basic example)
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
	workflowData.Set("user_id", 12345)

	wCtx := core.NewWorkflowContext(context.Background(), "basic-workflow-001", "Basic Workflow Demo", workflowData, nil)

	// Define workflow steps
	steps := []core.Step{
		core.NewFunctionStep("validate_input", "Validates input data", validateInputStep),
		core.NewFunctionStep("process_data", "Processes the input data", processDataStep),
		core.NewFunctionStep("finalize", "Finalizes the workflow", finalizeStep),
	}

	// Execute workflow
	ctx := context.Background()
	err = engine.Execute(ctx, wCtx.WorkflowID, steps, wCtx.Data)
	if err != nil {
		log.Fatalf("Workflow execution failed: %v", err)
	}

	fmt.Printf("Workflow completed successfully!\n")
	fmt.Printf("Start time: %v\n", wCtx.StartTime)
	fmt.Printf("Final status: %v\n", wCtx.Status)
}

// Step handlers
func validateInputStep(workflowCtx *core.WorkflowContext) (*string, error) {
	input, err := core.GetString(workflowCtx.Data, "input")
	if err != nil {
		return nil, fmt.Errorf("input validation failed: %v", err)
	}

	if len(input) == 0 {
		return nil, fmt.Errorf("input cannot be empty")
	}

	workflowCtx.SetStepResult("validate_input", map[string]interface{}{
		"validated": true,
		"input_length": len(input),
	})

	return nil, nil
}

func processDataStep(workflowCtx *core.WorkflowContext) (*string, error) {
	input, _ := workflowCtx.GetData("input")
	userID, _ := workflowCtx.GetData("user_id")

	processedData := fmt.Sprintf("Processed: %v for user %v", input, userID)
	workflowCtx.SetData("processed_result", processedData)

	workflowCtx.SetStepResult("process_data", map[string]interface{}{
		"processed_at": time.Now(),
		"result_length": len(processedData),
	})

	return nil, nil
}

func finalizeStep(workflowCtx *core.WorkflowContext) (*string, error) {
	processedResult, exists := workflowCtx.GetData("processed_result")
	if !exists {
		return nil, fmt.Errorf("processed result not found")
	}

	fmt.Printf("Final result: %v\n", processedResult)

	workflowCtx.SetStepResult("finalize", map[string]interface{}{
		"completed_at": time.Now(),
		"success": true,
	})

	return nil, nil
}