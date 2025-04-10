// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "github.com/spf13/cobra"

var (
	typeOption = Option{Name: "type", Abbrev: "t", Usage: "Type of role"}
)

type RoleCreateOptions struct {
	AllowedMethods []string
	AllowedViews   []string
}

func (o *RoleCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.AllowedMethods, "method", o.AllowedMethods, "Role allowed method")
	cmd.Flags().StringArrayVar(&o.AllowedViews, "view", o.AllowedViews, "Role allowed view")
}

type RoleGetOptions struct {
	All  bool
	Type string
}

func (o *RoleGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Display all roles including builtin roles")
	cmd.Flags().StringVar(&o.Type, "type", o.Type, "Display only roles of a certain type")
}

type RoleDescribeOptions struct {
	Type string
}

func (o *RoleDescribeOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Type, typeOption.Name, typeOption.Abbrev, o.Type, typeOption.Usage)
}

type RoleExportOptions struct {
	Type string
}

func (o *RoleExportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Type, typeOption.Name, typeOption.Abbrev, o.Type, typeOption.Usage)
}
