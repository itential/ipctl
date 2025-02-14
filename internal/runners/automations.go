// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
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

type AutomationRunner struct {
	config    *config.Config
	client    client.Client
	service   *services.AutomationService
	workflows *services.WorkflowService
	triggers  *services.TriggerService
}

func NewAutomationRunner(c client.Client, cfg *config.Config) *AutomationRunner {
	return &AutomationRunner{
		config:    cfg,
		client:    c,
		service:   services.NewAutomationService(c),
		workflows: services.NewWorkflowService(c),
		triggers:  services.NewTriggerService(c),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get is the implementation of the command `get automations`
func (r *AutomationRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	automations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var display = []string{"NAME\tDESCRIPTION"}
	for _, ele := range automations {
		desc := strings.Replace(ele.Description, "\n", " ", -1)
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, desc))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(automations),
	), nil
}

func (r *AutomationRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	res, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		"",
		WithJson(res),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

func (r *AutomationRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.AutomationCreateOptions
	utils.LoadObject(in.Options, &options)

	if options.Replace {
		existing, err := r.service.GetByName(name)

		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if err.Error() != "automation not found" {
				return nil, err
			}
		}
	}

	res, err := r.service.Create(services.NewAutomation(name, options.Description))
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created automation `%s`", res.Name),
		WithJson(res),
	), nil
}

func (r *AutomationRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	automations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var selected *services.Automation

	for _, ele := range automations {
		if ele.Name == name {
			selected = ele
			break
		}
	}

	if selected != nil {
		if err := r.service.Delete(selected.Id); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted automation `%s`", name),
	), nil
}

// Clear implements the `clear automations` command
func (r *AutomationRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	automations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range automations {
		if err := r.service.Delete(ele.Id); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Deleted %v automations(s)", len(automations)),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

func (r *AutomationRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "automation"}, r)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied automation `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
}

func (r *AutomationRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewAutomationService(client).GetByName(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *AutomationRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewAutomationService(client)

	name := in.(services.Automation).Name

	if exists, err := svc.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("automation `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
	}

	res, err := svc.Import(in.(services.Automation))

	if err != nil {
		return nil, errors.New(r.formatImportErrorMessage(err))
	}

	return res, err
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the `import automation <name>` command
func (r *AutomationRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common *flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	var options *flags.AutomationImportOptions
	utils.LoadObject(in.Options, &options)

	currentUser, err := services.GetCurrentUser(r.client)
	if err != nil {
		return nil, err
	}

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	content, err := utils.ReadFromFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	utils.UnmarshalData(content, &data)

	var automation services.Automation

	if err := json.Unmarshal(content, &automation); err != nil {
		return nil, err
	}

	if !options.DisableComponentCheck && automation.ComponentType == "workflows" {
		if !r.WorkflowExists(automation.ComponentName) {
			return nil, errors.New(
				fmt.Sprintf(
					"worklow `%s` does not exist, cannot import automation", automation.ComponentName,
				),
			)
		}
	}

	if !options.DisableGroupExistsCheck {
		logger.Info("Starting group exists check")

		groupService := services.NewGroupService(r.client)

		groups, err := groupService.GetAll()
		if err != nil {
			return nil, err
		}

		var readExists bool = true
		var writeExists bool = true

		var name string

		for _, ele := range automation.Gbac.Read {
			name = ele.(map[string]interface{})["name"].(string)
			for _, g := range groups {
				readExists = g.Name == name
				if readExists {
					break
				}
			}
		}

		for _, ele := range automation.Gbac.Write {
			name = ele.(map[string]interface{})["name"].(string)
			for _, g := range groups {
				writeExists = g.Name == name
				if writeExists {
					break
				}
			}
		}

		if !readExists {
			return nil, errors.New("configured read group not found on the server")
		}

		if !writeExists {
			return nil, errors.New("configured write group not found on the server")
		}

		logger.Info("Group exists check completely successfully")
	}

	var readMember bool = true
	var writeMember bool = true

	if !options.DisableGroupReadCheck {
		for _, ele := range automation.Gbac.Read {
			for _, grp := range currentUser.Groups {
				readMember = grp.Name == ele.(map[string]interface{})["name"].(string)
				if readMember {
					break
				}
			}
		}
	}

	if !options.DisableGroupWriteCheck {
		for _, ele := range automation.Gbac.Write {
			for _, grp := range currentUser.Groups {
				writeMember = grp.Name == ele.(map[string]interface{})["name"].(string)
				if writeMember {
					break
				}
			}
		}
	}

	if !writeMember {
		if !readMember {
			return nil, errors.New("current user must have read or write access to the automation")
		}
	}

	var triggers []services.Trigger

	if value, exists := data["triggers"]; exists {
		if value != nil {
			for _, ele := range value.([]interface{}) {
				b, err := json.Marshal(ele)
				if err != nil {
					return nil, err
				}

				item := ele.(map[string]interface{})

				switch item["type"].(string) {
				case "endpoint":
					var t services.EndpointTrigger
					if err := json.Unmarshal(b, &t); err != nil {
						return nil, err
					}
					triggers = append(triggers, t)
				case "eventSystem":
					var t services.EventTrigger
					if err := json.Unmarshal(b, &t); err != nil {
						return nil, err
					}
					triggers = append(triggers, t)
				case "manual":
					var t services.ManualTrigger
					if err := json.Unmarshal(b, &t); err != nil {
						return nil, err
					}
					triggers = append(triggers, t)
				case "schedule":
					var t services.ScheduleTrigger
					if err := json.Unmarshal(b, &t); err != nil {
						return nil, err
					}
					triggers = append(triggers, t)
				}
			}
		}
	}

	automation.Triggers = triggers

	existing, err := r.service.GetByName(automation.Name)
	if existing != nil {
		if common.Replace {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("automation `%s` already exists", data["name"].(string)))
		}

	}
	if err != nil {
		if err.Error() != "automation not found" {
			return nil, err
		}
	}

	res, err := r.service.Import(automation)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported automation `%s` with %v trigger(s)", res.Name, len(triggers)),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export implements the `export automation <name>` command
func (r *AutomationRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetExportCommon)

	name := in.Args[0]

	automation, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	res, err := r.service.Export(automation.Id)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.automation.json", name)

	if err := utils.WriteJsonToDisk(res, fn, common.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported automation `%s`", name),
	), nil

}

func (r *AutomationRunner) Exists(name string) bool {
	logger.Trace()
	res, err := r.service.GetAll()
	if err != nil {
		logger.Fatal(err, "error attempting to get all automations")
	}

	for _, ele := range res {
		if ele.Name == name {
			return true
		}
	}

	return false
}

func (r *AutomationRunner) WorkflowExists(name string) bool {
	logger.Trace()
	_, err := r.workflows.Get(name)
	if err != nil {
		return false
	}
	return true
}

func (r *AutomationRunner) formatImportErrorMessage(e error) string {
	logger.Trace()

	type ResponseError struct {
		Message  string `json:"message"`
		Data     any    `json:"data"`
		Metadata struct {
			Errors []struct {
				Success bool                   `json:"success"`
				Reason  string                 `json:"reason"`
				Data    map[string]interface{} `json:"data"`
			} `json:"errors"`
		} `json:"metadata"`
	}

	var res ResponseError

	if err := json.Unmarshal([]byte(e.Error()), &res); err != nil {
		logger.Fatal(err, "failed to unmarshal error message")
	}

	var output = []string{
		fmt.Sprintf("%s (See details below)", res.Message),
	}

	for _, ele := range res.Metadata.Errors {
		output = append(output, fmt.Sprintf("- %s", ele.Reason))
	}

	return strings.Join(output, "\n")

}
