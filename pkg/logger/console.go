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

// debugOut and errorOut are console writers configured for different output streams.
// debugOut writes to stdout for debug, info, and warn levels.
// errorOut writes to stderr for error and fatal levels.
var (
	debugOut zerolog.ConsoleWriter
	errorOut zerolog.ConsoleWriter
)

// customConsoleWriter implements zerolog.LevelWriter to route log messages
// to appropriate output streams (stdout/stderr) based on log level.
type customConsoleWriter struct{}

// customConsoleJsonWriter implements zerolog.LevelWriter for JSON formatted console output
// with level-based routing to stdout/stderr.
type customConsoleJsonWriter struct{}

// EnableConsoleLogs configures console-based logging output to stdout and stderr.
// It supports both human-readable console format and JSON format based on configuration.
// Log messages are routed to stdout (debug, info, warn) or stderr (error, fatal) based on level.
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

// Write implements io.Writer interface for customConsoleWriter.
// This method routes all output to stdout and should not be called directly.
// Use WriteLevel instead for proper level-based routing.
func (l customConsoleWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// Write implements io.Writer interface for customConsoleJsonWriter.
// This method routes all output to stdout and should not be called directly.
// Use WriteLevel instead for proper level-based routing.
func (l customConsoleJsonWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// WriteLevel implements zerolog.LevelWriter interface for customConsoleWriter.
// It routes log messages to stdout (for levels <= warn) or stderr (for error and fatal).
// The output is formatted using the configured console writers with proper timestamps and colors.
func (l customConsoleWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return debugOut.Write(p)
	} else {
		return errorOut.Write(p)
	}
}

// WriteLevel implements zerolog.LevelWriter interface for customConsoleJsonWriter.
// It routes JSON-formatted log messages to stdout (for levels <= warn) or stderr (for error and fatal).
// This provides structured logging output while maintaining proper stream separation.
func (l customConsoleJsonWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return os.Stdout.Write(p)
	} else {
		return os.Stderr.Write(p)
	}
}
