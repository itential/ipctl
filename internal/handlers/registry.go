// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

var (
	registry  []any
	readers   []Reader
	writers   []Writer
	copiers   []Copier
	editors   []Editor
	importers []Importer
	exporters []Exporter
	controllers []Controller
	inspectors []Inspector
	dumpers   []Dumper
	loaders   []Loader
)

func register(types ...any) {
	// FIXME (privateip) this is a short term work around for the fact that the
	// handler is now loaded multiple times.
	if len(registry) > 0 {
		return
	}

	registry = types

	// Pre-compute type-specific slices once at initialization
	for _, ele := range types {
		if r, ok := ele.(Reader); ok {
			readers = append(readers, r)
		}
		if w, ok := ele.(Writer); ok {
			writers = append(writers, w)
		}
		if c, ok := ele.(Copier); ok {
			copiers = append(copiers, c)
		}
		if e, ok := ele.(Editor); ok {
			editors = append(editors, e)
		}
		if i, ok := ele.(Importer); ok {
			importers = append(importers, i)
		}
		if e, ok := ele.(Exporter); ok {
			exporters = append(exporters, e)
		}
		if c, ok := ele.(Controller); ok {
			controllers = append(controllers, c)
		}
		if i, ok := ele.(Inspector); ok {
			inspectors = append(inspectors, i)
		}
		if d, ok := ele.(Dumper); ok {
			dumpers = append(dumpers, d)
		}
		if l, ok := ele.(Loader); ok {
			loaders = append(loaders, l)
		}
	}
}

func Readers() []Reader {
	return readers
}

func Writers() []Writer {
	return writers
}

func Copiers() []Copier {
	return copiers
}

func Editors() []Editor {
	return editors
}

func Importers() []Importer {
	return importers
}

func Exporters() []Exporter {
	return exporters
}

func Controllers() []Controller {
	return controllers
}

func Inspectors() []Inspector {
	return inspectors
}

func Dumpers() []Dumper {
	return dumpers
}

func Loaders() []Loader {
	return loaders
}
