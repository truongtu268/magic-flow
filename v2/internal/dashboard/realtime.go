package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RealtimeManager handles real-time updates for dashboard clients
type RealtimeManager struct {
	clients    map[string]*Client
	clientsMux sync.RWMutex
	broadcast  chan RealtimeUpdate
	register   chan *Client
	unregister chan *Client
	done       chan struct{}
	wg         sync.WaitGroup
}

// Client represents a connected dashboard client
type Client struct {
	ID       string
	UserID   uuid.UUID
	Updates  chan RealtimeUpdate
	Done     chan struct{}
	LastSeen time.Time
}

// RealtimeUpdate represents a real-time update message
type RealtimeUpdate struct {
	Type      string                 `json:"type"`
	Data      interface{}            `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	ClientID  string                 `json:"client_id,omitempty"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateType constants for different types of real-time updates
const (
	UpdateTypeExecutionStarted   = "execution_started"
	UpdateTypeExecutionCompleted = "execution_completed"
	UpdateTypeExecutionFailed    = "execution_failed"
	UpdateTypeExecutionCancelled = "execution_cancelled"
	UpdateTypeStepStarted        = "step_started"
	UpdateTypeStepCompleted      = "step_completed"
	UpdateTypeStepFailed         = "step_failed"
	UpdateTypeWorkflowCreated    = "workflow_created"
	UpdateTypeWorkflowUpdated    = "workflow_updated"
	UpdateTypeWorkflowDeleted    = "workflow_deleted"
	UpdateTypeMetricsUpdated     = "metrics_updated"
	UpdateTypeAlertTriggered     = "alert_triggered"
	UpdateTypeAlertResolved      = "alert_resolved"
	UpdateTypeSystemStatus       = "system_status"
	UpdateTypeUserActivity       = "user_activity"
	UpdateTypeHeartbeat          = "heartbeat"
	UpdateTypeError              = "error"
)

// NewRealtimeManager creates a new realtime manager
func NewRealtimeManager() *RealtimeManager {
	rm := &RealtimeManager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan RealtimeUpdate, 256),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		done:       make(chan struct{}),
	}

	// Start the manager
	rm.wg.Add(1)
	go rm.run()

	return rm
}

// Subscribe subscribes a client to real-time updates
func (rm *RealtimeManager) Subscribe(ctx context.Context, clientID string) (<-chan RealtimeUpdate, error) {
	client := &Client{
		ID:       clientID,
		Updates:  make(chan RealtimeUpdate, 64),
		Done:     make(chan struct{}),
		LastSeen: time.Now(),
	}

	// Extract user ID from context if available
	if userID, ok := ctx.Value("user_id").(uuid.UUID); ok {
		client.UserID = userID
	}

	// Register the client
	select {
	case rm.register <- client:
		return client.Updates, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rm.done:
		return nil, fmt.Errorf("realtime manager is shutting down")
	}
}

// Unsubscribe unsubscribes a client from real-time updates
func (rm *RealtimeManager) Unsubscribe(clientID string) {
	rm.clientsMux.RLock()
	client, exists := rm.clients[clientID]
	rm.clientsMux.RUnlock()

	if exists {
		select {
		case rm.unregister <- client:
		case <-rm.done:
		}
	}
}

// Publish publishes a real-time update to all subscribed clients
func (rm *RealtimeManager) Publish(update RealtimeUpdate) {
	update.Timestamp = time.Now()

	select {
	case rm.broadcast <- update:
	case <-rm.done:
	}
}

// PublishToClient publishes a real-time update to a specific client
func (rm *RealtimeManager) PublishToClient(clientID string, update RealtimeUpdate) {
	update.Timestamp = time.Now()
	update.ClientID = clientID

	rm.clientsMux.RLock()
	client, exists := rm.clients[clientID]
	rm.clientsMux.RUnlock()

	if exists {
		select {
		case client.Updates <- update:
		default:
			// Client channel is full, skip this update
		}
	}
}

// PublishToUser publishes a real-time update to all clients of a specific user
func (rm *RealtimeManager) PublishToUser(userID uuid.UUID, update RealtimeUpdate) {
	update.Timestamp = time.Now()
	update.UserID = &userID

	rm.clientsMux.RLock()
	defer rm.clientsMux.RUnlock()

	for _, client := range rm.clients {
		if client.UserID == userID {
			select {
			case client.Updates <- update:
			default:
				// Client channel is full, skip this update
			}
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (rm *RealtimeManager) GetConnectedClients() int {
	rm.clientsMux.RLock()
	defer rm.clientsMux.RUnlock()
	return len(rm.clients)
}

// GetClientsByUser returns the number of clients for a specific user
func (rm *RealtimeManager) GetClientsByUser(userID uuid.UUID) int {
	rm.clientsMux.RLock()
	defer rm.clientsMux.RUnlock()

	count := 0
	for _, client := range rm.clients {
		if client.UserID == userID {
			count++
		}
	}
	return count
}

// Close closes the realtime manager and all client connections
func (rm *RealtimeManager) Close() {
	close(rm.done)
	rm.wg.Wait()
}

// run is the main loop for the realtime manager
func (rm *RealtimeManager) run() {
	defer rm.wg.Done()

	// Start heartbeat ticker
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	// Start cleanup ticker
	cleanupTicker := time.NewTicker(5 * time.Minute)
	defer cleanupTicker.Stop()

	for {
		select {
		case client := <-rm.register:
			rm.registerClient(client)

		case client := <-rm.unregister:
			rm.unregisterClient(client)

		case update := <-rm.broadcast:
			rm.broadcastUpdate(update)

		case <-heartbeatTicker.C:
			rm.sendHeartbeat()

		case <-cleanupTicker.C:
			rm.cleanupStaleClients()

		case <-rm.done:
			rm.cleanup()
			return
		}
	}
}

// registerClient registers a new client
func (rm *RealtimeManager) registerClient(client *Client) {
	rm.clientsMux.Lock()
	rm.clients[client.ID] = client
	rm.clientsMux.Unlock()

	// Send welcome message
	welcomeUpdate := RealtimeUpdate{
		Type: "connected",
		Data: map[string]interface{}{
			"client_id": client.ID,
			"message":   "Connected to real-time updates",
		},
		Timestamp: time.Now(),
	}

	select {
	case client.Updates <- welcomeUpdate:
	default:
	}
}

// unregisterClient unregisters a client
func (rm *RealtimeManager) unregisterClient(client *Client) {
	rm.clientsMux.Lock()
	delete(rm.clients, client.ID)
	rm.clientsMux.Unlock()

	// Close client channels
	close(client.Updates)
	close(client.Done)
}

// broadcastUpdate broadcasts an update to all clients
func (rm *RealtimeManager) broadcastUpdate(update RealtimeUpdate) {
	rm.clientsMux.RLock()
	defer rm.clientsMux.RUnlock()

	for _, client := range rm.clients {
		// Skip if update is for a specific user and client doesn't match
		if update.UserID != nil && client.UserID != *update.UserID {
			continue
		}

		// Skip if update is for a specific client and this isn't it
		if update.ClientID != "" && client.ID != update.ClientID {
			continue
		}

		select {
		case client.Updates <- update:
			client.LastSeen = time.Now()
		default:
			// Client channel is full, skip this update
		}
	}
}

// sendHeartbeat sends heartbeat to all clients
func (rm *RealtimeManager) sendHeartbeat() {
	heartbeatUpdate := RealtimeUpdate{
		Type: UpdateTypeHeartbeat,
		Data: map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "alive",
		},
		Timestamp: time.Now(),
	}

	rm.broadcastUpdate(heartbeatUpdate)
}

// cleanupStaleClients removes clients that haven't been seen recently
func (rm *RealtimeManager) cleanupStaleClients() {
	staleThreshold := time.Now().Add(-10 * time.Minute)

	rm.clientsMux.Lock()
	staleClients := make([]*Client, 0)

	for id, client := range rm.clients {
		if client.LastSeen.Before(staleThreshold) {
			staleClients = append(staleClients, client)
			delete(rm.clients, id)
		}
	}
	rm.clientsMux.Unlock()

	// Close stale client channels
	for _, client := range staleClients {
		close(client.Updates)
		close(client.Done)
	}
}

// cleanup closes all client connections
func (rm *RealtimeManager) cleanup() {
	rm.clientsMux.Lock()
	defer rm.clientsMux.Unlock()

	for _, client := range rm.clients {
		close(client.Updates)
		close(client.Done)
	}

	rm.clients = make(map[string]*Client)
}

// Helper functions for creating specific update types

// CreateExecutionUpdate creates an execution-related update
func CreateExecutionUpdate(updateType string, executionID uuid.UUID, workflowID uuid.UUID, data interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: updateType,
		Data: data,
		Metadata: map[string]interface{}{
			"execution_id": executionID,
			"workflow_id":  workflowID,
		},
	}
}

// CreateWorkflowUpdate creates a workflow-related update
func CreateWorkflowUpdate(updateType string, workflowID uuid.UUID, data interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: updateType,
		Data: data,
		Metadata: map[string]interface{}{
			"workflow_id": workflowID,
		},
	}
}

// CreateMetricsUpdate creates a metrics-related update
func CreateMetricsUpdate(metricsType string, data interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: UpdateTypeMetricsUpdated,
		Data: data,
		Metadata: map[string]interface{}{
			"metrics_type": metricsType,
		},
	}
}

// CreateAlertUpdate creates an alert-related update
func CreateAlertUpdate(updateType string, alertID uuid.UUID, data interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: updateType,
		Data: data,
		Metadata: map[string]interface{}{
			"alert_id": alertID,
		},
	}
}

// CreateSystemUpdate creates a system-related update
func CreateSystemUpdate(data interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: UpdateTypeSystemStatus,
		Data: data,
	}
}

// CreateErrorUpdate creates an error update
func CreateErrorUpdate(errorMsg string, details interface{}) RealtimeUpdate {
	return RealtimeUpdate{
		Type: UpdateTypeError,
		Data: map[string]interface{}{
			"error":   errorMsg,
			"details": details,
		},
	}
}

// MarshalUpdate marshals an update to JSON
func MarshalUpdate(update RealtimeUpdate) ([]byte, error) {
	return json.Marshal(update)
}

// UnmarshalUpdate unmarshals an update from JSON
func UnmarshalUpdate(data []byte) (*RealtimeUpdate, error) {
	var update RealtimeUpdate
	err := json.Unmarshal(data, &update)
	if err != nil {
		return nil, err
	}
	return &update, nil
}