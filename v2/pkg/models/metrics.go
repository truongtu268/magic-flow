package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// MetricCategory represents the category of metric
type MetricCategory string

const (
	MetricCategoryWorkflow MetricCategory = "workflow"
	MetricCategorySystem   MetricCategory = "system"
	MetricCategoryBusiness MetricCategory = "business"
	MetricCategoryCustom   MetricCategory = "custom"
)

// WorkflowMetric represents workflow execution metrics
type WorkflowMetric struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkflowID   uuid.UUID `json:"workflow_id" gorm:"type:uuid;not null;index"`
	ExecutionID  uuid.UUID `json:"execution_id" gorm:"type:uuid;not null;index"`
	VersionID    uuid.UUID `json:"version_id" gorm:"type:uuid;index"`
	
	// Metric information
	Name        string         `json:"name" gorm:"not null;index"`
	Type        MetricType     `json:"type" gorm:"not null"`
	Category    MetricCategory `json:"category" gorm:"not null;index"`
	Description string         `json:"description"`
	
	// Metric values
	Value       float64                `json:"value"`
	Unit        string                 `json:"unit"`
	Labels      map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timing information
	Timestamp   time.Time `json:"timestamp" gorm:"index"`
	Duration    int64     `json:"duration"` // in milliseconds
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	
	// Context information
	StepName    string `json:"step_name,omitempty"`
	StepType    string `json:"step_type,omitempty"`
	Environment string `json:"environment"`
	Region      string `json:"region,omitempty"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Workflow  Workflow        `json:"-" gorm:"foreignKey:WorkflowID"`
	Execution Execution       `json:"-" gorm:"foreignKey:ExecutionID"`
	Version   WorkflowVersion `json:"-" gorm:"foreignKey:VersionID"`
}

// SystemMetric represents system-level metrics
type SystemMetric struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Metric information
	Name        string         `json:"name" gorm:"not null;index"`
	Type        MetricType     `json:"type" gorm:"not null"`
	Category    MetricCategory `json:"category" gorm:"not null;index"`
	Description string         `json:"description"`
	
	// Metric values
	Value       float64                `json:"value"`
	Unit        string                 `json:"unit"`
	Labels      map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// System information
	Component   string `json:"component"` // api, engine, database, cache
	Instance    string `json:"instance"`
	Environment string `json:"environment"`
	Region      string `json:"region,omitempty"`
	
	// Timing information
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BusinessMetric represents business-level metrics
type BusinessMetric struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Metric information
	Name        string         `json:"name" gorm:"not null;index"`
	Type        MetricType     `json:"type" gorm:"not null"`
	Category    MetricCategory `json:"category" gorm:"not null;index"`
	Description string         `json:"description"`
	
	// Metric values
	Value       float64                `json:"value"`
	Unit        string                 `json:"unit"`
	Labels      map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Business context
	WorkflowID  *uuid.UUID `json:"workflow_id,omitempty" gorm:"type:uuid;index"`
	ExecutionID *uuid.UUID `json:"execution_id,omitempty" gorm:"type:uuid;index"`
	CustomerID  string     `json:"customer_id,omitempty" gorm:"index"`
	TenantID    string     `json:"tenant_id,omitempty" gorm:"index"`
	Domain      string     `json:"domain,omitempty"`
	Service     string     `json:"service,omitempty"`
	
	// Timing information
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Workflow  *Workflow  `json:"-" gorm:"foreignKey:WorkflowID"`
	Execution *Execution `json:"-" gorm:"foreignKey:ExecutionID"`
}

// MetricAggregation represents aggregated metrics
type MetricAggregation struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Aggregation information
	Name         string                 `json:"name" gorm:"not null;index"`
	Type         MetricType             `json:"type" gorm:"not null"`
	Category     MetricCategory         `json:"category" gorm:"not null;index"`
	Aggregation  string                 `json:"aggregation"` // sum, avg, min, max, count
	Interval     string                 `json:"interval"`    // 1m, 5m, 1h, 1d
	GroupBy      []string               `json:"group_by" gorm:"type:jsonb"`
	Filters      map[string]interface{} `json:"filters" gorm:"type:jsonb"`
	
	// Aggregated values
	Value       float64           `json:"value"`
	Count       int64             `json:"count"`
	Min         float64           `json:"min"`
	Max         float64           `json:"max"`
	Sum         float64           `json:"sum"`
	Avg         float64           `json:"avg"`
	Percentiles map[string]float64 `json:"percentiles" gorm:"type:jsonb"`
	
	// Time range
	StartTime time.Time `json:"start_time" gorm:"index"`
	EndTime   time.Time `json:"end_time" gorm:"index"`
	
	// Labels and metadata
	Labels   map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Alert represents an alert configuration
type Alert struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Alert information
	Name        string      `json:"name" gorm:"not null;unique"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled" gorm:"default:true"`
	Severity    string      `json:"severity"` // critical, warning, info
	Status      AlertStatus `json:"status" gorm:"default:'active'"`
	
	// Alert conditions
	Conditions AlertConditions `json:"conditions" gorm:"type:jsonb"`
	
	// Notification configuration
	Notifications []AlertNotification `json:"notifications" gorm:"type:jsonb"`
	
	// Alert metadata
	Labels   map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timing
	LastTriggered *time.Time `json:"last_triggered"`
	LastResolved  *time.Time `json:"last_resolved"`
	TriggerCount  int64      `json:"trigger_count" gorm:"default:0"`
	
	// Configuration
	CreatedBy string `json:"created_by"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	AlertStatusActive    AlertStatus = "active"
	AlertStatusTriggered AlertStatus = "triggered"
	AlertStatusResolved  AlertStatus = "resolved"
	AlertStatusSuppressed AlertStatus = "suppressed"
	AlertStatusDisabled  AlertStatus = "disabled"
)

// AlertConditions represents alert conditions
type AlertConditions struct {
	Metric      string                 `json:"metric"`
	Aggregation string                 `json:"aggregation"`
	Threshold   float64                `json:"threshold"`
	Comparison  string                 `json:"comparison"` // gt, lt, eq, gte, lte
	Duration    string                 `json:"duration"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	GroupBy     []string               `json:"group_by,omitempty"`
}

// AlertNotification represents alert notification configuration
type AlertNotification struct {
	Type     string                 `json:"type"` // email, slack, webhook, sms
	Target   string                 `json:"target"`
	Template string                 `json:"template,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

// AlertEvent represents an alert event
type AlertEvent struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AlertID uuid.UUID `json:"alert_id" gorm:"type:uuid;not null;index"`
	
	// Event information
	Type        string                 `json:"type"` // triggered, resolved, suppressed
	Message     string                 `json:"message"`
	Value       float64                `json:"value"`
	Threshold   float64                `json:"threshold"`
	Labels      map[string]string      `json:"labels" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timing
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Alert Alert `json:"-" gorm:"foreignKey:AlertID"`
}

// Dashboard represents a metrics dashboard
type Dashboard struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Dashboard information
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`
	Public      bool   `json:"public" gorm:"default:false"`
	
	// Dashboard configuration
	Layout   DashboardLayout        `json:"layout" gorm:"type:jsonb"`
	Widgets  []DashboardWidget      `json:"widgets" gorm:"type:jsonb"`
	Filters  []DashboardFilter      `json:"filters" gorm:"type:jsonb"`
	Settings map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	
	// Access control
	CreatedBy   string   `json:"created_by"`
	SharedWith  []string `json:"shared_with" gorm:"type:jsonb"`
	Permissions string   `json:"permissions"` // read, write, admin
	
	// Metadata
	Tags     []string               `json:"tags" gorm:"type:jsonb"`
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// DashboardLayout represents dashboard layout configuration
type DashboardLayout struct {
	Type    string `json:"type"` // grid, flex
	Columns int    `json:"columns"`
	Rows    int    `json:"rows"`
	Gap     int    `json:"gap"`
}

// DashboardWidget represents a dashboard widget
type DashboardWidget struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // chart, table, metric, alert
	Title    string                 `json:"title"`
	Position WidgetPosition         `json:"position"`
	Size     WidgetSize             `json:"size"`
	Query    WidgetQuery            `json:"query"`
	Config   map[string]interface{} `json:"config"`
}

// WidgetPosition represents widget position
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize represents widget size
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WidgetQuery represents widget query configuration
type WidgetQuery struct {
	Metric      string                 `json:"metric"`
	Aggregation string                 `json:"aggregation"`
	Interval    string                 `json:"interval"`
	TimeRange   string                 `json:"time_range"`
	Filters     map[string]interface{} `json:"filters"`
	GroupBy     []string               `json:"group_by"`
}

// DashboardFilter represents dashboard filter
type DashboardFilter struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"` // select, multiselect, date, text
	Field    string   `json:"field"`
	Options  []string `json:"options"`
	Default  string   `json:"default"`
	Required bool     `json:"required"`
}

// BeforeCreate hooks
func (wm *WorkflowMetric) BeforeCreate(tx *gorm.DB) error {
	if wm.ID == uuid.Nil {
		wm.ID = uuid.New()
	}
	return nil
}

func (sm *SystemMetric) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}

func (bm *BusinessMetric) BeforeCreate(tx *gorm.DB) error {
	if bm.ID == uuid.Nil {
		bm.ID = uuid.New()
	}
	return nil
}

func (ma *MetricAggregation) BeforeCreate(tx *gorm.DB) error {
	if ma.ID == uuid.Nil {
		ma.ID = uuid.New()
	}
	return nil
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (ae *AlertEvent) BeforeCreate(tx *gorm.DB) error {
	if ae.ID == uuid.Nil {
		ae.ID = uuid.New()
	}
	return nil
}

func (d *Dashboard) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// Table names
func (WorkflowMetric) TableName() string {
	return "workflow_metrics"
}

func (SystemMetric) TableName() string {
	return "system_metrics"
}

func (BusinessMetric) TableName() string {
	return "business_metrics"
}

func (MetricAggregation) TableName() string {
	return "metric_aggregations"
}

func (Alert) TableName() string {
	return "alerts"
}

func (AlertEvent) TableName() string {
	return "alert_events"
}

func (Dashboard) TableName() string {
	return "dashboards"
}

// Utility methods for WorkflowMetric
func (wm *WorkflowMetric) ToJSON() ([]byte, error) {
	return json.Marshal(wm)
}

func (wm *WorkflowMetric) FromJSON(data []byte) error {
	return json.Unmarshal(data, wm)
}

// Utility methods for Alert
func (a *Alert) Trigger() {
	now := time.Now()
	a.Status = AlertStatusTriggered
	a.LastTriggered = &now
	a.TriggerCount++
}

func (a *Alert) Resolve() {
	now := time.Now()
	a.Status = AlertStatusResolved
	a.LastResolved = &now
}

func (a *Alert) Suppress() {
	a.Status = AlertStatusSuppressed
}

func (a *Alert) IsActive() bool {
	return a.Enabled && a.Status == AlertStatusActive
}

func (a *Alert) IsTriggered() bool {
	return a.Status == AlertStatusTriggered
}

// Utility methods for Dashboard
func (d *Dashboard) AddWidget(widget DashboardWidget) {
	d.Widgets = append(d.Widgets, widget)
}

func (d *Dashboard) RemoveWidget(widgetID string) {
	for i, widget := range d.Widgets {
		if widget.ID == widgetID {
			d.Widgets = append(d.Widgets[:i], d.Widgets[i+1:]...)
			break
		}
	}
}

func (d *Dashboard) GetWidget(widgetID string) *DashboardWidget {
	for _, widget := range d.Widgets {
		if widget.ID == widgetID {
			return &widget
		}
	}
	return nil
}

func (d *Dashboard) IsPublic() bool {
	return d.Public
}

func (d *Dashboard) CanAccess(userID string) bool {
	if d.Public || d.CreatedBy == userID {
		return true
	}
	
	for _, sharedUser := range d.SharedWith {
		if sharedUser == userID {
			return true
		}
	}
	
	return false
}