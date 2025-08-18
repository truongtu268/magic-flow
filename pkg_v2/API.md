# Magic Flow v2 - API Documentation

## Overview

Magic Flow v2 provides a comprehensive REST API for managing workflows, executing them, monitoring performance, generating code, and handling versioning. All APIs follow RESTful conventions and return JSON responses.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints require authentication using Bearer tokens:

```bash
Authorization: Bearer your-api-token
```

## API Endpoints

### 1. Workflow Management API

#### Create Workflow

```http
POST /workflows
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "name": "order_processing",
  "version": "1.0.0",
  "description": "E-commerce order processing workflow",
  "yaml_content": "name: order_processing\nversion: 1.0.0\n...",
  "metadata": {
    "tags": ["ecommerce", "orders"],
    "owner": "team@company.com",
    "environment": "production"
  }
}
```

**Response:**
```json
{
  "id": "wf-uuid-123",
  "name": "order_processing",
  "version": "1.0.0",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "validation_status": "valid",
  "endpoints": {
    "execute": "/api/v1/workflows/order_processing/execute",
    "metrics": "/api/v1/workflows/order_processing/metrics",
    "code_generation": "/api/v1/workflows/order_processing/generate"
  }
}
```

#### List Workflows

```http
GET /workflows?page=1&limit=20&tags=ecommerce&status=active
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "workflows": [
    {
      "id": "wf-uuid-123",
      "name": "order_processing",
      "version": "1.0.0",
      "description": "E-commerce order processing workflow",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "last_execution": "2024-01-15T14:25:00Z",
      "execution_count": 1247,
      "success_rate": 0.987,
      "avg_duration": "45.2s"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

#### Get Workflow Details

```http
GET /workflows/{workflow_id}
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "id": "wf-uuid-123",
  "name": "order_processing",
  "version": "1.0.0",
  "description": "E-commerce order processing workflow",
  "yaml_content": "name: order_processing\nversion: 1.0.0\n...",
  "status": "active",
  "metadata": {
    "tags": ["ecommerce", "orders"],
    "owner": "team@company.com",
    "environment": "production"
  },
  "schema": {
    "input_schema": {...},
    "output_schema": {...}
  },
  "steps": [
    {
      "name": "validate_order",
      "type": "service_call",
      "service": "order_service",
      "timeout": "30s"
    }
  ],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Update Workflow

```http
PUT /workflows/{workflow_id}
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "description": "Updated description",
  "yaml_content": "updated YAML content",
  "metadata": {
    "tags": ["ecommerce", "orders", "v2"]
  }
}
```

#### Delete Workflow

```http
DELETE /workflows/{workflow_id}
Authorization: Bearer your-api-token
```

### 2. Workflow Execution API

#### Execute Workflow

```http
POST /workflows/{workflow_id}/execute
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "input": {
    "order_id": "order-123",
    "customer_id": "customer-456",
    "items": [
      {
        "product_id": "product-789",
        "quantity": 2,
        "price": 29.99
      }
    ]
  },
  "metadata": {
    "source": "web",
    "user_agent": "Mozilla/5.0...",
    "correlation_id": "corr-123"
  },
  "options": {
    "async": true,
    "timeout": "300s",
    "priority": "high"
  }
}
```

**Response:**
```json
{
  "execution_id": "exec-uuid-123",
  "workflow_id": "wf-uuid-123",
  "workflow_name": "order_processing",
  "status": "running",
  "created_at": "2024-01-15T10:30:00Z",
  "estimated_completion": "2024-01-15T10:32:00Z",
  "progress": {
    "current_step": "validate_order",
    "completed_steps": 0,
    "total_steps": 4,
    "percentage": 0
  },
  "tracking_url": "http://localhost:9090/executions/exec-uuid-123",
  "stream_url": "/api/v1/executions/exec-uuid-123/stream"
}
```

#### Get Execution Status

```http
GET /executions/{execution_id}
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "execution_id": "exec-uuid-123",
  "workflow_id": "wf-uuid-123",
  "workflow_name": "order_processing",
  "status": "completed",
  "created_at": "2024-01-15T10:30:00Z",
  "started_at": "2024-01-15T10:30:05Z",
  "completed_at": "2024-01-15T10:31:42Z",
  "duration": "97.3s",
  "progress": {
    "current_step": "completed",
    "completed_steps": 4,
    "total_steps": 4,
    "percentage": 100
  },
  "input": {...},
  "output": {
    "order_status": "confirmed",
    "tracking_number": "TRK-789456123",
    "estimated_delivery": "2024-01-17T18:00:00Z"
  },
  "steps": [
    {
      "name": "validate_order",
      "status": "completed",
      "started_at": "2024-01-15T10:30:05Z",
      "completed_at": "2024-01-15T10:30:15Z",
      "duration": "10.2s",
      "output": {...}
    }
  ],
  "metadata": {...}
}
```

#### Get Execution Results

```http
GET /executions/{execution_id}/results
Authorization: Bearer your-api-token
```

#### Stream Execution Events

```http
GET /executions/{execution_id}/stream
Authorization: Bearer your-api-token
Accept: text/event-stream
```

**Response (Server-Sent Events):**
```
data: {"event": "step_started", "step": "validate_order", "timestamp": "2024-01-15T10:30:05Z"}

data: {"event": "step_completed", "step": "validate_order", "status": "success", "timestamp": "2024-01-15T10:30:15Z"}

data: {"event": "workflow_completed", "status": "success", "timestamp": "2024-01-15T10:31:42Z"}
```

#### List Executions

```http
GET /executions?workflow_id=wf-uuid-123&status=completed&page=1&limit=20
Authorization: Bearer your-api-token
```

#### Cancel Execution

```http
POST /executions/{execution_id}/cancel
Authorization: Bearer your-api-token
```

### 3. Metrics and Monitoring API

#### Get Workflow Metrics

```http
GET /workflows/{workflow_id}/metrics?period=24h&granularity=1h
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "workflow_id": "wf-uuid-123",
  "period": "24h",
  "metrics": {
    "execution_count": 1247,
    "success_count": 1230,
    "failure_count": 17,
    "success_rate": 0.987,
    "avg_duration": "45.2s",
    "p50_duration": "42.1s",
    "p95_duration": "78.5s",
    "p99_duration": "125.3s",
    "throughput_per_hour": 52
  },
  "time_series": [
    {
      "timestamp": "2024-01-15T10:00:00Z",
      "execution_count": 45,
      "success_rate": 0.978,
      "avg_duration": "43.2s"
    }
  ],
  "step_metrics": {
    "validate_order": {
      "avg_duration": "10.2s",
      "success_rate": 0.995,
      "error_rate": 0.005
    }
  }
}
```

#### Get System Metrics

```http
GET /metrics/system
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "system": {
    "cpu_usage": 0.45,
    "memory_usage": 0.67,
    "disk_usage": 0.23,
    "network_io": {
      "bytes_in": 1024000,
      "bytes_out": 2048000
    }
  },
  "application": {
    "active_executions": 23,
    "queued_executions": 5,
    "total_workflows": 45,
    "cache_hit_rate": 0.89,
    "database_connections": 15
  },
  "performance": {
    "requests_per_second": 125,
    "avg_response_time": "45ms",
    "error_rate": 0.002
  }
}
```

#### Get Custom Business Metrics

```http
GET /metrics/business?metric=order_value&period=7d&group_by=customer_segment
Authorization: Bearer your-api-token
```

### 4. Code Generation API

#### Generate Client Code

```http
POST /workflows/{workflow_id}/generate
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "languages": ["go", "typescript", "python", "java"],
  "options": {
    "include_tests": true,
    "include_docs": true,
    "package_name": "order-processing-client",
    "namespace": "com.company.workflows"
  },
  "templates": {
    "go": "enterprise",
    "typescript": "react",
    "python": "fastapi",
    "java": "spring-boot"
  }
}
```

**Response:**
```json
{
  "generation_id": "gen-uuid-123",
  "status": "generating",
  "languages": ["go", "typescript", "python", "java"],
  "estimated_completion": "2024-01-15T10:32:00Z",
  "download_urls": {
    "go": "/api/v1/code-generation/gen-uuid-123/download/go",
    "typescript": "/api/v1/code-generation/gen-uuid-123/download/typescript",
    "python": "/api/v1/code-generation/gen-uuid-123/download/python",
    "java": "/api/v1/code-generation/gen-uuid-123/download/java"
  }
}
```

#### Get Code Generation Status

```http
GET /code-generation/{generation_id}
Authorization: Bearer your-api-token
```

#### Download Generated Code

```http
GET /code-generation/{generation_id}/download/{language}
Authorization: Bearer your-api-token
```

#### List Available Templates

```http
GET /code-generation/templates?language=go
Authorization: Bearer your-api-token
```

### 5. Version Management API

#### Create New Version

```http
POST /workflows/{workflow_id}/versions
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "version": "1.1.0",
  "description": "Added error handling improvements",
  "yaml_content": "updated YAML content",
  "migration_strategy": "gradual",
  "rollback_enabled": true
}
```

#### List Workflow Versions

```http
GET /workflows/{workflow_id}/versions
Authorization: Bearer your-api-token
```

**Response:**
```json
{
  "workflow_id": "wf-uuid-123",
  "versions": [
    {
      "version": "1.1.0",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "description": "Added error handling improvements",
      "execution_count": 245,
      "success_rate": 0.992
    },
    {
      "version": "1.0.0",
      "status": "deprecated",
      "created_at": "2024-01-10T10:30:00Z",
      "description": "Initial version",
      "execution_count": 1002,
      "success_rate": 0.987
    }
  ]
}
```

#### Rollback to Previous Version

```http
POST /workflows/{workflow_id}/versions/{version}/rollback
Authorization: Bearer your-api-token
```

#### Compare Versions

```http
GET /workflows/{workflow_id}/versions/compare?from=1.0.0&to=1.1.0
Authorization: Bearer your-api-token
```

## Error Handling

All API endpoints return standard HTTP status codes and JSON error responses:

```json
{
  "error": {
    "code": "WORKFLOW_NOT_FOUND",
    "message": "Workflow with ID 'wf-uuid-123' not found",
    "details": {
      "workflow_id": "wf-uuid-123",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    "request_id": "req-uuid-456"
  }
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate workflow name)
- `422 Unprocessable Entity`: Validation errors
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service temporarily unavailable

## Rate Limiting

API requests are rate-limited per API key:

- **Standard**: 1000 requests per hour
- **Premium**: 10000 requests per hour
- **Enterprise**: Unlimited

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642248000
```

## Webhooks

Magic Flow v2 supports webhooks for real-time notifications:

### Configure Webhooks

```http
POST /webhooks
Content-Type: application/json
Authorization: Bearer your-api-token

{
  "url": "https://your-app.com/webhooks/magicflow",
  "events": ["workflow.execution.completed", "workflow.execution.failed"],
  "secret": "your-webhook-secret",
  "active": true
}
```

### Webhook Events

- `workflow.created`
- `workflow.updated`
- `workflow.deleted`
- `workflow.execution.started`
- `workflow.execution.completed`
- `workflow.execution.failed`
- `workflow.execution.cancelled`
- `workflow.step.completed`
- `workflow.step.failed`

### Webhook Payload Example

```json
{
  "event": "workflow.execution.completed",
  "timestamp": "2024-01-15T10:31:42Z",
  "data": {
    "execution_id": "exec-uuid-123",
    "workflow_id": "wf-uuid-123",
    "workflow_name": "order_processing",
    "status": "completed",
    "duration": "97.3s",
    "input": {...},
    "output": {...}
  }
}
```

## SDK Examples

### Go SDK

```go
package main

import (
    "context"
    "github.com/magicflow/go-sdk"
)

func main() {
    client := magicflow.NewClient("your-api-token")
    
    // Execute workflow
    execution, err := client.ExecuteWorkflow(context.Background(), &magicflow.ExecuteRequest{
        WorkflowID: "order_processing",
        Input: map[string]interface{}{
            "order_id": "order-123",
            "customer_id": "customer-456",
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    // Wait for completion
    result, err := client.WaitForExecution(context.Background(), execution.ExecutionID)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Workflow completed: %+v\n", result.Output)
}
```

### TypeScript SDK

```typescript
import { MagicFlowClient } from '@magicflow/typescript-sdk';

const client = new MagicFlowClient({
  apiToken: 'your-api-token',
  baseURL: 'http://localhost:8080/api/v1'
});

// Execute workflow
const execution = await client.executeWorkflow({
  workflowId: 'order_processing',
  input: {
    order_id: 'order-123',
    customer_id: 'customer-456'
  }
});

// Stream execution events
const stream = client.streamExecution(execution.executionId);
stream.on('step_completed', (event) => {
  console.log('Step completed:', event.step);
});

stream.on('workflow_completed', (event) => {
  console.log('Workflow completed:', event.output);
});
```

### Python SDK

```python
from magicflow import MagicFlowClient

client = MagicFlowClient(api_token='your-api-token')

# Execute workflow
execution = client.execute_workflow(
    workflow_id='order_processing',
    input={
        'order_id': 'order-123',
        'customer_id': 'customer-456'
    }
)

# Wait for completion
result = client.wait_for_execution(execution.execution_id)
print(f"Workflow completed: {result.output}")
```

## OpenAPI Specification

The complete OpenAPI 3.0 specification is available at:

```
http://localhost:8080/api/docs/openapi.json
```

Interactive API documentation (Swagger UI) is available at:

```
http://localhost:8080/api/docs
```