package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/magic-flow/v2/internal/api"
	"github.com/magic-flow/v2/internal/database"
	"github.com/magic-flow/v2/internal/engine"
	"github.com/magic-flow/v2/internal/metrics"
	"github.com/magic-flow/v2/internal/services"
	"github.com/magic-flow/v2/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile string
	logLevel   string
	port       int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "magic-flow-server",
		Short: "Magic Flow v2 API Server",
		Long:  "Magic Flow v2 is a powerful workflow orchestration platform with visual design capabilities.",
		Run:   runServer,
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Configuration file path")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override port if specified
	if port != 8080 {
		cfg.Server.Port = port
	}

	// Setup logging
	setupLogging(logLevel, cfg.Logging)

	logrus.Info("Starting Magic Flow v2 Server...")

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize metrics
	metricsCollector := metrics.NewCollector(cfg.Metrics)
	if err := metricsCollector.Start(); err != nil {
		logrus.Fatalf("Failed to start metrics collector: %v", err)
	}

	// Initialize services
	serviceContainer := services.NewContainer(db, cfg)

	// Initialize workflow engine
	workflowEngine := engine.New(serviceContainer, cfg.Engine)
	if err := workflowEngine.Start(); err != nil {
		logrus.Fatalf("Failed to start workflow engine: %v", err)
	}

	// Setup Gin router
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup API routes
	apiHandler := api.NewHandler(serviceContainer, workflowEngine, metricsCollector)
	apiHandler.SetupRoutes(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown workflow engine
	if err := workflowEngine.Stop(); err != nil {
		logrus.Errorf("Error stopping workflow engine: %v", err)
	}

	// Shutdown metrics collector
	if err := metricsCollector.Stop(); err != nil {
		logrus.Errorf("Error stopping metrics collector: %v", err)
	}

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

func setupLogging(level string, cfg config.LoggingConfig) {
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	// Set log format
	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output
	if cfg.Output != "" && cfg.Output != "stdout" {
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Warnf("Failed to open log file %s, using stdout: %v", cfg.Output, err)
		} else {
			logrus.SetOutput(file)
		}
	}
}