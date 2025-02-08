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

type TemplateRunner struct {
	service *services.TemplateService
	config  *config.Config
}

func NewTemplateRunner(c client.Client, cfg *config.Config) *TemplateRunner {
	return &TemplateRunner{
		service: services.NewTemplateService(c),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get command-templates` command
func (r *TemplateRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.TemplateGetOptions
	utils.LoadObject(in.Options, &options)

	templates, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME"}
	for _, ele := range templates {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			display = append(display, ele.Name)
		} else if !strings.HasPrefix(ele.Name, "@") {
			display = append(display, ele.Name)
		}
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(templates),
	), nil

}

// Describe implements the `describe command-template <name>` command
func (r *TemplateRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	template, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", template.Name),
		WithJson(template),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create implements the `create template ...` command
func (r *TemplateRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.TemplateCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.GetByName(name)

		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "template not found" {
				return nil, err
			}
		}
	}

	res, err := r.service.Create(services.NewTemplate(
		name,
		options.Group,
		options.Description,
		options.Type,
	))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created template `%s`", res.Name),
		WithJson(res),
	), nil
}

func (r *TemplateRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	t, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(t.Id); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted template `%s`", t.Name),
	), nil
}

func (r *TemplateRunner) Clear(in Request) (*Response, error) {
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

	return NewResponse(fmt.Sprintf("Deleted %v template(s)", len(elements))), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

func (r *TemplateRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "template"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied template `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *TemplateRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewTemplateService(client).Export(name)
	if err != nil {
		return nil, err
	}
	return *res, nil

}

func (r *TemplateRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewTemplateService(client)

	name := in.(services.Template).Name

	if exists, err := svc.Get(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("template `%s` exists on the destination server", name))
		} else if err != nil {
			return nil, err
		}
		logger.Info("Deleting existing template `%s` from `%s`", name, profile)
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := services.NewTemplateService(client).Import(in.(services.Template))
	if err != nil {
		return nil, err
	}

	return res, nil

}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

func (r *TemplateRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var template services.Template
	if err := ReadImportFromFile(in, &template); err != nil {
		return nil, err
	}

	res, err := r.service.Import(template)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported command template `%s`", res.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

func (r *TemplateRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	exported, err := r.service.Export(res.Id)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.template.json", exported.Name)

	if err := utils.WriteJsonToDisk(exported, fn, common.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported template `%s`", exported.Name),
	), nil

}
