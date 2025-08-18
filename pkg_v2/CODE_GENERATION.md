# Magic Flow v2 - Code Generation Documentation

## Overview

Magic Flow v2 provides powerful code generation capabilities that automatically create client libraries, service implementations, and integration code from YAML workflow definitions. This document covers the complete code generation process, supported languages, templates, customization options, and integration patterns.

## Code Generation Architecture

### Generation Pipeline

```
YAML Workflow Definition
         ↓
    Schema Parser
         ↓
   Template Engine
         ↓
  Language Generator
         ↓
   Code Formatter
         ↓
  Generated Artifacts
```

### Core Components

1. **Schema Parser**: Analyzes YAML workflow definitions and extracts metadata
2. **Template Engine**: Processes templates with workflow data
3. **Language Generator**: Generates language-specific code
4. **Code Formatter**: Applies language-specific formatting and linting
5. **Artifact Manager**: Packages and delivers generated code

## Supported Languages

### Go

**Features:**
- Native workflow execution engine
- Type-safe client libraries
- Context-aware error handling
- Built-in observability
- Graceful shutdown support

**Generated Artifacts:**
- Workflow executor service
- Client SDK
- Data models (structs)
- API interfaces
- Configuration types
- Test scaffolding

**Example Generation:**
```bash
# Generate Go code
magicflow generate go \
  --workflow order_processing.yaml \
  --output ./generated/go \
  --package github.com/company/order-service
```

**Generated Structure:**
```
generated/go/
├── cmd/
│   └── executor/
│       └── main.go                 # Service entry point
├── pkg/
│   ├── client/
│   │   ├── client.go              # Workflow client
│   │   └── types.go               # Request/response types
│   ├── executor/
│   │   ├── executor.go            # Workflow executor
│   │   ├── steps.go               # Step implementations
│   │   └── handlers.go            # HTTP handlers
│   ├── models/
│   │   ├── workflow.go            # Workflow models
│   │   └── context.go             # Execution context
│   └── config/
│       └── config.go              # Configuration
├── api/
│   └── openapi.yaml               # OpenAPI specification
├── docker/
│   └── Dockerfile                 # Container image
├── k8s/
│   ├── deployment.yaml            # Kubernetes deployment
│   └── service.yaml               # Kubernetes service
├── go.mod                         # Go module
├── go.sum                         # Dependencies
└── README.md                      # Usage documentation
```

### TypeScript/Node.js

**Features:**
- Promise-based async execution
- Type definitions with TypeScript
- Express.js integration
- Comprehensive error handling
- Built-in logging and metrics

**Generated Artifacts:**
- Express.js service
- TypeScript client SDK
- Type definitions
- API routes
- Middleware
- Test suites

**Example Generation:**
```bash
# Generate TypeScript code
magicflow generate typescript \
  --workflow order_processing.yaml \
  --output ./generated/typescript \
  --package @company/order-service
```

**Generated Structure:**
```
generated/typescript/
├── src/
│   ├── server.ts                  # Express server
│   ├── executor/
│   │   ├── workflow-executor.ts   # Workflow executor
│   │   └── step-handlers.ts       # Step implementations
│   ├── client/
│   │   ├── workflow-client.ts     # Client SDK
│   │   └── types.ts               # TypeScript types
│   ├── models/
│   │   ├── workflow.ts            # Workflow models
│   │   └── context.ts             # Execution context
│   ├── routes/
│   │   └── workflow.ts            # API routes
│   └── middleware/
│       ├── auth.ts                # Authentication
│       └── logging.ts             # Request logging
├── tests/
│   ├── unit/                      # Unit tests
│   └── integration/               # Integration tests
├── docker/
│   └── Dockerfile                 # Container image
├── k8s/
│   ├── deployment.yaml            # Kubernetes deployment
│   └── service.yaml               # Kubernetes service
├── package.json                   # NPM package
├── tsconfig.json                  # TypeScript config
├── jest.config.js                 # Test configuration
└── README.md                      # Usage documentation
```

### Python

**Features:**
- FastAPI-based service
- Pydantic data validation
- Async/await support
- Type hints
- Comprehensive testing

**Generated Artifacts:**
- FastAPI service
- Python client library
- Pydantic models
- API endpoints
- Background tasks
- Test framework

**Example Generation:**
```bash
# Generate Python code
magicflow generate python \
  --workflow order_processing.yaml \
  --output ./generated/python \
  --package order_service
```

**Generated Structure:**
```
generated/python/
├── order_service/
│   ├── __init__.py
│   ├── main.py                    # FastAPI application
│   ├── executor/
│   │   ├── __init__.py
│   │   ├── workflow_executor.py   # Workflow executor
│   │   └── step_handlers.py       # Step implementations
│   ├── client/
│   │   ├── __init__.py
│   │   ├── workflow_client.py     # Client SDK
│   │   └── types.py               # Type definitions
│   ├── models/
│   │   ├── __init__.py
│   │   ├── workflow.py            # Pydantic models
│   │   └── context.py             # Execution context
│   ├── api/
│   │   ├── __init__.py
│   │   └── workflow.py            # API endpoints
│   └── config/
│       ├── __init__.py
│       └── settings.py            # Configuration
├── tests/
│   ├── __init__.py
│   ├── test_executor.py           # Executor tests
│   └── test_client.py             # Client tests
├── docker/
│   └── Dockerfile                 # Container image
├── k8s/
│   ├── deployment.yaml            # Kubernetes deployment
│   └── service.yaml               # Kubernetes service
├── requirements.txt               # Dependencies
├── pyproject.toml                 # Project configuration
├── pytest.ini                    # Test configuration
└── README.md                      # Usage documentation
```

### Java

**Features:**
- Spring Boot integration
- Maven/Gradle support
- Jackson JSON processing
- Comprehensive validation
- Enterprise-ready patterns

**Generated Artifacts:**
- Spring Boot application
- Java client library
- POJOs with validation
- REST controllers
- Service layer
- Test classes

**Example Generation:**
```bash
# Generate Java code
magicflow generate java \
  --workflow order_processing.yaml \
  --output ./generated/java \
  --package com.company.orderservice
```

### C#/.NET

**Features:**
- ASP.NET Core integration
- Entity Framework support
- Dependency injection
- Comprehensive validation
- NuGet package generation

**Generated Artifacts:**
- ASP.NET Core application
- .NET client library
- Data models
- Controllers
- Services
- Unit tests

## Template System

### Template Structure

Magic Flow uses a hierarchical template system with language-specific templates and shared components.

```
templates/
├── shared/
│   ├── docker/
│   │   └── Dockerfile.tmpl
│   ├── k8s/
│   │   ├── deployment.yaml.tmpl
│   │   └── service.yaml.tmpl
│   └── docs/
│       └── README.md.tmpl
├── go/
│   ├── cmd/
│   │   └── main.go.tmpl
│   ├── pkg/
│   │   ├── client.go.tmpl
│   │   ├── executor.go.tmpl
│   │   └── models.go.tmpl
│   └── go.mod.tmpl
├── typescript/
│   ├── src/
│   │   ├── server.ts.tmpl
│   │   ├── client.ts.tmpl
│   │   └── types.ts.tmpl
│   └── package.json.tmpl
└── python/
    ├── main.py.tmpl
    ├── client.py.tmpl
    └── requirements.txt.tmpl
```

### Template Variables

Templates have access to rich workflow metadata:

```yaml
# Available template variables
workflow:
  name: "order_processing"
  version: "1.0.0"
  description: "Process customer orders"
  metadata:
    owner: "ecommerce-team"
    tags: ["ecommerce", "orders"]

steps:
  - name: "validate_order"
    type: "service_call"
    config:
      service: "validation-service"
      endpoint: "/api/v1/validate"
    input_schema: { ... }
    output_schema: { ... }

input_schema:
  type: "object"
  properties: { ... }

output_schema:
  type: "object"
  properties: { ... }

generation:
  language: "go"
  package: "github.com/company/order-service"
  output_dir: "./generated/go"
  options:
    include_tests: true
    include_docker: true
    include_k8s: true
```

### Template Syntax

Magic Flow uses Go's `text/template` syntax with additional helper functions:

```go
// Basic variable substitution
package {{.generation.package}}

// Conditional generation
{{if .generation.options.include_tests}}
// Test code here
{{end}}

// Loop over steps
{{range .steps}}
func handle{{.name | title}}(ctx context.Context, input {{.input_type}}) ({{.output_type}}, error) {
    // Step implementation
}
{{end}}

// Helper functions
{{.workflow.name | camelCase}}     // orderProcessing
{{.workflow.name | pascalCase}}    // OrderProcessing
{{.workflow.name | kebabCase}}     // order-processing
{{.workflow.name | snakeCase}}     // order_processing
```

### Custom Templates

Users can provide custom templates for specific use cases:

```bash
# Use custom template directory
magicflow generate go \
  --workflow order_processing.yaml \
  --templates ./custom-templates \
  --output ./generated/go
```

**Custom Template Structure:**
```
custom-templates/
├── go/
│   ├── cmd/
│   │   └── main.go.tmpl           # Custom main template
│   └── pkg/
│       └── custom-handler.go.tmpl # Additional custom file
└── shared/
    └── custom-config.yaml.tmpl    # Custom configuration
```

## Code Generation Options

### CLI Options

```bash
magicflow generate [language] [options]

Options:
  --workflow, -w        Workflow YAML file (required)
  --output, -o          Output directory (required)
  --package, -p         Package name/namespace
  --templates, -t       Custom template directory
  --config, -c          Generation configuration file
  --include-tests       Generate test files (default: true)
  --include-docker      Generate Dockerfile (default: true)
  --include-k8s         Generate Kubernetes manifests (default: true)
  --include-docs        Generate documentation (default: true)
  --format              Format generated code (default: true)
  --validate            Validate generated code (default: true)
  --overwrite           Overwrite existing files (default: false)
  --dry-run             Show what would be generated without creating files
  --verbose, -v         Verbose output
```

### Configuration File

```yaml
# generation-config.yaml
generation:
  # Global settings
  author: "Platform Team"
  license: "MIT"
  copyright: "2024 Company Inc."
  
  # Language-specific settings
  go:
    module: "github.com/company/workflows"
    go_version: "1.21"
    build_tags: ["integration"]
    
  typescript:
    registry: "@company"
    node_version: "18"
    
  python:
    python_version: "3.11"
    package_index: "https://pypi.company.com"
    
  java:
    group_id: "com.company"
    artifact_id: "workflow-service"
    java_version: "17"
    
  # Feature flags
  features:
    include_metrics: true
    include_tracing: true
    include_health_checks: true
    include_graceful_shutdown: true
    include_rate_limiting: true
    
  # Docker settings
  docker:
    base_image: "alpine:3.18"
    registry: "registry.company.com"
    
  # Kubernetes settings
  kubernetes:
    namespace: "workflows"
    ingress_class: "nginx"
    resource_limits:
      cpu: "500m"
      memory: "512Mi"
```

## Generated Code Examples

### Go Client Example

```go
// Generated Go client
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// OrderProcessingClient represents the workflow client
type OrderProcessingClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

// NewOrderProcessingClient creates a new workflow client
func NewOrderProcessingClient(baseURL, apiKey string) *OrderProcessingClient {
    return &OrderProcessingClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        apiKey: apiKey,
    }
}

// OrderProcessingInput represents the workflow input
type OrderProcessingInput struct {
    CustomerID      string                 `json:"customer_id" validate:"required,uuid"`
    Items           []OrderItem            `json:"items" validate:"required,min=1,max=50"`
    PaymentMethod   PaymentMethod          `json:"payment_method" validate:"required"`
    ShippingAddress Address                `json:"shipping_address" validate:"required"`
    Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// OrderProcessingOutput represents the workflow output
type OrderProcessingOutput struct {
    OrderID           string    `json:"order_id"`
    Status            string    `json:"status"`
    TotalAmount       float64   `json:"total_amount"`
    TransactionID     string    `json:"transaction_id"`
    TrackingNumber    string    `json:"tracking_number,omitempty"`
    EstimatedDelivery time.Time `json:"estimated_delivery,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
}

// ExecuteWorkflow executes the order processing workflow
func (c *OrderProcessingClient) ExecuteWorkflow(ctx context.Context, input OrderProcessingInput) (*OrderProcessingOutput, error) {
    // Validate input
    if err := c.validateInput(input); err != nil {
        return nil, fmt.Errorf("invalid input: %w", err)
    }
    
    // Execute workflow
    resp, err := c.executeHTTPRequest(ctx, "POST", "/api/v1/workflows/order-processing/execute", input)
    if err != nil {
        return nil, fmt.Errorf("failed to execute workflow: %w", err)
    }
    
    var output OrderProcessingOutput
    if err := json.Unmarshal(resp, &output); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    return &output, nil
}

// ExecuteWorkflowAsync executes the workflow asynchronously
func (c *OrderProcessingClient) ExecuteWorkflowAsync(ctx context.Context, input OrderProcessingInput) (string, error) {
    // Implementation for async execution
    // Returns execution ID for tracking
}

// GetExecutionStatus gets the status of an async execution
func (c *OrderProcessingClient) GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error) {
    // Implementation for status checking
}

// GetExecutionResult gets the result of a completed execution
func (c *OrderProcessingClient) GetExecutionResult(ctx context.Context, executionID string) (*OrderProcessingOutput, error) {
    // Implementation for result retrieval
}
```

### TypeScript Client Example

```typescript
// Generated TypeScript client
import axios, { AxiosInstance, AxiosResponse } from 'axios';

export interface OrderProcessingInput {
  customer_id: string;
  items: OrderItem[];
  payment_method: PaymentMethod;
  shipping_address: Address;
  metadata?: Record<string, any>;
}

export interface OrderProcessingOutput {
  order_id: string;
  status: string;
  total_amount: number;
  transaction_id: string;
  tracking_number?: string;
  estimated_delivery?: Date;
  created_at: Date;
}

export interface ExecutionStatus {
  execution_id: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progress: number;
  current_step?: string;
  error?: string;
  started_at: Date;
  completed_at?: Date;
}

export class OrderProcessingClient {
  private client: AxiosInstance;

  constructor(baseURL: string, apiKey: string) {
    this.client = axios.create({
      baseURL,
      headers: {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json',
      },
      timeout: 30000,
    });
  }

  /**
   * Execute the order processing workflow synchronously
   */
  async executeWorkflow(input: OrderProcessingInput): Promise<OrderProcessingOutput> {
    try {
      const response: AxiosResponse<OrderProcessingOutput> = await this.client.post(
        '/api/v1/workflows/order-processing/execute',
        input
      );
      return response.data;
    } catch (error) {
      throw new Error(`Failed to execute workflow: ${error.message}`);
    }
  }

  /**
   * Execute the workflow asynchronously
   */
  async executeWorkflowAsync(input: OrderProcessingInput): Promise<string> {
    try {
      const response: AxiosResponse<{ execution_id: string }> = await this.client.post(
        '/api/v1/workflows/order-processing/execute-async',
        input
      );
      return response.data.execution_id;
    } catch (error) {
      throw new Error(`Failed to execute workflow async: ${error.message}`);
    }
  }

  /**
   * Get execution status
   */
  async getExecutionStatus(executionId: string): Promise<ExecutionStatus> {
    try {
      const response: AxiosResponse<ExecutionStatus> = await this.client.get(
        `/api/v1/executions/${executionId}/status`
      );
      return response.data;
    } catch (error) {
      throw new Error(`Failed to get execution status: ${error.message}`);
    }
  }

  /**
   * Get execution result
   */
  async getExecutionResult(executionId: string): Promise<OrderProcessingOutput> {
    try {
      const response: AxiosResponse<OrderProcessingOutput> = await this.client.get(
        `/api/v1/executions/${executionId}/result`
      );
      return response.data;
    } catch (error) {
      throw new Error(`Failed to get execution result: ${error.message}`);
    }
  }

  /**
   * Stream execution events
   */
  async *streamExecutionEvents(executionId: string): AsyncGenerator<ExecutionEvent> {
    // Implementation for server-sent events streaming
  }
}
```

### Python Client Example

```python
# Generated Python client
from typing import Dict, Any, Optional, AsyncGenerator
from datetime import datetime
from dataclasses import dataclass
from enum import Enum
import httpx
import asyncio

class ExecutionStatus(str, Enum):
    PENDING = "pending"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class OrderProcessingInput:
    customer_id: str
    items: list[Dict[str, Any]]
    payment_method: Dict[str, Any]
    shipping_address: Dict[str, Any]
    metadata: Optional[Dict[str, Any]] = None

@dataclass
class OrderProcessingOutput:
    order_id: str
    status: str
    total_amount: float
    transaction_id: str
    tracking_number: Optional[str] = None
    estimated_delivery: Optional[datetime] = None
    created_at: datetime

@dataclass
class ExecutionStatusResponse:
    execution_id: str
    status: ExecutionStatus
    progress: float
    current_step: Optional[str] = None
    error: Optional[str] = None
    started_at: datetime
    completed_at: Optional[datetime] = None

class OrderProcessingClient:
    """Client for the Order Processing workflow."""
    
    def __init__(self, base_url: str, api_key: str):
        self.base_url = base_url.rstrip('/')
        self.client = httpx.AsyncClient(
            headers={
                "Authorization": f"Bearer {api_key}",
                "Content-Type": "application/json",
            },
            timeout=30.0,
        )
    
    async def execute_workflow(self, input_data: OrderProcessingInput) -> OrderProcessingOutput:
        """Execute the order processing workflow synchronously."""
        try:
            response = await self.client.post(
                f"{self.base_url}/api/v1/workflows/order-processing/execute",
                json=input_data.__dict__,
            )
            response.raise_for_status()
            data = response.json()
            return OrderProcessingOutput(**data)
        except httpx.HTTPError as e:
            raise Exception(f"Failed to execute workflow: {e}")
    
    async def execute_workflow_async(self, input_data: OrderProcessingInput) -> str:
        """Execute the workflow asynchronously."""
        try:
            response = await self.client.post(
                f"{self.base_url}/api/v1/workflows/order-processing/execute-async",
                json=input_data.__dict__,
            )
            response.raise_for_status()
            data = response.json()
            return data["execution_id"]
        except httpx.HTTPError as e:
            raise Exception(f"Failed to execute workflow async: {e}")
    
    async def get_execution_status(self, execution_id: str) -> ExecutionStatusResponse:
        """Get execution status."""
        try:
            response = await self.client.get(
                f"{self.base_url}/api/v1/executions/{execution_id}/status"
            )
            response.raise_for_status()
            data = response.json()
            return ExecutionStatusResponse(**data)
        except httpx.HTTPError as e:
            raise Exception(f"Failed to get execution status: {e}")
    
    async def get_execution_result(self, execution_id: str) -> OrderProcessingOutput:
        """Get execution result."""
        try:
            response = await self.client.get(
                f"{self.base_url}/api/v1/executions/{execution_id}/result"
            )
            response.raise_for_status()
            data = response.json()
            return OrderProcessingOutput(**data)
        except httpx.HTTPError as e:
            raise Exception(f"Failed to get execution result: {e}")
    
    async def stream_execution_events(self, execution_id: str) -> AsyncGenerator[Dict[str, Any], None]:
        """Stream execution events."""
        async with self.client.stream(
            "GET",
            f"{self.base_url}/api/v1/executions/{execution_id}/events",
            headers={"Accept": "text/event-stream"},
        ) as response:
            async for line in response.aiter_lines():
                if line.startswith("data: "):
                    yield json.loads(line[6:])
    
    async def close(self):
        """Close the HTTP client."""
        await self.client.aclose()
```

## Integration Patterns

### Microservices Integration

```yaml
# microservices-integration.yaml
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: "microservices-order-flow"
  version: "1.0.0"

spec:
  steps:
    - name: "user_service_validation"
      type: "service_call"
      config:
        service: "user-service"
        endpoint: "/api/v1/users/${input.user_id}/validate"
        discovery: "consul"  # Service discovery
        circuit_breaker:
          enabled: true
          failure_threshold: 5
          timeout: "30s"
    
    - name: "inventory_service_check"
      type: "service_call"
      config:
        service: "inventory-service"
        endpoint: "/api/v1/inventory/check"
        load_balancer: "round_robin"
        retry_policy:
          max_attempts: 3
          backoff: "exponential"
    
    - name: "payment_service_charge"
      type: "service_call"
      config:
        service: "payment-service"
        endpoint: "/api/v1/payments/charge"
        timeout: "45s"
        idempotency_key: "${input.order_id}-${now('unix')}"
```

### Event-Driven Integration

```yaml
# event-driven-integration.yaml
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: "event-driven-order-flow"
  version: "1.0.0"

spec:
  triggers:
    - type: "event"
      source: "order-events"
      event_type: "order.created"
      filter: "$.data.amount > 100"
  
  steps:
    - name: "publish_order_validated"
      type: "message_queue"
      config:
        exchange: "order-events"
        routing_key: "order.validated"
        message_format: "cloudevents"
      input:
        event:
          type: "order.validated"
          source: "order-processing-workflow"
          data: "$.context.order"
    
    - name: "publish_payment_processed"
      type: "message_queue"
      config:
        exchange: "payment-events"
        routing_key: "payment.processed"
      input:
        event:
          type: "payment.processed"
          source: "order-processing-workflow"
          data:
            order_id: "$.context.order.id"
            payment_id: "$.steps.process_payment.output.payment_id"
            amount: "$.steps.process_payment.output.amount"
```

### Database Integration

```yaml
# database-integration.yaml
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: "database-order-flow"
  version: "1.0.0"

spec:
  steps:
    - name: "create_order_record"
      type: "database"
      config:
        connection: "primary_postgres"
        operation: "insert"
        table: "orders"
        transaction: true
      input:
        data:
          customer_id: "$.input.customer_id"
          total_amount: "$.steps.calculate_total.output.amount"
          status: "pending"
          created_at: "${now()}"
    
    - name: "update_inventory"
      type: "database"
      config:
        connection: "inventory_mysql"
        operation: "update"
        table: "inventory"
        transaction: true
      input:
        data:
          quantity: "$.input.quantity"
        where:
          product_id: "$.input.product_id"
    
    - name: "log_audit_trail"
      type: "database"
      config:
        connection: "audit_mongodb"
        operation: "insert"
        collection: "audit_logs"
      input:
        document:
          workflow_id: "$.execution.workflow_id"
          execution_id: "$.execution.id"
          action: "order_created"
          user_id: "$.input.customer_id"
          timestamp: "${now()}"
          metadata: "$.context"
```

## Advanced Features

### Custom Code Injection

```yaml
# Custom code injection points
generation:
  custom_code:
    # Inject custom imports
    imports:
      go: |
        import (
            "github.com/company/custom-lib"
            "github.com/company/metrics"
        )
      
      typescript: |
        import { CustomValidator } from '@company/validation';
        import { MetricsCollector } from '@company/metrics';
    
    # Inject custom middleware
    middleware:
      go: |
        func CustomAuthMiddleware(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Custom authentication logic
                next.ServeHTTP(w, r)
            })
        }
      
      typescript: |
        export const customAuthMiddleware = (req: Request, res: Response, next: NextFunction) => {
            // Custom authentication logic
            next();
        };
    
    # Inject custom step handlers
    step_handlers:
      validate_business_rules: |
        func validateBusinessRules(ctx context.Context, input BusinessRulesInput) (BusinessRulesOutput, error) {
            // Custom business validation logic
            return BusinessRulesOutput{Valid: true}, nil
        }
```

### Plugin System

```yaml
# Plugin configuration
generation:
  plugins:
    - name: "metrics-plugin"
      version: "1.0.0"
      config:
        provider: "prometheus"
        namespace: "workflow"
    
    - name: "tracing-plugin"
      version: "1.0.0"
      config:
        provider: "jaeger"
        service_name: "order-processing"
    
    - name: "auth-plugin"
      version: "1.0.0"
      config:
        provider: "oauth2"
        issuer: "https://auth.company.com"
```

### Code Quality and Testing

```yaml
# Quality configuration
generation:
  quality:
    # Code formatting
    formatting:
      enabled: true
      tools:
        go: ["gofmt", "goimports"]
        typescript: ["prettier", "eslint"]
        python: ["black", "isort"]
        java: ["google-java-format"]
    
    # Linting
    linting:
      enabled: true
      tools:
        go: ["golangci-lint"]
        typescript: ["eslint", "tslint"]
        python: ["flake8", "pylint"]
        java: ["checkstyle", "spotbugs"]
    
    # Testing
    testing:
      enabled: true
      coverage_threshold: 80
      test_types:
        - "unit"
        - "integration"
        - "contract"
      
      frameworks:
        go: "testify"
        typescript: "jest"
        python: "pytest"
        java: "junit5"
```

## Best Practices

### 1. Workflow Design for Code Generation

**Keep workflows language-agnostic:**
```yaml
# Good: Generic step definition
- name: "validate_input"
  type: "validation"
  config:
    schema: "$.input_schema"
    strict: true

# Avoid: Language-specific implementations
- name: "validate_input"
  type: "custom"
  config:
    language: "go"
    code: "func validate() { ... }"
```

**Use clear, descriptive names:**
```yaml
# Good: Clear naming
- name: "calculate_order_total"
- name: "validate_payment_method"
- name: "reserve_inventory_items"

# Avoid: Generic naming
- name: "step1"
- name: "process"
- name: "handle"
```

### 2. Template Customization

**Extend rather than replace:**
```yaml
# custom-templates/go/pkg/client.go.tmpl
{{/* Extend base template */}}
{{template "base-client" .}}

{{/* Add custom methods */}}
// Custom method for batch processing
func (c *{{.workflow.name | pascalCase}}Client) ExecuteBatch(ctx context.Context, inputs []{{.workflow.name | pascalCase}}Input) ([]{{.workflow.name | pascalCase}}Output, error) {
    // Custom batch implementation
}
```

**Use template inheritance:**
```yaml
# templates/shared/base-client.tmpl
{{define "base-client"}}
type {{.workflow.name | pascalCase}}Client struct {
    baseURL string
    client  *http.Client
}
{{end}}

# templates/go/pkg/client.go.tmpl
{{template "base-client" .}}

// Language-specific additions
func (c *{{.workflow.name | pascalCase}}Client) Close() error {
    // Go-specific cleanup
}
```

### 3. Error Handling

**Generate comprehensive error handling:**
```yaml
generation:
  error_handling:
    strategy: "comprehensive"
    include_retry_logic: true
    include_circuit_breaker: true
    include_timeout_handling: true
    custom_error_types: true
```

### 4. Documentation Generation

**Include comprehensive documentation:**
```yaml
generation:
  documentation:
    include_api_docs: true
    include_examples: true
    include_troubleshooting: true
    format: "markdown"
    generate_openapi: true
```

### 5. Versioning and Compatibility

**Plan for version compatibility:**
```yaml
generation:
  versioning:
    strategy: "semantic"
    backward_compatibility: true
    deprecation_warnings: true
    migration_guides: true
```

This comprehensive code generation documentation provides everything needed to understand and effectively use Magic Flow v2's powerful code generation capabilities.