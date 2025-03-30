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

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

func (r *GroupRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: groups,
	}, nil

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

	return &Response{
		Object: grp,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

func (r *GroupRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	var options flags.GroupCreateOptions
	utils.LoadObject(in.Options, &options)

	group := services.NewGroup(in.Args[0], options.Description)

	res, err := r.service.Create(group)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created group `%s`", in.Args[0]),
		Object: res,
	}, nil
}

func (r *GroupRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	group, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	if group.Provenance != "Pronghorn" {
		return nil, errors.New("cannot delete non-local group")
	}

	if err := r.service.Delete(group.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted group `%s`", group.Name),
	}, nil
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

	return &Response{
		Text: fmt.Sprintf("Successfully deleted %v group(s)", cnt),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

func (r *GroupRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "group"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied group `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
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

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

func (r *GroupRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var grp services.Group

	if err := importUnmarshalFromRequest(in, &grp); err != nil {
		return nil, err
	}

	if err := r.importGroup(grp, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported group `%s`", grp.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

func (r *GroupRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	grp, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.group.json", name)

	if err := exportAssetFromRequest(in, grp, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported gropu `%s` (%s)", grp.Name, grp.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

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
