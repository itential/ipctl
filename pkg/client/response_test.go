// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package client

import (
	"net/http"
	"testing"
)

func TestResponseStruct(t *testing.T) {
	testCases := []struct {
		name       string
		response   Response
		wantMethod string
		wantUrl    string
		wantStatus string
		wantCode   int
		wantBody   string
	}{
		{
			name: "Complete response",
			response: Response{
				Method:     "GET",
				Url:        "https://api.example.com/users",
				Status:     "200 OK",
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"users": []}`),
			},
			wantMethod: "GET",
			wantUrl:    "https://api.example.com/users",
			wantStatus: "200 OK",
			wantCode:   200,
			wantBody:   `{"users": []}`,
		},
		{
			name: "POST response with error",
			response: Response{
				Method:     "POST",
				Url:        "https://api.example.com/users",
				Status:     "400 Bad Request",
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"error": "invalid input"}`),
			},
			wantMethod: "POST",
			wantUrl:    "https://api.example.com/users",
			wantStatus: "400 Bad Request",
			wantCode:   400,
			wantBody:   `{"error": "invalid input"}`,
		},
		{
			name: "Empty body response",
			response: Response{
				Method:     "DELETE",
				Url:        "https://api.example.com/users/123",
				Status:     "204 No Content",
				StatusCode: 204,
				Headers:    map[string]string{"Content-Length": "0"},
				Body:       []byte{},
			},
			wantMethod: "DELETE",
			wantUrl:    "https://api.example.com/users/123",
			wantStatus: "204 No Content",
			wantCode:   204,
			wantBody:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.response.Method != tc.wantMethod {
				t.Errorf("Expected Method %s, got %s", tc.wantMethod, tc.response.Method)
			}

			if tc.response.Url != tc.wantUrl {
				t.Errorf("Expected Url %s, got %s", tc.wantUrl, tc.response.Url)
			}

			if tc.response.Status != tc.wantStatus {
				t.Errorf("Expected Status %s, got %s", tc.wantStatus, tc.response.Status)
			}

			if tc.response.StatusCode != tc.wantCode {
				t.Errorf("Expected StatusCode %d, got %d", tc.wantCode, tc.response.StatusCode)
			}

			if string(tc.response.Body) != tc.wantBody {
				t.Errorf("Expected Body %s, got %s", tc.wantBody, string(tc.response.Body))
			}

			if tc.response.Headers == nil {
				t.Error("Expected Headers to be initialized")
			}
		})
	}
}

func TestResponseHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Authorization":  "Bearer token123",
		"X-Custom-Field": "custom-value",
	}

	response := Response{
		Method:     "GET",
		Url:        "https://example.com/api",
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       []byte(`{"test": "data"}`),
	}

	for key, expectedValue := range headers {
		if actualValue, exists := response.Headers[key]; !exists {
			t.Errorf("Expected header %s to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, actualValue)
		}
	}
}

func TestResponseStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		status     string
	}{
		{"Success", http.StatusOK, "200 OK"},
		{"Created", http.StatusCreated, "201 Created"},
		{"No Content", http.StatusNoContent, "204 No Content"},
		{"Bad Request", http.StatusBadRequest, "400 Bad Request"},
		{"Unauthorized", http.StatusUnauthorized, "401 Unauthorized"},
		{"Forbidden", http.StatusForbidden, "403 Forbidden"},
		{"Not Found", http.StatusNotFound, "404 Not Found"},
		{"Internal Server Error", http.StatusInternalServerError, "500 Internal Server Error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := Response{
				Method:     "GET",
				Url:        "https://example.com/api",
				Status:     tc.status,
				StatusCode: tc.statusCode,
				Headers:    make(map[string]string),
				Body:       []byte{},
			}

			if response.StatusCode != tc.statusCode {
				t.Errorf("Expected StatusCode %d, got %d", tc.statusCode, response.StatusCode)
			}

			if response.Status != tc.status {
				t.Errorf("Expected Status %s, got %s", tc.status, response.Status)
			}
		})
	}
}

func TestResponseBodyTypes(t *testing.T) {
	testCases := []struct {
		name string
		body []byte
		want string
	}{
		{
			name: "JSON body",
			body: []byte(`{"message": "success", "data": {"id": 123}}`),
			want: `{"message": "success", "data": {"id": 123}}`,
		},
		{
			name: "Plain text body",
			body: []byte("Simple text response"),
			want: "Simple text response",
		},
		{
			name: "Empty body",
			body: []byte{},
			want: "",
		},
		{
			name: "XML body",
			body: []byte(`<?xml version="1.0"?><root><item>test</item></root>`),
			want: `<?xml version="1.0"?><root><item>test</item></root>`,
		},
		{
			name: "HTML body",
			body: []byte(`<html><body><h1>Test Page</h1></body></html>`),
			want: `<html><body><h1>Test Page</h1></body></html>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := Response{
				Method:     "GET",
				Url:        "https://example.com/api",
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Headers:    make(map[string]string),
				Body:       tc.body,
			}

			if string(response.Body) != tc.want {
				t.Errorf("Expected body %s, got %s", tc.want, string(response.Body))
			}
		})
	}
}
