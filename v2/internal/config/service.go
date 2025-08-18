package config

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Service provides a unified interface for configuration operations
type Service struct {
	manager         *Manager
	templateManager *TemplateManager
	mu              sync.RWMutex
	listeners       []ConfigChangeListener
	metrics         *ServiceMetrics
	ctx             context.Context
	cancel          context.CancelFunc
}

// ConfigChangeListener defines the interface for configuration change notifications
type ConfigChangeListener interface {
	OnConfigChanged(oldConfig, newConfig *Config) error
}

// ServiceMetrics tracks configuration service metrics
type ServiceMetrics struct {
	ConfigReloads     int64     `json:"config_reloads"`
	TemplateApplies   int64     `json:"template_applies"`
	ValidationErrors  int64     `json:"validation_errors"`
	FeatureFlagChanges int64    `json:"feature_flag_changes"`
	LastReload        time.Time `json:"last_reload"`
	LastValidation    time.Time `json:"last_validation"`
	Uptime            time.Time `json:"uptime"`
}

// NewService creates a new configuration service
func NewService(configPath string) (*Service, error) {
	manager, err := NewManager(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	templateManager, err := NewTemplateManager("configs/templates")
	if err != nil {
		return nil, fmt.Errorf("failed to create template manager: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		manager:         manager,
		templateManager: templateManager,
		listeners:       make([]ConfigChangeListener, 0),
		metrics: &ServiceMetrics{
			Uptime: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// Add service as a change listener to the manager
	manager.AddChangeListener(service.onConfigChange)

	return service, nil
}

// Start starts the configuration service
func (s *Service) Start() error {
	// Start any background processes if needed
	return nil
}

// Stop stops the configuration service
func (s *Service) Stop() error {
	s.cancel()
	return nil
}

// Configuration operations

// GetConfig returns the current configuration
func (s *Service) GetConfig() *Config {
	return s.manager.GetConfig()
}

// UpdateConfig updates the configuration with validation and notifications
func (s *Service) UpdateConfig(newConfig *Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate the new configuration
	if err := s.validateConfigWithMetrics(newConfig); err != nil {
		s.metrics.ValidationErrors++
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	oldConfig := s.manager.GetConfig()

	// Update the configuration
	if err := s.manager.UpdateConfig(newConfig); err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	// Notify listeners
	for _, listener := range s.listeners {
		if err := listener.OnConfigChanged(oldConfig, newConfig); err != nil {
			// Log error but don't fail the update
			fmt.Printf("Config change listener error: %v\n", err)
		}
	}

	return nil
}

// ReloadConfig reloads configuration from file
func (s *Service) ReloadConfig() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.manager.ReloadConfig(); err != nil {
		return fmt.Errorf("failed to reload configuration: %w", err)
	}

	s.metrics.ConfigReloads++
	s.metrics.LastReload = time.Now()

	return nil
}

// ValidateConfig validates a configuration without applying it
func (s *Service) ValidateConfig(config *Config) *ValidationResult {
	result := &ValidationResult{
		Valid:       true,
		Errors:      []string{},
		Warnings:    []string{},
		ValidatedAt: time.Now(),
	}

	if err := validateConfig(config); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, err.Error())
	}

	// Perform additional validation
	s.performServiceValidation(config, result)

	s.metrics.LastValidation = time.Now()
	if !result.Valid {
		s.metrics.ValidationErrors++
	}

	return result
}

// GetConfigSummary returns a summary of the current configuration
func (s *Service) GetConfigSummary() *ConfigSummary {
	return s.manager.GetConfigSummary()
}

// Environment operations

// SetEnvironment switches to a different environment
func (s *Service) SetEnvironment(environment string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.manager.SetEnvironment(environment)
}

// GetEnvironmentConfig returns configuration for a specific environment
func (s *Service) GetEnvironmentConfig(environment string) (*Config, error) {
	return s.manager.GetEnvironmentConfig(environment)
}

// ListEnvironments returns available environments
func (s *Service) ListEnvironments() []string {
	return []string{"development", "staging", "production", "test"}
}

// Feature flag operations

// GetFeatureFlag returns the value of a feature flag
func (s *Service) GetFeatureFlag(flag string) bool {
	return s.manager.GetFeatureFlag(flag)
}

// SetFeatureFlag sets a feature flag value
func (s *Service) SetFeatureFlag(flag string, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.manager.SetFeatureFlag(flag, enabled); err != nil {
		return fmt.Errorf("failed to set feature flag: %w", err)
	}

	s.metrics.FeatureFlagChanges++
	return nil
}

// GetAllFeatureFlags returns all feature flags
func (s *Service) GetAllFeatureFlags() map[string]bool {
	return s.manager.GetAllFeatureFlags()
}

// UpdateFeatureFlags updates multiple feature flags
func (s *Service) UpdateFeatureFlags(flags map[string]bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for flag, enabled := range flags {
		if err := s.manager.SetFeatureFlag(flag, enabled); err != nil {
			return fmt.Errorf("failed to set feature flag %s: %w", flag, err)
		}
		s.metrics.FeatureFlagChanges++
	}

	return nil
}

// Template operations

// ListTemplates returns all available configuration templates
func (s *Service) ListTemplates() []*ConfigTemplate {
	return s.templateManager.ListTemplates()
}

// GetTemplate returns a specific configuration template
func (s *Service) GetTemplate(name string) (*ConfigTemplate, error) {
	return s.templateManager.GetTemplate(name)
}

// ApplyTemplate applies a template to generate a configuration
func (s *Service) ApplyTemplate(name string, overrides map[string]interface{}) (*Config, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	config, err := s.templateManager.ApplyTemplate(name, overrides)
	if err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	s.metrics.TemplateApplies++
	return config, nil
}

// ApplyTemplateAndUpdate applies a template and updates the current configuration
func (s *Service) ApplyTemplateAndUpdate(name string, overrides map[string]interface{}) error {
	config, err := s.ApplyTemplate(name, overrides)
	if err != nil {
		return err
	}

	return s.UpdateConfig(config)
}

// SaveTemplate saves a configuration template
func (s *Service) SaveTemplate(template *ConfigTemplate) error {
	return s.templateManager.SaveTemplate(template)
}

// Listener management

// AddChangeListener adds a configuration change listener
func (s *Service) AddChangeListener(listener ConfigChangeListener) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listeners = append(s.listeners, listener)
}

// RemoveChangeListener removes a configuration change listener
func (s *Service) RemoveChangeListener(listener ConfigChangeListener) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, l := range s.listeners {
		if l == listener {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			break
		}
	}
}

// Metrics and monitoring

// GetMetrics returns service metrics
func (s *Service) GetMetrics() *ServiceMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &ServiceMetrics{
		ConfigReloads:      s.metrics.ConfigReloads,
		TemplateApplies:    s.metrics.TemplateApplies,
		ValidationErrors:   s.metrics.ValidationErrors,
		FeatureFlagChanges: s.metrics.FeatureFlagChanges,
		LastReload:         s.metrics.LastReload,
		LastValidation:     s.metrics.LastValidation,
		Uptime:             s.metrics.Uptime,
	}
}

// GetHealthStatus returns the health status of the configuration service
func (s *Service) GetHealthStatus() *HealthStatus {
	config := s.GetConfig()
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    make(map[string]interface{}),
	}

	// Check configuration validity
	validationResult := s.ValidateConfig(config)
	if !validationResult.Valid {
		status.Status = "unhealthy"
		status.Checks["config_validation"] = map[string]interface{}{
			"status": "failed",
			"errors": validationResult.Errors,
		}
	} else {
		status.Checks["config_validation"] = map[string]interface{}{
			"status": "passed",
		}
	}

	// Check file watcher status
	if s.manager.watcherActive {
		status.Checks["file_watcher"] = map[string]interface{}{
			"status": "active",
		}
	} else {
		status.Checks["file_watcher"] = map[string]interface{}{
			"status": "inactive",
		}
	}

	// Check template availability
	templates := s.ListTemplates()
	status.Checks["templates"] = map[string]interface{}{
		"status": "available",
		"count":  len(templates),
	}

	// Add metrics
	status.Checks["metrics"] = s.GetMetrics()

	return status
}

// ResetMetrics resets service metrics
func (s *Service) ResetMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics = &ServiceMetrics{
		Uptime: s.metrics.Uptime, // Keep uptime
	}
}

// Backup and restore operations

// BackupConfig creates a backup of the current configuration
func (s *Service) BackupConfig() (*ConfigBackup, error) {
	config := s.GetConfig()
	backup := &ConfigBackup{
		Config:    *config,
		Timestamp: time.Now(),
		Version:   "1.0",
		Metadata: map[string]interface{}{
			"environment": config.Environment,
			"created_by":  "config-service",
		},
	}

	return backup, nil
}

// RestoreConfig restores configuration from a backup
func (s *Service) RestoreConfig(backup *ConfigBackup) error {
	return s.UpdateConfig(&backup.Config)
}

// Private methods

func (s *Service) onConfigChange(oldConfig, newConfig *Config) {
	// This is called when the manager detects a configuration change
	// We can add service-level logic here if needed
}

func (s *Service) validateConfigWithMetrics(config *Config) error {
	result := s.ValidateConfig(config)
	if !result.Valid {
		return fmt.Errorf("validation failed: %v", result.Errors)
	}
	return nil
}

func (s *Service) performServiceValidation(config *Config, result *ValidationResult) {
	// Service-level validation logic
	
	// Check for required services based on feature flags
	if config.Features.Dashboard && config.Dashboard.Port == 0 {
		result.Errors = append(result.Errors, "Dashboard feature enabled but no port configured")
		result.Valid = false
	}

	if config.Features.CodeGeneration && config.CodeGeneration.OutputDir == "" {
		result.Errors = append(result.Errors, "Code generation feature enabled but no output directory configured")
		result.Valid = false
	}

	if config.Features.Versioning && config.Versioning.RetentionDays == 0 {
		result.Warnings = append(result.Warnings, "Versioning feature enabled but no retention policy configured")
	}

	// Check for resource limits
	if config.Engine.MaxConcurrentWorkflows == 0 {
		result.Warnings = append(result.Warnings, "No limit set for concurrent workflows")
	}

	if config.Database.MaxOpenConns == 0 {
		result.Warnings = append(result.Warnings, "No limit set for database connections")
	}

	// Check for security configurations
	if config.Environment == "production" {
		if config.Security.RateLimit.RequestsPerMinute == 0 {
			result.Warnings = append(result.Warnings, "No rate limiting configured for production")
		}

		if len(config.Security.CORS.AllowedOrigins) == 0 {
			result.Warnings = append(result.Warnings, "No CORS origins configured for production")
		}
	}
}

// Helper types

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
}

type ConfigBackup struct {
	Config    Config                 `json:"config"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Metadata  map[string]interface{} `json:"metadata"`
}