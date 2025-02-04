// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmdutils

import (
	"strings"

	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/logger"
	"gopkg.in/yaml.v2"
)

type Descriptor struct {
	Use string `yaml:"use"`

	Group string `yaml:"group"`

	Description string `yaml:"description"`

	Example string `yaml:"example"`

	IncludeGroups bool `yaml:"include_groups"`

	ExactArgs int `yaml:"exact_args"`

	Disabled bool `yaml:"disabled"`

	Hidden bool `yaml:"hidden"`
}

func (d Descriptor) Short() string {
	return strings.Split(d.Description, "\n")[0]
}

func LoadDescriptor(data []byte) map[string]Descriptor {
	var s map[string]Descriptor
	if err := yaml.Unmarshal(data, &s); err != nil {
		logger.Fatal(err, "failed to unmarshal data")
	}
	return s
}

func LoadDescriptorFromString(s string) map[string]Descriptor {
	return LoadDescriptor([]byte(s))
}

// checkError functions very similarly to cobra's builtin cobra.checkErr
// function. Our function adds extra functionality such as colored outputs
// and logging the error to the loggers
func CheckError(err error, terminalNoColor bool) {
	if err != nil {
		terminal.Error(err, terminalNoColor)
		terminal.Display("")
		logger.Fatal(err, "")
	}
}
