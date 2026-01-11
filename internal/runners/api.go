// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
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

	res, err := r.service.Get(in.Args[0], 0)
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

	res, err := r.service.Delete(in.Args[0], options.ExpectedStatusCode)
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

	res, err := r.service.Post(in.Args[0], body, options.ExpectedStatusCode)
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

	res, err := r.service.Put(in.Args[0], body, options.ExpectedStatusCode)
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

	res, err := r.service.Patch(in.Args[0], body, options.ExpectedStatusCode)
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
