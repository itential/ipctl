// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/itential/ipctl/internal/logging"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

var (
	displayBuffer   bytes.Buffer
	displayToStdout bool = true
)

// Display Write a message to stdout.
func Display(format string, args ...interface{}) {
	displayBuffer.Reset()
	displayBuffer.WriteString(fmt.Sprintf(format, args...))
	if displayToStdout {
		fmt.Printf("%s\n", displayBuffer.String())
	}
}

// DisplayError prints a formatted error message to the terminal and prints the same message to the logger
// which will exit torero
func Error(err error, terminalNoColor bool) {
	if terminalNoColor {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "\x1b[31;1mError:\x1b[0m %s\n", err)
	}
}

func Warning(format string, args ...interface{}) {
	Display(fmt.Sprintf("WARNING: %s", format), args...)
}

// DisplayTabWriter takes in a string for a table that already has tab and newlines set
// and prints a properly spaced table
// E.g. DisplayTabWriter("COL1\tCOL2\nVal1\tVal2\n", 3, false)
func DisplayTabWriter(tabbedMsg string, maxlen, padding int, limitColLen bool) {
	tw := tabwriter.NewWriter(os.Stdout, maxlen, 1, padding, ' ', 0)

	if limitColLen {
		tabbedMsg = truncateOutput(tabbedMsg)
	}

	fmt.Fprintln(tw, tabbedMsg)
	tw.Flush()
}

func DisplayTabWriterString(msg []string, maxlen, padding int, limitColLen bool) {
	DisplayTabWriter(strings.Join(msg, "\n"), maxlen, padding, limitColLen)
}

func DisplayTabWriterStringWithPager(msg []string, maxlen, padding int, limitColLen bool) {
	DisplayTabWriterWithPager(strings.Join(msg, "\n"), maxlen, padding, limitColLen)
}

// Confirm Write a message to stdout and prompt for a confirmation of yes or no..
// Returns a bool value based on the response
func Confirm(prompt string, preamble string) bool {
	var ans string

	fmt.Println(preamble)

	for {
		fmt.Printf("%s(y/n)? ", prompt)
		fmt.Scanln(&ans)
		for _, item := range []string{"y", "n"} {
			if strings.ToLower(ans) == item {
				return ans == "y"
			}
		}
	}
}

// FormatTimestamp Returns a correctly formatted timestamp with a proper timezone
// which is typically based on the config variable TORERO_TERMINAL_TIMESTAMP_TIMEZONE.
func FormatTimestamp(timestamp time.Time, timezone *time.Location) string {
	return timestamp.In(timezone).Format(time.RFC3339)
}

// Password will prompt the user to enter a new password and mask the input so
// it does not echo back to the terminal.   It will return what text is entered
// by the user.
func Password() string {
	logging.Trace()

	prompt := promptui.Prompt{
		Label: "Password",
		Mask:  1,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{.}}: ",
			Valid:   "{{.}}: ",
			Invalid: "{{.}}: ",
		},
	}

	value, err := prompt.Run()
	if err != nil {
		logging.Fatal(err, "failed to get password")
	}

	return value
}

// DisplayJson accepts any interface object and will marshal the object to JSON
// format and display it to stdout.   This function will return an error if it
// cannot marhsal the object.
func DisplayJson(o any) error {
	logging.Trace()

	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return err
	}

	Display("%s\n", b)

	return nil
}

// DisplayYaml accepts any interface object and will marshal the object to YAML
// format and display it to stdout.   This function will return an error if it
// cannot marhsal the object.
func DisplayYaml(o any) error {
	logging.Trace()

	b, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	Display("%s\n", b)

	return nil
}
