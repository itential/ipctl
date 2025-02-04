// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"
)

func TestAnalyticTemplatesGetOptions(t *testing.T) {
	checkFlags(t, &AnalyticTemplateGetOptions{}, []string{"all"})
}

func TestAnalyticTemplatesCreateOptions(t *testing.T) {
	checkFlags(t, &AnalyticTemplateCreateOptions{}, []string{"replace"})
}
