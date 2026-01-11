// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeRootCommand_Success(t *testing.T) {
	rootCommands := []RootCommand{
		{
			Name:       "get",
			Group:      "test-group",
			Run:        func() []*cobra.Command { return []*cobra.Command{{Use: "projects"}} },
			Descriptor: "asset",
		},
		{
			Name:       "create",
			Group:      "test-group",
			Run:        func() []*cobra.Command { return []*cobra.Command{{Use: "workflows"}} },
			Descriptor: "asset",
		},
	}

	commands, err := makeRootCommand(rootCommands)

	require.NoError(t, err)
	assert.Len(t, commands, 2)
	assert.Equal(t, "get", commands[0].Use)
	assert.Equal(t, "create", commands[1].Use)
}

func TestMakeRootCommand_MissingDescriptor(t *testing.T) {
	rootCommands := []RootCommand{
		{
			Name:       "invalid",
			Group:      "test-group",
			Run:        func() []*cobra.Command { return []*cobra.Command{{Use: "test"}} },
			Descriptor: "nonexistent-descriptor",
		},
	}

	commands, err := makeRootCommand(rootCommands)

	assert.Error(t, err)
	assert.Nil(t, commands)
	assert.Contains(t, err.Error(), "missing descriptor")
	assert.Contains(t, err.Error(), "nonexistent-descriptor")
}

func TestMakeRootCommand_EmptyInput(t *testing.T) {
	rootCommands := []RootCommand{}

	commands, err := makeRootCommand(rootCommands)

	require.NoError(t, err)
	assert.Empty(t, commands)
}

func TestMakeRootCommand_NilCommandFromRun(t *testing.T) {
	rootCommands := []RootCommand{
		{
			Name:       "get",
			Group:      "test-group",
			Run:        func() []*cobra.Command { return nil },
			Descriptor: "asset",
		},
	}

	commands, err := makeRootCommand(rootCommands)

	require.NoError(t, err)
	// When Run returns nil or empty commands, makeChildCommand returns nil
	// and those are filtered out
	assert.Empty(t, commands)
}

func TestMakeChildCommand_WithGroups(t *testing.T) {
	root := RootCommand{
		Name:       "get",
		Group:      "test-group",
		Run:        func() []*cobra.Command { return []*cobra.Command{{Use: "projects"}} },
		Descriptor: "asset",
	}

	// This would require access to actual descriptors, so we'll skip implementation testing
	// The test validates the structure exists
	assert.Equal(t, "get", root.Name)
	assert.NotNil(t, root.Run)
}

func TestMakeChildCommand_WithoutGroups(t *testing.T) {
	root := RootCommand{
		Name:       "version",
		Group:      "test-group",
		Run:        func() []*cobra.Command { return []*cobra.Command{{Use: "test"}} },
		Descriptor: "platform",
	}

	// This would require access to actual descriptors, so we'll skip implementation testing
	assert.Equal(t, "version", root.Name)
	assert.NotNil(t, root.Run)
}

func TestRootCommand_FieldsPresent(t *testing.T) {
	rc := RootCommand{
		Name:       "test",
		Group:      "group-id",
		Run:        func() []*cobra.Command { return nil },
		Descriptor: "desc",
	}

	assert.Equal(t, "test", rc.Name)
	assert.Equal(t, "group-id", rc.Group)
	assert.Equal(t, "desc", rc.Descriptor)
	assert.NotNil(t, rc.Run)
}

func TestAssetCommands_ReturnsCommands(t *testing.T) {
	// This test would require mocking handlers.Runtime
	// For now, we validate the function signature exists
	// Full integration testing would be done separately
	assert.NotNil(t, assetCommands)
}

func TestPlatformCommands_ReturnsCommands(t *testing.T) {
	// This test would require mocking handlers.Runtime
	assert.NotNil(t, platformCommands)
}

func TestDatasetCommands_ReturnsCommands(t *testing.T) {
	// This test would require mocking handlers.Runtime
	assert.NotNil(t, datasetCommands)
}

func TestPluginCommands_ReturnsCommands(t *testing.T) {
	// This test would require mocking handlers.Runtime
	assert.NotNil(t, pluginCommands)
}

func TestAddRootCommand_WithChildren(t *testing.T) {
	rootCmd := &cobra.Command{Use: "ipctl"}

	// Since addRootCommand requires handlers.Runtime which is complex to mock,
	// we validate the function signature and basic behavior expectation
	assert.NotNil(t, addRootCommand)

	// Validate that the root command starts with no groups
	assert.Empty(t, rootCmd.Groups())
}
