// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"net/http"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

// BaseService provides common HTTP operations for all service types.
// It embeds the client.Client directly to avoid unnecessary wrapper indirection.
type BaseService struct {
	client client.Client
}

// NewBaseService creates a new BaseService with the provided client.
func NewBaseService(c client.Client) BaseService {
	return BaseService{client: c}
}

// Do executes an HTTP request and unmarshals the response into ptr if provided.
func (svc *BaseService) Do(req *Request, ptr any) error {
	logging.Trace()

	req.client = svc.client

	res, err := Do(req)
	if err != nil {
		return err
	}

	if ptr != nil {
		if err := json.Unmarshal(res.Body, &ptr); err != nil {
			return err
		}
	}

	return nil
}

// Get performs a GET request to the specified URI.
func (svc *BaseService) Get(uri string, ptr any) error {
	logging.Trace()
	return svc.GetRequest(&Request{
		uri: uri,
	}, ptr)
}

// GetRequest performs a GET request with the provided request configuration.
func (svc *BaseService) GetRequest(in *Request, ptr any) error {
	logging.Trace()
	in.method = http.MethodGet
	return svc.Do(in, ptr)
}

// Post performs a POST request to the specified URI with the given body.
func (svc *BaseService) Post(uri string, body any, ptr any) error {
	logging.Trace()
	return svc.PostRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}

// PostRequest performs a POST request with the provided request configuration.
func (svc *BaseService) PostRequest(in *Request, ptr any) error {
	logging.Trace()
	in.method = http.MethodPost
	return svc.Do(in, ptr)
}

// Delete performs a DELETE request to the specified URI.
func (svc *BaseService) Delete(uri string) error {
	logging.Trace()
	return svc.DeleteRequest(&Request{
		uri: uri,
	}, nil)
}

// DeleteRequest performs a DELETE request with the provided request configuration.
func (svc *BaseService) DeleteRequest(in *Request, ptr any) error {
	logging.Trace()
	in.method = http.MethodDelete
	return svc.Do(in, ptr)
}

// Patch performs a PATCH request to the specified URI with the given body.
func (svc *BaseService) Patch(uri string, body any, ptr any) error {
	return svc.PatchRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}

// PatchRequest performs a PATCH request with the provided request configuration.
func (svc *BaseService) PatchRequest(in *Request, ptr any) error {
	logging.Trace()
	in.method = http.MethodPatch
	return svc.Do(in, ptr)
}

// PutRequest performs a PUT request with the provided request configuration.
func (svc *BaseService) PutRequest(in *Request, ptr any) error {
	logging.Trace()
	in.method = http.MethodPut
	return svc.Do(in, ptr)
}

// Put performs a PUT request to the specified URI with the given body.
func (svc *BaseService) Put(uri string, body any, ptr any) error {
	logging.Trace()
	return svc.PutRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}

// DEPRECATED: ServiceClient is deprecated, use BaseService instead.
// Kept for backward compatibility during migration.
type ServiceClient = BaseService

// DEPRECATED: NewServiceClient is deprecated, use NewBaseService instead.
// Kept for backward compatibility during migration.
func NewServiceClient(c client.Client) *ServiceClient {
	bs := NewBaseService(c)
	return &bs
}
