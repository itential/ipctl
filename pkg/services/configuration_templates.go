// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ConfigurationTemplateCollection struct {
	List  []ConfigurationTemplate `json:"list"`
	Total int                     `json:"total"`
}

type ConfigurationTemplate struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Template      string                 `json:"template"`
	Variables     map[string]interface{} `json:"variables"`
	Created       string                 `json:"created"`
	CreatedBy     string                 `json:"createdBy"`
	Updated       string                 `json:"updated"`
	UpdatedBy     string                 `json:"updatedBy"`
	Gbac          Gbac                   `json:"gbac"`
	DeviceOsTypes []string               `json:"deviceOSTypes"`
}

type ConfigurationTemplateService struct {
	client *ServiceClient
}

func NewConfigurationTemplateService(c client.Client) *ConfigurationTemplateService {
	return &ConfigurationTemplateService{client: NewServiceClient(c)}
}

func (svc *ConfigurationTemplateService) GetAll() ([]ConfigurationTemplate, error) {
	logger.Trace()
	// FIXME (privateip) need to implement full paging
	body := map[string]interface{}{
		"name": "",
		"options": map[string]interface{}{
			"limit": 100,
			"sort": map[string]interface{}{
				"name": 1,
			},
			"start": 0,
		},
	}
	var res ConfigurationTemplateCollection
	var uri = "/configuration_manager/templates/search"
	if err := svc.client.Post(uri, &body, &res); err != nil {
		return nil, err
	}
	return res.List, nil
}

func (svc *ConfigurationTemplateService) Get(id string) (*ConfigurationTemplate, error) {
	logger.Trace()
	body := map[string]interface{}{
		"name":    id,
		"options": make(map[string]interface{}),
	}
	var res ConfigurationTemplateCollection
	var uri = "/configuration_manager/templates/search"
	if err := svc.client.Post(uri, &body, &res); err != nil {
		return nil, err
	}
	return &res.List[0], nil
}

func (svc *ConfigurationTemplateService) GetByName(name string) (*ConfigurationTemplate, error) {
	logger.Trace()
	var res *ConfigurationTemplate
	var err error
	elements, err := svc.GetAll()
	if err != nil {
		return nil, err
	}
	for _, ele := range elements {
		if ele.Name == name {
			res, err = svc.Get(ele.Id)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return res, nil
}
