// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAppendQueryParams tests the appendQueryParams helper function
func TestAppendQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		urlPath  string
		params   map[string]string
		expected string
	}{
		{
			name:     "empty params",
			urlPath:  "/api/projects",
			params:   map[string]string{},
			expected: "/api/projects",
		},
		{
			name:     "nil params",
			urlPath:  "/api/projects",
			params:   nil,
			expected: "/api/projects",
		},
		{
			name:    "single param",
			urlPath: "/api/projects",
			params: map[string]string{
				"limit": "10",
			},
			expected: "/api/projects?limit=10",
		},
		{
			name:    "multiple params",
			urlPath: "/api/workflows",
			params: map[string]string{
				"limit":  "25",
				"offset": "50",
			},
			// Note: order may vary due to map iteration
			// We'll check both are present separately
		},
		{
			name:    "append to existing query string",
			urlPath: "/api/projects?status=active",
			params: map[string]string{
				"limit": "10",
			},
			expected: "/api/projects?limit=10&status=active",
		},
		{
			name:    "special characters in values",
			urlPath: "/api/search",
			params: map[string]string{
				"query": "hello world",
			},
			expected: "/api/search?query=hello+world",
		},
		{
			name:    "url encoding",
			urlPath: "/api/data",
			params: map[string]string{
				"filter": "name eq 'test'",
			},
			expected: "/api/data?filter=name+eq+%27test%27",
		},
		{
			name:    "empty value",
			urlPath: "/api/items",
			params: map[string]string{
				"key": "",
			},
			expected: "/api/items?key=",
		},
		{
			name:    "complex url with path and query",
			urlPath: "/api/v2/automations?type=workflow",
			params: map[string]string{
				"limit": "100",
			},
			expected: "/api/v2/automations?limit=100&type=workflow",
		},
		{
			name:    "override existing param",
			urlPath: "/api/projects?limit=10",
			params: map[string]string{
				"limit": "50",
			},
			expected: "/api/projects?limit=50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendQueryParams(tt.urlPath, tt.params)

			// For tests with multiple params, we can't guarantee order
			if tt.name == "multiple params" {
				assert.Contains(t, result, "limit=25")
				assert.Contains(t, result, "offset=50")
				assert.Contains(t, result, "/api/workflows?")
			} else if tt.expected != "" {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestAppendQueryParams_InvalidURL tests handling of invalid URLs
func TestAppendQueryParams_InvalidURL(t *testing.T) {
	// Invalid URLs should return the original path
	invalidURL := "://invalid"
	params := map[string]string{"key": "value"}

	result := appendQueryParams(invalidURL, params)

	// Should return original when URL parsing fails
	assert.Equal(t, invalidURL, result)
}

// TestAppendQueryParams_PreservesFragment tests that URL fragments are preserved
func TestAppendQueryParams_PreservesFragment(t *testing.T) {
	urlPath := "/api/docs#section"
	params := map[string]string{"page": "1"}

	result := appendQueryParams(urlPath, params)

	assert.Contains(t, result, "page=1")
	assert.Contains(t, result, "#section")
}

// TestAppendQueryParams_EmptyPath tests empty path handling
func TestAppendQueryParams_EmptyPath(t *testing.T) {
	urlPath := ""
	params := map[string]string{"key": "value"}

	result := appendQueryParams(urlPath, params)

	assert.Equal(t, "?key=value", result)
}
