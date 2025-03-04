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
	"github.com/itential/ipctl/pkg/editor"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type WorkflowRunner struct {
	config  *config.Config
	service *services.WorkflowService
}

func NewWorkflowRunner(c client.Client, cfg *config.Config) *WorkflowRunner {
	return &WorkflowRunner{
		config:  cfg,
		service: services.NewWorkflowService(c),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get workflows` command
func (r *WorkflowRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.WorkflowGetOptions
	utils.LoadObject(in.Options, &options)

	workflows, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range workflows {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
		} else if !strings.HasPrefix(ele.Name, "@") {
			display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
		}
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(workflows),
	), nil

}

// Describe implements the `describe workflow ...` command
func (r *WorkflowRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	workflow, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	human := strings.Join([]string{
		"Name: %s (ID: %s)",
		"Url: %s",
		"Description:\n%s",
	}, "\n")

	return NewResponse(
		fmt.Sprintf(human, workflow.Name, workflow.Id, r.makeUrl(name), workflow.Description),
		WithJson(workflow),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer interface
//

// Create is the implementation of the `create workflow ...` command
func (r *WorkflowRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	wf, err := r.service.Create(services.NewWorkflow(name))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created workflow `%s` (%s)", name, r.makeUrl(name)),
		WithJson(wf),
	), nil
}

// Delete is the implementation of the `delete workflow ...` command
func (r *WorkflowRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	if err := r.service.Delete(in.Args[0]); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted workflow `%s`", in.Args[0]),
	), nil
}

// Clear is the implementation of the `clear workflows` command
func (r *WorkflowRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	cnt := 0

	workflows, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range workflows {
		r.service.Delete(ele.Id)
		cnt++
	}

	return NewResponse(
		fmt.Sprintf("Deleted %v workflow(s)", cnt),
	), nil
}

func (r *WorkflowRunner) Edit(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	current, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	var updated services.Workflow

	if err := editor.Run(current, &updated); err != nil {
		return nil, err
	}

	if _, err := r.service.Update(updated); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully updated workflow `%s`", name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier interface
//

// Copy implements the `copy workflow ...` command
func (r *WorkflowRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "workflow"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied workflow `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *WorkflowRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewWorkflowService(client).Export(name)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (r *WorkflowRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewWorkflowService(client)

	name := in.(services.Workflow).Name

	if exists, err := svc.Get(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("workflow `%s` exists on the destination server", name))
		} else if err != nil {
			return nil, err
		}
		logger.Info("Deleting existing workflow `%s` from `%s`", name, profile)
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := services.NewWorkflowService(client).Import(in.(services.Workflow))
	if err != nil {
		return nil, err
	}

	return res, nil
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the command `import workflow <path>`
func (r *WorkflowRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var workflow services.Workflow

	if err := importUnmarshalFromRequest(in, &workflow); err != nil {
		return nil, err
	}

	if err := r.importWorkflow(workflow, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported workflow `%s`", workflow.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export is the implementation of the `export workflow ...` command
func (r *WorkflowRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	workflow, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.workflow.json", name)

	if err := exportAssetFromRequest(in, workflow, fn); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported workflow `%s` (%s)", workflow.Name, workflow.Id),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r *WorkflowRunner) importWorkflow(in services.Workflow, replace bool) error {
	logger.Trace()

	res, err := r.service.Get(in.Name)
	if err == nil {
		if res != nil {
			if replace {
				if err := r.service.Delete(res.Name); err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("workflow with name `%s` already exists, use `--replace` to overwrite", res.Name))
			}
		} else {
			return err
		}
	}

	_, err = r.service.Import(in)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkflowRunner) makeUrl(name string) string {
	logger.Trace()
	return makeUrl(r.config, "/automation-studio/#/edit?tab=0&workflow=%s", name)
}
