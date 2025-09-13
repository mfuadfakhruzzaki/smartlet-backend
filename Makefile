# Swiflet Backend Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=swiflet-backend
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server

# Build the application
build:
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)

# Run the application
run:
	$(GOCMD) run $(MAIN_PATH)/main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
test:
	$(GOTEST) -v ./...

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH)-linux $(MAIN_PATH)

# Build for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH).exe $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Database migration (requires migrate tool)
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/swiflet_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/swiflet_db?sslmode=disable" down

# Docker commands
docker-build:
	docker build -t swiflet-backend .

docker-run:
	docker run -p 8080:8080 swiflet-backend

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  test          - Run tests"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  dev           - Run with hot reload"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help"

.PHONY: build run clean deps test build-linux build-windows dev fmt lint migrate-up migrate-down docker-build docker-run help