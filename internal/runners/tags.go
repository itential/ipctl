// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"strings"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type TagRunner struct {
	config  *config.Config
	service *services.TagService
}

func NewTagRunner(client client.Client, cfg *config.Config) *TagRunner {
	return &TagRunner{
		config:  cfg,
		service: services.NewTagService(client),
	}
}

// Get implements the `get tags` command
func (r *TagRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	tags, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range tags {
		lines := []string{ele.Name, ele.Description}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(tags),
	), nil
}

// Describe implements the `describe tag <name>` command
func (r *TagRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	tag, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		"",
		WithJson(tag),
	), nil
}
