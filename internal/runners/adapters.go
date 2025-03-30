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
	"github.com/itential/ipctl/pkg/editor"
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

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get is the implementation of the command `get adapters`
func (r *AdapterRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	adapters, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "model"},
		Object: adapters,
	}, nil
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

	return &Response{
		Text:   string(b),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

func (r *AdapterRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	options := in.Options.(*flags.AdapterCreateOptions)

	adapter := services.Adapter{
		Name: name,
		Type: "Adapter",
		Properties: services.AdapterProperties{
			Id: name,
			//Type: options.Type,
		},
	}

	if options.Model != "" {
		adapter.Model = options.Model
	}

	if options.Template != "" {
		content, err := utils.ReadFromFile(options.Template)
		if err != nil {
			return nil, err
		}

		tmpl := template.Must(
			template.New("adapter").
				Option("missingkey=error").
				Parse(string(content)))

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

		var props map[string]interface{}

		if err := json.Unmarshal([]byte(buf.String()), &props); err != nil {
			return nil, err
		}

		adapter.Properties.Properties = props
	}

	res, err := r.service.Create(adapter)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully created adapter `%s`", res.Name),
	}, nil
}

// Delete is the implementation of `delete adatper <name>`
func (r *AdapterRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	if err := r.service.Delete(in.Args[0]); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted adapter `%s`", in.Args[0]),
	}, nil
}

func (r *AdapterRunner) Clear(in Request) (*Response, error) {
	return notImplemented(in)
}

//////////////////////////////////////////////////////////////////////////////
// Editor Interface
//

// Edit implements the `edti adapter ...` command
func (r *AdapterRunner) Edit(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	current, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	var updated services.Adapter

	if err := editor.Run(current, &updated); err != nil {
		return nil, err
	}

	if _, err := r.service.Update(updated); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully updated adapter `%s`", name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

func (r *AdapterRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "adapter"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied adapter `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
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

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import provides the implementation for `import adapter <filepath>`
func (r *AdapterRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var adapter services.Adapter

	if err := importUnmarshalFromRequest(in, &adapter); err != nil {
		return nil, err
	}

	if err := r.importAdapter(adapter, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported adapter `%s`", adapter.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export) provides the implementation for `export adapter <name>`
func (r *AdapterRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	adapter, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.adapter.json", name)

	if err := exportAssetFromRequest(in, adapter, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported adapter `%s`", adapter.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Inspector interface
//

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

	return &Response{
		Keys:   []string{"name", "status", "version"},
		Object: adapters,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Controller interface
//

func (r *AdapterRunner) Start(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Start(name); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully started adapter `%s`", name),
	}, nil
}

func (r *AdapterRunner) Stop(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Stop(name); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully stopped adapter `%s`", name),
	}, nil
}

func (r *AdapterRunner) Restart(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Stop(name); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully restarted adapter `%s`", name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Dumper Interface
//

// Dump implements the `dump adapters...` command
func (r *AdapterRunner) Dump(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var assets = map[string]interface{}{}
	for _, ele := range res {
		key := fmt.Sprintf("%s.adapter.json", ele.Name)
		assets[key] = ele
	}

	if err := dumpAssets(in, assets); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Dumped %v adapter(s)", len(assets)),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Loader Interface
//

// Load implements the `load adapters ...` command
func (r *AdapterRunner) Load(in Request) (*Response, error) {
	logger.Trace()

	elements, err := loadAssets(in)
	if err != nil {
		return nil, err
	}

	var loaded int
	var skipped int

	var output []string

	for fn, ele := range elements {
		var adapter services.Adapter

		if err := loadUnmarshalAsset(ele, &adapter); err != nil {
			output = append(output, fmt.Sprintf("Failed to load adapter from `%s`, skipping", fn))
			skipped++
		} else {
			if _, err := r.service.Import(adapter); err != nil {
				if !strings.HasSuffix(err.Error(), "already exists!\"") {
					return nil, err
				}
				output = append(output, fmt.Sprintf("Skipping `%s`, adapter `%s` already exists", fn, adapter.Name))
				skipped++
			} else {
				output = append(output, fmt.Sprintf("Loaded adapter `%s` successfully from `%s`", adapter.Name, fn))
				loaded++
			}
		}
	}

	output = append(output, fmt.Sprintf(
		"\nSuccessfully loaded %v and skipped %v files from `%s`", loaded, skipped, in.Args[0],
	))

	return &Response{
		Text: strings.Join(output, "\n"),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

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
