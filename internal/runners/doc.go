// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package runners provides command execution layer that implements business
// logic for CLI operations.
//
// Runners sit between handlers and services, coordinating multiple service
// calls, handling data transformations, and implementing complex workflows.
//
// # Architecture
//
// The runner layer sits between handlers and services:
//
//	Handlers → Runners → Services → Client → API
//
// Runners receive Request objects from handlers, execute business logic using
// services, and return Response objects for display.
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
// Runners are created with a client and configuration:
//
//	type ProjectRunner struct {
//	    config       *config.Config
//	    service      *services.ProjectService
//	    accounts     *services.AccountService
//	    groups       *services.GroupService
//	    userSettings *services.UserSettingsService
//	}
//
//	func NewProjectRunner(client client.Client, cfg *config.Config) *ProjectRunner {
//	    return &ProjectRunner{
//	        config:       cfg,
//	        service:      services.NewProjectService(client),
//	        accounts:     services.NewAccountService(client),
//	        groups:       services.NewGroupService(client),
//	        userSettings: services.NewUserSettingsService(client),
//	    }
//	}
//
// Runners often need multiple services to implement complex operations.
//
// # Implementing Reader
//
// The Reader interface retrieves and describes resources:
//
//	func (r *ProjectRunner) Get(in Request) (*Response, error) {
//	    projects, err := r.service.GetAll()
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
//	    project, err := r.service.GetByName(in.Args[0])
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
//	    project, err := r.service.Create(name)
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
//	    project, err := r.service.GetByName(in.Args[0])
//	    if err != nil {
//	        return nil, err
//	    }
//	    if err := r.service.Delete(project.Id); err != nil {
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
//	    project, err := r.service.GetByName(in.Args[0])
//	    if err != nil {
//	        return nil, err
//	    }
//	    exported, err := r.service.Export(project.Id)
//	    if err != nil {
//	        return nil, err
//	    }
//	    // Write to file or Git repository
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
// Get resource by name:
//
//	resource, err := r.service.GetByName(in.Args[0])
//	if err != nil {
//	    return nil, err
//	}
//
// Check for existing resource:
//
//	existing, err := r.service.GetByName(name)
//	if existing != nil {
//	    return nil, fmt.Errorf("resource already exists")
//	}
//
// Iterate with pagination:
//
//	resources, err := r.service.GetAll()
//	if err != nil {
//	    return nil, err
//	}
//	for _, res := range resources {
//	    // Process each resource
//	}
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
