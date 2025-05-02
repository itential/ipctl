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
	Runner     runners.LocalAAARunner
	Runtime    Runtime
	Descriptor DescriptorMap
}

func NewLocalAAAHandler(r Runtime) LocalAAAHandler {
	return LocalAAAHandler{
		Runner:     runners.NewLocalAAARunner(r.Client, r.Config),
		Runtime:    r,
		Descriptor: r.Descriptors[localAAADescriptor],
	}

}

func (h LocalAAAHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	r := NewCommandRunner(
		key,
		h.Descriptor,
		runner,
		&h.Runtime,
		nil,
		opts...,
	)
	r.Options = options
	return NewCommand(r)
}

// Commands returns a list of commands that are attached to the root command
// for this handler
func (h LocalAAAHandler) Commands() []*cobra.Command {
	logger.Trace()

	p, err := h.Runtime.Config.ActiveProfile()
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
	return h.newCommand("get-accounts", h.Runner.GetAccounts, nil)
}

func (h LocalAAAHandler) getGroups() *cobra.Command {
	logger.Trace()
	return h.newCommand("get-groups", h.Runner.GetGroups, nil)
}

/*
*******************************************************************************
Create commands
*******************************************************************************
*/

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
	cmd := h.newCommand("create-account", h.Runner.CreateAccount, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

func (h LocalAAAHandler) createGroup() *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("create-group", h.Runner.CreateGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

/*
*******************************************************************************
Delete commands
*******************************************************************************
*/

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
	cmd := h.newCommand("delete-account", h.Runner.DeleteAccount, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func (h LocalAAAHandler) deleteGroup() *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("delete-group", h.Runner.DeleteGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}
