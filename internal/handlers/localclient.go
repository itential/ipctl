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
	Runtime    Runtime
	Descriptor DescriptorMap
}

func NewLocalClientHandler(r Runtime) LocalClientHandler {
	return LocalClientHandler{
		Runner:     runners.NewLocalClientRunner(r.Client, r.Config),
		Runtime:    r,
		Descriptor: r.Descriptors[localClientDescriptor],
	}

}

func (h LocalClientHandler) newCommand(key string, runner runners.RunnerFunc, options flags.Flagger, opts ...CommandRunnerOption) *cobra.Command {
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

func (h LocalClientHandler) Commands() []*cobra.Command {
	logger.Trace()
	return []*cobra.Command{
		h.Show(),
	}
}

func (h LocalClientHandler) Show() *cobra.Command {
	logger.Trace()
	return h.newCommand("show-config", h.Runner.ShowConfig, nil)
}
