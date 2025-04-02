// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ConfigManagerJinja2Template struct {
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables"`
	Options   map[string]interface{} `json:"options"`
}

type ConfigManagerService struct {
	client *ServiceClient
}

func NewConfigManagerService(c client.Client) *ConfigManagerService {
	return &ConfigManagerService{client: NewServiceClient(c)}
}

func (svc *ConfigManagerService) Render(in ConfigManagerJinja2Template) {
	logger.Trace()

	if err := svc.client.Post("/configuration_manager/jinja2", in, nil); err != nil {
		panic(err)
	}

}
