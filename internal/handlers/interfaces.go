// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/spf13/cobra"
)

type Reader interface {
	Get(*Runtime) *cobra.Command
	Describe(*Runtime) *cobra.Command
}

type Writer interface {
	Create(*Runtime) *cobra.Command
	Delete(*Runtime) *cobra.Command
	Clear(*Runtime) *cobra.Command
}

type Copier interface {
	Copy(*Runtime) *cobra.Command
}

type Editor interface {
	Edit(*Runtime) *cobra.Command
}

type Importer interface {
	Import(*Runtime) *cobra.Command
}

type Exporter interface {
	Export(*Runtime) *cobra.Command
}

type Controller interface {
	Start(*Runtime) *cobra.Command
	Stop(*Runtime) *cobra.Command
	Restart(*Runtime) *cobra.Command
}

type Inspector interface {
	Inspect(*Runtime) *cobra.Command
}
