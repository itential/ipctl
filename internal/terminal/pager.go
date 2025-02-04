// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/itential/ipctl/pkg/logger"
)

func truncateOutput(tabbedMsg string) string {
	logger.Trace()

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

// DisplayTabWriterWithPager taks in a string for a table that already has tabs
// and newlines set and printes a propertly spaced table with pagination
func DisplayTabWriterWithPager(tabbedMsg string, maxlen, padding int, limitColLen bool) {
	logger.Trace()

	cmd := exec.Command("less")
	r, stdin := io.Pipe()

	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//tw := tabwriter.NewWriter(os.Stdout, maxlen, 1, padding, ' ', 0)
	tw := tabwriter.NewWriter(stdin, maxlen, 1, padding, ' ', 0)

	if limitColLen {
		tabbedMsg = truncateOutput(tabbedMsg)
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		cmd.Run()
	}()

	fmt.Fprintln(tw, tabbedMsg)
	tw.Flush()
	stdin.Close()

	<-c
}
