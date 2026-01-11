// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"regexp"
	"strings"
)

// redactionPatterns holds compiled regex patterns for detecting sensitive information.
// Each pattern is designed to match common formats of credentials and sensitive data.
var redactionPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	// API Keys and Bearer tokens
	{
		name:    "bearer_token",
		pattern: regexp.MustCompile(`(?i)(bearer\s+)[a-zA-Z0-9\-._~+/]+=*`),
	},
	{
		name:    "api_key",
		pattern: regexp.MustCompile(`(?i)(api[_\-\s]?key\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},
	{
		name:    "api_key_header",
		pattern: regexp.MustCompile(`(?i)(x-api-key\s*:\s*)[a-zA-Z0-9\-._~+/]{16,}`),
	},

	// JWT tokens (three base64 segments separated by dots)
	{
		name:    "jwt_token",
		pattern: regexp.MustCompile(`eyJ[a-zA-Z0-9\-._~+/]*\.eyJ[a-zA-Z0-9\-._~+/]*\.[a-zA-Z0-9\-._~+/]*`),
	},

	// OAuth tokens
	{
		name:    "oauth_token",
		pattern: regexp.MustCompile(`(?i)(oauth[_-]?token\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},
	{
		name:    "access_token",
		pattern: regexp.MustCompile(`(?i)(access[_-]?token\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},
	{
		name:    "refresh_token",
		pattern: regexp.MustCompile(`(?i)(refresh[_-]?token\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},

	// Passwords (minimum 6 characters to catch common passwords)
	{
		name:    "password",
		pattern: regexp.MustCompile(`(?i)(password\s*[:=]\s*['"]?)[^\s'"]{6,}['"]?`),
	},
	{
		name:    "passwd",
		pattern: regexp.MustCompile(`(?i)(passwd\s*[:=]\s*['"]?)[^\s'"]{6,}['"]?`),
	},
	{
		name:    "pwd",
		pattern: regexp.MustCompile(`(?i)(pwd\s*[:=]\s*['"]?)[^\s'"]{6,}['"]?`),
	},

	// AWS credentials
	{
		name:    "aws_access_key",
		pattern: regexp.MustCompile(`(?:A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`),
	},
	{
		name:    "aws_secret_key",
		pattern: regexp.MustCompile(`(?i)(aws[_-]?secret[_-]?access[_-]?key\s*[:=]\s*['"]?)[a-zA-Z0-9/+=]{40}['"]?`),
	},

	// GitHub tokens
	{
		name:    "github_token",
		pattern: regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
	},
	{
		name:    "github_oauth",
		pattern: regexp.MustCompile(`gho_[a-zA-Z0-9]{36}`),
	},
	{
		name:    "github_app_token",
		pattern: regexp.MustCompile(`(?:ghu|ghs)_[a-zA-Z0-9]{36}`),
	},

	// Generic secrets
	{
		name:    "secret",
		pattern: regexp.MustCompile(`(?i)(secret\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},
	{
		name:    "client_secret",
		pattern: regexp.MustCompile(`(?i)(client[_-]?secret\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},

	// Session tokens
	{
		name:    "session_token",
		pattern: regexp.MustCompile(`(?i)(session[_-]?token\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},
	{
		name:    "session_id",
		pattern: regexp.MustCompile(`(?i)(session[_-]?id\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{16,}['"]?`),
	},

	// Database connection strings
	{
		name:    "mongodb_uri",
		pattern: regexp.MustCompile(`mongodb(?:\+srv)?://[^:]+:[^@]+@[^\s]+`),
	},
	{
		name:    "postgres_uri",
		pattern: regexp.MustCompile(`postgres(?:ql)?://[^:]+:[^@]+@[^\s]+`),
	},

	// SSH private keys
	{
		name:    "ssh_private_key",
		pattern: regexp.MustCompile(`-----BEGIN\s+(?:RSA|DSA|EC|OPENSSH)\s+PRIVATE\s+KEY-----[\s\S]*?-----END\s+(?:RSA|DSA|EC|OPENSSH)\s+PRIVATE\s+KEY-----`),
	},

	// Authorization headers - Basic auth specifically
	{
		name:    "authorization_basic",
		pattern: regexp.MustCompile(`(?i)(authorization\s*:\s*basic\s+)[a-zA-Z0-9+/=]+`),
	},

	// Generic token patterns (last resort)
	{
		name:    "generic_token",
		pattern: regexp.MustCompile(`(?i)(token\s*[:=]\s*['"]?)[a-zA-Z0-9\-._~+/]{20,}['"]?`),
	},
}

// Redactor handles scanning and redacting sensitive information from log output.
type Redactor struct {
	enabled bool
}

// NewRedactor creates a new Redactor instance.
func NewRedactor(enabled bool) *Redactor {
	return &Redactor{
		enabled: enabled,
	}
}

// Redact scans the input text for sensitive information and replaces it with <REDACTED>.
// It uses heuristics-based pattern matching to identify credentials, tokens, and other secrets.
// If redaction is disabled, returns the original text unchanged.
func (r *Redactor) Redact(text string) string {
	if !r.enabled {
		return text
	}

	result := text

	// Apply each redaction pattern
	for _, rp := range redactionPatterns {
		result = rp.pattern.ReplaceAllStringFunc(result, func(match string) string {
			// For patterns with capture groups (prefix), preserve the prefix
			// and only redact the sensitive portion
			if rp.pattern.NumSubexp() > 0 {
				submatches := rp.pattern.FindStringSubmatch(match)
				if len(submatches) > 1 {
					// Keep the prefix (capture group 1) and redact the rest
					prefix := submatches[1]
					return prefix + "<REDACTED>"
				}
			}
			// For patterns without capture groups, redact the entire match
			return "<REDACTED>"
		})
	}

	return result
}

// RedactBytes is a convenience method that wraps Redact for byte slices.
// It converts the input to a string, performs redaction, and returns the result as bytes.
func (r *Redactor) RedactBytes(data []byte) []byte {
	if !r.enabled {
		return data
	}

	text := string(data)
	redacted := r.Redact(text)
	return []byte(redacted)
}

// ShouldRedact checks if a specific string contains patterns that would trigger redaction.
// This can be used for performance optimization to skip redaction on clean text.
func (r *Redactor) ShouldRedact(text string) bool {
	if !r.enabled {
		return false
	}

	// Quick checks for common indicators before running expensive regex
	lowerText := strings.ToLower(text)
	indicators := []string{
		"token", "password", "secret", "api", "bearer", "authorization",
		"oauth", "session", "key", "jwt", "aws", "github", "mongodb",
		"postgres", "mysql", "-----begin", "passwd", "pwd",
	}

	hasIndicator := false
	for _, indicator := range indicators {
		if strings.Contains(lowerText, indicator) {
			hasIndicator = true
			break
		}
	}

	if !hasIndicator {
		return false
	}

	// Check if any pattern matches
	for _, rp := range redactionPatterns {
		if rp.pattern.MatchString(text) {
			return true
		}
	}

	return false
}
