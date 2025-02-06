// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type AssetImportCommon struct {
	Force   bool
	Replace bool
}

func (o *AssetImportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Force, "force", o.Force, "Force overwriting existing assets")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace an existing asset (if it exists)")
}

type AssetExportCommon struct {
	Path string
}

func (o *AssetExportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path where asset should be exported to")
}

type AssetCopyCommon struct {
	To      string
	From    string
	Replace bool
}

func (o *AssetCopyCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.To, "to", o.To, "Destination server to copy the asset to (REQUIRED)")
	cmd.MarkFlagRequired("to")

	cmd.Flags().StringVar(&o.From, "from", o.From, "Source server to copy the asset from (REQUIRED)")
	cmd.MarkFlagRequired("from")

	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace asset on destination server if it exists")
}

type AssetPullCommon struct {
	Path      string
	Replace   bool
	Reference string
}

func (o *AssetPullCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path to the file in the repository")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace an existing asset (if it exists)")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
}

type AssetPushCommon struct {
	Path      string
	Message   string
	Reference string
}

func (o *AssetPushCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path to the file in the repository")
	cmd.Flags().StringVar(&o.Message, "message", o.Message, "Commit message")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
}
