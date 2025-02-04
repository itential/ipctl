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

type ViewRunner struct {
	service *services.ViewService
	config  *config.Config
}

func NewViewRunner(client client.Client, cfg *config.Config) *ViewRunner {
	return &ViewRunner{
		service: services.NewViewService(client),
	}
}

// Get implements the `get views` command
func (r *ViewRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	views, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"PATH\tTYPE"}
	for _, ele := range views {
		line := fmt.Sprintf("%s\t%s", ele.Path, ele.Provenance)
		display = append(display, line)
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(views),
	), nil
}

func (r *ViewRunner) Describe(in Request) (*Response, error) {
	return NotImplemented(in)
}
