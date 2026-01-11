# Logging Reference

This document provides comprehensive information about the logging capabilities
in `ipctl`, including configuration options, sensitive data redaction, and best
practices for production deployments.

## Overview

`ipctl` includes structured logging with automatic sensitive data redaction to
protect credentials and secrets from accidental exposure in logs. The logging
system is built on [zerolog](https://github.com/rs/zerolog) and provides both
human-readable console output and machine-parseable JSON output.

## Key Features

- **Automatic sensitive data redaction** - Credentials, tokens, and secrets are
  automatically detected and redacted before output
- **Multiple log levels** - Control verbosity with standard log levels (TRACE,
  DEBUG, INFO, WARN, ERROR, FATAL)
- **Flexible output formats** - Choose between human-readable console format or
  structured JSON
- **Stream separation** - Logs automatically route to stdout (INFO/DEBUG/WARN)
  or stderr (ERROR/FATAL)
- **Timezone support** - Configure timestamp timezones for log entries
- **Enabled by default** - Security-first design with redaction on by default

## Configuration

Logging can be configured through environment variables or the configuration
file. Environment variables take precedence over file settings.

### Environment Variables

All logging configuration uses the `IPCTL_LOG_` prefix:

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `IPCTL_LOG_LEVEL` | TRACE, DEBUG, INFO, WARN, ERROR, FATAL, DISABLED | INFO | Minimum log level to output |
| `IPCTL_LOG_CONSOLE_JSON` | true, false | false | Enable JSON format for console output |
| `IPCTL_LOG_TIMESTAMP_TIMEZONE` | UTC, Local, or IANA timezone | UTC | Timezone for log timestamps |
| `IPCTL_LOG_REDACT_SENSITIVE_DATA` | true, false | true | Enable automatic sensitive data redaction |

### Configuration File

Logging settings are defined under the `[log]` section of the configuration
file (default: `~/.platform.d/config`):

```ini
[log]
level = info
console_json = false
timestamp_timezone = UTC
redact_sensitive_data = true
```

### Command Line

The `--verbose` flag enables console logging with DEBUG level output:

```bash
ipctl --verbose get projects
```

Without `--verbose`, logs are suppressed unless there's an error.

## Log Levels

Log levels control the verbosity of output. Messages below the configured level
are filtered out.

### TRACE

Extremely verbose output including function names, file names, and line numbers.
Only use for detailed debugging of specific issues.

```bash
export IPCTL_LOG_LEVEL=TRACE
ipctl --verbose get projects
```

Example output:
```
TRC github.com/itential/ipctl/internal/runners.ProjectRunner.Get.projects.go.142
```

### DEBUG

Detailed debugging information useful for troubleshooting. Includes API calls,
parameter values, and internal state.

```bash
export IPCTL_LOG_LEVEL=DEBUG
ipctl --verbose get projects
```

Example output:
```
DBG Fetching projects from server host=platform.example.com
```

### INFO (Default)

Important operational information such as successful operations and significant
state changes.

```bash
export IPCTL_LOG_LEVEL=INFO
ipctl --verbose get projects
```

Example output:
```
INF Retrieved 42 projects from server
```

### WARN

Warning messages for concerning but recoverable situations. The operation
continues but may require attention.

Example output:
```
WRN Rate limit approaching requests=950 limit=1000
```

### ERROR

Error conditions that occurred but allow the application to continue.

Example output:
```
ERR Failed to fetch projects error="connection timeout"
```

### FATAL

Critical errors that require immediate application shutdown. After logging,
the application calls `os.Exit(1)`.

Example output:
```
FTL Cannot connect to database error="authentication failed"
```

### DISABLED

Completely disable all logging output:

```bash
export IPCTL_LOG_LEVEL=DISABLED
```

## Output Formats

### Console Format (Default)

Human-readable output with timestamps, log levels, and optional colors:

```
2024-01-15T10:30:45Z INF Server started on port 8080
2024-01-15T10:30:46Z DBG Processing request user=john
2024-01-15T10:30:47Z WRN Slow query duration=5.2s
```

Disable colors when redirecting to a file:

```bash
ipctl --no-color --verbose get projects > output.log
```

### JSON Format

Structured JSON output suitable for log aggregation and analysis tools:

```bash
export IPCTL_LOG_CONSOLE_JSON=true
ipctl --verbose get projects
```

Example output:
```json
{"level":"info","time":"2024-01-15T10:30:45Z","message":"Server started on port 8080"}
{"level":"debug","time":"2024-01-15T10:30:46Z","user":"john","message":"Processing request"}
{"level":"warn","time":"2024-01-15T10:30:47Z","duration":"5.2s","message":"Slow query"}
```

JSON format is recommended for:
- Production environments
- Log aggregation systems (ELK, Splunk, etc.)
- Automated log parsing
- Cloud logging services

## Sensitive Data Redaction

The logging system automatically scans all output for sensitive information and
replaces it with `<REDACTED>` before writing to stdout or stderr. This protects
against accidental credential exposure in logs.

### Protected Data Types

The redactor detects and removes:

**Authentication & Authorization:**
- API keys (api_key=, api-key:, X-API-Key:)
- Bearer tokens (Authorization: Bearer ...)
- OAuth tokens (oauth_token, access_token, refresh_token)
- JWT tokens (eyJhbGc... format)
- Basic auth headers (Authorization: Basic ...)
- Session tokens and IDs

**Cloud Provider Credentials:**
- AWS access keys (AKIA..., ASIA..., etc.)
- AWS secret keys
- GitHub personal access tokens (ghp_...)
- GitHub OAuth tokens (gho_...)
- GitHub app tokens (ghu_..., ghs_...)

**Secrets & Passwords:**
- Passwords (password=, passwd=, pwd=)
- Generic secrets (secret=, client_secret=)

**Database & Infrastructure:**
- MongoDB connection strings (mongodb://user:pass@...)
- PostgreSQL connection strings (postgresql://user:pass@...)
- SSH private keys (-----BEGIN ... PRIVATE KEY-----)

### Redaction Examples

Before redaction:
```
DBG Connecting to API token=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
INF Using API key: sk_live_1234567890abcdefghijklmnop
DBG Database URI: mongodb://admin:password123@localhost:27017/db
```

After redaction:
```
DBG Connecting to API token=Bearer <REDACTED>
INF Using API key: <REDACTED>
DBG Database URI: <REDACTED>
```

### Disabling Redaction

Redaction is **enabled by default** for security. To disable (not recommended
for production):

```bash
export IPCTL_LOG_REDACT_SENSITIVE_DATA=false
```

Or in the configuration file:

```ini
[log]
redact_sensitive_data = false
```

**Warning:** Only disable redaction in isolated development environments where
logs are not persisted or shared. Production logs should always have redaction
enabled.

## Stream Routing

Log messages are automatically routed to the appropriate output stream based on
their severity level:

**stdout** (standard output):
- TRACE
- DEBUG
- INFO
- WARN

**stderr** (standard error):
- ERROR
- FATAL

This separation allows easy filtering in shell pipelines:

```bash
# Capture only normal logs
ipctl --verbose get projects > normal.log

# Capture only errors
ipctl --verbose get projects 2> errors.log

# Separate streams
ipctl --verbose get projects > normal.log 2> errors.log

# Combine streams
ipctl --verbose get projects &> all.log
```

## Timezone Configuration

Timestamps in logs can be configured to use any timezone:

```bash
# UTC (default)
export IPCTL_LOG_TIMESTAMP_TIMEZONE=UTC

# Local system timezone
export IPCTL_LOG_TIMESTAMP_TIMEZONE=Local

# Specific IANA timezone
export IPCTL_LOG_TIMESTAMP_TIMEZONE=America/New_York
export IPCTL_LOG_TIMESTAMP_TIMEZONE=Europe/London
export IPCTL_LOG_TIMESTAMP_TIMEZONE=Asia/Tokyo
```

Configuration file:

```ini
[log]
timestamp_timezone = America/New_York
```

All timestamps use RFC3339 format: `2024-01-15T10:30:45-05:00`

## Best Practices

### Development

- Use `--verbose` flag for detailed output during development
- Use DEBUG level for troubleshooting: `IPCTL_LOG_LEVEL=DEBUG`
- Use console format (default) for readability
- Disable colors when redirecting: `--no-color`

### Production

- Use INFO or WARN level to reduce noise
- Enable JSON format for log aggregation: `IPCTL_LOG_CONSOLE_JSON=true`
- **Always keep redaction enabled**: `IPCTL_LOG_REDACT_SENSITIVE_DATA=true`
- Configure appropriate timezone for your region
- Route stdout and stderr to separate log files
- Set up log rotation to manage disk space

### CI/CD Pipelines

- Use INFO level for build logs
- Enable JSON format for parsing: `IPCTL_LOG_CONSOLE_JSON=true`
- Disable colors: `--no-color`
- Keep redaction enabled to protect secrets in CI logs

### Debugging

- Start with DEBUG level: `IPCTL_LOG_LEVEL=DEBUG`
- If more detail needed, use TRACE: `IPCTL_LOG_LEVEL=TRACE`
- Use console format for readability
- Grep for specific patterns: `ipctl --verbose get projects 2>&1 | grep ERROR`

## Common Scenarios

### Troubleshooting Connection Issues

```bash
export IPCTL_LOG_LEVEL=DEBUG
ipctl --verbose --profile prod get projects
```

### Capturing All Logs to File

```bash
ipctl --verbose --no-color get projects &> ipctl.log
```

### JSON Logs for ELK Stack

```bash
export IPCTL_LOG_LEVEL=INFO
export IPCTL_LOG_CONSOLE_JSON=true
ipctl --verbose get projects | tee -a ipctl-json.log
```

### Separating Normal and Error Logs

```bash
ipctl --verbose get projects \
  > >(tee -a normal.log) \
  2> >(tee -a errors.log >&2)
```

### Testing Redaction

Create a test script to verify redaction works:

```bash
# This should show <REDACTED> instead of the actual token
export IPCTL_LOG_LEVEL=DEBUG
export TEST_TOKEN="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"
echo "Test token: $TEST_TOKEN" | ipctl --verbose get projects
```

## Performance Considerations

The redaction system adds minimal overhead to log output:

- **Clean text** (no sensitive data): ~29 microseconds per message
- **Text with sensitive data**: ~34 microseconds per message
- **Check if redaction needed**: ~809 nanoseconds per message

For most applications, this overhead is negligible. However, if you're logging
at extremely high volumes (>100,000 messages/second), consider:

1. Using JSON format (faster than console formatting)
2. Raising the log level to INFO or WARN
3. Benchmarking with your specific workload

## Security Considerations

### Defense in Depth

The redactor provides an additional layer of security but should not be the
only defense:

1. **Avoid logging credentials** - Don't log sensitive data in the first place
2. **Use structured logging** - Separate sensitive from non-sensitive data
3. **Secure log files** - Set appropriate permissions: `chmod 600`
4. **Rotate logs** - Don't keep logs indefinitely
5. **Encrypt logs** - Use encryption for log storage and transmission
6. **Audit access** - Monitor who accesses log files

### Limitations

The redactor uses heuristic patterns and may not catch:

- Custom or proprietary credential formats
- Obfuscated credentials
- Credentials split across multiple log lines
- Binary data encoded in unusual ways

### False Positives

The redactor may occasionally redact non-sensitive data that matches patterns:

- Long alphanumeric strings that look like tokens
- Encoded data that resembles credentials
- Test data with credential-like format

This is acceptable as the cost of false positives (losing some context) is much
lower than false negatives (exposing credentials).

## Troubleshooting

### Logs Not Appearing

**Problem:** Running `ipctl` commands shows no log output.

**Solution:** Use the `--verbose` flag:
```bash
ipctl --verbose get projects
```

### Too Much Output

**Problem:** Logs are too verbose and cluttering the terminal.

**Solution:** Raise the log level:
```bash
export IPCTL_LOG_LEVEL=WARN
ipctl --verbose get projects
```

### Can't See Credentials for Debugging

**Problem:** Need to see actual credentials for debugging (development only).

**Solution:** Temporarily disable redaction:
```bash
export IPCTL_LOG_REDACT_SENSITIVE_DATA=false
ipctl --verbose get projects
```

**Important:** Re-enable redaction when done and never use this in production.

### Timestamps in Wrong Timezone

**Problem:** Log timestamps don't match your local time.

**Solution:** Configure the timezone:
```bash
export IPCTL_LOG_TIMESTAMP_TIMEZONE=Local
```

Or specify your timezone:
```bash
export IPCTL_LOG_TIMESTAMP_TIMEZONE=America/New_York
```

### Colors Interfering with Log Files

**Problem:** Log files contain color codes when redirected.

**Solution:** Use the `--no-color` flag:
```bash
ipctl --no-color --verbose get projects > output.log
```

## Examples

### Basic Logging

```bash
# Default INFO level with console format
ipctl --verbose get projects

# DEBUG level for troubleshooting
export IPCTL_LOG_LEVEL=DEBUG
ipctl --verbose get projects

# JSON output for machines
export IPCTL_LOG_CONSOLE_JSON=true
ipctl --verbose get projects
```

### Production Configuration

```bash
# Production environment variables
export IPCTL_LOG_LEVEL=INFO
export IPCTL_LOG_CONSOLE_JSON=true
export IPCTL_LOG_REDACT_SENSITIVE_DATA=true
export IPCTL_LOG_TIMESTAMP_TIMEZONE=UTC

# Run with logging
ipctl --verbose get projects | tee -a /var/log/ipctl/ipctl.log
```

### CI/CD Pipeline

```bash
#!/bin/bash
# CI pipeline logging configuration

export IPCTL_LOG_LEVEL=INFO
export IPCTL_LOG_CONSOLE_JSON=true
export IPCTL_LOG_REDACT_SENSITIVE_DATA=true

# Run command and parse JSON logs
ipctl --verbose --no-color deploy | \
  jq -r 'select(.level=="error") | .message'
```

## Related Documentation

- [Configuration Reference](configuration-reference.md) - Complete configuration options
- [Command Quick Reference](commands-quick-reference.md) - Available commands
- [Running from Source](running-from-source.md) - Development setup

## Additional Resources

- [zerolog Documentation](https://github.com/rs/zerolog) - Underlying logging library
- [RFC3339 Timestamp Format](https://www.rfc-editor.org/rfc/rfc3339) - Timestamp standard
- [IANA Timezone Database](https://www.iana.org/time-zones) - Valid timezone names
