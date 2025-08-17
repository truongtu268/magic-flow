package messaging

import (
	"context"
	"time"

	"github.com/truongtu268/magic-flow/pkg/events"
)

// Message represents a message in the queue
type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Payload   map[string]interface{} `json:"payload"`
	Headers   map[string]string      `json:"headers"`
	Timestamp time.Time              `json:"timestamp"`
	RetryCount int                   `json:"retry_count"`
	MaxRetries int                   `json:"max_retries"`
}

// MessageQueue defines the interface for message queue operations
type MessageQueue interface {
	// Publish sends a message to a topic
	Publish(ctx context.Context, topic string, message *Message) error
	// Subscribe subscribes to a topic and processes messages
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	// Unsubscribe unsubscribes from a topic
	Unsubscribe(ctx context.Context, topic string) error
	// GetMessage retrieves a message by ID
	GetMessage(ctx context.Context, messageID string) (*Message, error)
	// AckMessage acknowledges message processing
	AckMessage(ctx context.Context, messageID string) error
	// NackMessage negatively acknowledges message processing
	NackMessage(ctx context.Context, messageID string, requeue bool) error
	// Close closes the message queue connection
	Close() error
}

// MessageHandler defines the interface for handling messages
type MessageHandler interface {
	// Handle processes a message
	Handle(ctx context.Context, message *Message) error
	// GetTopic returns the topic this handler processes
	GetTopic() string
}

// PubSubService defines the interface for pub/sub operations
type PubSubService interface {
	// Publish publishes an event
	Publish(ctx context.Context, event *events.WorkflowEvent) error
	// Subscribe subscribes to workflow events
	Subscribe(ctx context.Context, eventType events.WorkflowEventType, handler events.WorkflowEventHandler) error
	// Unsubscribe unsubscribes from workflow events
	Unsubscribe(ctx context.Context, eventType events.WorkflowEventType) error
	// PublishWorkflowEvent publishes a workflow-specific event
	PublishWorkflowEvent(ctx context.Context, workflowID string, event *events.WorkflowEvent) error
	// SubscribeToWorkflow subscribes to events for a specific workflow
	SubscribeToWorkflow(ctx context.Context, workflowID string, handler events.WorkflowEventHandler) error
	// UnsubscribeFromWorkflow unsubscribes from events for a specific workflow
	UnsubscribeFromWorkflow(ctx context.Context, workflowID string) error
	// Close closes the pub/sub service connection
	Close() error
}

// TriggerMessage represents a trigger message for waiting workflows
type TriggerMessage struct {
	TriggerKey string                 `json:"trigger_key"`
	WorkflowID string                 `json:"workflow_id"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  time.Time              `json:"timestamp"`
}

// TriggerService defines the interface for handling workflow triggers
type TriggerService interface {
	// SendTrigger sends a trigger to resume waiting workflows
	SendTrigger(ctx context.Context, trigger *TriggerMessage) error
	// RegisterTriggerHandler registers a handler for trigger processing
	RegisterTriggerHandler(triggerKey string, handler TriggerHandler) error
	// UnregisterTriggerHandler unregisters a trigger handler
	UnregisterTriggerHandler(triggerKey string) error
	// ListenForTriggers starts listening for trigger messages
	ListenForTriggers(ctx context.Context) error
	// Stop stops the trigger service
	Stop() error
}

// TriggerHandler defines the interface for handling trigger messages
type TriggerHandler interface {
	// Handle processes a trigger message
	Handle(ctx context.Context, trigger *TriggerMessage) error
	// GetTriggerKey returns the trigger key this handler processes
	GetTriggerKey() string
}

// NotificationMessage represents a notification message
type NotificationMessage struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Recipient   string                 `json:"recipient"`
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	SentAt      *time.Time             `json:"sent_at,omitempty"`
	Status      string                 `json:"status"`
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	// SendNotification sends a notification
	SendNotification(ctx context.Context, notification *NotificationMessage) error
	// ScheduleNotification schedules a notification for later delivery
	ScheduleNotification(ctx context.Context, notification *NotificationMessage, scheduleAt time.Time) error
	// CancelNotification cancels a scheduled notification
	CancelNotification(ctx context.Context, notificationID string) error
	// GetNotification retrieves a notification by ID
	GetNotification(ctx context.Context, notificationID string) (*NotificationMessage, error)
	// RegisterNotificationHandler registers a handler for notification types
	RegisterNotificationHandler(notificationType string, handler NotificationHandler) error
	// Close closes the notification service
	Close() error
}

// NotificationHandler defines the interface for handling notifications
type NotificationHandler interface {
	// Handle processes a notification
	Handle(ctx context.Context, notification *NotificationMessage) error
	// GetNotificationType returns the notification type this handler processes
	GetNotificationType() string
}

// MessagingConfig defines configuration for messaging services
type MessagingConfig struct {
	// Message Queue Configuration
	QueueURL         string        `json:"queue_url"`
	QueueType        string        `json:"queue_type"` // redis, rabbitmq, kafka, etc.
	MaxConnections   int           `json:"max_connections"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	RetryAttempts    int           `json:"retry_attempts"`
	RetryDelay       time.Duration `json:"retry_delay"`
	
	// Pub/Sub Configuration
	PubSubURL        string        `json:"pubsub_url"`
	PubSubType       string        `json:"pubsub_type"` // redis, nats, kafka, etc.
	ChannelBuffer    int           `json:"channel_buffer"`
	SubscriberTimeout time.Duration `json:"subscriber_timeout"`
	
	// Trigger Configuration
	TriggerTimeout   time.Duration `json:"trigger_timeout"`
	TriggerRetries   int           `json:"trigger_retries"`
	
	// Notification Configuration
	NotificationURL  string        `json:"notification_url"`
	NotificationType string        `json:"notification_type"` // email, sms, webhook, etc.
	BatchSize        int           `json:"batch_size"`
	BatchTimeout     time.Duration `json:"batch_timeout"`
	
	// General Configuration
	EnableMetrics    bool          `json:"enable_metrics"`
	EnableTracing    bool          `json:"enable_tracing"`
	LogLevel         string        `json:"log_level"`
}

// DefaultMessagingConfig returns default messaging configuration
func DefaultMessagingConfig() *MessagingConfig {
	return &MessagingConfig{
		QueueType:         "redis",
		MaxConnections:    10,
		ConnectionTimeout: 30 * time.Second,
		RetryAttempts:     3,
		RetryDelay:        1 * time.Second,
		PubSubType:        "redis",
		ChannelBuffer:     100,
		SubscriberTimeout: 30 * time.Second,
		TriggerTimeout:    10 * time.Second,
		TriggerRetries:    3,
		NotificationType:  "webhook",
		BatchSize:         10,
		BatchTimeout:      5 * time.Second,
		EnableMetrics:     true,
		EnableTracing:     false,
		LogLevel:          "info",
	}
}

// MessageQueueFactory defines the interface for creating message queue instances
type MessageQueueFactory interface {
	// CreateMessageQueue creates a new message queue instance
	CreateMessageQueue(config *MessagingConfig) (MessageQueue, error)
	// GetSupportedTypes returns supported message queue types
	GetSupportedTypes() []string
}

// PubSubFactory defines the interface for creating pub/sub service instances
type PubSubFactory interface {
	// CreatePubSubService creates a new pub/sub service instance
	CreatePubSubService(config *MessagingConfig) (PubSubService, error)
	// GetSupportedTypes returns supported pub/sub types
	GetSupportedTypes() []string
}