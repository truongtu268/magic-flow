# Magic Flow v2 - Complete Workflow Service Platform

> A comprehensive workflow service platform with drag-and-drop designer, automatic code generation, database-backed storage, high-performance caching, and Redis/Elasticsearch-like deployment simplicity.

## 🚀 Overview

Magic Flow v2 is a complete workflow service platform that transforms how you build, deploy, and manage workflows. Like Redis or Elasticsearch, it's designed as a standalone service that can be easily deployed and integrated into any architecture. With its visual drag-and-drop designer, automatic code generation, and comprehensive API ecosystem, Magic Flow v2 makes workflow development accessible to both developers and business users.

### 🌟 Platform Highlights

- **🎨 Visual Workflow Designer**: Drag-and-drop interface for creating complex workflows
- **🤖 Automatic Code Generation**: Generate type-safe client libraries in multiple languages
- **🗄️ Database-Backed Storage**: Persistent workflow metadata with versioning and history
- **⚡ High-Performance Caching**: Redis-based caching for lightning-fast execution
- **📡 API-First Design**: RESTful APIs for workflow management, execution, and monitoring
- **📊 Real-Time Dashboard**: Live workflow monitoring with custom metrics and alerting
- **🚀 Easy Deployment**: Single-binary deployment like Redis/Elasticsearch
- **🔄 Workflow Versioning**: Complete version management with rollback capabilities
- **🌐 Multi-Language Support**: Client libraries for Go, TypeScript, Python, and more

### 🔧 Core Features

- **🔄 Producer-Consumer Architecture**: Decoupled workflow orchestration and execution
- **🚀 Event-Driven & API-Triggered**: Support for both event-driven and manual API triggers
- **💓 Comprehensive Health Monitoring**: Service-level and logic-level health checks
- **⏱️ Advanced Timeout Management**: Hierarchical timeouts with configurable actions
- **🔄 Circuit Breaker Pattern**: Prevent cascading failures with configurable thresholds
- **🔁 Intelligent Retry**: Multiple retry strategies with exponential backoff and jitter
- **🛡️ Fallback Mechanisms**: Graceful degradation with fallback steps and default responses
- **📝 YAML-Based Definitions**: Declarative workflow and context definitions
- **🔗 Multiple Message Brokers**: Support for Kafka, Redis, and RabbitMQ
- **💾 Flexible State Storage**: Redis (in-memory) and PostgreSQL backends
- **🛠️ Advanced Code Generation**: Automatic generation of service and executor code
- **🔧 Full v1 Compatibility**: All features from v1 with enhanced capabilities
- **📊 Built-in Monitoring**: Metrics, logging, and distributed tracing
- **🔒 Security First**: Authentication, authorization, and encryption support
- **🛡️ Traffic Protection**: Anti-spike protection through message queue buffering

## 🏗️ Platform Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           Magic Flow v2 Platform                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Web Frontend   │  │   REST APIs     │  │   Dashboard     │  │ Admin Panel │ │
│  │ • Drag & Drop   │  │ • Workflow CRUD │  │ • Real-time     │  │ • Config    │ │
│  │ • YAML Editor   │  │ • Execution     │  │ • Metrics       │  │ • Users     │ │
│  │ • Code Gen      │  │ • Monitoring    │  │ • Alerting      │  │ • Security  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ Workflow Engine │  │ Code Generator  │  │ Version Manager │  │ Metadata    │ │
│  │ • Orchestration │  │ • Multi-lang    │  │ • Versioning    │  │ Store       │ │
│  │ • Execution     │  │ • Type Safety   │  │ • Rollback      │  │ • Database  │ │
│  │ • State Mgmt    │  │ • Templates     │  │ • Migration     │  │ • Search    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ Message Broker  │  │   Cache Layer   │  │   Monitoring    │  │  Security   │ │
│  │ • Task Queue    │  │ • Redis Cache   │  │ • Metrics       │  │ • Auth      │ │
│  │ • Event Stream  │  │ • Fast Access   │  │ • Tracing       │  │ • AuthZ     │ │
│  │ • Reliability   │  │ • Invalidation  │  │ • Logging       │  │ • Encryption│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Core Platform Components

#### 1. **Visual Workflow Designer**
   - **Drag-and-Drop Interface**: Intuitive workflow creation with visual components
   - **Real-Time YAML Generation**: Automatic translation from visual to YAML
   - **Template Library**: Pre-built workflow templates and patterns
   - **Validation Engine**: Real-time workflow validation and error checking
   - **Collaboration Features**: Multi-user editing and workflow sharing

#### 2. **Workflow Service Platform**
   - **Database-Backed Storage**: Persistent workflow metadata with PostgreSQL/MongoDB
   - **High-Performance Caching**: Redis-based caching for sub-millisecond access
   - **API Gateway**: RESTful APIs for all platform operations
   - **Version Management**: Complete workflow versioning with rollback capabilities
   - **Multi-Tenant Support**: Isolated workspaces for different teams/projects

#### 3. **Code Generation Engine**
   - **Multi-Language Support**: Generate clients for Go, TypeScript, Python, Java
   - **Type-Safe Generation**: Strongly-typed interfaces and data structures
   - **SDK Creation**: Complete SDKs with authentication and error handling
   - **Documentation Generation**: Automatic API documentation and examples
   - **Custom Templates**: Extensible template system for custom code generation

#### 4. **Real-Time Dashboard**
   - **Live Monitoring**: Real-time workflow execution tracking
   - **Custom Metrics**: Configurable business and technical metrics
   - **Interactive Visualizations**: Workflow flow diagrams and execution timelines
   - **Alerting System**: Configurable alerts with multiple notification channels
   - **Performance Analytics**: Bottleneck detection and optimization suggestions

#### 5. **Execution Engine**
   - **Distributed Execution**: Scalable workflow execution across multiple nodes
   - **Event-Driven & API-Triggered**: Support for both execution models
   - **Resilience Patterns**: Circuit breakers, retries, and fallback mechanisms
   - **State Management**: Persistent execution state with recovery capabilities
   - **Resource Optimization**: Intelligent scheduling and resource allocation

#### 6. **Data Flow & Integration**
   - **Data Retrieval Workflows**: Specialized patterns for data processing
   - **Streaming Support**: Real-time data processing capabilities
   - **External Integrations**: Connectors for databases, APIs, and services
   - **Data Transformation**: Built-in data mapping and transformation tools
   - **Quality Assurance**: Data validation and quality monitoring

## 🛡️ Resilience Patterns

Magic Flow v2 provides comprehensive resilience patterns to handle failures gracefully:

### Timeout Actions
- **`cancel_and_cleanup`**: Stop execution and clean up resources (default for workflows)
- **`retry`**: Retry the step with configurable backoff strategies
- **`fallback`**: Execute an alternative step or provide default response
- **`ignore`**: Continue workflow execution (useful for non-critical steps)

### Circuit Breaker States
- **Closed**: Normal operation, requests pass through
- **Open**: Failures exceeded threshold, requests fail fast
- **Half-Open**: Testing recovery, limited requests allowed

### Retry Strategies
- **Linear**: Fixed delay between retries
- **Exponential**: Exponentially increasing delays
- **Fixed**: Same delay for all retries
- **Jitter**: Add randomness to prevent thundering herd

### Best Practices
- Use **circuit breakers** for external service calls
- Apply **retry** for transient failures (network, temporary unavailability)
- Implement **fallbacks** for non-critical functionality
- Set **ignore** timeout action for optional steps (logging, analytics)
- Configure **jitter** in high-concurrency scenarios

## 📋 Quick Start

### 1. Easy Deployment

#### Single Binary Deployment (Redis/Elasticsearch-like)

```bash
# Download and run Magic Flow v2
wget https://releases.magicflow.dev/v2/magicflow-v2-linux-amd64.tar.gz
tar -xzf magicflow-v2-linux-amd64.tar.gz

# Start with default configuration
./magicflow start

# Or with custom configuration
./magicflow start --config config.yaml --port 8080
```

#### Docker Deployment

```bash
# Quick start with Docker
docker run -d \
  --name magicflow \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/data:/data \
  magicflow/magicflow:v2

# Access the platform
open http://localhost:8080  # Web Interface
open http://localhost:9090  # Dashboard
```

#### Kubernetes Deployment

```yaml
# k8s/magicflow.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: magicflow
spec:
  replicas: 3
  selector:
    matchLabels:
      app: magicflow
  template:
    metadata:
      labels:
        app: magicflow
    spec:
      containers:
      - name: magicflow
        image: magicflow/magicflow:v2
        ports:
        - containerPort: 8080
        - containerPort: 9090
        env:
        - name: MAGICFLOW_MODE
          value: "cluster"
        - name: MAGICFLOW_DATABASE_URL
          value: "postgresql://user:pass@postgres:5432/magicflow"
---
apiVersion: v1
kind: Service
metadata:
  name: magicflow-service
spec:
  selector:
    app: magicflow
  ports:
  - name: web
    port: 8080
    targetPort: 8080
  - name: dashboard
    port: 9090
    targetPort: 9090
  type: LoadBalancer
```

### 2. Platform Configuration

#### Core Configuration

```yaml
# config/magicflow.yaml
server:
  port: 8080
  dashboard_port: 9090
  mode: "standalone"  # standalone, cluster

database:
  type: "postgresql"  # postgresql, mongodb
  url: "postgresql://user:pass@localhost:5432/magicflow"
  max_connections: 100
  migration: true

cache:
  type: "redis"
  url: "redis://localhost:6379/0"
  ttl: "1h"
  max_memory: "512mb"

messaging:
  type: "redis"  # redis, rabbitmq, kafka, nats
  url: "redis://localhost:6379/1"
  queue_prefix: "magicflow"

security:
  auth_enabled: true
  jwt_secret: "your-secret-key"
  session_timeout: "24h"

monitoring:
  metrics_enabled: true
  tracing_enabled: true
  log_level: "info"
```

### 3. Visual Workflow Designer

#### Drag-and-Drop Interface

Access the visual designer at `http://localhost:8080/designer`:

1. **Create New Workflow**
   - Drag components from the palette
   - Connect steps with visual connectors
   - Configure step properties in the sidebar
   - Real-time YAML generation

2. **Component Library**
   - **Service Calls**: HTTP APIs, gRPC services, database operations
   - **Data Processing**: Transformations, validations, aggregations
   - **Control Flow**: Conditions, loops, parallel execution
   - **Integrations**: Third-party services, webhooks, notifications

3. **Template Gallery**
   - E-commerce order processing
   - Data pipeline workflows
   - Microservices orchestration
   - CI/CD automation

#### YAML Workflow Definition

The visual designer generates YAML like this:

```yaml
# workflow.yaml
apiVersion: workflow.magic-flow.io/v2
kind: Workflow
metadata:
  name: user-onboarding
  version: v1.0.0
  description: "Complete user onboarding process"
  tags: ["user-management", "onboarding"]
  owner: "team@company.com"
  created_by: "visual_designer"

spec:
  timeout: 10m
  timeout_action: cancel_and_cleanup
  retry_policy:
    max_attempts: 3
    backoff: exponential
  
  steps:
    - name: validate-user
      type: function
      timeout: 30s
      timeout_action: retry
      retry_config:
        max_attempts: 3
        backoff_strategy: exponential
        initial_delay: 1s
        max_delay: 10s
        jitter: true
      executor: user-service
      function: validateUserData
      data_mapping:
        input:
          user_data: "$.input"
        output:
          validation_result: "$.response.valid"
          validation_errors: "$.response.errors"
      on_success: create-account
      on_failure: notify-validation-error
    
    - name: create-account
      type: function
      timeout: 1m
      timeout_action: fallback
      fallback_step: create-guest-account
      condition: "$.validate-user.validation_result == true"
      executor: account-service
      function: createUserAccount
      data_mapping:
        input:
          email: "$.input.email"
          name: "$.input.name"
          preferences: "$.input.preferences"
      on_success: send-welcome-email
      on_failure: cleanup-partial-data
    
    - name: external-verification
      type: function
      timeout: 15s
      timeout_action: fallback
      fallback_step: skip-verification
      circuit_breaker:
        failure_threshold: 5
        recovery_timeout: 60s
        half_open_max_calls: 3
        failure_rate_threshold: 0.5
      executor: verification-service
      function: verifyUser
      on_success: send-welcome-email
      on_failure: send-welcome-email
    
    - name: send-welcome-email
      type: function
      timeout: 30s
      timeout_action: ignore  # Non-critical step
      executor: notification-service
      function: sendWelcomeEmail
      data_mapping:
        input:
          user_id: "$.create-account.user_id"
          email: "$.input.email"
          name: "$.input.name"

  error_handling:
    - condition: "$.validate-user.validation_result == false"
      action:
        type: "return_error"
        message: "User validation failed"
        details: "$.validate-user.validation_errors"
    
    - condition: "$.create-account.status == 'failed'"
      action:
        type: "compensate"
        steps:
          - name: "cleanup-partial-data"
            service: "user-service"
            method: "cleanup"
```

### 4. Define Context Schema

```yaml
# context.yaml
apiVersion: workflow.magic-flow.io/v2
kind: WorkflowContext
metadata:
  name: user-onboarding-context

spec:
  input_schema:
    type: object
    properties:
      email:
        type: string
        format: email
        description: "User's email address"
        required: true
      name:
        type: string
        description: "User's full name"
        required: true
      preferences:
        type: object
        description: "User preferences"
        properties:
          newsletter:
            type: boolean
            default: false
          notifications:
            type: boolean
            default: true
      metadata:
        type: object
        description: "Additional user metadata"
        properties:
          source:
            type: string
            enum: ["web", "mobile", "api"]
          campaign_id:
            type: string
  
  output_schema:
    type: object
    properties:
      user_id:
        type: string
        description: "Generated user ID"
      account_status:
        type: string
        enum: ["active", "pending", "failed"]
        description: "Final account status"
      created_at:
        type: string
        format: date-time
        description: "Account creation timestamp"
      verification_status:
        type: string
        enum: ["verified", "pending", "skipped"]
      welcome_email_sent:
        type: boolean
        description: "Whether welcome email was sent successfully"
```

### 5. Generate Code

```bash
# Install Magic Flow v2 CLI
go install github.com/truongtu268/magic-flow/cmd/magic-flow@v2

# Generate workflow service
magic-flow generate server \
  --workflow workflow.yaml \
  --context context.yaml \
  --config config/magic-flow.yaml \
  --output ./generated/server \
  --language go

# Generate workflow executor
magic-flow generate executor \
  --workflow workflow.yaml \
  --context context.yaml \
  --config config/magic-flow.yaml \
  --output ./generated/executor \
  --language go

# Generate client SDKs for multiple languages
magic-flow generate client \
  --workflow workflow.yaml \
  --context context.yaml \
  --output ./generated/clients \
  --languages go,typescript,python,java

# Generate API documentation
magic-flow generate docs \
  --workflow workflow.yaml \
  --context context.yaml \
  --output ./docs/api \
  --format openapi
```

### 6. Implement Business Logic

```go
// In your executor service
package main

import (
    "context"
    "github.com/truongtu268/magic-flow/pkg_v2/runtime/executor"
    "./generated/executor"
)

func main() {
    // Create executor with business logic
    exec := executor.New()
    
    // Register your business functions using generated constants
    exec.RegisterFunction(StepValidateUser, validateUserData)
    exec.RegisterFunction(StepCreateAccount, createUserAccount)
    exec.RegisterFunction(StepSendWelcomeEmail, sendWelcomeEmail)
    
    // Start the executor
    exec.Start()
}

func validateUserData(ctx context.Context, input *UserOnboardingInput) (*ValidateUserOutput, error) {
    // Your validation logic here
    email := input.Email
    name := input.Name
    
    // Validate email format, check for duplicates, etc.
    if !isValidEmail(email) {
        return nil, errors.New("invalid email format")
    }
    
    return &ValidateUserOutput{
        Validated: true,
        UserData:  input,
    }, nil
}
```

### 7. Trigger Workflows from Business Logic

Send workflow initialization events from your business services:

```go
package main

import (
    "context"
    "time"
    "github.com/truongtu268/magic-flow/generated"
)

func main() {
    // Initialize message broker client
    broker := generated.NewMessageBroker()
    
    // Send workflow initialization event
    event := &generated.WorkflowInitEvent{
        WorkflowName: "user-onboarding",
        InitialData: &generated.UserOnboardingInput{
            Email:    "user@example.com",
            Name:     "John Doe",
            Preferences: map[string]bool{
                "newsletter":    true,
                "notifications": false,
            },
        },
        Source:    "user-service",
        Timestamp: time.Now(),
    }
    
    err := broker.PublishWorkflowEvent(context.Background(), event)
    if err != nil {
        panic(err)
    }
}
```

### 8. API-Triggered Workflows

#### Manual Workflow Execution

```bash
# Execute workflow via REST API
curl -X POST http://localhost:8080/api/v1/workflows/order_processing/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-token" \
  -d '{
    "input": {
      "order_id": "order-123",
      "customer_id": "customer-456",
      "items": [
        {
          "product_id": "product-789",
          "quantity": 2,
          "price": 29.99
        }
      ],
      "payment_method": "credit_card",
      "shipping_address": {
        "street": "123 Main St",
        "city": "San Francisco",
        "state": "CA",
        "zip": "94105"
      }
    },
    "metadata": {
      "source": "web",
      "user_agent": "Mozilla/5.0..."
    }
  }'

# Response
{
  "execution_id": "exec-uuid-123",
  "workflow_id": "order_processing",
  "status": "running",
  "created_at": "2024-01-15T10:30:00Z",
  "estimated_completion": "2024-01-15T10:32:00Z",
  "tracking_url": "http://localhost:9090/executions/exec-uuid-123"
}
```

#### Data Retrieval Workflows

```bash
# Get workflow execution status
curl -X GET http://localhost:8080/api/v1/executions/exec-uuid-123 \
  -H "Authorization: Bearer your-api-token"

# Get workflow execution results
curl -X GET http://localhost:8080/api/v1/executions/exec-uuid-123/results \
  -H "Authorization: Bearer your-api-token"

# Get workflow metrics
curl -X GET http://localhost:8080/api/v1/workflows/order_processing/metrics \
  -H "Authorization: Bearer your-api-token"

# Stream real-time execution events
curl -X GET http://localhost:8080/api/v1/executions/exec-uuid-123/stream \
  -H "Authorization: Bearer your-api-token" \
  -H "Accept: text/event-stream"
```

### 9. Real-Time Dashboard

#### Access Dashboard

Open `http://localhost:9090` to access the real-time dashboard:

1. **Workflow Overview**
   - Active workflows and execution counts
   - Success/failure rates
   - Average execution times
   - Resource utilization

2. **Live Execution Monitoring**
   - Real-time workflow execution tracking
   - Step-by-step progress visualization
   - Error detection and alerting
   - Performance bottleneck identification

3. **Custom Metrics Dashboard**
   - Business metrics from workflow data
   - Custom charts and visualizations
   - Configurable alerts and notifications
   - Historical trend analysis

#### Custom Metrics Configuration

```yaml
# config/metrics.yaml
metrics:
  business_metrics:
    - name: "order_value"
      description: "Total order value processed"
      type: "counter"
      source: "$.output.total_amount"
      labels:
        - "customer_segment"
        - "product_category"
    
    - name: "processing_time"
      description: "Order processing duration"
      type: "histogram"
      source: "$.execution.duration"
      buckets: [1, 5, 10, 30, 60, 300]
    
    - name: "payment_success_rate"
      description: "Payment processing success rate"
      type: "gauge"
      source: "$.steps.process_payment.success_rate"
      window: "1h"

  alerts:
    - name: "high_failure_rate"
      condition: "payment_success_rate < 0.95"
      severity: "critical"
      channels: ["slack", "email", "pagerduty"]
    
    - name: "slow_processing"
      condition: "avg(processing_time) > 60"
      severity: "warning"
      channels: ["slack"]
```

### 10. Start the Platform

```bash
# Start Magic Flow v2 platform
./magicflow start --config config/magicflow.yaml

# Or using Docker
docker-compose up -d

# Access the platform
open http://localhost:8080  # Main interface
open http://localhost:9090  # Dashboard
open http://localhost:8080/designer  # Visual designer
open http://localhost:8080/api/docs  # API documentation
```

Your complete workflow service platform is now running! You can:
- Design workflows visually
- Execute workflows via API
- Monitor in real-time
- Generate client code
- Scale horizontally

## 📚 Documentation

- **[Architecture Guide](./ARCHITECTURE.md)** - Detailed system architecture and design principles
- **[Implementation Planning](./PLANNING.md)** - Development roadmap and technical specifications
- **[API Reference](./docs/api/)** - Complete API documentation
- **[User Guide](./docs/guides/)** - Step-by-step usage guides
- **[Examples](./examples/)** - Working examples and patterns

## 🔧 Configuration

### Workflow Service Configuration

```yaml
# config/server.yaml
server:
  host: "0.0.0.0"
  port: 8080
  timeout: 30s

messaging:
  backend: "redis"
  redis:
    host: "localhost"
    port: 6379
    db: 0

health:
  check_interval: 30s
  timeout: 5s
  failure_threshold: 3

storage:
  backend: "postgres"
  postgres:
    host: "localhost"
    port: 5432
    database: "workflows"
    username: "workflow_user"
    password: "${DB_PASSWORD}"

monitoring:
  metrics_enabled: true
  tracing_enabled: true
  log_level: "info"
```

### Executor Configuration

```yaml
# config/executor.yaml
executor:
  name: "user-service-executor"
  concurrency: 10
  timeout: 5m

messaging:
  backend: "redis"
  redis:
    host: "localhost"
    port: 6379
    db: 0

health:
  endpoint: "/health"
  port: 8081
  checks:
    - name: "database"
      type: "postgres"
      config:
        host: "localhost"
        port: 5432
        database: "users"
    - name: "external_api"
      type: "http"
      config:
        url: "https://api.external-service.com/health"

monitoring:
  metrics_enabled: true
  log_level: "info"
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...

# Run load tests
go test -tags=load ./...

# Benchmark tests
go test -bench=. ./...
```

## 🚀 Deployment

### Docker Deployment

```dockerfile
# Dockerfile.server
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o workflow-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/workflow-service .
CMD ["./workflow-service"]
```

```bash
# Build and run with Docker Compose
docker-compose up -d
```

### Kubernetes Deployment

```yaml
# k8s/workflow-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: workflow-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: workflow-service
  template:
    metadata:
      labels:
        app: workflow-service
    spec:
      containers:
      - name: workflow-service
        image: magic-flow/workflow-service:v2.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 🔄 Migration from v1

Magic Flow v2 provides comprehensive migration tools and compatibility layers:

```bash
# Analyze existing v1 workflows
magic-flow migrate analyze --source ./v1-workflows

# Convert v1 configurations to v2 YAML
magic-flow migrate convert \
  --source ./v1-workflows \
  --output ./v2-workflows

# Validate converted workflows
magic-flow validate --workflows ./v2-workflows

# Generate migration report
magic-flow migrate report \
  --source ./v1-workflows \
  --target ./v2-workflows
```

## 📊 Monitoring & Observability

### Metrics

- **Workflow Metrics**: Execution count, duration, success rate
- **Step Metrics**: Step execution time, failure rate, retry count
- **System Metrics**: CPU, memory, queue depth, connection count
- **Business Metrics**: Custom metrics from your business logic

### Health Checks

- **Service Health**: Overall service availability
- **Component Health**: Individual component status
- **Dependency Health**: External service availability
- **Custom Health**: Business logic health checks

### Distributed Tracing

- **Request Tracing**: End-to-end request tracking
- **Workflow Tracing**: Complete workflow execution traces
- **Step Tracing**: Individual step execution details
- **Error Tracing**: Detailed error propagation tracking

## 🧪 Test-Driven Development (TDD)

Magic Flow v2 follows strict TDD principles throughout development:

### TDD Workflow

1. **🔴 Red Phase**: Write failing tests first
   - Define expected behavior through tests
   - Ensure tests fail for the right reasons
   - Cover edge cases and error scenarios

2. **🟢 Green Phase**: Write minimal code to make tests pass
   - Implement just enough functionality
   - Focus on making tests pass, not perfection
   - Maintain simplicity and clarity

3. **🔄 Refactor Phase**: Improve code while keeping tests green
   - Optimize performance and structure
   - Enhance readability and maintainability
   - Ensure all tests continue to pass

### Testing Strategy

#### Unit Tests (90%+ Coverage)
```go
// Example: Testing workflow execution engine
func TestWorkflowEngine_Execute(t *testing.T) {
    engine := NewWorkflowEngine()
    workflow := &Workflow{
        Steps: []Step{
            {ID: StepValidateUser, Function: "validateUser"},
            {ID: StepCreateAccount, Function: "createAccount"},
        },
    }
    
    result, err := engine.Execute(workflow, &UserOnboardingInput{
        Email: "test@example.com",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "completed", result.Status)
}
```

#### Integration Tests
```go
// Example: Testing message broker integration
func TestKafkaMessageBroker_Integration(t *testing.T) {
    broker := NewKafkaMessageBroker(&KafkaConfig{
        Brokers: []string{"localhost:9092"},
        Topic: "workflow-events-test",
    })
    
    event := &WorkflowInitEvent{
        WorkflowID: "test-workflow-123",
        WorkflowType: "user-onboarding",
        Input: &UserOnboardingInput{Email: "test@example.com"},
    }
    
    err := broker.Publish(event)
    assert.NoError(t, err)
    
    // Test message consumption...
}
```

#### End-to-End Tests
```go
// Example: Complete workflow execution test
func TestWorkflowService_E2E(t *testing.T) {
    service := setupTestWorkflowService(t)
    defer teardownTestWorkflowService(t, service)
    
    // Register workflow, trigger execution, verify results
    result := executeCompleteWorkflow(t, service, "user-onboarding")
    assert.Equal(t, "completed", result.Status)
    assert.Len(t, result.Steps, 3)
}
```

### Test Infrastructure

- **Mock Implementations**: Comprehensive mocking for external dependencies
- **Test Utilities**: Helper functions for common test scenarios
- **Test Data Builders**: Fluent builders for creating test data
- **Integration Test Environment**: Docker-based test environment setup
- **Continuous Testing**: Automated test execution on every commit

### Quality Metrics

- **Code Coverage**: Minimum 90% for unit tests
- **Integration Coverage**: All external dependencies tested
- **Performance Testing**: Automated benchmarks and load testing
- **Security Testing**: Automated security scans and penetration testing
- **Chaos Testing**: Resilience validation under failure conditions

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/truongtu268/magic-flow.git
cd magic-flow/pkg_v2

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run

# Build all components
go build ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🆘 Support

- **Documentation**: [docs.magic-flow.io](https://docs.magic-flow.io)
- **Issues**: [GitHub Issues](https://github.com/truongtu268/magic-flow/issues)
- **Discussions**: [GitHub Discussions](https://github.com/truongtu268/magic-flow/discussions)
- **Community**: [Discord Server](https://discord.gg/magic-flow)

## 🗺️ Roadmap

### v2.1 (Q2 2024)
- [ ] Visual workflow designer
- [ ] Advanced scheduling capabilities
- [ ] Workflow versioning and blue-green deployments
- [ ] Enhanced monitoring dashboards

### v2.2 (Q3 2024)
- [ ] Multi-tenancy support
- [ ] Workflow analytics and optimization
- [ ] Cloud provider integrations
- [ ] Kubernetes operator

### v2.3 (Q4 2024)
- [ ] Advanced workflow patterns
- [ ] Machine learning integration
- [ ] Performance optimizations
- [ ] Enterprise features

---

**Magic Flow v2** - Building the future of distributed workflow orchestration 🚀