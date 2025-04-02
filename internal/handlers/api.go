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

type ApiHandler struct {
	Runner     runners.ApiRunner
	Runtime    Runtime
	Descriptor DescriptorMap
}

func NewApiHandler(r Runtime) ApiHandler {
	return ApiHandler{
		Runner:     runners.NewApiRunner(r.Client),
		Runtime:    r,
		Descriptor: r.Descriptors[apiDescriptor],
	}

}

func (h ApiHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
	r := NewCommandRunner(
		key,
		h.Descriptor,
		runner,
		&h.Runtime,
		options,
		opts...,
	)
	return NewCommand(r)
}

func (h ApiHandler) Commands() []*cobra.Command {
	logger.Trace()
	return []*cobra.Command{
		h.Get(),
		h.Delete(),
		h.Put(),
		h.Post(),
		h.Patch(),
	}
}

// Adds the `api get <path> ...` command
func (h ApiHandler) Get() *cobra.Command {
	cmd := h.newCommand("get", h.Runner.Get, nil)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

// Adds the `api delete <path> ...` command
func (h ApiHandler) Delete() *cobra.Command {
	options := flags.ApiDeleteOptions{}
	cmd := h.newCommand("delete", h.Runner.Delete, &options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api put <path> ...` command
func (h ApiHandler) Put() *cobra.Command {
	options := &flags.ApiPutOptions{}
	cmd := h.newCommand("put", h.Runner.Put, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api post <path> ...` command
func (h ApiHandler) Post() *cobra.Command {
	options := &flags.ApiPostOptions{}
	cmd := h.newCommand("post", h.Runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}

// Adds the `api path <path> ...` command
func (h ApiHandler) Patch() *cobra.Command {
	options := &flags.ApiPatchOptions{}
	cmd := h.newCommand("patch", h.Runner.Post, options)
	cmd.Args = cobra.ExactArgs(1)
	options.Flags(cmd)
	return cmd
}
