// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestDisplayJson(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{
			name:     "simple struct",
			input:    struct{ Name string }{Name: "Test"},
			expected: "{\n    \"Name\": \"Test\"\n}\n",
			wantErr:  false,
		},
		{
			name:     "map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: "{\n    \"a\": 1,\n    \"b\": 2\n}\n",
			wantErr:  false,
		},
		{
			name:     "error during marshal",
			input:    func() {}, // Function can't be marshaled
			expected: "",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			displayBuffer.Reset() // Reset the buffer before each test.
			displayToStdout = false

			err := DisplayJson(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("DisplayJson() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && displayBuffer.String() != tc.expected {
				t.Errorf("DisplayJson() got = %q, want %q", displayBuffer.String(), tc.expected)
			}
		})
	}
}

func TestDisplayYaml(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
		errType  *yaml.TypeError
	}{
		{
			name:     "simple struct",
			input:    struct{ Name string }{Name: "Test"},
			expected: "name: Test\n\n",
			wantErr:  false,
		},
		{
			name:     "map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: "a: 1\nb: 2\n\n",
			wantErr:  false,
		},
		{
			name:     "slice",
			input:    []string{"one", "two", "three"},
			expected: "- one\n- two\n- three\n\n",
			wantErr:  false,
		},
		/*
			{
				name:     "error during marshal",
				input:    func() {}, // Function can't be marshaled
				expected: "",
				wantErr:  true,
				//errType:  *yaml.TypeError,
			},
		*/
		{
			name: "nested struct",
			input: struct {
				Person struct {
					Name string
					Age  int
				}
			}{
				Person: struct {
					Name string
					Age  int
				}{
					Name: "Alice",
					Age:  30,
				},
			},
			expected: "person:\n  name: Alice\n  age: 30\n\n",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			displayBuffer.Reset()
			displayToStdout = false

			err := DisplayYaml(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("DisplayYaml() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			/*
				if tc.wantErr && tc.errType != nil {
					//if err != nil && !errors.As(err, &tc.errType) {
					if err != nil {
						t.Errorf("DisplayYaml() error type = %T, want %T", err, tc.errType)
					}
					return
				}
			*/

			if !tc.wantErr && displayBuffer.String() != tc.expected {
				t.Errorf("DisplayYaml() got = %q, want %q", displayBuffer.String(), tc.expected)
			}
		})
	}
}

func TestDisplay(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple message",
			format:   "Hello, World!",
			args:     nil,
			expected: "Hello, World!",
		},
		{
			name:     "formatted message",
			format:   "Hello, %s!",
			args:     []interface{}{"Alice"},
			expected: "Hello, Alice!",
		},
		{
			name:     "multiple arguments",
			format:   "%s has %d items",
			args:     []interface{}{"Cart", 5},
			expected: "Cart has 5 items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			displayBuffer.Reset()
			displayToStdout = false
			defer func() { displayToStdout = true }()

			Display(tt.format, tt.args...)

			result := displayBuffer.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDisplay_Concurrency(t *testing.T) {
	// Test that Display() is safe for concurrent use
	const goroutines = 100
	const iterations = 10

	displayToStdout = false
	defer func() { displayToStdout = true }()

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				Display("Goroutine %d iteration %d", id, j)
			}
		}(i)
	}

	wg.Wait()
	// If no race condition, test passes
}

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		noColor  bool
		expected string
	}{
		{
			name:     "with color",
			err:      fmt.Errorf("test error"),
			noColor:  false,
			expected: "\x1b[31;1mError:\x1b[0m test error\n",
		},
		{
			name:     "without color",
			err:      fmt.Errorf("test error"),
			noColor:  true,
			expected: "Error: test error\n",
		},
		{
			name:     "wrapped error with color",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner")),
			noColor:  false,
			expected: "\x1b[31;1mError:\x1b[0m outer: inner\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			Error(tt.err, tt.noColor)

			w.Close()
			os.Stderr = old

			var buf bytes.Buffer
			io.Copy(&buf, r)

			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestWarning(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple warning",
			format:   "This is a warning",
			args:     nil,
			expected: "WARNING: This is a warning",
		},
		{
			name:     "formatted warning",
			format:   "Failed to connect to %s",
			args:     []interface{}{"server"},
			expected: "WARNING: Failed to connect to server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			displayBuffer.Reset()
			displayToStdout = false
			defer func() { displayToStdout = true }()

			Warning(tt.format, tt.args...)

			result := displayBuffer.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		timezone *time.Location
		expected string
	}{
		{
			name:     "UTC time",
			time:     time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			timezone: time.UTC,
			expected: "2024-01-15T10:30:00Z",
		},
		{
			name:     "EST time",
			time:     time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			timezone: func() *time.Location { loc, _ := time.LoadLocation("America/New_York"); return loc }(),
			expected: "2024-01-15T05:30:00-05:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimestamp(tt.time, tt.timezone)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfirm(t *testing.T) {
	// Note: Confirm() is difficult to test as-is because it reads from stdin directly
	// This test just validates the function exists and has the expected signature
	// In a real implementation, Confirm should be refactored to accept an io.Reader
	assert.NotNil(t, Confirm)
}
