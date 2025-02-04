// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package main

import (
	"os"

	"github.com/itential/ipctl/cmd"
)

func main() {
	os.Exit(
		cmd.Execute(),
	)
}
