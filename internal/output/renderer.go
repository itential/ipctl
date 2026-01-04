// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package output

import (
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/config"
)

// Renderer coordinates formatting and displaying output.
//
// The Renderer uses a Formatter to convert data into the appropriate format,
// then delegates to the terminal package for actual display. This separation
// allows command handlers to remain agnostic of output format details.
type Renderer struct {
	formatter Formatter
	config    *config.Config
}

// NewRenderer creates a new Renderer with the specified output format and configuration.
//
// The format parameter should be one of "json", "yaml", or "human".
// The config provides additional display settings like pager usage.
//
// Returns an error if the format is not supported.
func NewRenderer(format string, cfg *config.Config) (*Renderer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	usePager := cfg.TerminalPager
	formatter, err := NewFormatter(format, usePager)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		formatter: formatter,
		config:    cfg,
	}, nil
}

// Render formats and displays the response data.
//
// The rendering strategy depends on the configured output format:
//
//   - JSON: Renders Response.Object as formatted JSON
//   - YAML: Renders Response.Object as formatted YAML
//   - Human: Renders Response.String() output, optionally with pager for tabular data
//
// Returns an error if the response cannot be rendered or if the response
// data is not suitable for the configured format.
func (r *Renderer) Render(resp *runners.Response) error {
	if resp == nil {
		return fmt.Errorf("cannot render nil response")
	}

	switch f := r.formatter.(type) {
	case *JSONFormatter:
		return r.renderJSON(resp)
	case *YAMLFormatter:
		return r.renderYAML(resp)
	case *HumanFormatter:
		return r.renderHuman(resp, f)
	default:
		return fmt.Errorf("unknown formatter type: %T", f)
	}
}

// renderJSON formats and displays the response as JSON.
func (r *Renderer) renderJSON(resp *runners.Response) error {
	if resp.Object == nil {
		return fmt.Errorf("unable to display response: no object data available for JSON output")
	}

	formatted, err := r.formatter.Format(resp.Object)
	if err != nil {
		return err
	}

	terminal.Display("%s", formatted)
	return nil
}

// renderYAML formats and displays the response as YAML.
func (r *Renderer) renderYAML(resp *runners.Response) error {
	if resp.Object == nil {
		return fmt.Errorf("unable to display response: no object data available for YAML output")
	}

	formatted, err := r.formatter.Format(resp.Object)
	if err != nil {
		return err
	}

	terminal.Display("%s", formatted)
	return nil
}

// renderHuman formats and displays the response in human-readable format.
//
// For tabular data (when resp.Keys is set), the output is displayed using
// a tab writer, optionally with a pager. For other output types (templates
// or text), the formatted string is displayed directly.
func (r *Renderer) renderHuman(resp *runners.Response, formatter *HumanFormatter) error {
	// Convert the response to its string representation
	output := resp.String()

	// Check if we got an empty or error string
	if output == "" || output == "error formating response object" {
		return fmt.Errorf("unable to display response: no displayable data available")
	}

	// Format the output
	formatted, err := formatter.Format(output)
	if err != nil {
		return err
	}

	// Display based on whether we have tabular data
	if len(resp.Keys) > 0 {
		// Tabular output - use tab writer
		lines := strings.Split(formatted, "\n")
		if formatter.UsePager() {
			terminal.DisplayTabWriterStringWithPager(lines, 3, 3, true)
		} else {
			terminal.DisplayTabWriterString(lines, 3, 3, true)
		}
	} else {
		// Non-tabular output - display directly
		terminal.Display("%s", formatted)
	}

	return nil
}
