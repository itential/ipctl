// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type ActionRunner struct {
	config_manager *services.ConfigManagerService
}

func NewActionRunner(client client.Client) *ActionRunner {
	return &ActionRunner{
		config_manager: services.NewConfigManagerService(client),
	}
}

func (r *ActionRunner) Render(in Request) (*Response, error) {
	logger.Trace()

	r.config_manager.Render(services.ConfigManagerJinja2Template{
		Template:  "hostname {{ foo }}",
		Variables: map[string]interface{}{"foo": "bar"},
		Options:   map[string]interface{}{},
	})

	return notImplemented(in)
}
