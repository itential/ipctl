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

type ProjectRunner struct {
	config       *config.Config
	service      *services.ProjectService
	accounts     *services.AccountService
	groups       *services.GroupService
	userSettings *services.UserSettingsService
}

func NewProjectRunner(client client.Client, cfg *config.Config) *ProjectRunner {
	return &ProjectRunner{
		config:       cfg,
		service:      services.NewProjectService(client),
		accounts:     services.NewAccountService(client),
		groups:       services.NewGroupService(client),
		userSettings: services.NewUserSettingsService(client),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get projects ...` command
func (r *ProjectRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: projects,
	}, nil

}

// Describe implements the `describe project ...` command
func (r *ProjectRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	var res *services.Project

	res, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	createdBy := res.CreatedBy.(map[string]interface{})["username"].(string)
	updatedBy := res.LastUpdatedBy.(map[string]interface{})["username"].(string)

	output := []string{
		fmt.Sprintf("Name: %s (%s)", res.Name, res.Id),
		fmt.Sprintf("Description: %s", res.Description),
		fmt.Sprintf("Created: %s, by: %s", res.Created, createdBy),
		fmt.Sprintf("Updated: %s, by: %s", res.LastUpdated, updatedBy),
	}

	return &Response{
		Text:   strings.Join(output, "\n"),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create is the implementation of the command `ccreate project <name>`
func (r *ProjectRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	existing, err := r.service.GetByName(name)
	if existing != nil {
		return nil, errors.New(fmt.Sprintf("project `%s` already exists", name))
	}

	project, err := r.service.Create(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created project `%s` (%s)", project.Name, project.Id),
		Object: project,
	}, nil
}

// Delete is the implementation of the command `delete project <name>`
func (r *ProjectRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	project, err := r.service.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(project.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted project `%s` (%s)", project.Name, project.Id),
	}, nil
}

// Clear is the implementation of the command `clear projects`
func (r *ProjectRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range projects {
		if err := r.service.Delete(ele.Id); err != nil {
			logger.Debug("failed to delete project `%s` (%s)", ele.Name, ele.Id)
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v project(s)", len(projects)),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

// Copy implements the `copy project <name> <dst>` command
func (r *ProjectRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "project"}, r)
	if err != nil {
		return nil, err
	}

	client, cancel, err := NewClient(in.Common.(*flags.AssetCopyCommon).To, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	projectsvc := services.NewProjectService(client)
	accounts := services.NewAccountService(client)
	groups := services.NewGroupService(client)
	userSettings := services.NewUserSettingsService(client)

	var members []services.ProjectMember

	activeUser, err := userSettings.Get()
	if err != nil {
		return nil, err
	}

	for _, ele := range in.Options.(*flags.ProjectCopyOptions).Members {
		member, err := parseMember(ele)
		if err != nil {
			return nil, err
		}

		if !(member.Type == "account" && member.Name == activeUser.Username) {
			if member.Type == "account" {
				account, err := accounts.GetByName(member.Name)
				if err != nil {
					return nil, err
				}

				members = append(members, services.ProjectMember{
					Provenance: account.Provenance,
					Reference:  account.Id,
					Role:       member.Access,
					Type:       "account",
					Username:   account.Username,
				})
			} else {
				group, err := groups.GetByName(member.Name)
				if err != nil {
					return nil, err
				}

				members = append(members, services.ProjectMember{
					Provenance: group.Provenance,
					Reference:  group.Id,
					Role:       member.Access,
					Type:       "group",
					Name:       group.Name,
				})
			}
		} else {
			logger.Info("skipping %s", member.Name)
		}
	}

	if len(members) > 0 {
		if err := projectsvc.AddMembers(res.CopyToData.(services.Project).Id, members); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied project `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *ProjectRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewProjectService(client)

	project, err := svc.GetByName(name)
	if err != nil {
		return nil, err
	}

	res, err := svc.Export(project.Id)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *ProjectRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewProjectService(client)

	name := in.(services.Project).Name

	if exists, err := svc.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("project `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
	}

	return svc.Import(in.(services.Project))
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the command `import project <path>`
func (r *ProjectRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)
	options := in.Options.(*flags.ProjectImportOptions)

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	wd := filepath.Dir(path)

	if common.Repository != "" {
		defer os.RemoveAll(wd)
	}

	var project services.Project

	if err := importLoadFromDisk(path, &project); err != nil {
		return nil, err
	}

	imported, err := r.importProject(project, path, common.Replace)
	if err != nil {
		return nil, err
	}

	if err := r.updateMembers(imported.Id, options.Members); err != nil {
		// Delete the project
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported project `%s` (%s)", project.Name, project.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export is the implementation of the command `export project <name>`
func (r *ProjectRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetExportCommon)
	options := in.Options.(*flags.ProjectExportOptions)

	name := in.Args[0]

	p, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	project, err := r.service.Export(p.Id)
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

		if err := expandProject(in, project, path); err != nil {
			return nil, err
		}

		if common.Repository != "" {
			logger.Info("commting %s", repoPath)
			if err := repo.CommitAndPush(repoPath, common.Message); err != nil {
				return nil, err
			}
		}

	} else {
		b, err := json.Marshal(project)
		if err != nil {
			return nil, err
		}

		var exported map[string]interface{}
		if err := json.Unmarshal(b, &exported); err != nil {
			return nil, err
		}

		// NOTE (privateip) need to remove these fields to comply with the same
		// export format returned from the UI
		delete(exported, "members")
		delete(exported, "accessControl")

		fn := fmt.Sprintf("%s.project.json", strings.Replace(name, "/", "_", -1))

		if err := exportAssetFromRequest(in, exported, fn); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported project `%s`", project.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

type Member struct {
	Type   string
	Name   string
	Access string
}

// This function will attempt to import the project to the server.  It is
// responsible for reconstructing the project file if the project was exported
// using the `--expand` command line option.
func (r *ProjectRunner) importProject(project services.Project, path string, replace bool) (*services.Project, error) {
	logger.Trace()

	var projectMap map[string]interface{}
	if err := importLoadFromDisk(path, &projectMap); err != nil {
		return nil, err
	}

	components := projectMap["components"].([]interface{})
	basepath := filepath.Dir(path)

	for idx, ele := range project.Components {
		if val, exists := components[idx].(map[string]interface{})["filename"]; exists {
			fp := filepath.Join(basepath, ele.Folder[1:len(ele.Folder)], val.(string))

			doc, err := os.ReadFile(fp)
			if err != nil {
				return nil, err
			}

			var document map[string]interface{}
			utils.UnmarshalData(doc, &document)

			project.Components[idx].Document = document
		}
	}

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range projects {
		if ele.Name == project.Name {
			if replace {
				if err := r.service.Delete(ele.Id); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New(fmt.Sprintf("project `%s` already exists, use `--replace` to overwrite it", project.Name))
			}
		}
	}

	return r.service.Import(project)
}

func (r *ProjectRunner) updateMembers(projectId string, projectMembers []string) error {
	logger.Trace()

	var members []services.ProjectMember

	activeUser, err := r.userSettings.Get()
	if err != nil {
		return err
	}

	for _, ele := range projectMembers {
		member, err := parseMember(ele)
		if err != nil {
			return err
		}

		logger.Info("checking member %v", member)

		if !(member.Type == "account" && member.Name == activeUser.Username) {
			if member.Type == "account" {
				account, err := r.accounts.GetByName(member.Name)
				if err != nil {
					return err
				}

				members = append(members, services.ProjectMember{
					Provenance: account.Provenance,
					Reference:  account.Id,
					Role:       member.Access,
					Type:       "account",
					Username:   account.Username,
				})
			} else {
				group, err := r.groups.GetByName(member.Name)
				if err != nil {
					return err
				}

				members = append(members, services.ProjectMember{
					Provenance: group.Provenance,
					Reference:  group.Id,
					Role:       member.Access,
					Type:       "group",
					Name:       group.Name,
				})
			}
		} else {
			logger.Info("skipping %s", member.Name)
		}
	}

	if len(members) > 0 {
		if err := r.service.AddMembers(projectId, members); err != nil {
			return err
		}
	}

	return nil

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func makeFolder(p string, f services.ProjectFolder) {
	path := filepath.Join(p, f.Name)
	if !utils.PathExists(path) {
		utils.EnsurePathExists(path)
	}
	for _, ele := range f.Children {
		makeFolder(path, ele)
	}
}

func parseMember(member string) (*Member, error) {
	parts := strings.Split(member, ",")

	m := &Member{Access: "editor"}

	for _, ele := range parts {
		tokens := strings.Split(ele, "=")
		switch tokens[0] {
		case "type":
			m.Type = tokens[1]
		case "name":
			m.Name = tokens[1]
		case "access":
			m.Access = tokens[1]
		}
	}

	if !stringInSlice(m.Access, []string{"owner", "editor", "operator", "viewer"}) {
		return nil, errors.New(fmt.Sprintf("invalid value for access: `%s`", m.Access))
	}

	return m, nil
}

func expandProject(in Request, project *services.Project, path string) error {
	logger.Trace()

	for _, ele := range project.Folders {
		if ele.NodeType == "folder" {
			makeFolder(path, ele)
		}
	}

	var projectMap map[string]interface{}
	if err := utils.ToMap(project, &projectMap); err != nil {
		return err
	}

	// NOTE (privateip) need to remove these fields to comply with the same
	// export format returned from the UI
	delete(projectMap, "members")
	delete(projectMap, "accessControl")

	components := projectMap["components"].([]interface{})

	for idx, ele := range project.Components {
		p := path
		if ele.Folder != "/" {
			p = filepath.Join(path, ele.Folder[1:len(ele.Folder)])
		}

		docName := strings.Replace(ele.Document["name"].(string), "/", "_", -1)
		fn := fmt.Sprintf("%s.%s.json", docName, strings.ToLower(ele.Type))
		if err := utils.WriteJsonToDisk(ele.Document, fn, p); err != nil {
			return err
		}

		delete(components[idx].(map[string]interface{}), "document")
		components[idx].(map[string]interface{})["filename"] = fn
	}

	fn := fmt.Sprintf("%s.project.json", strings.Replace(project.Name, "/", "_", -1))

	return utils.WriteJsonToDisk(project, fn, path)
}
