# Magic Flow v2 - Distributed Workflow Service Architecture

## Overview

Magic Flow v2 is a comprehensive workflow service platform that provides a complete ecosystem for defining, managing, executing, and monitoring workflows. The platform combines a powerful workflow engine with an intuitive management interface, automatic code generation, and enterprise-grade deployment capabilities.

### Key Features

- **YAML-Based Workflow Definitions**: Define workflows using declarative YAML syntax
- **Drag-and-Drop Visual Designer**: Create workflows visually with an intuitive web interface
- **Automatic Code Generation**: Generate type-safe client code and execution services
- **API-Driven Workflow Triggers**: Manual and automated workflow execution via REST APIs
- **Real-Time Dashboard**: Monitor workflow execution, metrics, and performance
- **Database-Backed Metadata Storage**: Persistent workflow definitions and execution history
- **High-Performance Caching**: Fast workflow execution with intelligent caching
- **Easy Deployment**: Redis/Elasticsearch-like deployment with internal dependency management
- **Workflow Versioning**: Complete version control for workflow definitions and client code
- **Enterprise Scalability**: Distributed architecture supporting high-throughput execution

## Architecture Principles

### 1. Producer-Consumer Pattern
- **Workflow Service (Producer)**: Orchestrates workflow definitions and manages execution state
- **Workflow Executor (Consumer)**: Executes individual workflow steps and reports back results
- **Message Broker**: Decouples services and provides reliable communication with multiple backend options

### 2. Health Monitoring
- **Service-Level Health**: Workflow service monitors executor service availability
- **Logic-Level Health**: Executor services monitor individual workflow step health
- **Bidirectional Monitoring**: Both services implement health check endpoints

### 3. Timeout Management & Resilience Patterns
- **Unified Timeout Pattern**: Consistent timeout handling across all components
- **Configurable Timeouts**: Per-workflow, per-step, and global timeout settings
- **Configurable Timeout Actions**: Define specific behaviors for timeout scenarios
- **Circuit Breaker Pattern**: Prevent cascading failures and resource exhaustion
- **Retry Mechanisms**: Configurable retry strategies with exponential backoff
- **Graceful Degradation**: Proper cleanup and error handling on timeouts
- **Timeout Propagation**: Cascading timeout behavior across workflow steps

### 4. Data Flow & Workflow Management
- **Data Flow Workflows**: Structured data processing pipelines with type safety
- **API-Triggered Execution**: RESTful APIs for manual and automated workflow triggers
- **Real-Time Monitoring**: Live dashboard for workflow execution and performance metrics
- **Metadata Management**: Database-backed storage for workflow definitions and execution history
- **Intelligent Caching**: Multi-layer caching for fast workflow execution and data retrieval

### 5. Visual Workflow Design & Code Generation
- **Drag-and-Drop Interface**: Web-based visual workflow designer
- **YAML Translation**: Automatic conversion from visual design to YAML definitions
- **Type-Safe Code Generation**: Generate client libraries and execution services
- **Version Control Integration**: Git-like versioning for workflow definitions
- **Backward Compatibility**: Support for multiple workflow versions simultaneously

## Workflow Service Platform Architecture

### Platform Components Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           Magic Flow v2 Platform                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Web Frontend  │  │   REST APIs     │  │   Dashboard     │  │   Admin     │ │
│  │                 │  │                 │  │                 │  │   Panel     │ │
│  │ • Drag & Drop   │  │ • Workflow CRUD │  │ • Real-time     │  │ • User Mgmt │ │
│  │ • Visual Design │  │ • Execution     │  │   Monitoring    │  │ • Config    │ │
│  │ • YAML Editor   │  │ • Metrics       │  │ • Performance   │  │ • Security  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Workflow       │  │  Code Generator │  │  Version        │  │  Metadata   │ │
│  │  Engine         │  │                 │  │  Manager        │  │  Store      │ │
│  │                 │  │ • Client Code   │  │                 │  │             │ │
│  │ • Orchestration │  │ • Server Code   │  │ • Git-like      │  │ • PostgreSQL│ │
│  │ • Execution     │  │ • Type Safety   │  │   Versioning    │  │ • MongoDB   │ │
│  │ • State Mgmt    │  │ • Templates     │  │ • Migration     │  │ • Schemas   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Message        │  │  Cache Layer    │  │  Monitoring     │  │  Security   │ │
│  │  Broker         │  │                 │  │                 │  │             │ │
│  │                 │  │ • Redis Cache   │  │ • Prometheus    │  │ • OAuth 2.0 │ │
│  │ • Kafka/Redis   │  │ • Workflow Data │  │ • Grafana       │  │ • JWT Auth  │ │
│  │ • RabbitMQ      │  │ • Execution     │  │ • Alerting      │  │ • RBAC      │ │
│  │ • Pub/Sub       │  │   Results       │  │ • Tracing       │  │ • Audit Log │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Data Flow Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Workflow      │    │   Execution     │
│   Interface     │    │                 │    │   Service       │    │   Engine        │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │                       │
         │ 1. Design Workflow    │                       │                       │
         ├──────────────────────►│                       │                       │
         │                       │ 2. Save Definition   │                       │
         │                       ├──────────────────────►│                       │
         │                       │                       │ 3. Generate Code     │
         │                       │                       ├──────────────────────►│
         │                       │                       │                       │
         │ 4. Trigger Execution  │                       │                       │
         ├──────────────────────►│                       │                       │
         │                       │ 5. Execute Workflow  │                       │
         │                       ├──────────────────────►│                       │
         │                       │                       │ 6. Process Steps     │
         │                       │                       ├──────────────────────►│
         │                       │                       │                       │
         │ 7. Monitor Progress   │                       │                       │
         │◄──────────────────────┼───────────────────────┼───────────────────────┤
         │                       │                       │                       │

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Database      │    │   Cache Layer   │    │   Message       │
│   Storage       │    │                 │    │   Broker        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ Workflow Metadata     │ Execution Cache       │ Event Streaming
         │ Execution History     │ Performance Data      │ Real-time Updates
         │ User Management       │ Session Data          │ Notifications
         │ Version Control       │ Temporary Results     │ Async Processing
```

## Message Broker Configuration

The system supports multiple message broker backends for different use cases:

### Apache Kafka
- **Use Case**: High-throughput, distributed environments
- **Features**: Partitioning, replication, stream processing
- **Configuration**: Bootstrap servers, topic configuration, consumer groups

### Redis
- **Use Case**: Low-latency messaging with optional persistence
- **Features**: Pub/Sub, streams, clustering support
- **Configuration**: Connection pooling, persistence settings, cluster mode

### RabbitMQ
- **Use Case**: Feature-rich messaging with complex routing
- **Features**: Exchanges, queues, routing keys, dead letter queues
- **Configuration**: Virtual hosts, exchanges, queue durability

## State Storage Configuration

Workflow state can be persisted using different storage backends:

### Redis (In-Memory)
- **Use Case**: High-performance, low-latency state access
- **Features**: Atomic operations, expiration, clustering
- **Configuration**: Memory policies, persistence (RDB/AOF), clustering

### PostgreSQL
- **Use Case**: Complex state queries, ACID compliance, reporting
- **Features**: JSON/JSONB support, transactions, indexing
- **Configuration**: Connection pooling, schema management, partitioning

## API Specifications

### Workflow Management APIs

#### Workflow CRUD Operations
```http
# Create new workflow
POST /api/v2/workflows
Content-Type: application/json

{
  "name": "user-onboarding",
  "description": "Complete user onboarding process",
  "version": "1.0.0",
  "yaml_definition": "...",
  "tags": ["user", "onboarding"],
  "metadata": {
    "owner": "team-backend",
    "environment": "production"
  }
}

# Get workflow by ID
GET /api/v2/workflows/{workflow_id}

# Update workflow
PUT /api/v2/workflows/{workflow_id}

# Delete workflow
DELETE /api/v2/workflows/{workflow_id}

# List workflows with filtering
GET /api/v2/workflows?tags=user&status=active&page=1&limit=20
```

#### Workflow Execution APIs
```http
# Trigger workflow execution
POST /api/v2/workflows/{workflow_id}/execute
Content-Type: application/json

{
  "input_data": {
    "user_email": "user@example.com",
    "user_name": "John Doe"
  },
  "priority": "high",
  "timeout_override": "30m",
  "metadata": {
    "source": "web-app",
    "request_id": "req-123"
  }
}

# Get execution status
GET /api/v2/executions/{execution_id}

# Cancel execution
POST /api/v2/executions/{execution_id}/cancel

# Retry failed execution
POST /api/v2/executions/{execution_id}/retry

# Get execution logs
GET /api/v2/executions/{execution_id}/logs?level=info&limit=100
```

#### Metrics and Monitoring APIs
```http
# Get workflow metrics
GET /api/v2/workflows/{workflow_id}/metrics?period=24h&granularity=1h

# Get system metrics
GET /api/v2/metrics/system?metrics=cpu,memory,throughput

# Get execution statistics
GET /api/v2/metrics/executions?start_time=2024-01-01&end_time=2024-01-31

# Create custom metric
POST /api/v2/metrics/custom
{
  "name": "user_conversion_rate",
  "description": "Rate of successful user onboarding",
  "query": "SELECT COUNT(*) FROM executions WHERE workflow_name='user-onboarding' AND status='completed'",
  "aggregation": "rate",
  "tags": ["business", "conversion"]
}
```

### Code Generation APIs
```http
# Generate client code
POST /api/v2/workflows/{workflow_id}/generate/client
{
  "language": "go",
  "package_name": "workflows",
  "output_format": "zip",
  "include_tests": true,
  "version": "1.0.0"
}

# Generate server code
POST /api/v2/workflows/{workflow_id}/generate/server
{
  "language": "go",
  "framework": "gin",
  "include_docker": true,
  "include_k8s": true
}

# Get generated code status
GET /api/v2/generation/{generation_id}/status

# Download generated code
GET /api/v2/generation/{generation_id}/download
```

### Version Management APIs
```http
# Create new version
POST /api/v2/workflows/{workflow_id}/versions
{
  "version": "1.1.0",
  "changelog": "Added error handling for payment failures",
  "yaml_definition": "...",
  "migration_notes": "No breaking changes"
}

# List versions
GET /api/v2/workflows/{workflow_id}/versions

# Compare versions
GET /api/v2/workflows/{workflow_id}/versions/compare?from=1.0.0&to=1.1.0

# Rollback to version
POST /api/v2/workflows/{workflow_id}/rollback
{
  "target_version": "1.0.0",
  "force": false
}
```

## Dashboard and Visualization Features

### Real-Time Monitoring Dashboard

#### Workflow Execution Overview
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          Workflow Execution Dashboard                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Active        │  │   Completed     │  │   Failed        │  │   Queued    │ │
│  │   Workflows     │  │   Today         │  │   Last Hour     │  │   Pending   │ │
│  │                 │  │                 │  │                 │  │             │ │
│  │      1,247      │  │      8,932      │  │        23       │  │      156    │ │
│  │   ↑ 12% (24h)   │  │   ↑ 5% (24h)    │  │   ↓ 15% (24h)   │  │  ↑ 8% (1h)  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐  ┌─────────────────────────────────────┐ │
│  │        Execution Timeline           │  │         Performance Metrics         │ │
│  │                                     │  │                                     │ │
│  │  ████████████████████████████████   │  │  Avg Response Time: 2.3s           │ │
│  │  ████████████████████████████████   │  │  95th Percentile: 8.1s             │ │
│  │  ████████████████████████████████   │  │  Throughput: 450 workflows/min     │ │
│  │  ████████████████████████████████   │  │  Error Rate: 0.8%                  │ │
│  │                                     │  │  Cache Hit Rate: 94.2%             │ │
│  └─────────────────────────────────────┘  └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Workflow Visual Flow Monitoring
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        User Onboarding Workflow - Live View                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │
│  │  Validate   │───►│   Create    │───►│    Send     │───►│   Complete  │      │
│  │    User     │    │   Account   │    │   Welcome   │    │  Onboarding │      │
│  │             │    │             │    │    Email    │    │             │      │
│  │   ✅ 1,234   │    │   ✅ 1,198   │    │   ✅ 1,156   │    │   ✅ 1,142   │      │
│  │   ❌ 12      │    │   ❌ 36      │    │   ❌ 42      │    │   ❌ 14      │      │
│  │   ⏳ 45      │    │   ⏳ 23      │    │   ⏳ 18      │    │   ⏳ 8       │      │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘      │
│                                                                                 │
│  Step Performance:                                                              │
│  • Validate User: Avg 0.8s (Target: <1s) ✅                                    │
│  • Create Account: Avg 2.1s (Target: <3s) ✅                                   │
│  • Send Welcome Email: Avg 1.2s (Target: <2s) ✅                               │
│  • Complete Onboarding: Avg 0.3s (Target: <0.5s) ✅                            │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Custom Metrics Dashboard
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            Custom Business Metrics                             │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐  ┌─────────────────────────────────────┐ │
│  │        User Conversion Rate         │  │         Revenue Impact             │ │
│  │                                     │  │                                     │ │
│  │         94.2%                       │  │        $127,450                     │ │
│  │      ↑ 2.1% (week)                  │  │     ↑ $12,340 (week)               │ │
│  │                                     │  │                                     │ │
│  │  Target: 95% ████████████████████▓  │  │  Monthly Target: $500K ██████▓▓▓▓  │ │
│  └─────────────────────────────────────┘  └─────────────────────────────────────┘ │
│  ┌─────────────────────────────────────┐  ┌─────────────────────────────────────┐ │
│  │      Payment Success Rate           │  │        Support Ticket Volume       │ │
│  │                                     │  │                                     │ │
│  │         98.7%                       │  │           23 tickets                │ │
│  │      ↑ 0.3% (day)                   │  │        ↓ 8 tickets (day)           │ │
│  └─────────────────────────────────────┘  └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Alerting and Notification System
```yaml
# Alert Configuration
alerting:
  rules:
    - name: high_failure_rate
      condition: "workflow_failure_rate > 5%"
      duration: "5m"
      severity: "critical"
      notifications:
        - type: "slack"
          channel: "#alerts"
        - type: "email"
          recipients: ["team@company.com"]
        - type: "pagerduty"
          service_key: "workflow-service"
    
    - name: slow_execution
      condition: "avg_execution_time > 30s"
      duration: "10m"
      severity: "warning"
      notifications:
        - type: "slack"
          channel: "#monitoring"
    
    - name: queue_backlog
      condition: "pending_workflows > 1000"
      duration: "2m"
      severity: "warning"
      auto_scale:
        enabled: true
        max_instances: 10
```

## Easy Deployment Architecture

### Single Binary Deployment (Redis/Elasticsearch-like)

#### All-in-One Deployment
```bash
# Download and run Magic Flow v2
wget https://releases.magic-flow.io/v2.0.0/magic-flow-linux-amd64
chmod +x magic-flow-linux-amd64

# Start with embedded dependencies
./magic-flow-linux-amd64 start --mode=standalone

# Or with custom configuration
./magic-flow-linux-amd64 start --config=config.yaml
```

#### Internal Dependency Management
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Magic Flow v2 - Single Binary                           │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Web Server    │  │   API Server    │  │   Workflow      │  │   Admin     │ │
│  │                 │  │                 │  │   Engine        │  │   Interface │ │
│  │ • Static Files  │  │ • REST APIs     │  │                 │  │             │ │
│  │ • React App     │  │ • GraphQL       │  │ • Execution     │  │ • Config UI │ │
│  │ • WebSocket     │  │ • WebSocket     │  │ • Orchestration │  │ • Monitoring│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Embedded      │  │   Embedded      │  │   Embedded      │  │   Embedded  │ │
│  │   Database      │  │   Cache         │  │   Message       │  │   Monitoring│ │
│  │                 │  │                 │  │   Queue         │  │             │ │
│  │ • SQLite/       │  │ • In-Memory     │  │ • In-Memory     │  │ • Metrics   │ │
│  │   BadgerDB      │  │   Cache         │  │   Queue         │  │ • Logging   │ │
│  │ • File Storage  │  │ • LRU/LFU       │  │ • Pub/Sub       │  │ • Tracing   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

#### Configuration Options
```yaml
# config.yaml - Standalone Mode
mode: standalone

server:
  host: "0.0.0.0"
  port: 8080
  tls:
    enabled: false
    cert_file: ""
    key_file: ""

storage:
  type: "embedded"  # embedded, postgres, mysql
  embedded:
    engine: "badger"  # badger, sqlite
    path: "./data"
    backup:
      enabled: true
      interval: "1h"
      retention: "7d"

cache:
  type: "embedded"  # embedded, redis
  embedded:
    max_memory: "512MB"
    eviction_policy: "lru"
    persistence: true

messaging:
  type: "embedded"  # embedded, kafka, redis, rabbitmq
  embedded:
    max_queue_size: 10000
    persistence: true
    durability: "memory"  # memory, disk

monitoring:
  metrics:
    enabled: true
    endpoint: "/metrics"
  logging:
    level: "info"
    format: "json"
    output: "stdout"
  tracing:
    enabled: false
    jaeger_endpoint: ""

auth:
  enabled: false
  type: "jwt"  # jwt, oauth2, basic
  jwt:
    secret: "your-secret-key"
    expiration: "24h"

workflow:
  max_concurrent_executions: 1000
  default_timeout: "30m"
  cleanup:
    completed_retention: "30d"
    failed_retention: "90d"
```

#### Docker Deployment
```dockerfile
# Dockerfile
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY magic-flow-linux-amd64 /app/magic-flow
COPY config.yaml /app/config.yaml

EXPOSE 8080

CMD ["./magic-flow", "start", "--config=config.yaml"]
```

```bash
# Build and run
docker build -t magic-flow:v2.0.0 .
docker run -p 8080:8080 -v $(pwd)/data:/app/data magic-flow:v2.0.0
```

#### Kubernetes Deployment
```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: magic-flow
  namespace: workflows
spec:
  replicas: 3
  selector:
    matchLabels:
      app: magic-flow
  template:
    metadata:
      labels:
        app: magic-flow
    spec:
      containers:
      - name: magic-flow
        image: magic-flow:v2.0.0
        ports:
        - containerPort: 8080
        env:
        - name: MODE
          value: "cluster"
        - name: STORAGE_TYPE
          value: "postgres"
        - name: POSTGRES_URL
          valueFrom:
            secretKeyRef:
              name: magic-flow-secrets
              key: postgres-url
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
        - name: data
          mountPath: /app/data
      volumes:
      - name: config
        configMap:
          name: magic-flow-config
      - name: data
        persistentVolumeClaim:
          claimName: magic-flow-data

---
apiVersion: v1
kind: Service
metadata:
  name: magic-flow-service
  namespace: workflows
spec:
  selector:
    app: magic-flow
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Compatibility with Magic Flow v1

Magic Flow v2 maintains compatibility with the existing v1 state storage patterns found in `/pkg/storage/`. Key compatibility features:

### State Structure Compatibility
- **Workflow State**: Compatible with v1 `WorkflowState` structure
- **Step State**: Maintains v1 `StepState` format for seamless migration
- **Context Storage**: Supports both v1 `map[string]interface{}` and v2 strongly-typed contexts

### Migration Support
- **Gradual Migration**: Run v1 and v2 workflows side-by-side
- **State Converter**: Automatic conversion between v1 and v2 state formats
- **Backward Compatibility**: v2 can read and update v1 workflow states

### Storage Interface Reuse
- **Interface Compatibility**: Reuses v1 storage interfaces where applicable
- **Implementation Sharing**: Leverages existing v1 Redis and database implementations
- **Configuration Migration**: Easy migration of v1 storage configurations to v2

## System Components

### Workflow Service (Producer)
```
┌─────────────────────────────────────┐
│           Workflow Service          │
├─────────────────────────────────────┤
│ • Workflow Definition Management    │
│ • Execution Orchestration           │
│ • State Management                  │
│ • Health Check Coordinator          │
│ • Timeout Management                │
│ • Message Queue Producer            │
└─────────────────────────────────────┘
```

**Responsibilities:**
- Load and validate workflow definitions from YAML
- Listen for workflow initialization events from business logic
- Orchestrate workflow execution across multiple executors
- Maintain workflow state and progress tracking
- Monitor executor service health
- Handle workflow-level timeouts and error recovery
- Publish execution tasks to message broker
- Persist workflow state to configured storage backend

### Workflow Executor (Consumer)
```
┌─────────────────────────────────────┐
│         Workflow Executor           │
├─────────────────────────────────────┤
│ • Step Execution Engine             │
│ • Business Logic Integration        │
│ • Health Check Provider             │
│ • Timeout Handler                   │
│ • Message Queue Consumer            │
│ • Result Reporter                   │
└─────────────────────────────────────┘
```

**Responsibilities:**
- Execute individual workflow steps
- Integrate with core business logic
- Provide health status for executed steps
- Handle step-level timeouts
- Consume tasks from message queue
- Report execution results back to workflow service

## Communication Flow

### Event-Driven Workflow Initialization
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Business       │    │   Message       │    │   Workflow      │
│  Logic          │    │   Broker        │    │   Service       │
│  Services       │    │                 │    │   (Producer)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ 1. Workflow Event     │                       │
         ├──────────────────────►│                       │
         │                       │ 2. Event Consumed    │
         │                       ├──────────────────────►│
         │                       │                       │
         │                       │ 3. Workflow Started  │
         │                       │                   ┌───┤
         │                       │                   │   │
         │                       │                   └──►│
```

### Workflow Execution Flow
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Workflow       │    │   Message       │    │   Workflow      │
│  Service        │    │   Broker        │    │   Executor      │
│  (Producer)     │    │                 │    │   (Consumer)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ 1. Publish Task       │                       │
         ├──────────────────────►│                       │
         │                       │ 2. Consume Task       │
         │                       ├──────────────────────►│
         │                       │                       │
         │                       │ 3. Execute Step       │
         │                       │                   ┌───┤
         │                       │                   │   │
         │                       │                   └──►│
         │                       │                       │
         │ 4. Report Result      │                       │
         │◄──────────────────────┼───────────────────────┤
         │                       │                       │
         │ 5. Health Check       │                       │
         │◄─────────────────────────────────────────────►│
```

## Health Check System

### Service Health Monitoring
```yaml
health_check:
  service_level:
    endpoint: "/health"
    interval: 30s
    timeout: 5s
    retries: 3
    failure_threshold: 3
  
  logic_level:
    endpoint: "/health/workflow"
    interval: 60s
    timeout: 10s
    custom_checks:
      - database_connectivity
      - external_api_availability
      - resource_utilization
```

### Health Check Implementation
- **Workflow Service Health Checks:**
  - Executor service availability
  - Message queue connectivity
  - Database connectivity
  - Resource utilization

- **Executor Service Health Checks:**
  - Business logic component health
  - External dependency availability
  - Step execution capacity
  - Resource constraints

## Timeout Management

### Timeout Hierarchy
```yaml
timeouts:
  global:
    workflow_max_duration: 1h
    step_max_duration: 10m
    health_check_timeout: 5s
  
  workflow_level:
    payment_processing: 30m
    data_migration: 2h
  
  step_level:
    api_call: 30s
    database_operation: 60s
    file_processing: 5m
```

### Timeout Patterns
- **Cascading Timeouts**: Parent timeouts encompass child timeouts
- **Graceful Shutdown**: Proper cleanup on timeout expiration
- **Retry Logic**: Configurable retry attempts with exponential backoff
- **Circuit Breaker**: Prevent cascading failures

## Timeout & Resilience Configuration

### Timeout Actions

Define specific behaviors when timeouts occur:

```yaml
# In workflow YAML definition
workflow:
  name: user_onboarding
  timeout: 300s
  timeout_action: cancel_and_cleanup  # cancel_and_cleanup, retry, fallback
  
steps:
  - name: validate_user
    timeout: 30s
    timeout_action: retry
    retry_config:
      max_attempts: 3
      backoff_strategy: exponential
      initial_delay: 1s
      max_delay: 10s
      
  - name: external_api_call
    timeout: 15s
    timeout_action: fallback
    fallback_step: default_response
    circuit_breaker:
      failure_threshold: 5
      recovery_timeout: 60s
      half_open_max_calls: 3
```

### Circuit Breaker Configuration

```go
type CircuitBreakerConfig struct {
    FailureThreshold   int           `yaml:"failure_threshold"`   // Number of failures before opening
    RecoveryTimeout    time.Duration `yaml:"recovery_timeout"`    // Time before attempting recovery
    HalfOpenMaxCalls   int           `yaml:"half_open_max_calls"` // Max calls in half-open state
    FailureRateThreshold float64     `yaml:"failure_rate_threshold"` // Percentage threshold
}

type RetryConfig struct {
    MaxAttempts      int           `yaml:"max_attempts"`
    BackoffStrategy  string        `yaml:"backoff_strategy"`  // linear, exponential, fixed
    InitialDelay     time.Duration `yaml:"initial_delay"`
    MaxDelay         time.Duration `yaml:"max_delay"`
    Jitter           bool          `yaml:"jitter"`           // Add randomness to prevent thundering herd
}

type TimeoutAction string

const (
    TimeoutActionCancel   TimeoutAction = "cancel_and_cleanup"
    TimeoutActionRetry    TimeoutAction = "retry"
    TimeoutActionFallback TimeoutAction = "fallback"
    TimeoutActionIgnore   TimeoutAction = "ignore"
)
```

## YAML Workflow Definition

### Workflow Schema
```yaml
apiVersion: workflow.magic-flow.io/v2
kind: Workflow
metadata:
  name: payment-processing
  version: v1.0.0
  description: "Process payment transactions"

spec:
  timeout: 30m
  retry_policy:
    max_attempts: 3
    backoff: exponential
  
  error_handling:
    on_failure: rollback
    notification:
      - email: admin@company.com
      - slack: "#alerts"
  
  steps:
    - name: validate-payment
      type: function
      timeout: 30s
      executor: payment-service
      function: validatePaymentRequest
      retry:
        max_attempts: 2
      on_success: process-payment
      on_failure: notify-failure
    
    - name: process-payment
      type: parallel
      timeout: 5m
      branches:
        - name: charge-card
          executor: payment-gateway
          function: chargeCard
        - name: update-inventory
          executor: inventory-service
          function: reserveItems
      on_success: confirm-order
      on_failure: rollback-payment
    
    - name: confirm-order
      type: function
      executor: order-service
      function: confirmOrder
      timeout: 1m
```

### Context Definition
```yaml
apiVersion: workflow.magic-flow.io/v2
kind: WorkflowContext
metadata:
  name: payment-context

spec:
  input_schema:
    type: object
    properties:
      user_id:
        type: string
        required: true
      amount:
        type: number
        minimum: 0.01
      currency:
        type: string
        enum: ["USD", "EUR", "GBP"]
      payment_method:
        type: object
        properties:
          type:
            type: string
            enum: ["card", "bank_transfer"]
          details:
            type: object
  
  output_schema:
    type: object
    properties:
      transaction_id:
        type: string
      status:
        type: string
        enum: ["success", "failed", "pending"]
      timestamp:
        type: string
        format: date-time
```

## Code Generation

### Generator Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   YAML Files    │    │   Code          │    │   Generated     │
│                 │    │   Generator     │    │   Code          │
│ • workflow.yaml │───►│                 │───►│ • server.go     │
│ • context.yaml  │    │ • Parser        │    │ • client.go     │
│ • config.yaml   │    │ • Validator     │    │ • types.go      │
└─────────────────┘    │ • Templates     │    │ • handlers.go   │
                       └─────────────────┘    └─────────────────┘
```

### Generated Components

#### Strongly-Typed Context Structures
```go
// Generated from context.yaml
type UserOnboardingInput struct {
    Email       string            `json:"email" validate:"required,email"`
    Name        string            `json:"name" validate:"required"`
    Preferences UserPreferences   `json:"preferences"`
}

type UserPreferences struct {
    Newsletter    bool `json:"newsletter"`
    Notifications bool `json:"notifications"`
}

type ValidateUserOutput struct {
    Validated bool                  `json:"validated"`
    UserData  *UserOnboardingInput  `json:"user_data"`
}

type CreateAccountOutput struct {
    UserID    string    `json:"user_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Generated Step Constants
```go
// Generated from workflow.yaml
type StepConstant string

const (
    StepValidateUser     StepConstant = "validate-user"
    StepCreateAccount    StepConstant = "create-account"
    StepSendWelcomeEmail StepConstant = "send-welcome-email"
    StepNotifyValidationError StepConstant = "notify-validation-error"
    StepCleanupPartialData    StepConstant = "cleanup-partial-data"
)

// Workflow initialization event from business logic
type WorkflowInitEvent struct {
    WorkflowName string      `json:"workflow_name"`
    InitialData  interface{} `json:"initial_data"`
    Priority     int         `json:"priority,omitempty"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    Source       string      `json:"source"` // Business service identifier
    Timestamp    time.Time   `json:"timestamp"`
}

// Task message for workflow execution
type TaskMessage struct {
    WorkflowID       string              `json:"workflow_id"`
    StepID           StepConstant        `json:"step_id"`
    Function         StepConstant        `json:"function"`
    Context          interface{}         `json:"context"`
    Timeout          time.Duration       `json:"timeout"`
    TimeoutAction    TimeoutAction       `json:"timeout_action"`
    RetryCount       int                `json:"retry_count"`
    MaxRetries       int                `json:"max_retries"`
    RetryConfig      *RetryConfig        `json:"retry_config,omitempty"`
    CircuitBreaker   *CircuitBreakerConfig `json:"circuit_breaker,omitempty"`
    FallbackStep     *StepConstant       `json:"fallback_step,omitempty"`
    Priority         int                `json:"priority"`
}

// Result message from step execution
type ResultMessage struct {
    WorkflowID string        `json:"workflow_id"`
    StepID     StepConstant  `json:"step_id"`
    Success    bool         `json:"success"`
    Result     interface{}   `json:"result"`
    Error      string       `json:"error,omitempty"`
    Duration   time.Duration `json:"duration"`
}

// Workflow state storage (compatible with v1 patterns)
type WorkflowState struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Status      string                 `json:"status"`
    CurrentStep StepConstant           `json:"current_step"`
    Context     interface{}            `json:"context"`
    Steps       map[StepConstant]*StepState `json:"steps"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    Metadata    map[string]string      `json:"metadata"`
}

type StepState struct {
    Status           string        `json:"status"`
    Result           interface{}   `json:"result,omitempty"`
    Error            string        `json:"error,omitempty"`
    StartedAt        *time.Time    `json:"started_at,omitempty"`
    EndedAt          *time.Time    `json:"ended_at,omitempty"`
    Duration         time.Duration `json:"duration,omitempty"`
    RetryAttempts    int          `json:"retry_attempts"`
    CircuitBreakerState string    `json:"circuit_breaker_state,omitempty"` // closed, open, half-open
    TimeoutOccurred  bool         `json:"timeout_occurred"`
    FallbackUsed     bool         `json:"fallback_used"`
    LastFailureTime  *time.Time   `json:"last_failure_time,omitempty"`
}

// Step function signatures
type ValidateUserFunc func(ctx context.Context, input *UserOnboardingInput) (*ValidateUserOutput, error)
type CreateAccountFunc func(ctx context.Context, input *ValidateUserOutput) (*CreateAccountOutput, error)
type SendWelcomeEmailFunc func(ctx context.Context, input *CreateAccountOutput) error
```

#### Workflow Service Code
- **Server Implementation**: HTTP/gRPC server with workflow endpoints
- **State Management**: Workflow state persistence and retrieval
- **Message Queue Integration**: Producer implementation with strongly-typed messages
- **Health Check Handlers**: Service and logic health endpoints
- **Monitoring**: Metrics and logging integration
- **Type Definitions**: Generated step constants and context types
- **Input Validation**: Automatic validation based on context schema

#### Workflow Executor Code
- **Consumer Implementation**: Message queue consumer with step execution
- **Business Logic Integration**: Generated interfaces for custom logic with type safety
- **Health Check Providers**: Step-level health monitoring
- **Error Handling**: Timeout and retry logic implementation
- **Result Reporting**: Structured result communication with type validation
- **Step Constants**: Generated constants for step identification and validation
- **Function Registry**: Type-safe function registration and validation

### Code Generation Command
```bash
# Generate workflow service with type-safe messaging
magic-flow generate server \
  --workflow workflow.yaml \
  --context context.yaml \
  --output ./generated/server \
  --type-safe

# Generate workflow executor with step constants
magic-flow generate executor \
  --workflow workflow.yaml \
  --context context.yaml \
  --output ./generated/executor \
  --type-safe

# Generate shared types and constants
magic-flow generate types \
  --workflow workflow.yaml \
  --context context.yaml \
  --output ./generated/types
```

## Package Structure

```
pkg_v2/
├── generator/
│   ├── parser/          # YAML parsing and validation
│   ├── templates/       # Code generation templates
│   ├── validator/       # Schema validation
│   └── cmd/            # CLI commands
├── runtime/
│   ├── server/         # Workflow service runtime
│   ├── executor/       # Workflow executor runtime
│   ├── messaging/      # Message queue abstraction
│   └── health/         # Health check framework
├── core/
│   ├── types/          # Core type definitions and generated constants
│   ├── interfaces/     # Runtime interfaces with type safety
│   ├── messaging/      # Strongly-typed message definitions
│   └── config/         # Configuration management
├── monitoring/
│   ├── metrics/        # Metrics collection
│   ├── logging/        # Structured logging
│   └── tracing/        # Distributed tracing
└── examples/
    ├── simple/         # Basic workflow example
    ├── complex/        # Advanced workflow patterns
    └── integration/    # Integration examples
```

## Migration from v1

### Compatibility Layer
- **API Compatibility**: v1 APIs wrapped in v2 runtime
- **Gradual Migration**: Step-by-step migration path
- **Feature Parity**: All v1 features available in v2
- **Configuration Migration**: Tools to convert v1 configs to v2 YAML

### Migration Steps
1. **Assessment**: Analyze existing v1 workflows
2. **YAML Conversion**: Convert v1 configurations to v2 YAML
3. **Code Generation**: Generate v2 service and executor code
4. **Testing**: Validate functionality with existing test suites
5. **Deployment**: Gradual rollout with monitoring
6. **Optimization**: Performance tuning and optimization

## Deployment Considerations

### Infrastructure Requirements
- **Message Queue**: Redis, RabbitMQ, or Apache Kafka
- **Service Discovery**: Consul, etcd, or Kubernetes DNS
- **Load Balancing**: HAProxy, NGINX, or cloud load balancers
- **Monitoring**: Prometheus, Grafana, and alerting systems

### Scalability Patterns
- **Horizontal Scaling**: Multiple executor instances per service type
- **Auto-scaling**: Dynamic scaling based on queue depth and CPU usage
- **Circuit Breakers**: Prevent cascade failures across services
- **Bulkhead Pattern**: Isolate critical workflows from non-critical ones

## Security Considerations

### Authentication & Authorization
- **Service-to-Service**: mTLS or JWT-based authentication
- **API Security**: OAuth 2.0 or API key authentication
- **Message Queue Security**: Encrypted communication and access controls

### Data Protection
- **Encryption**: At-rest and in-transit encryption
- **Secrets Management**: Integration with vault systems
- **Audit Logging**: Comprehensive audit trails

## Future Enhancements

### Planned Features
- **Visual Workflow Designer**: Web-based workflow design tool
- **Advanced Scheduling**: Cron-based and event-driven scheduling
- **Workflow Versioning**: Blue-green deployments for workflows
- **Multi-tenancy**: Isolated workflow execution per tenant
- **Workflow Analytics**: Performance metrics and optimization insights

### Test Plan & TDD Implementation

#### Testing Strategy Overview

Magic Flow v2 follows Test-Driven Development (TDD) principles to ensure high code quality, maintainability, and reliability. The testing strategy covers all layers of the system with comprehensive unit, integration, and end-to-end tests.

#### TDD Workflow

1. **Red Phase**: Write failing tests first
2. **Green Phase**: Write minimal code to make tests pass
3. **Refactor Phase**: Improve code while keeping tests green

#### Test Categories

##### 1. Unit Tests

**Core Engine Tests**
```go
// Test workflow execution engine
func TestWorkflowEngine_Execute(t *testing.T) {
    tests := []struct {
        name     string
        workflow *Workflow
        input    interface{}
        want     *WorkflowResult
        wantErr  bool
    }{
        {
            name: "successful_linear_workflow",
            workflow: &Workflow{
                Steps: []Step{
                    {ID: StepValidateUser, Function: "validateUser"},
                    {ID: StepCreateAccount, Function: "createAccount"},
                },
            },
            input: &UserOnboardingInput{Email: "test@example.com"},
            want: &WorkflowResult{Status: "completed"},
            wantErr: false,
        },
        {
            name: "workflow_with_timeout",
            workflow: &Workflow{
                TimeoutAction: TimeoutAction{Action: "cancel_and_cleanup", TimeoutMs: 5000},
                Steps: []Step{{ID: StepValidateUser, Function: "slowValidation"}},
            },
            input: &UserOnboardingInput{Email: "test@example.com"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := NewWorkflowEngine()
            got, err := engine.Execute(tt.workflow, tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Execute() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Circuit Breaker Tests**
```go
func TestCircuitBreaker_StateTransitions(t *testing.T) {
    cb := NewCircuitBreaker(&CircuitBreakerConfig{
        FailureThreshold: 3,
        RecoveryTimeoutMs: 1000,
    })
    
    // Test Closed -> Open transition
    for i := 0; i < 3; i++ {
        cb.RecordFailure()
    }
    assert.Equal(t, CircuitBreakerStateOpen, cb.GetState())
    
    // Test Open -> Half-Open transition after timeout
    time.Sleep(1100 * time.Millisecond)
    assert.Equal(t, CircuitBreakerStateHalfOpen, cb.GetState())
}
```

**Retry Mechanism Tests**
```go
func TestRetryConfig_CalculateDelay(t *testing.T) {
    tests := []struct {
        name     string
        config   RetryConfig
        attempt  int
        expected time.Duration
    }{
        {
            name: "exponential_backoff",
            config: RetryConfig{
                Strategy: "exponential",
                BaseDelayMs: 100,
                MaxDelayMs: 5000,
            },
            attempt: 3,
            expected: 800 * time.Millisecond, // 100 * 2^3
        },
        {
            name: "linear_backoff",
            config: RetryConfig{
                Strategy: "linear",
                BaseDelayMs: 100,
            },
            attempt: 3,
            expected: 300 * time.Millisecond, // 100 * 3
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            delay := tt.config.CalculateDelay(tt.attempt)
            assert.Equal(t, tt.expected, delay)
        })
    }
}
```

##### 2. Integration Tests

**Message Broker Integration**
```go
func TestKafkaMessageBroker_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    broker := NewKafkaMessageBroker(&KafkaConfig{
        Brokers: []string{"localhost:9092"},
        Topic: "workflow-events-test",
    })
    
    event := &WorkflowInitEvent{
        WorkflowID: "test-workflow-123",
        WorkflowType: "user-onboarding",
        Input: &UserOnboardingInput{Email: "test@example.com"},
    }
    
    // Test publish
    err := broker.Publish(event)
    assert.NoError(t, err)
    
    // Test consume
    received := make(chan *WorkflowInitEvent, 1)
    go broker.Subscribe(func(e *WorkflowInitEvent) {
        received <- e
    })
    
    select {
    case receivedEvent := <-received:
        assert.Equal(t, event.WorkflowID, receivedEvent.WorkflowID)
    case <-time.After(5 * time.Second):
        t.Fatal("Timeout waiting for message")
    }
}
```

**State Storage Integration**
```go
func TestRedisStateStorage_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    storage := NewRedisStateStorage(&RedisConfig{
        Address: "localhost:6379",
        DB: 1, // Use test database
    })
    
    state := &WorkflowState{
        WorkflowID: "test-workflow-123",
        Status: "running",
        Steps: map[string]*StepState{
            "validate-user": {
                Status: "completed",
                Output: &ValidateUserOutput{Valid: true},
            },
        },
    }
    
    // Test save
    err := storage.SaveWorkflowState(state)
    assert.NoError(t, err)
    
    // Test load
    loaded, err := storage.LoadWorkflowState("test-workflow-123")
    assert.NoError(t, err)
    assert.Equal(t, state.Status, loaded.Status)
    
    // Test cleanup
    err = storage.DeleteWorkflowState("test-workflow-123")
    assert.NoError(t, err)
}
```

##### 3. End-to-End Tests

**Complete Workflow Execution**
```go
func TestWorkflowService_E2E(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test")
    }
    
    // Setup test environment
    service := setupTestWorkflowService(t)
    defer teardownTestWorkflowService(t, service)
    
    // Define test workflow
    workflow := &Workflow{
        ID: "user-onboarding-e2e",
        Steps: []Step{
            {
                ID: StepValidateUser,
                Function: "validateUser",
                TimeoutAction: TimeoutAction{Action: "retry", TimeoutMs: 5000},
                RetryConfig: &RetryConfig{
                    MaxAttempts: 3,
                    Strategy: "exponential",
                    BaseDelayMs: 100,
                },
            },
            {
                ID: StepCreateAccount,
                Function: "createAccount",
                CircuitBreaker: &CircuitBreakerConfig{
                    FailureThreshold: 2,
                    RecoveryTimeoutMs: 1000,
                },
            },
            {
                ID: StepSendWelcomeEmail,
                Function: "sendWelcomeEmail",
                TimeoutAction: TimeoutAction{Action: "fallback", FallbackStep: "logWelcomeFailure"},
            },
        },
    }
    
    // Register workflow
    err := service.RegisterWorkflow(workflow)
    assert.NoError(t, err)
    
    // Trigger workflow via message broker
    event := &WorkflowInitEvent{
        WorkflowID: "test-run-" + uuid.New().String(),
        WorkflowType: "user-onboarding-e2e",
        Input: &UserOnboardingInput{
            Email: "e2e-test@example.com",
            Name: "E2E Test User",
        },
    }
    
    err = service.messageBroker.Publish(event)
    assert.NoError(t, err)
    
    // Wait for completion and verify results
    result := waitForWorkflowCompletion(t, service, event.WorkflowID, 30*time.Second)
    assert.Equal(t, "completed", result.Status)
    assert.Len(t, result.Steps, 3)
    
    // Verify each step completed successfully
    for stepID, stepState := range result.Steps {
        assert.Equal(t, "completed", stepState.Status, "Step %s should be completed", stepID)
        assert.Zero(t, stepState.RetryAttempts, "Step %s should not require retries", stepID)
    }
}
```

#### Test Infrastructure

##### Test Utilities
```go
// Test helpers for mocking and setup
type MockBusinessLogicRegistry struct {
    functions map[StepConstant]BusinessLogicFunction
}

func (m *MockBusinessLogicRegistry) Register(step StepConstant, fn BusinessLogicFunction) {
    m.functions[step] = fn
}

func (m *MockBusinessLogicRegistry) Get(step StepConstant) (BusinessLogicFunction, bool) {
    fn, exists := m.functions[step]
    return fn, exists
}

// Test workflow factory
func CreateTestWorkflow(steps ...Step) *Workflow {
    return &Workflow{
        ID: "test-workflow-" + uuid.New().String(),
        Steps: steps,
        TimeoutAction: TimeoutAction{Action: "cancel_and_cleanup", TimeoutMs: 10000},
    }
}

// Test data builders
func NewUserOnboardingInputBuilder() *UserOnboardingInputBuilder {
    return &UserOnboardingInputBuilder{
        input: &UserOnboardingInput{},
    }
}

func (b *UserOnboardingInputBuilder) WithEmail(email string) *UserOnboardingInputBuilder {
    b.input.Email = email
    return b
}

func (b *UserOnboardingInputBuilder) Build() *UserOnboardingInput {
    return b.input
}
```

##### Test Configuration
```yaml
# test-config.yaml
testing:
  message_broker:
    type: "redis"
    redis:
      address: "localhost:6379"
      db: 15  # Use separate DB for tests
  
  state_storage:
    type: "redis"
    redis:
      address: "localhost:6379"
      db: 14  # Use separate DB for tests
      ttl: "1h"  # Shorter TTL for tests
  
  timeouts:
    default_step_timeout: 5000
    default_workflow_timeout: 30000
  
  circuit_breaker:
    default_failure_threshold: 2
    default_recovery_timeout: 1000
```

#### TDD Implementation Phases

##### Phase 1: Core Engine TDD
1. **Write Tests First**:
   - Workflow execution tests
   - Step execution tests
   - Context management tests
   - Error handling tests

2. **Implement Core Engine**:
   - Basic workflow engine
   - Step executor
   - Context manager
   - Error handlers

3. **Refactor**:
   - Optimize performance
   - Improve code structure
   - Add logging and metrics

##### Phase 2: Resilience Features TDD
1. **Write Tests First**:
   - Timeout management tests
   - Circuit breaker tests
   - Retry mechanism tests
   - Fallback strategy tests

2. **Implement Resilience**:
   - Timeout handlers
   - Circuit breaker implementation
   - Retry logic
   - Fallback mechanisms

3. **Refactor**:
   - Optimize resilience patterns
   - Add configuration flexibility
   - Improve error reporting

##### Phase 3: Integration TDD
1. **Write Tests First**:
   - Message broker integration tests
   - State storage integration tests
   - End-to-end workflow tests

2. **Implement Integration**:
   - Message broker adapters
   - State storage implementations
   - Workflow service orchestration

3. **Refactor**:
   - Optimize integrations
   - Add monitoring
   - Improve error handling

#### Test Coverage Requirements

- **Unit Tests**: Minimum 90% code coverage
- **Integration Tests**: Cover all external dependencies
- **E2E Tests**: Cover critical user workflows
- **Performance Tests**: Validate timeout and throughput requirements
- **Chaos Tests**: Test resilience under failure conditions

#### Continuous Testing

```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
  
  integration-tests:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7
        ports:
          - 6379:6379
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: go test -v -tags=integration ./...
  
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: docker-compose -f docker-compose.test.yml up -d
      - run: go test -v -tags=e2e ./...
      - run: docker-compose -f docker-compose.test.yml down
```

### Integration Roadmap
- **Cloud Providers**: AWS Step Functions, Azure Logic Apps integration
- **Kubernetes**: Native Kubernetes operator
- **Observability**: OpenTelemetry integration
- **CI/CD**: GitOps workflow deployment patterns

This architecture provides a solid foundation for building scalable, resilient, and maintainable distributed workflow systems while maintaining compatibility with the existing v1 implementation.