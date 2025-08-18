package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"magic-flow/v2/pkg/models"
)

// GoHandler implements LanguageHandler for Go code generation
type GoHandler struct {
	templateManager *TemplateManager
}

// NewGoHandler creates a new Go language handler
func NewGoHandler(templateManager *TemplateManager) *GoHandler {
	return &GoHandler{
		templateManager: templateManager,
	}
}

// Generate generates Go code for a workflow
func (h *GoHandler) Generate(workflow *models.Workflow, request *GenerationRequest, templateData *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate client file
	clientFile, err := h.generateClientFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client file: %w", err)
	}
	files = append(files, clientFile)

	// Generate models file
	modelsFile, err := h.generateModelsFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate models file: %w", err)
	}
	files = append(files, modelsFile)

	// Generate types file
	typesFile, err := h.generateTypesFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate types file: %w", err)
	}
	files = append(files, typesFile)

	// Generate test file if requested
	if request.IncludeTests {
		testFile, err := h.generateTestFile(templateData)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test file: %w", err)
		}
		files = append(files, testFile)
	}

	// Generate go.mod file
	goModFile, err := h.generateGoModFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate go.mod file: %w", err)
	}
	files = append(files, goModFile)

	// Generate README file
	readmeFile, err := h.generateReadmeFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate README file: %w", err)
	}
	files = append(files, readmeFile)

	return files, nil
}

// ValidateRequest validates Go-specific generation request
func (h *GoHandler) ValidateRequest(request *GenerationRequest) error {
	if request.PackageName == "" {
		request.PackageName = h.GetDefaultPackageName()
	}

	// Validate package name format
	if !isValidGoPackageName(request.PackageName) {
		return fmt.Errorf("invalid Go package name: %s", request.PackageName)
	}

	return nil
}

// PrepareTemplateData prepares template data for Go code generation
func (h *GoHandler) PrepareTemplateData(workflow *models.Workflow, request *GenerationRequest) (*TemplateData, error) {
	packageName := request.PackageName
	if packageName == "" {
		packageName = h.GetDefaultPackageName()
	}

	className := ToPascalCase(workflow.Name) + "Client"

	// Extract methods from workflow steps
	methods := ExtractStepMethods(workflow)

	// Generate imports
	imports := h.generateImports(workflow, request)

	// Generate models
	models := h.generateModels(workflow)

	templateData := &TemplateData{
		Workflow:    workflow,
		PackageName: packageName,
		ClassName:   className,
		Imports:     imports,
		Methods:     methods,
		Models:      models,
		Options:     request.Options,
		GeneratedAt: workflow.CreatedAt,
	}

	return templateData, nil
}

// GetFileExtension returns the file extension for Go files
func (h *GoHandler) GetFileExtension() string {
	return ".go"
}

// GetDefaultPackageName returns the default package name for Go
func (h *GoHandler) GetDefaultPackageName() string {
	return "magicflow"
}

// generateClientFile generates the main client file
func (h *GoHandler) generateClientFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("go", "client")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "client.go"),
		Content:  content,
		Language: "go",
		Type:     "client",
	}, nil
}

// generateModelsFile generates the models file
func (h *GoHandler) generateModelsFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("go", "models")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "models.go"),
		Content:  content,
		Language: "go",
		Type:     "models",
	}, nil
}

// generateTypesFile generates the types file
func (h *GoHandler) generateTypesFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("go", "types")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "types.go"),
		Content:  content,
		Language: "go",
		Type:     "types",
	}, nil
}

// generateTestFile generates the test file
func (h *GoHandler) generateTestFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("go", "test")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "client_test.go"),
		Content:  content,
		Language: "go",
		Type:     "test",
	}, nil
}

// generateGoModFile generates the go.mod file
func (h *GoHandler) generateGoModFile(data *TemplateData) (GeneratedFile, error) {
	moduleName := data.PackageName
	if data.Options != nil {
		if module, ok := data.Options["module_name"].(string); ok && module != "" {
			moduleName = module
		}
	}

	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`, moduleName)

	return GeneratedFile{
		Path:     "go.mod",
		Content:  content,
		Language: "go",
		Type:     "config",
	}, nil
}

// generateReadmeFile generates the README file
func (h *GoHandler) generateReadmeFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`# %s Go Client

Generated Go client library for the %s workflow.

## Installation

` + "```bash" + `
go get %s
` + "```" + `

## Usage

` + "```go" + `
package main

import (
	"context"
	"fmt"
	"log"

	"%s"
)

func main() {
	client := %s.New%s("http://localhost:8080", "your-api-key")
	
	ctx := context.Background()
	input := map[string]interface{}{
		"key": "value",
	}
	
	result, err := client.ExecuteWorkflow(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Execution ID: %%s\n", result.ID)
	fmt.Printf("Status: %%s\n", result.Status)
}
` + "```" + `

## API Reference

### Client Methods

#### ExecuteWorkflow

Executes the %s workflow with the provided input.

` + "```go" + `
func (c *%s) ExecuteWorkflow(ctx context.Context, input map[string]interface{}) (*ExecutionResult, error)
` + "```" + `

#### GetExecutionStatus

Retrieves the status of a workflow execution.

` + "```go" + `
func (c *%s) GetExecutionStatus(ctx context.Context, executionID uuid.UUID) (*ExecutionStatus, error)
` + "```" + `

%s

## Models

### ExecutionResult

Represents the result of a workflow execution.

### ExecutionStatus

Represents the status of a workflow execution.

### StepStatus

Represents the status of a workflow step.

## Constants

- ` + "`WORKFLOW_ID`" + `: The ID of the workflow
- ` + "`WORKFLOW_NAME`" + `: The name of the workflow
- ` + "`Status.*`" + `: Execution status constants
- ` + "`Steps.*`" + `: Step ID constants

## Error Handling

All methods return an error as the second return value. Always check for errors:

` + "```go" + `
result, err := client.ExecuteWorkflow(ctx, input)
if err != nil {
	// Handle error
	log.Printf("Error executing workflow: %%v", err)
	return
}
` + "```" + `

## License

Generated code - see original workflow license.
`,
		data.Workflow.Name,
		data.Workflow.Name,
		data.PackageName,
		data.PackageName,
		data.PackageName,
		data.ClassName,
		data.Workflow.Name,
		data.ClassName,
		data.ClassName,
		h.generateMethodDocs(data.Methods),
	)

	return GeneratedFile{
		Path:     "README.md",
		Content:  content,
		Language: "markdown",
		Type:     "documentation",
	}, nil
}

// generateImports generates the list of imports needed
func (h *GoHandler) generateImports(workflow *models.Workflow, request *GenerationRequest) []string {
	imports := []string{
		"bytes",
		"context",
		"encoding/json",
		"fmt",
		"net/http",
		"time",
		"github.com/google/uuid",
	}

	// Add test imports if tests are included
	if request.IncludeTests {
		imports = append(imports, "testing", "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require")
	}

	return imports
}

// generateModels generates model definitions from workflow
func (h *GoHandler) generateModels(workflow *models.Workflow) []ModelData {
	var models []ModelData

	// Generate models based on workflow inputs/outputs
	if workflow.Definition.Input != nil {
		for key, schema := range workflow.Definition.Input {
			model := ModelData{
				Name:        ToPascalCase(key) + "Input",
				Description: fmt.Sprintf("Input model for %s", key),
				Fields:      h.generateFieldsFromSchema(schema),
			}
			models = append(models, model)
		}
	}

	if workflow.Definition.Output != nil {
		for key, schema := range workflow.Definition.Output {
			model := ModelData{
				Name:        ToPascalCase(key) + "Output",
				Description: fmt.Sprintf("Output model for %s", key),
				Fields:      h.generateFieldsFromSchema(schema),
			}
			models = append(models, model)
		}
	}

	return models
}

// generateFieldsFromSchema generates field definitions from schema
func (h *GoHandler) generateFieldsFromSchema(schema interface{}) []FieldData {
	var fields []FieldData

	// This is a simplified implementation
	// In a real scenario, you would parse the JSON schema properly
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
			for fieldName, fieldSchema := range properties {
				field := FieldData{
					Name:        ToPascalCase(fieldName),
					Type:        h.mapSchemaTypeToGoType(fieldSchema),
					Description: h.getSchemaDescription(fieldSchema),
					Required:    h.isFieldRequired(fieldName, schemaMap),
					Tags: map[string]string{
						"json": ToSnakeCase(fieldName),
					},
				}
				fields = append(fields, field)
			}
		}
	}

	return fields
}

// mapSchemaTypeToGoType maps JSON schema types to Go types
func (h *GoHandler) mapSchemaTypeToGoType(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if schemaType, ok := schemaMap["type"].(string); ok {
			switch schemaType {
			case "string":
				return "string"
			case "integer":
				return "int64"
			case "number":
				return "float64"
			case "boolean":
				return "bool"
			case "array":
				return "[]interface{}"
			case "object":
				return "map[string]interface{}"
			}
		}
	}
	return "interface{}"
}

// getSchemaDescription extracts description from schema
func (h *GoHandler) getSchemaDescription(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if desc, ok := schemaMap["description"].(string); ok {
			return desc
		}
	}
	return ""
}

// isFieldRequired checks if a field is required
func (h *GoHandler) isFieldRequired(fieldName string, schema map[string]interface{}) bool {
	if required, ok := schema["required"].([]interface{}); ok {
		for _, req := range required {
			if reqStr, ok := req.(string); ok && reqStr == fieldName {
				return true
			}
		}
	}
	return false
}

// generateMethodDocs generates documentation for methods
func (h *GoHandler) generateMethodDocs(methods []MethodData) string {
	if len(methods) == 0 {
		return ""
	}

	var docs strings.Builder
	docs.WriteString("\n### Step Methods\n\n")

	for _, method := range methods {
		docs.WriteString(fmt.Sprintf("#### %s\n\n", method.Name))
		if method.Description != "" {
			docs.WriteString(fmt.Sprintf("%s\n\n", method.Description))
		}

		docs.WriteString("```go\n")
		docs.WriteString(fmt.Sprintf("func (c *%sClient) %s(ctx context.Context", ToPascalCase(method.StepID), method.Name))
		for _, param := range method.Parameters {
			docs.WriteString(fmt.Sprintf(", %s %s", param.Name, param.Type))
		}
		docs.WriteString(fmt.Sprintf(") (%s, error)\n", method.ReturnType))
		docs.WriteString("```\n\n")
	}

	return docs.String()
}

// isValidGoPackageName validates Go package name
func isValidGoPackageName(name string) bool {
	if name == "" {
		return false
	}

	// Go package names should be lowercase and contain only letters, numbers, and underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	// Should not start with a number
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	return true
}