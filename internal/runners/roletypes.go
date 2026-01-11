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

type RoleTypeRunner struct {
	service *services.RoleService
	BaseRunner
	client client.Client
}

func NewRoleTypeRunner(client client.Client, cfg *config.Config) *RoleTypeRunner {
	return &RoleTypeRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		client:     client,
		service:    services.NewRoleService(client),
	}
}

// Get implements the `get role-types` command
func (r *RoleTypeRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "provenance"},
		Object: roles,
	}, nil

}

func (r *RoleTypeRunner) Describe(in Request) (*Response, error) {
	return notImplemented(in)
}
