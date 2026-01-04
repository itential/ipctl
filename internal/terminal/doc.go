// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package terminal provides output formatting and display functions for the
// ipctl CLI application.
//
// This package handles rendering data in various formats (human-readable,
// JSON, YAML, tables) with support for color output, pagination, and
// templates.
//
// # Output Functions
//
// Display renders output based on format and configuration:
//
//	terminal.Display("Resource created successfully")
//	terminal.DisplayJSON(resource)
//	terminal.DisplayYAML(resource)
//	terminal.DisplayTable(rows, []string{"name", "status"})
//
// # Output Formats
//
// Human-readable (default):
//   - Plain text with optional color
//   - Tables for structured data
//   - Custom templates for complex layouts
//
// JSON:
//   - Indented JSON output
//   - Suitable for programmatic processing
//   - Preserves all data fields
//
// YAML:
//   - YAML formatted output
//   - Human-friendly structured format
//   - Preserves data hierarchy
//
// # Display Function
//
// The primary display function formats strings with color support:
//
//	terminal.Display("Success: Resource created")
//	terminal.Display("Error: %s", err)
//	terminal.Display("Found %d items", count)
//
// Color output is automatically disabled when:
//   - NO_COLOR environment variable is set
//   - terminal.no_color config option is true
//   - Output is redirected (not a TTY)
//
// # Error Display
//
// Display errors with color highlighting:
//
//	terminal.Error(err, noColor)
//	// Output: Error: resource not found (in red)
//
// Warning messages:
//
//	terminal.Warn("Resource already exists", noColor)
//	// Output: Warning: resource already exists (in yellow)
//
// Success messages:
//
//	terminal.Success("Operation completed", noColor)
//	// Output: Success: operation completed (in green)
//
// # JSON Output
//
// Display structured data as JSON:
//
//	data := map[string]interface{}{
//	    "name": "my-project",
//	    "status": "active",
//	}
//	terminal.DisplayJSON(data)
//
// Output:
//
//	{
//	  "name": "my-project",
//	  "status": "active"
//	}
//
// # YAML Output
//
// Display structured data as YAML:
//
//	terminal.DisplayYAML(data)
//
// Output:
//
//	name: my-project
//	status: active
//
// # Table Display
//
// Render data as formatted tables:
//
//	rows := []map[string]interface{}{
//	    {"name": "project1", "status": "active"},
//	    {"name": "project2", "status": "inactive"},
//	}
//	keys := []string{"name", "status"}
//	terminal.DisplayTable(rows, keys)
//
// Output:
//
//	NAME       STATUS
//	project1   active
//	project2   inactive
//
// Tables automatically:
//   - Calculate column widths
//   - Align headers and data
//   - Truncate long values
//   - Add borders and separators
//
// # Template Display
//
// Use Go templates for custom formatting:
//
//	tmpl := `
//	Name: {{.Name}}
//	Status: {{.Status}}
//	Created: {{.Created}}
//	`
//	terminal.DisplayTemplate(tmpl, data)
//
// Templates support:
//   - All Go template syntax
//   - Custom functions (sprig library)
//   - Conditional rendering
//   - Loops and ranges
//
// # Pagination
//
// Long output can be piped to a pager:
//
//	terminal.EnablePager()
//	terminal.Display(longOutput)
//
// Pagination is enabled when:
//   - terminal.pager config option is true
//   - Output exceeds terminal height
//   - Running in interactive terminal (TTY)
//
// Pager commands used (in order of preference):
//   - $PAGER environment variable
//   - less -R
//   - more
//
// # Color Support
//
// Control color output:
//
//	terminal.SetNoColor(true)  // Disable colors
//	terminal.SetNoColor(false) // Enable colors
//
// Color codes:
//   - Red: Errors and failures
//   - Yellow: Warnings and deprecations
//   - Green: Success messages
//   - Blue: Informational messages
//   - Gray: Metadata and timestamps
//
// # Output Redirection
//
// Output can be captured for testing:
//
//	terminal.CaptureOutput(true)
//	terminal.Display("test output")
//	output := terminal.GetCapturedOutput()
//	terminal.CaptureOutput(false)
//
// This is used in unit tests to verify command output without displaying
// to stdout.
//
// # Format Detection
//
// The output format is determined by:
//
//  1. --output flag on command line
//  2. IPCTL_TERMINAL_DEFAULT_OUTPUT environment variable
//  3. terminal.default_output config option
//  4. Default: human
//
// Format values:
//   - human: Human-readable with tables and color
//   - json: JSON formatted
//   - yaml: YAML formatted
//
// # Progress Indicators
//
// Display progress for long operations:
//
//	terminal.StartProgress("Exporting project...")
//	// Perform operation
//	terminal.StopProgress()
//	terminal.Display("Export complete")
//
// Progress indicators:
//   - Show spinner animation
//   - Update status text
//   - Automatically stop on completion
//
// # Status Messages
//
// Display operation status:
//
//	terminal.Status("Processing...", noColor)
//	// Perform operation
//	terminal.StatusComplete("Done", noColor)
//
// Status updates appear on the same line, providing real-time feedback
// without cluttering output.
//
// # Width Detection
//
// Terminal width is automatically detected:
//
//	width := terminal.GetWidth()
//
// Used for:
//   - Wrapping long lines
//   - Calculating table column widths
//   - Centering output
//
// Default width is 80 columns when detection fails.
//
// # Best Practices
//
// Use Display for all output:
//   - Centralized formatting
//   - Consistent color handling
//   - Respects user preferences
//
// Provide structured data to DisplayJSON/YAML:
//   - Enables programmatic processing
//   - Preserves all data fields
//   - Supports automation
//
// Use tables for list operations:
//   - Easier to scan visually
//   - Aligned columns
//   - Consistent formatting
//
// Respect NoColor configuration:
//   - Check config before adding color
//   - Fallback to plain text
//   - Support automation tools
//
// # Thread Safety
//
// Terminal functions use synchronized output to prevent interleaved writes
// from concurrent goroutines. All Display* functions are safe for concurrent use.
package terminal
