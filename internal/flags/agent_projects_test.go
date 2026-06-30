// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAgentProjectImportOptions_Flags(t *testing.T) {
	opts := &AgentProjectImportOptions{}
	cmd := &cobra.Command{}
	opts.Flags(cmd)
	// No extra flags defined; just verify the struct satisfies Flagger
	assert.NotNil(t, opts)
}

func TestAgentProjectExportOptions_Flags(t *testing.T) {
	opts := &AgentProjectExportOptions{}
	cmd := &cobra.Command{}
	opts.Flags(cmd)
	assert.NotNil(t, opts)
}
