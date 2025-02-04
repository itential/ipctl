// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logger

import (
	"os"

	"github.com/itential/ipctl/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Used by customConsoleWriter methods
var (
	debugOut zerolog.ConsoleWriter
	errorOut zerolog.ConsoleWriter
)

// Define our own writers so that we can
type customConsoleWriter struct{}
type customConsoleJsonWriter struct{}

// EnableConsoleLogs Enables console logs that will be displayed to the CLI user via stdout and
// stderr.
func EnableConsoleLogs(cfg *config.Config) {
	if cfg.LogConsoleJSON {
		iowriters = append(iowriters, customConsoleJsonWriter{})
	} else {
		debugOut = zerolog.NewConsoleWriter()
		errorOut = zerolog.NewConsoleWriter()

		// format/timezone is set globally for zerolog in InitializeLogger, but we need to explicit set it here
		// as NewConsoleWriter overrides that assigment
		debugOut.FormatTimestamp = timestampFormatter(cfg.LogTimestampTimezone)
		errorOut.FormatTimestamp = timestampFormatter(cfg.LogTimestampTimezone)

		debugOut.NoColor = cfg.TerminalNoColor
		errorOut.NoColor = cfg.TerminalNoColor

		debugOut.Out = os.Stdout
		errorOut.Out = os.Stderr

		iowriters = append(iowriters, customConsoleWriter{})
	}

	writers := zerolog.MultiLevelWriter(iowriters...)
	log.Logger = zerolog.New(writers).With().Timestamp().Logger()
}

// Write is required to implement io.Writer and should not be called
func (l customConsoleWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// Write is required to implement io.Writer and should not be called
func (l customConsoleJsonWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// WriteLevel Determines the correct output destination depending on the level of the
// request for the customConsoleWriter
func (l customConsoleWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return debugOut.Write(p)
	} else {
		return errorOut.Write(p)
	}
}

// WriteLevel Determines the correct output destination depending on the level of the
// request for the customConsoleJsonWriter
func (l customConsoleJsonWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return os.Stdout.Write(p)
	} else {
		return os.Stderr.Write(p)
	}
}
