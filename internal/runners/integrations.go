// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type IntegrationRunner struct {
	config  *config.Config
	service *services.IntegrationService
}

func NewIntegrationRunner(c client.Client, cfg *config.Config) *IntegrationRunner {
	return &IntegrationRunner{
		service: services.NewIntegrationService(c),
		config:  cfg,
	}
}

// Get implements the `get integration-models` command
func (r *IntegrationRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tMODEL"}
	for _, ele := range res {
		model := strings.Replace(ele.Model, "@itential/adapter_", "", -1)
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, model))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(res),
	), nil

}

// Describe implements the `describe integration-model <name>` command
func (r *IntegrationRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", res.Name),
		WithJson(res),
	), nil
}

// Create implements the `create integration <name>` command
func (r *IntegrationRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.IntegrationCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.Get(name)

		if existing != nil {
			if err := r.service.Delete(existing.Name); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "integration not found" {
				return nil, err
			}
		}
	}

	res, err := r.service.Create(services.NewIntegration(name, options.Model))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created integration `%s`", name),
		WithJson(res),
	), nil
}

// Delete implements the `delete integration <name>` command
func (r *IntegrationRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Delete(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted integration `%s`", name),
	), nil
}

// Clear implements the `clear integrations` command
func (r *IntegrationRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range elements {
		if err := r.service.Delete(ele.Name); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Deleted %v integration(s)", len(elements)),
	), nil
}

// Copy implements the `copy integration <name>` command
func (r *IntegrationRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common *flags.AssetCopyCommon
	utils.LoadObject(in.Common, &common)

	if common.From == common.To {
		return nil, errors.New("source and destination servers must be different values")
	}

	fromClient, cancel, err := NewClient(common.From, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	fromService := services.NewIntegrationService(fromClient)

	res, err := fromService.Get(name)
	if err != nil {
		return nil, err
	}

	toClient, cancel, err := NewClient(common.To, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	toService := services.NewIntegrationService(toClient)

	if _, err := toService.Get(name); err != nil {
		if err.Error() != "integration not found" {
			return nil, errors.New(fmt.Sprintf("integration `%s` already exists on the destination server", name))
		}
	}

	_, err = toService.Create(*res)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied integration `%s` from `%s` to `%s`", name, common.From, common.To),
	), nil
}
