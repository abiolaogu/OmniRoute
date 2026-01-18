# OmniRoute - Harvester HCI Deployment
## Kubernetes Manifests for SUSE Harvester

---

## Overview

This directory contains Kubernetes manifests optimized for deployment on SUSE Harvester HCI (Hyper-Converged Infrastructure). The manifests leverage Harvester's VM and container orchestration capabilities.

## Prerequisites

- SUSE Harvester 1.3.x or later
- kubectl configured for Harvester cluster
- StorageClass with Longhorn provisioner
- LoadBalancer (MetalLB or Harvester's built-in)

---

## Quick Deploy

```bash
# Create namespace
kubectl create namespace omniroute

# Apply secrets
kubectl apply -f secrets.yaml -n omniroute

# Deploy database layer
kubectl apply -f database/ -n omniroute

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=yugabytedb -n omniroute --timeout=300s

# Deploy services
kubectl apply -f services/ -n omniroute

# Deploy ingress
kubectl apply -f ingress.yaml -n omniroute
```

---

## Directory Structure

```
infrastructure/kubernetes/harvester/
├── README.md
├── namespace.yaml
├── secrets.yaml.template
├── storage/
│   ├── storageclass.yaml       # Longhorn storage class
│   └── pvc.yaml                # Persistent volume claims
├── database/
│   ├── yugabytedb.yaml         # Distributed SQL
│   ├── dragonflydb.yaml        # High-performance cache
│   └── redpanda.yaml           # Event streaming
├── services/
│   ├── auth-service.yaml
│   ├── order-service.yaml
│   ├── payment-service.yaml
│   ├── inventory-service.yaml
│   ├── catalog-service.yaml
│   ├── customer-service.yaml
│   ├── pricing-engine.yaml
│   ├── gig-platform.yaml
│   ├── route-optimizer.yaml
│   ├── wms-service.yaml
│   ├── fleet-service.yaml
│   ├── analytics-service.yaml
│   ├── credit-scoring.yaml
│   ├── forecasting.yaml
│   ├── market-intel.yaml
│   ├── recommendations.yaml
│   ├── notification.yaml
│   ├── atc-service.yaml
│   ├── sce-service.yaml
│   ├── mcp-server.yaml
│   ├── bank-gateway.yaml
│   ├── workflow-compiler.yaml
│   └── ai-gateway.yaml
├── orchestration/
│   ├── temporal.yaml           # Workflow engine
│   ├── n8n.yaml                # Automation
│   └── hasura.yaml             # GraphQL API
├── observability/
│   ├── prometheus.yaml
│   ├── grafana.yaml
│   └── loki.yaml
└── ingress.yaml
```
