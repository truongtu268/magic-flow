package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"magic-flow/v2/pkg/models"
)

// Service provides code generation functionality
type Service struct {
	templateManager *TemplateManager
	handlers        map[Language]LanguageHandler
}

// NewService creates a new code generation service
func NewService() (*Service, error) {
	templateManager, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create template manager: %w", err)
	}

	service := &Service{
		templateManager: templateManager,
		handlers:        make(map[Language]LanguageHandler),
	}

	// Register language handlers
	service.handlers[LanguageGo] = NewGoHandler(templateManager)
	service.handlers[LanguageTypeScript] = NewTypeScriptHandler(templateManager)
	service.handlers[LanguagePython] = NewPythonHandler(templateManager)
	service.handlers[LanguageJava] = NewJavaHandler(templateManager)

	return service, nil
}

// GenerateCode generates code for a workflow in the specified language
func (s *Service) GenerateCode(workflow *models.Workflow, request *GenerationRequest) (*GenerationResult, error) {
	if workflow == nil {
		return nil, fmt.Errorf("workflow cannot be nil")
	}

	if request == nil {
		return nil, fmt.Errorf("generation request cannot be nil")
	}

	// Get language handler
	handler, exists := s.handlers[request.Language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", request.Language)
	}

	// Validate request
	if err := handler.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Prepare template data
	templateData, err := handler.PrepareTemplateData(workflow, request)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare template data: %w", err)
	}

	// Generate files
	files, err := handler.Generate(workflow, request, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	// Create result
	result := &GenerationResult{
		WorkflowID:   workflow.ID,
		WorkflowName: workflow.Name,
		Language:     request.Language,
		PackageName:  request.PackageName,
		Files:        files,
		Metadata: map[string]interface{}{
			"generated_at":    templateData.GeneratedAt,
			"workflow_version": workflow.Version,
			"file_count":       len(files),
			"include_tests":    request.IncludeTests,
			"options":          request.Options,
		},
	}

	return result, nil
}

// GetSupportedLanguages returns the list of supported languages
func (s *Service) GetSupportedLanguages() []Language {
	languages := make([]Language, 0, len(s.handlers))
	for lang := range s.handlers {
		languages = append(languages, lang)
	}
	return languages
}

// GetLanguageHandler returns the handler for a specific language
func (s *Service) GetLanguageHandler(language Language) (LanguageHandler, error) {
	handler, exists := s.handlers[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return handler, nil
}

// ValidateWorkflow validates that a workflow can be used for code generation
func (s *Service) ValidateWorkflow(workflow *models.Workflow) error {
	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	if workflow.Name == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}

	if workflow.Definition.Steps == nil || len(workflow.Definition.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	// Validate workflow steps
	for i, step := range workflow.Definition.Steps {
		if step.ID == "" {
			return fmt.Errorf("step %d must have an ID", i)
		}
		if step.Type == "" {
			return fmt.Errorf("step %d (%s) must have a type", i, step.ID)
		}
	}

	return nil
}

// GetDefaultRequest creates a default generation request for a language
func (s *Service) GetDefaultRequest(language Language, workflowName string) (*GenerationRequest, error) {
	handler, exists := s.handlers[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	request := &GenerationRequest{
		Language:     language,
		PackageName:  handler.GetDefaultPackageName(),
		IncludeTests: true,
		Options:      make(map[string]interface{}),
	}

	// Set language-specific defaults
	switch language {
	case LanguageGo:
		request.Options["module_name"] = "github.com/your-org/" + ToSnakeCase(workflowName) + "-client"
	case LanguageTypeScript:
		request.Options["npm_package_name"] = "@your-org/" + ToKebabCase(workflowName) + "-client"
		request.Options["version"] = "1.0.0"
	case LanguagePython:
		request.Options["package_name"] = ToSnakeCase(workflowName) + "_client"
		request.Options["version"] = "1.0.0"
	case LanguageJava:
		request.Options["group_id"] = "com.magicflow"
		request.Options["artifact_id"] = ToSnakeCase(workflowName) + "-client"
		request.Options["version"] = "1.0.0"
	}

	return request, nil
}

// GetFileStructure returns the expected file structure for a language
func (s *Service) GetFileStructure(language Language, packageName string) (map[string]string, error) {
	handler, exists := s.handlers[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	structure := make(map[string]string)
	extension := handler.GetFileExtension()

	switch language {
	case LanguageGo:
		structure["client"+extension] = "Main client implementation"
		structure["models"+extension] = "Data models and types"
		structure["types"+extension] = "Type definitions"
		structure["go.mod"] = "Go module definition"
		structure["README.md"] = "Documentation"
		if packageName != "" {
			structure["client_test"+extension] = "Client tests"
		}

	case LanguageTypeScript:
		structure["client"+extension] = "Main client implementation"
		structure["types"+extension] = "Type definitions"
		structure["models"+extension] = "Data models"
		structure["index"+extension] = "Package entry point"
		structure["package.json"] = "NPM package definition"
		structure["tsconfig.json"] = "TypeScript configuration"
		structure["README.md"] = "Documentation"
		if packageName != "" {
			structure["client.test"+extension] = "Client tests"
		}

	case LanguagePython:
		structure["__init__.py"] = "Package initialization"
		structure["client.py"] = "Main client implementation"
		structure["models.py"] = "Data models"
		structure["types.py"] = "Type definitions"
		structure["exceptions.py"] = "Exception classes"
		structure["setup.py"] = "Package setup script"
		structure["requirements.txt"] = "Dependencies"
		structure["pyproject.toml"] = "Modern Python packaging"
		structure["README.md"] = "Documentation"
		if packageName != "" {
			structure["test_client.py"] = "Client tests"
		}

	case LanguageJava:
		packagePath := strings.ReplaceAll(packageName, ".", "/")
		structure[filepath.Join("src/main/java", packagePath, "Client.java")] = "Main client implementation"
		structure[filepath.Join("src/main/java", packagePath, "models")] = "Data models directory"
		structure[filepath.Join("src/main/java", packagePath, "exceptions")] = "Exception classes directory"
		structure[filepath.Join("src/main/java", packagePath, "config")] = "Configuration classes directory"
		structure["pom.xml"] = "Maven project definition"
		structure["build.gradle"] = "Gradle build script"
		structure["README.md"] = "Documentation"
		if packageName != "" {
			structure[filepath.Join("src/test/java", packagePath, "ClientTest.java")] = "Client tests"
		}
	}

	return structure, nil
}

// GetTemplateNames returns available template names for a language
func (s *Service) GetTemplateNames(language Language) ([]string, error) {
	return s.templateManager.GetTemplateNames(string(language))
}

// GetTemplate returns a specific template for a language
func (s *Service) GetTemplate(language Language, templateName string) (string, error) {
	return s.templateManager.GetTemplate(string(language), templateName)
}

// AddCustomTemplate adds a custom template for a language
func (s *Service) AddCustomTemplate(language Language, templateName, content string) error {
	return s.templateManager.AddCustomTemplate(string(language), templateName, content)
}

// GeneratePreview generates a preview of the code without creating files
func (s *Service) GeneratePreview(workflow *models.Workflow, request *GenerationRequest, maxFiles int) (*GenerationResult, error) {
	result, err := s.GenerateCode(workflow, request)
	if err != nil {
		return nil, err
	}

	// Limit the number of files in preview
	if maxFiles > 0 && len(result.Files) > maxFiles {
		result.Files = result.Files[:maxFiles]
		result.Metadata["preview_truncated"] = true
		result.Metadata["total_files"] = len(result.Files)
	}

	// Truncate file content for preview
	for i := range result.Files {
		if len(result.Files[i].Content) > 2000 {
			result.Files[i].Content = result.Files[i].Content[:2000] + "\n\n// ... content truncated for preview ..."
		}
	}

	return result, nil
}

// GetGenerationStats returns statistics about code generation
func (s *Service) GetGenerationStats(workflow *models.Workflow, language Language) (map[string]interface{}, error) {
	handler, exists := s.handlers[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Create a minimal request for analysis
	request := &GenerationRequest{
		Language:     language,
		PackageName:  handler.GetDefaultPackageName(),
		IncludeTests: true,
	}

	templateData, err := handler.PrepareTemplateData(workflow, request)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare template data: %w", err)
	}

	stats := map[string]interface{}{
		"language":        language,
		"workflow_name":   workflow.Name,
		"workflow_steps":  len(workflow.Definition.Steps),
		"methods_count":   len(templateData.Methods),
		"models_count":    len(templateData.Models),
		"imports_count":   len(templateData.Imports),
		"package_name":    templateData.PackageName,
		"class_name":      templateData.ClassName,
		"file_extension": handler.GetFileExtension(),
	}

	// Estimate file count
	files, err := handler.Generate(workflow, request, templateData)
	if err == nil {
		stats["estimated_files"] = len(files)
		stats["file_types"] = s.getFileTypes(files)
	}

	return stats, nil
}

// getFileTypes returns a summary of file types
func (s *Service) getFileTypes(files []GeneratedFile) map[string]int {
	types := make(map[string]int)
	for _, file := range files {
		types[file.Type]++
	}
	return types
}

// ValidateGenerationRequest validates a generation request
func (s *Service) ValidateGenerationRequest(request *GenerationRequest) error {
	if request == nil {
		return fmt.Errorf("generation request cannot be nil")
	}

	if request.Language == "" {
		return fmt.Errorf("language must be specified")
	}

	// Check if language is supported
	if _, exists := s.handlers[request.Language]; !exists {
		return fmt.Errorf("unsupported language: %s", request.Language)
	}

	// Validate with language-specific handler
	handler := s.handlers[request.Language]
	return handler.ValidateRequest(request)
}

// GetLanguageInfo returns information about a specific language
func (s *Service) GetLanguageInfo(language Language) (map[string]interface{}, error) {
	handler, exists := s.handlers[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	templateNames, err := s.templateManager.GetTemplateNames(string(language))
	if err != nil {
		templateNames = []string{}
	}

	info := map[string]interface{}{
		"language":            language,
		"file_extension":      handler.GetFileExtension(),
		"default_package":     handler.GetDefaultPackageName(),
		"available_templates": templateNames,
		"supported_features": s.getSupportedFeatures(language),
	}

	return info, nil
}

// getSupportedFeatures returns the features supported by a language
func (s *Service) getSupportedFeatures(language Language) []string {
	features := []string{"client_generation", "model_generation", "documentation"}

	switch language {
	case LanguageGo:
		features = append(features, "modules", "interfaces", "error_handling")
	case LanguageTypeScript:
		features = append(features, "npm_package", "type_definitions", "async_await")
	case LanguagePython:
		features = append(features, "pip_package", "type_hints", "async_support")
	case LanguageJava:
		features = append(features, "maven_support", "gradle_support", "annotations")
	}

	return features
}