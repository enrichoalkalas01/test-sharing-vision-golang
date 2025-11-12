# Makefile

# Development with hot-reload
dev:
	@air

# Production build (disable CGO)
build:
	@set CGO_ENABLED=0 && go build -o bin/server.exe cmd/api/server.go

# Run without hot-reload (disable CGO)
start:
	@set CGO_ENABLED=0 && go run cmd/api/server.go

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@if exist tmp rmdir /s /q tmp
	@if exist bin rmdir /s /q bin

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Format code
fmt:
	@go fmt ./...

# Run linter
lint:
	@golangci-lint run

.PHONY: dev build start test clean deps fmt lint