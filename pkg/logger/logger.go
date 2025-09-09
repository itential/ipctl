// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/itential/ipctl/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

// iowriters holds all configured io.Writer instances for logging output.
// Writers are appended by EnableConsoleLogs and EnableFileLogs functions.
var iowriters []io.Writer

// InitializeLogger sets up the global logger configuration based on the provided config.
// It configures console and/or file logging, sets the log level, and handles timezone formatting.
// The logger is initialized only once - subsequent calls are ignored.
func InitializeLogger(cfg *config.Config) {
	if iowriters != nil {
		return
	}
	var verbose bool
	log.Logger = zerolog.New(io.Discard)
	loggerFlags := pflag.NewFlagSet("loggerFlags", pflag.ContinueOnError)
	loggerFlags.Usage = func() {} // Prevents help message
	loggerFlags.BoolVar(&verbose, "verbose", false, "")

	// Prevents errors when parsing flags not set in 'loggerFlags'
	loggerFlags.ParseErrorsWhitelist.UnknownFlags = true
	// Nil won't trigger as 0 is executable. This would just return an empty slice.
	if err := loggerFlags.Parse(os.Args[1:]); err != nil && err != pflag.ErrHelp {
		fmt.Fprintf(os.Stderr, "%s. Failed to parse verbose command line argument\n", err.Error())
		os.Exit(1)
	}

	if verbose {
		EnableConsoleLogs(cfg)
	}
	if cfg.LogFileEnabled {
		EnableFileLogs(cfg)
	}

	// Only sets timezone on JSON loggers as they reference default zerolog settings
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(cfg.LogTimestampTimezone)
	}

	zerolog.SetGlobalLevel(getLogLevel(cfg.LogLevel))
}

// Trace creates an extremely verbose log message for debugging purposes.
// It automatically captures runtime information including package, method, file, and line number.
// The message format is: pkg.<rcvr>.Method.file.lineNumber
// Only generates output when the global log level is set to TraceLevel.
func Trace() {
	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		pc, _, _, _ := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		f, l := details.FileLine(pc)
		c := strings.Split(f, "/")
		// Nil won't trigger as trace will always have a path with /
		n := strings.Split(c[len(c)-1], ".")[0]
		log.Trace().Msgf("%s.%v.%v",
			strings.Join(strings.Split(details.Name(), "/")[3:], "/"), n, l)
	}
}

// Debug creates a log message with detailed information for debugging purposes.
// Use this for verbose output that helps with troubleshooting but isn't needed in production.
func Debug(format string, args ...any) {
	log.Debug().Msgf(format, args...)
}

// Info creates a log message for important operational information.
// Use this for significant application events like API calls, startup information,
// and major state changes that are relevant for monitoring and troubleshooting.
func Info(format string, args ...any) {
	log.Info().Msgf(format, args...)
}

// Warn creates a log message for concerning situations that don't constitute errors.
// Use this for conditions that are unusual but recoverable, or that may indicate
// potential problems that should be monitored.
func Warn(format string, args ...any) {
	log.Warn().Msgf(format, args...)
}

// Error creates a log message for error conditions that have occurred.
// If an error instance is available, it should be provided for structured logging.
// If no error instance exists, pass nil and the function will create one from the message.
// Use this for recoverable errors and operational failures.
func Error(err error, format string, args ...any) {
	if err == nil {
		err = errors.New(fmt.Sprintf(format, args...))
	}
	log.Error().Err(err).Msgf(format, args...)
}

// Fatal creates a log message for critical errors that require immediate application shutdown.
// This function logs the error and terminates the application by calling os.Exit(1).
// If an error instance is available, it should be provided for structured logging.
// If no error instance exists, pass nil and the function will create one from the message.
// Use this only for unrecoverable errors that make continued operation impossible.
func Fatal(err error, format string, args ...any) {
	if err == nil {
		err = errors.New(fmt.Sprintf(format, args...))
	}
	log.Fatal().Err(err).Msgf(format, args...)
}

// timestampFormatter returns a function that formats timestamps according to the specified timezone.
// It converts RFC3339 formatted timestamps to the configured timezone while maintaining RFC3339 format.
// If timestamp parsing fails, it returns an error message with the original timestamp.
func timestampFormatter(loc *time.Location) func(interface{}) string {
	return func(timestamp interface{}) string {
		t, err := time.Parse(time.RFC3339, timestamp.(string))
		if err != nil { // If we run into an error just return an unformatted timestamp.
			return fmt.Sprintf("%v: error formatting timestamp: %s:", timestamp, err.Error())
		}
		return t.In(loc).Format(time.RFC3339)
	}
}
