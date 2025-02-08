// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type AdapterRunner struct {
	config  *config.Config
	client  client.Client
	service *services.AdapterService
}

func NewAdapterRunner(c client.Client, cfg *config.Config) *AdapterRunner {
	return &AdapterRunner{
		service: services.NewAdapterService(c),
		config:  cfg,
		client:  c,
	}
}

// Get is the implementation of the command `get adapters`
func (r *AdapterRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	adapters, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var display = []string{"NAME\tMODEL"}

	for _, ele := range adapters {
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Model))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(adapters),
	), nil
}

// Describe is the implementation of the `describe adapters` command
func (r *AdapterRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.Get(in.Args[0])
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		logger.Fatal(err, "failed to marshal data")
	}

	return NewResponse(
		string(b),
		WithJson(res),
	), nil
}

// Delete is the implementation of `delete adatper <name>`
func (r *AdapterRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	if err := r.service.Delete(in.Args[0]); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted adapter `%s`", in.Args[0]),
	), nil
}

func (r *AdapterRunner) Clear(in Request) (*Response, error) {
	return NotImplemented(in)
}

func (r *AdapterRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "adapter"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied adapter `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *AdapterRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewAdapterService(client).Get(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *AdapterRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	models, err := services.NewAdapterModelService(client).GetAll()
	if err != nil {
		return nil, err
	}

	exists := func(model string) bool {
		for _, ele := range models {
			if strings.ToLower(ele) == strings.ToLower(model) {
				return true
			}
		}
		return false
	}

	modelName := strings.Split(in.(services.Adapter).Model, "-")[1]

	if !exists(modelName) {
		return nil, errors.New("adapter model does not exist on destination server")
	}

	svc := services.NewAdapterService(client)

	name := in.(services.Adapter).Name

	if exists, err := svc.Get(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("adapter `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := svc.Create(in.(services.Adapter))
	if err != nil {
		return nil, err
	}

	return res, nil

}
func (r *AdapterRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options *flags.AdapterCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Template != "" {
		content, err := utils.ReadFromFile(options.Template)
		if err != nil {
			return nil, err
		}

		tmpl := template.Must(template.New("adapter").Parse(string(content)))

		var variables = map[string]interface{}{
			"name": name,
		}

		for _, ele := range options.Variables {
			parts := strings.Split(ele, "=")
			if len(parts) != 2 {
				return nil, errors.New("invalid variable")
			}
			variables[parts[0]] = parts[1]
		}

		buf := &bytes.Buffer{}

		if err := tmpl.Execute(buf, variables); err != nil {
			return nil, err
		}
	}

	fmt.Println(name)
	fmt.Println(options)

	return NotImplemented(in)
}

// Import provides the implementation for `import adapter <filepath>`
func (r *AdapterRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	var adapter services.Adapter
	if err := utils.ReadObjectFromDisk(path, &adapter); err != nil {
		return nil, err
	}

	if err := r.importAdapter(adapter, false); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported adapter `%s`", in.Args[0]),
	), nil

}

// Export) provides the implementation for `export adapter <name>`
func (r *AdapterRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	adapter, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.adapter.json", name)

	if err := utils.WriteJsonToDisk(adapter, fn, options.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported adapter `%s` to `%s`", adapter.Name, fn),
	), nil
}

func (r *AdapterRunner) Inspect(in Request) (*Response, error) {
	logger.Trace()

	svc := services.NewHealthService(r.client)
	adapters, err := svc.GetAdapterHealth()
	if err != nil {
		return nil, err
	}

	var display = []string{"NAME\tSTATUS\tVERSION"}

	for _, ele := range adapters {
		display = append(display, fmt.Sprintf(
			"%s\t%s\t%s", ele.Id, ele.State, ele.Version,
		))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(adapters),
	), nil
}

// Pull implements the command `pull adapter <repo>`
func (r *AdapterRunner) Pull(in Request) (*Response, error) {
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

	var adapter services.Adapter
	utils.UnmarshalData(data, &adapter)

	if err := r.importAdapter(adapter, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled adapter `%s`", adapter.Name),
	), nil
}

// Push implements the command `push adapter <repo>`
func (r *AdapterRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPushCommon
	utils.LoadObject(in.Common, &common)

	res, err := r.service.Export(in.Args[0])
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:     in.Args[1],
		Filename: fmt.Sprintf("%s.adapter.json", in.Args[0]),
		Options:  common,
		Config:   r.config,
		Data:     res,
	}

	if err := push.Do(); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed adapter `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}

func (r *AdapterRunner) Start(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Start(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully started adapter `%s`", name),
	), nil
}

func (r *AdapterRunner) Stop(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Stop(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully stopped adapter `%s`", name),
	), nil
}

func (r *AdapterRunner) Restart(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Stop(name); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully restarted adapter `%s`", name),
	), nil
}

func (r *AdapterRunner) importAdapter(in services.Adapter, replace bool) error {
	logger.Trace()

	adapters, err := r.service.GetAll()
	if err != nil {
		return err
	}

	for _, ele := range adapters {
		if ele.Name == in.Name {
			if replace {
				if err := r.service.Delete(ele.Name); err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("adapter `%s` already exists", ele.Name))
			}
		}
	}

	if _, err := r.service.Import(in); err != nil {
		return err
	}

	return nil
}
