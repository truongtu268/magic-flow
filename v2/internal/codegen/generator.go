package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"

	"magic-flow/v2/pkg/models"
)

// Language represents supported programming languages
type Language string

const (
	LanguageGo         Language = "go"
	LanguageTypeScript Language = "typescript"
	LanguagePython     Language = "python"
	LanguageJava       Language = "java"
)

// GenerationRequest represents a code generation request
type GenerationRequest struct {
	WorkflowID   uuid.UUID `json:"workflow_id"`
	Language     Language  `json:"language"`
	PackageName  string    `json:"package_name,omitempty"`
	Namespace    string    `json:"namespace,omitempty"`
	OutputDir    string    `json:"output_dir,omitempty"`
	IncludeTests bool      `json:"include_tests,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
}

// GenerationResult represents the result of code generation
type GenerationResult struct {
	ID          uuid.UUID              `json:"id"`
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	Language    Language               `json:"language"`
	Files       []GeneratedFile        `json:"files"`
	Metadata    map[string]interface{} `json:"metadata"`
	GeneratedAt time.Time              `json:"generated_at"`
	Status      string                 `json:"status"`
	Error       string                 `json:"error,omitempty"`
}

// GeneratedFile represents a generated code file
type GeneratedFile struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Language string `json:"language"`
	Type     string `json:"type"` // "client", "model", "test", etc.
}

// TemplateData represents data passed to templates
type TemplateData struct {
	Workflow    *models.Workflow
	PackageName string
	Namespace   string
	ClassName   string
	Imports     []string
	Methods     []MethodData
	Models      []ModelData
	Options     map[string]interface{}
	GeneratedAt time.Time
}

// MethodData represents method information for templates
type MethodData struct {
	Name        string
	Description string
	Parameters  []ParameterData
	ReturnType  string
	StepID      string
	StepType    string
}

// ParameterData represents parameter information for templates
type ParameterData struct {
	Name        string
	Type        string
	Description string
	Required    bool
	DefaultValue string
}

// ModelData represents model information for templates
type ModelData struct {
	Name        string
	Description string
	Fields      []FieldData
}

// FieldData represents field information for templates
type FieldData struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Tags        map[string]string
}

// Generator interface defines code generation capabilities
type Generator interface {
	Generate(workflow *models.Workflow, request *GenerationRequest) (*GenerationResult, error)
	ValidateRequest(request *GenerationRequest) error
	GetSupportedLanguages() []Language
	GetTemplates(language Language) (map[string]string, error)
}

// CodeGenerator implements the Generator interface
type CodeGenerator struct {
	templateManager *TemplateManager
	languageHandlers map[Language]LanguageHandler
}

// LanguageHandler interface for language-specific code generation
type LanguageHandler interface {
	Generate(workflow *models.Workflow, request *GenerationRequest, templateData *TemplateData) ([]GeneratedFile, error)
	ValidateRequest(request *GenerationRequest) error
	PrepareTemplateData(workflow *models.Workflow, request *GenerationRequest) (*TemplateData, error)
	GetFileExtension() string
	GetDefaultPackageName() string
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	templateManager := NewTemplateManager()
	
	generator := &CodeGenerator{
		templateManager:  templateManager,
		languageHandlers: make(map[Language]LanguageHandler),
	}

	// Register language handlers
	generator.languageHandlers[LanguageGo] = NewGoHandler(templateManager)
	generator.languageHandlers[LanguageTypeScript] = NewTypeScriptHandler(templateManager)
	generator.languageHandlers[LanguagePython] = NewPythonHandler(templateManager)
	generator.languageHandlers[LanguageJava] = NewJavaHandler(templateManager)

	return generator
}

// Generate generates code for a workflow
func (g *CodeGenerator) Generate(workflow *models.Workflow, request *GenerationRequest) (*GenerationResult, error) {
	// Validate the request
	if err := g.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Get the language handler
	handler, exists := g.languageHandlers[request.Language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", request.Language)
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
		ID:          uuid.New(),
		WorkflowID:  request.WorkflowID,
		Language:    request.Language,
		Files:       files,
		GeneratedAt: time.Now().UTC(),
		Status:      "success",
		Metadata: map[string]interface{}{
			"package_name":   templateData.PackageName,
			"namespace":      templateData.Namespace,
			"class_name":     templateData.ClassName,
			"file_count":     len(files),
			"include_tests":  request.IncludeTests,
			"workflow_name":  workflow.Name,
			"workflow_steps": len(workflow.Definition.Steps),
		},
	}

	return result, nil
}

// ValidateRequest validates a generation request
func (g *CodeGenerator) ValidateRequest(request *GenerationRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if request.WorkflowID == uuid.Nil {
		return fmt.Errorf("workflow ID is required")
	}

	if request.Language == "" {
		return fmt.Errorf("language is required")
	}

	// Check if language is supported
	if _, exists := g.languageHandlers[request.Language]; !exists {
		return fmt.Errorf("unsupported language: %s", request.Language)
	}

	// Validate language-specific requirements
	handler := g.languageHandlers[request.Language]
	if err := handler.ValidateRequest(request); err != nil {
		return fmt.Errorf("language validation failed: %w", err)
	}

	return nil
}

// GetSupportedLanguages returns the list of supported languages
func (g *CodeGenerator) GetSupportedLanguages() []Language {
	languages := make([]Language, 0, len(g.languageHandlers))
	for lang := range g.languageHandlers {
		languages = append(languages, lang)
	}
	return languages
}

// GetTemplates returns available templates for a language
func (g *CodeGenerator) GetTemplates(language Language) (map[string]string, error) {
	return g.templateManager.GetTemplatesForLanguage(string(language))
}

// Helper functions

// SanitizeIdentifier sanitizes a string to be used as an identifier
func SanitizeIdentifier(input string) string {
	// Remove special characters and replace with underscores
	sanitized := strings.ReplaceAll(input, "-", "_")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, ".", "_")
	
	// Remove any non-alphanumeric characters except underscores
	var result strings.Builder
	for _, r := range sanitized {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(input string) string {
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '.'
	})
	
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			if len(word) > 1 {
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}
	
	return result.String()
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(input string) string {
	pascal := ToPascalCase(input)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(input string) string {
	var result strings.Builder
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// GetFileNameForLanguage generates appropriate file name for the language
func GetFileNameForLanguage(baseName string, language Language) string {
	switch language {
	case LanguageGo:
		return ToSnakeCase(baseName) + ".go"
	case LanguageTypeScript:
		return ToCamelCase(baseName) + ".ts"
	case LanguagePython:
		return ToSnakeCase(baseName) + ".py"
	case LanguageJava:
		return ToPascalCase(baseName) + ".java"
	default:
		return baseName
	}
}

// RenderTemplate renders a template with the given data
func RenderTemplate(templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New("code").Funcs(template.FuncMap{
		"toPascalCase": ToPascalCase,
		"toCamelCase":  ToCamelCase,
		"toSnakeCase":  ToSnakeCase,
		"sanitize":     SanitizeIdentifier,
		"join":         strings.Join,
		"title":        strings.Title,
		"lower":        strings.ToLower,
		"upper":        strings.ToUpper,
	}).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ExtractStepMethods extracts method information from workflow steps
func ExtractStepMethods(workflow *models.Workflow) []MethodData {
	var methods []MethodData

	for _, step := range workflow.Definition.Steps {
		method := MethodData{
			Name:        ToCamelCase(step.ID),
			Description: step.Name,
			StepID:      step.ID,
			StepType:    step.Type,
			ReturnType:  "interface{}", // Default return type
		}

		// Extract parameters from step configuration
		if step.Config != nil {
			for key, value := range step.Config {
				param := ParameterData{
					Name:        ToCamelCase(key),
					Type:        inferTypeFromValue(value),
					Description: fmt.Sprintf("Parameter for %s", key),
					Required:    true,
				}
				method.Parameters = append(method.Parameters, param)
			}
		}

		methods = append(methods, method)
	}

	return methods
}

// inferTypeFromValue infers the type from a value
func inferTypeFromValue(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "number"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "any"
	}
}