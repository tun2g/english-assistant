#!/bin/bash

# Deployment script for app-backend
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Check if environment file exists
ENV_FILE="$SCRIPT_DIR/.env.$ENVIRONMENT"
if [[ ! -f "$ENV_FILE" ]]; then
    error "Environment file $ENV_FILE not found!"
    exit 1
fi

# Load environment variables
log "Loading environment variables from $ENV_FILE"
set -a
source "$ENV_FILE"
set +a

# Validate required environment variables
REQUIRED_VARS=("JWT_SECRET" "DB_PASSWORD")
for var in "${REQUIRED_VARS[@]}"; do
    if [[ -z "${!var}" ]]; then
        error "Required environment variable $var is not set!"
        exit 1
    fi
done

# Build and deploy based on deployment type
case "$DEPLOYMENT_TYPE" in
    "docker")
        log "Deploying with Docker Compose..."
        cd "$PROJECT_ROOT/devops/docker"
        
        # Pull latest images
        log "Pulling latest base images..."
        docker-compose -f docker-compose.prod.yml pull postgres redis nginx
        
        # Build application image
        log "Building application image..."
        docker-compose -f docker-compose.prod.yml build app
        
        # Stop existing containers
        log "Stopping existing containers..."
        docker-compose -f docker-compose.prod.yml down
        
        # Start services
        log "Starting services..."
        docker-compose -f docker-compose.prod.yml up -d
        
        # Wait for services to be healthy
        log "Waiting for services to be healthy..."
        timeout 120 bash -c 'until docker-compose -f docker-compose.prod.yml ps | grep -q "healthy"; do sleep 5; done'
        
        success "Docker deployment completed successfully!"
        ;;
        
    "kubernetes")
        log "Deploying to Kubernetes..."
        cd "$PROJECT_ROOT/devops/k8s"
        
        # Apply configurations
        log "Applying Kubernetes configurations..."
        kubectl apply -f namespace.yml
        kubectl apply -f configmap.yml
        kubectl apply -f secrets.yml
        kubectl apply -f postgres.yml
        kubectl apply -f redis.yml
        kubectl apply -f app.yml
        kubectl apply -f nginx.yml
        
        # Wait for rollout
        log "Waiting for deployment rollout..."
        kubectl rollout status deployment/app-backend -n app-backend --timeout=300s
        
        success "Kubernetes deployment completed successfully!"
        ;;
        
    *)
        error "Unknown deployment type: $DEPLOYMENT_TYPE"
        error "Supported types: docker, kubernetes"
        exit 1
        ;;
esac

# Run health check
log "Running health check..."
HEALTH_URL="http://localhost/health"
if curl -f -s "$HEALTH_URL" > /dev/null; then
    success "Health check passed!"
else
    warning "Health check failed. Please check the logs."
fi

# Show service status
case "$DEPLOYMENT_TYPE" in
    "docker")
        log "Service status:"
        docker-compose -f "$PROJECT_ROOT/devops/docker/docker-compose.prod.yml" ps
        ;;
    "kubernetes")
        log "Pod status:"
        kubectl get pods -n app-backend
        ;;
esac

success "Deployment completed!"
log "Application should be available at: https://yourdomain.com"