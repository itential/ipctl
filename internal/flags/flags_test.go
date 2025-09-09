// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type mockFlagger struct {
	called bool
}

func (m *mockFlagger) Flags(cmd *cobra.Command) {
	m.called = true
	cmd.Flags().String("test", "default", "test flag")
}

func TestFlaggerInterface(t *testing.T) {
	mock := &mockFlagger{}
	cmd := &cobra.Command{}
	
	mock.Flags(cmd)
	
	assert.True(t, mock.called, "Flags method should be called")
	assert.True(t, cmd.Flags().HasFlags(), "Command should have flags after calling Flags method")
	
	flag := cmd.Flag("test")
	assert.NotNil(t, flag, "Test flag should exist")
	assert.Equal(t, "test", flag.Name)
	assert.Equal(t, "test flag", flag.Usage)
}

func TestOptionStruct(t *testing.T) {
	option := Option{
		Name:   "test-option",
		Abbrev: "t",
		Usage:  "Test option for testing",
	}
	
	assert.Equal(t, "test-option", option.Name)
	assert.Equal(t, "t", option.Abbrev)
	assert.Equal(t, "Test option for testing", option.Usage)
}

func TestOptionStructZeroValues(t *testing.T) {
	option := Option{}
	
	assert.Empty(t, option.Name)
	assert.Empty(t, option.Abbrev)
	assert.Empty(t, option.Usage)
}