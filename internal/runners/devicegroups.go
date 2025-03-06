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

type DeviceGroupRunner struct {
	config  *config.Config
	service *services.DeviceGroupService
}

func NewDeviceGroupRunner(client client.Client, cfg *config.Config) *DeviceGroupRunner {
	return &DeviceGroupRunner{
		config:  cfg,
		service: services.NewDeviceGroupService(client),
	}
}

// Get is the implementation of the command `get devices`
func (r *DeviceGroupRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	devices, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range devices {
		lines := []string{ele.Name, ele.Description}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithObject(devices),
	), nil

}

func (r *DeviceGroupRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", res.Name),
		WithObject(res),
	), nil
}
