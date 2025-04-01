// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/google/uuid"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/spf13/cobra"
)

type Handler struct {
	Runtime     *Runtime
	Descriptors Descriptors
}

type Runtime struct {
	Client      client.Client
	Config      *config.Config
	Descriptors Descriptors
	Verbose     bool
}

var commands []any

func NewRuntime(c client.Client, cfg *config.Config) Runtime {
	descriptors := loadDescriptors()
	return Runtime{
		Client:      c,
		Config:      cfg,
		Descriptors: descriptors,
	}
}

func NewHandler(r Runtime) Handler {
	descriptors := loadDescriptors()

	register(
		// Automation Studio handlers
		NewProjectHandler(r, descriptors),
		NewWorkflowHandler(r, descriptors),
		NewTransformationHandler(r, descriptors),
		NewJsonFormHandler(r, descriptors),
		NewCommandTemplateHandler(r, descriptors),
		NewAnalyticTemplateHandler(r, descriptors),
		NewTemplateHandler(r, descriptors),

		// Operations Manager Handlers
		NewAutomationHandler(r, descriptors),

		// Admin Essentials handlers
		NewAccountHandler(r, descriptors),
		NewProfileHandler(r, descriptors),
		NewRoleHandler(r, descriptors),
		NewRoleTypesHandler(r, descriptors),
		NewGroupHandler(r, descriptors),
		NewMethodHandler(r, descriptors),
		NewViewHandler(r, descriptors),
		NewPrebuiltHandler(r, descriptors),
		NewIntegrationModelHandler(r, descriptors),
		NewIntegrationHandler(r, descriptors),
		NewAdapterHandler(r, descriptors),
		NewAdapterModelHandler(r, descriptors),
		NewTagHandler(r, descriptors),
		NewApplicationHandler(r, descriptors),

		// Configuration Manager handlers
		NewDeviceHandler(r, descriptors),
		NewDeviceGroupHandler(r, descriptors),
		NewConfigurationParserHandler(r, descriptors),
		NewGoldenConfigHandler(r, descriptors),

		// Lifecycle Manager handlers
		NewModelHandler(r, descriptors),

		NewServerHandler(r, descriptors),
	)

	return Handler{
		Runtime:     &r,
		Descriptors: descriptors,
	}
}

func (h Handler) GetCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Readers() {
		cmd := ele.Get(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) DescribeCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Readers() {
		cmd := ele.Describe(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) AddCommandGroup(cmd *cobra.Command, title string, f func(Handler, string) []*cobra.Command) {
	id := uuid.New().String()
	cmd.AddGroup(&cobra.Group{ID: id, Title: title})
	for _, ele := range f(h, id) {
		cmd.AddCommand(ele)
	}
}

func (h Handler) CreateCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Writers() {
		cmd := ele.Create(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) DeleteCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Writers() {
		cmd := ele.Delete(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) CopyCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Copiers() {
		cmd := ele.Copy(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) ClearCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Writers() {
		cmd := ele.Clear(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) ImportCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Importers() {
		cmd := ele.Import(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) ExportCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Exporters() {
		cmd := ele.Export(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) StartCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Controllers() {
		cmd := ele.Start(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) StopCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Controllers() {
		cmd := ele.Stop(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) RestartCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Controllers() {
		cmd := ele.Restart(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) InspectCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Inspectors() {
		cmd := ele.Inspect(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) EditCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Editors() {
		cmd := ele.Edit(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) DumpCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Dumpers() {
		cmd := ele.Dump(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (h Handler) LoadCommands() []*cobra.Command {
	var commands []*cobra.Command
	for _, ele := range Loaders() {
		cmd := ele.Load(h.Runtime)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}
