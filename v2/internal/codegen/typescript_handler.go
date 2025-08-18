package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"magic-flow/v2/pkg/models"
)

// TypeScriptHandler implements LanguageHandler for TypeScript code generation
type TypeScriptHandler struct {
	templateManager *TemplateManager
}

// NewTypeScriptHandler creates a new TypeScript language handler
func NewTypeScriptHandler(templateManager *TemplateManager) *TypeScriptHandler {
	return &TypeScriptHandler{
		templateManager: templateManager,
	}
}

// Generate generates TypeScript code for a workflow
func (h *TypeScriptHandler) Generate(workflow *models.Workflow, request *GenerationRequest, templateData *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate client file
	clientFile, err := h.generateClientFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client file: %w", err)
	}
	files = append(files, clientFile)

	// Generate types file
	typesFile, err := h.generateTypesFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate types file: %w", err)
	}
	files = append(files, typesFile)

	// Generate models file
	modelsFile, err := h.generateModelsFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate models file: %w", err)
	}
	files = append(files, modelsFile)

	// Generate index file
	indexFile, err := h.generateIndexFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate index file: %w", err)
	}
	files = append(files, indexFile)

	// Generate test file if requested
	if request.IncludeTests {
		testFile, err := h.generateTestFile(templateData)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test file: %w", err)
		}
		files = append(files, testFile)
	}

	// Generate package.json file
	packageFile, err := h.generatePackageJsonFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate package.json file: %w", err)
	}
	files = append(files, packageFile)

	// Generate tsconfig.json file
	tsconfigFile, err := h.generateTsConfigFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tsconfig.json file: %w", err)
	}
	files = append(files, tsconfigFile)

	// Generate README file
	readmeFile, err := h.generateReadmeFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate README file: %w", err)
	}
	files = append(files, readmeFile)

	return files, nil
}

// ValidateRequest validates TypeScript-specific generation request
func (h *TypeScriptHandler) ValidateRequest(request *GenerationRequest) error {
	if request.PackageName == "" {
		request.PackageName = h.GetDefaultPackageName()
	}

	// Validate package name format (npm package naming rules)
	if !isValidNpmPackageName(request.PackageName) {
		return fmt.Errorf("invalid npm package name: %s", request.PackageName)
	}

	return nil
}

// PrepareTemplateData prepares template data for TypeScript code generation
func (h *TypeScriptHandler) PrepareTemplateData(workflow *models.Workflow, request *GenerationRequest) (*TemplateData, error) {
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

// GetFileExtension returns the file extension for TypeScript files
func (h *TypeScriptHandler) GetFileExtension() string {
	return ".ts"
}

// GetDefaultPackageName returns the default package name for TypeScript
func (h *TypeScriptHandler) GetDefaultPackageName() string {
	return "@magicflow/client"
}

// generateClientFile generates the main client file
func (h *TypeScriptHandler) generateClientFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("typescript", "client")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     "src/client.ts",
		Content:  content,
		Language: "typescript",
		Type:     "client",
	}, nil
}

// generateTypesFile generates the types file
func (h *TypeScriptHandler) generateTypesFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("typescript", "types")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     "src/types.ts",
		Content:  content,
		Language: "typescript",
		Type:     "types",
	}, nil
}

// generateModelsFile generates the models file
func (h *TypeScriptHandler) generateModelsFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("typescript", "models")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     "src/models.ts",
		Content:  content,
		Language: "typescript",
		Type:     "models",
	}, nil
}

// generateIndexFile generates the index file
func (h *TypeScriptHandler) generateIndexFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`/**
 * %s TypeScript Client
 * Generated at: %s
 */

export { %s } from './client';
export * from './types';
export * from './models';

// Re-export for convenience
export default %s;
`,
		data.Workflow.Name,
		data.GeneratedAt.Format("2006-01-02 15:04:05"),
		data.ClassName,
		data.ClassName,
	)

	return GeneratedFile{
		Path:     "src/index.ts",
		Content:  content,
		Language: "typescript",
		Type:     "index",
	}, nil
}

// generateTestFile generates the test file
func (h *TypeScriptHandler) generateTestFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("typescript", "test")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     "src/client.test.ts",
		Content:  content,
		Language: "typescript",
		Type:     "test",
	}, nil
}

// generatePackageJsonFile generates the package.json file
func (h *TypeScriptHandler) generatePackageJsonFile(data *TemplateData) (GeneratedFile, error) {
	version := "1.0.0"
	if data.Options != nil {
		if v, ok := data.Options["version"].(string); ok && v != "" {
			version = v
		}
	}

	author := "Magic Flow"
	if data.Options != nil {
		if a, ok := data.Options["author"].(string); ok && a != "" {
			author = a
		}
	}

	content := fmt.Sprintf(`{
  "name": "%s",
  "version": "%s",
  "description": "TypeScript client for %s workflow",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "test": "jest",
    "test:watch": "jest --watch",
    "lint": "eslint src/**/*.ts",
    "lint:fix": "eslint src/**/*.ts --fix",
    "prepublishOnly": "npm run build"
  },
  "keywords": [
    "magic-flow",
    "workflow",
    "client",
    "typescript"
  ],
  "author": "%s",
  "license": "MIT",
  "dependencies": {
    "axios": "^1.6.0",
    "uuid": "^9.0.0"
  },
  "devDependencies": {
    "@types/jest": "^29.5.0",
    "@types/node": "^20.0.0",
    "@types/uuid": "^9.0.0",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "eslint": "^8.0.0",
    "jest": "^29.5.0",
    "ts-jest": "^29.1.0",
    "typescript": "^5.0.0"
  },
  "files": [
    "dist/**/*"
  ],
  "repository": {
    "type": "git",
    "url": "https://github.com/your-org/your-repo.git"
  },
  "bugs": {
    "url": "https://github.com/your-org/your-repo/issues"
  },
  "homepage": "https://github.com/your-org/your-repo#readme"
}
`,
		data.PackageName,
		version,
		data.Workflow.Name,
		author,
	)

	return GeneratedFile{
		Path:     "package.json",
		Content:  content,
		Language: "json",
		Type:     "config",
	}, nil
}

// generateTsConfigFile generates the tsconfig.json file
func (h *TypeScriptHandler) generateTsConfigFile(data *TemplateData) (GeneratedFile, error) {
	content := `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "removeComments": false,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": false
  },
  "include": [
    "src/**/*"
  ],
  "exclude": [
    "node_modules",
    "dist",
    "**/*.test.ts"
  ]
}
`

	return GeneratedFile{
		Path:     "tsconfig.json",
		Content:  content,
		Language: "json",
		Type:     "config",
	}, nil
}

// generateReadmeFile generates the README file
func (h *TypeScriptHandler) generateReadmeFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`# %s TypeScript Client

Generated TypeScript client library for the %s workflow.

## Installation

` + "```bash" + `
npm install %s
# or
yarn add %s
` + "```" + `

## Usage

` + "```typescript" + `
import { %s } from '%s';

const client = new %s('http://localhost:8080', 'your-api-key');

async function executeWorkflow() {
  try {
    const input = {
      key: 'value'
    };
    
    const result = await client.executeWorkflow(input);
    console.log('Execution ID:', result.id);
    console.log('Status:', result.status);
    
    // Check execution status
    const status = await client.getExecutionStatus(result.id);
    console.log('Current status:', status.status);
    console.log('Progress:', status.progress);
  } catch (error) {
    console.error('Error executing workflow:', error);
  }
}

executeWorkflow();
` + "```" + `

## API Reference

### Client Methods

#### executeWorkflow

Executes the %s workflow with the provided input.

` + "```typescript" + `
executeWorkflow(input: Record<string, any>): Promise<ExecutionResult>
` + "```" + `

#### getExecutionStatus

Retrieves the status of a workflow execution.

` + "```typescript" + `
getExecutionStatus(executionId: string): Promise<ExecutionStatus>
` + "```" + `

#### cancelExecution

Cancels a running workflow execution.

` + "```typescript" + `
cancelExecution(executionId: string): Promise<void>
` + "```" + `

#### getExecutionResult

Retrieves the result of a completed workflow execution.

` + "```typescript" + `
getExecutionResult(executionId: string): Promise<ExecutionResult>
` + "```" + `

%s

## Types

### ExecutionResult

Represents the result of a workflow execution.

` + "```typescript" + `
interface ExecutionResult {
  id: string;
  workflowId: string;
  status: ExecutionStatus;
  input: Record<string, any>;
  output?: Record<string, any>;
  error?: string;
  startedAt: Date;
  completedAt?: Date;
  duration?: number;
}
` + "```" + `

### ExecutionStatus

Represents the status of a workflow execution.

` + "```typescript" + `
interface ExecutionStatus {
  id: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  progress: number;
  currentStep?: string;
  steps: StepStatus[];
  startedAt: Date;
  updatedAt: Date;
}
` + "```" + `

### StepStatus

Represents the status of a workflow step.

` + "```typescript" + `
interface StepStatus {
  id: string;
  name: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'skipped';
  input?: Record<string, any>;
  output?: Record<string, any>;
  error?: string;
  startedAt?: Date;
  completedAt?: Date;
  duration?: number;
}
` + "```" + `

## Constants

- ` + "`WORKFLOW_ID`" + `: The ID of the workflow
- ` + "`WORKFLOW_NAME`" + `: The name of the workflow
- ` + "`ExecutionStatus`" + `: Execution status enum
- ` + "`StepIds`" + `: Step ID constants

## Error Handling

All methods return promises that may reject with errors. Always use try-catch blocks:

` + "```typescript" + `
try {
  const result = await client.executeWorkflow(input);
  // Handle success
} catch (error) {
  // Handle error
  console.error('Error executing workflow:', error.message);
}
` + "```" + `

## Development

` + "```bash" + `
# Install dependencies
npm install

# Build the project
npm run build

# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Lint code
npm run lint

# Fix linting issues
npm run lint:fix
` + "```" + `

## License

Generated code - see original workflow license.
`,
		data.Workflow.Name,
		data.Workflow.Name,
		data.PackageName,
		data.PackageName,
		data.ClassName,
		data.PackageName,
		data.ClassName,
		data.Workflow.Name,
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
func (h *TypeScriptHandler) generateImports(workflow *models.Workflow, request *GenerationRequest) []string {
	imports := []string{
		"axios",
		"{ v4 as uuidv4 } from 'uuid'",
	}

	// Add test imports if tests are included
	if request.IncludeTests {
		imports = append(imports, "@jest/globals")
	}

	return imports
}

// generateModels generates model definitions from workflow
func (h *TypeScriptHandler) generateModels(workflow *models.Workflow) []ModelData {
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
func (h *TypeScriptHandler) generateFieldsFromSchema(schema interface{}) []FieldData {
	var fields []FieldData

	// This is a simplified implementation
	// In a real scenario, you would parse the JSON schema properly
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
			for fieldName, fieldSchema := range properties {
				field := FieldData{
					Name:        fieldName,
					Type:        h.mapSchemaTypeToTSType(fieldSchema),
					Description: h.getSchemaDescription(fieldSchema),
					Required:    h.isFieldRequired(fieldName, schemaMap),
				}
				fields = append(fields, field)
			}
		}
	}

	return fields
}

// mapSchemaTypeToTSType maps JSON schema types to TypeScript types
func (h *TypeScriptHandler) mapSchemaTypeToTSType(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if schemaType, ok := schemaMap["type"].(string); ok {
			switch schemaType {
			case "string":
				return "string"
			case "integer":
				return "number"
			case "number":
				return "number"
			case "boolean":
				return "boolean"
			case "array":
				return "any[]"
			case "object":
				return "Record<string, any>"
			}
		}
	}
	return "any"
}

// getSchemaDescription extracts description from schema
func (h *TypeScriptHandler) getSchemaDescription(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if desc, ok := schemaMap["description"].(string); ok {
			return desc
		}
	}
	return ""
}

// isFieldRequired checks if a field is required
func (h *TypeScriptHandler) isFieldRequired(fieldName string, schema map[string]interface{}) bool {
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
func (h *TypeScriptHandler) generateMethodDocs(methods []MethodData) string {
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

		docs.WriteString("```typescript\n")
		docs.WriteString(fmt.Sprintf("%s(", method.Name))
		for i, param := range method.Parameters {
			if i > 0 {
				docs.WriteString(", ")
			}
			docs.WriteString(fmt.Sprintf("%s: %s", param.Name, param.Type))
		}
		docs.WriteString(fmt.Sprintf("): Promise<%s>\n", method.ReturnType))
		docs.WriteString("```\n\n")
	}

	return docs.String()
}

// isValidNpmPackageName validates npm package name
func isValidNpmPackageName(name string) bool {
	if name == "" {
		return false
	}

	// Handle scoped packages
	if strings.HasPrefix(name, "@") {
		parts := strings.Split(name, "/")
		if len(parts) != 2 {
			return false
		}
		// Validate scope and package name separately
		return isValidNpmName(parts[0][1:]) && isValidNpmName(parts[1])
	}

	return isValidNpmName(name)
}

// isValidNpmName validates individual npm name component
func isValidNpmName(name string) bool {
	if name == "" || len(name) > 214 {
		return false
	}

	// Must be lowercase
	if strings.ToLower(name) != name {
		return false
	}

	// Cannot start with . or _
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return false
	}

	// Can only contain URL-safe characters
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.') {
			return false
		}
	}

	return true
}