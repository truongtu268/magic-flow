#!/bin/bash

# Magic Flow v2 Deployment Script
# This script automates the deployment of Magic Flow to different environments

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEPLOYMENT_DIR="$PROJECT_ROOT/deployments"
K8S_DIR="$DEPLOYMENT_DIR/k8s"

# Default values
ENVIRONMENT="development"
NAMESPACE="magic-flow"
IMAGE_TAG="latest"
REGISTRY=""
DRY_RUN=false
VERBOSE=false
SKIP_BUILD=false
SKIP_TESTS=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Help function
show_help() {
    cat << EOF
Magic Flow v2 Deployment Script

Usage: $0 [OPTIONS]

Options:
    -e, --environment ENV    Target environment (development|staging|production) [default: development]
    -n, --namespace NS       Kubernetes namespace [default: magic-flow]
    -t, --tag TAG           Docker image tag [default: latest]
    -r, --registry REG      Docker registry URL
    --dry-run               Show what would be deployed without actually deploying
    --skip-build            Skip Docker image build
    --skip-tests            Skip running tests
    -v, --verbose           Enable verbose output
    -h, --help              Show this help message

Examples:
    $0 --environment production --tag v2.1.0
    $0 --dry-run --verbose
    $0 --skip-build --environment staging

Environments:
    development: Local development with minimal resources
    staging:     Staging environment with production-like setup
    production:  Production environment with full security and monitoring
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
    log_error "Invalid environment: $ENVIRONMENT"
    log_error "Valid environments: development, staging, production"
    exit 1
fi

# Set image name
IMAGE_NAME="magic-flow"
if [[ -n "$REGISTRY" ]]; then
    FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
else
    FULL_IMAGE_NAME="$IMAGE_NAME:$IMAGE_TAG"
fi

log_info "Starting Magic Flow v2 deployment"
log_info "Environment: $ENVIRONMENT"
log_info "Namespace: $NAMESPACE"
log_info "Image: $FULL_IMAGE_NAME"
log_info "Dry run: $DRY_RUN"

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if docker is installed (unless skipping build)
    if [[ "$SKIP_BUILD" == false ]] && ! command -v docker &> /dev/null; then
        log_error "docker is not installed or not in PATH"
        exit 1
    fi
    
    # Check if we can connect to Kubernetes cluster
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster"
        log_error "Please check your kubeconfig and cluster connectivity"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Run tests
run_tests() {
    if [[ "$SKIP_TESTS" == true ]]; then
        log_warning "Skipping tests"
        return
    fi
    
    log_info "Running tests..."
    cd "$PROJECT_ROOT"
    
    if [[ "$DRY_RUN" == false ]]; then
        if ! go test ./...; then
            log_error "Tests failed"
            exit 1
        fi
        log_success "All tests passed"
    else
        log_info "[DRY RUN] Would run: go test ./..."
    fi
}

# Build Docker image
build_image() {
    if [[ "$SKIP_BUILD" == true ]]; then
        log_warning "Skipping Docker image build"
        return
    fi
    
    log_info "Building Docker image: $FULL_IMAGE_NAME"
    cd "$PROJECT_ROOT"
    
    if [[ "$DRY_RUN" == false ]]; then
        if ! docker build -t "$FULL_IMAGE_NAME" .; then
            log_error "Docker build failed"
            exit 1
        fi
        
        # Push to registry if specified
        if [[ -n "$REGISTRY" ]]; then
            log_info "Pushing image to registry..."
            if ! docker push "$FULL_IMAGE_NAME"; then
                log_error "Docker push failed"
                exit 1
            fi
        fi
        
        log_success "Docker image built successfully"
    else
        log_info "[DRY RUN] Would run: docker build -t $FULL_IMAGE_NAME ."
        if [[ -n "$REGISTRY" ]]; then
            log_info "[DRY RUN] Would run: docker push $FULL_IMAGE_NAME"
        fi
    fi
}

# Create namespace
create_namespace() {
    log_info "Creating namespace: $NAMESPACE"
    
    if [[ "$DRY_RUN" == false ]]; then
        if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
            kubectl apply -f "$K8S_DIR/namespace.yaml"
            log_success "Namespace created"
        else
            log_info "Namespace already exists"
        fi
    else
        log_info "[DRY RUN] Would create namespace: $NAMESPACE"
    fi
}

# Deploy PostgreSQL
deploy_postgres() {
    log_info "Deploying PostgreSQL..."
    
    if [[ "$DRY_RUN" == false ]]; then
        kubectl apply -f "$K8S_DIR/postgres.yaml" -n "$NAMESPACE"
        
        # Wait for PostgreSQL to be ready
        log_info "Waiting for PostgreSQL to be ready..."
        kubectl wait --for=condition=ready pod -l app=postgres -n "$NAMESPACE" --timeout=300s
        
        log_success "PostgreSQL deployed successfully"
    else
        log_info "[DRY RUN] Would deploy PostgreSQL"
    fi
}

# Deploy Redis
deploy_redis() {
    log_info "Deploying Redis..."
    
    if [[ "$DRY_RUN" == false ]]; then
        kubectl apply -f "$K8S_DIR/redis.yaml" -n "$NAMESPACE"
        
        # Wait for Redis to be ready
        log_info "Waiting for Redis to be ready..."
        kubectl wait --for=condition=ready pod -l app=redis -n "$NAMESPACE" --timeout=300s
        
        log_success "Redis deployed successfully"
    else
        log_info "[DRY RUN] Would deploy Redis"
    fi
}

# Deploy Magic Flow application
deploy_application() {
    log_info "Deploying Magic Flow application..."
    
    # Update image in deployment manifest
    local temp_manifest="/tmp/magic-flow-deployment.yaml"
    sed "s|image: magic-flow:v2|image: $FULL_IMAGE_NAME|g" "$K8S_DIR/magic-flow.yaml" > "$temp_manifest"
    
    if [[ "$DRY_RUN" == false ]]; then
        kubectl apply -f "$temp_manifest" -n "$NAMESPACE"
        
        # Wait for deployment to be ready
        log_info "Waiting for Magic Flow deployment to be ready..."
        kubectl wait --for=condition=available deployment/magic-flow -n "$NAMESPACE" --timeout=600s
        
        # Clean up temp file
        rm -f "$temp_manifest"
        
        log_success "Magic Flow application deployed successfully"
    else
        log_info "[DRY RUN] Would deploy Magic Flow application with image: $FULL_IMAGE_NAME"
        rm -f "$temp_manifest"
    fi
}

# Deploy ingress
deploy_ingress() {
    if [[ "$ENVIRONMENT" == "development" ]]; then
        log_info "Skipping ingress deployment for development environment"
        return
    fi
    
    log_info "Deploying ingress..."
    
    if [[ "$DRY_RUN" == false ]]; then
        kubectl apply -f "$K8S_DIR/ingress.yaml" -n "$NAMESPACE"
        log_success "Ingress deployed successfully"
    else
        log_info "[DRY RUN] Would deploy ingress"
    fi
}

# Verify deployment
verify_deployment() {
    log_info "Verifying deployment..."
    
    if [[ "$DRY_RUN" == false ]]; then
        # Check pod status
        log_info "Pod status:"
        kubectl get pods -n "$NAMESPACE"
        
        # Check service status
        log_info "Service status:"
        kubectl get services -n "$NAMESPACE"
        
        # Check if application is responding
        log_info "Checking application health..."
        if kubectl get pods -l app=magic-flow -n "$NAMESPACE" -o jsonpath='{.items[0].status.phase}' | grep -q "Running"; then
            log_success "Application is running"
            
            # Port forward for health check in development
            if [[ "$ENVIRONMENT" == "development" ]]; then
                log_info "You can access the application at:"
                log_info "  API: kubectl port-forward svc/magic-flow-api-service 8080:8080 -n $NAMESPACE"
                log_info "  Dashboard: kubectl port-forward svc/magic-flow-dashboard-service 8081:8081 -n $NAMESPACE"
            fi
        else
            log_warning "Application may not be fully ready yet"
        fi
    else
        log_info "[DRY RUN] Would verify deployment status"
    fi
}

# Cleanup function
cleanup() {
    log_info "Cleaning up temporary files..."
    # Add any cleanup logic here
}

# Set trap for cleanup
trap cleanup EXIT

# Main deployment flow
main() {
    check_prerequisites
    run_tests
    build_image
    create_namespace
    deploy_postgres
    deploy_redis
    deploy_application
    deploy_ingress
    verify_deployment
    
    log_success "Magic Flow v2 deployment completed successfully!"
    
    if [[ "$ENVIRONMENT" != "development" ]]; then
        log_info "Don't forget to:"
        log_info "  1. Update DNS records to point to your ingress"
        log_info "  2. Configure SSL certificates"
        log_info "  3. Set up monitoring and alerting"
        log_info "  4. Configure backup procedures"
    fi
}

# Run main function
main