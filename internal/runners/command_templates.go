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

type CommandTemplateRunner struct {
	service *services.CommandTemplateService
	config  *config.Config
}

func NewCommandTemplateRunner(c client.Client, cfg *config.Config) *CommandTemplateRunner {
	return &CommandTemplateRunner{
		service: services.NewCommandTemplateService(c),
		config:  cfg,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get command-templates` command
func (r *CommandTemplateRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.CommandTemplateGetOptions
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
func (r *CommandTemplateRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	output := []string{
		fmt.Sprintf("Name: %s (%s)", res.Name, res.Id),
	}

	return NewResponse(
		strings.Join(output, "\n"),
		WithJson(res),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

func (r *CommandTemplateRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.Create(services.NewCommandTemplate(name))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created command template `%s` (%s)", res.Name, res.Id),
		WithJson(res),
	), nil
}

func (r *CommandTemplateRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Delete(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted command-template `%s`", name),
	), nil
}

func (r *CommandTemplateRunner) Clear(in Request) (*Response, error) {
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

	return NewResponse(fmt.Sprintf("Deleted %v command template(s)", len(elements))), nil
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

func (r *CommandTemplateRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var ct services.CommandTemplate

	if err := importUnmarshalFromRequest(in, &ct); err != nil {
		return nil, err
	}

	if err := r.importCommandTemplate(ct, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported command template `%s`", ct.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

func (r *CommandTemplateRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	ct, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.command_template.json", name)

	if err := exportAssetFromRequest(in, ct, fn); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported command template `%s`", ct.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

func (r *CommandTemplateRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "command-template"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied command template `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *CommandTemplateRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewCommandTemplateService(client).Export(name)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (r *CommandTemplateRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewCommandTemplateService(client)

	name := in.(services.CommandTemplate).Name

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

	if err := services.NewCommandTemplateService(client).Import(in.(services.CommandTemplate)); err != nil {
		return nil, err
	}

	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r *CommandTemplateRunner) importCommandTemplate(in services.CommandTemplate, replace bool) error {
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
				return errors.New(fmt.Sprintf("command template `%s` already exists", ele.Name))
			}
		}
	}

	return r.service.Import(in)
}
