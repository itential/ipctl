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

type Integration struct {
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	Model              string                 `json:"model"`
	Properties         map[string]interface{} `json:"properties"`
	IsEncrypted        bool                   `json:"isEncrypted"`
	LoggerProperties   map[string]interface{} `json:"loggerProps"`
	Virtual            bool                   `json:"virtual"`
	Profiling          bool                   `json:"profiling"`
	SystemProperties   map[string]interface{} `json:"systemProps"`
	EventDeduplciation map[string]interface{} `json:"eventDeduplication"`
}

type IntegrationService struct {
	client *ServiceClient
}

func NewIntegrationService(iapClient client.Client) *IntegrationService {
	return &IntegrationService{client: NewServiceClient(iapClient)}
}

func NewIntegration(name, integrationType string) Integration {
	logger.Trace()
	return Integration{Name: name, Type: integrationType}
}

func (svc *IntegrationService) Create(in Integration) (*Integration, error) {
	logger.Trace()

	body := map[string]interface{}{
		"properties": map[string]interface{}{
			"name": in.Name,
			"properties": map[string]interface{}{
				"id":   in.Name,
				"type": in.Type,
			},
			"type":    "Adapter",
			"virtual": true,
		},
	}

	type Response struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Data    *Integration `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/integrations",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil

}

func (svc *IntegrationService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/integrations/%s", name))
}

func (svc *IntegrationService) Get(name string) (*Integration, error) {
	logger.Trace()

	type Response struct {
		Data     *Integration           `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var res Response
	var uri = fmt.Sprintf("/integrations/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (svc *IntegrationService) GetAll() ([]Integration, error) {
	logger.Trace()

	type Result struct {
		Data     *Integration           `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	type Response struct {
		Results []Result `json:"results"`
		Total   int      `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/integrations", &res); err != nil {
		return nil, err
	}

	var elements []Integration

	for _, ele := range res.Results {
		elements = append(elements, *ele.Data)
	}

	return elements, nil
}
