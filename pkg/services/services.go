// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

func tobytes(obj any) ([]byte, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}
	return b, nil
}

type Request struct {
	client                 client.Client
	method                 string
	uri                    string
	params                 Params
	query                  map[string]string
	body                   any
	response               any
	expectedStatusCode     int
	disableStatusCodeCheck bool
}

func Do(r *Request) (*client.Response, error) {
	logging.Trace()

	req := &client.Request{
		Path: r.uri,
	}

	if r.query != nil {
		req.Params = r.query
	} else if r.params != nil {
		req.Params = r.params.Query()
	}

	if r.body != nil {
		body, err := tobytes(r.body)
		if err != nil {
			return nil, err
		}
		req.Body = body
	}

	var resp *client.Response
	var err error

	switch r.method {
	case http.MethodGet:
		if r.expectedStatusCode == 0 {
			r.expectedStatusCode = http.StatusOK
		}
		resp, err = r.client.Get(req)
	case http.MethodPost:
		if r.expectedStatusCode == 0 {
			r.expectedStatusCode = http.StatusCreated
		}
		resp, err = r.client.Post(req)
	case http.MethodPut:
		if r.expectedStatusCode == 0 {
			r.expectedStatusCode = http.StatusOK
		}
		resp, err = r.client.Put(req)
	case http.MethodPatch:
		if r.expectedStatusCode == 0 {
			r.expectedStatusCode = http.StatusOK
		}
		resp, err = r.client.Patch(req)
	case http.MethodDelete:
		if r.expectedStatusCode == 0 {
			r.expectedStatusCode = http.StatusOK
		}
		resp, err = r.client.Delete(req)
	}

	if err != nil {
		return resp, err
	}

	if err := checkResponseForError(r, resp); err != nil {
		return nil, err
	}

	if resp.Body != nil {
		logging.Debug("%s", string(resp.Body))
	}

	if r.response != nil {
		if err := json.Unmarshal(resp.Body, r.response); err != nil {
			logging.Error(err, "failed to unmarshal reponse")
			return resp, err
		}
	}

	return resp, err
}

func checkResponseForError(r *Request, resp *client.Response) error {
	if !r.disableStatusCodeCheck {
		if r.expectedStatusCode != 0 && r.expectedStatusCode != resp.StatusCode {
			logging.Error(nil, "status code = %v, expected status code = %v", resp.StatusCode, r.expectedStatusCode)
			logging.Error(nil, "%s", string(resp.Body))
			if resp != nil {
				return errors.New(string(resp.Body))
			} else {
				return errors.New(
					fmt.Sprintf(
						"expected status code %v, got %v", r.expectedStatusCode, resp.StatusCode,
					),
				)
			}
		}

		if resp.StatusCode > 299 {
			var errMsg ErrorMessage

			var body map[string]interface{}
			json.Unmarshal(resp.Body, &body)
			b, _ := json.MarshalIndent(body, "", "    ")
			logging.Debug("\n%s", string(b))

			logging.Error(errors.New(resp.Status), "%s", errMsg.Message)
		}
	}

	if value, exists := resp.Headers["Content-Type"]; exists {
		parts := strings.Split(value, ";")
		if parts[0] != "application/json" {
			logging.Warn("expected response header Content-Type=application/json, got %s", parts[0])
		}
	} else {
		logging.Warn("missing Content-Type response header")
	}

	return nil
}
