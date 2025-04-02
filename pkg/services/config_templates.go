// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ConfigTemplate struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Template      string                 `json:"template"`
	Variables     map[string]interface{} `json:"variables"`
	Created       string                 `json:"created"`
	Updated       string                 `json:"updated"`
	CreatedBy     string                 `json:"createdBy"`
	UpdatedBy     string                 `json:"updatedBy"`
	Gbac          Gbac                   `json:"gbac"`
	DeviceOsTypes []string               `json:"deviceOSTypes"`
}

type ConfigTemplateService struct {
	client *ServiceClient
}

func NewConfigTemplateService(c client.Client) *ConfigTemplateService {
	return &ConfigTemplateService{client: NewServiceClient(c)}
}

func NewConfigTemplate(name string) ConfigTemplate {
	return ConfigTemplate{Name: name}
}

func (svc *ConfigTemplateService) Create(in ConfigTemplate) (*ConfigTemplate, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":     in.Name,
		"template": in.Template,
	}

	type Response struct {
		Result string          `json:"result"`
		Data   *ConfigTemplate `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/configuration_manager/templates",
		body:               body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
