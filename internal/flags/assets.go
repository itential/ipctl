// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type Gitter interface {
	GetRepository() string
	GetReference() string
	GetPrivateKeyFile() string
}

type Committer interface {
	GetMessage() string
	GetPath() string
}

type AssetImportCommon struct {
	Replace        bool
	Repository     string
	Reference      string
	PrivateKeyFile string
}

func (o *AssetImportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace an existing asset (if it exists)")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
}

func (o *AssetImportCommon) GetPath() string {
	return ""
}

func (o *AssetImportCommon) GetRepository() string {
	return o.Repository
}

func (o *AssetImportCommon) GetReference() string {
	return o.Reference
}

func (o *AssetImportCommon) GetPrivateKeyFile() string {
	return o.PrivateKeyFile
}

type AssetExportCommon struct {
	Path           string
	Repository     string
	Reference      string
	PrivateKeyFile string
	Message        string
}

func (o *AssetExportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path where asset should be exported to")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringVar(&o.Message, "message", o.Message, "Git commit message")
}

func (o *AssetExportCommon) GetPath() string {
	return o.Path
}

func (o *AssetExportCommon) GetRepository() string {
	return o.Repository
}

func (o *AssetExportCommon) GetReference() string {
	return o.Reference
}

func (o *AssetExportCommon) GetPrivateKeyFile() string {
	return o.PrivateKeyFile
}

func (o *AssetExportCommon) GetMessage() string {
	return o.Message
}

type AssetDumpCommon struct {
	Path           string
	Repository     string
	Reference      string
	PrivateKeyFile string
	Message        string
}

func (o *AssetDumpCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path where asset should be exported to")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringVar(&o.Message, "message", o.Message, "Git commit message")
}

func (o *AssetDumpCommon) GetPath() string {
	return o.Path
}

func (o *AssetDumpCommon) GetRepository() string {
	return o.Repository
}

func (o *AssetDumpCommon) GetReference() string {
	return o.Reference
}

func (o *AssetDumpCommon) GetPrivateKeyFile() string {
	return o.PrivateKeyFile
}

func (o *AssetDumpCommon) GetMessage() string {
	return o.Message
}

type AssetLoadCommon struct {
	Repository     string
	Reference      string
	PrivateKeyFile string
}

func (o *AssetLoadCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
}

func (o *AssetLoadCommon) GetRepository() string {
	return o.Repository
}

func (o *AssetLoadCommon) GetReference() string {
	return o.Reference
}

func (o *AssetLoadCommon) GetPrivateKeyFile() string {
	return o.PrivateKeyFile
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
