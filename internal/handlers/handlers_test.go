// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"strings"
	"testing"

	"github.com/itential/ipctl/internal/terminal"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	assert.NotNil(t, handler.runtime)
	assert.NotNil(t, handler.registry)
	assert.NotNil(t, handler.descriptors)
}

func TestHandler_GetCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.GetCommands()

	// Should have commands from all registered handlers that implement Reader
	assert.NotEmpty(t, commands)

	// Verify all commands are non-nil
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}

	// Check that we have expected commands (projects, workflows, etc.)
	commandUses := make(map[string]bool)
	for _, cmd := range commands {
		commandUses[cmd.Use] = true
	}

	// These resources should all have get commands
	assert.Contains(t, commandUses, "projects")
	assert.Contains(t, commandUses, "workflows")
	assert.Contains(t, commandUses, "automations")
}

func TestHandler_DescribeCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.DescribeCommands()

	// Should have describe commands from all registered handlers that implement Reader
	assert.NotEmpty(t, commands)

	// Verify all commands are non-nil
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
		// Describe commands should require exactly 1 argument
		assert.NotNil(t, cmd.Args)
	}
}

func TestHandler_CreateCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.CreateCommands()

	// Should have create commands from all registered handlers that implement Writer
	assert.NotEmpty(t, commands)

	// Verify all commands are non-nil
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_DeleteCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.DeleteCommands()

	// Should have delete commands from all registered handlers that implement Writer
	assert.NotEmpty(t, commands)

	// Verify all commands are non-nil
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_CopyCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.CopyCommands()

	// Should have copy commands from registered handlers that implement Copier
	// Not all handlers implement Copier, so this may be empty or have fewer commands
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_ClearCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.ClearCommands()

	// Should have clear commands from registered handlers that implement Writer
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_ImportCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.ImportCommands()

	// Should have import commands from registered handlers that implement Importer
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_ExportCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.ExportCommands()

	// Should have export commands from registered handlers that implement Exporter
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_StartCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.StartCommands()

	// Should have start commands from registered handlers that implement Controller
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_StopCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.StopCommands()

	// Should have stop commands from registered handlers that implement Controller
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_RestartCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.RestartCommands()

	// Should have restart commands from registered handlers that implement Controller
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_InspectCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.InspectCommands()

	// Should have inspect commands from registered handlers that implement Inspector
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_EditCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.EditCommands()

	// Should have edit commands from registered handlers that implement Editor
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_DumpCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.DumpCommands()

	// Should have dump commands from registered handlers that implement Dumper
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_LoadCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	commands := handler.LoadCommands()

	// Should have load commands from registered handlers that implement Loader
	for _, cmd := range commands {
		assert.NotNil(t, cmd)
	}
}

func TestHandler_AddCommandGroup(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	rootCmd := &cobra.Command{
		Use: "root",
	}

	// Add a command group
	handler.AddCommandGroup(rootCmd, "Test Group", func(h Handler, groupID string) []*cobra.Command {
		return []*cobra.Command{
			{Use: "cmd1"},
			{Use: "cmd2"},
		}
	})

	// Verify group was added
	groups := rootCmd.Groups()
	require.Len(t, groups, 1)
	assert.Equal(t, "test-group", groups[0].ID)
	assert.Equal(t, "Test Group", groups[0].Title)

	// Verify commands were added
	commands := rootCmd.Commands()
	assert.Len(t, commands, 2)
}

func TestHandler_AddCommandGroup_IDFormatting(t *testing.T) {
	tests := []struct {
		title      string
		expectedID string
	}{
		{
			title:      "Simple Title",
			expectedID: "simple-title",
		},
		{
			title:      "Title With Multiple Words",
			expectedID: "title-with-multiple-words",
		},
		{
			title:      "UPPERCASE TITLE",
			expectedID: "uppercase-title",
		},
		{
			title:      "Mixed Case Title",
			expectedID: "mixed-case-title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
			require.NoError(t, err)

			handler := NewHandler(rt)

			rootCmd := &cobra.Command{
				Use: "root",
			}

			handler.AddCommandGroup(rootCmd, tt.title, func(h Handler, groupID string) []*cobra.Command {
				// Verify the groupID passed to the function matches expected ID
				assert.Equal(t, tt.expectedID, groupID)
				return []*cobra.Command{{Use: "test"}}
			})

			groups := rootCmd.Groups()
			require.Len(t, groups, 1)
			assert.Equal(t, tt.expectedID, groups[0].ID)
		})
	}
}

func TestHandler_AllCommandMethods_ReturnCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	// Test all command getter methods return non-nil slices
	assert.NotNil(t, handler.GetCommands())
	assert.NotNil(t, handler.DescribeCommands())
	assert.NotNil(t, handler.CreateCommands())
	assert.NotNil(t, handler.DeleteCommands())
	assert.NotNil(t, handler.CopyCommands())
	assert.NotNil(t, handler.ClearCommands())
	assert.NotNil(t, handler.ImportCommands())
	assert.NotNil(t, handler.ExportCommands())
	assert.NotNil(t, handler.StartCommands())
	assert.NotNil(t, handler.StopCommands())
	assert.NotNil(t, handler.RestartCommands())
	assert.NotNil(t, handler.InspectCommands())
	assert.NotNil(t, handler.EditCommands())
	assert.NotNil(t, handler.DumpCommands())
	assert.NotNil(t, handler.LoadCommands())
}

func TestHandler_CommandsHaveRunE(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	// Test that all returned commands have RunE set
	commands := handler.GetCommands()
	for _, cmd := range commands {
		assert.NotNil(t, cmd.RunE, "command %s should have RunE set", cmd.Use)
	}
}

func TestHandler_DescriptorsSharedWithRuntime(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	// Descriptors should be the same as runtime's descriptors
	assert.Equal(t, rt.GetDescriptors(), handler.descriptors)
}

func TestHandler_MultipleHandlers_IndependentCommands(t *testing.T) {
	rt1, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	rt2, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler1 := NewHandler(rt1)
	handler2 := NewHandler(rt2)

	// Both handlers should have the same number of commands
	commands1 := handler1.GetCommands()
	commands2 := handler2.GetCommands()

	assert.Equal(t, len(commands1), len(commands2))
}

func TestToGroupID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Asset Commands",
			expected: "asset-commands",
		},
		{
			input:    "Platform Commands",
			expected: "platform-commands",
		},
		{
			input:    "UPPERCASE",
			expected: "uppercase",
		},
		{
			input:    "multiple word title here",
			expected: "multiple-word-title-here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Test the ID transformation by checking the result in AddCommandGroup
			rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
			require.NoError(t, err)

			handler := NewHandler(rt)
			rootCmd := &cobra.Command{Use: "root"}

			handler.AddCommandGroup(rootCmd, tt.input, func(h Handler, groupID string) []*cobra.Command {
				assert.Equal(t, tt.expected, groupID)
				return nil
			})
		})
	}
}

func TestHandler_CommandGroup_EmptyCommands(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)
	rootCmd := &cobra.Command{Use: "root"}

	// Add command group that returns empty command list
	handler.AddCommandGroup(rootCmd, "Test Group", func(h Handler, groupID string) []*cobra.Command {
		return []*cobra.Command{}
	})

	// Should handle empty list gracefully
	groups := rootCmd.Groups()
	assert.Len(t, groups, 1)
	assert.Equal(t, "test-group", groups[0].ID)
}

func TestHandler_RegistryPopulated(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	// Verify registry is populated with handlers
	assert.NotEmpty(t, handler.registry.Readers())
	assert.NotEmpty(t, handler.registry.Writers())

	// Check specific handler types are registered
	readers := handler.registry.Readers()
	writers := handler.registry.Writers()

	// Should have multiple handlers registered
	assert.Greater(t, len(readers), 10, "should have at least 10 reader handlers")
	assert.Greater(t, len(writers), 10, "should have at least 10 writer handlers")
}

func TestHandler_AllHandlersRegistered(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	handler := NewHandler(rt)

	// Count total handlers by checking command availability
	resourceTypes := []string{
		"projects", "workflows", "transformations", "jsonforms",
		"command_templates", "analytic_templates", "templates",
		"automations", "accounts", "profiles", "roles", "groups",
		"methods", "views", "prebuilts", "integration_models",
		"integrations", "adapters", "adapter_models", "tags",
		"applications", "devices", "devicegroups", "configuration_parsers",
		"gctrees", "models", "server",
	}

	// Get all get commands
	commands := handler.GetCommands()
	commandMap := make(map[string]*cobra.Command)
	for _, cmd := range commands {
		// Extract base resource name from command use
		parts := strings.Split(cmd.Use, " ")
		if len(parts) > 0 {
			commandMap[parts[0]] = cmd
		}
	}

	// Verify we have commands for core resources
	foundResources := 0
	for _, resource := range resourceTypes {
		if _, found := commandMap[resource]; found {
			foundResources++
		}
	}

	// Should have found a significant number of expected resources
	assert.Greater(t, foundResources, 15, "should have registered handlers for most resources")
}
