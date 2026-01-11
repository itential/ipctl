// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Config holds the configuration options for the logging package.
// It defines log levels, output formats, and timestamp handling.
type Config struct {
	// Level specifies the minimum log level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, DISABLED)
	Level string
	// ConsoleJSON enables JSON format for console output (default: false for human-readable)
	ConsoleJSON bool
	// TimestampTimezone is the location for formatting log timestamps
	TimestampTimezone *time.Location
	// RedactSensitiveData enables automatic redaction of sensitive information in logs
	RedactSensitiveData bool
}

// DefaultConfig returns a Config with sensible defaults.
// Default level is INFO, console output is human-readable, timestamps use UTC,
// and sensitive data redaction is enabled for security.
func DefaultConfig() Config {
	return Config{
		Level:               "INFO",
		ConsoleJSON:         false,
		TimestampTimezone:   time.UTC,
		RedactSensitiveData: true,
	}
}

// LoadFromEnv creates a Config by loading values from environment variables.
// Supported environment variables:
//   - IPCTL_LOG_LEVEL: Log level (DEBUG, INFO, WARN, ERROR, FATAL, DISABLED, TRACE)
//   - IPCTL_LOG_CONSOLE_JSON: Enable JSON console output (true/false)
//   - IPCTL_LOG_TIMESTAMP_TIMEZONE: Timezone for timestamps (UTC, Local, or IANA timezone)
//   - IPCTL_LOG_REDACT_SENSITIVE_DATA: Enable sensitive data redaction (true/false, default: true)
//
// Returns a Config with defaults for any unset environment variables.
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	if level := os.Getenv("IPCTL_LOG_LEVEL"); level != "" {
		cfg.Level = level
	}

	if consoleJSON := os.Getenv("IPCTL_LOG_CONSOLE_JSON"); consoleJSON == "true" {
		cfg.ConsoleJSON = true
	}

	if tz := os.Getenv("IPCTL_LOG_TIMESTAMP_TIMEZONE"); tz != "" {
		if loc := parseTimezone(tz); loc != nil {
			cfg.TimestampTimezone = loc
		}
	}

	if redact := os.Getenv("IPCTL_LOG_REDACT_SENSITIVE_DATA"); redact == "false" {
		cfg.RedactSensitiveData = false
	}

	return cfg
}

// supportedLogLevels defines the valid log level strings that can be parsed by getLogLevel.
// These correspond to zerolog's built-in log levels plus a DISABLED option.
var supportedLogLevels = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "DISABLED", "TRACE"}

// parseTimezone converts a timezone string to a time.Location.
// It handles "UTC", "Local", and IANA timezone identifiers.
// Returns nil if the timezone is invalid.
func parseTimezone(tz string) *time.Location {
	var parsedTz string
	switch strings.ToLower(tz) {
	case "utc":
		parsedTz = "UTC"
	case "local":
		parsedTz = "Local"
	default:
		parsedTz = tz
	}

	location, err := time.LoadLocation(parsedTz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "# Warning: failed to load timezone '%s': %v. Defaulting to UTC\n", tz, err)
		return time.UTC
	}
	return location
}

// getLogLevel converts a string log level to a zerolog.Level constant.
// Valid levels are defined in supportedLogLevels. If an invalid level is provided,
// returns an error instead of exiting the application.
func getLogLevel(level string) zerolog.Level {
	s := strings.ToUpper(level)

	// Not using internal/functions.ContainsString to avoid cyclic dependency
	supported := false
	for _, ele := range supportedLogLevels {
		if ele == s {
			supported = true
			break
		}
	}

	if !supported {
		fmt.Fprintf(
			os.Stderr,
			"invalid value for IPCTL_LOG_LEVEL, got %s, expected one of %s. Defaulting to INFO\n",
			s, strings.Join(supportedLogLevels, ","),
		)
		return zerolog.InfoLevel
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
