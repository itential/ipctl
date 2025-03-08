// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"reflect"
)

var registry []any

func register(types ...any) {
	for _, ele := range types {
		registry = append(registry, ele)
	}
}

func Readers() []Reader {
	resource := reflect.TypeOf((*Reader)(nil)).Elem()
	var resources []Reader
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Reader))
		}
	}
	return resources
}

func Writers() []Writer {
	resource := reflect.TypeOf((*Writer)(nil)).Elem()
	var resources []Writer
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Writer))
		}
	}
	return resources
}

func Copiers() []Copier {
	resource := reflect.TypeOf((*Copier)(nil)).Elem()
	var resources []Copier
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Copier))
		}
	}
	return resources
}

func Editors() []Editor {
	resource := reflect.TypeOf((*Editor)(nil)).Elem()
	var resources []Editor
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Editor))
		}
	}
	return resources
}

func Importers() []Importer {
	resource := reflect.TypeOf((*Importer)(nil)).Elem()
	var resources []Importer
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Importer))
		}
	}
	return resources
}

func Exporters() []Exporter {
	resource := reflect.TypeOf((*Exporter)(nil)).Elem()
	var resources []Exporter
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Exporter))
		}
	}
	return resources
}

func Controllers() []Controller {
	resource := reflect.TypeOf((*Controller)(nil)).Elem()
	var resources []Controller
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Controller))
		}
	}
	return resources
}

func Inspectors() []Inspector {
	resource := reflect.TypeOf((*Inspector)(nil)).Elem()
	var resources []Inspector
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(resource) {
			resources = append(resources, ele.(Inspector))
		}
	}
	return resources
}

func Dumpers() []Dumper {
	res := reflect.TypeOf((*Dumper)(nil)).Elem()
	var resources []Dumper
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(res) {
			resources = append(resources, ele.(Dumper))
		}
	}
	return resources
}

func Loaders() []Loader {
	res := reflect.TypeOf((*Loader)(nil)).Elem()
	var resources []Loader
	for _, ele := range registry {
		if reflect.TypeOf(ele).Implements(res) {
			resources = append(resources, ele.(Loader))
		}
	}
	return resources
}
