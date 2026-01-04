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

type IntegrationModelRunner struct {
	BaseRunner
	service *services.IntegrationModelService
}

func NewIntegrationModelRunner(c client.Client, cfg *config.Config) *IntegrationModelRunner {
	return &IntegrationModelRunner{
		service: services.NewIntegrationModelService(c),
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

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

	return &Response{
		Keys:   []string{"model"},
		Object: res,
	}, nil

}

// Describe implements the `describe integration-model <name>` command
func (r *IntegrationModelRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s", res.Model),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

func (r *IntegrationModelRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	var schema map[string]interface{}

	if err := importUnmarshalFromRequest(in, &schema); err != nil {
		return nil, err
	}

	res, err := r.service.Create(schema)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created integration model `%s`", res.Model),
		Object: res,
	}, nil

}

func (r *IntegrationModelRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Delete(in.Args[0]); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted integration model `%s`", name),
	}, nil
}

func (r *IntegrationModelRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range res {
		if err := r.service.Delete(ele.Model); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Suuccessfully deleted %v integration models", len(res)),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the `import integration-model ...` command
func (r *IntegrationModelRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var schema map[string]interface{}

	if err := importUnmarshalFromRequest(in, &schema); err != nil {
		return nil, err
	}

	res, err := r.service.Create(schema)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully imported integration model `%s`", res.VersionId),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export implements the  `export integration-model ...` command
func (r *IntegrationModelRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	//common := in.Common.(*flags.AssetExportCommon)

	name := in.Args[0]

	res, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.integration_model.json", name)

	if err := exportAssetFromRequest(in, res, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported integration_model `%s`", name),
	}, nil
}
