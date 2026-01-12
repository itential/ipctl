// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type ApiRunner struct {
	client  client.Client
	service *services.ApiService
}

func NewApiRunner(client client.Client) ApiRunner {
	return ApiRunner{
		service: services.NewApiService(client),
	}
}

// appendQueryParams appends query parameters to a URL path.
// If the path already contains query parameters, the new params are appended.
// Parameters are automatically URL encoded.
//
// Example:
//
//	appendQueryParams("/api/projects", map[string]string{"limit": "10"})
//	// Returns: "/api/projects?limit=10"
//
//	appendQueryParams("/api/projects?status=active", map[string]string{"limit": "10"})
//	// Returns: "/api/projects?status=active&limit=10"
func appendQueryParams(urlPath string, params map[string]string) string {
	logging.Trace()

	if len(params) == 0 {
		return urlPath
	}

	// Parse the URL to extract existing query parameters
	u, err := url.Parse(urlPath)
	if err != nil {
		// If parsing fails, just return the original path
		return urlPath
	}

	// Get existing query values
	q := u.Query()

	// Add new parameters
	for key, value := range params {
		q.Set(key, value)
	}

	// Update the query string
	u.RawQuery = q.Encode()

	return u.String()
}

func (r *ApiRunner) readData(in string) (map[string]interface{}, error) {
	logging.Trace()

	var data map[string]interface{}

	var content []byte
	var err error

	if strings.HasPrefix(in, "@") {
		content, err = os.ReadFile(in[1:len(in)])
		if err != nil {
			return nil, err
		}
	} else if in != "" {
		content = []byte(in)
	}

	if in != "" {
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}
	}

	return data, nil

}
func (r *ApiRunner) jsonResponse(in string) (interface{}, error) {
	logging.Trace()

	var response interface{}

	if err := json.Unmarshal([]byte(in), &response); err != nil {
		return nil, err
	}

	return response, nil

}

func (r *ApiRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.ApiGetOptions
	utils.LoadObject(in.Common, &options)

	// Build URL with query parameters
	urlPath := in.Args[0]
	queryParams, err := options.ParseParams()
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		urlPath = appendQueryParams(urlPath, queryParams)
	}

	res, err := r.service.Get(urlPath, options.ExpectedStatusCode)
	if err != nil {
		return nil, err
	}

	jsonRes, err := r.jsonResponse(res)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: jsonRes,
	}, nil
}

func (r *ApiRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.ApiDeleteOptions
	utils.LoadObject(in.Common, &options)

	// Build URL with query parameters
	urlPath := in.Args[0]
	queryParams, err := options.ParseParams()
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		urlPath = appendQueryParams(urlPath, queryParams)
	}

	res, err := r.service.Delete(urlPath, options.ExpectedStatusCode)
	if err != nil {
		return nil, err
	}

	jsonRes, err := r.jsonResponse(res)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: jsonRes,
	}, nil
}

func (r *ApiRunner) Post(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.ApiPostOptions
	utils.LoadObject(in.Common, &options)

	body, err := r.readData(options.Data)
	if err != nil {
		return nil, err
	}

	// Build URL with query parameters
	urlPath := in.Args[0]
	queryParams, err := options.ParseParams()
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		urlPath = appendQueryParams(urlPath, queryParams)
	}

	res, err := r.service.Post(urlPath, body, options.ExpectedStatusCode)
	if err != nil {
		return nil, err
	}

	jsonRes, err := r.jsonResponse(res)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: jsonRes,
	}, nil
}

func (r *ApiRunner) Put(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.ApiPutOptions
	utils.LoadObject(in.Common, &options)

	body, err := r.readData(options.Data)
	if err != nil {
		return nil, err
	}

	// Build URL with query parameters
	urlPath := in.Args[0]
	queryParams, err := options.ParseParams()
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		urlPath = appendQueryParams(urlPath, queryParams)
	}

	res, err := r.service.Put(urlPath, body, options.ExpectedStatusCode)
	if err != nil {
		return nil, err
	}

	jsonRes, err := r.jsonResponse(res)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: jsonRes,
	}, nil
}

func (r *ApiRunner) Patch(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.ApiPatchOptions
	utils.LoadObject(in.Common, &options)

	body, err := r.readData(options.Data)
	if err != nil {
		return nil, err
	}

	// Build URL with query parameters
	urlPath := in.Args[0]
	queryParams, err := options.ParseParams()
	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		urlPath = appendQueryParams(urlPath, queryParams)
	}

	res, err := r.service.Patch(urlPath, body, options.ExpectedStatusCode)
	if err != nil {
		return nil, err
	}

	jsonRes, err := r.jsonResponse(res)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: jsonRes,
	}, nil
}
