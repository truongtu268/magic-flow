package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"magic-flow/v2/pkg/models"
)

// HTTPExecutor executes HTTP requests
type HTTPExecutor struct {
	client *resty.Client
	logger *logrus.Logger
}

// NewHTTPExecutor creates a new HTTP executor
func NewHTTPExecutor(logger *logrus.Logger) *HTTPExecutor {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	return &HTTPExecutor{
		client: client,
		logger: logger,
	}
}

func (e *HTTPExecutor) Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Extract HTTP configuration
	config, ok := step.Config["http"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid HTTP configuration")
	}

	url, ok := config["url"].(string)
	if !ok {
		return nil, fmt.Errorf("URL is required for HTTP step")
	}

	method := "GET"
	if m, ok := config["method"].(string); ok {
		method = strings.ToUpper(m)
	}

	// Prepare request
	req := e.client.R().SetContext(ctx)

	// Set headers
	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.SetHeader(key, strValue)
			}
		}
	}

	// Set query parameters
	if params, ok := config["params"].(map[string]interface{}); ok {
		for key, value := range params {
			req.SetQueryParam(key, fmt.Sprintf("%v", value))
		}
	}

	// Set body for POST/PUT/PATCH requests
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if body, ok := config["body"]; ok {
			req.SetBody(body)
		} else if len(input) > 0 {
			req.SetBody(input)
		}
	}

	// Execute request
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(url)
	case "POST":
		resp, err = req.Post(url)
	case "PUT":
		resp, err = req.Put(url)
	case "PATCH":
		resp, err = req.Patch(url)
	case "DELETE":
		resp, err = req.Delete(url)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Check status code
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Parse response
	result := map[string]interface{}{
		"status_code": resp.StatusCode(),
		"headers":     resp.Header(),
		"body":        resp.String(),
	}

	// Try to parse JSON response
	if strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		var jsonBody interface{}
		if err := json.Unmarshal(resp.Body(), &jsonBody); err == nil {
			result["json"] = jsonBody
		}
	}

	e.logger.WithFields(logrus.Fields{
		"step_id":     step.ID,
		"method":      method,
		"url":         url,
		"status_code": resp.StatusCode(),
		"duration":    resp.Time().Milliseconds(),
	}).Info("HTTP request completed")

	return result, nil
}

func (e *HTTPExecutor) Validate(step *models.WorkflowStep) error {
	config, ok := step.Config["http"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("HTTP configuration is required")
	}

	if _, ok := config["url"].(string); !ok {
		return fmt.Errorf("URL is required for HTTP step")
	}

	return nil
}

func (e *HTTPExecutor) GetType() string {
	return "http"
}

// ScriptExecutor executes shell scripts
type ScriptExecutor struct {
	logger *logrus.Logger
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor(logger *logrus.Logger) *ScriptExecutor {
	return &ScriptExecutor{
		logger: logger,
	}
}

func (e *ScriptExecutor) Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Extract script configuration
	config, ok := step.Config["script"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid script configuration")
	}

	script, ok := config["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command is required for script step")
	}

	// Get shell (default to bash)
	shell := "bash"
	if s, ok := config["shell"].(string); ok {
		shell = s
	}

	// Get working directory
	workDir := ""
	if wd, ok := config["working_directory"].(string); ok {
		workDir = wd
	}

	// Prepare command
	cmd := exec.CommandContext(ctx, shell, "-c", script)
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Set environment variables
	if env, ok := config["environment"].(map[string]interface{}); ok {
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", key, value))
		}
	}

	// Add input as environment variables
	for key, value := range input {
		cmd.Env = append(cmd.Env, fmt.Sprintf("INPUT_%s=%v", strings.ToUpper(key), value))
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	e.logger.WithFields(logrus.Fields{
		"step_id":   step.ID,
		"command":   script,
		"shell":     shell,
		"work_dir":  workDir,
	}).Info("Executing script")

	// Execute command
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := map[string]interface{}{
		"stdout":      stdout.String(),
		"stderr":      stderr.String(),
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result["exit_code"] = exitError.ExitCode()
		} else {
			result["exit_code"] = -1
		}
		return result, fmt.Errorf("script execution failed: %w", err)
	}

	result["exit_code"] = 0

	e.logger.WithFields(logrus.Fields{
		"step_id":  step.ID,
		"duration": duration.Milliseconds(),
	}).Info("Script execution completed")

	return result, nil
}

func (e *ScriptExecutor) Validate(step *models.WorkflowStep) error {
	config, ok := step.Config["script"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("script configuration is required")
	}

	if _, ok := config["command"].(string); !ok {
		return fmt.Errorf("command is required for script step")
	}

	return nil
}

func (e *ScriptExecutor) GetType() string {
	return "script"
}

// TransformExecutor executes data transformations
type TransformExecutor struct {
	logger *logrus.Logger
}

// NewTransformExecutor creates a new transform executor
func NewTransformExecutor(logger *logrus.Logger) *TransformExecutor {
	return &TransformExecutor{
		logger: logger,
	}
}

func (e *TransformExecutor) Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Extract transform configuration
	config, ok := step.Config["transform"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid transform configuration")
	}

	transformType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("transform type is required")
	}

	switch transformType {
	case "json":
		return e.executeJSONTransform(config, input)
	case "filter":
		return e.executeFilter(config, input)
	case "map":
		return e.executeMap(config, input)
	case "aggregate":
		return e.executeAggregate(config, input)
	default:
		return nil, fmt.Errorf("unsupported transform type: %s", transformType)
	}
}

func (e *TransformExecutor) executeJSONTransform(config map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple JSON path extraction/transformation
	operations, ok := config["operations"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("operations are required for JSON transform")
	}

	result := make(map[string]interface{})
	for key, value := range input {
		result[key] = value
	}

	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}

		action, ok := opMap["action"].(string)
		if !ok {
			continue
		}

		switch action {
		case "extract":
			if path, ok := opMap["path"].(string); ok {
				if target, ok := opMap["target"].(string); ok {
					if value, exists := result[path]; exists {
						result[target] = value
					}
				}
			}
		case "remove":
			if path, ok := opMap["path"].(string); ok {
				delete(result, path)
			}
		case "rename":
			if from, ok := opMap["from"].(string); ok {
				if to, ok := opMap["to"].(string); ok {
					if value, exists := result[from]; exists {
						result[to] = value
						delete(result, from)
					}
				}
			}
		}
	}

	return result, nil
}

func (e *TransformExecutor) executeFilter(config map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple filtering based on conditions
	conditions, ok := config["conditions"].([]interface{})
	if !ok {
		return input, nil
	}

	result := make(map[string]interface{})

	for key, value := range input {
		include := true
		for _, cond := range conditions {
			condMap, ok := cond.(map[string]interface{})
			if !ok {
				continue
			}

			field, ok := condMap["field"].(string)
			if !ok || field != key {
				continue
			}

			operator, ok := condMap["operator"].(string)
			if !ok {
				continue
			}

			expected := condMap["value"]

			switch operator {
			case "equals":
				if value != expected {
					include = false
				}
			case "not_equals":
				if value == expected {
					include = false
				}
			case "exists":
				// Field exists, so include it
			case "not_exists":
				include = false
			}
		}

		if include {
			result[key] = value
		}
	}

	return result, nil
}

func (e *TransformExecutor) executeMap(config map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple field mapping
	mapping, ok := config["mapping"].(map[string]interface{})
	if !ok {
		return input, nil
	}

	result := make(map[string]interface{})

	for targetField, sourceField := range mapping {
		if sourceStr, ok := sourceField.(string); ok {
			if value, exists := input[sourceStr]; exists {
				result[targetField] = value
			}
		}
	}

	return result, nil
}

func (e *TransformExecutor) executeAggregate(config map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple aggregation operations
	operations, ok := config["operations"].([]interface{})
	if !ok {
		return input, nil
	}

	result := make(map[string]interface{})

	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}

		function, ok := opMap["function"].(string)
		if !ok {
			continue
		}

		field, ok := opMap["field"].(string)
		if !ok {
			continue
		}

		target, ok := opMap["target"].(string)
		if !ok {
			target = field + "_" + function
		}

		value, exists := input[field]
		if !exists {
			continue
		}

		switch function {
		case "count":
			result[target] = 1
		case "sum":
			if numValue, ok := value.(float64); ok {
				result[target] = numValue
			}
		case "avg":
			if numValue, ok := value.(float64); ok {
				result[target] = numValue
			}
		case "min":
			result[target] = value
		case "max":
			result[target] = value
		}
	}

	return result, nil
}

func (e *TransformExecutor) Validate(step *models.WorkflowStep) error {
	config, ok := step.Config["transform"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("transform configuration is required")
	}

	if _, ok := config["type"].(string); !ok {
		return fmt.Errorf("transform type is required")
	}

	return nil
}

func (e *TransformExecutor) GetType() string {
	return "transform"
}

// DelayExecutor executes delay/wait steps
type DelayExecutor struct {
	logger *logrus.Logger
}

// NewDelayExecutor creates a new delay executor
func NewDelayExecutor(logger *logrus.Logger) *DelayExecutor {
	return &DelayExecutor{
		logger: logger,
	}
}

func (e *DelayExecutor) Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Extract delay configuration
	config, ok := step.Config["delay"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid delay configuration")
	}

	durationValue, ok := config["duration"]
	if !ok {
		return nil, fmt.Errorf("duration is required for delay step")
	}

	var duration time.Duration
	var err error

	switch v := durationValue.(type) {
	case string:
		duration, err = time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid duration format: %w", err)
		}
	case int:
		duration = time.Duration(v) * time.Second
	case float64:
		duration = time.Duration(v) * time.Second
	default:
		return nil, fmt.Errorf("invalid duration type")
	}

	e.logger.WithFields(logrus.Fields{
		"step_id":  step.ID,
		"duration": duration.String(),
	}).Info("Starting delay")

	// Wait for the specified duration
	select {
	case <-time.After(duration):
		e.logger.WithFields(logrus.Fields{
			"step_id":  step.ID,
			"duration": duration.String(),
		}).Info("Delay completed")
		return map[string]interface{}{
			"duration": duration.String(),
			"waited":   true,
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (e *DelayExecutor) Validate(step *models.WorkflowStep) error {
	config, ok := step.Config["delay"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("delay configuration is required")
	}

	if _, ok := config["duration"]; !ok {
		return fmt.Errorf("duration is required for delay step")
	}

	return nil
}

func (e *DelayExecutor) GetType() string {
	return "delay"
}

// ConditionalExecutor executes conditional logic
type ConditionalExecutor struct {
	logger *logrus.Logger
}

// NewConditionalExecutor creates a new conditional executor
func NewConditionalExecutor(logger *logrus.Logger) *ConditionalExecutor {
	return &ConditionalExecutor{
		logger: logger,
	}
}

func (e *ConditionalExecutor) Execute(ctx context.Context, step *models.WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Extract conditional configuration
	config, ok := step.Config["conditional"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid conditional configuration")
	}

	condition, ok := config["condition"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("condition is required for conditional step")
	}

	// Evaluate condition
	result := e.evaluateCondition(condition, input)

	e.logger.WithFields(logrus.Fields{
		"step_id":   step.ID,
		"condition": condition,
		"result":    result,
	}).Info("Conditional evaluation completed")

	return map[string]interface{}{
		"condition_result": result,
		"input":            input,
	}, nil
}

func (e *ConditionalExecutor) evaluateCondition(condition map[string]interface{}, input map[string]interface{}) bool {
	operator, ok := condition["operator"].(string)
	if !ok {
		return false
	}

	switch operator {
	case "equals":
		field, ok := condition["field"].(string)
		if !ok {
			return false
		}
		expected := condition["value"]
		actual, exists := input[field]
		return exists && actual == expected

	case "not_equals":
		field, ok := condition["field"].(string)
		if !ok {
			return false
		}
		expected := condition["value"]
		actual, exists := input[field]
		return !exists || actual != expected

	case "exists":
		field, ok := condition["field"].(string)
		if !ok {
			return false
		}
		_, exists := input[field]
		return exists

	case "not_exists":
		field, ok := condition["field"].(string)
		if !ok {
			return false
		}
		_, exists := input[field]
		return !exists

	case "and":
		conditions, ok := condition["conditions"].([]interface{})
		if !ok {
			return false
		}
		for _, cond := range conditions {
			if condMap, ok := cond.(map[string]interface{}); ok {
				if !e.evaluateCondition(condMap, input) {
					return false
				}
			}
		}
		return true

	case "or":
		conditions, ok := condition["conditions"].([]interface{})
		if !ok {
			return false
		}
		for _, cond := range conditions {
			if condMap, ok := cond.(map[string]interface{}); ok {
				if e.evaluateCondition(condMap, input) {
					return true
				}
			}
		}
		return false

	default:
		return false
	}
}

func (e *ConditionalExecutor) Validate(step *models.WorkflowStep) error {
	config, ok := step.Config["conditional"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("conditional configuration is required")
	}

	if _, ok := config["condition"].(map[string]interface{}); !ok {
		return fmt.Errorf("condition is required for conditional step")
	}

	return nil
}

func (e *ConditionalExecutor) GetType() string {
	return "conditional"
}