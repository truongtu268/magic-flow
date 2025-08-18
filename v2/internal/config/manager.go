package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Manager handles configuration management with hot reloading and environment-specific configs
type Manager struct {
	mu           sync.RWMutex
	config       *Config
	configPath   string
	watcher      *fsnotify.Watcher
	listeners    []ConfigChangeListener
	ctx          context.Context
	cancel       context.CancelFunc
	environments map[string]*Config
}

// ConfigChangeListener is called when configuration changes
type ConfigChangeListener func(oldConfig, newConfig *Config) error

// NewManager creates a new configuration manager
func NewManager(configPath string) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &Manager{
		configPath:   configPath,
		ctx:          ctx,
		cancel:       cancel,
		listeners:    make([]ConfigChangeListener, 0),
		environments: make(map[string]*Config),
	}

	// Load initial configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}
	manager.config = config

	// Load environment-specific configurations
	if err := manager.loadEnvironmentConfigs(); err != nil {
		return nil, fmt.Errorf("failed to load environment configs: %w", err)
	}

	// Set up file watcher for hot reloading
	if err := manager.setupWatcher(); err != nil {
		return nil, fmt.Errorf("failed to setup config watcher: %w", err)
	}

	return manager, nil
}

// GetConfig returns the current configuration (thread-safe)
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// GetEnvironmentConfig returns configuration for a specific environment
func (m *Manager) GetEnvironmentConfig(env string) (*Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if config, exists := m.environments[env]; exists {
		return config, nil
	}

	return nil, fmt.Errorf("configuration for environment '%s' not found", env)
}

// UpdateConfig updates the configuration and notifies listeners
func (m *Manager) UpdateConfig(newConfig *Config) error {
	m.mu.Lock()
	oldConfig := m.config
	m.mu.Unlock()

	// Validate new configuration
	if err := validateConfig(newConfig); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Notify listeners before updating
	for _, listener := range m.listeners {
		if err := listener(oldConfig, newConfig); err != nil {
			return fmt.Errorf("config change listener failed: %w", err)
		}
	}

	// Update configuration
	m.mu.Lock()
	m.config = newConfig
	m.mu.Unlock()

	return nil
}

// ReloadConfig reloads configuration from file
func (m *Manager) ReloadConfig() error {
	newConfig, err := LoadConfig(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	return m.UpdateConfig(newConfig)
}

// AddConfigChangeListener adds a listener for configuration changes
func (m *Manager) AddConfigChangeListener(listener ConfigChangeListener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = append(m.listeners, listener)
}

// SaveCurrentConfig saves the current configuration to file
func (m *Manager) SaveCurrentConfig() error {
	m.mu.RLock()
	config := m.config
	m.mu.RUnlock()

	return SaveConfig(config, m.configPath)
}

// SetEnvironment switches to a different environment configuration
func (m *Manager) SetEnvironment(env string) error {
	envConfig, err := m.GetEnvironmentConfig(env)
	if err != nil {
		return err
	}

	// Create a copy and update environment
	newConfig := *envConfig
	newConfig.Environment = env

	return m.UpdateConfig(&newConfig)
}

// GetFeatureFlag returns the value of a feature flag
func (m *Manager) GetFeatureFlag(flag string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.IsFeatureEnabled(flag)
}

// SetFeatureFlag sets a feature flag value
func (m *Manager) SetFeatureFlag(flag string, enabled bool) error {
	m.mu.Lock()
	newConfig := *m.config
	m.mu.Unlock()

	// Update the specific feature flag
	switch flag {
	case "workflow_versioning":
		newConfig.Features.WorkflowVersioning = enabled
	case "code_generation":
		newConfig.Features.CodeGeneration = enabled
	case "dashboard":
		newConfig.Features.Dashboard = enabled
	case "metrics":
		newConfig.Features.Metrics = enabled
	case "authentication":
		newConfig.Features.Authentication = enabled
	case "authorization":
		newConfig.Features.Authorization = enabled
	case "rate_limit":
		newConfig.Features.RateLimit = enabled
	case "encryption":
		newConfig.Features.Encryption = enabled
	case "audit_log":
		newConfig.Features.AuditLog = enabled
	case "backup":
		newConfig.Features.Backup = enabled
	case "clustering":
		newConfig.Features.Clustering = enabled
	case "advanced_workflows":
		newConfig.Features.AdvancedWorkflows = enabled
	default:
		return fmt.Errorf("unknown feature flag: %s", flag)
	}

	return m.UpdateConfig(&newConfig)
}

// GetAllFeatureFlags returns all feature flags
func (m *Manager) GetAllFeatureFlags() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]bool{
		"workflow_versioning": m.config.Features.WorkflowVersioning,
		"code_generation":     m.config.Features.CodeGeneration,
		"dashboard":           m.config.Features.Dashboard,
		"metrics":             m.config.Features.Metrics,
		"authentication":      m.config.Features.Authentication,
		"authorization":       m.config.Features.Authorization,
		"rate_limit":          m.config.Features.RateLimit,
		"encryption":          m.config.Features.Encryption,
		"audit_log":           m.config.Features.AuditLog,
		"backup":              m.config.Features.Backup,
		"clustering":          m.config.Features.Clustering,
		"advanced_workflows":  m.config.Features.AdvancedWorkflows,
	}
}

// ValidateCurrentConfig validates the current configuration
func (m *Manager) ValidateCurrentConfig() error {
	m.mu.RLock()
	config := m.config
	m.mu.RUnlock()

	return validateConfig(config)
}

// GetConfigSummary returns a summary of the current configuration
func (m *Manager) GetConfigSummary() *ConfigSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return &ConfigSummary{
		Environment:     m.config.Environment,
		ServerAddress:   m.config.Server.GetAddress(),
		DatabaseDriver:  m.config.Database.Driver,
		DatabaseHost:    m.config.Database.Host,
		TLSEnabled:      m.config.Server.TLS.Enabled,
		AuthEnabled:     m.config.Security.Authentication.Enabled,
		MetricsEnabled:  m.config.Metrics.Enabled,
		DashboardEnabled: m.config.Dashboard.Enabled,
		FeatureFlags:    m.GetAllFeatureFlags(),
		LoadedAt:        time.Now(),
	}
}

// Close stops the configuration manager and cleans up resources
func (m *Manager) Close() error {
	m.cancel()

	if m.watcher != nil {
		return m.watcher.Close()
	}

	return nil
}

// Private methods

func (m *Manager) setupWatcher() error {
	if m.configPath == "" {
		return nil // No file to watch
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	m.watcher = watcher

	// Watch the config file
	if err := watcher.Add(m.configPath); err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	// Start watching in a goroutine
	go m.watchConfigFile()

	return nil
}

func (m *Manager) watchConfigFile() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}

			// Only reload on write events
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Add a small delay to avoid multiple rapid reloads
				time.Sleep(100 * time.Millisecond)

				if err := m.ReloadConfig(); err != nil {
					// Log error but don't stop watching
					// In a real implementation, this would use a proper logger
					fmt.Printf("Failed to reload config: %v\n", err)
				}
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			// Log error but continue watching
			fmt.Printf("Config watcher error: %v\n", err)
		}
	}
}

func (m *Manager) loadEnvironmentConfigs() error {
	if m.configPath == "" {
		return nil
	}

	configDir := filepath.Dir(m.configPath)
	baseConfigName := filepath.Base(m.configPath)
	extension := filepath.Ext(baseConfigName)
	baseName := baseConfigName[:len(baseConfigName)-len(extension)]

	// Look for environment-specific config files
	environments := []string{"development", "staging", "production", "test"}

	for _, env := range environments {
		envConfigPath := filepath.Join(configDir, fmt.Sprintf("%s.%s%s", baseName, env, extension))

		if fileExists(envConfigPath) {
			envConfig, err := LoadConfig(envConfigPath)
			if err != nil {
				return fmt.Errorf("failed to load %s config: %w", env, err)
			}
			envConfig.Environment = env
			m.environments[env] = envConfig
		}
	}

	return nil
}

// ConfigSummary provides a summary of the current configuration
type ConfigSummary struct {
	Environment      string            `json:"environment"`
	ServerAddress    string            `json:"server_address"`
	DatabaseDriver   string            `json:"database_driver"`
	DatabaseHost     string            `json:"database_host"`
	TLSEnabled       bool              `json:"tls_enabled"`
	AuthEnabled      bool              `json:"auth_enabled"`
	MetricsEnabled   bool              `json:"metrics_enabled"`
	DashboardEnabled bool              `json:"dashboard_enabled"`
	FeatureFlags     map[string]bool   `json:"feature_flags"`
	LoadedAt         time.Time         `json:"loaded_at"`
}

// ConfigTemplate represents a configuration template
type ConfigTemplate struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	Environment string                 `yaml:"environment" json:"environment"`
	Config      map[string]interface{} `yaml:"config" json:"config"`
	Tags        []string               `yaml:"tags" json:"tags"`
	CreatedAt   time.Time              `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `yaml:"updated_at" json:"updated_at"`
}

// TemplateManager manages configuration templates
type TemplateManager struct {
	templatesDir string
	templates    map[string]*ConfigTemplate
	mu           sync.RWMutex
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templatesDir string) *TemplateManager {
	return &TemplateManager{
		templatesDir: templatesDir,
		templates:    make(map[string]*ConfigTemplate),
	}
}

// LoadTemplates loads all configuration templates
func (tm *TemplateManager) LoadTemplates() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if !dirExists(tm.templatesDir) {
		return fmt.Errorf("templates directory does not exist: %s", tm.templatesDir)
	}

	return filepath.Walk(tm.templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			template, err := tm.loadTemplate(path)
			if err != nil {
				return fmt.Errorf("failed to load template %s: %w", path, err)
			}
			tm.templates[template.Name] = template
		}

		return nil
	})
}

// GetTemplate returns a configuration template by name
func (tm *TemplateManager) GetTemplate(name string) (*ConfigTemplate, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	template, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}

	return template, nil
}

// ListTemplates returns all available templates
func (tm *TemplateManager) ListTemplates() []*ConfigTemplate {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	templates := make([]*ConfigTemplate, 0, len(tm.templates))
	for _, template := range tm.templates {
		templates = append(templates, template)
	}

	return templates
}

// ApplyTemplate applies a template to create a new configuration
func (tm *TemplateManager) ApplyTemplate(templateName string, overrides map[string]interface{}) (*Config, error) {
	template, err := tm.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}

	// Start with default config
	config := DefaultConfig()

	// Apply template configuration
	if err := tm.applyConfigMap(config, template.Config); err != nil {
		return nil, fmt.Errorf("failed to apply template config: %w", err)
	}

	// Apply overrides
	if overrides != nil {
		if err := tm.applyConfigMap(config, overrides); err != nil {
			return nil, fmt.Errorf("failed to apply config overrides: %w", err)
		}
	}

	// Set environment from template
	config.Environment = template.Environment

	return config, nil
}

// SaveTemplate saves a configuration template
func (tm *TemplateManager) SaveTemplate(template *ConfigTemplate) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	template.UpdatedAt = time.Now()
	if template.CreatedAt.IsZero() {
		template.CreatedAt = template.UpdatedAt
	}

	// Save to file
	templatePath := filepath.Join(tm.templatesDir, template.Name+".yaml")
	if err := tm.saveTemplateToFile(template, templatePath); err != nil {
		return fmt.Errorf("failed to save template to file: %w", err)
	}

	// Update in memory
	tm.templates[template.Name] = template

	return nil
}

// Private methods for TemplateManager

func (tm *TemplateManager) loadTemplate(path string) (*ConfigTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	var template ConfigTemplate
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template file: %w", err)
	}

	return &template, nil
}

func (tm *TemplateManager) saveTemplateToFile(template *ConfigTemplate, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	// Marshal template to YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	return nil
}

func (tm *TemplateManager) applyConfigMap(config *Config, configMap map[string]interface{}) error {
	// This is a simplified implementation
	// In a real implementation, you would use reflection or a more sophisticated
	// mapping mechanism to apply the configuration map to the config struct

	// For now, we'll handle a few common cases
	if env, ok := configMap["environment"]; ok {
		if envStr, ok := env.(string); ok {
			config.Environment = envStr
		}
	}

	if server, ok := configMap["server"]; ok {
		if serverMap, ok := server.(map[string]interface{}); ok {
			if host, ok := serverMap["host"]; ok {
				if hostStr, ok := host.(string); ok {
					config.Server.Host = hostStr
				}
			}
			if port, ok := serverMap["port"]; ok {
				if portInt, ok := port.(int); ok {
					config.Server.Port = portInt
				}
			}
		}
	}

	// Add more field mappings as needed

	return nil
}