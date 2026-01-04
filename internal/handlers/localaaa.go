// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/spf13/cobra"
)

type LocalAAAHandler struct {
	runner     runners.LocalAAARunner
	runtime    *Runtime
	descriptor DescriptorMap
}

func NewLocalAAAHandler(rt *Runtime) LocalAAAHandler {
	return LocalAAAHandler{
		runner:     runners.NewLocalAAARunner(rt.GetClient(), rt.GetConfig()),
		runtime:    rt,
		descriptor: rt.GetDescriptors()[localAAADescriptor],
	}
}

func (h LocalAAAHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	cr := NewCommandRunner(
		key,
		h.descriptor,
		runner,
		h.runtime,
		nil,
		opts...,
	)
	cr.Options = options
	return NewCommand(cr)
}

// Commands returns a list of commands that are attached to the root command
// for this handler.
func (h LocalAAAHandler) Commands() []*cobra.Command {
	logger.Trace()

	p, err := h.runtime.GetConfig().ActiveProfile()
	if err != nil {
		logger.Warn("failed to load active profile, using defaults")
	}

	if p.MongoUrl != "" {
		return []*cobra.Command{
			h.Get(),
			h.Create(),
			h.Delete(),
		}
	}

	return nil
}

/*
*******************************************************************************
Get commands
*******************************************************************************
*/

// Get returns the get command group for local-aaa resources.
func (h LocalAAAHandler) Get() *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more local-aaa resources",
	}

	cmd.AddCommand(
		h.getAccounts(),
		h.getGroups(),
	)

	return cmd
}

func (h LocalAAAHandler) getAccounts() *cobra.Command {
	logger.Trace()
	return h.newCommand("get-accounts", h.runner.GetAccounts, nil)
}

func (h LocalAAAHandler) getGroups() *cobra.Command {
	logger.Trace()
	return h.newCommand("get-groups", h.runner.GetGroups, nil)
}

/*
*******************************************************************************
Create commands
*******************************************************************************
*/

// Create returns the create command group for local-aaa resources.
func (h LocalAAAHandler) Create() *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new local-aaa resource",
	}

	cmd.AddCommand(
		h.createAccount(),
		h.createGroup(),
	)

	return cmd
}

func (h LocalAAAHandler) createAccount() *cobra.Command {
	logger.Trace()
	options := &flags.LocalAAAOptions{}
	cmd := h.newCommand("create-account", h.runner.CreateAccount, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

func (h LocalAAAHandler) createGroup() *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("create-group", h.runner.CreateGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

/*
*******************************************************************************
Delete commands
*******************************************************************************
*/

// Delete returns the delete command group for local-aaa resources.
func (h LocalAAAHandler) Delete() *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a local AAA resource",
	}

	cmd.AddCommand(
		h.deleteAccount(),
		h.deleteGroup(),
	)

	return cmd
}

func (h LocalAAAHandler) deleteAccount() *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("delete-account", h.runner.DeleteAccount, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func (h LocalAAAHandler) deleteGroup() *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("delete-group", h.runner.DeleteGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}
