// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type DeviceGroupRunner struct {
	BaseRunner
	service *services.DeviceGroupService
}

func NewDeviceGroupRunner(client client.Client, cfg config.Provider) *DeviceGroupRunner {
	return &DeviceGroupRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		service:    services.NewDeviceGroupService(client),
	}
}

/*
*******************************************************************************
Reader interface
*******************************************************************************
*/

// Get implements the `get devices ...` command
func (r *DeviceGroupRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	devices, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range devices {
		lines := []string{ele.Name, ele.Description}
		display = append(display, strings.Join(lines, "\t"))
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: devices,
	}, nil

}

// Describe implements the `describe device-group ...` command
func (r *DeviceGroupRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s", res.Name),
		Object: res,
	}, nil
}

/*
*******************************************************************************
Writer interface
*******************************************************************************
*/

// Create implements the `create device-group ...` command
func (r *DeviceGroupRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	options := in.Options.(*flags.DeviceGroupCreateOptions)

	name := in.Args[0]

	res, err := r.service.Create(services.NewDeviceGroup(name, options.Description))
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created new device-group `%s` (%s)", res.Name, res.Id),
		Object: res,
	}, nil
}

// Delete implements the `delete device-group ...` command
func (r *DeviceGroupRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	deviceGroup, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(deviceGroup.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully deleted device-group `%s` (%s)", deviceGroup.Name, deviceGroup.Id),
		Object: deviceGroup,
	}, nil
}

// Clear implements the `clear device-group ...` command
func (r *DeviceGroupRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range groups {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v device-groups", len(groups)),
	}, nil
}
