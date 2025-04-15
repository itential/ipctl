// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestTemplateCreateOptions(t *testing.T) {
	checkFlags(t, &TemplateCreateOptions{}, []string{"description", "replace", "group", "type"})
}

func TestTemplateGetOptions(t *testing.T) {
	checkFlags(t, &TemplateGetOptions{}, []string{"all"})
}

func TestTemplateLoadOptions(t *testing.T) {
	checkFlags(t, &TemplateLoadOptions{}, []string{"type", "group", "include"})
}
