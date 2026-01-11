// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logging

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedactor(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "enabled redactor",
			enabled: true,
		},
		{
			name:    "disabled redactor",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedactor(tt.enabled)
			assert.NotNil(t, r)
			assert.Equal(t, tt.enabled, r.enabled)
		})
	}
}

func TestRedactor_Redact_BearerTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bearer token in authorization header",
			input:    "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected: "Authorization: Bearer <REDACTED>",
		},
		{
			name:     "bearer token lowercase",
			input:    "authorization: bearer abc123def456ghi789",
			expected: "authorization: bearer <REDACTED>",
		},
		{
			name:     "bearer token in log message",
			input:    "Using Bearer token12345678901234567890 for authentication",
			expected: "Using Bearer <REDACTED> for authentication",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_APIKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "api key with equals",
			input:    "api_key=sk_test_1234567890abcdefghijklmnop",
			expected: "api_key=<REDACTED>",
		},
		{
			name:     "api key with colon",
			input:    "API-KEY: 1234567890abcdefghijklmnop",
			expected: "API-KEY: <REDACTED>",
		},
		{
			name:     "x-api-key header",
			input:    "X-API-Key: abcdef1234567890",
			expected: "X-API-Key: <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_JWTTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid jwt token",
			input:    "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expected: "Token: <REDACTED>",
		},
		{
			name:     "jwt in json",
			input:    `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abcdef"}`,
			expected: `{"token":"<REDACTED>"}`,
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_OAuthTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "oauth token",
			input:    "oauth_token=ya29.a0AfH6SMBx1234567890",
			expected: "oauth_token=<REDACTED>",
		},
		{
			name:     "access token",
			input:    "access_token: ghu_1234567890abcdefghij",
			expected: "access_token: <REDACTED>",
		},
		{
			name:     "refresh token",
			input:    "refresh-token=1234567890abcdefghijklmnop",
			expected: "refresh-token=<REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_Passwords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "password with equals",
			input:    "password=MySecretP@ssw0rd",
			expected: "password=<REDACTED>",
		},
		{
			name:     "password with colon",
			input:    "PASSWORD: SuperSecret123!",
			expected: "PASSWORD: <REDACTED>",
		},
		{
			name:     "passwd field",
			input:    "passwd=hunter2",
			expected: "passwd=<REDACTED>",
		},
		{
			name:     "pwd field",
			input:    "pwd: password123",
			expected: "pwd: <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_AWSCredentials(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "aws access key",
			input:    "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			expected: "AWS_ACCESS_KEY_ID=<REDACTED>",
		},
		{
			name:     "aws secret key",
			input:    "aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			expected: "aws_secret_access_key=<REDACTED>",
		},
		{
			name:     "aws in log message",
			input:    "Using credentials AKIAIOSFODNN7EXAMPLE for AWS",
			expected: "Using credentials <REDACTED> for AWS",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_GitHubTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "github personal access token",
			input:    "GITHUB_TOKEN=ghp_1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "GITHUB_TOKEN=<REDACTED>",
		},
		{
			name:     "github oauth token",
			input:    "token: gho_1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "token: <REDACTED>",
		},
		{
			name:     "github app token",
			input:    "Using ghu_1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "Using <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_Secrets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "generic secret",
			input:    "secret=abcdef1234567890",
			expected: "secret=<REDACTED>",
		},
		{
			name:     "client secret",
			input:    "client_secret: xyz789abc123def456",
			expected: "client_secret: <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_SessionTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "session token",
			input:    "session_token=abc123def456ghi789jkl012",
			expected: "session_token=<REDACTED>",
		},
		{
			name:     "session id",
			input:    "session-id: xyz789abc123def456ghi",
			expected: "session-id: <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_DatabaseURIs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mongodb uri",
			input:    "mongodb://user:password@localhost:27017/database",
			expected: "<REDACTED>",
		},
		{
			name:     "mongodb srv uri",
			input:    "mongodb+srv://admin:secret123@cluster.mongodb.net/mydb",
			expected: "<REDACTED>",
		},
		{
			name:     "postgres uri",
			input:    "postgresql://user:pass@localhost:5432/db",
			expected: "<REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_SSHKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "rsa private key",
			input: `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890
-----END RSA PRIVATE KEY-----`,
			expected: "<REDACTED>",
		},
		{
			name: "openssh private key",
			input: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAA
-----END OPENSSH PRIVATE KEY-----`,
			expected: "<REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_AuthorizationHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "authorization header basic",
			input:    "Authorization: Basic dXNlcjpwYXNzd29yZA==",
			expected: "Authorization: Basic <REDACTED>",
		},
		{
			name:     "authorization header basic lowercase",
			input:    "authorization: basic dXNlcjpwYXNzd29yZA==",
			expected: "authorization: basic <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_GenericTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "generic token field",
			input:    "token=abcdefghijklmnopqrstuvwxyz1234567890",
			expected: "token=<REDACTED>",
		},
		{
			name:     "token with colon",
			input:    "TOKEN: xyz123abc456def789ghi012jkl345mno678",
			expected: "TOKEN: <REDACTED>",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_MultipleSecrets(t *testing.T) {
	input := `
	API Key: sk_test_1234567890abcdef
	Password: MySecretPass123
	Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abc
	AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
	mongodb://user:pass@localhost/db
	`

	expected := `
	API Key: <REDACTED>
	Password: <REDACTED>
	Bearer <REDACTED>
	AWS_ACCESS_KEY_ID=<REDACTED>
	<REDACTED>
	`

	r := NewRedactor(true)
	result := r.Redact(input)
	assert.Equal(t, expected, result)
}

func TestRedactor_Redact_Disabled(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "bearer token should not be redacted",
			input: "Authorization: Bearer secret123token",
		},
		{
			name:  "password should not be redacted",
			input: "password=MySecretPassword",
		},
		{
			name:  "api key should not be redacted",
			input: "api_key=sk_test_1234567890",
		},
	}

	r := NewRedactor(false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.input, result, "redactor should not modify input when disabled")
		})
	}
}

func TestRedactor_RedactBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "redact bytes with api key",
			input:    []byte("api_key=sk_test_1234567890abcdef"),
			expected: []byte("api_key=<REDACTED>"),
		},
		{
			name:     "disabled redactor returns original",
			input:    []byte("password=secret"),
			expected: []byte("password=secret"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var enabled bool
			if strings.Contains(tt.name, "disabled") {
				enabled = false
			} else {
				enabled = true
			}
			r := NewRedactor(enabled)
			result := r.RedactBytes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_ShouldRedact(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		enabled  bool
	}{
		{
			name:     "should redact with password",
			input:    "password=secret",
			expected: true,
			enabled:  true,
		},
		{
			name:     "should redact with token",
			input:    "token=abcdefghijklmnopqrstuvwxyz1234567890",
			expected: true,
			enabled:  true,
		},
		{
			name:     "should not redact clean text",
			input:    "Hello, world! This is a normal log message.",
			expected: false,
			enabled:  true,
		},
		{
			name:     "should not redact when disabled",
			input:    "password=secret",
			expected: false,
			enabled:  false,
		},
		{
			name:     "should not redact short token",
			input:    "token=abc",
			expected: false,
			enabled:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedactor(tt.enabled)
			result := r.ShouldRedact(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactor_Redact_SafeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "normal log message",
			input: "Application started successfully",
		},
		{
			name:  "message with numbers",
			input: "Processing 1234 items in queue",
		},
		{
			name:  "message with email",
			input: "User logged in: user@example.com",
		},
		{
			name:  "message with url",
			input: "Connecting to https://api.example.com/v1/users",
		},
	}

	r := NewRedactor(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Redact(tt.input)
			assert.Equal(t, tt.input, result, "safe text should not be modified")
		})
	}
}

func BenchmarkRedactor_Redact(b *testing.B) {
	r := NewRedactor(true)
	text := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abc password=secret123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Redact(text)
	}
}

func BenchmarkRedactor_ShouldRedact(b *testing.B) {
	r := NewRedactor(true)
	text := "Authorization: Bearer token123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ShouldRedact(text)
	}
}

func BenchmarkRedactor_Redact_CleanText(b *testing.B) {
	r := NewRedactor(true)
	text := "This is a normal log message without any sensitive data"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Redact(text)
	}
}
