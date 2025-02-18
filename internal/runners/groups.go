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
		config:  cfg,
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
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "group"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied group `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *GroupRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewGroupService(client).GetByName(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *GroupRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewGroupService(client)

	name := in.(services.Group).Name

	if exists, err := svc.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("group `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := svc.Create(in.(services.Group))
	if err != nil {
		return nil, err
	}

	return res, nil

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
}

func (r *GroupRunner) importGroup(in services.Group, replace bool) error {
	logger.Trace()

	existing, err := r.service.GetByName(in.Name)

	if err != nil {
		if err.Error() != "group does not exist" {
			return err
		}
	}

	if existing != nil {
		if replace {
			if err := r.service.Delete(existing.Id); err != nil {
				return err
			}
		} else {
			return errors.New(
				fmt.Sprintf("group `%s` already exists, use --replace to overwrite it", in.Name),
			)
		}
	}

	_, err = r.service.Create(in)

	return err

}

func (r *GroupRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	var group services.Group

	if err := importFile(in, &group); err != nil {
		return nil, err
	}

	if err := r.importGroup(group, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported group `%s`", group.Name),
	), nil
}

// Pull implements the command `pull profile <repo>`
func (r *GroupRunner) Pull(in Request) (*Response, error) {
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

	var group services.Group
	utils.UnmarshalData(data, &group)

	if err := r.importGroup(group, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled group `%s`", group.Name),
	), nil
}

// Push implements the command `push profile <repo>`
func (r *GroupRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPushCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:     in.Args[1],
		Filename: fmt.Sprintf("%s.group.json", in.Args[0]),
		Options:  common,
		Config:   r.config,
		Data:     res,
	}

	if err := push.Do(); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed group `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}
