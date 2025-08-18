module github.com/magic-flow/demo

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.0.5
	github.com/stretchr/testify v1.8.4
	github.com/testcontainers/testcontainers-go v0.24.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
	github.com/gorilla/websocket v1.5.0
	github.com/prometheus/client_golang v1.16.0
	go.uber.org/zap v1.25.0
	github.com/spf13/viper v1.16.0
	github.com/golang-migrate/migrate/v4 v4.16.2
	github.com/shopspring/decimal v1.3.1
)

replace github.com/magic-flow/v2 => ../v2