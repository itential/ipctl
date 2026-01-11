// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

// RequestOption is a functional option for configuring a Request.
//
// RequestOption functions modify a Request instance during construction,
// allowing for flexible and composable request configuration using the
// functional options pattern.
type RequestOption func(r *Request)

// Request represents an HTTP request to be sent to the Itential Platform.
//
// Request instances are created using NewRequest and configured with
// functional options. The Path field is required; all other fields are optional.
//
// Example:
//
//	req := NewRequest("/api/v1/users",
//	    WithParams(map[string]string{"page": "1"}),
//	    WithBody([]byte(`{"name": "John"}`)),
//	)
type Request struct {
	// Path is the request path appended to the base URL.
	// This should start with a forward slash (e.g., "/api/v1/users").
	Path string

	// Params are query parameters appended to the URL as a query string.
	// Example: {"page": "1", "limit": "10"} becomes "?page=1&limit=10"
	Params map[string]string

	// Headers are custom HTTP headers to send with the request.
	// Note: Default headers (Content-Type, Accept) are set automatically.
	Headers map[string]string

	// Body is the request body sent to the server.
	// Typically contains JSON-encoded data for POST, PUT, and PATCH requests.
	Body []byte

	// NoLog suppresses request body logging when true.
	// Use this for requests containing sensitive data like passwords or tokens.
	NoLog bool
}

// NewRequest creates a new HTTP request with the specified path and options.
//
// The path parameter is required and specifies the API endpoint path.
// Optional parameters can be configured using functional options like
// WithParams, WithBody, WithHeaders, and WithNoLog.
//
// Example:
//
//	// Simple GET request
//	req := NewRequest("/api/v1/users")
//
//	// POST request with body and query parameters
//	req := NewRequest("/api/v1/users",
//	    WithParams(map[string]string{"notify": "true"}),
//	    WithBody([]byte(`{"name": "Alice", "email": "alice@example.com"}`)),
//	)
//
//	// Request with sensitive data (no logging)
//	req := NewRequest("/login",
//	    WithBody([]byte(`{"username": "admin", "password": "secret"}`)),
//	    WithNoLog(true),
//	)
func NewRequest(path string, opts ...RequestOption) *Request {
	req := &Request{Path: path}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

// WithHeaders returns a RequestOption that sets custom HTTP headers.
//
// The provided headers will be added to the request. Note that some headers
// like Content-Type and Accept are set automatically by the HTTP client.
//
// Example:
//
//	req := NewRequest("/api/v1/data",
//	    WithHeaders(map[string]string{
//	        "X-Custom-Header": "value",
//	        "Authorization": "Bearer token",
//	    }),
//	)
func WithHeaders(v map[string]string) RequestOption {
	return func(r *Request) {
		r.Headers = v
	}
}

// WithParams returns a RequestOption that sets query string parameters.
//
// The parameters are appended to the URL as a query string. Parameter keys
// and values are automatically URL-encoded.
//
// Example:
//
//	req := NewRequest("/api/v1/users",
//	    WithParams(map[string]string{
//	        "page": "1",
//	        "limit": "50",
//	        "status": "active",
//	    }),
//	)
//	// Results in: /api/v1/users?page=1&limit=50&status=active
func WithParams(v map[string]string) RequestOption {
	return func(r *Request) {
		r.Params = v
	}
}

// WithBody returns a RequestOption that sets the request body.
//
// The body should typically contain JSON-encoded data for POST, PUT, and
// PATCH requests. For GET and DELETE requests, the body is usually empty.
//
// Example:
//
//	data := map[string]interface{}{
//	    "name": "New Project",
//	    "description": "A project description",
//	}
//	body, _ := json.Marshal(data)
//	req := NewRequest("/api/v1/projects", WithBody(body))
func WithBody(v []byte) RequestOption {
	return func(r *Request) {
		r.Body = v
	}
}

// WithNoLog returns a RequestOption that disables request body logging.
//
// When set to true, the request body will not be logged even when debug
// logging is enabled. This is important for requests containing sensitive
// information like passwords, API keys, or tokens.
//
// Example:
//
//	req := NewRequest("/login",
//	    WithBody([]byte(`{"username": "admin", "password": "secret123"}`)),
//	    WithNoLog(true), // Password won't appear in logs
//	)
func WithNoLog(v bool) RequestOption {
	return func(r *Request) {
		r.NoLog = v
	}
}
