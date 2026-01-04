// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package handlers provides command handlers that coordinate between the CLI
// layer and the business logic layer.
//
// Handlers are responsible for:
//   - Registering Cobra commands
//   - Parsing command-line flags
//   - Delegating to runners for execution
//   - Formatting and displaying output
//
// # Architecture
//
// The handler layer sits between the Cobra CLI framework and the runner layer:
//
//	CLI (Cobra) → Handlers → Runners → Services → Client → API
//
// Handlers convert user input into runner requests and format runner responses
// for terminal display.
//
// # Handler Types
//
// Handlers implement various interfaces based on the operations they support:
//
//	type Reader interface {
//	    Get(runtime *Runtime) *cobra.Command
//	    Describe(runtime *Runtime) *cobra.Command
//	}
//
//	type Writer interface {
//	    Create(runtime *Runtime) *cobra.Command
//	    Delete(runtime *Runtime) *cobra.Command
//	    Clear(runtime *Runtime) *cobra.Command
//	}
//
//	type Copier interface {
//	    Copy(runtime *Runtime) *cobra.Command
//	}
//
//	type Editor interface {
//	    Edit(runtime *Runtime) *cobra.Command
//	}
//
//	type Importer interface {
//	    Import(runtime *Runtime) *cobra.Command
//	}
//
//	type Exporter interface {
//	    Export(runtime *Runtime) *cobra.Command
//	}
//
//	type Controller interface {
//	    Start(runtime *Runtime) *cobra.Command
//	    Stop(runtime *Runtime) *cobra.Command
//	    Restart(runtime *Runtime) *cobra.Command
//	}
//
//	type Inspector interface {
//	    Inspect(runtime *Runtime) *cobra.Command
//	}
//
//	type Dumper interface {
//	    Dump(runtime *Runtime) *cobra.Command
//	}
//
//	type Loader interface {
//	    Load(runtime *Runtime) *cobra.Command
//	}
//
// # Handler Registry
//
// Handlers register themselves at initialization using interface-based discovery:
//
//	func NewHandler(r Runtime) Handler {
//	    register(
//	        NewProjectHandler(r, descriptors),
//	        NewWorkflowHandler(r, descriptors),
//	        NewAutomationHandler(r, descriptors),
//	        // ... 30+ handlers
//	    )
//	    return Handler{Runtime: &r, Descriptors: descriptors}
//	}
//
// The registry pre-computes handler lists by interface type for O(1) lookup:
//
//	readers := Readers()     // All handlers implementing Reader
//	writers := Writers()     // All handlers implementing Writer
//	exporters := Exporters() // All handlers implementing Exporter
//
// # Creating Handlers
//
// Handlers are created with a Runtime and Descriptors:
//
//	type ProjectHandler struct {
//	    *AssetHandler
//	}
//
//	func NewProjectHandler(r Runtime, d Descriptors) *ProjectHandler {
//	    runner := runners.NewProjectRunner(r.Client, r.Config)
//	    return &ProjectHandler{
//	        AssetHandler: NewAssetHandler(r, d, runner, "project", "projects"),
//	    }
//	}
//
// # AssetHandler
//
// Most handlers use AssetHandler which provides default implementations:
//
//	type AssetHandler struct {
//	    runtime     *Runtime
//	    descriptors Descriptors
//	    runner      runners.Runner
//	    singular    string
//	    plural      string
//	    flags       *AssetHandlerFlags
//	    // Interface capabilities discovered via reflection
//	    isReader    bool
//	    isWriter    bool
//	    isCopier    bool
//	    // ... etc
//	}
//
// AssetHandler automatically generates Cobra commands based on which interfaces
// the runner implements.
//
// # Descriptors
//
// Command metadata is defined in YAML descriptor files:
//
//	project:
//	  get:
//	    use: projects
//	    short: Get projects from the server
//	    long: |
//	      Retrieves all projects configured on the Itential Platform server.
//	    examples:
//	      - Get all projects: ipctl get projects
//	      - Filter by name: ipctl get projects --filter name=MyProject
//
// Descriptors are loaded at startup and used to populate command help text.
//
// # Command Generation
//
// Commands are generated dynamically based on handler capabilities:
//
//	func (h *AssetHandler) Get(r *Runtime) *cobra.Command {
//	    if !h.isReader {
//	        return nil // Handler doesn't support Get
//	    }
//	    cmd := &cobra.Command{
//	        Use:   h.descriptors.Get("get", "use"),
//	        Short: h.descriptors.Get("get", "short"),
//	        Long:  h.descriptors.Get("get", "long"),
//	        Run:   h.createGetRunner(),
//	    }
//	    // Add flags
//	    return cmd
//	}
//
// # Execution Flow
//
// When a command executes:
//
//  1. Cobra calls the command's Run function
//  2. Handler extracts flags and arguments
//  3. Handler creates a runners.Request
//  4. Handler calls runner method (e.g., runner.Get())
//  5. Runner returns runners.Response
//  6. Handler formats output using terminal package
//  7. Handler displays result to user
//
// Example:
//
//	Run: func(cmd *cobra.Command, args []string) {
//	    req := runners.Request{
//	        Args:    args,
//	        Common:  h.flags.Common,
//	        Options: h.flags.GetOptions,
//	        Runner:  h.runner,
//	        Config:  h.runtime.Config,
//	    }
//	    resp, err := h.runner.(runners.Reader).Get(req)
//	    cmdutils.CheckError(err, h.runtime.Config.TerminalNoColor)
//	    terminal.Display(resp.Text)
//	}
//
// # Special Handlers
//
// Some handlers have custom implementations:
//
// ServerHandler:
//   - Manages server-level operations (health, info, status)
//
// LocalaaaHandler:
//   - Provides local AAA operations
//
// ApiHandler:
//   - Generic HTTP API access
//
// # Output Formatting
//
// Handlers use the terminal package for output:
//
//	terminal.Display(text)              // Human-readable
//	terminal.DisplayJSON(data)          // JSON format
//	terminal.DisplayYAML(data)          // YAML format
//	terminal.DisplayTable(rows, cols)   // Tabular format
//	terminal.DisplayTemplate(tmpl, data) // Custom template
//
// Output format is controlled by flags (--output json/yaml/human) and
// configuration (terminal.default_output).
//
// # Error Handling
//
// Handlers use cmdutils.CheckError for consistent error reporting:
//
//	cmdutils.CheckError(err, cfg.TerminalNoColor)
//
// This displays the error and exits with non-zero status code.
//
// # Thread Safety
//
// Handlers are created once at startup and used concurrently for command
// execution. Handlers maintain no mutable state and are safe for concurrent use.
package handlers
