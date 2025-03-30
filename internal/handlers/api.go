// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/spf13/cobra"
)

type ApiHandler struct {
	Runner     runners.ApiRunner
	Descriptor DescriptorMap
}

func NewApiHandler(r Runtime, desc Descriptors) ApiHandler {
	return ApiHandler{
		Runner:     runners.NewApiRunner(r.Client),
		Descriptor: desc[apiDescriptor],
	}

}

func (h ApiHandler) newCommand(key string, runtime *Runtime, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	r := NewCommandRunner(
		key,
		h.Descriptor,
		runner,
		runtime,
		options,
		opts...,
	)
	return NewCommand(r)
}

// Adds the `api get <path>` command
func (h ApiHandler) Get(runtime *Runtime) *cobra.Command {
	cmd := h.newCommand("get", runtime, h.Runner.Get, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

// Adds the `api delete <path>` command
func (h ApiHandler) Delete(runtime *Runtime) *cobra.Command {
	options := flags.ApiDeleteOptions{}
	cmd := h.newCommand("delete", runtime, h.Runner.Delete, &options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api put <path>` command
func (h ApiHandler) Put(runtime *Runtime) *cobra.Command {
	options := &flags.ApiPutOptions{}
	cmd := h.newCommand("put", runtime, h.Runner.Put, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api post <path>` command
func (h ApiHandler) Post(runtime *Runtime) *cobra.Command {
	options := &flags.ApiPostOptions{}
	cmd := h.newCommand("post", runtime, h.Runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api path <path>` command
func (h ApiHandler) Patch(runtime *Runtime) *cobra.Command {
	options := &flags.ApiPatchOptions{}
	cmd := h.newCommand("patch", runtime, h.Runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}
