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

type LocalClientHandler struct {
	runner     runners.LocalClientRunner
	runtime    *Runtime
	descriptor DescriptorMap
}

func NewLocalClientHandler(rt *Runtime) LocalClientHandler {
	return LocalClientHandler{
		runner:     runners.NewLocalClientRunner(rt.GetClient(), rt.GetConfig()),
		runtime:    rt,
		descriptor: rt.GetDescriptors()[localClientDescriptor],
	}
}

func (h LocalClientHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
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

// Commands returns a list of commands for the local client handler.
func (h LocalClientHandler) Commands() []*cobra.Command {
	logging.Trace()
	return []*cobra.Command{
		h.Show(),
	}
}

// Show returns the show-config command.
func (h LocalClientHandler) Show() *cobra.Command {
	logging.Trace()
	return h.newCommand("show-config", h.runner.ShowConfig, nil)
}
