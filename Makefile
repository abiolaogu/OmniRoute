# OmniRoute Commerce Platform Makefile
# The Commerce Operating System for Emerging Markets

.PHONY: help build test run clean docker-build docker-up docker-down lint fmt

# Default target
help:
	@echo "OmniRoute Commerce Platform - Build Commands"
	@echo "============================================="
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build all services"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start all services"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-logs    - View logs"
	@echo ""
	@echo "Individual Services:"
	@echo "  make run-pricing    - Run pricing engine"
	@echo "  make run-gig        - Run gig platform"
	@echo "  make run-notify     - Run notification service"
	@echo "  make run-payment    - Run payment service"
	@echo ""
	@echo "Database:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo ""

# =============================================================================
# Build Commands
# =============================================================================

build: build-pricing build-gig build-notify build-payment
	@echo "All services built successfully"

build-pricing:
	@echo "Building pricing engine..."
	cd services/pricing-engine && go build -o ../../bin/pricing-engine ./cmd/server

build-gig:
	@echo "Building gig platform..."
	cd services/gig-platform && go build -o ../../bin/gig-platform ./cmd/server

build-notify:
	@echo "Building notification service..."
	cd services/notification-service && go build -o ../../bin/notification-service ./cmd/server

build-payment:
	@echo "Building payment service..."
	cd services/payment-service && go build -o ../../bin/payment-service ./cmd/server

# =============================================================================
# Test Commands
# =============================================================================

test: test-pricing test-gig test-notify test-payment
	@echo "All tests passed"

test-pricing:
	@echo "Testing pricing engine..."
	cd services/pricing-engine && go test -v ./...

test-gig:
	@echo "Testing gig platform..."
	cd services/gig-platform && go test -v ./...

test-notify:
	@echo "Testing notification service..."
	cd services/notification-service && go test -v ./...

test-payment:
	@echo "Testing payment service..."
	cd services/payment-service && go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		echo "Testing services/$$service..."; \
		cd services/$$service && go test -coverprofile=coverage.out ./... && cd ../..; \
	done

test-integration:
	@echo "Running integration tests..."
	docker-compose -f docker-compose.test.yml up -d
	go test -v -tags=integration ./...
	docker-compose -f docker-compose.test.yml down

# =============================================================================
# Run Commands
# =============================================================================

run-pricing:
	cd services/pricing-engine && go run ./cmd/server

run-gig:
	cd services/gig-platform && go run ./cmd/server

run-notify:
	cd services/notification-service && go run ./cmd/server

run-payment:
	cd services/payment-service && go run ./cmd/server

# =============================================================================
# Docker Commands
# =============================================================================

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --rmi all

# Infrastructure only
infra-up:
	docker-compose up -d postgres redis kafka minio

infra-down:
	docker-compose stop postgres redis kafka minio

# =============================================================================
# Database Commands
# =============================================================================

migrate:
	@echo "Running database migrations..."
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		echo "Migrating $$service..."; \
		cd services/$$service && go run ./cmd/migrate up && cd ../..; \
	done

migrate-down:
	@echo "Rolling back migrations..."
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		cd services/$$service && go run ./cmd/migrate down && cd ../..; \
	done

migrate-status:
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		echo "Migration status for $$service:"; \
		cd services/$$service && go run ./cmd/migrate status && cd ../..; \
	done

# =============================================================================
# Code Quality
# =============================================================================

lint:
	@echo "Running linters..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

vet:
	@echo "Running go vet..."
	go vet ./...

# =============================================================================
# Clean Commands
# =============================================================================

clean:
	rm -rf bin/
	rm -rf coverage.out
	go clean -cache

clean-all: clean docker-clean
	rm -rf vendor/

# =============================================================================
# Dependencies
# =============================================================================

deps:
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		echo "Getting dependencies for $$service..."; \
		cd services/$$service && go mod download && cd ../..; \
	done

tidy:
	@for service in pricing-engine gig-platform notification-service payment-service; do \
		echo "Tidying $$service..."; \
		cd services/$$service && go mod tidy && cd ../..; \
	done

# =============================================================================
# Code Generation
# =============================================================================

generate:
	go generate ./...

proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go-grpc_out=. proto/*.proto

swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go -d services/pricing-engine -o services/pricing-engine/docs

# =============================================================================
# Release
# =============================================================================

VERSION ?= $(shell git describe --tags --always --dirty)

release: build
	@echo "Creating release $(VERSION)..."
	mkdir -p release/$(VERSION)
	cp bin/* release/$(VERSION)/
	tar -czvf release/omniroute-$(VERSION).tar.gz -C release/$(VERSION) .

# =============================================================================
# Help for Development
# =============================================================================

info:
	@echo "OmniRoute Commerce Platform"
	@echo "==========================="
	@echo "Version: $(VERSION)"
	@echo "Go Version: $(shell go version)"
	@echo ""
	@echo "Services:"
	@echo "  - Pricing Engine:      services/pricing-engine"
	@echo "  - Gig Platform:        services/gig-platform"
	@echo "  - Notification Service: services/notification-service"
	@echo "  - Payment Service:     services/payment-service"
