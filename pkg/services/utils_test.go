// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	// Test with a simple struct
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
		Flag  bool   `json:"flag"`
	}

	input := map[string]interface{}{
		"name":  "test",
		"value": 42,
		"flag":  true,
	}

	var result TestStruct
	err := Unmarshal(input, &result)

	assert.Nil(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
	assert.True(t, result.Flag)
}

func TestUnmarshal_WithNestedStruct(t *testing.T) {
	type NestedStruct struct {
		Field string `json:"field"`
	}

	type TestStruct struct {
		Name   string       `json:"name"`
		Nested NestedStruct `json:"nested"`
	}

	input := map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"field": "nested_value",
		},
	}

	var result TestStruct
	err := Unmarshal(input, &result)

	assert.Nil(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, "nested_value", result.Nested.Field)
}

func TestUnmarshal_WithSlice(t *testing.T) {
	type TestStruct struct {
		Names []string `json:"names"`
		IDs   []int    `json:"ids"`
	}

	input := map[string]interface{}{
		"names": []interface{}{"alice", "bob", "charlie"},
		"ids":   []interface{}{1, 2, 3},
	}

	var result TestStruct
	err := Unmarshal(input, &result)

	assert.Nil(t, err)
	assert.Equal(t, []string{"alice", "bob", "charlie"}, result.Names)
	assert.Equal(t, []int{1, 2, 3}, result.IDs)
}

func TestUnmarshal_WithEmptyInput(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	input := map[string]interface{}{}

	var result TestStruct
	err := Unmarshal(input, &result)

	assert.Nil(t, err)
	assert.Equal(t, "", result.Name)
}

func TestUnmarshal_WithNilInput(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	var input map[string]interface{}
	var result TestStruct
	err := Unmarshal(input, &result)

	assert.Nil(t, err)
	assert.Equal(t, "", result.Name)
}
