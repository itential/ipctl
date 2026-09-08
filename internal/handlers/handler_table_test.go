// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"testing"

	"github.com/itential/ipctl/internal/terminal"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestHandler creates a Handler backed by mock dependencies for use in tests.
func newTestHandler(t *testing.T) Handler {
	t.Helper()
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)
	return NewHandler(rt)
}

// countTypes returns the number of specs in handlerTable that satisfy pred.
func countTypes(pred func(handlerTypes) bool) int {
	n := 0
	for _, spec := range handlerTable {
		if pred(spec.types) {
			n++
		}
	}
	return n
}

// handlerTypes declares which handler archetypes a registered handler should implement.
// Set a field to true to assert that the corresponding command is generated; false
// asserts that it is not. Adding a new handler requires only a new entry in handlerTable.
type handlerTypes struct {
	Reader     bool
	Writer     bool
	Copier     bool
	Editor     bool
	Importer   bool
	Exporter   bool
	Controller bool
	Inspector  bool
	Dumper     bool
	Loader     bool
}

// handlerSpec pairs a human-readable name with the archetypes its handler should implement.
// The name is used only for t.Run labelling — it is not matched against command Use strings.
type handlerSpec struct {
	name  string
	types handlerTypes
}

// handlerTable is the ground truth for every handler registered in NewHandler.
// To add a new handler: register it in NewHandler() and add a row here.
// commandSpec pairs a command name with its getter and whether it requires positional args.
type commandSpec struct {
	name        string
	getter      func(Handler) []*cobra.Command
	requiresArg bool
}

// interfaceTable defines which commands each handler archetype exposes and which take positional args.
// hasArgs=true asserts cmd.Args is set (cobra enforces the arg count); false skips that check.
var interfaceTable = []struct {
	typeName string
	pred     func(handlerTypes) bool
	commands []commandSpec
}{
	{
		typeName: "Reader",
		pred:     func(ht handlerTypes) bool { return ht.Reader },
		commands: []commandSpec{
			{"get", Handler.GetCommands, false},
			{"describe", Handler.DescribeCommands, true},
		},
	},
	{
		typeName: "Writer",
		pred:     func(ht handlerTypes) bool { return ht.Writer },
		commands: []commandSpec{
			{"create", Handler.CreateCommands, false},
			{"delete", Handler.DeleteCommands, true},
			{"clear", Handler.ClearCommands, false},
		},
	},
	{
		typeName: "Copier",
		pred:     func(ht handlerTypes) bool { return ht.Copier },
		commands: []commandSpec{
			{"copy", Handler.CopyCommands, true},
		},
	},
	{
		typeName: "Editor",
		pred:     func(ht handlerTypes) bool { return ht.Editor },
		commands: []commandSpec{
			{"edit", Handler.EditCommands, true},
		},
	},
	{
		typeName: "Importer",
		pred:     func(ht handlerTypes) bool { return ht.Importer },
		commands: []commandSpec{
			{"import", Handler.ImportCommands, true},
		},
	},
	{
		typeName: "Exporter",
		pred:     func(ht handlerTypes) bool { return ht.Exporter },
		commands: []commandSpec{
			{"export", Handler.ExportCommands, true},
		},
	},
	{
		typeName: "Controller",
		pred:     func(ht handlerTypes) bool { return ht.Controller },
		commands: []commandSpec{
			{"start", Handler.StartCommands, true},
			{"stop", Handler.StopCommands, true},
			{"restart", Handler.RestartCommands, true},
		},
	},
	{
		typeName: "Inspector",
		pred:     func(ht handlerTypes) bool { return ht.Inspector },
		commands: []commandSpec{
			{"inspect", Handler.InspectCommands, false},
		},
	},
	{
		typeName: "Dumper",
		pred:     func(ht handlerTypes) bool { return ht.Dumper },
		commands: []commandSpec{
			{"dump", Handler.DumpCommands, false},
		},
	},
	{
		typeName: "Loader",
		pred:     func(ht handlerTypes) bool { return ht.Loader },
		commands: []commandSpec{
			{"load", Handler.LoadCommands, true},
		},
	},
}

func TestHandlers_Interfaces(t *testing.T) {
	h := newTestHandler(t)

	for _, iface := range interfaceTable {
		expected := countTypes(iface.pred)
		t.Run(iface.typeName, func(t *testing.T) {
			for _, spec := range iface.commands {
				cmds := spec.getter(h)
				t.Run(spec.name, func(t *testing.T) {
					t.Run("Count", func(t *testing.T) {
						assert.Len(t, cmds, expected)
					})
					t.Run("HaveRunE", func(t *testing.T) {
						for _, cmd := range cmds {
							t.Run(cmd.Use, func(t *testing.T) {
								assert.NotNil(t, cmd.RunE)
							})
						}
					})
					if spec.requiresArg {
						t.Run("HasArg", func(t *testing.T) {
							for _, cmd := range cmds {
								t.Run(cmd.Use, func(t *testing.T) {
									assert.NotNil(t, cmd.Args)
								})
							}
						})
					}
				})
			}
		})
	}
}

var handlerTable = []handlerSpec{
	// Automation Studio
	{
		name: "projects",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "workflows",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Editor:   true,
			Importer: true,
			Exporter: true,
			Dumper:   true,
			Loader:   true,
		},
	},
	{
		name: "transformations",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "jsonforms",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "command-templates",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "analytic-templates",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "templates",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
			Dumper:   true,
			Loader:   true,
		},
	},

	// Operations Manager
	{
		name: "automations",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
			Dumper:   true,
			Loader:   true,
		},
	},

	// Admin Essentials
	{
		name: "accounts",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "profiles",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
			Dumper:   true,
			Loader:   true,
		},
	},
	{
		name: "roles",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "role-types",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "groups",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "methods",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "views",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "prebuilts",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "integration-models",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "integrations",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},
	{
		name: "adapters",
		types: handlerTypes{
			Reader:     true,
			Writer:     true,
			Copier:     true,
			Editor:     true,
			Importer:   true,
			Exporter:   true,
			Controller: true,
			Inspector:  true,
			Dumper:     true,
			Loader:     true,
		},
	},
	{
		name: "adapter-models",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "tags",
		types: handlerTypes{
			Reader: true,
			Writer: true,
			Copier: true,
		},
	},
	{
		name: "applications",
		types: handlerTypes{
			Reader:     true,
			Controller: true,
			Inspector:  true,
		},
	},

	// Configuration Manager
	{
		name: "devices",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "device-groups",
		types: handlerTypes{
			Reader: true,
			Writer: true,
		},
	},
	{
		name: "configuration-parsers",
		types: handlerTypes{
			Reader: true,
		},
	},
	{
		name: "gctrees",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Importer: true,
			Exporter: true,
		},
	},

	// Lifecycle Manager
	{
		name: "models",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},

	// Flow Agent
	{
		name: "agent-projects",
		types: handlerTypes{
			Reader:   true,
			Writer:   true,
			Copier:   true,
			Importer: true,
			Exporter: true,
		},
	},

	// Platform
	{
		name: "server",
		types: handlerTypes{
			Inspector: true,
		},
	},
}
