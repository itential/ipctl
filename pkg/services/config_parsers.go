// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ConfigurationParser struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Template string     `json:"template"`
	LexRules [][]string `json:"lexRules"`
	Updated  any        `json:"updated"`
	Created  any        `json:"created"`
	Gbac     Gbac       `json:"gbac"`
}

type ConfigurationParserService struct {
	client *ServiceClient
}

func NewConfigurationParser(name string) ConfigurationParser {
	logger.Trace()
	return ConfigurationParser{Name: name}
}

func NewConfigurationParserService(iapClient client.Client) *ConfigurationParserService {
	return &ConfigurationParserService{client: NewServiceClient(iapClient)}
}

func (svc *ConfigurationParserService) GetAll() ([]ConfigurationParser, error) {
	logger.Trace()

	type Response struct {
		List  []ConfigurationParser `json:"list"`
		Total int                   `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/configuration_manager/configurations/parser", &res); err != nil {
		return nil, err
	}

	return res.List, nil
}
