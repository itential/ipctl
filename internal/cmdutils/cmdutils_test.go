// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmdutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testYAML = `
testcmd:
  use: test
  group: testing
  description: |
    Short description line.
    This is the long description.
  example: testcmd --flag
  include_groups: true
  exact_args: 1
  disabled: false
  hidden: false
`

func TestLoadDescriptorFromString(t *testing.T) {
	descs := LoadDescriptorFromString(testYAML)

	require.Len(t, descs, 1)
	desc, ok := descs["testcmd"]
	require.True(t, ok)

	assert.Equal(t, "test", desc.Use)
	assert.Equal(t, "testing", desc.Group)
	assert.Equal(t, "testcmd --flag", desc.Example)
	assert.Equal(t, true, desc.IncludeGroups)
	assert.Equal(t, 1, desc.ExactArgs)
	assert.Equal(t, false, desc.Disabled)
	assert.Equal(t, false, desc.Hidden)
	assert.Equal(t, "Short description line.", desc.Short())
}

func TestDescriptorShort(t *testing.T) {
	desc := Descriptor{Description: "Hello there!\nLonger details here."}
	assert.Equal(t, "Hello there!", desc.Short())
}

func TestCheckErrorNil(t *testing.T) {
	assert.NotPanics(t, func() {
		CheckError(nil, false)
	})
}
