// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// supportedLogLevels defines the valid log level strings that can be parsed by getLogLevel.
// These correspond to zerolog's built-in log levels plus a DISABLED option.
var supportedLogLevels = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "DISABLED", "TRACE"}

// getLogLevel converts a string log level to a zerolog.Level constant.
// Valid levels are defined in supportedLogLevels. If an invalid level is provided,
// the function prints an error message and exits the application with status code 1.
func getLogLevel(level string) zerolog.Level {
	s := strings.ToUpper(level)

	// Not using internal/functions.ContainsString to avoid cyclic dependency
	supported := false
	for _, ele := range supportedLogLevels {
		if ele == s {
			supported = true
		}
	}

	if !supported {
		fmt.Fprintf(
			os.Stderr,
			"invalid value for IAGCTL_LOG_LEVEL, got %s, expected one of %s",
			s, strings.Join(supportedLogLevels, ","),
		)
		os.Exit(1)
	}

	var z zerolog.Level
	switch s {
	case "TRACE":
		z = zerolog.TraceLevel
	case "DEBUG":
		z = zerolog.DebugLevel
	case "INFO":
		z = zerolog.InfoLevel
	case "WARN":
		z = zerolog.WarnLevel
	case "ERROR":
		z = zerolog.ErrorLevel
	case "FATAL":
		z = zerolog.FatalLevel
	case "DISABLED":
		z = zerolog.Disabled
	}
	return z
}
