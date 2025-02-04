// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestPrebuiltDeleteOptions(t *testing.T) {
	checkFlags(t, &PrebuiltDeleteOptions{}, []string{"all"})
}

func TestPrebuiltImportOptions(t *testing.T) {
	checkFlags(t, &PrebuiltImportOptions{}, []string{"path"})
}

func TestPrebuiltExportOptions(t *testing.T) {
	checkFlags(t, &PrebuiltExportOptions{}, []string{"expand"})
}
