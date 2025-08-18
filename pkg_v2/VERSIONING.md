# Magic Flow v2 - Workflow Versioning Documentation

## Overview

Magic Flow v2 provides comprehensive workflow versioning capabilities that enable safe evolution of workflows in production environments. This document covers version management, migration strategies, rollback procedures, and compatibility handling.

## Versioning Strategy

### Semantic Versioning

Magic Flow v2 follows semantic versioning (SemVer) for workflow versions:

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes that require migration
- **MINOR**: New features that are backward compatible
- **PATCH**: Bug fixes and minor improvements

### Version Lifecycle

```
Development → Testing → Staging → Production → Deprecated → Archived
```

**Lifecycle States:**
- **Development**: Active development, frequent changes
- **Testing**: Feature complete, undergoing testing
- **Staging**: Pre-production validation
- **Production**: Live production use
- **Deprecated**: Marked for removal, migration encouraged
- **Archived**: No longer available for new executions

## Version Management

### Creating Versions

#### CLI Commands

```bash
# Create a new version
magicflow version create \
  --workflow order_processing.yaml \
  --version 2.1.0 \
  --description "Added fraud detection step"

# Create from existing version
magicflow version create \
  --from-version 2.0.0 \
  --version 2.1.0 \
  --workflow order_processing.yaml

# Auto-increment version
magicflow version create \
  --workflow order_processing.yaml \
  --auto-increment minor
```

#### API Endpoints

```bash
# Create new version
curl -X POST "http://localhost:8080/api/v1/workflows/order-processing/versions" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "2.1.0",
    "description": "Added fraud detection step",
    "workflow_definition": {...},
    "migration_notes": "New fraud detection step added after payment validation"
  }'

# List versions
curl "http://localhost:8080/api/v1/workflows/order-processing/versions"

# Get specific version
curl "http://localhost:8080/api/v1/workflows/order-processing/versions/2.1.0"
```

### Version Metadata

```yaml
# Workflow version metadata
apiVersion: magicflow.io/v2
kind: Workflow
metadata:
  name: "order_processing"
  version: "2.1.0"
  description: "Process customer orders with fraud detection"
  
  # Version-specific metadata
  version_info:
    created_at: "2024-01-15T10:30:00Z"
    created_by: "john.doe@company.com"
    changelog: "Added fraud detection step after payment validation"
    breaking_changes: false
    migration_required: false
    
    # Compatibility information
    compatibility:
      min_platform_version: "2.0.0"
      max_platform_version: "2.9.x"
      deprecated_features: []
      removed_features: []
    
    # Dependencies
    dependencies:
      - name: "fraud-detection-service"
        version: ">=1.2.0"
        required: true
      - name: "payment-gateway"
        version: ">=2.0.0"
        required: true
    
    # Rollback information
    rollback:
      safe_rollback_versions: ["2.0.0", "2.0.1"]
      rollback_notes: "Safe to rollback to 2.0.x versions"
      data_migration_required: false

spec:
  # Workflow specification
  ...
```

### Version Comparison

```bash
# Compare versions
magicflow version compare \
  --workflow order_processing \
  --from 2.0.0 \
  --to 2.1.0

# Generate diff
magicflow version diff \
  --workflow order_processing \
  --from 2.0.0 \
  --to 2.1.0 \
  --format yaml
```

**Comparison Output:**
```yaml
comparison:
  workflow: "order_processing"
  from_version: "2.0.0"
  to_version: "2.1.0"
  
  changes:
    added_steps:
      - name: "fraud_detection"
        type: "service_call"
        position: 3
        description: "Added fraud detection validation"
    
    modified_steps:
      - name: "process_payment"
        changes:
          - field: "depends_on"
            old_value: ["validate_order"]
            new_value: ["validate_order", "fraud_detection"]
    
    removed_steps: []
    
    schema_changes:
      input_schema:
        added_fields: []
        removed_fields: []
        modified_fields: []
      
      output_schema:
        added_fields:
          - name: "fraud_score"
            type: "number"
            description: "Fraud detection score"
    
  compatibility:
    breaking_changes: false
    migration_required: false
    rollback_safe: true
```

## Migration Strategies

### Blue-Green Deployment

```yaml
# Blue-Green migration configuration
migration:
  strategy: "blue_green"
  
  # Current production version (Blue)
  blue:
    version: "2.0.0"
    traffic_percentage: 100
    instances: 3
  
  # New version (Green)
  green:
    version: "2.1.0"
    traffic_percentage: 0
    instances: 3
    
  # Migration phases
  phases:
    - name: "deploy_green"
      description: "Deploy new version alongside current"
      actions:
        - deploy_version: "2.1.0"
        - health_check: true
        - smoke_test: true
    
    - name: "canary_traffic"
      description: "Route 10% traffic to new version"
      actions:
        - route_traffic:
            blue: 90
            green: 10
        - monitor_metrics: true
        - duration: "30m"
    
    - name: "full_migration"
      description: "Route all traffic to new version"
      actions:
        - route_traffic:
            blue: 0
            green: 100
        - monitor_metrics: true
        - duration: "1h"
    
    - name: "cleanup"
      description: "Remove old version"
      actions:
        - undeploy_version: "2.0.0"
        - cleanup_resources: true
```

### Canary Deployment

```yaml
# Canary migration configuration
migration:
  strategy: "canary"
  
  # Traffic routing configuration
  traffic_routing:
    - version: "2.0.0"
      percentage: 90
      criteria:
        - type: "default"
    
    - version: "2.1.0"
      percentage: 10
      criteria:
        - type: "header"
          header: "X-Canary"
          value: "true"
        - type: "customer_tier"
          value: "premium"
  
  # Monitoring and rollback
  monitoring:
    metrics:
      - name: "error_rate"
        threshold: 0.05
        action: "rollback"
      - name: "response_time_p95"
        threshold: "2s"
        action: "alert"
    
    duration: "2h"
    auto_promote: true
    auto_rollback: true
```

### Rolling Deployment

```yaml
# Rolling migration configuration
migration:
  strategy: "rolling"
  
  # Rolling update configuration
  rolling_update:
    max_unavailable: 1
    max_surge: 1
    batch_size: 1
    batch_interval: "5m"
  
  # Health checks
  health_checks:
    readiness_probe:
      path: "/health/ready"
      initial_delay: "30s"
      period: "10s"
    
    liveness_probe:
      path: "/health/live"
      initial_delay: "60s"
      period: "30s"
  
  # Rollback configuration
  rollback:
    auto_rollback: true
    failure_threshold: 2
    rollback_timeout: "10m"
```

### Data Migration

```yaml
# Data migration configuration
migration:
  data_migration:
    required: true
    
    # Migration scripts
    scripts:
      - name: "migrate_order_schema"
        type: "sql"
        script: |
          ALTER TABLE orders ADD COLUMN fraud_score DECIMAL(3,2);
          UPDATE orders SET fraud_score = 0.0 WHERE fraud_score IS NULL;
        rollback_script: |
          ALTER TABLE orders DROP COLUMN fraud_score;
      
      - name: "migrate_execution_data"
        type: "go"
        script: "./migrations/migrate_execution_data.go"
        rollback_script: "./migrations/rollback_execution_data.go"
    
    # Migration validation
    validation:
      - name: "validate_schema"
        type: "sql"
        query: "SELECT COUNT(*) FROM orders WHERE fraud_score IS NOT NULL"
        expected_result: "> 0"
      
      - name: "validate_data_integrity"
        type: "custom"
        function: "validateDataIntegrity"
    
    # Backup configuration
    backup:
      enabled: true
      retention: "30d"
      storage: "s3://backups/workflow-migrations/"
```

## Rollback Procedures

### Automatic Rollback

```yaml
# Automatic rollback configuration
rollback:
  auto_rollback:
    enabled: true
    
    # Trigger conditions
    triggers:
      - metric: "error_rate"
        threshold: 0.05
        duration: "5m"
        action: "rollback"
      
      - metric: "response_time_p95"
        threshold: "5s"
        duration: "10m"
        action: "rollback"
      
      - metric: "success_rate"
        threshold: 0.95
        comparison: "less_than"
        duration: "5m"
        action: "rollback"
    
    # Rollback strategy
    strategy: "immediate"
    target_version: "last_stable"
    
    # Notifications
    notifications:
      - type: "slack"
        channel: "#platform-alerts"
        message: "Automatic rollback triggered for workflow {{.workflow_name}}"
      
      - type: "email"
        recipients: ["platform-team@company.com"]
        subject: "Workflow Rollback Alert"
```

### Manual Rollback

```bash
# Manual rollback commands

# Rollback to previous version
magicflow rollback \
  --workflow order_processing \
  --to-version 2.0.0 \
  --reason "High error rate in production"

# Rollback with data migration
magicflow rollback \
  --workflow order_processing \
  --to-version 2.0.0 \
  --migrate-data \
  --backup-current

# Emergency rollback (skip validations)
magicflow rollback \
  --workflow order_processing \
  --to-version 2.0.0 \
  --emergency \
  --force
```

### Rollback Validation

```yaml
# Rollback validation configuration
rollback:
  validation:
    pre_rollback:
      - name: "check_dependencies"
        type: "dependency_check"
        description: "Verify dependent services compatibility"
      
      - name: "backup_current_state"
        type: "backup"
        description: "Create backup of current version"
      
      - name: "validate_rollback_target"
        type: "version_check"
        description: "Ensure target version is valid for rollback"
    
    post_rollback:
      - name: "health_check"
        type: "health"
        timeout: "5m"
        description: "Verify system health after rollback"
      
      - name: "smoke_test"
        type: "test"
        test_suite: "smoke_tests"
        description: "Run smoke tests to verify functionality"
      
      - name: "data_integrity_check"
        type: "data_validation"
        description: "Verify data integrity after rollback"
```

## Compatibility Handling

### Backward Compatibility

```yaml
# Backward compatibility configuration
compatibility:
  backward_compatibility:
    # API compatibility
    api:
      maintain_endpoints: true
      deprecation_period: "6m"
      version_header: "X-Workflow-Version"
    
    # Schema compatibility
    schema:
      input_schema:
        allow_additional_fields: true
        ignore_unknown_fields: true
        default_values:
          new_field: "default_value"
      
      output_schema:
        maintain_existing_fields: true
        mark_deprecated_fields: true
        deprecation_warnings: true
    
    # Step compatibility
    steps:
      maintain_step_interfaces: true
      allow_step_additions: true
      prevent_step_removal: true
      step_deprecation_period: "3m"
```

### Forward Compatibility

```yaml
# Forward compatibility configuration
compatibility:
  forward_compatibility:
    # Version negotiation
    version_negotiation:
      enabled: true
      strategy: "highest_compatible"
      fallback_version: "2.0.0"
    
    # Feature flags
    feature_flags:
      - name: "fraud_detection"
        enabled_versions: [">=2.1.0"]
        fallback_behavior: "skip"
      
      - name: "advanced_analytics"
        enabled_versions: [">=2.2.0"]
        fallback_behavior: "basic_analytics"
    
    # Graceful degradation
    graceful_degradation:
      enabled: true
      strategies:
        - feature: "fraud_detection"
          fallback: "basic_validation"
        - feature: "real_time_analytics"
          fallback: "batch_analytics"
```

### Version Negotiation

```go
// Version negotiation example
type VersionNegotiator struct {
    supportedVersions []string
    defaultVersion    string
}

func (vn *VersionNegotiator) NegotiateVersion(clientVersion string) (string, error) {
    // Parse client version
    clientVer, err := semver.Parse(clientVersion)
    if err != nil {
        return vn.defaultVersion, nil
    }
    
    // Find highest compatible version
    var bestMatch string
    var bestVer semver.Version
    
    for _, supportedVersion := range vn.supportedVersions {
        supportedVer, err := semver.Parse(supportedVersion)
        if err != nil {
            continue
        }
        
        // Check compatibility
        if vn.isCompatible(clientVer, supportedVer) {
            if bestMatch == "" || supportedVer.GT(bestVer) {
                bestMatch = supportedVersion
                bestVer = supportedVer
            }
        }
    }
    
    if bestMatch == "" {
        return vn.defaultVersion, nil
    }
    
    return bestMatch, nil
}

func (vn *VersionNegotiator) isCompatible(client, server semver.Version) bool {
    // Major version must match for compatibility
    if client.Major != server.Major {
        return false
    }
    
    // Server minor version must be >= client minor version
    return server.Minor >= client.Minor
}
```

## Version Lifecycle Management

### Deprecation Process

```yaml
# Deprecation configuration
deprecation:
  policy:
    notice_period: "6m"  # 6 months notice
    support_period: "12m"  # 12 months support after deprecation
    
  # Deprecation phases
  phases:
    - name: "announcement"
      duration: "1m"
      actions:
        - announce_deprecation
        - update_documentation
        - notify_users
    
    - name: "warning_period"
      duration: "3m"
      actions:
        - add_deprecation_warnings
        - provide_migration_guides
        - offer_migration_support
    
    - name: "migration_period"
      duration: "2m"
      actions:
        - enforce_migration_timeline
        - provide_migration_tools
        - monitor_usage
    
    - name: "end_of_life"
      actions:
        - disable_new_executions
        - archive_version
        - cleanup_resources

# Deprecation notifications
notifications:
  deprecation_warning:
    message: |
      WARNING: Workflow version {{.version}} is deprecated.
      Please migrate to version {{.recommended_version}} by {{.end_of_life_date}}.
      Migration guide: {{.migration_guide_url}}
    
    channels:
      - type: "api_response_header"
        header: "X-Deprecation-Warning"
      - type: "log"
        level: "warn"
      - type: "metrics"
        metric: "deprecated_version_usage"
```

### Archive Process

```yaml
# Archive configuration
archive:
  policy:
    retention_period: "2y"  # Keep archived versions for 2 years
    
  # Archive criteria
  criteria:
    - no_active_executions: "30d"
    - deprecated_for: "12m"
    - manual_archive: true
  
  # Archive process
  process:
    - backup_version_data
    - export_execution_history
    - remove_from_active_catalog
    - update_documentation
    - notify_stakeholders
  
  # Archive storage
  storage:
    location: "s3://archives/workflow-versions/"
    encryption: true
    compression: true
    metadata_retention: "5y"
```

## Monitoring and Observability

### Version Metrics

```yaml
# Version-specific metrics
metrics:
  version_usage:
    - name: "workflow_executions_by_version"
      type: "counter"
      labels: ["workflow_name", "version", "status"]
    
    - name: "version_adoption_rate"
      type: "gauge"
      labels: ["workflow_name", "version"]
    
    - name: "migration_success_rate"
      type: "histogram"
      labels: ["from_version", "to_version"]
  
  version_health:
    - name: "version_error_rate"
      type: "gauge"
      labels: ["workflow_name", "version"]
    
    - name: "version_performance"
      type: "histogram"
      labels: ["workflow_name", "version"]
      buckets: [0.1, 0.5, 1, 2, 5, 10, 30]
```

### Version Dashboards

```yaml
# Dashboard configuration
dashboards:
  version_overview:
    panels:
      - title: "Version Distribution"
        type: "pie_chart"
        query: "workflow_executions_by_version"
        time_range: "24h"
      
      - title: "Migration Progress"
        type: "bar_chart"
        query: "version_adoption_rate"
        time_range: "7d"
      
      - title: "Version Health"
        type: "heatmap"
        query: "version_error_rate"
        time_range: "24h"
  
  migration_tracking:
    panels:
      - title: "Migration Timeline"
        type: "timeline"
        query: "migration_events"
        time_range: "30d"
      
      - title: "Rollback Events"
        type: "table"
        query: "rollback_events"
        time_range: "7d"
```

## Best Practices

### 1. Version Planning

**Plan version releases:**
```yaml
# Version release plan
release_plan:
  v2.1.0:
    target_date: "2024-02-01"
    features: ["fraud_detection", "enhanced_logging"]
    breaking_changes: false
    migration_required: false
  
  v2.2.0:
    target_date: "2024-03-15"
    features: ["real_time_analytics", "advanced_routing"]
    breaking_changes: false
    migration_required: true
  
  v3.0.0:
    target_date: "2024-06-01"
    features: ["new_execution_engine", "improved_api"]
    breaking_changes: true
    migration_required: true
```

**Use feature flags for gradual rollouts:**
```yaml
spec:
  feature_flags:
    fraud_detection:
      enabled: true
      rollout_percentage: 10
      target_versions: [">=2.1.0"]
    
    advanced_analytics:
      enabled: false
      rollout_percentage: 0
      target_versions: [">=2.2.0"]
```

### 2. Testing Strategy

**Comprehensive testing for versions:**
```yaml
testing:
  version_testing:
    # Backward compatibility tests
    backward_compatibility:
      - test_old_clients_new_version
      - test_schema_compatibility
      - test_api_compatibility
    
    # Forward compatibility tests
    forward_compatibility:
      - test_new_clients_old_version
      - test_graceful_degradation
      - test_feature_flags
    
    # Migration tests
    migration:
      - test_blue_green_deployment
      - test_canary_deployment
      - test_rollback_procedures
      - test_data_migration
```

### 3. Documentation

**Maintain comprehensive version documentation:**
```markdown
# Version 2.1.0 Release Notes

## New Features
- Added fraud detection step
- Enhanced logging and monitoring
- Improved error handling

## Breaking Changes
None

## Migration Guide
No migration required. The new fraud detection step is optional and will be skipped for existing workflows.

## Compatibility
- Backward compatible with v2.0.x
- Forward compatible with v2.x.x
- Minimum platform version: 2.0.0

## Deprecations
None

## Known Issues
- Fraud detection may add 100-200ms latency
- New metrics require Prometheus 2.30+
```

### 4. Monitoring and Alerting

**Set up version-specific monitoring:**
```yaml
alerting:
  version_alerts:
    - name: "high_error_rate_new_version"
      condition: "version_error_rate{version=~'2.1.*'} > 0.05"
      severity: "critical"
      action: "auto_rollback"
    
    - name: "slow_migration_adoption"
      condition: "version_adoption_rate{version='2.1.0'} < 0.5"
      duration: "7d"
      severity: "warning"
      action: "notify_team"
```

This comprehensive versioning documentation provides everything needed to safely manage workflow versions in Magic Flow v2, ensuring smooth evolution and reliable operations in production environments.