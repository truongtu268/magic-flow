package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"magic-flow/v2/pkg/config"
	"magic-flow/v2/pkg/models"
)

// Database represents the database connection and configuration
type Database struct {
	DB     *gorm.DB
	Config *config.DatabaseConfig
	Logger *logrus.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.DatabaseConfig, log *logrus.Logger) (*Database, error) {
	db := &Database{
		Config: cfg,
		Logger: log,
	}

	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Connect establishes a connection to the database
func (d *Database) Connect() error {
	var dialector gorm.Dialector

	switch d.Config.Driver {
	case "postgres":
		dialector = postgres.Open(d.Config.GetConnectionString())
	case "mysql":
		dialector = mysql.Open(d.Config.GetConnectionString())
	default:
		return fmt.Errorf("unsupported database driver: %s", d.Config.Driver)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		d.Logger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  d.getLogLevel(),
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:           true,
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(d.Config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.Config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(d.Config.ConnMaxLifetime) * time.Second)

	d.DB = db
	d.Logger.Info("Database connection established")

	return nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	d.Logger.Info("Running database migrations...")

	err := d.DB.AutoMigrate(
		&models.Workflow{},
		&models.Execution{},
		&models.StepExecution{},
		&models.ExecutionEvent{},
		&models.WorkflowVersion{},
		&models.Deployment{},
		&models.WorkflowMetric{},
		&models.SystemMetric{},
		&models.BusinessMetric{},
		&models.MetricAggregation{},
		&models.Alert{},
		&models.AlertEvent{},
		&models.Dashboard{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.Logger.Info("Database migrations completed")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// Health checks the database connection health
func (d *Database) Health() error {
	if d.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Ping()
}

// Transaction executes a function within a database transaction
func (d *Database) Transaction(fn func(*gorm.DB) error) error {
	return d.DB.Transaction(fn)
}

// GetDB returns the GORM database instance
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

func (d *Database) getLogLevel() logger.LogLevel {
	switch d.Config.LogLevel {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

// CreateIndexes creates database indexes for better performance
func (d *Database) CreateIndexes() error {
	d.Logger.Info("Creating database indexes...")

	indexes := []string{
		// Workflow indexes
		"CREATE INDEX IF NOT EXISTS idx_workflows_name ON workflows(name);",
		"CREATE INDEX IF NOT EXISTS idx_workflows_status ON workflows(status);",
		"CREATE INDEX IF NOT EXISTS idx_workflows_created_at ON workflows(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_workflows_updated_at ON workflows(updated_at);",

		// Execution indexes
		"CREATE INDEX IF NOT EXISTS idx_executions_workflow_id ON executions(workflow_id);",
		"CREATE INDEX IF NOT EXISTS idx_executions_status ON executions(status);",
		"CREATE INDEX IF NOT EXISTS idx_executions_started_at ON executions(started_at);",
		"CREATE INDEX IF NOT EXISTS idx_executions_completed_at ON executions(completed_at);",
		"CREATE INDEX IF NOT EXISTS idx_executions_trigger_type ON executions(trigger_type);",

		// Step execution indexes
		"CREATE INDEX IF NOT EXISTS idx_step_executions_execution_id ON step_executions(execution_id);",
		"CREATE INDEX IF NOT EXISTS idx_step_executions_step_id ON step_executions(step_id);",
		"CREATE INDEX IF NOT EXISTS idx_step_executions_status ON step_executions(status);",

		// Execution event indexes
		"CREATE INDEX IF NOT EXISTS idx_execution_events_execution_id ON execution_events(execution_id);",
		"CREATE INDEX IF NOT EXISTS idx_execution_events_type ON execution_events(type);",
		"CREATE INDEX IF NOT EXISTS idx_execution_events_timestamp ON execution_events(timestamp);",

		// Workflow version indexes
		"CREATE INDEX IF NOT EXISTS idx_workflow_versions_workflow_id ON workflow_versions(workflow_id);",
		"CREATE INDEX IF NOT EXISTS idx_workflow_versions_version ON workflow_versions(version);",
		"CREATE INDEX IF NOT EXISTS idx_workflow_versions_status ON workflow_versions(status);",

		// Deployment indexes
		"CREATE INDEX IF NOT EXISTS idx_deployments_workflow_version_id ON deployments(workflow_version_id);",
		"CREATE INDEX IF NOT EXISTS idx_deployments_environment ON deployments(environment);",
		"CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);",

		// Metric indexes
		"CREATE INDEX IF NOT EXISTS idx_workflow_metrics_name ON workflow_metrics(name);",
		"CREATE INDEX IF NOT EXISTS idx_workflow_metrics_timestamp ON workflow_metrics(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_workflow_metrics_workflow_id ON workflow_metrics((labels->>'workflow_id'));",

		"CREATE INDEX IF NOT EXISTS idx_system_metrics_name ON system_metrics(name);",
		"CREATE INDEX IF NOT EXISTS idx_system_metrics_timestamp ON system_metrics(timestamp);",

		"CREATE INDEX IF NOT EXISTS idx_business_metrics_name ON business_metrics(name);",
		"CREATE INDEX IF NOT EXISTS idx_business_metrics_timestamp ON business_metrics(timestamp);",

		// Alert indexes
		"CREATE INDEX IF NOT EXISTS idx_alerts_name ON alerts(name);",
		"CREATE INDEX IF NOT EXISTS idx_alerts_enabled ON alerts(enabled);",
		"CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity);",

		"CREATE INDEX IF NOT EXISTS idx_alert_events_alert_id ON alert_events(alert_id);",
		"CREATE INDEX IF NOT EXISTS idx_alert_events_timestamp ON alert_events(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_alert_events_status ON alert_events(status);",

		// Dashboard indexes
		"CREATE INDEX IF NOT EXISTS idx_dashboards_name ON dashboards(name);",
		"CREATE INDEX IF NOT EXISTS idx_dashboards_created_by ON dashboards(created_by);",
		"CREATE INDEX IF NOT EXISTS idx_dashboards_is_public ON dashboards(is_public);",
	}

	for _, indexSQL := range indexes {
		if err := d.DB.Exec(indexSQL).Error; err != nil {
			d.Logger.WithError(err).Warnf("Failed to create index: %s", indexSQL)
			// Continue with other indexes even if one fails
		}
	}

	d.Logger.Info("Database indexes created")
	return nil
}

// SeedData inserts initial data into the database
func (d *Database) SeedData() error {
	d.Logger.Info("Seeding initial data...")

	// Check if data already exists
	var count int64
	d.DB.Model(&models.Workflow{}).Count(&count)
	if count > 0 {
		d.Logger.Info("Data already exists, skipping seed")
		return nil
	}

	// Create sample workflow
	sampleWorkflow := &models.Workflow{
		Name:        "Sample HTTP Workflow",
		Description: "A sample workflow that makes HTTP requests",
		Status:      models.WorkflowStatusActive,
		Definition: &models.WorkflowDefinition{
			Version: "1.0.0",
			Metadata: models.WorkflowMetadata{
				Name:        "Sample HTTP Workflow",
				Description: "A sample workflow that makes HTTP requests",
				Version:     "1.0.0",
				Labels: map[string]string{
					"category": "sample",
					"type":     "http",
				},
			},
			Spec: models.WorkflowSpec{
				Steps: []models.WorkflowStep{
					{
						ID:          "fetch-data",
						Name:        "Fetch Data",
						Description: "Fetch data from API",
						Type:        "http",
						Config: map[string]interface{}{
							"url":    "https://jsonplaceholder.typicode.com/posts/1",
							"method": "GET",
						},
						Timeout: "30s",
					},
					{
						ID:          "process-data",
						Name:        "Process Data",
						Description: "Process the fetched data",
						Type:        "transform",
						Config: map[string]interface{}{
							"type":       "json",
							"expression": ".title",
						},
						DependsOn: []string{"fetch-data"},
						Timeout:   "10s",
					},
				},
				Triggers: []models.WorkflowTrigger{
					{
						Type: "manual",
						Config: map[string]interface{}{
							"description": "Manual trigger for testing",
						},
					},
				},
			},
		},
	}

	if err := d.DB.Create(sampleWorkflow).Error; err != nil {
		return fmt.Errorf("failed to create sample workflow: %w", err)
	}

	// Create sample dashboard
	sampleDashboard := &models.Dashboard{
		Name:        "System Overview",
		Description: "Overview of system metrics and workflow status",
		Config: map[string]interface{}{
			"layout": "grid",
			"widgets": []map[string]interface{}{
				{
					"type":  "metric",
					"title": "Active Workflows",
					"query": "workflow_executions_active",
				},
				{
					"type":  "chart",
					"title": "Execution Success Rate",
					"query": "workflow_executions_success_rate",
				},
			},
		},
		IsPublic:  true,
		CreatedBy: "system",
	}

	if err := d.DB.Create(sampleDashboard).Error; err != nil {
		return fmt.Errorf("failed to create sample dashboard: %w", err)
	}

	d.Logger.Info("Initial data seeded successfully")
	return nil
}