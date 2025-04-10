// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

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
