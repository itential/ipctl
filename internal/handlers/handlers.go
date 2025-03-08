// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/google/uuid"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/spf13/cobra"
)

type Handler struct {
	Runtime     *Runtime
	Config      *config.Config
	Client      client.Client
	Descriptors Descriptors
}

type Runtime struct {
	Config  *config.Config
	Verbose bool
}

var commands []any

func NewHandler(iapClient client.Client, cfg *config.Config) Handler {
	descriptors := loadDescriptors()

	register(
		// Automation Studio handlers
		NewProjectHandler(iapClient, cfg, descriptors),
		NewWorkflowHandler(iapClient, cfg, descriptors),
		NewTransformationHandler(iapClient, cfg, descriptors),
		NewJsonFormHandler(iapClient, cfg, descriptors),
		NewCommandTemplateHandler(iapClient, cfg, descriptors),
		NewAnalyticTemplateHandler(iapClient, cfg, descriptors),
		NewTemplateHandler(iapClient, cfg, descriptors),

		// Operations Manager Handlers
		NewAutomationHandler(iapClient, cfg, descriptors),

		// Admin Essentials handlers
		NewAccountHandler(iapClient, cfg, descriptors),
		NewProfileHandler(iapClient, cfg, descriptors),
		NewRoleHandler(iapClient, cfg, descriptors),
		NewRoleTypesHandler(iapClient, cfg, descriptors),
		NewGroupHandler(iapClient, cfg, descriptors),
		NewMethodHandler(iapClient, cfg, descriptors),
		NewViewHandler(iapClient, cfg, descriptors),
		NewPrebuiltHandler(iapClient, cfg, descriptors),
		NewIntegrationModelHandler(iapClient, cfg, descriptors),
		NewIntegrationHandler(iapClient, cfg, descriptors),
		NewAdapterHandler(iapClient, cfg, descriptors),
		NewAdapterModelHandler(iapClient, cfg, descriptors),
		NewTagHandler(iapClient, cfg, descriptors),
		NewApplicationHandler(iapClient, cfg, descriptors),

		// Configuration Manager handlers
		NewDeviceHandler(iapClient, cfg, descriptors),
		NewDeviceGroupHandler(iapClient, cfg, descriptors),
		NewConfigurationParserHandler(iapClient, cfg, descriptors),
		NewGoldenConfigHandler(iapClient, cfg, descriptors),

		// Lifecycle Manager handlers
		NewModelHandler(iapClient, cfg, descriptors),

		NewServerHandler(iapClient, cfg, descriptors),

		NewLocalClientHandler(iapClient, cfg, descriptors),
	)

	if cfg.MongoUri != "" {
		NewLocalAAAHandler(iapClient, cfg, descriptors)
	}

	return Handler{
		Runtime: &Runtime{
			Config: cfg,
		},
		Config:      cfg,
		Client:      iapClient,
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

func (h Handler) ApiCommands() []*cobra.Command {
	logger.Trace()
	handler := NewApiHandler(h.Client, h.Config, h.Descriptors)
	var commands = []*cobra.Command{
		handler.Get(h.Runtime),
		handler.Delete(h.Runtime),
		handler.Put(h.Runtime),
		handler.Post(h.Runtime),
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

func (h Handler) LocalAAACommands() []*cobra.Command {
	p, err := h.Config.ActiveProfile()
	if err != nil {
		logger.Fatal(err, "")
	}

	if p.MongoUrl != "" {
		handler := NewLocalAAAHandler(h.Client, h.Config, h.Descriptors)
		logger.Info("adding LocalAAA commands")
		return []*cobra.Command{
			handler.Get(h.Runtime),
			handler.Create(h.Runtime),
			handler.Delete(h.Runtime),
		}
	}

	return nil
}

func (h Handler) LocalClientCommands() []*cobra.Command {
	handler := NewLocalClientHandler(h.Client, h.Config, h.Descriptors)
	var commands = []*cobra.Command{
		handler.Show(h.Runtime),
	}
	return commands
}
