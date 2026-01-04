// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package cmdutils

import (
	"embed"
	"strings"

	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/logger"
	"gopkg.in/yaml.v2"
)

type Descriptors map[string]map[string]Descriptor

// Descriptor describes a command in the application.
type Descriptor struct {
	// Use defines the command
	Use string `yaml:"use"`

	// Group name to put the command into
	Group string `yaml:"group"`

	// The full description of the command. This text is split into a short
	// description and a long description when the command is created. The
	// short description is the first line of the command ending with `\n`.
	// The long form of the description is the entire description text.
	Description string `yaml:"description"`

	// Example provides examples of how to run the command.
	Example string `yaml:"example"`

	// IncludeGroups specifies whether or not to include the default groups as
	// part of the command. This needs to be refactored out of the descriptor.
	IncludeGroups bool `yaml:"include_groups"`

	ExactArgs int `yaml:"exact_args"`

	// Disabled specifies whether or not the command is enabled or disabled.
	// When the command is disabled, it is not added to the command tree.
	Disabled bool `yaml:"disabled"`

	// Hidden specifies whether or not the command is seen when using the
	// `--help` flag.
	Hidden bool `yaml:"hidden"`
}

func (d Descriptor) Short() string {
	return strings.Split(d.Description, "\n")[0]
}

// LoadDescriptor will unmarshal a byte stream into a map of descriptor
// entries.  The `data` argument must be a YAML file that has string based keys
// with each key is a descriptor.
func LoadDescriptor(data []byte) map[string]Descriptor {
	var s map[string]Descriptor
	if err := yaml.Unmarshal(data, &s); err != nil {
		logger.Fatal(err, "failed to unmarshal data")
	}
	return s
}

// LoadDescriptorFromString will unmarhsal a string into a map of descriptor
// entries.
func LoadDescriptorFromString(s string) map[string]Descriptor {
	return LoadDescriptor([]byte(s))
}

// CheckError functions very similarly to cobra's builtin cobra.checkErr
// function. Our function adds extra functionality such as colored outputs
// and logging the error to the loggers
func CheckError(err error, terminalNoColor bool) {
	if err != nil {
		terminal.Error(err, terminalNoColor)
		terminal.Display("")
		logger.Fatal(err, "")
	}
}

// LoadDescriptorsFromContent iterates over the files provided in `content`
// and unmarshales them into maps of descriptors.  Those maps are then added to
// create the `Descriptors` object.  The file map will look like the following:
//
// <filename>:
//	  <key>: <descriptor>
//
// Give a the following filename in the content folder:
// assets.yaml
// ---
// get:
//	use: get
//	description: this is for documentation
//
// The resulting descriptor object would be
//
// assets:
// 	 get: <descriptor>

func LoadDescriptorsFromContent(path string, content *embed.FS) Descriptors {
	logger.Trace()

	descriptors := map[string]map[string]Descriptor{}

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

		descriptors[name] = LoadDescriptor(data)
	}

	return descriptors
}
