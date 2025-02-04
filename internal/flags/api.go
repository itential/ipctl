// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type ApiPutOptions struct {
	Data               string
	ExpectedStatusCode int
}

func (o *ApiPutOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
}

type ApiDeleteOptions struct {
	ExpectedStatusCode int
}

func (o *ApiDeleteOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
}

type ApiPostOptions struct {
	Data               string
	ExpectedStatusCode int
}

func (o *ApiPostOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
}

type ApiPatchOptions struct {
	Data               string
	ExpectedStatusCode int
}

func (o *ApiPatchOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
}
