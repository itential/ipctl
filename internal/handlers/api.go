// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/runners"
	"github.com/spf13/cobra"
)

type ApiHandler struct {
	runner     runners.ApiRunner
	runtime    *Runtime
	descriptor DescriptorMap
}

func NewApiHandler(rt *Runtime) ApiHandler {
	return ApiHandler{
		runner:     runners.NewApiRunner(rt.GetClient()),
		runtime:    rt,
		descriptor: rt.GetDescriptors()[apiDescriptor],
	}
}

func (h ApiHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	cr := NewCommandRunner(
		key,
		h.descriptor,
		runner,
		h.runtime,
		options,
		opts...,
	)
	return NewCommand(cr)
}

func (h ApiHandler) Commands() []*cobra.Command {
	logging.Trace()
	return []*cobra.Command{
		h.Get(),
		h.Delete(),
		h.Put(),
		h.Post(),
		h.Patch(),
	}
}

// Get adds the `api get <path> ...` command.
func (h ApiHandler) Get() *cobra.Command {
	cmd := h.newCommand("get", h.runner.Get, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

// Delete adds the `api delete <path> ...` command.
func (h ApiHandler) Delete() *cobra.Command {
	options := flags.ApiDeleteOptions{}
	cmd := h.newCommand("delete", h.runner.Delete, &options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Put adds the `api put <path> ...` command.
func (h ApiHandler) Put() *cobra.Command {
	options := &flags.ApiPutOptions{}
	cmd := h.newCommand("put", h.runner.Put, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Post adds the `api post <path> ...` command.
func (h ApiHandler) Post() *cobra.Command {
	options := &flags.ApiPostOptions{}
	cmd := h.newCommand("post", h.runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Patch adds the `api patch <path> ...` command.
func (h ApiHandler) Patch() *cobra.Command {
	options := &flags.ApiPatchOptions{}
	cmd := h.newCommand("patch", h.runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}
