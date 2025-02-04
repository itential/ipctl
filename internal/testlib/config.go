// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package testlib

import (
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultFileName = "config"

	defaultAppWorkingDir     = "~/.iap.d/test"
	defaultAppDefaultProfile = ""

	defaultLogLevel             = "INFO"
	defaultLogFileJson          = false
	defaultLogConsoleJson       = false
	defaultLogFileEnabled       = false
	defaultLogTimestampTimezone = "utc"

	defaultTerminalNoColor           = false
	defaultTerminalTimestampTimezone = "utc"
)

var defaultValues = map[string]interface{}{
	"application.working_dir":     defaultAppWorkingDir,
	"application.default_profile": defaultAppDefaultProfile,

	"log.level":              defaultLogLevel,
	"log.file_json":          defaultLogFileJson,
	"log.console_json":       defaultLogConsoleJson,
	"log.file_enabled":       defaultLogFileEnabled,
	"log.timestamp_timezone": defaultLogTimestampTimezone,

	"terminal.no_color":           defaultTerminalNoColor,
	"terminal.timestamp_timezone": defaultTerminalTimestampTimezone,
}

// DefaultConfig is mainly used for configuration testing. This will
// set a sane default that then can be overridden within the specific test.
func DefaultConfig() *config.Config {
	var ac config.Config

	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}

	var err error

	ac.WorkingDir, err = homedir.Expand(ac.WorkingDir)
	if err != nil {
		logger.Fatal(err, "error attemping to expand home directory")
	}

	return &ac
}
