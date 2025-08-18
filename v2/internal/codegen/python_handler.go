package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"magic-flow/v2/pkg/models"
)

// PythonHandler implements LanguageHandler for Python code generation
type PythonHandler struct {
	templateManager *TemplateManager
}

// NewPythonHandler creates a new Python language handler
func NewPythonHandler(templateManager *TemplateManager) *PythonHandler {
	return &PythonHandler{
		templateManager: templateManager,
	}
}

// Generate generates Python code for a workflow
func (h *PythonHandler) Generate(workflow *models.Workflow, request *GenerationRequest, templateData *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate __init__.py file
	initFile, err := h.generateInitFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate __init__.py file: %w", err)
	}
	files = append(files, initFile)

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

	// Generate exceptions file
	exceptionsFile, err := h.generateExceptionsFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate exceptions file: %w", err)
	}
	files = append(files, exceptionsFile)

	// Generate test file if requested
	if request.IncludeTests {
		testFile, err := h.generateTestFile(templateData)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test file: %w", err)
		}
		files = append(files, testFile)
	}

	// Generate setup.py file
	setupFile, err := h.generateSetupFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate setup.py file: %w", err)
	}
	files = append(files, setupFile)

	// Generate requirements.txt file
	requirementsFile, err := h.generateRequirementsFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate requirements.txt file: %w", err)
	}
	files = append(files, requirementsFile)

	// Generate pyproject.toml file
	pyprojectFile, err := h.generatePyprojectFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pyproject.toml file: %w", err)
	}
	files = append(files, pyprojectFile)

	// Generate README file
	readmeFile, err := h.generateReadmeFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate README file: %w", err)
	}
	files = append(files, readmeFile)

	return files, nil
}

// ValidateRequest validates Python-specific generation request
func (h *PythonHandler) ValidateRequest(request *GenerationRequest) error {
	if request.PackageName == "" {
		request.PackageName = h.GetDefaultPackageName()
	}

	// Validate package name format (Python package naming rules)
	if !isValidPythonPackageName(request.PackageName) {
		return fmt.Errorf("invalid Python package name: %s", request.PackageName)
	}

	return nil
}

// PrepareTemplateData prepares template data for Python code generation
func (h *PythonHandler) PrepareTemplateData(workflow *models.Workflow, request *GenerationRequest) (*TemplateData, error) {
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

// GetFileExtension returns the file extension for Python files
func (h *PythonHandler) GetFileExtension() string {
	return ".py"
}

// GetDefaultPackageName returns the default package name for Python
func (h *PythonHandler) GetDefaultPackageName() string {
	return "magicflow_client"
}

// generateInitFile generates the __init__.py file
func (h *PythonHandler) generateInitFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`"""
%s Python Client
Generated at: %s
"""

from .client import %s
from .models import *
from .types import *
from .exceptions import *

__version__ = "1.0.0"
__all__ = [
    "%s",
    # Add other exports here
]

# For convenience
Client = %s
`,
		data.Workflow.Name,
		data.GeneratedAt.Format("2006-01-02 15:04:05"),
		data.ClassName,
		data.ClassName,
		data.ClassName,
	)

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "__init__.py"),
		Content:  content,
		Language: "python",
		Type:     "init",
	}, nil
}

// generateClientFile generates the main client file
func (h *PythonHandler) generateClientFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("python", "client")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "client.py"),
		Content:  content,
		Language: "python",
		Type:     "client",
	}, nil
}

// generateModelsFile generates the models file
func (h *PythonHandler) generateModelsFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("python", "models")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "models.py"),
		Content:  content,
		Language: "python",
		Type:     "models",
	}, nil
}

// generateTypesFile generates the types file
func (h *PythonHandler) generateTypesFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("python", "types")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "types.py"),
		Content:  content,
		Language: "python",
		Type:     "types",
	}, nil
}

// generateExceptionsFile generates the exceptions file
func (h *PythonHandler) generateExceptionsFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`"""
Exceptions for %s Python Client
Generated at: %s
"""


class MagicFlowError(Exception):
    """Base exception for Magic Flow client."""
    pass


class APIError(MagicFlowError):
    """Exception raised for API errors."""
    
    def __init__(self, message: str, status_code: int = None, response_data: dict = None):
        super().__init__(message)
        self.status_code = status_code
        self.response_data = response_data or {}


class AuthenticationError(APIError):
    """Exception raised for authentication errors."""
    pass


class ValidationError(MagicFlowError):
    """Exception raised for validation errors."""
    pass


class ExecutionError(MagicFlowError):
    """Exception raised for workflow execution errors."""
    
    def __init__(self, message: str, execution_id: str = None, step_id: str = None):
        super().__init__(message)
        self.execution_id = execution_id
        self.step_id = step_id


class TimeoutError(MagicFlowError):
    """Exception raised for timeout errors."""
    pass


class NetworkError(MagicFlowError):
    """Exception raised for network errors."""
    pass
`,
		data.Workflow.Name,
		data.GeneratedAt.Format("2006-01-02 15:04:05"),
	)

	return GeneratedFile{
		Path:     filepath.Join(data.PackageName, "exceptions.py"),
		Content:  content,
		Language: "python",
		Type:     "exceptions",
	}, nil
}

// generateTestFile generates the test file
func (h *PythonHandler) generateTestFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("python", "test")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Path:     filepath.Join("tests", "test_client.py"),
		Content:  content,
		Language: "python",
		Type:     "test",
	}, nil
}

// generateSetupFile generates the setup.py file
func (h *PythonHandler) generateSetupFile(data *TemplateData) (GeneratedFile, error) {
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

	email := "contact@magicflow.dev"
	if data.Options != nil {
		if e, ok := data.Options["email"].(string); ok && e != "" {
			email = e
		}
	}

	content := fmt.Sprintf(`#!/usr/bin/env python3
"""
Setup script for %s Python Client
"""

from setuptools import setup, find_packages
import os

# Read README file
with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

# Read requirements
with open("requirements.txt", "r", encoding="utf-8") as fh:
    requirements = [line.strip() for line in fh if line.strip() and not line.startswith("#")]

setup(
    name="%s",
    version="%s",
    author="%s",
    author_email="%s",
    description="Python client for %s workflow",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/your-org/your-repo",
    project_urls={
        "Bug Tracker": "https://github.com/your-org/your-repo/issues",
        "Documentation": "https://docs.magicflow.dev",
        "Source Code": "https://github.com/your-org/your-repo",
    },
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Internet :: WWW/HTTP :: HTTP Servers",
        "Topic :: System :: Distributed Computing",
    ],
    python_requires=">=3.8",
    install_requires=requirements,
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "pytest-cov>=4.0.0",
            "pytest-asyncio>=0.21.0",
            "black>=23.0.0",
            "flake8>=6.0.0",
            "mypy>=1.0.0",
            "isort>=5.12.0",
        ],
        "docs": [
            "sphinx>=6.0.0",
            "sphinx-rtd-theme>=1.2.0",
        ],
    },
    keywords=["magic-flow", "workflow", "client", "api"],
    include_package_data=True,
    zip_safe=False,
)
`,
		data.Workflow.Name,
		data.PackageName,
		version,
		author,
		email,
		data.Workflow.Name,
	)

	return GeneratedFile{
		Path:     "setup.py",
		Content:  content,
		Language: "python",
		Type:     "config",
	}, nil
}

// generateRequirementsFile generates the requirements.txt file
func (h *PythonHandler) generateRequirementsFile(data *TemplateData) (GeneratedFile, error) {
	content := `# Core dependencies
requests>=2.28.0
httpx>=0.24.0
pydantic>=2.0.0
typing-extensions>=4.5.0

# Optional dependencies for async support
aiohttp>=3.8.0

# Development dependencies (install with pip install -e ".[dev]")
# pytest>=7.0.0
# pytest-cov>=4.0.0
# pytest-asyncio>=0.21.0
# black>=23.0.0
# flake8>=6.0.0
# mypy>=1.0.0
# isort>=5.12.0
`

	return GeneratedFile{
		Path:     "requirements.txt",
		Content:  content,
		Language: "text",
		Type:     "config",
	}, nil
}

// generatePyprojectFile generates the pyproject.toml file
func (h *PythonHandler) generatePyprojectFile(data *TemplateData) (GeneratedFile, error) {
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

	email := "contact@magicflow.dev"
	if data.Options != nil {
		if e, ok := data.Options["email"].(string); ok && e != "" {
			email = e
		}
	}

	content := fmt.Sprintf(`[build-system]
requires = ["setuptools>=61.0", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "%s"
version = "%s"
authors = [
    {name = "%s", email = "%s"},
]
description = "Python client for %s workflow"
readme = "README.md"
requires-python = ">=3.8"
classifiers = [
    "Development Status :: 4 - Beta",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Operating System :: OS Independent",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
    "Programming Language :: Python :: 3.12",
]
keywords = ["magic-flow", "workflow", "client", "api"]
dependencies = [
    "requests>=2.28.0",
    "httpx>=0.24.0",
    "pydantic>=2.0.0",
    "typing-extensions>=4.5.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "pytest-asyncio>=0.21.0",
    "black>=23.0.0",
    "flake8>=6.0.0",
    "mypy>=1.0.0",
    "isort>=5.12.0",
]
docs = [
    "sphinx>=6.0.0",
    "sphinx-rtd-theme>=1.2.0",
]
async = [
    "aiohttp>=3.8.0",
]

[project.urls]
"Homepage" = "https://github.com/your-org/your-repo"
"Bug Tracker" = "https://github.com/your-org/your-repo/issues"
"Documentation" = "https://docs.magicflow.dev"
"Source Code" = "https://github.com/your-org/your-repo"

[tool.setuptools.packages.find]
where = ["."]
include = ["%s*"]

[tool.black]
line-length = 88
target-version = ['py38']
include = '\.pyi?$'
extend-exclude = '''
(
  /(
      \.eggs
    | \.git
    | \.hg
    | \.mypy_cache
    | \.tox
    | \.venv
    | _build
    | buck-out
    | build
    | dist
  )/
)
'''

[tool.isort]
profile = "black"
line_length = 88
multi_line_output = 3
include_trailing_comma = true
force_grid_wrap = 0
use_parentheses = true
ensure_newline_before_comments = true

[tool.mypy]
python_version = "3.8"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true
disallow_incomplete_defs = true
check_untyped_defs = true
disallow_untyped_decorators = true
no_implicit_optional = true
warn_redundant_casts = true
warn_unused_ignores = true
warn_no_return = true
warn_unreachable = true
strict_equality = true

[tool.pytest.ini_options]
minversion = "7.0"
addopts = "-ra -q --strict-markers --strict-config"
testpaths = [
    "tests",
]
python_files = [
    "test_*.py",
    "*_test.py",
]
python_classes = [
    "Test*",
]
python_functions = [
    "test_*",
]
markers = [
    "slow: marks tests as slow (deselect with '-m "not slow"')",
    "integration: marks tests as integration tests",
]

[tool.coverage.run]
source = ["%s"]
omit = [
    "*/tests/*",
    "*/test_*.py",
    "*/*_test.py",
]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "if self.debug:",
    "if settings.DEBUG",
    "raise AssertionError",
    "raise NotImplementedError",
    "if 0:",
    "if __name__ == .__main__.:",
    "class .*\\bProtocol\\):",
    "@(abc\\.)?abstractmethod",
]
`,
		data.PackageName,
		version,
		author,
		email,
		data.Workflow.Name,
		data.PackageName,
		data.PackageName,
	)

	return GeneratedFile{
		Path:     "pyproject.toml",
		Content:  content,
		Language: "toml",
		Type:     "config",
	}, nil
}

// generateReadmeFile generates the README file
func (h *PythonHandler) generateReadmeFile(data *TemplateData) (GeneratedFile, error) {
	content := fmt.Sprintf(`# %s Python Client

Generated Python client library for the %s workflow.

## Installation

` + "```bash" + `
pip install %s
` + "```" + `

## Quick Start

` + "```python" + `
from %s import %s

# Initialize the client
client = %s(
    base_url="http://localhost:8080",
    api_key="your-api-key"
)

# Execute workflow
input_data = {
    "key": "value"
}

try:
    result = client.execute_workflow(input_data)
    print(f"Execution ID: {result.id}")
    print(f"Status: {result.status}")
    
    # Check execution status
    status = client.get_execution_status(result.id)
    print(f"Current status: {status.status}")
    print(f"Progress: {status.progress}%%")
    
except Exception as e:
    print(f"Error executing workflow: {e}")
` + "```" + `

## Async Usage

` + "```python" + `
import asyncio
from %s import %s

async def main():
    async with %s(
        base_url="http://localhost:8080",
        api_key="your-api-key"
    ) as client:
        input_data = {
            "key": "value"
        }
        
        try:
            result = await client.execute_workflow(input_data)
            print(f"Execution ID: {result.id}")
            
            # Wait for completion
            final_result = await client.wait_for_completion(result.id)
            print(f"Final status: {final_result.status}")
            
        except Exception as e:
            print(f"Error: {e}")

# Run async example
asyncio.run(main())
` + "```" + `

## API Reference

### Client Methods

#### execute_workflow

Executes the %s workflow with the provided input.

` + "```python" + `
def execute_workflow(self, input_data: Dict[str, Any]) -> ExecutionResult:
    """Execute workflow with input data."""
    pass

async def execute_workflow(self, input_data: Dict[str, Any]) -> ExecutionResult:
    """Execute workflow with input data (async version)."""
    pass
` + "```" + `

#### get_execution_status

Retrieves the status of a workflow execution.

` + "```python" + `
def get_execution_status(self, execution_id: str) -> ExecutionStatus:
    """Get execution status by ID."""
    pass

async def get_execution_status(self, execution_id: str) -> ExecutionStatus:
    """Get execution status by ID (async version)."""
    pass
` + "```" + `

#### cancel_execution

Cancels a running workflow execution.

` + "```python" + `
def cancel_execution(self, execution_id: str) -> None:
    """Cancel a running execution."""
    pass

async def cancel_execution(self, execution_id: str) -> None:
    """Cancel a running execution (async version)."""
    pass
` + "```" + `

#### wait_for_completion

Waits for a workflow execution to complete.

` + "```python" + `
def wait_for_completion(
    self, 
    execution_id: str, 
    timeout: int = 300,
    poll_interval: int = 5
) -> ExecutionResult:
    """Wait for execution to complete."""
    pass

async def wait_for_completion(
    self, 
    execution_id: str, 
    timeout: int = 300,
    poll_interval: int = 5
) -> ExecutionResult:
    """Wait for execution to complete (async version)."""
    pass
` + "```" + `

%s

## Models

### ExecutionResult

Represents the result of a workflow execution.

` + "```python" + `
class ExecutionResult:
    id: str
    workflow_id: str
    status: ExecutionStatus
    input: Dict[str, Any]
    output: Optional[Dict[str, Any]]
    error: Optional[str]
    started_at: datetime
    completed_at: Optional[datetime]
    duration: Optional[int]
` + "```" + `

### ExecutionStatus

Represents the status of a workflow execution.

` + "```python" + `
class ExecutionStatus:
    id: str
    status: Literal['pending', 'running', 'completed', 'failed', 'cancelled']
    progress: int
    current_step: Optional[str]
    steps: List[StepStatus]
    started_at: datetime
    updated_at: datetime
` + "```" + `

### StepStatus

Represents the status of a workflow step.

` + "```python" + `
class StepStatus:
    id: str
    name: str
    status: Literal['pending', 'running', 'completed', 'failed', 'skipped']
    input: Optional[Dict[str, Any]]
    output: Optional[Dict[str, Any]]
    error: Optional[str]
    started_at: Optional[datetime]
    completed_at: Optional[datetime]
    duration: Optional[int]
` + "```" + `

## Constants

` + "```python" + `
# Workflow information
WORKFLOW_ID = "%s"
WORKFLOW_NAME = "%s"

# Execution statuses
class ExecutionStatus:
    PENDING = "pending"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"

# Step IDs
class StepIds:
    # Add step constants here based on workflow definition
    pass
` + "```" + `

## Error Handling

The client provides several exception types for different error scenarios:

` + "```python" + `
from %s.exceptions import (
    MagicFlowError,
    APIError,
    AuthenticationError,
    ValidationError,
    ExecutionError,
    TimeoutError,
    NetworkError
)

try:
    result = client.execute_workflow(input_data)
except AuthenticationError:
    print("Invalid API key")
except ValidationError as e:
    print(f"Invalid input: {e}")
except ExecutionError as e:
    print(f"Execution failed: {e}")
except APIError as e:
    print(f"API error: {e} (status: {e.status_code})")
except NetworkError:
    print("Network connection failed")
except MagicFlowError as e:
    print(f"General error: {e}")
` + "```" + `

## Configuration

### Environment Variables

` + "```bash" + `
# Set default base URL
export MAGICFLOW_BASE_URL="http://localhost:8080"

# Set default API key
export MAGICFLOW_API_KEY="your-api-key"

# Set default timeout (seconds)
export MAGICFLOW_TIMEOUT="30"

# Enable debug logging
export MAGICFLOW_DEBUG="true"
` + "```" + `

### Client Configuration

` + "```python" + `
client = %s(
    base_url="http://localhost:8080",
    api_key="your-api-key",
    timeout=30,
    retry_attempts=3,
    retry_delay=1.0,
    debug=False
)
` + "```" + `

## Development

### Setup Development Environment

` + "```bash" + `
# Clone the repository
git clone https://github.com/your-org/your-repo.git
cd your-repo

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install development dependencies
pip install -e ".[dev]"
` + "```" + `

### Running Tests

` + "```bash" + `
# Run all tests
pytest

# Run with coverage
pytest --cov=%s --cov-report=html

# Run specific test file
pytest tests/test_client.py

# Run with verbose output
pytest -v
` + "```" + `

### Code Quality

` + "```bash" + `
# Format code
black %s

# Sort imports
isort %s

# Lint code
flake8 %s

# Type checking
mypy %s
` + "```" + `

## License

Generated code - see original workflow license.
`,
		data.Workflow.Name,
		data.Workflow.Name,
		data.PackageName,
		data.PackageName,
		data.ClassName,
		data.ClassName,
		data.PackageName,
		data.ClassName,
		data.ClassName,
		data.Workflow.Name,
		h.generateMethodDocs(data.Methods),
		data.Workflow.ID.String(),
		data.Workflow.Name,
		data.PackageName,
		data.ClassName,
		data.PackageName,
		data.PackageName,
		data.PackageName,
		data.PackageName,
		data.PackageName,
	)

	return GeneratedFile{
		Path:     "README.md",
		Content:  content,
		Language: "markdown",
		Type:     "documentation",
	}, nil
}

// generateImports generates the list of imports needed
func (h *PythonHandler) generateImports(workflow *models.Workflow, request *GenerationRequest) []string {
	imports := []string{
		"import json",
		"import time",
		"from datetime import datetime",
		"from typing import Dict, Any, Optional, List, Union",
		"from uuid import UUID",
		"import requests",
		"from pydantic import BaseModel, Field",
	}

	// Add test imports if tests are included
	if request.IncludeTests {
		imports = append(imports, "import pytest", "from unittest.mock import Mock, patch")
	}

	return imports
}

// generateModels generates model definitions from workflow
func (h *PythonHandler) generateModels(workflow *models.Workflow) []ModelData {
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
func (h *PythonHandler) generateFieldsFromSchema(schema interface{}) []FieldData {
	var fields []FieldData

	// This is a simplified implementation
	// In a real scenario, you would parse the JSON schema properly
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
			for fieldName, fieldSchema := range properties {
				field := FieldData{
					Name:        ToSnakeCase(fieldName),
					Type:        h.mapSchemaTypeToPythonType(fieldSchema),
					Description: h.getSchemaDescription(fieldSchema),
					Required:    h.isFieldRequired(fieldName, schemaMap),
				}
				fields = append(fields, field)
			}
		}
	}

	return fields
}

// mapSchemaTypeToPythonType maps JSON schema types to Python types
func (h *PythonHandler) mapSchemaTypeToPythonType(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if schemaType, ok := schemaMap["type"].(string); ok {
			switch schemaType {
			case "string":
				return "str"
			case "integer":
				return "int"
			case "number":
				return "float"
			case "boolean":
				return "bool"
			case "array":
				return "List[Any]"
			case "object":
				return "Dict[str, Any]"
			}
		}
	}
	return "Any"
}

// getSchemaDescription extracts description from schema
func (h *PythonHandler) getSchemaDescription(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if desc, ok := schemaMap["description"].(string); ok {
			return desc
		}
	}
	return ""
}

// isFieldRequired checks if a field is required
func (h *PythonHandler) isFieldRequired(fieldName string, schema map[string]interface{}) bool {
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
func (h *PythonHandler) generateMethodDocs(methods []MethodData) string {
	if len(methods) == 0 {
		return ""
	}

	var docs strings.Builder
	docs.WriteString("\n### Step Methods\n\n")

	for _, method := range methods {
		docs.WriteString(fmt.Sprintf("#### %s\n\n", ToSnakeCase(method.Name)))
		if method.Description != "" {
			docs.WriteString(fmt.Sprintf("%s\n\n", method.Description))
		}

		docs.WriteString("```python\n")
		docs.WriteString(fmt.Sprintf("def %s(self", ToSnakeCase(method.Name)))
		for _, param := range method.Parameters {
			docs.WriteString(fmt.Sprintf(", %s: %s", ToSnakeCase(param.Name), param.Type))
		}
		docs.WriteString(fmt.Sprintf(") -> %s:\n", method.ReturnType))
		docs.WriteString(fmt.Sprintf("    \"\"\"Execute %s step.\"\"\"\n", method.Name))
		docs.WriteString("    pass\n\n")

		// Async version
		docs.WriteString(fmt.Sprintf("async def %s(self", ToSnakeCase(method.Name)))
		for _, param := range method.Parameters {
			docs.WriteString(fmt.Sprintf(", %s: %s", ToSnakeCase(param.Name), param.Type))
		}
		docs.WriteString(fmt.Sprintf(") -> %s:\n", method.ReturnType))
		docs.WriteString(fmt.Sprintf("    \"\"\"Execute %s step (async version).\"\"\"\n", method.Name))
		docs.WriteString("    pass\n")
		docs.WriteString("```\n\n")
	}

	return docs.String()
}

// isValidPythonPackageName validates Python package name
func isValidPythonPackageName(name string) bool {
	if name == "" {
		return false
	}

	// Python package names should be lowercase and contain only letters, numbers, and underscores
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