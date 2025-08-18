package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `yaml:"server" json:"server"`

	// Database configuration
	Database DatabaseConfig `yaml:"database" json:"database"`

	// Engine configuration
	Engine EngineConfig `yaml:"engine" json:"engine"`

	// Dashboard configuration
	Dashboard DashboardConfig `yaml:"dashboard" json:"dashboard"`

	// Code generation configuration
	CodeGen CodeGenConfig `yaml:"codegen" json:"codegen"`

	// Versioning configuration
	Versioning VersioningConfig `yaml:"versioning" json:"versioning"`

	// Security configuration
	Security SecurityConfig `yaml:"security" json:"security"`

	// Logging configuration
	Logging LoggingConfig `yaml:"logging" json:"logging"`

	// Metrics configuration
	Metrics MetricsConfig `yaml:"metrics" json:"metrics"`

	// Feature flags
	Features FeatureFlags `yaml:"features" json:"features"`

	// Environment-specific settings
	Environment string `yaml:"environment" json:"environment"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	TLS          TLSConfig     `yaml:"tls" json:"tls"`
	CORS         CORSConfig    `yaml:"cors" json:"cors"`
}

// TLSConfig contains TLS/SSL configuration
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	CertFile string `yaml:"cert_file" json:"cert_file"`
	KeyFile  string `yaml:"key_file" json:"key_file"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled" json:"enabled"`
	AllowedOrigins   []string `yaml:"allowed_origins" json:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" json:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers" json:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials" json:"allow_credentials"`
	MaxAge           int      `yaml:"max_age" json:"max_age"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Driver          string        `yaml:"driver" json:"driver"`
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	Database        string        `yaml:"database" json:"database"`
	Username        string        `yaml:"username" json:"username"`
	Password        string        `yaml:"password" json:"password"`
	SSLMode         string        `yaml:"ssl_mode" json:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
	Migrations      MigrationConfig `yaml:"migrations" json:"migrations"`
}

// MigrationConfig contains database migration configuration
type MigrationConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Directory string `yaml:"directory" json:"directory"`
	Table     string `yaml:"table" json:"table"`
}

// EngineConfig contains workflow engine configuration
type EngineConfig struct {
	MaxConcurrentWorkflows int           `yaml:"max_concurrent_workflows" json:"max_concurrent_workflows"`
	MaxConcurrentSteps     int           `yaml:"max_concurrent_steps" json:"max_concurrent_steps"`
	StepTimeout            time.Duration `yaml:"step_timeout" json:"step_timeout"`
	WorkflowTimeout        time.Duration `yaml:"workflow_timeout" json:"workflow_timeout"`
	RetryPolicy            RetryPolicy   `yaml:"retry_policy" json:"retry_policy"`
	Storage                StorageConfig `yaml:"storage" json:"storage"`
}

// RetryPolicy contains retry configuration
type RetryPolicy struct {
	MaxRetries      int           `yaml:"max_retries" json:"max_retries"`
	InitialDelay    time.Duration `yaml:"initial_delay" json:"initial_delay"`
	MaxDelay        time.Duration `yaml:"max_delay" json:"max_delay"`
	BackoffFactor   float64       `yaml:"backoff_factor" json:"backoff_factor"`
	RetryableErrors []string      `yaml:"retryable_errors" json:"retryable_errors"`
}

// StorageConfig contains storage configuration
type StorageConfig struct {
	Type   string                 `yaml:"type" json:"type"`
	Config map[string]interface{} `yaml:"config" json:"config"`
}

// DashboardConfig contains dashboard configuration
type DashboardConfig struct {
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	RefreshInterval time.Duration `yaml:"refresh_interval" json:"refresh_interval"`
	MaxConnections  int           `yaml:"max_connections" json:"max_connections"`
	WebSocket       WebSocketConfig `yaml:"websocket" json:"websocket"`
	UI              UIConfig      `yaml:"ui" json:"ui"`
}

// WebSocketConfig contains WebSocket configuration
type WebSocketConfig struct {
	Enabled       bool          `yaml:"enabled" json:"enabled"`
	PingInterval  time.Duration `yaml:"ping_interval" json:"ping_interval"`
	PongTimeout   time.Duration `yaml:"pong_timeout" json:"pong_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout" json:"write_timeout"`
	ReadTimeout   time.Duration `yaml:"read_timeout" json:"read_timeout"`
	BufferSize    int           `yaml:"buffer_size" json:"buffer_size"`
	MaxMessageSize int64        `yaml:"max_message_size" json:"max_message_size"`
}

// UIConfig contains UI configuration
type UIConfig struct {
	Theme         string            `yaml:"theme" json:"theme"`
	Title         string            `yaml:"title" json:"title"`
	Logo          string            `yaml:"logo" json:"logo"`
	CustomCSS     string            `yaml:"custom_css" json:"custom_css"`
	CustomJS      string            `yaml:"custom_js" json:"custom_js"`
	Features      map[string]bool   `yaml:"features" json:"features"`
	Settings      map[string]string `yaml:"settings" json:"settings"`
}

// CodeGenConfig contains code generation configuration
type CodeGenConfig struct {
	Enabled           bool              `yaml:"enabled" json:"enabled"`
	TemplatesDir      string            `yaml:"templates_dir" json:"templates_dir"`
	OutputDir         string            `yaml:"output_dir" json:"output_dir"`
	SupportedLanguages []string         `yaml:"supported_languages" json:"supported_languages"`
	LanguageConfigs   map[string]LanguageConfig `yaml:"language_configs" json:"language_configs"`
}

// LanguageConfig contains language-specific configuration
type LanguageConfig struct {
	Enabled      bool              `yaml:"enabled" json:"enabled"`
	TemplateDir  string            `yaml:"template_dir" json:"template_dir"`
	FileExtension string           `yaml:"file_extension" json:"file_extension"`
	PackageFormat string           `yaml:"package_format" json:"package_format"`
	Options      map[string]string `yaml:"options" json:"options"`
}

// VersioningConfig contains versioning configuration
type VersioningConfig struct {
	Enabled               bool          `yaml:"enabled" json:"enabled"`
	AutoVersioning        bool          `yaml:"auto_versioning" json:"auto_versioning"`
	VersioningStrategy    string        `yaml:"versioning_strategy" json:"versioning_strategy"`
	MigrationTimeout      time.Duration `yaml:"migration_timeout" json:"migration_timeout"`
	MaxRollbackDepth      int           `yaml:"max_rollback_depth" json:"max_rollback_depth"`
	BackupBeforeMigration bool          `yaml:"backup_before_migration" json:"backup_before_migration"`
	RetentionPolicy       RetentionPolicy `yaml:"retention_policy" json:"retention_policy"`
}

// RetentionPolicy contains version retention configuration
type RetentionPolicy struct {
	MaxVersions        int           `yaml:"max_versions" json:"max_versions"`
	RetentionPeriod    time.Duration `yaml:"retention_period" json:"retention_period"`
	KeepActiveVersions bool          `yaml:"keep_active_versions" json:"keep_active_versions"`
	KeepTaggedVersions bool          `yaml:"keep_tagged_versions" json:"keep_tagged_versions"`
	ArchiveOldVersions bool          `yaml:"archive_old_versions" json:"archive_old_versions"`
}

// SecurityConfig contains security configuration
type SecurityConfig struct {
	Authentication AuthConfig `yaml:"authentication" json:"authentication"`
	Authorization  AuthzConfig `yaml:"authorization" json:"authorization"`
	Encryption     EncryptionConfig `yaml:"encryption" json:"encryption"`
	RateLimit      RateLimitConfig `yaml:"rate_limit" json:"rate_limit"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Enabled  bool              `yaml:"enabled" json:"enabled"`
	Provider string            `yaml:"provider" json:"provider"`
	JWT      JWTConfig         `yaml:"jwt" json:"jwt"`
	OAuth    OAuthConfig       `yaml:"oauth" json:"oauth"`
	LDAP     LDAPConfig        `yaml:"ldap" json:"ldap"`
	Options  map[string]string `yaml:"options" json:"options"`
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret         string        `yaml:"secret" json:"secret"`
	Expiration     time.Duration `yaml:"expiration" json:"expiration"`
	RefreshEnabled bool          `yaml:"refresh_enabled" json:"refresh_enabled"`
	RefreshExpiration time.Duration `yaml:"refresh_expiration" json:"refresh_expiration"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
	Audience       string        `yaml:"audience" json:"audience"`
}

// OAuthConfig contains OAuth configuration
type OAuthConfig struct {
	ClientID     string   `yaml:"client_id" json:"client_id"`
	ClientSecret string   `yaml:"client_secret" json:"client_secret"`
	RedirectURL  string   `yaml:"redirect_url" json:"redirect_url"`
	Scopes       []string `yaml:"scopes" json:"scopes"`
	AuthURL      string   `yaml:"auth_url" json:"auth_url"`
	TokenURL     string   `yaml:"token_url" json:"token_url"`
	UserInfoURL  string   `yaml:"user_info_url" json:"user_info_url"`
}

// LDAPConfig contains LDAP configuration
type LDAPConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	BindDN       string `yaml:"bind_dn" json:"bind_dn"`
	BindPassword string `yaml:"bind_password" json:"bind_password"`
	BaseDN       string `yaml:"base_dn" json:"base_dn"`
	UserFilter   string `yaml:"user_filter" json:"user_filter"`
	GroupFilter  string `yaml:"group_filter" json:"group_filter"`
	TLS          bool   `yaml:"tls" json:"tls"`
}

// AuthzConfig contains authorization configuration
type AuthzConfig struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Provider    string            `yaml:"provider" json:"provider"`
	PolicyFile  string            `yaml:"policy_file" json:"policy_file"`
	DefaultRole string            `yaml:"default_role" json:"default_role"`
	Roles       map[string]Role   `yaml:"roles" json:"roles"`
	Permissions map[string]string `yaml:"permissions" json:"permissions"`
}

// Role contains role configuration
type Role struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Permissions []string `yaml:"permissions" json:"permissions"`
	Inherits    []string `yaml:"inherits" json:"inherits"`
}

// EncryptionConfig contains encryption configuration
type EncryptionConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Algorithm string `yaml:"algorithm" json:"algorithm"`
	KeyFile   string `yaml:"key_file" json:"key_file"`
	KeySize   int    `yaml:"key_size" json:"key_size"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled    bool          `yaml:"enabled" json:"enabled"`
	Requests   int           `yaml:"requests" json:"requests"`
	Window     time.Duration `yaml:"window" json:"window"`
	Burst      int           `yaml:"burst" json:"burst"`
	SkipPaths  []string      `yaml:"skip_paths" json:"skip_paths"`
	Headers    []string      `yaml:"headers" json:"headers"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string        `yaml:"level" json:"level"`
	Format     string        `yaml:"format" json:"format"`
	Output     string        `yaml:"output" json:"output"`
	File       FileLogConfig `yaml:"file" json:"file"`
	Structured bool          `yaml:"structured" json:"structured"`
	Fields     map[string]string `yaml:"fields" json:"fields"`
}

// FileLogConfig contains file logging configuration
type FileLogConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Path       string `yaml:"path" json:"path"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled    bool              `yaml:"enabled" json:"enabled"`
	Provider   string            `yaml:"provider" json:"provider"`
	Endpoint   string            `yaml:"endpoint" json:"endpoint"`
	Interval   time.Duration     `yaml:"interval" json:"interval"`
	Namespace  string            `yaml:"namespace" json:"namespace"`
	Labels     map[string]string `yaml:"labels" json:"labels"`
	Prometheus PrometheusConfig  `yaml:"prometheus" json:"prometheus"`
}

// PrometheusConfig contains Prometheus-specific configuration
type PrometheusConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Path      string `yaml:"path" json:"path"`
	Namespace string `yaml:"namespace" json:"namespace"`
	Subsystem string `yaml:"subsystem" json:"subsystem"`
}

// FeatureFlags contains feature flag configuration
type FeatureFlags struct {
	WorkflowVersioning bool `yaml:"workflow_versioning" json:"workflow_versioning"`
	CodeGeneration     bool `yaml:"code_generation" json:"code_generation"`
	Dashboard          bool `yaml:"dashboard" json:"dashboard"`
	Metrics            bool `yaml:"metrics" json:"metrics"`
	Authentication     bool `yaml:"authentication" json:"authentication"`
	Authorization      bool `yaml:"authorization" json:"authorization"`
	RateLimit          bool `yaml:"rate_limit" json:"rate_limit"`
	Encryption         bool `yaml:"encryption" json:"encryption"`
	AuditLog           bool `yaml:"audit_log" json:"audit_log"`
	Backup             bool `yaml:"backup" json:"backup"`
	Clustering         bool `yaml:"clustering" json:"clustering"`
	AdvancedWorkflows  bool `yaml:"advanced_workflows" json:"advanced_workflows"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLS: TLSConfig{
				Enabled: false,
			},
			CORS: CORSConfig{
				Enabled:        true,
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"*"},
				MaxAge:         86400,
			},
		},
		Database: DatabaseConfig{
			Driver:          "postgres",
			Host:            "localhost",
			Port:            5432,
			Database:        "magicflow",
			Username:        "postgres",
			Password:        "password",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			Migrations: MigrationConfig{
				Enabled:   true,
				Directory: "migrations",
				Table:     "schema_migrations",
			},
		},
		Engine: EngineConfig{
			MaxConcurrentWorkflows: 100,
			MaxConcurrentSteps:     1000,
			StepTimeout:            5 * time.Minute,
			WorkflowTimeout:        30 * time.Minute,
			RetryPolicy: RetryPolicy{
				MaxRetries:    3,
				InitialDelay:  1 * time.Second,
				MaxDelay:      30 * time.Second,
				BackoffFactor: 2.0,
			},
			Storage: StorageConfig{
				Type: "database",
			},
		},
		Dashboard: DashboardConfig{
			Enabled:         true,
			RefreshInterval: 5 * time.Second,
			MaxConnections:  1000,
			WebSocket: WebSocketConfig{
				Enabled:        true,
				PingInterval:   30 * time.Second,
				PongTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				ReadTimeout:    60 * time.Second,
				BufferSize:     1024,
				MaxMessageSize: 1024 * 1024, // 1MB
			},
			UI: UIConfig{
				Theme: "default",
				Title: "Magic Flow Dashboard",
				Features: map[string]bool{
					"realtime_updates": true,
					"dark_mode":       true,
					"export":          true,
				},
			},
		},
		CodeGen: CodeGenConfig{
			Enabled:            true,
			TemplatesDir:       "internal/codegen/templates",
			OutputDir:          "generated",
			SupportedLanguages: []string{"go", "typescript", "python", "java"},
			LanguageConfigs: map[string]LanguageConfig{
				"go": {
					Enabled:       true,
					TemplateDir:   "go",
					FileExtension: ".go",
					PackageFormat: "module",
				},
				"typescript": {
					Enabled:       true,
					TemplateDir:   "typescript",
					FileExtension: ".ts",
					PackageFormat: "npm",
				},
				"python": {
					Enabled:       true,
					TemplateDir:   "python",
					FileExtension: ".py",
					PackageFormat: "pip",
				},
				"java": {
					Enabled:       true,
					TemplateDir:   "java",
					FileExtension: ".java",
					PackageFormat: "maven",
				},
			},
		},
		Versioning: VersioningConfig{
			Enabled:               true,
			AutoVersioning:        false,
			VersioningStrategy:    "semantic",
			MigrationTimeout:      30 * time.Minute,
			MaxRollbackDepth:      10,
			BackupBeforeMigration: true,
			RetentionPolicy: RetentionPolicy{
				MaxVersions:        50,
				RetentionPeriod:    365 * 24 * time.Hour,
				KeepActiveVersions: true,
				KeepTaggedVersions: true,
				ArchiveOldVersions: true,
			},
		},
		Security: SecurityConfig{
			Authentication: AuthConfig{
				Enabled:  false,
				Provider: "jwt",
				JWT: JWTConfig{
					Expiration:        24 * time.Hour,
					RefreshEnabled:    true,
					RefreshExpiration: 7 * 24 * time.Hour,
					Issuer:            "magic-flow",
				},
			},
			Authorization: AuthzConfig{
				Enabled:     false,
				Provider:    "rbac",
				DefaultRole: "user",
			},
			Encryption: EncryptionConfig{
				Enabled:   false,
				Algorithm: "AES-256-GCM",
				KeySize:   256,
			},
			RateLimit: RateLimitConfig{
				Enabled:  false,
				Requests: 1000,
				Window:   time.Hour,
				Burst:    100,
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			Structured: true,
			File: FileLogConfig{
				Enabled:    false,
				Path:       "logs/app.log",
				MaxSize:    100, // MB
				MaxBackups: 3,
				MaxAge:     28, // days
				Compress:   true,
			},
		},
		Metrics: MetricsConfig{
			Enabled:   true,
			Provider:  "prometheus",
			Interval:  15 * time.Second,
			Namespace: "magicflow",
			Prometheus: PrometheusConfig{
				Enabled:   true,
				Path:      "/metrics",
				Namespace: "magicflow",
				Subsystem: "api",
			},
		},
		Features: FeatureFlags{
			WorkflowVersioning: true,
			CodeGeneration:     true,
			Dashboard:          true,
			Metrics:            true,
			Authentication:     false,
			Authorization:      false,
			RateLimit:          false,
			Encryption:         false,
			AuditLog:           false,
			Backup:             false,
			Clustering:         false,
			AdvancedWorkflows:  true,
		},
		Environment: "development",
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// Load from file if provided
	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	if err := loadFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from YAML file
func loadFromFile(config *Config, configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) error {
	// Environment
	if env := os.Getenv("MAGIC_FLOW_ENV"); env != "" {
		config.Environment = env
	}

	// Server configuration
	if host := os.Getenv("MAGIC_FLOW_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("MAGIC_FLOW_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// Database configuration
	if dbHost := os.Getenv("MAGIC_FLOW_DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("MAGIC_FLOW_DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = p
		}
	}
	if dbName := os.Getenv("MAGIC_FLOW_DB_NAME"); dbName != "" {
		config.Database.Database = dbName
	}
	if dbUser := os.Getenv("MAGIC_FLOW_DB_USER"); dbUser != "" {
		config.Database.Username = dbUser
	}
	if dbPass := os.Getenv("MAGIC_FLOW_DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}
	if dbSSL := os.Getenv("MAGIC_FLOW_DB_SSL_MODE"); dbSSL != "" {
		config.Database.SSLMode = dbSSL
	}

	// Security configuration
	if jwtSecret := os.Getenv("MAGIC_FLOW_JWT_SECRET"); jwtSecret != "" {
		config.Security.Authentication.JWT.Secret = jwtSecret
	}

	// TLS configuration
	if tlsEnabled := os.Getenv("MAGIC_FLOW_TLS_ENABLED"); tlsEnabled != "" {
		config.Server.TLS.Enabled = strings.ToLower(tlsEnabled) == "true"
	}
	if certFile := os.Getenv("MAGIC_FLOW_TLS_CERT_FILE"); certFile != "" {
		config.Server.TLS.CertFile = certFile
	}
	if keyFile := os.Getenv("MAGIC_FLOW_TLS_KEY_FILE"); keyFile != "" {
		config.Server.TLS.KeyFile = keyFile
	}

	// Feature flags
	if auth := os.Getenv("MAGIC_FLOW_FEATURE_AUTH"); auth != "" {
		config.Features.Authentication = strings.ToLower(auth) == "true"
		config.Security.Authentication.Enabled = config.Features.Authentication
	}
	if authz := os.Getenv("MAGIC_FLOW_FEATURE_AUTHZ"); authz != "" {
		config.Features.Authorization = strings.ToLower(authz) == "true"
		config.Security.Authorization.Enabled = config.Features.Authorization
	}
	if metrics := os.Getenv("MAGIC_FLOW_FEATURE_METRICS"); metrics != "" {
		config.Features.Metrics = strings.ToLower(metrics) == "true"
		config.Metrics.Enabled = config.Features.Metrics
	}
	if dashboard := os.Getenv("MAGIC_FLOW_FEATURE_DASHBOARD"); dashboard != "" {
		config.Features.Dashboard = strings.ToLower(dashboard) == "true"
		config.Dashboard.Enabled = config.Features.Dashboard
	}

	// Logging configuration
	if logLevel := os.Getenv("MAGIC_FLOW_LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	if logFormat := os.Getenv("MAGIC_FLOW_LOG_FORMAT"); logFormat != "" {
		config.Logging.Format = logFormat
	}

	return nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate database configuration
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate TLS configuration
	if config.Server.TLS.Enabled {
		if config.Server.TLS.CertFile == "" {
			return fmt.Errorf("TLS cert file is required when TLS is enabled")
		}
		if config.Server.TLS.KeyFile == "" {
			return fmt.Errorf("TLS key file is required when TLS is enabled")
		}
		if !fileExists(config.Server.TLS.CertFile) {
			return fmt.Errorf("TLS cert file does not exist: %s", config.Server.TLS.CertFile)
		}
		if !fileExists(config.Server.TLS.KeyFile) {
			return fmt.Errorf("TLS key file does not exist: %s", config.Server.TLS.KeyFile)
		}
	}

	// Validate JWT configuration
	if config.Security.Authentication.Enabled && config.Security.Authentication.Provider == "jwt" {
		if config.Security.Authentication.JWT.Secret == "" {
			return fmt.Errorf("JWT secret is required when JWT authentication is enabled")
		}
	}

	// Validate code generation configuration
	if config.CodeGen.Enabled {
		if config.CodeGen.TemplatesDir == "" {
			return fmt.Errorf("code generation templates directory is required")
		}
		if !dirExists(config.CodeGen.TemplatesDir) {
			return fmt.Errorf("code generation templates directory does not exist: %s", config.CodeGen.TemplatesDir)
		}
	}

	// Validate logging configuration
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, config.Logging.Level) {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, config.Logging.Format) {
		return fmt.Errorf("invalid log format: %s", config.Logging.Format)
	}

	return nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal configuration to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Helper functions

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// GetAddress returns the server address
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsFeatureEnabled checks if a feature is enabled
func (c *Config) IsFeatureEnabled(feature string) bool {
	switch feature {
	case "workflow_versioning":
		return c.Features.WorkflowVersioning
	case "code_generation":
		return c.Features.CodeGeneration
	case "dashboard":
		return c.Features.Dashboard
	case "metrics":
		return c.Features.Metrics
	case "authentication":
		return c.Features.Authentication
	case "authorization":
		return c.Features.Authorization
	case "rate_limit":
		return c.Features.RateLimit
	case "encryption":
		return c.Features.Encryption
	case "audit_log":
		return c.Features.AuditLog
	case "backup":
		return c.Features.Backup
	case "clustering":
		return c.Features.Clustering
	case "advanced_workflows":
		return c.Features.AdvancedWorkflows
	default:
		return false
	}
}