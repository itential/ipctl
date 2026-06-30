// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

const (
	agentProjectsBasePath    = "/agent-project-service/projects"
	agentProjectBundlesPath  = "/agent-project-service/project-bundles"
	defaultAgentProjectLimit = 100
)

// AgentProjectComponent represents a single component (agent) within an agent project.
type AgentProjectComponent struct {
	Type      string `json:"type"`
	Iid       int    `json:"iid"`
	Reference string `json:"reference"`
	Folder    string `json:"folder"`
}

// AgentProjectMember represents a user or group member of an agent project.
type AgentProjectMember struct {
	Type       string `json:"type"`
	Role       string `json:"role"`
	Reference  string `json:"reference"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
	Provenance string `json:"provenance,omitempty"`
}

// AgentProject represents an agent project in the agent-project-service.
type AgentProject struct {
	Id          string                  `json:"_id"`
	Iid         int                     `json:"iid"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Components  []AgentProjectComponent `json:"components"`
	Created     string                  `json:"created"`
	CreatedBy   any                     `json:"createdBy"`
	LastUpdated string                  `json:"lastUpdated"`
	LastUpdatedBy any                   `json:"lastUpdatedBy"`
	Members     []AgentProjectMember    `json:"members"`
}

// AgentProjectBundle represents the export bundle for an agent project.
// It contains the project metadata along with all agent definitions.
type AgentProjectBundle struct {
	Id                       string                 `json:"_id"`
	Name                     string                 `json:"name"`
	Description              string                 `json:"description"`
	AgentProjectBundleVersion int                   `json:"agentProjectBundleVersion"`
	Created                  string                 `json:"created,omitempty"`
	CreatedBy                any                    `json:"createdBy,omitempty"`
	Agents                   []map[string]interface{} `json:"agents"`
}

// AgentProjectService provides operations for managing agent projects.
type AgentProjectService struct {
	BaseService
}

// NewAgentProjectService creates a new AgentProjectService with the provided client.
func NewAgentProjectService(c client.Client) *AgentProjectService {
	return &AgentProjectService{BaseService: NewBaseService(c)}
}

type getAgentProjectsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Items []AgentProject `json:"items"`
		Total int            `json:"total"`
		Skip  int            `json:"skip"`
		Limit int            `json:"limit"`
	} `json:"data"`
}

type agentProjectResponse struct {
	Message string       `json:"message"`
	Data    AgentProject `json:"data"`
}

type agentProjectBundleResponse struct {
	Message string             `json:"message"`
	Data    AgentProjectBundle `json:"data"`
}

// GetAll retrieves all agent projects, handling pagination automatically.
func (svc *AgentProjectService) GetAll() ([]AgentProject, error) {
	logging.Trace()

	var res getAgentProjectsResponse
	projects := make([]AgentProject, 0, defaultAgentProjectLimit)
	limit := defaultAgentProjectLimit
	skip := 0

	for {
		if err := svc.GetRequest(&Request{
			uri:    agentProjectsBasePath,
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, fmt.Errorf("failed to retrieve agent projects (skip=%d, limit=%d): %w", skip, limit, err)
		}

		projects = append(projects, res.Data.Items...)

		if len(projects) >= res.Data.Total {
			break
		}

		skip += limit
	}

	logging.Info("Found %d agent project(s)", len(projects))

	return projects, nil
}

// Get retrieves a single agent project by its ID.
func (svc *AgentProjectService) Get(id string) (*AgentProject, error) {
	logging.Trace()

	if id == "" {
		return nil, fmt.Errorf("agent project id cannot be empty")
	}

	var res agentProjectResponse

	uri := fmt.Sprintf("%s/%s", agentProjectsBasePath, id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, fmt.Errorf("failed to get agent project %s: %w", id, err)
	}

	return &res.Data, nil
}

// GetByName retrieves an agent project by name using client-side filtering.
func (svc *AgentProjectService) GetByName(name string) (*AgentProject, error) {
	logging.Trace()

	projects, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if projects[i].Name == name {
			return &projects[i], nil
		}
	}

	return nil, errors.New("agent project not found")
}

// Export exports an agent project bundle by project ID.
// The bundle contains the project metadata and all agent definitions.
func (svc *AgentProjectService) Export(id string) (*AgentProjectBundle, error) {
	logging.Trace()

	if id == "" {
		return nil, fmt.Errorf("agent project id cannot be empty")
	}

	var res agentProjectBundleResponse

	uri := fmt.Sprintf("%s/%s/export", agentProjectBundlesPath, id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, fmt.Errorf("failed to export agent project %s: %w", id, err)
	}

	return &res.Data, nil
}

// Import imports an agent project bundle into the platform.
func (svc *AgentProjectService) Import(bundle AgentProjectBundle) (*AgentProjectBundle, error) {
	logging.Trace()

	type importResponse struct {
		Message string             `json:"message"`
		Data    AgentProjectBundle `json:"data"`
	}

	body := map[string]interface{}{
		"bundle": bundle,
	}

	var res importResponse

	uri := fmt.Sprintf("%s/import", agentProjectBundlesPath)

	if err := svc.PostRequest(&Request{
		uri:                uri,
		body:               body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, fmt.Errorf("failed to import agent project: %w", err)
	}

	return &res.Data, nil
}
