// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCustomConsoleWriterWrite tests the Write method of customConsoleWriter
func TestCustomConsoleWriterWrite(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	writer := customConsoleWriter{}
	testData := []byte("test message")

	n, err := writer.Write(testData)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	_, copyErr := buf.ReadFrom(r)
	require.NoError(t, copyErr)

	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Contains(t, buf.String(), "test message")
}

// TestCustomConsoleJsonWriterWrite tests the Write method of customConsoleJsonWriter
func TestCustomConsoleJsonWriterWrite(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	writer := customConsoleJsonWriter{}
	testData := []byte(`{"level":"info","message":"test"}`)

	n, err := writer.Write(testData)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	_, copyErr := buf.ReadFrom(r)
	require.NoError(t, copyErr)

	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Contains(t, buf.String(), `"level":"info"`)
}

// TestCustomConsoleWriterWriteLevel tests WriteLevel with different log levels
func TestCustomConsoleWriterWriteLevel(t *testing.T) {
	writer := customConsoleWriter{}

	// Use properly formatted JSON log message that console writer can handle
	testData := []byte(`{"level":"info","message":"test message"}`)

	testCases := []struct {
		level zerolog.Level
		name  string
	}{
		{zerolog.TraceLevel, "trace"},
		{zerolog.DebugLevel, "debug"},
		{zerolog.InfoLevel, "info"},
		{zerolog.WarnLevel, "warn"},
		{zerolog.ErrorLevel, "error"},
		{zerolog.FatalLevel, "fatal"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var debugBuf, errorBuf bytes.Buffer
			debugOut = zerolog.ConsoleWriter{Out: &debugBuf, NoColor: true}
			errorOut = zerolog.ConsoleWriter{Out: &errorBuf, NoColor: true}

			n, err := writer.WriteLevel(tc.level, testData)

			assert.NoError(t, err)
			assert.Equal(t, len(testData), n)

			// The actual implementation checks level <= zerolog.WarnLevel for debugOut
			if tc.level <= zerolog.WarnLevel {
				assert.NotEmpty(t, debugBuf.String(), "Expected output to debugOut for level %s", tc.name)
				assert.Empty(t, errorBuf.String(), "Expected no output to errorOut for level %s", tc.name)
			} else {
				assert.Empty(t, debugBuf.String(), "Expected no output to debugOut for level %s", tc.name)
				assert.NotEmpty(t, errorBuf.String(), "Expected output to errorOut for level %s", tc.name)
			}
		})
	}
}

// TestCustomConsoleJsonWriterWriteLevel tests WriteLevel for JSON writer
func TestCustomConsoleJsonWriterWriteLevel(t *testing.T) {
	writer := customConsoleJsonWriter{}
	testData := []byte(`{"level":"info","message":"test"}`)

	testCases := []struct {
		level        zerolog.Level
		expectStdout bool
		name         string
	}{
		{zerolog.TraceLevel, true, "trace"},
		{zerolog.DebugLevel, true, "debug"},
		{zerolog.InfoLevel, true, "info"},
		{zerolog.WarnLevel, true, "warn"},
		{zerolog.ErrorLevel, false, "error"},
		{zerolog.FatalLevel, false, "fatal"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr

			rOut, wOut, err := os.Pipe()
			require.NoError(t, err)
			rErr, wErr, err := os.Pipe()
			require.NoError(t, err)

			os.Stdout = wOut
			os.Stderr = wErr

			n, writeErr := writer.WriteLevel(tc.level, testData)

			// Close writers and restore
			wOut.Close()
			wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// Read captured output
			var stdoutBuf, stderrBuf bytes.Buffer
			_, copyErr1 := stdoutBuf.ReadFrom(rOut)
			_, copyErr2 := stderrBuf.ReadFrom(rErr)
			require.NoError(t, copyErr1)
			require.NoError(t, copyErr2)

			assert.NoError(t, writeErr)
			assert.Equal(t, len(testData), n)

			if tc.expectStdout {
				assert.NotEmpty(t, stdoutBuf.String(), "Expected output to stdout for level %s", tc.name)
				assert.Empty(t, stderrBuf.String(), "Expected no output to stderr for level %s", tc.name)
			} else {
				assert.Empty(t, stdoutBuf.String(), "Expected no output to stdout for level %s", tc.name)
				assert.NotEmpty(t, stderrBuf.String(), "Expected output to stderr for level %s", tc.name)
			}
		})
	}
}

// TestEnableConsoleLogsJSON tests EnableConsoleLogs with JSON format
func TestEnableConsoleLogsJSON(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		ConsoleJSON:       true,
		TimestampTimezone: time.UTC,
	}

	EnableConsoleLogs(cfg, false)

	assert.Len(t, iowriters, 1)
	// Check that the writer is of the correct type
	_, ok := iowriters[0].(customConsoleJsonWriter)
	assert.True(t, ok, "Expected customConsoleJsonWriter")

	// Verify that logger is configured
	assert.NotNil(t, log.Logger)
}

// TestEnableConsoleLogsConsole tests EnableConsoleLogs with console format
func TestEnableConsoleLogsConsole(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		ConsoleJSON:       false,
		TimestampTimezone: time.UTC,
	}

	EnableConsoleLogs(cfg, true) // Test no-color setting

	assert.Len(t, iowriters, 1)
	// Check that the writer is of the correct type
	_, ok := iowriters[0].(customConsoleWriter)
	assert.True(t, ok, "Expected customConsoleWriter")

	// Verify console writers are configured
	assert.NotNil(t, debugOut)
	assert.NotNil(t, errorOut)
	assert.True(t, debugOut.NoColor, "Expected NoColor to be true")
	assert.True(t, errorOut.NoColor, "Expected NoColor to be true")
	assert.Equal(t, os.Stdout, debugOut.Out)
	assert.Equal(t, os.Stderr, errorOut.Out)

	// Verify that logger is configured
	assert.NotNil(t, log.Logger)
}

// TestEnableConsoleLogsWithColor tests EnableConsoleLogs with color enabled
func TestEnableConsoleLogsWithColor(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		ConsoleJSON:       false,
		TimestampTimezone: time.UTC, // Test with color enabled
	}

	EnableConsoleLogs(cfg, false)

	assert.Len(t, iowriters, 1)

	// Verify console writers have color enabled
	assert.False(t, debugOut.NoColor, "Expected NoColor to be false")
	assert.False(t, errorOut.NoColor, "Expected NoColor to be false")
}

// TestEnableConsoleLogsTimestamp tests timestamp formatting configuration
func TestEnableConsoleLogsTimestamp(t *testing.T) {
	// Reset global state
	iowriters = nil

	// Test with specific timezone
	est, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	cfg := Config{
		ConsoleJSON:       false,
		TimestampTimezone: est,
	}

	EnableConsoleLogs(cfg, false)

	// Verify timestamp formatters are set
	assert.NotNil(t, debugOut.FormatTimestamp)
	assert.NotNil(t, errorOut.FormatTimestamp)

	// Test the timestamp formatter
	testTimestamp := "2023-01-01T12:00:00Z"
	formattedDebug := debugOut.FormatTimestamp(testTimestamp)
	formattedError := errorOut.FormatTimestamp(testTimestamp)

	// Should be converted to EST
	assert.Contains(t, formattedDebug, "07:00:00")
	assert.Contains(t, formattedError, "07:00:00")
}

// TestConsoleWriterTypeImplementations tests that custom writers implement required interfaces
func TestConsoleWriterTypeImplementations(t *testing.T) {
	// Test that customConsoleWriter implements required interfaces
	var consoleWriter interface{} = customConsoleWriter{}
	_, implementsWriter := consoleWriter.(zerolog.LevelWriter)
	assert.True(t, implementsWriter, "customConsoleWriter should implement zerolog.LevelWriter")

	// Test that customConsoleJsonWriter implements required interfaces
	var jsonWriter interface{} = customConsoleJsonWriter{}
	_, implementsWriter = jsonWriter.(zerolog.LevelWriter)
	assert.True(t, implementsWriter, "customConsoleJsonWriter should implement zerolog.LevelWriter")
}

// TestMultipleEnableConsoleLogsCalls tests multiple calls to EnableConsoleLogs
func TestMultipleEnableConsoleLogsCalls(t *testing.T) {
	// Reset global state
	iowriters = nil

	cfg := Config{
		ConsoleJSON:       false,
		TimestampTimezone: time.UTC,
	}

	// First call
	EnableConsoleLogs(cfg, false)
	firstLength := len(iowriters)

	// Second call
	EnableConsoleLogs(cfg, false)
	secondLength := len(iowriters)

	// Should append another writer
	assert.Equal(t, firstLength+1, secondLength)
	assert.Equal(t, 2, secondLength)
}
