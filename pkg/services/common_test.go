// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParams_Query(t *testing.T) {
	params := &QueryParams{
		Contains:        "test",
		ContainsField:   "name",
		Equals:          "value",
		EqualsField:     "id",
		StartsWith:      "prefix",
		StartsWithField: "title",
		Skip:            10,
		Limit:           50,
		Sort:            "name",
		Order:           1,
		SkipActiveSync:  true,
		Raw: map[string]string{
			"custom": "value",
		},
	}

	result := params.Query()

	assert.Equal(t, "test", result["contains"])
	assert.Equal(t, "name", result["containsField"])
	assert.Equal(t, "value", result["equals"])
	assert.Equal(t, "id", result["equalsField"])
	assert.Equal(t, "prefix", result["startsWith"])
	assert.Equal(t, "title", result["startsWithField"])
	assert.Equal(t, "10", result["skip"])
	assert.Equal(t, "50", result["limit"])
	assert.Equal(t, "name", result["sort"])
	assert.Equal(t, "1", result["order"])
	assert.Equal(t, "true", result["skipActiveSync"])
	assert.Equal(t, "value", result["custom"])
}

func TestQueryParams_Query_EmptyValues(t *testing.T) {
	params := &QueryParams{}

	result := params.Query()

	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestQueryParams_Query_PartialValues(t *testing.T) {
	params := &QueryParams{
		Contains: "test",
		Skip:     5,
		Raw: map[string]string{
			"custom1": "value1",
			"custom2": "value2",
		},
	}

	result := params.Query()

	assert.Equal(t, "test", result["contains"])
	assert.Equal(t, "5", result["skip"])
	assert.Equal(t, "value1", result["custom1"])
	assert.Equal(t, "value2", result["custom2"])
	assert.Equal(t, 4, len(result))
}

func TestRawParams_Query(t *testing.T) {
	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("param2", "value2")
	values.Add("param3", "value3")

	params := &RawParams{Values: values}
	result := params.Query()

	assert.Equal(t, "value1", result["param1"])
	assert.Equal(t, "value2", result["param2"])
	assert.Equal(t, "value3", result["param3"])
}

func TestRawParams_Query_Empty(t *testing.T) {
	params := &RawParams{Values: url.Values{}}
	result := params.Query()

	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestRawParams_Query_MultipleValues(t *testing.T) {
	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("param1", "value2") // Only first value should be used

	params := &RawParams{Values: values}
	result := params.Query()

	assert.Equal(t, "value1", result["param1"])
}
