// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type ServerRunner struct {
	BaseRunner
	service *services.HealthService
}

func NewServerRunner(c client.Client, cfg config.Provider) *ServerRunner {
	return &ServerRunner{
		service:    services.NewHealthService(c),
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

func (r *ServerRunner) Inspect(in Request) (*Response, error) {
	logging.Trace()

	status, err := r.service.GetStatus()
	if err != nil {
		return nil, err
	}

	system, err := r.service.GetSystemHealth()
	if err != nil {
		return nil, err
	}

	server, err := r.service.GetServerHealth()
	if err != nil {
		return nil, err
	}

	app, err := r.service.GetApplicationHealth()
	if err != nil {
		return nil, err
	}

	adapter, err := r.service.GetAdapterHealth()
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{
		"status":      status,
		"system":      system,
		"server":      server,
		"application": app,
		"adapters":    adapter,
	}

	return &Response{
		Object: res,
	}, nil
}
