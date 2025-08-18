package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure for Magic Flow v2
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Security SecurityConfig `mapstructure:"security"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Features FeatureConfig  `mapstructure:"features"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host" default:"0.0.0.0"`
	Port         int           `mapstructure:"port" default:"8080"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" default:"60s"`
	TLS          TLSConfig     `mapstructure:"tls"`
	CORS         CORSConfig    `mapstructure:"cors"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" default:"false"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	Enabled        bool     `mapstructure:"enabled" default:"true"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" default:"postgres"`
	Host            string        `mapstructure:"host" default:"localhost"`
	Port            int           `mapstructure:"port" default:"5432"`
	Database        string        `mapstructure:"database" default:"magicflow"`
	Username        string        `mapstructure:"username" default:"postgres"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode" default:"disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"5m"`
	Migrations      MigrationConfig `mapstructure:"migrations"`
}

// MigrationConfig contains database migration configuration
type MigrationConfig struct {
	Enabled   bool   `mapstructure:"enabled" default:"true"`
	Directory string `mapstructure:"directory" default:"./migrations"`
	AutoRun   bool   `mapstructure:"auto_run" default:"false"`
}

// CacheConfig contains Redis cache configuration
type CacheConfig struct {
	Enabled  bool          `mapstructure:"enabled" default:"true"`
	Host     string        `mapstructure:"host" default:"localhost"`
	Port     int           `mapstructure:"port" default:"6379"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db" default:"0"`
	TTL      time.Duration `mapstructure:"ttl" default:"1h"`
	Prefix   string        `mapstructure:"prefix" default:"magicflow:"`
}

// SecurityConfig contains security-related configuration
type SecurityConfig struct {
	JWT      JWTConfig      `mapstructure:"jwt"`
	API      APIKeyConfig   `mapstructure:"api"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration" default:"24h"`
	Issuer     string        `mapstructure:"issuer" default:"magicflow"`
}

// APIKeyConfig contains API key configuration
type APIKeyConfig struct {
	Enabled bool     `mapstructure:"enabled" default:"false"`
	Keys    []string `mapstructure:"keys"`
	Header  string   `mapstructure:"header" default:"X-API-Key"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled" default:"true"`
	RPS     int  `mapstructure:"rps" default:"100"`
	Burst   int  `mapstructure:"burst" default:"200"`
}

// MetricsConfig contains metrics and monitoring configuration
type MetricsConfig struct {
	Enabled    bool          `mapstructure:"enabled" default:"true"`
	Path       string        `mapstructure:"path" default:"/metrics"`
	Interval   time.Duration `mapstructure:"interval" default:"15s"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

// PrometheusConfig contains Prometheus-specific configuration
type PrometheusConfig struct {
	Enabled   bool   `mapstructure:"enabled" default:"true"`
	Namespace string `mapstructure:"namespace" default:"magicflow"`
	Subsystem string `mapstructure:"subsystem" default:"v2"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level" default:"info"`
	Format string `mapstructure:"format" default:"json"`
	Output string `mapstructure:"output" default:"stdout"`
	File   string `mapstructure:"file"`
}

// FeatureConfig contains feature flags
type FeatureConfig struct {
	CodeGeneration bool `mapstructure:"code_generation" default:"true"`
	Dashboard      bool `mapstructure:"dashboard" default:"true"`
	Versioning     bool `mapstructure:"versioning" default:"true"`
	Webhooks       bool `mapstructure:"webhooks" default:"true"`
	Metrics        bool `mapstructure:"metrics" default:"true"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/etc/magicflow")
	}
	
	// Environment variable support
	viper.SetEnvPrefix("MAGICFLOW")
	viper.AutomaticEnv()
	
	// Set defaults
	setDefaults()
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults and environment variables
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")
	
	// Database defaults
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "magicflow")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	
	// Cache defaults
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.host", "localhost")
	viper.SetDefault("cache.port", 6379)
	viper.SetDefault("cache.db", 0)
	viper.SetDefault("cache.ttl", "1h")
	viper.SetDefault("cache.prefix", "magicflow:")
	
	// Security defaults
	viper.SetDefault("security.jwt.expiration", "24h")
	viper.SetDefault("security.jwt.issuer", "magicflow")
	viper.SetDefault("security.rate_limit.enabled", true)
	viper.SetDefault("security.rate_limit.rps", 100)
	viper.SetDefault("security.rate_limit.burst", 200)
	
	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.interval", "15s")
	viper.SetDefault("metrics.prometheus.enabled", true)
	viper.SetDefault("metrics.prometheus.namespace", "magicflow")
	viper.SetDefault("metrics.prometheus.subsystem", "v2")
	
	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	
	// Feature flags defaults
	viper.SetDefault("features.code_generation", true)
	viper.SetDefault("features.dashboard", true)
	viper.SetDefault("features.versioning", true)
	viper.SetDefault("features.webhooks", true)
	viper.SetDefault("features.metrics", true)
}

// validate validates the configuration
func validate(config *Config) error {
	// Validate server configuration
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	
	// Validate database configuration
	if config.Database.Driver != "postgres" && config.Database.Driver != "mysql" {
		return fmt.Errorf("unsupported database driver: %s", config.Database.Driver)
	}
	
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	
	// Validate JWT secret if JWT is used
	if config.Security.JWT.Secret == "" {
		config.Security.JWT.Secret = os.Getenv("JWT_SECRET")
		if config.Security.JWT.Secret == "" {
			return fmt.Errorf("JWT secret is required")
		}
	}
	
	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	default:
		return ""
	}
}

// GetRedisAddr returns the Redis connection address
func (c *CacheConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetServerAddr returns the server address
func (c *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}