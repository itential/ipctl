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

// All loggers will be appended to this slice elsewhere in this package
var iowriters []io.Writer

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

// Trace Creates a log message intended to be used to generate extremely verbose
// messages for debugging purposes. The Format of the trace message can be
// changed how seen fit but it as follows:
// pkg.<rcvr>.Method.file.lineNumber
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

// Debug Creates a log message with fairly verbose, specific information for debugging
// purposes
func Debug(format string, args ...any) {
	log.Debug().Msgf(format, args...)
}

// Info Creates a log message that contains important information about major actions
// that iap performs such as API calls being made, startup info, etc
func Info(format string, args ...any) {
	log.Info().Msgf(format, args...)
}

// Warn Creates a log message that contains information about an occurrence that is
// considered concerning but has not necessarily caused an error
func Warn(format string, args ...any) {
	log.Warn().Msgf(format, args...)
}

// Error Creates a log message with information about the errors that has occurred.
// If an error type is present, it should be provided in addition to a custom
// message for each error. If no relevant error type is present, pass nil for
// err.
func Error(err error, format string, args ...any) {
	if err == nil {
		err = errors.New(fmt.Sprintf(format, args...))
	}
	log.Error().Err(err).Msgf(format, args...)
}

// Fatal Creates an error message with information about the error that has occurred
// which is sever enough that torero will need to immediately shut down. This
// function will log the message and shut torero down by calling `os.Exit(1)`.
// If an error type is present, it should be provided in addition to a custom
// message for each Fatal call. If no relevant error type is present, pass nil
// for err.
func Fatal(err error, format string, args ...any) {
	if err == nil {
		err = errors.New(fmt.Sprintf(format, args...))
	}
	log.Fatal().Err(err).Msgf(format, args...)
}

// timestampFormatter can be passed to a zerolog logger and will reformat any timestamp that it is given to match
// what is given in the config variable TORERO_LOG_TIMESTAMP_TIMEZONE
func timestampFormatter(loc *time.Location) func(interface{}) string {
	return func(timestamp interface{}) string {
		t, err := time.Parse(time.RFC3339, timestamp.(string))
		if err != nil { // If we run into an error just return an unformatted timestamp.
			return fmt.Sprintf("%v: error formatting timestamp: %s:", timestamp, err.Error())
		}
		return t.In(loc).Format(time.RFC3339)
	}
}
