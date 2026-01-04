// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type TransformationRunner struct {
	BaseRunner
	service *services.TransformationService
}

func NewTransformationRunner(c client.Client, cfg *config.Config) *TransformationRunner {
	return &TransformationRunner{
		service:    services.NewTransformationService(c),
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader interface
//

// Get implements the `get transformations` command
func (r *TransformationRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	options := in.Options.(*flags.TransformationGetOptions)

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var transformations []services.Transformation

	for _, ele := range res {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			transformations = append(transformations, ele)
		} else if !strings.HasPrefix(ele.Name, "@") {
			transformations = append(transformations, ele)
		}
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: transformations,
	}, nil

}

// Describe implements the `describe transformation ...` command
func (r *TransformationRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	output := []string{
		fmt.Sprintf("Name: %s (%s)", res.Name, res.Id),
		fmt.Sprintf("Description: %s", res.Description),
		fmt.Sprintf("Created: %s", res.Created),
		fmt.Sprintf("Updated: %s", res.LastUpdated),
	}

	return &Response{
		Text:   strings.Join(output, "\n"),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer interface
//

// Create implements the `create transformation <name>` command
func (r *TransformationRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	options := in.Options.(*flags.TransformationCreateOptions)

	name := in.Args[0]

	res, err := r.service.Create(
		services.NewTransformation(name, options.Description),
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created transformation `%s` (%s)", res.Name, res.Id),
		Object: res,
	}, nil
}

// Delete is the implementation of the command `delete transformation <name>`
func (r *TransformationRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(res.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted transformation `%s` (%s)", res.Name, res.Id),
	}, nil
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
			logger.Debug("failed to delete transformation `%s` (%s)", ele.Name, ele.Id)
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v transformation(s)", len(transformations)),
	}, nil
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

	return &Response{
		Text: fmt.Sprintf("Successfully copied transformation `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil

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

	common := in.Common.(*flags.AssetImportCommon)

	var res services.Transformation

	if err := importUnmarshalFromRequest(in, &res); err != nil {
		return nil, err
	}

	if err := r.importTransformation(res, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported transformation `%s` (%s)", res.Name, res.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter interface
//

func (r *TransformationRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.transformation.json", name)

	if err := exportAssetFromRequest(in, res, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported transformation `%s` (%s)", res.Name, res.Id),
	}, nil
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
