// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type ModelRunner struct {
	config  *config.Config
	service *services.ModelService
}

func NewModelRunner(client client.Client, cfg *config.Config) *ModelRunner {
	return &ModelRunner{
		config:  cfg,
		service: services.NewModelService(client),
	}
}

// GetModels() is the implementation of the command `get models`
func (r *ModelRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	models, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range models {
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(models),
	), nil

}

func (r *ModelRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	model, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", model.Name),
		WithJson(model),
	), nil
}

// Create implements the `create model <name>` command
func (r *ModelRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.ModelCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.GetByName(name)

		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "model not found" {
				return nil, err
			}
		}
	}

	model := services.NewModel(name, options.Description)

	if options.Schema != "" {
		data, err := os.ReadFile(options.Schema)
		if err != nil {
			return nil, err
		}

		var schema map[string]interface{}

		if err := json.Unmarshal(data, &schema); err != nil {
			return nil, err
		}

		model.Schema = schema
	}

	jf, err := r.service.Create(model)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created model `%s`", jf.Name),
		WithJson(jf),
	), nil
}

// Delete implements the `delete model <name>` command
func (r *ModelRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	elements, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var model *services.Model

	for _, ele := range elements {
		if ele.Name == name {
			model = &ele
			break
		}
	}

	if model == nil {
		return nil, errors.New(fmt.Sprintf("model `%s` not found", name))
	}

	if err := r.service.Delete(model.Id); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted model `%s`", name),
	), nil
}

// Clear implements the `clear models` command
func (r *ModelRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	models, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range models {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Deleted %v model(s)", len(models)),
	), nil
}

// Copy implements the `copy model <name>` command
func (r *ModelRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common *flags.AssetCopyCommon
	utils.LoadObject(in.Common, &common)

	if common.From == common.To {
		return nil, errors.New("source and destination servers must be different values")
	}

	src, err := r.CopyFrom(common.From, name)
	if err != nil {
		return nil, err
	}

	_, err = r.CopyTo(common.To, *src)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied model `%s` from `%s` to `%s`", name, common.From, common.To),
	), nil
}

// Import implements the command `import model <path>`
func (r *ModelRunner) Import(in Request) (*Response, error) {
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

	var model services.Model
	utils.UnmarshalData(data, &model)

	if err := r.importModel(model, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported model `%s`", model.Name),
		WithJson(model),
	), nil
}

// Export implements the `export model ...` command
func (r *ModelRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	var options *flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	model, err := r.service.Export(res.Id)
	if err != nil {
		return nil, err
	}

	out, err := Export(
		model,
		WithExportName(name),
		WithExportType("model"),
		WithExportEncoding("json"),
		WithExportPath(options.Path),
	)

	return NewResponse(
		fmt.Sprintf("Successfully exported model `%s` to `%s`", model.Name, out.AbsPath),
		WithJson(out),
	), nil
}

func (r *ModelRunner) importModel(in services.Model, replace bool) error {
	logger.Trace()

	res, err := r.service.Get(in.Name)
	if err == nil {
		if replace {
			if err := r.service.Delete(res.Id); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("model with name `%s` already exists", res.Name))
		}
	}

	_, err = r.service.Import(in)
	if err != nil {
		return err
	}

	return nil
}

func (r *ModelRunner) CopyFrom(profile, name string) (*services.Model, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewModelService(client)

	res, err := svc.GetByName(name)
	if err != nil {
		return nil, err
	}

	return svc.Export(res.Id)
}

func (r *ModelRunner) CopyTo(profile string, in services.Model) (*services.Model, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewModelService(client)

	if _, err := svc.GetByName(in.Name); err != nil {
		if err.Error() != "model not found" {
			return nil, errors.New(fmt.Sprintf("model `%s` already exists on the destination server", in.Name))
		}
	}

	return svc.Import(in)
}
