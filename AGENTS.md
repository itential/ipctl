# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`ipctl` is a Go CLI tool for managing Itential Platform servers. It provides a command-line interface for working with 35+ resource types including projects, workflows, automations, adapters, accounts, groups, and more. The tool uses a layered architecture with clean separation between CLI concerns, business logic, and API communication.

## What This Does

- **Resource Management**: CRUD operations (get, create, update, delete) for Itential Platform resources
- **Import/Export**: Git-based asset import/export with SSH authentication
- **Profile Management**: Multiple server configurations with precedence-based settings
- **Local AAA**: MongoDB-backed local authentication/authorization for development
- **Asset Orchestration**: Complex operations like copying automations between environments

## Architecture & Request Flow

```
cmd/ipctl/main.go
  ↓
internal/cli/root.go (Execute)
  ├── Load config (internal/config + internal/profile)
  ├── Create HTTP client (pkg/client)
  ├── Load YAML descriptors (internal/cli/descriptors/*.yaml)
  └── Build Cobra command tree via loadCommands()
    ↓
internal/handlers/{resource}.go
  ├── NewHandler() creates 35 resource handlers
  ├── Registry uses interface-based discovery (Reader, Writer, Copier, etc.)
  └── AssetHandler generates Cobra commands with flags
    ↓
internal/handlers/runner.go (CommandRunner)
  ├── Parse flags from cobra.Command
  ├── Create runners.Request{Args, Common, Options, Config}
  └── Call runner method (Get/Create/Delete/Import/Export)
    ↓
internal/runners/{resource}.go
  ├── Parse common flags and options
  ├── Call resource methods for business logic
  └── Return runners.Response{Object, Text, Template, Keys}
    ↓
pkg/resources/{resource}.go
  ├── Apply client-side filtering (GetByName)
  ├── Business logic and orchestration
  └── Call service methods
    ↓
pkg/services/{resource}.go
  ├── Pure API operations (GET/POST/PUT/DELETE)
  ├── Handle pagination and transformation
  └── Return structured data or errors
    ↓
pkg/client/client.go (HttpClient)
  ├── Build URL with scheme, host, port, path, query params
  ├── Handle authentication (OAuth2 client credentials or basic auth)
  ├── Send HTTP request with timeout and context
  └── Return response or error
```

**Key Design Principles**:
- Each layer has a single responsibility
- Dependencies flow downward (no circular dependencies)
- Interfaces enable testing with mocks
- Business logic lives in resources, not services or runners

## Package Organization

### Core Packages (internal/)

| Package | Lines | What It Does | Issues |
|---------|-------|--------------|--------|
| **cli/** | ~400 | Cobra setup, descriptor loading, command tree construction | None |
| **handlers/** | ~2,156 | Handler registration, command generation, runtime management | **NO TESTS** (0 _test.go files) |
| **runners/** | ~10,000 | CLI orchestration, flag parsing, Git integration, response formatting | Some business logic leakage |
| **flags/** | ~4,000 | Flag definitions for all commands, validation | Well structured |
| **config/** | ~500 | Configuration loading, feature flags | Clean design |
| **profile/** | ~400 | Profile definition, profile manager, precedence-based loading | Clean design |
| **logging/** | ~500 | Zerolog integration, configurable levels | Clean, no file logging |
| **terminal/** | ~1,500 | Output formatting (JSON/YAML/table/human), color support | Well done |
| **cmdutils/** | ~350 | Descriptor loading, error handling | Has 1 FIXME |
| **utils/** | ~200 | Encoding/type conversion | **CRITICAL: UnmarshalData() calls fatal()** |

### API Packages (pkg/)

| Package | Lines | What It Does | Issues |
|---------|-------|--------------|--------|
| **client/** | ~15,000 | HTTP client, OAuth2/basic auth, cookies, TLS, timeouts | Excellent implementation |
| **services/** | ~3,500 | Pure API operations, pagination, error handling | Has deprecated methods |
| **resources/** | ~2,500 | Business logic, filtering, validation, orchestration | Good separation |
| **validators/** | ~300 | Input validation for complex operations | Limited coverage |
| **localaaa/** | ~800 | MongoDB-backed local AAA | **Uses context.TODO()** |
| **repositories/** | ~400 | Git operations (clone, commit, push) | Well implemented |
| **editor/** | ~100 | File editing via $EDITOR | Simple, clean |

## What's Well Done

### 1. Architecture & Layering ✓
- Clean separation: CLI → Handler → Runner → Resource → Service → Client
- Each layer tested independently with mocked dependencies
- No circular dependencies
- Interface-based design enables testability

### 2. Handler Discovery Pattern ✓
```go
// interfaces.go defines capabilities
type Reader interface {
    Get(*Runtime) *cobra.Command
    Describe(*Runtime) *cobra.Command
}

type Writer interface {
    Create(*Runtime) *cobra.Command
    Delete(*Runtime) *cobra.Command
}

// Registry uses type assertions for auto-discovery
func NewRegistry(handlers []Handler) *Registry {
    for _, h := range handlers {
        if reader, ok := h.(Reader); ok {
            registry.readers = append(registry.readers, reader)
        }
        // ... similar for Writer, Copier, etc.
    }
}
```
**Why it's good**: Adding a new handler requires zero registration boilerplate. Just implement the interfaces you need.

### 3. Configuration System ✓
- Profile-based multi-instance support
- Precedence: CLI flags > env vars > config file > defaults
- Type-safe accessors (getString, getInt, getBool)
- Clean separation: Config (app settings), Profile (server connection), Features (flags)

### 4. HTTP Client ✓
- Comprehensive authentication (OAuth2 client credentials, basic auth)
- Custom cookie jar for session management
- Context-aware with timeout and cancellation
- TLS with optional certificate verification
- Proper error handling with status codes

### 5. YAML Descriptors ✓
```yaml
# internal/cli/descriptors/asset.yaml
get:
  use: projects
  description: Get projects from the server
  example: |
    ipctl get projects
  include_groups: true
```
**Why it's good**: Command metadata lives separate from code. Easy to update help text without code changes.

### 6. Request/Response Pattern ✓
```go
type Request struct {
    Args    []string  // Positional arguments
    Common  any       // Shared flags
    Options any       // Operation-specific flags
    Runner  Runner    // The runner instance
    Config  Provider  // Configuration
}

type Response struct {
    Object   any      // For JSON/YAML
    Text     string   // For human output
    Template string   // Custom formatting
    Keys     []string // Table columns
}
```
**Why it's good**: Consistent interface across all 35+ handlers. Output format handled uniformly.

### 7. Git Integration ✓
- Full Git support (clone, commit, push)
- SSH key authentication
- Reference selection (branch/tag)
- Proper cleanup of temp clones
- Used for import/export workflows

### 8. Package Documentation ✓
- Every major package has doc.go with comprehensive documentation
- Explains architecture, patterns, interfaces, usage
- Examples included

## Technical Debt & Problems

### CRITICAL Issues (Fix Immediately)

#### 1. No Handler Tests
- **Problem**: internal/handlers/ has 38 .go files, ZERO _test.go files
- **Impact**: 2,156 lines of untested code handling command generation and dispatch
- **Location**: `internal/handlers/` - completely untested
- **What to test**: AssetHandler command generation, Registry interface dispatch, flag binding

#### 2. context.TODO() in Production Code
- **Problem**: MongoDB operations use context.TODO() instead of proper cancellation
- **Locations**:
  - `pkg/localaaa/accounts.go` (6 uses)
  - `pkg/localaaa/groups.go` (4 uses)
  - `pkg/localaaa/localaaa.go` (1 use)
- **Impact**: Operations can hang indefinitely, no timeout control
- **Fix**: Pass context through function signatures, use request context

#### 3. UnmarshalData() Calls Fatal in Library Code
- **Problem**: `internal/utils/encoding.go` line 27 calls `logging.Fatal()` on unmarshal error
- **Impact**: Untestable, exits entire process on error instead of returning error
- **FIXME Comment**: "This function should be refactored to return an error"
- **Fix**: Return error, update all callers

#### 4. Deprecated Service Methods Still Used
- **Problem**: Methods marked DEPRECATED still called throughout codebase
- **Examples**:
  - `pkg/services/workflows.go`: GetById(), Clear()
  - `pkg/services/transformations.go`: GetByName(), Clear()
  - `pkg/services/groups.go`: GetByName()
  - `pkg/services/adapters.go`: Restart()
- **Message**: "Business logic method - prefer using resources.XyzResource.GetByName"
- **Impact**: Business logic leaking into service layer, violates layering
- **Fix**: Move to resources layer, remove from services

### Design Issues (Refactor When Possible)

#### 5. Business Logic in Runners
- **Problem**: `internal/runners/automations.go` has `updateTriggers()` with complex type routing
- **Lines**: Lines 280+ - JSON unmarshaling, type detection, transformation
- **Issue**: This is business logic, not CLI concern
- **Fix**: Move to resources/automations.go

#### 6. Error Parsing Complexity
- **Problem**: `automations.go` has `formatImportErrorMessage()` parsing JSON error responses
- **Lines**: ~30 lines (373-406) of brittle error parsing
- **Issue**: Depends on exact response structure, duplicated logic
- **Fix**: Move to service/resource layer for reuse

#### 7. GetByName() Duplication
- **Problem**: Every resource implements its own GetByName() with similar logic
- **Locations**: accounts.go, adapters.go, groups.go, models.go, templates.go, transformations.go, workflows.go
- **Fix**: Create generic helper: `FindByName[T](items []T, name string, extractor func(T) string) (T, error)`

#### 8. handlers.go Size (336 lines)
- **Content**: Runtime struct, Handler struct, NewHandler(), 14 command getter methods
- **Issue**: Single file handling multiple concerns
- **Better**: Split into runtime.go, registry.go (exists!), commands.go

### Missing Functionality

#### 9. No Structured Logging Context
- **Problem**: Minimal context in logs, no request tracing or correlation IDs
- **Impact**: Difficult to debug production issues
- **Fix**: Add request ID to context, structured fields for JSON output

#### 10. No Retry Logic
- **Problem**: Network failures cause immediate exit, no exponential backoff
- **Impact**: Transient failures not handled gracefully
- **Fix**: Implement retry with exponential backoff, circuit breaker

#### 11. Inconsistent Error Messages
- Some: "automation not found"
- Others: "account `%s` does not exist"
- **Fix**: Standardize format and capitalization

#### 12. Type Assertions Without Validation
- **Problem**: Runners cast `in.Options` to specific types without checking
- **Location**: Throughout runners/*.go
- **Impact**: Could panic if wrong type passed
- **Fix**: Add validation or use type parameters

## Refactoring Priorities

### High Priority (Do First)
1. **Fix context.TODO()** (1 hour) - Critical for production stability
2. **Add handler tests** (1-2 days) - 2,156 lines untested
3. **Fix encoding.go UnmarshalData()** (1 hour) - Return error instead of fatal
4. **Remove deprecated service methods** (4 hours) - Move to resources

### Medium Priority (Next Sprint)
5. **Consolidate GetByName()** (1 day) - Generic helper reduces duplication
6. **Move error parsing to service layer** (2 hours) - Better reuse
7. **Add structured logging** (1 day) - Request IDs, timing, JSON fields
8. **Split handlers.go** (2 hours) - Better organization

### Lower Priority (Technical Debt)
9. **Add retry logic** (2 days) - Exponential backoff, circuit breaker
10. **Consolidate magic numbers** (4 hours) - Config constants
11. **Improve error messages** (1 day) - Standardize format
12. **Add input validation** (2 days) - Schema validation layer

## Code Conventions

### File Naming
- `{resource}.go` - Main implementation
- `{resource}_test.go` - Tests with testify/assert
- `descriptors/` - YAML command metadata
- `templates/` - Output templates for describe commands

### Function Naming
- `New{Type}()` - Factory functions (constructors)
- `New{Resource}Handler()` - Handler factories (35 total)
- `New{Resource}Runner()` - Runner factories
- `New{Resource}Service()` - Service factories

### Package Structure Convention
```
pkg/          # Public, reusable, no internal dependencies
  services/   # API operations only, no business logic
  resources/  # Business logic, orchestration
  client/     # HTTP transport

internal/     # Private, CLI-specific
  cli/        # Cobra setup
  handlers/   # Command handlers
  runners/    # CLI orchestration
  flags/      # CLI argument parsing
```

### Interface Pattern
```go
// Define capability interfaces
type Reader interface {
    Get(*Runtime) *cobra.Command
    Describe(*Runtime) *cobra.Command
}

// Handlers implement what they support
type ProjectHandler struct {
    // AssetHandler embeds common functionality
}

func (h *ProjectHandler) Get(rt *Runtime) *cobra.Command {
    // Implementation
}

// Registry discovers capabilities via type assertions
```

### Error Handling Pattern
- **RULE**: Return errors from functions, never call logging.Fatal() except in CLI root
- Services return errors for API failures
- Resources wrap with business context
- Runners propagate to Cobra RunE
- CLI root uses CheckError() for top-level exit

Example:
```go
// GOOD
func DoWork() error {
    if err := doThing(); err != nil {
        return fmt.Errorf("failed to do thing: %w", err)
    }
    return nil
}

// BAD - don't do this
func DoWork() {
    if err := doThing(); err != nil {
        logging.Fatal(err, "failed to do thing")  // WRONG: exits process
    }
}
```

### Testing Pattern
- Mock interfaces, not concrete types
- Use testdata/ for fixtures
- Service tests verify API call structure
- Resource tests verify business logic
- Client tests use mock cookie jar
- **Handler tests**: MISSING - need to add

### Configuration Pattern
- Provider interface for dependency injection
- Loaded at startup via Viper
- Passed to Runtime → handlers → runners
- Precedence: flags > env vars > config file > defaults
- Profiles for multi-instance support

## Key Dependencies

- **Cobra** (v1.10.1): CLI framework for command tree
- **Viper** (v1.19.0): Configuration management (INI/TOML/YAML)
- **Zerolog** (v1.34.0): Structured logging
- **OAuth2** (golang.org/x/oauth2): Client credentials flow
- **go-git** (v5.16.4): Git operations (clone, commit, push)
- **MongoDB Driver** (v1.17.6): Local AAA storage
- **Testify** (v1.11.1): Test assertions

## Development Commands

### Build and Test
```bash
make build      # Build ipctl binary
make test       # Run tests with go fmt and go vet
make coverage   # Generate coverage report in cover/
make install    # Install dependencies
make clean      # Remove build artifacts
make snapshot   # Create goreleaser snapshot build

# Test scripts
scripts/test.sh unittest   # Run unit tests
scripts/test.sh coverage   # Run with coverage
scripts/test.sh debugtest  # Debug tests
```

### Adding a New Resource

1. Create handler in `internal/handlers/{resource}.go`:
```go
func NewMyResourceHandler(runner runners.MyResourceRunner) Handler {
    return &AssetHandler{
        name:   "myresource",
        runner: runner,
    }
}
```

2. Create runner in `internal/runners/{resource}.go`:
```go
type MyResourceRunner struct {
    BaseRunner
    resource resources.MyResourceResource
}

func (r *MyResourceRunner) Get(in *Request) (*Response, error) {
    // Implementation
}
```

3. Create resource in `pkg/resources/{resource}.go`:
```go
type MyResourceResource interface {
    GetByName(name string) (*services.MyResource, error)
}

type myResourceResource struct {
    service services.MyResourceService
}
```

4. Create service in `pkg/services/{resource}.go`:
```go
type MyResourceService interface {
    Get() ([]MyResource, error)
}

type myResourceService struct {
    BaseService
}
```

5. Create flags in `internal/flags/{resource}.go`
6. Create descriptor in `internal/cli/descriptors/{resource}.yaml`
7. Register handler in `internal/handlers/handlers.go` NewHandler()

## Important Files

- `cmd/ipctl/main.go`: Entry point (13 lines)
- `internal/cli/root.go`: Cobra setup, command tree construction
- `internal/handlers/handlers.go`: Handler initialization (35 handlers)
- `internal/handlers/registry.go`: Interface-based handler discovery
- `internal/handlers/runner.go`: CommandRunner for executing operations
- `internal/runners/base.go`: BaseRunner with common functionality
- `pkg/client/client.go`: HTTP client (OAuth2, basic auth, cookies, TLS)
- `internal/config/config.go`: Configuration loading and management
- `internal/profile/manager.go`: Profile management
- `Makefile`: Build, test, and development commands

## Production Readiness Assessment

### Ready for Production ✓
- Context-aware HTTP client with timeouts
- TLS support with certificate verification option
- OAuth2 and basic auth
- Configuration via environment variables
- Extensive testing for services and resources (73 test files)
- Structured logging with configurable levels
- Git integration for import/export

### Needs Work Before Production ✗
- Handler layer completely untested (2,156 lines)
- context.TODO() in MongoDB operations (localaaa)
- No retry logic for transient failures
- Inconsistent error messages
- No request tracing or correlation IDs
- Minimal input validation
- Deprecated methods still in use

## Quick Start for New Developers

1. **Understand the flow**: CLI → Handler → Runner → Resource → Service → Client
2. **Start at the layer you care about**:
   - CLI changes? Look at `internal/cli/` and `internal/handlers/`
   - Business logic? Look at `pkg/resources/`
   - API issues? Look at `pkg/services/` and `pkg/client/`
3. **Find examples**: Every resource follows the same pattern. Look at `automations.go` files in each layer.
4. **Read package docs**: Each major package has doc.go with comprehensive documentation.
5. **Run tests**: `make test` runs all tests with coverage
6. **Check descriptors**: YAML files in `internal/cli/descriptors/` define command structure

## Common Pitfalls

1. **Don't call logging.Fatal() in library code** - Return errors instead
2. **Don't put business logic in services** - Services are for API operations only
3. **Don't use context.TODO()** - Pass context through function signatures
4. **Don't implement Reader/Writer interfaces without testing** - Handler tests missing
5. **Don't duplicate GetByName() logic** - Use shared helper or base implementation
6. **Don't parse error responses in runners** - Move to service/resource layer
7. **Don't use deprecated service methods** - Use resource methods instead

## Summary

`ipctl` is a well-architected Go CLI with clean layering and separation of concerns. The handler/runner/resource/service structure is textbook clean architecture. The HTTP client is production-grade, configuration system is flexible, and Git integration is comprehensive.

**Main strengths**: Architecture, layering, HTTP client, configuration, package documentation

**Main issues**: Untested handler layer (critical), context.TODO() in production code (critical), deprecated methods still in use, some business logic leakage into runners

**Quick wins**: Add handler tests (1-2 days), fix context.TODO() (1 hour), fix UnmarshalData() (1 hour), consolidate GetByName() (1 day)

Overall assessment: **Production-ready architecture with maintenance issues that should be addressed before scaling.**
