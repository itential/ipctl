// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/services"
)

type ViewRunner struct {
	service *services.ViewService
	BaseRunner
}

func NewViewRunner(client client.Client, cfg *config.Config) *ViewRunner {
	return &ViewRunner{
		service: services.NewViewService(client),
	}
}

// Get implements the `get views` command
func (r *ViewRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	views, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"path", "provenance"},
		Object: views,
	}, nil
}

func (r *ViewRunner) Describe(in Request) (*Response, error) {
	return notImplemented(in)
}
