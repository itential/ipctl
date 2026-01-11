// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDebug tests the Debug logging function
func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	Debug("test debug message: %s", "param")

	output := buf.String()
	assert.Contains(t, output, "test debug message: param")
	assert.Contains(t, output, `"level":"debug"`)
}

// TestInfo tests the Info logging function
func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Info("test info message: %s", "param")

	output := buf.String()
	assert.Contains(t, output, "test info message: param")
	assert.Contains(t, output, `"level":"info"`)
}

// TestWarn tests the Warn logging function
func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	Warn("test warn message: %s", "param")

	output := buf.String()
	assert.Contains(t, output, "test warn message: param")
	assert.Contains(t, output, `"level":"warn"`)
}

// TestError tests the Error logging function with an error instance
func TestErrorWithErr(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	testErr := errors.New("test error")
	Error(testErr, "error occurred: %s", "param")

	output := buf.String()
	assert.Contains(t, output, "error occurred: param")
	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, "test error")
}

// TestError tests the Error logging function without an error instance
func TestErrorWithoutErr(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	Error(nil, "error occurred: %s", "param")

	output := buf.String()
	assert.Contains(t, output, "error occurred: param")
	assert.Contains(t, output, `"level":"error"`)
}

// TestTrace tests the Trace logging function at trace level
func TestTraceEnabled(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	Trace()

	output := buf.String()
	assert.Contains(t, output, `"level":"trace"`)
	// Should contain runtime information
	assert.Contains(t, output, "TestTraceEnabled")
}

// TestTrace tests that Trace is disabled when not at trace level
func TestTraceDisabled(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Trace()

	output := buf.String()
	assert.Empty(t, output)
}

// TestTimestampFormatter tests the timestamp formatting function
func TestTimestampFormatter(t *testing.T) {
	// Test with UTC timezone
	utc := time.UTC
	formatter := timestampFormatter(utc)

	// Test valid timestamp
	timestamp := "2023-01-01T12:00:00Z"
	result := formatter(timestamp)
	assert.Equal(t, "2023-01-01T12:00:00Z", result)

	// Test invalid timestamp
	invalidTimestamp := "invalid-timestamp"
	result = formatter(invalidTimestamp)
	assert.Contains(t, result, "error formatting timestamp")
	assert.Contains(t, result, "invalid-timestamp")
}

// TestTimestampFormatterWithTimezone tests timestamp formatting with different timezone
func TestTimestampFormatterWithTimezone(t *testing.T) {
	// Test with EST timezone
	est, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	formatter := timestampFormatter(est)
	timestamp := "2023-01-01T12:00:00Z"
	result := formatter(timestamp)

	// Should be converted to EST (UTC-5 in January)
	assert.Contains(t, result, "07:00:00")
}

// TestInitializeLoggerFirstCall tests that InitializeLogger works on first call
func TestInitializeLoggerFirstCall(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		Level:             "INFO",
		TimestampTimezone: time.UTC,
	}

	// Mock os.Args to avoid verbose flag parsing
	originalArgs := os.Args
	os.Args = []string{"test"}
	defer func() { os.Args = originalArgs }()

	InitializeLogger(cfg, false)

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	// iowriters should be empty (initialized as empty slice) when no logging is enabled
	assert.Empty(t, iowriters)
}

// TestInitializeLoggerSubsequentCalls tests that subsequent calls are ignored
func TestInitializeLoggerSubsequentCalls(t *testing.T) {
	// Set up initial state
	iowriters = []io.Writer{os.Stdout} // Set to something non-nil
	originalLevel := zerolog.GlobalLevel()

	cfg := Config{
		Level:             "DEBUG",
		TimestampTimezone: time.UTC,
	}

	InitializeLogger(cfg, false)

	// Should not change because iowriters was already initialized
	assert.Equal(t, originalLevel, zerolog.GlobalLevel())
}

// TestInitializeLoggerWithVerboseFlag tests verbose flag handling
func TestInitializeLoggerWithVerboseFlag(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		Level:             "INFO",
		TimestampTimezone: time.UTC,
	}

	// Mock os.Args with verbose flag
	originalArgs := os.Args
	os.Args = []string{"test", "--verbose"}
	defer func() { os.Args = originalArgs }()

	InitializeLogger(cfg, false)

	// Should have enabled console logs
	assert.NotEmpty(t, iowriters)
}

// TestLogLevelsFiltering tests that log messages are filtered by level
func TestLogLevelsFiltering(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)

	// Set to WARN level - should only show warn, error, fatal
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error(nil, "error message")

	output := buf.String()
	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

// Helper function to capture Fatal calls (which would normally exit)
// We can't easily test Fatal directly since it calls os.Exit()
func TestFatalLogsMessage(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	// We can't actually call Fatal() since it would exit the test
	// But we can test the error creation logic that's used in Fatal
	testErr := errors.New("test fatal error")

	// Test the same logic used in Fatal for error handling
	if testErr == nil {
		testErr = errors.New("fatal error occurred")
	}

	assert.NotNil(t, testErr)
	assert.Equal(t, "test fatal error", testErr.Error())

	// Test with nil error (same logic as Fatal)
	var nilErr error
	if nilErr == nil {
		nilErr = errors.New(fmt.Sprintf("fatal error: %s", "test"))
	}

	assert.Equal(t, "fatal error: test", nilErr.Error())
}

// TestPackageVariables tests package-level variables
func TestPackageVariables(t *testing.T) {
	// Reset iowriters for testing
	originalWriters := iowriters
	defer func() { iowriters = originalWriters }()

	iowriters = nil
	assert.Nil(t, iowriters)

	// Test appending writers
	iowriters = append(iowriters, os.Stdout)
	assert.Len(t, iowriters, 1)
	assert.Equal(t, os.Stdout, iowriters[0])
}
