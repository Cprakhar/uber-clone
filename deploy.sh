#!/bin/bash

# Build and Deploy Uber Clone to Minikube

set -e

echo "🚀 Building and deploying Uber Clone to Minikube..."

# Set Docker environment to minikube
eval $(minikube docker-env)

# Build Docker images
echo "📦 Building Docker images..."

# Build API Gateway
docker build -f services/api-gateway/Dockerfile -t uber-clone/api-gateway:latest .

# Build Trip Service
docker build -f services/trip-service/Dockerfile -t uber-clone/trip-service:latest .

# Change directory to k8s directory
cd deployments/k8s/

# Apply Kubernetes manifests
echo "🎯 Deploying to Kubernetes..."

# Apply namespace and config first
kubectl apply -f namespace.yaml

# Wait a moment for namespace to be ready
sleep 5

# Apply remaining resources (excluding kustomization.yaml)
kubectl apply -f api-gateway.yaml
kubectl apply -f trip-service.yaml
kubectl apply -f mongodb.yaml

# Wait for deployments to be ready
echo "⏳ Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/api-gateway -n uber-clone
kubectl wait --for=condition=available --timeout=300s deployment/trip-service -n uber-clone
kubectl wait --for=condition=available --timeout=300s deployment/mongodb -n uber-clone

# Get the service URLs
echo "🌐 Service URLs:"
echo "API Gateway: http://$(minikube ip):30080"

# Show pod status
echo "📋 Pod Status:"
kubectl get pods -n uber-clone

echo "✅ Deployment complete!"
echo ""
echo "🔗 Access the API Gateway at: http://$(minikube ip):30080"
echo "📊 Monitor with: kubectl get pods -n uber-clone -w"