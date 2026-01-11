// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/services"
)

type TagRunner struct {
	BaseRunner
	service *services.TagService
}

func NewTagRunner(client client.Client, cfg *config.Config) *TagRunner {
	return &TagRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		service:    services.NewTagService(client),
	}
}

// Get implements the `get tags` command
func (r *TagRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	tags, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: tags,
	}, nil
}

// Describe implements the `describe tag <name>` command
func (r *TagRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	tag, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: tag,
	}, nil
}

func (r *TagRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	var options flags.TagCreateOptions
	utils.LoadObject(in.Options, &options)

	if _, err := r.service.Create(
		services.NewTag(name, options.Description),
	); err != nil {
		if strings.HasPrefix(err.Error(), "\"E11000 duplicate key error") {
			return nil, errors.New(fmt.Sprintf("tag `%s` already exists", name))
		} else {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfully created new tag `%s`", name),
	}, nil
}

func (r *TagRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	tag, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(tag.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted tag `%s`", name),
	}, nil
}

func (r *TagRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	tags, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range tags {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted %v tag(s)", len(tags)),
	}, nil
}

func (r *TagRunner) Copy(in Request) (*Response, error) {
	logging.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "tag"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied tag `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil

}

func (r *TagRunner) CopyFrom(profile, name string) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewTagService(client).GetByName(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *TagRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewTagService(client)

	name := in.(services.Tag).Name

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

	res, err := svc.Create(in.(services.Tag))
	if err != nil {
		return nil, err
	}

	return res, nil

}
