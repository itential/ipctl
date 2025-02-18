// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestAssetImportCommon(t *testing.T) {
	checkFlags(t, &AssetImportCommon{}, []string{"replace"})
}

func TestAssetExportCommon(t *testing.T) {
	checkFlags(t, &AssetExportCommon{}, []string{"path"})
}

func TestAssetCopyCommon(t *testing.T) {
	checkFlags(t, &AssetCopyCommon{}, []string{"to", "from", "replace"})
}

func TestAssetPullCommon(t *testing.T) {
	checkFlags(t, &AssetPullCommon{}, []string{"path", "replace", "reference"})
}

func TestAssetPushCommon(t *testing.T) {
	checkFlags(t, &AssetPushCommon{}, []string{"path", "message", "reference"})
}
