// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestJsonFormCreateOptions(t *testing.T) {
	checkFlags(t, &JsonFormCreateOptions{}, []string{"description", "replace"})
}

func TestJsonFormGetOptions(t *testing.T) {
	checkFlags(t, &JsonFormGetOptions{}, []string{"all"})
}
