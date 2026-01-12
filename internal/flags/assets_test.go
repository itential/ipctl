// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetImportCommon(t *testing.T) {
	checkFlags(t, &AssetImportCommon{}, []string{"replace", "repository", "reference", "private-key-file", "params"})
}

func TestAssetExportCommon(t *testing.T) {
	checkFlags(t, &AssetExportCommon{}, []string{"path", "repository", "reference", "private-key-file", "message", "params"})
}

func TestAssetDumpCommon(t *testing.T) {
	checkFlags(t, &AssetDumpCommon{}, []string{"path", "repository", "reference", "private-key-file", "message", "params"})
}

func TestAssetLoadCommon(t *testing.T) {
	checkFlags(t, &AssetLoadCommon{}, []string{"repository", "reference", "private-key-file", "params"})
}

func TestAssetCopyCommon(t *testing.T) {
	checkFlags(t, &AssetCopyCommon{}, []string{"to", "from", "replace", "params"})
}

func TestAssetImportCommonGetters(t *testing.T) {
	common := &AssetImportCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "main",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "main", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "", common.GetPath())
}

func TestAssetExportCommonGetters(t *testing.T) {
	common := &AssetExportCommon{
		Path:           "/export/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "develop",
		PrivateKeyFile: "/path/to/key",
		Message:        "Export commit message",
	}

	assert.Equal(t, "/export/path", common.GetPath())
	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "develop", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "Export commit message", common.GetMessage())
}

func TestAssetDumpCommonGetters(t *testing.T) {
	common := &AssetDumpCommon{
		Path:           "/dump/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "feature",
		PrivateKeyFile: "/path/to/key",
		Message:        "Dump commit message",
	}

	assert.Equal(t, "/dump/path", common.GetPath())
	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "feature", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
	assert.Equal(t, "Dump commit message", common.GetMessage())
}

func TestAssetLoadCommonGetters(t *testing.T) {
	common := &AssetLoadCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "v1.0.0",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", common.GetRepository())
	assert.Equal(t, "v1.0.0", common.GetReference())
	assert.Equal(t, "/path/to/key", common.GetPrivateKeyFile())
}

func TestAssetCopyCommonFields(t *testing.T) {
	common := &AssetCopyCommon{
		To:      "destination",
		From:    "source",
		Replace: true,
	}

	assert.Equal(t, "destination", common.To)
	assert.Equal(t, "source", common.From)
	assert.True(t, common.Replace)
}

func TestGitterInterface(t *testing.T) {
	var gitter Gitter = &AssetImportCommon{
		Repository:     "https://github.com/example/repo",
		Reference:      "main",
		PrivateKeyFile: "/path/to/key",
	}

	assert.Equal(t, "https://github.com/example/repo", gitter.GetRepository())
	assert.Equal(t, "main", gitter.GetReference())
	assert.Equal(t, "/path/to/key", gitter.GetPrivateKeyFile())
}

func TestCommitterInterface(t *testing.T) {
	var committer Committer = &AssetExportCommon{
		Path:           "/export/path",
		Repository:     "https://github.com/example/repo",
		Reference:      "develop",
		PrivateKeyFile: "/path/to/key",
		Message:        "Export commit message",
	}

	assert.Equal(t, "/export/path", committer.GetPath())
	assert.Equal(t, "Export commit message", committer.GetMessage())
}

// TestParseParams tests the ParseParams function with various inputs
func TestParseParams(t *testing.T) {
	tests := []struct {
		name        string
		params      []string
		expected    map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty params",
			params:      []string{},
			expected:    nil,
			expectError: false,
		},
		{
			name:        "nil params",
			params:      nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:     "single valid param",
			params:   []string{"key=value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "multiple valid params",
			params:   []string{"limit=10", "offset=20", "status=active"},
			expected: map[string]string{"limit": "10", "offset": "20", "status": "active"},
		},
		{
			name:     "param with empty value",
			params:   []string{"key="},
			expected: map[string]string{"key": ""},
		},
		{
			name:     "param with equals in value",
			params:   []string{"key=value=with=equals"},
			expected: map[string]string{"key": "value=with=equals"},
		},
		{
			name:     "param with special characters in value",
			params:   []string{"url=https://example.com/path?query=1"},
			expected: map[string]string{"url": "https://example.com/path?query=1"},
		},
		{
			name:     "param with spaces in value",
			params:   []string{"message=hello world"},
			expected: map[string]string{"message": "hello world"},
		},
		{
			name:     "param with unicode characters",
			params:   []string{"name=José", "city=München"},
			expected: map[string]string{"name": "José", "city": "München"},
		},
		{
			name:     "duplicate keys (last one wins)",
			params:   []string{"key=first", "key=second"},
			expected: map[string]string{"key": "second"},
		},
		{
			name:        "missing equals",
			params:      []string{"invalidparam"},
			expectError: true,
			errorMsg:    `invalid param format "invalidparam": expected key=value`,
		},
		{
			name:        "empty key",
			params:      []string{"=value"},
			expectError: true,
			errorMsg:    `invalid param format "=value": key cannot be empty`,
		},
		{
			name:        "only equals",
			params:      []string{"="},
			expectError: true,
			errorMsg:    `invalid param format "=": key cannot be empty`,
		},
		{
			name:        "multiple errors (returns first)",
			params:      []string{"valid=1", "invalid", "also=bad"},
			expectError: true,
			errorMsg:    `invalid param format "invalid": expected key=value`,
		},
		{
			name:     "complex real-world example",
			params:   []string{"filter=name eq 'test'", "limit=25", "offset=0", "sort=-createdAt"},
			expected: map[string]string{"filter": "name eq 'test'", "limit": "25", "offset": "0", "sort": "-createdAt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseParams(tt.params)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestAssetImportCommon_ParseParams tests AssetImportCommon.ParseParams
func TestAssetImportCommon_ParseParams(t *testing.T) {
	tests := []struct {
		name        string
		params      []string
		expected    map[string]string
		expectError bool
	}{
		{
			name:        "empty params",
			params:      []string{},
			expected:    nil,
			expectError: false,
		},
		{
			name:     "valid params",
			params:   []string{"limit=10", "offset=20"},
			expected: map[string]string{"limit": "10", "offset": "20"},
		},
		{
			name:        "invalid params",
			params:      []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common := &AssetImportCommon{
				Params: tt.params,
			}

			result, err := common.ParseParams()

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestAssetExportCommon_ParseParams tests AssetExportCommon.ParseParams
func TestAssetExportCommon_ParseParams(t *testing.T) {
	common := &AssetExportCommon{
		Params: []string{"key1=value1", "key2=value2"},
	}

	result, err := common.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, result)
}

// TestAssetDumpCommon_ParseParams tests AssetDumpCommon.ParseParams
func TestAssetDumpCommon_ParseParams(t *testing.T) {
	common := &AssetDumpCommon{
		Params: []string{"format=json"},
	}

	result, err := common.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"format": "json"}, result)
}

// TestAssetLoadCommon_ParseParams tests AssetLoadCommon.ParseParams
func TestAssetLoadCommon_ParseParams(t *testing.T) {
	common := &AssetLoadCommon{
		Params: []string{"validate=true"},
	}

	result, err := common.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"validate": "true"}, result)
}

// TestAssetCopyCommon_ParseParams tests AssetCopyCommon.ParseParams
func TestAssetCopyCommon_ParseParams(t *testing.T) {
	common := &AssetCopyCommon{
		Params: []string{"dryRun=false"},
	}

	result, err := common.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"dryRun": "false"}, result)
}

// TestAssetImportCommon_GetParams tests AssetImportCommon.GetParams
func TestAssetImportCommon_GetParams(t *testing.T) {
	params := []string{"key=value"}
	common := &AssetImportCommon{
		Params: params,
	}

	result := common.GetParams()

	assert.Equal(t, params, result)
}

// TestAssetExportCommon_GetParams tests AssetExportCommon.GetParams
func TestAssetExportCommon_GetParams(t *testing.T) {
	params := []string{"key=value"}
	common := &AssetExportCommon{
		Params: params,
	}

	result := common.GetParams()

	assert.Equal(t, params, result)
}

// TestAssetDumpCommon_GetParams tests AssetDumpCommon.GetParams
func TestAssetDumpCommon_GetParams(t *testing.T) {
	params := []string{"key=value"}
	common := &AssetDumpCommon{
		Params: params,
	}

	result := common.GetParams()

	assert.Equal(t, params, result)
}

// TestAssetLoadCommon_GetParams tests AssetLoadCommon.GetParams
func TestAssetLoadCommon_GetParams(t *testing.T) {
	params := []string{"key=value"}
	common := &AssetLoadCommon{
		Params: params,
	}

	result := common.GetParams()

	assert.Equal(t, params, result)
}

// TestAssetCopyCommon_GetParams tests AssetCopyCommon.GetParams
func TestAssetCopyCommon_GetParams(t *testing.T) {
	params := []string{"key=value"}
	common := &AssetCopyCommon{
		Params: params,
	}

	result := common.GetParams()

	assert.Equal(t, params, result)
}

// TestParamer_Interface verifies that all common structs implement the Paramer interface
func TestParamer_Interface(t *testing.T) {
	var _ Paramer = (*AssetImportCommon)(nil)
	var _ Paramer = (*AssetExportCommon)(nil)
	var _ Paramer = (*AssetDumpCommon)(nil)
	var _ Paramer = (*AssetLoadCommon)(nil)
	var _ Paramer = (*AssetCopyCommon)(nil)
}
