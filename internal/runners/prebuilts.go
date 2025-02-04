// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type PrebuiltRunner struct {
	config  *config.Config
	service *services.PrebuiltService
	client  client.Client
}

func NewPrebuiltRunner(client client.Client, cfg *config.Config) *PrebuiltRunner {
	return &PrebuiltRunner{
		config:  cfg,
		client:  client,
		service: services.NewPrebuiltService(client),
	}
}

// Get() implements the "get prebuilts" command
func (r *PrebuiltRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	prebuilts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range prebuilts {
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(prebuilts),
	), nil

}

// Describe() implements the `describe prebuilt <name>` command
func (r *PrebuiltRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	prebuilt, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", prebuilt.Name),
		WithJson(prebuilt),
	), nil
}

// Create is the implementation of the command `ccreate prebuilt <name>`
func (r *PrebuiltRunner) Create(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
}

// Delete is the implementation of the command `delete prebuilt <name>`
func (r *PrebuiltRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options *flags.PrebuiltDeleteOptions
	utils.LoadObject(in.Options, &options)

	prebuilt, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	if options.All {
		for _, ele := range prebuilt.Components {
			switch ele.Type {
			case "workflow":
				logger.Info("Checking for workflow: %s\n", ele.Name)
				svc := services.NewWorkflowService(r.client)
				exists, err := svc.Get(ele.Name)
				if err != nil {
					if !strings.HasPrefix(err.Error(), "could not find workflow") {
						return nil, err
					}
					logger.Info("Workflow %s not found, skipping", ele.Name)
				}
				if exists != nil {
					if err := svc.Delete(ele.Name); err != nil {
						return nil, err
					}
				}
			case "transformation":
				logger.Info("Checking transformation: %s", ele.Name)
				svc := services.NewTransformationService(r.client)
				exists, err := svc.Get(ele.Id)
				if err != nil {
					var message map[string]interface{}
					if err := json.Unmarshal([]byte(err.Error()), &message); err != nil {
						return nil, err
					}
					httpMessage := message["error"].(map[string]interface{})
					if httpMessage["code"].(float64) != 404 {
						return nil, err
					}
					logger.Info("Transformation `%s` does not exist, skipping", ele.Name)
				}
				if exists != nil {
					if err := svc.Delete(ele.Id); err != nil {
						return nil, err
					}
				}
			case "json-forms":
				logger.Info("JSON Form: %s\n", ele.Name)
				svc := services.NewJsonFormService(r.client)
				exists, err := svc.Get(ele.Id)
				if err != nil {
					if err.Error() != "\"Form not found\"" {
						return nil, err
					}
				}
				if exists != nil {
					if err := svc.Delete([]string{ele.Id}); err != nil {
						return nil, err
					}
				}
			case "automation":
				logger.Info("Automation: %s\n", ele.Name)
				svc := services.NewAutomationService(r.client)
				exists, err := svc.Get(ele.Id)
				if err != nil {
					var message map[string]interface{}
					if err := json.Unmarshal([]byte(err.Error()), &message); err != nil {
						return nil, err
					}
					if !strings.HasPrefix(message["message"].(string), "Cannot find Automation") {
						return nil, err
					}
					logger.Info("Automation `%s` not found, skipping delete", ele.Name)
				}
				if exists != nil {
					if err := svc.Delete(ele.Id); err != nil {
						return nil, err
					}
				}

			case "mop-template":
				logger.Info("MOP Template: %s\n", ele.Name)
				svc := services.NewCommandTemplateService(r.client)
				exists, err := svc.Get(ele.Id)
				if err != nil {
					if err.Error() != "command template not found" {
						return nil, err
					}
				}
				if exists != nil {
					if err := svc.Delete(ele.Id); err != nil {
						return nil, err
					}
				}
			case "template":
				logger.Info("Template: %s", ele.Name)
				svc := services.NewTemplateService(r.client)
				exists, err := svc.GetByName(ele.Name)
				if err != nil {
					if err.Error() != "template not found" {
						return nil, err
					}
				}
				if exists != nil {
					if err := svc.Delete(exists.Id); err != nil {
						return nil, err
					}
				}

			}
		}
	}

	r.service.Delete(prebuilt.Id)

	return &Response{
		Text: fmt.Sprintf("Successfully deleted prebuilt `%s`", name),
	}, nil
}

// Clear is the implementation of the command `clear prebuilts`
func (r *PrebuiltRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	cnt := 0

	prebuilts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range prebuilts {
		r.service.Delete(ele.Id)
		cnt++
	}

	return NewResponse(fmt.Sprintf("Deleted %v prebuilt(s)", cnt)), nil
}

// Copy implements the `copy prebuilt <name> <dst>` command
func (r *PrebuiltRunner) Copy(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
}

// Import implements the command `import prebuilt <path>`
func (r *PrebuiltRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	data, err := utils.ReadFromFile(path)
	if err != nil {
		return nil, err
	}

	var pkg services.PrebuiltPackage
	utils.UnmarshalData(data, &pkg)

	if !common.Force {
		if err := r.validatePackage(pkg); err != nil {
			return nil, err
		}
	}

	var prebuilt map[string]interface{}
	utils.UnmarshalData(data, &prebuilt)

	imported, err := r.service.ImportRaw(prebuilt, common.Force)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported prebuilt `%s`", imported.Name),
	), nil

}

// Export is the implementation of the command `export prebuilt <name>`
func (r *PrebuiltRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	var options flags.PrebuiltExportOptions
	utils.LoadObject(in.Options, &options)

	pb, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if common.Path != "" {
		path = common.Path
	}

	if options.Expand {
		for _, ele := range pb.Components {
			outputPath := filepath.Join(path, fmt.Sprintf("bundles/%ss", ele.Type))

			if err := utils.EnsurePathExists(outputPath); err != nil {
				return nil, err
			}

			var res any
			var err error
			var fn string

			switch ele.Type {
			case "template":
				res, err = services.NewTemplateService(r.client).Get(ele.Id)
				if res != nil {
					fn = res.(*services.Template).Name
				}
			case "mop-template":
				res, err = services.NewCommandTemplateService(r.client).Get(ele.Id)
				if err != nil {
					return nil, err
				}
				if res != nil {
					fn = res.(*services.CommandTemplate).Name
				}
			case "workflow":
				res, err = services.NewWorkflowService(r.client).Get(ele.Name)
				if res != nil {
					fn = res.(*services.Workflow).Name
				}
			case "json-forms":
				res, err = services.NewJsonFormService(r.client).Get(ele.Id)
				if res != nil {
					fn = res.(*services.JsonForm).Name
				}
			case "transformation":
				res, err = services.NewTransformationService(r.client).Get(ele.Id)
				if res != nil {
					fn = res.(*services.Transformation).Name
				}
			case "automation":
				res, err = services.NewAutomationService(r.client).Get(ele.Id)
				if res != nil {
					fn = res.(*services.Automation).Name
				}
			default:
				return nil, errors.New(fmt.Sprintf("unknown prebuilt component: %s", ele.Type))
			}

			if res == nil {
				return nil, errors.New(fmt.Sprintf(
					"unable to find prebuit component %s, type %s", ele.Name, ele.Type,
				))
			}

			if err != nil {
				return nil, err
			}

			if err := utils.WriteJsonToDisk(res, fn, outputPath); err != nil {
				return nil, err
			}
		}

		fn := fmt.Sprintf("%s.prebuilt.json", pb.Name)
		if err := utils.WriteJsonToDisk(pb, fn, path); err != nil {
			return nil, err
		}

	} else {
		res, err := r.service.Export(pb.Id)
		if err != nil {
			return nil, err
		}

		fn := fmt.Sprintf("%s.prebuilt.json", strings.Replace(res.Metadata.Name, "/", "_", 1))

		if err := utils.WriteJsonToDisk(res, fn, common.Path); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported prebuilt `%s`", name),
	), nil
}

func (r *PrebuiltRunner) validatePackage(in services.PrebuiltPackage) error {
	logger.Trace()

	for _, ele := range in.Bundles {
		logger.Debug("validating prebuilt asset of type %s", ele.Type)
		switch ele.Type {
		case "automation":
			if err := r.validateAutomation(ele.Data); err != nil {
				return err
			}
		case "workflow":
			if err := r.validateWorkflow(ele.Data); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateWorkflow will check if the workflow already exists on the
// destination IAP and return an error if it does.
func (r *PrebuiltRunner) validateWorkflow(in map[string]interface{}) error {
	logger.Trace()

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	var wf services.Workflow
	if err := json.Unmarshal(b, &wf); err != nil {
		return err
	}

	svc := services.NewWorkflowService(r.client)

	exists, _ := svc.Get(wf.Name)
	if exists != nil {
		return errors.New(fmt.Sprintf("workflow `%s` already exists", wf.Name))
	}

	return nil
}

func (r *PrebuiltRunner) validateAutomation(in map[string]interface{}) error {
	logger.Trace()

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	var automation services.Automation
	if err := json.Unmarshal(b, &automation); err != nil {
		return err
	}

	svc := services.NewAutomationService(r.client)

	exists, _ := svc.Get(automation.Id)
	if exists != nil {
		return errors.New(fmt.Sprintf("automation `%s` already exists", automation.Name))
	}

	/*
		for _, ele := range automation.Gbac.Write {
			item := ele.(map[string]interface{})
			if item["provenance"].(string) == "LDAP" {
				return errors.New("cannot import automation with LDAP group")
			}
		}

		for _, ele := range automation.Gbac.Read {
			item := ele.(map[string]interface{})
			if item["provenance"].(string) == "LDAP" {
				return errors.New("cannot import automation with LDAP group")
			}
		}
	*/

	return nil
}

func (r *PrebuiltRunner) GetAutomationIdFromName(name string) (string, error) {
	logger.Trace()

	svc := services.NewAutomationService(r.client)

	automations, err := svc.GetAll()
	if err != nil {
		return "", err
	}

	var automationId string

	for _, ele := range automations {
		if ele.Name == name {
			automationId = ele.Id
			break
		}
	}

	if automationId == "" {
		return "", errors.New(fmt.Sprintf("automation `%s` not found", name))
	}

	return automationId, nil
}

func (r *PrebuiltRunner) GetByName(name string) (*services.Prebuilt, error) {
	logger.Trace()

	prebuilts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var prebuiltId string
	for _, ele := range prebuilts {
		if ele.Name == name {
			prebuiltId = ele.Id
		}
	}

	if prebuiltId == "" {
		return nil, errors.New(fmt.Sprintf("prebuilt `%s` does not exist", name))
	}

	prebuilt, err := r.service.Get(prebuiltId)
	if err != nil {
		return nil, err
	}

	return prebuilt, nil
}
