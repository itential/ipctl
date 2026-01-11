// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

/*
Package logging provides structured logging functionality with automatic sensitive data redaction.

The package wraps zerolog to provide a simple, consistent logging interface with configurable
output formats (console or JSON), log levels, and timezone handling. It includes built-in
protection against accidentally logging sensitive information such as API keys, passwords,
tokens, and other credentials.

# Basic Usage

Initialize the logger at application startup:

	cfg := logging.LoadFromEnv()
	logging.InitializeLogger(cfg, false)

Log messages at different levels:

	logging.Info("Server started on port %d", 8080)
	logging.Debug("Processing request for user %s", userID)
	logging.Warn("Rate limit approaching for client %s", clientID)
	logging.Error(err, "Failed to connect to database")
	logging.Fatal(err, "Critical error - cannot continue")

# Configuration

The package supports configuration through environment variables:

  - IPCTL_LOG_LEVEL: Set log level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, DISABLED)
  - IPCTL_LOG_CONSOLE_JSON: Enable JSON console output (true/false)
  - IPCTL_LOG_TIMESTAMP_TIMEZONE: Set timezone for timestamps (UTC, Local, or IANA timezone)
  - IPCTL_LOG_REDACT_SENSITIVE_DATA: Enable sensitive data redaction (true/false, default: true)

Example:

	export IPCTL_LOG_LEVEL=DEBUG
	export IPCTL_LOG_CONSOLE_JSON=true
	export IPCTL_LOG_REDACT_SENSITIVE_DATA=true
	./ipctl

Alternatively, create a Config programmatically:

	cfg := logging.Config{
		Level:               "DEBUG",
		ConsoleJSON:         false,
		TimestampTimezone:   time.UTC,
		RedactSensitiveData: true,
	}
	logging.InitializeLogger(cfg, false)

# Sensitive Data Redaction

By default, the logging package automatically scans all output for sensitive information
and replaces it with "<REDACTED>" before writing to stdout or stderr. This protects against
accidental exposure of credentials in logs.

The redactor uses heuristics-based pattern matching to identify:

  - API keys and access tokens
  - Bearer tokens and OAuth tokens
  - JWT tokens
  - Passwords (password=, passwd=, pwd=)
  - AWS credentials (access keys, secret keys)
  - GitHub tokens (personal access tokens, OAuth tokens, app tokens)
  - Client secrets
  - Session tokens and IDs
  - Database connection strings (MongoDB, PostgreSQL)
  - SSH private keys
  - Authorization headers

Example of automatic redaction:

	logging.Info("Connecting with token: Bearer eyJhbGc...")
	// Output: Connecting with token: Bearer <REDACTED>

	logging.Debug("Using password=%s", userPassword)
	// Output: Using password=<REDACTED>

To disable redaction (not recommended for production):

	cfg := logging.Config{
		Level:               "INFO",
		RedactSensitiveData: false,
	}

# Output Routing

Log messages are automatically routed to appropriate streams based on level:

  - stdout: TRACE, DEBUG, INFO, WARN
  - stderr: ERROR, FATAL

This follows Unix conventions and makes it easy to separate normal operational
logs from error logs in production environments.

# Log Levels

The package supports standard log levels with filtering:

  - TRACE: Extremely verbose debugging with runtime information
  - DEBUG: Detailed debugging information
  - INFO: Important operational information (default)
  - WARN: Warning messages for concerning but recoverable situations
  - ERROR: Error conditions that occurred
  - FATAL: Critical errors that require immediate shutdown (calls os.Exit(1))
  - DISABLED: Disable all logging

# Output Formats

Console format (default):
Human-readable output with colors and formatting for terminal display:

	2024-01-15T10:30:45Z INF Server started on port 8080

JSON format (for production):
Structured JSON output for log aggregation and analysis:

	{"level":"info","time":"2024-01-15T10:30:45Z","message":"Server started on port 8080"}

# Architecture

The package consists of several components:

1. Configuration (config.go): Manages logging configuration and environment variables
2. Logger (logger.go): Provides the logging API (Debug, Info, Warn, Error, Fatal)
3. Console Writers (console.go): Route output to stdout/stderr based on log level
4. Redactor (redactor.go): Scans and redacts sensitive information
5. Tests: Comprehensive test coverage for all components

# Performance Considerations

The redactor uses optimized pattern matching with early exit conditions:

1. Quick indicator check before running expensive regex patterns
2. Compiled regex patterns (no runtime compilation)
3. Optional ShouldRedact() method to check if redaction is needed

For applications with very high logging volume, consider:

  - Using JSON format (faster than console formatting)
  - Setting appropriate log level (INFO or WARN in production)
  - Benchmarking with redaction enabled vs disabled

Benchmark results show redaction adds minimal overhead for clean text and
acceptable overhead for text with sensitive data.

# Thread Safety

The package is thread-safe. Zerolog's underlying writer is safe for concurrent use,
and the redactor can be called from multiple goroutines simultaneously.

# Examples

Basic logging:

	logging.Info("Application started successfully")
	logging.Debug("User %s logged in from %s", username, ipAddr)

Error handling:

	if err := doSomething(); err != nil {
		logging.Error(err, "Failed to do something")
		return err
	}

Critical failures:

	if err := connectDB(); err != nil {
		logging.Fatal(err, "Cannot connect to database")
		// Application exits here
	}

Trace logging for detailed debugging:

	func processRequest() {
		logging.Trace() // Logs function name, file, line number
		// ... implementation
	}

# Best Practices

1. Initialize the logger once at application startup
2. Use appropriate log levels (avoid DEBUG in production)
3. Keep redaction enabled in production
4. Use structured fields when possible
5. Avoid logging large payloads or binary data
6. Test log output in development to verify redaction
7. Use JSON format for production log aggregation

# Security Considerations

The redactor provides defense-in-depth against accidental credential exposure:

1. Enabled by default (opt-out rather than opt-in)
2. Comprehensive pattern coverage for common credential formats
3. Regularly updated patterns for new token types
4. Works transparently without code changes
5. Cannot be bypassed without explicitly disabling

However, the redactor is not perfect:

1. Custom or unusual credential formats may not be detected
2. Heuristics may have false positives (redacting non-sensitive data)
3. Performance overhead may not be suitable for extreme high-volume logging

Always follow security best practices:

1. Avoid logging credentials if possible
2. Use structured logging to separate sensitive from non-sensitive data
3. Implement proper access controls for log files
4. Rotate and encrypt logs
5. Use secret management systems for credentials
*/
package logging
