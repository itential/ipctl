// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ServiceClient struct {
	client client.Client
}

func NewServiceClient(c client.Client) *ServiceClient {
	return &ServiceClient{client: c}
}

func (http *ServiceClient) Do(req *Request, ptr any) error {
	logger.Trace()

	req.client = http.client

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

func (svc *ServiceClient) Get(uri string, ptr any) error {
	logger.Trace()
	return svc.GetRequest(&Request{
		uri: uri,
	}, ptr)
}

func (svc *ServiceClient) GetRequest(in *Request, ptr any) error {
	logger.Trace()
	in.method = http.MethodGet
	return svc.Do(in, ptr)
}

func (svc *ServiceClient) Post(uri string, body any, ptr any) error {
	logger.Trace()
	return svc.PostRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}

func (svc *ServiceClient) PostRequest(in *Request, ptr any) error {
	logger.Trace()
	in.method = http.MethodPost
	return svc.Do(in, ptr)
}

func (svc *ServiceClient) Delete(uri string) error {
	logger.Trace()
	return svc.DeleteRequest(&Request{
		uri: uri,
	}, nil)
}

func (svc *ServiceClient) DeleteRequest(in *Request, ptr any) error {
	logger.Trace()
	in.method = http.MethodDelete
	return svc.Do(in, ptr)
}

func (svc *ServiceClient) Patch(uri string, body any, ptr any) error {
	return svc.PatchRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}

func (svc *ServiceClient) PatchRequest(in *Request, ptr any) error {
	logger.Trace()
	in.method = http.MethodPatch
	return svc.Do(in, ptr)
}

func (svc *ServiceClient) PutRequest(in *Request, ptr any) error {
	logger.Trace()
	in.method = http.MethodPut
	return svc.Do(in, ptr)
}

func (svc *ServiceClient) Put(uri string, body any, ptr any) error {
	logger.Trace()
	return svc.PutRequest(&Request{
		uri:  uri,
		body: body,
	}, ptr)
}
