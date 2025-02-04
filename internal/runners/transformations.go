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

type TransformationRunner struct {
	config  *config.Config
	service *services.TransformationService
}

func NewTransformationRunner(c client.Client, cfg *config.Config) *TransformationRunner {
	return &TransformationRunner{
		service: services.NewTransformationService(c),
	}
}

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

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var trans *services.Transformation

	for _, ele := range elements {
		if ele.Name == name {
			trans = &ele
			break
		}
	}

	if trans == nil {
		return nil, errors.New(fmt.Sprintf("transformation not found"))
	}

	if err := r.service.Delete(trans.Id); err != nil {
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

func (r *TransformationRunner) Copy(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
}

func (r *TransformationRunner) importTransformation(in services.Transformation, replace bool) error {
	logger.Trace()

	p, err := r.service.Get(in.Name)
	if err == nil {
		if replace {
			if err := r.service.Delete(p.Id); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("transformation with name `%s` already exists", p.Name))
		}
	}

	_, err = r.service.Import(in)
	if err != nil {
		return err
	}

	return err
}
