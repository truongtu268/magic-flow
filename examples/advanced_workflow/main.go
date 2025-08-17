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
	// Create configuration with custom settings
	cfg := config.DefaultConfig()
	cfg.Engine.StepTimeout = 30 * time.Second
	cfg.Engine.WorkflowTimeout = 5 * time.Minute
	cfg.Engine.EnableMetrics = true

	// Create workflow engine with custom logger
	logger := &CustomLogger{}
	engineConfig := &core.EngineConfig{
		Config: cfg,
		Logger: logger,
	}

	engine, err := core.NewWorkflowEngine(engineConfig)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	// Add custom middleware
	engine.AddMiddleware(&TimingMiddleware{Logger: logger})
	engine.AddMiddleware(&ValidationMiddleware{})

	// Create workflow context with metadata
	workflowData := core.NewDefaultWorkflowData()
	workflowData.Set("order_id", "ORD-12345")
	workflowData.Set("customer_id", "CUST-67890")
	workflowData.Set("amount", 99.99)
	workflowData.Set("currency", "USD")

	metadata := core.NewDefaultWorkflowMetadata()
	metadata.AddTag("payment")
	metadata.AddTag("production")
	metadata.SetCustomField("priority", "high")

	wCtx := core.NewWorkflowContext(
		context.Background(),
		"payment-workflow-001",
		"Payment Processing Workflow",
		workflowData,
		metadata,
	)

	// Define workflow steps with conditional logic
	steps := []core.Step{
		core.NewFunctionStep("validate_payment", "Validates payment details", validatePaymentStep),
		core.NewConditionalStep("check_amount", "Checks payment amount", checkAmountCondition, "process_large_payment", "process_small_payment"),
		core.NewFunctionStep("process_large_payment", "Processes large payments", processLargePaymentStep),
		core.NewFunctionStep("process_small_payment", "Processes small payments", processSmallPaymentStep),
		core.NewFunctionStep("send_confirmation", "Sends payment confirmation", sendConfirmationStep),
	}

	// Execute workflow
	ctx := context.Background()
	err = engine.Execute(ctx, wCtx.WorkflowID, steps, wCtx.Data)
	if err != nil {
		log.Fatalf("Workflow execution failed: %v", err)
	}

	fmt.Printf("\n=== Workflow Completed Successfully! ===\n")
	fmt.Printf("Workflow ID: %s\n", wCtx.WorkflowID)
	fmt.Printf("Start time: %v\n", wCtx.StartTime)
	fmt.Printf("Final status: %v\n", wCtx.Status)
	fmt.Printf("Tags: %v\n", wCtx.Metadata.GetTags())

	// Print step results
	fmt.Printf("\n=== Step Results ===\n")
	for stepName, result := range wCtx.GetAllStepResults() {
		fmt.Printf("%s: %v\n", stepName, result)
	}
}

// Custom Logger implementation
type CustomLogger struct{}

func (l *CustomLogger) Debug(message string, fields map[string]interface{}) {
	fmt.Printf("[DEBUG] %s - %v\n", message, fields)
}

func (l *CustomLogger) Info(message string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s - %v\n", message, fields)
}

func (l *CustomLogger) Warn(message string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s - %v\n", message, fields)
}

func (l *CustomLogger) Error(message string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s - %v\n", message, fields)
}

// Custom Middleware implementations
type TimingMiddleware struct {
	Logger core.Logger
}

func (tm *TimingMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	start := time.Now()
	result, err := next(ctx)
	duration := time.Since(start)

	tm.Logger.Info("Step timing", map[string]interface{}{
		"workflow_id": ctx.GetWorkflowID(),
		"step_order": ctx.GetStepOrder(),
		"duration_ms": duration.Milliseconds(),
	})

	return result, err
}

type ValidationMiddleware struct{}

func (vm *ValidationMiddleware) Handle(ctx *core.WorkflowContext, next core.StepHandler) (*string, error) {
	// Validate workflow context before step execution
	if ctx.GetWorkflowID() == "" {
		return nil, fmt.Errorf("workflow ID is required")
	}

	if ctx.Data == nil {
		return nil, fmt.Errorf("workflow data is required")
	}

	return next(ctx)
}

// Step handlers
func validatePaymentStep(workflowCtx *core.WorkflowContext) (*string, error) {
	orderID, err := core.GetString(workflowCtx.Data, "order_id")
	if err != nil {
		return nil, fmt.Errorf("order_id is required: %v", err)
	}

	amount, err := core.GetFloat64(workflowCtx.Data, "amount")
	if err != nil {
		return nil, fmt.Errorf("amount is required: %v", err)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	workflowCtx.SetStepResult("validate_payment", map[string]interface{}{
		"validated": true,
		"order_id": orderID,
		"amount": amount,
		"validated_at": time.Now(),
	})

	nextStep := "check_amount"
	return &nextStep, nil
}

func checkAmountCondition(workflowCtx *core.WorkflowContext) (bool, error) {
	amount, err := core.GetFloat64(workflowCtx.Data, "amount")
	if err != nil {
		return false, err
	}

	// Return true for large payments (>= 100), false for small payments
	return amount >= 100.0, nil
}

func processLargePaymentStep(workflowCtx *core.WorkflowContext) (*string, error) {
	amount, _ := core.GetFloat64(workflowCtx.Data, "amount")
	orderID, _ := core.GetString(workflowCtx.Data, "order_id")

	// Simulate additional verification for large payments
	time.Sleep(100 * time.Millisecond)

	workflowCtx.SetData("processing_fee", amount*0.03) // 3% fee for large payments
	workflowCtx.SetData("requires_approval", true)

	workflowCtx.SetStepResult("process_large_payment", map[string]interface{}{
		"payment_type": "large",
		"order_id": orderID,
		"amount": amount,
		"fee": amount * 0.03,
		"processed_at": time.Now(),
	})

	nextStep := "send_confirmation"
	return &nextStep, nil
}

func processSmallPaymentStep(workflowCtx *core.WorkflowContext) (*string, error) {
	amount, _ := core.GetFloat64(workflowCtx.Data, "amount")
	orderID, _ := core.GetString(workflowCtx.Data, "order_id")

	// Simulate faster processing for small payments
	time.Sleep(50 * time.Millisecond)

	workflowCtx.SetData("processing_fee", amount*0.01) // 1% fee for small payments
	workflowCtx.SetData("requires_approval", false)

	workflowCtx.SetStepResult("process_small_payment", map[string]interface{}{
		"payment_type": "small",
		"order_id": orderID,
		"amount": amount,
		"fee": amount * 0.01,
		"processed_at": time.Now(),
	})

	nextStep := "send_confirmation"
	return &nextStep, nil
}

func sendConfirmationStep(workflowCtx *core.WorkflowContext) (*string, error) {
	orderID, _ := core.GetString(workflowCtx.Data, "order_id")
	customerID, _ := core.GetString(workflowCtx.Data, "customer_id")
	amount, _ := core.GetFloat64(workflowCtx.Data, "amount")
	fee, _ := workflowCtx.GetData("processing_fee")

	// Simulate sending confirmation
	confirmationID := fmt.Sprintf("CONF-%d", time.Now().Unix())
	workflowCtx.SetData("confirmation_id", confirmationID)

	fmt.Printf("\n=== Payment Confirmation ===\n")
	fmt.Printf("Confirmation ID: %s\n", confirmationID)
	fmt.Printf("Order ID: %s\n", orderID)
	fmt.Printf("Customer ID: %s\n", customerID)
	fmt.Printf("Amount: $%.2f\n", amount)
	fmt.Printf("Processing Fee: $%.2f\n", fee)

	workflowCtx.SetStepResult("send_confirmation", map[string]interface{}{
		"confirmation_id": confirmationID,
		"sent_at": time.Now(),
		"success": true,
	})

	return nil, nil
}