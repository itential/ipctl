// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
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
