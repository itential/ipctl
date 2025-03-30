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

type AnalyticTemplateRunner struct {
	service *services.AnalyticTemplateService
	config  *config.Config
}

func NewAnalyticTemplateRunner(c client.Client, cfg *config.Config) *AnalyticTemplateRunner {
	return &AnalyticTemplateRunner{
		service: services.NewAnalyticTemplateService(c),
		config:  cfg,
	}
}

/*
******************************************************************************
Reader interface
******************************************************************************
*/

// Get implements the `get command-templates` command
func (r *AnalyticTemplateRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	options := in.Options.(*flags.AnalyticTemplateGetOptions)

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var templates []services.AnalyticTemplate
	for _, ele := range res {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			templates = append(templates, ele)
		} else if !strings.HasPrefix(ele.Name, "@") {
			templates = append(templates, ele)
		}
	}

	return &Response{
		Keys:   []string{"name"},
		Object: templates,
	}, nil

}

// Describe implements the `describe command-template <name>` command
func (r *AnalyticTemplateRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	template, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s", template.Name),
		Object: template,
	}, nil
}

/*
******************************************************************************
Writer Interface
******************************************************************************
*/

func (r *AnalyticTemplateRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.AnalyticTemplateCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.Get(name)

		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "analytic template not found" {
				return nil, err
			}
		}
	}

	res, err := r.service.Create(services.NewAnalyticTemplate(name))
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created analytic template `%s`", name),
		Object: res,
	}, nil
}

func (r *AnalyticTemplateRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	template, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(template.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted analytic-template `%s`", name),
	}, nil
}

func (r *AnalyticTemplateRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range elements {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v analytic-template(s)", len(elements)),
	}, nil
}

/*
******************************************************************************
Importer Interface
******************************************************************************
*/

func (r *AnalyticTemplateRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var ct services.AnalyticTemplate

	if err := importUnmarshalFromRequest(in, &ct); err != nil {
		return nil, err
	}

	if err := r.importAnalyticTemplate(ct, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported analytic template `%s`", ct.Name),
	}, nil
}

/*
******************************************************************************
Exporter Interface
******************************************************************************
*/

func (r *AnalyticTemplateRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	ct, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.analytic_template.json", name)

	if err := exportAssetFromRequest(in, ct, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported analytic template `%s`", ct.Name),
	}, nil
}

/*
******************************************************************************
Copier Interface
******************************************************************************
*/

func (r *AnalyticTemplateRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "analytic-template"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied analytic template `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *AnalyticTemplateRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewAnalyticTemplateService(client).Export(name)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (r *AnalyticTemplateRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewAnalyticTemplateService(client)

	name := in.(services.AnalyticTemplate).Name

	if exists, err := svc.Get(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("command template `%s` exists on the destination server", name))
		} else if err != nil {
			return nil, err
		}
		logger.Info("Deleting existing command template `%s` from `%s`", name, profile)
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	if err := services.NewAnalyticTemplateService(client).Import(in.(services.AnalyticTemplate)); err != nil {
		return nil, err
	}

	return nil, nil
}

/*
******************************************************************************
Private functions
******************************************************************************
*/

func (r *AnalyticTemplateRunner) importAnalyticTemplate(in services.AnalyticTemplate, replace bool) error {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return err
	}

	for _, ele := range res {
		if ele.Name == in.Name {
			if replace {
				if err := r.service.Delete(ele.Name); err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("analytic-template `%s` already exists", ele.Name))
			}
		}
	}

	return r.service.Import(in)
}
