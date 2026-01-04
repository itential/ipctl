// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package cmdutils provides utility functions and helpers for command
// implementation and execution.
//
// This package contains common utilities used across handlers and runners
// including error handling, validation, descriptors loading, and command
// metadata management.
//
// # Error Handling
//
// CheckError provides centralized error handling for commands:
//
//	err := operation()
//	cmdutils.CheckError(err, cfg.TerminalNoColor)
//
// When an error occurs:
//   - Error message is displayed to stderr
//   - Color formatting applied (unless disabled)
//   - Process exits with non-zero status code
//   - Stack trace logged at debug level
//
// This ensures consistent error reporting across all commands.
//
// # Descriptors
//
// Command descriptors define metadata in YAML format:
//
//	project:
//	  get:
//	    use: projects
//	    short: Get projects from the server
//	    long: |
//	      Retrieves all projects configured on the Itential Platform.
//	    examples:
//	      - Get all projects: ipctl get projects
//	      - Filter by name: ipctl get projects --filter name=test
//	    disabled: false
//
// Loading descriptors:
//
//	descriptors, err := cmdutils.LoadDescriptors("cmd/descriptors")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Accessing descriptor values:
//
//	use := descriptors.Get("project", "get", "use")
//	short := descriptors.Get("project", "get", "short")
//	long := descriptors.Get("project", "get", "long")
//	examples := descriptors.GetExamples("project", "get")
//
// # Descriptor Structure
//
// Descriptors are organized hierarchically:
//
//	Resource → Command → Field
//	project  → get     → use, short, long, examples
//
// Common fields:
//   - use: Command usage string for Cobra
//   - short: Short description (one line)
//   - long: Long description (multiple lines)
//   - examples: List of usage examples
//   - disabled: Whether command is disabled
//
// # Validation
//
// Validate command arguments and flags:
//
//	err := cmdutils.ValidateArgs(args, 1, "project name required")
//	cmdutils.CheckError(err, noColor)
//
//	err = cmdutils.ValidateFlags(flags)
//	cmdutils.CheckError(err, noColor)
//
// Common validations:
//   - Argument count
//   - Required flags
//   - Flag value formats
//   - Mutual exclusivity
//
// # Command Builders
//
// Helper functions for building Cobra commands:
//
//	cmd := cmdutils.NewCommand(
//	    "projects",
//	    "Get projects from the server",
//	    runFunc,
//	)
//	cmdutils.AddOutputFlag(cmd)
//	cmdutils.AddFilterFlag(cmd)
//
// Standard command setup:
//   - Usage and description
//   - Run function
//   - Flags registration
//   - Example text
//
// # Flag Helpers
//
// Add common flags to commands:
//
//	cmdutils.AddOutputFlag(cmd)      // --output, -o
//	cmdutils.AddFilterFlag(cmd)      // --filter
//	cmdutils.AddFormatFlag(cmd)      // --format
//	cmdutils.AddVerboseFlag(cmd)     // --verbose, -v
//
// These ensure consistent flag naming and behavior across commands.
//
// # Exit Codes
//
// Standard exit codes used by the application:
//
//	const (
//	    ExitSuccess         = 0  // Command succeeded
//	    ExitError           = 1  // General error
//	    ExitUsageError      = 2  // Invalid command usage
//	    ExitNotFound        = 3  // Resource not found
//	    ExitPermissionError = 4  // Permission denied
//	    ExitTimeout         = 5  // Operation timeout
//	)
//
// Use appropriate exit codes:
//
//	cmdutils.ExitWithCode(cmdutils.ExitNotFound)
//
// # Resource Names
//
// Utilities for working with resource names:
//
//	singular := cmdutils.Singularize("projects")  // "project"
//	plural := cmdutils.Pluralize("project")       // "projects"
//	title := cmdutils.Title("project")            // "Project"
//
// Used for generating consistent command names and descriptions.
//
// # Path Resolution
//
// Resolve and expand file paths:
//
//	path, err := cmdutils.ExpandPath("~/.platform.d/config")
//	if err != nil {
//	    return err
//	}
//
//	absolute, err := cmdutils.AbsPath("./config")
//	if err != nil {
//	    return err
//	}
//
// Handles:
//   - Tilde (~) expansion to home directory
//   - Relative to absolute path conversion
//   - Environment variable expansion
//   - Symlink resolution
//
// # String Formatting
//
// Format strings for display:
//
//	wrapped := cmdutils.Wrap(text, 80)        // Wrap at 80 columns
//	truncated := cmdutils.Truncate(text, 50)  // Truncate to 50 chars
//	padded := cmdutils.Pad(text, 20)          // Pad to 20 chars
//
// Used for table formatting and terminal output.
//
// # Confirmation Prompts
//
// Prompt user for confirmation:
//
//	confirmed, err := cmdutils.Confirm("Delete project?", false)
//	if err != nil {
//	    return err
//	}
//	if !confirmed {
//	    return fmt.Errorf("operation cancelled")
//	}
//
// Prompts accept y/yes/Y/YES as confirmation, n/no/N/NO as rejection.
// Default value is used when user presses Enter without input.
//
// # Input Validation
//
// Validate user input:
//
//	err := cmdutils.ValidateEmail(email)
//	err = cmdutils.ValidateURL(url)
//	err = cmdutils.ValidateHostname(host)
//	err = cmdutils.ValidatePort(port)
//
// Common validation patterns:
//   - Email addresses
//   - URLs and URIs
//   - Hostnames and IP addresses
//   - Port numbers
//   - Identifiers and names
//
// # Debug Helpers
//
// Debug command execution:
//
//	cmdutils.DebugPrint("Processing resource: %s", name)
//	cmdutils.DebugDump(data)  // Pretty-print structure
//
// Debug output is controlled by:
//   - --verbose flag
//   - IPCTL_LOG_LEVEL=DEBUG
//   - log.level=DEBUG in config
//
// # Best Practices
//
// Use CheckError for all command errors:
//   - Consistent error formatting
//   - Proper exit codes
//   - Stack traces in debug mode
//
// Load descriptors at startup:
//   - Cache descriptor data
//   - Fail fast on invalid YAML
//   - Validate all required fields
//
// Validate early:
//   - Check arguments before execution
//   - Validate flags before API calls
//   - Return descriptive error messages
//
// Use helpers for common operations:
//   - Path expansion
//   - String formatting
//   - User prompts
//
// # Example Usage
//
//	// Load descriptors
//	descriptors, err := cmdutils.LoadDescriptors("cmd/descriptors")
//	cmdutils.CheckError(err, false)
//
//	// Create command
//	cmd := &cobra.Command{
//	    Use:   descriptors.Get("project", "get", "use"),
//	    Short: descriptors.Get("project", "get", "short"),
//	    Long:  descriptors.Get("project", "get", "long"),
//	    Run: func(cmd *cobra.Command, args []string) {
//	        // Validate
//	        err := cmdutils.ValidateArgs(args, 0, "")
//	        cmdutils.CheckError(err, cfg.TerminalNoColor)
//
//	        // Execute
//	        result, err := execute()
//	        cmdutils.CheckError(err, cfg.TerminalNoColor)
//
//	        // Display
//	        terminal.DisplayJSON(result)
//	    },
//	}
//
//	// Add flags
//	cmdutils.AddOutputFlag(cmd)
//
//	return cmd
package cmdutils
