// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type AdapterModelRunner struct {
	config  *config.Config
	service *services.AdapterModelService
}

func NewAdapterModelRunner(c client.Client, cfg *config.Config) *AdapterModelRunner {
	return &AdapterModelRunner{
		service: services.NewAdapterModelService(c),
		config:  cfg,
	}
}

// Get implements the `get adapter-models` command
func (r *AdapterModelRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"MODEL"}
	for _, ele := range res {
		display = append(display, ele)
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(res),
	), nil

}

func (r *AdapterModelRunner) Describe(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
}
