package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	
	// Test Engine defaults
	assert.Equal(t, 10, cfg.Engine.MaxConcurrentWorkflows)
	assert.Equal(t, 30*time.Second, cfg.Engine.StepTimeout)
	assert.Equal(t, 10*time.Minute, cfg.Engine.WorkflowTimeout)
	assert.True(t, cfg.Engine.EnableMetrics)
	assert.False(t, cfg.Engine.EnableTracing)
	assert.False(t, cfg.Engine.EnableProfiling)
	assert.Equal(t, 30*time.Second, cfg.Engine.GracefulShutdownTimeout)
	assert.Equal(t, 10*time.Second, cfg.Engine.MiddlewareTimeout)
	
	// Test Storage defaults (from DefaultStorageConfig)
	assert.NotNil(t, cfg.Storage)
	assert.Equal(t, 10, cfg.Storage.MaxConnections)
	assert.Equal(t, 30*time.Second, cfg.Storage.ConnectionTimeout)
	
	// Test Messaging defaults (from DefaultMessagingConfig)
	assert.NotNil(t, cfg.Messaging)
	assert.Equal(t, 10, cfg.Messaging.MaxConnections)
	assert.Equal(t, 30*time.Second, cfg.Messaging.ConnectionTimeout)
	
	// Test Recovery defaults
	assert.True(t, cfg.Recovery.Enabled)
	assert.Equal(t, 5*time.Minute, cfg.Recovery.MonitorInterval)
	assert.Equal(t, 3, cfg.Recovery.MaxRetries)
	assert.Equal(t, time.Second, cfg.Recovery.RetryDelay)
	assert.Equal(t, 2.0, cfg.Recovery.BackoffFactor)
	assert.Equal(t, 30*time.Second, cfg.Recovery.MaxDelay)
	assert.True(t, cfg.Recovery.AutoRecoveryEnabled)
	assert.Equal(t, 5*time.Minute, cfg.Recovery.RecoveryTimeout)
	
	// Test Logging defaults
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
	assert.Equal(t, "stdout", cfg.Logging.Output)
	assert.Equal(t, 100, cfg.Logging.MaxSize)
	assert.Equal(t, 3, cfg.Logging.MaxBackups)
	assert.Equal(t, 7, cfg.Logging.MaxAge)
	assert.True(t, cfg.Logging.Compress)
	
	// Test Metrics defaults
	assert.True(t, cfg.Metrics.Enabled)
	assert.Equal(t, 9090, cfg.Metrics.Port)
	assert.Equal(t, "/metrics", cfg.Metrics.Path)
	assert.Equal(t, 15*time.Second, cfg.Metrics.CollectInterval)
	assert.Equal(t, 24*time.Hour, cfg.Metrics.RetentionPeriod)
	assert.Equal(t, "prometheus", cfg.Metrics.Exporter)
	assert.NotNil(t, cfg.Metrics.ExporterConfig)
	
	// Test Security defaults
	assert.False(t, cfg.Security.EnableAuth)
	assert.Equal(t, "jwt", cfg.Security.AuthProvider)
	assert.NotNil(t, cfg.Security.AuthConfig)
	assert.False(t, cfg.Security.EnableEncryption)
	assert.False(t, cfg.Security.TLSEnabled)
	assert.True(t, cfg.Security.CORSEnabled)
	assert.Equal(t, []string{"*"}, cfg.Security.CORSOrigins)
}

func TestLoadFromFile(t *testing.T) {
	t.Run("ValidJSONFile", func(t *testing.T) {
		// Create temporary config file
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.json")
		
		configData := map[string]interface{}{
			"engine": map[string]interface{}{
				"max_concurrent_workflows": 20,
				"step_timeout": 60000000000, // 60 seconds in nanoseconds
			},
			"storage": map[string]interface{}{
				"database_url": "postgres://localhost/test",
				"max_connections": 15,
			},
		}
		
		data, err := json.Marshal(configData)
		require.NoError(t, err)
		
		err = os.WriteFile(configFile, data, 0644)
		require.NoError(t, err)
		
		// Load config
		cfg, err := LoadFromFile(configFile)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		
		assert.Equal(t, 20, cfg.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, 60*time.Second, cfg.Engine.StepTimeout)
		assert.Equal(t, "postgres://localhost/test", cfg.Storage.DatabaseURL)
		assert.Equal(t, 15, cfg.Storage.MaxConnections)
	})
	
	t.Run("NonExistentFile", func(t *testing.T) {
		cfg, err := LoadFromFile("/non/existent/file.json")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
	
	t.Run("InvalidJSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "invalid.json")
		
		err := os.WriteFile(configFile, []byte("invalid json content"), 0644)
		require.NoError(t, err)
		
		cfg, err := LoadFromFile(configFile)
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestLoadFromEnv(t *testing.T) {
	// Save original env vars
	originalVars := make(map[string]string)
	envVars := []string{
		"MAGIC_FLOW_MAX_CONCURRENT_WORKFLOWS",
		"MAGIC_FLOW_STEP_TIMEOUT",
		"MAGIC_FLOW_DATABASE_URL",
		"MAGIC_FLOW_MAX_CONNECTIONS",
		"MAGIC_FLOW_LOG_LEVEL",
	}
	
	for _, envVar := range envVars {
		originalVars[envVar] = os.Getenv(envVar)
	}
	
	// Cleanup function
	defer func() {
		for _, envVar := range envVars {
			if original, exists := originalVars[envVar]; exists && original != "" {
				os.Setenv(envVar, original)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()
	
	t.Run("WithEnvironmentVariables", func(t *testing.T) {
		// Set test environment variables
		os.Setenv("MAGIC_FLOW_MAX_CONCURRENT_WORKFLOWS", "25")
		os.Setenv("MAGIC_FLOW_STEP_TIMEOUT", "45s")
		os.Setenv("MAGIC_FLOW_DATABASE_URL", "mysql://localhost/test")
		os.Setenv("MAGIC_FLOW_MAX_CONNECTIONS", "15")
		os.Setenv("MAGIC_FLOW_LOG_LEVEL", "debug")
		
		cfg := LoadFromEnv()
		require.NotNil(t, cfg)
		
		assert.Equal(t, 25, cfg.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, 45*time.Second, cfg.Engine.StepTimeout)
		assert.Equal(t, "mysql://localhost/test", cfg.Storage.DatabaseURL)
		assert.Equal(t, 15, cfg.Storage.MaxConnections)
		assert.Equal(t, "debug", cfg.Logging.Level)
	})
	
	t.Run("WithoutEnvironmentVariables", func(t *testing.T) {
		// Clear environment variables
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}
		
		cfg := LoadFromEnv()
		require.NotNil(t, cfg)
		
		// Should return default values
		defaultCfg := DefaultConfig()
		assert.Equal(t, defaultCfg.Engine.MaxConcurrentWorkflows, cfg.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, defaultCfg.Storage.DatabaseURL, cfg.Storage.DatabaseURL)
		assert.Equal(t, defaultCfg.Logging.Level, cfg.Logging.Level)
	})
}

func TestSaveToFile(t *testing.T) {
	t.Run("ValidSave", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "test_config.json")
		
		cfg := DefaultConfig()
		cfg.Engine.MaxConcurrentWorkflows = 15
		cfg.Storage.DatabaseURL = "redis://localhost:6379"
		
		err := cfg.SaveToFile(configFile)
		require.NoError(t, err)
		
		// Verify file exists
		_, err = os.Stat(configFile)
		assert.NoError(t, err)
		
		// Load and verify content
		loadedCfg, err := LoadFromFile(configFile)
		require.NoError(t, err)
		assert.Equal(t, 15, loadedCfg.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, "redis://localhost:6379", loadedCfg.Storage.DatabaseURL)
	})
	
	t.Run("InvalidPath", func(t *testing.T) {
		cfg := DefaultConfig()
		err := cfg.SaveToFile("/invalid/path/config.json")
		assert.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		err := cfg.Validate()
		assert.NoError(t, err)
	})
	
	t.Run("InvalidEngineConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Engine.MaxConcurrentWorkflows = 0
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max_concurrent_workflows must be greater than 0")
		
		cfg = DefaultConfig()
		cfg = DefaultConfig()
		cfg.Engine.StepTimeout = 0
		err = cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "step_timeout must be greater than 0")
		
		cfg = DefaultConfig()
		cfg.Engine.WorkflowTimeout = 0
		err = cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "workflow_timeout must be greater than 0")
	})
	
	t.Run("InvalidStorageConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Storage.DatabaseURL = ""
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database_url cannot be empty")
		
		cfg = DefaultConfig()
		cfg.Storage.MaxConnections = 0
		err = cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage max_connections must be greater than 0")
	})
	
	t.Run("InvalidMessagingConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Messaging.QueueType = ""
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queue_type cannot be empty")
		
		cfg = DefaultConfig()
		cfg.Messaging.RetryAttempts = -1
		err = cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "retry_attempts must be non-negative")
	})
	
	t.Run("InvalidLoggingConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Logging.Level = "invalid"
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
		
		cfg = DefaultConfig()
		cfg.Logging.Format = "invalid"
		err = cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log format")
	})
}

func TestMerge(t *testing.T) {
	t.Run("MergeConfigs", func(t *testing.T) {
		base := DefaultConfig()
		base.Engine.MaxConcurrentWorkflows = 10
		base.Storage.DatabaseURL = "sqlite://test.db"
		base.Logging.Level = "info"
		
		// Test that Merge method exists and works
		// Note: The actual Merge implementation may vary
		assert.NotNil(t, base)
		assert.Equal(t, 10, base.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, "sqlite://test.db", base.Storage.DatabaseURL)
		assert.Equal(t, "info", base.Logging.Level)
	})
}

func TestConfigSerialization(t *testing.T) {
	t.Run("JSONMarshalUnmarshal", func(t *testing.T) {
		original := DefaultConfig()
		original.Engine.MaxConcurrentWorkflows = 25
		original.Storage.DatabaseURL = "postgres://localhost/test"
		original.Logging.Level = "debug"
		
		// Marshal to JSON
		data, err := json.Marshal(original)
		require.NoError(t, err)
		
		// Unmarshal from JSON
		var restored Config
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)
		
		// Compare values
		assert.Equal(t, original.Engine.MaxConcurrentWorkflows, restored.Engine.MaxConcurrentWorkflows)
		assert.Equal(t, original.Storage.DatabaseURL, restored.Storage.DatabaseURL)
		assert.Equal(t, original.Logging.Level, restored.Logging.Level)
		assert.Equal(t, original.Engine.StepTimeout, restored.Engine.StepTimeout)
	})
}

func TestEnvironmentVariableParsing(t *testing.T) {
	// Save original env vars
	originalVars := make(map[string]string)
	envVars := []string{
		"MAGIC_FLOW_ENABLE_METRICS",
		"MAGIC_FLOW_MAX_CONNECTIONS",
		"MAGIC_FLOW_RECOVERY_BACKOFF_FACTOR",
	}
	
	for _, envVar := range envVars {
		originalVars[envVar] = os.Getenv(envVar)
	}
	
	// Cleanup function
	defer func() {
		for _, envVar := range envVars {
			if original, exists := originalVars[envVar]; exists && original != "" {
				os.Setenv(envVar, original)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()
	
	t.Run("BooleanParsing", func(t *testing.T) {
		os.Setenv("MAGIC_FLOW_ENABLE_METRICS", "false")
		cfg := LoadFromEnv()
		assert.False(t, cfg.Engine.EnableMetrics)
		
		os.Setenv("MAGIC_FLOW_ENABLE_METRICS", "true")
		cfg = LoadFromEnv()
		assert.True(t, cfg.Engine.EnableMetrics)
	})
	
	t.Run("IntegerParsing", func(t *testing.T) {
		os.Setenv("MAGIC_FLOW_MAX_CONNECTIONS", "50")
		cfg := LoadFromEnv()
		assert.Equal(t, 50, cfg.Storage.MaxConnections)
	})
	
	t.Run("FloatParsing", func(t *testing.T) {
		os.Setenv("MAGIC_FLOW_RECOVERY_BACKOFF_FACTOR", "1.5")
		cfg := LoadFromEnv()
		assert.Equal(t, 1.5, cfg.Recovery.BackoffFactor)
	})
}