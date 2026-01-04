// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

// Registry holds categorized handlers for different command types.
// This is an instance-based registry that avoids global mutable state,
// making the code thread-safe and testable.
type Registry struct {
	readers     []Reader
	writers     []Writer
	copiers     []Copier
	editors     []Editor
	importers   []Importer
	exporters   []Exporter
	controllers []Controller
	inspectors  []Inspector
	dumpers     []Dumper
	loaders     []Loader
}

// NewRegistry creates and populates a new handler registry.
// It categorizes handlers by their supported interfaces using type assertions.
func NewRegistry(handlers []any) *Registry {
	r := &Registry{}

	for _, handler := range handlers {
		if reader, ok := handler.(Reader); ok {
			r.readers = append(r.readers, reader)
		}
		if writer, ok := handler.(Writer); ok {
			r.writers = append(r.writers, writer)
		}
		if copier, ok := handler.(Copier); ok {
			r.copiers = append(r.copiers, copier)
		}
		if editor, ok := handler.(Editor); ok {
			r.editors = append(r.editors, editor)
		}
		if importer, ok := handler.(Importer); ok {
			r.importers = append(r.importers, importer)
		}
		if exporter, ok := handler.(Exporter); ok {
			r.exporters = append(r.exporters, exporter)
		}
		if controller, ok := handler.(Controller); ok {
			r.controllers = append(r.controllers, controller)
		}
		if inspector, ok := handler.(Inspector); ok {
			r.inspectors = append(r.inspectors, inspector)
		}
		if dumper, ok := handler.(Dumper); ok {
			r.dumpers = append(r.dumpers, dumper)
		}
		if loader, ok := handler.(Loader); ok {
			r.loaders = append(r.loaders, loader)
		}
	}

	return r
}

// Readers returns a copy of all registered Reader handlers.
// Returns a copy to prevent external mutation of the registry.
func (r *Registry) Readers() []Reader {
	return append([]Reader(nil), r.readers...)
}

// Writers returns a copy of all registered Writer handlers.
func (r *Registry) Writers() []Writer {
	return append([]Writer(nil), r.writers...)
}

// Copiers returns a copy of all registered Copier handlers.
func (r *Registry) Copiers() []Copier {
	return append([]Copier(nil), r.copiers...)
}

// Editors returns a copy of all registered Editor handlers.
func (r *Registry) Editors() []Editor {
	return append([]Editor(nil), r.editors...)
}

// Importers returns a copy of all registered Importer handlers.
func (r *Registry) Importers() []Importer {
	return append([]Importer(nil), r.importers...)
}

// Exporters returns a copy of all registered Exporter handlers.
func (r *Registry) Exporters() []Exporter {
	return append([]Exporter(nil), r.exporters...)
}

// Controllers returns a copy of all registered Controller handlers.
func (r *Registry) Controllers() []Controller {
	return append([]Controller(nil), r.controllers...)
}

// Inspectors returns a copy of all registered Inspector handlers.
func (r *Registry) Inspectors() []Inspector {
	return append([]Inspector(nil), r.inspectors...)
}

// Dumpers returns a copy of all registered Dumper handlers.
func (r *Registry) Dumpers() []Dumper {
	return append([]Dumper(nil), r.dumpers...)
}

// Loaders returns a copy of all registered Loader handlers.
func (r *Registry) Loaders() []Loader {
	return append([]Loader(nil), r.loaders...)
}
