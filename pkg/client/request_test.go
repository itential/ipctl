// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		opts         []RequestOption
		expectedPath string
	}{
		{
			name:         "Simple request without options",
			path:         "/api/v1/test",
			opts:         nil,
			expectedPath: "/api/v1/test",
		},
		{
			name: "Request with headers option",
			path: "/api/v1/users",
			opts: []RequestOption{
				WithHeaders(map[string]string{"Authorization": "Bearer token"}),
			},
			expectedPath: "/api/v1/users",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := NewRequest(tc.path, tc.opts...)

			if req.Path != tc.expectedPath {
				t.Errorf("Expected path %s, got %s", tc.expectedPath, req.Path)
			}
		})
	}
}

func TestWithHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	req := NewRequest("/test", WithHeaders(headers))

	if !reflect.DeepEqual(req.Headers, headers) {
		t.Errorf("Expected headers %v, got %v", headers, req.Headers)
	}
}

func TestWithParams(t *testing.T) {
	params := map[string]string{
		"page":   "1",
		"limit":  "10",
		"filter": "active",
	}

	req := NewRequest("/test", WithParams(params))

	if !reflect.DeepEqual(req.Params, params) {
		t.Errorf("Expected params %v, got %v", params, req.Params)
	}
}

func TestWithBody(t *testing.T) {
	body := []byte(`{"name": "test", "value": 123}`)

	req := NewRequest("/test", WithBody(body))

	if !reflect.DeepEqual(req.Body, body) {
		t.Errorf("Expected body %s, got %s", string(body), string(req.Body))
	}
}

func TestWithNoLog(t *testing.T) {
	testCases := []struct {
		name     string
		noLog    bool
		expected bool
	}{
		{
			name:     "NoLog set to true",
			noLog:    true,
			expected: true,
		},
		{
			name:     "NoLog set to false",
			noLog:    false,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := NewRequest("/test", WithNoLog(tc.noLog))

			if req.NoLog != tc.expected {
				t.Errorf("Expected NoLog %t, got %t", tc.expected, req.NoLog)
			}
		})
	}
}

func TestMultipleOptions(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	params := map[string]string{"page": "1"}
	body := []byte(`{"test": "data"}`)

	req := NewRequest("/api/test",
		WithHeaders(headers),
		WithParams(params),
		WithBody(body),
		WithNoLog(true),
	)

	if req.Path != "/api/test" {
		t.Errorf("Expected path /api/test, got %s", req.Path)
	}

	if !reflect.DeepEqual(req.Headers, headers) {
		t.Errorf("Expected headers %v, got %v", headers, req.Headers)
	}

	if !reflect.DeepEqual(req.Params, params) {
		t.Errorf("Expected params %v, got %v", params, req.Params)
	}

	if !reflect.DeepEqual(req.Body, body) {
		t.Errorf("Expected body %s, got %s", string(body), string(req.Body))
	}

	if !req.NoLog {
		t.Error("Expected NoLog to be true")
	}
}

func TestRequestOptionFunctional(t *testing.T) {
	req := &Request{Path: "/initial"}

	headerOption := WithHeaders(map[string]string{"Test": "Value"})
	headerOption(req)

	if req.Headers["Test"] != "Value" {
		t.Errorf("Expected header Test to be Value, got %s", req.Headers["Test"])
	}

	paramOption := WithParams(map[string]string{"param": "value"})
	paramOption(req)

	if req.Params["param"] != "value" {
		t.Errorf("Expected param param to be value, got %s", req.Params["param"])
	}

	bodyOption := WithBody([]byte("test body"))
	bodyOption(req)

	if string(req.Body) != "test body" {
		t.Errorf("Expected body 'test body', got '%s'", string(req.Body))
	}

	noLogOption := WithNoLog(true)
	noLogOption(req)

	if !req.NoLog {
		t.Error("Expected NoLog to be true")
	}
}
