// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"sort"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type RoleTypeRunner struct {
	service *services.RoleService
	config  *config.Config
	client  client.Client
}

func NewRoleTypeRunner(client client.Client, cfg *config.Config) *RoleTypeRunner {
	return &RoleTypeRunner{
		config:  cfg,
		client:  client,
		service: services.NewRoleService(client),
	}
}

// Get implements the `get role-types` command
func (r *RoleTypeRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var display []string

	for _, ele := range roles {
		if !utils.StringExists(display, ele.Provenance) {
			display = append(display, ele.Provenance)
		}
	}

	sort.Strings(display)

	utils.StringInsert(display, 0, "TYPE")

	return NewResponse(
		"",
		WithTable(display),
	), nil

}

func (r *RoleTypeRunner) Describe(in Request) (*Response, error) {
	return NotImplemented(in)
}
