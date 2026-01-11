# Config Package Refactoring - Phase 6: Dependency Injection

## Overview

Phase 6 implements dependency injection using interfaces to decouple components from the concrete `Config` type. This improves testability, reduces coupling, and follows the Interface Segregation Principle.

## Goals Achieved

1. ✅ Defined minimal interfaces for configuration access
2. ✅ Updated packages to accept interfaces instead of concrete types
3. ✅ Improved testability with mock implementations
4. ✅ Reduced coupling between components
5. ✅ Maintained backward compatibility

## Changes Made

### 1. Interface Definition (`internal/config/interfaces.go`)

Created small, focused interfaces:

```go
type Provider interface {
    ProfileProvider
    RepositoryProvider
    ApplicationProvider
    FeaturesProvider
    GitProvider
}

type ProfileProvider interface {
    GetProfile(name string) (*profile.Profile, error)
    ActiveProfile() (*profile.Profile, error)
}

type ApplicationProvider interface {
    GetWorkingDir() string
    GetDefaultProfile() string
    GetDefaultRepository() string
}

type FeaturesProvider interface {
    IsDatasetsEnabled() bool
}

type GitProvider interface {
    GetGitName() string
    GetGitEmail() string
    GetGitUser() string
}

type RepositoryProvider interface {
    GetRepository(name string) (*repository.Repository, error)
}
```

### 2. Config Implementation (`internal/config/config.go`)

Added getter methods to implement the interfaces:

```go
// Ensure compile-time interface compliance
var _ Provider = (*Config)(nil)

func (c *Config) GetWorkingDir() string { return c.WorkingDir }
func (c *Config) GetDefaultProfile() string { return c.DefaultProfile }
func (c *Config) GetDefaultRepository() string { return c.DefaultRepository }
func (c *Config) IsDatasetsEnabled() bool { return c.Features.DatasetsEnabled }
func (c *Config) GetGitName() string { return c.Git.Name }
func (c *Config) GetGitEmail() string { return c.Git.Email }
func (c *Config) GetGitUser() string { return c.Git.User }
// GetProfile, ActiveProfile, GetRepository already existed
```

### 3. Handlers Package Updates

Updated `RuntimeContext` and `Runtime` to use `config.Provider`:

```go
type RuntimeContext interface {
    GetClient() client.Client
    GetConfig() config.Provider  // Changed from *config.Config
    GetDescriptors() Descriptors
    GetTerminalConfig() *terminal.Config
    IsVerbose() bool
}

type Runtime struct {
    config config.Provider  // Changed from *config.Config
    // ... other fields
}

func NewRuntime(c client.Client, cfg config.Provider, termCfg *terminal.Config) *Runtime
```

### 4. Runners Package Updates

Updated `BaseRunner` and all runner constructors:

```go
type BaseRunner struct {
    config config.Provider  // Changed from *config.Config
    client client.Client
}

func NewBaseRunner(client client.Client, cfg config.Provider) BaseRunner

// Updated all 30+ runner constructors
func NewProjectRunner(client client.Client, cfg config.Provider) *ProjectRunner
func NewWorkflowRunner(c client.Client, cfg config.Provider) *WorkflowRunner
// ... etc.
```

Updated `Request` struct:

```go
type Request struct {
    Args    []string
    Common  any
    Options any
    Runner  Runner
    Config  config.Provider  // Changed from *config.Config
}
```

Updated utility functions to use specific interfaces:

```go
func GetProfile(name string, cfg config.ProfileProvider) (*profile.Profile, error)
func NewClient(name string, cfg config.ProfileProvider) (client.Client, context.CancelFunc, error)
```

Updated exporter/importer to use interface methods:

```go
// Before: in.Config.Git.Name
// After: in.Config.GetGitName()
```

### 5. CLI Root Updates

Updated to use `config.Provider`:

```go
func runCli(c client.Client, cfg config.Provider, termCfg *terminal.Config) *cobra.Command

// Access feature flags via interface
if runtime.GetConfig().IsDatasetsEnabled() {
    // ...
}
```

### 6. Testing Infrastructure

Created mock provider (`internal/config/mock_test.go`):

```go
type mockProvider struct {
    workingDir        string
    defaultProfile    string
    defaultRepository string
    datasetsEnabled   bool
    gitName           string
    gitEmail          string
    gitUser           string
    profiles          map[string]*profile.Profile
    repositories      map[string]*repository.Repository
}

// Implements all Provider interfaces
var _ config.Provider = (*mockProvider)(nil)
```

Created comprehensive interface tests (`internal/config/interfaces_test.go`):

- Compile-time interface compliance tests
- Runtime behavior tests
- Interface segregation demonstration
- Mock provider functionality tests

## Benefits

### 1. Improved Testability

**Before:**
```go
func TestMyFunction(t *testing.T) {
    // Need to create full Config with file loading
    cfg := config.NewConfig(nil, nil, "", "", "")
    // Hard to control what config returns
}
```

**After:**
```go
func TestMyFunction(t *testing.T) {
    mock := newMockProvider()
    mock.profiles["test"] = &profile.Profile{Host: "testhost"}
    // Easy to control behavior
}
```

### 2. Reduced Coupling

Components now depend on interfaces, not concrete types:

```go
// Component only needs profiles
func Connect(profiles config.ProfileProvider) error {
    profile, _ := profiles.ActiveProfile()
    // Can't accidentally access Git, Features, etc.
}
```

### 3. Clearer Dependencies

Interface names document what's needed:

```go
func ConfigureGit(git config.GitProvider) error {
    // Obviously needs git configuration
}

func LoadProfiles(profiles config.ProfileProvider) error {
    // Obviously needs profile access
}
```

### 4. Interface Segregation

Small interfaces mean:
- Components depend only on what they use
- Easier to mock (implement fewer methods)
- Changes to unused interfaces don't affect component

### 5. Backward Compatibility

`*config.Config` implements `config.Provider`, so existing code works:

```go
// Old code still works
cfg := config.NewConfig(nil, nil, "", "", "")
runtime := handlers.NewRuntime(client, cfg, termCfg)  // cfg is Provider
```

## Files Modified

### New Files
- `internal/config/interfaces.go` - Interface definitions
- `internal/config/mock_test.go` - Mock provider for testing
- `internal/config/interfaces_test.go` - Interface compliance tests
- `docs/config-refactoring-phase6.md` - This document

### Modified Files
- `internal/config/config.go` - Added getter methods
- `internal/config/README.md` - Added interface documentation
- `internal/handlers/handlers.go` - Use Provider interface
- `internal/cli/root.go` - Use Provider interface
- `internal/runners/runners.go` - Use Provider interface
- `internal/runners/request.go` - Use Provider interface
- `internal/runners/utils.go` - Use specific interfaces
- `internal/runners/exporter.go` - Use interface methods
- `internal/runners/importer.go` - Use interface methods
- `internal/runners/*.go` - 30+ runner constructors updated

## Testing

All tests pass:
```bash
go test ./...
# All packages: PASS

go test ./internal/config/...
# 25 tests, all passing
# Including new interface compliance tests
```

Application builds successfully:
```bash
go build ./cmd/ipctl
# Build successful
```

## Usage Examples

### Using Full Provider Interface

```go
func Initialize(cfg config.Provider) error {
    // Access multiple aspects
    profile, _ := cfg.ActiveProfile()
    if cfg.IsDatasetsEnabled() {
        // Initialize datasets
    }
    return nil
}
```

### Using Specific Interfaces (Preferred)

```go
// Only needs profiles
func ConnectToServer(profiles config.ProfileProvider) error {
    profile, err := profiles.ActiveProfile()
    if err != nil {
        return err
    }
    return client.Connect(profile)
}

// Only needs git config
func CreateCommit(git config.GitProvider) error {
    name := git.GetGitName()
    email := git.GetGitEmail()
    return git.Commit(name, email)
}
```

### Testing with Mock

```go
func TestConnect(t *testing.T) {
    mock := newMockProvider()
    mock.profiles["default"] = &profile.Profile{
        Host: "localhost",
        Port: 8080,
    }

    err := ConnectToServer(mock)
    assert.NoError(t, err)
}
```

## Best Practices

1. **Prefer Specific Interfaces**: Use `ProfileProvider`, `GitProvider`, etc. instead of full `Provider`
2. **Use Interface Methods**: Call `cfg.GetGitName()` instead of `cfg.Git.Name`
3. **Mock for Testing**: Use `mockProvider` instead of real `Config` in tests
4. **Document Dependencies**: Interface parameter names make dependencies clear

## Next Steps

Future improvements to consider:

1. **More Granular Interfaces**: Split `ProfileProvider` into read/write operations
2. **Context-Aware Methods**: Add `context.Context` to interface methods
3. **Validation Interfaces**: Add `Validator` interface for config validation
4. **Observer Pattern**: Add interfaces for config change notifications

## Conclusion

Phase 6 successfully implements dependency injection using interfaces, improving testability and reducing coupling throughout the codebase. All existing functionality is preserved while enabling better testing and more flexible architecture.

The implementation follows Go best practices:
- Accept interfaces, return structs
- Small, focused interfaces
- Interface Segregation Principle
- Compile-time interface compliance checks
- Comprehensive test coverage
