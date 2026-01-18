# OmniRoute - Deployment Guide

## Overview

This guide covers deployment of OmniRoute to production environments using Kubernetes.

---

## Deployment Environments

| Environment | Purpose | Infrastructure |
|-------------|---------|----------------|
| **Development** | Local development | Docker Compose |
| **Staging** | Pre-production testing | Kubernetes (single region) |
| **Production** | Live environment | Kubernetes (multi-region) |

---

## Prerequisites

### Tools

```bash
# Kubernetes CLI
brew install kubectl

# Helm package manager
brew install helm

# Terraform (IaC)
brew install terraform
```

### Access

- Kubernetes cluster admin access
- Docker registry credentials
- Cloud provider credentials (AWS/GCP/Azure)

---

## Docker Images

### Build Images

```bash
# Build all services
make docker-build-all

# Build specific service
make docker-build SERVICE=pricing-engine

# With version tag
make docker-build SERVICE=pricing-engine VERSION=1.0.0
```

### Push to Registry

```bash
# Login to registry
docker login registry.omniroute.io

# Push all images
make docker-push-all

# Push specific image
make docker-push SERVICE=pricing-engine VERSION=1.0.0
```

---

## Kubernetes Deployment

### Namespace Setup

```yaml
# namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: omniroute
  labels:
    environment: production
---
apiVersion: v1
kind: Namespace
metadata:
  name: omniroute-database
---
apiVersion: v1
kind: Namespace
metadata:
  name: omniroute-observability
```

```bash
kubectl apply -f infrastructure/kubernetes/namespaces.yaml
```

### Secrets Management

```bash
# Create secrets from .env file
kubectl create secret generic omniroute-secrets \
  --from-env-file=.env.production \
  -n omniroute

# Or use external secrets operator
kubectl apply -f infrastructure/kubernetes/external-secrets.yaml
```

### Deploy Infrastructure

```bash
# YugabyteDB
helm install yugabyte yugabytedb/yugabyte \
  --namespace omniroute-database \
  --values infrastructure/helm/yugabyte-values.yaml

# DragonflyDB
kubectl apply -f infrastructure/kubernetes/dragonfly.yaml

# Redpanda
helm install redpanda redpanda/redpanda \
  --namespace omniroute-messaging \
  --values infrastructure/helm/redpanda-values.yaml

# Hasura
kubectl apply -f infrastructure/kubernetes/hasura.yaml

# Temporal
helm install temporal temporal/temporal \
  --namespace omniroute \
  --values infrastructure/helm/temporal-values.yaml
```

### Deploy Services

```bash
# Deploy all services
kubectl apply -f infrastructure/kubernetes/services/

# Or deploy individually
kubectl apply -f infrastructure/kubernetes/services/pricing-engine.yaml
kubectl apply -f infrastructure/kubernetes/services/gig-platform.yaml
kubectl apply -f infrastructure/kubernetes/services/payment-service.yaml
# ... etc
```

### Service Manifest Example

```yaml
# pricing-engine.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pricing-engine
  namespace: omniroute
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pricing-engine
  template:
    metadata:
      labels:
        app: pricing-engine
    spec:
      containers:
      - name: pricing-engine
        image: registry.omniroute.io/pricing-engine:1.0.0
        ports:
        - containerPort: 8081
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: omniroute-secrets
              key: DATABASE_URL
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: omniroute-secrets
              key: REDIS_URL
        resources:
          requests:
            cpu: "200m"
            memory: "256Mi"
          limits:
            cpu: "1000m"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: pricing-engine
  namespace: omniroute
spec:
  selector:
    app: pricing-engine
  ports:
  - port: 8081
    targetPort: 8081
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: pricing-engine-hpa
  namespace: omniroute
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: pricing-engine
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## Ingress Configuration

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: omniroute-ingress
  namespace: omniroute
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.omniroute.io
    secretName: omniroute-tls
  rules:
  - host: api.omniroute.io
    http:
      paths:
      - path: /graphql
        pathType: Prefix
        backend:
          service:
            name: hasura
            port:
              number: 8080
      - path: /v1
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 8080
```

---

## Database Migrations

```bash
# Run migrations (via Job)
kubectl create job --from=cronjob/migrate-job migrate-manual -n omniroute

# Check status
kubectl get jobs -n omniroute
kubectl logs job/migrate-manual -n omniroute
```

---

## Monitoring & Alerts

### Prometheus Stack

```bash
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace omniroute-observability \
  --values infrastructure/helm/prometheus-values.yaml
```

### Alerting Rules

```yaml
# alerts.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: omniroute-alerts
spec:
  groups:
  - name: omniroute
    rules:
    - alert: ServiceDown
      expr: up{job=~"omniroute-.*"} == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Service {{ $labels.job }} is down"
    - alert: HighLatency
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High latency on {{ $labels.service }}"
```

---

## Rollback Procedures

```bash
# View deployment history
kubectl rollout history deployment/pricing-engine -n omniroute

# Rollback to previous version
kubectl rollout undo deployment/pricing-engine -n omniroute

# Rollback to specific revision
kubectl rollout undo deployment/pricing-engine --to-revision=2 -n omniroute
```

---

## Disaster Recovery

### Database Backup

```bash
# Via CronJob (automated)
kubectl apply -f infrastructure/kubernetes/backup-cronjob.yaml

# Manual backup
kubectl exec -it yugabyte-0 -n omniroute-database -- \
  ysql_dump -h localhost -U yugabyte omniroute > backup.sql
```

### Restore

```bash
# Stop services
kubectl scale deployment --replicas=0 -l app.kubernetes.io/part-of=omniroute

# Restore database
kubectl exec -i yugabyte-0 -n omniroute-database -- \
  ysqlsh -h localhost -U yugabyte omniroute < backup.sql

# Restart services
kubectl scale deployment --replicas=3 -l app.kubernetes.io/part-of=omniroute
```

---

## Health Checks

```bash
# Check all pods
kubectl get pods -n omniroute

# Check service endpoints
kubectl get endpoints -n omniroute

# Check ingress
kubectl describe ingress omniroute-ingress -n omniroute

# Check logs
kubectl logs -f deployment/pricing-engine -n omniroute
```
