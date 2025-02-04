// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type IntegrationModelRunner struct {
	config  *config.Config
	service *services.IntegrationModelService
}

func NewIntegrationModelRunner(c client.Client, cfg *config.Config) *IntegrationModelRunner {
	return &IntegrationModelRunner{
		service: services.NewIntegrationModelService(c),
		config:  cfg,
	}
}

// Get implements the `get integration-models` command
func (r *IntegrationModelRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"MODEL"}
	for _, ele := range res {
		display = append(display, ele.VersionId)
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(res),
	), nil

}

// Describe implements the `describe integration-model <name>` command
func (r *IntegrationModelRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", res.Model),
		WithJson(res),
	), nil
}

// Import implements the `import integration-model ...` command
func (r *IntegrationModelRunner) Import(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)

	/*
		path, err := NormalizePath(in)
		if err != nil {
			return nil, err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var integration_model services.IntegrationModel
		utils.UnmarshalData(data, &integration_model)

		p, err := r.service.Get(integration_model.Name)
		if err == nil {
			if common.Force {
				r.service.Delete([]string{p.Id})
			} else {
				return nil, errors.New(fmt.Sprintf("integration_model with name `%s` already exists", p.Name))
			}
		}

		err = r.service.Import(integration_model)
		if err != nil {
			return nil, err
		}

		return NewResponse(
			fmt.Sprintf("Successfully imported integration_model `%s`", integration_model.Name),
		), nil
	*/
}

// Export implements the  `export integration-model ...` command
func (r *IntegrationModelRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options *flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	res, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.integration_model.json", name)

	if err := utils.WriteJsonToDisk(res, fn, options.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported integration_model `%s`", name),
	), nil
}
