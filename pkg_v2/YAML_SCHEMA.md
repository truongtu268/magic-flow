# Magic Flow v2 - YAML Schema Documentation

## Overview

Magic Flow v2 uses YAML files to define workflows in a declarative, human-readable format. This document provides comprehensive documentation of the YAML schema, validation rules, examples, and best practices.

## Schema Structure

### Root Schema

```yaml
# Complete workflow definition schema
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: string                    # Required: Workflow name (must be unique)
  version: string                 # Required: Semantic version (e.g., "1.0.0")
  description: string             # Optional: Human-readable description
  tags: [string]                  # Optional: Tags for categorization
  owner: string                   # Optional: Owner/team responsible
  created_by: string              # Optional: Creator information
  created_at: string              # Auto-generated: ISO 8601 timestamp
  updated_at: string              # Auto-generated: ISO 8601 timestamp
  labels:                         # Optional: Key-value labels
    key: value
  annotations:                    # Optional: Extended metadata
    key: value

spec:
  # Workflow specification
  input_schema: object            # Required: Input data schema
  output_schema: object           # Required: Output data schema
  timeout: string                 # Optional: Overall workflow timeout
  retry_policy: object            # Optional: Retry configuration
  error_handling: object          # Optional: Error handling strategy
  
  # Workflow steps
  steps: [object]                 # Required: Array of workflow steps
  
  # Advanced features
  conditions: [object]            # Optional: Conditional execution
  parallel_groups: [object]       # Optional: Parallel execution groups
  data_mapping: object            # Optional: Data transformation rules
  notifications: object           # Optional: Notification configuration
  monitoring: object              # Optional: Monitoring and metrics
```

## Metadata Section

### Basic Metadata

```yaml
metadata:
  name: "order_processing"          # Workflow identifier (kebab-case recommended)
  version: "2.1.0"                  # Semantic versioning
  description: "Process customer orders from cart to fulfillment"
  tags:
    - "ecommerce"
    - "order-management"
    - "payment"
  owner: "ecommerce-team"
  created_by: "john.doe@company.com"
```

**Validation Rules:**
- `name`: Must be 1-63 characters, lowercase alphanumeric with hyphens
- `version`: Must follow semantic versioning (MAJOR.MINOR.PATCH)
- `description`: Maximum 500 characters
- `tags`: Array of strings, each 1-50 characters
- `owner`: Valid email or team identifier

### Labels and Annotations

```yaml
metadata:
  labels:
    environment: "production"
    team: "platform"
    cost-center: "engineering"
    compliance: "pci-dss"
  
  annotations:
    documentation: "https://wiki.company.com/workflows/order-processing"
    runbook: "https://runbooks.company.com/order-processing"
    slack-channel: "#ecommerce-alerts"
    oncall-team: "ecommerce-oncall"
    sla: "99.9%"
    data-classification: "confidential"
```

**Usage:**
- **Labels**: Used for filtering, grouping, and selection
- **Annotations**: Extended metadata for documentation and tooling

## Input/Output Schema

### JSON Schema Format

```yaml
spec:
  input_schema:
    type: "object"
    required: ["customer_id", "items", "payment_method"]
    properties:
      customer_id:
        type: "string"
        pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
        description: "UUID of the customer"
      
      items:
        type: "array"
        minItems: 1
        maxItems: 100
        items:
          type: "object"
          required: ["product_id", "quantity", "price"]
          properties:
            product_id:
              type: "string"
              description: "Product identifier"
            quantity:
              type: "integer"
              minimum: 1
              maximum: 999
            price:
              type: "number"
              minimum: 0
              multipleOf: 0.01
      
      payment_method:
        type: "object"
        required: ["type"]
        properties:
          type:
            type: "string"
            enum: ["credit_card", "debit_card", "paypal", "bank_transfer"]
          card_token:
            type: "string"
            description: "Tokenized card information"
          billing_address:
            $ref: "#/definitions/Address"
      
      shipping_address:
        $ref: "#/definitions/Address"
      
      metadata:
        type: "object"
        additionalProperties: true
        description: "Additional order metadata"
    
    definitions:
      Address:
        type: "object"
        required: ["street", "city", "country", "postal_code"]
        properties:
          street:
            type: "string"
            maxLength: 200
          city:
            type: "string"
            maxLength: 100
          state:
            type: "string"
            maxLength: 100
          country:
            type: "string"
            pattern: "^[A-Z]{2}$"
            description: "ISO 3166-1 alpha-2 country code"
          postal_code:
            type: "string"
            maxLength: 20

  output_schema:
    type: "object"
    required: ["order_id", "status", "total_amount"]
    properties:
      order_id:
        type: "string"
        description: "Generated order identifier"
      
      status:
        type: "string"
        enum: ["confirmed", "processing", "shipped", "delivered", "cancelled", "failed"]
      
      total_amount:
        type: "number"
        minimum: 0
        description: "Total order amount including taxes and shipping"
      
      payment_status:
        type: "string"
        enum: ["pending", "authorized", "captured", "failed", "refunded"]
      
      tracking_number:
        type: "string"
        description: "Shipping tracking number"
      
      estimated_delivery:
        type: "string"
        format: "date-time"
        description: "Estimated delivery date"
      
      items:
        type: "array"
        items:
          type: "object"
          properties:
            product_id:
              type: "string"
            quantity_fulfilled:
              type: "integer"
            unit_price:
              type: "number"
            total_price:
              type: "number"
      
      timestamps:
        type: "object"
        properties:
          created_at:
            type: "string"
            format: "date-time"
          confirmed_at:
            type: "string"
            format: "date-time"
          shipped_at:
            type: "string"
            format: "date-time"
```

### Schema Validation

**Built-in Validators:**
- **Type validation**: string, number, integer, boolean, array, object
- **Format validation**: date-time, email, uri, uuid, etc.
- **Range validation**: minimum, maximum, minLength, maxLength
- **Pattern validation**: Regular expressions
- **Enum validation**: Predefined value lists
- **Array validation**: minItems, maxItems, uniqueItems
- **Object validation**: required fields, additionalProperties

**Custom Validators:**
```yaml
spec:
  input_schema:
    type: "object"
    properties:
      email:
        type: "string"
        format: "email"
        x-validator: "business_email"  # Custom validator
      
      phone:
        type: "string"
        x-validator: "international_phone"
      
      credit_card:
        type: "string"
        x-validator: "luhn_checksum"
    
    x-validators:
      business_email:
        function: "validateBusinessEmail"
        params:
          allowed_domains: ["company.com", "partner.com"]
      
      international_phone:
        function: "validatePhoneNumber"
        params:
          format: "E164"
      
      luhn_checksum:
        function: "validateCreditCard"
        params:
          algorithm: "luhn"
```

## Workflow Steps

### Basic Step Structure

```yaml
spec:
  steps:
    - name: "validate_order"              # Required: Step identifier
      type: "service_call"                # Required: Step type
      description: "Validate order data"   # Optional: Step description
      
      # Step configuration
      config:
        service: "order-validation-service"
        endpoint: "/api/v1/validate"
        method: "POST"
        timeout: "30s"
      
      # Input/output mapping
      input:
        order_data: "$.input"             # JSONPath expression
        validation_rules: "strict"
      
      output:
        validation_result: "$.response.valid"
        errors: "$.response.errors"
      
      # Error handling
      on_error:
        action: "retry"
        max_attempts: 3
        backoff: "exponential"
      
      # Conditions
      when: "$.input.amount > 100"         # Execute only if condition is true
      
      # Dependencies
      depends_on: []                       # Steps that must complete first
```

### Step Types

#### 1. Service Call

```yaml
- name: "call_payment_service"
  type: "service_call"
  config:
    service: "payment-gateway"
    endpoint: "/api/v1/charge"
    method: "POST"
    headers:
      Authorization: "Bearer ${secrets.payment_api_key}"
      Content-Type: "application/json"
    timeout: "45s"
    retry_policy:
      max_attempts: 3
      backoff: "exponential"
      backoff_multiplier: 2
      max_backoff: "60s"
  
  input:
    amount: "$.input.total_amount"
    currency: "USD"
    payment_method: "$.input.payment_method"
    customer_id: "$.input.customer_id"
  
  output:
    transaction_id: "$.response.transaction_id"
    status: "$.response.status"
    authorization_code: "$.response.auth_code"
```

#### 2. Database Operation

```yaml
- name: "save_order"
  type: "database"
  config:
    connection: "primary_db"
    operation: "insert"
    table: "orders"
    timeout: "10s"
  
  input:
    data:
      customer_id: "$.input.customer_id"
      total_amount: "$.steps.calculate_total.output.amount"
      status: "confirmed"
      created_at: "${now()}"
  
  output:
    order_id: "$.result.id"
    created_at: "$.result.created_at"
```

#### 3. Message Queue

```yaml
- name: "notify_fulfillment"
  type: "message_queue"
  config:
    queue: "fulfillment-queue"
    exchange: "orders"
    routing_key: "order.created"
    timeout: "5s"
  
  input:
    message:
      order_id: "$.steps.save_order.output.order_id"
      customer_id: "$.input.customer_id"
      items: "$.input.items"
      shipping_address: "$.input.shipping_address"
      priority: "normal"
  
  output:
    message_id: "$.result.message_id"
    published_at: "$.result.timestamp"
```

#### 4. Data Transformation

```yaml
- name: "calculate_totals"
  type: "transform"
  config:
    language: "javascript"  # or "python", "go"
    timeout: "10s"
  
  script: |
    function transform(input) {
      const items = input.items;
      let subtotal = 0;
      
      for (const item of items) {
        subtotal += item.quantity * item.price;
      }
      
      const tax = subtotal * 0.08;  // 8% tax
      const shipping = subtotal > 50 ? 0 : 9.99;
      const total = subtotal + tax + shipping;
      
      return {
        subtotal: Math.round(subtotal * 100) / 100,
        tax: Math.round(tax * 100) / 100,
        shipping: Math.round(shipping * 100) / 100,
        total: Math.round(total * 100) / 100
      };
    }
  
  input:
    items: "$.input.items"
  
  output:
    subtotal: "$.result.subtotal"
    tax: "$.result.tax"
    shipping: "$.result.shipping"
    total: "$.result.total"
```

#### 5. Conditional Logic

```yaml
- name: "fraud_check"
  type: "condition"
  config:
    timeout: "5s"
  
  condition: "$.input.total_amount > 1000 || $.input.customer.risk_score > 0.7"
  
  then:
    - name: "manual_review"
      type: "human_task"
      config:
        assignee: "fraud-team"
        timeout: "24h"
        form:
          - field: "approve"
            type: "boolean"
            label: "Approve this transaction?"
          - field: "notes"
            type: "text"
            label: "Review notes"
      
      output:
        approved: "$.form.approve"
        review_notes: "$.form.notes"
  
  else:
    - name: "auto_approve"
      type: "transform"
      script: |
        function transform(input) {
          return { approved: true, review_notes: "Auto-approved" };
        }
      
      output:
        approved: "$.result.approved"
        review_notes: "$.result.review_notes"
```

#### 6. Parallel Execution

```yaml
- name: "parallel_processing"
  type: "parallel"
  config:
    timeout: "60s"
    fail_fast: false  # Continue even if some branches fail
  
  branches:
    - name: "inventory_check"
      steps:
        - name: "check_stock"
          type: "service_call"
          config:
            service: "inventory-service"
            endpoint: "/api/v1/check-stock"
          input:
            items: "$.input.items"
          output:
            stock_status: "$.response.status"
    
    - name: "customer_validation"
      steps:
        - name: "validate_customer"
          type: "service_call"
          config:
            service: "customer-service"
            endpoint: "/api/v1/validate"
          input:
            customer_id: "$.input.customer_id"
          output:
            customer_valid: "$.response.valid"
    
    - name: "address_validation"
      steps:
        - name: "validate_address"
          type: "service_call"
          config:
            service: "address-service"
            endpoint: "/api/v1/validate"
          input:
            address: "$.input.shipping_address"
          output:
            address_valid: "$.response.valid"
  
  output:
    inventory_status: "$.branches.inventory_check.stock_status"
    customer_status: "$.branches.customer_validation.customer_valid"
    address_status: "$.branches.address_validation.address_valid"
```

#### 7. Loop/Iteration

```yaml
- name: "process_items"
  type: "loop"
  config:
    timeout: "300s"
    max_iterations: 100
    parallel: true  # Process items in parallel
    batch_size: 10  # Process 10 items at a time
  
  iterate_over: "$.input.items"
  item_variable: "current_item"
  
  steps:
    - name: "reserve_inventory"
      type: "service_call"
      config:
        service: "inventory-service"
        endpoint: "/api/v1/reserve"
      input:
        product_id: "$.current_item.product_id"
        quantity: "$.current_item.quantity"
      output:
        reservation_id: "$.response.reservation_id"
    
    - name: "calculate_item_total"
      type: "transform"
      script: |
        function transform(input) {
          return {
            item_total: input.current_item.quantity * input.current_item.price
          };
        }
      input:
        current_item: "$.current_item"
      output:
        item_total: "$.result.item_total"
  
  output:
    processed_items: "$.iterations[*].steps.reserve_inventory.output"
    total_reservations: "$.iterations.length"
```

### Data Mapping and Transformation

#### JSONPath Expressions

```yaml
# Basic JSONPath examples
input:
  # Root input
  customer_id: "$.input.customer_id"
  
  # Nested properties
  street_address: "$.input.shipping_address.street"
  
  # Array elements
  first_item: "$.input.items[0]"
  last_item: "$.input.items[-1]"
  
  # Array filtering
  expensive_items: "$.input.items[?(@.price > 100)]"
  
  # Array mapping
  product_ids: "$.input.items[*].product_id"
  
  # Conditional selection
  priority_customer: "$.input.customer.tier == 'premium'"
  
  # Previous step outputs
  payment_result: "$.steps.process_payment.output.status"
  
  # Built-in functions
  current_time: "${now()}"
  random_id: "${uuid()}"
  
  # Environment variables
  api_endpoint: "${env.PAYMENT_API_URL}"
  
  # Secrets
  api_key: "${secrets.payment_api_key}"
```

#### Advanced Data Mapping

```yaml
spec:
  data_mapping:
    # Global data transformations
    transforms:
      - name: "normalize_phone"
        input: "$.input.customer.phone"
        script: |
          function transform(phone) {
            return phone.replace(/[^0-9]/g, '');
          }
        output: "normalized_phone"
      
      - name: "calculate_discount"
        input:
          total: "$.steps.calculate_totals.output.total"
          customer_tier: "$.input.customer.tier"
        script: |
          function transform(input) {
            const discounts = {
              'bronze': 0.05,
              'silver': 0.10,
              'gold': 0.15,
              'platinum': 0.20
            };
            
            const discount_rate = discounts[input.customer_tier] || 0;
            const discount_amount = input.total * discount_rate;
            
            return {
              discount_rate: discount_rate,
              discount_amount: Math.round(discount_amount * 100) / 100,
              final_total: Math.round((input.total - discount_amount) * 100) / 100
            };
          }
        output: "discount_calculation"
    
    # Data validation rules
    validations:
      - name: "validate_email"
        field: "$.input.customer.email"
        rules:
          - type: "format"
            format: "email"
          - type: "custom"
            function: "validateBusinessEmail"
      
      - name: "validate_amount"
        field: "$.input.total_amount"
        rules:
          - type: "range"
            min: 0.01
            max: 10000.00
    
    # Data enrichment
    enrichments:
      - name: "customer_profile"
        source: "customer-service"
        endpoint: "/api/v1/customers/${input.customer_id}/profile"
        cache_ttl: "1h"
        output: "customer_profile"
      
      - name: "product_details"
        source: "product-service"
        endpoint: "/api/v1/products/batch"
        input:
          product_ids: "$.input.items[*].product_id"
        cache_ttl: "30m"
        output: "product_details"
```

## Error Handling and Retry Policies

### Global Error Handling

```yaml
spec:
  error_handling:
    # Default error strategy
    default_strategy: "retry_then_fail"
    
    # Global retry policy
    retry_policy:
      max_attempts: 3
      backoff: "exponential"
      backoff_multiplier: 2
      initial_backoff: "1s"
      max_backoff: "60s"
      jitter: true
    
    # Error categorization
    error_categories:
      - name: "transient_errors"
        patterns:
          - "connection_timeout"
          - "service_unavailable"
          - "rate_limit_exceeded"
        strategy: "retry"
        max_attempts: 5
      
      - name: "validation_errors"
        patterns:
          - "invalid_input"
          - "schema_validation_failed"
        strategy: "fail_immediately"
      
      - name: "business_errors"
        patterns:
          - "insufficient_funds"
          - "product_out_of_stock"
        strategy: "custom_handler"
        handler: "business_error_handler"
    
    # Compensation actions
    compensation:
      enabled: true
      steps:
        - name: "rollback_inventory"
          condition: "$.error.step == 'reserve_inventory'"
          action:
            type: "service_call"
            config:
              service: "inventory-service"
              endpoint: "/api/v1/release"
            input:
              reservation_ids: "$.context.reservation_ids"
        
        - name: "refund_payment"
          condition: "$.error.step_index > 3"  # After payment step
          action:
            type: "service_call"
            config:
              service: "payment-service"
              endpoint: "/api/v1/refund"
            input:
              transaction_id: "$.context.transaction_id"
    
    # Dead letter queue
    dead_letter_queue:
      enabled: true
      queue: "failed-workflows"
      max_retries: 3
      retention_period: "7d"
```

### Step-Level Error Handling

```yaml
steps:
  - name: "critical_payment_step"
    type: "service_call"
    config:
      service: "payment-gateway"
      endpoint: "/api/v1/charge"
    
    # Step-specific error handling
    on_error:
      # Immediate retry for specific errors
      retry:
        conditions:
          - "$.error.code == 'TIMEOUT'"
          - "$.error.code == 'RATE_LIMIT'"
        max_attempts: 5
        backoff: "linear"
        backoff_interval: "2s"
      
      # Fallback to alternative service
      fallback:
        condition: "$.error.code == 'SERVICE_UNAVAILABLE'"
        steps:
          - name: "backup_payment_service"
            type: "service_call"
            config:
              service: "backup-payment-gateway"
              endpoint: "/api/v1/process"
      
      # Circuit breaker
      circuit_breaker:
        enabled: true
        failure_threshold: 5
        timeout: "60s"
        half_open_max_calls: 3
      
      # Custom error handling
      custom_handler:
        condition: "$.error.code == 'FRAUD_DETECTED'"
        action:
          type: "human_task"
          config:
            assignee: "fraud-team"
            priority: "high"
            timeout: "2h"
```

## Advanced Features

### Conditional Execution

```yaml
spec:
  conditions:
    # Global conditions
    - name: "high_value_order"
      expression: "$.input.total_amount > 1000"
      description: "Orders over $1000 require additional verification"
    
    - name: "international_order"
      expression: "$.input.shipping_address.country != 'US'"
      description: "International orders have different processing"
    
    - name: "premium_customer"
      expression: "$.input.customer.tier in ['gold', 'platinum']"
      description: "Premium customers get expedited processing"
  
  steps:
    - name: "fraud_screening"
      type: "service_call"
      when: "conditions.high_value_order || conditions.international_order"
      config:
        service: "fraud-detection"
        endpoint: "/api/v1/screen"
    
    - name: "expedited_processing"
      type: "parallel"
      when: "conditions.premium_customer"
      branches:
        - name: "priority_fulfillment"
          steps:
            - name: "reserve_premium_inventory"
              type: "service_call"
              config:
                service: "inventory-service"
                endpoint: "/api/v1/reserve-premium"
```

### Workflow Composition

```yaml
spec:
  # Include other workflows
  includes:
    - workflow: "customer_validation"
      version: "1.2.0"
      alias: "validate_customer"
    
    - workflow: "payment_processing"
      version: "2.0.0"
      alias: "process_payment"
    
    - workflow: "inventory_management"
      version: "1.5.0"
      alias: "manage_inventory"
  
  steps:
    # Use included workflows as steps
    - name: "validate_customer_data"
      type: "workflow"
      workflow: "validate_customer"
      input:
        customer_id: "$.input.customer_id"
        validation_level: "strict"
      output:
        customer_valid: "$.result.valid"
        customer_profile: "$.result.profile"
    
    - name: "process_order_payment"
      type: "workflow"
      workflow: "process_payment"
      input:
        amount: "$.steps.calculate_totals.output.total"
        payment_method: "$.input.payment_method"
        customer_id: "$.input.customer_id"
      output:
        transaction_id: "$.result.transaction_id"
        payment_status: "$.result.status"
```

### Monitoring and Observability

```yaml
spec:
  monitoring:
    # Metrics collection
    metrics:
      - name: "order_value"
        type: "gauge"
        value: "$.steps.calculate_totals.output.total"
        labels:
          customer_tier: "$.input.customer.tier"
          order_type: "$.input.order_type"
      
      - name: "processing_duration"
        type: "histogram"
        value: "$.execution.duration"
        buckets: [1, 5, 10, 30, 60, 300]
      
      - name: "payment_success_rate"
        type: "counter"
        increment_on: "$.steps.process_payment.output.status == 'success'"
        labels:
          payment_method: "$.input.payment_method.type"
    
    # Distributed tracing
    tracing:
      enabled: true
      service_name: "order-processing-workflow"
      tags:
        workflow_version: "$.metadata.version"
        customer_id: "$.input.customer_id"
        order_type: "$.input.order_type"
    
    # Logging
    logging:
      level: "info"
      structured: true
      fields:
        workflow_name: "$.metadata.name"
        execution_id: "$.execution.id"
        customer_id: "$.input.customer_id"
        order_value: "$.steps.calculate_totals.output.total"
    
    # Health checks
    health_checks:
      - name: "payment_service_health"
        type: "http"
        url: "http://payment-service/health"
        interval: "30s"
        timeout: "5s"
      
      - name: "database_health"
        type: "database"
        connection: "primary_db"
        query: "SELECT 1"
        interval: "60s"
```

### Security and Compliance

```yaml
spec:
  security:
    # Data classification
    data_classification: "confidential"
    
    # PII handling
    pii_fields:
      - "$.input.customer.email"
      - "$.input.customer.phone"
      - "$.input.shipping_address"
      - "$.input.payment_method.card_token"
    
    # Encryption
    encryption:
      at_rest: true
      in_transit: true
      key_rotation: "90d"
    
    # Access control
    rbac:
      required_roles: ["order_processor"]
      required_permissions: ["orders:create", "payments:process"]
    
    # Audit logging
    audit:
      enabled: true
      events:
        - "workflow_started"
        - "payment_processed"
        - "order_created"
        - "workflow_completed"
        - "workflow_failed"
    
    # Compliance
    compliance:
      frameworks: ["PCI-DSS", "GDPR", "SOX"]
      data_retention: "7y"
      data_residency: "US"
```

## Validation Rules

### Schema Validation

```yaml
# Built-in validation rules
validation:
  # Required fields
  required:
    - "metadata.name"
    - "metadata.version"
    - "spec.input_schema"
    - "spec.output_schema"
    - "spec.steps"
  
  # Field constraints
  constraints:
    metadata.name:
      pattern: "^[a-z0-9-]+$"
      min_length: 1
      max_length: 63
    
    metadata.version:
      pattern: "^\d+\.\d+\.\d+$"
    
    metadata.description:
      max_length: 500
    
    metadata.tags:
      max_items: 20
      item_pattern: "^[a-z0-9-]+$"
    
    spec.steps:
      min_items: 1
      max_items: 100
    
    spec.timeout:
      pattern: "^\d+[smh]$"
      max_value: "24h"
  
  # Custom validation functions
  custom_validators:
    - name: "unique_step_names"
      function: "validateUniqueStepNames"
      message: "Step names must be unique within a workflow"
    
    - name: "valid_dependencies"
      function: "validateStepDependencies"
      message: "Step dependencies must reference existing steps"
    
    - name: "no_circular_dependencies"
      function: "validateNoCycles"
      message: "Workflow must not contain circular dependencies"
```

### Runtime Validation

```yaml
spec:
  runtime_validation:
    # Input validation
    input_validation:
      enabled: true
      strict_mode: true  # Reject additional properties
      coerce_types: false  # Don't automatically convert types
    
    # Output validation
    output_validation:
      enabled: true
      log_violations: true
      fail_on_violation: false  # Log but don't fail
    
    # Step validation
    step_validation:
      validate_service_endpoints: true
      validate_database_connections: true
      validate_message_queues: true
    
    # Performance validation
    performance_validation:
      max_execution_time: "1h"
      max_memory_usage: "2GB"
      max_cpu_usage: "80%"
```

## Best Practices

### Workflow Design

1. **Keep workflows focused and single-purpose**
   ```yaml
   # Good: Focused on order processing
   metadata:
     name: "order_processing"
     description: "Process customer orders from validation to fulfillment"
   
   # Avoid: Too broad scope
   metadata:
     name: "customer_management"
     description: "Handle all customer-related operations"
   ```

2. **Use meaningful names and descriptions**
   ```yaml
   steps:
     - name: "validate_payment_method"  # Clear and descriptive
       description: "Validate customer payment method and billing address"
     
     # Avoid generic names
     - name: "step1"  # Not descriptive
       description: "Do something"
   ```

3. **Design for idempotency**
   ```yaml
   steps:
     - name: "create_order"
       type: "service_call"
       config:
         service: "order-service"
         endpoint: "/api/v1/orders"
         idempotency_key: "${input.customer_id}-${input.cart_id}-${now('date')}"
   ```

4. **Handle errors gracefully**
   ```yaml
   steps:
     - name: "charge_payment"
       type: "service_call"
       on_error:
         retry:
           max_attempts: 3
           backoff: "exponential"
         fallback:
           steps:
             - name: "notify_payment_failure"
               type: "message_queue"
   ```

### Performance Optimization

1. **Use parallel execution when possible**
   ```yaml
   - name: "parallel_validations"
     type: "parallel"
     branches:
       - name: "validate_customer"
       - name: "validate_inventory"
       - name: "validate_address"
   ```

2. **Optimize data flow**
   ```yaml
   # Good: Pass only necessary data
   input:
     customer_id: "$.input.customer_id"
     order_total: "$.steps.calculate_total.output.amount"
   
   # Avoid: Passing entire context
   input:
     everything: "$"  # Inefficient
   ```

3. **Use caching for expensive operations**
   ```yaml
   - name: "get_customer_profile"
     type: "service_call"
     config:
       cache:
         enabled: true
         ttl: "1h"
         key: "customer-${input.customer_id}"
   ```

### Security Best Practices

1. **Never expose sensitive data in logs**
   ```yaml
   input:
     customer_id: "$.input.customer_id"
     # Don't log sensitive fields
     payment_token: "$.input.payment_method.token"  # This will be redacted
   
   security:
     pii_fields:
       - "$.input.payment_method.token"
       - "$.input.customer.ssn"
   ```

2. **Use secrets management**
   ```yaml
   config:
     headers:
       Authorization: "Bearer ${secrets.api_key}"  # From secrets store
       # Avoid hardcoding
       # Authorization: "Bearer abc123"  # Don't do this
   ```

3. **Validate all inputs**
   ```yaml
   spec:
     input_schema:
       type: "object"
       required: ["customer_id", "amount"]
       properties:
         customer_id:
           type: "string"
           pattern: "^[0-9a-f-]{36}$"  # UUID format
         amount:
           type: "number"
           minimum: 0.01
           maximum: 10000
   ```

### Maintainability

1. **Version your workflows**
   ```yaml
   metadata:
     version: "2.1.0"  # Semantic versioning
   
   # Document breaking changes
   annotations:
     changelog: "v2.1.0: Added fraud detection step"
     migration_guide: "https://docs.company.com/workflows/order-processing/migration"
   ```

2. **Use composition for reusability**
   ```yaml
   # Create reusable sub-workflows
   includes:
     - workflow: "common_validations"
       version: "1.0.0"
     - workflow: "payment_processing"
       version: "2.0.0"
   ```

3. **Document your workflows**
   ```yaml
   metadata:
     description: "Process customer orders from cart to fulfillment"
     annotations:
       documentation: "https://wiki.company.com/workflows/order-processing"
       runbook: "https://runbooks.company.com/order-processing"
       owner_team: "ecommerce"
       slack_channel: "#ecommerce-support"
   ```

## Examples

### Complete E-commerce Order Processing Workflow

```yaml
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: "ecommerce-order-processing"
  version: "2.1.0"
  description: "Complete e-commerce order processing from cart to fulfillment"
  tags: ["ecommerce", "orders", "payments", "fulfillment"]
  owner: "ecommerce-team"
  created_by: "platform-team"
  labels:
    environment: "production"
    team: "ecommerce"
    compliance: "pci-dss"
  annotations:
    documentation: "https://wiki.company.com/workflows/order-processing"
    runbook: "https://runbooks.company.com/order-processing"
    slack_channel: "#ecommerce-alerts"

spec:
  input_schema:
    type: "object"
    required: ["customer_id", "items", "payment_method", "shipping_address"]
    properties:
      customer_id:
        type: "string"
        pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
      items:
        type: "array"
        minItems: 1
        maxItems: 50
        items:
          type: "object"
          required: ["product_id", "quantity", "price"]
          properties:
            product_id: { type: "string" }
            quantity: { type: "integer", minimum: 1, maximum: 99 }
            price: { type: "number", minimum: 0 }
      payment_method:
        type: "object"
        required: ["type", "token"]
        properties:
          type: { type: "string", enum: ["credit_card", "debit_card", "paypal"] }
          token: { type: "string" }
      shipping_address:
        type: "object"
        required: ["street", "city", "state", "country", "postal_code"]
        properties:
          street: { type: "string", maxLength: 200 }
          city: { type: "string", maxLength: 100 }
          state: { type: "string", maxLength: 100 }
          country: { type: "string", pattern: "^[A-Z]{2}$" }
          postal_code: { type: "string", maxLength: 20 }

  output_schema:
    type: "object"
    required: ["order_id", "status", "total_amount"]
    properties:
      order_id: { type: "string" }
      status: { type: "string", enum: ["confirmed", "processing", "shipped", "failed"] }
      total_amount: { type: "number", minimum: 0 }
      tracking_number: { type: "string" }
      estimated_delivery: { type: "string", format: "date-time" }

  timeout: "30m"
  
  retry_policy:
    max_attempts: 3
    backoff: "exponential"
    backoff_multiplier: 2
    max_backoff: "60s"

  steps:
    # Parallel validation phase
    - name: "validation_phase"
      type: "parallel"
      config:
        timeout: "60s"
        fail_fast: true
      branches:
        - name: "customer_validation"
          steps:
            - name: "validate_customer"
              type: "service_call"
              config:
                service: "customer-service"
                endpoint: "/api/v1/customers/${input.customer_id}/validate"
                timeout: "10s"
              output:
                customer_valid: "$.response.valid"
                customer_tier: "$.response.tier"
        
        - name: "inventory_validation"
          steps:
            - name: "check_inventory"
              type: "service_call"
              config:
                service: "inventory-service"
                endpoint: "/api/v1/check-availability"
                timeout: "15s"
              input:
                items: "$.input.items"
              output:
                availability: "$.response.availability"
        
        - name: "address_validation"
          steps:
            - name: "validate_address"
              type: "service_call"
              config:
                service: "address-service"
                endpoint: "/api/v1/validate"
                timeout: "10s"
              input:
                address: "$.input.shipping_address"
              output:
                address_valid: "$.response.valid"
                normalized_address: "$.response.normalized"

    # Calculate totals
    - name: "calculate_totals"
      type: "transform"
      config:
        language: "javascript"
        timeout: "5s"
      script: |
        function transform(input) {
          const items = input.items;
          const customerTier = input.customer_tier;
          
          let subtotal = 0;
          for (const item of items) {
            subtotal += item.quantity * item.price;
          }
          
          // Calculate tax (8%)
          const tax = subtotal * 0.08;
          
          // Calculate shipping
          const shipping = subtotal > 50 ? 0 : 9.99;
          
          // Calculate discount based on customer tier
          const discountRates = {
            'bronze': 0.05,
            'silver': 0.10,
            'gold': 0.15,
            'platinum': 0.20
          };
          const discountRate = discountRates[customerTier] || 0;
          const discount = subtotal * discountRate;
          
          const total = subtotal + tax + shipping - discount;
          
          return {
            subtotal: Math.round(subtotal * 100) / 100,
            tax: Math.round(tax * 100) / 100,
            shipping: Math.round(shipping * 100) / 100,
            discount: Math.round(discount * 100) / 100,
            total: Math.round(total * 100) / 100
          };
        }
      input:
        items: "$.input.items"
        customer_tier: "$.steps.validation_phase.branches.customer_validation.customer_tier"
      output:
        subtotal: "$.result.subtotal"
        tax: "$.result.tax"
        shipping: "$.result.shipping"
        discount: "$.result.discount"
        total: "$.result.total"

    # Fraud detection for high-value orders
    - name: "fraud_detection"
      type: "service_call"
      when: "$.steps.calculate_totals.output.total > 500"
      config:
        service: "fraud-detection-service"
        endpoint: "/api/v1/screen"
        timeout: "30s"
      input:
        customer_id: "$.input.customer_id"
        order_amount: "$.steps.calculate_totals.output.total"
        payment_method: "$.input.payment_method"
        shipping_address: "$.input.shipping_address"
      output:
        risk_score: "$.response.risk_score"
        approved: "$.response.approved"
      on_error:
        action: "continue"  # Continue if fraud service is down
        default_output:
          risk_score: 0.0
          approved: true

    # Reserve inventory
    - name: "reserve_inventory"
      type: "service_call"
      config:
        service: "inventory-service"
        endpoint: "/api/v1/reserve"
        timeout: "20s"
      input:
        items: "$.input.items"
        customer_id: "$.input.customer_id"
        reservation_ttl: "15m"
      output:
        reservation_id: "$.response.reservation_id"
        reserved_items: "$.response.items"
      on_error:
        action: "fail"
        message: "Unable to reserve inventory"

    # Process payment
    - name: "process_payment"
      type: "service_call"
      config:
        service: "payment-gateway"
        endpoint: "/api/v1/charge"
        timeout: "45s"
        retry_policy:
          max_attempts: 3
          backoff: "exponential"
      input:
        amount: "$.steps.calculate_totals.output.total"
        currency: "USD"
        payment_method: "$.input.payment_method"
        customer_id: "$.input.customer_id"
        order_reference: "${uuid()}"
      output:
        transaction_id: "$.response.transaction_id"
        payment_status: "$.response.status"
        authorization_code: "$.response.auth_code"
      on_error:
        compensation:
          - name: "release_inventory"
            type: "service_call"
            config:
              service: "inventory-service"
              endpoint: "/api/v1/release"
            input:
              reservation_id: "$.steps.reserve_inventory.output.reservation_id"

    # Create order record
    - name: "create_order"
      type: "database"
      config:
        connection: "primary_db"
        operation: "insert"
        table: "orders"
        timeout: "10s"
      input:
        data:
          customer_id: "$.input.customer_id"
          items: "$.input.items"
          subtotal: "$.steps.calculate_totals.output.subtotal"
          tax: "$.steps.calculate_totals.output.tax"
          shipping: "$.steps.calculate_totals.output.shipping"
          discount: "$.steps.calculate_totals.output.discount"
          total_amount: "$.steps.calculate_totals.output.total"
          payment_method: "$.input.payment_method.type"
          transaction_id: "$.steps.process_payment.output.transaction_id"
          shipping_address: "$.steps.validation_phase.branches.address_validation.normalized_address"
          status: "confirmed"
          created_at: "${now()}"
      output:
        order_id: "$.result.id"
        created_at: "$.result.created_at"

    # Parallel fulfillment notifications
    - name: "fulfillment_notifications"
      type: "parallel"
      config:
        timeout: "30s"
        fail_fast: false
      branches:
        - name: "notify_fulfillment_center"
          steps:
            - name: "send_fulfillment_request"
              type: "message_queue"
              config:
                queue: "fulfillment-requests"
                exchange: "orders"
                routing_key: "order.created"
              input:
                message:
                  order_id: "$.steps.create_order.output.order_id"
                  customer_id: "$.input.customer_id"
                  items: "$.steps.reserve_inventory.output.reserved_items"
                  shipping_address: "$.steps.validation_phase.branches.address_validation.normalized_address"
                  priority: "normal"
        
        - name: "notify_customer"
          steps:
            - name: "send_confirmation_email"
              type: "service_call"
              config:
                service: "notification-service"
                endpoint: "/api/v1/email/send"
              input:
                template: "order_confirmation"
                recipient: "$.input.customer_id"
                data:
                  order_id: "$.steps.create_order.output.order_id"
                  total_amount: "$.steps.calculate_totals.output.total"
                  items: "$.input.items"
        
        - name: "update_analytics"
          steps:
            - name: "track_order_event"
              type: "service_call"
              config:
                service: "analytics-service"
                endpoint: "/api/v1/events"
              input:
                event: "order_created"
                properties:
                  order_id: "$.steps.create_order.output.order_id"
                  customer_id: "$.input.customer_id"
                  order_value: "$.steps.calculate_totals.output.total"
                  customer_tier: "$.steps.validation_phase.branches.customer_validation.customer_tier"
                  payment_method: "$.input.payment_method.type"

  # Final output mapping
  output:
    order_id: "$.steps.create_order.output.order_id"
    status: "confirmed"
    total_amount: "$.steps.calculate_totals.output.total"
    transaction_id: "$.steps.process_payment.output.transaction_id"
    created_at: "$.steps.create_order.output.created_at"

  # Monitoring configuration
  monitoring:
    metrics:
      - name: "order_value"
        type: "gauge"
        value: "$.steps.calculate_totals.output.total"
        labels:
          customer_tier: "$.steps.validation_phase.branches.customer_validation.customer_tier"
          payment_method: "$.input.payment_method.type"
      
      - name: "order_processing_duration"
        type: "histogram"
        value: "$.execution.duration"
        buckets: [5, 10, 30, 60, 120, 300]
    
    alerts:
      - name: "high_failure_rate"
        condition: "error_rate > 0.05"
        severity: "warning"
      
      - name: "payment_failures"
        condition: "$.steps.process_payment.status == 'failed'"
        severity: "critical"

  # Security configuration
  security:
    data_classification: "confidential"
    pii_fields:
      - "$.input.payment_method.token"
      - "$.input.shipping_address"
    audit:
      enabled: true
      events: ["workflow_started", "payment_processed", "order_created"]
```

This comprehensive YAML schema documentation provides everything needed to create, validate, and maintain Magic Flow v2 workflows effectively.