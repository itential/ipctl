// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmd

import (
	"embed"
	"strings"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/pkg/logger"
)

//go:embed descriptors/*.yaml
var content embed.FS

type Descriptors map[string]map[string]cmdutils.Descriptor

func loadDescriptors(path string) Descriptors {
	descriptors := map[string]map[string]cmdutils.Descriptor{}

	entries, err := content.ReadDir(path)
	if err != nil {
		logger.Fatal(err, "failed to read descriptors directory")
	}

	for _, ele := range entries {
		name := strings.Split(ele.Name(), ".")[0]
		fn := strings.Join([]string{path, ele.Name()}, "/")

		data, err := content.ReadFile(fn)
		if err != nil {
			logger.Fatal(err, "failed to read descriptor")
		}

		descriptors[name] = cmdutils.LoadDescriptor(data)
	}

	return descriptors
}
