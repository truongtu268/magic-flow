# Magic Flow v2 - Dashboard Documentation

## Overview

The Magic Flow v2 Dashboard provides a comprehensive real-time monitoring and visualization platform for workflow management. Access the dashboard at `http://localhost:9090` after starting your Magic Flow v2 instance.

## Dashboard Features

### 1. Real-Time Monitoring

#### Workflow Overview Dashboard

The main dashboard provides an at-a-glance view of your entire workflow ecosystem:

**Key Metrics Panel:**
- **Active Workflows**: Currently running workflow executions
- **Total Executions Today**: Number of workflow executions in the last 24 hours
- **Success Rate**: Overall success percentage across all workflows
- **Average Duration**: Mean execution time across all workflows
- **Throughput**: Executions per hour/minute
- **Error Rate**: Percentage of failed executions

**System Health Panel:**
- **CPU Usage**: Real-time CPU utilization
- **Memory Usage**: Current memory consumption
- **Database Connections**: Active database connection count
- **Cache Hit Rate**: Redis cache performance metrics
- **Queue Depth**: Number of pending workflow executions

**Recent Activity Feed:**
- Live stream of workflow executions
- Real-time status updates
- Error notifications
- Performance alerts

#### Live Execution Monitoring

Monitor individual workflow executions in real-time:

**Execution List View:**
```
┌─────────────────┬──────────────────┬─────────────┬──────────────┬─────────────┐
│ Execution ID    │ Workflow         │ Status      │ Duration     │ Progress    │
├─────────────────┼──────────────────┼─────────────┼──────────────┼─────────────┤
│ exec-uuid-123   │ order_processing │ Running     │ 00:01:23     │ ████░░░ 60% │
│ exec-uuid-124   │ user_onboarding  │ Completed   │ 00:00:45     │ ███████ 100%│
│ exec-uuid-125   │ data_pipeline    │ Failed      │ 00:02:15     │ ███░░░░ 40% │
│ exec-uuid-126   │ order_processing │ Queued      │ -            │ ░░░░░░░ 0%  │
└─────────────────┴──────────────────┴─────────────┴──────────────┴─────────────┘
```

**Execution Detail View:**
- Step-by-step execution progress
- Real-time logs and outputs
- Performance metrics per step
- Error details and stack traces
- Input/output data inspection

### 2. Visual Workflow Flow Monitoring

#### Interactive Workflow Diagrams

Visual representation of workflow execution with real-time updates:

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          Order Processing Workflow                             │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │
│  │ Validate    │───▶│ Process     │───▶│ Update      │───▶│ Ship Order  │      │
│  │ Order ✓     │    │ Payment ⚡   │    │ Inventory   │    │             │      │
│  │ 10.2s       │    │ 25.1s       │    │ Queued      │    │ Pending     │      │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘      │
│                                                                                 │
│  Legend: ✓ Completed  ⚡ Running  ⏳ Queued  ❌ Failed                          │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**Visual Elements:**
- **Step Status Indicators**: Color-coded status (green=completed, blue=running, red=failed)
- **Execution Times**: Duration display for each step
- **Data Flow Arrows**: Show data dependencies between steps
- **Parallel Execution**: Visual representation of concurrent steps
- **Error Highlighting**: Failed steps highlighted with error details
- **Real-Time Updates**: Live status changes without page refresh

#### Workflow Topology View

High-level view of all workflows and their relationships:

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            Workflow Topology                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│     ┌─────────────────┐                    ┌─────────────────┐                 │
│     │ User Onboarding │                    │ Data Pipeline   │                 │
│     │ ✓ 245 today     │                    │ ⚡ 12 running   │                 │
│     │ 98.2% success   │                    │ 94.5% success   │                 │
│     └─────────────────┘                    └─────────────────┘                 │
│              │                                       │                         │
│              ▼                                       ▼                         │
│     ┌─────────────────┐                    ┌─────────────────┐                 │
│     │ Order Processing│◄───────────────────┤ Inventory Sync  │                 │
│     │ ⚡ 23 running   │                    │ ✓ 156 today     │                 │
│     │ 96.8% success   │                    │ 99.1% success   │                 │
│     └─────────────────┘                    └─────────────────┘                 │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3. Custom Metrics Dashboard

#### Business Metrics Configuration

Create custom dashboards for business-specific metrics:

**Configuration Example:**
```yaml
# config/dashboard_metrics.yaml
dashboards:
  - name: "E-commerce Analytics"
    description: "Order processing and revenue metrics"
    panels:
      - title: "Revenue Metrics"
        type: "time_series"
        metrics:
          - name: "total_revenue"
            query: "sum($.output.total_amount) by workflow"
            color: "#2E8B57"
          - name: "average_order_value"
            query: "avg($.output.total_amount) by hour"
            color: "#4169E1"
        
      - title: "Order Status Distribution"
        type: "pie_chart"
        metric: "order_status"
        query: "count($.output.order_status) group by status"
        
      - title: "Processing Time Trends"
        type: "histogram"
        metric: "processing_duration"
        query: "histogram($.execution.duration) buckets [1,5,10,30,60,300]"
        
      - title: "Geographic Distribution"
        type: "map"
        metric: "order_location"
        query: "count($.input.shipping_address.state) group by state"

  - name: "System Performance"
    description: "Technical performance metrics"
    panels:
      - title: "Throughput"
        type: "gauge"
        metric: "executions_per_second"
        query: "rate(workflow_executions_total[1m])"
        thresholds:
          warning: 50
          critical: 100
        
      - title: "Error Rate"
        type: "stat"
        metric: "error_percentage"
        query: "(failed_executions / total_executions) * 100"
        unit: "%"
        color_mode: "threshold"
        thresholds:
          - value: 0
            color: "green"
          - value: 5
            color: "yellow"
          - value: 10
            color: "red"
```

#### Dashboard Widgets

**Available Widget Types:**

1. **Time Series Charts**
   - Line charts for trending data
   - Multiple metrics on single chart
   - Zoom and pan capabilities
   - Custom time ranges

2. **Statistical Panels**
   - Single value displays
   - Percentage indicators
   - Trend arrows (up/down)
   - Color-coded thresholds

3. **Distribution Charts**
   - Pie charts for categorical data
   - Bar charts for comparisons
   - Histograms for distributions
   - Heatmaps for correlation analysis

4. **Geographic Visualizations**
   - World/country maps
   - Regional data overlays
   - Choropleth visualizations
   - Marker clustering

5. **Table Views**
   - Sortable data tables
   - Pagination support
   - Export capabilities
   - Drill-down functionality

#### Custom Query Language

Use JSONPath-like syntax to extract metrics from workflow data:

```javascript
// Revenue calculation
sum($.output.total_amount) where $.output.order_status == 'completed'

// Average processing time by workflow type
avg($.execution.duration) group by $.workflow.name

// Error rate calculation
(count($.execution.status == 'failed') / count(*)) * 100

// Customer segment analysis
count($.input.customer_id) group by $.input.customer_segment

// Geographic distribution
count(*) group by $.input.shipping_address.country

// Time-based aggregation
sum($.output.items[*].quantity) bucket by hour

// Percentile calculations
percentile($.execution.duration, 95) over last 24h
```

### 4. Alerting and Notification System

#### Alert Configuration

Set up intelligent alerts based on workflow performance and business metrics:

```yaml
# config/alerts.yaml
alerts:
  - name: "High Error Rate"
    description: "Alert when workflow error rate exceeds threshold"
    condition: "error_rate > 5%"
    query: "(failed_executions / total_executions) * 100 > 5"
    evaluation_interval: "1m"
    for: "5m"  # Alert after condition persists for 5 minutes
    severity: "critical"
    labels:
      team: "platform"
      service: "workflow-engine"
    annotations:
      summary: "Workflow error rate is {{ $value }}%"
      description: "The error rate for workflows has exceeded 5% for more than 5 minutes"
    
  - name: "Slow Workflow Execution"
    description: "Alert when workflow execution time is unusually high"
    condition: "avg_duration > 300s"
    query: "avg(execution_duration) > 300"
    evaluation_interval: "30s"
    for: "2m"
    severity: "warning"
    labels:
      team: "platform"
      service: "workflow-engine"
    
  - name: "Queue Backlog"
    description: "Alert when workflow queue has too many pending executions"
    condition: "queue_depth > 100"
    query: "pending_executions > 100"
    evaluation_interval: "30s"
    for: "1m"
    severity: "warning"
    
  - name: "Revenue Drop"
    description: "Business alert for significant revenue decrease"
    condition: "revenue_change < -20%"
    query: "(current_hour_revenue - previous_hour_revenue) / previous_hour_revenue * 100 < -20"
    evaluation_interval: "5m"
    for: "10m"
    severity: "critical"
    labels:
      team: "business"
      type: "revenue"

  - name: "System Resource Usage"
    description: "Alert when system resources are running low"
    conditions:
      - "cpu_usage > 80%"
      - "memory_usage > 85%"
      - "disk_usage > 90%"
    evaluation_interval: "30s"
    for: "2m"
    severity: "warning"
```

#### Notification Channels

Configure multiple notification channels for different alert types:

```yaml
# config/notifications.yaml
notification_channels:
  - name: "slack-platform-team"
    type: "slack"
    config:
      webhook_url: "https://hooks.slack.com/services/..."
      channel: "#platform-alerts"
      username: "MagicFlow Bot"
      icon_emoji: ":warning:"
      title_template: "{{ .Alert.Name }}"
      text_template: |
        *Severity:* {{ .Alert.Severity }}
        *Description:* {{ .Alert.Description }}
        *Value:* {{ .Alert.Value }}
        *Dashboard:* <{{ .DashboardURL }}|View Dashboard>
    
  - name: "email-oncall"
    type: "email"
    config:
      smtp_server: "smtp.company.com:587"
      from: "alerts@company.com"
      to: ["oncall@company.com", "platform-team@company.com"]
      subject_template: "[{{ .Alert.Severity | upper }}] {{ .Alert.Name }}"
      body_template: |
        Alert: {{ .Alert.Name }}
        Severity: {{ .Alert.Severity }}
        Description: {{ .Alert.Description }}
        
        Current Value: {{ .Alert.Value }}
        Threshold: {{ .Alert.Threshold }}
        
        Dashboard: {{ .DashboardURL }}
        Runbook: {{ .RunbookURL }}
    
  - name: "pagerduty-critical"
    type: "pagerduty"
    config:
      integration_key: "your-pagerduty-integration-key"
      severity_mapping:
        critical: "critical"
        warning: "warning"
        info: "info"
    
  - name: "webhook-custom"
    type: "webhook"
    config:
      url: "https://your-app.com/webhooks/alerts"
      method: "POST"
      headers:
        Authorization: "Bearer your-token"
        Content-Type: "application/json"
      body_template: |
        {
          "alert_name": "{{ .Alert.Name }}",
          "severity": "{{ .Alert.Severity }}",
          "value": {{ .Alert.Value }},
          "timestamp": "{{ .Alert.Timestamp }}",
          "dashboard_url": "{{ .DashboardURL }}"
        }

# Alert routing rules
routing:
  - match:
      severity: "critical"
    channels: ["slack-platform-team", "email-oncall", "pagerduty-critical"]
    
  - match:
      severity: "warning"
      team: "platform"
    channels: ["slack-platform-team"]
    
  - match:
      severity: "warning"
      team: "business"
    channels: ["email-oncall"]
    
  - match:
      type: "revenue"
    channels: ["slack-platform-team", "email-oncall"]
```

#### Alert Templates

Customize alert messages with dynamic templates:

**Slack Template Example:**
```json
{
  "blocks": [
    {
      "type": "header",
      "text": {
        "type": "plain_text",
        "text": "🚨 {{ .Alert.Name }}"
      }
    },
    {
      "type": "section",
      "fields": [
        {
          "type": "mrkdwn",
          "text": "*Severity:*\n{{ .Alert.Severity | upper }}"
        },
        {
          "type": "mrkdwn",
          "text": "*Current Value:*\n{{ .Alert.Value }}"
        },
        {
          "type": "mrkdwn",
          "text": "*Threshold:*\n{{ .Alert.Threshold }}"
        },
        {
          "type": "mrkdwn",
          "text": "*Duration:*\n{{ .Alert.Duration }}"
        }
      ]
    },
    {
      "type": "section",
      "text": {
        "type": "mrkdwn",
        "text": "{{ .Alert.Description }}"
      }
    },
    {
      "type": "actions",
      "elements": [
        {
          "type": "button",
          "text": {
            "type": "plain_text",
            "text": "View Dashboard"
          },
          "url": "{{ .DashboardURL }}"
        },
        {
          "type": "button",
          "text": {
            "type": "plain_text",
            "text": "View Runbook"
          },
          "url": "{{ .RunbookURL }}"
        }
      ]
    }
  ]
}
```

### 5. Dashboard Configuration

#### Access Control

Configure role-based access to different dashboard sections:

```yaml
# config/dashboard_access.yaml
roles:
  - name: "admin"
    permissions:
      - "dashboard:read"
      - "dashboard:write"
      - "alerts:read"
      - "alerts:write"
      - "metrics:read"
      - "metrics:write"
      - "system:read"
    
  - name: "operator"
    permissions:
      - "dashboard:read"
      - "alerts:read"
      - "metrics:read"
      - "system:read"
    
  - name: "business_user"
    permissions:
      - "dashboard:read"
      - "metrics:read"
    dashboards:
      - "E-commerce Analytics"
      - "Revenue Metrics"
    
  - name: "developer"
    permissions:
      - "dashboard:read"
      - "metrics:read"
      - "system:read"
    dashboards:
      - "System Performance"
      - "Workflow Execution"

users:
  - username: "admin@company.com"
    roles: ["admin"]
  - username: "ops@company.com"
    roles: ["operator"]
  - username: "business@company.com"
    roles: ["business_user"]
```

#### Theme Customization

Customize dashboard appearance:

```yaml
# config/dashboard_theme.yaml
theme:
  name: "company_theme"
  colors:
    primary: "#2E8B57"
    secondary: "#4169E1"
    success: "#28a745"
    warning: "#ffc107"
    danger: "#dc3545"
    info: "#17a2b8"
    light: "#f8f9fa"
    dark: "#343a40"
  
  fonts:
    family: "Inter, sans-serif"
    sizes:
      small: "12px"
      medium: "14px"
      large: "16px"
      xlarge: "20px"
  
  layout:
    sidebar_width: "250px"
    header_height: "60px"
    panel_spacing: "16px"
    border_radius: "8px"
  
  branding:
    logo_url: "/assets/company-logo.png"
    favicon_url: "/assets/favicon.ico"
    title: "Company Workflow Dashboard"
```

### 6. Dashboard API

#### Programmatic Dashboard Management

Manage dashboards programmatically via REST API:

```bash
# Create custom dashboard
curl -X POST http://localhost:9090/api/v1/dashboards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "name": "Custom Business Dashboard",
    "description": "Custom metrics for business team",
    "panels": [...]
  }'

# Get dashboard configuration
curl -X GET http://localhost:9090/api/v1/dashboards/dashboard-id \
  -H "Authorization: Bearer your-token"

# Update dashboard
curl -X PUT http://localhost:9090/api/v1/dashboards/dashboard-id \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"panels": [...]}'

# Export dashboard
curl -X GET http://localhost:9090/api/v1/dashboards/dashboard-id/export \
  -H "Authorization: Bearer your-token"

# Import dashboard
curl -X POST http://localhost:9090/api/v1/dashboards/import \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d @dashboard-export.json
```

### 7. Mobile Dashboard

Access key metrics on mobile devices:

**Mobile-Optimized Views:**
- Responsive design for tablets and phones
- Touch-friendly navigation
- Simplified metric displays
- Push notifications for critical alerts
- Offline capability for cached data

**Mobile App Features:**
- Real-time workflow status
- Push notifications for alerts
- Quick action buttons (cancel execution, acknowledge alert)
- Biometric authentication
- Dark mode support

### 8. Dashboard Integrations

#### Third-Party Integrations

**Grafana Integration:**
```yaml
# Export metrics to Grafana
grafana:
  enabled: true
  datasource:
    url: "http://localhost:9090/api/v1/metrics/prometheus"
    type: "prometheus"
  dashboards:
    auto_import: true
    folder: "Magic Flow"
```

**Datadog Integration:**
```yaml
# Send metrics to Datadog
datadog:
  enabled: true
  api_key: "your-datadog-api-key"
  tags:
    - "service:magicflow"
    - "environment:production"
  metrics_prefix: "magicflow."
```

**Prometheus Integration:**
```yaml
# Expose Prometheus metrics
prometheus:
  enabled: true
  endpoint: "/metrics"
  port: 9091
  metrics:
    - workflow_executions_total
    - workflow_duration_seconds
    - workflow_errors_total
    - system_cpu_usage
    - system_memory_usage
```

## Getting Started

### 1. Access the Dashboard

```bash
# Start Magic Flow v2
./magicflow start

# Open dashboard in browser
open http://localhost:9090
```

### 2. Initial Setup

1. **Configure Authentication**: Set up user accounts and roles
2. **Create Custom Dashboards**: Design dashboards for your specific needs
3. **Set Up Alerts**: Configure alerts for critical metrics
4. **Configure Notifications**: Set up Slack, email, or other notification channels
5. **Customize Theme**: Apply your company branding

### 3. Best Practices

**Dashboard Design:**
- Keep dashboards focused on specific use cases
- Use consistent color schemes and layouts
- Prioritize the most important metrics at the top
- Include context and descriptions for complex metrics

**Alert Configuration:**
- Set appropriate thresholds based on historical data
- Use different severity levels effectively
- Avoid alert fatigue with proper routing
- Include actionable information in alert messages

**Performance Optimization:**
- Use appropriate time ranges for queries
- Cache frequently accessed metrics
- Optimize complex queries
- Use sampling for high-volume metrics

The Magic Flow v2 Dashboard provides everything you need to monitor, analyze, and optimize your workflow operations in real-time.