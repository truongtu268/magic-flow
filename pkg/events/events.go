package events

import (
	"time"
)

// WorkflowEventType represents the type of workflow event
type WorkflowEventType string

const (
	WorkflowEventStarted      WorkflowEventType = "started"
	WorkflowEventCompleted    WorkflowEventType = "completed"
	WorkflowEventFailed       WorkflowEventType = "failed"
	WorkflowEventCancelled    WorkflowEventType = "cancelled"
	WorkflowEventPaused       WorkflowEventType = "paused"
	WorkflowEventResumed      WorkflowEventType = "resumed"
	WorkflowEventStepStarted  WorkflowEventType = "step_started"
	WorkflowEventStepCompleted WorkflowEventType = "step_completed"
	WorkflowEventStepFailed   WorkflowEventType = "step_failed"
)

// WorkflowEvent represents an event in the workflow lifecycle
type WorkflowEvent struct {
	ID         string                 `json:"id"`
	Type       WorkflowEventType      `json:"type"`
	WorkflowID string                 `json:"workflow_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

// WorkflowEventHandler is a function type for handling workflow events
type WorkflowEventHandler func(event *WorkflowEvent) error

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
	WorkflowStatusPaused    WorkflowStatus = "paused"
	WorkflowStatusUnknown   WorkflowStatus = "unknown"
)

// JobStatus represents the status of a background job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)