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
	BaseRunner
	service *services.PrebuiltService
	client  client.Client
}

func NewPrebuiltRunner(client client.Client, cfg *config.Config) *PrebuiltRunner {
	return &PrebuiltRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		client:     client,
		service:    services.NewPrebuiltService(client),
	}
}

/*
******************************************************************************
Reader interface
******************************************************************************
*/

// Get implements the "get prebuilts ..." command
func (r *PrebuiltRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	prebuilts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: prebuilts,
	}, nil

}

// Describe implements the `describe prebuilt ...` command
func (r *PrebuiltRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	prebuilt, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s", prebuilt.Name),
		Object: prebuilt,
	}, nil
}

/*
******************************************************************************
Writer interface
******************************************************************************
*/

// Create implements the `create prebuilt ...` command
func (r *PrebuiltRunner) Create(in Request) (*Response, error) {
	logger.Trace()
	return notImplemented(in)
}

// Delete implementes the `delete prebuilt ...` command
func (r *PrebuiltRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	options := in.Options.(*flags.PrebuiltDeleteOptions)

	prebuilt, err := r.service.GetByName(name)
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
					if !strings.HasPrefix(err.Error(), "workflow not found") {
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

	if err := r.service.Delete(prebuilt.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted prebuilt `%s` (%s)", name, prebuilt.Id),
	}, nil
}

// Clear implements the `clear prebuilts ...` command
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

	return &Response{
		Text: fmt.Sprintf("Deleted %v prebuilt(s)", cnt),
	}, nil
}

/*
******************************************************************************
Importer interface
******************************************************************************
*/

// Import implements the `import prebuilt ...` command
func (r *PrebuiltRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	wd := filepath.Dir(path)

	if common.Repository != "" {
		defer os.RemoveAll(wd)
	}

	var mPkg map[string]interface{}

	if err := importLoadFromDisk(path, &mPkg); err != nil {
		return nil, err
	}

	var bundles []map[string]interface{}

	for _, ele := range mPkg["bundles"].([]interface{}) {
		item := ele.(map[string]interface{})
		data := item["data"].(map[string]interface{})

		if f, exists := data["filename"]; exists {
			if strings.HasPrefix(f.(string), "@") {
				fp := filepath.Join(wd, f.(string)[1:])

				var b map[string]interface{}

				if err := importLoadFromDisk(fp, &b); err != nil {
					return nil, err
				}

				item = map[string]interface{}{
					"type": item["type"].(string),
					"data": b,
				}
			}
		}

		bundles = append(bundles, item)
	}

	mPkg["bundles"] = bundles

	b, err := json.Marshal(mPkg)
	if err != nil {
		return nil, err
	}

	var pkg services.PrebuiltPackage
	if err := json.Unmarshal(b, &pkg); err != nil {
		return nil, err
	}

	pb, err := r.service.Import(pkg, false)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported prebuilt `%s` (%s)", pb.Name, pb.Id),
	}, nil

}

/*
******************************************************************************
Exporter interface
******************************************************************************
*/

// Export implements the `export prebuilt ...` command
func (r *PrebuiltRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	common := in.Common.(*flags.AssetExportCommon)
	options := in.Options.(*flags.PrebuiltExportOptions)

	pb, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	pkg, err := r.service.Export(pb.Id)
	if err != nil {
		return nil, err
	}

	if options.Expand {
		path := common.Path

		var repo *Repository
		var repoPath string

		if common.Repository != "" {
			repo, err = exportNewRepositoryFromRequest(in)
			if err != nil {
				return nil, err
			}

			var e error

			repoPath, e = repo.Clone(&FileReaderImpl{}, &ClonerImpl{})
			if e != nil {
				return nil, e
			}
			defer os.RemoveAll(repoPath)

			path = filepath.Join(repoPath, common.Path)
		}

		if err := r.expandPrebuilt(pkg, path); err != nil {
			return nil, err
		}

		if common.Repository != "" {
			if err := repo.CommitAndPush(repoPath, common.Message); err != nil {
				return nil, err
			}
		}

	} else {
		fn := fmt.Sprintf("%s.prebuilt.json", strings.Replace(pkg.Metadata.Name, "/", "_", 1))

		if exportAssetFromRequest(in, pkg, fn); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported prebuilt `%s`", name),
	}, nil
}

/*
******************************************************************************
Private functions
******************************************************************************
*/

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

func (r *PrebuiltRunner) expandPrebuilt(pkg *services.PrebuiltPackage, path string) error {
	logger.Trace()

	var bundles []map[string]interface{}

	for _, ele := range pkg.Bundles {
		outputPath := filepath.Join(path, fmt.Sprintf("bundles/%ss", ele.Type))

		if err := utils.EnsurePathExists(outputPath); err != nil {
			return err
		}

		fn := fmt.Sprintf("%s.%s.json", ele.Data["name"], ele.Type)

		if err := utils.WriteJsonToDisk(ele.Data, fn, outputPath); err != nil {
			return err
		}

		filename := fmt.Sprintf("@%s", filepath.Join(fmt.Sprintf("bundles/%ss", ele.Type), fn))

		bundles = append(bundles, map[string]interface{}{
			"type": ele.Type,
			"data": map[string]interface{}{"filename": filename},
		})
	}

	res, err := toMap(pkg)
	if err != nil {
		return err
	}

	res["bundles"] = bundles

	fn := fmt.Sprintf("%s.prebuilt.json", pkg.Metadata.Name)

	if err := utils.WriteJsonToDisk(res, fn, path); err != nil {
		return err
	}

	return nil

}
