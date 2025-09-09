// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetImportCommon(t *testing.T) {
	checkFlags(t, &AssetImportCommon{}, []string{"replace", "repository", "reference", "private-key-file"})
}

func TestAssetExportCommon(t *testing.T) {
	checkFlags(t, &AssetExportCommon{}, []string{"path", "repository", "reference", "private-key-file", "message"})
}

func TestAssetDumpCommon(t *testing.T) {
	checkFlags(t, &AssetDumpCommon{}, []string{"path", "repository", "reference", "private-key-file", "message"})
}

func TestAssetLoadCommon(t *testing.T) {
	checkFlags(t, &AssetLoadCommon{}, []string{"repository", "reference", "private-key-file"})
}

func TestAssetCopyCommon(t *testing.T) {
	checkFlags(t, &AssetCopyCommon{}, []string{"to", "from", "replace"})
}

func TestAssetImportCommonGetters(t *testing.T) {
	common := &AssetImportCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "main",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "main", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "", common.GetPath())
}

func TestAssetExportCommonGetters(t *testing.T) {
	common := &AssetExportCommon{
		Path:           "/export/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "develop",
		PrivateKeyFile: "/path/to/key",
		Message:        "Export commit message",
	}

	assert.Equal(t, "/export/path", common.GetPath())
	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "develop", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "Export commit message", common.GetMessage())
}

func TestAssetDumpCommonGetters(t *testing.T) {
	common := &AssetDumpCommon{
		Path:           "/dump/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "feature",
		PrivateKeyFile: "/path/to/key",
		Message:        "Dump commit message",
	}

	assert.Equal(t, "/dump/path", common.GetPath())
	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "feature", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "Dump commit message", common.GetMessage())
}

func TestAssetLoadCommonGetters(t *testing.T) {
	common := &AssetLoadCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "v1.0.0",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "v1.0.0", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
}

func TestAssetCopyCommonFields(t *testing.T) {
	common := &AssetCopyCommon{
		To:      "destination",
		From:    "source",
		Replace: true,
	}

	assert.Equal(t, "destination", common.To)
	assert.Equal(t, "source", common.From)
	assert.True(t, common.Replace)
}

func TestGitterInterface(t *testing.T) {
	var gitter Gitter = &AssetImportCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "main",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", gitter.GetRepository())
	assert.Equal(t, "main", gitter.GetReference())
	assert.Equal(t, "/path/to/key", gitter.GetPrivateKeyFile())
}

func TestCommitterInterface(t *testing.T) {
	var committer Committer = &AssetExportCommon{
		Path:           "/export/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "develop",
		PrivateKeyFile: "/path/to/key",
		Message:        "Export commit message",
	}

	assert.Equal(t, "/export/path", committer.GetPath())
	assert.Equal(t, "Export commit message", committer.GetMessage())
}
