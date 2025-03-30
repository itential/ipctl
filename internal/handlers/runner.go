// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/spf13/cobra"
)

type RunCommand func(*cobra.Command, []string)

type CommandOptions func(*cobra.Command)

type CommandRunnerOption func(*CommandRunner)

type CommandRunner struct {
	Key         string
	Descriptors DescriptorMap
	Run         runners.RunnerFunc
	Common      flags.Flagger
	Options     flags.Flagger
	Runtime     *Runtime
	Runner      runners.Runner
	Confirm     bool
	Flags       *AssetHandlerFlags
	PreRun      func(args []string) error
	PostRun     func(args []string)
}

func NewCommandRunner(
	key string,
	desc DescriptorMap,
	run runners.RunnerFunc,
	runtime *Runtime,
	options flags.Flagger,
	opts ...CommandRunnerOption,
) *CommandRunner {

	cr := &CommandRunner{
		Key:         key,
		Descriptors: desc,
		Run:         run,
		Runtime:     runtime,
		Common:      options,
	}

	for _, opt := range opts {
		opt(cr)
	}

	return cr
}

func withOptions(f *AssetHandlerFlags) CommandRunnerOption {
	return func(c *CommandRunner) {
		switch c.Key {
		case "create":
			c.Options = f.Create
		case "delete":
			c.Options = f.Delete
		case "get":
			c.Options = f.Get
		case "describe":
			c.Options = f.Describe
		case "copy":
			c.Options = f.Copy
		case "clear":
			c.Options = f.Clear
		case "import":
			c.Options = f.Import
		case "export":
			c.Options = f.Export
		}
	}
}

func checkError(err error, runtime *Runtime) {
	cmdutils.CheckError(err, runtime.Config.TerminalNoColor)
}

func unableToDisplayOutput(nocolor bool) {
	terminal.Error(fmt.Errorf("unable to display response"), nocolor)
}

func NewCommand(c *CommandRunner) *cobra.Command {
	desc, exists := c.Descriptors[c.Key]

	if !exists || desc.Disabled {
		return nil
	}

	var example string

	if desc.Example != "" {
		var lines []string
		for _, ele := range strings.Split(desc.Example, "\n") {
			lines = append(lines, fmt.Sprintf("  %s", ele))
		}
		example = strings.Join(lines, "\n")
	}

	cmd := &cobra.Command{
		Use:     desc.Use,
		GroupID: desc.Group,

		Short: desc.Short(),
		Long:  desc.Description,

		Example: example,

		Hidden: desc.Hidden,

		Run: func(cmd *cobra.Command, args []string) {

			req := runners.Request{
				Args:    args,
				Options: c.Options,
				Common:  c.Common,
				Runner:  c.Runner,
				Config:  c.Runtime.Config,
			}

			resp, err := c.Run(req)
			checkError(err, c.Runtime)

			switch c.Runtime.Config.TerminalDefaultOutput {
			case "json":
				if resp.Object != nil {
					checkError(terminal.DisplayJson(resp.Object), c.Runtime)
				} else {
					unableToDisplayOutput(c.Runtime.Config.TerminalNoColor)
				}
			case "yaml":
				if resp.Object != nil {
					checkError(terminal.DisplayYaml(resp.Object), c.Runtime)
				} else {
					unableToDisplayOutput(c.Runtime.Config.TerminalNoColor)
				}
			case "human":
				if len(resp.Keys) > 0 {
					output := strings.Split(resp.String(), "\n")
					if c.Runtime.Config.TerminalPager {
						terminal.DisplayTabWriterStringWithPager(output, 3, 3, true)
					} else {
						terminal.DisplayTabWriterString(output, 3, 3, true)
					}
				} else {
					output := resp.String()
					if output == "" {
						unableToDisplayOutput(c.Runtime.Config.TerminalNoColor)
					}
					terminal.Display(fmt.Sprintf("%s\n", output))
				}
			}
		},
	}

	if desc.ExactArgs > 0 {
		cmd.Args = cobra.ExactArgs(desc.ExactArgs)
	}

	return cmd
}
