# Magic Flow v2

A powerful workflow automation platform built with Go, featuring a modern dashboard, code generation capabilities, and comprehensive API management.

## Features

- **Workflow Engine**: Execute complex workflows with YAML-based definitions
- **REST API**: Comprehensive API for workflow management and execution
- **Real-time Dashboard**: Monitor workflows, executions, and system metrics
- **Code Generation**: Generate client libraries in multiple languages (Go, TypeScript, Python, Java)
- **Versioning System**: Manage workflow versions with migration and rollback capabilities
- **Configuration Management**: Environment-specific configurations with hot reloading
- **Security**: JWT authentication, rate limiting, and CORS support
- **Monitoring**: Prometheus metrics and health checks
- **Containerized Deployment**: Docker and Kubernetes support

## Quick Start

### Prerequisites

- Go 1.21 or later
- PostgreSQL 15+
- Redis 7+
- Docker (for containerized deployment)
- Kubernetes (for production deployment)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/magic-flow.git
   cd magic-flow/v2
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   # Start PostgreSQL and Redis using Docker Compose
   docker-compose up -d postgres redis
   
   # Run database migrations
   go run cmd/migrate/main.go up
   ```

4. **Configure the application**
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   # Edit configs/config.yaml with your settings
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

6. **Access the application**
   - API: http://localhost:8080
   - Dashboard: http://localhost:8081
   - Health Check: http://localhost:8080/health
   - Metrics: http://localhost:9090/metrics

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **Access the application**
   - API: http://localhost:8080
   - Dashboard: http://localhost:8081
   - Grafana: http://localhost:3000 (admin/admin)
   - Prometheus: http://localhost:9090

### Kubernetes Deployment

1. **Deploy using the deployment script**
   ```bash
   # Development environment
   ./scripts/deploy.sh --environment development
   
   # Production environment
   ./scripts/deploy.sh --environment production --tag v2.1.0 --registry your-registry.com
   ```

2. **Manual deployment**
   ```bash
   # Apply Kubernetes manifests
   kubectl apply -f deployments/k8s/namespace.yaml
   kubectl apply -f deployments/k8s/postgres.yaml
   kubectl apply -f deployments/k8s/redis.yaml
   kubectl apply -f deployments/k8s/magic-flow.yaml
   kubectl apply -f deployments/k8s/ingress.yaml
   ```

## API Documentation

The API documentation is available at `/api/docs` when the server is running, or you can view the [API.md](API.md) file.

### Key Endpoints

- **Workflows**: `/api/v1/workflows`
- **Executions**: `/api/v1/executions`
- **Code Generation**: `/api/v1/codegen`
- **Metrics**: `/api/v1/metrics`
- **Dashboard**: `/api/v1/dashboard`
- **Versioning**: `/api/v1/versions`
- **Configuration**: `/api/v1/config`

## Configuration

The application uses YAML configuration files located in the `configs/` directory. Key configuration sections:

### Server Configuration
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  tls:
    enabled: false
    cert_file: ""
    key_file: ""
```

### Database Configuration
```yaml
database:
  host: "localhost"
  port: 5432
  name: "magic_flow"
  user: "magic_flow"
  password: "password"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
```

### Feature Flags
```yaml
features:
  dashboard: true
  code_generation: true
  versioning: true
  metrics: true
  authentication: true
  rate_limiting: true
```

## Development

### Project Structure

```
v2/
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── api/              # API handlers and routes
│   ├── config/           # Configuration management
│   ├── dashboard/        # Dashboard backend
│   ├── database/         # Database layer
│   ├── engine/           # Workflow execution engine
│   ├── models/           # Data models
│   ├── codegen/          # Code generation engine
│   └── versioning/       # Version management
├── configs/              # Configuration files
├── templates/            # Code generation templates
├── migrations/           # Database migrations
├── deployments/          # Deployment configurations
│   ├── k8s/             # Kubernetes manifests
│   ├── nginx/           # Nginx configurations
│   └── grafana/         # Grafana dashboards
└── scripts/              # Deployment and utility scripts
```

### Building

```bash
# Build the server
go build -o bin/magic-flow cmd/server/main.go

# Build with Docker
docker build -t magic-flow:latest .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/magic-flow-linux-amd64 cmd/server/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test package
go test ./internal/engine/...
```

### Code Generation

Generate client libraries for different languages:

```bash
# Generate Go client
curl -X POST http://localhost:8080/api/v1/codegen/generate \
  -H "Content-Type: application/json" \
  -d '{"language": "go", "package_name": "magicflow"}'

# Generate TypeScript client
curl -X POST http://localhost:8080/api/v1/codegen/generate \
  -H "Content-Type: application/json" \
  -d '{"language": "typescript", "package_name": "magic-flow-client"}'
```

## Monitoring and Observability

### Metrics

The application exposes Prometheus metrics at `/metrics`. Key metrics include:

- `magic_flow_workflows_total`: Total number of workflows
- `magic_flow_executions_total`: Total number of executions
- `magic_flow_execution_duration_seconds`: Execution duration histogram
- `magic_flow_api_requests_total`: API request counter
- `magic_flow_database_connections`: Database connection pool metrics

### Health Checks

- **Liveness**: `/health` - Basic application health
- **Readiness**: `/ready` - Application readiness for traffic
- **Detailed**: `/health/detailed` - Comprehensive health status

### Logging

The application uses structured JSON logging. Log levels can be configured:

- `debug`: Detailed debugging information
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors that cause application exit

## Security

### Authentication

The application supports JWT-based authentication:

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'

# Use token in subsequent requests
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/workflows
```

### Rate Limiting

API endpoints are protected by rate limiting (configurable):

- Default: 100 requests per minute per IP
- Burst: 200 requests
- Configurable per endpoint

### CORS

Cross-Origin Resource Sharing is configurable:

```yaml
security:
  cors:
    enabled: true
    allowed_origins: ["https://your-frontend.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["*"]
```

## Deployment Environments

### Development
- Minimal resource requirements
- Debug logging enabled
- Hot reloading for configuration
- Local file storage

### Staging
- Production-like setup
- Reduced resource allocation
- Integration testing environment
- External database connections

### Production
- Full security features enabled
- High availability setup
- Monitoring and alerting
- Backup and disaster recovery

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   ```bash
   # Check database connectivity
   psql -h localhost -p 5432 -U magic_flow -d magic_flow
   
   # Verify connection string in config
   grep -A 10 "database:" configs/config.yaml
   ```

2. **Redis Connection Issues**
   ```bash
   # Test Redis connectivity
   redis-cli -h localhost -p 6379 ping
   
   # Check Redis logs
   docker logs magic-flow-redis
   ```

3. **Application Won't Start**
   ```bash
   # Check application logs
   docker logs magic-flow-app
   
   # Verify configuration
   go run cmd/server/main.go --config configs/config.yaml --validate
   ```

4. **Kubernetes Deployment Issues**
   ```bash
   # Check pod status
   kubectl get pods -n magic-flow
   
   # View pod logs
   kubectl logs -f deployment/magic-flow -n magic-flow
   
   # Describe problematic resources
   kubectl describe pod <pod-name> -n magic-flow
   ```

### Performance Tuning

1. **Database Optimization**
   - Adjust connection pool settings
   - Monitor slow queries
   - Optimize indexes

2. **Memory Usage**
   - Configure Go garbage collector
   - Monitor heap usage
   - Adjust container limits

3. **Concurrency**
   - Tune max concurrent workflows
   - Adjust worker pool sizes
   - Monitor goroutine counts

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for API changes
- Use conventional commit messages
- Ensure all CI checks pass

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [API.md](API.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/magic-flow/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/magic-flow/discussions)
- **Email**: support@magic-flow.com

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.