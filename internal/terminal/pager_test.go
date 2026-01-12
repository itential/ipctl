// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "short columns",
			input:    "Col1\tCol2\nValue1\tValue2\n",
			expected: "Col1\tCol2\t\nValue1\tValue2\t\n",
		},
		{
			name:     "long column gets truncated",
			input:    "Column1\t" + strings.Repeat("x", 60) + "\n",
			expected: "Column1\t" + strings.Repeat("x", 49) + "…\t\n",
		},
		{
			name:     "multiple long columns",
			input:    strings.Repeat("a", 60) + "\t" + strings.Repeat("b", 60) + "\n",
			expected: strings.Repeat("a", 49) + "…\t" + strings.Repeat("b", 49) + "…\t\n",
		},
		{
			name:     "empty line preserved",
			input:    "Col1\tCol2\n\nValue1\tValue2\n",
			expected: "Col1\tCol2\t\nValue1\tValue2\t\n",
		},
		{
			name:     "exactly 50 characters not truncated",
			input:    strings.Repeat("x", 50) + "\n",
			expected: strings.Repeat("x", 50) + "\t\n",
		},
		{
			name:     "49 characters not truncated",
			input:    strings.Repeat("x", 49) + "\n",
			expected: strings.Repeat("x", 49) + "\t\n",
		},
		{
			name:     "51 characters truncated to 49 plus ellipsis",
			input:    strings.Repeat("x", 51) + "\n",
			expected: strings.Repeat("x", 49) + "…\t\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateOutput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncateOutput_MultipleRows(t *testing.T) {
	input := "Header1\tHeader2\tHeader3\n" +
		"Value1\tValue2\tValue3\n" +
		strings.Repeat("a", 60) + "\t" + strings.Repeat("b", 60) + "\t" + strings.Repeat("c", 60) + "\n"

	result := truncateOutput(input)

	lines := strings.Split(result, "\n")
	assert.Equal(t, 4, len(lines)) // 3 data lines + 1 empty at end

	// Check that long values are truncated
	lastLine := lines[2]
	assert.Contains(t, lastLine, "…")
}

func TestTruncateOutput_EmptyInput(t *testing.T) {
	result := truncateOutput("")
	assert.Equal(t, "", result)
}

func TestTruncateOutput_OnlyNewlines(t *testing.T) {
	result := truncateOutput("\n\n\n")
	assert.Equal(t, "", result)
}
