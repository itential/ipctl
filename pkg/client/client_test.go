// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/itential/ipctl/pkg/config"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Profile{
		Host:         "test.example.com",
		Port:         8080,
		Username:     "testuser",
		Password:     "testpass",
		UseTLS:       true,
		Verify:       false,
		ClientID:     "testclient",
		ClientSecret: "testsecret",
	}

	client := New(ctx, cfg)

	if client.Host != cfg.Host {
		t.Errorf("Expected Host %s, got %s", cfg.Host, client.Host)
	}
	if client.Port != cfg.Port {
		t.Errorf("Expected Port %d, got %d", cfg.Port, client.Port)
	}
	if client.Username != cfg.Username {
		t.Errorf("Expected Username %s, got %s", cfg.Username, client.Username)
	}
	if client.Password != cfg.Password {
		t.Errorf("Expected Password %s, got %s", cfg.Password, client.Password)
	}
	if client.UseTls != cfg.UseTLS {
		t.Errorf("Expected UseTls %t, got %t", cfg.UseTLS, client.UseTls)
	}
	if client.Verify != cfg.Verify {
		t.Errorf("Expected Verify %t, got %t", cfg.Verify, client.Verify)
	}
	if client.ClientId != cfg.ClientID {
		t.Errorf("Expected ClientId %s, got %s", cfg.ClientID, client.ClientId)
	}
	if client.ClientSecret != cfg.ClientSecret {
		t.Errorf("Expected ClientSecret %s, got %s", cfg.ClientSecret, client.ClientSecret)
	}
	if client.context != ctx {
		t.Error("Expected context to be set")
	}
	if client.jar == nil {
		t.Error("Expected jar to be initialized")
	}
	if client.authenticated {
		t.Error("Expected authenticated to be false")
	}
}

func TestIsAuthenticated(t *testing.T) {
	client := &HttpClient{authenticated: false}
	if client.IsAuthenticated() {
		t.Error("Expected IsAuthenticated to return false")
	}

	client.authenticated = true
	if !client.IsAuthenticated() {
		t.Error("Expected IsAuthenticated to return true")
	}
}

func TestSetScheme(t *testing.T) {
	testCases := []struct {
		name     string
		useTls   bool
		expected string
	}{
		{
			name:     "HTTPS when UseTls is true",
			useTls:   true,
			expected: ProtocolHttps,
		},
		{
			name:     "HTTP when UseTls is false",
			useTls:   false,
			expected: ProtocolHttp,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &HttpClient{UseTls: tc.useTls}
			result := client.setScheme()
			if result != tc.expected {
				t.Errorf("Expected scheme %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestSetPort(t *testing.T) {
	testCases := []struct {
		name     string
		port     int
		useTls   bool
		expected int
	}{
		{
			name:     "Use configured port",
			port:     8080,
			useTls:   false,
			expected: 8080,
		},
		{
			name:     "Auto-set HTTPS port",
			port:     0,
			useTls:   true,
			expected: 443,
		},
		{
			name:     "Auto-set HTTP port",
			port:     0,
			useTls:   false,
			expected: 80,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &HttpClient{Port: tc.port, UseTls: tc.useTls}
			result := client.setPort()
			if result != tc.expected {
				t.Errorf("Expected port %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestNewUrl(t *testing.T) {
	client := &HttpClient{}

	testCases := []struct {
		name        string
		scheme      string
		host        string
		path        string
		params      map[string]string
		expectedUrl string
	}{
		{
			name:        "Simple URL without params",
			scheme:      "https",
			host:        "example.com:443",
			path:        "/api/v1/test",
			params:      nil,
			expectedUrl: "https://example.com:443/api/v1/test",
		},
		{
			name:        "URL with query parameters",
			scheme:      "http",
			host:        "localhost:8080",
			path:        "/api/data",
			params:      map[string]string{"page": "1", "limit": "10"},
			expectedUrl: "http://localhost:8080/api/data?limit=10&page=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.newUrl(tc.scheme, tc.host, tc.path, tc.params)
			if result.String() != tc.expectedUrl {
				t.Errorf("Expected URL %s, got %s", tc.expectedUrl, result.String())
			}
		})
	}
}

func TestNewHttpRequest(t *testing.T) {
	ctx := context.Background()
	client := &HttpClient{context: ctx}

	testCases := []struct {
		name    string
		method  string
		url     string
		body    []byte
		wantErr bool
	}{
		{
			name:    "Valid GET request",
			method:  "GET",
			url:     "https://example.com/api",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "Valid POST request with body",
			method:  "POST",
			url:     "https://example.com/api",
			body:    []byte(`{"test": "data"}`),
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			method:  "GET",
			url:     "://invalid-url",
			body:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := client.newHttpRequest(tc.method, tc.url, tc.body)

			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if req.Method != tc.method {
				t.Errorf("Expected method %s, got %s", tc.method, req.Method)
			}

			if req.Header.Get("Content-Type") != "application/json" {
				t.Error("Expected Content-Type header to be application/json")
			}

			if req.Header.Get("Accept") != "application/json" {
				t.Error("Expected Accept header to be application/json")
			}
		})
	}
}

func TestNewResponse(t *testing.T) {
	client := &HttpClient{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Test-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "test"}`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make test request: %v", err)
	}
	defer resp.Body.Close()

	response, err := client.newResponse(resp, "GET", server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if response.Method != "GET" {
		t.Errorf("Expected method GET, got %s", response.Method)
	}

	if response.Url != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, response.Url)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if string(response.Body) != `{"message": "test"}` {
		t.Errorf("Expected body %s, got %s", `{"message": "test"}`, string(response.Body))
	}

	if response.Headers["Test-Header"] != "test-value" {
		t.Errorf("Expected header Test-Header to be test-value, got %s", response.Headers["Test-Header"])
	}
}

func TestHttpClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"method": "` + r.Method + `"}`))
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	testCases := []struct {
		name           string
		method         func(*Request) (*Response, error)
		expectedMethod string
	}{
		{"GET", client.Get, "GET"},
		{"POST", client.Post, "POST"},
		{"PUT", client.Put, "PUT"},
		{"DELETE", client.Delete, "DELETE"},
		{"PATCH", client.Patch, "PATCH"},
		{"TRACE", client.Trace, "TRACE"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := NewRequest("/test")
			resp, err := tc.method(req)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}

			expectedBody := `{"method": "` + tc.expectedMethod + `"}`
			if string(resp.Body) != expectedBody {
				t.Errorf("Expected body %s, got %s", expectedBody, string(resp.Body))
			}
		})
	}
}

func TestAuthenticateUsingOAuth(t *testing.T) {
	client := &HttpClient{
		ClientId:     "test-client",
		ClientSecret: "test-secret",
		context:      context.Background(),
	}

	httpClient := &http.Client{}
	scheme := "https"
	remoteHost := "example.com:443"

	result, err := client.authenticateUsingOAuth(httpClient, scheme, remoteHost)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected OAuth client to be returned")
	}
}

func TestAuthenticateUsingOAuthMissingCredentials(t *testing.T) {
	testCases := []struct {
		name         string
		clientId     string
		clientSecret string
	}{
		{"Missing ClientId", "", "secret"},
		{"Missing ClientSecret", "client", ""},
		{"Missing both", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &HttpClient{
				ClientId:     tc.clientId,
				ClientSecret: tc.clientSecret,
				context:      context.Background(),
			}

			httpClient := &http.Client{}
			scheme := "https"
			remoteHost := "example.com:443"

			_, err := client.authenticateUsingOAuth(httpClient, scheme, remoteHost)
			if err == nil {
				t.Error("Expected error for missing credentials")
			}
		})
	}
}

// TestNewWithTimeout verifies that timeout configuration is properly applied
func TestNewWithTimeout(t *testing.T) {
	testCases := []struct {
		name            string
		timeout         int
		expectedTimeout int
	}{
		{
			name:            "Timeout set to 30 seconds",
			timeout:         30,
			expectedTimeout: 30,
		},
		{
			name:            "Timeout set to 60 seconds",
			timeout:         60,
			expectedTimeout: 60,
		},
		{
			name:            "Timeout set to 0 (no timeout)",
			timeout:         0,
			expectedTimeout: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := &config.Profile{
				Host:    "test.example.com",
				Port:    8080,
				Timeout: tc.timeout,
			}

			client := New(ctx, cfg)

			if client.Timeout != tc.expectedTimeout {
				t.Errorf("Expected Timeout %d, got %d", tc.expectedTimeout, client.Timeout)
			}
		})
	}
}

// TestHttpClientTimeout verifies that HTTP client respects timeout configuration
func TestHttpClientTimeout(t *testing.T) {
	// Create a server that delays response beyond timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		Timeout:       1, // 1 second timeout
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	req := NewRequest("/test")
	_, err := client.Get(req)

	if err == nil {
		t.Error("Expected timeout error, but got none")
	}
}

// TestHttpClientNoTimeout verifies that requests work when timeout is 0
func TestHttpClientNoTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		Timeout:       0, // No timeout
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	req := NewRequest("/test")
	resp, err := client.Get(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestHttpClientWithContextCancellation verifies context cancellation behavior
func TestHttpClientWithContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	// Create context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		context:       ctx,
		jar:           NewJar(),
		authenticated: true,
	}

	// Cancel context immediately
	cancel()

	req := NewRequest("/test")
	_, err := client.Get(req)

	if err == nil {
		t.Error("Expected error due to context cancellation, but got none")
	}
}

// TestHttpClientTLSConfiguration verifies TLS configuration behavior
func TestHttpClientTLSConfiguration(t *testing.T) {
	testCases := []struct {
		name   string
		useTls bool
		verify bool
	}{
		{
			name:   "TLS enabled with verification",
			useTls: true,
			verify: true,
		},
		{
			name:   "TLS enabled without verification",
			useTls: true,
			verify: false,
		},
		{
			name:   "TLS disabled",
			useTls: false,
			verify: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := &config.Profile{
				Host:   "test.example.com",
				Port:   443,
				UseTLS: tc.useTls,
				Verify: tc.verify,
			}

			client := New(ctx, cfg)

			if client.UseTls != tc.useTls {
				t.Errorf("Expected UseTls %t, got %t", tc.useTls, client.UseTls)
			}

			if client.Verify != tc.verify {
				t.Errorf("Expected Verify %t, got %t", tc.verify, client.Verify)
			}
		})
	}
}

// TestHttpClientErrorScenarios tests various error conditions
func TestHttpClientErrorScenarios(t *testing.T) {
	testCases := []struct {
		name        string
		setupServer func() *httptest.Server
		expectError bool
	}{
		{
			name: "Server returns 404",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "not found"}`))
				}))
			},
			expectError: false, // 404 is valid HTTP response, not client error
		},
		{
			name: "Server returns 500",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "internal error"}`))
				}))
			},
			expectError: false, // 500 is valid HTTP response, not client error
		},
		{
			name: "Server returns empty body",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}))
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := tc.setupServer()
			defer server.Close()

			serverURL, _ := url.Parse(server.URL)
			port := serverURL.Port()
			portInt := 80
			if port != "" {
				if p, err := strconv.ParseInt(port, 10, 32); err == nil {
					portInt = int(p)
				}
			}

			client := &HttpClient{
				Host:          serverURL.Hostname(),
				Port:          portInt,
				UseTls:        false,
				Verify:        true,
				context:       context.Background(),
				jar:           NewJar(),
				authenticated: true,
			}

			req := NewRequest("/test")
			resp, err := client.Get(req)

			if tc.expectError && err == nil {
				t.Error("Expected error, but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify we get a response even for error status codes
			if !tc.expectError && resp == nil {
				t.Error("Expected response, but got nil")
			}
		})
	}
}

// TestHttpClientConnectionRefused tests behavior when server is unreachable
func TestHttpClientConnectionRefused(t *testing.T) {
	client := &HttpClient{
		Host:          "localhost",
		Port:          9999, // Port that nothing is listening on
		UseTls:        false,
		Verify:        true,
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	req := NewRequest("/test")
	_, err := client.Get(req)

	if err == nil {
		t.Error("Expected connection error, but got none")
	}
}

// TestRequestWithQueryParameters verifies query parameter handling
func TestRequestWithQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters are present
		if r.URL.Query().Get("page") != "1" {
			t.Error("Expected page parameter to be 1")
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Error("Expected limit parameter to be 10")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	req := NewRequest("/test", WithParams(map[string]string{
		"page":  "1",
		"limit": "10",
	}))

	_, err := client.Get(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRequestWithBody verifies request body handling
func TestRequestWithBody(t *testing.T) {
	expectedBody := `{"name":"test","value":123}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, len(expectedBody))
		r.Body.Read(body)

		if string(body) != expectedBody {
			t.Errorf("Expected body %s, got %s", expectedBody, string(body))
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	port := serverURL.Port()
	portInt := 80
	if port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			portInt = int(p)
		}
	}

	client := &HttpClient{
		Host:          serverURL.Hostname(),
		Port:          portInt,
		UseTls:        false,
		Verify:        true,
		context:       context.Background(),
		jar:           NewJar(),
		authenticated: true,
	}

	req := NewRequest("/test", WithBody([]byte(expectedBody)))

	resp, err := client.Post(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}
