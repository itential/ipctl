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

type IntegrationModelTls struct {
	Enabled             bool `json:"enabled"`
	RefjectUnauthorized bool `json:"rejectUnauthroized"`
}

type IntegrationModelServer struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	BasePath string `json:"base_path"`
}

type IntegrationModelProperties struct {
	Authentication map[string]interface{} `json:"authentication"`
	Server         IntegrationModelServer `json:"server"`
	Tls            IntegrationModelTls    `json:"tls"`
	Version        string                 `json:"version"`
}

type IntegrationModel struct {
	Model       string                     `json:"model,omitempty"`
	VersionId   string                     `json:"versionId"`
	Description string                     `json:"description"`
	Properties  IntegrationModelProperties `json:"properties"`
}

type IntegrationModelService struct {
	client *ServiceClient
}

func NewIntegrationModelService(iapClient client.Client) *IntegrationModelService {
	return &IntegrationModelService{client: NewServiceClient(iapClient)}
}

func (svc *IntegrationModelService) GetAll() ([]IntegrationModel, error) {
	logger.Trace()

	type Response struct {
		IntegrationModels []IntegrationModel `json:"integrationModels"`
		Total             int                `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/integration-models", &res); err != nil {
		return nil, err
	}

	return res.IntegrationModels, nil
}

func (svc *IntegrationModelService) Get(name string) (*IntegrationModel, error) {
	logger.Trace()

	var res *IntegrationModel
	var uri = fmt.Sprintf("/integration-models/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *IntegrationModelService) Create(in map[string]interface{}) (*IntegrationModel, error) {
	logger.Trace()

	body := map[string]interface{}{"model": in}

	type Response struct {
		Message string                 `json:"message"`
		Status  string                 `json:"status"`
		Data    map[string]interface{} `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/integration-models",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	model, err := svc.Get(res.Data["versionId"].(string))
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (svc *IntegrationModelService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/integration-models/%s", name))
}

func (svc *IntegrationModelService) Export(name string) (map[string]interface{}, error) {
	logger.Trace()

	var res map[string]interface{}
	var uri = fmt.Sprintf("/integration-models/%s/export", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}
