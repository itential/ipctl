// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type JsonFormRunner struct {
	config  *config.Config
	service *services.JsonFormService
}

func NewJsonFormRunner(c client.Client, cfg *config.Config) *JsonFormRunner {
	return &JsonFormRunner{
		service: services.NewJsonFormService(c),
		config:  cfg,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get json_forms` command
func (r *JsonFormRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.WorkflowGetOptions
	utils.LoadObject(in.Options, &options)

	json_forms, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME"}
	for _, ele := range json_forms {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			display = append(display, ele.Name)
		} else if !strings.HasPrefix(ele.Name, "@") {
			display = append(display, ele.Name)
		}
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(json_forms),
	), nil

}

// Describe implements the `describe json_form <name>` command
func (r *JsonFormRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	json_form, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", json_form.Name),
		WithJson(json_form),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create implements the `create jsonform <name>` command
func (r *JsonFormRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.JsonFormCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.GetByName(name)

		if existing != nil {
			if err := r.service.Delete([]string{existing.Id}); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "jsonform not found" {
				return nil, err
			}
		}
	}

	jf, err := r.service.Create(services.NewJsonForm(name, options.Description))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created jsonform `%s`", jf.Name),
		WithJson(jf),
	), nil
}

// Delete implements the `delete jsonform <name>` command
func (r *JsonFormRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var jf *services.JsonForm

	for _, ele := range elements {
		if ele.Name == name {
			jf = &ele
			break
		}
	}

	if jf == nil {
		return nil, errors.New(fmt.Sprintf("JSON form `%s` not found", name))
	}

	if err := r.service.Delete([]string{jf.Id}); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted jsonform `%s`", name),
	), nil
}

// Clear implements the `clear jsonforms` command
func (r *JsonFormRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	jsonforms, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var ids []string

	for _, ele := range jsonforms {
		ids = append(ids, ele.Id)
	}

	if err := r.service.Delete(ids); err != nil {
		return nil, err
	}

	return NewResponse(fmt.Sprintf("Deleted %v jsonform(s)", len(ids))), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

// Copy implements the `copy jsonform <name>` command
func (r *JsonFormRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "jsonform"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied jsonform `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *JsonFormRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewJsonFormService(client).GetByName(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *JsonFormRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewJsonFormService(client)

	name := in.(services.JsonForm).Name

	if exists, err := svc.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("jsonform `%s` exists on the destination server, use --replace to overwrite"))
		} else if err != nil {
			return nil, err
		}
	}

	err = svc.Import(in.(services.JsonForm))

	return nil, err
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the command `import jsonform <path>`
func (r *JsonFormRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common *flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var jsonform services.JsonForm
	utils.UnmarshalData(data, &jsonform)

	if err := r.importJsonForm(jsonform, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported jsonform `%s`", jsonform.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export is the implementation of the command `export jsonform <name>`
func (r *JsonFormRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	var options *flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	name := in.Args[0]

	jsonform, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.jsonform.json", name)

	if err := utils.WriteJsonToDisk(jsonform, fn, options.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported jsonform `%s`", jsonform.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Gitter Interface
//

// Pull implements the command `pull jsonform <repo>`
func (r *JsonFormRunner) Pull(in Request) (*Response, error) {
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

	var jsonform services.JsonForm
	utils.UnmarshalData(data, &jsonform)

	if err := r.importJsonForm(jsonform, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled jsonform `%s`", jsonform.Name),
	), nil
}

// Push implements the command `push jsonform <repo>`
func (r *JsonFormRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPushCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:     in.Args[1],
		Filename: fmt.Sprintf("%s.jsonform.json", in.Args[0]),
		Options:  common,
		Config:   r.config,
		Data:     res,
	}

	if err := push.Do(); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed jsonform `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r JsonFormRunner) importJsonForm(in services.JsonForm, replace bool) error {
	logger.Trace()

	p, err := r.service.Get(in.Name)
	if err == nil {
		if replace {
			r.service.Delete([]string{p.Id})
		} else {
			return errors.New(fmt.Sprintf("jsonform with name `%s` already exists", p.Name))
		}
	}

	err = r.service.Import(in)
	if err != nil {
		return err
	}

	return nil
}
