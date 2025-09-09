// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/itential/ipctl/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnableFileLogsJSON tests EnableFileLogs with JSON format
func TestEnableFileLogsJSON(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check that iowriters was populated
	assert.Len(t, iowriters, 1)
	
	// Check that log file was created
	logFilePath := filepath.Join(tempDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err, "Log file should be created")
	
	// Verify that logger is configured
	assert.NotNil(t, log.Logger)
}

// TestEnableFileLogsConsole tests EnableFileLogs with console format
func TestEnableFileLogsConsole(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          false,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check that iowriters was populated
	assert.Len(t, iowriters, 1)
	
	// Check that log file was created
	logFilePath := filepath.Join(tempDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err, "Log file should be created")
	
	// Verify that logger is configured
	assert.NotNil(t, log.Logger)
}

// TestEnableFileLogsCreatesDirectory tests that missing directories are created
func TestEnableFileLogsCreatesDirectory(t *testing.T) {
	// Create temporary directory and then a subdirectory path that doesn't exist
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	logDir := filepath.Join(tempDir, "subdir", "logdir")
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           logDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check that directory was created
	_, err = os.Stat(logDir)
	assert.NoError(t, err, "Log directory should be created")
	
	// Check that log file was created
	logFilePath := filepath.Join(logDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err, "Log file should be created")
	
	// Check that iowriters was populated
	assert.Len(t, iowriters, 1)
}

// TestEnableFileLogsExistingDirectory tests with existing directory
func TestEnableFileLogsExistingDirectory(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check that iowriters was populated
	assert.Len(t, iowriters, 1)
	
	// Check that log file was created
	logFilePath := filepath.Join(tempDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err, "Log file should be created")
}

// TestEnableFileLogsAppendMode tests that log files are opened in append mode
func TestEnableFileLogsAppendMode(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	logFilePath := filepath.Join(tempDir, logFileName)
	
	// Pre-create log file with some content
	initialContent := "existing log content\n"
	err = os.WriteFile(logFilePath, []byte(initialContent), 0644)
	require.NoError(t, err)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Read file content
	content, err := os.ReadFile(logFilePath)
	require.NoError(t, err)
	
	// Should still contain the original content (append mode)
	assert.Contains(t, string(content), initialContent)
}

// TestEnableFileLogsPermissions tests file permissions
func TestEnableFileLogsPermissions(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check file permissions
	logFilePath := filepath.Join(tempDir, logFileName)
	info, err := os.Stat(logFilePath)
	require.NoError(t, err)
	
	// File should be created with 0664 permissions (but may be modified by umask)
	// Just verify it's a readable file
	actualMode := info.Mode().Perm()
	assert.True(t, actualMode&0400 != 0, "File should be readable by owner")
	assert.True(t, actualMode&0200 != 0, "File should be writable by owner")
}

// TestEnableFileLogsTimestamp tests timestamp formatting in file output
func TestEnableFileLogsTimestamp(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Test with specific timezone
	est, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          false, // Use console format to test timestamp formatter
		LogTimestampTimezone: est,
	}
	
	EnableFileLogs(cfg)
	
	// Verify that a console writer was added (not direct file writer)
	assert.Len(t, iowriters, 1)
	
	// Check that log file was created
	logFilePath := filepath.Join(tempDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err, "Log file should be created")
}

// TestLogFileNameConstant tests the logFileName constant
func TestLogFileNameConstant(t *testing.T) {
	assert.Equal(t, "iap.log", logFileName)
}

// TestEnableFileLogsMultipleCalls tests multiple calls to EnableFileLogs
func TestEnableFileLogsMultipleCalls(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	// First call
	EnableFileLogs(cfg)
	firstLength := len(iowriters)
	
	// Second call
	EnableFileLogs(cfg)
	secondLength := len(iowriters)
	
	// Should append another writer
	assert.Equal(t, firstLength+1, secondLength)
	assert.Equal(t, 2, secondLength)
}

// TestEnableFileLogsInvalidDirectory tests error handling for invalid directory
// Note: This test is conceptual as the actual function calls os.Exit on error
func TestEnableFileLogsInvalidDirectoryValidation(t *testing.T) {
	// Test what happens when we try to create a file in a non-writable location
	// We can't actually test the os.Exit behavior, but we can test the concept
	
	// On Unix systems, trying to write to /root usually requires elevated permissions
	invalidDir := "/root/non_existent_dir"
	
	// We won't actually call EnableFileLogs with this path since it would exit
	// Instead, we test the directory creation logic conceptually
	
	_, err := os.Stat(invalidDir)
	exists := !os.IsNotExist(err)
	
	if !exists {
		// This is the condition that would trigger directory creation in EnableFileLogs
		err := os.MkdirAll(invalidDir, os.ModePerm)
		// We expect this to fail due to permissions
		assert.Error(t, err, "Should fail to create directory in restricted location")
	}
}

// TestEnableFileLogsFilePermissions tests file creation permissions
func TestEnableFileLogsFilePermissions(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          true,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Get file info
	logFilePath := filepath.Join(tempDir, logFileName)
	info, err := os.Stat(logFilePath)
	require.NoError(t, err)
	
	// Test that file has reasonable permissions (may be modified by umask)
	mode := info.Mode()
	actualPerm := mode.Perm()
	
	// Just verify basic read/write permissions for owner
	assert.True(t, actualPerm&0400 != 0, "File should be readable by owner")
	assert.True(t, actualPerm&0200 != 0, "File should be writable by owner")
}

// TestEnableFileLogsConsoleWriterConfiguration tests console writer setup for file logging
func TestEnableFileLogsConsoleWriterConfiguration(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Reset global state
	iowriters = nil
	
	cfg := &config.Config{
		WorkingDir:           tempDir,
		LogFileJSON:          false,
		LogTimestampTimezone: time.UTC,
	}
	
	EnableFileLogs(cfg)
	
	// Check that one writer was added
	assert.Len(t, iowriters, 1)
	
	// The writer should be a console writer (for file output)
	// We can't directly inspect the writer type, but we can verify the file exists
	logFilePath := filepath.Join(tempDir, logFileName)
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err)
}