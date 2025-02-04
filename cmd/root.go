// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmd

import (
	"context"
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

const description = `Manage Itential Platform

  Find more information at: https://docs.itential.com
`

func loadCommands(cmd *cobra.Command, h handlers.Handler, cfg *config.Config) {
	addCommandGroup(cmd, h, "Asset Commands:", assetCommands)
	addCommandGroup(cmd, h, "Application Commands:", applicationCommands)
	addCommandGroup(cmd, h, "Import/Export Commands:", importExportCommands)
	addCommandGroup(cmd, h, "Platform Commands:", platformCommands)
	addCommandGroup(cmd, h, "Repository Commands:", repoCommands)
	addCommandGroup(cmd, h, "Plugin Commands:", pluginCommands)
}

func versionCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the client version information",
		Run: func(cmd *cobra.Command, args []string) {
			terminal.Display("version: %s", metadata.Version)
			terminal.Display("commit: %s", metadata.Sha)
			e, _ := os.Executable()
			terminal.Display("executable: %s", path.Dir(e))
			terminal.Display("")
		},
	}
	return cmd
}

func Do(iapClient *client.IapClient, cfg *config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ipctl",
		Short: strings.Split(description, "\n")[0],
		Long:  description,
	}

	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	h := handlers.NewHandler(iapClient, cfg)

	cmd.PersistentFlags().BoolVar(&h.Runtime.Verbose, "verbose", h.Runtime.Verbose, "Enable verbose output")
	cmd.PersistentFlags().StringVar(&h.Runtime.Config.DefaultOutput, "output", h.Runtime.Config.DefaultOutput, "Output format")

	// Note: Values are read during /pkg/config's initialization
	cmd.PersistentFlags().String("config", "", "Path to the configuration file")
	cmd.PersistentFlags().String("profile", "", "Connection profile to use")

	loadCommands(cmd, h, cfg)

	cmd.AddCommand(versionCommand())

	return cmd
}

func Execute() int {
	cfg := config.NewConfig(nil, nil, "", "", "")
	logger.InitializeLogger(cfg)

	profile, err := cfg.ActiveProfile()
	if err != nil {
		cmdutils.CheckError(err, cfg.TerminalNoColor)
	}

	logger.Info("connection timeout is %v second(s)", profile.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(profile.Timeout)*time.Second)
	defer cancel()

	c := client.New(ctx, profile)

	if err := Do(c, cfg).Execute(); err != nil {
		cmdutils.CheckError(err, cfg.TerminalNoColor)
	}

	return 0
}
