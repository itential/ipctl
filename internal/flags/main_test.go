// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkFlags(t *testing.T, flag Flagger, flags []string) {
	cmd := &cobra.Command{}
	flag.Flags(cmd)
	assert.True(t, cmd.Flags().HasFlags())
	for _, ele := range flags {
		f := cmd.Flag(ele)
		assert.NotEmpty(t, f.Name)
		assert.NotEmpty(t, f.Usage)
	}
}
