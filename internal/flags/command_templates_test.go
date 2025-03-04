// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestCommandTemplateGetOptions(t *testing.T) {
	checkFlags(t, &CommandTemplateGetOptions{}, []string{"all"})
}
