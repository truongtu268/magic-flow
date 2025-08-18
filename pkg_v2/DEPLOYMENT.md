# Magic Flow v2 - Deployment Guide

## Overview

Magic Flow v2 is designed for easy deployment across various environments, from local development to large-scale production clusters. This guide covers all deployment options and best practices.

## Quick Start Deployment

### Single Binary Deployment

The simplest way to get started with Magic Flow v2:

```bash
# Download the latest release
wget https://github.com/your-org/magic-flow/releases/latest/download/magicflow-v2-linux-amd64.tar.gz
tar -xzf magicflow-v2-linux-amd64.tar.gz
cd magicflow-v2

# Initialize configuration
./magicflow init

# Start the platform
./magicflow start
```

**What's included in the single binary:**
- Workflow Service Platform
- Visual Workflow Designer
- Real-Time Dashboard
- Code Generation Engine
- Execution Engine
- Built-in SQLite database
- Built-in Redis cache
- Web UI and API server

**Default URLs:**
- Main Interface: `http://localhost:9090`
- Dashboard: `http://localhost:9090/dashboard`
- Visual Designer: `http://localhost:9090/designer`
- API Documentation: `http://localhost:9090/docs`

### Configuration Files

After initialization, you'll find these configuration files:

```
magicflow-v2/
├── config/
│   ├── server.yaml          # Server configuration
│   ├── database.yaml        # Database settings
│   ├── cache.yaml          # Cache configuration
│   ├── messaging.yaml      # Message broker settings
│   ├── security.yaml       # Authentication & authorization
│   ├── monitoring.yaml     # Metrics and logging
│   ├── dashboard.yaml      # Dashboard customization
│   └── workflows/          # Workflow definitions
├── data/                   # SQLite database files
├── cache/                  # Redis data files
├── logs/                   # Application logs
└── templates/              # Code generation templates
```

## Docker Deployment

### Single Container

Run Magic Flow v2 in a single Docker container:

```bash
# Pull the latest image
docker pull magicflow/magicflow-v2:latest

# Run with default configuration
docker run -d \
  --name magicflow-v2 \
  -p 9090:9090 \
  -v magicflow-data:/app/data \
  -v magicflow-config:/app/config \
  magicflow/magicflow-v2:latest

# Access the platform
open http://localhost:9090
```

### Docker Compose

For production-like setup with external dependencies:

```yaml
# docker-compose.yml
version: '3.8'

services:
  magicflow:
    image: magicflow/magicflow-v2:latest
    ports:
      - "9090:9090"
    environment:
      - DATABASE_URL=postgresql://postgres:password@postgres:5432/magicflow
      - REDIS_URL=redis://redis:6379
      - KAFKA_BROKERS=kafka:9092
    volumes:
      - ./config:/app/config
      - ./workflows:/app/workflows
      - magicflow-logs:/app/logs
    depends_on:
      - postgres
      - redis
      - kafka
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: magicflow
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped

  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    volumes:
      - kafka-data:/var/lib/kafka/data
    depends_on:
      - zookeeper
    restart: unless-stopped

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - zookeeper-data:/var/lib/zookeeper/data
    restart: unless-stopped

volumes:
  postgres-data:
  redis-data:
  kafka-data:
  zookeeper-data:
  magicflow-logs:

networks:
  default:
    name: magicflow-network
```

```bash
# Start the complete stack
docker-compose up -d

# View logs
docker-compose logs -f magicflow

# Stop the stack
docker-compose down
```

### Custom Docker Image

Build a custom image with your configurations:

```dockerfile
# Dockerfile.custom
FROM magicflow/magicflow-v2:latest

# Copy custom configuration
COPY config/ /app/config/
COPY workflows/ /app/workflows/
COPY templates/ /app/templates/

# Set custom environment variables
ENV MAGICFLOW_ENV=production
ENV LOG_LEVEL=info

# Expose additional ports if needed
EXPOSE 9091 9092

# Custom entrypoint
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
```

```bash
# Build custom image
docker build -f Dockerfile.custom -t my-company/magicflow-v2:latest .

# Run custom image
docker run -d \
  --name my-magicflow \
  -p 9090:9090 \
  my-company/magicflow-v2:latest
```

## Kubernetes Deployment

### Basic Kubernetes Deployment

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: magicflow
  labels:
    name: magicflow
---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: magicflow-config
  namespace: magicflow
data:
  server.yaml: |
    server:
      host: "0.0.0.0"
      port: 9090
      read_timeout: 30s
      write_timeout: 30s
      idle_timeout: 120s
    
  database.yaml: |
    database:
      type: "postgresql"
      host: "postgres-service"
      port: 5432
      name: "magicflow"
      user: "postgres"
      password: "password"
      ssl_mode: "disable"
      max_connections: 100
      max_idle_connections: 10
    
  cache.yaml: |
    cache:
      type: "redis"
      host: "redis-service"
      port: 6379
      database: 0
      max_connections: 100
      max_idle_connections: 10
---
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: magicflow-secrets
  namespace: magicflow
type: Opaque
data:
  database-password: cGFzc3dvcmQ=  # base64 encoded "password"
  jwt-secret: eW91ci1qd3Qtc2VjcmV0LWtleQ==  # base64 encoded secret
  api-key: eW91ci1hcGkta2V5  # base64 encoded API key
---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: magicflow
  namespace: magicflow
  labels:
    app: magicflow
spec:
  replicas: 3
  selector:
    matchLabels:
      app: magicflow
  template:
    metadata:
      labels:
        app: magicflow
    spec:
      containers:
      - name: magicflow
        image: magicflow/magicflow-v2:latest
        ports:
        - containerPort: 9090
          name: http
        env:
        - name: MAGICFLOW_ENV
          value: "production"
        - name: DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: magicflow-secrets
              key: database-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: magicflow-secrets
              key: jwt-secret
        volumeMounts:
        - name: config
          mountPath: /app/config
          readOnly: true
        - name: workflows
          mountPath: /app/workflows
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 9090
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 9090
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: magicflow-config
      - name: workflows
        persistentVolumeClaim:
          claimName: magicflow-workflows-pvc
---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: magicflow-service
  namespace: magicflow
  labels:
    app: magicflow
spec:
  selector:
    app: magicflow
  ports:
  - name: http
    port: 80
    targetPort: 9090
    protocol: TCP
  type: ClusterIP
---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: magicflow-ingress
  namespace: magicflow
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
spec:
  tls:
  - hosts:
    - magicflow.your-domain.com
    secretName: magicflow-tls
  rules:
  - host: magicflow.your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: magicflow-service
            port:
              number: 80
---
# k8s/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: magicflow-workflows-pvc
  namespace: magicflow
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 10Gi
  storageClassName: nfs-client
```

### Deploy to Kubernetes

```bash
# Apply all configurations
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n magicflow
kubectl get services -n magicflow
kubectl get ingress -n magicflow

# View logs
kubectl logs -f deployment/magicflow -n magicflow

# Scale deployment
kubectl scale deployment magicflow --replicas=5 -n magicflow
```

### Helm Chart Deployment

For easier Kubernetes management, use the official Helm chart:

```bash
# Add Helm repository
helm repo add magicflow https://charts.magicflow.io
helm repo update

# Install with default values
helm install magicflow magicflow/magicflow-v2 \
  --namespace magicflow \
  --create-namespace

# Install with custom values
helm install magicflow magicflow/magicflow-v2 \
  --namespace magicflow \
  --create-namespace \
  --values values.yaml
```

**Custom values.yaml:**
```yaml
# values.yaml
replicaCount: 3

image:
  repository: magicflow/magicflow-v2
  tag: "latest"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: magicflow.your-domain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: magicflow-tls
      hosts:
        - magicflow.your-domain.com

postgresql:
  enabled: true
  auth:
    postgresPassword: "secure-password"
    database: "magicflow"
  primary:
    persistence:
      enabled: true
      size: 20Gi

redis:
  enabled: true
  auth:
    enabled: false
  master:
    persistence:
      enabled: true
      size: 5Gi

resources:
  limits:
    cpu: 1000m
    memory: 2Gi
  requests:
    cpu: 250m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

monitoring:
  enabled: true
  serviceMonitor:
    enabled: true
  grafana:
    enabled: true
  prometheus:
    enabled: true
```

## Cloud Provider Deployments

### AWS EKS

```bash
# Create EKS cluster
eksctl create cluster \
  --name magicflow-cluster \
  --region us-west-2 \
  --nodegroup-name magicflow-nodes \
  --node-type m5.large \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 10 \
  --managed

# Install AWS Load Balancer Controller
kubectl apply -k "github.com/aws/eks-charts/stable/aws-load-balancer-controller//crds?ref=master"

# Deploy Magic Flow v2
helm install magicflow magicflow/magicflow-v2 \
  --namespace magicflow \
  --create-namespace \
  --set ingress.className=alb \
  --set ingress.annotations."kubernetes\.io/ingress\.class"=alb \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/scheme"=internet-facing
```

### Google GKE

```bash
# Create GKE cluster
gcloud container clusters create magicflow-cluster \
  --zone us-central1-a \
  --num-nodes 3 \
  --enable-autoscaling \
  --min-nodes 1 \
  --max-nodes 10 \
  --machine-type n1-standard-2

# Get credentials
gcloud container clusters get-credentials magicflow-cluster --zone us-central1-a

# Deploy Magic Flow v2
helm install magicflow magicflow/magicflow-v2 \
  --namespace magicflow \
  --create-namespace \
  --set ingress.className=gce
```

### Azure AKS

```bash
# Create resource group
az group create --name magicflow-rg --location eastus

# Create AKS cluster
az aks create \
  --resource-group magicflow-rg \
  --name magicflow-cluster \
  --node-count 3 \
  --enable-addons monitoring \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --resource-group magicflow-rg --name magicflow-cluster

# Deploy Magic Flow v2
helm install magicflow magicflow/magicflow-v2 \
  --namespace magicflow \
  --create-namespace
```

## Configuration Management

### Environment-Specific Configurations

**Development Environment:**
```yaml
# config/environments/development.yaml
server:
  debug: true
  log_level: debug
  cors:
    enabled: true
    origins: ["*"]

database:
  type: sqlite
  path: "./data/magicflow-dev.db"

cache:
  type: memory

messaging:
  type: memory

security:
  jwt_secret: "dev-secret-key"
  auth_required: false

monitoring:
  metrics_enabled: true
  tracing_enabled: false
```

**Production Environment:**
```yaml
# config/environments/production.yaml
server:
  debug: false
  log_level: info
  cors:
    enabled: true
    origins: ["https://your-domain.com"]

database:
  type: postgresql
  host: "${DATABASE_HOST}"
  port: 5432
  name: "${DATABASE_NAME}"
  user: "${DATABASE_USER}"
  password: "${DATABASE_PASSWORD}"
  ssl_mode: require
  max_connections: 100

cache:
  type: redis
  host: "${REDIS_HOST}"
  port: 6379
  password: "${REDIS_PASSWORD}"
  database: 0

messaging:
  type: kafka
  brokers: ["${KAFKA_BROKERS}"]
  security:
    protocol: SASL_SSL
    mechanism: PLAIN
    username: "${KAFKA_USERNAME}"
    password: "${KAFKA_PASSWORD}"

security:
  jwt_secret: "${JWT_SECRET}"
  auth_required: true
  session_timeout: 24h

monitoring:
  metrics_enabled: true
  tracing_enabled: true
  log_format: json
```

### Configuration Validation

```bash
# Validate configuration
./magicflow config validate --env production

# Test database connection
./magicflow config test-db --env production

# Test cache connection
./magicflow config test-cache --env production

# Test messaging connection
./magicflow config test-messaging --env production
```

## Scaling and Performance

### Horizontal Scaling

**Load Balancer Configuration:**
```yaml
# nginx.conf for load balancing
upstream magicflow_backend {
    least_conn;
    server magicflow-1:9090 max_fails=3 fail_timeout=30s;
    server magicflow-2:9090 max_fails=3 fail_timeout=30s;
    server magicflow-3:9090 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name magicflow.your-domain.com;
    
    location / {
        proxy_pass http://magicflow_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://magicflow_backend/health;
    }
}
```

**Auto-scaling Configuration:**
```yaml
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: magicflow-hpa
  namespace: magicflow
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: magicflow
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: workflow_executions_per_second
      target:
        type: AverageValue
        averageValue: "10"
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
      - type: Pods
        value: 2
        periodSeconds: 60
```

### Vertical Scaling

**Resource Recommendations:**

| Workload Size | CPU | Memory | Storage | Replicas |
|---------------|-----|--------|---------|----------|
| Small (< 100 workflows/day) | 250m | 512Mi | 10Gi | 1-2 |
| Medium (< 1000 workflows/day) | 500m | 1Gi | 50Gi | 2-3 |
| Large (< 10000 workflows/day) | 1000m | 2Gi | 100Gi | 3-5 |
| Enterprise (> 10000 workflows/day) | 2000m | 4Gi | 500Gi | 5-10 |

**Performance Tuning:**
```yaml
# config/performance.yaml
performance:
  # Database connection pool
  database:
    max_connections: 100
    max_idle_connections: 10
    connection_max_lifetime: 1h
    
  # Cache settings
  cache:
    max_connections: 100
    max_idle_connections: 20
    key_prefix: "magicflow:"
    default_ttl: 1h
    
  # Workflow execution
  execution:
    max_concurrent_workflows: 100
    max_concurrent_steps: 500
    default_timeout: 30m
    retry_attempts: 3
    
  # API rate limiting
  rate_limiting:
    enabled: true
    requests_per_minute: 1000
    burst_size: 100
    
  # Background jobs
  background_jobs:
    max_workers: 10
    queue_size: 1000
    batch_size: 50
```

## Monitoring and Observability

### Metrics Collection

**Prometheus Configuration:**
```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'magicflow'
    static_configs:
      - targets: ['magicflow:9090']
    metrics_path: '/metrics'
    scrape_interval: 10s
    
  - job_name: 'magicflow-kubernetes'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
```

**Grafana Dashboard:**
```json
{
  "dashboard": {
    "title": "Magic Flow v2 - System Overview",
    "panels": [
      {
        "title": "Workflow Executions",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(magicflow_workflow_executions_total[5m])",
            "legendFormat": "Executions/sec"
          }
        ]
      },
      {
        "title": "Success Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "(rate(magicflow_workflow_executions_total{status=\"success\"}[5m]) / rate(magicflow_workflow_executions_total[5m])) * 100",
            "legendFormat": "Success %"
          }
        ]
      },
      {
        "title": "Average Execution Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(magicflow_workflow_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### Logging

**Structured Logging Configuration:**
```yaml
# config/logging.yaml
logging:
  level: info
  format: json
  output: stdout
  
  # Log rotation
  rotation:
    enabled: true
    max_size: 100MB
    max_files: 10
    max_age: 30d
    
  # Log shipping
  shipping:
    enabled: true
    type: elasticsearch
    endpoint: "https://elasticsearch.your-domain.com:9200"
    index: "magicflow-logs"
    
  # Log levels by component
  components:
    workflow_engine: debug
    api_server: info
    database: warn
    cache: warn
```

### Distributed Tracing

**Jaeger Configuration:**
```yaml
# config/tracing.yaml
tracing:
  enabled: true
  service_name: "magicflow-v2"
  
  jaeger:
    endpoint: "http://jaeger-collector:14268/api/traces"
    sampler:
      type: probabilistic
      param: 0.1  # Sample 10% of traces
    
  # Trace specific operations
  operations:
    workflow_execution: true
    api_requests: true
    database_queries: false
    cache_operations: false
```

## Security

### SSL/TLS Configuration

```yaml
# config/tls.yaml
tls:
  enabled: true
  cert_file: "/etc/ssl/certs/magicflow.crt"
  key_file: "/etc/ssl/private/magicflow.key"
  ca_file: "/etc/ssl/certs/ca.crt"
  
  # Minimum TLS version
  min_version: "1.2"
  
  # Cipher suites
  cipher_suites:
    - "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
    - "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
    - "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"
```

### Authentication and Authorization

```yaml
# config/auth.yaml
authentication:
  providers:
    - name: "local"
      type: "local"
      enabled: true
      
    - name: "oauth2"
      type: "oauth2"
      enabled: true
      client_id: "your-oauth-client-id"
      client_secret: "your-oauth-client-secret"
      auth_url: "https://auth.your-domain.com/oauth/authorize"
      token_url: "https://auth.your-domain.com/oauth/token"
      user_info_url: "https://auth.your-domain.com/oauth/userinfo"
      
    - name: "ldap"
      type: "ldap"
      enabled: false
      server: "ldap://ldap.your-domain.com:389"
      bind_dn: "cn=admin,dc=your-domain,dc=com"
      bind_password: "admin-password"
      user_base: "ou=users,dc=your-domain,dc=com"
      user_filter: "(uid=%s)"

authorization:
  rbac:
    enabled: true
    roles:
      - name: "admin"
        permissions: ["*"]
      - name: "operator"
        permissions: ["workflow:read", "workflow:execute", "dashboard:read"]
      - name: "viewer"
        permissions: ["workflow:read", "dashboard:read"]
```

### Network Security

```yaml
# k8s/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: magicflow-network-policy
  namespace: magicflow
spec:
  podSelector:
    matchLabels:
      app: magicflow
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    - podSelector:
        matchLabels:
          app: prometheus
    ports:
    - protocol: TCP
      port: 9090
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  - to: []
    ports:
    - protocol: TCP
      port: 53
    - protocol: UDP
      port: 53
```

## Backup and Disaster Recovery

### Database Backup

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/magicflow"
DATABASE_URL="postgresql://user:pass@localhost:5432/magicflow"

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
pg_dump $DATABASE_URL | gzip > $BACKUP_DIR/magicflow_db_$DATE.sql.gz

# Workflow definitions backup
tar -czf $BACKUP_DIR/workflows_$DATE.tar.gz /app/workflows/

# Configuration backup
tar -czf $BACKUP_DIR/config_$DATE.tar.gz /app/config/

# Upload to S3 (optional)
aws s3 cp $BACKUP_DIR/magicflow_db_$DATE.sql.gz s3://your-backup-bucket/magicflow/
aws s3 cp $BACKUP_DIR/workflows_$DATE.tar.gz s3://your-backup-bucket/magicflow/
aws s3 cp $BACKUP_DIR/config_$DATE.tar.gz s3://your-backup-bucket/magicflow/

# Cleanup old backups (keep last 30 days)
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete
```

### Disaster Recovery

```bash
#!/bin/bash
# restore.sh

BACKUP_DATE=$1
BACKUP_DIR="/backups/magicflow"
DATABASE_URL="postgresql://user:pass@localhost:5432/magicflow"

if [ -z "$BACKUP_DATE" ]; then
    echo "Usage: $0 <backup_date>"
    echo "Example: $0 20240115_143000"
    exit 1
fi

# Restore database
echo "Restoring database..."
gunzip -c $BACKUP_DIR/magicflow_db_$BACKUP_DATE.sql.gz | psql $DATABASE_URL

# Restore workflows
echo "Restoring workflows..."
tar -xzf $BACKUP_DIR/workflows_$BACKUP_DATE.tar.gz -C /

# Restore configuration
echo "Restoring configuration..."
tar -xzf $BACKUP_DIR/config_$BACKUP_DATE.tar.gz -C /

echo "Restore completed. Please restart Magic Flow v2."
```

## Troubleshooting

### Common Issues

**1. Database Connection Issues**
```bash
# Check database connectivity
./magicflow config test-db

# Check database logs
kubectl logs deployment/postgres -n magicflow

# Verify database credentials
kubectl get secret magicflow-secrets -n magicflow -o yaml
```

**2. High Memory Usage**
```bash
# Check memory usage
kubectl top pods -n magicflow

# Analyze memory profile
curl http://localhost:9090/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Adjust memory limits
kubectl patch deployment magicflow -n magicflow -p '{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "magicflow",
          "resources": {
            "limits": {
              "memory": "4Gi"
            }
          }
        }]
      }
    }
  }
}'
```

**3. Workflow Execution Failures**
```bash
# Check workflow logs
curl -X GET "http://localhost:9090/api/v1/workflows/executions/{execution_id}/logs"

# Check system metrics
curl -X GET "http://localhost:9090/api/v1/metrics/system"

# Restart failed workflows
curl -X POST "http://localhost:9090/api/v1/workflows/executions/{execution_id}/retry"
```

### Performance Optimization

**Database Optimization:**
```sql
-- Create indexes for better performance
CREATE INDEX CONCURRENTLY idx_workflows_status ON workflows(status);
CREATE INDEX CONCURRENTLY idx_executions_created_at ON executions(created_at);
CREATE INDEX CONCURRENTLY idx_executions_workflow_id ON executions(workflow_id);

-- Analyze query performance
EXPLAIN ANALYZE SELECT * FROM executions WHERE status = 'running';

-- Update table statistics
ANALYZE workflows;
ANALYZE executions;
```

**Cache Optimization:**
```bash
# Monitor Redis performance
redis-cli --latency-history -i 1

# Check cache hit rate
redis-cli info stats | grep keyspace

# Optimize cache configuration
redis-cli config set maxmemory-policy allkeys-lru
redis-cli config set maxmemory 2gb
```

This comprehensive deployment guide covers all aspects of deploying Magic Flow v2 from development to production environments. Choose the deployment method that best fits your infrastructure and requirements.