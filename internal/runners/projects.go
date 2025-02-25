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

// GetProjects() is the implementation of the command `get projects`
func (r *ProjectRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range projects {
		display = append(display, fmt.Sprintf("%s\t%s", ele.Name, ele.Description))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(projects),
	), nil

}

// Describe implements the `describe project <name>` command
func (r *ProjectRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var project *services.Project
	var err error

	project, err = r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Name: %s (%s)", project.Name, project.Id),
		WithJson(project),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create is the implementation of the command `ccreate project <name>`
func (r *ProjectRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	existing, err := r.service.GetByName(name)
	if err != nil {
		if err.Error() != "project not found" {
			return nil, errors.New(fmt.Sprintf("project `%s` already exists", name))
		}
	}
	if existing != nil {
		return nil, errors.New(fmt.Sprintf("project `%s` already exists", name))
	}

	project, err := r.service.Create(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully created project `%s`", name),
		WithJson(project),
	), nil
}

// Delete is the implementation of the command `delete project <name>`
func (r *ProjectRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	project, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	r.service.Delete(project.Id)

	return &Response{
		Text: fmt.Sprintf("Successfully deleted project `%s`", name),
	}, nil
}

// Clear is the implementation of the command `clear projects`
func (r *ProjectRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	cnt := 0

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range projects {
		r.service.Delete(ele.Id)
		cnt++
	}

	return NewResponse(fmt.Sprintf("Deleted %v project(s)", cnt)), nil
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

	return NewResponse(
		fmt.Sprintf("Successfully copied project `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	), nil
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

	// Handle loading the file from disk
	var project services.Project
	if err := importFile(in, &project); err != nil {
		return nil, err
	}

	common := in.Common.(*flags.AssetImportCommon)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	imported, err := r.importProject(project, path, common.Replace)
	if err != nil {
		return nil, err
	}

	options := in.Options.(*flags.ProjectImportOptions)

	if err := r.updateMembers(imported.Id, options.Members); err != nil {
		// Delete the project
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported project `%s`", project.Name),
	), nil
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
		if err := expandProject(project, common.Path); err != nil {
			return nil, err
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
		if err := utils.WriteJsonToDisk(exported, fn, common.Path); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported project `%s`", project.Name),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Gitter functions
//

// Pull implements the command `pull jsonform <repo>`
func (r *ProjectRunner) Pull(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetPullCommon)

	pull := PullAction{
		Name:    in.Args[1],
		Config:  r.config,
		Options: *common,
	}

	path, err := pull.Clone()
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	file := filepath.Join(path, common.Path, in.Args[0])

	var res services.Project

	if err := importFromPath(file, &res); err != nil {
		return nil, err
	}

	imported, err := r.importProject(res, file, common.Replace)
	if err != nil {
		return nil, err
	}

	options := in.Options.(*flags.ProjectPullOptions)

	if err := r.updateMembers(imported.Id, options.Members); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled project `%s`", res.Name),
	), nil
}

// Push implements the command `push jsonform <repo>`
func (r *ProjectRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	common := in.Common.(*flags.AssetPushCommon)
	options := in.Options.(*flags.ProjectPushOptions)

	project, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	res, err := r.service.Export(project.Id)
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:    in.Args[1],
		Options: *common,
		Config:  r.config,
	}

	if options.Expand {
		path, err := push.Clone()
		if err != nil {
			return nil, err
		}

		if err := expandProject(res, filepath.Join(path, common.Path)); err != nil {
			return nil, err
		}

		if err := push.Commit(path); err != nil {
			return nil, err
		}

	} else {
		b, err := json.Marshal(res)
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

		push.Data = exported
		push.Filename = fmt.Sprintf("%s.project.json", strings.Replace(name, "/", "_", -1))
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed project `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

type Member struct {
	Type   string
	Name   string
	Access string
}

func (r *ProjectRunner) importProject(project services.Project, path string, replace bool) (*services.Project, error) {
	logger.Trace()

	var projectMap map[string]interface{}
	if err := importFromPath(path, &projectMap); err != nil {
		return nil, err
	}

	components := projectMap["components"].([]interface{})

	for idx, ele := range project.Components {
		if val, exists := components[idx].(map[string]interface{})["filename"]; exists {
			basepath := filepath.Dir(path)
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
				r.service.Delete(ele.Id)
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

func expandProject(project *services.Project, path string) error {
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

	return utils.WriteJsonToDisk(projectMap, fn, path)
}
