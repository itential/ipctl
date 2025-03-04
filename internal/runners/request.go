// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import "github.com/itential/ipctl/pkg/config"

type Request struct {
	Args    []string
	Common  any
	Options any
	Runner  Runner
	Config  *config.Config
}
