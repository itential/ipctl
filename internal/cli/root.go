// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cli

import (
	"context"
	"embed"
	"os"
	"path"
	"strings"
	"time"

	"github.com/itential/ipctl/internal/app"
	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/handlers"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/client"
	"github.com/spf13/cobra"
)

// descriptorFiles embeds all YAML descriptor files at compile time.
// These descriptors define command structure, help text, and examples.
//
//go:embed descriptors/*.yaml
var descriptorFiles embed.FS

// Provides the general description of the application.  This should be moved
// into descriptors.
const description = `Manage Itential Platform

  Find more information at: https://docs.itential.com
`

// loadCommands will load the command tree for the application. All top level
// commands are defined by this function except for the `version` command which
// is defined below.
func loadCommands(cmd *cobra.Command, runtime *handlers.Runtime) {
	addRootCommand(cmd, runtime, "Asset Commands:", assetCommands)
	if runtime.GetConfig().IsDatasetsEnabled() {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			info := app.GetInfo()
			terminal.Display("version: %s", info.Version)
			terminal.Display("commit: %s", info.Build)
			e, err := os.Executable()
			if err != nil {
				terminal.Display("executable: <unknown> (error: %v)", err)
			} else {
				terminal.Display("executable: %s", path.Dir(e))
			}
			terminal.Display("")
			return nil
		},
	}
	return cmd
}

// runCli builds and runs the CLI command. It will create the command tree
// using cobra and execute the command.
func runCli(c client.Client, cfg config.Provider, termCfg *terminal.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ipctl",
		Short: strings.Split(description, "\n")[0],
		Long:  description,
	}

	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	runtime := handlers.NewRuntime(c, cfg, termCfg)

	cmd.PersistentFlags().BoolVar(&runtime.Verbose, "verbose", runtime.Verbose, "Enable verbose output")
	cmd.PersistentFlags().StringVar(&termCfg.DefaultOutput, "output", termCfg.DefaultOutput, "Output format")

	// Note: Values are read during /pkg/config's initialization
	cmd.PersistentFlags().String("config", "", "Path to the configuration file")
	cmd.PersistentFlags().String("profile", "", "Connection profile to use")

	loadCommands(cmd, runtime)

	cmd.AddCommand(versionCommand())

	return cmd
}

// Execute is the entrypoint to the CLI called from main. This function will
// load the configuration file, initialize the logger, create the client and
// run the application. It will return an int that is to be used as the return
// code.
func Execute() int {
	cfg := config.NewConfig(nil, nil, "", "", "")

	// Initialize logging with domain-specific config from environment
	logCfg := logging.LoadFromEnv()
	termCfg := terminal.LoadFromEnv()
	logging.InitializeLogger(logCfg, termCfg.NoColor)

	info := app.GetInfo()
	if info.IsRelease() {
		logging.Info("ipctl %s (%s)", info.Version, info.Build)
	} else {
		sha, err := app.GetCurrentSha()
		if err == nil {
			logging.Info("ipctl running from commit %s", sha)
		} else {
			logging.Info("ipctl unable to determine source")
		}
	}

	profile, err := cfg.ActiveProfile()
	if err != nil {
		terminal.Error(err, termCfg.NoColor)
		return 1
	}

	logging.Info("connection timeout is %v second(s)", profile.Timeout)

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

	// Execute the CLI command tree. This is the only place in the application
	// where CheckError should be called. All command handlers use RunE to return
	// errors which are propagated up through Cobra and handled here.
	if err := runCli(c, cfg, &termCfg).Execute(); err != nil {
		cmdutils.CheckError(err, termCfg.NoColor)
	}

	return 0
}
