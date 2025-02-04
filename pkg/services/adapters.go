// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type AdapterOperationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UpdateAdapterResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    Adapter `json:"data"`
}

type CreateAdapterResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    Adapter `json:"data"`
}

type Adapter struct {
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	Model            string                 `json:"model"`
	Properties       map[string]interface{} `json:"properties"`
	IsEncrypted      bool                   `json:"isEncrypted"`
	LoggerProperties map[string]interface{} `json:"loggerProps"`
}

type AdapterService struct {
	client *ServiceClient
}

func NewAdapterService(iapClient client.Client) *AdapterService {
	return &AdapterService{client: NewServiceClient(iapClient)}
}

// GetAll will retrieve all configured adapter instances and return them to the
// calling function as an array of type Adapter.  If there are no configured
// adapters, this function will return an empty array.
func (svc *AdapterService) GetAll() ([]Adapter, error) {
	logger.Trace()

	type Metadata struct {
		ActiveSync bool `json:"activeSync"`
		IsAlive    bool `json:"isAlive"`
	}

	type Results struct {
		Metadata Metadata `json:"metadata"`
		Data     Adapter  `json:"data"`
		Virtual  bool     `json:"virtual"`
	}

	type Collection struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Collection

	if err := svc.client.Get("/adapters", &res); err != nil {
		return nil, err
	}

	var adapters []Adapter

	for _, ele := range res.Results {
		adapters = append(adapters, ele.Data)
	}

	logger.Info("Found %v adapter(s)", res.Total)

	return adapters, nil
}

// Get attempts to retrieve the adapter as specified by the name argument.  If
// the adapter exists, it is returned to the calling function.  If the
// specified adapter does not exist, an error is returned.
func (svc *AdapterService) Get(name string) (*Adapter, error) {
	logger.Trace()

	type Metadata struct {
		ActiveSync bool `json:"activeSync"`
		IsAlive    bool `json:"isAlive"`
	}

	type Response struct {
		Metadata Metadata `json:"metadata"`
		Data     Adapter  `json:"data"`
		Virtual  bool     `json:"virtual"`
	}

	var res Response
	var uri = fmt.Sprintf("/adapters/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (svc *AdapterService) Create(in Adapter) (*CreateAdapterResponse, error) {
	logger.Trace()

	body := map[string]interface{}{"properties": in}

	var res CreateAdapterResponse

	if err := svc.client.Post("/adapters", &body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (svc *AdapterService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/adapters/%s", name))
}

func (svc *AdapterService) Import(in Adapter) (*Adapter, error) {
	logger.Trace()

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Adapter `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/adapters/import",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

func (svc *AdapterService) Export(name string) (*Adapter, error) {
	logger.Trace()
	return svc.Get(name)
}

func (svc *AdapterService) Start(name string) error {
	logger.Trace()
	return svc.client.Put(fmt.Sprintf("/adapters/%s/start", name), nil, nil)
}

func (svc *AdapterService) Stop(name string) error {
	logger.Trace()
	return svc.client.Put(fmt.Sprintf("/adapters/%s/stop", name), nil, nil)
}

func (svc *AdapterService) Restart(name string) error {
	logger.Trace()

	if err := svc.Stop(name); err != nil {
		return err
	}

	return svc.Start(name)
	//return svc.client.Put(fmt.Sprintf("/adapters/%s/restart", name), nil, nil)
}
