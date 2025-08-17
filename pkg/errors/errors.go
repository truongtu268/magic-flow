package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents different types of errors in the system
type ErrorCode string

const (
	// Workflow errors
	ErrWorkflowNotFound     ErrorCode = "WORKFLOW_NOT_FOUND"
	ErrWorkflowInvalidState ErrorCode = "WORKFLOW_INVALID_STATE"
	ErrWorkflowTimeout      ErrorCode = "WORKFLOW_TIMEOUT"
	ErrWorkflowCancelled    ErrorCode = "WORKFLOW_CANCELLED"
	ErrWorkflowFailed       ErrorCode = "WORKFLOW_FAILED"
	ErrWorkflowExists       ErrorCode = "WORKFLOW_EXISTS"
	
	// Step errors
	ErrStepNotFound      ErrorCode = "STEP_NOT_FOUND"
	ErrStepFailed        ErrorCode = "STEP_FAILED"
	ErrStepTimeout       ErrorCode = "STEP_TIMEOUT"
	ErrStepInvalidInput  ErrorCode = "STEP_INVALID_INPUT"
	ErrStepInvalidOutput ErrorCode = "STEP_INVALID_OUTPUT"
	ErrStepPanic         ErrorCode = "STEP_PANIC"
	
	// Storage errors
	ErrStorageConnection ErrorCode = "STORAGE_CONNECTION"
	ErrStorageQuery      ErrorCode = "STORAGE_QUERY"
	ErrStorageTransaction ErrorCode = "STORAGE_TRANSACTION"
	ErrStorageNotFound   ErrorCode = "STORAGE_NOT_FOUND"
	ErrStorageConflict   ErrorCode = "STORAGE_CONFLICT"
	
	// Messaging errors
	ErrMessagingConnection ErrorCode = "MESSAGING_CONNECTION"
	ErrMessagingPublish    ErrorCode = "MESSAGING_PUBLISH"
	ErrMessagingSubscribe  ErrorCode = "MESSAGING_SUBSCRIBE"
	ErrMessagingTimeout    ErrorCode = "MESSAGING_TIMEOUT"
	
	// Configuration errors
	ErrConfigInvalid   ErrorCode = "CONFIG_INVALID"
	ErrConfigMissing   ErrorCode = "CONFIG_MISSING"
	ErrConfigLoad      ErrorCode = "CONFIG_LOAD"
	ErrConfigSave      ErrorCode = "CONFIG_SAVE"
	
	// Validation errors
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrValidationRule   ErrorCode = "VALIDATION_RULE"
	
	// Security errors
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrAuthFailed   ErrorCode = "AUTH_FAILED"
	
	// System errors
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrNotSupported ErrorCode = "NOT_SUPPORTED"
	ErrRateLimit    ErrorCode = "RATE_LIMIT"
	ErrResourceLimit ErrorCode = "RESOURCE_LIMIT"
)

// Severity represents the severity level of an error
type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// MagicFlowError represents a structured error in the magic-flow system
type MagicFlowError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Cause     error                  `json:"cause,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Severity  Severity               `json:"severity"`
	StackTrace string                `json:"stack_trace,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (e *MagicFlowError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *MagicFlowError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error code
func (e *MagicFlowError) Is(target error) bool {
	if targetErr, ok := target.(*MagicFlowError); ok {
		return e.Code == targetErr.Code
	}
	return false
}

// WithDetail adds a detail to the error
func (e *MagicFlowError) WithDetail(key string, value interface{}) *MagicFlowError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithContext adds context information to the error
func (e *MagicFlowError) WithContext(key string, value interface{}) *MagicFlowError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithCause sets the underlying cause of the error
func (e *MagicFlowError) WithCause(cause error) *MagicFlowError {
	e.Cause = cause
	return e
}

// WithSeverity sets the severity of the error
func (e *MagicFlowError) WithSeverity(severity Severity) *MagicFlowError {
	e.Severity = severity
	return e
}

// WithStackTrace captures the current stack trace
func (e *MagicFlowError) WithStackTrace() *MagicFlowError {
	e.StackTrace = captureStackTrace()
	return e
}

// New creates a new MagicFlowError
func New(code ErrorCode, message string) *MagicFlowError {
	return &MagicFlowError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  SeverityMedium,
	}
}

// Newf creates a new MagicFlowError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *MagicFlowError {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap wraps an existing error with a MagicFlowError
func Wrap(code ErrorCode, message string, cause error) *MagicFlowError {
	return &MagicFlowError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Severity:  SeverityMedium,
	}
}

// Wrapf wraps an existing error with a formatted MagicFlowError
func Wrapf(code ErrorCode, cause error, format string, args ...interface{}) *MagicFlowError {
	return Wrap(code, fmt.Sprintf(format, args...), cause)
}

// Is checks if an error is of a specific error code
func Is(err error, code ErrorCode) bool {
	if magicErr, ok := err.(*MagicFlowError); ok {
		return magicErr.Code == code
	}
	return false
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	if magicErr, ok := err.(*MagicFlowError); ok {
		return magicErr.Code
	}
	return ErrInternal
}

// GetSeverity extracts the severity from an error
func GetSeverity(err error) Severity {
	if magicErr, ok := err.(*MagicFlowError); ok {
		return magicErr.Severity
	}
	return SeverityMedium
}

// GetDetails extracts the details from an error
func GetDetails(err error) map[string]interface{} {
	if magicErr, ok := err.(*MagicFlowError); ok {
		return magicErr.Details
	}
	return nil
}

// GetContext extracts the context from an error
func GetContext(err error) map[string]interface{} {
	if magicErr, ok := err.(*MagicFlowError); ok {
		return magicErr.Context
	}
	return nil
}

// Common error constructors

// NewWorkflowNotFoundError creates a workflow not found error
func NewWorkflowNotFoundError(workflowID string) *MagicFlowError {
	return New(ErrWorkflowNotFound, "workflow not found").
		WithDetail("workflow_id", workflowID).
		WithSeverity(SeverityMedium)
}

// NewWorkflowTimeoutError creates a workflow timeout error
func NewWorkflowTimeoutError(workflowID string, timeout time.Duration) *MagicFlowError {
	return New(ErrWorkflowTimeout, "workflow execution timed out").
		WithDetail("workflow_id", workflowID).
		WithDetail("timeout", timeout.String()).
		WithSeverity(SeverityHigh)
}

// NewStepFailedError creates a step failed error
func NewStepFailedError(stepName string, cause error) *MagicFlowError {
	return Wrap(ErrStepFailed, "step execution failed", cause).
		WithDetail("step_name", stepName).
		WithSeverity(SeverityHigh)
}

// NewStepTimeoutError creates a step timeout error
func NewStepTimeoutError(stepName string, timeout time.Duration) *MagicFlowError {
	return New(ErrStepTimeout, "step execution timed out").
		WithDetail("step_name", stepName).
		WithDetail("timeout", timeout.String()).
		WithSeverity(SeverityHigh)
}

// NewStepPanicError creates a step panic error
func NewStepPanicError(stepName string, panicValue interface{}) *MagicFlowError {
	return New(ErrStepPanic, "step panicked during execution").
		WithDetail("step_name", stepName).
		WithDetail("panic_value", panicValue).
		WithSeverity(SeverityCritical).
		WithStackTrace()
}

// NewValidationError creates a validation error
func NewValidationError(field string, message string) *MagicFlowError {
	return New(ErrValidationFailed, "validation failed").
		WithDetail("field", field).
		WithDetail("validation_message", message).
		WithSeverity(SeverityMedium)
}

// NewStorageError creates a storage error
func NewStorageError(operation string, cause error) *MagicFlowError {
	return Wrap(ErrStorageQuery, "storage operation failed", cause).
		WithDetail("operation", operation).
		WithSeverity(SeverityHigh)
}

// NewMessagingError creates a messaging error
func NewMessagingError(operation string, cause error) *MagicFlowError {
	return Wrap(ErrMessagingConnection, "messaging operation failed", cause).
		WithDetail("operation", operation).
		WithSeverity(SeverityHigh)
}

// NewConfigError creates a configuration error
func NewConfigError(configType string, message string) *MagicFlowError {
	return New(ErrConfigInvalid, "configuration error").
		WithDetail("config_type", configType).
		WithDetail("config_message", message).
		WithSeverity(SeverityHigh)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(resource string) *MagicFlowError {
	return New(ErrUnauthorized, "access denied").
		WithDetail("resource", resource).
		WithSeverity(SeverityHigh)
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(limit int, window time.Duration) *MagicFlowError {
	return New(ErrRateLimit, "rate limit exceeded").
		WithDetail("limit", limit).
		WithDetail("window", window.String()).
		WithSeverity(SeverityMedium)
}

// ErrorList represents a collection of errors
type ErrorList struct {
	Errors []*MagicFlowError `json:"errors"`
}

// Error implements the error interface for ErrorList
func (el *ErrorList) Error() string {
	if len(el.Errors) == 0 {
		return "no errors"
	}
	
	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}
	
	var messages []string
	for _, err := range el.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple errors: [%s]", strings.Join(messages, "; "))
}

// Add adds an error to the list
func (el *ErrorList) Add(err *MagicFlowError) {
	el.Errors = append(el.Errors, err)
}

// AddError adds a generic error to the list
func (el *ErrorList) AddError(err error) {
	if magicErr, ok := err.(*MagicFlowError); ok {
		el.Add(magicErr)
	} else {
		el.Add(Wrap(ErrInternal, "unknown error", err))
	}
}

// HasErrors returns true if the list contains any errors
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// Count returns the number of errors in the list
func (el *ErrorList) Count() int {
	return len(el.Errors)
}

// GetBySeverity returns errors filtered by severity
func (el *ErrorList) GetBySeverity(severity Severity) []*MagicFlowError {
	var filtered []*MagicFlowError
	for _, err := range el.Errors {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// GetByCode returns errors filtered by error code
func (el *ErrorList) GetByCode(code ErrorCode) []*MagicFlowError {
	var filtered []*MagicFlowError
	for _, err := range el.Errors {
		if err.Code == code {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		Errors: make([]*MagicFlowError, 0),
	}
}

// Helper functions

func captureStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	
	var frames []string
	for i := 0; i < n; i++ {
		pc := pcs[i]
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		file, line := fn.FileLine(pc)
		frames = append(frames, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	
	return strings.Join(frames, "\n")
}

// Recovery function for handling panics
func RecoverToError(stepName string) *MagicFlowError {
	if r := recover(); r != nil {
		return NewStepPanicError(stepName, r)
	}
	return nil
}

// SafeExecute executes a function and converts panics to errors
func SafeExecute(stepName string, fn func() error) error {
	defer func() {
		if err := RecoverToError(stepName); err != nil {
			panic(err) // Re-panic with structured error
		}
	}()
	
	return fn()
}