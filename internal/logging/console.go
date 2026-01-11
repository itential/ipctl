// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// debugOut and errorOut are console writers configured for different output streams.
// debugOut writes to stdout for debug, info, and warn levels.
// errorOut writes to stderr for error and fatal levels.
var (
	debugOut zerolog.ConsoleWriter
	errorOut zerolog.ConsoleWriter
	redactor *Redactor
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
// Sensitive data redaction is applied if enabled in the configuration.
//
// The noColor parameter controls whether console output should disable color codes.
func EnableConsoleLogs(cfg Config, noColor bool) {
	// Initialize the redactor with configuration
	redactor = NewRedactor(cfg.RedactSensitiveData)

	if cfg.ConsoleJSON {
		iowriters = append(iowriters, customConsoleJsonWriter{})
	} else {
		debugOut = zerolog.NewConsoleWriter()
		errorOut = zerolog.NewConsoleWriter()

		// format/timezone is set globally for zerolog in InitializeLogger, but we need to explicit set it here
		// as NewConsoleWriter overrides that assigment
		debugOut.FormatTimestamp = timestampFormatter(cfg.TimestampTimezone)
		errorOut.FormatTimestamp = timestampFormatter(cfg.TimestampTimezone)

		debugOut.NoColor = noColor
		errorOut.NoColor = noColor

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
// Sensitive data is redacted before output if redaction is enabled.
func (l customConsoleWriter) Write(p []byte) (n int, err error) {
	if redactor != nil {
		p = redactor.RedactBytes(p)
	}
	return os.Stdout.Write(p)
}

// Write implements io.Writer interface for customConsoleJsonWriter.
// This method routes all output to stdout and should not be called directly.
// Use WriteLevel instead for proper level-based routing.
// Sensitive data is redacted before output if redaction is enabled.
func (l customConsoleJsonWriter) Write(p []byte) (n int, err error) {
	if redactor != nil {
		p = redactor.RedactBytes(p)
	}
	return os.Stdout.Write(p)
}

// WriteLevel implements zerolog.LevelWriter interface for customConsoleWriter.
// It routes log messages to stdout (for levels <= warn) or stderr (for error and fatal).
// The output is formatted using the configured console writers with proper timestamps and colors.
// Sensitive data is redacted before output if redaction is enabled.
func (l customConsoleWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if redactor != nil {
		p = redactor.RedactBytes(p)
	}
	if level <= zerolog.WarnLevel {
		return debugOut.Write(p)
	} else {
		return errorOut.Write(p)
	}
}

// WriteLevel implements zerolog.LevelWriter interface for customConsoleJsonWriter.
// It routes JSON-formatted log messages to stdout (for levels <= warn) or stderr (for error and fatal).
// This provides structured logging output while maintaining proper stream separation.
// Sensitive data is redacted before output if redaction is enabled.
func (l customConsoleJsonWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if redactor != nil {
		p = redactor.RedactBytes(p)
	}
	if level <= zerolog.WarnLevel {
		return os.Stdout.Write(p)
	} else {
		return os.Stderr.Write(p)
	}
}
