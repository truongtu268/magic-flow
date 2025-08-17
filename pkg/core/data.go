package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// DefaultWorkflowData is the default implementation of WorkflowData
type DefaultWorkflowData struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewDefaultWorkflowData creates a new default workflow data instance
func NewDefaultWorkflowData() WorkflowData {
	return &DefaultWorkflowData{
		data: make(map[string]interface{}),
	}
}

// NewDefaultWorkflowDataWithMap creates a new default workflow data instance with initial data
func NewDefaultWorkflowDataWithMap(data map[string]interface{}) WorkflowData {
	wd := &DefaultWorkflowData{
		data: make(map[string]interface{}),
	}
	for k, v := range data {
		wd.data[k] = v
	}
	return wd
}

// Validate checks if the data structure is valid
func (d *DefaultWorkflowData) Validate() error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Basic validation - ensure data is not nil
	if d.data == nil {
		return fmt.Errorf("workflow data is nil")
	}
	return nil
}

// Convert converts the data to the target type
func (d *DefaultWorkflowData) Convert(target interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Convert to JSON first, then unmarshal to target
	jsonData, err := json.Marshal(d.data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	err = json.Unmarshal(jsonData, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	
	return nil
}

// GetAll returns all data as a map
func (d *DefaultWorkflowData) GetAll() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range d.data {
		result[k] = v
	}
	return result
}

// Get retrieves a value by key
func (d *DefaultWorkflowData) Get(key string) (interface{}, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	value, exists := d.data[key]
	return value, exists
}

// Set stores a value by key
func (d *DefaultWorkflowData) Set(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if d.data == nil {
		d.data = make(map[string]interface{})
	}
	d.data[key] = value
}

// Delete removes a value by key
func (d *DefaultWorkflowData) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.data, key)
}

// Has checks if a key exists
func (d *DefaultWorkflowData) Has(key string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.data[key]
	return exists
}

// Keys returns all keys
func (d *DefaultWorkflowData) Keys() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keys := make([]string, 0, len(d.data))
	for k := range d.data {
		keys = append(keys, k)
	}
	return keys
}

// Clear removes all data
func (d *DefaultWorkflowData) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = make(map[string]interface{})
}

// Size returns the number of items
func (d *DefaultWorkflowData) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.data)
}

// ToMap returns data as a map
func (d *DefaultWorkflowData) ToMap() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range d.data {
		result[k] = v
	}
	return result
}

// FromMap loads data from a map
func (d *DefaultWorkflowData) FromMap(data map[string]interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = make(map[string]interface{})
	for k, v := range data {
		d.data[k] = v
	}
}

// MustGet retrieves a value by key, panics if not found
func (d *DefaultWorkflowData) MustGet(key string) interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, exists := d.data[key]
	if !exists {
		panic(fmt.Sprintf("key '%s' not found in workflow data", key))
	}
	return value
}

// DefaultWorkflowMetadata is the default implementation of WorkflowMetadata
type DefaultWorkflowMetadata struct {
	metrics      map[string]interface{}
	tags         map[string]bool
	customFields map[string]interface{}
	mu           sync.RWMutex
}

// NewDefaultWorkflowMetadata creates a new default workflow metadata instance
func NewDefaultWorkflowMetadata() WorkflowMetadata {
	return &DefaultWorkflowMetadata{
		metrics:      make(map[string]interface{}),
		tags:         make(map[string]bool),
		customFields: make(map[string]interface{}),
	}
}

// NewDefaultWorkflowMetadataWithMap creates a new default workflow metadata instance with initial metrics
func NewDefaultWorkflowMetadataWithMap(metrics map[string]interface{}) WorkflowMetadata {
	md := &DefaultWorkflowMetadata{
		metrics:      make(map[string]interface{}),
		tags:         make(map[string]bool),
		customFields: make(map[string]interface{}),
	}
	for k, v := range metrics {
		md.metrics[k] = v
	}
	return md
}

// GetExecutionMetrics returns execution-related metadata
func (m *DefaultWorkflowMetadata) GetExecutionMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range m.metrics {
		result[k] = v
	}
	return result
}

// SetExecutionMetric sets an execution metric
func (m *DefaultWorkflowMetadata) SetExecutionMetric(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.metrics == nil {
		m.metrics = make(map[string]interface{})
	}
	m.metrics[key] = value
}

// GetExecutionMetric gets an execution metric
func (m *DefaultWorkflowMetadata) GetExecutionMetric(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.metrics[key]
	return value, exists
}

// AddTag adds a tag
func (m *DefaultWorkflowMetadata) AddTag(tag string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tags[tag] = true
}

// HasTag checks if a tag exists
func (m *DefaultWorkflowMetadata) HasTag(tag string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tags[tag]
}

// GetTags returns all tags
func (m *DefaultWorkflowMetadata) GetTags() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tags := make([]string, 0, len(m.tags))
	for tag := range m.tags {
		tags = append(tags, tag)
	}
	return tags
}

// RemoveTag removes a tag
func (m *DefaultWorkflowMetadata) RemoveTag(tag string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tags, tag)
}

// SetCustomField sets a custom field
func (m *DefaultWorkflowMetadata) SetCustomField(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customFields[key] = value
}

// GetCustomField gets a custom field
func (m *DefaultWorkflowMetadata) GetCustomField(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.customFields[key]
	return value, exists
}

// GetCustomFields returns all custom fields
func (m *DefaultWorkflowMetadata) GetCustomFields() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range m.customFields {
		result[k] = v
	}
	return result
}

// ToMap returns metadata as a map
func (m *DefaultWorkflowMetadata) ToMap() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := map[string]interface{}{
		"tags":              m.GetTags(),
		"execution_metrics": m.GetExecutionMetrics(),
		"custom_fields":     m.GetCustomFields(),
	}
	return result
}

// Validate checks if the metadata is valid
func (m *DefaultWorkflowMetadata) Validate() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Basic validation - ensure maps are not nil
	if m.metrics == nil || m.tags == nil || m.customFields == nil {
		return fmt.Errorf("workflow metadata is not properly initialized")
	}
	return nil
}

// GetString safely gets a string value from workflow data
func GetString(data WorkflowData, key string) (string, error) {
	value, exists := data.Get(key)
	if !exists {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("value for key '%s' is not a string", key)
}

// GetInt safely gets an int value from workflow data
func GetInt(data WorkflowData, key string) (int, error) {
	value, exists := data.Get(key)
	if !exists {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	
	// Handle different numeric types
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	default:
		return 0, fmt.Errorf("value for key '%s' is not an int", key)
	}
}

// GetBool safely gets a bool value from workflow data
func GetBool(data WorkflowData, key string) (bool, error) {
	value, exists := data.Get(key)
	if !exists {
		return false, fmt.Errorf("key '%s' not found", key)
	}
	
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("value for key '%s' is not a bool", key)
}

// GetFloat64 safely gets a float64 value from workflow data
func GetFloat64(data WorkflowData, key string) (float64, error) {
	value, exists := data.Get(key)
	if !exists {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	
	// Handle different numeric types
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("value for key '%s' is not a float64", key)
	}
}

// GetSlice safely gets a slice value from workflow data
func GetSlice(data WorkflowData, key string) ([]interface{}, error) {
	value, exists := data.Get(key)
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	
	// Use reflection to handle different slice types
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value for key '%s' is not a slice", key)
	}
	
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result, nil
}

// GetMap safely gets a map value from workflow data
func GetMap(data WorkflowData, key string) (map[string]interface{}, error) {
	value, exists := data.Get(key)
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	
	if m, ok := value.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, fmt.Errorf("value for key '%s' is not a map", key)
}