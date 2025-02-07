// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type GroupRunner struct {
	service *services.GroupService
	config  *config.Config
}

func NewGroupRunner(c client.Client, cfg *config.Config) *GroupRunner {
	return &GroupRunner{
		service: services.NewGroupService(c),
	}
}

func (r *GroupRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range groups {
		line := fmt.Sprintf("%s\t%s", ele.Name, ele.Description)
		display = append(display, line)
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(groups),
	), nil

}

func (r *GroupRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var grp *services.Group

	for _, ele := range groups {
		if ele.Name == name {
			grp = &ele
			break
		}
	}

	if grp == nil {
		return nil, errors.New(
			fmt.Sprintf("Group with name `%s` does not exist", name),
		)
	}

	return NewResponse(
		"",
		WithJson(grp),
	), nil
}

func (r *GroupRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	var options flags.GroupCreateOptions
	utils.LoadObject(in.Options, &options)

	group := services.NewGroup(in.Args[0], options.Description)

	res, err := r.service.Create(group)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created group `%s`", in.Args[0]),
		WithJson(res),
	), nil
}

func (r *GroupRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	group, err := r.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	if group.Provenance != "Pronghorn" {
		return nil, errors.New("cannot delete non-local group")
	}

	if err := r.service.Delete(group.Id); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted group `%s`", group.Name),
	), nil
}

func (r *GroupRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var cnt int = 0

	for _, ele := range groups {
		if ele.Provenance == "Pronghorn" {
			if err := r.service.Delete(ele.Id); err != nil {
				return nil, err
			}
			cnt++
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted %v group(s)", cnt),
	), nil
}

func (r *GroupRunner) Copy(in Request) (*Response, error) {
	return NotImplemented(in)
}

// GetByName will retrive the group from the server by name.  If the gorup does
// not exist, an error is returned
func (r *GroupRunner) GetByName(name string) (*services.Group, error) {
	logger.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range groups {
		if ele.Name == name {
			return &ele, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("group with name `%s` does not exist", name))
}

func (r *GroupRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.group.json", name)

	if err := utils.WriteJsonToDisk(res, fn, common.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported group `%s` to `%s`", res.Name, fn),
	), nil

	return NewResponse(
		"",
		WithJson(res),
	), nil
}

func (r *GroupRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	var group services.Group

	if err := importFile(in, &group); err != nil {
		return nil, err
	}

	existing, err := r.service.GetByName(group.Name)

	if err != nil {
		if err.Error() != "group does not exist" {
			return nil, err
		}
	}

	if existing != nil {
		if common.Replace {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(
				fmt.Sprintf("group `%` already exists, use --replace to overwrite it"),
			)
		}
	}

	_, err = r.service.Create(group)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported group `%s`", group.Name),
	), nil
}
