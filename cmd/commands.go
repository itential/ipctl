// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/handlers"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	Name       string
	Group      string
	Run        func() []*cobra.Command
	Descriptor string
}

// addRootCommand adds a new top level command to the application.  Root
// commands typically do not implement functionality, rather provide a command
// tree for more specific commands.
func addRootCommand(cmd *cobra.Command, r handlers.Runtime, title string, f func(handlers.Runtime, string) []*cobra.Command) {
	id := uuid.New().String()
	children := f(r, id)
	if len(children) > 0 {
		cmd.AddGroup(&cobra.Group{ID: id, Title: title})
		for _, ele := range children {
			cmd.AddCommand(ele)
		}
	}
}

// makeRootCommand will create a new root command for the application.  Root
// commands are top level commands that implement addiitonal subcommands and
// therefore do not direclty perform any actions.  The function accepts a
// single argument `rootCommands` which is an array of RootCommand instances.
func makeRootCommand(rootCommands []RootCommand) []*cobra.Command {
	descriptors := loadDescriptors("descriptors")

	var commands []*cobra.Command

	for _, ele := range rootCommands {
		if desc, exists := descriptors[ele.Descriptor]; exists {
			c := makeChildCommand(ele, desc)
			if c != nil {
				commands = append(commands, c)
			}
		} else {
			logger.Fatal(fmt.Errorf("missing root command descriptor: %s", ele.Descriptor), "")
		}
	}

	return commands
}

// makeChildCommand creates a single command attached to a root command.  Child
// commands are typically handed off to a handler for further implementation of
// the command action.
func makeChildCommand(root RootCommand, desc map[string]cmdutils.Descriptor) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     root.Name,
		GroupID: root.Group,

		Short: strings.Split(desc[root.Name].Description, "\n")[0],
		Long:  desc[root.Name].Description,

		Example: desc[root.Name].Example,

		Hidden: desc[root.Name].Hidden,
	}

	if desc[root.Name].IncludeGroups {
		cmd.AddGroup(
			&cobra.Group{ID: "admin-essentials", Title: "Admin Essentials Commands:"},
			&cobra.Group{ID: "automation-studio", Title: "Automation Studio Commands:"},
			&cobra.Group{ID: "configuration-manager", Title: "Configuration Manager Commands:"},
			&cobra.Group{ID: "operations-manager", Title: "Operations Manager Commands:"},
			&cobra.Group{ID: "lifecycle-manager", Title: "Lifecycle Manager Commands:"},
		)
	}

	children := root.Run()

	if len(children) == 0 {
		return nil
	}

	cmd.AddCommand(children...)

	return cmd
}

// assetCommands define the aggregate set of commands for working with assets
func assetCommands(r handlers.Runtime, id string) []*cobra.Command {
	h := handlers.NewHandler(r)
	return makeRootCommand([]RootCommand{
		RootCommand{"get", id, h.GetCommands, "asset"},
		RootCommand{"describe", id, h.DescribeCommands, "asset"},

		RootCommand{"create", id, h.CreateCommands, "asset"},
		RootCommand{"delete", id, h.DeleteCommands, "asset"},

		RootCommand{"copy", id, h.CopyCommands, "asset"},

		RootCommand{"clear", id, h.ClearCommands, "asset"},

		RootCommand{"edit", id, h.EditCommands, "asset"},

		RootCommand{"import", id, h.ImportCommands, "asset"},
		RootCommand{"export", id, h.ExportCommands, "asset"},
	})
}

// platformCommands define the set of commands that can be performed on a
// speific server instance.
func platformCommands(r handlers.Runtime, id string) []*cobra.Command {
	h := handlers.NewHandler(r)
	return makeRootCommand([]RootCommand{
		RootCommand{"api", id, h.ApiCommands, "platform"},
		RootCommand{"inspect", id, h.InspectCommands, "platform"},
		RootCommand{"start", id, h.StartCommands, "platform"},
		RootCommand{"stop", id, h.StopCommands, "platform"},
		RootCommand{"restart", id, h.RestartCommands, "platform"},
	})
}

// dataserCommands provide a set of commands for performing batch operations on
// specific asset types
func datasetCommands(r handlers.Runtime, id string) []*cobra.Command {
	h := handlers.NewHandler(r)
	return makeRootCommand([]RootCommand{
		RootCommand{"load", id, h.LoadCommands, "dataset"},
		RootCommand{"dump", id, h.DumpCommands, "dataset"},
	})
}

// pluginCommands are commands that extend the functionality of the
// application.
func pluginCommands(r handlers.Runtime, id string) []*cobra.Command {
	localAAAHandler := handlers.NewLocalAAAHandler(r)
	localClientHandler := handlers.NewLocalClientHandler(r)

	return makeRootCommand([]RootCommand{
		RootCommand{"local-aaa", id, localAAAHandler.Commands, "localaaa"},
		RootCommand{"client", id, localClientHandler.Commands, "localclient"},
	})
}
