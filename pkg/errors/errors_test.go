package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMagicFlowError(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "test message")
		require.NotNil(t, err)
		assert.Equal(t, ErrWorkflowNotFound, err.Code)
		assert.Equal(t, "test message", err.Message)
		assert.Equal(t, SeverityMedium, err.Severity)
		assert.False(t, err.Timestamp.IsZero())
	})

	t.Run("Newf", func(t *testing.T) {
		err := Newf(ErrStepFailed, "step %s failed with code %d", "test-step", 500)
		require.NotNil(t, err)
		assert.Equal(t, ErrStepFailed, err.Code)
		assert.Equal(t, "step test-step failed with code 500", err.Message)
	})

	t.Run("Wrap", func(t *testing.T) {
		cause := errors.New("original error")
		err := Wrap(ErrStorageQuery, "storage operation failed", cause)
		require.NotNil(t, err)
		assert.Equal(t, ErrStorageQuery, err.Code)
		assert.Equal(t, "storage operation failed", err.Message)
		assert.Equal(t, cause, err.Cause)
	})

	t.Run("Error", func(t *testing.T) {
		// Test without cause
		err := New(ErrWorkflowNotFound, "workflow not found")
		expected := "WORKFLOW_NOT_FOUND: workflow not found"
		assert.Equal(t, expected, err.Error())
		
		// Test with cause
		cause := errors.New("database connection failed")
		err = Wrap(ErrStorageConnection, "storage error", cause)
		expected = "STORAGE_CONNECTION: storage error (caused by: database connection failed)"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("original error")
		err := Wrap(ErrStorageQuery, "storage operation failed", cause)
		
		unwrapped := err.Unwrap()
		assert.Equal(t, cause, unwrapped)
		
		// Test without cause
		err = New(ErrWorkflowNotFound, "workflow not found")
		unwrapped = err.Unwrap()
		assert.Nil(t, unwrapped)
	})

	t.Run("Is", func(t *testing.T) {
		err1 := New(ErrWorkflowNotFound, "workflow not found")
		err2 := New(ErrWorkflowNotFound, "another message")
		err3 := New(ErrStepFailed, "step failed")
		
		assert.True(t, err1.Is(err2))
		assert.False(t, err1.Is(err3))
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		assert.False(t, err1.Is(stdErr))
	})

	t.Run("WithDetail", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		err.WithDetail("workflow_id", "test-123")
		err.WithDetail("user_id", "user-456")
		
		assert.Equal(t, "test-123", err.Details["workflow_id"])
		assert.Equal(t, "user-456", err.Details["user_id"])
	})

	t.Run("WithContext", func(t *testing.T) {
		err := New(ErrStepFailed, "step failed")
		err.WithContext("step_name", "validation")
		err.WithContext("retry_count", 3)
		
		assert.Equal(t, "validation", err.Context["step_name"])
		assert.Equal(t, 3, err.Context["retry_count"])
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("original error")
		err := New(ErrStepFailed, "step failed")
		err.WithCause(cause)
		
		assert.Equal(t, cause, err.Cause)
	})

	t.Run("WithSeverity", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		err.WithSeverity(SeverityHigh)
		
		assert.Equal(t, SeverityHigh, err.Severity)
	})

	t.Run("WithStackTrace", func(t *testing.T) {
		err := New(ErrStepPanic, "step panicked")
		err.WithStackTrace()
		
		assert.NotEmpty(t, err.StackTrace)
		assert.Contains(t, err.StackTrace, "WithStackTrace")
	})
}

func TestErrorHelpers(t *testing.T) {
	t.Run("Is", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		
		assert.True(t, Is(err, ErrWorkflowNotFound))
		assert.False(t, Is(err, ErrStepFailed))
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		assert.False(t, Is(stdErr, ErrWorkflowNotFound))
	})

	t.Run("GetCode", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		code := GetCode(err)
		assert.Equal(t, ErrWorkflowNotFound, code)
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		code = GetCode(stdErr)
		assert.Equal(t, ErrInternal, code)
	})

	t.Run("GetSeverity", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		err.WithSeverity(SeverityHigh)
		severity := GetSeverity(err)
		assert.Equal(t, SeverityHigh, severity)
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		severity = GetSeverity(stdErr)
		assert.Equal(t, SeverityMedium, severity)
	})

	t.Run("GetDetails", func(t *testing.T) {
		err := New(ErrWorkflowNotFound, "workflow not found")
		err.WithDetail("workflow_id", "test-123")
		details := GetDetails(err)
		assert.Equal(t, "test-123", details["workflow_id"])
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		details = GetDetails(stdErr)
		assert.Nil(t, details)
	})

	t.Run("GetContext", func(t *testing.T) {
		err := New(ErrStepFailed, "step failed")
		err.WithContext("step_name", "validation")
		context := GetContext(err)
		assert.Equal(t, "validation", context["step_name"])
		
		// Test with non-MagicFlowError
		stdErr := errors.New("standard error")
		context = GetContext(stdErr)
		assert.Nil(t, context)
	})
}

func TestCommonErrorConstructors(t *testing.T) {
	t.Run("NewWorkflowNotFoundError", func(t *testing.T) {
		err := NewWorkflowNotFoundError("test-workflow-123")
		assert.Equal(t, ErrWorkflowNotFound, err.Code)
		assert.Equal(t, "workflow not found", err.Message)
		assert.Equal(t, "test-workflow-123", err.Details["workflow_id"])
		assert.Equal(t, SeverityMedium, err.Severity)
	})

	t.Run("NewWorkflowTimeoutError", func(t *testing.T) {
		timeout := 30 * time.Second
		err := NewWorkflowTimeoutError("test-workflow-123", timeout)
		assert.Equal(t, ErrWorkflowTimeout, err.Code)
		assert.Equal(t, "workflow execution timed out", err.Message)
		assert.Equal(t, "test-workflow-123", err.Details["workflow_id"])
		assert.Equal(t, timeout.String(), err.Details["timeout"])
		assert.Equal(t, SeverityHigh, err.Severity)
	})

	t.Run("NewStepFailedError", func(t *testing.T) {
		cause := errors.New("database connection failed")
		err := NewStepFailedError("validation-step", cause)
		assert.Equal(t, ErrStepFailed, err.Code)
		assert.Equal(t, "step execution failed", err.Message)
		assert.Equal(t, "validation-step", err.Details["step_name"])
		assert.Equal(t, cause, err.Cause)
		assert.Equal(t, SeverityHigh, err.Severity)
	})

	t.Run("NewStepTimeoutError", func(t *testing.T) {
		timeout := 10 * time.Second
		err := NewStepTimeoutError("processing-step", timeout)
		assert.Equal(t, ErrStepTimeout, err.Code)
		assert.Equal(t, "step execution timed out", err.Message)
		assert.Equal(t, "processing-step", err.Details["step_name"])
		assert.Equal(t, timeout.String(), err.Details["timeout"])
		assert.Equal(t, SeverityHigh, err.Severity)
	})

	t.Run("NewStepPanicError", func(t *testing.T) {
		panicValue := "null pointer dereference"
		err := NewStepPanicError("calculation-step", panicValue)
		assert.Equal(t, ErrStepPanic, err.Code)
		assert.Equal(t, "step panicked during execution", err.Message)
		assert.Equal(t, "calculation-step", err.Details["step_name"])
		assert.Equal(t, panicValue, err.Details["panic_value"])
		assert.Equal(t, SeverityCritical, err.Severity)
		assert.NotEmpty(t, err.StackTrace)
	})

	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError("email", "invalid email format")
		assert.Equal(t, ErrValidationFailed, err.Code)
		assert.Equal(t, "validation failed", err.Message)
		assert.Equal(t, "email", err.Details["field"])
		assert.Equal(t, "invalid email format", err.Details["validation_message"])
		assert.Equal(t, SeverityMedium, err.Severity)
	})

	t.Run("NewRateLimitError", func(t *testing.T) {
		limit := 100
		window := time.Minute
		err := NewRateLimitError(limit, window)
		assert.Equal(t, ErrRateLimit, err.Code)
		assert.Equal(t, "rate limit exceeded", err.Message)
		assert.Equal(t, limit, err.Details["limit"])
		assert.Equal(t, window.String(), err.Details["window"])
		assert.Equal(t, SeverityMedium, err.Severity)
	})
}

func TestErrorList(t *testing.T) {
	t.Run("NewErrorList", func(t *testing.T) {
		el := NewErrorList()
		require.NotNil(t, el)
		assert.Equal(t, 0, el.Count())
		assert.False(t, el.HasErrors())
	})

	t.Run("Add", func(t *testing.T) {
		el := NewErrorList()
		err1 := New(ErrWorkflowNotFound, "workflow not found")
		err2 := New(ErrStepFailed, "step failed")
		
		el.Add(err1)
		assert.Equal(t, 1, el.Count())
		assert.True(t, el.HasErrors())
		
		el.Add(err2)
		assert.Equal(t, 2, el.Count())
	})

	t.Run("AddError", func(t *testing.T) {
		el := NewErrorList()
		
		// Add MagicFlowError
		magicErr := New(ErrWorkflowNotFound, "workflow not found")
		el.AddError(magicErr)
		assert.Equal(t, 1, el.Count())
		
		// Add standard error
		stdErr := errors.New("standard error")
		el.AddError(stdErr)
		assert.Equal(t, 2, el.Count())
		
		// Check that standard error was wrapped
		assert.Equal(t, ErrInternal, el.Errors[1].Code)
	})

	t.Run("Error", func(t *testing.T) {
		// Test empty list
		el := NewErrorList()
		assert.Equal(t, "no errors", el.Error())
		
		// Test single error
		err1 := New(ErrWorkflowNotFound, "workflow not found")
		el.Add(err1)
		assert.Equal(t, err1.Error(), el.Error())
		
		// Test multiple errors
		err2 := New(ErrStepFailed, "step failed")
		el.Add(err2)
		expected := "multiple errors: [WORKFLOW_NOT_FOUND: workflow not found; STEP_FAILED: step failed]"
		assert.Equal(t, expected, el.Error())
	})

	t.Run("GetBySeverity", func(t *testing.T) {
		el := NewErrorList()
		err1 := New(ErrWorkflowNotFound, "workflow not found").WithSeverity(SeverityMedium)
		err2 := New(ErrStepPanic, "step panicked").WithSeverity(SeverityCritical)
		err3 := New(ErrValidationFailed, "validation failed").WithSeverity(SeverityMedium)
		
		el.Add(err1)
		el.Add(err2)
		el.Add(err3)
		
		mediumErrors := el.GetBySeverity(SeverityMedium)
		assert.Len(t, mediumErrors, 2)
		
		criticalErrors := el.GetBySeverity(SeverityCritical)
		assert.Len(t, criticalErrors, 1)
		assert.Equal(t, err2, criticalErrors[0])
	})

	t.Run("GetByCode", func(t *testing.T) {
		el := NewErrorList()
		err1 := New(ErrWorkflowNotFound, "workflow not found")
		err2 := New(ErrStepFailed, "step failed")
		err3 := New(ErrWorkflowNotFound, "another workflow not found")
		
		el.Add(err1)
		el.Add(err2)
		el.Add(err3)
		
		workflowErrors := el.GetByCode(ErrWorkflowNotFound)
		assert.Len(t, workflowErrors, 2)
		
		stepErrors := el.GetByCode(ErrStepFailed)
		assert.Len(t, stepErrors, 1)
		assert.Equal(t, err2, stepErrors[0])
	})
}

func TestSafeExecute(t *testing.T) {
	t.Run("SuccessfulExecution", func(t *testing.T) {
		err := SafeExecute("test-step", func() error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("ErrorExecution", func(t *testing.T) {
		expectedErr := errors.New("test error")
		err := SafeExecute("test-step", func() error {
			return expectedErr
		})
		assert.Equal(t, expectedErr, err)
	})

	t.Run("PanicExecution", func(t *testing.T) {
		assert.Panics(t, func() {
			SafeExecute("test-step", func() error {
				panic("test panic")
			})
		})
	})
}

func TestRecoverToError(t *testing.T) {
	t.Run("NoPanic", func(t *testing.T) {
		func() {
			defer func() {
				err := RecoverToError("test-step")
				assert.Nil(t, err)
			}()
			// Normal execution, no panic
		}()
	})

	t.Run("WithPanic", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r != nil {
				err := NewStepPanicError("test-step", r)
				require.NotNil(t, err)
				assert.Equal(t, ErrStepPanic, err.Code)
				assert.Equal(t, "test-step", err.Details["step_name"])
				assert.Equal(t, "test panic", err.Details["panic_value"])
			}
			}()
			panic("test panic")
		}()
	})
}