// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/runners"
	"github.com/spf13/cobra"
)

// AssetHandlerFlags holds optional flag definitions for asset handler commands.
type AssetHandlerFlags struct {
	Create   flags.Flagger
	Delete   flags.Flagger
	Get      flags.Flagger
	Describe flags.Flagger
	Copy     flags.Flagger
	Clear    flags.Flagger

	Edit flags.Flagger

	Import flags.Flagger
	Export flags.Flagger

	Start   flags.Flagger
	Stop    flags.Flagger
	Restart flags.Flagger

	Inspect flags.Flagger

	Dump flags.Flagger
	Load flags.Flagger
}

// AssetHandler provides command generation for asset-based resources.
// It uses type assertions instead of reflection to determine which
// operations a runner supports.
type AssetHandler struct {
	// Typed runner fields for each operation
	reader     runners.Reader
	writer     runners.Writer
	copier     runners.Copier
	editor     runners.Editor
	importer   runners.Importer
	exporter   runners.Exporter
	controller runners.Controller
	inspector  runners.Inspector
	dumper     runners.Dumper
	loader     runners.Loader

	descriptor DescriptorMap
	flags      *AssetHandlerFlags
}

// NewAssetHandler creates a new AssetHandler from a runner.
// It uses type assertions at construction time to determine which
// operations the runner supports, avoiding runtime reflection.
func NewAssetHandler(runner runners.Runner, dm DescriptorMap, flagsArg *AssetHandlerFlags) AssetHandler {
	if flagsArg == nil {
		flagsArg = &AssetHandlerFlags{}
	}

	handler := AssetHandler{
		descriptor: dm,
		flags:      flagsArg,
	}

	// Use type assertions instead of reflection
	if reader, ok := runner.(runners.Reader); ok {
		handler.reader = reader
	}
	if writer, ok := runner.(runners.Writer); ok {
		handler.writer = writer
	}
	if copier, ok := runner.(runners.Copier); ok {
		handler.copier = copier
	}
	if editor, ok := runner.(runners.Editor); ok {
		handler.editor = editor
	}
	if importer, ok := runner.(runners.Importer); ok {
		handler.importer = importer
	}
	if exporter, ok := runner.(runners.Exporter); ok {
		handler.exporter = exporter
	}
	if controller, ok := runner.(runners.Controller); ok {
		handler.controller = controller
	}
	if inspector, ok := runner.(runners.Inspector); ok {
		handler.inspector = inspector
	}
	if dumper, ok := runner.(runners.Dumper); ok {
		handler.dumper = dumper
	}
	if loader, ok := runner.(runners.Loader); ok {
		handler.loader = loader
	}

	return handler
}

func (h AssetHandler) newCommand(key string, runtime *Runtime, runner runners.RunnerFunc, common flags.Flagger) *cobra.Command {
	cmd := NewCommandRunner(
		key,
		h.descriptor,
		runner,
		runtime,
		common,
		withOptions(h.flags),
	)
	return NewCommand(cmd)
}

// Get returns the 'get' command if the runner supports the Reader interface.
func (h AssetHandler) Get(runtime *Runtime) *cobra.Command {
	logging.Trace()
	if h.reader == nil {
		return nil
	}
	cmd := h.newCommand("get", runtime, h.reader.Get, nil)
	if cmd != nil && h.flags.Get != nil {
		h.flags.Get.Flags(cmd)
	}
	return cmd
}

// Describe returns the 'describe' command if the runner supports the Reader interface.
func (h AssetHandler) Describe(runtime *Runtime) *cobra.Command {
	if h.reader == nil {
		return nil
	}
	cmd := h.newCommand("describe", runtime, h.reader.Describe, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Describe != nil {
			h.flags.Describe.Flags(cmd)
		}
	}
	return cmd
}

// Create returns the 'create' command if the runner supports the Writer interface.
func (h AssetHandler) Create(runtime *Runtime) *cobra.Command {
	if h.writer == nil {
		return nil
	}
	cmd := h.newCommand("create", runtime, h.writer.Create, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Create != nil {
			h.flags.Create.Flags(cmd)
		}
	}
	return cmd
}

// Delete returns the 'delete' command if the runner supports the Writer interface.
func (h AssetHandler) Delete(runtime *Runtime) *cobra.Command {
	if h.writer == nil {
		return nil
	}
	cmd := h.newCommand("delete", runtime, h.writer.Delete, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Delete != nil {
			h.flags.Delete.Flags(cmd)
		}
	}
	return cmd
}

// Copy returns the 'copy' command if the runner supports the Copier interface.
func (h AssetHandler) Copy(runtime *Runtime) *cobra.Command {
	if h.copier == nil {
		return nil
	}
	common := &flags.AssetCopyCommon{}
	cmd := h.newCommand("copy", runtime, h.copier.Copy, common)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		common.Flags(cmd)
		if h.flags.Copy != nil {
			h.flags.Copy.Flags(cmd)
		}
	}
	return cmd
}

// Clear returns the 'clear' command if the runner supports the Writer interface.
func (h AssetHandler) Clear(runtime *Runtime) *cobra.Command {
	if h.writer == nil {
		return nil
	}
	cmd := h.newCommand("clear", runtime, h.writer.Clear, nil)
	if cmd != nil && h.flags.Clear != nil {
		h.flags.Clear.Flags(cmd)
	}
	return cmd
}

// Import returns the 'import' command if the runner supports the Importer interface.
func (h AssetHandler) Import(runtime *Runtime) *cobra.Command {
	if h.importer == nil {
		return nil
	}
	common := &flags.AssetImportCommon{}
	cmd := h.newCommand("import", runtime, h.importer.Import, common)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		common.Flags(cmd)
		if h.flags.Import != nil {
			h.flags.Import.Flags(cmd)
		}
	}
	return cmd
}

// Export returns the 'export' command if the runner supports the Exporter interface.
func (h AssetHandler) Export(runtime *Runtime) *cobra.Command {
	if h.exporter == nil {
		return nil
	}
	common := &flags.AssetExportCommon{}
	cmd := h.newCommand("export", runtime, h.exporter.Export, common)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		common.Flags(cmd)
		if h.flags.Export != nil {
			h.flags.Export.Flags(cmd)
		}
	}
	return cmd
}

// Start returns the 'start' command if the runner supports the Controller interface.
func (h AssetHandler) Start(runtime *Runtime) *cobra.Command {
	if h.controller == nil {
		return nil
	}
	cmd := h.newCommand("start", runtime, h.controller.Start, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Start != nil {
			h.flags.Start.Flags(cmd)
		}
	}
	return cmd
}

// Stop returns the 'stop' command if the runner supports the Controller interface.
func (h AssetHandler) Stop(runtime *Runtime) *cobra.Command {
	if h.controller == nil {
		return nil
	}
	cmd := h.newCommand("stop", runtime, h.controller.Stop, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Stop != nil {
			h.flags.Stop.Flags(cmd)
		}
	}
	return cmd
}

// Restart returns the 'restart' command if the runner supports the Controller interface.
func (h AssetHandler) Restart(runtime *Runtime) *cobra.Command {
	if h.controller == nil {
		return nil
	}
	cmd := h.newCommand("restart", runtime, h.controller.Restart, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Restart != nil {
			h.flags.Restart.Flags(cmd)
		}
	}
	return cmd
}

// Inspect returns the 'inspect' command if the runner supports the Inspector interface.
func (h AssetHandler) Inspect(runtime *Runtime) *cobra.Command {
	if h.inspector == nil {
		return nil
	}
	cmd := h.newCommand("inspect", runtime, h.inspector.Inspect, nil)
	if cmd != nil && h.flags.Inspect != nil {
		h.flags.Inspect.Flags(cmd)
	}
	return cmd
}

// Edit returns the 'edit' command if the runner supports the Editor interface.
func (h AssetHandler) Edit(runtime *Runtime) *cobra.Command {
	if h.editor == nil {
		return nil
	}
	cmd := h.newCommand("edit", runtime, h.editor.Edit, nil)
	if cmd != nil {
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Edit != nil {
			h.flags.Edit.Flags(cmd)
		}
	}
	return cmd
}

// Dump returns the 'dump' command if the runner supports the Dumper interface.
func (h AssetHandler) Dump(runtime *Runtime) *cobra.Command {
	if h.dumper == nil {
		return nil
	}
	common := &flags.AssetDumpCommon{}
	cmd := h.newCommand("dump", runtime, h.dumper.Dump, common)
	if cmd != nil {
		common.Flags(cmd)
		if h.flags.Dump != nil {
			h.flags.Dump.Flags(cmd)
		}
	}
	return cmd
}

// Load returns the 'load' command if the runner supports the Loader interface.
func (h AssetHandler) Load(runtime *Runtime) *cobra.Command {
	if h.loader == nil {
		return nil
	}
	common := &flags.AssetLoadCommon{}
	cmd := h.newCommand("load", runtime, h.loader.Load, common)
	if cmd != nil {
		common.Flags(cmd)
		cmd.Args = cobra.ExactArgs(1)
		if h.flags.Load != nil {
			h.flags.Load.Flags(cmd)
		}
	}
	return cmd
}
