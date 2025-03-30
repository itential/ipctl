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
	Descriptor DescriptorMap
}

func NewLocalAAAHandler(r Runtime, desc Descriptors) LocalAAAHandler {
	return LocalAAAHandler{
		Runner:     runners.NewLocalAAARunner(r.Client, r.Config),
		Descriptor: desc[localAAADescriptor],
	}

}

func (h LocalAAAHandler) newCommand(key string, runtime *Runtime, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	r := NewCommandRunner(
		key,
		h.Descriptor,
		runner,
		runtime,
		nil,
		opts...,
	)
	r.Options = options
	return NewCommand(r)
}

func (h LocalAAAHandler) Get(runtime *Runtime) *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more local-aaa resources",
	}

	cmd.AddCommand(
		h.getAccounts(runtime),
		h.getGroups(runtime),
	)

	return cmd
}

func (h LocalAAAHandler) getAccounts(runtime *Runtime) *cobra.Command {
	logger.Trace()
	return h.newCommand("get-accounts", runtime, h.Runner.GetAccounts, nil)
}

func (h LocalAAAHandler) getGroups(runtime *Runtime) *cobra.Command {
	logger.Trace()
	return h.newCommand("get-groups", runtime, h.Runner.GetGroups, nil)
}

func (h LocalAAAHandler) Create(runtime *Runtime) *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new local-aaa resource",
	}

	cmd.AddCommand(
		h.createAccount(runtime),
		h.createGroup(runtime),
	)

	return cmd
}

func (h LocalAAAHandler) createAccount(runtime *Runtime) *cobra.Command {
	logger.Trace()
	options := &flags.LocalAAAOptions{}
	cmd := h.newCommand("create-account", runtime, h.Runner.CreateAccount, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

func (h LocalAAAHandler) createGroup(runtime *Runtime) *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("create-group", runtime, h.Runner.CreateGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func (h LocalAAAHandler) Delete(runtime *Runtime) *cobra.Command {
	logger.Trace()

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a local AAA resource",
	}

	cmd.AddCommand(
		h.deleteAccount(runtime),
		h.deleteGroup(runtime),
	)

	return cmd
}

func (h LocalAAAHandler) deleteAccount(runtime *Runtime) *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("delete-account", runtime, h.Runner.DeleteAccount, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func (h LocalAAAHandler) deleteGroup(runtime *Runtime) *cobra.Command {
	logger.Trace()
	cmd := h.newCommand("delete-group", runtime, h.Runner.DeleteGroup, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}
