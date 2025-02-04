// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestCommandTemplateCreateOptions(t *testing.T) {
	checkFlags(t, &CommandTemplateCreateOptions{}, []string{"replace"})
}

func TestCommandTemplateGetOptions(t *testing.T) {
	checkFlags(t, &CommandTemplateGetOptions{}, []string{"all"})
}
