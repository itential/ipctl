// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"testing"

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
