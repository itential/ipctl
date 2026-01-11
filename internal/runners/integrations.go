// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type IntegrationRunner struct {
	BaseRunner
	client  client.Client
	service *services.IntegrationService
}

func NewIntegrationRunner(c client.Client, cfg config.Provider) *IntegrationRunner {
	return &IntegrationRunner{
		service:    services.NewIntegrationService(c),
		client:     c,
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

/*
******************************************************************************
Reader interface
******************************************************************************
*/

// Get implements the `get integration-models` command
func (r *IntegrationRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "model"},
		Object: res,
	}, nil

}

// Describe implements the `describe integration-model <name>` command
func (r *IntegrationRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s", res.Name),
		Object: res,
	}, nil
}

/*
******************************************************************************
Writer interface
******************************************************************************
*/

// Create implements the `create integration <name>` command
func (r *IntegrationRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	options := in.Options.(*flags.IntegrationCreateOptions)

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

	return &Response{
		Text:   fmt.Sprintf("Successfully created integration `%s`", name),
		Object: res,
	}, nil
}

// Delete implements the `delete integration <name>` command
func (r *IntegrationRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	if err := r.service.Delete(name); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted integration `%s`", name),
	}, nil
}

// Clear implements the `clear integrations` command
func (r *IntegrationRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range elements {
		if err := r.service.Delete(ele.Name); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v integration(s)", len(elements)),
	}, nil
}

/*
******************************************************************************
Copier interface
******************************************************************************
*/

// Copy implements the `copy integration <name>` command
func (r *IntegrationRunner) Copy(in Request) (*Response, error) {
	logging.Trace()

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

	return &Response{
		Text: fmt.Sprintf("Successfully copied integration `%s` from `%s` to `%s`", name, common.From, common.To),
	}, nil
}

/*
*******************************************************************************
Exporter interface
*******************************************************************************
*/

// Export implements the `export integration ...` command
func (r *IntegrationRunner) Export(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.integration.json", name)

	if err := exportAssetFromRequest(in, res, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully exported integration `%s`", res.Name),
		Object: res,
	}, nil
}

/*
*******************************************************************************
Importer interface
*******************************************************************************
*/

// Import implements the `import integration ...` command
func (r *IntegrationRunner) Import(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var integration services.Integration

	if err := importUnmarshalFromRequest(in, &integration); err != nil {
		return nil, err
	}

	existing, err := r.service.Get(integration.Name)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "does not exist.\"") {
			return nil, errors.New(
				fmt.Sprintf("integration `%s` already exists", integration.Name),
			)
		}
	}

	if existing != nil {
		if common.Replace {
			if err := r.service.Delete(integration.Name); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("integration `%s` already exists, use `--replace` to overwrite it", integration.Name)
		}
	}

	modelName := strings.Split(integration.Model, "adapter_")[1]
	if err := r.checkIfModelExists(modelName); err != nil {
		return nil, err
	}

	res, err := r.service.Create(integration)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully imported integration `%s`", res.Name),
		Object: res,
	}, nil
}

func (r *IntegrationRunner) checkIfModelExists(name string) error {
	logging.Trace()
	_, err := services.NewIntegrationModelService(r.client).Get(name)
	if err != nil {
		return err
	}
	return nil
}
