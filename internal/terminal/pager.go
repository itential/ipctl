// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/itential/ipctl/internal/logging"
)

func truncateOutput(tabbedMsg string) string {
	logging.Trace()

	const maxColLen int = 50

	lines := strings.Split(tabbedMsg, "\n")
	colInLine := make([][]string, len(lines))
	maxColMsg := ""
	for i, line := range lines {
		if line != "" {
			colInLine[i] = strings.Split(line, "\t")
			for j := range colInLine[i] {
				if len(colInLine[i][j]) > maxColLen {
					colInLine[i][j] = fmt.Sprintf("%s…", colInLine[i][j][:maxColLen-1])
				}
				maxColMsg += fmt.Sprintf("%s\t", colInLine[i][j])
			}
			maxColMsg += "\n"
		}
	}
	return maxColMsg
}

// DisplayTabWriterWithPager takes in a string for a table that already has tabs
// and newlines set and prints a properly spaced table with pagination.
// It respects the $PAGER environment variable and falls back to direct output
// if the pager is not available or encounters an error.
func DisplayTabWriterWithPager(ctx context.Context, tabbedMsg string, maxlen, padding int, limitColLen bool) error {
	logging.Trace()

	// Respect $PAGER environment variable
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	// Check if pager exists
	if _, err := exec.LookPath(pagerCmd); err != nil {
		// Fallback to direct output if pager not available
		logging.Warn("pager '%s' not found, using direct output", pagerCmd)
		DisplayTabWriter(tabbedMsg, maxlen, padding, limitColLen)
		return nil
	}

	cmd := exec.CommandContext(ctx, pagerCmd, "-R")
	r, stdin := io.Pipe()
	defer r.Close()

	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	tw := tabwriter.NewWriter(stdin, maxlen, 1, padding, ' ', 0)

	if limitColLen {
		tabbedMsg = truncateOutput(tabbedMsg)
	}

	// Start the pager
	if err := cmd.Start(); err != nil {
		logging.Warn("failed to start pager: %v, using direct output", err)
		DisplayTabWriter(tabbedMsg, maxlen, padding, limitColLen)
		return nil
	}

	// Write data
	if _, err := fmt.Fprintln(tw, tabbedMsg); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to write to pager: %w", err)
	}
	if err := tw.Flush(); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to flush table writer: %w", err)
	}
	stdin.Close()

	// Wait for pager to exit
	if err := cmd.Wait(); err != nil {
		// Don't return error if user just quit the pager
		if ctx.Err() == nil {
			logging.Debug("pager exited: %v", err)
		}
	}

	return nil
}
