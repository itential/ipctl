// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"strings"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type ConfigurationParserRunner struct {
	config  *config.Config
	service *services.ConfigurationParserService
}

func NewConfigurationParserRunner(client client.Client, cfg *config.Config) *ConfigurationParserRunner {
	return &ConfigurationParserRunner{
		config:  cfg,
		service: services.NewConfigurationParserService(client),
	}
}

// Get is the implementation of the command `get devices`
func (r *ConfigurationParserRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	devices, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME"}
	for _, ele := range devices {
		lines := []string{ele.Name}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithObject(devices),
	), nil

}

func (r *ConfigurationParserRunner) Describe(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
	/*

		name := in.Args[0]

		res, err := r.service.Get(name)
		if err != nil {
			return nil, err
		}

		return NewResponse(
			fmt.Sprintf("Name: %s", res.Name),
			WithUrl(fmt.Sprintf("/configuration_manager/#/devices/%s", res.Name)),
			WithObject(res),
		), nil
	*/
}
