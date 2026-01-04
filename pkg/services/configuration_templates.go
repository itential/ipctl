// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"

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
	BaseService
}

func NewConfigurationTemplateService(c client.Client) *ConfigurationTemplateService {
	return &ConfigurationTemplateService{BaseService: NewBaseService(c)}
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
	if err := svc.Post(uri, &body, &res); err != nil {
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
	if err := svc.Post(uri, &body, &res); err != nil {
		return nil, err
	}
	return &res.List[0], nil
}

// GetByName retrieves a configuration template by name using client-side filtering.
// DEPRECATED: Business logic method - prefer using resources.ConfigurationTemplateResource.GetByName
func (svc *ConfigurationTemplateService) GetByName(name string) (*ConfigurationTemplate, error) {
	logger.Trace()

	templates, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		if template.Name == name {
			return svc.Get(template.Id)
		}
	}

	return nil, errors.New("configuration template not found")
}
