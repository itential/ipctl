// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package flags provides command-line flag definitions and parsing for the
// ipctl CLI application.
//
// This package defines flag structures for all commands, implements flag
// validation, and provides interfaces for accessing flag values in handlers
// and runners.
//
// # Flag Structure
//
// Flags are organized into two categories:
//
// Common Flags:
//   - Shared across multiple commands
//   - Examples: --output, --filter, --repository, --path
//
// Command-Specific Flags:
//   - Unique to specific operations
//   - Examples: --expand (export), --replace (import), --members (project)
//
// # Flagger Interface
//
// The Flagger interface provides type-safe flag access:
//
//	type Flagger interface {
//	    Flags() *pflag.FlagSet
//	}
//
// Additional interfaces extend Flagger for specific capabilities:
//
//	type Outputter interface {
//	    Flagger
//	    GetOutput() string
//	}
//
//	type Filterer interface {
//	    Flagger
//	    GetFilter() []string
//	}
//
//	type Gitter interface {
//	    Flagger
//	    GetRepository() string
//	    GetReference() string
//	    GetPrivateKeyFile() string
//	}
//
//	type Committer interface {
//	    Flagger
//	    GetPath() string
//	    GetMessage() string
//	}
//
//	type Paramer interface {
//	    GetParams() []string
//	    ParseParams() (map[string]string, error)
//	}
//
// # Common Flag Types
//
// AssetCommon provides standard flags for asset operations:
//
//	type AssetCommon struct {
//	    Output string
//	    Filter []string
//	}
//
//	func (f *AssetCommon) AddFlags(flags *pflag.FlagSet) {
//	    flags.StringVarP(&f.Output, "output", "o", "human",
//	        "Output format (human|json|yaml)")
//	    flags.StringSliceVar(&f.Filter, "filter", []string{},
//	        "Filter results by field=value")
//	}
//
// AssetImportCommon provides flags for import operations:
//
//	type AssetImportCommon struct {
//	    Replace         bool
//	    Repository      string
//	    Reference       string
//	    PrivateKeyFile  string
//	}
//
// AssetExportCommon provides flags for export operations:
//
//	type AssetExportCommon struct {
//	    Path            string
//	    Repository      string
//	    Reference       string
//	    PrivateKeyFile  string
//	    Message         string
//	}
//
// # Resource-Specific Flags
//
// Each resource type defines specific flag structures:
//
//	type ProjectImportOptions struct {
//	    Members []string
//	}
//
//	func (f *ProjectImportOptions) AddFlags(flags *pflag.FlagSet) {
//	    flags.StringSliceVar(&f.Members, "members", []string{},
//	        "Project members in format type=account,name=user,access=editor")
//	}
//
//	type ProjectExportOptions struct {
//	    Expand bool
//	}
//
//	func (f *ProjectExportOptions) AddFlags(flags *pflag.FlagSet) {
//	    flags.BoolVar(&f.Expand, "expand", false,
//	        "Expand project components into separate files")
//	}
//
// # Flag Registration
//
// Handlers register flags during command creation:
//
//	func (h *ProjectHandler) Import(r *Runtime) *cobra.Command {
//	    common := &flags.AssetImportCommon{}
//	    options := &flags.ProjectImportOptions{}
//
//	    cmd := &cobra.Command{
//	        Use:   "project [path]",
//	        Short: "Import project from file or repository",
//	        Run:   h.createImportRunner(common, options),
//	    }
//
//	    // Register flags
//	    common.AddFlags(cmd.Flags())
//	    options.AddFlags(cmd.Flags())
//
//	    return cmd
//	}
//
// # Flag Validation
//
// Flags can implement validation logic:
//
//	func (f *AssetImportCommon) Validate() error {
//	    if f.Repository != "" && f.Reference == "" {
//	        return fmt.Errorf("--reference required with --repository")
//	    }
//	    return nil
//	}
//
// # Accessing Flags in Runners
//
// Runners receive flags through the Request structure:
//
//	func (r *ProjectRunner) Import(in Request) (*Response, error) {
//	    common := in.Common.(*flags.AssetImportCommon)
//	    options := in.Options.(*flags.ProjectImportOptions)
//
//	    if common.Replace {
//	        // Handle replace logic
//	    }
//
//	    // Parse custom query parameters
//	    queryParams, err := common.ParseParams()
//	    if err != nil {
//	        return nil, fmt.Errorf("invalid params: %w", err)
//	    }
//
//	    // Pass query params to service layer
//	    req := &services.Request{
//	        query: queryParams,
//	    }
//
//	    for _, member := range options.Members {
//	        // Process member definitions
//	    }
//	}
//
// # Standard Flags
//
// Output Format (--output, -o):
//   - human: Human-readable table format (default)
//   - json: JSON output
//   - yaml: YAML output
//
// Filter (--filter):
//   - Filter results by field values
//   - Format: field=value
//   - Multiple filters: --filter name=test --filter status=active
//
// Repository (--repository):
//   - Git repository URL for import/export
//   - Supports: https://, git://, ssh://, file://
//   - Can use named repositories from config: file://@myrepo
//
// Reference (--reference):
//   - Git branch, tag, or commit to use
//   - Required when --repository is specified
//
// Private Key (--private-key-file):
//   - SSH private key for Git authentication
//   - Path to key file: ~/.ssh/id_rsa
//
// Path (--path):
//   - Local filesystem path for import/export
//   - Supports: absolute paths, relative paths, ~ expansion
//
// Message (--message):
//   - Commit message for Git operations
//   - Used with --repository during export
//
// Replace (--replace):
//   - Replace existing resource during import
//   - Default: false (fail if resource exists)
//
// Query Parameters (--params):
//   - Pass custom query parameters to API requests
//   - Format: key=value
//   - Can be specified multiple times
//   - URL encoded automatically
//   - Example: --params limit=10 --params offset=20 --params filter=active
//   - Values can contain spaces, special characters, and URLs
//   - Available on all asset operations (get, import, export, etc.) and API commands
//
// # Flag Persistence
//
// Persistent flags are available to all subcommands:
//
//	rootCmd.PersistentFlags().StringVar(&profile, "profile", "",
//	    "Connection profile to use")
//	rootCmd.PersistentFlags().StringVar(&config, "config", "",
//	    "Path to config file")
//	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false,
//	    "Enable verbose output")
//
// # Flag Aliases
//
// Common flags support short forms:
//
//	flags.StringVarP(&output, "output", "o", "human", "Output format")
//	// Can use: --output json OR -o json
//
// # Flag Groups
//
// Related flags can be grouped:
//
//	gitGroup := "Git Options"
//	cmd.Flags().SetAnnotation("repository", cobra.BashCompCustom, []string{gitGroup})
//	cmd.Flags().SetAnnotation("reference", cobra.BashCompCustom, []string{gitGroup})
//	cmd.Flags().SetAnnotation("private-key-file", cobra.BashCompCustom, []string{gitGroup})
//
// # Environment Variables
//
// Flags can be overridden with environment variables:
//
//	IPCTL_OUTPUT=json ipctl get projects
//	IPCTL_REPOSITORY=git@github.com:user/repo.git ipctl export project myproject
//
// # Best Practices
//
// Define flags in resource-specific files:
//   - projects.go: ProjectImportOptions, ProjectExportOptions
//   - workflows.go: WorkflowImportOptions, WorkflowExportOptions
//
// Use common flag types when possible:
//   - AssetCommon for standard operations
//   - AssetImportCommon for import operations
//   - AssetExportCommon for export operations
//
// Implement Validate() for complex validation:
//   - Cross-field validation
//   - Format validation
//   - Required field checks
//
// Document flags with clear descriptions:
//   - Explain what the flag does
//   - Show format examples
//   - List valid values
//
// # Type Safety
//
// Use type assertions with validation:
//
//	common, ok := in.Common.(*flags.AssetImportCommon)
//	if !ok {
//	    return nil, fmt.Errorf("invalid common flags type")
//	}
//
// This prevents panics from incorrect type casts.
package flags
