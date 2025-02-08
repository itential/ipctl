// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"reflect"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/spf13/cobra"
)

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

	Pull flags.Flagger
	Push flags.Flagger
}

type AssetHandler struct {
	Runner       runners.Runner
	isReader     bool
	isWriter     bool
	isCopier     bool
	isEditor     bool
	isImporter   bool
	isExporter   bool
	isController bool
	isInspector  bool
	isGitter     bool
	Descriptor   DescriptorMap
	Flags        *AssetHandlerFlags
}

func NewAssetHandler(r runners.Runner, dm DescriptorMap, flags *AssetHandlerFlags) AssetHandler {
	if flags == nil {
		flags = &AssetHandlerFlags{}
	}

	assetHandler := AssetHandler{
		Descriptor: dm,
		Flags:      flags,
		Runner:     r,
	}

	assetHandler.isReader = implements(r, (*runners.Reader)(nil))
	assetHandler.isWriter = implements(r, (*runners.Writer)(nil))
	assetHandler.isCopier = implements(r, (*runners.Copier)(nil))
	assetHandler.isEditor = implements(r, (*runners.Editor)(nil))
	assetHandler.isImporter = implements(r, (*runners.Importer)(nil))
	assetHandler.isExporter = implements(r, (*runners.Exporter)(nil))
	assetHandler.isController = implements(r, (*runners.Controller)(nil))
	assetHandler.isInspector = implements(r, (*runners.Inspector)(nil))
	assetHandler.isGitter = implements(r, (*runners.Gitter)(nil))

	return assetHandler
}

func implements(r any, i any) bool {
	elem := reflect.TypeOf(i).Elem()
	return reflect.TypeOf(r).Implements(elem)
}

func (h AssetHandler) newCommand(key string, runtime *Runtime, runner runners.RunnerFunc, common flags.Flagger) *cobra.Command {
	cmd := NewCommandRunner(
		key,
		h.Descriptor,
		runner,
		runtime,
		common,
		withOptions(h.Flags),
	)
	return NewCommand(cmd)
}

func (h AssetHandler) Get(runtime *Runtime) *cobra.Command {
	logger.Trace()
	var cmd *cobra.Command
	if h.isReader {
		cmd = h.newCommand("get", runtime, h.Runner.(runners.Reader).Get, nil)
		if cmd != nil {
			if h.Flags.Get != nil {
				h.Flags.Get.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Describe(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isReader {
		cmd = h.newCommand("describe", runtime, h.Runner.(runners.Reader).Describe, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Describe != nil {
				h.Flags.Describe.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Create(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isWriter {
		cmd = h.newCommand("create", runtime, h.Runner.(runners.Writer).Create, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Create != nil {
				h.Flags.Create.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Delete(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isWriter {
		cmd = h.newCommand("delete", runtime, h.Runner.(runners.Writer).Delete, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Delete != nil {
				h.Flags.Delete.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Copy(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isCopier {
		common := &flags.AssetCopyCommon{}
		cmd = h.newCommand("copy", runtime, h.Runner.(runners.Copier).Copy, common)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			common.Flags(cmd)
			if h.Flags.Copy != nil {
				h.Flags.Copy.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Clear(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isWriter {
		cmd = h.newCommand("clear", runtime, h.Runner.(runners.Writer).Clear, nil)
		if cmd != nil {
			if h.Flags.Clear != nil {
				h.Flags.Clear.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Import(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isImporter {
		common := &flags.AssetImportCommon{}
		cmd = h.newCommand("import", runtime, h.Runner.(runners.Importer).Import, common)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			common.Flags(cmd)
			if h.Flags.Import != nil {
				h.Flags.Import.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Export(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isExporter {
		common := &flags.AssetExportCommon{}
		cmd = h.newCommand("export", runtime, h.Runner.(runners.Exporter).Export, common)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			common.Flags(cmd)
			if h.Flags.Export != nil {
				h.Flags.Export.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Start(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isController {
		cmd = h.newCommand("start", runtime, h.Runner.(runners.Controller).Start, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Start != nil {
				h.Flags.Start.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Stop(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isController {
		cmd = h.newCommand("stop", runtime, h.Runner.(runners.Controller).Stop, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Stop != nil {
				h.Flags.Stop.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Restart(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isController {
		cmd = h.newCommand("restart", runtime, h.Runner.(runners.Controller).Restart, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Restart != nil {
				h.Flags.Restart.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Inspect(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isInspector {
		cmd = h.newCommand("inspect", runtime, h.Runner.(runners.Inspector).Inspect, nil)
		if cmd != nil {
			if h.Flags.Inspect != nil {
				h.Flags.Inspect.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Edit(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isEditor {
		cmd = h.newCommand("edit", runtime, h.Runner.(runners.Editor).Edit, nil)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(1)
			if h.Flags.Edit != nil {
				h.Flags.Edit.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Push(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isGitter {
		common := &flags.AssetPushCommon{}
		cmd = h.newCommand("push", runtime, h.Runner.(runners.Gitter).Push, common)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(2)
			common.Flags(cmd)
			if h.Flags.Push != nil {
				h.Flags.Push.Flags(cmd)
			}
		}
	}
	return cmd
}

func (h AssetHandler) Pull(runtime *Runtime) *cobra.Command {
	var cmd *cobra.Command
	if h.isGitter {
		common := &flags.AssetPullCommon{}
		cmd = h.newCommand("pull", runtime, h.Runner.(runners.Gitter).Pull, common)
		if cmd != nil {
			cmd.Args = cobra.ExactArgs(2)
			common.Flags(cmd)
			if h.Flags.Pull != nil {
				h.Flags.Pull.Flags(cmd)
			}
		}
	}
	return cmd
}
