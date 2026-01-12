// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiGetOptions(t *testing.T) {
	checkFlags(t, &ApiGetOptions{}, []string{"expected-status-code", "params"})
}

func TestApiPutOptions(t *testing.T) {
	checkFlags(t, &ApiPutOptions{}, []string{"data", "expected-status-code", "params"})
}

func TestApiDeleteOptions(t *testing.T) {
	checkFlags(t, &ApiDeleteOptions{}, []string{"expected-status-code", "params"})
}

func TestApiPostOptions(t *testing.T) {
	checkFlags(t, &ApiPostOptions{}, []string{"data", "expected-status-code", "params"})
}

func TestApiPatchOptions(t *testing.T) {
	checkFlags(t, &ApiPatchOptions{}, []string{"data", "expected-status-code", "params"})
}

// TestApiGetOptions_ParseParams tests ApiGetOptions.ParseParams
func TestApiGetOptions_ParseParams(t *testing.T) {
	options := &ApiGetOptions{
		Params: []string{"limit=50", "offset=100"},
	}

	result, err := options.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"limit": "50", "offset": "100"}, result)
}

// TestApiPutOptions_ParseParams tests ApiPutOptions.ParseParams
func TestApiPutOptions_ParseParams(t *testing.T) {
	options := &ApiPutOptions{
		Params: []string{"key=value", "limit=10"},
	}

	result, err := options.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"key": "value", "limit": "10"}, result)
}

// TestApiDeleteOptions_ParseParams tests ApiDeleteOptions.ParseParams
func TestApiDeleteOptions_ParseParams(t *testing.T) {
	options := &ApiDeleteOptions{
		Params: []string{"force=true"},
	}

	result, err := options.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"force": "true"}, result)
}

// TestApiPostOptions_ParseParams tests ApiPostOptions.ParseParams
func TestApiPostOptions_ParseParams(t *testing.T) {
	options := &ApiPostOptions{
		Params: []string{"validate=true"},
	}

	result, err := options.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"validate": "true"}, result)
}

// TestApiPatchOptions_ParseParams tests ApiPatchOptions.ParseParams
func TestApiPatchOptions_ParseParams(t *testing.T) {
	options := &ApiPatchOptions{
		Params: []string{"partial=true"},
	}

	result, err := options.ParseParams()

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"partial": "true"}, result)
}

// TestApiOptions_GetParams tests GetParams methods for API options
func TestApiOptions_GetParams(t *testing.T) {
	params := []string{"key=value"}

	t.Run("ApiGetOptions", func(t *testing.T) {
		options := &ApiGetOptions{Params: params}
		assert.Equal(t, params, options.GetParams())
	})

	t.Run("ApiPutOptions", func(t *testing.T) {
		options := &ApiPutOptions{Params: params}
		assert.Equal(t, params, options.GetParams())
	})

	t.Run("ApiDeleteOptions", func(t *testing.T) {
		options := &ApiDeleteOptions{Params: params}
		assert.Equal(t, params, options.GetParams())
	})

	t.Run("ApiPostOptions", func(t *testing.T) {
		options := &ApiPostOptions{Params: params}
		assert.Equal(t, params, options.GetParams())
	})

	t.Run("ApiPatchOptions", func(t *testing.T) {
		options := &ApiPatchOptions{Params: params}
		assert.Equal(t, params, options.GetParams())
	})
}
