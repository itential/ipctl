// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

const defaultEncoding = "json"

type ExportObject struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Path     string `json:"path"`
	AbsPath  string `json:"abspath"`
}

type ExportOption func(*ExportObject)

func Export(in any, opts ...ExportOption) (*ExportObject, error) {
	logger.Trace()

	obj := new(ExportObject)

	for _, opt := range opts {
		opt(obj)
	}

	if obj.Type == "" {
		return nil, errors.New("export type not specified")
	}

	if obj.Path == "" {
		path, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		obj.Path = path
	}

	if obj.Encoding == "" {
		obj.Encoding = defaultEncoding
	}

	fn := fmt.Sprintf("%s.%s.json", obj.Name, obj.Type)

	obj.AbsPath = filepath.Join(obj.Path, fn)

	if err := utils.WriteJsonToDisk(in, fn, obj.Path); err != nil {
		return nil, err
	}

	return obj, nil
}

func WithExportName(v string) ExportOption {
	return func(obj *ExportObject) {
		obj.Name = v
	}
}

func WithExportType(v string) ExportOption {
	return func(obj *ExportObject) {
		obj.Type = v
	}
}

func WithExportEncoding(v string) ExportOption {
	return func(obj *ExportObject) {
		obj.Encoding = v
	}
}

func WithExportPath(v string) ExportOption {
	return func(obj *ExportObject) {
		obj.Path = v
	}
}
