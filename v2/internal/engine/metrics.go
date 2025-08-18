package engine

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"magic-flow/v2/pkg/models"
)

// PrometheusMetricsCollector implements MetricsCollector using Prometheus
type PrometheusMetricsCollector struct {
	registry *prometheus.Registry
	metrics  map[string]prometheus.Collector
	mutex    sync.RWMutex
	logger   *logrus.Logger
}

// NewPrometheusMetricsCollector creates a new Prometheus metrics collector
func NewPrometheusMetricsCollector(logger *logrus.Logger) *PrometheusMetricsCollector {
	registry := prometheus.NewRegistry()
	c := &PrometheusMetricsCollector{
		registry: registry,
		metrics:  make(map[string]prometheus.Collector),
		logger:   logger,
	}

	// Register default metrics
	c.registerDefaultMetrics()

	return c
}

func (c *PrometheusMetricsCollector) registerDefaultMetrics() {
	// Workflow execution metrics
	c.registerCounter("workflow_executions_started_total", "Total number of workflow executions started", []string{"workflow_id", "event_type"})
	c.registerCounter("workflow_executions_completed_total", "Total number of workflow executions completed", []string{"workflow_id", "event_type"})
	c.registerCounter("workflow_executions_failed_total", "Total number of workflow executions failed", []string{"workflow_id", "event_type"})
	c.registerCounter("workflow_executions_cancelled_total", "Total number of workflow executions cancelled", []string{"workflow_id", "event_type"})

	// Workflow step metrics
	c.registerCounter("workflow_steps_started_total", "Total number of workflow steps started", []string{"workflow_id", "step_id", "event_type"})
	c.registerCounter("workflow_steps_completed_total", "Total number of workflow steps completed", []string{"workflow_id", "step_id", "event_type"})
	c.registerCounter("workflow_steps_failed_total", "Total number of workflow steps failed", []string{"workflow_id", "step_id", "event_type"})

	// Duration metrics
	c.registerHistogram("workflow_execution_duration_seconds", "Duration of workflow executions in seconds", []string{"workflow_id", "event_type"}, []float64{0.1, 0.5, 1, 5, 10, 30, 60, 300, 600, 1800, 3600})
	c.registerHistogram("workflow_step_duration_seconds", "Duration of workflow steps in seconds", []string{"workflow_id", "step_id", "event_type"}, []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60})

	// Engine metrics
	c.registerGauge("workflow_engine_active_executions", "Number of currently active workflow executions", []string{})
	c.registerGauge("workflow_engine_queue_size", "Number of workflow executions in queue", []string{})
	c.registerCounter("workflow_engine_errors_total", "Total number of workflow engine errors", []string{"error_type"})
}

func (c *PrometheusMetricsCollector) registerCounter(name, help string, labels []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	counter := promauto.With(c.registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	c.metrics[name] = counter
}

func (c *PrometheusMetricsCollector) registerGauge(name, help string, labels []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	gauge := promauto.With(c.registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	c.metrics[name] = gauge
}

func (c *PrometheusMetricsCollector) registerHistogram(name, help string, labels []string, buckets []float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	histogram := promauto.With(c.registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
		labels,
	)
	c.metrics[name] = histogram
}

func (c *PrometheusMetricsCollector) RecordMetric(name string, value float64, labels map[string]string) {
	c.mutex.RLock()
	metric, exists := c.metrics[name]
	c.mutex.RUnlock()

	if !exists {
		c.logger.WithFields(logrus.Fields{
			"metric_name": name,
			"labels":      labels,
		}).Warn("Metric not found")
		return
	}

	switch m := metric.(type) {
	case *prometheus.CounterVec:
		m.With(labels).Add(value)
	case *prometheus.GaugeVec:
		m.With(labels).Set(value)
	case *prometheus.HistogramVec:
		m.With(labels).Observe(value)
	default:
		c.logger.WithFields(logrus.Fields{
			"metric_name": name,
			"metric_type": m,
		}).Warn("Unknown metric type")
	}
}

func (c *PrometheusMetricsCollector) IncrementCounter(name string, labels map[string]string) {
	c.RecordMetric(name, 1, labels)
}

func (c *PrometheusMetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	c.RecordMetric(name, value, labels)
}

func (c *PrometheusMetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	c.RecordMetric(name, value, labels)
}

func (c *PrometheusMetricsCollector) GetRegistry() *prometheus.Registry {
	return c.registry
}

// DatabaseMetricsCollector implements MetricsCollector by storing metrics in database
type DatabaseMetricsCollector struct {
	db     *gorm.DB
	logger *logrus.Logger
	buffer chan *models.WorkflowMetric
	stop   chan struct{}
	wg     sync.WaitGroup
}

// NewDatabaseMetricsCollector creates a new database metrics collector
func NewDatabaseMetricsCollector(db *gorm.DB, logger *logrus.Logger) *DatabaseMetricsCollector {
	c := &DatabaseMetricsCollector{
		db:     db,
		logger: logger,
		buffer: make(chan *models.WorkflowMetric, 1000),
		stop:   make(chan struct{}),
	}

	// Start background worker to flush metrics
	c.wg.Add(1)
	go c.flushWorker()

	return c
}

func (c *DatabaseMetricsCollector) RecordMetric(name string, value float64, labels map[string]string) {
	metric := &models.WorkflowMetric{
		Name:      name,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	// Try to send to buffer, drop if full
	select {
	case c.buffer <- metric:
	default:
		c.logger.WithFields(logrus.Fields{
			"metric_name": name,
			"buffer_size": len(c.buffer),
		}).Warn("Metrics buffer full, dropping metric")
	}
}

func (c *DatabaseMetricsCollector) IncrementCounter(name string, labels map[string]string) {
	c.RecordMetric(name, 1, labels)
}

func (c *DatabaseMetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	c.RecordMetric(name, value, labels)
}

func (c *DatabaseMetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	c.RecordMetric(name, value, labels)
}

func (c *DatabaseMetricsCollector) flushWorker() {
	defer c.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	batch := make([]*models.WorkflowMetric, 0, 100)

	for {
		select {
		case <-c.stop:
			// Flush remaining metrics
			c.flushBatch(batch)
			return

		case metric := <-c.buffer:
			batch = append(batch, metric)
			if len(batch) >= 100 {
				c.flushBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				c.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

func (c *DatabaseMetricsCollector) flushBatch(metrics []*models.WorkflowMetric) {
	if len(metrics) == 0 {
		return
	}

	if err := c.db.CreateInBatches(metrics, 100).Error; err != nil {
		c.logger.WithFields(logrus.Fields{
			"batch_size": len(metrics),
			"error":      err.Error(),
		}).Error("Failed to flush metrics batch")
	} else {
		c.logger.WithFields(logrus.Fields{
			"batch_size": len(metrics),
		}).Debug("Flushed metrics batch")
	}
}

func (c *DatabaseMetricsCollector) Close() error {
	close(c.stop)
	c.wg.Wait()
	close(c.buffer)
	return nil
}

// CompositeMetricsCollector combines multiple metrics collectors
type CompositeMetricsCollector struct {
	collectors []MetricsCollector
	logger     *logrus.Logger
}

// NewCompositeMetricsCollector creates a new composite metrics collector
func NewCompositeMetricsCollector(collectors []MetricsCollector, logger *logrus.Logger) *CompositeMetricsCollector {
	return &CompositeMetricsCollector{
		collectors: collectors,
		logger:     logger,
	}
}

func (c *CompositeMetricsCollector) RecordMetric(name string, value float64, labels map[string]string) {
	for _, collector := range c.collectors {
		collector.RecordMetric(name, value, labels)
	}
}

func (c *CompositeMetricsCollector) IncrementCounter(name string, labels map[string]string) {
	for _, collector := range c.collectors {
		collector.IncrementCounter(name, labels)
	}
}

func (c *CompositeMetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	for _, collector := range c.collectors {
		collector.SetGauge(name, value, labels)
	}
}

func (c *CompositeMetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	for _, collector := range c.collectors {
		collector.ObserveHistogram(name, value, labels)
	}
}

// NoOpMetricsCollector is a no-op implementation for testing
type NoOpMetricsCollector struct{}

// NewNoOpMetricsCollector creates a new no-op metrics collector
func NewNoOpMetricsCollector() *NoOpMetricsCollector {
	return &NoOpMetricsCollector{}
}

func (c *NoOpMetricsCollector) RecordMetric(name string, value float64, labels map[string]string) {}
func (c *NoOpMetricsCollector) IncrementCounter(name string, labels map[string]string)         {}
func (c *NoOpMetricsCollector) SetGauge(name string, value float64, labels map[string]string)  {}
func (c *NoOpMetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {}