// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package client provides HTTP client implementations for communicating with
// Itential Platform REST APIs.
//
// This package implements the Client interface with support for multiple
// authentication methods, TLS configuration, cookie management, and
// request/response handling.
//
// # Client Interface
//
// The Client interface defines the contract for HTTP operations:
//
//	type Client interface {
//	    Get(req *Request) (*Response, error)
//	    Post(req *Request) (*Response, error)
//	    Put(req *Request) (*Response, error)
//	    Patch(req *Request) (*Response, error)
//	    Delete(req *Request) (*Response, error)
//	    Trace(req *Request) (*Response, error)
//	}
//
// # HttpClient Implementation
//
// The primary implementation is HttpClient, which uses the standard library's
// net/http package with additional features:
//
//   - Basic Authentication
//   - OAuth2 Client Credentials flow
//   - Custom TLS configuration with certificate verification control
//   - Thread-safe cookie jar for session management
//   - Configurable request timeouts
//   - Context-aware request cancellation
//   - Request/response logging with optional suppression
//   - Automatic retry logic (future enhancement)
//
// # Creating a Client
//
// Create an HTTP client using a configuration profile:
//
//	profile := &config.Profile{
//	    Host:     "localhost",
//	    Port:     3000,
//	    UseTLS:   true,
//	    Verify:   false,
//	    Username: "admin@pronghorn",
//	    Password: "password",
//	    Timeout:  30,
//	}
//
//	client, err := client.NewHttpClient(profile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Authentication Methods
//
// The client supports two authentication methods:
//
// Basic Authentication (username/password):
//
//	profile.Username = "admin@pronghorn"
//	profile.Password = "password"
//
// OAuth2 Client Credentials:
//
//	profile.ClientID = "your-client-id"
//	profile.ClientSecret = "your-client-secret"
//
// OAuth2 is preferred when both are configured.
//
// # TLS Configuration
//
// Control TLS behavior with profile settings:
//
//	profile.UseTLS = true   // Use HTTPS
//	profile.Verify = false  // Skip certificate verification (development only)
//
// For production, always set Verify to true to ensure secure connections.
//
// # Making Requests
//
// Create requests using the Request struct:
//
//	req := &client.Request{
//	    Path:   "/automation-studio/projects",
//	    Params: map[string]string{"limit": "100"},
//	    Body:   []byte(`{"name": "my-project"}`),
//	}
//
//	resp, err := client.Post(req)
//	if err != nil {
//	    return err
//	}
//
// # Response Handling
//
// The Response struct contains:
//
//	type Response struct {
//	    StatusCode int                 // HTTP status code
//	    Status     string              // HTTP status text
//	    Headers    map[string]string   // Response headers
//	    Body       []byte              // Response body
//	}
//
// Check status codes and unmarshal response bodies:
//
//	if resp.StatusCode != http.StatusOK {
//	    return fmt.Errorf("unexpected status: %d", resp.StatusCode)
//	}
//
//	var data MyStruct
//	if err := json.Unmarshal(resp.Body, &data); err != nil {
//	    return err
//	}
//
// # Error Handling
//
// The client returns errors for:
//   - Network connectivity failures
//   - DNS resolution failures
//   - TLS handshake failures
//   - Connection timeouts
//   - Invalid URLs or parameters
//
// HTTP error status codes (4xx, 5xx) do not cause errors at the client level.
// The caller must check Response.StatusCode to handle HTTP-level errors.
//
// # Timeouts and Context
//
// Configure request timeouts via the profile:
//
//	profile.Timeout = 30  // 30 seconds (0 = no timeout)
//
// Timeouts apply to the entire request/response cycle including connection
// establishment, TLS handshake, request transmission, and response reading.
// Setting Timeout to 0 disables the timeout, allowing requests to run indefinitely.
//
// The client also supports context-based cancellation. Pass a context with
// cancellation or deadline to enable external control over request lifecycle:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	client := New(ctx, profile)
//
// # Cookie Management
//
// The client automatically manages cookies using a cookie jar. Session cookies
// and authentication tokens are preserved across requests to the same host.
//
// # Logging
//
// Request and response details are logged when the logger is configured:
//
//	logger.Debug("Request: %s %s", req.Method, req.URL)
//	logger.Debug("Response: %d %s", resp.StatusCode, resp.Status)
//
// Request bodies can be suppressed from logs by setting:
//
//	req.NoLog = true
//
// # Thread Safety
//
// HttpClient is safe for concurrent use from multiple goroutines. The underlying
// http.Client handles connection pooling and request serialization automatically.
// The cookie jar implementation uses sync.RWMutex to ensure thread-safe access
// to cookies across concurrent requests.
//
// Multiple goroutines can safely make requests using the same HttpClient instance:
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//	    wg.Add(1)
//	    go func() {
//	        defer wg.Done()
//	        req := NewRequest("/api/endpoint")
//	        resp, err := client.Get(req)
//	        // handle response
//	    }()
//	}
//	wg.Wait()
//
// # Connection Pooling
//
// The client uses the default http.Transport which maintains a pool of
// persistent connections. Connections are reused when possible to reduce
// latency and resource consumption. The transport automatically manages
// connection lifecycle, keepalives, and idle connection cleanup.
package client
