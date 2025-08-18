-- Drop triggers
DROP TRIGGER IF EXISTS update_dashboards_updated_at ON dashboards;
DROP TRIGGER IF EXISTS update_alerts_updated_at ON alerts;
DROP TRIGGER IF EXISTS update_deployments_updated_at ON deployments;
DROP TRIGGER IF EXISTS update_workflow_versions_updated_at ON workflow_versions;
DROP TRIGGER IF EXISTS update_step_executions_updated_at ON step_executions;
DROP TRIGGER IF EXISTS update_executions_updated_at ON executions;
DROP TRIGGER IF EXISTS update_workflows_updated_at ON workflows;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_dashboards_is_public;
DROP INDEX IF EXISTS idx_dashboards_created_by;
DROP INDEX IF EXISTS idx_dashboards_name;

DROP INDEX IF EXISTS idx_alert_events_status;
DROP INDEX IF EXISTS idx_alert_events_timestamp;
DROP INDEX IF EXISTS idx_alert_events_alert_id;

DROP INDEX IF EXISTS idx_alerts_severity;
DROP INDEX IF EXISTS idx_alerts_enabled;
DROP INDEX IF EXISTS idx_alerts_name;

DROP INDEX IF EXISTS idx_metric_aggregations_start_time;
DROP INDEX IF EXISTS idx_metric_aggregations_time_window;
DROP INDEX IF EXISTS idx_metric_aggregations_metric_name;

DROP INDEX IF EXISTS idx_business_metrics_timestamp;
DROP INDEX IF EXISTS idx_business_metrics_name;

DROP INDEX IF EXISTS idx_system_metrics_timestamp;
DROP INDEX IF EXISTS idx_system_metrics_name;

DROP INDEX IF EXISTS idx_workflow_metrics_workflow_id;
DROP INDEX IF EXISTS idx_workflow_metrics_timestamp;
DROP INDEX IF EXISTS idx_workflow_metrics_name;

DROP INDEX IF EXISTS idx_deployments_status;
DROP INDEX IF EXISTS idx_deployments_environment;
DROP INDEX IF EXISTS idx_deployments_workflow_version_id;

DROP INDEX IF EXISTS idx_workflow_versions_status;
DROP INDEX IF EXISTS idx_workflow_versions_version;
DROP INDEX IF EXISTS idx_workflow_versions_workflow_id;

DROP INDEX IF EXISTS idx_execution_events_timestamp;
DROP INDEX IF EXISTS idx_execution_events_type;
DROP INDEX IF EXISTS idx_execution_events_execution_id;

DROP INDEX IF EXISTS idx_step_executions_status;
DROP INDEX IF EXISTS idx_step_executions_step_id;
DROP INDEX IF EXISTS idx_step_executions_execution_id;

DROP INDEX IF EXISTS idx_executions_trigger_type;
DROP INDEX IF EXISTS idx_executions_completed_at;
DROP INDEX IF EXISTS idx_executions_started_at;
DROP INDEX IF EXISTS idx_executions_status;
DROP INDEX IF EXISTS idx_executions_workflow_id;

DROP INDEX IF EXISTS idx_workflows_updated_at;
DROP INDEX IF EXISTS idx_workflows_created_at;
DROP INDEX IF EXISTS idx_workflows_status;
DROP INDEX IF EXISTS idx_workflows_name;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS dashboards;
DROP TABLE IF EXISTS alert_events;
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS metric_aggregations;
DROP TABLE IF EXISTS business_metrics;
DROP TABLE IF EXISTS system_metrics;
DROP TABLE IF EXISTS workflow_metrics;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS workflow_versions;
DROP TABLE IF EXISTS execution_events;
DROP TABLE IF EXISTS step_executions;
DROP TABLE IF EXISTS executions;
DROP TABLE IF EXISTS workflows;