module github.com/magic-flow/v2

go 1.21

require (
	// Web framework
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/websocket v1.5.0
	
	// Database
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	gorm.io/driver/mysql v1.5.2
	
	// Redis for caching
	github.com/redis/go-redis/v9 v9.3.0
	
	// Configuration
	github.com/spf13/viper v1.17.0
	github.com/spf13/cobra v1.8.0
	
	// YAML processing
	gopkg.in/yaml.v3 v3.0.1
	
	// Validation
	github.com/go-playground/validator/v10 v10.16.0
	
	// UUID generation
	github.com/google/uuid v1.4.0
	
	// Logging
	github.com/sirupsen/logrus v1.9.3
	
	// Metrics and monitoring
	github.com/prometheus/client_golang v1.17.0
	
	// Template engine for code generation
	text/template
	
	// HTTP client
	github.com/go-resty/resty/v2 v2.10.0
	
	// JWT authentication
	github.com/golang-jwt/jwt/v5 v5.2.0
	
	// Testing
	github.com/stretchr/testify v1.8.4
	github.com/testcontainers/testcontainers-go v0.26.0
	
	// Migration
	github.com/golang-migrate/migrate/v4 v4.16.2
	
	// JSON processing
	github.com/tidwall/gjson v1.17.0
	github.com/tidwall/sjson v1.2.5
	
	// Context and cancellation
	context
	time
	fmt
	errors
	sync
	os
	path/filepath
	strings
	bytes
	io
	net/http
	encoding/json
	regexp
)