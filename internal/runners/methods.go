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

type MethodRunner struct {
	service *services.MethodService
	config  *config.Config
}

func NewMethodRunner(client client.Client, cfg *config.Config) *MethodRunner {
	return &MethodRunner{
		service: services.NewMethodService(client),
	}
}

// Get implements the `get methods` command
func (r *MethodRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	methods, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "provenance"},
		Object: methods,
	}, nil
}

func (r *MethodRunner) Describe(in Request) (*Response, error) {
	return notImplemented(in)
}
