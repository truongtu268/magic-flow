package dashboard

import (
	"time"

	"github.com/google/uuid"
)

// DashboardConfig represents dashboard configuration for a user
type DashboardConfig struct {
	ID              uuid.UUID      `json:"id"`
	UserID          uuid.UUID      `json:"user_id"`
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	Theme           string         `json:"theme"`           // light, dark, auto
	Language        string         `json:"language"`        // en, es, fr, etc.
	Timezone        string         `json:"timezone"`        // UTC, America/New_York, etc.
	RefreshInterval int            `json:"refresh_interval"` // seconds
	Widgets         []WidgetConfig `json:"widgets"`
	Layout          LayoutConfig   `json:"layout"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// WidgetConfig represents configuration for a dashboard widget
type WidgetConfig struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`        // chart, metric, table, alert, etc.
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Position    WidgetPosition         `json:"position"`
	Size        WidgetSize             `json:"size"`
	Config      map[string]interface{} `json:"config"`
	DataSource  DataSourceConfig       `json:"data_source"`
	Visible     bool                   `json:"visible"`
	RefreshRate int                    `json:"refresh_rate"` // seconds, 0 = no auto refresh
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// WidgetPosition represents the position of a widget on the dashboard
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize represents the size of a widget
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// LayoutConfig represents the layout configuration for the dashboard
type LayoutConfig struct {
	Columns    int    `json:"columns"`
	RowHeight  int    `json:"row_height"`
	Margin     int    `json:"margin"`
	Padding    int    `json:"padding"`
	Breakpoint string `json:"breakpoint"` // xs, sm, md, lg, xl
}

// DataSourceConfig represents the data source configuration for a widget
type DataSourceConfig struct {
	Type       string                 `json:"type"`       // metrics, executions, workflows, alerts, etc.
	Endpoint   string                 `json:"endpoint"`   // API endpoint to fetch data
	Params     map[string]interface{} `json:"params"`     // Query parameters
	Filters    map[string]interface{} `json:"filters"`    // Data filters
	TimeRange  string                 `json:"time_range"` // 1h, 24h, 7d, etc.
	GroupBy    string                 `json:"group_by"`   // Field to group by
	Aggregation string                `json:"aggregation"` // sum, avg, count, etc.
}

// Widget type constants
const (
	WidgetTypeMetric           = "metric"
	WidgetTypeChart            = "chart"
	WidgetTypeTable            = "table"
	WidgetTypeAlert            = "alert"
	WidgetTypeActivity         = "activity"
	WidgetTypeStatus           = "status"
	WidgetTypeProgress         = "progress"
	WidgetTypeGauge            = "gauge"
	WidgetTypeHeatmap          = "heatmap"
	WidgetTypeTimeline         = "timeline"
	WidgetTypeKPI              = "kpi"
	WidgetTypeWorkflowList     = "workflow_list"
	WidgetTypeExecutionList    = "execution_list"
	WidgetTypeSystemHealth     = "system_health"
	WidgetTypeResourceUsage    = "resource_usage"
	WidgetTypePerformanceTrend = "performance_trend"
)

// Chart type constants
const (
	ChartTypeLine      = "line"
	ChartTypeBar       = "bar"
	ChartTypePie       = "pie"
	ChartTypeDoughnut  = "doughnut"
	ChartTypeArea      = "area"
	ChartTypeScatter   = "scatter"
	ChartTypeHistogram = "histogram"
	ChartTypeBoxPlot   = "box_plot"
)

// HealthStatus represents the health status of the dashboard service
type HealthStatus struct {
	Status    string                    `json:"status"`    // healthy, degraded, unhealthy
	Timestamp time.Time               `json:"timestamp"`
	Services  map[string]ServiceHealth `json:"services"`
	Uptime    time.Duration            `json:"uptime"`
	Version   string                   `json:"version"`
}

// ServiceHealth represents the health status of a specific service
type ServiceHealth struct {
	Status      string        `json:"status"`      // healthy, unhealthy
	Latency     time.Duration `json:"latency"`     // Response time
	Error       string        `json:"error,omitempty"`
	LastChecked time.Time     `json:"last_checked"`
}

// AlertConfiguration represents alert configuration for dashboard monitoring
type AlertConfiguration struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Metric      string                 `json:"metric"`      // execution_failure_rate, response_time, etc.
	Threshold   float64                `json:"threshold"`   // Alert threshold value
	Operator    string                 `json:"operator"`    // >, <, >=, <=, ==, !=
	Severity    string                 `json:"severity"`    // low, medium, high, critical
	Enabled     bool                   `json:"enabled"`
	Conditions  []AlertCondition       `json:"conditions"`
	Actions     []AlertAction          `json:"actions"`
	Cooldown    time.Duration          `json:"cooldown"`    // Minimum time between alerts
	Tags        map[string]string      `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// AlertCondition represents a condition for triggering an alert
type AlertCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Logic    string      `json:"logic"` // AND, OR
}

// AlertAction represents an action to take when an alert is triggered
type AlertAction struct {
	Type   string                 `json:"type"`   // email, webhook, slack, etc.
	Target string                 `json:"target"` // email address, webhook URL, etc.
	Config map[string]interface{} `json:"config"` // Additional configuration
}

// DashboardTemplate represents a pre-defined dashboard template
type DashboardTemplate struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`    // monitoring, analytics, operations, etc.
	Tags        []string       `json:"tags"`
	Widgets     []WidgetConfig `json:"widgets"`
	Layout      LayoutConfig   `json:"layout"`
	Preview     string         `json:"preview"`     // URL to preview image
	Popular     bool           `json:"popular"`     // Is this a popular template?
	CreatedBy   uuid.UUID      `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UserPreferences represents user preferences for the dashboard
type UserPreferences struct {
	UserID              uuid.UUID         `json:"user_id"`
	DefaultDashboard    *uuid.UUID        `json:"default_dashboard,omitempty"`
	Theme               string            `json:"theme"`
	Language            string            `json:"language"`
	Timezone            string            `json:"timezone"`
	DateFormat          string            `json:"date_format"`
	TimeFormat          string            `json:"time_format"`
	Notifications       NotificationPrefs `json:"notifications"`
	AutoRefresh         bool              `json:"auto_refresh"`
	RefreshInterval     int               `json:"refresh_interval"`
	CompactMode         bool              `json:"compact_mode"`
	ShowTooltips        bool              `json:"show_tooltips"`
	AnimationsEnabled   bool              `json:"animations_enabled"`
	KeyboardShortcuts   bool              `json:"keyboard_shortcuts"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// NotificationPrefs represents notification preferences
type NotificationPrefs struct {
	Email     bool `json:"email"`
	Browser   bool `json:"browser"`
	Slack     bool `json:"slack"`
	Webhook   bool `json:"webhook"`
	Alerts    bool `json:"alerts"`
	Executions bool `json:"executions"`
	Workflows bool `json:"workflows"`
	System    bool `json:"system"`
}

// DashboardShare represents a shared dashboard configuration
type DashboardShare struct {
	ID           uuid.UUID `json:"id"`
	DashboardID  uuid.UUID `json:"dashboard_id"`
	ShareToken   string    `json:"share_token"`
	SharedBy     uuid.UUID `json:"shared_by"`
	SharedWith   *uuid.UUID `json:"shared_with,omitempty"` // nil for public shares
	Permissions  []string  `json:"permissions"`            // read, write, admin
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	PasswordHash *string   `json:"password_hash,omitempty"`
	AccessCount  int64     `json:"access_count"`
	LastAccessed *time.Time `json:"last_accessed,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DashboardExport represents an exported dashboard configuration
type DashboardExport struct {
	Version     string           `json:"version"`
	ExportedAt  time.Time        `json:"exported_at"`
	ExportedBy  uuid.UUID        `json:"exported_by"`
	Dashboard   DashboardConfig  `json:"dashboard"`
	Templates   []WidgetTemplate `json:"templates,omitempty"`
	Dependencies []string        `json:"dependencies,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// WidgetTemplate represents a reusable widget template
type WidgetTemplate struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	Config      map[string]interface{} `json:"config"`
	DataSource  DataSourceConfig       `json:"data_source"`
	Preview     string                 `json:"preview"`
	Popular     bool                   `json:"popular"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// DashboardAnalytics represents analytics data for dashboard usage
type DashboardAnalytics struct {
	DashboardID     uuid.UUID `json:"dashboard_id"`
	ViewCount       int64     `json:"view_count"`
	UniqueViewers   int64     `json:"unique_viewers"`
	AverageViewTime time.Duration `json:"average_view_time"`
	BounceRate      float64   `json:"bounce_rate"`
	MostViewedWidget string   `json:"most_viewed_widget"`
	LeastViewedWidget string  `json:"least_viewed_widget"`
	PeakUsageHour   int       `json:"peak_usage_hour"`
	LastViewed      time.Time `json:"last_viewed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DashboardPermission represents permissions for dashboard access
type DashboardPermission struct {
	ID          uuid.UUID `json:"id"`
	DashboardID uuid.UUID `json:"dashboard_id"`
	UserID      uuid.UUID `json:"user_id"`
	Role        string    `json:"role"`        // viewer, editor, admin
	Permissions []string  `json:"permissions"` // read, write, delete, share
	GrantedBy   uuid.UUID `json:"granted_by"`
	GrantedAt   time.Time `json:"granted_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission constants
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
	PermissionShare  = "share"
	PermissionAdmin  = "admin"
)

// Role constants
const (
	RoleViewer = "viewer"
	RoleEditor = "editor"
	RoleAdmin  = "admin"
	RoleOwner  = "owner"
)

// Theme constants
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
	ThemeAuto  = "auto"
)

// Status constants
const (
	StatusHealthy   = "healthy"
	StatusDegraded  = "degraded"
	StatusUnhealthy = "unhealthy"
)

// Severity constants
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Operator constants
const (
	OperatorGreaterThan      = ">"
	OperatorLessThan         = "<"
	OperatorGreaterThanEqual = ">="
	OperatorLessThanEqual    = "<="
	OperatorEqual            = "=="
	OperatorNotEqual         = "!="
	OperatorContains         = "contains"
	OperatorStartsWith       = "starts_with"
	OperatorEndsWith         = "ends_with"
)