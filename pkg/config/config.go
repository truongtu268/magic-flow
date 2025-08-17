package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/truongtu268/magic-flow/pkg/messaging"
	"github.com/truongtu268/magic-flow/pkg/storage"
)

// Config represents the main configuration for the magic-flow library
type Config struct {
	Engine    *EngineConfig    `json:"engine"`
	Storage   *storage.StorageConfig `json:"storage"`
	Messaging *messaging.MessagingConfig `json:"messaging"`
	Recovery  *RecoveryConfig  `json:"recovery"`
	Logging   *LoggingConfig   `json:"logging"`
	Metrics   *MetricsConfig   `json:"metrics"`
	Security  *SecurityConfig  `json:"security"`
}

// EngineConfig defines configuration for the workflow engine
type EngineConfig struct {
	MaxConcurrentWorkflows int           `json:"max_concurrent_workflows"`
	StepTimeout            time.Duration `json:"step_timeout"`
	WorkflowTimeout        time.Duration `json:"workflow_timeout"`
	EnableMetrics          bool          `json:"enable_metrics"`
	EnableTracing          bool          `json:"enable_tracing"`
	EnableProfiling        bool          `json:"enable_profiling"`
	GracefulShutdownTimeout time.Duration `json:"graceful_shutdown_timeout"`
	MiddlewareTimeout      time.Duration `json:"middleware_timeout"`
}

// RecoveryConfig defines configuration for workflow recovery
type RecoveryConfig struct {
	Enabled              bool          `json:"enabled"`
	MonitorInterval      time.Duration `json:"monitor_interval"`
	MaxRetries           int           `json:"max_retries"`
	RetryDelay           time.Duration `json:"retry_delay"`
	BackoffFactor        float64       `json:"backoff_factor"`
	MaxDelay             time.Duration `json:"max_delay"`
	AutoRecoveryEnabled  bool          `json:"auto_recovery_enabled"`
	RecoveryTimeout      time.Duration `json:"recovery_timeout"`
}

// LoggingConfig defines configuration for logging
type LoggingConfig struct {
	Level      string `json:"level"`      // debug, info, warn, error
	Format     string `json:"format"`     // json, text
	Output     string `json:"output"`     // stdout, stderr, file
	FilePath   string `json:"file_path"`  // path to log file if output is file
	MaxSize    int    `json:"max_size"`   // max size in MB
	MaxBackups int    `json:"max_backups"` // max number of backup files
	MaxAge     int    `json:"max_age"`    // max age in days
	Compress   bool   `json:"compress"`   // compress backup files
}

// MetricsConfig defines configuration for metrics collection
type MetricsConfig struct {
	Enabled        bool          `json:"enabled"`
	Port           int           `json:"port"`
	Path           string        `json:"path"`
	CollectInterval time.Duration `json:"collect_interval"`
	RetentionPeriod time.Duration `json:"retention_period"`
	Exporter       string        `json:"exporter"` // prometheus, statsd, etc.
	ExporterConfig map[string]interface{} `json:"exporter_config"`
}

// SecurityConfig defines security-related configuration
type SecurityConfig struct {
	EnableAuth       bool              `json:"enable_auth"`
	AuthProvider     string            `json:"auth_provider"` // jwt, oauth2, etc.
	AuthConfig       map[string]interface{} `json:"auth_config"`
	EnableEncryption bool              `json:"enable_encryption"`
	EncryptionKey    string            `json:"encryption_key"`
	TLSEnabled       bool              `json:"tls_enabled"`
	TLSCertFile      string            `json:"tls_cert_file"`
	TLSKeyFile       string            `json:"tls_key_file"`
	CORSEnabled      bool              `json:"cors_enabled"`
	CORSOrigins      []string          `json:"cors_origins"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Engine: &EngineConfig{
			MaxConcurrentWorkflows:  10,
			StepTimeout:             30 * time.Second,
			WorkflowTimeout:         10 * time.Minute,
			EnableMetrics:           true,
			EnableTracing:           false,
			EnableProfiling:         false,
			GracefulShutdownTimeout: 30 * time.Second,
			MiddlewareTimeout:       10 * time.Second,
		},
		Storage:   storage.DefaultStorageConfig(),
		Messaging: messaging.DefaultMessagingConfig(),
		Recovery: &RecoveryConfig{
			Enabled:             true,
			MonitorInterval:     5 * time.Minute,
			MaxRetries:          3,
			RetryDelay:          1 * time.Second,
			BackoffFactor:       2.0,
			MaxDelay:            30 * time.Second,
			AutoRecoveryEnabled: true,
			RecoveryTimeout:     5 * time.Minute,
		},
		Logging: &LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		},
		Metrics: &MetricsConfig{
			Enabled:         true,
			Port:            9090,
			Path:            "/metrics",
			CollectInterval: 15 * time.Second,
			RetentionPeriod: 24 * time.Hour,
			Exporter:        "prometheus",
			ExporterConfig:  make(map[string]interface{}),
		},
		Security: &SecurityConfig{
			EnableAuth:       false,
			AuthProvider:     "jwt",
			AuthConfig:       make(map[string]interface{}),
			EnableEncryption: false,
			TLSEnabled:       false,
			CORSEnabled:      true,
			CORSOrigins:      []string{"*"},
		},
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	config := DefaultConfig()
	
	// Engine configuration
	if val := getEnvInt("MAGIC_FLOW_MAX_CONCURRENT_WORKFLOWS", config.Engine.MaxConcurrentWorkflows); val != config.Engine.MaxConcurrentWorkflows {
		config.Engine.MaxConcurrentWorkflows = val
	}
	if val := getEnvDuration("MAGIC_FLOW_STEP_TIMEOUT", config.Engine.StepTimeout); val != config.Engine.StepTimeout {
		config.Engine.StepTimeout = val
	}
	if val := getEnvDuration("MAGIC_FLOW_WORKFLOW_TIMEOUT", config.Engine.WorkflowTimeout); val != config.Engine.WorkflowTimeout {
		config.Engine.WorkflowTimeout = val
	}
	config.Engine.EnableMetrics = getEnvBool("MAGIC_FLOW_ENABLE_METRICS", config.Engine.EnableMetrics)
	config.Engine.EnableTracing = getEnvBool("MAGIC_FLOW_ENABLE_TRACING", config.Engine.EnableTracing)
	config.Engine.EnableProfiling = getEnvBool("MAGIC_FLOW_ENABLE_PROFILING", config.Engine.EnableProfiling)
	
	// Storage configuration
	if val := getEnvString("MAGIC_FLOW_DATABASE_URL", config.Storage.DatabaseURL); val != config.Storage.DatabaseURL {
		config.Storage.DatabaseURL = val
	}
	if val := getEnvInt("MAGIC_FLOW_MAX_CONNECTIONS", config.Storage.MaxConnections); val != config.Storage.MaxConnections {
		config.Storage.MaxConnections = val
	}
	
	// Messaging configuration
	if val := getEnvString("MAGIC_FLOW_QUEUE_URL", config.Messaging.QueueURL); val != config.Messaging.QueueURL {
		config.Messaging.QueueURL = val
	}
	if val := getEnvString("MAGIC_FLOW_QUEUE_TYPE", config.Messaging.QueueType); val != config.Messaging.QueueType {
		config.Messaging.QueueType = val
	}
	if val := getEnvString("MAGIC_FLOW_PUBSUB_URL", config.Messaging.PubSubURL); val != config.Messaging.PubSubURL {
		config.Messaging.PubSubURL = val
	}
	
	// Recovery configuration
	config.Recovery.Enabled = getEnvBool("MAGIC_FLOW_RECOVERY_ENABLED", config.Recovery.Enabled)
	if val := getEnvDuration("MAGIC_FLOW_RECOVERY_MONITOR_INTERVAL", config.Recovery.MonitorInterval); val != config.Recovery.MonitorInterval {
		config.Recovery.MonitorInterval = val
	}
	if val := getEnvInt("MAGIC_FLOW_RECOVERY_MAX_RETRIES", config.Recovery.MaxRetries); val != config.Recovery.MaxRetries {
		config.Recovery.MaxRetries = val
	}
	if val := getEnvFloat("MAGIC_FLOW_RECOVERY_BACKOFF_FACTOR", config.Recovery.BackoffFactor); val != config.Recovery.BackoffFactor {
		config.Recovery.BackoffFactor = val
	}
	config.Recovery.AutoRecoveryEnabled = getEnvBool("MAGIC_FLOW_AUTO_RECOVERY_ENABLED", config.Recovery.AutoRecoveryEnabled)
	
	// Logging configuration
	if val := getEnvString("MAGIC_FLOW_LOG_LEVEL", config.Logging.Level); val != config.Logging.Level {
		config.Logging.Level = val
	}
	if val := getEnvString("MAGIC_FLOW_LOG_FORMAT", config.Logging.Format); val != config.Logging.Format {
		config.Logging.Format = val
	}
	if val := getEnvString("MAGIC_FLOW_LOG_OUTPUT", config.Logging.Output); val != config.Logging.Output {
		config.Logging.Output = val
	}
	if val := getEnvString("MAGIC_FLOW_LOG_FILE_PATH", config.Logging.FilePath); val != config.Logging.FilePath {
		config.Logging.FilePath = val
	}
	
	// Metrics configuration
	config.Metrics.Enabled = getEnvBool("MAGIC_FLOW_METRICS_ENABLED", config.Metrics.Enabled)
	if val := getEnvInt("MAGIC_FLOW_METRICS_PORT", config.Metrics.Port); val != config.Metrics.Port {
		config.Metrics.Port = val
	}
	if val := getEnvString("MAGIC_FLOW_METRICS_PATH", config.Metrics.Path); val != config.Metrics.Path {
		config.Metrics.Path = val
	}
	
	// Security configuration
	config.Security.EnableAuth = getEnvBool("MAGIC_FLOW_ENABLE_AUTH", config.Security.EnableAuth)
	if val := getEnvString("MAGIC_FLOW_AUTH_PROVIDER", config.Security.AuthProvider); val != config.Security.AuthProvider {
		config.Security.AuthProvider = val
	}
	config.Security.EnableEncryption = getEnvBool("MAGIC_FLOW_ENABLE_ENCRYPTION", config.Security.EnableEncryption)
	if val := getEnvString("MAGIC_FLOW_ENCRYPTION_KEY", config.Security.EncryptionKey); val != config.Security.EncryptionKey {
		config.Security.EncryptionKey = val
	}
	config.Security.TLSEnabled = getEnvBool("MAGIC_FLOW_TLS_ENABLED", config.Security.TLSEnabled)
	if val := getEnvString("MAGIC_FLOW_TLS_CERT_FILE", config.Security.TLSCertFile); val != config.Security.TLSCertFile {
		config.Security.TLSCertFile = val
	}
	if val := getEnvString("MAGIC_FLOW_TLS_KEY_FILE", config.Security.TLSKeyFile); val != config.Security.TLSKeyFile {
		config.Security.TLSKeyFile = val
	}
	config.Security.CORSEnabled = getEnvBool("MAGIC_FLOW_CORS_ENABLED", config.Security.CORSEnabled)
	if val := getEnvString("MAGIC_FLOW_CORS_ORIGINS", ""); val != "" {
		config.Security.CORSOrigins = strings.Split(val, ",")
	}
	
	return config
}

// SaveToFile saves configuration to a JSON file
func (c *Config) SaveToFile(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Engine == nil {
		return fmt.Errorf("engine configuration is required")
	}
	
	if c.Engine.MaxConcurrentWorkflows <= 0 {
		return fmt.Errorf("max_concurrent_workflows must be greater than 0")
	}
	
	if c.Engine.StepTimeout <= 0 {
		return fmt.Errorf("step_timeout must be greater than 0")
	}
	
	if c.Engine.WorkflowTimeout <= 0 {
		return fmt.Errorf("workflow_timeout must be greater than 0")
	}
	
	if c.Storage == nil {
		return fmt.Errorf("storage configuration is required")
	}

	if c.Storage.DatabaseURL == "" {
		return fmt.Errorf("database_url cannot be empty")
	}

	if c.Storage.MaxConnections <= 0 {
		return fmt.Errorf("storage max_connections must be greater than 0")
	}

	if c.Messaging == nil {
		return fmt.Errorf("messaging configuration is required")
	}

	if c.Messaging.QueueType == "" {
		return fmt.Errorf("queue_type cannot be empty")
	}

	if c.Messaging.RetryAttempts < 0 {
		return fmt.Errorf("messaging retry_attempts must be non-negative")
	}
	
	if c.Recovery == nil {
		return fmt.Errorf("recovery configuration is required")
	}
	
	if c.Recovery.MaxRetries < 0 {
		return fmt.Errorf("recovery max_retries must be non-negative")
	}
	
	if c.Recovery.BackoffFactor <= 0 {
		return fmt.Errorf("recovery backoff_factor must be greater than 0")
	}
	
	if c.Logging == nil {
		return fmt.Errorf("logging configuration is required")
	}
	
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.Logging.Level) {
		return fmt.Errorf("invalid log level: %s, must be one of %v", c.Logging.Level, validLogLevels)
	}
	
	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, c.Logging.Format) {
		return fmt.Errorf("invalid log format: %s, must be one of %v", c.Logging.Format, validLogFormats)
	}
	
	validLogOutputs := []string{"stdout", "stderr", "file"}
	if !contains(validLogOutputs, c.Logging.Output) {
		return fmt.Errorf("invalid log output: %s, must be one of %v", c.Logging.Output, validLogOutputs)
	}
	
	if c.Logging.Output == "file" && c.Logging.FilePath == "" {
		return fmt.Errorf("file_path is required when log output is 'file'")
	}
	
	if c.Metrics == nil {
		return fmt.Errorf("metrics configuration is required")
	}
	
	if c.Metrics.Enabled && (c.Metrics.Port <= 0 || c.Metrics.Port > 65535) {
		return fmt.Errorf("metrics port must be between 1 and 65535")
	}
	
	if c.Security == nil {
		return fmt.Errorf("security configuration is required")
	}
	
	if c.Security.TLSEnabled {
		if c.Security.TLSCertFile == "" {
			return fmt.Errorf("tls_cert_file is required when TLS is enabled")
		}
		if c.Security.TLSKeyFile == "" {
			return fmt.Errorf("tls_key_file is required when TLS is enabled")
		}
	}
	
	return nil
}

// Merge merges another configuration into this one
func (c *Config) Merge(other *Config) {
	if other.Engine != nil {
		if c.Engine == nil {
			c.Engine = &EngineConfig{}
		}
		mergeEngineConfig(c.Engine, other.Engine)
	}
	
	if other.Storage != nil {
		if c.Storage == nil {
			c.Storage = &storage.StorageConfig{}
		}
		mergeStorageConfig(c.Storage, other.Storage)
	}
	
	if other.Messaging != nil {
		if c.Messaging == nil {
			c.Messaging = &messaging.MessagingConfig{}
		}
		mergeMessagingConfig(c.Messaging, other.Messaging)
	}
	
	if other.Recovery != nil {
		if c.Recovery == nil {
			c.Recovery = &RecoveryConfig{}
		}
		mergeRecoveryConfig(c.Recovery, other.Recovery)
	}
	
	if other.Logging != nil {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		mergeLoggingConfig(c.Logging, other.Logging)
	}
	
	if other.Metrics != nil {
		if c.Metrics == nil {
			c.Metrics = &MetricsConfig{}
		}
		mergeMetricsConfig(c.Metrics, other.Metrics)
	}
	
	if other.Security != nil {
		if c.Security == nil {
			c.Security = &SecurityConfig{}
		}
		mergeSecurityConfig(c.Security, other.Security)
	}
}

// Helper functions

func getEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func mergeEngineConfig(dst, src *EngineConfig) {
	if src.MaxConcurrentWorkflows != 0 {
		dst.MaxConcurrentWorkflows = src.MaxConcurrentWorkflows
	}
	if src.StepTimeout != 0 {
		dst.StepTimeout = src.StepTimeout
	}
	if src.WorkflowTimeout != 0 {
		dst.WorkflowTimeout = src.WorkflowTimeout
	}
	dst.EnableMetrics = src.EnableMetrics
	dst.EnableTracing = src.EnableTracing
	dst.EnableProfiling = src.EnableProfiling
	if src.GracefulShutdownTimeout != 0 {
		dst.GracefulShutdownTimeout = src.GracefulShutdownTimeout
	}
	if src.MiddlewareTimeout != 0 {
		dst.MiddlewareTimeout = src.MiddlewareTimeout
	}
}

func mergeStorageConfig(dst, src *storage.StorageConfig) {
	if src.DatabaseURL != "" {
		dst.DatabaseURL = src.DatabaseURL
	}
	if src.MaxConnections != 0 {
		dst.MaxConnections = src.MaxConnections
	}
	if src.ConnectionTimeout != 0 {
		dst.ConnectionTimeout = src.ConnectionTimeout
	}
	if src.QueryTimeout != 0 {
		dst.QueryTimeout = src.QueryTimeout
	}
	if src.RetryAttempts != 0 {
		dst.RetryAttempts = src.RetryAttempts
	}
	if src.RetryDelay != 0 {
		dst.RetryDelay = src.RetryDelay
	}
	dst.EnableMetrics = src.EnableMetrics
	if src.TablePrefix != "" {
		dst.TablePrefix = src.TablePrefix
	}
}

func mergeMessagingConfig(dst, src *messaging.MessagingConfig) {
	if src.QueueURL != "" {
		dst.QueueURL = src.QueueURL
	}
	if src.QueueType != "" {
		dst.QueueType = src.QueueType
	}
	if src.MaxConnections != 0 {
		dst.MaxConnections = src.MaxConnections
	}
	if src.ConnectionTimeout != 0 {
		dst.ConnectionTimeout = src.ConnectionTimeout
	}
	if src.PubSubURL != "" {
		dst.PubSubURL = src.PubSubURL
	}
	if src.PubSubType != "" {
		dst.PubSubType = src.PubSubType
	}
	dst.EnableMetrics = src.EnableMetrics
	dst.EnableTracing = src.EnableTracing
	if src.LogLevel != "" {
		dst.LogLevel = src.LogLevel
	}
}

func mergeRecoveryConfig(dst, src *RecoveryConfig) {
	dst.Enabled = src.Enabled
	if src.MonitorInterval != 0 {
		dst.MonitorInterval = src.MonitorInterval
	}
	if src.MaxRetries != 0 {
		dst.MaxRetries = src.MaxRetries
	}
	if src.RetryDelay != 0 {
		dst.RetryDelay = src.RetryDelay
	}
	if src.BackoffFactor != 0 {
		dst.BackoffFactor = src.BackoffFactor
	}
	if src.MaxDelay != 0 {
		dst.MaxDelay = src.MaxDelay
	}
	dst.AutoRecoveryEnabled = src.AutoRecoveryEnabled
	if src.RecoveryTimeout != 0 {
		dst.RecoveryTimeout = src.RecoveryTimeout
	}
}

func mergeLoggingConfig(dst, src *LoggingConfig) {
	if src.Level != "" {
		dst.Level = src.Level
	}
	if src.Format != "" {
		dst.Format = src.Format
	}
	if src.Output != "" {
		dst.Output = src.Output
	}
	if src.FilePath != "" {
		dst.FilePath = src.FilePath
	}
	if src.MaxSize != 0 {
		dst.MaxSize = src.MaxSize
	}
	if src.MaxBackups != 0 {
		dst.MaxBackups = src.MaxBackups
	}
	if src.MaxAge != 0 {
		dst.MaxAge = src.MaxAge
	}
	dst.Compress = src.Compress
}

func mergeMetricsConfig(dst, src *MetricsConfig) {
	dst.Enabled = src.Enabled
	if src.Port != 0 {
		dst.Port = src.Port
	}
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.CollectInterval != 0 {
		dst.CollectInterval = src.CollectInterval
	}
	if src.RetentionPeriod != 0 {
		dst.RetentionPeriod = src.RetentionPeriod
	}
	if src.Exporter != "" {
		dst.Exporter = src.Exporter
	}
	if src.ExporterConfig != nil {
		if dst.ExporterConfig == nil {
			dst.ExporterConfig = make(map[string]interface{})
		}
		for k, v := range src.ExporterConfig {
			dst.ExporterConfig[k] = v
		}
	}
}

func mergeSecurityConfig(dst, src *SecurityConfig) {
	dst.EnableAuth = src.EnableAuth
	if src.AuthProvider != "" {
		dst.AuthProvider = src.AuthProvider
	}
	if src.AuthConfig != nil {
		if dst.AuthConfig == nil {
			dst.AuthConfig = make(map[string]interface{})
		}
		for k, v := range src.AuthConfig {
			dst.AuthConfig[k] = v
		}
	}
	dst.EnableEncryption = src.EnableEncryption
	if src.EncryptionKey != "" {
		dst.EncryptionKey = src.EncryptionKey
	}
	dst.TLSEnabled = src.TLSEnabled
	if src.TLSCertFile != "" {
		dst.TLSCertFile = src.TLSCertFile
	}
	if src.TLSKeyFile != "" {
		dst.TLSKeyFile = src.TLSKeyFile
	}
	dst.CORSEnabled = src.CORSEnabled
	if len(src.CORSOrigins) > 0 {
		dst.CORSOrigins = make([]string, len(src.CORSOrigins))
		copy(dst.CORSOrigins, src.CORSOrigins)
	}
}