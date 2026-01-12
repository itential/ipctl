# !make

# Copyright 2024 Itential Inc. All Rights Reserved
# Unauthorized copying of this file, via any medium is strictly prohibited
# Proprietary and confidential

# ============================================================================
# Configuration
# ============================================================================

export GOOS        := $(shell uname | tr '[:upper:]' '[:lower:]')
export GOARCH      := amd64
export CGO_ENABLED := 0

# Build metadata
GIT_COMMIT  := $(shell git rev-parse --short HEAD)
GIT_VERSION := $(shell git tag --sort=-v:refname | head -n 1)

# Directories
BIN_DIR   := bin
DIST_DIR  := dist
COVER_DIR := cover

# Binary name
BINARY := ipctl

# ============================================================================
# Phony targets
# ============================================================================

.PHONY: build \
	clean \
	config \
	coverage \
	dependencies \
	help \
	install \
	licenses \
	snapshot \
	test \
	vulncheck

.DEFAULT_GOAL := help

# ============================================================================
# Targets
# ============================================================================

## build: Build the application binary
build: install
	@echo "Building $(BINARY)..."
	@go build \
		-v \
		-o $(BIN_DIR)/$(BINARY) \
		-ldflags="-X 'github.com/itential/ipctl/internal/app.Build=$(GIT_COMMIT)' \
		          -X 'github.com/itential/ipctl/internal/app.Version=$(GIT_VERSION)'" \
		./cmd/ipctl

## clean: Remove build artifacts and generated files
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) $(DIST_DIR) $(COVER_DIR)

## config: Display the current build configuration
config:
	@echo "Build Configuration:"
	@echo "  GOOS        = $(GOOS)"
	@echo "  GOARCH      = $(GOARCH)"
	@echo "  CGO_ENABLED = $(CGO_ENABLED)"
	@echo "  GIT_COMMIT  = $(GIT_COMMIT)"
	@echo "  GIT_VERSION = $(GIT_VERSION)"
	@echo ""

## coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	@scripts/test.sh coverage

## dependencies: Install required development tools
dependencies:
	@echo "Installing development dependencies..."
	@go install github.com/google/go-licenses/v2@latest

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ": "} \
		/^## / { \
			sub(/^## /, "", $$0); \
			printf "  \033[36m%-14s\033[0m %s\n", $$1, substr($$0, length($$1) + 3) \
		}' $(MAKEFILE_LIST)
	@echo ""

## install: Download and tidy Go module dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

## licenses: Ensure license headers exist and generate NOTICE file
licenses:
	@echo "Checking license headers and generating NOTICE file..."
	@go-licenses report . \
		--template ./scripts/license-attributions/template.tpl \
		--ignore github.com/itential > NOTICE
	@go run ./scripts/copyrighter/main.go

## snapshot: Create a development snapshot build with goreleaser
snapshot:
	@echo "Creating snapshot build..."
	@BUILD=$(GIT_COMMIT) goreleaser release --snapshot --clean

## test: Run the test suite with license checks
test:
	@echo "Running tests..."
	@go run ./scripts/copyrighter/main.go -check
	@scripts/test.sh unittest

## vulncheck: Check for known vulnerabilities in dependencies
vulncheck:
	@echo "Checking for vulnerabilities..."
	@govulncheck ./...
