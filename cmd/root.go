// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmd

import (
	"context"
	"embed"
	"os"
	"path"
	"strings"
	"time"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/handlers"
	"github.com/itential/ipctl/internal/metadata"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/spf13/cobra"
)

//go:embed descriptors/*.yaml
var content embed.FS

// Provides the general description of the application.  This should be moved
// into descriptors.
const description = `Manage Itential Platform

  Find more information at: https://docs.itential.com
`

// loadCommands will load the command tree for the application.  All top level
// comamnds are defined by this function except for the `version` command which
// is defined below.
func loadCommands(cmd *cobra.Command, runtime handlers.Runtime) {
	addRootCommand(cmd, runtime, "Asset Commands:", assetCommands)
	if runtime.Config.FeaturesDatasetsEnabled {
		addRootCommand(cmd, runtime, "Dataset Commands:", datasetCommands)
	}
	addRootCommand(cmd, runtime, "Platform Commands:", platformCommands)
	addRootCommand(cmd, runtime, "Plugin Commands:", pluginCommands)
}

// versionCommand is a top level command that displays the current application
// version.
func versionCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			terminal.Display("version: %s", metadata.Version)
			terminal.Display("commit: %s", metadata.Build)
			e, _ := os.Executable()
			terminal.Display("executable: %s", path.Dir(e))
			terminal.Display("")
		},
	}
	return cmd
}

// runCli builds and runs the CLI command.   It will create the command tree
// using cobra and execute the command.
func runCli(c client.Client, cfg *config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ipctl",
		Short: strings.Split(description, "\n")[0],
		Long:  description,
	}

	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	runtime := handlers.NewRuntime(c, cfg)

	cmd.PersistentFlags().BoolVar(&runtime.Verbose, "verbose", runtime.Verbose, "Enable verbose output")
	cmd.PersistentFlags().StringVar(&runtime.Config.TerminalDefaultOutput, "output", runtime.Config.TerminalDefaultOutput, "Output format")

	// Note: Values are read during /pkg/config's initialization
	cmd.PersistentFlags().String("config", "", "Path to the configuration file")
	cmd.PersistentFlags().String("profile", "", "Connection profile to use")

	loadCommands(cmd, runtime)

	cmd.AddCommand(versionCommand())

	return cmd
}

// Execute is the entrypoint to the CLI called from main.  This fucntion will
// load the configuration file, initialize the logger, create the client and
// run the application.  It will return an int that is to be used as the return
// code.
func Execute() int {
	cfg := config.NewConfig(nil, nil, "", "", "")
	logger.InitializeLogger(cfg)

	if metadata.Version != "" && metadata.Build != "" {
		logger.Info("ipctl %s (%s)", metadata.Version, metadata.Build)
	} else {
		sha, err := metadata.GetCurrentSha()
		if err == nil {
			logger.Info("ipctl running from commit %s", sha)
		} else {
			logger.Info("ipctl unable to determine source")
		}
	}

	profile, err := cfg.ActiveProfile()
	if err != nil {
		cmdutils.CheckError(err, cfg.TerminalNoColor)
	}

	logger.Info("connection timeout is %v second(s)", profile.Timeout)

	var ctx context.Context
	var cancel context.CancelFunc

	if profile.Timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			time.Duration(profile.Timeout)*time.Second,
		)
	}
	defer cancel()

	c := client.New(ctx, profile)

	if err := runCli(c, cfg).Execute(); err != nil {
		cmdutils.CheckError(err, cfg.TerminalNoColor)
	}

	return 0
}
