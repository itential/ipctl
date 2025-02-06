// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/repositories"
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

	project, err = r.service.Get(name)
	if err != nil {
		project, err = r.GetByName(name)
		if err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Name: %s", project.Name),
		WithJson(project),
	), nil
}

// Create is the implementation of the command `ccreate project <name>`
func (r *ProjectRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	_, err := r.GetByName(name)
	if err == nil {
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

	project, err := r.GetByName(name)
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

// Copy implements the `copy project <name> <dst>` command
func (r *ProjectRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common *flags.AssetCopyCommon
	utils.LoadObject(in.Common, &common)

	var options *flags.ProjectImportOptions
	utils.LoadObject(in.Options, &options)

	if common.From == common.To {
		return nil, errors.New("source and destination servers must be different values")
	}

	from, cancel, err := newClient(common.From, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	fromService := services.NewProjectService(from)

	fromProjects, err := fromService.GetAll()
	if err != nil {
		return nil, err
	}

	var projectId string

	for _, ele := range fromProjects {
		if ele.Name == name {
			projectId = ele.Id
		}
	}

	if projectId == "" {
		return nil, errors.New(fmt.Sprintf("could not find project named `%s`", name))
	}

	project, err := fromService.Export(projectId)
	if err != nil {
		return nil, err
	}

	to, cancel, err := newClient(common.To, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	toService := services.NewProjectService(to)

	importedProject, err := toService.Import(*project)
	if err != nil {
		return nil, err
	}

	accounts := services.NewAccountService(to)
	groups := services.NewGroupService(to)
	userSettings := services.NewUserSettingsService(to)

	var members []services.ProjectMember

	activeUser, err := userSettings.Get()
	if err != nil {
		return nil, err
	}

	for _, ele := range options.Members {
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
		if err := toService.AddMembers(importedProject.Id, members); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully copied project `%s` from `%s` to `%s`", name, common.From, common.To),
	), nil
}

// Import implements the command `import project <path>`
func (r *ProjectRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	var options flags.ProjectImportOptions
	utils.LoadObject(in.Options, &options)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	data, err := utils.ReadFromFile(path)
	if err != nil {
		return nil, err
	}

	///////
	var projectMap map[string]interface{}
	utils.UnmarshalData(data, &projectMap)

	var project services.Project
	utils.UnmarshalData(data, &project)

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
			if common.Replace {
				r.service.Delete(ele.Id)
			} else {
				return nil, errors.New(fmt.Sprintf("project `%s` already exists", project.Name))
			}
		}
	}

	imported, err := r.service.Import(project)
	if err != nil {
		return nil, err
	}

	var members []services.ProjectMember

	activeUser, err := r.userSettings.Get()
	if err != nil {
		return nil, err
	}

	for _, ele := range options.Members {
		member, err := parseMember(ele)
		if err != nil {
			return nil, err
		}

		logger.Info("checking member %v", member)

		if !(member.Type == "account" && member.Name == activeUser.Username) {
			if member.Type == "account" {
				account, err := r.accounts.GetByName(member.Name)
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
				group, err := r.groups.GetByName(member.Name)
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
		if err := r.service.AddMembers(imported.Id, members); err != nil {
			return nil, err
		}
	}
	//////

	return NewResponse(
		fmt.Sprintf("Successfully imported project `%s`", imported.Name),
	), nil
}

// Export is the implementation of the command `export project <name>`
func (r *ProjectRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	var options flags.ProjectExportOptions
	utils.LoadObject(in.Options, &options)

	if common.Path != "" {
		utils.EnsurePathExists(common.Path)
	}

	name := in.Args[0]

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var id string
	for _, ele := range projects {
		if ele.Name == name {
			id = ele.Id
		}
	}

	if id == "" {
		return nil, errors.New(fmt.Sprintf("Unable to find project with name %s", name))
	}

	project, err := r.service.Export(id)
	if err != nil {
		return nil, err
	}

	if options.Expand {
		for _, ele := range project.Folders {
			if ele.NodeType == "folder" {
				makeFolder(common.Path, ele)
			}
		}

		var projectMap map[string]interface{}
		if err := utils.ToMap(project, &projectMap); err != nil {
			return nil, err
		}

		components := projectMap["components"].([]interface{})

		for idx, ele := range project.Components {
			p := common.Path
			if ele.Folder != "/" {
				p = filepath.Join(common.Path, ele.Folder[1:len(ele.Folder)])
			}

			docName := strings.Replace(ele.Document["name"].(string), "/", "_", -1)
			fn := fmt.Sprintf("%s.%s.json", docName, strings.ToLower(ele.Type))
			if err := utils.WriteJsonToDisk(ele.Document, fn, p); err != nil {
				return nil, err
			}

			delete(components[idx].(map[string]interface{}), "document")
			components[idx].(map[string]interface{})["filename"] = fn
		}

		fn := fmt.Sprintf("%s.project.json", strings.Replace(name, "/", "_", -1))
		if err := utils.WriteJsonToDisk(projectMap, fn, common.Path); err != nil {
			return nil, err
		}
	} else {
		fn := fmt.Sprintf("%s.project.json", strings.Replace(name, "/", "_", -1))
		if err := utils.WriteJsonToDisk(project, fn, common.Path); err != nil {
			return nil, err
		}
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported project `%s`", project.Name),
	), nil
}

// This function will attempt to get a project wiht the specified name.  If the
// project does not exist, this function will return an error
func (r *ProjectRunner) GetByName(name string) (*services.Project, error) {
	logger.Trace()

	projects, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var projectId string
	for _, ele := range projects {
		if ele.Name == name {
			projectId = ele.Id
		}
	}

	if projectId == "" {
		return nil, errors.New(fmt.Sprintf("project `%s` does not exist", name))
	}

	project, err := r.service.Get(projectId)
	if err != nil {
		return nil, err
	}

	return project, nil
}

type Member struct {
	Type   string
	Name   string
	Access string
}

func cloneRepository(in *Repository) (string, error) {
	logger.Trace()

	r := repositories.Repository{
		Url: in.Url,
	}
	if in.PrivateKeyFile != "" {
		pk, err := utils.ReadStringFromFile(in.PrivateKeyFile)
		if err != nil {
			return "", err
		}
		r.PrivateKey = pk
	}
	if in.Reference != "" {
		r.Reference = in.Reference
	}
	p, err := r.Clone()
	if err != nil {
		return "", err
	}

	return p, nil
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

func newClient(name string, cfg *config.Config) (client.Client, context.CancelFunc, error) {
	logger.Trace()

	profile, err := cfg.GetProfile(name)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(profile.Timeout)*time.Second,
	)

	return client.New(ctx, profile), cancel, nil
}
