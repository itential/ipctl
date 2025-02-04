// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/pkg/logger"
)

type CommandRunner struct {
	runner Runner
}

func NewCommandRunner(runner Runner) CommandRunner {
	return CommandRunner{runner: runner}
}

func (c CommandRunner) Get(req Request) (*Response, error) {
	logger.Trace()
	return c.runner.(Reader).Get(req)
}

func (c CommandRunner) Describe(req Request) (*Response, error) {
	logger.Trace()
	return c.runner.(Reader).Describe(req)
}
