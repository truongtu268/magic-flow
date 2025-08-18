package database

import (
	"gorm.io/gorm"
)

// repositoryManager implements RepositoryManager interface
type repositoryManager struct {
	db                       *gorm.DB
	workflowRepo             WorkflowRepository
	executionRepo            ExecutionRepository
	stepExecutionRepo        StepExecutionRepository
	executionEventRepo       ExecutionEventRepository
	workflowVersionRepo      WorkflowVersionRepository
	metricsRepo              MetricsRepository
	alertRepo                AlertRepository
	dashboardRepo            DashboardRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) RepositoryManager {
	return &repositoryManager{
		db:                  db,
		workflowRepo:        NewWorkflowRepository(db),
		executionRepo:       NewExecutionRepository(db),
		stepExecutionRepo:   NewStepExecutionRepository(db),
		executionEventRepo:  NewExecutionEventRepository(db),
		workflowVersionRepo: NewWorkflowVersionRepository(db),
		metricsRepo:         NewMetricsRepository(db),
		alertRepo:           NewAlertRepository(db),
		dashboardRepo:       NewDashboardRepository(db),
	}
}

// Workflow returns the workflow repository
func (rm *repositoryManager) Workflow() WorkflowRepository {
	return rm.workflowRepo
}

// Execution returns the execution repository
func (rm *repositoryManager) Execution() ExecutionRepository {
	return rm.executionRepo
}

// StepExecution returns the step execution repository
func (rm *repositoryManager) StepExecution() StepExecutionRepository {
	return rm.stepExecutionRepo
}

// ExecutionEvent returns the execution event repository
func (rm *repositoryManager) ExecutionEvent() ExecutionEventRepository {
	return rm.executionEventRepo
}

// WorkflowVersion returns the workflow version repository
func (rm *repositoryManager) WorkflowVersion() WorkflowVersionRepository {
	return rm.workflowVersionRepo
}

// Metrics returns the metrics repository
func (rm *repositoryManager) Metrics() MetricsRepository {
	return rm.metricsRepo
}

// Alert returns the alert repository
func (rm *repositoryManager) Alert() AlertRepository {
	return rm.alertRepo
}

// Dashboard returns the dashboard repository
func (rm *repositoryManager) Dashboard() DashboardRepository {
	return rm.dashboardRepo
}

// Transaction executes a function within a database transaction
func (rm *repositoryManager) Transaction(fn func(RepositoryManager) error) error {
	return rm.db.Transaction(func(tx *gorm.DB) error {
		// Create a new repository manager with the transaction
		txManager := NewRepositoryManager(tx)
		return fn(txManager)
	})
}

// GetDB returns the underlying database connection
func (rm *repositoryManager) GetDB() *gorm.DB {
	return rm.db
}

// Health checks the health of all repositories
func (rm *repositoryManager) Health() error {
	// Check database connection
	sqlDB, err := rm.db.DB()
	if err != nil {
		return err
	}

	// Ping the database
	return sqlDB.Ping()
}

// Close closes all repository connections
func (rm *repositoryManager) Close() error {
	sqlDB, err := rm.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}