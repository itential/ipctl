// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"
	"strings"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type DeviceRunner struct {
	config  *config.Config
	service *services.DeviceService
}

func NewDeviceRunner(client client.Client, cfg *config.Config) *DeviceRunner {
	return &DeviceRunner{
		config:  cfg,
		service: services.NewDeviceService(client),
	}
}

// Get is the implementation of the command `get devices`
func (r *DeviceRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	devices, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tOSTYPE"}
	for _, ele := range devices {
		lines := []string{ele.Name, ele.OsType}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(devices),
	), nil

}

func (r *DeviceRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", res.Name),
		WithUrl(fmt.Sprintf("/configuration_manager/#/devices/%s", res.Name)),
		WithJson(res),
	), nil
}
