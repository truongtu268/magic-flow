package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultWorkflowData(t *testing.T) {
	t.Run("NewDefaultWorkflowData", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		require.NotNil(t, data)
		assert.Equal(t, 0, data.Size())
	})

	t.Run("Set and Get", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		
		// Test setting and getting a value
		data.Set("key1", "value1")
		value, exists := data.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", value)
		
		// Test getting non-existent key
		value, exists = data.Get("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("Has", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		
		assert.True(t, data.Has("key1"))
		assert.False(t, data.Has("nonexistent"))
	})

	t.Run("Delete", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		
		assert.True(t, data.Has("key1"))
		data.Delete("key1")
		assert.False(t, data.Has("key1"))
	})

	t.Run("Keys", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		data.Set("key2", "value2")
		
		keys := data.Keys()
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
	})

	t.Run("Clear", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		data.Set("key2", "value2")
		
		assert.Len(t, data.Keys(), 2)
		data.Clear()
		assert.Len(t, data.Keys(), 0)
	})

	t.Run("Size", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		assert.Equal(t, 0, data.Size())
		
		data.Set("key1", "value1")
		assert.Equal(t, 1, data.Size())
		
		data.Set("key2", "value2")
		assert.Equal(t, 2, data.Size())
	})

	t.Run("ToMap", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		data.Set("key2", 42)
		
		result := data.ToMap()
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		}
		assert.Equal(t, expected, result)
	})

	t.Run("FromMap", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		input := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		}
		
		data.FromMap(input)
		assert.Equal(t, "value1", data.MustGet("key1"))
		assert.Equal(t, 42, data.MustGet("key2"))
	})

	t.Run("MustGet", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("key1", "value1")
		
		// Test existing key
		value := data.MustGet("key1")
		assert.Equal(t, "value1", value)
		
		// Test non-existent key should panic
		assert.Panics(t, func() {
			data.MustGet("nonexistent")
		})
	})

	t.Run("Validate", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		
		// Empty data should be valid
		err := data.Validate()
		assert.NoError(t, err)
		
		// Data with values should be valid
		data.Set("key1", "value1")
		err = data.Validate()
		assert.NoError(t, err)
	})
}

func TestDefaultWorkflowMetadata(t *testing.T) {
	t.Run("NewDefaultWorkflowMetadata", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		require.NotNil(t, meta)
		assert.NoError(t, meta.Validate())
		assert.Empty(t, meta.GetTags())
		assert.Empty(t, meta.GetExecutionMetrics())
		assert.Empty(t, meta.GetCustomFields())
	})

	t.Run("Tags", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		
		// Test adding tags
		meta.AddTag("tag1")
		meta.AddTag("tag2")
		assert.True(t, meta.HasTag("tag1"))
		assert.True(t, meta.HasTag("tag2"))
		assert.False(t, meta.HasTag("tag3"))
		
		// Test getting tags
		tags := meta.GetTags()
		assert.Len(t, tags, 2)
		assert.Contains(t, tags, "tag1")
		assert.Contains(t, tags, "tag2")
		
		// Test removing tags
		meta.RemoveTag("tag1")
		assert.False(t, meta.HasTag("tag1"))
		assert.True(t, meta.HasTag("tag2"))
	})

	t.Run("ExecutionMetrics", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		
		// Test setting and getting metrics
		meta.SetExecutionMetric("duration", time.Second)
		meta.SetExecutionMetric("steps_count", 5)
		
		value, exists := meta.GetExecutionMetric("duration")
		assert.True(t, exists)
		assert.Equal(t, time.Second, value)
		
		value, exists = meta.GetExecutionMetric("steps_count")
		assert.True(t, exists)
		assert.Equal(t, 5, value)
		
		// Test non-existent metric
		value, exists = meta.GetExecutionMetric("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
		
		// Test getting all metrics
		metrics := meta.GetExecutionMetrics()
		assert.Len(t, metrics, 2)
		assert.Equal(t, time.Second, metrics["duration"])
		assert.Equal(t, 5, metrics["steps_count"])
	})

	t.Run("CustomFields", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		
		// Test setting and getting custom fields
		meta.SetCustomField("priority", "high")
		meta.SetCustomField("department", "engineering")
		
		value, exists := meta.GetCustomField("priority")
		assert.True(t, exists)
		assert.Equal(t, "high", value)
		
		value, exists = meta.GetCustomField("department")
		assert.True(t, exists)
		assert.Equal(t, "engineering", value)
		
		// Test non-existent field
		value, exists = meta.GetCustomField("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
		
		// Test getting all custom fields
		fields := meta.GetCustomFields()
		assert.Len(t, fields, 2)
		assert.Equal(t, "high", fields["priority"])
		assert.Equal(t, "engineering", fields["department"])
	})

	t.Run("ToMap", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		meta.AddTag("tag1")
		meta.SetExecutionMetric("duration", time.Second)
		meta.SetCustomField("priority", "high")
		
		result := meta.ToMap()
		assert.Contains(t, result, "tags")
		assert.Contains(t, result, "execution_metrics")
		assert.Contains(t, result, "custom_fields")
	})

	t.Run("Validate", func(t *testing.T) {
		meta := NewDefaultWorkflowMetadata()
		
		// Empty metadata should be valid
		err := meta.Validate()
		assert.NoError(t, err)
		
		// Metadata with values should be valid
		meta.AddTag("tag1")
		meta.SetExecutionMetric("duration", time.Second)
		err = meta.Validate()
		assert.NoError(t, err)
	})
}

func TestWorkflowDataHelpers(t *testing.T) {
	t.Run("GetString", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("string_key", "hello")
		data.Set("int_key", 42)
		
		// Test valid string
		value, err := GetString(data, "string_key")
		assert.NoError(t, err)
		assert.Equal(t, "hello", value)
		
		// Test invalid type
		value, err = GetString(data, "int_key")
		assert.Error(t, err)
		assert.Empty(t, value)
		
		// Test non-existent key
		value, err = GetString(data, "nonexistent")
		assert.Error(t, err)
		assert.Empty(t, value)
	})

	t.Run("GetInt", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("int_key", 42)
		data.Set("string_key", "hello")
		
		// Test valid int
		value, err := GetInt(data, "int_key")
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
		
		// Test invalid type
		value, err = GetInt(data, "string_key")
		assert.Error(t, err)
		assert.Equal(t, 0, value)
		
		// Test non-existent key
		value, err = GetInt(data, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("GetBool", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("bool_key", true)
		data.Set("string_key", "hello")
		
		// Test valid bool
		value, err := GetBool(data, "bool_key")
		assert.NoError(t, err)
		assert.True(t, value)
		
		// Test invalid type
		value, err = GetBool(data, "string_key")
		assert.Error(t, err)
		assert.False(t, value)
		
		// Test non-existent key
		value, err = GetBool(data, "nonexistent")
		assert.Error(t, err)
		assert.False(t, value)
	})

	t.Run("GetFloat64", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		data.Set("float_key", 3.14)
		data.Set("string_key", "hello")
		
		// Test valid float64
		value, err := GetFloat64(data, "float_key")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, value)
		
		// Test invalid type
		value, err = GetFloat64(data, "string_key")
		assert.Error(t, err)
		assert.Equal(t, 0.0, value)
		
		// Test non-existent key
		value, err = GetFloat64(data, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, 0.0, value)
	})

	t.Run("GetSlice", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		slice := []interface{}{"a", "b", "c"}
		data.Set("slice_key", slice)
		data.Set("string_key", "hello")
		
		// Test valid slice
		value, err := GetSlice(data, "slice_key")
		assert.NoError(t, err)
		assert.Equal(t, slice, value)
		
		// Test invalid type
		value, err = GetSlice(data, "string_key")
		assert.Error(t, err)
		assert.Nil(t, value)
		
		// Test non-existent key
		value, err = GetSlice(data, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, value)
	})

	t.Run("GetMap", func(t *testing.T) {
		data := NewDefaultWorkflowData()
		mapValue := map[string]interface{}{"key": "value"}
		data.Set("map_key", mapValue)
		data.Set("string_key", "hello")
		
		// Test valid map
		value, err := GetMap(data, "map_key")
		assert.NoError(t, err)
		assert.Equal(t, mapValue, value)
		
		// Test invalid type
		value, err = GetMap(data, "string_key")
		assert.Error(t, err)
		assert.Nil(t, value)
		
		// Test non-existent key
		value, err = GetMap(data, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, value)
	})
}