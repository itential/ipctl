// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type IntegrationCreateOptions struct {
	Model   string
	Replace bool
}

func (o *IntegrationCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Model, "model", o.Model, "Model of integration (REQUIRED)")
	cmd.MarkFlagRequired("model")

	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist integration if it exists")
}
