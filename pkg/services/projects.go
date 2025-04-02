// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type ProjectComponent struct {
	Iid       int                    `json:"iid"`
	Type      string                 `json:"type"`
	Folder    string                 `json:"folder"`
	Reference string                 `json:"reference"`
	Document  map[string]interface{} `json:"document"`
}

type ProjectFolder struct {
	Iid      int             `json:"iid"`
	Name     string          `json:"name"`
	NodeType string          `json:"nodeType"`
	Children []ProjectFolder `json:"children"`
}

type ProjectOperation struct {
	Message  string   `json:"message"`
	Data     Project  `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ProjectMember struct {
	Provenance string `json:"provenance"`
	Reference  string `json:"reference"`
	Role       string `json:"role"`
	Type       string `json:"type"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
}

type ProjectAccessControl struct {
	Manage  []string `json:"manage"`
	Write   []string `json:"write"`
	Execute []string `json:"execute"`
	Read    []string `json:"read"`
}

type Project struct {
	Id              string             `json:"_id"`
	Name            string             `json:"name"`
	BackgroundColor string             `json:"backgroundColor"`
	Components      []ProjectComponent `json:"components"`
	Created         string             `json:"created"`
	CreatedBy       any                `json:"createdBy"`
	Description     string             `json:"description"`
	Folders         []ProjectFolder    `json:"folders"`
	Iid             int                `json:"iid"`
	// Not supported for import
	ComponentIidIndex int    `json:"componentIidIndex"`
	LastUpdated       string `json:"lastUpdated"`
	LastUpdatedBy     any    `json:"lastUpdatedBy"`
	Thumbnail         string `json:"thumbnail,omitempty"`
	// NOt supported for import
	Members []ProjectMember `json:"members"`
	// Not supported for import
	AccessControl ProjectAccessControl `json:"accessControl"`
}

func (p Project) Import() map[string]interface{} {
	logger.Trace()
	return map[string]interface{}{
		"_id":             p.Id,
		"name":            p.Name,
		"backgroundColor": p.BackgroundColor,
		"components":      p.Components,
		"created":         p.Created,
		"createdBy":       p.CreatedBy,
		"description":     p.Description,
		"folders":         p.Folders,
		"iid":             p.Iid,
		"lastUpdated":     p.LastUpdated,
		"lastUpdatedBy":   p.LastUpdatedBy,
		"thumbnail":       p.Thumbnail,
	}
}

type ProjectService struct {
	client *ServiceClient
}

func NewProjectService(c client.Client) *ProjectService {
	return &ProjectService{
		client: NewServiceClient(c),
	}
}

// GetAll will retrieve all of the configured projects from the server and
// return them as an array of Projects.  If there are no configured projects on
// the server, this function will return an empty array.
func (svc *ProjectService) GetAll() ([]Project, error) {
	logger.Trace()

	type Response struct {
		Message  string    `json:"message"`
		Data     []Project `json:"data"`
		Metadata Metadata  `json:"metadata"`
	}

	var res Response
	var projects []Project

	var limit = 100
	var skip = 0

	for {
		if err := svc.client.GetRequest(&Request{
			uri:    "/automation-studio/projects",
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.Data {
			projects = append(projects, ele)
		}

		if len(projects) == res.Metadata.Total {
			break
		}

		skip += limit
	}

	logger.Info("Found %v project(s)", len(projects))

	return projects, nil
}

// Get will attempt to retreive the project by its identifier.  If the project
// does not exist, this function will return an error.
func (svc *ProjectService) Get(id string) (*Project, error) {
	logger.Trace()

	type Response struct {
		Message  string   `json:"message"`
		Data     *Project `json:"data"`
		Metadata Metadata `json:"metadata"`
	}

	var res Response

	var uri = fmt.Sprintf("/automation-studio/projects/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

// GetByName will attempt to get a project wiht the specified name.  If the
// project does not exist, this function will return an error
func (svc *ProjectService) GetByName(name string) (*Project, error) {
	logger.Trace()

	projects, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var res *Project

	for _, ele := range projects {
		if ele.Name == name {
			res = &ele
			break
		}
	}

	if res == nil {
		return nil, errors.New("project not found")
	}

	return res, nil
}

// Create implement `http.MethodPost /automation-studio/projects`
func (svc *ProjectService) Create(name string) (*Project, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":       name,
		"components": []string{},
	}

	type Response struct {
		Message  string   `json:"message"`
		Data     *Project `json:"data"`
		Metadata any      `json:"metadata"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/automation-studio/projects",
		body:               body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

// Delete implements `http.MethodDelete /automation-studio/projects/{id}`
func (svc *ProjectService) Delete(id string) error {
	logger.Trace()
	return svc.client.Delete(
		fmt.Sprintf("/automation-studio/projects/%s", id),
	)
}

// This function will recusively iterate over folders in a project schema and
// remove keys in order for the body to be accepted by the server
func (svc *ProjectService) transformImport(in map[string]interface{}) {
	if in["nodeType"].(string) == "folder" {
		delete(in, "iid")
	}

	if in["nodeType"].(string) == "component" {
		delete(in, "name")
	}

	if in["children"] != nil {
		for _, ele := range in["children"].([]interface{}) {
			svc.transformImport(ele.(map[string]interface{}))
		}
	} else if in["children"] == nil {
		delete(in, "children")
	}
}

// Import implements `http.MethodPost /automation-studio/projects/import`
func (svc *ProjectService) Import(in Project) (*Project, error) {
	logger.Trace()

	body := map[string]interface{}{
		"conflictMode": "insert-new",
		"project":      in.Import(),
	}

	b, _ := json.Marshal(body)

	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	project := data["project"].(map[string]interface{})
	folders := project["folders"].([]interface{})

	if folders != nil {
		for _, ele := range folders {
			svc.transformImport(ele.(map[string]interface{}))
		}
	}

	type Response struct {
		Message  string                 `json:"message"`
		Data     *Project               `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/automation-studio/projects/import",
		body:               data,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

// Export implements `http.MethodGet /automation-studio/projects/{id}/export`
func (svc *ProjectService) Export(id string) (*Project, error) {
	logger.Trace()

	type Response struct {
		Message  string   `json:"message"`
		Data     *Project `json:"data"`
		Metadata Metadata `json:"metadata"`
	}

	var res Response
	var uri = fmt.Sprintf("/automation-studio/projects/%s/export", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

func (svc *ProjectService) AddMembers(projectId string, members []ProjectMember) error {
	logger.Trace()

	project, err := svc.Get(projectId)
	if err != nil {
		return err
	}

	for _, ele := range project.Members {
		members = append(members, ele)
	}

	body := map[string]interface{}{
		"members": members,
	}

	uri := fmt.Sprintf("/automation-studio/projects/%s", projectId)

	return svc.client.Patch(uri, body, nil)
}
