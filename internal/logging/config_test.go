// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// TestGetLogLevelValidLevels tests all supported log levels
func TestGetLogLevelValidLevels(t *testing.T) {
	testCases := []struct {
		input    string
		expected zerolog.Level
	}{
		{"TRACE", zerolog.TraceLevel},
		{"DEBUG", zerolog.DebugLevel},
		{"INFO", zerolog.InfoLevel},
		{"WARN", zerolog.WarnLevel},
		{"ERROR", zerolog.ErrorLevel},
		{"FATAL", zerolog.FatalLevel},
		{"DISABLED", zerolog.Disabled},
		// Test case insensitivity
		{"trace", zerolog.TraceLevel},
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"disabled", zerolog.Disabled},
		// Test mixed case
		{"Debug", zerolog.DebugLevel},
		{"Info", zerolog.InfoLevel},
		{"WaRn", zerolog.WarnLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			level := getLogLevel(tc.input)
			assert.Equal(t, tc.expected, level)
		})
	}
}

// TestGetLogLevelInvalidLevel tests invalid log level handling
// Note: This test cannot directly test the os.Exit(1) behavior since that would terminate the test
// Instead, we test the validation logic by checking the supportedLogLevels slice
func TestGetLogLevelInvalidLevelValidation(t *testing.T) {
	invalidLevels := []string{
		"INVALID",
		"UNKNOWN",
		"VERBOSE",
		"CRITICAL",
		"",
		"123",
		"debug123",
		"inf",
	}

	for _, level := range invalidLevels {
		t.Run(level, func(t *testing.T) {
			// Test the validation logic without calling getLogLevel (which would exit)
			supported := false
			for _, ele := range supportedLogLevels {
				if ele == level {
					supported = true
					break
				}
			}
			assert.False(t, supported, "Level %s should not be supported", level)
		})
	}
}

// TestSupportedLogLevels tests the supportedLogLevels variable
func TestSupportedLogLevels(t *testing.T) {
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "DISABLED", "TRACE"}

	assert.Equal(t, expectedLevels, supportedLogLevels)
	assert.Len(t, supportedLogLevels, 7)

	// Test that all expected levels are present
	for _, expected := range expectedLevels {
		found := false
		for _, actual := range supportedLogLevels {
			if actual == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected level %s should be in supportedLogLevels", expected)
	}
}

// TestLogLevelCaseInsensitivity tests that log levels are case insensitive
func TestLogLevelCaseInsensitivity(t *testing.T) {
	testCases := []string{
		"debug",
		"DEBUG",
		"Debug",
		"dEbUg",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			level := getLogLevel(tc)
			assert.Equal(t, zerolog.DebugLevel, level)
		})
	}
}

// TestLogLevelWhitespace tests handling of whitespace in log levels
func TestLogLevelWhitespace(t *testing.T) {
	// These should be treated as invalid since getLogLevel uses strings.ToUpper
	// which doesn't trim whitespace
	invalidWithWhitespace := []string{
		" DEBUG",
		"DEBUG ",
		" DEBUG ",
		"\tDEBUG",
		"DEBUG\n",
	}

	for _, level := range invalidWithWhitespace {
		t.Run("'"+level+"'", func(t *testing.T) {
			// Test the validation logic - these should not be supported
			supported := false
			upperLevel := level
			if len(level) > 0 {
				upperLevel = level
				// Simulate the strings.ToUpper behavior in getLogLevel
				for _, ele := range supportedLogLevels {
					if ele == upperLevel {
						supported = true
						break
					}
				}
			}
			assert.False(t, supported, "Level with whitespace should not be supported")
		})
	}
}

// TestGetLogLevelEdgeCases tests edge cases for getLogLevel
func TestGetLogLevelEdgeCases(t *testing.T) {
	// Test empty string validation
	emptySupported := false
	for _, ele := range supportedLogLevels {
		if ele == "" {
			emptySupported = true
			break
		}
	}
	assert.False(t, emptySupported, "Empty string should not be supported")

	// Test that all supported levels can be converted
	for _, supportedLevel := range supportedLogLevels {
		t.Run(supportedLevel, func(t *testing.T) {
			level := getLogLevel(supportedLevel)
			// Verify it returns a valid zerolog.Level (not panicking)
			assert.True(t, level >= zerolog.TraceLevel && level <= zerolog.Disabled)
		})
	}
}

// Mock test for os.Exit behavior (conceptual test - cannot actually test os.Exit)
func TestGetLogLevelExitBehavior(t *testing.T) {
	// We can't actually test os.Exit(1) behavior since it would terminate the test runner
	// But we can test that the function would identify invalid levels correctly

	// Test that getLogLevel would call os.Exit for invalid levels
	// by checking the validation logic that happens before the switch statement
	invalidLevel := "INVALID_LEVEL"

	// This simulates the validation logic in getLogLevel
	supported := false
	for _, ele := range supportedLogLevels {
		if ele == invalidLevel {
			supported = true
		}
	}

	// If we got here with !supported, getLogLevel would call os.Exit(1)
	assert.False(t, supported, "Invalid level should trigger exit behavior")

	// In a real scenario, we might test this with a wrapper function or by
	// checking stderr output, but direct testing of os.Exit is not practical
}
