// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package runners provides the command execution layer for CLI operations.
//
// Runners sit between handlers and resources, orchestrating CLI-specific concerns
// such as formatting output, handling flags, and coordinating Git operations.
// Business logic resides in the resources layer, not in runners.
//
// # Architecture
//
// The runner layer sits between handlers and resources:
//
//	Handler → Runner → Resource → Service → API
//
// Runners receive Request objects from handlers, delegate business logic to
// resources, and return Response objects formatted for display.
//
// # Interface-Based Design
//
// Runners depend on resource interfaces, not concrete types. This enables:
//   - Easy unit testing with mocked resources
//   - Dependency injection
//   - Clean separation from business logic
//
// Example:
//
//	type AccountRunner struct {
//	    BaseRunner
//	    resource resources.AccountResourcer  // Interface, not concrete type
//	}
//
//	func NewAccountRunner(client client.Client, cfg *config.Config) *AccountRunner {
//	    return &AccountRunner{
//	        BaseRunner: NewBaseRunner(client, cfg),
//	        resource:   resources.NewAccountResource(services.NewAccountService(client)),
//	    }
//	}
//
// # Runner Interfaces
//
// Runners implement various interfaces based on supported operations:
//
//	type Reader interface {
//	    Get(Request) (*Response, error)
//	    Describe(Request) (*Response, error)
//	}
//
//	type Writer interface {
//	    Create(Request) (*Response, error)
//	    Delete(Request) (*Response, error)
//	    Clear(Request) (*Response, error)
//	}
//
//	type Copier interface {
//	    Copy(Request) (*Response, error)
//	    CopyFrom(profile, name string) (any, error)
//	    CopyTo(profile string, data any, replace bool) (any, error)
//	}
//
//	type Editor interface {
//	    Edit(Request) (*Response, error)
//	}
//
//	type Importer interface {
//	    Import(Request) (*Response, error)
//	}
//
//	type Exporter interface {
//	    Export(Request) (*Response, error)
//	}
//
//	type Controller interface {
//	    Start(Request) (*Response, error)
//	    Stop(Request) (*Response, error)
//	    Restart(Request) (*Response, error)
//	}
//
//	type Inspector interface {
//	    Inspect(Request) (*Response, error)
//	}
//
//	type Dumper interface {
//	    Dump(Request) (*Response, error)
//	}
//
//	type Loader interface {
//	    Load(Request) (*Response, error)
//	}
//
// # Request Structure
//
// Handlers pass Request objects to runners:
//
//	type Request struct {
//	    Args    []string       // Command arguments
//	    Common  any            // Common flags (Flagger interface)
//	    Options any            // Command-specific flags
//	    Runner  Runner         // Runner instance
//	    Config  *config.Config // Global configuration
//	}
//
// The Args field contains positional arguments (resource names, IDs, etc.).
// Common contains flags shared across commands (--output, --verbose, etc.).
// Options contains command-specific flags parsed by the handler.
//
// # Response Structure
//
// Runners return Response objects to handlers:
//
//	type Response struct {
//	    Object   any      // Data for structured output (JSON/YAML)
//	    Text     string   // Human-readable text output
//	    Template string   // Template string for custom formatting
//	    Keys     []string // Table column keys for tabular output
//	}
//
// Handlers use the Response to format output based on user preferences.
//
// # Creating Runners
//
// Runners are created with a client and configuration, and initialize resources:
//
//	type ProjectRunner struct {
//	    BaseRunner
//	    resource resources.ProjectResourcer  // Resource interface
//	}
//
//	func NewProjectRunner(client client.Client, cfg *config.Config) *ProjectRunner {
//	    return &ProjectRunner{
//	        BaseRunner: NewBaseRunner(client, cfg),
//	        resource:   resources.NewProjectResource(services.NewProjectService(client)),
//	    }
//	}
//
// Runners delegate all business logic to resources. They focus on:
//   - Parsing CLI flags and arguments
//   - Calling resource methods with appropriate parameters
//   - Formatting responses for output
//   - Handling Git operations for import/export
//
// # Implementing Reader
//
// The Reader interface retrieves and describes resources:
//
//	func (r *ProjectRunner) Get(in Request) (*Response, error) {
//	    projects, err := r.resource.GetAll()  // Delegate to resource
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &Response{
//	        Keys:   []string{"name", "description"},
//	        Object: projects,
//	    }, nil
//	}
//
//	func (r *ProjectRunner) Describe(in Request) (*Response, error) {
//	    // Resource handles business logic (client-side filtering)
//	    project, err := r.resource.GetByName(in.Args[0])
//	    if err != nil {
//	        return nil, err
//	    }
//	    text := fmt.Sprintf("Name: %s\nDescription: %s",
//	        project.Name, project.Description)
//	    return &Response{
//	        Text:   text,
//	        Object: project,
//	    }, nil
//	}
//
// # Implementing Writer
//
// The Writer interface creates, deletes, and clears resources:
//
//	func (r *ProjectRunner) Create(in Request) (*Response, error) {
//	    name := in.Args[0]
//	    project, err := r.resource.Create(name)  // Delegate to resource
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &Response{
//	        Text: fmt.Sprintf("Created project: %s", project.Name),
//	        Object: project,
//	    }, nil
//	}
//
//	func (r *ProjectRunner) Delete(in Request) (*Response, error) {
//	    // Resource handles finding by name
//	    project, err := r.resource.GetByName(in.Args[0])
//	    if err != nil {
//	        return nil, err
//	    }
//	    // Resource handles deletion
//	    if err := r.resource.Delete(project.Id); err != nil {
//	        return nil, err
//	    }
//	    return &Response{
//	        Text: fmt.Sprintf("Deleted project: %s", project.Name),
//	    }, nil
//	}
//
// # Implementing Importer/Exporter
//
// Import and export operations support Git repositories and local filesystems:
//
//	func (r *ProjectRunner) Export(in Request) (*Response, error) {
//	    // Resource handles business logic
//	    project, err := r.resource.GetByName(in.Args[0])
//	    if err != nil {
//	        return nil, err
//	    }
//	    exported, err := r.resource.Export(project.Id)
//	    if err != nil {
//	        return nil, err
//	    }
//	    // Runner handles CLI concerns (file/Git operations)
//	    if err := exportAssetFromRequest(in, exported, filename); err != nil {
//	        return nil, err
//	    }
//	    return &Response{
//	        Text: fmt.Sprintf("Exported project: %s", project.Name),
//	    }, nil
//	}
//
// Export helpers support writing to:
//   - Local filesystem paths
//   - Git repositories (clone, modify, commit, push)
//
// # Implementing Copier
//
// Copy operations transfer resources between profiles/servers:
//
//	func (r *ProjectRunner) Copy(in Request) (*Response, error) {
//	    // Implement copy logic between profiles
//	    fromData, err := r.CopyFrom(fromProfile, name)
//	    if err != nil {
//	        return nil, err
//	    }
//	    toData, err := r.CopyTo(toProfile, fromData, replace)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &Response{
//	        Text: fmt.Sprintf("Copied from %s to %s", fromProfile, toProfile),
//	    }, nil
//	}
//
// # Git Integration
//
// Runners support Git operations for import/export:
//
//	repo := NewRepository(url,
//	    WithReference("main"),
//	    WithPrivateKeyFile("~/.ssh/id_rsa"),
//	)
//	path, err := repo.Clone(&FileReaderImpl{}, &ClonerImpl{})
//	if err != nil {
//	    return nil, err
//	}
//	defer os.RemoveAll(path)
//
//	// Modify files
//	// ...
//
//	err = repo.CommitAndPush(path, "Update from ipctl")
//	if err != nil {
//	    return nil, err
//	}
//
// # Error Handling
//
// Runners return errors for:
//   - Service operation failures
//   - Resource not found conditions
//   - Validation errors
//   - File I/O errors
//   - Git operation failures
//
// Always return descriptive errors with context:
//
//	return nil, fmt.Errorf("project %q not found", name)
//	return nil, fmt.Errorf("failed to export project: %w", err)
//
// # Common Patterns
//
// Get resource by name (business logic in resource layer):
//
//	resource, err := r.resource.GetByName(in.Args[0])
//	if err != nil {
//	    return nil, err
//	}
//
// Check for existing resource:
//
//	existing, err := r.resource.GetByName(name)
//	if existing != nil {
//	    return nil, fmt.Errorf("resource already exists")
//	}
//
// Iterate with pagination (handled by resource/service):
//
//	resources, err := r.resource.GetAll()
//	if err != nil {
//	    return nil, err
//	}
//	for _, res := range resources {
//	    // Format for output
//	}
//
// # Separation of Concerns
//
// Runners should handle:
//   - CLI flag parsing and validation
//   - Formatting responses for display
//   - Git operations for import/export
//   - File I/O for local operations
//
// Runners should NOT handle:
//   - Business logic (delegate to resources)
//   - Client-side filtering (delegate to resources)
//   - Data transformation (delegate to resources)
//   - API calls (resources delegate to services)
//
// # Type Safety
//
// Use type assertions carefully with common flags:
//
//	common, ok := in.Common.(*flags.AssetImportCommon)
//	if !ok {
//	    return nil, fmt.Errorf("invalid flags type")
//	}
//
// Always validate that the type assertion succeeded before dereferencing.
//
// # Thread Safety
//
// Runner instances are created per-handler and may be used concurrently.
// Runners should maintain no mutable state and use only the services
// provided at initialization.
package runners
