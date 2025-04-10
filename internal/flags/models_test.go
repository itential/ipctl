// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestModelCreateOptions(t *testing.T) {
	checkFlags(t, &ModelCreateOptions{}, []string{"description", "schema", "replace"})
}

func TestModelDeleteOptions(t *testing.T) {
	checkFlags(t, &ModelDeleteOptions{}, []string{"all", "delete-instances"})
}

func TestModelExportOptions(t *testing.T) {
	checkFlags(t, &ModelExportOptions{}, []string{"expand"})
}

func TestModelImportOptions(t *testing.T) {
	checkFlags(t, &ModelImportOptions{}, []string{"all", "skip-checks"})
}
