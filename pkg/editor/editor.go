// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package editor

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"

	"github.com/itential/ipctl/internal/logging"
	"github.com/mitchellh/mapstructure"
)

const defaultEditor = "vi"

func Run(in interface{}, ptr any) error {
	logging.Trace()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	// Create the temp directory
	tempFile, err := os.CreateTemp("", "iag*.json")
	if err != nil {
		logging.Fatal(err, "could not create temp file")
	}
	defer os.Remove(tempFile.Name())

	// Marshal the object to bytes and write it to a file in the temp directory
	b, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		logging.Fatal(err, "failed to marshal object")
	}

	if _, err := tempFile.Write(b); err != nil {
		logging.Fatal(err, "could not write temp file")

	}
	if err := tempFile.Close(); err != nil {
		logging.Fatal(err, "could not close temp file")
	}

	// Create the editor command and launch it
	editorCmd := exec.Command(editor, tempFile.Name())

	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		logging.Fatal(err, "could not open file in editor")
	}

	// Read the contents of the updated file back
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		logging.Fatal(err, "failed to temp file")
	}

	if ptr != nil {
		err = json.Unmarshal(data, &ptr)
		if err != nil {
			logging.Fatal(err, "failed to unmarshal temp file")
		}

		var out = reflect.Zero(reflect.TypeOf(in)).Interface()
		mapstructure.Decode(ptr, &out)
	}

	return nil
}
