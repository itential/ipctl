// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type ApiGetOptions struct {
	ExpectedStatusCode int
	Params             []string
}

func (o *ApiGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *ApiGetOptions) GetParams() []string {
	return o.Params
}

func (o *ApiGetOptions) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type ApiPutOptions struct {
	Data               string
	ExpectedStatusCode int
	Params             []string
}

func (o *ApiPutOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *ApiPutOptions) GetParams() []string {
	return o.Params
}

func (o *ApiPutOptions) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type ApiDeleteOptions struct {
	ExpectedStatusCode int
	Params             []string
}

func (o *ApiDeleteOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *ApiDeleteOptions) GetParams() []string {
	return o.Params
}

func (o *ApiDeleteOptions) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type ApiPostOptions struct {
	Data               string
	ExpectedStatusCode int
	Params             []string
}

func (o *ApiPostOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *ApiPostOptions) GetParams() []string {
	return o.Params
}

func (o *ApiPostOptions) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type ApiPatchOptions struct {
	Data               string
	ExpectedStatusCode int
	Params             []string
}

func (o *ApiPatchOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Data, "data", "d", o.Data, "HTTP data to include in the request")
	cmd.Flags().IntVar(&o.ExpectedStatusCode, "expected-status-code", o.ExpectedStatusCode, "Expected response status code")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *ApiPatchOptions) GetParams() []string {
	return o.Params
}

func (o *ApiPatchOptions) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}
