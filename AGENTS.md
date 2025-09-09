# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`ipctl` is a CLI tool for managing Itential Platform servers. It's written in Go and provides a command-line interface for working with various Itential Platform resources like projects, profiles, workflows, automations, and assets.

## Architecture

- **Entry Point**: `main.go` → `cmd/root.go` → command execution
- **Command Structure**: Uses Cobra CLI framework with commands organized by resource type
- **Configuration**: Profile-based configuration in `~/.platform.d/config` with environment variable overrides
- **Client**: HTTP client in `pkg/client/` handles API communication with Itential Platform servers
- **Handlers**: Located in `internal/handlers/` - coordinate between CLI commands and service layer
- **Services**: In `pkg/services/` - business logic for interacting with platform APIs
- **Runners**: In `internal/runners/` - execute specific operations (import/export/CRUD)
- **Flags**: In `internal/flags/` - command-line argument parsing and validation

## Key Components

### Resource Types
The CLI manages these primary resources:
- **Assets**: Projects, workflows, automations, templates, transformations
- **Platform**: Accounts, groups, roles, profiles, adapters, integrations
- **Datasets**: Data management (when feature flag enabled)

### Configuration System
- Profile-based configuration supporting multiple Itential Platform instances
- Configuration precedence: CLI flags → environment variables → config file
- Profiles stored in `~/.platform.d/config` with INI format

### Command Descriptors
YAML descriptors in `cmd/descriptors/` and `internal/handlers/descriptors/` define command structure, validation rules, and API mappings.

## Development Commands

### Build and Test
```bash
# Build the application
make build

# Run unit tests with formatting and vetting
make test

# Run tests with coverage report
make coverage

# Install dependencies
make install

# Clean build artifacts
make clean
```

### Test Scripts
```bash
# Run unit tests
scripts/test.sh unittest

# Run tests with coverage
scripts/test.sh coverage

# Debug tests
scripts/test.sh debugtest
```

### Development Build
```bash
# Create snapshot build (requires goreleaser)
make snapshot
```

## Code Patterns

### Command Structure
Commands follow a consistent pattern:
1. Flag definition and parsing in `internal/flags/`
2. Handler coordination in `internal/handlers/`
3. Business logic in `pkg/services/`
4. API operations via `pkg/client/`
5. Output formatting via `internal/terminal/`

### Error Handling
- Use `cmdutils.CheckError()` for CLI error handling
- Service layer returns structured errors
- Client handles HTTP-specific errors and retries

### Testing
- Unit tests use mock data in `testdata/` directories
- Test files follow `*_test.go` naming convention
- Coverage reports generated in `cover/` directory

## Important Files

- `cmd/root.go`: CLI application entry point and command tree construction
- `pkg/config/config.go`: Configuration management and profile handling
- `pkg/client/client.go`: HTTP client for Itential Platform APIs
- `internal/handlers/registry.go`: Command registration and routing
- `Makefile`: Build, test, and development commands