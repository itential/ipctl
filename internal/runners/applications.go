// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type ApplicationRunner struct {
	config  *config.Config
	client  client.Client
	service *services.ApplicationService
}

func NewApplicationRunner(c client.Client, cfg *config.Config) *ApplicationRunner {
	return &ApplicationRunner{
		service: services.NewApplicationService(c),
		config:  cfg,
		client:  c,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Describe implements the `get applications` command
func (r *ApplicationRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var display = []string{"NAME\tMODEL"}

	for _, ele := range res {
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Model))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(res),
	), nil

}

// Describe implements the `describe applications ...` command
func (r *ApplicationRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.Get(in.Args[0])
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		logger.Fatal(err, "failed to marshal data")
	}

	return NewResponse(
		string(b),
		WithJson(res),
	), nil
}

func (r *ApplicationRunner) Start(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Start(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully started application `%s`", name),
	), nil
}

func (r *ApplicationRunner) Stop(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Stop(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully stopped application `%s`", name),
	), nil
}

func (r *ApplicationRunner) Restart(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Restart(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully restarted application `%s`", name),
	), nil
}

func (r *ApplicationRunner) Inspect(in Request) (*Response, error) {
	logger.Trace()

	svc := services.NewHealthService(r.client)
	res, err := svc.GetApplicationHealth()
	if err != nil {
		return nil, err
	}

	var display = []string{"NAME\tSTATUS\tVERSION"}

	for _, ele := range res {
		display = append(display, fmt.Sprintf(
			"%s\t%s\t%s", ele.Id, ele.State, ele.Version,
		))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(res),
	), nil
}
