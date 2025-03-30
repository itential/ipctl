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

type LocalClientHandler struct {
	Runner     runners.LocalClientRunner
	Descriptor DescriptorMap
}

func NewLocalClientHandler(r Runtime, desc Descriptors) LocalClientHandler {
	return LocalClientHandler{
		Runner:     runners.NewLocalClientRunner(r.Client, r.Config),
		Descriptor: desc[localClientDescriptor],
	}

}

func (h LocalClientHandler) newCommand(key string, runtime *Runtime, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
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

func (h LocalClientHandler) Show(runtime *Runtime) *cobra.Command {
	logger.Trace()
	return h.newCommand("show-config", runtime, h.Runner.ShowConfig, nil)
}
