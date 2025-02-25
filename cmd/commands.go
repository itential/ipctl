// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmd

import (
	"strings"

	"github.com/google/uuid"
	"github.com/itential/ipctl/internal/handlers"
	"github.com/spf13/cobra"
)

type GroupCommand struct {
	Use        string
	Id         string
	Run        func() []*cobra.Command
	Descriptor string
}

func addCommandGroup(cmd *cobra.Command, h handlers.Handler, title string, f func(handlers.Handler, string) []*cobra.Command) {
	id := uuid.New().String()
	children := f(h, id)
	if len(children) > 0 {
		cmd.AddGroup(&cobra.Group{ID: id, Title: title})
		for _, ele := range children {
			cmd.AddCommand(ele)
		}
	}
}

func makeGroupCommand(groupCommands []GroupCommand) []*cobra.Command {
	var commands []*cobra.Command

	for _, ele := range groupCommands {
		c := makeCommand(
			ele.Use,
			ele.Id,
			ele.Run,
			ele.Descriptor,
		)

		if c != nil {
			commands = append(commands, c)
		}
	}

	return commands
}

func makeCommand(name, group string, f func() []*cobra.Command, desc string) *cobra.Command {
	descriptors := loadDescriptors("descriptors")[desc]

	var cmd = &cobra.Command{
		Use:     name,
		GroupID: group,

		Short: strings.Split(descriptors[name].Description, "\n")[0],
		Long:  descriptors[name].Description,

		Example: descriptors[name].Example,

		Hidden: descriptors[name].Hidden,
	}

	if descriptors[name].IncludeGroups {
		cmd.AddGroup(
			&cobra.Group{ID: "admin-essentials", Title: "Admin Essentials Commands:"},
			&cobra.Group{ID: "automation-studio", Title: "Automation Studio Commands:"},
			&cobra.Group{ID: "configuration-manager", Title: "Configuration Manager Commands:"},
			&cobra.Group{ID: "operations-manager", Title: "Operations Manager Commands:"},
			&cobra.Group{ID: "lifecycle-manager", Title: "Lifecycle Manager Commands:"},
		)
	}

	children := f()

	if len(children) == 0 {
		return nil
	}

	cmd.AddCommand(f()...)

	return cmd
}

func assetCommands(h handlers.Handler, id string) []*cobra.Command {
	return makeGroupCommand([]GroupCommand{
		GroupCommand{"get", id, h.GetCommands, "asset"},
		GroupCommand{"describe", id, h.DescribeCommands, "asset"},

		GroupCommand{"create", id, h.CreateCommands, "asset"},
		GroupCommand{"delete", id, h.DeleteCommands, "asset"},

		GroupCommand{"copy", id, h.CopyCommands, "asset"},

		GroupCommand{"clear", id, h.ClearCommands, "asset"},

		GroupCommand{"edit", id, h.EditCommands, "asset"},

		GroupCommand{"import", id, h.ImportCommands, "asset"},
		GroupCommand{"export", id, h.ExportCommands, "asset"},
	})
}

func platformCommands(h handlers.Handler, id string) []*cobra.Command {
	return makeGroupCommand([]GroupCommand{
		GroupCommand{"api", id, h.ApiCommands, "platform"},
		GroupCommand{"inspect", id, h.InspectCommands, "platform"},
		GroupCommand{"start", id, h.StartCommands, "platform"},
		GroupCommand{"stop", id, h.StopCommands, "platform"},
		GroupCommand{"restart", id, h.RestartCommands, "platform"},
	})
}

func repoCommands(h handlers.Handler, id string) []*cobra.Command {
	return makeGroupCommand([]GroupCommand{
		GroupCommand{"push", id, h.PushCommands, "repo"},
		GroupCommand{"pull", id, h.PullCommands, "repo"},
	})
}

func pluginCommands(h handlers.Handler, id string) []*cobra.Command {
	return makeGroupCommand([]GroupCommand{
		GroupCommand{"local-aaa", id, h.LocalAAACommands, "localaaa"},
		GroupCommand{"client", id, h.LocalClientCommands, "localclient"},
	})
}
