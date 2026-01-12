// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Paramer provides access to custom query parameters.
// This interface is implemented by all common flag structures that support
// the --params flag for passing additional query parameters to API requests.
type Paramer interface {
	GetParams() []string
	ParseParams() (map[string]string, error)
}

type Gitter interface {
	GetRepository() string
	GetReference() string
	GetPrivateKeyFile() string
}

type Committer interface {
	GetMessage() string
	GetPath() string
}

// ParseParams parses a slice of "key=value" strings into a map of query parameters.
// Each element in the params slice must be in the format "key=value".
//
// Parameters:
//   - params: A slice of strings in "key=value" format
//
// Returns:
//   - A map of parsed key-value pairs, or an error if any param is malformed
//
// Example:
//
//	params := []string{"limit=10", "offset=20", "status=active"}
//	result, err := ParseParams(params)
//	// result: map[string]string{"limit": "10", "offset": "20", "status": "active"}
func ParseParams(params []string) (map[string]string, error) {
	if len(params) == 0 {
		return nil, nil
	}

	result := make(map[string]string, len(params))
	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid param format %q: expected key=value", param)
		}
		if parts[0] == "" {
			return nil, fmt.Errorf("invalid param format %q: key cannot be empty", param)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

type AssetImportCommon struct {
	Replace        bool
	Repository     string
	Reference      string
	PrivateKeyFile string
	Params         []string
}

func (o *AssetImportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace an existing asset (if it exists)")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
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

func (o *AssetImportCommon) GetParams() []string {
	return o.Params
}

func (o *AssetImportCommon) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type AssetExportCommon struct {
	Path           string
	Repository     string
	Reference      string
	PrivateKeyFile string
	Message        string
	Params         []string
}

func (o *AssetExportCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path where asset should be exported to")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringVar(&o.Message, "message", o.Message, "Git commit message")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
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

func (o *AssetExportCommon) GetParams() []string {
	return o.Params
}

func (o *AssetExportCommon) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type AssetDumpCommon struct {
	Path           string
	Repository     string
	Reference      string
	PrivateKeyFile string
	Message        string
	Params         []string
}

func (o *AssetDumpCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path where asset should be exported to")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringVar(&o.Message, "message", o.Message, "Git commit message")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
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

func (o *AssetDumpCommon) GetParams() []string {
	return o.Params
}

func (o *AssetDumpCommon) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type AssetLoadCommon struct {
	Repository     string
	Reference      string
	PrivateKeyFile string
	Params         []string
}

func (o *AssetLoadCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "private-key-file", o.PrivateKeyFile, "Path to Git private key")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
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

func (o *AssetLoadCommon) GetParams() []string {
	return o.Params
}

func (o *AssetLoadCommon) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}

type AssetCopyCommon struct {
	To      string
	From    string
	Replace bool
	Params  []string
}

func (o *AssetCopyCommon) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.To, "to", o.To, "Destination server to copy the asset to (REQUIRED)")
	cmd.MarkFlagRequired("to")

	cmd.Flags().StringVar(&o.From, "from", o.From, "Source server to copy the asset from (REQUIRED)")
	cmd.MarkFlagRequired("from")

	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace asset on destination server if it exists")
	cmd.Flags().StringArrayVar(&o.Params, "params", o.Params, "Query parameters in key=value format (can be specified multiple times)")
}

func (o *AssetCopyCommon) GetParams() []string {
	return o.Params
}

func (o *AssetCopyCommon) ParseParams() (map[string]string, error) {
	return ParseParams(o.Params)
}
