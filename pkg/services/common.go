// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/url"
	"strconv"
)

// Metadata contains pagination and response metadata
type Metadata struct {
	Skip  int `json:"skip,omitempty"`
	Limit int `json:"limit,omitempty"`
	Total int `json:"total,omitempty"`

	// NOTE (privateip) commenting these fields out as they are current not
	// being used by the application and the data returned from the API is
	// inconsistent
	//NextPageSkip     string `json:"nextPageSkip,omitempty"`
	//PreviousPageSkip string `json:"previousPageSkip,omitempty"`
	//CurrentPageSize  int    `json:"currentPageSize,omitempty"`
}

// Gbac represents Group-Based Access Control permissions
type Gbac struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items    []interface{} `json:"items"`
	Count    int           `json:"count"`
	End      int           `json:"end"`
	Limit    int           `json:"limit"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Skip     int           `json:"skip"`
	Total    int           `json:"total"`
}

// ErrorMessage represents an API error response
type ErrorMessage struct {
	Message  string                 `json:"message"`
	Data     string                 `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Params defines the interface for query parameters
type Params interface {
	Query() map[string]string
}

// RawParams wraps url.Values for query parameters
type RawParams struct {
	Values url.Values
}

// Query converts RawParams to a map of query parameters
func (p *RawParams) Query() map[string]string {
	params := make(map[string]string)
	for key, value := range p.Values {
		params[key] = value[0]
	}
	return params
}

// QueryParams provides structured query parameter handling
type QueryParams struct {
	Contains        string
	ContainsField   string
	Equals          string
	EqualsField     string
	StartsWith      string
	StartsWithField string
	Skip            int
	Limit           int
	Sort            string
	Order           int
	SkipActiveSync  bool
	Raw             map[string]string
}

// Query converts QueryParams to a map of query parameters
func (p *QueryParams) Query() map[string]string {
	m := make(map[string]string)

	for key, value := range p.Raw {
		m[key] = value
	}

	if p.Contains != "" {
		m["contains"] = p.Contains
	}
	if p.ContainsField != "" {
		m["containsField"] = p.ContainsField
	}
	if p.Equals != "" {
		m["equals"] = p.Equals
	}
	if p.EqualsField != "" {
		m["equalsField"] = p.EqualsField
	}
	if p.StartsWith != "" {
		m["startsWith"] = p.StartsWith
	}
	if p.StartsWithField != "" {
		m["startsWithField"] = p.StartsWithField
	}
	if p.Skip != 0 {
		m["skip"] = strconv.Itoa(p.Skip)
	}
	if p.Limit != 0 {
		m["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Sort != "" {
		m["sort"] = p.Sort
	}
	if p.Order != 0 {
		m["order"] = strconv.Itoa(p.Order)
	}
	if p.SkipActiveSync {
		m["skipActiveSync"] = "true"
	}

	return m
}
