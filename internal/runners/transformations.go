// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type TransformationRunner struct {
	config  *config.Config
	service *services.TransformationService
}

func NewTransformationRunner(c client.Client, cfg *config.Config) *TransformationRunner {
	return &TransformationRunner{
		service: services.NewTransformationService(c),
		config:  cfg,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader interface
//

// Get implements the `get transformations` command
func (r *TransformationRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.WorkflowGetOptions
	utils.LoadObject(in.Options, &options)

	transformations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range transformations {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
		} else if !strings.HasPrefix(ele.Name, "@") {
			display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
		}
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(transformations),
	), nil

}

// Describe implements the `describe transformation <name>` command
func (r *TransformationRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	transformation, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", transformation.Name),
		WithJson(transformation),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer interface
//

// Create implements the `create transformation <name>` command
func (r *TransformationRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.TransformationCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.GetByName(name)

		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "transformation not found" {
				return nil, err
			}
		}
	}

	res, err := r.service.Create(services.NewTransformation(name, options.Description))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created transformation `%s`", res.Name),
		WithJson(res),
	), nil
}

// Delete is the implementation of the command `delete transformation <name>`
func (r *TransformationRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New(fmt.Sprintf("transformation not found"))
	}

	if err := r.service.Delete(res.Id); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted transformation `%s`", name),
	), nil
}

// Clear is the implementation of the command `clear transformations`
func (r *TransformationRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	transformations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range transformations {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Deleted %v transformation(s)", len(transformations)),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier interface
//

func (r *TransformationRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "transformation"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied transformation `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil

}

func (r *TransformationRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewTransformationService(client).GetByName(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *TransformationRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewTransformationService(client)

	name := in.(services.Transformation).Name

	if exists, err := svc.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("transformation `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := svc.Import(in.(services.Transformation))
	if err != nil {
		return nil, err
	}

	return res, nil
}

//////////////////////////////////////////////////////////////////////////////
// Importer interface
//

func (r *TransformationRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common *flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	var res services.Transformation
	if err := importFile(in, &res); err != nil {
		return nil, err
	}

	if err := r.importTransformation(res, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported transformation `%s`", res.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter interface
//

func (r *TransformationRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]
	common := in.Common.(*flags.AssetExportCommon)

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.transformation.json", name)

	if err := NewExportAction(res, fn, common).Do(); err != nil {
		return nil, err
	}

	fp := filepath.Join(common.Path, fn)

	return NewResponse(
		fmt.Sprintf(exportSuccessMessage, "transformation", res.Name, fp),
	), nil

}

//////////////////////////////////////////////////////////////////////////////
// Gitter interface
//

// Pull implements the command `pull transformation <repo>`
func (r *TransformationRunner) Pull(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPullCommon
	utils.LoadObject(in.Common, &common)

	pull := PullAction{
		Name:     in.Args[1],
		Filename: in.Args[0],
		Config:   r.config,
		Options:  common,
	}

	data, err := pull.Do()
	if err != nil {
		return nil, err
	}

	var res services.Transformation
	utils.UnmarshalData(data, &res)

	if err := r.importTransformation(res, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled transformation `%s`", res.Name),
	), nil
}

// Push implements the command `push transformation <repo>`
func (r *TransformationRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPushCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:     in.Args[1],
		Filename: fmt.Sprintf("%s.transformation.json", in.Args[0]),
		Options:  common,
		Config:   r.config,
		Data:     res,
	}

	if err := push.Do(); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed transformation `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r *TransformationRunner) importTransformation(in services.Transformation, replace bool) error {
	logger.Trace()

	p, err := r.service.GetByName(in.Name)
	if err == nil {
		if replace {
			if err := r.service.Delete(p.Id); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("transformation `%s` already exists, use `--replace` to overwrite", p.Name))
		}
	}

	_, err = r.service.Import(in)
	if err != nil {
		return err
	}

	return err
}
