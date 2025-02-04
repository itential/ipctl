// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestAutomationCreateOptions(t *testing.T) {
	checkFlags(t, &AutomationCreateOptions{}, []string{"description", "replace"})
}

func TestAutomationImportOptions(t *testing.T) {
	checkFlags(t, &AutomationImportOptions{}, []string{
		"disable-component-check",
		"disable-group-read-check",
		"disable-group-write-check",
		"disable-group-exists-check",
	})
}
