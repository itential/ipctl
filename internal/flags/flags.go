// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type Flagger interface {
	Flags(cmd *cobra.Command)
}

type Option struct {
	Name   string
	Abbrev string
	Usage  string
}
