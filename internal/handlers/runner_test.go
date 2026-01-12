// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"errors"
	"testing"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRunner implements runners.Runner for testing
type mockRunner struct{}

func TestNewCommandRunner(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:         "resources",
			Description: "Get resources",
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	options := &mockFlagger{}

	cr := NewCommandRunner("get", desc, runFunc, rt, options)

	require.NotNil(t, cr)
	assert.Equal(t, "get", cr.Key)
	assert.Equal(t, desc, cr.Descriptors)
	assert.NotNil(t, cr.Run)
	assert.Equal(t, rt, cr.Runtime)
	assert.Equal(t, options, cr.Common)
}

func TestNewCommandRunner_WithOptions(t *testing.T) {
	desc := DescriptorMap{
		"create": cmdutils.Descriptor{Use: "resource"},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "created"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	assetFlags := &AssetHandlerFlags{
		Create: &mockFlagger{},
		Get:    &mockFlagger{},
		Copy:   &mockFlagger{},
	}

	cr := NewCommandRunner(
		"create",
		desc,
		runFunc,
		rt,
		&mockFlagger{},
		withOptions(assetFlags),
	)

	require.NotNil(t, cr)
	assert.Equal(t, assetFlags.Create, cr.Options)
}

func TestWithOptions_AllKeys(t *testing.T) {
	tests := []struct {
		key      string
		expected flags.Flagger
	}{
		{"create", &mockFlagger{}},
		{"delete", &mockFlagger{}},
		{"get", &mockFlagger{}},
		{"describe", &mockFlagger{}},
		{"copy", &mockFlagger{}},
		{"clear", &mockFlagger{}},
		{"import", &mockFlagger{}},
		{"export", &mockFlagger{}},
		{"load", &mockFlagger{}},
		{"dump", &mockFlagger{}},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			assetFlags := &AssetHandlerFlags{
				Create:   &mockFlagger{},
				Delete:   &mockFlagger{},
				Get:      &mockFlagger{},
				Describe: &mockFlagger{},
				Copy:     &mockFlagger{},
				Clear:    &mockFlagger{},
				Import:   &mockFlagger{},
				Export:   &mockFlagger{},
				Load:     &mockFlagger{},
				Dump:     &mockFlagger{},
			}

			cr := &CommandRunner{Key: tt.key}
			opt := withOptions(assetFlags)
			opt(cr)

			assert.NotNil(t, cr.Options)
		})
	}
}

func TestWithOptions_UnknownKey(t *testing.T) {
	assetFlags := &AssetHandlerFlags{
		Create: &mockFlagger{},
	}

	cr := &CommandRunner{Key: "unknown"}
	opt := withOptions(assetFlags)
	opt(cr)

	// Options should remain nil for unknown keys
	assert.Nil(t, cr.Options)
}

func TestNewCommand_Success(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:         "resources",
			Description: "Get all resources",
			Example:     "ipctl get resources\nipctl get resources --format json",
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{
			Text: "success",
		}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{
		DefaultOutput: "human",
		Pager:         false,
	})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
		Options:     &mockFlagger{},
		Runner:      &mockRunner{},
	}

	cmd := NewCommand(cr)

	require.NotNil(t, cmd)
	assert.Equal(t, "resources", cmd.Use)
	assert.Equal(t, "Get all resources", cmd.Short)
	assert.Equal(t, "Get all resources", cmd.Long)
	assert.Contains(t, cmd.Example, "ipctl get resources")
	assert.NotNil(t, cmd.RunE)
}

func TestNewCommand_DisabledDescriptor(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:      "resources",
			Disabled: true,
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
	}

	cmd := NewCommand(cr)

	// Should return nil for disabled descriptor
	assert.Nil(t, cmd)
}

func TestNewCommand_MissingDescriptor(t *testing.T) {
	desc := DescriptorMap{}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
	}

	cmd := NewCommand(cr)

	// Should return nil for missing descriptor
	assert.Nil(t, cmd)
}

func TestNewCommand_WithExactArgs(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:       "resources",
			ExactArgs: 2,
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
	}

	cmd := NewCommand(cr)

	require.NotNil(t, cmd)
	assert.NotNil(t, cmd.Args)

	// Test that Args validator is set correctly
	err = cmd.Args(cmd, []string{"arg1", "arg2"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"arg1"})
	assert.Error(t, err)
}

func TestNewCommand_HiddenFlag(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:    "resources",
			Hidden: true,
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
	}

	cmd := NewCommand(cr)

	require.NotNil(t, cmd)
	assert.True(t, cmd.Hidden)
}

func TestNewCommand_WithGroupID(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:   "resources",
			Group: "resource-group",
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
	}

	cmd := NewCommand(cr)

	require.NotNil(t, cmd)
	assert.Equal(t, "resource-group", cmd.GroupID)
}

func TestNewCommand_RunE_Success(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use: "resources",
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{
			Text: "operation successful",
		}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{
		DefaultOutput: "human",
		Pager:         false,
	})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
		Options:     &mockFlagger{},
		Runner:      &mockRunner{},
	}

	cmd := NewCommand(cr)
	require.NotNil(t, cmd)

	// Execute the command
	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestNewCommand_RunE_Error(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use: "resources",
		},
	}

	expectedErr := errors.New("runner error")
	runFunc := func(req runners.Request) (*runners.Response, error) {
		return nil, expectedErr
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
		Options:     &mockFlagger{},
		Runner:      &mockRunner{},
	}

	cmd := NewCommand(cr)
	require.NotNil(t, cmd)

	// Execute the command - should return error from runner
	err = cmd.RunE(cmd, []string{})
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestNewCommand_ExampleFormatting(t *testing.T) {
	desc := DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:     "resources",
			Example: "line1\nline2\nline3",
		},
	}

	runFunc := func(req runners.Request) (*runners.Response, error) {
		return &runners.Response{Text: "success"}, nil
	}

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cr := &CommandRunner{
		Key:         "get",
		Descriptors: desc,
		Run:         runFunc,
		Runtime:     rt,
		Common:      &mockFlagger{},
	}

	cmd := NewCommand(cr)

	require.NotNil(t, cmd)
	// Example should be indented with 2 spaces
	assert.Contains(t, cmd.Example, "  line1")
	assert.Contains(t, cmd.Example, "  line2")
	assert.Contains(t, cmd.Example, "  line3")
}

func TestCommandRunner_FieldsRemoved(t *testing.T) {
	// Verify that unused fields were removed from CommandRunner
	cr := &CommandRunner{}

	// These fields should not exist (will fail to compile if they do)
	// cr.Confirm = true     // Should not compile
	// cr.PreRun = func(args []string) error { return nil }  // Should not compile
	// cr.PostRun = func(args []string) {}  // Should not compile

	// Only valid fields should exist
	assert.NotNil(t, &cr.Key)
	assert.NotNil(t, &cr.Descriptors)
	assert.NotNil(t, &cr.Run)
	assert.NotNil(t, &cr.Common)
	assert.NotNil(t, &cr.Options)
	assert.NotNil(t, &cr.Runtime)
	assert.NotNil(t, &cr.Runner)
	assert.NotNil(t, &cr.Flags)
}
