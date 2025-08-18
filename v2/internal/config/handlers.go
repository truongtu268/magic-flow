package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Handlers provides HTTP handlers for configuration management
type Handlers struct {
	manager         *Manager
	templateManager *TemplateManager
}

// NewHandlers creates new configuration handlers
func NewHandlers(manager *Manager, templateManager *TemplateManager) *Handlers {
	return &Handlers{
		manager:         manager,
		templateManager: templateManager,
	}
}

// RegisterRoutes registers configuration routes
func (h *Handlers) RegisterRoutes(router *mux.Router) {
	configRouter := router.PathPrefix("/config").Subrouter()

	// Configuration endpoints
	configRouter.HandleFunc("/current", h.GetCurrentConfig).Methods("GET")
	configRouter.HandleFunc("/current", h.UpdateConfig).Methods("PUT")
	configRouter.HandleFunc("/summary", h.GetConfigSummary).Methods("GET")
	configRouter.HandleFunc("/validate", h.ValidateConfig).Methods("POST")
	configRouter.HandleFunc("/reload", h.ReloadConfig).Methods("POST")
	configRouter.HandleFunc("/save", h.SaveConfig).Methods("POST")

	// Environment endpoints
	configRouter.HandleFunc("/environments", h.ListEnvironments).Methods("GET")
	configRouter.HandleFunc("/environments/{env}", h.GetEnvironmentConfig).Methods("GET")
	configRouter.HandleFunc("/environments/{env}/activate", h.SetEnvironment).Methods("POST")

	// Feature flag endpoints
	configRouter.HandleFunc("/features", h.GetFeatureFlags).Methods("GET")
	configRouter.HandleFunc("/features/{flag}", h.GetFeatureFlag).Methods("GET")
	configRouter.HandleFunc("/features/{flag}", h.SetFeatureFlag).Methods("PUT")
	configRouter.HandleFunc("/features/bulk", h.UpdateFeatureFlags).Methods("PUT")

	// Template endpoints
	templateRouter := configRouter.PathPrefix("/templates").Subrouter()
	templateRouter.HandleFunc("", h.ListTemplates).Methods("GET")
	templateRouter.HandleFunc("/{name}", h.GetTemplate).Methods("GET")
	templateRouter.HandleFunc("/{name}", h.SaveTemplate).Methods("PUT")
	templateRouter.HandleFunc("/{name}", h.DeleteTemplate).Methods("DELETE")
	templateRouter.HandleFunc("/{name}/apply", h.ApplyTemplate).Methods("POST")
	templateRouter.HandleFunc("/reload", h.ReloadTemplates).Methods("POST")
}

// Configuration handlers

// GetCurrentConfig returns the current configuration
func (h *Handlers) GetCurrentConfig(w http.ResponseWriter, r *http.Request) {
	config := h.manager.GetConfig()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode config: %v", err), http.StatusInternalServerError)
		return
	}
}

// UpdateConfig updates the current configuration
func (h *Handlers) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var newConfig Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.manager.UpdateConfig(&newConfig); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update config: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration updated successfully"})
}

// GetConfigSummary returns a summary of the current configuration
func (h *Handlers) GetConfigSummary(w http.ResponseWriter, r *http.Request) {
	summary := h.manager.GetConfigSummary()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode summary: %v", err), http.StatusInternalServerError)
		return
	}
}

// ValidateConfig validates a configuration
func (h *Handlers) ValidateConfig(w http.ResponseWriter, r *http.Request) {
	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	validationResult := &ValidationResult{
		Valid:      true,
		Errors:     []string{},
		Warnings:   []string{},
		ValidatedAt: time.Now(),
	}

	if err := validateConfig(&config); err != nil {
		validationResult.Valid = false
		validationResult.Errors = append(validationResult.Errors, err.Error())
	}

	// Add additional validation checks
	h.performAdditionalValidation(&config, validationResult)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(validationResult); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode validation result: %v", err), http.StatusInternalServerError)
		return
	}
}

// ReloadConfig reloads configuration from file
func (h *Handlers) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.ReloadConfig(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to reload config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration reloaded successfully"})
}

// SaveConfig saves the current configuration to file
func (h *Handlers) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.SaveCurrentConfig(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration saved successfully"})
}

// Environment handlers

// ListEnvironments returns available environments
func (h *Handlers) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	// This would typically query available environment configurations
	environments := []EnvironmentInfo{
		{Name: "development", Description: "Development environment", Active: h.manager.GetConfig().Environment == "development"},
		{Name: "staging", Description: "Staging environment", Active: h.manager.GetConfig().Environment == "staging"},
		{Name: "production", Description: "Production environment", Active: h.manager.GetConfig().Environment == "production"},
		{Name: "test", Description: "Test environment", Active: h.manager.GetConfig().Environment == "test"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(environments); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode environments: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetEnvironmentConfig returns configuration for a specific environment
func (h *Handlers) GetEnvironmentConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	env := vars["env"]

	config, err := h.manager.GetEnvironmentConfig(env)
	if err != nil {
		http.Error(w, fmt.Sprintf("Environment not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode config: %v", err), http.StatusInternalServerError)
		return
	}
}

// SetEnvironment switches to a different environment
func (h *Handlers) SetEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	env := vars["env"]

	if err := h.manager.SetEnvironment(env); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set environment: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "success",
		"message":     fmt.Sprintf("Environment switched to %s", env),
		"environment": env,
	})
}

// Feature flag handlers

// GetFeatureFlags returns all feature flags
func (h *Handlers) GetFeatureFlags(w http.ResponseWriter, r *http.Request) {
	flags := h.manager.GetAllFeatureFlags()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flags); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode feature flags: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetFeatureFlag returns a specific feature flag
func (h *Handlers) GetFeatureFlag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flag := vars["flag"]

	enabled := h.manager.GetFeatureFlag(flag)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"flag":    flag,
		"enabled": enabled,
	})
}

// SetFeatureFlag sets a feature flag value
func (h *Handlers) SetFeatureFlag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flag := vars["flag"]

	var request struct {
		Enabled bool `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.manager.SetFeatureFlag(flag, request.Enabled); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set feature flag: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Feature flag %s set to %v", flag, request.Enabled),
		"flag":    flag,
		"enabled": request.Enabled,
	})
}

// UpdateFeatureFlags updates multiple feature flags
func (h *Handlers) UpdateFeatureFlags(w http.ResponseWriter, r *http.Request) {
	var request map[string]bool
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	updated := make(map[string]bool)
	errors := make(map[string]string)

	for flag, enabled := range request {
		if err := h.manager.SetFeatureFlag(flag, enabled); err != nil {
			errors[flag] = err.Error()
		} else {
			updated[flag] = enabled
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"updated": updated,
		"errors":  errors,
	})
}

// Template handlers

// ListTemplates returns all configuration templates
func (h *Handlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates := h.templateManager.ListTemplates()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(templates); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode templates: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetTemplate returns a specific configuration template
func (h *Handlers) GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	template, err := h.templateManager.GetTemplate(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(template); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode template: %v", err), http.StatusInternalServerError)
		return
	}
}

// SaveTemplate saves a configuration template
func (h *Handlers) SaveTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var template ConfigTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	template.Name = name // Ensure name matches URL parameter

	if err := h.templateManager.SaveTemplate(&template); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Template %s saved successfully", name),
	})
}

// DeleteTemplate deletes a configuration template
func (h *Handlers) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// This would implement template deletion
	// For now, return not implemented
	http.Error(w, "Template deletion not implemented", http.StatusNotImplemented)
}

// ApplyTemplate applies a template to create a new configuration
func (h *Handlers) ApplyTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var request struct {
		Overrides map[string]interface{} `json:"overrides,omitempty"`
		Apply     bool                   `json:"apply,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	config, err := h.templateManager.ApplyTemplate(name, request.Overrides)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply template: %v", err), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"config": config,
	}

	// If apply is true, update the current configuration
	if request.Apply {
		if err := h.manager.UpdateConfig(config); err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply template config: %v", err), http.StatusBadRequest)
			return
		}
		response["message"] = fmt.Sprintf("Template %s applied successfully", name)
		response["applied"] = true
	} else {
		response["message"] = fmt.Sprintf("Template %s generated successfully (not applied)", name)
		response["applied"] = false
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// ReloadTemplates reloads all configuration templates
func (h *Handlers) ReloadTemplates(w http.ResponseWriter, r *http.Request) {
	if err := h.templateManager.LoadTemplates(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to reload templates: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Templates reloaded successfully",
	})
}

// Helper methods

func (h *Handlers) performAdditionalValidation(config *Config, result *ValidationResult) {
	// Check for potential security issues
	if config.Security.Authentication.Enabled && config.Security.Authentication.JWT.Secret == "" {
		result.Errors = append(result.Errors, "JWT secret is required when authentication is enabled")
		result.Valid = false
	}

	// Check for performance warnings
	if config.Database.MaxOpenConns > 100 {
		result.Warnings = append(result.Warnings, "High number of database connections may impact performance")
	}

	if config.Engine.MaxConcurrentWorkflows > 1000 {
		result.Warnings = append(result.Warnings, "High number of concurrent workflows may impact system resources")
	}

	// Check for development vs production settings
	if config.Environment == "production" {
		if config.Logging.Level == "debug" {
			result.Warnings = append(result.Warnings, "Debug logging enabled in production environment")
		}
		if !config.Server.TLS.Enabled {
			result.Warnings = append(result.Warnings, "TLS not enabled in production environment")
		}
		if !config.Security.Authentication.Enabled {
			result.Warnings = append(result.Warnings, "Authentication not enabled in production environment")
		}
	}

	// Check for feature flag consistency
	if config.Features.Authentication && !config.Security.Authentication.Enabled {
		result.Errors = append(result.Errors, "Authentication feature flag enabled but authentication not configured")
		result.Valid = false
	}

	if config.Features.Metrics && !config.Metrics.Enabled {
		result.Errors = append(result.Errors, "Metrics feature flag enabled but metrics not configured")
		result.Valid = false
	}
}

func getUserID(r *http.Request) string {
	// Extract user ID from request context or headers
	// This would typically be set by authentication middleware
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}
	return "anonymous"
}

func parseQueryBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// Response types

type ValidationResult struct {
	Valid       bool      `json:"valid"`
	Errors      []string  `json:"errors"`
	Warnings    []string  `json:"warnings"`
	ValidatedAt time.Time `json:"validated_at"`
}

type EnvironmentInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}